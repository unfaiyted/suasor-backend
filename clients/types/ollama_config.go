package types

// @Description Ollama local AI service configuration
type OllamaConfig struct {
	AIClientConfig
	Model       string  `json:"model" mapstructure:"model" example:"llama2"`
	Temperature float64 `json:"temperature" mapstructure:"temperature" example:"0.7"`
}

func NewOllamaConfig(baseURL string, model string, temperature float64, enabled bool, validateConn bool) OllamaConfig {
	clientConfig := NewClientAIConfig(AIClientTypeOllama, ClientCategoryAI, "Ollama", baseURL, "", enabled, validateConn)
	return OllamaConfig{
		AIClientConfig: clientConfig,
		Model:          model,
		Temperature:    temperature,
	}
}

func (c *OllamaConfig) GetModel() string {
	return c.Model
}

func (c *OllamaConfig) GetTemperature() float64 {
	return c.Temperature
}

func (OllamaConfig) GetClientType() AIClientType {
	return AIClientTypeOllama
}

func (OllamaConfig) GetCategory() ClientCategory {
	return ClientCategoryAI
}

func (c *OllamaConfig) UnmarshalJSON(data []byte) error {
	return UnmarshalConfigJSON(data, c)
}
