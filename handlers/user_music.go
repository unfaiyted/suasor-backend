// handlers/user_music.go
package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	mediatypes "suasor/clients/media/types"
	"suasor/services"
	"suasor/types/models"
	"suasor/types/responses"
	"suasor/utils/logger"
)

// UserMusicHandler handles operations for music items related to users
type UserMusicHandler struct {
	trackService  services.UserMediaItemService[*mediatypes.Track]
	albumService  services.UserMediaItemService[*mediatypes.Album]
	artistService services.UserMediaItemService[*mediatypes.Artist]
}

// NewUserMusicHandler creates a new user music handler
func NewUserMusicHandler(
	trackService services.UserMediaItemService[*mediatypes.Track],
	albumService services.UserMediaItemService[*mediatypes.Album],
	artistService services.UserMediaItemService[*mediatypes.Artist],
) *UserMusicHandler {
	return &UserMusicHandler{
		trackService:  trackService,
		albumService:  albumService,
		artistService: artistService,
	}
}

// GetFavoriteTracks godoc
// @Summary Get user favorite tracks
// @Description Retrieves tracks that a user has marked as favorites
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Maximum number of tracks to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Track]] "Tracks retrieved successfully"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /user/music/tracks/favorites [get]
func (h *UserMusicHandler) GetFavoriteTracks(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access favorites without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	log.Debug().
		Uint64("userID", uid).
		Int("limit", limit).
		Msg("Getting favorite tracks")

	// This is a placeholder for a real implementation
	// In a real implementation, you would query for tracks specifically marked as favorites
	options := mediatypes.QueryOptions{
		MediaType: mediatypes.MediaTypeTrack,
		OwnerID:   uid,
		Favorites: true,
		Limit:     limit,
	}

	tracks, err := h.trackService.SearchUserContent(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Msg("Failed to retrieve favorite tracks")
		responses.RespondInternalError(c, err, "Failed to retrieve favorite tracks")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", len(tracks)).
		Msg("Favorite tracks retrieved successfully")
	responses.RespondOK(c, tracks, "Favorite tracks retrieved successfully")
}

// GetFavoriteAlbums godoc
// @Summary Get user favorite albums
// @Description Retrieves albums that a user has marked as favorites
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Maximum number of albums to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Album]] "Albums retrieved successfully"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /user/music/albums/favorites [get]
func (h *UserMusicHandler) GetFavoriteAlbums(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access favorites without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	log.Debug().
		Uint64("userID", uid).
		Int("limit", limit).
		Msg("Getting favorite albums")

	// This is a placeholder for a real implementation
	// In a real implementation, you would query for albums specifically marked as favorites
	options := mediatypes.QueryOptions{
		MediaType: mediatypes.MediaTypeAlbum,
		OwnerID:   uid,
		Favorites: true,
		Limit:     limit,
	}

	albums, err := h.albumService.SearchUserContent(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Msg("Failed to retrieve favorite albums")
		responses.RespondInternalError(c, err, "Failed to retrieve favorite albums")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", len(albums)).
		Msg("Favorite albums retrieved successfully")
	responses.RespondOK(c, albums, "Favorite albums retrieved successfully")
}

// GetFavoriteArtists godoc
// @Summary Get user favorite artists
// @Description Retrieves artists that a user has marked as favorites
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Maximum number of artists to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Artist]] "Artists retrieved successfully"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /user/music/artists/favorites [get]
func (h *UserMusicHandler) GetFavoriteArtists(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access favorites without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	log.Debug().
		Uint64("userID", uid).
		Int("limit", limit).
		Msg("Getting favorite artists")

	// This is a placeholder for a real implementation
	// In a real implementation, you would query for artists specifically marked as favorites
	options := mediatypes.QueryOptions{
		MediaType: mediatypes.MediaTypeArtist,
		OwnerID:   uid,
		Favorites: true,
		Limit:     limit,
	}

	artists, err := h.artistService.SearchUserContent(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Msg("Failed to retrieve favorite artists")
		responses.RespondInternalError(c, err, "Failed to retrieve favorite artists")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", len(artists)).
		Msg("Favorite artists retrieved successfully")
	responses.RespondOK(c, artists, "Favorite artists retrieved successfully")
}

// GetRecentlyPlayedTracks godoc
// @Summary Get recently played tracks
// @Description Retrieves tracks that a user has recently played
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Maximum number of tracks to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Track]] "Tracks retrieved successfully"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /user/music/tracks/recently-played [get]
func (h *UserMusicHandler) GetRecentlyPlayedTracks(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access recently played without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	log.Debug().
		Uint64("userID", uid).
		Int("limit", limit).
		Msg("Getting recently played tracks")

	// This is a placeholder for a real implementation
	// In a real implementation, you would query play history to find recently played tracks
	// For now, we'll use a simplified approach
	options := mediatypes.QueryOptions{
		MediaType: mediatypes.MediaTypeTrack,
		OwnerID:   uid,
		Sort:      "last_played",
		SortOrder: "desc",
		Limit:     limit,
	}

	tracks, err := h.trackService.SearchUserContent(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Msg("Failed to retrieve recently played tracks")
		responses.RespondInternalError(c, err, "Failed to retrieve recently played tracks")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", len(tracks)).
		Msg("Recently played tracks retrieved successfully")
	responses.RespondOK(c, tracks, "Recently played tracks retrieved successfully")
}

// UpdateTrackUserData godoc
// @Summary Update user data for a track
// @Description Updates user-specific data for a track (favorite, rating, etc.)
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Track ID"
// @Param data body requests.UserMediaItemDataRequest true "Updated user data"
// @Success 200 {object} responses.APIResponse[models.MediaItem[mediatypes.Track]] "Track updated successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 404 {object} responses.ErrorResponse[any] "Track not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /user/music/tracks/{id} [patch]
func (h *UserMusicHandler) UpdateTrackUserData(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to update track data without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	trackID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid track ID")
		responses.RespondBadRequest(c, err, "Invalid track ID")
		return
	}

	// Parse request body
	var userData models.UserMediaItemData[*mediatypes.Track]
	if err := c.ShouldBindJSON(&userData); err != nil {
		log.Warn().Err(err).Msg("Invalid request body")
		responses.RespondBadRequest(c, err, "Invalid request body")
		return
	}

	log.Debug().
		Uint64("userID", uid).
		Uint64("trackID", trackID).
		Interface("userData", userData).
		Msg("Updating track user data")

	// Get the existing track first
	track, err := h.trackService.GetByID(ctx, trackID)
	if err != nil {
		log.Error().Err(err).
			Uint64("trackID", trackID).
			Msg("Failed to retrieve track")
		responses.RespondNotFound(c, err, "Track not found")
		return
	}

	// Update user data
	// TODO: Track user data
	// track.UserData = userData

	// Update the track
	updatedTrack, err := h.trackService.Update(ctx, track)
	if err != nil {
		log.Error().Err(err).
			Uint64("trackID", trackID).
			Msg("Failed to update track")
		responses.RespondInternalError(c, err, "Failed to update track")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("trackID", trackID).
		Msg("Track user data updated successfully")
	responses.RespondOK(c, updatedTrack, "Track updated successfully")
}
