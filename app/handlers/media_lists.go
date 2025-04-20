package handlers

import (
	mediatypes "suasor/client/media/types"
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
