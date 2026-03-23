package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("room-booking-secret-key")

const (
	AdminRole = "admin"
	UserRole  = "user"

	AdminID = "11111111-1111-1111-1111-111111111111"
	UserID  = "22222222-2222-2222-2222-222222222222"
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func FixedUserIDByRole(role string) (string, error) {
	if role == AdminRole {
		return AdminID, nil
	}
	if role == UserRole {
		return UserID, nil
	}
	return "", errors.New("invalid role")
}

func GenerateToken(role string) (string, error) {
	userID, err := FixedUserIDByRole(role)
	if err != nil {
		return "", err
	}

	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ParseToken(raw string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(raw, &Claims{}, func(token *jwt.Token) (any, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
