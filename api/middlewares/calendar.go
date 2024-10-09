package middlewares

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/google/calendar"
)

type CalendarMiddleware struct {
	middleware *Middleware
}

func NewCalendarMiddleware(middleware *Middleware) *CalendarMiddleware {
	return &CalendarMiddleware{middleware: middleware}
}

func (cm *CalendarMiddleware) SyncGoogleCalendars() gin.HandlerFunc {
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
		accountRepo := cm.middleware.Server.AccountRepo
		calendarRepo := cm.middleware.Server.CalendarRepo
		authManager := cm.middleware.Server.AuthManager

		userAccounts, err := accountRepo.FilterByUserID(ctx, nil, userid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user accounts"})
			c.Abort()
			return
		}

		calendarList, err := calendar.RegisterCalendarList(ctx, authManager, userid, userAccounts, calendarRepo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register calendar list"})
			c.Abort()
			return
		}

		c.Set("calendarList", calendarList)
		c.Next()
	}
}
