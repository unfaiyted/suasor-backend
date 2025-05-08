package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// @Description Ollama local AI service configuration
type OllamaConfig struct {
	AIClientConfig `json:"details"`
	Model          string  `json:"model" mapstructure:"model" example:"llama2"`
	Temperature    float64 `json:"temperature" mapstructure:"temperature" example:"0.7"`
}

func NewOllamaConfig(baseURL string, model string, temperature float64, enabled bool, validateConn bool) OllamaConfig {
	// Default reasonable values for tokens since Ollama doesn't explicitly set them
	maxTokens := 2048
	maxContextTokens := 4096
	
	clientConfig := NewClientAIConfig(AIClientTypeOllama, ClientCategoryAI, "Ollama", baseURL, "", model, temperature, maxTokens, maxContextTokens, enabled, validateConn)
	return OllamaConfig{
		AIClientConfig: clientConfig,
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

func (c *OllamaConfig) Value() (driver.Value, error) {
	// Serialize the entire item to JSON for storage
	return json.Marshal(c)
}

// Scan implements sql.Scanner for database retrieval
func (m *OllamaConfig) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	// Use the same custom unmarshaling logic we defined in UnmarshalJSON
	return m.UnmarshalJSON(bytes)
}
