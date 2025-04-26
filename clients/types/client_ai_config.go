package types

type AIClientConfig interface {
	ClientConfig
	GetClientType() AIClientType
	GetAPIKey() string
}

type clientAIConfig struct {
	ClientConfig
	ClientType       AIClientType `json:"clientType"`
	BaseURL          string       `json:"baseURL" mapstructure:"baseURL" example:"http://localhost:8096"`
	APIKey           string       `json:"apiKey" mapstructure:"apiKey" example:"your-api-key" binding:"required_if=Enabled true"`
	Model            string       `json:"model" mapstructure:"model" example:"claude-2"`
	Temperature      float64      `json:"temperature" mapstructure:"temperature" example:"0.5"`
	MaxTokens        int          `json:"maxTokens" mapstructure:"maxTokens" example:"100"`
	MaxContextTokens int          `json:"maxContextTokens" mapstructure:"maxContextTokens" example:"1000"`
}

func NewClientAIConfig(clientType AIClientType, category ClientCategory, name string, baseURL string, apiKey string, enabled bool, validateConn bool) AIClientConfig {
	return &clientAIConfig{
		ClientConfig: NewClientConfig(clientType.AsGenericClient(), category, name, baseURL, enabled, validateConn),
		ClientType:   clientType,
		APIKey:       apiKey,
	}
}

func (c *clientAIConfig) GetClientType() AIClientType {
	return c.ClientType
}

func (c *clientAIConfig) GetCategory() ClientCategory {
	return ClientCategoryAI
}

func (c *clientAIConfig) GetAPIKey() string {
	return c.APIKey
}
