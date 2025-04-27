package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// TMDBConfig represents the configuration for a TMDB client
type TMDBConfig struct {
	ClientMetadataConfig `json:"details"`
	ApiKey               string `json:"apiKey" mapstructure:"apiKey" example:"your-tmdb-api-key" binding:"required_if=Enabled true"`
}

// Validate checks if the provided config is valid
func (t *TMDBConfig) Validate() error {
	if t.ApiKey == "" {
		return fmt.Errorf("missing required field: apiKey")
	}

	return nil
}

// GetApiKey returns the API key
func (t *TMDBConfig) GetApiKey() string {
	return t.ApiKey
}

// GetMetadataClientType returns the metadata client type
func (t *TMDBConfig) GetMetadataClientType() MetadataClientType {
	return MetadataClientTypeTMDB
}

// GetClientType returns the client type
func (t *TMDBConfig) GetClientType() ClientType {
	return ClientTypeTMDB
}

// GetCategory returns the client category
func (t *TMDBConfig) GetCategory() ClientCategory {
	return ClientCategoryMetadata
}

// Feature support methods
func (t *TMDBConfig) SupportsMovieMetadata() bool {
	return true
}

func (t *TMDBConfig) SupportsTVMetadata() bool {
	return true
}

func (t *TMDBConfig) SupportsPersonMetadata() bool {
	return true
}

func (t *TMDBConfig) SupportsCollectionMetadata() bool {
	return true
}

// NewTMDBConfig creates a new TMDB client configuration
func NewTMDBConfig(apiKey string, baseURL string, enabled bool, validateConn bool) *TMDBConfig {
	// Create a basic client config
	clientConfig := NewClientConfig(ClientTypeTMDB, ClientCategoryMetadata, "TMDB", baseURL, enabled, validateConn)

	// Create a client metadata config with the features supported
	metadataConfig := &clientMetadataConfig{
		ClientConfig:        clientConfig,
		SupportsMovies:      true,
		SupportsTV:          true,
		SupportsPersons:     true,
		SupportsCollections: true,
	}

	return &TMDBConfig{
		ClientMetadataConfig: metadataConfig,
		ApiKey:               apiKey,
	}
}

func (c *TMDBConfig) UnmarshalJSON(data []byte) error {
	return UnmarshalConfigJSON(data, c)
}

// Value implements driver.Valuer for database storage
func (c *TMDBConfig) Value() (driver.Value, error) {
	// Serialize the entire item to JSON for storage
	return json.Marshal(c)
}

// Scan implements sql.Scanner for database retrieval
func (m *TMDBConfig) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	// Use the same custom unmarshaling logic we defined in UnmarshalJSON
	return m.UnmarshalJSON(bytes)
}
