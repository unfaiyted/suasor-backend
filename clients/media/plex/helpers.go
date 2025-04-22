package plex

import (
	"context"
	"fmt"
	"strings"
	media "suasor/clients/media"
	mediatypes "suasor/clients/media/types"
	"suasor/clients/types"
	"suasor/types/models"
	"suasor/utils/logger"

	"github.com/LukeHagar/plexgo/models/operations"
)

// Helper functions for item conversion using the factory pattern

func GetItemFromLibraryMetadata[T mediatypes.MediaData](
	ctx context.Context,
	client *PlexClient,
	item *operations.GetLibraryItemsMetadata,
) (T, error) {
	return media.ConvertTo[*PlexClient, *operations.GetLibraryItemsMetadata, T](
		client, ctx, item)
}

//tracks list GetMetadataChildrenMetadata

func GetItemFromMetadata[T mediatypes.MediaData](
	ctx context.Context,
	client *PlexClient,
	item *operations.GetMediaMetaDataMetadata,
) (T, error) {
	return media.ConvertTo[*PlexClient, *operations.GetMediaMetaDataMetadata, T](
		client, ctx, item)
}

func GetItemFromPlaylist[T mediatypes.MediaData](
	ctx context.Context,
	client *PlexClient,
	item *operations.GetPlaylistsMetadata,
) (T, error) {
	return media.ConvertTo[*PlexClient, *operations.GetPlaylistsMetadata, T](
		client, ctx, item)
}

func GetChildItem[T mediatypes.MediaData](
	ctx context.Context,
	client *PlexClient,
	item *operations.GetMetadataChildrenMetadata,
) (T, error) {
	return media.ConvertTo[*PlexClient, *operations.GetMetadataChildrenMetadata, T](
		client, ctx, item)
}

// func GetMediaItemFromPlaylist[T mediatypes.MediaData](
// 	ctx context.Context,
// 	client *PlexClient,
// 	item T,
// 	itemID string,
// ) (*models.MediaItem[T], error) {
// 	mediaItem := models.NewMediaItem[T](item.GetMediaType(), item)
// 	mediaItem.SetClientInfo(client.ClientID, client.ClientType, itemID)
//
// 	return mediaItem, nil
// }

func GetChildItemsList[T mediatypes.MediaData](
	ctx context.Context,
	client *PlexClient,
	items []operations.GetMetadataChildrenMetadata,
) ([]T, error) {

	mediaItems := make([]T, 0, len(items))

	for _, item := range items {
		itemT, err := GetChildItem[T](ctx, client, &item)
		if err != nil {
			return nil, err
		}
		mediaItems = append(mediaItems, itemT)
	}

	return mediaItems, nil

}

func GetMediaItem[T mediatypes.MediaData](
	ctx context.Context,
	client *PlexClient,
	item T,
	itemID string,
) (*models.MediaItem[T], error) {
	mediaItem := models.NewMediaItem[T](item.GetMediaType(), item)
	mediaItem.SetClientInfo(client.ClientID, client.ClientType, itemID)

	return mediaItem, nil
}

func GetChildMediaItem[T mediatypes.MediaData](
	ctx context.Context,
	client *PlexClient,
	item T,
	itemID string,
) (*models.MediaItem[T], error) {
	mediaItem := models.NewMediaItem[T](item.GetMediaType(), item)
	mediaItem.SetClientInfo(client.ClientID, client.ClientType, itemID)

	return mediaItem, nil
}

func GetChildMediaItemsList[T mediatypes.MediaData](
	ctx context.Context,
	client *PlexClient,
	items []operations.GetMetadataChildrenMetadata,
) ([]*models.MediaItem[T], error) {
	var mediaItems []*models.MediaItem[T]
	for _, item := range items {
		itemT, err := GetChildItem[T](ctx, client, &item)
		if err != nil {
			return nil, err
		}
		mediaItem, err := GetChildMediaItem[T](ctx, client, itemT, *item.RatingKey)
		if err != nil {
			return nil, err
		}
		mediaItem.SetClientInfo(client.ClientID, client.ClientType, *item.RatingKey)
		mediaItems = append(mediaItems, mediaItem)
	}

	return mediaItems, nil
}

func GetMediaItemList[T mediatypes.MediaData](
	ctx context.Context,
	client *PlexClient,
	items []operations.GetLibraryItemsMetadata,
) ([]*models.MediaItem[T], error) {
	var mediaItems []*models.MediaItem[T]
	for _, item := range items {
		itemT, err := GetItemFromLibraryMetadata[T](ctx, client, &item)
		if err != nil {
			return nil, err
		}
		mediaItem, err := GetMediaItem[T](ctx, client, itemT, item.RatingKey)
		if err != nil {
			return nil, err
		}
		mediaItem.SetClientInfo(client.ClientID, client.ClientType, item.RatingKey)
		mediaItems = append(mediaItems, mediaItem)
	}

	return mediaItems, nil
}

func GetMediaItemListFromPlaylist[T mediatypes.MediaData](
	ctx context.Context,
	client *PlexClient,
	items []operations.GetPlaylistsMetadata,
) ([]*models.MediaItem[T], error) {
	var mediaItems []*models.MediaItem[T]
	for _, item := range items {
		itemT, err := GetItemFromPlaylist[T](ctx, client, &item)
		if err != nil {
			return nil, err
		}
		mediaItem, err := GetMediaItem[T](ctx, client, itemT, *item.RatingKey)
		if err != nil {
			return nil, err
		}
		mediaItem.SetClientInfo(client.ClientID, client.ClientType, *item.RatingKey)
		mediaItems = append(mediaItems, mediaItem)
	}

	return mediaItems, nil
}

func (c *PlexClient) findLibrarySectionByType(ctx context.Context, sectionType string) (string, error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Debug().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("sectionType", sectionType).
		Msg("Finding library section by type")

	libraries, err := c.plexAPI.Library.GetAllLibraries(ctx)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("sectionType", sectionType).
			Msg("Failed to get libraries from Plex")
		return "", fmt.Errorf("failed to get libraries: %w", err)
	}

	log.Debug().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("libraryCount", len(libraries.Object.MediaContainer.Directory)).
		Msg("Retrieved libraries from Plex")

	for _, dir := range libraries.Object.MediaContainer.Directory {
		if dir.Type == sectionType {
			log.Debug().
				Uint64("clientID", c.ClientID).
				Str("clientType", string(c.ClientType)).
				Str("sectionType", sectionType).
				Str("sectionKey", dir.Key).
				Str("sectionTitle", dir.Title).
				Msg("Found matching library section")
			return dir.Key, nil
		}
	}

	log.Debug().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("sectionType", sectionType).
		Msg("No matching library section found")

	return "", nil
}

// makeFullURL creates a complete URL from a resource path
func (c *PlexClient) makeFullURL(resourcePath string) string {
	if resourcePath == "" {
		return ""
	}

	plexConfig := c.Config.(*types.PlexConfig)

	if strings.HasPrefix(resourcePath, "http") {
		return resourcePath
	}

	return fmt.Sprintf("%s%s", plexConfig.BaseURL, resourcePath)
}
