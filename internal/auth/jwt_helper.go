package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtProvider struct {
	secretKey []byte
}

func NewJwtProvider(secretKey string) (*JwtProvider, error) {
	if secretKey == "" {
		return nil, fmt.Errorf("JWT secret key is not set")
	}
	return &JwtProvider{secretKey: []byte(secretKey)}, nil
}

type VerificationTokenClaims struct {
	ID               string `json:"id"`
	VerificationCode string `json:"verification_code"`
	jwt.RegisteredClaims
}

func (tp *JwtProvider) GenerateVerificationToken(id string, verificationCode string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &VerificationTokenClaims{
		ID:               id,
		VerificationCode: verificationCode,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(tp.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}
