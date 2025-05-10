package plex

import (
	"context"
	"fmt"
	"strings"
	media "suasor/clients/media"
	mediatypes "suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/logger"
	"time"

	"github.com/unfaiyted/plexgo/models/operations"
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

func GetItemFromPlaylistMetadata[T mediatypes.MediaData](
	ctx context.Context,
	client *PlexClient,
	item *operations.GetPlaylistMetadata,
) (T, error) {
	return media.ConvertTo[*PlexClient, *operations.GetPlaylistMetadata, T](
		client, ctx, item)
}

func GetItemFromPlaylistCreate[T mediatypes.MediaData](
	ctx context.Context,
	client *PlexClient,
	item *operations.CreatePlaylistMetadata,
) (T, error) {
	newItem := mapPlaylistCreateMetadataToPlaylistMetadata(item)
	return GetItemFromPlaylistMetadata[T](ctx, client, newItem)

}

func GetItemFromPlaylistContents[T mediatypes.MediaData](
	ctx context.Context,
	client *PlexClient,
	item *operations.GetPlaylistContentsMetadata,
) (T, error) {

	newItem := mapPlaylistContentsMetadataToChildrenMetadata(item)

	return media.ConvertTo[*PlexClient, *operations.GetMetadataChildrenMetadata, T](
		client, ctx, newItem)
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
	mediaItem.SetClientInfo(client.GetClientID(), client.GetClientType(), itemID)

	return mediaItem, nil
}

func GetChildMediaItem[T mediatypes.MediaData](
	ctx context.Context,
	client *PlexClient,
	item T,
	itemID string,
) (*models.MediaItem[T], error) {
	mediaItem := models.NewMediaItem[T](item.GetMediaType(), item)
	mediaItem.SetClientInfo(client.GetClientID(), client.GetClientType(), itemID)

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
		mediaItem.SetClientInfo(client.GetClientID(), client.GetClientType(), *item.RatingKey)
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
		mediaItem.SetClientInfo(client.GetClientID(), client.GetClientType(), item.RatingKey)
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
		mediaItem.SetClientInfo(client.GetClientID(), client.GetClientType(), *item.RatingKey)
		mediaItems = append(mediaItems, mediaItem)
	}

	return mediaItems, nil
}

func (c *PlexClient) findLibrarySectionByType(ctx context.Context, sectionType operations.GetAllLibrariesType) (string, error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Debug().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("sectionType", string(sectionType)).
		Msg("Finding library section by type")

	libraries, err := c.plexAPI.Library.GetAllLibraries(ctx)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Str("sectionType", string(sectionType)).
			Msg("Failed to get libraries from Plex")
		return "", fmt.Errorf("failed to get libraries: %w", err)
	}

	log.Debug().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Int("libraryCount", len(libraries.Object.MediaContainer.Directory)).
		Msg("Retrieved libraries from Plex")

	for _, dir := range libraries.Object.MediaContainer.Directory {
		if dir.Type == sectionType {
			log.Debug().
				Uint64("clientID", c.GetClientID()).
				Str("clientType", string(c.GetClientType())).
				Str("sectionType", string(sectionType)).
				Str("sectionKey", dir.Key).
				Str("sectionTitle", dir.Title).
				Msg("Found matching library section")
			return dir.Key, nil
		}
	}

	log.Debug().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("sectionType", string(sectionType)).
		Msg("No matching library section found")

	return "", nil
}

// makeFullURL creates a complete URL from a resource path
func (c *PlexClient) makeFullURL(resourcePath string) string {
	if resourcePath == "" {
		return ""
	}

	if strings.HasPrefix(resourcePath, "http") {
		return resourcePath
	}

	return fmt.Sprintf("%s%s", c.plexConfig().GetBaseURL(), resourcePath)
}

func mapPlaylistContentsMetadataToChildrenMetadata(item *operations.GetPlaylistContentsMetadata) *operations.GetMetadataChildrenMetadata {
	newItem := &operations.GetMetadataChildrenMetadata{
		RatingKey: item.RatingKey,
		Key:       item.Key,
		GUID:      item.GUID,
		Type:      item.Type,
		Title:     item.Title,
		Summary:   item.Summary,
		Thumb:     item.Thumb,
		Art:       item.Art,
		AddedAt:   item.AddedAt,
		UpdatedAt: item.UpdatedAt,
	}

	return newItem
}
func mapPlaylistCreateMetadataToPlaylistMetadata(item *operations.CreatePlaylistMetadata) *operations.GetPlaylistMetadata {
	newItem := &operations.GetPlaylistMetadata{
		RatingKey:    item.RatingKey,
		Key:          item.Key,
		GUID:         item.GUID,
		Type:         item.Type,
		Title:        item.Title,
		Summary:      item.Summary,
		AddedAt:      item.AddedAt,
		UpdatedAt:    item.UpdatedAt,
		PlaylistType: item.PlaylistType,
		Smart:        item.Smart,
		Icon:         item.Icon,
		Duration:     item.Duration,
		LeafCount:    item.LeafCount,
		Composite:    item.Composite,
	}

	return newItem
}

func setMediaListToItemList(list *models.MediaItemList, itemList *mediatypes.ItemList) {

	// Loop over all the items in the MediaItemList
	list.ForEach(func(uuid string, mediaType mediatypes.MediaType, item any) bool {
		// Cast the item to the correct type
		itemT, ok := item.(*models.MediaItem[mediatypes.MediaData])
		if !ok {
			return true
		}
		// Find the corresponding item in the ItemList
		foundItem, index, found := itemList.FindItemByID(itemT.ID)
		if found {
			// Update the item in the ItemList
			itemList.Items[index].Position = foundItem.Position
			itemList.Items[index].LastChanged = foundItem.LastChanged
			itemList.Items[index].ChangeHistory = foundItem.ChangeHistory
		} else {
			// Add the item to the ItemList
			itemList.AddItem(mediatypes.ListItem{
				ItemID:        itemT.ID,
				Position:      0,
				LastChanged:   time.Now(),
				ChangeHistory: []mediatypes.ChangeRecord{},
			})
		}
		return true
	})

}
