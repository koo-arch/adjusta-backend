package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/models"
	"github.com/koo-arch/adjusta-backend/queryparser"
)

type CalendarHandler struct {
	handler *Handler
}

func NewCalendarHandler(handler *Handler) *CalendarHandler {
	return &CalendarHandler{handler: handler}
}

func (ch *CalendarHandler) FetchEventListHandler() gin.HandlerFunc {
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

		eventFetchingManager := ch.handler.Server.EventFetchingManager

		accountsEvents, err := eventFetchingManager.FetchAllGoogleEvents(ctx, userid, email.(string))
		if err != nil {
			fmt.Printf("failed to fetch events: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch events"})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, accountsEvents)
	}
}

func (ch *CalendarHandler) FetchAllEventDraftListHandler() gin.HandlerFunc {
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

		eventFetchingManager := ch.handler.Server.EventFetchingManager

		draftedEvents, err := eventFetchingManager.FetchAllDraftedEvents(ctx, userid, emailStr)
		if err != nil {
			fmt.Printf("failed to fetch events: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch events"})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, draftedEvents)
	}
}

func (ch *CalendarHandler) SearchEventDraftHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// クエリパラメータの取得
		queryparser := queryparser.NewQueryParser(c)

		// クエリパラメータの解析
		query, err := queryparser.ParseSearchEventQuery()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse query"})
			c.Abort()
			return
		}
		

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

		eventFetchingManager := ch.handler.Server.EventFetchingManager

		draftedEvents, err := eventFetchingManager.SearchDraftedEvents(ctx, userid, emailStr, *query)
		if err != nil {
			fmt.Printf("failed to fetch events: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch events"})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, draftedEvents)
	}
}

func (ch *CalendarHandler) FetchUpcomingEventsHandler() gin.HandlerFunc {
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

		eventFetchingManager := ch.handler.Server.EventFetchingManager

		daysBefore := 3
		upcomingEvents, err := eventFetchingManager.FetchUpcomingEvents(ctx, userid, email.(string), daysBefore)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch events"})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, upcomingEvents)
	}
}

func (ch *CalendarHandler) FetchEventDraftDetailHandler() gin.HandlerFunc {
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

		eventFetchingManager := ch.handler.Server.EventFetchingManager

		draftedEvent, err := eventFetchingManager.FetchDraftedEventDetail(ctx, userid, emailStr, eventID)
		if err != nil {
			fmt.Printf("failed to fetch events: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch events"})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, draftedEvent)
	}
}

func (ch *CalendarHandler) CreateEventDraftHandler() gin.HandlerFunc {
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

		eventCreationManager := ch.handler.Server.EventCreationManager

		err = eventCreationManager.CreateDraftedEvents(ctx, userid, emailStr, eventDraft)
		if err != nil {
			fmt.Printf("failed to fetch events: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch events"})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}
}

func (ch *CalendarHandler) EventFinalizeHandler() gin.HandlerFunc {
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

		eventManager := ch.handler.Server.EventManager

		err = eventManager.FinalizeProposedDate(ctx, userid, eventID, emailStr, confirmEvent)
		if err != nil {
			fmt.Printf("failed to finalize event: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to finalize event"})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}
}

func (ch *CalendarHandler) UpdateEventDraftHandler() gin.HandlerFunc {
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

		var eventDraft *models.EventDraftDetail
		if err := c.ShouldBindJSON(&eventDraft); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind json"})
			c.Abort()
			return
		}

		eventUpdateManager := ch.handler.Server.EventUpdateManager

		err = eventUpdateManager.UpdateDraftedEvents(ctx, userid, eventID, emailStr, eventDraft)
		if err != nil {
			fmt.Printf("failed to fetch events: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch events"})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}
}

func (ch *CalendarHandler) DeleteEventDraftHandler() gin.HandlerFunc {
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

		var eventDraft *models.EventDraftDetail
		if err := c.ShouldBindJSON(&eventDraft); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind json"})
			c.Abort()
			return
		}

		eventDeleteManager := ch.handler.Server.EventDeleteManager

		err = eventDeleteManager.DeleteDraftedEvents(ctx, userid, emailStr, eventDraft)
		if err != nil {
			fmt.Printf("failed to fetch events: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch events"})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}
}