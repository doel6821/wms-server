package models

import (
	"github.com/golang-jwt/jwt"
)

type LoginRequest struct {
	// Tenant   string `json:"tenant"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ChangePasswordRequest struct {
	// Tenant   string `json:"tenant"`
	Email    string `json:"email"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Tenant   string `json:"tenant"`
	Role     string `json:"role"`
}

type ClaimsUserDto struct {
	jwt.StandardClaims
	Scope ClaimsUserScopeDto `json:"scopes"`
}

type ClaimsUserScopeDto struct {
	Email  string `json:"email"`
	Role   string    `json:"role"`
	Tenant string `json:"tenant"`
}
