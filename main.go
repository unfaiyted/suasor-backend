package main

import (
	"context"
	"suasor/app"
	"suasor/database"
	"suasor/repository"
	"suasor/router"
	"suasor/router/middleware"
	"suasor/services"
	"suasor/types"
	logger "suasor/utils"

	_ "suasor/client/ai/claude"  // Force init() to run
	_ "suasor/client/media/emby" // Force init() to run
	_ "suasor/client/media/jellyfin"
	_ "suasor/client/media/plex"
	_ "suasor/client/media/subsonic"
	_ "suasor/client/metadata/tmdb"
	_ "suasor/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/rs/zerolog/log"
)

//	@title			Suasor API
//	@version		1.0
//	@description	API Server for Suasor
//	@termsOfService	http://swagger.io/terms/
//	@contact.name	Dane Miller
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:8080
// @BasePath	/api/v1
// @schemes	http
// @openapi	3.0.0
func main() {
	logger.Initialize()

	ctx := context.Background()

	configRepo := repository.NewConfigRepository()
	configService := services.NewConfigService(configRepo)

	if err := configService.InitConfig(ctx); err != nil {
		log.Fatal().Err(err).Msg("Failed to init config")
	}

	appConfig := configService.GetConfig()

	dbConfig := types.DatabaseConfig{
		Host:     appConfig.Db.Host,
		User:     appConfig.Db.User,
		Password: appConfig.Db.Password,
		Name:     appConfig.Db.Name,
		Port:     appConfig.Db.Port,
	}

	db, err := database.Initialize(ctx, dbConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database:")
	}

	deps := app.InitializeDependencies(db, configService)

	// Start the job scheduler
	log.Info().Msg("Starting job scheduler")
	if err := deps.JobServices.JobService().StartScheduler(); err != nil {
		log.Error().Err(err).Msg("Failed to start job scheduler")
	}

	// Sync job schedules from database
	if err := deps.JobServices.JobService().SyncDatabaseJobs(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to sync job schedules from database")
	}

	r := router.Setup(ctx, deps)

	r.Use(middleware.LoggerMiddleware())

	// Swagger Docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	port := ":" + appConfig.HTTP.Port
	log.Info().Str("port", port).Msg("Starting server")
	r.Run(port)
}
