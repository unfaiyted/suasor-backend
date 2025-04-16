package handlers

import (
	"suasor/client/types"
	"suasor/handlers"
)

type ClientHandlers interface {
	EmbyClientHandler() *handlers.ClientHandler[*types.EmbyConfig]
	JellyfinClientHandler() *handlers.ClientHandler[*types.JellyfinConfig]
	PlexClientHandler() *handlers.ClientHandler[*types.PlexConfig]
	SubsonicClientHandler() *handlers.ClientHandler[*types.SubsonicConfig]
	RadarrClientHandler() *handlers.ClientHandler[*types.RadarrConfig]
	LidarrClientHandler() *handlers.ClientHandler[*types.LidarrConfig]
	SonarrClientHandler() *handlers.ClientHandler[*types.SonarrConfig]
	ClaudeClientHandler() *handlers.ClientHandler[*types.ClaudeConfig]
	OpenAIClientHandler() *handlers.ClientHandler[*types.OpenAIConfig]
	OllamaClientHandler() *handlers.ClientHandler[*types.OllamaConfig]
}

type AIClientHandlers interface {
	ClaudeAIHandler() *handlers.AIHandler[*types.ClaudeConfig]
	OpenAIHandler() *handlers.AIHandler[*types.OpenAIConfig]
	OllamaHandler() *handlers.AIHandler[*types.OllamaConfig]
}
