package types

import (
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"sort"
	"time"
)

type ListType string

const (
	ListTypePlaylist   ListType = "playlist"
	ListTypeCollection ListType = "collection"
)

// ChangeRecord tracks when and where an item was changed
type ChangeRecord struct {
	ClientID   uint64    `json:"clientID"` // 0 = internal client
	ItemID     string    `json:"itemID,omitempty"`
	ChangeType string    `json:"changeType"` // "add", "remove", "update", "reorder", "sync"
	Timestamp  time.Time `json:"timestamp"`
}

// For internal use
type ListItem struct {
	ItemID        uint64         `json:"itemID"`
	Type          MediaType      `json:"type"`
	Position      int            `json:"position"`
	LastChanged   time.Time      `json:"lastChanged"`
	ChangeHistory []ChangeRecord `json:"changeHistory,omitempty"`
}

// For external sync
type SyncListItem struct {
	ItemID        string         `json:"itemID"`
	Position      int            `json:"position"`
	LastChanged   time.Time      `json:"lastChanged"`
	ChangeHistory []ChangeRecord `json:"changeHistory,omitempty"`
}

type SyncListItems []SyncListItem
type ListItems []ListItem

func (item *ListItem) AddChangeRecord(clientID uint64, changeType string) {
	item.ChangeHistory = append(item.ChangeHistory, ChangeRecord{
		ClientID:   clientID,
		ItemID:     fmt.Sprintf("%d", item.ItemID), // Convert to string for consistency
		ChangeType: changeType,
		Timestamp:  time.Now(),
	})
	item.LastChanged = time.Now()
}

func (item *SyncListItem) ToSyncItems(idMapper func(string) string) SyncListItems {
	return SyncListItems{
		SyncListItem{
			ItemID:        idMapper(item.ItemID),
			Position:      item.Position,
			LastChanged:   item.LastChanged,
			ChangeHistory: item.ChangeHistory,
		},
	}
}
func (items SyncListItems) ToListItems(idMapper func(string) string) SyncListItems {
	result := make(SyncListItems, len(items))
	for i, item := range items {
		result[i] = SyncListItem{
			ItemID:        idMapper(item.ItemID),
			Position:      item.Position,
			LastChanged:   item.LastChanged,
			ChangeHistory: item.ChangeHistory,
		}
	}
	return result
}

func (items ListItems) ToSyncItems(idMapper func(uint64) string) SyncListItems {
	result := make(SyncListItems, len(items))
	for i, item := range items {
		result[i] = SyncListItem{
			ItemID:        idMapper(item.ItemID),
			Position:      item.Position,
			LastChanged:   item.LastChanged,
			ChangeHistory: item.ChangeHistory,
		}
	}
	return result
}

// SyncClientState represents the state of a collection or playlist on a particular client
type SyncClientState struct {
	ClientID     uint64 `json:"clientID"`
	ClientListID string `json:"clientListID,omitempty"`

	// Integration Client's Internal IDs for the items
	Items SyncListItems `json:"items"`

	// Time last synced to this client
	LastSynced time.Time `json:"lastSynced,omitempty"`
}

type SyncClientStates []SyncClientState

// Add validation method
func (state SyncClientState) ValidateItemOrdering() bool {
	// Check that positions match array indices
	for i, item := range state.Items {
		if item.Position != i {
			return false
		}
	}
	return true
}

// GetSyncClientState returns the sync state for a specific client
func (states SyncClientStates) GetSyncClientState(clientID uint64) *SyncClientState {
	for i, state := range states {
		if state.ClientID == clientID {
			return &states[i]
		}
	}
	return nil
}

// Find a SyncClientState by its ClientListID
func (states SyncClientStates) FindByClientListID(clientListID string) *SyncClientState {
	for i, state := range states {
		if state.ClientListID == clientListID {
			return &states[i]
		}
	}
	return nil
}

// MergeItemsIntoSyncState merges new items with existing ones in a sync state
func (states *SyncClientStates) MergeItemsIntoSyncState(clientID uint64, newItems SyncListItems, clientListID string) {
	now := time.Now()
	state := states.GetSyncClientState(clientID)

	// Update timestamps and add sync records for all incoming items
	for i := range newItems {
		if newItems[i].LastChanged.IsZero() {
			newItems[i].LastChanged = now
		}

		// Add sync record if not present
		newItems[i].ChangeHistory = append(newItems[i].ChangeHistory, ChangeRecord{
			ClientID:   clientID,
			ItemID:     newItems[i].ItemID,
			ChangeType: "sync",
			Timestamp:  now,
		})
	}

	if state == nil {
		// No existing state, just add a new one with all items
		*states = append(*states, SyncClientState{
			ClientID:     clientID,
			Items:        newItems,
			ClientListID: clientListID,
			LastSynced:   now,
		})
		return
	}

	// Create a map of existing items by ID for quick lookup
	existingItemsMap := make(map[string]int) // Maps ID to index
	for i, item := range state.Items {
		existingItemsMap[item.ItemID] = i
	}

	// Process each new item
	for _, newItem := range newItems {
		if existingIndex, exists := existingItemsMap[newItem.ItemID]; exists {
			// Item exists - update if newer
			if newItem.LastChanged.After(state.Items[existingIndex].LastChanged) {
				state.Items[existingIndex] = newItem
			}
		} else {
			// New item - add it
			state.Items = append(state.Items, newItem)
		}
	}

	// Update state metadata
	state.LastSynced = now
	state.ClientListID = clientListID
}

// ItemList is the common base for collections and playlists
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

// Interface for ID mapping service
type IDMappingService interface {
	ExternalToInternal(externalID string, serviceType string) (uint64, error)
	InternalToExternal(internalID uint64, serviceType string) (string, error)
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
			log.Printf("Failed to map external ID %s: %v", item.ItemID, err)
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

// Collection represents an unordered set of unique media items
type Collection struct {
	ItemList
	itemMap map[uint64]int // Maps IDs to indices for fast lookup
}

// Initialize the collection's item map
func (c *Collection) initItemMap() {
	if c.itemMap == nil {
		c.itemMap = make(map[uint64]int, len(c.Items))
		for i, item := range c.Items {
			c.itemMap[item.ItemID] = i
		}
	}
}

// FindItemByID optimized for Collections using the item map
func (c *Collection) FindItemByID(id uint64) (ListItem, int, bool) {
	c.initItemMap()
	if idx, exists := c.itemMap[id]; exists && idx < len(c.Items) && c.Items[idx].ItemID == id {
		return c.Items[idx], idx, true
	}
	// Fallback to linear search in case map is out of sync
	return c.ItemList.FindItemByID(id)
}

// AddItem overriden to ensure uniqueness in collections
func (c *Collection) AddItem(item ListItem, clientID uint64) {
	c.initItemMap()

	// Check for existing item first - ignore duplicates
	if _, exists := c.itemMap[item.ItemID]; exists {
		return
	}

	// Call parent implementation
	c.ItemList.AddItemWithClientID(item, clientID)

	// Update the map
	c.itemMap[item.ItemID] = len(c.Items) - 1
}

// RemoveItem overriden to update the map
// func (c *Collection) RemoveItem(itemID uint64, clientID uint64) error {
// 	err := c.ItemList.RemoveItem(itemID, clientID)
// 	if err == nil {
// 		// Item was removed, update the map
// 		delete(c.itemMap, itemID)
//
// 		// Rebuild the map if needed since indices have changed
// 		c.itemMap = make(map[uint64]int, len(c.Items))
// 		for i, item := range c.Items {
// 			c.itemMap[item.ItemID] = i
// 		}
// 	}
// 	return err
// }

// EnsureNoDuplicates is now much more efficient with the item map
func (c *Collection) EnsureNoDuplicates() {
	c.initItemMap()
	seen := make(map[uint64]bool)
	uniqueItems := make([]ListItem, 0, len(c.Items))

	for _, item := range c.Items {
		if !seen[item.ItemID] {
			seen[item.ItemID] = true
			uniqueItems = append(uniqueItems, item)
		}
	}

	c.Items = uniqueItems
	c.ItemCount = len(c.Items)

	// Rebuild the map
	c.itemMap = make(map[uint64]int, len(c.Items))
	for i, item := range c.Items {
		c.itemMap[item.ItemID] = i
	}
}

// Playlist represents an ordered collection that can contain duplicates
type Playlist struct {
	ItemList
}

// Scan implements the sql.Scanner interface
func (p *Playlist) Scan(value any) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, got %T", value)
	}

	// Create a new playlist to unmarshal into
	var playlist Playlist
	if err := json.Unmarshal(bytes, &playlist); err != nil {
		return fmt.Errorf("failed to unmarshal Playlist data: %w", err)
	}

	// Copy the unmarshaled data to the receiver
	*p = playlist
	return nil
}

// DetectItemOrderConflicts finds conflicts in item ordering
func (p *Playlist) DetectItemOrderConflicts(clientState SyncClientState,
	mappingService IDMappingService,
	serviceType string) []string {
	conflicts := []string{}

	// Create position map for local items
	localPositions := make(map[uint64]int)
	for _, item := range p.Items {
		localPositions[item.ItemID] = item.Position
	}

	// Check each client item
	for _, clientItem := range clientState.Items {
		internalID, err := mappingService.ExternalToInternal(clientItem.ItemID, serviceType)
		if err != nil {
			continue // Skip items we can't map
		}

		// If item exists locally with different position, it's a conflict
		if localPos, exists := localPositions[internalID]; exists {
			if localPos != clientItem.Position {
				conflicts = append(conflicts, clientItem.ItemID)
			}
		}
	}

	return conflicts
}

// ReorderItem changes an item's position within the playlist
func (p *Playlist) ReorderItem(itemID uint64, newPosition int, clientID uint64) error {
	// Find the item
	item, oldIndex, found := p.FindItemByID(itemID)
	if !found {
		return fmt.Errorf("item %d not found", itemID)
	}

	oldPosition := item.Position
	if oldPosition == newPosition {
		return nil // No change needed
	}

	// Adjust positions of other items
	for i := range p.Items {
		if oldPosition < newPosition {
			// Moving forward: decrement items in between
			if p.Items[i].Position > oldPosition && p.Items[i].Position <= newPosition {
				p.Items[i].Position--
				p.Items[i].AddChangeRecord(clientID, "reorder")
			}
		} else {
			// Moving backward: increment items in between
			if p.Items[i].Position >= newPosition && p.Items[i].Position < oldPosition {
				p.Items[i].Position++
				p.Items[i].AddChangeRecord(clientID, "reorder")
			}
		}
	}

	// Update the moved item's position
	p.Items[oldIndex].Position = newPosition
	p.Items[oldIndex].AddChangeRecord(clientID, "reorder")

	p.ensureItemsOrdered()
	p.LastModified = time.Now()
	p.ModifiedBy = clientID
	return nil
}

// Create transformation methods between ID types
func ExternalToInternalID(externalID string, mapping map[string]uint64) (uint64, error) {
	id, exists := mapping[externalID]
	if !exists {
		return 0, fmt.Errorf("no mapping for external ID: %s", externalID)
	}
	return id, nil
}

// MergePlaylists combines two playlists with conflict resolution
func MergePlaylists(primary, secondary *Playlist) *Playlist {
	result := *primary // Make a copy

	// Map of items in primary by ID
	primaryItems := make(map[uint64]bool)
	for _, item := range primary.Items {
		primaryItems[item.ItemID] = true
	}

	// Add items from secondary that aren't in primary
	for _, item := range secondary.Items {
		if !primaryItems[item.ItemID] {
			result.Items = append(result.Items, item)
		}
	}

	result.NormalizePositions()
	result.ItemCount = len(result.Items)
	result.LastModified = time.Now()

	return &result
}

func NewList[T ListData](details MediaDetails, itemList ItemList) T {
	var zero T
	zero.SetDetails(details)
	zero.SetItemList(itemList)
	return zero
}
