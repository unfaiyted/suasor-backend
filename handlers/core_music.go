// handlers/core_music.go
package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	mediatypes "suasor/client/media/types"
	"suasor/services"
	"suasor/types/responses"
	"suasor/utils"
)

// CoreMusicHandler handles operations for music items in the database
type CoreMusicHandler struct {
	trackService  services.CoreMediaItemService[*mediatypes.Track]
	albumService  services.CoreMediaItemService[*mediatypes.Album]
	artistService services.CoreMediaItemService[*mediatypes.Artist]
}

// NewCoreMusicHandler creates a new core music handler
func NewCoreMusicHandler(
	trackService services.CoreMediaItemService[*mediatypes.Track],
	albumService services.CoreMediaItemService[*mediatypes.Album],
	artistService services.CoreMediaItemService[*mediatypes.Artist],
) *CoreMusicHandler {
	return &CoreMusicHandler{
		trackService:  trackService,
		albumService:  albumService,
		artistService: artistService,
	}
}

// GetAlbumTracks godoc
// @Summary Get tracks for an album
// @Description Retrieves all tracks for a specific album
// @Tags music
// @Accept json
// @Produce json
// @Param id path int true "Album ID"
// @Param userId query int true "User ID"
// @Success 200 {object} responses.APIResponse[[]mediatypes.Track] "Tracks retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Album not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /music/albums/{id}/tracks [get]
func (h *CoreMusicHandler) GetAlbumTracks(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	albumID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid album ID")
		responses.RespondBadRequest(c, err, "Invalid album ID")
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	log.Debug().
		Uint64("albumID", albumID).
		Uint64("userID", userID).
		Msg("Getting tracks for album")

	// Get the album first to ensure it exists
	album, err := h.albumService.GetByID(ctx, albumID)
	if err != nil {
		log.Error().Err(err).
			Uint64("albumID", albumID).
			Msg("Failed to retrieve album")
		responses.RespondNotFound(c, err, "Album not found")
		return
	}

	// Get tracks from the album data
	tracks := album.Data.Tracks
	if tracks == nil {
		tracks = []*mediatypes.Track{}
	}

	log.Info().
		Uint64("albumID", albumID).
		Int("trackCount", len(tracks)).
		Msg("Tracks retrieved successfully")
	responses.RespondOK(c, tracks, "Tracks retrieved successfully")
}

// GetArtistAlbums godoc
// @Summary Get albums for an artist
// @Description Retrieves all albums for a specific artist
// @Tags music
// @Accept json
// @Produce json
// @Param id path int true "Artist ID"
// @Param userId query int true "User ID"
// @Success 200 {object} responses.APIResponse[[]mediatypes.Album] "Albums retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Artist not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /music/artists/{id}/albums [get]
func (h *CoreMusicHandler) GetArtistAlbums(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	artistID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid artist ID")
		responses.RespondBadRequest(c, err, "Invalid artist ID")
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	log.Debug().
		Uint64("artistID", artistID).
		Uint64("userID", userID).
		Msg("Getting albums for artist")

	// Get the artist first to ensure it exists
	artist, err := h.artistService.GetByID(ctx, artistID)
	if err != nil {
		log.Error().Err(err).
			Uint64("artistID", artistID).
			Msg("Failed to retrieve artist")
		responses.RespondNotFound(c, err, "Artist not found")
		return
	}

	// Get albums for the artist
	albums := artist.Data.Albums
	if albums == nil {
		albums = []*mediatypes.Album{}
	}

	log.Info().
		Uint64("artistID", artistID).
		Int("albumCount", len(albums)).
		Msg("Albums retrieved successfully")
	responses.RespondOK(c, albums, "Albums retrieved successfully")
}

// GetTopTracks godoc
// @Summary Get top tracks
// @Description Retrieves the top tracks based on play count, ratings, etc.
// @Tags music
// @Accept json
// @Produce json
// @Param userId query int true "User ID"
// @Param limit query int false "Maximum number of tracks to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]mediatypes.Track] "Tracks retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /music/tracks/top [get]
func (h *CoreMusicHandler) GetTopTracks(c *gin.Context) {
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
		Msg("Getting top tracks")

	options := &mediatypes.QueryOptions{
		Limit:         limit,
		Sort:          "rating",
		MinimumRating: 7,
		MaximumRating: 10,
		SortOrder:     mediatypes.SortOrderDesc,
		// TODO: Add more logic to better filter
	}

	allTracks, err := h.trackService.Search(ctx, *options)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to retrieve tracks")
		responses.RespondInternalError(c, err, "Failed to retrieve tracks")
		return
	}

	// In a real implementation, you'd sort by play count, rating, etc.
	// For now we'll just limit the results
	var topTracks []*mediatypes.Track
	for i, track := range allTracks {
		if i >= limit {
			break
		}
		topTracks = append(topTracks, track.Data)
	}

	log.Info().
		Uint64("userID", userID).
		Int("trackCount", len(topTracks)).
		Msg("Top tracks retrieved successfully")
	responses.RespondOK(c, topTracks, "Top tracks retrieved successfully")
}

// GetRecentlyAddedTracks godoc
// @Summary Get recently added tracks
// @Description Retrieves tracks that were recently added to the library
// @Tags music
// @Accept json
// @Produce json
// @Param userId query int true "User ID"
// @Param limit query int false "Maximum number of tracks to return (default 10)"
// @Param days query int false "Number of days to look back (default 30)"
// @Success 200 {object} responses.APIResponse[[]mediatypes.Track] "Tracks retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /music/tracks/recently-added [get]
func (h *CoreMusicHandler) GetRecentlyAddedTracks(c *gin.Context) {
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

	days, err := strconv.Atoi(c.DefaultQuery("days", "30"))
	if err != nil {
		days = 30
	}

	log.Debug().
		Uint64("userID", userID).
		Int("limit", limit).
		Int("days", days).
		Msg("Getting recently added tracks")

	// In a real implementation, you'd query tracks ordered by added date
	recentTracks, err := h.trackService.GetRecentItems(ctx, days, limit)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to retrieve recently added tracks")
		responses.RespondInternalError(c, err, "Failed to retrieve tracks")
		return
	}

	// Extract the track data from the media items
	var tracks []*mediatypes.Track
	for _, item := range recentTracks {
		tracks = append(tracks, item.Data)
	}

	log.Info().
		Uint64("userID", userID).
		Int("trackCount", len(tracks)).
		Msg("Recently added tracks retrieved successfully")
	responses.RespondOK(c, tracks, "Recently added tracks retrieved successfully")
}

// GetTopAlbums godoc
// @Summary Get top albums
// @Description Retrieves the top albums based on play count, ratings, etc.
// @Tags music
// @Accept json
// @Produce json
// @Param userId query int true "User ID"
// @Param limit query int false "Maximum number of albums to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]mediatypes.Album] "Albums retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /music/albums/top [get]
func (h *CoreMusicHandler) GetTopAlbums(c *gin.Context) {
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
		Msg("Getting top albums")

	options := &mediatypes.QueryOptions{
		Limit:         limit,
		Sort:          "rating",
		MinimumRating: 7,
		MaximumRating: 10,
		SortOrder:     mediatypes.SortOrderDesc,
		// TODO: Add more logic to better filter
	}
	// This is a placeholder - in a real implementation you'd query albums by popularity
	// Here we'll just get all albums and sort/limit them
	allAlbums, err := h.albumService.Search(ctx, *options)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to retrieve albums")
		responses.RespondInternalError(c, err, "Failed to retrieve albums")
		return
	}

	// In a real implementation, you'd sort by play count, rating, etc.
	// For now we'll just limit the results
	var topAlbums []*mediatypes.Album
	for i, album := range allAlbums {
		if i >= limit {
			break
		}
		topAlbums = append(topAlbums, album.Data)
	}

	log.Info().
		Uint64("userID", userID).
		Int("albumCount", len(topAlbums)).
		Msg("Top albums retrieved successfully")
	responses.RespondOK(c, topAlbums, "Top albums retrieved successfully")
}
