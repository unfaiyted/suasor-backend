package types

import (
	"suasor/clients/metadata/types"
)

// MetadataClientConfig is the interface for metadata client configurations
type ClientMetadataConfig interface {
	ClientConfig

	// Supported features
	SupportsMovieMetadata() bool
	SupportsTVMetadata() bool
	SupportsPersonMetadata() bool
	SupportsCollectionMetadata() bool
	GetClientType() ClientType

	SupportsMetadataType(metadataType types.MetadataType) bool
}

// clientMetadataConfig provides a base implementation of MetadataClientConfig
type clientMetadataConfig struct {
	ClientConfig        `json:"core"`
	Type                ClientType `json:"type"`
	SupportsMovies      bool       `json:"supportsMovies"`
	SupportsTV          bool       `json:"supportsTV"`
	SupportsPersons     bool       `json:"supportsPersons"`
	SupportsCollections bool       `json:"supportsCollections"`
}

func (m *clientMetadataConfig) SupportsMovieMetadata() bool {
	return m.SupportsMovies
}

func (m *clientMetadataConfig) SupportsTVMetadata() bool {
	return m.SupportsTV
}

func (m *clientMetadataConfig) SupportsPersonMetadata() bool {
	return m.SupportsPersons
}

func (m *clientMetadataConfig) SupportsCollectionMetadata() bool {
	return m.SupportsCollections
}

func (m *clientMetadataConfig) GetClientType() ClientType {
	return m.Type
}

func (m *clientMetadataConfig) SupportsMetadataType(metadataType types.MetadataType) bool {
	switch metadataType {
	case types.MetadataTypeMovie:
		return m.SupportsMovieMetadata()
	case types.MetadataTypeTV:
		return m.SupportsTVMetadata()
	case types.MetadataTypePerson:
		return m.SupportsPersonMetadata()
	case types.MetadataTypeCollection:
		return m.SupportsCollectionMetadata()
	default:
		return false
	}
}
