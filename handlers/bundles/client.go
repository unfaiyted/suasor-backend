package bundles

import (
	"suasor/clients/types"
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

// ClientHandlers contains client-related handlers
// type clientHandlers struct {
// 	// Master Handler
// 	ClientsHandler *handlers.ClientsHandler
//
// 	// Media Clients
// 	EmbyHandler     *handlers.ClientHandler
// 	JellyfinHandler *handlers.ClientHandler
// 	PlexHandler     *handlers.ClientHandler
// 	SubsonicHandler *handlers.ClientHandler
//
// 	// Automation Clients
// 	RadarrHandler *handlers.ClientHandler
// 	SonarrHandler *handlers.ClientHandler
// 	LidarrHandler *handlers.ClientHandler
//
// 	// AI Clients
// 	AIHandler     *handlers.AIHandler
// 	ClaudeHandler *handlers.ClientHandler
// 	OpenAIHandler *handlers.ClientHandler
// 	OllamaHandler *handlers.ClientHandler
//
// 	// Metadata
// 	MetadataHandler *handlers.MetadataClientHandler
// 	CreditHandler   *handlers.CreditHandler
// }

type AIClientHandlers interface {
	ClaudeAIHandler() *handlers.AIHandler[*types.ClaudeConfig]
	OpenAIHandler() *handlers.AIHandler[*types.OpenAIConfig]
	OllamaHandler() *handlers.AIHandler[*types.OllamaConfig]
}
