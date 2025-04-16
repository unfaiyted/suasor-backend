// handlers/client_media_music.go
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

// ClientMusicHandler handles music-related operations for media clients
type ClientMusicHandler[T clienttypes.ClientMediaConfig] struct {
	musicService services.ClientMusicService[T]
}

// NewClientMusicHandler creates a new media client music handler
func NewClientMusicHandler[T clienttypes.ClientMediaConfig](musicService services.ClientMusicService[T]) *ClientMusicHandler[T] {
	return &ClientMusicHandler[T]{
		musicService: musicService,
	}
}

// GetTrackByID godoc
// @Summary Get track by ID
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
func (h *ClientMusicHandler[T]) GetTrackByID(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
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

	track, err := h.musicService.GetTrackByID(ctx, uid, clientID, trackID)
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
func (h *ClientMusicHandler[T]) GetAlbumByID(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
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

	album, err := h.musicService.GetAlbumByID(ctx, uid, clientID, albumID)
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
func (h *ClientMusicHandler[T]) GetArtistByID(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
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

	artist, err := h.musicService.GetArtistByID(ctx, uid, clientID, artistID)
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
func (h *ClientMusicHandler[T]) GetTracksByAlbum(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
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

	tracks, err := h.musicService.GetTracksByAlbum(ctx, uid, clientID, albumID)
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
func (h *ClientMusicHandler[T]) GetAlbumsByArtist(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
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

	albums, err := h.musicService.GetAlbumsByArtist(ctx, uid, clientID, artistID)
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
func (h *ClientMusicHandler[T]) GetArtistsByGenre(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
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

	artists, err := h.musicService.GetArtistsByGenre(ctx, uid, genre)
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
// @Router /music/albums/genre/{genre} [get]
func (h *ClientMusicHandler[T]) GetAlbumsByGenre(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
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

	albums, err := h.musicService.GetAlbumsByGenre(ctx, uid, genre)
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
// @Router /music/tracks/genre/{genre} [get]
func (h *ClientMusicHandler[T]) GetTracksByGenre(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
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

	tracks, err := h.musicService.GetTracksByGenre(ctx, uid, genre)
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
func (h *ClientMusicHandler[T]) GetAlbumsByYear(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
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

	albums, err := h.musicService.GetAlbumsByYear(ctx, uid, year)
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
func (h *ClientMusicHandler[T]) GetLatestAlbumsByAdded(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
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

	albums, err := h.musicService.GetLatestAlbumsByAdded(ctx, uid, count)
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
func (h *ClientMusicHandler[T]) GetPopularAlbums(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
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

	albums, err := h.musicService.GetPopularAlbums(ctx, uid, count)
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
func (h *ClientMusicHandler[T]) GetPopularArtists(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
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

	artists, err := h.musicService.GetPopularArtists(ctx, uid, count)
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
func (h *ClientMusicHandler[T]) SearchMusic(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
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

	results, err := h.musicService.SearchMusic(ctx, uid, query)
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
