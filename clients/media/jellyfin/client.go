// client.go
package jellyfin

import (
	"context"
	"fmt"
	"strings"

	jellyfin "github.com/sj14/jellyfin-go/api"
	"suasor/clients/media"
	mediatypes "suasor/clients/media/types"
	clienttypes "suasor/clients/types"
	"suasor/types/responses"
	"suasor/utils/logger"
)

type JellyfinClient struct {
	media.ClientMedia
	client *jellyfin.APIClient
	config *clienttypes.JellyfinConfig
}

// NewJellyfinClient creates a new Jellyfin client instance
func NewJellyfinClient(ctx context.Context, registry *media.ClientItemRegistry, clientID uint64, cfg *clienttypes.JellyfinConfig) (media.ClientMedia, error) {
	// Create API client configuration
	apiConfig := &jellyfin.Configuration{
		Servers:       jellyfin.ServerConfigurations{{URL: cfg.GetBaseURL()}},
		DefaultHeader: map[string]string{"Authorization": fmt.Sprintf(`MediaBrowser Token="%s"`, cfg.GetAPIKey())},
	}

	client := jellyfin.NewAPIClient(apiConfig)

	clientMedia, err := media.NewClientMedia(ctx, clientID, clienttypes.ClientMediaTypeJellyfin, registry, cfg)
	if err != nil {
		return nil, err
	}

	jellyfinClient := &JellyfinClient{
		ClientMedia: clientMedia,
		client:      client,
		config:      cfg,
	}

	// Resolve user ID if username is provided
	if cfg.GetUsername() != "" && cfg.GetUserID() == "" {
		if err := jellyfinClient.resolveUserID(ctx); err != nil {
			// Log but don't fail - some operations might work without a user ID
			log := logger.LoggerFromContext(ctx)
			log.Warn().
				Err(err).
				Str("username", cfg.GetUsername()).
				Msg("Failed to resolve Jellyfin user ID, some operations may be limited")
		}
	}
	return jellyfinClient, nil
}

// Capability methods

func (j *JellyfinClient) SupportsMusic() bool   { return true }
func (j *JellyfinClient) SupportsHistory() bool { return true }

func (j *JellyfinClient) jellyfinConfig() *clienttypes.JellyfinConfig {
	cfg := j.GetConfig().(*clienttypes.JellyfinConfig)
	return cfg
}

// getUserID returns the Jellyfin user ID - either directly from config or resolved from username
func (j *JellyfinClient) getUserID() string {
	// Defensive programming - check for nil pointer
	if j == nil || j.jellyfinConfig() == nil {
		return ""
	}

	// Return existing user ID if available
	if j.jellyfinConfig().GetUserID() != "" {
		return j.jellyfinConfig().GetUserID()
	}

	// If no user ID is available, return empty string
	return ""
}

// resolveUserID resolves the user ID from the username
func (j *JellyfinClient) resolveUserID(ctx context.Context) error {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Str("username", j.jellyfinConfig().GetUsername()).
		Msg("Resolving Jellyfin user ID from username")

	// Get the list of public users
	publicUsersReq := j.client.UserAPI.GetUsers(ctx)
	users, resp, err := publicUsersReq.Execute()

	if err != nil {
		log.Error().
			Err(err).
			Str("username", j.jellyfinConfig().GetUsername()).
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
			if strings.EqualFold(*user.Name.Get(), j.jellyfinConfig().GetUsername()) {
				// TODO: Use proper setter method once added to JellyfinConfig
				j.jellyfinConfig().UserID = *user.Id
				log.Info().
					Str("username", j.jellyfinConfig().GetUsername()).
					Str("userID", j.jellyfinConfig().GetUserID()).
					Msg("Successfully resolved Jellyfin user ID")
				return nil
			}
		}
	}

	log.Warn().
		Str("username", j.jellyfinConfig().GetUsername()).
		Msg("Could not find matching user in Jellyfin")
	return fmt.Errorf("user '%s' not found in Jellyfin", j.jellyfinConfig().GetUsername())
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

// Search implements the Search method for JellyfinClient
// This overrides the default implementation to provide better search results
func (j *JellyfinClient) Search(ctx context.Context, options *mediatypes.QueryOptions) (responses.SearchResults, error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Str("query", options.Query).
		Str("mediaType", string(options.MediaType)).
		Msg("Searching media items in Jellyfin")

	var results responses.SearchResults

	// If no query, return empty results
	if options.Query == "" {
		return results, fmt.Errorf("search query cannot be empty")
	}

	// Check if mediaType is empty and if it is, set it to ALL
	if options.MediaType == "" {
		options.MediaType = mediatypes.MediaTypeAll
	}

	// Process the search based on media type
	switch options.MediaType {
	case mediatypes.MediaTypeMovie:
		movies, err := j.SearchMovies(ctx, options)
		if err != nil {
			log.Error().Err(err).Msg("Failed to search movies")
			return results, fmt.Errorf("failed to search through client: %w", err)
		}
		results.Movies = movies

	case mediatypes.MediaTypeSeries:
		// TODO: Implement series search when available
		return results, media.ErrFeatureNotSupported

	case mediatypes.MediaTypeTrack, mediatypes.MediaTypeAlbum, mediatypes.MediaTypeArtist:
		// TODO: Implement music search when available
		return results, media.ErrFeatureNotSupported

	case mediatypes.MediaTypePlaylist:
		// TODO: Implement playlist search when available
		return results, media.ErrFeatureNotSupported

	case mediatypes.MediaTypeCollection:
		// TODO: Implement collection search when available
		return results, media.ErrFeatureNotSupported

	default:
		// For generic searches, try all supported types
		// For now we just search movies since that's what we've implemented
		movies, err := j.SearchMovies(ctx, options)
		if err != nil {
			log.Error().Err(err).Msg("Failed to search movies")
			// Continue with other types even if movies fail
		} else {
			results.Movies = movies
		}

		// TODO: Add series, music, etc. when those search methods are implemented
	}

	return results, nil
}
