// handlers/client_media_movie.go
package handlers

import (
	"strconv"
	"strings"
	"suasor/clients/media/types"
	clienttypes "suasor/clients/types"
	"suasor/services"
	models "suasor/types/models"
	"suasor/types/responses"
	"suasor/utils/logger"

	"github.com/gin-gonic/gin"
)

type ClientMovieHandler[T clienttypes.ClientMediaConfig] interface {
	// CoreMovieHandler has the basic database retreival operations
	CoreMovieHandler

	// Client-Based Movie Operations
	GetClientMovieByID(c *gin.Context)
	GetClientMovieByExternalID(c *gin.Context) //imdb, tmdb, etc.

	GetClientMoviesByGenre(c *gin.Context)
	GetClientMoviesByYear(c *gin.Context)
	GetClientMoviesByActor(c *gin.Context)
	GetClientMoviesByDirector(c *gin.Context)
	GetClientMoviesByRating(c *gin.Context)
	GetClientMoviesLatestByAdded(c *gin.Context)
	GetClientMoviesPopular(c *gin.Context)
	GetClientMoviesTopRated(c *gin.Context)

	SearchClientMovies(c *gin.Context)
}

// clientMovieHandler handles movie-related operations for media clients
type clientMovieHandler[T clienttypes.ClientMediaConfig] struct {
	CoreMovieHandler
	movieService services.ClientMovieService[T]
}

// NewclientMovieHandler creates a new media client movie handler
func NewClientMovieHandler[T clienttypes.ClientMediaConfig](
	coreHandler CoreMovieHandler,
	movieService services.ClientMovieService[T]) ClientMovieHandler[T] {
	return &clientMovieHandler[T]{
		CoreMovieHandler: coreHandler,
		movieService:     movieService,
	}
}

// GetMovieByID godoc
//
//	@Summary		Get movie by ID
//	@Description	Retrieves a specific movie from the client by ID
//	@Tags			movies, clients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			clientID		path		int															true	"Client ID"
//	@Param			clientItemID	path		string														true	"Movie ID"
//	@Success		200				{object}	responses.APIResponse[models.MediaItem[types.Movie]]	"Movies retrieved"
//	@Failure		400				{object}	responses.ErrorResponse[responses.ErrorDetails]				"Invalid client ID"
//	@Failure		401				{object}	responses.ErrorResponse[responses.ErrorDetails]				"Unauthorized"
//	@Failure		500				{object}	responses.ErrorResponse[responses.ErrorDetails]				"Server error"
//	@Router			/client/{clientID}/movie/{clientItemId} [get]
func (h *clientMovieHandler[T]) GetClientMovieByID(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting movie by ID")

	// Get authenticated user ID
	uid, ok := checkUserAccess(c)
	if !ok {
		return
	}

	// Parse client ID from URL
	clientID, err := checkItemID(c, "clientID")
	if err != nil {
		return
	}

	movieID, _ := checkClientItemID(c, "clientItemID")

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("movieID", movieID).
		Msg("Retrieving movie by ID")

	movie, err := h.movieService.GetClientMovieByItemID(ctx, clientID, movieID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Uint64("clientID", clientID).
			Str("movieID", movieID).
			Msg("Failed to retrieve movie")
		responses.RespondInternalError(c, err, "Failed to retrieve movie")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("clientID", clientID).
		Str("movieID", movieID).
		Msg("Movie retrieved successfully")
	responses.RespondOK(c, movie, "Movie retrieved successfully")
}

// GetMoviesByGenre godoc
//
//	@Summary		Get movies by genre
//	@Description	Retrieves movies from all connected clients that match the specified genre
//	@Tags			movies, clients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			genre	path		string														true	"Genre name"
//	@Success		200		{object}	responses.APIResponse[responses.MediaItemList[types.Movie]]	"Movies retrieved"
//	@Failure		401		{object}	responses.ErrorResponse[responses.ErrorDetails]				"Unauthorized"
//	@Failure		500		{object}	responses.ErrorResponse[responses.ErrorDetails]				"Server error"
//	@Router			/client/{clientID}/movie/genre/{genre} [get]
func (h *clientMovieHandler[T]) GetClientMoviesByGenre(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting movies by genre")

	// Get authenticated user ID
	uid, ok := checkUserAccess(c)
	if !ok {
		return
	}
	genre := c.Param("genre")

	log.Info().
		Uint64("userID", uid).
		Str("genre", genre).
		Msg("Retrieving movies by genre")

	movies, err := h.movieService.GetClientMoviesByGenre(ctx, uid, genre)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Str("genre", genre).
			Msg("Failed to retrieve movies by genre")
		responses.RespondInternalError(c, err, "Failed to retrieve movies")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Str("genre", genre).
		Int("count", len(movies)).
		Msg("Movies retrieved successfully")
	responses.RespondMediaItemListOK(c, movies, "Movies retrieved successfully")
}

// GetMoviesByYear godoc
//
//	@Summary		Get movies by release year
//	@Description	Retrieves movies from all connected clients that were released in the specified year
//	@Tags			movies, clients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			year	path		int															true	"Release year"
//	@Success		200		{object}	responses.APIResponse[responses.MediaItemList[types.Movie]]	"Movies retrieved"
//	@Failure		400		{object}	responses.ErrorResponse[responses.ErrorDetails]				"Invalid year"
//	@Failure		401		{object}	responses.ErrorResponse[responses.ErrorDetails]				"Unauthorized"
//	@Failure		500		{object}	responses.ErrorResponse[responses.ErrorDetails]				"Server error"
//	@Router			/client/{clientID}/movie/year/{year} [get]
func (h *clientMovieHandler[T]) GetClientMoviesByYear(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting movies by year")

	// Get authenticated user ID
	uid, ok := checkUserAccess(c)
	if !ok {
		return
	}

	year, ok := checkYear(c, "year")
	if !ok {
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("year", year).
		Msg("Retrieving movies by year")

	movies, err := h.movieService.GetClientMoviesByYear(ctx, uid, year)
	if handleServiceError(c, err,
		"Failed to retrieve movies by year",
		"No movies found for this year",
		"Failed to retrieve movies") {
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("year", year).
		Int("count", len(movies)).
		Msg("Movies retrieved successfully")
	responses.RespondMediaItemListOK(c, movies, "Movies retrieved successfully")
}

// GetMoviesByActor godoc
//
//	@Summary		Get movies by actor
//	@Description	Retrieves movies from all connected clients featuring the specified actor
//	@Tags			movies, clients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			actor	path		string														true	"Actor name"
//	@Success		200		{object}	responses.APIResponse[responses.MediaItemList[types.Movie]]	"Movies retrieved"
//	@Failure		401		{object}	responses.ErrorResponse[responses.ErrorDetails]				"Unauthorized"
//	@Failure		500		{object}	responses.ErrorResponse[responses.ErrorDetails]				"Server error"
//	@Router			/client/{clientID}/movie/actor/{actor} [get]
func (h *clientMovieHandler[T]) GetClientMoviesByActor(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting movies by actor")

	// Get authenticated user ID
	uid, ok := checkUserAccess(c)
	if !ok {
		return
	}

	actor, ok := checkRequiredStringParam(c, "actor", "Actor name is required")
	if !ok {
		return
	}

	log.Info().
		Uint64("userID", uid).
		Str("actor", actor).
		Msg("Retrieving movies by actor")

	movies, err := h.movieService.GetClientMoviesByActor(ctx, uid, actor)
	if handleServiceError(c, err,
		"Failed to retrieve movies by actor",
		"No movies found for this actor",
		"Failed to retrieve movies") {
		return
	}

	log.Info().
		Uint64("userID", uid).
		Str("actor", actor).
		Int("count", len(movies)).
		Msg("Movies retrieved successfully")
	responses.RespondMediaItemListOK(c, movies, "Movies retrieved successfully")
}

// GetMoviesByDirector godoc
//
//	@Summary		Get movies by director
//	@Description	Retrieves movies from all connected clients directed by the specified director
//	@Tags			movies, clients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			director	path		string														true	"Director name"
//	@Success		200			{object}	responses.APIResponse[responses.MediaItemList[types.Movie]]	"Movies retrieved"
//	@Failure		401			{object}	responses.ErrorResponse[responses.ErrorDetails]				"Unauthorized"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]				"Server error"
//	@Router			/client/{clientID}/movie/director/{director} [get]
func (h *clientMovieHandler[T]) GetClientMoviesByDirector(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting movies by director")

	// Get authenticated user ID
	uid, ok := checkUserAccess(c)
	if !ok {
		return
	}

	director, ok := checkRequiredStringParam(c, "director", "Director name is required")
	if !ok {
		return
	}

	log.Info().
		Uint64("userID", uid).
		Str("director", director).
		Msg("Retrieving movies by director")

	movies, err := h.movieService.GetClientMoviesByDirector(ctx, uid, director)
	if handleServiceError(c, err,
		"Failed to retrieve movies by director",
		"No movies found for this director",
		"Failed to retrieve movies") {
		return
	}

	log.Info().
		Uint64("userID", uid).
		Str("director", director).
		Int("count", len(movies)).
		Msg("Movies retrieved successfully")
	responses.RespondMediaItemListOK(c, movies, "Movies retrieved successfully")
}

// GetMoviesByRating godoc
//
//	@Summary		Get movies by rating range
//	@Description	Retrieves movies from all connected clients with ratings in the specified range
//	@Tags			movies, clients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			min	query		number														true	"Minimum rating (e.g. 7.5)"
//	@Param			max	query		number														true	"Maximum rating (e.g. 10.0)"
//	@Success		200	{object}	responses.APIResponse[responses.MediaItemList[types.Movie]]	"Movies retrieved"
//	@Failure		400	{object}	responses.ErrorResponse[responses.ErrorDetails]				"Invalid rating format"
//	@Failure		401	{object}	responses.ErrorResponse[responses.ErrorDetails]				"Unauthorized"
//	@Failure		500	{object}	responses.ErrorResponse[responses.ErrorDetails]				"Server error"
//	@Router			/client/{clientID}/movie/rating [get]
func (h *clientMovieHandler[T]) GetClientMoviesByRating(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting movies by rating")

	// Get authenticated user ID
	uid, ok := checkUserAccess(c)
	if !ok {
		return
	}

	minRatingStr, ok := checkRequiredQueryParam(c, "min", "Minimum rating is required")
	if !ok {
		return
	}

	maxRatingStr, ok := checkRequiredQueryParam(c, "max", "Maximum rating is required")
	if !ok {
		return
	}

	minRating, err := strconv.ParseFloat(minRatingStr, 64)
	if err != nil {
		log.Error().Err(err).Str("min", minRatingStr).Msg("Invalid minimum rating format")
		responses.RespondBadRequest(c, err, "Invalid minimum rating")
		return
	}

	maxRating, err := strconv.ParseFloat(maxRatingStr, 64)
	if err != nil {
		log.Error().Err(err).Str("max", maxRatingStr).Msg("Invalid maximum rating format")
		responses.RespondBadRequest(c, err, "Invalid maximum rating")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Float64("minRating", minRating).
		Float64("maxRating", maxRating).
		Msg("Retrieving movies by rating range")

	movies, err := h.movieService.GetClientMoviesByRating(ctx, uid, minRating, maxRating)
	if handleServiceError(c, err,
		"Failed to retrieve movies by rating",
		"No movies found in this rating range",
		"Failed to retrieve movies") {
		return
	}

	log.Info().
		Uint64("userID", uid).
		Float64("minRating", minRating).
		Float64("maxRating", maxRating).
		Int("count", len(movies)).
		Msg("Movies retrieved successfully")
	responses.RespondMediaItemListOK(c, movies, "Movies retrieved successfully")
}

// GetLatestMoviesByAdded godoc
//
//	@Summary		Get latest added movies
//	@Description	Retrieves the most recently added movies from all connected clients
//	@Tags			movies, clients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			count	path		int															true	"Number of movies to retrieve"
//	@Success		200		{object}	responses.APIResponse[responses.MediaItemList[types.Movie]]	"Movies retrieved"
//	@Failure		400		{object}	responses.ErrorResponse[responses.ErrorDetails]				"Invalid count format"
//	@Failure		401		{object}	responses.ErrorResponse[responses.ErrorDetails]				"Unauthorized"
//	@Failure		500		{object}	responses.ErrorResponse[responses.ErrorDetails]				"Server error"
//	@Router			/client/{clientID}/movie/latest/{count} [get]
func (h *clientMovieHandler[T]) GetClientMoviesLatestByAdded(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting latest movies by added date")

	// Get authenticated user ID
	uid, ok := checkUserAccess(c)
	if !ok {
		return
	}

	countStr, ok := checkRequiredStringParam(c, "count", "Count is required")
	if !ok {
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		log.Error().Err(err).Str("count", countStr).Msg("Invalid count format")
		responses.RespondBadRequest(c, err, "Invalid count")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", count).
		Msg("Retrieving latest movies by added date")

	movies, err := h.movieService.GetClientMoviesLatestByAdded(ctx, uid, count)
	if handleServiceError(c, err,
		"Failed to retrieve latest movies",
		"No recent movies found",
		"Failed to retrieve movies") {
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", count).
		Int("moviesReturned", len(movies)).
		Msg("Latest movies retrieved successfully")
	responses.RespondMediaItemListOK(c, movies, "Movies retrieved successfully")
}

// GetPopularMovies godoc
//
//	@Summary		Get popular movies
//	@Description	Retrieves the most popular movies from all connected clients
//	@Tags			movies, clients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			count	path		int															true	"Number of movies to retrieve"
//	@Success		200		{object}	responses.APIResponse[responses.MediaItemList[types.Movie]]	"Movies retrieved"
//	@Failure		400		{object}	responses.ErrorResponse[responses.ErrorDetails]				"Invalid count format"
//	@Failure		401		{object}	responses.ErrorResponse[responses.ErrorDetails]				"Unauthorized"
//	@Failure		500		{object}	responses.ErrorResponse[responses.ErrorDetails]				"Server error"
//	@Router			/client/{clientID}/movie/popular/{count} [get]
func (h *clientMovieHandler[T]) GetClientMoviesPopular(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting popular movies")

	// Get authenticated user ID
	uid, ok := checkUserAccess(c)
	if !ok {
		return
	}

	countStr, ok := checkRequiredStringParam(c, "count", "Count is required")
	if !ok {
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		log.Error().Err(err).Str("count", countStr).Msg("Invalid count format")
		responses.RespondBadRequest(c, err, "Invalid count")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", count).
		Msg("Retrieving popular movies")

	movies, err := h.movieService.GetClientPopularMovies(ctx, uid, count)
	if handleServiceError(c, err,
		"Failed to retrieve popular movies",
		"No popular movies found",
		"Failed to retrieve movies") {
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", count).
		Int("moviesReturned", len(movies)).
		Msg("Popular movies retrieved successfully")
	responses.RespondMediaItemListOK(c, movies, "Movies retrieved successfully")
}

// GetTopRatedMovies godoc
//
//	@Summary		Get top rated movies
//	@Description	Retrieves the highest rated movies from all connected clients
//	@Tags			movies, clients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			count		path		int															true	"Number of movies to retrieve"
//	@Param			clientID	path		int															true	"Client ID"
//	@Success		200			{object}	responses.APIResponse[responses.MediaItemList[types.Movie]]	"Movies retrieved"
//	@Failure		400			{object}	responses.ErrorResponse[responses.ErrorDetails]				"Invalid count format"
//	@Failure		401			{object}	responses.ErrorResponse[responses.ErrorDetails]				"Unauthorized"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]				"Server error"
//	@Router			/client/{clientID}/movie/top-rated/{count} [get]
func (h *clientMovieHandler[T]) GetClientMoviesTopRated(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Getting top rated movies")

	// Get authenticated user ID
	uid, ok := checkUserAccess(c)
	if !ok {
		return
	}

	countStr, ok := checkRequiredStringParam(c, "count", "Count is required")
	if !ok {
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		log.Error().Err(err).Str("count", countStr).Msg("Invalid count format")
		responses.RespondBadRequest(c, err, "Invalid count")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", count).
		Msg("Retrieving top rated movies")

	movies, err := h.movieService.GetClientTopRatedMovies(ctx, uid, count)
	if handleServiceError(c, err,
		"Failed to retrieve top rated movies",
		"No top rated movies found",
		"Failed to retrieve movies") {
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", count).
		Int("moviesReturned", len(movies)).
		Msg("Top rated movies retrieved successfully")
	responses.RespondMediaItemListOK(c, movies, "Movies retrieved successfully")
}

// SearchMovies godoc
//
//	@Summary		Search for movies
//	@Description	Searches for movies across all connected clients matching the query
//	@Tags			movies, clients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			q			query		string														true	"Search query"
//	@Param			clientID	path		int															true	"Client ID"
//	@Success		200			{object}	responses.APIResponse[responses.MediaItemList[types.Movie]]	"Movies retrieved"
//	@Failure		400			{object}	responses.ErrorResponse[responses.ErrorDetails]				"Missing search query"
//	@Failure		401			{object}	responses.ErrorResponse[responses.ErrorDetails]				"Unauthorized"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]				"Server error"
//	@Router			/client/{clientID}/movie/search [get]
func (h *clientMovieHandler[T]) SearchClientMovies(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Searching movies")

	// Get authenticated user ID
	uid, ok := checkUserAccess(c)
	if !ok {
		return
	}

	query, ok := checkRequiredQueryParam(c, "q", "Search query is required")
	if !ok {
		return
	}

	log.Info().
		Uint64("userID", uid).
		Str("query", query).
		Msg("Searching movies")

	options := types.QueryOptions{
		Query: query,
	}

	movies, err := h.movieService.SearchClientMovies(ctx, uid, &options)
	if handleServiceError(c, err,
		"Failed to search movies",
		"No movies found matching the search query",
		"Failed to search movies") {
		return
	}

	log.Info().
		Uint64("userID", uid).
		Str("query", query).
		Int("resultsCount", len(movies)).
		Msg("Movie search completed successfully")
	responses.RespondMediaItemListOK(c, movies, "Movies retrieved successfully")
}

// GetMovieByExternalID godoc
//
//	@Summary		Get movie by external ID
//	@Description	Retrieves a movie from all connected clients by external ID
//	@Tags			movies, clients
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			source		path		int															true	"Source"
//	@Param			externalID	path		string														true	"External ID"
//	@Param			clientID	path		int															true	"Client ID"
//	@Success		200			{object}	responses.APIResponse[models.MediaItem[types.Movie]]	"Movies retrieved"
//	@Failure		400			{object}	responses.ErrorResponse[responses.ErrorDetails]				"Invalid client ID"
//	@Failure		401			{object}	responses.ErrorResponse[responses.ErrorDetails]				"Unauthorized"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]				"Server error"
//	@Router			/client/{clientID}/movie/external/{source}/{externalID} [get]
func (h *clientMovieHandler[T]) GetClientMovieByExternalID(c *gin.Context) {

}

// GetClientByActor godoc
//
//	@Summary		Get movies by actor
//	@Description	Retrieves movies featuring a specific actor
//	@Tags			movies, clients
//	@Accept			json
//	@Produce		json
//	@Param			actor	path		string														true	"Actor name"
//	@Param			limit	query		int															false	"Maximum number of movies to return (default 20)"
//	@Success		200		{object}	responses.APIResponse[responses.MediaItemList[types.Movie]]	"Movies retrieved successfully"
//	@Failure		400		{object}	responses.ErrorResponse[any]								"Invalid request"
//	@Failure		404		{object}	responses.ErrorResponse[any]								"Movie not found"
//	@Failure		500		{object}	responses.ErrorResponse[any]								"Server error"
//
// Note: This functionality is redundant, covered by GetClientMoviesByActor
func (h *clientMovieHandler[T]) GetClientByActor(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

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
	clientIDStr := c.Param("clientID")
	clientID, err := strconv.ParseUint(clientIDStr, 10, 64)
	if err != nil {
		log.Error().Err(err).Str("clientID", clientIDStr).Msg("Invalid client ID format")
		responses.RespondBadRequest(c, err, "Invalid client ID")
		return
	}

	log.Debug().
		Str("actor", actor).
		Int("limit", limit).
		Msg("Getting movies by actor")

	// Create query options
	options := types.QueryOptions{
		Actor: actor,
		Limit: limit,
	}

	// Search movies by actor
	movies, err := h.movieService.SearchClientMovies(ctx, clientID, &options)
	if err != nil {
		log.Error().Err(err).
			Str("actor", actor).
			Msg("Failed to get movies by actor")
		responses.RespondInternalError(c, err, "Failed to get movies")
		return
	}

	// Filter for items with the specified actor
	var filtered []*models.MediaItem[*types.Movie]
	for _, movie := range movies {
		if strings.EqualFold(movie.Data.Credits.GetCast()[0].Name, actor) {
			filtered = append(filtered, movie)
		}

		if len(filtered) >= limit {
			break
		}
	}

	log.Info().
		Str("actor", actor).
		Int("count", len(filtered)).
		Msg("Movies by actor retrieved successfully")
	responses.RespondMediaItemListOK(c, filtered, "Movies retrieved successfully")
}
