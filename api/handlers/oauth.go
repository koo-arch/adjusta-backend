package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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

func LogoutHandler(c *gin.Context) {
	// セッションをクリア
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save session"})
		return
	}

	// クッキーを削除
	c.SetCookie("access_token", "", -1, "/", "", false, true)

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

		// クッキーにトークンをセット
		c.SetCookie("access_token", jwtToken.AccessToken, 256, "/", "", false, true)

		// セッションにユーザー情報をセット
		session.Set("googleid", userInfo.GoogleID)
		session.Set("userid", u.ID)
		if err := session.Save(); err != nil {
			fmt.Println("failed to save session")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save session"})
			return
		}

		// リダイレクト
		c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000")
	}
}
