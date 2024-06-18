package account

import (
	"context"

	"github.com/koo-arch/adjusta-backend/ent"
	"golang.org/x/oauth2"
)

type AccountRepository interface {
	Read(ctx context.Context, tx *ent.Tx, id int) (*ent.Account, error)
	FindByEmail(ctx context.Context, tx *ent.Tx, email string) (*ent.Account, error)
	FilterByUserID(ctx context.Context, tx *ent.Tx, userID int) ([]*ent.Account, error)
	Create(ctx context.Context, tx *ent.Tx, email, googleID string, oauthToken *oauth2.Token, user *ent.User) (*ent.Account, error)
	Update(ctx context.Context, tx *ent.Tx, id int, oauthToken *oauth2.Token) (*ent.Account, error)
	Delete(ctx context.Context, tx *ent.Tx, id int) error
}