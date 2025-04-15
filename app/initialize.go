// app/initialize.go
package app

import (
	"context"
	"suasor/client"
	"suasor/services"
	
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// Initialize creates and initializes all application dependencies
// This is a cleaner, more modular approach compared to InitializeDependencies
func Initialize(ctx context.Context, db *gorm.DB, configService services.ConfigService) *AppDependencies {
	log.Info().Msg("Initializing application dependencies")
	
	// Create empty dependencies structure
	deps := &AppDependencies{
		db: db,
	}
	
	// Get the client factory service
	clientFactory := client.GetClientFactoryService()
	deps.ClientFactoryService = clientFactory
	
	// Create the service registrar
	registrar := NewServiceRegistrar(db, clientFactory, deps)
	
	// Register all services
	registrar.RegisterAllServices(configService)
	
	// Initialize job service related components
	initializeJobServices(ctx, deps)
	
	log.Info().Msg("Application dependencies initialized successfully")
	return deps
}

// initializeJobServices initializes and registers all job services
func initializeJobServices(ctx context.Context, deps *AppDependencies) {
	// Initialize job services
	jobService := deps.JobServices.JobService()
	
	// Register jobs with the job service if available
	if jobService != nil {
		// Register recommendation job if available
		if recommendationJob := deps.JobServices.RecommendationJob(); recommendationJob != nil {
			jobService.RegisterJob(recommendationJob)
		}
		
		// Register media sync job if available
		if mediaSyncJob := deps.JobServices.MediaSyncJob(); mediaSyncJob != nil {
			jobService.RegisterJob(mediaSyncJob)
		}
		
		// Register watch history sync job if available
		if watchHistorySyncJob := deps.JobServices.WatchHistorySyncJob(); watchHistorySyncJob != nil {
			jobService.RegisterJob(watchHistorySyncJob)
		}
		
		// Register favorites sync job if available
		if favoritesSyncJob := deps.JobServices.FavoritesSyncJob(); favoritesSyncJob != nil {
			jobService.RegisterJob(favoritesSyncJob)
		}
	}
}