package types

import (
	"suasor/clients/media/types"
)

// @Description Plex media server configuration
type PlexConfig struct {
	ClientMediaConfig
	Token string `json:"token" mapstructure:"token" example:"your-plex-token" binding:"required_if=Enabled true"`
}

func NewPlexConfig(host string, token string, enabled bool, validateConn bool) PlexConfig {
	clientConfig := NewClientMediaConfig(ClientMediaTypePlex, ClientCategoryMedia, "Plex", host, "", enabled, validateConn)
	return PlexConfig{
		ClientMediaConfig: clientConfig,
		Token:             token,
	}
}

func (c *PlexConfig) GetToken() string {
	return c.Token
}

func (PlexConfig) GetClientType() ClientMediaType {
	return ClientMediaTypePlex
}
func (PlexConfig) GetCategory() ClientCategory {
	return ClientCategoryMedia
}

func (PlexConfig) SupportsMovies() bool {
	return true
}
func (PlexConfig) SupportsSeries() bool {
	return true
}
func (PlexConfig) SupportsMusic() bool {
	return true
}
func (PlexConfig) SupportsPlaylists() bool {
	return true
}
func (PlexConfig) SupportsCollections() bool {
	return true
}
func (PlexConfig) SupportsHistory() bool {
	return true
}

func (c *PlexConfig) UnmarshalJSON(data []byte) error {
	return UnmarshalConfigJSON(data, c)
}
func (c *PlexConfig) SupportsMediaType(mediaType types.MediaType) bool {
	return DoesClientSupportMediaType(c, mediaType)
}
