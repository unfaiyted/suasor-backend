// handler/media_movie_handler.go
package handlers

import (
	"net/http"
	"strconv"
	"suasor/services"

	"github.com/gin-gonic/gin"
)

type MediaMovieHandler struct {
	movieService services.MediaMovieService
}

func NewMediaMovieHandler(movieService services.MediaMovieService) *MediaMovieHandler {
	return &MediaMovieHandler{
		movieService: movieService,
	}
}

func (h *MediaMovieHandler) RegisterRoutes(router *gin.RouterGroup) {
	movies := router.Group("/movies")
	{
		movies.GET("/:clientID/:movieID", h.GetMovieByID)
		movies.GET("/genre/:genre", h.GetMoviesByGenre)
		movies.GET("/year/:year", h.GetMoviesByYear)
		movies.GET("/actor/:actor", h.GetMoviesByActor)
		movies.GET("/director/:director", h.GetMoviesByDirector)
		movies.GET("/rating", h.GetMoviesByRating)
		movies.GET("/latest/:count", h.GetLatestMoviesByAdded)
		movies.GET("/popular/:count", h.GetPopularMovies)
		movies.GET("/top-rated/:count", h.GetTopRatedMovies)
		movies.GET("/search", h.SearchMovies)
	}
}

func (h *MediaMovieHandler) GetMovieByID(c *gin.Context) {
	userID := getUserIDFromContext(c)

	clientID, err := strconv.ParseUint(c.Param("clientID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client ID"})
		return
	}

	movieID := c.Param("movieID")

	movie, err := h.movieService.GetMovieByID(c.Request.Context(), userID, clientID, movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movie)
}

func (h *MediaMovieHandler) GetMoviesByGenre(c *gin.Context) {
	userID := getUserIDFromContext(c)
	genre := c.Param("genre")

	movies, err := h.movieService.GetMoviesByGenre(c.Request.Context(), userID, genre)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movies)
}

func (h *MediaMovieHandler) GetMoviesByYear(c *gin.Context) {
	userID := getUserIDFromContext(c)

	year, err := strconv.Atoi(c.Param("year"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year"})
		return
	}

	movies, err := h.movieService.GetMoviesByYear(c.Request.Context(), userID, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movies)
}

func (h *MediaMovieHandler) GetMoviesByActor(c *gin.Context) {
	userID := getUserIDFromContext(c)
	actor := c.Param("actor")

	movies, err := h.movieService.GetMoviesByActor(c.Request.Context(), userID, actor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movies)
}

func (h *MediaMovieHandler) GetMoviesByDirector(c *gin.Context) {
	userID := getUserIDFromContext(c)
	director := c.Param("director")

	movies, err := h.movieService.GetMoviesByDirector(c.Request.Context(), userID, director)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movies)
}

func (h *MediaMovieHandler) GetMoviesByRating(c *gin.Context) {
	userID := getUserIDFromContext(c)

	minRating, err := strconv.ParseFloat(c.Query("min"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid minimum rating"})
		return
	}

	maxRating, err := strconv.ParseFloat(c.Query("max"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid maximum rating"})
		return
	}

	movies, err := h.movieService.GetMoviesByRating(c.Request.Context(), userID, minRating, maxRating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movies)
}

func (h *MediaMovieHandler) GetLatestMoviesByAdded(c *gin.Context) {
	userID := getUserIDFromContext(c)

	count, err := strconv.Atoi(c.Param("count"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid count"})
		return
	}

	movies, err := h.movieService.GetLatestMoviesByAdded(c.Request.Context(), userID, count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movies)
}

func (h *MediaMovieHandler) GetPopularMovies(c *gin.Context) {
	userID := getUserIDFromContext(c)

	count, err := strconv.Atoi(c.Param("count"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid count"})
		return
	}

	movies, err := h.movieService.GetPopularMovies(c.Request.Context(), userID, count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movies)
}

func (h *MediaMovieHandler) GetTopRatedMovies(c *gin.Context) {
	userID := getUserIDFromContext(c)

	count, err := strconv.Atoi(c.Param("count"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid count"})
		return
	}

	movies, err := h.movieService.GetTopRatedMovies(c.Request.Context(), userID, count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movies)
}

func (h *MediaMovieHandler) SearchMovies(c *gin.Context) {
	userID := getUserIDFromContext(c)
	query := c.Query("q")

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	movies, err := h.movieService.SearchMovies(c.Request.Context(), userID, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movies)
}

// Helper function to get user ID from context
func getUserIDFromContext(c *gin.Context) uint64 {
	// This should be set by your authentication middleware
	userIDValue, exists := c.Get("userID")
	if !exists {
		return 0
	}

	userID, ok := userIDValue.(uint64)
	if !ok {
		return 0
	}

	return userID
}
