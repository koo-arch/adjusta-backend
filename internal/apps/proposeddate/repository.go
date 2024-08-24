package proposeddate

import (
	"context"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/models"
)

type ProposedDateRepository interface {
	Read(ctx context.Context, tx *ent.Tx, id uuid.UUID) (*ent.ProposedDate, error)
	FilterByEventID(ctx context.Context, tx *ent.Tx, eventID uuid.UUID) ([]*ent.ProposedDate, error)
	Create(ctx context.Context, tx *ent.Tx, selectedDates models.SelectedDate, entEvent *ent.Event, priority int) (*ent.ProposedDate, error)
	Update(ctx context.Context, tx *ent.Tx, id uuid.UUID, selectedDates *models.SelectedDate, priority *int, isFinalized *bool) (*ent.ProposedDate, error)
	Delete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error
	CreateBulk(ctx context.Context, tx *ent.Tx, selectedDates []models.SelectedDate, entEvent *ent.Event) ([]*ent.ProposedDate, error)
}
