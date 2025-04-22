# Claude AI Client for Suasor

This package provides integration with Claude AI models from Anthropic, enabling natural language generation capabilities in the Suasor application.

## Features

- Text generation with Claude models
- Structured output (JSON) generation
- Conversational interactions with memory
- Support for multiple Claude models

## Configuration

The Claude client requires the following configuration:

```go
config := types.ClaudeConfig{
    BaseAIClientConfig: types.BaseAIClientConfig{
        BaseClientConfig: types.BaseClientConfig{
            Type:    types.ClientTypeClaude,
            Name:    "Your Claude Client Name",
            Enabled: true,
        },
        ClientType:       types.AIClientTypeClaude,
        APIKey:           "your-api-key",               // Required
        Model:            "claude-3-haiku-20240307",    // Required
        Temperature:      0.7,                          // Optional, defaults to model default
        MaxTokens:        150,                          // Optional, defaults to model default
        MaxContextTokens: 2000,                         // Optional, defaults to model default
    },
}
```

## Usage

### Creating a Client

```go
import (
    "suasor/client/ai/claude"
    "suasor/client/types"
)

// Create a Claude client
ctx := context.Background()
client, err := claude.NewClaudeClient(ctx, clientID, config)
if err != nil {
    log.Fatalf("Failed to create Claude client: %v", err)
}
```

### Basic Text Generation

```go
response, err := client.GenerateText(ctx, "Tell me a joke about programming", nil)
if err != nil {
    log.Fatalf("Failed to generate text: %v", err)
}
fmt.Printf("Claude says: %s\n", response)
```

### Structured Output Generation

```go
type MovieRecommendation struct {
    Title       string   `json:"title"`
    Description string   `json:"description"`
    Year        int      `json:"year"`
    Genres      []string `json:"genres"`
    Rating      float64  `json:"rating"`
}

var movie MovieRecommendation
err := client.GenerateStructured(ctx, 
    "Suggest a sci-fi movie from the 1980s with a brief description", 
    &movie, 
    nil)
```

### Conversational Interactions

```go
// Start a conversation with system instructions
conversationID, err := client.StartConversation(ctx, "You are a friendly assistant.")
if err != nil {
    log.Fatalf("Failed to start conversation: %v", err)
}

// Send a message in the conversation
response1, err := client.SendMessage(ctx, conversationID, "Hi, who are you?")

// Send a follow-up message
response2, err := client.SendMessage(ctx, conversationID, "What can you help me with?")
```

## Available Models

- `claude-3-5-sonnet-20240620`
- `claude-3-opus-20240229`
- `claude-3-sonnet-20240229`
- `claude-3-haiku-20240307`
- `claude-2.1`

## Testing

Set the `CLAUDE_API_KEY` environment variable before running tests:

```bash
export CLAUDE_API_KEY=your-api-key
go test ./client/ai/claude -v
```

## Dependencies

This client uses the [gollm](https://github.com/teilomillet/gollm) library to communicate with the Claude API.