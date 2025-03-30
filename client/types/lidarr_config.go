package types

import "encoding/json"

// @Description Jellyfin media server configuration
type LidarrConfig struct {
	BaseAutomationClientConfig
}

func NewLidarrConfig() LidarrConfig {
	return LidarrConfig{
		BaseAutomationClientConfig: BaseAutomationClientConfig{
			BaseClientConfig: BaseClientConfig{
				Type: ClientTypeLidarr,
			},
			Type: AutomationClientTypeLidarr,
		},
	}
}
func (LidarrConfig) isClientConfig()           {}
func (LidarrConfig) isAutomationClientConfig() {}

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

	// Ensure Type is always the correct constant
	c.Type = AutomationClientTypeLidarr
	return nil
}
