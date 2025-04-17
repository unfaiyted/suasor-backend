// handlers/core_movie.go
package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	mediatypes "suasor/client/media/types"
	"suasor/services"
	"suasor/types/responses"
	"suasor/utils"
)

// CoreMovieHandler handles operations for movies in the database
type CoreMovieHandler struct {
	CoreMediaItemHandler[mediatypes.Movie]
	movieService services.CoreMediaItemService[*mediatypes.Movie]
}

// NewCoreMovieHandler creates a new core movie handler
func NewCoreMovieHandler(
	movieService services.CoreMediaItemService[*mediatypes.Movie],
) *CoreMovieHandler {
	return &CoreMovieHandler{
		movieService: movieService,
	}
}

// GetAll godoc
// @Summary Get all movies
// @Description Retrieves all movies in the database
// @Tags movies
// @Accept json
// @Produce json
// @Param limit query int false "Maximum number of movies to return (default 20)"
// @Param offset query int false "Offset for pagination (default 0)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[*mediatypes.Movie]] "Movies retrieved successfully"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /movies [get]
func (h *CoreMovieHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	log.Debug().Msg("Getting all movies")
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		limit = 20
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		offset = 0
	}

	// Get all movies
	movies, err := h.movieService.GetAll(ctx, limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve movies")
		responses.RespondInternalError(c, err, "Failed to retrieve movies")
		return
	}

	log.Info().
		Int("count", len(movies)).
		Msg("Movies retrieved successfully")
	responses.RespondOK(c, movies, "Movies retrieved successfully")
}

// GetByID godoc
// @Summary Get movie by ID
// @Description Retrieves a specific movie by ID
// @Tags movies
// @Accept json
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} responses.APIResponse[models.MediaItem[*mediatypes.Movie]] "Movie retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 404 {object} responses.ErrorResponse[any] "Movie not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /movies/{id} [get]
func (h *CoreMovieHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Warn().Err(err).Str("id", c.Param("id")).Msg("Invalid movie ID")
		responses.RespondBadRequest(c, err, "Invalid movie ID")
		return
	}

	log.Debug().
		Uint64("id", id).
		Msg("Getting movie by ID")

	movie, err := h.movieService.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("id", id).
			Msg("Failed to retrieve movie")
		responses.RespondNotFound(c, err, "Movie not found")
		return
	}

	log.Info().
		Uint64("id", id).
		Msg("Movie retrieved successfully")
	responses.RespondOK(c, movie, "Movie retrieved successfully")
}

// GetByGenre godoc
// @Summary Get movies by genre
// @Description Retrieves movies that match a specific genre
// @Tags movies
// @Accept json
// @Produce json
// @Param genre path string true "Genre name"
// @Param limit query int false "Maximum number of movies to return (default 20)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[*mediatypes.Movie]] "Movies retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /movies/genre/{genre} [get]
func (h *CoreMovieHandler) GetByGenre(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	genre := c.Param("genre")
	if genre == "" {
		log.Warn().Msg("Genre is required")
		responses.RespondBadRequest(c, nil, "Genre is required")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		limit = 20
	}

	log.Debug().
		Str("genre", genre).
		Int("limit", limit).
		Msg("Getting movies by genre")

	// Create query options
	options := mediatypes.QueryOptions{
		Genre:     genre,
		MediaType: mediatypes.MediaTypeMovie,
		Limit:     limit,
	}

	// Search movies by genre
	movies, err := h.movieService.Search(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Str("genre", genre).
			Msg("Failed to retrieve movies by genre")
		responses.RespondInternalError(c, err, "Failed to retrieve movies")
		return
	}

	log.Info().
		Str("genre", genre).
		Int("count", len(movies)).
		Msg("Movies by genre retrieved successfully")
	responses.RespondOK(c, movies, "Movies retrieved successfully")
}

// GetByYear godoc
// @Summary Get movies by year
// @Description Retrieves movies released in a specific year
// @Tags movies
// @Accept json
// @Produce json
// @Param year path int true "Release year"
// @Param limit query int false "Maximum number of movies to return (default 20)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[*mediatypes.Movie]] "Movies retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /movies/year/{year} [get]
func (h *CoreMovieHandler) GetByYear(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	yearStr := c.Param("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		log.Warn().Err(err).Str("year", yearStr).Msg("Invalid year format")
		responses.RespondBadRequest(c, err, "Invalid year format")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		limit = 20
	}

	log.Debug().
		Int("year", year).
		Int("limit", limit).
		Msg("Getting movies by year")

	// Create query options
	options := mediatypes.QueryOptions{
		Year:      year,
		MediaType: mediatypes.MediaTypeMovie,
		Limit:     limit,
	}

	// Search movies by year
	movies, err := h.movieService.Search(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Int("year", year).
			Msg("Failed to retrieve movies by year")
		responses.RespondInternalError(c, err, "Failed to retrieve movies")
		return
	}

	log.Info().
		Int("year", year).
		Int("count", len(movies)).
		Msg("Movies by year retrieved successfully")
	responses.RespondOK(c, movies, "Movies retrieved successfully")
}

// GetByActor godoc
// @Summary Get movies by actor
// @Description Retrieves movies featuring a specific actor
// @Tags movies
// @Accept json
// @Produce json
// @Param actor path string true "Actor name"
// @Param limit query int false "Maximum number of movies to return (default 20)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[*mediatypes.Movie]] "Movies retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /movies/actor/{actor} [get]
func (h *CoreMovieHandler) GetByActor(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	actor := c.Param("actor")
	if actor == "" {
		log.Warn().Msg("Actor name is required")
		responses.RespondBadRequest(c, nil, "Actor name is required")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		limit = 20
	}

	log.Debug().
		Str("actor", actor).
		Int("limit", limit).
		Msg("Getting movies by actor")

	// Create query options
	options := mediatypes.QueryOptions{
		Actor:     actor,
		MediaType: mediatypes.MediaTypeMovie,
		Limit:     limit,
	}

	// Search movies by actor
	movies, err := h.movieService.Search(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Str("actor", actor).
			Msg("Failed to retrieve movies by actor")
		responses.RespondInternalError(c, err, "Failed to retrieve movies")
		return
	}

	log.Info().
		Str("actor", actor).
		Int("count", len(movies)).
		Msg("Movies by actor retrieved successfully")
	responses.RespondOK(c, movies, "Movies retrieved successfully")
}

// GetByDirector godoc
// @Summary Get movies by director
// @Description Retrieves movies directed by a specific director
// @Tags movies
// @Accept json
// @Produce json
// @Param director path string true "Director name"
// @Param limit query int false "Maximum number of movies to return (default 20)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[*mediatypes.Movie]] "Movies retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /movies/director/{director} [get]
func (h *CoreMovieHandler) GetByDirector(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	director := c.Param("director")
	if director == "" {
		log.Warn().Msg("Director name is required")
		responses.RespondBadRequest(c, nil, "Director name is required")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		limit = 20
	}

	log.Debug().
		Str("director", director).
		Int("limit", limit).
		Msg("Getting movies by director")

	// Create query options
	options := mediatypes.QueryOptions{
		Director:  director,
		MediaType: mediatypes.MediaTypeMovie,
		Limit:     limit,
	}

	// Search movies by director
	movies, err := h.movieService.Search(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Str("director", director).
			Msg("Failed to retrieve movies by director")
		responses.RespondInternalError(c, err, "Failed to retrieve movies")
		return
	}

	log.Info().
		Str("director", director).
		Int("count", len(movies)).
		Msg("Movies by director retrieved successfully")
	responses.RespondOK(c, movies, "Movies retrieved successfully")
}

// Search godoc
// @Summary Search movies
// @Description Searches for movies that match the query
// @Tags movies
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Maximum number of movies to return (default 20)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[*mediatypes.Movie]] "Movies retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /movies/search [get]
func (h *CoreMovieHandler) Search(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	query := c.Query("q")
	if query == "" {
		log.Warn().Msg("Search query is required")
		responses.RespondBadRequest(c, nil, "Search query is required")
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		limit = 20
	}

	log.Debug().
		Str("query", query).
		Int("limit", limit).
		Msg("Searching movies")

	// Create query options
	options := mediatypes.QueryOptions{
		Query:     query,
		MediaType: mediatypes.MediaTypeMovie,
		Limit:     limit,
	}

	// Search movies
	movies, err := h.movieService.Search(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Str("query", query).
			Msg("Failed to search movies")
		responses.RespondInternalError(c, err, "Failed to search movies")
		return
	}

	log.Info().
		Str("query", query).
		Int("count", len(movies)).
		Msg("Movies search completed successfully")
	responses.RespondOK(c, movies, "Movies retrieved successfully")
}

// GetTopRated godoc
// @Summary Get top rated movies
// @Description Retrieves the highest rated movies
// @Tags movies
// @Accept json
// @Produce json
// @Param limit query int false "Maximum number of movies to return (default 20)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[*mediatypes.Movie]] "Movies retrieved successfully"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /movies/top-rated [get]
func (h *CoreMovieHandler) GetTopRated(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		limit = 20
	}

	log.Debug().
		Int("limit", limit).
		Msg("Getting top rated movies")

	// Create query options
	options := mediatypes.QueryOptions{
		MediaType: mediatypes.MediaTypeMovie,
		Sort:      "rating",
		SortOrder: mediatypes.SortOrderDesc,
		Limit:     limit,
	}

	// Get top rated movies
	movies, err := h.movieService.Search(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Msg("Failed to retrieve top rated movies")
		responses.RespondInternalError(c, err, "Failed to retrieve movies")
		return
	}

	log.Info().
		Int("count", len(movies)).
		Msg("Top rated movies retrieved successfully")
	responses.RespondOK(c, movies, "Movies retrieved successfully")
}

// GetRecentlyAdded godoc
// @Summary Get recently added movies
// @Description Retrieves the most recently added movies
// @Tags movies
// @Accept json
// @Produce json
// @Param limit query int false "Maximum number of movies to return (default 20)"
// @Param days query int false "Number of days to look back (default 30)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[*mediatypes.Movie]] "Movies retrieved successfully"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /movies/recently-added [get]
func (h *CoreMovieHandler) GetRecentlyAdded(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		limit = 20
	}

	days, err := strconv.Atoi(c.DefaultQuery("days", "30"))
	if err != nil {
		days = 30
	}

	log.Debug().
		Int("limit", limit).
		Int("days", days).
		Msg("Getting recently added movies")

	// Get recently added movies
	movies, err := h.movieService.GetRecentItems(ctx, days, limit)
	if err != nil {
		log.Error().Err(err).
			Msg("Failed to retrieve recently added movies")
		responses.RespondInternalError(c, err, "Failed to retrieve movies")
		return
	}

	log.Info().
		Int("count", len(movies)).
		Msg("Recently added movies retrieved successfully")
	responses.RespondOK(c, movies, "Movies retrieved successfully")
}
