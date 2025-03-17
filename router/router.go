// router/router.go
package router

import (
	"context"
	"suasor/repository"
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

	// Register all routes
	RegisterConfigRoutes(v1, configService)
	RegisterHealthRoutes(v1, healthService)
	RegisterAuthRoutes(v1, authService)
	RegisterUserRoutes(v1, userService)
	RegisterUserConfigRoutes(v1, userConfigService)

	return r
}
