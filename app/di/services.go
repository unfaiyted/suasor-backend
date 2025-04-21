// app/di/services.go
package di

import (
	"context"
	"gorm.io/gorm"
	"suasor/app/container"
	"suasor/app/di/factories"
	apprepository "suasor/app/repository"
	appservices "suasor/app/services"
	"suasor/client"
	mediatypes "suasor/client/media/types"
	"suasor/client/types"
	clienttypes "suasor/client/types"
	"suasor/repository"
	"suasor/services"
)

// RegisterServices registers all service dependencies
func RegisterServices(ctx context.Context, c *container.Container) {
	// Register system services
	registerSystemServices(ctx, c)

	// Register client services
	registerClientServices(ctx, c)

	// Register three-pronged architecture services
	registerThreeProngedServices(ctx, c)
}

// Register system-level services
func registerSystemServices(ctx context.Context, c *container.Container) {
	// Health service
	container.RegisterFactory[services.HealthService](c, func(c *container.Container) services.HealthService {
		db := container.MustGet[*gorm.DB](c)
		return services.NewHealthService(db)
	})

	// Media services
	container.RegisterFactory[appservices.PeopleServices](c, func(c *container.Container) appservices.PeopleServices {
		personService := container.MustGet[services.PersonService](c)
		creditService := container.MustGet[services.CreditService](c)
		return appservices.NewPeopleServices(&personService, &creditService)
	})
}

// Register client-specific services
func registerClientServices(ctx context.Context, c *container.Container) {
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
func registerThreeProngedServices(ctx context.Context, c *container.Container) {
	// Core media item services
	container.RegisterFactory[appservices.CoreMediaItemServices](c, func(c *container.Container) appservices.CoreMediaItemServices {
		factory := container.MustGet[factories.MediaDataFactory](c)
		repos := container.MustGet[apprepository.CoreMediaItemRepositories](c)
		return factory.CreateCoreServices(repos)
	})

	// User media item services
	container.RegisterFactory[appservices.UserMediaItemServices](c, func(c *container.Container) appservices.UserMediaItemServices {
		factory := container.MustGet[factories.MediaDataFactory](c)
		coreServices := container.MustGet[appservices.CoreMediaItemServices](c)
		userRepos := container.MustGet[apprepository.UserMediaItemRepositories](c)
		return factory.CreateUserServices(coreServices, userRepos)
	})

	// Client media item services
	container.RegisterFactory[appservices.ClientMediaItemServices[clienttypes.ClientMediaConfig]](c, func(c *container.Container) appservices.ClientMediaItemServices[clienttypes.ClientMediaConfig] {
		factory := container.MustGet[factories.MediaDataFactory](c)
		coreServices := container.MustGet[appservices.CoreMediaItemServices](c)
		clientRepo := container.MustGet[repository.ClientRepository[clienttypes.ClientMediaConfig]](c)
		clientRepos := container.MustGet[apprepository.ClientMediaItemRepositories](c)
		return factory.CreateClientServices(coreServices, clientRepo, clientRepos)
	})

	// Collection services
	container.RegisterFactory[services.CoreListService[*mediatypes.Collection]](c, func(c *container.Container) services.CoreListService[*mediatypes.Collection] {
		repos := container.MustGet[apprepository.CoreMediaItemRepositories](c)
		return services.NewCoreListService(repos.CollectionRepo())
	})

	container.RegisterFactory[services.UserListService[*mediatypes.Collection]](c, func(c *container.Container) services.UserListService[*mediatypes.Collection] {
		coreService := container.MustGet[services.CoreListService[*mediatypes.Collection]](c)
		userItemRepos := container.MustGet[apprepository.UserMediaItemRepositories](c)
		userDataRepo := container.MustGet[repository.UserMediaItemDataRepository[*mediatypes.Collection]](c)

		return services.NewUserListService(coreService, userItemRepos.CollectionUserRepo(), userDataRepo)
	})

	container.RegisterFactory[services.ClientListService[*types.EmbyConfig, *mediatypes.Collection]](c, func(c *container.Container) services.ClientListService[*types.EmbyConfig, *mediatypes.Collection] {
		coreListService := container.MustGet[services.CoreListService[*mediatypes.Collection]](c)
		clientRepo := container.MustGet[repository.ClientRepository[*types.EmbyConfig]](c)
		clientFactory := container.MustGet[client.ClientFactoryService](c)
		return services.NewClientListService[*types.EmbyConfig, *mediatypes.Collection](coreListService, clientRepo, &clientFactory)
	})

	container.RegisterFactory[services.ClientListService[*types.JellyfinConfig, *mediatypes.Collection]](c, func(c *container.Container) services.ClientListService[*types.JellyfinConfig, *mediatypes.Collection] {
		coreListService := container.MustGet[services.CoreListService[*mediatypes.Collection]](c)
		clientRepo := container.MustGet[repository.ClientRepository[*types.JellyfinConfig]](c)
		clientFactory := container.MustGet[client.ClientFactoryService](c)
		return services.NewClientListService[*types.JellyfinConfig, *mediatypes.Collection](coreListService, clientRepo, &clientFactory)
	})

	container.RegisterFactory[services.ClientListService[*types.PlexConfig, *mediatypes.Collection]](c, func(c *container.Container) services.ClientListService[*types.PlexConfig, *mediatypes.Collection] {
		coreListService := container.MustGet[services.CoreListService[*mediatypes.Collection]](c)
		clientRepo := container.MustGet[repository.ClientRepository[*types.PlexConfig]](c)
		clientFactory := container.MustGet[client.ClientFactoryService](c)
		return services.NewClientListService[*types.PlexConfig, *mediatypes.Collection](coreListService, clientRepo, &clientFactory)
	})

	container.RegisterFactory[services.ClientListService[*types.SubsonicConfig, *mediatypes.Collection]](c, func(c *container.Container) services.ClientListService[*types.SubsonicConfig, *mediatypes.Collection] {
		coreListService := container.MustGet[services.CoreListService[*mediatypes.Collection]](c)
		clientRepo := container.MustGet[repository.ClientRepository[*types.SubsonicConfig]](c)
		clientFactory := container.MustGet[client.ClientFactoryService](c)
		return services.NewClientListService[*types.SubsonicConfig, *mediatypes.Collection](coreListService, clientRepo, &clientFactory)
	})

}
