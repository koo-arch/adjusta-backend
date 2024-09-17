package account

import (
	"fmt"
)

type Account struct {
	ID           int    `json:"id"`
	Email        string `json:"google_account_email"`
	GoogleID     string `json:"google_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewAccount(id int, Email, googleID, accessToken, refreshToken string) *Account {
	return &Account{
		ID:           id,
		GoogleID:     googleID,
		Email:        Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

func (a *Account) Validate() error {
	if a.Email == "" {
		return fmt.Errorf("google account email is required")
	}
	return nil
}
