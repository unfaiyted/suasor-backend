package repository

import (
	"suasor/client/types"
	"suasor/repository"
)

// These repositories store client configurations by client type
type ClientRepositories interface {
	EmbyRepo() repository.ClientRepository[*types.EmbyConfig]
	JellyfinRepo() repository.ClientRepository[*types.JellyfinConfig]
	PlexRepo() repository.ClientRepository[*types.PlexConfig]
	SubsonicRepo() repository.ClientRepository[*types.SubsonicConfig]
	SonarrRepo() repository.ClientRepository[*types.SonarrConfig]
	RadarrRepo() repository.ClientRepository[*types.RadarrConfig]
	LidarrRepo() repository.ClientRepository[*types.LidarrConfig]
	ClaudeRepo() repository.ClientRepository[*types.ClaudeConfig]
	OpenAIRepo() repository.ClientRepository[*types.OpenAIConfig]
	OllamaRepo() repository.ClientRepository[*types.OllamaConfig]
}

func NewClientRepositories(
	embyRepo repository.ClientRepository[*types.EmbyConfig],
	jellyfinRepo repository.ClientRepository[*types.JellyfinConfig],
	plexRepo repository.ClientRepository[*types.PlexConfig],
	subsonicRepo repository.ClientRepository[*types.SubsonicConfig],
	sonarrRepo repository.ClientRepository[*types.SonarrConfig],
	radarrRepo repository.ClientRepository[*types.RadarrConfig],
	lidarrRepo repository.ClientRepository[*types.LidarrConfig],
	claudeRepo repository.ClientRepository[*types.ClaudeConfig],
	openaiRepo repository.ClientRepository[*types.OpenAIConfig],
	ollamaRepo repository.ClientRepository[*types.OllamaConfig],

) ClientRepositories {
	return &clientRepositories{
		embyRepo:     embyRepo,
		jellyfinRepo: jellyfinRepo,
		plexRepo:     plexRepo,
		subsonicRepo: subsonicRepo,
		sonarrRepo:   sonarrRepo,
		radarrRepo:   radarrRepo,
		lidarrRepo:   lidarrRepo,
		claudeRepo:   claudeRepo,
		openaiRepo:   openaiRepo,
		ollamaRepo:   ollamaRepo,
	}
}

// clientRepositories implements ClientRepositories
type clientRepositories struct {
	embyRepo     repository.ClientRepository[*types.EmbyConfig]
	jellyfinRepo repository.ClientRepository[*types.JellyfinConfig]
	plexRepo     repository.ClientRepository[*types.PlexConfig]
	subsonicRepo repository.ClientRepository[*types.SubsonicConfig]
	sonarrRepo   repository.ClientRepository[*types.SonarrConfig]
	radarrRepo   repository.ClientRepository[*types.RadarrConfig]
	lidarrRepo   repository.ClientRepository[*types.LidarrConfig]
	claudeRepo   repository.ClientRepository[*types.ClaudeConfig]
	openaiRepo   repository.ClientRepository[*types.OpenAIConfig]
	ollamaRepo   repository.ClientRepository[*types.OllamaConfig]
}

func (c *clientRepositories) EmbyRepo() repository.ClientRepository[*types.EmbyConfig] {
	return c.embyRepo
}

func (c *clientRepositories) JellyfinRepo() repository.ClientRepository[*types.JellyfinConfig] {
	return c.jellyfinRepo
}

func (c *clientRepositories) PlexRepo() repository.ClientRepository[*types.PlexConfig] {
	return c.plexRepo
}

func (c *clientRepositories) SubsonicRepo() repository.ClientRepository[*types.SubsonicConfig] {
	return c.subsonicRepo
}

func (c *clientRepositories) SonarrRepo() repository.ClientRepository[*types.SonarrConfig] {
	return c.sonarrRepo
}

func (c *clientRepositories) RadarrRepo() repository.ClientRepository[*types.RadarrConfig] {
	return c.radarrRepo
}

func (c *clientRepositories) LidarrRepo() repository.ClientRepository[*types.LidarrConfig] {
	return c.lidarrRepo
}

func (c *clientRepositories) ClaudeRepo() repository.ClientRepository[*types.ClaudeConfig] {
	return c.claudeRepo
}

func (c *clientRepositories) OpenAIRepo() repository.ClientRepository[*types.OpenAIConfig] {
	return c.openaiRepo
}

func (c *clientRepositories) OllamaRepo() repository.ClientRepository[*types.OllamaConfig] {
	return c.ollamaRepo
}
