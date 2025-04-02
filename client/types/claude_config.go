package types

import "encoding/json"

// @Description Claude media server configuration
type ClaudeConfig struct {
	BaseAIClientConfig
}

func NewClaudeConfig() ClaudeConfig {
	return ClaudeConfig{
		BaseAIClientConfig: BaseAIClientConfig{
			BaseClientConfig: BaseClientConfig{
				Type: ClientTypeClaude,
			},
			ClientType: AIClientTypeClaude,
		},
	}
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

	// Ensure Type is always the correct constant
	c.ClientType = AIClientTypeClaude
	c.Type = ClientTypeClaude
	return nil
}
