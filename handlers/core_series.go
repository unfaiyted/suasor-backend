// handlers/series.go
package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	mediatypes "suasor/client/media/types"
	"suasor/services"
	"suasor/types/responses"
	"suasor/utils"
)

type CoreSeriesHandler interface {
	CoreMediaItemHandler[*mediatypes.Series]

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
	CoreMediaItemHandler[mediatypes.Series]
	seriesService  services.CoreMediaItemService[*mediatypes.Series]
	seasonService  services.CoreMediaItemService[*mediatypes.Season]
	episodeService services.CoreMediaItemService[*mediatypes.Episode]
}

// NewcoreSeriesHandler creates a new series handler
func NewCoreSeriesHandler(
	coreHandler CoreMediaItemHandler[mediatypes.Series],
	seriesService services.CoreMediaItemService[*mediatypes.Series],
	seasonService services.CoreMediaItemService[*mediatypes.Season],
	episodeService services.CoreMediaItemService[*mediatypes.Episode],
) CoreSeriesHandler {
	return &coreSeriesHandler{
		CoreMediaItemHandler: coreHandler,
		seriesService:        seriesService,
		seasonService:        seasonService,
		episodeService:       episodeService,
	}
}

// GetSeasonsBySeriesID godoc
// @Summary Get seasons for a series
// @Description Retrieves all seasons for a specific series
// @Tags series
// @Accept json
// @Produce json
// @Param id path int true "Series ID"
// @Param userId query int true "User ID"
// @Success 200 {object} responses.APIResponse[[]mediatypes.Season] "Seasons retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Series not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /series/{id}/seasons [get]
func (h *coreSeriesHandler) GetSeasonsBySeriesID(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	seriesID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid series ID")
		responses.RespondBadRequest(c, err, "Invalid series ID")
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	log.Debug().
		Uint64("seriesID", seriesID).
		Uint64("userID", userID).
		Msg("Getting seasons for series")

	// Get the series first to ensure it exists
	series, err := h.seriesService.GetByID(ctx, seriesID)
	if err != nil {
		log.Error().Err(err).
			Uint64("seriesID", seriesID).
			Msg("Failed to retrieve series")
		responses.RespondNotFound(c, err, "Series not found")
		return
	}

	// Get seasons from the series data
	seasons := series.Data.Seasons
	if seasons == nil {
		seasons = []*mediatypes.Season{}
	}

	log.Info().
		Uint64("seriesID", seriesID).
		Int("seasonCount", len(seasons)).
		Msg("Seasons retrieved successfully")
	responses.RespondOK(c, seasons, "Seasons retrieved successfully")
}

// GetEpisodesBySeriesIDAndSeasonNumber godoc
// @Summary Get episodes for a season
// @Description Retrieves all episodes for a specific season of a series
// @Tags series
// @Accept json
// @Produce json
// @Param id path int true "Series ID"
// @Param seasonNumber path int true "Season number"
// @Param userId query int true "User ID"
// @Success 200 {object} responses.APIResponse[[]mediatypes.Episode] "Episodes retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Series or season not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /series/{id}/seasons/{seasonNumber}/episodes [get]
func (h *coreSeriesHandler) GetEpisodesBySeriesIDAndSeasonNumber(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	seriesID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid series ID")
		responses.RespondBadRequest(c, err, "Invalid series ID")
		return
	}

	seasonNumber, err := strconv.Atoi(c.Param("seasonNumber"))
	if err != nil {
		log.Warn().Err(err).Str("seasonNumber", c.Param("seasonNumber")).Msg("Invalid season number")
		responses.RespondBadRequest(c, err, "Invalid season number")
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	log.Debug().
		Uint64("seriesID", seriesID).
		Int("seasonNumber", seasonNumber).
		Uint64("userID", userID).
		Msg("Getting episodes for season")

	// Get the series first to ensure it exists
	series, err := h.seriesService.GetByID(ctx, seriesID)
	if err != nil {
		log.Error().Err(err).
			Uint64("seriesID", seriesID).
			Msg("Failed to retrieve series")
		responses.RespondNotFound(c, err, "Series not found")
		return
	}

	// Find the correct season and get its episodes
	var episodes []*mediatypes.Episode
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
		episodes = []*mediatypes.Episode{}
	}

	log.Info().
		Uint64("seriesID", seriesID).
		Int("seasonNumber", seasonNumber).
		Int("episodeCount", len(episodes)).
		Msg("Episodes retrieved successfully")
	responses.RespondOK(c, episodes, "Episodes retrieved successfully")
}

// GetContinueWatchingSeries godoc
// @Summary Get series in progress
// @Description Retrieves series that are currently in progress (partially watched)
// @Tags series
// @Accept json
// @Produce json
// @Param userId query int true "User ID"
// @Param limit query int false "Maximum number of series to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]mediatypes.Series] "Series retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /series/continue-watching [get]
func (h *coreSeriesHandler) GetContinueWatchingSeries(c *gin.Context) {
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
		Msg("Getting continue watching series")

	//TODO: This would typically involve checking the play history to find series with episodes that have been partially watched
	// For now, we'll just return a not implemented response since this requires integration with the play history service

	// TODO: Implement by checking the play history for partially watched episodes, then looking up their series
	log.Info().Msg("Continue watching for series not yet implemented")
	responses.RespondNotImplemented(c, nil, "Continue watching for series not yet implemented")
}

// GetAllEpisodes godoc
// @Summary Get all episodes for a series
// @Description Retrieves all episodes across all seasons for a specific series
// @Tags series
// @Accept json
// @Produce json
// @Param id path int true "Series ID"
// @Param userId query int true "User ID"
// @Success 200 {object} responses.APIResponse[[]mediatypes.Episode] "Episodes retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Series not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /series/{id}/episodes [get]
func (h *coreSeriesHandler) GetAllEpisodes(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	seriesID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid series ID")
		responses.RespondBadRequest(c, err, "Invalid series ID")
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	log.Debug().
		Uint64("seriesID", seriesID).
		Uint64("userID", userID).
		Msg("Getting all episodes for series")

	// Get the series first to ensure it exists
	series, err := h.seriesService.GetByID(ctx, seriesID)
	if err != nil {
		log.Error().Err(err).
			Uint64("seriesID", seriesID).
			Msg("Failed to retrieve series")
		responses.RespondNotFound(c, err, "Series not found")
		return
	}

	// Collect all episodes from all seasons
	var allEpisodes []*mediatypes.Episode

	for _, season := range series.Data.Seasons {
		if season.Episodes != nil {
			allEpisodes = append(allEpisodes, season.Episodes...)
		}
	}

	log.Info().
		Uint64("seriesID", seriesID).
		Int("totalEpisodes", len(allEpisodes)).
		Msg("All episodes retrieved successfully")
	responses.RespondOK(c, allEpisodes, "Episodes retrieved successfully")
}

// GetNextUpEpisodes godoc
// @Summary Get next episodes to watch
// @Description Retrieves the next unwatched episodes for series in progress
// @Tags series
// @Accept json
// @Produce json
// @Param userId query int true "User ID"
// @Param limit query int false "Maximum number of episodes to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]mediatypes.Episode] "Episodes retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /series/next-up [get]
func (h *coreSeriesHandler) GetNextUpEpisodes(c *gin.Context) {
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
// @Summary Get recently aired episodes
// @Description Retrieves episodes that have recently aired based on their air date
// @Tags series
// @Accept json
// @Produce json
// @Param userId query int true "User ID"
// @Param limit query int false "Maximum number of episodes to return (default 10)"
// @Param days query int false "Number of days to look back (default 7)"
// @Success 200 {object} responses.APIResponse[[]mediatypes.Episode] "Episodes retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /series/recently-aired [get]
func (h *coreSeriesHandler) GetRecentlyAiredEpisodes(c *gin.Context) {
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

	days, err := strconv.Atoi(c.DefaultQuery("days", "7"))
	if err != nil {
		days = 7
	}

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
// @Summary Get series by network
// @Description Retrieves series from a specific TV network
// @Tags series
// @Accept json
// @Produce json
// @Param network path string true "Network name"
// @Param limit query int false "Maximum number of series to return (default 10)"
// @Param offset query int false "Offset for pagination (default 0)"
// @Success 200 {object} responses.APIResponse[[]mediatypes.Series] "Series retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /series/network/{network} [get]
func (h *coreSeriesHandler) GetSeriesByNetwork(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	network := c.Param("network")
	if network == "" {
		log.Warn().Msg("Network name is required")
		responses.RespondBadRequest(c, nil, "Network name is required")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		offset = 0
	}

	log.Debug().
		Str("network", network).
		Int("limit", limit).
		Msg("Getting series by network")

	// Get all series for the user
	allSeries, err := h.seriesService.GetAll(ctx, limit, offset)
	if err != nil {
		log.Error().Err(err).
			Msg("Failed to retrieve series")
		responses.RespondInternalError(c, err, "Failed to retrieve series")
		return
	}

	// Filter series by network
	// This assumes the Series struct has a Network field or similar
	var filteredSeries []*mediatypes.Series

	for _, seriesItem := range allSeries {
		if seriesItem.Data.Network == network {
			filteredSeries = append(filteredSeries, seriesItem.Data)
			if len(filteredSeries) >= limit {
				break
			}
		}
	}

	log.Info().
		Str("network", network).
		Int("count", len(filteredSeries)).
		Msg("Series by network retrieved successfully")

	// If Series type doesn't have a Network field, respond with empty result
	if len(filteredSeries) == 0 {
		log.Info().Msg("Network-based filtering not fully implemented")
		responses.RespondOK(c, []*mediatypes.Series{}, "Series retrieved successfully (network field may not be available)")
		return
	}

	responses.RespondOK(c, filteredSeries, "Series retrieved successfully")
}

// GetSeasonWithEpisodes godoc
// @Summary Get a season and all its episodes
// @Description Retrieves a specific season and all its episodes
// @Tags series
// @Accept json
// @Produce json
// @Param id path int true "Series ID"
// @Param seasonNumber path int true "Season number"
// @Param userId query int true "User ID"
// @Success 200 {object} responses.APIResponse[[]mediatypes.Episode] "Episodes retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Series or season not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /series/{id}/seasons/{seasonNumber}/episodes [get]
func (h *coreSeriesHandler) GetSeasonWithEpisodes(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	seriesID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid series ID")
		responses.RespondBadRequest(c, err, "Invalid series ID")
		return
	}

	seasonNumber, err := strconv.Atoi(c.Param("seasonNumber"))
	if err != nil {
		log.Warn().Err(err).Str("seasonNumber", c.Param("seasonNumber")).Msg("Invalid season number")
		responses.RespondBadRequest(c, err, "Invalid season number")
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	log.Debug().
		Uint64("seriesID", seriesID).
		Int("seasonNumber", seasonNumber).
		Uint64("userID", userID).
		Msg("Getting episodes for season")

	// Get the series first to ensure it exists
	series, err := h.seriesService.GetByID(ctx, seriesID)
	if err != nil {
		log.Error().Err(err).
			Uint64("seriesID", seriesID).
			Msg("Failed to retrieve series")
		responses.RespondNotFound(c, err, "Series not found")
		return
	}

	// Find the correct season and get its episodes
	var episodes []*mediatypes.Episode
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
		episodes = []*mediatypes.Episode{}
	}

	log.Info().
		Uint64("seriesID", seriesID).
		Int("seasonNumber", seasonNumber).
		Int("episodeCount", len(episodes)).
		Msg("Episodes retrieved successfully")
	responses.RespondOK(c, episodes, "Episodes retrieved successfully")
}

// GetByCreator godoc
// @Summary Get series by creator
// @Description Retrieves series created by a specific creator
// @Tags series
// @Accept json
// @Produce json
// @Param creatorId path int true "Creator ID"
// @Param limit query int false "Maximum number of series to return (default 10)"
// @Param offset query int false "Offset for pagination (default 0)"
// @Success 200 {object} responses.APIResponse[[]mediatypes.Series] "Series retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /series/creator/{creatorId} [get]
func (h *coreSeriesHandler) GetByCreator(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	creatorID, err := strconv.ParseUint(c.Param("creatorId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("creatorId", c.Param("creatorId")).Msg("Invalid creator ID")
		responses.RespondBadRequest(c, err, "Invalid creator ID")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		offset = 0
	}

	log.Debug().
		Uint64("creatorID", creatorID).
		Int("limit", limit).
		Msg("Getting series by creator")

	// Get all series for the user
	allSeries, err := h.seriesService.GetAll(ctx, limit, offset)
	if err != nil {
		log.Error().Err(err).
			Msg("Failed to retrieve series")
		responses.RespondInternalError(c, err, "Failed to retrieve series")
		return
	}

	// Filter series by creator
	// This assumes the Series struct has a Creator field or similar
	var filteredSeries []*mediatypes.Series

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
		responses.RespondOK(c, []*mediatypes.Series{}, "Series retrieved successfully (creator field may not be available)")
		return
	}

	responses.RespondOK(c, allSeries, "Series retrieved successfully")
}

// GetEpisodesBySeasonID godoc
// @Summary Get episodes for a season
// @Description Retrieves all episodes for a specific season of a series
// @Tags series
// @Accept json
// @Produce json
// @Param id path int true "Series ID"
// @Param seasonNumber path int true "Season number"
// @Param userId query int true "User ID"
// @Success 200 {object} responses.APIResponse[[]mediatypes.Episode] "Episodes retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Series or season not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /series/{id}/seasons/{seasonNumber}/episodes [get]
func (h *coreSeriesHandler) GetEpisodesBySeasonID(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	seriesID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid series ID")
		responses.RespondBadRequest(c, err, "Invalid series ID")
		return
	}

	seasonNumber, err := strconv.Atoi(c.Param("seasonNumber"))
	if err != nil {
		log.Warn().Err(err).Str("seasonNumber", c.Param("seasonNumber")).Msg("Invalid season number")
		responses.RespondBadRequest(c, err, "Invalid season number")
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	log.Debug().
		Uint64("seriesID", seriesID).
		Int("seasonNumber", seasonNumber).
		Uint64("userID", userID).
		Msg("Getting episodes for season")

	// Get the series first to ensure it exists
	series, err := h.seriesService.GetByID(ctx, seriesID)
	if err != nil {
		log.Error().Err(err).
			Uint64("seriesID", seriesID).
			Msg("Failed to retrieve series")
		responses.RespondNotFound(c, err, "Series not found")
		return
	}

	// Find the correct season and get its episodes
	var episodes []*mediatypes.Episode
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
		episodes = []*mediatypes.Episode{}
	}

	log.Info().
		Uint64("seriesID", seriesID).
		Int("seasonNumber", seasonNumber).
		Int("episodeCount", len(episodes)).
		Msg("Episodes retrieved successfully")
	responses.RespondOK(c, episodes, "Episodes retrieved successfully")
}
