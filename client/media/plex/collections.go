package plex

import (
	"context"
	"fmt"
	mediatypes "suasor/client/media/types"
	"suasor/types/models"
	"suasor/utils"

	"github.com/LukeHagar/plexgo/models/operations"
)

// GetCollections retrieves collections from a Plex server
func (c *PlexClient) GetCollections(ctx context.Context, options *mediatypes.QueryOptions) ([]models.MediaItem[*mediatypes.Collection], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Msg("Retrieving collections from Plex server")

	request := operations.GetLibraryItemsRequest{
		IncludeMeta: operations.GetLibraryItemsQueryParamIncludeMetaEnable.ToPointer(),
		Tag:         "collection",
	}
	// Make API call to get collections
	// For Plex, collections are directories with type="collection"
	res, err := c.plexAPI.Library.GetLibraryItems(ctx, request)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("Failed to get collections from Plex")
		return nil, fmt.Errorf("failed to get collections: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil {
		log.Info().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("No collections found in Plex")
		return []models.MediaItem[*mediatypes.Collection]{}, nil
	}

	// Convert Plex directories to Collection models
	collections := make([]models.MediaItem[*mediatypes.Collection], 0, len(res.Object.MediaContainer.Metadata))
	for _, dir := range res.Object.MediaContainer.Metadata {

		collection := models.MediaItem[*mediatypes.Collection]{
			Data: &mediatypes.Collection{
				ItemList: mediatypes.ItemList{
					Details: mediatypes.MediaDetails{
						Title: dir.Title,
						Artwork: mediatypes.Artwork{
							Thumbnail: c.makeFullURL(*dir.Thumb),
						},
						ExternalIDs: mediatypes.ExternalIDs{mediatypes.ExternalID{
							Source: "plex",
							ID:     dir.Key,
						}},
					},
				},
			},
			Type: mediatypes.MediaTypeCollection,
		}
		collection.SetClientInfo(c.ClientID, c.ClientType, dir.Key)

		collections = append(collections, collection)
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("collectionsReturned", len(collections)).
		Msg("Completed GetCollections request")

	return collections, nil
}

// Function to get a single collection by ID
func (c *PlexClient) GetCollectionByID(ctx context.Context, collectionID string) (*models.MediaItem[*mediatypes.Collection], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("collectionID", collectionID).
		Msg("Retrieving collection from Plex server")

	// TODO: Implement fetching a single collection from Plex
	// This would typically involve:
	// 1. Making an API call to fetch the collection by ID
	// 2. Converting the response to a Collection model
	// 3. Returning the model

	return nil, fmt.Errorf("GetCollectionByID not implemented for Plex")
}
