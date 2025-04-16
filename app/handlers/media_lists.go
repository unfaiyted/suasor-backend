package handlers

import (
	"suasor/handlers"
)

type MediaListHandlers interface {
	CoreMediaListHandlers
	UserMediaListHandlers
}

type CoreMediaListHandlers interface {
	CorePlaylistsHandler() *handlers.CorePlaylistHandler
	CoreCollectionsHandler() *handlers.CoreCollectionHandler
}

type UserMediaListHandlers interface {
	UserPlaylistsHandler() *handlers.UserPlaylistHandler
	UserCollectionsHandler() *handlers.UserCollectionHandler
}

