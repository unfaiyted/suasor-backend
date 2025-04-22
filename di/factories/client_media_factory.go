package factories

import (
	"context"
	"fmt"

	"suasor/clients"
	"suasor/clients/media"
	emby "suasor/clients/media/emby"
	jellyfin "suasor/clients/media/jellyfin"
	plex "suasor/clients/media/plex"
	subsonic "suasor/clients/media/subsonic"
	clienttypes "suasor/clients/types"
	"suasor/container"
	"suasor/utils/logger"
)

func RegisterClientFactories(ctx context.Context, c *container.Container) {

	// Server to register clients
	// Get the service directly to avoid dependency issues
	service := clients.GetClientFactoryService()

	// Register it for others to use
	container.RegisterFactory[*clients.ClientFactoryService](c, func(c *container.Container) *clients.ClientFactoryService {
		return service
	})

	// Get the registry directly to avoid dependency issues
	registry := container.MustGet[media.ClientItemRegistry](c)

	// EMBY CLIENT
	service.RegisterClientFactory(clienttypes.ClientTypeEmby, func(ctx context.Context, clientID uint64, config clienttypes.ClientConfig) (clients.Client, error) {
		log := logger.LoggerFromContext(ctx)
		// Use the provided config (should be an EmbyConfig)
		embyConfig, ok := config.(*clienttypes.EmbyConfig)
		log.Debug().
			Bool("ok", ok).
			Msg("Checking config type")
		if !ok {
			log.Error().
				Err(fmt.Errorf("expected *config.EmbyConfig, got %T", config)).
				Msg("Expected *config.EmbyConfig, got %T")
			return nil, fmt.Errorf("expected *config.EmbyConfig, got %T", config)
		}

		fmt.Printf("Factory called for Emby client with ID: %d\n", clientID)
		return emby.NewEmbyClient(ctx, &registry, clientID, *embyConfig)
	})

	// JELLYFIN CLIENT
	service.RegisterClientFactory(clienttypes.ClientTypeJellyfin, func(ctx context.Context, clientID uint64, config clienttypes.ClientConfig) (clients.Client, error) {
		log := logger.LoggerFromContext(ctx)
		// Use the provided config (should be an EmbyConfig)
		jellyfinConfig, ok := config.(*clienttypes.JellyfinConfig)
		log.Debug().
			Bool("ok", ok).
			Msg("Checking config type")
		if !ok {
			log.Error().
				Err(fmt.Errorf("expected *config.JellyfinConfig, got %T", config)).
				Msg("Expected *config.JellyfinConfig, got %T")
			return nil, fmt.Errorf("expected *config.JellyfinConfig, got %T", config)
		}

		fmt.Printf("Factory called for Jellyfin client with ID: %d\n", clientID)
		return jellyfin.NewJellyfinClient(ctx, &registry, clientID, *jellyfinConfig)
	})
	// PLEX CLIENT
	service.RegisterClientFactory(clienttypes.ClientTypePlex, func(ctx context.Context, clientID uint64, config clienttypes.ClientConfig) (clients.Client, error) {
		log := logger.LoggerFromContext(ctx)
		// Use the provided config (should be an EmbyConfig)
		plexConfig, ok := config.(*clienttypes.PlexConfig)
		log.Debug().
			Bool("ok", ok).
			Msg("Checking config type")
		if !ok {
			log.Error().
				Err(fmt.Errorf("expected *config.PlexConfig, got %T", config)).
				Msg("Expected *config.PlexConfig, got %T")
			return nil, fmt.Errorf("expected *config.PlexConfig, got %T", config)
		}

		fmt.Printf("Factory called for Plex client with ID: %d\n", clientID)
		return plex.NewPlexClient(ctx, &registry, clientID, *plexConfig)
	})

	// SUBSONIC CLIENT
	service.RegisterClientFactory(clienttypes.ClientTypeSubsonic, func(ctx context.Context, clientID uint64, config clienttypes.ClientConfig) (clients.Client, error) {
		log := logger.LoggerFromContext(ctx)
		// Use the provided config (should be an EmbyConfig)
		subsonicConfig, ok := config.(*clienttypes.SubsonicConfig)
		log.Debug().
			Bool("ok", ok).
			Msg("Checking config type")
		if !ok {
			log.Error().
				Err(fmt.Errorf("expected *config.SubsonicConfig, got %T", config)).
				Msg("Expected *config.SubsonicConfig, got %T")
			return nil, fmt.Errorf("expected *config.SubsonicConfig, got %T", config)
		}

		fmt.Printf("Factory called for Subsonic client with ID: %d\n", clientID)
		return subsonic.NewSubsonicClient(ctx, &registry, clientID, *subsonicConfig)
	})

}
