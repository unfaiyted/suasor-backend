// app/di/handlers/media_data.go
package handlers

import (
	"context"
	mediatypes "suasor/clients/media/types"
	clienttypes "suasor/clients/types"
	"suasor/di/container"
	"suasor/handlers"
	"suasor/services"
)

// RegisterMediaDataHandlers registers all media data-related handlers
func RegisterMediaDataHandlers(ctx context.Context, c *container.Container) {

	registerCoreUserMediaItemDataHandler[*mediatypes.Movie](c)
	registerCoreUserMediaItemDataHandler[*mediatypes.Series](c)
	registerCoreUserMediaItemDataHandler[*mediatypes.Season](c)
	registerCoreUserMediaItemDataHandler[*mediatypes.Episode](c)
	registerCoreUserMediaItemDataHandler[*mediatypes.Track](c)
	registerCoreUserMediaItemDataHandler[*mediatypes.Album](c)
	registerCoreUserMediaItemDataHandler[*mediatypes.Artist](c)
	registerCoreUserMediaItemDataHandler[*mediatypes.Collection](c)
	registerCoreUserMediaItemDataHandler[*mediatypes.Playlist](c)

	registerUserMediaItemDataHandler[*mediatypes.Movie](c)
	registerUserMediaItemDataHandler[*mediatypes.Series](c)
	registerUserMediaItemDataHandler[*mediatypes.Season](c)
	registerUserMediaItemDataHandler[*mediatypes.Episode](c)
	registerUserMediaItemDataHandler[*mediatypes.Track](c)
	registerUserMediaItemDataHandler[*mediatypes.Album](c)
	registerUserMediaItemDataHandler[*mediatypes.Artist](c)
	registerUserMediaItemDataHandler[*mediatypes.Collection](c)
	registerUserMediaItemDataHandler[*mediatypes.Playlist](c)

	registerClientUserMediaItemDataHandler[*clienttypes.EmbyConfig, *mediatypes.Movie](c)
	registerClientUserMediaItemDataHandler[*clienttypes.JellyfinConfig, *mediatypes.Movie](c)
	registerClientUserMediaItemDataHandler[*clienttypes.PlexConfig, *mediatypes.Movie](c)

	registerClientUserMediaItemDataHandler[*clienttypes.EmbyConfig, *mediatypes.Series](c)
	registerClientUserMediaItemDataHandler[*clienttypes.JellyfinConfig, *mediatypes.Series](c)
	registerClientUserMediaItemDataHandler[*clienttypes.PlexConfig, *mediatypes.Series](c)

	registerClientUserMediaItemDataHandler[*clienttypes.EmbyConfig, *mediatypes.Season](c)
	registerClientUserMediaItemDataHandler[*clienttypes.JellyfinConfig, *mediatypes.Season](c)
	registerClientUserMediaItemDataHandler[*clienttypes.PlexConfig, *mediatypes.Season](c)

	registerClientUserMediaItemDataHandler[*clienttypes.EmbyConfig, *mediatypes.Episode](c)
	registerClientUserMediaItemDataHandler[*clienttypes.JellyfinConfig, *mediatypes.Episode](c)
	registerClientUserMediaItemDataHandler[*clienttypes.PlexConfig, *mediatypes.Episode](c)

	registerClientUserMediaItemDataHandler[*clienttypes.EmbyConfig, *mediatypes.Track](c)
	registerClientUserMediaItemDataHandler[*clienttypes.JellyfinConfig, *mediatypes.Track](c)
	registerClientUserMediaItemDataHandler[*clienttypes.PlexConfig, *mediatypes.Track](c)
	registerClientUserMediaItemDataHandler[*clienttypes.SubsonicConfig, *mediatypes.Track](c)

	registerClientUserMediaItemDataHandler[*clienttypes.EmbyConfig, *mediatypes.Album](c)
	registerClientUserMediaItemDataHandler[*clienttypes.JellyfinConfig, *mediatypes.Album](c)
	registerClientUserMediaItemDataHandler[*clienttypes.PlexConfig, *mediatypes.Album](c)
	registerClientUserMediaItemDataHandler[*clienttypes.SubsonicConfig, *mediatypes.Album](c)

	registerClientUserMediaItemDataHandler[*clienttypes.EmbyConfig, *mediatypes.Artist](c)
	registerClientUserMediaItemDataHandler[*clienttypes.JellyfinConfig, *mediatypes.Artist](c)
	registerClientUserMediaItemDataHandler[*clienttypes.PlexConfig, *mediatypes.Artist](c)
	registerClientUserMediaItemDataHandler[*clienttypes.SubsonicConfig, *mediatypes.Artist](c)

	registerClientUserMediaItemDataHandler[*clienttypes.EmbyConfig, *mediatypes.Collection](c)
	registerClientUserMediaItemDataHandler[*clienttypes.JellyfinConfig, *mediatypes.Collection](c)
	registerClientUserMediaItemDataHandler[*clienttypes.PlexConfig, *mediatypes.Collection](c)

	registerClientUserMediaItemDataHandler[*clienttypes.EmbyConfig, *mediatypes.Playlist](c)
	registerClientUserMediaItemDataHandler[*clienttypes.JellyfinConfig, *mediatypes.Playlist](c)
	registerClientUserMediaItemDataHandler[*clienttypes.PlexConfig, *mediatypes.Playlist](c)
	registerClientUserMediaItemDataHandler[*clienttypes.SubsonicConfig, *mediatypes.Playlist](c)
}

func registerCoreUserMediaItemDataHandler[T mediatypes.MediaData](c *container.Container) {
	container.RegisterFactory[handlers.CoreUserMediaItemDataHandler[T]](c, func(c *container.Container) handlers.CoreUserMediaItemDataHandler[T] {
		coreService := container.MustGet[services.CoreUserMediaItemDataService[T]](c)
		return handlers.NewCoreUserMediaItemDataHandler[T](coreService)
	})
}

func registerUserMediaItemDataHandler[T mediatypes.MediaData](c *container.Container) {
	container.RegisterFactory[handlers.UserMediaItemDataHandler[T]](c, func(c *container.Container) handlers.UserMediaItemDataHandler[T] {
		coreHandler := container.MustGet[handlers.CoreUserMediaItemDataHandler[T]](c)
		userService := container.MustGet[services.UserMediaItemDataService[T]](c)
		return handlers.NewUserMediaItemDataHandler[T](coreHandler, userService)
	})
}

func registerClientUserMediaItemDataHandler[T clienttypes.ClientMediaConfig, U mediatypes.MediaData](c *container.Container) {
	container.RegisterFactory[handlers.ClientUserMediaItemDataHandler[T, U]](c, func(c *container.Container) handlers.ClientUserMediaItemDataHandler[T, U] {
		userHandler := container.MustGet[handlers.UserMediaItemDataHandler[U]](c)
		clientService := container.MustGet[services.ClientUserMediaItemDataService[T, U]](c)
		return handlers.NewClientUserMediaItemDataHandler[T, U](userHandler, clientService)
	})
}
