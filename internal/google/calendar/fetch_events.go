package calendar

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/auth"
	"github.com/koo-arch/adjusta-backend/internal/apps/calendar"
)

type AccountsEvents struct {
	AccountID uuid.UUID   `json:"account_id"`
	Email     string      `json:"email"`
	Events    []*Event    `json:"events"`
}

func FetchAllEvents(ctx context.Context, authManager *auth.AuthManager, userID uuid.UUID, userAccounts []*ent.Account, calendarRepo calendar.CalendarRepository) ([]*AccountsEvents, error) {
	var accountsEvents []*AccountsEvents

	for _, userAccount := range userAccounts {
		token, err := authManager.VerifyOAuthToken(ctx, userID, userAccount.Email)
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

		calendars, err := calendarRepo.FilterByAccountID(ctx, nil, userAccount.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get calendars from db for account: %s, error: %w", userAccount.Email, err)
		}

		events, err := fetchEventsFromCalendars(calendarService, calendars, startTime, endTime)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch events for account: %s, error: %w", userAccount.Email, err)
		}

		accountsEvents = append(accountsEvents, &AccountsEvents{
			AccountID: userAccount.ID,
			Email:     userAccount.Email,
			Events:    events,
		})
	}

	return accountsEvents, nil
}

func fetchEventsFromCalendars(calendarService *Calendar, calendars []*ent.Calendar, startTime, endTime time.Time) ([]*Event, error) {
	var events []*Event

	for _, cal := range calendars {
		calEvents, err := calendarService.FetchEvents(cal.CalendarID, startTime, endTime)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch events from calendar: %s, error: %w", cal.Summary, err)
		}

		events = append(events, calEvents...)
	}

	return events, nil
}