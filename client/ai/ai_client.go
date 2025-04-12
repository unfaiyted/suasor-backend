// client/ai/ai_client.go
package ai

import (
	"context"
	"errors"
	"suasor/client"
	aitypes "suasor/client/ai/types"
	types "suasor/client/types"
)

var ErrFeatureNotSupported = errors.New("feature not supported by this AI client")

// AIClient defines the interface for all AI providers
type AIClient interface {
	client.Client
	
	// Core text generation capabilities
	GenerateText(ctx context.Context, promptText string, options *aitypes.GenerationOptions) (string, error)
	GenerateStructured(ctx context.Context, promptText string, outputSchema interface{}, options *aitypes.GenerationOptions) error
	GenerateContent(ctx context.Context, systemPrompt string, userPrompt string, model string, options map[string]interface{}) (*aitypes.ContentResponse, error)
	
	// Conversational capabilities
	StartConversation(ctx context.Context, systemInstructions string) (string, error)
	SendMessage(ctx context.Context, conversationID string, message string) (string, error)
	CreateMessage(ctx context.Context, request aitypes.MessageRequest) (*aitypes.MessageResponse, error)
	
	// Recommendations capabilities
	GetRecommendations(ctx context.Context, request *aitypes.RecommendationRequest) (*aitypes.RecommendationResponse, error)
	
	// Information methods
	GetSupportedModels() []string
	GetCapabilities() *aitypes.AICapabilities
}

// BaseAIClient is a partial implementation of the AIClient interface
type BaseAIClient struct {
	client.BaseClient
	ClientType types.AIClientType
	config     *types.AIClientConfig
}

func NewAIClient(ctx context.Context, clientID uint64, clientType types.AIClientType, config types.AIClientConfig) (AIClient, error) {
	return &BaseAIClient{
		BaseClient: client.BaseClient{
			ClientID: clientID,
			Category: clientType.AsCategory(),
		},
		config:     &config,
		ClientType: clientType,
	}, nil
}

func (b *BaseAIClient) TestConnection(ctx context.Context) (bool, error) {
	return false, nil
}

// Provide default implementations for AIClient interface methods
// These methods will be overridden by concrete implementations

func (b *BaseAIClient) GenerateText(ctx context.Context, promptText string, options *aitypes.GenerationOptions) (string, error) {
	return "", ErrFeatureNotSupported
}

func (b *BaseAIClient) GenerateStructured(ctx context.Context, promptText string, outputSchema interface{}, options *aitypes.GenerationOptions) error {
	return ErrFeatureNotSupported
}

func (b *BaseAIClient) StartConversation(ctx context.Context, systemInstructions string) (string, error) {
	return "", ErrFeatureNotSupported
}

func (b *BaseAIClient) SendMessage(ctx context.Context, conversationID string, message string) (string, error) {
	return "", ErrFeatureNotSupported
}

func (b *BaseAIClient) GetSupportedModels() []string {
	return []string{}
}

func (b *BaseAIClient) GetCapabilities() *aitypes.AICapabilities {
	return &aitypes.AICapabilities{
		SupportsStructuredOutput: false,
		SupportsConversation:     false,
		SupportsStreaming:        false,
		MaxContextTokens:         0,
		DefaultMaxTokens:         0,
	}
}

func (b *BaseAIClient) GetRecommendations(ctx context.Context, request *aitypes.RecommendationRequest) (*aitypes.RecommendationResponse, error) {
	return nil, ErrFeatureNotSupported
}

func (b *BaseAIClient) GenerateContent(ctx context.Context, systemPrompt string, userPrompt string, model string, options map[string]interface{}) (*aitypes.ContentResponse, error) {
	return nil, ErrFeatureNotSupported
}

func (b *BaseAIClient) CreateMessage(ctx context.Context, request aitypes.MessageRequest) (*aitypes.MessageResponse, error) {
	return nil, ErrFeatureNotSupported
}
