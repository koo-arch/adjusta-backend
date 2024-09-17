package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/auth"
	"github.com/koo-arch/adjusta-backend/internal/google/userinfo"
	"github.com/koo-arch/adjusta-backend/internal/repo/account"
	"github.com/koo-arch/adjusta-backend/internal/repo/user"
)

func GetCurrentUserHandler(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		email, ok := c.Get("email")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to get email from context"})
			c.Abort()
			return
		}

		session := sessions.Default(c)
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

		ctx := c.Request.Context()

		userRepo := user.NewUserRepository(client)
		accountRepo := account.NewAccountRepository(client)
		authManager := auth.NewAuthManager(client, userRepo, accountRepo)

		token, err := authManager.VerifyOAuthToken(ctx, userid, email.(string))
		if err != nil {
			println("oauth期限切れ")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to verify token"})
			c.Abort()
			return
		}

		userInfo, err := userinfo.FetchGoogleUserInfo(ctx, token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user info"})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, userInfo)
	}
}
