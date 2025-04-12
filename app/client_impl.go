// app/dependencies.go
package app

import (
	mediatypes "suasor/client/media/types"
	"suasor/client/types"
	"suasor/handlers"

	"suasor/repository"
	"suasor/services"
)

// Concrete implementation of ClientServices
type clientServicesImpl struct {
	embyService     services.ClientService[*types.EmbyConfig]
	jellyfinService services.ClientService[*types.JellyfinConfig]
	plexService     services.ClientService[*types.PlexConfig]
	subsonicService services.ClientService[*types.SubsonicConfig]
	sonarrService   services.ClientService[*types.SonarrConfig]
	radarrService   services.ClientService[*types.RadarrConfig]
	lidarrService   services.ClientService[*types.LidarrConfig]
	claudeService   services.ClientService[*types.ClaudeConfig]
	openaiService   services.ClientService[*types.OpenAIConfig]
	ollamaService   services.ClientService[*types.OllamaConfig]
	allServices     map[string]services.ClientService[types.ClientConfig]
}

func (s *clientServicesImpl) AllServices() map[string]services.ClientService[types.ClientConfig] {
	if s.allServices == nil {
		s.allServices = map[string]services.ClientService[types.ClientConfig]{
			"emby":     s.embyService.(services.ClientService[types.ClientConfig]),
			"jellyfin": s.jellyfinService.(services.ClientService[types.ClientConfig]),
			"plex":     s.plexService.(services.ClientService[types.ClientConfig]),
			"subsonic": s.subsonicService.(services.ClientService[types.ClientConfig]),
			"sonarr":   s.sonarrService.(services.ClientService[types.ClientConfig]),
			"radarr":   s.radarrService.(services.ClientService[types.ClientConfig]),
			"lidarr":   s.lidarrService.(services.ClientService[types.ClientConfig]),
			"claude":   s.claudeService.(services.ClientService[types.ClientConfig]),
			"openai":   s.openaiService.(services.ClientService[types.ClientConfig]),
			"ollama":   s.ollamaService.(services.ClientService[types.ClientConfig]),
		}
	}

	return s.allServices
}

func (s *clientServicesImpl) ClaudeService() services.ClientService[*types.ClaudeConfig] {
	return s.claudeService
}

func (s *clientServicesImpl) EmbyService() services.ClientService[*types.EmbyConfig] {
	return s.embyService
}

func (s *clientServicesImpl) JellyfinService() services.ClientService[*types.JellyfinConfig] {
	return s.jellyfinService
}

func (s *clientServicesImpl) PlexService() services.ClientService[*types.PlexConfig] {
	return s.plexService
}

func (s *clientServicesImpl) SubsonicService() services.ClientService[*types.SubsonicConfig] {
	return s.subsonicService
}

func (s *clientServicesImpl) SonarrService() services.ClientService[*types.SonarrConfig] {
	return s.sonarrService
}

func (s *clientServicesImpl) RadarrService() services.ClientService[*types.RadarrConfig] {
	return s.radarrService
}

func (s *clientServicesImpl) LidarrService() services.ClientService[*types.LidarrConfig] {
	return s.lidarrService
}

func (s *clientServicesImpl) OpenAIService() services.ClientService[*types.OpenAIConfig] {
	return s.openaiService
}

func (s *clientServicesImpl) OllamaService() services.ClientService[*types.OllamaConfig] {
	return s.ollamaService
}

// Concrete implementation of ClientSeriesServices
type clientSeriesServicesImpl struct {
	embySeriesService     services.MediaClientSeriesService[*types.EmbyConfig]
	jellyfinSeriesService services.MediaClientSeriesService[*types.JellyfinConfig]
	plexSeriesService     services.MediaClientSeriesService[*types.PlexConfig]
	subsonicSeriesService services.MediaClientSeriesService[*types.SubsonicConfig]
}

func (s *clientSeriesServicesImpl) EmbySeriesService() services.MediaClientSeriesService[*types.EmbyConfig] {
	return s.embySeriesService
}

func (s *clientSeriesServicesImpl) JellyfinSeriesService() services.MediaClientSeriesService[*types.JellyfinConfig] {
	return s.jellyfinSeriesService
}

func (s *clientSeriesServicesImpl) PlexSeriesService() services.MediaClientSeriesService[*types.PlexConfig] {
	return s.plexSeriesService
}

func (s *clientSeriesServicesImpl) SubsonicSeriesService() services.MediaClientSeriesService[*types.SubsonicConfig] {
	return s.subsonicSeriesService
}

// Concrete implementation of MediaItemServices
type mediaItemServicesImpl struct {
	movieService      services.MediaItemService[*mediatypes.Movie]
	seriesService     services.MediaItemService[*mediatypes.Series]
	episodeService    services.MediaItemService[*mediatypes.Episode]
	trackService      services.MediaItemService[*mediatypes.Track]
	albumService      services.MediaItemService[*mediatypes.Album]
	artistService     services.MediaItemService[*mediatypes.Artist]
	collectionService services.MediaItemService[*mediatypes.Collection]
	playlistService   services.MediaItemService[*mediatypes.Playlist]
}

func (s *mediaItemServicesImpl) MovieService() services.MediaItemService[*mediatypes.Movie] {
	return s.movieService
}

func (s *mediaItemServicesImpl) SeriesService() services.MediaItemService[*mediatypes.Series] {
	return s.seriesService
}

func (s *mediaItemServicesImpl) EpisodeService() services.MediaItemService[*mediatypes.Episode] {
	return s.episodeService
}

func (s *mediaItemServicesImpl) TrackService() services.MediaItemService[*mediatypes.Track] {
	return s.trackService
}

func (s *mediaItemServicesImpl) AlbumService() services.MediaItemService[*mediatypes.Album] {
	return s.albumService
}

func (s *mediaItemServicesImpl) ArtistService() services.MediaItemService[*mediatypes.Artist] {
	return s.artistService
}

func (s *mediaItemServicesImpl) CollectionService() services.MediaItemService[*mediatypes.Collection] {
	return s.collectionService
}

func (s *mediaItemServicesImpl) PlaylistService() services.MediaItemService[*mediatypes.Playlist] {
	return s.playlistService
}

// Concrete implementation of ClientHandlers
type clientHandlersImpl struct {
	embyHandler     *handlers.ClientHandler[*types.EmbyConfig]
	jellyfinHandler *handlers.ClientHandler[*types.JellyfinConfig]
	plexHandler     *handlers.ClientHandler[*types.PlexConfig]
	subsonicHandler *handlers.ClientHandler[*types.SubsonicConfig]
	radarrHandler   *handlers.ClientHandler[*types.RadarrConfig]
	lidarrHandler   *handlers.ClientHandler[*types.LidarrConfig]
	sonarrHandler   *handlers.ClientHandler[*types.SonarrConfig]
	claudeHandler   *handlers.ClientHandler[*types.ClaudeConfig]
	openaiHandler   *handlers.ClientHandler[*types.OpenAIConfig]
	ollamaHandler   *handlers.ClientHandler[*types.OllamaConfig]
}

func (h *clientHandlersImpl) OpenAIHandler() *handlers.ClientHandler[*types.OpenAIConfig] {
	return h.openaiHandler
}

func (h *clientHandlersImpl) OllamaHandler() *handlers.ClientHandler[*types.OllamaConfig] {
	return h.ollamaHandler
}

func (h *clientHandlersImpl) ClaudeHandler() *handlers.ClientHandler[*types.ClaudeConfig] {
	return h.claudeHandler
}

func (h *clientHandlersImpl) EmbyHandler() *handlers.ClientHandler[*types.EmbyConfig] {
	return h.embyHandler
}

func (h *clientHandlersImpl) JellyfinHandler() *handlers.ClientHandler[*types.JellyfinConfig] {
	return h.jellyfinHandler
}

func (h *clientHandlersImpl) PlexHandler() *handlers.ClientHandler[*types.PlexConfig] {
	return h.plexHandler
}

func (h *clientHandlersImpl) SubsonicHandler() *handlers.ClientHandler[*types.SubsonicConfig] {
	return h.subsonicHandler
}

func (h *clientHandlersImpl) RadarrHandler() *handlers.ClientHandler[*types.RadarrConfig] {
	return h.radarrHandler
}

func (h *clientHandlersImpl) LidarrHandler() *handlers.ClientHandler[*types.LidarrConfig] {
	return h.lidarrHandler
}

func (h *clientHandlersImpl) SonarrHandler() *handlers.ClientHandler[*types.SonarrConfig] {
	return h.sonarrHandler
}

type repositoryCollectionsImpl struct {
	clientRepos repository.ClientRepositoryCollection
}

func (r *repositoryCollectionsImpl) ClientRepositories() repository.ClientRepositoryCollection {
	return r.clientRepos
}

type clientRepositoriesImpl struct {
	embyRepo     repository.ClientRepository[*types.EmbyConfig]
	jellyfinRepo repository.ClientRepository[*types.JellyfinConfig]
	plexRepo     repository.ClientRepository[*types.PlexConfig]
	subsonicRepo repository.ClientRepository[*types.SubsonicConfig]
	sonarrRepo   repository.ClientRepository[*types.SonarrConfig]
	radarrRepo   repository.ClientRepository[*types.RadarrConfig]
	lidarrRepo   repository.ClientRepository[*types.LidarrConfig]
	claudeRepo   repository.ClientRepository[*types.ClaudeConfig]
	openaiRepo   repository.ClientRepository[*types.OpenAIConfig]
	ollamaRepo   repository.ClientRepository[*types.OllamaConfig]
}

// AllRepos returns all client repositories in a type-safe struct
func (r *clientRepositoriesImpl) AllRepos() repository.ClientRepoCollection {
	return repository.ClientRepoCollection{
		EmbyRepo:     r.embyRepo,
		JellyfinRepo: r.jellyfinRepo,
		PlexRepo:     r.plexRepo,
		SubsonicRepo: r.subsonicRepo,
		SonarrRepo:   r.sonarrRepo,
		RadarrRepo:   r.radarrRepo,
		LidarrRepo:   r.lidarrRepo,
		ClaudeRepo:   r.claudeRepo,
		OpenAIRepo:   r.openaiRepo,
		OllamaRepo:   r.ollamaRepo,
	}
}

// GetAllByCategory returns repositories filtered by category
func (r *clientRepositoriesImpl) GetAllByCategory(category types.ClientCategory) repository.ClientRepoCollection {
	allRepos := r.AllRepos()
	filteredRepos := repository.ClientRepoCollection{}

	// Only populate repositories that match the category
	if category == types.ClientCategoryMedia {
		filteredRepos.EmbyRepo = allRepos.EmbyRepo
		filteredRepos.JellyfinRepo = allRepos.JellyfinRepo
		filteredRepos.PlexRepo = allRepos.PlexRepo
		filteredRepos.SubsonicRepo = allRepos.SubsonicRepo
	} else if category == types.ClientCategoryAutomation {
		filteredRepos.SonarrRepo = allRepos.SonarrRepo
		filteredRepos.RadarrRepo = allRepos.RadarrRepo
		filteredRepos.LidarrRepo = allRepos.LidarrRepo
	} else if category == types.ClientCategoryAI {
		filteredRepos.ClaudeRepo = allRepos.ClaudeRepo
		filteredRepos.OpenAIRepo = allRepos.OpenAIRepo
		filteredRepos.OllamaRepo = allRepos.OllamaRepo
	}

	return filteredRepos
}

func (r *clientRepositoriesImpl) ClaudeRepo() repository.ClientRepository[*types.ClaudeConfig] {
	return r.claudeRepo
}

func (r *clientRepositoriesImpl) EmbyRepo() repository.ClientRepository[*types.EmbyConfig] {
	return r.embyRepo
}

func (r *clientRepositoriesImpl) JellyfinRepo() repository.ClientRepository[*types.JellyfinConfig] {
	return r.jellyfinRepo
}

func (r *clientRepositoriesImpl) PlexRepo() repository.ClientRepository[*types.PlexConfig] {
	return r.plexRepo
}

func (r *clientRepositoriesImpl) SubsonicRepo() repository.ClientRepository[*types.SubsonicConfig] {
	return r.subsonicRepo
}

func (r *clientRepositoriesImpl) SonarrRepo() repository.ClientRepository[*types.SonarrConfig] {
	return r.sonarrRepo
}

func (r *clientRepositoriesImpl) RadarrRepo() repository.ClientRepository[*types.RadarrConfig] {
	return r.radarrRepo
}

func (r *clientRepositoriesImpl) LidarrRepo() repository.ClientRepository[*types.LidarrConfig] {
	return r.lidarrRepo
}

func (r *clientRepositoriesImpl) OpenAIRepo() repository.ClientRepository[*types.OpenAIConfig] {
	return r.openaiRepo
}

func (r *clientRepositoriesImpl) OllamaRepo() repository.ClientRepository[*types.OllamaConfig] {
	return r.ollamaRepo
}

type clientMediaServicesImpl struct {
	movieServices    clientMovieServicesImpl
	seriesServices   clientSeriesServicesImpl
	episodeServices  clientEpisodeServicesImpl
	musicServices    clientMusicServicesImpl
	playlistServices clientPlaylistServicesImpl
}

func (cms *clientMediaServicesImpl) EmbyMovieService() services.MediaClientMovieService[*types.EmbyConfig] {
	return cms.movieServices.EmbyMovieService()
}

func (cms *clientMediaServicesImpl) JellyfinMovieService() services.MediaClientMovieService[*types.JellyfinConfig] {
	return cms.movieServices.JellyfinMovieService()
}

func (cms *clientMediaServicesImpl) PlexMovieService() services.MediaClientMovieService[*types.PlexConfig] {
	return cms.movieServices.PlexMovieService()
}

func (cms *clientMediaServicesImpl) SubsonicMovieService() services.MediaClientMovieService[*types.SubsonicConfig] {
	return cms.movieServices.SubsonicMovieService()
}

func (cms *clientMediaServicesImpl) EmbySeriesService() services.MediaClientSeriesService[*types.EmbyConfig] {
	return cms.seriesServices.EmbySeriesService()
}

func (cms *clientMediaServicesImpl) JellyfinSeriesService() services.MediaClientSeriesService[*types.JellyfinConfig] {
	return cms.seriesServices.JellyfinSeriesService()
}

func (cms *clientMediaServicesImpl) PlexSeriesService() services.MediaClientSeriesService[*types.PlexConfig] {
	return cms.seriesServices.PlexSeriesService()
}

func (cms *clientMediaServicesImpl) SubsonicSeriesService() services.MediaClientSeriesService[*types.SubsonicConfig] {
	return cms.seriesServices.SubsonicSeriesService()
}

func (cms *clientMediaServicesImpl) EmbyMusicService() services.MediaClientMusicService[*types.EmbyConfig] {
	return cms.musicServices.EmbyMusicService()
}

func (cms *clientMediaServicesImpl) JellyfinMusicService() services.MediaClientMusicService[*types.JellyfinConfig] {
	return cms.musicServices.JellyfinMusicService()
}

func (cms *clientMediaServicesImpl) PlexMusicService() services.MediaClientMusicService[*types.PlexConfig] {
	return cms.musicServices.PlexMusicService()
}

func (cms *clientMediaServicesImpl) SubsonicMusicService() services.MediaClientMusicService[*types.SubsonicConfig] {
	return cms.musicServices.SubsonicMusicService()
}

func (cms *clientMediaServicesImpl) EmbyPlaylistService() services.MediaClientPlaylistService[*types.EmbyConfig] {
	return cms.playlistServices.EmbyPlaylistService()
}

func (cms *clientMediaServicesImpl) JellyfinPlaylistService() services.MediaClientPlaylistService[*types.JellyfinConfig] {
	return cms.playlistServices.JellyfinPlaylistService()
}

func (cms *clientMediaServicesImpl) PlexPlaylistService() services.MediaClientPlaylistService[*types.PlexConfig] {
	return cms.playlistServices.PlexPlaylistService()
}

func (cms *clientMediaServicesImpl) SubsonicPlaylistService() services.MediaClientPlaylistService[*types.SubsonicConfig] {
	return cms.playlistServices.SubsonicPlaylistService()
}

type clientEpisodeServicesImpl struct {
	// embyEpisodeService     services.MediaClientEpisodeService[*types.EmbyConfig]
	// jellyfinEpisodeService services.MediaClientEpisodeService[*types.JellyfinConfig]
	// plexEpisodeService     services.MediaClientEpisodeService[*types.PlexConfig]
}
type clientPlaylistServicesImpl struct {
	embyPlaylistService     services.MediaClientPlaylistService[*types.EmbyConfig]
	jellyfinPlaylistService services.MediaClientPlaylistService[*types.JellyfinConfig]
	plexPlaylistService     services.MediaClientPlaylistService[*types.PlexConfig]
	subsonicPlaylistService services.MediaClientPlaylistService[*types.SubsonicConfig]
}

func (s *clientPlaylistServicesImpl) EmbyPlaylistService() services.MediaClientPlaylistService[*types.EmbyConfig] {
	return s.embyPlaylistService
}

func (s *clientPlaylistServicesImpl) JellyfinPlaylistService() services.MediaClientPlaylistService[*types.JellyfinConfig] {
	return s.jellyfinPlaylistService
}

func (s *clientPlaylistServicesImpl) PlexPlaylistService() services.MediaClientPlaylistService[*types.PlexConfig] {
	return s.plexPlaylistService
}

func (s *clientPlaylistServicesImpl) SubsonicPlaylistService() services.MediaClientPlaylistService[*types.SubsonicConfig] {
	return s.subsonicPlaylistService
}

type clientMusicServicesImpl struct {
	embyMusicService     services.MediaClientMusicService[*types.EmbyConfig]
	jellyfinMusicService services.MediaClientMusicService[*types.JellyfinConfig]
	plexMusicService     services.MediaClientMusicService[*types.PlexConfig]
	subsonicMusicService services.MediaClientMusicService[*types.SubsonicConfig]
}

func (s *clientMusicServicesImpl) EmbyMusicService() services.MediaClientMusicService[*types.EmbyConfig] {
	return s.embyMusicService
}

func (s *clientMusicServicesImpl) JellyfinMusicService() services.MediaClientMusicService[*types.JellyfinConfig] {
	return s.jellyfinMusicService
}

func (s *clientMusicServicesImpl) PlexMusicService() services.MediaClientMusicService[*types.PlexConfig] {
	return s.plexMusicService
}

func (s *clientMusicServicesImpl) SubsonicMusicService() services.MediaClientMusicService[*types.SubsonicConfig] {
	return s.subsonicMusicService
}

type clientMovieServicesImpl struct {
	embyMovieService     services.MediaClientMovieService[*types.EmbyConfig]
	jellyfinMovieService services.MediaClientMovieService[*types.JellyfinConfig]
	plexMovieService     services.MediaClientMovieService[*types.PlexConfig]
	subsonicMovieService services.MediaClientMovieService[*types.SubsonicConfig]
}

func (s *clientMovieServicesImpl) EmbyMovieService() services.MediaClientMovieService[*types.EmbyConfig] {
	return s.embyMovieService
}

func (s *clientMovieServicesImpl) JellyfinMovieService() services.MediaClientMovieService[*types.JellyfinConfig] {
	return s.jellyfinMovieService
}

func (s *clientMovieServicesImpl) PlexMovieService() services.MediaClientMovieService[*types.PlexConfig] {
	return s.plexMovieService
}

func (s *clientMovieServicesImpl) SubsonicMovieService() services.MediaClientMovieService[*types.SubsonicConfig] {
	return s.subsonicMovieService
}

type clientMediaHandlersImpl struct {
	movieHandlers    *clientMediaMovieHandlersImpl
	seriesHandlers   *clientMediaSeriesHandlersImpl
	episodeHandlers  *clientMediaEpisodeHandlersImpl
	musicHandlers    *clientMediaMusicHandlersImpl
	playlistHandlers *clientMediaPlaylistHandlersImpl
}

func (h *clientMediaHandlersImpl) MovieHandlers() *clientMediaMovieHandlersImpl {
	return h.movieHandlers
}

func (h *clientMediaHandlersImpl) SeriesHandlers() *clientMediaSeriesHandlersImpl {
	return h.seriesHandlers
}

func (h *clientMediaHandlersImpl) EpisodeHandlers() *clientMediaEpisodeHandlersImpl {
	return h.episodeHandlers
}

func (h *clientMediaHandlersImpl) MusicHandlers() *clientMediaMusicHandlersImpl {
	return h.musicHandlers
}

func (h *clientMediaHandlersImpl) EmbyMovieHandler() *handlers.MediaClientMovieHandler[*types.EmbyConfig] {
	return h.movieHandlers.embyMovieHandler
}

func (h *clientMediaHandlersImpl) JellyfinMovieHandler() *handlers.MediaClientMovieHandler[*types.JellyfinConfig] {
	return h.movieHandlers.jellyfinMovieHandler
}

func (h *clientMediaHandlersImpl) PlexMovieHandler() *handlers.MediaClientMovieHandler[*types.PlexConfig] {
	return h.movieHandlers.plexMovieHandler
}

func (h *clientMediaHandlersImpl) EmbySeriesHandler() *handlers.MediaClientSeriesHandler[*types.EmbyConfig] {
	return h.seriesHandlers.embySeriesHandler
}

func (h *clientMediaHandlersImpl) JellyfinSeriesHandler() *handlers.MediaClientSeriesHandler[*types.JellyfinConfig] {
	return h.seriesHandlers.jellyfinSeriesHandler
}

func (h *clientMediaHandlersImpl) PlexSeriesHandler() *handlers.MediaClientSeriesHandler[*types.PlexConfig] {
	return h.seriesHandlers.plexSeriesHandler
}

func (h *clientMediaHandlersImpl) EmbyMusicHandler() *handlers.MediaClientMusicHandler[*types.EmbyConfig] {
	return h.musicHandlers.embyMusicHandler
}

func (h *clientMediaHandlersImpl) JellyfinMusicHandler() *handlers.MediaClientMusicHandler[*types.JellyfinConfig] {
	return h.musicHandlers.jellyfinMusicHandler
}

func (h *clientMediaHandlersImpl) PlexMusicHandler() *handlers.MediaClientMusicHandler[*types.PlexConfig] {
	return h.musicHandlers.plexMusicHandler
}

func (h *clientMediaHandlersImpl) SubsonicMusicHandler() *handlers.MediaClientMusicHandler[*types.SubsonicConfig] {
	return h.musicHandlers.subsonicMusicHandler
}

func (h *clientMediaHandlersImpl) PlaylistHandlers() *clientMediaPlaylistHandlersImpl {
	return h.playlistHandlers
}

func (h *clientMediaHandlersImpl) EmbyPlaylistHandler() *handlers.MediaClientPlaylistHandler[*types.EmbyConfig] {
	return h.playlistHandlers.embyPlaylistHandler
}

func (h *clientMediaHandlersImpl) JellyfinPlaylistHandler() *handlers.MediaClientPlaylistHandler[*types.JellyfinConfig] {
	return h.playlistHandlers.jellyfinPlaylistHandler
}

func (h *clientMediaHandlersImpl) PlexPlaylistHandler() *handlers.MediaClientPlaylistHandler[*types.PlexConfig] {
	return h.playlistHandlers.plexPlaylistHandler
}

func (h *clientMediaHandlersImpl) SubsonicPlaylistHandler() *handlers.MediaClientPlaylistHandler[*types.SubsonicConfig] {
	return h.playlistHandlers.subsonicPlaylistHandler
}

type clientMediaSeriesHandlersImpl struct {
	embySeriesHandler     *handlers.MediaClientSeriesHandler[*types.EmbyConfig]
	jellyfinSeriesHandler *handlers.MediaClientSeriesHandler[*types.JellyfinConfig]
	plexSeriesHandler     *handlers.MediaClientSeriesHandler[*types.PlexConfig]
}

func (h *clientMediaSeriesHandlersImpl) EmbySeriesHandler() *handlers.MediaClientSeriesHandler[*types.EmbyConfig] {
	return h.embySeriesHandler
}

func (h *clientMediaSeriesHandlersImpl) JellyfinSeriesHandler() *handlers.MediaClientSeriesHandler[*types.JellyfinConfig] {
	return h.jellyfinSeriesHandler
}

func (h *clientMediaSeriesHandlersImpl) PlexSeriesHandler() *handlers.MediaClientSeriesHandler[*types.PlexConfig] {
	return h.plexSeriesHandler
}

type clientMediaEpisodeHandlersImpl struct {
	// embyEpisodeHandler     *handlers.MediaClientEpisodeHandler[*types.EmbyConfig]
	// jellyfinEpisodeHandler *handlers.MediaClientEpisodeHandler[*types.JellyfinConfig]
	// plexEpisodeHandler     *handlers.MediaClientEpisodeHandler[*types.PlexConfig]
}

type clientMediaMusicHandlersImpl struct {
	embyMusicHandler     *handlers.MediaClientMusicHandler[*types.EmbyConfig]
	jellyfinMusicHandler *handlers.MediaClientMusicHandler[*types.JellyfinConfig]
	plexMusicHandler     *handlers.MediaClientMusicHandler[*types.PlexConfig]
	subsonicMusicHandler *handlers.MediaClientMusicHandler[*types.SubsonicConfig]
}

func (h *clientMediaMusicHandlersImpl) EmbyMusicHandler() *handlers.MediaClientMusicHandler[*types.EmbyConfig] {
	return h.embyMusicHandler
}

func (h *clientMediaMusicHandlersImpl) JellyfinMusicHandler() *handlers.MediaClientMusicHandler[*types.JellyfinConfig] {
	return h.jellyfinMusicHandler
}

func (h *clientMediaMusicHandlersImpl) PlexMusicHandler() *handlers.MediaClientMusicHandler[*types.PlexConfig] {
	return h.plexMusicHandler
}

func (h *clientMediaMusicHandlersImpl) SubsonicMusicHandler() *handlers.MediaClientMusicHandler[*types.SubsonicConfig] {
	return h.subsonicMusicHandler
}

type clientMediaMovieHandlersImpl struct {
	embyMovieHandler     *handlers.MediaClientMovieHandler[*types.EmbyConfig]
	jellyfinMovieHandler *handlers.MediaClientMovieHandler[*types.JellyfinConfig]
	plexMovieHandler     *handlers.MediaClientMovieHandler[*types.PlexConfig]
}

func (h *clientMediaMovieHandlersImpl) EmbyMovieHandler() *handlers.MediaClientMovieHandler[*types.EmbyConfig] {
	return h.embyMovieHandler
}

func (h *clientMediaMovieHandlersImpl) JellyfinMovieHandler() *handlers.MediaClientMovieHandler[*types.JellyfinConfig] {
	return h.jellyfinMovieHandler
}

func (h *clientMediaMovieHandlersImpl) PlexMovieHandler() *handlers.MediaClientMovieHandler[*types.PlexConfig] {
	return h.plexMovieHandler
}

type clientMediaPlaylistHandlersImpl struct {
	embyPlaylistHandler     *handlers.MediaClientPlaylistHandler[*types.EmbyConfig]
	jellyfinPlaylistHandler *handlers.MediaClientPlaylistHandler[*types.JellyfinConfig]
	plexPlaylistHandler     *handlers.MediaClientPlaylistHandler[*types.PlexConfig]
	subsonicPlaylistHandler *handlers.MediaClientPlaylistHandler[*types.SubsonicConfig]
}

func (h *clientMediaPlaylistHandlersImpl) EmbyPlaylistHandler() *handlers.MediaClientPlaylistHandler[*types.EmbyConfig] {
	return h.embyPlaylistHandler
}

func (h *clientMediaPlaylistHandlersImpl) JellyfinPlaylistHandler() *handlers.MediaClientPlaylistHandler[*types.JellyfinConfig] {
	return h.jellyfinPlaylistHandler
}

func (h *clientMediaPlaylistHandlersImpl) PlexPlaylistHandler() *handlers.MediaClientPlaylistHandler[*types.PlexConfig] {
	return h.plexPlaylistHandler
}

func (h *clientMediaPlaylistHandlersImpl) SubsonicPlaylistHandler() *handlers.MediaClientPlaylistHandler[*types.SubsonicConfig] {
	return h.subsonicPlaylistHandler
}
