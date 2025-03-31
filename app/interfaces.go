// app/dependencies.go
package app

import (
	mediatypes "suasor/client/media/types"
	"suasor/client/types"
	"suasor/handlers"
	"suasor/repository"
	"suasor/services"
)

// Interface definitions (unchanged)
type ClientServices interface {
	EmbyService() services.ClientService[*types.EmbyConfig]
	JellyfinService() services.ClientService[*types.JellyfinConfig]
	PlexService() services.ClientService[*types.PlexConfig]
	SubsonicService() services.ClientService[*types.SubsonicConfig]
	SonarrService() services.ClientService[*types.SonarrConfig]
	RadarrService() services.ClientService[*types.RadarrConfig]
	LidarrService() services.ClientService[*types.LidarrConfig]
}

type ClientRepositories interface {
	EmbyRepo() repository.ClientRepository[*types.EmbyConfig]
	JellyfinRepo() repository.ClientRepository[*types.JellyfinConfig]
	PlexRepo() repository.ClientRepository[*types.PlexConfig]
	SubsonicRepo() repository.ClientRepository[*types.SubsonicConfig]
	SonarrRepo() repository.ClientRepository[*types.SonarrConfig]
	RadarrRepo() repository.ClientRepository[*types.RadarrConfig]
	LidarrRepo() repository.ClientRepository[*types.LidarrConfig]
}

type ClientMediaServices interface {
	ClientMovieServies
	ClientSeriesServices
	ClientEpisodeServices
	ClientPlaylistServices
}

type ClientAutomationServices interface {
	// TODO: implement
}

type ClientMovieServies interface {
	EmbyMovieService() services.MediaClientMovieService[*types.EmbyConfig]
	JellyfinMovieService() services.MediaClientMovieService[*types.JellyfinConfig]
	PlexMovieService() services.MediaClientMovieService[*types.PlexConfig]
	SubsonicMovieService() services.MediaClientMovieService[*types.SubsonicConfig]
}

type ClientSeriesServices interface {
	// TODO: implement
}

type ClientEpisodeServices interface {
	// TODO: implement
}

type ClientPlaylistServices interface {
	// TODO: implement
}

type MediaItemServices interface {
	MovieService() services.MediaItemService[*mediatypes.Movie]
	SeriesService() services.MediaItemService[*mediatypes.Series]
	EpisodeService() services.MediaItemService[*mediatypes.Episode]
	TrackService() services.MediaItemService[*mediatypes.Track]
	AlbumService() services.MediaItemService[*mediatypes.Album]
	ArtistService() services.MediaItemService[*mediatypes.Artist]
	CollectionService() services.MediaItemService[*mediatypes.Collection]
	PlaylistService() services.MediaItemService[*mediatypes.Playlist]
}

type UserRepositories interface {
	UserRepo() repository.UserRepository
	UserConfigRepo() repository.UserConfigRepository
	SessionRepo() repository.SessionRepository
}

type MediaItemRepositories interface {
	MovieRepo() repository.MediaItemRepository[*mediatypes.Movie]
	SeriesRepo() repository.MediaItemRepository[*mediatypes.Series]
	EpisodeRepo() repository.MediaItemRepository[*mediatypes.Episode]
	TrackRepo() repository.MediaItemRepository[*mediatypes.Track]
	AlbumRepo() repository.MediaItemRepository[*mediatypes.Album]
	ArtistRepo() repository.MediaItemRepository[*mediatypes.Artist]
	CollectionRepo() repository.MediaItemRepository[*mediatypes.Collection]
	PlaylistRepo() repository.MediaItemRepository[*mediatypes.Playlist]
}

type UserServices interface {
	UserService() services.UserService
	UserConfigService() services.UserConfigService
	AuthService() services.AuthService
}

type ClientHandlers interface {
	EmbyHandler() *handlers.ClientHandler[*types.EmbyConfig]
	JellyfinHandler() *handlers.ClientHandler[*types.JellyfinConfig]
	PlexHandler() *handlers.ClientHandler[*types.PlexConfig]
	SubsonicHandler() *handlers.ClientHandler[*types.SubsonicConfig]
	RadarrHandler() *handlers.ClientHandler[*types.RadarrConfig]
	LidarrHandler() *handlers.ClientHandler[*types.LidarrConfig]
	SonarrHandler() *handlers.ClientHandler[*types.SonarrConfig]
}

type MediaItemHandlers interface {
	MovieHandler() *handlers.MediaItemHandler[*mediatypes.Movie]
	SeriesHandler() *handlers.MediaItemHandler[*mediatypes.Series]
	EpisodeHandler() *handlers.MediaItemHandler[*mediatypes.Episode]
	TrackHandler() *handlers.MediaItemHandler[*mediatypes.Track]
	AlbumHandler() *handlers.MediaItemHandler[*mediatypes.Album]
	ArtistHandler() *handlers.MediaItemHandler[*mediatypes.Artist]
	CollectionHandler() *handlers.MediaItemHandler[*mediatypes.Collection]
	PlaylistHandler() *handlers.MediaItemHandler[*mediatypes.Playlist]
}

type ClientMediaHandlers interface {
	ClientMediaMovieHandlers
	ClientMediaSeriesHandlers
	ClientMediaEpisodeHandlers
}

type ClientMediaMovieHandlers interface {
	EmbyMovieHandler() *handlers.MediaClientMovieHandler[*types.EmbyConfig]
	JellyfinMovieHandler() *handlers.MediaClientMovieHandler[*types.JellyfinConfig]
	PlexMovieHandler() *handlers.MediaClientMovieHandler[*types.PlexConfig]
}

type ClientMediaSeriesHandlers interface {
	// EmbySeriesHandler() *handlers.MediaClientSeriesHandler[*types.EmbyConfig]
	// JellyfinSeriesHandler() *handlers.MediaClientSeriesHandler[*types.JellyfinConfig]
	// PlexSeriesHandler() *handlers.MediaClientSeriesHandler[*types.PlexConfig]
}

type ClientMediaEpisodeHandlers interface {
	// EmbyEpisodeHandler() *handlers.MediaClientEpisodeHandler[*types.EmbyConfig]
	// JellyfinEpisodeHandler() *handlers.MediaClientEpisodeHandler[*types.JellyfinConfig]
	// PlexEpisodeHandler() *handlers.MediaClientEpisodeHandler[*types.PlexConfig]
}

type UserHandlers interface {
	AuthHandler() *handlers.AuthHandler
	UserHandler() *handlers.UserHandler
	UserConfigHandler() *handlers.UserConfigHandler
}

type SystemHandlers interface {
	ConfigHandler() *handlers.ConfigHandler
	HealthHandler() *handlers.HealthHandler
}

type SystemServices interface {
	HealthService() services.HealthService
	ConfigService() services.ConfigService
}

type SystemRepositories interface {
	ConfigRepo() repository.ConfigRepository
}
