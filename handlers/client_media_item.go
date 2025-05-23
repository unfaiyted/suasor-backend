// handlers/client_media_item.go
package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"suasor/utils"

	"github.com/gin-gonic/gin"

	"suasor/clients/media/types"
	clientTypes "suasor/clients/types"
	"suasor/services"
	"suasor/types/models"
	"suasor/types/requests"
	"suasor/types/responses"
	"suasor/utils/logger"
)

type ClientMediaItemHandler[T clientTypes.ClientMediaConfig, U types.MediaData] interface {
	UserMediaItemHandler[U]

	GetAllClientItems(c *gin.Context)
	GetClientItemByItemID(c *gin.Context)

	SearchClient(c *gin.Context)
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
//
// @Summary		Create a new media item associated with a client
// @Description	Creates a new media item in the database with client association
// @Tags			media, clients
// @Accept			json
// @Produce		json
// @Param			mediaItem	body		models.MediaItem[any]							true	"Media item data with type, client info, and type-specific data"
// @Success		201			{object}	responses.APIResponse[models.MediaItem[types.MediaData]]	"Media item created successfully"
// @Failure		400			{object}	responses.ErrorResponse[any]					"Invalid request"
// @Failure		500			{object}	responses.ErrorResponse[any]					"Server error"
// @Router			/client/{clientID}/media/{mediaType} [post]
func (h *clientMediaItemHandler[T, U]) CreateClientItem(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	var req requests.ClientMediaItemCreateRequest[U]

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid request body for CreateMediaItem")
		responses.RespondBadRequest(c, err, "Invalid request format")
		return
	}

	mediaType := types.MediaType(req.Type)
	log.Debug().
		Str("mediaType", string(mediaType)).
		Uint64("clientID", req.ClientID).
		Str("clientType", string(req.ClientType)).
		Msg("Creating client media item")

	// Create the media item
	mediaItem := models.NewMediaItem(req.Data)

	// Set client info
	mediaItem.SetClientInfo(req.ClientID, req.ClientType, req.ExternalID)

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

// UpdateClientMediaItem godoc
//
// @Summary		Update an existing client media item
// @Description	Updates a client media item in the database by ID
// @Tags			media, clients
// @Accept			json
// @Produce		json
// @Param			clientItemID		path		string											true	"Media item ID"
// @Param			mediaItem	body		models.MediaItem[any]							true	"Media item data to update"
// @Success		200			{object}	responses.APIResponse[models.MediaItem[types.MediaData]]	"Media item updated successfully"
// @Failure		400			{object}	responses.ErrorResponse[any]					"Invalid request"
// @Failure		404			{object}	responses.ErrorResponse[any]					"Media item not found"
// @Failure		500			{object}	responses.ErrorResponse[any]					"Server error"
// @Router			/client/{clientID}/media/{clientItemID} [put]
func (h *clientMediaItemHandler[T, U]) UpdateClientItem(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	id := c.Param("clientItemID")

	var req requests.ClientMediaItemUpdateRequest[U]

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Str("id", id).Msg("Invalid request body for UpdateMediaItem")
		responses.RespondValidationError(c, err)
		return
	}

	log.Debug().
		Str("id", id).
		Str("mediaType", string(req.Type)).
		Uint64("clientID", req.ClientID).
		Str("clientType", string(req.ClientType)).
		Msg("Updating client media item")

	mediaItem := models.NewMediaItem(req.Data)
	mediaItem.SetClientInfo(req.ClientID, req.ClientType, req.ExternalID)
	// Set client info

	// Only add external ID if provided
	if req.ExternalID != "" {
		mediaItem.AddExternalID("client", req.ExternalID)
	}

	result, err := h.clientService.Update(ctx, mediaItem)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to update client media item")
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
//
//		@Summary		Get media items by client
//		@Description	Retrieves all media items for a specific client
//		@Tags			media, clients
//		@Accept			json
//		@Produce		json
//		@Param			clientID  path		int												true	"Client ID"
//		@Param			mediaType path		string											false	"Media type filter"
//	 @Param			limit     query		int												false	"Maximum number of items to return (default 20)"
//	 @Param			offset    query		int												false	"Offset for pagination (default 0)"
//		@Success		200			{object}	responses.APIResponse[responses.MediaItemList[types.MediaData]]	"Media items retrieved successfully"
//		@Failure		400			{object}	responses.ErrorResponse[any]					"Invalid client ID"
//		@Failure		500			{object}	responses.ErrorResponse[any]					"Server error"
//		@Router			/client/{clientID}/media/{mediaType} [get]
func (h *clientMediaItemHandler[T, U]) GetAllClientItems(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	clientID, exists := checkClientID(c)
	if !exists {
		return
	}
	var zero U
	mediaType := types.GetMediaTypeFromTypeName(zero)

	limit := utils.GetLimit(c, 20, 100, false)
	offset := utils.GetOffset(c, 0)

	log.Debug().
		Uint64("clientID", clientID).
		Str("mediaType", string(mediaType)).
		Int("limit", limit).
		Int("offset", offset).
		Msg("Getting media items by client")

	var items []*models.MediaItem[U]

	items, err := h.clientService.GetByClientID(ctx, clientID, mediaType, limit, offset)

	if err != nil {
		log.Error().Err(err).
			Uint64("clientID", clientID).
			Str("mediaType", string(mediaType)).
			Msg("Failed to retrieve media items by client")
		responses.RespondInternalError(c, err, "Failed to retrieve media items by client")
		return
	}

	log.Info().
		Uint64("clientID", clientID).
		Int("count", len(items)).
		Msg("Media items retrieved by client successfully")

	responses.RespondOK(c, items, "Media items retrieved successfully")
}

// GetMediaItemByClientItemID godoc
//
// @Summary		Get media item by client-specific ID
// @Description	Retrieves a media item using its client-specific ID
// @Tags			media, clients
// @Accept			json
// @Produce		json
// @Param			clientID	path		int												true	"Client ID"
// @Param			mediaType	path		string											true	"Media type"
// @Param			clientItemID		path		string											true	"Client-specific item ID"
// @Success		200			{object}	responses.APIResponse[models.MediaItem[types.MediaData]]	"Media item retrieved successfully"
// @Failure		400			{object}	responses.ErrorResponse[any]					"Invalid request"
// @Failure		404			{object}	responses.ErrorResponse[any]					"Media item not found"
// @Failure		500			{object}	responses.ErrorResponse[any]					"Server error"
// @Router			/client/{clientID}/media/{mediaType}/{clientItemID} [get]
func (h *clientMediaItemHandler[T, U]) GetClientItemByItemID(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	clientID, exists := checkClientID(c)
	if !exists {
		return
	}
	clientItemID := c.Param("clientItemID")
	if clientItemID == "" {
		log.Warn().Msg("Client item ID is required")
		responses.RespondBadRequest(c, nil, "Client item ID is required")
		return
	}

	log.Debug().
		Uint64("clientID", clientID).
		Str("clientItemID", clientItemID).
		Msg("Getting media item by client item ID")

	item, err := h.clientService.GetByClientItemID(ctx, clientID, clientItemID)
	if err != nil {
		log.Error().Err(err).
			Uint64("clientID", clientID).
			Str("clientItemID", clientItemID).
			Msg("Failed to retrieve media item by client item ID")
		responses.RespondNotFound(c, err, "Media item not found")
		return
	}

	log.Info().
		Uint64("clientID", clientID).
		Str("clientItemID", clientItemID).
		Uint64("id", item.ID).
		Msg("Media item retrieved by client item ID successfully")

	responses.RespondOK(c, item, "Media item retrieved successfully")
}

// GetMediaItemsByMultipleClients godoc
//
// @Summary		Get media items from multiple clients
// @Description	Retrieves media items associated with any of the specified clients
// @Tags			media, clients
// @Accept			json
// @Produce		json
// @Param			clientIDs	query		string											true	"Comma-separated list of client IDs"
// @Success		200			{object}	responses.APIResponse[responses.MediaItemList[types.MediaData]]	"Media items retrieved successfully"
// @Failure		400			{object}	responses.ErrorResponse[any]					"Invalid request"
// @Failure		500			{object}	responses.ErrorResponse[any]					"Server error"
// @Router			/client/media/multi [get]
func (h *clientMediaItemHandler[T, U]) GetItemsByMultipleClients(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	clientIDsStr := c.Query("clientIDs")
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
			log.Warn().Err(err).Str("clientID", idStr).Msg("Invalid client ID format")
			responses.RespondBadRequest(c, err, "Invalid client ID format")
			return
		}
		clientIDs = append(clientIDs, id)
	}

	log.Debug().
		Interface("clientIDs", clientIDs).
		Msg("Getting media items by multiple clients")

	items, err := h.clientService.GetByMultipleClients(ctx, clientIDs)
	if err != nil {
		log.Error().Err(err).
			Interface("clientIDs", clientIDs).
			Msg("Failed to retrieve media items by multiple clients")
		responses.RespondInternalError(c, err, "Failed to retrieve media items")
		return
	}

	log.Info().
		Interface("clientIDs", clientIDs).
		Int("count", len(items)).
		Msg("Media items retrieved by multiple clients successfully")

	responses.RespondMediaItemListOK(c, items, "Media items retrieved successfully")
}

// SearchAcrossClients godoc
//
// @Summary		Search for media items across multiple clients
// @Description	Searches for media items across multiple clients based on query parameters
// @Tags			media, clients
// @Accept			json
// @Produce		json
// @Param			q			query		string														true	"Search query"
// @Param			clientIDs	query		string														true	"Comma-separated list of client IDs"
// @Param			type		query		string														false	"Media type filter"
// @Success		200			{object}	responses.APIResponse[map[string]models.MediaItem[types.MediaData]]	"Media items retrieved successfully"
// @Failure		400			{object}	responses.ErrorResponse[any]								"Invalid request"
// @Failure		500			{object}	responses.ErrorResponse[any]								"Server error"
// @Router			/client/media/search [get]
func (h *clientMediaItemHandler[T, U]) SearchAcrossClients(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	query := c.Query("q")
	if query == "" {
		log.Warn().Msg("Search query is required")
		responses.RespondBadRequest(c, nil, "Search query is required")
		return
	}

	clientIDsStr := c.Query("clientIDs")
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
			log.Warn().Err(err).Str("clientID", idStr).Msg("Invalid client ID format")
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
		Interface("clientIDs", clientIDs).
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
			Interface("clientIDs", clientIDs).
			Msg("Failed to search media items across clients")
		responses.RespondInternalError(c, err, "Failed to search media items")
		return
	}

	log.Info().
		Str("query", query).
		Interface("clientIDs", clientIDs).
		Int("clientCount", len(results)).
		Msg("Media items search across clients completed successfully")

	responses.RespondOK(c, results, "Media items retrieved successfully")
}

// SyncItemBetweenClients godoc
//
// @Summary		Sync a media item between clients
// @Description	Creates or updates a mapping between a media item and a target client
// @Tags			media, clients
// @Accept			json
// @Produce		json
// @Param			syncRequest	body		object							true	"Sync request with source and target client info"
// @Success		200			{object}	responses.APIResponse[types.MediaData]		"Item synced successfully"
// @Failure		400			{object}	responses.ErrorResponse[any]	"Invalid request"
// @Failure		404			{object}	responses.ErrorResponse[any]	"Media item not found"
// @Failure		500			{object}	responses.ErrorResponse[any]	"Server error"
// @Router			/client/media/sync [post]
func (h *clientMediaItemHandler[T, U]) SyncItemBetweenClients(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	var req struct {
		ItemID         uint64 `json:"clientItemID" binding:"required"`
		SourceClientID uint64 `json:"sourceClientID" binding:"required"`
		TargetClientID uint64 `json:"targetClientID" binding:"required"`
		TargetItemID   string `json:"targetItemID" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid request body for SyncItemBetweenClients")
		responses.RespondBadRequest(c, err, "Invalid request format")
		return
	}

	log.Debug().
		Uint64("clientItemID", req.ItemID).
		Uint64("sourceClientID", req.SourceClientID).
		Uint64("targetClientID", req.TargetClientID).
		Str("targetItemID", req.TargetItemID).
		Msg("Syncing item between clients")

	err := h.clientService.SyncItemBetweenClients(ctx, req.ItemID, req.SourceClientID, req.TargetClientID, req.TargetItemID)
	if err != nil {
		log.Error().Err(err).
			Uint64("clientItemID", req.ItemID).
			Uint64("sourceClientID", req.SourceClientID).
			Uint64("targetClientID", req.TargetClientID).
			Msg("Failed to sync item between clients")
		responses.RespondInternalError(c, err, "Failed to sync item between clients")
		return
	}

	log.Info().
		Uint64("clientItemID", req.ItemID).
		Uint64("sourceClientID", req.SourceClientID).
		Uint64("targetClientID", req.TargetClientID).
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

	userID, err := strconv.ParseUint(c.Query("userID"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userID", c.Query("userID")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	log.Debug().
		Int("year", year).
		Uint64("userID", userID).
		Int("limit", limit).
		Msg("Getting media items by year")

	// Filter by year using the search functionality
	options := &types.QueryOptions{
		Year:  year,
		Limit: limit,
	}

	items, err := h.clientService.Search(ctx, *options)
	if err != nil {
		log.Error().Err(err).Int("year", year).Uint64("userID", userID).Msg("Failed to retrieve media items")
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

	log.Info().Int("year", year).Uint64("userID", userID).Int("count", len(filtered)).Msg("Media items retrieved by year")
	responses.RespondOK(c, filtered, "Media items retrieved successfully")
}

// DeleteClientItem godoc
//
// @Summary		Delete a media item from a client
// @Description	Deletes a media item from a client
// @Tags			media, clients
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			clientID	path		int								true	"Client ID"
// @Param			clientItemID		path		string							true	"Item ID"
// @Success		200			{object}	responses.SuccessResponse									"Item deleted"
// @Failure		400			{object}	responses.ErrorResponse[error]	"Invalid request"
// @Failure		401			{object}	responses.ErrorResponse[error]	"Unauthorized"
// @Failure		500			{object}	responses.ErrorResponse[error]	"Server error"
// @Router			/client/{clientID}/media/item/{clientItemID} [delete]
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

	clientItemID := c.Param("clientItemID")

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("clientItemID", clientItemID).
		Msg("Deleting client item")

	err = h.clientService.DeleteClientItem(ctx, clientID, clientItemID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Str("clientItemID", clientItemID).
			Msg("Failed to delete client item")
		responses.RespondInternalError(c, err, "Failed to delete client item")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("clientItemID", clientItemID).
		Msg("Client item deleted successfully")
	responses.RespondOK(c, "Item deleted successfully", "Item deleted successfully")
}

// GetByClientItemID godoc
//
// @Summary		Get media item by client-specific ID
// @Description	Retrieves a media item using its client-specific ID
// @Tags			media, clients
// @Accept			json
// @Produce		json
// @Param			clientID		path		int												true	"Client ID"
// @Param			clientItemID	path		string											true	"Client-specific item ID"
// @Param			mediaType		path		string											true	"Media type"
// @Success		200				{object}	responses.APIResponse[models.MediaItem[types.MediaData]]	"Media item retrieved successfully"
// @Failure		400				{object}	responses.ErrorResponse[any]					"Invalid request"
// @Failure		404				{object}	responses.ErrorResponse[any]					"Media item not found"
// @Failure		500				{object}	responses.ErrorResponse[any]					"Server error"
// @Router			/client/{clientID}/media/{mediaType}/{clientItemID} [get]
func (h *clientMediaItemHandler[T, U]) GetByClientItemID(c *gin.Context) {

}

// SearchClient godoc
//
// @Summary		Search for media items in a specific client
// @Description	Searches for media items in a specific client based on query parameters
// @Tags			media, clients
// @Accept			json
// @Produce		json
// @Param			q				 query		string														false	"Search query"
// @Param 		options  body		  types.QueryOptions											false	"Search options"
// @Param			clientID path					string														true	"Client ID"
// @Param			mediaType		 query		string														false	"Media type filter"
// @Success		200			{object}	responses.APIResponse[responses.MediaItemList[types.MediaData]]	"Media items retrieved successfully"
// @Failure		400			{object}	responses.ErrorResponse[any]								"Invalid request"
// @Failure		500			{object}	responses.ErrorResponse[any]								"Server error"
// @Router			/client/{clientID}/media/{mediaType}/search [get]
func (h *clientMediaItemHandler[T, U]) SearchClient(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	clientID, exists := checkClientID(c)
	if !exists {
		return
	}

	var options types.QueryOptions
	if err := c.ShouldBindJSON(&options); err != nil {
		// JSON binding failed, try to get parameters from query string
		log.Debug().Err(err).Msg("JSON binding failed, using query parameters instead")
	}

	// Get query from URL parameter if not in JSON body
	query := c.Query("q")
	if query == "" && options.Query == "" {
		log.Warn().Msg("Search query is required")
		responses.RespondBadRequest(c, nil, "Search query is required")
		return
	}
	if query != "" {
		options.Query = query
	}
	mediaType := types.GetMediaType[U]()

	// Set the client ID in the options
	options.WithClientID(clientID)
	options.WithMediaType(mediaType)

	// if there is no limit set, set it to 20
	if options.Limit == 0 {
		options.Limit = 20
	}
	// offset is set to 0 by default
	if options.Offset == 0 {
		options.Offset = 0
	}

	log.Debug().
		Uint64("clientID", clientID).
		Str("query", options.Query).
		Msg("Searching for media items in client")

	// Perform the search
	results, err := h.clientService.SearchClient(ctx, clientID, options)
	if err != nil {
		log.Error().Err(err).
			Uint64("clientID", clientID).
			Str("query", options.Query).
			Msg("Failed to search media items in client")
		responses.RespondInternalError(c, err, "Failed to search media items")
		return
	}

	log.Info().
		Uint64("clientID", clientID).
		Str("query", options.Query).
		Int("resultCount", len(results)).
		Msg("Media items search in client completed successfully")

	responses.RespondMediaItemListOK(c, results, "Media items retrieved successfully")
}
