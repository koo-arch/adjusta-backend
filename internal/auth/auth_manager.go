package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/google/oauth"
	"github.com/koo-arch/adjusta-backend/internal/google/userinfo"
	"github.com/koo-arch/adjusta-backend/internal/models"
	"github.com/koo-arch/adjusta-backend/internal/repo/oauthtoken"
	"github.com/koo-arch/adjusta-backend/internal/repo/user"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
	"golang.org/x/oauth2"
)

type AuthManager struct {
	client      *ent.Client
	userRepo    user.UserRepository
	oauthRepo	oauthtoken.OAuthTokenRepository
}

func NewAuthManager(client *ent.Client, userRepo user.UserRepository, oauthRepo oauthtoken.OAuthTokenRepository) *AuthManager {
	return &AuthManager{
		client:      client,
		userRepo:    userRepo,
		oauthRepo:   oauthRepo,
	}
}

func (am *AuthManager) ProcessUserSignIn(ctx context.Context, userInfo *userinfo.UserInfo, jwtToken *models.JWTToken, oauthToken *oauth2.Token) (*ent.User, error) {
	userOtions := user.UserQueryOptions{
		WithOAuthToken: true,
	}
	u, err := am.userRepo.FindByEmail(ctx, nil, userInfo.Email, userOtions)
	if err != nil {
		if ent.IsNotFound(err) {
			// ユーザーが存在しない場合は新規作成
			return am.CreateUser(ctx, userInfo, jwtToken, oauthToken)
		}
		// エラーが発生した場合
		return nil, fmt.Errorf("error querying user: %w", err)
	}
	// ユーザーが存在する場合はトークンを更新
	println("updating login token")
	return am.UpdateJWTAndOAuth(ctx, u.ID, u.Edges.OauthToken.ID, jwtToken, oauthToken)
}

func (am *AuthManager) CreateUser(ctx context.Context, userInfo *userinfo.UserInfo, jwtToken *models.JWTToken, oauthToken *oauth2.Token) (*ent.User, error) {
	// トランザクションを開始
	tx, err := am.client.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed starting transaction: %w", err)
	}

	// トランザクションエラー用の変数を別に定義
	defer transaction.HandleTransaction(tx, &err)

	u, err := am.userRepo.Create(ctx, tx, userInfo.Email, jwtToken)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	_, err = am.oauthRepo.Create(ctx, tx, oauthToken, u)
	if err != nil {
		return nil, fmt.Errorf("error creating oauthtoken: %w", err)
	}

	return u, nil
}

func (am *AuthManager) UpdateJWTAndOAuth(ctx context.Context, userID, oauthTokenID uuid.UUID, jwtToken *models.JWTToken, oauthToken *oauth2.Token) (*ent.User, error) {
	tx, err := am.client.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed starting transaction: %w", err)
	}

	defer transaction.HandleTransaction(tx, &err)

	u, err := am.userRepo.Update(ctx, tx, userID, jwtToken)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	oauthOptions := oauthtoken.OAuthTokenQuertOptions{
		AccessToken:  &oauthToken.AccessToken,
		RefreshToken: &oauthToken.RefreshToken,
		Expiry:       &oauthToken.Expiry,
	}

	_, err = am.oauthRepo.Update(ctx, tx, oauthTokenID, oauthOptions)
	if err != nil {
		return nil, fmt.Errorf("error updating oauthtoken: %w", err)
	}

	return u, nil

}


func (am *AuthManager) VerifyOAuthToken(ctx context.Context, userID uuid.UUID) (*oauth2.Token, error) {
	userOptions := user.UserQueryOptions{
		WithOAuthToken: true,
	}
	entUser, err := am.userRepo.Read(ctx, nil, userID, userOptions)
	if err != nil {
		return nil, fmt.Errorf("error reading user: %w", err)
	}

	token := &oauth2.Token{
		AccessToken:  entUser.Edges.OauthToken.AccessToken,
		TokenType:    "Bearer",
		RefreshToken: entUser.Edges.OauthToken.RefreshToken,
		Expiry:       entUser.Edges.OauthToken.Expiry,
	}

	// トークンが期限切れの場合は再取得
	newToken, err := oauth2.ReuseTokenSource(token, oauth.GoogleOAuthConfig.TokenSource(ctx, token)).Token()
	if err != nil {
		fmt.Printf("error getting token: %v\n", err)
		return nil, fmt.Errorf("error getting token: %w", err)
	}

	// トークンが再発行された場合、データベースを更新
	if newToken.AccessToken != token.AccessToken {
		println("updating token")
		oauthtokenOptions := oauthtoken.OAuthTokenQuertOptions{
			AccessToken:  &newToken.AccessToken,
			RefreshToken: &newToken.RefreshToken,
			Expiry:       &newToken.Expiry,
		}
		_, err = am.oauthRepo.Update(ctx, nil, entUser.Edges.OauthToken.ID, oauthtokenOptions)
		if err != nil {
			return nil, fmt.Errorf("error updating token: %w", err)
		}
	}

	return newToken, nil
}
