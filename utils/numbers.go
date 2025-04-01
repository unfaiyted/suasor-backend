package utils

import (
	"math/rand"
	"time"
)

// GenerateRandomID generates a random ID of the specified length
func GenerateRandomID(length int) string {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = byte(rand.Intn(256))
	}
	return string(bytes)
}
