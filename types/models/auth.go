package models

import (
	"time"

	"gorm.io/gorm"
)

// Session represents an active user session
// @Description Active user login session
type Session struct {
	gorm.Model

	// UserID is the foreign key to the user
	// @Description ID of the user this session belongs to
	UserID uint64 `json:"userId" gorm:"not null"`

	// RefreshToken is used to generate new access tokens
	// @Description Refresh token value (hashed in database)
	RefreshToken string `json:"-" gorm:"not null"`

	// UserAgent records the client that created this session
	// @Description Browser/client user agent string
	// @Example "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
	UserAgent string `json:"userAgent" gorm:"not null" example:"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"`

	// IP records the IP address of the client
	// @Description IP address of the client
	// @Example "192.168.1.1"
	IP string `json:"ip" gorm:"not null" example:"192.168.1.1"`

	// ExpiresAt indicates when this session should be invalidated
	// @Description When this session expires
	ExpiresAt time.Time `json:"expiresAt" gorm:"not null"`

	// LastUsed tracks the most recent usage of this session
	// @Description Time the session was last used for authentication
	LastUsed time.Time `json:"lastUsed" gorm:"not null"`
}
