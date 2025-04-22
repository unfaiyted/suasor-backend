package factories

import (
	"context"

	"suasor/clients/media"
	"suasor/clients/media/emby"
	"suasor/clients/media/jellyfin"
	"suasor/clients/media/plex"
	"suasor/clients/media/subsonic"
	"suasor/container"
)

func RegisterClientMediaItemFactories(ctx context.Context, c *container.Container) {

	// Registry to allow a client to create different Items factories
	container.RegisterFactory[media.ClientItemRegistry](c, func(c *container.Container) media.ClientItemRegistry {
		return *media.NewClientItemRegistry()
	})

	emby.RegisterMediaItemFactories(c)
	plex.RegisterMediaItemFactories(c)
	jellyfin.RegisterMediaItemFactories(c)
	subsonic.RegisterMediaItemFactories(c)

}
