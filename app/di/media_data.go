// app/di/media_data.go
package di

import (
	"context"
	"gorm.io/gorm"
	"suasor/app/container"
	"suasor/app/di/factories"
	"suasor/app/handlers"
	"suasor/app/repository"
	"suasor/app/services"
	"suasor/client"
	clienttypes "suasor/client/types"
	repo "suasor/repository"
)

// RegisterMediaData registers the media data factory and all media-related repositories
func RegisterMediaData(ctx context.Context, c *container.Container) {

	//  (Factory) MediaDataFactory
	container.RegisterFactory[factories.MediaDataFactory](c, func(c *container.Container) factories.MediaDataFactory {
		db := container.MustGet[*gorm.DB](c)
		clientFactory := container.MustGet[*client.ClientFactoryService](c)
		// We implement our factory in media_factory.go
		return createMediaDataFactory(db, clientFactory)
	})

	// --- REPOSITORY FACTORIES --- //

	//  MediaItem
	//  (Repositories) Core MediaItem Repositories
	container.RegisterFactory[repository.CoreMediaItemRepositories](c, func(c *container.Container) repository.CoreMediaItemRepositories {
		factory := container.MustGet[factories.MediaDataFactory](c)
		return factory.CreateCoreRepositories()
	})
	//  (Repositories) User MediaItem Repositories
	container.RegisterFactory[repository.UserMediaDataRepositories](c, func(c *container.Container) repository.UserMediaDataRepositories {
		factory := container.MustGet[factories.MediaDataFactory](c)
		return factory.CreateUserDataRepositories()
	})
	//  (Repositories) Client MediaItem Repositories
	container.RegisterFactory[repository.ClientMediaItemRepositories](c, func(c *container.Container) repository.ClientMediaItemRepositories {
		factory := container.MustGet[factories.MediaDataFactory](c)
		return factory.CreateClientMediaItemRepositories()
	})

	//  Core Services Registration
	//  (Services) Core MediaItem Services
	container.RegisterFactory[services.CoreMediaItemServices](c, func(c *container.Container) services.CoreMediaItemServices {
		factory := container.MustGet[factories.MediaDataFactory](c)
		repos := container.MustGet[repository.CoreMediaItemRepositories](c)
		return factory.CreateCoreServices(repos)
	})

	//  UserMediaItemData
	//  (Repositories) Core UserMediaItemData Repositories
	container.RegisterFactory[repository.CoreUserMediaItemDataRepositories](c, func(c *container.Container) repository.CoreUserMediaItemDataRepositories {
		factory := container.MustGet[factories.MediaDataFactory](c)
		return factory.CreateCoreDataRepositories()
	})
	//  (Repositories) User UserMediaItemData Repositories
	container.RegisterFactory[repository.UserMediaDataRepositories](c, func(c *container.Container) repository.UserMediaDataRepositories {
		factory := container.MustGet[factories.MediaDataFactory](c)
		return factory.CreateUserDataRepositories()
	})
	//  (Repositories) Client UserMediaItemData Repositories
	container.RegisterFactory[repository.ClientUserMediaDataRepositories](c, func(c *container.Container) repository.ClientUserMediaDataRepositories {
		factory := container.MustGet[factories.MediaDataFactory](c)
		return factory.CreateClientDataRepositories()
	})

	// --- SERVICE FACTORIES --- //

	//  MediaItem
	//  (Services) Core MediaItem Services
	container.RegisterFactory[services.CoreMediaItemServices](c, func(c *container.Container) services.CoreMediaItemServices {
		factory := container.MustGet[factories.MediaDataFactory](c)
		repos := container.MustGet[repository.CoreMediaItemRepositories](c)
		return factory.CreateCoreServices(repos)
	})
	//  (Services) User MediaItem Services
	container.RegisterFactory[services.UserMediaItemServices](c, func(c *container.Container) services.UserMediaItemServices {
		factory := container.MustGet[factories.MediaDataFactory](c)
		coreServices := container.MustGet[services.CoreMediaItemServices](c)
		userRepos := container.MustGet[repository.UserMediaItemRepositories](c)
		return factory.CreateUserServices(coreServices, userRepos)
	})
	//  (Services) Client MediaItem Services
	container.RegisterFactory[services.ClientMediaItemServices[clienttypes.ClientMediaConfig]](c, func(c *container.Container) services.ClientMediaItemServices[clienttypes.ClientMediaConfig] {
		factory := container.MustGet[factories.MediaDataFactory](c)
		coreServices := container.MustGet[services.CoreMediaItemServices](c)
		clientRepo := container.MustGet[repo.ClientRepository[clienttypes.ClientMediaConfig]](c)
		clientItemRepos := container.MustGet[repository.ClientMediaItemRepositories](c)
		return factory.CreateClientServices(coreServices, clientRepo, clientItemRepos)
	})
	//  UserMediaItemData
	//  (Services) Core UserMediaItemData Services
	container.RegisterFactory[services.CoreUserMediaItemDataServices](c, func(c *container.Container) services.CoreUserMediaItemDataServices {
		factory := container.MustGet[factories.MediaDataFactory](c)
		repos := container.MustGet[repository.CoreMediaItemRepositories](c)
		return factory.CreateCoreDataServices(repos)
	})
	//  (Services) User UserMediaItemData Services
	container.RegisterFactory[services.UserMediaItemDataServices](c, func(c *container.Container) services.UserMediaItemDataServices {
		factory := container.MustGet[factories.MediaDataFactory](c)
		coreDataServices := container.MustGet[services.CoreUserMediaItemDataServices](c)
		userRepos := container.MustGet[repository.UserMediaDataRepositories](c)
		return factory.CreateUserDataServices(coreDataServices, userRepos)
	})

	// (Services) Client UserMediaItemData Services
	container.RegisterFactory[services.ClientUserMediaItemDataServices](c, func(c *container.Container) services.ClientUserMediaItemDataServices {
		factory := container.MustGet[factories.MediaDataFactory](c)
		coreDataServices := container.MustGet[services.UserMediaItemDataServices](c)
		clientRepos := container.MustGet[repository.ClientUserMediaDataRepositories](c)
		return factory.CreateClientDataServices(coreDataServices, clientRepos)
	})

	// --- HANDLER --- //
	// (Handlers) Core MediaItem Handlers
	container.RegisterFactory[handlers.CoreMediaItemHandlers](c, func(c *container.Container) handlers.CoreMediaItemHandlers {
		factory := container.MustGet[factories.MediaDataFactory](c)
		coreServices := container.MustGet[services.CoreMediaItemServices](c)
		return factory.CreateCoreMediaItemHandlers(coreServices)
	})
	// (Handlers) User MediaItem Handlers
	container.RegisterFactory[handlers.UserMediaItemHandlers](c, func(c *container.Container) handlers.UserMediaItemHandlers {
		factory := container.MustGet[factories.MediaDataFactory](c)
		userServices := container.MustGet[services.UserMediaItemServices](c)
		coreHandlers := container.MustGet[handlers.CoreMediaItemHandlers](c)
		return factory.CreateUserMediaItemHandlers(userServices, coreHandlers)
	})
	// (Handlers) Client MediaItem Handlers
	container.RegisterFactory[handlers.ClientMediaItemHandlers[clienttypes.ClientMediaConfig]](c, func(c *container.Container) handlers.ClientMediaItemHandlers[clienttypes.ClientMediaConfig] {
		factory := container.MustGet[factories.MediaDataFactory](c)
		clientServices := container.MustGet[services.ClientMediaItemServices[clienttypes.ClientMediaConfig]](c)
		userMediaItemServices := container.MustGet[services.UserMediaItemServices](c)
		coreHandlers := container.MustGet[handlers.UserMediaItemHandlers](c)
		return factory.CreateClientMediaItemHandlers(clientServices, userMediaItemServices, coreHandlers)
	})
	// (Handlers) Core UserMediaItemData Handlers
	container.RegisterFactory[handlers.CoreMediaItemDataHandlers](c, func(c *container.Container) handlers.CoreMediaItemDataHandlers {
		factory := container.MustGet[factories.MediaDataFactory](c)
		coreServices := container.MustGet[services.CoreUserMediaItemDataServices](c)
		return factory.CreateCoreDataHandlers(coreServices)
	})
	// (Handlers) User UserMediaItemData Handlers
	container.RegisterFactory[handlers.UserMediaItemDataHandlers](c, func(c *container.Container) handlers.UserMediaItemDataHandlers {
		factory := container.MustGet[factories.MediaDataFactory](c)
		userServices := container.MustGet[services.UserMediaItemDataServices](c)
		coreHandlers := container.MustGet[handlers.CoreMediaItemDataHandlers](c)
		return factory.CreateUserDataHandlers(userServices, coreHandlers)
	})
	// (Handlers) Client UserMediaItemData Handlers
	container.RegisterFactory[handlers.ClientMediaItemDataHandlers](c, func(c *container.Container) handlers.ClientMediaItemDataHandlers {
		factory := container.MustGet[factories.MediaDataFactory](c)
		dataServices := container.MustGet[services.ClientUserMediaItemDataServices](c)
		userDataHandlers := container.MustGet[handlers.UserMediaItemDataHandlers](c)
		return factory.CreateClientDataHandlers(dataServices, userDataHandlers)
	})
}
