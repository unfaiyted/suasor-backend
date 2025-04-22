package types

import (
	"context"
)

type AIClientConfig interface {
	ClientConfig
	GetClientType() AIClientType
	GetAPIKey() string
}

// AiClient interface for interacting with AI services
type AiClient interface {
	// Basic recommendation and analysis
	GetRecommendations(ctx context.Context, contentType string, filters map[string]any, count int) ([]map[string]interface{}, error)
	AnalyzeContent(ctx context.Context, contentType string, content string, options map[string]any) (map[string]interface{}, error)

	// Conversational recommendation
	StartRecommendationConversation(ctx context.Context, contentType string, preferences map[string]interface{}, systemInstructions string) (string, string, error)
	ContinueRecommendationConversation(ctx context.Context, conversationID string, message string, context map[string]interface{}) (string, []map[string]interface{}, error)
}

type BaseAIClientConfig struct {
	BaseClientConfig
	ClientType       AIClientType `json:"clientType"`
	BaseURL          string       `json:"baseURL" mapstructure:"baseURL" example:"http://localhost:8096"`
	APIKey           string       `json:"apiKey" mapstructure:"apiKey" example:"your-api-key" binding:"required_if=Enabled true"`
	Model            string       `json:"model" mapstructure:"model" example:"claude-2"`
	Temperature      float64      `json:"temperature" mapstructure:"temperature" example:"0.5"`
	MaxTokens        int          `json:"maxTokens" mapstructure:"maxTokens" example:"100"`
	MaxContextTokens int          `json:"maxContextTokens" mapstructure:"maxContextTokens" example:"1000"`
}

func (c *BaseAIClientConfig) GetClientType() AIClientType {
	return c.ClientType
}

func (c *BaseAIClientConfig) GetCategory() ClientCategory {
	c.Category = ClientCategoryAI
	return c.Category
}

func (c *BaseAIClientConfig) GetAPIKey() string {
	return c.APIKey
}
