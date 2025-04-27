package types

import (
	"fmt"
	"time"
)

// For internal use
type ListItem struct {
	ItemID        uint64         `json:"itemID"`
	Type          MediaType      `json:"type"`
	Position      int            `json:"position"`
	LastChanged   time.Time      `json:"lastChanged"`
	ChangeHistory []ChangeRecord `json:"changeHistory,omitempty"`
}

// For external sync id references
type ClientListItem struct {
	ItemID        string         `json:"itemID"`
	Position      int            `json:"position"`
	LastChanged   time.Time      `json:"lastChanged"`
	ChangeHistory []ChangeRecord `json:"changeHistory,omitempty"`
}

// ChangeRecord tracks when and where an item was changed
type ChangeRecord struct {
	ClientID   uint64    `json:"clientID"` // 0 = internal client
	ItemID     string    `json:"itemID,omitempty"`
	ChangeType string    `json:"changeType"` // "add", "remove", "update", "reorder", "sync"
	Timestamp  time.Time `json:"timestamp"`
}

type ClientListItems []ClientListItem
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
func (items ListItems) ToSyncItems(idMapper func(uint64) string) ClientListItems {
	result := make(ClientListItems, len(items))
	for i, item := range items {
		result[i] = ClientListItem{
			ItemID:        idMapper(item.ItemID),
			Position:      item.Position,
			LastChanged:   item.LastChanged,
			ChangeHistory: item.ChangeHistory,
		}
	}
	return result
}

func (item *ClientListItem) ToSyncItems(idMapper func(string) string) ClientListItems {
	return ClientListItems{
		ClientListItem{
			ItemID:        idMapper(item.ItemID),
			Position:      item.Position,
			LastChanged:   item.LastChanged,
			ChangeHistory: item.ChangeHistory,
		},
	}
}
func (items ClientListItems) ToListItems(idMapper func(string) string) ClientListItems {
	result := make(ClientListItems, len(items))
	for i, item := range items {
		result[i] = ClientListItem{
			ItemID:        idMapper(item.ItemID),
			Position:      item.Position,
			LastChanged:   item.LastChanged,
			ChangeHistory: item.ChangeHistory,
		}
	}
	return result
}
