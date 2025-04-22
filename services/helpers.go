package services

import (
	mediatypes "suasor/clients/media/types"
	"suasor/types/models"
)

// createStubItems creates stub media items for all items in a collection
func createStubItems(listItems []mediatypes.ListItem) []*models.MediaItem[mediatypes.MediaData] {
	var items []*models.MediaItem[mediatypes.MediaData]
	for _, listItem := range listItems {
		items = append(items, createStubItem(listItem))
	}
	return items
}

// createStubItem creates a stub media item for a list item
func createStubItem(listItem mediatypes.ListItem) *models.MediaItem[mediatypes.MediaData] {
	return &models.MediaItem[mediatypes.MediaData]{
		BaseModel: models.BaseModel{
			ID:        listItem.ItemID,
			CreatedAt: listItem.LastChanged,
			UpdatedAt: listItem.LastChanged,
		},
		Type:        mediatypes.MediaTypeTrack, // Placeholder, would get actual type from DB
		ReleaseYear: 0,                         // Would be populated from DB
	}
}
