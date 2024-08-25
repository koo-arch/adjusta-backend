package calendar

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/api/calendar/v3"
	"github.com/koo-arch/adjusta-backend/ent"
	dbCalendar "github.com/koo-arch/adjusta-backend/internal/apps/calendar"
	"github.com/koo-arch/adjusta-backend/internal/apps/event"
	"github.com/koo-arch/adjusta-backend/internal/auth"
	"github.com/koo-arch/adjusta-backend/internal/models"
	"github.com/koo-arch/adjusta-backend/internal/apps/proposeddate"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
)

type EventManager struct {
	client *ent.Client
	authManager *auth.AuthManager
	calendarRepo dbCalendar.CalendarRepository
	eventRepo event.EventRepository
	dateRepo proposeddate.ProposedDateRepository
}

func NewEventManager(client *ent.Client, authManager *auth.AuthManager, calendarRepo dbCalendar.CalendarRepository, eventRepo event.EventRepository, dateRepo proposeddate.ProposedDateRepository) *EventManager {
	return &EventManager{
		client: client,
		authManager: authManager,
		calendarRepo: calendarRepo,
		eventRepo: eventRepo,
		dateRepo: dateRepo,
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

func (em *EventManager) CreateDraftedEvents(ctx context.Context, userID, accountID uuid.UUID, email string, eventReq *models.EventDraft) error {

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

	_, err = em.dateRepo.CreateBulk(ctx, tx, eventReq.SelectedDates, entEvent)
	if err != nil {
		return fmt.Errorf("failed to create proposed dates for account: %s, error: %w", email, err)
	}
	 // Googleカレンダーに登録されたイベントを追跡するスライス
    var insertedGoogleEvents []*calendar.Event

	// 並列処理でイベントを登録
	var wg sync.WaitGroup
	var mu sync.Mutex
	errCh := make(chan error, len(eventReq.SelectedDates))
	
	wg.Add(len(eventReq.SelectedDates))
	for _, date := range eventReq.SelectedDates {
		go func(date models.SelectedDate) {
			defer wg.Done()

			event := em.convertToCalendarEvent(eventReq.Title, eventReq.Location, eventReq.Description, date.Start, date.End)

			insertEvent, err := calendarService.InsertEvent(event)
			if err != nil {
				errCh <- fmt.Errorf("failed to insert event to Google Calendar: %w", err)
				return
			}

			mu.Lock()
			insertedGoogleEvents = append(insertedGoogleEvents, insertEvent)
			mu.Unlock()
		}(date)
	}

	wg.Wait()
	close(errCh)

	// エラーが発生していた場合、登録したイベントを削除
	for err := range errCh {
		if err != nil {
			for _, event := range insertedGoogleEvents {
				if err := calendarService.DeleteEvent(event.Id); err != nil {
					return fmt.Errorf("failed to delete event from Google Calendar: %w", err)
				}
			}
			// トランザクションをロールバック
			return err
		}
	}

	// トランザクションをコミット
	return nil
}

func (em *EventManager) convertToCalendarEvent (title, location, description string, start, end time.Time) *calendar.Event {
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