// router/router.go
package router

import (
	"context"
	"suasor/di/container"
	"suasor/router/middleware"
	"suasor/services"
	"suasor/utils/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(ctx context.Context, c *container.Container) *gin.Engine {
	r := gin.Default()
	log := logger.LoggerFromContext(ctx)

	log.Info().Msg("Setting up CORS middleware")
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
	log.Info().Msg("Setting up API v1 routes")
	v1 := r.Group("/api/v1")

	// {base}/health
	healthService := container.MustGet[services.HealthService](c)
	log.Info().Msg("Setting up health routes")
	RegisterHealthRoutes(v1, healthService)
	// {base}/auth
	log.Info().Msg("Setting up auth routes")
	authService := container.MustGet[services.AuthService](c)
	RegisterAuthRoutes(v1, authService)

	// Serve static avatar files
	r.Static("/uploads/avatars", appConfig.App.AvatarPath)

	// Protected Routes
	authenticated := v1.Group("")
	authenticated.Use(middleware.VerifyToken(authService))
	{
		// User Centric Data
		// {base}/user/
		RegisterUserRoutes(authenticated, c)

		// {base}/user-config/
		RegisterUserConfigRoutes(authenticated, c)

		// {base}/user-data/
		RegisterMediaItemDataRoutes(authenticated, c) // Register media play history routes

		// {base}/item/
		RegisterMediaItemRoutes(ctx, authenticated, c)

		// {base}/people/
		RegisterPeopleBasedRoutes(authenticated, c)

		// {base}/metadata/
		RegisterMetadataRoutes(authenticated, c) // Register metadata routes

		// {base}/playlist or {base}/collection
		RegisterLocalMediaListRoutes(authenticated, c)

		// {base}/history/
		RegisterMediaPlayHistoryRoutes(authenticated, c)

		// {base}/ai/
		RegisterAIRoutes(authenticated, c) // Register AI routes

		// {base}/jobs/
		RegisterJobRoutes(authenticated, c) // Register job routes

		// {base}/recommendations/
		RegisterRecommendationRoutes(authenticated, c) // Register recommendation routes

		// {base}/search/
		RegisterSearchRoutes(authenticated, c) // Register search routes
	}

	//Admin Routes
	adminRoutes := v1.Group("/admin")
	adminRoutes.Use(middleware.VerifyToken(authService), middleware.RequireRole("admin"))
	{
		// {base}/admin/config/
		RegisterConfigRoutes(adminRoutes, configService)
		// {base}/admin/client/
		RegisterClientRoutes(ctx, adminRoutes, c)
		// {base}/admin/clients/
		RegisterClientsRoutes(authenticated, c) // Register all clients route
	}

	return r
}
