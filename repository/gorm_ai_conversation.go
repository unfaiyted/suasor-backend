package repository

import (
	"context"
	"errors"
	"fmt"
	"suasor/types/models"
	"suasor/utils/logger"
	"time"

	"gorm.io/gorm"
)

// gormAIConversationRepository implements the AIConversationRepository interface using GORM
type gormAIConversationRepository struct {
	db *gorm.DB
}

// NewGormAIConversationRepository creates a new AI conversation repository with GORM
func NewGormAIConversationRepository(db *gorm.DB) AIConversationRepository {
	return &gormAIConversationRepository{
		db: db,
	}
}

// CreateConversation creates a new AI conversation in the database
func (r *gormAIConversationRepository) CreateConversation(ctx context.Context, conversation *models.AIConversation) (string, error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Str("conversationID", conversation.ID).
		Uint64("userID", conversation.UserID).
		Uint64("clientID", conversation.ClientID).
		Str("contentType", conversation.ContentType).
		Msg("Creating new AI conversation")

	result := r.db.WithContext(ctx).Create(conversation)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to create AI conversation")
		return "", fmt.Errorf("failed to create conversation: %w", result.Error)
	}

	return conversation.ID, nil
}

// GetConversationByID retrieves an AI conversation by ID
func (r *gormAIConversationRepository) GetConversationByID(ctx context.Context, conversationID string) (*models.AIConversation, error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().Str("conversationID", conversationID).Msg("Getting AI conversation by ID")

	var conversation models.AIConversation
	result := r.db.WithContext(ctx).First(&conversation, "id = ?", conversationID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Not found
		}
		log.Error().Err(result.Error).Msg("Failed to get AI conversation")
		return nil, fmt.Errorf("failed to get conversation: %w", result.Error)
	}

	return &conversation, nil
}

// GetConversationsByUserID retrieves conversations for a user with pagination
func (r *gormAIConversationRepository) GetConversationsByUserID(ctx context.Context, userID uint64, limit, offset int) ([]*models.AIConversation, int, error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", userID).
		Int("limit", limit).
		Int("offset", offset).
		Msg("Getting AI conversations for user")

	// Set default limit if not provided
	if limit <= 0 {
		limit = 20
	}

	// Get total count first
	var count int64
	result := r.db.WithContext(ctx).Model(&models.AIConversation{}).Where("user_id = ?", userID).Count(&count)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to count AI conversations")
		return nil, 0, fmt.Errorf("failed to count conversations: %w", result.Error)
	}

	// Now get the actual conversations
	var conversations []*models.AIConversation
	result = r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("updated_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&conversations)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to get AI conversations")
		return nil, 0, fmt.Errorf("failed to get conversations: %w", result.Error)
	}

	return conversations, int(count), nil
}

// UpdateConversationStatus updates the status of a conversation
func (r *gormAIConversationRepository) UpdateConversationStatus(ctx context.Context, conversationID string, status string) error {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Str("conversationID", conversationID).
		Str("status", status).
		Msg("Updating AI conversation status")

	result := r.db.WithContext(ctx).
		Model(&models.AIConversation{}).
		Where("id = ?", conversationID).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to update AI conversation status")
		return fmt.Errorf("failed to update conversation status: %w", result.Error)
	}

	return nil
}

// DeleteConversation deletes a conversation and all related data
func (r *gormAIConversationRepository) DeleteConversation(ctx context.Context, conversationID string) error {
	log := logger.LoggerFromContext(ctx)
	log.Info().Str("conversationID", conversationID).Msg("Deleting AI conversation")

	// Start a transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		log.Error().Err(tx.Error).Msg("Failed to begin transaction")
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Delete related analytics (if exists)
	if err := tx.Where("conversation_id = ?", conversationID).Delete(&models.AIConversationAnalytics{}).Error; err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("Failed to delete AI conversation analytics")
		return fmt.Errorf("failed to delete conversation analytics: %w", err)
	}

	// Delete related recommendations
	if err := tx.Where("conversation_id = ?", conversationID).Delete(&models.AIRecommendation{}).Error; err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("Failed to delete AI recommendations")
		return fmt.Errorf("failed to delete recommendations: %w", err)
	}

	// Delete related messages
	if err := tx.Where("conversation_id = ?", conversationID).Delete(&models.AIMessage{}).Error; err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("Failed to delete AI messages")
		return fmt.Errorf("failed to delete messages: %w", err)
	}

	// Delete the conversation
	if err := tx.Delete(&models.AIConversation{}, "id = ?", conversationID).Error; err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("Failed to delete AI conversation")
		return fmt.Errorf("failed to delete conversation: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ArchiveOldConversations archives conversations older than the specified duration
func (r *gormAIConversationRepository) ArchiveOldConversations(ctx context.Context, olderThan time.Duration) (int, error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().Str("olderThan", olderThan.String()).Msg("Archiving old AI conversations")

	cutoffTime := time.Now().Add(-olderThan)

	result := r.db.WithContext(ctx).
		Model(&models.AIConversation{}).
		Where("status = ? AND updated_at < ?", "active", cutoffTime).
		Updates(map[string]interface{}{
			"status":     "archived",
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to archive old AI conversations")
		return 0, fmt.Errorf("failed to archive old conversations: %w", result.Error)
	}

	return int(result.RowsAffected), nil
}

// AddMessage adds a new message to a conversation
func (r *gormAIConversationRepository) AddMessage(ctx context.Context, message *models.AIMessage) (string, error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Str("messageID", message.ID).
		Str("conversationID", message.ConversationID).
		Str("role", message.Role).
		Msg("Adding new AI message")

	// Start a transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		log.Error().Err(tx.Error).Msg("Failed to begin transaction")
		return "", fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Insert the message
	if err := tx.Create(message).Error; err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("Failed to add AI message")
		return "", fmt.Errorf("failed to add message: %w", err)
	}

	// Update the conversation's message count and last message time
	if err := tx.Model(&models.AIConversation{}).
		Where("id = ?", message.ConversationID).
		Updates(map[string]interface{}{
			"message_count":     gorm.Expr("message_count + 1"),
			"last_message_time": message.Timestamp,
			"updated_at":        time.Now(),
		}).Error; err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("Failed to update conversation")
		return "", fmt.Errorf("failed to update conversation: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return message.ID, nil
}

// GetMessagesByConversationID retrieves messages for a conversation with pagination
func (r *gormAIConversationRepository) GetMessagesByConversationID(ctx context.Context, conversationID string, limit, offset int) ([]*models.AIMessage, int, error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Str("conversationID", conversationID).
		Int("limit", limit).
		Int("offset", offset).
		Msg("Getting messages for conversation")

	// Set default limit if not provided
	if limit <= 0 {
		limit = 50
	}

	// Get total count first
	var count int64
	result := r.db.WithContext(ctx).
		Model(&models.AIMessage{}).
		Where("conversation_id = ?", conversationID).
		Count(&count)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to count AI messages")
		return nil, 0, fmt.Errorf("failed to count messages: %w", result.Error)
	}

	// Now get the actual messages
	var messages []*models.AIMessage
	result = r.db.WithContext(ctx).
		Where("conversation_id = ?", conversationID).
		Order("timestamp ASC").
		Limit(limit).
		Offset(offset).
		Find(&messages)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to get AI messages")
		return nil, 0, fmt.Errorf("failed to get messages: %w", result.Error)
	}

	return messages, int(count), nil
}

// GetConversationHistory retrieves the full history of a conversation
func (r *gormAIConversationRepository) GetConversationHistory(ctx context.Context, conversationID string) ([]*models.AIMessage, error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().Str("conversationID", conversationID).Msg("Getting full conversation history")

	var messages []*models.AIMessage
	result := r.db.WithContext(ctx).
		Where("conversation_id = ?", conversationID).
		Order("timestamp ASC").
		Find(&messages)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to get conversation history")
		return nil, fmt.Errorf("failed to get conversation history: %w", result.Error)
	}

	return messages, nil
}

// AddRecommendation adds a new recommendation extracted from a conversation
func (r *gormAIConversationRepository) AddRecommendation(ctx context.Context, recommendation *models.AIRecommendation) (string, error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Str("recommendationID", recommendation.ID).
		Str("conversationID", recommendation.ConversationID).
		Str("itemType", recommendation.ItemType).
		Str("title", recommendation.Title).
		Msg("Adding new AI recommendation")

	result := r.db.WithContext(ctx).Create(recommendation)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to add AI recommendation")
		return "", fmt.Errorf("failed to add recommendation: %w", result.Error)
	}

	return recommendation.ID, nil
}

// GetRecommendationsByConversationID retrieves all recommendations for a conversation
func (r *gormAIConversationRepository) GetRecommendationsByConversationID(ctx context.Context, conversationID string) ([]*models.AIRecommendation, error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().Str("conversationID", conversationID).Msg("Getting recommendations for conversation")

	var recommendations []*models.AIRecommendation
	result := r.db.WithContext(ctx).
		Where("conversation_id = ?", conversationID).
		Order("created_at DESC").
		Find(&recommendations)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to get recommendations")
		return nil, fmt.Errorf("failed to get recommendations: %w", result.Error)
	}

	return recommendations, nil
}

// GetRecommendationsByUserID retrieves recommendations for a user with pagination
func (r *gormAIConversationRepository) GetRecommendationsByUserID(ctx context.Context, userID uint64, limit, offset int) ([]*models.AIRecommendation, int, error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", userID).
		Int("limit", limit).
		Int("offset", offset).
		Msg("Getting recommendations for user")

	// Set default limit if not provided
	if limit <= 0 {
		limit = 20
	}

	// Get total count first
	var count int64
	result := r.db.WithContext(ctx).
		Model(&models.AIRecommendation{}).
		Where("user_id = ?", userID).
		Count(&count)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to count recommendations")
		return nil, 0, fmt.Errorf("failed to count recommendations: %w", result.Error)
	}

	// Now get the actual recommendations
	var recommendations []*models.AIRecommendation
	result = r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&recommendations)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to get recommendations")
		return nil, 0, fmt.Errorf("failed to get recommendations: %w", result.Error)
	}

	return recommendations, int(count), nil
}

// UpdateRecommendationSelection updates the selection status of a recommendation
func (r *gormAIConversationRepository) UpdateRecommendationSelection(ctx context.Context, recommendationID string, selected bool) error {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Str("recommendationID", recommendationID).
		Bool("selected", selected).
		Msg("Updating recommendation selection")

	updates := map[string]interface{}{
		"selected": selected,
	}

	if selected {
		updates["selected_at"] = time.Now()
	} else {
		updates["selected_at"] = nil
	}

	result := r.db.WithContext(ctx).
		Model(&models.AIRecommendation{}).
		Where("id = ?", recommendationID).
		Updates(updates)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to update recommendation selection")
		return fmt.Errorf("failed to update recommendation selection: %w", result.Error)
	}

	return nil
}

// GetConversationAnalytics retrieves analytics for a conversation
func (r *gormAIConversationRepository) GetConversationAnalytics(ctx context.Context, conversationID string) (*models.AIConversationAnalytics, error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().Str("conversationID", conversationID).Msg("Getting conversation analytics")

	var analytics models.AIConversationAnalytics
	result := r.db.WithContext(ctx).
		Where("conversation_id = ?", conversationID).
		First(&analytics)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// If no analytics exist yet, create a new one
			return models.NewAIConversationAnalytics(conversationID), nil
		}
		log.Error().Err(result.Error).Msg("Failed to get conversation analytics")
		return nil, fmt.Errorf("failed to get conversation analytics: %w", result.Error)
	}

	return &analytics, nil
}

// UpdateConversationAnalytics updates analytics for a conversation
func (r *gormAIConversationRepository) UpdateConversationAnalytics(ctx context.Context, analytics *models.AIConversationAnalytics) error {
	log := logger.LoggerFromContext(ctx)
	log.Info().Str("conversationID", analytics.ConversationID).Msg("Updating conversation analytics")

	// Check if analytics already exist
	var count int64
	result := r.db.WithContext(ctx).
		Model(&models.AIConversationAnalytics{}).
		Where("conversation_id = ?", analytics.ConversationID).
		Count(&count)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to check if analytics exist")
		return fmt.Errorf("failed to check if analytics exist: %w", result.Error)
	}

	if count > 0 {
		// Update existing analytics
		result = r.db.WithContext(ctx).
			Where("conversation_id = ?", analytics.ConversationID).
			Updates(map[string]interface{}{
				"total_user_messages":       analytics.TotalUserMessages,
				"total_assistant_messages":  analytics.TotalAssistantMessages,
				"total_recommendations":     analytics.TotalRecommendations,
				"selected_recommendations":  analytics.SelectedRecommendations,
				"total_tokens_used":         analytics.TotalTokensUsed,
				"average_response_time":     analytics.AverageResponseTime,
				"conversation_duration":     analytics.ConversationDuration,
				"last_updated_at":           time.Now(),
			})
	} else {
		// Insert new analytics
		result = r.db.WithContext(ctx).Create(analytics)
	}

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to update conversation analytics")
		return fmt.Errorf("failed to update conversation analytics: %w", result.Error)
	}

	return nil
}

// GetUserConversationStats retrieves aggregated statistics for a user's AI interactions
func (r *gormAIConversationRepository) GetUserConversationStats(ctx context.Context, userID uint64) (*models.UserConversationStats, error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().Uint64("userID", userID).Msg("Getting user conversation stats")

	stats := &models.UserConversationStats{
		UserID: userID,
	}

	// Get total conversations count
	var totalConversations int64
	result := r.db.WithContext(ctx).
		Model(&models.AIConversation{}).
		Where("user_id = ?", userID).
		Count(&totalConversations)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to count user conversations")
		return nil, fmt.Errorf("failed to count user conversations: %w", result.Error)
	}

	stats.TotalConversations = int(totalConversations)

	// Get total messages
	var totalMessages int64
	result = r.db.WithContext(ctx).
		Model(&models.AIConversation{}).
		Where("user_id = ?", userID).
		Select("SUM(message_count)").
		Scan(&totalMessages)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to count user messages")
		return nil, fmt.Errorf("failed to count user messages: %w", result.Error)
	}

	stats.TotalMessages = int(totalMessages)

	// Calculate average messages per conversation
	if totalConversations > 0 {
		stats.AverageMessagesPerConv = float64(totalMessages) / float64(totalConversations)
	}

	// Get total recommendations
	var totalRecommendations int64
	result = r.db.WithContext(ctx).
		Model(&models.AIRecommendation{}).
		Where("user_id = ?", userID).
		Count(&totalRecommendations)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to count user recommendations")
		return nil, fmt.Errorf("failed to count user recommendations: %w", result.Error)
	}

	stats.TotalRecommendations = int(totalRecommendations)

	// Get selected recommendations
	var selectedRecommendations int64
	result = r.db.WithContext(ctx).
		Model(&models.AIRecommendation{}).
		Where("user_id = ? AND selected = true", userID).
		Count(&selectedRecommendations)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to count selected recommendations")
		return nil, fmt.Errorf("failed to count selected recommendations: %w", result.Error)
	}

	stats.SelectedRecommendations = int(selectedRecommendations)

	// Get favorite content type (most used)
	type contentTypeCount struct {
		ContentType string
		Count       int
	}

	var favoriteContentType contentTypeCount
	result = r.db.WithContext(ctx).
		Model(&models.AIConversation{}).
		Select("content_type, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("content_type").
		Order("count DESC").
		Limit(1).
		Scan(&favoriteContentType)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to get favorite content type")
		return nil, fmt.Errorf("failed to get favorite content type: %w", result.Error)
	}

	stats.FavoriteContentType = favoriteContentType.ContentType

	// Get last conversation time
	var lastConversation models.AIConversation
	result = r.db.WithContext(ctx).
		Model(&models.AIConversation{}).
		Where("user_id = ?", userID).
		Order("updated_at DESC").
		Limit(1).
		Select("updated_at").
		Scan(&lastConversation)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Error().Err(result.Error).Msg("Failed to get last conversation time")
		return nil, fmt.Errorf("failed to get last conversation time: %w", result.Error)
	}

	if !lastConversation.UpdatedAt.IsZero() {
		stats.LastConversation = lastConversation.UpdatedAt
	} else {
		stats.LastConversation = time.Now()
	}

	return stats, nil
}