package event_operations

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/apps/events"
	"github.com/koo-arch/adjusta-backend/internal/models"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/repo/calendar"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
)

type EventCreationManager struct {
	event *events.EventManager
}

func NewEventCreationManager(event *events.EventManager) *EventCreationManager {
	return &EventCreationManager{
		event: event,
	}
}

func (ecm *EventCreationManager) CreateDraftedEvents(ctx context.Context, userID uuid.UUID, email string, eventReq *models.EventDraftCreation) error {

	tx, err := ecm.event.Client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed starting transaction: %w", err)
	}

	defer transaction.HandleTransaction(tx, &err)

	isPrimary := true
	findOptions := repoCalendar.CalendarQueryOptions{
		IsPrimary: &isPrimary,
	}
	entCalendar, err := ecm.event.CalendarRepo.FindByFields(ctx, tx, userID, findOptions)
	if err != nil {
		return fmt.Errorf("failed to get primary calendar for account: %s, error: %w", email, err)
	}

	convEvent := ecm.event.CalendarApp.ConvertToCalendarEvent(nil, eventReq.Title, eventReq.Location, eventReq.Description, eventReq.SelectedDates[0].Start, eventReq.SelectedDates[0].End)

	entEvent, err := ecm.event.EventRepo.Create(ctx, tx, convEvent, entCalendar)
	if err != nil {
		return fmt.Errorf("failed to get event for account: %s, error: %w", email, err)
	}

	_, err = ecm.event.DateRepo.CreateBulk(ctx, tx, eventReq.SelectedDates, nil, entEvent)
	if err != nil {
		return fmt.Errorf("failed to create proposed dates for account: %s, error: %w", email, err)
	}

	// トランザクションをコミット
	return nil
}
