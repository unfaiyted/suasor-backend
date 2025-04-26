package requests

import (
	mediatypes "suasor/clients/media/types"
	"suasor/clients/types"
	client "suasor/clients/types"
)

// ClientTestRequest is used for testing a client connection
type ClientTestRequest[T types.ClientConfig] struct {
	ClientType types.ClientType `json:"clientType" binding:"required,oneof=radarr sonarr lidarr emby jellyfin subsonic plex claude openai ollama"`
	Client     T                `json:"client" gorm:"serializer:json"`
}

// AutomationClientRequest is used for creating/updating a download client
type AutomationClientRequest[T types.ClientAutomationConfig] struct {
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
	Config     T                 `json:"config" gorm:"serializer:jsonb"`
}

// ClientMediaRequest is used for testing a media client connection
type ClientMediaRequest[T types.ClientMediaConfig] struct {
	Name       string                 `json:"name" binding:"required"`
	ClientType client.ClientMediaType `json:"clientType" binding:"required,oneof=plex jellyfin emby subsonic"`
	Config     T                      `json:"config" gorm:"serializer:jsonb"`
}

type ClientMediaItemUpdateRequest[T mediatypes.MediaData] struct {
	ID         uint64               `json:"id" binding:"required"`
	Type       mediatypes.MediaType `json:"type" binding:"required"`
	ClientID   uint64               `json:"clientID" binding:"required"`
	ClientType client.ClientType    `json:"clientType" binding:"required"`
	ExternalID string               `json:"externalId" binding:"required"`
	Data       T                    `json:"data" binding:"required"`
}

type ClientMediaItemCreateRequest[T mediatypes.MediaData] struct {
	Type       mediatypes.MediaType `json:"type" binding:"required"`
	ClientID   uint64               `json:"clientID" binding:"required"`
	ClientType client.ClientType    `json:"clientType" binding:"required"`
	ExternalID string               `json:"externalId" binding:"required"`
	Data       T                    `json:"data" binding:"required"`
}
