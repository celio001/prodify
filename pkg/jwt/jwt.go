package pkg_jwt

import (
	"errors"
	"time"

	"github.com/celio001/prodify/config"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrInvalidClaims = errors.New("invalid token claims")
	ErrUserNotFound = errors.New("user not found in token")
	ErrTokenTypeNotFound = errors.New("token type not found")
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

// CreateRefreshToken - Big token (7 days)
func CreateRefreshToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": userID,
			"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
			"iat":     time.Now().Unix(),
			"type":    "refresh",
		})

	return token.SignedString([]byte(config.GetString("JWT_SECRET")))
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetString("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return token, nil
}

func GetUserIDFromToken(token *jwt.Token) (string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrInvalidClaims
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", ErrUserNotFound
	}

	return userID, nil
}

func IsAccessToken(token *jwt.Token) (string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrInvalidClaims
	}

	tokenType, ok := claims["type"].(string)
	if !ok {
		return "", ErrTokenTypeNotFound
	}
	return tokenType, nil
}