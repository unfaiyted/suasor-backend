// router/client.go
package router

import (
	"context"
	clienttypes "suasor/clients/types"
	"suasor/di/container"
	"suasor/handlers"
	"suasor/router/middleware"
	"suasor/types/responses"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterClientRoutes configures routes for client endpoints
// These endpoints are specific to a single client instance.
func RegisterClientRoutes(ctx context.Context, r *gin.RouterGroup, c *container.Container) {
	db := container.MustGet[*gorm.DB](c)

	clientGroup := r.Group("/client/:clientID")
	clientGroup.Use(middleware.ClientTypeMiddleware(db))
	{
		clientGroup.GET("", func(g *gin.Context) {
			clientType := g.Param("clientType")
			getClientHandler(g, c, clientType).GetClient(g)
		})
		clientGroup.PUT("", func(g *gin.Context) {
			clientType := g.Param("clientType")
			getClientHandler(g, c, clientType).UpdateClient(g)
		})
		clientGroup.DELETE("", func(g *gin.Context) {
			clientType := g.Param("clientType")
			getClientHandler(g, c, clientType).DeleteClient(g)
		})

		clientGroup.GET("/:clientID/test", func(g *gin.Context) {
			clientType := g.Param("clientType")
			getClientHandler(g, c, clientType).TestConnection(g)
		})
	}

	automationHandler := container.MustGet[*handlers.ClientAutomationHandler](c)
	automationClientGroup := r.Group("client/:clientID/automation")
	{
		automationClientGroup.GET("/status", automationHandler.GetSystemStatus)
		automationClientGroup.POST("/library", automationHandler.GetLibraryItems)
		automationClientGroup.GET("/item/:itemID", automationHandler.GetMediaByID)
		automationClientGroup.POST("/item", automationHandler.AddMedia)
		automationClientGroup.PUT("/item/:itemID", automationHandler.UpdateMedia)
		automationClientGroup.DELETE("/item/:itemID", automationHandler.DeleteMedia)
		automationClientGroup.GET("/command", automationHandler.ExecuteCommand)
		automationClientGroup.GET("/calendar", automationHandler.GetCalendar)
		automationClientGroup.GET("/search", automationHandler.SearchMedia)

	}
}

func getClientHandler(g *gin.Context, c *container.Container, clientType string) handlers.ClientHandler[clienttypes.ClientConfig] {
	handlers := map[string]handlers.ClientHandler[clienttypes.ClientConfig]{
		"emby":     container.MustGet[handlers.ClientHandler[*clienttypes.EmbyConfig]](c),
		"jellyfin": container.MustGet[handlers.ClientHandler[*clienttypes.JellyfinConfig]](c),
		"plex":     container.MustGet[handlers.ClientHandler[*clienttypes.PlexConfig]](c),
		"subsonic": container.MustGet[handlers.ClientHandler[*clienttypes.SubsonicConfig]](c),
		"sonarr":   container.MustGet[handlers.ClientHandler[*clienttypes.SonarrConfig]](c),
		"radarr":   container.MustGet[handlers.ClientHandler[*clienttypes.RadarrConfig]](c),
		"lidarr":   container.MustGet[handlers.ClientHandler[*clienttypes.LidarrConfig]](c),
		"claude":   container.MustGet[handlers.ClientHandler[*clienttypes.ClaudeConfig]](c),
		"openai":   container.MustGet[handlers.ClientHandler[*clienttypes.OpenAIConfig]](c),
		"ollama":   container.MustGet[handlers.ClientHandler[*clienttypes.OllamaConfig]](c),
	}
	handler, exists := handlers[clientType]
	if !exists {
		responses.RespondInternalError(g, nil, "Client handler not found")
		return nil
	}
	return handler
}
