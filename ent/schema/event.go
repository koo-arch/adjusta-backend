package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/edge"
	"github.com/google/uuid"
)

// Event holds the schema definition for the Event entity.
type Event struct {
	ent.Schema
}

// Fields of the Event.
func (Event) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("event_id").Unique(),
		field.String("summary").Optional(),
		field.String("description").Optional(),
		field.String("location").Optional(),
		field.Time("start_time").Optional(),
		field.Time("end_time").Optional(),
	}
}

// Edges of the Event.
func (Event) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("calendar", Calendar.Type).Ref("events").Unique(),
	}
}
