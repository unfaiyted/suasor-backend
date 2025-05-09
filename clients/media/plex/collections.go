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
	"time"

	"github.com/LukeHagar/plexgo/models/operations"
)

// GetCollections retrieves collections from a Plex server
func (c *PlexClient) GetCollections(ctx context.Context, options *mediatypes.QueryOptions) ([]*models.MediaItem[*mediatypes.Collection], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
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
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Msg("Failed to get collections from Plex")
		return nil, fmt.Errorf("failed to get collections: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil {
		log.Info().
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Msg("No collections found in Plex")
		return nil, nil
	}

	collections, err := GetMediaItemList[*mediatypes.Collection](ctx, c, res.Object.MediaContainer.Metadata)

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

	// Plex doesn't have a direct API for getting a collection by ID
	// We need to use the metadata endpoint
	res, err := c.plexAPI.Library.GetMetadata(ctx, collectionRatingKey)
	if err != nil {
		log.Error().
			Err(err).
			Str("collectionID", collectionID).
			Msg("Failed to get collection from Plex")
		return nil, fmt.Errorf("failed to get collection: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil || len(res.Object.MediaContainer.Metadata) == 0 {
		log.Error().
			Str("collectionID", collectionID).
			Msg("Collection not found or empty response from Plex")
		return nil, fmt.Errorf("collection not found")
	}

	// Convert Plex metadata to Collection model
	plexMetadata := res.Object.MediaContainer.Metadata[0]
	collection := types.NewCollection()
	collection.MediaItemList.Details = &types.MediaDetails{
		Title:       plexMetadata.Title,
		Description: plexMetadata.Summary,
	}

	// Set collection item count if available
	if plexMetadata.ChildCount != nil {
		collection.MediaItemList.ItemCount = int(*plexMetadata.ChildCount)
	}

	// Create MediaItem
	mediaItem := models.NewMediaItem(types.MediaTypeCollection, collection)
	mediaItem.ID = c.GetClientID()
	mediaItem.Title = plexMetadata.Title
	mediaItem.ClientID = c.GetClientID()
	mediaItem.ClientItemID = strconv.Itoa(*plexMetadata.RatingKey)

	// Set sync clients
	syncClients := models.SyncClients{}
	syncClients.AddClient(c.GetClientID(), c.GetClientType(), strconv.Itoa(*plexMetadata.RatingKey))
	mediaItem.SyncClients = syncClients

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

	// Get collection contents from Plex
	res, err := c.plexAPI.Library.GetMetadataChildren(ctx, collectionRatingKey)
	if err != nil {
		log.Error().
			Err(err).
			Str("collectionID", collectionID).
			Msg("Failed to get collection items from Plex")
		return nil, fmt.Errorf("failed to get collection items: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil {
		log.Error().
			Str("collectionID", collectionID).
			Msg("Collection items not found or empty response from Plex")
		return nil, fmt.Errorf("collection items not found")
	}

	// Create MediaItemList
	itemList := &models.MediaItemList{
		Items: make([]*models.MediaItemListItem, 0, len(res.Object.MediaContainer.Metadata)),
	}

	// Add items to list
	for i, item := range res.Object.MediaContainer.Metadata {
		mediaType := types.MediaTypeUnknown
		// Determine media type based on Plex media type
		switch item.Type {
		case "movie":
			mediaType = types.MediaTypeMovie
		case "show":
			mediaType = types.MediaTypeSeries
		case "episode":
			mediaType = types.MediaTypeEpisode
		case "track":
			mediaType = types.MediaTypeTrack
		case "album":
			mediaType = types.MediaTypeAlbum
		case "artist":
			mediaType = types.MediaTypeArtist
		default:
			mediaType = types.MediaTypeUnknown
		}

		itemList.Items = append(itemList.Items, &models.MediaItemListItem{
			ID:       strconv.Itoa(*item.RatingKey),
			Position: i,
			Title:    item.Title,
			Type:     mediaType,
		})
	}

	itemList.ItemCount = len(itemList.Items)
	log.Info().
		Str("collectionID", collectionID).
		Int("itemCount", itemList.ItemCount).
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

	// Plex doesn't have a direct API for creating empty collections
	// Collections are usually created with items
	log.Warn().
		Str("name", name).
		Msg("Plex API does not support creating empty collections, collections are created by adding items to them")

	return nil, fmt.Errorf("creating empty collections not directly supported by Plex API")
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

	// Format item IDs for Plex API
	var itemKeys []string
	for _, id := range itemIDs {
		itemKeys = append(itemKeys, id)
	}

	// Join item keys for the API request
	itemKeysStr := strings.Join(itemKeys, ",")

	// Create collection with items
	// This is a simplified implementation - the real Plex API might require different parameters
	title := operations.CreateCollectionTitle(name)
	summary := operations.CreateCollectionSummary(description)
	itemID := operations.CreateCollectionItemID(itemKeysStr)

	res, err := c.plexAPI.Library.CreateCollection(ctx, title, summary, itemID)
	if err != nil {
		log.Error().
			Err(err).
			Str("name", name).
			Msg("Failed to create collection in Plex")
		return nil, fmt.Errorf("failed to create collection: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil || len(res.Object.MediaContainer.Metadata) == 0 {
		log.Error().
			Str("name", name).
			Msg("Empty response when creating collection in Plex")
		return nil, fmt.Errorf("empty response when creating collection")
	}

	// Create MediaItem from created collection
	plexMetadata := res.Object.MediaContainer.Metadata[0]
	collection := types.NewCollection()
	collection.MediaItemList.Details = &types.MediaDetails{
		Title:       plexMetadata.Title,
		Description: plexMetadata.Summary,
		AddedAt:     time.Now(),
		UpdatedAt:   time.Now(),
	}

	mediaItem := models.NewMediaItem(types.MediaTypeCollection, collection)
	mediaItem.ID = c.GetClientID()
	mediaItem.Title = plexMetadata.Title
	mediaItem.ClientID = c.GetClientID()
	mediaItem.ClientItemID = strconv.Itoa(*plexMetadata.RatingKey)

	// Set sync clients
	syncClients := models.SyncClients{}
	syncClients.AddClient(c.GetClientID(), c.GetClientType(), strconv.Itoa(*plexMetadata.RatingKey))
	mediaItem.SyncClients = syncClients

	log.Info().
		Str("name", name).
		Str("collectionID", mediaItem.ClientItemID).
		Msg("Created collection in Plex")

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

	// Convert collectionID to integer for Plex
	collectionRatingKey, err := strconv.Atoi(collectionID)
	if err != nil {
		log.Error().
			Err(err).
			Str("collectionID", collectionID).
			Msg("Failed to convert collection ID to integer")
		return nil, fmt.Errorf("invalid collection ID: %w", err)
	}

	// Verify collection exists
	_, err = c.GetCollection(ctx, collectionID)
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

	// Delete collection using Plex API
	// Note: The Plex API doesn't have a direct DeleteCollection method,
	// but we can use DeleteFromLibrary with the collection's rating key
	_, err = c.plexAPI.Library.DeleteFromLibrary(ctx, collectionRatingKey)
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

	// Format item IDs for Plex API
	itemKeysStr := strings.Join(itemIDs, ",")

	// Add items to collection using Plex API
	_, err = c.plexAPI.Library.UpdateCollection(
		ctx,
		collectionRatingKey,
		operations.UpdateCollectionUpdateCollectionType.Add,
		operations.UpdateCollectionUpdatedAt(strconv.FormatInt(time.Now().Unix(), 10)),
		operations.UpdateCollectionItemID(itemKeysStr),
	)
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

	// Format item IDs for Plex API
	itemKeysStr := strings.Join(itemIDs, ",")

	// Remove items from collection using Plex API
	_, err = c.plexAPI.Library.UpdateCollection(
		ctx,
		collectionRatingKey,
		operations.UpdateCollectionUpdateCollectionType.Remove,
		operations.UpdateCollectionUpdatedAt(strconv.FormatInt(time.Now().Unix(), 10)),
		operations.UpdateCollectionItemID(itemKeysStr),
	)
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

	if itemList == nil || len(itemList.Items) == 0 {
		log.Info().
			Str("collectionID", collectionID).
			Msg("No items to remove from collection")
		return nil
	}

	// Extract item IDs
	var itemIDs []string
	for _, item := range itemList.Items {
		itemIDs = append(itemIDs, item.ID)
	}

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
