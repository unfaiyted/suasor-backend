package types

import "encoding/json"

// @Description Claude media server configuration
type OllamaConfig struct {
	BaseAIClientConfig
}

func NewOllamaConfig() OllamaConfig {
	return OllamaConfig{
		BaseAIClientConfig: BaseAIClientConfig{
			BaseClientConfig: BaseClientConfig{
				Type: ClientTypeOllama,
			},
			ClientType: AIClientTypeOllama,
		},
	}
}

func (OllamaConfig) GetClientType() AIClientType {
	return AIClientTypeClaude
}
func (OllamaConfig) GetCategory() ClientCategory {
	return ClientCategoryMedia
}

func (c *OllamaConfig) UnmarshalJSON(data []byte) error {
	// Create a temporary type to avoid recursion
	type Alias OllamaConfig
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Ensure Type is always the correct constant
	c.ClientType = AIClientTypeOllama
	c.Type = ClientTypeOllama
	return nil
}
