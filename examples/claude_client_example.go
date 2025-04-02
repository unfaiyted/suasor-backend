package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"path/filepath"
	"suasor/client/ai/claude"
	aitypes "suasor/client/ai/types"
	"suasor/client/types"
)

func init() {
	// Try to load .env file from several possible locations
	locations := []string{
		".env",    // Current directory
		"../.env", // Project root
		filepath.Join(os.Getenv("HOME"), "claude_test.env"), // Home directory
	}

	for _, location := range locations {
		err := godotenv.Load(location)
		if err == nil {
			fmt.Printf("Loaded environment from: %s\n", location)
			break
		}
	}
}

func main() {
	// Get API key from environment
	apiKey := os.Getenv("CLAUDE_API_KEY")
	if apiKey == "" {
		log.Fatal("CLAUDE_API_KEY environment variable not set")
	}

	// Create a Claude config
	config := types.ClaudeConfig{
		BaseAIClientConfig: types.BaseAIClientConfig{
			BaseClientConfig: types.BaseClientConfig{
				Type:     types.ClientTypeClaude,
				Category: types.ClientCategoryAI,
			},
			ClientType:       types.AIClientTypeClaude,
			APIKey:           apiKey,
			Model:            "claude-3-haiku-20240307", // Use smallest model for example
			Temperature:      0.7,
			MaxTokens:        150,
			MaxContextTokens: 2000,
		},
	}

	// Create a Claude client
	ctx := context.Background()
	client, err := claude.NewClaudeClient(ctx, 1, config)
	if err != nil {
		log.Fatalf("Failed to create Claude client: %v", err)
	}

	// Test connection
	ok, err := client.TestConnection(ctx)
	if err != nil {
		log.Fatalf("Failed to test connection: %v", err)
	}
	if !ok {
		log.Fatal("Connection test failed")
	}
	fmt.Println("✅ Connection to Claude API successful")

	// Basic text generation
	prompt := "Explain the concept of neural networks in 2 sentences."
	response, err := client.GenerateText(ctx, prompt, &aitypes.GenerationOptions{
		Temperature: 0.5,
		MaxTokens:   100,
	})
	if err != nil {
		log.Fatalf("Failed to generate text: %v", err)
	}
	fmt.Printf("\n🤖 Text Generation:\n%s\n\n", response)

	// Structured output - simpler example
	type MathProblem struct {
		Problem string   `json:"problem"`
		Answer  int      `json:"answer"`
		Steps   []string `json:"steps"`
	}

	var math MathProblem
	structuredPrompt := "Solve the following math problem: What is 12 * 5 + 3? Return a JSON object with the problem, answer (as a number), and steps fields."

	err = client.GenerateStructured(ctx, structuredPrompt, &math, nil)
	if err != nil {
		log.Fatalf("Failed to generate structured output: %v", err)
	}

	fmt.Println("🧮 Math Problem Solution:")
	fmt.Printf("  Problem: %s\n", math.Problem)
	fmt.Printf("  Answer: %d\n", math.Answer)
	fmt.Printf("  Steps: %s\n", math.Steps)
	fmt.Println()

	// For conversation example, we'll use separate calls instead of conversation API
	fmt.Println("💬 Sample Questions:")

	// First question
	msg1 := "What is the best programming language for a beginner to learn?"
	fmt.Printf("🧑 User: %s\n", msg1)

	resp1, err := client.GenerateText(ctx, msg1, nil)
	if err != nil {
		log.Fatalf("Failed to generate text: %v", err)
	}
	fmt.Printf("🤖 Claude: %s\n\n", resp1)

	// Second question
	msg2 := "What are some good resources for learning programming?"
	fmt.Printf("🧑 User: %s\n", msg2)

	resp2, err := client.GenerateText(ctx, msg2, nil)
	if err != nil {
		log.Fatalf("Failed to generate text: %v", err)
	}
	fmt.Printf("🤖 Claude: %s\n", resp2)

	// Print capabilities
	capabilities := client.GetCapabilities()
	fmt.Printf("\n⚙️ Claude Capabilities:\n")
	fmt.Printf("  Structured Output: %t\n", capabilities.SupportsStructuredOutput)
	fmt.Printf("  Conversation: %t\n", capabilities.SupportsConversation)
	fmt.Printf("  Streaming: %t\n", capabilities.SupportsStreaming)
	fmt.Printf("  Max Context Tokens: %d\n", capabilities.MaxContextTokens)
	fmt.Printf("  Default Max Tokens: %d\n", capabilities.DefaultMaxTokens)

	// Print supported models
	models := client.GetSupportedModels()
	fmt.Printf("\n📋 Supported Claude Models:\n")
	for _, model := range models {
		fmt.Printf("  - %s\n", model)
	}
}
