package services

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"time"

	"suasor/client"
	"suasor/client/media"
	"suasor/client/media/providers"
	mediatypes "suasor/client/media/types"
	"suasor/client/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/utils"
)

// ClientMediaPlaylistService defines operations for interacting with playlist clients
type ClientMediaPlaylistService[T types.ClientConfig] interface {
	// Basic playlist operations
	GetPlaylistByID(ctx context.Context, userID uint64, clientID uint64, playlistID string) (*models.MediaItem[*mediatypes.Playlist], error)
	GetPlaylists(ctx context.Context, userID uint64, count int) ([]models.MediaItem[*mediatypes.Playlist], error)
	CreatePlaylist(ctx context.Context, userID uint64, clientID uint64, name string, description string) (*models.MediaItem[*mediatypes.Playlist], error)
	UpdatePlaylist(ctx context.Context, userID uint64, clientID uint64, playlistID string, name string, description string) (*models.MediaItem[*mediatypes.Playlist], error)
	DeletePlaylist(ctx context.Context, userID uint64, clientID uint64, playlistID string) error

	// Playlist item operations
	GetPlaylistItems(ctx context.Context, userID uint64, clientID uint64, playlistID string) ([]models.MediaItem[mediatypes.MediaData], error)
	AddItemToPlaylist(ctx context.Context, userID uint64, clientID uint64, playlistID string, itemID string) error
	RemoveItemFromPlaylist(ctx context.Context, userID uint64, clientID uint64, playlistID string, itemID string) error
	ReorderPlaylistItems(ctx context.Context, userID uint64, clientID uint64, playlistID string, itemIDs []string) error

	// Search and sync operations
	SearchPlaylists(ctx context.Context, userID uint64, query string) ([]models.MediaItem[*mediatypes.Playlist], error)
	SyncPlaylist(ctx context.Context, userID uint64, clientID uint64, playlistID string) error
}

type mediaPlaylistService[T types.ClientMediaConfig] struct {
	repo    repository.ClientRepository[T]
	factory *client.ClientFactoryService
}

// NewClientMediaPlaylistService creates a new media playlist service
func NewClientMediaPlaylistService[T types.ClientMediaConfig](
	repo repository.ClientRepository[T],
	factory *client.ClientFactoryService,
) ClientMediaPlaylistService[T] {
	return &mediaPlaylistService[T]{
		repo:    repo,
		factory: factory,
	}
}

// getPlaylistClients gets all playlist clients for a user
func (s *mediaPlaylistService[T]) getPlaylistClients(ctx context.Context, userID uint64) ([]media.ClientMedia, error) {
	repo := s.repo
	// Get all media clients for the user
	clients, err := repo.GetByCategory(ctx, types.ClientCategoryMedia, userID)
	if err != nil {
		return nil, err
	}

	var playlistClients []media.ClientMedia

	// Filter and instantiate clients that support playlists
	for _, clientConfig := range clients {
		if clientConfig.Config.Data.SupportsPlaylists() {
			clientId := clientConfig.GetID()
			client, err := s.factory.GetClient(ctx, clientId, clientConfig.Config.Data)
			if err != nil {
				// Log error but continue with other clients
				continue
			}
			playlistClients = append(playlistClients, client.(media.ClientMedia))
		}
	}

	return playlistClients, nil
}

// getSpecificPlaylistClient gets a specific playlist client
func (s *mediaPlaylistService[T]) getSpecificPlaylistClient(ctx context.Context, userID, clientID uint64) (media.ClientMedia, error) {
	log := utils.LoggerFromContext(ctx)

	clientConfig, err := (s.repo).GetByID(ctx, clientID)
	if err != nil {
		return nil, err
	}
	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientType", clientConfig.Config.Data.GetType().String()).
		Msg("Retrieved client config")

	if !clientConfig.Config.Data.SupportsPlaylists() {
		log.Warn().
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("clientType", clientConfig.Config.Data.GetType().String()).
			Msg("Client does not support playlists")
		return nil, ErrUnsupportedFeature
	}

	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientType", clientConfig.Config.Data.GetType().String()).
		Msg("Client supports playlists")

	client, err := s.factory.GetClient(ctx, clientID, clientConfig.Config.Data)
	if err != nil {
		return nil, err
	}
	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientType", clientConfig.Config.Data.GetType().String()).
		Msg("Retrieved client")
	return client.(media.ClientMedia), nil
}

func (s *mediaPlaylistService[T]) GetPlaylistByID(ctx context.Context, userID uint64, clientID uint64, playlistID string) (*models.MediaItem[*mediatypes.Playlist], error) {
	client, err := s.getSpecificPlaylistClient(ctx, userID, clientID)
	log := utils.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Msg("Retrieved client")
	if err != nil {
		return nil, err
	}
	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Msg("Retrieving playlist")

	playlistProvider, ok := client.(providers.PlaylistProvider)
	if !ok {
		log.Warn().
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("playlistID", playlistID).
			Msg("Client does not support playlists")
		return nil, ErrUnsupportedFeature
	}

	// Check if the client supports getting playlist by ID
	if !playlistProvider.SupportsPlaylists() {
		log.Warn().
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("playlistID", playlistID).
			Msg("Client does not support playlists")
		return nil, ErrUnsupportedFeature
	}

	// Get all playlists and find by ID
	options := &mediatypes.QueryOptions{
		ExternalSourceID: playlistID,
	}

	playlists, err := playlistProvider.GetPlaylists(ctx, options)
	if err != nil {
		return nil, err
	}

	// Check if we found any playlists
	if len(playlists) == 0 {
		return nil, errors.New("playlist not found")
	}

	// Return the first matching playlist
	return &playlists[0], nil
}

func (s *mediaPlaylistService[T]) GetPlaylists(ctx context.Context, userID uint64, count int) ([]models.MediaItem[*mediatypes.Playlist], error) {
	clients, err := s.getPlaylistClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allPlaylists []models.MediaItem[*mediatypes.Playlist]

	for _, client := range clients {
		playlistProvider, ok := client.(providers.PlaylistProvider)
		if !ok || !playlistProvider.SupportsPlaylists() {
			continue
		}

		options := &mediatypes.QueryOptions{
			Limit: count,
		}

		playlists, err := playlistProvider.GetPlaylists(ctx, options)
		if err != nil {
			continue
		}

		allPlaylists = append(allPlaylists, playlists...)
	}

	// Sort by added date
	sort.Slice(allPlaylists, func(i, j int) bool {
		return allPlaylists[i].Data.GetDetails().AddedAt.After(allPlaylists[j].Data.GetDetails().AddedAt)
	})

	// Limit to requested count if specified
	if count > 0 && len(allPlaylists) > count {
		allPlaylists = allPlaylists[:count]
	}

	return allPlaylists, nil
}

func (s *mediaPlaylistService[T]) CreatePlaylist(ctx context.Context, userID uint64, clientID uint64, name string, description string) (*models.MediaItem[*mediatypes.Playlist], error) {
	log := utils.LoggerFromContext(ctx)
	client, err := s.getSpecificPlaylistClient(ctx, userID, clientID)
	if err != nil {
		return nil, err
	}

	playlistProvider, ok := client.(providers.PlaylistProvider)
	if !ok {
		log.Warn().
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("name", name).
			Msg("Client does not support playlists")
		return nil, ErrUnsupportedFeature
	}

	// Create the playlist
	playlist, err := playlistProvider.CreatePlaylist(ctx, name, description)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("name", name).
			Msg("Failed to create playlist")
		return nil, err
	}

	// Update modification timestamp and client ID
	now := time.Now()
	playlist.Data.ItemList.LastModified = now
	playlist.Data.ItemList.ModifiedBy = clientID

	// Get the client's ID for this playlist - ModifiedBy is just a uint64 client ID
	clientItemID, found := playlist.GetClientItemID(clientID) // Default to Plex as a placeholder
	playlistIDStr := "unknown"
	if found {
		playlistIDStr = clientItemID
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("playlistID", playlistIDStr).
		Str("name", name).
		Msg("Created playlist")

	return playlist, nil
}

func (s *mediaPlaylistService[T]) UpdatePlaylist(ctx context.Context, userID uint64, clientID uint64, playlistID string, name string, description string) (*models.MediaItem[*mediatypes.Playlist], error) {
	log := utils.LoggerFromContext(ctx)
	client, err := s.getSpecificPlaylistClient(ctx, userID, clientID)
	if err != nil {
		return nil, err
	}

	playlistProvider, ok := client.(providers.PlaylistProvider)
	if !ok {
		log.Warn().
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("playlistID", playlistID).
			Msg("Client does not support playlists")
		return nil, ErrUnsupportedFeature
	}

	// Update the playlist
	playlist, err := playlistProvider.UpdatePlaylist(ctx, playlistID, name, description)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("playlistID", playlistID).
			Msg("Failed to update playlist")
		return nil, err
	}

	// Update modification timestamp and client ID
	now := time.Now()
	playlist.Data.ItemList.LastModified = now
	playlist.Data.ItemList.ModifiedBy = clientID

	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Str("name", name).
		Msg("Updated playlist")

	return playlist, nil
}

func (s *mediaPlaylistService[T]) DeletePlaylist(ctx context.Context, userID uint64, clientID uint64, playlistID string) error {
	log := utils.LoggerFromContext(ctx)
	client, err := s.getSpecificPlaylistClient(ctx, userID, clientID)
	if err != nil {
		return err
	}

	playlistProvider, ok := client.(providers.PlaylistProvider)
	if !ok {
		log.Warn().
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("playlistID", playlistID).
			Msg("Client does not support playlists")
		return ErrUnsupportedFeature
	}

	// Delete the playlist
	err = playlistProvider.DeletePlaylist(ctx, playlistID)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("playlistID", playlistID).
			Msg("Failed to delete playlist")
		return err
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Msg("Deleted playlist")

	return nil
}

func (s *mediaPlaylistService[T]) AddItemToPlaylist(ctx context.Context, userID uint64, clientID uint64, playlistID string, itemID string) error {
	log := utils.LoggerFromContext(ctx)
	client, err := s.getSpecificPlaylistClient(ctx, userID, clientID)
	if err != nil {
		return err
	}

	playlistProvider, ok := client.(providers.PlaylistProvider)
	if !ok {
		log.Warn().
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("playlistID", playlistID).
			Msg("Client does not support playlists")
		return ErrUnsupportedFeature
	}

	// Add item to playlist
	err = playlistProvider.AddItemToPlaylist(ctx, playlistID, itemID)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("playlistID", playlistID).
			Str("itemID", itemID).
			Msg("Failed to add item to playlist")
		return err
	}

	// Get the updated playlist to update its metadata
	options := &mediatypes.QueryOptions{
		ExternalSourceID: playlistID,
	}
	playlists, err := playlistProvider.GetPlaylists(ctx, options)
	if err == nil && len(playlists) > 0 {
		// Record the change in the playlist metadata
		now := time.Now()
		playlist := &playlists[0]
		playlist.Data.LastModified = now
		playlist.Data.ModifiedBy = clientID

		// Add to change history if Items array is used
		if len(playlist.Data.Items) > 0 {
			// Try to convert string itemID to uint64 for comparison
			var numericItemID uint64
			if id, err := strconv.ParseUint(itemID, 10, 64); err == nil {
				numericItemID = id
			}

			// Find the item and update its change history
			for i, item := range playlist.Data.Items {
				if item.ItemID == numericItemID {
					// Update existing item
					playlist.Data.Items[i].LastChanged = now
					playlist.Data.Items[i].ChangeHistory = append(playlist.Data.Items[i].ChangeHistory,
						mediatypes.ChangeRecord{
							ClientID:   clientID,
							ItemID:     itemID,
							ChangeType: "add",
							Timestamp:  now,
						})
					break
				}
			}
		}
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Str("itemID", itemID).
		Msg("Added item to playlist")

	return nil
}

func (s *mediaPlaylistService[T]) RemoveItemFromPlaylist(ctx context.Context, userID uint64, clientID uint64, playlistID string, itemID string) error {
	log := utils.LoggerFromContext(ctx)
	client, err := s.getSpecificPlaylistClient(ctx, userID, clientID)
	if err != nil {
		return err
	}

	playlistProvider, ok := client.(providers.PlaylistProvider)
	if !ok {
		log.Warn().
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("playlistID", playlistID).
			Msg("Client does not support playlists")
		return ErrUnsupportedFeature
	}

	// Remove item from playlist
	err = playlistProvider.RemoveItemFromPlaylist(ctx, playlistID, itemID)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("playlistID", playlistID).
			Str("itemID", itemID).
			Msg("Failed to remove item from playlist")
		return err
	}

	// The change is already recorded in the client, but we should update our change history
	// This would require additional code to maintain a history in our database

	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Str("itemID", itemID).
		Msg("Removed item from playlist")

	return nil
}

// GetPlaylistItems gets all items in a playlist
func (s *mediaPlaylistService[T]) GetPlaylistItems(ctx context.Context, userID uint64, clientID uint64, playlistID string) ([]models.MediaItem[mediatypes.MediaData], error) {
	log := utils.LoggerFromContext(ctx)
	client, err := s.getSpecificPlaylistClient(ctx, userID, clientID)
	if err != nil {
		return nil, err
	}

	playlistProvider, ok := client.(providers.PlaylistProvider)
	if !ok {
		log.Warn().
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("playlistID", playlistID).
			Msg("Client does not support playlists")
		return nil, ErrUnsupportedFeature
	}

	// Get all items in the playlist
	items, err := playlistProvider.GetPlaylistItems(ctx, playlistID, nil)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("playlistID", playlistID).
			Msg("Failed to get playlist items")
		return nil, err
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Int("itemCount", len(items)).
		Msg("Retrieved playlist items")

	return items, nil
}

// ReorderPlaylistItems reorders items in a playlist
func (s *mediaPlaylistService[T]) ReorderPlaylistItems(ctx context.Context, userID uint64, clientID uint64, playlistID string, itemIDs []string) error {
	log := utils.LoggerFromContext(ctx)
	client, err := s.getSpecificPlaylistClient(ctx, userID, clientID)
	if err != nil {
		return err
	}

	playlistProvider, ok := client.(providers.PlaylistProvider)
	if !ok {
		log.Warn().
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("playlistID", playlistID).
			Msg("Client does not support playlists")
		return ErrUnsupportedFeature
	}

	// Reorder items in the playlist
	err = playlistProvider.ReorderPlaylistItems(ctx, playlistID, itemIDs)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("playlistID", playlistID).
			Msg("Failed to reorder playlist items")
		return err
	}

	// Get the updated playlist to update its metadata
	options := &mediatypes.QueryOptions{
		ExternalSourceID: playlistID,
	}
	playlists, err := playlistProvider.GetPlaylists(ctx, options)
	if err == nil && len(playlists) > 0 {
		// Record the change in the playlist metadata
		now := time.Now()
		playlist := &playlists[0]
		playlist.Data.LastModified = now
		playlist.Data.ModifiedBy = clientID

		// Add to change history
		if len(playlist.Data.Items) > 0 {
			// Update the positions based on the new order
			for i, idStr := range itemIDs {
				// Try to convert string ID to uint64 for comparison
				var numericID uint64
				if id, err := strconv.ParseUint(idStr, 10, 64); err == nil {
					numericID = id
				}

				for j, item := range playlist.Data.Items {
					if item.ItemID == numericID {
						// Update position
						playlist.Data.Items[j].Position = i
						playlist.Data.Items[j].LastChanged = now
						playlist.Data.Items[j].ChangeHistory = append(playlist.Data.Items[j].ChangeHistory,
							mediatypes.ChangeRecord{
								ClientID:   clientID,
								ItemID:     idStr,
								ChangeType: "reorder",
								Timestamp:  now,
							})
						break
					}
				}
			}
		}
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Int("itemCount", len(itemIDs)).
		Msg("Reordered playlist items")

	return nil
}

// SyncPlaylist syncs a playlist with all other clients
func (s *mediaPlaylistService[T]) SyncPlaylist(ctx context.Context, userID uint64, clientID uint64, playlistID string) error {
	log := utils.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Msg("Syncing playlist with other clients")

	// TODO: This is a placeholder. In a real implementation, we would implement
	// the logic to sync this playlist with all other clients for this user.
	// This could include:
	// 1. Getting all media clients for the user
	// 2. For each client that supports playlists and is not the source client
	//    a. Find or create the equivalent playlist
	//    b. Map media items between clients using ClientIDs arrays
	//    c. Update the target playlist to match the source playlist
	//    d. Track changes and handle conflicts based on user preferences

	return errors.New("sync playlist not fully implemented")
}

func (s *mediaPlaylistService[T]) SearchPlaylists(ctx context.Context, userID uint64, query string) ([]models.MediaItem[*mediatypes.Playlist], error) {
	clients, err := s.getPlaylistClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allPlaylists []models.MediaItem[*mediatypes.Playlist]

	for _, client := range clients {
		playlistProvider, ok := client.(providers.PlaylistProvider)
		if !ok || !playlistProvider.SupportsPlaylists() {
			continue
		}

		options := &mediatypes.QueryOptions{
			Query: query,
		}

		playlists, err := playlistProvider.GetPlaylists(ctx, options)
		if err != nil {
			continue
		}

		allPlaylists = append(allPlaylists, playlists...)
	}

	return allPlaylists, nil
}

// EnhancedClientMediaPlaylistService extends the basic playlist service with advanced item ID mapping
// This version is aware of the MediaItemRepository and can translate IDs between clients
type EnhancedClientMediaPlaylistService[T types.ClientMediaConfig] struct {
	mediaPlaylistService[T]
	mediaItemRepo repository.ClientMediaItemRepository[mediatypes.MediaData]
}

// NewEnhancedClientMediaPlaylistService creates a new enhanced playlist service with item ID mapping
func NewEnhancedClientMediaPlaylistService[T types.ClientMediaConfig](
	repo repository.ClientRepository[T],
	factory *client.ClientFactoryService,
	mediaItemRepo repository.ClientMediaItemRepository[mediatypes.MediaData],
) ClientMediaPlaylistService[T] {
	return &EnhancedClientMediaPlaylistService[T]{
		mediaPlaylistService: mediaPlaylistService[T]{
			repo:    repo,
			factory: factory,
		},
		mediaItemRepo: mediaItemRepo,
	}
}

// AddItemToPlaylist adds an item to a playlist with proper ID translation
func (s *EnhancedClientMediaPlaylistService[T]) AddItemToPlaylist(ctx context.Context, userID uint64, clientID uint64, playlistID string, itemID string) error {
	log := utils.LoggerFromContext(ctx)
	client, err := s.getSpecificPlaylistClient(ctx, userID, clientID)
	if err != nil {
		return err
	}

	playlistProvider, ok := client.(providers.PlaylistProvider)
	if !ok {
		log.Warn().
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("playlistID", playlistID).
			Msg("Client does not support playlists")
		return ErrUnsupportedFeature
	}

	// Translate the item ID to the client's format if needed
	clientItemID := itemID

	// Check if the itemID looks like a system ID (as a number or UUID)
	// For this example, we check if it might be a numeric ID from our system
	if s.mediaItemRepo != nil {
		// Try to parse as uint64 (internal ID)
		if numericID, err := strconv.ParseUint(itemID, 10, 64); err == nil {
			// Lookup item by internal ID
			mediaItem, err := s.mediaItemRepo.GetByID(ctx, numericID)
			if err == nil {
				// Find the client-specific ID
				for _, cid := range mediaItem.SyncClients {
					if cid.ID == clientID {
						clientItemID = cid.ItemID
						log.Debug().
							Uint64("userID", userID).
							Uint64("clientID", clientID).
							Str("internalID", itemID).
							Str("clientItemID", clientItemID).
							Msg("Translated internal ID to client item ID")
						break
					}
				}
			}
		} else {
			// This might be a client-specific ID from a different client
			// Loop through all clients to find a match
			clientConfigs, err := s.repo.GetByCategory(ctx, types.ClientCategoryMedia, userID)
			if err == nil {
				for _, config := range clientConfigs {
					otherClientID := config.GetID()
					if otherClientID == clientID {
						continue // Skip the current client
					}

					// Try to find a media item with this ID in another client
					mediaItem, err := s.mediaItemRepo.GetByClientItemID(ctx, itemID, otherClientID)
					if err == nil {
						// Found the item, now get its ID for the target client
						for _, cid := range mediaItem.SyncClients {
							if cid.ID == clientID {
								clientItemID = cid.ItemID
								log.Debug().
									Uint64("userID", userID).
									Uint64("clientID", clientID).
									Uint64("sourceClientID", otherClientID).
									Str("sourceItemID", itemID).
									Str("targetItemID", clientItemID).
									Msg("Translated between client IDs")
								break
							}
						}
						break
					}
				}
			}
		}
	}

	// Add item to playlist
	err = playlistProvider.AddItemToPlaylist(ctx, playlistID, clientItemID)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("playlistID", playlistID).
			Str("originalItemID", itemID).
			Str("clientItemID", clientItemID).
			Msg("Failed to add item to playlist")
		return err
	}

	// Get the updated playlist to update its metadata
	options := &mediatypes.QueryOptions{
		ExternalSourceID: playlistID,
	}
	playlists, err := playlistProvider.GetPlaylists(ctx, options)
	if err == nil && len(playlists) > 0 {
		// Record the change in the playlist metadata
		now := time.Now()
		playlist := &playlists[0]
		playlist.Data.LastModified = now
		playlist.Data.ModifiedBy = clientID

		// Add to change history if Items array is used
		if len(playlist.Data.Items) > 0 {
			// Find the item and update its change history
			found := false

			// Try to convert to uint64 for comparison
			var numericClientItemID uint64
			if id, err := strconv.ParseUint(clientItemID, 10, 64); err == nil {
				numericClientItemID = id
			}

			for i, item := range playlist.Data.Items {
				if item.ItemID == numericClientItemID {
					// Update existing item
					playlist.Data.Items[i].LastChanged = now
					playlist.Data.Items[i].ChangeHistory = append(playlist.Data.Items[i].ChangeHistory,
						mediatypes.ChangeRecord{
							ClientID:   clientID,
							ItemID:     clientItemID,
							ChangeType: "add",
							Timestamp:  now,
						})
					found = true
					break
				}
			}

			// If item not found, add it to the Items array
			if !found {
				// Convert string clientItemID to uint64
				numericID, _ := strconv.ParseUint(clientItemID, 10, 64)
				playlist.Data.Items = append(playlist.Data.Items, mediatypes.ListItem{
					ItemID:      numericID,
					Position:    len(playlist.Data.Items),
					LastChanged: now,
					ChangeHistory: []mediatypes.ChangeRecord{
						{
							ClientID:   clientID,
							ItemID:     clientItemID,
							ChangeType: "add",
							Timestamp:  now,
						},
					},
				})
			}
		}
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Str("originalItemID", itemID).
		Str("clientItemID", clientItemID).
		Msg("Added item to playlist")

	return nil
}
