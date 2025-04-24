package services

import (
	"context"
	"gorm.io/gorm"
	mediatypes "suasor/clients/media/types"
	clienttypes "suasor/clients/types"

	"suasor/di/container"
	"suasor/repository"
	"suasor/services"
	svcbundles "suasor/services/bundles"
	"suasor/utils/logger"
)

func registerMediaItemServices(ctx context.Context, c *container.Container) {
	log := logger.LoggerFromContext(ctx)

	log.Info().Msg("Registering media item services")
	registerMediaItemService[*mediatypes.Movie](c)
	registerMediaItemService[*mediatypes.Series](c)
	registerMediaItemService[*mediatypes.Season](c)
	registerMediaItemService[*mediatypes.Episode](c)
	registerMediaItemService[*mediatypes.Track](c)
	registerMediaItemService[*mediatypes.Album](c)
	registerMediaItemService[*mediatypes.Artist](c)
	registerMediaItemService[*mediatypes.Collection](c)
	registerMediaItemService[*mediatypes.Playlist](c)

	log.Info().Msg("Registering client media item services")
	registerClientMediaItemService[*clienttypes.EmbyConfig, *mediatypes.Movie](c)
	registerClientMediaItemService[*clienttypes.JellyfinConfig, *mediatypes.Movie](c)
	registerClientMediaItemService[*clienttypes.PlexConfig, *mediatypes.Movie](c)

	registerClientMediaItemService[*clienttypes.EmbyConfig, *mediatypes.Series](c)
	registerClientMediaItemService[*clienttypes.JellyfinConfig, *mediatypes.Series](c)
	registerClientMediaItemService[*clienttypes.PlexConfig, *mediatypes.Series](c)

	registerClientMediaItemService[*clienttypes.EmbyConfig, *mediatypes.Season](c)
	registerClientMediaItemService[*clienttypes.JellyfinConfig, *mediatypes.Season](c)
	registerClientMediaItemService[*clienttypes.PlexConfig, *mediatypes.Season](c)

	registerClientMediaItemService[*clienttypes.EmbyConfig, *mediatypes.Episode](c)
	registerClientMediaItemService[*clienttypes.JellyfinConfig, *mediatypes.Episode](c)
	registerClientMediaItemService[*clienttypes.PlexConfig, *mediatypes.Episode](c)

	registerClientMediaItemService[*clienttypes.EmbyConfig, *mediatypes.Track](c)
	registerClientMediaItemService[*clienttypes.JellyfinConfig, *mediatypes.Track](c)
	registerClientMediaItemService[*clienttypes.PlexConfig, *mediatypes.Track](c)
	registerClientMediaItemService[*clienttypes.SubsonicConfig, *mediatypes.Track](c)

	registerClientMediaItemService[*clienttypes.EmbyConfig, *mediatypes.Album](c)
	registerClientMediaItemService[*clienttypes.JellyfinConfig, *mediatypes.Album](c)
	registerClientMediaItemService[*clienttypes.PlexConfig, *mediatypes.Album](c)
	registerClientMediaItemService[*clienttypes.SubsonicConfig, *mediatypes.Album](c)

	registerClientMediaItemService[*clienttypes.EmbyConfig, *mediatypes.Artist](c)
	registerClientMediaItemService[*clienttypes.JellyfinConfig, *mediatypes.Artist](c)
	registerClientMediaItemService[*clienttypes.PlexConfig, *mediatypes.Artist](c)
	registerClientMediaItemService[*clienttypes.SubsonicConfig, *mediatypes.Artist](c)

	registerClientMediaItemService[*clienttypes.EmbyConfig, *mediatypes.Collection](c)
	registerClientMediaItemService[*clienttypes.JellyfinConfig, *mediatypes.Collection](c)
	registerClientMediaItemService[*clienttypes.PlexConfig, *mediatypes.Collection](c)

	registerClientMediaItemService[*clienttypes.EmbyConfig, *mediatypes.Playlist](c)
	registerClientMediaItemService[*clienttypes.JellyfinConfig, *mediatypes.Playlist](c)
	registerClientMediaItemService[*clienttypes.PlexConfig, *mediatypes.Playlist](c)
	registerClientMediaItemService[*clienttypes.SubsonicConfig, *mediatypes.Playlist](c)

	// People and credits services
	container.RegisterFactory[*services.PersonService](c, func(c *container.Container) *services.PersonService {
		personRepo := container.MustGet[repository.PersonRepository](c)
		creditRepo := container.MustGet[repository.CreditRepository](c)
		return services.NewPersonService(personRepo, creditRepo)
	})

	container.RegisterFactory[*services.CreditService](c, func(c *container.Container) *services.CreditService {
		creditRepo := container.MustGet[repository.CreditRepository](c)
		personRepo := container.MustGet[repository.PersonRepository](c)
		return services.NewCreditService(creditRepo, personRepo)
	})

	// Media services
	container.RegisterFactory[svcbundles.PeopleServices](c, func(c *container.Container) svcbundles.PeopleServices {
		personService := container.MustGet[services.PersonService](c)
		creditService := container.MustGet[services.CreditService](c)
		return svcbundles.NewPeopleServices(&personService, &creditService)
	})

	// Register CoreMediaItemServices bundle
	container.RegisterFactory[svcbundles.CoreMediaItemServices](c, func(c *container.Container) svcbundles.CoreMediaItemServices {
		return &coreMediaItemServicesImpl{
			movieCoreService:      container.MustGet[services.CoreMediaItemService[*mediatypes.Movie]](c),
			seriesCoreService:     container.MustGet[services.CoreMediaItemService[*mediatypes.Series]](c),
			seasonCoreService:     container.MustGet[services.CoreMediaItemService[*mediatypes.Season]](c),
			episodeCoreService:    container.MustGet[services.CoreMediaItemService[*mediatypes.Episode]](c),
			trackCoreService:      container.MustGet[services.CoreMediaItemService[*mediatypes.Track]](c),
			albumCoreService:      container.MustGet[services.CoreMediaItemService[*mediatypes.Album]](c),
			artistCoreService:     container.MustGet[services.CoreMediaItemService[*mediatypes.Artist]](c),
			collectionCoreService: container.MustGet[services.CoreMediaItemService[*mediatypes.Collection]](c),
			playlistCoreService:   container.MustGet[services.CoreMediaItemService[*mediatypes.Playlist]](c),
		}
	})
	
	// Register MusicRepository
	container.RegisterFactory[repository.MusicRepository](c, func(c *container.Container) repository.MusicRepository {
		db := container.MustGet[*gorm.DB](c)
		trackRepo := container.MustGet[repository.MediaItemRepository[*mediatypes.Track]](c)
		albumRepo := container.MustGet[repository.MediaItemRepository[*mediatypes.Album]](c)
		artistRepo := container.MustGet[repository.MediaItemRepository[*mediatypes.Artist]](c)
		return repository.NewMusicRepository(db, trackRepo, albumRepo, artistRepo)
	})
	
	// Register CoreMusicService
	container.RegisterFactory[services.CoreMusicService](c, func(c *container.Container) services.CoreMusicService {
		musicRepo := container.MustGet[repository.MusicRepository](c)
		return services.NewCoreMusicService(musicRepo)
	})

}

func registerMediaItemService[T mediatypes.MediaData](c *container.Container) {
	container.RegisterFactory[services.CoreMediaItemService[T]](c, func(c *container.Container) services.CoreMediaItemService[T] {
		repo := container.MustGet[repository.MediaItemRepository[T]](c)
		return services.NewCoreMediaItemService[T](repo)
	})
	container.RegisterFactory[services.UserMediaItemService[T]](c, func(c *container.Container) services.UserMediaItemService[T] {
		coreService := container.MustGet[services.CoreMediaItemService[T]](c)
		userRepos := container.MustGet[repository.UserMediaItemRepository[T]](c)
		return services.NewUserMediaItemService[T](coreService, userRepos)
	})
}

func registerClientMediaItemService[T clienttypes.ClientMediaConfig, U mediatypes.MediaData](c *container.Container) {
	container.RegisterFactory[services.ClientMediaItemService[T, U]](c, func(c *container.Container) services.ClientMediaItemService[T, U] {
		coreService := container.MustGet[services.CoreMediaItemService[U]](c)
		clientRepo := container.MustGet[repository.ClientRepository[T]](c)
		itemRepo := container.MustGet[repository.ClientMediaItemRepository[U]](c)
		return services.NewClientMediaItemService[T, U](coreService, clientRepo, itemRepo)
	})
}

// CoreMediaItemServices implementation
type coreMediaItemServicesImpl struct {
	movieCoreService      services.CoreMediaItemService[*mediatypes.Movie]
	seriesCoreService     services.CoreMediaItemService[*mediatypes.Series]
	seasonCoreService     services.CoreMediaItemService[*mediatypes.Season]
	episodeCoreService    services.CoreMediaItemService[*mediatypes.Episode]
	trackCoreService      services.CoreMediaItemService[*mediatypes.Track]
	albumCoreService      services.CoreMediaItemService[*mediatypes.Album]
	artistCoreService     services.CoreMediaItemService[*mediatypes.Artist]
	collectionCoreService services.CoreMediaItemService[*mediatypes.Collection]
	playlistCoreService   services.CoreMediaItemService[*mediatypes.Playlist]
}

func (s *coreMediaItemServicesImpl) MovieCoreService() services.CoreMediaItemService[*mediatypes.Movie] {
	return s.movieCoreService
}

func (s *coreMediaItemServicesImpl) SeriesCoreService() services.CoreMediaItemService[*mediatypes.Series] {
	return s.seriesCoreService
}

func (s *coreMediaItemServicesImpl) SeasonCoreService() services.CoreMediaItemService[*mediatypes.Season] {
	return s.seasonCoreService
}

func (s *coreMediaItemServicesImpl) EpisodeCoreService() services.CoreMediaItemService[*mediatypes.Episode] {
	return s.episodeCoreService
}

func (s *coreMediaItemServicesImpl) TrackCoreService() services.CoreMediaItemService[*mediatypes.Track] {
	return s.trackCoreService
}

func (s *coreMediaItemServicesImpl) AlbumCoreService() services.CoreMediaItemService[*mediatypes.Album] {
	return s.albumCoreService
}

func (s *coreMediaItemServicesImpl) ArtistCoreService() services.CoreMediaItemService[*mediatypes.Artist] {
	return s.artistCoreService
}

func (s *coreMediaItemServicesImpl) CollectionCoreService() services.CoreMediaItemService[*mediatypes.Collection] {
	return s.collectionCoreService
}

func (s *coreMediaItemServicesImpl) PlaylistCoreService() services.CoreMediaItemService[*mediatypes.Playlist] {
	return s.playlistCoreService
}
