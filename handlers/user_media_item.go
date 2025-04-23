// handlers/user_media_item.go
package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"net/http"

	mediatypes "suasor/clients/media/types"
	"suasor/services"
	models "suasor/types/models"
	"suasor/types/responses"
	"suasor/utils/logger"
)

type UserMediaItemHandler[T mediatypes.MediaData] interface {
	CoreMediaItemHandler[T]

	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByUserID(c *gin.Context)
	GetUserContent(c *gin.Context)
	SearchUserContent(c *gin.Context)
	GetRecentUserContent(c *gin.Context)
	CreateUserMediaItem(c *gin.Context)
	UpdateUserMediaItem(c *gin.Context)
	DeleteUserMediaItem(c *gin.Context)
}

// userMediaItemHandler handles operations for user-owned media items
// This handler extends CoreMediaItemHandler with operations specific to media items
// that are owned by users, such as playlists and collections
type userMediaItemHandler[T mediatypes.MediaData] struct {
	CoreMediaItemHandler[T] // Embed the core handler
	userService             services.UserMediaItemService[T]
}

// NewuserMediaItemHandler creates a new user media item handler
func NewUserMediaItemHandler[T mediatypes.MediaData](
	coreHandler CoreMediaItemHandler[T],
	userService services.UserMediaItemService[T],
) UserMediaItemHandler[T] {
	return &userMediaItemHandler[T]{
		CoreMediaItemHandler: coreHandler,
		userService:          userService,
	}
}

// GetByUserID godoc
// @Summary Get media items by user ID
// @Description Retrieves media items owned by a specific user
// @Tags user-media
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Param limit query int false "Maximum number of items to return (default 20)"
// @Param offset query int false "Offset for pagination (default 0)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Movie]] "User media items retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[responses.ErrorDetails] "User not found"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/user/{userID}/media [get]
func (h *userMediaItemHandler[T]) GetByUserID(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	userID, err := strconv.ParseUint(c.Param("userID"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userID", c.Param("userID")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	log.Debug().
		Uint64("userID", userID).
		Msg("Getting media items by user ID")

	// Get media items by user ID
	items, err := h.userService.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to retrieve user's media items")
		responses.RespondInternalError(c, err, "Failed to retrieve user's media items")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Int("count", len(items)).
		Msg("User's media items retrieved successfully")
	responses.RespondOK(c, items, "User's media items retrieved successfully")
}

// GetUserContent godoc
// @Summary Get all user-owned content
// @Description Retrieves all types of content owned by a user (playlists, collections, etc.)
// @Tags user-media
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Param limit query int false "Maximum number of items to return (default 20)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Movie]] "User content retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[responses.ErrorDetails] "User not found"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/user/{userID}/content [get]
func (h *userMediaItemHandler[T]) GetUserContent(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	userID, err := strconv.ParseUint(c.Param("userID"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userID", c.Param("userID")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		limit = 20
	}

	log.Debug().
		Uint64("userID", userID).
		Int("limit", limit).
		Msg("Getting all user-owned content")

	// Get all user-owned content
	items, err := h.userService.GetUserContent(ctx, userID, limit)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to retrieve user's content")
		responses.RespondInternalError(c, err, "Failed to retrieve user's content")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Int("count", len(items)).
		Msg("User's content retrieved successfully")
	responses.RespondOK(c, items, "User's content retrieved successfully")
}

// SearchUserContent godoc
// @Summary Search user-owned content
// @Description Searches for content owned by a user based on query parameters
// @Tags user-media
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Param q query string true "Search query"
// @Param type query string false "Media type filter"
// @Param limit query int false "Maximum number of items to return (default 20)"
// @Param offset query int false "Offset for pagination (default 0)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Movie]] "User content found successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[responses.ErrorDetails] "User not found"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/user/{userID}/content/search [get]
func (h *userMediaItemHandler[T]) SearchUserContent(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	userID, err := strconv.ParseUint(c.Param("userID"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userID", c.Param("userID")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	query := c.Query("q")
	if query == "" {
		log.Warn().Msg("Search query is required")
		responses.RespondBadRequest(c, nil, "Search query is required")
		return
	}

	// Get media type from query parameters
	mediaTypeStr := c.Query("type")
	var mediaType mediatypes.MediaType
	if mediaTypeStr != "" {
		mediaType = mediatypes.MediaType(mediaTypeStr)
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		limit = 20
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		offset = 0
	}

	log.Debug().
		Uint64("userID", userID).
		Str("query", query).
		Str("type", string(mediaType)).
		Int("limit", limit).
		Int("offset", offset).
		Msg("Searching user-owned content")

	// Create query options
	options := mediatypes.QueryOptions{
		Query:     query,
		MediaType: mediaType,
		OwnerID:   userID,
		Limit:     limit,
		Offset:    offset,
	}

	// Search user-owned content
	items, err := h.userService.SearchUserContent(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Str("query", query).
			Msg("Failed to search user's content")
		responses.RespondInternalError(c, err, "Failed to search user's content")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Str("query", query).
		Int("count", len(items)).
		Msg("User's content search completed successfully")
	responses.RespondOK(c, items, "User's content found successfully")
}

// GetRecentUserContent godoc
// @Summary Get recently created or updated user content
// @Description Retrieves recently created or updated content owned by a user
// @Tags user-media
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Param limit query int false "Maximum number of items to return (default 20)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Movie]] "Recent user content retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[responses.ErrorDetails] "User not found"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/user/{userID}/content/recent [get]
func (h *userMediaItemHandler[T]) GetRecentUserContent(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	userID, err := strconv.ParseUint(c.Param("userID"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userID", c.Param("userID")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		limit = 20
	}

	log.Debug().
		Uint64("userID", userID).
		Int("limit", limit).
		Msg("Getting recent user-owned content")

	// Get recent user content
	items, err := h.userService.GetRecentUserContent(ctx, userID, limit)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to retrieve recent user content")
		responses.RespondInternalError(c, err, "Failed to retrieve recent user content")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Int("count", len(items)).
		Msg("Recent user content retrieved successfully")
	responses.RespondOK(c, items, "Recent user content retrieved successfully")
}

// CreateUserMediaItem godoc
// @Summary Create a new user-owned media item
// @Description Creates a new media item owned by a user
// @Tags user-media
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Param mediaItem body requests.MediaItemCreateRequest true "Media item to create"
// @Success 201 {object} responses.APIResponse[models.MediaItem[mediatypes.Movie]] "Media item created successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[responses.ErrorDetails] "User not found"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/user/{userID}/media [post]
func (h *userMediaItemHandler[T]) CreateUserMediaItem(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	userID, err := strconv.ParseUint(c.Param("userID"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userID", c.Param("userID")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	var zero T
	mediaType := mediatypes.GetMediaTypeFromTypeName(zero)
	// Bind the request body to a media item struct

	mediaData := mediatypes.NewItem[T]()

	mediaItem := models.NewMediaItem(mediaType, mediaData)
	if err := c.ShouldBindJSON(&mediaItem); err != nil {
		log.Warn().Err(err).Msg("Invalid media item data")
		responses.RespondBadRequest(c, err, "Invalid media item data")
		return
	}

	// Ensure the media item is associated with the user
	// This will depend on your data structure, but generally for user-owned content
	// you'll need to set owner ID in the appropriate field within the data property
	// For example, if ItemList is the structure for playlists/collections:
	if &mediaItem.Data != nil {
		// Assuming your media data might have an ItemList property for collections/playlists
		// Check if we can set the owner field
		// Playlist and collections have an ItemList property
		// TODO move logic to playlist handling
		if itemList, ok := h.hasItemList(mediaItem.Data); ok {
			itemList.OwnerID = userID
		}
	}

	log.Debug().
		Uint64("userID", userID).
		Str("type", string(mediaItem.Type)).
		Msg("Creating user-owned media item")

	// Create the media item
	createdItem, err := h.userService.Create(ctx, mediaItem)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to create user-owned media item")
		responses.RespondInternalError(c, err, "Failed to create media item")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("itemID", createdItem.ID).
		Str("type", string(createdItem.Type)).
		Msg("User-owned media item created successfully")
	responses.RespondCreated(c, createdItem, "Media item created successfully")
}

// UpdateUserMediaItem godoc
// @Summary Update a user-owned media item
// @Description Updates an existing media item owned by a user
// @Tags user-media
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Param id path int true "Media Item ID"
// @Param mediaItem body requests.MediaItemUpdateRequest true "Updated media item data"
// @Success 200 {object} responses.APIResponse[models.MediaItem[mediatypes.MediaData]] "Media item updated successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[responses.ErrorDetails] "Media item not found"
// @Failure 403 {object} responses.ErrorResponse[responses.ErrorDetails] "Not authorized to update this media item"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/user/{userID}/media/{id} [put]
func (h *userMediaItemHandler[T]) UpdateUserMediaItem(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	userID, err := strconv.ParseUint(c.Param("userID"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userID", c.Param("userID")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	itemID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid media item ID")
		responses.RespondBadRequest(c, err, "Invalid media item ID")
		return
	}

	// Bind the request body to a media item struct
	var mediaItem models.MediaItem[T]
	if err := c.ShouldBindJSON(&mediaItem); err != nil {
		log.Warn().Err(err).Msg("Invalid media item data")
		responses.RespondBadRequest(c, err, "Invalid media item data")
		return
	}

	// Ensure the ID in the path matches the ID in the body
	mediaItem.ID = itemID

	// First, get the existing item to verify ownership
	existingItem, err := h.userService.GetByID(ctx, itemID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("itemID", itemID).
			Msg("Failed to retrieve existing media item")
		responses.RespondNotFound(c, err, "Media item not found")
		return
	}

	// Verify that the user owns this item
	// This will depend on your data structure
	if !h.isUserOwned(existingItem.Data, userID) {
		log.Warn().
			Uint64("userID", userID).
			Uint64("itemID", itemID).
			Msg("User not authorized to update this media item")
		responses.RespondForbidden(c, nil, "Not authorized to update this media item")
		return
	}

	// Ensure the item maintains the same owner
	if &mediaItem.Data != nil {
		// Assuming your media data might have an ItemList property for collections/playlists
		// Check if we can set the owner field
		if itemList, ok := h.hasItemList(mediaItem.Data); ok {
			itemList.OwnerID = userID
		}
	}

	log.Debug().
		Uint64("userID", userID).
		Uint64("itemID", itemID).
		Str("type", string(mediaItem.Type)).
		Msg("Updating user-owned media item")

	// Update the media item
	updatedItem, err := h.userService.Update(ctx, &mediaItem)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("itemID", itemID).
			Msg("Failed to update user-owned media item")
		responses.RespondInternalError(c, err, "Failed to update media item")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("itemID", updatedItem.ID).
		Str("type", string(updatedItem.Type)).
		Msg("User-owned media item updated successfully")
	responses.RespondOK(c, updatedItem, "Media item updated successfully")
}

// DeleteUserMediaItem godoc
// @Summary Delete a user-owned media item
// @Description Deletes a media item owned by a user
// @Tags user-media
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Param id path int true "Media Item ID"
// @Success 200 {object} responses.APIResponse[any] "Media item deleted successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[responses.ErrorDetails] "Media item not found"
// @Failure 403 {object} responses.ErrorResponse[responses.ErrorDetails] "Not authorized to delete this media item"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/user/{userID}/media/{id} [delete]
func (h *userMediaItemHandler[T]) DeleteUserMediaItem(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	userID, err := strconv.ParseUint(c.Param("userID"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userID", c.Param("userID")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	itemID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid media item ID")
		responses.RespondBadRequest(c, err, "Invalid media item ID")
		return
	}

	// First, get the existing item to verify ownership
	existingItem, err := h.userService.GetByID(ctx, itemID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("itemID", itemID).
			Msg("Failed to retrieve existing media item")
		responses.RespondNotFound(c, err, "Media item not found")
		return
	}

	// Verify that the user owns this item
	if !h.isUserOwned(existingItem.Data, userID) {
		log.Warn().
			Uint64("userID", userID).
			Uint64("itemID", itemID).
			Msg("User not authorized to delete this media item")
		responses.RespondForbidden(c, nil, "Not authorized to delete this media item")
		return
	}

	log.Debug().
		Uint64("userID", userID).
		Uint64("itemID", itemID).
		Str("type", string(existingItem.Type)).
		Msg("Deleting user-owned media item")

	// Delete the media item
	if err := h.userService.Delete(ctx, itemID); err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("itemID", itemID).
			Msg("Failed to delete user-owned media item")
		responses.RespondInternalError(c, err, "Failed to delete media item")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("itemID", itemID).
		Msg("User-owned media item deleted successfully")
	responses.RespondOK(c, http.StatusOK, "Media item deleted successfully")
}

// Helper function to check if a mediaData has an ItemList property
// and returns the ItemList for modification
func (h *userMediaItemHandler[T]) hasItemList(mediaData T) (*mediatypes.ItemList, bool) {
	// Implementation depends on your specific mediatypes.MediaData structure
	// This is just a placeholder - you'll need to implement based on your actual structure

	// For playlist type
	if playlist, ok := any(mediaData).(*mediatypes.Playlist); ok && playlist != nil {
		return &playlist.ItemList, true
	}

	// For collection type
	if collection, ok := any(mediaData).(*mediatypes.Collection); ok && collection != nil {
		return &collection.ItemList, true
	}

	return nil, false
}

// Helper function to check if a mediaData item is owned by a specific user
func (h *userMediaItemHandler[T]) isUserOwned(mediaData T, userID uint64) bool {
	// Implementation depends on your specific mediatypes.MediaData structure
	// This is just a placeholder - you'll need to implement based on your actual structure

	// Check for playlist ownership
	if playlist, ok := any(mediaData).(*mediatypes.Playlist); ok && playlist != nil {
		return playlist.ItemList.OwnerID == userID
	}

	// Check for collection ownership
	if collection, ok := any(mediaData).(*mediatypes.Collection); ok && collection != nil {
		return collection.ItemList.OwnerID == userID
	}

	return false
}

// Create godoc
// @Summary Create a new user-owned media item
// @Description Creates a new media item owned by a user
// @Tags UserMediaItem
// @Accept json
// @Produce json
// @Param mediaItem body requests.MediaItemCreateRequest true "Media item data with type, client info, and type-specific data"
// @Success 201 {object} responses.APIResponse[models.MediaItem[any]] "Media item created successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[responses.ErrorDetails] "Media item not found"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /user-media [post]
func (h *userMediaItemHandler[T]) Create(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	var zero T
	mediaType := mediatypes.GetMediaTypeFromTypeName(zero)
	// Bind the request body to a media item struct

	mediaData := mediatypes.NewItem[T]()
	mediaItem := models.NewMediaItem(mediaType, mediaData)
	if err := c.ShouldBindJSON(&mediaItem); err != nil {
		log.Warn().Err(err).Msg("Invalid media item data")
		responses.RespondBadRequest(c, err, "Invalid media item data")
		return
	}

	userIDStr, exists := c.Get("userID")
	userID, err := strconv.ParseUint(userIDStr.(string), 10, 64)
	if !exists {
		log.Warn().Msg("Attempt to create media item without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	// Ensure the media item is associated with the user
	// This will depend on your data structure, but generally for user-owned content
	// you'll need to set owner ID in the appropriate field within the data property
	// For example, if ItemList is the structure for playlists/collections:
	if &mediaItem.Data != nil {
		// Assuming your media data might have an ItemList property for collections/playlists
		// Check if we can set the owner field
		// Playlist and collections have an ItemList property
		// TODO move logic to playlist handling
		if itemList, ok := h.hasItemList(mediaItem.Data); ok {
			itemList.OwnerID = userID
		}
	}

	log.Debug().
		Uint64("userID", userID).
		Str("type", string(mediaItem.Type)).
		Msg("Creating user-owned media item")

	// Create the media item
	createdItem, err := h.userService.Create(ctx, mediaItem)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to create user-owned media item")
		responses.RespondInternalError(c, err, "Failed to create media item")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("id", createdItem.ID).
		Str("type", string(createdItem.Type)).
		Msg("User-owned media item created successfully")

	responses.RespondCreated(c, createdItem, "Media item created successfully")
}

// Update godoc
// @Summary Update an existing user-owned media item
// @Description Updates an existing media item owned by a user
// @Tags UserMediaItem
// @Accept json
// @Produce json
// @Param id path int true "Media Item ID"
// @Param mediaItem body requests.MediaItemUpdateRequest true "Updated media item data"
// @Success 200 {object} responses.APIResponse[models.MediaItem[any]] "Media item updated successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[responses.ErrorDetails] "Media item not found"
// @Failure 403 {object} responses.ErrorResponse[responses.ErrorDetails] "Not authorized to update this media item"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /user-media/{id} [put]
func (h *userMediaItemHandler[T]) Update(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid media item ID")
		responses.RespondBadRequest(c, err, "Invalid media item ID")
		return
	}

	// var zero T
	// mediaType := mediatypes.GetMediaTypeFromTypeName(zero)
	// Bind the request body to a media item struct
	var mediaItem models.MediaItem[T]
	if err := c.ShouldBindJSON(&mediaItem); err != nil {
		log.Warn().Err(err).Msg("Invalid media item data")
		responses.RespondBadRequest(c, err, "Invalid media item data")
		return
	}
	userIDStr, exists := c.Get("userID")
	userID, err := strconv.ParseUint(userIDStr.(string), 10, 64)
	if !exists {
		log.Warn().Msg("Attempt to create media item without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	// Ensure the ID in the path matches the ID in the body
	mediaItem.ID = id

	// First, get the existing item to verify ownership
	existingItem, err := h.userService.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("id", id).
			Msg("Failed to get existing media item")
		responses.RespondNotFound(c, err, "Media item not found")
		return
	}

	// Verify that the user owns this item
	// This will depend on your data structure
	if !h.isUserOwned(existingItem.Data, userID) {
		log.Warn().
			Uint64("userID", userID).
			Uint64("id", id).
			Msg("User not authorized to update this media item")
		responses.RespondForbidden(c, nil, "Not authorized to update this media item")
		return
	}

	// Ensure the item maintains the same owner
	if &mediaItem.Data != nil {
		// Assuming your media data might have an ItemList property for collections/playlists
		// Check if we can set the owner field
		// Playlist and collections have an ItemList property
		// TODO move logic to playlist handling
		if itemList, ok := h.hasItemList(mediaItem.Data); ok {
			itemList.OwnerID = userID
		}
	}

	log.Debug().
		Uint64("userID", userID).
		Uint64("id", id).
		Str("type", string(mediaItem.Type)).
		Msg("Updating user-owned media item")

	// Update the media item
	updatedItem, err := h.userService.Update(ctx, &mediaItem)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("id", id).
			Msg("Failed to update user-owned media item")
		responses.RespondInternalError(c, err, "Failed to update media item")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("id", updatedItem.ID).
		Str("type", string(updatedItem.Type)).
		Msg("User-owned media item updated successfully")

	responses.RespondOK(c, updatedItem, "Media item updated successfully")
}

// Delete godoc
// @Summary Delete a user-owned media item
// @Description Deletes a user-owned media item by its ID
// @Tags UserMediaItem
// @Accept json
// @Produce json
// @Param id path int true "User Media Item ID"
// @Success 200 {object} responses.APIResponse[any] "Successfully deleted user media item"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Bad request"
// @Failure 404 {object} responses.ErrorResponse[responses.ErrorDetails] "Not found"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Internal server error"
// @Router /user-media/{id} [delete]
func (h *userMediaItemHandler[T]) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid user media item ID")
		responses.RespondBadRequest(c, err, "Invalid user media item ID")
		return
	}

	// log.Debug.Uint64("id", id).Msg("Deleting user media item")

	err = h.userService.Delete(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("id", id).
			Msg("Failed to delete user media item")
		responses.RespondInternalError(c, err, "Failed to delete user media item")
		return
	}

	log.Info().
		Uint64("id", id).
		Msg("User media item deleted successfully")
	responses.RespondOK(c, responses.EmptyResponse{Success: true}, "User media item deleted successfully")
}