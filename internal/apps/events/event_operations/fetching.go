package event_operations

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"

	"github.com/koo-arch/adjusta-backend/internal/apps/events"
	customCalendar "github.com/koo-arch/adjusta-backend/internal/google/calendar"
	"github.com/koo-arch/adjusta-backend/internal/models"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/repo/calendar"
	"github.com/koo-arch/adjusta-backend/internal/repo/event"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
)

type EventFetchingManager struct {
	event *events.EventManager
}

func NewEventFetchingManager(event *events.EventManager) *EventFetchingManager {
	return &EventFetchingManager{
		event: event,
	}
}

func (efm *EventFetchingManager) FetchAllGoogleEvents(ctx context.Context, userID uuid.UUID, email string) ([]*models.Event, error) {

	token, err := efm.event.AuthManager.VerifyOAuthToken(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token for account: %s, error: %w", email, err)
	}

	calendarService, err := customCalendar.NewCalendar(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar service for account: %s, error: %w", email, err)
	}

	now := time.Now()
	startTime := now.AddDate(0, -2, 0)
	endTime := now.AddDate(1, 0, 0)

	calendarOptions := repoCalendar.CalendarQueryOptions{
		WithGoogleCalendarInfo: true,
	}
	calendars, err := efm.event.CalendarRepo.FilterByFields(ctx, nil, userID, calendarOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to get calendars from db for account: %s, error: %w", email, err)
	}

	var entGoogleCalendars []*ent.GoogleCalendarInfo
	for _, cal := range calendars {
		if cal.Edges.GoogleCalendarInfos != nil {
			entGoogleCalendars = append(entGoogleCalendars, cal.Edges.GoogleCalendarInfos...)
		}
	}

	events, err := efm.event.CalendarApp.FetchEventsFromCalendars(calendarService, entGoogleCalendars, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch events for account: %s, error: %w", email, err)
	}

	return events, nil
}

func (efm *EventFetchingManager) FetchDraftedEvents(ctx context.Context, userID uuid.UUID, email string) ([]*models.EventDraftDetail, error) {
	isPrimary := true
	findOptions := repoCalendar.CalendarQueryOptions{
		IsPrimary:         &isPrimary,
		WithEvents:        true,
		WithProposedDates: true,
	}
	entCalendar, err := efm.event.CalendarRepo.FindByFields(ctx, nil, userID, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to get primary calendar for account: %s, error: %w", email, err)
	}

	if entCalendar.Edges.Events == nil {
		return nil, fmt.Errorf("failed to get events for account: %s", email)
	}

	var draftedEvents []*models.EventDraftDetail
	var wg sync.WaitGroup
	var mu sync.Mutex
	errCh := make(chan error, 1)

	for _, entEvent := range entCalendar.Edges.Events {
		wg.Add(1)

		go func(entEvent *ent.Event) {
			defer wg.Done()
			var proposedDates []models.ProposedDate

			if entEvent.Edges.ProposedDates == nil {
				return
			}
			for _, entDate := range entEvent.Edges.ProposedDates {
				proposedDates = append(proposedDates, models.ProposedDate{
					ID:            &entDate.ID,
					GoogleEventID: entDate.GoogleEventID,
					Start:         &entDate.StartTime,
					End:           &entDate.EndTime,
					Priority:      entDate.Priority,
				})
			}

			// Priorityに基づいてProposedDatesを昇順にソート
			sort.Slice(proposedDates, func(i, j int) bool {
				return proposedDates[i].Priority < proposedDates[j].Priority
			})

			// 同時に書き込むことがないようにミューテックスを使う
			mu.Lock()
			draftedEvents = append(draftedEvents, &models.EventDraftDetail{
				ID:            entEvent.ID,
				Title:         entEvent.Summary,
				Location:      entEvent.Location,
				Description:   entEvent.Description,
				Status:		   models.EventStatus(entEvent.Status),
				ConfirmedDateID: &entEvent.ConfirmedDateID,
				ProposedDates: proposedDates,
			})
			mu.Unlock()
		}(entEvent)
	}

	wg.Wait()
	close(errCh)

	if len(errCh) > 0 {
		return nil, <-errCh
	}

	return draftedEvents, nil
}

func (efm *EventFetchingManager) FetchDraftedEventDetail(ctx context.Context, userID uuid.UUID, email string, eventID uuid.UUID) (*models.EventDraftDetail, error) {
	tx, err := efm.event.Client.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed starting transaction: %w", err)
	}

	defer transaction.HandleTransaction(tx, &err)

	queryOpt := event.EventQueryOptions{
		WithProposedDates: true,
	}
	entEvent, err := efm.event.EventRepo.Read(ctx, tx, eventID, queryOpt)
	if err != nil {
		return nil, fmt.Errorf("failed to get event for account: %s, error: %w", email, err)
	}

	if entEvent.Edges.ProposedDates == nil {
		return nil, fmt.Errorf("failed to get proposed dates for account: %s", email)
	}

	var proposedDates []models.ProposedDate
	for _, entDate := range entEvent.Edges.ProposedDates {
		proposedDates = append(proposedDates, models.ProposedDate{
			ID:            &entDate.ID,
			GoogleEventID: entDate.GoogleEventID,
			Start:         &entDate.StartTime,
			End:           &entDate.EndTime,
			Priority:      entDate.Priority,
		})
	}

	// Priorityに基づいてProposedDatesを昇順にソート
	sort.Slice(proposedDates, func(i, j int) bool {
		return proposedDates[i].Priority < proposedDates[j].Priority
	})

	return &models.EventDraftDetail{
		ID:            entEvent.ID,
		Title:         entEvent.Summary,
		Location:      entEvent.Location,
		Description:   entEvent.Description,
		Status:        models.EventStatus(entEvent.Status),
		ConfirmedDateID: &entEvent.ConfirmedDateID,
		ProposedDates: proposedDates,
	}, nil
}
