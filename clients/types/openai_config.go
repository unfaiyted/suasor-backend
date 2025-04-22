package types

import "encoding/json"

// @Description Claude media server configuration
type OpenAIConfig struct {
	BaseAIClientConfig
}

func NewOpenAIConfig() OpenAIConfig {
	return OpenAIConfig{
		BaseAIClientConfig: BaseAIClientConfig{
			BaseClientConfig: BaseClientConfig{
				Type: ClientTypeOpenAI,
			},
			ClientType: AIClientTypeOpenAI,
		},
	}
}

func (OpenAIConfig) GetClientType() AIClientType {
	return AIClientTypeClaude
}
func (OpenAIConfig) GetCategory() ClientCategory {
	return ClientCategoryMedia
}

func (c *OpenAIConfig) UnmarshalJSON(data []byte) error {
	// Create a temporary type to avoid recursion
	type Alias OpenAIConfig
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Ensure Type is always the correct constant
	c.ClientType = AIClientTypeOpenAI
	c.Type = ClientTypeOpenAI
	return nil
}
