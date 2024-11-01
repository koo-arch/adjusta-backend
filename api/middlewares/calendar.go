package middlewares

import (
	"net/http"
	"context"
	"fmt"
	"time"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/repo/user"
	customCalendar "github.com/koo-arch/adjusta-backend/internal/google/calendar"
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
	calendarRepo := cm.middleware.Server.CalendarRepo
	dbCalendars, err := calendarRepo.FilterByUserID(ctx, nil, entUser.ID)
	if err != nil {
		return fmt.Errorf("failed to get calendars from db for user: %s, error: %w", entUser.Email, err)
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
		if _, err := calendarRepo.Create(ctx, nil, cal.CalendarID, cal.Summary, cal.Primary, entUser); err != nil {
			if strings.Contains(err.Error(), "is already in use by another account of the same user") {
				fmt.Printf("calendar already exists: %s\n", cal.Summary)
				continue
			}
			return fmt.Errorf("failed to insert calendar to google calendar: %s, error: %w", cal.Summary, err)
		}
	}

	return nil
}