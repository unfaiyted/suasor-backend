package services

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"time"

	"suasor/client"
	"suasor/client/media/providers"
	mediatypes "suasor/client/media/types"
	"suasor/client/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/utils"
)

// ClientListService defines operations for interacting with playlist clients
// Its designed to help create and sync with the clients as well as provide a unified interface
// for playlist operations on the integrations side
// Every Get Operations should also save a copy of the MediaItem when it syncs or updates. If the item already exists.
// It will update the existing item, ensuring that the item is has the appropraite IDs and other metadata to sync and keep updated.
type ClientListService[T types.ClientMediaConfig, U mediatypes.ListData] interface {
	CoreListService[U]
	// Client-specific operations
	GetClientList(ctx context.Context, clientID uint64, clientListID string) (*models.MediaItem[U], error)
	GetClientLists(ctx context.Context, clientID uint64, count int) ([]*models.MediaItem[U], error)
	GetClientListsByUserID(ctx context.Context, userID uint64, count int) ([]*models.MediaItem[U], error)

	CreateClientList(ctx context.Context, clientID uint64, name string, description string) (*models.MediaItem[U], error)
	UpdateClientList(ctx context.Context, clientID uint64, clientListID string, name string, description string) (*models.MediaItem[U], error)
	DeleteClientList(ctx context.Context, clientID uint64, clientListID string) error

	// Playlist item operations
	GetClientItems(ctx context.Context, clientID uint64, clientListID string) ([]*models.MediaItem[U], error)
	AddClientItem(ctx context.Context, clientID uint64, clientListID string, itemID string) error
	RemoveClientItem(ctx context.Context, clientID uint64, clientListID string, itemID string) error
	ReorderClientItems(ctx context.Context, clientID uint64, clientListID string, itemIDs []string) error

	// Search and sync operations
	SyncClientList(ctx context.Context, clientID uint64, clientListID string) error
	SearchClientLists(ctx context.Context, clientID uint64, query mediatypes.QueryOptions) ([]*models.MediaItem[U], error)
	SearchUsersClientsLists(ctx context.Context, userID uint64, query mediatypes.QueryOptions) ([]*models.MediaItem[U], error)
	ImportClientList(ctx context.Context, clientID uint64, clientPlaylistID string) (*models.MediaItem[U], error)

	GetSyncStatus(ctx context.Context, clientListID string) (*models.ListSyncStatus, error)
	// SyncToClients(ctx context.Context, clientListID string, clientIDs []uint64) error
	// SyncClientList(ctx context.Context, clientID uint64, clientListID string) error

}

type clientListService[T types.ClientMediaConfig, U mediatypes.ListData] struct {
	listService   CoreListService[U]
	clientRepo    repository.ClientRepository[T]
	clientFactory *client.ClientFactoryService
}

// NewClientPlaylistService creates a new media playlist service
func NewClientListService[T types.ClientMediaConfig, U mediatypes.ListData](
	listService CoreListService[U],
	clientRepo repository.ClientRepository[T],
	clientFactory *client.ClientFactoryService,
) ClientListService[T, U] {
	return &clientListService[T, U]{
		listService:   listService,
		clientRepo:    clientRepo,
		clientFactory: clientFactory,
	}
}

func (s *clientListService[T, U]) GetAll(ctx context.Context, limit int, offset int) ([]*models.MediaItem[U], error) {
	return s.listService.GetAll(ctx, limit, offset)
}

// UpdateItems
func (s *clientListService[T, U]) UpdateItems(ctx context.Context, listID uint64, items []*models.MediaItem[U]) error {
	return s.listService.UpdateItems(ctx, listID, items)
}

// GetRecent
func (s *clientListService[T, U]) GetRecent(ctx context.Context, days int, limit int) ([]*models.MediaItem[U], error) {
	return s.listService.GetRecent(ctx, days, limit)
}

// Sync
func (s *clientListService[T, U]) Sync(ctx context.Context, listID uint64, targetClientIDs []uint64) error {
	return s.listService.Sync(ctx, listID, targetClientIDs)
}

// Search
func (s *clientListService[T, U]) Search(ctx context.Context, query mediatypes.QueryOptions) ([]*models.MediaItem[U], error) {
	return s.listService.Search(ctx, query)
}

// RemoveItem
func (s *clientListService[T, U]) RemoveItem(ctx context.Context, listID uint64, itemID uint64) error {
	return s.listService.RemoveItem(ctx, listID, itemID)
}

// ReorderItems
func (s *clientListService[T, U]) ReorderItems(ctx context.Context, listID uint64, itemIDs []uint64) error {
	return s.listService.ReorderItems(ctx, listID, itemIDs)
}

// GetItems
func (s *clientListService[T, U]) GetItems(ctx context.Context, listID uint64) (*models.MediaItems, error) {
	return s.listService.GetItems(ctx, listID)
}

// Create
func (s *clientListService[T, U]) Create(ctx context.Context, list *models.MediaItem[U]) (*models.MediaItem[U], error) {
	return s.listService.Create(ctx, list)
}

// Update
func (s *clientListService[T, U]) Update(ctx context.Context, list *models.MediaItem[U]) (*models.MediaItem[U], error) {
	return s.listService.Update(ctx, list)
}

// GetByID
func (s *clientListService[T, U]) GetByID(ctx context.Context, id uint64) (*models.MediaItem[U], error) {
	return s.listService.GetByID(ctx, id)
}

// GetByUserID
func (s *clientListService[T, U]) GetByUserID(ctx context.Context, userID uint64, limit int, offset int) ([]*models.MediaItem[U], error) {
	return s.listService.GetByUserID(ctx, userID, limit, offset)
}

// Delete
func (s *clientListService[T, U]) Delete(ctx context.Context, id uint64) error {
	return s.listService.Delete(ctx, id)
}

// AddItem
func (s *clientListService[T, U]) AddItem(ctx context.Context, listID uint64, itemID uint64) error {
	return s.listService.AddItem(ctx, listID, itemID)
}
func (s *clientListService[T, U]) GetClientList(ctx context.Context, clientID uint64, clientListID string) (*models.MediaItem[U], error) {
	log := utils.LoggerFromContext(ctx)
	provider, err := s.getListProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}
	log.Info().
		Uint64("clientID", clientID).
		Msg("Retrieved provider")

	// Get all playlists and find by ID
	options := &mediatypes.QueryOptions{
		ExternalSourceID: clientListID,
	}

	playlists, err := provider.Search(ctx, options)
	if err != nil {
		return nil, err
	}

	// Check if we found any playlists
	if len(playlists) == 0 {
		return nil, errors.New("playlist not found")
	}

	// Return the first matching playlist
	return playlists[0], nil
}
func (s *clientListService[T, U]) GetClientLists(ctx context.Context, clientID uint64, limit int) ([]*models.MediaItem[U], error) {

	provider, err := s.getListProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	options := &mediatypes.QueryOptions{
		Limit: limit,
	}

	playlists, err := provider.Search(ctx, options)
	if err != nil {
		return nil, err
	}

	return playlists, nil
}
func (s *clientListService[T, U]) GetClientListsByUserID(ctx context.Context, userID uint64, count int) ([]*models.MediaItem[U], error) {

	listProviders, err := s.getUserListProviders(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allPlaylists []*models.MediaItem[U]

	for _, provider := range listProviders {
		options := &mediatypes.QueryOptions{
			Limit: count,
		}

		playlists, err := provider.Search(ctx, options)
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
func (s *clientListService[T, U]) CreateClientList(ctx context.Context, clientID uint64, name string, description string) (*models.MediaItem[U], error) {
	log := utils.LoggerFromContext(ctx)
	provider, err := s.getListProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	// Create the playlist
	playlist, err := provider.Create(ctx, name, description)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", clientID).
			Str("name", name).
			Msg("Failed to create playlist")
		return nil, err
	}

	itemList := playlist.Data.GetItemList()
	// Update modification timestamp and client ID
	now := time.Now()
	itemList.LastModified = now
	itemList.ModifiedBy = clientID

	playlist.Data.SetItemList(itemList)

	// Get the client's ID for this playlist - ModifiedBy is just a uint64 client ID
	clientItemID, found := playlist.GetClientItemID(clientID) // Default to Plex as a placeholder
	clientListIDStr := "unknown"
	if found {
		clientListIDStr = clientItemID
	}

	log.Info().
		Uint64("clientID", clientID).
		Str("clientListID", clientListIDStr).
		Str("name", name).
		Msg("Created playlist")

	return playlist, nil
}
func (s *clientListService[T, U]) UpdateClientList(ctx context.Context, clientID uint64, clientListID string, name string, description string) (*models.MediaItem[U], error) {
	log := utils.LoggerFromContext(ctx)
	provider, err := s.getListProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	// Update the playlist
	playlist, err := provider.Update(ctx, clientListID, name, description)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", clientID).
			Str("clientListID", clientListID).
			Msg("Failed to update playlist")
		return nil, err
	}
	itemList := playlist.Data.GetItemList()

	// Update modification timestamp and client ID
	now := time.Now()
	itemList.LastModified = now
	itemList.ModifiedBy = clientID
	playlist.Data.SetItemList(itemList)

	log.Info().
		Uint64("clientID", clientID).
		Str("clientListID", clientListID).
		Str("name", name).
		Msg("Updated playlist")

	return playlist, nil
}
func (s *clientListService[T, U]) DeleteClientList(ctx context.Context, clientID uint64, clientListID string) error {
	log := utils.LoggerFromContext(ctx)
	provider, err := s.getListProvider(ctx, clientID)
	if err != nil {
		return err
	}

	// Delete the playlist
	err = provider.Delete(ctx, clientListID)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", clientID).
			Str("clientListID", clientListID).
			Msg("Failed to delete playlist")
		return err
	}

	log.Info().
		Uint64("clientID", clientID).
		Str("clientListID", clientListID).
		Msg("Deleted playlist")

	return nil
}
func (s *clientListService[T, U]) AddClientItem(ctx context.Context, clientID uint64, clientListID string, itemID string) error {
	log := utils.LoggerFromContext(ctx)
	provider, err := s.getListProvider(ctx, clientID)
	if err != nil {
		return err
	}

	err = provider.AddItem(ctx, clientListID, itemID)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", clientID).
			Str("clientListID", clientListID).
			Str("itemID", itemID).
			Msg("Failed to add item to playlist")
		return err
	}

	// Get the updated playlist to update its metadata
	options := &mediatypes.QueryOptions{
		ExternalSourceID: clientListID,
	}
	lists, err := provider.Search(ctx, options)
	if err == nil && len(lists) > 0 {
		// Record the change in the playlist metadata
		now := time.Now()

		list := lists[0]
		itemList := list.Data.GetItemList()
		itemList.LastModified = now
		itemList.ModifiedBy = clientID

		// list.SetData(list.Data)
		// Add to change history if Items array is used
		if len(itemList.Items) > 0 {
			// Try to convert string itemID to uint64 for comparison
			var numericItemID uint64
			if id, err := strconv.ParseUint(itemID, 10, 64); err == nil {
				numericItemID = id
			}

			// Find the item and update its change history
			for i, item := range itemList.Items {
				if item.ItemID == numericItemID {
					// Update existing item
					itemList.Items[i].LastChanged = now
					itemList.Items[i].ChangeHistory = append(itemList.Items[i].ChangeHistory,
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
		Uint64("clientID", clientID).
		Str("clientListID", clientListID).
		Str("itemID", itemID).
		Msg("Added item to playlist")

	return nil
}
func (s *clientListService[T, U]) RemoveClientItem(ctx context.Context, clientID uint64, clientListID string, itemID string) error {
	log := utils.LoggerFromContext(ctx)
	provider, err := s.getListProvider(ctx, clientID)
	if err != nil {
		return err
	}

	// Remove item from playlist
	err = provider.RemoveItem(ctx, clientListID, itemID)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", clientID).
			Str("clientListID", clientListID).
			Str("itemID", itemID).
			Msg("Failed to remove item from playlist")
		return err
	}

	// The change is already recorded in the client, but we should update our change history
	// This would require additional code to maintain a history in our database

	log.Info().
		Uint64("clientID", clientID).
		Str("clientListID", clientListID).
		Str("itemID", itemID).
		Msg("Removed item from playlist")

	return nil
}

// GetPlaylistItems gets all items in a playlist
func (s *clientListService[T, U]) GetClientItems(ctx context.Context, clientID uint64, clientListID string) ([]*models.MediaItem[U], error) {
	log := utils.LoggerFromContext(ctx)
	provider, err := s.getListProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	// Get all items in the playlist
	items, err := provider.GetItems(ctx, clientListID, nil)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", clientID).
			Str("clientListID", clientListID).
			Msg("Failed to get playlist items")
		return nil, err
	}

	log.Info().
		Uint64("clientID", clientID).
		Str("clientListID", clientListID).
		Int("itemCount", len(items)).
		Msg("Retrieved playlist items")

	return items, nil
}

// ReorderPlaylistItems reorders items in a playlist
func (s *clientListService[T, U]) ReorderClientItems(ctx context.Context, clientID uint64, clientListID string, itemIDs []string) error {
	log := utils.LoggerFromContext(ctx)
	provider, err := s.getListProvider(ctx, clientID)
	if err != nil {
		return err
	}

	// Reorder items in the playlist
	err = provider.ReorderItems(ctx, clientListID, itemIDs)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", clientID).
			Str("clientListID", clientListID).
			Msg("Failed to reorder playlist items")
		return err
	}

	// Get the updated playlist to update its metadata
	options := &mediatypes.QueryOptions{
		ExternalSourceID: clientListID,
	}
	playlists, err := provider.Search(ctx, options)
	if err == nil && len(playlists) > 0 {
		// Record the change in the playlist metadata
		now := time.Now()
		playlist := playlists[0]
		itemList := playlist.Data.GetItemList()
		itemList.LastModified = now
		itemList.ModifiedBy = clientID

		// Add to change history
		if len(itemList.Items) > 0 {
			// Update the positions based on the new order
			for i, idStr := range itemIDs {
				// Try to convert string ID to uint64 for comparison
				var numericID uint64
				if id, err := strconv.ParseUint(idStr, 10, 64); err == nil {
					numericID = id
				}

				for j, item := range itemList.Items {
					if item.ItemID == numericID {
						// Update position
						itemList.Items[j].Position = i
						itemList.Items[j].LastChanged = now
						itemList.Items[j].ChangeHistory = append(itemList.Items[j].ChangeHistory,
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
		Uint64("clientID", clientID).
		Str("clientListID", clientListID).
		Int("itemCount", len(itemIDs)).
		Msg("Reordered playlist items")

	return nil
}

// Sync list syncs a list with the client
func (s *clientListService[T, U]) SyncClientList(ctx context.Context, clientID uint64, clientListID string) error {
	log := utils.LoggerFromContext(ctx)
	log.Info().
		Uint64("clientID", clientID).
		Str("clientListID", clientListID).
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

// GetSyncStatus retrieves the sync status of a playlist across clients
func (s *clientListService[T, U]) GetSyncStatus(ctx context.Context, clientListID string) (*models.ListSyncStatus, error) {
	// log := utils.LoggerFromContext(ctx)
	// log.Debug().
	// 	Str("clientListID", clientListID).
	// 	Msg("Getting playlist sync status")
	//
	// // Get the playlist
	// playlist, err := s.GetClientList(ctx, clientListID)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get playlist sync status: %w", err)
	// }
	//
	// // Verify user has permission to view this playlist
	// userID := ctx.Value("userID").(uint64)
	// if !s.hasPlaylistReadPermission(ctx, userID, playlist) {
	// 	log.Warn().
	// 		Uint64("playlistID", playlist.ID).
	// 		Uint64("ownerID", playlist.Data.OwnerID).
	// 		Uint64("requestingUserID", userID).
	// 		Msg("User attempting to view playlist sync status without permission")
	// 	return nil, errors.New("you don't have permission to view this playlist's sync status")
	//
	// }
	// // Get the sync status
	// status, err := s.repo.GetPlaylistSyncStatus(ctx, playlist.ID)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get playlist sync status: %w", err)
	// }
	// TODO: Implement and test playlist sync status
	return nil, nil
}

// ImportClientList imports a playlist from a different client
func (s *clientListService[T, U]) ImportClientList(ctx context.Context, clientID uint64, clientPlaylistID string) (*models.MediaItem[U], error) {
	// log := utils.LoggerFromContext(ctx)
	// provider, err := s.getListProvider(ctx, clientID)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// playlist, err := provider.ImportPlaylist(ctx, clientPlaylistID)
	// if err != nil {
	// 	log.Error().
	// 		Err(err).
	// 		Uint64("clientID", clientID).
	// 		Str("clientListID", clientPlaylistID).
	// 		Msg("Failed to import playlist")
	// 	return nil, fmt.Errorf("failed to import playlist: %w", err)
	// }
	//
	// return playlist, nil
	return nil, nil
}
func (s *clientListService[T, U]) SearchClientLists(ctx context.Context, clientID uint64, query mediatypes.QueryOptions) ([]*models.MediaItem[U], error) {

	listProvider, err := s.getListProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}
	lists, err := listProvider.Search(ctx, &query)
	if err != nil {
		return nil, err
	}

	return lists, nil
}
func (s *clientListService[T, U]) SearchUsersClientsLists(ctx context.Context, userID uint64, query mediatypes.QueryOptions) ([]*models.MediaItem[U], error) {

	providers, err := s.getUserListProviders(ctx, query.OwnerID)
	if err != nil {
		return nil, err
	}

	var allPlaylists []*models.MediaItem[U]

	for _, listProvider := range providers {

		playlists, err := listProvider.Search(ctx, &query)
		if err != nil {
			continue
		}

		allPlaylists = append(allPlaylists, playlists...)
	}

	return allPlaylists, nil
}

// getSpecificPlaylistClient gets a specific playlist client
func (s *clientListService[T, U]) getListProvider(ctx context.Context, clientID uint64) (providers.ListProvider[U], error) {
	log := utils.LoggerFromContext(ctx)

	clientConfig, err := s.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return nil, err
	}
	log.Debug().
		Uint64("clientID", clientID).
		Str("clientType", clientConfig.Type.String()).
		Msg("Retrieved client config")

	if !clientConfig.Config.Data.SupportsPlaylists() && !clientConfig.Config.Data.SupportsCollections() {
		log.Warn().
			Uint64("clientID", clientID).
			Str("clientType", clientConfig.Config.Data.GetType().String()).
			Msg("Client does not support lists")
		return nil, ErrUnsupportedFeature
	}

	log.Debug().
		Uint64("clientID", clientID).
		Str("clientType", clientConfig.Config.Data.GetType().String()).
		Msg("Client supports lists")

	client, err := s.clientFactory.GetClient(ctx, clientID, clientConfig.Config.Data)
	if err != nil {
		return nil, err
	}
	log.Debug().
		Uint64("clientID", clientID).
		Str("clientType", clientConfig.Config.Data.GetType().String()).
		Msg("Retrieved client")
	return client.(providers.ListProvider[U]), nil
}
func (s *clientListService[T, U]) getUserListProviders(ctx context.Context, userID uint64) ([]providers.ListProvider[U], error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("userID", userID).
		Msg("Retrieving playlist providers for user")

	// Get all media clients for the user
	clients, err := s.clientRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var providers []providers.ListProvider[U]
	for _, clientConfig := range clients {
		if clientConfig.Config.Data.SupportsPlaylists() {
			clientID := clientConfig.GetID()
			provider, err := s.getListProvider(ctx, clientID)
			if err != nil {
				// Log error but continue with other clients
				continue
			}
			providers = append(providers, provider)
		}
	}
	log.Info().
		Uint64("userID", userID).
		Int("count", len(providers)).
		Msg("Retrieved playlist providers for user")
	return providers, nil

}
