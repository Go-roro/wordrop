package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJwtProvider_GenerateVerificationToken(t *testing.T) {
	t.Run("GenerateVerificationToken", func(t *testing.T) {
		provider, err := NewJwtProvider("a-string-secret-at-least-256-bits-long")
		if err != nil {
			t.Fatalf("Failed to create JWT provider: %v", err)
		}

		token, err := provider.GenerateVerificationToken("test-id", "test-code")
		assert.NoError(t, err, "Expected no error when creating JWT provider")
		assert.NotEmpty(t, token, "Expected token to be generated")
	})
}

func TestJwtProvider_ParseVerificationToken(t *testing.T) {
	secret := "a-string-secret-at-least-256-bits-long"
	tokenGenerator, err := NewJwtProvider(secret)
	require.NoError(t, err)

	t.Run("Failure-Expired Token", func(t *testing.T) {
		expiredToken := generateExpiredToken(secret)
		_, err := tokenGenerator.ParseVerificationToken(expiredToken)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, jwt.ErrTokenExpired))
	})

	t.Run("Failure-Invalid Signature", func(t *testing.T) {
		wrongSecret := "different-secret"
		otherTokenGenerator, _ := NewJwtProvider(wrongSecret)

		tokenString, err := tokenGenerator.GenerateVerificationToken("user-123", "code")
		require.NoError(t, err)

		_, err = otherTokenGenerator.ParseVerificationToken(tokenString)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, jwt.ErrTokenSignatureInvalid))
	})

	t.Run("Failure-Malformed Token", func(t *testing.T) {
		malformedToken := "not.a.real.jwt"

		_, err := tokenGenerator.ParseVerificationToken(malformedToken)

		assert.Error(t, err)
	})
}

func generateExpiredToken(secret string) string {
	expirationTime := time.Now()
	claims := &VerificationTokenClaims{
		ID:               "user12",
		VerificationCode: "code123",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}
