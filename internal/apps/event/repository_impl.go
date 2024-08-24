package event

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/api/calendar/v3"
	"github.com/koo-arch/adjusta-backend/ent"
	dbCalendar "github.com/koo-arch/adjusta-backend/ent/calendar"
	"github.com/koo-arch/adjusta-backend/ent/event"
)

type EventRepositoryImpl struct {
	client *ent.Client
}

func NewEventRepository(client *ent.Client) *EventRepositoryImpl {
	return &EventRepositoryImpl{
		client: client,
	}
}

func (r *EventRepositoryImpl) Read(ctx context.Context, tx *ent.Tx, id uuid.UUID) (*ent.Event, error) {
	if tx != nil {
		return tx.Event.Get(ctx, id)
	}
	return r.client.Event.Get(ctx, id)
}

func (r *EventRepositoryImpl) FilterByCalendarID(ctx context.Context, tx *ent.Tx, calendarID string) ([]*ent.Event, error) {
	filterEvent := r.client.Event.Query()
	if tx != nil {
		filterEvent = tx.Event.Query()
	}
	return filterEvent.
		Where(event.HasCalendarWith(dbCalendar.CalendarIDEQ(calendarID))).
		All(ctx)
}

func (r *EventRepositoryImpl) FindByCalendarIDAndEventID(ctx context.Context, tx *ent.Tx, calendarID, eventID string) (*ent.Event, error) {
	findEvent := r.client.Event.Query()
	if tx != nil {
		findEvent = tx.Event.Query()
	}
	return findEvent.
		Where(
			event.HasCalendarWith(dbCalendar.CalendarIDEQ(calendarID)),
			event.EventIDEQ(eventID),
		).
		Only(ctx)
}

func (r *EventRepositoryImpl) Create(ctx context.Context, tx *ent.Tx, event *calendar.Event, entCalendar *ent.Calendar) (*ent.Event, error) {
	eventCreate := r.client.Event.Create()
	if tx != nil {
		eventCreate = tx.Event.Create()
	}

	eventCreate = eventCreate.
		SetEventID(event.Id).
		SetSummary(event.Summary).
		SetDescription(event.Description).
		SetLocation(event.Location).
		SetCalendar(entCalendar)

	return eventCreate.Save(ctx)
}

func (r *EventRepositoryImpl) Update(ctx context.Context, tx *ent.Tx, id uuid.UUID, event *calendar.Event) (*ent.Event, error) {
	eventUpdate := r.client.Event.UpdateOneID(id)
	if tx != nil {
		eventUpdate = tx.Event.UpdateOneID(id)
	}

	eventUpdate = eventUpdate.
		SetSummary(event.Summary).
		SetDescription(event.Description).
		SetLocation(event.Location)

	return eventUpdate.Save(ctx)
}

func (r *EventRepositoryImpl) Delete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error {
	if tx != nil {
		return tx.Event.DeleteOneID(id).Exec(ctx)
	}
	return r.client.Event.DeleteOneID(id).Exec(ctx)
}