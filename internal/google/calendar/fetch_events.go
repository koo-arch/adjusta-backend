package calendar

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/auth"
)

type AccountsEvents struct {
	AccountID uuid.UUID   `json:"account_id"`
	Email     string      `json:"email"`
	Events    []*Event    `json:"events"`
}

func FetchAllEvents(ctx context.Context, authManager *auth.AuthManager, userAccounts []*ent.Account) ([]*AccountsEvents, error) {
	var accountsEvents []*AccountsEvents

	for _, userAccount := range userAccounts {
		token, err := authManager.VerifyOAuthToken(ctx, userAccount.Email)
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

		events, err := calendarService.FetchEvents(startTime, endTime)
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