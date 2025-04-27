package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// @Description Claude AI service configuration
type ClaudeConfig struct {
	AIClientConfig
	Model            string  `json:"model" mapstructure:"model" example:"claude-3-opus-20240229"`
	Temperature      float64 `json:"temperature" mapstructure:"temperature" example:"0.7"`
	MaxTokens        int     `json:"maxTokens" mapstructure:"maxTokens" example:"2000"`
	MaxContextTokens int     `json:"maxContextTokens" mapstructure:"maxContextTokens" example:"100000"`
}

func NewClaudeConfig(apiKey string, baseURL string, model string, temperature float64, maxTokens int, maxContextTokens int, enabled bool, validateConn bool) ClaudeConfig {
	clientConfig := NewClientAIConfig(AIClientTypeClaude, ClientCategoryAI, "Claude", baseURL, apiKey, enabled, validateConn)
	return ClaudeConfig{
		AIClientConfig:   clientConfig,
		Model:            model,
		Temperature:      temperature,
		MaxTokens:        maxTokens,
		MaxContextTokens: maxContextTokens,
	}
}

func (c *ClaudeConfig) GetModel() string {
	return c.Model
}

func (c *ClaudeConfig) GetTemperature() float64 {
	return c.Temperature
}

func (c *ClaudeConfig) GetMaxTokens() int {
	return c.MaxTokens
}

func (c *ClaudeConfig) GetMaxContextTokens() int {
	return c.MaxContextTokens
}

func (ClaudeConfig) GetClientType() AIClientType {
	return AIClientTypeClaude
}
func (ClaudeConfig) GetCategory() ClientCategory {
	return ClientCategoryAI
}

func (c *ClaudeConfig) UnmarshalJSON(data []byte) error {
	// Create a temporary type to avoid recursion
	type Alias ClaudeConfig
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	return nil
}

// Value implements driver.Valuer for database storage
func (c *ClaudeConfig) Value() (driver.Value, error) {
	// Serialize the entire item to JSON for storage
	return json.Marshal(c)
}

// Scan implements sql.Scanner for database retrieval
func (m *ClaudeConfig) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	// Use the same custom unmarshaling logic we defined in UnmarshalJSON
	return m.UnmarshalJSON(bytes)
}
