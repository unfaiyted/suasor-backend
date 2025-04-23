package router

import (
	"github.com/gin-gonic/gin"
	"suasor/di/container"
	"suasor/handlers"
)

// RegisterClientsRoutes sets up routes for client-related operations
func RegisterClientsRoutes(r *gin.RouterGroup, c *container.Container) {
	// Route to get all clients of all types
	clientsHandler := container.MustGet[*handlers.ClientsHandler](c)

	clientGroup := r.Group("/clients")
	{
		clientGroup.GET("", clientsHandler.GetAllClients)
		clientGroup.GET("/:clientType", clientsHandler.GetClientsByType)
		clientGroup.POST("/:clientType/test", clientsHandler.TestNewConnection)
	}
}
