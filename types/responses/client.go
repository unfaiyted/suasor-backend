package responses

import (
	"suasor/client/types"
	"time"
)

// MediaClientResponse is a non-generic representation of MediaClient for API responses
type ClientResponse struct {
	ID         uint64                `json:"id"`
	UserID     uint64                `json:"userId"`
	Name       string                `json:"name"`
	ClientType types.MediaClientType `json:"clientType"`
	Client     any                   `json:"client"`
	CreatedAt  time.Time             `json:"createdAt"`
	UpdatedAt  time.Time             `json:"updatedAt"`
}

// MediaClientTestResponse is the response for a media client connection test
type ClientTestResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Version string `json:"version,omitempty"`
}
