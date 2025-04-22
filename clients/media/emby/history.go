package emby

import (
	"context"
	"fmt"

	"github.com/antihax/optional"
	"suasor/clients/media/types"
	embyclient "suasor/internal/clients/embyAPI"
	"suasor/types/models"
	"suasor/utils/logger"
)

// GetWatchHistory retrieves watch history from the Emby server
func (e *EmbyClient) GetPlayHistory(ctx context.Context, options *types.QueryOptions) (*models.MediaItemDataList, error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.ClientID).
		Str("clientType", string(e.ClientType)).
		Msg("Retrieving watch history from Emby server")

	if e.embyConfig().UserID == "" {
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

	items, resp, err := e.client.ItemsServiceApi.GetUsersByUseridItems(ctx, e.embyConfig().UserID, &queryParams)
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", e.embyConfig().BaseURL).
			Str("apiEndpoint", "/Users/"+e.embyConfig().UserID+"/Items").
			Msg("Failed to fetch watch history from Emby")
		return nil, fmt.Errorf("failed to fetch watch history: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(items.Items)).
		Int("totalRecordCount", int(items.TotalRecordCount)).
		Msg("Successfully retrieved watch history from Emby")

	history, err := GetMixedMediaItemsData(e, ctx, items.Items)
	if err != nil {
		return nil, err
	}

	return history, nil
}
