package types

// GenerationOptions contains options for AI text generation
type GenerationOptions struct {
	Temperature        float64 // Controls randomness (0.0-1.0)
	MaxTokens          int     // Maximum response length
	SystemInstructions string  // System instructions for the AI
	ResponseFormat     string  // Optional format like "json" if supported
}

// AICapabilities describes what features an AI model supports
type AICapabilities struct {
	SupportsStructuredOutput bool     // Whether the model can output structured data like JSON
	SupportsConversation     bool     // Whether the model supports conversational mode
	SupportsStreaming        bool     // Whether the model supports streaming responses
	MaxContextTokens         int      // Maximum context window size
	DefaultMaxTokens         int      // Default maximum tokens for responses
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

// RecommendationRequest contains parameters for generating recommendations
type RecommendationRequest struct {
	MediaType         string             // Type of media ("movie", "series", "music")
	UserPreferences   map[string]any     // User preferences to consider
	ExcludeIDs        []string           // IDs to exclude from recommendations
	Count             int                // Number of recommendations to generate
	IncludeSimilarTo  []string           // Include items similar to these IDs
	AdditionalContext string             // Additional context or instructions
	GenerationOptions *GenerationOptions // Options for the generation process
}

// RecommendationItem represents a single recommended item
type RecommendationItem struct {
	Title       string   `json:"title"`
	Year        int      `json:"year,omitempty"`
	Genres      []string `json:"genres,omitempty"`
	Reason      string   `json:"reason,omitempty"`
	ExternalID  string   `json:"externalId,omitempty"`
	Rating      float32  `json:"rating,omitempty"`
	Popularity  int      `json:"popularity,omitempty"`
	SourceNames []string `json:"sourceNames,omitempty"`
	PosterURL   string   `json:"posterUrl,omitempty"`
	BackdropURL string   `json:"backdropUrl,omitempty"`
	ReleaseDate string   `json:"releaseDate,omitempty"`
	Directors   []string `json:"directors,omitempty"`
	Actors      []string `json:"actors,omitempty"`
	Description string   `json:"description,omitempty"`
}

// RecommendationResponse contains the list of recommendations
type RecommendationResponse struct {
	Items       []RecommendationItem `json:"items"`
	Explanation string               `json:"explanation,omitempty"`
	TokenUsage  TokenUsage           `json:"tokenUsage,omitempty"`
}

// ContentResponse represents a structured response from the AI model
type ContentResponse struct {
	Text       string     `json:"text"`
	TokenUsage TokenUsage `json:"tokenUsage,omitempty"`
}

// Message represents a message in a conversation
type Message struct {
	Role    string           `json:"role"`
	Content []MessageContent `json:"content"`
}

// MessageContent represents the content of a message
type MessageContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ResponseFormat describes the format of the response
type ResponseFormat struct {
	Type string `json:"type"`
}

// MessageRequest represents a request to create a message
type MessageRequest struct {
	Model          string         `json:"model"`
	Messages       []Message      `json:"messages"`
	MaxTokens      int            `json:"max_tokens,omitempty"`
	Temperature    float64        `json:"temperature,omitempty"`
	TopP           float64        `json:"top_p,omitempty"`
	ResponseFormat ResponseFormat `json:"response_format,omitempty"`
	ConversationID string         `json:"conversation_id,omitempty"`
	System         string         `json:"system,omitempty"`
	Tools          any            `json:"tools,omitempty"`
	ToolChoice     any            `json:"tool_choice,omitempty"`
	Stop           []string       `json:"stop,omitempty"`
	Stream         bool           `json:"stream,omitempty"`
	TopK           int            `json:"top_k,omitempty"`
	Seed           int            `json:"seed,omitempty"`
}

// MessageResponse represents a response to a message request
type MessageResponse struct {
	ID         string           `json:"id"`
	Object     string           `json:"object"`
	Created    int64            `json:"created"`
	Model      string           `json:"model"`
	Content    []MessageContent `json:"content"`
	Role       string           `json:"role"`
	TokenUsage TokenUsage       `json:"usage,omitempty"`
}

