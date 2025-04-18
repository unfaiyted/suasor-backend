package services

// import (
// 	"context"
// 	"errors"
// 	"sort"
// 	"strconv"
// 	"time"
//
// 	"suasor/client"
// 	"suasor/client/media/providers"
// 	mediatypes "suasor/client/media/types"
// 	"suasor/client/types"
// 	"suasor/repository"
// 	"suasor/types/models"
// 	"suasor/utils"
// )

// ClientPlaylistService defines operations for interacting with playlist clients
// Its designed to help create and sync with the clients as well as provide a unified interface
// for playlist operations on the integrations side
// Every Get Operations should also save a copy of the MediaItem when it syncs or updates. If the item already exists.
// It will update the existing item, ensuring that the item is has the appropraite IDs and other metadata to sync and keep updated.
// type ClientPlaylistService[T types.ClientConfig] interface {
// 	ClientListService[T, mediatypes.Playlist]
// }
//
// type mediaPlaylistService[T types.ClientMediaConfig] struct {
// 	playlistService PlaylistService
// 	clientRepo      repository.ClientRepository[T]
// 	clientFactory   *client.ClientFactoryService
// }
//
// // NewClientPlaylistService creates a new media playlist service
// func NewClientPlaylistService[T types.ClientMediaConfig](
// 	 PlaylistService,
// 	clientRepo repository.ClientRepository[T],
// 	clientFactory *client.ClientFactoryService,
// ) ClientPlaylistService[T] {
// 	return &mediaPlaylistService[T]{
// 		playlistService: playlistService,
// 		clientRepo:      clientRepo,
// 		clientFactory:   clientFactory,
// 	}
// }
//
// // getSpecificPlaylistClient gets a specific playlist client
// func (s *mediaPlaylistService[T]) getPlaylistProvider(ctx context.Context, clientID uint64) (providers.PlaylistProvider, error) {
// 	log := utils.LoggerFromContext(ctx)
//
// 	clientConfig, err := s.clientRepo.GetByID(ctx, clientID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	log.Debug().
// 		Uint64("clientID", clientID).
// 		Str("clientType", clientConfig.Type.String()).
// 		Msg("Retrieved client config")
//
// 	if !clientConfig.Config.Data.SupportsPlaylists() {
// 		log.Warn().
// 			Uint64("clientID", clientID).
// 			Str("clientType", clientConfig.Config.Data.GetType().String()).
// 			Msg("Client does not support playlists")
// 		return nil, ErrUnsupportedFeature
// 	}
//
// 	log.Debug().
// 		Uint64("clientID", clientID).
// 		Str("clientType", clientConfig.Config.Data.GetType().String()).
// 		Msg("Client supports playlists")
//
// 	client, err := s.clientFactory.GetClient(ctx, clientID, clientConfig.Config.Data)
// 	if err != nil {
// 		return nil, err
// 	}
// 	log.Debug().
// 		Uint64("clientID", clientID).
// 		Str("clientType", clientConfig.Config.Data.GetType().String()).
// 		Msg("Retrieved client")
// 	return client.(providers.PlaylistProvider), nil
// }
//
// func (s *mediaPlaylistService[T]) getUserPlaylistProviders(ctx context.Context, userID uint64) ([]providers.PlaylistProvider, error) {
// 	log := utils.LoggerFromContext(ctx)
//
// 	log.Info().
// 		Uint64("userID", userID).
// 		Msg("Retrieving playlist providers for user")
//
// 	// Get all media clients for the user
// 	clients, err := s.clientRepo.GetByUserID(ctx, userID)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	var providers []providers.PlaylistProvider
// 	for _, clientConfig := range clients {
// 		if clientConfig.Config.Data.SupportsPlaylists() {
// 			clientID := clientConfig.GetID()
// 			provider, err := s.getPlaylistProvider(ctx, clientID)
// 			if err != nil {
// 				// Log error but continue with other clients
// 				continue
// 			}
// 			providers = append(providers, provider)
// 		}
// 	}
// 	log.Info().
// 		Uint64("userID", userID).
// 		Int("count", len(providers)).
// 		Msg("Retrieved playlist providers for user")
// 	return providers, nil
//
// }
// func (s *mediaPlaylistService[T]) GetClientList(ctx context.Context, clientID uint64, clientListID string) (*models.MediaItem[*mediatypes.Playlist], error) {
// 	log := utils.LoggerFromContext(ctx)
// 	provider, err := s.getPlaylistProvider(ctx, clientID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	log.Info().
// 		Uint64("clientID", clientID).
// 		Msg("Retrieved provider")
//
// 	// Get all playlists and find by ID
// 	options := &mediatypes.QueryOptions{
// 		ExternalSourceID: clientListID,
// 	}
//
// 	playlists, err := provider.Search(ctx, options)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	// Check if we found any playlists
// 	if len(playlists) == 0 {
// 		return nil, errors.New("playlist not found")
// 	}
//
// 	// Return the first matching playlist
// 	return &playlists[0], nil
// }
//
// func (s *mediaPlaylistService[T]) GetClientLists(ctx context.Context, clientID uint64, limit int) ([]models.MediaItem[*mediatypes.Playlist], error) {
//
// 	provider, err := s.getPlaylistProvider(ctx, clientID)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	options := &mediatypes.QueryOptions{
// 		Limit: limit,
// 	}
//
// 	playlists, err := provider.Search(ctx, options)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return playlists, nil
// }
//
// func (s *mediaPlaylistService[T]) GetClientListsByUserID(ctx context.Context, userID uint64, count int) ([]models.MediaItem[*mediatypes.Playlist], error) {
//
// 	listProviders, err := s.getUserPlaylistProviders(ctx, userID)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	var allPlaylists []models.MediaItem[*mediatypes.Playlist]
//
// 	for _, provider := range listProviders {
// 		options := &mediatypes.QueryOptions{
// 			Limit: count,
// 		}
//
// 		playlists, err := provider.Search(ctx, options)
// 		if err != nil {
// 			continue
// 		}
//
// 		allPlaylists = append(allPlaylists, playlists...)
// 	}
//
// 	// Sort by added date
// 	sort.Slice(allPlaylists, func(i, j int) bool {
// 		return allPlaylists[i].Data.GetDetails().AddedAt.After(allPlaylists[j].Data.GetDetails().AddedAt)
// 	})
//
// 	// Limit to requested count if specified
// 	if count > 0 && len(allPlaylists) > count {
// 		allPlaylists = allPlaylists[:count]
// 	}
//
// 	return allPlaylists, nil
// }
//
// func (s *mediaPlaylistService[T]) CreateClientList(ctx context.Context, clientID uint64, name string, description string) (*models.MediaItem[*mediatypes.Playlist], error) {
// 	log := utils.LoggerFromContext(ctx)
// 	provider, err := s.getPlaylistProvider(ctx, clientID)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	// Create the playlist
// 	playlist, err := provider.CreatePlaylist(ctx, name, description)
// 	if err != nil {
// 		log.Error().
// 			Err(err).
// 			Uint64("clientID", clientID).
// 			Str("name", name).
// 			Msg("Failed to create playlist")
// 		return nil, err
// 	}
//
// 	// Update modification timestamp and client ID
// 	now := time.Now()
// 	playlist.Data.ItemList.LastModified = now
// 	playlist.Data.ItemList.ModifiedBy = clientID
//
// 	// Get the client's ID for this playlist - ModifiedBy is just a uint64 client ID
// 	clientItemID, found := playlist.GetClientItemID(clientID) // Default to Plex as a placeholder
// 	clientListIDStr := "unknown"
// 	if found {
// 		clientListIDStr = clientItemID
// 	}
//
// 	log.Info().
// 		Uint64("clientID", clientID).
// 		Str("clientListID", clientListIDStr).
// 		Str("name", name).
// 		Msg("Created playlist")
//
// 	return playlist, nil
// }
//
// func (s *mediaPlaylistService[T]) UpdateClientList(ctx context.Context, clientID uint64, clientListID string, name string, description string) (*models.MediaItem[*mediatypes.Playlist], error) {
// 	log := utils.LoggerFromContext(ctx)
// 	provider, err := s.getPlaylistProvider(ctx, clientID)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	// Update the playlist
// 	playlist, err := provider.UpdatePlaylist(ctx, clientListID, name, description)
// 	if err != nil {
// 		log.Error().
// 			Err(err).
// 			Uint64("clientID", clientID).
// 			Str("clientListID", clientListID).
// 			Msg("Failed to update playlist")
// 		return nil, err
// 	}
//
// 	// Update modification timestamp and client ID
// 	now := time.Now()
// 	playlist.Data.ItemList.LastModified = now
// 	playlist.Data.ItemList.ModifiedBy = clientID
//
// 	log.Info().
// 		Uint64("clientID", clientID).
// 		Str("clientListID", clientListID).
// 		Str("name", name).
// 		Msg("Updated playlist")
//
// 	return playlist, nil
// }
//
// func (s *mediaPlaylistService[T]) DeleteClientList(ctx context.Context, clientID uint64, clientListID string) error {
// 	log := utils.LoggerFromContext(ctx)
// 	provider, err := s.getPlaylistProvider(ctx, clientID)
// 	if err != nil {
// 		return err
// 	}
//
// 	// Delete the playlist
// 	err = provider.DeletePlaylist(ctx, clientListID)
// 	if err != nil {
// 		log.Error().
// 			Err(err).
// 			Uint64("clientID", clientID).
// 			Str("clientListID", clientListID).
// 			Msg("Failed to delete playlist")
// 		return err
// 	}
//
// 	log.Info().
// 		Uint64("clientID", clientID).
// 		Str("clientListID", clientListID).
// 		Msg("Deleted playlist")
//
// 	return nil
// }
//
// func (s *mediaPlaylistService[T]) AddClientItem(ctx context.Context, clientID uint64, clientListID string, itemID string) error {
// 	log := utils.LoggerFromContext(ctx)
// 	provider, err := s.getPlaylistProvider(ctx, clientID)
// 	if err != nil {
// 		return err
// 	}
//
// 	err = provider.AddItemToPlaylist(ctx, clientListID, itemID)
// 	if err != nil {
// 		log.Error().
// 			Err(err).
// 			Uint64("clientID", clientID).
// 			Str("clientListID", clientListID).
// 			Str("itemID", itemID).
// 			Msg("Failed to add item to playlist")
// 		return err
// 	}
//
// 	// Get the updated playlist to update its metadata
// 	options := &mediatypes.QueryOptions{
// 		ExternalSourceID: clientListID,
// 	}
// 	playlists, err := provider.Search(ctx, options)
// 	if err == nil && len(playlists) > 0 {
// 		// Record the change in the playlist metadata
// 		now := time.Now()
// 		playlist := &playlists[0]
// 		playlist.Data.LastModified = now
// 		playlist.Data.ModifiedBy = clientID
//
// 		// Add to change history if Items array is used
// 		if len(playlist.Data.Items) > 0 {
// 			// Try to convert string itemID to uint64 for comparison
// 			var numericItemID uint64
// 			if id, err := strconv.ParseUint(itemID, 10, 64); err == nil {
// 				numericItemID = id
// 			}
//
// 			// Find the item and update its change history
// 			for i, item := range playlist.Data.Items {
// 				if item.ItemID == numericItemID {
// 					// Update existing item
// 					playlist.Data.Items[i].LastChanged = now
// 					playlist.Data.Items[i].ChangeHistory = append(playlist.Data.Items[i].ChangeHistory,
// 						mediatypes.ChangeRecord{
// 							ClientID:   clientID,
// 							ItemID:     itemID,
// 							ChangeType: "add",
// 							Timestamp:  now,
// 						})
// 					break
// 				}
// 			}
// 		}
// 	}
//
// 	log.Info().
// 		Uint64("clientID", clientID).
// 		Str("clientListID", clientListID).
// 		Str("itemID", itemID).
// 		Msg("Added item to playlist")
//
// 	return nil
// }
//
// func (s *mediaPlaylistService[T]) RemoveClientItem(ctx context.Context, clientID uint64, clientListID string, itemID string) error {
// 	log := utils.LoggerFromContext(ctx)
// 	provider, err := s.getPlaylistProvider(ctx, clientID)
// 	if err != nil {
// 		return err
// 	}
//
// 	// Remove item from playlist
// 	err = provider.RemoveItemFromPlaylist(ctx, clientListID, itemID)
// 	if err != nil {
// 		log.Error().
// 			Err(err).
// 			Uint64("clientID", clientID).
// 			Str("clientListID", clientListID).
// 			Str("itemID", itemID).
// 			Msg("Failed to remove item from playlist")
// 		return err
// 	}
//
// 	// The change is already recorded in the client, but we should update our change history
// 	// This would require additional code to maintain a history in our database
//
// 	log.Info().
// 		Uint64("clientID", clientID).
// 		Str("clientListID", clientListID).
// 		Str("itemID", itemID).
// 		Msg("Removed item from playlist")
//
// 	return nil
// }
//
// // GetPlaylistItems gets all items in a playlist
// func (s *mediaPlaylistService[T]) GetClientItems(ctx context.Context, clientID uint64, clientListID string) ([]models.MediaItem[mediatypes.MediaData], error) {
// 	log := utils.LoggerFromContext(ctx)
// 	provider, err := s.getPlaylistProvider(ctx, clientID)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	// Get all items in the playlist
// 	items, err := provider.GetPlaylistItems(ctx, clientListID, nil)
// 	if err != nil {
// 		log.Error().
// 			Err(err).
// 			Uint64("clientID", clientID).
// 			Str("clientListID", clientListID).
// 			Msg("Failed to get playlist items")
// 		return nil, err
// 	}
//
// 	log.Info().
// 		Uint64("clientID", clientID).
// 		Str("clientListID", clientListID).
// 		Int("itemCount", len(items)).
// 		Msg("Retrieved playlist items")
//
// 	return items, nil
// }
//
// // ReorderPlaylistItems reorders items in a playlist
// func (s *mediaPlaylistService[T]) ReorderClientItems(ctx context.Context, clientID uint64, clientListID string, itemIDs []string) error {
// 	log := utils.LoggerFromContext(ctx)
// 	provider, err := s.getPlaylistProvider(ctx, clientID)
// 	if err != nil {
// 		return err
// 	}
//
// 	// Reorder items in the playlist
// 	err = provider.ReorderPlaylistItems(ctx, clientListID, itemIDs)
// 	if err != nil {
// 		log.Error().
// 			Err(err).
// 			Uint64("clientID", clientID).
// 			Str("clientListID", clientListID).
// 			Msg("Failed to reorder playlist items")
// 		return err
// 	}
//
// 	// Get the updated playlist to update its metadata
// 	options := &mediatypes.QueryOptions{
// 		ExternalSourceID: clientListID,
// 	}
// 	playlists, err := provider.Search(ctx, options)
// 	if err == nil && len(playlists) > 0 {
// 		// Record the change in the playlist metadata
// 		now := time.Now()
// 		playlist := &playlists[0]
// 		playlist.Data.LastModified = now
// 		playlist.Data.ModifiedBy = clientID
//
// 		// Add to change history
// 		if len(playlist.Data.Items) > 0 {
// 			// Update the positions based on the new order
// 			for i, idStr := range itemIDs {
// 				// Try to convert string ID to uint64 for comparison
// 				var numericID uint64
// 				if id, err := strconv.ParseUint(idStr, 10, 64); err == nil {
// 					numericID = id
// 				}
//
// 				for j, item := range playlist.Data.Items {
// 					if item.ItemID == numericID {
// 						// Update position
// 						playlist.Data.Items[j].Position = i
// 						playlist.Data.Items[j].LastChanged = now
// 						playlist.Data.Items[j].ChangeHistory = append(playlist.Data.Items[j].ChangeHistory,
// 							mediatypes.ChangeRecord{
// 								ClientID:   clientID,
// 								ItemID:     idStr,
// 								ChangeType: "reorder",
// 								Timestamp:  now,
// 							})
// 						break
// 					}
// 				}
// 			}
// 		}
// 	}
//
// 	log.Info().
// 		Uint64("clientID", clientID).
// 		Str("clientListID", clientListID).
// 		Int("itemCount", len(itemIDs)).
// 		Msg("Reordered playlist items")
//
// 	return nil
// }
//
// // SyncPlaylist syncs a playlist with all other clients
// func (s *mediaPlaylistService[T]) Sync(ctx context.Context, clientID uint64, clientListID string) error {
// 	log := utils.LoggerFromContext(ctx)
// 	log.Info().
// 		Uint64("clientID", clientID).
// 		Str("clientListID", clientListID).
// 		Msg("Syncing playlist with other clients")
//
// 	// TODO: This is a placeholder. In a real implementation, we would implement
// 	// the logic to sync this playlist with all other clients for this user.
// 	// This could include:
// 	// 1. Getting all media clients for the user
// 	// 2. For each client that supports playlists and is not the source client
// 	//    a. Find or create the equivalent playlist
// 	//    b. Map media items between clients using ClientIDs arrays
// 	//    c. Update the target playlist to match the source playlist
// 	//    d. Track changes and handle conflicts based on user preferences
//
// 	return errors.New("sync playlist not fully implemented")
// }
//
// // GetSyncStatus retrieves the sync status of a playlist across clients
// func (s *mediaPlaylistService[T]) GetSyncStatus(ctx context.Context, clientListID string) (*models.PlaylistSyncStatus, error) {
// 	// log := utils.LoggerFromContext(ctx)
// 	// log.Debug().
// 	// 	Str("clientListID", clientListID).
// 	// 	Msg("Getting playlist sync status")
// 	//
// 	// // Get the playlist
// 	// playlist, err := s.GetClientList(ctx, clientListID)
// 	// if err != nil {
// 	// 	return nil, fmt.Errorf("failed to get playlist sync status: %w", err)
// 	// }
// 	//
// 	// // Verify user has permission to view this playlist
// 	// userID := ctx.Value("userID").(uint64)
// 	// if !s.hasPlaylistReadPermission(ctx, userID, playlist) {
// 	// 	log.Warn().
// 	// 		Uint64("playlistID", playlist.ID).
// 	// 		Uint64("ownerID", playlist.Data.OwnerID).
// 	// 		Uint64("requestingUserID", userID).
// 	// 		Msg("User attempting to view playlist sync status without permission")
// 	// 	return nil, errors.New("you don't have permission to view this playlist's sync status")
// 	//
// 	// }
// 	// // Get the sync status
// 	// status, err := s.repo.GetPlaylistSyncStatus(ctx, playlist.ID)
// 	// if err != nil {
// 	// 	return nil, fmt.Errorf("failed to get playlist sync status: %w", err)
// 	// }
// 	// TODO: Implement and test playlist sync status
// 	return nil, nil
// }
//
// // ImportClientList imports a playlist from a different client
// func (s *mediaPlaylistService[T]) ImportClientList(ctx context.Context, clientID uint64, clientPlaylistID string) (*models.MediaItem[*mediatypes.Playlist], error) {
// 	// log := utils.LoggerFromContext(ctx)
// 	// provider, err := s.getPlaylistProvider(ctx, clientID)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	//
// 	// playlist, err := provider.ImportPlaylist(ctx, clientPlaylistID)
// 	// if err != nil {
// 	// 	log.Error().
// 	// 		Err(err).
// 	// 		Uint64("clientID", clientID).
// 	// 		Str("clientListID", clientPlaylistID).
// 	// 		Msg("Failed to import playlist")
// 	// 	return nil, fmt.Errorf("failed to import playlist: %w", err)
// 	// }
// 	//
// 	// return playlist, nil
// 	return nil, nil
// }
//
// func (s *mediaPlaylistService[T]) SearchClientLists(ctx context.Context, clientID uint64, query mediatypes.QueryOptions) ([]models.MediaItem[*mediatypes.Playlist], error) {
//
// 	listProvider, err := s.getPlaylistProvider(ctx, clientID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	playlists, err := listProvider.Search(ctx, &query)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return playlists, nil
// }
//
// func (s *mediaPlaylistService[T]) SearchUsersClientsLists(ctx context.Context, userID uint64, query mediatypes.QueryOptions) ([]models.MediaItem[*mediatypes.Playlist], error) {
//
// 	providers, err := s.getUserPlaylistProviders(ctx, query.OwnerID)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	var allPlaylists []models.MediaItem[*mediatypes.Playlist]
//
// 	for _, listProvider := range providers {
//
// 		playlists, err := listProvider.Search(ctx, &query)
// 		if err != nil {
// 			continue
// 		}
//
// 		allPlaylists = append(allPlaylists, playlists...)
// 	}
//
// 	return allPlaylists, nil
// }
