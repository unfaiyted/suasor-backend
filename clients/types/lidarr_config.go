package types

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
	return UnmarshalConfigJSON(data, c)
}
