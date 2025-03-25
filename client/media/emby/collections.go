// collections.go
package emby

import (
	"context"
	"fmt"

	"github.com/antihax/optional"
	"suasor/client/media/types"
	embyclient "suasor/internal/clients/embyAPI"
	"suasor/utils"
)

// GetCollections retrieves collections from the Emby server
func (e *EmbyClient) GetCollections(ctx context.Context, options *types.QueryOptions) ([]types.MediaItem[types.Collection], error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.ClientID).
		Str("clientType", string(e.ClientType)).
		Msg("Retrieving collections from Emby server")

	queryParams := embyclient.ItemsServiceApiGetItemsOpts{
		IncludeItemTypes: optional.NewString("BoxSet"),
		Recursive:        optional.NewBool(true),
	}

	applyQueryOptions(&queryParams, options)

	items, resp, err := e.client.ItemsServiceApi.GetItems(ctx, &queryParams)
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", e.config.BaseURL).
			Str("apiEndpoint", "/Items").
			Msg("Failed to fetch collections from Emby")
		return nil, fmt.Errorf("failed to fetch collections: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(items.Items)).
		Int("totalRecordCount", int(items.TotalRecordCount)).
		Msg("Successfully retrieved collections from Emby")

	collections := make([]types.MediaItem[types.Collection], 0)
	for _, item := range items.Items {
		if item.Type_ == "BoxSet" {
			collection, err := e.convertToCollection(&item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("collectionID", item.Id).
					Str("collectionName", item.Name).
					Msg("Error converting Emby item to collection format")
				continue
			}
			collections = append(collections, collection)
		}
	}

	log.Info().
		Int("collectionsReturned", len(collections)).
		Msg("Completed GetCollections request")

	return collections, nil
}

// Add this converter method to converter.go
func (e *EmbyClient) convertToCollection(item *embyclient.BaseItemDto) (types.MediaItem[types.Collection], error) {
	if item == nil {
		return types.MediaItem[types.Collection]{}, fmt.Errorf("cannot convert nil item to collection")
	}

	collection := types.MediaItem[types.Collection]{
		Data: types.Collection{
			Details: types.MediaMetadata{
				Title:       item.Name,
				Description: item.Overview,
				Artwork:     e.getArtworkURLs(item),
			}},
		Type: "collection",
	}
	collection.SetClientInfo(e.ClientID, e.ClientType, item.Id)

	return collection, nil
}
