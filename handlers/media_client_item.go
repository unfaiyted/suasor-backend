package handlers

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"suasor/client/media/types"
	clientTypes "suasor/client/types"
	"suasor/services"
	"suasor/types/models"
	"suasor/types/responses"
	"suasor/utils"
)

// MediaClientItemHandler handles operations on media items stored in the database
// It provides access to different types of media (movies, series, music, etc.)
// using a generic type parameter
type MediaClientItemHandler[T types.MediaData] struct {
	service services.MediaItemService[T]
}

// NewMediaClientItemHandler creates a new media item handler
func NewMediaClientItemHandler[T types.MediaData](service services.MediaItemService[T]) *MediaClientItemHandler[T] {
	return &MediaClientItemHandler[T]{service: service}
}

// CreateMediaItem godoc
// @Summary Create a new media item
// @Description Creates a new media item in the database
// @Tags media-items
// @Accept json
// @Produce json
// @Param mediaItem body object true "Media item data with type, client info, and type-specific data"
// @Success 201 {object} responses.APIResponse[models.MediaItem[any]] "Media item created successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /item/media [post]
func (h *MediaClientItemHandler[T]) CreateMediaItem(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	var req struct {
		Type       string          `json:"type" binding:"required"`
		ClientID   uint64          `json:"clientId" binding:"required"`
		ClientType string          `json:"clientType" binding:"required"`
		ExternalID string          `json:"externalId" binding:"required"`
		Data       json.RawMessage `json:"data" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid request body for CreateMediaItem")
		responses.RespondBadRequest(c, err, "Invalid request format")
		return
	}

	mediaType := types.MediaType(req.Type)
	log.Debug().
		Str("mediaType", string(mediaType)).
		Uint64("clientId", req.ClientID).
		Str("clientType", req.ClientType).
		Msg("Creating media item")

	// Create the media item
	var mediaData T
	if err := json.Unmarshal(req.Data, &mediaData); err != nil {
		log.Warn().Err(err).Msg("Failed to unmarshal media data")
		responses.RespondBadRequest(c, err, "Invalid media data format")
		return
	}

	mediaItem := models.MediaItem[T]{
		Type:        mediaType,
		SyncClients: []models.SyncClient{},
		ExternalIDs: []models.ExternalID{},
		Data:        mediaData,
	}

	// Set client info
	mediaItem.SetClientInfo(req.ClientID, clientTypes.MediaClientType(req.ClientType), req.ExternalID)

	// Only add external ID if provided
	if req.ExternalID != "" {
		mediaItem.AddExternalID("client", req.ExternalID)
	}

	result, err := h.service.Create(ctx, mediaItem)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create media item")
		responses.RespondInternalError(c, err, "Failed to create media item")
		return
	}

	log.Info().
		Uint64("id", result.ID).
		Str("type", string(result.Type)).
		Msg("Media item created successfully")

	responses.RespondCreated(c, result, "Media item created successfully")
}

// UpdateMediaItem godoc
// @Summary Update an existing media item
// @Description Updates a media item in the database by ID
// @Tags media-items
// @Accept json
// @Produce json
// @Param id path int true "Media item ID"
// @Param mediaItem body object true "Media item data to update"
// @Success 200 {object} responses.APIResponse[models.MediaItem[any]] "Media item updated successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Media item not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /item/media/{id} [put]
func (h *MediaClientItemHandler[T]) UpdateMediaItem(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid media item ID")
		responses.RespondBadRequest(c, err, "Invalid media item ID")
		return
	}

	var req struct {
		Type       string          `json:"type" binding:"required"`
		ClientID   uint64          `json:"clientId" binding:"required"`
		ClientType string          `json:"clientType" binding:"required"`
		ExternalID string          `json:"externalId" binding:"required"`
		Data       json.RawMessage `json:"data" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Uint64("id", id).Msg("Invalid request body for UpdateMediaItem")
		responses.RespondValidationError(c, err)
		return
	}

	log.Debug().
		Uint64("id", id).
		Str("mediaType", req.Type).
		Uint64("clientId", req.ClientID).
		Str("clientType", req.ClientType).
		Msg("Updating media item")

	// Update the media item
	mediaType := types.MediaType(req.Type)
	var mediaData T
	if err := json.Unmarshal(req.Data, &mediaData); err != nil {
		log.Warn().Err(err).Uint64("id", id).Msg("Failed to unmarshal media data")
		responses.RespondBadRequest(c, err, "Invalid media data format")
		return
	}

	mediaItem := models.MediaItem[T]{
		ID:          id,
		Type:        mediaType,
		SyncClients: []models.SyncClient{},
		ExternalIDs: []models.ExternalID{},
		Data:        mediaData,
	}

	// Set client info
	mediaItem.SetClientInfo(req.ClientID, clientTypes.MediaClientType(req.ClientType), req.ExternalID)

	// Only add external ID if provided
	if req.ExternalID != "" {
		mediaItem.AddExternalID("client", req.ExternalID)
	}

	result, err := h.service.Update(ctx, mediaItem)
	if err != nil {
		log.Error().Err(err).Uint64("id", id).Msg("Failed to update media item")
		responses.RespondInternalError(c, err, "Failed to update media item")
	}

	log.Info().
		Uint64("id", result.ID).
		Str("type", string(result.Type)).
		Msg("Media item updated successfully")

	responses.RespondOK(c, result, "Media item updated successfully")
}

// GetMediaItem godoc
// @Summary Get a media item by ID
// @Description Retrieves a media item from the database by its ID
// @Tags media-items
// @Accept json
// @Produce json
// @Param id path int true "Media item ID"
// @Success 200 {object} responses.APIResponse[models.MediaItem[any]] "Media item retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid media item ID"
// @Failure 404 {object} responses.ErrorResponse[any] "Media item not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /item/media/{id} [get]
func (h *MediaClientItemHandler[T]) GetMediaItem(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid media item ID")
		responses.RespondBadRequest(c, err, "Invalid media item ID")
		return
	}

	log.Debug().Uint64("id", id).Msg("Retrieving media item by ID")
	h.getMediaItem(c, id)
}

// Type-specific media item retrieval helper
func (h *MediaClientItemHandler[T]) getMediaItem(c *gin.Context, id uint64) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	item, err := h.service.GetByID(ctx, id)
	if err != nil {
		log.Warn().Err(err).Uint64("id", id).Msg("Media item not found")
		responses.RespondNotFound(c, err, "Media item not found")
		return
	}

	log.Info().Uint64("id", id).Str("type", string(item.Type)).Msg("Media item retrieved successfully")
	responses.RespondOK(c, item, "Media item retrieved successfully")
}

// GetMediaItemsByClient godoc
// @Summary Get media items by client
// @Description Retrieves all media items for a specific client
// @Tags media-items
// @Accept json
// @Produce json
// @Param clientId path int true "Client ID"
// @Param type query string false "Media type filter"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[any]] "Media items retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid client ID"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Failure 501 {object} responses.ErrorResponse[any] "Not implemented"
// @Router /item/media/client/{clientId} [get]
func (h *MediaClientItemHandler[T]) GetMediaItemsByClient(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	clientID, err := strconv.ParseUint(c.Param("clientId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("clientId", c.Param("clientId")).Msg("Invalid client ID")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	mediaType := c.Query("type") // Optional media type filter
	log.Debug().
		Uint64("clientId", clientID).
		Str("mediaType", mediaType).
		Msg("Getting media items by client")

	if mediaType != "" {
		h.getMediaItemsByClient(c, clientID)
		return
	}

	// If no specific type, we'll need to collect all types and merge
	log.Warn().Uint64("clientId", clientID).Msg("Fetching all media types not yet implemented")
	responses.RespondNotImplemented(c, nil, "Fetching all media types not yet implemented")
}

// Type-specific media items by client retrieval helper
func (h *MediaClientItemHandler[T]) getMediaItemsByClient(c *gin.Context, clientID uint64) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	items, err := h.service.GetByClientID(ctx, clientID)
	if err != nil {
		log.Error().Err(err).Uint64("clientId", clientID).Msg("Failed to retrieve media items by client")
		responses.RespondInternalError(c, err, "Failed to retrieve media items by client")
		return
	}

	log.Info().Uint64("clientId", clientID).Int("count", len(items)).Msg("Media items retrieved by client successfully")
	responses.RespondOK(c, items, "Media items retrieved successfully")
}

// DeleteMediaItem godoc
// @Summary Delete a media item
// @Description Deletes a media item from the database by ID
// @Tags media-items
// @Accept json
// @Produce json
// @Param id path int true "Media item ID"
// @Success 200 {object} responses.APIResponse[any] "Media item deleted successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid media item ID"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /item/media/{id} [delete]
func (h *MediaClientItemHandler[T]) DeleteMediaItem(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid media item ID")
		responses.RespondBadRequest(c, err, "Invalid media item ID")
		return
	}

	log.Debug().Uint64("id", id).Msg("Deleting media item")
	h.deleteMediaItem(c, id)
}

// Type-specific media item deletion helper
func (h *MediaClientItemHandler[T]) deleteMediaItem(c *gin.Context, id uint64) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	err := h.service.Delete(ctx, id)
	if err != nil {
		log.Error().Err(err).Uint64("id", id).Msg("Failed to delete media item")
		responses.RespondInternalError(c, err, "Failed to delete media item")
		return
	}

	log.Info().Uint64("id", id).Msg("Media item deleted successfully")
	responses.RespondOK(c, responses.EmptyResponse{Success: true}, "Media item deleted successfully")
}

// SearchMediaItems godoc
// @Summary Search for media items
// @Description Searches for media items by title or other criteria
// @Tags media-items
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param userId query int true "User ID"
// @Param type query string false "Media type filter"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[any]] "Media items retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /item/media/search [get]
func (h *MediaClientItemHandler[T]) SearchMediaItems(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	query := c.Query("q")
	if query == "" {
		log.Warn().Msg("Search query is required")
		responses.RespondBadRequest(c, nil, "Search query is required")
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	mediaType := c.Query("type") // Optional media type filter
	log.Debug().
		Str("query", query).
		Uint64("userId", userID).
		Str("mediaType", mediaType).
		Msg("Searching media items")

	if mediaType != "" {
		h.searchMediaItems(c, query, userID)
		return
	}

	// If no specific type, search across all types (not implemented here)
	log.Warn().Msg("Searching across all media types not yet implemented")
	responses.RespondNotImplemented(c, nil, "Searching across all media types not yet implemented")
}

// Type-specific media items search helper
func (h *MediaClientItemHandler[T]) searchMediaItems(c *gin.Context, query string, userID uint64) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	items, err := h.service.SearchByTitle(ctx, query, userID)
	if err != nil {
		log.Error().Err(err).Str("query", query).Uint64("userId", userID).Msg("Failed to search media items")
		responses.RespondInternalError(c, err, "Failed to search media items")
		return
	}

	log.Info().Str("query", query).Uint64("userId", userID).Int("count", len(items)).Msg("Media items search completed successfully")
	responses.RespondOK(c, items, "Media items retrieved successfully")
}

// GetRecentMediaItems godoc
// @Summary Get recently added media items
// @Description Retrieves recently added media items for a user
// @Tags media-items
// @Accept json
// @Produce json
// @Param userId query int true "User ID"
// @Param limit query int false "Maximum number of items to return"
// @Param type query string false "Media type filter"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[any]] "Recent media items retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /item/media/recent [get]
func (h *MediaClientItemHandler[T]) GetRecentMediaItems(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	mediaType := c.Query("type") // Optional media type filter
	log.Debug().
		Uint64("userId", userID).
		Int("limit", limit).
		Str("mediaType", mediaType).
		Msg("Getting recent media items")

	if mediaType != "" {
		h.getRecentMediaItems(c, userID, limit)
		return
	}

	// If no specific type, get recent items across all types (not implemented here)
	log.Warn().Msg("Fetching recent items across all media types not yet implemented")
	responses.RespondNotImplemented(c, nil, "Fetching recent items across all media types not yet implemented")
}

// Type-specific recent media items retrieval helper
func (h *MediaClientItemHandler[T]) getRecentMediaItems(c *gin.Context, userID uint64, limit int) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	items, err := h.service.GetRecentItems(ctx, userID, limit)
	if err != nil {
		log.Error().Err(err).Uint64("userId", userID).Int("limit", limit).Msg("Failed to retrieve recent media items")
		responses.RespondInternalError(c, err, "Failed to retrieve recent media items")
		return
	}

	log.Info().Uint64("userId", userID).Int("limit", limit).Int("count", len(items)).Msg("Recent media items retrieved successfully")
	responses.RespondOK(c, items, "Recent media items retrieved successfully")
}

// GetMediaItemByExternalSourceID godoc
// @Summary Get a media item by external source ID
// @Description Retrieves a media item using its external source ID (e.g., TMDB ID)
// @Tags media-items
// @Accept json
// @Produce json
// @Param source path string true "External source name (e.g., tmdb, imdb)"
// @Param externalId path string true "External ID from the source"
// @Success 200 {object} responses.APIResponse[models.MediaItem[any]] "Media item retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Media item not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /item/media/external/{source}/{externalId} [get]
func (h *MediaClientItemHandler[T]) GetMediaItemByExternalSourceID(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	source := c.Param("source")
	externalID := c.Param("externalId")

	if source == "" || externalID == "" {
		log.Warn().Str("source", source).Str("externalId", externalID).Msg("Source and externalId are required")
		responses.RespondBadRequest(c, nil, "Source and externalId are required")
		return
	}

	log.Debug().Str("source", source).Str("externalId", externalID).Msg("Retrieving media item by external ID")
	item, err := h.service.GetByExternalID(ctx, source, externalID)
	if err != nil {
		log.Warn().Err(err).Str("source", source).Str("externalId", externalID).Msg("Media item not found by external ID")
		responses.RespondNotFound(c, err, "Media item not found")
		return
	}

	log.Info().Uint64("id", item.ID).Str("source", source).Str("externalId", externalID).Msg("Media item retrieved by external ID")
	responses.RespondOK(c, item, "Media item retrieved successfully")
}

// GetMediaItemsByGenre godoc
// @Summary Get media items by genre
// @Description Retrieves media items that belong to a specific genre
// @Tags media-items
// @Accept json
// @Produce json
// @Param genre path string true "Genre name"
// @Param userId query int true "User ID"
// @Param limit query int false "Maximum number of items to return"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[any]] "Media items retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /item/media/genre/{genre} [get]
func (h *MediaClientItemHandler[T]) GetMediaItemsByGenre(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	genre := c.Param("genre")
	if genre == "" {
		log.Warn().Msg("Genre is required")
		responses.RespondBadRequest(c, nil, "Genre is required")
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	log.Debug().
		Str("genre", genre).
		Uint64("userId", userID).
		Int("limit", limit).
		Msg("Getting media items by genre")

	// Filter by genre using the search functionality
	// This is a basic implementation - ideally this would be a more efficient query
	items, err := h.service.GetByUserID(ctx, userID)
	if err != nil {
		log.Error().Err(err).Str("genre", genre).Uint64("userId", userID).Msg("Failed to retrieve media items")
		responses.RespondInternalError(c, err, "Failed to retrieve media items")
		return
	}

	// Filter for items with the specified genre
	var filtered []*models.MediaItem[T]
	for _, item := range items {
		details := item.Data.GetDetails()
		for _, g := range details.Genres {
			if strings.EqualFold(g, genre) {
				filtered = append(filtered, item)
				break
			}
		}

		if len(filtered) >= limit {
			break
		}
	}

	log.Info().Str("genre", genre).Uint64("userId", userID).Int("count", len(filtered)).Msg("Media items retrieved by genre")
	responses.RespondOK(c, filtered, "Media items retrieved successfully")
}

// GetMediaItemsByYear godoc
// @Summary Get media items by release year
// @Description Retrieves media items released in a specific year
// @Tags media-items
// @Accept json
// @Produce json
// @Param year path int true "Release year"
// @Param userId query int true "User ID"
// @Param limit query int false "Maximum number of items to return"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[any]] "Media items retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /item/media/year/{year} [get]
func (h *MediaClientItemHandler[T]) GetMediaItemsByYear(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	yearStr := c.Param("year")
	if yearStr == "" {
		log.Warn().Msg("Year is required")
		responses.RespondBadRequest(c, nil, "Year is required")
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		log.Warn().Err(err).Str("year", yearStr).Msg("Invalid year format")
		responses.RespondBadRequest(c, err, "Invalid year format")
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	log.Debug().
		Int("year", year).
		Uint64("userId", userID).
		Int("limit", limit).
		Msg("Getting media items by year")

	// Filter by year using the search functionality
	// This is a basic implementation - ideally this would be a more efficient query
	items, err := h.service.GetByUserID(ctx, userID)
	if err != nil {
		log.Error().Err(err).Int("year", year).Uint64("userId", userID).Msg("Failed to retrieve media items")
		responses.RespondInternalError(c, err, "Failed to retrieve media items")
		return
	}

	// Filter for items with the specified year
	var filtered []*models.MediaItem[T]
	for _, item := range items {
		if item.ReleaseYear == year {
			filtered = append(filtered, item)
		}

		if len(filtered) >= limit {
			break
		}
	}

	log.Info().Int("year", year).Uint64("userId", userID).Int("count", len(filtered)).Msg("Media items retrieved by year")
	responses.RespondOK(c, filtered, "Media items retrieved successfully")
}

// GetMediaItemsByPerson godoc
// @Summary Get media items by person
// @Description Retrieves media items associated with a specific person (actor, director, etc.)
// @Tags media-items
// @Accept json
// @Produce json
// @Param personId path int true "Person ID"
// @Param role query string false "Role filter (actor, director, etc.)"
// @Param userId query int true "User ID"
// @Param limit query int false "Maximum number of items to return"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[any]] "Media items retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Failure 501 {object} responses.ErrorResponse[any] "Not implemented"
// @Router /item/media/person/{personId} [get]
func (h *MediaClientItemHandler[T]) GetMediaItemsByPerson(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	personID := c.Param("personId")
	role := c.Query("role") // Optional role filter (actor, director, etc.)

	if personID == "" {
		log.Warn().Msg("Person ID is required")
		responses.RespondBadRequest(c, nil, "Person ID is required")
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	log.Debug().
		Str("personId", personID).
		Str("role", role).
		Uint64("userId", userID).
		Int("limit", limit).
		Msg("Request for media items by person")

	// This would ideally be implemented at the repository level with a proper join to the credits table
	// For now, we'll respond with not implemented
	log.Info().Str("personId", personID).Msg("Person-based filtering not yet implemented")
	responses.RespondNotImplemented(c, nil, "Person-based filtering not yet implemented")
}

// GetPopularMediaItems godoc
// @Summary Get popular media items
// @Description Retrieves popular media items based on play counts, ratings, etc.
// @Tags media-items
// @Accept json
// @Produce json
// @Param userId query int true "User ID"
// @Param limit query int false "Maximum number of items to return"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[any]] "Popular media items retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /item/media/popular [get]
func (h *MediaClientItemHandler[T]) GetPopularMediaItems(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	count, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		count = 10
	}

	log.Debug().
		Uint64("userId", userID).
		Int("limit", count).
		Msg("Getting popular media items")

	// This could be implemented based on play counts, ratings, etc.
	// For now, we'll just return recent items as a fallback
	items, err := h.service.GetRecentItems(ctx, userID, count)
	if err != nil {
		log.Error().Err(err).Uint64("userId", userID).Msg("Failed to retrieve popular media items")
		responses.RespondInternalError(c, err, "Failed to retrieve popular media items")
		return
	}

	log.Info().Uint64("userId", userID).Int("count", len(items)).Msg("Popular media items retrieved successfully")
	responses.RespondOK(c, items, "Popular media items retrieved successfully")
}

// GetTopRatedMediaItems godoc
// @Summary Get top rated media items
// @Description Retrieves media items with the highest ratings
// @Tags media-items
// @Accept json
// @Produce json
// @Param userId query int true "User ID"
// @Param limit query int false "Maximum number of items to return"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[any]] "Top rated media items retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /item/media/top-rated [get]
func (h *MediaClientItemHandler[T]) GetTopRatedMediaItems(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	count, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		count = 10
	}

	log.Debug().
		Uint64("userId", userID).
		Int("limit", count).
		Msg("Getting top rated media items")

	// This should be implemented with a proper query that sorts by ratings
	// For now, we'll just return recent items as a fallback
	items, err := h.service.GetRecentItems(ctx, userID, count)
	if err != nil {
		log.Error().Err(err).Uint64("userId", userID).Msg("Failed to retrieve top rated media items")
		responses.RespondInternalError(c, err, "Failed to retrieve top rated media items")
		return
	}

	log.Info().Uint64("userId", userID).Int("count", len(items)).Msg("Top rated media items retrieved successfully")
	responses.RespondOK(c, items, "Top rated media items retrieved successfully")
}

// GetAllMediaItems godoc
// @Summary Get all media items
// @Description Retrieves all media items with optional filtering
// @Tags media-items
// @Accept json
// @Produce json
// @Param userId query int true "User ID"
// @Param limit query int false "Maximum number of items to return"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[any]] "Media items retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /item/media/all [get]
func (h *MediaClientItemHandler[T]) GetAllMediaItems(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	log.Debug().
		Uint64("userId", userID).
		Int("limit", limit).
		Msg("Getting all media items")

	items, err := h.service.GetByUserID(ctx, userID)
	if err != nil {
		log.Error().Err(err).Uint64("userId", userID).Msg("Failed to retrieve media items")
		responses.RespondInternalError(c, err, "Failed to retrieve media items")
		return
	}

	// Limit the results
	totalCount := len(items)
	if totalCount > limit {
		items = items[:limit]
	}

	log.Info().Uint64("userId", userID).Int("totalCount", totalCount).Int("returnedCount", len(items)).Msg("Media items retrieved successfully")
	responses.RespondOK(c, items, "Media items retrieved successfully")
}

// GetLatestMediaItemsByAdded godoc
// @Summary Get latest media items by added
// @Description Retrieves recently added media items for a user
// @Tags media-items
// @Accept json
// @Produce json
// @Param userId query int true "User ID"
// @Param limit query int false "Maximum number of items to return"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[any]] "Media items retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /item/media/latest [get]
func (h *MediaClientItemHandler[T]) GetLatestMediaItemsByAdded(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	log.Debug().
		Uint64("userId", userID).
		Int("limit", limit).
		Msg("Getting latest media items by added")

	// Get all media items for the user
	recentMediaItems, err := h.service.GetRecentItems(ctx, userID, limit)
	if err != nil {
		log.Error().Err(err).
			Uint64("userId", userID).
			Msg("Failed to retrieve media items")
		responses.RespondInternalError(c, err, "Failed to retrieve media items")
		return
	}
	log.Info().
		Uint64("userId", userID).
		Int("count", len(recentMediaItems)).
		Msg("Media items retrieved successfully")

	responses.RespondOK(c, recentMediaItems, "Media items retrieved successfully")
}
