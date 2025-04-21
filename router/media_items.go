package router // router/media_mediaItems.go

import (
	"suasor/app/container"

	"context"
	"github.com/gin-gonic/gin"
)

func RegisterMediaItemRoutes(ctx context.Context, rg *gin.RouterGroup, c *container.Container) {

	mediaItems := rg.Group("/item")
	{

		// {base}/item/:mediaType/ example: /item/movie/
		RegisterLocalMediaItemRoutes(ctx, mediaItems, c) // Register direct media item routes (non-client specific)
	}

	clientMediaItems := rg.Group("/client")
	{

		// {base}/client/:id/item/:mediaType/ example: /client/11/movies
		RegisterClientMediaItemRoutes(ctx, clientMediaItems, c)
		// {base}/client/:id/item/playlist/ example: /client/11/movies
		RegisterClientListRoutes(ctx, clientMediaItems, c)

	}
	// userMediaItems := rg.Group("/user/item")
	// {
	// {base}/user/item/:mediaType/ example: /user/movies/
	// RegisterUserMediaItemRoutes(userMediaItems, c)
	// }

}
