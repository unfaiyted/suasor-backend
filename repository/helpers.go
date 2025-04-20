package repository

import (
	"context"
	"suasor/client/media/types"
	// clienttypes "suasor/client/types"
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

// func GetClientTypeFromID(ctx context.Context, db *gorm.DB, clientID uint64) (clienttypes.ClientType, error) {
// 	var client clienttypes.ClientType
// 	if err := db.WithContext(ctx).
// 		Where("id = ?", clientID).
// 		Select("config -> 'data' ->> 'type'").
// 		Scan(&client).Error; err != nil {
// 		return "", err
// 	}
// 	return client, nil
// }
