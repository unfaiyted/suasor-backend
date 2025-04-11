package plex

import (
	"context"
	"fmt"
	"suasor/client/media/types"
	"suasor/types/models"
	"suasor/utils"
)

// GetCollections retrieves collections from Plex
func (c *PlexClient) GetCollections(ctx context.Context, options *types.QueryOptions) ([]models.MediaItem[*types.Collection], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Msg("Retrieving collections from Plex server")

	log.Debug().Msg("Making API request to Plex server for collections")
	res, err := c.plexAPI.Library.GetAllLibraries(ctx)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("Failed to get collections from Plex")
		return nil, fmt.Errorf("failed to get collections: %w", err)
	}

	directories := res.Object.MediaContainer.GetDirectory()

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("totalDirectories", len(directories)).
		Msg("Successfully retrieved library directories from Plex")

	collections := make([]models.MediaItem[*types.Collection], 0, len(directories))

	for _, dir := range directories {
		collection := models.MediaItem[*types.Collection]{
			Data: &types.Collection{
				Details: types.MediaDetails{
					Title: dir.Title,
					Artwork: types.Artwork{
						Thumbnail: c.makeFullURL(dir.Thumb),
					},
					ExternalIDs: types.ExternalIDs{types.ExternalID{
						Source: "plex",
						ID:     dir.Key,
					}},
				}},
		}

		collection.SetClientInfo(c.ClientID, c.ClientType, dir.Key)

		collections = append(collections, collection)

		log.Debug().
			Str("collectionID", dir.Key).
			Str("collectionName", dir.Title).
			Msg("Added collection to result list")
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("collectionsReturned", len(collections)).
		Msg("Completed GetCollections request")

	return collections, nil
}
