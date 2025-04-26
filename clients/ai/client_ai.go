// client/ai/ai_clients.go
package ai

import (
	"context"
	"errors"
	"suasor/clients"
	aitypes "suasor/clients/ai/types"
	types "suasor/clients/types"
)

var ErrFeatureNotSupported = errors.New("feature not supported by this AI client")

// AIClient defines the interface for all AI providers
type ClientAI interface {
	clients.Client

	// Core text generation capabilities
	GenerateText(ctx context.Context, promptText string, options *aitypes.GenerationOptions) (string, error)
	GenerateStructured(ctx context.Context, promptText string, outputSchema any, options *aitypes.GenerationOptions) error
	GenerateContent(ctx context.Context, systemPrompt string, userPrompt string, model string, options map[string]any) (*aitypes.ContentResponse, error)

	// Conversational capabilities
	StartConversation(ctx context.Context, systemInstructions string) (string, error)
	SendMessage(ctx context.Context, conversationID string, message string) (string, error)
	CreateMessage(ctx context.Context, request aitypes.MessageRequest) (*aitypes.MessageResponse, error)

	// Recommendations capabilities
	GetRecommendations(ctx context.Context, request *aitypes.RecommendationRequest) (*aitypes.RecommendationResponse, error)

	AnalyzeContent(ctx context.Context, contentType string, content string, options map[string]any) (map[string]any, error)

	StartRecommendationConversation(ctx context.Context, contentType string, preferences map[string]any, systemInstructions string) (string, string, error)
	ContinueRecommendationConversation(ctx context.Context, conversationID string, message string, context map[string]any) (string, []map[string]any, error)

	// Information methods
	GetSupportedModels() []string
	GetCapabilities() *aitypes.AICapabilities
}

// clientAI is a partial implementation of the AIClient interface
type clientAI struct {
	clients.Client
	ClientType types.AIClientType
	config     *types.AIClientConfig
}

func NewAIClient(ctx context.Context, clientID uint64, clientType types.AIClientType, config types.AIClientConfig) (ClientAI, error) {

	client := clients.NewClient(clientID, clientType.AsCategory(), config)
	return &clientAI{
		Client:     client,
		config:     &config,
		ClientType: clientType,
	}, nil
}

func (b *clientAI) TestConnection(ctx context.Context) (bool, error) {
	return false, nil
}

// Provide default implementations for AIClient interface methods
// These methods will be overridden by concrete implementations

func (b *clientAI) GenerateText(ctx context.Context, promptText string, options *aitypes.GenerationOptions) (string, error) {
	return "", ErrFeatureNotSupported
}

func (b *clientAI) GenerateStructured(ctx context.Context, promptText string, outputSchema any, options *aitypes.GenerationOptions) error {
	return ErrFeatureNotSupported
}

func (b *clientAI) StartConversation(ctx context.Context, systemInstructions string) (string, error) {
	return "", ErrFeatureNotSupported
}

func (b *clientAI) SendMessage(ctx context.Context, conversationID string, message string) (string, error) {
	return "", ErrFeatureNotSupported
}

func (b *clientAI) GetSupportedModels() []string {
	return []string{}
}

func (b *clientAI) GetCapabilities() *aitypes.AICapabilities {
	return &aitypes.AICapabilities{
		SupportsStructuredOutput: false,
		SupportsConversation:     false,
		SupportsStreaming:        false,
		MaxContextTokens:         0,
		DefaultMaxTokens:         0,
	}
}

func (b *clientAI) GetRecommendations(ctx context.Context, request *aitypes.RecommendationRequest) (*aitypes.RecommendationResponse, error) {
	return nil, ErrFeatureNotSupported
}

func (b *clientAI) GenerateContent(ctx context.Context, systemPrompt string, userPrompt string, model string, options map[string]any) (*aitypes.ContentResponse, error) {
	return nil, ErrFeatureNotSupported
}

func (b *clientAI) CreateMessage(ctx context.Context, request aitypes.MessageRequest) (*aitypes.MessageResponse, error) {
	return nil, ErrFeatureNotSupported
}

func (b *clientAI) AnalyzeContent(ctx context.Context, contentType string, content string, options map[string]any) (map[string]any, error) {
	return nil, ErrFeatureNotSupported
}

func (b *clientAI) StartRecommendationConversation(ctx context.Context, contentType string, preferences map[string]any, systemInstructions string) (string, string, error) {
	return "", "", ErrFeatureNotSupported
}
func (b *clientAI) ContinueRecommendationConversation(ctx context.Context, conversationID string, message string, context map[string]any) (string, []map[string]any, error) {
	return "", nil, ErrFeatureNotSupported
}
