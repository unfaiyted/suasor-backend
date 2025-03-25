package jellyfin

import (
	"context"
	"fmt"

	jellyfin "github.com/sj14/jellyfin-go/api"
	t "suasor/client/media/types"
	"suasor/utils"
)

func (j *JellyfinClient) GetWatchHistory(ctx context.Context, options *t.QueryOptions) ([]t.WatchHistoryItem[t.MediaData], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.ClientID).
		Str("clientType", string(j.ClientType)).
		Str("baseURL", j.config.BaseURL).
		Msg("Retrieving watch history from Jellyfin server")

	limit, startIndex, sortBy, sortOrder := j.getQueryParameters(options)

	// Call the Jellyfin API to get resumed items
	log.Debug().Msg("Making API request to Jellyfin server for resume items")
	watchedItemsReq := j.client.ItemsAPI.GetItems(ctx).
		Limit(*limit).
		StartIndex(*startIndex).
		SortBy(sortBy).
		SortOrder(sortOrder).
		UserId(j.config.UserID).
		IsPlayed(true)

	result, resp, err := watchedItemsReq.Execute()

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.BaseURL).
			Str("apiEndpoint", "/UserItems/Resume").
			Int("statusCode", 0).
			Msg("Failed to fetch watch history from Jellyfin")
		return nil, fmt.Errorf("failed to fetch watch history: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(result.Items)).
		Int("totalRecordCount", int(*result.TotalRecordCount)).
		Msg("Successfully retrieved watch history from Jellyfin")

	// Convert results to expected format
	historyItems := make([]t.WatchHistoryItem[t.MediaData], 0)
	for _, item := range result.Items {

		userDataReq := j.client.ItemsAPI.GetItemUserData(ctx, *item.Id)
		userData, resp, err := userDataReq.Execute()

		if err != nil {
			continue
		}

		log.Info().
			Int("statusCode", resp.StatusCode).
			Int32("playCount", userData.GetPlayCount()).
			Msg("Successfully retrieved user item data from Jellyfin")

		historyItem := t.WatchHistoryItem[t.MediaData]{
			// Item: t.MediaData{
			// 	Details: t.MediaMetadata{
			// 		Title:       *item.Name.Get(),
			// 		Description: *item.Overview.Get(),
			// 		Artwork:     j.getArtworkURLs(&item),
			// 	},
			// },
			PlayedPercentage: *userData.PlayedPercentage.Get(),
			LastWatchedAt:    *userData.LastPlayedDate.Get(), // Default to now if not available
		}
		historyItem.Item.SetClientInfo(j.ClientID, j.ClientType, *item.Id)

		// Set type based on item type
		switch *item.Type {
		case jellyfin.BASEITEMKIND_MOVIE:
			historyItem.Item.Type = t.MEDIATYPE_MOVIE
		case jellyfin.BASEITEMKIND_SERIES:
			historyItem.Item.Type = t.MEDIATYPE_SHOW
		case jellyfin.BASEITEMKIND_EPISODE:
			historyItem.Item.Type = t.MEDIATYPE_EPISODE

			// Add additional episode info if available
			if item.SeriesName.IsSet() {
				historyItem.SeriesName = *item.SeriesName.Get()
			}
			if item.ParentIndexNumber.IsSet() {
				historyItem.SeasonNumber = int(*item.ParentIndexNumber.Get())
			}
			if item.IndexNumber.IsSet() {
				historyItem.EpisodeNumber = int(*item.IndexNumber.Get())
			}
		}

		// Set last played date if available
		if userData.LastPlayedDate.IsSet() {
			historyItem.LastWatchedAt = *userData.LastPlayedDate.Get()
		}

		historyItems = append(historyItems, historyItem)
	}

	log.Info().
		Int("historyItemsReturned", len(historyItems)).
		Msg("Completed GetWatchHistory request")

	return historyItems, nil
}
