// handlers/playlists.go
package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	mediatypes "suasor/client/media/types"
	"suasor/services"
	"suasor/types/models"
	"suasor/types/responses"
	"suasor/utils"
)

// PlaylistHandler provides handlers for playlist operations
type PlaylistHandler struct {
	*MediaItemHandler[*mediatypes.Playlist]
	service services.MediaClientPlaylistService[any]
}

// NewPlaylistHandler creates a new playlist handler
func NewPlaylistHandler(
	mediaItemService services.MediaItemService[*mediatypes.Playlist],
	playlistService services.MediaClientPlaylistService[any],
) *PlaylistHandler {
	return &PlaylistHandler{
		MediaItemHandler: NewMediaItemHandler(mediaItemService),
		service:          playlistService,
	}
}

// GetPlaylists godoc
// @Summary Get all playlists
// @Description Retrieves all playlists for a user
// @Tags playlists
// @Accept json
// @Produce json
// @Param userId query int true "User ID"
// @Param limit query int false "Maximum number of playlists to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Playlist]] "Playlists retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /playlists [get]
func (h *PlaylistHandler) GetPlaylists(c *gin.Context) {
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
		Uint64("userID", userID).
		Int("limit", limit).
		Msg("Getting playlists")

	playlists, err := h.service.GetPlaylists(ctx, userID, limit)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to retrieve playlists")
		responses.RespondInternalError(c, err, "Failed to retrieve playlists")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Int("count", len(playlists)).
		Msg("Playlists retrieved successfully")

	responses.RespondOK(c, playlists, "Playlists retrieved successfully")
}

// GetPlaylistByID godoc
// @Summary Get a playlist by ID
// @Description Retrieves a specific playlist by ID
// @Tags playlists
// @Accept json
// @Produce json
// @Param id path int true "Playlist ID"
// @Param userId query int true "User ID"
// @Param clientId query int true "Client ID"
// @Success 200 {object} responses.APIResponse[models.MediaItem[mediatypes.Playlist]] "Playlist retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Playlist not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /playlists/{id} [get]
func (h *PlaylistHandler) GetPlaylistByID(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	playlistID := c.Param("id")
	if playlistID == "" {
		log.Warn().Msg("Playlist ID is required")
		responses.RespondBadRequest(c, nil, "Playlist ID is required")
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	clientID, err := strconv.ParseUint(c.Query("clientId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("clientId", c.Query("clientId")).Msg("Invalid client ID")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Msg("Getting playlist by ID")

	playlist, err := h.service.GetPlaylistByID(ctx, userID, clientID, playlistID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("playlistID", playlistID).
			Msg("Failed to retrieve playlist")
		responses.RespondNotFound(c, err, "Playlist not found")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Msg("Playlist retrieved successfully")

	responses.RespondOK(c, playlist, "Playlist retrieved successfully")
}

// GetPlaylistItems godoc
// @Summary Get items in a playlist
// @Description Retrieves all items in a specific playlist
// @Tags playlists
// @Accept json
// @Produce json
// @Param id path int true "Playlist ID"
// @Param userId query int true "User ID"
// @Param clientId query int true "Client ID"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[any]] "Playlist items retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Playlist not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /playlists/{id}/items [get]
func (h *PlaylistHandler) GetPlaylistItems(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	playlistID := c.Param("id")
	if playlistID == "" {
		log.Warn().Msg("Playlist ID is required")
		responses.RespondBadRequest(c, nil, "Playlist ID is required")
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	clientID, err := strconv.ParseUint(c.Query("clientId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("clientId", c.Query("clientId")).Msg("Invalid client ID")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Msg("Getting playlist items")

	items, err := h.service.GetPlaylistItems(ctx, userID, clientID, playlistID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("playlistID", playlistID).
			Msg("Failed to retrieve playlist items")
		responses.RespondInternalError(c, err, "Failed to retrieve playlist items")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Int("itemCount", len(items)).
		Msg("Playlist items retrieved successfully")

	responses.RespondOK(c, items, "Playlist items retrieved successfully")
}

// CreatePlaylist godoc
// @Summary Create a new playlist
// @Description Creates a new playlist in the specified client
// @Tags playlists
// @Accept json
// @Produce json
// @Param userId query int true "User ID"
// @Param clientId query int true "Client ID"
// @Param playlist body object true "Playlist data including name and description"
// @Success 201 {object} responses.APIResponse[models.MediaItem[mediatypes.Playlist]] "Playlist created successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /playlists [post]
func (h *PlaylistHandler) CreatePlaylist(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	clientID, err := strconv.ParseUint(c.Query("clientId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("clientId", c.Query("clientId")).Msg("Invalid client ID")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid request body for CreatePlaylist")
		responses.RespondValidationError(c, err)
		return
	}

	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("name", req.Name).
		Msg("Creating playlist")

	playlist, err := h.service.CreatePlaylist(ctx, userID, clientID, req.Name, req.Description)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("name", req.Name).
			Msg("Failed to create playlist")
		responses.RespondInternalError(c, err, "Failed to create playlist")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("name", req.Name).
		Msg("Playlist created successfully")

	responses.RespondCreated(c, playlist, "Playlist created successfully")
}

// UpdatePlaylist godoc
// @Summary Update a playlist
// @Description Updates an existing playlist
// @Tags playlists
// @Accept json
// @Produce json
// @Param id path int true "Playlist ID"
// @Param userId query int true "User ID"
// @Param clientId query int true "Client ID"
// @Param playlist body object true "Updated playlist data"
// @Success 200 {object} responses.APIResponse[models.MediaItem[mediatypes.Playlist]] "Playlist updated successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Playlist not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /playlists/{id} [put]
func (h *PlaylistHandler) UpdatePlaylist(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	playlistID := c.Param("id")
	if playlistID == "" {
		log.Warn().Msg("Playlist ID is required")
		responses.RespondBadRequest(c, nil, "Playlist ID is required")
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	clientID, err := strconv.ParseUint(c.Query("clientId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("clientId", c.Query("clientId")).Msg("Invalid client ID")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid request body for UpdatePlaylist")
		responses.RespondValidationError(c, err)
		return
	}

	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Str("name", req.Name).
		Msg("Updating playlist")

	playlist, err := h.service.UpdatePlaylist(ctx, userID, clientID, playlistID, req.Name, req.Description)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("playlistID", playlistID).
			Msg("Failed to update playlist")
		responses.RespondInternalError(c, err, "Failed to update playlist")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Msg("Playlist updated successfully")

	responses.RespondOK(c, playlist, "Playlist updated successfully")
}

// DeletePlaylist godoc
// @Summary Delete a playlist
// @Description Deletes an existing playlist
// @Tags playlists
// @Accept json
// @Produce json
// @Param id path int true "Playlist ID"
// @Param userId query int true "User ID"
// @Param clientId query int true "Client ID"
// @Success 200 {object} responses.APIResponse[any] "Playlist deleted successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /playlists/{id} [delete]
func (h *PlaylistHandler) DeletePlaylist(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	playlistID := c.Param("id")
	if playlistID == "" {
		log.Warn().Msg("Playlist ID is required")
		responses.RespondBadRequest(c, nil, "Playlist ID is required")
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	clientID, err := strconv.ParseUint(c.Query("clientId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("clientId", c.Query("clientId")).Msg("Invalid client ID")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Msg("Deleting playlist")

	err = h.service.DeletePlaylist(ctx, userID, clientID, playlistID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("playlistID", playlistID).
			Msg("Failed to delete playlist")
		responses.RespondInternalError(c, err, "Failed to delete playlist")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Msg("Playlist deleted successfully")

	responses.RespondOK(c, responses.EmptyResponse{Success: true}, "Playlist deleted successfully")
}

// AddItemToPlaylist godoc
// @Summary Add an item to a playlist
// @Description Adds a media item to an existing playlist
// @Tags playlists
// @Accept json
// @Produce json
// @Param id path int true "Playlist ID"
// @Param userId query int true "User ID"
// @Param clientId query int true "Client ID"
// @Param itemData body object true "Item data including media item ID"
// @Success 200 {object} responses.APIResponse[any] "Item added to playlist successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /playlists/{id}/items [post]
func (h *PlaylistHandler) AddItemToPlaylist(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	playlistID := c.Param("id")
	if playlistID == "" {
		log.Warn().Msg("Playlist ID is required")
		responses.RespondBadRequest(c, nil, "Playlist ID is required")
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	clientID, err := strconv.ParseUint(c.Query("clientId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("clientId", c.Query("clientId")).Msg("Invalid client ID")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	var req struct {
		ItemID string `json:"itemId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid request body for AddItemToPlaylist")
		responses.RespondValidationError(c, err)
		return
	}

	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Str("itemID", req.ItemID).
		Msg("Adding item to playlist")

	err = h.service.AddItemToPlaylist(ctx, userID, clientID, playlistID, req.ItemID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("playlistID", playlistID).
			Str("itemID", req.ItemID).
			Msg("Failed to add item to playlist")
		responses.RespondInternalError(c, err, "Failed to add item to playlist")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Str("itemID", req.ItemID).
		Msg("Item added to playlist successfully")

	responses.RespondOK(c, responses.EmptyResponse{Success: true}, "Item added to playlist successfully")
}

// RemoveItemFromPlaylist godoc
// @Summary Remove an item from a playlist
// @Description Removes a media item from an existing playlist
// @Tags playlists
// @Accept json
// @Produce json
// @Param id path int true "Playlist ID"
// @Param itemId path int true "Item ID to remove"
// @Param userId query int true "User ID"
// @Param clientId query int true "Client ID"
// @Success 200 {object} responses.APIResponse[any] "Item removed from playlist successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /playlists/{id}/items/{itemId} [delete]
func (h *PlaylistHandler) RemoveItemFromPlaylist(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	playlistID := c.Param("id")
	if playlistID == "" {
		log.Warn().Msg("Playlist ID is required")
		responses.RespondBadRequest(c, nil, "Playlist ID is required")
		return
	}

	itemID := c.Param("itemId")
	if itemID == "" {
		log.Warn().Msg("Item ID is required")
		responses.RespondBadRequest(c, nil, "Item ID is required")
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	clientID, err := strconv.ParseUint(c.Query("clientId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("clientId", c.Query("clientId")).Msg("Invalid client ID")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Str("itemID", itemID).
		Msg("Removing item from playlist")

	err = h.service.RemoveItemFromPlaylist(ctx, userID, clientID, playlistID, itemID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("playlistID", playlistID).
			Str("itemID", itemID).
			Msg("Failed to remove item from playlist")
		responses.RespondInternalError(c, err, "Failed to remove item from playlist")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Str("itemID", itemID).
		Msg("Item removed from playlist successfully")

	responses.RespondOK(c, responses.EmptyResponse{Success: true}, "Item removed from playlist successfully")
}

// ReorderPlaylistItems godoc
// @Summary Reorder items in a playlist
// @Description Changes the order of items in a playlist
// @Tags playlists
// @Accept json
// @Produce json
// @Param id path int true "Playlist ID"
// @Param userId query int true "User ID"
// @Param clientId query int true "Client ID"
// @Param orderData body object true "New item order data"
// @Success 200 {object} responses.APIResponse[any] "Playlist items reordered successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /playlists/{id}/reorder [post]
func (h *PlaylistHandler) ReorderPlaylistItems(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	playlistID := c.Param("id")
	if playlistID == "" {
		log.Warn().Msg("Playlist ID is required")
		responses.RespondBadRequest(c, nil, "Playlist ID is required")
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	clientID, err := strconv.ParseUint(c.Query("clientId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("clientId", c.Query("clientId")).Msg("Invalid client ID")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	var req struct {
		ItemIDs []string `json:"itemIds" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid request body for ReorderPlaylistItems")
		responses.RespondValidationError(c, err)
		return
	}

	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Int("itemCount", len(req.ItemIDs)).
		Msg("Reordering playlist items")

	err = h.service.ReorderPlaylistItems(ctx, userID, clientID, playlistID, req.ItemIDs)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("playlistID", playlistID).
			Msg("Failed to reorder playlist items")
		responses.RespondInternalError(c, err, "Failed to reorder playlist items")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Int("itemCount", len(req.ItemIDs)).
		Msg("Playlist items reordered successfully")

	responses.RespondOK(c, responses.EmptyResponse{Success: true}, "Playlist items reordered successfully")
}

// SearchPlaylists godoc
// @Summary Search for playlists
// @Description Searches for playlists matching a query string
// @Tags playlists
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param userId query int true "User ID"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Playlist]] "Playlists found"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /playlists/search [get]
func (h *PlaylistHandler) SearchPlaylists(c *gin.Context) {
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

	log.Debug().
		Uint64("userID", userID).
		Str("query", query).
		Msg("Searching playlists")

	playlists, err := h.service.SearchPlaylists(ctx, userID, query)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Str("query", query).
			Msg("Failed to search playlists")
		responses.RespondInternalError(c, err, "Failed to search playlists")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Str("query", query).
		Int("count", len(playlists)).
		Msg("Playlists search completed successfully")

	responses.RespondOK(c, playlists, "Playlists found")
}

// SyncPlaylist godoc
// @Summary Sync a playlist across clients
// @Description Synchronizes a playlist's content across all compatible clients
// @Tags playlists
// @Accept json
// @Produce json
// @Param id path int true "Playlist ID"
// @Param userId query int true "User ID"
// @Param clientId query int true "Source client ID"
// @Success 200 {object} responses.APIResponse[any] "Playlist sync initiated"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /playlists/{id}/sync [post]
func (h *PlaylistHandler) SyncPlaylist(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	playlistID := c.Param("id")
	if playlistID == "" {
		log.Warn().Msg("Playlist ID is required")
		responses.RespondBadRequest(c, nil, "Playlist ID is required")
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	clientID, err := strconv.ParseUint(c.Query("clientId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("clientId", c.Query("clientId")).Msg("Invalid client ID")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Msg("Syncing playlist across clients")

	err = h.service.SyncPlaylist(ctx, userID, clientID, playlistID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("playlistID", playlistID).
			Msg("Failed to sync playlist")
		responses.RespondInternalError(c, err, "Failed to sync playlist")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Msg("Playlist sync initiated successfully")

	responses.RespondOK(c, responses.EmptyResponse{Success: true}, "Playlist sync initiated")
}
