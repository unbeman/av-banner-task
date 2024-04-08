package utils

import (
	"errors"
	"fmt"
)
import "github.com/dgrijalva/jwt-go"

var (
	ErrInvalidToken = errors.New("invalid token")
)

type JWTManager struct {
	privateKey string
}

func NewJWTManager(privateKey string) (*JWTManager, error) {
	return &JWTManager{privateKey: privateKey}, nil
}

type UserClaims struct {
	jwt.StandardClaims
	UserRole int
}

func (m *JWTManager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}
			return m.privateKey, nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", err, ErrInvalidToken)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("couldn't get user claims: %w", ErrInvalidToken)
	}

	return claims, nil
}
