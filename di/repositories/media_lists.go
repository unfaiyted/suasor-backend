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
func registerMediaListRepositories(c *container.Container) {
	db := container.MustGet[*gorm.DB](c)

}	

	// Register CoreMediaItemRepositories container
	// container.RegisterFactory[repobundle.CoreMediaListRepositories](c, func(c *container.Container) repobundle.CoreMediaItemRepositories {
	//
	// 	return repobundle.NewCoreMediaItemRepositories(
	// 		movieRepo, seriesRepo, seasonRepo, episodeRepo,
	// 		trackRepo, albumRepo, artistRepo,
	// 		collectionRepo, playlistRepo,
	// 	)
	// })

	// Register ClientMediaItemRepositories container
	// container.RegisterFactory[repobundle.ClientMediaItemRepositories](c, func(c *container.Container) repobundle.ClientMediaItemRepositories {
	// 	movieRepo := container.MustGet[repository.ClientMediaItemRepository[*mediatypes.Movie]](c)
	// 	seriesRepo := container.MustGet[repository.ClientMediaItemRepository[*mediatypes.Series]](c)
	// 	seasonRepo := container.MustGet[repository.ClientMediaItemRepository[*mediatypes.Season]](c)
	// 	episodeRepo := container.MustGet[repository.ClientMediaItemRepository[*mediatypes.Episode]](c)
	// 	trackRepo := container.MustGet[repository.ClientMediaItemRepository[*mediatypes.Track]](c)
	// 	albumRepo := container.MustGet[repository.ClientMediaItemRepository[*mediatypes.Album]](c)
	// 	artistRepo := container.MustGet[repository.ClientMediaItemRepository[*mediatypes.Artist]](c)
	// 	collectionRepo := container.MustGet[repository.ClientMediaItemRepository[*mediatypes.Collection]](c)
	// 	playlistRepo := container.MustGet[repository.ClientMediaItemRepository[*mediatypes.Playlist]](c)
	//
	// 	return repobundle.NewClientMediaItemRepositories(
	// 		movieRepo, seriesRepo, seasonRepo, episodeRepo,
	// 		trackRepo, albumRepo, artistRepo,
	// 		collectionRepo, playlistRepo,
	// 	)
	// })
}

// Reigisters all 3 types of media item repositories
func registerMediaListRepository[T mediatypes.ListData](c *container.Container, db *gorm.DB) {

	container.RegisterFactory[repository.MediaListRepository[T]](c, func(c *container.Container) repository.MediaItemRepository[T] {
		return repository.NewMediaItemRepository[T](db)
	})

	container.RegisterFactory[repository.UserMediaItemRepository[T]](c, func(c *container.Container) repository.UserMediaItemRepository[T] {
		itemRepo := container.MustGet[repository.MediaItemRepository[T]](c)
		return repository.NewUserMediaItemRepository[T](db, itemRepo)
	})

	container.RegisterFactory[repository.ClientMediaItemRepository[T]](c, func(c *container.Container) repository.ClientMediaItemRepository[T] {
		coreRepo := container.MustGet[repository.MediaItemRepository[T]](c)
		return repository.NewClientMediaItemRepository[T](db, coreRepo)
	})

}
