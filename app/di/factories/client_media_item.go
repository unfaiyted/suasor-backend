package factories

import (
	"context"

	"suasor/app/container"
	"suasor/client/media"
	"suasor/client/media/emby"
	"suasor/client/media/jellyfin"
	"suasor/client/media/plex"
	"suasor/client/media/subsonic"
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
