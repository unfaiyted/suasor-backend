package types

// MOVES to /home/faiyt/codebase/suasor/backend/types/models/client_states.go
//
// import (
// 	"time"
// )
//
// // SyncClientState represents the state of a collection or playlist on a particular client
// type SyncClientState struct {
// 	ClientID     uint64 `json:"clientID"`
// 	ClientListID string `json:"clientListID,omitempty"`
//
// 	// Integration Client's Internal IDs for the items
// 	Items SyncListItems `json:"items"`
//
// 	// Time last synced to this client
// 	LastSynced time.Time `json:"lastSynced,omitempty"`
// }
// type SyncClientStates []SyncClientState
//
// // Add validation method
// func (state SyncClientState) ValidateItemOrdering() bool {
// 	// Check that positions match array indices
// 	for i, item := range state.Items {
// 		if item.Position != i {
// 			return false
// 		}
// 	}
// 	return true
// }
// func (states SyncClientStates) GetSyncClientState(clientID uint64) *SyncClientState {
// 	for i, state := range states {
// 		if state.ClientID == clientID {
// 			return &states[i]
// 		}
// 	}
// 	return nil
// }
//
// func (states SyncClientStates) FindByClientListID(clientListID string) *SyncClientState {
// 	for i, state := range states {
// 		if state.ClientListID == clientListID {
// 			return &states[i]
// 		}
// 	}
// 	return nil
// }
//
// // MergeItemsIntoSyncState merges new items with existing ones in a sync state
// func (states *SyncClientStates) MergeItemsIntoSyncState(clientID uint64, newItems SyncListItems, clientListID string) {
// 	now := time.Now()
// 	state := states.GetSyncClientState(clientID)
//
// 	// Update timestamps and add sync records for all incoming items
// 	for i := range newItems {
// 		if newItems[i].LastChanged.IsZero() {
// 			newItems[i].LastChanged = now
// 		}
//
// 		// Add sync record if not present
// 		newItems[i].ChangeHistory = append(newItems[i].ChangeHistory, ChangeRecord{
// 			ClientID:   clientID,
// 			ItemID:     newItems[i].ItemID,
// 			ChangeType: "sync",
// 			Timestamp:  now,
// 		})
// 	}
//
// 	if state == nil {
// 		// No existing state, just add a new one with all items
// 		*states = append(*states, SyncClientState{
// 			ClientID:     clientID,
// 			Items:        newItems,
// 			ClientListID: clientListID,
// 			LastSynced:   now,
// 		})
// 		return
// 	}
//
// 	// Create a map of existing items by ID for quick lookup
// 	existingItemsMap := make(map[string]int) // Maps ID to index
// 	for i, item := range state.Items {
// 		existingItemsMap[item.ItemID] = i
// 	}
//
// 	// Process each new item
// 	for _, newItem := range newItems {
// 		if existingIndex, exists := existingItemsMap[newItem.ItemID]; exists {
// 			// Item exists - update if newer
// 			if newItem.LastChanged.After(state.Items[existingIndex].LastChanged) {
// 				state.Items[existingIndex] = newItem
// 			}
// 		} else {
// 			// New item - add it
// 			state.Items = append(state.Items, newItem)
// 		}
// 	}
//
// 	// Update state metadata
// 	state.LastSynced = now
// 	state.ClientListID = clientListID
// }
