package event_operations

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/apps/events"
	customCalendar "github.com/koo-arch/adjusta-backend/internal/google/calendar"
	"github.com/koo-arch/adjusta-backend/internal/models"
	dbCalendar "github.com/koo-arch/adjusta-backend/internal/repo/calendar"
	"github.com/koo-arch/adjusta-backend/internal/repo/proposeddate"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
)

type EventUpdateManager struct {
	event *events.EventManager
}

func NewEventUpdateManager(event *events.EventManager) *EventUpdateManager {
	return &EventUpdateManager{
		event: event,
	}
}

func (eum *EventUpdateManager) UpdateDraftedEvents(ctx context.Context, userID, accountID uuid.UUID, email string, eventReq *models.EventDraftDetail) error {
	tx, err := eum.event.Client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed starting transaction: %w", err)
	}

	// OAuthトークンを検証
	token, err := eum.event.AuthManager.VerifyOAuthToken(ctx, userID, email)
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

	// プライマリカレンダーを取得
	isPrimary := true
	findOptions := dbCalendar.CalendarQueryOptions{
		IsPrimary: &isPrimary,
	}
	_, err = eum.event.CalendarRepo.FindByFields(ctx, tx, accountID, findOptions)
	if err != nil {
		return fmt.Errorf("failed to get primary calendar for account: %s, error: %w", email, err)
	}
	// calendar.Event型に変換
	convEvent := eum.event.CalendarApp.ConvertToCalendarEvent(nil, eventReq.Title, eventReq.Location, eventReq.Description, time.Time{}, time.Time{})

	// イベントの詳細を更新
	entEvent, err := eum.event.EventRepo.Update(ctx, tx, eventReq.ID, convEvent)
	if err != nil {
		return fmt.Errorf("failed to get event for account: %s, error: %w", email, err)
	}

	// 日程候補の変更処理
	err = eum.updateProposedDates(ctx, tx, eventReq, entEvent)
	if err != nil {
		return fmt.Errorf("failed to update proposed dates for account: %s, error: %w", email, err)
	}

	// Googleカレンダーのイベントを更新
	err = eum.event.CalendarApp.UpdateGoogleCalendarEvents(calendarService, eventReq)
	if err != nil {
		return fmt.Errorf("failed to update events for account: %s, error: %w", email, err)
	}

	// トランザクションをコミット
	return nil
}

func (eum *EventUpdateManager) updateProposedDates(ctx context.Context, tx *ent.Tx, eventReq *models.EventDraftDetail, entEvent *ent.Event) error {
	// 日程候補を取得
	existingDates, err := eum.event.DateRepo.FilterByEventID(ctx, tx, eventReq.ID)
	if err != nil {
		return fmt.Errorf("failed to get proposed dates, error: %w", err)
	}

	// ハッシュテーブルを作成
	updateDateMap := make(map[uuid.UUID]models.ProposedDate)
	for _, date := range eventReq.ProposedDates {
		updateDateMap[date.ID] = date
	}

	// 日程候補を更新または削除
	for _, date := range existingDates {
		if updateDate, ok := updateDateMap[date.ID]; ok {

			dateOptions := proposeddate.ProposedDateQueryOptions{
				GoogleEventID: &updateDate.GoogleEventID,
				StartTime:     updateDate.Start,
				EndTime:       updateDate.End,
				Priority:      &updateDate.Priority,
				IsFinalized:   &updateDate.IsFinalized,
			}
			_, err = eum.event.DateRepo.Update(ctx, tx, date.ID, dateOptions)
			if err != nil {
				return fmt.Errorf("failed to update proposed dates, error: %w", err)
			}
			// 更新した日程候補を削除
			delete(updateDateMap, date.ID)
		} else {
			err = eum.event.DateRepo.Delete(ctx, tx, date.ID)
			if err != nil {
				return fmt.Errorf("failed to delete proposed dates, error: %w", err)
			}
		}
	}

	// DBに存在しない日程候補を追加
	for _, date := range updateDateMap {
		dateOptions := proposeddate.ProposedDateQueryOptions{
			StartTime: date.Start,
			EndTime:   date.End,
			Priority:  &date.Priority,
		}
		_, err = eum.event.DateRepo.Create(ctx, tx, &date.GoogleEventID, dateOptions, entEvent)
		if err != nil {
			return fmt.Errorf("failed to create proposed dates, error: %w", err)
		}
	}

	return nil
}
