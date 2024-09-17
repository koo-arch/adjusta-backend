package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/cookie"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/auth"
	"github.com/koo-arch/adjusta-backend/internal/models"
	"github.com/koo-arch/adjusta-backend/internal/repo/user"
)

func AuthMiddleware(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie("access_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to get access token"})
			c.Abort()
			return
		}

		session := sessions.Default(c)
		googleid := session.Get("googleid").(string)
		println(googleid)

		ctx := c.Request.Context()

		// トークンの有効性を確認
		jwtManager := auth.NewJWTManager(client, auth.NewKeyManager(client))
		email, err := jwtManager.VerifyToken(ctx, client, accessToken, "access")
		if err != nil {
			// トークンの有効期限が切れている場合はリフレッシュトークンを利用してトークンを再発行
			if strings.Contains(err.Error(), "token is expired") {
				token, err := tokenRefresh(c, client, jwtManager)
				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to refresh token"})
					c.Abort()
					return
				}

				accessToken = token.AccessToken
				email, err = jwtManager.VerifyToken(ctx, client, accessToken, "access")
				if err != nil {
					println("アクセストークンの再発行に失敗")
					c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to verify token"})
					c.Abort()
					return
				}
			} else {
				println(err.Error())
				println("アクセストークンの有効性確認に失敗")
				c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to verify token"})
				c.Abort()
				return
			}
		}

		c.Set("email", email)
		c.Set("access_token", accessToken)
		c.Next()
	}
}

func tokenRefresh(c *gin.Context, client *ent.Client, jwtManager *auth.JWTManager) (*models.JWTToken, error) {
	ctx := c.Request.Context()
	session := sessions.Default(c)
	useridStr, ok := session.Get("userid").(string)
	if !ok {
		return nil, fmt.Errorf("failed to get userid")
	}

	userid, err := uuid.Parse(useridStr)
	if err != nil {
		return nil, err
	}

	// リフレッシュトークンの取得
	userRepo := user.NewUserRepository(client)
	u, err := userRepo.Read(ctx, nil, userid)
	if err != nil {
		return nil, err
	}

	// リフレッシュトークンの有効性を確認
	email, err := jwtManager.VerifyToken(ctx, client, u.RefreshToken, "refresh")
	if err != nil {
		return nil, err
	}

	// トークンの再発行
	token, err := jwtManager.GenerateTokens(ctx, client, email)
	if err != nil {
		return nil, err
	}

	maxAge := int(token.RefreshExpiration.Sub(time.Now()).Seconds())
	cookie.SetCookie(c, "access_token", token.AccessToken, maxAge)

	// リフレッシュトークンの更新
	_, err = userRepo.Update(ctx, nil, userid, token)
	if err != nil {
		return nil, err
	}

	return token, nil
}
