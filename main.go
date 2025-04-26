package main

import (
	"context"
	"suasor/di"
	"suasor/di/container"
	"suasor/repository"
	"suasor/router"
	"suasor/router/middleware"
	"suasor/services"
	"suasor/services/jobs"
	"suasor/types"
	"suasor/utils/db"
	logger "suasor/utils/logger"

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

	log.Info().Msg("Initializing database")
	db, err := database.Initialize(ctx, dbConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database:")
	}

	log.Info().Msg("Initializing application dependencies")
	deps := di.InitializeDependencies(ctx, db, configService)

	log.Info().Msg("Initializing job service")
	jobService := container.MustGet[jobs.JobService](deps.GetContainer())
	log.Info().Msg("Job service initialized")

	// Start the job scheduler
	log.Info().Msg("Starting job scheduler")
	if err := jobService.StartScheduler(); err != nil {
		log.Error().Err(err).Msg("Failed to start job scheduler")
	}

	// Sync job schedules from database
	if err := jobService.SyncDatabaseJobs(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to sync job schedules from database")
	}
	log.Info().Msg("Job scheduler disabled temporarily")

	log.Info().Msg("Setting up router")
	r := router.Setup(ctx, deps.GetContainer())

	log.Info().Msg("Setting up middleware")
	r.Use(middleware.LoggerMiddleware())

	// Swagger Docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	port := ":" + appConfig.HTTP.Port
	log.Info().Str("port", port).Msg("Starting server")
	r.Run(port)
}
