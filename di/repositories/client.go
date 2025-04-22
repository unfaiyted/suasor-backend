// app/di/repositories/client.go
package repositories

import (
	"context"
	"gorm.io/gorm"
	clienttypes "suasor/clients/types"
	"suasor/di/container"
	"suasor/repository"
	repobundles "suasor/repository/bundles"
	"suasor/utils/logger"
)

// Register all client-type specific repositories
func registerClientRepositories(ctx context.Context, c *container.Container) {
	log := logger.LoggerFromContext(ctx)
	db := container.MustGet[*gorm.DB](c)

	// Media client repositories
	log.Info().Msg("Registering client repositories")
	registerClientRepository[*clienttypes.EmbyConfig](c, db)
	registerClientRepository[*clienttypes.JellyfinConfig](c, db)
	registerClientRepository[*clienttypes.PlexConfig](c, db)
	registerClientRepository[*clienttypes.SubsonicConfig](c, db)

	// Automation client repositories
	log.Info().Msg("Registering automation client repositories")
	registerClientRepository[*clienttypes.SonarrConfig](c, db)
	registerClientRepository[*clienttypes.RadarrConfig](c, db)
	registerClientRepository[*clienttypes.LidarrConfig](c, db)

	// AI client repositories
	log.Info().Msg("Registering AI client repositories")
	registerClientRepository[*clienttypes.ClaudeConfig](c, db)
	registerClientRepository[*clienttypes.OpenAIConfig](c, db)
	registerClientRepository[*clienttypes.OllamaConfig](c, db)

	// Metadata client service
	registerClientRepository[*clienttypes.TMDBConfig](c, db)

	// Client repository collection
	container.RegisterFactory[repobundles.ClientRepositories](c, func(c *container.Container) repobundles.ClientRepositories {
		embyRepo := container.MustGet[repository.ClientRepository[*clienttypes.EmbyConfig]](c)
		jellyfinRepo := container.MustGet[repository.ClientRepository[*clienttypes.JellyfinConfig]](c)
		plexRepo := container.MustGet[repository.ClientRepository[*clienttypes.PlexConfig]](c)
		subsonicRepo := container.MustGet[repository.ClientRepository[*clienttypes.SubsonicConfig]](c)
		sonarrRepo := container.MustGet[repository.ClientRepository[*clienttypes.SonarrConfig]](c)
		radarrRepo := container.MustGet[repository.ClientRepository[*clienttypes.RadarrConfig]](c)
		lidarrRepo := container.MustGet[repository.ClientRepository[*clienttypes.LidarrConfig]](c)
		claudeRepo := container.MustGet[repository.ClientRepository[*clienttypes.ClaudeConfig]](c)
		openaiRepo := container.MustGet[repository.ClientRepository[*clienttypes.OpenAIConfig]](c)
		ollamaRepo := container.MustGet[repository.ClientRepository[*clienttypes.OllamaConfig]](c)

		return repobundles.NewClientRepositories(
			embyRepo, jellyfinRepo,
			plexRepo, subsonicRepo,
			sonarrRepo, radarrRepo,
			lidarrRepo, claudeRepo,
			openaiRepo, ollamaRepo,
		)
	})
}

func registerClientRepository[T clienttypes.ClientConfig](c *container.Container, db *gorm.DB) {
	// Register core user media item data repositories
	container.RegisterFactory[repository.ClientRepository[T]](c, func(c *container.Container) repository.ClientRepository[T] {
		return repository.NewClientRepository[T](db)
	})
}
