package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// @Description Lidarr automation server configuration
type LidarrConfig struct {
	ClientAutomationConfig `json:"details"`
}

func NewLidarrConfig(baseURL string, apiKey string, enabled bool, validateConn bool) LidarrConfig {

	clientConfig := NewClientAutomationConfig(
		AutomationClientTypeLidarr, ClientCategoryAutomation,
		"Lidarr", baseURL,
		apiKey, enabled, validateConn)

	return LidarrConfig{
		ClientAutomationConfig: clientConfig,
	}
}

func (LidarrConfig) GetClientType() AutomationClientType {
	return AutomationClientTypeLidarr
}

func (LidarrConfig) GetCategory() ClientCategory {
	return ClientCategoryAutomation
}

func (LidarrConfig) SupportsMovies() bool {
	return false
}
func (LidarrConfig) SupportsSeries() bool {
	return false
}
func (LidarrConfig) SupportsMusic() bool {
	return true
}

func (c *LidarrConfig) UnmarshalJSON(data []byte) error {
	return UnmarshalConfigJSON(data, c)
}

// Value implements driver.Valuer for database storage
func (c *LidarrConfig) Value() (driver.Value, error) {
	// Serialize the entire item to JSON for storage
	return json.Marshal(c)
}

// Scan implements sql.Scanner for database retrieval
func (m *LidarrConfig) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	// Use the same custom unmarshaling logic we defined in UnmarshalJSON
	return m.UnmarshalJSON(bytes)
}
