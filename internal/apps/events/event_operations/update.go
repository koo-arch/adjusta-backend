package event_operations

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/repo/event"
	"github.com/koo-arch/adjusta-backend/internal/apps/events"
	customCalendar "github.com/koo-arch/adjusta-backend/internal/google/calendar"
	"github.com/koo-arch/adjusta-backend/internal/models"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/repo/calendar"
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

func (eum *EventUpdateManager) UpdateDraftedEvents(ctx context.Context, userID, eventID uuid.UUID, email string, eventReq *models.EventDraftDetail) error {
	tx, err := eum.event.Client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed starting transaction: %w", err)
	}

	// トランザクションをデファーで処理
	defer transaction.HandleTransaction(tx, &err)

	// プライマリカレンダーを取得
	isPrimary := true
	findOptions := repoCalendar.CalendarQueryOptions{
		IsPrimary: &isPrimary,
	}
	_, err = eum.event.CalendarRepo.FindByFields(ctx, tx, userID, findOptions)
	if err != nil {
		return fmt.Errorf("failed to get primary calendar for account: %s, error: %w", email, err)
	}

	// イベントの詳細を更新
	eventOptions := event.EventQueryOptions{
		Summary:      &eventReq.Title,
		Location:     &eventReq.Location,
		Description:  &eventReq.Description,
		Status:       &eventReq.Status,
	}
	entEvent, err := eum.event.EventRepo.Update(ctx, tx, eventID, eventOptions)
	if err != nil {
		return fmt.Errorf("failed to get event for account: %s, error: %w", email, err)
	}

	// DB上の日程候補を取得
	existingDates, err := eum.event.DateRepo.FilterByEventID(ctx, tx, eventID)
	if err != nil {
		return fmt.Errorf("failed to get proposed dates, error: %w", err)
	}

	// イベントのステータスによって処理を分岐
	if eventReq.Status == models.StatusConfirmed {
		// 確定済みの場合
		confirmDate := models.ConfirmDate{
			ID: 			eventReq.ProposedDates[0].ID,
			GoogleEventID: 	eventReq.ProposedDates[0].GoogleEventID,
			Start:  		eventReq.ProposedDates[0].Start,
			End:    		eventReq.ProposedDates[0].End,
			Priority: 		eventReq.ProposedDates[0].Priority,
		}
		confirmEvent := models.ConfirmEvent{
			ConfirmDate: confirmDate,
		}

		// OAuthトークンを検証
		token, err := eum.event.AuthManager.VerifyOAuthToken(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to verify token for account: %s, error: %w", email, err)
		}

		// Google Calendarサービスを作成
		calendarService, err := customCalendar.NewCalendar(ctx, token)
		if err != nil {
			return fmt.Errorf("failed to create calendar service for account: %s, error: %w", email, err)
		}

		// Googleカレンダーイベントの新規登録または既存イベントのIDチェック
		googleEventID, err := eum.event.HandleGoogleEvent(calendarService, entEvent, &confirmEvent)
		if err != nil {
			return fmt.Errorf("failed to handle google event for account: %s, error: %w", email, err)
		}

		// いずれかの日程候補を確定
		err = eum.event.ConfirmEventDate(ctx, tx, googleEventID, &confirmEvent, entEvent)
		if err != nil {
			return fmt.Errorf("failed to confirm event date for account: %s, error: %w", email, err)
		}
	} 
	
	// DB上の日程候補を更新
	err = eum.updateProposedDates(ctx, tx, eventReq, entEvent, existingDates)
	if err != nil {
		return fmt.Errorf("failed to update proposed dates, error: %w", err)
	}

	// トランザクションをコミット
	return nil
}


func (eum *EventUpdateManager) updateProposedDates(ctx context.Context, tx *ent.Tx, eventReq *models.EventDraftDetail, entEvent *ent.Event, existingDates []*ent.ProposedDate) error {
	// 提案された日程候補のハッシュテーブルを作成
	updateDateMap := make(map[uuid.UUID]models.ProposedDate)
	for _, date := range eventReq.ProposedDates {
		if date.ID != nil {
			updateDateMap[*date.ID] = date
		} else {
			// 新規の日程候補の場合、なんでもいいので一意なマップキーを生成
			updateDateMap[uuid.New()] = date
		}
	}

	// 既存の日程候補を更新または削除
	for _, date := range existingDates {
		if updateDate, ok := updateDateMap[date.ID]; ok {
			// 既存のイベントを更新
			dateOptions := proposeddate.ProposedDateQueryOptions{
				GoogleEventID: &updateDate.GoogleEventID,
				StartTime:     updateDate.Start,
				EndTime:       updateDate.End,
				Priority:      &updateDate.Priority,
			}
			_, err := eum.event.DateRepo.Update(ctx, tx, date.ID, dateOptions)
			if err != nil {
				return fmt.Errorf("failed to update proposed dates, error: %w", err)
			}
			// 更新済みの候補をハッシュマップから削除
			delete(updateDateMap, date.ID)
		} else {
			// DB上の不要な日程候補を削除
			err := eum.event.DateRepo.Delete(ctx, tx, date.ID)
			if err != nil {
				return fmt.Errorf("failed to delete proposed dates, error: %w", err)
			}
		}
	}

	// DBに存在しない新しい日程候補を追加
	for _, date := range updateDateMap {

		// 新しい日程候補をDBに追加
		dateOptions := proposeddate.ProposedDateQueryOptions{
			StartTime:     date.Start,
			EndTime:       date.End,
			Priority:      &date.Priority,
			GoogleEventID: &date.GoogleEventID,
		}
		_, err := eum.event.DateRepo.Create(ctx, tx, &date.GoogleEventID, dateOptions, entEvent)
		if err != nil {
			return fmt.Errorf("failed to create proposed dates, error: %w", err)
		}
	}

	return nil
}
