package handlers

import (
	"net/http"
	
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/google/userinfo"
)

type AccountHandler struct {
	handler *Handler
}

func NewAccountHandler(handler *Handler) *AccountHandler {
	return &AccountHandler{handler: handler}
}

func (ah *AccountHandler) FetchAccountsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

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

		authManager := ah.handler.Server.AuthManager

		token, err := authManager.VerifyOAuthToken(ctx, userid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify token"})
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
