package router

import (
	"suasor/handlers"
	"suasor/services"

	"github.com/gin-gonic/gin"
)

func RegisterShortenRoutes(rg *gin.RouterGroup, service services.ShortenService) {
	shortenHandlers := handlers.NewShortenHandler(service)
	shorts := rg.Group("/shorten")
	{

		shorts.POST("", shortenHandlers.Create)
		shorts.POST("lookup", shortenHandlers.GetByOriginalURL)
		shorts.PUT("/:code", shortenHandlers.Update)
		shorts.DELETE("/:code", shortenHandlers.Delete)
		shorts.GET("/:code", shortenHandlers.Redirect)

	}
}
