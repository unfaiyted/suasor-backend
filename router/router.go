// router/router.go
package router

import (
	"context"
	"suasor/app"
	"suasor/router/middleware"
	"suasor/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(ctx context.Context, deps *app.AppDependencies) *gin.Engine {
	r := gin.Default()
	log := utils.LoggerFromContext(ctx)

	appConfig := deps.SystemServices.ConfigService().GetConfig()
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

	RegisterHealthRoutes(v1, deps.SystemServices.HealthService())
	RegisterAuthRoutes(v1, deps.AuthService())
	
	// Serve static avatar files
	avatarPath := deps.SystemServices.ConfigService().GetConfig().App.AvatarPath
	r.Static("/uploads/avatars", avatarPath)

	// Protected Routes
	authenticated := v1.Group("")
	authenticated.Use(middleware.VerifyToken(deps.AuthService()))
	{
		// Register all routes
		RegisterUserRoutes(authenticated, deps)
		RegisterUserConfigRoutes(authenticated, deps)
		RegisterMediaItemRoutes(authenticated, deps)
		RegisterMediaClientRoutes(authenticated, deps)
		RegisterMetadataRoutes(authenticated)      // Register metadata routes
		RegisterAIRoutes(authenticated, deps)      // Register AI routes
		RegisterClientsRoutes(authenticated, deps) // Register all clients route
		RegisterJobRoutes(authenticated, deps.JobServices.JobService()) // Register job routes
		RegisterRecommendationRoutes(authenticated, deps) // Register recommendation routes
		RegisterSearchRoutes(authenticated, deps.SearchHandler()) // Register search routes
	}

	//Admin Routes
	adminRoutes := v1.Group("/admin")
	adminRoutes.Use(middleware.VerifyToken(deps.AuthService()), middleware.RequireRole("admin"))
	{
		RegisterConfigRoutes(adminRoutes, deps.SystemServices.ConfigService())
		RegisterClientRoutes(adminRoutes, deps)
	}

	return r
}