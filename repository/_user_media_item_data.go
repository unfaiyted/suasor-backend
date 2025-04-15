package repository

// import (
// 	"context"
// 	"fmt"
// 	"suasor/client/media/types"
// 	"suasor/types/models"
// 	"time"
//
// 	"gorm.io/gorm"
// )
//
// // UserMediaItemDataRepository defines the interface for media item data operations
// // This is a facade that combines core, user, and client functionality
// type UserMediaItemDataRepository[T types.MediaData] interface {
// 	// Core operations
// 	Create(ctx context.Context, data models.UserMediaItemData[T]) (*models.UserMediaItemData[T], error)
// 	GetByID(ctx context.Context, id uint64) (*models.UserMediaItemData[T], error)
// 	Update(ctx context.Context, data models.UserMediaItemData[T]) (*models.UserMediaItemData[T], error)
// 	Delete(ctx context.Context, id uint64) error
// 	GetByUserIDAndMediaItemID(ctx context.Context, userID, mediaItemID uint64) (*models.UserMediaItemData[T], error)
// 	HasUserMediaItemData(ctx context.Context, userID, mediaItemID uint64) (bool, error)
//
// 	// User operations
// 	GetUserHistory(ctx context.Context, userID uint64, limit, offset int, mediaType *types.MediaType) ([]*models.UserMediaItemData[T], error)
// 	GetRecentHistory(ctx context.Context, userID uint64, limit int, mediaType *types.MediaType) ([]*models.UserMediaItemData[T], error)
// 	GetUserPlayHistory(ctx context.Context, userID uint64, limit, offset int, mediaType *types.MediaType, completed *bool) ([]*models.UserMediaItemData[T], error)
// 	GetContinueWatching(ctx context.Context, userID uint64, limit int, mediaType *types.MediaType) ([]*models.UserMediaItemData[T], error)
// 	RecordPlay(ctx context.Context, data *models.UserMediaItemData[T]) (*models.UserMediaItemData[T], error)
// 	ToggleFavorite(ctx context.Context, mediaItemID, userID uint64, favorite bool) (*models.UserMediaItemData[T], error)
// 	UpdateRating(ctx context.Context, mediaItemID, userID uint64, rating float32) (*models.UserMediaItemData[T], error)
// 	GetFavorites(ctx context.Context, userID uint64, limit, offset int) ([]*models.UserMediaItemData[T], error)
// 	ClearUserHistory(ctx context.Context, userID uint64) error
//
// 	// Client operations
// 	SyncClientItemData(ctx context.Context, userID uint64, clientID string, items []models.UserMediaItemData[T]) error
// 	GetClientItemData(ctx context.Context, userID uint64, clientID string, since *string) ([]*models.UserMediaItemData[T], error)
// 	GetByClientID(ctx context.Context, userID uint64, clientID string, clientItemID string) (*models.UserMediaItemData[T], error)
// 	RecordClientPlay(ctx context.Context, userID uint64, clientID string, clientItemID string, data *models.UserMediaItemData[T]) (*models.UserMediaItemData[T], error)
// 	GetPlaybackState(ctx context.Context, userID uint64, clientID string, clientItemID string) (*models.UserMediaItemData[T], error)
// 	UpdatePlaybackState(ctx context.Context, userID uint64, clientID string, clientItemID string, position int, duration int, percentage float64) (*models.UserMediaItemData[T], error)
// }
//
// // userMediaItemDataRepository implements UserMediaItemDataRepository
// type userMediaItemDataRepository[T types.MediaData] struct {
// 	coreRepo   CoreUserMediaItemDataRepository[T]
// 	userRepo   UserUserMediaItemDataRepository[T]
// 	clientRepo ClientUserMediaItemDataRepository[T]
// }
//
// // NewUserMediaItemDataRepository creates a new media item data repository facade
// func NewUserMediaItemDataRepository[T types.MediaData](db *gorm.DB) UserMediaItemDataRepository[T] {
// 	return &userMediaItemDataRepository[T]{
// 		coreRepo:   NewCoreUserMediaItemDataRepository[T](db),
// 		userRepo:   NewUserUserMediaItemDataRepository[T](db),
// 		clientRepo: NewClientUserMediaItemDataRepository[T](db),
// 	}
// }
//
// // === Core operations ===
//
// func (r *userMediaItemDataRepository[T]) Create(ctx context.Context, data models.UserMediaItemData[T]) (*models.UserMediaItemData[T], error) {
// 	return r.coreRepo.Create(ctx, data)
// }
//
// func (r *userMediaItemDataRepository[T]) GetByID(ctx context.Context, id uint64) (*models.UserMediaItemData[T], error) {
// 	return r.coreRepo.GetByID(ctx, id)
// }
//
// func (r *userMediaItemDataRepository[T]) Update(ctx context.Context, data models.UserMediaItemData[T]) (*models.UserMediaItemData[T], error) {
// 	return r.coreRepo.Update(ctx, data)
// }
//
// func (r *userMediaItemDataRepository[T]) Delete(ctx context.Context, id uint64) error {
// 	return r.coreRepo.Delete(ctx, id)
// }
//
// func (r *userMediaItemDataRepository[T]) GetByUserIDAndMediaItemID(ctx context.Context, userID, mediaItemID uint64) (*models.UserMediaItemData[T], error) {
// 	return r.coreRepo.GetByUserIDAndMediaItemID(ctx, userID, mediaItemID)
// }
//
// func (r *userMediaItemDataRepository[T]) HasUserMediaItemData(ctx context.Context, userID, mediaItemID uint64) (bool, error) {
// 	return r.coreRepo.HasUserMediaItemData(ctx, userID, mediaItemID)
// }
//
// // === User operations ===
//
// func (r *userMediaItemDataRepository[T]) GetUserHistory(ctx context.Context, userID uint64, limit, offset int, mediaType *types.MediaType) ([]*models.UserMediaItemData[T], error) {
// 	return r.userRepo.GetUserHistory(ctx, userID, limit, offset, mediaType)
// }
//
// func (r *userMediaItemDataRepository[T]) GetRecentHistory(ctx context.Context, userID uint64, limit int, mediaType *types.MediaType) ([]*models.UserMediaItemData[T], error) {
// 	return r.userRepo.GetRecentHistory(ctx, userID, limit, mediaType)
// }
//
// func (r *userMediaItemDataRepository[T]) GetUserPlayHistory(ctx context.Context, userID uint64, limit, offset int, mediaType *types.MediaType, completed *bool) ([]*models.UserMediaItemData[T], error) {
// 	return r.userRepo.GetUserPlayHistory(ctx, userID, limit, offset, mediaType, completed)
// }
//
// func (r *userMediaItemDataRepository[T]) GetContinueWatching(ctx context.Context, userID uint64, limit int, mediaType *types.MediaType) ([]*models.UserMediaItemData[T], error) {
// 	return r.userRepo.GetContinueWatching(ctx, userID, limit, mediaType)
// }
//
// func (r *userMediaItemDataRepository[T]) RecordPlay(ctx context.Context, data *models.UserMediaItemData[T]) (*models.UserMediaItemData[T], error) {
// 	return r.userRepo.RecordPlay(ctx, data)
// }
//
// func (r *userMediaItemDataRepository[T]) ToggleFavorite(ctx context.Context, mediaItemID, userID uint64, favorite bool) (*models.UserMediaItemData[T], error) {
// 	return r.userRepo.ToggleFavorite(ctx, mediaItemID, userID, favorite)
// }
//
// func (r *userMediaItemDataRepository[T]) UpdateRating(ctx context.Context, mediaItemID, userID uint64, rating float32) (*models.UserMediaItemData[T], error) {
// 	return r.userRepo.UpdateRating(ctx, mediaItemID, userID, rating)
// }
//
// func (r *userMediaItemDataRepository[T]) GetFavorites(ctx context.Context, userID uint64, limit, offset int) ([]*models.UserMediaItemData[T], error) {
// 	return r.userRepo.GetFavorites(ctx, userID, limit, offset)
// }
//
// func (r *userMediaItemDataRepository[T]) ClearUserHistory(ctx context.Context, userID uint64) error {
// 	return r.userRepo.ClearUserHistory(ctx, userID)
// }
//
// // === Client operations ===
//
// func (r *userMediaItemDataRepository[T]) SyncClientItemData(ctx context.Context, userID uint64, clientID string, items []models.UserMediaItemData[T]) error {
// 	return r.clientRepo.SyncClientItemData(ctx, userID, clientID, items)
// }
//
// func (r *userMediaItemDataRepository[T]) GetClientItemData(ctx context.Context, userID uint64, clientID string, sinceDateStr *string) ([]*models.UserMediaItemData[T], error) {
// 	// If sinceDateStr is nil or empty, use a default date (e.g., 24 hours ago)
// 	var since string
// 	if sinceDateStr == nil || *sinceDateStr == "" {
// 		// Use default date
// 		since = fmt.Sprintf("%d hours ago", 24)
// 	} else {
// 		since = *sinceDateStr
// 	}
//
// 	// Parse since as a date
// 	sinceTime, err := parseTime(since)
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid since date: %w", err)
// 	}
//
// 	return r.clientRepo.GetClientItemData(ctx, userID, clientID, sinceTime)
// }
//
// func (r *userMediaItemDataRepository[T]) GetByClientID(ctx context.Context, userID uint64, clientID string, clientItemID string) (*models.UserMediaItemData[T], error) {
// 	return r.clientRepo.GetByClientID(ctx, userID, clientID, clientItemID)
// }
//
// func (r *userMediaItemDataRepository[T]) RecordClientPlay(ctx context.Context, userID uint64, clientID string, clientItemID string, data *models.UserMediaItemData[T]) (*models.UserMediaItemData[T], error) {
// 	return r.clientRepo.RecordClientPlay(ctx, userID, clientID, clientItemID, data)
// }
//
// func (r *userMediaItemDataRepository[T]) GetPlaybackState(ctx context.Context, userID uint64, clientID string, clientItemID string) (*models.UserMediaItemData[T], error) {
// 	return r.clientRepo.GetPlaybackState(ctx, userID, clientID, clientItemID)
// }
//
// func (r *userMediaItemDataRepository[T]) UpdatePlaybackState(ctx context.Context, userID uint64, clientID string, clientItemID string, position int, duration int, percentage float64) (*models.UserMediaItemData[T], error) {
// 	return r.clientRepo.UpdatePlaybackState(ctx, userID, clientID, clientItemID, position, duration, percentage)
// }
//
// // Helper function to parse time strings like "24 hours ago", "7 days ago", etc.
// func parseTime(timeStr string) (time.Time, error) {
// 	// For now, just return a simple placeholder
// 	// In a real implementation, this would parse the time string properly
// 	return time.Now().AddDate(0, 0, -1), nil
// }

