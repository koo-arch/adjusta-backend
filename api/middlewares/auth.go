package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/apps/user"
	"github.com/koo-arch/adjusta-backend/internal/auth"
	"github.com/koo-arch/adjusta-backend/internal/models"
)

func AuthMiddleware(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization header"})
			c.Abort()
			return
		}

		ctx := c.Request.Context()

		// トークンの有効性を確認
		accessToken := tokenStr
		jwtManager := auth.NewJWTManager(client, auth.NewKeyManager(client))
		email, err := jwtManager.VerifyToken(ctx, client, accessToken, "access")
		if err != nil{
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
					c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to verify token"})
					c.Abort()
					return
				}
			} else {
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
	session := sessions.Default(c)

	userID, ok := session.Get("userid").(int)
	if !ok {
		return nil, fmt.Errorf("session not found")
	}
	userRepo := user.NewUserRepository(client)

	ctx := c.Request.Context()
	// ユーザーの検索
	u, err := userRepo.Read(ctx, nil, userID)
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

	// ユーザーのリフレッシュトークンを更新
	_, err = userRepo.Update(ctx, nil, userID, token)
	if err != nil {
		return nil, err
	}

	c.SetCookie("access_token", token.AccessToken, 256, "/", "", false, true)

	return token, nil
}