package bundles

import (
	"context"
	"fmt"
	"suasor/client/types"
	"suasor/repository"
	"suasor/types/models"
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

	// Helpers
	GetAllClientsForUser(ctx context.Context, userID uint64) (*models.Clients, error)
	GetAllMediaClientsForUser(ctx context.Context, userID uint64) (*models.Clients, error)
	GetAllMetadataClientsForUser(ctx context.Context, userID uint64) (*models.MetadataClients, error)
	GetAllAutomationClientsForUser(ctx context.Context, userID uint64) (*models.AutomationClients, error)
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

func (c *clientRepositories) GetAllClientsForUser(ctx context.Context, userID uint64) (*models.Clients, error) {
	var clients *models.Clients
	clients.Total = 0

	// Emby clients
	embyClients, err := c.embyRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clients, fmt.Errorf("failed to get emby client configs for user: %w", err)
	}

	clients.Emby = embyClients
	clients.Total += len(embyClients)
	// Jellyfin clients
	jellyfinClients, err := c.jellyfinRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clients, fmt.Errorf("failed to get jellyfin client configs for user: %w", err)
	}
	clients.Jellyfin = jellyfinClients
	clients.Total += len(jellyfinClients)
	// Plex clients
	plexClients, err := c.plexRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clients, fmt.Errorf("failed to get plex client configs for user: %w", err)
	}
	clients.Plex = plexClients
	clients.Total += len(plexClients)
	// Subsonic clients (primarily for music)
	subsonicClients, err := c.subsonicRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clients, fmt.Errorf("failed to get subsonic client configs for user: %w", err)
	}
	clients.Subsonic = subsonicClients
	clients.Total += len(subsonicClients)
	// Sonarr clients
	sonarrClients, err := c.sonarrRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clients, fmt.Errorf("failed to get sonarr client configs for user: %w", err)
	}
	clients.Sonarr = sonarrClients
	clients.Total += len(sonarrClients)
	// Radarr clients
	radarrClients, err := c.radarrRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clients, fmt.Errorf("failed to get radarr client configs for user: %w", err)
	}
	clients.Radarr = radarrClients
	clients.Total += len(radarrClients)
	// Lidarr clients
	lidarrClients, err := c.lidarrRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clients, fmt.Errorf("failed to get lidarr client configs for user: %w", err)
	}
	clients.Lidarr = lidarrClients
	clients.Total += len(lidarrClients)
	// Claude clients
	claudeClients, err := c.claudeRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clients, fmt.Errorf("failed to get claude client configs for user: %w", err)
	}
	clients.Claude = claudeClients
	clients.Total += len(claudeClients)
	// OpenAI clients
	openaiClients, err := c.openaiRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clients, fmt.Errorf("failed to get openai client configs for user: %w", err)
	}
	clients.OpenAI = openaiClients
	clients.Total += len(openaiClients)
	// Ollama clients
	ollamaClients, err := c.ollamaRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clients, fmt.Errorf("failed to get ollama client configs for user: %w", err)
	}
	clients.Ollama = ollamaClients

	return clients, nil
}

func (c *clientRepositories) GetAllMediaClientsForUser(ctx context.Context, userID uint64) (*models.Clients, error) {
	var clients *models.Clients
	clients.Total = 0

	// Emby clients
	embyClients, err := c.embyRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clients, fmt.Errorf("failed to get emby client configs for user: %w", err)
	}
	clients.Emby = embyClients
	clients.Total += len(embyClients)
	// Jellyfin clients
	jellyfinClients, err := c.jellyfinRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clients, fmt.Errorf("failed to get jellyfin client configs for user: %w", err)
	}
	clients.Jellyfin = jellyfinClients
	clients.Total += len(jellyfinClients)
	// Plex clients
	plexClients, err := c.plexRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clients, fmt.Errorf("failed to get plex client configs for user: %w", err)
	}
	clients.Plex = plexClients
	clients.Total += len(plexClients)
	// Subsonic clients (primarily for music)
	subsonicClients, err := c.subsonicRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clients, fmt.Errorf("failed to get subsonic client configs for user: %w", err)
	}
	clients.Subsonic = subsonicClients
	clients.Total += len(subsonicClients)
	return clients, nil
}

func (c *clientRepositories) GetAllMetadataClientsForUser(ctx context.Context, userID uint64) (*models.MetadataClients, error) {
	var metadataClients *models.MetadataClients

	// Tmdb clients
	// tmdbClients, err := c.tmdbRepo.GetByUserID(ctx, userID)
	// if err != nil {
	// 	return metadataClients, fmt.Errorf("failed to get tmdb client configs for user: %w", err)
	// }
	// metadataClients.Tmdb = tmdbClients
	// // Trakt clients
	// traktClients, err := c.traktRepo.GetByUserID(ctx, userID)
	// if err != nil {
	// 	return metadataClients, fmt.Errorf("failed to get trakt client configs for user: %w", err)
	// }
	// metadataClients.Trakt = traktClients

	return metadataClients, nil
}

func (c *clientRepositories) GetAllAutomationClientsForUser(ctx context.Context, userID uint64) (*models.AutomationClients, error) {
	var automationClients *models.AutomationClients
	automationClients.Total = 0

	// Sonarr clients
	sonarrClients, err := c.sonarrRepo.GetByUserID(ctx, userID)

	if err != nil {
		return automationClients, fmt.Errorf("failed to get sonarr client configs for user: %w", err)
	}
	automationClients.Sonarr = sonarrClients
	automationClients.Total += len(sonarrClients)
	// Radarr clients
	radarrClients, err := c.radarrRepo.GetByUserID(ctx, userID)
	if err != nil {
		return automationClients, fmt.Errorf("failed to get radarr client configs for user: %w", err)
	}
	automationClients.Radarr = radarrClients
	automationClients.Total += len(radarrClients)
	// Lidarr clients
	lidarrClients, err := c.lidarrRepo.GetByUserID(ctx, userID)
	if err != nil {
		return automationClients, fmt.Errorf("failed to get lidarr client configs for user: %w", err)
	}
	automationClients.Lidarr = lidarrClients
	automationClients.Total += len(lidarrClients)
	return automationClients, nil
}
