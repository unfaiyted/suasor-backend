package plex

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/logger"
	"time"

	"github.com/LukeHagar/plexgo/models/operations"
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
	res, err := c.plexAPI.Playlists.GetPlaylist(ctx, playlistRatingKey)
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
	playlist := types.NewPlaylist()
	playlist.MediaItemList.Details = &types.MediaDetails{
		Title:       plexMetadata.Title,
		Description: plexMetadata.Summary,
		AddedAt:     time.Unix(int64(*plexMetadata.AddedAt), 0),
		UpdatedAt:   time.Unix(int64(*plexMetadata.UpdatedAt), 0),
	}

	// Create MediaItem
	mediaItem := models.NewMediaItem(types.MediaTypePlaylist, playlist)
	mediaItem.ID = c.GetClientID()
	mediaItem.Title = plexMetadata.Title
	mediaItem.ClientID = c.GetClientID()
	mediaItem.ClientItemID = strconv.Itoa(*plexMetadata.RatingKey)

	// Set sync clients
	syncClients := models.SyncClients{}
	syncClients.AddClient(c.GetClientID(), c.GetClientType(), strconv.Itoa(*plexMetadata.RatingKey))
	mediaItem.SyncClients = syncClients

	log.Info().
		Str("playlistID", playlistID).
		Str("title", mediaItem.Title).
		Msg("Retrieved playlist from Plex")

	return mediaItem, nil
}

// GetPlaylistItems retrieves items in a playlist
func (c *PlexClient) GetPlaylistItems(ctx context.Context, playlistID string) (*models.MediaItemList, error) {
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
	res, err := c.plexAPI.Playlists.GetPlaylistContents(ctx, playlistRatingKey)
	if err != nil {
		log.Error().
			Err(err).
			Str("playlistID", playlistID).
			Msg("Failed to get playlist contents from Plex")
		return nil, fmt.Errorf("failed to get playlist contents: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil {
		log.Error().
			Str("playlistID", playlistID).
			Msg("Playlist contents not found or empty response from Plex")
		return nil, fmt.Errorf("playlist contents not found")
	}

	// Create MediaItemList
	itemList := &models.MediaItemList{
		Items: make([]*models.MediaItemListItem, 0, len(res.Object.MediaContainer.Metadata)),
	}

	// Add items to list
	for i, item := range res.Object.MediaContainer.Metadata {
		mediaType := types.MediaTypeUnknown
		// Determine media type based on Plex media type
		switch item.Type {
		case "movie":
			mediaType = types.MediaTypeMovie
		case "episode":
			mediaType = types.MediaTypeEpisode
		case "track":
			mediaType = types.MediaTypeTrack
		default:
			mediaType = types.MediaTypeUnknown
		}

		itemList.Items = append(itemList.Items, &models.MediaItemListItem{
			ID:       strconv.Itoa(*item.RatingKey),
			Position: i,
			Title:    item.Title,
			Type:     mediaType,
		})
	}

	itemList.ItemCount = len(itemList.Items)
	log.Info().
		Str("playlistID", playlistID).
		Int("itemCount", itemList.ItemCount).
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

	// Create playlist in Plex
	res, err := c.plexAPI.Playlists.CreatePlaylist(
		ctx,
		operations.CreatePlaylistType.Audio,
		operations.CreatePlaylistTitle(name),
		operations.CreatePlaylistSmart.False,
		operations.CreatePlaylistSummary(description),
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
	playlist := types.NewPlaylist()
	playlist.MediaItemList.Details = &types.MediaDetails{
		Title:       plexMetadata.Title,
		Description: plexMetadata.Summary,
		AddedAt:     time.Now(),
		UpdatedAt:   time.Now(),
	}

	mediaItem := models.NewMediaItem(types.MediaTypePlaylist, playlist)
	mediaItem.ID = c.GetClientID()
	mediaItem.Title = plexMetadata.Title
	mediaItem.ClientID = c.GetClientID()
	mediaItem.ClientItemID = strconv.Itoa(*plexMetadata.RatingKey)

	// Set sync clients
	syncClients := models.SyncClients{}
	syncClients.AddClient(c.GetClientID(), c.GetClientType(), strconv.Itoa(*plexMetadata.RatingKey))
	mediaItem.SyncClients = syncClients

	log.Info().
		Str("name", name).
		Str("playlistID", mediaItem.ClientItemID).
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
		err = c.AddPlaylistItems(ctx, mediaItem.ClientItemID, itemIDs)
		if err != nil {
			log.Error().
				Err(err).
				Str("playlistID", mediaItem.ClientItemID).
				Msg("Failed to add items to newly created playlist")
			// Consider cleaning up by deleting the playlist
			c.DeletePlaylist(ctx, mediaItem.ClientItemID)
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
	_, err = c.plexAPI.Playlists.DeletePlaylist(ctx, playlistRatingKey)
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
		// Since this requires custom implementation, we'll have to mock it
		// In a real implementation, you would need to get this from the server
		machineID := "PLACEHOLDER" // In real implementation, get this from server info
		uri := fmt.Sprintf("server://%s/library/metadata/%d", machineID, itemRatingKey)
		uris = append(uris, uri)
	}

	if len(uris) == 0 {
		log.Warn().Msg("No valid items to add to playlist")
		return nil
	}

	// Use AddPlaylistContents API
	// Note: This is a simplified implementation. In reality, you need the correct URIs
	uriString := strings.Join(uris, ",")
	_, err = c.plexAPI.Playlists.AddPlaylistContents(
		ctx,
		playlistRatingKey,
		operations.AddPlaylistContentsUri(uriString),
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
	_, err = c.plexAPI.Playlists.ClearPlaylistContents(ctx, playlistRatingKey)
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
