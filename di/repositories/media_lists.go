package repositories

import (
	"context"
	"gorm.io/gorm"
	"suasor/clients/media/types"
	"suasor/di/container"
	"suasor/repository"
	"suasor/utils/logger"
)

// Register core repositories for all media types
func registerMediaListRepositories(ctx context.Context, c *container.Container) {
	db := container.MustGet[*gorm.DB](c)
	log := logger.LoggerFromContext(ctx)

	log.Info().Msg("Registering list repositories")

	container.RegisterFactory[repository.CoreListRepository[*types.Collection]](c, func(c *container.Container) repository.CoreListRepository[*types.Collection] {
		mediaItemRepo := container.MustGet[repository.CoreMediaItemRepository[*types.Collection]](c)
		return repository.NewCoreListRepository(db, mediaItemRepo)
	})

	container.RegisterFactory[repository.CoreListRepository[*types.Playlist]](c, func(c *container.Container) repository.CoreListRepository[*types.Playlist] {
		mediaItemRepo := container.MustGet[repository.CoreMediaItemRepository[*types.Playlist]](c)
		return repository.NewCoreListRepository(db, mediaItemRepo)
	})

}
