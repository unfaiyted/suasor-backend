// app/di/handlers/media.go
package handlers

import (
	"context"
	mediatypes "suasor/clients/media/types"
	clienttypes "suasor/clients/types"
	"suasor/di/container"
	"suasor/handlers"
	"suasor/services"
	svcbundles "suasor/services/bundles"
)

// RegisterMediaHandlers registers all media-related handlers
func RegisterMediaItemHandlers(ctx context.Context, c *container.Container) {

	registerCoreMediaItemHandler[*mediatypes.Movie](c)
	registerCoreMediaItemHandler[*mediatypes.Series](c)
	registerCoreMediaItemHandler[*mediatypes.Season](c)
	registerCoreMediaItemHandler[*mediatypes.Episode](c)
	registerCoreMediaItemHandler[*mediatypes.Track](c)
	registerCoreMediaItemHandler[*mediatypes.Album](c)
	registerCoreMediaItemHandler[*mediatypes.Artist](c)
	registerCoreMediaItemHandler[*mediatypes.Collection](c)
	registerCoreMediaItemHandler[*mediatypes.Playlist](c)

	registerUserMediaItemHandler[*mediatypes.Movie](c)
	registerUserMediaItemHandler[*mediatypes.Series](c)
	registerUserMediaItemHandler[*mediatypes.Season](c)
	registerUserMediaItemHandler[*mediatypes.Episode](c)
	registerUserMediaItemHandler[*mediatypes.Track](c)
	registerUserMediaItemHandler[*mediatypes.Album](c)
	registerUserMediaItemHandler[*mediatypes.Artist](c)
	registerUserMediaItemHandler[*mediatypes.Collection](c)
	registerUserMediaItemHandler[*mediatypes.Playlist](c)

	registerClientMediaItemHandler[*clienttypes.EmbyConfig, *mediatypes.Movie](c)
	registerClientMediaItemHandler[*clienttypes.JellyfinConfig, *mediatypes.Movie](c)
	registerClientMediaItemHandler[*clienttypes.PlexConfig, *mediatypes.Movie](c)

	registerClientMediaItemHandler[*clienttypes.EmbyConfig, *mediatypes.Series](c)
	registerClientMediaItemHandler[*clienttypes.JellyfinConfig, *mediatypes.Series](c)
	registerClientMediaItemHandler[*clienttypes.PlexConfig, *mediatypes.Series](c)

	registerClientMediaItemHandler[*clienttypes.EmbyConfig, *mediatypes.Season](c)
	registerClientMediaItemHandler[*clienttypes.JellyfinConfig, *mediatypes.Season](c)
	registerClientMediaItemHandler[*clienttypes.PlexConfig, *mediatypes.Season](c)

	registerClientMediaItemHandler[*clienttypes.EmbyConfig, *mediatypes.Episode](c)
	registerClientMediaItemHandler[*clienttypes.JellyfinConfig, *mediatypes.Episode](c)
	registerClientMediaItemHandler[*clienttypes.PlexConfig, *mediatypes.Episode](c)

	registerClientMediaItemHandler[*clienttypes.EmbyConfig, *mediatypes.Track](c)
	registerClientMediaItemHandler[*clienttypes.JellyfinConfig, *mediatypes.Track](c)
	registerClientMediaItemHandler[*clienttypes.PlexConfig, *mediatypes.Track](c)
	registerClientMediaItemHandler[*clienttypes.SubsonicConfig, *mediatypes.Track](c)

	registerClientMediaItemHandler[*clienttypes.EmbyConfig, *mediatypes.Album](c)
	registerClientMediaItemHandler[*clienttypes.JellyfinConfig, *mediatypes.Album](c)
	registerClientMediaItemHandler[*clienttypes.PlexConfig, *mediatypes.Album](c)
	registerClientMediaItemHandler[*clienttypes.SubsonicConfig, *mediatypes.Album](c)

	registerClientMediaItemHandler[*clienttypes.EmbyConfig, *mediatypes.Artist](c)
	registerClientMediaItemHandler[*clienttypes.JellyfinConfig, *mediatypes.Artist](c)
	registerClientMediaItemHandler[*clienttypes.PlexConfig, *mediatypes.Artist](c)
	registerClientMediaItemHandler[*clienttypes.SubsonicConfig, *mediatypes.Artist](c)

	registerClientMediaItemHandler[*clienttypes.EmbyConfig, *mediatypes.Collection](c)
	registerClientMediaItemHandler[*clienttypes.JellyfinConfig, *mediatypes.Collection](c)
	registerClientMediaItemHandler[*clienttypes.PlexConfig, *mediatypes.Collection](c)

	registerClientMediaItemHandler[*clienttypes.EmbyConfig, *mediatypes.Playlist](c)
	registerClientMediaItemHandler[*clienttypes.JellyfinConfig, *mediatypes.Playlist](c)
	registerClientMediaItemHandler[*clienttypes.PlexConfig, *mediatypes.Playlist](c)

	// Register specialized handlers
	registerSpecializedMediaHandlers(c)
}

// Register specialized media handlers for specific domains like music, movies, etc.
func registerSpecializedMediaHandlers(c *container.Container) {
	// People handler
	container.RegisterFactory[*handlers.PeopleHandler](c, func(c *container.Container) *handlers.PeopleHandler {
		personService := container.MustGet[*services.PersonService](c)
		return handlers.NewPeopleHandler(personService)
	})

	// Credit handler
	container.RegisterFactory[*handlers.CreditHandler](c, func(c *container.Container) *handlers.CreditHandler {
		creditService := container.MustGet[*services.CreditService](c)
		return handlers.NewCreditHandler(creditService)
	})

	// Core music handler
	container.RegisterFactory[handlers.CoreMusicHandler](c, func(c *container.Container) handlers.CoreMusicHandler {
		clientServices := container.MustGet[svcbundles.CoreMediaItemServices](c)
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
		coreHandler := container.MustGet[handlers.CoreMediaItemHandler[*mediatypes.Movie]](c)
		itemService := container.MustGet[services.CoreMediaItemService[*mediatypes.Movie]](c)
		return handlers.NewCoreMovieHandler(coreHandler, itemService)
	})

	// Core series handler
	container.RegisterFactory[handlers.CoreSeriesHandler](c, func(c *container.Container) handlers.CoreSeriesHandler {
		coreHandler := container.MustGet[handlers.CoreMediaItemHandler[*mediatypes.Series]](c)
		seriesService := container.MustGet[services.CoreMediaItemService[*mediatypes.Series]](c)
		seasonService := container.MustGet[services.CoreMediaItemService[*mediatypes.Season]](c)
		episodeService := container.MustGet[services.CoreMediaItemService[*mediatypes.Episode]](c)
		return handlers.NewCoreSeriesHandler(coreHandler, seriesService, seasonService, episodeService)
	})
	
	// User series handler
	container.RegisterFactory[handlers.UserSeriesHandler](c, func(c *container.Container) handlers.UserSeriesHandler {
		coreHandler := container.MustGet[handlers.CoreSeriesHandler](c)
		
		// Item Services
		seriesService := container.MustGet[services.UserMediaItemService[*mediatypes.Series]](c)
		seasonService := container.MustGet[services.UserMediaItemService[*mediatypes.Season]](c)
		episodeService := container.MustGet[services.UserMediaItemService[*mediatypes.Episode]](c)
		
		// Data Services
		seriesDataService := container.MustGet[services.UserMediaItemDataService[*mediatypes.Series]](c)
		seasonDataService := container.MustGet[services.UserMediaItemDataService[*mediatypes.Season]](c)
		episodeDataService := container.MustGet[services.UserMediaItemDataService[*mediatypes.Episode]](c)
		
		return handlers.NewUserSeriesHandler(
			coreHandler,
			seriesService,
			seasonService,
			episodeService,
			seriesDataService,
			seasonDataService,
			episodeDataService,
		)
	})
}

func registerCoreMediaItemHandler[T mediatypes.MediaData](c *container.Container) {
	container.RegisterFactory[handlers.CoreMediaItemHandler[T]](c, func(c *container.Container) handlers.CoreMediaItemHandler[T] {
		coreServices := container.MustGet[services.CoreMediaItemService[T]](c)
		return handlers.NewCoreMediaItemHandler[T](coreServices)
	})
}
func registerUserMediaItemHandler[T mediatypes.MediaData](c *container.Container) {
	container.RegisterFactory[handlers.UserMediaItemHandler[T]](c, func(c *container.Container) handlers.UserMediaItemHandler[T] {
		coreHandlers := container.MustGet[handlers.CoreMediaItemHandler[T]](c)
		userServices := container.MustGet[services.UserMediaItemService[T]](c)
		return handlers.NewUserMediaItemHandler[T](coreHandlers, userServices)
	})
}
func registerClientMediaItemHandler[T clienttypes.ClientMediaConfig, U mediatypes.MediaData](c *container.Container) {
	container.RegisterFactory[handlers.ClientMediaItemHandler[T, U]](c, func(c *container.Container) handlers.ClientMediaItemHandler[T, U] {
		userHander := container.MustGet[handlers.UserMediaItemHandler[U]](c)
		clientService := container.MustGet[services.ClientMediaItemService[T, U]](c)
		return handlers.NewClientMediaItemHandler[T, U](userHander, clientService)
	})

}
