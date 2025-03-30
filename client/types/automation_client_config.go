package types

type AutomationClientConfig interface {
	ClientConfig
	isAutomationClientConfig()
	GetClientType() AutomationClientType

	SupportsMovies() bool
	SupportsSeries() bool
	SupportsMusic() bool
}

type BaseAutomationClientConfig struct {
	BaseClientConfig
	BaseURL string `json:"baseURL" mapstructure:"baseURL" example:"http://localhost:8096"`
	APIKey  string `json:"apiKey" mapstructure:"apiKey" example:"your-api-key" binding:"required_if=Enabled true"`
	SSL     bool   `json:"ssl" mapstructure:"ssl" example:"false"`

	Type AutomationClientType
}

func (c *BaseAutomationClientConfig) GetClientType() AutomationClientType {
	return c.Type
}

func (c *BaseAutomationClientConfig) GetName() string {
	return c.Name
}

func (c *BaseAutomationClientConfig) GetCategory() ClientCategory {
	return ClientCategoryAutomation
}

func (c *BaseAutomationClientConfig) SupportsMovies() bool {
	return false
}
func (c *BaseAutomationClientConfig) SupportsSeries() bool {
	return false
}
func (c *BaseAutomationClientConfig) SupportsMusic() bool {
	return false
}
