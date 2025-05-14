package models

import (
	"suasor/clients/media/types"
)

type MediaItemList[T types.ListData] struct {
	List *MediaItem[T] `json:"list"`

	Items *MediaItemResults `json:"items"`

	ListType           types.ListType `json:"listType"`
	ListOriginClientID uint64         `json:"listOriginClientID"` // 0 for internal db, otherwise external client/ProviderID
	OwnerID            uint64         `json:"ownerID"`
}

func NewMediaItemList[T types.ListData](list *MediaItem[T], listOriginClientID uint64, ownerID uint64) *MediaItemList[T] {
	listType := types.GetListType[T]()

	return &MediaItemList[T]{
		List:               list,
		ListType:           listType,
		ListOriginClientID: listOriginClientID,
		OwnerID:            ownerID,
		Items: &MediaItemResults{
			Order:      ListItems{},
			TotalItems: 0,
		},
	}
}

func (m *MediaItemList[T]) Len() int {
	return m.Items.TotalItems
}

func (m *MediaItemList[T]) GetSyncClientItemIDs(clientID uint64) []string {
	ids := make([]string, 0)
	m.ForEach(func(uuid string, mediaType types.MediaType, item any) bool {
		if mediaItem, ok := item.(*MediaItem[types.MediaData]); ok && mediaItem != nil {
			clientItemID := mediaItem.SyncClients.GetClientItemID(clientID)
			// Only add non-empty client IDs
			if clientItemID != "" {
				ids = append(ids, clientItemID)
			}
		}
		return true
	})

	return ids
}

func (m *MediaItemList[T]) GetListSyncClients() []*SyncClient {
	return m.List.SyncClients.GetSyncClients()
}

func (m *MediaItemList[T]) GetOriginalClientID() uint64 {
	return m.ListOriginClientID
}

func (m *MediaItemList[T]) GetOwnerID() uint64 {
	return m.OwnerID
}

// ForEach iterates over all media items in the list in the specified order.
// The callback function receives the UUID, media type, and the item itself.
// If the callback returns false, iteration stops early.
func (m *MediaItemList[T]) ForEach(callback func(uuid string, mediaType types.MediaType, item any) bool) {
	m.Items.ForEach(callback)
}

// IsItemAtPosition checks if a media item is at a specific position
func (m *MediaItemList[T]) IsItemAtPosition(uuid string, position int) bool {
	return m.Items.IsItemAtPosition(uuid, position)
}
