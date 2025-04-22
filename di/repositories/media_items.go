// app/di/repositories/media_repositories.go
package repositories

import (
	"context"
	"gorm.io/gorm"
	mediatypes "suasor/client/media/types"
	"suasor/container"
	"suasor/repository"
	repobundle "suasor/repository/bundles"
)

func registerMediaItemRepository[T mediatypes.MediaData](c *container.Container, db *gorm.DB) {

	container.RegisterFactory[repository.MediaItemRepository[T]](c, func(c *container.Container) repository.MediaItemRepository[T] {
		return repository.NewMediaItemRepository[T](db)
	})

	container.RegisterFactory[repository.ClientMediaItemRepository[T]](c, func(c *container.Container) repository.ClientMediaItemRepository[T] {
		itemRepo := container.MustGet[repository.MediaItemRepository[T]](c)
		return repository.NewClientMediaItemRepository[T](db, itemRepo)
	})

	container.RegisterFactory[repository.ClientMediaItemRepository[T]](c, func(c *container.Container) repository.ClientMediaItemRepository[T] {
		coreRepo := container.MustGet[repository.MediaItemRepository[T]](c)
		return repository.NewClientMediaItemRepository[T](db, coreRepo)
	})

}

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
		movieRepo := container.MustGet[repository.MediaItemRepository[*mediatypes.Movie]](c)
		seriesRepo := container.MustGet[repository.MediaItemRepository[*mediatypes.Series]](c)
		seasonRepo := container.MustGet[repository.MediaItemRepository[*mediatypes.Season]](c)
		episodeRepo := container.MustGet[repository.MediaItemRepository[*mediatypes.Episode]](c)
		trackRepo := container.MustGet[repository.MediaItemRepository[*mediatypes.Track]](c)
		albumRepo := container.MustGet[repository.MediaItemRepository[*mediatypes.Album]](c)
		artistRepo := container.MustGet[repository.MediaItemRepository[*mediatypes.Artist]](c)
		collectionRepo := container.MustGet[repository.MediaItemRepository[*mediatypes.Collection]](c)
		playlistRepo := container.MustGet[repository.MediaItemRepository[*mediatypes.Playlist]](c)

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
