package types

type AIClientConfig interface {
	ClientConfig
	isAutomationClientConfig()
	GetClientType() AIClientType
	GetAPIKey() string
}

type BaseAIClientConfig struct {
	BaseClientConfig
	ClientType       AIClientType `json:"clientType"`
	BaseURL          string       `json:"baseURL" mapstructure:"baseURL" example:"http://localhost:8096"`
	APIKey           string       `json:"apiKey" mapstructure:"apiKey" example:"your-api-key" binding:"required_if=Enabled true"`
	Model            string       `json:"model" mapstructure:"model" example:"claude-2"`
	Temperature      float64      `json:"temperature" mapstructure:"temperature" example:"0.5"`
	MaxTokens        int          `json:"maxTokens" mapstructure:"maxTokens" example:"100"`
	MaxContextTokens int          `json:"maxContextTokens" mapstructure:"maxContextTokens" example:"1000"`
}

func (c *BaseAIClientConfig) GetClientType() AIClientType {
	return c.ClientType
}

func (c *BaseAIClientConfig) GetCategory() ClientCategory {
	c.Category = ClientCategoryAI
	return c.Category
}

func (c *BaseAIClientConfig) GetAPIKey() string {
	return c.APIKey
}
