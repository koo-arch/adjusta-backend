package user

import (
	"context"

	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/ent/user"
	"github.com/koo-arch/adjusta-backend/internal/models"
	"github.com/google/uuid"
)

type UserRepositoryImpl struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		client: client,
	}
}

func (r *UserRepositoryImpl) Read(ctx context.Context, tx *ent.Tx, id uuid.UUID) (*ent.User, error) {
	if tx != nil {
		return tx.User.Get(ctx, id)
	}
	return r.client.User.Get(ctx, id)
}

func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, tx *ent.Tx, email string) (*ent.User, error) {
	findUser := r.client.User.Query()
	if tx != nil {
		findUser = tx.User.Query()
	}
	return findUser.
		Where(user.EmailEQ(email)).
		Only(ctx)
}

func (r *UserRepositoryImpl) Create(ctx context.Context, tx *ent.Tx, email string, jwtToken *models.JWTToken) (*ent.User, error) {
	userCreate := r.client.User.Create()
	if tx != nil {
		userCreate = tx.User.Create()
	}
	return userCreate.
		SetEmail(email).
		SetNillableRefreshToken(&jwtToken.RefreshToken).
		SetNillableRefreshTokenExpiry(&jwtToken.RefreshExpiration).
		Save(ctx)
}

func (r *UserRepositoryImpl) Update(ctx context.Context,tx *ent.Tx, id uuid.UUID, jwtToken *models.JWTToken) (*ent.User, error) {
	userUpdate := r.client.User.UpdateOneID(id)
	if tx != nil {
		userUpdate = tx.User.UpdateOneID(id)
	}
	return userUpdate.
		SetNillableRefreshToken(&jwtToken.RefreshToken).
		SetNillableRefreshTokenExpiry(&jwtToken.RefreshExpiration).
		Save(ctx)
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error {
	if tx != nil {
		return tx.User.DeleteOneID(id).Exec(ctx)
	}
	return r.client.User.DeleteOneID(id).Exec(ctx)
}
	