package requests

import (
	"suasor/client/types"
	client "suasor/client/types"
)

// ClientTestRequest is used for testing a client connection
type ClientTestRequest struct {
	BaseURL    string           `json:"baseUrl" binding:"required,url"`
	APIKey     string           `json:"apiKey" binding:"required"`
	ClientType types.ClientType `json:"clientType" binding:"required,oneof=radarr sonarr lidarr"`
}

// AutomationClientRequest is used for creating/updating a download client

// type AutomationClientRequest struct {
// 	Name       string                      `json:"name" binding:"required"`
// 	ClientType client.AutomationClientType `json:"clientType" binding:"required,oneof=radarr sonarr lidarr"`
// 	BaseURL    string                      `json:"baseUrl" binding:"required,url"`
// 	APIKey     string                      `json:"apiKey" binding:"required"`
// 	IsEnabled  bool                        `json:"isEnabled"`
// }

// MediaClientRequest is used for creating/updating a media client
type ClientRequest struct {
	Name       string                 `json:"name" binding:"required"`
	ClientType client.MediaClientType `json:"clientType" binding:"required,oneof=plex jellyfin emby subsonic"`
	Client     any                    `json:"client" gorm:"serializer:json"`
}

// // MediaClientTestRequest is used for testing a media client connection
// type MediaClientTestRequest struct {
// 	ClientType client.MediaClientType `json:"clientType" binding:"required,oneof=plex jellyfin emby subsonic"`
// 	Client     any                    `json:"client" gorm:"serializer:json"`
// }
