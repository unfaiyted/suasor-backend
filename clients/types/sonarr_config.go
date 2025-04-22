package types

import "encoding/json"

// @Description Emby media server configuration
type SonarrConfig struct {
	BaseAutomationClientConfig
}

func NewSonarrConfig() SonarrConfig {
	return SonarrConfig{
		BaseAutomationClientConfig: BaseAutomationClientConfig{
			BaseClientConfig: BaseClientConfig{
				Type: ClientTypeSonarr,
			},
			ClientType: AutomationClientTypeSonarr,
		},
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
	// Create a temporary type to avoid recursion
	type Alias SonarrConfig
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Ensure Type is always the correct constant
	c.ClientType = AutomationClientTypeSonarr
	c.Type = ClientTypeSonarr
	return nil
}
