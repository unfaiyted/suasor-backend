package router

import (
	"suasor/handlers"
	"suasor/services"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup, service services.UserService) {
	userHandlers := handlers.NewUserHandler(service)
	users := rg.Group("/users")
	{
		// Public routes
		users.POST("/register", userHandlers.Register)

		// Authenticated user routes
		users.GET("/profile", userHandlers.GetProfile)
		users.PUT("/profile", userHandlers.UpdateProfile)
		users.PUT("/password", userHandlers.ChangePassword)

		// Admin routes
		users.GET("/:id", userHandlers.GetByID)
		users.PUT("/:id/role", userHandlers.ChangeRole)
		users.POST("/:id/activate", userHandlers.ActivateUser)
		users.POST("/:id/deactivate", userHandlers.DeactivateUser)
		users.DELETE("/:id", userHandlers.Delete)
	}
}
