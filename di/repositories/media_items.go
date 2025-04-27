// app/di/repositories/media_repositories.go
package repositories

import (
	"gorm.io/gorm"
	mediatypes "suasor/clients/media/types"
	"suasor/di/container"
	"suasor/repository"
	repobundle "suasor/repository/bundles"
)

// Register core repositories for all media types
func registerMediaItemRepositories(c *container.Container) {
	db := container.MustGet[*gorm.DB](c)

	registerMediaItemRepository[*mediatypes.Movie](c, db)
	registerMediaItemRepository[*mediatypes.Series](c, db)
	registerMediaItemRepository[*mediatypes.Season](c, db)
	registerMediaItemRepository[*mediatypes.Episode](c, db)
	registerMediaItemRepository[*mediatypes.Track](c, db)
	registerMediaItemRepository[*mediatypes.Album](c, db)
	registerMediaItemRepository[*mediatypes.Artist](c, db)
	registerMediaItemRepository[*mediatypes.Collection](c, db)
	registerMediaItemRepository[*mediatypes.Playlist](c, db)

	// Register CoreMediaItemRepositories container
	container.RegisterFactory[repobundle.CoreMediaItemRepositories](c, func(c *container.Container) repobundle.CoreMediaItemRepositories {
		movieRepo := container.MustGet[repository.CoreMediaItemRepository[*mediatypes.Movie]](c)
		seriesRepo := container.MustGet[repository.CoreMediaItemRepository[*mediatypes.Series]](c)
		seasonRepo := container.MustGet[repository.CoreMediaItemRepository[*mediatypes.Season]](c)
		episodeRepo := container.MustGet[repository.CoreMediaItemRepository[*mediatypes.Episode]](c)
		trackRepo := container.MustGet[repository.CoreMediaItemRepository[*mediatypes.Track]](c)
		albumRepo := container.MustGet[repository.CoreMediaItemRepository[*mediatypes.Album]](c)
		artistRepo := container.MustGet[repository.CoreMediaItemRepository[*mediatypes.Artist]](c)
		collectionRepo := container.MustGet[repository.CoreMediaItemRepository[*mediatypes.Collection]](c)
		playlistRepo := container.MustGet[repository.CoreMediaItemRepository[*mediatypes.Playlist]](c)

		return repobundle.NewCoreMediaItemRepositories(
			movieRepo, seriesRepo, seasonRepo, episodeRepo,
			trackRepo, albumRepo, artistRepo,
			collectionRepo, playlistRepo,
		)
	})

	// Register UserMediaItemRepositories container
	container.RegisterFactory[repobundle.UserMediaItemRepositories](c, func(c *container.Container) repobundle.UserMediaItemRepositories {
		movieRepo := container.MustGet[repository.UserMediaItemRepository[*mediatypes.Movie]](c)
		seriesRepo := container.MustGet[repository.UserMediaItemRepository[*mediatypes.Series]](c)
		seasonRepo := container.MustGet[repository.UserMediaItemRepository[*mediatypes.Season]](c)
		episodeRepo := container.MustGet[repository.UserMediaItemRepository[*mediatypes.Episode]](c)
		trackRepo := container.MustGet[repository.UserMediaItemRepository[*mediatypes.Track]](c)
		albumRepo := container.MustGet[repository.UserMediaItemRepository[*mediatypes.Album]](c)
		artistRepo := container.MustGet[repository.UserMediaItemRepository[*mediatypes.Artist]](c)
		collectionRepo := container.MustGet[repository.UserMediaItemRepository[*mediatypes.Collection]](c)
		playlistRepo := container.MustGet[repository.UserMediaItemRepository[*mediatypes.Playlist]](c)

		return repobundle.NewUserMediaItemRepositories(
			movieRepo, seriesRepo, seasonRepo, episodeRepo,
			trackRepo, albumRepo, artistRepo,
			collectionRepo, playlistRepo,
		)
	})

	// Register ClientMediaItemRepositories container
	container.RegisterFactory[repobundle.ClientMediaItemRepositories](c, func(c *container.Container) repobundle.ClientMediaItemRepositories {
		movieRepo := container.MustGet[repository.ClientMediaItemRepository[*mediatypes.Movie]](c)
		seriesRepo := container.MustGet[repository.ClientMediaItemRepository[*mediatypes.Series]](c)
		seasonRepo := container.MustGet[repository.ClientMediaItemRepository[*mediatypes.Season]](c)
		episodeRepo := container.MustGet[repository.ClientMediaItemRepository[*mediatypes.Episode]](c)
		trackRepo := container.MustGet[repository.ClientMediaItemRepository[*mediatypes.Track]](c)
		albumRepo := container.MustGet[repository.ClientMediaItemRepository[*mediatypes.Album]](c)
		artistRepo := container.MustGet[repository.ClientMediaItemRepository[*mediatypes.Artist]](c)
		collectionRepo := container.MustGet[repository.ClientMediaItemRepository[*mediatypes.Collection]](c)
		playlistRepo := container.MustGet[repository.ClientMediaItemRepository[*mediatypes.Playlist]](c)

		return repobundle.NewClientMediaItemRepositories(
			movieRepo, seriesRepo, seasonRepo, episodeRepo,
			trackRepo, albumRepo, artistRepo,
			collectionRepo, playlistRepo,
		)
	})
}

// Registers all 3 types of media item repositories
func registerMediaItemRepository[T mediatypes.MediaData](c *container.Container, db *gorm.DB) {

	container.RegisterFactory[repository.CoreMediaItemRepository[T]](c, func(c *container.Container) repository.CoreMediaItemRepository[T] {
		return repository.NewMediaItemRepository[T](db)
	})

	container.RegisterFactory[repository.UserMediaItemRepository[T]](c, func(c *container.Container) repository.UserMediaItemRepository[T] {
		itemRepo := container.MustGet[repository.CoreMediaItemRepository[T]](c)
		return repository.NewUserMediaItemRepository[T](db, itemRepo)
	})

	container.RegisterFactory[repository.ClientMediaItemRepository[T]](c, func(c *container.Container) repository.ClientMediaItemRepository[T] {
		userRepo := container.MustGet[repository.UserMediaItemRepository[T]](c)
		return repository.NewClientMediaItemRepository[T](db, userRepo)
	})

}
