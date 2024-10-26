package middlewares

import (
	"net/http"
	"context"
	"fmt"
	"time"
	"strings"
	"sync"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/auth"
	customCalendar "github.com/koo-arch/adjusta-backend/internal/google/calendar"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/repo/calendar"
)

type CalendarMiddleware struct {
	middleware *Middleware
}

func NewCalendarMiddleware(middleware *Middleware) *CalendarMiddleware {
	return &CalendarMiddleware{
		middleware: middleware,
	}
}

func (cm *CalendarMiddleware) SyncGoogleCalendars() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		session := sessions.Default(c)
		useridStr, ok := session.Get("userid").(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to get userid from session"})
			c.Abort()
			return
		}

		userid, err := uuid.Parse(useridStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userid format"})
			c.Abort()
			return
		}
		accountRepo := cm.middleware.Server.AccountRepo
		calendarRepo := cm.middleware.Server.CalendarRepo
		authManager := cm.middleware.Server.AuthManager

		// キャッシュにある場合はそれを使う
		cache := cm.middleware.Server.Cache
		cacheKey := fmt.Sprintf("calendars:%s", useridStr)
		if cacheCalendar, found := cache.CalendarCache.Get(cacheKey); found {
			println("use cache")
			c.Set("calendarList", cacheCalendar)
			c.Next()
			c.Abort()
			return
		}

		userAccounts, err := accountRepo.FilterByUserID(ctx, nil, userid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user accounts"})
			c.Abort()
			return
		}

		calendarList, err := cm.fetchAndCacheCalendars(ctx, authManager, userid, userAccounts, calendarRepo)
		if err != nil {
			fmt.Printf("failed to register calendar list: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register calendar list"})
			c.Abort()
			return
		}

		c.Set("calendarList", calendarList)
		c.Next()
	}
}

type AccountsCalendars struct {
	AccountID uuid.UUID       `json:"account_id"`
	Email     string          `json:"email"`
	Calendars []*customCalendar.CalendarList `json:"calendars"`
}

func(cm *CalendarMiddleware) fetchAndCacheCalendars(ctx context.Context, authManager *auth.AuthManager, userID uuid.UUID, userAccounts []*ent.Account, calendarRepo repoCalendar.CalendarRepository) ([]*AccountsCalendars, error) {
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

			calendarService, err := customCalendar.NewCalendar(ctx, token)
			if err != nil {
				errCh <- fmt.Errorf("failed to create calendar service for account: %s, error: %w", userAccount.Email, err)
				return
			}

			calendars, err := calendarService.FetchCalendarList()
			if err != nil {
				errCh <- fmt.Errorf("failed to fetch calendars for account: %s, error: %w", userAccount.Email, err)
				return
			}

			if err := cm.syncCalendar(ctx, calendars, userAccount, calendarRepo); err != nil {
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

	// キャッシュに保存
	calendarCache := cm.middleware.Server.Cache.CalendarCache
	cacheKey := fmt.Sprintf("calendars:%s", userID.String())
	calendarCache.Set(cacheKey, accountsCalendars, 5*time.Hour)

	return accountsCalendars, nil
}

func(cm *CalendarMiddleware) syncCalendar(ctx context.Context, calendars []*customCalendar.CalendarList, userAccount *ent.Account, calendarRepo repoCalendar.CalendarRepository) error {
	dbCalendars, err := calendarRepo.FilterByAccountID(ctx, nil, userAccount.ID)
	if err != nil {
		return fmt.Errorf("failed to get calendars from db for account: %s, error: %w", userAccount.Email, err)
	}

	// Googleから取得したカレンダーとデータベースのカレンダーを比較
	calendarMap := make(map[string]*customCalendar.CalendarList)
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
			if strings.Contains(err.Error(), "is already in use by another account of the same user") {
				fmt.Printf("calendar already exists: %s\n", cal.Summary)
				continue
			}
			return fmt.Errorf("failed to insert calendar to google calendar: %s, error: %w", cal.Summary, err)
		}
	}

	return nil
}