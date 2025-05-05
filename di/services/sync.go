package services

import (
	"suasor/clients"
	"suasor/clients/media/types"
	"suasor/di/container"
	"suasor/repository"
	repobundle "suasor/repository/bundles"
	"suasor/services"
)

// Register list sync services in the dependency injection container
func RegisterListSyncServices(c *container.Container) {
	// Register playlist sync service
	container.RegisterFactory[services.ListSyncService[*types.Playlist]](c, func(c *container.Container) services.ListSyncService[*types.Playlist] {
		clientRepos := container.MustGet[repobundle.ClientRepositories](c)
		clientFactory := container.MustGet[*clients.ClientProviderFactoryService](c)
		listService := container.MustGet[services.UserListService[*types.Playlist]](c)
		mediaItemRepo := container.MustGet[repository.UserMediaItemRepository[*types.Playlist]](c)

		return services.NewListSyncService[*types.Playlist](
			clientRepos,
			clientFactory,
			listService,
			mediaItemRepo,
		)
	})

	// Register collection sync service
	container.RegisterFactory[services.ListSyncService[*types.Collection]](c, func(c *container.Container) services.ListSyncService[*types.Collection] {
		clientRepos := container.MustGet[repobundle.ClientRepositories](c)
		clientFactory := container.MustGet[*clients.ClientProviderFactoryService](c)
		listService := container.MustGet[services.UserListService[*types.Collection]](c)
		mediaItemRepo := container.MustGet[repository.UserMediaItemRepository[*types.Collection]](c)

		return services.NewListSyncService[*types.Collection](
			clientRepos,
			clientFactory,
			listService,
			mediaItemRepo,
		)
	})
}