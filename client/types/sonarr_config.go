package types

// @Description Emby media server configuration
type SonarrConfig struct {
	Enabled bool                 `json:"enabled" mapstructure:"enabled" example:"false"`
	BaseURL string               `json:"baseURL" mapstructure:"baseURL" example:"http://localhost:8096"`
	APIKey  string               `json:"apiKey" mapstructure:"apiKey" example:"your-api-key" binding:"required_if=Enabled true"`
	SSL     bool                 `json:"ssl" mapstructure:"ssl" example:"false"`
	Type    AutomationClientType `json:"type" mapstructure:"type" default:"sonarr"`
}

func NewSonarrConfig() SonarrConfig {
	return SonarrConfig{
		Type: AutomationClientTypeSonarr,
	}
}
func (SonarrConfig) isClientConfig()           {}
func (SonarrConfig) isAutomationClientConfig() {}

func (SonarrConfig) GetClientType() AutomationClientType {
	return AutomationClientTypeSonarr
}

func (SonarrConfig) GetCategory() ClientCategory {
	return ClientCategoryAutomation
}

func (SonarrConfig) SupportsMovies() bool {
	return false
}
func (SonarrConfig) SupportsSeries() bool {
	return true
}
func (SonarrConfig) SupportsMusic() bool {
	return false
}
