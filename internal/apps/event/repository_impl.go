package event

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/ent/calendar"
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
		Where(event.HasCalendarWith(calendar.CalendarIDEQ(calendarID))).
		All(ctx)
}

func (r *EventRepositoryImpl) FindByCalendarIDAndEventID(ctx context.Context, tx *ent.Tx, calendarID, eventID string) (*ent.Event, error) {
	findEvent := r.client.Event.Query()
	if tx != nil {
		findEvent = tx.Event.Query()
	}
	return findEvent.
		Where(
			event.HasCalendarWith(calendar.CalendarIDEQ(calendarID)),
			event.EventIDEQ(eventID),
		).
		Only(ctx)
}

func (r *EventRepositoryImpl) Create(ctx context.Context, tx *ent.Tx, eventID string, summary, description, location *string, calendar *ent.Calendar, startTime, endTime time.Time) (*ent.Event, error) {
	eventCreate := r.client.Event.Create()
	if tx != nil {
		eventCreate = tx.Event.Create()
	}

	if summary != nil {
		eventCreate = eventCreate.SetSummary(*summary)
	}

	if description != nil {
		eventCreate = eventCreate.SetDescription(*description)
	}

	eventCreate = eventCreate.
		SetEventID(eventID).
		SetStartTime(startTime).
		SetEndTime(endTime).
		SetCalendar(calendar)

	return eventCreate.Save(ctx)
}

func (r *EventRepositoryImpl) Update(ctx context.Context, tx *ent.Tx, id uuid.UUID, summary, description *string, startTime, endTime *time.Time) (*ent.Event, error) {
	eventUpdate := r.client.Event.UpdateOneID(id)
	if tx != nil {
		eventUpdate = tx.Event.UpdateOneID(id)
	}

	if summary != nil {
		eventUpdate = eventUpdate.SetSummary(*summary)
	}

	if description != nil {
		eventUpdate = eventUpdate.SetDescription(*description)
	}

	if startTime != nil {
		eventUpdate = eventUpdate.SetStartTime(*startTime)
	}

	if endTime != nil {
		eventUpdate = eventUpdate.SetEndTime(*endTime)
	}

	return eventUpdate.Save(ctx)
}

func (r *EventRepositoryImpl) Delete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error {
	if tx != nil {
		return tx.Event.DeleteOneID(id).Exec(ctx)
	}
	return r.client.Event.DeleteOneID(id).Exec(ctx)
}