// handlers/collections.go
package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	mediatypes "suasor/client/media/types"
	"suasor/services"
	"suasor/types/responses"
	"suasor/utils"
)

// CollectionHandler provides handlers for collection operations
type CollectionHandler struct {
	*MediaItemHandler[*mediatypes.Collection]
	service services.CollectionService
}

// NewCollectionHandler creates a new collection handler
func NewCollectionHandler(
	mediaItemService services.MediaItemService[*mediatypes.Collection],
	collectionService services.CollectionService,
) *CollectionHandler {
	return &CollectionHandler{
		MediaItemHandler: NewMediaItemHandler(mediaItemService),
		service:          collectionService,
	}
}

// GetCollections godoc
// @Summary Get all collections
// @Description Retrieves all collections for a user
// @Tags collections
// @Accept json
// @Produce json
// @Param userId query int true "User ID"
// @Param limit query int false "Maximum number of collections to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Collection]] "Collections retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /collections [get]
func (h *CollectionHandler) GetCollections(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	log.Debug().
		Uint64("userID", userID).
		Int("limit", limit).
		Msg("Getting collections")

	collections, err := h.service.GetCollections(ctx, userID, limit)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to retrieve collections")
		responses.RespondInternalError(c, err, "Failed to retrieve collections")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Int("count", len(collections)).
		Msg("Collections retrieved successfully")

	responses.RespondOK(c, collections, "Collections retrieved successfully")
}

// GetCollectionByID godoc
// @Summary Get a collection by ID
// @Description Retrieves a specific collection by ID
// @Tags collections
// @Accept json
// @Produce json
// @Param id path int true "Collection ID"
// @Param userId query int true "User ID"
// @Param clientId query int true "Client ID"
// @Success 200 {object} responses.APIResponse[models.MediaItem[mediatypes.Collection]] "Collection retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Collection not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /collections/{id} [get]
func (h *CollectionHandler) GetCollectionByID(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	collectionID, err := utils.ParseUint64(c.Param("id"))
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid collection ID")
		responses.RespondBadRequest(c, err, "Invalid collection ID")
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	clientID, err := strconv.ParseUint(c.Query("clientId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("clientId", c.Query("clientId")).Msg("Invalid client ID")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Uint64("collectionID", collectionID).
		Msg("Getting collection by ID")

	collection, err := h.service.GetCollectionByID(ctx, userID, clientID, collectionID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Uint64("collectionID", collectionID).
			Msg("Failed to retrieve collection")
		responses.RespondNotFound(c, err, "Collection not found")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Uint64("collectionID", collectionID).
		Msg("Collection retrieved successfully")

	responses.RespondOK(c, collection, "Collection retrieved successfully")
}

// GetCollectionItems godoc
// @Summary Get items in a collection
// @Description Retrieves all items in a specific collection
// @Tags collections
// @Accept json
// @Produce json
// @Param id path int true "Collection ID"
// @Param userId query int true "User ID"
// @Param clientId query int true "Client ID"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.MediaData]] "Collection items retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Collection not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /collections/{id}/items [get]
func (h *CollectionHandler) GetCollectionItems(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	collectionID, err := utils.ParseUint64(c.Param("id"))
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid collection ID")
		responses.RespondBadRequest(c, err, "Invalid collection ID")
		return
	}

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	clientID, err := strconv.ParseUint(c.Query("clientId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("clientId", c.Query("clientId")).Msg("Invalid client ID")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Uint64("collectionID", collectionID).
		Msg("Getting collection items")

	// First, verify the collection exists
	collection, err := h.service.GetCollectionByID(ctx, userID, clientID, collectionID)
	if err != nil || collection == nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Uint64("collectionID", collectionID).
			Msg("Failed to retrieve collection")
		responses.RespondNotFound(c, err, "Collection not found")
		return
	}

	// Get the collection ID as uint64
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Uint64("collectionID", collectionID).
			Msg("Failed to parse collection ID")
		responses.RespondBadRequest(c, err, "Invalid collection ID format")
		return
	}

	// Get items from the collection using the service
	items, err := h.service.GetCollectionItems(ctx, collectionID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Uint64("collectionID", collectionID).
			Msg("Failed to retrieve collection items")
		responses.RespondInternalError(c, err, "Failed to retrieve collection items")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Uint64("collectionID", collectionID).
		// Int("itemCount", len(items)).
		Msg("Collection items retrieved successfully")

	responses.RespondOK(c, items, "Collection items retrieved successfully")
}

// GetSmartCollections godoc
// @Summary Get smart collections
// @Description Retrieves smart collections that are dynamically generated based on rules
// @Tags collections
// @Accept json
// @Produce json
// @Param userId query int true "User ID"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Collection]] "Smart collections retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /collections/smart [get]
func (h *CollectionHandler) GetSmartCollections(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	log.Debug().
		Uint64("userID", userID).
		Msg("Getting smart collections")

	// Smart collections implementation would be more complex and likely
	// require additional business logic beyond the scope of this handler.
	// For now, we return a not implemented response.
	log.Info().
		Uint64("userID", userID).
		Msg("Smart collections feature not yet implemented")
	responses.RespondNotImplemented(c, nil, "Smart collections feature not yet implemented")
}

// GetFeaturedCollections godoc
// @Summary Get featured collections
// @Description Retrieves featured collections recommended to the user
// @Tags collections
// @Accept json
// @Produce json
// @Param userId query int true "User ID"
// @Param limit query int false "Maximum number of collections to return (default 5)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Collection]] "Featured collections retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /collections/featured [get]
func (h *CollectionHandler) GetFeaturedCollections(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	userID, err := strconv.ParseUint(c.Query("userId"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("userId", c.Query("userId")).Msg("Invalid user ID")
		responses.RespondBadRequest(c, err, "Invalid user ID")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if err != nil {
		limit = 5
	}

	log.Debug().
		Uint64("userID", userID).
		Int("limit", limit).
		Msg("Getting featured collections")

	// Featured collections implementation would typically involve some additional
	// business logic to select collections to feature. For now, we'll use the
	// standard get collections with a smaller limit as a fallback.
	collections, err := h.service.GetCollections(ctx, userID, limit)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to retrieve featured collections")
		responses.RespondInternalError(c, err, "Failed to retrieve featured collections")
		return
	}

	log.Info().
		Uint64("userID", userID).
		Int("count", len(collections)).
		Msg("Featured collections retrieved successfully")

	responses.RespondOK(c, collections, "Featured collections retrieved successfully")
}

