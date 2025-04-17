// router/client.go
package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"suasor/app/container"
	"suasor/app/handlers"
	"suasor/types/responses"
	"suasor/utils"
)

// ClientHandlerInterface defines the common operations for all client handlers
type ClientHandlerInterface interface {
	CreateClient(c *gin.Context)
	GetAllClients(c *gin.Context)
	GetClientsByType(c *gin.Context)
	GetClient(c *gin.Context)
	UpdateClient(c *gin.Context)
	DeleteClient(c *gin.Context)
	TestConnection(c *gin.Context)
	TestNewConnection(c *gin.Context)
}

// SetupClientRoutes configures routes for client endpoints
func RegisterClientRoutes(r *gin.RouterGroup, c *container.Container) {

	clientHandlers := container.MustGet[handlers.ClientHandlers](c)

	// Create a map of client type to handler using the interface
	handlerMap := map[string]ClientHandlerInterface{
		"emby":     clientHandlers.EmbyClientHandler(),
		"jellyfin": clientHandlers.JellyfinClientHandler(),
		"subsonic": clientHandlers.SubsonicClientHandler(),
		"plex":     clientHandlers.PlexClientHandler(),

		"sonarr": clientHandlers.SonarrClientHandler(),
		"radarr": clientHandlers.RadarrClientHandler(),
		"lidarr": clientHandlers.LidarrClientHandler(),

		"claude": clientHandlers.ClaudeClientHandler(),
	}

	// Helper function to get the appropriate handler
	getHandler := func(c *gin.Context) ClientHandlerInterface {
		clientType := c.Param("clientType")
		handler, exists := handlerMap[clientType]
		if !exists {
			err := fmt.Errorf("unsupported client type: %s", clientType)
			responses.RespondBadRequest(c, err, "Unsupported client type")
			return nil
		}
		return handler
	}

	clientGroup := r.Group("/client")

	client := clientGroup.Group("/:clientType")
	{
		client.POST("", func(c *gin.Context) {
			log := utils.LoggerFromContext(c.Request.Context())
			log.Info().Msg("Creating new media client")
			if handler := getHandler(c); handler != nil {
				handler.CreateClient(c)
			}
		})

		client.GET("", func(c *gin.Context) {
			if handler := getHandler(c); handler != nil {
				handler.GetClientsByType(c)
			}
		})

		client.GET("/:id", func(c *gin.Context) {
			if handler := getHandler(c); handler != nil {
				handler.GetClient(c)
			}
		})

		client.PUT("/:id", func(c *gin.Context) {
			if handler := getHandler(c); handler != nil {
				handler.UpdateClient(c)
			}
		})

		client.DELETE("/:id", func(c *gin.Context) {
			if handler := getHandler(c); handler != nil {
				handler.DeleteClient(c)
			}
		})
		client.POST("/test", func(c *gin.Context) {
			if handler := getHandler(c); handler != nil {
				handler.TestNewConnection(c)
			}
		})

		client.GET("/:id/test", func(c *gin.Context) {
			if handler := getHandler(c); handler != nil {
				handler.TestConnection(c)
			}
		})
	}
}
