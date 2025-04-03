package utils

import (
	"crypto/rand"
	"math/big"
)

// GenerateRandomID generates a URL-safe random ID of the specified length
func GenerateRandomID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
	charsetLength := big.NewInt(int64(len(charset)))
	
	b := make([]byte, length)
	for i := range b {
		// Use crypto/rand for better randomness
		randomIndex, _ := rand.Int(rand.Reader, charsetLength)
		b[i] = charset[randomIndex.Int64()]
	}
	return string(b)
}
