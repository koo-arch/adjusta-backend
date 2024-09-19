package handlers

import (
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/auth"
	"github.com/koo-arch/adjusta-backend/internal/repo/account"
	dbCalendar "github.com/koo-arch/adjusta-backend/internal/repo/calendar"
	"github.com/koo-arch/adjusta-backend/internal/repo/event"
	"github.com/koo-arch/adjusta-backend/internal/repo/proposeddate"
	"github.com/koo-arch/adjusta-backend/internal/repo/user"
	appEvents "github.com/koo-arch/adjusta-backend/internal/apps/events"
	"github.com/koo-arch/adjusta-backend/internal/apps/events/event_operations"
	appCalendar "github.com/koo-arch/adjusta-backend/internal/apps/calendar"
)

type Server struct {
	client        				*ent.Client
	userRepo      				user.UserRepository
	accountRepo   				account.AccountRepository
	calendarRepo  				dbCalendar.CalendarRepository
	eventRepo     				event.EventRepository
	dateRepo      				proposeddate.ProposedDateRepository
	authManager  		 		*auth.AuthManager
	jwtManager   				*auth.JWTManager
	keyManager					*auth.KeyManager
	eventCreationManager 		*event_operations.EventCreationManager
	eventFetchingManager 		*event_operations.EventFetchingManager
	eventFinalizationManager 	*event_operations.EventFinalizationManager
	eventUpdateManager 			*event_operations.EventUpdateManager
}

// Server の初期化関数
func NewServer(client *ent.Client) *Server {
	userRepo := user.NewUserRepository(client)
	accountRepo := account.NewAccountRepository(client)
	calendarRepo := dbCalendar.NewCalendarRepository(client)
	eventRepo := event.NewEventRepository(client)
	dateRepo := proposeddate.NewProposedDateRepository(client)
	calendarApp := appCalendar.NewGoogleCalendarManager(client) // Google Calendar APIのマネージャー
	keyManager := auth.NewKeyManager(client)
	jwtManager := auth.NewJWTManager(client, keyManager)
	authManager := auth.NewAuthManager(client, userRepo, accountRepo)

	eventManager := appEvents.NewEventManager(client, authManager, calendarRepo, eventRepo, dateRepo, calendarApp)

	return &Server{
		client:       				client,
		userRepo:     				userRepo,
		accountRepo:  				accountRepo,
		calendarRepo: 				calendarRepo,
		eventRepo:    				eventRepo,
		dateRepo:     				dateRepo,
		authManager:  				authManager,
		jwtManager:   			  	jwtManager,
		keyManager:   			  	keyManager,
		eventCreationManager: 	  	event_operations.NewEventCreationManager(eventManager),
		eventFetchingManager: 	  	event_operations.NewEventFetchingManager(eventManager),
		eventFinalizationManager:	event_operations.NewEventFinalizationManager(eventManager),
		eventUpdateManager: 	  	event_operations.NewEventUpdateManager(eventManager),
	}
}