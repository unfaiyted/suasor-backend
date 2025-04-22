// handlers/client_media_music.go
package handlers

import (
	"strconv"
	mediatypes "suasor/clients/media/types"
	clienttypes "suasor/clients/types"
	"suasor/services"
	"suasor/types/models"
	"suasor/types/responses"
	"suasor/utils/logger"

	"github.com/gin-gonic/gin"
)

type ClientMusicHandler[T clienttypes.ClientMediaConfig] interface {
	CoreMusicHandler
	// ID-based retrieval operations
	GetClientTrackByID(c *gin.Context)
	GetClientAlbumByID(c *gin.Context)
	GetClientArtistByID(c *gin.Context)

	// Relationship operations
	GetClientTracksByAlbum(c *gin.Context)
	GetClientAlbumsByArtist(c *gin.Context)
	GetClientSimilarArtists(c *gin.Context)
	GetClientSimilarTracks(c *gin.Context)

	// Filter/category operations
	GetClientArtistsByGenre(c *gin.Context)
	GetClientAlbumsByGenre(c *gin.Context)
	GetClientTracksByGenre(c *gin.Context)
	GetClientAlbumsByYear(c *gin.Context)

	// Recommendations & discovery
	GetClientLatestAlbumsByAdded(c *gin.Context)
	GetClientRecentlyAddedTracks(c *gin.Context)
	GetClientRecentlyPlayedTracks(c *gin.Context)
	GetClientPopularAlbums(c *gin.Context)
	GetClientPopularArtists(c *gin.Context)
	GetClientTopTracks(c *gin.Context)
	GetClientTopAlbums(c *gin.Context)
	GetClientTopArtists(c *gin.Context)

	// User-specific operations
	GetClientFavoriteArtists(c *gin.Context)
	GetClientFavoriteTracks(c *gin.Context)
	GetClientFavoriteAlbums(c *gin.Context)

	// Search operations
	SearchClientMusic(c *gin.Context)

	// Playback operations
	// StartClientTrackPlayback(c *gin.Context)
	// GetClientPlaybackState(c *gin.Context)
	// GetClientPlaybackInfo(c *gin.Context)

	// Playlist operations
	// GetClientUserPlaylists(c *gin.Context)
	// GetClientPlaylistTracks(c *gin.Context)
}

// clientMusicHandler handles music-related operations for media clients
type clientMusicHandler[T clienttypes.ClientMediaConfig] struct {
	CoreMusicHandler
	musicService services.ClientMusicService[T]
}

// NewClientMusicHandler creates a new media client music handler
func NewClientMusicHandler[T clienttypes.ClientMediaConfig](
	musicService services.ClientMusicService[T],
	coreHandler CoreMusicHandler,
) ClientMusicHandler[T] {
	return &clientMusicHandler[T]{
		CoreMusicHandler: coreHandler,
		musicService:     musicService,
	}
}

// GetClientTrackByID godoc
// @Summary Get track by ID from client
// @Description Retrieves a specific music track from the client by ID
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param trackID path string true "Track ID"
// @Success 200 {object} responses.APIResponse[models.MediaItem[mediatypes.Track]] "Track retrieved"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid client ID"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /clients/media/{clientID}/music/tracks/{trackID} [get]
func (h *clientMusicHandler[T]) GetClientTrackByID(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting track by ID")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access track without authentication")
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

	trackID := c.Param("trackID")

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("trackID", trackID).
		Msg("Retrieving track by ID")

	track, err := h.musicService.GetClientTrackByID(ctx, clientID, trackID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Str("trackID", trackID).
			Msg("Failed to retrieve track")
		responses.RespondInternalError(c, err, "Failed to retrieve track")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("trackID", trackID).
		Msg("Track retrieved successfully")
	responses.RespondOK(c, track, "Track retrieved successfully")
}

// GetAlbumByID godoc
// @Summary Get album by ID
// @Description Retrieves a specific music album from the client by ID
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param albumID path string true "Album ID"
// @Success 200 {object} responses.APIResponse[models.MediaItem[mediatypes.Album]] "Album retrieved"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid client ID"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /clients/media/{clientID}/music/albums/{albumID} [get]
func (h *clientMusicHandler[T]) GetClientAlbumByID(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting album by ID")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access album without authentication")
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

	albumID := c.Param("albumID")

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("albumID", albumID).
		Msg("Retrieving album by ID")

	album, err := h.musicService.GetClientAlbumByID(ctx, clientID, albumID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Str("albumID", albumID).
			Msg("Failed to retrieve album")
		responses.RespondInternalError(c, err, "Failed to retrieve album")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("albumID", albumID).
		Msg("Album retrieved successfully")
	responses.RespondOK(c, album, "Album retrieved successfully")
}

// GetArtistByID godoc
// @Summary Get artist by ID
// @Description Retrieves a specific music artist from the client by ID
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param artistID path string true "Artist ID"
// @Success 200 {object} responses.APIResponse[models.MediaItem[mediatypes.Artist]] "Artist retrieved"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid client ID"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /clients/media/{clientID}/music/artists/{artistID} [get]
func (h *clientMusicHandler[T]) GetClientArtistByID(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting artist by ID")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access artist without authentication")
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

	artistID := c.Param("artistID")

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("artistID", artistID).
		Msg("Retrieving artist by ID")

	artist, err := h.musicService.GetClientArtistByID(ctx, clientID, artistID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Str("artistID", artistID).
			Msg("Failed to retrieve artist")
		responses.RespondInternalError(c, err, "Failed to retrieve artist")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("artistID", artistID).
		Msg("Artist retrieved successfully")
	responses.RespondOK(c, artist, "Artist retrieved successfully")
}

// GetTracksByAlbum godoc
// @Summary Get tracks by album
// @Description Retrieves all tracks for a specific album
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param albumID path string true "Album ID"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Track]] "Tracks retrieved"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid client ID"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /clients/media/{clientID}/music/albums/{albumID}/tracks [get]
func (h *clientMusicHandler[T]) GetClientTracksByAlbum(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting tracks by album")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access tracks without authentication")
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

	albumID := c.Param("albumID")

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("albumID", albumID).
		Msg("Retrieving tracks by album")

	tracks, err := h.musicService.GetClientTracksByAlbum(ctx, clientID, albumID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Str("albumID", albumID).
			Msg("Failed to retrieve tracks by album")
		responses.RespondInternalError(c, err, "Failed to retrieve tracks")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("albumID", albumID).
		Int("trackCount", len(tracks)).
		Msg("Tracks retrieved successfully")
	responses.RespondOK(c, tracks, "Tracks retrieved successfully")
}

// GetAlbumsByArtist godoc
// @Summary Get albums by artist
// @Description Retrieves all albums for a specific artist
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param artistID path string true "Artist ID"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Album]] "Albums retrieved"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid client ID"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /clients/media/{clientID}/music/artists/{artistID}/albums [get]
func (h *clientMusicHandler[T]) GetClientAlbumsByArtist(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting albums by artist")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access albums without authentication")
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

	artistID := c.Param("artistID")

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("artistID", artistID).
		Msg("Retrieving albums by artist")

	albums, err := h.musicService.GetClientAlbumsByArtist(ctx, clientID, artistID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Str("artistID", artistID).
			Msg("Failed to retrieve albums by artist")
		responses.RespondInternalError(c, err, "Failed to retrieve albums")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("artistID", artistID).
		Int("albumCount", len(albums)).
		Msg("Albums retrieved successfully")
	responses.RespondOK(c, albums, "Albums retrieved successfully")
}

// GetArtistsByGenre godoc
// @Summary Get artists by genre
// @Description Retrieves artists from all connected clients that match the specified genre
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param genre path string true "Genre name"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Artist]] "Artists retrieved"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /music/artists/genre/{genre} [get]
func (h *clientMusicHandler[T]) GetClientArtistsByGenre(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting artists by genre")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access artists without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)
	genre := c.Param("genre")

	log.Info().
		Uint64("userID", uid).
		Str("genre", genre).
		Msg("Retrieving artists by genre")

	artists, err := h.musicService.GetClientArtistsByGenre(ctx, uid, genre)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Str("genre", genre).
			Msg("Failed to retrieve artists by genre")
		responses.RespondInternalError(c, err, "Failed to retrieve artists")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Str("genre", genre).
		Int("count", len(artists)).
		Msg("Artists retrieved successfully")
	responses.RespondOK(c, artists, "Artists retrieved successfully")
}

// GetAlbumsByGenre godoc
// @Summary Get albums by genre
// @Description Retrieves albums from all connected clients that match the specified genre
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param genre path string true "Genre name"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Album]] "Albums retrieved"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /client/{clientID}/music/albums/genre/{genre} [get]
func (h *clientMusicHandler[T]) GetClientAlbumsByGenre(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting albums by genre")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access albums without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)
	genre := c.Param("genre")

	log.Info().
		Uint64("userID", uid).
		Str("genre", genre).
		Msg("Retrieving albums by genre")
	clientIDStr := c.Param("clientID")
	clientID, err := strconv.ParseUint(clientIDStr, 10, 64)
	if err != nil {
		log.Error().Err(err).Str("clientID", clientIDStr).Msg("Invalid client ID format")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	albums, err := h.musicService.GetClientAlbumsByGenre(ctx, clientID, genre)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Str("genre", genre).
			Msg("Failed to retrieve albums by genre")
		responses.RespondInternalError(c, err, "Failed to retrieve albums")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Str("genre", genre).
		Int("count", len(albums)).
		Msg("Albums retrieved successfully")
	responses.RespondOK(c, albums, "Albums retrieved successfully")
}

// GetTracksByGenre godoc
// @Summary Get tracks by genre
// @Description Retrieves tracks from all connected clients that match the specified genre
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param genre path string true "Genre name"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Track]] "Tracks retrieved"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /client/{clientID}/music/tracks/genre/{genre} [get]
func (h *clientMusicHandler[T]) GetClientTracksByGenre(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting tracks by genre")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access tracks without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)
	genre := c.Param("genre")

	log.Info().
		Uint64("userID", uid).
		Str("genre", genre).
		Msg("Retrieving tracks by genre")
	clientIDStr := c.Param("clientID")
	clientID, err := strconv.ParseUint(clientIDStr, 10, 64)
	if err != nil {
		log.Error().Err(err).Str("clientID", clientIDStr).Msg("Invalid client ID format")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	tracks, err := h.musicService.GetClientTracksByGenre(ctx, clientID, genre)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Str("genre", genre).
			Msg("Failed to retrieve tracks by genre")
		responses.RespondInternalError(c, err, "Failed to retrieve tracks")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Str("genre", genre).
		Int("count", len(tracks)).
		Msg("Tracks retrieved successfully")
	responses.RespondOK(c, tracks, "Tracks retrieved successfully")
}

// GetAlbumsByYear godoc
// @Summary Get albums by release year
// @Description Retrieves albums from all connected clients that were released in the specified year
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param year path int true "Release year"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Album]] "Albums retrieved"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid year"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /music/albums/year/{year} [get]
func (h *clientMusicHandler[T]) GetClientAlbumsByYear(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting albums by year")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access albums without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		log.Error().Err(err).Str("year", c.Param("year")).Msg("Invalid year format")
		responses.RespondBadRequest(c, err, "Invalid year")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("year", year).
		Msg("Retrieving albums by year")
	clientIDStr := c.Param("clientID")
	clientID, err := strconv.ParseUint(clientIDStr, 10, 64)
	if err != nil {
		log.Error().Err(err).Str("clientID", clientIDStr).Msg("Invalid client ID format")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	albums, err := h.musicService.GetClientAlbumsByYear(ctx, clientID, year)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Int("year", year).
			Msg("Failed to retrieve albums by year")
		responses.RespondInternalError(c, err, "Failed to retrieve albums")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("year", year).
		Int("count", len(albums)).
		Msg("Albums retrieved successfully")
	responses.RespondOK(c, albums, "Albums retrieved successfully")
}

// GetLatestAlbumsByAdded godoc
// @Summary Get latest albums by added date
// @Description Retrieves the most recently added albums
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param count path int true "Number of albums to retrieve"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Album]] "Albums retrieved"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid count"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /music/albums/latest/{count} [get]
func (h *clientMusicHandler[T]) GetClientLatestAlbumsByAdded(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting latest albums by added date")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access albums without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	count, err := strconv.Atoi(c.Param("count"))
	if err != nil {
		log.Error().Err(err).Str("count", c.Param("count")).Msg("Invalid count format")
		responses.RespondBadRequest(c, err, "Invalid count")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", count).
		Msg("Retrieving latest albums by added date")
	clientIDStr := c.Param("clientID")
	clientID, err := strconv.ParseUint(clientIDStr, 10, 64)
	if err != nil {
		log.Error().Err(err).Str("clientID", clientIDStr).Msg("Invalid client ID format")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	albums, err := h.musicService.GetClientRecentlyAddedAlbums(ctx, clientID, count)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Int("count", count).
			Msg("Failed to retrieve latest albums")
		responses.RespondInternalError(c, err, "Failed to retrieve albums")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", count).
		Int("albumsReturned", len(albums)).
		Msg("Latest albums retrieved successfully")
	responses.RespondOK(c, albums, "Albums retrieved successfully")
}

// GetPopularAlbums godoc
// @Summary Get popular albums
// @Description Retrieves most popular albums
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param count path int true "Number of albums to retrieve"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Album]] "Albums retrieved"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid count"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /music/albums/popular/{count} [get]
func (h *clientMusicHandler[T]) GetClientPopularAlbums(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting popular albums")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access albums without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	count, err := strconv.Atoi(c.Param("count"))
	if err != nil {
		log.Error().Err(err).Str("count", c.Param("count")).Msg("Invalid count format")
		responses.RespondBadRequest(c, err, "Invalid count")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", count).
		Msg("Retrieving popular albums")

	albums, err := h.musicService.GetClientTopAlbums(ctx, uid, count)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Int("count", count).
			Msg("Failed to retrieve popular albums")
		responses.RespondInternalError(c, err, "Failed to retrieve albums")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", count).
		Int("albumsReturned", len(albums)).
		Msg("Popular albums retrieved successfully")
	responses.RespondOK(c, albums, "Albums retrieved successfully")
}

// GetPopularArtists godoc
// @Summary Get popular artists
// @Description Retrieves most popular artists
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param count path int true "Number of artists to retrieve"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Artist]] "Artists retrieved"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid count"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /music/artists/popular/{count} [get]
func (h *clientMusicHandler[T]) GetClientPopularArtists(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting popular artists")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access artists without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	count, err := strconv.Atoi(c.Param("count"))
	if err != nil {
		log.Error().Err(err).Str("count", c.Param("count")).Msg("Invalid count format")
		responses.RespondBadRequest(c, err, "Invalid count")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", count).
		Msg("Retrieving popular artists")

	artists, err := h.musicService.GetClientTopArtists(ctx, uid, count)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Int("count", count).
			Msg("Failed to retrieve popular artists")
		responses.RespondInternalError(c, err, "Failed to retrieve artists")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", count).
		Int("artistsReturned", len(artists)).
		Msg("Popular artists retrieved successfully")
	responses.RespondOK(c, artists, "Artists retrieved successfully")
}

// SearchMusic godoc
// @Summary Search music (artists, albums, tracks)
// @Description Search for music across all connected clients
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param q query string true "Search query"
// @Success 200 {object} responses.APIResponse[responses.MediaItemResponse] "Music search results retrieved"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid query"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /music/search [get]
func (h *clientMusicHandler[T]) SearchClientMusic(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Searching music")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to search music without authentication")
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
		Msg("Searching music")
	clientIDStr := c.Param("clientID")
	clientID, err := strconv.ParseUint(clientIDStr, 10, 64)
	if err != nil {
		log.Error().Err(err).Str("clientID", clientIDStr).Msg("Invalid client ID format")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}
	options := mediatypes.QueryOptions{
		Query: query,
	}

	results, err := h.musicService.SearchMusic(ctx, clientID, &options)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Str("query", query).
			Msg("Failed to search music")
		responses.RespondInternalError(c, err, "Failed to search music")
		return
	}

	// Create a response with counts for each type
	response := map[string]interface{}{
		"artists": results.Artists,
		"albums":  results.Albums,
		"tracks":  results.Tracks,
	}

	log.Info().
		Uint64("userID", uid).
		Str("query", query).
		Int("artistsCount", len(results.Artists)).
		Int("albumsCount", len(results.Albums)).
		Int("tracksCount", len(results.Tracks)).
		Msg("Music search completed successfully")
	responses.RespondOK(c, response, "Music retrieved successfully")
}

// GetTopTracks godoc
// @Summary Get top tracks from a client
// @Description Retrieves the most popular tracks from a client
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientType path string true "Client Type"
// @Param clientID path int true "Client ID"
// @Param limit query int false "Number of tracks to retrieve (default 10)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Track]] "Tracks retrieved"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid client ID"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /clients/{clientType}/{clientID}/music/tracks/top [get]
func (h *clientMusicHandler[T]) GetClientTopTracks(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting top tracks")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access top tracks without authentication")
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

	// Parse limit parameter (optional)
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Error().Err(err).Str("limit", limitStr).Msg("Invalid limit format")
		responses.RespondBadRequest(c, err, "Invalid limit parameter")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Int("limit", limit).
		Msg("Retrieving top tracks")

	tracks, err := h.musicService.GetClientTopTracks(ctx, clientID, limit)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Int("limit", limit).
			Msg("Failed to retrieve top tracks")
		responses.RespondInternalError(c, err, "Failed to retrieve top tracks")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Int("trackCount", len(tracks)).
		Msg("Top tracks retrieved successfully")
	responses.RespondOK(c, tracks, "Top tracks retrieved successfully")
}

// GetRecentlyAddedTracks godoc
// @Summary Get recently added tracks from a client
// @Description Retrieves the most recently added tracks from a client
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientType path string true "Client Type"
// @Param clientID path int true "Client ID"
// @Param limit query int false "Number of tracks to retrieve (default 10)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Track]] "Tracks retrieved"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid client ID"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /clients/{clientType}/{clientID}/music/tracks/recently-added [get]
func (h *clientMusicHandler[T]) GetClientRecentlyAddedTracks(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting recently added tracks")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access recently added tracks without authentication")
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

	// Parse limit parameter (optional)
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Error().Err(err).Str("limit", limitStr).Msg("Invalid limit format")
		responses.RespondBadRequest(c, err, "Invalid limit parameter")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Int("limit", limit).
		Msg("Retrieving recently added tracks")

	tracks, err := h.musicService.GetClientRecentlyAddedTracks(ctx, clientID, limit)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Int("limit", limit).
			Msg("Failed to retrieve recently added tracks")
		responses.RespondInternalError(c, err, "Failed to retrieve recently added tracks")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Int("trackCount", len(tracks)).
		Msg("Recently added tracks retrieved successfully")
	responses.RespondOK(c, tracks, "Recently added tracks retrieved successfully")
}

// GetTopAlbums godoc
// @Summary Get top albums from a client
// @Description Retrieves the most popular albums from a client
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientType path string true "Client Type"
// @Param clientID path int true "Client ID"
// @Param limit query int false "Number of albums to retrieve (default 10)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Album]] "Albums retrieved"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid client ID"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /clients/{clientType}/{clientID}/music/albums/top [get]
// This method name remains unchanged to match the interface in the router
func (h *clientMusicHandler[T]) GetClientTopAlbums(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting top albums")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access top albums without authentication")
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

	// Parse limit parameter (optional)
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Error().Err(err).Str("limit", limitStr).Msg("Invalid limit format")
		responses.RespondBadRequest(c, err, "Invalid limit parameter")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Int("limit", limit).
		Msg("Retrieving top albums")

	albums, err := h.musicService.GetClientTopAlbums(ctx, clientID, limit)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Int("limit", limit).
			Msg("Failed to retrieve top albums")
		responses.RespondInternalError(c, err, "Failed to retrieve top albums")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Int("albumCount", len(albums)).
		Msg("Top albums retrieved successfully")
	responses.RespondOK(c, albums, "Top albums retrieved successfully")
}

// GetTopArtists godoc
// @Summary Get top artists from a client
// @Description Retrieves the most popular artists from a client
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientType path string true "Client Type"
// @Param clientID path int true "Client ID"
// @Param limit query int false "Number of artists to retrieve (default 10)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Artist]] "Artists retrieved"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid client ID"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /clients/{clientType}/{clientID}/music/artists/top [get]
func (h *clientMusicHandler[T]) GetClientTopArtists(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting top artists")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access top artists without authentication")
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

	// Parse limit parameter (optional)
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Error().Err(err).Str("limit", limitStr).Msg("Invalid limit format")
		responses.RespondBadRequest(c, err, "Invalid limit parameter")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Int("limit", limit).
		Msg("Retrieving top artists")

	artists, err := h.musicService.GetClientTopArtists(ctx, clientID, limit)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Int("limit", limit).
			Msg("Failed to retrieve top artists")
		responses.RespondInternalError(c, err, "Failed to retrieve top artists")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Int("artistCount", len(artists)).
		Msg("Top artists retrieved successfully")
	responses.RespondOK(c, artists, "Top artists retrieved successfully")
}

// GetFavoriteArtists godoc
// @Summary Get favorite artists from a client
// @Description Retrieves the user's favorite artists from a client
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientType path string true "Client Type"
// @Param clientID path int true "Client ID"
// @Param limit query int false "Number of artists to retrieve (default 10)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Artist]] "Artists retrieved"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid client ID"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /clients/{clientType}/{clientID}/music/artists/favorites [get]
func (h *clientMusicHandler[T]) GetClientFavoriteArtists(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting favorite artists")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access favorite artists without authentication")
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

	// Parse limit parameter (optional)
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Error().Err(err).Str("limit", limitStr).Msg("Invalid limit format")
		responses.RespondBadRequest(c, err, "Invalid limit parameter")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Int("limit", limit).
		Msg("Retrieving favorite artists")

	artists, err := h.musicService.GetClientFavoriteArtists(ctx, clientID, limit)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Int("limit", limit).
			Msg("Failed to retrieve favorite artists")
		responses.RespondInternalError(c, err, "Failed to retrieve favorite artists")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Int("artistCount", len(artists)).
		Msg("Favorite artists retrieved successfully")
	responses.RespondOK(c, artists, "Favorite artists retrieved successfully")
}

func createTrackMediaItem(clientID uint64, clientType clienttypes.ClientMediaType, externalID string, data mediatypes.Track) models.MediaItem[mediatypes.Track] {
	mediaItem := models.MediaItem[mediatypes.Track]{
		Type:        mediatypes.MediaTypeTrack,
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

func createAlbumMediaItem(clientID uint64, clientType clienttypes.ClientMediaType, externalID string, data mediatypes.Album) models.MediaItem[mediatypes.Album] {
	mediaItem := models.MediaItem[mediatypes.Album]{
		Type:        mediatypes.MediaTypeAlbum,
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

func createArtistMediaItem(clientID uint64, clientType clienttypes.ClientMediaType, externalID string, data mediatypes.Artist) models.MediaItem[mediatypes.Artist] {
	mediaItem := models.MediaItem[mediatypes.Artist]{
		Type:        mediatypes.MediaTypeArtist,
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

// GetSimilarTracks godoc
// @Summary Get similar tracks
// @Description Retrieves tracks similar to a specific track from a client
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param trackID path string true "Track ID"
// @Param limit query int false "Maximum number of tracks to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Track]] "Similar tracks retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /clients/media/{clientID}/music/tracks/{trackID}/similar [get]
func (h *clientMusicHandler[T]) GetClientSimilarTracks(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting similar tracks")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access similar tracks without authentication")
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

	trackID := c.Param("trackID")

	// Parse limit parameter (optional)
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Error().Err(err).Str("limit", limitStr).Msg("Invalid limit format")
		responses.RespondBadRequest(c, err, "Invalid limit parameter")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("trackID", trackID).
		Int("limit", limit).
		Msg("Retrieving similar tracks")

	tracks, err := h.musicService.GetClientSimilarTracks(ctx, clientID, trackID, limit)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Str("trackID", trackID).
			Msg("Failed to retrieve similar tracks")
		responses.RespondInternalError(c, err, "Failed to retrieve similar tracks")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("trackID", trackID).
		Int("trackCount", len(tracks)).
		Msg("Similar tracks retrieved successfully")
	responses.RespondOK(c, tracks, "Similar tracks retrieved successfully")
}

// GetRecentlyPlayedTracks godoc
// @Summary Get recently played tracks
// @Description Retrieves the user's recently played tracks from a client
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param limit query int false "Maximum number of tracks to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Track]] "Recently played tracks retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /clients/media/{clientID}/music/tracks/recently-played [get]
func (h *clientMusicHandler[T]) GetClientRecentlyPlayedTracks(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting recently played tracks")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access recently played tracks without authentication")
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

	// Parse limit parameter (optional)
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Error().Err(err).Str("limit", limitStr).Msg("Invalid limit format")
		responses.RespondBadRequest(c, err, "Invalid limit parameter")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Int("limit", limit).
		Msg("Retrieving recently played tracks")

	tracks, err := h.musicService.GetClientRecentlyPlayedTracks(ctx, uid, limit)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Int("limit", limit).
			Msg("Failed to retrieve recently played tracks")
		responses.RespondInternalError(c, err, "Failed to retrieve recently played tracks")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Int("trackCount", len(tracks)).
		Msg("Recently played tracks retrieved successfully")
	responses.RespondOK(c, tracks, "Recently played tracks retrieved successfully")
}

// GetFavoriteTracks godoc
// @Summary Get favorite tracks
// @Description Retrieves the user's favorite tracks from a client
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param limit query int false "Maximum number of tracks to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Track]] "Favorite tracks retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /clients/media/{clientID}/music/tracks/favorites [get]
func (h *clientMusicHandler[T]) GetClientFavoriteTracks(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting favorite tracks")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access favorite tracks without authentication")
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

	// Parse limit parameter (optional)
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Error().Err(err).Str("limit", limitStr).Msg("Invalid limit format")
		responses.RespondBadRequest(c, err, "Invalid limit parameter")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Int("limit", limit).
		Msg("Retrieving favorite tracks")

	tracks, err := h.musicService.GetClientFavoriteTracks(ctx, clientID, limit)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Int("limit", limit).
			Msg("Failed to retrieve favorite tracks")
		responses.RespondInternalError(c, err, "Failed to retrieve favorite tracks")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Int("trackCount", len(tracks)).
		Msg("Favorite tracks retrieved successfully")
	responses.RespondOK(c, tracks, "Favorite tracks retrieved successfully")
}

// GetFavoriteAlbums godoc
// @Summary Get favorite albums
// @Description Retrieves the user's favorite albums from a client
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param limit query int false "Maximum number of albums to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Album]] "Favorite albums retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /clients/media/{clientID}/music/albums/favorites [get]
func (h *clientMusicHandler[T]) GetClientFavoriteAlbums(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting favorite albums")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access favorite albums without authentication")
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

	// Parse limit parameter (optional)
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Error().Err(err).Str("limit", limitStr).Msg("Invalid limit format")
		responses.RespondBadRequest(c, err, "Invalid limit parameter")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Int("limit", limit).
		Msg("Retrieving favorite albums")

	albums, err := h.musicService.GetClientFavoriteAlbums(ctx, clientID, limit)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Int("limit", limit).
			Msg("Failed to retrieve favorite albums")
		responses.RespondInternalError(c, err, "Failed to retrieve favorite albums")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Int("albumCount", len(albums)).
		Msg("Favorite albums retrieved successfully")
	responses.RespondOK(c, albums, "Favorite albums retrieved successfully")
}

// StartTrackPlayback godoc
// @Summary Start track playback
// @Description Start playback of a specific track on the client
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param trackID path string true "Track ID"
// @Success 200 {object} responses.APIResponse[any] "Playback started successfully"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /clients/media/{clientID}/music/tracks/{trackID}/play [post]
// func (h *clientMusicHandler[T]) StartClientTrackPlayback(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	log := logger.LoggerFromContext(ctx)
// 	log.Info().Msg("Starting track playback")
//
// 	// Get authenticated user ID
// 	userID, exists := c.Get("userID")
// 	if !exists {
// 		log.Warn().Msg("Attempt to start playback without authentication")
// 		responses.RespondUnauthorized(c, nil, "Authentication required")
// 		return
// 	}
//
// 	uid := userID.(uint64)
//
// 	// Parse client ID from URL
// 	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)
// 	if err != nil {
// 		log.Error().Err(err).Str("clientID", c.Param("clientID")).Msg("Invalid client ID format")
// 		responses.RespondBadRequest(c, err, "Invalid client ID")
// 		return
// 	}
//
// 	trackID := c.Param("trackID")
//
// 	log.Info().
// 		Uint64("userID", uid).
// 		Uint64("clientID", clientID).
// 		Str("trackID", trackID).
// 		Msg("Starting track playback")
//
// 	err = h.musicService.StartTrackPlayback(ctx, uid, clientID, trackID)
// 	if err != nil {
// 		log.Error().Err(err).
// 			Uint64("userID", uid).
// 			Uint64("clientID", clientID).
// 			Str("trackID", trackID).
// 			Msg("Failed to start track playback")
// 		responses.RespondInternalError(c, err, "Failed to start track playback")
// 		return
// 	}
//
// 	log.Info().
// 		Uint64("userID", uid).
// 		Uint64("clientID", clientID).
// 		Str("trackID", trackID).
// 		Msg("Track playback started successfully")
// 	responses.RespondOK(c, http.StatusOK, "Track playback started successfully")
// }

// GetPlaybackState godoc
// @Summary Get playback state
// @Description Get the current playback state of the client
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Success 200 {object} responses.APIResponse[map[string]interface{}] "Playback state retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /clients/media/{clientID}/music/playback [get]
// func (h *clientMusicHandler[T]) GetClientPlaybackState(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	log := logger.LoggerFromContext(ctx)
// 	log.Info().Msg("Getting playback state")
//
// 	// Get authenticated user ID
// 	userID, exists := c.Get("userID")
// 	if !exists {
// 		log.Warn().Msg("Attempt to get playback state without authentication")
// 		responses.RespondUnauthorized(c, nil, "Authentication required")
// 		return
// 	}
//
// 	uid := userID.(uint64)
//
// 	// Parse client ID from URL
// 	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)
// 	if err != nil {
// 		log.Error().Err(err).Str("clientID", c.Param("clientID")).Msg("Invalid client ID format")
// 		responses.RespondBadRequest(c, err, "Invalid client ID")
// 		return
// 	}
//
// 	log.Info().
// 		Uint64("userID", uid).
// 		Uint64("clientID", clientID).
// 		Msg("Getting playback state")
//
// 	state, err := h.musicService.GetPlaybackState(ctx, uid, clientID)
// 	if err != nil {
// 		log.Error().Err(err).
// 			Uint64("userID", uid).
// 			Uint64("clientID", clientID).
// 			Msg("Failed to get playback state")
// 		responses.RespondInternalError(c, err, "Failed to get playback state")
// 		return
// 	}
//
// 	log.Info().
// 		Uint64("userID", uid).
// 		Uint64("clientID", clientID).
// 		Msg("Playback state retrieved successfully")
// 	responses.RespondOK(c, state, "Playback state retrieved successfully")
// }

// GetPlaybackInfo godoc
// @Summary Get playback info
// @Description Get detailed information about the current playback
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Success 200 {object} responses.APIResponse[map[string]interface{}] "Playback info retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /clients/media/{clientID}/music/playback/info [get]
// func (h *clientMusicHandler[T]) GetClientPlaybackInfo(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	log := logger.LoggerFromContext(ctx)
// 	log.Info().Msg("Getting playback info")
//
// 	// Get authenticated user ID
// 	userID, exists := c.Get("userID")
// 	if !exists {
// 		log.Warn().Msg("Attempt to get playback info without authentication")
// 		responses.RespondUnauthorized(c, nil, "Authentication required")
// 		return
// 	}
//
// 	uid := userID.(uint64)
//
// 	// Parse client ID from URL
// 	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)
// 	if err != nil {
// 		log.Error().Err(err).Str("clientID", c.Param("clientID")).Msg("Invalid client ID format")
// 		responses.RespondBadRequest(c, err, "Invalid client ID")
// 		return
// 	}
//
// 	log.Info().
// 		Uint64("userID", uid).
// 		Uint64("clientID", clientID).
// 		Msg("Getting playback info")
//
// 	info, err := h.musicService.GetPlaybackInfo(ctx, uid, clientID)
// 	if err != nil {
// 		log.Error().Err(err).
// 			Uint64("userID", uid).
// 			Uint64("clientID", clientID).
// 			Msg("Failed to get playback info")
// 		responses.RespondInternalError(c, err, "Failed to get playback info")
// 		return
// 	}
//
// 	log.Info().
// 		Uint64("userID", uid).
// 		Uint64("clientID", clientID).
// 		Msg("Playback info retrieved successfully")
// 	responses.RespondOK(c, info, "Playback info retrieved successfully")
// }

// GetUserPlaylists godoc
// @Summary Get user playlists
// @Description Retrieve the user's playlists from a client
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Playlist]] "User playlists retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /clients/media/{clientID}/music/playlists [get]

// GetPlaylistTracks godoc
// @Summary Get playlist tracks
// @Description Retrieve tracks from a specific playlist
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param playlistID path string true "Playlist ID"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Track]] "Playlist tracks retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /clients/media/{clientID}/music/playlists/{playlistID}/tracks [get]
// func (h *clientMusicHandler[T]) GetClientPlaylistTracks(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	log := logger.LoggerFromContext(ctx)
// 	log.Info().Msg("Getting playlist tracks")
//
// 	// Get authenticated user ID
// 	userID, exists := c.Get("userID")
// 	if !exists {
// 		log.Warn().Msg("Attempt to get playlist tracks without authentication")
// 		responses.RespondUnauthorized(c, nil, "Authentication required")
// 		return
// 	}
//
// 	uid := userID.(uint64)
//
// 	// Parse client ID from URL
// 	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)
// 	if err != nil {
// 		log.Error().Err(err).Str("clientID", c.Param("clientID")).Msg("Invalid client ID format")
// 		responses.RespondBadRequest(c, err, "Invalid client ID")
// 		return
// 	}
//
// 	playlistID := c.Param("playlistID")
//
// 	log.Info().
// 		Uint64("userID", uid).
// 		Uint64("clientID", clientID).
// 		Str("playlistID", playlistID).
// 		Msg("Getting playlist tracks")
//
// 	tracks, err := h.musicService.GetPlaylistTracks(ctx, uid, clientID, playlistID)
// 	if err != nil {
// 		log.Error().Err(err).
// 			Uint64("userID", uid).
// 			Uint64("clientID", clientID).
// 			Str("playlistID", playlistID).
// 			Msg("Failed to get playlist tracks")
// 		responses.RespondInternalError(c, err, "Failed to get playlist tracks")
// 		return
// 	}
//
// 	log.Info().
// 		Uint64("userID", uid).
// 		Uint64("clientID", clientID).
// 		Str("playlistID", playlistID).
// 		Int("trackCount", len(tracks)).
// 		Msg("Playlist tracks retrieved successfully")
// 	responses.RespondOK(c, tracks, "Playlist tracks retrieved successfully")
// }

// GetClientSimilarArtists godoc
// @Summary Get similar artists
// @Description Retrieves artists similar to a specific artist from a client
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param artistID path string true "Artist ID"
// @Param limit query int false "Maximum number of artists to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Artist]] "Similar artists retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /clients/media/{clientID}/music/artists/{artistID}/similar [get]
func (h *clientMusicHandler[T]) GetClientSimilarArtists(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting similar artists")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access similar artists without authentication")
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

	artistID := c.Param("artistID")

	// Parse limit parameter (optional)
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Error().Err(err).Str("limit", limitStr).Msg("Invalid limit format")
		responses.RespondBadRequest(c, err, "Invalid limit parameter")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("artistID", artistID).
		Int("limit", limit).
		Msg("Retrieving similar artists")

	artists, err := h.musicService.GetClientSimilarArtists(ctx, clientID, artistID, limit)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Str("artistID", artistID).
			Msg("Failed to retrieve similar artists")
		responses.RespondInternalError(c, err, "Failed to retrieve similar artists")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("artistID", artistID).
		Int("artistCount", len(artists)).
		Msg("Similar artists retrieved successfully")
	responses.RespondOK(c, artists, "Similar artists retrieved successfully")
}
