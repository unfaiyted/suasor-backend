// app/di/handlers/media.go
package handlers

import (
	"context"
	"suasor/app/container"
	"suasor/app/di/factories"
	apphandlers "suasor/app/handlers"
	appservices "suasor/app/services"
	"suasor/handlers"
	"suasor/services"
)

// RegisterMediaHandlers registers all media-related handlers
func RegisterMediaHandlers(ctx context.Context, c *container.Container) {
	// Register core media item handlers
	container.RegisterFactory[apphandlers.CoreMediaItemHandlers](c, func(c *container.Container) apphandlers.CoreMediaItemHandlers {
		factory := container.MustGet[factories.MediaDataFactory](c)
		coreServices := container.MustGet[appservices.CoreMediaItemServices](c)
		return factory.CreateCoreMediaItemHandlers(coreServices)
	})

	// Register user media item handlers
	container.RegisterFactory[apphandlers.UserMediaItemHandlers](c, func(c *container.Container) apphandlers.UserMediaItemHandlers {
		factory := container.MustGet[factories.MediaDataFactory](c)
		userServices := container.MustGet[appservices.UserMediaItemServices](c)
		coreHandlers := container.MustGet[apphandlers.CoreMediaItemHandlers](c)
		return factory.CreateUserMediaItemHandlers(userServices, coreHandlers)
	})

	// Register client media item handlers
	container.RegisterFactory[apphandlers.ClientMediaItemHandlers](c, func(c *container.Container) apphandlers.ClientMediaItemHandlers {
		factory := container.MustGet[factories.MediaDataFactory](c)
		clientServices := container.MustGet[appservices.ClientMediaItemServices](c)
		userServices := container.MustGet[appservices.UserMediaItemServices](c)
		userHandlers := container.MustGet[apphandlers.UserMediaItemHandlers](c)
		return factory.CreateClientMediaItemHandlers(clientServices, userServices, userHandlers)
	})

	// Register specialized handlers
	registerSpecializedMediaHandlers(c)
}

// Register specialized media handlers for specific domains like music, movies, etc.
func registerSpecializedMediaHandlers(c *container.Container) {
	// Core music handler
	container.RegisterFactory[handlers.CoreMusicHandler](c, func(c *container.Container) handlers.CoreMusicHandler {
		clientServices := container.MustGet[appservices.CoreMediaItemServices](c)
		coreMusicService := container.MustGet[services.CoreMusicService](c)
		return handlers.NewCoreMusicHandler(
			coreMusicService,
			clientServices.TrackCoreService(),
			clientServices.AlbumCoreService(),
			clientServices.ArtistCoreService(),
		)
	})

	// Core movie handler
	container.RegisterFactory[handlers.CoreMovieHandler](c, func(c *container.Container) handlers.CoreMovieHandler {
		coreHandlers := container.MustGet[apphandlers.CoreMediaItemHandlers](c).MovieCoreHandler()
		itemService := container.MustGet[appservices.CoreMediaItemServices](c).MovieCoreService()
		return handlers.NewCoreMovieHandler(coreHandlers, itemService)
	})

	// Core series handler
	container.RegisterFactory[handlers.CoreSeriesHandler](c, func(c *container.Container) handlers.CoreSeriesHandler {
		coreServices := container.MustGet[appservices.CoreMediaItemServices](c)
		return handlers.NewCoreSeriesHandler(
			coreServices.SeriesCoreService(),
			coreServices.SeasonCoreService(),
			coreServices.EpisodeCoreService(),
		)
	})
}
