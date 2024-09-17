package event

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/api/calendar/v3"
	"github.com/koo-arch/adjusta-backend/ent"
)

type EventQueryOptions struct {
	WithProposedDates bool
	EventOffset       int
	EventLimit        int
	ProposedDateOffset int
	ProposedDateLimit  int
}

type EventRepository interface {
	Read(ctx context.Context, tx *ent.Tx, id uuid.UUID, opt EventQueryOptions) (*ent.Event, error)
	FilterByCalendarID(ctx context.Context, tx *ent.Tx, calendarID string, opt EventQueryOptions) ([]*ent.Event, error)
	Create(ctx context.Context, tx *ent.Tx, event *calendar.Event, entCalendar *ent.Calendar) (*ent.Event, error)
	Update(ctx context.Context, tx *ent.Tx, id uuid.UUID, event *calendar.Event) (*ent.Event, error)
	Delete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error
}