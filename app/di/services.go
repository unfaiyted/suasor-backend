// app/di/services.go
package di

import (
	"gorm.io/gorm"
	"suasor/app/container"
	"suasor/client"
	mediatypes "suasor/client/media/types"
	"suasor/client/types"
	"suasor/repository"
	"suasor/services"
)

// RegisterServices registers all service dependencies
func RegisterServices(c *container.Container) {
	// Register system services
	registerSystemServices(c)

	// Register client services
	registerClientServices(c)

	// Register three-pronged architecture services
	registerThreeProngedServices(c)
}

// Register system-level services
func registerSystemServices(c *container.Container) {
	// Health service
	container.RegisterFactory[services.HealthService](c, func(c *container.Container) services.HealthService {
		db := container.MustGet[*gorm.DB](c)
		return services.NewHealthService(db)
	})

	// Media services
	container.RegisterFactory[interfaces.MediaServices](c, func(c *container.Container) interfaces.MediaServices {
		personRepo := container.MustGet[repository.PersonRepository](c)
		creditRepo := container.MustGet[repository.CreditRepository](c)
		return &mediaServicesImpl{
			personRepo: personRepo,
			creditRepo: creditRepo,
		}
	})
}

// Register client-specific services
func registerClientServices(c *container.Container) {
	// Media clients
	container.RegisterFactory[services.ClientService[*types.EmbyConfig]](c, func(c *container.Container) services.ClientService[*types.EmbyConfig] {
		clientFactory := container.MustGet[client.ClientFactoryService](c)
		repo := container.MustGet[repository.ClientRepository[*types.EmbyConfig]](c)
		return services.NewClientService[*types.EmbyConfig](&clientFactory, repo)
	})

	container.RegisterFactory[services.ClientService[*types.JellyfinConfig]](c, func(c *container.Container) services.ClientService[*types.JellyfinConfig] {
		clientFactory := container.MustGet[client.ClientFactoryService](c)
		repo := container.MustGet[repository.ClientRepository[*types.JellyfinConfig]](c)
		return services.NewClientService[*types.JellyfinConfig](&clientFactory, repo)
	})

	container.RegisterFactory[services.ClientService[*types.PlexConfig]](c, func(c *container.Container) services.ClientService[*types.PlexConfig] {
		clientFactory := container.MustGet[client.ClientFactoryService](c)
		repo := container.MustGet[repository.ClientRepository[*types.PlexConfig]](c)
		return services.NewClientService[*types.PlexConfig](&clientFactory, repo)
	})

	container.RegisterFactory[services.ClientService[*types.SubsonicConfig]](c, func(c *container.Container) services.ClientService[*types.SubsonicConfig] {
		clientFactory := container.MustGet[client.ClientFactoryService](c)
		repo := container.MustGet[repository.ClientRepository[*types.SubsonicConfig]](c)
		return services.NewClientService[*types.SubsonicConfig](&clientFactory, repo)
	})

	// Automation clients
	container.RegisterFactory[services.ClientService[*types.SonarrConfig]](c, func(c *container.Container) services.ClientService[*types.SonarrConfig] {
		clientFactory := container.MustGet[client.ClientFactoryService](c)
		repo := container.MustGet[repository.ClientRepository[*types.SonarrConfig]](c)
		return services.NewClientService[*types.SonarrConfig](&clientFactory, repo)
	})

	container.RegisterFactory[services.ClientService[*types.RadarrConfig]](c, func(c *container.Container) services.ClientService[*types.RadarrConfig] {
		clientFactory := container.MustGet[client.ClientFactoryService](c)
		repo := container.MustGet[repository.ClientRepository[*types.RadarrConfig]](c)
		return services.NewClientService[*types.RadarrConfig](&clientFactory, repo)
	})

	container.RegisterFactory[services.ClientService[*types.LidarrConfig]](c, func(c *container.Container) services.ClientService[*types.LidarrConfig] {
		clientFactory := container.MustGet[client.ClientFactoryService](c)
		repo := container.MustGet[repository.ClientRepository[*types.LidarrConfig]](c)
		return services.NewClientService[*types.LidarrConfig](&clientFactory, repo)
	})

	// AI clients
	container.RegisterFactory[services.ClientService[*types.ClaudeConfig]](c, func(c *container.Container) services.ClientService[*types.ClaudeConfig] {
		clientFactory := container.MustGet[client.ClientFactoryService](c)
		repo := container.MustGet[repository.ClientRepository[*types.ClaudeConfig]](c)
		return services.NewClientService[*types.ClaudeConfig](&clientFactory, repo)
	})

	container.RegisterFactory[services.ClientService[*types.OpenAIConfig]](c, func(c *container.Container) services.ClientService[*types.OpenAIConfig] {
		clientFactory := container.MustGet[client.ClientFactoryService](c)
		repo := container.MustGet[repository.ClientRepository[*types.OpenAIConfig]](c)
		return services.NewClientService[*types.OpenAIConfig](&clientFactory, repo)
	})

	container.RegisterFactory[services.ClientService[*types.OllamaConfig]](c, func(c *container.Container) services.ClientService[*types.OllamaConfig] {
		clientFactory := container.MustGet[client.ClientFactoryService](c)
		repo := container.MustGet[repository.ClientRepository[*types.OllamaConfig]](c)
		return services.NewClientService[*types.OllamaConfig](&clientFactory, repo)
	})
}

// Register services for the three-pronged architecture
func registerThreeProngedServices(c *container.Container) {
	// Core media item services
	container.RegisterFactory[interfaces.CoreMediaItemServices](c, func(c *container.Container) interfaces.CoreMediaItemServices {
		factory := container.MustGet[interfaces.MediaDataFactory](c)
		repos := container.MustGet[interfaces.CoreMediaItemRepositories](c)
		return factory.CreateCoreServices(repos)
	})

	// User media item services
	container.RegisterFactory[interfaces.UserMediaItemServices](c, func(c *container.Container) interfaces.UserMediaItemServices {
		factory := container.MustGet[interfaces.MediaDataFactory](c)
		coreServices := container.MustGet[interfaces.CoreMediaItemServices](c)
		userRepos := container.MustGet[interfaces.UserRepositoryFactories](c)
		return factory.CreateUserServices(coreServices, userRepos)
	})

	// Client media item services
	container.RegisterFactory[interfaces.ClientMediaItemServices](c, func(c *container.Container) interfaces.ClientMediaItemServices {
		factory := container.MustGet[interfaces.MediaDataFactory](c)
		coreServices := container.MustGet[interfaces.CoreMediaItemServices](c)
		clientRepos := container.MustGet[interfaces.ClientRepositoryFactories](c)
		return factory.CreateClientServices(coreServices, clientRepos)
	})

	// Collection services
	container.RegisterFactory[services.CoreCollectionService](c, func(c *container.Container) services.CoreCollectionService {
		repos := container.MustGet[interfaces.CoreMediaItemRepositories](c)
		return services.NewCoreCollectionService(repos.CollectionRepo())
	})

	container.RegisterFactory[services.UserCollectionService](c, func(c *container.Container) services.UserCollectionService {
		coreService := container.MustGet[services.CoreCollectionService](c)
		userRepos := container.MustGet[interfaces.UserRepositoryFactories](c)
		mediaDataRepo := container.MustGet[repository.MediaItemRepository[mediatypes.MediaData]](c)
		return services.NewUserCollectionService(coreService, userRepos.CollectionUserRepo(), mediaDataRepo)
	})

	container.RegisterFactory[services.ClientMediaCollectionService](c, func(c *container.Container) services.ClientMediaCollectionService {
		clientRepo := container.MustGet[repository.ClientRepositoryCollection](c)
		return services.NewClientMediaCollectionService(clientRepo.AllRepos())
	})
}

