// client.go
package emby

import (
	"context"
	"fmt"
	"strings"

	base "suasor/client"
	media "suasor/client/media"
	types "suasor/client/media/types"
	config "suasor/client/types"
	embyclient "suasor/internal/clients/embyAPI"
	"suasor/utils"
)

// EmbyClient implements the MediaContentProvider interface
type EmbyClient struct {
	media.BaseMediaClient
	client *embyclient.APIClient
	config *config.EmbyConfig
}

// NewEmbyClient creates a new Emby client instance
func NewEmbyClient(ctx context.Context, clientID uint64, cfg config.EmbyConfig) (media.MediaClient, error) {

	// Create API client configuration
	apiConfig := embyclient.NewConfiguration()
	apiConfig.BasePath = cfg.BaseURL

	// Set up API key in default headers
	apiConfig.DefaultHeader = map[string]string{
		"X-Emby-Token": cfg.APIKey,
	}

	client := embyclient.NewAPIClient(apiConfig)

	embyClient := &EmbyClient{
		BaseMediaClient: media.BaseMediaClient{
			BaseClient: base.BaseClient{
				ClientID: clientID,
				Category: config.ClientTypeEmby.AsCategory(),
				Type:     config.ClientTypeEmby,
				Config:   &cfg,
			},
		},
		client: client,
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
// 	media.RegisterClient(config.MediaClientTypeEmby, NewEmbyClient)
// }

// Capability methods
func (e *EmbyClient) SupportsMovies() bool      { return true }
func (e *EmbyClient) SupportsTVShows() bool     { return true }
func (e *EmbyClient) SupportsMusic() bool       { return true }
func (e *EmbyClient) SupportsPlaylists() bool   { return true }
func (e *EmbyClient) SupportsCollections() bool { return true }
func (e *EmbyClient) SupportsHistory() bool     { return true }

// resolveUserID resolves the Emby user ID from username
func (e *EmbyClient) resolveUserID(ctx context.Context) error {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Str("username", e.config.Username).
		Msg("Resolving Emby user ID from username")

	// Get the list of public users
	users, resp, err := e.client.UserServiceApi.GetUsersPublic(ctx)
	if err != nil {
		log.Error().
			Err(err).
			Str("username", e.config.Username).
			Msg("Failed to fetch Emby users")
		return fmt.Errorf("failed to fetch users: %w", err)
	}

	log.Debug().
		Int("statusCode", resp.StatusCode).
		Int("userCount", len(users)).
		Msg("Retrieved public users from Emby")

	// Find the user with matching username
	for _, user := range users {
		if strings.EqualFold(user.Name, e.config.Username) {
			e.config.UserID = user.Id
			log.Info().
				Str("username", e.config.Username).
				Str("userID", e.config.UserID).
				Msg("Successfully resolved Emby user ID")
			return nil
		}
	}

	log.Warn().
		Str("username", e.config.Username).
		Msg("Could not find matching user in Emby")
	return fmt.Errorf("user '%s' not found in Emby", e.config.Username)
}

func (e *EmbyClient) getArtworkURLs(item *embyclient.BaseItemDto) types.Artwork {
	imageURLs := types.Artwork{}

	if item == nil {
		return imageURLs
	}

	baseURL := strings.TrimSuffix(e.config.BaseURL, "/")

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
