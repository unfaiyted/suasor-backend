// handlers/user_lists.go
package handlers

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"suasor/utils"

	"suasor/clients/media/types"
	"suasor/services"
	"suasor/types/models"
	"suasor/types/requests"
	"suasor/types/responses"
	"suasor/utils/logger"
)

// Define list handler interface
type UserListHandler[T types.ListData] interface {
	CoreListHandler[T]

	// User needs to be authenticated to access these operations
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	AddItem(c *gin.Context)
	RemoveItem(c *gin.Context)
	RemoveItemAtPosition(c *gin.Context)
	ReorderItems(c *gin.Context)

	// User-specific operations
	GetFavorite(c *gin.Context)
	GetUserLists(c *gin.Context)

	// Sync local list with remote list
	Sync(c *gin.Context)
}

// userListHandler handles user-specific operations for lists
type userListHandler[T types.ListData] struct {
	CoreListHandler[T]

	itemService services.UserMediaItemService[T]
	listService services.UserListService[T]
	syncService services.ListSyncService[T]
}

// NewuserListHandler creates a new user list handler
func NewUserListHandler[T types.ListData](
	coreHandler CoreListHandler[T],
	itemService services.UserMediaItemService[T],
	listService services.UserListService[T],
	syncService services.ListSyncService[T],
) UserListHandler[T] {
	return &userListHandler[T]{
		CoreListHandler: coreHandler,
		itemService:     itemService,
		listService:     listService,
		syncService:     syncService,
	}
}

// GetUserLists godoc
//
//	@Summary		Get user's lists
//	@Description	Retrieves all lists owned by the authenticated user
//	@Tags			lists
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			limit	query		int																false	"Maximum number of lists to return (default 20)"
//	@Param			offset	query		int																false	"Offset for pagination (default 0)"
//	@Param			userID	query		uint64															false	"User ID"
//	@Success		200		{object}	responses.APIResponse[[]models.MediaItem[types.ListData]]	"Lists retrieved successfully"
//	@Failure		401		{object}	responses.ErrorResponse[any]									"Unauthorized"
//	@Failure		500		{object}	responses.ErrorResponse[any]									"Server error"
//	@Router			/{listType}/user [get]
func (h *userListHandler[T]) GetUserLists(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, ok := checkUserAccess(c)
	if !ok {
		return
	}

	limit := utils.GetLimit(c, 20, 100, true)
	offset := utils.GetOffset(c, 0)

	log.Debug().
		Uint64("userID", userID).
		Msg("Getting user lists")

	// Get user lists
	lists, err := h.listService.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to retrieve user lists")
		responses.RespondInternalError(c, err, "Failed to retrieve lists")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Int("count", len(lists)).
		Msg("User lists retrieved successfully")
	responses.RespondOK(c, lists, "Lists retrieved successfully")
}

// CreateList godoc
//
//	@Summary		Create a new list
//	@Description	Creates a new list for the authenticated user
//	@Tags			lists
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			list	body		models.MediaItem[types.ListData]										true	"List details"
//	@Success		201		{object}	responses.APIResponse[models.MediaItem[types.Playlist]]	"List created successfully"
//	@Failure		400		{object}	responses.ErrorResponse[any]									"Invalid request"
//	@Failure		401		{object}	responses.ErrorResponse[any]									"Unauthorized"
//	@Failure		500		{object}	responses.ErrorResponse[any]									"Server error"
//	@Router			/{listType}/ [post]
func (h *userListHandler[T]) Create(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := checkUserAccess(c)
	if !exists {
		return
	}

	// Parse request body
	var req models.MediaItem[T]
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid request body")
		responses.RespondBadRequest(c, err, "Invalid request body")
		return
	}

	log.Debug().
		Uint64("userID", userID).
		Str("name", req.Title).
		Interface("data", req.Data).
		Msg("Creating new list")

	// itemList := types.ItemList{
	// 	ItemCount: 0,
	// 	OwnerID:   userID,
	// 	Details:   details,
	// 	IsPublic:  req.IsPublic,
	// 	IsSmart:   req.IsSmart,
	// 	SmartCriteria: map[string]any{
	// 		"genre":    req.Genre,
	// 		"year":     req.Year,
	// 		"rating":   req.Rating,
	// 		"duration": req.Duration,
	// 	},
	// }
	req.UUID = uuid.New().String()

	// list := types.NewList[T](details, itemList)

	// list.AddListItem(types.ListItem{
	// 	ItemID:        0,
	// 	Position:      0,
	// 	LastChanged:   time.Now(),
	// 	ChangeHistory: []types.ChangeRecord{},
	// })
	// var zero T
	// mediaType := types.GetMediaTypeFromTypeName(zero)
	//
	// // Create media item
	// mediaItem := models.NewMediaItem[T](mediaType, list)
	// mediaItem.Title = list.GetTitle()
	// mediaItem.Type = mediaType
	//
	// Create list
	createdList, err := h.listService.Create(ctx, userID, &req)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to create list")
		responses.RespondInternalError(c, err, "Failed to create list")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("listID", createdList.ID).
		Msg("List created successfully")
	responses.RespondCreated(c, createdList, "List created successfully")
}

// UpdateList godoc
//
//	@Summary		Update a list
//	@Description	Updates an existing list owned by the authenticated user
//	@Tags			lists
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			listID		path		int																true	"List ID"
//	@Param			listType	path		string															true	"List type (e.g. 'playlist', 'collection')"
//	@Param			list		body		requests.ListUpdateRequest										true	"Updated list details"
//	@Success		200			{object}	responses.APIResponse[models.MediaItem[types.Playlist]]	"List updated successfully"
//	@Failure		400			{object}	responses.ErrorResponse[any]									"Invalid request"
//	@Failure		401			{object}	responses.ErrorResponse[any]									"Unauthorized"
//	@Failure		403			{object}	responses.ErrorResponse[any]									"Forbidden"
//	@Failure		404			{object}	responses.ErrorResponse[any]									"List not found"
//	@Failure		500			{object}	responses.ErrorResponse[any]									"Server error"
//	@Router			/{listType}/{listID} [put]
func (h *userListHandler[T]) Update(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, isAdmin := checkAdminAccess(c)

	// Parse list ID
	listID, _ := checkItemID(c, "listID")

	// Parse request body
	var req requests.ListUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Err(err).Msg("Invalid request body")
		responses.RespondBadRequest(c, err, "Invalid request body")
		return
	}

	log.Debug().
		Uint64("userID", userID).
		Bool("isAdmin", isAdmin).
		Uint64("listID", listID).
		Msg("Updating list")

	// Get existing list
	existingList, err := h.listService.GetByID(ctx, listID)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Msg("Failed to retrieve list")
		responses.RespondNotFound(c, err, "List not found")
		return
	}

	itemList := existingList.GetData().GetItemList()
	itemList.Details.Title = req.Name
	itemList.Details.Description = req.Description
	itemList.IsPublic = req.IsPublic
	existingList.GetData().SetItemList(*itemList)

	// Update list
	existingList.Title = req.Name

	// Save updated list
	updatedList, err := h.listService.Update(ctx, userID, existingList)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Msg("Failed to update list")
		responses.RespondInternalError(c, err, "Failed to update list")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("listID", listID).
		Msg("List updated successfully")
	responses.RespondOK(c, updatedList, "List updated successfully")
}

// DeleteList godoc
//
//	@Summary		Delete a list
//	@Description	Deletes a list owned by the authenticated user
//	@Tags			lists
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			listID	path		int								true	"List ID"
//	@Success		200	{object}	responses.APIResponse[any]		"List deleted successfully"
//	@Failure		400	{object}	responses.ErrorResponse[any]	"Invalid request"
//	@Failure		401	{object}	responses.ErrorResponse[any]	"Unauthorized"
//	@Failure		403	{object}	responses.ErrorResponse[any]	"Forbidden"
//	@Failure		404	{object}	responses.ErrorResponse[any]	"List not found"
//	@Failure		500	{object}	responses.ErrorResponse[any]	"Server error"
//	@Router			/{listType}/{listID} [delete]
func (h *userListHandler[T]) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, _ := checkUserAccess(c)

	// Parse list ID
	listID, _ := checkItemID(c, "listID")

	log.Debug().
		Uint64("userID", userID).
		Uint64("listID", listID).
		Msg("Deleting list")

	// Get existing list
	existingList, err := h.listService.GetByID(ctx, listID)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Msg("Failed to retrieve list")
		responses.RespondNotFound(c, err, "List not found")
		return
	}

	// Check if user owns the list
	if existingList.OwnerID != userID {
		log.Warn().
			Uint64("userID", userID).
			Uint64("ownerID", existingList.OwnerID).
			Uint64("listID", listID).
			Msg("User does not own the list")
		responses.RespondForbidden(c, nil, "You do not have permission to delete this list")
		return
	}

	// Delete list
	err = h.listService.Delete(ctx, userID, listID)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Msg("Failed to delete list")
		responses.RespondInternalError(c, err, "Failed to delete list")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("listID", listID).
		Msg("List deleted successfully")
	responses.RespondOK(c, http.StatusOK, "List deleted successfully")
}

// AddItemToList godoc
//
//	@Summary		Add a track to a list
//	@Description	Adds a track to a list owned by the authenticated user
//	@Tags			lists
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			listID	 path		int																true	"List ID"
//	@Param			itemID	 path		int																true	"Item ID"
//	@Success		200		{object}	responses.APIResponse[models.MediaItem[types.Playlist]]	"Track added successfully"
//	@Failure		400		{object}	responses.ErrorResponse[any]									"Invalid request"
//	@Failure		401		{object}	responses.ErrorResponse[any]									"Unauthorized"
//	@Failure		403		{object}	responses.ErrorResponse[any]									"Forbidden"
//	@Failure		404		{object}	responses.ErrorResponse[any]									"List not found"
//	@Failure		500		{object}	responses.ErrorResponse[any]									"Server error"
//	@Router			/{listType}/{listID}/add/{itemID} [post]
func (h *userListHandler[T]) AddItem(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, _ := checkUserAccess(c)

	// Parse list ID
	listID, _ := checkItemID(c, "listID")
	itemID, _ := checkItemID(c, "itemID")

	log.Debug().
		Uint64("userID", userID).
		Uint64("listID", listID).
		Uint64("itemID", itemID).
		Msg("Adding item to list")

	// Add track to list
	err := h.listService.AddItem(ctx, userID, listID, itemID)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Uint64("itemID", itemID).
			Msg("Failed to add track to list")
		responses.RespondInternalError(c, err, "Failed to add track to list")
		return
	}

	// Get updated list
	updatedList, err := h.listService.GetByID(ctx, listID)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Msg("Failed to retrieve updated list")
		responses.RespondInternalError(c, err, "Failed to retrieve updated list")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("listID", listID).
		Uint64("itemID", itemID).
		Msg("Track added to list successfully")
	responses.RespondOK(c, updatedList, "Track added to list successfully")
}

// RemoveItemFromList godoc
//
//	@Summary		Remove a item from a list
//	@Description	Removes a item from a list owned by the authenticated user
//	@Tags			lists
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			listID		path		int																true	"List ID"
//	@Param			itemID	path		int																true	"Track ID"
//	@Success		200		{object}	responses.APIResponse[models.MediaItem[types.ListData]]	"Track removed successfully"
//	@Failure		400		{object}	responses.ErrorResponse[any]									"Invalid request"
//	@Failure		401		{object}	responses.ErrorResponse[any]									"Unauthorized"
//	@Failure		403		{object}	responses.ErrorResponse[any]									"Forbidden"
//	@Failure		404		{object}	responses.ErrorResponse[any]									"List not found"
//	@Failure		500		{object}	responses.ErrorResponse[any]									"Server error"
//	@Router			/{listType}/{listID}/item/{itemID} [delete]
func (h *userListHandler[T]) RemoveItem(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, _ := checkUserAccess(c)

	// Parse list ID
	listID, _ := checkItemID(c, "listID")
	itemID, _ := checkItemID(c, "itemID")

	log.Debug().
		Uint64("userID", userID).
		Uint64("listID", listID).
		Uint64("trackID", itemID).
		Msg("Removing track from list")

	// Get existing list
	existingList, err := h.listService.GetByID(ctx, listID)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Msg("Failed to retrieve list")
		responses.RespondNotFound(c, err, "List not found")
		return
	}

	// Check if user owns the list
	if existingList.OwnerID != userID {
		log.Warn().
			Uint64("userID", userID).
			Uint64("ownerID", existingList.OwnerID).
			Uint64("listID", listID).
			Msg("User does not own the list")
		responses.RespondForbidden(c, nil, "You do not have permission to modify this list")
		return
	}

	// Remove track from list
	err = h.listService.RemoveItem(ctx, userID, listID, itemID)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Uint64("itemID", itemID).
			Msg("Failed to remove track from list")
		responses.RespondInternalError(c, err, "Failed to remove track from list")
		return
	}

	// Get updated list
	updatedList, err := h.listService.GetByID(ctx, listID)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Msg("Failed to retrieve updated list")
		responses.RespondInternalError(c, err, "Failed to retrieve updated list")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("listID", listID).
		Uint64("trackID", itemID).
		Msg("Track removed from list successfully")
	responses.RespondOK(c, updatedList, "Track removed from list successfully")
}

// RemoveItemAtPosition godoc
//
//	@Summary		Remove an item from a list at a specific position
//	@Description	Removes an item from a list owned by the authenticated user
//	@Tags			lists
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			listID		path		int																true	"List ID"
//	@Param			itemID	path		int																true	"Item ID"
//	@Param			position	path		int																true	"Position of item to remove"
//	@Success		200		{object}	responses.APIResponse[models.MediaItem[types.ListData]]	"Item removed successfully"
//	@Failure		400		{object}	responses.ErrorResponse[any]									"Invalid request"
//	@Failure		401		{object}	responses.ErrorResponse[any]									"Unauthorized"
//	@Failure		403		{object}	responses.ErrorResponse[any]									"Forbidden"
//	@Failure		404		{object}	responses.ErrorResponse[any]									"List not found"
//	@Failure		500		{object}	responses.ErrorResponse[any]									"Server error"
//	@Router			/{listType}/{listID}/item/{itemID}/position/{position} [delete]
func (h *userListHandler[T]) RemoveItemAtPosition(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, _ := checkUserAccess(c)

	// Parse list ID
	listID, _ := checkItemID(c, "listID")
	itemID, _ := checkItemID(c, "itemID")
	position, _ := checkItemID(c, "position")

	log.Debug().
		Uint64("userID", userID).
		Uint64("listID", listID).
		Uint64("itemID", itemID).
		Uint64("position", position).
		Msg("Removing item from list")

	// Get existing list
	existingList, err := h.listService.GetByID(ctx, listID)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Msg("Failed to retrieve list")
		responses.RespondNotFound(c, err, "List not found")
		return
	}

	// Check if user owns the list
	if existingList.OwnerID != userID {
		log.Warn().
			Uint64("userID", userID).
			Uint64("ownerID", existingList.OwnerID).
			Uint64("listID", listID).
			Msg("User does not own the list")
		responses.RespondForbidden(c, nil, "You do not have permission to modify this list")
		return
	}

	// Remove item from list
	err = h.listService.RemoveItemAtPosition(ctx, userID, listID, itemID, int(position))
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Uint64("itemID", itemID).
			Msg("Failed to remove item from list")
		responses.RespondInternalError(c, err, "Failed to remove item from list")
		return
	}

	// Get updated list
	updatedList, err := h.listService.GetByID(ctx, listID)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Msg("Failed to retrieve updated list")
		responses.RespondInternalError(c, err, "Failed to retrieve updated list")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("listID", listID).
		Uint64("itemID", itemID).
		Msg("Item removed from list successfully")
	responses.RespondOK(c, updatedList, "Item removed from list successfully")
}

// Reorder godoc
//
//	@Summary		Reorder list items
//	@Description	Reorders the items in a list
//	@Tags			lists
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			listID		path		int																true	"List ID"
//	@Param			request		body		requests.ListReorderRequest										true	"Reorder request"
//	@Param			listType	path		string															true	"List type (e.g. 'playlist', 'collection')"
//	@Success		200			{object}	responses.APIResponse[models.MediaItem[types.ListData]]	"List reordered successfully"
//	@Failure		400			{object}	responses.ErrorResponse[any]									"Invalid request"
//	@Failure		401			{object}	responses.ErrorResponse[any]									"Unauthorized"
//	@Failure		403			{object}	responses.ErrorResponse[any]									"Forbidden"
//	@Failure		404			{object}	responses.ErrorResponse[any]									"List not found"
//	@Failure		500			{object}	responses.ErrorResponse[any]									"Server error"
//	@Router			/{listType}/{listID}/reorder [post]
func (h *userListHandler[T]) ReorderItems(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, _ := checkUserAccess(c)

	// Parse list ID
	id, err := strconv.ParseUint(c.Param("listID"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("listID")).Msg("Invalid list ID")
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
		Uint64("userID", userID).
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
	if existingList.OwnerID != userID {
		log.Warn().
			Uint64("userID", userID).
			Uint64("ownerID", existingList.OwnerID).
			Uint64("listID", id).
			Msg("User does not own the list")
		responses.RespondForbidden(c, nil, "You do not have permission to reorder this list")
		return
	}

	// Reorder list items
	err = h.listService.ReorderItems(ctx, userID, id, req.ItemIDs)
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
		Uint64("userID", userID).
		Uint64("listID", id).
		Msg("List reordered successfully")
	responses.RespondOK(c, updatedList, "List reordered successfully")
}

// GetFavorites godoc
//
//	@Summary		Get favorites
//	@Description	Retrieves the favorites for the authenticated user
//	@Tags			lists
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			limit		query		int																false	"Maximum number of lists to return (default 20)"
//	@Param			offset		query		int																false	"Offset for pagination (default 0)"
//	@Param			userID		path		uint64															false	"User ID"
//	@Param			listType	path		string															true	"List type (e.g. 'playlist', 'collection')"
//	@Success		200			{object}	responses.APIResponse[[]models.MediaItem[types.ListData]]	"Lists retrieved successfully"
//	@Failure		401			{object}	responses.ErrorResponse[any]									"Unauthorized"
//	@Failure		500			{object}	responses.ErrorResponse[any]									"Server error"
//	@Router			/{listType}/favorites [get]
func (h *userListHandler[T]) GetFavorite(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, _ := checkUserAccess(c)
	limit := utils.GetLimit(c, 20, 100, true)
	offset := utils.GetOffset(c, 0)

	log.Debug().
		Uint64("userID", userID).
		Msg("Retrieving favorites")

	// Get favorites
	favorites, err := h.listService.GetFavorite(ctx, userID, limit, offset)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to retrieve favorites")
		responses.RespondInternalError(c, err, "Failed to retrieve favorites")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Int("count", len(favorites)).
		Msg("Favorites retrieved successfully")
	responses.RespondOK(c, favorites, "Favorites retrieved successfully")
}

// Sync godoc
//
// @Summary		Sync local list with remote list
// @Description	Synchronizes the local list with the remote list
// @Tags			lists
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			listID		path		int																true	"List ID"
// @Param			listType	path		string															true	"List type (e.g. 'playlist', 'collection')"
// @Param			clientID	path		int																true	"Client ID"
// @Success		200			{object}	responses.SuccessResponse									"List synced successfully"
// @Failure		400			{object}	responses.ErrorResponse[any]									"Invalid request"
// @Failure		401			{object}	responses.ErrorResponse[any]									"Unauthorized"
// @Failure		403			{object}	responses.ErrorResponse[any]									"Forbidden"
// @Failure		404			{object}	responses.ErrorResponse[any]									"List not found"
// @Failure		500			{object}	responses.ErrorResponse[any]									"Server error"
// @Router			/{listType}/{listID}/sync/{clientID} [post]
func (h *userListHandler[T]) Sync(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, _ := checkUserAccess(c)
	if userID == 0 {
		return
	}

	// Parse list ID and client ID
	listID, _ := checkItemID(c, "listID")
	clientID, _ := checkItemID(c, "clientID")

	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Uint64("listID", listID).
		Msg("Syncing list to client")

	// Verify list exists and user has access
	list, err := h.listService.GetByID(ctx, listID)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Msg("Failed to retrieve list")
		responses.RespondNotFound(c, err, "List not found")
		return
	}

	// Check if user owns the list
	if list.OwnerID != userID {
		log.Warn().
			Uint64("userID", userID).
			Uint64("ownerID", list.OwnerID).
			Uint64("listID", listID).
			Msg("User does not own the list")
		responses.RespondForbidden(c, nil, "You do not have permission to sync this list")
		return
	}

	// Use the sync service to sync the list to the client
	err = h.syncService.SyncToClient(ctx, userID, listID, clientID)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Uint64("clientID", clientID).
			Msg("Failed to sync list to client")
		responses.RespondInternalError(c, err, fmt.Sprintf("Failed to sync list: %s", err.Error()))
		return
	}

	// Return success response
	log.Info().
		Uint64("userID", userID).
		Uint64("listID", listID).
		Uint64("clientID", clientID).
		Msg("List synced successfully")
	responses.RespondOK(c, nil, "List synced successfully")
}
