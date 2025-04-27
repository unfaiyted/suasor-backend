package services

import (
	"context"
	"suasor/clients"
	mediatypes "suasor/clients/media/types"
	types "suasor/clients/types"
	"suasor/di/container"
	"suasor/repository"
	"suasor/services"
)

func registerClientServices(ctx context.Context, c *container.Container) {
	// Media clients
	registerClientService[*types.JellyfinConfig](c)
	registerClientService[*types.EmbyConfig](c)
	registerClientService[*types.PlexConfig](c)
	registerClientService[*types.SubsonicConfig](c)
	registerClientService[*types.RadarrConfig](c)
	registerClientService[*types.LidarrConfig](c)
	registerClientService[*types.SonarrConfig](c)
	registerClientService[*types.ClaudeConfig](c)
	registerClientService[*types.OpenAIConfig](c)
	registerClientService[*types.OllamaConfig](c)

	// Register AutomationClientService
	container.RegisterFactory[services.AutomationClientService](c, func(c *container.Container) services.AutomationClientService {
		clientFactory := container.MustGet[*clients.ClientProviderFactoryService](c)
		repo := container.MustGet[repository.ClientRepository[types.ClientAutomationConfig]](c)
		return services.NewAutomationClientService(repo, clientFactory)
	})

}

func registerClientService[T types.ClientConfig](c *container.Container) {
	container.RegisterFactory[services.ClientService[T]](c, func(c *container.Container) services.ClientService[T] {
		clientFactory := container.MustGet[*clients.ClientProviderFactoryService](c)
		repo := container.MustGet[repository.ClientRepository[T]](c)
		return services.NewClientService[T](clientFactory, repo)
	})
}

func registerMediaTypeService[T types.ClientMediaConfig, U mediatypes.MediaData](c *container.Container) {
	container.RegisterFactory[services.ClientMediaItemService[T, U]](c, func(c *container.Container) services.ClientMediaItemService[T, U] {
		// Dependencies
		clientFactory := container.MustGet[*clients.ClientProviderFactoryService](c)
		clientRepo := container.MustGet[repository.ClientRepository[T]](c)
		itemRepo := container.MustGet[repository.ClientMediaItemRepository[U]](c)
		userService := container.MustGet[services.UserMediaItemService[U]](c)

		return services.NewClientMediaItemService[T, U](
			userService,
			clientRepo,
			itemRepo,
			clientFactory,
		)
	})
}
