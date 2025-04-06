package types

// MetadataClientConfig is the interface for metadata client configurations
type MetadataClientConfig interface {
	ClientConfig

	// Supported features
	SupportsMovieMetadata() bool
	SupportsTVMetadata() bool
	SupportsPersonMetadata() bool
	SupportsCollectionMetadata() bool
	
	// GetClientType returns the specific metadata client type
	GetMetadataClientType() MetadataClientType
}

// BaseMetadataClientConfig provides a base implementation of MetadataClientConfig
type BaseMetadataClientConfig struct {
	BaseClientConfig
	SupportsMovies         bool `json:"supportsMovies"`
	SupportsTV             bool `json:"supportsTV"`
	SupportsPersons        bool `json:"supportsPersons"`
	SupportsCollections    bool `json:"supportsCollections"`
}

// Validate validates the base metadata client config
func (m *BaseMetadataClientConfig) Validate() error {
	return m.BaseClientConfig.Validate()
}

func (m *BaseMetadataClientConfig) SupportsMovieMetadata() bool {
	return m.SupportsMovies
}

func (m *BaseMetadataClientConfig) SupportsTVMetadata() bool {
	return m.SupportsTV
}

func (m *BaseMetadataClientConfig) SupportsPersonMetadata() bool {
	return m.SupportsPersons
}

func (m *BaseMetadataClientConfig) SupportsCollectionMetadata() bool {
	return m.SupportsCollections
}