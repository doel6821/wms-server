package helpers

import (
	"time"
	cModels "wms-server/controllers/v1/models"
	pModels "wms-server/databases/postgre/models"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	authorizationSecretKey = GetEnv("SECRET_KEY")
	secretKey              = []byte(authorizationSecretKey)
)


// GenerateToken ..
func GenerateToken(user pModels.User) (string, int64, error) {
	PrintHeader()
	var claims cModels.ClaimsUserDto

	claims.Scope.Role = user.Role
	claims.Scope.Email = user.Email
	claims.Scope.Tenant = user.Tenant
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(8 * time.Hour).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", 0, err
	}

	return tokenString, claims.ExpiresAt, nil
}


