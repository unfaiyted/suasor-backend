// app/di/services/lists.go
package services

import (
	"context"
	mediatypes "suasor/clients/media/types"
	"suasor/di/container"
	apprepository "suasor/repository/bundles"
	repobundles "suasor/repository/bundles"

	"suasor/services"

	"suasor/clients"
	"suasor/clients/types"
	"suasor/repository"
)

func registerListServices(ctx context.Context, c *container.Container) {

	registerMediaListServices(ctx, c)

	registerClientListService[*types.JellyfinConfig, *mediatypes.Collection](c)
	registerClientListService[*types.EmbyConfig, *mediatypes.Collection](c)
	registerClientListService[*types.PlexConfig, *mediatypes.Collection](c)
	registerClientListService[*types.SubsonicConfig, *mediatypes.Collection](c)

	registerClientListService[*types.EmbyConfig, *mediatypes.Playlist](c)
	registerClientListService[*types.JellyfinConfig, *mediatypes.Playlist](c)
	registerClientListService[*types.PlexConfig, *mediatypes.Playlist](c)
	registerClientListService[*types.SubsonicConfig, *mediatypes.Playlist](c)

}

// RegisterMediaListServices registers services for media lists
func registerMediaListServices(ctx context.Context, c *container.Container) {

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
		coreListService := container.MustGet[services.CoreListService[*mediatypes.Playlist]](c)
		userRepo := container.MustGet[repository.UserRepository](c)
		userItemRepos := container.MustGet[repobundles.UserMediaItemRepositories](c)
		userDataRepos := container.MustGet[repobundles.UserMediaDataRepositories](c)
		listRepo := container.MustGet[repository.CoreListRepository[*mediatypes.Playlist]](c)

		// Get the specific repositories
		userItemRepo := userItemRepos.PlaylistUserRepo()
		userDataRepo := userDataRepos.PlaylistDataRepo()

		// coreListService CoreListService[T],
		// userRepo repository.UserRepository,
		// listRepo repository.CoreListRepository[T],
		// userItemRepo repository.UserMediaItemRepository[T],
		// userDataRepo repository.UserMediaItemDataRepository[T],

		return services.NewUserListService[*mediatypes.Playlist](coreListService, userRepo, listRepo, userItemRepo, userDataRepo)
	})

	// Register UserListService for Collections
	container.RegisterFactory[services.UserListService[*mediatypes.Collection]](c, func(c *container.Container) services.UserListService[*mediatypes.Collection] {
		coreListService := container.MustGet[services.CoreListService[*mediatypes.Collection]](c)
		userRepo := container.MustGet[repository.UserRepository](c)
		userItemRepos := container.MustGet[repobundles.UserMediaItemRepositories](c)
		userDataRepos := container.MustGet[repobundles.UserMediaDataRepositories](c)
		listRepo := container.MustGet[repository.CoreListRepository[*mediatypes.Collection]](c)

		// Get the specific repositories
		userItemRepo := userItemRepos.CollectionUserRepo()
		userDataRepo := userDataRepos.CollectionDataRepo()
		return services.NewUserListService[*mediatypes.Collection](coreListService, userRepo, listRepo, userItemRepo, userDataRepo)
	})
}

func registerClientListService[T types.ClientMediaConfig, U mediatypes.ListData](c *container.Container) {
	container.RegisterFactory[services.ClientListService[T, U]](c, func(c *container.Container) services.ClientListService[T, U] {
		userListService := container.MustGet[services.UserListService[U]](c)
		userItemRepo := container.MustGet[repository.ClientMediaItemRepository[U]](c)
		clientRepo := container.MustGet[repository.ClientRepository[T]](c)
		clientFactory := container.MustGet[*clients.ClientProviderFactoryService](c)
		return services.NewClientListService[T, U](userListService, userItemRepo, clientRepo, clientFactory)
	})
}
