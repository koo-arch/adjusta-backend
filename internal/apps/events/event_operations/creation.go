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

func (ecm *EventCreationManager) CreateDraftedEvents(ctx context.Context, userID uuid.UUID, email string, eventReq *models.EventDraftCreation) (*models.EventDraftDetail, error) {

	tx, err := ecm.event.Client.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed starting transaction: %w", err)
	}

	defer transaction.HandleTransaction(tx, &err)

	isPrimary := true
	findOptions := repoCalendar.CalendarQueryOptions{
		IsPrimary: &isPrimary,
	}
	entCalendar, err := ecm.event.CalendarRepo.FindByFields(ctx, tx, userID, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to get primary calendar for account: %s, error: %w", email, err)
	}

	convEvent := ecm.event.CalendarApp.ConvertToCalendarEvent(nil, eventReq.Title, eventReq.Location, eventReq.Description, eventReq.SelectedDates[0].Start, eventReq.SelectedDates[0].End)

	entEvent, err := ecm.event.EventRepo.Create(ctx, tx, convEvent, entCalendar)
	if err != nil {
		return nil, fmt.Errorf("failed to get event for account: %s, error: %w", email, err)
	}

	entDate, err := ecm.event.DateRepo.CreateBulk(ctx, tx, eventReq.SelectedDates, entEvent)
	if err != nil {
		return nil, fmt.Errorf("failed to create proposed dates for account: %s, error: %w", email, err)
	}

	eventDates := make([]models.ProposedDate, 0)
	for _, date := range entDate {
		eventDates = append(eventDates, models.ProposedDate{
			ID:        &date.ID,
			Start:     &date.StartTime,
			End:       &date.EndTime,
			Priority: date.Priority,
		})
	}
	response := &models.EventDraftDetail{
		ID:            entEvent.ID,
		Title:         entEvent.Summary,
		Location:      entEvent.Location,
		Description:   entEvent.Description,
		Status:	  	   models.EventStatus(entEvent.Status),
		ProposedDates: eventDates,
	}

	// トランザクションをコミット
	return response, nil
}
