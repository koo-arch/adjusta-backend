package user

import (
	"context"

	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/ent/user"
	"github.com/koo-arch/adjusta-backend/internal/models"
)

type UserRepositoryImpl struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		client: client,
	}
}

func (r *UserRepositoryImpl) Read(ctx context.Context, tx *ent.Tx, id int) (*ent.User, error) {
	if tx != nil {
		return tx.User.Get(ctx, id)
	}
	return r.client.User.Get(ctx, id)
}

func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, tx *ent.Tx, email string) (*ent.User, error) {
	if tx != nil {
		return tx.User.Query().
			Where(user.EmailEQ(email)).
			Only(ctx)
	}
	return r.client.User.Query().
		Where(user.EmailEQ(email)).
		Only(ctx)
}

func (r *UserRepositoryImpl) Create(ctx context.Context, tx *ent.Tx, email string, jwtToken *models.JWTToken) (*ent.User, error) {
	if tx != nil {
		return tx.User.Create().
			SetEmail(email).
			SetNillableRefreshToken(&jwtToken.RefreshToken).
			SetNillableRefreshTokenExpiry(&jwtToken.RefreshExpiration).
			Save(ctx)
	}
	return r.client.User.Create().
		SetEmail(email).
		SetNillableRefreshToken(&jwtToken.RefreshToken).
		SetNillableRefreshTokenExpiry(&jwtToken.RefreshExpiration).
		Save(ctx)
}

func (r *UserRepositoryImpl) Update(ctx context.Context,tx *ent.Tx, id int, jwtToken *models.JWTToken) (*ent.User, error) {
	if tx != nil {
		return tx.User.UpdateOneID(id).
			SetNillableRefreshToken(&jwtToken.RefreshToken).
			SetNillableRefreshTokenExpiry(&jwtToken.RefreshExpiration).
			Save(ctx)
	}
	return r.client.User.UpdateOneID(id).
		SetNillableRefreshToken(&jwtToken.RefreshToken).
		SetNillableRefreshTokenExpiry(&jwtToken.RefreshExpiration).
		Save(ctx)
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, tx *ent.Tx, id int) error {
	if tx != nil {
		return tx.User.DeleteOneID(id).Exec(ctx)
	}
	return r.client.User.DeleteOneID(id).Exec(ctx)
}
	