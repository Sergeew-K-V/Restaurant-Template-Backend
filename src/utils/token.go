package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func GenerateRandomToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)

	return hex.EncodeToString(bytes)
}

func HashString(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}
