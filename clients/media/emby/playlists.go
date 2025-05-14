// playlists.go
package emby

import (
	"context"
	"fmt"
	"strings"
	"suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/logger"

	"github.com/antihax/optional"

	embyclient "suasor/internal/clients/embyAPI"
)

// SupportsPlaylists returns true since Emby supports playlists
func (e *EmbyClient) SupportsPlaylists() bool {
	return true
}

// Search retrieves playlists from the Emby server matching the query options
// Implements ListProvider[*types.Playlist] interface method
func (e *EmbyClient) SearchPlaylists(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Playlist], error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.GetClientID()).
		Str("clientType", string(e.GetClientType())).
		Msg("Searching playlists from Emby server")

	queryParams := embyclient.ItemsServiceApiGetItemsOpts{
		IncludeItemTypes: optional.NewString("Playlist"),
		Recursive:        optional.NewBool(true),
	}

	ApplyClientQueryOptions(ctx, &queryParams, options)

	// Get user ID
	userID := e.getUserID()
	if userID != "" {
		queryParams.UserId = optional.NewString(userID)
	}

	items, resp, err := e.client.ItemsServiceApi.GetItems(ctx, &queryParams)
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", e.embyConfig().GetBaseURL()).
			Str("apiEndpoint", "/Items").
			Msg("Failed to fetch playlists from Emby")
		return nil, fmt.Errorf("failed to fetch playlists: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(items.Items)).
		Int("totalRecordCount", int(items.TotalRecordCount)).
		Msg("Successfully retrieved playlists from Emby")

	playlists := make([]*models.MediaItem[*types.Playlist], 0)
	for _, item := range items.Items {
		if item.Type_ == "Playlist" {
			itemPlaylist, err := GetItem[*types.Playlist](ctx, e, &item)
			playlist, err := GetMediaItem[*types.Playlist](ctx, e, itemPlaylist, item.Id)
			if err != nil {
				log.Warn().
					Err(err).
					Str("playlistID", item.Id).
					Str("playlistName", item.Name).
					Msg("Error converting Emby item to playlist format")
				continue
			}
			playlists = append(playlists, playlist)
		}
	}

	log.Info().
		Int("playlistsReturned", len(playlists)).
		Msg("Completed Search playlists request")

	return playlists, nil
}

func (e *EmbyClient) SearchPlaylistItems(ctx context.Context, playlistID string, options *types.QueryOptions) ([]*models.MediaItem[*types.Playlist], error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Str("playlistID", playlistID).
		Uint64("clientID", e.GetClientID()).
		Msg("Searching playlist items in Emby")
	// Get user ID
	userID := e.getUserID()
	if userID == "" {
		log.Error().Msg("User ID is required for Emby queries but was not provided or resolved")
		return nil, fmt.Errorf("failed to search playlist items: missing user ID")
	}
	// Query for playlist items
	queryParams := embyclient.ItemsServiceApiGetItemsOpts{
		ParentId:  optional.NewString(playlistID),
		UserId:    optional.NewString(userID),
		Recursive: optional.NewBool(false),
	}
	ApplyClientQueryOptions(ctx, &queryParams, options)
	// Make the API call
	response, resp, err := e.client.ItemsServiceApi.GetItems(ctx, &queryParams)
	if err != nil {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Msg("Failed to fetch playlist items from Emby")
		return nil, fmt.Errorf("failed to fetch playlist items: %w", err)
	}
	log.Info().Int("statusCode", resp.StatusCode).
		Int("itemCount", len(response.Items)).
		Msg("Successfully retrieved playlist items from Emby")
	// Process each item
	items := make([]*models.MediaItem[*types.Playlist], 0, len(response.Items))
	for _, item := range response.Items {
		// Convert to playlist item
		playlistItem, err := GetItem[*types.Playlist](ctx, e, &item)
		if err != nil {
			log.Warn().
				Err(err).
				Str("itemID", item.Id).
				Str("itemName", item.Name).
				Msg("Error converting Emby item to playlist format")
			continue
		}
		mediaItem, err := GetMediaItem[*types.Playlist](ctx, e, playlistItem, item.Id)
		if err != nil {
			log.Warn().
				Err(err).
				Str("itemID", item.Id).
				Str("itemName", item.Name).
				Msg("Error creating media item for playlist item")
			continue
		}
		items = append(items, mediaItem)
	}
	log.Info().
		Int("itemsReturned", len(items)).
		Str("playlistID", playlistID).
		Msg("Completed getting items from playlist")
	return items, nil
}

func (e *EmbyClient) GetPlaylist(ctx context.Context, playlistID string) (*models.MediaItem[*types.Playlist], error) {
	return e.GetPlaylistByID(ctx, playlistID)
}

func (e *EmbyClient) GetPlaylistByID(ctx context.Context, playlistID string) (*models.MediaItem[*types.Playlist], error) {
	// opts := embclient.ItemsServiceApiGetItemsOpts{
	// 	Ids:    optional.NewString(collectionID),
	// 	Fields: optional.NewString("PrimaryImageAspectRatio,BasicSyncInfo,CanDelete,Container,DateCreated,PremiereDate,Genres,MediaSources,Overview,ParentId,Path,SortName,Studios,Taglines"),
	// }
	opts := types.QueryOptions{
		Limit:   1,
		ItemIDs: playlistID,
	}
	playlists, err := e.SearchPlaylists(ctx, &opts)
	if err != nil {
		return nil, err
	}
	if len(playlists) == 0 {
		return nil, fmt.Errorf("collection not found")
	}
	playlist := playlists[0]

	return playlist, err
}

// GetItems retrieves the items in a specific playlist
// Implements ListProvider[*types.Playlist] interface method
func (e *EmbyClient) GetPlaylistItems(ctx context.Context, playlistID string) (*models.MediaItemList[*types.Playlist], error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Str("playlistID", playlistID).
		Uint64("clientID", e.GetClientID()).
		Msg("Getting items from Emby playlist")

	// Get user ID
	userID := e.getUserID()
	if userID == "" {
		log.Error().Msg("User ID is required for Emby queries but was not provided or resolved")
		return nil, fmt.Errorf("failed to get playlist items: missing user ID")
	}

	// Query for playlist items
	queryParams := embyclient.ItemsServiceApiGetItemsOpts{
		ParentId:  optional.NewString(playlistID),
		UserId:    optional.NewString(userID),
		Recursive: optional.NewBool(false),
	}

	// Make the API call
	response, resp, err := e.client.ItemsServiceApi.GetItems(ctx, &queryParams)
	if err != nil {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Msg("Failed to get playlist items from Emby")
		return nil, fmt.Errorf("failed to get playlist items: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("itemCount", len(response.Items)).
		Msg("Successfully retrieved playlist items from Emby")

	// Convert string ID to uint64 for media item list
	var listIDUint uint64 = 0
	// Don't worry about conversion errors, we'll use 0 as default

	playlist, err := e.GetPlaylistByID(ctx, playlistID)
	if err != nil {
		return nil, err
	}
	// Create new media item list
	itemList := models.NewMediaItemList[*types.Playlist](playlist, listIDUint, 0)

	// Initialize the maps
	itemList.Playlists = make(map[string]*models.MediaItem[*types.Playlist])

	// Process each item
	for _, item := range response.Items {
		// Convert to playlist item
		playlistItem, err := GetItem[*types.Playlist](ctx, e, &item)
		if err != nil {
			log.Warn().
				Err(err).
				Str("itemID", item.Id).
				Str("itemName", item.Name).
				Msg("Error converting Emby item to playlist format")
			continue
		}

		mediaItem, err := GetMediaItem[*types.Playlist](ctx, e, playlistItem, item.Id)
		if err != nil {
			log.Warn().
				Err(err).
				Str("itemID", item.Id).
				Str("itemName", item.Name).
				Msg("Error creating media item for playlist item")
			continue
		}

		itemList.AddPlaylist(mediaItem)
	}

	log.Info().
		Int("itemsReturned", itemList.TotalItems).
		Str("playlistID", playlistID).
		Msg("Completed getting items from playlist")

	return itemList, nil
}

// Create creates a new playlist in Emby
// Implements ListProvider[*types.Playlist] interface method
func (e *EmbyClient) CreatePlaylist(ctx context.Context, name string, description string) (*models.MediaItem[*types.Playlist], error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Str("name", name).
		Uint64("clientID", e.GetClientID()).
		Msg("Creating new playlist in Emby")

	// Get user ID
	userID := e.getUserID()
	if userID == "" {
		log.Error().Msg("User ID is required for Emby queries but was not provided or resolved")
		return nil, fmt.Errorf("failed to create playlist: missing user ID")
	}

	// Create the playlist
	opts := embyclient.PlaylistServiceApiPostPlaylistsOpts{
		Name:      optional.NewString(name),
		MediaType: optional.NewString("Mixed"),
	}

	playlist, resp, err := e.client.PlaylistServiceApi.PostPlaylists(ctx, &opts)
	if err != nil {
		log.Error().
			Err(err).
			Str("name", name).
			Msg("Failed to create playlist in Emby")
		return nil, fmt.Errorf("failed to create playlist: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Str("playlistID", playlist.Id).
		Str("playlistName", playlist.Name).
		Msg("Successfully created playlist in Emby")

	// Convert to our playlist format
	playlistItem, err := GetPlaylistItem(ctx, e, &playlist)
	if err != nil {
		return nil, fmt.Errorf("error converting created playlist: %w", err)
	}

	// Create a media item for the playlist
	mediaItem, err := GetMediaItem[*types.Playlist](ctx, e, playlistItem, playlist.Id)
	if err != nil {
		return nil, fmt.Errorf("error creating media item for playlist: %w", err)
	}

	return mediaItem, nil
}

// CreatePlaylistWithItems creates a new playlist with items in Emby
// Implements ListProvider[*types.Playlist] interface method - THIS WAS THE MISSING METHOD
func (e *EmbyClient) CreatePlaylistWithItems(ctx context.Context, name string, description string, itemIDs []string) (*models.MediaItem[*types.Playlist], error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Str("name", name).
		Uint64("clientID", e.GetClientID()).
		Msg("Creating new playlist with items in Emby")

	// First create an empty playlist
	playlist, err := e.CreatePlaylist(ctx, name, description)
	if err != nil {
		return nil, fmt.Errorf("failed to create base playlist: %w", err)
	}

	// If we have items to add, add them
	if len(itemIDs) > 0 {
		// Get the playlistID from the SyncClients field which contains the client's ID for this item
		playlistID := ""
		if playlist.SyncClients != nil {
			playlistID = playlist.SyncClients.GetClientItemID(e.GetClientID())
		}

		// If we can't get it from SyncClients, try using the UUID instead
		if playlistID == "" {
			playlistID = playlist.UUID
		}

		err = e.AddPlaylistItems(ctx, playlistID, itemIDs)
		if err != nil {
			// If we fail to add items, still return the playlist but log the error
			log.Error().
				Err(err).
				Str("playlistID", playlistID).
				Strs("itemIDs", itemIDs).
				Msg("Failed to add items to newly created playlist")
		}
	}

	return playlist, nil
}

// Update updates an existing playlist in Emby
// Implements ListProvider[*types.Playlist] interface method
func (e *EmbyClient) UpdatePlaylist(ctx context.Context, playlistID string, name string, description string) (*models.MediaItem[*types.Playlist], error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Str("playlistID", playlistID).
		Str("name", name).
		Uint64("clientID", e.GetClientID()).
		Msg("Updating playlist in Emby")

	// Get user ID
	userID := e.getUserID()
	if userID == "" {
		log.Error().Msg("User ID is required for Emby queries but was not provided or resolved")
		return nil, fmt.Errorf("failed to update playlist: missing user ID")
	}

	// QueryResultBaseItemDto
	// First get the existing playlist
	existingPlaylists, resp, err := e.client.ItemsServiceApi.GetItems(ctx,
		&embyclient.ItemsServiceApiGetItemsOpts{
			UserId: optional.NewString(userID),
			Ids:    optional.NewString(playlistID),
		})
	existingPlaylist := existingPlaylists.Items[0]
	if err != nil {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Msg("Failed to get existing playlist from Emby")
		return nil, fmt.Errorf("failed to get existing playlist: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Str("playlistID", existingPlaylist.Id).
		Msg("Retrieved existing playlist from Emby")

	// Update the playlist - Emby doesn't have a dedicated update playlist endpoint
	// We'd need to use the item update endpoint with the updated data
	updateBody := embyclient.BaseItemDto{
		Id:       playlistID,
		Name:     name,
		Overview: description,
	}

	_, err = e.client.ItemUpdateServiceApi.PostItemsByItemid(ctx, updateBody, playlistID)
	if err != nil {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Msg("Failed to update playlist in Emby")
		return nil, fmt.Errorf("failed to update playlist: %w", err)
	}

	// Get the updated playlist
	updatedPlaylists, resp, err := e.client.ItemsServiceApi.GetItems(ctx, &embyclient.ItemsServiceApiGetItemsOpts{
		UserId: optional.NewString(userID),
		Ids:    optional.NewString(playlistID),
	})
	updatedPlaylist := updatedPlaylists.Items[0]
	if err != nil {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Msg("Failed to get updated playlist from Emby")
		return nil, fmt.Errorf("failed to get updated playlist: %w", err)
	}

	// Convert to our playlist format
	playlistItem, err := GetItem[*types.Playlist](ctx, e, &updatedPlaylist)
	if err != nil {
		return nil, fmt.Errorf("error converting updated playlist: %w", err)
	}

	// Create a media item for the playlist
	mediaItem, err := GetMediaItem[*types.Playlist](ctx, e, playlistItem, updatedPlaylist.Id)
	if err != nil {
		return nil, fmt.Errorf("error creating media item for updated playlist: %w", err)
	}

	log.Info().
		Str("playlistID", playlistID).
		Msg("Successfully updated playlist in Emby")

	return mediaItem, nil
}

// Delete removes a playlist from Emby
// Implements ListProvider[*types.Playlist] interface method
func (e *EmbyClient) DeletePlaylist(ctx context.Context, playlistID string) error {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Str("playlistID", playlistID).
		Uint64("clientID", e.GetClientID()).
		Msg("Deleting playlist from Emby")

	entryIDs := ""

	// TODO: Implement playlist deletion for Emby, the API doesn't support it, but we should
	// be able to delete all of the items in the playlist and delete the playlist itself.
	// check old python code for example on deleting playlists

	// Delete the item (playlist)
	resp, err := e.client.PlaylistServiceApi.DeletePlaylistsByIdItems(ctx, playlistID, entryIDs)
	if err != nil {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Msg("Failed to delete playlist from Emby")
		return fmt.Errorf("failed to delete playlist: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Str("playlistID", playlistID).
		Msg("Successfully deleted playlist from Emby")

	return nil
}

// AddPlaylistItem adds an item to a playlist in Emby
// Implementation detail - called by AddItemPlaylist to match the PlaylistProvider interface
func (e *EmbyClient) AddPlaylistItem(ctx context.Context, playlistID string, itemID string) error {
	return e.AddPlaylistItems(ctx, playlistID, []string{itemID})
}

// AddPlaylistItems adds items to a playlist in Emby
// Implements PlaylistProvider interface method
func (e *EmbyClient) AddPlaylistItems(ctx context.Context, playlistID string, itemIDs []string) error {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Str("playlistID", playlistID).
		Strs("itemIDs", itemIDs).
		Uint64("clientID", e.GetClientID()).
		Msg("Adding item to playlist in Emby")

	// Get user ID
	userID := e.getUserID()
	if userID == "" {
		log.Error().Msg("User ID is required for Emby queries but was not provided or resolved")
		return fmt.Errorf("failed to add item to playlist: missing user ID")
	}

	// If there are no items to add, return immediately without error
	if len(itemIDs) == 0 {
		log.Info().
			Str("playlistID", playlistID).
			Msg("No items to add to playlist in Emby")
		return nil
	}

	// Use the Emby API to add items to playlist
	opts := embyclient.PlaylistServiceApiPostPlaylistsByIdItemsOpts{
		UserId: optional.NewString(userID),
	}

	strItemIDs := strings.Join(itemIDs, ",")

	result, resp, err := e.client.PlaylistServiceApi.PostPlaylistsByIdItems(ctx, strItemIDs, playlistID, &opts)
	if err != nil {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Strs("itemIDs", itemIDs).
			Msg("Failed to add item to playlist in Emby")
		return fmt.Errorf("failed to add item to playlist: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Str("playlistID", playlistID).
		Strs("itemIDs", itemIDs).
		Msg("Successfully added item to playlist in Emby")

	if len(itemIDs) > 0 && result.ItemAddedCount == 0 {
		return fmt.Errorf("failed to add item to playlist: no items added")
	}

	return nil
}

// Implements ListProvider[*types.Playlist] interface method
func (e *EmbyClient) RemovePlaylistItem(ctx context.Context, playlistID string, itemID string) error {
	return e.RemovePlaylistItems(ctx, playlistID, []string{itemID})
}

func (e *EmbyClient) RemovePlaylistItems(ctx context.Context, playlistID string, itemIDs []string) error {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Str("playlistID", playlistID).
		Strs("itemIDs", itemIDs).
		Uint64("clientID", e.GetClientID()).
		Msg("Removing item from playlist in Emby")

	// Get user ID
	userID := e.getUserID()
	if userID == "" {
		log.Error().Msg("User ID is required for Emby queries but was not provided or resolved")
		return fmt.Errorf("failed to remove item from playlist: missing user ID")
	}
	strItemIDs := strings.Join(itemIDs, ",")

	resp, err := e.client.PlaylistServiceApi.DeletePlaylistsByIdItems(ctx, playlistID, strItemIDs)
	if err != nil {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Strs("itemIDs", itemIDs).
			Msg("Failed to remove items from playlist in Emby")
		return fmt.Errorf("failed to remove item from playlist: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Str("playlistID", playlistID).
		Strs("itemIDs", itemIDs).
		Msg("Successfully removed items from playlist in Emby")

	return nil
}

func (e *EmbyClient) RemoveAllPlaylistItems(ctx context.Context, playlistID string) error {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Str("playlistID", playlistID).
		Uint64("clientID", e.GetClientID()).
		Msg("Removing item from playlist in Emby")

	// Get user ID
	userID := e.getUserID()
	if userID == "" {
		log.Error().Msg("User ID is required for Emby queries but was not provided or resolved")
		return fmt.Errorf("failed to remove item from playlist: missing user ID")
	}

	// Use the Emby API to remove item from playlist
	// We need to get the position of the item in the playlist first
	queryParams := embyclient.ItemsServiceApiGetItemsOpts{
		ParentId:  optional.NewString(playlistID),
		UserId:    optional.NewString(userID),
		Recursive: optional.NewBool(false),
	}

	// Get all items in the playlist
	items, _, err := e.client.ItemsServiceApi.GetItems(ctx, &queryParams)
	if err != nil {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Msg("Failed to get playlist items for removal")
		return fmt.Errorf("failed to get playlist items for removal: %w", err)
	}
	itemIDs := make([]string, 0, len(items.Items))
	for _, item := range items.Items {
		itemIDs = append(itemIDs, item.Id)
	}

	strItemIDs := strings.Join(itemIDs, ",")

	resp, err := e.client.PlaylistServiceApi.DeletePlaylistsByIdItems(ctx, playlistID, strItemIDs)
	if err != nil {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Strs("itemIDs", itemIDs).
			Msg("Failed to remove item from playlist in Emby")
		return fmt.Errorf("failed to remove item from playlist: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Str("playlistID", playlistID).
		Strs("itemIDs", itemIDs).
		Msg("Successfully removed item from playlist in Emby")

	return nil
}

// ReorderItems reorders items in a playlist in Emby
// Implements ListProvider[*types.Playlist] interface method
func (e *EmbyClient) ReorderPlaylistItems(ctx context.Context, playlistID string, itemIDs []string) error {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Str("playlistID", playlistID).
		Uint64("clientID", e.GetClientID()).
		Msg("Reordering items in playlist in Emby")

	// Get user ID
	userID := e.getUserID()
	if userID == "" {
		log.Error().Msg("User ID is required for Emby queries but was not provided or resolved")
		return fmt.Errorf("failed to reorder playlist items: missing user ID")
	}

	// Emby API doesn't have a direct method to reorder playlist items
	// We would need to get the original playlist, remove all items, then add them in the desired order

	// Get the original playlist
	originalPlaylists, resp, err := e.client.ItemsServiceApi.GetItems(ctx, &embyclient.ItemsServiceApiGetItemsOpts{
		Ids: optional.NewString(playlistID),
	})
	originalPlaylist := originalPlaylists.Items[0]
	if err != nil {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Msg("Failed to get original playlist from Emby")
		return fmt.Errorf("failed to get original playlist: %w", err)
	}

	opts := embyclient.PlaylistServiceApiGetPlaylistsByIdItemsOpts{
		UserId: optional.NewString(userID),
	}

	// get playlist items
	result, resp, err := e.client.PlaylistServiceApi.GetPlaylistsByIdItems(ctx, originalPlaylist.Id, &opts)
	playlistItems := result.Items

	if err != nil {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Msg("Failed to get playlist items from Emby")
		return fmt.Errorf("failed to get playlist items: %w", err)
	}

	entryIDs := ""
	for _, item := range playlistItems {
		entryIDs += item.Id
		entryIDs += ","
	}

	// Delete the original playlist items
	_, err = e.client.PlaylistServiceApi.DeletePlaylistsByIdItems(ctx, playlistID, originalPlaylist.Id)
	if err != nil {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Msg("Failed to delete original playlist items from Emby")
		return fmt.Errorf("failed to delete original playlist items: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Str("playlistName", originalPlaylist.Name).
		Msg("Retrieved original playlist from Emby")

	// As a workaround to reorder the playlist, we can:
	// 1. Clear the playlist
	// 2. Add items back in the desired order

	// First, get all current items to remove them
	queryParams := embyclient.ItemsServiceApiGetItemsOpts{
		ParentId:  optional.NewString(playlistID),
		UserId:    optional.NewString(userID),
		Recursive: optional.NewBool(false),
	}

	items, _, err := e.client.ItemsServiceApi.GetItems(ctx, &queryParams)
	if err != nil {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Msg("Failed to get playlist items for reordering")
		return fmt.Errorf("failed to get playlist items for reordering: %w", err)
	}

	// Remove all items
	entryIds := ""
	for i, item := range items.Items {
		if i > 0 {
			entryIds += ","
		}
		entryIds += item.Id
	}

	_, err = e.client.PlaylistServiceApi.DeletePlaylistsByIdItems(ctx, playlistID, entryIds)
	if err != nil {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Msg("Failed to clear playlist for reordering")
		return fmt.Errorf("failed to clear playlist for reordering: %w", err)
	}

	addItemIDs := ""
	// Add items back in the desired order
	for _, itemID := range itemIDs {
		addItemIDs += itemID
		addItemIDs += ","
	}
	updateOpts := embyclient.PlaylistServiceApiPostPlaylistsByIdItemsOpts{
		UserId: optional.NewString(userID),
	}
	_, respItems, err := e.client.PlaylistServiceApi.PostPlaylistsByIdItems(ctx, addItemIDs, playlistID, &updateOpts)
	if err != nil {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Strs("itemIDs", itemIDs).
			Msg("Failed to add item to playlist during reordering")
		return fmt.Errorf("failed to add item during reordering: %w", err)
	}
	if respItems.StatusCode != 200 {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Strs("itemIDs", itemIDs).
			Msg("Failed to add item to playlist during reordering")
		return fmt.Errorf("failed to add item during reordering: %w", err)
	}

	log.Info().
		Str("playlistID", playlistID).
		Int("itemCount", len(itemIDs)).
		Msg("Successfully reordered items in playlist in Emby")

	return nil
}
