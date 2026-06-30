package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"

	"golang.org/x/crypto/blake2b"
)

const tokenBytes = 32

func GenerateToken() (string, error) {
	bytes := make([]byte, tokenBytes)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func HashToken(token string) string {
	hash := blake2b.Sum256([]byte(token))

	return hex.EncodeToString(hash[:])
}
