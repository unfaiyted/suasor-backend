package responses

import (
	"suasor/client/types"
	"time"
)

// MediaClientResponse is a non-generic representation of MediaClient for API responses
type ClientResponse struct {
	ID         uint64                `json:"id" example:"1"`
	UserID     uint64                `json:"userId" example:"123"`
	Name       string                `json:"name" example:"My Plex Server"`
	ClientType types.MediaClientType `json:"clientType" example:"plex"`
	IsEnabled  bool                  `json:"isEnabled"`
	Client     any                   `json:"client"` // Can be any of the config types
	CreatedAt  time.Time             `json:"createdAt" example:"2023-01-01T12:00:00Z"`
	UpdatedAt  time.Time             `json:"updatedAt" example:"2023-01-01T12:00:00Z"`
}

// MediaClientTestResponse is the response for a media client connection test
type ClientTestResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Version string `json:"version,omitempty"`
}
