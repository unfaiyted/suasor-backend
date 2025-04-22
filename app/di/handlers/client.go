// app/di/handlers/client.go
package handlers

import (
	"context"
	"suasor/app/container"
	apphandlers "suasor/app/handlers"
	"suasor/client"
	"suasor/client/types"
	"suasor/handlers"
	"suasor/services"
)

// RegisterClientHandlers registers the client-related handlers
func RegisterClientHandlers(ctx context.Context, c *container.Container) {
	// Register client type handlers
	container.RegisterFactory[*handlers.ClientHandler[*types.EmbyConfig]](c, func(c *container.Container) *handlers.ClientHandler[*types.EmbyConfig] {
		service := container.MustGet[services.ClientService[*types.EmbyConfig]](c)
		return handlers.NewClientHandler[*types.EmbyConfig](service)
	})

	container.RegisterFactory[*handlers.ClientHandler[*types.JellyfinConfig]](c, func(c *container.Container) *handlers.ClientHandler[*types.JellyfinConfig] {
		service := container.MustGet[services.ClientService[*types.JellyfinConfig]](c)
		return handlers.NewClientHandler[*types.JellyfinConfig](service)
	})

	container.RegisterFactory[*handlers.ClientHandler[*types.PlexConfig]](c, func(c *container.Container) *handlers.ClientHandler[*types.PlexConfig] {
		service := container.MustGet[services.ClientService[*types.PlexConfig]](c)
		return handlers.NewClientHandler[*types.PlexConfig](service)
	})

	container.RegisterFactory[*handlers.ClientHandler[*types.SubsonicConfig]](c, func(c *container.Container) *handlers.ClientHandler[*types.SubsonicConfig] {
		service := container.MustGet[services.ClientService[*types.SubsonicConfig]](c)
		return handlers.NewClientHandler[*types.SubsonicConfig](service)
	})

	container.RegisterFactory[*handlers.ClientHandler[*types.SonarrConfig]](c, func(c *container.Container) *handlers.ClientHandler[*types.SonarrConfig] {
		service := container.MustGet[services.ClientService[*types.SonarrConfig]](c)
		return handlers.NewClientHandler[*types.SonarrConfig](service)
	})

	container.RegisterFactory[*handlers.ClientHandler[*types.RadarrConfig]](c, func(c *container.Container) *handlers.ClientHandler[*types.RadarrConfig] {
		service := container.MustGet[services.ClientService[*types.RadarrConfig]](c)
		return handlers.NewClientHandler[*types.RadarrConfig](service)
	})

	container.RegisterFactory[*handlers.ClientHandler[*types.LidarrConfig]](c, func(c *container.Container) *handlers.ClientHandler[*types.LidarrConfig] {
		service := container.MustGet[services.ClientService[*types.LidarrConfig]](c)
		return handlers.NewClientHandler[*types.LidarrConfig](service)
	})

	container.RegisterFactory[*handlers.ClientHandler[*types.ClaudeConfig]](c, func(c *container.Container) *handlers.ClientHandler[*types.ClaudeConfig] {
		service := container.MustGet[services.ClientService[*types.ClaudeConfig]](c)
		return handlers.NewClientHandler[*types.ClaudeConfig](service)
	})

	container.RegisterFactory[*handlers.ClientHandler[*types.OpenAIConfig]](c, func(c *container.Container) *handlers.ClientHandler[*types.OpenAIConfig] {
		service := container.MustGet[services.ClientService[*types.OpenAIConfig]](c)
		return handlers.NewClientHandler[*types.OpenAIConfig](service)
	})

	container.RegisterFactory[*handlers.ClientHandler[*types.OllamaConfig]](c, func(c *container.Container) *handlers.ClientHandler[*types.OllamaConfig] {
		service := container.MustGet[services.ClientService[*types.OllamaConfig]](c)
		return handlers.NewClientHandler[*types.OllamaConfig](service)
	})

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

	// Register AI handlers
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
		
		return &AIClientHandlersImpl{
			claudeHandler:  claudeHandler,
			openaiHandler:  openaiHandler,
			ollamaHandler:  ollamaHandler,
		}
	})
}

// AIClientHandlersImpl implements the AIClientHandlers interface
type AIClientHandlersImpl struct {
	claudeHandler  *handlers.AIHandler[*types.ClaudeConfig]
	openaiHandler  *handlers.AIHandler[*types.OpenAIConfig]
	ollamaHandler  *handlers.AIHandler[*types.OllamaConfig]
}

func (h *AIClientHandlersImpl) ClaudeAIHandler() *handlers.AIHandler[*types.ClaudeConfig] {
	return h.claudeHandler
}

func (h *AIClientHandlersImpl) OpenAIHandler() *handlers.AIHandler[*types.OpenAIConfig] {
	return h.openaiHandler
}

func (h *AIClientHandlersImpl) OllamaHandler() *handlers.AIHandler[*types.OllamaConfig] {
	return h.ollamaHandler
}