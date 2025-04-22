package router

import (
	"github.com/gin-gonic/gin"
	"suasor/app/container"
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
		users.GET("/:id", userHandlers.GetByID)
		users.PUT("/:id/role", userHandlers.ChangeRole)
		users.POST("/:id/activate", userHandlers.ActivateUser)
		users.POST("/:id/deactivate", userHandlers.DeactivateUser)
		users.DELETE("/:id", userHandlers.Delete)
	}

	// Avatar files are served statically from the main router
}
