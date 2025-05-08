// handlers/client_media_playlist.go
package handlers

import (
	"strconv"
	"suasor/clients/media/types"
	clienttypes "suasor/clients/types"
	"suasor/services"
	_ "suasor/types/models"
	"suasor/types/responses"
	"suasor/utils"
	"suasor/utils/logger"

	"github.com/gin-gonic/gin"
)

type ClientListHandler[T clienttypes.ClientMediaConfig, U types.ListData] interface {
	CoreListHandler[U]

	GetClientLists(c *gin.Context)
	GetClientListItems(c *gin.Context)
	GetClientListByID(c *gin.Context)
	GetClientListsByGenre(c *gin.Context)
	GetClientListsByYear(c *gin.Context)
	GetClientListsByActor(c *gin.Context)
	GetClientListsByCreator(c *gin.Context)
	GetClientListsByRating(c *gin.Context)
	GetClientLatestListsByAdded(c *gin.Context)
	GetClientPopularLists(c *gin.Context)
	GetClientTopRatedLists(c *gin.Context)
	SearchClientLists(c *gin.Context)
}

// clientListHandler handles playlist-related operations for media clients
type clientListHandler[T clienttypes.ClientMediaConfig, U types.ListData] struct {
	CoreListHandler[U]
	listService services.ClientListService[T, U]
}

// NewClientListHandler creates a new media client playlist handler
func NewClientListHandler[T clienttypes.ClientMediaConfig, U types.ListData](
	coreHandler CoreListHandler[U],
	listService services.ClientListService[T, U]) ClientListHandler[T, U] {
	return &clientListHandler[T, U]{
		CoreListHandler: coreHandler,
		listService:     listService,
	}
}

// GetListByID godoc
//
//	@Summary		Get list by ID
//	@Description	Retrieves a specific list from the client by ID
//	@Tags			lists, clients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			clientID	path		int																true	"Client ID"
//	@Param			listID		path		string															true	"List ID"
//	@Param			listType	path		string															true	"List type (e.g. 'playlist', 'collection')"
//	@Success		200			{object}	responses.APIResponse[models.MediaItem[types.ListData]]	"List retrieved"
//	@Failure		400			{object}	responses.ErrorResponse[responses.ErrorDetails]					"Invalid client ID"
//	@Failure		401			{object}	responses.ErrorResponse[responses.ErrorDetails]					"Unauthorized"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]					"Server error"
//	@Router			/client/{clientID}/{listType}/{listID} [get]
func (h *clientListHandler[T, U]) GetClientListByID(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting list by ID")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access playlist without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	// Parse client ID from URL
	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)

	listID := c.Param("listId")

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("playlistID", listID).
		Msg("Retrieving playlist by ID")

	list, err := h.listService.GetClientList(ctx, uid, listID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Str("playlistID", listID).
			Msg("Failed to retrieve playlist")
		responses.RespondInternalError(c, err, "Failed to retrieve playlist")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("playlistID", listID).
		Msg("List retrieved successfully")
	responses.RespondOK(c, list, "List retrieved successfully")
}

// GetClientListItems godoc
//
//	@Summary		Get tracks in a list
//	@Description	Retrieves all tracks in a specific list
//	@Tags			lists, clients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			clientID	path		int											true	"Client ID"
//	@Param			listID		path		string										true	"List ID"
//	@Param			listType	path		string										true	"List type (e.g. 'playlist', 'collection')"
//	@Success		200			{object}	responses.APIResponse[models.MediaItemList]	"Tracks retrieved successfully"
//	@Failure		400			{object}	responses.ErrorResponse[any]				"Invalid request"
//	@Failure		404			{object}	responses.ErrorResponse[any]				"List not found"
//	@Failure		500			{object}	responses.ErrorResponse[any]				"Server error"
//	@Router			/client/{clientID}/{listType}/{listID}/items [get]
func (h *clientListHandler[T, U]) GetClientListItems(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	listID, _ := checkClientItemID(c, "listID")
	clientID, _ := checkClientID(c)

	log.Debug().
		Str("listID", listID).
		Msg("Getting tracks for list")

	list, err := h.listService.GetClientItems(ctx, clientID, listID)
	if handleServiceError(c, err, "Failed to retrieve list", "List not found", "List not found") {
		return
	}

	log.Info().
		Str("listID", listID).
		Int("itemCount", list.GetTotalItems()).
		Msg("List tracks retrieved successfully")
	responses.RespondOK(c, list, "Items retrieved successfully")
}

// GetLists godoc
//
//	@Summary		Get all lists
//	@Description	Retrieves all lists from the client
//	@Tags			lists, clients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			clientID	path		int																true	"Client ID"
//	@Param			listType	path		string															true	"List type (e.g. 'playlist', 'collection')"
//	@Param			limit		query		int																false	"Maximum number of lists to return"
//	@Success		200			{object}	responses.APIResponse[responses.MediaItemList[types.ListData]]	"Lists retrieved"
//	@Failure		400			{object}	responses.ErrorResponse[responses.ErrorDetails]					"Invalid client ID"
//	@Failure		401			{object}	responses.ErrorResponse[responses.ErrorDetails]					"Unauthorized"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]					"Server error"
//	@Router			/client/{clientID}/{listType} [get]
func (h *clientListHandler[T, U]) GetClientLists(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting all lists")

	// Get authenticated user ID
	uid, _ := checkUserAccess(c)
	limit := utils.GetLimit(c, 10, 100, true)

	clientID, _ := checkClientID(c)

	log.Info().
		Uint64("userID", uid).
		Int("count", limit).
		Msg("Retrieving lists")

	lists, err := h.listService.GetClientLists(ctx, clientID, limit)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Int("count", limit).
			Msg("Failed to retrieve lists")
		responses.RespondInternalError(c, err, "Failed to retrieve lists")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", limit).
		Int("listsReturned", len(lists)).
		Msg("Lists retrieved successfully")
	responses.RespondMediaItemListOK(c, lists, "Lists retrieved successfully")
}

// GetListsByGenre godoc
//
//	@Summary		Get lists by genre
//	@Description	Retrieves lists from the client filtered by genre
//	@Tags			lists, clients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			clientID	path		int																true	"Client ID"
//	@Param			listType	path		string															true	"List type (e.g. 'playlist', 'collection')"
//	@Param			genre		path		string															true	"Genre"
//	@Success		200			{object}	responses.APIResponse[responses.MediaItemList[types.ListData]]	"Lists retrieved"
//	@Failure		400			{object}	responses.ErrorResponse[responses.ErrorDetails]					"Invalid client ID"
//	@Failure		401			{object}	responses.ErrorResponse[responses.ErrorDetails]					"Unauthorized"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]					"Server error"
//	@Router			/client/{clientID}/{listType}/genre/{genre} [get]
func (h *clientListHandler[T, U]) GetClientListsByGenre(c *gin.Context) {
	// This would typically query the client with a genre filter
	// For now, just use the SearchClientLists method with a genre query
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	_, _ = checkUserAccess(c)

	// Parse client ID
	clientID, _ := checkClientID(c)

	genre := c.Param("genre")

	options := types.QueryOptions{
		Genre: genre,
	}

	lists, err := h.listService.SearchClientLists(ctx, clientID, options)
	if err != nil {
		log.Error().Err(err).
			Msg("Failed to retrieve lists by genre")
		responses.RespondInternalError(c, err, "Failed to retrieve lists")
		return
	}

	responses.RespondMediaItemListOK(c, lists, "Lists retrieved successfully")
}

// GetListsByYear godoc
//
//	@Summary		Get lists by year
//	@Description	Retrieves lists from the client filtered by year
//	@Tags			lists, clients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			clientID	path		int																true	"Client ID"
//	@Param			listType	path		string															true	"List type (e.g. 'playlist', 'collection')"
//	@Param			year		path		int																true	"Year"
//	@Success		200			{object}	responses.APIResponse[responses.MediaItemList[types.ListData]]	"Lists retrieved"
//	@Failure		400			{object}	responses.ErrorResponse[responses.ErrorDetails]					"Invalid client ID"
//	@Failure		401			{object}	responses.ErrorResponse[responses.ErrorDetails]					"Unauthorized"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]					"Server error"
//	@Router			/client/{clientID}/{listType}/year/{year} [get]
func (h *clientListHandler[T, U]) GetClientListsByYear(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	uid, _ := checkUserAccess(c)
	log.Info().
		Uint64("userID", uid).
		Msg("Retrieving lists by year")

	// Parse client ID
	clientID, exists := checkClientID(c)
	if !exists {
		responses.RespondBadRequest(c, nil, "Invalid client ID")
		return
	}
	// Parse year
	year, exists := checkYear(c, "year")
	if !exists {
		responses.RespondBadRequest(c, nil, "Invalid year")
		return
	}

	options := types.QueryOptions{
		Year: year,
	}

	lists, err := h.listService.SearchClientLists(ctx, clientID, options)
	if err != nil {
		log.Error().Err(err).
			Msg("Failed to retrieve lists by year")
		responses.RespondInternalError(c, err, "Failed to retrieve lists")
		return
	}

	responses.RespondMediaItemListOK(c, lists, "Lists retrieved successfully")
}

// GetListsByActor godoc
//
//	@Summary		Get lists by actor
//	@Description	Retrieves lists from the client filtered by actor
//	@Tags			lists, clients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			clientID	path		int																	true	"Client ID"
//	@Param			listType	path		string															true	"List type (e.g. 'playlist', 'collection')"
//	@Param			actorID		path		string															true	"Actor name"
//	@Success		200			{object}	responses.APIResponse[responses.MediaItemList[types.ListData]]	"Lists retrieved"
//	@Failure		400			{object}	responses.ErrorResponse[responses.ErrorDetails]					"Invalid client ID"
//	@Failure		401			{object}	responses.ErrorResponse[responses.ErrorDetails]					"Unauthorized"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]					"Server error"
//	@Router			/client/{clientID}/{listType}/actor/{actorID} [get]
func (h *clientListHandler[T, U]) GetClientListsByActor(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	uid, _ := checkUserAccess(c)
	log.Debug().
		Uint64("userID", uid).
		Msg("Retrieving lists by actor")

	// Parse client ID
	clientID, _ := checkClientID(c)
	actorID, _ := checkItemID(c, "actorID")

	options := types.QueryOptions{
		PersonID:   actorID,
		PersonType: "Actor",
	}

	lists, err := h.listService.SearchClientLists(ctx, clientID, options)
	if err != nil {
		log.Error().Err(err).
			Msg("Failed to retrieve lists by actor")
		responses.RespondInternalError(c, err, "Failed to retrieve lists")
		return
	}

	responses.RespondMediaItemListOK(c, lists, "Lists retrieved successfully")
}

// GetListsByCreator godoc
//
//	@Summary		Get lists by creator
//	@Description	Retrieves lists from the client filtered by creator
//	@Tags			lists, clients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			clientID	path		int																true	"Client ID"
//	@Param			listType	path		string															true	"List type (e.g. 'playlist', 'collection')"
//	@Param			creatorID		path		string															true	"Creator name"
//	@Success		200			{object}	responses.APIResponse[responses.MediaItemList[types.ListData]]	"Lists retrieved"
//	@Failure		400			{object}	responses.ErrorResponse[responses.ErrorDetails]					"Invalid client ID"
//	@Failure		401			{object}	responses.ErrorResponse[responses.ErrorDetails]					"Unauthorized"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]					"Server error"
//	@Router			/client/{clientID}/{listType}/creator/{creatorID} [get]
func (h *clientListHandler[T, U]) GetClientListsByCreator(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	uid, _ := checkUserAccess(c)
	if uid == 0 {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	// Parse client ID
	clientID, exists := checkClientID(c)
	if !exists {
		responses.RespondBadRequest(c, nil, "Invalid client ID")
		return
	}

	creatorID, _ := checkClientItemID(c, "creatorID")

	options := types.QueryOptions{
		ClientPersonID: creatorID,
		PersonType:     "Creator",
	}

	lists, err := h.listService.SearchClientLists(ctx, clientID, options)
	if err != nil {
		log.Error().Err(err).
			Msg("Failed to retrieve lists by creator")
		responses.RespondInternalError(c, err, "Failed to retrieve lists")
		return
	}

	responses.RespondMediaItemListOK(c, lists, "Lists retrieved successfully")
}

// GetListsByRating godoc
//
//	@Summary		Get lists by rating
//	@Description	Retrieves lists from the client filtered by rating
//	@Tags			lists, clients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			clientID	path		int																true	"Client ID"
//	@Param			listType	path		string															true	"List type (e.g. 'playlist', 'collection')"
//	@Param			rating		query		float															true	"Minimum rating"
//	@Success		200			{object}	responses.APIResponse[responses.MediaItemList[types.ListData]]	"Lists retrieved"
//	@Failure		400			{object}	responses.ErrorResponse[responses.ErrorDetails]					"Invalid client ID"
//	@Failure		401			{object}	responses.ErrorResponse[responses.ErrorDetails]					"Unauthorized"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]					"Server error"
//	@Router			/client/{clientID}/{listType}/rating [get]
func (h *clientListHandler[T, U]) GetClientListsByRating(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	uid, _ := checkUserAccess(c)
	if uid == 0 {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	// Parse client ID
	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)
	if err != nil {
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	// Parse minimum rating
	minRating, err := strconv.ParseFloat(c.Query("rating"), 32)
	if err != nil {
		responses.RespondBadRequest(c, err, "Invalid rating format")
		return
	}

	options := types.QueryOptions{
		MinimumRating: float32(minRating),
	}

	lists, err := h.listService.SearchClientLists(ctx, clientID, options)
	if err != nil {
		log.Error().Err(err).
			Msg("Failed to retrieve lists by rating")
		responses.RespondInternalError(c, err, "Failed to retrieve lists")
		return
	}

	responses.RespondMediaItemListOK(c, lists, "Lists retrieved successfully")
}

// GetLatestListsByAdded godoc
//
//	@Summary		Get latest lists by date added
//	@Description	Retrieves the latest lists from the client sorted by date added
//	@Tags			lists, clients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			clientID	path		int																true	"Client ID"
//	@Param			listType	path		string															true	"List type (e.g. 'playlist', 'collection')"
//	@Param			count		path		int																true	"Number of lists to return"
//	@Success		200			{object}	responses.APIResponse[responses.MediaItemList[types.ListData]]	"Lists retrieved"
//	@Failure		400			{object}	responses.ErrorResponse[responses.ErrorDetails]					"Invalid client ID"
//	@Failure		401			{object}	responses.ErrorResponse[responses.ErrorDetails]					"Unauthorized"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]					"Server error"
//	@Router			/client/{clientID}/{listType}/latest/{count} [get]
func (h *clientListHandler[T, U]) GetClientLatestListsByAdded(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	uid, _ := checkUserAccess(c)
	if uid == 0 {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	// Parse client ID
	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)
	if err != nil {
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	// Parse count
	count, err := strconv.Atoi(c.Param("count"))
	if err != nil {
		responses.RespondBadRequest(c, err, "Invalid count")
		return
	}

	options := types.QueryOptions{
		Limit:     count,
		Sort:      types.SortTypeAddedAt,
		SortOrder: "Desc",
	}

	lists, err := h.listService.SearchClientLists(ctx, clientID, options)
	if err != nil {
		log.Error().Err(err).
			Msg("Failed to retrieve latest lists")
		responses.RespondInternalError(c, err, "Failed to retrieve lists")
		return
	}

	responses.RespondMediaItemListOK(c, lists, "Lists retrieved successfully")
}

// GetPopularLists godoc
//
//	@Summary		Get popular lists
//	@Description	Retrieves the most popular lists from the client
//	@Tags			lists, clients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			clientID	path		int																true	"Client ID"
//	@Param			listType	path		string															true	"List type (e.g. 'playlist', 'collection')"
//	@Param			count		path		int																true	"Number of lists to return"
//	@Success		200			{object}	responses.APIResponse[responses.MediaItemList[types.ListData]]	"Lists retrieved"
//	@Failure		400			{object}	responses.ErrorResponse[responses.ErrorDetails]					"Invalid client ID"
//	@Failure		401			{object}	responses.ErrorResponse[responses.ErrorDetails]					"Unauthorized"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]					"Server error"
//	@Router			/client/{clientID}/{listType}/popular/{count} [get]
func (h *clientListHandler[T, U]) GetClientPopularLists(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	uid, _ := checkUserAccess(c)
	if uid == 0 {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	// Parse client ID
	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)
	if err != nil {
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	// Parse count
	count, err := strconv.Atoi(c.Param("count"))
	if err != nil {
		responses.RespondBadRequest(c, err, "Invalid count")
		return
	}

	options := types.QueryOptions{
		Limit:     count,
		Sort:      types.SortTypePopularity,
		SortOrder: "Desc",
	}

	lists, err := h.listService.SearchClientLists(ctx, clientID, options)
	if err != nil {
		log.Error().Err(err).
			Msg("Failed to retrieve popular lists")
		responses.RespondInternalError(c, err, "Failed to retrieve lists")
		return
	}

	responses.RespondMediaItemListOK(c, lists, "Lists retrieved successfully")
}

// GetTopRatedLists godoc
//
//	@Summary		Get top rated lists
//	@Description	Retrieves the highest rated lists from the client
//	@Tags			lists, clients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			clientID	path		int																true	"Client ID"
//	@Param			listType	path		string															true	"List type (e.g. 'playlist', 'collection')"
//	@Param			count		path		int																true	"Number of lists to return"
//	@Success		200			{object}	responses.APIResponse[responses.MediaItemList[types.ListData]]	"Lists retrieved"
//	@Failure		400			{object}	responses.ErrorResponse[responses.ErrorDetails]					"Invalid client ID"
//	@Failure		401			{object}	responses.ErrorResponse[responses.ErrorDetails]					"Unauthorized"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]					"Server error"
//	@Router			/client/{clientID}/{listType}/top-rated/{count} [get]
func (h *clientListHandler[T, U]) GetClientTopRatedLists(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	_, _ = checkUserAccess(c)

	// Parse client ID
	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)
	if err != nil {
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	// Parse count
	count, err := strconv.Atoi(c.Param("count"))
	if err != nil {
		responses.RespondBadRequest(c, err, "Invalid count")
		return
	}

	options := types.QueryOptions{
		Limit:     count,
		Sort:      types.SortTypeRating,
		SortOrder: types.SortOrderDesc,
	}

	lists, err := h.listService.SearchClientLists(ctx, clientID, options)
	if err != nil {
		log.Error().Err(err).
			Msg("Failed to retrieve top rated lists")
		responses.RespondInternalError(c, err, "Failed to retrieve lists")
		return
	}

	responses.RespondMediaItemListOK(c, lists, "Lists retrieved successfully")
}

// SearchLists godoc
//
//	@Summary		Search lists
//	@Description	Searches for lists matching the given query
//	@Tags			lists, clients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			clientID	path		int																true	"Client ID"
//	@Param			listType	path		string															true	"List type (e.g. 'playlist', 'collection')"
//	@Param			q			query		string															true	"Search query"
//	@Success		200			{object}	responses.APIResponse[responses.MediaItemList[types.ListData]]	"Lists found"
//	@Failure		400			{object}	responses.ErrorResponse[responses.ErrorDetails]					"Invalid request"
//	@Failure		401			{object}	responses.ErrorResponse[responses.ErrorDetails]					"Unauthorized"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]					"Server error"
//	@Router			/client/{clientID}/{listType}/search [get]
func (h *clientListHandler[T, U]) SearchClientLists(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Searching lists")

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to search lists without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)
	query := c.Query("q")
	if query == "" {
		log.Warn().Uint64("userID", uid).Msg("Empty search query provided")
		responses.RespondBadRequest(c, nil, "Search query is required")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Str("query", query).
		Msg("Searching lists")

		// Parse client ID from URL
	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("clientID", c.Param("clientID")).Msg("Invalid client ID format")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}
	options := types.QueryOptions{
		Query: query,
	}

	lists, err := h.listService.SearchClientLists(ctx, clientID, options)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Str("query", query).
			Msg("Failed to search lists")
		responses.RespondInternalError(c, err, "Failed to search lists")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Str("query", query).
		Int("resultsCount", len(lists)).
		Msg("List search completed successfully")
	responses.RespondMediaItemListOK(c, lists, "Lists retrieved successfully")
}
