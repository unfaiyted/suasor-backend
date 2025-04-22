package types

import "encoding/json"

// @Description Emby media server configuration
type RadarrConfig struct {
	BaseAutomationClientConfig
}

func NewRadarrConfig() RadarrConfig {
	return RadarrConfig{
		BaseAutomationClientConfig: BaseAutomationClientConfig{
			BaseClientConfig: BaseClientConfig{
				Type: ClientTypeRadarr,
			},
			ClientType: AutomationClientTypeRadarr,
		},
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
	c.ClientType = AutomationClientTypeRadarr
	c.Type = ClientTypeRadarr
	return nil
}
