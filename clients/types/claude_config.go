package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// @Description Claude AI service configuration
type ClaudeConfig struct {
	AIClientConfig `json:"details"`
}

func NewClaudeConfig(apiKey string, baseURL string, model string, temperature float64, maxTokens int, maxContextTokens int, enabled bool, validateConn bool) ClaudeConfig {
	clientConfig := NewClientAIConfig(AIClientTypeClaude, ClientCategoryAI, "Claude", baseURL, apiKey, model, temperature, maxTokens, maxContextTokens, enabled, validateConn)
	return ClaudeConfig{
		AIClientConfig: clientConfig,
	}
}

func (ClaudeConfig) GetClientType() AIClientType {
	return AIClientTypeClaude
}
func (ClaudeConfig) GetCategory() ClientCategory {
	return ClientCategoryAI
}

func (c *ClaudeConfig) UnmarshalJSON(data []byte) error {
	return UnmarshalConfigJSON(data, c)
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
