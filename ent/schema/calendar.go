package schema

import (
	"entgo.io/ent"
	"github.com/google/uuid"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/index"
)

// Calendar holds the schema definition for the Calendar entity.
type Calendar struct {
	ent.Schema
}

// Fields of the Calendar.
func (Calendar) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.String("calendar_id"),
		field.String("summary"),
	}
}

// Edges of the Calendar.
func (Calendar) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).Ref("calendars").Unique(),
		edge.To("events", Event.Type),
	}
}

// Indexes of the Calendar.
func (Calendar) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("calendar_id").Edges("account").Unique(),
	}
}