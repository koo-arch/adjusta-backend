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

func (r *EventRepositoryImpl) Read(ctx context.Context, tx *ent.Tx, id uuid.UUID, opt EventQueryOptions) (*ent.Event, error) {
	query := r.client.Event.Query()
	if tx != nil {
		query = tx.Event.Query()
	}

	if opt.WithProposedDates {
		query = query.WithProposedDates()
	}

	return query.Where(event.IDEQ(id)).Only(ctx)
}

func (r *EventRepositoryImpl) FilterByCalendarID(ctx context.Context, tx *ent.Tx, calendarID string, opt EventQueryOptions) ([]*ent.Event, error) {
	filterEvent := r.client.Event.Query()
	if tx != nil {
		filterEvent = tx.Event.Query()
	}

	filterEvent = filterEvent.Where(event.HasCalendarWith(dbCalendar.CalendarIDEQ(calendarID)))

	// イベントに対するオフセットとリミットを適用
	if opt.EventOffset > 0 {
		filterEvent = filterEvent.Offset(opt.EventOffset)
	}
	if opt.EventLimit > 0 {
		filterEvent = filterEvent.Limit(opt.EventLimit)
	}

	// イベントの提案日に対するオフセットとリミットを適用
	if opt.WithProposedDates {
		filterEvent = filterEvent.WithProposedDates(func(query *ent.ProposedDateQuery) {
			if opt.ProposedDateOffset > 0 {
				query.Offset(opt.ProposedDateOffset)
			}
			if opt.ProposedDateLimit > 0 {
				query.Limit(opt.ProposedDateLimit)
			}
		})
	}

	return filterEvent.All(ctx)
}


func (r *EventRepositoryImpl) Create(ctx context.Context, tx *ent.Tx, googleEvent *calendar.Event, entCalendar *ent.Calendar) (*ent.Event, error) {
	eventCreate := r.client.Event.Create()
	if tx != nil {
		eventCreate = tx.Event.Create()
	}

	eventCreate = eventCreate.
		SetSummary(googleEvent.Summary).
		SetDescription(googleEvent.Description).
		SetLocation(googleEvent.Location).
		SetCalendar(entCalendar)

	return eventCreate.Save(ctx)
}

func (r *EventRepositoryImpl) Update(ctx context.Context, tx *ent.Tx, id uuid.UUID, opt EventQueryOptions) (*ent.Event, error) {
	eventUpdate := r.client.Event.UpdateOneID(id)
	if tx != nil {
		eventUpdate = tx.Event.UpdateOneID(id)
	}

	if opt.Summary != nil {
		eventUpdate = eventUpdate.SetSummary(*opt.Summary)
	}

	if opt.Location != nil {
		eventUpdate = eventUpdate.SetLocation(*opt.Location)
	}

	if opt.Description != nil {
		eventUpdate = eventUpdate.SetDescription(*opt.Description)
	}

	if opt.Status != nil {
		status := event.Status(*opt.Status)
		eventUpdate = eventUpdate.SetStatus(status)
	}

	if opt.ConfirmedDateID != nil {
		eventUpdate = eventUpdate.SetConfirmedDateID(*opt.ConfirmedDateID)
	}

	return eventUpdate.Save(ctx)
}

func (r *EventRepositoryImpl) Delete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error {
	if tx != nil {
		return tx.Event.DeleteOneID(id).Exec(ctx)
	}
	return r.client.Event.DeleteOneID(id).Exec(ctx)
}