package handlers

import (
	"context"
	"fmt"
	"suasor/clients"
	"suasor/clients/types"
	"suasor/services"
	"suasor/types/requests"
	"suasor/types/responses"
	"suasor/utils/logger"

	"github.com/gin-gonic/gin"
)

// AIHandler implements AI-related handlers with support for multiple AI client types
type AIHandler[T types.AIClientConfig] struct {
	factory *clients.ClientProviderFactoryService
	service services.ClientService[T]
	// Map to track active conversations by conversationID
	activeConversations map[string]uint64 // conversationID -> userID
}

// NewAIHandler creates a new AI handler
func NewAIHandler[T types.AIClientConfig](
	factory *clients.ClientProviderFactoryService,
	service services.ClientService[T],
) *AIHandler[T] {
	return &AIHandler[T]{
		factory:             factory,
		service:             service,
		activeConversations: make(map[string]uint64),
	}
}

// RequestRecommendation godoc
//
//	@Summary		Get AI-powered content recommendations
//	@Description	Get content recommendations from an AI service
//	@Tags			ai
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		requests.AiRecommendationRequest							true	"Recommendation request"
//	@Success		200		{object}	responses.APIResponse[responses.AiRecommendationResponse]	"Recommendation response"
//	@Failure		400		{object}	responses.ErrorResponse[responses.ErrorDetails]				"Invalid request"
//	@Failure		401		{object}	responses.ErrorResponse[responses.ErrorDetails]				"Unauthorized"
//	@Failure		500		{object}	responses.ErrorResponse[responses.ErrorDetails]				"Server error"
//	@Router			/ai/recommendations [post]
func (h *AIHandler[T]) RequestRecommendation(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	var req requests.AiRecommendationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.RespondValidationError(c, err)
		return
	}

	log.Info().
		Uint64("userID", userID.(uint64)).
		Str("contentType", req.ContentType).
		Interface("filters", req.Filters).
		Msg("Requesting AI recommendations")

	clientType := types.ClientType(c.Param("clientType"))

	// Get available AI client based on specified type or default
	aiClient, err := h.getAIClient(ctx, userID.(uint64), clientType, req.ClientID)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to initialize AI client")
		return
	}

	// Get recommendations
	recommendations, err := aiClient.GetRecommendations(ctx, req.ContentType, req.Filters, req.Count)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to get recommendations")
		return
	}

	response := responses.AiRecommendationResponse{
		Items: recommendations,
	}

	responses.RespondOK(c, response, "Recommendations retrieved successfully")
}

// AnalyzeContent godoc
//
//	@Summary		Analyze content with AI
//	@Description	Use AI to analyze provided content
//	@Tags			ai
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		requests.AiContentAnalysisRequest							true	"Content analysis request"
//	@Success		200		{object}	responses.APIResponse[responses.AiContentAnalysisResponse]	"Analysis response"
//	@Failure		400		{object}	responses.ErrorResponse[responses.ErrorDetails]				"Invalid request"
//	@Failure		401		{object}	responses.ErrorResponse[responses.ErrorDetails]				"Unauthorized"
//	@Failure		500		{object}	responses.ErrorResponse[responses.ErrorDetails]				"Server error"
//	@Router			/ai/analyze [post]
func (h *AIHandler[T]) AnalyzeContent(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	var req requests.AiContentAnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.RespondValidationError(c, err)
		return
	}

	log.Info().
		Uint64("userID", userID.(uint64)).
		Str("contentType", req.ContentType).
		Int("contentLength", len(req.Content)).
		Msg("Requesting AI content analysis")

	clientType := types.ClientType(c.Param("clientType"))

	// Get available AI client based on specified type or default
	aiClient, err := h.getAIClient(ctx, userID.(uint64), clientType, req.ClientID)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to initialize AI client")
		return
	}

	// Analyze content
	analysis, err := aiClient.AnalyzeContent(ctx, req.ContentType, req.Content, req.Options)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to analyze content")
		return
	}

	response := responses.AiContentAnalysisResponse{
		Analysis: analysis,
	}

	responses.RespondOK(c, response, "Content analyzed successfully")
}

// StartConversation godoc
//
//	@Summary		Start a new AI conversation for recommendations
//	@Description	Initialize a new conversational session with the AI for personalized recommendations
//	@Tags			ai
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		requests.StartConversationRequest						true	"Conversation initialization request"
//	@Success		200		{object}	responses.APIResponse[responses.ConversationResponse]	"Conversation started"
//	@Failure		400		{object}	responses.ErrorResponse[responses.ErrorDetails]			"Invalid request"
//	@Failure		401		{object}	responses.ErrorResponse[responses.ErrorDetails]			"Unauthorized"
//	@Failure		500		{object}	responses.ErrorResponse[responses.ErrorDetails]			"Server error"
//	@Router			/ai/conversation/start [post]
func (h *AIHandler[T]) StartConversation(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	var req requests.StartConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.RespondValidationError(c, err)
		return
	}
	clientType := types.ClientType(c.Param("clientType"))

	log.Info().
		Uint64("userID", userID.(uint64)).
		Str("contentType", req.ContentType).
		Uint64("clientID", req.ClientID).
		Str("clientType", clientType.String()).
		Msg("Starting AI recommendation conversation")

	// Get available AI client based on specified type or default
	aiClient, err := h.getAIClient(ctx, userID.(uint64), clientType, req.ClientID)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to initialize AI client")
		return
	}

	// Start the conversation
	conversationID, welcomeMessage, err := aiClient.StartRecommendationConversation(
		ctx,
		req.ContentType,
		req.Preferences,
		req.SystemInstructions,
	)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to start recommendation conversation")
		return
	}

	// Save the conversation in our tracking map
	h.activeConversations[conversationID] = userID.(uint64)

	// Prepare and send response
	response := responses.ConversationResponse{
		ConversationID: conversationID,
		Welcome:        welcomeMessage,
		Context: map[string]interface{}{
			"contentType": req.ContentType,
		},
	}

	responses.RespondOK(c, response, "Conversation started successfully")
}

func (h *AIHandler[T]) getAIClient(ctx context.Context, userID uint64, clientType types.ClientType, clientID uint64) (types.AiClient, error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("userID", userID).
		Str("clientType", clientType.String()).
		Uint64("clientID", clientID).
		Msg("Retrieving client")
	// Get all AI clients for the user
	clientModel, err := h.service.GetByID(ctx, clientID, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get AI client")
		return nil, err
	}
	log.Info().
		Uint64("userID", userID).
		Str("clientType", clientType.String()).
		Uint64("clientID", clientID).
		Msg("Client Model retrieved")
	// from factory
	aiClient, err := h.factory.GetClient(ctx, clientModel.ID, clientModel.Config.Data)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get AI client")
		return nil, err
	}
	aiAiClient, ok := aiClient.(types.AiClient)
	if !ok {
		return nil, fmt.Errorf("client is not an AI client")
	}

	return aiAiClient, nil

}

// SendConversationMessage godoc
//
//	@Summary		Send a message in an existing AI conversation
//	@Description	Continue a conversation with the AI by sending a new message
//	@Tags			ai
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		requests.ConversationMessageRequest								true	"Message request"
//	@Success		200		{object}	responses.APIResponse[responses.ConversationMessageResponse]	"AI response"
//	@Failure		400		{object}	responses.ErrorResponse[responses.ErrorDetails]					"Invalid request"
//	@Failure		401		{object}	responses.ErrorResponse[responses.ErrorDetails]					"Unauthorized"
//	@Failure		403		{object}	responses.ErrorResponse[responses.ErrorDetails]					"Conversation not owned by user"
//	@Failure		404		{object}	responses.ErrorResponse[responses.ErrorDetails]					"Conversation not found"
//	@Failure		500		{object}	responses.ErrorResponse[responses.ErrorDetails]					"Server error"
//	@Router			/ai/conversation/message [post]
func (h *AIHandler[T]) SendConversationMessage(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	var req requests.ConversationMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.RespondValidationError(c, err)
		return
	}

	log.Info().
		Uint64("userID", userID.(uint64)).
		Str("conversationID", req.ConversationID).
		Str("message", req.Message).
		Msg("Sending message in AI conversation")

	// Verify conversation exists and belongs to this user
	conversationOwnerID, exists := h.activeConversations[req.ConversationID]
	if !exists {
		responses.RespondNotFound(c, nil, "Conversation not found")
		return
	}

	if conversationOwnerID != userID.(uint64) {
		responses.RespondForbidden(c, nil, "You do not have access to this conversation")
		return
	}
	clientType := types.ClientType(c.Param("clientType"))

	// Get available AI client
	// Note: We don't need to specify a client type here as the conversation is already
	// associated with a specific AI client from the start conversation request
	aiClient, err := h.getAIClient(ctx, userID.(uint64), clientType, req.ClientID)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to initialize AI client")
		return
	}

	// Set extractRecommendations to true by default if not specified
	context := req.Context
	if context == nil {
		context = make(map[string]interface{})
	}
	if _, exists := context["extractRecommendations"]; !exists {
		context["extractRecommendations"] = true
	}

	// Continue the conversation
	message, recommendations, err := aiClient.ContinueRecommendationConversation(
		ctx,
		req.ConversationID,
		req.Message,
		context,
	)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to continue conversation")
		return
	}

	// Prepare and send response
	response := responses.ConversationMessageResponse{
		Message:         message,
		Recommendations: recommendations,
		Context:         make(map[string]interface{}),
	}

	responses.RespondOK(c, response, "Message sent successfully")
}
