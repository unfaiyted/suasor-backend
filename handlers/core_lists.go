// handlers/core_playlists.go
package handlers

import (
	"github.com/gin-gonic/gin"
	"strconv"

	mediatypes "suasor/clients/media/types"
	"suasor/services"
	models "suasor/types/models"
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

// coreListHandler[T] handles operations for playlists in the database
type coreListHandler[T mediatypes.ListData] struct {
	CoreMediaItemHandler[T]
	listService services.CoreListService[T]
}

// NewcoreListHandler[T] creates a new core playlist handler
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
// @Summary Get all playlists
// @Description Retrieves all playlists in the database
// @Tags playlists
// @Accept json
// @Produce json
// @Param limit query int false "Maximum number of playlists to return (default 10)"
// @Param offset query int false "Offset for pagination (default 0)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Playlist]] "Lists retrieved successfully"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /playlists [get]
func (h *coreListHandler[T]) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Debug().Msg("Getting all playlists")

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		offset = 0
	}

	// Get all playlists
	playlists, err := h.listService.GetAll(ctx, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve playlists")
		responses.RespondInternalError(c, err, "Failed to retrieve playlists")
		return
	}

	log.Info().
		Int("count", len(playlists)).
		Msg("Lists retrieved successfully")
	responses.RespondOK(c, playlists, "Lists retrieved successfully")
}

// GetByID godoc
// @Summary Get playlist by ID
// @Description Retrieves a specific playlist by ID
// @Tags playlists
// @Accept json
// @Produce json
// @Param id path int true "List ID"
// @Success 200 {object} responses.APIResponse[models.MediaItem[mediatypes.Playlist]] "List retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "List not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /playlists/{id} [get]
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

// GetListTracks godoc
// @Summary Get tracks in a playlist
// @Description Retrieves all tracks in a specific playlist
// @Tags playlists
// @Accept json
// @Produce json
// @Param id path int true "List ID"
// @Success 200 {object} responses.APIResponse[[]mediatypes.Track] "Tracks retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "List not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /playlists/{id}/tracks [get]
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
// @Summary Get playlists by genre
// @Description Retrieves playlists that match a specific genre
// @Tags playlists
// @Accept json
// @Produce json
// @Param genre path string true "Genre name"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Playlist]] "Lists retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /playlists/genre/{genre} [get]
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
		Msg("Getting playlists by genre")

	var zero T
	mediaType := mediatypes.GetMediaTypeFromTypeName(zero)
	// Create query options
	options := mediatypes.QueryOptions{
		Genre:     genre,
		MediaType: mediaType,
	}

	// Search playlists by genre
	playlists, err := h.listService.Search(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Str("genre", genre).
			Msg("Failed to retrieve playlists by genre")
		responses.RespondInternalError(c, err, "Failed to retrieve playlists")
		return
	}

	log.Info().
		Str("genre", genre).
		Int("count", len(playlists)).
		Msg("Lists by genre retrieved successfully")
	responses.RespondOK(c, playlists, "Lists retrieved successfully")
}

// Search godoc
// @Summary Search playlists
// @Description Searches for playlists that match the query
// @Tags playlists
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Playlist]] "Lists retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /playlists/search [get]
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
		Msg("Searching playlists")

	var zero T
	mediaType := mediatypes.GetMediaTypeFromTypeName(zero)

	// Create query options
	options := mediatypes.QueryOptions{
		Query:     query,
		MediaType: mediaType,
	}

	// Search playlists
	playlists, err := h.listService.Search(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Str("query", query).
			Msg("Failed to search playlists")
		responses.RespondInternalError(c, err, "Failed to search playlists")
		return
	}

	log.Info().
		Str("query", query).
		Int("count", len(playlists)).
		Msg("Lists search completed successfully")
	responses.RespondOK(c, playlists, "Lists retrieved successfully")
}

// AddItem godoc
// @Summary Add an item to a playlist
// @Description Adds a media item to an existing playlist
// @Tags playlists
// @Accept json
// @Produce json
// @Param id path int true "List ID"
// @Param itemID path string true "Item ID to add"
// @Success 200 {object} responses.APIResponse[string] "Item added to playlist"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "List not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /playlists/{id}/items/{itemID} [post]
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
