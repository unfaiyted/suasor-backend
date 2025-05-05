package services

import (
	"context"
	"fmt"
	"time"

	"suasor/clients"
	"suasor/clients/media/providers"
	mediatypes "suasor/clients/media/types"
	clienttypes "suasor/clients/types"
	"suasor/repository"
	repobundle "suasor/repository/bundles"
	"suasor/types/models"
	"suasor/utils/logger"
)

// ListSyncService defines the interface for synchronizing lists between local storage and media clients
type ListSyncService[T mediatypes.ListData] interface {
	// SyncToClient syncs a local list to a remote client
	SyncToClient(ctx context.Context, userID uint64, listID uint64, clientID uint64) error

	// SyncFromClient syncs a remote client list to local storage
	SyncFromClient(ctx context.Context, userID uint64, clientID uint64, clientListID string) (*models.MediaItem[T], error)

	// GetSyncStatus gets the sync status for a list
	GetSyncClients(ctx context.Context, listID uint64) (*models.SyncClients, error)

	// UpdateSyncStatus updates the sync status for a list
	UpdateSyncStatus(ctx context.Context, userID uint64, listID uint64, clientID uint64, syncState models.SyncStatus) error
}

// listSyncService implements the ListSyncService interface
type listSyncService[T mediatypes.ListData] struct {
	clientRepos     repobundle.ClientRepositories
	clientFactories *clients.ClientProviderFactoryService
	listService     UserListService[T]
	mediaItemRepo   repository.UserMediaItemRepository[T]
}

// NewListSyncService creates a new list sync service
func NewListSyncService[T mediatypes.ListData](
	clientRepos repobundle.ClientRepositories,
	clientFactories *clients.ClientProviderFactoryService,
	listService UserListService[T],
	mediaItemRepo repository.UserMediaItemRepository[T],
) ListSyncService[T] {
	return &listSyncService[T]{
		clientRepos:     clientRepos,
		clientFactories: clientFactories,
		listService:     listService,
		mediaItemRepo:   mediaItemRepo,
	}
}

// SyncToClient syncs a local list to a remote client
func (s *listSyncService[T]) SyncToClient(ctx context.Context, userID uint64, listID uint64, clientID uint64) error {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Uint64("listID", listID).
		Uint64("clientID", clientID).
		Msg("Syncing list to client")

	// Get the local list
	localList, err := s.listService.GetByID(ctx, listID)
	if err != nil {
		return fmt.Errorf("failed to get local list: %w", err)
	}
	listProvider, config, err := s.getProvider(ctx, clientID)
	if err != nil {
		return fmt.Errorf("failed to get provider: %w", err)
	}

	// Get list details
	itemList := localList.GetData().GetItemList()

	// Check if the list is already synced to this client
	clientListID := ""
	if localList.SyncClients != nil {
		if localList.SyncClients.IsClientPresent(clientID) {
			log.Info().
				Str("clientListID", clientListID).
				Msg("List already exists on client, will update")
			clientListID = localList.SyncClients.GetClientItemID(clientID)
		}

	}

	// Get the items in the list that need to be synced
	items := make([]mediatypes.ListItem, 0, len(itemList.Items))
	strItemIDs := make([]string, 0, len(itemList.Items))
	for _, item := range itemList.Items {
		items = append(items, item)
		strItemIDs = append(strItemIDs, string(item.ItemID))
	}

	// If clientListID is empty, create a new list on the client
	// Otherwise, update the existing list
	var resultClientListID string
	var resultErr error

	if clientListID == "" {
		// Create new list on client
		log.Info().Msg("Creating new list on client")
		mediaItem, err := listProvider.CreateListWithItems(ctx,
			itemList.Details.Title,
			itemList.Details.Description,
			strItemIDs,
		)
		if err != nil {
			return fmt.Errorf("failed to create list on client: %w", err)
		}

		localList.SyncClients.AddClient(clientID, config.GetType(), mediaItem.SyncClients.GetClientItemID(clientID))

	} else {
		// Update existing list on client
		log.Info().
			Str("clientListID", clientListID).
			Msg("Updating existing list on client")
		_, err := listProvider.UpdateList(ctx, clientListID,
			itemList.Details.Title,
			itemList.Details.Description,
			strItemIDs,
		)
		if err != nil {
			return fmt.Errorf("failed to update list on client: %w", err)
		}

		resultClientListID = clientListID
	}

	if resultErr != nil {
		return fmt.Errorf("failed to sync list to client: %w", resultErr)
	}

	// Update local list with sync status
	localList.SyncClients.AddClient(clientID, config.GetType(), resultClientListID)

	// Save updated list
	localList.GetData().SetItemList(*itemList)
	_, err = s.mediaItemRepo.Update(ctx, localList)
	if err != nil {
		return fmt.Errorf("failed to update local list with sync status: %w", err)
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("listID", listID).
		Uint64("clientID", clientID).
		Str("clientListID", resultClientListID).
		Msg("List synchronized successfully to client")

	return nil
}

// SyncFromClient syncs a remote client list to local storage
func (s *listSyncService[T]) SyncFromClient(ctx context.Context, userID uint64, clientID uint64, clientListID string) (*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientListID", clientListID).
		Msg("Syncing list from client")

	// Check client capabilities based on list type
	var zero T
	mediaType := mediatypes.GetMediaTypeFromTypeName(zero)

	listProvider, _, err := s.getProvider(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	// Get list from client
	clientList, err := listProvider.GetList(ctx, clientListID)
	if err != nil {
		return nil, fmt.Errorf("failed to get list from client: %w", err)
	}

	// Check if list already exists locally (by checking if any local list has this clientListID in its sync states)
	existingLists, err := s.mediaItemRepo.Search(ctx, mediatypes.QueryOptions{
		MediaType: mediaType,
		OwnerID:   userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search for existing lists: %w", err)
	}

	// Look for existing list with this client ID
	var existingList *models.MediaItem[T]
	for _, list := range existingLists {
		if list.SyncClients != nil {
			if syncState, exists := list.SyncClients.GetByClientID(clientID); exists && syncState.ItemID == clientListID {
				existingList = list
				break
			}
		}
	}

	if existingList != nil {
		// Update existing list
		log.Info().
			Uint64("listID", existingList.ID).
			Msg("Updating existing local list")

		// Update list details
		existingItemList := existingList.GetData().GetItemList()
		existingItemList.Details.Title = clientList.Title
		existingItemList.Details.Description = clientList.GetData().GetDetails().Description
		existingItemList.IsPublic = clientList.IsPublic

		// Update items
		existingItemList.Items = []mediatypes.ListItem{}

		for i, item := range clientList.GetData().GetItemList().Items {
			existingItemList.Items = append(existingItemList.Items, mediatypes.ListItem{
				ItemID:      item.ItemID,
				Position:    i,
				LastChanged: time.Now(),
				ChangeHistory: []mediatypes.ChangeRecord{
					{
						ClientID:   clientID,
						ItemID:     fmt.Sprintf("%d", item.ItemID),
						ChangeType: "sync",
						Timestamp:  time.Now(),
					},
				},
			})
		}

		// Update sync state
		existingList.SyncClients.UpdateSyncStatus(clientID, models.SyncStatusSuccess)

		// Save updates
		existingList.GetData().SetItemList(*existingItemList)
		result, err := s.mediaItemRepo.Update(ctx, existingList)
		if err != nil {
			return nil, fmt.Errorf("failed to update local list: %w", err)
		}

		return result, nil
	} else {
		// Create new local list
		log.Info().Msg("Creating new local list from client list")

		// Create item list
		itemList := mediatypes.NewItemList(&mediatypes.MediaDetails{
			Title:       clientList.Title,
			Description: clientList.GetData().GetDetails().Description,
			AddedAt:     time.Now(),
			UpdatedAt:   time.Now(),
		})
		itemList.IsPublic = clientList.IsPublic
		itemList.ItemCount = len(clientList.GetData().GetItemList().Items)
		itemList.OwnerID = userID
		itemList.ModifiedBy = userID
		itemList.LastModified = time.Now()
		itemList.Items = clientList.Data.GetItemList().Items

		// Add items
		for i, item := range clientList.Data.GetItemList().Items {
			itemList.Items = append(itemList.Items, mediatypes.ListItem{
				ItemID:      item.ItemID,
				Position:    i,
				LastChanged: time.Now(),
				ChangeHistory: []mediatypes.ChangeRecord{
					{
						ClientID:   clientID,
						ItemID:     fmt.Sprintf("%d", item.ItemID),
						ChangeType: "sync",
						Timestamp:  time.Now(),
					},
				},
			})
		}

		// Create data object
		var data T
		data.SetItemList(itemList)

		// Create media item
		newList := models.NewMediaItem[T](mediaType, data)
		newList.Title = clientList.Title
		newList.OwnerID = userID

		// Save new list
		result, err := s.listService.Create(ctx, userID, newList)
		if err != nil {
			return nil, fmt.Errorf("failed to create local list: %w", err)
		}

		return result, nil
	}
}

// GetSyncStatus gets the sync status for a list
func (s *listSyncService[T]) GetSyncClients(ctx context.Context, listID uint64) (*models.SyncClients, error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("listID", listID).
		Msg("Getting sync status for list")

	// Get the local list
	localList, err := s.listService.GetByID(ctx, listID)
	if err != nil {
		return nil, fmt.Errorf("failed to get local list: %w", err)
	}

	// Get sync states
	if localList.SyncClients == nil {
		return &models.SyncClients{}, nil
	}

	return &localList.SyncClients, nil
}

func (s *listSyncService[T]) GetSyncStatus(ctx context.Context, clientID, listID uint64) (*models.SyncStatus, error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("clientID", clientID).
		Uint64("listID", listID).
		Msg("Getting sync status for list")

	syncClients, err := s.GetSyncClients(ctx, listID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sync clients: %w", err)
	}
	if syncClients == nil {
		return nil, nil
	}

	// Get the sync status for this client
	syncStatus := syncClients.GetSyncStatus(clientID)

	return &syncStatus, nil
}

// UpdateSyncStatus updates the sync status for a list
func (s *listSyncService[T]) UpdateSyncStatus(ctx context.Context, userID uint64, listID uint64, clientID uint64, syncState models.SyncStatus) error {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Uint64("listID", listID).
		Uint64("clientID", clientID).
		Msg("Updating sync status for list")

	// Get the local list
	localList, err := s.listService.GetByID(ctx, listID)
	if err != nil {
		return fmt.Errorf("failed to get local list: %w", err)
	}

	localList.SyncClients.UpdateSyncStatus(clientID, syncState)

	// Save updated list
	_, err = s.mediaItemRepo.Update(ctx, localList)
	if err != nil {
		return fmt.Errorf("failed to update local list with sync status: %w", err)
	}

	log.Info().
		Uint64("userID", userID).
		Uint64("listID", listID).
		Uint64("clientID", clientID).
		Msg("List sync status updated successfully")

	return nil
}

func (s *listSyncService[T]) getProvider(ctx context.Context, clientID uint64) (providers.ListProvider[T], clienttypes.ClientConfig, error) {
	// Get client
	allMediaClients, err := s.clientRepos.GetAllMediaClients(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get client: %w", err)
	}

	// Get client config
	config := allMediaClients.GetClientConfig(clientID)
	if config == nil {
		return nil, nil, fmt.Errorf("client not found")
	}

	var listProvider providers.ListProvider[T]
	var zero T

	mediaType := mediatypes.GetMediaTypeFromTypeName(zero)

	if mediaType == mediatypes.MediaTypePlaylist {
		playlistProvider, err := s.clientFactories.GetListProviderPlaylist(ctx, clientID, config)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get playlist provider: %w", err)
		}
		listProvider = playlistProvider.(providers.ListProvider[T])
	} else if mediaType == mediatypes.MediaTypeCollection {
		collectionProvider, err := s.clientFactories.GetListProviderCollection(ctx, clientID, config)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get collection provider: %w", err)
		}
		listProvider = collectionProvider.(providers.ListProvider[T])
	} else {
		return nil, nil, fmt.Errorf("unsupported list type: %s", mediaType)
	}

	return listProvider, config, nil
}
