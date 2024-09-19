package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/models"
)

func (s *Server) FetchEventListHandler(client *ent.Client) gin.HandlerFunc {
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

		accountsEvents, err := s.eventFetchingManager.FetchAllEvents(ctx, userid, userAccounts)
		if err != nil {
			fmt.Printf("failed to fetch events: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch events"})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, accountsEvents)
	}
}

func (s *Server) FetchEventDraftListHandler(client *ent.Client) gin.HandlerFunc {
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

		email, ok := c.Get("email")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to get email from context"})
			c.Abort()
			return
		}

		emailStr, ok := email.(string)

		userAccount, err := s.accountRepo.FindByUserIDAndEmail(ctx, nil, userid, emailStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user accounts"})
			c.Abort()
			return
		}

		draftedEvents, err := s.eventFetchingManager.FetchDraftedEvents(ctx, userid, userAccount.ID, emailStr)
		if err != nil {
			fmt.Printf("failed to fetch events: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch events"})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, draftedEvents)
	}
}

func (s *Server) FetchEventDraftDetailHandler(client *ent.Client) gin.HandlerFunc {
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

		email, ok := c.Get("email")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to get email from context"})
			c.Abort()
			return
		}

		emailStr, ok := email.(string)

		eventIDParam := c.Param("eventID")
		if eventIDParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing eventID"})
			c.Abort()
			return
		}

		eventID, err := uuid.Parse(eventIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid eventID format"})
			c.Abort()
			return
		}

		userAccount, err := s.accountRepo.FindByUserIDAndEmail(ctx, nil, userid, emailStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user accounts"})
			c.Abort()
			return
		}

		draftedEvent, err := s.eventFetchingManager.FetchDraftedEventDetail(ctx, userid, userAccount.ID, emailStr, eventID)
		if err != nil {
			fmt.Printf("failed to fetch events: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch events"})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, draftedEvent)
	}
}

func (s *Server) CreateEventDraftHandler(client *ent.Client) gin.HandlerFunc {
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

		email, ok := c.Get("email")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to get email from context"})
			c.Abort()
			return
		}
		emailStr, ok := email.(string)

		var eventDraft *models.EventDraftCreation
		if err := c.ShouldBindJSON(&eventDraft); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind json"})
			c.Abort()
			return
		}

		a, err := s.accountRepo.FindByUserIDAndEmail(ctx, nil, userid, emailStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get account"})
			c.Abort()
			return
		}

		err = s.eventCreationManager.CreateDraftedEvents(ctx, userid, a.ID, emailStr, eventDraft)
		if err != nil {
			fmt.Printf("failed to fetch events: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch events"})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}
}

func (s *Server) EventFinalizeHandler(client *ent.Client) gin.HandlerFunc {
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

		email, ok := c.Get("email")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to get email from context"})
			c.Abort()
			return
		}
		emailStr, ok := email.(string)

		eventIDParam := c.Param("eventID")
		if eventIDParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing eventID"})
			c.Abort()
			return
		}

		eventID, err := uuid.Parse(eventIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid eventID format"})
			c.Abort()
			return
		}

		var confirmEvent *models.ConfirmEvent
		if err := c.ShouldBindJSON(&confirmEvent); err != nil {
			fmt.Printf("failed to bind json: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind json"})
			c.Abort()
			return
		}

		userAccount, err := s.accountRepo.FindByUserIDAndEmail(ctx, nil, userid, emailStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user accounts"})
			c.Abort()
			return
		}

		err = s.eventFinalizationManager.FinalizeProposedDate(ctx, userid, userAccount.ID, eventID, emailStr, confirmEvent)
		if err != nil {
			fmt.Printf("failed to finalize event: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to finalize event"})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}
}
