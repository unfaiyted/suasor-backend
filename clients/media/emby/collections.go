package emby

import (
	"context"
	"fmt"
	"strings"
	embyclient "suasor/internal/clients/embyAPI"
	"suasor/types/models"

	"github.com/antihax/optional"
	"suasor/clients/media/types"
	"suasor/utils/logger"
)

func (e *EmbyClient) SupportsCollections() bool { return true }

func (e *EmbyClient) SearchCollections(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Collection], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.GetClientID()).
		Str("clientType", string(e.GetClientType())).
		Msg("Retrieving collections from Emby")

	// Prepare the query options
	includeItemTypes := "BoxSet"

	queryParams := embyclient.ItemsServiceApiGetItemsOpts{
		IncludeItemTypes: optional.NewString(includeItemTypes),
		Recursive:        optional.NewBool(true),
	}

	ApplyClientQueryOptions(ctx, &queryParams, options)

	// Get collections
	results, _, err := e.client.ItemsServiceApi.GetItems(ctx, &queryParams)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to fetch collections from Emby")
		return nil, err
	}

	if results.Items == nil || len(results.Items) == 0 {
		log.Info().Msg("No collections returned from Emby")
		return []*models.MediaItem[*types.Collection]{}, nil
	}

	collections := make([]*models.MediaItem[*types.Collection], 0, len(results.Items))

	for _, item := range results.Items {
		itemCollection, err := GetItem[*types.Collection](ctx, e, &item)
		collection, err := GetMediaItem[*types.Collection](ctx, e, itemCollection, item.Id)
		if err != nil {
			log.Warn().
				Err(err).
				Str("itemID", item.Id).
				Str("itemName", item.Name).
				Msg("Error converting Emby item to collection format")
			continue
		}
		if collection != nil {
			collection.SetClientInfo(e.GetClientID(), e.GetClientType(), item.Id)
			collections = append(collections, collection)
		}
	}

	log.Info().
		Int("collectionCount", len(collections)).
		Msg("Successfully retrieved collections from Emby")

	return collections, nil
}

func (e *EmbyClient) SearchCollectionItems(ctx context.Context, collectionID string, options *types.QueryOptions) (*models.MediaItemResults, error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.GetClientID()).
		Str("clientType", string(e.GetClientType())).
		Str("collectionID", collectionID).
		Msg("Searching collection items from Emby")

	// Prepare the query options
	includeItemTypes := "BoxSet"

	queryParams := embyclient.ItemsServiceApiGetItemsOpts{
		IncludeItemTypes: optional.NewString(includeItemTypes),
		Recursive:        optional.NewBool(true),
	}

	ApplyClientQueryOptions(ctx, &queryParams, options)

	// Get collection
	results, _, err := e.client.ItemsServiceApi.GetItems(ctx, &queryParams)
	if err != nil {
		log.Error().
			Err(err).
			Str("collectionID", collectionID).
			Msg("Failed to fetch collection items from Emby")
		return nil, err
	}

	if results.Items == nil || len(results.Items) == 0 {
		log.Info().
			Str("collectionID", collectionID).
			Msg("No items found in collection")
		return nil, nil
	}

	items, err := GetMixedMediaItems(e, ctx, results.Items)
	if err != nil {
		return nil, err
	}

	log.Info().
		Str("collectionID", collectionID).
		Int("itemCount", items.Len()).
		Msg("Successfully retrieved collection items from Emby")

	return items, nil
}

func (e *EmbyClient) GetCollection(ctx context.Context, collectionID string) (*models.MediaItem[*types.Collection], error) {
	return e.GetCollectionByID(ctx, collectionID)
}

func (e *EmbyClient) GetCollectionByID(ctx context.Context, collectionID string) (*models.MediaItem[*types.Collection], error) {

	opts := types.QueryOptions{
		Limit:   1,
		ItemIDs: collectionID,
	}
	collections, err := e.SearchCollections(ctx, &opts)
	if err != nil {
		return nil, err
	}
	if len(collections) == 0 {
		return nil, fmt.Errorf("collection not found")
	}
	collection := collections[0]

	return collection, err
}

func (e *EmbyClient) GetCollectionItems(ctx context.Context, collectionID string) (*models.MediaItemList[*types.Collection], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.GetClientID()).
		Str("clientType", string(e.GetClientType())).
		Str("collectionID", collectionID).
		Msg("Retrieving collection items from Emby")

	opts := embyclient.ItemsServiceApiGetItemsOpts{
		ParentId:     optional.NewString(collectionID),
		SortBy:       optional.NewString("SortName"),
		SortOrder:    optional.NewString("Ascending"),
		Recursive:    optional.NewBool(true),
		Fields:       optional.NewString("Overview,Path,Genres,Tags"),
		EnableImages: optional.NewBool(true),
	}

	itemsResponse, _, err := e.client.ItemsServiceApi.GetItems(ctx, &opts)
	if err != nil {
		log.Error().
			Err(err).
			Str("collectionID", collectionID).
			Msg("Failed to fetch collection items from Emby")
		return nil, err
	}

	if itemsResponse.Items == nil || len(itemsResponse.Items) == 0 {
		log.Info().
			Str("collectionID", collectionID).
			Msg("No items found in collection")
		return nil, nil
	}

	items, err := GetMixedMediaItems(e, ctx, itemsResponse.Items)
	if err != nil {
		return nil, err
	}

	log.Info().
		Str("collectionID", collectionID).
		Int("itemCount", items.Len()).
		Msg("Successfully retrieved collection items from Emby")

	list, err := e.GetCollectionByID(ctx, collectionID)
	if err != nil {
		return nil, err
	}
	mediaItemList := models.NewMediaItemList[*types.Collection](list, e.GetClientID(), 0)

	return mediaItemList, nil
}

func (e *EmbyClient) CreateCollection(ctx context.Context, name string, description string, itemIDs []string) (*models.MediaItem[*types.Collection], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.GetClientID()).
		Str("clientType", string(e.GetClientType())).
		Str("name", name).
		Msg("Creating collection in Emby")

	ops := embyclient.CollectionServiceApiPostCollectionsOpts{
		Name: optional.NewString(name),
		Ids:  optional.NewString(strings.Join(itemIDs, ",")),
	}

	// // Create collection
	newCollection, _, err := e.client.CollectionServiceApi.PostCollections(ctx, &ops)
	if err != nil {
		log.Error().
			Err(err).
			Str("name", name).
			Msg("Failed to create collection in Emby")
		return nil, err
	}

	// Get the collection ID
	collectionID := fmt.Sprintf("%v", newCollection.Id)
	if collectionID == "" {
		log.Error().
			Msg("Empty collection ID returned from Emby")
		return nil, fmt.Errorf("empty collection ID")
	}

	// Add items to collection if needed
	if len(itemIDs) > 0 {
		// TODO: Add items to collection
	}

	collectionOpts := &embyclient.ItemsServiceApiGetItemsOpts{
		Ids:    optional.NewString(newCollection.Id),
		Fields: optional.NewString("PrimaryImageAspectRatio,BasicSyncInfo,CanDelete,Container,DateCreated,PremiereDate,Genres,MediaSources,Overview,ParentId,Path,SortName,Studios,Taglines"),
	}
	// Get the collection details
	collectionResponse, _, err := e.client.ItemsServiceApi.GetItems(ctx, collectionOpts)
	if err != nil {
		log.Error().
			Err(err).
			Str("collectionId", collectionID).
			Msg("Failed to get collection details from Emby")
		return nil, err
	}

	// Convert to Collection model
	collectionItem, err := GetItem[*types.Collection](ctx, e, &collectionResponse.Items[0])
	collection, err := GetMediaItem[*types.Collection](ctx, e, collectionItem, collectionResponse.Items[0].Id)
	if err != nil {
		log.Warn().
			Err(err).
			Str("collectionId", collectionID).
			Str("name", name).
			Msg("Error converting Emby item to collection format")
		return nil, err
	}
	// collection.SetClientInfo(e.GetClientID(), e.GetClientType(), newCollection.Id)

	log.Info().
		Str("collectionId", collectionID).
		Str("name", name).
		Msg("Successfully created collection in Emby")

	return collection, nil
}

func (e *EmbyClient) UpdateCollection(ctx context.Context, collectionID string, name string, description string) (*models.MediaItem[*types.Collection], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.GetClientID()).
		Str("clientType", string(e.GetClientType())).
		Str("collectionID", collectionID).
		Str("name", name).
		Msg("Updating collection in Emby")

	updatedCollection := embyclient.BaseItemDto{
		Name:      name,
		Id:        collectionID,
		Overview:  description,
		MediaType: "BoxSet",
	}

	response, err := e.client.ItemUpdateServiceApi.PostItemsByItemid(ctx, updatedCollection, collectionID)
	if err != nil {
		log.Error().
			Err(err).
			Str("collectionID", collectionID).
			Msg("Failed to get collection from Emby")
		return nil, err
	}

	if response.StatusCode != 200 {
		log.Error().
			Err(err).
			Str("collectionID", collectionID).
			Msg("Failed to update the collection from Emby")
		return nil, fmt.Errorf("failed to get collection from Emby")
	}

	// Get the updated collection
	finalCollection, err := e.GetCollectionByID(ctx, collectionID)
	if err != nil {
		log.Error().
			Err(err).
			Str("collectionID", collectionID).
			Msg("Failed to get updated collection from Emby")
		return nil, err
	}

	log.Info().
		Str("collectionID", collectionID).
		Str("name", finalCollection.Data.Details.Title).
		Msg("Successfully updated collection in Emby")

	return finalCollection, nil
}

func (e *EmbyClient) DeleteCollection(ctx context.Context, collectionID string) error {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.GetClientID()).
		Str("clientType", string(e.GetClientType())).
		Str("collectionID", collectionID).
		Msg("Deleting collection from Emby")

	// TODO: Dont see direct delete method for a whole collection. I can empty it. I know I did this before in python somehow. I know empty collections get hidden by the UI as well, but it would probably be best to just have a delete. this might be about safty or maybe im just missing the endpoint hiding in the docs.
	// Delete the collection
	// _, err := e.client.CollectionServiceApi.DeleteCollectionsByIdItems()(ctx, collectionID)
	// if err != nil {
	// 	log.Error().
	// 		Err(err).
	// 		Str("collectionID", collectionID).
	// 		Msg("Failed to delete collection from Emby")
	// 	return err
	// }
	//
	// log.Info().
	// 	Str("collectionID", collectionID).
	// 	Msg("Successfully deleted collection from Emby")

	return nil
}

func (e *EmbyClient) AddCollectionItem(ctx context.Context, collectionID string, itemID string) error {
	return e.AddCollectionItems(ctx, collectionID, []string{itemID})
}

func (e *EmbyClient) AddCollectionItems(ctx context.Context, collectionID string, itemIDs []string) error {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.GetClientID()).
		Str("clientType", string(e.GetClientType())).
		Str("collectionID", collectionID).
		Int("itemCount", len(itemIDs)).
		Msg("Adding items to collection in Emby")

	// Emby requires a separate call for each item
	strItemIDs := strings.Join(itemIDs, ",")

	// Add item to collection
	_, err := e.client.CollectionServiceApi.PostCollectionsByIdItems(ctx, strItemIDs, collectionID)
	if err != nil {
		log.Error().
			Err(err).
			Str("collectionID", collectionID).
			Strs("itemIDs", itemIDs).
			Msg("Failed to add item to collection in Emby")
		return err
	}

	log.Info().
		Str("collectionID", collectionID).
		Int("itemCount", len(itemIDs)).
		Msg("Successfully added items to collection in Emby")

	return nil
}

func (e *EmbyClient) RemoveCollectionItem(ctx context.Context, collectionID string, itemID string) error {
	return e.RemoveCollectionItems(ctx, collectionID, []string{itemID})
}

func (e *EmbyClient) RemoveCollectionItems(ctx context.Context, collectionID string, itemIDs []string) error {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.GetClientID()).
		Str("clientType", string(e.GetClientType())).
		Str("collectionID", collectionID).
		Int("itemCount", len(itemIDs)).
		Msg("Removing items from collection in Emby")

	strItemIDs := strings.Join(itemIDs, ",")

	_, err := e.client.CollectionServiceApi.PostCollectionsByIdItemsDelete(ctx, strItemIDs, collectionID)
	if err != nil {
		log.Error().
			Err(err).
			Str("collectionID", collectionID).
			Strs("itemIDs", itemIDs).
			Msg("Failed to remove item from collection in Emby")
		return err
	}

	log.Info().
		Str("collectionID", collectionID).
		Int("itemCount", len(itemIDs)).
		Msg("Successfully removed items from collection in Emby")

	return nil
}

func (e *EmbyClient) RemoveAllCollectionItems(ctx context.Context, collectionID string) error {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.GetClientID()).
		Str("clientType", string(e.GetClientType())).
		Str("collectionID", collectionID).
		Msg("Removing all items from collection in Emby")
	// Get all the items in the collection
	items, _, err := e.client.ItemsServiceApi.GetItems(ctx, &embyclient.ItemsServiceApiGetItemsOpts{
		ParentId:  optional.NewString(collectionID),
		Recursive: optional.NewBool(true),
		SortBy:    optional.NewString("SortName"),
		SortOrder: optional.NewString("Ascending"),
	})

	if err != nil {
		log.Error().
			Err(err).
			Str("collectionID", collectionID).
			Msg("Failed to fetch collection items from Emby")
		return err
	}

	if len(items.Items) == 0 {
		log.Info().
			Str("collectionID", collectionID).
			Msg("No items found in collection")
		return nil
	}

	itemIDs := make([]string, 0, len(items.Items))
	for _, item := range items.Items {
		itemIDs = append(itemIDs, item.Id)
	}
	strItemIDs := strings.Join(itemIDs, ",")

	// Emby requires a separate call for each item
	_, err = e.client.CollectionServiceApi.PostCollectionsByIdItemsDelete(ctx, strItemIDs, collectionID)
	if err != nil {
		log.Error().
			Err(err).
			Str("collectionID", collectionID).
			Msg("Failed to remove all items from collection in Emby")
		return err
	}

	log.Info().
		Str("collectionID", collectionID).
		Msg("Successfully removed all items from collection in Emby")

	return nil
}
