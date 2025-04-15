// app/interfaces.go
package app

import (
	mediatypes "suasor/client/media/types"
	"suasor/client/types"
	"suasor/handlers"
	"suasor/repository"
	"suasor/services"
	"suasor/services/jobs"
	"suasor/services/jobs/recommendation"
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
	ClaudeService() services.ClientService[*types.ClaudeConfig]
	OpenAIService() services.ClientService[*types.OpenAIConfig]
	OllamaService() services.ClientService[*types.OllamaConfig]
	AllServices() map[string]services.ClientService[types.ClientConfig]
}

// ClientRepositories provides access to client repositories
// These repositories store client configurations by client type
type ClientRepositories interface {
	EmbyRepo() repository.ClientRepository[*types.EmbyConfig]
	JellyfinRepo() repository.ClientRepository[*types.JellyfinConfig]
	PlexRepo() repository.ClientRepository[*types.PlexConfig]
	SubsonicRepo() repository.ClientRepository[*types.SubsonicConfig]
	SonarrRepo() repository.ClientRepository[*types.SonarrConfig]
	RadarrRepo() repository.ClientRepository[*types.RadarrConfig]
	LidarrRepo() repository.ClientRepository[*types.LidarrConfig]
	ClaudeRepo() repository.ClientRepository[*types.ClaudeConfig]
	OpenAIRepo() repository.ClientRepository[*types.OpenAIConfig]
	OllamaRepo() repository.ClientRepository[*types.OllamaConfig]
}

type RepositoryCollections interface {
	ClientRepositories() repository.ClientRepositoryCollection
}

// Three-pronged repository interfaces
type CoreMediaItemRepositories interface {
	MovieRepo() repository.MediaItemRepository[*mediatypes.Movie]
	SeriesRepo() repository.MediaItemRepository[*mediatypes.Series]
	EpisodeRepo() repository.MediaItemRepository[*mediatypes.Episode]
	TrackRepo() repository.MediaItemRepository[*mediatypes.Track]
	AlbumRepo() repository.MediaItemRepository[*mediatypes.Album]
	ArtistRepo() repository.MediaItemRepository[*mediatypes.Artist]
	CollectionRepo() repository.MediaItemRepository[*mediatypes.Collection]
	PlaylistRepo() repository.MediaItemRepository[*mediatypes.Playlist]
}

type CoreUserMediaItemDataRepositories interface {
	MovieDataRepo() repository.CoreUserMediaItemDataRepository[*mediatypes.Movie]
	SeriesDataRepo() repository.CoreUserMediaItemDataRepository[*mediatypes.Series]
	EpisodeDataRepo() repository.CoreUserMediaItemDataRepository[*mediatypes.Episode]
	TrackDataRepo() repository.CoreUserMediaItemDataRepository[*mediatypes.Track]
	AlbumDataRepo() repository.CoreUserMediaItemDataRepository[*mediatypes.Album]
	ArtistDataRepo() repository.CoreUserMediaItemDataRepository[*mediatypes.Artist]
	CollectionDataRepo() repository.CoreUserMediaItemDataRepository[*mediatypes.Collection]
	PlaylistDataRepo() repository.CoreUserMediaItemDataRepository[*mediatypes.Playlist]
}

type UserRepositoryFactories interface {
	MovieUserRepo() repository.UserMediaItemRepository[*mediatypes.Movie]
	SeriesUserRepo() repository.UserMediaItemRepository[*mediatypes.Series]
	EpisodeUserRepo() repository.UserMediaItemRepository[*mediatypes.Episode]
	TrackUserRepo() repository.UserMediaItemRepository[*mediatypes.Track]
	AlbumUserRepo() repository.UserMediaItemRepository[*mediatypes.Album]
	ArtistUserRepo() repository.UserMediaItemRepository[*mediatypes.Artist]
	CollectionUserRepo() repository.UserMediaItemRepository[*mediatypes.Collection]
	PlaylistUserRepo() repository.UserMediaItemRepository[*mediatypes.Playlist]
}

type ClientUserDataRepositories interface {
	MovieDataRepo() repository.ClientUserMediaItemDataRepository[*mediatypes.Movie]
	SeriesDataRepo() repository.ClientUserMediaItemDataRepository[*mediatypes.Series]
	EpisodeDataRepo() repository.ClientUserMediaItemDataRepository[*mediatypes.Episode]
	TrackDataRepo() repository.ClientUserMediaItemDataRepository[*mediatypes.Track]
	AlbumDataRepo() repository.ClientUserMediaItemDataRepository[*mediatypes.Album]
	ArtistDataRepo() repository.ClientUserMediaItemDataRepository[*mediatypes.Artist]
	CollectionDataRepo() repository.ClientUserMediaItemDataRepository[*mediatypes.Collection]
	PlaylistDataRepo() repository.ClientUserMediaItemDataRepository[*mediatypes.Playlist]
}

type ClientRepositoryFactories interface {
	MovieClientRepo() repository.ClientMediaItemRepository[*mediatypes.Movie]
	SeriesClientRepo() repository.ClientMediaItemRepository[*mediatypes.Series]
	EpisodeClientRepo() repository.ClientMediaItemRepository[*mediatypes.Episode]
	TrackClientRepo() repository.ClientMediaItemRepository[*mediatypes.Track]
	AlbumClientRepo() repository.ClientMediaItemRepository[*mediatypes.Album]
	ArtistClientRepo() repository.ClientMediaItemRepository[*mediatypes.Artist]
	CollectionClientRepo() repository.ClientMediaItemRepository[*mediatypes.Collection]
	PlaylistClientRepo() repository.ClientMediaItemRepository[*mediatypes.Playlist]
}

type UserDataFactories interface {
	MovieDataRepo() repository.UserMediaItemDataRepository[*mediatypes.Movie]
	SeriesDataRepo() repository.UserMediaItemDataRepository[*mediatypes.Series]
	EpisodeDataRepo() repository.UserMediaItemDataRepository[*mediatypes.Episode]
	TrackDataRepo() repository.UserMediaItemDataRepository[*mediatypes.Track]
	AlbumDataRepo() repository.UserMediaItemDataRepository[*mediatypes.Album]
	ArtistDataRepo() repository.UserMediaItemDataRepository[*mediatypes.Artist]
	CollectionDataRepo() repository.UserMediaItemDataRepository[*mediatypes.Collection]
	PlaylistDataRepo() repository.UserMediaItemDataRepository[*mediatypes.Playlist]
}

type ClientMediaServices interface {
	ClientMovieServies
	ClientSeriesServices
	ClientEpisodeServices
	ClientMusicServices
	ClientPlaylistServices
}

type ClientAutomationServices interface {
	// TODO: implement
}

type ClientMovieServies interface {
	EmbyMovieService() services.ClientMediaMovieService[*types.EmbyConfig]
	JellyfinMovieService() services.ClientMediaMovieService[*types.JellyfinConfig]
	PlexMovieService() services.ClientMediaMovieService[*types.PlexConfig]
	SubsonicMovieService() services.ClientMediaMovieService[*types.SubsonicConfig]
}

type ClientSeriesServices interface {
	EmbySeriesService() services.ClientMediaSeriesService[*types.EmbyConfig]
	JellyfinSeriesService() services.ClientMediaSeriesService[*types.JellyfinConfig]
	PlexSeriesService() services.ClientMediaSeriesService[*types.PlexConfig]
	SubsonicSeriesService() services.ClientMediaSeriesService[*types.SubsonicConfig]
}

type ClientEpisodeServices interface {
	// TODO: implement
}

type ClientMusicServices interface {
	EmbyMusicService() services.ClientMediaMusicService[*types.EmbyConfig]
	JellyfinMusicService() services.ClientMediaMusicService[*types.JellyfinConfig]
	PlexMusicService() services.ClientMediaMusicService[*types.PlexConfig]
	SubsonicMusicService() services.ClientMediaMusicService[*types.SubsonicConfig]
}

type ClientPlaylistServices interface {
	EmbyPlaylistService() services.ClientMediaPlaylistService[*types.EmbyConfig]
	JellyfinPlaylistService() services.ClientMediaPlaylistService[*types.JellyfinConfig]
	PlexPlaylistService() services.ClientMediaPlaylistService[*types.PlexConfig]
	SubsonicPlaylistService() services.ClientMediaPlaylistService[*types.SubsonicConfig]
}

// Three-pronged architecture for service interfaces

// Core-layer services
type CoreMediaItemServices interface {
	MovieCoreService() services.CoreMediaItemService[*mediatypes.Movie]
	SeriesCoreService() services.CoreMediaItemService[*mediatypes.Series]
	EpisodeCoreService() services.CoreMediaItemService[*mediatypes.Episode]
	TrackCoreService() services.CoreMediaItemService[*mediatypes.Track]
	AlbumCoreService() services.CoreMediaItemService[*mediatypes.Album]
	ArtistCoreService() services.CoreMediaItemService[*mediatypes.Artist]
	CollectionCoreService() services.CoreMediaItemService[*mediatypes.Collection]
	PlaylistCoreService() services.CoreMediaItemService[*mediatypes.Playlist]
}

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

// User-layer services (extend core)
type UserMediaItemServices interface {
	MovieUserService() services.UserMediaItemService[*mediatypes.Movie]
	SeriesUserService() services.UserMediaItemService[*mediatypes.Series]
	EpisodeUserService() services.UserMediaItemService[*mediatypes.Episode]
	TrackUserService() services.UserMediaItemService[*mediatypes.Track]
	AlbumUserService() services.UserMediaItemService[*mediatypes.Album]
	ArtistUserService() services.UserMediaItemService[*mediatypes.Artist]
	CollectionUserService() services.UserMediaItemService[*mediatypes.Collection]
	PlaylistUserService() services.UserMediaItemService[*mediatypes.Playlist]
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

// Client-layer services (extend core)
type ClientMediaItemServices interface {
	MovieClientService() services.ClientMediaItemService[*mediatypes.Movie]
	SeriesClientService() services.ClientMediaItemService[*mediatypes.Series]
	EpisodeClientService() services.ClientMediaItemService[*mediatypes.Episode]
	TrackClientService() services.ClientMediaItemService[*mediatypes.Track]
	AlbumClientService() services.ClientMediaItemService[*mediatypes.Album]
	ArtistClientService() services.ClientMediaItemService[*mediatypes.Artist]
	CollectionClientService() services.ClientMediaItemService[*mediatypes.Collection]
	PlaylistClientService() services.ClientMediaItemService[*mediatypes.Playlist]
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

// Specialized collection and playlist services
type MediaCollectionServices interface {
	CoreCollectionService() services.CoreCollectionService
	UserCollectionService() services.UserCollectionService
	ClientCollectionService() services.ClientMediaCollectionService

	CorePlaylistService() services.CoreMediaItemService[*mediatypes.Playlist]
	UserPlaylistService() services.UserMediaItemService[*mediatypes.Playlist]
	ClientPlaylistService() services.ClientMediaItemService[*mediatypes.Playlist]

	// Extended services with additional functionality
	PlaylistService() services.PlaylistService
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

	// User-owned media repositories
	UserMediaPlaylistRepo() repository.UserMediaItemRepository[*mediatypes.Playlist]
}

type UserServices interface {
	UserService() services.UserService
	UserConfigService() services.UserConfigService
	AuthService() services.AuthService
}

type ClientHandlers interface {
	// Media
	EmbyHandler() *handlers.ClientHandler[*types.EmbyConfig]
	JellyfinHandler() *handlers.ClientHandler[*types.JellyfinConfig]
	PlexHandler() *handlers.ClientHandler[*types.PlexConfig]
	SubsonicHandler() *handlers.ClientHandler[*types.SubsonicConfig]
	// Automation
	RadarrHandler() *handlers.ClientHandler[*types.RadarrConfig]
	LidarrHandler() *handlers.ClientHandler[*types.LidarrConfig]
	SonarrHandler() *handlers.ClientHandler[*types.SonarrConfig]
	// AI
	ClaudeHandler() *handlers.ClientHandler[*types.ClaudeConfig]
	OpenAIHandler() *handlers.ClientHandler[*types.OpenAIConfig]
	OllamaHandler() *handlers.ClientHandler[*types.OllamaConfig]
}

// Three-pronged architecture for handler interfaces

// Core-layer handlers
type CoreMediaItemHandlers interface {
	MovieCoreHandler() *handlers.CoreUserMediaItemDataHandler[*mediatypes.Movie]
	SeriesCoreHandler() *handlers.CoreUserMediaItemDataHandler[*mediatypes.Series]
	EpisodeCoreHandler() *handlers.CoreUserMediaItemDataHandler[*mediatypes.Episode]
	TrackCoreHandler() *handlers.CoreUserMediaItemDataHandler[*mediatypes.Track]
	AlbumCoreHandler() *handlers.CoreUserMediaItemDataHandler[*mediatypes.Album]
	ArtistCoreHandler() *handlers.CoreUserMediaItemDataHandler[*mediatypes.Artist]
	CollectionCoreHandler() *handlers.CoreUserMediaItemDataHandler[*mediatypes.Collection]
	PlaylistCoreHandler() *handlers.CoreUserMediaItemDataHandler[*mediatypes.Playlist]
}

// Core-Data-layer handlers (extend core)
type CoreMediaItemDataHandlers interface {
	MovieCoreDataHandler() *handlers.CoreUserMediaItemDataHandler[*mediatypes.Movie]
	SeriesCoreDataHandler() *handlers.CoreUserMediaItemDataHandler[*mediatypes.Series]
	EpisodeCoreDataHandler() *handlers.CoreUserMediaItemDataHandler[*mediatypes.Episode]
	TrackCoreDataHandler() *handlers.CoreUserMediaItemDataHandler[*mediatypes.Track]
	AlbumCoreDataHandler() *handlers.CoreUserMediaItemDataHandler[*mediatypes.Album]
	ArtistCoreDataHandler() *handlers.CoreUserMediaItemDataHandler[*mediatypes.Artist]
	CollectionCoreDataHandler() *handlers.CoreUserMediaItemDataHandler[*mediatypes.Collection]
	PlaylistCoreDataHandler() *handlers.CoreUserMediaItemDataHandler[*mediatypes.Playlist]
}

// User-layer handlers (extend core)
type UserMediaItemHandlers interface {
	MovieUserHandler() *handlers.UserUserMediaItemDataHandler[*mediatypes.Movie]
	SeriesUserHandler() *handlers.UserUserMediaItemDataHandler[*mediatypes.Series]
	EpisodeUserHandler() *handlers.UserUserMediaItemDataHandler[*mediatypes.Episode]
	TrackUserHandler() *handlers.UserUserMediaItemDataHandler[*mediatypes.Track]
	AlbumUserHandler() *handlers.UserUserMediaItemDataHandler[*mediatypes.Album]
	ArtistUserHandler() *handlers.UserUserMediaItemDataHandler[*mediatypes.Artist]
	CollectionUserHandler() *handlers.UserUserMediaItemDataHandler[*mediatypes.Collection]
	PlaylistUserHandler() *handlers.UserUserMediaItemDataHandler[*mediatypes.Playlist]
}

// Client-layer handlers (extend user)
type ClientMediaItemHandlers interface {
	MovieClientHandler() *handlers.ClientUserMediaItemDataHandler[*mediatypes.Movie]
	SeriesClientHandler() *handlers.ClientUserMediaItemDataHandler[*mediatypes.Series]
	EpisodeClientHandler() *handlers.ClientUserMediaItemDataHandler[*mediatypes.Episode]
	TrackClientHandler() *handlers.ClientUserMediaItemDataHandler[*mediatypes.Track]
	AlbumClientHandler() *handlers.ClientUserMediaItemDataHandler[*mediatypes.Album]
	ArtistClientHandler() *handlers.ClientUserMediaItemDataHandler[*mediatypes.Artist]
	CollectionClientHandler() *handlers.ClientUserMediaItemDataHandler[*mediatypes.Collection]
	PlaylistClientHandler() *handlers.ClientUserMediaItemDataHandler[*mediatypes.Playlist]
}

// Specialized media handlers for specific domains
type SpecializedMediaHandlers interface {
	// Domain-specific handlers
	MusicHandler() *handlers.CoreMusicHandler
	SeriesSpecificHandler() *handlers.ClientMediaSeriesHandler[*types.JellyfinConfig] // Using JellyfinConfig as default

	// Season handler (special case)
	SeasonHandler() *handlers.CoreUserMediaItemDataHandler[*mediatypes.Season]
}

type ClientMediaHandlers interface {
	ClientMediaMovieHandlers
	ClientMediaSeriesHandlers
	ClientMediaEpisodeHandlers
	ClientMediaMusicHandlers
	ClientMediaPlaylistHandlers
}

type ClientMediaMovieHandlers interface {
	EmbyMovieHandler() *handlers.ClientMediaMovieHandler[*types.EmbyConfig]
	JellyfinMovieHandler() *handlers.ClientMediaMovieHandler[*types.JellyfinConfig]
	PlexMovieHandler() *handlers.ClientMediaMovieHandler[*types.PlexConfig]
}

type ClientMediaSeriesHandlers interface {
	EmbySeriesHandler() *handlers.ClientMediaSeriesHandler[*types.EmbyConfig]
	JellyfinSeriesHandler() *handlers.ClientMediaSeriesHandler[*types.JellyfinConfig]
	PlexSeriesHandler() *handlers.ClientMediaSeriesHandler[*types.PlexConfig]
}

type ClientMediaEpisodeHandlers interface {
	// EmbyEpisodeHandler() *handlers.ClientMediaEpisodeHandler[*types.EmbyConfig]
	// JellyfinEpisodeHandler() *handlers.ClientMediaEpisodeHandler[*types.JellyfinConfig]
	// PlexEpisodeHandler() *handlers.ClientMediaEpisodeHandler[*types.PlexConfig]
}

type ClientMediaMusicHandlers interface {
	EmbyMusicHandler() *handlers.ClientMediaMusicHandler[*types.EmbyConfig]
	JellyfinMusicHandler() *handlers.ClientMediaMusicHandler[*types.JellyfinConfig]
	PlexMusicHandler() *handlers.ClientMediaMusicHandler[*types.PlexConfig]
	SubsonicMusicHandler() *handlers.ClientMediaMusicHandler[*types.SubsonicConfig]
}

type ClientMediaPlaylistHandlers interface {
	EmbyPlaylistHandler() *handlers.ClientMediaPlaylistHandler[*types.EmbyConfig]
	JellyfinPlaylistHandler() *handlers.ClientMediaPlaylistHandler[*types.JellyfinConfig]
	PlexPlaylistHandler() *handlers.ClientMediaPlaylistHandler[*types.PlexConfig]
	SubsonicPlaylistHandler() *handlers.ClientMediaPlaylistHandler[*types.SubsonicConfig]
}

type UserHandlers interface {
	AuthHandler() *handlers.AuthHandler
	UserHandler() *handlers.UserHandler
	UserConfigHandler() *handlers.UserConfigHandler
}

type SystemHandlers interface {
	ConfigHandler() *handlers.ConfigHandler
	HealthHandler() *handlers.HealthHandler
	ClientsHandler() *handlers.ClientsHandler
}

type SystemServices interface {
	HealthService() services.HealthService
	ConfigService() services.ConfigService
}

type SystemRepositories interface {
	ConfigRepo() repository.ConfigRepository
}

type JobRepositories interface {
	JobRepo() repository.JobRepository
}

type JobServices interface {
	JobService() services.JobService
	RecommendationJob() *recommendation.RecommendationJob
	MediaSyncJob() *jobs.MediaSyncJob
	WatchHistorySyncJob() *jobs.WatchHistorySyncJob
	FavoritesSyncJob() *jobs.FavoritesSyncJob
}

type AIHandlers interface {
	ClaudeAIHandler() *handlers.AIHandler[*types.ClaudeConfig]
	OpenAIHandler() *handlers.AIHandler[*types.OpenAIConfig]
	OllamaHandler() *handlers.AIHandler[*types.OllamaConfig]
}

type JobHandlers interface {
	JobHandler() *handlers.JobHandler
	RecommendationHandler() *handlers.RecommendationHandler
}
