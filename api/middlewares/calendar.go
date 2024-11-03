package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	customCalendar "github.com/koo-arch/adjusta-backend/internal/google/calendar"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/repo/calendar"
	"github.com/koo-arch/adjusta-backend/internal/repo/googlecalendarinfo"
	"github.com/koo-arch/adjusta-backend/internal/repo/user"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
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

		userRepo := cm.middleware.Server.UserRepo
		entUser, err := userRepo.Read(ctx, nil, userid, user.UserQueryOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user from db"})
			c.Abort()
			return
		}

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

		calendarList, err := cm.fetchAndCacheCalendars(ctx, entUser)
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

func(cm *CalendarMiddleware) fetchAndCacheCalendars(ctx context.Context, entUser *ent.User) ([]*customCalendar.CalendarList, error) {

	authManager := cm.middleware.Server.AuthManager
	token, err := authManager.VerifyOAuthToken(ctx, entUser.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token for account: %s, error: %w", entUser.Email, err)
	}

	calendarService, err := customCalendar.NewCalendar(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar service for account: %s, error: %w", entUser.Email, err)
	}

	calendars, err := calendarService.FetchCalendarList()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch calendars for account: %s, error: %w", entUser.Email, err)
	}

	if err := cm.syncCalendar(ctx, calendars, entUser); err != nil {
		return nil, fmt.Errorf("failed to sync calendars for account: %s, error: %w", entUser.Email, err)
	}

	// キャッシュに保存
	calendarCache := cm.middleware.Server.Cache.CalendarCache
	cacheKey := fmt.Sprintf("calendars:%s", entUser.ID.String())
	calendarCache.Set(cacheKey, calendars, 5*time.Hour)

	return calendars, nil
}

func(cm *CalendarMiddleware) syncCalendar(ctx context.Context, calendars []*customCalendar.CalendarList, entUser *ent.User) error {
	// トランザクションを開始
	tx, err := cm.middleware.Server.Client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed starting transaction: %w", err)
	}

	// トランザクションをデファーで処理
	defer transaction.HandleTransaction(tx, &err)

	calendarRepo := cm.middleware.Server.CalendarRepo
	googleCalendarRepo := cm.middleware.Server.GoogleCalendarRepo

	
	// カレンダーがDBに存在するか確認
	for _, cal := range calendars {
		// カレンダーがDBに存在するか確認
		repoCalOptions := repoCalendar.CalendarQueryOptions{
			WithGoogleCalendarInfo: true,
			GoogleCalendarID: &cal.CalendarID,
		}
		entCalendar, err := calendarRepo.FindByFields(ctx, tx, entUser.ID, repoCalOptions)
		if err != nil {
			if !ent.IsNotFound(err) {
				return fmt.Errorf("failed to find calendar: %w", err)
			}
			
			// カレンダーが存在しない場合は新規作成
			entCalendar, err = calendarRepo.Create(ctx, tx, entUser, nil)
			if err != nil {
				return fmt.Errorf("failed to create calendar: %w", err)
			}
		}

		// GoogleCalendarInfoにカレンダーが存在するか確認
		gCalOptions := googlecalendarinfo.GoogleCalendarInfoQueryOptions{
			GoogleCalendarID : &cal.CalendarID,
		}
		entGoogleCalendar, err := googleCalendarRepo.FindByFields(ctx, tx, gCalOptions)
		if err != nil {
			if !ent.IsNotFound(err) {
				return fmt.Errorf("failed to find google calendar info: %w", err)
			}

			// 存在しない場合は新規作成
			gCalOptions := googlecalendarinfo.GoogleCalendarInfoQueryOptions{
				GoogleCalendarID: &cal.CalendarID,
				Summary: &cal.Summary,
				IsPrimary: &cal.Primary,
			}
			_, err := googleCalendarRepo.Create(ctx, tx, gCalOptions, entCalendar)
			if err != nil {
				return fmt.Errorf("failed to create google calendar info: %w", err)
			}
		}

		// カレンダーが存在する場合は関連付け
		if entGoogleCalendar != nil {
			_, err := googleCalendarRepo.Update(ctx, tx, entGoogleCalendar.ID, googlecalendarinfo.GoogleCalendarInfoQueryOptions{}, entCalendar)
			if err != nil {
				return fmt.Errorf("failed to update google calendar info: %w", err)
			}
		}
	}

	// トランザクションをコミット
	return nil
}