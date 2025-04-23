package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"suasor/services"
	"suasor/types/requests"
	"suasor/types/responses"
)

// RecommendationHandler interface defines the methods for the recommendation handler
// type RecommendationHandler interface {
// 	// GetRecommendations retrieves recommendations for a user with optional filtering
// 	GetRecommendations(c *gin.Context)
// 	// GetRecommendationByID retrieves a specific recommendation by ID
// 	GetRecommendationByID(c *gin.Context)
// 	// GetRecentRecommendations retrieves recent recommendations for a user
// 	GetRecentRecommendations(c *gin.Context)
// 	// GetTopRecommendations retrieves top-scored recommendations for a user
// 	GetTopRecommendations(c *gin.Context)
// 	// MarkRecommendationAsViewed marks a recommendation as viewed
// 	MarkRecommendationAsViewed(c *gin.Context)
// 	// RateRecommendation allows a user to rate a recommendation
// 	RateRecommendation(c *gin.Context)
// }

// RecommendationHandler handles API requests for recommendations
type RecommendationHandler struct {
	recommendationService services.RecommendationService
}

// NewRecommendationHandler creates a new handler for recommendations
func NewRecommendationHandler(recommendationService services.RecommendationService) *RecommendationHandler {
	return &RecommendationHandler{
		recommendationService: recommendationService,
	}
}

// GetRecommendations godoc
// @Summary Get recommendations for the current user
// @Description Retrieves a list of recommendations for the authenticated user
// @Tags recommendations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param mediaType query string false "Filter by media type (movie, series, music)"
// @Param limit query int false "Number of recommendations to return (default: 20)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {object} responses.APIResponse[responses.RecommendationsListResponse] "Recommendations retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request parameters"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/recommendations [get]
func (h *RecommendationHandler) GetRecommendations(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	// Parse query parameters
	var req requests.GetRecommendationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		responses.RespondValidationError(c, err)
		return
	}

	// Set defaults if not provided
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	// Get recommendations from service
	recommendations, err := h.recommendationService.GetRecommendations(
		c.Request.Context(),
		userID.(uint64),
		req.MediaType,
		req.Limit,
		req.Offset,
	)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to retrieve recommendations")
		return
	}

	// Get total count for pagination
	total, err := h.recommendationService.GetRecommendationCount(
		c.Request.Context(),
		userID.(uint64),
		req.MediaType,
	)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to retrieve recommendation count")
		return
	}

	// Convert to response format
	resp := responses.ConvertToRecommendationsListResponse(
		recommendations,
		total,
		req.Limit,
		req.Offset,
		req.MediaType,
	)

	responses.RespondOK(c, resp, "Recommendations retrieved successfully")
}

// GetRecentRecommendations godoc
// @Summary Get recent recommendations for the current user
// @Description Retrieves a list of recently created recommendations for the authenticated user
// @Tags recommendations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param days query int false "Number of days to look back (default: 7)"
// @Param mediaType query string false "Filter by media type (movie, series, music)"
// @Param limit query int false "Number of recommendations to return (default: 20)"
// @Success 200 {object} responses.APIResponse[responses.RecommendationsListResponse] "Recent recommendations retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request parameters"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/recommendations/recent [get]
func (h *RecommendationHandler) GetRecentRecommendations(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	// Parse query parameters
	var req requests.GetRecentRecommendationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		responses.RespondValidationError(c, err)
		return
	}

	// Set defaults if not provided
	if req.Days <= 0 {
		req.Days = 7
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	// Get recent recommendations from service
	recommendations, err := h.recommendationService.GetRecentRecommendations(
		c.Request.Context(),
		userID.(uint64),
		req.Days,
		req.Limit,
	)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to retrieve recent recommendations")
		return
	}

	// Get total count for pagination
	total, err := h.recommendationService.GetRecommendationCount(
		c.Request.Context(),
		userID.(uint64),
		req.MediaType,
	)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to retrieve recommendation count")
		return
	}

	// Convert to response format
	resp := responses.ConvertToRecommendationsListResponse(
		recommendations,
		total,
		req.Limit,
		0, // Offset is always 0 for recent recommendations
		req.MediaType,
	)

	responses.RespondOK(c, resp, "Recent recommendations retrieved successfully")
}

// GetTopRecommendations godoc
// @Summary Get top-scored recommendations for the current user
// @Description Retrieves a list of top-scored recommendations for the authenticated user
// @Tags recommendations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param minScore query string false "Minimum score (0-1) for recommendations (default: 0.7)"
// @Param mediaType query string false "Filter by media type (movie, series, music)"
// @Param limit query int false "Number of recommendations to return (default: 20)"
// @Success 200 {object} responses.APIResponse[responses.RecommendationsListResponse] "Top recommendations retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request parameters"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/recommendations/top [get]
func (h *RecommendationHandler) GetTopRecommendations(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	// Parse query parameters
	var req requests.GetTopRecommendationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		responses.RespondValidationError(c, err)
		return
	}

	// Set defaults if not provided
	if req.MinScore <= 0 {
		req.MinScore = 0.7
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	// Get top recommendations from service
	recommendations, err := h.recommendationService.GetTopRecommendations(
		c.Request.Context(),
		userID.(uint64),
		req.MinScore,
		req.Limit,
	)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to retrieve top recommendations")
		return
	}

	// Get total count for pagination
	total, err := h.recommendationService.GetRecommendationCount(
		c.Request.Context(),
		userID.(uint64),
		req.MediaType,
	)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to retrieve recommendation count")
		return
	}

	// Convert to response format
	resp := responses.ConvertToRecommendationsListResponse(
		recommendations,
		total,
		req.Limit,
		0, // Offset is always 0 for top recommendations
		req.MediaType,
	)

	responses.RespondOK(c, resp, "Top recommendations retrieved successfully")
}

// GetRecommendationByID godoc
// @Summary Get a specific recommendation by ID
// @Description Retrieves a specific recommendation by its ID for the authenticated user
// @Tags recommendations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Recommendation ID"
// @Success 200 {object} responses.APIResponse[responses.RecommendationResponse] "Recommendation retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid recommendation ID"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 404 {object} responses.ErrorResponse[responses.ErrorDetails] "Recommendation not found"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/recommendations/{id} [get]
func (h *RecommendationHandler) GetRecommendationByID(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	// Parse recommendation ID from path
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		responses.RespondValidationError(c, err)
		return
	}

	// Get recommendation from service
	recommendation, err := h.recommendationService.GetRecommendationByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "record not found" {
			responses.RespondNotFound(c, nil, "Recommendation not found")
			return
		}
		responses.RespondInternalError(c, err, "Failed to retrieve recommendation")
		return
	}

	// Check if recommendation belongs to the user
	if recommendation.UserID != userID.(uint64) {
		responses.RespondNotFound(c, nil, "Recommendation not found")
		return
	}

	// Convert to response format
	resp := responses.ConvertToRecommendationResponse(*recommendation)

	responses.RespondOK(c, resp, "Recommendation retrieved successfully")
}

// MarkRecommendationAsViewed godoc
// @Summary Mark a recommendation as viewed
// @Description Marks a specific recommendation as viewed for the authenticated user
// @Tags recommendations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body requests.MarkRecommendationAsViewedRequest true "Recommendation ID to mark as viewed"
// @Success 200 {object} responses.APIResponse[any] "Recommendation marked as viewed successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request parameters"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 404 {object} responses.ErrorResponse[responses.ErrorDetails] "Recommendation not found"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/recommendations/view [post]
func (h *RecommendationHandler) MarkRecommendationAsViewed(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	// Parse request body
	var req requests.MarkRecommendationAsViewedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.RespondValidationError(c, err)
		return
	}

	// Mark recommendation as viewed
	err := h.recommendationService.MarkRecommendationAsViewed(
		c.Request.Context(),
		req.RecommendationID,
		userID.(uint64),
	)
	if err != nil {
		if err.Error() == "record not found" {
			responses.RespondNotFound(c, nil, "Recommendation not found")
			return
		}
		if err.Error() == "recommendation does not belong to the user" {
			responses.RespondNotFound(c, nil, "Recommendation not found")
			return
		}
		responses.RespondInternalError(c, err, "Failed to mark recommendation as viewed")
		return
	}

	responses.RespondOK(c, http.StatusOK, "Recommendation marked as viewed successfully")
}

// RateRecommendation godoc
// @Summary Rate a recommendation
// @Description Sets a user rating for a specific recommendation
// @Tags recommendations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body requests.RateRecommendationRequest true "Recommendation ID and rating"
// @Success 200 {object} responses.APIResponse[any] "Recommendation rated successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request parameters"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 404 {object} responses.ErrorResponse[responses.ErrorDetails] "Recommendation not found"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/recommendations/rate [post]
func (h *RecommendationHandler) RateRecommendation(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	// Parse request body
	var req requests.RateRecommendationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.RespondValidationError(c, err)
		return
	}

	// Rate recommendation
	err := h.recommendationService.RateRecommendation(
		c.Request.Context(),
		req.RecommendationID,
		userID.(uint64),
		req.Rating,
	)
	if err != nil {
		if err.Error() == "record not found" {
			responses.RespondNotFound(c, nil, "Recommendation not found")
			return
		}
		if err.Error() == "recommendation does not belong to the user" {
			responses.RespondNotFound(c, nil, "Recommendation not found")
			return
		}
		responses.RespondInternalError(c, err, "Failed to rate recommendation")
		return
	}

	responses.RespondOK(c, http.StatusOK, "Recommendation rated successfully")
}
