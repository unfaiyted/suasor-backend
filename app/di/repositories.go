// app/di/repositories.go
package di

import (
	"gorm.io/gorm"
	"suasor/app/container"
	apprepos "suasor/app/repository"
	"suasor/client/types"
	"suasor/repository"
)

// RegisterRepositories registers all repository dependencies
func RegisterRepositories(c *container.Container) {
	// User repositories
	container.RegisterFactory[repository.UserRepository](c, func(c *container.Container) repository.UserRepository {
		db := container.MustGet[*gorm.DB](c)
		return repository.NewUserRepository(db)
	})

	// User config repository
	container.RegisterFactory[repository.UserConfigRepository](c, func(c *container.Container) repository.UserConfigRepository {
		db := container.MustGet[*gorm.DB](c)
		return repository.NewUserConfigRepository(db)
	})

	// Session repository
	container.RegisterFactory[repository.SessionRepository](c, func(c *container.Container) repository.SessionRepository {
		db := container.MustGet[*gorm.DB](c)
		return repository.NewSessionRepository(db)
	})

	// Register client repositories
	registerClientRepositories(c)

	// Job repository
	container.RegisterFactory[repository.JobRepository](c, func(c *container.Container) repository.JobRepository {
		db := container.MustGet[*gorm.DB](c)
		return repository.NewJobRepository(db)
	})

	// Person and credit repositories
	container.RegisterFactory[repository.PersonRepository](c, func(c *container.Container) repository.PersonRepository {
		db := container.MustGet[*gorm.DB](c)
		return repository.NewPersonRepository(db)
	})

	container.RegisterFactory[repository.CreditRepository](c, func(c *container.Container) repository.CreditRepository {
		db := container.MustGet[*gorm.DB](c)
		return repository.NewCreditRepository(db)
	})

	// Search repository
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

