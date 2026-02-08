package pkg_jwt

import (
	"time"

	"github.com/celio001/prodify/config"
	"github.com/golang-jwt/jwt/v5"
)

// token 15min
func CreateAccessToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": userID,
			"exp":     time.Now().Add(15 * time.Minute).Unix(),
			"iat":     time.Now().Unix(),
			"type":    "access",
		})

	return token.SignedString([]byte(config.GetString("JWT_SECRET")))
}