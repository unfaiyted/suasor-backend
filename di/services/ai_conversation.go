package services

import (
	"context"
	"fmt"
	"suasor/clients/types"
	"suasor/repository"
	"suasor/services"

	"github.com/sarulabs/di"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RegisterAIConversationService registers the AI conversation service in the DI container
func RegisterAIConversationService(builder *di.Builder) {
	builder.Add(di.Def{
		Name: "ai_conversation_service",
		Build: func(ctn di.Container) (interface{}, error) {
			logger := ctn.Get("logger").(*zap.SugaredLogger)
			logger.Info("Building AI conversation service")

			ctx := context.Background()

			// Get the DB from the container
			db, ok := ctn.Get("db").(*gorm.DB)
			if !ok {
				return nil, fmt.Errorf("failed to get DB from container")
			}

			// Get the client factory from the container
			clientFactory, ok := ctn.Get("client_factory").(*services.ClientProviderFactoryService)
			if !ok {
				return nil, fmt.Errorf("failed to get client factory from container")
			}

			// Get the AI client service from the container
			clientService, ok := ctn.Get("ai_client_service").(services.ClientService[types.AIClientConfig])
			if !ok {
				return nil, fmt.Errorf("failed to get AI client service from container")
			}

			// Create the AI conversation repository
			repo := repository.NewGormAIConversationRepository(db)

			// Create the AI conversation service
			service := services.NewAIConversationService(
				repo,
				clientService,
				clientFactory,
			)

			logger.Info("AI conversation service built successfully")
			return service, nil
		},
		Close: func(obj interface{}) error {
			// Nothing to close
			return nil
		},
	})
}