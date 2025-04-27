// handlers/core_lists.go
package handlers

import (
	"github.com/gin-gonic/gin"
	"suasor/utils"

	"suasor/clients/media/types"
	"suasor/services"
	_ "suasor/types/models"
	"suasor/types/responses"
	"suasor/utils/logger"
)

type CoreListHandler[T types.ListData] interface {
	CoreMediaItemHandler[T]

	GetAll(c *gin.Context)
	GetByID(c *gin.Context)
	GetItemsByListID(c *gin.Context)
	GetByGenre(c *gin.Context)
	Search(c *gin.Context)
}

// coreListHandler[T] handles operations for lists in the database
type coreListHandler[T types.ListData] struct {
	CoreMediaItemHandler[T]
	listService services.CoreListService[T]
}

// NewCoreListHandler[T] creates a new core playlist handler
func NewCoreListHandler[T types.ListData](
	CoreMediaItemHandler CoreMediaItemHandler[T],
	listService services.CoreListService[T],
) CoreListHandler[T] {
	return &coreListHandler[T]{
		CoreMediaItemHandler: CoreMediaItemHandler,
		listService:          listService,
	}
}

// GetAll godoc
//
//	@Summary		Get all lists
//	@Description	Retrieves all lists in the database
//	@Tags			lists
//	@Accept			json
//	@Produce    json
//	@Param      limit		  query		int							false	"Maximum number of lists to return (default 10)"
//	@Param			offset		query		int						false	"Offset for pagination (default 0)"
//	@Param			listType	path		string				true	"List type (e.g. 'playlist', 'collection')"
//	@Success		200			{object}	responses.APIResponse[[]models.MediaItem[types.ListData]]	"Lists retrieved successfully"
//	@Failure		500			{object}	responses.ErrorResponse[any]									"Server error"
//	@Router			/{listType} [get]
func (h *coreListHandler[T]) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Debug().Msg("Getting all lists")

	limit := utils.GetLimit(c, 10, 100, true)
	offset := utils.GetOffset(c, 0)

	// Get all lists
	lists, err := h.listService.GetAll(ctx, limit, offset)
	if handleServiceError(c, err, "Failed to retrieve lists", "", "Failed to retrieve lists") {
		return
	}

	log.Info().
		Int("count", len(lists)).
		Msg("Lists retrieved successfully")
	responses.RespondOK(c, lists, "Lists retrieved successfully")
}

// GetByID godoc
//
//	@Summary		Get playlist by ID
//	@Description	Retrieves a specific playlist by ID
//	@Tags			lists
//	@Accept			json
//	@Produce		json
//	@Param			listID		path		int																true	"List ID"
//	@Param			listType	path		string															true	"List type (e.g. 'playlist', 'collection')"
//	@Success		200			{object}	responses.APIResponse[models.MediaItem[types.ListData]]	"List retrieved successfully"
//	@Failure		400			{object}	responses.ErrorResponse[any]									"Invalid request"
//	@Failure		404			{object}	responses.ErrorResponse[any]									"List not found"
//	@Failure		500			{object}	responses.ErrorResponse[any]									"Server error"
//	@Router			/{listType}/{listID} [get]
func (h *coreListHandler[T]) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	listID, err := checkItemID(c, "listID")
	if err != nil {
		return
	}

	log.Debug().
		Uint64("listID", listID).
		Msg("Getting playlist by ID")

	playlist, err := h.listService.GetByID(ctx, listID)
	if handleServiceError(c, err, "Failed to retrieve playlist", "List not found", "List not found") {
		return
	}

	log.Info().
		Uint64("id", listID).
		Msg("List retrieved successfully")
	responses.RespondOK(c, playlist, "List retrieved successfully")
}

// GetListItems godoc
//
//	@Summary		Get tracks in a playlist
//	@Description	Retrieves all tracks in a specific playlist
//	@Tags			lists
//	@Accept			json
//	@Produce		json
//	@Param			listID		path		int											true	"List ID"
//	@Param			listType	path		string										true	"List type (e.g. 'playlist', 'collection')"
//	@Success		200			{object}	responses.APIResponse[models.MediaItemList]	"Tracks retrieved successfully"
//	@Failure		400			{object}	responses.ErrorResponse[any]				"Invalid request"
//	@Failure		404			{object}	responses.ErrorResponse[any]				"List not found"
//	@Failure		500			{object}	responses.ErrorResponse[any]				"Server error"
//	@Router			/{listType}/{listID}/items [get]
func (h *coreListHandler[T]) GetItemsByListID(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	listID, err := checkItemID(c, "listID")

	log.Debug().
		Uint64("id", listID).
		Msg("Getting tracks for playlist")

	playlist, err := h.listService.GetByID(ctx, listID)
	if handleServiceError(c, err, "Failed to retrieve playlist", "List not found", "List not found") {
		return
	}

	itemList := playlist.GetData().GetItemList()

	log.Info().
		Uint64("listID", listID).
		Int("itemCount", len(itemList.Items)).
		Msg("List tracks retrieved successfully")
	responses.RespondOK(c, itemList, "Items retrieved successfully")
}

// GetByGenre godoc
//
//	@Summary		Get lists by genre
//	@Description	Retrieves lists that match a specific genre
//	@Tags			lists
//	@Accept			json
//	@Produce		json
//	@Param			genre		path		string															true	"Genre name"
//	@Param			listType	path		string															true	"List type (e.g. 'playlist', 'collection')"
//	@Success		200			{object}	responses.APIResponse[[]models.MediaItem[types.ListData]]	"Lists retrieved successfully"
//	@Failure		400			{object}	responses.ErrorResponse[any]									"Invalid request"
//	@Failure		500			{object}	responses.ErrorResponse[any]									"Server error"
//	@Router			/{listType}/genre/{genre} [get]
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
	mediaType := types.GetMediaTypeFromTypeName(zero)
	// Create query options
	options := types.QueryOptions{
		Genre:     genre,
		MediaType: mediaType,
	}

	// Search lists by genre
	lists, err := h.listService.Search(ctx, options)
	if handleServiceError(c, err, "Failed to retrieve lists by genre", "", "Failed to retrieve lists") {
		return
	}

	log.Info().
		Str("genre", genre).
		Int("count", len(lists)).
		Msg("Lists by genre retrieved successfully")
	responses.RespondOK(c, lists, "Lists retrieved successfully")
}

// Search godoc
//
//	@Summary		Search lists
//	@Description	Searches for lists that match the query
//	@Tags			lists
//	@Accept			json
//	@Produce		json
//	@Param			listType	path		string															true	"List type (e.g. 'playlist', 'collection')"
//	@Param			q			query		string															true	"Search query"
//	@Success		200			{object}	responses.APIResponse[[]models.MediaItem[types.ListData]]	"Lists retrieved successfully"
//	@Failure		400			{object}	responses.ErrorResponse[any]									"Invalid request"
//	@Failure		500			{object}	responses.ErrorResponse[any]									"Server error"
//	@Router			/{listType}/search [get]
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
	mediaType := types.GetMediaTypeFromTypeName(zero)

	// Create query options
	options := types.QueryOptions{
		Query:     query,
		MediaType: mediaType,
	}

	// Search lists
	lists, err := h.listService.Search(ctx, options)
	if handleServiceError(c, err, "Failed to search lists", "", "Failed to search lists") {
		return
	}

	log.Info().
		Str("query", query).
		Int("count", len(lists)).
		Msg("Lists search completed successfully")
	responses.RespondOK(c, lists, "Lists retrieved successfully")
}
