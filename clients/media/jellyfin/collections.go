package jellyfin

import (
	"context"
	"fmt"

	jellyfin "github.com/sj14/jellyfin-go/api"
	media "suasor/clients/media"
	t "suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/logger"
)

func (j *JellyfinClient) GetCollections(ctx context.Context, options *t.QueryOptions) ([]*models.MediaItem[*t.Collection], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.GetClientID()).
		Str("clientType", string(j.GetClientType())).
		Str("baseURL", j.config.GetBaseURL()).
		Msg("Retrieving collections from Jellyfin server")

		// Set up query parameters

	// Call the Jellyfin API
	log.Debug().Msg("Making API request to Jellyfin server for collections")
	itemsReq := j.client.ItemsAPI.GetItems(ctx)

	NewJellyfinQueryOptions(ctx, options).
		SetItemsRequest(ctx, &itemsReq)

	itemsReq.IncludeItemTypes([]jellyfin.BaseItemKind{jellyfin.BASEITEMKIND_BOX_SET})

	result, resp, err := itemsReq.Execute()

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.GetBaseURL()).
			Str("apiEndpoint", "/Items").
			Int("statusCode", 0).
			Msg("Failed to fetch collections from Jellyfin")
		return nil, fmt.Errorf("failed to fetch collections: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(result.Items)).
		Int("totalRecordCount", int(*result.TotalRecordCount)).
		Msg("Successfully retrieved collections from Jellyfin")

	// Convert results to expected format
	collections, err := GetMediaItemList[*t.Collection](ctx, j, result.Items)
	if err != nil {
		return nil, err
	}

	log.Info().
		Int("collectionsReturned", len(collections)).
		Msg("Completed GetCollections request")

	return collections, nil
}

func (j *JellyfinClient) GetCollectionByID(ctx context.Context, collectionID string) (*models.MediaItem[*t.Collection], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.GetClientID()).
		Str("clientType", string(j.GetClientType())).
		Str("collectionID", collectionID).
		Msg("Retrieving collection by ID from Jellyfin")

	options := &t.QueryOptions{
		ItemIDs:   collectionID,
		MediaType: t.MediaTypeCollection,
		Limit:     1,
	}
	collections, err := j.GetCollections(ctx, options)

	if err != nil {
		log.Error().Err(err).
			Str("collectionID", collectionID).
			Msg("Failed to get collection by ID from Jellyfin")
		return nil, err
	}
	if len(collections) == 0 {
		log.Info().
			Str("collectionID", collectionID).
			Msg("No collections found in Jellyfin")
		return nil, nil
	}

	return collections[0], nil
}

// GetCollectionItems retrieves all items in a collection from Jellyfin
func (j *JellyfinClient) GetCollectionItems(ctx context.Context, collectionID string, options *t.QueryOptions) (*models.MediaItemList[*t.Collection], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.GetClientID()).
		Str("clientType", string(j.GetClientType())).
		Str("collectionID", collectionID).
		Msg("Retrieving collection items from Jellyfin")

	// Set up query parameters

	var fields []jellyfin.ItemFields

	fields = append(fields, jellyfin.ITEMFIELDS_OVERVIEW)
	fields = append(fields, jellyfin.ITEMFIELDS_PATH)
	fields = append(fields, jellyfin.ITEMFIELDS_GENRES)
	fields = append(fields, jellyfin.ITEMFIELDS_TAGS)

	// Call the Jellyfin API to get items in the collection
	itemsReq := j.client.ItemsAPI.GetItems(ctx).
		ParentId(collectionID).Fields(fields).EnableImages(true)

	NewJellyfinQueryOptions(ctx, options).
		SetItemsRequest(ctx, &itemsReq)

	result, resp, err := itemsReq.Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("collectionID", collectionID).
			Str("baseURL", j.config.GetBaseURL()).
			Str("apiEndpoint", "/Items").
			Int("statusCode", 0).
			Msg("Failed to fetch collection items from Jellyfin")
		return nil, fmt.Errorf("failed to fetch collection items: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(result.Items)).
		Int("totalRecordCount", int(*result.TotalRecordCount)).
		Msg("Successfully retrieved collection items from Jellyfin")

	if len(result.Items) == 0 {
		log.Info().
			Str("collectionID", collectionID).
			Msg("No items found in collection")
		return &models.MediaItemList[*t.Collection]{}, nil
	}

	mediaItemResults, err := GetMixedMediaItems(j, ctx, result.Items)
	if err != nil {
		return nil, err
	}

	log.Info().
		Str("collectionID", collectionID).
		Int("itemCount", mediaItemResults.Len()).
		Msg("Successfully retrieved collection items from Jellyfin")

	list, err := j.GetCollectionByID(ctx, collectionID)
	if err != nil {
		return nil, err
	}

	mediaItemList := models.NewMediaItemList[*t.Collection](list, j.GetClientID(), 0)
	mediaItemList.Items = mediaItemResults
	return mediaItemList, nil

}

// SupportsCollections indicates if this client supports collections
func (j *JellyfinClient) SupportsCollections() bool {
	return true
}

// CreateCollection creates a new collection in Jellyfin
func (j *JellyfinClient) CreateCollection(ctx context.Context, name string, description string, collectionType string) (*models.MediaItem[*t.Collection], error) {
	// TODO: Implement collection creation for Jellyfin
	// This would involve:
	// 1. Creating a proper request to the Jellyfin API
	// 2. Converting the response to our internal model
	// 3. Using the new ItemList structure
	return nil, media.ErrFeatureNotSupported
}

// UpdateCollection updates an existing collection in Jellyfin
func (j *JellyfinClient) UpdateCollection(ctx context.Context, collectionID string, name string, description string) (*models.MediaItem[*t.Collection], error) {
	// TODO: Implement collection update for Jellyfin
	return nil, media.ErrFeatureNotSupported
}

// DeleteCollection deletes a collection from Jellyfin
func (j *JellyfinClient) DeleteCollection(ctx context.Context, collectionID string) error {
	// TODO: Implement collection deletion for Jellyfin
	return media.ErrFeatureNotSupported
}

// AddItemToCollection adds an item to a collection in Jellyfin
func (j *JellyfinClient) AddItemToCollection(ctx context.Context, collectionID string, itemID string) error {
	// TODO: Implement adding items to collections for Jellyfin
	return media.ErrFeatureNotSupported
}

// RemoveItemFromCollection removes an item from a collection in Jellyfin
func (j *JellyfinClient) RemoveItemFromCollection(ctx context.Context, collectionID string, itemID string) error {
	// TODO: Implement removing items from collections for Jellyfin
	return media.ErrFeatureNotSupported
}
