package api

import (
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/cache"
	appCalendar "github.com/koo-arch/adjusta-backend/internal/apps/calendar"
	appEvents "github.com/koo-arch/adjusta-backend/internal/apps/events"
	"github.com/koo-arch/adjusta-backend/internal/apps/events/event_operations"
	"github.com/koo-arch/adjusta-backend/internal/auth"
	"github.com/koo-arch/adjusta-backend/internal/repo/account"
	dbCalendar "github.com/koo-arch/adjusta-backend/internal/repo/calendar"
	"github.com/koo-arch/adjusta-backend/internal/repo/event"
	"github.com/koo-arch/adjusta-backend/internal/repo/proposeddate"
	"github.com/koo-arch/adjusta-backend/internal/repo/user"
)

type Server struct {
	Client 	      			 *ent.Client
	Cache 	   			 	 *cache.Cache
	UserRepo      			 user.UserRepository
	AccountRepo   			 account.AccountRepository
	CalendarRepo  			 dbCalendar.CalendarRepository
	EventRepo     			 event.EventRepository
	DateRepo      			 proposeddate.ProposedDateRepository
	AuthManager   			 *auth.AuthManager
	JWTManager    			 *auth.JWTManager
	KeyManager    			 *auth.KeyManager
	EventCreationManager     *event_operations.EventCreationManager
	EventFetchingManager     *event_operations.EventFetchingManager
	EventFinalizationManager *event_operations.EventFinalizationManager
	EventUpdateManager       *event_operations.EventUpdateManager
	EventDeleteManager 	 	 *event_operations.EventDeleteManager
}

func NewServer(client *ent.Client) *Server {
	cache := cache.NewCache()

	userRepo := user.NewUserRepository(client)
	accountRepo := account.NewAccountRepository(client)
	calendarRepo := dbCalendar.NewCalendarRepository(client)
	eventRepo := event.NewEventRepository(client)
	dateRepo := proposeddate.NewProposedDateRepository(client)
	calendarApp := appCalendar.NewGoogleCalendarManager(client) // Google Calendar API manager
	keyManager := auth.NewKeyManager(client, cache)
	jwtManager := auth.NewJWTManager(client, keyManager)
	authManager := auth.NewAuthManager(client, userRepo, accountRepo)

	eventManager := appEvents.NewEventManager(client, authManager, calendarRepo, eventRepo, dateRepo, calendarApp)

	return &Server{
		Client:        			  client,
		Cache:         			  cache,
		UserRepo:      			  userRepo,
		AccountRepo:   			  accountRepo,
		CalendarRepo:  			  calendarRepo,
		EventRepo:     			  eventRepo,
		DateRepo:      			  dateRepo,
		AuthManager:   			  authManager,
		JWTManager:    			  jwtManager,
		KeyManager:    			  keyManager,
		EventCreationManager:     event_operations.NewEventCreationManager(eventManager),
		EventFetchingManager:     event_operations.NewEventFetchingManager(eventManager),
		EventFinalizationManager: event_operations.NewEventFinalizationManager(eventManager),
		EventUpdateManager:       event_operations.NewEventUpdateManager(eventManager),
		EventDeleteManager:       event_operations.NewEventDeleteManager(eventManager),
	}
}