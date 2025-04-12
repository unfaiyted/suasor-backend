package repository

import (
	"context"
	"fmt"
	"suasor/client/types"
	"suasor/types/models"
)

// ClientRepositoryCollection defines an interface for accessing multiple client repositories
type ClientRepositoryCollection interface {
	AllRepos() ClientRepoCollection
	GetAllByCategory(ctx context.Context, category types.ClientCategory) ClientRepoCollection
	GetAllMediaClientsForUser(ctx context.Context, userID uint64) (MediaClients, error)
	GetAllMetadataClientsForUser(ctx context.Context, userID uint64) (MetadataClients, error)
	GetAllAutomationClientsForUser(ctx context.Context, userID uint64) (AutomationClients, error)
}

// clientRepoCollection implements ClientRepositoryCollection
type clientRepoCollection struct {
	embyRepo     ClientRepository[*types.EmbyConfig]
	jellyfinRepo ClientRepository[*types.JellyfinConfig]
	plexRepo     ClientRepository[*types.PlexConfig]
	subsonicRepo ClientRepository[*types.SubsonicConfig]
	sonarrRepo   ClientRepository[*types.SonarrConfig]
	radarrRepo   ClientRepository[*types.RadarrConfig]
	lidarrRepo   ClientRepository[*types.LidarrConfig]
	claudeRepo   ClientRepository[*types.ClaudeConfig]
	openaiRepo   ClientRepository[*types.OpenAIConfig]
	ollamaRepo   ClientRepository[*types.OllamaConfig]
}

// NewClientRepositoryCollection creates a new ClientRepositoryCollection
func NewClientRepositoryCollection(
	embyRepo ClientRepository[*types.EmbyConfig],
	jellyfinRepo ClientRepository[*types.JellyfinConfig],
	plexRepo ClientRepository[*types.PlexConfig],
	subsonicRepo ClientRepository[*types.SubsonicConfig],
	sonarrRepo ClientRepository[*types.SonarrConfig],
	radarrRepo ClientRepository[*types.RadarrConfig],
	lidarrRepo ClientRepository[*types.LidarrConfig],
	claudeRepo ClientRepository[*types.ClaudeConfig],
	openaiRepo ClientRepository[*types.OpenAIConfig],
	ollamaRepo ClientRepository[*types.OllamaConfig],

) ClientRepositoryCollection {
	return &clientRepoCollection{
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

// AllRepos returns all client repositories
func (c *clientRepoCollection) AllRepos() ClientRepoCollection {
	return ClientRepoCollection{}
}

// GetAllByCategory returns all client repositories for a given category
func (c *clientRepoCollection) GetAllByCategory(ctx context.Context, category types.ClientCategory) ClientRepoCollection {
	// Just return all for now
	return c.AllRepos()
}

type MediaClients struct {
	Emby     []*models.Client[*types.EmbyConfig]
	Jellyfin []*models.Client[*types.JellyfinConfig]
	Plex     []*models.Client[*types.PlexConfig]
	Subsonic []*models.Client[*types.SubsonicConfig]
}

type AIClients struct {
	Claude []*models.Client[*types.ClaudeConfig]
	OpenAI []*models.Client[*types.OpenAIConfig]
	Ollama []*models.Client[*types.OllamaConfig]
}

type AutomationClients struct {
	Sonarr []*models.Client[*types.SonarrConfig]
	Radarr []*models.Client[*types.RadarrConfig]
	Lidarr []*models.Client[*types.LidarrConfig]
}

type MetadataClients struct {
	// Tmdb []*models.Client[*types.TmdbConfig]
}

func (c *clientRepoCollection) GetAllMediaClientsForUser(ctx context.Context, userID uint64) (MediaClients, error) {
	var clientConfigs = MediaClients{}

	embyRepoClients, err := c.embyRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clientConfigs, fmt.Errorf("failed to get emby client configs for user: %w", err)
	}
	clientConfigs.Emby = embyRepoClients
	jellyfinRepoClients, err := c.jellyfinRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clientConfigs, fmt.Errorf("failed to get jellyfin client configs for user: %w", err)
	}
	clientConfigs.Jellyfin = jellyfinRepoClients
	plexRepoClients, err := c.plexRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clientConfigs, fmt.Errorf("failed to get plex client configs for user: %w", err)
	}
	clientConfigs.Plex = plexRepoClients
	subsonicRepoClients, err := c.subsonicRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clientConfigs, fmt.Errorf("failed to get subsonic client configs for user: %w", err)
	}
	clientConfigs.Subsonic = subsonicRepoClients

	return clientConfigs, nil
}
func (c *clientRepoCollection) GetAllMetadataClientsForUser(ctx context.Context, userID uint64) (MetadataClients, error) {
	var clientConfigs = MetadataClients{}
	return clientConfigs, nil
}
func (c *clientRepoCollection) GetAllAutomationClientsForUser(ctx context.Context, userID uint64) (AutomationClients, error) {
	clientConfigs := AutomationClients{}
	sonarr, err := c.sonarrRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clientConfigs, fmt.Errorf("failed to get sonarr client configs for user: %w", err)
	}
	clientConfigs.Sonarr = sonarr
	radarr, err := c.radarrRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clientConfigs, fmt.Errorf("failed to get radarr client configs for user: %w", err)
	}
	clientConfigs.Radarr = radarr
	lidarr, err := c.lidarrRepo.GetByUserID(ctx, userID)
	if err != nil {
		return clientConfigs, fmt.Errorf("failed to get lidarr client configs for user: %w", err)
	}
	clientConfigs.Lidarr = lidarr
	return clientConfigs, nil
}
