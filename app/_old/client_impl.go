// app/dependencies.go
package app

//
// import (
// 	mediatypes "suasor/client/media/types"
// 	"suasor/client/types"
// 	"suasor/handlers"
//
// 	"suasor/repository"
// 	"suasor/services"
// )
//
// // Concrete implementation of ClientServices
// type clientServicesImpl struct {
// 	embyService     services.ClientService[*types.EmbyConfig]
// 	jellyfinService services.ClientService[*types.JellyfinConfig]
// 	plexService     services.ClientService[*types.PlexConfig]
// 	subsonicService services.ClientService[*types.SubsonicConfig]
// 	sonarrService   services.ClientService[*types.SonarrConfig]
// 	radarrService   services.ClientService[*types.RadarrConfig]
// 	lidarrService   services.ClientService[*types.LidarrConfig]
// 	claudeService   services.ClientService[*types.ClaudeConfig]
// 	openaiService   services.ClientService[*types.OpenAIConfig]
// 	ollamaService   services.ClientService[*types.OllamaConfig]
// 	allServices     map[string]services.ClientService[types.ClientConfig]
// }
//
// func (s *clientServicesImpl) AllServices() map[string]services.ClientService[types.ClientConfig] {
// 	if s.allServices == nil {
// 		s.allServices = map[string]services.ClientService[types.ClientConfig]{
// 			"emby":     s.embyService.(services.ClientService[types.ClientConfig]),
// 			"jellyfin": s.jellyfinService.(services.ClientService[types.ClientConfig]),
// 			"plex":     s.plexService.(services.ClientService[types.ClientConfig]),
// 			"subsonic": s.subsonicService.(services.ClientService[types.ClientConfig]),
// 			"sonarr":   s.sonarrService.(services.ClientService[types.ClientConfig]),
// 			"radarr":   s.radarrService.(services.ClientService[types.ClientConfig]),
// 			"lidarr":   s.lidarrService.(services.ClientService[types.ClientConfig]),
// 			"claude":   s.claudeService.(services.ClientService[types.ClientConfig]),
// 			"openai":   s.openaiService.(services.ClientService[types.ClientConfig]),
// 			"ollama":   s.ollamaService.(services.ClientService[types.ClientConfig]),
// 		}
// 	}
//
// 	return s.allServices
// }
//
// func (s *clientServicesImpl) ClaudeService() services.ClientService[*types.ClaudeConfig] {
// 	return s.claudeService
// }
//
// func (s *clientServicesImpl) EmbyService() services.ClientService[*types.EmbyConfig] {
// 	return s.embyService
// }
//
// func (s *clientServicesImpl) JellyfinService() services.ClientService[*types.JellyfinConfig] {
// 	return s.jellyfinService
// }
//
// func (s *clientServicesImpl) PlexService() services.ClientService[*types.PlexConfig] {
// 	return s.plexService
// }
//
// func (s *clientServicesImpl) SubsonicService() services.ClientService[*types.SubsonicConfig] {
// 	return s.subsonicService
// }
//
// func (s *clientServicesImpl) SonarrService() services.ClientService[*types.SonarrConfig] {
// 	return s.sonarrService
// }
//
// func (s *clientServicesImpl) RadarrService() services.ClientService[*types.RadarrConfig] {
// 	return s.radarrService
// }
//
// func (s *clientServicesImpl) LidarrService() services.ClientService[*types.LidarrConfig] {
// 	return s.lidarrService
// }
//
// func (s *clientServicesImpl) OpenAIService() services.ClientService[*types.OpenAIConfig] {
// 	return s.openaiService
// }
//
// func (s *clientServicesImpl) OllamaService() services.ClientService[*types.OllamaConfig] {
// 	return s.ollamaService
// }
//
// // Concrete implementation of ClientSeriesServices
// type clientSeriesServicesImpl struct {
// 	embySeriesService     services.ClientMediaSeriesService[*types.EmbyConfig]
// 	jellyfinSeriesService services.ClientMediaSeriesService[*types.JellyfinConfig]
// 	plexSeriesService     services.ClientMediaSeriesService[*types.PlexConfig]
// 	subsonicSeriesService services.ClientMediaSeriesService[*types.SubsonicConfig]
// }
//
// func (s *clientSeriesServicesImpl) EmbySeriesService() services.ClientMediaSeriesService[*types.EmbyConfig] {
// 	return s.embySeriesService
// }
//
// func (s *clientSeriesServicesImpl) JellyfinSeriesService() services.ClientMediaSeriesService[*types.JellyfinConfig] {
// 	return s.jellyfinSeriesService
// }
//
// func (s *clientSeriesServicesImpl) PlexSeriesService() services.ClientMediaSeriesService[*types.PlexConfig] {
// 	return s.plexSeriesService
// }
//
// func (s *clientSeriesServicesImpl) SubsonicSeriesService() services.ClientMediaSeriesService[*types.SubsonicConfig] {
// 	return s.subsonicSeriesService
// }
//
// // Concrete implementation of MediaItemServices
// type mediaItemServicesImpl struct {
// 	movieService      services.CoreMediaItemService[*mediatypes.Movie]
// 	seriesService     services.CoreMediaItemService[*mediatypes.Series]
// 	episodeService    services.CoreMediaItemService[*mediatypes.Episode]
// 	trackService      services.CoreMediaItemService[*mediatypes.Track]
// 	albumService      services.CoreMediaItemService[*mediatypes.Album]
// 	artistService     services.CoreMediaItemService[*mediatypes.Artist]
// 	collectionService services.CoreMediaItemService[*mediatypes.Collection]
// 	playlistService   services.CoreMediaItemService[*mediatypes.Playlist]
//
// 	// Three-pronged architecture services for collections
// 	coreCollectionService     services.CoreCollectionService
// 	clientCollectionService   services.ClientMediaCollectionService
// 	collectionExtendedService services.UserCollectionService
//
// 	// Playlist services
// 	playlistExtendedService services.PlaylistService
// }
//
// func (s *mediaItemServicesImpl) CoreCollectionService() services.CoreCollectionService {
// 	return s.coreCollectionService
// }
//
// func (s *mediaItemServicesImpl) ClientCollectionService() services.ClientMediaCollectionService {
// 	return s.clientCollectionService
// }
//
// func (s *mediaItemServicesImpl) CollectionExtendedService() services.UserCollectionService {
// 	return s.collectionExtendedService
// }
//
// func (s *mediaItemServicesImpl) PlaylistExtendedService() services.PlaylistService {
// 	return s.playlistExtendedService
// }
//
// func (s *mediaItemServicesImpl) MovieService() services.CoreMediaItemService[*mediatypes.Movie] {
// 	return s.movieService
// }
//
// func (s *mediaItemServicesImpl) SeriesService() services.CoreMediaItemService[*mediatypes.Series] {
// 	return s.seriesService
// }
//
// func (s *mediaItemServicesImpl) EpisodeService() services.CoreMediaItemService[*mediatypes.Episode] {
// 	return s.episodeService
// }
//
// func (s *mediaItemServicesImpl) TrackService() services.CoreMediaItemService[*mediatypes.Track] {
// 	return s.trackService
// }
//
// func (s *mediaItemServicesImpl) AlbumService() services.CoreMediaItemService[*mediatypes.Album] {
// 	return s.albumService
// }
//
// func (s *mediaItemServicesImpl) ArtistService() services.CoreMediaItemService[*mediatypes.Artist] {
// 	return s.artistService
// }
//
// func (s *mediaItemServicesImpl) CollectionService() services.CoreMediaItemService[*mediatypes.Collection] {
// 	return s.collectionService
// }
//
// func (s *mediaItemServicesImpl) PlaylistService() services.CoreMediaItemService[*mediatypes.Playlist] {
// 	return s.playlistService
// }
//
// // Concrete implementation of ClientHandlers
// type clientHandlersImpl struct {
// 	embyHandler     *handlers.ClientHandler[*types.EmbyConfig]
// 	jellyfinHandler *handlers.ClientHandler[*types.JellyfinConfig]
// 	plexHandler     *handlers.ClientHandler[*types.PlexConfig]
// 	subsonicHandler *handlers.ClientHandler[*types.SubsonicConfig]
// 	radarrHandler   *handlers.ClientHandler[*types.RadarrConfig]
// 	lidarrHandler   *handlers.ClientHandler[*types.LidarrConfig]
// 	sonarrHandler   *handlers.ClientHandler[*types.SonarrConfig]
// 	claudeHandler   *handlers.ClientHandler[*types.ClaudeConfig]
// 	openaiHandler   *handlers.ClientHandler[*types.OpenAIConfig]
// 	ollamaHandler   *handlers.ClientHandler[*types.OllamaConfig]
// }
//
// func (h *clientHandlersImpl) OpenAIHandler() *handlers.ClientHandler[*types.OpenAIConfig] {
// 	return h.openaiHandler
// }
//
// func (h *clientHandlersImpl) OllamaHandler() *handlers.ClientHandler[*types.OllamaConfig] {
// 	return h.ollamaHandler
// }
//
// func (h *clientHandlersImpl) ClaudeHandler() *handlers.ClientHandler[*types.ClaudeConfig] {
// 	return h.claudeHandler
// }
//
// func (h *clientHandlersImpl) EmbyHandler() *handlers.ClientHandler[*types.EmbyConfig] {
// 	return h.embyHandler
// }
//
// func (h *clientHandlersImpl) JellyfinHandler() *handlers.ClientHandler[*types.JellyfinConfig] {
// 	return h.jellyfinHandler
// }
//
// func (h *clientHandlersImpl) PlexHandler() *handlers.ClientHandler[*types.PlexConfig] {
// 	return h.plexHandler
// }
//
// func (h *clientHandlersImpl) SubsonicHandler() *handlers.ClientHandler[*types.SubsonicConfig] {
// 	return h.subsonicHandler
// }
//
// func (h *clientHandlersImpl) RadarrHandler() *handlers.ClientHandler[*types.RadarrConfig] {
// 	return h.radarrHandler
// }
//
// func (h *clientHandlersImpl) LidarrHandler() *handlers.ClientHandler[*types.LidarrConfig] {
// 	return h.lidarrHandler
// }
//
// func (h *clientHandlersImpl) SonarrHandler() *handlers.ClientHandler[*types.SonarrConfig] {
// 	return h.sonarrHandler
// }
//
// type repositoryCollectionsImpl struct {
// 	clientRepos repository.ClientRepositoryCollection
// }
//
// func (r *repositoryCollectionsImpl) ClientRepositories() repository.ClientRepositoryCollection {
// 	return r.clientRepos
// }
//
// type clientRepositoriesImpl struct {
// 	embyRepo     repository.ClientRepository[*types.EmbyConfig]
// 	jellyfinRepo repository.ClientRepository[*types.JellyfinConfig]
// 	plexRepo     repository.ClientRepository[*types.PlexConfig]
// 	subsonicRepo repository.ClientRepository[*types.SubsonicConfig]
// 	sonarrRepo   repository.ClientRepository[*types.SonarrConfig]
// 	radarrRepo   repository.ClientRepository[*types.RadarrConfig]
// 	lidarrRepo   repository.ClientRepository[*types.LidarrConfig]
// 	claudeRepo   repository.ClientRepository[*types.ClaudeConfig]
// 	openaiRepo   repository.ClientRepository[*types.OpenAIConfig]
// 	ollamaRepo   repository.ClientRepository[*types.OllamaConfig]
// }
//
// // AllRepos returns all client repositories in a type-safe struct
// func (r *clientRepositoriesImpl) AllRepos() repository.ClientRepoCollection {
// 	return repository.ClientRepoCollection{
// 		EmbyRepo:     r.embyRepo,
// 		JellyfinRepo: r.jellyfinRepo,
// 		PlexRepo:     r.plexRepo,
// 		SubsonicRepo: r.subsonicRepo,
// 		SonarrRepo:   r.sonarrRepo,
// 		RadarrRepo:   r.radarrRepo,
// 		LidarrRepo:   r.lidarrRepo,
// 		ClaudeRepo:   r.claudeRepo,
// 		OpenAIRepo:   r.openaiRepo,
// 		OllamaRepo:   r.ollamaRepo,
// 	}
// }
//
// // GetAllByCategory returns repositories filtered by category
// func (r *clientRepositoriesImpl) GetAllByCategory(category types.ClientCategory) repository.ClientRepoCollection {
// 	allRepos := r.AllRepos()
// 	filteredRepos := repository.ClientRepoCollection{}
//
// 	// Only populate repositories that match the category
// 	if category == types.ClientCategoryMedia {
// 		filteredRepos.EmbyRepo = allRepos.EmbyRepo
// 		filteredRepos.JellyfinRepo = allRepos.JellyfinRepo
// 		filteredRepos.PlexRepo = allRepos.PlexRepo
// 		filteredRepos.SubsonicRepo = allRepos.SubsonicRepo
// 	} else if category == types.ClientCategoryAutomation {
// 		filteredRepos.SonarrRepo = allRepos.SonarrRepo
// 		filteredRepos.RadarrRepo = allRepos.RadarrRepo
// 		filteredRepos.LidarrRepo = allRepos.LidarrRepo
// 	} else if category == types.ClientCategoryAI {
// 		filteredRepos.ClaudeRepo = allRepos.ClaudeRepo
// 		filteredRepos.OpenAIRepo = allRepos.OpenAIRepo
// 		filteredRepos.OllamaRepo = allRepos.OllamaRepo
// 	}
//
// 	return filteredRepos
// }
//
// func (r *clientRepositoriesImpl) ClaudeRepo() repository.ClientRepository[*types.ClaudeConfig] {
// 	return r.claudeRepo
// }
//
// func (r *clientRepositoriesImpl) EmbyRepo() repository.ClientRepository[*types.EmbyConfig] {
// 	return r.embyRepo
// }
//
// func (r *clientRepositoriesImpl) JellyfinRepo() repository.ClientRepository[*types.JellyfinConfig] {
// 	return r.jellyfinRepo
// }
//
// func (r *clientRepositoriesImpl) PlexRepo() repository.ClientRepository[*types.PlexConfig] {
// 	return r.plexRepo
// }
//
// func (r *clientRepositoriesImpl) SubsonicRepo() repository.ClientRepository[*types.SubsonicConfig] {
// 	return r.subsonicRepo
// }
//
// func (r *clientRepositoriesImpl) SonarrRepo() repository.ClientRepository[*types.SonarrConfig] {
// 	return r.sonarrRepo
// }
//
// func (r *clientRepositoriesImpl) RadarrRepo() repository.ClientRepository[*types.RadarrConfig] {
// 	return r.radarrRepo
// }
//
// func (r *clientRepositoriesImpl) LidarrRepo() repository.ClientRepository[*types.LidarrConfig] {
// 	return r.lidarrRepo
// }
//
// func (r *clientRepositoriesImpl) OpenAIRepo() repository.ClientRepository[*types.OpenAIConfig] {
// 	return r.openaiRepo
// }
//
// func (r *clientRepositoriesImpl) OllamaRepo() repository.ClientRepository[*types.OllamaConfig] {
// 	return r.ollamaRepo
// }
//
// type clientMediaServicesImpl struct {
// 	movieServices    clientMovieServicesImpl
// 	seriesServices   clientSeriesServicesImpl
// 	episodeServices  clientEpisodeServicesImpl
// 	musicServices    clientMusicServicesImpl
// 	playlistServices clientPlaylistServicesImpl
// }
//
// func (cms *clientMediaServicesImpl) EmbyMovieService() services.ClientMediaMovieService[*types.EmbyConfig] {
// 	return cms.movieServices.EmbyMovieService()
// }
//
// func (cms *clientMediaServicesImpl) JellyfinMovieService() services.ClientMediaMovieService[*types.JellyfinConfig] {
// 	return cms.movieServices.JellyfinMovieService()
// }
//
// func (cms *clientMediaServicesImpl) PlexMovieService() services.ClientMediaMovieService[*types.PlexConfig] {
// 	return cms.movieServices.PlexMovieService()
// }
//
// func (cms *clientMediaServicesImpl) SubsonicMovieService() services.ClientMediaMovieService[*types.SubsonicConfig] {
// 	return cms.movieServices.SubsonicMovieService()
// }
//
// func (cms *clientMediaServicesImpl) EmbySeriesService() services.ClientMediaSeriesService[*types.EmbyConfig] {
// 	return cms.seriesServices.EmbySeriesService()
// }
//
// func (cms *clientMediaServicesImpl) JellyfinSeriesService() services.ClientMediaSeriesService[*types.JellyfinConfig] {
// 	return cms.seriesServices.JellyfinSeriesService()
// }
//
// func (cms *clientMediaServicesImpl) PlexSeriesService() services.ClientMediaSeriesService[*types.PlexConfig] {
// 	return cms.seriesServices.PlexSeriesService()
// }
//
// func (cms *clientMediaServicesImpl) SubsonicSeriesService() services.ClientMediaSeriesService[*types.SubsonicConfig] {
// 	return cms.seriesServices.SubsonicSeriesService()
// }
//
// func (cms *clientMediaServicesImpl) EmbyMusicService() services.ClientMediaMusicService[*types.EmbyConfig] {
// 	return cms.musicServices.EmbyMusicService()
// }
//
// func (cms *clientMediaServicesImpl) JellyfinMusicService() services.ClientMediaMusicService[*types.JellyfinConfig] {
// 	return cms.musicServices.JellyfinMusicService()
// }
//
// func (cms *clientMediaServicesImpl) PlexMusicService() services.ClientMediaMusicService[*types.PlexConfig] {
// 	return cms.musicServices.PlexMusicService()
// }
//
// func (cms *clientMediaServicesImpl) SubsonicMusicService() services.ClientMediaMusicService[*types.SubsonicConfig] {
// 	return cms.musicServices.SubsonicMusicService()
// }
//
// func (cms *clientMediaServicesImpl) EmbyPlaylistService() services.ClientMediaPlaylistService[*types.EmbyConfig] {
// 	return cms.playlistServices.EmbyPlaylistService()
// }
//
// func (cms *clientMediaServicesImpl) JellyfinPlaylistService() services.ClientMediaPlaylistService[*types.JellyfinConfig] {
// 	return cms.playlistServices.JellyfinPlaylistService()
// }
//
// func (cms *clientMediaServicesImpl) PlexPlaylistService() services.ClientMediaPlaylistService[*types.PlexConfig] {
// 	return cms.playlistServices.PlexPlaylistService()
// }
//
// func (cms *clientMediaServicesImpl) SubsonicPlaylistService() services.ClientMediaPlaylistService[*types.SubsonicConfig] {
// 	return cms.playlistServices.SubsonicPlaylistService()
// }
//
// type clientEpisodeServicesImpl struct {
// 	// embyEpisodeService     services.ClientMediaEpisodeService[*types.EmbyConfig]
// 	// jellyfinEpisodeService services.ClientMediaEpisodeService[*types.JellyfinConfig]
// 	// plexEpisodeService     services.ClientMediaEpisodeService[*types.PlexConfig]
// }
// type clientPlaylistServicesImpl struct {
// 	embyPlaylistService     services.ClientMediaPlaylistService[*types.EmbyConfig]
// 	jellyfinPlaylistService services.ClientMediaPlaylistService[*types.JellyfinConfig]
// 	plexPlaylistService     services.ClientMediaPlaylistService[*types.PlexConfig]
// 	subsonicPlaylistService services.ClientMediaPlaylistService[*types.SubsonicConfig]
// }
//
// func (s *clientPlaylistServicesImpl) EmbyPlaylistService() services.ClientMediaPlaylistService[*types.EmbyConfig] {
// 	return s.embyPlaylistService
// }
//
// func (s *clientPlaylistServicesImpl) JellyfinPlaylistService() services.ClientMediaPlaylistService[*types.JellyfinConfig] {
// 	return s.jellyfinPlaylistService
// }
//
// func (s *clientPlaylistServicesImpl) PlexPlaylistService() services.ClientMediaPlaylistService[*types.PlexConfig] {
// 	return s.plexPlaylistService
// }
//
// func (s *clientPlaylistServicesImpl) SubsonicPlaylistService() services.ClientMediaPlaylistService[*types.SubsonicConfig] {
// 	return s.subsonicPlaylistService
// }
//
// type clientMusicServicesImpl struct {
// 	embyMusicService     services.ClientMediaMusicService[*types.EmbyConfig]
// 	jellyfinMusicService services.ClientMediaMusicService[*types.JellyfinConfig]
// 	plexMusicService     services.ClientMediaMusicService[*types.PlexConfig]
// 	subsonicMusicService services.ClientMediaMusicService[*types.SubsonicConfig]
// }
//
// func (s *clientMusicServicesImpl) EmbyMusicService() services.ClientMediaMusicService[*types.EmbyConfig] {
// 	return s.embyMusicService
// }
//
// func (s *clientMusicServicesImpl) JellyfinMusicService() services.ClientMediaMusicService[*types.JellyfinConfig] {
// 	return s.jellyfinMusicService
// }
//
// func (s *clientMusicServicesImpl) PlexMusicService() services.ClientMediaMusicService[*types.PlexConfig] {
// 	return s.plexMusicService
// }
//
// func (s *clientMusicServicesImpl) SubsonicMusicService() services.ClientMediaMusicService[*types.SubsonicConfig] {
// 	return s.subsonicMusicService
// }
//
// type clientMovieServicesImpl struct {
// 	embyMovieService     services.ClientMediaMovieService[*types.EmbyConfig]
// 	jellyfinMovieService services.ClientMediaMovieService[*types.JellyfinConfig]
// 	plexMovieService     services.ClientMediaMovieService[*types.PlexConfig]
// 	subsonicMovieService services.ClientMediaMovieService[*types.SubsonicConfig]
// }
//
// func (s *clientMovieServicesImpl) EmbyMovieService() services.ClientMediaMovieService[*types.EmbyConfig] {
// 	return s.embyMovieService
// }
//
// func (s *clientMovieServicesImpl) JellyfinMovieService() services.ClientMediaMovieService[*types.JellyfinConfig] {
// 	return s.jellyfinMovieService
// }
//
// func (s *clientMovieServicesImpl) PlexMovieService() services.ClientMediaMovieService[*types.PlexConfig] {
// 	return s.plexMovieService
// }
//
// func (s *clientMovieServicesImpl) SubsonicMovieService() services.ClientMediaMovieService[*types.SubsonicConfig] {
// 	return s.subsonicMovieService
// }
//
// type clientMediaHandlersImpl struct {
// 	movieHandlers    *clientMediaMovieHandlersImpl
// 	seriesHandlers   *clientMediaSeriesHandlersImpl
// 	episodeHandlers  *clientMediaEpisodeHandlersImpl
// 	musicHandlers    *clientMediaMusicHandlersImpl
// 	playlistHandlers *clientMediaPlaylistHandlersImpl
// }
//
// func (h *clientMediaHandlersImpl) MovieHandlers() *clientMediaMovieHandlersImpl {
// 	return h.movieHandlers
// }
//
// func (h *clientMediaHandlersImpl) SeriesHandlers() *clientMediaSeriesHandlersImpl {
// 	return h.seriesHandlers
// }
//
// func (h *clientMediaHandlersImpl) EpisodeHandlers() *clientMediaEpisodeHandlersImpl {
// 	return h.episodeHandlers
// }
//
// func (h *clientMediaHandlersImpl) MusicHandlers() *clientMediaMusicHandlersImpl {
// 	return h.musicHandlers
// }
//
// func (h *clientMediaHandlersImpl) EmbyMovieHandler() *handlers.ClientMediaMovieHandler[*types.EmbyConfig] {
// 	return h.movieHandlers.embyMovieHandler
// }
//
// func (h *clientMediaHandlersImpl) JellyfinMovieHandler() *handlers.ClientMediaMovieHandler[*types.JellyfinConfig] {
// 	return h.movieHandlers.jellyfinMovieHandler
// }
//
// func (h *clientMediaHandlersImpl) PlexMovieHandler() *handlers.ClientMediaMovieHandler[*types.PlexConfig] {
// 	return h.movieHandlers.plexMovieHandler
// }
//
// func (h *clientMediaHandlersImpl) EmbySeriesHandler() *handlers.ClientMediaSeriesHandler[*types.EmbyConfig] {
// 	return h.seriesHandlers.embySeriesHandler
// }
//
// func (h *clientMediaHandlersImpl) JellyfinSeriesHandler() *handlers.ClientMediaSeriesHandler[*types.JellyfinConfig] {
// 	return h.seriesHandlers.jellyfinSeriesHandler
// }
//
// func (h *clientMediaHandlersImpl) PlexSeriesHandler() *handlers.ClientMediaSeriesHandler[*types.PlexConfig] {
// 	return h.seriesHandlers.plexSeriesHandler
// }
//
// func (h *clientMediaHandlersImpl) EmbyMusicHandler() *handlers.ClientMediaMusicHandler[*types.EmbyConfig] {
// 	return h.musicHandlers.embyMusicHandler
// }
//
// func (h *clientMediaHandlersImpl) JellyfinMusicHandler() *handlers.ClientMediaMusicHandler[*types.JellyfinConfig] {
// 	return h.musicHandlers.jellyfinMusicHandler
// }
//
// func (h *clientMediaHandlersImpl) PlexMusicHandler() *handlers.ClientMediaMusicHandler[*types.PlexConfig] {
// 	return h.musicHandlers.plexMusicHandler
// }
//
// func (h *clientMediaHandlersImpl) SubsonicMusicHandler() *handlers.ClientMediaMusicHandler[*types.SubsonicConfig] {
// 	return h.musicHandlers.subsonicMusicHandler
// }
//
// func (h *clientMediaHandlersImpl) PlaylistHandlers() *clientMediaPlaylistHandlersImpl {
// 	return h.playlistHandlers
// }
//
// func (h *clientMediaHandlersImpl) EmbyPlaylistHandler() *handlers.ClientMediaPlaylistHandler[*types.EmbyConfig] {
// 	return h.playlistHandlers.embyPlaylistHandler
// }
//
// func (h *clientMediaHandlersImpl) JellyfinPlaylistHandler() *handlers.ClientMediaPlaylistHandler[*types.JellyfinConfig] {
// 	return h.playlistHandlers.jellyfinPlaylistHandler
// }
//
// func (h *clientMediaHandlersImpl) PlexPlaylistHandler() *handlers.ClientMediaPlaylistHandler[*types.PlexConfig] {
// 	return h.playlistHandlers.plexPlaylistHandler
// }
//
// func (h *clientMediaHandlersImpl) SubsonicPlaylistHandler() *handlers.ClientMediaPlaylistHandler[*types.SubsonicConfig] {
// 	return h.playlistHandlers.subsonicPlaylistHandler
// }
//
// type clientMediaSeriesHandlersImpl struct {
// 	embySeriesHandler     *handlers.ClientMediaSeriesHandler[*types.EmbyConfig]
// 	jellyfinSeriesHandler *handlers.ClientMediaSeriesHandler[*types.JellyfinConfig]
// 	plexSeriesHandler     *handlers.ClientMediaSeriesHandler[*types.PlexConfig]
// }
//
// func (h *clientMediaSeriesHandlersImpl) EmbySeriesHandler() *handlers.ClientMediaSeriesHandler[*types.EmbyConfig] {
// 	return h.embySeriesHandler
// }
//
// func (h *clientMediaSeriesHandlersImpl) JellyfinSeriesHandler() *handlers.ClientMediaSeriesHandler[*types.JellyfinConfig] {
// 	return h.jellyfinSeriesHandler
// }
//
// func (h *clientMediaSeriesHandlersImpl) PlexSeriesHandler() *handlers.ClientMediaSeriesHandler[*types.PlexConfig] {
// 	return h.plexSeriesHandler
// }
//
// type clientMediaEpisodeHandlersImpl struct {
// 	// embyEpisodeHandler     *handlers.ClientMediaEpisodeHandler[*types.EmbyConfig]
// 	// jellyfinEpisodeHandler *handlers.ClientMediaEpisodeHandler[*types.JellyfinConfig]
// 	// plexEpisodeHandler     *handlers.ClientMediaEpisodeHandler[*types.PlexConfig]
// }
//
// type clientMediaMusicHandlersImpl struct {
// 	embyMusicHandler     *handlers.ClientMediaMusicHandler[*types.EmbyConfig]
// 	jellyfinMusicHandler *handlers.ClientMediaMusicHandler[*types.JellyfinConfig]
// 	plexMusicHandler     *handlers.ClientMediaMusicHandler[*types.PlexConfig]
// 	subsonicMusicHandler *handlers.ClientMediaMusicHandler[*types.SubsonicConfig]
// }
//
// func (h *clientMediaMusicHandlersImpl) EmbyMusicHandler() *handlers.ClientMediaMusicHandler[*types.EmbyConfig] {
// 	return h.embyMusicHandler
// }
//
// func (h *clientMediaMusicHandlersImpl) JellyfinMusicHandler() *handlers.ClientMediaMusicHandler[*types.JellyfinConfig] {
// 	return h.jellyfinMusicHandler
// }
//
// func (h *clientMediaMusicHandlersImpl) PlexMusicHandler() *handlers.ClientMediaMusicHandler[*types.PlexConfig] {
// 	return h.plexMusicHandler
// }
//
// func (h *clientMediaMusicHandlersImpl) SubsonicMusicHandler() *handlers.ClientMediaMusicHandler[*types.SubsonicConfig] {
// 	return h.subsonicMusicHandler
// }
//
// type clientMediaMovieHandlersImpl struct {
// 	embyMovieHandler     *handlers.ClientMediaMovieHandler[*types.EmbyConfig]
// 	jellyfinMovieHandler *handlers.ClientMediaMovieHandler[*types.JellyfinConfig]
// 	plexMovieHandler     *handlers.ClientMediaMovieHandler[*types.PlexConfig]
// }
//
// func (h *clientMediaMovieHandlersImpl) EmbyMovieHandler() *handlers.ClientMediaMovieHandler[*types.EmbyConfig] {
// 	return h.embyMovieHandler
// }
//
// func (h *clientMediaMovieHandlersImpl) JellyfinMovieHandler() *handlers.ClientMediaMovieHandler[*types.JellyfinConfig] {
// 	return h.jellyfinMovieHandler
// }
//
// func (h *clientMediaMovieHandlersImpl) PlexMovieHandler() *handlers.ClientMediaMovieHandler[*types.PlexConfig] {
// 	return h.plexMovieHandler
// }
//
// type clientMediaPlaylistHandlersImpl struct {
// 	embyPlaylistHandler     *handlers.ClientMediaPlaylistHandler[*types.EmbyConfig]
// 	jellyfinPlaylistHandler *handlers.ClientMediaPlaylistHandler[*types.JellyfinConfig]
// 	plexPlaylistHandler     *handlers.ClientMediaPlaylistHandler[*types.PlexConfig]
// 	subsonicPlaylistHandler *handlers.ClientMediaPlaylistHandler[*types.SubsonicConfig]
// }
//
// func (h *clientMediaPlaylistHandlersImpl) EmbyPlaylistHandler() *handlers.ClientMediaPlaylistHandler[*types.EmbyConfig] {
// 	return h.embyPlaylistHandler
// }
//
// func (h *clientMediaPlaylistHandlersImpl) JellyfinPlaylistHandler() *handlers.ClientMediaPlaylistHandler[*types.JellyfinConfig] {
// 	return h.jellyfinPlaylistHandler
// }
//
// func (h *clientMediaPlaylistHandlersImpl) PlexPlaylistHandler() *handlers.ClientMediaPlaylistHandler[*types.PlexConfig] {
// 	return h.plexPlaylistHandler
// }
//
// func (h *clientMediaPlaylistHandlersImpl) SubsonicPlaylistHandler() *handlers.ClientMediaPlaylistHandler[*types.SubsonicConfig] {
// 	return h.subsonicPlaylistHandler
// }
//
// // Core user media item data services
// type coreUserMediaItemDataServicesImpl struct {
// 	movieCoreService services.CoreUserMediaItemDataService[*mediatypes.Movie]
// 	seriesCoreService services.CoreUserMediaItemDataService[*mediatypes.Series]
// 	episodeCoreService services.CoreUserMediaItemDataService[*mediatypes.Episode]
// 	trackCoreService services.CoreUserMediaItemDataService[*mediatypes.Track]
// 	albumCoreService services.CoreUserMediaItemDataService[*mediatypes.Album]
// 	artistCoreService services.CoreUserMediaItemDataService[*mediatypes.Artist]
// 	collectionCoreService services.CoreUserMediaItemDataService[*mediatypes.Collection]
// 	playlistCoreService services.CoreUserMediaItemDataService[*mediatypes.Playlist]
// }
//
// func (s *coreUserMediaItemDataServicesImpl) MovieCoreService() services.CoreUserMediaItemDataService[*mediatypes.Movie] {
// 	return s.movieCoreService
// }
//
// func (s *coreUserMediaItemDataServicesImpl) SeriesCoreService() services.CoreUserMediaItemDataService[*mediatypes.Series] {
// 	return s.seriesCoreService
// }
//
// func (s *coreUserMediaItemDataServicesImpl) EpisodeCoreService() services.CoreUserMediaItemDataService[*mediatypes.Episode] {
// 	return s.episodeCoreService
// }
//
// func (s *coreUserMediaItemDataServicesImpl) TrackCoreService() services.CoreUserMediaItemDataService[*mediatypes.Track] {
// 	return s.trackCoreService
// }
//
// func (s *coreUserMediaItemDataServicesImpl) AlbumCoreService() services.CoreUserMediaItemDataService[*mediatypes.Album] {
// 	return s.albumCoreService
// }
//
// func (s *coreUserMediaItemDataServicesImpl) ArtistCoreService() services.CoreUserMediaItemDataService[*mediatypes.Artist] {
// 	return s.artistCoreService
// }
//
// func (s *coreUserMediaItemDataServicesImpl) CollectionCoreService() services.CoreUserMediaItemDataService[*mediatypes.Collection] {
// 	return s.collectionCoreService
// }
//
// func (s *coreUserMediaItemDataServicesImpl) PlaylistCoreService() services.CoreUserMediaItemDataService[*mediatypes.Playlist] {
// 	return s.playlistCoreService
// }
//
// type clientUserMediaItemDataServicesImpl struct {
// 	movieClientService      services.ClientUserMediaItemDataService[*mediatypes.Movie]
// 	seriesClientService     services.ClientUserMediaItemDataService[*mediatypes.Series]
// 	episodeClientService    services.ClientUserMediaItemDataService[*mediatypes.Episode]
// 	trackClientService      services.ClientUserMediaItemDataService[*mediatypes.Track]
// 	albumClientService      services.ClientUserMediaItemDataService[*mediatypes.Album]
// 	artistClientService     services.ClientUserMediaItemDataService[*mediatypes.Artist]
// 	collectionClientService services.ClientUserMediaItemDataService[*mediatypes.Collection]
// 	playlistClientService   services.ClientUserMediaItemDataService[*mediatypes.Playlist]
// }
//
// func (s *clientUserMediaItemDataServicesImpl) MovieDataService() services.ClientUserMediaItemDataService[*mediatypes.Movie] {
// 	return s.movieClientService
// }
//
// func (s *clientUserMediaItemDataServicesImpl) SeriesDataService() services.ClientUserMediaItemDataService[*mediatypes.Series] {
// 	return s.seriesClientService
// }
//
// func (s *clientUserMediaItemDataServicesImpl) EpisodeDataService() services.ClientUserMediaItemDataService[*mediatypes.Episode] {
// 	return s.episodeClientService
// }
//
// func (s *clientUserMediaItemDataServicesImpl) TrackDataService() services.ClientUserMediaItemDataService[*mediatypes.Track] {
// 	return s.trackClientService
// }
//
// func (s *clientUserMediaItemDataServicesImpl) AlbumDataService() services.ClientUserMediaItemDataService[*mediatypes.Album] {
// 	return s.albumClientService
// }
//
// func (s *clientUserMediaItemDataServicesImpl) ArtistDataService() services.ClientUserMediaItemDataService[*mediatypes.Artist] {
// 	return s.artistClientService
// }
//
// func (s *clientUserMediaItemDataServicesImpl) CollectionDataService() services.ClientUserMediaItemDataService[*mediatypes.Collection] {
// 	return s.collectionClientService
// }
//
// func (s *clientUserMediaItemDataServicesImpl) PlaylistDataService() services.ClientUserMediaItemDataService[*mediatypes.Playlist] {
// 	return s.playlistClientService
// }
