package models

import (
	"time"

	"github.com/google/uuid"
)

// OAuth saugo Google prisijungimo duomenis
type OAuth struct {
	OAuthProviderID uuid.UUID `db:"oauth_provider_id"`
	UserID          uuid.UUID `db:"user_id"`
	Provider        string    `db:"provider"`    // 'google'
	ProviderID      string    `db:"provider_id"` // Google ID
	ProviderEmail   string    `db:"provider_email"`
	AccessToken     string    `db:"access_token"`
	RefreshToken    string    `db:"refresh_token"`
	CreatedAt       time.Time `db:"created_at"`
}

// OAuthResponse - Google API grąžinami duomenys
type OAuthResponse struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

// TokenResponse - OAuth token atsakymas
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	IDToken      string `json:"id_token"`
}
