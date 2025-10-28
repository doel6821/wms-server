package models

// AuthTokenData represents the authentication token data structure
type AuthTokenData struct {
	Token     string `json:"token"`
	TokenType string `json:"tokenType"`
	ExpiresIn int64  `json:"expiresIn"`
}

// AuthResponse wraps the token data in the standard response format
// This follows the same pattern as other API endpoints that use hModels.Response
type AuthResponse struct {
	Token     string `json:"token"`
	Role      string `json:"role"`
	TokenType string `json:"tokenType"`
	ExpiresIn int64  `json:"expiresIn"`
}
