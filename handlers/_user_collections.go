// handlers/user_collections.go
package handlers

// import (
// 	"net/http"
// 	"strconv"
//
// 	"github.com/gin-gonic/gin"
//
// 	mediatypes "suasor/client/media/types"
// 	"suasor/services"
// 	"suasor/types/models"
// 	"suasor/types/requests"
// 	"suasor/types/responses"
// 	"suasor/utils"
// )
//
// // UserCollectionHandler handles user-specific operations for collections
// type UserCollectionHandler struct {
// 	userCollectionService services.UserMediaItemService[*mediatypes.Collection]
// 	collectionService     services.UserListService[*mediatypes.Collection]
// }
//
// // NewUserCollectionHandler creates a new user collection handler
// func NewUserCollectionHandler(
// 	userCollectionService services.UserMediaItemService[*mediatypes.Collection],
// 	collectionService services.UserListService[*mediatypes.Collection],
// ) *UserCollectionHandler {
// 	return &UserCollectionHandler{
// 		userCollectionService: userCollectionService,
// 		collectionService:     collectionService,
// 	}
// }
//
// // GetUserCollections godoc
// // @Summary Get user's collections
// // @Description Retrieves all collections owned by the authenticated user
// // @Tags collections
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param limit query int false "Maximum number of collections to return (default 10)"
// // @Param offset query int false "Offset for pagination (default 0)"
// // @Success 200 {object} responses.APIResponse[[]models.MediaItem[*mediatypes.Collection]] "Collections retrieved successfully"
// // @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// // @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// // @Router /user/collections [get]
// func (h *UserCollectionHandler) GetUser(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	log := utils.LoggerFromContext(ctx)
//
// 	// Get authenticated user ID
// 	userID, exists := c.Get("userID")
// 	if !exists {
// 		log.Warn().Msg("Attempt to access collections without authentication")
// 		responses.RespondUnauthorized(c, nil, "Authentication required")
// 		return
// 	}
// 	limitStr := c.DefaultQuery("limit", "10")
// 	offsetStr := c.DefaultQuery("offset", "0")
// 	limit, err := strconv.Atoi(limitStr)
// 	if err != nil {
// 		limit = 10
// 	}
// 	offset, err := strconv.Atoi(offsetStr)
// 	if err != nil {
// 		offset = 0
// 	}
//
// 	uid := userID.(uint64)
//
// 	log.Debug().
// 		Uint64("userID", uid).
// 		Msg("Getting user collections")
//
// 	// Get user collections
// 	collections, err := h.userCollectionService.GetByUserID(ctx, uid, limit, offset)
// 	if err != nil {
// 		log.Error().Err(err).
// 			Uint64("userID", uid).
// 			Msg("Failed to retrieve user collections")
// 		responses.RespondInternalError(c, err, "Failed to retrieve collections")
// 		return
// 	}
//
// 	log.Info().
// 		Uint64("userID", uid).
// 		Int("count", len(collections)).
// 		Msg("User collections retrieved successfully")
// 	responses.RespondOK(c, collections, "Collections retrieved successfully")
// }
//
// // CreateCollection godoc
// // @Summary Create a new collection
// // @Description Creates a new collection for the authenticated user
// // @Tags collections
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param collection body requests.CollectionCreateRequest true "Collection details"
// // @Success 201 {object} responses.APIResponse[models.MediaItem[*mediatypes.Collection]] "Collection created successfully"
// // @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// // @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// // @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// // @Router /user/collections [post]
// func (h *UserCollectionHandler) Create(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	log := utils.LoggerFromContext(ctx)
//
// 	// Get authenticated user ID
// 	userID, exists := c.Get("userID")
// 	if !exists {
// 		log.Warn().Msg("Attempt to create collection without authentication")
// 		responses.RespondUnauthorized(c, nil, "Authentication required")
// 		return
// 	}
//
// 	uid := userID.(uint64)
//
// 	// Parse request body
// 	var req requests.CollectionCreateRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		log.Warn().Err(err).Msg("Invalid request body")
// 		responses.RespondBadRequest(c, err, "Invalid request body")
// 		return
// 	}
//
// 	log.Debug().
// 		Uint64("userID", uid).
// 		Str("name", req.Name).
// 		Msg("Creating new collection")
//
// 	// Create collection data
// 	collection := &mediatypes.Collection{
// 		ItemList: mediatypes.ItemList{
// 			Details: mediatypes.MediaDetails{
// 				Title:       req.Name,
// 				Description: req.Description,
// 			},
// 			OwnerID:  uid,
// 			IsPublic: req.IsPublic,
// 			Items:    []mediatypes.ListItem{},
// 		},
// 	}
// 	// Create media item
// 	mediaItem := models.NewMediaItem(mediatypes.MediaTypeCollection, collection)
//
// 	// Create collection
// 	createdCollection, err := h.userCollectionService.Create(ctx, mediaItem)
// 	if err != nil {
// 		log.Error().Err(err).
// 			Uint64("userID", uid).
// 			Msg("Failed to create collection")
// 		responses.RespondInternalError(c, err, "Failed to create collection")
// 		return
// 	}
//
// 	log.Info().
// 		Uint64("userID", uid).
// 		Uint64("collectionID", createdCollection.ID).
// 		Msg("Collection created successfully")
// 	responses.RespondCreated(c, createdCollection, "Collection created successfully")
// }
//
// // UpdateCollection godoc
// // @Summary Update a collection
// // @Description Updates an existing collection owned by the authenticated user
// // @Tags collections
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int true "Collection ID"
// // @Param collection body requests.CollectionUpdateRequest true "Updated collection details"
// // @Success 200 {object} responses.APIResponse[models.MediaItem[*mediatypes.Collection]] "Collection updated successfully"
// // @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// // @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// // @Failure 403 {object} responses.ErrorResponse[any] "Forbidden"
// // @Failure 404 {object} responses.ErrorResponse[any] "Collection not found"
// // @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// // @Router /user/collections/{id} [put]
// func (h *UserCollectionHandler) Update(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	log := utils.LoggerFromContext(ctx)
//
// 	// Get authenticated user ID
// 	userID, exists := c.Get("userID")
// 	if !exists {
// 		log.Warn().Msg("Attempt to update collection without authentication")
// 		responses.RespondUnauthorized(c, nil, "Authentication required")
// 		return
// 	}
//
// 	uid := userID.(uint64)
//
// 	// Parse collection ID
// 	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
// 	if err != nil {
// 		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid collection ID")
// 		responses.RespondBadRequest(c, err, "Invalid collection ID")
// 		return
// 	}
//
// 	// Parse request body
// 	var req requests.CollectionUpdateRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		log.Warn().Err(err).Msg("Invalid request body")
// 		responses.RespondBadRequest(c, err, "Invalid request body")
// 		return
// 	}
//
// 	log.Debug().
// 		Uint64("userID", uid).
// 		Uint64("collectionID", id).
// 		Msg("Updating collection")
//
// 	// Get existing collection
// 	existingCollection, err := h.userCollectionService.GetByID(ctx, id)
// 	if err != nil {
// 		log.Error().Err(err).
// 			Uint64("collectionID", id).
// 			Msg("Failed to retrieve collection")
// 		responses.RespondNotFound(c, err, "Collection not found")
// 		return
// 	}
//
// 	// Check if user owns the collection
// 	if existingCollection.Data.OwnerID != uid {
// 		log.Warn().
// 			Uint64("userID", uid).
// 			Uint64("ownerID", existingCollection.Data.OwnerID).
// 			Uint64("collectionID", id).
// 			Msg("User does not own the collection")
// 		responses.RespondForbidden(c, nil, "You do not have permission to update this collection")
// 		return
// 	}
//
// 	// Update collection
// 	existingCollection.Data.Details.Title = req.Name
// 	existingCollection.Data.Details.Description = req.Description
// 	existingCollection.Type = mediatypes.MediaTypeCollection
// 	existingCollection.Data.IsPublic = req.IsPublic
//
// 	// Save updated collection
// 	updatedCollection, err := h.userCollectionService.Update(ctx, existingCollection)
// 	if err != nil {
// 		log.Error().Err(err).
// 			Uint64("collectionID", id).
// 			Msg("Failed to update collection")
// 		responses.RespondInternalError(c, err, "Failed to update collection")
// 		return
// 	}
//
// 	log.Info().
// 		Uint64("userID", uid).
// 		Uint64("collectionID", id).
// 		Msg("Collection updated successfully")
// 	responses.RespondOK(c, updatedCollection, "Collection updated successfully")
// }
//
// // DeleteCollection godoc
// // @Summary Delete a collection
// // @Description Deletes a collection owned by the authenticated user
// // @Tags collections
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int true "Collection ID"
// // @Success 200 {object} responses.APIResponse[any] "Collection deleted successfully"
// // @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// // @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// // @Failure 403 {object} responses.ErrorResponse[any] "Forbidden"
// // @Failure 404 {object} responses.ErrorResponse[any] "Collection not found"
// // @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// // @Router /user/collections/{id} [delete]
// func (h *UserCollectionHandler) Delete(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	log := utils.LoggerFromContext(ctx)
//
// 	// Get authenticated user ID
// 	userID, exists := c.Get("userID")
// 	if !exists {
// 		log.Warn().Msg("Attempt to delete collection without authentication")
// 		responses.RespondUnauthorized(c, nil, "Authentication required")
// 		return
// 	}
//
// 	uid := userID.(uint64)
//
// 	// Parse collection ID
// 	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
// 	if err != nil {
// 		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid collection ID")
// 		responses.RespondBadRequest(c, err, "Invalid collection ID")
// 		return
// 	}
//
// 	log.Debug().
// 		Uint64("userID", uid).
// 		Uint64("collectionID", id).
// 		Msg("Deleting collection")
//
// 	// Get existing collection
// 	existingCollection, err := h.userCollectionService.GetByID(ctx, id)
// 	if err != nil {
// 		log.Error().Err(err).
// 			Uint64("collectionID", id).
// 			Msg("Failed to retrieve collection")
// 		responses.RespondNotFound(c, err, "Collection not found")
// 		return
// 	}
//
// 	// Check if user owns the collection
// 	if existingCollection.Data.OwnerID != uid {
// 		log.Warn().
// 			Uint64("userID", uid).
// 			Uint64("ownerID", existingCollection.Data.OwnerID).
// 			Uint64("collectionID", id).
// 			Msg("User does not own the collection")
// 		responses.RespondForbidden(c, nil, "You do not have permission to delete this collection")
// 		return
// 	}
//
// 	// Delete collection
// 	err = h.userCollectionService.Delete(ctx, id)
// 	if err != nil {
// 		log.Error().Err(err).
// 			Uint64("collectionID", id).
// 			Msg("Failed to delete collection")
// 		responses.RespondInternalError(c, err, "Failed to delete collection")
// 		return
// 	}
//
// 	log.Info().
// 		Uint64("userID", uid).
// 		Uint64("collectionID", id).
// 		Msg("Collection deleted successfully")
// 	responses.RespondOK(c, http.StatusOK, "Collection deleted successfully")
// }
//
// // AddItemToCollection godoc
// // @Summary Add an item to a collection
// // @Description Adds a media item to a collection owned by the authenticated user
// // @Tags collections
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int true "Collection ID"
// // @Param item body requests.CollectionAddItemRequest true "Item details"
// // @Success 200 {object} responses.APIResponse[models.MediaItem[*mediatypes.Collection]] "Item added successfully"
// // @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// // @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// // @Failure 403 {object} responses.ErrorResponse[any] "Forbidden"
// // @Failure 404 {object} responses.ErrorResponse[any] "Collection not found"
// // @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// // @Router /user/collections/{id}/items [post]
// func (h *UserCollectionHandler) AddItem(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	log := utils.LoggerFromContext(ctx)
//
// 	// Get authenticated user ID
// 	userID, exists := c.Get("userID")
// 	if !exists {
// 		log.Warn().Msg("Attempt to add item to collection without authentication")
// 		responses.RespondUnauthorized(c, nil, "Authentication required")
// 		return
// 	}
//
// 	uid := userID.(uint64)
//
// 	// Parse collection ID
// 	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
// 	if err != nil {
// 		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid collection ID")
// 		responses.RespondBadRequest(c, err, "Invalid collection ID")
// 		return
// 	}
//
// 	// Parse request body
// 	var req requests.CollectionAddItemRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		log.Warn().Err(err).Msg("Invalid request body")
// 		responses.RespondBadRequest(c, err, "Invalid request body")
// 		return
// 	}
//
// 	log.Debug().
// 		Uint64("userID", uid).
// 		Uint64("collectionID", id).
// 		Uint64("itemID", req.ItemID).
// 		// Str("itemType", string(req.ItemType)).
// 		Msg("Adding item to collection")
//
// 	// Get existing collection
// 	existingCollection, err := h.userCollectionService.GetByID(ctx, id)
// 	if err != nil {
// 		log.Error().Err(err).
// 			Uint64("collectionID", id).
// 			Msg("Failed to retrieve collection")
// 		responses.RespondNotFound(c, err, "Collection not found")
// 		return
// 	}
//
// 	// Check if user owns the collection
// 	if existingCollection.Data.OwnerID != uid {
// 		log.Warn().
// 			Uint64("userID", uid).
// 			Uint64("ownerID", existingCollection.Data.OwnerID).
// 			Uint64("collectionID", id).
// 			Msg("User does not own the collection")
// 		responses.RespondForbidden(c, nil, "You do not have permission to modify this collection")
// 		return
// 	}
//
// 	// Add item to collection
// 	err = h.collectionService.AddItem(ctx, id, req.ItemID)
// 	if err != nil {
// 		log.Error().Err(err).
// 			Uint64("collectionID", id).
// 			Uint64("itemID", req.ItemID).
// 			Msg("Failed to add item to collection")
// 		responses.RespondInternalError(c, err, "Failed to add item to collection")
// 		return
// 	}
//
// 	// Get updated collection
// 	updatedCollection, err := h.userCollectionService.GetByID(ctx, id)
// 	if err != nil {
// 		log.Error().Err(err).
// 			Uint64("collectionID", id).
// 			Msg("Failed to retrieve updated collection")
// 		responses.RespondInternalError(c, err, "Failed to retrieve updated collection")
// 		return
// 	}
//
// 	log.Info().
// 		Uint64("userID", uid).
// 		Uint64("collectionID", id).
// 		Uint64("itemID", req.ItemID).
// 		Msg("Item added to collection successfully")
// 	responses.RespondOK(c, updatedCollection, "Item added to collection successfully")
// }
//
// // RemoveItemFromCollection godoc
// // @Summary Remove an item from a collection
// // @Description Removes a media item from a collection owned by the authenticated user
// // @Tags collections
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int true "Collection ID"
// // @Param itemId path int true "Item ID"
// // @Success 200 {object} responses.APIResponse[models.MediaItem[*mediatypes.Collection]] "Item removed successfully"
// // @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// // @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// // @Failure 403 {object} responses.ErrorResponse[any] "Forbidden"
// // @Failure 404 {object} responses.ErrorResponse[any] "Collection not found"
// // @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// // @Router /user/collections/{id}/items/{itemId} [delete]
// func (h *UserCollectionHandler) RemoveItem(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	log := utils.LoggerFromContext(ctx)
//
// 	// Get authenticated user ID
// 	userID, exists := c.Get("userID")
// 	if !exists {
// 		log.Warn().Msg("Attempt to remove item from collection without authentication")
// 		responses.RespondUnauthorized(c, nil, "Authentication required")
// 		return
// 	}
//
// 	uid := userID.(uint64)
//
// 	// Parse collection ID
// 	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
// 	if err != nil {
// 		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid collection ID")
// 		responses.RespondBadRequest(c, err, "Invalid collection ID")
// 		return
// 	}
//
// 	// Parse item ID
// 	itemID, err := strconv.ParseUint(c.Param("itemId"), 10, 64)
// 	if err != nil {
// 		log.Warn().Err(err).Str("itemId", c.Param("itemId")).Msg("Invalid item ID")
// 		responses.RespondBadRequest(c, err, "Invalid item ID")
// 		return
// 	}
//
// 	log.Debug().
// 		Uint64("userID", uid).
// 		Uint64("collectionID", id).
// 		Uint64("itemID", itemID).
// 		Msg("Removing item from collection")
//
// 	// Get existing collection
// 	existingCollection, err := h.userCollectionService.GetByID(ctx, id)
// 	if err != nil {
// 		log.Error().Err(err).
// 			Uint64("collectionID", id).
// 			Msg("Failed to retrieve collection")
// 		responses.RespondNotFound(c, err, "Collection not found")
// 		return
// 	}
//
// 	// Check if user owns the collection
// 	if existingCollection.Data.OwnerID != uid {
// 		log.Warn().
// 			Uint64("userID", uid).
// 			Uint64("ownerID", existingCollection.Data.OwnerID).
// 			Uint64("collectionID", id).
// 			Msg("User does not own the collection")
// 		responses.RespondForbidden(c, nil, "You do not have permission to modify this collection")
// 		return
// 	}
//
// 	// Remove item from collection
// 	err = h.collectionService.RemoveItem(ctx, id, itemID)
// 	if err != nil {
// 		log.Error().Err(err).
// 			Uint64("collectionID", id).
// 			Uint64("itemID", itemID).
// 			Msg("Failed to remove item from collection")
// 		responses.RespondInternalError(c, err, "Failed to remove item from collection")
// 		return
// 	}
//
// 	// Get updated collection
// 	updatedCollection, err := h.userCollectionService.GetByID(ctx, id)
// 	if err != nil {
// 		log.Error().Err(err).
// 			Uint64("collectionID", id).
// 			Msg("Failed to retrieve updated collection")
// 		responses.RespondInternalError(c, err, "Failed to retrieve updated collection")
// 		return
// 	}
//
// 	log.Info().
// 		Uint64("userID", uid).
// 		Uint64("collectionID", id).
// 		Uint64("itemID", itemID).
// 		Msg("Item removed from collection successfully")
// 	responses.RespondOK(c, updatedCollection, "Item removed from collection successfully")
// }
