package services

import (
	"context"
	"suasor/clients"
	"suasor/clients/types"
	"suasor/di/container"
	"suasor/repository"
	"suasor/services"
	"suasor/utils/logger"

	"gorm.io/gorm"
)

// registerAIConversationService registers AI conversation related services
func registerAIConversationService(ctx context.Context, c *container.Container) {
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Registering AI conversation repository")

	// Register the AI conversation repository
	container.RegisterFactory[repository.AIConversationRepository](c, func(c *container.Container) repository.AIConversationRepository {
		db := container.MustGet[*gorm.DB](c)
		return repository.NewAIConversationRepository(db)
	})

	// Register the AI conversation service
	log.Info().Msg("Registering AI conversation service")
	container.RegisterFactory[services.AIConversationService](c, func(c *container.Container) services.AIConversationService {
		repo := container.MustGet[repository.AIConversationRepository](c)
		clientFactory := container.MustGet[*clients.ClientProviderFactoryService](c)
		claudeClientService := container.MustGet[services.ClientService[*types.ClaudeConfig]](c)
		openaiClientService := container.MustGet[services.ClientService[*types.OpenAIConfig]](c)
		ollamaClientService := container.MustGet[services.ClientService[*types.OllamaConfig]](c)
		clientHelper := container.MustGet[repository.ClientHelper](c)

		return services.NewAIConversationService(repo, claudeClientService, openaiClientService, ollamaClientService, clientHelper, clientFactory)
	})
}
