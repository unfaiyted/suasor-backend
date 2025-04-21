package router // router/media_mediaItems.go
//
// import (
// 	"suasor/app/container"
//
// 	"github.com/gin-gonic/gin"
// )
//
// func RegisterMediaRoutes(rg *gin.RouterGroup, c *container.Container) {
//
// 	mediaItems := rg.Group("/media")
// 	{
// 		// {base}/media/series/
// 		// Client-specific: {base}/clients/:clientType/:clientID/media/series/
// 		RegisterSeriesRoutes(mediaItems, c)
//
// 		// {base}/media/movies/
// 		// Client-specific: {base}/clients/:clientType/:clientID/media/movies/
// 		RegisterMovieRoutes(mediaItems, c)
//
// 		// Base: {base}/media/music/
// 		// Client-specific: {base}/clients/:clientType/:clientID/media/music/
// 		RegisterMusicRoutes(mediaItems, c)
// 	}
//
// }
