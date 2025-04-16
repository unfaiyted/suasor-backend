// app/di/handlers/media.go
package handlers

import (
	"suasor/app/container"
	"suasor/handlers"
)

// RegisterMediaHandlers registers all media-related handlers
func RegisterMediaHandlers(c *container.Container) {
	// Register core media item handlers
	container.RegisterFactory[interfaces.CoreMediaItemHandlers](c, func(c *container.Container) interfaces.CoreMediaItemHandlers {
		factory := container.MustGet[interfaces.MediaDataFactory](c)
		coreServices := container.MustGet[interfaces.CoreMediaItemServices](c)
		return factory.CreateCoreHandlers(coreServices)
	})

	// Register user media item handlers
	container.RegisterFactory[interfaces.UserMediaItemHandlers](c, func(c *container.Container) interfaces.UserMediaItemHandlers {
		factory := container.MustGet[interfaces.MediaDataFactory](c)
		userServices := container.MustGet[interfaces.UserMediaItemServices](c)
		userDataServices := container.MustGet[interfaces.UserMediaItemDataServices](c)
		coreHandlers := container.MustGet[interfaces.CoreMediaItemHandlers](c)
		return factory.CreateUserHandlers(userServices, userDataServices, coreHandlers)
	})

	// Register client media item handlers
	container.RegisterFactory[interfaces.ClientMediaItemHandlers](c, func(c *container.Container) interfaces.ClientMediaItemHandlers {
		factory := container.MustGet[interfaces.MediaDataFactory](c)
		clientServices := container.MustGet[interfaces.ClientMediaItemServices](c)
		clientDataServices := container.MustGet[interfaces.ClientUserMediaItemDataServices](c)
		userHandlers := container.MustGet[interfaces.UserMediaItemHandlers](c)
		return factory.CreateClientHandlers(clientServices, clientDataServices, userHandlers)
	})

	// Register specialized handlers
	registerSpecializedMediaHandlers(c)
}

// Register specialized media handlers for specific domains like music, movies, etc.
func registerSpecializedMediaHandlers(c *container.Container) {
	// Core music handler
	container.RegisterFactory[*handlers.CoreMusicHandler](c, func(c *container.Container) *handlers.CoreMusicHandler {
		clientServices := container.MustGet[interfaces.ClientMediaItemServices](c)
		return handlers.NewCoreMusicHandler(
			clientServices.TrackClientService(),
			clientServices.AlbumClientService(),
			clientServices.ArtistClientService(),
		)
	})

	// Core movie handler
	container.RegisterFactory[*handlers.CoreMovieHandler](c, func(c *container.Container) *handlers.CoreMovieHandler {
		coreMovieService := container.MustGet[interfaces.CoreMediaItemServices](c).MovieCoreService()
		return handlers.NewCoreMovieHandler(coreMovieService)
	})

	// Core series handler
	container.RegisterFactory[*handlers.CoreSeriesHandler](c, func(c *container.Container) *handlers.CoreSeriesHandler {
		coreSeriesService := container.MustGet[interfaces.CoreMediaItemServices](c).SeriesCoreService()
		return handlers.NewCoreSeriesHandler(coreSeriesService)
	})
}
