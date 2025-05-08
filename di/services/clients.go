package services

import (
	"context"
	"suasor/clients"
	mediatypes "suasor/clients/media/types"
	types "suasor/clients/types"
	"suasor/di/container"
	"suasor/repository"
	"suasor/services"
	svcbundles "suasor/services/bundles"
)

func registerClientServices(ctx context.Context, c *container.Container) {
	// Media clients
	registerClientService[*types.JellyfinConfig](c)
	registerClientService[*types.EmbyConfig](c)
	registerClientService[*types.PlexConfig](c)
	registerClientService[*types.SubsonicConfig](c)
	registerClientService[*types.RadarrConfig](c)
	registerClientService[*types.LidarrConfig](c)
	registerClientService[*types.SonarrConfig](c)
	registerClientService[*types.ClaudeConfig](c)
	registerClientService[*types.OpenAIConfig](c)
	registerClientService[*types.OllamaConfig](c)
	
	// AI client registration
	registerClientService[types.AIClientConfig](c)

	// Register ClientSeriesService for each media client type
	registerClientSeriesService[*types.JellyfinConfig](c)
	registerClientSeriesService[*types.EmbyConfig](c)
	registerClientSeriesService[*types.PlexConfig](c)

	// Register ClientMusicService for each media client type
	registerClientMusicService[*types.JellyfinConfig](c)
	registerClientMusicService[*types.EmbyConfig](c)
	registerClientMusicService[*types.PlexConfig](c)
	registerClientMusicService[*types.SubsonicConfig](c)

	// Register client music services bundle
	registerClientMusicServicesBundle(c)

	// Register AutomationClientService
	container.RegisterFactory[services.AutomationClientService](c, func(c *container.Container) services.AutomationClientService {
		clientFactory := container.MustGet[*clients.ClientProviderFactoryService](c)
		repo := container.MustGet[repository.ClientRepository[types.ClientAutomationConfig]](c)
		return services.NewAutomationClientService(repo, clientFactory)
	})

}

func registerClientService[T types.ClientConfig](c *container.Container) {
	container.RegisterFactory[services.ClientService[T]](c, func(c *container.Container) services.ClientService[T] {
		clientFactory := container.MustGet[*clients.ClientProviderFactoryService](c)
		repo := container.MustGet[repository.ClientRepository[T]](c)
		return services.NewClientService[T](clientFactory, repo)
	})
}

func registerMediaTypeService[T types.ClientMediaConfig, U mediatypes.MediaData](c *container.Container) {
	container.RegisterFactory[services.ClientMediaItemService[T, U]](c, func(c *container.Container) services.ClientMediaItemService[T, U] {
		// Dependencies
		clientFactory := container.MustGet[*clients.ClientProviderFactoryService](c)
		clientRepo := container.MustGet[repository.ClientRepository[T]](c)
		itemRepo := container.MustGet[repository.ClientMediaItemRepository[U]](c)
		userService := container.MustGet[services.UserMediaItemService[U]](c)

		return services.NewClientMediaItemService[T, U](
			userService,
			clientRepo,
			itemRepo,
			clientFactory,
		)
	})
}

// registerClientSeriesService registers a specialized series service for a given client config type
func registerClientSeriesService[T types.ClientMediaConfig](c *container.Container) {
	container.RegisterFactory[services.ClientSeriesService[T]](c, func(c *container.Container) services.ClientSeriesService[T] {
		clientRepo := container.MustGet[repository.ClientRepository[T]](c)
		clientFactory := container.MustGet[*clients.ClientProviderFactoryService](c)
		return services.NewClientSeriesService[T](clientRepo, clientFactory)
	})
}

// registerClientMusicService registers a specialized music service for a given client config type
func registerClientMusicService[T types.ClientMediaConfig](c *container.Container) {
	container.RegisterFactory[services.ClientMusicService[T]](c, func(c *container.Container) services.ClientMusicService[T] {
		clientRepo := container.MustGet[repository.ClientRepository[T]](c)
		clientFactory := container.MustGet[*clients.ClientProviderFactoryService](c)
		coreService := container.MustGet[services.CoreMusicService](c)
		return services.NewClientMusicService[T](coreService, clientRepo, clientFactory)
	})
}

// ClientMusicServicesImpl is a bundled implementation of the ClientMusicServices interface
type ClientMusicServicesImpl struct {
	embyService     services.ClientMusicService[*types.EmbyConfig]
	jellyfinService services.ClientMusicService[*types.JellyfinConfig]
	plexService     services.ClientMusicService[*types.PlexConfig]
	subsonicService services.ClientMusicService[*types.SubsonicConfig]
}

// EmbyMusicService returns the Emby music service
func (s *ClientMusicServicesImpl) EmbyMusicService() services.ClientMusicService[*types.EmbyConfig] {
	return s.embyService
}

// JellyfinMusicService returns the Jellyfin music service
func (s *ClientMusicServicesImpl) JellyfinMusicService() services.ClientMusicService[*types.JellyfinConfig] {
	return s.jellyfinService
}

// PlexMusicService returns the Plex music service
func (s *ClientMusicServicesImpl) PlexMusicService() services.ClientMusicService[*types.PlexConfig] {
	return s.plexService
}

// SubsonicMusicService returns the Subsonic music service
func (s *ClientMusicServicesImpl) SubsonicMusicService() services.ClientMusicService[*types.SubsonicConfig] {
	return s.subsonicService
}

// registerClientMusicServicesBundle registers the ClientMusicServices implementation
func registerClientMusicServicesBundle(c *container.Container) {
	container.RegisterFactory[svcbundles.ClientMusicServices](c, func(c *container.Container) svcbundles.ClientMusicServices {
		return &ClientMusicServicesImpl{
			embyService:     container.MustGet[services.ClientMusicService[*types.EmbyConfig]](c),
			jellyfinService: container.MustGet[services.ClientMusicService[*types.JellyfinConfig]](c),
			plexService:     container.MustGet[services.ClientMusicService[*types.PlexConfig]](c),
			subsonicService: container.MustGet[services.ClientMusicService[*types.SubsonicConfig]](c),
		}
	})
}
