package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"suasor/client/media/types"
	clientTypes "suasor/client/types"
	"suasor/services"
	"suasor/types/models"
	"suasor/types/responses"
)

// TODO: Clean this up and merge the duplicated code. There are 2 of every method a local. Also need to add in godoc api documentation like we have done with our other handlers.
// TODO: Add in logging like we have done with our other handlers.

// MediaItemHandler handles all media item operations
type MediaItemHandler[T types.MediaData] struct {
	service services.MediaItemService[T]
}

// NewMediaItemHandler creates a new media item handler
func NewMediaItemHandler[T types.MediaData](service services.MediaItemService[T]) *MediaItemHandler[T] {
	return &MediaItemHandler[T]{service: service}
}

// CreateMediaItem handles creating any type of media item
func (h *MediaItemHandler[T]) CreateMediaItem(c *gin.Context) {
	var req struct {
		Type       string          `json:"type" binding:"required"`
		ClientID   uint64          `json:"clientId" binding:"required"`
		ClientType string          `json:"clientType" binding:"required"`
		ExternalID string          `json:"externalId" binding:"required"`
		Data       json.RawMessage `json:"data" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mediaType := types.MediaType(req.Type)

	h.createMediaItem(c, mediaType, req.ClientID, clientTypes.MediaClientType(req.ClientType), req.ExternalID, req.Data)
}

// Type-specific media item creation helper
func (h *MediaItemHandler[T]) createMediaItem(
	c *gin.Context,
	mediaType types.MediaType,
	clientID uint64,
	clientType clientTypes.MediaClientType,
	externalID string,
	data json.RawMessage,
) {
	var mediaData T
	if err := json.Unmarshal(data, &mediaData); err != nil {
		responses.RespondBadRequest(c, err, "Invalid media data format")
		return
	}

	mediaItem := models.MediaItem[T]{
		Type:        mediaType,
		ClientIDs:   []models.ClientID{},
		ExternalIDs: []models.ExternalID{},
		Data:        mediaData,
	}

	// Set client info
	mediaItem.SetClientInfo(clientID, clientType, externalID)

	// Only add external ID if provided
	if externalID != "" {
		mediaItem.AddExternalID("client", externalID)
	}

	result, err := h.service.Create(c.Request.Context(), mediaItem)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to create media item")
		return
	}

	responses.RespondCreated(c, result, "Media item created successfully")
}

// UpdateMediaItem handles updating any type of media item
func (h *MediaItemHandler[T]) UpdateMediaItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		responses.RespondBadRequest(c, err, "Invalid media item ID")

		return
	}

	var req struct {
		Type       string          `json:"type" binding:"required"`
		ClientID   uint64          `json:"clientId" binding:"required"`
		ClientType string          `json:"clientType" binding:"required"`
		ExternalID string          `json:"externalId" binding:"required"`
		Data       json.RawMessage `json:"data" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		responses.RespondValidationError(c, err)
		return
	}

	mediaType := types.MediaType(req.Type)

	h.updateMediaItem(c, id, mediaType, req.ClientID, clientTypes.MediaClientType(req.ClientType), req.ExternalID, req.Data)
}

// Type-specific media item update helper
func (h *MediaItemHandler[T]) updateMediaItem(
	c *gin.Context,
	id uint64,
	mediaType types.MediaType,
	clientID uint64,
	clientType clientTypes.MediaClientType,
	externalID string,
	data json.RawMessage,
) {
	var mediaData T
	if err := json.Unmarshal(data, &mediaData); err != nil {
		responses.RespondBadRequest(c, err, "Invalid media data format")
		return
	}

	mediaItem := models.MediaItem[T]{
		ID:          id,
		Type:        mediaType,
		ClientIDs:   []models.ClientID{},
		ExternalIDs: []models.ExternalID{},
		Data:        mediaData,
	}

	// Set client info
	mediaItem.SetClientInfo(clientID, clientType, externalID)

	// Only add external ID if provided
	if externalID != "" {
		mediaItem.AddExternalID("client", externalID)
	}

	result, err := h.service.Update(c.Request.Context(), mediaItem)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to update media item")
		return
	}

	responses.RespondOK(c, result, "Media item updated successfully")
}

// GetMediaItem retrieves a media item by ID
func (h *MediaItemHandler[T]) GetMediaItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		responses.RespondBadRequest(c, err, "Invalid media item ID")
		return
	}

	h.getMediaItem(c, id)
}

// Type-specific media item retrieval helper
func (h *MediaItemHandler[T]) getMediaItem(c *gin.Context, id uint64) {
	item, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		responses.RespondNotFound(c, err, "Media item not found")
		return
	}

	responses.RespondOK(c, item, "Media item retrieved successfully")
}

// GetMediaItemsByClient retrieves all media items for a specific client
func (h *MediaItemHandler[T]) GetMediaItemsByClient(c *gin.Context) {
	clientID, err := strconv.ParseUint(c.Param("clientId"), 10, 64)
	if err != nil {
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	mediaType := c.Query("type") // Optional media type filter

	if mediaType != "" {
		h.getMediaItemsByClient(c, clientID)
	}

	// If no specific type, we'll need to collect all types and merge
	responses.RespondNotImplemented(c, err, "Fetching all media types not yet implemented")

}

// Type-specific media items by client retrieval helper
func (h *MediaItemHandler[T]) getMediaItemsByClient(c *gin.Context, clientID uint64) {
	items, err := h.service.GetByClientID(c.Request.Context(), clientID)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to retrieve media items by client")
		return
	}

	responses.RespondOK(c, items, "Media items retrieved successfully")
}

// DeleteMediaItem deletes a media item by ID
func (h *MediaItemHandler[T]) DeleteMediaItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		responses.RespondBadRequest(c, err, "Invalid media item ID")
		return
	}

	h.deleteMediaItem(c, id)
}

// Type-specific media item deletion helper
func (h *MediaItemHandler[T]) deleteMediaItem(c *gin.Context, id uint64) {

	err := h.service.Delete(c.Request.Context(), id)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to delete media item")
		return
	}

	responses.RespondOK(c, responses.EmptyResponse{Success: true}, "Media item deleted successfully")
}

// SearchMediaItems searches for media items by title
func (h *MediaItemHandler[T]) SearchMediaItems(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		responses.RespondBadRequest(c, nil, "Search query is required")
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		responses.RespondBadRequest(c, err, "Invalid user ID")

		return
	}

	mediaType := c.Query("type") // Optional media type filter

	if mediaType != "" {
		h.searchMediaItems(c, query, userID)
	}

	// If no specific type, search across all types (not implemented here)
	responses.RespondNotImplemented(c, err, "Searching across all media types not yet implemented")
}

// Type-specific media items search helper
func (h *MediaItemHandler[T]) searchMediaItems(c *gin.Context, query string, userID uint64) {

	items, err := h.service.SearchByTitle(c.Request.Context(), query, userID)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to search media items")
		return
	}

	responses.RespondOK(c, items, "Media items retrieved successfully")
}

// GetRecentMediaItems retrieves recently added media items
func (h *MediaItemHandler[T]) GetRecentMediaItems(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	mediaType := c.Query("type") // Optional media type filter

	if mediaType != "" {
		h.getRecentMediaItems(c, userID, limit)
	}

	// If no specific type, get recent items across all types (not implemented here)
	responses.RespondNotImplemented(c, err, "Fetching recent items across all media types not yet implemented")
}

// Type-specific recent media items retrieval helper
func (h *MediaItemHandler[T]) getRecentMediaItems(c *gin.Context, userID uint64, limit int) {
	items, err := h.service.GetRecentItems(c.Request.Context(), userID, limit)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to retrieve recent media items")
		return
	}

	responses.RespondOK(c, items, "Recent media items retrieved successfully")
}

// GetMediaItemByExternalSourceID retrieves a media item by external source ID

func (h *MediaItemHandler[T]) GetMediaItemByExternalSourceID(c *gin.Context) {

	// TODO: Implement being able to get a media item by external source ID
	// You should be able to get a media item by the external source ID that is stored in the database if it exists.

}
