package models

import "time"

// MediaClientType represents different types of media clients
type MediaClientType string

const (
	MediaClientTypePlex     MediaClientType = "plex"
	MediaClientTypeJellyfin MediaClientType = "jellyfin"
	MediaClientTypeEmby     MediaClientType = "emby"
	MediaClientTypeSubsonic MediaClientType = "subsonic"
)

type ClientConfig interface {
	PlexConfig | EmbyConfig | JellyfinConfig | NavidromeConfig
}

// MediaClient represents a media client configuration
type MediaClient[T ClientConfig] struct {
	ID         uint64          `json:"id" gorm:"primaryKey"`
	UserID     uint64          `json:"userId" gorm:"not null"`
	Name       string          `json:"name" gorm:"not null"`
	ClientType MediaClientType `json:"clientType" gorm:"not null"`
	Client     T               `json:"client" gorm:"serializer:json"`
	CreatedAt  time.Time       `json:"createdAt"`
	UpdatedAt  time.Time       `json:"updatedAt"`
}

// MediaClientResponse is a non-generic representation of MediaClient for API responses
type MediaClientResponse struct {
	ID         uint64          `json:"id"`
	UserID     uint64          `json:"userId"`
	Name       string          `json:"name"`
	ClientType MediaClientType `json:"clientType"`
	Client     any             `json:"client"`
	CreatedAt  time.Time       `json:"createdAt"`
	UpdatedAt  time.Time       `json:"updatedAt"`
}

// ToResponse converts a generic MediaClient[T] to a MediaClientResponse
func ToResponse[T ClientConfig](client *MediaClient[T]) MediaClientResponse {
	return MediaClientResponse{
		ID:         client.ID,
		UserID:     client.UserID,
		Name:       client.Name,
		ClientType: client.ClientType,
		Client:     client.Client,
		CreatedAt:  client.CreatedAt,
		UpdatedAt:  client.UpdatedAt,
	}
}

// MediaClientRequest is used for creating/updating a media client
type MediaClientRequest struct {
	Name       string          `json:"name" binding:"required"`
	ClientType MediaClientType `json:"clientType" binding:"required,oneof=plex jellyfin emby subsonic"`
	Client     any             `json:"client" gorm:"serializer:json"`
}

// MediaClientTestRequest is used for testing a media client connection
type MediaClientTestRequest struct {
	ClientType MediaClientType `json:"clientType" binding:"required,oneof=plex jellyfin emby subsonic"`
	Client     any             `json:"client" gorm:"serializer:json"`
}

// MediaClientTestResponse is the response for a media client connection test
type MediaClientTestResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Version string `json:"version,omitempty"`
}
