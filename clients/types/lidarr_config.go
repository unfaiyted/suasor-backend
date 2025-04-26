package types

import "encoding/json"

// @Description Lidarr automation server configuration
type LidarrConfig struct {
	ClientAutomationConfig
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
	// Create a temporary type to avoid recursion
	type Alias LidarrConfig
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	return nil
}
