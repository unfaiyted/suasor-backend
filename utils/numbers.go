package utils

import (
	"crypto/rand"
	"math/big"
	"strconv"
	"time"
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

// Helper functions

// parseUint64 parses a string to uint64
func ParseUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

func GetInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

func GetUint64(s string) uint64 {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

// GetStringPtr returns a string pointer from a string or nil if empty
func GetStringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// GetTimeFromUnix converts a unix timestamp string to time.Time
func GetTimeFromUnix(s string) time.Time {
	ts := GetInt64(s)
	if ts == 0 {
		return time.Time{}
	}
	return time.Unix(ts, 0)
}
