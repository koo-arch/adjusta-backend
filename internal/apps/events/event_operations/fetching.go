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

func (efm *EventFetchingManager) FetchAllEvents(ctx context.Context, userID uuid.UUID, userAccounts []*ent.Account) ([]*models.AccountsEvents, error) {
	var accountsEvents []*models.AccountsEvents

	for _, userAccount := range userAccounts {
		token, err := efm.event.AuthManager.VerifyOAuthToken(ctx, userID, userAccount.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to verify token for account: %s, error: %w", userAccount.Email, err)
		}

		calendarService, err := customCalendar.NewCalendar(ctx, token)
		if err != nil {
			return nil, fmt.Errorf("failed to create calendar service for account: %s, error: %w", userAccount.Email, err)
		}

		now := time.Now()
		startTime := now.AddDate(0, -2, 0)
		endTime := now.AddDate(1, 0, 0)

		calendars, err := efm.event.CalendarRepo.FilterByAccountID(ctx, nil, userAccount.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get calendars from db for account: %s, error: %w", userAccount.Email, err)
		}

		events, err := efm.event.CalendarApp.FetchEventsFromCalendars(calendarService, calendars, startTime, endTime)
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

func (efm *EventFetchingManager) FetchDraftedEvents(ctx context.Context, userID, accountID uuid.UUID, email string) ([]*models.EventDraftDetail, error) {
	isPrimary := true
	findOptions := dbCalendar.CalendarQueryOptions{
		IsPrimary:         &isPrimary,
		WithEvents:        true,
		WithProposedDates: true,
	}
	entCalendar, err := efm.event.CalendarRepo.FindByFields(ctx, nil, accountID, findOptions)
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

func (efm *EventFetchingManager) FetchDraftedEventDetail(ctx context.Context, userID, accountID uuid.UUID, email string, eventID uuid.UUID) (*models.EventDraftDetail, error) {
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