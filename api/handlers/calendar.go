package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/apps/account"
	"github.com/koo-arch/adjusta-backend/internal/apps/user"
	"github.com/koo-arch/adjusta-backend/internal/auth"
	"github.com/koo-arch/adjusta-backend/internal/google/calendar"
)

func GetCalendarHandler(client *ent.Client) gin.HandlerFunc {
	return func (c *gin.Context) {
		email, ok := c.Get("email")
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "missing email"})
			c.Abort()
			return
		}

		ctx := c.Request.Context()

		userRepo := user.NewUserRepository(client)
		accountRepo := account.NewAccountRepository(client)
		authManager := auth.NewAuthManager(client, userRepo, accountRepo)

		token, err := authManager.VerifyOAuthToken(ctx, email.(string))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to verify token"})
			c.Abort()
			return
		}

		calendarService, err := calendar.NewCalendar(ctx, token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create calendar service"})
			c.Abort()
			return
		}
		now := time.Now()
		startTime := now.AddDate(0, -2, 0)
		endTime := now.AddDate(1, 0, 0)
		
		events, err := calendarService.FetchEvents(startTime, endTime)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch events"})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, events)
	}
}