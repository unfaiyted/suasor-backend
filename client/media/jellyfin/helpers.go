// helpers.go
package jellyfin

import (
	"context"
	"fmt"
	"strings"
	"time"

	jellyfin "github.com/sj14/jellyfin-go/api"
	media "suasor/client/media"
	"suasor/client/media/types"
	"suasor/types/models"
	"suasor/utils"
)

func GetItem[T types.MediaData](
	ctx context.Context,
	client *JellyfinClient,
	item *jellyfin.BaseItemDto,
) (T, error) {
	return media.ConvertTo[*JellyfinClient, *jellyfin.BaseItemDto, T](
		client, ctx, item)
}

func GetMediaItem[T types.MediaData](
	ctx context.Context,
	client *JellyfinClient,
	item T,
	itemID string,
) (*models.MediaItem[T], error) {
	mediaItem := models.NewMediaItem[T](item.GetMediaType(), item)
	mediaItem.SetClientInfo(client.ClientID, client.ClientType, itemID)

	return mediaItem, nil
}

func GetMediaItemList[T types.MediaData](
	ctx context.Context,
	client *JellyfinClient,
	items []jellyfin.BaseItemDto,
) ([]*models.MediaItem[T], error) {
	var mediaItems []*models.MediaItem[T]
	for _, item := range items {
		if item.Id == nil {
			continue
		}
		itemT, err := GetItem[T](ctx, client, &item)
		if err != nil {
			return nil, err
		}
		mediaItem, err := GetMediaItem[T](ctx, client, itemT, *item.Id)
		if err != nil {
			return nil, err
		}
		mediaItem.SetClientInfo(client.ClientID, client.ClientType, *item.Id)
		mediaItems = append(mediaItems, mediaItem)
	}

	return mediaItems, nil
}

func GetMediaItemData[T types.MediaData](
	ctx context.Context,
	j *JellyfinClient,
	item *jellyfin.BaseItemDto,
	itemData *jellyfin.UserItemDataDto,
) (*models.UserMediaItemData[T], error) {
	if itemData.ItemId == nil {
		return nil, fmt.Errorf("item data has no ID")
	}

	baseItem, err := GetItem[T](ctx, j, item)
	if err != nil {
		return nil, err
	}
	mediaItem, err := GetMediaItem[T](ctx, j, baseItem, *item.Id)
	if err != nil {
		return nil, err
	}

	mediaItemData := models.UserMediaItemData[T]{
		Type: types.MediaType(*item.Type),
	}

	// Set played data if available
	if item.UserData.IsSet() {
		userData := item.UserData.Get()

		if userData.LastPlayedDate.IsSet() {
			mediaItemData.PlayedAt = *userData.LastPlayedDate.Get()
		}

		mediaItemData.PositionSeconds = int(*userData.PlaybackPositionTicks / 10000000)
		if userData.PlayedPercentage.IsSet() {
			mediaItemData.PlayedPercentage = float64(*userData.PlayedPercentage.Get())
		}

		mediaItemData.PlayCount = int32(*userData.PlayCount)
		mediaItemData.IsFavorite = *userData.IsFavorite
	}

	mediaItemData.Item.SetClientInfo(j.ClientID, j.ClientType, *item.Id)
	mediaItemData.Associate(mediaItem)

	return &mediaItemData, nil
}

func GetMediaItemDataList[T types.MediaData](
	ctx context.Context,
	j *JellyfinClient,
	items []jellyfin.BaseItemDto,
) ([]*models.UserMediaItemData[T], error) {
	var mediaItems []*models.UserMediaItemData[T]
	for _, item := range items {

		userDataReq := j.client.ItemsAPI.GetItemUserData(ctx, *item.Id)
		userData, _, err := userDataReq.Execute()
		if err != nil {
			continue
		}
		if item.Id == nil {
			continue
		}
		mediaItemData, err := GetMediaItemData[T](ctx, j, &item, userData)
		if err != nil {
			continue
		}
		mediaItems = append(mediaItems, mediaItemData)
	}

	return mediaItems, nil
}

func GetMixedMediaItems(
	j *JellyfinClient,
	ctx context.Context,
	items []jellyfin.BaseItemDto,
) (*models.MediaItems, error) {
	mediaItems := models.MediaItems{}
	for _, item := range items {
		if item.Id == nil || item.Type == nil {
			continue
		}

		switch *item.Type {
		case jellyfin.BASEITEMKIND_MOVIE:
			movie, err := GetItem[*types.Movie](ctx, j, &item)
			if err != nil {
				return nil, err
			}
			movieItem, err := GetMediaItem[*types.Movie](ctx, j, movie, *item.Id)
			if err != nil {
				return nil, err
			}
			mediaItems.AddMovie(movieItem)
		case jellyfin.BASEITEMKIND_EPISODE:
			episode, err := GetItem[*types.Episode](ctx, j, &item)
			if err != nil {
				return nil, err
			}
			episodeItem, err := GetMediaItem[*types.Episode](ctx, j, episode, *item.Id)
			if err != nil {
				return nil, err
			}
			mediaItems.AddEpisode(episodeItem)
		case jellyfin.BASEITEMKIND_AUDIO:
			track, err := GetItem[*types.Track](ctx, j, &item)
			if err != nil {
				return nil, err
			}
			trackItem, err := GetMediaItem[*types.Track](ctx, j, track, *item.Id)
			if err != nil {
				return nil, err
			}
			mediaItems.AddTrack(trackItem)
		case jellyfin.BASEITEMKIND_PLAYLIST:
			playlist, err := GetItem[*types.Playlist](ctx, j, &item)
			if err != nil {
				return nil, err
			}
			playlistItem, err := GetMediaItem[*types.Playlist](ctx, j, playlist, *item.Id)
			if err != nil {
				return nil, err
			}
			mediaItems.AddPlaylist(playlistItem)
		case jellyfin.BASEITEMKIND_SERIES:
			series, err := GetItem[*types.Series](ctx, j, &item)
			if err != nil {
				return nil, err
			}
			seriesItem, err := GetMediaItem[*types.Series](ctx, j, series, *item.Id)
			if err != nil {
				return nil, err
			}
			mediaItems.AddSeries(seriesItem)
		case jellyfin.BASEITEMKIND_SEASON:
			season, err := GetItem[*types.Season](ctx, j, &item)
			if err != nil {
				return nil, err
			}
			mediaItem, err := GetMediaItem[*types.Season](ctx, j, season, *item.Id)
			if err != nil {
				return nil, err
			}
			mediaItems.AddSeason(mediaItem)
		case jellyfin.BASEITEMKIND_COLLECTION_FOLDER:
			collection, err := GetItem[*types.Collection](ctx, j, &item)
			if err != nil {
				return nil, err
			}
			mediaItem, err := GetMediaItem[*types.Collection](ctx, j, collection, *item.Id)
			if err != nil {
				return nil, err
			}
			mediaItems.AddCollection(mediaItem)
		}
	}

	return &mediaItems, nil
}

func GetMixedMediaItemsData(
	j *JellyfinClient,
	ctx context.Context,
	items []jellyfin.BaseItemDto,
) (*models.MediaItemDatas, error) {
	log := utils.LoggerFromContext(ctx)
	datas := models.MediaItemDatas{
		TotalItems: 0,
	}

	for _, item := range items {
		if item.Id == nil || item.Type == nil {
			continue
		}

		userDataReq := j.client.ItemsAPI.GetItemUserData(ctx, *item.Id)
		userData, _, err := userDataReq.Execute()
		if err != nil {
			continue
		}

		switch *item.Type {
		case jellyfin.BASEITEMKIND_MOVIE:
			movie, err := GetMediaItemData[*types.Movie](ctx, j, &item, userData)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", *item.Id).
					Str("itemName", getItemName(&item)).
					Msg("Error converting Jellyfin item to media data format")
				continue
			}
			datas.AddMovie(movie)
		case jellyfin.BASEITEMKIND_EPISODE:
			episode, err := GetMediaItemData[*types.Episode](ctx, j, &item, userData)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", *item.Id).
					Str("itemName", getItemName(&item)).
					Msg("Error converting Jellyfin item to media data format")
				continue
			}
			datas.AddEpisode(episode)
		case jellyfin.BASEITEMKIND_AUDIO:
			track, err := GetMediaItemData[*types.Track](ctx, j, &item, userData)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", *item.Id).
					Str("itemName", getItemName(&item)).
					Msg("Error converting Jellyfin item to media data format")
				continue
			}
			datas.AddTrack(track)
		case jellyfin.BASEITEMKIND_PLAYLIST:
			playlist, err := GetMediaItemData[*types.Playlist](ctx, j, &item, userData)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", *item.Id).
					Str("itemName", getItemName(&item)).
					Msg("Error converting Jellyfin item to media data format")
				continue
			}
			datas.AddPlaylist(playlist)
		case jellyfin.BASEITEMKIND_SERIES:
			series, err := GetMediaItemData[*types.Series](ctx, j, &item, userData)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", *item.Id).
					Str("itemName", getItemName(&item)).
					Msg("Error converting Jellyfin item to media data format")
				continue
			}
			datas.AddSeries(series)
		case jellyfin.BASEITEMKIND_SEASON:
			season, err := GetMediaItemData[*types.Season](ctx, j, &item, userData)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", *item.Id).
					Str("itemName", getItemName(&item)).
					Msg("Error converting Jellyfin item to media data format")
				continue
			}
			datas.AddSeason(season)
		case jellyfin.BASEITEMKIND_COLLECTION_FOLDER:
			collection, err := GetMediaItemData[*types.Collection](ctx, j, &item, userData)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", *item.Id).
					Str("itemName", getItemName(&item)).
					Msg("Error converting Jellyfin item to media data format")
				continue
			}
			datas.AddCollection(collection)
		}
	}

	return &datas, nil
}

// getItemName safely gets an item name with nil checking
func getItemName(item *jellyfin.BaseItemDto) string {
	if item == nil {
		return ""
	}

	if !item.Name.IsSet() {
		return ""
	}

	return *item.Name.Get()
}

func (j *JellyfinClient) getArtworkURLs(item *jellyfin.BaseItemDto) *types.Artwork {

	imageURLs := types.Artwork{}

	if item == nil || item.Id == nil {
		return &imageURLs
	}

	baseURL := strings.TrimSuffix(j.config.BaseURL, "/")
	itemID := *item.Id

	// Primary image (poster)
	if item.ImageTags != nil {
		if tag, ok := (item.ImageTags)["Primary"]; ok {
			imageURLs.Poster = fmt.Sprintf("%s/Items/%s/Images/Primary?tag=%s", baseURL, itemID, tag)
		}
	}

	// Backdrop image
	if item.BackdropImageTags != nil && len(item.BackdropImageTags) > 0 {
		imageURLs.Background = fmt.Sprintf("%s/Items/%s/Images/Backdrop?tag=%s", baseURL, itemID, item.BackdropImageTags[0])
	}

	// Other image types
	if item.ImageTags != nil {
		if tag, ok := (item.ImageTags)["Logo"]; ok {
			imageURLs.Logo = fmt.Sprintf("%s/Items/%s/Images/Logo?tag=%s", baseURL, itemID, tag)
		}

		if tag, ok := (item.ImageTags)["Thumb"]; ok {
			imageURLs.Thumbnail = fmt.Sprintf("%s/Items/%s/Images/Thumb?tag=%s", baseURL, itemID, tag)
		}

		if tag, ok := (item.ImageTags)["Banner"]; ok {
			imageURLs.Banner = fmt.Sprintf("%s/Items/%s/Images/Banner?tag=%s", baseURL, itemID, tag)
		}
	}

	return &imageURLs

}

// extractProviderIDs adds external IDs from the Jellyfin provider IDs map to the metadata
func extractProviderIDs(providerIds *map[string]string, externalIDs *types.ExternalIDs) {
	if providerIds == nil {
		return
	}

	// Common media identifier mappings
	idMappings := map[string]string{
		"Imdb":              "imdb",
		"Tmdb":              "tmdb",
		"Tvdb":              "tvdb",
		"MusicBrainzTrack":  "musicbrainz",
		"MusicBrainzAlbum":  "musicbrainz",
		"MusicBrainzArtist": "musicbrainz",
	}

	// Extract all available IDs based on the mappings
	for jellyfinKey, externalKey := range idMappings {
		if id, ok := (*providerIds)[jellyfinKey]; ok {
			externalIDs.AddOrUpdate(externalKey, id)
		}
	}

}

// Helper function to get duration seconds from ticks pointer
func getDurationFromTicks(ticks *int64) int64 {
	if ticks == nil {
		return 0
	}
	duration := time.Duration(*ticks/10000000) * time.Second
	return int64(duration.Seconds())
}
