package emby

import (
	"context"
	"fmt"
	"strings"
	embyclient "suasor/internal/clients/embyAPI"
	"suasor/types/models"

	"github.com/antihax/optional"
	"suasor/client/media/types"
	"suasor/utils"
)

func (e *EmbyClient) CreateCollection(ctx context.Context, name string, description string, itemIDs []string) (*models.MediaItem[*types.Collection], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.ClientID).
		Str("clientType", string(e.ClientType)).
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
	collection := e.ConvertItemToCollection(collectionResponse.Items[0])
	collection.SetClientInfo(e.ClientID, e.ClientType, newCollection.Id)

	log.Info().
		Str("collectionId", collectionID).
		Str("name", name).
		Msg("Successfully created collection in Emby")

	return collection, nil
}

func (e *EmbyClient) ConvertItemToCollection(item embyclient.BaseItemDto) *models.MediaItem[*types.Collection] {

	collection := models.MediaItem[*types.Collection]{
		Data: &types.Collection{
			ItemList: types.ItemList{
				Details: types.MediaDetails{
					Title:       item.Name,
					Description: item.Overview,
					Artwork:     e.getArtworkURLs(&item),
				},
			},
		},
		Type: "collection",
	}

	if collection.Data.ItemList.Details.Title == "" {
		collection.Data.ItemList.Details.Title = "Unnamed Collection"
	}

	return &collection
}

func (e *EmbyClient) GetCollections(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Collection], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.ClientID).
		Str("clientType", string(e.ClientType)).
		Msg("Retrieving collections from Emby")

	// Prepare the query options
	includeItemTypes := "BoxSet"

	queryParams := embyclient.ItemsServiceApiGetItemsOpts{
		IncludeItemTypes: optional.NewString(includeItemTypes),
		Recursive:        optional.NewBool(true),
	}
	applyQueryOptions(&queryParams, options)
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
		collection := e.ConvertItemToCollection(item)
		if collection != nil {
			collection.SetClientInfo(e.ClientID, e.ClientType, item.Id)
			collections = append(collections, collection)
		}
	}

	log.Info().
		Int("collectionCount", len(collections)).
		Msg("Successfully retrieved collections from Emby")

	return collections, nil
}

func (e *EmbyClient) GetCollectionByID(ctx context.Context, collectionID string) (*models.MediaItem[*types.Collection], error) {

	// opts := embclient.ItemsServiceApiGetItemsOpts{
	// 	Ids:    optional.NewString(collectionID),
	// 	Fields: optional.NewString("PrimaryImageAspectRatio,BasicSyncInfo,CanDelete,Container,DateCreated,PremiereDate,Genres,MediaSources,Overview,ParentId,Path,SortName,Studios,Taglines"),
	// }
	opts := types.QueryOptions{
		Limit:   1,
		ItemIDs: collectionID,
	}
	collections, err := e.GetCollections(ctx, &opts)
	if err != nil {
		return nil, err
	}
	if len(collections) == 0 {
		return nil, fmt.Errorf("collection not found")
	}
	collection := collections[0]

	return collection, err
}

func (e *EmbyClient) UpdateCollection(ctx context.Context, collectionID string, name string, description string) (*models.MediaItem[*types.Collection], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.ClientID).
		Str("clientType", string(e.ClientType)).
		Str("collectionID", collectionID).
		Str("name", name).
		Msg("Updating collection in Emby")

	// originalCollection, err := e.GetCollectionByID(ctx, collectionID)

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
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.ClientID).
		Str("clientType", string(e.ClientType)).
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

func (e *EmbyClient) AddItemsToCollection(ctx context.Context, collectionID string, itemIDs []string) error {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.ClientID).
		Str("clientType", string(e.ClientType)).
		Str("collectionID", collectionID).
		Int("itemCount", len(itemIDs)).
		Msg("Adding items to collection in Emby")

	// Emby requires a separate call for each item
	for _, itemID := range itemIDs {
		// Add item to collection
		_, err := e.client.CollectionServiceApi.PostCollectionsByIdItems(ctx, collectionID, itemID)
		if err != nil {
			log.Error().
				Err(err).
				Str("collectionID", collectionID).
				Str("itemID", itemID).
				Msg("Failed to add item to collection in Emby")
			return err
		}
	}

	log.Info().
		Str("collectionID", collectionID).
		Int("itemCount", len(itemIDs)).
		Msg("Successfully added items to collection in Emby")

	return nil
}

func (e *EmbyClient) RemoveItemsFromCollection(ctx context.Context, collectionID string, itemIDs []string) error {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.ClientID).
		Str("clientType", string(e.ClientType)).
		Str("collectionID", collectionID).
		Int("itemCount", len(itemIDs)).
		Msg("Removing items from collection in Emby")

	// Emby requires a separate call for each item
	for _, itemID := range itemIDs {
		// Remove item from collection
		_, err := e.client.CollectionServiceApi.PostCollectionsByIdItemsDelete(ctx, collectionID, itemID)
		if err != nil {
			log.Error().
				Err(err).
				Str("collectionID", collectionID).
				Str("itemID", itemID).
				Msg("Failed to remove item from collection in Emby")
			return err
		}
	}

	log.Info().
		Str("collectionID", collectionID).
		Int("itemCount", len(itemIDs)).
		Msg("Successfully removed items from collection in Emby")

	return nil
}

func (e *EmbyClient) GetCollectionItems(ctx context.Context, collectionID string) (*models.MediaItems, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.ClientID).
		Str("clientType", string(e.ClientType)).
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

	response, _, err := e.client.ItemsServiceApi.GetItems(ctx, &opts)
	if err != nil {
		log.Error().
			Err(err).
			Str("collectionID", collectionID).
			Msg("Failed to fetch collection items from Emby")
		return nil, err
	}

	if response.Items == nil || len(response.Items) == 0 {
		log.Info().
			Str("collectionID", collectionID).
			Msg("No items found in collection")
		return nil, nil
	}

	// Convert items to model format
	items := models.MediaItems{}

	for _, item := range response.Items {
		e.convertItemToMediaItems(ctx, &item, &items)
	}

	log.Info().
		Str("collectionID", collectionID).
		Int("itemCount", items.TotalItems).
		Msg("Successfully retrieved collection items from Emby")

	return &items, nil
}
