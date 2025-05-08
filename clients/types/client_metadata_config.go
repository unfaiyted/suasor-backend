package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
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
	ClientType          ClientType `json:"type"`
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
	return m.ClientType
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

// Value implements driver.Valuer for database storage
func (c *clientMetadataConfig) Value() (driver.Value, error) {
	// Serialize the entire item to JSON for storage
	return json.Marshal(c)
}

// Scan implements sql.Scanner for database retrieval
func (m *clientMetadataConfig) Scan(value any) error {
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
func (m *clientMetadataConfig) UnmarshalJSON(data []byte) error {
	// Create a temporary struct without the embedded interface
	type Alias clientMetadataConfig
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
			Type:     m.ClientType,
			Category: ClientCategoryMedia,
			Name:     "Default Client",
			Enabled:  true,
		}
	}

	return nil
}
