package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// Playlist represents an ordered collection that can contain duplicates
type Playlist struct {
	ItemList
}

// DetectItemOrderConflicts finds conflicts in item ordering
func (p *Playlist) DetectItemOrderConflicts(clientState ListSyncState,
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
func (p *Playlist) isMediaData()             {}
func (p *Playlist) isListData()              {}
func (p *Playlist) GetDetails() MediaDetails { return p.ItemList.Details }
func (p *Playlist) SetDetails(details MediaDetails) {
	p.ItemList.Details = details
}
func (p *Playlist) GetItemList() ItemList         { return p.ItemList }
func (p *Playlist) SetItemList(itemList ItemList) { p.ItemList = itemList }
func (p *Playlist) GetMediaType() MediaType       { return MediaTypePlaylist }

func (p *Playlist) AddListItem(item ListItem) {
	p.ItemList.AddItem(item)
}

func (p *Playlist) AddListItemWithClientID(item ListItem, clientID uint64) {
	p.ItemList.AddItemWithClientID(item, clientID)
}

func (p *Playlist) RemoveItem(itemID uint64, clientID uint64) error {
	return p.ItemList.RemoveItem(itemID, clientID)
}

// GetTitle
func (p *Playlist) GetTitle() string { return p.Details.Title }

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

// Value
func (p *Playlist) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}
	return json.Marshal(p)
}

func ExternalToInternalID(externalID string, mapping map[string]uint64) (uint64, error) {
	id, exists := mapping[externalID]
	if !exists {
		return 0, fmt.Errorf("no mapping for external ID: %s", externalID)
	}
	return id, nil
}

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
