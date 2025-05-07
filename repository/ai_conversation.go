package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"suasor/types/models"
	"suasor/utils/logger"
	"time"

	"github.com/jmoiron/sqlx"
)

// AIConversationRepository defines the interface for managing AI conversations in the database
type AIConversationRepository interface {
	// Conversation methods
	CreateConversation(ctx context.Context, conversation *models.AIConversation) (string, error)
	GetConversationByID(ctx context.Context, conversationID string) (*models.AIConversation, error)
	GetConversationsByUserID(ctx context.Context, userID uint64, limit, offset int) ([]*models.AIConversation, int, error)
	UpdateConversationStatus(ctx context.Context, conversationID string, status string) error
	DeleteConversation(ctx context.Context, conversationID string) error
	ArchiveOldConversations(ctx context.Context, olderThan time.Duration) (int, error)
	
	// Message methods
	AddMessage(ctx context.Context, message *models.AIMessage) (string, error)
	GetMessagesByConversationID(ctx context.Context, conversationID string, limit, offset int) ([]*models.AIMessage, int, error)
	GetConversationHistory(ctx context.Context, conversationID string) ([]*models.AIMessage, error)
	
	// Recommendation methods
	AddRecommendation(ctx context.Context, recommendation *models.AIRecommendation) (string, error)
	GetRecommendationsByConversationID(ctx context.Context, conversationID string) ([]*models.AIRecommendation, error)
	GetRecommendationsByUserID(ctx context.Context, userID uint64, limit, offset int) ([]*models.AIRecommendation, int, error)
	UpdateRecommendationSelection(ctx context.Context, recommendationID string, selected bool) error
	
	// Analytics methods
	GetConversationAnalytics(ctx context.Context, conversationID string) (*models.AIConversationAnalytics, error)
	UpdateConversationAnalytics(ctx context.Context, analytics *models.AIConversationAnalytics) error
	GetUserConversationStats(ctx context.Context, userID uint64) (*models.UserConversationStats, error)
}

// sqlAIConversationRepository implements the AIConversationRepository interface using SQL
type sqlAIConversationRepository struct {
	db *sqlx.DB
}

// NewAIConversationRepository creates a new AI conversation repository
func NewAIConversationRepository(db *sqlx.DB) AIConversationRepository {
	return &sqlAIConversationRepository{
		db: db,
	}
}

// CreateConversation creates a new AI conversation in the database
func (r *sqlAIConversationRepository) CreateConversation(ctx context.Context, conversation *models.AIConversation) (string, error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Str("conversationID", conversation.ID).
		Uint64("userID", conversation.UserID).
		Uint64("clientID", conversation.ClientID).
		Str("contentType", conversation.ContentType).
		Msg("Creating new AI conversation")

	query := `
		INSERT INTO ai_conversations (
			id, user_id, client_id, content_type, status, system_prompt, 
			user_preferences, created_at, updated_at, expires_at, message_count
		) VALUES (
			:id, :user_id, :client_id, :content_type, :status, :system_prompt,
			:user_preferences, :created_at, :updated_at, :expires_at, :message_count
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, conversation)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create AI conversation")
		return "", fmt.Errorf("failed to create conversation: %w", err)
	}

	return conversation.ID, nil
}

// GetConversationByID retrieves an AI conversation by ID
func (r *sqlAIConversationRepository) GetConversationByID(ctx context.Context, conversationID string) (*models.AIConversation, error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().Str("conversationID", conversationID).Msg("Getting AI conversation by ID")

	query := `
		SELECT * FROM ai_conversations
		WHERE id = ?
	`

	var conversation models.AIConversation
	err := r.db.GetContext(ctx, &conversation, query, conversationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found
		}
		log.Error().Err(err).Msg("Failed to get AI conversation")
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	return &conversation, nil
}

// GetConversationsByUserID retrieves conversations for a user with pagination
func (r *sqlAIConversationRepository) GetConversationsByUserID(ctx context.Context, userID uint64, limit, offset int) ([]*models.AIConversation, int, error) {
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
	var count int
	countQuery := `SELECT COUNT(*) FROM ai_conversations WHERE user_id = ?`
	err := r.db.GetContext(ctx, &count, countQuery, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to count AI conversations")
		return nil, 0, fmt.Errorf("failed to count conversations: %w", err)
	}

	// Now get the actual conversations
	query := `
		SELECT * FROM ai_conversations
		WHERE user_id = ?
		ORDER BY updated_at DESC
		LIMIT ? OFFSET ?
	`

	var conversations []*models.AIConversation
	err = r.db.SelectContext(ctx, &conversations, query, userID, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get AI conversations")
		return nil, 0, fmt.Errorf("failed to get conversations: %w", err)
	}

	return conversations, count, nil
}

// UpdateConversationStatus updates the status of a conversation
func (r *sqlAIConversationRepository) UpdateConversationStatus(ctx context.Context, conversationID string, status string) error {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Str("conversationID", conversationID).
		Str("status", status).
		Msg("Updating AI conversation status")

	query := `
		UPDATE ai_conversations
		SET status = ?, updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query, status, conversationID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update AI conversation status")
		return fmt.Errorf("failed to update conversation status: %w", err)
	}

	return nil
}

// DeleteConversation deletes a conversation and all related data
func (r *sqlAIConversationRepository) DeleteConversation(ctx context.Context, conversationID string) error {
	log := logger.LoggerFromContext(ctx)
	log.Info().Str("conversationID", conversationID).Msg("Deleting AI conversation")

	// Start a transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to begin transaction")
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Delete related analytics (if exists)
	_, err = tx.ExecContext(ctx, "DELETE FROM ai_conversation_analytics WHERE conversation_id = ?", conversationID)
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("Failed to delete AI conversation analytics")
		return fmt.Errorf("failed to delete conversation analytics: %w", err)
	}

	// Delete related recommendations
	_, err = tx.ExecContext(ctx, "DELETE FROM ai_recommendations WHERE conversation_id = ?", conversationID)
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("Failed to delete AI recommendations")
		return fmt.Errorf("failed to delete recommendations: %w", err)
	}

	// Delete related messages
	_, err = tx.ExecContext(ctx, "DELETE FROM ai_messages WHERE conversation_id = ?", conversationID)
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("Failed to delete AI messages")
		return fmt.Errorf("failed to delete messages: %w", err)
	}

	// Delete the conversation
	_, err = tx.ExecContext(ctx, "DELETE FROM ai_conversations WHERE id = ?", conversationID)
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("Failed to delete AI conversation")
		return fmt.Errorf("failed to delete conversation: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ArchiveOldConversations archives conversations older than the specified duration
func (r *sqlAIConversationRepository) ArchiveOldConversations(ctx context.Context, olderThan time.Duration) (int, error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().Str("olderThan", olderThan.String()).Msg("Archiving old AI conversations")

	cutoffTime := time.Now().Add(-olderThan)

	query := `
		UPDATE ai_conversations
		SET status = 'archived', updated_at = NOW()
		WHERE status = 'active' AND updated_at < ?
	`

	result, err := r.db.ExecContext(ctx, query, cutoffTime)
	if err != nil {
		log.Error().Err(err).Msg("Failed to archive old AI conversations")
		return 0, fmt.Errorf("failed to archive old conversations: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get rows affected")
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return int(affected), nil
}

// AddMessage adds a new message to a conversation
func (r *sqlAIConversationRepository) AddMessage(ctx context.Context, message *models.AIMessage) (string, error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Str("messageID", message.ID).
		Str("conversationID", message.ConversationID).
		Str("role", message.Role).
		Msg("Adding new AI message")

	// Start a transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to begin transaction")
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Insert the message
	insertQuery := `
		INSERT INTO ai_messages (
			id, conversation_id, role, content, timestamp, metadata, token_usage
		) VALUES (
			:id, :conversation_id, :role, :content, :timestamp, :metadata, :token_usage
		)
	`

	_, err = tx.NamedExecContext(ctx, insertQuery, message)
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("Failed to add AI message")
		return "", fmt.Errorf("failed to add message: %w", err)
	}

	// Update the conversation's message count and last message time
	updateQuery := `
		UPDATE ai_conversations
		SET message_count = message_count + 1, 
		    last_message_time = ?,
		    updated_at = NOW()
		WHERE id = ?
	`

	_, err = tx.ExecContext(ctx, updateQuery, message.Timestamp, message.ConversationID)
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("Failed to update conversation")
		return "", fmt.Errorf("failed to update conversation: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return message.ID, nil
}

// GetMessagesByConversationID retrieves messages for a conversation with pagination
func (r *sqlAIConversationRepository) GetMessagesByConversationID(ctx context.Context, conversationID string, limit, offset int) ([]*models.AIMessage, int, error) {
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
	var count int
	countQuery := `SELECT COUNT(*) FROM ai_messages WHERE conversation_id = ?`
	err := r.db.GetContext(ctx, &count, countQuery, conversationID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to count AI messages")
		return nil, 0, fmt.Errorf("failed to count messages: %w", err)
	}

	// Now get the actual messages
	query := `
		SELECT * FROM ai_messages
		WHERE conversation_id = ?
		ORDER BY timestamp
		LIMIT ? OFFSET ?
	`

	var messages []*models.AIMessage
	err = r.db.SelectContext(ctx, &messages, query, conversationID, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get AI messages")
		return nil, 0, fmt.Errorf("failed to get messages: %w", err)
	}

	return messages, count, nil
}

// GetConversationHistory retrieves the full history of a conversation
func (r *sqlAIConversationRepository) GetConversationHistory(ctx context.Context, conversationID string) ([]*models.AIMessage, error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().Str("conversationID", conversationID).Msg("Getting full conversation history")

	query := `
		SELECT * FROM ai_messages
		WHERE conversation_id = ?
		ORDER BY timestamp
	`

	var messages []*models.AIMessage
	err := r.db.SelectContext(ctx, &messages, query, conversationID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get conversation history")
		return nil, fmt.Errorf("failed to get conversation history: %w", err)
	}

	return messages, nil
}

// AddRecommendation adds a new recommendation extracted from a conversation
func (r *sqlAIConversationRepository) AddRecommendation(ctx context.Context, recommendation *models.AIRecommendation) (string, error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Str("recommendationID", recommendation.ID).
		Str("conversationID", recommendation.ConversationID).
		Str("itemType", recommendation.ItemType).
		Str("title", recommendation.Title).
		Msg("Adding new AI recommendation")

	query := `
		INSERT INTO ai_recommendations (
			id, message_id, conversation_id, user_id, item_type, 
			title, external_id, data, reason, created_at, 
			selected, selected_at
		) VALUES (
			:id, :message_id, :conversation_id, :user_id, :item_type,
			:title, :external_id, :data, :reason, :created_at,
			:selected, :selected_at
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, recommendation)
	if err != nil {
		log.Error().Err(err).Msg("Failed to add AI recommendation")
		return "", fmt.Errorf("failed to add recommendation: %w", err)
	}

	return recommendation.ID, nil
}

// GetRecommendationsByConversationID retrieves all recommendations for a conversation
func (r *sqlAIConversationRepository) GetRecommendationsByConversationID(ctx context.Context, conversationID string) ([]*models.AIRecommendation, error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().Str("conversationID", conversationID).Msg("Getting recommendations for conversation")

	query := `
		SELECT * FROM ai_recommendations
		WHERE conversation_id = ?
		ORDER BY created_at DESC
	`

	var recommendations []*models.AIRecommendation
	err := r.db.SelectContext(ctx, &recommendations, query, conversationID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get recommendations")
		return nil, fmt.Errorf("failed to get recommendations: %w", err)
	}

	return recommendations, nil
}

// GetRecommendationsByUserID retrieves recommendations for a user with pagination
func (r *sqlAIConversationRepository) GetRecommendationsByUserID(ctx context.Context, userID uint64, limit, offset int) ([]*models.AIRecommendation, int, error) {
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
	var count int
	countQuery := `SELECT COUNT(*) FROM ai_recommendations WHERE user_id = ?`
	err := r.db.GetContext(ctx, &count, countQuery, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to count recommendations")
		return nil, 0, fmt.Errorf("failed to count recommendations: %w", err)
	}

	// Now get the actual recommendations
	query := `
		SELECT * FROM ai_recommendations
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	var recommendations []*models.AIRecommendation
	err = r.db.SelectContext(ctx, &recommendations, query, userID, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get recommendations")
		return nil, 0, fmt.Errorf("failed to get recommendations: %w", err)
	}

	return recommendations, count, nil
}

// UpdateRecommendationSelection updates the selection status of a recommendation
func (r *sqlAIConversationRepository) UpdateRecommendationSelection(ctx context.Context, recommendationID string, selected bool) error {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Str("recommendationID", recommendationID).
		Bool("selected", selected).
		Msg("Updating recommendation selection")

	var query string
	var args []interface{}

	if selected {
		query = `
			UPDATE ai_recommendations
			SET selected = TRUE, selected_at = NOW()
			WHERE id = ?
		`
		args = []interface{}{recommendationID}
	} else {
		query = `
			UPDATE ai_recommendations
			SET selected = FALSE, selected_at = NULL
			WHERE id = ?
		`
		args = []interface{}{recommendationID}
	}

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update recommendation selection")
		return fmt.Errorf("failed to update recommendation selection: %w", err)
	}

	return nil
}

// GetConversationAnalytics retrieves analytics for a conversation
func (r *sqlAIConversationRepository) GetConversationAnalytics(ctx context.Context, conversationID string) (*models.AIConversationAnalytics, error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().Str("conversationID", conversationID).Msg("Getting conversation analytics")

	query := `
		SELECT * FROM ai_conversation_analytics
		WHERE conversation_id = ?
	`

	var analytics models.AIConversationAnalytics
	err := r.db.GetContext(ctx, &analytics, query, conversationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// If no analytics exist yet, create a new one
			return models.NewAIConversationAnalytics(conversationID), nil
		}
		log.Error().Err(err).Msg("Failed to get conversation analytics")
		return nil, fmt.Errorf("failed to get conversation analytics: %w", err)
	}

	return &analytics, nil
}

// UpdateConversationAnalytics updates analytics for a conversation
func (r *sqlAIConversationRepository) UpdateConversationAnalytics(ctx context.Context, analytics *models.AIConversationAnalytics) error {
	log := logger.LoggerFromContext(ctx)
	log.Info().Str("conversationID", analytics.ConversationID).Msg("Updating conversation analytics")

	// Check if analytics already exist
	checkQuery := `SELECT COUNT(*) FROM ai_conversation_analytics WHERE conversation_id = ?`
	var count int
	err := r.db.GetContext(ctx, &count, checkQuery, analytics.ConversationID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check if analytics exist")
		return fmt.Errorf("failed to check if analytics exist: %w", err)
	}

	var query string
	if count > 0 {
		// Update existing analytics
		query = `
			UPDATE ai_conversation_analytics
			SET 
				total_user_messages = :total_user_messages,
				total_assistant_messages = :total_assistant_messages,
				total_recommendations = :total_recommendations,
				selected_recommendations = :selected_recommendations,
				total_tokens_used = :total_tokens_used,
				average_response_time = :average_response_time,
				conversation_duration = :conversation_duration,
				last_updated_at = NOW()
			WHERE conversation_id = :conversation_id
		`
	} else {
		// Insert new analytics
		query = `
			INSERT INTO ai_conversation_analytics (
				conversation_id, total_user_messages, total_assistant_messages,
				total_recommendations, selected_recommendations, total_tokens_used,
				average_response_time, conversation_duration, last_updated_at
			) VALUES (
				:conversation_id, :total_user_messages, :total_assistant_messages,
				:total_recommendations, :selected_recommendations, :total_tokens_used,
				:average_response_time, :conversation_duration, NOW()
			)
		`
	}

	_, err = r.db.NamedExecContext(ctx, query, analytics)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update conversation analytics")
		return fmt.Errorf("failed to update conversation analytics: %w", err)
	}

	return nil
}

// GetUserConversationStats retrieves aggregated statistics for a user's AI interactions
func (r *sqlAIConversationRepository) GetUserConversationStats(ctx context.Context, userID uint64) (*models.UserConversationStats, error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().Uint64("userID", userID).Msg("Getting user conversation stats")

	query := `
		SELECT
			? AS user_id,
			COUNT(DISTINCT c.id) AS total_conversations,
			SUM(c.message_count) AS total_messages,
			(SELECT COUNT(*) FROM ai_recommendations r WHERE r.user_id = ?) AS total_recommendations,
			(SELECT COUNT(*) FROM ai_recommendations r WHERE r.user_id = ? AND r.selected = TRUE) AS selected_recommendations,
			IFNULL(AVG(c.message_count), 0) AS average_messages_per_conv,
			(
				SELECT content_type
				FROM ai_conversations
				WHERE user_id = ?
				GROUP BY content_type
				ORDER BY COUNT(*) DESC
				LIMIT 1
			) AS favorite_content_type,
			IFNULL(MAX(c.updated_at), NOW()) AS last_conversation
		FROM ai_conversations c
		WHERE c.user_id = ?
	`

	var stats models.UserConversationStats
	err := r.db.GetContext(ctx, &stats, query, userID, userID, userID, userID, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user conversation stats")
		return nil, fmt.Errorf("failed to get user conversation stats: %w", err)
	}

	return &stats, nil
}