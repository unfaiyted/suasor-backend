package handlers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"net/http"

	"suasor/client/media/types"
	"suasor/services"
	"suasor/types/models"
	"suasor/types/responses"
	"suasor/utils"
)

// UserUserMediaItemDataHandler handles user-specific operations for user media item data
// This is the user layer of the three-pronged architecture
type UserMediaItemDataHandler[T types.MediaData] struct {
	service     services.UserMediaItemDataService[T]
	coreHandler *CoreUserMediaItemDataHandler[T]
}

// NewUserUserMediaItemDataHandler creates a new user user media item data handler
func NewUserUserMediaItemDataHandler[T types.MediaData](
	service services.UserMediaItemDataService[T],
	coreHandler *CoreUserMediaItemDataHandler[T],
) *UserMediaItemDataHandler[T] {
	return &UserMediaItemDataHandler[T]{
		service:     service,
		coreHandler: coreHandler,
	}
}

// GetMediaPlayHistory godoc
// @Summary Get a user's media play history
// @Description Get a user's media play history with optional filtering
// @Tags History
// @Accept json
// @Produce json
// @Param userId query int true "User ID"
// @Param limit query int false "Number of items to return (default 10)"
// @Param offset query int false "Number of items to skip (default 0)"
// @Param type query string false "Media type filter (movie, series, episode, track, etc.)"
// @Param completed query bool false "Filter by completion status"
// @Success 200 {object} responses.APIResponse[[]models.UserMediaItemData[any]] "Successfully retrieved play history"
// @Failure 400 {object} responses.ErrorResponse[any] "Bad request"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[any] "Internal server error"
// @Router /user-media-data/history [get]
func (h *UserMediaItemDataHandler[T]) GetMediaPlayHistory(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	mediaType := c.Query("type")
	completedStr := c.Query("completed")

	var completed *bool
	if completedStr != "" {
		completedBool, err := strconv.ParseBool(completedStr)
		if err != nil {
			log.Warn().Err(err).Str("completed", completedStr).Msg("Invalid completed value")
			responses.RespondBadRequest(c, err, "Invalid completed value")
			return
		}
		completed = &completedBool
	}

	// Filter by media type if provided
	var typedMediaType *types.MediaType
	if mediaType != "" {
		mt := types.MediaType(mediaType)
		typedMediaType = &mt
	}

	log.Debug().
		Uint64("userId", userID).
		Int("limit", limit).
		Int("offset", offset).
		Str("mediaType", mediaType).
		Msg("Getting user media play history")

	history, err := h.service.GetUserPlayHistory(ctx, userID, limit, offset, typedMediaType, completed)
	if err != nil {
		log.Error().Err(err).Uint64("userId", userID).Msg("Failed to retrieve play history")
		responses.RespondInternalError(c, err, "Failed to retrieve play history")
		return
	}

	log.Info().
		Uint64("userId", userID).
		Int("count", len(history)).
		Msg("Play history retrieved successfully")

	responses.RespondOK(c, history, "Play history retrieved successfully")
}

// GetContinuePlaying godoc
// @Summary Get a user's continue watching list
// @Description Get media items that a user has started but not completed
// @Tags History
// @Accept json
// @Produce json
// @Param userId query int true "User ID"
// @Param limit query int false "Number of items to return (default 10)"
// @Param type query string false "Media type filter (movie, series, episode, track, etc.)"
// @Success 200 {object} responses.APIResponse[[]models.UserMediaItemData[any]] "Successfully retrieved continue watching items"
// @Failure 400 {object} responses.ErrorResponse[any] "Bad request"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[any] "Internal server error"
// @Router /user-media-data/continue-watching [get]
func (h *UserMediaItemDataHandler[T]) GetContinuePlaying(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	mediaType := c.Query("type")

	// Filter by media type if provided
	var typedMediaType *types.MediaType
	if mediaType != "" {
		mt := types.MediaType(mediaType)
		typedMediaType = &mt
	}

	log.Debug().
		Uint64("userId", userID).
		Int("limit", limit).
		Str("mediaType", mediaType).
		Msg("Getting continue watching items")

	// Get items that are not completed and have been played recently
	items, err := h.service.GetContinueWatching(ctx, userID, limit, typedMediaType)
	if err != nil {
		log.Error().Err(err).Uint64("userId", userID).Msg("Failed to retrieve continue watching items")
		responses.RespondInternalError(c, err, "Failed to retrieve continue watching items")
		return
	}

	log.Info().
		Uint64("userId", userID).
		Int("count", len(items)).
		Msg("Continue watching items retrieved successfully")

	responses.RespondOK(c, items, "Continue watching items retrieved successfully")
}

// GetRecentHistory godoc
// @Summary Get a user's recent media history
// @Description Get a user's recent media history
// @Tags History
// @Accept json
// @Produce json
// @Param userId query int true "User ID"
// @Param limit query int false "Number of items to return (default 10)"
// @Param type query string false "Media type filter (movie, series, episode, track, etc.)"
// @Success 200 {object} responses.APIResponse[[]models.UserMediaItemData[any]] "Successfully retrieved recent history"
// @Failure 400 {object} responses.ErrorResponse[any] "Bad request"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[any] "Internal server error"
// @Router /user-media-data/recent [get]
func (h *UserMediaItemDataHandler[T]) GetRecentHistory(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	mediaType := c.Query("type")

	// Filter by media type if provided
	var typedMediaType *types.MediaType
	if mediaType != "" {
		mt := types.MediaType(mediaType)
		typedMediaType = &mt
	}

	log.Debug().
		Uint64("userId", userID).
		Int("limit", limit).
		Str("mediaType", mediaType).
		Msg("Getting recent user media history")

	history, err := h.service.GetRecentHistory(ctx, userID, limit, typedMediaType)
	if err != nil {
		log.Error().Err(err).Uint64("userId", userID).Msg("Failed to retrieve recent history")
		responses.RespondInternalError(c, err, "Failed to retrieve recent history")
		return
	}

	log.Info().
		Uint64("userId", userID).
		Int("count", len(history)).
		Msg("Recent history retrieved successfully")

	responses.RespondOK(c, history, "Recent history retrieved successfully")
}

// RecordMediaPlay godoc
// @Summary Record a media play event
// @Description Record a new play event for a media item
// @Tags History
// @Accept json
// @Produce json
// @Param mediaPlay body models.UserMediaItemDataRequest true "Media play information"
// @Success 201 {object} responses.APIResponse[models.UserMediaItemData[any]] "Play event recorded successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Bad request"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[any] "Internal server error"
// @Router /user-media-data/record [post]
func (h *UserMediaItemDataHandler[T]) RecordMediaPlay(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	var req models.UserMediaItemDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid request body")
		responses.RespondBadRequest(c, err, "Invalid request body")
		return
	}

	log.Debug().
		Uint64("userId", req.UserID).
		Uint64("mediaItemId", req.MediaItemID).
		Str("type", string(req.Type)).
		Msg("Recording media play event")

	// Create a play history record
	playHistory := &models.UserMediaItemData[T]{
		UserID:           req.UserID,
		MediaItemID:      req.MediaItemID,
		Type:             req.Type,
		PlayedAt:         time.Now(),
		LastPlayedAt:     time.Now(),
		IsFavorite:       req.IsFavorite,
		UserRating:       req.UserRating,
		PlayedPercentage: req.PlayedPercentage,
		PositionSeconds:  req.PositionSeconds,
		DurationSeconds:  req.DurationSeconds,
		Completed:        req.Completed,
	}

	// If this is a continuation, increment the play count
	if req.Continued {
		playHistory.PlayCount = 1
	}

	result, err := h.service.RecordPlay(ctx, playHistory)
	if err != nil {
		log.Error().Err(err).
			Uint64("userId", req.UserID).
			Uint64("mediaItemId", req.MediaItemID).
			Msg("Failed to record play event")
		responses.RespondInternalError(c, err, "Failed to record play event")
		return
	}

	log.Info().
		Uint64("userId", req.UserID).
		Uint64("mediaItemId", req.MediaItemID).
		Msg("Play event recorded successfully")

	responses.RespondCreated(c, result, "Play event recorded successfully")
}

// ToggleFavorite godoc
// @Summary Toggle favorite status for a media item
// @Description Mark or unmark a media item as a favorite
// @Tags History
// @Accept json
// @Produce json
// @Param mediaItemId path int true "Media Item ID"
// @Param userId query int true "User ID"
// @Param favorite query bool true "Favorite status"
// @Success 200 {object} responses.APIResponse[models.UserMediaItemData[any]] "Favorite status updated successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Bad request"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[any] "Internal server error"
// @Router /user-media-data/media/{mediaItemId}/favorite [put]
func (h *UserMediaItemDataHandler[T]) ToggleFavorite(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	mediaItemID, err := strconv.ParseUint(c.Param("mediaItemId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("mediaItemId", c.Param("mediaItemId")).Msg("Invalid media item ID")
		responses.RespondBadRequest(c, err, "Invalid media item ID")
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	favorite, err := strconv.ParseBool(c.Query("favorite"))
	if err != nil {
		log.Warn().Err(err).Str("favorite", c.Query("favorite")).Msg("Invalid favorite value")
		responses.RespondBadRequest(c, err, "Invalid favorite value")
		return
	}

	log.Debug().
		Uint64("userId", userID).
		Uint64("mediaItemId", mediaItemID).
		Bool("favorite", favorite).
		Msg("Toggling favorite status")

	err = h.service.ToggleFavorite(ctx, mediaItemID, userID, favorite)
	if err != nil {
		log.Error().Err(err).
			Uint64("userId", userID).
			Uint64("mediaItemId", mediaItemID).
			Msg("Failed to update favorite status")
		responses.RespondInternalError(c, err, "Failed to update favorite status")
		return
	}

	log.Info().
		Uint64("userId", userID).
		Uint64("mediaItemId", mediaItemID).
		Bool("favorite", favorite).
		Msg("Favorite status updated successfully")

	responses.RespondOK(c, http.StatusOK, "Favorite status updated successfully")
}

// UpdateUserRating godoc
// @Summary Update user rating for a media item
// @Description Set a user's rating for a media item
// @Tags History
// @Accept json
// @Produce json
// @Param mediaItemId path int true "Media Item ID"
// @Param userId query int true "User ID"
// @Param rating query number true "User rating (0-10)"
// @Success 200 {object} responses.APIResponse[models.UserMediaItemData[any]] "Rating updated successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Bad request"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[any] "Internal server error"
// @Router /user-media-data/media/{mediaItemId}/rating [put]
func (h *UserMediaItemDataHandler[T]) UpdateUserRating(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	mediaItemID, err := strconv.ParseUint(c.Param("mediaItemId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("mediaItemId", c.Param("mediaItemId")).Msg("Invalid media item ID")
		responses.RespondBadRequest(c, err, "Invalid media item ID")
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	rating, err := strconv.ParseFloat(c.Query("rating"), 32)
	if err != nil {
		log.Warn().Err(err).Str("rating", c.Query("rating")).Msg("Invalid rating value")
		responses.RespondBadRequest(c, err, "Invalid rating value")
		return
	}

	if rating < 0 || rating > 10 {
		log.Warn().Float64("rating", rating).Msg("Rating must be between 0 and 10")
		responses.RespondBadRequest(c, nil, "Rating must be between 0 and 10")
		return
	}

	log.Debug().
		Uint64("userId", userID).
		Uint64("mediaItemId", mediaItemID).
		Float64("rating", rating).
		Msg("Updating user rating")

	err = h.service.UpdateRating(ctx, mediaItemID, userID, float32(rating))
	if err != nil {
		log.Error().Err(err).
			Uint64("userId", userID).
			Uint64("mediaItemId", mediaItemID).
			Msg("Failed to update rating")
		responses.RespondInternalError(c, err, "Failed to update rating")
		return
	}

	log.Info().
		Uint64("userId", userID).
		Uint64("mediaItemId", mediaItemID).
		Float64("rating", rating).
		Msg("Rating updated successfully")

	responses.RespondOK(c, http.StatusOK, "Rating updated successfully")
}

// GetFavorites godoc
// @Summary Get a user's favorite media items
// @Description Get all media items marked as favorites by a user
// @Tags History
// @Accept json
// @Produce json
// @Param userId query int true "User ID"
// @Param limit query int false "Number of items to return (default 10)"
// @Param offset query int false "Number of items to skip (default 0)"
// @Success 200 {object} responses.APIResponse[[]models.UserMediaItemData[any]] "Successfully retrieved favorites"
// @Failure 400 {object} responses.ErrorResponse[any] "Bad request"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[any] "Internal server error"
// @Router /user-media-data/favorites [get]
func (h *UserMediaItemDataHandler[T]) GetFavorites(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	log.Debug().
		Uint64("userId", userID).
		Int("limit", limit).
		Int("offset", offset).
		Msg("Getting user favorites")

	favorites, err := h.service.GetFavorites(ctx, userID, limit, offset)
	if err != nil {
		log.Error().Err(err).Uint64("userId", userID).Msg("Failed to retrieve favorites")
		responses.RespondInternalError(c, err, "Failed to retrieve favorites")
		return
	}

	log.Info().
		Uint64("userId", userID).
		Int("count", len(favorites)).
		Msg("Favorites retrieved successfully")

	responses.RespondOK(c, favorites, "Favorites retrieved successfully")
}

// ClearUserHistory godoc
// @Summary Clear a user's play history
// @Description Delete all play history entries for a user
// @Tags History
// @Accept json
// @Produce json
// @Param userId query int true "User ID"
// @Success 200 {object} responses.APIResponse[any] "History cleared successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Bad request"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[any] "Internal server error"
// @Router /user-media-data/clear [delete]
func (h *UserMediaItemDataHandler[T]) ClearUserHistory(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	log.Debug().
		Uint64("userId", userID).
		Msg("Clearing user history")

	err = h.service.ClearUserHistory(ctx, userID)
	if err != nil {
		log.Error().Err(err).Uint64("userId", userID).Msg("Failed to clear history")
		responses.RespondInternalError(c, err, "Failed to clear history")
		return
	}

	log.Info().
		Uint64("userId", userID).
		Msg("History cleared successfully")

	responses.RespondOK(c, responses.EmptyResponse{Success: true}, "History cleared successfully")
}

// Core handler method forwarding

func (h *UserMediaItemDataHandler[T]) GetMediaItemDataByID(c *gin.Context) {
	h.coreHandler.GetMediaItemDataByID(c)
}

func (h *UserMediaItemDataHandler[T]) CheckUserMediaItemData(c *gin.Context) {
	h.coreHandler.CheckUserMediaItemData(c)
}

func (h *UserMediaItemDataHandler[T]) GetMediaItemDataByUserAndMedia(c *gin.Context) {
	h.coreHandler.GetMediaItemDataByUserAndMedia(c)
}

func (h *UserMediaItemDataHandler[T]) DeleteMediaItemData(c *gin.Context) {
	h.coreHandler.DeleteMediaItemData(c)
}
