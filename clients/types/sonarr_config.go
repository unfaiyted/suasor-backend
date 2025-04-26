package types

// @Description Sonarr TV series automation server configuration
type SonarrConfig struct {
	ClientAutomationConfig
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
