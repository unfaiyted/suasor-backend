package jellyfin

import (
	"context"
	"fmt"

	t "suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/logger"
)

func (j *JellyfinClient) GetPlayHistory(ctx context.Context, options *t.QueryOptions) (*models.MediaItemDataList, error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.GetClientID()).
		Str("clientType", string(j.GetClientType())).
		Str("baseURL", j.config.GetBaseURL()).
		Msg("Retrieving watch history from Jellyfin server")

	// Call the Jellyfin API to get resumed items
	log.Debug().Msg("Making API request to Jellyfin server for resume items")
	userItemData := j.client.ItemsAPI.GetItems(ctx).
		UserId(j.config.UserID).
		IsPlayed(true)

	NewJellyfinQueryOptions(options).
		SetItemsRequest(&userItemData)

	results, resp, err := userItemData.Execute()

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.GetBaseURL()).
			Str("apiEndpoint", "/UserItems/Resume").
			Int("statusCode", 0).
			Msg("Failed to fetch watch history from Jellyfin")
		return nil, fmt.Errorf("failed to fetch watch history: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(results.Items)).
		Int("totalRecordCount", int(*results.TotalRecordCount)).
		Msg("Successfully retrieved watch history from Jellyfin")

	userHistoryDatas, err := GetMixedMediaItemsData(j, ctx, results.Items)

	log.Info().
		Int("statusCode", resp.StatusCode).
		Msg("Successfully retrieved user item data from Jellyfin")

	log.Info().
		Int("historyItemsReturned", userHistoryDatas.GetTotalItems()).
		Msg("Completed GetWatchHistory request")

	return userHistoryDatas, nil
}
