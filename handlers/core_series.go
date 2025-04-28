// handlers/series.go
package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"suasor/clients/media/types"
	"suasor/services"
	"suasor/types/models"
	"suasor/types/responses"
	"suasor/utils"
	"suasor/utils/logger"
)

type CoreSeriesHandler interface {
	CoreMediaItemHandler[*types.Series]

	GetAllEpisodes(c *gin.Context)
	GetByCreator(c *gin.Context)

	GetSeasonsBySeriesID(c *gin.Context)
	GetEpisodesBySeasonID(c *gin.Context)
	GetEpisodesBySeriesIDAndSeasonNumber(c *gin.Context)

	GetSeasonWithEpisodes(c *gin.Context)
	GetRecentlyAiredEpisodes(c *gin.Context)
	GetNextUpEpisodes(c *gin.Context)
	GetSeriesByNetwork(c *gin.Context)
}

// coreSeriesHandler handles operations for series items in the database
type coreSeriesHandler struct {
	CoreMediaItemHandler[*types.Series]
	seriesService  services.CoreMediaItemService[*types.Series]
	seasonService  services.CoreMediaItemService[*types.Season]
	episodeService services.CoreMediaItemService[*types.Episode]
}

// NewcoreSeriesHandler creates a new series handler
func NewCoreSeriesHandler(
	coreHandler CoreMediaItemHandler[*types.Series],
	seriesService services.CoreMediaItemService[*types.Series],
	seasonService services.CoreMediaItemService[*types.Season],
	episodeService services.CoreMediaItemService[*types.Episode],
) CoreSeriesHandler {
	return &coreSeriesHandler{
		CoreMediaItemHandler: coreHandler,
		seriesService:        seriesService,
		seasonService:        seasonService,
		episodeService:       episodeService,
	}
}

// GetSeasonsBySeriesID godoc
//
//	@Summary		Get seasons for a series
//	@Description	Retrieves all seasons for a specific series
//	@Tags			series, core
//	@Accept			json
//	@Produce		json
//	@Param			seriesID	path		int											true	"Series ID"
//	@Param			userId		query		int											true	"User ID"
//	@Success		200				{object}	responses.APIResponse[responses.MediaItemList[types.Season]]	"Episodes retrieved successfully"
//	@Failure		400			{object}	responses.ErrorResponse[any]				"Invalid request"
//	@Failure		404			{object}	responses.ErrorResponse[any]				"Series not found"
//	@Failure		500			{object}	responses.ErrorResponse[any]				"Server error"
//	@Router			/media/series/{seriesID}/seasons [get]
func (h *coreSeriesHandler) GetSeasonsBySeriesID(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	seriesID, _ := checkItemID(c, "seriesID")
	userID, _ := checkUserAccess(c)

	log.Debug().
		Uint64("seriesID", seriesID).
		Uint64("userID", userID).
		Msg("Getting seasons for series")

	// Get the series first to ensure it exists
	series, err := h.seriesService.GetByID(ctx, seriesID)
	if handleServiceError(c, err, "Failed to retrieve series", "Series not found", "Series not found") {
		return
	}

	// Get seasons from the series data
	seasons := series.Data.Seasons
	if seasons == nil {
		seasons = []*types.Season{}
	}

	// wrap all in media item list
	var seasonList []*models.MediaItem[*types.Season]
	for _, season := range seasons {
		seasonList = append(seasonList, models.NewMediaItem[*types.Season](types.MediaTypeSeason, season))
	}

	log.Info().
		Uint64("seriesID", seriesID).
		Int("seasonCount", len(seasons)).
		Msg("Seasons retrieved successfully")
	responses.RespondMediaItemListOK(c, seasonList, "Seasons retrieved successfully")
}

// GetEpisodesBySeriesIDAndSeasonNumber godoc
//
//	@Summary		Get episodes for a season
//	@Description	Retrieves all episodes for a specific season of a series
//	@Tags			series, core
//	@Accept			json
//	@Produce		json
//	@Param			seriesID		path		int											true	"Series ID"
//	@Param			seasonNumber	path		int											true	"Season number"
//	@Param			userId			query		int											true	"User ID"
//	@Success		200				{object}	responses.APIResponse[responses.MediaItemList[types.Episode]]	"Episodes retrieved successfully"
//	@Failure		400				{object}	responses.ErrorResponse[any]				"Invalid request"
//	@Failure		404				{object}	responses.ErrorResponse[any]				"Series or season not found"
//	@Failure		500				{object}	responses.ErrorResponse[any]				"Server error"
//	@Router			/media/series/{seriesID}/seasons/{seasonNumber}/episodes [get]
func (h *coreSeriesHandler) GetEpisodesBySeriesIDAndSeasonNumber(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	seriesID, _ := checkItemID(c, "seriesID")
	userID, _ := checkUserAccess(c)

	seasonNumber, ok := checkSeasonNumber(c, "seasonNumber")
	if !ok {
		return
	}

	log.Debug().
		Uint64("seriesID", seriesID).
		Int("seasonNumber", seasonNumber).
		Uint64("userID", userID).
		Msg("Getting episodes for season")

	// Get the series first to ensure it exists
	series, err := h.seriesService.GetByID(ctx, seriesID)
	if handleServiceError(c, err, "Failed to retrieve series", "Series not found", "Series not found") {
		return
	}

	// Find the correct season and get its episodes
	var episodes []*types.Episode
	seasonFound := false

	for _, season := range series.Data.Seasons {
		if season.Number == seasonNumber {
			episodes = season.Episodes
			seasonFound = true
			break
		}
	}

	if !seasonFound {
		log.Warn().
			Uint64("seriesID", seriesID).
			Int("seasonNumber", seasonNumber).
			Msg("Season not found for series")
		responses.RespondNotFound(c, nil, "Season not found")
		return
	}

	if episodes == nil {
		episodes = []*types.Episode{}
	}
	// wrap all in media item list
	var episodeList []*models.MediaItem[*types.Episode]
	for _, episode := range episodes {
		episodeList = append(episodeList, models.NewMediaItem[*types.Episode](types.MediaTypeEpisode, episode))
	}

	log.Info().
		Uint64("seriesID", seriesID).
		Int("seasonNumber", seasonNumber).
		Int("episodeCount", len(episodes)).
		Msg("Episodes retrieved successfully")
	responses.RespondMediaItemListOK(c, episodeList, "Episodes retrieved successfully")
}

// GetContinueWatchingSeries godoc
//
//	@Summary		Get series in progress
//	@Description	Retrieves series that are currently in progress (partially watched)
//	@Tags			series, core
//	@Accept			json
//	@Produce		json
//	@Param			userId	query		int											true	"User ID"
//	@Param			limit	query		int											false	"Maximum number of series to return (default 10)"
//	@Success		200				{object}	responses.APIResponse[responses.MediaItemList[types.Series]]	"Episodes retrieved successfully"
//	@Failure		400		{object}	responses.ErrorResponse[any]				"Invalid request"
//	@Failure		500		{object}	responses.ErrorResponse[any]				"Server error"
//	@Router			/media/series/continue-watching [get]
func (h *coreSeriesHandler) GetContinueWatchingSeries(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

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
		Msg("Getting continue watching series")

	//TODO: This would typically involve checking the play history to find series with episodes that have been partially watched
	// For now, we'll just return a not implemented response since this requires integration with the play history service

	// TODO: Implement by checking the play history for partially watched episodes, then looking up their series
	log.Info().Msg("Continue watching for series not yet implemented")
	responses.RespondNotImplemented(c, nil, "Continue watching for series not yet implemented")
}

// GetAllEpisodes godoc
//
//	@Summary		Get all episodes for a series
//	@Description	Retrieves all episodes across all seasons for a specific series
//	@Tags			series, core
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int											true	"Series ID"
//	@Param			userId	query		int											true	"User ID"
//	@Success		200		{object}	responses.APIResponse[responses.MediaItemList[types.Episode]]	"Episodes retrieved successfully"
//	@Failure		400		{object}	responses.ErrorResponse[any]				"Invalid request"
//	@Failure		404		{object}	responses.ErrorResponse[any]				"Series not found"
//	@Failure		500		{object}	responses.ErrorResponse[any]				"Server error"
//	@Router			/media/series/{id}/episodes [get]
func (h *coreSeriesHandler) GetAllEpisodes(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	seriesID, err := checkItemID(c, "id")
	if err != nil {
		return
	}

	userID, _ := checkUserAccess(c)

	log.Debug().
		Uint64("seriesID", seriesID).
		Uint64("userID", userID).
		Msg("Getting all episodes for series")

	// Get the series first to ensure it exists
	series, err := h.seriesService.GetByID(ctx, seriesID)
	if handleServiceError(c, err, "Failed to retrieve series", "Series not found", "Series not found") {
		return
	}

	// Collect all episodes from all seasons
	var allEpisodes []*types.Episode

	for _, season := range series.Data.Seasons {
		if season.Episodes != nil {
			allEpisodes = append(allEpisodes, season.Episodes...)
		}
	}
	// wrap all in media item list
	var episodeList []*models.MediaItem[*types.Episode]
	for _, episode := range allEpisodes {
		episodeList = append(episodeList, models.NewMediaItem[*types.Episode](types.MediaTypeEpisode, episode))
	}

	log.Info().
		Uint64("seriesID", seriesID).
		Int("totalEpisodes", len(allEpisodes)).
		Msg("All episodes retrieved successfully")
	responses.RespondMediaItemListOK(c, episodeList, "Episodes retrieved successfully")
}

// GetNextUpEpisodes godoc
//
//	@Summary		Get next episodes to watch
//	@Description	Retrieves the next unwatched episodes for series in progress
//	@Tags			series, core
//	@Accept			json
//	@Produce		json
//	@Param			userId	query		int											true	"User ID"
//	@Param			limit	query		int											false	"Maximum number of episodes to return (default 10)"
//	@Success		200		{object}	responses.APIResponse[responses.MediaItemList[types.Episode]]	"Episodes retrieved successfully"
//	@Failure		400		{object}	responses.ErrorResponse[any]				"Invalid request"
//	@Failure		500		{object}	responses.ErrorResponse[any]				"Server error"
//	@Router			/media/series/next-up [get]
func (h *coreSeriesHandler) GetNextUpEpisodes(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	userID, _ := checkUserAccess(c)
	limit := utils.GetLimit(c, 10, 100, true)

	log.Debug().
		Uint64("userID", userID).
		Int("limit", limit).
		Msg("Getting next up episodes")

	// This requires integration with the play history service to determine which episodes have been watched
	// and which ones are next in the sequence
	// For now, we'll just return a not implemented response

	// TODO: Implement this by checking play history for each series, finding the last watched episode,
	// and then determining the next episode in the sequence
	log.Info().Msg("Next up episodes feature not yet implemented")
	responses.RespondNotImplemented(c, nil, "Next up episodes feature not yet implemented")
}

// GetRecentlyAiredEpisodes godoc
//
//	@Summary		Get recently aired episodes
//	@Description	Retrieves episodes that have recently aired based on their air date
//	@Tags			series, core
//	@Accept			json
//	@Produce		json
//	@Param			userId	query		int											true	"User ID"
//	@Param			limit	query		int											false	"Maximum number of episodes to return (default 10)"
//	@Param			days	query		int											false	"Number of days to look back (default 7)"
//	@Success		200		{object}	responses.APIResponse[responses.MediaItemList[types.Episode]]	"Episodes retrieved successfully"
//	@Failure		400		{object}	responses.ErrorResponse[any]				"Invalid request"
//	@Failure		500		{object}	responses.ErrorResponse[any]				"Server error"
//	@Router			/media/series/recently-aired [get]
func (h *coreSeriesHandler) GetRecentlyAiredEpisodes(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	userID, _ := checkUserAccess(c)
	limit := utils.GetLimit(c, 10, 100, true)
	days := checkDaysParam(c, 7)

	log.Debug().
		Uint64("userID", userID).
		Int("limit", limit).
		Int("days", days).
		Msg("Getting recently aired episodes")

	// This would require checking the air dates of episodes across all series
	// For now, we'll just return a not implemented response

	// TODO: Implement this by querying episodes with air dates within the specified time window
	log.Info().Msg("Recently aired episodes feature not yet implemented")
	responses.RespondNotImplemented(c, nil, "Recently aired episodes feature not yet implemented")
}

// GetSeriesByNetwork godoc
//
//	@Summary		Get series by network
//	@Description	Retrieves series from a specific TV network
//	@Tags			series, core
//	@Accept			json
//	@Produce		json
//	@Param			network	path		string										true	"Network name"
//	@Param			limit	query		int											false	"Maximum number of series to return (default 10)"
//	@Param			offset	query		int											false	"Offset for pagination (default 0)"
//	@Success		200		{object}	responses.APIResponse[responses.MediaItemList[types.Series]]	"Series retrieved successfully"
//	@Failure		400		{object}	responses.ErrorResponse[any]				"Invalid request"
//	@Failure		500		{object}	responses.ErrorResponse[any]				"Server error"
//	@Router			/media/series/network/{network} [get]
func (h *coreSeriesHandler) GetSeriesByNetwork(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	network, ok := checkRequiredStringParam(c, "network", "Network name is required")
	if !ok {
		return
	}

	limit := utils.GetLimit(c, 10, 100, true)
	offset := utils.GetOffset(c, 0)

	log.Debug().
		Str("network", network).
		Int("limit", limit).
		Msg("Getting series by network")

	options := types.QueryOptions{
		Network: network,
		Limit:   limit,
		Offset:  offset,
	}

	// Get all series for the user
	allSeries, err := h.seriesService.Search(ctx, options)
	if handleServiceError(c, err, "Failed to retrieve series", "", "Failed to retrieve series") {
		return
	}

	log.Info().
		Str("network", network).
		Msg("Series by network retrieved successfully")

	responses.RespondMediaItemListOK(c, allSeries, "Series retrieved successfully")
}

// GetSeasonWithEpisodes godoc
//
//	@Summary		Get a season and all its episodes
//	@Description	Retrieves a specific season and all its episodes
//	@Tags			series, core
//	@Accept			json
//	@Produce		json
//	@Param			id				path		int											true	"Series ID"
//	@Param			seasonNumber	path		int											true	"Season number"
//	@Param			userId			query		int											true	"User ID"
//	@Success		200				{object}	responses.APIResponse[responses.MediaItemList[types.Episode]]	"Episodes retrieved successfully"
//	@Failure		400				{object}	responses.ErrorResponse[any]				"Invalid request"
//	@Failure		404				{object}	responses.ErrorResponse[any]				"Series or season not found"
//	@Failure		500				{object}	responses.ErrorResponse[any]				"Server error"
//
// Note: This functionality is implemented by GetEpisodesBySeriesIDAndSeasonNumber
func (h *coreSeriesHandler) GetSeasonWithEpisodes(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	seriesID, err := checkItemID(c, "id")
	if err != nil {
		return
	}

	seasonNumber, ok := checkSeasonNumber(c, "seasonNumber")
	if !ok {
		return
	}

	userID, _ := checkUserAccess(c)

	log.Debug().
		Uint64("seriesID", seriesID).
		Int("seasonNumber", seasonNumber).
		Uint64("userID", userID).
		Msg("Getting episodes for season")

	// Get the series first to ensure it exists
	series, err := h.seriesService.GetByID(ctx, seriesID)
	if handleServiceError(c, err, "Failed to retrieve series", "Series not found", "Series not found") {
		return
	}

	// Find the correct season and get its episodes
	var episodes []*types.Episode
	seasonFound := false

	for _, season := range series.Data.Seasons {
		if season.Number == seasonNumber {
			episodes = season.Episodes
			seasonFound = true
			break
		}
	}

	if !seasonFound {
		log.Warn().
			Uint64("seriesID", seriesID).
			Int("seasonNumber", seasonNumber).
			Msg("Season not found for series")
		responses.RespondNotFound(c, nil, "Season not found")
		return
	}

	if episodes == nil {
		episodes = []*types.Episode{}
	}
	// wrap all in media item list
	var episodeList []*models.MediaItem[*types.Episode]
	for _, episode := range episodes {
		episodeList = append(episodeList, models.NewMediaItem[*types.Episode](types.MediaTypeEpisode, episode))
	}

	log.Info().
		Uint64("seriesID", seriesID).
		Int("seasonNumber", seasonNumber).
		Int("episodeCount", len(episodes)).
		Msg("Episodes retrieved successfully")
	responses.RespondMediaItemListOK(c, episodeList, "Episodes retrieved successfully")
}

// GetByCreator godoc
//
//	@Summary		Get series by creator
//	@Description	Retrieves series created by a specific creator
//	@Tags			series, core
//	@Accept			json
//	@Produce		json
//	@Param			creatorId	path		int											true	"Creator ID"
//	@Param			limit		query		int											false	"Maximum number of series to return (default 10)"
//	@Param			offset		query		int											false	"Offset for pagination (default 0)"
//	@Success		200			{object}	responses.APIResponse[responses.MediaItemList[types.Series]]	"Series retrieved successfully"
//	@Failure		400			{object}	responses.ErrorResponse[any]				"Invalid request"
//	@Failure		500			{object}	responses.ErrorResponse[any]				"Server error"
//	@Router			/media/series/creator/{creatorId} [get]
func (h *coreSeriesHandler) GetByCreator(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	creatorID, err := checkItemID(c, "creatorId")
	if err != nil {
		return
	}

	limit := utils.GetLimit(c, 10, 100, true)
	offset := utils.GetOffset(c, 0)

	log.Debug().
		Uint64("creatorID", creatorID).
		Int("limit", limit).
		Msg("Getting series by creator")

	// Get all series for the user
	allSeries, err := h.seriesService.GetAll(ctx, limit, offset)
	if handleServiceError(c, err, "Failed to retrieve series", "", "Failed to retrieve series") {
		return
	}

	// Filter series by creator
	// This assumes the Series struct has a Creator field or similar
	var filteredSeries []*types.Series

	// for _, seriesItem := range allSeries {
	// if seriesItem.Data.Creator == creatorID {
	// filteredSeries = append(filteredSeries, seriesItem.Data)
	// if len(filteredSeries) >= limit {
	// break
	// }
	// }
	// }

	log.Info().
		Uint64("creatorID", creatorID).
		Int("count", len(filteredSeries)).
		Msg("Series by creator retrieved successfully")

	// If Series type doesn't have a Creator field, respond with empty result
	if len(filteredSeries) == 0 {
		log.Info().Msg("Creator-based filtering not fully implemented")
		var emptyList []*models.MediaItem[*types.Series]
		responses.RespondMediaItemListOK(c, emptyList, "Series retrieved successfully (creator field may not be available)")
		return
	}

	responses.RespondMediaItemListOK(c, allSeries, "Series retrieved successfully")
}

// GetEpisodesBySeasonID godoc
//
//	@Summary		Get episodes for a season
//	@Description	Retrieves all episodes for a specific season of a series
//	@Tags			series, core
//	@Accept			json
//	@Produce		json
//	@Param			id				path		int											true	"Series ID"
//	@Param			seasonNumber	path		int											true	"Season number"
//	@Param			userId			query		int											true	"User ID"
//	@Success		200				{object}	responses.APIResponse[responses.MediaItemList[types.Episode]]	"Episodes retrieved successfully"
//	@Failure		400				{object}	responses.ErrorResponse[any]				"Invalid request"
//	@Failure		404				{object}	responses.ErrorResponse[any]				"Series or season not found"
//	@Failure		500				{object}	responses.ErrorResponse[any]				"Server error"
//
// Note: This functionality is implemented by GetEpisodesBySeriesIDAndSeasonNumber
func (h *coreSeriesHandler) GetEpisodesBySeasonID(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	seriesID, err := checkItemID(c, "id")
	if err != nil {
		return
	}

	seasonNumber, ok := checkSeasonNumber(c, "seasonNumber")
	if !ok {
		return
	}

	userID, _ := checkUserAccess(c)

	log.Debug().
		Uint64("seriesID", seriesID).
		Int("seasonNumber", seasonNumber).
		Uint64("userID", userID).
		Msg("Getting episodes for season")

	// Get the series first to ensure it exists
	series, err := h.seriesService.GetByID(ctx, seriesID)
	if handleServiceError(c, err, "Failed to retrieve series", "Series not found", "Series not found") {
		return
	}

	// Find the correct season and get its episodes
	var episodes []*types.Episode
	seasonFound := false

	for _, season := range series.Data.Seasons {
		if season.Number == seasonNumber {
			episodes = season.Episodes
			seasonFound = true
			break
		}
	}

	if !seasonFound {
		log.Warn().
			Uint64("seriesID", seriesID).
			Int("seasonNumber", seasonNumber).
			Msg("Season not found for series")
		responses.RespondNotFound(c, nil, "Season not found")
		return
	}

	if episodes == nil {
		episodes = []*types.Episode{}
	}
	// wrap all in media item list
	var episodeList []*models.MediaItem[*types.Episode]
	for _, episode := range episodes {
		episodeList = append(episodeList, models.NewMediaItem[*types.Episode](types.MediaTypeEpisode, episode))
	}

	log.Info().
		Uint64("seriesID", seriesID).
		Int("seasonNumber", seasonNumber).
		Int("episodeCount", len(episodes)).
		Msg("Episodes retrieved successfully")
	responses.RespondMediaItemListOK(c, episodeList, "Episodes retrieved successfully")
}
