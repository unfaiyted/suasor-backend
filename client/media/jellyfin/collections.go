package jellyfin

import (
	"context"
	"fmt"

	jellyfin "github.com/sj14/jellyfin-go/api"
	t "suasor/client/media/types"
	"suasor/types/models"
	"suasor/utils"
)

func (j *JellyfinClient) GetCollections(ctx context.Context, options *t.QueryOptions) ([]models.MediaItem[t.Collection], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.ClientID).
		Str("clientType", string(j.ClientType)).
		Str("baseURL", j.config.BaseURL).
		Msg("Retrieving collections from Jellyfin server")

	// Set up query parameters
	includeItemTypes := []jellyfin.BaseItemKind{jellyfin.BASEITEMKIND_BOX_SET}

	limit, startIndex, sortBy, sortOrder := j.getQueryParameters(options)

	// Call the Jellyfin API
	log.Debug().Msg("Making API request to Jellyfin server for collections")
	itemsReq := j.client.ItemsAPI.GetItems(ctx).
		IncludeItemTypes(includeItemTypes).
		Recursive(true).
		Limit(*limit).
		StartIndex(*startIndex).
		SortBy(sortBy).
		SortOrder(sortOrder)

	result, resp, err := itemsReq.Execute()

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.BaseURL).
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
	collections := make([]models.MediaItem[t.Collection], 0)
	for _, item := range result.Items {
		if *item.Type == "BoxSet" {
			collection, err := j.convertToCollection(ctx, &item)
			if err != nil {
				// Log error but continue
				log.Warn().
					Err(err).
					Str("collectionID", *item.Id).
					Msg("Error converting Jellyfin item to collection format")
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
