package handlers

//
// import (
// 	"strconv"
// 	"time"
//
// 	"github.com/gin-gonic/gin"
//
// 	"suasor/client/media/types"
// 	"suasor/services"
// 	"suasor/types/models"
// 	"suasor/types/responses"
// )
//
// // UserMediaItemDataHandler handles all media play history operations
// type UserMediaItemDataHandler[T types.MediaData] struct {
// 	service services.UserMediaItemDataService[T]
// }
//
// // NewUserMediaItemDataHandler creates a new media play history handler
// func NewUserMediaItemDataHandler[T types.MediaData](service services.UserMediaItemDataService[T]) *UserMediaItemDataHandler[T] {
// 	return &UserMediaItemDataHandler[T]{service: service}
// }
//
// // GetMediaPlayHistory godoc
// // @Summary Get a user's media play history
// // @Description Get a user's media play history with optional filtering
// // @Tags History
// // @Accept json
// // @Produce json
// // @Param userId query int true "User ID"
// // @Param limit query int false "Number of items to return (default 10)"
// // @Param offset query int false "Number of items to skip (default 0)"
// // @Param type query string false "Media type filter (movie, series, episode, track, etc.)"
// // @Param completed query bool false "Filter by completion status"
// // @Success 200 {object} responses.APIResponse[[]models.MediaPlayHistory[any]] "Successfully retrieved play history"
// // @Failure 400 {object} responses.ErrorResponse[any] "Bad request"
// // @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// // @Failure 500 {object} responses.ErrorResponse[any] "Internal server error"
// // @Router /history [get]
// func (h *UserMediaItemDataHandler) GetMediaPlayHistory(c *gin.Context) {
// 	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
// 	if err != nil {
// 		responses.RespondBadRequest(c, err, "Invalid user ID")
// 		return
// 	}
//
// 	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
// 	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
// 	mediaType := c.Query("type")
// 	completedStr := c.Query("completed")
//
// 	var completed *bool
// 	if completedStr != "" {
// 		completedBool, err := strconv.ParseBool(completedStr)
// 		if err != nil {
// 			responses.RespondBadRequest(c, err, "Invalid completed value")
// 			return
// 		}
// 		completed = &completedBool
// 	}
//
// 	// Filter by media type if provided
// 	var typedMediaType *types.MediaType
// 	if mediaType != "" {
// 		mt := types.MediaType(mediaType)
// 		typedMediaType = &mt
// 	}
//
// 	history, err := h.service.GetUserPlayHistory(c.Request.Context(), userID, limit, offset, typedMediaType, completed)
// 	if err != nil {
// 		responses.RespondInternalError(c, err, "Failed to retrieve play history")
// 		return
// 	}
//
// 	responses.RespondOK(c, history, "Play history retrieved successfully")
// }
//
// // GetContinueWatching godoc
// // @Summary Get a user's continue watching list
// // @Description Get media items that a user has started but not completed
// // @Tags History
// // @Accept json
// // @Produce json
// // @Param userId query int true "User ID"
// // @Param limit query int false "Number of items to return (default 10)"
// // @Success 200 {object} responses.APIResponse[[]models.MediaItem[any]] "Successfully retrieved continue watching items"
// // @Failure 400 {object} responses.ErrorResponse[any] "Bad request"
// // @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// // @Failure 500 {object} responses.ErrorResponse[any] "Internal server error"
// // @Router /history/continue-watching [get]
// func (h *UserMediaItemDataHandler) GetContinueWatching(c *gin.Context) {
// 	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
// 	if err != nil {
// 		responses.RespondBadRequest(c, err, "Invalid user ID")
// 		return
// 	}
//
// 	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
//
// 	// Get items that are not completed and have been played recently
// 	items, err := h.service.GetContinueWatching(c.Request.Context(), userID, limit)
// 	if err != nil {
// 		responses.RespondInternalError(c, err, "Failed to retrieve continue watching items")
// 		return
// 	}
//
// 	responses.RespondOK(c, items, "Continue watching items retrieved successfully")
// }
//
// // GetMediaPlayHistoryByID godoc
// // @Summary Get a specific media play history entry
// // @Description Get a specific media play history entry by ID
// // @Tags History
// // @Accept json
// // @Produce json
// // @Param id path int true "History ID"
// // @Success 200 {object} responses.APIResponse[models.MediaPlayHistory[any]] "Successfully retrieved play history entry"
// // @Failure 400 {object} responses.ErrorResponse[any] "Bad request"
// // @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// // @Failure 404 {object} responses.ErrorResponse[any] "Not found"
// // @Failure 500 {object} responses.ErrorResponse[any] "Internal server error"
// // @Router /history/{id} [get]
// func (h *UserMediaItemDataHandler) GetMediaPlayHistoryByID(c *gin.Context) {
// 	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
// 	if err != nil {
// 		responses.RespondBadRequest(c, err, "Invalid history ID")
// 		return
// 	}
//
// 	history, err := h.service.GetByID(c.Request.Context(), id)
// 	if err != nil {
// 		responses.RespondNotFound(c, err, "History entry not found")
// 		return
// 	}
//
// 	responses.RespondOK(c, history, "History entry retrieved successfully")
// }
//
// // GetMediaPlayHistoryByMediaItem godoc
// // @Summary Get play history for a specific media item
// // @Description Get play history entries for a specific media item
// // @Tags History
// // @Accept json
// // @Produce json
// // @Param mediaItemId path int true "Media Item ID"
// // @Param userId query int true "User ID"
// // @Success 200 {object} responses.APIResponse[[]models.MediaPlayHistory[any]] "Successfully retrieved play history for media item"
// // @Failure 400 {object} responses.ErrorResponse[any] "Bad request"
// // @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// // @Failure 500 {object} responses.ErrorResponse[any] "Internal server error"
// // @Router /history/media/{mediaItemId} [get]
// func (h *UserMediaItemDataHandler) GetMediaPlayHistoryByMediaItem(c *gin.Context) {
// 	mediaItemID, err := strconv.ParseUint(c.Param("mediaItemId"), 10, 64)
// 	if err != nil {
// 		responses.RespondBadRequest(c, err, "Invalid media item ID")
// 		return
// 	}
//
// 	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
// 	if err != nil {
// 		responses.RespondBadRequest(c, err, "Invalid user ID")
// 		return
// 	}
//
// 	history, err := h.service.GetByMediaItemID(c.Request.Context(), mediaItemID, userID)
// 	if err != nil {
// 		responses.RespondInternalError(c, err, "Failed to retrieve play history for media item")
// 		return
// 	}
//
// 	responses.RespondOK(c, history, "Play history for media item retrieved successfully")
// }
//
// // RecordMediaPlay godoc
// // @Summary Record a media play event
// // @Description Record a new play event for a media item
// // @Tags History
// // @Accept json
// // @Produce json
// // @Param mediaPlay body models.UserMediaItemDataRequest true "Media play information"
// // @Success 201 {object} responses.APIResponse[models.MediaPlayHistory[any]] "Play event recorded successfully"
// // @Failure 400 {object} responses.ErrorResponse[any] "Bad request"
// // @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// // @Failure 500 {object} responses.ErrorResponse[any] "Internal server error"
// // @Router /history [post]
// func (h *UserMediaItemDataHandler) RecordMediaPlay(c *gin.Context) {
// 	var req models.UserMediaItemDataRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		responses.RespondBadRequest(c, err, "Invalid request body")
// 		return
// 	}
//
// 	// Create a generic play history record
// 	playHistory := &models.UserMediaItemData[mediatypes.Movie]{
// 		UserID:           req.UserID,
// 		MediaItemID:      req.MediaItemID,
// 		Type:             req.Type,
// 		PlayedAt:         time.Now(),
// 		LastPlayedAt:     time.Now(),
// 		IsFavorite:       req.IsFavorite,
// 		UserRating:       req.UserRating,
// 		PlayedPercentage: req.PlayedPercentage,
// 		PositionSeconds:  req.PositionSeconds,
// 		DurationSeconds:  req.DurationSeconds,
// 		Completed:        req.Completed,
// 	}
//
// 	// If this is a continuation, increment the play count
// 	if req.Continued {
// 		playHistory.PlayCount += 1
// 	}
//
// 	result, err := h.service.RecordPlay(c.Request.Context(), playHistory)
// 	if err != nil {
// 		responses.RespondInternalError(c, err, "Failed to record play event")
// 		return
// 	}
//
// 	responses.RespondCreated(c, result, "Play event recorded successfully")
// }
//
// // ToggleFavorite godoc
// // @Summary Toggle favorite status for a media item
// // @Description Mark or unmark a media item as a favorite
// // @Tags History
// // @Accept json
// // @Produce json
// // @Param mediaItemId path int true "Media Item ID"
// // @Param userId query int true "User ID"
// // @Param favorite query bool true "Favorite status"
// // @Success 200 {object} responses.APIResponse[models.MediaPlayHistory[any]] "Favorite status updated successfully"
// // @Failure 400 {object} responses.ErrorResponse[any] "Bad request"
// // @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// // @Failure 500 {object} responses.ErrorResponse[any] "Internal server error"
// // @Router /history/media/{mediaItemId}/favorite [put]
// func (h *UserMediaItemDataHandler) ToggleFavorite(c *gin.Context) {
// 	mediaItemID, err := strconv.ParseUint(c.Param("mediaItemId"), 10, 64)
// 	if err != nil {
// 		responses.RespondBadRequest(c, err, "Invalid media item ID")
// 		return
// 	}
//
// 	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
// 	if err != nil {
// 		responses.RespondBadRequest(c, err, "Invalid user ID")
// 		return
// 	}
//
// 	favorite, err := strconv.ParseBool(c.Query("favorite"))
// 	if err != nil {
// 		responses.RespondBadRequest(c, err, "Invalid favorite value")
// 		return
// 	}
//
// 	err = h.service.ToggleFavorite(c.Request.Context(), mediaItemID, userID, favorite)
// 	if err != nil {
// 		responses.RespondInternalError(c, err, "Failed to update favorite status")
// 		return
// 	}
//
// 	responses.RespondOK(c, responses.EmptyResponse{Success: true}, "Favorite status updated successfully")
// }
//
// // UpdateUserRating godoc
// // @Summary Update user rating for a media item
// // @Description Set a user's rating for a media item
// // @Tags History
// // @Accept json
// // @Produce json
// // @Param mediaItemId path int true "Media Item ID"
// // @Param userId query int true "User ID"
// // @Param rating query number true "User rating (0-10)"
// // @Success 200 {object} responses.APIResponse[models.MediaPlayHistory[any]] "Rating updated successfully"
// // @Failure 400 {object} responses.ErrorResponse[any] "Bad request"
// // @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// // @Failure 500 {object} responses.ErrorResponse[any] "Internal server error"
// // @Router /history/media/{mediaItemId}/rating [put]
// func (h *UserMediaItemDataHandler) UpdateUserRating(c *gin.Context) {
// 	mediaItemID, err := strconv.ParseUint(c.Param("mediaItemId"), 10, 64)
// 	if err != nil {
// 		responses.RespondBadRequest(c, err, "Invalid media item ID")
// 		return
// 	}
//
// 	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
// 	if err != nil {
// 		responses.RespondBadRequest(c, err, "Invalid user ID")
// 		return
// 	}
//
// 	rating, err := strconv.ParseFloat(c.Query("rating"), 32)
// 	if err != nil {
// 		responses.RespondBadRequest(c, err, "Invalid rating value")
// 		return
// 	}
//
// 	if rating < 0 || rating > 10 {
// 		responses.RespondBadRequest(c, nil, "Rating must be between 0 and 10")
// 		return
// 	}
//
// 	err = h.service.UpdateRating(c.Request.Context(), mediaItemID, userID, float32(rating))
// 	if err != nil {
// 		responses.RespondInternalError(c, err, "Failed to update rating")
// 		return
// 	}
//
// 	responses.RespondOK(c, responses.EmptyResponse{Success: true}, "Rating updated successfully")
// }
//
// // GetFavorites godoc
// // @Summary Get a user's favorite media items
// // @Description Get all media items marked as favorites by a user
// // @Tags History
// // @Accept json
// // @Produce json
// // @Param userId query int true "User ID"
// // @Param type query string false "Media type filter (movie, series, episode, track, etc.)"
// // @Param limit query int false "Number of items to return (default 10)"
// // @Param offset query int false "Number of items to skip (default 0)"
// // @Success 200 {object} responses.APIResponse[[]models.MediaItem[any]] "Successfully retrieved favorites"
// // @Failure 400 {object} responses.ErrorResponse[any] "Bad request"
// // @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// // @Failure 500 {object} responses.ErrorResponse[any] "Internal server error"
// // @Router /favorites [get]
// func (h *UserMediaItemDataHandler) GetFavorites(c *gin.Context) {
// 	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
// 	if err != nil {
// 		responses.RespondBadRequest(c, err, "Invalid user ID")
// 		return
// 	}
//
// 	mediaType := c.Query("type")
// 	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
// 	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
//
// 	// Filter by media type if provided
// 	var typedMediaType *types.MediaType
// 	if mediaType != "" {
// 		mt := types.MediaType(mediaType)
// 		typedMediaType = &mt
// 	}
//
// 	favorites, err := h.service.GetFavorites(c.Request.Context(), userID, typedMediaType, limit, offset)
// 	if err != nil {
// 		responses.RespondInternalError(c, err, "Failed to retrieve favorites")
// 		return
// 	}
//
// 	responses.RespondOK(c, favorites, "Favorites retrieved successfully")
// }
//
// // DeleteHistory godoc
// // @Summary Delete a play history entry
// // @Description Delete a specific play history entry by ID
// // @Tags History
// // @Accept json
// // @Produce json
// // @Param id path int true "History ID"
// // @Success 200 {object} responses.APIResponse[models.MediaPlayHistory[any]] "History entry deleted successfully"
// // @Failure 400 {object} responses.ErrorResponse[any] "Bad request"
// // @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// // @Failure 404 {object} responses.ErrorResponse[any] "Not found"
// // @Failure 500 {object} responses.ErrorResponse[any] "Internal server error"
// // @Router /history/{id} [delete]
// func (h *UserMediaItemDataHandler) DeleteHistory(c *gin.Context) {
// 	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
// 	if err != nil {
// 		responses.RespondBadRequest(c, err, "Invalid history ID")
// 		return
// 	}
//
// 	err = h.service.Delete(c.Request.Context(), id)
// 	if err != nil {
// 		responses.RespondInternalError(c, err, "Failed to delete history entry")
// 		return
// 	}
//
// 	responses.RespondOK(c, responses.EmptyResponse{Success: true}, "History entry deleted successfully")
// }
//
// // ClearUserHistory godoc
// // @Summary Clear a user's play history
// // @Description Delete all play history entries for a user
// // @Tags History
// // @Accept json
// // @Produce json
// // @Param userId query int true "User ID"
// // @Param type query string false "Media type filter (movie, series, episode, track, etc.)"
// // @Success 200 {object} responses.APIResponse[any] "History cleared successfully"
// // @Failure 400 {object} responses.ErrorResponse[any] "Bad request"
// // @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// // @Failure 500 {object} responses.ErrorResponse[any] "Internal server error"
// // @Router /history/clear [delete]
// func (h *UserMediaItemDataHandler) ClearUserHistory(c *gin.Context) {
// 	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
// 	if err != nil {
// 		responses.RespondBadRequest(c, err, "Invalid user ID")
// 		return
// 	}
//
// 	mediaType := c.Query("type")
//
// 	// Filter by media type if provided
// 	var typedMediaType *types.MediaType
// 	if mediaType != "" {
// 		mt := types.MediaType(mediaType)
// 		typedMediaType = &mt
// 	}
//
// 	err = h.service.ClearUserHistory(c.Request.Context(), userID, typedMediaType)
// 	if err != nil {
// 		responses.RespondInternalError(c, err, "Failed to clear history")
// 		return
// 	}
//
// 	responses.RespondOK(c, responses.EmptyResponse{Success: true}, "History cleared successfully")
// }
