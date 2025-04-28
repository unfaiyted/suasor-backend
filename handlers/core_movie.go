// handlers/core_movie.go
package handlers

import (
	"github.com/gin-gonic/gin"

	"suasor/clients/media/types"
	"suasor/services"
	_ "suasor/types/models"
	"suasor/types/responses"
	"suasor/utils"
	"suasor/utils/logger"
)

// CoreMovieHandler handles operations for movies in the database
type CoreMovieHandler interface {
	CoreMediaItemHandler[*types.Movie]

	GetByActor(c *gin.Context)
	GetByDirector(c *gin.Context)
}

// coreMovieHandler handles operations for movies in the database
type coreMovieHandler struct {
	CoreMediaItemHandler[*types.Movie]
	itemService services.CoreMediaItemService[*types.Movie]
}

// NewcoreMovieHandler creates a new core movie handler
func NewCoreMovieHandler(
	coreHandler CoreMediaItemHandler[*types.Movie],
	itemService services.CoreMediaItemService[*types.Movie],
) CoreMovieHandler {
	return &coreMovieHandler{
		CoreMediaItemHandler: coreHandler,
		itemService:          itemService,
	}
}

// GetAll godoc
//
//	@Summary		Get all movies
//	@Description	Retrieves all movies in the database
//	@Tags			movies, core
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int															false	"Maximum number of movies to return (default 20)"
//	@Param			offset	query		int															false	"Offset for pagination (default 0)"
//	@Success		200		{object}	responses.APIResponse[responses.MediaItemList[types.Movie]]	"Movies retrieved successfully"
//	@Failure		500		{object}	responses.ErrorResponse[any]								"Server error"
//	@Router			/media/movies [get]
func (h *coreMovieHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	log.Debug().Msg("Getting all movies")
	limit := utils.GetLimit(c, 20, 100, true)
	offset := utils.GetOffset(c, 0)

	// Get all movies
	movies, err := h.itemService.GetAll(ctx, limit, offset)
	if handleServiceError(c, err, "Failed to retrieve movies", "", "Failed to retrieve movies") {
		return
	}

	log.Info().
		Int("count", len(movies)).
		Msg("Movies retrieved successfully")
	responses.RespondMediaItemListOK(c, movies, "Movies retrieved successfully")
}

// GetByID godoc
//
//	@Summary		Get movie by ID
//	@Description	Retrieves a specific movie by ID
//	@Tags			movies, core
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int															true	"Movie ID"
//	@Success		200	{object}	responses.APIResponse[models.MediaItem[types.Movie]]	"Movie retrieved successfully"
//	@Failure		400	{object}	responses.ErrorResponse[any]								"Invalid request"
//	@Failure		404	{object}	responses.ErrorResponse[any]								"Movie not found"
//	@Failure		500	{object}	responses.ErrorResponse[any]								"Server error"
//	@Router			/media/movie/{id} [get]
func (h *coreMovieHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	id, err := checkItemID(c, "id")
	if err != nil {
		return
	}

	log.Debug().
		Uint64("id", id).
		Msg("Getting movie by ID")

	movie, err := h.itemService.GetByID(ctx, id)
	if handleServiceError(c, err, "Failed to retrieve movie", "Movie not found", "Movie not found") {
		return
	}

	log.Info().
		Uint64("id", id).
		Msg("Movie retrieved successfully")
	responses.RespondOK(c, movie, "Movie retrieved successfully")
}

// GetByGenre godoc
//
//	@Summary		Get movies by genre
//	@Description	Retrieves movies that match a specific genre
//	@Tags			movies, core
//	@Accept			json
//	@Produce		json
//	@Param			genre	path		string														true	"Genre name"
//	@Param			limit	query		int															false	"Maximum number of movies to return (default 20)"
//	@Success		200		{object}	responses.APIResponse[responses.MediaItemList[types.Movie]]	"Movies retrieved successfully"
//	@Failure		400		{object}	responses.ErrorResponse[any]								"Invalid request"
//	@Failure		500		{object}	responses.ErrorResponse[any]								"Server error"
//	@Router			/media/movie/genre/{genre} [get]
func (h *coreMovieHandler) GetByGenre(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	genre, ok := checkRequiredStringParam(c, "genre", "Genre is required")
	if !ok {
		return
	}

	limit := utils.GetLimit(c, 20, 100, true)

	log.Debug().
		Str("genre", genre).
		Int("limit", limit).
		Msg("Getting movies by genre")

	// Create query options
	options := types.QueryOptions{
		Genre:     genre,
		MediaType: types.MediaTypeMovie,
		Limit:     limit,
	}

	// Search movies by genre
	movies, err := h.itemService.Search(ctx, options)
	if handleServiceError(c, err, "Failed to retrieve movies by genre", "", "Failed to retrieve movies") {
		return
	}

	log.Info().
		Str("genre", genre).
		Int("count", len(movies)).
		Msg("Movies by genre retrieved successfully")
	responses.RespondMediaItemListOK(c, movies, "Movies retrieved successfully")
}

// GetByYear godoc
//
//	@Summary		Get movies by year
//	@Description	Retrieves movies released in a specific year
//	@Tags			movies, core
//	@Accept			json
//	@Produce		json
//	@Param			year	path		int															true	"Release year"
//	@Param			limit	query		int															false	"Maximum number of movies to return (default 20)"
//	@Success		200		{object}	responses.APIResponse[responses.MediaItemList[types.Movie]]	"Movies retrieved successfully"
//	@Failure		400		{object}	responses.ErrorResponse[any]								"Invalid request"
//	@Failure		500		{object}	responses.ErrorResponse[any]								"Server error"
//	@Router			/media/movie/year/{year} [get]
func (h *coreMovieHandler) GetByYear(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	year, ok := checkYear(c, "year")
	if !ok {
		return
	}

	limit := utils.GetLimit(c, 20, 100, true)

	log.Debug().
		Int("year", year).
		Int("limit", limit).
		Msg("Getting movies by year")

	// Create query options
	options := types.QueryOptions{
		Year:      year,
		MediaType: types.MediaTypeMovie,
		Limit:     limit,
	}

	// Search movies by year
	movies, err := h.itemService.Search(ctx, options)
	if handleServiceError(c, err, "Failed to retrieve movies by year", "", "Failed to retrieve movies") {
		return
	}

	log.Info().
		Int("year", year).
		Int("count", len(movies)).
		Msg("Movies by year retrieved successfully")
	responses.RespondMediaItemListOK(c, movies, "Movies retrieved successfully")
}

// GetByActor godoc
//
//	@Summary		Get movies by actor
//	@Description	Retrieves movies featuring a specific actor
//	@Tags			movies, core
//	@Accept			json
//	@Produce		json
//	@Param			actor	path		string														true	"Actor name"
//	@Param			limit	query		int															false	"Maximum number of movies to return (default 20)"
//	@Success		200		{object}	responses.APIResponse[responses.MediaItemList[types.Movie]]	"Movies retrieved successfully"
//	@Failure		400		{object}	responses.ErrorResponse[any]								"Invalid request"
//	@Failure		500		{object}	responses.ErrorResponse[any]								"Server error"
//	@Router			/media/movie/actor/{actor} [get]
func (h *coreMovieHandler) GetByActor(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	actor, ok := checkRequiredStringParam(c, "actor", "Actor name is required")
	if !ok {
		return
	}

	limit := utils.GetLimit(c, 20, 100, true)

	log.Debug().
		Str("actor", actor).
		Int("limit", limit).
		Msg("Getting movies by actor")

	// Create query options
	options := types.QueryOptions{
		Actor:     actor,
		MediaType: types.MediaTypeMovie,
		Limit:     limit,
	}

	// Search movies by actor
	movies, err := h.itemService.Search(ctx, options)
	if handleServiceError(c, err, "Failed to retrieve movies by actor", "", "Failed to retrieve movies") {
		return
	}

	log.Info().
		Str("actor", actor).
		Int("count", len(movies)).
		Msg("Movies by actor retrieved successfully")
	responses.RespondMediaItemListOK(c, movies, "Movies retrieved successfully")
}

// GetByDirector godoc
//
//	@Summary		Get movies by director
//	@Description	Retrieves movies directed by a specific director
//	@Tags			movies, core
//	@Accept			json
//	@Produce		json
//	@Param			director	path		string														true	"Director name"
//	@Param			limit		query		int															false	"Maximum number of movies to return (default 20)"
//	@Success		200			{object}	responses.APIResponse[responses.MediaItemList[types.Movie]]	"Movies retrieved successfully"
//	@Failure		400			{object}	responses.ErrorResponse[any]								"Invalid request"
//	@Failure		500			{object}	responses.ErrorResponse[any]								"Server error"
//	@Router			/media/movie/director/{director} [get]
func (h *coreMovieHandler) GetByDirector(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	director, ok := checkRequiredStringParam(c, "director", "Director name is required")
	if !ok {
		return
	}

	limit := utils.GetLimit(c, 20, 100, true)

	log.Debug().
		Str("director", director).
		Int("limit", limit).
		Msg("Getting movies by director")

	// Create query options
	options := types.QueryOptions{
		Director:  director,
		MediaType: types.MediaTypeMovie,
		Limit:     limit,
	}

	// Search movies by director
	movies, err := h.itemService.Search(ctx, options)
	if handleServiceError(c, err, "Failed to retrieve movies by director", "", "Failed to retrieve movies") {
		return
	}

	log.Info().
		Str("director", director).
		Int("count", len(movies)).
		Msg("Movies by director retrieved successfully")
	responses.RespondMediaItemListOK(c, movies, "Movies retrieved successfully")
}

// Search godoc
//
//	@Summary		Search movies
//	@Description	Searches for movies that match the query
//	@Tags			movies, core
//	@Accept			json
//	@Produce		json
//	@Param			q		query		string														true	"Search query"
//	@Param			limit	query		int															false	"Maximum number of movies to return (default 20)"
//	@Success		200		{object}	responses.APIResponse[responses.MediaItemList[types.Movie]]	"Movies retrieved successfully"
//	@Failure		400		{object}	responses.ErrorResponse[any]								"Invalid request"
//	@Failure		500		{object}	responses.ErrorResponse[any]								"Server error"
//	@Router			/media/movie/search [get]
func (h *coreMovieHandler) Search(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	query, ok := checkRequiredQueryParam(c, "q", "Search query is required")
	if !ok {
		return
	}

	limit := utils.GetLimit(c, 20, 100, true)

	log.Debug().
		Str("query", query).
		Int("limit", limit).
		Msg("Searching movies")

	// Create query options
	options := types.QueryOptions{
		Query:     query,
		MediaType: types.MediaTypeMovie,
		Limit:     limit,
	}

	// Search movies
	movies, err := h.itemService.Search(ctx, options)
	if handleServiceError(c, err, "Failed to search movies", "", "Failed to search movies") {
		return
	}

	log.Info().
		Str("query", query).
		Int("count", len(movies)).
		Msg("Movies search completed successfully")
	responses.RespondMediaItemListOK(c, movies, "Movies retrieved successfully")
}

// GetTopRated godoc
//
//	@Summary		Get top rated movies
//	@Description	Retrieves the highest rated movies
//	@Tags			movies, core
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int															false	"Maximum number of movies to return (default 20)"
//	@Success		200		{object}	responses.APIResponse[responses.MediaItemList[types.Movie]]	"Movies retrieved successfully"
//	@Failure		500		{object}	responses.ErrorResponse[any]								"Server error"
//	@Router			/media/movie/top-rated [get]
func (h *coreMovieHandler) GetTopRated(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	limit := utils.GetLimit(c, 20, 100, true)

	log.Debug().
		Int("limit", limit).
		Msg("Getting top rated movies")

	// Create query options
	options := types.QueryOptions{
		MediaType: types.MediaTypeMovie,
		Sort:      "rating",
		SortOrder: types.SortOrderDesc,
		Limit:     limit,
	}

	// Get top rated movies
	movies, err := h.itemService.Search(ctx, options)
	if handleServiceError(c, err, "Failed to retrieve top rated movies", "", "Failed to retrieve movies") {
		return
	}

	log.Info().
		Int("count", len(movies)).
		Msg("Top rated movies retrieved successfully")
	responses.RespondMediaItemListOK(c, movies, "Movies retrieved successfully")
}

// GetRecentlyAdded godoc
//
//	@Summary		Get recently added movies
//	@Description	Retrieves the most recently added movies
//	@Tags			movies, core
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int															false	"Maximum number of movies to return (default 20)"
//	@Param			days	query		int															false	"Number of days to look back (default 30)"
//	@Success		200		{object}	responses.APIResponse[responses.MediaItemList[types.Movie]]	"Movies retrieved successfully"
//	@Failure		500		{object}	responses.ErrorResponse[any]								"Server error"
//	@Router			/media/movie/recently-added [get]
func (h *coreMovieHandler) GetRecentlyAdded(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	limit := utils.GetLimit(c, 20, 100, true)
	days := checkDaysParam(c, 30)

	log.Debug().
		Int("limit", limit).
		Int("days", days).
		Msg("Getting recently added movies")

	// Get recently added movies
	movies, err := h.itemService.GetRecentItems(ctx, days, limit)
	if handleServiceError(c, err, "Failed to retrieve recently added movies", "", "Failed to retrieve movies") {
		return
	}

	log.Info().
		Int("count", len(movies)).
		Msg("Recently added movies retrieved successfully")
	responses.RespondMediaItemListOK(c, movies, "Movies retrieved successfully")
}

// GetByRating godoc
//
//	@Summary		Get movies by rating
//	@Description	Retrieves movies that match a specific rating
//	@Tags			movies, core
//	@Accept			json
//	@Produce		json
//	@Param			rating	path		number														true	"Rating"
//	@Param			limit	query		int															false	"Maximum number of movies to return (default 20)"
//	@Success		200		{object}	responses.APIResponse[responses.MediaItemList[types.Movie]]	"Movies retrieved successfully"
//	@Failure		400		{object}	responses.ErrorResponse[any]								"Invalid request"
//	@Failure		500		{object}	responses.ErrorResponse[any]								"Server error"
//	@Router			/media/movie/rating/{rating} [get]
func (h *coreMovieHandler) GetByRating(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	rating, ok := checkRating(c, "rating")
	if !ok {
		return
	}

	limit := utils.GetLimit(c, 20, 100, true)

	log.Debug().
		Float64("rating", rating).
		Int("limit", limit).
		Msg("Getting movies by rating")

	// Create query options
	options := types.QueryOptions{
		MinimumRating: float32(rating),
		MediaType:     types.MediaTypeMovie,
		Limit:         limit,
	}

	// Search movies by rating
	movies, err := h.itemService.Search(ctx, options)
	if handleServiceError(c, err, "Failed to retrieve movies by rating", "", "Failed to retrieve movies") {
		return
	}

	log.Info().
		Float64("rating", rating).
		Int("count", len(movies)).
		Msg("Movies by rating retrieved successfully")
	responses.RespondMediaItemListOK(c, movies, "Movies retrieved successfully")
}

// GetLatestByAdded godoc
//
//	@Summary		Get latest added movies
//	@Description	Retrieves the most recently added movies
//	@Tags			movies, core
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int															false	"Maximum number of movies to return (default 20)"
//	@Param			days	query		int															false	"Number of days to look back (default 30)"
//	@Success		200		{object}	responses.APIResponse[responses.MediaItemList[types.Movie]]	"Movies retrieved successfully"
//	@Failure		500		{object}	responses.ErrorResponse[any]								"Server error"
//	@Router			/media/movie/latest [get]
func (h *coreMovieHandler) GetLatestByAdded(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	limit := utils.GetLimit(c, 20, 100, true)
	days := checkDaysParam(c, 30)

	log.Debug().
		Int("limit", limit).
		Int("days", days).
		Msg("Getting latest added movies")

	// Get latest added movies
	movies, err := h.itemService.GetRecentItems(ctx, days, limit)
	if handleServiceError(c, err, "Failed to retrieve latest added movies", "", "Failed to retrieve movies") {
		return
	}

	log.Info().
		Int("count", len(movies)).
		Msg("Latest added movies retrieved successfully")
	responses.RespondMediaItemListOK(c, movies, "Movies retrieved successfully")
}

// GetPopularByID godoc
//
//	@Summary		Get popular movies
//	@Description	Retrieves the most popular movies
//	@Tags			movies, core
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int															false	"Maximum number of movies to return (default 20)"
//	@Success		200		{object}	responses.APIResponse[responses.MediaItemList[types.Movie]]	"Movies retrieved successfully"
//	@Failure		500		{object}	responses.ErrorResponse[any]								"Server error"
//	@Router			/media/movie/popular [get]
func (h *coreMovieHandler) GetPopular(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	limit := utils.GetLimit(c, 20, 100, true)

	log.Debug().
		Int("limit", limit).
		Msg("Getting popular movies")

	// Get popular movies
	movies, err := h.itemService.GetMostPlayed(ctx, limit)
	if handleServiceError(c, err, "Failed to retrieve popular movies", "", "Failed to retrieve movies") {
		return
	}

	log.Info().
		Int("count", len(movies)).
		Msg("Popular movies retrieved successfully")
	responses.RespondMediaItemListOK(c, movies, "Movies retrieved successfully")
}

// GetByClientItemID godoc
//
//	@Summary		Get movies by client-specific ID
//	@Description	Retrieves movies associated with a specific client
//	@Tags			movies, core
//	@Accept			json
//	@Produce		json
//	@Param			clientID		path		int															true	"Client ID"
//	@Param			clientItemID	path		string														true	"Client Item ID"
//	@Param			limit			query		int															false	"Maximum number of movies to return (default 20)"
//	@Success		200				{object}	responses.APIResponse[responses.MediaItemList[types.Movie]]	"Movies retrieved successfully"
//	@Failure		400				{object}	responses.ErrorResponse[any]								"Invalid request"
//	@Failure		404				{object}	responses.ErrorResponse[any]								"Movie not found"
//	@Failure		500				{object}	responses.ErrorResponse[any]								"Server error"
//	@Router			/media/movie/client/{clientID}/item/{clientItemID} [get]
func (h *coreMovieHandler) GetByClientItemID(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	clientID, err := checkItemID(c, "clientID")
	if err != nil {
		return
	}

	clientItemID, err := checkClientItemID(c, "clientItemID")
	if err != nil {
		return
	}

	limit := utils.GetLimit(c, 20, 100, true)

	log.Debug().
		Uint64("clientID", clientID).
		Str("clientItemID", clientItemID).
		Int("limit", limit).
		Msg("Getting movies by client ID")

	// Get movies by client ID
	movie, err := h.itemService.GetByClientItemID(ctx, clientID, clientItemID)
	if handleServiceError(c, err, "Failed to retrieve movies by client ID", "", "Failed to retrieve movies") {
		return
	}

	log.Info().
		Uint64("clientID", clientID).
		Str("clientItemID", clientItemID).
		Msg("Movies by client ID retrieved successfully")
	responses.RespondOK(c, movie, "Movies retrieved successfully")
}
