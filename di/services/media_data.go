package services

import (
	"context"
	"suasor/clients"
	mediatypes "suasor/clients/media/types"
	clienttypes "suasor/clients/types"
	"suasor/di/container"
	"suasor/repository"
	"suasor/services"
	"suasor/utils/logger"
)

func registerMediaDataServices(ctx context.Context, c *container.Container) {
	log := logger.LoggerFromContext(ctx)

	log.Info().Msg("Registering media Data services")
	registerMediaDataService[*mediatypes.Movie](c)
	registerMediaDataService[*mediatypes.Series](c)
	registerMediaDataService[*mediatypes.Season](c)
	registerMediaDataService[*mediatypes.Episode](c)
	registerMediaDataService[*mediatypes.Track](c)
	registerMediaDataService[*mediatypes.Album](c)
	registerMediaDataService[*mediatypes.Artist](c)
	registerMediaDataService[*mediatypes.Collection](c)
	registerMediaDataService[*mediatypes.Playlist](c)

	log.Info().Msg("Registering client media item services")
	registerClientMediaDataService[*clienttypes.EmbyConfig, *mediatypes.Movie](c)
	registerClientMediaDataService[*clienttypes.JellyfinConfig, *mediatypes.Movie](c)
	registerClientMediaDataService[*clienttypes.PlexConfig, *mediatypes.Movie](c)

	registerClientMediaDataService[*clienttypes.EmbyConfig, *mediatypes.Series](c)
	registerClientMediaDataService[*clienttypes.JellyfinConfig, *mediatypes.Series](c)
	registerClientMediaDataService[*clienttypes.PlexConfig, *mediatypes.Series](c)

	registerClientMediaDataService[*clienttypes.EmbyConfig, *mediatypes.Season](c)
	registerClientMediaDataService[*clienttypes.JellyfinConfig, *mediatypes.Season](c)
	registerClientMediaDataService[*clienttypes.PlexConfig, *mediatypes.Season](c)

	registerClientMediaDataService[*clienttypes.EmbyConfig, *mediatypes.Episode](c)
	registerClientMediaDataService[*clienttypes.JellyfinConfig, *mediatypes.Episode](c)
	registerClientMediaDataService[*clienttypes.PlexConfig, *mediatypes.Episode](c)

	registerClientMediaDataService[*clienttypes.EmbyConfig, *mediatypes.Track](c)
	registerClientMediaDataService[*clienttypes.JellyfinConfig, *mediatypes.Track](c)
	registerClientMediaDataService[*clienttypes.PlexConfig, *mediatypes.Track](c)
	registerClientMediaDataService[*clienttypes.SubsonicConfig, *mediatypes.Track](c)

	registerClientMediaDataService[*clienttypes.EmbyConfig, *mediatypes.Album](c)
	registerClientMediaDataService[*clienttypes.JellyfinConfig, *mediatypes.Album](c)
	registerClientMediaDataService[*clienttypes.PlexConfig, *mediatypes.Album](c)
	registerClientMediaDataService[*clienttypes.SubsonicConfig, *mediatypes.Album](c)

	registerClientMediaDataService[*clienttypes.EmbyConfig, *mediatypes.Artist](c)
	registerClientMediaDataService[*clienttypes.JellyfinConfig, *mediatypes.Artist](c)
	registerClientMediaDataService[*clienttypes.PlexConfig, *mediatypes.Artist](c)
	registerClientMediaDataService[*clienttypes.SubsonicConfig, *mediatypes.Artist](c)

	registerClientMediaDataService[*clienttypes.EmbyConfig, *mediatypes.Collection](c)
	registerClientMediaDataService[*clienttypes.JellyfinConfig, *mediatypes.Collection](c)
	registerClientMediaDataService[*clienttypes.PlexConfig, *mediatypes.Collection](c)

	registerClientMediaDataService[*clienttypes.EmbyConfig, *mediatypes.Playlist](c)
	registerClientMediaDataService[*clienttypes.JellyfinConfig, *mediatypes.Playlist](c)
	registerClientMediaDataService[*clienttypes.PlexConfig, *mediatypes.Playlist](c)
	registerClientMediaDataService[*clienttypes.SubsonicConfig, *mediatypes.Playlist](c)
}

func registerMediaDataService[T mediatypes.MediaData](c *container.Container) {
	container.RegisterFactory[services.CoreUserMediaItemDataService[T]](c, func(c *container.Container) services.CoreUserMediaItemDataService[T] {
		itemService := container.MustGet[services.CoreMediaItemService[T]](c)
		repo := container.MustGet[repository.CoreUserMediaItemDataRepository[T]](c)
		return services.NewCoreUserMediaItemDataService[T](itemService, repo)
	})
	container.RegisterFactory[services.UserMediaItemDataService[T]](c, func(c *container.Container) services.UserMediaItemDataService[T] {
		coreService := container.MustGet[services.CoreUserMediaItemDataService[T]](c)
		repo := container.MustGet[repository.UserMediaItemDataRepository[T]](c)
		return services.NewUserMediaItemDataService[T](coreService, repo)
	})
}

func registerClientMediaDataService[T clienttypes.ClientMediaConfig, U mediatypes.MediaData](c *container.Container) {
	container.RegisterFactory[services.ClientUserMediaItemDataService[T, U]](c, func(c *container.Container) services.ClientUserMediaItemDataService[T, U] {
		userService := container.MustGet[services.UserMediaItemDataService[U]](c)
		dataRepo := container.MustGet[repository.ClientUserMediaItemDataRepository[U]](c)
		clientRepo := container.MustGet[repository.ClientRepository[T]](c)
		clientFactory := container.MustGet[*clients.ClientProviderFactoryService](c)
		return services.NewClientUserMediaItemDataService(userService, dataRepo, clientRepo, clientFactory)
	})
}
