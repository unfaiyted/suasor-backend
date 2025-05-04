package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"suasor/clients/media"
	mediatypes "suasor/clients/media/types"
	"suasor/clients/media/providers"
	"suasor/repository"
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
	GetSyncStatus(ctx context.Context, listID uint64) (mediatypes.ListSyncStates, error)
	
	// UpdateSyncStatus updates the sync status for a list
	UpdateSyncStatus(ctx context.Context, userID uint64, listID uint64, clientID uint64, syncState mediatypes.SyncState) error
}

// listSyncService implements the ListSyncService interface
type listSyncService[T mediatypes.ListData] struct {
	clientRepo    repository.ClientRepository
	listService   UserListService[T]
	mediaItemRepo repository.UserMediaItemRepository[T]
}

// NewListSyncService creates a new list sync service
func NewListSyncService[T mediatypes.ListData](
	clientRepo repository.ClientRepository,
	listService UserListService[T],
	mediaItemRepo repository.UserMediaItemRepository[T],
) ListSyncService[T] {
	return &listSyncService[T]{
		clientRepo:    clientRepo,
		listService:   listService,
		mediaItemRepo: mediaItemRepo,
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

	// Get client
	client, err := s.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	// Check if client is a media client that supports lists
	mediaClient, ok := client.Client.(media.ClientMedia)
	if !ok {
		return errors.New("client does not support media operations")
	}

	var listProvider providers.ListProvider
	
	// Check client capabilities based on list type
	var zero T
	mediaType := mediatypes.GetMediaTypeFromTypeName(zero)
	
	if mediaType == mediatypes.MediaTypePlaylist {
		if !mediaClient.SupportsPlaylists() {
			return errors.New("client does not support playlists")
		}
		// Get playlist provider
		var err error
		listProvider, err = mediaClient.GetPlaylistProvider()
		if err != nil {
			return fmt.Errorf("failed to get playlist provider: %w", err)
		}
	} else if mediaType == mediatypes.MediaTypeCollection {
		if !mediaClient.SupportsCollections() {
			return errors.New("client does not support collections")
		}
		// Get collection provider
		var err error
		listProvider, err = mediaClient.GetCollectionProvider()
		if err != nil {
			return fmt.Errorf("failed to get collection provider: %w", err)
		}
	} else {
		return fmt.Errorf("unsupported list type: %s", mediaType)
	}

	// Get list details
	itemList := localList.GetData().GetItemList()
	
	// Check if the list is already synced to this client
	clientListID := ""
	if itemList.SyncStates != nil {
		if syncState, exists := itemList.SyncStates[clientID]; exists && syncState.ClientListID != "" {
			clientListID = syncState.ClientListID
			log.Info().
				Str("clientListID", clientListID).
				Msg("List already exists on client, will update")
		}
	}

	// Get the items in the list that need to be synced
	items := make([]mediatypes.ListItem, 0, len(itemList.Items))
	for _, item := range itemList.Items {
		items = append(items, item)
	}

	// If clientListID is empty, create a new list on the client
	// Otherwise, update the existing list
	var resultClientListID string
	var resultErr error
	
	if clientListID == "" {
		// Create new list on client
		log.Info().Msg("Creating new list on client")
		resultClientListID, resultErr = listProvider.CreateList(ctx, &mediatypes.List{
			Name:        itemList.Details.Title,
			Description: itemList.Details.Description,
			Items:       items,
			IsPublic:    itemList.IsPublic,
		})
	} else {
		// Update existing list on client
		log.Info().
			Str("clientListID", clientListID).
			Msg("Updating existing list on client")
		resultErr = listProvider.UpdateList(ctx, clientListID, &mediatypes.List{
			Name:        itemList.Details.Title,
			Description: itemList.Details.Description,
			Items:       items,
			IsPublic:    itemList.IsPublic,
		})
		resultClientListID = clientListID
	}

	if resultErr != nil {
		return fmt.Errorf("failed to sync list to client: %w", resultErr)
	}

	// Update local list with sync status
	if itemList.SyncStates == nil {
		itemList.SyncStates = make(mediatypes.ListSyncStates)
	}
	
	itemList.SyncStates[clientID] = mediatypes.SyncState{
		ClientListID: resultClientListID,
		LastSynced:   time.Now(),
		Status:       mediatypes.SyncStatusSuccess,
	}
	
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

	// Get client
	client, err := s.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	// Check if client is a media client that supports lists
	mediaClient, ok := client.Client.(media.ClientMedia)
	if !ok {
		return nil, errors.New("client does not support media operations")
	}

	var listProvider providers.ListProvider
	
	// Check client capabilities based on list type
	var zero T
	mediaType := mediatypes.GetMediaTypeFromTypeName(zero)
	
	if mediaType == mediatypes.MediaTypePlaylist {
		if !mediaClient.SupportsPlaylists() {
			return nil, errors.New("client does not support playlists")
		}
		// Get playlist provider
		var err error
		listProvider, err = mediaClient.GetPlaylistProvider()
		if err != nil {
			return nil, fmt.Errorf("failed to get playlist provider: %w", err)
		}
	} else if mediaType == mediatypes.MediaTypeCollection {
		if !mediaClient.SupportsCollections() {
			return nil, errors.New("client does not support collections")
		}
		// Get collection provider
		var err error
		listProvider, err = mediaClient.GetCollectionProvider()
		if err != nil {
			return nil, fmt.Errorf("failed to get collection provider: %w", err)
		}
	} else {
		return nil, fmt.Errorf("unsupported list type: %s", mediaType)
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
		itemList := list.GetData().GetItemList()
		if itemList.SyncStates != nil {
			if syncState, exists := itemList.SyncStates[clientID]; exists && syncState.ClientListID == clientListID {
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
		itemList := existingList.GetData().GetItemList()
		itemList.Details.Title = clientList.Name
		itemList.Details.Description = clientList.Description
		itemList.IsPublic = clientList.IsPublic
		
		// Update items
		itemList.Items = []mediatypes.ListItem{}
		for i, item := range clientList.Items {
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
		
		// Update sync state
		if itemList.SyncStates == nil {
			itemList.SyncStates = make(mediatypes.ListSyncStates)
		}
		itemList.SyncStates[clientID] = mediatypes.SyncState{
			ClientListID: clientListID,
			LastSynced:   time.Now(),
			Status:       mediatypes.SyncStatusSuccess,
		}
		
		// Save updates
		existingList.GetData().SetItemList(*itemList)
		result, err := s.mediaItemRepo.Update(ctx, existingList)
		if err != nil {
			return nil, fmt.Errorf("failed to update local list: %w", err)
		}
		
		return result, nil
	} else {
		// Create new local list
		log.Info().Msg("Creating new local list from client list")
		
		// Create item list
		itemList := mediatypes.ItemList{
			Details: &mediatypes.MediaDetails{
				Title:       clientList.Name,
				Description: clientList.Description,
				AddedAt:     time.Now(),
				UpdatedAt:   time.Now(),
			},
			IsPublic:     clientList.IsPublic,
			ItemCount:    len(clientList.Items),
			OwnerID:      userID,
			ModifiedBy:   userID,
			LastModified: time.Now(),
			Items:        []mediatypes.ListItem{},
			SyncStates: mediatypes.ListSyncStates{
				clientID: mediatypes.SyncState{
					ClientListID: clientListID,
					LastSynced:   time.Now(),
					Status:       mediatypes.SyncStatusSuccess,
				},
			},
		}
		
		// Add items
		for i, item := range clientList.Items {
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
		newList := models.NewMediaItem[T](mediaType, &data)
		newList.Title = clientList.Name
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
func (s *listSyncService[T]) GetSyncStatus(ctx context.Context, listID uint64) (mediatypes.ListSyncStates, error) {
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
	itemList := localList.GetData().GetItemList()
	if itemList.SyncStates == nil {
		return mediatypes.ListSyncStates{}, nil
	}

	return itemList.SyncStates, nil
}

// UpdateSyncStatus updates the sync status for a list
func (s *listSyncService[T]) UpdateSyncStatus(ctx context.Context, userID uint64, listID uint64, clientID uint64, syncState mediatypes.SyncState) error {
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

	// Update sync status
	itemList := localList.GetData().GetItemList()
	if itemList.SyncStates == nil {
		itemList.SyncStates = make(mediatypes.ListSyncStates)
	}
	
	itemList.SyncStates[clientID] = syncState
	
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
		Msg("List sync status updated successfully")

	return nil
}