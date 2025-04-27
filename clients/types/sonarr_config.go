package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// @Description Sonarr TV series automation server configuration
type SonarrConfig struct {
	ClientAutomationConfig `json:"details"`
	// Add any Sonarr-specific configuration fields here
}

func NewSonarrConfig(baseURL string, apiKey string, enabled bool, validateConn bool) SonarrConfig {
	clientConfig := NewClientAutomationConfig(AutomationClientTypeSonarr, ClientCategoryAutomation, "Sonarr", baseURL, apiKey, enabled, validateConn)
	return SonarrConfig{
		ClientAutomationConfig: clientConfig,
	}
}

func (SonarrConfig) GetClientType() AutomationClientType {
	return AutomationClientTypeSonarr
}

func (SonarrConfig) GetCategory() ClientCategory {
	return ClientCategoryAutomation
}

func (SonarrConfig) SupportsSeries() bool {
	return true
}

func (c *SonarrConfig) UnmarshalJSON(data []byte) error {
	return UnmarshalConfigJSON(data, c)
}

// Value implements driver.Valuer for database storage
func (c *SonarrConfig) Value() (driver.Value, error) {
	// Serialize the entire item to JSON for storage
	return json.Marshal(c)
}

// Scan implements sql.Scanner for database retrieval
func (m *SonarrConfig) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	// Use the same custom unmarshaling logic we defined in UnmarshalJSON
	return m.UnmarshalJSON(bytes)
}
