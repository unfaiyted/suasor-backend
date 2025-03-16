package models

import "time"

type Shorten struct {
	ID          uint64    `json:"id" example:"1"`
	OriginalURL string    `json:"originalUrl" binding:"required" example:"https://example.com/some/long/path"`
	ShortCode   string    `json:"shortCode" example:"abc123"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	ClickCount  uint64    `json:"clickCount" example:"0"`
	ExpiresAt   time.Time `json:"expiresAt,omitempty"`
}

// ShortenRequest represents the request to create a shortened URL
type ShortenRequest struct {
	OriginalURL  string `json:"originalUrl" binding:"required"`
	CustomCode   string `json:"customCode,omitempty"`
	ExpiresAfter int    `json:"expiresAfter,omitempty"` // In days
}

type ShortenData struct {
	Shorten  *Shorten `json:"shorten"`
	ShortURL string   `json:"shortUrl"`
}

type GetByOriginalURLRequest struct {
	OriginalURL       string `json:"originalUrl" binding:"required"`
	CreateIfNotExists bool   `json:"createIfNotExists"`
	ExpiresAfter      int    `json:"expiresAfter,omitempty"`
	CustomCode        string `json:"customCode,omitempty"`
	// TODO: allow duplicates? like create more copies if someone wants multiple short urls to the same domain. ??
}
