package userinfo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/koo-arch/adjusta-backend/internal/google/oauth"
	"golang.org/x/oauth2"
)

type UserInfo struct {
	GoogleID string `json:"sub"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
}

func FetchGoogleUserInfo(ctx context.Context, token *oauth2.Token) (*UserInfo, error) {
	client := oauth.GoogleOAuthConfig.Client(ctx, token)
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