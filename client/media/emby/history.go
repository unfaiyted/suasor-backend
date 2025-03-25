package emby

import (
	"context"
	"fmt"

	"github.com/antihax/optional"
	"suasor/client/media/types"
	embyclient "suasor/internal/clients/embyAPI"
	"suasor/utils"
)

// GetWatchHistory retrieves watch history from the Emby server
func (e *EmbyClient) GetWatchHistory(ctx context.Context, options *types.QueryOptions) ([]types.WatchHistoryItem[types.MediaData], error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.ClientID).
		Str("clientType", string(e.ClientType)).
		Msg("Retrieving watch history from Emby server")

	if e.config.UserID == "" {
		return nil, fmt.Errorf("user ID is required for watch history")
	}

	queryParams := embyclient.ItemsServiceApiGetUsersByUseridItemsOpts{
		IsPlayed:  optional.NewBool(true),
		Recursive: optional.NewBool(true),
	}

	// Apply options for pagination
	if options != nil {
		if options.Limit > 0 {
			queryParams.Limit = optional.NewInt32(int32(options.Limit))
		}
		if options.Offset > 0 {
			queryParams.StartIndex = optional.NewInt32(int32(options.Offset))
		}
	}

	items, resp, err := e.client.ItemsServiceApi.GetUsersByUseridItems(ctx, e.config.UserID, &queryParams)
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", e.config.BaseURL).
			Str("apiEndpoint", "/Users/"+e.config.UserID+"/Items").
			Msg("Failed to fetch watch history from Emby")
		return nil, fmt.Errorf("failed to fetch watch history: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(items.Items)).
		Int("totalRecordCount", int(items.TotalRecordCount)).
		Msg("Successfully retrieved watch history from Emby")

	history := make([]types.WatchHistoryItem[types.MediaData], 0)
	for _, item := range items.Items {
		if item.UserData != nil && item.UserData.Played {
			watchItem, err := e.convertToWatchHistoryItem(&item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", item.Id).
					Str("itemName", item.Name).
					Msg("Error converting Emby item to watch history format")
				continue
			}
			history = append(history, watchItem)
		}
	}

	log.Info().
		Int("historyItemsReturned", len(history)).
		Msg("Completed GetWatchHistory request")

	return history, nil
}
