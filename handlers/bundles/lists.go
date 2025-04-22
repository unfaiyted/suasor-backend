package bundles

import (
	mediatypes "suasor/clients/media/types"
	"suasor/handlers"
)

type MediaListHandlers interface {
	CoreMediaListHandlers
	UserMediaListHandlers
}

type CoreMediaListHandlers interface {
	CorePlaylistsHandler() handlers.CoreListHandler[*mediatypes.Playlist]
	CoreCollectionsHandler() handlers.CoreListHandler[*mediatypes.Collection]
}

type UserMediaListHandlers interface {
	UserPlaylistsHandler() handlers.UserListHandler[*mediatypes.Playlist]
	UserCollectionsHandler() handlers.UserListHandler[*mediatypes.Collection]
}

func NewUserMediaListHandlers(
	userPlaylistHandler handlers.UserListHandler[*mediatypes.Playlist],
	userCollectionHandler handlers.UserListHandler[*mediatypes.Collection],
) UserMediaListHandlers {
	return &UserMediaListHandlersImpl{
		userPlaylistHandler:   userPlaylistHandler,
		userCollectionHandler: userCollectionHandler,
	}
}

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

type MediaListHandlersImpl struct {
	corePlaylistHandler   handlers.CoreListHandler[*mediatypes.Playlist]
	coreCollectionHandler handlers.CoreListHandler[*mediatypes.Collection]
	userPlaylistHandler   handlers.UserListHandler[*mediatypes.Playlist]
	userCollectionHandler handlers.UserListHandler[*mediatypes.Collection]
}

func NewMediaListHandlers(
	corePlaylistHandler handlers.CoreListHandler[*mediatypes.Playlist],
	coreCollectionHandler handlers.CoreListHandler[*mediatypes.Collection],
	userPlaylistHandler handlers.UserListHandler[*mediatypes.Playlist],
	userCollectionHandler handlers.UserListHandler[*mediatypes.Collection],
) MediaListHandlers {
	return &MediaListHandlersImpl{
		corePlaylistHandler:   corePlaylistHandler,
		coreCollectionHandler: coreCollectionHandler,
		userPlaylistHandler:   userPlaylistHandler,
		userCollectionHandler: userCollectionHandler,
	}
}

// Implementation of CoreMediaListHandlers
func (h *MediaListHandlersImpl) CorePlaylistsHandler() handlers.CoreListHandler[*mediatypes.Playlist] {
	return h.corePlaylistHandler
}
func (h *MediaListHandlersImpl) CoreCollectionsHandler() handlers.CoreListHandler[*mediatypes.Collection] {
	return h.coreCollectionHandler
}
func (h *MediaListHandlersImpl) UserPlaylistsHandler() handlers.UserListHandler[*mediatypes.Playlist] {
	return h.userPlaylistHandler
}
func (h *MediaListHandlersImpl) UserCollectionsHandler() handlers.UserListHandler[*mediatypes.Collection] {
	return h.userCollectionHandler
}
