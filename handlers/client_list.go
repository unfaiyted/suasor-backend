// handlers/client_media_playlist.go
package handlers

import (
	"strconv"
	mediatypes "suasor/clients/media/types"
	clienttypes "suasor/clients/types"
	"suasor/services"
	_ "suasor/types/models"
	"suasor/types/responses"
	"suasor/utils"
	"suasor/utils/logger"

	"github.com/gin-gonic/gin"
)

type ClientListHandler[T clienttypes.ClientMediaConfig, U mediatypes.ListData] interface {
	CoreListHandler[U]

	GetListByID(c *gin.Context)
	GetListsByGenre(c *gin.Context)
	GetListsByYear(c *gin.Context)
	GetListsByActor(c *gin.Context)
	GetListsByCreator(c *gin.Context)
	GetListsByRating(c *gin.Context)
	GetLatestListsByAdded(c *gin.Context)
	GetPopularLists(c *gin.Context)
	GetTopRatedLists(c *gin.Context)
	SearchLists(c *gin.Context)
}

// clientListHandler handles playlist-related operations for media clients
type clientListHandler[T clienttypes.ClientMediaConfig, U mediatypes.ListData] struct {
	CoreListHandler[U]
	listService services.ClientListService[T, U]
}

// NewclientListHandler creates a new media client playlist handler
func NewClientListHandler[T clienttypes.ClientMediaConfig, U mediatypes.ListData](
	coreHandler CoreListHandler[U],
	listService services.ClientListService[T, U]) *clientListHandler[T, U] {
	return &clientListHandler[T, U]{
		CoreListHandler: coreHandler,
		listService:     listService,
	}
}

// GetListByID godoc
// @Summary Get list by ID
// @Description Retrieves a specific list from the client by ID
// @Tags lists, clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param listID path string true "List ID"
// @Param listType path string true "List type (e.g. 'playlist', 'collection')"
// @Success 200 {object} responses.APIResponse[models.MediaItem[mediatypes.Playlist]] "List retrieved"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid client ID"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/client/{clientID}/{listType}/{listID} [get]
func (h *clientListHandler[T, U]) GetListByID(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting list by ID")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access playlist without authentication")
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

	playlistID := c.Param("id")

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Msg("Retrieving playlist by ID")

	playlist, err := h.listService.GetClientList(ctx, uid, playlistID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Str("playlistID", playlistID).
			Msg("Failed to retrieve playlist")
		responses.RespondInternalError(c, err, "Failed to retrieve playlist")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Msg("List retrieved successfully")
	responses.RespondOK(c, playlist, "List retrieved successfully")
}

// GetLists godoc
// @Summary Get all playlists
// @Description Retrieves all playlists from the client
// @Tags lists, clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param listType path string true "List type (e.g. 'playlist', 'collection')"
// @Param limit query int false "Maximum number of playlists to return"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Playlist]] "Lists retrieved"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid client ID"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/client/{clientID}/{listType} [get]
func (h *clientListHandler[T, U]) GetLists(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting all playlists")

	// Get authenticated user ID
	uid, _ := checkUserAccess(c)
	limit := utils.GetLimit(c, 10, 100, true)

	log.Info().
		Uint64("userID", uid).
		Int("count", limit).
		Msg("Retrieving playlists")

	playlists, err := h.listService.GetClientLists(ctx, uid, limit)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Int("count", limit).
			Msg("Failed to retrieve playlists")
		responses.RespondInternalError(c, err, "Failed to retrieve playlists")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", limit).
		Int("playlistsReturned", len(playlists)).
		Msg("Lists retrieved successfully")
	responses.RespondOK(c, playlists, "Lists retrieved successfully")
}

// CreateList godoc
// @Summary Create a new playlist
// @Description Creates a new playlist on the specified client
// @Tags lists, clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param listType path string true "List type (e.g. 'playlist', 'collection')"
// @Param list body object true "List creation data"
// @Success 201 {object} responses.APIResponse[models.MediaItem[mediatypes.Playlist]] "List created"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/client/{clientID}/{listType} [post]
func (h *clientListHandler[T, U]) CreateList(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Creating list")

	// Get authenticated user ID
	uid, _ := checkUserAccess(c)
	// Parse client ID from URL
	clientID, _ := checkClientID(c)

	// Parse request body
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request body")
		responses.RespondBadRequest(c, err, "Invalid request body")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("name", req.Name).
		Msg("Creating playlist")

	list, err := h.listService.CreateClientList(ctx, clientID, req.Name, req.Description)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Str("name", req.Name).
			Msg("Failed to create playlist")
		responses.RespondInternalError(c, err, "Failed to create playlist")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("name", req.Name).
		Msg("List created successfully")
	responses.RespondCreated(c, list, "List created successfully")
}

// UpdateList godoc
// @Summary Update a playlist
// @Description Updates an existing playlist on the specified client
// @Tags lists, clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param listType path string true "List type (e.g. 'playlist', 'collection')"
// @Param listID path string true "List ID"
// @Param list body object true "List update data"
// @Success 200 {object} responses.APIResponse[models.MediaItem[mediatypes.Playlist]] "List updated"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/client/{clientID}/{listType}/{listID} [put]
func (h *clientListHandler[T, U]) UpdateList(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Updating playlist")

	// Get authenticated user ID
	uid, _ := checkUserAccess(c)
	// Parse client ID from URL
	clientID, _ := checkClientID(c)
	listID := c.Param("listID")

	// Parse request body
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request body")
		responses.RespondBadRequest(c, err, "Invalid request body")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("playlistID", listID).
		Str("name", req.Name).
		Msg("Updating playlist")

	list, err := h.listService.UpdateClientList(ctx, clientID, listID, req.Name, req.Description)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Str("playlistID", listID).
			Msg("Failed to update playlist")
		responses.RespondInternalError(c, err, "Failed to update playlist")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("playlistID", listID).
		Msg("List updated successfully")
	responses.RespondOK(c, list, "List updated successfully")
}

// DeleteList godoc
// @Summary Delete a playlist
// @Description Deletes a playlist from the specified client
// @Tags lists, clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param listType path string true "List type (e.g. 'playlist', 'collection')"
// @Param listID path string true "List ID"
// @Success 200 {object} responses.APIResponse[string] "List deleted"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/client/{clientID}/{listType}/{listID} [delete]
func (h *clientListHandler[T, U]) DeleteList(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Deleting list")

	// Get authenticated user ID
	uid, _ := checkUserAccess(c)

	clientID, _ := checkClientID(c)

	listID := c.Param("id")

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("playlistID", listID).
		Msg("Deleting list")

	err := h.listService.DeleteClientList(ctx, clientID, listID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Str("playlistID", listID).
			Msg("Failed to delete list")
		responses.RespondInternalError(c, err, "Failed to delete list")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("listID", listID).
		Msg("List deleted successfully")
	responses.RespondOK(c, "List deleted successfully", "List deleted successfully")
}

// AddItemToList godoc
// @Summary Add an item to a playlist
// @Description Adds a media item to an existing playlist
// @Tags lists, clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param listType path string true "List type (e.g. 'playlist', 'collection')"
// @Param listID path string true "List ID"
// @Param item body object true "Item to add"
// @Success 200 {object} responses.APIResponse[string] "Item added to playlist"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/client/{clientID}/{listType}/{listID}/items [post]
func (h *clientListHandler[T, U]) AddItemToList(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Adding item to playlist")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to modify playlist without authentication")
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

	playlistID := c.Param("id")

	// Parse request body
	var req struct {
		ItemID string `json:"itemId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request body")
		responses.RespondBadRequest(c, err, "Invalid request body")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Str("itemID", req.ItemID).
		Msg("Adding item to playlist")

	err = h.listService.AddClientItem(ctx, clientID, playlistID, req.ItemID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Str("playlistID", playlistID).
			Str("itemID", req.ItemID).
			Msg("Failed to add item to playlist")
		responses.RespondInternalError(c, err, "Failed to add item to playlist")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Str("itemID", req.ItemID).
		Msg("Item added to playlist successfully")
	responses.RespondOK(c, "Item added to playlist", "Item added to playlist successfully")
}

// RemoveItemFromList godoc
// @Summary Remove an item from a playlist
// @Description Removes a media item from an existing playlist
// @Tags lists, clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param playlistID path string true "List ID"
// @Param itemID path string true "Item ID to remove"
// @Success 200 {object} responses.APIResponse[string] "Item removed from playlist"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/client/{clientID}/{listType}/{playlistID}/items/{itemID} [delete]
func (h *clientListHandler[T, U]) RemoveItemFromList(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Removing item from playlist")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to modify playlist without authentication")
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

	playlistID := c.Param("id")
	itemID := c.Param("itemID")

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Str("itemID", itemID).
		Msg("Removing item from playlist")

	err = h.listService.RemoveClientItem(ctx, clientID, playlistID, itemID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Str("playlistID", playlistID).
			Str("itemID", itemID).
			Msg("Failed to remove item from playlist")
		responses.RespondInternalError(c, err, "Failed to remove item from playlist")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Str("itemID", itemID).
		Msg("Item removed from playlist successfully")
	responses.RespondOK(c, "Item removed from playlist", "Item removed from playlist successfully")
}

// SearchLists godoc
// @Summary Search playlists
// @Description Searches for playlists matching the given query
// @Tags lists, clients
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param listType path string true "List type (e.g. 'playlist', 'collection')"
// @Param q query string true "Search query"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Playlist]] "Lists found"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/client/{clientID}/{listType}/search [get]
func (h *clientListHandler[T, U]) SearchLists(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Searching playlists")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to search playlists without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)
	query := c.Query("q")
	if query == "" {
		log.Warn().Uint64("userID", uid).Msg("Empty search query provided")
		responses.RespondBadRequest(c, nil, "Search query is required")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Str("query", query).
		Msg("Searching playlists")

		// Parse client ID from URL
	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("clientID", c.Param("clientID")).Msg("Invalid client ID format")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}
	options := mediatypes.QueryOptions{
		Query: query,
	}

	playlists, err := h.listService.SearchClientLists(ctx, clientID, options)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Str("query", query).
			Msg("Failed to search playlists")
		responses.RespondInternalError(c, err, "Failed to search playlists")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Str("query", query).
		Int("resultsCount", len(playlists)).
		Msg("List search completed successfully")
	responses.RespondOK(c, playlists, "Lists retrieved successfully")
}
