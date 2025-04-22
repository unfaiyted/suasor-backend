// app/di/media_data.go
package di

import (
	"context"
	"gorm.io/gorm"
	"suasor/client"
	mediatypes "suasor/client/media/types"
	clienttypes "suasor/client/types"
	"suasor/container"
	"suasor/di/factories"
	"suasor/handlers"
	handlerbundles "suasor/handlers/bundles"
	"suasor/repository"
	repobundles "suasor/repository/bundles"
	"suasor/services"
	svcbundles "suasor/services/bundles"
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

	// --- repobundles FACTORIES --- //

	//  MediaItem
	//  (Repositories) Core MediaItem Repositories
	container.RegisterFactory[repobundles.CoreMediaItemRepositories](c, func(c *container.Container) repobundles.CoreMediaItemRepositories {
		factory := container.MustGet[factories.MediaDataFactory](c)
		return factory.CreateCoreRepositories()
	})
	//  (Repositories) User MediaItem Repositories
	container.RegisterFactory[repobundles.UserMediaItemRepositories](c, func(c *container.Container) repobundles.UserMediaItemRepositories {
		factory := container.MustGet[factories.MediaDataFactory](c)
		return factory.CreateUserRepositories()
	})
	//  (Repositories) Client MediaItem Repositories
	container.RegisterFactory[repobundles.ClientMediaItemRepositories](c, func(c *container.Container) repobundles.ClientMediaItemRepositories {
		factory := container.MustGet[factories.MediaDataFactory](c)
		return factory.CreateClientMediaItemRepositories()
	})

	//  Core Services Registration
	//  (Services) Core MediaItem Services
	container.RegisterFactory[svcbundles.CoreMediaItemServices](c, func(c *container.Container) svcbundles.CoreMediaItemServices {
		factory := container.MustGet[factories.MediaDataFactory](c)
		repos := container.MustGet[repobundles.CoreMediaItemRepositories](c)
		return factory.CreateCoreServices(repos)
	})

	//  UserMediaItemData
	//  (Repositories) Core UserMediaItemData Repositories
	container.RegisterFactory[repobundles.CoreUserMediaItemDataRepositories](c, func(c *container.Container) repobundles.CoreUserMediaItemDataRepositories {
		factory := container.MustGet[factories.MediaDataFactory](c)
		return factory.CreateCoreDataRepositories()
	})
	//  (Repositories) User UserMediaItemData Repositories
	container.RegisterFactory[repobundles.UserMediaDataRepositories](c, func(c *container.Container) repobundles.UserMediaDataRepositories {
		factory := container.MustGet[factories.MediaDataFactory](c)
		return factory.CreateUserDataRepositories()
	})
	//  (Repositories) Client UserMediaItemData Repositories
	container.RegisterFactory[repobundles.ClientUserMediaDataRepositories](c, func(c *container.Container) repobundles.ClientUserMediaDataRepositories {
		factory := container.MustGet[factories.MediaDataFactory](c)
		return factory.CreateClientDataRepositories()
	})

	// --- SERVICE FACTORIES --- //

	//  MediaItem
	//  (Services) Core MediaItem Services
	container.RegisterFactory[svcbundles.CoreMediaItemServices](c, func(c *container.Container) svcbundles.CoreMediaItemServices {
		factory := container.MustGet[factories.MediaDataFactory](c)
		repos := container.MustGet[repobundles.CoreMediaItemRepositories](c)
		return factory.CreateCoreServices(repos)
	})
	//  (Services) User MediaItem Services
	container.RegisterFactory[svcbundles.UserMediaItemServices](c, func(c *container.Container) svcbundles.UserMediaItemServices {
		factory := container.MustGet[factories.MediaDataFactory](c)
		coreServices := container.MustGet[svcbundles.CoreMediaItemServices](c)
		userRepos := container.MustGet[repobundles.UserMediaItemRepositories](c)
		return factory.CreateUserServices(coreServices, userRepos)
	})
	//  (Services) Client MediaItem Services
	container.RegisterFactory[svcbundles.ClientMediaItemServices[clienttypes.ClientMediaConfig]](c, func(c *container.Container) svcbundles.ClientMediaItemServices[clienttypes.ClientMediaConfig] {
		factory := container.MustGet[factories.MediaDataFactory](c)
		coreServices := container.MustGet[svcbundles.CoreMediaItemServices](c)
		clientRepo := container.MustGet[repository.ClientRepository[clienttypes.ClientMediaConfig]](c)
		clientItemRepos := container.MustGet[repobundles.ClientMediaItemRepositories](c)
		return factory.CreateClientServices(coreServices, clientRepo, clientItemRepos)
	})
	//  UserMediaItemData
	//  (Services) Core UserMediaItemData Services
	container.RegisterFactory[svcbundles.CoreUserMediaItemDataServices](c, func(c *container.Container) svcbundles.CoreUserMediaItemDataServices {
		factory := container.MustGet[factories.MediaDataFactory](c)
		repos := container.MustGet[repobundles.CoreMediaItemRepositories](c)
		return factory.CreateCoreDataServices(repos)
	})
	//  (Services) User UserMediaItemData Services
	container.RegisterFactory[svcbundles.UserMediaItemDataServices](c, func(c *container.Container) svcbundles.UserMediaItemDataServices {
		factory := container.MustGet[factories.MediaDataFactory](c)
		coreDataServices := container.MustGet[svcbundles.CoreUserMediaItemDataServices](c)
		userRepos := container.MustGet[repobundles.UserMediaDataRepositories](c)
		return factory.CreateUserDataServices(coreDataServices, userRepos)
	})

	// (Services) Client UserMediaItemData Services
	container.RegisterFactory[svcbundles.ClientUserMediaItemDataServices](c, func(c *container.Container) svcbundles.ClientUserMediaItemDataServices {
		factory := container.MustGet[factories.MediaDataFactory](c)
		coreDataServices := container.MustGet[svcbundles.UserMediaItemDataServices](c)
		clientRepos := container.MustGet[repobundles.ClientUserMediaDataRepositories](c)
		return factory.CreateClientDataServices(coreDataServices, clientRepos)
	})

	// --- HANDLER --- //
	// (Handlers) Core MediaItem Handlers
	container.RegisterFactory[handlerbundles.CoreMediaItemHandlers](c, func(c *container.Container) handlerbundles.CoreMediaItemHandlers {
		factory := container.MustGet[factories.MediaDataFactory](c)
		coreServices := container.MustGet[svcbundles.CoreMediaItemServices](c)
		return factory.CreateCoreMediaItemHandlers(coreServices)
	})
	// (Handlers) User MediaItem Handlers
	container.RegisterFactory[handlerbundles.UserMediaItemHandlers](c, func(c *container.Container) apphandlers.UserMediaItemHandlers {
		factory := container.MustGet[factories.MediaDataFactory](c)
		userServices := container.MustGet[svcbundles.UserMediaItemServices](c)
		coreHandlers := container.MustGet[handlerbundles.CoreMediaItemHandlers](c)
		return factory.CreateUserMediaItemHandlers(userServices, coreHandlers)
	})
	// (Handlers) Client MediaItem Handlers
	container.RegisterFactory[handlerbundles.ClientMediaItemHandlers[clienttypes.ClientMediaConfig]](c, func(c *container.Container) apphandlers.ClientMediaItemHandlers[clienttypes.ClientMediaConfig] {
		factory := container.MustGet[factories.MediaDataFactory](c)
		clientServices := container.MustGet[svcbundles.ClientMediaItemServices[clienttypes.ClientMediaConfig]](c)
		userMediaItemServices := container.MustGet[svcbundles.UserMediaItemServices](c)
		coreHandlers := container.MustGet[handlerbundles.UserMediaItemHandlers](c)
		return factory.CreateClientMediaItemHandlers(clientServices, userMediaItemServices, coreHandlers)
	})
	// (Handlers) Core UserMediaItemData Handlers
	container.RegisterFactory[handlerbundles.CoreMediaItemDataHandlers](c, func(c *container.Container) apphandlers.CoreMediaItemDataHandlers {
		factory := container.MustGet[factories.MediaDataFactory](c)
		coreServices := container.MustGet[svcbundles.CoreUserMediaItemDataServices](c)
		return factory.CreateCoreDataHandlers(coreServices)
	})
	// (Handlers) User UserMediaItemData Handlers
	container.RegisterFactory[handlerbundles.UserMediaItemDataHandlers](c, func(c *container.Container) apphandlers.UserMediaItemDataHandlers {
		factory := container.MustGet[factories.MediaDataFactory](c)
		userServices := container.MustGet[svcbundles.UserMediaItemDataServices](c)
		coreHandlers := container.MustGet[handlerbundles.CoreMediaItemDataHandlers](c)
		return factory.CreateUserDataHandlers(userServices, coreHandlers)
	})
	// (Handlers) Client UserMediaItemData Handlers
	container.RegisterFactory[handlerbundles.ClientMediaItemDataHandlers](c, func(c *container.Container) handlerbundles.ClientMediaItemDataHandlers {
		factory := container.MustGet[factories.MediaDataFactory](c)
		dataServices := container.MustGet[svcbundles.ClientUserMediaItemDataServices](c)
		userDataHandlers := container.MustGet[handlerbundles.UserMediaItemDataHandlers](c)
		return factory.CreateClientDataHandlers(dataServices, userDataHandlers)
	})

	// Register individual media type handlers for direct use in routes

	// Movie handlers
	container.RegisterFactory[handlers.UserMediaItemDataHandler[*mediatypes.Movie]](c, func(c *container.Container) handlers.UserMediaItemDataHandler[*mediatypes.Movie] {
		userHandlers := container.MustGet[handlerbundles.UserMediaItemDataHandlers](c)
		return userHandlers.MovieUserDataHandler()
	})

	container.RegisterFactory[handlers.ClientUserMediaItemDataHandler[*mediatypes.Movie]](c, func(c *container.Container) handlers.ClientUserMediaItemDataHandler[*mediatypes.Movie] {
		clientHandlers := container.MustGet[handlerbundles.ClientMediaItemDataHandlers](c)
		return clientHandlers.MovieClientDataHandler()
	})

	// Series handlers
	container.RegisterFactory[handlers.UserMediaItemDataHandler[*mediatypes.Series]](c, func(c *container.Container) handlers.UserMediaItemDataHandler[*mediatypes.Series] {
		userHandlers := container.MustGet[handlerbundles.UserMediaItemDataHandlers](c)
		return userHandlers.SeriesUserDataHandler()
	})

	container.RegisterFactory[handlers.ClientUserMediaItemDataHandler[*mediatypes.Series]](c, func(c *container.Container) handlers.ClientUserMediaItemDataHandler[*mediatypes.Series] {
		clientHandlers := container.MustGet[handlerbundles.ClientMediaItemDataHandlers](c)
		return clientHandlers.SeriesClientDataHandler()
	})

	// Track handlers
	container.RegisterFactory[handlers.UserMediaItemDataHandler[*mediatypes.Track]](c, func(c *container.Container) handlers.UserMediaItemDataHandler[*mediatypes.Track] {
		userHandlers := container.MustGet[handlerbundles.UserMediaItemDataHandlers](c)
		return userHandlers.TrackUserDataHandler()
	})

	container.RegisterFactory[handlers.ClientUserMediaItemDataHandler[*mediatypes.Track]](c, func(c *container.Container) handlers.ClientUserMediaItemDataHandler[*mediatypes.Track] {
		clientHandlers := container.MustGet[handlerbundles.ClientMediaItemDataHandlers](c)
		return clientHandlers.TrackClientDataHandler()
	})

	// Album handlers
	container.RegisterFactory[handlers.UserMediaItemDataHandler[*mediatypes.Album]](c, func(c *container.Container) handlers.UserMediaItemDataHandler[*mediatypes.Album] {
		userHandlers := container.MustGet[handlerbundles.UserMediaItemDataHandlers](c)
		return userHandlers.AlbumUserDataHandler()
	})

	container.RegisterFactory[handlers.ClientUserMediaItemDataHandler[*mediatypes.Album]](c, func(c *container.Container) handlers.ClientUserMediaItemDataHandler[*mediatypes.Album] {
		clientHandlers := container.MustGet[handlerbundles.ClientMediaItemDataHandlers](c)
		return clientHandlers.AlbumClientDataHandler()
	})

	// Artist handlers
	container.RegisterFactory[handlers.UserMediaItemDataHandler[*mediatypes.Artist]](c, func(c *container.Container) handlers.UserMediaItemDataHandler[*mediatypes.Artist] {
		userHandlers := container.MustGet[handlerbundles.UserMediaItemDataHandlers](c)
		return userHandlers.ArtistUserDataHandler()
	})

	container.RegisterFactory[handlers.ClientUserMediaItemDataHandler[*mediatypes.Artist]](c, func(c *container.Container) handlers.ClientUserMediaItemDataHandler[*mediatypes.Artist] {
		clientHandlers := container.MustGet[handlerbundles.ClientMediaItemDataHandlers](c)
		return clientHandlers.ArtistClientDataHandler()
	})
}
