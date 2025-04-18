package services

//
// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"time"
//
// 	mediatypes "suasor/client/media/types"
// 	"suasor/repository"
// 	"suasor/types/models"
// 	"suasor/utils"
// )
//
// // UserCollectionService defines the interface for user-owned collection operations
// // This service extends CoreCollectionService with operations specific to user-owned collections
// type UserCollectionService interface {
// 	// Include all core collection service methods
// 	CoreListService[*mediatypes.Collection]
//
// 	// User-specific operations
// 	GetByUserID(ctx context.Context, userID uint64, limit int, offset int) ([]*models.MediaItem[*mediatypes.Collection], error)
// 	GetRecent(ctx context.Context, days int, limit int) ([]*models.MediaItem[*mediatypes.Collection], error)
//
// 	// Smart collections and sharing
// 	CreateSmartCollection(ctx context.Context, userID uint64, name string, criteria map[string]interface{}) (*models.MediaItem[*mediatypes.Collection], error)
// 	ShareCollectionWithUser(ctx context.Context, collectionID uint64, targetUserID uint64) error
// }
//
// type userCollectionService struct {
// 	userRepo      repository.UserMediaItemRepository[*mediatypes.Collection]
// 	coreMediaRepo repository.MediaItemRepository[mediatypes.MediaData] // For fetching collection items
// }
//
// // NewUserCollectionService creates a new user collection service
// func NewUserCollectionService(
// 	userRepo repository.UserMediaItemRepository[*mediatypes.Collection],
// 	coreMediaRepo repository.MediaItemRepository[mediatypes.MediaData],
// ) UserCollectionService {
// 	return &userCollectionService{
// 		userRepo:      userRepo,
// 		coreMediaRepo: coreMediaRepo,
// 	}
// }
//
// // Implement all methods from CoreCollectionService through delegation
//
// // Create adds a new collection
// func (s *userCollectionService) Create(ctx context.Context, collection models.MediaItem[*mediatypes.Collection]) (*models.MediaItem[*mediatypes.Collection], error) {
// 	log := utils.LoggerFromContext(ctx)
// 	log.Debug().
// 		Str("title", collection.Title).
// 		Msg("Creating user collection")
//
// 	// Ensure collection-specific validation
// 	if collection.Type != mediatypes.MediaTypeCollection {
// 		collection.Type = mediatypes.MediaTypeCollection
// 	}
//
// 	// get user id from context
// 	userID := ctx.Value("userID").(uint64)
//
// 	// Ensure collection has a valid name
// 	if collection.Data == nil || collection.Data.Details.Title == "" {
// 		return nil, errors.New("collection must have a title")
// 	}
//
// 	// Initialize Items array if nil
// 	if collection.Data.Items == nil {
// 		collection.Data.Items = []mediatypes.ListItem{}
// 	}
//
// 	// Add user ownership metadata
// 	// TODO:
// 	if collection.Data.OwnerID == 0 && userID != 0 {
// 		collection.Data.OwnerID = userID
// 		collection.Data.ModifiedBy = userID
// 	}
//
// 	// Create using the core service
// 	result, err := s.coreService.Create(ctx, collection)
// 	if err != nil {
// 		log.Error().Err(err).
// 			Str("title", collection.Title).
// 			Msg("Failed to create user collection")
// 		return nil, fmt.Errorf("failed to create user collection: %w", err)
// 	}
//
// 	log.Info().
// 		Uint64("id", result.ID).
// 		Str("title", result.Title).
// 		Uint64("userID", collection.Data.OwnerID).
// 		Msg("User collection created successfully")
//
// 	return result, nil
// }
//
// // Update modifies an existing collection
// func (s *userCollectionService) Update(ctx context.Context, collection models.MediaItem[*mediatypes.Collection]) (*models.MediaItem[*mediatypes.Collection], error) {
// 	log := utils.LoggerFromContext(ctx)
// 	log.Debug().
// 		Uint64("id", collection.ID).
// 		Str("title", collection.Title).
// 		Msg("Updating user collection")
//
// 	userID := ctx.Value("userID").(uint64)
//
// 	// Ensure the collection exists
// 	existing, err := s.GetByID(ctx, collection.ID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to update user collection: %w", err)
// 	}
//
// 	// Verify ownership - user can only modify their own collections or shared ones
// 	if existing.Data.OwnerID != 0 && collection.Data.OwnerID != 0 && existing.Data.OwnerID != collection.Data.OwnerID {
// 		// TODO: Check if the collection is shared with this user before denying
// 		log.Warn().
// 			Uint64("collectionID", collection.ID).
// 			Uint64("ownerID", existing.Data.OwnerID).
// 			Uint64("requestingUserID", userID).
// 			Msg("User attempting to update a collection they don't own")
// 		return nil, errors.New("you don't have permission to update this collection")
// 	}
//
// 	// Preserve owner information
// 	collection.Data.OwnerID = existing.Data.OwnerID
//
// 	// Update ModifiedBy field
// 	collection.Data.ModifiedBy = userID
//
// 	// Delegate to core service for update
// 	result, err := s.coreService.Update(ctx, collection)
// 	if err != nil {
// 		log.Error().Err(err).
// 			Uint64("id", collection.ID).
// 			Msg("Failed to update user collection")
// 		return nil, fmt.Errorf("failed to update user collection: %w", err)
// 	}
//
// 	log.Info().
// 		Uint64("id", result.ID).
// 		Str("title", result.Title).
// 		Uint64("userID", collection.Data.OwnerID).
// 		Msg("User collection updated successfully")
//
// 	return result, nil
// }
//
// // GetByID retrieves a collection by its ID
// func (s *userCollectionService) GetByID(ctx context.Context, id uint64) (*models.MediaItem[*mediatypes.Collection], error) {
// 	return s.coreService.GetByID(ctx, id)
// }
//
// func (s *userCollectionService) GetAll(ctx context.Context, limit int, offset int) ([]*models.MediaItem[*mediatypes.Collection], error) {
// 	return s.coreService.GetAll(ctx, limit, offset)
// }
//
// // Delete removes a collection
// func (s *userCollectionService) Delete(ctx context.Context, id uint64) error {
// 	log := utils.LoggerFromContext(ctx)
// 	log.Debug().
// 		Uint64("id", id).
// 		Msg("Deleting user collection")
//
// 	// First verify this is a collection and the user has permission
// 	collection, err := s.GetByID(ctx, id)
// 	if err != nil {
// 		return fmt.Errorf("failed to delete user collection: %w", err)
// 	}
//
// 	// TODO: Extract current user ID from context and verify permission
// 	// For now, just check that it's a user-owned collection
// 	if collection.Data.OwnerID == 0 {
// 		log.Warn().
// 			Uint64("id", id).
// 			Msg("Attempting to delete a collection that isn't user-owned")
// 		return errors.New("collection is not user-owned")
// 	}
//
// 	// Delegate to core service for deletion
// 	err = s.Delete(ctx, id)
// 	if err != nil {
// 		log.Error().Err(err).
// 			Uint64("id", id).
// 			Msg("Failed to delete user collection")
// 		return fmt.Errorf("failed to delete user collection: %w", err)
// 	}
//
// 	log.Info().
// 		Uint64("id", id).
// 		Str("title", collection.Title).
// 		Uint64("userID", collection.Data.OwnerID).
// 		Msg("User collection deleted successfully")
//
// 	return nil
// }
//
// // GetByType retrieves all collections of a specific type
// func (s *userCollectionService) GetByType(ctx context.Context, mediaType mediatypes.MediaType) ([]*models.MediaItem[*mediatypes.Collection], error) {
// 	return s.GetByType(ctx, mediaType)
// }
//
// // GetByExternalID retrieves a collection by its external ID
// func (s *userCollectionService) GetByExternalID(ctx context.Context, source string, externalID string) (*models.MediaItem[*mediatypes.Collection], error) {
// 	return s.GetByExternalID(ctx, source, externalID)
// }
//
// // Search finds collections based on a query string
// func (s *userCollectionService) Search(ctx context.Context, query mediatypes.QueryOptions) ([]*models.MediaItem[*mediatypes.Collection], error) {
// 	return s.Search(ctx, query)
// }
//
// // GetRecentItems retrieves recently added collections
// func (s *userCollectionService) GetRecentItems(ctx context.Context, days int, limit int) ([]*models.MediaItem[*mediatypes.Collection], error) {
// 	return s.GetRecentItems(ctx, days, limit)
// }
//
// // GetCollectionItems retrieves all items in a collection
// func (s *userCollectionService) GetCollectionItems(ctx context.Context, collectionID uint64) (*models.MediaItems, error) {
// 	return s.GetCollectionItems(ctx, collectionID)
// }
//
// // AddItemToCollection adds an item to a collection
// func (s *userCollectionService) AddItemToCollection(ctx context.Context, collectionID uint64, itemID uint64) error {
// 	log := utils.LoggerFromContext(ctx)
// 	log.Debug().
// 		Uint64("collectionID", collectionID).
// 		Uint64("itemID", itemID).
// 		Msg("Adding item to user collection")
//
// 	// Get the collection
// 	collection, err := s.GetByID(ctx, collectionID)
// 	if err != nil {
// 		return fmt.Errorf("failed to add item to user collection: %w", err)
// 	}
//
// 	log.Debug().
// 		Uint64("collectionID", collection.ID).
// 		Uint64("itemID", itemID).
// 		Msg("Adding item to user collection")
//
// 	// TODO: Verify user has permission to modify this collection
// 	// Extract user ID from context and check against collection.UserID
// 	// or check if the collection is shared with this user
//
// 	// Delegate to core service
// 	err = s.AddItemToCollection(ctx, collectionID, itemID)
// 	if err != nil {
// 		log.Error().Err(err).
// 			Uint64("collectionID", collectionID).
// 			Uint64("itemID", itemID).
// 			Msg("Failed to add item to user collection")
// 		return fmt.Errorf("failed to add item to user collection: %w", err)
// 	}
//
// 	log.Info().
// 		Uint64("collectionID", collectionID).
// 		Uint64("itemID", itemID).
// 		Msg("Item added to user collection successfully")
//
// 	return nil
// }
//
// // RemoveItemFromCollection removes an item from a collection
// func (s *userCollectionService) RemoveItemFromCollection(ctx context.Context, collectionID uint64, itemID uint64) error {
// 	log := utils.LoggerFromContext(ctx)
// 	log.Debug().
// 		Uint64("collectionID", collectionID).
// 		Uint64("itemID", itemID).
// 		Msg("Removing item from user collection")
//
// 	// Get the collection
// 	collection, err := s.GetByID(ctx, collectionID)
// 	if err != nil {
// 		return fmt.Errorf("failed to remove item from user collection: %w", err)
// 	}
//
// 	// TODO: Verify user has permission to modify this collection
// 	// Extract user ID from context and check against collection.UserID
// 	// or check if the collection is shared with this user
// 	log.Debug().
// 		Uint64("collectionID", collection.ID).
// 		Uint64("itemID", itemID).
// 		Msg("Removing item from user collection")
//
// 	// Delegate to core service
// 	err = s.RemoveItemFromCollection(ctx, collectionID, itemID)
// 	if err != nil {
// 		log.Error().Err(err).
// 			Uint64("collectionID", collectionID).
// 			Uint64("itemID", itemID).
// 			Msg("Failed to remove item from user collection")
// 		return fmt.Errorf("failed to remove item from user collection: %w", err)
// 	}
//
// 	log.Info().
// 		Uint64("collectionID", collectionID).
// 		Uint64("itemID", itemID).
// 		Msg("Item removed from user collection successfully")
//
// 	return nil
// }
//
// // UpdateCollectionItems replaces all items in a collection
// func (s *userCollectionService) UpdateCollectionItems(ctx context.Context, collectionID uint64, items []models.MediaItem[mediatypes.MediaData]) error {
// 	return s.UpdateCollectionItems(ctx, collectionID, items)
// }
//
// // User-specific methods
//
// // GetByUserID retrieves all collections for a user
// func (s *userCollectionService) GetByUserID(ctx context.Context, userID uint64, limit int, offset int) ([]*models.MediaItem[*mediatypes.Collection], error) {
// 	log := utils.LoggerFromContext(ctx)
// 	log.Debug().
// 		Uint64("userID", userID).
// 		Msg("Getting collections by user ID")
//
// 	// Use the repository directly for efficiency
// 	collections, err := s.userRepo.GetByUserID(ctx, userID, limit, offset)
// 	if err != nil {
// 		log.Error().Err(err).
// 			Uint64("userID", userID).
// 			Msg("Failed to get collections for user")
// 		return nil, fmt.Errorf("failed to get collections for user: %w", err)
// 	}
//
// 	log.Info().
// 		Uint64("userID", userID).
// 		Int("count", len(collections)).
// 		Msg("Retrieved collections for user")
//
// 	return collections, nil
// }
//
// // SearchUserCollections searches for collections owned by a specific user
// func (s *userCollectionService) SearchUserCollections(ctx context.Context, query mediatypes.QueryOptions, userID uint64) ([]*models.MediaItem[*mediatypes.Collection], error) {
// 	log := utils.LoggerFromContext(ctx)
// 	log.Debug().
// 		Str("query", query.Query).
// 		Uint64("userID", userID).
// 		Msg("Searching user collections")
//
// 	options := mediatypes.QueryOptions{
// 		MediaType: mediatypes.MediaTypeCollection,
// 		Query:     query.Query,
// 		OwnerID:   userID,
// 	}
//
// 	// Use the repository directly for user-specific search
// 	results, err := s.userRepo.Search(ctx, options)
// 	if err != nil {
// 		log.Error().Err(err).
// 			Str("query", query.Query).
// 			Uint64("userID", userID).
// 			Msg("Failed to search user collections")
// 		return nil, fmt.Errorf("failed to search user collections: %w", err)
// 	}
//
// 	log.Info().
// 		Str("query", query.Query).
// 		Uint64("userID", userID).
// 		Int("count", len(results)).
// 		Msg("User collections found")
//
// 	return results, nil
// }
//
// // GetRecentUserCollections retrieves recently updated collections for a user
// func (s *userCollectionService) GetRecent(ctx context.Context, days int, limit int) ([]*models.MediaItem[*mediatypes.Collection], error) {
// 	log := utils.LoggerFromContext(ctx)
//
// 	userID := ctx.Value("userID").(uint64)
//
// 	log.Debug().
// 		Uint64("userID", userID).
// 		Int("limit", limit).
// 		Msg("Getting recent user collections")
//
// 	// Use the repository directly for user-specific recent content
// 	results, err := s.userRepo.GetRecentItems(ctx, days, limit)
// 	if err != nil {
// 		log.Error().Err(err).
// 			Uint64("userID", userID).
// 			Msg("Failed to get recent user collections")
// 		return nil, fmt.Errorf("failed to get recent user collections: %w", err)
// 	}
//
// 	log.Info().
// 		Uint64("userID", userID).
// 		Int("count", len(results)).
// 		Msg("Recent user collections retrieved")
//
// 	return results, nil
// }
//
// // CreateSmartCollection creates a collection that updates automatically based on criteria
// func (s *userCollectionService) CreateSmartCollection(ctx context.Context, userID uint64, name string, criteria map[string]interface{}) (*models.MediaItem[*mediatypes.Collection], error) {
// 	log := utils.LoggerFromContext(ctx)
// 	log.Debug().
// 		Uint64("userID", userID).
// 		Str("name", name).
// 		Interface("criteria", criteria).
// 		Msg("Creating smart collection")
//
// 	// Create a new collection
// 	collection := models.MediaItem[*mediatypes.Collection]{
// 		Title: name,
// 		Type:  mediatypes.MediaTypeCollection,
// 		Data: &mediatypes.Collection{
// 			ItemList: mediatypes.ItemList{
// 				Details: mediatypes.MediaDetails{
// 					Title:       name,
// 					Description: "Smart collection - automatically updates based on criteria",
// 					AddedAt:     time.Now(),
// 				},
// 				OwnerID:    userID,
// 				ModifiedBy: userID,
// 				Items:      []mediatypes.ListItem{},
// 				ItemCount:  0,
// 				// Add smart collection metadata
// 				IsSmart:        true,
// 				SmartCriteria:  criteria,
// 				AutoUpdateTime: time.Now(),
// 			},
// 		},
// 	}
//
// 	// Create the collection
// 	result, err := s.Create(ctx, collection)
// 	if err != nil {
// 		log.Error().Err(err).
// 			Uint64("userID", userID).
// 			Str("name", name).
// 			Msg("Failed to create smart collection")
// 		return nil, fmt.Errorf("failed to create smart collection: %w", err)
// 	}
//
// 	// Now populate the collection based on criteria
// 	// This would typically call a method to evaluate the criteria and find matching items
// 	// For now, we'll just return the empty collection
// 	log.Info().
// 		Uint64("id", result.ID).
// 		Str("title", result.Title).
// 		Uint64("userID", userID).
// 		Msg("Smart collection created successfully")
//
// 	return result, nil
// }
//
// // ShareCollectionWithUser shares a collection with another user
// func (s *userCollectionService) ShareCollectionWithUser(ctx context.Context, collectionID uint64, targetUserID uint64) error {
// 	log := utils.LoggerFromContext(ctx)
// 	log.Debug().
// 		Uint64("collectionID", collectionID).
// 		Uint64("targetUserID", targetUserID).
// 		Msg("Sharing collection with user")
//
// 	// Get the collection
// 	collection, err := s.GetByID(ctx, collectionID)
// 	if err != nil {
// 		return fmt.Errorf("failed to share collection: %w", err)
// 	}
//
// 	// TODO: Verify current user has permission to share this collection
// 	// Extract current user ID from context and check against collection.UserID
//
// 	// Add the target user to the shared users list if not already present
// 	// if collection.Data.SharedWith == nil {
// 	// 	collection.Data.SharedWith = []uint64{targetUserID}
// 	// } else {
// 	// 	// Check if already shared
// 	// 	alreadyShared := false
// 	// 	for _, uid := range collection.Data.SharedWith {
// 	// 		if uid == targetUserID {
// 	// 			alreadyShared = true
// 	// 			break
// 	// 		}
// 	// 	}
// 	//
// 	// 	if !alreadyShared {
// 	// 		collection.Data.SharedWith = append(collection.Data.SharedWith, targetUserID)
// 	// 	} else {
// 	// 		log.Info().
// 	// 			Uint64("collectionID", collectionID).
// 	// 			Uint64("targetUserID", targetUserID).
// 	// 			Msg("Collection already shared with this user")
// 	// 		return nil
// 	// 	}
// 	// }
//
// 	// Update the collection
// 	_, err = s.Update(ctx, *collection)
// 	if err != nil {
// 		log.Error().Err(err).
// 			Uint64("collectionID", collectionID).
// 			Uint64("targetUserID", targetUserID).
// 			Msg("Failed to update collection sharing information")
// 		return fmt.Errorf("failed to update collection sharing information: %w", err)
// 	}
//
// 	log.Info().
// 		Uint64("collectionID", collectionID).
// 		Uint64("targetUserID", targetUserID).
// 		Msg("Collection shared successfully")
//
// 	return nil
// }
