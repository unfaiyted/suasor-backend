// app/di/repositories.go
package di

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"suasor/app/container"
	apprepos "suasor/app/repository"
	"suasor/client/types"
	"suasor/repository"
	"suasor/services/jobs"
	"suasor/services/jobs/recommendation"
	"suasor/utils"
)

// RegisterRepositories registers all repository dependencies
func RegisterRepositories(ctx context.Context, c *container.Container) {
	log := utils.LoggerFromContext(ctx)
	// User repositories
	log.Info().Msg("Registering user repository")
	// Use RegisterSingleton to ensure it's only initialized once
	container.RegisterSingleton[repository.UserRepository](c, func(c *container.Container) repository.UserRepository {
		fmt.Println("Creating UserRepository")
		db := container.MustGet[*gorm.DB](c)
		repo := repository.NewUserRepository(db)
		fmt.Println("UserRepository created successfully")
		return repo
	})

	// User config repository
	log.Info().Msg("Registering user config repository")
	container.RegisterFactory[repository.UserConfigRepository](c, func(c *container.Container) repository.UserConfigRepository {
		db := container.MustGet[*gorm.DB](c)
		return repository.NewUserConfigRepository(db)
	})

	// Session repository
	log.Info().Msg("Registering session repository")
	// Use RegisterSingleton to ensure it's only initialized once
	container.RegisterSingleton[repository.SessionRepository](c, func(c *container.Container) repository.SessionRepository {
		fmt.Println("Creating SessionRepository")
		db := container.MustGet[*gorm.DB](c)
		repo := repository.NewSessionRepository(db)
		fmt.Println("SessionRepository created successfully")
		return repo
	})

	// Register client repositories
	log.Info().Msg("Registering client repositories")
	registerClientRepositories(c)

	// Job repository
	log.Info().Msg("Registering job repository")

	registerJobRepositories(ctx, c)

	// Person and credit repositories
	log.Info().Msg("Registering person and credit repositories")
	container.RegisterFactory[repository.PersonRepository](c, func(c *container.Container) repository.PersonRepository {
		db := container.MustGet[*gorm.DB](c)
		return repository.NewPersonRepository(db)
	})

	log.Info().Msg("Registering credit repository")
	container.RegisterFactory[repository.CreditRepository](c, func(c *container.Container) repository.CreditRepository {
		db := container.MustGet[*gorm.DB](c)
		return repository.NewCreditRepository(db)
	})
	// Search repository
	log.Info().Msg("Registering search repository")
	container.RegisterFactory[repository.SearchRepository](c, func(c *container.Container) repository.SearchRepository {
		db := container.MustGet[*gorm.DB](c)
		return repository.NewSearchRepository(db)
	})
}

// Register all client-type specific repositories
func registerClientRepositories(c *container.Container) {
	// Media client repositories
	container.RegisterFactory[repository.ClientRepository[*types.EmbyConfig]](c, func(c *container.Container) repository.ClientRepository[*types.EmbyConfig] {
		db := container.MustGet[*gorm.DB](c)
		return repository.NewClientRepository[*types.EmbyConfig](db)
	})

	container.RegisterFactory[repository.ClientRepository[*types.JellyfinConfig]](c, func(c *container.Container) repository.ClientRepository[*types.JellyfinConfig] {
		db := container.MustGet[*gorm.DB](c)
		return repository.NewClientRepository[*types.JellyfinConfig](db)
	})

	container.RegisterFactory[repository.ClientRepository[*types.PlexConfig]](c, func(c *container.Container) repository.ClientRepository[*types.PlexConfig] {
		db := container.MustGet[*gorm.DB](c)
		return repository.NewClientRepository[*types.PlexConfig](db)
	})

	container.RegisterFactory[repository.ClientRepository[*types.SubsonicConfig]](c, func(c *container.Container) repository.ClientRepository[*types.SubsonicConfig] {
		db := container.MustGet[*gorm.DB](c)
		return repository.NewClientRepository[*types.SubsonicConfig](db)
	})

	// Automation client repositories
	container.RegisterFactory[repository.ClientRepository[*types.SonarrConfig]](c, func(c *container.Container) repository.ClientRepository[*types.SonarrConfig] {
		db := container.MustGet[*gorm.DB](c)
		return repository.NewClientRepository[*types.SonarrConfig](db)
	})

	container.RegisterFactory[repository.ClientRepository[*types.RadarrConfig]](c, func(c *container.Container) repository.ClientRepository[*types.RadarrConfig] {
		db := container.MustGet[*gorm.DB](c)
		return repository.NewClientRepository[*types.RadarrConfig](db)
	})

	container.RegisterFactory[repository.ClientRepository[*types.LidarrConfig]](c, func(c *container.Container) repository.ClientRepository[*types.LidarrConfig] {
		db := container.MustGet[*gorm.DB](c)
		return repository.NewClientRepository[*types.LidarrConfig](db)
	})

	// AI client repositories
	container.RegisterFactory[repository.ClientRepository[*types.ClaudeConfig]](c, func(c *container.Container) repository.ClientRepository[*types.ClaudeConfig] {
		db := container.MustGet[*gorm.DB](c)
		return repository.NewClientRepository[*types.ClaudeConfig](db)
	})

	container.RegisterFactory[repository.ClientRepository[*types.OpenAIConfig]](c, func(c *container.Container) repository.ClientRepository[*types.OpenAIConfig] {
		db := container.MustGet[*gorm.DB](c)
		return repository.NewClientRepository[*types.OpenAIConfig](db)
	})

	container.RegisterFactory[repository.ClientRepository[*types.OllamaConfig]](c, func(c *container.Container) repository.ClientRepository[*types.OllamaConfig] {
		db := container.MustGet[*gorm.DB](c)
		return repository.NewClientRepository[*types.OllamaConfig](db)
	})

	// Client repository collection
	container.RegisterFactory[apprepos.ClientRepositories](c, func(c *container.Container) apprepos.ClientRepositories {
		embyRepo := container.MustGet[repository.ClientRepository[*types.EmbyConfig]](c)
		jellyfinRepo := container.MustGet[repository.ClientRepository[*types.JellyfinConfig]](c)
		plexRepo := container.MustGet[repository.ClientRepository[*types.PlexConfig]](c)
		subsonicRepo := container.MustGet[repository.ClientRepository[*types.SubsonicConfig]](c)
		sonarrRepo := container.MustGet[repository.ClientRepository[*types.SonarrConfig]](c)
		radarrRepo := container.MustGet[repository.ClientRepository[*types.RadarrConfig]](c)
		lidarrRepo := container.MustGet[repository.ClientRepository[*types.LidarrConfig]](c)
		claudeRepo := container.MustGet[repository.ClientRepository[*types.ClaudeConfig]](c)
		openaiRepo := container.MustGet[repository.ClientRepository[*types.OpenAIConfig]](c)
		ollamaRepo := container.MustGet[repository.ClientRepository[*types.OllamaConfig]](c)

		return apprepos.NewClientRepositories(
			embyRepo, jellyfinRepo,
			plexRepo, subsonicRepo,
			sonarrRepo, radarrRepo,
			lidarrRepo, claudeRepo,
			openaiRepo, ollamaRepo,
		)
	})
}

func registerJobRepositories(ctx context.Context, c *container.Container) {
	log := utils.LoggerFromContext(ctx)
	log.Info().Msg("Registering job repository")
	container.RegisterFactory[repository.JobRepository](c, func(c *container.Container) repository.JobRepository {
		db := container.MustGet[*gorm.DB](c)
		return repository.NewJobRepository(db)
	})

	// Recommendation Repo
	log.Info().Msg("Registering recommendation repository")
	container.RegisterFactory[repository.RecommendationRepository](c, func(c *container.Container) repository.RecommendationRepository {
		db := container.MustGet[*gorm.DB](c)
		return repository.NewRecommendationRepository(db)
	})

	// Watch History Sync Job
	log.Info().Msg("Registering watch history sync job repository")
	container.RegisterFactory[jobs.WatchHistorySyncJob](c, func(c *container.Container) jobs.WatchHistorySyncJob {
		return *jobs.NewWatchHistorySyncJob(ctx, c)
	})
	// Favorites Sync Job
	log.Info().Msg("Registering favorites sync job repository")
	container.RegisterFactory[jobs.FavoritesSyncJob](c, func(c *container.Container) jobs.FavoritesSyncJob {
		return *jobs.NewFavoritesSyncJob(ctx, c)
	})

	// Media Sync Job
	log.Info().Msg("Registering media sync job repository")
	container.RegisterFactory[jobs.MediaSyncJob](c, func(c *container.Container) jobs.MediaSyncJob {
		return *jobs.NewMediaSyncJob(ctx, c)
	})

	// Recommendation Job
	log.Info().Msg("Registering recommendation job repository")
	container.RegisterFactory[recommendation.RecommendationJob](c, func(c *container.Container) recommendation.RecommendationJob {
		return *recommendation.NewRecommendationJob(ctx, c)
	})

}
