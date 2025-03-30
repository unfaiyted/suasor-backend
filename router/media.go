package router

import (
	"suasor/handlers"
	"suasor/repository"
	"suasor/services"

	"suasor/client"
	"suasor/client/types"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterMediaClientRoutes(rg *gin.RouterGroup, db *gorm.DB) {

	clientRepo := repository.NewClientRepository[types.MediaClientConfig](db)
	clientFactory := client.NewClientFactoryService()
	movieService := services.NewMediaClientMovieService(clientRepo, *clientFactory)
	movieHandlers := handlers.NewMediaClientMovieHandler(movieService)

	clientGroup := rg.Group("/client")

	client := clientGroup.Group("/:clientID")
	{
		client.GET("/movie/:movieID", movieHandlers.GetMovieByID)
		client.GET("/genre/:genre", movieHandlers.GetMoviesByGenre)
		client.GET("/year/:year", movieHandlers.GetMoviesByYear)
		client.GET("/actor/:actor", movieHandlers.GetMoviesByActor)
		client.GET("/director/:director", movieHandlers.GetMoviesByDirector)
		client.GET("/rating", movieHandlers.GetMoviesByRating)
		client.GET("/latest/:count", movieHandlers.GetLatestMoviesByAdded)
		client.GET("/popular/:count", movieHandlers.GetPopularMovies)
		client.GET("/top-rated/:count", movieHandlers.GetTopRatedMovies)
		client.GET("/search", movieHandlers.SearchMovies)
	}
	// showHandlers := handlers.NewMediaClientSeriesHandler(movieService)
	// shows := client.Group("/shows")
	// {
	// 	shows.GET("/:clientID/:showID", showHandlers.GetShowByID)
	// 	shows.GET("/genre/:genre", showHandlers.GetShowsByGenre)
	// 	shows.GET("/year/:year", showHandlers.GetShowsByYear)
	// 	shows.GET("/actor/:actor", showHandlers.GetShowsByActor)
	// 	shows.GET("/director/:director", showHandlers.GetShowsByDirector)
	// 	shows.GET("/rating", showHandlers.GetShowsByRating)
	// 	shows.GET("/latest/:count", showHandlers.GetLatestShowsByAdded)
	// 	shows.GET("/popular/:count", showHandlers.GetPopularShows)
	// 	shows.GET("/top-rated/:count", showHandlers.GetTopRatedShows)
	// 	shows.GET("/search", showHandlers.SearchShows)
	// }
	// musicHandlers := handlers.NewMediaMusicHandler(movieService)
	// music := rg.Group("/music")
	// {
	// 	music.GET("/:clientID/:musicID", musicHandlers.GetMusicByID)
	// 	music.GET("/genre/:genre", musicHandlers.GetMusicsByGenre)
	// 	music.GET("/year/:year", musicHandlers.GetMusicsByYear)
	// 	music.GET("/artist/:artist", musicHandlers.GetMusicsByArtist)
	// 	music.GET("/director/:director", musicHandlers.GetMusicsByDirector)
	// 	music.GET("/rating", musicHandlers.GetMusicsByRating)
	// 	music.GET("/latest/:count", musicHandlers.GetLatestMusicsByAdded)
	// 	music.GET("/popular/:count", musicHandlers.GetPopularMusics)
	// 	music.GET("/top-rated/:count", musicHandlers.GetTopRatedMusics)
	// 	music.GET("/search", musicHandlers.SearchMusics)
	// }
}
