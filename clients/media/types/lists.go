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

	OwnerID        uint64 `json:"ownerId"`
	OriginClientID uint64 `json:"originClientId"` // 0 for internal db, otherwise external clientID
	IsPublic       bool   `json:"isPublic"`
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

func (il *ItemList) ForEach(callback func(item ListItem) bool) {
	for _, item := range il.Items {
		if !callback(item) {
			return
		}
	}
}
