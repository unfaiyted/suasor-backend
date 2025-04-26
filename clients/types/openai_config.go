package types

// @Description OpenAI service configuration
type OpenAIConfig struct {
	AIClientConfig
	Model            string  `json:"model" mapstructure:"model" example:"gpt-4-turbo"`
	Temperature      float64 `json:"temperature" mapstructure:"temperature" example:"0.7"`
	MaxTokens        int     `json:"maxTokens" mapstructure:"maxTokens" example:"1000"`
	MaxContextTokens int     `json:"maxContextTokens" mapstructure:"maxContextTokens" example:"8192"`
}

func NewOpenAIConfig(apiKey string, baseURL string, model string, temperature float64, maxTokens int, maxContextTokens int, enabled bool, validateConn bool) OpenAIConfig {
	clientConfig := NewClientAIConfig(AIClientTypeOpenAI, ClientCategoryAI, "OpenAI", baseURL, apiKey, enabled, validateConn)
	return OpenAIConfig{
		AIClientConfig:   clientConfig,
		Model:            model,
		Temperature:      temperature,
		MaxTokens:        maxTokens,
		MaxContextTokens: maxContextTokens,
	}
}

func (c *OpenAIConfig) GetModel() string {
	return c.Model
}

func (c *OpenAIConfig) GetTemperature() float64 {
	return c.Temperature
}

func (c *OpenAIConfig) GetMaxTokens() int {
	return c.MaxTokens
}

func (c *OpenAIConfig) GetMaxContextTokens() int {
	return c.MaxContextTokens
}

func (OpenAIConfig) GetClientType() AIClientType {
	return AIClientTypeOpenAI
}

func (OpenAIConfig) GetCategory() ClientCategory {
	return ClientCategoryAI
}

func (c *OpenAIConfig) UnmarshalJSON(data []byte) error {
	return UnmarshalConfigJSON(data, c)
}
