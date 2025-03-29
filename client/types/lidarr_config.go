package types

// @Description Jellyfin media server configuration
type LidarrConfig struct {
	Enabled bool                 `json:"enabled" mapstructure:"enabled" example:"false"`
	BaseURL string               `json:"baseURL" mapstructure:"host" example:"http://localhost:8096" binding:"required_if=Enabled true"`
	APIKey  string               `json:"apiKey" mapstructure:"apiKey" example:"your-api-key" binding:"required_if=Enabled true"`
	SSL     bool                 `json:"ssl" mapstructure:"ssl" example:"false"`
	Type    AutomationClientType `json:"type" mapstructure:"type" default:"lidarr"`
}

func NewLidarrConfig() LidarrConfig {
	return LidarrConfig{
		Type: AutomationClientTypeLidarr,
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
