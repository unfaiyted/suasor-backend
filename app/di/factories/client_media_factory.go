package factories

import (
	"context"
	"fmt"

	"suasor/app/container"
	"suasor/client"
	"suasor/client/media"
	emby "suasor/client/media/emby"
	jellyfin "suasor/client/media/jellyfin"
	plex "suasor/client/media/plex"
	subsonic "suasor/client/media/subsonic"
	clienttypes "suasor/client/types"
	"suasor/utils"
)

func RegisterClientFactories(ctx context.Context, c *container.Container) {

	// Server to register clients
	container.RegisterFactory[*client.ClientFactoryService](c, func(c *container.Container) *client.ClientFactoryService {
		return client.GetClientFactoryService()
	})

	registry := container.MustGet[media.ClientItemRegistry](c)
	service := container.MustGet[client.ClientFactoryService](c)

	// EMBY CLIENT
	service.RegisterClientFactory(clienttypes.ClientTypeEmby, func(ctx context.Context, clientID uint64, config clienttypes.ClientConfig) (client.Client, error) {
		log := utils.LoggerFromContext(ctx)
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
	service.RegisterClientFactory(clienttypes.ClientTypeJellyfin, func(ctx context.Context, clientID uint64, config clienttypes.ClientConfig) (client.Client, error) {
		log := utils.LoggerFromContext(ctx)
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
	service.RegisterClientFactory(clienttypes.ClientTypePlex, func(ctx context.Context, clientID uint64, config clienttypes.ClientConfig) (client.Client, error) {
		log := utils.LoggerFromContext(ctx)
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
	service.RegisterClientFactory(clienttypes.ClientTypeSubsonic, func(ctx context.Context, clientID uint64, config clienttypes.ClientConfig) (client.Client, error) {
		log := utils.LoggerFromContext(ctx)
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
