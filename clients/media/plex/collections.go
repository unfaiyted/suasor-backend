package plex

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"suasor/clients/media/types"
	mediatypes "suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/logger"

	"github.com/unfaiyted/plexgo/models/operations"
)

// GetCollections retrieves collections from a Plex server
func (c *PlexClient) GetCollections(ctx context.Context, options *mediatypes.QueryOptions) ([]*models.MediaItem[*mediatypes.Collection], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Msg("Retrieving collections from Plex server")

	// First, get a library section ID to query (movie section as default)
	sectionID, err := c.getLibrarySectionID(ctx, operations.GetAllLibrariesTypeMovie)
	if err != nil {
		// Try to get TV show section if movie section not found
		sectionID, err = c.getLibrarySectionID(ctx, operations.GetAllLibrariesTypeTvShow)
		if err != nil {
			log.Error().
				Err(err).
				Msg("No suitable library section found for collections")
			return nil, fmt.Errorf("no suitable library section found: %w", err)
		}
	}

	// Use the new GetAllCollections method from plexgo
	plexCollections, err := c.plexAPI.Collections.GetAllCollections(ctx, sectionID)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Msg("Failed to get collections from Plex")
		return nil, fmt.Errorf("failed to get collections: %w", err)
	}

	if len(plexCollections) == 0 {
		log.Info().
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Msg("No collections found in Plex")
		return nil, nil
	}

	// Convert Plex collections to MediaItems
	var collections []*models.MediaItem[*mediatypes.Collection]
	for _, plexCollection := range plexCollections {

		collection := mediatypes.NewCollection(&types.MediaDetails{
			Title:       plexCollection.Title,
			Description: plexCollection.Summary,
		})

		mediaItem, err := GetMediaItem[*types.Collection](ctx, c, collection, plexCollection.RatingKey)
		if err != nil {
			log.Error().
				Err(err).
				Str("collectionID", plexCollection.RatingKey).
				Msg("Failed to convert collection to MediaItem")
			continue
		}

		mediaItem.ID = c.GetClientID()
		mediaItem.Title = plexCollection.Title

		// Set sync clients
		syncClients := models.SyncClients{}
		syncClients.AddClient(c.GetClientID(), c.GetClientType(), plexCollection.RatingKey)
		mediaItem.SyncClients = syncClients

		collections = append(collections, mediaItem)
	}

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Int("collectionsReturned", len(collections)).
		Msg("Completed GetCollections request")

	return collections, nil
}

// GetCollection retrieves a single collection by ID
func (c *PlexClient) GetCollection(ctx context.Context, collectionID string) (*models.MediaItem[*mediatypes.Collection], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("collectionID", collectionID).
		Msg("Retrieving collection from Plex server")

	// Convert collectionID to integer for Plex
	collectionRatingKey, err := strconv.Atoi(collectionID)
	if err != nil {
		log.Error().
			Err(err).
			Str("collectionID", collectionID).
			Msg("Failed to convert collection ID to integer")
		return nil, fmt.Errorf("invalid collection ID: %w", err)
	}

	// Use the new GetCollection method from plexgo
	plexCollection, err := c.plexAPI.Collections.GetCollection(ctx, collectionRatingKey)
	if err != nil {
		log.Error().
			Err(err).
			Str("collectionID", collectionID).
			Msg("Failed to get collection from Plex")
		return nil, fmt.Errorf("failed to get collection: %w", err)
	}

	// Create collection from Plex data
	mediaCollection := types.NewCollection(&types.MediaDetails{
		Title:       plexCollection.Title,
		Description: plexCollection.Summary,
	})

	mediaItem, err := GetMediaItem[*types.Collection](ctx, c, mediaCollection, plexCollection.RatingKey)
	if err != nil {
		log.Error().
			Err(err).
			Str("collectionID", plexCollection.RatingKey).
			Msg("Failed to convert collection to MediaItem")
		return nil, fmt.Errorf("failed to convert collection to MediaItem: %w", err)
	}

	// Set sync clients
	syncClients := models.SyncClients{}
	syncClients.AddClient(c.GetClientID(), c.GetClientType(), plexCollection.RatingKey)
	mediaItem.SyncClients = syncClients

	// Set collection items if available
	if plexCollection.ChildCount > 0 {
		// Get collection items
		collectionItemList, err := c.GetCollectionItems(ctx, plexCollection.RatingKey)
		if err != nil {
			log.Error().
				Err(err).
				Str("collectionID", plexCollection.RatingKey).
				Msg("Failed to get collection items")
		} else {
			setMediaListToItemList(collectionItemList, &mediaItem.Data.ItemList)
		}
	}

	log.Info().
		Str("collectionID", collectionID).
		Str("title", mediaItem.Title).
		Msg("Retrieved collection from Plex")

	return mediaItem, nil
}

// GetCollectionItems retrieves items in a collection
func (c *PlexClient) GetCollectionItems(ctx context.Context, collectionID string) (*models.MediaItemList, error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("collectionID", collectionID).
		Msg("Retrieving collection items from Plex server")

	// Convert collectionID to integer for Plex
	collectionRatingKey, err := strconv.Atoi(collectionID)
	if err != nil {
		log.Error().
			Err(err).
			Str("collectionID", collectionID).
			Msg("Failed to convert collection ID to integer")
		return nil, fmt.Errorf("invalid collection ID: %w", err)
	}

	// Use the new GetCollectionItems method from plexgo
	itemIDs, err := c.plexAPI.Collections.GetCollectionItems(ctx, collectionRatingKey)
	if err != nil {
		log.Error().
			Err(err).
			Str("collectionID", collectionID).
			Msg("Failed to get collection items from Plex")
		return nil, fmt.Errorf("failed to get collection items: %w", err)
	}

	if len(itemIDs) == 0 {
		log.Info().
			Str("collectionID", collectionID).
			Msg("Collection contains no items")
		return models.NewMediaItemList[*types.Collection](c.GetClientID(), 0), nil
	}

	// Create MediaItemList
	itemList := models.NewMediaItemList[*types.Collection](c.GetClientID(), 0)

	// Get details for each item
	for _, itemID := range itemIDs {
		// Convert itemID to int
		itemRatingKey, err := strconv.Atoi(itemID)
		if err != nil {
			log.Warn().
				Err(err).
				Str("itemID", itemID).
				Msg("Failed to convert item ID to integer, skipping")
			continue
		}

		// Get metadata for this item
		res, err := c.plexAPI.Library.GetMediaMetaData(ctx, operations.GetMediaMetaDataRequest{
			RatingKey: int64(itemRatingKey),
		})
		if err != nil {
			log.Warn().
				Err(err).
				Str("itemID", itemID).
				Msg("Failed to get item details, skipping")
			continue
		}

		if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil || len(res.Object.MediaContainer.Metadata) == 0 {
			log.Warn().
				Str("itemID", itemID).
				Msg("Empty response for item, skipping")
			continue
		}

		// Get the metadata for the item
		item := res.Object.MediaContainer.Metadata[0]
		if &item.Type == nil {
			log.Warn().
				Str("itemID", itemID).
				Msg("Item has no type, skipping")
			continue
		}

		// Determine media type based on Plex media type
		switch item.Type {
		case operations.GetMediaMetaDataTypeMovie:
			movie, err := GetItemFromMetadata[*types.Movie](ctx, c, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", itemID).
					Msg("Failed to convert item to movie, skipping")
				continue
			}
			mediaItem, err := GetChildMediaItem[*types.Movie](ctx, c, movie, item.RatingKey)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", itemID).
					Msg("Failed to create media item for movie, skipping")
				continue
			}

			itemList.AddMovie(mediaItem)
		case "show":
			show, err := GetItemFromMetadata[*types.Series](ctx, c, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", itemID).
					Msg("Failed to convert item to show, skipping")
				continue
			}
			mediaItem, err := GetChildMediaItem[*types.Series](ctx, c, show, item.RatingKey)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", itemID).
					Msg("Failed to create media item for show, skipping")
				continue
			}

			itemList.AddSeries(mediaItem)
		case "episode":
			episode, err := GetItemFromMetadata[*types.Episode](ctx, c, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", itemID).
					Msg("Failed to convert item to episode, skipping")
				continue
			}
			mediaItem, err := GetChildMediaItem[*types.Episode](ctx, c, episode, item.RatingKey)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", itemID).
					Msg("Failed to create media item for episode, skipping")
				continue
			}
			itemList.AddEpisode(mediaItem)
		case "track":
			track, err := GetItemFromMetadata[*types.Track](ctx, c, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", itemID).
					Msg("Failed to convert item to track, skipping")
				continue
			}
			mediaItem, err := GetChildMediaItem[*types.Track](ctx, c, track, item.RatingKey)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", itemID).
					Msg("Failed to create media item for track, skipping")
				continue
			}
			itemList.AddTrack(mediaItem)
		case "album":
			album, err := GetItemFromMetadata[*types.Album](ctx, c, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", itemID).
					Msg("Failed to convert item to album, skipping")
				continue
			}
			mediaItem, err := GetChildMediaItem[*types.Album](ctx, c, album, item.RatingKey)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", itemID).
					Msg("Failed to create media item for album, skipping")
				continue
			}
			itemList.AddAlbum(mediaItem)
		case "artist":
			artist, err := GetItemFromMetadata[*types.Artist](ctx, c, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", itemID).
					Msg("Failed to convert item to artist, skipping")
				continue
			}
			mediaItem, err := GetMediaItem[*types.Artist](ctx, c, artist, item.RatingKey)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", itemID).
					Msg("Failed to create media item for artist, skipping")
				continue
			}
			itemList.AddArtist(mediaItem)
		default:
			log.Warn().
				Str("collectionID", collectionID).
				Str("itemID", itemID).
				Str("type", string(item.Type)).
				Msg("Unknown media type in collection")
		}
	}

	log.Info().
		Str("collectionID", collectionID).
		Int("itemCount", itemList.GetTotalItems()).
		Msg("Retrieved collection items from Plex")

	return itemList, nil
}

// CreateCollection creates a new collection
func (c *PlexClient) CreateCollection(ctx context.Context, name string, description string) (*models.MediaItem[*mediatypes.Collection], error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("name", name).
		Msg("Creating collection in Plex server")

	// Get the movie library section ID (type 1 is movie)
	sectionID, err := c.getLibrarySectionID(ctx, operations.GetAllLibrariesTypeMovie)
	if err != nil {
		// Try to get TV show section if movie section not found
		sectionID, err = c.getLibrarySectionID(ctx, operations.GetAllLibrariesTypeTvShow)
		if err != nil {
			log.Error().
				Err(err).
				Msg("No suitable library section found for collection")
			return nil, fmt.Errorf("no suitable library section found: %w", err)
		}
	}

	// Create an empty collection
	plexCollection, err := c.plexAPI.Collections.CreateCollection(ctx, sectionID, name, []string{})
	if err != nil {
		log.Error().
			Err(err).
			Str("name", name).
			Int("sectionID", sectionID).
			Msg("Failed to create collection in Plex")
		return nil, fmt.Errorf("failed to create collection: %w", err)
	}

	collection := mediatypes.NewCollection(&types.MediaDetails{
		Title:       plexCollection.Title,
		Description: description, // Add the description even though Plex might not use it
	})

	mediaItem, err := GetMediaItem[*types.Collection](ctx, c, collection, plexCollection.RatingKey)
	if err != nil {
		log.Error().
			Err(err).
			Str("collectionID", plexCollection.RatingKey).
			Msg("Failed to convert collection to MediaItem")
		return nil, fmt.Errorf("failed to convert collection to MediaItem: %w", err)
	}
	// Create a MediaItem from the created collection

	log.Info().
		Str("name", name).
		Str("collectionID", plexCollection.RatingKey).
		Msg("Created collection in Plex")

	return mediaItem, nil
}

// CreateCollectionWithItems creates a new collection with items
func (c *PlexClient) CreateCollectionWithItems(ctx context.Context, name string, description string, itemIDs []string) (*models.MediaItem[*mediatypes.Collection], error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("name", name).
		Int("itemCount", len(itemIDs)).
		Msg("Creating collection with items in Plex server")

	if len(itemIDs) == 0 {
		return c.CreateCollection(ctx, name, description)
	}

	// Get the movie library section ID (type 1 is movie)
	sectionID, err := c.getLibrarySectionID(ctx, operations.GetAllLibrariesTypeMovie)
	if err != nil {
		// Try to get TV show section if movie section not found
		sectionID, err = c.getLibrarySectionID(ctx, operations.GetAllLibrariesTypeTvShow)
		if err != nil {
			log.Error().
				Err(err).
				Msg("No suitable library section found for collection")
			return nil, fmt.Errorf("no suitable library section found: %w", err)
		}
	}

	// Create collection with items
	plexCollection, err := c.plexAPI.Collections.CreateCollection(ctx, sectionID, name, itemIDs)
	if err != nil {
		log.Error().
			Err(err).
			Str("name", name).
			Int("sectionID", sectionID).
			Int("itemCount", len(itemIDs)).
			Msg("Failed to create collection in Plex")
		return nil, fmt.Errorf("failed to create collection: %w", err)
	}

	// Create a MediaItem from the created collection
	collection := types.NewCollection(&types.MediaDetails{
		Title:       plexCollection.Title,
		Description: description, // Add the description even though Plex might not use it
	})

	// Create the media item
	mediaItem, err := GetMediaItem[*types.Collection](ctx, c, collection, plexCollection.RatingKey)
	if err != nil {
		log.Error().
			Err(err).
			Str("collectionID", plexCollection.RatingKey).
			Msg("Failed to convert collection to MediaItem")
		return nil, fmt.Errorf("failed to convert collection to MediaItem: %w", err)
	}

	log.Info().
		Str("name", name).
		Int("itemCount", len(itemIDs)).
		Msg("Created collection with items in Plex")

	return mediaItem, nil
}

// UpdateCollection updates a collection's metadata
func (c *PlexClient) UpdateCollection(ctx context.Context, collectionID string, name string, description string) (*models.MediaItem[*mediatypes.Collection], error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("collectionID", collectionID).
		Str("name", name).
		Msg("Updating collection in Plex server")

	// Verify collection exists
	_, err := c.GetCollection(ctx, collectionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find collection to update: %w", err)
	}

	// Update collection using the Plex API
	// Note: Plex doesn't provide a direct method for updating collection metadata
	// through the API wrapper we're using
	log.Warn().
		Str("collectionID", collectionID).
		Msg("Plex API does not fully support collection metadata updates through the API wrapper")

	// For a real implementation, we would need to use a custom HTTP request
	// to update the collection metadata
	// Since this is not available in the API wrapper, we'll return the existing collection

	// Try to get the updated collection
	return c.GetCollection(ctx, collectionID)
}

// DeleteCollection deletes a collection
func (c *PlexClient) DeleteCollection(ctx context.Context, collectionID string) error {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("collectionID", collectionID).
		Msg("Deleting collection from Plex server")

	// Convert collectionID to integer for Plex
	collectionRatingKey, err := strconv.Atoi(collectionID)
	if err != nil {
		log.Error().
			Err(err).
			Str("collectionID", collectionID).
			Msg("Failed to convert collection ID to integer")
		return fmt.Errorf("invalid collection ID: %w", err)
	}

	// Use the new DeleteCollection method from plexgo
	err = c.plexAPI.Collections.DeleteCollection(ctx, collectionRatingKey)
	if err != nil {
		log.Error().
			Err(err).
			Str("collectionID", collectionID).
			Msg("Failed to delete collection from Plex")
		return fmt.Errorf("failed to delete collection: %w", err)
	}

	log.Info().
		Str("collectionID", collectionID).
		Msg("Collection deleted from Plex")

	return nil
}

// AddCollectionItem adds an item to a collection
func (c *PlexClient) AddCollectionItem(ctx context.Context, collectionID string, itemID string) error {
	return c.AddCollectionItems(ctx, collectionID, []string{itemID})
}

// AddCollectionItems adds multiple items to a collection
func (c *PlexClient) AddCollectionItems(ctx context.Context, collectionID string, itemIDs []string) error {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("collectionID", collectionID).
		Int("itemCount", len(itemIDs)).
		Msg("Adding items to collection in Plex server")

	if len(itemIDs) == 0 {
		log.Warn().Msg("No items to add to collection")
		return nil
	}

	// Convert collectionID to integer for Plex
	collectionRatingKey, err := strconv.Atoi(collectionID)
	if err != nil {
		log.Error().
			Err(err).
			Str("collectionID", collectionID).
			Msg("Failed to convert collection ID to integer")
		return fmt.Errorf("invalid collection ID: %w", err)
	}

	// Use the new AddToCollection method from plexgo
	err = c.plexAPI.Collections.AddToCollection(ctx, collectionRatingKey, itemIDs)
	if err != nil {
		log.Error().
			Err(err).
			Str("collectionID", collectionID).
			Msg("Failed to add items to collection")
		return fmt.Errorf("failed to add items to collection: %w", err)
	}

	log.Info().
		Str("collectionID", collectionID).
		Int("itemCount", len(itemIDs)).
		Msg("Items added to collection")

	return nil
}

// RemoveCollectionItem removes an item from a collection
func (c *PlexClient) RemoveCollectionItem(ctx context.Context, collectionID string, itemID string) error {
	return c.RemoveCollectionItems(ctx, collectionID, []string{itemID})
}

// RemoveCollectionItems removes multiple items from a collection
func (c *PlexClient) RemoveCollectionItems(ctx context.Context, collectionID string, itemIDs []string) error {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("collectionID", collectionID).
		Int("itemCount", len(itemIDs)).
		Msg("Removing items from collection in Plex server")

	if len(itemIDs) == 0 {
		log.Warn().Msg("No items to remove from collection")
		return nil
	}

	// Convert collectionID to integer for Plex
	collectionRatingKey, err := strconv.Atoi(collectionID)
	if err != nil {
		log.Error().
			Err(err).
			Str("collectionID", collectionID).
			Msg("Failed to convert collection ID to integer")
		return fmt.Errorf("invalid collection ID: %w", err)
	}

	// Use the new RemoveFromCollection method from plexgo
	err = c.plexAPI.Collections.RemoveFromCollection(ctx, collectionRatingKey, itemIDs)
	if err != nil {
		log.Error().
			Err(err).
			Str("collectionID", collectionID).
			Msg("Failed to remove items from collection")
		return fmt.Errorf("failed to remove items from collection: %w", err)
	}

	log.Info().
		Str("collectionID", collectionID).
		Int("itemCount", len(itemIDs)).
		Msg("Items removed from collection")

	return nil
}

// RemoveAllCollectionItems removes all items from a collection
func (c *PlexClient) RemoveAllCollectionItems(ctx context.Context, collectionID string) error {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("collectionID", collectionID).
		Msg("Removing all items from collection in Plex server")

	// Get all items in the collection
	itemList, err := c.GetCollectionItems(ctx, collectionID)
	if err != nil {
		log.Error().
			Err(err).
			Str("collectionID", collectionID).
			Msg("Failed to get collection items to remove")
		return fmt.Errorf("failed to get collection items: %w", err)
	}

	if itemList == nil || itemList.GetTotalItems() == 0 {
		log.Info().
			Str("collectionID", collectionID).
			Msg("No items to remove from collection")
		return nil
	}

	// Extract item IDs
	var itemIDs []string

	itemList.ForEach(func(uuid string, mediaType mediatypes.MediaType, item any) bool {
		itemIDs = append(itemIDs, uuid)
		return true
	})

	// Remove all items from the collection
	return c.RemoveCollectionItems(ctx, collectionID, itemIDs)
}

// ReorderCollectionItems reorders items in a collection
func (c *PlexClient) ReorderCollectionItems(ctx context.Context, collectionID string, itemIDs []string) error {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("collectionID", collectionID).
		Int("itemCount", len(itemIDs)).
		Msg("Reordering items in collection in Plex server")

	// Plex doesn't support reordering collection items through the API
	log.Warn().
		Str("collectionID", collectionID).
		Msg("Plex API does not support reordering collection items")

	return fmt.Errorf("reordering collection items not supported by Plex API")
}

// SearchCollections searches collections
func (c *PlexClient) SearchCollections(ctx context.Context, options *mediatypes.QueryOptions) ([]*models.MediaItem[*mediatypes.Collection], error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Msg("Searching collections in Plex server")

	// Get all collections and filter by search query
	collections, err := c.GetCollections(ctx, options)
	if err != nil {
		return nil, err
	}

	// If no search query, return all collections
	if options == nil || options.Query == "" {
		return collections, nil
	}

	// Filter collections by search query
	query := strings.ToLower(options.Query)
	var filteredCollections []*models.MediaItem[*mediatypes.Collection]
	for _, collection := range collections {
		if strings.Contains(strings.ToLower(collection.Title), query) {
			filteredCollections = append(filteredCollections, collection)
		}
	}

	log.Info().
		Str("query", options.Query).
		Int("resultCount", len(filteredCollections)).
		Msg("Filtered collections by search query")

	return filteredCollections, nil
}

// SearchCollectionItems searches items in a collection
func (c *PlexClient) SearchCollectionItems(ctx context.Context, collectionID string, options *mediatypes.QueryOptions) ([]*models.MediaItem[*mediatypes.Collection], error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("collectionID", collectionID).
		Msg("Searching items in collection in Plex server")

	// Plex doesn't have a direct API for searching within collection items
	log.Warn().
		Str("collectionID", collectionID).
		Msg("Plex API does not support searching within collection items")

	return nil, fmt.Errorf("searching within collection items not supported by Plex API")
}

// SupportsCollections returns whether the client supports collections
func (c *PlexClient) SupportsCollections() bool {
	return true
}

// getLibrarySectionID returns the first library section ID of the specified type
// type corresponds to the Plex library type: 1=movie, 2=show, 8=music, etc.
func (c *PlexClient) getLibrarySectionID(ctx context.Context, libraryType operations.GetAllLibrariesType) (int, error) {
	log := logger.LoggerFromContext(ctx)

	// Get all library sections
	res, err := c.plexAPI.Library.GetAllLibraries(ctx)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to get library sections from Plex")
		return 0, fmt.Errorf("failed to get library sections: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Directory == nil {
		return 0, fmt.Errorf("no library sections found")
	}

	// Find the first section of the specified type
	for _, dir := range res.Object.MediaContainer.Directory {
		if dir.Type == libraryType && &dir.Key != nil {
			sectionID, err := strconv.Atoi(dir.Key)
			if err != nil {
				log.Error().
					Err(err).
					Str("sectionKey", dir.Key).
					Msg("Failed to convert section key to integer")
				continue
			}
			log.Debug().
				Int("sectionID", sectionID).
				Str("title", dir.Title).
				Msg("Found library section")
			return sectionID, nil
		}
	}

	return 0, fmt.Errorf("no library section of type %d found", libraryType)
}
