package router // router/media_mediaItems.go

import (
	"suasor/app/container"

	"github.com/gin-gonic/gin"
)

func RegisterMediaItemRoutes(rg *gin.RouterGroup, c *container.Container) {

	mediaItems := rg.Group("/item")
	{

		// {base}/item/:mediaType/ example: /movie/
		RegisterLocalMediaItemRoutes(mediaItems, c) // Register direct media item routes (non-client specific)
	}

	clientMediaItems := rg.Group("/client")
	{
		// {base}/client/:id/item/:mediaType/ example: /client/11/movies
		RegisterClientMediaItemRoutes(clientMediaItems, c)

	}
	// userMediaItems := rg.Group("/user/item")
	// {
	// {base}/user/item/:mediaType/ example: /user/movies/
	// RegisterUserMediaItemRoutes(userMediaItems, c)
	// }

}
