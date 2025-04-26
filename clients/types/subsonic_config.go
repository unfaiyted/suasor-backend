package types

import (
	"suasor/clients/media/types"
)

// @Description Supersonic music server configuration
type SubsonicConfig struct {
	ClientMediaConfig
	Username string `json:"username" mapstructure:"username" example:"admin" binding:"required_if=Enabled true"`
	Password string `json:"password" mapstructure:"password" example:"your-password" binding:"required_if=Enabled true"`
}

func NewSubsonicConfig(username string, password string, baseURL string, enabled bool, validateConn bool) SubsonicConfig {
	clientConfig := NewClientMediaConfig(ClientMediaTypeSubsonic, ClientCategoryMedia, "Subsonic", baseURL, "", enabled, validateConn)
	return SubsonicConfig{
		ClientMediaConfig: clientConfig,
		Username:          username,
		Password:          password,
	}
}

func (c *SubsonicConfig) GetUsername() string {
	return c.Username
}

func (c *SubsonicConfig) GetPassword() string {
	return c.Password
}
func (SubsonicConfig) GetClientType() ClientMediaType {
	return ClientMediaTypeSubsonic
}
func (SubsonicConfig) GetCategory() ClientCategory {
	return ClientCategoryMedia
}

func (SubsonicConfig) SupportsMovies() bool {
	return false
}
func (SubsonicConfig) SupportsSeries() bool {
	return false
}
func (SubsonicConfig) SupportsMusic() bool {
	return true
}
func (SubsonicConfig) SupportsPlaylists() bool {
	return true
}
func (SubsonicConfig) SupportsCollections() bool {
	return false
}
func (SubsonicConfig) SupportsHistory() bool {
	return true
}

func (c *SubsonicConfig) UnmarshalJSON(data []byte) error {
	return UnmarshalConfigJSON(data, c)
}

func (c *SubsonicConfig) SupportsMediaType(mediaType types.MediaType) bool {
	return DoesClientSupportMediaType(c, mediaType)
}
