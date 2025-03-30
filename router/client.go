// router/client.go
package router

import (
	"fmt"
	"suasor/client/types"
	"suasor/handlers"
	"suasor/services"
	"suasor/types/responses"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ClientHandlerInterface defines the common operations for all client handlers
type ClientHandlerInterface interface {
	CreateClient(c *gin.Context)
	GetAllClients(c *gin.Context)
	GetClient(c *gin.Context)
	UpdateClient(c *gin.Context)
	DeleteClient(c *gin.Context)
	TestConnection(c *gin.Context)
}

// SetupClientRoutes configures routes for client endpoints
func RegisterClientRoutes(r *gin.RouterGroup, db *gorm.DB) {
	embyService := services.NewClientService[*types.EmbyConfig](db)
	jellyfinService := services.NewClientService[*types.JellyfinConfig](db)
	subsonicService := services.NewClientService[*types.SubsonicConfig](db)
	plexService := services.NewClientService[*types.PlexConfig](db)

	sonarrService := services.NewClientService[*types.SonarrConfig](db)
	lidarrService := services.NewClientService[*types.LidarrConfig](db)
	radarrService := services.NewClientService[*types.RadarrConfig](db)

	// Initialize all handlers
	embyClientHandler := handlers.NewClientHandler[*types.EmbyConfig](embyService)
	jellyfinClientHandler := handlers.NewClientHandler[*types.JellyfinConfig](jellyfinService)
	lidarrClientHandler := handlers.NewClientHandler[*types.LidarrConfig](lidarrService)
	subsonicClientHandler := handlers.NewClientHandler[*types.SubsonicConfig](subsonicService)
	plexClientHandler := handlers.NewClientHandler[*types.PlexConfig](plexService)
	sonarrClientHandler := handlers.NewClientHandler[*types.SonarrConfig](sonarrService)
	radarrClientHandler := handlers.NewClientHandler[*types.RadarrConfig](radarrService)

	// Create a map of client type to handler using the interface
	handlerMap := map[string]ClientHandlerInterface{
		"emby":     embyClientHandler,
		"jellyfin": jellyfinClientHandler,
		"subsonic": subsonicClientHandler,
		"plex":     plexClientHandler,

		"sonarr": sonarrClientHandler,
		"radarr": radarrClientHandler,
		"lidarr": lidarrClientHandler,
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
	client := clients.Group("/:clientType")
	{
		client.POST("", func(c *gin.Context) {
			if handler := getHandler(c); handler != nil {
				handler.CreateClient(c)
			}
		})

		client.GET("", func(c *gin.Context) {
			if handler := getHandler(c); handler != nil {
				handler.GetAllClients(c)
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
				handler.TestConnection(c)
			}
		})
	}
}
