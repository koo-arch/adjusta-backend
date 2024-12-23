package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/models"
	"github.com/koo-arch/adjusta-backend/queryparser"
	"github.com/koo-arch/adjusta-backend/utils"
)


type CalendarHandler struct {
	handler *Handler
}

func NewCalendarHandler(handler *Handler) *CalendarHandler {
	return &CalendarHandler{handler: handler}
}

var extractErrorMessage = "ユーザー情報確認時にエラーが発生しました。"

func (ch *CalendarHandler) FetchEventListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userid, email, err := utils.ExtractUserIDAndEmail(c)
		if err != nil {
			utils.HandleAPIError(c, err, extractErrorMessage)
			return
		}

		eventFetchingManager := ch.handler.Server.EventFetchingManager

		accountsEvents, err := eventFetchingManager.FetchAllGoogleEvents(ctx, userid, email)
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

		userid, email,  err := utils.ExtractUserIDAndEmail(c)
		if err != nil {
			utils.HandleAPIError(c, err, extractErrorMessage)
			return
		}

		eventFetchingManager := ch.handler.Server.EventFetchingManager

		draftedEvents, err := eventFetchingManager.FetchAllDraftedEvents(ctx, userid, email)
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

		userid, email, err := utils.ExtractUserIDAndEmail(c)
		if err != nil {
			utils.HandleAPIError(c, err, extractErrorMessage)
			return
		}

		eventFetchingManager := ch.handler.Server.EventFetchingManager

		draftedEvents, err := eventFetchingManager.SearchDraftedEvents(ctx, userid, email, *query)
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

		userid, email, err := utils.ExtractUserIDAndEmail(c)
		if err != nil {
			utils.HandleAPIError(c, err, extractErrorMessage)
			return
		}

		eventFetchingManager := ch.handler.Server.EventFetchingManager

		daysBefore := 3
		upcomingEvents, err := eventFetchingManager.FetchUpcomingEvents(ctx, userid, email, daysBefore)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch events"})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, upcomingEvents)
	}
}

func (ch *CalendarHandler) FetchNeedsActionDraftsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userid, email, err := utils.ExtractUserIDAndEmail(c)
		if err != nil {
			utils.HandleAPIError(c, err, extractErrorMessage)
			return
		}

		eventFetchingManager := ch.handler.Server.EventFetchingManager

		daysBefore := 3
		upcomingEvents, err := eventFetchingManager.FetchNeedsActionDrafts(ctx, userid, email, daysBefore)
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

		userid, email, err := utils.ExtractUserIDAndEmail(c)
		if err != nil {
			utils.HandleAPIError(c, err, extractErrorMessage)
			return
		}

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

		draftedEvent, err := eventFetchingManager.FetchDraftedEventDetail(ctx, userid, email, eventID)
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

		userid, email, err := utils.ExtractUserIDAndEmail(c)
		if err != nil {
			utils.HandleAPIError(c, err, extractErrorMessage)
			return
		}

		var eventDraft *models.EventDraftCreation
		if err := c.ShouldBindJSON(&eventDraft); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind json"})
			c.Abort()
			return
		}

		eventCreationManager := ch.handler.Server.EventCreationManager

		response, err := eventCreationManager.CreateDraftedEvents(ctx, userid, email, eventDraft)
		if err != nil {
			fmt.Printf("failed to fetch events: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch events"})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

func (ch *CalendarHandler) EventFinalizeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userid, email, err := utils.ExtractUserIDAndEmail(c)
		if err != nil {
			utils.HandleAPIError(c, err, extractErrorMessage)
			return
		}

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

		err = eventManager.FinalizeProposedDate(ctx, userid, eventID, email, confirmEvent)
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

		userid, email, err := utils.ExtractUserIDAndEmail(c)
		if err != nil {
			utils.HandleAPIError(c, err, extractErrorMessage)
			return
		}

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

		err = eventUpdateManager.UpdateDraftedEvents(ctx, userid, eventID, email, eventDraft)
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

		userid, email, err := utils.ExtractUserIDAndEmail(c)
		if err != nil {
			utils.HandleAPIError(c, err, extractErrorMessage)
			return
		}

		var eventDraft *models.EventDraftDetail
		if err := c.ShouldBindJSON(&eventDraft); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to bind json"})
			c.Abort()
			return
		}

		eventDeleteManager := ch.handler.Server.EventDeleteManager

		err = eventDeleteManager.DeleteDraftedEvents(ctx, userid, email, eventDraft)
		if err != nil {
			fmt.Printf("failed to fetch events: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch events"})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}
}
