package repositories

import (
	"context"
	"gorm.io/gorm"
	"suasor/di/container"
	"suasor/repository"
	"suasor/utils/logger"
)

// RegisterMediaRepositories registers all media-related repositories for the three-pronged architecture
func RegisterMediaRepositories(ctx context.Context, c *container.Container) {
	log := logger.LoggerFromContext(ctx)
	// Core Media Item Repositories
	registerMediaItemRepositories(c)
	// Media Item Data Repositories
	registerMediaItemDataRepositories(c)

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
