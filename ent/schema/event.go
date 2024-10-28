package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/dialect/entsql"
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
		field.String("summary").Optional(),
		field.String("description").Optional(),
		field.String("location").Optional(),
	}
}

// Edges of the Event.
func (Event) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("calendar", Calendar.Type).Ref("events").Unique(),
		edge.To("proposed_dates", ProposedDate.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
