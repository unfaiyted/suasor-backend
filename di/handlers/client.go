// app/di/handlers/client.go
package handlers

import (
	"context"
	"suasor/app/container"
	"suasor/client"
	"suasor/client/types"
	"suasor/handlers"
	apphandlers "suasor/handlers/bundles"
	"suasor/services"
)

// RegisterClientHandlers registers the client-related handlers
func RegisterClientHandlers(ctx context.Context, c *container.Container) {
	// Media client handlers
	registerClientHandler[*types.EmbyConfig](c)
	registerClientHandler[*types.JellyfinConfig](c)
	registerClientHandler[*types.PlexConfig](c)
	registerClientHandler[*types.SubsonicConfig](c)
	registerClientHandler[*types.RadarrConfig](c)
	registerClientHandler[*types.LidarrConfig](c)
	registerClientHandler[*types.SonarrConfig](c)
	registerClientHandler[*types.ClaudeConfig](c)
	registerClientHandler[*types.OpenAIConfig](c)
	registerClientHandler[*types.OllamaConfig](c)

	// Register clients collection handler
	container.RegisterFactory[*handlers.ClientsHandler](c, func(c *container.Container) *handlers.ClientsHandler {
		embyService := container.MustGet[services.ClientService[*types.EmbyConfig]](c)
		jellyfinService := container.MustGet[services.ClientService[*types.JellyfinConfig]](c)
		plexService := container.MustGet[services.ClientService[*types.PlexConfig]](c)
		subsonicService := container.MustGet[services.ClientService[*types.SubsonicConfig]](c)
		sonarrService := container.MustGet[services.ClientService[*types.SonarrConfig]](c)
		radarrService := container.MustGet[services.ClientService[*types.RadarrConfig]](c)
		lidarrService := container.MustGet[services.ClientService[*types.LidarrConfig]](c)
		claudeService := container.MustGet[services.ClientService[*types.ClaudeConfig]](c)
		openaiService := container.MustGet[services.ClientService[*types.OpenAIConfig]](c)
		ollamaService := container.MustGet[services.ClientService[*types.OllamaConfig]](c)

		return handlers.NewClientsHandler(
			embyService, jellyfinService, plexService, subsonicService,
			sonarrService, radarrService, lidarrService,
			claudeService, openaiService, ollamaService,
		)
	})

	// Register AI handlers TODO: move to AI bundle
	container.RegisterFactory[*handlers.AIHandler[*types.ClaudeConfig]](c, func(c *container.Container) *handlers.AIHandler[*types.ClaudeConfig] {
		clientFactory := container.MustGet[*client.ClientFactoryService](c)
		clientService := container.MustGet[services.ClientService[*types.ClaudeConfig]](c)

		return handlers.NewAIHandler(
			clientFactory,
			clientService,
		)
	})
	container.RegisterFactory[*handlers.AIHandler[*types.OpenAIConfig]](c, func(c *container.Container) *handlers.AIHandler[*types.OpenAIConfig] {
		clientFactory := container.MustGet[*client.ClientFactoryService](c)
		clientService := container.MustGet[services.ClientService[*types.OpenAIConfig]](c)

		return handlers.NewAIHandler(
			clientFactory,
			clientService,
		)
	})
	container.RegisterFactory[*handlers.AIHandler[*types.OllamaConfig]](c, func(c *container.Container) *handlers.AIHandler[*types.OllamaConfig] {
		clientFactory := container.MustGet[*client.ClientFactoryService](c)
		clientService := container.MustGet[services.ClientService[*types.OllamaConfig]](c)

		return handlers.NewAIHandler(
			clientFactory,
			clientService,
		)
	})

	// Register AIClientHandlers for the router
	container.RegisterFactory[apphandlers.AIClientHandlers](c, func(c *container.Container) apphandlers.AIClientHandlers {
		claudeHandler := container.MustGet[*handlers.AIHandler[*types.ClaudeConfig]](c)
		openaiHandler := container.MustGet[*handlers.AIHandler[*types.OpenAIConfig]](c)
		ollamaHandler := container.MustGet[*handlers.AIHandler[*types.OllamaConfig]](c)

		return apphandlers.NewAIClientHandlers(
			claudeHandler,
			openaiHandler,
			ollamaHandler,
		)
	})
}

func registerClientHandler[T types.ClientConfig](c *container.Container) {
	container.RegisterFactory[handlers.ClientHandler[T]](c, func(c *container.Container) handlers.ClientHandler[T] {
		service := container.MustGet[services.ClientService[T]](c)
		return *handlers.NewClientHandler[T](service)
	})
}
