// app/di/handlers/media.go
package handlers

import (
	"context"
	"suasor/app/container"
	"suasor/app/di/factories"
	apphandlers "suasor/app/handlers"
	"suasor/app/services"
	"suasor/handlers"
)

// RegisterMediaHandlers registers all media-related handlers
func RegisterMediaHandlers(ctx context.Context, c *container.Container) {
	// Register core media item handlers
	container.RegisterFactory[apphandlers.CoreMediaItemHandlers](c, func(c *container.Container) apphandlers.CoreMediaItemHandlers {
		factory := container.MustGet[factories.MediaDataFactory](c)
		coreServices := container.MustGet[services.CoreMediaItemServices](c)
		return factory.CreateCoreMediaItemHandlers(coreServices)
	})

	// Register user media item handlers
	container.RegisterFactory[apphandlers.UserMediaItemHandlers](c, func(c *container.Container) apphandlers.UserMediaItemHandlers {
		factory := container.MustGet[factories.MediaDataFactory](c)
		userServices := container.MustGet[services.UserMediaItemServices](c)
		coreHandlers := container.MustGet[apphandlers.CoreMediaItemHandlers](c)
		return factory.CreateUserMediaItemHandlers(userServices, coreHandlers)
	})

	// Register client media item handlers
	container.RegisterFactory[apphandlers.ClientMediaItemHandlers](c, func(c *container.Container) apphandlers.ClientMediaItemHandlers {
		factory := container.MustGet[factories.MediaDataFactory](c)
		clientServices := container.MustGet[services.ClientMediaItemServices](c)
		userServices := container.MustGet[services.UserMediaItemServices](c)
		userHandlers := container.MustGet[apphandlers.UserMediaItemHandlers](c)
		return factory.CreateClientMediaItemHandlers(clientServices, userServices, userHandlers)
	})

	// Register specialized handlers
	registerSpecializedMediaHandlers(c)
}

// Register specialized media handlers for specific domains like music, movies, etc.
func registerSpecializedMediaHandlers(c *container.Container) {
	// Core music handler
	container.RegisterFactory[*handlers.CoreMusicHandler](c, func(c *container.Container) *handlers.CoreMusicHandler {
		clientServices := container.MustGet[services.ClientMediaItemServices](c)
		return handlers.NewCoreMusicHandler(
			clientServices.TrackClientService(),
			clientServices.AlbumClientService(),
			clientServices.ArtistClientService(),
		)
	})

	// Core movie handler
	container.RegisterFactory[*handlers.CoreMovieHandler](c, func(c *container.Container) *handlers.CoreMovieHandler {
		coreMovieService := container.MustGet[services.CoreMediaItemServices](c).MovieCoreService()
		return handlers.NewCoreMovieHandler(coreMovieService)
	})

	// Core series handler
	container.RegisterFactory[*handlers.CoreSeriesHandler](c, func(c *container.Container) *handlers.CoreSeriesHandler {
		coreServices := container.MustGet[services.CoreMediaItemServices](c)
		return handlers.NewCoreSeriesHandler(
			coreServices.SeriesCoreService(),
			nil, // We don't have a Season service yet defined in the interfaces
			coreServices.EpisodeCoreService(),
		)
	})
}
