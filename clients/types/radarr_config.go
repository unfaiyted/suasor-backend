package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// @Description Radarr movie automation server configuration
type RadarrConfig struct {
	ClientAutomationConfig `json:"details"`
	// Add any Radarr-specific configuration fields here
}

func NewRadarrConfig(baseURL string, apiKey string, enabled bool, validateConn bool) RadarrConfig {
	clientConfig := NewClientAutomationConfig(AutomationClientTypeRadarr, ClientCategoryAutomation, "Radarr", baseURL, apiKey, enabled, validateConn)
	return RadarrConfig{
		ClientAutomationConfig: clientConfig,
	}
}

func (RadarrConfig) GetClientType() AutomationClientType {
	return AutomationClientTypeRadarr
}

func (RadarrConfig) GetCategory() ClientCategory {
	return ClientCategoryAutomation
}

func (RadarrConfig) SupportsMovies() bool {
	return true
}

func (c *RadarrConfig) UnmarshalJSON(data []byte) error {
	return UnmarshalConfigJSON(data, c)
}

// Value implements driver.Valuer for database storage
func (c *RadarrConfig) Value() (driver.Value, error) {
	// Serialize the entire item to JSON for storage
	return json.Marshal(c)
}

// Scan implements sql.Scanner for database retrieval
func (m *RadarrConfig) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	// Use the same custom unmarshaling logic we defined in UnmarshalJSON
	return m.UnmarshalJSON(bytes)
}
