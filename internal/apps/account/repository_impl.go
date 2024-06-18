package account

import (
	"context"

	"golang.org/x/oauth2"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/ent/user"
	"github.com/koo-arch/adjusta-backend/ent/account"
)

type AccountRepositoryImpl struct {
	client *ent.Client
}

func NewAccountRepository(client *ent.Client) *AccountRepositoryImpl {
	return &AccountRepositoryImpl{
		client: client,
	}
}

func (r *AccountRepositoryImpl) Read(ctx context.Context, tx *ent.Tx, id int) (*ent.Account, error) {
	if tx != nil {
		return tx.Account.Get(ctx, id)
	}
	return r.client.Account.Get(ctx, id)
}

func (r *AccountRepositoryImpl) FindByEmail(ctx context.Context, tx *ent.Tx, email string) (*ent.Account, error) {
	if tx != nil {
		return tx.Account.Query().
			Where(account.EmailEQ(email)).
			Only(ctx)
	}
	return r.client.Account.Query().
		Where(account.EmailEQ(email)).
		Only(ctx)
}

func (r *AccountRepositoryImpl) FilterByUserID(ctx context.Context, tx *ent.Tx, userID int) ([]*ent.Account, error) {
	if tx != nil {
		return tx.Account.Query().
			Where(account.HasUserWith(user.ID(userID))).
			All(ctx)
	}
	return r.client.Account.Query().
		Where(account.HasUserWith(user.ID(userID))).
		All(ctx)
}

func (r *AccountRepositoryImpl) Create(ctx context.Context, tx *ent.Tx, email, googleID string, oauthToken *oauth2.Token, user *ent.User) (*ent.Account, error) {
	if tx != nil {
		return tx.Account.Create().
			SetEmail(email).
			SetGoogleID(googleID).
			SetAccessToken(oauthToken.AccessToken).
			SetRefreshToken(oauthToken.RefreshToken).
			SetAccessTokenExpiry(oauthToken.Expiry).
			SetUser(user).
			Save(ctx)
	}
	return r.client.Account.Create().
		SetEmail(email).
		SetGoogleID(googleID).
		SetAccessToken(oauthToken.AccessToken).
		SetRefreshToken(oauthToken.RefreshToken).
		SetAccessTokenExpiry(oauthToken.Expiry).
		SetUser(user).
		Save(ctx)
}

func (r *AccountRepositoryImpl) Update(ctx context.Context, tx *ent.Tx, id int, oauthToken *oauth2.Token) (*ent.Account, error) {
	if tx != nil {
		return tx.Account.UpdateOneID(id).
			SetNillableAccessToken(&oauthToken.AccessToken).
			SetNillableRefreshToken(&oauthToken.RefreshToken).
			SetNillableAccessTokenExpiry(&oauthToken.Expiry).
			Save(ctx)
	}
	return r.client.Account.UpdateOneID(id).
		SetNillableAccessToken(&oauthToken.AccessToken).
		SetNillableRefreshToken(&oauthToken.RefreshToken).
		SetNillableAccessTokenExpiry(&oauthToken.Expiry).
		Save(ctx)
}

func (r *AccountRepositoryImpl) Delete(ctx context.Context, tx *ent.Tx, id int) error {
	if tx != nil {
		return tx.Account.DeleteOneID(id).Exec(ctx)
	}
	return r.client.Account.DeleteOneID(id).Exec(ctx)
}
