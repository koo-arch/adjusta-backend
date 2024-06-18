package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// JWTkey holds the schema definition for the JWTkey entity.
type JWTKey struct {
	ent.Schema
}

// Fields of the JWTkey.
func (JWTKey) Fields() []ent.Field {
	return []ent.Field{
		field.String("key").NotEmpty().Sensitive(),
		field.String("type").NotEmpty().Default("access"), // access or refresh
		field.Time("created_at").Immutable(),
		field.Time("expires_at").Immutable(),
	}
}

// Edges of the JWTkey.
func (JWTKey) Edges() []ent.Edge {
	return nil
}

func (JWTKey) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("key", "expires_at").StorageKey("idx_type_expires"),
	}
}