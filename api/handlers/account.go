package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/google/userinfo"
)

type AccountsInfo struct {
	AccountID string             `json:"account_id"`
	UserInfo  *userinfo.UserInfo `json:"user_info"`
}

func (s *Server) FetchAccountsHandler(client *ent.Client) gin.HandlerFunc {
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

		userAccounts, err := s.accountRepo.FilterByUserID(ctx, nil, userid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user accounts"})
			c.Abort()
			return
		}

		var accountsInfo []AccountsInfo

		for _, userAccount := range userAccounts {
			token, err := s.authManager.VerifyOAuthToken(ctx, userid, userAccount.Email)
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

			accountsInfo = append(accountsInfo, AccountsInfo{
				AccountID: userAccount.ID.String(),
				UserInfo:  userInfo,
			})

		}

		c.JSON(http.StatusOK, accountsInfo)
	}
}
