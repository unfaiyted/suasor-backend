package jellyfin

import (
	"context"
	"fmt"
	"time"

	"suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/logger"

	clienttypes "suasor/clients/types"

	jellyfin "github.com/sj14/jellyfin-go/api"
)

// SupportsPlaylists indicates if this client supports playlists
func (j *JellyfinClient) SupportsPlaylists() bool {
	return true
}

// GetPlaylists retrieves playlists from Jellyfin
func (j *JellyfinClient) GetPlaylists(ctx context.Context, options *types.QueryOptions) ([]models.MediaItem[*types.Playlist], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.GetClientID()).
		Str("clientType", string(j.GetClientType())).
		Msg("Retrieving playlists from Jellyfin")

	includeItemTypes := []jellyfin.BaseItemKind{jellyfin.BASEITEMKIND_PLAYLIST}

	// Construct filter string for playlists
	itemsReq := j.client.ItemsAPI.GetItems(ctx).
		IncludeItemTypes(includeItemTypes).
		Recursive(true)

	NewJellyfinQueryOptions(ctx, options).
		SetItemsRequest(ctx, &itemsReq)

	// Get playlists from Jellyfin
	response, _, err := itemsReq.Execute()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to fetch playlists from Jellyfin")
		return nil, err
	}

	if response.Items == nil || len(response.Items) == 0 {
		log.Info().Msg("No playlists returned from Jellyfin")
		return []models.MediaItem[*types.Playlist]{}, nil
	}

	// Convert Jellyfin items to playlist models
	playlists := make([]models.MediaItem[*types.Playlist], 0, len(response.Items))

	for _, item := range response.Items {
		// Safely handle name/title
		title := ""
		if item.Name.IsSet() {
			title = *item.Name.Get()
		}

		// Safely handle description
		description := ""
		if item.Overview.IsSet() {
			description = *item.Overview.Get()
		}

		// Safely handle item count
		itemCount := 0
		if item.ChildCount.IsSet() {
			itemCount = int(*item.ChildCount.Get())
		}

		// Convert to our playlist model
		playlist := *models.NewMediaItem[*types.Playlist](&types.Playlist{
			ItemList: types.ItemList{
				Details: &types.MediaDetails{
					Title:       title,
					Description: description,
					Artwork:     *j.getArtworkURLs(&item),
				},
				ItemCount: itemCount,
				IsPublic:  true, // Assume public by default in Jellyfin
			},
		})
		playlist.SetClientInfo(j.GetClientID(), j.GetClientType(), *item.Id)
		playlists = append(playlists, playlist)
	}

	log.Info().
		Int("playlistCount", len(playlists)).
		Msg("Successfully retrieved playlists from Jellyfin")

	return playlists, nil
}

// Getplaylist retrieves a playlist from Jellyfin
func (j *JellyfinClient) GetPlaylist(ctx context.Context, playlistID string) *models.MediaItem[*types.Playlist] {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.GetClientID()).
		Str("clientType", string(j.GetClientType())).
		Str("playlistID", playlistID).
		Msg("Retrieving playlist from Jellyfin")

	// Get playlist from Jellyfin
	getReq := j.client.PlaylistsAPI.GetPlaylist(ctx, playlistID)
	getReq.Execute()

	return nil
}

// GetPlaylistItems retrieves items in a playlist from Jellyfin
func (j *JellyfinClient) GetPlaylistItems(ctx context.Context, playlistID string, options *types.QueryOptions) (*models.MediaItemList[*types.Playlist], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.GetClientID()).
		Str("clientType", string(j.GetClientType())).
		Str("playlistID", playlistID).
		Msg("Retrieving playlist items from Jellyfin")

	playlistDetails := j.GetPlaylist(ctx, playlistID)

	playlist := models.NewMediaItemList[*types.Playlist](playlistDetails, 0, 0)

	playlistRes := j.client.PlaylistsAPI.GetPlaylistItems(ctx, playlistID)
	jellyfinPlaylist, _, err := playlistRes.Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Msg("Failed to fetch playlist items from Jellyfin")
		return nil, err
	}

	if jellyfinPlaylist.Items == nil || len(jellyfinPlaylist.Items) == 0 {
		log.Info().
			Str("playlistID", playlistID).
			Msg("No items found in playlist")
		return nil, nil
	}

	// Convert Jellyfin items to models
	mediaResults, err := GetMixedMediaItems(j, ctx, jellyfinPlaylist.Items)
	if err != nil {
		return nil, err
	}

	for _, item := range jellyfinPlaylist.Items {
		if item.Type == jellyfin.BASEITEMKIND_MOVIE.Ptr() {
			movieItem, err := GetItem[*types.Movie](ctx, j, &item)
			if err != nil {
				return nil, err
			}
			movie, err := GetMediaItem[*types.Movie](ctx, j, movieItem, *item.Id)
			playlist.Items.AddMovie(movie)
		} else if item.Type == jellyfin.BASEITEMKIND_EPISODE.Ptr() {
			episodeItem, err := GetItem[*types.Episode](ctx, j, &item)
			if err != nil {
				return nil, err
			}
			episode, err := GetMediaItem[*types.Episode](ctx, j, episodeItem, *item.Id)
			if err != nil {
				return nil, err
			}
			playlist.Items.AddEpisode(episode)
		} else if item.Type == jellyfin.BASEITEMKIND_AUDIO.Ptr() {
			trackItem, err := GetItem[*types.Track](ctx, j, &item)
			if err != nil {
				return nil, err
			}
			track, err := GetMediaItem[*types.Track](ctx, j, trackItem, *item.Id)
			if err != nil {
				return nil, err
			}
			playlist.Items.AddTrack(track)
		} else if item.Type == jellyfin.BASEITEMKIND_PLAYLIST.Ptr() {
			// playlistItem, err := GetItem[*types.Playlist](ctx, j, &item)
			// if err != nil {
			// 	return nil, err
			// }
			// playlist, err := GetMediaItem[*types.Playlist](ctx, j, playlistItem, *item.Id)
			// if err != nil {
			// 	return nil, err
			// }
			// playlist.AddPlaylist(playlist)
		} else if item.Type == jellyfin.BASEITEMKIND_SERIES.Ptr() {
			seriesItem, err := GetItem[*types.Series](ctx, j, &item)
			if err != nil {
				return nil, err
			}
			series, err := GetMediaItem[*types.Series](ctx, j, seriesItem, *item.Id)
			if err != nil {
				return nil, err
			}
			playlist.Items.AddSeries(series)
		} else if item.Type == jellyfin.BASEITEMKIND_SEASON.Ptr() {
			seasonItem, err := GetItem[*types.Season](ctx, j, &item)
			if err != nil {
				return nil, err
			}
			season, err := GetMediaItem[*types.Season](ctx, j, seasonItem, *item.Id)
			if err != nil {
				return nil, err
			}
			playlist.Items.AddSeason(season)
		} else if item.Type == jellyfin.BASEITEMKIND_COLLECTION_FOLDER.Ptr() {
			// collection, err := GetItem[*types.Collection](ctx, j, &item)
			// if err != nil {
			// 	return nil, err
			// }
			// playlist.AddCollection(collection)

		}

	}

	log.Info().
		Str("playlistID", playlistID).
		Int("itemCount", mediaResults.Len()).
		Msg("Successfully retrieved playlist items from Jellyfin")

	return playlist, nil
}

// CreatePlaylist creates a new playlist in Jellyfin
func (j *JellyfinClient) CreatePlaylist(ctx context.Context, name string, description string) (*models.MediaItem[*types.Playlist], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.GetClientID()).
		Str("clientType", string(j.GetClientType())).
		Str("name", name).
		Msg("Creating new playlist in Jellyfin")

	// Create a new playlist using Items API
	// In Jellyfin, playlists are created as special items
	createReq := j.client.PlaylistsAPI.CreatePlaylist(ctx).
		Name(name).
		UserId(j.config.UserID).
		// TODO: Use the correct media type if possible?
		MediaType(jellyfin.MEDIATYPE_UNKNOWN)

	// // Add description if provided
	// if description != "" {
	// 	createReq = createReq.Overview(description)
	// }

	response, _, err := createReq.Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("name", name).
			Msg("Failed to create playlist in Jellyfin")
		return nil, err
	}

	if response.Id == nil || *response.Id == "" {
		return nil, fmt.Errorf("created playlist has no ID")
	}

	// Convert response to our internal playlist model
	playlist := &models.MediaItem[*types.Playlist]{
		Data: &types.Playlist{
			ItemList: types.ItemList{
				Details: &types.MediaDetails{
					Title:       name,
					Description: description,
					// Use default artwork and fill in when items are added
				},
				ItemCount: 0, // New playlist has no items
				IsPublic:  true,
				// Set creation timestamp
				LastModified: time.Now(),
				ModifiedBy:   j.GetClientID(),
			},
		},
		Type: "playlist",
	}

	// Set client info
	playlist.SetClientInfo(j.GetClientID(), j.GetClientType(), *response.Id)

	log.Info().
		Str("playlistID", *response.Id).
		Str("name", name).
		Msg("Successfully created playlist in Jellyfin")

	return playlist, nil
}

// UpdatePlaylist updates an existing playlist in Jellyfin
func (j *JellyfinClient) UpdatePlaylist(ctx context.Context, playlistID string, name string, description string) (*models.MediaItem[*types.Playlist], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.GetClientID()).
		Str("clientType", string(j.GetClientType())).
		Str("playlistID", playlistID).
		Str("name", name).
		Msg("Updating playlist in Jellyfin")

	updateReq := jellyfin.UpdatePlaylistDto{
		Name: *jellyfin.NewNullableString(&name),
	}

	// First, get the current item to make sure it exists and is a playlist
	getReq := j.client.PlaylistsAPI.UpdatePlaylist(ctx, playlistID)
	getReq.UpdatePlaylistDto(updateReq)

	resp, err := getReq.Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Int("statusCode", resp.StatusCode).
			Msg("Failed to get playlist from Jellyfin")
		return nil, err
	}

	// Create our internal model with updated info
	playlist := models.NewMediaItem[*types.Playlist](&types.Playlist{
		ItemList: types.ItemList{
			Details: &types.MediaDetails{
				Title:       name,
				Description: description,
				// TODO: need to look into playlist artwork handling
				// May need to use the main items api if the playlists have cover at all.
				// Artwork:     j.getArtworkURLs(item),
			},
			LastModified: time.Now(),
			ModifiedBy:   j.GetClientID(),
			// Preserve existing items
			ItemCount: 0,
			IsPublic:  true,
		},
	})

	playlist.SyncClients.AddClient(j.GetClientID(), clienttypes.ClientTypeJellyfin, playlistID)

	// get items from the playlist
	itemReq := j.client.PlaylistsAPI.GetPlaylistItems(ctx, playlistID)
	playlistItems, resp, err := itemReq.Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Msg("Failed to get playlist items from Jellyfin")
		return nil, err
	}

	playlist.Data.ItemCount = int(*playlistItems.TotalRecordCount)

	playlist.SetClientInfo(j.GetClientID(), j.GetClientType(), playlistID)

	log.Info().
		Str("playlistID", playlistID).
		Str("name", name).
		Msg("Successfully updated playlist in Jellyfin")

	return playlist, nil
}

// DeletePlaylist deletes a playlist from Jellyfin
func (j *JellyfinClient) DeletePlaylist(ctx context.Context, playlistID string) error {
	// Get logger from context
	// log := logger.LoggerFromContext(ctx)
	//
	// log.Info().
	// 	Uint64("clientID", j.GetClientID()).
	// 	Str("clientType", string(j.GetClientType())).
	// 	Str("playlistID", playlistID).
	// 	Msg("Deleting playlist from Jellyfin")
	//
	// // First, verify the item exists and is a playlist
	// getReq := j.client.ItemsAPI.DeletePlaylist(ctx, playlistID)
	// item, resp, err := getReq.Execute()
	// if err != nil {
	// 	// If 404, consider it already deleted
	// 	if resp != nil && resp.StatusCode == 404 {
	// 		log.Warn().
	// 			Str("playlistID", playlistID).
	// 			Msg("Playlist not found in Jellyfin, considering it already deleted")
	// 		return nil
	// 	}
	//
	// 	log.Error().
	// 		Err(err).
	// 		Str("playlistID", playlistID).
	// 		Int("statusCode", resp.StatusCode).
	// 		Msg("Failed to get playlist from Jellyfin")
	// 	return err
	// }
	//
	// // Check that this is actually a playlist
	// if *item.Type != jellyfin.BASEITEMKIND_PLAYLIST {
	// 	log.Error().
	// 		Str("playlistID", playlistID).
	// 		Str("actualType", string(*item.Type)).
	// 		Msg("Item is not a playlist")
	// 	return fmt.Errorf("item %s is not a playlist", playlistID)
	// }
	//
	// // Delete the playlist using the ItemsAPI
	// deleteReq := j.client.PlaylistsAPI.RemoveItemFromPlaylist(ctx, playlistID)
	// deleteReq = deleteReq.EntryIds([]string{playlistID})
	// resp, err = deleteReq.Execute()
	// if err != nil {
	// 	log.Error().
	// 		Err(err).
	// 		Str("playlistID", playlistID).
	// 		Int("statusCode", resp.StatusCode).
	// 		Msg("Failed to delete playlist from Jellyfin")
	// 	return err
	// }
	//
	// log.Info().
	// 	Str("playlistID", playlistID).
	// 	Msg("Successfully deleted playlist from Jellyfin")
	// TODO: Implement playlist deletion for Jellyfin

	return nil
}

// AddItemToPlaylist adds an item to a playlist in Jellyfin
func (j *JellyfinClient) AddItemToPlaylist(ctx context.Context, playlistID string, itemID string) error {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.GetClientID()).
		Str("clientType", string(j.GetClientType())).
		Str("playlistID", playlistID).
		Str("itemID", itemID).
		Msg("Adding item to playlist in Jellyfin")

	// Use PlaylistsAPI to add item to playlist if available
	// If not, fall back to generic ItemsAPI
	request := j.client.PlaylistsAPI.AddItemToPlaylist(ctx, playlistID)
	// Add the item ID to the request
	// Note: We're adding a single item, but the API expects an array of IDs
	request = request.Ids([]string{itemID})

	// Execute the request
	resp, err := request.Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Str("itemID", itemID).
			Int("statusCode", resp.StatusCode).
			Msg("Failed to add item to playlist in Jellyfin")
		return err
	}

	log.Info().
		Str("playlistID", playlistID).
		Str("itemID", itemID).
		Msg("Successfully added item to playlist in Jellyfin")

	return nil
}

// RemoveItemFromPlaylist removes an item from a playlist in Jellyfin
func (j *JellyfinClient) RemoveItemFromPlaylist(ctx context.Context, playlistID string, itemID string) error {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.GetClientID()).
		Str("clientType", string(j.GetClientType())).
		Str("playlistID", playlistID).
		Str("itemID", itemID).
		Msg("Removing item from playlist in Jellyfin")

	// Use PlaylistsAPI to remove item from playlist
	request := j.client.PlaylistsAPI.RemoveItemFromPlaylist(ctx, playlistID)
	request.EntryIds([]string{itemID})

	// Try to find the EntryId or PlaylistItemId that corresponds to this item
	// First, get the playlist items to find the entry ID
	// playlistItems, err := j.GetPlaylistItems(ctx, playlistID, nil)
	// if err != nil {
	// 	log.Error().
	// 		Err(err).
	// 		Str("playlistID", playlistID).
	// 		Str("itemID", itemID).
	// 		Msg("Failed to get playlist items to find entry ID")
	// 	return err
	// }

	// Find the entry ID for this item
	var entryID string = ""
	// for _, orderItem := range playlistItems.ClientOrder {
	//
	// 	orderItem.ItemID
	//
	// 	client, exists := item.SyncClients.GetByClientID(j.GetClientID())
	// 	if exists && client.ItemID == itemID {
	// 		// This i the playlistID
	// 		// Found the item, try to get its entry ID
	// 		// TODO:
	// 		// Jellyfin may store the entry ID in the item's metadata
	// 		// if entryItemID, ok := item.Data; ok {
	// 		// entryID = entryItemID
	// 		// break
	// 	}
	// 	// If we can't find a specific entry ID, use the item ID
	// 	// entryID = id
	// 	break
	// }

	if entryID == "" {
		log.Warn().
			Str("playlistID", playlistID).
			Str("itemID", itemID).
			Msg("Item not found in playlist, cannot remove")
		return nil
	}

	// Add the entry ID to the request
	request = request.EntryIds([]string{entryID})

	// Execute the request
	resp, err := request.Execute()
	if err != nil {
		// If we get a 404, the item might already be removed or doesn't exist
		if resp != nil && resp.StatusCode == 404 {
			log.Warn().
				Str("playlistID", playlistID).
				Str("itemID", itemID).
				Msg("Item not found in playlist, may have been already removed")
			return nil
		}

		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Str("itemID", itemID).
			Int("statusCode", resp.StatusCode).
			Msg("Failed to remove item from playlist in Jellyfin")
		return err
	}

	log.Info().
		Str("playlistID", playlistID).
		Str("itemID", itemID).
		Msg("Successfully removed item from playlist in Jellyfin")

	return nil
}

// ReorderPlaylistItems reorders items in a playlist in Jellyfin
func (j *JellyfinClient) ReorderPlaylistItems(ctx context.Context, playlistID string, itemIDs []string) error {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.GetClientID()).
		Str("clientType", string(j.GetClientType())).
		Str("playlistID", playlistID).
		Int("itemCount", len(itemIDs)).
		Msg("Reordering playlist items in Jellyfin")

	if len(itemIDs) == 0 {
		log.Warn().
			Str("playlistID", playlistID).
			Msg("No items provided for reordering")
		return nil
	}

	return nil
}

// First, get the current playlist items to determine their entry IDs
// Jellyfin assigns special IDs to items within a playlist that are different from the media item IDs
// We need to map our media item IDs to these entry IDs
// currentItems, err := j.GetPlaylistItems(ctx, playlistID, nil)
// if err != nil {
// 	log.Error().
// 		Err(err).
// 		Str("playlistID", playlistID).
// 		Msg("Failed to get current playlist items for reordering")
// 	return err
// }
//
// // Create maps for item IDs, entry IDs, and current positions
// itemToEntryID := make(map[string]string)
// currentPositions := make(map[string]int, len(currentItems))
//
// // Build the mappings
// for i, item := range currentItems {
// 	client, exists := item.SyncClients.GetByClientID(j.GetClientID())
// 	if !exists {
// 		continue
// 	}
//
// 	// Store the current position
// 	currentPositions[client.ItemID] = i
//
// // Try to find entry ID in external IDs
// if entryItemID, ok := item.Data.Get; ok {
// 	itemToEntryID[clientID] = entryItemID
// } else {
// 	// Use the regular ID if no special entry ID is found
// 	itemToEntryID[clientID] = clientID
// }
// }

// For each item in the new order, move it to its new position
// We need to do this one by one since Jellyfin doesn't support reordering the entire playlist at once
// for newIndex, itemID := range itemIDs {
// If the item is already at the correct position, skip it
// if currentPos, exists := currentPositions[itemID]; exists && currentPos == newIndex {
// 	continue
// }
//
// Get the entry ID for this item
// entryID, exists := itemToEntryID[itemID]
// if !exists {
// 	log.Warn().
// 		Str("playlistID", playlistID).
// 		Str("itemID", itemID).
// 		Msg("Item not found in playlist, cannot reorder")
// 	continue
// }

// Move the item to its new position
// request := j.client.PlaylistsAPI.MoveItem(ctx, playlistID, entryID, int32(newIndex))

// Execute the request
// resp, err := request.Execute()
// if err != nil {
// 	log.Error().
// 		Err(err).
// 		Str("playlistID", playlistID).
// 		Str("itemID", itemID).
// 		Str("entryID", entryID).
// 		Int("newPosition", newIndex).
// 		Int("statusCode", resp.StatusCode).
// 		Msg("Failed to reorder item in playlist")
// 	return err
// }
//
// // Update the current positions map for subsequent operations
// currentPositions[itemID] = newIndex
// }

// log.Info().
// 	Str("playlistID", playlistID).
// 	Int("itemCount", len(itemIDs)).
// 	Msg("Successfully reordered playlist items in Jellyfin")
// return nil
