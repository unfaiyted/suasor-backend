// router/router.go
package router

import (
	"context"
	"suasor/app/container"
	"suasor/router/middleware"
	"suasor/services"
	"suasor/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(ctx context.Context, c *container.Container) *gin.Engine {
	r := gin.Default()
	log := utils.LoggerFromContext(ctx)

	configService := container.MustGet[services.ConfigService](c)
	appConfig := configService.GetConfig()

	// CORS config
	config := cors.DefaultConfig()
	config.AllowOrigins = appConfig.Auth.AllowedOrigins

	log.Info().
		Strs("AllowedOrigins", config.AllowOrigins).
		Msg("Allowed Origins set.")

	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Authorization", "Content-Type"}
	r.Use(cors.New(config))

	// Setup API v1 routes
	v1 := r.Group("/api/v1")

	// TODO: should I fix this? It doesent technically need a repo, but ti does interact with the database?
	healthService := container.MustGet[services.HealthService](c)
	authService := container.MustGet[services.AuthService](c)

	RegisterHealthRoutes(v1, healthService)
	RegisterAuthRoutes(v1, authService)

	// Serve static avatar files
	r.Static("/uploads/avatars", appConfig.App.AvatarPath)

	// Protected Routes
	authenticated := v1.Group("")
	authenticated.Use(middleware.VerifyToken(authService))
	{
		// Register all routes
		RegisterUserRoutes(authenticated, c)
		RegisterUserConfigRoutes(authenticated, c)
		RegisterMediaItemRoutes(authenticated, c)
		RegisterClientMediaRoutes(authenticated, c)
		RegisterLocalMediaItemRoutes(authenticated, c)    // Register direct media item routes (non-client specific)
		RegisterUserMediaItemDataRoutes(authenticated, c) // Register media play history routes
		RegisterMetadataRoutes(authenticated)             // Register metadata routes
		RegisterAIRoutes(authenticated, c)                // Register AI routes
		RegisterClientsRoutes(authenticated, c)           // Register all clients route
		RegisterJobRoutes(authenticated, c)               // Register job routes
		RegisterRecommendationRoutes(authenticated, c)    // Register recommendation routes
		RegisterSearchRoutes(authenticated, c)            // Register search routes
	}

	//Admin Routes
	adminRoutes := v1.Group("/admin")
	adminRoutes.Use(middleware.VerifyToken(authService), middleware.RequireRole("admin"))
	{
		RegisterConfigRoutes(adminRoutes, configService)
		RegisterClientRoutes(adminRoutes, c)
	}

	return r
}
