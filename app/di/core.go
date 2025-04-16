// app/di/core.go
package di

import (
	"gorm.io/gorm"
	"suasor/app/container"
	"suasor/client"
	"suasor/repository"
	"suasor/services"
)

// Register core dependencies that are used throughout the application
func RegisterCore(c *container.Container, db *gorm.DB, configService services.ConfigService) {
	// Register core components
	c.Register(db)
	c.Register(configService)

	// Register config service as a factory for reuse
	container.RegisterFactory[services.ConfigService](c,
		func(c *container.Container) services.ConfigService {
			return configService
		})

	// Register client factory service
	// This is responsible for creating clients based on the client type and client ID
	// This is a singleton service that is shared across the application
	// It ensures that our external clients are created only once and reused throughout the application
	container.RegisterFactory[*client.ClientFactoryService](c,
		func(c *container.Container) *client.ClientFactoryService {
			return client.GetClientFactoryService()
		})

	// Register config repository
	container.RegisterFactory[repository.ConfigRepository](c, func(c *container.Container) repository.ConfigRepository {
		return container.MustGet[services.ConfigService](c).GetRepo()
	})
}
