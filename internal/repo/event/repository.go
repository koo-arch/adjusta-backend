package event

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/api/calendar/v3"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/models"
)

type EventQueryOptions struct {
	Summary		   *string
	Location	   *string
	Description	   *string
	Status 		     *models.EventStatus
	ConfirmedDateID *uuid.UUID
	WithProposedDates bool
	EventOffset       int
	EventLimit        int
	ProposedDateOffset int
	ProposedDateLimit  int
}

type EventRepository interface {
	Read(ctx context.Context, tx *ent.Tx, id uuid.UUID, opt EventQueryOptions) (*ent.Event, error)
	FilterByCalendarID(ctx context.Context, tx *ent.Tx, calendarID string, opt EventQueryOptions) ([]*ent.Event, error)
	Create(ctx context.Context, tx *ent.Tx, googleEvent *calendar.Event, entCalendar *ent.Calendar) (*ent.Event, error)
	Update(ctx context.Context, tx *ent.Tx, id uuid.UUID, opt EventQueryOptions) (*ent.Event, error)
	Delete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error
}