package types

import "fmt"

// TMDBConfig represents the configuration for a TMDB client
type TMDBConfig struct {
	BaseMetadataClientConfig
	ApiKey string `json:"apiKey"`
}

// Validate checks if the provided config is valid
func (t *TMDBConfig) Validate() error {
	if err := t.BaseMetadataClientConfig.Validate(); err != nil {
		return err
	}

	if t.ApiKey == "" {
		return fmt.Errorf("missing required field: apiKey")
	}

	return nil
}

// GetMetadataClientType returns the metadata client type
func (t *TMDBConfig) GetMetadataClientType() MetadataClientType {
	return MetadataClientTypeTMDB
}

// GetClientType returns the client type
func (t *TMDBConfig) GetClientType() ClientType {
	return ClientTypeTMDB
}

// NewTMDBConfig creates a new TMDB client configuration
func NewTMDBConfig() *TMDBConfig {
	return &TMDBConfig{
		BaseMetadataClientConfig: BaseMetadataClientConfig{
			SupportsMovies:      true,
			SupportsTV:          true,
			SupportsPersons:     true,
			SupportsCollections: true,
		},
	}
}