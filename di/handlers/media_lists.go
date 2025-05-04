// app/di/handlers/media_lists.go
package handlers

import (
	"context"
	mediatypes "suasor/clients/media/types"
	"suasor/di/container"
	"suasor/handlers"
	apphandlers "suasor/handlers/bundles"
	"suasor/services"
)

// RegisterMediaListHandlers registers handlers for media lists
func RegisterMediaListHandlers(ctx context.Context, c *container.Container) {

	registerMediaListHandler[*mediatypes.Playlist](c)
	registerMediaListHandler[*mediatypes.Collection](c)

	// Register the UserMediaListHandlers implementation
	container.RegisterFactory[apphandlers.UserMediaListHandlers](c, func(c *container.Container) apphandlers.UserMediaListHandlers {
		userPlaylistHandler := container.MustGet[handlers.UserListHandler[*mediatypes.Playlist]](c)
		userCollectionHandler := container.MustGet[handlers.UserListHandler[*mediatypes.Collection]](c)

		return apphandlers.NewUserMediaListHandlers(
			userPlaylistHandler,
			userCollectionHandler,
		)
	})

	// Register the complete MediaListHandlers implementation
	container.RegisterFactory[apphandlers.MediaListHandlers](c, func(c *container.Container) apphandlers.MediaListHandlers {
		corePlaylistHandler := container.MustGet[handlers.CoreListHandler[*mediatypes.Playlist]](c)
		coreCollectionHandler := container.MustGet[handlers.CoreListHandler[*mediatypes.Collection]](c)
		userPlaylistHandler := container.MustGet[handlers.UserListHandler[*mediatypes.Playlist]](c)
		userCollectionHandler := container.MustGet[handlers.UserListHandler[*mediatypes.Collection]](c)

		return apphandlers.NewMediaListHandlers(
			corePlaylistHandler,
			coreCollectionHandler,
			userPlaylistHandler,
			userCollectionHandler,
		)
	})
}

func registerMediaListHandler[T mediatypes.ListData](c *container.Container) {

	container.RegisterFactory[handlers.CoreListHandler[T]](c, func(c *container.Container) handlers.CoreListHandler[T] {
		coreItemHandler := container.MustGet[handlers.CoreMediaItemHandler[T]](c)
		listService := container.MustGet[services.CoreListService[T]](c)
		return handlers.NewCoreListHandler[T](coreItemHandler, listService)
	})

	container.RegisterFactory[handlers.UserListHandler[T]](c, func(c *container.Container) handlers.UserListHandler[T] {
		coreHandler := container.MustGet[handlers.CoreListHandler[T]](c)
		itemService := container.MustGet[services.UserMediaItemService[T]](c)
		listService := container.MustGet[services.UserListService[T]](c)
		syncService := container.MustGet[services.ListSyncService[T]](c)
		return handlers.NewUserListHandler[T](coreHandler, itemService, listService, syncService)
	})
}
