package services

//
// import (
// 	"context"
// 	"suasor/client/media/types"
// 	"suasor/repository"
// 	"suasor/types/models"
// )
//
// // This service provides operations for getting User-specific media history and
// // user favorites and ratings, also any other user specific data that is stored on the server
// // we are connecting to
// type ClientUserMediaItemDataService[T types.MediaData] interface {
// 	// GetUserPlayHistory retrieves play history for a user with optional filtering
// 	GetUserPlayHistory(ctx context.Context, userID uint64, limit, offset int, mediaType *types.MediaType, completed *bool) (interface{}, error)
//
// 	// GetContinueWatching retrieves items that a user has started but not completed
// 	GetContinueWatching(ctx context.Context, userID uint64, limit int) (interface{}, error)
//
// 	// GetByID retrieves a specific play history entry by ID
// 	GetByID(ctx context.Context, id uint64) (interface{}, error)
//
// 	// GetByMediaItemID retrieves play history for a specific media item
// 	GetByMediaItemID(ctx context.Context, mediaItemID, userID uint64) (interface{}, error)
//
// 	// RecordPlay records a new play event
// 	RecordPlay(ctx context.Context, history *models.UserMediaItemData[T]) (interface{}, error)
//
// 	// ToggleFavorite marks or unmarks a media item as a favorite
// 	ToggleFavorite(ctx context.Context, mediaItemID, userID uint64, favorite bool) error
//
// 	// UpdateRating sets a user's rating for a media item
// 	UpdateRating(ctx context.Context, mediaItemID, userID uint64, rating float32) error
//
// 	// GetFavorites retrieves favorite media items for a user
// 	GetFavorites(ctx context.Context, userID uint64, mediaType *types.MediaType, limit, offset int) (interface{}, error)
//
// 	// Delete removes a specific play history entry
// 	Delete(ctx context.Context, id uint64) error
//
// 	// ClearUserHistory removes all play history for a user
// 	ClearUserHistory(ctx context.Context, userID uint64, mediaType *types.MediaType) error
// }
//
// // clientUserMediaItemDataService[T] implements ClientUserMediaItemDataService
// type clientUserMediaItemDataService[T types.MediaData] struct {
// 	repo repository.UserMediaItemDataRepository[T]
// }
//
// // NewClientUserMediaItemDataService creates a new media play history service
// func NewClientUserMediaItemDataService[T types.MediaData](repo repository.UserMediaItemDataRepository[T]) ClientUserMediaItemDataService[T] {
// 	return &clientUserMediaItemDataService[T]{
// 		repo: repo,
// 	}
// }
//
// // GetUserPlayHistory retrieves play history for a user with optional filtering
// func (s *clientUserMediaItemDataService[T]) GetUserPlayHistory(ctx context.Context, userID uint64, limit, offset int, mediaType *types.MediaType, completed *bool) (interface{}, error) {
// 	// Placeholder - this would be implemented with repository calls
// 	return s.repo.GetUserPlayHistory(ctx, userID, limit, offset, completed)
// }
//
// // GetContinueWatching retrieves items that a user has started but not completed
// func (s *clientUserMediaItemDataService[T]) GetContinueWatching(ctx context.Context, userID uint64, limit int) (interface{}, error) {
// 	// Items that are not completed and have been played recently
// 	return s.repo.GetContinueWatching(ctx, userID, limit)
// }
//
// // GetByID retrieves a specific play history entry by ID
// func (s *clientUserMediaItemDataService[T]) GetByID(ctx context.Context, id uint64) (interface{}, error) {
// 	return s.repo.GetByID(ctx, id)
// }
//
// // GetByMediaItemID retrieves play history for a specific media item
// func (s *clientUserMediaItemDataService[T]) GetByMediaItemID(ctx context.Context, mediaItemID, userID uint64) (interface{}, error) {
// 	return s.repo.GetByID(ctx, mediaItemID)
// }
//
// // RecordPlay records a new play event
// func (s *clientUserMediaItemDataService[T]) RecordPlay(ctx context.Context, history *models.UserMediaItemData[T]) (interface{}, error) {
// 	return s.repo.RecordPlay(ctx, history)
// }
//
// // ToggleFavorite marks or unmarks a media item as a favorite
// func (s *clientUserMediaItemDataService[T]) ToggleFavorite(ctx context.Context, mediaItemID, userID uint64, favorite bool) error {
// 	return s.repo.ToggleFavorite(ctx, mediaItemID, userID, favorite)
// }
//
// // UpdateRating sets a user's rating for a media item
// func (s *clientUserMediaItemDataService[T]) UpdateRating(ctx context.Context, mediaItemID, userID uint64, rating float32) error {
// 	return s.repo.UpdateRating(ctx, mediaItemID, userID, rating)
// }
//
// // GetFavorites retrieves favorite media items for a user
// func (s *clientUserMediaItemDataService[T]) GetFavorites(ctx context.Context, userID uint64, mediaType *types.MediaType, limit, offset int) (interface{}, error) {
// 	return s.repo.GetFavorites(ctx, userID, limit, offset)
// }
//
// // Delete removes a specific play history entry
// func (s *clientUserMediaItemDataService[T]) Delete(ctx context.Context, id uint64) error {
// 	return s.repo.Delete(ctx, id)
// }
//
// // ClearUserHistory removes all play history for a user
// func (s *clientUserMediaItemDataService[T]) ClearUserHistory(ctx context.Context, userID uint64, mediaType *types.MediaType) error {
// 	return s.repo.ClearUserHistory(ctx, userID)
// }
