package types

// @Description Emby media server configuration
type RadarrConfig struct {
	Enabled bool                 `json:"enabled" mapstructure:"enabled" example:"false"`
	BaseURL string               `json:"baseURL" mapstructure:"baseURL" example:"http://localhost:8096"`
	APIKey  string               `json:"apiKey" mapstructure:"apiKey" example:"your-api-key" binding:"required_if=Enabled true"`
	SSL     bool                 `json:"ssl" mapstructure:"ssl" example:"false"`
	Type    AutomationClientType `json:"type" mapstructure:"type" default:"radarr"`
}

func NewRadarrConfig() RadarrConfig {
	return RadarrConfig{
		Type: AutomationClientTypeRadarr,
	}
}
func (RadarrConfig) isAutomationClientConfig() {}
func (RadarrConfig) isClientConfig()           {}

func (RadarrConfig) GetClientType() AutomationClientType {
	return AutomationClientTypeRadarr
}

func (RadarrConfig) GetCategory() ClientCategory {
	return ClientCategoryAutomation
}

func (RadarrConfig) SupportsMovies() bool {
	return true
}
func (RadarrConfig) SupportsSeries() bool {
	return false
}
func (RadarrConfig) SupportsMusic() bool {
	return false
}
