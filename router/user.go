package router

import (
	"github.com/gin-gonic/gin"
	"suasor/app"
)

func RegisterUserRoutes(rg *gin.RouterGroup, deps *app.AppDependencies) {
	userHandlers := deps.UserHandlers.UserHandler()
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
