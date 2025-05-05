package types

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"sort"
	"suasor/utils/logger"
	"time"
)

type ListType string

const (
	ListTypePlaylist   ListType = "playlist"
	ListTypeCollection ListType = "collection"
	ListTypeUnknown    ListType = "unknown"
)

type ListData interface {
	MediaData
	isListData()

	GetTitle() string
	GetItemList() *ItemList

	SetItemList(ItemList)
	AddListItem(ListItem)

	GetMediaType() MediaType
}

type ItemList struct {
	Details   *MediaDetails `json:"details"`
	Items     []ListItem    `json:"items"`
	ItemCount int           `json:"itemCount"`

	OwnerID  uint64 `json:"ownerId"`
	IsPublic bool   `json:"isPublic"`
	// ListCollaboratorIDs
	SharedWith []uint64 `json:"sharedWith"`

	// SyncStates ListSyncStates `json:"syncStates" gorm:"type:jsonb"` // List states for this item (mapping client to their IDs
	LastSynced time.Time `json:"lastSynced"`

	// Track when and which client last modified this playlist
	LastModified time.Time `json:"lastModified"`
	ModifiedBy   uint64    `json:"modifiedBy"` // client ID

	// Smart lists
	IsSmart        bool           `json:"isSmart"`
	SmartCriteria  map[string]any `json:"smartCriteria"`
	AutoUpdateTime time.Time      `json:"autoUpdateTime"`
}

func NewItemList(details *MediaDetails) ItemList {
	return ItemList{
		Details:   details,
		Items:     []ListItem{},
		ItemCount: 0,
	}
}

func NewList[T ListData](details *MediaDetails, itemList ItemList) T {
	var result T

	// Create a concrete type based on T
	switch any(result).(type) {
	case *Playlist:
		playlist := &Playlist{ItemList: itemList}
		playlist.SetDetails(details)
		return any(playlist).(T)
	case *Collection:
		collection := &Collection{ItemList: itemList}
		collection.SetDetails(details)
		return any(collection).(T)
	default:
		// Fallback (should not reach here in practice)
		result.SetDetails(details)
		result.SetItemList(itemList)
		return result
	}
}

// Find an item by ID
func (il *ItemList) FindItemByID(id uint64) (ListItem, int, bool) {
	for i, item := range il.Items {
		if item.ItemID == id {
			return item, i, true
		}
	}
	return ListItem{}, -1, false
}

// Get ordered list of item IDs
func (il *ItemList) GetItemIDs() []uint64 {
	ids := make([]uint64, len(il.Items))
	for i, item := range il.Items {
		ids[i] = item.ItemID
	}
	return ids
}

// Ensure items are ordered by Position
func (il *ItemList) ensureItemsOrdered() {
	sort.Slice(il.Items, func(i, j int) bool {
		return il.Items[i].Position < il.Items[j].Position
	})
}

// Normalize positions to ensure they're sequential from 0
func (il *ItemList) NormalizePositions() {
	il.ensureItemsOrdered()
	for i := range il.Items {
		il.Items[i].Position = i
	}
}

// Add a new item
func (il *ItemList) AddItem(item ListItem) {
	// Set position to end if not specified
	if item.Position < 0 || item.Position > len(il.Items) {
		item.Position = len(il.Items)
	}

	// Shift positions of items that come after
	for i := range il.Items {
		if il.Items[i].Position >= item.Position {
			il.Items[i].Position++
			il.Items[i].LastChanged = time.Now()
		}
	}

	// Add change record if not present
	if len(item.ChangeHistory) == 0 {
		item.AddChangeRecord(0, "add")
	}

	il.Items = append(il.Items, item)
	il.ensureItemsOrdered()
	il.ItemCount = len(il.Items)
	il.LastModified = time.Now()
	il.ModifiedBy = 0
}

func (il *ItemList) AddItemWithClientID(item ListItem, clientID uint64) {
	// Set position to end if not specified
	if item.Position < 0 || item.Position > len(il.Items) {
		item.Position = len(il.Items)
	}

	// Shift positions of items that come after
	for i := range il.Items {
		if il.Items[i].Position >= item.Position {
			il.Items[i].Position++
			il.Items[i].LastChanged = time.Now()
		}
	}

	// Add change record if not present
	if len(item.ChangeHistory) == 0 {
		item.AddChangeRecord(0, "add")
	}

	il.Items = append(il.Items, item)
	il.ensureItemsOrdered()
	il.ItemCount = len(il.Items)
	il.LastModified = time.Now()
	il.ModifiedBy = clientID
}

// Remove an item by ID
func (il *ItemList) RemoveItem(itemID uint64, clientID uint64) error {
	for i, item := range il.Items {
		if item.ItemID == itemID {
			// Remove the item
			// il.Items = append(il.Items[:i], il.Items[i+1:]...)
			il.Items = slices.Delete(il.Items, i, i+1)

			// Update positions for items after the removed one
			for j := range il.Items {
				if il.Items[j].Position > item.Position {
					il.Items[j].Position--
					il.Items[j].AddChangeRecord(clientID, "reposition")
				}
			}

			il.ItemCount = len(il.Items)
			il.LastModified = time.Now()
			il.ModifiedBy = clientID
			return nil
		}
	}
	return fmt.Errorf("item %d not found", itemID)
}

// ValidateItems checks for integrity issues
func (il *ItemList) ValidateItems() []string {
	issues := []string{}

	// Check position uniqueness
	posMap := make(map[int]bool)
	for _, item := range il.Items {
		if posMap[item.Position] {
			issues = append(issues, fmt.Sprintf("duplicate position: %d", item.Position))
		}
		posMap[item.Position] = true
	}

	// Check for position gaps
	for i := range il.Items {
		if !posMap[i] {
			issues = append(issues, fmt.Sprintf("missing position: %d", i))
		}
	}

	// Check ItemCount matches actual count
	if il.ItemCount != len(il.Items) {
		issues = append(issues, "ItemCount doesn't match actual item count")
	}

	return issues
}

// SearchItems filters items based on a predicate
// GetPage returns a paginated view of the items
func (il *ItemList) GetPage(page, pageSize int) []ListItem {
	if page < 0 || pageSize <= 0 {
		return []ListItem{}
	}

	start := page * pageSize
	if start >= len(il.Items) {
		return []ListItem{}
	}

	end := min(start+pageSize, len(il.Items))

	return il.Items[start:end]
}

// Apply changes from client with proper ID mapping
// func (il *ItemList) ApplyClientChanges(clientID uint64, clientItems ClientListItems,
// 	mappingService IDMappingService, serviceType string) error {
// 	now := time.Now()
// 	internalItems := make([]ListItem, 0, len(clientItems))
//
// 	clientListID := "" // Get this from somewhere or add a parameter
// 	existingState := il.SyncStates.GetListSyncState(clientID)
// 	if existingState != nil {
// 		clientListID = existingState.ClientListID
// 	}
//
// 	// Map external IDs to internal ones
// 	for _, item := range clientItems {
// 		internalID, err := mappingService.ExternalToInternal(item.ItemID, serviceType)
// 		if err != nil {
// 			// log.Printf("Failed to map external ID %s: %v", item.ItemID, err)
// 			continue
// 		}
//
// 		internalItems = append(internalItems, ListItem{
// 			ItemID:        internalID,
// 			Position:      item.Position,
// 			LastChanged:   item.LastChanged,
// 			ChangeHistory: item.ChangeHistory,
// 		})
// 	}
//
// 	// Create lookup map for internal items
// 	existingItemsMap := make(map[uint64]int) // Maps ID to index
// 	for i, item := range il.Items {
// 		existingItemsMap[item.ItemID] = i
// 	}
//
// 	// Apply changes
// 	for _, newItem := range internalItems {
// 		if existingIndex, exists := existingItemsMap[newItem.ItemID]; exists {
// 			// Item exists - update if newer
// 			if newItem.LastChanged.After(il.Items[existingIndex].LastChanged) {
// 				il.Items[existingIndex].Position = newItem.Position
// 				il.Items[existingIndex].LastChanged = now
// 				il.Items[existingIndex].ChangeHistory = append(
// 					il.Items[existingIndex].ChangeHistory,
// 					ChangeRecord{
// 						ClientID:   clientID,
// 						ItemID:     fmt.Sprintf("%d", newItem.ItemID),
// 						ChangeType: "update",
// 						Timestamp:  now,
// 					})
// 			}
// 		} else {
// 			// New item - add it
// 			il.AddItemWithClientID(newItem, clientID)
// 		}
// 	}
//
// 	// Handle removals (items in playlist but not in client items)
// 	if len(internalItems) > 0 {
// 		clientItemsMap := make(map[uint64]bool)
// 		for _, item := range internalItems {
// 			clientItemsMap[item.ItemID] = true
// 		}
//
// 		// Find items to remove
// 		var itemsToRemove []int
// 		for i, item := range il.Items {
// 			if !clientItemsMap[item.ItemID] {
// 				itemsToRemove = append(itemsToRemove, i)
// 			}
// 		}
//
// 		// Remove items (in reverse order to maintain indices)
// 		for i := len(itemsToRemove) - 1; i >= 0; i-- {
// 			idx := itemsToRemove[i]
// 			// il.Items = append(il.Items[:idx], il.Items[idx+1:]...)
// 			il.Items = slices.Delete(il.Items, idx, idx+1)
// 		}
// 	}
//
// 	il.NormalizePositions()
// 	il.updateClientState(clientID, clientItems, clientListID)
// 	il.ItemCount = len(il.Items)
// 	il.LastModified = now
// 	il.ModifiedBy = clientID
//
// 	return nil
// }

// SynchronizeWithClient compares local state with a client state and resolves differences
// func (il *ItemList) SynchronizeListWithClient(clientID uint64, mappingService IDMappingService, serviceType string) (bool, error) {
// 	state := il.SyncStates.GetListSyncState(clientID)
// 	if state == nil {
// 		return false, fmt.Errorf("no state exists for client %d", clientID)
// 	}
//
// 	// Generate new sync payload
// 	currentItems, err := il.GenerateSyncPayload(clientID, mappingService, serviceType)
// 	if err != nil {
// 		return false, err
// 	}
//
// 	// Compare with existing client state
// 	changes := false
// 	localItemMap := make(map[string]ClientListItem)
// 	for _, item := range currentItems {
// 		localItemMap[item.ItemID] = item
// 	}
//
// 	// Check for additions or changes
// 	for _, clientItem := range state.Items {
// 		localItem, exists := localItemMap[clientItem.ItemID]
// 		if !exists || clientItem.Position != localItem.Position {
// 			changes = true
// 			break
// 		}
// 	}
//
// 	// Check for removals
// 	if len(currentItems) != len(state.Items) {
// 		changes = true
// 	}
//
// 	if changes {
// 		state.Items = currentItems
// 		state.LastSynced = time.Now()
// 	}
//
// 	return changes, nil
// }

// SynchronizeWithClient compares local state with a client state and resolves differences
// func (mi *ItemList) SynchronizeWithClient(clientID uint64, mappingService IDMappingService, serviceType string) (bool, error) {
// 	state := mi.SyncStates.GetListSyncState(clientID)
// 	if state == nil {
// 		return false, fmt.Errorf("no state exists for client %d", clientID)
// 	}
//
// 	// Generate new sync payload
// 	currentItems, err := mi.GenerateSyncPayload(clientID, mappingService, serviceType)
// 	if err != nil {
// 		return false, err
// 	}
//
// 	// Compare with existing client state
// 	changes := false
// 	localItemMap := make(map[string]ClientListItem)
// 	for _, item := range currentItems {
// 		localItemMap[item.ItemID] = item
// 	}
//
// 	// Check for additions or changes
// 	for _, clientItem := range state.Items {
// 		localItem, exists := localItemMap[clientItem.ItemID]
// 		if !exists || clientItem.Position != localItem.Position {
// 			changes = true
// 			break
// 		}
// 	}
//
// 	// Check for removals
// 	if len(currentItems) != len(state.Items) {
// 		changes = true
// 	}
//
// 	if changes {
// 		state.Items = currentItems
// 		state.LastSynced = time.Now()
// 	}
//
// 	return changes, nil
// }

// ApplyChangesFromMultipleClients safely applies changes from multiple clients
// func (mi *ItemList) ApplyChangesFromMultipleClients(
// 	clientChanges map[uint64]ClientListItems,
// 	mappingService IDMappingService,
// 	serviceType string) error {
//
// 	// First collect and convert all changes
// 	allChanges := make(map[uint64][]ListItem)
// 	now := time.Now()
//
// 	for cID, items := range clientChanges {
// 		internalItems := make([]ListItem, 0, len(items))
//
// 		for _, item := range items {
// 			internaID, err := mappingService.ExternalToInternal(item.ItemID, serviceType)
// 			if err != nil {
// 				continue
// 			}
// 			internalItems = append(internalItems, ListItem{
// 				ItemID:      internaID,
// 				Position:    item.Position,
// 				LastChanged: item.LastChanged,
// 				// ChangeHistory: item.ChangeHistory,
// 			})
// 		}
//
// 		allChanges[cID] = internalItems
// 	}
//
// 	// Create a map to track the latest version of each item
// 	latestItems := make(map[uint64]ListItem)
// 	latestTimestamps := make(map[uint64]time.Time)
//
// 	// Find the latest version of each item across all clients
// 	for _, items := range allChanges {
// 		for _, item := range items {
// 			currentLatestTime, exists := latestTimestamps[item.ItemID]
// 			if !exists || item.LastChanged.After(currentLatestTime) {
// 				latestItems[item.ItemID] = item
// 				latestTimestamps[item.ItemID] = item.LastChanged
// 			}
// 		}
// 	}
//
// 	// Apply the latest version of each item
// 	for itemID, item := range latestItems {
// 		existingItem, idx, found := mi.FindItemByID(itemID)
// 		if found {
// 			if item.LastChanged.After(existingItem.LastChanged) {
// 				// Update existing item
// 				mi.Items[idx].Position = item.Position
// 				mi.Items[idx].LastChanged = now
// 				mi.Items[idx].ChangeHistory = append(
// 					mi.Items[idx].ChangeHistory,
// 					ChangeRecord{
// 						ClientID:   mi.ModifiedBy,
// 						ItemID:     fmt.Sprintf("%d", itemID),
// 						ChangeType: "multi-client-update",
// 						Timestamp:  now,
// 					})
// 			}
// 		} else {
// 			// Add new item
// 			mi.AddItemWithClientID(item, mi.ModifiedBy)
// 		}
// 	}
//
// 	mi.NormalizePositions()
// 	mi.LastModified = now
// 	mi.ItemCount = len(mi.Items)
//
// 	return nil
// }

// Generate sync payload with proper ID mapping
// func (mi *ItemList) GenerateSyncPayload(clientID uint64,
// 	mappingService IDMappingService,
// 	serviceType string) (ClientListItems, error) {
// 	payload := make(ClientListItems, 0, len(mi.Items))
//
// 	for _, item := range mi.Items {
// 		// Map internal ID to external ID for this service
// 		externalID, err := mappingService.InternalToExternal(item.ItemID, serviceType)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to map item %d: %w", item.ItemID, err)
// 		}
//
// 		payload = append(payload, ClientListItem{
// 			ItemID:        externalID,
// 			Position:      item.Position,
// 			LastChanged:   item.LastChanged,
// 			ChangeHistory: item.ChangeHistory,
// 		})
// 	}
//
// 	return payload, nil
// }

// Update client state
// func (mi *ItemList) updateClientState(clientID uint64, clientItems ClientListItems, clientListID string) {
// 	now := time.Now()
// 	state := mi.SyncStates.GetListSyncState(clientID)
//
// 	if clientListID == "" {
// 		if state := mi.SyncStates.GetListSyncState(clientID); state != nil {
// 			clientListID = state.ClientListID
// 		}
// 	}
//
// 	if state != nil {
// 		// Update existing state
// 		state.Items = clientItems
// 		state.LastSynced = now
// 		state.ClientListID = clientListID
// 	} else {
// 		// Add new state
// 		mi.SyncStates = append(mi.SyncStates, ListSyncState{
// 			ClientID:     clientID,
// 			Items:        clientItems,
// 			ClientListID: clientListID,
// 			LastSynced:   now,
// 		})
// 	}
//
// 	mi.LastSynced = now
//
// }
//
// // Interface for ID mapping service
// type IDMappingService interface {
// 	ExternalToInternal(externalID string, serviceType string) (uint64, error)
// 	InternalToExternal(internalID uint64, serviceType string) (string, error)
// }

// Scan
func (m *ItemList) Scan(value any) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, m)
}

func (m *ItemList) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

func (m *ItemList) Merge(other *ItemList) {
	m.Details.Merge(other.Details)
	// TODO: Additional merge logic for fields

}

// RemoveItemAtPosition removes an item from a list at a specific position
func (il *ItemList) RemoveItemAtPosition(itemID uint64, position int, clientID uint64) error {
	log := logger.LoggerFromContext(context.Background())
	log.Debug().
		Uint64("itemID", itemID).
		Int("position", position).
		Msg("Removing item from list")

	// Find item to remove
	var item *ListItem
	for i, existingItem := range il.Items {
		if existingItem.ItemID == itemID && existingItem.Position == position {
			item = &il.Items[i]
			break
		}
	}

	if item == nil {
		log.Warn().
			Uint64("itemID", itemID).
			Msg("Item not found in list at specified position")
		return fmt.Errorf("item not found")
	}

	// Remove item
	il.Items = slices.Delete(il.Items, position, position+1)

	// Update client state
	if item.ChangeHistory == nil {
		item.ChangeHistory = make([]ChangeRecord, 0)
	}

	item.ChangeHistory = append(item.ChangeHistory, ChangeRecord{
		ClientID:   clientID,
		ItemID:     fmt.Sprintf("%d", itemID),
		ChangeType: "remove",
		Timestamp:  time.Now(),
	})

	// Normalize positions
	il.normalizePositions()

	log.Info().
		Uint64("itemID", itemID).
		Int("position", position).
		Msg("Item removed from list successfully")

	return nil
}

func (il *ItemList) normalizePositions() {
	for i := range il.Items {
		il.Items[i].Position = i
	}
}
