// app/di/handlers/media_lists.go
package handlers

import (
	"context"
	"suasor/app/container"
	apphandlers "suasor/app/handlers"
	mediatypes "suasor/client/media/types"
	"suasor/handlers"
	"suasor/services"
)

type UserMediaListHandlersImpl struct {
	userPlaylistHandler   handlers.UserListHandler[*mediatypes.Playlist]
	userCollectionHandler handlers.UserListHandler[*mediatypes.Collection]
}

func (h *UserMediaListHandlersImpl) UserPlaylistsHandler() handlers.UserListHandler[*mediatypes.Playlist] {
	return h.userPlaylistHandler
}

func (h *UserMediaListHandlersImpl) UserCollectionsHandler() handlers.UserListHandler[*mediatypes.Collection] {
	return h.userCollectionHandler
}

// MediaListHandlersImpl implements the MediaListHandlers interface
type MediaListHandlersImpl struct {
	corePlaylistHandler   handlers.CoreListHandler[*mediatypes.Playlist]
	coreCollectionHandler handlers.CoreListHandler[*mediatypes.Collection]
	userPlaylistHandler   handlers.UserListHandler[*mediatypes.Playlist]
	userCollectionHandler handlers.UserListHandler[*mediatypes.Collection]
}

// Implementation of CoreMediaListHandlers
func (h *MediaListHandlersImpl) CorePlaylistsHandler() handlers.CoreListHandler[*mediatypes.Playlist] {
	return h.corePlaylistHandler
}

func (h *MediaListHandlersImpl) CoreCollectionsHandler() handlers.CoreListHandler[*mediatypes.Collection] {
	return h.coreCollectionHandler
}

// Implementation of UserMediaListHandlers
func (h *MediaListHandlersImpl) UserPlaylistsHandler() handlers.UserListHandler[*mediatypes.Playlist] {
	return h.userPlaylistHandler
}

func (h *MediaListHandlersImpl) UserCollectionsHandler() handlers.UserListHandler[*mediatypes.Collection] {
	return h.userCollectionHandler
}

// RegisterMediaListHandlers registers handlers for media lists
func RegisterMediaListHandlers(ctx context.Context, c *container.Container) {
	// Register CoreMediaItemHandler for Playlists
	container.RegisterFactory[handlers.CoreMediaItemHandler[*mediatypes.Playlist]](c, func(c *container.Container) handlers.CoreMediaItemHandler[*mediatypes.Playlist] {
		service := container.MustGet[services.CoreMediaItemService[*mediatypes.Playlist]](c)
		return handlers.NewCoreMediaItemHandler[*mediatypes.Playlist](service)
	})

	// Register CoreMediaItemHandler for Collections
	container.RegisterFactory[handlers.CoreMediaItemHandler[*mediatypes.Collection]](c, func(c *container.Container) handlers.CoreMediaItemHandler[*mediatypes.Collection] {
		service := container.MustGet[services.CoreMediaItemService[*mediatypes.Collection]](c)
		return handlers.NewCoreMediaItemHandler[*mediatypes.Collection](service)
	})

	// Register the individual handlers
	container.RegisterFactory[handlers.CoreListHandler[*mediatypes.Playlist]](c, func(c *container.Container) handlers.CoreListHandler[*mediatypes.Playlist] {
		coreItemHandler := container.MustGet[handlers.CoreMediaItemHandler[*mediatypes.Playlist]](c)
		service := container.MustGet[services.CoreListService[*mediatypes.Playlist]](c)
		return handlers.NewCoreListHandler[*mediatypes.Playlist](coreItemHandler, service)
	})

	container.RegisterFactory[handlers.CoreListHandler[*mediatypes.Collection]](c, func(c *container.Container) handlers.CoreListHandler[*mediatypes.Collection] {
		coreItemHandler := container.MustGet[handlers.CoreMediaItemHandler[*mediatypes.Collection]](c)
		service := container.MustGet[services.CoreListService[*mediatypes.Collection]](c)
		return handlers.NewCoreListHandler[*mediatypes.Collection](coreItemHandler, service)
	})

	container.RegisterFactory[handlers.UserListHandler[*mediatypes.Playlist]](c, func(c *container.Container) handlers.UserListHandler[*mediatypes.Playlist] {
		coreHandler := container.MustGet[handlers.CoreListHandler[*mediatypes.Playlist]](c)
		itemService := container.MustGet[services.UserMediaItemService[*mediatypes.Playlist]](c)
		listService := container.MustGet[services.UserListService[*mediatypes.Playlist]](c)

		return handlers.NewUserListHandler(coreHandler, itemService, listService)
	})

	container.RegisterFactory[handlers.UserListHandler[*mediatypes.Collection]](c, func(c *container.Container) handlers.UserListHandler[*mediatypes.Collection] {
		coreHandler := container.MustGet[handlers.CoreListHandler[*mediatypes.Collection]](c)
		itemService := container.MustGet[services.UserMediaItemService[*mediatypes.Collection]](c)
		listService := container.MustGet[services.UserListService[*mediatypes.Collection]](c)

		return handlers.NewUserListHandler[*mediatypes.Collection](coreHandler, itemService, listService)
	})

	// Register the UserMediaListHandlers implementation
	container.RegisterFactory[apphandlers.UserMediaListHandlers](c, func(c *container.Container) apphandlers.UserMediaListHandlers {
		userPlaylistHandler := container.MustGet[handlers.UserListHandler[*mediatypes.Playlist]](c)
		userCollectionHandler := container.MustGet[handlers.UserListHandler[*mediatypes.Collection]](c)

		return &UserMediaListHandlersImpl{
			userPlaylistHandler:   userPlaylistHandler,
			userCollectionHandler: userCollectionHandler,
		}
	})

	// Register the complete MediaListHandlers implementation
	container.RegisterFactory[apphandlers.MediaListHandlers](c, func(c *container.Container) apphandlers.MediaListHandlers {
		corePlaylistHandler := container.MustGet[handlers.CoreListHandler[*mediatypes.Playlist]](c)
		coreCollectionHandler := container.MustGet[handlers.CoreListHandler[*mediatypes.Collection]](c)
		userPlaylistHandler := container.MustGet[handlers.UserListHandler[*mediatypes.Playlist]](c)
		userCollectionHandler := container.MustGet[handlers.UserListHandler[*mediatypes.Collection]](c)

		return &MediaListHandlersImpl{
			corePlaylistHandler:   corePlaylistHandler,
			coreCollectionHandler: coreCollectionHandler,
			userPlaylistHandler:   userPlaylistHandler,
			userCollectionHandler: userCollectionHandler,
		}
	})
}

