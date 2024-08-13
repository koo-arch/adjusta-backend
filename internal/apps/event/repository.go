package event

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
)

type EventRepository interface {
	Read(ctx context.Context, tx *ent.Tx, id uuid.UUID) (*ent.Event, error)
	FilterByCalendarID(ctx context.Context, tx *ent.Tx, calendarID string) ([]*ent.Event, error)
	FindByCalendarIDAndEventID(ctx context.Context, tx *ent.Tx, calendarID, eventID string) (*ent.Event, error)
	Create(ctx context.Context, tx *ent.Tx, eventID string, summary, description, location *string, calendar *ent.Calendar, startTime, endTime time.Time) (*ent.Event, error)
	Update(ctx context.Context, tx *ent.Tx, id uuid.UUID, summary, description, location *string, startTime, endTime time.Time) (*ent.Event, error)
	Delete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error
}