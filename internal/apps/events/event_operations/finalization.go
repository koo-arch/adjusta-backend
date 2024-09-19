package event_operations

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/apps/events"
	customCalendar "github.com/koo-arch/adjusta-backend/internal/google/calendar"
	"github.com/koo-arch/adjusta-backend/internal/models"
	dbCalendar "github.com/koo-arch/adjusta-backend/internal/repo/calendar"
	"github.com/koo-arch/adjusta-backend/internal/repo/event"
	"github.com/koo-arch/adjusta-backend/internal/repo/proposeddate"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
	"google.golang.org/api/calendar/v3"
)

type EventFinalizationManager struct {
	event *events.EventManager
}

func NewEventFinalizationManager(event *events.EventManager) *EventFinalizationManager {
	return &EventFinalizationManager{
		event: event,
	}
}

func (efm *EventFinalizationManager) FinalizeProposedDate(ctx context.Context, userID, accountID, eventID uuid.UUID, email string, eventReq *models.ConfirmEvent) error {
	tx, err := efm.event.Client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed starting transaction: %w", err)
	}

	// OAuthトークンを検証
	token, err := efm.event.AuthManager.VerifyOAuthToken(ctx, userID, email)
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

	// is_Finalizedをfalseにリセット
	err = efm.resetAllFinalized(ctx, tx, eventID, accountID, email)

	entEvent, err := efm.event.EventRepo.Read(ctx, tx, eventID, event.EventQueryOptions{})
	if err != nil {
		return fmt.Errorf("failed to get event for account: %s, error: %w", email, err)
	}

	// Googleカレンダーイベントの新規登録または既存イベントのIDチェック
	googleEventID, err := efm.handleGoogleEvent(calendarService, entEvent, eventReq)
	if err != nil {
		return fmt.Errorf("failed to handle google event for account: %s, error: %w", email, err)
	}

	// いずれかの日程候補を確定
	err = efm.confirmEventDate(ctx, tx, googleEventID, eventReq, entEvent)
	if err != nil {
		return fmt.Errorf("failed to confirm event date for account: %s, error: %w", email, err)
	}

	// is_finalizedがfalseの日程をGoogleカレンダーから削除
	err = efm.cleanupNotFinalizedDates(ctx, tx, calendarService, eventID)
	if err != nil {
		return fmt.Errorf("failed to cleanup not finalized dates for account: %s, error: %w", email, err)
	}

	// トランザクションをコミット
	return nil
}

func (efm *EventFinalizationManager) resetAllFinalized(ctx context.Context, tx *ent.Tx, eventID, accountID uuid.UUID, email string) error {
	// プライマリカレンダーを取得
	isPrimary := true
	findOptions := dbCalendar.CalendarQueryOptions{
		IsPrimary: &isPrimary,
	}
	_, err := efm.event.CalendarRepo.FindByFields(ctx, tx, accountID, findOptions)
	if err != nil {
		return fmt.Errorf("failed to get primary calendar for account: %s, error: %w", email, err)
	}

	// イベントの日程候補のis_Finalizedを全てfalseに更新
	err = efm.event.DateRepo.ResetFinalized(ctx, tx, eventID)
	if err != nil {
		return fmt.Errorf("failed to reset is_finalized for account: %s, error: %w", email, err)
	}

	return nil
}

func (efm *EventFinalizationManager) handleGoogleEvent(calendarService *customCalendar.Calendar, entEvent *ent.Event, eventReq *models.ConfirmEvent) (*string, error) {
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
		googleEvents, err := efm.event.CalendarApp.CreateGoogleEvents(calendarService, &eventDraftCreate)
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

func (efm *EventFinalizationManager) confirmEventDate(ctx context.Context, tx *ent.Tx, googleEventID *string, eventReq *models.ConfirmEvent, entEvent *ent.Event) error {
	isFinalized := true
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
		IsFinalized:   &isFinalized,
	}

	// 日程候補の新規作成または更新
	if eventReq.ConfirmDate.ID == nil {
		_, err := efm.event.DateRepo.Create(ctx, tx, googleEventID, dateOptions, entEvent)
		if err != nil {
			return fmt.Errorf("failed to create proposed date error: %w", err)
		}
	} else {
		// 日程候補のis_Finalizedを更新
		entDate, err := efm.event.DateRepo.Read(ctx, tx, *eventReq.ConfirmDate.ID)
		if err != nil {
			return fmt.Errorf("failed to get proposed date error: %w", err)
		}

		_, err = efm.event.DateRepo.Update(ctx, tx, entDate.ID, dateOptions)
		if err != nil {
			return fmt.Errorf("failed to update proposed error: %w", err)
		}
	}

	return nil
}

func (efm *EventFinalizationManager) cleanupNotFinalizedDates(ctx context.Context, tx *ent.Tx, calendarService *customCalendar.Calendar, eventID uuid.UUID) error {
	// is_finalizedがfalseの日程候補を検索
	notFinalizedDates, err := efm.event.DateRepo.FilterByEventIDWithFinalized(ctx, tx, eventID, false)
	if err != nil {
		return fmt.Errorf("failed to get not finalized proposed dates: %w", err)
	}

	convEvents := make([]*calendar.Event, len(notFinalizedDates))
	for i, date := range notFinalizedDates {
		convEvents[i] = efm.event.CalendarApp.ConvertToCalendarEvent(&date.GoogleEventID, "", "", "", date.StartTime, date.EndTime)
	}
	// Googleカレンダーのイベントを削除
	err = efm.event.CalendarApp.DeleteGoogleEvents(calendarService, convEvents)
	if err != nil {
		return fmt.Errorf("failed to delete events error: %w", err)
	}

	empty := ""
	// Googleイベントを削除した日程候補のgoogle_event_idを削除
	for _, date := range notFinalizedDates {
		dateOptions := proposeddate.ProposedDateQueryOptions{
			GoogleEventID: &empty,
		}
		_, err = efm.event.DateRepo.Update(ctx, tx, date.ID, dateOptions)
		if err != nil {
			return fmt.Errorf("failed to update proposed date error: %w", err)
		}
	}

	return nil
}
