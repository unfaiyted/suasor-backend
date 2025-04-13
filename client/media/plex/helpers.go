package plex

import (
	"context"
	"fmt"
	"strings"
	"suasor/client/types"
	"suasor/utils"
)

// makeFullURL creates a complete URL from a resource path
func (c *PlexClient) makeFullURL(resourcePath string) string {
	if resourcePath == "" {
		return ""
	}

	plexConfig := c.Config.(*types.PlexConfig)

	if strings.HasPrefix(resourcePath, "http") {
		return resourcePath
	}

	return fmt.Sprintf("%s%s", plexConfig.BaseURL, resourcePath)
}

// // makeFullURL converts a Plex path to a full URL with authentication
// func (c *PlexClient) makeFullURL(path *string) string {
// 	if path == nil || *path == "" {
// 		return ""
// 	}
// 	return fmt.Sprintf("%s%s?X-Plex-Token=%s", c.BaseURL, *path, c.Token)
// }

// findLibrarySectionByType returns the section key for the specified type
func (c *PlexClient) findLibrarySectionByType(ctx context.Context, sectionType string) (string, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Debug().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("sectionType", sectionType).
		Msg("Finding library section by type")

	libraries, err := c.plexAPI.Library.GetAllLibraries(ctx)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("sectionType", sectionType).
			Msg("Failed to get libraries from Plex")
		return "", fmt.Errorf("failed to get libraries: %w", err)
	}

	log.Debug().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("libraryCount", len(libraries.Object.MediaContainer.Directory)).
		Msg("Retrieved libraries from Plex")

	for _, dir := range libraries.Object.MediaContainer.Directory {
		if dir.Type == sectionType {
			log.Debug().
				Uint64("clientID", c.ClientID).
				Str("clientType", string(c.ClientType)).
				Str("sectionType", sectionType).
				Str("sectionKey", dir.Key).
				Str("sectionTitle", dir.Title).
				Msg("Found matching library section")
			return dir.Key, nil
		}
	}

	log.Debug().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("sectionType", sectionType).
		Msg("No matching library section found")

	return "", nil
}
