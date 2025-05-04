package services

import (
	"suasor/clients/media/types"
	"suasor/di/container"
	"suasor/di/repositories"
	"suasor/services"
)

// Register list sync services in the dependency injection container
func RegisterListSyncServices(c *container.Container) {
	// Register playlist sync service
	c.Register(func(
		clientRepo container.LazyValue[repositories.ClientRepository],
		listService container.LazyValue[services.UserListService[*types.Playlist]],
		mediaItemRepo container.LazyValue[repositories.UserMediaItemRepository[*types.Playlist]],
	) services.ListSyncService[*types.Playlist] {
		return services.NewListSyncService[*types.Playlist](
			clientRepo.Value(),
			listService.Value(),
			mediaItemRepo.Value(),
		)
	})

	// Register collection sync service
	c.Register(func(
		clientRepo container.LazyValue[repositories.ClientRepository],
		listService container.LazyValue[services.UserListService[*types.Collection]],
		mediaItemRepo container.LazyValue[repositories.UserMediaItemRepository[*types.Collection]],
	) services.ListSyncService[*types.Collection] {
		return services.NewListSyncService[*types.Collection](
			clientRepo.Value(),
			listService.Value(),
			mediaItemRepo.Value(),
		)
	})
}