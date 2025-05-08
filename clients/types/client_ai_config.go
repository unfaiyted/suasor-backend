package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type AIClientConfig interface {
	ClientConfig
	GetClientType() AIClientType
	GetAPIKey() string
	GetModel() string
	GetTemperature() float64
	GetMaxTokens() int
	GetMaxContextTokens() int
}

type clientAIConfig struct {
	ClientConfig     `json:"core"`
	ClientType       AIClientType `json:"clientType"`
	APIKey           string       `json:"apiKey" mapstructure:"apiKey" example:"your-api-key" binding:"required_if=Enabled true"`
	Model            string       `json:"model" mapstructure:"model" example:"claude-2"`
	Temperature      float64      `json:"temperature" mapstructure:"temperature" example:"0.5"`
	MaxTokens        int          `json:"maxTokens" mapstructure:"maxTokens" example:"100"`
	MaxContextTokens int          `json:"maxContextTokens" mapstructure:"maxContextTokens" example:"1000"`
}

func NewClientAIConfig(clientType AIClientType, category ClientCategory, name string, baseURL string, apiKey string, model string, temperature float64, maxTokens int, maxContextTokens int, enabled bool, validateConn bool) AIClientConfig {
	return &clientAIConfig{
		ClientConfig:     NewClientConfig(clientType.AsGenericClient(), category, name, baseURL, enabled, validateConn),
		ClientType:       clientType,
		APIKey:           apiKey,
		Model:            model,
		Temperature:      temperature,
		MaxTokens:        maxTokens,
		MaxContextTokens: maxContextTokens,
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

func (c *clientAIConfig) GetModel() string {
	return c.Model
}

func (c *clientAIConfig) GetTemperature() float64 {
	return c.Temperature
}

func (c *clientAIConfig) GetMaxTokens() int {
	return c.MaxTokens
}

func (c *clientAIConfig) GetMaxContextTokens() int {
	return c.MaxContextTokens
}

// Value implements driver.Valuer for database storage
func (c *clientAIConfig) Value() (driver.Value, error) {
	// Serialize the entire item to JSON for storage
	return json.Marshal(c)
}

// Scan implements sql.Scanner for database retrieval
func (m *clientAIConfig) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	// Use our custom unmarshaling
	err := m.UnmarshalJSON(bytes)
	if err != nil {
		return err
	}

	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (m *clientAIConfig) UnmarshalJSON(data []byte) error {
	// Create a temporary struct without the embedded interface
	type Alias clientAIConfig
	temp := struct {
		Core json.RawMessage `json:"core"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	// Unmarshal the basic fields
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Handle the ClientConfig by creating a concrete instance
	if len(temp.Core) > 0 {
		baseConfig := clientConfig{}
		if err := json.Unmarshal(temp.Core, &baseConfig); err != nil {
			return err
		}
		m.ClientConfig = &baseConfig
	} else {
		// If no base config provided, create a default one
		m.ClientConfig = &clientConfig{
			Type:     m.ClientType.AsGenericClient(),
			Category: ClientCategoryAI,
			Name:     "Default Client",
			Enabled:  true,
		}
	}

	return nil
}
