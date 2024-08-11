package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/cookie"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/apps/account"
	"github.com/koo-arch/adjusta-backend/internal/apps/user"
	"github.com/koo-arch/adjusta-backend/internal/auth"
	"github.com/koo-arch/adjusta-backend/internal/google/oauth"
	"github.com/koo-arch/adjusta-backend/internal/google/userinfo"
	"golang.org/x/oauth2"
)

func GoogleLoginHandler(c *gin.Context) {
	url := oauth.GetGoogleAuthConfig().AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func AddAccountHandler(c *gin.Context) {
	url := oauth.GetAddAccountAuthConfig().AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func LogoutHandler(c *gin.Context) {
	// セッションをクリア
	session := sessions.Default(c)
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1, Path: "/"})
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save session"})
		return
	}

	// クッキーを削除
	cookie.DeleteCookie(c, "access_token")
	cookie.DeleteCookie(c, "refresh_token")

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func GoogleCallbackHandler(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		// クエリパラメータからcodeを取得
		code := c.Query("code")
		if code == "" {
			fmt.Println("missing code")
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing code"})
			return
		}

		ctx := c.Request.Context()

		// Googleからトークンを取得
		oauthToken, err := oauth.GetGoogleAuthConfig().Exchange(ctx, code)
		if err != nil {
			fmt.Println("failed to exchange oauthToken")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to exchange oauthToken"})
			return
		}

		// Googleからユーザー情報を取得
		userInfo, err := userinfo.FetchGoogleUserInfo(ctx, oauthToken)
		if err != nil {
			fmt.Println("failed to fetch user info")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user info"})
			return
		}

		// アプリ独自のJWTトークンを生成
		jwtManager := auth.NewJWTManager(client, auth.NewKeyManager(client))
		jwtToken, err := jwtManager.GenerateTokens(ctx, client, userInfo.Email)

		userRepo := user.NewUserRepository(client)
		accountRepo := account.NewAccountRepository(client)
		authManager := auth.NewAuthManager(client, userRepo, accountRepo)

		// ユーザー情報をデータベースに保存
		u, err := authManager.ProcessUserSignIn(ctx, userInfo, jwtToken, oauthToken)
		if err != nil {
			fmt.Printf("failed to create or update user: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create or update user"})
			return
		}

		// アカウントのoauthトークンを検証
		accounts, err := accountRepo.FilterByUserID(ctx, nil, u.ID)
		if err != nil {
			fmt.Printf("failed to get accounts: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get accounts"})
			return
		}

		for _, account := range accounts {
			_, err := authManager.VerifyOAuthToken(ctx, u.ID, account.Email)
			if err != nil {
				fmt.Printf("failed to reuse token source: %s", err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reuse token source"})
				return
			}
		}

		// クッキーにトークンをセット
		maxAge := int(jwtToken.RefreshExpiration.Sub(time.Now()).Seconds())
		cookie.SetCookie(c, "access_token", jwtToken.AccessToken, maxAge)

		// セッションにユーザー情報をセット
		session.Set("googleid", userInfo.GoogleID)
		session.Set("userid", u.ID.String())
		if err := session.Save(); err != nil {
			fmt.Printf("failed to save session:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save session"})
			return
		}

		// リダイレクト
		c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000")
	}
}

func AddAccountCallbackHandler(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		// クエリパラメータからcodeを取得
		code := c.Query("code")
		if code == "" {
			fmt.Println("missing code")
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing code"})
			return
		}

		ctx := c.Request.Context()

		// Googleからトークンを取得
		oauthToken, err := oauth.GetAddAccountAuthConfig().Exchange(ctx, code)
		if err != nil {
			fmt.Println("failed to exchange oauthToken")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to exchange oauthToken"})
			return
		}

		// Googleからユーザー情報を取得
		userInfo, err := userinfo.FetchGoogleUserInfo(ctx, oauthToken)
		if err != nil {
			fmt.Println("failed to fetch user info")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user info"})
			return
		}

		userRepo := user.NewUserRepository(client)
		accountRepo := account.NewAccountRepository(client)
		authManager := auth.NewAuthManager(client, userRepo, accountRepo)

		// 現在のユーザーにアカウントを追加
		useridStr, ok := session.Get("userid").(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to get userid from session"})
			c.Abort()
			return
		}

		userid, err := uuid.Parse(useridStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userid format"})
			c.Abort()
			return
		}

		_, err = authManager.AddAccountToUser(ctx, userid, userInfo, oauthToken)
		if err != nil {
			fmt.Printf("failed to add account: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add account"})
			return
		}

		// リダイレクト
		c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000")
	}
}