# AI Clients for Suasor

This package provides a collection of AI integrations for the Suasor application, enabling natural language processing and generation capabilities.

## Available AI Clients

- **Claude**: Integration with Anthropic's Claude models
- **OpenAI**: (In development) Integration with OpenAI models

## Architecture

The AI client architecture follows a common interface pattern, allowing different AI providers to be used interchangeably:

```
AIClient (interface)
 ├── BaseAIClient (partial implementation)
 │    ├── ClaudeClient (concrete implementation)
 │    └── OpenAIClient (concrete implementation)
 └── ...
```

### Core Interface

All AI clients implement the `AIClient` interface defined in `ai_client.go`:

```go
type AIClient interface {
    client.Client
    
    // Core text generation capabilities
    GenerateText(ctx context.Context, promptText string, options *aitypes.GenerationOptions) (string, error)
    GenerateStructured(ctx context.Context, promptText string, outputSchema interface{}, options *aitypes.GenerationOptions) error
    
    // Conversational capabilities
    StartConversation(ctx context.Context, systemInstructions string) (string, error)
    SendMessage(ctx context.Context, conversationID string, message string) (string, error)
    
    // Information methods
    GetSupportedModels() []string
    GetCapabilities() *aitypes.AICapabilities
}
```

## Common Usage Pattern

All AI clients follow a similar usage pattern:

1. Create configuration with appropriate client-specific settings
2. Instantiate client with configuration and client ID
3. Use client methods for text generation, structured output, or conversation

Example:

```go
// Create configuration
config := types.ClaudeConfig{...}

// Create client
client, err := claude.NewClaudeClient(ctx, clientID, config)

// Use client
response, err := client.GenerateText(ctx, "Write a haiku about programming", nil)
```

## Adding New AI Clients

To add a new AI client:

1. Create a new subdirectory for the provider
2. Implement the `AIClient` interface
3. Create appropriate config types in `client/types`
4. Add tests and documentation

## Running Tests

To run tests for all AI clients:

```bash
go test ./client/ai/... -v
```

Tests will be skipped if the required API keys are not provided as environment variables.