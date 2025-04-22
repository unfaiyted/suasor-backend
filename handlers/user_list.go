// handlers/user_lists.go
package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	mediatypes "suasor/clients/media/types"
	"suasor/services"
	"suasor/types/models"
	"suasor/types/requests"
	"suasor/types/responses"
	"suasor/utils/logger"
)

// Define list handler interface
type UserListHandler[T mediatypes.ListData] interface {
	CoreListHandler[T]

	// User needs to be authenticated to access these operations
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	AddItem(c *gin.Context)
	RemoveItem(c *gin.Context)
	ReorderItems(c *gin.Context)

	// User-specific operations
	GetFavorite(c *gin.Context)
	GetUserLists(c *gin.Context)
}

// userListHandler handles user-specific operations for lists
type userListHandler[T mediatypes.ListData] struct {
	CoreListHandler[T]

	itemService services.UserMediaItemService[T]
	listService services.UserListService[T]
}

// NewuserListHandler creates a new user list handler
func NewUserListHandler[T mediatypes.ListData](
	coreHandler CoreListHandler[T],
	itemService services.UserMediaItemService[T],
	listService services.UserListService[T],
) UserListHandler[T] {
	return &userListHandler[T]{
		CoreListHandler: coreHandler,
		itemService:     itemService,
		listService:     listService,
	}
}

// GetUserLists godoc
// @Summary Get user's lists
// @Description Retrieves all lists owned by the authenticated user
// @Tags lists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Maximum number of lists to return (default 20)"
// @Param offset query int false "Offset for pagination (default 0)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[*mediatypes.List]] "Lists retrieved successfully"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /user/lists [get]
func (h *userListHandler[T]) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access lists without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 20
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	uid := userID.(uint64)

	log.Debug().
		Uint64("userID", uid).
		Msg("Getting user lists")

	// Get user lists
	lists, err := h.listService.GetByUserID(ctx, uid, limit, offset)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Msg("Failed to retrieve user lists")
		responses.RespondInternalError(c, err, "Failed to retrieve lists")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", len(lists)).
		Msg("User lists retrieved successfully")
	responses.RespondOK(c, lists, "Lists retrieved successfully")
}

// CreateList godoc
// @Summary Create a new list
// @Description Creates a new list for the authenticated user
// @Tags lists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param list body requests.ListCreateRequest true "List details"
// @Success 201 {object} responses.APIResponse[models.MediaItem[*mediatypes.List]] "List created successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /user/lists [post]
func (h *userListHandler[T]) Create(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to create list without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	// Parse request body
	var req requests.ListCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid request body")
		responses.RespondBadRequest(c, err, "Invalid request body")
		return
	}

	log.Debug().
		Uint64("userID", uid).
		Str("name", req.Name).
		Msg("Creating new list")

	itemList := mediatypes.ItemList{
		ItemCount: 0,
		OwnerID:   uid,
		Details: mediatypes.MediaDetails{
			Title:       req.Name,
			Description: req.Description,
		},
		IsPublic: req.IsPublic,
		IsSmart:  req.IsSmart,
		SmartCriteria: map[string]any{
			"genre":    req.Genre,
			"year":     req.Year,
			"rating":   req.Rating,
			"duration": req.Duration,
		},
	}

	list := mediatypes.NewList[T](mediatypes.MediaDetails{}, itemList)

	list.AddListItem(mediatypes.ListItem{
		ItemID:        0,
		Position:      0,
		LastChanged:   time.Now(),
		ChangeHistory: []mediatypes.ChangeRecord{},
	})
	var zero T
	mediaType := mediatypes.GetMediaTypeFromTypeName(zero)

	// Create media item
	mediaItem := models.NewMediaItem[T](mediaType, list)
	mediaItem.Title = list.GetTitle()
	mediaItem.Type = mediaType

	// Create list
	createdList, err := h.listService.Create(ctx, mediaItem)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Msg("Failed to create list")
		responses.RespondInternalError(c, err, "Failed to create list")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("listID", createdList.ID).
		Msg("List created successfully")
	responses.RespondCreated(c, createdList, "List created successfully")
}

// UpdateList godoc
// @Summary Update a list
// @Description Updates an existing list owned by the authenticated user
// @Tags lists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "List ID"
// @Param list body requests.ListUpdateRequest true "Updated list details"
// @Success 200 {object} responses.APIResponse[models.MediaItem[*mediatypes.List]] "List updated successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 403 {object} responses.ErrorResponse[any] "Forbidden"
// @Failure 404 {object} responses.ErrorResponse[any] "List not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /user/lists/{id} [put]
func (h *userListHandler[T]) Update(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to update list without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	// Parse list ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid list ID")
		responses.RespondBadRequest(c, err, "Invalid list ID")
		return
	}

	// Parse request body
	var req requests.ListUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid request body")
		responses.RespondBadRequest(c, err, "Invalid request body")
		return
	}

	log.Debug().
		Uint64("userID", uid).
		Uint64("listID", id).
		Msg("Updating list")

	// Get existing list
	existingList, err := h.listService.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", id).
			Msg("Failed to retrieve list")
		responses.RespondNotFound(c, err, "List not found")
		return
	}

	// Check if user owns the list
	if existingList.OwnerID != uid {
		log.Warn().
			Uint64("userID", uid).
			Uint64("ownerID", existingList.OwnerID).
			Uint64("listID", id).
			Msg("User does not own the list")
		responses.RespondForbidden(c, nil, "You do not have permission to update this list")
		return
	}

	itemList := existingList.GetData().GetItemList()
	itemList.Details.Title = req.Name
	itemList.Details.Description = req.Description
	itemList.IsPublic = req.IsPublic
	existingList.GetData().SetItemList(itemList)

	// Update list
	existingList.Title = req.Name

	// Save updated list
	updatedList, err := h.listService.Update(ctx, existingList)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", id).
			Msg("Failed to update list")
		responses.RespondInternalError(c, err, "Failed to update list")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("listID", id).
		Msg("List updated successfully")
	responses.RespondOK(c, updatedList, "List updated successfully")
}

// DeleteList godoc
// @Summary Delete a list
// @Description Deletes a list owned by the authenticated user
// @Tags lists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "List ID"
// @Success 200 {object} responses.APIResponse[any] "List deleted successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 403 {object} responses.ErrorResponse[any] "Forbidden"
// @Failure 404 {object} responses.ErrorResponse[any] "List not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /user/lists/{id} [delete]
func (h *userListHandler[T]) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to delete list without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	// Parse list ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid list ID")
		responses.RespondBadRequest(c, err, "Invalid list ID")
		return
	}

	log.Debug().
		Uint64("userID", uid).
		Uint64("listID", id).
		Msg("Deleting list")

	// Get existing list
	existingList, err := h.listService.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", id).
			Msg("Failed to retrieve list")
		responses.RespondNotFound(c, err, "List not found")
		return
	}

	// Check if user owns the list
	if existingList.OwnerID != uid {
		log.Warn().
			Uint64("userID", uid).
			Uint64("ownerID", existingList.OwnerID).
			Uint64("listID", id).
			Msg("User does not own the list")
		responses.RespondForbidden(c, nil, "You do not have permission to delete this list")
		return
	}

	// Delete list
	err = h.listService.Delete(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", id).
			Msg("Failed to delete list")
		responses.RespondInternalError(c, err, "Failed to delete list")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("listID", id).
		Msg("List deleted successfully")
	responses.RespondOK(c, http.StatusOK, "List deleted successfully")
}

// AddTrackToList godoc
// @Summary Add a track to a list
// @Description Adds a track to a list owned by the authenticated user
// @Tags lists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "List ID"
// @Param track body requests.ListAddTrackRequest true "Track details"
// @Success 200 {object} responses.APIResponse[models.MediaItem[*mediatypes.List]] "Track added successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 403 {object} responses.ErrorResponse[any] "Forbidden"
// @Failure 404 {object} responses.ErrorResponse[any] "List not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /user/lists/{id}/tracks [post]
func (h *userListHandler[T]) AddItem(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to add track to list without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	// Parse list ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid list ID")
		responses.RespondBadRequest(c, err, "Invalid list ID")
		return
	}

	// Parse request body
	var req requests.ListAddTrackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid request body")
		responses.RespondBadRequest(c, err, "Invalid request body")
		return
	}

	log.Debug().
		Uint64("userID", uid).
		Uint64("listID", id).
		Uint64("trackID", req.TrackID).
		Msg("Adding track to list")

	// Get existing list
	existingList, err := h.listService.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", id).
			Msg("Failed to retrieve list")
		responses.RespondNotFound(c, err, "List not found")
		return
	}

	// Check if user owns the list
	if existingList.OwnerID != uid {
		log.Warn().
			Uint64("userID", uid).
			Uint64("ownerID", existingList.OwnerID).
			Uint64("listID", id).
			Msg("User does not own the list")
		responses.RespondForbidden(c, nil, "You do not have permission to modify this list")
		return
	}

	// Add track to list
	err = h.listService.AddItem(ctx, id, req.TrackID)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", id).
			Uint64("trackID", req.TrackID).
			Msg("Failed to add track to list")
		responses.RespondInternalError(c, err, "Failed to add track to list")
		return
	}

	// Get updated list
	updatedList, err := h.listService.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", id).
			Msg("Failed to retrieve updated list")
		responses.RespondInternalError(c, err, "Failed to retrieve updated list")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("listID", id).
		Uint64("trackID", req.TrackID).
		Msg("Track added to list successfully")
	responses.RespondOK(c, updatedList, "Track added to list successfully")
}

// RemoveTrackFromList godoc
// @Summary Remove a track from a list
// @Description Removes a track from a list owned by the authenticated user
// @Tags lists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "List ID"
// @Param trackId path int true "Track ID"
// @Success 200 {object} responses.APIResponse[models.MediaItem[*mediatypes.List]] "Track removed successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 403 {object} responses.ErrorResponse[any] "Forbidden"
// @Failure 404 {object} responses.ErrorResponse[any] "List not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /user/lists/{id}/tracks/{trackId} [delete]
func (h *userListHandler[T]) RemoveItem(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to remove track from list without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	// Parse list ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid list ID")
		responses.RespondBadRequest(c, err, "Invalid list ID")
		return
	}

	// Parse track ID
	trackID, err := strconv.ParseUint(c.Param("trackId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("trackId", c.Param("trackId")).Msg("Invalid track ID")
		responses.RespondBadRequest(c, err, "Invalid track ID")
		return
	}

	log.Debug().
		Uint64("userID", uid).
		Uint64("listID", id).
		Uint64("trackID", trackID).
		Msg("Removing track from list")

	// Get existing list
	existingList, err := h.listService.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", id).
			Msg("Failed to retrieve list")
		responses.RespondNotFound(c, err, "List not found")
		return
	}

	// Check if user owns the list
	if existingList.OwnerID != uid {
		log.Warn().
			Uint64("userID", uid).
			Uint64("ownerID", existingList.OwnerID).
			Uint64("listID", id).
			Msg("User does not own the list")
		responses.RespondForbidden(c, nil, "You do not have permission to modify this list")
		return
	}

	// Remove track from list
	err = h.listService.RemoveItem(ctx, id, trackID)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", id).
			Uint64("trackID", trackID).
			Msg("Failed to remove track from list")
		responses.RespondInternalError(c, err, "Failed to remove track from list")
		return
	}

	// Get updated list
	updatedList, err := h.listService.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", id).
			Msg("Failed to retrieve updated list")
		responses.RespondInternalError(c, err, "Failed to retrieve updated list")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("listID", id).
		Uint64("trackID", trackID).
		Msg("Track removed from list successfully")
	responses.RespondOK(c, updatedList, "Track removed from list successfully")
}

// Reorder godoc
// @Summary Reorder list items
// @Description Reorders the items in a list
// @Tags lists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "List ID"
// @Param request body requests.ListReorderRequest true "Reorder request"
// @Success 200 {object} responses.APIResponse[models.MediaItem[*types.List]] "List reordered successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 403 {object} responses.ErrorResponse[any] "Forbidden"
// @Failure 404 {object} responses.ErrorResponse[any] "List not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /user/lists/{id}/reorder [post]
func (h *userListHandler[T]) ReorderItems(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to reorder list without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	// Parse list ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid list ID")
		responses.RespondBadRequest(c, err, "Invalid list ID")
		return
	}

	// Parse request body
	var req requests.ListReorderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid request body")
		responses.RespondBadRequest(c, err, "Invalid request body")
		return
	}

	log.Debug().
		Uint64("userID", uid).
		Uint64("listID", id).
		Msg("Reordering list items")

	// Get existing list
	existingList, err := h.listService.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", id).
			Msg("Failed to retrieve list")
		responses.RespondNotFound(c, err, "List not found")
		return
	}

	// Check if user owns the list
	if existingList.OwnerID != uid {
		log.Warn().
			Uint64("userID", uid).
			Uint64("ownerID", existingList.OwnerID).
			Uint64("listID", id).
			Msg("User does not own the list")
		responses.RespondForbidden(c, nil, "You do not have permission to reorder this list")
		return
	}

	// Reorder list items
	err = h.listService.ReorderItems(ctx, id, req.ItemIDs)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", id).
			Msg("Failed to reorder list items")
		responses.RespondInternalError(c, err, "Failed to reorder list items")
		return
	}

	// Get updated list
	updatedList, err := h.listService.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", id).
			Msg("Failed to retrieve updated list")
		responses.RespondInternalError(c, err, "Failed to retrieve updated list")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("listID", id).
		Msg("List reordered successfully")
	responses.RespondOK(c, updatedList, "List reordered successfully")
}

// GetFavorite godoc
// @Summary Get favorites
// @Description Retrieves the favorites for the authenticated user
// @Tags lists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Maximum number of lists to return (default 20)"
// @Param offset query int false "Offset for pagination (default 0)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[*mediatypes.List]] "Lists retrieved successfully"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /user/favorites [get]
func (h *userListHandler[T]) GetFavorite(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to retrieve favorites without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 20
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	uid := userID.(uint64)

	log.Debug().
		Uint64("userID", uid).
		Msg("Retrieving favorites")

	// Get favorites
	favorites, err := h.listService.GetFavorite(ctx, uid, limit, offset)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Msg("Failed to retrieve favorites")
		responses.RespondInternalError(c, err, "Failed to retrieve favorites")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", len(favorites)).
		Msg("Favorites retrieved successfully")
	responses.RespondOK(c, favorites, "Favorites retrieved successfully")
}

// GetUserLists godoc
// @Summary Get user's lists
// @Description Retrieves all lists owned by the authenticated user
// @Tags lists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Maximum number of lists to return (default 20)"
// @Param offset query int false "Offset for pagination (default 0)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[*mediatypes.List]] "Lists retrieved successfully"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /user/lists [get]
func (h *userListHandler[T]) GetUserLists(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to retrieve lists without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 20
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	uid := userID.(uint64)

	log.Debug().
		Uint64("userID", uid).
		Msg("Getting user lists")

	// Get lists
	lists, err := h.listService.GetByUserID(ctx, uid, limit, offset)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Msg("Failed to retrieve user lists")
		responses.RespondInternalError(c, err, "Failed to retrieve lists")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", len(lists)).
		Msg("Lists retrieved successfully")
	responses.RespondOK(c, lists, "Lists retrieved successfully")
}
