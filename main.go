package main

import (
	"context"
	"suasor/app"
	"suasor/client"
	"suasor/database"
	"suasor/repository"
	"suasor/router"
	"suasor/router/middleware"
	"suasor/services"
	"suasor/types"
	logger "suasor/utils"

	_ "suasor/client/media/emby" // Force init() to run
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

	factory := client.GetClientFactoryService()
	deps := app.InitializeDependencies(db, configService, factory)
	r := router.Setup(ctx, deps, configService)

	r.Use(middleware.LoggerMiddleware())

	// Swagger Docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Start server
	r.Run(":8080")
}
