// handlers/user_playlists.go
package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"

	mediatypes "suasor/client/media/types"
	"suasor/services"
	"suasor/types/models"
	"suasor/types/requests"
	"suasor/types/responses"
	"suasor/utils"
)

// UserPlaylistHandler handles user-specific operations for playlists
type UserPlaylistHandler struct {
	mediaItemService    services.UserMediaItemService[*mediatypes.Playlist]
	userPlaylistService services.UserPlaylistService
	playlistService     services.PlaylistService
}

// NewUserPlaylistHandler creates a new user playlist handler
func NewUserPlaylistHandler(
	userPlaylistService services.UserMediaItemService[*mediatypes.Playlist],
	playlistService services.PlaylistService,
) *UserPlaylistHandler {
	return &UserPlaylistHandler{
		userPlaylistService: userPlaylistService,
		playlistService:     playlistService,
	}
}

// GetUserPlaylists godoc
// @Summary Get user's playlists
// @Description Retrieves all playlists owned by the authenticated user
// @Tags playlists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[*mediatypes.Playlist]] "Playlists retrieved successfully"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /user/playlists [get]
func (h *UserPlaylistHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access playlists without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	log.Debug().
		Uint64("userID", uid).
		Msg("Getting user playlists")

	// Get user playlists
	playlists, err := h.userPlaylistService.GetByUserID(ctx, uid)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Msg("Failed to retrieve user playlists")
		responses.RespondInternalError(c, err, "Failed to retrieve playlists")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", len(playlists)).
		Msg("User playlists retrieved successfully")
	responses.RespondOK(c, playlists, "Playlists retrieved successfully")
}

// CreatePlaylist godoc
// @Summary Create a new playlist
// @Description Creates a new playlist for the authenticated user
// @Tags playlists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param playlist body requests.PlaylistCreateRequest true "Playlist details"
// @Success 201 {object} responses.APIResponse[models.MediaItem[*mediatypes.Playlist]] "Playlist created successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /user/playlists [post]
func (h *UserPlaylistHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to create playlist without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	// Parse request body
	var req requests.PlaylistCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid request body")
		responses.RespondBadRequest(c, err, "Invalid request body")
		return
	}

	log.Debug().
		Uint64("userID", uid).
		Str("name", req.Name).
		Msg("Creating new playlist")

	// Create playlist data
	playlist := &mediatypes.Playlist{
		ItemList: mediatypes.ItemList{
			ItemCount: 0,
			OwnerID:   uid,
			Details: mediatypes.MediaDetails{
				Title:       req.Name,
				Description: req.Description,
			},
			IsPublic: req.IsPublic,
			IsSmart:  req.IsSmart,
			SmartCriteria: map[string]any{
				"genre":    req.Genre,
				"year":     req.Year,
				"rating":   req.Rating,
				"duration": req.Duration,
			},
			// AutoUpdateTime: ,
		},
	}

	playlist.ItemList.AddItem(mediatypes.ListItem{
		ItemID:        0,
		Position:      0,
		LastChanged:   time.Now(),
		ChangeHistory: []mediatypes.ChangeRecord{},
	})

	// Create media item
	mediaItem := models.MediaItem[*mediatypes.Playlist]{
		Type:      mediatypes.MediaTypePlaylist,
		Title:     playlist.ItemList.Details.Title,
		Data:      playlist,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create playlist
	createdPlaylist, err := h.userPlaylistService.Create(ctx, mediaItem)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Msg("Failed to create playlist")
		responses.RespondInternalError(c, err, "Failed to create playlist")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("playlistID", createdPlaylist.ID).
		Msg("Playlist created successfully")
	responses.RespondCreated(c, createdPlaylist, "Playlist created successfully")
}

// UpdatePlaylist godoc
// @Summary Update a playlist
// @Description Updates an existing playlist owned by the authenticated user
// @Tags playlists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Playlist ID"
// @Param playlist body requests.PlaylistUpdateRequest true "Updated playlist details"
// @Success 200 {object} responses.APIResponse[models.MediaItem[*mediatypes.Playlist]] "Playlist updated successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 403 {object} responses.ErrorResponse[any] "Forbidden"
// @Failure 404 {object} responses.ErrorResponse[any] "Playlist not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /user/playlists/{id} [put]
func (h *UserPlaylistHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to update playlist without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	// Parse playlist ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid playlist ID")
		responses.RespondBadRequest(c, err, "Invalid playlist ID")
		return
	}

	// Parse request body
	var req requests.PlaylistUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid request body")
		responses.RespondBadRequest(c, err, "Invalid request body")
		return
	}

	log.Debug().
		Uint64("userID", uid).
		Uint64("playlistID", id).
		Msg("Updating playlist")

	// Get existing playlist
	existingPlaylist, err := h.userPlaylistService.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("playlistID", id).
			Msg("Failed to retrieve playlist")
		responses.RespondNotFound(c, err, "Playlist not found")
		return
	}

	// Check if user owns the playlist
	if existingPlaylist.Data.OwnerID != uid {
		log.Warn().
			Uint64("userID", uid).
			Uint64("ownerID", existingPlaylist.Data.OwnerID).
			Uint64("playlistID", id).
			Msg("User does not own the playlist")
		responses.RespondForbidden(c, nil, "You do not have permission to update this playlist")
		return
	}

	// Update playlist
	existingPlaylist.Title = req.Name
	existingPlaylist.Data.ItemList.Details.Title = req.Name
	existingPlaylist.Data.ItemList.Details.Description = req.Description
	existingPlaylist.Data.IsPublic = req.IsPublic

	// Save updated playlist
	updatedPlaylist, err := h.userPlaylistService.Update(ctx, *existingPlaylist)
	if err != nil {
		log.Error().Err(err).
			Uint64("playlistID", id).
			Msg("Failed to update playlist")
		responses.RespondInternalError(c, err, "Failed to update playlist")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("playlistID", id).
		Msg("Playlist updated successfully")
	responses.RespondOK(c, updatedPlaylist, "Playlist updated successfully")
}

// DeletePlaylist godoc
// @Summary Delete a playlist
// @Description Deletes a playlist owned by the authenticated user
// @Tags playlists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Playlist ID"
// @Success 200 {object} responses.APIResponse[any] "Playlist deleted successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 403 {object} responses.ErrorResponse[any] "Forbidden"
// @Failure 404 {object} responses.ErrorResponse[any] "Playlist not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /user/playlists/{id} [delete]
func (h *UserPlaylistHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to delete playlist without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	// Parse playlist ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid playlist ID")
		responses.RespondBadRequest(c, err, "Invalid playlist ID")
		return
	}

	log.Debug().
		Uint64("userID", uid).
		Uint64("playlistID", id).
		Msg("Deleting playlist")

	// Get existing playlist
	existingPlaylist, err := h.userPlaylistService.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("playlistID", id).
			Msg("Failed to retrieve playlist")
		responses.RespondNotFound(c, err, "Playlist not found")
		return
	}

	// Check if user owns the playlist
	if existingPlaylist.Data.OwnerID != uid {
		log.Warn().
			Uint64("userID", uid).
			Uint64("ownerID", existingPlaylist.Data.OwnerID).
			Uint64("playlistID", id).
			Msg("User does not own the playlist")
		responses.RespondForbidden(c, nil, "You do not have permission to delete this playlist")
		return
	}

	// Delete playlist
	err = h.userPlaylistService.Delete(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("playlistID", id).
			Msg("Failed to delete playlist")
		responses.RespondInternalError(c, err, "Failed to delete playlist")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("playlistID", id).
		Msg("Playlist deleted successfully")
	responses.RespondOK(c, http.StatusOK, "Playlist deleted successfully")
}

// AddTrackToPlaylist godoc
// @Summary Add a track to a playlist
// @Description Adds a track to a playlist owned by the authenticated user
// @Tags playlists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Playlist ID"
// @Param track body requests.PlaylistAddTrackRequest true "Track details"
// @Success 200 {object} responses.APIResponse[models.MediaItem[*mediatypes.Playlist]] "Track added successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 403 {object} responses.ErrorResponse[any] "Forbidden"
// @Failure 404 {object} responses.ErrorResponse[any] "Playlist not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /user/playlists/{id}/tracks [post]
func (h *UserPlaylistHandler) AddItem(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to add track to playlist without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	// Parse playlist ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid playlist ID")
		responses.RespondBadRequest(c, err, "Invalid playlist ID")
		return
	}

	// Parse request body
	var req requests.PlaylistAddTrackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid request body")
		responses.RespondBadRequest(c, err, "Invalid request body")
		return
	}

	log.Debug().
		Uint64("userID", uid).
		Uint64("playlistID", id).
		Uint64("trackID", req.TrackID).
		Msg("Adding track to playlist")

	// Get existing playlist
	existingPlaylist, err := h.userPlaylistService.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("playlistID", id).
			Msg("Failed to retrieve playlist")
		responses.RespondNotFound(c, err, "Playlist not found")
		return
	}

	// Check if user owns the playlist
	if existingPlaylist.Data.OwnerID != uid {
		log.Warn().
			Uint64("userID", uid).
			Uint64("ownerID", existingPlaylist.Data.OwnerID).
			Uint64("playlistID", id).
			Msg("User does not own the playlist")
		responses.RespondForbidden(c, nil, "You do not have permission to modify this playlist")
		return
	}

	// Add track to playlist
	err = h.playlistService.AddItemToPlaylist(ctx, id, req.TrackID)
	if err != nil {
		log.Error().Err(err).
			Uint64("playlistID", id).
			Uint64("trackID", req.TrackID).
			Msg("Failed to add track to playlist")
		responses.RespondInternalError(c, err, "Failed to add track to playlist")
		return
	}

	// Get updated playlist
	updatedPlaylist, err := h.userPlaylistService.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("playlistID", id).
			Msg("Failed to retrieve updated playlist")
		responses.RespondInternalError(c, err, "Failed to retrieve updated playlist")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("playlistID", id).
		Uint64("trackID", req.TrackID).
		Msg("Track added to playlist successfully")
	responses.RespondOK(c, updatedPlaylist, "Track added to playlist successfully")
}

// RemoveTrackFromPlaylist godoc
// @Summary Remove a track from a playlist
// @Description Removes a track from a playlist owned by the authenticated user
// @Tags playlists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Playlist ID"
// @Param trackId path int true "Track ID"
// @Success 200 {object} responses.APIResponse[models.MediaItem[*mediatypes.Playlist]] "Track removed successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 403 {object} responses.ErrorResponse[any] "Forbidden"
// @Failure 404 {object} responses.ErrorResponse[any] "Playlist not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /user/playlists/{id}/tracks/{trackId} [delete]
func (h *UserPlaylistHandler) RemoveItem(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to remove track from playlist without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	// Parse playlist ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid playlist ID")
		responses.RespondBadRequest(c, err, "Invalid playlist ID")
		return
	}

	// Parse track ID
	trackID, err := strconv.ParseUint(c.Param("trackId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("trackId", c.Param("trackId")).Msg("Invalid track ID")
		responses.RespondBadRequest(c, err, "Invalid track ID")
		return
	}

	log.Debug().
		Uint64("userID", uid).
		Uint64("playlistID", id).
		Uint64("trackID", trackID).
		Msg("Removing track from playlist")

	// Get existing playlist
	existingPlaylist, err := h.userPlaylistService.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("playlistID", id).
			Msg("Failed to retrieve playlist")
		responses.RespondNotFound(c, err, "Playlist not found")
		return
	}

	// Check if user owns the playlist
	if existingPlaylist.Data.OwnerID != uid {
		log.Warn().
			Uint64("userID", uid).
			Uint64("ownerID", existingPlaylist.Data.OwnerID).
			Uint64("playlistID", id).
			Msg("User does not own the playlist")
		responses.RespondForbidden(c, nil, "You do not have permission to modify this playlist")
		return
	}

	// Remove track from playlist
	err = h.playlistService.RemoveItemFromPlaylist(ctx, id, trackID)
	if err != nil {
		log.Error().Err(err).
			Uint64("playlistID", id).
			Uint64("trackID", trackID).
			Msg("Failed to remove track from playlist")
		responses.RespondInternalError(c, err, "Failed to remove track from playlist")
		return
	}

	// Get updated playlist
	updatedPlaylist, err := h.userPlaylistService.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("playlistID", id).
			Msg("Failed to retrieve updated playlist")
		responses.RespondInternalError(c, err, "Failed to retrieve updated playlist")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("playlistID", id).
		Uint64("trackID", trackID).
		Msg("Track removed from playlist successfully")
	responses.RespondOK(c, updatedPlaylist, "Track removed from playlist successfully")
}

// Reorder godoc
// @Summary Reorder playlist items
// @Description Reorders the items in a playlist
// @Tags playlists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Playlist ID"
// @Param request body requests.PlaylistReorderRequest true "Reorder request"
// @Success 200 {object} responses.APIResponse[models.MediaItem[*types.Playlist]] "Playlist reordered successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 403 {object} responses.ErrorResponse[any] "Forbidden"
// @Failure 404 {object} responses.ErrorResponse[any] "Playlist not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /user/playlists/{id}/reorder [post]
func (h *UserPlaylistHandler) ReorderItems(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to reorder playlist without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	// Parse playlist ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid playlist ID")
		responses.RespondBadRequest(c, err, "Invalid playlist ID")
		return
	}

	// Parse request body
	var req requests.PlaylistReorderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid request body")
		responses.RespondBadRequest(c, err, "Invalid request body")
		return
	}

	log.Debug().
		Uint64("userID", uid).
		Uint64("playlistID", id).
		Msg("Reordering playlist items")

	// Get existing playlist
	existingPlaylist, err := h.userPlaylistService.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("playlistID", id).
			Msg("Failed to retrieve playlist")
		responses.RespondNotFound(c, err, "Playlist not found")
		return
	}

	// Check if user owns the playlist
	if existingPlaylist.Data.OwnerID != uid {
		log.Warn().
			Uint64("userID", uid).
			Uint64("ownerID", existingPlaylist.Data.OwnerID).
			Uint64("playlistID", id).
			Msg("User does not own the playlist")
		responses.RespondForbidden(c, nil, "You do not have permission to reorder this playlist")
		return
	}

	// Reorder playlist items
	err = h.userPlaylistService.ReorderItems(ctx, id, req.ItemIDs)
	if err != nil {
		log.Error().Err(err).
			Uint64("playlistID", id).
			Msg("Failed to reorder playlist items")
		responses.RespondInternalError(c, err, "Failed to reorder playlist items")
		return
	}

	// Get updated playlist
	updatedPlaylist, err := h.userPlaylistService.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("playlistID", id).
			Msg("Failed to retrieve updated playlist")
		responses.RespondInternalError(c, err, "Failed to retrieve updated playlist")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("playlistID", id).
		Msg("Playlist reordered successfully")
	responses.RespondOK(c, updatedPlaylist, "Playlist reordered successfully")
}
