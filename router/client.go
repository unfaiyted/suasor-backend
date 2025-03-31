// router/client.go
package router

import (
	"fmt"
	"suasor/app"
	"suasor/types/responses"

	"github.com/gin-gonic/gin"
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
}

// SetupClientRoutes configures routes for client endpoints
func RegisterClientRoutes(r *gin.RouterGroup, deps *app.AppDependencies) {

	// Create a map of client type to handler using the interface
	handlerMap := map[string]ClientHandlerInterface{
		"emby":     deps.ClientHandlers.EmbyHandler(),
		"jellyfin": deps.ClientHandlers.JellyfinHandler(),
		"subsonic": deps.ClientHandlers.SubsonicHandler(),
		"plex":     deps.ClientHandlers.PlexHandler(),

		"sonarr": deps.ClientHandlers.SonarrHandler(),
		"radarr": deps.ClientHandlers.RadarrHandler(),
		"lidarr": deps.ClientHandlers.LidarrHandler(),
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

	clients := r.Group("/client")

	// clients.GET("", hander.GetAllClients)

	client := clients.Group("/:clientType")
	{
		client.POST("", func(c *gin.Context) {
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

		client.GET("/:id/test", func(c *gin.Context) {
			if handler := getHandler(c); handler != nil {
				handler.TestConnection(c)
			}
		})
	}
}
