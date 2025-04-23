package router

import (
	"github.com/gin-gonic/gin"
	"suasor/di/container"
	"suasor/handlers"
)

func RegisterUserRoutes(rg *gin.RouterGroup, c *container.Container) {
	userHandlers := container.MustGet[*handlers.UserHandler](c)
	users := rg.Group("/user")
	{
		// Public routes
		users.POST("/register", userHandlers.Register)

		// Authenticated user routes
		users.GET("/profile", userHandlers.GetProfile)
		users.PUT("/profile", userHandlers.UpdateProfile)
		users.PUT("/password", userHandlers.ChangePassword)
		users.POST("/avatar", userHandlers.UploadAvatar)

		// Admin routes
		users.GET("/:userID", userHandlers.GetByID)
		users.PUT("/:userID/role", userHandlers.ChangeRole)
		users.POST("/:userID/activate", userHandlers.ActivateUser)
		users.POST("/:userID/deactivate", userHandlers.DeactivateUser)
		users.DELETE("/:userID", userHandlers.Delete)
	}

	// Avatar files are served statically from the main router
}
