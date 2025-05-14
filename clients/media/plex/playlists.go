package plex

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/logger"

	"github.com/unfaiyted/plexgo/models/operations"
)

// GetPlaylists retrieves playlists from Plex
func (c *PlexClient) GetPlaylists(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Playlist], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Msg("Retrieving playlists from Plex server")

	log.Debug().Msg("Making API request to Plex server for playlists")
	res, err := c.plexAPI.Playlists.GetPlaylists(ctx, operations.PlaylistTypeAudio.ToPointer(), operations.QueryParamSmartOne.ToPointer())
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Msg("Failed to get playlists from Plex")
		return nil, fmt.Errorf("failed to get playlists: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil {
		log.Info().
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Msg("No playlists found in Plex")
		return nil, nil
	}

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Int("totalItems", len(res.Object.MediaContainer.Metadata)).
		Msg("Successfully retrieved playlists from Plex")

	playlists, err := GetMediaItemListFromPlaylist[*types.Playlist](ctx, c, res.Object.MediaContainer.Metadata)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Int("playlistsReturned", len(playlists)).
		Msg("Completed GetPlaylists request")

	return playlists, nil
}

// GetPlaylist retrieves a single playlist by ID
func (c *PlexClient) GetPlaylist(ctx context.Context, playlistID string) (*models.MediaItem[*types.Playlist], error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("playlistID", playlistID).
		Msg("Retrieving playlist from Plex server")

	// Convert playlistID to integer for Plex
	playlistRatingKey, err := strconv.Atoi(playlistID)
	if err != nil {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Msg("Failed to convert playlist ID to integer")
		return nil, fmt.Errorf("invalid playlist ID: %w", err)
	}

	// Get playlist metadata from Plex
	res, err := c.plexAPI.Playlists.GetPlaylist(ctx, float64(playlistRatingKey))
	if err != nil {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Msg("Failed to get playlist from Plex")
		return nil, fmt.Errorf("failed to get playlist: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil || len(res.Object.MediaContainer.Metadata) == 0 {
		log.Error().
			Str("playlistID", playlistID).
			Msg("Playlist not found or empty response from Plex")
		return nil, fmt.Errorf("playlist not found")
	}

	// Convert single Plex playlist to MediaItem
	plexMetadata := res.Object.MediaContainer.Metadata[0]
	playlist, err := GetItemFromPlaylistMetadata[*types.Playlist](ctx, c, &plexMetadata)
	if err != nil {
		return nil, err
	}

	// Create MediaItem
	mediaItem := models.NewMediaItem(playlist)
	mediaItem.ID = c.GetClientID()
	mediaItem.Title = playlist.GetTitle()

	mediaItem.SetClientInfo(c.GetClientID(), c.GetClientType(), *plexMetadata.RatingKey)

	log.Info().
		Str("playlistID", playlistID).
		Str("title", mediaItem.Title).
		Msg("Retrieved playlist from Plex")

	return mediaItem, nil
}

// GetPlaylistItems retrieves items in a playlist
func (c *PlexClient) GetPlaylistItems(ctx context.Context, playlistID string) (*models.MediaItemList[*types.Playlist], error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("playlistID", playlistID).
		Msg("Retrieving playlist items from Plex server")

	// Convert playlistID to integer for Plex
	playlistRatingKey, err := strconv.Atoi(playlistID)
	if err != nil {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Msg("Failed to convert playlist ID to integer")
		return nil, fmt.Errorf("invalid playlist ID: %w", err)
	}

	// Get playlist contents from Plex
	itemList, err := c.GetAllPlaylistContentsTypes(ctx, float64(playlistRatingKey))
	if err != nil {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Msg("Failed to get playlist contents from Plex")
		return nil, fmt.Errorf("failed to get playlist contents: %w", err)
	}

	log.Info().
		Str("playlistID", playlistID).
		Int("itemCount", itemList.Len()).
		Msg("Retrieved playlist items from Plex")

	return itemList, nil
}

// CreatePlaylist creates a new playlist
func (c *PlexClient) CreatePlaylist(ctx context.Context, name string, description string) (*models.MediaItem[*types.Playlist], error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("name", name).
		Msg("Creating playlist in Plex server")

	createRequest := &operations.CreatePlaylistRequest{
		Title: name,
		Type:  operations.CreatePlaylistQueryParamTypeAudio,
		Smart: operations.SmartZero,
	}

	// Create playlist in Plex
	res, err := c.plexAPI.Playlists.CreatePlaylist(
		ctx,
		*createRequest,
	)
	if err != nil {
		log.Error().
			Err(err).
			Str("name", name).
			Msg("Failed to create playlist in Plex")
		return nil, fmt.Errorf("failed to create playlist: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil || len(res.Object.MediaContainer.Metadata) == 0 {
		log.Error().
			Str("name", name).
			Msg("Empty response when creating playlist in Plex")
		return nil, fmt.Errorf("empty response when creating playlist")
	}

	// Create MediaItem from created playlist
	plexMetadata := res.Object.MediaContainer.Metadata[0]

	playlist, err := GetItemFromPlaylistCreate[*types.Playlist](ctx, c, &plexMetadata)

	mediaItem := models.NewMediaItem(playlist)
	mediaItem.ID = c.GetClientID()
	mediaItem.Title = *plexMetadata.Title
	mediaItem.SetClientInfo(c.GetClientID(), c.GetClientType(), *plexMetadata.RatingKey)

	log.Info().
		Str("name", name).
		Str("playlistID", mediaItem.SyncClients.GetClientItemID(c.GetClientID())).
		Msg("Created playlist in Plex")

	return mediaItem, nil
}

// CreatePlaylistWithItems creates a new playlist with items
func (c *PlexClient) CreatePlaylistWithItems(ctx context.Context, name string, description string, itemIDs []string) (*models.MediaItem[*types.Playlist], error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("name", name).
		Int("itemCount", len(itemIDs)).
		Msg("Creating playlist with items in Plex server")

	// First create the playlist
	mediaItem, err := c.CreatePlaylist(ctx, name, description)
	if err != nil {
		return nil, err
	}

	// Then add items if there are any
	if len(itemIDs) > 0 {
		clientItemID := mediaItem.SyncClients.GetClientItemID(c.GetClientID())
		err = c.AddPlaylistItems(ctx, clientItemID, itemIDs)
		if err != nil {
			log.Error().
				Err(err).
				Str("playlistID", clientItemID).
				Msg("Failed to add items to newly created playlist")
			// Consider cleaning up by deleting the playlist
			c.DeletePlaylist(ctx, clientItemID)
			return nil, fmt.Errorf("failed to add items to newly created playlist: %w", err)
		}
	}

	return mediaItem, nil
}

// UpdatePlaylist updates a playlist's metadata
func (c *PlexClient) UpdatePlaylist(ctx context.Context, playlistID string, name string, description string) (*models.MediaItem[*types.Playlist], error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("playlistID", playlistID).
		Str("name", name).
		Msg("Updating playlist in Plex server")

	// Plex doesn't have a direct API for updating playlist metadata
	// We need to use the PUT method with the /library/metadata/{id} endpoint
	// This is a custom implementation as the plexAPI doesn't provide this directly

	// First get the current playlist to ensure it exists
	_, err := c.GetPlaylist(ctx, playlistID)
	if err != nil {
		return nil, fmt.Errorf("failed to find playlist to update: %w", err)
	}

	// Plex doesn't support direct playlist updating through the API wrapper
	// In a real implementation, this would need to use a custom HTTP request to Plex
	log.Warn().
		Str("playlistID", playlistID).
		Msg("Plex API does not support direct playlist metadata updates, only title and summary changes through custom API calls")

	// For now, we'll return the playlist as is
	// A complete solution would need to implement a custom HTTP request
	return c.GetPlaylist(ctx, playlistID)
}

// DeletePlaylist deletes a playlist
func (c *PlexClient) DeletePlaylist(ctx context.Context, playlistID string) error {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("playlistID", playlistID).
		Msg("Deleting playlist from Plex server")

	// Convert playlistID to integer for Plex
	playlistRatingKey, err := strconv.Atoi(playlistID)
	if err != nil {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Msg("Failed to convert playlist ID to integer")
		return fmt.Errorf("invalid playlist ID: %w", err)
	}

	// Delete playlist in Plex
	_, err = c.plexAPI.Playlists.DeletePlaylist(ctx, float64(playlistRatingKey))
	if err != nil {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Msg("Failed to delete playlist from Plex")
		return fmt.Errorf("failed to delete playlist: %w", err)
	}

	log.Info().
		Str("playlistID", playlistID).
		Msg("Playlist deleted from Plex")

	return nil
}

// AddPlaylistItem adds an item to a playlist
func (c *PlexClient) AddPlaylistItem(ctx context.Context, playlistID string, itemID string) error {
	return c.AddPlaylistItems(ctx, playlistID, []string{itemID})
}

// AddPlaylistItems adds multiple items to a playlist
func (c *PlexClient) AddPlaylistItems(ctx context.Context, playlistID string, itemIDs []string) error {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("playlistID", playlistID).
		Int("itemCount", len(itemIDs)).
		Msg("Adding items to playlist in Plex server")

	if len(itemIDs) == 0 {
		log.Warn().Msg("No items to add to playlist")
		return nil
	}

	// Convert playlistID to integer for Plex
	playlistRatingKey, err := strconv.Atoi(playlistID)
	if err != nil {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Msg("Failed to convert playlist ID to integer")
		return fmt.Errorf("invalid playlist ID: %w", err)
	}

	// Format item IDs as required by Plex API
	// Plex expects URIs in the format "server://machineIdentifier/library/metadata/itemID"
	var uris []string
	for _, id := range itemIDs {
		itemRatingKey, err := strconv.Atoi(id)
		if err != nil {
			log.Warn().
				Err(err).
				Str("itemID", id).
				Msg("Invalid item ID, skipping")
			continue
		}

		// We need the server's machine identifier to construct the URI
		serverInfo, err := c.plexAPI.Server.GetServerIdentity(ctx)
		if err != nil {
			log.Error().
				Err(err).
				Str("itemID", id).
				Msg("Failed to get server info to construct URI")
			continue
		}
		machineID := serverInfo.Object.MediaContainer.MachineIdentifier
		uri := fmt.Sprintf("server://%s/library/metadata/%d", machineID, itemRatingKey)
		uris = append(uris, uri)
	}

	if len(uris) == 0 {
		log.Warn().Msg("No valid items to add to playlist")
		return nil
	}

	// Use AddPlaylistContents API
	uriString := strings.Join(uris, ",")
	_, err = c.plexAPI.Playlists.AddPlaylistContents(
		ctx,
		float64(playlistRatingKey),
		uriString,
		nil,
	)
	if err != nil {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Msg("Failed to add items to playlist")
		return fmt.Errorf("failed to add items to playlist: %w", err)
	}

	log.Info().
		Str("playlistID", playlistID).
		Int("itemCount", len(itemIDs)).
		Msg("Items added to playlist")

	return nil
}

// RemovePlaylistItem removes an item from a playlist
func (c *PlexClient) RemovePlaylistItem(ctx context.Context, playlistID string, itemID string) error {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("playlistID", playlistID).
		Str("itemID", itemID).
		Msg("Removing item from playlist in Plex server")

	// Plex doesn't have a direct API for removing a specific item from a playlist
	// This would require getting all items, filtering out the one to remove, and reconstructing the playlist
	// For now, we'll log a warning
	log.Warn().
		Str("playlistID", playlistID).
		Str("itemID", itemID).
		Msg("Plex API does not support direct individual item removal from playlists, this requires custom implementation")

	return fmt.Errorf("operation not supported by Plex API")
}

// RemovePlaylistItems removes multiple items from a playlist
func (c *PlexClient) RemovePlaylistItems(ctx context.Context, playlistID string, itemIDs []string) error {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("playlistID", playlistID).
		Int("itemCount", len(itemIDs)).
		Msg("Removing items from playlist in Plex server")

	// Similar to RemovePlaylistItem, Plex doesn't have a direct API for this
	log.Warn().
		Str("playlistID", playlistID).
		Msg("Plex API does not support direct removal of multiple items from playlists, this requires custom implementation")

	return fmt.Errorf("operation not supported by Plex API")
}

// RemoveAllPlaylistItems removes all items from a playlist
func (c *PlexClient) RemoveAllPlaylistItems(ctx context.Context, playlistID string) error {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("playlistID", playlistID).
		Msg("Removing all items from playlist in Plex server")

	// Convert playlistID to integer for Plex
	playlistRatingKey, err := strconv.Atoi(playlistID)
	if err != nil {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Msg("Failed to convert playlist ID to integer")
		return fmt.Errorf("invalid playlist ID: %w", err)
	}

	// Clear playlist contents
	_, err = c.plexAPI.Playlists.ClearPlaylistContents(ctx, float64(playlistRatingKey))
	if err != nil {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Msg("Failed to clear playlist contents")
		return fmt.Errorf("failed to clear playlist contents: %w", err)
	}

	log.Info().
		Str("playlistID", playlistID).
		Msg("All items removed from playlist")

	return nil
}

// ReorderPlaylistItems reorders items in a playlist
func (c *PlexClient) ReorderPlaylistItems(ctx context.Context, playlistID string, itemIDs []string) error {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("playlistID", playlistID).
		Int("itemCount", len(itemIDs)).
		Msg("Reordering items in playlist in Plex server")

	// Plex doesn't have a direct API for reordering playlist items
	log.Warn().
		Str("playlistID", playlistID).
		Msg("Plex API does not support direct reordering of playlist items, this requires custom implementation")

	return fmt.Errorf("operation not supported by Plex API")
}

// SearchPlaylists searches playlists
func (c *PlexClient) SearchPlaylists(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Playlist], error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Msg("Searching playlists in Plex server")

	// Plex doesn't have a specific API for searching playlists
	// We'll reuse GetPlaylists and filter the results
	playlists, err := c.GetPlaylists(ctx, options)
	if err != nil {
		return nil, err
	}

	// If no search query, return all playlists
	if options == nil || options.Query == "" {
		return playlists, nil
	}

	// Filter playlists by search query
	query := strings.ToLower(options.Query)
	var filteredPlaylists []*models.MediaItem[*types.Playlist]
	for _, playlist := range playlists {
		if strings.Contains(strings.ToLower(playlist.Title), query) {
			filteredPlaylists = append(filteredPlaylists, playlist)
		}
	}

	log.Info().
		Str("query", options.Query).
		Int("resultCount", len(filteredPlaylists)).
		Msg("Filtered playlists by search query")

	return filteredPlaylists, nil
}

// SearchPlaylistItems searches items in a playlist
func (c *PlexClient) SearchPlaylistItems(ctx context.Context, playlistID string, options *types.QueryOptions) ([]*models.MediaItem[*types.Playlist], error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("playlistID", playlistID).
		Msg("Searching items in playlist in Plex server")

	// This operation isn't directly supported by Plex API
	// We would need to get all items and filter them
	log.Warn().
		Str("playlistID", playlistID).
		Msg("Plex API does not support searching within playlist items, this requires custom implementation")

	return nil, fmt.Errorf("operation not supported by Plex API")
}

// SupportsPlaylists returns whether the client supports playlists
func (c *PlexClient) SupportsPlaylists() bool {
	return true
}

func (c *PlexClient) GetAllPlaylistContentsTypes(ctx context.Context, playlistID float64) (*models.MediaItemList[*types.Playlist], error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Float64("playlistID", playlistID).
		Msg("Retrieving all playlist contents types from Plex server")

	strPlaylistID := fmt.Sprintf("%d", playlistID)

	playlist, err := c.GetPlaylist(ctx, strPlaylistID)
	if err != nil {
		return nil, err
	}
	itemList := models.NewMediaItemList[*types.Playlist](playlist, 0, c.GetClientID())

	playlistContentsTypes := make([]operations.GetPlaylistContentsQueryParamType, 0, 10)
	playlistContentsTypes = append(playlistContentsTypes,
		operations.GetPlaylistContentsQueryParamTypeMovie,
		operations.GetPlaylistContentsQueryParamTypeTvShow,
		operations.GetPlaylistContentsQueryParamTypeSeason,
		operations.GetPlaylistContentsQueryParamTypeEpisode,
		operations.GetPlaylistContentsQueryParamTypeAudio,
		operations.GetPlaylistContentsQueryParamTypeAlbum,
		operations.GetPlaylistContentsQueryParamTypeTrack,
	)

	for _, playlistType := range playlistContentsTypes {
		log.Debug().
			Float64("playlistID", playlistID).
			Str("playlistType", string(playlistType)).
			Msg("Retrieving playlist contents of type")
		res, err := c.plexAPI.Playlists.GetPlaylistContents(ctx, playlistID, playlistType)
		if err != nil {
			return nil, err
		}

		// Add the playlist items to the list
		if res.Object.MediaContainer != nil && res.Object.MediaContainer.Metadata != nil {
			for _, item := range res.Object.MediaContainer.Metadata {
				log.Debug().
					Float64("playlistID", playlistID).
					Str("playlistType", string(playlistType)).
					Msg("Adding playlist item to list")

				switch playlistType {
				case operations.GetPlaylistContentsQueryParamTypeMovie:
					rawItem, err := GetItemFromPlaylistContents[*types.Movie](ctx, c, &item)
					if err != nil {
						return nil, err
					}
					mediaItem, err := GetMediaItem[*types.Movie](ctx, c, rawItem, *item.RatingKey)
					if err != nil {
						return nil, err
					}
					mediaItem.SetClientInfo(c.GetClientID(), c.GetClientType(), *item.RatingKey)
					itemList.Items.AddMovie(mediaItem)
				case operations.GetPlaylistContentsQueryParamTypeTvShow:
					rawItem, err := GetItemFromPlaylistContents[*types.Series](ctx, c, &item)
					if err != nil {
						return nil, err
					}
					mediaItem, err := GetMediaItem[*types.Series](ctx, c, rawItem, *item.RatingKey)
					if err != nil {
						return nil, err
					}
					mediaItem.SetClientInfo(c.GetClientID(), c.GetClientType(), *item.RatingKey)
					itemList.Items.AddSeries(mediaItem)
				case operations.GetPlaylistContentsQueryParamTypeSeason:
					rawItem, err := GetItemFromPlaylistContents[*types.Season](ctx, c, &item)
					if err != nil {
						return nil, err
					}
					mediaItem, err := GetMediaItem[*types.Season](ctx, c, rawItem, *item.RatingKey)
					if err != nil {
						return nil, err
					}
					mediaItem.SetClientInfo(c.GetClientID(), c.GetClientType(), *item.RatingKey)
					itemList.Items.AddSeason(mediaItem)
				case operations.GetPlaylistContentsQueryParamTypeEpisode:
					rawItem, err := GetItemFromPlaylistContents[*types.Episode](ctx, c, &item)
					if err != nil {
						return nil, err
					}
					mediaItem, err := GetMediaItem[*types.Episode](ctx, c, rawItem, *item.RatingKey)
					if err != nil {
						return nil, err
					}
					mediaItem.SetClientInfo(c.GetClientID(), c.GetClientType(), *item.RatingKey)
					itemList.Items.AddEpisode(mediaItem)
				case operations.GetPlaylistContentsQueryParamTypeAudio:
					rawItem, err := GetItemFromPlaylistContents[*types.Track](ctx, c, &item)
					if err != nil {
						return nil, err
					}
					mediaItem, err := GetMediaItem[*types.Track](ctx, c, rawItem, *item.RatingKey)
					if err != nil {
						return nil, err
					}
					mediaItem.SetClientInfo(c.GetClientID(), c.GetClientType(), *item.RatingKey)
					itemList.Items.AddTrack(mediaItem)
				case operations.GetPlaylistContentsQueryParamTypeAlbum:
					rawItem, err := GetItemFromPlaylistContents[*types.Album](ctx, c, &item)
					if err != nil {
						return nil, err
					}
					mediaItem, err := GetMediaItem[*types.Album](ctx, c, rawItem, *item.RatingKey)
					if err != nil {
						return nil, err
					}
					mediaItem.SetClientInfo(c.GetClientID(), c.GetClientType(), *item.RatingKey)
					itemList.Items.AddAlbum(mediaItem)
				case operations.GetPlaylistContentsQueryParamTypeTrack:
					rawItem, err := GetItemFromPlaylistContents[*types.Track](ctx, c, &item)
					if err != nil {
						return nil, err
					}
					mediaItem, err := GetMediaItem[*types.Track](ctx, c, rawItem, *item.RatingKey)
					if err != nil {
						return nil, err
					}
					mediaItem.SetClientInfo(c.GetClientID(), c.GetClientType(), *item.RatingKey)
					itemList.Items.AddTrack(mediaItem)
				default:
					log.Error().
						Float64("playlistID", playlistID).
						Str("playlistType", string(playlistType)).
						Msg("Unknown playlist item type")
					return nil, fmt.Errorf("unknown playlist item type: %s", playlistType)
				}

			}
		}
	}
	return itemList, nil
}
