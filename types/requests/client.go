package requests

import (
	"suasor/clients/types"
	client "suasor/clients/types"
)

// ClientTestRequest is used for testing a client connection
type ClientTestRequest[T types.ClientConfig] struct {
	ClientType types.ClientType `json:"clientType" binding:"required,oneof=radarr sonarr lidarr emby jellyfin subsonic plex claude openai ollama"`
	Client     T                `json:"client" gorm:"serializer:json"`
}

// AutomationClientRequest is used for creating/updating a download client
type AutomationClientRequest[T types.AutomationClientConfig] struct {
	Name       string                      `json:"name" binding:"required"`
	ClientType client.AutomationClientType `json:"clientType" binding:"required,oneof=radarr sonarr lidarr"`
	IsEnabled  bool                        `json:"isEnabled"`
	Client     T                           `json:"client" gorm:"serializer:json"`
}

// ClientRequest is used for creating/updating a media client
type ClientRequest[T types.ClientConfig] struct {
	Name       string            `json:"name" binding:"required"`
	ClientID   uint64            `json:"clientID,omitempty"`
	ClientType client.ClientType `json:"clientType" binding:"required"`
	IsEnabled  bool              `json:"isEnabled"`
	Client     T                 `json:"client" gorm:"serializer:json"`
}

// ClientMediaRequest is used for testing a media client connection
type ClientMediaRequest[T types.ClientMediaConfig] struct {
	Name       string                 `json:"name" binding:"required"`
	ClientType client.ClientMediaType `json:"clientType" binding:"required,oneof=plex jellyfin emby subsonic"`
	Client     T                      `json:"client" gorm:"serializer:json"`
}
