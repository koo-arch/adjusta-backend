package proposeddate

import (
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/api/calendar/v3"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/models"
)

type ProposedDateQueryOptions struct {
	GoogleEventID *string
	StartTime     *time.Time
	EndTime       *time.Time
	Priority      *int
}

type ProposedDateRepository interface {
	Read(ctx context.Context, tx *ent.Tx, id uuid.UUID) (*ent.ProposedDate, error)
	FilterByEventID(ctx context.Context, tx *ent.Tx, eventID uuid.UUID) ([]*ent.ProposedDate, error)
	ExclusionEventID(ctx context.Context, tx *ent.Tx, eventID uuid.UUID) ([]*ent.ProposedDate, error)
	Create(ctx context.Context, tx *ent.Tx, googleEventID *string, opt ProposedDateQueryOptions, entEvent *ent.Event) (*ent.ProposedDate, error)
	Update(ctx context.Context, tx *ent.Tx, id uuid.UUID, opt ProposedDateQueryOptions) (*ent.ProposedDate, error)
	Delete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error
	CreateBulk(ctx context.Context, tx *ent.Tx, selectedDates []models.SelectedDate, googleEvents []*calendar.Event, entEvent *ent.Event) ([]*ent.ProposedDate, error)
	UpdateByGoogleEventID(ctx context.Context, tx *ent.Tx, googleEventID *string, opt ProposedDateQueryOptions) error
}
