package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
