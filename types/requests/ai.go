package requests

// GenerateTextRequest defines the request body for generating text with an AI client
// @Description Request for generating text with an AI client
type GenerateTextRequest struct {
	// The prompt text to generate from
	// required: true
	// example: Write a short poem about programming
	Prompt string `json:"prompt" binding:"required"`

	// Optional system instructions to guide the AI
	// example: You are a helpful assistant specializing in creative writing
	SystemInstructions string `json:"systemInstructions,omitempty"`

	// Maximum number of tokens to generate
	// example: 500
	MaxTokens uint `json:"maxTokens,omitempty"`

	// Temperature for AI response (0.0 to 1.0)
	// example: 0.7
	Temperature float32 `json:"temperature,omitempty"`
}

// GenerateStructuredRequest defines the request body for generating structured data with an AI client
// @Description Request for generating structured JSON data with an AI client
type GenerateStructuredRequest struct {
	// The prompt text to generate structured data from
	// required: true
	// example: Create a JSON object representing three programming languages
	Prompt string `json:"prompt" binding:"required"`

	// Optional system instructions to guide the AI
	// example: You are a helpful assistant specializing in structured data generation
	SystemInstructions string `json:"systemInstructions,omitempty"`

	// Maximum number of tokens to generate
	// example: 500
	MaxTokens uint `json:"maxTokens,omitempty"`

	// Temperature for AI response (0.0 to 1.0)
	// example: 0.2
	Temperature float32 `json:"temperature,omitempty"`
}

// AiRecommendationRequest defines the request for AI-powered recommendations
// @Description Request for AI-powered content recommendations
type AiRecommendationRequest struct {
	// Type of content to recommend (movie, tv, music, etc)
	// required: true
	// example: movie
	ContentType string `json:"contentType" binding:"required"`

	// Number of recommendations to return
	// example: 5
	Count int `json:"count,omitempty"`

	// Optional filters to apply to recommendations
	// example: {"genre": "sci-fi", "year": "2020-2023"}
	Filters map[string]any `json:"filters,omitempty"`

	// Specific AI client type to use (claude, openai, ollama)
	// example: claude
	ClientType string `json:"clientType,omitempty"`

	// Client ID to use for the conversation
	ClientID uint64 `json:"clientId,omitempty"`
}

// AiContentAnalysisRequest defines the request for AI content analysis
// @Description Request for AI-powered content analysis
type AiContentAnalysisRequest struct {
	// Type of content being analyzed (text, movie, etc)
	// required: true
	// example: text
	ContentType string `json:"contentType" binding:"required"`

	// The content to analyze
	// required: true
	// example: This is a sample text that needs analysis for sentiment and themes.
	Content string `json:"content" binding:"required"`

	// Optional analysis options
	// example: {"includeThemes": true, "includeSentiment": true}
	Options map[string]any `json:"options,omitempty"`

	// Specific AI client type to use (claude, openai, ollama)
	// example: claude
	ClientType string `json:"clientType,omitempty"`

	// Client ID to use for the conversation
	ClientID uint64 `json:"clientId,omitempty"`
}

// StartConversationRequest defines the request to start a new AI conversation
// @Description Request to start a new AI-powered conversation for recommendations
type StartConversationRequest struct {
	// Type of content to discuss (movie, tv, music, etc)
	// required: true
	// example: movie
	ContentType string `json:"contentType" binding:"required"`

	// Optional user preferences to initialize the conversation
	// example: {"favoriteGenres": ["sci-fi", "thriller"], "recentlyWatched": ["Inception", "Tenet"]}
	Preferences map[string]any `json:"preferences,omitempty"`

	// Optional custom system instructions
	// example: You are a helpful movie recommendation assistant
	SystemInstructions string `json:"systemInstructions,omitempty"`

	// Client ID to use for the conversation
	ClientID uint64 `json:"clientId,omitempty"`
}

// ConversationMessageRequest defines a message in an existing conversation
// @Description Request to send a message in an existing AI conversation
type ConversationMessageRequest struct {
	// The conversation ID from a previous StartConversation call
	// required: true
	// example: conv-123-abcdef
	ConversationID string `json:"conversationId" binding:"required"`

	// The user's message to the AI
	// required: true
	// example: I'm looking for sci-fi movies similar to Interstellar
	Message string `json:"message" binding:"required"`

	// Optional context information for this message
	// example: {"includeRecommendations": true, "maxResults": 3}
	Context map[string]any `json:"context,omitempty"`

	// Client ID to use for the conversation
	ClientID uint64 `json:"clientId,omitempty"`
}

