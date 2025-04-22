package bundles

import (
	mediatypes "suasor/client/media/types"
	"suasor/services"
)

// Core-Data-layer services (extend core)
type CoreUserMediaItemDataServices interface {
	MovieCoreService() services.CoreUserMediaItemDataService[*mediatypes.Movie]
	SeriesCoreService() services.CoreUserMediaItemDataService[*mediatypes.Series]
	EpisodeCoreService() services.CoreUserMediaItemDataService[*mediatypes.Episode]
	TrackCoreService() services.CoreUserMediaItemDataService[*mediatypes.Track]
	AlbumCoreService() services.CoreUserMediaItemDataService[*mediatypes.Album]
	ArtistCoreService() services.CoreUserMediaItemDataService[*mediatypes.Artist]
	CollectionCoreService() services.CoreUserMediaItemDataService[*mediatypes.Collection]
	PlaylistCoreService() services.CoreUserMediaItemDataService[*mediatypes.Playlist]
}

type UserMediaItemDataServices interface {
	MovieDataService() services.UserMediaItemDataService[*mediatypes.Movie]
	SeriesDataService() services.UserMediaItemDataService[*mediatypes.Series]
	EpisodeDataService() services.UserMediaItemDataService[*mediatypes.Episode]
	TrackDataService() services.UserMediaItemDataService[*mediatypes.Track]
	AlbumDataService() services.UserMediaItemDataService[*mediatypes.Album]
	ArtistDataService() services.UserMediaItemDataService[*mediatypes.Artist]
	CollectionDataService() services.UserMediaItemDataService[*mediatypes.Collection]
	PlaylistDataService() services.UserMediaItemDataService[*mediatypes.Playlist]
}

type ClientUserMediaItemDataServices interface {
	MovieDataService() services.ClientUserMediaItemDataService[*mediatypes.Movie]
	SeriesDataService() services.ClientUserMediaItemDataService[*mediatypes.Series]
	EpisodeDataService() services.ClientUserMediaItemDataService[*mediatypes.Episode]
	TrackDataService() services.ClientUserMediaItemDataService[*mediatypes.Track]
	AlbumDataService() services.ClientUserMediaItemDataService[*mediatypes.Album]
	ArtistDataService() services.ClientUserMediaItemDataService[*mediatypes.Artist]
	CollectionDataService() services.ClientUserMediaItemDataService[*mediatypes.Collection]
	PlaylistDataService() services.ClientUserMediaItemDataService[*mediatypes.Playlist]
}
