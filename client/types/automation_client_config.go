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

	ClientType AutomationClientType `json:"clientType"`
}

func (c *BaseAutomationClientConfig) GetClientType() AutomationClientType {
	return c.ClientType
}

func (c *BaseAutomationClientConfig) GetCategory() ClientCategory {
	c.Category = ClientCategoryAutomation
	return c.Category
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
