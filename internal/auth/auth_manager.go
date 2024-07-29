package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/apps/account"
	"github.com/koo-arch/adjusta-backend/internal/apps/user"
	"github.com/koo-arch/adjusta-backend/internal/google/oauth"
	"github.com/koo-arch/adjusta-backend/internal/google/userinfo"
	"github.com/koo-arch/adjusta-backend/internal/models"
	"golang.org/x/oauth2"
	"github.com/google/uuid"
)

type AuthManager struct {
	client      *ent.Client
	userRepo    user.UserRepository
	accountRepo account.AccountRepository
}

func NewAuthManager(client *ent.Client, userRepo user.UserRepository, accountRepo account.AccountRepository) *AuthManager {
	return &AuthManager{
		client:      client,
		userRepo:    userRepo,
		accountRepo: accountRepo,
	}
}

func (am *AuthManager) ProcessUserSignIn(ctx context.Context, userInfo *userinfo.UserInfo, jwtToken *models.JWTToken, oauthToken *oauth2.Token) (*ent.User, error) {
	u, err := am.userRepo.FindByEmail(ctx, nil, userInfo.Email)
	if err != nil {
		if ent.IsNotFound(err) {
			// ユーザーが存在しない場合は新規作成
			return am.CreateUserAndAccount(ctx, userInfo, jwtToken, oauthToken)
		}
		// エラーが発生した場合
		return nil, fmt.Errorf("error querying user: %w", err)
	}
	// ユーザーが存在する場合はリフレッシュトークンを更新
	return am.UpdateTokens(ctx, u.ID, jwtToken, oauthToken)
}

func (am *AuthManager) CreateUserAndAccount(ctx context.Context, userInfo *userinfo.UserInfo, jwtToken *models.JWTToken, oauthToken *oauth2.Token) (*ent.User, error) {
	tx, err := am.client.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed starting transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			if err := tx.Rollback(); err != nil {
				log.Printf("failed rolling back transaction: %v", err)
			}
			panic(p)
		} else if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Printf("failed rolling back transaction: %v", err)
			}
		} else {
			if err := tx.Commit(); err != nil {
				log.Printf("failed committing transaction: %v", err)
			}
		}
	}()

	u, err := am.userRepo.Create(ctx, tx, userInfo.Email, jwtToken)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	_, err = am.accountRepo.Create(ctx, tx, userInfo.Email, userInfo.GoogleID, oauthToken, u)
	if err != nil {
		return nil, fmt.Errorf("error creating account: %w", err)
	}

	return u, nil
}

func (am *AuthManager) CreateAccount(ctx context.Context, userInfo *userinfo.UserInfo, oauthToken *oauth2.Token) (*ent.Account, error) {
	// トランザクションを開始
	tx, err := am.client.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed starting transaction: %w", err)
	}

	// トランザクションの終了時にコミットまたはロールバックを確実に実行
	defer func() {
		if p := recover(); p != nil {
			if err := tx.Rollback(); err != nil {
				log.Printf("failed rolling back transaction: %v", err)
			}
			panic(p) // パニック発生時はロールバック後に再スロー
		} else if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Printf("failed rolling back transaction: %v", err)
			}
		} else {
			if err := tx.Commit(); err != nil {
				log.Printf("failed committing transaction: %v", err)
			}
		}
	}()

	// ユーザーの検索
	u, err := am.userRepo.FindByEmail(ctx, tx, userInfo.Email)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, fmt.Errorf("error querying user: %w", err)
		}
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// アカウントを作成
	a, err := am.accountRepo.Create(ctx, tx, userInfo.Email, userInfo.GoogleID, oauthToken, u)
	if err != nil {
		return nil, fmt.Errorf("error creating account: %w", err)
	}
	return a, nil
}

func (am *AuthManager) UpdateTokens(ctx context.Context, id uuid.UUID, jwtToken *models.JWTToken, oauthToken *oauth2.Token) (*ent.User, error) {
	u, err := am.userRepo.Update(ctx, nil, id, jwtToken)
	if err != nil {
		return nil, fmt.Errorf("error updating refresh token: %w", err)
	}

	a, err := am.accountRepo.FilterByUserID(ctx, nil, u.ID)
	if err != nil {
		return nil, fmt.Errorf("error querying account: %w", err)
	}

	for _, account := range a {
		if account.Email == u.Email {
			_, err = am.accountRepo.Update(ctx, nil, account.ID, oauthToken)
			if err != nil {
				return nil, fmt.Errorf("error updating account: %w", err)
			}
			return u, nil
		}
	}

	return nil, fmt.Errorf("account not found")
}

func (am *AuthManager) VerifyOAuthToken(ctx context.Context, email string) (*oauth2.Token, error) {
	account, err := am.accountRepo.FindByEmail(ctx, nil, email)
	if err != nil {
		return nil, fmt.Errorf("error querying account: %w", err)
	}

	token := &oauth2.Token{
		AccessToken:  account.AccessToken,
		TokenType:    "Bearer",
		RefreshToken: account.RefreshToken,
		Expiry:       account.AccessTokenExpiry,
	}

	IsAccessTokenExpired := account.AccessTokenExpiry.Before(time.Now())

	// アクセストークンが期限切れの場合はリフレッシュ
	if IsAccessTokenExpired {
		// トークンのリフレッシュ
		token, err = oauth.RefreshOAuthToken(ctx, token)
		if err != nil {
			return nil, fmt.Errorf("error refreshing token: %w", err)
		}

		// トークンを保存
		_, err = am.accountRepo.Update(ctx, nil, account.ID, token)
		if err != nil {
			return nil, fmt.Errorf("error updating account: %w", err)
		}
	}

	return token, nil
}
