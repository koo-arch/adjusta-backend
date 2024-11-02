package events

import (
	"context"
	"fmt"

	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/google/uuid"
	appCalendar "github.com/koo-arch/adjusta-backend/internal/apps/calendar"
	"github.com/koo-arch/adjusta-backend/internal/auth"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/repo/calendar"
	"github.com/koo-arch/adjusta-backend/internal/repo/event"
	"github.com/koo-arch/adjusta-backend/internal/repo/proposeddate"
	customCalendar "github.com/koo-arch/adjusta-backend/internal/google/calendar"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
	"github.com/koo-arch/adjusta-backend/internal/models"
)

type EventManager struct {
	Client       *ent.Client
	AuthManager  *auth.AuthManager
	CalendarRepo repoCalendar.CalendarRepository
	EventRepo    event.EventRepository
	DateRepo     proposeddate.ProposedDateRepository
	CalendarApp  *appCalendar.GoogleCalendarManager
}

func NewEventManager(
	client *ent.Client,
	authManager *auth.AuthManager,
	calendarRepo repoCalendar.CalendarRepository,
	eventRepo event.EventRepository,
	dateRepo proposeddate.ProposedDateRepository,
	calendarApp *appCalendar.GoogleCalendarManager,
) *EventManager {
	return &EventManager{
		Client:       client,
		AuthManager:  authManager,
		CalendarRepo: calendarRepo,
		EventRepo:    eventRepo,
		DateRepo:     dateRepo,
		CalendarApp:  calendarApp,
	}
}

func (em *EventManager) FinalizeProposedDate(ctx context.Context, userID, eventID uuid.UUID, email string, eventReq *models.ConfirmEvent) error {
	tx, err := em.Client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed starting transaction: %w", err)
	}

	// OAuthトークンを検証
	token, err := em.AuthManager.VerifyOAuthToken(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to verify token for account: %s, error: %w", email, err)
	}

	// Google Calendarサービスを作成
	calendarService, err := customCalendar.NewCalendar(ctx, token)
	if err != nil {
		return fmt.Errorf("failed to create calendar service for account: %s, error: %w", email, err)
	}

	// トランザクションをデファーで処理
	defer transaction.HandleTransaction(tx, &err)

	entEvent, err := em.EventRepo.Read(ctx, tx, eventID, event.EventQueryOptions{})
	if err != nil {
		return fmt.Errorf("failed to get event for account: %s, error: %w", email, err)
	}

	// Googleカレンダーイベントの新規登録または既存イベントのIDチェック
	googleEventID, err := em.HandleGoogleEvent(calendarService, entEvent, eventReq)
	if err != nil {
		return fmt.Errorf("failed to handle google event for account: %s, error: %w", email, err)
	}

	// いずれかの日程候補を確定
	err = em.ConfirmEventDate(ctx, tx, googleEventID, eventReq, entEvent)
	if err != nil {
		return fmt.Errorf("failed to confirm event date for account: %s, error: %w", email, err)
	}

	// トランザクションをコミット
	return nil
}

func (em *EventManager) HandleGoogleEvent(calendarService *customCalendar.Calendar, entEvent *ent.Event, eventReq *models.ConfirmEvent) (*string, error) {
	// Googleカレンダーイベントの新規登録または既存イベントのIDチェック
	var googleEventID *string
	if eventReq.ConfirmDate.ID == nil || eventReq.ConfirmDate.GoogleEventID == "" {
		// 登録するイベントの情報を作成
		eventDraftCreate := models.EventDraftCreation{
			Title:       entEvent.Summary,
			Location:    entEvent.Location,
			Description: entEvent.Description,
			SelectedDates: []models.SelectedDate{
				{
					Start: *eventReq.ConfirmDate.Start,
					End:   *eventReq.ConfirmDate.End,
				},
			},
		}
		// googleカレンダーにイベントを登録
		googleEvents, err := em.CalendarApp.CreateGoogleEvents(calendarService, &eventDraftCreate)
		if err != nil {
			return nil, fmt.Errorf("failed to insert events, error: %w", err)
		}
		googleEventID = &googleEvents[0].Id

	} else {
		// 既存のGoogleカレンダーイベントIDを使用
		googleEventID = &eventReq.ConfirmDate.GoogleEventID
	}

	return googleEventID, nil
}

func (em *EventManager) ConfirmEventDate(ctx context.Context, tx *ent.Tx, googleEventID *string, eventReq *models.ConfirmEvent, entEvent *ent.Event) error {
	priority := 0
	// 優先度が設定されている場合は設定
	if eventReq.ConfirmDate.Priority > 0 {
		priority = eventReq.ConfirmDate.Priority
	}
	dateOptions := proposeddate.ProposedDateQueryOptions{
		GoogleEventID: googleEventID,
		StartTime:     eventReq.ConfirmDate.Start,
		EndTime:       eventReq.ConfirmDate.End,
		Priority:      &priority,
	}

	// 確定日程がDBに存在しない場合は新規作成
	if eventReq.ConfirmDate.ID == nil {
		_, err := em.DateRepo.Create(ctx, tx, googleEventID, dateOptions, entEvent)
		if err != nil {
			return fmt.Errorf("failed to create proposed date error: %w", err)
		}
	}

	// イベントステータスと確定日程IDを更新
	status := models.StatusConfirmed
	eventOptions := event.EventQueryOptions{
		Status: &status,
		ConfirmedDateID: eventReq.ConfirmDate.ID,
	}
	_, err := em.EventRepo.Update(ctx, tx, entEvent.ID, eventOptions)
	if err != nil {
		return fmt.Errorf("failed to update event status error: %w", err)
	}

	return nil
}
