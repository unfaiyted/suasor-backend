// client.go
package emby

import (
	"context"
	"fmt"
	"strings"

	"suasor/client"
	base "suasor/client"
	media "suasor/client/media"
	types "suasor/client/media/types"
	config "suasor/client/types"
	embyclient "suasor/internal/clients/embyAPI"
	"suasor/utils"
)

// Add this init function to register the Emby client factory
func init() {

	fmt.Printf("Registering factory for client type: %s (value: %v)\n",
		config.ClientTypeEmby.String(), config.ClientTypeEmby)

	fmt.Println("Registering Emby client factory...")
	client.GetClientFactoryService().RegisterClientFactory(config.ClientTypeEmby,
		func(ctx context.Context, clientID uint64, configData config.ClientConfig) (base.Client, error) {
			log := utils.LoggerFromContext(ctx)
			// Use the provided config (should be an EmbyConfig)
			embyConfig, ok := configData.(*config.EmbyConfig)
			log.Debug().
				Bool("ok", ok).
				Msg("Checking config type")
			if !ok {
				log.Error().
					Err(fmt.Errorf("expected *config.EmbyConfig, got %T", configData)).
					Msg("Expected *config.EmbyConfig, got %T")
				return nil, fmt.Errorf("expected *config.EmbyConfig, got %T", configData)
			}

			fmt.Printf("Factory called for Emby client with ID: %d\n", clientID)
			return NewEmbyClient(ctx, clientID, *embyConfig)
		})
}

// EmbyClient implements the MediaContentProvider interface
type EmbyClient struct {
	media.BaseClientMedia
	client *embyclient.APIClient
	// config *config.EmbyConfig
}

// NewEmbyClient creates a new Emby client instance
func NewEmbyClient(ctx context.Context, clientID uint64, cfg config.EmbyConfig) (media.ClientMedia, error) {

	// Create API client configuration
	apiConfig := embyclient.NewConfiguration()
	apiConfig.BasePath = cfg.BaseURL

	// Set up API key in default headers
	apiConfig.DefaultHeader = map[string]string{
		"X-Emby-Token": cfg.APIKey,
	}

	client := embyclient.NewAPIClient(apiConfig)

	embyClient := &EmbyClient{
		BaseClientMedia: media.BaseClientMedia{
			BaseClient: base.BaseClient{
				ClientID: clientID,
				Category: config.ClientTypeEmby.AsCategory(),
				Type:     config.ClientTypeEmby,
				Config:   &cfg,
			},
		},
		client: client,
		// config: &cfg,
	}
	// Resolve user ID if username is provided
	if cfg.Username != "" && cfg.UserID == "" {
		if err := embyClient.resolveUserID(ctx); err != nil {
			// Log but don't fail - some operations might work without a user ID
			log := utils.LoggerFromContext(ctx)
			log.Warn().
				Err(err).
				Str("username", cfg.Username).
				Msg("Failed to resolve Emby user ID, some operations may be limited")
		}
	}
	return embyClient, nil
}

// Register the provider factory
// func init() {
// 	media.RegisterClient(config.ClientMediaTypeEmby, NewEmbyClient)
// }

// Capability methods
func (e *EmbyClient) SupportsMovies() bool      { return true }
func (e *EmbyClient) SupportsTVShows() bool     { return true }
func (e *EmbyClient) SupportsMusic() bool       { return true }
func (e *EmbyClient) SupportsPlaylists() bool   { return true }
func (e *EmbyClient) SupportsCollections() bool { return true }
func (e *EmbyClient) SupportsHistory() bool     { return true }

func (e *EmbyClient) embyConfig() *config.EmbyConfig {
	cfg, _ := e.Config.(*config.EmbyConfig)
	return cfg
}

// getUserID returns the Emby user ID - either directly from config or resolved from username
func (e *EmbyClient) getUserID() string {
	if e.embyConfig() == nil {
		return ""
	}

	// Return existing user ID if available
	if e.embyConfig().UserID != "" {
		return e.embyConfig().UserID
	}

	// Try to infer it from username, but this won't work in this context
	// since we'd need to make API call which requires context
	// log error and return empty
	return ""
}

// resolveUserID resolves the Emby user ID from username
func (e *EmbyClient) resolveUserID(ctx context.Context) error {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Str("username", e.embyConfig().Username).
		Msg("Resolving Emby user ID from username")

	// Get the list of public users
	users, resp, err := e.client.UserServiceApi.GetUsersPublic(ctx)
	if err != nil {
		log.Error().
			Err(err).
			Str("username", e.embyConfig().Username).
			Msg("Failed to fetch Emby users")
		return fmt.Errorf("failed to fetch users: %w", err)
	}

	log.Debug().
		Int("statusCode", resp.StatusCode).
		Int("userCount", len(users)).
		Msg("Retrieved public users from Emby")

	// Find the user with matching username
	for _, user := range users {
		if strings.EqualFold(user.Name, e.embyConfig().Username) {
			e.embyConfig().UserID = user.Id
			log.Info().
				Str("username", e.embyConfig().Username).
				Str("userID", e.embyConfig().UserID).
				Msg("Successfully resolved Emby user ID")
			return nil
		}
	}

	log.Warn().
		Str("username", e.embyConfig().Username).
		Msg("Could not find matching user in Emby")
	return fmt.Errorf("user '%s' not found in Emby", e.embyConfig().Username)
}

func (e *EmbyClient) getArtworkURLs(item *embyclient.BaseItemDto) types.Artwork {
	imageURLs := types.Artwork{}

	if item == nil {
		return imageURLs
	}

	baseURL := strings.TrimSuffix(e.embyConfig().BaseURL, "/")

	// Primary image (poster) - with nil check
	if item.ImageTags != nil {
		if tag, ok := item.ImageTags["Primary"]; ok {
			imageURLs.Poster = fmt.Sprintf("%s/Items/%s/Images/Primary?tag=%s", baseURL, item.Id, tag)
		}
	}

	// Backdrop image - with nil and length check
	if item.BackdropImageTags != nil && len(item.BackdropImageTags) > 0 {
		imageURLs.Background = fmt.Sprintf("%s/Items/%s/Images/Backdrop?tag=%s", baseURL, item.Id, item.BackdropImageTags[0])
	}

	// Other image types - with nil check
	if item.ImageTags != nil {
		if tag, ok := item.ImageTags["Logo"]; ok {
			imageURLs.Logo = fmt.Sprintf("%s/Items/%s/Images/Logo?tag=%s", baseURL, item.Id, tag)
		}

		if tag, ok := item.ImageTags["Thumb"]; ok {
			imageURLs.Thumbnail = fmt.Sprintf("%s/Items/%s/Images/Thumb?tag=%s", baseURL, item.Id, tag)
		}

		if tag, ok := item.ImageTags["Banner"]; ok {
			imageURLs.Banner = fmt.Sprintf("%s/Items/%s/Images/Banner?tag=%s", baseURL, item.Id, tag)
		}
	}

	return imageURLs
}

func (c *EmbyClient) TestConnection(ctx context.Context) (bool, error) {
	sysInfo, _, err := c.client.SystemServiceApi.GetSystemInfo(ctx)
	if err != nil {
		return false, err
	}
	if sysInfo.Version == "" {
		return false, fmt.Errorf("failed to retrieve Emby server version")
	}
	return true, nil
}
