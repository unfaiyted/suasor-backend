// handlers/user_series.go
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

type UserSeriesHandler interface {
	CoreSeriesHandler

	GetFavoriteSeries(c *gin.Context)
	GetWatchedSeries(c *gin.Context)
	GetWatchlistSeries(c *gin.Context)
	GetRecentlyWatchedEpisodes(c *gin.Context)
	GetNextUpEpisodes(c *gin.Context)
	GetContinueWatchingSeries(c *gin.Context)
	UpdateSeriesUserData(c *gin.Context)
}

// userSeriesHandler handles operations for series items related to users
type userSeriesHandler struct {
	CoreSeriesHandler

	seriesItemService  services.UserMediaItemService[*mediatypes.Series]
	seasonItemService  services.UserMediaItemService[*mediatypes.Season]
	episodeItemService services.UserMediaItemService[*mediatypes.Episode]

	seriesDataService  services.UserMediaItemDataService[*mediatypes.Series]
	seasonDataService  services.UserMediaItemDataService[*mediatypes.Season]
	episodeDataService services.UserMediaItemDataService[*mediatypes.Episode]
}

// NewuserSeriesHandler creates a new user series handler
func NewUserSeriesHandler(
	coreHandler CoreSeriesHandler,

	// Items
	seriesService services.UserMediaItemService[*mediatypes.Series],
	seasonService services.UserMediaItemService[*mediatypes.Season],
	episodeService services.UserMediaItemService[*mediatypes.Episode],

	// Item Data
	seriesDataService services.UserMediaItemDataService[*mediatypes.Series],
	seasonDataService services.UserMediaItemDataService[*mediatypes.Season],
	episodeDataService services.UserMediaItemDataService[*mediatypes.Episode],

) UserSeriesHandler {
	return &userSeriesHandler{
		CoreSeriesHandler:  coreHandler,
		seriesItemService:  seriesService,
		seasonItemService:  seasonService,
		episodeItemService: episodeService,
		seriesDataService:  seriesDataService,
		seasonDataService:  seasonDataService,
		episodeDataService: episodeDataService,
	}
}

// GetFavoriteSeries godoc
// @Summary Get user favorite series
// @Description Retrieves series that a user has marked as favorites
// @Tags series
// @Accept json
// @Produce json
// @Security BearerAuth

// @Param limit query int false "Maximum number of series to return (default 10)"
// @Param offset query int false "Offset for pagination (default 0)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Series]] "Series retrieved successfully"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /user/series/favorites [get]
func (h *userSeriesHandler) GetFavoriteSeries(c *gin.Context) {
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
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		offset = 0
	}

	log.Debug().
		Uint64("userID", uid).
		Int("limit", limit).
		Msg("Getting favorite series")

	series, err := h.seriesDataService.GetFavorites(ctx, uid, limit, offset)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Msg("Failed to retrieve favorite series")
		responses.RespondInternalError(c, err, "Failed to retrieve favorite series")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", len(series)).
		Msg("Favorite series retrieved successfully")
	responses.RespondOK(c, series, "Favorite series retrieved successfully")
}

// GetWatchedSeries godoc
// @Summary Get user watched series
// @Description Retrieves series that a user has watched
// @Tags series
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Maximum number of series to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Series]] "Series retrieved successfully"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /user/series/watched [get]
func (h *userSeriesHandler) GetWatchedSeries(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access watched series without authentication")
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
		Msg("Getting watched series")

	// This is a placeholder for a real implementation
	// In a real implementation, you would query play history to find watched series
	// For now, we'll use SearchUserContent with a watched filter
	options := mediatypes.QueryOptions{
		MediaType: mediatypes.MediaTypeSeries,
		OwnerID:   uid,
		Watched:   true,
		Limit:     limit,
		Sort:      "last_watched",
		SortOrder: "desc",
	}

	series, err := h.seriesDataService.GetUserPlayHistory(ctx, options.OwnerID, &options)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Msg("Failed to retrieve watched series")
		responses.RespondInternalError(c, err, "Failed to retrieve watched series")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", len(series)).
		Msg("Watched series retrieved successfully")
	responses.RespondOK(c, series, "Watched series retrieved successfully")
}

// GetWatchlistSeries godoc
// @Summary Get series in user watchlist
// @Description Retrieves series that a user has added to their watchlist
// @Tags series
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Maximum number of series to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Series]] "Series retrieved successfully"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /user/series/watchlist [get]
func (h *userSeriesHandler) GetWatchlistSeries(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access watchlist without authentication")
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
		Msg("Getting watchlist series")

	// This is a placeholder for a real implementation
	// In a real implementation, you would query for series specifically marked for watchlist
	// For now, we'll use SearchUserContent with a watchlist filter
	options := mediatypes.QueryOptions{
		MediaType: mediatypes.MediaTypeSeries,
		OwnerID:   uid,
		Watchlist: true,
		Limit:     limit,
	}

	series, err := h.seriesDataService.Search(ctx, &options)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Msg("Failed to retrieve watchlist series")
		responses.RespondInternalError(c, err, "Failed to retrieve watchlist series")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", len(series)).
		Msg("Watchlist series retrieved successfully")
	responses.RespondOK(c, series, "Watchlist series retrieved successfully")
}

// UpdateSeriesUserData godoc
// @Summary Update user data for a series
// @Description Updates user-specific data for a series (favorite, watched status, rating, etc.)
// @Tags series
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Series ID"
// @Param data body models.UserMediaItemData[mediatypes.Series] true "Updated user data"
// @Success 200 {object} responses.APIResponse[models.MediaItem[mediatypes.Series]] "Series updated successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 404 {object} responses.ErrorResponse[responses.ErrorDetails] "Series not found"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /user/series/{id} [patch]
func (h *userSeriesHandler) UpdateSeriesUserData(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to update series data without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	seriesID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid series ID")
		responses.RespondBadRequest(c, err, "Invalid series ID")
		return
	}

	// Parse request body
	var userData models.UserMediaItemData[*mediatypes.Series]
	if err := c.ShouldBindJSON(&userData); err != nil {
		log.Warn().Err(err).Msg("Invalid request body")
		responses.RespondBadRequest(c, err, "Invalid request body")
		return
	}

	log.Debug().
		Uint64("userID", uid).
		Uint64("seriesID", seriesID).
		Interface("userData", userData).
		Msg("Updating series user data")

	// Get the existing series first
	series, err := h.seriesDataService.GetByID(ctx, seriesID)
	if err != nil {
		log.Error().Err(err).
			Uint64("seriesID", seriesID).
			Msg("Failed to retrieve series")
		responses.RespondNotFound(c, err, "Series not found")
		return
	}

	// Update user data
	// TODO: Proper add of user data to database here
	// In a real implementation, you would have a method like UpdateUserData
	// For now, we'll simulate by updating and using the regular Update method
	// series.UserData = userData

	// Update the series
	updatedSeries, err := h.seriesDataService.Update(ctx, series)
	if err != nil {
		log.Error().Err(err).
			Uint64("seriesID", seriesID).
			Msg("Failed to update series")
		responses.RespondInternalError(c, err, "Failed to update series")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("seriesID", seriesID).
		Msg("Series user data updated successfully")
	responses.RespondOK(c, updatedSeries, "Series updated successfully")
}

// GetContinueWatchingSeries godoc
// @Summary Get series in progress
// @Description Retrieves series that are currently in progress (partially watched)
// @Tags series
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId query int true "User ID"
// @Param limit query int false "Maximum number of series to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Series]] "Series retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /user/series/continue-watching [get]
func (h *userSeriesHandler) GetContinueWatchingSeries(c *gin.Context) {
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

	// This is a placeholder for a real implementation
	// In a real implementation, you would query for series specifically marked as favorites
	// For now, we'll just return a not implemented response since this requires integration with the play history service

	// TODO: Implement this by checking the play history for each series, finding the last watched episode,
	// and then determining the next episode in the sequence
	log.Info().Msg("Continue watching for series not yet implemented")
	responses.RespondNotImplemented(c, nil, "Continue watching for series not yet implemented")
}

// GetNextUpEpisodes godoc
// @Summary Get next episodes to watch
// @Description Retrieves the next unwatched episodes for series in progress
// @Tags series
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId query int true "User ID"
// @Param limit query int false "Maximum number of episodes to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Episode]] "Episodes retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /user/series/next-up [get]
func (h *userSeriesHandler) GetNextUpEpisodes(c *gin.Context) {
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
		Msg("Getting next up episodes")

	// This is a placeholder for a real implementation
	// In a real implementation, you would query play history to find series with episodes that have been partially watched
	// For now, we'll just return a not implemented response

	// TODO: Implement this by checking play history for each series, finding the last watched episode,
	// and then determining the next episode in the sequence
	log.Info().Msg("Next up episodes feature not yet implemented")
	responses.RespondNotImplemented(c, nil, "Next up episodes feature not yet implemented")
}

// GetRecentlyWatchedEpisodes godoc
// @Summary Get recently watched episodes
// @Description Retrieves the user's recently watched episodes
// @Tags series
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userId query int true "User ID"
// @Param days query int false "Number of days to look back (default 7)"
// @Param limit query int false "Maximum number of episodes to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Episode]] "Episodes retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /user/series/recently-watched [get]
func (h *userSeriesHandler) GetRecentlyWatchedEpisodes(c *gin.Context) {
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
	days, err := strconv.Atoi(c.DefaultQuery("days", "7"))
	if err != nil {
		days = 7
	}

	log.Debug().
		Uint64("userID", userID).
		Int("limit", limit).
		Msg("Getting recently watched series")

	// Get recently watched series
	series, err := h.episodeDataService.GetRecentHistory(ctx, userID, days, limit)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to retrieve recently watched series")
		responses.RespondInternalError(c, err, "Failed to retrieve series")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Int("count", len(series)).
		Msg("Recently watched series retrieved successfully")

	responses.RespondOK(c, series, "Recently watched series retrieved successfully")
}
