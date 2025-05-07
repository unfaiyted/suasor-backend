package handlers

import (
	"context"
	"fmt"
	"suasor/clients"
	"suasor/clients/ai"
	aitypes "suasor/clients/ai/types"
	"suasor/clients/types"
	"suasor/services"
	"suasor/types/requests"
	"suasor/types/responses"
	"suasor/utils/logger"

	"github.com/gin-gonic/gin"
)

type AIHandler[T types.AIClientConfig] interface {
	RequestRecommendation(c *gin.Context)
	AnalyzeContent(c *gin.Context)
	StartConversation(c *gin.Context)
	SendConversationMessage(c *gin.Context)
}

// aiHandler implements AI-related handlers with support for multiple AI client types
type aiHandler[T types.AIClientConfig] struct {
	factory             *clients.ClientProviderFactoryService
	service             services.ClientService[T]
	// Conversation service for persistent storage
	conversationService services.AIConversationService
	// Map to track active conversations by conversationID
	activeConversations map[string]uint64 // conversationID -> userID
}

// NewaiHandler creates a new AI handler
func NewAIHandler[T types.AIClientConfig](
	factory *clients.ClientProviderFactoryService,
	service services.ClientService[T],
	conversationService services.AIConversationService,
) AIHandler[T] {
	return &aiHandler[T]{
		factory:             factory,
		service:             service,
		conversationService: conversationService,
		activeConversations: make(map[string]uint64),
	}
}

// RequestRecommendation godoc
//
//		@Summary		Get AI-powered content recommendations
//		@Description	Get content recommendations from an AI service
//		@Tags			ai
//		@Accept			json
//		@Produce		json
//		@Security		BearerAuth
//	  @Param			clientID	path		int															true	"Client ID"
//		@Param			request	body		aitypes.RecommendationRequest							true	"Recommendation request"
//		@Success		200		{object}	responses.APIResponse[responses.AiRecommendationResponse]	"Recommendation response"
//		@Failure		400		{object}	responses.ErrorResponse[responses.ErrorDetails]				"Invalid request"
//		@Failure		401		{object}	responses.ErrorResponse[responses.ErrorDetails]				"Unauthorized"
//		@Failure		500		{object}	responses.ErrorResponse[responses.ErrorDetails]				"Server error"
//		@Router			/client/{clientID}/ai/recommendations [post]
func (h *aiHandler[T]) RequestRecommendation(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, _ := checkUserAccess(c)
	clientID, _ := checkItemID(c, "clientID")

	var req aitypes.RecommendationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.RespondValidationError(c, err)
		return
	}

	log.Info().
		Uint64("userID", userID).
		Msg("Requesting AI recommendations")

	clientType := types.ClientType(c.Param("clientType"))

	// Get available AI client based on specified type or default
	aiClient, err := h.getAIClient(ctx, userID, clientType, clientID)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to initialize AI client")
		return
	}

	// Get recommendations
	aiRecommendationResponse, err := aiClient.GetRecommendations(ctx, &req)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to get recommendations")
		return
	}

	responses.RespondOK(c, aiRecommendationResponse, "Recommendations retrieved successfully")
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
//	@Param			clientID	path		int															true	"Client ID"
//	@Success		200		{object}	responses.APIResponse[responses.AiContentAnalysisResponse]	"Analysis response"
//	@Failure		400		{object}	responses.ErrorResponse[responses.ErrorDetails]				"Invalid request"
//	@Failure		401		{object}	responses.ErrorResponse[responses.ErrorDetails]				"Unauthorized"
//	@Failure		500		{object}	responses.ErrorResponse[responses.ErrorDetails]				"Server error"
//	@Router			/client/{clientID}/ai/analyze [post]
func (h *aiHandler[T]) AnalyzeContent(c *gin.Context) {
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
//	@Router			/client/{clientID}/ai/conversation/start [post]
func (h *aiHandler[T]) StartConversation(c *gin.Context) {
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

	// Start the conversation with the AI client
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

	// Save the conversation both in our tracking map and in the database
	h.activeConversations[conversationID] = userID.(uint64)
	
	// Store it in the database through the conversation service
	_, err = h.conversationService.StartConversation(
		ctx,
		userID.(uint64),
		req.ClientID,
		req.ContentType,
		req.Preferences,
		req.SystemInstructions,
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to persist conversation to database, continuing anyway")
		// Not returning error as we can still continue with in-memory conversation
	}

	// Prepare and send response
	response := responses.ConversationResponse{
		ConversationID: conversationID,
		Welcome:        welcomeMessage,
		Context: map[string]any{
			"contentType": req.ContentType,
		},
	}

	responses.RespondOK(c, response, "Conversation started successfully")
}

func (h *aiHandler[T]) getAIClient(ctx context.Context, userID uint64, clientType types.ClientType, clientID uint64) (ai.ClientAI, error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("userID", userID).
		Str("clientType", clientType.String()).
		Uint64("clientID", clientID).
		Msg("Retrieving client")
	// Get all AI clients for the user
	clientModel, err := h.service.GetByID(ctx, clientID)
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
	aiClient, err := h.factory.GetClient(ctx, clientModel.ID, clientModel.Config)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get AI client")
		return nil, err
	}
	aiAiClient, ok := aiClient.(ai.ClientAI)
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
//	@Router			/client/{clientID}/ai/conversation/message [post]
func (h *aiHandler[T]) SendConversationMessage(c *gin.Context) {
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
		// Check if it exists in the database
		conversation, err := h.conversationService.GetConversationHistory(ctx, req.ConversationID, userID.(uint64))
		if err != nil || len(conversation) == 0 {
			responses.RespondNotFound(c, nil, "Conversation not found")
			return
		}
		// If it exists in the database but not in memory, add it to memory
		h.activeConversations[req.ConversationID] = userID.(uint64)
		conversationOwnerID = userID.(uint64)
	} else if conversationOwnerID != userID.(uint64) {
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
		context = make(map[string]any)
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

	// Persist the conversation message and recommendations to the database
	err = h.conversationService.SendMessage(
		ctx,
		req.ConversationID,
		userID.(uint64),
		req.Message,
		context,
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to persist message to database, continuing anyway")
		// Not returning error as we can still continue with in-memory conversation
	}

	// Prepare and send response
	response := responses.ConversationMessageResponse{
		Message:         message,
		Recommendations: recommendations,
		Context:         make(map[string]any),
	}

	responses.RespondOK(c, response, "Message sent successfully")
}