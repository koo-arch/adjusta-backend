package user

import (
	"context"

	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/models"
	"github.com/google/uuid"
)

type UserQueryOptions struct {
	WithOAuthToken bool
}

type UserRepository interface {
	Read(ctx context.Context, tx *ent.Tx, id uuid.UUID, opt UserQueryOptions) (*ent.User, error)
	FindByEmail(ctx context.Context, tx *ent.Tx, email string) (*ent.User, error)
	Create(ctx context.Context, tx *ent.Tx, email string, jwtToken *models.JWTToken) (*ent.User, error)
	Update(ctx context.Context, tx *ent.Tx, id uuid.UUID, jwtToken *models.JWTToken) (*ent.User, error)
	Delete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error
}