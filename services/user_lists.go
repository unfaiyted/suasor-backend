package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	mediatypes "suasor/clients/media/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/utils/logger"
)

// UserlistService defines the interface for user-owned list operations
// This service extends listService with operations specific to user-owned lists
type UserListService[T mediatypes.ListData] interface {
	// Include all core list service methods
	CoreListService[T]

	Create(ctx context.Context, userID uint64, list *models.MediaItem[T]) (*models.MediaItem[T], error)
	Update(ctx context.Context, userID uint64, list *models.MediaItem[T]) (*models.MediaItem[T], error)

	AddItem(ctx context.Context, userID uint64, listID uint64, itemID uint64) error
	RemoveItem(ctx context.Context, userID uint64, listID uint64, itemID uint64) error
	RemoveItemAtPosition(ctx context.Context, userID uint64, listID uint64, itemID uint64, position int) error
	ReorderItems(ctx context.Context, userID uint64, listID uint64, itemIDs []uint64) error
	UpdateItems(ctx context.Context, userID uint64, listID uint64, items []*models.MediaItem[T]) error

	Delete(ctx context.Context, userID uint64, id uint64) error

	GetFavorite(ctx context.Context, userID uint64, limit int, offset int) ([]*models.MediaItem[T], error)
	GetRecentByUser(ctx context.Context, userID uint64, days int, limit int) ([]*models.MediaItem[T], error)

	// Smart list operations
	CreateSmartList(ctx context.Context,
		userID uint64, name string,
		description string,
		criteria map[string]interface{},
	) (*models.MediaItem[T], error)
	UpdateSmartCriteria(ctx context.Context, userID uint64, listID uint64, criteria map[string]interface{}) (*models.MediaItem[T], error)
	RefreshSmartList(ctx context.Context, userID uint64, listID uint64) (*models.MediaItem[T], error)

	// list sharing and collaboration
	ShareWithUser(ctx context.Context, userID uint64, listID uint64, targetUserID uint64, permissionLevel string) error
	GetShared(ctx context.Context, userID uint64) ([]*models.MediaItem[T], error)
	GetCollaborators(ctx context.Context, userID uint64, listID uint64) ([]*models.ListCollaborator, error)
	RemoveCollaborator(ctx context.Context, userID uint64, listID uint64, targetUserID uint64) error
}

type userListService[T mediatypes.ListData] struct {
	CoreListService[T]

	userRepo     repository.UserRepository
	listRepo     repository.CoreListRepository[T]
	userItemRepo repository.UserMediaItemRepository[T]
	userDataRepo repository.UserMediaItemDataRepository[T]
}

// NewUserlistService creates a new user list service
func NewUserListService[T mediatypes.ListData](
	coreListService CoreListService[T],
	userRepo repository.UserRepository,
	listRepo repository.CoreListRepository[T],
	userItemRepo repository.UserMediaItemRepository[T],
	userDataRepo repository.UserMediaItemDataRepository[T],
) UserListService[T] {
	return &userListService[T]{
		CoreListService: coreListService,
		userRepo:        userRepo,
		listRepo:        listRepo,
		userItemRepo:    userItemRepo,
		userDataRepo:    userDataRepo,
	}
}

func (s *userListService[T]) Create(ctx context.Context, userID uint64, list *models.MediaItem[T]) (*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Str("title", list.Title).
		Msg("Creating user list")

	// Set ownership info if not already set
	itemList := list.GetData().GetItemList()

	if itemList.OwnerID == 0 {
		itemList.OwnerID = userID
	}
	if itemList.ModifiedBy == 0 {
		itemList.ModifiedBy = userID
	}

	// Ensure list has a valid name
	if list.Title == "" || itemList.Details.Title == "" {
		return nil, errors.New("list must have a title")
	}

	// Initialize items array if nil
	if itemList.Items == nil {
		itemList.Items = []mediatypes.ListItem{}
	}

	// Set creation time for LastModified
	itemList.Details.AddedAt = time.Now()
	itemList.Details.UpdatedAt = time.Now()
	itemList.LastModified = time.Now()

	// Initialize ItemCount
	itemList.ItemCount = len(itemList.Items)

	// Set owner if not set
	if itemList.OwnerID == 0 && list.OwnerID != 0 {
		itemList.OwnerID = list.OwnerID
	}

	// Set title at MediaItem level to match the Data.ItemList.Details.Title
	if list.Title == "" && itemList.Details.Title != "" {
		list.Title = itemList.Details.Title
	}

	// Use the underlying repository directly for better control over validation
	result, err := s.userItemRepo.Create(ctx, list)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create list")
		return nil, fmt.Errorf("failed to create list: %w", err)
	}

	log.Info().
		Uint64("id", result.ID).
		Str("title", result.Title).
		Msg("list created successfully")

	return result, nil
}
func (s *userListService[T]) Update(ctx context.Context, userID uint64, list *models.MediaItem[T]) (*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("id", list.ID).
		Str("title", list.Title).
		Msg("Updating user list")

	// Verify user has permission to update this list
	existing, err := s.GetByID(ctx, list.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update user list: %w", err)
	}
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to update user list: %w", err)
	}

	// Check ownership or collaboration permission
	if !s.hasListWritePermission(ctx, user, existing) {
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

	list.GetData().SetItemList(*itemList)

	// Ensure the list exists
	existingItemList := existing.GetData().GetItemList()

	// Preserve items if not provided in the update
	if itemList.Items == nil || len(itemList.Items) == 0 {
		itemList.Items = existingItemList.Items
	}

	// Preserve sync client states if not provided
	list.SyncClients = existing.SyncClients

	// Update last modified time
	itemList.LastModified = time.Now()

	// Update ItemCount
	itemList.ItemCount = len(itemList.Items)

	// Ensure positions are normalized
	itemList.NormalizePositions()

	// Set title at MediaItem level to match the Data.ItemList.Details.Title
	if list.Title != itemList.Details.Title {
		list.Title = itemList.Details.Title
	}

	// Run validation to check for issues
	issues := itemList.ValidateItems()
	if len(issues) > 0 {
		// Log the issues but continue with the update
		for _, issue := range issues {
			log.Warn().Str("issue", issue).Msg("list validation issue")
		}
	}

	// Update using the user service for consistent behavior
	result, err := s.userItemRepo.Update(ctx, list)
	if err != nil {
		log.Error().Err(err).
			Uint64("id", list.ID).
			Msg("Failed to update list")
		return nil, fmt.Errorf("failed to update list: %w", err)
	}

	log.Info().
		Uint64("id", result.ID).
		Str("title", result.Title).
		Msg("list updated successfully")

	return result, nil
}
func (s *userListService[T]) Delete(ctx context.Context, userID uint64, listID uint64) error {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("id", listID).
		Msg("Deleting user list")

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user list: %w", err)
	}
	isAdmin := user.Role == "admin"

	// Verify user has permission to delete this list
	list, err := s.GetByID(ctx, listID)
	if err != nil {
		return fmt.Errorf("failed to delete user list: %w", err)
	}

	// Only the owner can delete a list
	if list.OwnerID != userID && !isAdmin {
		log.Warn().
			Uint64("listID", listID).
			Uint64("ownerID", list.OwnerID).
			Uint64("requestingUserID", userID).
			Msg("User attempting to delete a list they don't own")
		return errors.New("only the owner can delete a list")
	}

	err = s.userItemRepo.Delete(ctx, listID)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Msg("Failed to delete list")
		return fmt.Errorf("failed to delete list: %w", err)
	}

	log.Info().
		Uint64("id", listID).
		Msg("List deleted successfully")

	return nil
}

func (s *userListService[T]) AddItem(ctx context.Context, userID uint64, listID uint64, itemID uint64) error {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("listID", listID).
		Uint64("itemID", itemID).
		Msg("Adding item to user list")

	// Verify user has permission to modify this list
	list, err := s.GetByID(ctx, listID)
	if err != nil {
		return fmt.Errorf("failed to add item to user list: %w", err)
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to add item to user list: %w", err)
	}

	// Check if user has write permission
	if !s.hasListWritePermission(ctx, user, list) {
		log.Warn().
			Uint64("listID", listID).
			Uint64("ownerID", list.OwnerID).
			Uint64("requestingUserID", userID).
			Msg("User attempting to modify a list without permission")
		return errors.New("you don't have permission to modify this list")
	}

	itemList := list.GetData().GetItemList()

	// Create a ListItem from the media item
	newItem := mediatypes.ListItem{
		ItemID:        itemID,
		Position:      len(itemList.Items),
		LastChanged:   time.Now(),
		ChangeHistory: []mediatypes.ChangeRecord{},
	}

	// Add the item using the built-in AddItem method
	// 0 indicates application level modification
	itemList.AddItem(newItem)

	// Store the update
	_, err = s.Update(ctx, userID, list)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Uint64("itemID", itemID).
			Msg("Failed to add item to list")
		return fmt.Errorf("failed to update list after adding item: %w", err)
	}

	log.Info().
		Uint64("listID", listID).
		Uint64("itemID", itemID).
		Msg("Item added to list successfully")

	return nil
}
func (s *userListService[T]) Search(ctx context.Context, query mediatypes.QueryOptions) ([]*models.MediaItem[T], error) {
	return s.Search(ctx, query)
}

func (s *userListService[T]) UpdateSmartCriteria(ctx context.Context, userID uint64, listID uint64, criteria map[string]interface{}) (*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("listID", listID).
		Interface("criteria", criteria).
		Msg("Updating smart list criteria")

	return nil, nil
	// return s.UpdateSmartCriteria(ctx, listID, criteria)
}

// User-specific operations

// GetUser lists retrieves lists owned by a specific user with pagination
func (s *userListService[T]) GetUser(ctx context.Context, userID uint64, limit int, offset int) ([]*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)
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
	log := logger.LoggerFromContext(ctx)
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
	log := logger.LoggerFromContext(ctx)
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
	log := logger.LoggerFromContext(ctx)
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

// CreateSmartlist creates a list that updates automatically based on criteria
func (s *userListService[T]) CreateSmartList(ctx context.Context, userID uint64, name string, description string, criteria map[string]interface{}) (*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)
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
	result, err := s.Create(ctx, userID, list)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Str("name", name).
			Msg("Failed to create smart list")
		return nil, fmt.Errorf("failed to create smart list: %w", err)
	}

	// Refresh the smart list to populate it based on criteria
	refreshed, err := s.RefreshSmartList(ctx, userID, result.ID)
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

// RefreshSmartlist updates a smart list based on its criteria
func (s *userListService[T]) RefreshSmartList(ctx context.Context, userID uint64, listID uint64) (*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)
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
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh smart list: %w", err)
	}

	// Verify user has permission to read this list (refresh is a read operation)
	if !s.hasListReadPermission(ctx, user, list) {
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
	updated, err := s.Update(ctx, userID, list)
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

// List sharing and collaboration

// SharelistWithUser shares a list with another user
func (s *userListService[T]) ShareWithUser(ctx context.Context, userID uint64, listID uint64, targetUserID uint64, permissionLevel string) error {
	log := logger.LoggerFromContext(ctx)
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
	_, err = s.Update(ctx, userID, list)
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
	log := logger.LoggerFromContext(ctx)
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

func (s *userListService[T]) GetCollaborators(ctx context.Context, userID uint64, listID uint64) ([]*models.ListCollaborator, error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("listID", listID).
		Msg("Getting list collaborators")

	// Get the list
	list, err := s.GetByID(ctx, listID)
	if err != nil {
		return nil, fmt.Errorf("failed to get list collaborators: %w", err)
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get list collaborators: %w", err)
	}

	if !s.hasListReadPermission(ctx, user, list) {
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
		return []*models.ListCollaborator{}, nil
	}

	collaborators, err := s.listRepo.GetCollaborators(ctx, listID)
	if err != nil {
		return nil, fmt.Errorf("failed to get list collaborators: %w", err)
	}

	return collaborators, nil
}
func (s *userListService[T]) RemoveCollaborator(ctx context.Context, userID uint64, listID uint64, collaboratorID uint64) error {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("listID", listID).
		Uint64("collaboratorID", collaboratorID).
		Msg("Removing list collaborator")

	// Get the list
	list, err := s.GetByID(ctx, listID)
	if err != nil {
		return fmt.Errorf("failed to remove list collaborator: %w", err)
	}
	if err != nil {
		return fmt.Errorf("failed to remove list collaborator: %w", err)
	}

	// Verify user has permission to modify sharing (only owner can do this)
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
	_, err = s.Update(ctx, userID, list)
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

func (s userListService[T]) RemoveItem(ctx context.Context, userID uint64, listID uint64, itemID uint64) error {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("listID", listID).
		Uint64("itemID", itemID).
		Msg("Removing item from list")

	// Get the list
	list, err := s.GetByID(ctx, listID)
	if err != nil {
		return fmt.Errorf("failed to remove item from list: %w", err)
	}
	itemList := list.GetData().GetItemList()

	// Use the RemoveItem method provided by ItemList
	// 0 indicates application level modification
	err = itemList.RemoveItem(itemID, 0)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Uint64("itemID", itemID).
			Msg("Failed to remove item from list")
		return err
	}

	// Store the update
	_, err = s.Update(ctx, userID, list)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Uint64("itemID", itemID).
			Msg("Failed to update list after removing item")
		return fmt.Errorf("failed to update list after removing item: %w", err)
	}

	log.Info().
		Uint64("listID", listID).
		Uint64("itemID", itemID).
		Msg("Item removed from list successfully")

	return nil
}
func (s *userListService[T]) ReorderItems(ctx context.Context, userID uint64, listID uint64, itemIDs []uint64) error {
	log := logger.LoggerFromContext(ctx)

	log.Debug().
		Uint64("listID", listID).
		Interface("itemIDs", itemIDs).
		Msg("Reordering list items")

	// Get the list
	list, err := s.GetByID(ctx, listID)
	if err != nil {
		return fmt.Errorf("failed to reorder list items: %w", err)
	}
	itemList := list.GetData().GetItemList()

	// Verify that the number of items matches
	if len(itemIDs) != len(itemList.Items) {
		log.Error().
			Int("providedCount", len(itemIDs)).
			Int("actualCount", len(itemList.Items)).
			Msg("Reorder operation must include all list items")
		return errors.New("reorder operation must include all list items")
	}

	// Create a new ordered list of items
	newOrder := make([]mediatypes.ListItem, len(itemIDs))
	tempItems := make(map[uint64]mediatypes.ListItem)

	// Create a map of existing items for quick lookup
	for _, item := range itemList.Items {
		tempItems[item.ItemID] = item
	}

	// First verify all items exist
	missingItems := []uint64{}
	for _, id := range itemIDs {
		// id, err := strconv.ParseUint(idStr, 10, 64)
		// if err != nil {
		// 	log.Error().
		// 		Str("itemID", idStr).
		// 		Msg("Invalid item ID format")
		// 	return fmt.Errorf("invalid item ID format: %s", idStr)
		// }

		if _, exists := tempItems[id]; !exists {
			missingItems = append(missingItems, id)
		}
	}

	if len(missingItems) > 0 {
		log.Error().
			Interface("missingItems", missingItems).
			Msg("Items not found in list")
		return fmt.Errorf("items not found in list: %v", missingItems)
	}

	// Now build the new order
	for i, id := range itemIDs {
		// id, _ := strconv.ParseUint(id, 10, 64)
		item := tempItems[id]

		// Update position
		item.Position = i

		// Add change record
		item.AddChangeRecord(0, "reorder") // 0 indicates application level change

		newOrder[i] = item
	}

	// Update the list with the new item order
	itemList.Items = newOrder
	itemList.LastModified = time.Now()
	itemList.ModifiedBy = 0 // 0 indicates application level modification

	// Normalize positions to ensure they're sequential
	itemList.NormalizePositions()

	_, err = s.Update(ctx, userID, list)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Msg("Failed to update list after reordering items")
		return fmt.Errorf("failed to update list after reordering: %w", err)
	}

	log.Info().
		Uint64("listID", listID).
		Msg("list items reordered successfully")

	return nil
}
func (s *userListService[T]) UpdateItems(ctx context.Context, userID uint64, listID uint64, items []*models.MediaItem[T]) error {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("listID", listID).
		Int("itemCount", len(items)).
		Msg("Updating list items")

	// Get the list
	list, err := s.GetByID(ctx, listID)
	if err != nil {
		return fmt.Errorf("failed to update list items: %w", err)
	}

	// Convert MediaItems to ListItems
	listItems := make([]mediatypes.ListItem, len(items))
	now := time.Now()

	for i, item := range items {
		listItems[i] = mediatypes.ListItem{
			ItemID:      item.ID,
			Position:    i,
			LastChanged: now,
			ChangeHistory: []mediatypes.ChangeRecord{
				{
					ClientID:   0, // 0 indicates application level change
					ItemID:     fmt.Sprintf("%d", item.ID),
					ChangeType: "update",
					Timestamp:  now,
				},
			},
		}
	}

	itemList := list.GetData().GetItemList()
	// Replace all items
	itemList.Items = listItems
	itemList.LastModified = now
	itemList.ModifiedBy = 0 // 0 indicates application level modification
	itemList.ItemCount = len(listItems)

	// Ensure positions are normalized
	itemList.NormalizePositions()

	// Update the list
	_, err = s.userItemRepo.Update(ctx, list)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Msg("Failed to update list items")
		return fmt.Errorf("failed to update list items: %w", err)
	}

	log.Info().
		Uint64("listID", listID).
		Int("itemCount", len(items)).
		Msg("list items updated successfully")

	return nil
}

// Helper functions
// haslistReadPermission checks if the user has read permission for a list
func (s *userListService[T]) hasListReadPermission(ctx context.Context, user *models.User, list *models.MediaItem[T]) bool {
	log := logger.LoggerFromContext(ctx)
	// The owner always has read permission
	if list.OwnerID == user.ID || user.Role == "admin" {
		return true
	}

	itemList := list.GetData().GetItemList()
	// Check if the list is shared with this user
	for _, collabID := range itemList.SharedWith {
		log.Info().
			Uint64("ownerID", list.OwnerID).
			Uint64("requestingUserID", user.ID).
			Uint64("sharedWith", collabID).
			Msg("User attempting to read list without permission")
		if collabID == user.ID {
			// Any permission level allows reading
			return true
		}
	}

	// No permission found
	return false
}

// haslistWritePermission checks if the user has write permission for a list
func (s *userListService[T]) hasListWritePermission(ctx context.Context, user *models.User, list *models.MediaItem[T]) bool {
	log := logger.LoggerFromContext(ctx)

	// The owner and admin always has write permission
	if list.OwnerID == user.ID || user.Role == "admin" {
		return true
	}

	itemList := list.GetData().GetItemList()
	// Check if the list is shared with this user with write permission
	for _, collabID := range itemList.SharedWith {
		if collabID == user.ID {
			// Check if the user has write permission
			collab, err := s.listRepo.GetCollaborator(ctx, list.ID, user.ID)
			if err != nil {
				log.Error().Err(err).Msg("Failed to get list collaborator")
				return false
			}
			if collab.Permission == models.CollaboratorPermissionWrite {
				return true
			}
		}
	}

	return false
}

func createList[T mediatypes.ListData](name string, description string, criteria map[string]interface{}, userID uint64) T {
	var list T
	now := time.Now()
	list.SetItemList(mediatypes.ItemList{
		Details: &mediatypes.MediaDetails{
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

func (s *userListService[T]) RemoveItemAtPosition(ctx context.Context, userID uint64, listID uint64, itemID uint64, position int) error {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("listID", listID).
		Uint64("itemID", itemID).
		Int("position", position).
		Msg("Removing item from list")

	// Verify user has permission to modify this list
	list, err := s.GetByID(ctx, listID)
	if err != nil {
		return fmt.Errorf("failed to remove item from list: %w", err)
	}
	itemList := list.GetData().GetItemList()

	// Use the RemoveItem method provided by ItemList
	// 0 indicates application level modification
	err = itemList.RemoveItemAtPosition(itemID, position, 0)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Uint64("itemID", itemID).
			Msg("Failed to remove item from list")
		return err
	}

	// Store the update
	_, err = s.Update(ctx, userID, list)
	if err != nil {
		log.Error().Err(err).
			Uint64("listID", listID).
			Uint64("itemID", itemID).
			Msg("Failed to update list after removing item")
		return fmt.Errorf("failed to update list after removing item: %w", err)
	}

	log.Info().
		Uint64("listID", listID).
		Uint64("itemID", itemID).
		Msg("Item removed from list successfully")

	return nil
}
