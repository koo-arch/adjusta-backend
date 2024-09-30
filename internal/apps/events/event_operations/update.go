package event_operations

import (
	"context"
	"log"
	"fmt"
	"time"
	"sync"

	"github.com/google/uuid"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/googleapi"
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

func (eum *EventUpdateManager) UpdateDraftedEvents(ctx context.Context, userID, accountID, eventID uuid.UUID, email string, eventReq *models.EventDraftDetail) error {
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
	entEvent, err := eum.event.EventRepo.Update(ctx, tx, eventID, convEvent)
	if err != nil {
		return fmt.Errorf("failed to get event for account: %s, error: %w", email, err)
	}

	// DB上の日程候補を取得
	existingDates, err := eum.event.DateRepo.FilterByEventID(ctx, tx, eventID)
	if err != nil {
		return fmt.Errorf("failed to get proposed dates, error: %w", err)
	}

	// Googleカレンダーイベントの新規登録または既存イベントのIDチェック
	operationMap, err := eum.SyncUpdateGoogleEvents(ctx, tx, calendarService, eventReq, existingDates)
	if err != nil {
		return fmt.Errorf("failed to handle google event for account: %s, error: %w", email, err)
	}

	// DB上の日程候補を更新
	err = eum.updateProposedDates(ctx, tx, eventReq, entEvent, existingDates)
	if err != nil {
		if rollbackErr := eum.RollbackGoogleEvents(ctx, tx, calendarService, operationMap); rollbackErr != nil {
			return fmt.Errorf("failed to rollback google events: %w", rollbackErr)
		}
	}

	// トランザクションをコミット
	return nil
}

type CalendarOperation struct {
	OperationType string
	BackupEvent   *calendar.Event
}

func (eum *EventUpdateManager) SyncUpdateGoogleEvents(ctx context.Context, tx *ent.Tx, calendarService *customCalendar.Calendar, eventReq *models.EventDraftDetail, existingDates []*ent.ProposedDate) (map[string]CalendarOperation, error) {
	// Googleカレンダーに対する操作のリスト
	operationsMap := make(map[string]CalendarOperation)

	// 提案された日程のGoogleEventIDをハッシュマップにして、削除対象を判定する
	proposedDateMap := make(map[string]struct{})
	for _, date := range eventReq.ProposedDates {
		if date.GoogleEventID != "" {
			proposedDateMap[date.GoogleEventID] = struct{}{}
		}
	}

	// 並列処理でイベントを登録
	var wg sync.WaitGroup
	var mu sync.Mutex
	errCh := make(chan error, len(eventReq.ProposedDates))

	// 更新・作成の処理
	wg.Add(len(eventReq.ProposedDates))
	for i, date := range eventReq.ProposedDates {
		go func(i int, date models.ProposedDate) {
			defer wg.Done()

			if date.GoogleEventID == "" {
				// 新規作成の場合
				event := eum.event.CalendarApp.ConvertToCalendarEvent(nil, eventReq.Title, eventReq.Location, eventReq.Description, *date.Start, *date.End)
				newEvent, err := calendarService.InsertEvent(event)
				if err != nil {
					errCh <- fmt.Errorf("failed to create event in Google Calendar: %w", err)
					return
				}

				// 新規作成されたGoogleEventIDを保存
				mu.Lock()
				eventReq.ProposedDates[i].GoogleEventID = newEvent.Id // GoogleEventIDを保存
				operationsMap[newEvent.Id] = CalendarOperation{OperationType: "create"}
				mu.Unlock()

			} else {
				// 更新の場合
				backupEvent, err := calendarService.FetchEvent(date.GoogleEventID)
				if err != nil {
					errCh <- fmt.Errorf("failed to fetch event from Google Calendar: %w", err)
					return
				}

				event := eum.event.CalendarApp.ConvertToCalendarEvent(&date.GoogleEventID, eventReq.Title, eventReq.Location, eventReq.Description, *date.Start, *date.End)
				_, err = calendarService.UpdateEvent(date.GoogleEventID, event)
				if err != nil {
					errCh <- fmt.Errorf("failed to update event in Google Calendar: %w", err)
					return
				}

				// 更新されたGoogleEventIDを保存
				mu.Lock()
				operationsMap[date.GoogleEventID] = CalendarOperation{OperationType: "update", BackupEvent: backupEvent}
				mu.Unlock()
			}
		}(i, date)
	}

	// 削除対象の判定と削除処理
	wg.Add(len(existingDates))
	for _, date := range existingDates {
		go func(date *ent.ProposedDate) {
			defer wg.Done()

			if _, exists := proposedDateMap[date.GoogleEventID]; !exists {
				backupEvent, err := calendarService.FetchEvent(date.GoogleEventID)
				if err != nil {
					if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
						// 404エラーは削除済みのイベントなのでスキップ
						log.Printf("Event %s not found (404), skipping fetch", date.GoogleEventID)
						return
					}
					errCh <- fmt.Errorf("failed to fetch event from Google Calendar: %w", err)
					return
				}

				err = calendarService.DeleteEvent(date.GoogleEventID)
				if err != nil {
					if gerr, ok := err.(*googleapi.Error); ok && (gerr.Code == 410 || gerr.Code == 404) {
						// 410エラーは削除済みのイベントなのでスキップ
						log.Printf("Event %s already deleted (410), skipping delete", date.GoogleEventID)
						return
					}
					errCh <- fmt.Errorf("failed to delete event in Google Calendar: %w", err)
					return
				}

				mu.Lock()
				operationsMap[date.GoogleEventID] = CalendarOperation{OperationType: "delete", BackupEvent: backupEvent} // 削除操作を記録
				mu.Unlock()
			}
		}(date)
	}

	// 全てのゴルーチンの終了を待機
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// エラーがあれば返す
	for err := range errCh {
		if err != nil {
			// エラーがあれば、ロールバックを実行
			if rollbackErr := eum.RollbackGoogleEvents(ctx, tx, calendarService, operationsMap); rollbackErr != nil {
				return nil, fmt.Errorf("failed to rollback google events: %w", rollbackErr)
			}
		}
	}

	// 正常終了なら、操作内容を返す
	return operationsMap, nil
}

func (eum *EventUpdateManager) updateProposedDates(ctx context.Context, tx *ent.Tx, eventReq *models.EventDraftDetail, entEvent *ent.Event, existingDates []*ent.ProposedDate) error {
	// 提案された日程候補のハッシュテーブルを作成
	updateDateMap := make(map[uuid.UUID]models.ProposedDate)
	for _, date := range eventReq.ProposedDates {
		updateDateMap[date.ID] = date
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
				IsFinalized:   &updateDate.IsFinalized,
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
		// 必要に応じてGoogleEventIDを生成または確認
		if date.GoogleEventID == "" {
			return fmt.Errorf("GoogleEventID is missing for new proposed date")
		}

		// 新しい日程候補をDBに追加
		dateOptions := proposeddate.ProposedDateQueryOptions{
			StartTime: date.Start,
			EndTime:   date.End,
			Priority:  &date.Priority,
			GoogleEventID: &date.GoogleEventID,
		}
		_, err := eum.event.DateRepo.Create(ctx, tx, &date.GoogleEventID, dateOptions, entEvent)
		if err != nil {
			return fmt.Errorf("failed to create proposed dates, error: %w", err)
		}
	}

	return nil
}

func (eum *EventUpdateManager) RollbackGoogleEvents(ctx context.Context, tx *ent.Tx, calendarService *customCalendar.Calendar, operationMap map[string]CalendarOperation) error {
	var rollbackErrors []error // エラーを蓄積するスライス

	for googleEventID, operation := range operationMap {
		switch operation.OperationType {
		case "create":
			// 新規作成されたイベントを削除
			if delErr := calendarService.DeleteEvent(googleEventID); delErr != nil {
				rollbackErrors = append(rollbackErrors, fmt.Errorf("failed to delete event %s from Google Calendar: %w", googleEventID, delErr))
			}

		case "update":
			// 更新されたイベントを元に戻す
			if _, err := calendarService.UpdateEvent(operation.BackupEvent.Id, operation.BackupEvent); err != nil {
				rollbackErrors = append(rollbackErrors, fmt.Errorf("failed to update event %s to original state: %w", googleEventID, err))
			}

		case "delete":
			// 削除されたイベントを再作成し、GoogleEventIDをDBに反映
			createEvent, err := calendarService.InsertEvent(operation.BackupEvent)
			if err != nil {
				rollbackErrors = append(rollbackErrors, fmt.Errorf("failed to insert event %s back to Google Calendar: %w", operation.BackupEvent.Id, err))
				continue // 次の処理へ
			}

			dateOperation := proposeddate.ProposedDateQueryOptions{
				GoogleEventID: &createEvent.Id,
			}

			// GoogleEventIDを更新
			err = eum.event.DateRepo.UpdateByGoogleEventID(ctx, tx, &operation.BackupEvent.Id, dateOperation)
			if err != nil {
				rollbackErrors = append(rollbackErrors, fmt.Errorf("failed to update proposed date with new GoogleEventID %s: %w", createEvent.Id, err))
			}
		}
	}

	// 蓄積されたエラーを返す（エラーがあればまとめて報告）
	if len(rollbackErrors) > 0 {
		return fmt.Errorf("rollback completed with errors: %v", rollbackErrors)
	}

	return nil
}