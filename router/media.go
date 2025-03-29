package router

import (
	"suasor/handlers"
	"suasor/services"

	"github.com/gin-gonic/gin"
)

func RegisterMediaRoutes(rg *gin.RouterGroup, movieService services.MediaMovieService) {
	movieHandlers := handlers.NewMediaMovieHandler(movieService)
	movies := rg.Group("/movies")
	{
		movies.GET("/:clientID/:movieID", movieHandlers.GetMovieByID)
		movies.GET("/genre/:genre", movieHandlers.GetMoviesByGenre)
		movies.GET("/year/:year", movieHandlers.GetMoviesByYear)
		movies.GET("/actor/:actor", movieHandlers.GetMoviesByActor)
		movies.GET("/director/:director", movieHandlers.GetMoviesByDirector)
		movies.GET("/rating", movieHandlers.GetMoviesByRating)
		movies.GET("/latest/:count", movieHandlers.GetLatestMoviesByAdded)
		movies.GET("/popular/:count", movieHandlers.GetPopularMovies)
		movies.GET("/top-rated/:count", movieHandlers.GetTopRatedMovies)
		movies.GET("/search", movieHandlers.SearchMovies)
	}
	showHandlers := handlers.NewMediaShowHandler(movieService)
	shows := rg.Group("/shows")
	{
		shows.GET("/:clientID/:showID", showHandlers.GetShowByID)
		shows.GET("/genre/:genre", showHandlers.GetShowsByGenre)
		shows.GET("/year/:year", showHandlers.GetShowsByYear)
		shows.GET("/actor/:actor", showHandlers.GetShowsByActor)
		shows.GET("/director/:director", showHandlers.GetShowsByDirector)
		shows.GET("/rating", showHandlers.GetShowsByRating)
		shows.GET("/latest/:count", showHandlers.GetLatestShowsByAdded)
		shows.GET("/popular/:count", showHandlers.GetPopularShows)
		shows.GET("/top-rated/:count", showHandlers.GetTopRatedShows)
		shows.GET("/search", showHandlers.SearchShows)
	}
	musicHandlers := handlers.NewMediaMusicHandler(movieService)
	music := rg.Group("/music")
	{
		music.GET("/:clientID/:musicID", musicHandlers.GetMusicByID)
		music.GET("/genre/:genre", musicHandlers.GetMusicsByGenre)
		music.GET("/year/:year", musicHandlers.GetMusicsByYear)
		music.GET("/artist/:artist", musicHandlers.GetMusicsByArtist)
		music.GET("/director/:director", musicHandlers.GetMusicsByDirector)
		music.GET("/rating", musicHandlers.GetMusicsByRating)
		music.GET("/latest/:count", musicHandlers.GetLatestMusicsByAdded)
		music.GET("/popular/:count", musicHandlers.GetPopularMusics)
		music.GET("/top-rated/:count", musicHandlers.GetTopRatedMusics)
		music.GET("/search", musicHandlers.SearchMusics)
	}
}
