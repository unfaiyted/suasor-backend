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

type CoreMusicHandler interface {
	GetTrackByID(c *gin.Context)
	GetAlbumByID(c *gin.Context)
	GetArtistByID(c *gin.Context)
	GetTracksByAlbum(c *gin.Context)
	GetAlbumsByArtist(c *gin.Context)
	GetArtistsByGenre(c *gin.Context)
	GetAlbumsByGenre(c *gin.Context)
	GetTracksByGenre(c *gin.Context)
	GetLatestAlbumsByAdded(c *gin.Context)
	GetPopularAlbums(c *gin.Context)
	GetPopularArtists(c *gin.Context)
	SearchMusic(c *gin.Context)
}

// coreMusicHandler handles operations for music items in the database
type coreMusicHandler struct {
	coreMusicService services.CoreMusicService
	artistService    services.CoreMediaItemService[*mediatypes.Artist]
	trackService     services.CoreMediaItemService[*mediatypes.Track]
	albumService     services.CoreMediaItemService[*mediatypes.Album]
}

// NewcoreMusicHandler creates a new core music handler
func NewCoreMusicHandler(
	coreMusicService services.CoreMusicService,
	trackService services.CoreMediaItemService[*mediatypes.Track],
	albumService services.CoreMediaItemService[*mediatypes.Album],
	artistService services.CoreMediaItemService[*mediatypes.Artist],
) CoreMusicHandler {
	return &coreMusicHandler{
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
func (h *coreMusicHandler) GetAlbumTracks(c *gin.Context) {
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
func (h *coreMusicHandler) GetArtistAlbums(c *gin.Context) {
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
func (h *coreMusicHandler) GetTopTracks(c *gin.Context) {
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
func (h *coreMusicHandler) GetRecentlyAddedTracks(c *gin.Context) {
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
func (h *coreMusicHandler) GetTopAlbums(c *gin.Context) {
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

// GetAlbumsByArtistID godoc
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
func (h *coreMusicHandler) GetAlbumsByArtistID(c *gin.Context) {
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

// GetSimilarArtists godoc
// @Summary Get similar artists
// @Description Retrieves the similar artists to a specific artist
// @Tags music
// @Accept json
// @Produce json
// @Param id path int true "Artist ID"
// @Param userId query int true "User ID"
// @Success 200 {object} responses.APIResponse[[]mediatypes.Artist] "Similar artists retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Artist not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /music/artists/{id}/similar [get]
func (h *coreMusicHandler) GetSimilarArtists(c *gin.Context) {
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
		Msg("Getting similar artists")

	// Get the artist first to ensure it exists
	artist, err := h.artistService.GetByID(ctx, artistID)
	if err != nil {
		log.Error().Err(err).
			Uint64("artistID", artistID).
			Msg("Failed to retrieve artist")
		responses.RespondNotFound(c, err, "Artist not found")
		return
	}

	// Get similar artists
	similarArtists := artist.Data.SimilarArtists
	if similarArtists == nil {
		similarArtists = []mediatypes.ArtistReference{}
	}

	log.Info().
		Uint64("artistID", artistID).
		Int("artistCount", len(similarArtists)).
		Msg("Similar artists retrieved successfully")
	responses.RespondOK(c, similarArtists, "Similar artists retrieved successfully")
}

// GetMediaItemByExternalSourceID godoc
// @Summary Get media item by external source ID
// @Description Retrieves a media item using its external source ID
// @Tags media
// @Accept json
// @Produce json
// @Param source path string true "External source name (e.g., tmdb, imdb)"
// @Param externalId path string true "External ID from the source"
// @Success 200 {object} responses.APIResponse[models.MediaItem[types.MediaData]] "Media item retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Media item not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /media/external/{source}/{externalId} [get]
func (h *coreMediaItemHandler[T]) GetMediaItemByExternalSourceID(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	source := c.Param("source")
	if source == "" {
		log.Warn().Msg("Source is required")
		responses.RespondBadRequest(c, nil, "Source is required")
		return
	}

	externalID := c.Param("externalId")
	if externalID == "" {
		log.Warn().Msg("External ID is required")
		responses.RespondBadRequest(c, nil, "External ID is required")
		return
	}

	log.Debug().
		Str("source", source).
		Str("externalID", externalID).
		Msg("Getting media item by external ID")

	item, err := h.mediaService.GetByExternalID(ctx, source, externalID)
	if err != nil {
		log.Error().Err(err).
			Str("source", source).
			Str("externalID", externalID).
			Msg("Failed to retrieve media item by external ID")
		responses.RespondNotFound(c, err, "Media item not found")
		return
	}

	log.Info().
		Str("source", source).
		Str("externalID", externalID).
		Msg("Media item retrieved successfully")
	responses.RespondOK(c, item, "Media item retrieved successfully")
}

// GetRecentlyAddedMusic godoc
// @Summary Get recently added music
// @Description Retrieves recently added music
// @Tags music
// @Accept json
// @Produce json
// @Param userId query int true "User ID"
// @Param limit query int false "Maximum number of music items to return (default 10)"
// @Param days query int false "Number of days to look back (default 30)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[types.MusicData]] "Music items retrieved successfully"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /music/recently-added [get]
func (h *coreMusicHandler) GetRecentlyAddedMusic(c *gin.Context) {
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
		Msg("Getting recently added music")

	// Get recently added music
	// TODO: get all 3 types of music
	items, err := h.coreMusicService.GetRecentlyAddedAlbums(ctx, days, limit)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to retrieve recently added music")
		responses.RespondInternalError(c, err, "Failed to retrieve recently added music")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Int("count", len(items)).
		Msg("Recently added music retrieved successfully")
	responses.RespondOK(c, items, "Recently added music retrieved successfully")
}

// GetGenreRecommendations godoc
// @Summary Get genre recommendations
// @Description Get music recommendations based on a genre
// @Tags music
// @Accept json
// @Produce json
// @Param genre path string true "Genre name"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[types.MusicData]] "Music items retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /music/genre/{genre} [get]
func (h *coreMusicHandler) GetGenreRecommendations(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	genre := c.Param("genre")
	if genre == "" {
		log.Warn().Msg("Genre is required")
		responses.RespondBadRequest(c, nil, "Genre is required")
		return
	}

	log.Debug().
		Str("genre", genre).
		Msg("Getting genre recommendations")

	// Create query options
	options := mediatypes.QueryOptions{
		Genre:     genre,
		MediaType: mediatypes.MediaTypeArtist,
	}

	// Search music items
	items, err := h.coreMusicService.SearchMusicLibrary(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Str("genre", genre).
			Msg("Failed to search music items")
		responses.RespondInternalError(c, err, "Failed to search music items")
		return
	}

	log.Info().
		Str("genre", genre).
		Int("count", items.GetTotalItems()).
		Msg("Genre recommendations retrieved successfully")
	responses.RespondOK(c, items, "Genre recommendations retrieved successfully")
}
