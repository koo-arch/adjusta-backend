package user

import (
	"context"

	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/models"
)

type UserRepository interface {
	Read(ctx context.Context, tx *ent.Tx, id int) (*ent.User, error)
	FindByEmail(ctx context.Context, tx *ent.Tx, email string) (*ent.User, error)
	Create(ctx context.Context, tx *ent.Tx, email string, jwtToken *models.JWTToken) (*ent.User, error)
	Update(ctx context.Context, tx *ent.Tx, id int, jwtToken *models.JWTToken) (*ent.User, error)
	Delete(ctx context.Context, tx *ent.Tx, id int) error
}