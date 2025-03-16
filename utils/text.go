package utils

import (
	"crypto/rand"
	"time"
)

// Helper function to truncate long strings for logging
func Truncate(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + "..."
}

func GenerateShortCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 6

	// Create a byte array with the code length
	shortCode := make([]byte, codeLength)

	// Use crypto/rand for secure random generation
	randomBytes := make([]byte, codeLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		// Fallback to less secure but working solution if crypto/rand fails
		for i := range shortCode {
			shortCode[i] = charset[time.Now().UnixNano()%int64(len(charset))]
			time.Sleep(1 * time.Nanosecond) // Add a tiny delay to change the seed
		}
		return string(shortCode)
	}

	// Map random bytes to characters in the charset
	for i, b := range randomBytes {
		shortCode[i] = charset[int(b)%len(charset)]
	}

	return string(shortCode)
}
