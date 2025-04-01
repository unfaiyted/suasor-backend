package claude

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/teilomillet/gollm"
	"suasor/client"
	aitypes "suasor/client/ai/types"
	"suasor/client/types"
	"suasor/utils"
)

// ClaudeClient implements the AI client interface
type ClaudeClient struct {
	client.BaseClient
	llm      gollm.LLM
	config   types.ClaudeConfig
	memoryID string // For conversation tracking
}

// NewClaudeClient creates a new Claude client instance using gollm
func NewClaudeClient(ctx context.Context, clientID uint64, cfg types.ClaudeConfig) (client.Client, error) {
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
		llm:    llm,
		config: cfg,
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
	structuredPrompt := promptText + "\n\nRespond only with valid JSON matching the required schema."

	prompt := gollm.NewPrompt(structuredPrompt,
		gollm.WithOutput("Respond with valid JSON that can be parsed into the required structure."))

	// Generate with JSON validation
	response, err := c.llm.Generate(ctx, prompt, gollm.WithJSONSchemaValidation())
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate structured output with Claude")
		return fmt.Errorf("claude structured generation failed: %w", err)
	}

	// Parse the JSON response into the provided schema
	if err := json.Unmarshal([]byte(response), outputSchema); err != nil {
		log.Error().Err(err).Msg("Failed to parse Claude response as JSON")
		return fmt.Errorf("failed to parse Claude JSON response: %w", err)
	}

	return nil
}

// StartConversation begins a new conversation with Claude
func (c *ClaudeClient) StartConversation(ctx context.Context, systemInstructions string) (string, error) {
	// Generate a unique conversation ID
	c.memoryID = fmt.Sprintf("conv-%d-%s", c.ClientID, utils.GenerateRandomID(8))

	// Enable memory context for this conversation
	c.llm.SetOption("memory_id", c.memoryID)

	// Set system instructions if provided
	if systemInstructions != "" {
		c.llm.SetOption("system_instructions", systemInstructions)
	}

	return c.memoryID, nil
}

// SendMessage sends a message in an existing conversation
func (c *ClaudeClient) SendMessage(ctx context.Context, conversationID string, message string) (string, error) {
	log := utils.LoggerFromContext(ctx)

	// Verify this is the active conversation
	if c.memoryID != conversationID {
		c.llm.SetOption("memory_id", conversationID)
		c.memoryID = conversationID
	}

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

