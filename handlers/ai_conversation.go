package handlers

import (
	"strconv"
	"suasor/services"
	"suasor/types/requests"
	"suasor/types/responses"
	"suasor/utils/logger"

	"github.com/gin-gonic/gin"
)

// AIConversationHandler defines the interface for conversation history handlers
type AIConversationHandler interface {
	// Conversation history endpoints
	GetUserConversations(c *gin.Context)
	GetConversationHistory(c *gin.Context)
	GetUserRecommendations(c *gin.Context)
	ContinueConversation(c *gin.Context)
	ArchiveConversation(c *gin.Context)
	DeleteConversation(c *gin.Context)
}

// aiConversationHandler implements AIConversationHandler
type aiConversationHandler struct {
	service services.AIConversationService
}

// NewAIConversationHandler creates a new conversation history handler
func NewAIConversationHandler(service services.AIConversationService) AIConversationHandler {
	return &aiConversationHandler{
		service: service,
	}
}

// GetUserConversations godoc
//
//	@Summary		Get user's conversation history
//	@Description	Retrieve a paginated list of the user's AI conversations
//	@Tags			ai, conversations
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			limit	query		int									false	"Number of items to return (default 20)"
//	@Param			offset	query		int									false	"Offset for pagination (default 0)"
//	@Success		200		{object}	responses.APIResponse				"Conversation list retrieved"
//	@Failure		401		{object}	responses.ErrorResponse				"Unauthorized"
//	@Failure		500		{object}	responses.ErrorResponse				"Server error"
//	@Router			/user/ai/conversations [get]
func (h *aiConversationHandler) GetUserConversations(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	// Parse pagination parameters
	limit, offset := getPaginationParams(c)

	log.Info().
		Uint64("userID", userID.(uint64)).
		Int("limit", limit).
		Int("offset", offset).
		Msg("Getting user conversation history")

	// Get conversations from service
	conversations, count, err := h.service.GetUserConversations(ctx, userID.(uint64), limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user conversations")
		responses.RespondInternalError(c, err, "Failed to retrieve conversation history")
		return
	}

	// Create pagination response
	pagination := responses.PaginationData{
		TotalCount: count,
		Limit:      limit,
		Offset:     offset,
	}

	responses.RespondOKWithPagination(c, conversations, pagination, "Conversation history retrieved successfully")
}

// GetConversationHistory godoc
//
//	@Summary		Get messages in a conversation
//	@Description	Retrieve all messages in a specific AI conversation
//	@Tags			ai, conversations
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			conversationId	path		string							true	"Conversation ID"
//	@Success		200				{object}	responses.APIResponse			"Messages retrieved"
//	@Failure		401				{object}	responses.ErrorResponse			"Unauthorized"
//	@Failure		403				{object}	responses.ErrorResponse			"Forbidden - conversation not owned by user"
//	@Failure		404				{object}	responses.ErrorResponse			"Conversation not found"
//	@Failure		500				{object}	responses.ErrorResponse			"Server error"
//	@Router			/user/ai/conversations/{conversationId}/messages [get]
func (h *aiConversationHandler) GetConversationHistory(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	// Get conversation ID from path
	conversationID := c.Param("conversationId")
	if conversationID == "" {
		responses.RespondBadRequest(c, nil, "Conversation ID is required")
		return
	}

	log.Info().
		Uint64("userID", userID.(uint64)).
		Str("conversationID", conversationID).
		Msg("Getting conversation history")

	// Get messages from service
	messages, err := h.service.GetConversationHistory(ctx, conversationID, userID.(uint64))
	if err != nil {
		log.Error().Err(err).Msg("Failed to get conversation history")
		responses.RespondInternalError(c, err, "Failed to retrieve conversation messages")
		return
	}

	if messages == nil || len(messages) == 0 {
		responses.RespondNotFound(c, nil, "Conversation not found or no messages available")
		return
	}

	responses.RespondOK(c, messages, "Conversation messages retrieved successfully")
}

// GetUserRecommendations godoc
//
//	@Summary		Get user's recommendation history
//	@Description	Retrieve a paginated list of AI recommendations made to the user
//	@Tags			ai, recommendations
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			itemType	query		string							false	"Filter by content type (e.g., movie, music)"
//	@Param			limit		query		int								false	"Number of items to return (default 20)"
//	@Param			offset		query		int								false	"Offset for pagination (default 0)"
//	@Success		200			{object}	responses.APIResponse			"Recommendations retrieved"
//	@Failure		401			{object}	responses.ErrorResponse			"Unauthorized"
//	@Failure		500			{object}	responses.ErrorResponse			"Server error"
//	@Router			/user/ai/recommendations [get]
func (h *aiConversationHandler) GetUserRecommendations(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	// Parse pagination parameters
	limit, offset := getPaginationParams(c)

	// Get item type filter if provided
	itemType := c.Query("itemType")

	log.Info().
		Uint64("userID", userID.(uint64)).
		Str("itemType", itemType).
		Int("limit", limit).
		Int("offset", offset).
		Msg("Getting user recommendation history")

	// Get recommendations from service
	recommendations, count, err := h.service.GetUserRecommendationHistory(ctx, userID.(uint64), itemType, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user recommendations")
		responses.RespondInternalError(c, err, "Failed to retrieve recommendation history")
		return
	}

	// Create pagination response
	pagination := responses.PaginationData{
		TotalCount: count,
		Limit:      limit,
		Offset:     offset,
	}

	responses.RespondOKWithPagination(c, recommendations, pagination, "Recommendation history retrieved successfully")
}

// ContinueConversation godoc
//
//	@Summary		Continue a previous conversation
//	@Description	Resume a previous AI conversation
//	@Tags			ai, conversations
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			conversationId	path		string							true	"Conversation ID"
//	@Param			request			body		requests.ContinueConversationRequest	true	"Continue request"
//	@Success		200				{object}	responses.APIResponse			"Conversation resumed"
//	@Failure		400				{object}	responses.ErrorResponse			"Invalid request"
//	@Failure		401				{object}	responses.ErrorResponse			"Unauthorized"
//	@Failure		403				{object}	responses.ErrorResponse			"Forbidden - conversation not owned by user"
//	@Failure		404				{object}	responses.ErrorResponse			"Conversation not found"
//	@Failure		500				{object}	responses.ErrorResponse			"Server error"
//	@Router			/user/ai/conversations/{conversationId}/continue [post]
func (h *aiConversationHandler) ContinueConversation(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	// Get conversation ID from path
	conversationID := c.Param("conversationId")
	if conversationID == "" {
		responses.RespondBadRequest(c, nil, "Conversation ID is required")
		return
	}

	var req requests.ContinueConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.RespondValidationError(c, err)
		return
	}

	log.Info().
		Uint64("userID", userID.(uint64)).
		Str("conversationID", conversationID).
		Str("message", req.Message).
		Msg("Continuing previous conversation")

	// Continue conversation using service
	response, recommendations, err := h.service.SendMessage(
		ctx,
		conversationID,
		userID.(uint64),
		req.Message,
		req.Context,
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to continue conversation")
		responses.RespondInternalError(c, err, "Failed to continue conversation")
		return
	}

	// Prepare response
	apiResponse := responses.ConversationMessageResponse{
		Message:         response,
		Recommendations: recommendations,
		Context:         make(map[string]any),
	}

	responses.RespondOK(c, apiResponse, "Conversation continued successfully")
}

// ArchiveConversation godoc
//
//	@Summary		Archive a conversation
//	@Description	Archive a conversation to remove it from active list
//	@Tags			ai, conversations
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			conversationId	path		string							true	"Conversation ID"
//	@Success		200				{object}	responses.APIResponse			"Conversation archived"
//	@Failure		401				{object}	responses.ErrorResponse			"Unauthorized"
//	@Failure		403				{object}	responses.ErrorResponse			"Forbidden - conversation not owned by user"
//	@Failure		404				{object}	responses.ErrorResponse			"Conversation not found"
//	@Failure		500				{object}	responses.ErrorResponse			"Server error"
//	@Router			/user/ai/conversations/{conversationId}/archive [put]
func (h *aiConversationHandler) ArchiveConversation(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	// Get conversation ID from path
	conversationID := c.Param("conversationId")
	if conversationID == "" {
		responses.RespondBadRequest(c, nil, "Conversation ID is required")
		return
	}

	log.Info().
		Uint64("userID", userID.(uint64)).
		Str("conversationID", conversationID).
		Msg("Archiving conversation")

	// Archive conversation
	err := h.service.ArchiveConversation(ctx, conversationID, userID.(uint64))
	if err != nil {
		log.Error().Err(err).Msg("Failed to archive conversation")
		responses.RespondInternalError(c, err, "Failed to archive conversation")
		return
	}

	responses.RespondOK[any](c, nil, "Conversation archived successfully")
}

// DeleteConversation godoc
//
//	@Summary		Delete a conversation
//	@Description	Permanently delete a conversation and all its messages
//	@Tags			ai, conversations
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			conversationId	path		string							true	"Conversation ID"
//	@Success		200				{object}	responses.APIResponse			"Conversation deleted"
//	@Failure		401				{object}	responses.ErrorResponse			"Unauthorized"
//	@Failure		403				{object}	responses.ErrorResponse			"Forbidden - conversation not owned by user"
//	@Failure		404				{object}	responses.ErrorResponse			"Conversation not found"
//	@Failure		500				{object}	responses.ErrorResponse			"Server error"
//	@Router			/user/ai/conversations/{conversationId} [delete]
func (h *aiConversationHandler) DeleteConversation(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	// Get conversation ID from path
	conversationID := c.Param("conversationId")
	if conversationID == "" {
		responses.RespondBadRequest(c, nil, "Conversation ID is required")
		return
	}

	log.Info().
		Uint64("userID", userID.(uint64)).
		Str("conversationID", conversationID).
		Msg("Deleting conversation")

	// Delete conversation
	err := h.service.DeleteConversation(ctx, conversationID, userID.(uint64))
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete conversation")
		responses.RespondInternalError(c, err, "Failed to delete conversation")
		return
	}

	responses.RespondOK[any](c, nil, "Conversation deleted successfully")
}

// Helper functions

// getPaginationParams extracts and validates pagination parameters
func getPaginationParams(c *gin.Context) (limit, offset int) {
	// Default values
	limit = 20
	offset = 0

	// Parse limit if provided
	if limitParam := c.Query("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil {
			if parsedLimit >= 1 && parsedLimit <= 100 {
				limit = parsedLimit
			}
		}
	}

	// Parse offset if provided
	if offsetParam := c.Query("offset"); offsetParam != "" {
		if parsedOffset, err := strconv.Atoi(offsetParam); err == nil {
			if parsedOffset >= 0 && parsedOffset <= 1000000 {
				offset = parsedOffset
			}
		}
	}

	return limit, offset
}