package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"suasor/clients/media/types"
	"suasor/services"
	_ "suasor/types/models"
	"suasor/types/responses"
	"suasor/utils"
	"suasor/utils/logger"
)

type CoreUserMediaItemDataHandler[T types.MediaData] interface {
	Service() services.CoreUserMediaItemDataService[T]

	GetMediaItemDataByID(c *gin.Context)
	CheckUserMediaItemData(c *gin.Context)
	GetUserMediaItemDataByItemID(c *gin.Context)
	DeleteMediaItemData(c *gin.Context)
}

// coreUserMediaItemDataHandler handles basic CRUD operations for user media item data
// This is the core layer of the three-pronged architecture
type coreUserMediaItemDataHandler[T types.MediaData] struct {
	service services.CoreUserMediaItemDataService[T]
}

// NewCoreUserMediaItemDataHandler creates a new core user media item data handler
func NewCoreUserMediaItemDataHandler[T types.MediaData](
	service services.CoreUserMediaItemDataService[T],
) CoreUserMediaItemDataHandler[T] {
	return &coreUserMediaItemDataHandler[T]{
		service: service,
	}
}

// Service returns the underlying service
func (h *coreUserMediaItemDataHandler[T]) Service() services.CoreUserMediaItemDataService[T] {
	return h.service
}

// GetMediaItemDataByID godoc
//
//	@Summary		Get a specific user media item data entry by ID
//	@Description	Retrieves a specific user media item data entry by its ID
//	@Tags			user-data, core
//	@Accept			json
//	@Produce		json
//	@Param			dataID	path		int																	true	"User Media Item Data ID"
//	@Success		200		{object}	responses.APIResponse[models.UserMediaItemData[types.MediaData]]	"Successfully retrieved user media item data"
//	@Failure		400		{object}	responses.ErrorResponse[responses.ErrorDetails]						"Bad request"
//	@Failure		404		{object}	responses.ErrorResponse[responses.ErrorDetails]						"Not found"
//	@Failure		500		{object}	responses.ErrorResponse[responses.ErrorDetails]						"Internal server error"
//	@Router			/user-data/data/{dataID} [get]
func (h *coreUserMediaItemDataHandler[T]) GetMediaItemDataByID(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	dataID, _ := checkItemID(c, "dataID")

	log.Debug().Uint64("id", dataID).Msg("Getting user media item data by ID")

	data, err := h.service.GetByID(ctx, dataID)
	if err != nil {
		log.Error().Err(err).Uint64("id", dataID).Msg("Failed to get user media item data")
		responses.RespondNotFound(c, err, "User media item data not found")
		return
	}

	log.Info().Uint64("id", dataID).Msg("User media item data retrieved successfully")
	responses.RespondOK(c, data, "User media item data retrieved successfully")
}

// CheckUserMediaItemData godoc
//
//	@Summary		Check if a user has data for a specific media item
//	@Description	Checks if a user has data for a specific media item
//	@Tags			user-data, core
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int												true	"Media Item ID"
//	@Param			mediaType	path		string											true	"Media type like movie, series, track, etc."
//	@Param			userID		query		int												false	"User ID (optional, uses authenticated user ID if not provided)"
//	@Success		200			{object}	responses.APIResponse[bool]						"Successfully checked user media item data"
//	@Failure		400			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Bad request"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Internal server error"
//	@Router			/user-data/{mediaType}/{id}/check [get]
func (h *coreUserMediaItemDataHandler[T]) CheckUserMediaItemData(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	userID, err := utils.GetUserID(c)
	if err != nil {
		log.Warn().Err(err).Str("userID", c.Query("userID")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	itemID, err := strconv.ParseUint(c.Query("itemID"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("itemID", c.Query("itemID")).Msg("Invalid media item ID")
		responses.RespondBadRequest(c, err, "Invalid media item ID")
		return
	}

	log.Debug().
		Uint64("userID", userID).
		Uint64("itemID", itemID).
		Msg("Checking if user has media item data")

	hasData, err := h.service.HasUserMediaItemData(ctx, userID, itemID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("itemID", itemID).
			Msg("Failed to check user media item data")
		responses.RespondInternalError(c, err, "Failed to check user media item data")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("itemID", itemID).
		Bool("hasData", hasData).
		Msg("User media item data check completed")

	responses.RespondOK(c, hasData, "User media item data check completed")
}

// GetUserMediaItemDataByItemID godoc
//
//	@Summary		Get user media item data for a specific user and media item
//	@Description	Retrieves user media item data for a specific user and media item
//	@Tags			user-data, core
//	@Accept			json
//	@Produce		json
//	@Param			itemID		path		int																	true	"Media Item ID"
//	@Param			mediaType	path		string																true	"Media type"
//	@Param			userID		query		int																	false	"User ID (optional, uses authenticated user ID if not provided)"
//	@Success		200			{object}	responses.APIResponse[models.UserMediaItemData[types.MediaData]]	"Successfully retrieved user media item data"
//	@Failure		400			{object}	responses.ErrorResponse[responses.ErrorDetails]						"Bad request"
//	@Failure		404			{object}	responses.ErrorResponse[responses.ErrorDetails]						"Not found"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]						"Internal server error"
//	@Router			/user-data/{mediaType}/{itemID} [get]
func (h *coreUserMediaItemDataHandler[T]) GetUserMediaItemDataByItemID(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	userID, _ := checkUserAccess(c)

	itemID, err := strconv.ParseUint(c.Query("itemID"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("itemID", c.Query("itemID")).Msg("Invalid media item ID")
		responses.RespondBadRequest(c, err, "Invalid media item ID")
		return
	}

	log.Debug().
		Uint64("userID", userID).
		Uint64("itemID", itemID).
		Msg("Getting user media item data")

	data, err := h.service.GetByUserIDAndMediaItemID(ctx, userID, itemID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("itemID", itemID).
			Msg("Failed to get user media item data")
		responses.RespondNotFound(c, err, "User media item data not found")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("itemID", itemID).
		Msg("User media item data retrieved successfully")

	responses.RespondOK(c, data, "User media item data retrieved successfully")
}

// DeleteMediaItemData godoc
//
//	@Summary		Delete a specific user media item data entry
//	@Description	Deletes a specific user media item data entry by its ID
//	@Tags			user-data, core
//	@Accept			json
//	@Produce		json
//	@Param			itemID		path		int																	true	"User Media Item ID"
//	@Param			mediaType	path		string																true	"Media type like movie, series, track, etc."
//	@Param			userID		query		int																	false	"User ID (optional, uses authenticated user ID if not provided)"
//	@Success		200			{object}	responses.APIResponse[models.UserMediaItemData[types.MediaData]]	"Successfully deleted user media item data"
//	@Failure		400			{object}	responses.ErrorResponse[responses.ErrorDetails]						"Bad request"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]						"Internal server error"
//	@Router			/user-data/{mediaType}/{itemID} [delete]
func (h *coreUserMediaItemDataHandler[T]) DeleteMediaItemData(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	userID, _ := checkUserAccess(c)
	itemID, _ := checkItemID(c, "itemitemID")

	log.Debug().Uint64("userID", userID).Uint64("itemID", itemID).Msg("Deleting user media item data")

	err := h.service.Delete(ctx, itemID)
	if err != nil {
		log.Error().Err(err).Uint64("itemID", itemID).Msg("Failed to delete user media item data")
		responses.RespondInternalError(c, err, "Failed to delete user media item data")
		return
	}

	log.Info().Uint64("itemID", itemID).Msg("User media item data deleted successfully")
	responses.RespondOK(c, responses.EmptyResponse{Success: true}, "User media item data deleted successfully")
}
