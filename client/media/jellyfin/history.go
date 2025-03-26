package jellyfin

import (
	"context"
	"fmt"

	jellyfin "github.com/sj14/jellyfin-go/api"
	t "suasor/client/media/types"
	"suasor/utils"
)

func (j *JellyfinClient) GetPlayHistory(ctx context.Context, options *t.QueryOptions) ([]t.MediaPlayHistory[t.MediaData], error) {
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
	historyItems := make([]t.MediaPlayHistory[t.MediaData], 0)
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

		historyItem := t.MediaPlayHistory[t.MediaData]{
			PlayedPercentage: *userData.PlayedPercentage.Get(),
			LastWatchedAt:    *userData.LastPlayedDate.Get(), // Default to now if not available
		}
		historyItem.Item.SetClientInfo(j.ClientID, j.ClientType, *item.Id)

		// Set type based on item type
		switch *item.Type {
		case jellyfin.BASEITEMKIND_MOVIE:
			historyItem.Item.Type = t.MEDIATYPE_MOVIE
			mediaItemMovie, err := j.convertToMovie(ctx, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("movieID", *item.Id).
					Str("movieName", *item.Name.Get()).
					Msg("Error converting Jellyfin item to movie format")
				continue
			}
			historyItem.Item.SetData(&historyItem.Item, mediaItemMovie.Data)
		case jellyfin.BASEITEMKIND_SERIES:
			historyItem.Item.Type = t.MEDIATYPE_SHOW
			mediaItemTVShow, err := j.convertToTVShow(ctx, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("showID", *item.Id).
					Str("showName", *item.Name.Get()).
					Msg("Error converting Jellyfin item to TV show format")
				continue
			}
			historyItem.Item.SetData(&historyItem.Item, mediaItemTVShow.Data)
		case jellyfin.BASEITEMKIND_EPISODE:
			historyItem.Item.Type = t.MEDIATYPE_EPISODE
			mediaItemEpisode, err := j.convertToEpisode(ctx, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("episodeID", *item.Id).
					Str("episodeName", *item.Name.Get()).
					Msg("Error converting Jellyfin item to episode format")
				continue
			}
			historyItem.Item.SetData(&historyItem.Item, mediaItemEpisode.Data)

		}

		if *item.Type == jellyfin.BASEITEMKIND_EPISODE {
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
