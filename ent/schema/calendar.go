package schema

import (
	"context"
	"fmt"

	"entgo.io/ent"
	"github.com/google/uuid"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/dialect/entsql"
	"github.com/koo-arch/adjusta-backend/ent/user"
	"github.com/koo-arch/adjusta-backend/ent/account"
	"github.com/koo-arch/adjusta-backend/ent/calendar"
	"github.com/koo-arch/adjusta-backend/ent/hook"
	gen "github.com/koo-arch/adjusta-backend/ent"
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
		field.Bool("is_primary").Default(false),
	}
}

// Edges of the Calendar.
func (Calendar) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).Ref("calendars").Unique(),
		edge.To("events", Event.Type).Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

// hooks of the Calendar.
func (Calendar) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(checkCalendarIDUniquePerUser, ent.OpCreate|ent.OpUpdate),
	}
}

// Mixin of the Calendar.
func checkCalendarIDUniquePerUser(next ent.Mutator) ent.Mutator {
	return hook.CalendarFunc(func(ctx context.Context, m *gen.CalendarMutation) (ent.Value, error) {
		// Check if the calendar_id is unique per user
		if calendarID, exists := m.CalendarID(); exists {
			accountID, exists := m.AccountID()
			if !exists {
				return nil, fmt.Errorf("account_id is required")
			}
	
			genUser, err := m.Client().Account.Query().Where(account.ID(accountID)).QueryUser().Only(ctx)
			if err != nil {
				return nil, fmt.Errorf("account not found: %w", err)
			}

			// Check if the calendar_id is already in use by another account of the same user
			exists, err = m.Client().Calendar.Query().Where(
				calendar.HasAccountWith(account.HasUserWith(user.ID(genUser.ID))),
				calendar.CalendarID(calendarID),
			).Exist(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to query calendars for user_id %v: %w", genUser.ID, err)
			}
			if exists {
				return nil, fmt.Errorf("calendar_id '%s' is already in use by another account of the same user", calendarID)
			}
		}

		return next.Mutate(ctx, m)
	})
}