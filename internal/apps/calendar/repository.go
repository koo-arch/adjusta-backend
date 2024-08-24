package calendar

import (
	"context"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
)

type CalendarRepository interface {
	Read(ctx context.Context, tx *ent.Tx, id uuid.UUID) (*ent.Calendar, error)
	FilterByAccountID(ctx context.Context, tx *ent.Tx, accountID uuid.UUID) ([]*ent.Calendar, error)
	FindByFields(ctx context.Context, tx *ent.Tx, accountID uuid.UUID, calendarID, summary *string, isPrimary *bool) (*ent.Calendar, error)
	FilterByFields(ctx context.Context, tx *ent.Tx, accountID uuid.UUID, calendarID, summary *string, isPrimary *bool) ([]*ent.Calendar, error)
	Create(ctx context.Context, tx *ent.Tx, calendarID string, summary string, is_primary bool, account *ent.Account) (*ent.Calendar, error)
	Update(ctx context.Context, tx *ent.Tx, id uuid.UUID, summary string) (*ent.Calendar, error)
	Delete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error
}