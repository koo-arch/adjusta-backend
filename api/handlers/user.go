package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/apps/account"
	"github.com/koo-arch/adjusta-backend/internal/auth"
)

func GetCurrentUserHandler(client *ent.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		email, ok := c.Get("email")
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "missing email"})
			c.Abort()
			return
		}

		ctx := c.Request.Context()
		
		accountRepo := account.NewAccountRepository(client)
		
		token, err := auth.VerifyOAuthToken(ctx, accountRepo, email.(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to verify token"})
			c.Abort()
			return
		}

		userInfo, err := auth.FetchGoogleUserInfo(ctx, token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user info"})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, userInfo)
	}
}