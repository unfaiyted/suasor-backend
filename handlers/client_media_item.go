// handlers/client_media_item.go
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"suasor/clients/media/types"
	clientTypes "suasor/clients/types"
	"suasor/services"
	"suasor/types/models"
	"suasor/types/responses"
	"suasor/utils/logger"
)

type ClientMediaItemHandler[T clientTypes.ClientMediaConfig, U types.MediaData] interface {
	UserMediaItemHandler[U]

	GetAllClientItems(c *gin.Context)
	GetClientItemByItemID(c *gin.Context)
	// RecordClientPlay(c *gin.Context)
	// GetClientPlaybackState(c *gin.Context)
	// UpdateClientPlaybackState(c *gin.Context)
	DeleteClientItem(c *gin.Context)
	// SyncClientItem(c *gin.Context)
}

// ClientMediaItemHandler handles operations on media items associated with external clients
// It extends UserMediaItemHandler to inherit both core and user-level functionality
// This is the third prong in the three-pronged architecture
type clientMediaItemHandler[T clientTypes.ClientMediaConfig, U types.MediaData] struct {
	UserMediaItemHandler[U] // Embed the user handler
	clientService           services.ClientMediaItemService[T, U]
}

// NewClientMediaItemHandler creates a new client media item handler
func NewClientMediaItemHandler[T clientTypes.ClientMediaConfig, U types.MediaData](
	userHandler UserMediaItemHandler[U],
	clientService services.ClientMediaItemService[T, U],
) ClientMediaItemHandler[T, U] {
	return &clientMediaItemHandler[T, U]{
		UserMediaItemHandler: userHandler,
		clientService:        clientService,
	}
}

// CreateMediaItem godoc
// @Summary Create a new media item associated with a client
// @Description Creates a new media item in the database with client association
// @Tags client-media
// @Accept json
// @Produce json
// @Param mediaItem body object true "Media item data with type, client info, and type-specific data"
// @Success 201 {object} responses.APIResponse[models.MediaItem[any]] "Media item created successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /item/media [post]
func (h *clientMediaItemHandler[T, U]) CreateClientItem(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

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
		Msg("Creating client media item")

	// Create the media item
	var mediaData U
	if err := json.Unmarshal(req.Data, &mediaData); err != nil {
		log.Warn().Err(err).Msg("Failed to unmarshal media data")
		responses.RespondBadRequest(c, err, "Invalid media data format")
		return
	}

	mediaItem := models.NewMediaItem[U](mediaType, mediaData)

	// Set client info
	mediaItem.SetClientInfo(req.ClientID, clientTypes.ClientMediaType(req.ClientType), req.ExternalID)

	// Only add external ID if provided
	if req.ExternalID != "" {
		mediaItem.AddExternalID("client", req.ExternalID)
	}

	result, err := h.clientService.Create(ctx, mediaItem)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create client media item")
		responses.RespondInternalError(c, err, "Failed to create media item")
		return
	}

	log.Info().
		Uint64("id", result.ID).
		Str("type", string(result.Type)).
		Msg("Client media item created successfully")

	responses.RespondCreated(c, result, "Media item created successfully")
}

// UpdateMediaItem godoc
// @Summary Update an existing client media item
// @Description Updates a client media item in the database by ID
// @Tags client-media
// @Accept json
// @Produce json
// @Param id path int true "Media item ID"
// @Param mediaItem body object true "Media item data to update"
// @Success 200 {object} responses.APIResponse[models.MediaItem[any]] "Media item updated successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Media item not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /item/media/{id} [put]
func (h *clientMediaItemHandler[T, U]) UpdateClientItem(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

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
		Msg("Updating client media item")

	// Update the media item
	mediaType := types.MediaType(req.Type)
	var mediaData U
	if err := json.Unmarshal(req.Data, &mediaData); err != nil {
		log.Warn().Err(err).Uint64("id", id).Msg("Failed to unmarshal media data")
		responses.RespondBadRequest(c, err, "Invalid media data format")
		return
	}
	mediaItem := models.NewMediaItem[U](mediaType, mediaData)
	mediaItem.ID = id
	mediaItem.SetClientInfo(req.ClientID, clientTypes.ClientMediaType(req.ClientType), req.ExternalID)
	// Set client info

	// Only add external ID if provided
	if req.ExternalID != "" {
		mediaItem.AddExternalID("client", req.ExternalID)
	}

	result, err := h.clientService.Update(ctx, mediaItem)
	if err != nil {
		log.Error().Err(err).Uint64("id", id).Msg("Failed to update client media item")
		responses.RespondInternalError(c, err, "Failed to update media item")
		return
	}

	log.Info().
		Uint64("id", result.ID).
		Str("type", string(result.Type)).
		Msg("Client media item updated successfully")

	responses.RespondOK(c, result, "Media item updated successfully")
}

// GetMediaItemsByClient godoc
// @Summary Get media items by client
// @Description Retrieves all media items for a specific client
// @Tags client-media
// @Accept json
// @Produce json
// @Param clientId path int true "Client ID"
// @Param type query string false "Media type filter"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[any]] "Media items retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid client ID"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /client/{clientId}/media [get]
func (h *clientMediaItemHandler[T, U]) GetAllClientItems(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	clientID, err := strconv.ParseUint(c.Param("clientId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("clientId", c.Param("clientId")).Msg("Invalid client ID")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	typeParam := c.Query("type") // Optional media type filter

	log.Debug().
		Uint64("clientId", clientID).
		Str("mediaType", typeParam).
		Msg("Getting media items by client")

	var items []*models.MediaItem[U]

	if typeParam != "" {
		// If media type is specified, get media items by client and type
		// TODO: Implement this
		// mediaType := types.MediaType(typeParam)
		// items, err = h.clientService.GetByClientID(ctx, clientID)
	} else {
		// Otherwise, get all media items for the client
		items, err = h.clientService.GetByClientID(ctx, clientID)
	}

	if err != nil {
		log.Error().Err(err).
			Uint64("clientId", clientID).
			Str("type", typeParam).
			Msg("Failed to retrieve media items by client")
		responses.RespondInternalError(c, err, "Failed to retrieve media items by client")
		return
	}

	log.Info().
		Uint64("clientId", clientID).
		Int("count", len(items)).
		Msg("Media items retrieved by client successfully")

	responses.RespondOK(c, items, "Media items retrieved successfully")
}

// GetMediaItemByClientItemID godoc
// @Summary Get media item by client-specific ID
// @Description Retrieves a media item using its client-specific ID
// @Tags client-media
// @Accept json
// @Produce json
// @Param clientId path int true "Client ID"
// @Param itemId path string true "Client-specific item ID"
// @Success 200 {object} responses.APIResponse[models.MediaItem[any]] "Media item retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Media item not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /client/{clientId}/media/item/{itemId} [get]
func (h *clientMediaItemHandler[T, U]) GetClientItemByItemID(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	clientID, err := strconv.ParseUint(c.Param("clientId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("clientId", c.Param("clientId")).Msg("Invalid client ID")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	itemID := c.Param("itemId")
	if itemID == "" {
		log.Warn().Msg("Client item ID is required")
		responses.RespondBadRequest(c, nil, "Client item ID is required")
		return
	}

	log.Debug().
		Uint64("clientId", clientID).
		Str("itemId", itemID).
		Msg("Getting media item by client item ID")

	item, err := h.clientService.GetByClientItemID(ctx, itemID, clientID)
	if err != nil {
		log.Error().Err(err).
			Uint64("clientId", clientID).
			Str("itemId", itemID).
			Msg("Failed to retrieve media item by client item ID")
		responses.RespondNotFound(c, err, "Media item not found")
		return
	}

	log.Info().
		Uint64("clientId", clientID).
		Str("itemId", itemID).
		Uint64("id", item.ID).
		Msg("Media item retrieved by client item ID successfully")

	responses.RespondOK(c, item, "Media item retrieved successfully")
}

// GetMediaItemsByMultipleClients godoc
// @Summary Get media items from multiple clients
// @Description Retrieves media items associated with any of the specified clients
// @Tags client-media
// @Accept json
// @Produce json
// @Param clientIds query string true "Comma-separated list of client IDs"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[any]] "Media items retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /client/media/multi [get]
func (h *clientMediaItemHandler[T, U]) GetItemsByMultipleClients(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	clientIDsStr := c.Query("clientIds")
	if clientIDsStr == "" {
		log.Warn().Msg("Client IDs parameter is required")
		responses.RespondBadRequest(c, nil, "Client IDs parameter is required")
		return
	}

	// Parse comma-separated list of client IDs
	clientIDStrs := strings.Split(clientIDsStr, ",")
	clientIDs := make([]uint64, 0, len(clientIDStrs))

	for _, idStr := range clientIDStrs {
		id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 64)
		if err != nil {
			log.Warn().Err(err).Str("clientId", idStr).Msg("Invalid client ID format")
			responses.RespondBadRequest(c, err, "Invalid client ID format")
			return
		}
		clientIDs = append(clientIDs, id)
	}

	log.Debug().
		Interface("clientIds", clientIDs).
		Msg("Getting media items by multiple clients")

	items, err := h.clientService.GetByMultipleClients(ctx, clientIDs)
	if err != nil {
		log.Error().Err(err).
			Interface("clientIds", clientIDs).
			Msg("Failed to retrieve media items by multiple clients")
		responses.RespondInternalError(c, err, "Failed to retrieve media items")
		return
	}

	log.Info().
		Interface("clientIds", clientIDs).
		Int("count", len(items)).
		Msg("Media items retrieved by multiple clients successfully")

	responses.RespondOK(c, items, "Media items retrieved successfully")
}

// SearchAcrossClients godoc
// @Summary Search for media items across multiple clients
// @Description Searches for media items across multiple clients based on query parameters
// @Tags client-media
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param clientIds query string true "Comma-separated list of client IDs"
// @Param type query string false "Media type filter"
// @Success 200 {object} responses.APIResponse[map[string][]models.MediaItem[any]] "Media items retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /client/media/search [get]
func (h *clientMediaItemHandler[T, U]) SearchAcrossClients(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	query := c.Query("q")
	if query == "" {
		log.Warn().Msg("Search query is required")
		responses.RespondBadRequest(c, nil, "Search query is required")
		return
	}

	clientIDsStr := c.Query("clientIds")
	if clientIDsStr == "" {
		log.Warn().Msg("Client IDs parameter is required")
		responses.RespondBadRequest(c, nil, "Client IDs parameter is required")
		return
	}

	// Parse comma-separated list of client IDs
	clientIDStrs := strings.Split(clientIDsStr, ",")
	clientIDs := make([]uint64, 0, len(clientIDStrs))

	for _, idStr := range clientIDStrs {
		id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 64)
		if err != nil {
			log.Warn().Err(err).Str("clientId", idStr).Msg("Invalid client ID format")
			responses.RespondBadRequest(c, err, "Invalid client ID format")
			return
		}
		clientIDs = append(clientIDs, id)
	}

	// Get media type from query parameters
	mediaTypeStr := c.Query("type")
	var mediaType types.MediaType
	if mediaTypeStr != "" {
		mediaType = types.MediaType(mediaTypeStr)
	}

	log.Debug().
		Str("query", query).
		Interface("clientIds", clientIDs).
		Str("type", string(mediaType)).
		Msg("Searching media items across clients")

	// Create query options
	options := types.QueryOptions{
		Query:     query,
		MediaType: mediaType,
	}

	// Search across clients
	results, err := h.clientService.SearchAcrossClients(ctx, options, clientIDs)
	if err != nil {
		log.Error().Err(err).
			Str("query", query).
			Interface("clientIds", clientIDs).
			Msg("Failed to search media items across clients")
		responses.RespondInternalError(c, err, "Failed to search media items")
		return
	}

	log.Info().
		Str("query", query).
		Interface("clientIds", clientIDs).
		Int("clientCount", len(results)).
		Msg("Media items search across clients completed successfully")

	responses.RespondOK(c, results, "Media items retrieved successfully")
}

// SyncItemBetweenClients godoc
// @Summary Sync a media item between clients
// @Description Creates or updates a mapping between a media item and a target client
// @Tags client-media
// @Accept json
// @Produce json
// @Param syncRequest body object true "Sync request with source and target client info"
// @Success 200 {object} responses.APIResponse[any] "Item synced successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Media item not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /client/media/sync [post]
func (h *clientMediaItemHandler[T, U]) SyncItemBetweenClients(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	var req struct {
		ItemID         uint64 `json:"itemId" binding:"required"`
		SourceClientID uint64 `json:"sourceClientId" binding:"required"`
		TargetClientID uint64 `json:"targetClientId" binding:"required"`
		TargetItemID   string `json:"targetItemId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid request body for SyncItemBetweenClients")
		responses.RespondBadRequest(c, err, "Invalid request format")
		return
	}

	log.Debug().
		Uint64("itemId", req.ItemID).
		Uint64("sourceClientId", req.SourceClientID).
		Uint64("targetClientId", req.TargetClientID).
		Str("targetItemId", req.TargetItemID).
		Msg("Syncing item between clients")

	err := h.clientService.SyncItemBetweenClients(ctx, req.ItemID, req.SourceClientID, req.TargetClientID, req.TargetItemID)
	if err != nil {
		log.Error().Err(err).
			Uint64("itemId", req.ItemID).
			Uint64("sourceClientId", req.SourceClientID).
			Uint64("targetClientId", req.TargetClientID).
			Msg("Failed to sync item between clients")
		responses.RespondInternalError(c, err, "Failed to sync item between clients")
		return
	}

	log.Info().
		Uint64("itemId", req.ItemID).
		Uint64("sourceClientId", req.SourceClientID).
		Uint64("targetClientId", req.TargetClientID).
		Msg("Item synced between clients successfully")

	responses.RespondOK(c, http.StatusOK, "Item synced successfully")
}

// Additional methods from the original implementation can be kept but should be
// organized to properly inherit from the embedded handlers and avoid duplication.
// Only client-specific methods that truly extend the functionality should remain here.

// GetMediaItemsByGenre - This method could be refactored to use the inherited Search method when appropriate
func (h *clientMediaItemHandler[T, U]) GetMediaItemsByGenre(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	genre := c.Param("genre")
	if genre == "" {
		log.Warn().Msg("Genre is required")
		responses.RespondBadRequest(c, nil, "Genre is required")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	log.Debug().
		Str("genre", genre).
		Int("limit", limit).
		Msg("Getting media items by genre")

	var zero T
	mediaType := types.GetMediaTypeFromTypeName(zero)

	options := types.QueryOptions{
		Genre:     genre,
		MediaType: mediaType,
		Limit:     limit,
	}
	// Filter by genre using the search functionality
	items, err := h.clientService.Search(ctx, options)
	if err != nil {
		log.Error().Err(err).Str("genre", genre).Msg("Failed to retrieve media items")
		responses.RespondInternalError(c, err, "Failed to retrieve media items")
		return
	}

	// Filter for items with the specified genre
	var filtered []*models.MediaItem[U]
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

	log.Info().Str("genre", genre).Int("count", len(filtered)).Msg("Media items retrieved by genre")
	responses.RespondOK(c, filtered, "Media items retrieved successfully")
}

// GetMediaItemsByYear keeps client-specific logic for year filtering
func (h *clientMediaItemHandler[T, U]) GetMediaItemsByYear(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

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
	options := &types.QueryOptions{
		Year:  year,
		Limit: limit,
	}

	items, err := h.clientService.Search(ctx, *options)
	if err != nil {
		log.Error().Err(err).Int("year", year).Uint64("userId", userID).Msg("Failed to retrieve media items")
		responses.RespondInternalError(c, err, "Failed to retrieve media items")
		return
	}

	// Filter for items with the specified year
	var filtered []*models.MediaItem[U]
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

// DeleteClientItem godoc
// @Summary Delete a media item from a client
// @Description Deletes a media item from a client
// @Tags client-media
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param itemID path string true "Item ID"
// @Success 200 {object} responses.APIResponse[string] "Item deleted"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /clients/media/{clientID}/item/{itemID} [delete]
func (h *clientMediaItemHandler[T, U]) DeleteClientItem(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Deleting client item")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to delete client item without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	// Parse client ID from URL
	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("clientID", c.Param("clientID")).Msg("Invalid client ID format")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	itemID := c.Param("itemID")

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("itemID", itemID).
		Msg("Deleting client item")

	err = h.clientService.DeleteClientItem(ctx, clientID, itemID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Str("itemID", itemID).
			Msg("Failed to delete client item")
		responses.RespondInternalError(c, err, "Failed to delete client item")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("itemID", itemID).
		Msg("Client item deleted successfully")
	responses.RespondOK(c, "Item deleted successfully", "Item deleted successfully")
}
