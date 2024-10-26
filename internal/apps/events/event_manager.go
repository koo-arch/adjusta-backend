package events

import (
	"github.com/koo-arch/adjusta-backend/ent"
	appCalendar "github.com/koo-arch/adjusta-backend/internal/apps/calendar"
	"github.com/koo-arch/adjusta-backend/internal/auth"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/repo/calendar"
	"github.com/koo-arch/adjusta-backend/internal/repo/event"
	"github.com/koo-arch/adjusta-backend/internal/repo/proposeddate"
)

type EventManager struct {
	Client       *ent.Client
	AuthManager  *auth.AuthManager
	CalendarRepo repoCalendar.CalendarRepository
	EventRepo    event.EventRepository
	DateRepo     proposeddate.ProposedDateRepository
	CalendarApp  *appCalendar.GoogleCalendarManager
}

func NewEventManager(
	client *ent.Client,
	authManager *auth.AuthManager,
	calendarRepo repoCalendar.CalendarRepository,
	eventRepo event.EventRepository,
	dateRepo proposeddate.ProposedDateRepository,
	calendarApp *appCalendar.GoogleCalendarManager,
) *EventManager {
	return &EventManager{
		Client:       client,
		AuthManager:  authManager,
		CalendarRepo: calendarRepo,
		EventRepo:    eventRepo,
		DateRepo:     dateRepo,
		CalendarApp:  calendarApp,
	}
}
