package router // router/media_mediaItems.go

import (
	"suasor/di/container"

	"context"
	"github.com/gin-gonic/gin"
)

func RegisterMediaItemRoutes(ctx context.Context, rg *gin.RouterGroup, c *container.Container) {

	mediaItems := rg.Group("/media")
	{

		// {base}/media/:mediaType/ example: /media/movie/
		RegisterLocalMediaItemRoutes(ctx, mediaItems, c) // Register direct media item routes (non-client specific)
	}

	clientMediaItems := rg.Group("/client")
	{

		// {base}/client/:id/:mediaType/ example: /client/11/movies
		RegisterClientMediaItemRoutes(ctx, clientMediaItems, c)
		// {base}/client/:id/playlist/ example: /client/11/playlist
		RegisterClientListRoutes(ctx, clientMediaItems, c)
	}

}
