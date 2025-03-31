// handlers/media_client_movie.go
package handlers

import (
	"strconv"
	mediatypes "suasor/client/media/types"
	clienttypes "suasor/client/types"
	"suasor/services"
	models "suasor/types/models"
	"suasor/types/responses"
	"suasor/utils"

	"github.com/gin-gonic/gin"
)

func createMovieMediaItem[T mediatypes.Movie](clientID uint64, clientType clienttypes.MediaClientType, externalID string, data mediatypes.Movie) models.MediaItem[mediatypes.Movie] {
	return models.MediaItem[mediatypes.Movie]{
		Type:       mediatypes.MediaTypeMovie,
		ClientID:   clientID,
		ClientType: clientType,
		ExternalID: externalID,
		Data:       data,
	}
}

// MediaClientMovieHandler handles movie-related operations for media clients
type MediaClientMovieHandler[T clienttypes.MediaClientConfig] struct {
	movieService services.MediaClientMovieService[T]
}

// NewMediaClientMovieHandler creates a new media client movie handler
func NewMediaClientMovieHandler[T clienttypes.MediaClientConfig](movieService services.MediaClientMovieService[T]) *MediaClientMovieHandler[T] {
	return &MediaClientMovieHandler[T]{
		movieService: movieService,
	}
}

// GetMovieByID godoc
// @Summary Get movie by ID
// @Description Retrieves a specific movie from the client by ID
// @Tags movies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param movieID path string true "Movie ID"
// @Success 200 {object} responses.APIResponse[models.MediaItem[mediatypes.Movie]] "Movie retrieved"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid client ID"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /clients/media/{clientID}/movies/{movieID} [get]
func (h *MediaClientMovieHandler[T]) GetMovieByID(c *gin.Context) {
	ctx := c.Request.Context()

	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	// Parse client ID from URL
	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)
	if err != nil {
		log.Error().Err(err).Str("clientID", c.Param("clientID")).Msg("Invalid client ID format")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	movieID := c.Param("movieID")

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("movieID", movieID).
		Msg("Retrieving movie by ID")

	movie, err := h.movieService.GetMovieByID(ctx, uid, clientID, movieID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Str("movieID", movieID).
			Msg("Failed to retrieve movie")
		responses.RespondInternalError(c, err, "Failed to retrieve movie")
		return
	}

	responses.RespondOK(c, movie, "Movie retrieved successfully")
}

// GetMoviesByGenre godoc
// @Summary Get movies by genre
// @Description Retrieves movies from all connected clients that match the specified genre
// @Tags movies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param genre path string true "Genre name"
// @Success 200 {object} models.APIResponse[[]models.MediaItem[mediatypes.Movie]] "Movies retrieved"
// @Failure 401 {object} models.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} models.ErrorResponse[error] "Server error"
// @Router /movies/genre/{genre} [get]
func (h *MediaClientMovieHandler[T]) GetMoviesByGenre(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)
	genre := c.Param("genre")

	log.Info().
		Uint64("userID", uid).
		Str("genre", genre).
		Msg("Retrieving movies by genre")

	movies, err := h.movieService.GetMoviesByGenre(ctx, uid, genre)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Str("genre", genre).
			Msg("Failed to retrieve movies by genre")
		responses.RespondInternalError(c, err, "Failed to retrieve movies")
		return
	}

	responses.RespondOK(c, movies, "Movies retrieved successfully")
}

// GetMoviesByYear godoc
// @Summary Get movies by release year
// @Description Retrieves movies from all connected clients that were released in the specified year
// @Tags movies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param year path int true "Release year"
// @Success 200 {object} models.APIResponse[[]models.Movie] "Movies retrieved"
// @Failure 400 {object} models.ErrorResponse[error] "Invalid year"
// @Failure 401 {object} models.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} models.ErrorResponse[error] "Server error"
// @Router /movies/year/{year} [get]
func (h *MediaClientMovieHandler[T]) GetMoviesByYear(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		log.Error().Err(err).Str("year", c.Param("year")).Msg("Invalid year format")
		responses.RespondBadRequest(c, err, "Invalid year")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("year", year).
		Msg("Retrieving movies by year")

	movies, err := h.movieService.GetMoviesByYear(ctx, uid, year)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Int("year", year).
			Msg("Failed to retrieve movies by year")
		responses.RespondInternalError(c, err, "Failed to retrieve movies")
		return
	}

	responses.RespondOK(c, movies, "Movies retrieved successfully")
}

func (h *MediaClientMovieHandler[T]) GetMoviesByActor(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)
	actor := c.Param("actor")

	log.Info().
		Uint64("userID", uid).
		Str("actor", actor).
		Msg("Retrieving movies by actor")

	movies, err := h.movieService.GetMoviesByActor(ctx, uid, actor)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Str("actor", actor).
			Msg("Failed to retrieve movies by actor")
		responses.RespondInternalError(c, err, "Failed to retrieve movies")
		return
	}

	responses.RespondOK(c, movies, "Movies retrieved successfully")
}

func (h *MediaClientMovieHandler[T]) GetMoviesByDirector(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)
	director := c.Param("director")

	log.Info().
		Uint64("userID", uid).
		Str("director", director).
		Msg("Retrieving movies by director")

	movies, err := h.movieService.GetMoviesByDirector(ctx, uid, director)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Str("director", director).
			Msg("Failed to retrieve movies by director")
		responses.RespondInternalError(c, err, "Failed to retrieve movies")
		return
	}

	responses.RespondOK(c, movies, "Movies retrieved successfully")
}

func (h *MediaClientMovieHandler[T]) GetMoviesByRating(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	minRating, err := strconv.ParseFloat(c.Query("min"), 64)
	if err != nil {
		log.Error().Err(err).Str("min", c.Query("min")).Msg("Invalid minimum rating format")
		responses.RespondBadRequest(c, err, "Invalid minimum rating")
		return
	}

	maxRating, err := strconv.ParseFloat(c.Query("max"), 64)
	if err != nil {
		log.Error().Err(err).Str("max", c.Query("max")).Msg("Invalid maximum rating format")
		responses.RespondBadRequest(c, err, "Invalid maximum rating")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Float64("minRating", minRating).
		Float64("maxRating", maxRating).
		Msg("Retrieving movies by rating range")

	movies, err := h.movieService.GetMoviesByRating(ctx, uid, minRating, maxRating)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Float64("minRating", minRating).
			Float64("maxRating", maxRating).
			Msg("Failed to retrieve movies by rating")
		responses.RespondInternalError(c, err, "Failed to retrieve movies")
		return
	}

	responses.RespondOK(c, movies, "Movies retrieved successfully")
}

func (h *MediaClientMovieHandler[T]) GetLatestMoviesByAdded(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	count, err := strconv.Atoi(c.Param("count"))
	if err != nil {
		log.Error().Err(err).Str("count", c.Param("count")).Msg("Invalid count format")
		responses.RespondBadRequest(c, err, "Invalid count")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", count).
		Msg("Retrieving latest movies by added date")

	movies, err := h.movieService.GetLatestMoviesByAdded(ctx, uid, count)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Int("count", count).
			Msg("Failed to retrieve latest movies")
		responses.RespondInternalError(c, err, "Failed to retrieve movies")
		return
	}

	responses.RespondOK(c, movies, "Movies retrieved successfully")
}

func (h *MediaClientMovieHandler[T]) GetPopularMovies(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	count, err := strconv.Atoi(c.Param("count"))
	if err != nil {
		log.Error().Err(err).Str("count", c.Param("count")).Msg("Invalid count format")
		responses.RespondBadRequest(c, err, "Invalid count")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", count).
		Msg("Retrieving popular movies")

	movies, err := h.movieService.GetPopularMovies(ctx, uid, count)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Int("count", count).
			Msg("Failed to retrieve popular movies")
		responses.RespondInternalError(c, err, "Failed to retrieve movies")
		return
	}

	responses.RespondOK(c, movies, "Movies retrieved successfully")
}

func (h *MediaClientMovieHandler[T]) GetTopRatedMovies(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	count, err := strconv.Atoi(c.Param("count"))
	if err != nil {
		log.Error().Err(err).Str("count", c.Param("count")).Msg("Invalid count format")
		responses.RespondBadRequest(c, err, "Invalid count")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", count).
		Msg("Retrieving top rated movies")

	movies, err := h.movieService.GetTopRatedMovies(ctx, uid, count)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Int("count", count).
			Msg("Failed to retrieve top rated movies")
		responses.RespondInternalError(c, err, "Failed to retrieve movies")
		return
	}

	responses.RespondOK(c, movies, "Movies retrieved successfully")
}

func (h *MediaClientMovieHandler[T]) SearchMovies(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
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
		Msg("Searching movies")

	movies, err := h.movieService.SearchMovies(ctx, uid, query)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Str("query", query).
			Msg("Failed to search movies")
		responses.RespondInternalError(c, err, "Failed to search movies")
		return
	}

	responses.RespondOK(c, movies, "Movies retrieved successfully")
}
