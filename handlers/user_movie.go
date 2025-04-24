// handlers/user_movie.go
package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	mediatypes "suasor/clients/media/types"
	"suasor/services"
	"suasor/types/models"
	"suasor/types/responses"
	"suasor/utils/logger"
)

// UserMovieHandler handles user-specific operations for movies
type UserMovieHandler struct {
	userMovieService services.UserMediaItemService[*mediatypes.Movie]
}

// NewUserMovieHandler creates a new user movie handler
func NewUserMovieHandler(
	userMovieService services.UserMediaItemService[*mediatypes.Movie],
) *UserMovieHandler {
	return &UserMovieHandler{
		userMovieService: userMovieService,
	}
}

// GetFavoriteMovies godoc
// @Summary Get user favorite movies
// @Description Retrieves movies that a user has marked as favorites
// @Tags movies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Maximum number of movies to return (default 20)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Movie]] "Movies retrieved successfully"
// @Failure 401 {object} responses.ErrorResponse[responses.ErrorDetails] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/media/movies/user/favorites [get]
func (h *UserMovieHandler) GetFavoriteMovies(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access favorites without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		limit = 20
	}

	log.Debug().
		Uint64("userID", uid).
		Int("limit", limit).
		Msg("Getting favorite movies")

	// This is a placeholder for a real implementation
	// In a real implementation, you would query for movies specifically marked as favorites
	options := mediatypes.QueryOptions{
		MediaType: mediatypes.MediaTypeMovie,
		OwnerID:   uid,
		Favorites: true,
		Limit:     limit,
	}

	movies, err := h.userMovieService.SearchUserContent(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Msg("Failed to retrieve favorite movies")
		responses.RespondInternalError(c, err, "Failed to retrieve favorite movies")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", len(movies)).
		Msg("Favorite movies retrieved successfully")
	responses.RespondOK(c, movies, "Favorite movies retrieved successfully")
}

// GetWatchedMovies godoc
// @Summary Get user watched movies
// @Description Retrieves movies that a user has watched
// @Tags movies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Maximum number of movies to return (default 20)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Movie]] "Movies retrieved successfully"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /api/v1/media/movies/user/watched [get]
func (h *UserMovieHandler) GetWatchedMovies(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access watched movies without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		limit = 20
	}

	log.Debug().
		Uint64("userID", uid).
		Int("limit", limit).
		Msg("Getting watched movies")

	// This is a placeholder for a real implementation
	// In a real implementation, you would query play history to find watched movies
	options := mediatypes.QueryOptions{
		MediaType: mediatypes.MediaTypeMovie,
		OwnerID:   uid,
		Watched:   true,
		Limit:     limit,
		Sort:      "last_watched",
		SortOrder: "desc",
	}

	movies, err := h.userMovieService.SearchUserContent(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Msg("Failed to retrieve watched movies")
		responses.RespondInternalError(c, err, "Failed to retrieve watched movies")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", len(movies)).
		Msg("Watched movies retrieved successfully")
	responses.RespondOK(c, movies, "Watched movies retrieved successfully")
}

// GetWatchlistMovies godoc
// @Summary Get movies in user watchlist
// @Description Retrieves movies that a user has added to their watchlist
// @Tags movies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Maximum number of movies to return (default 20)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Movie]] "Movies retrieved successfully"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /api/v1/media/movies/user/watchlist [get]
func (h *UserMovieHandler) GetWatchlistMovies(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access watchlist without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		limit = 20
	}

	log.Debug().
		Uint64("userID", uid).
		Int("limit", limit).
		Msg("Getting watchlist movies")

	// This is a placeholder for a real implementation
	// In a real implementation, you would query for movies specifically marked for watchlist
	options := mediatypes.QueryOptions{
		MediaType: mediatypes.MediaTypeMovie,
		OwnerID:   uid,
		Watchlist: true,
		Limit:     limit,
	}

	movies, err := h.userMovieService.SearchUserContent(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Msg("Failed to retrieve watchlist movies")
		responses.RespondInternalError(c, err, "Failed to retrieve watchlist movies")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", len(movies)).
		Msg("Watchlist movies retrieved successfully")
	responses.RespondOK(c, movies, "Watchlist movies retrieved successfully")
}

// GetRecommendedMovies godoc
// @Summary Get recommended movies for user
// @Description Retrieves movies recommended for the user based on their preferences and watch history
// @Tags movies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Maximum number of movies to return (default 20)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Movie]] "Movies retrieved successfully"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /api/v1/media/movies/user/recommended [get]
func (h *UserMovieHandler) GetRecommendedMovies(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	userID, exists := c.Get("userID")
	if !exists {
		log.Warn().Msg("Attempt to access recommendations without authentication")
		responses.RespondUnauthorized(c, nil, "Authentication required")
		return
	}

	uid := userID.(uint64)

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		limit = 20
	}

	log.Debug().
		Uint64("userID", uid).
		Int("limit", limit).
		Msg("Getting recommended movies")

	// This is a placeholder for a real implementation
	// In a real implementation, you would use a recommendation service
	// For now, we can return a basic set of movies
	options := mediatypes.QueryOptions{
		MediaType: mediatypes.MediaTypeMovie,
		OwnerID:   uid,
		Limit:     limit,
		Sort:      "rating",
		SortOrder: "desc",
	}

	movies, err := h.userMovieService.Search(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", uid).
			Msg("Failed to retrieve recommended movies")
		responses.RespondInternalError(c, err, "Failed to retrieve recommended movies")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Int("count", len(movies)).
		Msg("Recommended movies retrieved successfully")
	responses.RespondOK(c, movies, "Recommended movies retrieved successfully")
}

// UpdateMovie godoc
// @Summary Update user data for a movie
// @Description Updates user-specific data for a movie (favorite, watched status, rating, etc.)
// @Tags movies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param itemID path int true "Movie ID"
// @Param data body requests.UserMediaItemDataUpdateRequest true "Updated user data"
// @Success 200 {object} responses.APIResponse[models.MediaItem[mediatypes.Movie]] "Movie updated successfully"
// @Failure 400 {object} responses.ErrorResponse[any] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[any] "Unauthorized"
// @Failure 404 {object} responses.ErrorResponse[any] "Movie not found"
// @Failure 500 {object} responses.ErrorResponse[any] "Server error"
// @Router /api/v1/media/movie/{itemID} [patch]
func (h *UserMovieHandler) UpdateMovie(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get authenticated user ID
	uid, _ := checkUserAccess(c)
	movieID, _ := checkItemID(c, "itemID")

	// Parse request body
	var userData models.UserMediaItemData[mediatypes.Movie]
	if err := c.ShouldBindJSON(&userData); err != nil {
		log.Warn().Err(err).Msg("Invalid request body")
		responses.RespondBadRequest(c, err, "Invalid request body")
		return
	}

	log.Debug().
		Uint64("userID", uid).
		Uint64("movieID", movieID).
		Interface("userData", userData).
		Msg("Updating movie user data")

	// Get the existing movie first
	movie, err := h.userMovieService.GetByID(ctx, movieID)
	if err != nil {
		log.Error().Err(err).
			Uint64("movieID", movieID).
			Msg("Failed to retrieve movie")
		responses.RespondNotFound(c, err, "Movie not found")
		return
	}

	// Update user data
	// TODO: Go and save user data collected
	// movie.UserData = userData

	// Update the movie
	updatedMovie, err := h.userMovieService.Update(ctx, movie)
	if err != nil {
		log.Error().Err(err).
			Uint64("movieID", movieID).
			Msg("Failed to update movie")
		responses.RespondInternalError(c, err, "Failed to update movie")
		return
	}

	log.Info().
		Uint64("userID", uid).
		Uint64("movieID", movieID).
		Msg("Movie user data updated successfully")
	responses.RespondOK(c, updatedMovie, "Movie updated successfully")
}
