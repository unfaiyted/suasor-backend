// handlers/core_music.go
package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"suasor/types/models"

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
	GetAlbumsByYear(c *gin.Context)
	GetArtistsByGenre(c *gin.Context)
	GetAlbumsByGenre(c *gin.Context)
	GetTracksByGenre(c *gin.Context)
	GetLatestAlbumsByAdded(c *gin.Context)
	GetPopularAlbums(c *gin.Context)
	GetPopularArtists(c *gin.Context)
	SearchMusic(c *gin.Context)
	GetTopTracks(c *gin.Context)
	GetRecentlyAddedTracks(c *gin.Context)
	GetTopAlbums(c *gin.Context)
	GetAlbumTracks(c *gin.Context)
	GetArtistAlbums(c *gin.Context)
	GetAlbumsByArtistID(c *gin.Context)
	GetSimilarArtists(c *gin.Context)
	GetRecentlyAddedMusic(c *gin.Context)
	GetGenreRecommendations(c *gin.Context)
	GetSimilarTracks(c *gin.Context)
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
		coreMusicService: coreMusicService,
		trackService:     trackService,
		albumService:     albumService,
		artistService:    artistService,
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

// GetTrackByID godoc
// @Summary Get track by ID
// @Description Retrieves a track by its ID
// @Tags music
// @Accept json
// @Produce json
// @Param id path int true "Track ID"
// @Success 200 {object} responses.APIResponse[mediatypes.Track] "Track retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Track not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /music/tracks/{id} [get]
func (h *coreMusicHandler) GetTrackByID(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	trackID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid track ID")
		responses.RespondBadRequest(c, err, "Invalid track ID")
		return
	}

	log.Debug().
		Uint64("trackID", trackID).
		Msg("Getting track by ID")

	track, err := h.trackService.GetByID(ctx, trackID)
	if err != nil {
		log.Error().Err(err).
			Uint64("trackID", trackID).
			Msg("Failed to retrieve track")
		responses.RespondNotFound(c, err, "Track not found")
		return
	}

	log.Info().
		Uint64("trackID", trackID).
		Msg("Track retrieved successfully")
	responses.RespondOK(c, track.Data, "Track retrieved successfully")
}

// GetAlbumByID godoc
// @Summary Get album by ID
// @Description Retrieves an album by its ID
// @Tags music
// @Accept json
// @Produce json
// @Param id path int true "Album ID"
// @Success 200 {object} responses.APIResponse[mediatypes.Album] "Album retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Album not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /music/albums/{id} [get]
func (h *coreMusicHandler) GetAlbumByID(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	albumID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid album ID")
		responses.RespondBadRequest(c, err, "Invalid album ID")
		return
	}

	log.Debug().
		Uint64("albumID", albumID).
		Msg("Getting album by ID")

	album, err := h.albumService.GetByID(ctx, albumID)
	if err != nil {
		log.Error().Err(err).
			Uint64("albumID", albumID).
			Msg("Failed to retrieve album")
		responses.RespondNotFound(c, err, "Album not found")
		return
	}

	log.Info().
		Uint64("albumID", albumID).
		Msg("Album retrieved successfully")
	responses.RespondOK(c, album.Data, "Album retrieved successfully")
}

// GetArtistByID godoc
// @Summary Get artist by ID
// @Description Retrieves an artist by their ID
// @Tags music
// @Accept json
// @Produce json
// @Param id path int true "Artist ID"
// @Success 200 {object} responses.APIResponse[mediatypes.Artist] "Artist retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Artist not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /music/artists/{id} [get]
func (h *coreMusicHandler) GetArtistByID(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	artistID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid artist ID")
		responses.RespondBadRequest(c, err, "Invalid artist ID")
		return
	}

	log.Debug().
		Uint64("artistID", artistID).
		Msg("Getting artist by ID")

	artist, err := h.artistService.GetByID(ctx, artistID)
	if err != nil {
		log.Error().Err(err).
			Uint64("artistID", artistID).
			Msg("Failed to retrieve artist")
		responses.RespondNotFound(c, err, "Artist not found")
		return
	}

	log.Info().
		Uint64("artistID", artistID).
		Msg("Artist retrieved successfully")
	responses.RespondOK(c, artist.Data, "Artist retrieved successfully")
}

// GetTracksByAlbum godoc
// @Summary Get tracks by album ID
// @Description Retrieves all tracks for a specific album
// @Tags music
// @Accept json
// @Produce json
// @Param id path int true "Album ID"
// @Success 200 {object} responses.APIResponse[[]mediatypes.Track] "Tracks retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Album not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /music/albums/{id}/tracks [get]
func (h *coreMusicHandler) GetTracksByAlbum(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	albumID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid album ID")
		responses.RespondBadRequest(c, err, "Invalid album ID")
		return
	}

	log.Debug().
		Uint64("albumID", albumID).
		Msg("Getting tracks by album")

	tracks, err := h.coreMusicService.GetTracksByAlbumID(ctx, albumID)
	if err != nil {
		log.Error().Err(err).
			Uint64("albumID", albumID).
			Msg("Failed to get tracks for album")
		responses.RespondInternalError(c, err, "Failed to get tracks")
		return
	}

	// Extract the track data for response
	var trackData []*mediatypes.Track
	for _, track := range tracks {
		trackData = append(trackData, track.Data)
	}

	log.Info().
		Uint64("albumID", albumID).
		Int("trackCount", len(trackData)).
		Msg("Tracks retrieved successfully")
	responses.RespondOK(c, trackData, "Tracks retrieved successfully")
}

// GetAlbumsByArtist godoc
// @Summary Get albums by artist ID
// @Description Retrieves all albums for a specific artist
// @Tags music
// @Accept json
// @Produce json
// @Param id path int true "Artist ID"
// @Success 200 {object} responses.APIResponse[[]mediatypes.Album] "Albums retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Artist not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /music/artists/{id}/albums [get]
func (h *coreMusicHandler) GetAlbumsByArtist(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	artistID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid artist ID")
		responses.RespondBadRequest(c, err, "Invalid artist ID")
		return
	}

	log.Debug().
		Uint64("artistID", artistID).
		Msg("Getting albums by artist")

	albums, err := h.coreMusicService.GetAlbumsByArtistID(ctx, artistID)
	if err != nil {
		log.Error().Err(err).
			Uint64("artistID", artistID).
			Msg("Failed to get albums for artist")
		responses.RespondInternalError(c, err, "Failed to get albums")
		return
	}

	// Extract the album data for response
	var albumData []*mediatypes.Album
	for _, album := range albums {
		albumData = append(albumData, album.Data)
	}

	log.Info().
		Uint64("artistID", artistID).
		Int("albumCount", len(albumData)).
		Msg("Albums retrieved successfully")
	responses.RespondOK(c, albumData, "Albums retrieved successfully")
}

// GetArtistsByGenre godoc
// @Summary Get artists by genre
// @Description Retrieves artists by genre
// @Tags music
// @Accept json
// @Produce json
// @Param genre path string true "Genre name"
// @Param limit query int false "Maximum number of artists to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]mediatypes.Artist] "Artists retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /music/genres/{genre}/artists [get]
func (h *coreMusicHandler) GetArtistsByGenre(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

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
		Msg("Getting artists by genre")

	artists, err := h.coreMusicService.GetArtistsByGenre(ctx, genre, limit)
	if err != nil {
		log.Error().Err(err).
			Str("genre", genre).
			Msg("Failed to get artists by genre")
		responses.RespondInternalError(c, err, "Failed to get artists")
		return
	}

	// Extract the artist data for response
	var artistData []*mediatypes.Artist
	for _, artist := range artists {
		artistData = append(artistData, artist.Data)
	}

	log.Info().
		Str("genre", genre).
		Int("artistCount", len(artistData)).
		Msg("Artists retrieved successfully")
	responses.RespondOK(c, artistData, "Artists retrieved successfully")
}

// GetAlbumsByGenre godoc
// @Summary Get albums by genre
// @Description Retrieves albums by genre
// @Tags music
// @Accept json
// @Produce json
// @Param genre path string true "Genre name"
// @Param limit query int false "Maximum number of albums to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]mediatypes.Album] "Albums retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /music/genres/{genre}/albums [get]
func (h *coreMusicHandler) GetAlbumsByGenre(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

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
		Msg("Getting albums by genre")

	albums, err := h.coreMusicService.GetAlbumsByGenre(ctx, genre, limit)
	if err != nil {
		log.Error().Err(err).
			Str("genre", genre).
			Msg("Failed to get albums by genre")
		responses.RespondInternalError(c, err, "Failed to get albums")
		return
	}

	// Extract the album data for response
	var albumData []*mediatypes.Album
	for _, album := range albums {
		albumData = append(albumData, album.Data)
	}

	log.Info().
		Str("genre", genre).
		Int("albumCount", len(albumData)).
		Msg("Albums retrieved successfully")
	responses.RespondOK(c, albumData, "Albums retrieved successfully")
}

// GetTracksByGenre godoc
// @Summary Get tracks by genre
// @Description Retrieves tracks by genre
// @Tags music
// @Accept json
// @Produce json
// @Param genre path string true "Genre name"
// @Param limit query int false "Maximum number of tracks to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]mediatypes.Track] "Tracks retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /music/genres/{genre}/tracks [get]
func (h *coreMusicHandler) GetTracksByGenre(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

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
		Msg("Getting tracks by genre")

	tracks, err := h.coreMusicService.GetTracksByGenre(ctx, genre, limit)
	if err != nil {
		log.Error().Err(err).
			Str("genre", genre).
			Msg("Failed to get tracks by genre")
		responses.RespondInternalError(c, err, "Failed to get tracks")
		return
	}

	// Extract the track data for response
	var trackData []*mediatypes.Track
	for _, track := range tracks {
		trackData = append(trackData, track.Data)
	}

	log.Info().
		Str("genre", genre).
		Int("trackCount", len(trackData)).
		Msg("Tracks retrieved successfully")
	responses.RespondOK(c, trackData, "Tracks retrieved successfully")
}

// GetLatestAlbumsByAdded godoc
// @Summary Get latest albums by added date
// @Description Retrieves the latest albums added to the library
// @Tags music
// @Accept json
// @Produce json
// @Param limit query int false "Maximum number of albums to return (default 10)"
// @Param days query int false "Number of days to look back (default 30)"
// @Success 200 {object} responses.APIResponse[[]mediatypes.Album] "Albums retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /music/albums/latest [get]
func (h *coreMusicHandler) GetLatestAlbumsByAdded(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	days, err := strconv.Atoi(c.DefaultQuery("days", "30"))
	if err != nil {
		days = 30
	}

	log.Debug().
		Int("limit", limit).
		Int("days", days).
		Msg("Getting latest albums by added date")

	albums, err := h.coreMusicService.GetRecentlyAddedAlbums(ctx, days, limit)
	if err != nil {
		log.Error().Err(err).
			Msg("Failed to get latest albums")
		responses.RespondInternalError(c, err, "Failed to get latest albums")
		return
	}

	// Extract the album data for response
	var albumData []*mediatypes.Album
	for _, album := range albums {
		albumData = append(albumData, album.Data)
	}

	log.Info().
		Int("albumCount", len(albumData)).
		Msg("Latest albums retrieved successfully")
	responses.RespondOK(c, albumData, "Latest albums retrieved successfully")
}

// GetPopularAlbums godoc
// @Summary Get popular albums
// @Description Retrieves the most popular albums based on play count, ratings, etc.
// @Tags music
// @Accept json
// @Produce json
// @Param limit query int false "Maximum number of albums to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]mediatypes.Album] "Albums retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /music/albums/popular [get]
func (h *coreMusicHandler) GetPopularAlbums(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	log.Debug().
		Int("limit", limit).
		Msg("Getting popular albums")

	albums, err := h.coreMusicService.GetMostPlayedAlbums(ctx, limit)
	if err != nil {
		log.Error().Err(err).
			Msg("Failed to get popular albums")
		responses.RespondInternalError(c, err, "Failed to get popular albums")
		return
	}

	// Extract the album data for response
	var albumData []*mediatypes.Album
	for _, album := range albums {
		albumData = append(albumData, album.Data)
	}

	log.Info().
		Int("albumCount", len(albumData)).
		Msg("Popular albums retrieved successfully")
	responses.RespondOK(c, albumData, "Popular albums retrieved successfully")
}

// GetPopularArtists godoc
// @Summary Get popular artists
// @Description Retrieves the most popular artists based on play count, ratings, etc.
// @Tags music
// @Accept json
// @Produce json
// @Param limit query int false "Maximum number of artists to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]mediatypes.Artist] "Artists retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /music/artists/popular [get]
func (h *coreMusicHandler) GetPopularArtists(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	log.Debug().
		Int("limit", limit).
		Msg("Getting popular artists")

	artists, err := h.coreMusicService.GetTopArtists(ctx, limit)
	if err != nil {
		log.Error().Err(err).
			Msg("Failed to get popular artists")
		responses.RespondInternalError(c, err, "Failed to get popular artists")
		return
	}

	// Extract the artist data for response
	var artistData []*mediatypes.Artist
	for _, artist := range artists {
		artistData = append(artistData, artist.Data)
	}

	log.Info().
		Int("artistCount", len(artistData)).
		Msg("Popular artists retrieved successfully")
	responses.RespondOK(c, artistData, "Popular artists retrieved successfully")
}

// SearchMusic godoc
// @Summary Search music
// @Description Search for music items (tracks, albums, artists) by query
// @Tags music
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param type query string false "Media type to search for (track, album, artist)"
// @Param limit query int false "Maximum number of items to return (default 10)"
// @Success 200 {object} responses.APIResponse[models.MediaItemList] "Search results retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /music/search [get]
func (h *coreMusicHandler) SearchMusic(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	query := c.Query("q")
	if query == "" {
		log.Warn().Msg("Search query is required")
		responses.RespondBadRequest(c, nil, "Search query is required")
		return
	}

	mediaType := c.Query("type")
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	log.Debug().
		Str("query", query).
		Str("type", mediaType).
		Int("limit", limit).
		Msg("Searching music")

	// Create query options
	options := mediatypes.QueryOptions{
		Query:     query,
		Limit:     limit,
		MediaType: mediatypes.MediaType(mediaType),
	}

	// Search music library
	results, err := h.coreMusicService.SearchMusicLibrary(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Str("query", query).
			Msg("Failed to search music")
		responses.RespondInternalError(c, err, "Failed to search music")
		return
	}

	log.Info().
		Str("query", query).
		Int("count", results.GetTotalItems()).
		Msg("Music search completed successfully")
	responses.RespondOK(c, results, "Search results retrieved successfully")
}

// GetSimilarTracks godoc
// @Summary Get similar tracks
// @Description Retrieves tracks similar to a specific track
// @Tags music
// @Accept json
// @Produce json
// @Param id path int true "Track ID"
// @Param limit query int false "Maximum number of tracks to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]mediatypes.Track] "Similar tracks retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Track not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /music/tracks/{id}/similar [get]
func (h *coreMusicHandler) GetSimilarTracks(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	trackID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid track ID")
		responses.RespondBadRequest(c, err, "Invalid track ID")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	log.Debug().
		Uint64("trackID", trackID).
		Int("limit", limit).
		Msg("Getting similar tracks")

	// First, check if the track exists
	// track, err := h.trackService.GetByID(ctx, trackID)
	if err != nil {
		log.Error().Err(err).
			Uint64("trackID", trackID).
			Msg("Failed to retrieve track")
		responses.RespondNotFound(c, err, "Track not found")
		return
	}

	// Get similar tracks
	similarTracks, err := h.coreMusicService.GetSimilarTracks(ctx, trackID, limit)
	if err != nil {
		log.Error().Err(err).
			Uint64("trackID", trackID).
			Msg("Failed to get similar tracks")
		responses.RespondInternalError(c, err, "Failed to get similar tracks")
		return
	}

	// Extract the track data for response
	var trackData []*mediatypes.Track
	for _, track := range similarTracks {
		trackData = append(trackData, track.Data)
	}

	log.Info().
		Uint64("trackID", trackID).
		Int("trackCount", len(trackData)).
		Msg("Similar tracks retrieved successfully")
	responses.RespondOK(c, trackData, "Similar tracks retrieved successfully")
}

// GetAlbumsByYear godoc
// @Summary Get albums by release year
// @Description Retrieves albums released in a specific year
// @Tags music
// @Accept json
// @Produce json
// @Param year path int true "Release year"
// @Param limit query int false "Maximum number of albums to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[*mediatypes.Album]] "Albums retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /music/albums/year/{year} [get]
func (h *coreMusicHandler) GetAlbumsByYear(c *gin.Context) {
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

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	log.Debug().
		Int("year", year).
		Int("limit", limit).
		Msg("Getting albums by year")

	// Create query options
	options := &mediatypes.QueryOptions{
		Year:  year,
		Limit: limit,
	}

	// Search albums by year
	albums, err := h.albumService.Search(ctx, *options)
	if err != nil {
		log.Error().Err(err).
			Int("year", year).
			Msg("Failed to get albums by year")
		responses.RespondInternalError(c, err, "Failed to get albums")
		return
	}

	// Filter for items with the specified year
	var filtered []*models.MediaItem[*mediatypes.Album]
	for _, album := range albums {
		if album.ReleaseYear == year {
			filtered = append(filtered, album)
		}

		if len(filtered) >= limit {
			break
		}
	}

	log.Info().
		Int("year", year).
		Int("count", len(filtered)).
		Msg("Albums by year retrieved successfully")
	responses.RespondOK(c, filtered, "Albums retrieved successfully")
}
