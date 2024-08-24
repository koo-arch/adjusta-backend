package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/apps/account"
	"github.com/koo-arch/adjusta-backend/internal/apps/user"
	dbCalendar "github.com/koo-arch/adjusta-backend/internal/apps/calendar"
	"github.com/koo-arch/adjusta-backend/internal/apps/proposeddate"
	"github.com/koo-arch/adjusta-backend/internal/apps/event"
	"github.com/koo-arch/adjusta-backend/internal/auth"
	"github.com/koo-arch/adjusta-backend/internal/google/calendar"
	"github.com/koo-arch/adjusta-backend/internal/models"
)

func FetchEventListHandler(client *ent.Client) gin.HandlerFunc {
	return func (c *gin.Context) {
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

		userRepo := user.NewUserRepository(client)
		accountRepo := account.NewAccountRepository(client)
		authManager := auth.NewAuthManager(client, userRepo, accountRepo)

		userAccounts, err := accountRepo.FilterByUserID(ctx, nil, userid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user accounts"})
			c.Abort()
			return
		}

		calendarRepo := dbCalendar.NewCalendarRepository(client)
		eventRepo := event.NewEventRepository(client)
		dateRepo := proposeddate.NewProposedDateRepository(client)

		eventManager := calendar.NewEventManager(client, authManager, calendarRepo, eventRepo, dateRepo)


		accountsEvents, err := eventManager.FetchAllEvents(ctx, userid, userAccounts)
		if err != nil {
			fmt.Printf("failed to fetch events: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch events"})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, accountsEvents)
	}
}

func CreateEventDraftHandler(client *ent.Client) gin.HandlerFunc {
	return func (c *gin.Context) {
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

		email, ok := c.Get("email")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to get email from context"})
			c.Abort()
			return
		}
		emailStr, ok := email.(string)

		var eventDraft *models.EventDraft
		if err := c.ShouldBindJSON(&eventDraft); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind json"})
			c.Abort()
			return
		}

		userRepo := user.NewUserRepository(client)
		accountRepo := account.NewAccountRepository(client)
		authManager := auth.NewAuthManager(client, userRepo, accountRepo)
		
		a, err := accountRepo.FindByUserIDAndEmail(ctx, nil, userid, emailStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get account"})
			c.Abort()
			return
		}
		
		calendarRepo := dbCalendar.NewCalendarRepository(client)
		eventRepo := event.NewEventRepository(client)
		dateRepo := proposeddate.NewProposedDateRepository(client)
		eventManager := calendar.NewEventManager(client, authManager, calendarRepo, eventRepo, dateRepo)

		err = eventManager.CreateDraftedEvents(ctx, userid, a.ID, emailStr, eventDraft)
		if err != nil {
			fmt.Printf("failed to fetch events: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch events"})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}
}