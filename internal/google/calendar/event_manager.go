package calendar

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	dbCalendar "github.com/koo-arch/adjusta-backend/internal/apps/calendar"
	"github.com/koo-arch/adjusta-backend/internal/apps/event"
	"github.com/koo-arch/adjusta-backend/internal/apps/proposeddate"
	"github.com/koo-arch/adjusta-backend/internal/auth"
	"github.com/koo-arch/adjusta-backend/internal/models"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
	"google.golang.org/api/calendar/v3"
)

type EventManager struct {
	client       *ent.Client
	authManager  *auth.AuthManager
	calendarRepo dbCalendar.CalendarRepository
	eventRepo    event.EventRepository
	dateRepo     proposeddate.ProposedDateRepository
}

func NewEventManager(client *ent.Client, authManager *auth.AuthManager, calendarRepo dbCalendar.CalendarRepository, eventRepo event.EventRepository, dateRepo proposeddate.ProposedDateRepository) *EventManager {
	return &EventManager{
		client:       client,
		authManager:  authManager,
		calendarRepo: calendarRepo,
		eventRepo:    eventRepo,
		dateRepo:     dateRepo,
	}
}

func (em *EventManager) FetchAllEvents(ctx context.Context, userID uuid.UUID, userAccounts []*ent.Account) ([]*models.AccountsEvents, error) {
	var accountsEvents []*models.AccountsEvents

	for _, userAccount := range userAccounts {
		token, err := em.authManager.VerifyOAuthToken(ctx, userID, userAccount.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to verify token for account: %s, error: %w", userAccount.Email, err)
		}

		calendarService, err := NewCalendar(ctx, token)
		if err != nil {
			return nil, fmt.Errorf("failed to create calendar service for account: %s, error: %w", userAccount.Email, err)
		}

		now := time.Now()
		startTime := now.AddDate(0, -2, 0)
		endTime := now.AddDate(1, 0, 0)

		calendars, err := em.calendarRepo.FilterByAccountID(ctx, nil, userAccount.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get calendars from db for account: %s, error: %w", userAccount.Email, err)
		}

		events, err := em.fetchEventsFromCalendars(calendarService, calendars, startTime, endTime)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch events for account: %s, error: %w", userAccount.Email, err)
		}

		accountsEvents = append(accountsEvents, &models.AccountsEvents{
			AccountID: userAccount.ID,
			Email:     userAccount.Email,
			Events:    events,
		})
	}

	return accountsEvents, nil
}

func (em *EventManager) fetchEventsFromCalendars(calendarService *Calendar, calendars []*ent.Calendar, startTime, endTime time.Time) ([]*models.Event, error) {
	var events []*models.Event

	for _, cal := range calendars {
		calEvents, err := calendarService.FetchEvents(cal.CalendarID, startTime, endTime)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch events from calendar: %s, error: %w", cal.Summary, err)
		}

		events = append(events, calEvents...)
	}

	return events, nil
}

func (em *EventManager) FetchDraftedEvents(ctx context.Context, userID, accountID uuid.UUID, email string) ([]*models.EventDraftDetail, error) {
	tx, err := em.client.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed starting transaction: %w", err)
	}

	defer transaction.HandleTransaction(tx, &err)

	isPrimary := true
	entCalendar, err := em.calendarRepo.FindByFields(ctx, tx, accountID, nil, nil, &isPrimary)
	if err != nil {
		return nil, fmt.Errorf("failed to get primary calendar for account: %s, error: %w", email, err)
	}

	entEvents, err := em.eventRepo.FilterByCalendarID(ctx, tx, entCalendar.CalendarID)
	if err != nil {
		return nil, fmt.Errorf("failed to get events for account: %s, error: %w", email, err)
	}

	var events []*models.EventDraftDetail
	for _, entEvent := range entEvents {

		events = append(events, &models.EventDraftDetail{
			ID:            entEvent.ID,
			Title:         entEvent.Summary,
			Location:      entEvent.Location,
			Description:   entEvent.Description,
		})
	}

	return events, nil
}

func (em *EventManager) FetchDraftedEventDetail(ctx context.Context, userID, accountID uuid.UUID, email string, eventID uuid.UUID) (*models.EventDraftDetail, error) {
	tx, err := em.client.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed starting transaction: %w", err)
	}

	defer transaction.HandleTransaction(tx, &err)

	entEvent, err := em.eventRepo.Read(ctx, tx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get event for account: %s, error: %w", email, err)
	}

	entDates, err := em.dateRepo.FilterByEventID(ctx, tx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get proposed dates for account: %s, error: %w", email, err)
	}

	var proposedDates []models.ProposedDate
	for _, entDate := range entDates {
		proposedDates = append(proposedDates, models.ProposedDate{
			ID:          entDate.ID,
			EventID:     entDate.GoogleEventID,
			Start:       &entDate.StartTime,
			End:         &entDate.EndTime,
			Priority:    entDate.Priority,
			IsFinalized: entDate.IsFinalized,
		})
	}

	return &models.EventDraftDetail{
		ID:            entEvent.ID,
		Title:         entEvent.Summary,
		Location:      entEvent.Location,
		Description:   entEvent.Description,
		ProposedDates: proposedDates,
	}, nil
}

func (em *EventManager) CreateDraftedEvents(ctx context.Context, userID, accountID uuid.UUID, email string, eventReq *models.EventDraftCreation) error {

	tx, err := em.client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed starting transaction: %w", err)
	}

	token, err := em.authManager.VerifyOAuthToken(ctx, userID, email)
	if err != nil {
		return fmt.Errorf("failed to verify token for account: %s, error: %w", email, err)
	}

	calendarService, err := NewCalendar(ctx, token)
	if err != nil {
		return fmt.Errorf("failed to create calendar service for account: %s, error: %w", email, err)
	}

	defer transaction.HandleTransaction(tx, &err)

	isPrimary := true
	entCalendar, err := em.calendarRepo.FindByFields(ctx, tx, accountID, nil, nil, &isPrimary)
	if err != nil {
		return fmt.Errorf("failed to get primary calendar for account: %s, error: %w", email, err)
	}

	convEvent := em.convertToCalendarEvent(eventReq.Title, eventReq.Location, eventReq.Description, eventReq.SelectedDates[0].Start, eventReq.SelectedDates[0].End)

	entEvent, err := em.eventRepo.Create(ctx, tx, convEvent, entCalendar)
	if err != nil {
		return fmt.Errorf("failed to get event for account: %s, error: %w", email, err)
	}

	insertedGoogleEvents, err := em.createGoogleEvents(calendarService, eventReq)
	if err != nil {
		return fmt.Errorf("failed to insert events for account: %s, error: %w", email, err)
	}

	_, err = em.dateRepo.CreateBulk(ctx, tx, eventReq.SelectedDates , insertedGoogleEvents, entEvent)
	if err != nil {
		if delErr := em.deleteGoogleEvents(calendarService, insertedGoogleEvents); delErr != nil {
			return fmt.Errorf("failed to delete events from Google Calendar: %w", delErr)
		}
		return fmt.Errorf("failed to create proposed dates for account: %s, error: %w", email, err)
	}

	// トランザクションをコミット
	return nil
}

func (em *EventManager) UpdateDraftedEvents(ctx context.Context, userID, accountID uuid.UUID, email string, eventReq *models.EventDraftDetail) error {
	tx, err := em.client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed starting transaction: %w", err)
	}

	// OAuthトークンを検証
	token, err := em.authManager.VerifyOAuthToken(ctx, userID, email)
	if err != nil {
		return fmt.Errorf("failed to verify token for account: %s, error: %w", email, err)
	}

	// Google Calendarサービスを作成
	calendarService, err := NewCalendar(ctx, token)
	if err != nil {
		return fmt.Errorf("failed to create calendar service for account: %s, error: %w", email, err)
	}

	// トランザクションをデファーで処理
	defer transaction.HandleTransaction(tx, &err)

	// プライマリカレンダーを取得
	isPrimary := true
	_, err = em.calendarRepo.FindByFields(ctx, tx, accountID, nil, nil, &isPrimary)
	if err != nil {
		return fmt.Errorf("failed to get primary calendar for account: %s, error: %w", email, err)
	}

	// calendar.Event型に変換
	convEvent := em.convertToCalendarEvent(eventReq.Title, eventReq.Location, eventReq.Description, time.Time{}, time.Time{})

	// イベントの詳細を更新
	entEvent, err := em.eventRepo.Update(ctx, tx, eventReq.ID, convEvent)
	if err != nil {
		return fmt.Errorf("failed to get event for account: %s, error: %w", email, err)
	}

	// 日程候補の変更処理
	err = em.updateProposedDates(ctx, tx, eventReq, entEvent)
	if err != nil {
		return fmt.Errorf("failed to update proposed dates for account: %s, error: %w", email, err)
	}

	// Googleカレンダーのイベントを更新
	err = em.updateGoogleCalendarEvents(calendarService, eventReq)
	if err != nil {
		return fmt.Errorf("failed to update events for account: %s, error: %w", email, err)
	}

	// トランザクションをコミット
	return nil
}

func (em *EventManager) convertToCalendarEvent(title, location, description string, start, end time.Time) *calendar.Event {
	return &calendar.Event{
		Summary:     title,
		Location:    location,
		Description: description,
		Start: &calendar.EventDateTime{
			DateTime: start.Format(time.RFC3339),
			TimeZone: "Asia/Tokyo",
		},
		End: &calendar.EventDateTime{
			DateTime: end.Format(time.RFC3339),
			TimeZone: "Asia/Tokyo",
		},
	}
}

func (em *EventManager) updateProposedDates(ctx context.Context, tx *ent.Tx, eventReq *models.EventDraftDetail, entEvent *ent.Event) error {
	// 日程候補を取得
	existingDates, err := em.dateRepo.FilterByEventID(ctx, tx, eventReq.ID)
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
			_, err = em.dateRepo.Update(ctx, tx, date.ID, &updateDate.EventID, updateDate.Start, updateDate.End, &updateDate.Priority, &updateDate.IsFinalized)
			if err != nil {
				return fmt.Errorf("failed to update proposed dates, error: %w", err)
			}
			// 更新した日程候補を削除
			delete(updateDateMap, date.ID)
		} else {
			err = em.dateRepo.Delete(ctx, tx, date.ID)
			if err != nil {
				return fmt.Errorf("failed to delete proposed dates, error: %w", err)
			}
		}
	}

	// DBに存在しない日程候補を追加
	for _, date := range updateDateMap {
		_, err = em.dateRepo.Create(ctx, tx, &date.EventID, *date.Start, *date.End, date.Priority, entEvent)
		if err != nil {
			return fmt.Errorf("failed to create proposed dates, error: %w", err)
		}
	}

	return nil
}

func (em *EventManager) createGoogleEvents(calendarService *Calendar, eventReq *models.EventDraftCreation) ([]*calendar.Event, error) {
	// Googleカレンダーに登録されたイベントを追跡するスライス
	insertedGoogleEvents := make([]*calendar.Event, len(eventReq.SelectedDates))

	// 並列処理でイベントを登録
	var wg sync.WaitGroup
	var mu sync.Mutex
	errCh := make(chan error, len(eventReq.SelectedDates))

	wg.Add(len(eventReq.SelectedDates))
	for i, date := range eventReq.SelectedDates {
		go func(i int, date models.SelectedDate) {
			defer wg.Done()

			event := em.convertToCalendarEvent(eventReq.Title, eventReq.Location, eventReq.Description, date.Start, date.End)

			insertEvent, err := calendarService.InsertEvent(event)
			if err != nil {
				errCh <- fmt.Errorf("failed to insert event to Google Calendar: %w", err)
				return
			}

			mu.Lock()
			insertedGoogleEvents[i] = insertEvent
			mu.Unlock()
		}(i, date)
	}

	wg.Wait()
	close(errCh)

	// エラーが発生していた場合、登録したイベントを削除
	for err := range errCh {
		if err != nil {
			if delErr := em.deleteGoogleEvents(calendarService, insertedGoogleEvents); delErr != nil {
				return nil, fmt.Errorf("failed to delete events from Google Calendar: %w", delErr)
			}
			return nil, err
		}
	}
	return insertedGoogleEvents, nil
}

func (em *EventManager) updateGoogleCalendarEvents(calendarService *Calendar, eventReq *models.EventDraftDetail) error {
	// Googleカレンダーに登録されたイベントを追跡するスライス
	var backupGoogleEvents []*calendar.Event

	// 並列処理でイベントを登録
	var wg sync.WaitGroup
	var mu sync.Mutex
	errCh := make(chan error, len(eventReq.ProposedDates))

	wg.Add(len(eventReq.ProposedDates))
	for _, date := range eventReq.ProposedDates {
		go func(date models.ProposedDate) {
			defer wg.Done()

			event := em.convertToCalendarEvent(eventReq.Title, eventReq.Location, eventReq.Description, *date.Start, *date.End)

			// 更新前にイベントをバックアップできるように、Googleカレンダーからイベントを取得
			backupEvent, err := calendarService.FetchEvent(date.EventID)
			if err != nil {
				errCh <- fmt.Errorf("failed to fetch event from Google Calendar: %w", err)
				return
			}

			_, err = calendarService.UpdateEvent(date.EventID, event)
			if err != nil {
				errCh <- fmt.Errorf("failed to insert event to Google Calendar: %w", err)
				return
			}

			mu.Lock()
			backupGoogleEvents = append(backupGoogleEvents, backupEvent)
			mu.Unlock()
		}(date)
	}

	wg.Wait()
	close(errCh)

	// エラーが発生していた場合、更新したイベントを元に戻す
	for err := range errCh {
		if err != nil {
			for _, event := range backupGoogleEvents {
				if _, err := calendarService.UpdateEvent(event.Id, event); err != nil {
					return fmt.Errorf("failed to update event to Google Calendar: %w", err)
				}
			}
			return err
		}
	}

	return nil
}

func (em *EventManager) deleteGoogleEvents(calendarService *Calendar, googleEvents []*calendar.Event) error {
	for _, event := range googleEvents {
		err := calendarService.DeleteEvent(event.Id)
		if err != nil {
			return fmt.Errorf("failed to delete event from Google Calendar: %w", err)
		}
	}
	return nil
}