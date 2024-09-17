package calendar

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/auth"
	"github.com/koo-arch/adjusta-backend/internal/models"
	dbCalendar "github.com/koo-arch/adjusta-backend/internal/repo/calendar"
	"github.com/koo-arch/adjusta-backend/internal/repo/event"
	"github.com/koo-arch/adjusta-backend/internal/repo/proposeddate"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/googleapi"
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
	isPrimary := true
	findOptions := dbCalendar.CalendarQueryOptions{
		IsPrimary:         &isPrimary,
		WithEvents:        true,
		WithProposedDates: true,
	}
	entCalendar, err := em.calendarRepo.FindByFields(ctx, nil, accountID, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to get primary calendar for account: %s, error: %w", email, err)
	}

	if entCalendar.Edges.Events == nil {
		return nil, fmt.Errorf("failed to get events for account: %s", email)
	}

	var draftedEvents []*models.EventDraftDetail
	for _, entEvent := range entCalendar.Edges.Events {
		var proposedDates []models.ProposedDate

		if entEvent.Edges.ProposedDates == nil {
			continue
		}
		for _, entDate := range entEvent.Edges.ProposedDates {
			proposedDates = append(proposedDates, models.ProposedDate{
				ID:            entDate.ID,
				GoogleEventID: entDate.GoogleEventID,
				Start:         &entDate.StartTime,
				End:           &entDate.EndTime,
				Priority:      entDate.Priority,
				IsFinalized:   entDate.IsFinalized,
			})
		}

		draftedEvents = append(draftedEvents, &models.EventDraftDetail{
			ID:            entEvent.ID,
			Title:         entEvent.Summary,
			Location:      entEvent.Location,
			Description:   entEvent.Description,
			ProposedDates: proposedDates,
		})
	}

	return draftedEvents, nil
}

func (em *EventManager) FetchDraftedEventDetail(ctx context.Context, userID, accountID uuid.UUID, email string, eventID uuid.UUID) (*models.EventDraftDetail, error) {
	tx, err := em.client.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed starting transaction: %w", err)
	}

	defer transaction.HandleTransaction(tx, &err)

	queryOpt := event.EventQueryOptions{
		WithProposedDates: true,
	}
	entEvent, err := em.eventRepo.Read(ctx, tx, eventID, queryOpt)
	if err != nil {
		return nil, fmt.Errorf("failed to get event for account: %s, error: %w", email, err)
	}

	if entEvent.Edges.ProposedDates == nil {
		return nil, fmt.Errorf("failed to get proposed dates for account: %s", email)
	}

	var proposedDates []models.ProposedDate
	for _, entDate := range entEvent.Edges.ProposedDates {
		proposedDates = append(proposedDates, models.ProposedDate{
			ID:            entDate.ID,
			GoogleEventID: entDate.GoogleEventID,
			Start:         &entDate.StartTime,
			End:           &entDate.EndTime,
			Priority:      entDate.Priority,
			IsFinalized:   entDate.IsFinalized,
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
	findOptions := dbCalendar.CalendarQueryOptions{
		IsPrimary: &isPrimary,
	}
	entCalendar, err := em.calendarRepo.FindByFields(ctx, tx, accountID, findOptions)
	if err != nil {
		return fmt.Errorf("failed to get primary calendar for account: %s, error: %w", email, err)
	}

	convEvent := em.convertToCalendarEvent(nil, eventReq.Title, eventReq.Location, eventReq.Description, eventReq.SelectedDates[0].Start, eventReq.SelectedDates[0].End)

	entEvent, err := em.eventRepo.Create(ctx, tx, convEvent, entCalendar)
	if err != nil {
		return fmt.Errorf("failed to get event for account: %s, error: %w", email, err)
	}

	insertedGoogleEvents, err := em.createGoogleEvents(calendarService, eventReq)
	if err != nil {
		return fmt.Errorf("failed to insert events for account: %s, error: %w", email, err)
	}

	_, err = em.dateRepo.CreateBulk(ctx, tx, eventReq.SelectedDates, insertedGoogleEvents, entEvent)
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
	findOptions := dbCalendar.CalendarQueryOptions{
		IsPrimary: &isPrimary,
	}
	_, err = em.calendarRepo.FindByFields(ctx, tx, accountID, findOptions)
	if err != nil {
		return fmt.Errorf("failed to get primary calendar for account: %s, error: %w", email, err)
	}

	// calendar.Event型に変換
	convEvent := em.convertToCalendarEvent(nil, eventReq.Title, eventReq.Location, eventReq.Description, time.Time{}, time.Time{})

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

func (em *EventManager) FinalizeProposedDate(ctx context.Context, userID, accountID, eventID uuid.UUID, email string, eventReq *models.ConfrimEvent) error {
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
	findOptions := dbCalendar.CalendarQueryOptions{
		IsPrimary: &isPrimary,
	}
	_, err = em.calendarRepo.FindByFields(ctx, tx, accountID, findOptions)
	if err != nil {
		return fmt.Errorf("failed to get primary calendar for account: %s, error: %w", email, err)
	}

	// イベントの日程候補のis_Finalizedを全てfalseに更新
	err = em.dateRepo.ResetFinalized(ctx, tx, eventID)
	if err != nil {
		return fmt.Errorf("failed to reset is_finalized for account: %s, error: %w", email, err)
	}

	entEvent, err := em.eventRepo.Read(ctx, tx, eventID, event.EventQueryOptions{})
	if err != nil {
		return fmt.Errorf("failed to get event for account: %s, error: %w", email, err)
	}
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
		googleEvents, err := em.createGoogleEvents(calendarService, &eventDraftCreate)
		if err != nil {
			return fmt.Errorf("failed to insert events for account: %s, error: %w", email, err)
		}
		googleEventID = &googleEvents[0].Id

	} else {
		// 既存のGoogleカレンダーイベントIDを使用
		googleEventID = &eventReq.ConfirmDate.GoogleEventID
	}

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

	if eventReq.ConfirmDate.ID == nil {
		_, err = em.dateRepo.Create(ctx, tx, googleEventID, dateOptions, entEvent)
		if err != nil {
			return fmt.Errorf("failed to create proposed date for account: %s, error: %w", email, err)
		}
	} else {
		// 日程候補のis_Finalizedを更新
		entDate, err := em.dateRepo.Read(ctx, tx, *eventReq.ConfirmDate.ID)
		if err != nil {
			return fmt.Errorf("failed to get proposed date for account: %s, error: %w", email, err)
		}

		_, err = em.dateRepo.Update(ctx, tx, entDate.ID, dateOptions)
		if err != nil {
			return fmt.Errorf("failed to update proposed date for account: %s, error: %w", email, err)
		}
	}

	// is_finalizedがfalseの日程候補を検索
	notFinalizedDates, err := em.dateRepo.FilterByEventIDWithFinalized(ctx, tx, eventID, false)
	if err != nil {
		return fmt.Errorf("failed to get not finalized proposed dates for account: %s, error: %w", email, err)
	}

	fmt.Printf("notFinalizedDates: %v\n", notFinalizedDates)

	convEvents := make([]*calendar.Event, len(notFinalizedDates))
	for i, date := range notFinalizedDates {
		convEvents[i] = em.convertToCalendarEvent(&date.GoogleEventID, "", "", "", date.StartTime, date.EndTime)
	}
	// Googleカレンダーのイベントを削除
	err = em.deleteGoogleEvents(calendarService, convEvents)
	if err != nil {
		return fmt.Errorf("failed to delete events for account: %s, error: %w", email, err)
	}

	empty := ""
	// Googleイベントを削除した日程候補のgoogle_event_idを削除
	for _, date := range notFinalizedDates {
		dateOptions := proposeddate.ProposedDateQueryOptions{
			GoogleEventID: &empty,
		}
		_, err = em.dateRepo.Update(ctx, tx, date.ID, dateOptions)
		if err != nil {
			return fmt.Errorf("failed to update proposed date for account: %s, error: %w", email, err)
		}
	}

	// トランザクションをコミット
	return nil
}

func (em *EventManager) convertToCalendarEvent(ID *string, title, location, description string, start, end time.Time) *calendar.Event {
	event := &calendar.Event{
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
	// IDがnilでない場合のみ設定
	if ID != nil && *ID != "" {
		event.Id = *ID
	}

	return event
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

			dateOptions := proposeddate.ProposedDateQueryOptions{
				GoogleEventID: &updateDate.GoogleEventID,
				StartTime:     updateDate.Start,
				EndTime:       updateDate.End,
				Priority:      &updateDate.Priority,
				IsFinalized:   &updateDate.IsFinalized,
			}
			_, err = em.dateRepo.Update(ctx, tx, date.ID, dateOptions)
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
		dateOptions := proposeddate.ProposedDateQueryOptions{
			StartTime: date.Start,
			EndTime:   date.End,
			Priority:  &date.Priority,
		}
		_, err = em.dateRepo.Create(ctx, tx, &date.GoogleEventID, dateOptions, entEvent)
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(len(eventReq.SelectedDates))
	for i, date := range eventReq.SelectedDates {
		go func(i int, date models.SelectedDate) {
			defer wg.Done()

			select {
			case <-ctx.Done(): // エラー発生時に他のゴルーチンをキャンセル
				return
			default:
				event := em.convertToCalendarEvent(nil, eventReq.Title, eventReq.Location, eventReq.Description, date.Start, date.End)

				insertEvent, err := calendarService.InsertEvent(event)
				if err != nil {
					errCh <- fmt.Errorf("failed to insert event to Google Calendar: %w", err)
					cancel() // エラーが発生したら他のゴルーチンをキャンセル
					return
				}

				mu.Lock()
				insertedGoogleEvents[i] = insertEvent
				mu.Unlock()
			}
		}(i, date)
	}

	wg.Wait()
	close(errCh)

	// エラーが発生していた場合、登録したイベントを削除
	var errList []error
	for err := range errCh {
		errList = append(errList, err)
	}

	if len(errList) > 0 {
		if delErr := em.deleteGoogleEvents(calendarService, insertedGoogleEvents); delErr != nil {
			return nil, fmt.Errorf("failed to delete events from Google Calendar: %w", delErr)
		}
		return nil, fmt.Errorf("multiple errors occurred: %v", errList)
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

			event := em.convertToCalendarEvent(&date.GoogleEventID, eventReq.Title, eventReq.Location, eventReq.Description, *date.Start, *date.End)

			// 更新前にイベントをバックアップできるように、Googleカレンダーからイベントを取得
			backupEvent, err := calendarService.FetchEvent(date.GoogleEventID)
			if err != nil {
				errCh <- fmt.Errorf("failed to fetch event from Google Calendar: %w", err)
				return
			}

			_, err = calendarService.UpdateEvent(date.GoogleEventID, event)
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

func (em *EventManager) deleteGoogleEvents(calendarService *Calendar, events []*calendar.Event) error {
	for _, event := range events {
		if event == nil || event.Id == "" {
			continue // eventまたはIDがnilの場合はスキップ
		}

		err := calendarService.DeleteEvent(event.Id)
		if err != nil {
			if gErr, ok := err.(*googleapi.Error); ok && gErr.Code == 410 {
				// 410エラーはリソースが既に削除されているため、無視
				fmt.Printf("Warning: Event ID %s is already deleted.\n", event.Id)
				continue
			}
			return fmt.Errorf("failed to delete event with ID %s: %w", event.Id, err)
		}
	}
	return nil
}
