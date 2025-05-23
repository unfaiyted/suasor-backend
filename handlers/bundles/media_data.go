package bundles

import (
	mediatypes "suasor/clients/media/types"
	"suasor/handlers"
)

type MediaDataHandlers interface {
	CoreMediaItemDataHandlers
	UserMediaItemDataHandlers
	// ClientMediaItemDataHandlers
}

// Core-Data-layer handlers (extend core)
type CoreMediaItemDataHandlers interface {
	MovieCoreDataHandler() handlers.CoreUserMediaItemDataHandler[*mediatypes.Movie]
	SeriesCoreDataHandler() handlers.CoreUserMediaItemDataHandler[*mediatypes.Series]
	EpisodeCoreDataHandler() handlers.CoreUserMediaItemDataHandler[*mediatypes.Episode]
	TrackCoreDataHandler() handlers.CoreUserMediaItemDataHandler[*mediatypes.Track]
	AlbumCoreDataHandler() handlers.CoreUserMediaItemDataHandler[*mediatypes.Album]
	ArtistCoreDataHandler() handlers.CoreUserMediaItemDataHandler[*mediatypes.Artist]
	CollectionCoreDataHandler() handlers.CoreUserMediaItemDataHandler[*mediatypes.Collection]
	PlaylistCoreDataHandler() handlers.CoreUserMediaItemDataHandler[*mediatypes.Playlist]
}

type UserMediaItemDataHandlers interface {
	MovieUserDataHandler() handlers.UserMediaItemDataHandler[*mediatypes.Movie]
	SeriesUserDataHandler() handlers.UserMediaItemDataHandler[*mediatypes.Series]
	EpisodeUserDataHandler() handlers.UserMediaItemDataHandler[*mediatypes.Episode]
	SeasonUserDataHandler() handlers.UserMediaItemDataHandler[*mediatypes.Season]
	TrackUserDataHandler() handlers.UserMediaItemDataHandler[*mediatypes.Track]
	AlbumUserDataHandler() handlers.UserMediaItemDataHandler[*mediatypes.Album]
	ArtistUserDataHandler() handlers.UserMediaItemDataHandler[*mediatypes.Artist]
	CollectionUserDataHandler() handlers.UserMediaItemDataHandler[*mediatypes.Collection]
	PlaylistUserDataHandler() handlers.UserMediaItemDataHandler[*mediatypes.Playlist]
}

// type ClientMediaItemDataHandlers interface {
// 	MovieClientDataHandler() handlers.ClientUserMediaItemDataHandler[*clienttypes.EmbyConfig, *mediatypes.Movie]
// 	SeriesClientDataHandler() handlers.ClientUserMediaItemDataHandler[*clienttypes.EmbyConfig, *mediatypes.Series]
// 	EpisodeClientDataHandler() handlers.ClientUserMediaItemDataHandler[*clienttypes.EmbyConfig, *mediatypes.Episode]
// 	SeasonClientDataHandler() handlers.ClientUserMediaItemDataHandler[*clienttypes.EmbyConfig, *mediatypes.Season]
// 	TrackClientDataHandler() handlers.ClientUserMediaItemDataHandler[*clienttypes.EmbyConfig, *mediatypes.Track]
// 	AlbumClientDataHandler() handlers.ClientUserMediaItemDataHandler[*clienttypes.EmbyConfig, *mediatypes.Album]
// 	ArtistClientDataHandler() handlers.ClientUserMediaItemDataHandler[*clienttypes.EmbyConfig, *mediatypes.Artist]
// 	CollectionClientDataHandler() handlers.ClientUserMediaItemDataHandler[*clienttypes.EmbyConfig, *mediatypes.Collection]
// 	PlaylistClientDataHandler() handlers.ClientUserMediaItemDataHandler[*clienttypes.EmbyConfig, *mediatypes.Playlist]
//
// 	SeriesClientDataHandler() handlers.ClientUserMediaItemDataHandler[*mediatypes.Series]
// 	EpisodeClientDataHandler() handlers.ClientUserMediaItemDataHandler[*mediatypes.Episode]
// 	SeasonClientDataHandler() handlers.ClientUserMediaItemDataHandler[*mediatypes.Season]
// 	TrackClientDataHandler() handlers.ClientUserMediaItemDataHandler[*mediatypes.Track]
// 	AlbumClientDataHandler() handlers.ClientUserMediaItemDataHandler[*mediatypes.Album]
// 	ArtistClientDataHandler() handlers.ClientUserMediaItemDataHandler[*mediatypes.Artist]
// 	CollectionClientDataHandler() handlers.ClientUserMediaItemDataHandler[*mediatypes.Collection]
// 	PlaylistClientDataHandler() handlers.ClientUserMediaItemDataHandler[*mediatypes.Playlist]
// }
