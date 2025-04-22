package bundles

import (
	mediatypes "suasor/clients/media/types"
	clienttypes "suasor/clients/types"
	"suasor/services"
)

// CoreMediaItemServices defines the core services for media items
type CoreMediaItemServices interface {
	MovieCoreService() services.CoreMediaItemService[*mediatypes.Movie]
	SeriesCoreService() services.CoreMediaItemService[*mediatypes.Series]
	EpisodeCoreService() services.CoreMediaItemService[*mediatypes.Episode]
	SeasonCoreService() services.CoreMediaItemService[*mediatypes.Season]
	TrackCoreService() services.CoreMediaItemService[*mediatypes.Track]
	AlbumCoreService() services.CoreMediaItemService[*mediatypes.Album]
	ArtistCoreService() services.CoreMediaItemService[*mediatypes.Artist]
	CollectionCoreService() services.CoreMediaItemService[*mediatypes.Collection]
	PlaylistCoreService() services.CoreMediaItemService[*mediatypes.Playlist]
}

// UserMediaItemServices defines the user-specific services for media items
type UserMediaItemServices interface {
	MovieUserService() services.UserMediaItemService[*mediatypes.Movie]
	SeriesUserService() services.UserMediaItemService[*mediatypes.Series]
	EpisodeUserService() services.UserMediaItemService[*mediatypes.Episode]
	TrackUserService() services.UserMediaItemService[*mediatypes.Track]
	SeasonUserService() services.UserMediaItemService[*mediatypes.Season]
	AlbumUserService() services.UserMediaItemService[*mediatypes.Album]
	ArtistUserService() services.UserMediaItemService[*mediatypes.Artist]
	CollectionUserService() services.UserMediaItemService[*mediatypes.Collection]
	PlaylistUserService() services.UserMediaItemService[*mediatypes.Playlist]
}

// ClientMediaItemServices defines the client-specific services for media items
type ClientMediaItemServices[T clienttypes.ClientMediaConfig] interface {
	MovieClientService() services.ClientMediaItemService[T, *mediatypes.Movie]
	SeriesClientService() services.ClientMediaItemService[T, *mediatypes.Series]
	EpisodeClientService() services.ClientMediaItemService[T, *mediatypes.Episode]
	SeasonClientService() services.ClientMediaItemService[T, *mediatypes.Season]
	TrackClientService() services.ClientMediaItemService[T, *mediatypes.Track]
	AlbumClientService() services.ClientMediaItemService[T, *mediatypes.Album]
	ArtistClientService() services.ClientMediaItemService[T, *mediatypes.Artist]
	CollectionClientService() services.ClientMediaItemService[T, *mediatypes.Collection]
	PlaylistClientService() services.ClientMediaItemService[T, *mediatypes.Playlist]
}

type CoreListServices interface {
	CoreCollectionService() services.CoreListService[*mediatypes.Collection]
	CorePlaylistService() services.CoreListService[*mediatypes.Playlist]
}

type UserListServices interface {
	UserCollectionService() services.UserListService[*mediatypes.Collection]
	UserPlaylistService() services.UserListService[*mediatypes.Playlist]
}

type ClientListServices interface {
	EmbyClientCollectionService() services.ClientListService[*clienttypes.EmbyConfig, *mediatypes.Collection]
	EmbyClientPlaylistService() services.ClientListService[*clienttypes.EmbyConfig, *mediatypes.Playlist]
	JellyfinClientCollectionService() services.ClientListService[*clienttypes.JellyfinConfig, *mediatypes.Collection]
	JellyfinClientPlaylistService() services.ClientListService[*clienttypes.JellyfinConfig, *mediatypes.Playlist]
	PlexClientCollectionService() services.ClientListService[*clienttypes.PlexConfig, *mediatypes.Collection]
	PlexClientPlaylistService() services.ClientListService[*clienttypes.PlexConfig, *mediatypes.Playlist]
	SubsonicClientCollectionService() services.ClientListService[*clienttypes.SubsonicConfig, *mediatypes.Collection]
	SubsonicClientPlaylistService() services.ClientListService[*clienttypes.SubsonicConfig, *mediatypes.Playlist]
}

// MediaServices interface for media-related services
type PeopleServices interface {
	PersonService() *services.PersonService
	CreditService() *services.CreditService
}

type peopleServices struct {
	personService *services.PersonService
	creditService *services.CreditService
}

func NewPeopleServices(personService *services.PersonService, creditService *services.CreditService) *peopleServices {
	return &peopleServices{
		personService: personService,
		creditService: creditService,
	}
}

func (s *peopleServices) PersonService() *services.PersonService {
	return s.personService
}

func (s *peopleServices) CreditService() *services.CreditService {
	return s.creditService
}
