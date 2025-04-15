package requests

import (
	"suasor/client/types"
	client "suasor/client/types"
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

// SwaggerClientRequest is a concrete version of ClientRequest for swagger docs
type SwaggerClientRequest struct {
	Name       string            `json:"name" binding:"required"`
	ClientID   uint64            `json:"clientID,omitempty"`
	ClientType client.ClientType `json:"clientType" binding:"required"`
	IsEnabled  bool              `json:"isEnabled"`
	Client     interface{}       `json:"client"`
}

// ClientMediaRequest is used for testing a media client connection
type ClientMediaRequest[T types.ClientMediaConfig] struct {
	Name       string                 `json:"name" binding:"required"`
	ClientType client.ClientMediaType `json:"clientType" binding:"required,oneof=plex jellyfin emby subsonic"`
	Client     T                      `json:"client" gorm:"serializer:json"`
}

// SwaggerClientMediaRequest is a concrete version of ClientMediaRequest for swagger docs
type SwaggerClientMediaRequest struct {
	Name       string                 `json:"name" binding:"required"`
	ClientType client.ClientMediaType `json:"clientType" binding:"required,oneof=plex jellyfin emby subsonic"`
	Client     interface{}            `json:"client"`
}
