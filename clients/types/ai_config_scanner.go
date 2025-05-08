package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// AIConfigScanner is a concrete implementation of AIClientConfig for database scanning
// This is needed because interfaces can't be directly scanned from the database
type AIConfigScanner struct {
	configRaw json.RawMessage
	config    AIClientConfig
}

// UnmarshalJSON implements json.Unmarshaler
func (s *AIConfigScanner) UnmarshalJSON(data []byte) error {
	// Store the raw data
	s.configRaw = data

	// Try to determine the type of AI client config
	var typeObj struct {
		Details struct {
			Core struct {
				Type ClientType `json:"type"`
			} `json:"core"`
			ClientType AIClientType `json:"clientType"`
		} `json:"details"`
	}

	if err := json.Unmarshal(data, &typeObj); err != nil {
		return fmt.Errorf("failed to unmarshal AI config type: %w", err)
	}

	// Create the appropriate concrete type
	var config AIClientConfig
	switch typeObj.Details.ClientType {
	case AIClientTypeClaude:
		var claudeConfig ClaudeConfig
		if err := json.Unmarshal(data, &claudeConfig); err != nil {
			return fmt.Errorf("failed to unmarshal Claude config: %w", err)
		}
		config = &claudeConfig
	case AIClientTypeOpenAI:
		var openaiConfig OpenAIConfig
		if err := json.Unmarshal(data, &openaiConfig); err != nil {
			return fmt.Errorf("failed to unmarshal OpenAI config: %w", err)
		}
		config = &openaiConfig
	case AIClientTypeOllama:
		var ollamaConfig OllamaConfig
		if err := json.Unmarshal(data, &ollamaConfig); err != nil {
			return fmt.Errorf("failed to unmarshal Ollama config: %w", err)
		}
		config = &ollamaConfig
	default:
		return fmt.Errorf("unknown AI client type: %s", typeObj.Details.ClientType)
	}

	s.config = config
	return nil
}

// Value implements driver.Valuer for database storage
func (s AIConfigScanner) Value() (driver.Value, error) {
	if s.config == nil {
		return nil, errors.New("nil config cannot be stored in database")
	}
	return json.Marshal(s.config)
}

// Scan implements sql.Scanner for database retrieval
func (s *AIConfigScanner) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return s.UnmarshalJSON(bytes)
}

// GetConfig returns the actual AIClientConfig implementation
func (s *AIConfigScanner) GetConfig() AIClientConfig {
	return s.config
}

// GetType returns the ClientType (implementing ClientConfig)
func (s *AIConfigScanner) GetType() ClientType {
	if s.config == nil {
		return ClientTypeUnknown
	}
	return s.config.GetType()
}

// GetCategory returns the ClientCategory (implementing ClientConfig)
func (s *AIConfigScanner) GetCategory() ClientCategory {
	if s.config == nil {
		return ClientCategoryUnknown
	}
	return s.config.GetCategory()
}

// GetClientType returns the AIClientType (implementing AIClientConfig)
func (s *AIConfigScanner) GetClientType() AIClientType {
	if s.config == nil {
		return AIClientTypeUnknown
	}
	return s.config.GetClientType()
}

// GetAPIKey returns the API key (implementing AIClientConfig)
func (s *AIConfigScanner) GetAPIKey() string {
	if s.config == nil {
		return ""
	}
	return s.config.GetAPIKey()
}

// GetModel returns the model name (implementing AIClientConfig)
func (s *AIConfigScanner) GetModel() string {
	if s.config == nil {
		return ""
	}
	return s.config.GetModel()
}

// GetTemperature returns the temperature (implementing AIClientConfig)
func (s *AIConfigScanner) GetTemperature() float64 {
	if s.config == nil {
		return 0
	}
	return s.config.GetTemperature()
}

// GetMaxTokens returns the max tokens (implementing AIClientConfig)
func (s *AIConfigScanner) GetMaxTokens() int {
	if s.config == nil {
		return 0
	}
	return s.config.GetMaxTokens()
}

// GetMaxContextTokens returns the max context tokens (implementing AIClientConfig)
func (s *AIConfigScanner) GetMaxContextTokens() int {
	if s.config == nil {
		return 0
	}
	return s.config.GetMaxContextTokens()
}