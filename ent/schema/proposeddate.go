package schema

import (
	"context"
	"errors"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/edge"
	"github.com/koo-arch/adjusta-backend/ent/hook"
	gen "github.com/koo-arch/adjusta-backend/ent"
	"github.com/google/uuid"
)

// ProposedDate holds the schema definition for the ProposedDate entity.
type ProposedDate struct {
	ent.Schema
}

// Fields of the ProposedDate.
func (ProposedDate) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("google_event_id").Optional(),
		field.Time("start_time"),
		field.Time("end_time"),
		field.Bool("is_finalized").Default(false),
		field.Int("priority").Default(0),
	}
}

// Edges of the ProposedDate.
func (ProposedDate) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("event", Event.Type).Ref("proposed_dates").Unique(),
	}
}

func (ProposedDate) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(proposeddateHook, ent.OpCreate|ent.OpUpdate),
	}
}

func proposeddateHook(next ent.Mutator) ent.Mutator {
	return hook.ProposedDateFunc(func(ctx context.Context, m *gen.ProposedDateMutation) (ent.Value, error) {
		if startTime, ok := m.StartTime(); ok {
			if endTime, ok := m.EndTime(); ok {
				if startTime.After(endTime) {
					return nil, errors.New("start_time must be before end_time")
				}
			}
		}
		return next.Mutate(ctx, m)
	})
}