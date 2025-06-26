package models

type AuthResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"tokenType"`
	ExpiresIn int64  `json:"expiresIn"`
}
