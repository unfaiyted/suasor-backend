// app/dependencies.go
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

// Using repository.ClientRepoCollection

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
	EmbyMovieService() services.MediaClientMovieService[*types.EmbyConfig]
	JellyfinMovieService() services.MediaClientMovieService[*types.JellyfinConfig]
	PlexMovieService() services.MediaClientMovieService[*types.PlexConfig]
	SubsonicMovieService() services.MediaClientMovieService[*types.SubsonicConfig]
}

type ClientSeriesServices interface {
	EmbySeriesService() services.MediaClientSeriesService[*types.EmbyConfig]
	JellyfinSeriesService() services.MediaClientSeriesService[*types.JellyfinConfig]
	PlexSeriesService() services.MediaClientSeriesService[*types.PlexConfig]
	SubsonicSeriesService() services.MediaClientSeriesService[*types.SubsonicConfig]
}

type ClientEpisodeServices interface {
	// TODO: implement
}

type ClientMusicServices interface {
	EmbyMusicService() services.MediaClientMusicService[*types.EmbyConfig]
	JellyfinMusicService() services.MediaClientMusicService[*types.JellyfinConfig]
	PlexMusicService() services.MediaClientMusicService[*types.PlexConfig]
	SubsonicMusicService() services.MediaClientMusicService[*types.SubsonicConfig]
}

type ClientPlaylistServices interface {
	EmbyPlaylistService() services.MediaClientPlaylistService[*types.EmbyConfig]
	JellyfinPlaylistService() services.MediaClientPlaylistService[*types.JellyfinConfig]
	PlexPlaylistService() services.MediaClientPlaylistService[*types.PlexConfig]
	SubsonicPlaylistService() services.MediaClientPlaylistService[*types.SubsonicConfig]
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

	// Extended services (legacy interfaces)
	CollectionExtendedService() services.UserCollectionService
	PlaylistExtendedService() services.PlaylistService
	
	// Three-pronged architecture for collections
	CoreCollectionService() services.CoreCollectionService
	ClientCollectionService() services.ClientCollectionService
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

type MediaItemHandlers interface {
	MovieHandler() *handlers.MediaItemHandler[*mediatypes.Movie]
	SeriesHandler() *handlers.MediaItemHandler[*mediatypes.Series]
	EpisodeHandler() *handlers.MediaItemHandler[*mediatypes.Episode]
	SeasonHandler() *handlers.MediaItemHandler[*mediatypes.Season]
	TrackHandler() *handlers.MediaItemHandler[*mediatypes.Track]
	AlbumHandler() *handlers.MediaItemHandler[*mediatypes.Album]
	ArtistHandler() *handlers.MediaItemHandler[*mediatypes.Artist]
	CollectionHandler() *handlers.MediaItemHandler[*mediatypes.Collection]
	PlaylistHandler() *handlers.MediaItemHandler[*mediatypes.Playlist]

	// Specialized handlers
	MusicHandler() *handlers.MusicSpecificHandler
	SeriesSpecificHandler() *handlers.SeriesSpecificHandler
	PlaylistSpecificHandler() *handlers.PlaylistHandler
	CollectionSpecificHandler() *handlers.CollectionHandler
}

type ClientMediaHandlers interface {
	ClientMediaMovieHandlers
	ClientMediaSeriesHandlers
	ClientMediaEpisodeHandlers
	ClientMediaMusicHandlers
	ClientMediaPlaylistHandlers
}

type ClientMediaMovieHandlers interface {
	EmbyMovieHandler() *handlers.MediaClientMovieHandler[*types.EmbyConfig]
	JellyfinMovieHandler() *handlers.MediaClientMovieHandler[*types.JellyfinConfig]
	PlexMovieHandler() *handlers.MediaClientMovieHandler[*types.PlexConfig]
}

type ClientMediaSeriesHandlers interface {
	EmbySeriesHandler() *handlers.MediaClientSeriesHandler[*types.EmbyConfig]
	JellyfinSeriesHandler() *handlers.MediaClientSeriesHandler[*types.JellyfinConfig]
	PlexSeriesHandler() *handlers.MediaClientSeriesHandler[*types.PlexConfig]
}

type ClientMediaEpisodeHandlers interface {
	// EmbyEpisodeHandler() *handlers.MediaClientEpisodeHandler[*types.EmbyConfig]
	// JellyfinEpisodeHandler() *handlers.MediaClientEpisodeHandler[*types.JellyfinConfig]
	// PlexEpisodeHandler() *handlers.MediaClientEpisodeHandler[*types.PlexConfig]
}

type ClientMediaMusicHandlers interface {
	EmbyMusicHandler() *handlers.MediaClientMusicHandler[*types.EmbyConfig]
	JellyfinMusicHandler() *handlers.MediaClientMusicHandler[*types.JellyfinConfig]
	PlexMusicHandler() *handlers.MediaClientMusicHandler[*types.PlexConfig]
	SubsonicMusicHandler() *handlers.MediaClientMusicHandler[*types.SubsonicConfig]
}

type ClientMediaPlaylistHandlers interface {
	EmbyPlaylistHandler() *handlers.MediaClientPlaylistHandler[*types.EmbyConfig]
	JellyfinPlaylistHandler() *handlers.MediaClientPlaylistHandler[*types.JellyfinConfig]
	PlexPlaylistHandler() *handlers.MediaClientPlaylistHandler[*types.PlexConfig]
	SubsonicPlaylistHandler() *handlers.MediaClientPlaylistHandler[*types.SubsonicConfig]
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
