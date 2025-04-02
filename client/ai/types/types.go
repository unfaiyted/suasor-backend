package types

// GenerationOptions contains options for AI text generation
type GenerationOptions struct {
	Temperature       float64 // Controls randomness (0.0-1.0)
	MaxTokens         int     // Maximum response length
	SystemInstructions string  // System instructions for the AI
	ResponseFormat    string  // Optional format like "json" if supported
}

// AICapabilities describes what features an AI model supports
type AICapabilities struct {
	SupportsStructuredOutput bool  // Whether the model can output structured data like JSON
	SupportsConversation     bool  // Whether the model supports conversational mode
	SupportsStreaming        bool  // Whether the model supports streaming responses
	MaxContextTokens         int   // Maximum context window size
	DefaultMaxTokens         int   // Default maximum tokens for responses
	AvailableModels          []string // List of available models
}

// AIResponse represents a structured response from an AI client
type AIResponse struct {
	Content    string            // The primary content returned
	Metadata   map[string]string // Any additional metadata
	TokenUsage TokenUsage        // Token usage information if available
}

// TokenUsage tracks token consumption 
type TokenUsage struct {
	PromptTokens     int // Tokens used in the prompt
	CompletionTokens int // Tokens used in the completion
	TotalTokens      int // Total tokens used
}