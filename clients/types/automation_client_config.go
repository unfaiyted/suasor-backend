package types

type ClientAutomationConfig interface {
	ClientConfig
	isAutomationClientConfig()
	GetClientType() AutomationClientType

	SupportsMovies() bool
	SupportsSeries() bool
	SupportsMusic() bool
}

type clientAutomationConfig struct {
	ClientConfig
	BaseURL string `json:"baseURL" mapstructure:"baseURL" example:"http://localhost:8096"`
	APIKey  string `json:"apiKey" mapstructure:"apiKey" example:"your-api-key" binding:"required_if=Enabled true"`
	SSL     bool   `json:"ssl" mapstructure:"ssl" example:"false"`

	ClientType AutomationClientType `json:"clientType"`
}

func NewClientAutomationConfig(clientType AutomationClientType, category ClientCategory, name string, baseURL string, apiKey string, enabled bool, validateConn bool) ClientAutomationConfig {
	return &clientAutomationConfig{
		ClientConfig: NewClientConfig(clientType.AsGenericClient(), category, name, baseURL, enabled, validateConn),
		ClientType:   clientType,
		APIKey:       apiKey,
	}
}

func (clientAutomationConfig) isAutomationClientConfig() {}

func (c *clientAutomationConfig) GetClientType() AutomationClientType {
	return c.ClientType
}

func (c *clientAutomationConfig) GetCategory() ClientCategory {
	return ClientCategoryAutomation
}

func (c *clientAutomationConfig) GetBaseURL() string {
	return c.BaseURL
}

func (c *clientAutomationConfig) GetAPIKey() string {
	return c.APIKey
}

func (c *clientAutomationConfig) SupportsMovies() bool {
	return false
}
func (c *clientAutomationConfig) SupportsSeries() bool {
	return false
}
func (c *clientAutomationConfig) SupportsMusic() bool {
	return false
}
