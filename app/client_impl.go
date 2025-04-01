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

// Concrete implementation of ClientSeriesServices (placeholder)
type clientSeriesServicesImpl struct{}

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

type clientRepositoriesImpl struct {
	embyRepo     repository.ClientRepository[*types.EmbyConfig]
	jellyfinRepo repository.ClientRepository[*types.JellyfinConfig]
	plexRepo     repository.ClientRepository[*types.PlexConfig]
	subsonicRepo repository.ClientRepository[*types.SubsonicConfig]
	sonarrRepo   repository.ClientRepository[*types.SonarrConfig]
	radarrRepo   repository.ClientRepository[*types.RadarrConfig]
	lidarrRepo   repository.ClientRepository[*types.LidarrConfig]
	claudeRepo   repository.ClientRepository[*types.ClaudeConfig]
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

type clientMediaServicesImpl struct {
	movieServices    clientMovieServicesImpl
	seriesServices   clientSeriesServicesImpl
	episodeServices  clientEpisodeServicesImpl
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

type clientEpisodeServicesImpl struct {
	// embyEpisodeService     services.MediaClientEpisodeService[*types.EmbyConfig]
	// jellyfinEpisodeService services.MediaClientEpisodeService[*types.JellyfinConfig]
	// plexEpisodeService     services.MediaClientEpisodeService[*types.PlexConfig]
}
type clientPlaylistServicesImpl struct {
	// embyPlaylistService     services.MediaClientPlaylistService[*types.EmbyConfig]
	// jellyfinPlaylistService services.MediaClientPlaylistService[*types.JellyfinConfig]
	// plexPlaylistService     services.MediaClientPlaylistService[*types.PlexConfig]
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
	movieHandlers   *clientMediaMovieHandlersImpl
	seriesHandlers  *clientMediaSeriesHandlersImpl
	episodeHandlers *clientMediaEpisodeHandlersImpl
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

func (h *clientMediaHandlersImpl) EmbyMovieHandler() *handlers.MediaClientMovieHandler[*types.EmbyConfig] {
	return h.movieHandlers.embyMovieHandler
}

func (h *clientMediaHandlersImpl) JellyfinMovieHandler() *handlers.MediaClientMovieHandler[*types.JellyfinConfig] {
	return h.movieHandlers.jellyfinMovieHandler
}

func (h *clientMediaHandlersImpl) PlexMovieHandler() *handlers.MediaClientMovieHandler[*types.PlexConfig] {
	return h.movieHandlers.plexMovieHandler
}

type clientMediaSeriesHandlersImpl struct {
	// embySeriesHandler     *handlers.MediaClientSeriesHandler[*types.EmbyConfig]
	// jellyfinSeriesHandler *handlers.MediaClientSeriesHandler[*types.JellyfinConfig]
	// plexSeriesHandler     *handlers.MediaClientSeriesHandler[*types.PlexConfig]
}

type clientMediaEpisodeHandlersImpl struct {
	// embyEpisodeHandler     *handlers.MediaClientEpisodeHandler[*types.EmbyConfig]
	// jellyfinEpisodeHandler *handlers.MediaClientEpisodeHandler[*types.JellyfinConfig]
	// plexEpisodeHandler     *handlers.MediaClientEpisodeHandler[*types.PlexConfig]
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
