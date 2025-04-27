package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"sort"
	"time"
)

type ListType string

const (
	ListTypePlaylist   ListType = "playlist"
	ListTypeCollection ListType = "collection"
)

type ListData interface {
	MediaData
	isListData()

	GetTitle() string
	GetDetails() MediaDetails
	GetItemList() ItemList

	SetItemList(ItemList)
	SetDetails(MediaDetails)
	AddListItem(ListItem)

	GetMediaType() MediaType
}

type ItemList struct {
	Details MediaDetails `json:"details"`
	Items   []ListItem   `json:"items"`

	SyncClientStates SyncClientStates `json:"syncClientStates"`
	ItemCount        int              `json:"itemCount"`
	OwnerID          uint64           `json:"owner"`
	// ListCollaboratorIDs
	SharedWith []int64 `json:"sharedWith"`

	IsPublic   bool      `json:"isPublic"`
	LastSynced time.Time `json:"lastSynced"`

	// Track when and which client last modified this playlist
	LastModified time.Time `json:"lastModified"`
	ModifiedBy   uint64    `json:"modifiedBy"` // client ID

	// Smart lists
	IsSmart        bool           `json:"isSmart"`
	SmartCriteria  map[string]any `json:"smartCriteria"`
	AutoUpdateTime time.Time      `json:"autoUpdateTime"`
}

func NewList[T ListData](details MediaDetails, itemList ItemList) T {
	var zero T
	zero.SetDetails(details)
	zero.SetItemList(itemList)
	return zero
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

// Update client state
func (il *ItemList) updateClientState(clientID uint64, clientItems SyncListItems, clientListID string) {
	now := time.Now()
	state := il.SyncClientStates.GetSyncClientState(clientID)

	if clientListID == "" {
		if state := il.SyncClientStates.GetSyncClientState(clientID); state != nil {
			clientListID = state.ClientListID
		}
	}

	if state != nil {
		// Update existing state
		state.Items = clientItems
		state.LastSynced = now
		state.ClientListID = clientListID
	} else {
		// Add new state
		il.SyncClientStates = append(il.SyncClientStates, SyncClientState{
			ClientID:     clientID,
			Items:        clientItems,
			ClientListID: clientListID,
			LastSynced:   now,
		})
	}

	il.LastSynced = now
}

// Generate sync payload with proper ID mapping
func (il *ItemList) GenerateSyncPayload(clientID uint64,
	mappingService IDMappingService,
	serviceType string) (SyncListItems, error) {
	payload := make(SyncListItems, 0, len(il.Items))

	for _, item := range il.Items {
		// Map internal ID to external ID for this service
		externalID, err := mappingService.InternalToExternal(item.ItemID, serviceType)
		if err != nil {
			return nil, fmt.Errorf("failed to map item %d: %w", item.ItemID, err)
		}

		payload = append(payload, SyncListItem{
			ItemID:        externalID,
			Position:      item.Position,
			LastChanged:   item.LastChanged,
			ChangeHistory: item.ChangeHistory,
		})
	}

	return payload, nil
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

// Apply changes from client with proper ID mapping
func (il *ItemList) ApplyClientChanges(clientID uint64, clientItems SyncListItems,
	mappingService IDMappingService, serviceType string) error {
	now := time.Now()
	internalItems := make([]ListItem, 0, len(clientItems))

	clientListID := "" // Get this from somewhere or add a parameter
	existingState := il.SyncClientStates.GetSyncClientState(clientID)
	if existingState != nil {
		clientListID = existingState.ClientListID
	}

	// Map external IDs to internal ones
	for _, item := range clientItems {
		internalID, err := mappingService.ExternalToInternal(item.ItemID, serviceType)
		if err != nil {
			// log.Printf("Failed to map external ID %s: %v", item.ItemID, err)
			continue
		}

		internalItems = append(internalItems, ListItem{
			ItemID:        internalID,
			Position:      item.Position,
			LastChanged:   item.LastChanged,
			ChangeHistory: item.ChangeHistory,
		})
	}

	// Create lookup map for internal items
	existingItemsMap := make(map[uint64]int) // Maps ID to index
	for i, item := range il.Items {
		existingItemsMap[item.ItemID] = i
	}

	// Apply changes
	for _, newItem := range internalItems {
		if existingIndex, exists := existingItemsMap[newItem.ItemID]; exists {
			// Item exists - update if newer
			if newItem.LastChanged.After(il.Items[existingIndex].LastChanged) {
				il.Items[existingIndex].Position = newItem.Position
				il.Items[existingIndex].LastChanged = now
				il.Items[existingIndex].ChangeHistory = append(
					il.Items[existingIndex].ChangeHistory,
					ChangeRecord{
						ClientID:   clientID,
						ItemID:     fmt.Sprintf("%d", newItem.ItemID),
						ChangeType: "update",
						Timestamp:  now,
					})
			}
		} else {
			// New item - add it
			il.AddItemWithClientID(newItem, clientID)
		}
	}

	// Handle removals (items in playlist but not in client items)
	if len(internalItems) > 0 {
		clientItemsMap := make(map[uint64]bool)
		for _, item := range internalItems {
			clientItemsMap[item.ItemID] = true
		}

		// Find items to remove
		var itemsToRemove []int
		for i, item := range il.Items {
			if !clientItemsMap[item.ItemID] {
				itemsToRemove = append(itemsToRemove, i)
			}
		}

		// Remove items (in reverse order to maintain indices)
		for i := len(itemsToRemove) - 1; i >= 0; i-- {
			idx := itemsToRemove[i]
			// il.Items = append(il.Items[:idx], il.Items[idx+1:]...)
			il.Items = slices.Delete(il.Items, idx, idx+1)
		}
	}

	il.NormalizePositions()
	il.updateClientState(clientID, clientItems, clientListID)
	il.ItemCount = len(il.Items)
	il.LastModified = now
	il.ModifiedBy = clientID

	return nil
}

// SynchronizeWithClient compares local state with a client state and resolves differences
func (il *ItemList) SynchronizeWithClient(clientID uint64, mappingService IDMappingService, serviceType string) (bool, error) {
	state := il.SyncClientStates.GetSyncClientState(clientID)
	if state == nil {
		return false, fmt.Errorf("no state exists for client %d", clientID)
	}

	// Generate new sync payload
	currentItems, err := il.GenerateSyncPayload(clientID, mappingService, serviceType)
	if err != nil {
		return false, err
	}

	// Compare with existing client state
	changes := false
	localItemMap := make(map[string]SyncListItem)
	for _, item := range currentItems {
		localItemMap[item.ItemID] = item
	}

	// Check for additions or changes
	for _, clientItem := range state.Items {
		localItem, exists := localItemMap[clientItem.ItemID]
		if !exists || clientItem.Position != localItem.Position {
			changes = true
			break
		}
	}

	// Check for removals
	if len(currentItems) != len(state.Items) {
		changes = true
	}

	if changes {
		state.Items = currentItems
		state.LastSynced = time.Now()
	}

	return changes, nil
}

// ApplyChangesFromMultipleClients safely applies changes from multiple clients
func (il *ItemList) ApplyChangesFromMultipleClients(
	clientChanges map[uint64]SyncListItems,
	mappingService IDMappingService,
	serviceType string) error {

	// First collect and convert all changes
	allChanges := make(map[uint64][]ListItem)
	now := time.Now()

	for cID, items := range clientChanges {
		internalItems := make([]ListItem, 0, len(items))

		for _, item := range items {
			internalID, err := mappingService.ExternalToInternal(item.ItemID, serviceType)
			if err != nil {
				continue
			}
			internalItems = append(internalItems, ListItem{
				ItemID:        internalID,
				Position:      item.Position,
				LastChanged:   item.LastChanged,
				ChangeHistory: item.ChangeHistory,
			})
		}

		allChanges[cID] = internalItems
	}

	// Create a map to track the latest version of each item
	latestItems := make(map[uint64]ListItem)
	latestTimestamps := make(map[uint64]time.Time)

	// Find the latest version of each item across all clients
	for _, items := range allChanges {
		for _, item := range items {
			currentLatestTime, exists := latestTimestamps[item.ItemID]
			if !exists || item.LastChanged.After(currentLatestTime) {
				latestItems[item.ItemID] = item
				latestTimestamps[item.ItemID] = item.LastChanged
			}
		}
	}

	// Apply the latest version of each item
	for itemID, item := range latestItems {
		existingItem, idx, found := il.FindItemByID(itemID)
		if found {
			if item.LastChanged.After(existingItem.LastChanged) {
				// Update existing item
				il.Items[idx].Position = item.Position
				il.Items[idx].LastChanged = now
				il.Items[idx].ChangeHistory = append(
					il.Items[idx].ChangeHistory,
					ChangeRecord{
						ClientID:   il.ModifiedBy,
						ItemID:     fmt.Sprintf("%d", itemID),
						ChangeType: "multi-client-update",
						Timestamp:  now,
					})
			}
		} else {
			// Add new item
			il.AddItemWithClientID(item, il.ModifiedBy)
		}
	}

	il.NormalizePositions()
	il.LastModified = now
	il.ItemCount = len(il.Items)

	return nil
}

// SearchItems filters items based on a predicate
func (il *ItemList) SearchItems(predicate func(ListItem) bool) []ListItem {
	results := make([]ListItem, 0)
	for _, item := range il.Items {
		if predicate(item) {
			results = append(results, item)
		}
	}
	return results
}

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

// Interface for ID mapping service
type IDMappingService interface {
	ExternalToInternal(externalID string, serviceType string) (uint64, error)
	InternalToExternal(internalID uint64, serviceType string) (string, error)
}
