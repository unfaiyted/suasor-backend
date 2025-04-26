package types

import "encoding/json"

// @Description Radarr movie automation server configuration
type RadarrConfig struct {
	ClientAutomationConfig
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
	// Create a temporary type to avoid recursion
	type Alias RadarrConfig
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Ensure Type is always the correct constant
	return nil
}
