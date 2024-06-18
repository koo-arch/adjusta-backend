package auth

import (
	"context"
	"encoding/json"
	"time"
	"fmt"
	"log"
	"net/http"

	"github.com/koo-arch/adjusta-backend/configs"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/apps/user"
	"github.com/koo-arch/adjusta-backend/internal/apps/account"
	"github.com/koo-arch/adjusta-backend/internal/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleOAuthConfig oauth2.Config

type UserInfo struct {
	GoogleID string `json:"sub"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
}

func init() {
	configs.LoadEnv()

	GoogleOAuthConfig = oauth2.Config{
		ClientID:     configs.GetEnv("GOOGLE_CLIENT_ID"),
		ClientSecret: configs.GetEnv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  configs.GetEnv("GOOGLE_REDIRECT_URI"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/calendar", "openid"},
		Endpoint:     google.Endpoint,
	}
}

func GetGoogleAuthConfig() *oauth2.Config {
	return &GoogleOAuthConfig
}

func FetchGoogleUserInfo(ctx context.Context, token *oauth2.Token) (*UserInfo, error) {
	client := GoogleOAuthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, fmt.Errorf("request to Google userinfo API failed: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed closing response body: %v", err)
		}
	}()

	// レスポンスのステータスコードが200以外の場合はエラー
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request to Google userinfo API failed: %s", resp.Status)
	}

	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %w", err)
	}

	return &userInfo, nil
}

func CreateUserOrUpdateRefreshToken(ctx context.Context, client *ent.Client, userRepo user.UserRepository, accountRepo account.AccountRepository, userInfo *UserInfo, jwtToken *models.JWTToken, oauthToken *oauth2.Token) (*ent.User, error) {
	u, err := userRepo.FindByEmail(ctx, nil, userInfo.Email)
	if err != nil {
		if ent.IsNotFound(err) {
			// ユーザーが存在しない場合は新規作成
			return CreateUserandAccount(ctx, client, userRepo, accountRepo, userInfo, jwtToken, oauthToken)
		}
		// エラーが発生した場合
		return nil, fmt.Errorf("error querying user: %w", err)
	}
	// ユーザーが存在する場合はリフレッシュトークンを更新
	return UpdateTokens(ctx, userRepo, accountRepo, u.ID, jwtToken, oauthToken)
}

func CreateUserandAccount(ctx context.Context, client *ent.Client, userRepo user.UserRepository, accountRepo account.AccountRepository, userInfo *UserInfo, jwtToken *models.JWTToken, oauthToken *oauth2.Token) (*ent.User, error) {
	tx, err := client.Tx(ctx)
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

	u, err := userRepo.Create(ctx, tx, userInfo.Email, jwtToken)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	_, err = accountRepo.Create(ctx, tx, userInfo.Email, userInfo.GoogleID, oauthToken, u)
	if err != nil {
		return nil, fmt.Errorf("error creating account: %w", err)
	}

	return u, nil
}


func CreateAccount(ctx context.Context, client *ent.Client, userRepo user.UserRepository, accountRepo account.AccountRepository, userInfo *UserInfo, oauthToken *oauth2.Token) (*ent.Account, error) {
	// トランザクションを開始
	tx, err := client.Tx(ctx)
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
	u, err := userRepo.FindByEmail(ctx, tx, userInfo.Email)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, fmt.Errorf("error querying user: %w", err)
		}
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// アカウントを作成
	a, err := accountRepo.Create(ctx, tx, userInfo.Email, userInfo.GoogleID, oauthToken, u)
	if err != nil {
		return nil, fmt.Errorf("error creating account: %w", err)
	}
	return a, nil
}

func UpdateTokens(ctx context.Context, userRepo user.UserRepository, accountRepo account.AccountRepository, id int, jwtToken *models.JWTToken, oauthToken *oauth2.Token) (*ent.User, error) {
	u, err := userRepo.Update(ctx, nil, id, jwtToken)
	if err != nil {
		return nil, fmt.Errorf("error updating refresh token: %w", err)
	}

	a, err := accountRepo.FilterByUserID(ctx, nil, u.ID)
	if err != nil {
		return nil, fmt.Errorf("error querying account: %w", err)
	}

	for _, account := range a {
		if account.Email == u.Email {
			_, err = accountRepo.Update(ctx, nil, account.ID, oauthToken)
			if err != nil {
				return nil, fmt.Errorf("error updating account: %w", err)
			}
			return u, nil
		}
	}

	return nil, fmt.Errorf("account not found")
}

func RefreshOAuthToken(ctx context.Context, token *oauth2.Token) (*oauth2.Token, error) {
	newToken, err := GoogleOAuthConfig.TokenSource(ctx, token).Token()
	if err != nil {
		return nil, fmt.Errorf("error refreshing token: %w", err)
	}

	return newToken, nil
}

func VerifyOAuthToken(ctx context.Context, accountRepo account.AccountRepository, email string) (*oauth2.Token, error) {
	account, err := accountRepo.FindByEmail(ctx, nil, email)
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
		token, err = RefreshOAuthToken(ctx, token)
		if err != nil {
			return nil, fmt.Errorf("error refreshing token: %w", err)
		}

		// トークンを保存
		_, err = accountRepo.Update(ctx, nil, account.ID, token)
		if err != nil {
			return nil, fmt.Errorf("error updating account: %w", err)
		}
	}

	return token, nil
}