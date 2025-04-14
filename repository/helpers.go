package repository

import (
	"context"
	"suasor/client/media/types"
	"suasor/types/models"

	"gorm.io/gorm"
)

func fetchMediaItemsByType[T types.MediaData](ctx context.Context, db *gorm.DB, ids []uint64, mediaType types.MediaType) ([]*models.MediaItem[T], error) {

	// Fetch items by type
	var items []*models.MediaItem[T]
	if err := db.WithContext(ctx).Where("id IN ? AND type = ?", ids, mediaType).Find(&items).Error; err != nil {
		return nil, err

	}
	return items, nil
}
