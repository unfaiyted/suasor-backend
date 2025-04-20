package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	mediatypes "suasor/client/media/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/utils"
)

// UserlistService defines the interface for user-owned list operations
// This service extends listService with operations specific to user-owned lists
type UserListService[T mediatypes.ListData] interface {
	// Include all core list service methods
	CoreListService[T]

	// Create(c *gin.Context)
	// Update(c *gin.Context)
	// Delete(c *gin.Context)
	// AddItem(c *gin.Context)
	// RemoveItem(c *gin.Context)
	// ReorderItems(c *gin.Context)

	// User-specific operations
	GetFavorite(ctx context.Context, userID uint64, limit int, offset int) ([]*models.MediaItem[T], error)
	GetRecentByUser(ctx context.Context, userID uint64, days int, limit int) ([]*models.MediaItem[T], error)

	// Smart list operations
	CreateSmartList(ctx context.Context,
		userID uint64, name string,
		description string,
		criteria map[string]interface{},
	) (*models.MediaItem[T], error)
	UpdateSmartCriteria(ctx context.Context, listID uint64, criteria map[string]interface{}) (*models.MediaItem[T], error)
	RefreshSmartList(ctx context.Context, listID uint64) (*models.MediaItem[T], error)

	// list sharing and collaboration
	ShareWithUser(ctx context.Context, listID uint64, targetUserID uint64, permissionLevel string) error
	GetShared(ctx context.Context, userID uint64) ([]*models.MediaItem[T], error)
	// GetCollaborators(ctx context.Context, listID uint64) ([]models.ListCollaborator, error)
	// RemoveCollaborator(ctx context.Context, listID uint64, userID uint64) error

	// list sync
	SyncToClients(ctx context.Context, listID uint64, clientIDs []uint64) error
	GetSyncStatus(ctx context.Context, listID uint64) (*models.ListSyncStatus, error)
}

type userListService[T mediatypes.ListData] struct {
	userItemRepo    repository.UserMediaItemRepository[T]
	userDataRepo    repository.UserMediaItemDataRepository[T]
	coreListService CoreListService[T]
}

// NewUserlistService creates a new user list service
func NewUserListService[T mediatypes.ListData](
	coreListService CoreListService[T],
	userItemRepo repository.UserMediaItemRepository[T],
	userDataRepo repository.UserMediaItemDataRepository[T],
) UserListService[T] {
	return &userListService[T]{
		coreListService: coreListService,
		userItemRepo:    userItemRepo,
		userDataRepo:    userDataRepo,
	}
}

// Core methods delegated to the base listService
func (s *userListService[T]) GetAll(ctx context.Context, limit int, offset int) ([]*models.MediaItem[T], error) {
	return s.coreListService.GetAll(ctx, limit, offset)
}

func (s *userListService[T]) Create(ctx context.Context, list *models.MediaItem[T]) (*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Str("title", list.Title).
		Msg("Creating user list")

	// Get user ID from context
	userID := ctx.Value("userID").(uint64)

	// Set ownership info if not already set
	if list.OwnerID == 0 {
		list.OwnerID = userID
	}

	itemList := list.GetData().GetItemList()

	if itemList.OwnerID == 0 {
		itemList.OwnerID = userID
	}
	if itemList.ModifiedBy == 0 {
		itemList.ModifiedBy = userID
	}

	return s.coreListService.Create(ctx, list)
}

func (s *userListService[T]) Update(ctx context.Context, list *models.MediaItem[T]) (*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("id", list.ID).
		Str("title", list.Title).
		Msg("Updating user list")

	// Verify user has permission to update this list
	userID := ctx.Value("userID").(uint64)
	existing, err := s.GetByID(ctx, list.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update user list: %w", err)
	}

	// Check ownership or collaboration permission
	if !s.hasWritePermission(ctx, userID, existing) {
		log.Warn().
			Uint64("listID", list.ID).
			Uint64("ownerID", existing.OwnerID).
			Uint64("requestingUserID", userID).
			Msg("User attempting to update a list without permission")
		return nil, errors.New("you don't have permission to update this list")
	}

	itemList := list.GetData().GetItemList()

	// Set the modified by field to the current user
	itemList.ModifiedBy = userID
	itemList.LastModified = time.Now()

	list.GetData().SetItemList(itemList)

	// Delegate to core service
	return s.Update(ctx, list)
}

func (s *userListService[T]) GetByID(ctx context.Context, id uint64) (*models.MediaItem[T], error) {
	return s.GetByID(ctx, id)
}

func (s *userListService[T]) GetByUserID(ctx context.Context, userID uint64, limit int, offset int) ([]*models.MediaItem[T], error) {
	return s.coreListService.GetByUserID(ctx, userID, limit, offset)
}

func (s *userListService[T]) Delete(ctx context.Context, id uint64) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("id", id).
		Msg("Deleting user list")

	// Verify user has permission to delete this list
	userID := ctx.Value("userID").(uint64)
	list, err := s.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete user list: %w", err)
	}

	// Only the owner can delete a list
	if list.OwnerID != userID {
		log.Warn().
			Uint64("listID", id).
			Uint64("ownerID", list.OwnerID).
			Uint64("requestingUserID", userID).
			Msg("User attempting to delete a list they don't own")
		return errors.New("only the owner can delete a list")
	}

	// Delegate to core service for deletion
	return s.Delete(ctx, id)
}

func (s *userListService[T]) GetItems(ctx context.Context, listID uint64) (*models.MediaItemList, error) {
	return s.coreListService.GetItems(ctx, listID)
}

func (s *userListService[T]) AddItem(ctx context.Context, listID uint64, itemID uint64) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("listID", listID).
		Uint64("itemID", itemID).
		Msg("Adding item to user list")

	// Verify user has permission to modify this list
	userID := ctx.Value("userID").(uint64)
	list, err := s.GetByID(ctx, listID)
	if err != nil {
		return fmt.Errorf("failed to add item to user list: %w", err)
	}

	// Check if user has write permission
	if !s.hasWritePermission(ctx, userID, list) {
		log.Warn().
			Uint64("listID", listID).
			Uint64("ownerID", list.OwnerID).
			Uint64("requestingUserID", userID).
			Msg("User attempting to modify a list without permission")
		return errors.New("you don't have permission to modify this list")
	}

	// Delegate to core service
	return s.AddItem(ctx, listID, itemID)
}

func (s *userListService[T]) RemoveItem(ctx context.Context, listID uint64, itemID uint64) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("listID", listID).
		Uint64("itemID", itemID).
		Msg("Removing item from user list")

	// Verify user has permission to modify this list
	userID := ctx.Value("userID").(uint64)
	list, err := s.GetByID(ctx, listID)
	if err != nil {
		return fmt.Errorf("failed to remove item from user list: %w", err)
	}

	// Check if user has write permission
	if !s.hasWritePermission(ctx, userID, list) {
		log.Warn().
			Uint64("listID", listID).
			Uint64("ownerID", list.OwnerID).
			Uint64("requestingUserID", userID).
			Msg("User attempting to modify a list without permission")
		return errors.New("you don't have permission to modify this list")
	}

	// Delegate to core service
	return s.coreListService.RemoveItem(ctx, listID, itemID)
}

// func (s *userListService[T]) ReorderItems(ctx context.Context, listID uint64, itemIDs []uint64) error {
// 	log := utils.LoggerFromContext(ctx)
// 	log.Debug().
// 		Uint64("listID", listID).
// 		Interface("itemIDs", itemIDs).
// 		Msg("Reordering user list items")
//
// 	// Verify user has permission to modify this list
// 	userID := ctx.Value("userID").(uint64)
// 	list, err := s.GetByID(ctx, listID)
// 	if err != nil {
// 		return fmt.Errorf("failed to reorder user list items: %w", err)
// 	}
//
// 	// Check if user has write permission
// 	if !s.hasWritePermission(ctx, userID, list) {
// 		log.Warn().
// 			Uint64("listID", listID).
// 			Uint64("ownerID", list.OwnerID).
// 			Uint64("requestingUserID", userID).
// 			Msg("User attempting to modify a list without permission")
// 		return errors.New("you don't have permission to modify this list")
// 	}
//
// 	// Delegate to core service
// 	return s.coreListService.ReorderItems(ctx, listID, itemIDs)
// }

func (s *userListService[T]) UpdateItems(ctx context.Context, listID uint64, items []*models.MediaItem[T]) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("listID", listID).
		Int("itemCount", len(items)).
		Msg("Updating user list items")

	// Verify user has permission to modify this list
	userID := ctx.Value("userID").(uint64)
	list, err := s.GetByID(ctx, listID)
	if err != nil {
		return fmt.Errorf("failed to update user list items: %w", err)
	}

	// Check if user has write permission
	if !s.hasWritePermission(ctx, userID, list) {
		log.Warn().
			Uint64("listID", listID).
			Uint64("ownerID", list.OwnerID).
			Uint64("requestingUserID", userID).
			Msg("User attempting to modify a list without permission")
		return errors.New("you don't have permission to modify this list")
	}

	// Delegate to core service
	return s.UpdateItems(ctx, listID, items)
}

func (s *userListService[T]) Search(ctx context.Context, query mediatypes.QueryOptions) ([]*models.MediaItem[T], error) {
	return s.Search(ctx, query)
}

func (s *userListService[T]) GetRecent(ctx context.Context, days int, limit int) ([]*models.MediaItem[T], error) {
	return s.GetRecent(ctx, days, limit)
}

func (s *userListService[T]) Sync(ctx context.Context, listID uint64, targetClientIDs []uint64) error {
	// Verify user has permission to synchronize this list
	userID := ctx.Value("userID").(uint64)
	list, err := s.GetByID(ctx, listID)
	if err != nil {
		return fmt.Errorf("failed to sync user list: %w", err)
	}

	// Check if user has read permission (sync is a read operation followed by creation elsewhere)
	if !s.haslistReadPermission(ctx, userID, list) {
		return errors.New("you don't have permission to sync this list")
	}

	// Delegate to core service
	return s.Sync(ctx, listID, targetClientIDs)
}

func (s *userListService[T]) UpdateSmartCriteria(ctx context.Context, listID uint64, criteria map[string]interface{}) (*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("listID", listID).
		Interface("criteria", criteria).
		Msg("Updating smart list criteria")

	return nil, nil
	// return s.UpdateSmartCriteria(ctx, listID, criteria)
}

// User-specific operations

// GetUserlists retrieves lists owned by a specific user with pagination
func (s *userListService[T]) GetUser(ctx context.Context, userID uint64, limit int, offset int) ([]*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Int("limit", limit).
		Int("offset", offset).
		Msg("Getting user lists with pagination")

	// Get all lists for this user
	lists, err := s.userItemRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to get user lists")
		return nil, fmt.Errorf("failed to get user lists: %w", err)
	}

	log.Info().
		Uint64("userID", userID).
		Int("count", len(lists)).
		Msg("Retrieved user lists")

	return lists, nil
}

// SearchUserlists searches for lists owned by a specific user
func (s *userListService[T]) SearchUser(ctx context.Context, userID uint64, query string) ([]*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Str("query", query).
		Msg("Searching user lists")

	var zero T
	mediaType := mediatypes.GetMediaTypeFromTypeName(zero)

	// Create query options with user filter
	options := mediatypes.QueryOptions{
		MediaType: mediaType,
		Query:     query,
		OwnerID:   userID,
	}

	// Delegate to core service with owner filter
	return s.Search(ctx, options)
}

// GetRecentUserlists retrieves recently updated lists for a user
func (s *userListService[T]) GetRecentByUser(ctx context.Context, userID uint64, days int, limit int) ([]*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Int("limit", limit).
		Msg("Getting recent user lists")

	var zero T
	mediaType := mediatypes.GetMediaTypeFromTypeName(zero)

	options := mediatypes.QueryOptions{
		MediaType: mediaType,
		OwnerID:   userID,
		Limit:     limit,
		Sort:      "updated_at",
		SortOrder: "desc",
	}
	return s.userItemRepo.Search(ctx, options)

}

// GetFavoritelists retrieves lists marked as favorite by the user
func (s *userListService[T]) GetFavorite(ctx context.Context, userID uint64, limit int, offset int) ([]*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Msg("Getting favorite lists")

	// Use user data repository to get all favorites of type list
	userFavoritePlayData, err := s.userDataRepo.GetFavorites(ctx, userID, limit, offset)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to get favorite list IDs")
		return nil, fmt.Errorf("failed to get favorite lists: %w", err)
	}

	var ids []uint64
	for _, data := range userFavoritePlayData {
		ids = append(ids, data.MediaItemID)
	}

	// Fetch the lists by IDs
	lists, err := s.userItemRepo.GetByIDs(ctx, ids)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to get favorite lists")
		return nil, fmt.Errorf("failed to get favorite lists: %w", err)
	}

	log.Info().
		Uint64("userID", userID).
		Int("count", len(lists)).
		Msg("Retrieved favorite lists")

	return lists, nil
}

// Smart list operations

// CreateSmartlist creates a list that updates automatically based on criteria
func (s *userListService[T]) CreateSmartList(ctx context.Context, userID uint64, name string, description string, criteria map[string]interface{}) (*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Str("name", name).
		Interface("criteria", criteria).
		Msg("Creating smart list")

	var zero T
	mediaType := mediatypes.GetMediaTypeFromTypeName(zero)

	// Create a new list with smart flag enabled

	data := createList[T](name, description, criteria, userID)
	list := models.NewMediaItem[T](mediaType, data)

	// Create the list
	result, err := s.Create(ctx, list)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Str("name", name).
			Msg("Failed to create smart list")
		return nil, fmt.Errorf("failed to create smart list: %w", err)
	}

	// Refresh the smart list to populate it based on criteria
	refreshed, err := s.RefreshSmartList(ctx, result.ID)
	if err != nil {
		log.Warn().Err(err).
			Uint64("listID", result.ID).
			Msg("Failed to initially populate smart list")
		// Return the list anyway, just warn about the population failure
	} else if refreshed != nil {
		result = refreshed
	}

	itemList := result.GetData().GetItemList()

	log.Info().
		Uint64("listID", result.ID).
		Str("name", name).
		Uint64("userID", userID).
		Int("itemCount", itemList.ItemCount).
		Msg("Smart list created successfully")

	return result, nil
}

// UpdateSmartlistCriteria updates the criteria for a smart list
func (s *userListService[T]) UpdateSmartlistCriteria(ctx context.Context, listID uint64, criteria map[string]interface{}) (*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("listID", listID).
		Interface("criteria", criteria).
		Msg("Updating smart list criteria")

	// Get the list
	list, err := s.GetByID(ctx, listID)
	if err != nil {
		return nil, fmt.Errorf("failed to update smart list criteria: %w", err)
	}

	// Verify this is a smart list
	// if !list.Data.IsSmart {
	// 	log.Error().
	// 		Uint64("listID", listID).
	// 		Msg("Cannot update criteria of non-smart list")
	// 	return nil, errors.New("cannot update criteria of non-smart list")
	// }

	// Verify user has permission to modify this list
	userID := ctx.Value("userID").(uint64)
	if !s.hasWritePermission(ctx, userID, list) {
		log.Warn().
			Uint64("listID", listID).
			Uint64("ownerID", list.OwnerID).
			Uint64("requestingUserID", userID).
			Msg("User attempting to modify a list without permission")
		return nil, errors.New("you don't have permission to modify this list")
	}

	// Update the criteria
	// list.SmartCriteria = criteria
	// list.AutoUpdateTime = time.Now()
	// list.LastModified = time.Now()
	// list.ModifiedBy = userID

	// Update the list
	updated, err := s.Update(ctx, list)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Msg("Failed to update smart list criteria")
		return nil, fmt.Errorf("failed to update smart list criteria: %w", err)
	}

	// Refresh the list to apply the new criteria
	refreshed, err := s.RefreshSmartList(ctx, listID)
	if err != nil {
		log.Warn().Err(err).
			Uint64("listID", listID).
			Msg("Failed to refresh smart list after criteria update")
		// Return the updated list anyway
		return updated, nil
	}

	log.Info().
		Uint64("listID", listID).
		// Int("itemCount", refreshed.Data.ItemCount).
		Msg("Smart list criteria updated and refreshed successfully")

	return refreshed, nil
}

// RefreshSmartlist updates a smart list based on its criteria
func (s *userListService[T]) RefreshSmartList(ctx context.Context, listID uint64) (*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("listID", listID).
		Msg("Refreshing smart list")

	// Get the list
	list, err := s.GetByID(ctx, listID)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh smart list: %w", err)
	}
	itemList := list.GetData().GetItemList()

	// Verify this is a smart list
	if !itemList.IsSmart {
		log.Error().
			Uint64("listID", listID).
			Msg("Cannot refresh non-smart list")
		return nil, errors.New("cannot refresh non-smart list")
	}

	// Verify user has permission to read this list (refresh is a read operation)
	userID := ctx.Value("userID").(uint64)
	if !s.haslistReadPermission(ctx, userID, list) {
		log.Warn().
			Uint64("listID", listID).
			Uint64("ownerID", list.OwnerID).
			Uint64("requestingUserID", userID).
			Msg("User attempting to refresh a list without permission")
		return nil, errors.New("you don't have permission to refresh this list")
	}

	// Get the criteria
	// criteria := list.Data.SmartCriteria

	// In a real implementation, this would:
	// 1. Translate the criteria into a database query
	// 2. Execute the query to find all matching media items
	// 3. Replace the list items with the query results
	// 4. Update the list metadata

	// For now, just simulate the refresh by adding a note to the description
	// now := time.Now()
	// list.Data.Details.Description = fmt.Sprintf("%s\n\nLast refreshed: %s",
	// 	list.Data.Details.Description, now.Format(time.RFC3339))
	// list.Data.AutoUpdateTime = now
	// list.Data.LastModified = now
	// list.Data.ModifiedBy = userID

	// Update the list
	updated, err := s.Update(ctx, list)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Msg("Failed to update list after refresh")
		return nil, fmt.Errorf("failed to update list after refresh: %w", err)
	}

	log.Info().
		Uint64("listID", listID).
		Msg("Smart list refreshed successfully")

	return updated, nil
}

// list sharing and collaboration

// SharelistWithUser shares a list with another user
func (s *userListService[T]) ShareWithUser(ctx context.Context, listID uint64, targetUserID uint64, permissionLevel string) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("listID", listID).
		Uint64("targetUserID", targetUserID).
		Str("permissionLevel", permissionLevel).
		Msg("Sharing list with user")

	// Get the list
	list, err := s.GetByID(ctx, listID)
	if err != nil {
		return fmt.Errorf("failed to share list: %w", err)
	}

	// Verify user has permission to share this list (only owner can share)
	userID := ctx.Value("userID").(uint64)
	if list.OwnerID != userID {
		log.Warn().
			Uint64("listID", listID).
			Uint64("ownerID", list.OwnerID).
			Uint64("requestingUserID", userID).
			Msg("User attempting to share a list they don't own")
		return errors.New("only the owner can share a list")
	}

	// Validate permission level
	if permissionLevel != "read" && permissionLevel != "write" {
		return errors.New("invalid permission level: must be 'read' or 'write'")
	}

	// Add the target user to the shared users list if not already present
	// collaborator := models.ListCollaborator{
	// 	UserID:          targetUserID,
	// 	PermissionLevel: permissionLevel,
	// 	SharedAt:        time.Now(),
	// 	SharedBy:        userID,
	// }
	//
	// // Initialize the shared with array if it doesn't exist
	// if list.Data.SharedWith == nil {
	// 	list.Data.SharedWith = []models.ListCollaborator{collaborator}
	// } else {
	// 	// Check if already shared
	// 	alreadyShared := false
	// 	for i, collab := range list.Data.SharedWith {
	// 		if collab.UserID == targetUserID {
	// 			alreadyShared = true
	// 			// Update permission level if it's different
	// 			if collab.PermissionLevel != permissionLevel {
	// 				list.Data.SharedWith[i].PermissionLevel = permissionLevel
	// 				list.Data.SharedWith[i].SharedAt = time.Now()
	// 				list.Data.SharedWith[i].SharedBy = userID
	// 			}
	// 			break
	// 		}
	// 	}
	//
	// 	if !alreadyShared {
	// 		list.Data.SharedWith = append(list.Data.SharedWith, collaborator)
	// 	}
	// }

	// Update the list
	_, err = s.Update(ctx, list)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Uint64("targetUserID", targetUserID).
			Msg("Failed to update list sharing information")
		return fmt.Errorf("failed to update list sharing information: %w", err)
	}

	log.Info().
		Uint64("listID", listID).
		Uint64("targetUserID", targetUserID).
		Str("permissionLevel", permissionLevel).
		Msg("list shared successfully")

	return nil
}

// GetSharedlists retrieves lists shared with a user
func (s *userListService[T]) GetShared(ctx context.Context, userID uint64) ([]*models.MediaItem[T], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Msg("Getting lists shared with user")

	// TODO:
	// In a real implementation, this would query the database for lists where
	// the user is in the SharedWith array
	// For now, we'll scan through all lists to find ones shared with this user
	// alllists, err := s.userRepo.GetAll(ctx, 1000, 0)
	// if err != nil {
	// 	log.Error().Err(err).
	// 		Uint64("userID", userID).
	// 		Msg("Failed to get all lists")
	// 	return nil, fmt.Errorf("failed to get shared lists: %w", err)
	// }
	//
	// var sharedlists []*models.MediaItem[T]
	// for _, list := range alllists {
	// 	// Skip lists owned by this user (those are covered by GetUserlists)
	// 	if list.Data.OwnerID == userID {
	// 		continue
	// 	}
	//
	// 	// Check if this list is shared with the user
	// 	for _, collab := range list.Data.SharedWith {
	// 		if collab.UserID == userID {
	// 			sharedlists = append(sharedlists, list)
	// 			break
	// 		}
	// 	}
	// }
	//
	// log.Info().
	// 	Uint64("userID", userID).
	// 	Int("count", len(sharedlists)).
	// 	Msg("Retrieved lists shared with user")
	//
	// return sharedlists, nil
	return nil, nil
}

// GetlistCollaborators retrieves the list of users a list is shared with
func (s *userListService[T]) GetCollaborators(ctx context.Context, listID uint64) ([]models.ListCollaborator, error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("listID", listID).
		Msg("Getting list collaborators")

	// Get the list
	list, err := s.GetByID(ctx, listID)
	if err != nil {
		return nil, fmt.Errorf("failed to get list collaborators: %w", err)
	}

	// Verify user has permission to view this list
	userID := ctx.Value("userID").(uint64)
	if !s.haslistReadPermission(ctx, userID, list) {
		log.Warn().
			Uint64("listID", listID).
			Uint64("ownerID", list.OwnerID).
			Uint64("requestingUserID", userID).
			Msg("User attempting to view collaborators without permission")
		return nil, errors.New("you don't have permission to view this list's collaborators")
	}
	itemList := list.GetData().GetItemList()

	// Return the shared with array (may be nil)
	if itemList.SharedWith == nil {
		return []models.ListCollaborator{}, nil
	}

	// return list.Data.SharedWith, nil
	// TODO: implement all of the collaborator stuff
	return nil, nil
}

// RemovelistCollaborator removes a user from the list's collaborators
func (s *userListService[T]) RemoveCollaborator(ctx context.Context, listID uint64, collaboratorID uint64) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("listID", listID).
		Uint64("collaboratorID", collaboratorID).
		Msg("Removing list collaborator")

	// Get the list
	list, err := s.GetByID(ctx, listID)
	if err != nil {
		return fmt.Errorf("failed to remove list collaborator: %w", err)
	}

	// Verify user has permission to modify sharing (only owner can do this)
	userID := ctx.Value("userID").(uint64)
	if list.OwnerID != userID {
		log.Warn().
			Uint64("listID", listID).
			Uint64("ownerID", list.OwnerID).
			Uint64("requestingUserID", userID).
			Msg("User attempting to modify list sharing without permission")
		return errors.New("only the owner can modify list sharing")
	}
	itemList := list.GetData().GetItemList()

	// Check if the list has any collaborators
	if itemList.SharedWith == nil || len(itemList.SharedWith) == 0 {
		log.Info().
			Uint64("listID", listID).
			Msg("list has no collaborators to remove")
		return nil
	}

	// Find and remove the collaborator
	// var newCollaborators []models.ListCollaborator
	found := false
	// for _, collab := range list.Data.SharedWith {
	// if collab.UserID != collaboratorID {
	// 	newCollaborators = append(newCollaborators, collab)
	// } else {
	// 	found = true
	// }
	// }

	if !found {
		log.Info().
			Uint64("listID", listID).
			Uint64("collaboratorID", collaboratorID).
			Msg("Collaborator not found in list")
		return nil
	}

	// Update the list with the new collaborators list
	// list.Data.SharedWith = newCollaborators
	_, err = s.Update(ctx, list)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Uint64("collaboratorID", collaboratorID).
			Msg("Failed to update list after removing collaborator")
		return fmt.Errorf("failed to update list after removing collaborator: %w", err)
	}

	log.Info().
		Uint64("listID", listID).
		Uint64("collaboratorID", collaboratorID).
		Msg("list collaborator removed successfully")

	return nil
}

// Sync with media clients

// SynclistToClients synchronizes a list to specified media clients
func (s *userListService[T]) SyncToClients(ctx context.Context, listID uint64, clientIDs []uint64) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("listID", listID).
		Interface("clientIDs", clientIDs).
		Msg("Syncing list to clients")

	// This uses the same approach as the base service's Synclist method
	// But could be extended with user-specific validation and tracking
	// return s.Synclist(ctx, listID, clientIDs)
	return nil
}

// GetlistSyncStatus retrieves the sync status of a list across clients
func (s *userListService[T]) GetSyncStatus(ctx context.Context, listID uint64) (*models.ListSyncStatus, error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("listID", listID).
		Msg("Getting list sync status")

	// Get the list
	list, err := s.GetByID(ctx, listID)
	if err != nil {
		return nil, fmt.Errorf("failed to get list sync status: %w", err)
	}

	// Verify user has permission to view this list
	userID := ctx.Value("userID").(uint64)
	if !s.haslistReadPermission(ctx, userID, list) {
		log.Warn().
			Uint64("listID", listID).
			Uint64("ownerID", list.OwnerID).
			Uint64("requestingUserID", userID).
			Msg("User attempting to view list sync status without permission")
		return nil, errors.New("you don't have permission to view this list's sync status")
	}

	// Create a sync status object
	// status := &models.listSyncStatus{
	// 	listID:   listID,
	// 	LastSynced:   list.Data.LastSynced,
	// 	ClientStates: make(map[uint64]models.ClientSyncState),
	// }
	//
	// // Add state information for each client
	// for _, state := range list.Data.SyncClientStates {
	// 	clientState := models.ClientSyncState{
	// 		ClientID:     state.ClientID,
	// 		ClientListID: state.ClientListID,
	// 		LastSynced:   state.LastSynced,
	// 		SyncStatus:   "unknown",
	// 		ItemCount:    len(state.Items),
	// 	}
	//
	// 	// Determine sync status
	// 	if state.LastSynced.IsZero() {
	// 		clientState.SyncStatus = "never_synced"
	// 	} else if state.LastSynced.Before(list.Data.LastModified) {
	// 		clientState.SyncStatus = "out_of_sync"
	// 	} else {
	// 		clientState.SyncStatus = "in_sync"
	// 	}
	//
	// 	status.ClientStates[state.ClientID] = clientState
	// }

	// return status, nil
	// TODO: Implement and test list sync status
	return nil, nil

}

// func (s *userListService[T]) Delete(ctx context.Context, id uint64) error {
// 	log := utils.LoggerFromContext(ctx)
// 	log.Debug().
// 		Uint64("id", id).
// 		Msg("Deleting list")
//
// 	// First verify this is a list and the user has permission
// 	list, err := s.GetByID(ctx, id)
// 	if err != nil {
// 		return fmt.Errorf("failed to delete list: %w", err)
// 	}
//
// 	// TODO: Verify the current user has permission to delete this list
// 	// This would check if the current user ID matches list.Data.ItemList.OwnerID
//
// 	// Use the user service
// 	err = s.userItemRepo.Delete(ctx, id)
// 	if err != nil {
// 		log.Error().Err(err).
// 			Uint64("id", id).
// 			Msg("Failed to delete list")
// 		return fmt.Errorf("failed to delete list: %w", err)
// 	}
//
// 	log.Info().
// 		Uint64("id", id).
// 		Str("title", list.Title).
// 		Msg("list deleted successfully")
//
// 	return nil
// }

func (s *userListService[T]) ReorderItems(ctx context.Context, listID uint64, itemIDs []uint64) error {
	return s.coreListService.ReorderItems(ctx, listID, itemIDs)
}

// Helper functions

// haslistReadPermission checks if the user has read permission for a list
func (s *userListService[T]) haslistReadPermission(ctx context.Context, userID uint64, list *models.MediaItem[T]) bool {
	// The owner always has read permission
	if list.OwnerID == userID {
		return true
	}
	log := utils.LoggerFromContext(ctx)

	itemList := list.GetData().GetItemList()
	// Check if the list is shared with this user
	for _, collab := range itemList.SharedWith {
		log.Info().
			Uint64("ownerID", list.OwnerID).
			Uint64("requestingUserID", userID).
			Int64("sharedWith", collab).
			Msg("User attempting to view collaborators without permission")
		// if collab.UserID == userID {
		// Any permission level (read or write) allows reading
		// return true
		// }
	}

	// No permission found
	return false
}

// haslistWritePermission checks if the user has write permission for a list
func (s *userListService[T]) hasWritePermission(ctx context.Context, userID uint64, list *models.MediaItem[T]) bool {
	// The owner always has write permission
	if list.OwnerID == userID {
		return true
	}

	// itemList := list.GetData().GetItemList()
	// Check if the list is shared with this user with write permission
	// for _, collab := range itemList.SharedWith {
	// if collab.UserID == userID && collab.PermissionLevel == "write" {
	// 	return true
	// }
	// }

	// No write permission found
	return false
}

// func (s *userListService[T]) GetCollaborators(userID string, listID string) ([]UserCollaboration, error) {
// 	var collabs []UserCollaboration
// 	// note yet implemented
//
// 	return collabs, nil
//
// }

func createList[T mediatypes.ListData](name string, description string, criteria map[string]interface{}, userID uint64) T {
	var list T
	now := time.Now()
	list.SetItemList(mediatypes.ItemList{
		Details: mediatypes.MediaDetails{
			Title:       name,
			Description: description,
			AddedAt:     now,
		},
		OwnerID:    userID,
		ModifiedBy: userID,
		Items:      []mediatypes.ListItem{},
		ItemCount:  0,
		// Smart list specific fields
		IsSmart:        true,
		SmartCriteria:  criteria,
		AutoUpdateTime: now,
	})
	return list
}
