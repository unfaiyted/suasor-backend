// handlers/client_media_playlist.go
package handlers

import (
	"strconv"
	mediatypes "suasor/client/media/types"
	clienttypes "suasor/client/types"
	"suasor/services"
	models "suasor/types/models"
	"suasor/types/responses"
	"suasor/utils"

	"github.com/gin-gonic/gin"
)

func createPlaylistMediaItem[T mediatypes.Playlist](clientID uint64, clientType clienttypes.ClientMediaType, externalID string, data mediatypes.Playlist) models.MediaItem[mediatypes.Playlist] {
	mediaItem := models.MediaItem[mediatypes.Playlist]{
		Type:        mediatypes.MediaTypePlaylist,
		SyncClients: []models.SyncClient{},
		ExternalIDs: []models.ExternalID{},
		Data:        data,
	}

	// Set client info
	mediaItem.SetClientInfo(clientID, clientType, externalID)

	// Only add external ID if provided
	if externalID != "" {
		mediaItem.AddExternalID("client", externalID)
	}

	return mediaItem
}

// ClientPlaylistHandler handles playlist-related operations for media clients
type ClientPlaylistHandler[T clienttypes.ClientMediaConfig] struct {
	playlistService services.ClientPlaylistService[T]
}

// NewClientPlaylistHandler creates a new media client playlist handler
func NewClientPlaylistHandler[T clienttypes.ClientMediaConfig](playlistService services.ClientPlaylistService[T]) *ClientPlaylistHandler[T] {
	return &ClientPlaylistHandler[T]{
		playlistService: playlistService,
	}
}

// GetPlaylistByID godoc
// @Summary Get playlist by ID
// @Description Retrieves a specific playlist from the client by ID
// @Tags playlists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param playlistID path string true "Playlist ID"
// @Success 200 {object} responses.APIResponse[models.MediaItem[mediatypes.Playlist]] "Playlist retrieved"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid client ID"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /clients/media/{clientID}/playlists/{playlistID} [get]
func (h *ClientPlaylistHandler[T]) GetPlaylistByID(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	log.Info().Msg("Getting playlist by ID")

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

	playlist, err := h.playlistService.GetPlaylistByID(ctx, uid, clientID, playlistID)
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
		Msg("Playlist retrieved successfully")
	responses.RespondOK(c, playlist, "Playlist retrieved successfully")
}

// GetPlaylists godoc
// @Summary Get all playlists
// @Description Retrieves all playlists from the client
// @Tags playlists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param count query int false "Maximum number of playlists to return"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Playlist]] "Playlists retrieved"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid client ID"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /clients/media/{clientID}/playlists [get]
func (h *ClientPlaylistHandler[T]) GetPlaylists(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	log.Info().Msg("Getting all playlists")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access playlists without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	// Get count parameter
	count := 0
	countParam := c.Query("count")
	if countParam != "" {
		var err error
		count, err = strconv.Atoi(countParam)
		if err != nil {
			log.Error().Err(err).Str("count", countParam).Msg("Invalid count format")
			responses.RespondBadRequest(c, err, "Invalid count")
			return
		}
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", count).
		Msg("Retrieving playlists")

	playlists, err := h.playlistService.GetPlaylists(ctx, uid, count)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Int("count", count).
			Msg("Failed to retrieve playlists")
		responses.RespondInternalError(c, err, "Failed to retrieve playlists")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", count).
		Int("playlistsReturned", len(playlists)).
		Msg("Playlists retrieved successfully")
	responses.RespondOK(c, playlists, "Playlists retrieved successfully")
}

// CreatePlaylist godoc
// @Summary Create a new playlist
// @Description Creates a new playlist on the specified client
// @Tags playlists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param playlist body object true "Playlist creation data"
// @Success 201 {object} responses.APIResponse[models.MediaItem[mediatypes.Playlist]] "Playlist created"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /clients/media/{clientID}/playlists [post]
func (h *ClientPlaylistHandler[T]) CreatePlaylist(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	log.Info().Msg("Creating playlist")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to create playlist without authentication")
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

	playlist, err := h.playlistService.CreatePlaylist(ctx, uid, clientID, req.Name, req.Description)
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
		Msg("Playlist created successfully")
	responses.RespondCreated(c, playlist, "Playlist created successfully")
}

// UpdatePlaylist godoc
// @Summary Update a playlist
// @Description Updates an existing playlist on the specified client
// @Tags playlists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param playlistID path string true "Playlist ID"
// @Param playlist body object true "Playlist update data"
// @Success 200 {object} responses.APIResponse[models.MediaItem[mediatypes.Playlist]] "Playlist updated"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /clients/media/{clientID}/playlists/{playlistID} [put]
func (h *ClientPlaylistHandler[T]) UpdatePlaylist(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	log.Info().Msg("Updating playlist")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to update playlist without authentication")
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
		Str("playlistID", playlistID).
		Str("name", req.Name).
		Msg("Updating playlist")

	playlist, err := h.playlistService.UpdatePlaylist(ctx, uid, clientID, playlistID, req.Name, req.Description)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Str("playlistID", playlistID).
			Msg("Failed to update playlist")
		responses.RespondInternalError(c, err, "Failed to update playlist")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Msg("Playlist updated successfully")
	responses.RespondOK(c, playlist, "Playlist updated successfully")
}

// DeletePlaylist godoc
// @Summary Delete a playlist
// @Description Deletes a playlist from the specified client
// @Tags playlists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param playlistID path string true "Playlist ID"
// @Success 200 {object} responses.APIResponse[string] "Playlist deleted"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /clients/media/{clientID}/playlists/{playlistID} [delete]
func (h *ClientPlaylistHandler[T]) DeletePlaylist(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	log.Info().Msg("Deleting playlist")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to delete playlist without authentication")
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
		Msg("Deleting playlist")

	err = h.playlistService.DeletePlaylist(ctx, uid, clientID, playlistID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Str("playlistID", playlistID).
			Msg("Failed to delete playlist")
		responses.RespondInternalError(c, err, "Failed to delete playlist")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Msg("Playlist deleted successfully")
	responses.RespondOK(c, "Playlist deleted successfully", "Playlist deleted successfully")
}

// AddItemToPlaylist godoc
// @Summary Add an item to a playlist
// @Description Adds a media item to an existing playlist
// @Tags playlists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param playlistID path string true "Playlist ID"
// @Param item body object true "Item to add"
// @Success 200 {object} responses.APIResponse[string] "Item added to playlist"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /clients/media/{clientID}/playlists/{playlistID}/items [post]
func (h *ClientPlaylistHandler[T]) AddItemToPlaylist(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
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

	err = h.playlistService.AddItemToPlaylist(ctx, uid, clientID, playlistID, req.ItemID)
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

// RemoveItemFromPlaylist godoc
// @Summary Remove an item from a playlist
// @Description Removes a media item from an existing playlist
// @Tags playlists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param playlistID path string true "Playlist ID"
// @Param itemID path string true "Item ID to remove"
// @Success 200 {object} responses.APIResponse[string] "Item removed from playlist"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /clients/media/{clientID}/playlists/{playlistID}/items/{itemID} [delete]
func (h *ClientPlaylistHandler[T]) RemoveItemFromPlaylist(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
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

	err = h.playlistService.RemoveItemFromPlaylist(ctx, uid, clientID, playlistID, itemID)
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

// SearchPlaylists godoc
// @Summary Search playlists
// @Description Searches for playlists matching the given query
// @Tags playlists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param q query string true "Search query"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Playlist]] "Playlists found"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /playlists/search [get]
func (h *ClientPlaylistHandler[T]) SearchPlaylists(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
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

	playlists, err := h.playlistService.SearchPlaylists(ctx, uid, query)
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
		Msg("Playlist search completed successfully")
	responses.RespondOK(c, playlists, "Playlists retrieved successfully")
}
