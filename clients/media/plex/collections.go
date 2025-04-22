package plex

import (
	"context"
	"fmt"
	mediatypes "suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/logger"

	"github.com/LukeHagar/plexgo/models/operations"
)

// GetCollections retrieves collections from a Plex server
func (c *PlexClient) GetCollections(ctx context.Context, options *mediatypes.QueryOptions) ([]*models.MediaItem[*mediatypes.Collection], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

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
		return nil, nil
	}

	collections, err := GetMediaItemList[*mediatypes.Collection](ctx, c, res.Object.MediaContainer.Metadata)

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
	log := logger.LoggerFromContext(ctx)

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
