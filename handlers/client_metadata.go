package handlers

import (
	"net/http"
	"suasor/clients/types"
	"suasor/services"
	"suasor/types/requests"

	"github.com/gin-gonic/gin"
)

// ClientMetadataHandler handles requests for metadata clients
type ClientMetadataHandler[T types.ClientMetadataConfig] struct {
	service *services.ClientMetadataService[T]
}

// NewClientMetadataHandler creates a new ClientMetadataHandler
func NewClientMetadataHandler[T types.ClientMetadataConfig](service *services.ClientMetadataService[T]) *ClientMetadataHandler[T] {
	return &ClientMetadataHandler[T]{
		service: service,
	}
}

// Movies

// GetMovie retrieves a movie by ID
func (h *ClientMetadataHandler[T]) GetMovie(c *gin.Context) {
	var req requests.MetadataMovieRequest
	if !checkJSONBinding(c, &req) {
		return
	}

	movie, err := h.service.GetMovie(c.Request.Context(), req.ClientID, req.MovieID)
	if err != nil {
		handleServiceError(c, err, "Retrieving movie metadata", "Movie not found", "Failed to retrieve movie metadata")
		return
	}

	c.JSON(http.StatusOK, movie)
}

// SearchMovies searches for movies by query
func (h *ClientMetadataHandler[T]) SearchMovies(c *gin.Context) {
	var req requests.MetadataMovieSearchRequest
	if !checkJSONBinding(c, &req) {
		return
	}

	movies, err := h.service.SearchMovies(c.Request.Context(), req.ClientID, req.Query)
	if err != nil {
		handleServiceError(c, err, "Searching movies metadata", "", "Failed to search movies metadata")
		return
	}

	c.JSON(http.StatusOK, movies)
}

// GetMovieRecommendations gets movie recommendations based on a movie ID
func (h *ClientMetadataHandler[T]) GetMovieRecommendations(c *gin.Context) {
	var req requests.MetadataMovieRecommendationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	movies, err := h.service.GetMovieRecommendations(c.Request.Context(), req.ClientID, req.MovieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movies)
}

// GetPopularMovies gets popular movies
func (h *ClientMetadataHandler[T]) GetPopularMovies(c *gin.Context) {
	var req requests.MetadataPopularMoviesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	movies, err := h.service.GetPopularMovies(c.Request.Context(), req.ClientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movies)
}

// GetTrendingMovies gets trending movies
func (h *ClientMetadataHandler[T]) GetTrendingMovies(c *gin.Context) {
	var req requests.MetadataTrendingMoviesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	movies, err := h.service.GetTrendingMovies(c.Request.Context(), req.ClientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movies)
}

// TV Shows

// GetTVShow retrieves a TV show by ID
func (h *ClientMetadataHandler[T]) GetTVShow(c *gin.Context) {
	var req requests.MetadataTVShowRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	tvShow, err := h.service.GetTVShow(c.Request.Context(), req.ClientID, req.TVShowID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tvShow)
}

// SearchTVShows searches for TV shows by query
func (h *ClientMetadataHandler[T]) SearchTVShows(c *gin.Context) {
	var req requests.MetadataTVShowSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	tvShows, err := h.service.SearchTVShows(c.Request.Context(), req.ClientID, req.Query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tvShows)
}

// GetTVShowRecommendations gets TV show recommendations based on a TV show ID
func (h *ClientMetadataHandler[T]) GetTVShowRecommendations(c *gin.Context) {
	var req requests.MetadataTVShowRecommendationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	tvShows, err := h.service.GetTVShowRecommendations(c.Request.Context(), req.ClientID, req.TVShowID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tvShows)
}

// GetPopularTVShows gets popular TV shows
func (h *ClientMetadataHandler[T]) GetPopularTVShows(c *gin.Context) {
	var req requests.MetadataPopularTVShowsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	tvShows, err := h.service.GetPopularTVShows(c.Request.Context(), req.ClientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tvShows)
}

// GetTrendingTVShows gets trending TV shows
func (h *ClientMetadataHandler[T]) GetTrendingTVShows(c *gin.Context) {
	var req requests.MetadataTrendingTVShowsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	tvShows, err := h.service.GetTrendingTVShows(c.Request.Context(), req.ClientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tvShows)
}

// GetTVSeason retrieves a TV season by show ID and season number
func (h *ClientMetadataHandler[T]) GetTVSeason(c *gin.Context) {
	var req requests.MetadataTVSeasonRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	season, err := h.service.GetTVSeason(c.Request.Context(), req.ClientID, req.TVShowID, req.SeasonNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, season)
}

// GetTVEpisode retrieves a TV episode by show ID, season number, and episode number
func (h *ClientMetadataHandler[T]) GetTVEpisode(c *gin.Context) {
	var req requests.MetadataTVEpisodeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	episode, err := h.service.GetTVEpisode(c.Request.Context(), req.ClientID, req.TVShowID, req.SeasonNumber, req.EpisodeNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, episode)
}

// People

// GetPerson retrieves a person by ID
func (h *ClientMetadataHandler[T]) GetPerson(c *gin.Context) {
	var req requests.MetadataPersonRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	person, err := h.service.GetPerson(c.Request.Context(), req.ClientID, req.PersonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, person)
}

// SearchPeople searches for people by query
func (h *ClientMetadataHandler[T]) SearchPeople(c *gin.Context) {
	var req requests.MetadataPersonSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	people, err := h.service.SearchPeople(c.Request.Context(), req.ClientID, req.Query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, people)
}

// Collections

// GetCollection retrieves a collection by ID
func (h *ClientMetadataHandler[T]) GetCollection(c *gin.Context) {
	var req requests.MetadataCollectionRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	collection, err := h.service.GetCollection(c.Request.Context(), req.ClientID, req.CollectionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, collection)
}

// SearchCollections searches for collections by query
func (h *ClientMetadataHandler[T]) SearchCollections(c *gin.Context) {
	var req requests.MetadataCollectionSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	collections, err := h.service.SearchCollections(c.Request.Context(), req.ClientID, req.Query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, collections)
}
