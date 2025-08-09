package common

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRandomToken(length int) (string, error) {
	bytes := make([]byte, length/2) // since hex encoding doubles the byte size
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
