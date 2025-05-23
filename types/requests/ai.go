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
	ClientID uint64 `json:"clientID,omitempty"`
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
	ClientID uint64 `json:"clientID,omitempty"`
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
	ClientID uint64 `json:"clientID,omitempty"`
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
	ClientID uint64 `json:"clientID,omitempty"`
}

// ContinueConversationRequest defines a request to continue a past conversation
// @Description Request to continue a previous AI conversation from history
type ContinueConversationRequest struct {
	// The user's message to continue the conversation with
	// required: true
	// example: What other movies would you recommend based on our previous discussion?
	Message string `json:"message" binding:"required"`

	// Optional context information for this message
	// example: {"extractRecommendations": true}
	Context map[string]any `json:"context,omitempty"`
}

// GetUserConversationsRequest defines filters for retrieving conversation history
// @Description Filters for retrieving AI conversation history
type GetUserConversationsRequest struct {
	// Number of conversations to return (default: 20)
	// example: 10
	Limit int `form:"limit,default=20"`

	// Offset for pagination (default: 0)
	// example: 20
	Offset int `form:"offset,default=0"`

	// Filter by conversation status (active, archived, all)
	// example: active
	Status string `form:"status"`

	// Field to sort by (default: updatedAt)
	// example: createdAt
	SortBy string `form:"sortBy,default=updatedAt"`

	// Sort direction (asc or desc) (default: desc)
	// example: desc
	SortDir string `form:"sortDir,default=desc"`
}

// GetUserRecommendationsRequest defines filters for retrieving recommendation history
// @Description Filters for retrieving AI recommendation history
type GetUserRecommendationsRequest struct {
	// Number of recommendations to return (default: 20)
	// example: 10
	Limit int `form:"limit,default=20"`

	// Offset for pagination (default: 0)
	// example: 20
	Offset int `form:"offset,default=0"`

	// Filter by item type (movie, music, etc)
	// example: movie
	ItemType string `form:"itemType"`

	// Filter by selection status (true for selected, false for not selected)
	// example: true
	Selected *bool `form:"selected"`

	// Field to sort by (default: createdAt)
	// example: title
	SortBy string `form:"sortBy,default=createdAt"`

	// Sort direction (asc or desc) (default: desc)
	// example: desc
	SortDir string `form:"sortDir,default=desc"`
}
