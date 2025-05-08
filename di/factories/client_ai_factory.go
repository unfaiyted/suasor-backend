package factories

import (
	"context"
	"fmt"

	"suasor/clients"
	"suasor/clients/ai"
	claude "suasor/clients/ai/claude"
	clienttypes "suasor/clients/types"
	"suasor/di/container"
	"suasor/utils/logger"
)

// RegisterAIClientFactories registers all AI client factories
func RegisterAIClientFactories(ctx context.Context, c *container.Container) {
	// Get the service to register factories
	service := clients.GetClientProviderFactoryService()

	// Register Claude AI client
	registerAIClientProvider(ctx, service, clienttypes.ClientTypeClaude, func(ctx context.Context, clientID uint64, config *clienttypes.ClaudeConfig) (ai.ClientAI, error) {
		return claude.NewClaudeClient(ctx, clientID, *config)
	})

	// Register other AI clients here as needed (OpenAI, Ollama, etc.)
}

// registerAIClientProvider registers an AI client provider with the factory service
func registerAIClientProvider[T clienttypes.AIClientConfig](
	ctx context.Context,
	service *clients.ClientProviderFactoryService,
	clientType clienttypes.ClientType,
	createFn func(ctx context.Context, clientID uint64, config T) (ai.ClientAI, error),
) {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Str("clientType", clientType.String()).
		Msg("Registering AI client provider factory")

	service.RegisterClientProviderFactory(clientType, func(ctx context.Context, clientID uint64, config clienttypes.ClientConfig) (clients.Client, error) {
		log := logger.LoggerFromContext(ctx)
		
		// Type assertion to the specific config type
		typedConfig, ok := config.(T)
		if !ok {
			err := fmt.Errorf("invalid config type for %s: expected %T but got %T", clientType, *new(T), config)
			log.Error().Err(err).Msg("Client factory type mismatch")
			return nil, err
		}

		// Create the client using the provided creation function
		aiClient, err := createFn(ctx, clientID, typedConfig)
		if err != nil {
			log.Error().Err(err).
				Str("clientType", clientType.String()).
				Uint64("clientID", clientID).
				Msg("Failed to create AI client")
			return nil, err
		}

		return aiClient, nil
	})
}