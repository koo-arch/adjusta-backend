package calendar

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/auth"
	"github.com/koo-arch/adjusta-backend/internal/apps/calendar"
)

type AccountsCalendars struct {
	AccountID uuid.UUID   `json:"account_id"`
	Email     string      `json:"email"`
	Calendars []*CalendarList `json:"calendars"`
}

func RegisterCalendarList(ctx context.Context, authManager *auth.AuthManager, userID uuid.UUID, userAccounts []*ent.Account, calendarRepo calendar.CalendarRepository) ([]*AccountsCalendars, error) {
	var accountsCalendars []*AccountsCalendars

	for _, userAccount := range userAccounts {
		token, err := authManager.VerifyOAuthToken(ctx, userID, userAccount.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to verify token for account: %s, error: %w", userAccount.Email, err)
		}

		calendarService, err := NewCalendar(ctx, token)
		if err != nil {
			return nil, fmt.Errorf("failed to create calendar service for account: %s, error: %w", userAccount.Email, err)
		}

		calendars, err := calendarService.FetchCalendarList()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch calendars for account: %s, error: %w", userAccount.Email, err)
		}

		if err := syncCalendar(ctx, calendars, userAccount, calendarRepo); err != nil {
			return nil, fmt.Errorf("failed to sync calendars for account: %s, error: %w", userAccount.Email, err)
		}

			
		accountsCalendars = append(accountsCalendars, &AccountsCalendars{
			AccountID: userAccount.ID,
			Email:     userAccount.Email,
			Calendars: calendars,
		})
			
	}

	return accountsCalendars, nil
}

func syncCalendar(ctx context.Context, calendars []*CalendarList ,userAccount *ent.Account,  calendarRepo calendar.CalendarRepository) error {
	dbCalendars, err := calendarRepo.FilterByAccountID(ctx, nil, userAccount.ID)
	if err != nil {
		return fmt.Errorf("failed to get calendars from db for account: %s, error: %w", userAccount.Email, err)
	}

	// Googleから取得したカレンダーとデータベースのカレンダーを比較
	calendarMap := make(map[string]*CalendarList)
	for _, cal := range calendars {
		calendarMap[cal.CalendarID] = cal
	}

	// データベースに存在するカレンダーをマップから削除
	for _, dbCal := range dbCalendars {
		if _, ok := calendarMap[dbCal.CalendarID]; ok {
			delete(calendarMap, dbCal.CalendarID)
		}
	}

	// データベースに存在しないカレンダーを追加
	for _, cal := range calendarMap {
		if _, err := calendarRepo.Create(ctx, nil, cal.CalendarID, cal.Summary, cal.Primary, userAccount); err != nil {
			return fmt.Errorf("failed to insert calendar to google calendar: %s, error: %w", cal.Summary, err)
		}
	}

	return nil
}