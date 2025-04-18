// helpers.go
package emby

import (
	"context"
	"fmt"
	"github.com/antihax/optional"
	"strings"
	"suasor/utils"
	// "suasor/utils"

	media "suasor/client/media"
	"suasor/client/media/types"
	embyclient "suasor/internal/clients/embyAPI"
	"suasor/types/models"
	"time"
)

func GetItem[T types.MediaData](
	ctx context.Context,
	client *EmbyClient,
	item *embyclient.BaseItemDto,
) (T, error) {
	return media.ConvertTo[*EmbyClient, *embyclient.BaseItemDto, T](
		media.GlobalMediaRegistry, client, ctx, item)
}
func GetMediaItem[T types.MediaData](
	ctx context.Context,
	client *EmbyClient,
	item T,
	itemID string,
) (*models.MediaItem[T], error) {
	mediaItem := models.MediaItem[T]{
		Data: item,
		Type: item.GetMediaType(),
	}
	mediaItem.SetClientInfo(client.ClientID, client.ClientType, itemID)

	return &mediaItem, nil
}
func GetMediaItemList[T types.MediaData](
	ctx context.Context,
	client *EmbyClient,
	items []embyclient.BaseItemDto,
) (*[]*models.MediaItem[T], error) {
	var mediaItems []*models.MediaItem[T]
	for _, item := range items {
		itemT, err := GetItem[T](ctx, client, &item)
		if err != nil {
			return nil, err
		}
		mediaItem, err := GetMediaItem[T](ctx, client, itemT, item.Id)
		if err != nil {
			return nil, err
		}
		mediaItem.SetClientInfo(client.ClientID, client.ClientType, item.Id)
		mediaItems = append(mediaItems, mediaItem)
	}

	return &mediaItems, nil
}
func GetMediaItemData[T types.MediaData](e *EmbyClient, ctx context.Context, item *embyclient.BaseItemDto) (*models.UserMediaItemData[T], error) {

	baseItem, err := GetItem[T](ctx, e, item)
	mediaItem, err := GetMediaItem[T](ctx, e, baseItem, item.Id)

	if err != nil {
		return nil, err
	}
	mediaItemData := models.UserMediaItemData[T]{
		Type:             types.MediaType(item.Type_),
		PlayedAt:         item.UserData.LastPlayedDate,
		PlayedPercentage: item.UserData.PlayedPercentage,
		IsFavorite:       item.UserData.IsFavorite,
		PlayCount:        item.UserData.PlayCount,
		PositionSeconds:  int(item.UserData.PlaybackPositionTicks / 10000000),
	}
	mediaItemData.Item.SetClientInfo(e.ClientID, e.ClientType, item.Id)
	mediaItemData.Associate(mediaItem)

	return &mediaItemData, err
}
func GetMediaItemDataList[T types.MediaData](e *EmbyClient, ctx context.Context, items []embyclient.BaseItemDto) ([]*models.UserMediaItemData[T], error) {
	var mediaItems []*models.UserMediaItemData[T]
	for _, item := range items {
		mediaItemData, err := GetMediaItemData[T](e, ctx, &item)
		if err != nil {
			return nil, err
		}
		mediaItems = append(mediaItems, mediaItemData)
	}

	return mediaItems, nil
}

func GetMixedMediaItems(e *EmbyClient, ctx context.Context, items []embyclient.BaseItemDto) (*models.MediaItems, error) {
	mediaItems := models.MediaItems{}
	for _, item := range items {

		if item.Type_ == "Movie" {
			movie, err := GetItem[*types.Movie](ctx, e, &item)
			if err != nil {
				return nil, err
			}
			movieItem, err := GetMediaItem[*types.Movie](ctx, e, movie, item.Id)
			if err != nil {
				return nil, err
			}
			mediaItems.AddMovie(movieItem)
		} else if item.Type_ == "Episode" {
			episode, err := GetItem[*types.Episode](ctx, e, &item)
			if err != nil {
				return nil, err
			}
			episodeItem, err := GetMediaItem[*types.Episode](ctx, e, episode, item.Id)
			if err != nil {
				return nil, err
			}
			mediaItems.AddEpisode(episodeItem)
		} else if item.Type_ == "Audio" {
			track, err := GetItem[*types.Track](ctx, e, &item)
			if err != nil {
				return nil, err
			}
			trackItem, err := GetMediaItem[*types.Track](ctx, e, track, item.Id)
			if err != nil {
				return nil, err
			}
			mediaItems.AddTrack(trackItem)
		} else if item.Type_ == "Playlist" {
			playlist, err := GetItem[*types.Playlist](ctx, e, &item)
			if err != nil {
				return nil, err
			}
			playlistItem, err := GetMediaItem[*types.Playlist](ctx, e, playlist, item.Id)
			if err != nil {
				return nil, err
			}
			mediaItems.AddPlaylist(playlistItem)
		} else if item.Type_ == "Series" {
			series, err := GetItem[*types.Series](ctx, e, &item)
			if err != nil {
				return nil, err
			}
			seriesItem, err := GetMediaItem[*types.Series](ctx, e, series, item.Id)
			if err != nil {
				return nil, err
			}
			mediaItems.AddSeries(seriesItem)
		} else if item.Type_ == "Season" {
			season, err := GetItem[*types.Season](ctx, e, &item)
			if err != nil {
				return nil, err
			}
			mediaItem, err := GetMediaItem[*types.Season](ctx, e, season, item.Id)
			if err != nil {
				return nil, err
			}
			mediaItems.AddSeason(mediaItem)
		} else if item.Type_ == "Collection" {
			collection, err := GetItem[*types.Collection](ctx, e, &item)
			if err != nil {
				return nil, err
			}
			mediaItem, err := GetMediaItem[*types.Collection](ctx, e, collection, item.Id)
			if err != nil {
				return nil, err
			}
			mediaItems.AddCollection(mediaItem)
		}
	}

	return &mediaItems, nil

}
func GetMixedMediaItemsData(e *EmbyClient, ctx context.Context, items []embyclient.BaseItemDto) (*models.MediaItemDatas, error) {
	log := utils.LoggerFromContext(ctx)
	datas := models.MediaItemDatas{
		TotalItems: 0,
	}

	for _, item := range items {
		if item.Type_ == "Movie" {
			movie, err := GetMediaItemData[*types.Movie](e, ctx, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", item.Id).
					Str("itemName", item.Name).
					Msg("Error converting Emby item to media data format")
				continue
			}
			datas.AddMovie(movie)
		} else if item.Type_ == "Episode" {
			episode, err := GetMediaItemData[*types.Episode](e, ctx, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", item.Id).
					Str("itemName", item.Name).
					Msg("Error converting Emby item to media data format")
				continue
			}
			datas.AddEpisode(episode)
		} else if item.Type_ == "Audio" {
			track, err := GetMediaItemData[*types.Track](e, ctx, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", item.Id).
					Str("itemName", item.Name).
					Msg("Error converting Emby item to media data format")
				continue
			}
			datas.AddTrack(track)
		} else if item.Type_ == "Playlist" {
			playlist, err := GetMediaItemData[*types.Playlist](e, ctx, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", item.Id).
					Str("itemName", item.Name).
					Msg("Error converting Emby item to media data format")
				continue
			}
			datas.AddPlaylist(playlist)
		} else if item.Type_ == "Series" {
			series, err := GetMediaItemData[*types.Series](e, ctx, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", item.Id).
					Str("itemName", item.Name).
					Msg("Error converting Emby item to media data format")
				continue
			}
			datas.AddSeries(series)
		} else if item.Type_ == "Season" {
			season, err := GetMediaItemData[*types.Season](e, ctx, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", item.Id).
					Str("itemName", item.Name).
					Msg("Error converting Emby item to media data format")
				continue
			}
			datas.AddSeason(season)
		} else if item.Type_ == "Collection" {
			collection, err := GetMediaItemData[*types.Collection](e, ctx, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", item.Id).
					Str("itemName", item.Name).
					Msg("Error converting Emby item to media data format")
				continue
			}
			datas.AddCollection(collection)
		}

	}

	return &datas, nil
}

// Converts intenal mapped QueryOptions to external Emby API query options
func ApplyClientQueryOptions(queryParams *embyclient.ItemsServiceApiGetItemsOpts, options *types.QueryOptions) {
	if options == nil {
		return
	}

	if options.ItemIDs != "" {
		queryParams.Ids = optional.NewString(options.ItemIDs)
	}

	if options.Limit > 0 {
		queryParams.Limit = optional.NewInt32(int32(options.Limit))
	}

	if options.Offset > 0 {
		queryParams.StartIndex = optional.NewInt32(int32(options.Offset))
	}

	if options.Sort != "" {
		// TODO: Look into mapping the SortBy to emby definitions
		queryParams.SortBy = optional.NewString(string(options.Sort))
		if options.SortOrder == "desc" {
			queryParams.SortOrder = optional.NewString("Descending")
		} else {
			queryParams.SortOrder = optional.NewString("Ascending")
		}
	}

	// Apply search term (should be outside the filters check)
	if options.Query != "" {
		queryParams.SearchTerm = optional.NewString(options.Query)
		// Also enable recursive search when searching
		if !queryParams.Recursive.IsSet() {
			queryParams.Recursive = optional.NewBool(true)
		}
		// Increase the limit for search results if not explicitly set
		if options.Limit <= 0 && !queryParams.Limit.IsSet() {
			queryParams.Limit = optional.NewInt32(50) // Higher default for searches
		}
	}

	// Apply filters - use typed fields

	// Media type filter
	if options.MediaType != "" {
		queryParams.IncludeItemTypes = optional.NewString(string(options.MediaType))
	}

	// Genre filter
	if options.Genre != "" {
		queryParams.Genres = optional.NewString(options.Genre)
	}

	// Favorite filter
	if options.Favorites {
		queryParams.IsFavorite = optional.NewBool(true)
	}

	// Year filter
	if options.Year > 0 {
		queryParams.Years = optional.NewString(fmt.Sprintf("%d", options.Year))
	}

	// Person filters
	if options.Actor != "" {
		queryParams.Person = optional.NewString(options.Actor)
	}

	if options.Director != "" {
		queryParams.Person = optional.NewString(options.Director)
	}

	// Creator filter
	if options.Creator != "" {
		queryParams.Person = optional.NewString(options.Creator)
	}

	// Apply more advanced filters

	// Content rating filter
	if options.ContentRating != "" {
		queryParams.OfficialRatings = optional.NewString(options.ContentRating)
	}

	// Tags filter
	if len(options.Tags) > 0 {
		queryParams.Tags = optional.NewString(strings.Join(options.Tags, ","))
	}

	// Recently added filter
	if options.RecentlyAdded {
		queryParams.SortBy = optional.NewString("DateCreated,SortName")
		queryParams.SortOrder = optional.NewString("Descending")
	}

	// Recently played filter
	if options.RecentlyPlayed {
		queryParams.SortBy = optional.NewString("DatePlayed,SortName")
		queryParams.SortOrder = optional.NewString("Descending")
	}

	// Unwatched filter
	if options.Watched {
		queryParams.IsPlayed = optional.NewBool(true)
	}

	// Date filters
	if !options.DateAddedAfter.IsZero() {
		queryParams.MinDateLastSaved = optional.NewString(options.DateAddedAfter.Format(time.RFC3339))
	}

	if !options.DateAddedBefore.IsZero() {
		queryParams.MaxPremiereDate = optional.NewString(options.DateAddedBefore.Format(time.RFC3339))
	}

	if !options.ReleasedAfter.IsZero() {
		queryParams.MinPremiereDate = optional.NewString(options.ReleasedAfter.Format(time.RFC3339))
	}

	if !options.ReleasedBefore.IsZero() {
		queryParams.MaxPremiereDate = optional.NewString(options.ReleasedBefore.Format(time.RFC3339))
	}

	// Rating filter
	if options.MinimumRating > 0 {
		queryParams.MinCommunityRating = optional.NewFloat64(float64(options.MinimumRating))
	}

	// Debug logging removed to avoid logger dependency
}
