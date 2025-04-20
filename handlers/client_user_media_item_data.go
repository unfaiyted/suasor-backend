package handlers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"suasor/client/media/types"
	"suasor/services"
	"suasor/types/models"
	"suasor/types/requests"
	"suasor/types/responses"
	"suasor/utils"
)

type ClientUserMediaItemDataHandler[T types.MediaData] interface {
	UserMediaItemDataHandler[T]

	SyncClientItemData(c *gin.Context)
	GetClientItemData(c *gin.Context)
	GetMediaItemDataByClientID(c *gin.Context)
	RecordClientPlay(c *gin.Context)
	GetPlaybackState(c *gin.Context)
	UpdatePlaybackState(c *gin.Context)
}

// clientUserMediaItemDataHandler handles client-specific operations for user media item data
// This is the client layer of the three-pronged architecture
type clientUserMediaItemDataHandler[T types.MediaData] struct {
	UserMediaItemDataHandler[T]

	service services.ClientUserMediaItemDataService[T]
}

// NewclientUserMediaItemDataHandler creates a new client user media item data handler
func NewClientUserMediaItemDataHandler[T types.MediaData](
	userHandler *UserMediaItemDataHandler[T],
	service services.ClientUserMediaItemDataService[T],

) *clientUserMediaItemDataHandler[T] {
	return &clientUserMediaItemDataHandler[T]{
		UserMediaItemDataHandler: *userHandler,
		service:                  service,
	}
}

// SyncClientItemData godoc
// @Summary Synchronize user media item data from a client
// @Description Synchronizes user media item data from an external client
// @Tags History
// @Accept json
// @Produce json
// @Param clientId path int true "Client ID"
// @Param userId query int true "User ID"
// @Param items body []models.UserMediaItemData[any] true "Media item data to sync"
// @Success 200 {object} responses.APIResponse[any] "Successfully synchronized client media item data"
// @Failure 400 {object} responses.ErrorResponse[any] "Bad request"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[any] "Internal server error"
// @Router /user-media-data/client/{clientId}/sync [post]
func (h *clientUserMediaItemDataHandler[T]) SyncClientItemData(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	clientID, err := strconv.ParseUint(c.Param("clientId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("clientId", c.Param("clientId")).Msg("Invalid client ID")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	var items []models.UserMediaItemData[T]
	if err := c.ShouldBindJSON(&items); err != nil {
		log.Warn().Err(err).Msg("Invalid request body")
		responses.RespondBadRequest(c, err, "Invalid request body")
		return
	}

	log.Debug().
		Uint64("userId", userID).
		Uint64("clientId", clientID).
		Int("itemCount", len(items)).
		Msg("Synchronizing client media item data")

	err = h.service.SyncClientItemData(ctx, userID, clientID, items)
	if err != nil {
		log.Error().Err(err).
			Uint64("userId", userID).
			Uint64("clientId", clientID).
			Msg("Failed to synchronize client media item data")
		responses.RespondInternalError(c, err, "Failed to synchronize client media item data")
		return
	}

	log.Info().
		Uint64("userId", userID).
		Uint64("clientId", clientID).
		Int("itemCount", len(items)).
		Msg("Client media item data synchronized successfully")

	responses.RespondOK(c, responses.EmptyResponse{Success: true}, "Client media item data synchronized successfully")
}

// GetClientItemData godoc
// @Summary Get user media item data for a client
// @Description Retrieves user media item data for synchronization with a client
// @Tags History
// @Accept json
// @Produce json
// @Param clientId path int true "Client ID"
// @Param userId query int true "User ID"
// @Param since query string false "Since date (default 24 hours ago)"
// @Success 200 {object} responses.APIResponse[[]models.UserMediaItemData[any]] "Successfully retrieved client media item data"
// @Failure 400 {object} responses.ErrorResponse[any] "Bad request"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[any] "Internal server error"
// @Router /user-media-data/client/{clientId} [get]
func (h *clientUserMediaItemDataHandler[T]) GetClientItemData(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	clientID, err := strconv.ParseUint(c.Param("clientId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("clientId", c.Param("clientId")).Msg("Invalid client ID")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	sinceStr := c.Query("since")
	var since *string
	if sinceStr != "" {
		since = &sinceStr
	}

	log.Debug().
		Uint64("userId", userID).
		Uint64("clientId", clientID).
		Str("since", sinceStr).
		Msg("Getting client media item data")

	items, err := h.service.GetClientItemData(ctx, userID, clientID, since)
	if err != nil {
		log.Error().Err(err).
			Uint64("userId", userID).
			Uint64("clientId", clientID).
			Msg("Failed to get client media item data")
		responses.RespondInternalError(c, err, "Failed to get client media item data")
		return
	}

	log.Info().
		Uint64("userId", userID).
		Uint64("clientId", clientID).
		Int("count", len(items)).
		Msg("Client media item data retrieved successfully")

	responses.RespondOK(c, items, "Client media item data retrieved successfully")
}

// GetMediaItemDataByClientID godoc
// @Summary Get user media item data by client ID
// @Description Retrieves user media item data for a specific user and client item
// @Tags History
// @Accept json
// @Produce json
// @Param clientId path int true "Client ID"
// @Param clientItemId path string true "Client Item ID"
// @Param userId query int true "User ID"
// @Success 200 {object} responses.APIResponse[models.UserMediaItemData[any]] "Successfully retrieved user media item data"
// @Failure 400 {object} responses.ErrorResponse[any] "Bad request"
// @Failure 404 {object} responses.ErrorResponse[any] "Not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Internal server error"
// @Router /user-media-data/client/{clientId}/item/{clientItemId} [get]
func (h *clientUserMediaItemDataHandler[T]) GetMediaItemDataByClientID(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	clientID, err := strconv.ParseUint(c.Param("clientId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("clientId", c.Param("clientId")).Msg("Invalid client ID")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	clientItemID := c.Param("clientItemId")
	if clientItemID == "" {
		log.Warn().Msg("Client item ID is required")
		responses.RespondBadRequest(c, nil, "Client item ID is required")
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	log.Debug().
		Uint64("userId", userID).
		Uint64("clientId", clientID).
		Str("clientItemId", clientItemID).
		Msg("Getting user media item data by client ID")

	data, err := h.service.GetByClientID(ctx, userID, clientID, clientItemID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userId", userID).
			Uint64("clientId", clientID).
			Str("clientItemId", clientItemID).
			Msg("Failed to get user media item data by client ID")
		responses.RespondNotFound(c, err, "User media item data not found")
		return
	}

	log.Info().
		Uint64("userId", userID).
		Uint64("clientId", clientID).
		Str("clientItemId", clientItemID).
		Msg("User media item data retrieved successfully")

	responses.RespondOK(c, data, "User media item data retrieved successfully")
}

// RecordClientPlay godoc
// @Summary Record a client play event
// @Description Records a play event from a client
// @Tags History
// @Accept json
// @Produce json
// @Param clientId path int true "Client ID"
// @Param clientItemId path string true "Client Item ID"
// @Param userId query int true "User ID"
// @Param mediaPlay body requests.UserMediaItemDataRequest true "Media play information"
// @Success 201 {object} responses.APIResponse[models.UserMediaItemData[any]] "Play event recorded successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Bad request"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[any] "Internal server error"
// @Router /user-media-data/client/{clientId}/item/{clientItemId}/play [post]
func (h *clientUserMediaItemDataHandler[T]) RecordClientPlay(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	clientID, err := strconv.ParseUint(c.Param("clientId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("clientId", c.Param("clientId")).Msg("Invalid client ID")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	clientItemID := c.Param("clientItemId")
	if clientItemID == "" {
		log.Warn().Msg("Client item ID is required")
		responses.RespondBadRequest(c, nil, "Client item ID is required")
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
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
		Uint64("userId", userID).
		Uint64("clientId", clientID).
		Str("clientItemId", clientItemID).
		Msg("Recording client play event")

	// Create a play history record
	playHistory := &models.UserMediaItemData[T]{
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
			Uint64("userId", userID).
			Uint64("clientId", clientID).
			Str("clientItemId", clientItemID).
			Msg("Failed to record client play event")
		responses.RespondInternalError(c, err, "Failed to record client play event")
		return
	}

	log.Info().
		Uint64("userId", userID).
		Uint64("clientId", clientID).
		Str("clientItemId", clientItemID).
		Msg("Client play event recorded successfully")

	responses.RespondCreated(c, result, "Client play event recorded successfully")
}

// GetPlaybackState godoc
// @Summary Get playback state for a client item
// @Description Retrieves the current playback state for a client item
// @Tags History
// @Accept json
// @Produce json
// @Param clientId path int true "Client ID"
// @Param clientItemId path string true "Client Item ID"
// @Param userId query int true "User ID"
// @Success 200 {object} responses.APIResponse[models.UserMediaItemData[any]] "Successfully retrieved playback state"
// @Failure 400 {object} responses.ErrorResponse[any] "Bad request"
// @Failure 404 {object} responses.ErrorResponse[any] "Not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Internal server error"
// @Router /user-media-data/client/{clientId}/item/{clientItemId}/state [get]
func (h *clientUserMediaItemDataHandler[T]) GetPlaybackState(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	clientID, err := strconv.ParseUint(c.Param("clientId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("clientId", c.Param("clientId")).Msg("Invalid client ID")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	clientItemID := c.Param("clientItemId")
	if clientItemID == "" {
		log.Warn().Msg("Client item ID is required")
		responses.RespondBadRequest(c, nil, "Client item ID is required")
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	log.Debug().
		Uint64("userId", userID).
		Uint64("clientId", clientID).
		Str("clientItemId", clientItemID).
		Msg("Getting playback state")

	state, err := h.service.GetPlaybackState(ctx, userID, clientID, clientItemID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userId", userID).
			Uint64("clientId", clientID).
			Str("clientItemId", clientItemID).
			Msg("Failed to get playback state")
		responses.RespondNotFound(c, err, "Playback state not found")
		return
	}

	log.Info().
		Uint64("userId", userID).
		Uint64("clientId", clientID).
		Str("clientItemId", clientItemID).
		Msg("Playback state retrieved successfully")

	responses.RespondOK(c, state, "Playback state retrieved successfully")
}

// UpdatePlaybackState godoc
// @Summary Update playback state for a client item
// @Description Updates the playback state for a client item
// @Tags History
// @Accept json
// @Produce json
// @Param clientId path int true "Client ID"
// @Param clientItemId path string true "Client Item ID"
// @Param userId query int true "User ID"
// @Param state body object true "Playback state information"
// @Success 200 {object} responses.APIResponse[models.UserMediaItemData[any]] "Playback state updated successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Bad request"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[any] "Internal server error"
// @Router /user-media-data/client/{clientId}/item/{clientItemId}/state [put]
func (h *clientUserMediaItemDataHandler[T]) UpdatePlaybackState(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	clientID, err := strconv.ParseUint(c.Param("clientId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("clientId", c.Param("clientId")).Msg("Invalid client ID")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	clientItemID := c.Param("clientItemId")
	if clientItemID == "" {
		log.Warn().Msg("Client item ID is required")
		responses.RespondBadRequest(c, nil, "Client item ID is required")
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
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
		Uint64("userId", userID).
		Uint64("clientId", clientID).
		Str("clientItemId", clientItemID).
		Int("position", req.Position).
		Int("duration", req.Duration).
		Float64("percentage", req.Percentage).
		Msg("Updating playback state")

	result, err := h.service.UpdatePlaybackState(ctx, userID, clientID, clientItemID, req.Position, req.Duration, req.Percentage)
	if err != nil {
		log.Error().Err(err).
			Uint64("userId", userID).
			Uint64("clientId", clientID).
			Str("clientItemId", clientItemID).
			Msg("Failed to update playback state")
		responses.RespondInternalError(c, err, "Failed to update playback state")
		return
	}

	log.Info().
		Uint64("userId", userID).
		Uint64("clientId", clientID).
		Str("clientItemId", clientItemID).
		Msg("Playback state updated successfully")

	responses.RespondOK(c, result, "Playback state updated successfully")
}
