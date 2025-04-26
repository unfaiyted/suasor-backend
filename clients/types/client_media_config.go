package types

import (
	"suasor/clients/media/types"

	"fmt"
)

type ClientMediaConfig interface {
	ClientConfig
	isClientMediaConfig()
	GetClientType() ClientMediaType

	GetBaseURL() string
	SetBaseURL(baseURL string)
	GetAPIKey() string
	SetAPIKey(apiKey string)

	SupportsMovies() bool
	SupportsSeries() bool
	SupportsMusic() bool
	SupportsPlaylists() bool
	SupportsCollections() bool
	SupportsHistory() bool

	SupportsMediaType(mediaType types.MediaType) bool
}

type clientMediaConfig struct {
	ClientConfig
	ClientType ClientMediaType `json:"clientType"`
	BaseURL    string          `json:"baseURL" mapstructure:"baseURL" example:"http://localhost:8096"`
	APIKey     string          `json:"apiKey" mapstructure:"apiKey" example:"your-api-key" binding:"required_if=Enabled true"`
	SSL        bool            `json:"ssl" mapstructure:"ssl" example:"false"`
}

func NewClientMediaConfig(clientType ClientMediaType, category ClientCategory, name string, baseURL string, apiKey string, enabled bool, validateConn bool) ClientMediaConfig {
	return &clientMediaConfig{
		ClientConfig: NewClientConfig(clientType.AsGenericClient(), category, name, baseURL, enabled, validateConn),
		ClientType:   clientType,
		APIKey:       apiKey,
	}
}

func (clientMediaConfig) isClientMediaConfig() {}

func (c *clientMediaConfig) GetBaseURL() string {
	return c.BaseURL
}

// setBaseURL sets the base URL for the client
func (c *clientMediaConfig) SetBaseURL(baseURL string) {
	c.BaseURL = baseURL
}

func (c *clientMediaConfig) GetAPIKey() string {
	return c.APIKey
}

// setAPIKey sets the API key for the client
func (c *clientMediaConfig) SetAPIKey(apiKey string) {
	c.APIKey = apiKey
}

func (c *clientMediaConfig) GetClientType() ClientMediaType {
	return c.ClientType
}

func (c *clientMediaConfig) SetClientType(clientType ClientMediaType) {
	c.ClientType = clientType
}

func (c *clientMediaConfig) SupportsMovies() bool {
	return false
}
func (c *clientMediaConfig) SupportsSeries() bool {
	return false
}
func (c *clientMediaConfig) SupportsMusic() bool {
	return false
}
func (c *clientMediaConfig) SupportsPlaylists() bool {
	return false
}
func (c *clientMediaConfig) SupportsCollections() bool {
	return false
}
func (c *clientMediaConfig) SupportsHistory() bool {
	return false
}

func (c *clientMediaConfig) SupportsMediaType(mediaType types.MediaType) bool {

	switch mediaType {
	case types.MediaTypeMovie:
		fmt.Printf("SupportsMovies: %v\n", c.SupportsMovies())
		return c.SupportsMovies()
	case types.MediaTypeSeries, types.MediaTypeSeason, types.MediaTypeEpisode:
		return c.SupportsSeries()
	case types.MediaTypeTrack, types.MediaTypeArtist, types.MediaTypeAlbum:
		return c.SupportsMusic()
	case types.MediaTypePlaylist:
		return c.SupportsPlaylists()
	case types.MediaTypeCollection:
		return c.SupportsCollections()
	case types.MediaTypeAll:
		return c.SupportsMovies() || c.SupportsSeries() || c.SupportsMusic() || c.SupportsPlaylists() || c.SupportsCollections()
	default:
		return false
	}

}
