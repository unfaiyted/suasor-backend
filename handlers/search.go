package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"suasor/clients/media/types"
	"suasor/services"
	"suasor/types/responses"
	"suasor/utils/logger"
)

// SearchHandler handles all search operations
type SearchHandler struct {
	service services.SearchService
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(service services.SearchService) *SearchHandler {
	return &SearchHandler{service: service}
}

// Search godoc
// @Summary Search for content across all sources
// @Description Searches for content in the database, media clients, and metadata sources
// @Tags search
// @Accept json
// @Produce json
// @Param query query string true "Search query"
// @Param mediaType query string false "Limit search to specific media type (movie, series, music, person)"
// @Param limit query int false "Maximum number of results" default(20)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} responses.SearchResponse
// @Failure 400 {object} responses.ErrorResponse[any]
// @Failure 500 {object} responses.ErrorResponse[any]
// @Router /search [get]
func (h *SearchHandler) Search(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get current user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		err := errors.New("user not authenticated")
		responses.RespondUnauthorized(c, err, "User not authenticated")
		return
	}

	// Parse query parameters
	query := c.Query("query")
	if query == "" {
		err := errors.New("query parameter is required")
		responses.RespondBadRequest(c, err, "Query parameter is required")
		return
	}

	// Parse other query parameters
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")
	mediaType := c.Query("mediaType")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Create query options
	options := types.QueryOptions{
		Query:     query,
		Limit:     limit,
		Offset:    offset,
		MediaType: types.MediaType(mediaType),
	}

	// Perform search
	results, err := h.service.SearchAll(ctx, userID.(uint64), options)
	if err != nil {
		log.Error().Err(err).Msg("Error performing search")
		responses.RespondInternalError(c, err, "Error performing search")
		return
	}

	// Convert to response format
	response := responses.ConvertToSearchResponse(results)

	// Return results
	c.JSON(http.StatusOK, response)
}

// GetRecentSearches godoc
// @Summary Get recent searches for the current user
// @Description Returns a list of the user's recent searches
// @Tags search
// @Accept json
// @Produce json
// @Param limit query int false "Maximum number of results" default(10)
// @Success 200 {object} responses.RecentSearchesResponse
// @Failure 400 {object} responses.ErrorResponse[any]
// @Failure 500 {object} responses.ErrorResponse[any]
// @Router /search/recent [get]
func (h *SearchHandler) GetRecentSearches(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get current user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		err := errors.New("user not authenticated")
		responses.RespondUnauthorized(c, err, "User not authenticated")
		return
	}

	// Parse limit parameter
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		limit = 10
	}

	// Get recent searches
	searches, err := h.service.GetRecentSearches(ctx, userID.(uint64), limit)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving recent searches")
		responses.RespondInternalError(c, err, "Error retrieving recent searches")
		return
	}

	// Convert to response format
	response := responses.ConvertToRecentSearchesResponse(searches)

	// Return results
	c.JSON(http.StatusOK, response)
}

// GetTrendingSearches godoc
// @Summary Get trending searches across all users
// @Description Returns a list of popular searches across the platform
// @Tags search
// @Accept json
// @Produce json
// @Param limit query int false "Maximum number of results" default(10)
// @Success 200 {object} responses.TrendingSearchesResponse
// @Failure 400 {object} responses.ErrorResponse[any]
// @Failure 500 {object} responses.ErrorResponse[any]
// @Router /search/trending [get]
func (h *SearchHandler) GetTrendingSearches(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Parse limit parameter
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		limit = 10
	}

	// Get trending searches
	searches, err := h.service.GetTrendingSearches(ctx, limit)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving trending searches")
		responses.RespondInternalError(c, err, "Error retrieving trending searches")
		return
	}

	// Convert to response format
	response := responses.ConvertToTrendingSearchesResponse(searches)

	// Return results
	c.JSON(http.StatusOK, response)
}

// GetSearchSuggestions godoc
// @Summary Get search suggestions
// @Description Returns suggestions based on partial search input
// @Tags search
// @Accept json
// @Produce json
// @Param q query string true "Partial search query"
// @Param limit query int false "Maximum number of suggestions" default(5)
// @Success 200 {object} responses.SearchSuggestionsResponse
// @Failure 400 {object} responses.ErrorResponse[any]
// @Failure 500 {object} responses.ErrorResponse[any]
// @Router /search/suggestions [get]
func (h *SearchHandler) GetSearchSuggestions(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Parse query parameters
	partialQuery := c.Query("q")
	if partialQuery == "" {
		err := errors.New("query parameter 'q' is required")
		responses.RespondBadRequest(c, err, "Query parameter 'q' is required")
		return
	}

	// Parse limit parameter
	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		limit = 5
	}

	// Get search suggestions
	suggestions, err := h.service.GetSearchSuggestions(ctx, partialQuery, limit)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving search suggestions")
		responses.RespondInternalError(c, err, "Error retrieving search suggestions")
		return
	}

	// Create response
	response := responses.SearchSuggestionsResponse{
		Success:     true,
		Suggestions: suggestions,
	}

	// Return results
	c.JSON(http.StatusOK, response)
}
