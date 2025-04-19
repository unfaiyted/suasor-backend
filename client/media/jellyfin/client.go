package jellyfin

import (
	"context"
	"fmt"
	"strings"

	jellyfin "github.com/sj14/jellyfin-go/api"
	"suasor/client"
	"suasor/client/media"
	"suasor/client/types"
	"suasor/utils"
)

type JellyfinClient struct {
	media.BaseClientMedia
	client *jellyfin.APIClient
	config types.JellyfinConfig
}

// NewJellyfinClient creates a new Jellyfin client instance
func NewJellyfinClient(ctx context.Context, registry *media.ClientItemRegistry, clientID uint64, config types.JellyfinConfig) (media.ClientMedia, error) {
	// Get or create registry for media item factories

	// Create API client configuration
	apiConfig := &jellyfin.Configuration{
		Servers:       jellyfin.ServerConfigurations{{URL: config.BaseURL}},
		DefaultHeader: map[string]string{"Authorization": fmt.Sprintf(`MediaBrowser Token="%s"`, config.APIKey)},
	}

	jfClient := jellyfin.NewAPIClient(apiConfig)

	jellyfinClient := &JellyfinClient{
		BaseClientMedia: media.BaseClientMedia{
			ItemRegistry: registry,
			ClientType:   types.ClientMediaTypeJellyfin,
			BaseClient: client.BaseClient{
				ClientID: clientID,
				Category: types.ClientMediaTypeJellyfin.AsCategory(),
				Type:     types.ClientTypeJellyfin,
				Config:   &config,
			},
		},
		client: jfClient,
	}

	// Resolve user ID if username is provided
	if config.Username != "" && config.UserID == "" {
		if err := jellyfinClient.resolveUserID(ctx); err != nil {
			// Log but don't fail - some operations might work without a user ID
			log := utils.LoggerFromContext(ctx)
			log.Warn().
				Err(err).
				Str("username", config.Username).
				Msg("Failed to resolve Jellyfin user ID, some operations may be limited")
		}
	}
	return jellyfinClient, nil
}

// Register the client factory
// func init() {
// 	media.RegisterClient(types.ClientMediaTypeJellyfin, NewJellyfinClient)
// }

// Capability methods
func (j *JellyfinClient) SupportsMovies() bool  { return true }
func (j *JellyfinClient) SupportsTVShows() bool { return true }
func (j *JellyfinClient) SupportsMusic() bool   { return true }
func (j *JellyfinClient) SupportsHistory() bool { return true }

func (j *JellyfinClient) GetRegistry() *media.ClientItemRegistry {
	return j.ItemRegistry
}

// resolveUserID resolves the user ID from the username
func (j *JellyfinClient) resolveUserID(ctx context.Context) error {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Str("username", j.config.Username).
		Msg("Resolving Jellyfin user ID from username")

		// Get the list of public users
	publicUsersReq := j.client.UserAPI.GetUsers(ctx)
	users, resp, err := publicUsersReq.Execute()

	if err != nil {
		log.Error().
			Err(err).
			Str("username", j.config.Username).
			Msg("Failed to fetch Jellyfin users")
		return fmt.Errorf("failed to fetch users: %w", err)
	}

	log.Debug().
		Int("statusCode", resp.StatusCode).
		Int("userCount", len(users)).
		Msg("Retrieved public users from Jellyfin")

	// Find the user with matching username
	for _, user := range users {
		if user.Name.IsSet() {
			if strings.EqualFold(*user.Name.Get(), j.config.Username) {
				j.config.UserID = *user.Id
				log.Info().
					Str("username", j.config.Username).
					Str("userID", j.config.UserID).
					Msg("Successfully resolved Jellyfin user ID")
				return nil
			}
		}
	}

	log.Warn().
		Str("username", j.config.Username).
		Msg("Could not find matching user in Jellyfin")
	return fmt.Errorf("user '%s' not found in Jellyfin", j.config.Username)
}

func (j *JellyfinClient) TestConnection(ctx context.Context) (bool, error) {
	sysInfo, _, err := j.client.SystemAPI.GetSystemInfo(ctx).Execute()
	if err != nil {
		return false, err
	}
	if *sysInfo.Version.Get() == "" {
		return false, fmt.Errorf("failed to retrieve Jellyfin server version")
	}
	return true, nil
}
