package oauth

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"github.com/koo-arch/adjusta-backend/configs"
)

var (
	GoogleOAuthConfig   oauth2.Config
	AddAccountOAuthConfig oauth2.Config
)

func init() {
	configs.LoadEnv()

	GoogleOAuthConfig = oauth2.Config{
		ClientID:     configs.GetEnv("GOOGLE_CLIENT_ID"),
		ClientSecret: configs.GetEnv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  configs.GetEnv("GOOGLE_REDIRECT_URI"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/calendar", "openid"},
		Endpoint:     google.Endpoint,
	}

	// 別のRedirectURLを持つOAuthConfig
	AddAccountOAuthConfig = GoogleOAuthConfig
	AddAccountOAuthConfig.RedirectURL = configs.GetEnv("GOOGLE_ADD_ACCOUNT_REDIRECT_URI")
}

func GetGoogleAuthConfig() *oauth2.Config {
	return &GoogleOAuthConfig
}

func GetAddAccountAuthConfig() *oauth2.Config {
	return &AddAccountOAuthConfig
}

func RefreshOAuthToken(ctx context.Context, token *oauth2.Token) (*oauth2.Token, error) {
	newToken, err := GoogleOAuthConfig.TokenSource(ctx, token).Token()
	if err != nil {
		return nil, fmt.Errorf("error refreshing token: %w", err)
	}

	return newToken, nil
}