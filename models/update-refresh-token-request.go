package models

type UpdateRefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
	Email        string `json:"email"`
}
