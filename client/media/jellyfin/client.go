package jellyfin

import (
	"context"
	"fmt"
	"strings"

	jellyfin "github.com/sj14/jellyfin-go/api"
	base "suasor/client"
	media "suasor/client/media"
	t "suasor/client/media/types"
	client "suasor/client/types"
	"suasor/utils"
)

// JellyfinClient implements the MediaContentProvider interface
type JellyfinClient struct {
	base.BaseMediaClient
	client *jellyfin.APIClient
	config client.JellyfinConfig
}

// NewJellyfinClient creates a new Jellyfin client instance
func NewJellyfinClient(ctx context.Context, clientID uint64, config client.ClientConfig) (media.MediaClient, error) {
	cfg, ok := config.(client.JellyfinConfig)
	if !ok {
		return nil, fmt.Errorf("invalid configuration for Jellyfin client")
	}

	// Create API client configuration
	apiConfig := &jellyfin.Configuration{
		Servers:       jellyfin.ServerConfigurations{{URL: cfg.BaseURL}},
		DefaultHeader: map[string]string{"Authorization": fmt.Sprintf(`MediaBrowser Token="%s"`, cfg.APIKey)},
	}

	jfClient := jellyfin.NewAPIClient(apiConfig)

	jellyfinClient := &JellyfinClient{
		BaseMediaClient: base.BaseMediaClient{
			ClientID:   clientID,
			ClientType: client.MediaClientTypeJellyfin,
		},
		client: jfClient,
		config: cfg,
	}

	// Resolve user ID if username is provided
	if cfg.Username != "" && cfg.UserID == "" {
		if err := jellyfinClient.resolveUserID(ctx); err != nil {
			// Log but don't fail - some operations might work without a user ID
			log := utils.LoggerFromContext(ctx)
			log.Warn().
				Err(err).
				Str("username", cfg.Username).
				Msg("Failed to resolve Jellyfin user ID, some operations may be limited")
		}
	}
	return jellyfinClient, nil
}

// Register the provider factory
func init() {
	media.RegisterProvider(client.MediaClientTypeJellyfin, NewJellyfinClient)
}

// Capability methods
func (j *JellyfinClient) SupportsMovies() bool      { return true }
func (j *JellyfinClient) SupportsTVShows() bool     { return true }
func (j *JellyfinClient) SupportsMusic() bool       { return true }
func (j *JellyfinClient) SupportsPlaylists() bool   { return true }
func (j *JellyfinClient) SupportsCollections() bool { return true }

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

// GetMovies retrieves movies from the Jellyfin server

// getQueryParameters extracts common query parameters from QueryOptions
// and converts them to the format expected by the Jellyfin API
func (j *JellyfinClient) getQueryParameters(options *t.QueryOptions) (limit, startIndex *int32, sortBy []jellyfin.ItemSortBy, sortOrder []jellyfin.SortOrder) {

	// Default values
	defaultLimit := int32(100)
	defaultOffset := int32(0)
	limit = &defaultLimit
	startIndex = &defaultOffset

	if options != nil {
		if options.Limit > 0 {
			limitVal := int32(options.Limit)
			limit = &limitVal
		}
		if options.Offset > 0 {
			offsetVal := int32(options.Offset)
			startIndex = &offsetVal
		}
		if options.Sort != "" {
			// sortBy = &options.Sort
			sortBy = []jellyfin.ItemSortBy{jellyfin.ItemSortBy(options.Sort)}
			if options.SortOrder == "desc" {
				sortOrder = []jellyfin.SortOrder{jellyfin.SORTORDER_DESCENDING}
			} else {
				sortOrder = []jellyfin.SortOrder{jellyfin.SORTORDER_ASCENDING}
			}
		}
	}
	return
}

func (j *JellyfinClient) getArtworkURLs(item *jellyfin.BaseItemDto) t.Artwork {
	imageURLs := t.Artwork{}

	if item == nil || item.Id == nil {
		return imageURLs
	}

	baseURL := strings.TrimSuffix(j.config.BaseURL, "/")
	itemID := *item.Id

	// Primary image (poster)
	if item.ImageTags != nil {
		if tag, ok := (item.ImageTags)["Primary"]; ok {
			imageURLs.Poster = fmt.Sprintf("%s/Items/%s/Images/Primary?tag=%s", baseURL, itemID, tag)
		}
	}

	// Backdrop image
	if item.BackdropImageTags != nil && len(item.BackdropImageTags) > 0 {
		imageURLs.Background = fmt.Sprintf("%s/Items/%s/Images/Backdrop?tag=%s", baseURL, itemID, item.BackdropImageTags[0])
	}

	// Other image types
	if item.ImageTags != nil {
		if tag, ok := (item.ImageTags)["Logo"]; ok {
			imageURLs.Logo = fmt.Sprintf("%s/Items/%s/Images/Logo?tag=%s", baseURL, itemID, tag)
		}

		if tag, ok := (item.ImageTags)["Thumb"]; ok {
			imageURLs.Thumbnail = fmt.Sprintf("%s/Items/%s/Images/Thumb?tag=%s", baseURL, itemID, tag)
		}

		if tag, ok := (item.ImageTags)["Banner"]; ok {
			imageURLs.Banner = fmt.Sprintf("%s/Items/%s/Images/Banner?tag=%s", baseURL, itemID, tag)
		}
	}

	return imageURLs
}
