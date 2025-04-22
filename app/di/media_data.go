// app/di/media_data.go
package di

import (
	"context"
	"gorm.io/gorm"
	"suasor/app/container"
	"suasor/app/di/factories"
	"suasor/app/services"
	"suasor/client"
	mediatypes "suasor/client/media/types"
	clienttypes "suasor/client/types"
	"suasor/handlers"
	"suasor/repository"
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
	container.RegisterFactory[repository.UserMediaItemRepositories](c, func(c *container.Container) repository.UserMediaItemRepositories {
		factory := container.MustGet[factories.MediaDataFactory](c)
		return factory.CreateUserRepositories()
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
	container.RegisterFactory[apphandlers.CoreMediaItemHandlers](c, func(c *container.Container) apphandlers.CoreMediaItemHandlers {
		factory := container.MustGet[factories.MediaDataFactory](c)
		coreServices := container.MustGet[services.CoreMediaItemServices](c)
		return factory.CreateCoreMediaItemHandlers(coreServices)
	})
	// (Handlers) User MediaItem Handlers
	container.RegisterFactory[apphandlers.UserMediaItemHandlers](c, func(c *container.Container) apphandlers.UserMediaItemHandlers {
		factory := container.MustGet[factories.MediaDataFactory](c)
		userServices := container.MustGet[services.UserMediaItemServices](c)
		coreHandlers := container.MustGet[apphandlers.CoreMediaItemHandlers](c)
		return factory.CreateUserMediaItemHandlers(userServices, coreHandlers)
	})
	// (Handlers) Client MediaItem Handlers
	container.RegisterFactory[apphandlers.ClientMediaItemHandlers[clienttypes.ClientMediaConfig]](c, func(c *container.Container) apphandlers.ClientMediaItemHandlers[clienttypes.ClientMediaConfig] {
		factory := container.MustGet[factories.MediaDataFactory](c)
		clientServices := container.MustGet[services.ClientMediaItemServices[clienttypes.ClientMediaConfig]](c)
		userMediaItemServices := container.MustGet[services.UserMediaItemServices](c)
		coreHandlers := container.MustGet[apphandlers.UserMediaItemHandlers](c)
		return factory.CreateClientMediaItemHandlers(clientServices, userMediaItemServices, coreHandlers)
	})
	// (Handlers) Core UserMediaItemData Handlers
	container.RegisterFactory[apphandlers.CoreMediaItemDataHandlers](c, func(c *container.Container) apphandlers.CoreMediaItemDataHandlers {
		factory := container.MustGet[factories.MediaDataFactory](c)
		coreServices := container.MustGet[services.CoreUserMediaItemDataServices](c)
		return factory.CreateCoreDataHandlers(coreServices)
	})
	// (Handlers) User UserMediaItemData Handlers
	container.RegisterFactory[apphandlers.UserMediaItemDataHandlers](c, func(c *container.Container) apphandlers.UserMediaItemDataHandlers {
		factory := container.MustGet[factories.MediaDataFactory](c)
		userServices := container.MustGet[services.UserMediaItemDataServices](c)
		coreHandlers := container.MustGet[apphandlers.CoreMediaItemDataHandlers](c)
		return factory.CreateUserDataHandlers(userServices, coreHandlers)
	})
	// (Handlers) Client UserMediaItemData Handlers
	container.RegisterFactory[apphandlers.ClientMediaItemDataHandlers](c, func(c *container.Container) apphandlers.ClientMediaItemDataHandlers {
		factory := container.MustGet[factories.MediaDataFactory](c)
		dataServices := container.MustGet[services.ClientUserMediaItemDataServices](c)
		userDataHandlers := container.MustGet[apphandlers.UserMediaItemDataHandlers](c)
		return factory.CreateClientDataHandlers(dataServices, userDataHandlers)
	})

	// Register individual media type handlers for direct use in routes

	// Movie handlers
	container.RegisterFactory[handlers.UserMediaItemDataHandler[*mediatypes.Movie]](c, func(c *container.Container) handlers.UserMediaItemDataHandler[*mediatypes.Movie] {
		userHandlers := container.MustGet[apphandlers.UserMediaItemDataHandlers](c)
		return userHandlers.MovieUserDataHandler()
	})

	container.RegisterFactory[handlers.ClientUserMediaItemDataHandler[*mediatypes.Movie]](c, func(c *container.Container) handlers.ClientUserMediaItemDataHandler[*mediatypes.Movie] {
		clientHandlers := container.MustGet[apphandlers.ClientMediaItemDataHandlers](c)
		return clientHandlers.MovieClientDataHandler()
	})

	// Series handlers
	container.RegisterFactory[handlers.UserMediaItemDataHandler[*mediatypes.Series]](c, func(c *container.Container) handlers.UserMediaItemDataHandler[*mediatypes.Series] {
		userHandlers := container.MustGet[apphandlers.UserMediaItemDataHandlers](c)
		return userHandlers.SeriesUserDataHandler()
	})

	container.RegisterFactory[handlers.ClientUserMediaItemDataHandler[*mediatypes.Series]](c, func(c *container.Container) handlers.ClientUserMediaItemDataHandler[*mediatypes.Series] {
		clientHandlers := container.MustGet[apphandlers.ClientMediaItemDataHandlers](c)
		return clientHandlers.SeriesClientDataHandler()
	})

	// Track handlers
	container.RegisterFactory[handlers.UserMediaItemDataHandler[*mediatypes.Track]](c, func(c *container.Container) handlers.UserMediaItemDataHandler[*mediatypes.Track] {
		userHandlers := container.MustGet[apphandlers.UserMediaItemDataHandlers](c)
		return userHandlers.TrackUserDataHandler()
	})

	container.RegisterFactory[handlers.ClientUserMediaItemDataHandler[*mediatypes.Track]](c, func(c *container.Container) handlers.ClientUserMediaItemDataHandler[*mediatypes.Track] {
		clientHandlers := container.MustGet[apphandlers.ClientMediaItemDataHandlers](c)
		return clientHandlers.TrackClientDataHandler()
	})

	// Album handlers
	container.RegisterFactory[handlers.UserMediaItemDataHandler[*mediatypes.Album]](c, func(c *container.Container) handlers.UserMediaItemDataHandler[*mediatypes.Album] {
		userHandlers := container.MustGet[apphandlers.UserMediaItemDataHandlers](c)
		return userHandlers.AlbumUserDataHandler()
	})

	container.RegisterFactory[handlers.ClientUserMediaItemDataHandler[*mediatypes.Album]](c, func(c *container.Container) handlers.ClientUserMediaItemDataHandler[*mediatypes.Album] {
		clientHandlers := container.MustGet[apphandlers.ClientMediaItemDataHandlers](c)
		return clientHandlers.AlbumClientDataHandler()
	})

	// Artist handlers
	container.RegisterFactory[handlers.UserMediaItemDataHandler[*mediatypes.Artist]](c, func(c *container.Container) handlers.UserMediaItemDataHandler[*mediatypes.Artist] {
		userHandlers := container.MustGet[apphandlers.UserMediaItemDataHandlers](c)
		return userHandlers.ArtistUserDataHandler()
	})

	container.RegisterFactory[handlers.ClientUserMediaItemDataHandler[*mediatypes.Artist]](c, func(c *container.Container) handlers.ClientUserMediaItemDataHandler[*mediatypes.Artist] {
		clientHandlers := container.MustGet[apphandlers.ClientMediaItemDataHandlers](c)
		return clientHandlers.ArtistClientDataHandler()
	})
}
