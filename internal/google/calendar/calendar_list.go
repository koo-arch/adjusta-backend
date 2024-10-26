package calendar

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/auth"
	"github.com/koo-arch/adjusta-backend/internal/repo/calendar"
)

type AccountsCalendars struct {
	AccountID uuid.UUID       `json:"account_id"`
	Email     string          `json:"email"`
	Calendars []*CalendarList `json:"calendars"`
}

func RegisterCalendarList(ctx context.Context, authManager *auth.AuthManager, userID uuid.UUID, userAccounts []*ent.Account, calendarRepo calendar.CalendarRepository) ([]*AccountsCalendars, error) {
	var accountsCalendars []*AccountsCalendars
	var wg sync.WaitGroup
	var mu sync.Mutex
	errCh := make(chan error, len(userAccounts))

	for _, userAccount := range userAccounts {
		wg.Add(1)
		go func(userAccount *ent.Account) {
			defer wg.Done()

			token, err := authManager.VerifyOAuthToken(ctx, userID, userAccount.Email)
			if err != nil {
				errCh <- fmt.Errorf("failed to verify token for account: %s, error: %w", userAccount.Email, err)
				return
			}

			calendarService, err := NewCalendar(ctx, token)
			if err != nil {
				errCh <- fmt.Errorf("failed to create calendar service for account: %s, error: %w", userAccount.Email, err)
				return
			}

			calendars, err := calendarService.FetchCalendarList()
			if err != nil {
				errCh <- fmt.Errorf("failed to fetch calendars for account: %s, error: %w", userAccount.Email, err)
				return
			}

			if err := syncCalendar(ctx, calendars, userAccount, calendarRepo); err != nil {
				errCh <- fmt.Errorf("failed to sync calendars for account: %s, error: %w", userAccount.Email, err)
				return
			}

			mu.Lock() // accountsCalendarsにアクセスするためにロック
			accountsCalendars = append(accountsCalendars, &AccountsCalendars{
				AccountID: userAccount.ID,
				Email:     userAccount.Email,
				Calendars: calendars,
			})
			mu.Unlock()
		}(userAccount)
	}

	// 全てのgoroutineが終了するまで待機
	wg.Wait()
	close(errCh)

	// エラーが発生した場合はエラーを返す
	if len(errCh) > 0 {
		var errList []error
		for err := range errCh {
			errList = append(errList, err)
		}
		return nil, fmt.Errorf("multiple errors occurred: %v", errList)
	}

	return accountsCalendars, nil
}

func syncCalendar(ctx context.Context, calendars []*CalendarList, userAccount *ent.Account, calendarRepo calendar.CalendarRepository) error {
	repoCalendars, err := calendarRepo.FilterByAccountID(ctx, nil, userAccount.ID)
	if err != nil {
		return fmt.Errorf("failed to get calendars from db for account: %s, error: %w", userAccount.Email, err)
	}

	// Googleから取得したカレンダーとデータベースのカレンダーを比較
	calendarMap := make(map[string]*CalendarList)
	for _, cal := range calendars {
		calendarMap[cal.CalendarID] = cal
	}

	// データベースに存在するカレンダーをマップから削除
	for _, dbCal := range repoCalendars {
		if _, ok := calendarMap[dbCal.CalendarID]; ok {
			delete(calendarMap, dbCal.CalendarID)
		}
	}

	// データベースに存在しないカレンダーを追加
	for _, cal := range calendarMap {
		if _, err := calendarRepo.Create(ctx, nil, cal.CalendarID, cal.Summary, cal.Primary, userAccount); err != nil {
			if strings.Contains(err.Error(), "is already in use by another account of the same user") {
				fmt.Printf("calendar already exists: %s\n", cal.Summary)
				continue
			}
			return fmt.Errorf("failed to insert calendar to google calendar: %s, error: %w", cal.Summary, err)
		}
	}

	return nil
}
