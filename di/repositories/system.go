package repositories

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"suasor/di/container"
	"suasor/repository"
	"suasor/utils/logger"
)

// System repositories
func registerSystemRepositories(ctx context.Context, c *container.Container) {
	log := logger.LoggerFromContext(ctx)

	log.Info().Msg("Registering user repository")
	container.RegisterSingleton[repository.UserRepository](c, func(c *container.Container) repository.UserRepository {
		fmt.Println("Creating UserRepository")
		db := container.MustGet[*gorm.DB](c)
		repo := repository.NewUserRepository(db)
		fmt.Println("UserRepository created successfully")
		return repo
	})

	log.Info().Msg("Registering user config repository")
	container.RegisterFactory[repository.UserConfigRepository](c, func(c *container.Container) repository.UserConfigRepository {
		db := container.MustGet[*gorm.DB](c)
		return repository.NewUserConfigRepository(db)
	})

	log.Info().Msg("Registering session repository")
	container.RegisterSingleton[repository.SessionRepository](c, func(c *container.Container) repository.SessionRepository {
		fmt.Println("Creating SessionRepository")
		db := container.MustGet[*gorm.DB](c)
		repo := repository.NewSessionRepository(db)
		fmt.Println("SessionRepository created successfully")
		return repo
	})

}
