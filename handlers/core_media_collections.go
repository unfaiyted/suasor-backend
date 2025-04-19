// handlers/core_collections.go
package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	mediatypes "suasor/client/media/types"
	"suasor/services"
	// "suasor/types/models"
	"suasor/types/responses"
	"suasor/utils"
)

// CoreCollectionHandler handles operations for collections in the database
type CoreCollectionHandler struct {
	collectionService services.CoreMediaItemService[*mediatypes.Collection]
	coreService       services.CoreListService[*mediatypes.Collection]
}

// NewCoreCollectionHandler creates a new core collection handler
func NewCoreCollectionHandler(
	collectionService services.CoreMediaItemService[*mediatypes.Collection],
	coreService services.CoreListService[*mediatypes.Collection],
) *CoreCollectionHandler {
	return &CoreCollectionHandler{
		collectionService: collectionService,
		coreService:       coreService,
	}
}

// GetAll godoc
// @Summary Get all collections
// @Description Retrieves all collections in the database
// @Tags collections
// @Accept json
// @Produce json
// @Param limit query int false "Maximum number of collections to return (default 10)"
// @Param offset query int false "Offset for pagination (default 0)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[*mediatypes.Collection]] "Collections retrieved successfully"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /collections [get]
func (h *CoreCollectionHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	log.Debug().Msg("Getting all collections")

	// Get user ID from query if provided
	limitStr := c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}
	offsetStr := c.Query("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	// Get collections
	collections, err := h.coreService.GetAll(ctx, limit, offset)
	if err != nil {
		log.Error().Err(err).
			Msg("Failed to retrieve collections for user")
		responses.RespondInternalError(c, err, "Failed to retrieve collections")
		return
	}

	log.Info().
		Int("count", len(collections)).
		Msg("Collections retrieved successfully")
	responses.RespondOK(c, collections, "Collections retrieved successfully")
	return
}

// GetByID godoc
// @Summary Get collection by ID
// @Description Retrieves a specific collection by ID
// @Tags collections
// @Accept json
// @Produce json
// @Param id path int true "Collection ID"
// @Success 200 {object} responses.APIResponse[models.MediaItem[*mediatypes.Collection]] "Collection retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Collection not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /collections/{id} [get]
func (h *CoreCollectionHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid collection ID")
		responses.RespondBadRequest(c, err, "Invalid collection ID")
		return
	}

	log.Debug().
		Uint64("id", id).
		Msg("Getting collection by ID")

	collection, err := h.collectionService.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("id", id).
			Msg("Failed to retrieve collection")
		responses.RespondNotFound(c, err, "Collection not found")
		return
	}

	log.Info().
		Uint64("id", id).
		Msg("Collection retrieved successfully")
	responses.RespondOK(c, collection, "Collection retrieved successfully")
}

// GetCollectionItems godoc
// @Summary Get items in a collection
// @Description Retrieves all media items in a specific collection
// @Tags collections
// @Accept json
// @Produce json
// @Param id path int true "Collection ID"
// @Success 200 {object} responses.APIResponse[[]mediatypes.MediaItem] "Items retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Collection not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /collections/{id}/items [get]
func (h *CoreCollectionHandler) GetCollectionItems(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid collection ID")
		responses.RespondBadRequest(c, err, "Invalid collection ID")
		return
	}

	log.Debug().
		Uint64("id", id).
		Msg("Getting items for collection")

	items, err := h.coreService.GetItems(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("id", id).
			Msg("Failed to retrieve collection items")
		responses.RespondInternalError(c, err, "Failed to retrieve collection items")
		return
	}

	log.Info().
		Uint64("id", id).
		Int("itemCount", items.GetTotalItems()).
		Msg("Collection items retrieved successfully")
	responses.RespondOK(c, items, "Collection items retrieved successfully")
}

// GetByGenre godoc
// @Summary Get collections by genre
// @Description Retrieves collections that match a specific genre
// @Tags collections
// @Accept json
// @Produce json
// @Param genre path string true "Genre name"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[*mediatypes.Collection]] "Collections retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /collections/genre/{genre} [get]
func (h *CoreCollectionHandler) GetByGenre(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	genre := c.Param("genre")
	if genre == "" {
		log.Warn().Msg("Genre is required")
		responses.RespondBadRequest(c, nil, "Genre is required")
		return
	}

	log.Debug().
		Str("genre", genre).
		Msg("Getting collections by genre")

	// Create query options
	options := mediatypes.QueryOptions{
		Genre:     genre,
		MediaType: mediatypes.MediaTypeCollection,
	}

	// Search collections by genre
	collections, err := h.collectionService.Search(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Str("genre", genre).
			Msg("Failed to retrieve collections by genre")
		responses.RespondInternalError(c, err, "Failed to retrieve collections")
		return
	}

	log.Info().
		Str("genre", genre).
		Int("count", len(collections)).
		Msg("Collections by genre retrieved successfully")
	responses.RespondOK(c, collections, "Collections retrieved successfully")
}

// Search godoc
// @Summary Search collections
// @Description Searches for collections that match the query
// @Tags collections
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[*mediatypes.Collection]] "Collections retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /collections/search [get]
func (h *CoreCollectionHandler) Search(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	query := c.Query("q")
	if query == "" {
		log.Warn().Msg("Search query is required")
		responses.RespondBadRequest(c, nil, "Search query is required")
		return
	}

	log.Debug().
		Str("query", query).
		Msg("Searching collections")

	// Create query options
	options := mediatypes.QueryOptions{
		Query:     query,
		MediaType: mediatypes.MediaTypeCollection,
	}

	// Search collections
	collections, err := h.collectionService.Search(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Str("query", query).
			Msg("Failed to search collections")
		responses.RespondInternalError(c, err, "Failed to search collections")
		return
	}

	log.Info().
		Str("query", query).
		Int("count", len(collections)).
		Msg("Collections search completed successfully")
	responses.RespondOK(c, collections, "Collections retrieved successfully")
}

// GetPublicCollections godoc
// @Summary Get public collections
// @Description Retrieves collections that are marked as public
// @Tags collections
// @Accept json
// @Produce json
// @Param limit query int false "Maximum number of collections to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[*mediatypes.Collection]] "Collections retrieved successfully"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /collections/public [get]
func (h *CoreCollectionHandler) GetPublicCollections(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	log.Debug().
		Int("limit", limit).
		Msg("Getting public collections")

	// Create query options
	options := mediatypes.QueryOptions{
		MediaType: mediatypes.MediaTypeCollection,
		IsPublic:  true,
		Limit:     limit,
	}

	// Get public collections
	collections, err := h.collectionService.Search(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Msg("Failed to retrieve public collections")
		responses.RespondInternalError(c, err, "Failed to retrieve collections")
		return
	}

	log.Info().
		Int("count", len(collections)).
		Msg("Public collections retrieved successfully")
	responses.RespondOK(c, collections, "Collections retrieved successfully")
}
