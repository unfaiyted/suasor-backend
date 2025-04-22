package router

import (
	"github.com/gin-gonic/gin"
	"suasor/app/container"
	"suasor/handlers"
)

// RegisterClientsRoutes sets up routes for client-related operations
func RegisterClientsRoutes(r *gin.RouterGroup, c *container.Container) {
	// Route to get all clients of all types
	clientsHandler := container.MustGet[*handlers.ClientsHandler](c)
	r.GET("/clients", clientsHandler.ListAllClients)
}

