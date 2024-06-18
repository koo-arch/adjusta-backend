package schema

import (
	"context"
	"errors"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	gen "github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/ent/account"
	"github.com/koo-arch/adjusta-backend/ent/hook"
	"github.com/koo-arch/adjusta-backend/ent/user"
)

// Account holds the schema definition for the Account entity.
type Account struct {
	ent.Schema
}

// Fields of the Account.
func (Account) Fields() []ent.Field {
	return []ent.Field{
		field.String("email").NotEmpty(),
		field.String("google_id").NotEmpty(),
		field.String("access_token").Sensitive().Optional(),
		field.String("refresh_token").Sensitive().Optional(),
		field.Time("access_token_expiry").Optional(),
	}
}

// Edges of the Account.
func (Account) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("accounts").Unique(),
	}
}

func (Account) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(accounthook, ent.OpCreate|ent.OpUpdate),
		hook.On(checkEmailUniquePerUser, ent.OpCreate|ent.OpUpdate),
	}
}

func accounthook(next ent.Mutator) ent.Mutator {
	return hook.AccountFunc(func(ctx context.Context, m *gen.AccountMutation) (ent.Value, error) {
		if Email, ok := m.Email(); ok {
			if !emailRegex.MatchString(Email) {
				return nil, errors.New("invalid email address")
			}
		}

		return next.Mutate(ctx, m)
	})
}

func checkEmailUniquePerUser(next ent.Mutator) ent.Mutator {
	return hook.AccountFunc(func(ctx context.Context, m *gen.AccountMutation) (ent.Value, error) {
		if email, exists := m.Email(); exists {
			uid, exists := m.UserID()
			if !exists {
				return nil, errors.New("user id is required")
			}
			// Check if the email address is already in use by another account
			a, err := m.Client().Account.
				Query().
				Where(
					account.HasUserWith(user.ID(uid)),
					account.EmailEQ(email),
				).
				Only(ctx)
			if err != nil && !gen.IsNotFound(err) {
				return nil, err
			}
			if a != nil {
				return nil, errors.New("email address is already in use")
			}
		}
		return next.Mutate(ctx, m)
	})
}
