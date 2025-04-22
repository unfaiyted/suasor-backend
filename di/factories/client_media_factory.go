package factories

import (
	"context"
	"fmt"

	"suasor/clients"
	// lidarr "suasor/clients/automation/lidarr"
	// radarr "suasor/clients/automation/radarr"
	// sonarr "suasor/clients/automation/sonarr"
	"suasor/clients/media"
	emby "suasor/clients/media/emby"
	jellyfin "suasor/clients/media/jellyfin"
	plex "suasor/clients/media/plex"
	subsonic "suasor/clients/media/subsonic"
	clienttypes "suasor/clients/types"
	"suasor/di/container"
	"suasor/utils"
	"suasor/utils/logger"
)

func RegisterClientFactories(ctx context.Context, c *container.Container) {
	// Server to register clients
	// Get the service directly to avoid dependency issues
	service := clients.GetClientProviderFactoryService()

	// Register it for others to use
	// This provider factory provides the actual operations for creationg the workig clients and getting data from the clients.
	container.RegisterFactory[*clients.ClientProviderFactoryService](c, func(c *container.Container) *clients.ClientProviderFactoryService {
		return service
	})

	// Get the registry directly to avoid dependency issues
	registry := container.MustGet[media.ClientItemRegistry](c)

	// Register all clients using the generic approach
	registerClientProviders[*clienttypes.EmbyConfig](c, service, &registry, createMediaClientWrapper(emby.NewEmbyClient))
	registerClientProviders[*clienttypes.JellyfinConfig](c, service, &registry, createMediaClientWrapper(jellyfin.NewJellyfinClient))
	registerClientProviders[*clienttypes.PlexConfig](c, service, &registry, createMediaClientWrapper(plex.NewPlexClient))
	registerClientProviders[*clienttypes.SubsonicConfig](c, service, &registry, createMediaClientWrapper(subsonic.NewSubsonicClient))
	// registerClientProviders[*clienttypes.RadarrConfig](c, service, &registry, createAutomationClientWrapper(radarr.NewRadarrClient))
	// registerClientProviders[*clienttypes.LidarrConfig](c, service, &registry, createAutomationClientWrapper(lidarr.NewLidarrClient))
	// registerClientProviders[*clienttypes.SonarrConfig](c, service, &registry, createAutomationClientWrapper(sonarr.NewSonarrClient))

}

func registerClientProviders[T clienttypes.ClientMediaConfig](c *container.Container, service *clients.ClientProviderFactoryService, registry *media.ClientItemRegistry, factory func(ctx context.Context, registry *media.ClientItemRegistry, clientID uint64, config T) (clients.Client, error)) {
	// Register the client connection provider and its associated media providers

	var zero T
	typeName := utils.GetTypeName(zero)
	clientType := clienttypes.GetClientTypeFromTypeName(typeName)

	service.RegisterClientProviderFactory(clientType, func(ctx context.Context, clientID uint64, config clienttypes.ClientConfig) (clients.Client, error) {
		log := logger.LoggerFromContext(ctx)
		clientConfig, ok := config.(T)
		log.Debug().
			Bool("ok", ok).
			Msg("Checking config type")
		if !ok {
			log.Error().
				Err(fmt.Errorf("expected %s, got %T", typeName, config)).
				Msg("Expected %s, got %T")
			return nil, fmt.Errorf("expected %s, got %T", typeName, config)
		}

		fmt.Printf("Factory called for client with ID: %d\n", clientID)
		return factory(ctx, registry, clientID, clientConfig)
	})
}

// Convert a function returning media.ClientMedia to one returning clients.Client
func createMediaClientWrapper[T clienttypes.ClientMediaConfig](
	originalFactory func(ctx context.Context, registry *media.ClientItemRegistry, clientID uint64, config T) (media.ClientMedia, error),
) func(ctx context.Context, registry *media.ClientItemRegistry, clientID uint64, config T) (clients.Client, error) {
	return func(ctx context.Context, registry *media.ClientItemRegistry, clientID uint64, config T) (clients.Client, error) {
		return originalFactory(ctx, registry, clientID, config)
	}
}

func createAutomationClientWrapper[T clienttypes.AutomationClientConfig](
	originalFactory func(ctx context.Context, registry *media.ClientItemRegistry, clientID uint64, config T) (media.ClientMedia, error),
) func(ctx context.Context, registry *media.ClientItemRegistry, clientID uint64, config T) (clients.Client, error) {
	return func(ctx context.Context, registry *media.ClientItemRegistry, clientID uint64, config T) (clients.Client, error) {
		return originalFactory(ctx, registry, clientID, config)
	}
}

