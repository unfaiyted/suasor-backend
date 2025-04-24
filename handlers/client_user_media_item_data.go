package handlers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	mediatypes "suasor/clients/media/types"
	clienttypes "suasor/clients/types"
	"suasor/services"
	"suasor/types/models"
	"suasor/types/requests"
	"suasor/types/responses"
	"suasor/utils"
	"suasor/utils/logger"
)

// Note: Type parameters for interface and implementation
type ClientUserMediaItemDataHandler[T clienttypes.ClientMediaConfig, U mediatypes.MediaData] interface {
	UserMediaItemDataHandler[U]

	SyncClientItemData(c *gin.Context)
	GetClientItemData(c *gin.Context)
	GetMediaItemDataByClientID(c *gin.Context)
	RecordClientPlay(c *gin.Context)
	GetPlaybackState(c *gin.Context)
	UpdatePlaybackState(c *gin.Context)
}

// clientUserMediaItemDataHandler handles client-specific operations for user media item data
// This is the client layer of the three-pronged architecture
type clientUserMediaItemDataHandler[T clienttypes.ClientMediaConfig, U mediatypes.MediaData] struct {
	UserMediaItemDataHandler[U]

	service services.ClientUserMediaItemDataService[T, U]
}

// NewclientUserMediaItemDataHandler creates a new client user media item data handler
func NewClientUserMediaItemDataHandler[T clienttypes.ClientMediaConfig, U mediatypes.MediaData](
	userHandler UserMediaItemDataHandler[U],
	service services.ClientUserMediaItemDataService[T, U],

) *clientUserMediaItemDataHandler[T, U] {
	return &clientUserMediaItemDataHandler[T, U]{
		UserMediaItemDataHandler: userHandler,
		service:                  service,
	}
}

// SyncClientItemData godoc
// @Summary Synchronize user media item data from a client
// @Description Synchronizes user media item data from an external client
// @Tags user-data, clients
// @Accept json
// @Produce json
// @Param clientID path int true "Client ID"
// @Param userID query int false "User ID (optional, uses authenticated user ID if not provided)"
// @Param items body requests.UserMediaItemDataSyncRequest true "Media item data to synchronize"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Bad request"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Internal server error"
// @Router /api/v1/client/{clientID}/user-data/{mediaType}/sync [post]
func (h *clientUserMediaItemDataHandler[T, U]) SyncClientItemData(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("clientID", c.Param("clientID")).Msg("Invalid client ID")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		log.Warn().Err(err).Str("userID", c.Query("userID")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	var items []models.UserMediaItemData[U]
	if err := c.ShouldBindJSON(&items); err != nil {
		log.Warn().Err(err).Msg("Invalid request body")
		responses.RespondBadRequest(c, err, "Invalid request body")
		return
	}

	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Int("itemCount", len(items)).
		Msg("Synchronizing client media item data")

	err = h.service.SyncClientItemData(ctx, userID, clientID, items)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Msg("Failed to synchronize client media item data")
		responses.RespondInternalError(c, err, "Failed to synchronize client media item data")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Int("itemCount", len(items)).
		Msg("Client media item data synchronized successfully")

	responses.RespondOK(c, responses.EmptyResponse{Success: true}, "Client media item data synchronized successfully")
}

// GetClientItemData godoc
// @Summary Get user media item data for a client
// @Description Retrieves user media item data for synchronization with a client
// @Tags user-data, clients
// @Accept json
// @Produce json
// @Param clientID path int true "Client ID"
// @Param mediaType path string true "Media type like movie, series, track, etc."
// @Param userID query int false "User ID (optional, uses authenticated user ID if not provided)"
// @Param since query string false "Since date (default 24 hours ago)"
// @Success 200 {object} responses.APIResponse[[]models.UserMediaItemData[mediatypes.Movie]] "Successfully retrieved client media item data"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Bad request"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Internal server error"
// @Router /api/v1/client/{clientID}/user-data/{mediaType} [get]
func (h *clientUserMediaItemDataHandler[T, U]) GetClientItemData(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("clientID", c.Param("clientID")).Msg("Invalid client ID")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		log.Warn().Err(err).Str("userID", c.Query("userID")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	sinceStr := c.Query("since")
	var since *string
	if sinceStr != "" {
		since = &sinceStr
	}

	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("since", sinceStr).
		Msg("Getting client media item data")

	items, err := h.service.GetClientItemData(ctx, userID, clientID, since)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Msg("Failed to get client media item data")
		responses.RespondInternalError(c, err, "Failed to get client media item data")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Int("count", len(items)).
		Msg("Client media item data retrieved successfully")

	responses.RespondOK(c, items, "Client media item data retrieved successfully")
}

// GetMediaItemDataByClientID godoc
// @Summary Get user media item data by client ID
// @Description Retrieves user media item data for a specific user and client item
// @Tags user-data, clients
// @Accept json
// @Produce json
// @Param clientID path int true "Client ID"
// @Param clientItemID path string true "Client Item ID"
// @Param userID query int false "User ID (optional, uses authenticated user ID if not provided)"
// @Success 200 {object} responses.APIResponse[models.UserMediaItemData[mediatypes.Movie]] "Successfully retrieved user media item data"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Bad request"
// @Failure 404 {object} responses.ErrorResponse[responses.ErrorDetails] "Not found"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Internal server error"
// @Router /api/v1/client/{clientID}/user-data/{mediaType}/{clientItemID} [get]
func (h *clientUserMediaItemDataHandler[T, U]) GetMediaItemDataByClientID(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("clientID", c.Param("clientID")).Msg("Invalid client ID")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	clientItemID := c.Param("clientItemID")
	if clientItemID == "" {
		log.Warn().Msg("Client item ID is required")
		responses.RespondBadRequest(c, nil, "Client item ID is required")
		return
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		log.Warn().Err(err).Str("userID", c.Query("userID")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientItemID", clientItemID).
		Msg("Getting user media item data by client ID")

	data, err := h.service.GetByClientID(ctx, userID, clientID, clientItemID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("clientItemID", clientItemID).
			Msg("Failed to get user media item data by client ID")
		responses.RespondNotFound(c, err, "User media item data not found")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientItemID", clientItemID).
		Msg("User media item data retrieved successfully")

	responses.RespondOK(c, data, "User media item data retrieved successfully")
}

// RecordClientPlay godoc
// @Summary Record a client play event
// @Description Records a play event from a client
// @Tags user-data, clients
// @Accept json
// @Produce json
// @Param clientID path int true "Client ID"
// @Param clientItemID path string true "Client Item ID"
// @Param userID query int false "User ID (optional, uses authenticated user ID if not provided)"
// @Param mediaPlay body requests.UserMediaItemDataRequest true "Media play information"
// @Success 201 {object} responses.APIResponse[models.UserMediaItemData[mediatypes.Movie]] "Play event recorded successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Bad request"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Internal server error"
// @Router /api/v1/client/{clientID}/user-data/{mediaType}/{clientItemID}/play [post]
func (h *clientUserMediaItemDataHandler[T, U]) RecordClientPlay(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("clientID", c.Param("clientID")).Msg("Invalid client ID")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	clientItemID := c.Param("clientItemID")
	if clientItemID == "" {
		log.Warn().Msg("Client item ID is required")
		responses.RespondBadRequest(c, nil, "Client item ID is required")
		return
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		log.Warn().Err(err).Str("userID", c.Query("userID")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	var req requests.UserMediaItemDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid request body")
		responses.RespondBadRequest(c, err, "Invalid request body")
		return
	}

	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientItemID", clientItemID).
		Msg("Recording client play event")

	// Create a play history record
	playHistory := &models.UserMediaItemData[U]{
		UserID:           userID,
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

	result, err := h.service.RecordClientPlay(ctx, userID, clientID, clientItemID, playHistory)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("clientItemID", clientItemID).
			Msg("Failed to record client play event")
		responses.RespondInternalError(c, err, "Failed to record client play event")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientItemID", clientItemID).
		Msg("Client play event recorded successfully")

	responses.RespondCreated(c, result, "Client play event recorded successfully")
}

// GetPlaybackState godoc
// @Summary Get playback state for a client item
// @Description Retrieves the current playback state for a client item
// @Tags user-data, clients
// @Accept json
// @Produce json
// @Param userID query int false "User ID (optional, uses authenticated user ID if not provided)"
// @Param clientID path int true "Client ID"
// @Param clientItemID path string true "Client Item ID"
// @Param mediaType path string true "Media type like movie, series, track, etc."
// @Success 200 {object} responses.APIResponse[models.UserMediaItemData[mediatypes.Movie]] "Successfully retrieved playback state"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Bad request"
// @Failure 404 {object} responses.ErrorResponse[responses.ErrorDetails] "Not found"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Internal server error"
// @Router /api/v1/client/{clientID}/user-data/{mediaType}/{clientItemID}/state [get]
func (h *clientUserMediaItemDataHandler[T, U]) GetPlaybackState(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("clientID", c.Param("clientID")).Msg("Invalid client ID")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	clientItemID := c.Param("clientItemID")
	if clientItemID == "" {
		log.Warn().Msg("Client item ID is required")
		responses.RespondBadRequest(c, nil, "Client item ID is required")
		return
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		log.Warn().Err(err).Str("userID", c.Query("userID")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientItemID", clientItemID).
		Msg("Getting playback state")

	state, err := h.service.GetPlaybackState(ctx, userID, clientID, clientItemID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("clientItemID", clientItemID).
			Msg("Failed to get playback state")
		responses.RespondNotFound(c, err, "Playback state not found")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientItemID", clientItemID).
		Msg("Playback state retrieved successfully")

	responses.RespondOK(c, state, "Playback state retrieved successfully")
}

// UpdatePlaybackState godoc
// @Summary Update playback state for a client item
// @Description Updates the playback state for a client item
// @Tags user-data, clients
// @Accept json
// @Produce json
// @Param userID query int false "User ID (optional, uses authenticated user ID if not provided)"
// @Param clientID path int true "Client ID"
// @Param clientItemID path string true "Client Item ID"
// @Param state body object true "Playback state information"
// @Success 200 {object} responses.APIResponse[models.UserMediaItemData[mediatypes.Movie]] "Playback state updated successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Bad request"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Internal server error"
// @Router /api/v1/client/{clientID}/user-data/{mediaType}/{clientItemID}/state [put]
func (h *clientUserMediaItemDataHandler[T, U]) UpdatePlaybackState(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("clientID", c.Param("clientID")).Msg("Invalid client ID")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	clientItemID := c.Param("clientItemID")
	if clientItemID == "" {
		log.Warn().Msg("Client item ID is required")
		responses.RespondBadRequest(c, nil, "Client item ID is required")
		return
	}

	userID, err := utils.GetUserID(c)
	if err != nil {
		log.Warn().Err(err).Str("userID", c.Query("userID")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	var req struct {
		Position   int     `json:"positionSeconds" binding:"required"`
		Duration   int     `json:"durationSeconds" binding:"required"`
		Percentage float64 `json:"playedPercentage" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid request body")
		responses.RespondBadRequest(c, err, "Invalid request body")
		return
	}

	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientItemID", clientItemID).
		Int("position", req.Position).
		Int("duration", req.Duration).
		Float64("percentage", req.Percentage).
		Msg("Updating playback state")

	result, err := h.service.UpdatePlaybackState(ctx, userID, clientID, clientItemID, req.Position, req.Duration, req.Percentage)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("clientItemID", clientItemID).
			Msg("Failed to update playback state")
		responses.RespondInternalError(c, err, "Failed to update playback state")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientItemID", clientItemID).
		Msg("Playback state updated successfully")

	responses.RespondOK(c, result, "Playback state updated successfully")
}
