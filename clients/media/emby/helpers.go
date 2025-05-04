// helpers.go
package emby

import (
	"context"
	"fmt"
	"github.com/antihax/optional"
	"strings"
	"suasor/utils/logger"

	media "suasor/clients/media"
	"suasor/clients/media/types"
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
		client, ctx, item)
}

func GetPlaylistItem(
	ctx context.Context,
	client *EmbyClient,
	item *embyclient.PlaylistsPlaylistCreationResult,
) (*types.Playlist, error) {
	return media.ConvertTo[*EmbyClient, *embyclient.PlaylistsPlaylistCreationResult, *types.Playlist](
		client, ctx, item)
}

func GetMediaItem[T types.MediaData](
	ctx context.Context,
	client *EmbyClient,
	item T,
	itemID string,
) (*models.MediaItem[T], error) {
	mediaItem := models.NewMediaItem[T](item.GetMediaType(), item)
	mediaItem.SetClientInfo(client.GetClientID(), client.GetClientType(), itemID)

	return mediaItem, nil
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
		mediaItem.SetClientInfo(client.GetClientID(), client.GetClientType(), item.Id)
		mediaItems = append(mediaItems, mediaItem)
	}

	return &mediaItems, nil
}
func GetMediaItemData[T types.MediaData](
	ctx context.Context,
	e *EmbyClient,
	item *embyclient.BaseItemDto,
) (*models.UserMediaItemData[T], error) {

	baseItem, err := GetItem[T](ctx, e, item)
	mediaItem, err := GetMediaItem[T](ctx, e, baseItem, item.Id)

	mediaItemData := models.NewUserMediaItemData[T](mediaItem, 0)
	mediaItemData.PlayedAt = item.UserData.LastPlayedDate
	mediaItemData.PlayedPercentage = item.UserData.PlayedPercentage
	mediaItemData.IsFavorite = item.UserData.IsFavorite
	mediaItemData.PlayCount = item.UserData.PlayCount
	mediaItemData.PositionSeconds = int(item.UserData.PlaybackPositionTicks / 10000000)

	mediaItemData.Associate(mediaItem)

	return mediaItemData, err
}

func GetMediaItemDataList[T types.MediaData](
	ctx context.Context,
	e *EmbyClient,
	items []embyclient.BaseItemDto,
) ([]*models.UserMediaItemData[T], error) {

	var mediaItems []*models.UserMediaItemData[T]
	for _, item := range items {
		mediaItemData, err := GetMediaItemData[T](ctx, e, &item)
		if err != nil {
			return nil, err
		}
		mediaItems = append(mediaItems, mediaItemData)
	}

	return mediaItems, nil
}

func GetMixedMediaItems(
	e *EmbyClient,
	ctx context.Context,
	items []embyclient.BaseItemDto,
) (*models.MediaItemList, error) {
	mediaItems := models.MediaItemList{}
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
func GetMixedMediaItemsData(
	e *EmbyClient,
	ctx context.Context,
	items []embyclient.BaseItemDto,
) (*models.MediaItemDataList, error) {
	log := logger.LoggerFromContext(ctx)
	datas := models.NewMediaItemDataList()

	for _, item := range items {

		log.Debug().
			Str("itemID", item.Id).
			Str("itemName", item.Name).
			Str("itemType", item.Type_).
			Msg("Processing item")

		// TODO: Handle other item types, SWITCH, improve this to some sort of generic loop
		if item.Type_ == "Movie" {
			movie, err := GetMediaItemData[*types.Movie](ctx, e, &item)
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
			episode, err := GetMediaItemData[*types.Episode](ctx, e, &item)
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
			track, err := GetMediaItemData[*types.Track](ctx, e, &item)
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
			playlist, err := GetMediaItemData[*types.Playlist](ctx, e, &item)
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
			series, err := GetMediaItemData[*types.Series](ctx, e, &item)
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
			season, err := GetMediaItemData[*types.Season](ctx, e, &item)
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
			collection, err := GetMediaItemData[*types.Collection](ctx, e, &item)
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

	log.Info().
		Int("totalMovies", len(datas.Movies)).
		Int("totalSeries", len(datas.Series)).
		Int("totalEpisodes", len(datas.Episodes)).
		Int("totalTracks", len(datas.Tracks)).
		Int("totalAlbums", len(datas.Albums)).
		Int("totalArtists", len(datas.Artists)).
		Int("totalPlaylists", len(datas.Playlists)).
		Int("totalCollections", len(datas.Collections)).
		Int("totalItems", datas.GetTotalItems()).
		Msg("Completed GetMixedMediaItemsData request")

	return datas, nil
}

// Converts intenal mapped QueryOptions to external Emby API query options
func ApplyClientQueryOptions(ctx context.Context, queryParams *embyclient.ItemsServiceApiGetItemsOpts, options *types.QueryOptions) {
	log := logger.LoggerFromContext(ctx)
	if options == nil {
		log.Debug().Msg("No options provided, skipping query options")
		return
	}

	if options.ItemIDs != "" {
		log.Debug().Str("itemIDs", options.ItemIDs).Msg("Applying item IDs filter")
		queryParams.Ids = optional.NewString(options.ItemIDs)
	}

	if options.Limit > 0 {
		log.Debug().Int("limit", options.Limit).Msg("Applying limit filter")
		queryParams.Limit = optional.NewInt32(int32(options.Limit))
	}

	if options.Offset > 0 {
		log.Debug().Int("offset", options.Offset).Msg("Applying offset filter")
		queryParams.StartIndex = optional.NewInt32(int32(options.Offset))
	}

	if options.Sort != "" {
		log.Debug().Str("sort", string(options.Sort)).Msg("Applying sort filter")
		// TODO: Look into mapping the SortBy to emby definitions
		queryParams.SortBy = optional.NewString(string(options.Sort))
		if options.SortOrder == "desc" {
			log.Debug().Msg("Applying descending sort order")
			queryParams.SortOrder = optional.NewString("Descending")
		} else {
			log.Debug().Msg("Applying ascending sort order")
			queryParams.SortOrder = optional.NewString("Ascending")
		}
	}

	// Apply search term (should be outside the filters check)
	if options.Query != "" {
		log.Debug().Str("query", options.Query).Msg("Applying search term filter")
		queryParams.SearchTerm = optional.NewString(options.Query)
		// Also enable recursive search when searching
		if !queryParams.Recursive.IsSet() {
			log.Debug().Msg("Enabling recursive search")
			queryParams.Recursive = optional.NewBool(true)
		}
		// Increase the limit for search results if not explicitly set
		if options.Limit <= 0 && !queryParams.Limit.IsSet() {
			log.Debug().Msg("Increasing default search limit")
			queryParams.Limit = optional.NewInt32(50) // Higher default for searches
		}
	}

	// Apply filters - use typed fields

	// Media type filter
	if options.MediaType != "" {
		log.Debug().Str("mediaType", string(options.MediaType)).Msg("Applying media type filter")
		queryParams.IncludeItemTypes = optional.NewString(string(options.MediaType))
	}

	// Genre filter
	if options.Genre != "" {
		log.Debug().Str("genre", options.Genre).Msg("Applying genre filter")
		queryParams.Genres = optional.NewString(options.Genre)
	}

	// Favorite filter
	if options.Favorites {
		log.Debug().Msg("Applying favorite filter")
		queryParams.IsFavorite = optional.NewBool(true)
	}

	// Year filter
	if options.Year > 0 {
		log.Debug().Int("year", options.Year).Msg("Applying year filter")
		queryParams.Years = optional.NewString(fmt.Sprintf("%d", options.Year))
	}

	// Person filters
	if options.Actor != "" {
		log.Debug().Str("actor", options.Actor).Msg("Applying actor filter")
		queryParams.Person = optional.NewString(options.Actor)
	}

	if options.Director != "" {
		log.Debug().Str("director", options.Director).Msg("Applying director filter")
		queryParams.Person = optional.NewString(options.Director)
	}

	// Creator filter
	if options.Creator != "" {
		log.Debug().Str("creator", options.Creator).Msg("Applying creator filter")
		queryParams.Person = optional.NewString(options.Creator)
	}

	// Apply more advanced filters

	// Content rating filter
	if options.ContentRating != "" {
		log.Debug().Str("contentRating", options.ContentRating).Msg("Applying content rating filter")
		queryParams.OfficialRatings = optional.NewString(options.ContentRating)
	}

	// Tags filter
	if len(options.Tags) > 0 {
		log.Debug().Strs("tags", options.Tags).Msg("Applying tags filter")
		queryParams.Tags = optional.NewString(strings.Join(options.Tags, ","))
	}

	// Recently added filter
	if options.RecentlyAdded {
		log.Debug().Msg("Applying recently added filter")
		queryParams.SortBy = optional.NewString("DateCreated,SortName")
		queryParams.SortOrder = optional.NewString("Descending")
	}

	// Recently played filter
	if options.RecentlyPlayed {
		log.Debug().Msg("Applying recently played filter")
		queryParams.SortBy = optional.NewString("DatePlayed,SortName")
		queryParams.SortOrder = optional.NewString("Descending")
	}

	// Unwatched filter
	if options.Watched {
		log.Debug().Msg("Applying unwatched filter")
		queryParams.IsPlayed = optional.NewBool(true)
	}

	// Date filters
	if options.DateAddedAfter != nil {
		log.Debug().Time("dateAddedAfter", *options.DateAddedAfter).Msg("Applying date added after filter")
		queryParams.MinDateLastSaved = optional.NewString(options.DateAddedAfter.Format(time.RFC3339))
	}

	if options.DateAddedBefore != nil {
		log.Debug().Time("dateAddedBefore", *options.DateAddedBefore).Msg("Applying date added before filter")
		queryParams.MaxPremiereDate = optional.NewString(options.DateAddedBefore.Format(time.RFC3339))
	}

	if options.ReleasedAfter != nil {
		log.Debug().Time("releasedAfter", *options.ReleasedAfter).Msg("Applying released after filter")
		queryParams.MinPremiereDate = optional.NewString(options.ReleasedAfter.Format(time.RFC3339))
	}

	if options.ReleasedBefore != nil {
		log.Debug().Time("releasedBefore", *options.ReleasedBefore).Msg("Applying released before filter")
		queryParams.MaxPremiereDate = optional.NewString(options.ReleasedBefore.Format(time.RFC3339))
	}

	// Rating filter
	if options.MinimumRating > 0 {
		log.Debug().Float32("minimumRating", options.MinimumRating).Msg("Applying minimum rating filter")
		queryParams.MinCommunityRating = optional.NewFloat64(float64(options.MinimumRating))
	}

}

func convertToExternalIDs(providerIds *map[string]string) types.ExternalIDs {
	externalIDs := types.ExternalIDs{}
	if providerIds == nil {
		return externalIDs
	}
	for key, value := range *providerIds {
		externalIDs = append(externalIDs, types.ExternalID{Source: strings.ToLower(key), ID: value})
	}
	return externalIDs
}
