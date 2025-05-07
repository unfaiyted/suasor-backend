package database

import (
	"context"
	"fmt"
	"suasor/types/models"
	"suasor/utils/logger"

	"gorm.io/gorm"
)

// RegisterAIConversationModels registers the AI conversation models with GORM
func RegisterAIConversationModels(ctx context.Context, db *gorm.DB) error {
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Registering AI conversation models with GORM")

	// Auto migrate the AI conversation models
	if err := db.AutoMigrate(
		&models.AIConversation{},
		&models.AIMessage{},
		&models.AIRecommendation{},
		&models.AIConversationAnalytics{},
	); err != nil {
		return fmt.Errorf("failed to migrate AI conversation models: %w", err)
	}

	log.Info().Msg("AI conversation models registered successfully")
	return nil
}