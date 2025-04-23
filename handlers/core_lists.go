// handlers/core_lists.go
package handlers

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"suasor/utils"

	mediatypes "suasor/clients/media/types"
	"suasor/services"
	_ "suasor/types/models"
	"suasor/types/responses"
	"suasor/utils/logger"
)

type CoreListHandler[T mediatypes.ListData] interface {
	CoreMediaItemHandler[T]

	GetAll(c *gin.Context)
	GetByID(c *gin.Context)
	GetItemsByListID(c *gin.Context)
	GetByGenre(c *gin.Context)
	Search(c *gin.Context)
}

// coreListHandler[T] handles operations for lists in the database
type coreListHandler[T mediatypes.ListData] struct {
	CoreMediaItemHandler[T]
	listService services.CoreListService[T]
}

// NewCoreListHandler[T] creates a new core playlist handler
func NewCoreListHandler[T mediatypes.ListData](
	CoreMediaItemHandler CoreMediaItemHandler[T],
	listService services.CoreListService[T],
) CoreListHandler[T] {
	return &coreListHandler[T]{
		CoreMediaItemHandler: CoreMediaItemHandler,
		listService:          listService,
	}
}

// GetAll godoc
// @Summary Get all lists
// @Description Retrieves all lists in the database
// @Tags lists
// @Accept json
// @Produce json
// @Param limit query int false "Maximum number of lists to return (default 10)"
// @Param offset query int false "Offset for pagination (default 0)"
// @Param listType path string true "List type (e.g. 'playlist', 'collection')"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Playlist]] "Lists retrieved successfully"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /api/v1/{listType} [get]
func (h *coreListHandler[T]) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Debug().Msg("Getting all lists")

	limit := utils.GetLimit(c, 10, 100, true)
	offset := utils.GetOffset(c, 0)

	// Get all lists
	lists, err := h.listService.GetAll(ctx, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve lists")
		responses.RespondInternalError(c, err, "Failed to retrieve lists")
		return
	}

	log.Info().
		Int("count", len(lists)).
		Msg("Lists retrieved successfully")
	responses.RespondOK(c, lists, "Lists retrieved successfully")
}

// GetByID godoc
// @Summary Get playlist by ID
// @Description Retrieves a specific playlist by ID
// @Tags lists
// @Accept json
// @Produce json
// @Param listId path int true "List ID"
// @Param listType path string true "List type (e.g. 'playlist', 'collection')"
// @Success 200 {object} responses.APIResponse[models.MediaItem[mediatypes.Playlist]] "List retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "List not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /api/v1/{listType}/{listId} [get]
func (h *coreListHandler[T]) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid playlist ID")
		responses.RespondBadRequest(c, err, "Invalid playlist ID")
		return
	}

	log.Debug().
		Uint64("id", id).
		Msg("Getting playlist by ID")

	playlist, err := h.listService.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("id", id).
			Msg("Failed to retrieve playlist")
		responses.RespondNotFound(c, err, "List not found")
		return
	}

	log.Info().
		Uint64("id", id).
		Msg("List retrieved successfully")
	responses.RespondOK(c, playlist, "List retrieved successfully")
}

// GetListItems godoc
// @Summary Get tracks in a playlist
// @Description Retrieves all tracks in a specific playlist
// @Tags lists
// @Accept json
// @Produce json
// @Param listId path int true "List ID"
// @Param listType path string true "List type (e.g. 'playlist', 'collection')"
// @Success 200 {object} responses.APIResponse[[]mediatypes.Track] "Tracks retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "List not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /api/v1/{listType}/{listId}/items [get]
func (h *coreListHandler[T]) GetItemsByListID(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid playlist ID")
		responses.RespondBadRequest(c, err, "Invalid playlist ID")
		return
	}

	log.Debug().
		Uint64("id", id).
		Msg("Getting tracks for playlist")

	playlist, err := h.listService.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("id", id).
			Msg("Failed to retrieve playlist")
		responses.RespondNotFound(c, err, "List not found")
		return
	}

	itemList := playlist.GetData().GetItemList()

	log.Info().
		Uint64("id", id).
		Int("itemCount", len(itemList.Items)).
		Msg("List tracks retrieved successfully")
	responses.RespondOK(c, itemList.Items, "Items retrieved successfully")
}

// GetByGenre godoc
// @Summary Get lists by genre
// @Description Retrieves lists that match a specific genre
// @Tags lists
// @Accept json
// @Produce json
// @Param genre path string true "Genre name"
// @Param listType path string true "List type (e.g. 'playlist', 'collection')"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Playlist]] "Lists retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /api/v1/{listType}/genre/{genre} [get]
func (h *coreListHandler[T]) GetByGenre(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	genre := c.Param("genre")
	if genre == "" {
		log.Warn().Msg("Genre is required")
		responses.RespondBadRequest(c, nil, "Genre is required")
		return
	}

	log.Debug().
		Str("genre", genre).
		Msg("Getting lists by genre")

	var zero T
	mediaType := mediatypes.GetMediaTypeFromTypeName(zero)
	// Create query options
	options := mediatypes.QueryOptions{
		Genre:     genre,
		MediaType: mediaType,
	}

	// Search lists by genre
	lists, err := h.listService.Search(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Str("genre", genre).
			Msg("Failed to retrieve lists by genre")
		responses.RespondInternalError(c, err, "Failed to retrieve lists")
		return
	}

	log.Info().
		Str("genre", genre).
		Int("count", len(lists)).
		Msg("Lists by genre retrieved successfully")
	responses.RespondOK(c, lists, "Lists retrieved successfully")
}

// Search godoc
// @Summary Search lists
// @Description Searches for lists that match the query
// @Tags lists
// @Accept json
// @Produce json
// @Param listType path string true "List type (e.g. 'playlist', 'collection')"
// @Param q query string true "Search query"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Playlist]] "Lists retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /api/v1/{listType}/search [get]
func (h *coreListHandler[T]) Search(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	query := c.Query("q")
	if query == "" {
		log.Warn().Msg("Search query is required")
		responses.RespondBadRequest(c, nil, "Search query is required")
		return
	}

	log.Debug().
		Str("query", query).
		Msg("Searching lists")

	var zero T
	mediaType := mediatypes.GetMediaTypeFromTypeName(zero)

	// Create query options
	options := mediatypes.QueryOptions{
		Query:     query,
		MediaType: mediaType,
	}

	// Search lists
	lists, err := h.listService.Search(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Str("query", query).
			Msg("Failed to search lists")
		responses.RespondInternalError(c, err, "Failed to search lists")
		return
	}

	log.Info().
		Str("query", query).
		Int("count", len(lists)).
		Msg("Lists search completed successfully")
	responses.RespondOK(c, lists, "Lists retrieved successfully")
}

// AddItem godoc
// @Summary Add an item to a playlist
// @Description Adds a media item to an existing playlist
// @Tags lists
// @Accept json
// @Produce json
// @Param listId path int true "List ID"
// @Param itemID path string true "Item ID to add"
// @Param listType path string true "List type (e.g. 'playlist', 'collection')"
// @Success 200 {object} responses.APIResponse[string] "Item added to playlist"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "List not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /api/v1/{listType}/{listId}/items/{itemID} [post]
func (h *coreListHandler[T]) AddItem(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid playlist ID")
		responses.RespondBadRequest(c, err, "Invalid play		list ID")
		return
	}

	itemIDStr := c.Param("itemID")
	// Parse item ID
	itemID, err := strconv.ParseUint(itemIDStr, 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("itemID", itemIDStr).Msg("Invalid item ID format")
		responses.RespondBadRequest(c, err, "Invalid item ID format")
		return
	}
	listIDStr := c.Param("id")
	// Parse list ID
	listID, err := strconv.ParseUint(listIDStr, 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("listID", listIDStr).Msg("Invalid list ID format")
		responses.RespondBadRequest(c, err, "Invalid list ID format")
		return
	}

	err = h.listService.AddItem(ctx, listID, itemID)
	if err != nil {
		log.Error().Err(err).
			Uint64("id", id).
			Uint64("itemID", itemID).
			Msg("Failed to add item to playlist")
		responses.RespondInternalError(c, err, "Failed to add item to playlist")
		return
	}

	log.Info().
		Uint64("id", id).
		Uint64("itemID", itemID).
		Msg("Item added to playlist successfully")
	responses.RespondOK(c, "Item added to playlist", "Item added to playlist successfully")
}
