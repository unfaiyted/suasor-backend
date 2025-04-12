package claude

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/teilomillet/gollm"
	"suasor/client"
	"suasor/client/ai"
	aitypes "suasor/client/ai/types"
	"suasor/client/types"
	"suasor/utils"
)

// Add init function to register the Claude client factory
func init() {
	fmt.Printf("Registering factory for client type: %s (value: %v)\n",
		types.ClientTypeClaude.String(), types.ClientTypeClaude)

	fmt.Println("Registering Claude client factory...")
	client.RegisterClientFactory(types.ClientTypeClaude,
		func(ctx context.Context, clientID uint64, configData types.ClientConfig) (client.Client, error) {
			// Use the provided config (should be a ClaudeConfig)
			claudeConfig, ok := configData.(*types.ClaudeConfig)
			if !ok {
				return nil, fmt.Errorf("expected *types.ClaudeConfig, got %T", configData)
			}
			return NewClaudeClient(ctx, clientID, *claudeConfig)
		})
}

// ClaudeClient implements the AI client interface
type ClaudeClient struct {
	client.BaseClient
	llm           gollm.LLM
	config        types.ClaudeConfig
	memoryID      string // For conversation tracking
	conversations map[string]ConversationContext
}

// ConversationContext tracks the state of a conversation
type ConversationContext struct {
	ContentType        string
	UserPreferences    map[string]interface{}
	SystemInstructions string
	History            []ChatMessage
}

// ChatMessage represents a message in the conversation
type ChatMessage struct {
	Role    string // "user" or "assistant"
	Content string
}

// NewClaudeClient creates a new Claude client instance using gollm
func NewClaudeClient(ctx context.Context, clientID uint64, cfg types.ClaudeConfig) (ai.AIClient, error) {
	log := utils.LoggerFromContext(ctx)
	log.Info().
		Uint64("clientID", clientID).
		Str("model", cfg.Model).
		Msg("Creating new Claude client")

	// Create gollm LLM instance with Claude configuration
	llm, err := gollm.NewLLM(
		gollm.SetProvider("anthropic"),
		gollm.SetModel(cfg.Model),
		gollm.SetAPIKey(cfg.APIKey),
		gollm.SetMaxTokens(cfg.MaxTokens),
		gollm.SetTemperature(cfg.Temperature),
		gollm.SetMaxRetries(3),
		gollm.SetMemory(4096), // Enable memory for conversational context
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Claude LLM: %w", err)
	}

	return &ClaudeClient{
		BaseClient: client.BaseClient{
			ClientID: clientID,
			Category: types.ClientCategoryAI,
			Type:     types.ClientTypeClaude,
			Config:   &cfg,
		},
		llm:           llm,
		config:        cfg,
		conversations: make(map[string]ConversationContext),
	}, nil
}

// TestConnection tests the connection to Claude API
func (c *ClaudeClient) TestConnection(ctx context.Context) (bool, error) {
	log := utils.LoggerFromContext(ctx)
	log.Info().Msg("Testing Claude connection")

	// Simple prompt to test connectivity
	prompt := gollm.NewPrompt("Hello, are you available?")

	_, err := c.llm.Generate(ctx, prompt)
	if err != nil {
		log.Error().Err(err).Msg("Claude connection test failed")
		return false, err
	}

	log.Info().Msg("Claude connection test successful")
	return true, nil
}

// GenerateText sends a prompt to Claude and returns the response
func (c *ClaudeClient) GenerateText(ctx context.Context, promptText string, options *aitypes.GenerationOptions) (string, error) {
	log := utils.LoggerFromContext(ctx)

	// Apply custom options if provided
	if options != nil {
		if options.MaxTokens > 0 {
			c.llm.SetOption("max_tokens", options.MaxTokens)
		}
		if options.Temperature > 0 {
			c.llm.SetOption("temperature", options.Temperature)
		}
	}

	// Create the prompt
	prompt := gollm.NewPrompt(promptText)

	// Add system instructions if provided
	if options != nil && options.SystemInstructions != "" {
		prompt = gollm.NewPrompt(promptText,
			gollm.WithContext(options.SystemInstructions))
	}

	// Generate the response
	response, err := c.llm.Generate(ctx, prompt)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate text with Claude")
		return "", fmt.Errorf("claude text generation failed: %w", err)
	}

	return response, nil
}

// GenerateStructured generates structured output (like JSON) from Claude
func (c *ClaudeClient) GenerateStructured(ctx context.Context, promptText string, outputSchema interface{}, options *aitypes.GenerationOptions) error {
	log := utils.LoggerFromContext(ctx)

	// Add instructions for structured output
	structuredPrompt := promptText + "\n\nRespond ONLY with valid JSON matching the required schema, with no additional text before or after."

	// Set specific options for JSON generation
	jsonOptions := &aitypes.GenerationOptions{
		Temperature: 0.2, // Lower temperature for more predictable output
		MaxTokens:   500,
	}
	if options != nil {
		if options.MaxTokens > 0 {
			jsonOptions.MaxTokens = options.MaxTokens
		}
		if options.Temperature > 0 {
			jsonOptions.Temperature = options.Temperature
		}
		jsonOptions.SystemInstructions = options.SystemInstructions
	}

	// Add system instructions specifically for JSON if not provided
	if jsonOptions.SystemInstructions == "" {
		jsonOptions.SystemInstructions = "You are a helpful assistant that responds only with valid JSON. Do not include any explanations, markdown formatting, or text outside of the JSON structure."
	}

	// Generate the JSON response
	response, err := c.GenerateText(ctx, structuredPrompt, jsonOptions)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate structured output with Claude")
		return fmt.Errorf("claude structured generation failed: %w", err)
	}

	// Clean the response to ensure it only contains JSON
	response = cleanJSONResponse(response)

	// Parse the JSON response into the provided schema
	if err := json.Unmarshal([]byte(response), outputSchema); err != nil {
		log.Error().Err(err).Str("response", response).Msg("Failed to parse Claude response as JSON")
		return fmt.Errorf("failed to parse Claude JSON response: %w", err)
	}

	return nil
}

// cleanJSONResponse removes any text before the first { or [ and after the last } or ]
func cleanJSONResponse(input string) string {
	startIdx := 0
	for i, char := range input {
		if char == '{' || char == '[' {
			startIdx = i
			break
		}
	}

	endIdx := len(input)
	for i := len(input) - 1; i >= 0; i-- {
		if input[i] == '}' || input[i] == ']' {
			endIdx = i + 1
			break
		}
	}

	if startIdx < endIdx {
		return input[startIdx:endIdx]
	}
	return input
}

// StartConversation begins a new conversation with Claude
func (c *ClaudeClient) StartConversation(ctx context.Context, systemInstructions string) (string, error) {
	// Generate a unique, URL-safe conversation ID
	c.memoryID = fmt.Sprintf("conv-%d-%s", c.ClientID, utils.GenerateRandomID(12))

	// Set system instructions if provided
	if systemInstructions != "" {
		c.llm.SetOption("system", systemInstructions)
	}

	return c.memoryID, nil
}

// SendMessage sends a message in an existing conversation
func (c *ClaudeClient) SendMessage(ctx context.Context, conversationID string, message string) (string, error) {
	log := utils.LoggerFromContext(ctx)

	// Store conversation ID for future reference
	if c.memoryID != conversationID {
		c.memoryID = conversationID
	}

	// In this implementation we don't use a stateful conversation API
	// Instead we rely on the model's ability to follow the conversation
	prompt := gollm.NewPrompt(message)

	response, err := c.llm.Generate(ctx, prompt)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send message to Claude")
		return "", fmt.Errorf("claude message send failed: %w", err)
	}

	return response, nil
}

// GetSupportedModels returns a list of Claude models supported by this client
func (c *ClaudeClient) GetSupportedModels() []string {
	return []string{
		"claude-3-5-sonnet-20240620",
		"claude-3-opus-20240229",
		"claude-3-sonnet-20240229",
		"claude-3-haiku-20240307",
		"claude-2.1",
	}
}

// GetCapabilities returns the capabilities of this Claude client
func (c *ClaudeClient) GetCapabilities() *aitypes.AICapabilities {
	return &aitypes.AICapabilities{
		SupportsStructuredOutput: true,
		SupportsConversation:     true,
		SupportsStreaming:        false, // Implement if needed
		MaxContextTokens:         c.config.MaxContextTokens,
		DefaultMaxTokens:         c.config.MaxTokens,
	}
}

// GetRecommendations implements the AIClient interface to get content recommendations
func (c *ClaudeClient) GetRecommendations(ctx context.Context, request *aitypes.RecommendationRequest) (*aitypes.RecommendationResponse, error) {
	log := utils.LoggerFromContext(ctx)
	log.Info().
		Str("mediaType", request.MediaType).
		Int("count", request.Count).
		Msg("Getting recommendations from Claude")

	// Build the prompt based on content type and filters
	prompt := fmt.Sprintf("Please recommend %d %s items", request.Count, request.MediaType)

	// Add filter information
	if len(request.UserPreferences) > 0 {
		filterInfo := "\nConsider these preferences:\n"
		for k, v := range request.UserPreferences {
			filterInfo += fmt.Sprintf("- %s: %v\n", k, v)
		}
		prompt += filterInfo
	}

	// Add excluded IDs if any
	if len(request.ExcludeIDs) > 0 {
		prompt += "\nPlease exclude the following items:\n"
		for _, id := range request.ExcludeIDs {
			prompt += fmt.Sprintf("- %s\n", id)
		}
	}

	// Add additional context if provided
	if request.AdditionalContext != "" {
		prompt += "\n" + request.AdditionalContext
	}

	// Add output format instructions
	prompt += "\nPlease return the recommendations as a JSON array of objects with the following structure for each recommendation:" +
		"\n- title: The title of the item" +
		"\n- year: The release year (as a number)" +
		"\n- genres: An array of genre strings" +
		"\n- reason: A brief explanation of why this is recommended" +
		"\n- rating: The average rating (0-10 scale)" +
		"\n- popularity: A popularity rank (1-100, where 1 is most popular)" +
		"\n- description: A brief description of the item"

	// Generate the structured response
	var recommendations []aitypes.RecommendationItem

	// Use the provided generation options or create default ones
	genOptions := &aitypes.GenerationOptions{
		Temperature:        0.4,
		MaxTokens:          2000,
		SystemInstructions: fmt.Sprintf("You are a helpful recommendation system specialized in %s. Provide detailed and personalized recommendations based on the user's preferences.", request.MediaType),
	}

	if request.GenerationOptions != nil {
		if request.GenerationOptions.Temperature > 0 {
			genOptions.Temperature = request.GenerationOptions.Temperature
		}
		if request.GenerationOptions.MaxTokens > 0 {
			genOptions.MaxTokens = request.GenerationOptions.MaxTokens
		}
		if request.GenerationOptions.SystemInstructions != "" {
			genOptions.SystemInstructions = request.GenerationOptions.SystemInstructions
		}
	}

	// Generate the structured response
	err := c.GenerateStructured(ctx, prompt, &recommendations, genOptions)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get recommendations from Claude")
		return nil, err
	}

	// Ensure we have the requested number of recommendations if possible
	if len(recommendations) > request.Count {
		recommendations = recommendations[:request.Count]
	}

	// Create and return the response
	response := &aitypes.RecommendationResponse{
		Items:       recommendations,
		Explanation: "Recommendations generated based on user preferences using Claude",
	}

	return response, nil
}

// AnalyzeContent implements the AiClient interface to analyze content
func (c *ClaudeClient) AnalyzeContent(ctx context.Context, contentType string, content string, options map[string]interface{}) (map[string]interface{}, error) {
	log := utils.LoggerFromContext(ctx)
	log.Info().
		Str("contentType", contentType).
		Int("contentLength", len(content)).
		Interface("options", options).
		Msg("Analyzing content with Claude")

	// Build the prompt based on content type and options
	prompt := fmt.Sprintf("Please analyze the following %s content:\n\n%s\n\n", contentType, content)

	// Add specific analysis instructions based on options
	if options != nil {
		if include, ok := options["includeThemes"].(bool); ok && include {
			prompt += "Include main themes and motifs. "
		}
		if include, ok := options["includeSentiment"].(bool); ok && include {
			prompt += "Analyze the sentiment. "
		}
		if include, ok := options["includeStyleAnalysis"].(bool); ok && include {
			prompt += "Analyze the stylistic elements. "
		}
	}

	// Add output format instructions
	prompt += "\nPlease return the analysis as a JSON object with appropriate fields for the requested analysis."

	// Create the output schema to receive the data
	var analysis map[string]interface{}

	// Generate the structured response
	err := c.GenerateStructured(ctx, prompt, &analysis, &aitypes.GenerationOptions{
		Temperature:        0.3,
		MaxTokens:          2000,
		SystemInstructions: fmt.Sprintf("You are an expert at analyzing %s content. Provide detailed, insightful analysis.", contentType),
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to analyze content with Claude")
		return nil, err
	}

	return analysis, nil
}

// StartRecommendationConversation starts a conversational recommendation session
func (c *ClaudeClient) StartRecommendationConversation(ctx context.Context, contentType string, preferences map[string]interface{}, systemInstructions string) (string, string, error) {
	log := utils.LoggerFromContext(ctx)
	log.Info().
		Str("contentType", contentType).
		Interface("preferences", preferences).
		Msg("Starting recommendation conversation with Claude")

	// Generate a unique, URL-safe conversation ID
	conversationID := fmt.Sprintf("rec-%d-%s", c.ClientID, utils.GenerateRandomID(12))

	// Build system instructions if not provided
	if systemInstructions == "" {
		systemInstructions = fmt.Sprintf(
			"You are an expert %s recommendation assistant. Your goal is to help the user discover %s they'll love based on their preferences and interests. "+
				"Maintain a friendly, conversational tone. Ask questions to understand their preferences better. "+
				"When recommending items, provide a brief explanation of why you're recommending them based on the user's preferences.",
			contentType, contentType)
	}

	// Create a personalized welcome message based on content type and preferences
	welcomeMessage := fmt.Sprintf("Hi there! I'm your %s recommendation assistant. ", contentType)

	// Add information about preferences if provided
	if len(preferences) > 0 {
		welcomeMessage += "Based on your preferences, "

		if genres, ok := preferences["favoriteGenres"].([]interface{}); ok && len(genres) > 0 {
			welcomeMessage += fmt.Sprintf("I see you enjoy %s like ", contentType)
			for i, genre := range genres {
				if i > 0 {
					if i == len(genres)-1 {
						welcomeMessage += " and "
					} else {
						welcomeMessage += ", "
					}
				}
				welcomeMessage += fmt.Sprintf("%v", genre)
			}
			welcomeMessage += ". "
		}

		if recent, ok := preferences["recentlyWatched"].([]interface{}); ok && len(recent) > 0 {
			welcomeMessage += "You've recently enjoyed "
			for i, item := range recent {
				if i > 0 {
					if i == len(recent)-1 {
						welcomeMessage += " and "
					} else {
						welcomeMessage += ", "
					}
				}
				welcomeMessage += fmt.Sprintf("%v", item)
			}
			welcomeMessage += ". "
		}
	}

	welcomeMessage += fmt.Sprintf("What kind of %s are you in the mood for today?", contentType)

	// Initialize the conversation context
	c.conversations[conversationID] = ConversationContext{
		ContentType:        contentType,
		UserPreferences:    preferences,
		SystemInstructions: systemInstructions,
		History: []ChatMessage{
			{
				Role:    "assistant",
				Content: welcomeMessage,
			},
		},
	}

	return conversationID, welcomeMessage, nil
}

// ContinueRecommendationConversation continues an existing conversation with a new message
func (c *ClaudeClient) ContinueRecommendationConversation(ctx context.Context, conversationID string, message string, context map[string]interface{}) (string, []map[string]interface{}, error) {
	log := utils.LoggerFromContext(ctx)
	log.Info().
		Str("conversationID", conversationID).
		Str("message", message).
		Msg("Continuing recommendation conversation with Claude")

	// Check if the conversation exists
	conversation, exists := c.conversations[conversationID]
	if !exists {
		return "", nil, fmt.Errorf("conversation not found: %s", conversationID)
	}

	// Add the user message to history
	conversation.History = append(conversation.History, ChatMessage{
		Role:    "user",
		Content: message,
	})

	// Build a prompt that includes the conversation history and context
	var promptBuilder strings.Builder

	// Add system instructions
	promptBuilder.WriteString(conversation.SystemInstructions)
	promptBuilder.WriteString("\n\n")

	// Add user preferences
	if len(conversation.UserPreferences) > 0 {
		promptBuilder.WriteString("User preferences:\n")
		for k, v := range conversation.UserPreferences {
			promptBuilder.WriteString(fmt.Sprintf("- %s: %v\n", k, v))
		}
		promptBuilder.WriteString("\n")
	}

	// Add the conversation history
	promptBuilder.WriteString("Conversation history:\n")
	for _, msg := range conversation.History {
		promptBuilder.WriteString(fmt.Sprintf("%s: %s\n", msg.Role, msg.Content))
	}

	// Check if we should extract recommendations
	shouldExtractRecommendations := false
	if context != nil {
		if extract, ok := context["extractRecommendations"].(bool); ok {
			shouldExtractRecommendations = extract
		}
	}

	// Add special instructions for recommendations if needed
	if shouldExtractRecommendations {
		promptBuilder.WriteString("\nPlease include specific recommendations in your response. Format each recommendation as a clear item with relevant details.\n")
	}

	// Generate the AI response
	aiResponse, err := c.GenerateText(ctx, promptBuilder.String(), &aitypes.GenerationOptions{
		Temperature: 0.7,
		MaxTokens:   1000,
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to continue conversation with Claude")
		return "", nil, err
	}

	// Add the AI response to the conversation history
	conversation.History = append(conversation.History, ChatMessage{
		Role:    "assistant",
		Content: aiResponse,
	})

	// Update the conversation in the map
	c.conversations[conversationID] = conversation

	// Extract recommendations if requested
	var recommendations []map[string]interface{}

	if shouldExtractRecommendations {
		// Extract recommendations from the response text
		recommendations = c.extractRecommendationsFromText(ctx, aiResponse, conversation.ContentType)
	}

	return aiResponse, recommendations, nil
}

// extractRecommendationsFromText parses the AI response to extract structured recommendations
func (c *ClaudeClient) extractRecommendationsFromText(ctx context.Context, text string, contentType string) []map[string]interface{} {
	log := utils.LoggerFromContext(ctx)

	// Create a prompt to extract recommendations
	extractPrompt := fmt.Sprintf(
		"From the following assistant response, extract a list of %s recommendations as structured data:\n\n%s\n\n"+
			"Return ONLY a JSON array of recommendation objects. Each object should have appropriate fields for %s items.",
		contentType, text, contentType)

	// Create the output schema to receive the data
	var recommendations []map[string]interface{}

	// Generate the structured response
	err := c.GenerateStructured(ctx, extractPrompt, &recommendations, &aitypes.GenerationOptions{
		Temperature:        0.1, // Low temperature for deterministic extraction
		MaxTokens:          1000,
		SystemInstructions: "You are a helpful data extraction assistant. Your job is to extract structured recommendations from text.",
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to extract recommendations from text")
		return []map[string]interface{}{}
	}

	return recommendations
}

// CreateMessage implements the AIClient interface for Claude
func (c *ClaudeClient) CreateMessage(ctx context.Context, messageRequest aitypes.MessageRequest) (*aitypes.MessageResponse, error) {
	log := utils.LoggerFromContext(ctx)
	log.Info().
		Str("model", messageRequest.Model).
		Int("messageCount", len(messageRequest.Messages)).
		Msg("Creating message with Claude")

	// Extract system message if provided
	systemPrompt := messageRequest.System
	if systemPrompt == "" {
		// Default system prompt
		systemPrompt = "You are Claude, a helpful AI assistant by Anthropic."
	}

	// Build the conversation history from the messages
	var userPrompt strings.Builder
	for _, msg := range messageRequest.Messages {
		// Skip system messages as they're handled separately
		if msg.Role == "system" {
			continue
		}

		// Add each message content
		for _, content := range msg.Content {
			if content.Type == "text" {
				userPrompt.WriteString(content.Text)
				userPrompt.WriteString("\n\n")
			}
		}
	}

	// Set options for generation
	options := map[string]interface{}{
		"temperature": messageRequest.Temperature,
		"max_tokens":  messageRequest.MaxTokens,
	}

	// Set model or use default
	model := messageRequest.Model
	if model == "" {
		model = c.config.Model // Use default model from config
	}

	// Call GenerateContent to handle the actual request
	response, err := c.GenerateContent(ctx, systemPrompt, userPrompt.String(), model, options)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create message with Claude")
		return nil, err
	}

	// Build the response object
	messageContent := []aitypes.MessageContent{
		{
			Type: "text",
			Text: response.Text,
		},
	}

	msgResponse := &aitypes.MessageResponse{
		ID:         utils.GenerateRandomID(16),
		Object:     "message",
		Created:    time.Now().Unix(),
		Model:      model,
		Content:    messageContent,
		Role:       "assistant",
		TokenUsage: response.TokenUsage,
	}

	return msgResponse, nil
}

// GenerateContent implements the AIClient interface for generating content with Claude
func (c *ClaudeClient) GenerateContent(ctx context.Context, systemPrompt string, userPrompt string, model string, options map[string]interface{}) (*aitypes.ContentResponse, error) {
	log := utils.LoggerFromContext(ctx)
	log.Info().
		Str("model", model).
		Int("userPromptLength", len(userPrompt)).
		Int("systemPromptLength", len(systemPrompt)).
		Msg("Generating content with Claude")

	// Apply custom options if provided
	tempToUse := c.config.Temperature
	maxTokensToUse := c.config.MaxTokens

	if options != nil {
		if temp, ok := options["temperature"].(float64); ok && temp > 0 {
			tempToUse = temp
		}
		if maxTokens, ok := options["max_tokens"].(int); ok && maxTokens > 0 {
			maxTokensToUse = maxTokens
		}
	}

	c.llm.SetOption("temperature", tempToUse)
	c.llm.SetOption("max_tokens", maxTokensToUse)
	c.llm.SetOption("system", systemPrompt)
	c.llm.SetOption("retries", 3)

	prompt := c.llm.NewPrompt(userPrompt)
	// Generate the response
	responseText, err := c.llm.Generate(ctx, prompt)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate content with Claude")
		return nil, fmt.Errorf("claude content generation failed: %w", err)
	}

	// Estimate token usage (actual token count would require a separate API call)
	// This is a very rough estimate
	promptTokens := len(systemPrompt)/4 + len(userPrompt)/4
	completionTokens := len(responseText) / 4

	// Create the response
	response := &aitypes.ContentResponse{
		Text: responseText,
		TokenUsage: aitypes.TokenUsage{
			PromptTokens:     promptTokens,
			CompletionTokens: completionTokens,
			TotalTokens:      promptTokens + completionTokens,
		},
	}

	return response, nil
}
