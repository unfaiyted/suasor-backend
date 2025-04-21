// handlers/client_media_series.go
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

type ClientSeriesHandler[T clienttypes.ClientMediaConfig] interface {
	GetSeriesByID(c *gin.Context)
	GetSeriesByGenre(c *gin.Context)
	GetSeriesByYear(c *gin.Context)
	GetSeriesByActor(c *gin.Context)
	GetSeriesByCreator(c *gin.Context)
	GetSeriesByRating(c *gin.Context)
	GetLatestSeriesByAdded(c *gin.Context)
	GetPopularSeries(c *gin.Context)
	GetTopRatedSeries(c *gin.Context)
	SearchSeries(c *gin.Context)
	GetSeasonsBySeriesID(c *gin.Context)
}

// clientSeriesHandler handles series-related operations for media clients
type clientSeriesHandler[T clienttypes.ClientMediaConfig] struct {
	CoreSeriesHandler
	seriesService services.ClientSeriesService[T]
}

// NewclientSeriesHandler creates a new media client series handler
func NewClientSeriesHandler[T clienttypes.ClientMediaConfig](
	coreHandler CoreSeriesHandler,
	seriesService services.ClientSeriesService[T],
) *clientSeriesHandler[T] {
	return &clientSeriesHandler[T]{
		CoreSeriesHandler: coreHandler,
		seriesService:     seriesService,
	}
}

// GetSeriesByID godoc
// @Summary Get series by ID
// @Description Retrieves a specific TV series from the client by ID
// @Tags series
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param seriesID path string true "Series ID"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Series]] "Movies retrieved"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid client ID"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /clients/media/{clientID}/series/{seriesID} [get]
func (h *clientSeriesHandler[T]) GetSeriesByID(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	log.Info().Msg("Getting series by ID")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access series without authentication")
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

	seriesID := c.Param("seriesID")

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("seriesID", seriesID).
		Msg("Retrieving series by ID")

	series, err := h.seriesService.GetSeriesByID(ctx, clientID, seriesID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Str("seriesID", seriesID).
			Msg("Failed to retrieve series")
		responses.RespondInternalError(c, err, "Failed to retrieve series")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("seriesID", seriesID).
		Msg("Series retrieved successfully")
	responses.RespondOK(c, series, "Series retrieved successfully")
}

// GetSeriesByGenre godoc
// @Summary Get series by genre
// @Description Retrieves TV series from all connected clients that match the specified genre
// @Tags series
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param genre path string true "Genre name"
// @Success 200 {object} responses.APIResponse[[]responses.MediaItemResponse] "Series retrieved"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /series/genre/{genre} [get]
func (h *clientSeriesHandler[T]) GetSeriesByGenre(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	log.Info().Msg("Getting series by genre")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access series without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)
	genre := c.Param("genre")

	log.Info().
		Uint64("userID", uid).
		Str("genre", genre).
		Msg("Retrieving series by genre")

	series, err := h.seriesService.GetSeriesByGenre(ctx, uid, genre)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Str("genre", genre).
			Msg("Failed to retrieve series by genre")
		responses.RespondInternalError(c, err, "Failed to retrieve series")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Str("genre", genre).
		Int("count", len(series)).
		Msg("Series retrieved successfully")
	responses.RespondOK(c, series, "Series retrieved successfully")
}

// GetSeriesByYear godoc
// @Summary Get series by release year
// @Description Retrieves TV series from all connected clients that were released in the specified year
// @Tags series
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param year path int true "Release year"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Series]] "Series retrieved"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid year"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /series/year/{year} [get]
func (h *clientSeriesHandler[T]) GetSeriesByYear(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	log.Info().Msg("Getting series by year")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access series without authentication")
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
		Msg("Retrieving series by year")

	series, err := h.seriesService.GetSeriesByYear(ctx, uid, year)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Int("year", year).
			Msg("Failed to retrieve series by year")
		responses.RespondInternalError(c, err, "Failed to retrieve series")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("year", year).
		Int("count", len(series)).
		Msg("Series retrieved successfully")
	responses.RespondOK(c, series, "Series retrieved successfully")
}

// GetSeriesByActor godoc
// @Summary Get series by actor
// @Description Retrieves TV series featuring a specific actor
// @Tags series
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param actor path string true "Actor name"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Series]] "Series retrieved"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /series/actor/{actor} [get]
func (h *clientSeriesHandler[T]) GetSeriesByActor(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	log.Info().Msg("Getting series by actor")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access series without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)
	actor := c.Param("actor")

	log.Info().
		Uint64("userID", uid).
		Str("actor", actor).
		Msg("Retrieving series by actor")

	series, err := h.seriesService.GetSeriesByActor(ctx, uid, actor)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Str("actor", actor).
			Msg("Failed to retrieve series by actor")
		responses.RespondInternalError(c, err, "Failed to retrieve series")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Str("actor", actor).
		Int("count", len(series)).
		Msg("Series retrieved successfully")
	responses.RespondOK(c, series, "Series retrieved successfully")
}

// GetSeriesByCreator godoc
// @Summary Get series by creator
// @Description Retrieves TV series by a specific creator/director
// @Tags series
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param creator path string true "Creator name"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Series]] "Series retrieved"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /series/creator/{creator} [get]
func (h *clientSeriesHandler[T]) GetSeriesByCreator(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	log.Info().Msg("Getting series by creator")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access series without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)
	creator := c.Param("creator")

	log.Info().
		Uint64("userID", uid).
		Str("creator", creator).
		Msg("Retrieving series by creator")

	series, err := h.seriesService.GetSeriesByCreator(ctx, uid, creator)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Str("creator", creator).
			Msg("Failed to retrieve series by creator")
		responses.RespondInternalError(c, err, "Failed to retrieve series")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Str("creator", creator).
		Int("count", len(series)).
		Msg("Series retrieved successfully")
	responses.RespondOK(c, series, "Series retrieved successfully")
}

// GetSeriesByRating godoc
// @Summary Get series by rating range
// @Description Retrieves TV series with ratings within the specified range
// @Tags series
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param min query number true "Minimum rating"
// @Param max query number true "Maximum rating"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Series]] "Series retrieved"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid rating parameters"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /series/rating [get]
func (h *clientSeriesHandler[T]) GetSeriesByRating(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	log.Info().Msg("Getting series by rating")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access series without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	minRating, err := strconv.ParseFloat(c.Query("min"), 64)
	if err != nil {
		log.Error().Err(err).Str("min", c.Query("min")).Msg("Invalid minimum rating format")
		responses.RespondBadRequest(c, err, "Invalid minimum rating")
		return
	}

	maxRating, err := strconv.ParseFloat(c.Query("max"), 64)
	if err != nil {
		log.Error().Err(err).Str("max", c.Query("max")).Msg("Invalid maximum rating format")
		responses.RespondBadRequest(c, err, "Invalid maximum rating")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Float64("minRating", minRating).
		Float64("maxRating", maxRating).
		Msg("Retrieving series by rating range")

	series, err := h.seriesService.GetSeriesByRating(ctx, uid, minRating, maxRating)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Float64("minRating", minRating).
			Float64("maxRating", maxRating).
			Msg("Failed to retrieve series by rating")
		responses.RespondInternalError(c, err, "Failed to retrieve series")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Float64("minRating", minRating).
		Float64("maxRating", maxRating).
		Int("count", len(series)).
		Msg("Series retrieved successfully")
	responses.RespondOK(c, series, "Series retrieved successfully")
}

// GetLatestSeriesByAdded godoc
// @Summary Get latest series by added date
// @Description Retrieves the most recently added TV series
// @Tags series
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param count path int true "Number of series to retrieve"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Series]] "Series retrieved"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid count"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /series/latest/{count} [get]
func (h *clientSeriesHandler[T]) GetLatestSeriesByAdded(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	log.Info().Msg("Getting latest series by added date")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access series without authentication")
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
		Msg("Retrieving latest series by added date")

	series, err := h.seriesService.GetLatestSeriesByAdded(ctx, uid, count)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Int("count", count).
			Msg("Failed to retrieve latest series")
		responses.RespondInternalError(c, err, "Failed to retrieve series")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", count).
		Int("seriesReturned", len(series)).
		Msg("Latest series retrieved successfully")
	responses.RespondOK(c, series, "Series retrieved successfully")
}

// GetPopularSeries godoc
// @Summary Get popular series
// @Description Retrieves most popular TV series
// @Tags series
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param count path int true "Number of series to retrieve"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Series]] "Series retrieved"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid count"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /series/popular/{count} [get]
func (h *clientSeriesHandler[T]) GetPopularSeries(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	log.Info().Msg("Getting popular series")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access series without authentication")
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
		Msg("Retrieving popular series")

	series, err := h.seriesService.GetPopularSeries(ctx, uid, count)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Int("count", count).
			Msg("Failed to retrieve popular series")
		responses.RespondInternalError(c, err, "Failed to retrieve series")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", count).
		Int("seriesReturned", len(series)).
		Msg("Popular series retrieved successfully")
	responses.RespondOK(c, series, "Series retrieved successfully")
}

// GetTopRatedSeries godoc
// @Summary Get top rated series
// @Description Retrieves the highest rated TV series
// @Tags series
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param count path int true "Number of series to retrieve"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Series]] "Series retrieved"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid count"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /series/top-rated/{count} [get]
func (h *clientSeriesHandler[T]) GetTopRatedSeries(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	log.Info().Msg("Getting top rated series")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access series without authentication")
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
		Msg("Retrieving top rated series")

	series, err := h.seriesService.GetTopRatedSeries(ctx, uid, count)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Int("count", count).
			Msg("Failed to retrieve top rated series")
		responses.RespondInternalError(c, err, "Failed to retrieve series")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", count).
		Int("seriesReturned", len(series)).
		Msg("Top rated series retrieved successfully")
	responses.RespondOK(c, series, "Series retrieved successfully")
}

// SearchSeries godoc
// @Summary Search series
// @Description Search for TV series across all connected clients
// @Tags series
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param q query string true "Search query"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Series]] "Series retrieved"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid query"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /series/search [get]
func (h *clientSeriesHandler[T]) SearchSeries(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	log.Info().Msg("Searching series")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to search series without authentication")
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
		Msg("Searching series")
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

	series, err := h.seriesService.SearchSeries(ctx, clientID, &options)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Str("query", query).
			Msg("Failed to search series")
		responses.RespondInternalError(c, err, "Failed to search series")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Str("query", query).
		Int("resultsCount", len(series)).
		Msg("Series search completed successfully")
	responses.RespondOK(c, series, "Series retrieved successfully")
}

// GetSeasonsBySeriesID godoc
// @Summary Get seasons for a series
// @Description Retrieves all seasons for a specific TV series
// @Tags series
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param seriesID path string true "Series ID"
// @Success 200 {object} responses.APIResponse[models.MediaItem[mediatypes.Series]] "Series retrieved"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid client ID"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /clients/media/{clientID}/series/{seriesID}/seasons [get]
func (h *clientSeriesHandler[T]) GetSeasonsBySeriesID(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	log.Info().Msg("Getting seasons by series ID")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access seasons without authentication")
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

	seriesID := c.Param("seriesID")

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("seriesID", seriesID).
		Msg("Retrieving seasons for series")

	seasons, err := h.seriesService.GetSeasonsBySeriesID(ctx, clientID, seriesID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Str("seriesID", seriesID).
			Msg("Failed to retrieve seasons")
		responses.RespondInternalError(c, err, "Failed to retrieve seasons")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("seriesID", seriesID).
		Int("seasonsCount", len(seasons)).
		Msg("Seasons retrieved successfully")
	responses.RespondOK(c, seasons, "Seasons retrieved successfully")
}

func createSeriesMediaItem[T mediatypes.Series](clientID uint64, clientType clienttypes.ClientMediaType, externalID string, data mediatypes.Series) models.MediaItem[mediatypes.Series] {
	mediaItem := models.MediaItem[mediatypes.Series]{
		Type:        mediatypes.MediaTypeSeries,
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
