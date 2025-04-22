// app/di/repositories/media_repositories.go
package repositories

import (
	"context"
	"gorm.io/gorm"
	"suasor/app/container"
	apprepository "suasor/app/repository"
	mediatypes "suasor/client/media/types"
	"suasor/repository"
)

// RegisterMediaRepositories registers all media-related repositories for the three-pronged architecture
func RegisterMediaRepositories(ctx context.Context, c *container.Container) {
	// Core Media Item Repositories
	registerCoreMediaItemRepositories(c)

	// User Media Item Repositories
	registerUserMediaItemRepositories(c)

	// Client Media Item Repositories
	registerClientMediaItemRepositories(c)

	// Media Item Data Repositories
	registerMediaItemDataRepositories(c)
}

// Register core repositories for all media types
func registerCoreMediaItemRepositories(c *container.Container) {
	db := container.MustGet[*gorm.DB](c)

	// Register individual core repositories
	container.RegisterFactory[repository.MediaItemRepository[*mediatypes.Movie]](c, func(c *container.Container) repository.MediaItemRepository[*mediatypes.Movie] {
		return repository.NewMediaItemRepository[*mediatypes.Movie](db)
	})

	container.RegisterFactory[repository.MediaItemRepository[*mediatypes.Series]](c, func(c *container.Container) repository.MediaItemRepository[*mediatypes.Series] {
		return repository.NewMediaItemRepository[*mediatypes.Series](db)
	})

	container.RegisterFactory[repository.MediaItemRepository[*mediatypes.Season]](c, func(c *container.Container) repository.MediaItemRepository[*mediatypes.Season] {
		return repository.NewMediaItemRepository[*mediatypes.Season](db)
	})

	container.RegisterFactory[repository.MediaItemRepository[*mediatypes.Episode]](c, func(c *container.Container) repository.MediaItemRepository[*mediatypes.Episode] {
		return repository.NewMediaItemRepository[*mediatypes.Episode](db)
	})

	container.RegisterFactory[repository.MediaItemRepository[*mediatypes.Track]](c, func(c *container.Container) repository.MediaItemRepository[*mediatypes.Track] {
		return repository.NewMediaItemRepository[*mediatypes.Track](db)
	})

	container.RegisterFactory[repository.MediaItemRepository[*mediatypes.Album]](c, func(c *container.Container) repository.MediaItemRepository[*mediatypes.Album] {
		return repository.NewMediaItemRepository[*mediatypes.Album](db)
	})

	container.RegisterFactory[repository.MediaItemRepository[*mediatypes.Artist]](c, func(c *container.Container) repository.MediaItemRepository[*mediatypes.Artist] {
		return repository.NewMediaItemRepository[*mediatypes.Artist](db)
	})

	container.RegisterFactory[repository.MediaItemRepository[*mediatypes.Collection]](c, func(c *container.Container) repository.MediaItemRepository[*mediatypes.Collection] {
		return repository.NewMediaItemRepository[*mediatypes.Collection](db)
	})

	container.RegisterFactory[repository.MediaItemRepository[*mediatypes.Playlist]](c, func(c *container.Container) repository.MediaItemRepository[*mediatypes.Playlist] {
		return repository.NewMediaItemRepository[*mediatypes.Playlist](db)
	})

	// Register CoreMediaItemRepositories container
	container.RegisterFactory[apprepository.CoreMediaItemRepositories](c, func(c *container.Container) apprepository.CoreMediaItemRepositories {
		movieRepo := container.MustGet[repository.MediaItemRepository[*mediatypes.Movie]](c)
		seriesRepo := container.MustGet[repository.MediaItemRepository[*mediatypes.Series]](c)
		seasonRepo := container.MustGet[repository.MediaItemRepository[*mediatypes.Season]](c)
		episodeRepo := container.MustGet[repository.MediaItemRepository[*mediatypes.Episode]](c)
		trackRepo := container.MustGet[repository.MediaItemRepository[*mediatypes.Track]](c)
		albumRepo := container.MustGet[repository.MediaItemRepository[*mediatypes.Album]](c)
		artistRepo := container.MustGet[repository.MediaItemRepository[*mediatypes.Artist]](c)
		collectionRepo := container.MustGet[repository.MediaItemRepository[*mediatypes.Collection]](c)
		playlistRepo := container.MustGet[repository.MediaItemRepository[*mediatypes.Playlist]](c)

		return apprepository.NewCoreMediaItemRepositories(
			movieRepo, seriesRepo, seasonRepo, episodeRepo,
			trackRepo, albumRepo, artistRepo,
			collectionRepo, playlistRepo,
		)
	})
}

// Register user repositories for all media types
func registerUserMediaItemRepositories(c *container.Container) {
	db := container.MustGet[*gorm.DB](c)

	// Register individual user repositories
	container.RegisterFactory[repository.UserMediaItemRepository[*mediatypes.Movie]](c, func(c *container.Container) repository.UserMediaItemRepository[*mediatypes.Movie] {
		return repository.NewUserMediaItemRepository[*mediatypes.Movie](db)
	})

	container.RegisterFactory[repository.UserMediaItemRepository[*mediatypes.Series]](c, func(c *container.Container) repository.UserMediaItemRepository[*mediatypes.Series] {
		return repository.NewUserMediaItemRepository[*mediatypes.Series](db)
	})

	container.RegisterFactory[repository.UserMediaItemRepository[*mediatypes.Season]](c, func(c *container.Container) repository.UserMediaItemRepository[*mediatypes.Season] {
		return repository.NewUserMediaItemRepository[*mediatypes.Season](db)
	})

	container.RegisterFactory[repository.UserMediaItemRepository[*mediatypes.Episode]](c, func(c *container.Container) repository.UserMediaItemRepository[*mediatypes.Episode] {
		return repository.NewUserMediaItemRepository[*mediatypes.Episode](db)
	})

	container.RegisterFactory[repository.UserMediaItemRepository[*mediatypes.Track]](c, func(c *container.Container) repository.UserMediaItemRepository[*mediatypes.Track] {
		return repository.NewUserMediaItemRepository[*mediatypes.Track](db)
	})

	container.RegisterFactory[repository.UserMediaItemRepository[*mediatypes.Album]](c, func(c *container.Container) repository.UserMediaItemRepository[*mediatypes.Album] {
		return repository.NewUserMediaItemRepository[*mediatypes.Album](db)
	})

	container.RegisterFactory[repository.UserMediaItemRepository[*mediatypes.Artist]](c, func(c *container.Container) repository.UserMediaItemRepository[*mediatypes.Artist] {
		return repository.NewUserMediaItemRepository[*mediatypes.Artist](db)
	})

	container.RegisterFactory[repository.UserMediaItemRepository[*mediatypes.Collection]](c, func(c *container.Container) repository.UserMediaItemRepository[*mediatypes.Collection] {
		return repository.NewUserMediaItemRepository[*mediatypes.Collection](db)
	})

	container.RegisterFactory[repository.UserMediaItemRepository[*mediatypes.Playlist]](c, func(c *container.Container) repository.UserMediaItemRepository[*mediatypes.Playlist] {
		return repository.NewUserMediaItemRepository[*mediatypes.Playlist](db)
	})

	// Register UserMediaItemRepositories container
	container.RegisterFactory[apprepository.UserMediaItemRepositories](c, func(c *container.Container) apprepository.UserMediaItemRepositories {
		movieRepo := container.MustGet[repository.UserMediaItemRepository[*mediatypes.Movie]](c)
		seriesRepo := container.MustGet[repository.UserMediaItemRepository[*mediatypes.Series]](c)
		seasonRepo := container.MustGet[repository.UserMediaItemRepository[*mediatypes.Season]](c)
		episodeRepo := container.MustGet[repository.UserMediaItemRepository[*mediatypes.Episode]](c)
		trackRepo := container.MustGet[repository.UserMediaItemRepository[*mediatypes.Track]](c)
		albumRepo := container.MustGet[repository.UserMediaItemRepository[*mediatypes.Album]](c)
		artistRepo := container.MustGet[repository.UserMediaItemRepository[*mediatypes.Artist]](c)
		collectionRepo := container.MustGet[repository.UserMediaItemRepository[*mediatypes.Collection]](c)
		playlistRepo := container.MustGet[repository.UserMediaItemRepository[*mediatypes.Playlist]](c)

		return apprepository.NewUserMediaItemRepositories(
			movieRepo, seriesRepo, seasonRepo, episodeRepo,
			trackRepo, albumRepo, artistRepo,
			collectionRepo, playlistRepo,
		)
	})
}

// Register client repositories for all media types
func registerClientMediaItemRepositories(c *container.Container) {
	db := container.MustGet[*gorm.DB](c)

	// Register individual client repositories
	container.RegisterFactory[repository.ClientMediaItemRepository[*mediatypes.Movie]](c, func(c *container.Container) repository.ClientMediaItemRepository[*mediatypes.Movie] {
		coreRepo := container.MustGet[repository.MediaItemRepository[*mediatypes.Movie]](c)
		return repository.NewClientMediaItemRepository[*mediatypes.Movie](db, coreRepo)
	})

	container.RegisterFactory[repository.ClientMediaItemRepository[*mediatypes.Series]](c, func(c *container.Container) repository.ClientMediaItemRepository[*mediatypes.Series] {
		coreRepo := container.MustGet[repository.MediaItemRepository[*mediatypes.Series]](c)
		return repository.NewClientMediaItemRepository[*mediatypes.Series](db, coreRepo)
	})

	container.RegisterFactory[repository.ClientMediaItemRepository[*mediatypes.Season]](c, func(c *container.Container) repository.ClientMediaItemRepository[*mediatypes.Season] {
		coreRepo := container.MustGet[repository.MediaItemRepository[*mediatypes.Season]](c)
		return repository.NewClientMediaItemRepository[*mediatypes.Season](db, coreRepo)
	})

	container.RegisterFactory[repository.ClientMediaItemRepository[*mediatypes.Episode]](c, func(c *container.Container) repository.ClientMediaItemRepository[*mediatypes.Episode] {
		coreRepo := container.MustGet[repository.MediaItemRepository[*mediatypes.Episode]](c)
		return repository.NewClientMediaItemRepository[*mediatypes.Episode](db, coreRepo)
	})

	container.RegisterFactory[repository.ClientMediaItemRepository[*mediatypes.Track]](c, func(c *container.Container) repository.ClientMediaItemRepository[*mediatypes.Track] {
		coreRepo := container.MustGet[repository.MediaItemRepository[*mediatypes.Track]](c)
		return repository.NewClientMediaItemRepository[*mediatypes.Track](db, coreRepo)
	})

	container.RegisterFactory[repository.ClientMediaItemRepository[*mediatypes.Album]](c, func(c *container.Container) repository.ClientMediaItemRepository[*mediatypes.Album] {
		coreRepo := container.MustGet[repository.MediaItemRepository[*mediatypes.Album]](c)
		return repository.NewClientMediaItemRepository[*mediatypes.Album](db, coreRepo)
	})

	container.RegisterFactory[repository.ClientMediaItemRepository[*mediatypes.Artist]](c, func(c *container.Container) repository.ClientMediaItemRepository[*mediatypes.Artist] {
		coreRepo := container.MustGet[repository.MediaItemRepository[*mediatypes.Artist]](c)
		return repository.NewClientMediaItemRepository[*mediatypes.Artist](db, coreRepo)
	})

	container.RegisterFactory[repository.ClientMediaItemRepository[*mediatypes.Collection]](c, func(c *container.Container) repository.ClientMediaItemRepository[*mediatypes.Collection] {
		coreRepo := container.MustGet[repository.MediaItemRepository[*mediatypes.Collection]](c)
		return repository.NewClientMediaItemRepository[*mediatypes.Collection](db, coreRepo)
	})

	container.RegisterFactory[repository.ClientMediaItemRepository[*mediatypes.Playlist]](c, func(c *container.Container) repository.ClientMediaItemRepository[*mediatypes.Playlist] {
		coreRepo := container.MustGet[repository.MediaItemRepository[*mediatypes.Playlist]](c)
		return repository.NewClientMediaItemRepository[*mediatypes.Playlist](db, coreRepo)
	})

	// Register ClientMediaItemRepositories container
	container.RegisterFactory[apprepository.ClientMediaItemRepositories](c, func(c *container.Container) apprepository.ClientMediaItemRepositories {
		movieRepo := container.MustGet[repository.ClientMediaItemRepository[*mediatypes.Movie]](c)
		seriesRepo := container.MustGet[repository.ClientMediaItemRepository[*mediatypes.Series]](c)
		seasonRepo := container.MustGet[repository.ClientMediaItemRepository[*mediatypes.Season]](c)
		episodeRepo := container.MustGet[repository.ClientMediaItemRepository[*mediatypes.Episode]](c)
		trackRepo := container.MustGet[repository.ClientMediaItemRepository[*mediatypes.Track]](c)
		albumRepo := container.MustGet[repository.ClientMediaItemRepository[*mediatypes.Album]](c)
		artistRepo := container.MustGet[repository.ClientMediaItemRepository[*mediatypes.Artist]](c)
		collectionRepo := container.MustGet[repository.ClientMediaItemRepository[*mediatypes.Collection]](c)
		playlistRepo := container.MustGet[repository.ClientMediaItemRepository[*mediatypes.Playlist]](c)

		return apprepository.NewClientMediaItemRepositories(
			movieRepo, seriesRepo, seasonRepo, episodeRepo,
			trackRepo, albumRepo, artistRepo,
			collectionRepo, playlistRepo,
		)
	})
}

// Register media item data repositories
func registerMediaItemDataRepositories(c *container.Container) {
	db := container.MustGet[*gorm.DB](c)

	// Register core user media item data repositories
	container.RegisterFactory[repository.CoreUserMediaItemDataRepository[*mediatypes.Movie]](c, func(c *container.Container) repository.CoreUserMediaItemDataRepository[*mediatypes.Movie] {
		return repository.NewCoreUserMediaItemDataRepository[*mediatypes.Movie](db)
	})

	container.RegisterFactory[repository.CoreUserMediaItemDataRepository[*mediatypes.Series]](c, func(c *container.Container) repository.CoreUserMediaItemDataRepository[*mediatypes.Series] {
		return repository.NewCoreUserMediaItemDataRepository[*mediatypes.Series](db)
	})

	container.RegisterFactory[repository.CoreUserMediaItemDataRepository[*mediatypes.Episode]](c, func(c *container.Container) repository.CoreUserMediaItemDataRepository[*mediatypes.Episode] {
		return repository.NewCoreUserMediaItemDataRepository[*mediatypes.Episode](db)
	})

	container.RegisterFactory[repository.CoreUserMediaItemDataRepository[*mediatypes.Track]](c, func(c *container.Container) repository.CoreUserMediaItemDataRepository[*mediatypes.Track] {
		return repository.NewCoreUserMediaItemDataRepository[*mediatypes.Track](db)
	})

	container.RegisterFactory[repository.CoreUserMediaItemDataRepository[*mediatypes.Album]](c, func(c *container.Container) repository.CoreUserMediaItemDataRepository[*mediatypes.Album] {
		return repository.NewCoreUserMediaItemDataRepository[*mediatypes.Album](db)
	})

	container.RegisterFactory[repository.CoreUserMediaItemDataRepository[*mediatypes.Artist]](c, func(c *container.Container) repository.CoreUserMediaItemDataRepository[*mediatypes.Artist] {
		return repository.NewCoreUserMediaItemDataRepository[*mediatypes.Artist](db)
	})

	container.RegisterFactory[repository.CoreUserMediaItemDataRepository[*mediatypes.Collection]](c, func(c *container.Container) repository.CoreUserMediaItemDataRepository[*mediatypes.Collection] {
		return repository.NewCoreUserMediaItemDataRepository[*mediatypes.Collection](db)
	})

	container.RegisterFactory[repository.CoreUserMediaItemDataRepository[*mediatypes.Playlist]](c, func(c *container.Container) repository.CoreUserMediaItemDataRepository[*mediatypes.Playlist] {
		return repository.NewCoreUserMediaItemDataRepository[*mediatypes.Playlist](db)
	})

	// Register user media item data repositories
	container.RegisterFactory[repository.UserMediaItemDataRepository[*mediatypes.Movie]](c, func(c *container.Container) repository.UserMediaItemDataRepository[*mediatypes.Movie] {
		coreRepo := container.MustGet[repository.CoreUserMediaItemDataRepository[*mediatypes.Movie]](c)
		return repository.NewUserMediaItemDataRepository[*mediatypes.Movie](db, coreRepo)
	})

	container.RegisterFactory[repository.UserMediaItemDataRepository[*mediatypes.Series]](c, func(c *container.Container) repository.UserMediaItemDataRepository[*mediatypes.Series] {
		coreRepo := container.MustGet[repository.CoreUserMediaItemDataRepository[*mediatypes.Series]](c)
		return repository.NewUserMediaItemDataRepository[*mediatypes.Series](db, coreRepo)
	})

	container.RegisterFactory[repository.UserMediaItemDataRepository[*mediatypes.Episode]](c, func(c *container.Container) repository.UserMediaItemDataRepository[*mediatypes.Episode] {
		coreRepo := container.MustGet[repository.CoreUserMediaItemDataRepository[*mediatypes.Episode]](c)
		return repository.NewUserMediaItemDataRepository[*mediatypes.Episode](db, coreRepo)
	})

	container.RegisterFactory[repository.UserMediaItemDataRepository[*mediatypes.Track]](c, func(c *container.Container) repository.UserMediaItemDataRepository[*mediatypes.Track] {
		coreRepo := container.MustGet[repository.CoreUserMediaItemDataRepository[*mediatypes.Track]](c)
		return repository.NewUserMediaItemDataRepository[*mediatypes.Track](db, coreRepo)
	})

	container.RegisterFactory[repository.UserMediaItemDataRepository[*mediatypes.Album]](c, func(c *container.Container) repository.UserMediaItemDataRepository[*mediatypes.Album] {
		coreRepo := container.MustGet[repository.CoreUserMediaItemDataRepository[*mediatypes.Album]](c)
		return repository.NewUserMediaItemDataRepository[*mediatypes.Album](db, coreRepo)
	})

	container.RegisterFactory[repository.UserMediaItemDataRepository[*mediatypes.Artist]](c, func(c *container.Container) repository.UserMediaItemDataRepository[*mediatypes.Artist] {
		coreRepo := container.MustGet[repository.CoreUserMediaItemDataRepository[*mediatypes.Artist]](c)
		return repository.NewUserMediaItemDataRepository[*mediatypes.Artist](db, coreRepo)
	})

	container.RegisterFactory[repository.UserMediaItemDataRepository[*mediatypes.Collection]](c, func(c *container.Container) repository.UserMediaItemDataRepository[*mediatypes.Collection] {
		coreRepo := container.MustGet[repository.CoreUserMediaItemDataRepository[*mediatypes.Collection]](c)
		return repository.NewUserMediaItemDataRepository[*mediatypes.Collection](db, coreRepo)
	})

	container.RegisterFactory[repository.UserMediaItemDataRepository[*mediatypes.Playlist]](c, func(c *container.Container) repository.UserMediaItemDataRepository[*mediatypes.Playlist] {
		coreRepo := container.MustGet[repository.CoreUserMediaItemDataRepository[*mediatypes.Playlist]](c)
		return repository.NewUserMediaItemDataRepository[*mediatypes.Playlist](db, coreRepo)
	})

	// Register CoreUserMediaItemDataRepositories container
	container.RegisterFactory[apprepository.CoreUserMediaItemDataRepositories](c, func(c *container.Container) apprepository.CoreUserMediaItemDataRepositories {
		movieRepo := container.MustGet[repository.CoreUserMediaItemDataRepository[*mediatypes.Movie]](c)
		seriesRepo := container.MustGet[repository.CoreUserMediaItemDataRepository[*mediatypes.Series]](c)
		episodeRepo := container.MustGet[repository.CoreUserMediaItemDataRepository[*mediatypes.Episode]](c)
		trackRepo := container.MustGet[repository.CoreUserMediaItemDataRepository[*mediatypes.Track]](c)
		albumRepo := container.MustGet[repository.CoreUserMediaItemDataRepository[*mediatypes.Album]](c)
		artistRepo := container.MustGet[repository.CoreUserMediaItemDataRepository[*mediatypes.Artist]](c)
		collectionRepo := container.MustGet[repository.CoreUserMediaItemDataRepository[*mediatypes.Collection]](c)
		playlistRepo := container.MustGet[repository.CoreUserMediaItemDataRepository[*mediatypes.Playlist]](c)

		return apprepository.NewCoreUserMediaItemDataRepositories(
			movieRepo, seriesRepo, episodeRepo,
			trackRepo, albumRepo, artistRepo,
			collectionRepo, playlistRepo,
		)
	})

	// Register UserMediaDataRepositories container
	container.RegisterFactory[apprepository.UserMediaDataRepositories](c, func(c *container.Container) apprepository.UserMediaDataRepositories {
		movieRepo := container.MustGet[repository.UserMediaItemDataRepository[*mediatypes.Movie]](c)
		seriesRepo := container.MustGet[repository.UserMediaItemDataRepository[*mediatypes.Series]](c)
		episodeRepo := container.MustGet[repository.UserMediaItemDataRepository[*mediatypes.Episode]](c)
		trackRepo := container.MustGet[repository.UserMediaItemDataRepository[*mediatypes.Track]](c)
		albumRepo := container.MustGet[repository.UserMediaItemDataRepository[*mediatypes.Album]](c)
		artistRepo := container.MustGet[repository.UserMediaItemDataRepository[*mediatypes.Artist]](c)
		collectionRepo := container.MustGet[repository.UserMediaItemDataRepository[*mediatypes.Collection]](c)
		playlistRepo := container.MustGet[repository.UserMediaItemDataRepository[*mediatypes.Playlist]](c)

		return apprepository.NewUserMediaDataRepositories(
			movieRepo, seriesRepo, episodeRepo,
			trackRepo, albumRepo, artistRepo,
			collectionRepo, playlistRepo,
		)
	})
}