// handlers/client_media_series.go
package handlers

import (
	"strconv"
	"suasor/clients/media/types"
	clienttypes "suasor/clients/types"
	"suasor/services"
	_ "suasor/types/models"
	"suasor/types/responses"
	"suasor/utils/logger"

	"github.com/gin-gonic/gin"
)

type ClientSeriesHandler[T clienttypes.ClientMediaConfig] interface {
	GetSeriesByActor(c *gin.Context)
	GetSeriesByCreator(c *gin.Context)
	GetSeriesByRating(c *gin.Context)
	GetLatestSeriesByAdded(c *gin.Context)
	GetPopularSeries(c *gin.Context)
	GetTopRatedSeries(c *gin.Context)
	SearchSeries(c *gin.Context)
	GetSeasonsBySeriesID(c *gin.Context)
	GetEpisodesBySeriesID(c *gin.Context)
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

// GetSeriesByActor godoc
//
//	@Summary		Get series by actor
//	@Description	Retrieves TV series featuring a specific actor
//	@Tags			series, clients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			actor	path		string															true	"Actor name"
//	@Success		200		{object}	responses.APIResponse[responses.MediaItemList[types.Series]]	"Series retrieved"
//	@Failure		401		{object}	responses.ErrorResponse[responses.ErrorDetails]					"Unauthorized"
//	@Failure		500		{object}	responses.ErrorResponse[responses.ErrorDetails]					"Server error"
//	@Router			/client/{clientID}/media/series/actor/{actor} [get]
func (h *clientSeriesHandler[T]) GetSeriesByActor(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting series by actor")

	// Get authenticated user ID
	uid, ok := checkUserAccess(c)
	if !ok {
		return
	}

	actor, ok := checkRequiredStringParam(c, "actor", "Actor name is required")
	if !ok {
		return
	}

	log.Info().
		Uint64("userID", uid).
		Str("actor", actor).
		Msg("Retrieving series by actor")

	series, err := h.seriesService.GetSeriesByActor(ctx, uid, actor)
	if handleServiceError(c, err,
		"Failed to retrieve series by actor",
		"No series found with this actor",
		"Failed to retrieve series") {
		return
	}

	log.Info().
		Uint64("userID", uid).
		Str("actor", actor).
		Int("count", len(series)).
		Msg("Series retrieved successfully")
	responses.RespondMediaItemListOK(c, series, "Series retrieved successfully")
}

// GetSeriesByCreator godoc
//
//	@Summary		Get series by creator
//	@Description	Retrieves TV series by a specific creator/director
//	@Tags			series, clients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			creator	path		string															true	"Creator name"
//	@Success		200		{object}	responses.APIResponse[responses.MediaItemList[types.Series]]	"Series retrieved"
//	@Failure		401		{object}	responses.ErrorResponse[responses.ErrorDetails]					"Unauthorized"
//	@Failure		500		{object}	responses.ErrorResponse[responses.ErrorDetails]					"Server error"
//	@Router			/client/{clientID}/media/series/creator/{creator} [get]
func (h *clientSeriesHandler[T]) GetSeriesByCreator(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
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
	responses.RespondMediaItemListOK(c, series, "Series retrieved successfully")
}

// GetSeriesByRating godoc
//
//	@Summary		Get series by rating range
//	@Description	Retrieves TV series with ratings within the specified range
//	@Tags			series, clients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			min	query		number															true	"Minimum rating"
//	@Param			max	query		number															true	"Maximum rating"
//	@Success		200	{object}	responses.APIResponse[responses.MediaItemList[types.Series]]	"Series retrieved"
//	@Failure		400	{object}	responses.ErrorResponse[responses.ErrorDetails]					"Invalid rating parameters"
//	@Failure		401	{object}	responses.ErrorResponse[responses.ErrorDetails]					"Unauthorized"
//	@Failure		500	{object}	responses.ErrorResponse[responses.ErrorDetails]					"Server error"
//	@Router			/client/{clientID}/media/series/rating [get]
func (h *clientSeriesHandler[T]) GetSeriesByRating(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
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
	responses.RespondMediaItemListOK(c, series, "Series retrieved successfully")
}

// GetLatestSeriesByAdded godoc
//
//	@Summary		Get latest series by added date
//	@Description	Retrieves the most recently added TV series
//	@Tags			series, clients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			count	path		int																true	"Number of series to retrieve"
//	@Success		200		{object}	responses.APIResponse[responses.MediaItemList[types.Series]]	"Series retrieved"
//	@Failure		400		{object}	responses.ErrorResponse[responses.ErrorDetails]					"Invalid count"
//	@Failure		401		{object}	responses.ErrorResponse[responses.ErrorDetails]					"Unauthorized"
//	@Failure		500		{object}	responses.ErrorResponse[responses.ErrorDetails]					"Server error"
//	@Router			/client/{clientID}/media/series/latest/{count} [get]
func (h *clientSeriesHandler[T]) GetLatestSeriesByAdded(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
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
	responses.RespondMediaItemListOK(c, series, "Series retrieved successfully")
}

// GetPopularSeries godoc
//
//	@Summary		Get popular series
//	@Description	Retrieves most popular TV series
//	@Tags			series, clients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			count	path		int																true	"Number of series to retrieve"
//	@Success		200		{object}	responses.APIResponse[responses.MediaItemList[types.Series]]	"Series retrieved"
//	@Failure		400		{object}	responses.ErrorResponse[responses.ErrorDetails]					"Invalid count"
//	@Failure		401		{object}	responses.ErrorResponse[responses.ErrorDetails]					"Unauthorized"
//	@Failure		500		{object}	responses.ErrorResponse[responses.ErrorDetails]					"Server error"
//	@Router			/client/{clientID}/media/series/popular/{count} [get]
func (h *clientSeriesHandler[T]) GetPopularSeries(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
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
	responses.RespondMediaItemListOK(c, series, "Series retrieved successfully")
}

// GetTopRatedSeries godoc
//
//	@Summary		Get top rated series
//	@Description	Retrieves the highest rated TV series
//	@Tags			series, clients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			count	path		int																true	"Number of series to retrieve"
//	@Success		200		{object}	responses.APIResponse[responses.MediaItemList[types.Series]]	"Series retrieved"
//	@Failure		400		{object}	responses.ErrorResponse[responses.ErrorDetails]					"Invalid count"
//	@Failure		401		{object}	responses.ErrorResponse[responses.ErrorDetails]					"Unauthorized"
//	@Failure		500		{object}	responses.ErrorResponse[responses.ErrorDetails]					"Server error"
//	@Router			/client/{clientID}/media/series/top-rated/{count} [get]
func (h *clientSeriesHandler[T]) GetTopRatedSeries(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting top rated series")

	// Get authenticated user ID
	uid, _ := checkUserAccess(c)

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
	responses.RespondMediaItemListOK(c, series, "Series retrieved successfully")
}

// SearchSeries godoc
//
//	@Summary		Search series
//	@Description	Search for TV series across all connected clients
//	@Tags			series, clients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			q	query		string															true	"Search query"
//	@Success		200	{object}	responses.APIResponse[responses.MediaItemList[types.Series]]	"Series retrieved"
//	@Failure		400	{object}	responses.ErrorResponse[responses.ErrorDetails]					"Invalid query"
//	@Failure		401	{object}	responses.ErrorResponse[responses.ErrorDetails]					"Unauthorized"
//	@Failure		500	{object}	responses.ErrorResponse[responses.ErrorDetails]					"Server error"
//	@Router			/client/{clientID}/media/series/search [get]
func (h *clientSeriesHandler[T]) SearchSeries(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Searching series")

	// Get authenticated user ID
	uid, _ := checkUserAccess(c)

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

	options := types.QueryOptions{
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
	responses.RespondMediaItemListOK(c, series, "Series retrieved successfully")
}

// GetSeasonsBySeriesID godoc
//
//	@Summary		Get seasons for a series
//	@Description	Retrieves all seasons for a specific TV series
//	@Tags			series, clients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			clientID	path		int															true	"Client ID"
//	@Param			clientItemID	path		string														true	"Series ID"
//	@Success		200			{object}	responses.APIResponse[responses.MediaItemList[types.Series]]	"Series retrieved"
//	@Failure		400			{object}	responses.ErrorResponse[responses.ErrorDetails]				"Invalid client ID"
//	@Failure		401			{object}	responses.ErrorResponse[responses.ErrorDetails]				"Unauthorized"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]				"Server error"
//	@Router			/client/{clientID}/media/series/{clientItemID}/seasons [get]
func (h *clientSeriesHandler[T]) GetSeasonsBySeriesID(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
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
	clientID, _ := checkItemID(c, "clientID")
	seriesID := c.Param("clientItemID")

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
	responses.RespondMediaItemListOK(c, seasons, "Seasons retrieved successfully")
}

// GetEpisodesBySeriesID godoc
//
//	@Summary		Get episodes by series ID
//	@Description	Retrieves all episodes for a specific series
//	@Tags			series, client
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			clientID	path		int																true	"Client ID"
//	@Param			clientItemID	path		string														true	"Series ID"
//	@Param			userId			query		int																true	"User ID"
//	@Success		200			{object}	responses.APIResponse[responses.MediaItemList[types.Episode]]	"Episodes retrieved successfully"
//	@Failure		400			{object}	responses.ErrorResponse[error]								"Invalid client ID"
//	@Failure		401			{object}	responses.ErrorResponse[error]								"Unauthorized"
//	@Failure		500			{object}	responses.ErrorResponse[error]								"Server error"
//	@Router			/client/{clientID}/media/series/{clientItemID}/episodes [get]
func (h *clientSeriesHandler[T]) GetEpisodesBySeriesID(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting episodes by series ID")

	// Get authenticated user ID
	uid, ok := checkUserAccess(c)
	if !ok {
		return
	}

	// Parse client ID from URL
	clientID, _ := checkItemID(c, "clientID")
	seriesID := c.Param("clientItemID")

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("seriesID", seriesID).
		Msg("Retrieving episodes by series ID")

	episodes, err := h.seriesService.GetEpisodesBySeriesID(ctx, clientID, seriesID)
	if handleServiceError(c, err, "Failed to retrieve episodes by series ID", "", "Failed to retrieve episodes") {
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("seriesID", seriesID).
		Int("episodeCount", len(episodes)).
		Msg("Episodes retrieved successfully")
	responses.RespondMediaItemListOK(c, episodes, "Episodes retrieved successfully")
}

// GetEpisodesBySeasonID godoc
//
//	@Summary		Get episodes for a season
//	@Description	Retrieves all episodes for a specific season of a series
//	@Tags			series, client
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			clientID	path		int																true	"Client ID"
//	@Param			clientItemID	path		string														true	"Series ID"
//	@Param			seasonNumber	path		int																true	"Season number"
//	@Param			userId			query		int																true	"User ID"
//	@Success		200			{object}	responses.APIResponse[responses.MediaItemList[types.Episode]]	"Episodes retrieved successfully"
//	@Failure		400			{object}	responses.ErrorResponse[error]								"Invalid client ID"
//	@Failure		401			{object}	responses.ErrorResponse[error]								"Unauthorized"
//	@Failure		500			{object}	responses.ErrorResponse[error]								"Server error"
//	@Router			/client/{clientID}/media/series/{clientItemID}/season/{seasonNumber}/episodes [get]
func (h *clientSeriesHandler[T]) GetEpisodesBySeasonNbr(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting episodes by season ID")

	// Get authenticated user ID
	uid, ok := checkUserAccess(c)
	if !ok {
		return
	}

	// Parse client ID from URL
	clientID, _ := checkItemID(c, "clientID")
	seriesID := c.Param("clientItemID")
	seasonNumber, err := strconv.Atoi(c.Param("seasonNumber"))

	if err != nil {
		log.Error().Err(err).Str("seasonNumber", c.Param("seasonNumber")).Msg("Invalid season number format")
		responses.RespondBadRequest(c, err, "Invalid season number")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("seriesID", seriesID).
		Int("seasonNumber", seasonNumber).
		Msg("Retrieving episodes by season ID")

	episodes, err := h.seriesService.GetEpisodesBySeasonNbr(ctx, clientID, seriesID, seasonNumber)
	if handleServiceError(c, err, "Failed to retrieve episodes by season ID", "", "Failed to retrieve episodes") {
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("seriesID", seriesID).
		Int("seasonNumber", seasonNumber).
		Int("episodeCount", len(episodes)).
		Msg("Episodes retrieved successfully")
	responses.RespondMediaItemListOK(c, episodes, "Episodes retrieved successfully")
}
