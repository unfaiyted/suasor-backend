package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"suasor/clients/media/types"

	"fmt"
)

type ClientMediaConfig interface {
	ClientConfig
	isClientMediaConfig()
	GetClientType() ClientMediaType

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
	ClientConfig `json:"core"`
	ClientType   ClientMediaType `json:"clientType"`
	APIKey       string          `json:"apiKey,omitempty" mapstructure:"apiKey" example:"your-api-key" `
	SSL          bool            `json:"ssl" mapstructure:"ssl" example:"false"`
}

func NewClientMediaConfig(clientType ClientMediaType, category ClientCategory, name string, baseURL string, apiKey string, enabled bool, validateConn bool) ClientMediaConfig {
	return &clientMediaConfig{
		ClientConfig: NewClientConfig(clientType.AsGenericClient(), category, name, baseURL, enabled, validateConn),
		ClientType:   clientType,
		APIKey:       apiKey,
	}
}

func (clientMediaConfig) isClientMediaConfig() {}

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

// Value implements driver.Valuer for database storage
func (c *clientMediaConfig) Value() (driver.Value, error) {
	// Serialize the entire item to JSON for storage
	return json.Marshal(c)
}

// Scan implements sql.Scanner for database retrieval
func (m *clientMediaConfig) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	// Use our custom unmarshaling
	err := m.UnmarshalJSON(bytes)
	if err != nil {
		return err
	}

	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (m *clientMediaConfig) UnmarshalJSON(data []byte) error {
	// Create a temporary struct without the embedded interface
	type Alias clientMediaConfig
	temp := struct {
		Core json.RawMessage `json:"core"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	// Unmarshal the basic fields
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Handle the ClientConfig by creating a concrete instance
	if len(temp.Core) > 0 {
		baseConfig := clientConfig{}
		if err := json.Unmarshal(temp.Core, &baseConfig); err != nil {
			return err
		}
		m.ClientConfig = &baseConfig
	} else {
		// If no base config provided, create a default one
		m.ClientConfig = &clientConfig{
			Type:     m.ClientType.AsGenericClient(),
			Category: ClientCategoryMedia,
			Name:     "Default Client",
			Enabled:  true,
		}
	}

	return nil
}
