package services

import (
	"suasor/client/types"
	"suasor/services"
)

// Interface definitions (unchanged)
type ClientServices interface {
	EmbyService() services.ClientService[*types.EmbyConfig]
	JellyfinService() services.ClientService[*types.JellyfinConfig]
	PlexService() services.ClientService[*types.PlexConfig]
	SubsonicService() services.ClientService[*types.SubsonicConfig]
	SonarrService() services.ClientService[*types.SonarrConfig]
	RadarrService() services.ClientService[*types.RadarrConfig]
	LidarrService() services.ClientService[*types.LidarrConfig]
	ClaudeService() services.ClientService[*types.ClaudeConfig]
	OpenAIService() services.ClientService[*types.OpenAIConfig]
	OllamaService() services.ClientService[*types.OllamaConfig]
	AllServices() map[string]services.ClientService[types.ClientConfig]
}
