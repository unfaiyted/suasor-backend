package router

import (
	"suasor/app"

	"github.com/gin-gonic/gin"
)

// RegisterClientsRoutes sets up routes for client-related operations
func RegisterClientsRoutes(r *gin.RouterGroup, deps *app.AppDependencies) {
	// Route to get all clients of all types
	r.GET("/clients", deps.SystemHandlers.ClientsHandler().ListAllClients)
}