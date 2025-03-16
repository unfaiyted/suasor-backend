// router/router.go
package router

import (
	"context"
	"suasor/repository"
	"suasor/services"
	"suasor/utils"

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

	shortenRepo := repository.NewShortenRepository(db)
	shortenService := services.NewShortenService(shortenRepo, configService.GetConfig().App.AppURL)

	// Register all routes
	RegisterConfigRoutes(v1, configService)
	RegisterHealthRoutes(v1, healthService)
	RegisterShortenRoutes(v1, shortenService)

	return r
}
