package handlers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"net/http"

	mediatypes "suasor/clients/media/types"
	"suasor/services"
	models "suasor/types/models"
	"suasor/types/requests"
	"suasor/types/responses"
	"suasor/utils"
	"suasor/utils/logger"
)

type UserMediaItemDataHandler[T mediatypes.MediaData] interface {
	CoreUserMediaItemDataHandler[T]

	GetMediaPlayHistory(c *gin.Context)
	GetContinuePlaying(c *gin.Context)
	GetRecentHistory(c *gin.Context)
	RecordMediaPlay(c *gin.Context)
	ToggleFavorite(c *gin.Context)
	UpdateUserRating(c *gin.Context)
	GetFavorites(c *gin.Context)
	ClearUserHistory(c *gin.Context)
}

// UseruserMediaItemDataHandler handles user-specific operations for user media item data
// This is the user layer of the three-pronged architecture
type userMediaItemDataHandler[T mediatypes.MediaData] struct {
	CoreUserMediaItemDataHandler[T]
	service services.UserMediaItemDataService[T]
}

// NewUseruserMediaItemDataHandler creates a new user user media item data handler
func NewUserMediaItemDataHandler[T mediatypes.MediaData](
	coreHandler CoreUserMediaItemDataHandler[T],
	service services.UserMediaItemDataService[T],
) UserMediaItemDataHandler[T] {
	return &userMediaItemDataHandler[T]{
		CoreUserMediaItemDataHandler: coreHandler,
		service:                      service,
	}
}

// GetMediaPlayHistory godoc
//
//	@Summary		Get a user's media play user-data
//	@Description	Get a user's media play user-data with optional filtering
//	@Tags			user-data
//	@Accept			json
//	@Produce		json
//	@Param			userID		query		int																	false	"User ID (optional, uses authenticated user ID if not provided)"
//	@Param			limit		query		int																	false	"Number of items to return (default 10)"
//	@Param			offset		query		int																	false	"Number of items to skip (default 0)"
//	@Param			completed	query		bool																false	"Filter by completion status"
//	@Success		200			{object}	responses.APIResponse[[]models.UserMediaItemData[mediatypes.MediaData]]	"Successfully retrieved play user-data"
//	@Failure		400			{object}	responses.ErrorResponse[responses.ErrorDetails]						"Bad request"
//	@Failure		401			{object}	responses.ErrorResponse[responses.ErrorDetails]						"Unauthorized"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]						"Internal server error"
//	@Router			/user-data/{mediaType}/history [get]
func (h *userMediaItemDataHandler[T]) GetMediaPlayHistory(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	userID, err := utils.GetUserID(c)
	if err != nil {
		log.Warn().Err(err).Str("userID", c.Query("userID")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	limit := utils.GetLimit(c, 10, 100, true)
	offset := utils.GetOffset(c, 0)
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

	var zero T
	mediaType := mediatypes.GetMediaTypeFromTypeName(zero)

	log.Debug().
		Uint64("userID", userID).
		Int("limit", limit).
		Int("offset", offset).
		Str("mediaType", string(mediaType)).
		Msg("Getting user media play user-data")

	query := mediatypes.QueryOptions{
		MediaType: mediaType,
		Limit:     limit,
		Offset:    offset,
		Watched:   *completed,
	}

	history, err := h.service.GetUserPlayHistory(ctx, userID, &query)
	if err != nil {
		log.Error().Err(err).Uint64("userID", userID).Msg("Failed to retrieve play user-data")
		responses.RespondInternalError(c, err, "Failed to retrieve play user-data")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Int("count", len(history)).
		Msg("Play user-data retrieved successfully")

	responses.RespondOK(c, history, "Play history retrieved successfully")
}

// GetContinuePlaying godoc
//
//	@Summary		Get a user's continue watching list
//	@Description	Get media items that a user has started but not completed
//	@Tags			user-data
//	@Accept			json
//	@Produce		json
//	@Param			userID	query		int																	false	"User ID (optional, uses authenticated user ID if not provided)"
//	@Param			limit	query		int																	false	"Number of items to return (default 10)"
//	@Param			type	query		string																false	"Media type filter (movie, series, episode, track, etc.)"
//	@Success		200		{object}	responses.APIResponse[[]models.UserMediaItemData[mediatypes.MediaData]]	"Successfully retrieved continue watching items"
//	@Failure		400		{object}	responses.ErrorResponse[responses.ErrorDetails]						"Bad request"
//	@Failure		401		{object}	responses.ErrorResponse[responses.ErrorDetails]						"Unauthorized"
//	@Failure		500		{object}	responses.ErrorResponse[responses.ErrorDetails]						"Internal server error"
//	@Router			/user-data/{mediaType}/continue-watching [get]
func (h *userMediaItemDataHandler[T]) GetContinuePlaying(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	userID, err := utils.GetUserID(c)
	if err != nil {
		log.Warn().Err(err).Str("userID", c.Query("userID")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	limit := utils.GetLimit(c, 10, 100, true)

	log.Debug().
		Uint64("userID", userID).
		Int("limit", limit).
		Msg("Getting continue watching items")

	// Get items that are not completed and have been played recently
	items, err := h.service.GetContinueWatching(ctx, userID, limit)
	if err != nil {
		log.Error().Err(err).Uint64("userID", userID).Msg("Failed to retrieve continue watching items")
		responses.RespondInternalError(c, err, "Failed to retrieve continue watching items")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Int("count", len(items)).
		Msg("Continue watching items retrieved successfully")

	responses.RespondOK(c, items, "Continue watching items retrieved successfully")
}

// GetRecentHistory godoc
//
//	@Summary		Get a user's recent media user-data
//	@Description	Get a user's recent media user-data
//	@Tags			user-data
//	@Accept			json
//	@Produce		json
//	@Param			userID	query		int																	false	"User ID (optional, uses authenticated user ID if not provided)"
//	@Param			limit	query		int																	false	"Number of items to return (default 10)"
//	@Param			days	query		int																	false	"Number of days to look back (default 7)"
//	@Success		200		{object}	responses.APIResponse[[]models.UserMediaItemData[mediatypes.MediaData]]	"Successfully retrieved recent user-data"
//	@Failure		400		{object}	responses.ErrorResponse[responses.ErrorDetails]						"Bad request"
//	@Failure		401		{object}	responses.ErrorResponse[responses.ErrorDetails]						"Unauthorized"
//	@Failure		500		{object}	responses.ErrorResponse[responses.ErrorDetails]						"Internal server error"
//	@Router			/user-data/{mediaType}/recent [get]
func (h *userMediaItemDataHandler[T]) GetRecentHistory(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	userID, err := utils.GetUserID(c)
	if err != nil {
		log.Warn().Err(err).Str("userID", c.Query("userID")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	limit := utils.GetLimit(c, 10, 100, true)
	days, err := strconv.Atoi(c.DefaultQuery("days", "7"))
	if err != nil {
		days = 7
	}

	log.Debug().
		Uint64("userID", userID).
		Int("limit", limit).
		Msg("Getting recent user media user-data")

	history, err := h.service.GetRecentHistory(ctx, userID, days, limit)
	if err != nil {
		log.Error().Err(err).Uint64("userID", userID).Msg("Failed to retrieve recent user-data")
		responses.RespondInternalError(c, err, "Failed to retrieve recent user-data")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Int("count", len(history)).
		Msg("Recent user-data retrieved successfully")

	responses.RespondOK(c, history, "Recent history retrieved successfully")
}

// RecordMediaPlay godoc
//
//	@Summary		Record a media play event
//	@Description	Record a new play event for a media item
//	@Tags			user-data
//	@Accept			json
//	@Produce		json
//	@Param			mediaPlay	body		requests.UserMediaItemDataRequest									true	"Media play information"
//	@Param			itemID		path		int																	true	"Media Item ID"
//	@Success		201			{object}	responses.APIResponse[models.UserMediaItemData[mediatypes.MediaData]]	"Play event recorded successfully"
//	@Failure		400			{object}	responses.ErrorResponse[responses.ErrorDetails]						"Bad request"
//	@Failure		401			{object}	responses.ErrorResponse[responses.ErrorDetails]						"Unauthorized"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]						"Internal server error"
//	@Router			/user-data/{mediaType}/{itemID}/record [post]
func (h *userMediaItemDataHandler[T]) RecordMediaPlay(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	var req requests.UserMediaItemDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid request body")
		responses.RespondBadRequest(c, err, "Invalid request body")
		return
	}

	// If UserID is not provided in the request, use the authenticated user
	if req.UserID == 0 {
		userID, err := utils.GetUserID(c)
		if err != nil {
			log.Warn().Err(err).Msg("Could not determine user ID")
			responses.RespondBadRequest(c, err, "User ID is required")
			return
		}
		req.UserID = userID
	}

	log.Debug().
		Uint64("userID", req.UserID).
		Uint64("itemID", req.MediaItemID).
		Str("type", string(req.Type)).
		Msg("Recording media play event")

	// Create a play user-data record
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
			Uint64("userID", req.UserID).
			Uint64("itemID", req.MediaItemID).
			Msg("Failed to record play event")
		responses.RespondInternalError(c, err, "Failed to record play event")
		return
	}

	log.Info().
		Uint64("userID", req.UserID).
		Uint64("itemID", req.MediaItemID).
		Msg("Play event recorded successfully")

	responses.RespondCreated(c, result, "Play event recorded successfully")
}

// ToggleFavorite godoc
//
//	@Summary		Toggle favorite status for a media item
//	@Description	Mark or unmark a media item as a favorite
//	@Tags			user-data
//	@Accept			json
//	@Produce		json
//	@Param			favorite	query		bool																true	"Favorite status"
//	@Param			mediaType	path		string																true	"Media type like movie, series, track, etc."
//	@Param			itemID		path		int																	true	"Media Item ID"
//	@Param			userID		query		int																	false	"User ID (optional, uses authenticated user ID if not provided)"
//	@Success		200			{object}	responses.APIResponse[models.UserMediaItemData[mediatypes.MediaData]]	"Favorite status updated successfully"
//	@Failure		400			{object}	responses.ErrorResponse[responses.ErrorDetails]						"Bad request"
//	@Failure		401			{object}	responses.ErrorResponse[responses.ErrorDetails]						"Unauthorized"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]						"Internal server error"
//	@Router			/user-data/{mediaType}/{itemID}/favorite [put]
func (h *userMediaItemDataHandler[T]) ToggleFavorite(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	itemID, err := strconv.ParseUint(c.Param("itemID"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("itemID", c.Param("itemID")).Msg("Invalid media item ID")
		responses.RespondBadRequest(c, err, "Invalid media item ID")
		return
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		log.Warn().Err(err).Str("userID", c.Query("userID")).Msg("Invalid user ID")
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
		Uint64("userID", userID).
		Uint64("itemID", itemID).
		Bool("favorite", favorite).
		Msg("Toggling favorite status")

	err = h.service.ToggleFavorite(ctx, itemID, userID, favorite)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("itemID", itemID).
			Msg("Failed to update favorite status")
		responses.RespondInternalError(c, err, "Failed to update favorite status")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("itemID", itemID).
		Bool("favorite", favorite).
		Msg("Favorite status updated successfully")

	responses.RespondOK(c, http.StatusOK, "Favorite status updated successfully")
}

// UpdateUserRating godoc
//
//	@Summary		Update user rating for a media item
//	@Description	Set a user's rating for a media item
//	@Tags			user-data
//	@Accept			json
//	@Produce		json
//	@Param			rating		query		number																true	"User rating (0-10)"
//	@Param			mediaType	path		string																true	"Media type like movie, series, track, etc."
//	@Param			itemID		path		int																	true	"Media Item ID"
//	@Param			userID		query		int																	false	"User ID (optional, uses authenticated user ID if not provided)"
//	@Success		200			{object}	responses.APIResponse[models.UserMediaItemData[mediatypes.MediaData]]	"Rating updated successfully"
//	@Failure		400			{object}	responses.ErrorResponse[responses.ErrorDetails]						"Bad request"
//	@Failure		401			{object}	responses.ErrorResponse[responses.ErrorDetails]						"Unauthorized"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]						"Internal server error"
//	@Router			/user-data/{mediaType}/{itemID}/rating [put]
func (h *userMediaItemDataHandler[T]) UpdateUserRating(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	itemID, err := strconv.ParseUint(c.Param("itemID"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("itemID", c.Param("itemID")).Msg("Invalid media item ID")
		responses.RespondBadRequest(c, err, "Invalid media item ID")
		return
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		log.Warn().Err(err).Str("userID", c.Query("userID")).Msg("Invalid user ID")
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
		Uint64("userID", userID).
		Uint64("itemID", itemID).
		Float64("rating", rating).
		Msg("Updating user rating")

	err = h.service.UpdateRating(ctx, itemID, userID, float32(rating))
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("itemID", itemID).
			Msg("Failed to update rating")
		responses.RespondInternalError(c, err, "Failed to update rating")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("itemID", itemID).
		Float64("rating", rating).
		Msg("Rating updated successfully")

	responses.RespondOK(c, http.StatusOK, "Rating updated successfully")
}

// GetFavorites godoc
//
//	@Summary		Get a user's favorite media items
//	@Description	Get all media items marked as favorites by a user
//	@Tags			user-data
//	@Accept			json
//	@Produce		json
//	@Param			userID		query		int																	false	"User ID (optional, uses authenticated user ID if not provided)"
//	@Param			limit		query		int																	false	"Number of items to return (default 10)"
//	@Param			offset		query		int																	false	"Number of items to skip (default 0)"
//	@Param			mediaType	path		string																true	"Media type like movie, series, track, etc."
//	@Success		200			{object}	responses.APIResponse[[]models.UserMediaItemData[mediatypes.MediaData]]	"Successfully retrieved favorites"
//	@Failure		400			{object}	responses.ErrorResponse[responses.ErrorDetails]						"Bad request"
//	@Failure		401			{object}	responses.ErrorResponse[responses.ErrorDetails]						"Unauthorized"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]						"Internal server error"
//	@Router			/user-data/{mediaType}/favorites [get]
func (h *userMediaItemDataHandler[T]) GetFavorites(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	userID, err := utils.GetUserID(c)
	if err != nil {
		log.Warn().Err(err).Str("userID", c.Query("userID")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	limit := utils.GetLimit(c, 10, 100, true)
	offset := utils.GetOffset(c, 0)

	log.Debug().
		Uint64("userID", userID).
		Int("limit", limit).
		Int("offset", offset).
		Msg("Getting user favorites")

	favorites, err := h.service.GetFavorites(ctx, userID, limit, offset)
	if err != nil {
		log.Error().Err(err).Uint64("userID", userID).Msg("Failed to retrieve favorites")
		responses.RespondInternalError(c, err, "Failed to retrieve favorites")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Int("count", len(favorites)).
		Msg("Favorites retrieved successfully")

	responses.RespondOK(c, favorites, "Favorites retrieved successfully")
}

// ClearUserHistory godoc
//
//	@Summary		Clear a user's play user-data
//	@Description	Delete all play user-data entries for a user
//	@Tags			user-data
//	@Accept			json
//	@Produce		json
//	@Param			userID		query		int												false	"User ID (optional, uses authenticated user ID if not provided)"
//	@Param			mediaType	path		string											true	"Media type like movie, series, track, etc."
//	@Success		200			{object}	responses.APIResponse[any]						"History cleared successfully"
//	@Failure		400			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Bad request"
//	@Failure		401			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Unauthorized"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Internal server error"
//	@Router			/user-data/{mediaType}/clear [delete]
func (h *userMediaItemDataHandler[T]) ClearUserHistory(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	userID, err := utils.GetUserID(c)
	if err != nil {
		log.Warn().Err(err).Str("userID", c.Query("userID")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	log.Debug().
		Uint64("userID", userID).
		Msg("Clearing user user-data")

	err = h.service.ClearUserHistory(ctx, userID)
	if err != nil {
		log.Error().Err(err).Uint64("userID", userID).Msg("Failed to clear user-data")
		responses.RespondInternalError(c, err, "Failed to clear user-data")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Msg("History cleared successfully")

	responses.RespondOK(c, responses.EmptyResponse{Success: true}, "History cleared successfully")
}
