package repositories

import (
	"gorm.io/gorm"
	mediatypes "suasor/clients/media/types"
	"suasor/di/container"
	"suasor/repository"
	repobundles "suasor/repository/bundles"
)

// Register media item data repositories
func registerMediaItemDataRepositories(c *container.Container) {
	db := container.MustGet[*gorm.DB](c)

	registerDataRepository[*mediatypes.Movie](c, db)
	registerDataRepository[*mediatypes.Series](c, db)
	registerDataRepository[*mediatypes.Season](c, db)
	registerDataRepository[*mediatypes.Episode](c, db)
	registerDataRepository[*mediatypes.Track](c, db)
	registerDataRepository[*mediatypes.Album](c, db)
	registerDataRepository[*mediatypes.Artist](c, db)
	registerDataRepository[*mediatypes.Collection](c, db)
	registerDataRepository[*mediatypes.Playlist](c, db)

	// TODO: ClientDataRepositories

	container.RegisterFactory[repobundles.UserMediaDataRepositories](c, func(c *container.Container) repobundles.UserMediaDataRepositories {
		return repobundles.NewUserMediaDataRepositories(
			container.MustGet[repository.UserMediaItemDataRepository[*mediatypes.Movie]](c),
			container.MustGet[repository.UserMediaItemDataRepository[*mediatypes.Series]](c),
			container.MustGet[repository.UserMediaItemDataRepository[*mediatypes.Season]](c),
			container.MustGet[repository.UserMediaItemDataRepository[*mediatypes.Episode]](c),
			container.MustGet[repository.UserMediaItemDataRepository[*mediatypes.Track]](c),
			container.MustGet[repository.UserMediaItemDataRepository[*mediatypes.Album]](c),
			container.MustGet[repository.UserMediaItemDataRepository[*mediatypes.Artist]](c),
			container.MustGet[repository.UserMediaItemDataRepository[*mediatypes.Collection]](c),
			container.MustGet[repository.UserMediaItemDataRepository[*mediatypes.Playlist]](c),
		)
	})
}

func registerDataRepository[T mediatypes.MediaData](c *container.Container, db *gorm.DB) {
	// Register core user media item data repositories
	container.RegisterFactory[repository.CoreUserMediaItemDataRepository[T]](c, func(c *container.Container) repository.CoreUserMediaItemDataRepository[T] {
		return repository.NewCoreUserMediaItemDataRepository[T](db)
	})
	// Register user media item data repositories
	container.RegisterFactory[repository.UserMediaItemDataRepository[T]](c, func(c *container.Container) repository.UserMediaItemDataRepository[T] {
		coreRepo := container.MustGet[repository.CoreUserMediaItemDataRepository[T]](c)
		return repository.NewUserMediaItemDataRepository[T](db, coreRepo)
	})

}
