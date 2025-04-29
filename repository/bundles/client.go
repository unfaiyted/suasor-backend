package bundles

import (
	"context"
	"fmt"
	"suasor/clients/types"
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

	GetAllMediaClients(ctx context.Context) (*models.MediaClientList, error)

	// Helpers
	GetAllClientsForUser(ctx context.Context, userID uint64) (*models.ClientList, error)
	GetAllMediaClientsForUser(ctx context.Context, userID uint64) (*models.MediaClientList, error)
	GetAllMetadataClientsForUser(ctx context.Context, userID uint64) (*models.MetadataClientList, error)
	GetAllAutomationClientsForUser(ctx context.Context, userID uint64) (*models.AutomationClientList, error)

	// GetClientByID(ctx context.Context, clientID uint64) (*models.Client[types.ClientConfig], error)
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

func (c *clientRepositories) GetAllClients(ctx context.Context) (*models.ClientList, error) {

	clients := models.NewClientList()

	// Emby clients
	embyClients, err := c.embyRepo.GetAll(ctx)
	if err != nil {
		return clients, fmt.Errorf("failed to get emby client configs for user: %w", err)
	}
	clients.AddEmbyArray(embyClients)

	// Jellyfin clients
	jellyfinClients, err := c.jellyfinRepo.GetAll(ctx)
	if err != nil {
		return clients, fmt.Errorf("failed to get jellyfin client configs for user: %w", err)
	}
	clients.AddJellyfinArray(jellyfinClients)

	// Plex clients
	plexClients, err := c.plexRepo.GetAll(ctx)
	if err != nil {
		return clients, fmt.Errorf("failed to get plex client configs for user: %w", err)
	}
	clients.AddPlexArray(plexClients)

	// Subsonic clients (primarily for music)
	subsonicClients, err := c.subsonicRepo.GetAll(ctx)
	if err != nil {
		return clients, fmt.Errorf("failed to get subsonic client configs for user: %w", err)
	}
	clients.AddSubsonicArray(subsonicClients)

	// Sonarr clients
	sonarrClients, err := c.sonarrRepo.GetAll(ctx)
	if err != nil {
		return clients, fmt.Errorf("failed to get sonarr client configs for user: %w", err)
	}
	clients.AddSonarrArray(sonarrClients)

	// Radarr clients
	radarrClients, err := c.radarrRepo.GetAll(ctx)
	if err != nil {
		return clients, fmt.Errorf("failed to get radarr client configs for user: %w", err)
	}
	clients.AddRadarrArray(radarrClients)

	// Lidarr clients
	lidarrClients, err := c.lidarrRepo.GetAll(ctx)
	if err != nil {
		return clients, fmt.Errorf("failed to get lidarr client configs for user: %w", err)
	}
	clients.AddLidarrArray(lidarrClients)

	// Claude clients
	claudeClients, err := c.claudeRepo.GetAll(ctx)
	if err != nil {
		return clients, fmt.Errorf("failed to get claude client configs for user: %w", err)
	}
	clients.AddClaudeArray(claudeClients)

	// OpenAI clients
	openaiClients, err := c.openaiRepo.GetAll(ctx)
	if err != nil {
		return clients, fmt.Errorf("failed to get openai client configs for user: %w", err)
	}
	clients.AddOpenAIArray(openaiClients)

	// Ollama clients
	ollamaClients, err := c.ollamaRepo.GetAll(ctx)
	if err != nil {
		return clients, fmt.Errorf("failed to get ollama client configs for user: %w", err)
	}
	clients.AddOllamaArray(ollamaClients)

	return clients, nil
}

func (c *clientRepositories) GetAllMediaClients(ctx context.Context) (*models.MediaClientList, error) {
	clients := models.NewMediaClientList()

	// Emby clients
	embyClients, err := c.embyRepo.GetAll(ctx)
	if err != nil {
		return clients, fmt.Errorf("failed to get emby client configs for user: %w", err)
	}
	clients.AddEmbyArray(embyClients)

	// Jellyfin clients
	jellyfinClients, err := c.jellyfinRepo.GetAll(ctx)
	if err != nil {
		return clients, fmt.Errorf("failed to get jellyfin client configs for user: %w", err)
	}
	clients.AddJellyfinArray(jellyfinClients)

	// Plex clients
	plexClients, err := c.plexRepo.GetAll(ctx)
	if err != nil {
		return clients, fmt.Errorf("failed to get plex client configs for user: %w", err)
	}
	clients.AddPlexArray(plexClients)

	// Subsonic clients (primarily for music)
	subsonicClients, err := c.subsonicRepo.GetAll(ctx)
	if err != nil {
		return clients, fmt.Errorf("failed to get subsonic client configs for user: %w", err)
	}
	clients.AddSubsonicArray(subsonicClients)

	return clients, nil
}

func (c *clientRepositories) GetAllClientsForUser(ctx context.Context, userID uint64) (*models.ClientList, error) {

	clients := models.NewClientList()

	// Emby clients
	embyClients, err := c.embyRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clients, fmt.Errorf("failed to get emby client configs for user: %w", err)
	}

	clients.AddEmbyArray(embyClients)
	// Jellyfin clients
	jellyfinClients, err := c.jellyfinRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clients, fmt.Errorf("failed to get jellyfin client configs for user: %w", err)
	}
	clients.AddJellyfinArray(jellyfinClients)
	// Plex clients
	plexClients, err := c.plexRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clients, fmt.Errorf("failed to get plex client configs for user: %w", err)
	}
	clients.AddPlexArray(plexClients)
	// Subsonic clients (primarily for music)
	subsonicClients, err := c.subsonicRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clients, fmt.Errorf("failed to get subsonic client configs for user: %w", err)
	}
	clients.AddSubsonicArray(subsonicClients)
	// Sonarr clients
	sonarrClients, err := c.sonarrRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clients, fmt.Errorf("failed to get sonarr client configs for user: %w", err)
	}
	clients.AddSonarrArray(sonarrClients)
	// Radarr clients
	radarrClients, err := c.radarrRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clients, fmt.Errorf("failed to get radarr client configs for user: %w", err)
	}
	clients.AddRadarrArray(radarrClients)
	// Lidarr clients
	lidarrClients, err := c.lidarrRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clients, fmt.Errorf("failed to get lidarr client configs for user: %w", err)
	}
	clients.AddLidarrArray(lidarrClients)
	// Claude clients
	claudeClients, err := c.claudeRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clients, fmt.Errorf("failed to get claude client configs for user: %w", err)
	}
	clients.AddClaudeArray(claudeClients)
	// OpenAI clients
	openaiClients, err := c.openaiRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clients, fmt.Errorf("failed to get openai client configs for user: %w", err)
	}
	clients.AddOpenAIArray(openaiClients)
	// Ollama clients
	ollamaClients, err := c.ollamaRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clients, fmt.Errorf("failed to get ollama client configs for user: %w", err)
	}
	clients.AddOllamaArray(ollamaClients)

	return clients, nil
}

func (c *clientRepositories) GetAllMediaClientsForUser(ctx context.Context, userID uint64) (*models.MediaClientList, error) {
	clients := models.NewMediaClientList()

	// Emby clients
	embyClients, err := c.embyRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clients, fmt.Errorf("failed to get emby client configs for user: %w", err)
	}
	clients.AddEmbyArray(embyClients)

	// Jellyfin clients
	jellyfinClients, err := c.jellyfinRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clients, fmt.Errorf("failed to get jellyfin client configs for user: %w", err)
	}
	clients.AddJellyfinArray(jellyfinClients)

	// Plex clients
	plexClients, err := c.plexRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clients, fmt.Errorf("failed to get plex client configs for user: %w", err)
	}
	clients.AddPlexArray(plexClients)

	// Subsonic clients (primarily for music)
	subsonicClients, err := c.subsonicRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clients, fmt.Errorf("failed to get subsonic client configs for user: %w", err)
	}
	clients.AddSubsonicArray(subsonicClients)
	return clients, nil
}

func (c *clientRepositories) GetAllMetadataClientsForUser(ctx context.Context, userID uint64) (*models.MetadataClientList, error) {
	// Create an empty metadata client list
	metadataClients := &models.MetadataClientList{}

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

func (c *clientRepositories) GetAllAutomationClientsForUser(ctx context.Context, userID uint64) (*models.AutomationClientList, error) {
	automationClients := models.NewAutomationClientList()
	// Sonarr clients
	sonarrClients, err := c.sonarrRepo.GetByUserID(ctx, userID)
	if err != nil {
		return automationClients, fmt.Errorf("failed to get sonarr client configs for user: %w", err)
	}
	automationClients.AddSonarrArray(sonarrClients)

	// Radarr clients
	radarrClients, err := c.radarrRepo.GetByUserID(ctx, userID)
	if err != nil {
		return automationClients, fmt.Errorf("failed to get radarr client configs for user: %w", err)
	}
	automationClients.AddRadarrArray(radarrClients)

	// Lidarr clients
	lidarrClients, err := c.lidarrRepo.GetByUserID(ctx, userID)
	if err != nil {
		return automationClients, fmt.Errorf("failed to get lidarr client configs for user: %w", err)
	}
	automationClients.AddLidarrArray(lidarrClients)
	return automationClients, nil
}
