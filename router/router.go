// router/router.go
package router

import (
	"context"
	"suasor/repository"
	"suasor/router/middleware"
	"suasor/services"
	"suasor/utils"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(ctx context.Context, db *gorm.DB, configService services.ConfigService) *gin.Engine {
	r := gin.Default()
	log := utils.LoggerFromContext(ctx)

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
	healthService := services.NewHealthService(db)

	userRepo := repository.NewUserRepository(db)
	userConfigRepo := repository.NewUserConfigRepository(db)
	sessionRepo := repository.NewSessionRepository(db)

	userService := services.NewUserService(userRepo)
	userConfigService := services.NewUserConfigService(userConfigRepo)

	authService := services.NewAuthService(userRepo,
		sessionRepo,
		appConfig.Auth.JWTSecret,
		time.Duration(appConfig.Auth.AccessExpiryMinutes)*time.Minute,
		time.Duration(appConfig.Auth.RefreshExpiryDays)*24*time.Hour,
		appConfig.Auth.TokenIssuer,
		appConfig.Auth.TokenAudience,
	)

	RegisterHealthRoutes(v1, healthService)
	RegisterAuthRoutes(v1, authService)

	// Protected Routes
	authenticated := v1.Group("")
	authenticated.Use(middleware.VerifyToken(authService))
	{
		// Register all routes
		RegisterUserRoutes(authenticated, userService)
		RegisterUserConfigRoutes(authenticated, userConfigService)
		RegisterClientRoutes(authenticated, db)
	}

	//Admin Routes
	adminRoutes := v1.Group("")
	adminRoutes.Use(middleware.VerifyToken(authService), middleware.RequireRole("admin"))
	{
		RegisterConfigRoutes(adminRoutes, configService)
	}

	return r
}
