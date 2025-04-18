package emby

import (
	"context"
	"fmt"

	"github.com/antihax/optional"
	"suasor/client/media/types"
	embyclient "suasor/internal/clients/embyAPI"
	"suasor/types/models"
	"suasor/utils"
)

// GetWatchHistory retrieves watch history from the Emby server
func (e *EmbyClient) GetPlayHistory(ctx context.Context, options *types.QueryOptions) (*models.MediaItemDatas, error) {
	log := utils.LoggerFromContext(ctx)

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

	history := convertToMediaItemDatas(e, ctx, items.Items)

	return history, nil
}

func convertToMediaItemDatas(e *EmbyClient, ctx context.Context, items []embyclient.BaseItemDto) *models.MediaItemDatas {
	log := utils.LoggerFromContext(ctx)
	datas := &models.MediaItemDatas{}

	if items == nil {
		return datas
	}

	for _, item := range items {
		if item.Type_ == "Movie" {
			movie, err := convertToMediaItemData[*types.Movie](e, ctx, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", item.Id).
					Str("itemName", item.Name).
					Msg("Error converting Emby item to watch history format")
				continue
			}
			datas.AddMovie(movie)
		} else if item.Type_ == "Episode" {
			episode, err := convertToMediaItemData[*types.Episode](e, ctx, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", item.Id).
					Str("itemName", item.Name).
					Msg("Error converting Emby item to watch history format")
				continue
			}
			datas.AddEpisode(episode)
		} else if item.Type_ == "Audio" {
			track, err := convertToMediaItemData[*types.Track](e, ctx, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", item.Id).
					Str("itemName", item.Name).
					Msg("Error converting Emby item to watch history format")
				continue
			}
			datas.AddTrack(track)
		} else if item.Type_ == "Playlist" {
			playlist, err := convertToMediaItemData[*types.Playlist](e, ctx, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", item.Id).
					Str("itemName", item.Name).
					Msg("Error converting Emby item to watch history format")
				continue
			}
			datas.AddPlaylist(playlist)
		} else if item.Type_ == "Series" {
			series, err := convertToMediaItemData[*types.Series](e, ctx, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", item.Id).
					Str("itemName", item.Name).
					Msg("Error converting Emby item to watch history format")
				continue
			}
			datas.AddSeries(series)
		} else if item.Type_ == "Season" {
			season, err := convertToMediaItemData[*types.Season](e, ctx, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", item.Id).
					Str("itemName", item.Name).
					Msg("Error converting Emby item to watch history format")
				continue
			}
			datas.AddSeason(season)
		} else if item.Type_ == "Collection" {
			collection, err := convertToMediaItemData[*types.Collection](e, ctx, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", item.Id).
					Str("itemName", item.Name).
					Msg("Error converting Emby item to watch history format")
				continue
			}
			datas.AddCollection(collection)
		}

	}

	return datas
}

func convertToMediaItemData[T types.MediaData](e *EmbyClient, ctx context.Context, item *embyclient.BaseItemDto) (*models.UserMediaItemData[T], error) {

	var mediaItemData *models.UserMediaItemData[T]
	mediaItemData.IsFavorite = item.UserData.IsFavorite
	mediaItemData.PlayedPercentage = item.UserData.PlayedPercentage
	mediaItemData.Completed = item.UserData.Played
	mediaItemData.PlayCount = item.UserData.PlayCount
	// mediaItemData.UserRating = item.UserData.Rating
	// mediaItemData.Watchlist = item.UserData.Watchlist
	mediaItemData.PlayedAt = item.UserData.LastPlayedDate

	mediaItem, err := convertTo[T](e, ctx, item)
	if err != nil {
		return nil, err
	}

	mediaItemData.Associate(mediaItem)

	return mediaItemData, err
}
