package utils

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type JWTManager struct {
	privateKey []byte
}

func NewJWTManager(privateKey string) (*JWTManager, error) {
	return &JWTManager{privateKey: []byte(privateKey)}, nil
}

type UserClaims struct {
	jwt.StandardClaims
	UserRole int `json:"user_role"`
}

func (u UserClaims) Valid() error {
	return nil
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
