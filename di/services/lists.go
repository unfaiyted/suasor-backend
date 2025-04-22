// app/di/services/lists.go
package services

import (
	"context"
	mediatypes "suasor/client/media/types"
	"suasor/container"
	"suasor/repository"
	apprepository "suasor/repository/bundles"
	repobundles "suasor/repository/bundles"

	"suasor/services"
)

// RegisterMediaListServices registers services for media lists
func RegisterMediaListServices(ctx context.Context, c *container.Container) {
	// Register CoreMediaItemService for Playlists
	container.RegisterFactory[services.CoreMediaItemService[*mediatypes.Playlist]](c, func(c *container.Container) services.CoreMediaItemService[*mediatypes.Playlist] {
		repos := container.MustGet[apprepository.CoreMediaItemRepositories](c)
		return services.NewCoreMediaItemService[*mediatypes.Playlist](repos.PlaylistRepo())
	})

	// Register CoreMediaItemService for Collections
	container.RegisterFactory[services.CoreMediaItemService[*mediatypes.Collection]](c, func(c *container.Container) services.CoreMediaItemService[*mediatypes.Collection] {
		repos := container.MustGet[apprepository.CoreMediaItemRepositories](c)
		return services.NewCoreMediaItemService[*mediatypes.Collection](repos.CollectionRepo())
	})

	// Register UserMediaItemService for Playlists
	container.RegisterFactory[services.UserMediaItemService[*mediatypes.Playlist]](c, func(c *container.Container) services.UserMediaItemService[*mediatypes.Playlist] {
		coreService := container.MustGet[services.CoreMediaItemService[*mediatypes.Playlist]](c)
		userRepo := container.MustGet[repository.UserMediaItemRepository[*mediatypes.Playlist]](c)
		return services.NewUserMediaItemService[*mediatypes.Playlist](coreService, userRepo)
	})

	// Register UserMediaItemService for Collections
	container.RegisterFactory[services.UserMediaItemService[*mediatypes.Collection]](c, func(c *container.Container) services.UserMediaItemService[*mediatypes.Collection] {
		coreService := container.MustGet[services.CoreMediaItemService[*mediatypes.Collection]](c)
		userRepo := container.MustGet[repository.UserMediaItemRepository[*mediatypes.Collection]](c)
		return services.NewUserMediaItemService[*mediatypes.Collection](coreService, userRepo)
	})

	// Register CoreListService for Playlists
	container.RegisterFactory[services.CoreListService[*mediatypes.Playlist]](c, func(c *container.Container) services.CoreListService[*mediatypes.Playlist] {
		repos := container.MustGet[apprepository.CoreMediaItemRepositories](c)
		return services.NewCoreListService[*mediatypes.Playlist](repos.PlaylistRepo())
	})

	// Register CoreListService for Collections
	container.RegisterFactory[services.CoreListService[*mediatypes.Collection]](c, func(c *container.Container) services.CoreListService[*mediatypes.Collection] {
		repos := container.MustGet[apprepository.CoreMediaItemRepositories](c)
		return services.NewCoreListService[*mediatypes.Collection](repos.CollectionRepo())
	})

	// Register UserListService for Playlists
	container.RegisterFactory[services.UserListService[*mediatypes.Playlist]](c, func(c *container.Container) services.UserListService[*mediatypes.Playlist] {
		coreService := container.MustGet[services.CoreListService[*mediatypes.Playlist]](c)
		userRepos := container.MustGet[repobundles.UserMediaItemRepositories](c)
		userDataRepos := container.MustGet[repobundles.UserMediaDataRepositories](c)

		// Get the specific repositories
		userItemRepo := userRepos.PlaylistUserRepo()
		userDataRepo := userDataRepos.PlaylistDataRepo()

		return services.NewUserListService[*mediatypes.Playlist](coreService, userItemRepo, userDataRepo)
	})

	// Register UserListService for Collections
	container.RegisterFactory[services.UserListService[*mediatypes.Collection]](c, func(c *container.Container) services.UserListService[*mediatypes.Collection] {
		coreService := container.MustGet[services.CoreListService[*mediatypes.Collection]](c)
		userRepos := container.MustGet[repobundles.UserMediaItemRepositories](c)
		userDataRepos := container.MustGet[repobundles.UserMediaDataRepositories](c)

		// Get the specific repositories
		userItemRepo := userRepos.CollectionUserRepo()
		userDataRepo := userDataRepos.CollectionDataRepo()
		return services.NewUserListService[*mediatypes.Collection](coreService, userItemRepo, userDataRepo)
	})
}
