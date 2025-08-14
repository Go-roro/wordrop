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

func (p *JwtProvider) GenerateVerificationToken(id string, verificationCode string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &VerificationTokenClaims{
		ID:               id,
		VerificationCode: verificationCode,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(p.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (p *JwtProvider) ParseVerificationToken(tokenString string) (*VerificationTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &VerificationTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return p.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(*VerificationTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token: unable to parse claims")
}
