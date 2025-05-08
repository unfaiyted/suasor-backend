package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"suasor/clients"
	"suasor/clients/ai"
	clienttypes "suasor/clients/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/utils"
	"suasor/utils/logger"
	"time"
)

// AIConversationService defines the interface for AI conversation related operations
type AIConversationService interface {
	// Core conversation methods
	StartConversation(ctx context.Context, userID uint64, clientID uint64, contentType string,
		preferences map[string]any, systemInstructions string) (string, string, error)
	SendMessage(ctx context.Context, conversationID string, userID uint64, message string,
		context map[string]any) (string, []map[string]any, error)
	GetConversationHistory(ctx context.Context, conversationID string, userID uint64) ([]*models.AIMessage, error)

	// User history methods
	GetUserConversations(ctx context.Context, userID uint64, limit, offset int) ([]*models.AIConversation, int, error)
	GetUserRecommendationHistory(ctx context.Context, userID uint64, itemType string, limit, offset int) ([]*models.AIRecommendation, int, error)

	// Management methods
	ArchiveConversation(ctx context.Context, conversationID string, userID uint64) error
	DeleteConversation(ctx context.Context, conversationID string, userID uint64) error

	// Analytics methods
	GetConversationInsights(ctx context.Context, conversationID string, userID uint64) (*models.ConversationInsights, error)
	GetUserAIInteractionSummary(ctx context.Context, userID uint64) (*models.AIInteractionSummary, error)
}

// aiConversationService implements the AIConversationService interface
type aiConversationService struct {
	repo                repository.AIConversationRepository
	claudeClientService ClientService[*clienttypes.ClaudeConfig]
	openaiClientService ClientService[*clienttypes.OpenAIConfig]
	ollamaClientService ClientService[*clienttypes.OllamaConfig]
	clientHelper        repository.ClientHelper
	clientFactory       *clients.ClientProviderFactoryService
	activeClients       map[string]ai.ClientAI // conversationID -> client
}

// NewAIConversationService creates a new AI conversation service
func NewAIConversationService(
	repo repository.AIConversationRepository,
	claudeClientService ClientService[*clienttypes.ClaudeConfig],
	openaiClientService ClientService[*clienttypes.OpenAIConfig],
	ollamaClientService ClientService[*clienttypes.OllamaConfig],
	clientHelper repository.ClientHelper,
	clientFactory *clients.ClientProviderFactoryService,
) AIConversationService {
	return &aiConversationService{
		repo:                repo,
		claudeClientService: claudeClientService,
		openaiClientService: openaiClientService,
		ollamaClientService: ollamaClientService,
		clientHelper:        clientHelper,
		clientFactory:       clientFactory,
		activeClients:       make(map[string]ai.ClientAI),
	}
}

// StartConversation begins a new AI conversation
func (s *aiConversationService) StartConversation(
	ctx context.Context,
	userID uint64,
	clientID uint64,
	contentType string,
	preferences map[string]any,
	systemInstructions string,
) (string, string, error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("contentType", contentType).
		Msg("Starting new AI conversation")

	// Get AI client
	aiClient, err := s.getAIClient(ctx, clientID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get AI client")
		return "", "", err
	}

	// Start conversation with the AI provider
	conversationID, welcomeMessage, err := aiClient.StartRecommendationConversation(
		ctx,
		contentType,
		preferences,
		systemInstructions,
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to start AI conversation")
		return "", "", err
	}

	// Store conversation in database
	conversation := models.NewAIConversation(
		conversationID,
		userID,
		clientID,
		contentType,
		systemInstructions,
	)

	// Store preferences
	if err := conversation.SetUserPreferencesMap(preferences); err != nil {
		log.Error().Err(err).Msg("Failed to set user preferences")
		return "", "", err
	}

	// Save to database
	_, err = s.repo.CreateConversation(ctx, conversation)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create conversation in database")
		return "", "", err
	}

	// Save welcome message from AI
	welcomeMsg := models.NewAIMessage(
		utils.GenerateRandomID(16),
		conversationID,
		"assistant", // Role is assistant for the welcome message
		welcomeMessage,
	)
	_, err = s.repo.AddMessage(ctx, welcomeMsg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to save welcome message")
		// Continue even if saving the message fails
		log.Warn().Msg("Continuing despite message save failure")
	}

	// Store AI client in active clients map
	s.activeClients[conversationID] = aiClient

	return conversationID, welcomeMessage, nil
}

// SendMessage sends a message in an existing conversation
func (s *aiConversationService) SendMessage(
	ctx context.Context,
	conversationID string,
	userID uint64,
	message string,
	messageContext map[string]any,
) (string, []map[string]any, error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Str("conversationID", conversationID).
		Uint64("userID", userID).
		Str("message", message).
		Msg("Sending message in AI conversation")

	// Get the conversation
	conversation, err := s.repo.GetConversationByID(ctx, conversationID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get conversation")
		return "", nil, err
	}

	if conversation == nil {
		return "", nil, errors.New("conversation not found")
	}

	// Verify user owns the conversation
	if conversation.UserID != userID {
		return "", nil, errors.New("user does not own this conversation")
	}

	// Get the AI client
	aiClient, err := s.getClientForConversation(ctx, conversationID, conversation.ClientID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get AI client for conversation")
		return "", nil, err
	}

	// Save user message to database
	userMsg := models.NewAIMessage(
		utils.GenerateRandomID(16),
		conversationID,
		"user",
		message,
	)
	_, err = s.repo.AddMessage(ctx, userMsg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to save user message")
		// Continue despite the error
		log.Warn().Msg("Continuing despite message save failure")
	}

	// Ensure extract recommendations is set
	if messageContext == nil {
		messageContext = make(map[string]any)
	}
	if _, exists := messageContext["extractRecommendations"]; !exists {
		messageContext["extractRecommendations"] = true
	}

	// Send message to AI client
	startTime := time.Now()
	aiResponse, recommendations, err := aiClient.ContinueRecommendationConversation(
		ctx,
		conversationID,
		message,
		messageContext,
	)
	responseDuration := time.Since(startTime)
	if err != nil {
		log.Error().Err(err).Msg("Failed to continue conversation with AI")
		return "", nil, err
	}

	// Save AI response to database
	aiMsg := models.NewAIMessage(
		utils.GenerateRandomID(16),
		conversationID,
		"assistant",
		aiResponse,
	)

	// Add response time and other metadata
	metadata := map[string]any{
		"responseTime": responseDuration.Milliseconds(),
	}
	if err := aiMsg.SetMetadataMap(metadata); err != nil {
		log.Warn().Err(err).Msg("Failed to set message metadata")
	}

	msgID, err := s.repo.AddMessage(ctx, aiMsg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to save AI response")
		// Continue despite the error
		log.Warn().Msg("Continuing despite message save failure")
	}

	// Save recommendations if any
	if len(recommendations) > 0 {
		for _, rec := range recommendations {
			// Extract basic info
			title, _ := rec["title"].(string)
			if title == "" {
				title = "Unknown"
			}

			itemType := conversation.ContentType
			if specificType, ok := rec["type"].(string); ok && specificType != "" {
				itemType = specificType
			}

			// Create recommendation object
			recommendation := models.NewAIRecommendation(
				utils.GenerateRandomID(16),
				msgID,
				conversationID,
				userID,
				itemType,
				title,
			)

			// Save the full data
			if err := recommendation.SetDataMap(rec); err != nil {
				log.Warn().Err(err).Msg("Failed to set recommendation data")
				continue
			}

			// Set reason if available
			if reason, ok := rec["reason"].(string); ok && reason != "" {
				recommendation.Reason = sql.NullString{
					String: reason,
					Valid:  true,
				}
			}

			// Set external ID if available
			if externalID, ok := rec["externalId"].(string); ok && externalID != "" {
				recommendation.ExternalID = sql.NullString{
					String: externalID,
					Valid:  true,
				}
			}

			// Save to database
			_, err := s.repo.AddRecommendation(ctx, recommendation)
			if err != nil {
				log.Error().Err(err).Msg("Failed to save recommendation")
				// Continue with next recommendation despite the error
			}
		}
	}

	// Update analytics
	s.updateConversationAnalytics(ctx, conversationID, "user", responseDuration)

	return aiResponse, recommendations, nil
}

// GetConversationHistory retrieves the message history for a conversation
func (s *aiConversationService) GetConversationHistory(
	ctx context.Context,
	conversationID string,
	userID uint64,
) ([]*models.AIMessage, error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Str("conversationID", conversationID).
		Uint64("userID", userID).
		Msg("Getting conversation history")

	// Get the conversation
	conversation, err := s.repo.GetConversationByID(ctx, conversationID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get conversation")
		return nil, err
	}

	if conversation == nil {
		return nil, errors.New("conversation not found")
	}

	// Verify user owns the conversation
	if conversation.UserID != userID {
		return nil, errors.New("user does not own this conversation")
	}

	// Get messages
	messages, err := s.repo.GetConversationHistory(ctx, conversationID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get conversation history")
		return nil, err
	}

	return messages, nil
}

// GetUserConversations retrieves conversations for a user with pagination
func (s *aiConversationService) GetUserConversations(
	ctx context.Context,
	userID uint64,
	limit, offset int,
) ([]*models.AIConversation, int, error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", userID).
		Int("limit", limit).
		Int("offset", offset).
		Msg("Getting user conversations")

	return s.repo.GetConversationsByUserID(ctx, userID, limit, offset)
}

// GetUserRecommendationHistory retrieves recommendations for a user with pagination
func (s *aiConversationService) GetUserRecommendationHistory(
	ctx context.Context,
	userID uint64,
	itemType string,
	limit, offset int,
) ([]*models.AIRecommendation, int, error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", userID).
		Str("itemType", itemType).
		Int("limit", limit).
		Int("offset", offset).
		Msg("Getting user recommendation history")

	// If itemType is specified, filter by it
	if itemType != "" {
		// This would require a custom repository method
		// For now, we'll just get all recommendations and filter in memory
		recs, count, err := s.repo.GetRecommendationsByUserID(ctx, userID, limit, offset)
		if err != nil {
			return nil, 0, err
		}

		// Filter by itemType
		var filtered []*models.AIRecommendation
		for _, rec := range recs {
			if rec.ItemType == itemType {
				filtered = append(filtered, rec)
			}
		}

		return filtered, count, nil
	}

	return s.repo.GetRecommendationsByUserID(ctx, userID, limit, offset)
}

// ArchiveConversation archives a conversation
func (s *aiConversationService) ArchiveConversation(
	ctx context.Context,
	conversationID string,
	userID uint64,
) error {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Str("conversationID", conversationID).
		Uint64("userID", userID).
		Msg("Archiving conversation")

	// Get the conversation
	conversation, err := s.repo.GetConversationByID(ctx, conversationID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get conversation")
		return err
	}

	if conversation == nil {
		return errors.New("conversation not found")
	}

	// Verify user owns the conversation
	if conversation.UserID != userID {
		return errors.New("user does not own this conversation")
	}

	// Update status
	return s.repo.UpdateConversationStatus(ctx, conversationID, "archived")
}

// DeleteConversation deletes a conversation and all related data
func (s *aiConversationService) DeleteConversation(
	ctx context.Context,
	conversationID string,
	userID uint64,
) error {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Str("conversationID", conversationID).
		Uint64("userID", userID).
		Msg("Deleting conversation")

	// Get the conversation
	conversation, err := s.repo.GetConversationByID(ctx, conversationID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get conversation")
		return err
	}

	if conversation == nil {
		return errors.New("conversation not found")
	}

	// Verify user owns the conversation
	if conversation.UserID != userID {
		return errors.New("user does not own this conversation")
	}

	// Remove from active clients if present
	delete(s.activeClients, conversationID)

	// Delete from database
	return s.repo.DeleteConversation(ctx, conversationID)
}

// GetConversationInsights provides detailed analytics for a specific conversation
func (s *aiConversationService) GetConversationInsights(
	ctx context.Context,
	conversationID string,
	userID uint64,
) (*models.ConversationInsights, error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Str("conversationID", conversationID).
		Uint64("userID", userID).
		Msg("Getting conversation insights")

	// Get the conversation
	conversation, err := s.repo.GetConversationByID(ctx, conversationID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get conversation")
		return nil, err
	}

	if conversation == nil {
		return nil, errors.New("conversation not found")
	}

	// Verify user owns the conversation
	if conversation.UserID != userID {
		return nil, errors.New("user does not own this conversation")
	}

	// Get analytics
	analytics, err := s.repo.GetConversationAnalytics(ctx, conversationID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get conversation analytics")
		return nil, err
	}

	// Get messages to calculate timing data
	messages, err := s.repo.GetConversationHistory(ctx, conversationID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get conversation history")
		return nil, err
	}

	// Calculate message summary
	messageSummary := make(map[string]int)
	messageSummary["user"] = 0
	messageSummary["assistant"] = 0

	for _, msg := range messages {
		messageSummary[msg.Role]++
	}

	// Calculate timing data
	timingData := calculateTimingData(messages)

	// Get recommendations
	recommendations, err := s.repo.GetRecommendationsByConversationID(ctx, conversationID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get recommendations")
		return nil, err
	}

	// Calculate recommendation effectiveness
	var recommendationEffectiveness float64 = 0
	if len(recommendations) > 0 {
		selectedCount := 0
		for _, rec := range recommendations {
			if rec.Selected {
				selectedCount++
			}
		}
		recommendationEffectiveness = float64(selectedCount) / float64(len(recommendations))
	}

	// Build insights
	insights := &models.ConversationInsights{
		Conversation:                conversation,
		Analytics:                   analytics,
		MessageSummary:              messageSummary,
		TimingData:                  timingData,
		RecommendationEffectiveness: recommendationEffectiveness,
	}

	return insights, nil
}

// GetUserAIInteractionSummary provides a high-level view of user's AI usage
func (s *aiConversationService) GetUserAIInteractionSummary(
	ctx context.Context,
	userID uint64,
) (*models.AIInteractionSummary, error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", userID).
		Msg("Getting user AI interaction summary")

	// Get user conversation stats
	stats, err := s.repo.GetUserConversationStats(ctx, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user conversation stats")
		return nil, err
	}

	// Get recent conversations
	recentConversations, _, err := s.repo.GetConversationsByUserID(ctx, userID, 5, 0)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get recent conversations")
		return nil, err
	}

	// Get top recommendations (selected ones)
	// This would require a custom repository method
	// For now, we'll just get most recent recommendations
	topRecommendations, _, err := s.repo.GetRecommendationsByUserID(ctx, userID, 5, 0)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get top recommendations")
		return nil, err
	}

	// Get content type breakdown
	// This would require a custom repository method
	// For now, just use an empty map
	contentTypeBreakdown := make(map[string]int)

	// Build summary
	summary := &models.AIInteractionSummary{
		UserID:               userID,
		Stats:                *stats,
		RecentConversations:  recentConversations,
		TopRecommendations:   topRecommendations,
		ContentTypeBreakdown: contentTypeBreakdown,
	}

	return summary, nil
}

// Helper methods

// getAIClient retrieves or creates an AI client
func (s *aiConversationService) getAIClient(ctx context.Context, clientID uint64) (ai.ClientAI, error) {
	log := logger.LoggerFromContext(ctx)

	clientType, err := s.clientHelper.GetClientTypeByClientID(ctx, clientID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get AI client")
		return nil, err
	}

	var aiConfig clienttypes.ClientConfig
	if clientType == clienttypes.ClientTypeOpenAI {
		clientService, err := s.openaiClientService.GetByID(ctx, clientID)
		aiConfig = clientService.GetConfig()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get AI client")
			return nil, err
		}
	} else if clientType == clienttypes.ClientTypeOllama {
		clientService, err := s.ollamaClientService.GetByID(ctx, clientID)
		aiConfig = clientService.GetConfig()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get AI client")
			return nil, err
		}
	} else if clientType == clienttypes.ClientTypeClaude {
		clientService, err := s.claudeClientService.GetByID(ctx, clientID)
		aiConfig = clientService.GetConfig()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get AI client")
			return nil, err
		}
	}

	// Get client from factory using the actual config
	client, err := s.clientFactory.GetClient(ctx, clientID, aiConfig)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get AI client from factory")
		return nil, err
	}

	// Convert to AI client
	aiClient, ok := client.(ai.ClientAI)
	if !ok {
		return nil, fmt.Errorf("client is not an AI client")
	}

	return aiClient, nil
}

// getClientForConversation retrieves the client for a specific conversation
func (s *aiConversationService) getClientForConversation(ctx context.Context, conversationID string, clientID uint64) (ai.ClientAI, error) {
	// Check if we have an active client
	if client, exists := s.activeClients[conversationID]; exists {
		return client, nil
	}

	// Otherwise get a new client
	client, err := s.getAIClient(ctx, clientID)
	if err != nil {
		return nil, err
	}

	// Store for future use
	s.activeClients[conversationID] = client
	return client, nil
}

// updateConversationAnalytics updates analytics for a conversation
func (s *aiConversationService) updateConversationAnalytics(
	ctx context.Context,
	conversationID string,
	messageRole string,
	responseDuration time.Duration,
) {
	log := logger.LoggerFromContext(ctx)

	// Get existing analytics or create new ones
	analytics, err := s.repo.GetConversationAnalytics(ctx, conversationID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get conversation analytics")
		return
	}

	// Update based on message role
	if messageRole == "user" {
		analytics.TotalUserMessages++
	} else if messageRole == "assistant" {
		analytics.TotalAssistantMessages++
	}

	// Update response time if this is an assistant message
	if messageRole == "assistant" && responseDuration > 0 {
		// Calculate new average
		totalTime := analytics.AverageResponseTime * float64(analytics.TotalAssistantMessages-1)
		totalTime += float64(responseDuration.Milliseconds())
		analytics.AverageResponseTime = totalTime / float64(analytics.TotalAssistantMessages)
	}

	// Update last updated
	analytics.LastUpdatedAt = time.Now()

	// Save analytics
	err = s.repo.UpdateConversationAnalytics(ctx, analytics)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update conversation analytics")
	}
}

// calculateTimingData analyzes the message timing patterns
func calculateTimingData(messages []*models.AIMessage) map[string]any {
	if len(messages) < 2 {
		return map[string]any{
			"averageUserResponseTime": 0,
			"averageAIResponseTime":   0,
			"totalConversationTime":   0,
		}
	}

	var userResponseTimes []int64
	var aiResponseTimes []int64

	for i := 1; i < len(messages); i++ {
		currentMsg := messages[i]
		prevMsg := messages[i-1]

		// Time difference in milliseconds
		timeDiff := currentMsg.Timestamp.Sub(prevMsg.Timestamp).Milliseconds()

		if currentMsg.Role == "user" {
			userResponseTimes = append(userResponseTimes, timeDiff)
		} else {
			aiResponseTimes = append(aiResponseTimes, timeDiff)
		}
	}

	// Calculate averages
	var avgUserTime int64 = 0
	if len(userResponseTimes) > 0 {
		var sum int64 = 0
		for _, t := range userResponseTimes {
			sum += t
		}
		avgUserTime = sum / int64(len(userResponseTimes))
	}

	var avgAITime int64 = 0
	if len(aiResponseTimes) > 0 {
		var sum int64 = 0
		for _, t := range aiResponseTimes {
			sum += t
		}
		avgAITime = sum / int64(len(aiResponseTimes))
	}

	// Total conversation time
	totalTime := int64(0)
	if len(messages) >= 2 {
		first := messages[0]
		last := messages[len(messages)-1]
		totalTime = last.Timestamp.Sub(first.Timestamp).Milliseconds()
	}

	return map[string]any{
		"averageUserResponseTime": avgUserTime,
		"averageAIResponseTime":   avgAITime,
		"totalConversationTime":   totalTime,
		"userResponseTimes":       userResponseTimes,
		"aiResponseTimes":         aiResponseTimes,
	}
}
