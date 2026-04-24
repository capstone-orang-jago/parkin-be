package helpers

import (
	"errors"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(userID, role, nama, secret string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"nama": nama,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(duration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 64)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b), nil
}

func ValidateToken(token, secret string) (jwt.MapClaims, error) {
	tkn, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !tkn.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := tkn.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to parse claims")
	}
	return claims, nil
}
