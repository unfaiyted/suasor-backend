package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"suasor/app/container"
	mediatypes "suasor/client/media/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/utils"
)

// UserPlaylistService defines the interface for user-owned playlist operations
// This service extends PlaylistService with operations specific to user-owned playlists
type UserListService[T mediatypes.ListData] interface {
	// Include all core playlist service methods
	CoreListService[T]

	// User-specific operations
	GetFavorite(ctx context.Context, userID uint64) ([]*models.MediaItem[*mediatypes.Playlist], error)
	GetRecentByUser(ctx context.Context, userID uint64, days int, limit int) ([]*models.MediaItem[*mediatypes.Playlist], error)

	// Smart playlist operations
	CreateSmartList(ctx context.Context,
		userID uint64, name string,
		description string,
		criteria map[string]interface{},
	) (*models.MediaItem[*mediatypes.Playlist], error)
	UpdateSmartCriteria(ctx context.Context, playlistID uint64, criteria map[string]interface{}) (*models.MediaItem[*mediatypes.Playlist], error)
	RefreshSmartList(ctx context.Context, playlistID uint64) (*models.MediaItem[*mediatypes.Playlist], error)

	// Playlist sharing and collaboration
	ShareWithUser(ctx context.Context, playlistID uint64, targetUserID uint64, permissionLevel string) error
	GetShared(ctx context.Context, userID uint64) ([]*models.MediaItem[*mediatypes.Playlist], error)
	GetCollaborators(ctx context.Context, playlistID uint64) ([]models.ListCollaborator, error)
	RemoveCollaborator(ctx context.Context, playlistID uint64, userID uint64) error

	// Playlist sync
	SyncToClients(ctx context.Context, playlistID uint64, clientIDs []uint64) error
	GetSyncStatus(ctx context.Context, playlistID uint64) (*models.PlaylistSyncStatus, error)
}

type userPlaylistService struct {
	userItemRepo repository.UserMediaItemRepository[*mediatypes.Playlist]
	userDataRepo repository.UserMediaItemDataRepository[*mediatypes.Playlist]
	itemRepo     repository.MediaItemRepository[*mediatypes.Playlist]
}

// NewUserPlaylistService creates a new user playlist service
func NewUserPlaylistService(
	ctx context.Context,
	c *container.Container,
) UserPlaylistService {
	return &userPlaylistService{
		userItemRepo: container.MustGet[repository.UserMediaItemRepository[*mediatypes.Playlist]](c),
		userDataRepo: container.MustGet[repository.UserMediaItemDataRepository[*mediatypes.Playlist]](c),
		itemRepo:     container.MustGet[repository.MediaItemRepository[*mediatypes.Playlist]](c),
	}
}

// Core methods delegated to the base PlaylistService

func (s *userPlaylistService) Create(ctx context.Context, playlist *models.MediaItem[*mediatypes.Playlist]) (*models.MediaItem[*mediatypes.Playlist], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Str("title", playlist.Title).
		Msg("Creating user playlist")

	// Get user ID from context
	userID := ctx.Value("userID").(uint64)

	// Set ownership info if not already set
	if playlist.Data.OwnerID == 0 {
		playlist.Data.OwnerID = userID
	}
	if playlist.Data.ItemList.OwnerID == 0 {
		playlist.Data.ItemList.OwnerID = userID
	}
	if playlist.Data.ItemList.ModifiedBy == 0 {
		playlist.Data.ItemList.ModifiedBy = userID
	}

	// Delegate to core service
	return s.Create(ctx, playlist)
}

func (s *userPlaylistService) Update(ctx context.Context, playlist *models.MediaItem[*mediatypes.Playlist]) (*models.MediaItem[*mediatypes.Playlist], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("id", playlist.ID).
		Str("title", playlist.Title).
		Msg("Updating user playlist")

	// Verify user has permission to update this playlist
	userID := ctx.Value("userID").(uint64)
	existing, err := s.GetByID(ctx, playlist.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update user playlist: %w", err)
	}

	// Check ownership or collaboration permission
	if !s.hasWritePermission(ctx, userID, existing) {
		log.Warn().
			Uint64("playlistID", playlist.ID).
			Uint64("ownerID", existing.Data.OwnerID).
			Uint64("requestingUserID", userID).
			Msg("User attempting to update a playlist without permission")
		return nil, errors.New("you don't have permission to update this playlist")
	}

	// Set the modified by field to the current user
	playlist.Data.ItemList.ModifiedBy = userID
	playlist.Data.ItemList.LastModified = time.Now()

	// Delegate to core service
	return s.Update(ctx, playlist)
}

func (s *userPlaylistService) GetByID(ctx context.Context, id uint64) (*models.MediaItem[*mediatypes.Playlist], error) {
	return s.GetByID(ctx, id)
}

func (s *userPlaylistService) GetByUserID(ctx context.Context, userID uint64) ([]*models.MediaItem[*mediatypes.Playlist], error) {
	return s.GetByUserID(ctx, userID)
}

func (s *userPlaylistService) Delete(ctx context.Context, id uint64) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("id", id).
		Msg("Deleting user playlist")

	// Verify user has permission to delete this playlist
	userID := ctx.Value("userID").(uint64)
	playlist, err := s.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete user playlist: %w", err)
	}

	// Only the owner can delete a playlist
	if playlist.Data.OwnerID != userID {
		log.Warn().
			Uint64("playlistID", id).
			Uint64("ownerID", playlist.Data.OwnerID).
			Uint64("requestingUserID", userID).
			Msg("User attempting to delete a playlist they don't own")
		return errors.New("only the owner can delete a playlist")
	}

	// Delegate to core service for deletion
	return s.Delete(ctx, id)
}

func (s *userPlaylistService) GetPlaylistItems(ctx context.Context, playlistID uint64) ([]*models.MediaItem[mediatypes.MediaData], error) {
	return s.GetPlaylistItems(ctx, playlistID)
}

func (s *userPlaylistService) AddItemToPlaylist(ctx context.Context, playlistID uint64, itemID uint64) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("playlistID", playlistID).
		Uint64("itemID", itemID).
		Msg("Adding item to user playlist")

	// Verify user has permission to modify this playlist
	userID := ctx.Value("userID").(uint64)
	playlist, err := s.GetByID(ctx, playlistID)
	if err != nil {
		return fmt.Errorf("failed to add item to user playlist: %w", err)
	}

	// Check if user has write permission
	if !s.hasWritePermission(ctx, userID, playlist) {
		log.Warn().
			Uint64("playlistID", playlistID).
			Uint64("ownerID", playlist.Data.OwnerID).
			Uint64("requestingUserID", userID).
			Msg("User attempting to modify a playlist without permission")
		return errors.New("you don't have permission to modify this playlist")
	}

	// Delegate to core service
	return s.AddItemToPlaylist(ctx, playlistID, itemID)
}

func (s *userPlaylistService) RemoveItemFromPlaylist(ctx context.Context, playlistID uint64, itemID uint64) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("playlistID", playlistID).
		Uint64("itemID", itemID).
		Msg("Removing item from user playlist")

	// Verify user has permission to modify this playlist
	userID := ctx.Value("userID").(uint64)
	playlist, err := s.GetByID(ctx, playlistID)
	if err != nil {
		return fmt.Errorf("failed to remove item from user playlist: %w", err)
	}

	// Check if user has write permission
	if !s.hasWritePermission(ctx, userID, playlist) {
		log.Warn().
			Uint64("playlistID", playlistID).
			Uint64("ownerID", playlist.Data.OwnerID).
			Uint64("requestingUserID", userID).
			Msg("User attempting to modify a playlist without permission")
		return errors.New("you don't have permission to modify this playlist")
	}

	// Delegate to core service
	return s.RemoveItemFromPlaylist(ctx, playlistID, itemID)
}

func (s *userPlaylistService) ReorderPlaylistItems(ctx context.Context, playlistID uint64, itemIDs []string) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("playlistID", playlistID).
		Interface("itemIDs", itemIDs).
		Msg("Reordering user playlist items")

	// Verify user has permission to modify this playlist
	userID := ctx.Value("userID").(uint64)
	playlist, err := s.GetByID(ctx, playlistID)
	if err != nil {
		return fmt.Errorf("failed to reorder user playlist items: %w", err)
	}

	// Check if user has write permission
	if !s.hasWritePermission(ctx, userID, playlist) {
		log.Warn().
			Uint64("playlistID", playlistID).
			Uint64("ownerID", playlist.Data.OwnerID).
			Uint64("requestingUserID", userID).
			Msg("User attempting to modify a playlist without permission")
		return errors.New("you don't have permission to modify this playlist")
	}

	// Delegate to core service
	return s.ReorderItems(ctx, playlistID, itemIDs)
}

func (s *userPlaylistService) UpdateItems(ctx context.Context, playlistID uint64, items []*models.MediaItem[mediatypes.MediaData]) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("playlistID", playlistID).
		Int("itemCount", len(items)).
		Msg("Updating user playlist items")

	// Verify user has permission to modify this playlist
	userID := ctx.Value("userID").(uint64)
	playlist, err := s.GetByID(ctx, playlistID)
	if err != nil {
		return fmt.Errorf("failed to update user playlist items: %w", err)
	}

	// Check if user has write permission
	if !s.hasWritePermission(ctx, userID, playlist) {
		log.Warn().
			Uint64("playlistID", playlistID).
			Uint64("ownerID", playlist.Data.OwnerID).
			Uint64("requestingUserID", userID).
			Msg("User attempting to modify a playlist without permission")
		return errors.New("you don't have permission to modify this playlist")
	}

	// Delegate to core service
	return s.UpdateItems(ctx, playlistID, items)
}

func (s *userPlaylistService) Search(ctx context.Context, query mediatypes.QueryOptions) ([]*models.MediaItem[*mediatypes.Playlist], error) {
	return s.Search(ctx, query)
}

func (s *userPlaylistService) GetRecent(ctx context.Context, days int, limit int) ([]*models.MediaItem[*mediatypes.Playlist], error) {
	return s.GetRecent(ctx, days, limit)
}

func (s *userPlaylistService) Sync(ctx context.Context, playlistID uint64, targetClientIDs []uint64) error {
	// Verify user has permission to synchronize this playlist
	userID := ctx.Value("userID").(uint64)
	playlist, err := s.GetByID(ctx, playlistID)
	if err != nil {
		return fmt.Errorf("failed to sync user playlist: %w", err)
	}

	// Check if user has read permission (sync is a read operation followed by creation elsewhere)
	if !s.hasPlaylistReadPermission(ctx, userID, playlist) {
		return errors.New("you don't have permission to sync this playlist")
	}

	// Delegate to core service
	return s.Sync(ctx, playlistID, targetClientIDs)
}

// User-specific operations

// GetUserPlaylists retrieves playlists owned by a specific user with pagination
func (s *userPlaylistService) GetUser(ctx context.Context, userID uint64, limit int, offset int) ([]*models.MediaItem[*mediatypes.Playlist], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Int("limit", limit).
		Int("offset", offset).
		Msg("Getting user playlists with pagination")

	// Get all playlists for this user
	playlists, err := s.itemRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to get user playlists")
		return nil, fmt.Errorf("failed to get user playlists: %w", err)
	}

	log.Info().
		Uint64("userID", userID).
		Int("count", len(playlists)).
		Msg("Retrieved user playlists")

	return playlists, nil
}

// SearchUserPlaylists searches for playlists owned by a specific user
func (s *userPlaylistService) SearchUserPlaylists(ctx context.Context, userID uint64, query string) ([]*models.MediaItem[*mediatypes.Playlist], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Str("query", query).
		Msg("Searching user playlists")

	// Create query options with user filter
	options := mediatypes.QueryOptions{
		MediaType: mediatypes.MediaTypePlaylist,
		Query:     query,
		OwnerID:   userID,
	}

	// Delegate to core service with owner filter
	return s.Search(ctx, options)
}

// GetRecentUserPlaylists retrieves recently updated playlists for a user
func (s *userPlaylistService) GetRecentByUser(ctx context.Context, userID uint64, days int, limit int) ([]*models.MediaItem[*mediatypes.Playlist], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Int("limit", limit).
		Msg("Getting recent user playlists")

	options := mediatypes.QueryOptions{
		MediaType: mediatypes.MediaTypePlaylist,
		OwnerID:   userID,
		Limit:     limit,
		Sort:      "updated_at",
		SortOrder: "desc",
	}
	return s.itemRepo.Search(ctx, options)

}

// GetFavoritePlaylists retrieves playlists marked as favorite by the user
func (s *userPlaylistService) GetFavoritePlaylists(ctx context.Context, userID uint64) ([]*models.MediaItem[*mediatypes.Playlist], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Msg("Getting favorite playlists")

	limit := 0
	offset := 0

	// Use user data repository to get all favorites of type playlist
	userFavoritePlayData, err := s.userDataRepo.GetFavorites(ctx, userID, limit, offset)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to get favorite playlist IDs")
		return nil, fmt.Errorf("failed to get favorite playlists: %w", err)
	}

	var ids []uint64
	for _, data := range userFavoritePlayData {
		ids = append(ids, data.MediaItemID)
	}

	// Fetch the playlists by IDs
	playlists, err := s.itemRepo.GetByIDs(ctx, ids)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to get favorite playlists")
		return nil, fmt.Errorf("failed to get favorite playlists: %w", err)
	}

	log.Info().
		Uint64("userID", userID).
		Int("count", len(playlists)).
		Msg("Retrieved favorite playlists")

	return playlists, nil
}

// Smart playlist operations

// CreateSmartPlaylist creates a playlist that updates automatically based on criteria
func (s *userPlaylistService) CreateSmartList(ctx context.Context, userID uint64, name string, description string, criteria map[string]interface{}) (*models.MediaItem[*mediatypes.Playlist], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Str("name", name).
		Interface("criteria", criteria).
		Msg("Creating smart playlist")

	// Create a new playlist with smart flag enabled
	now := time.Now()
	playlist := &models.MediaItem[*mediatypes.Playlist]{
		Title: name,
		Type:  mediatypes.MediaTypePlaylist,
		Data: &mediatypes.Playlist{
			ItemList: mediatypes.ItemList{
				Details: mediatypes.MediaDetails{
					Title:       name,
					Description: description,
					AddedAt:     now,
				},
				OwnerID:    userID,
				ModifiedBy: userID,
				Items:      []mediatypes.ListItem{},
				ItemCount:  0,
				// Smart playlist specific fields
				IsSmart:        true,
				SmartCriteria:  criteria,
				AutoUpdateTime: now,
			},
		},
	}

	// Create the playlist
	result, err := s.Create(ctx, playlist)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Str("name", name).
			Msg("Failed to create smart playlist")
		return nil, fmt.Errorf("failed to create smart playlist: %w", err)
	}

	// Refresh the smart playlist to populate it based on criteria
	refreshed, err := s.RefreshSmartPlaylist(ctx, result.ID)
	if err != nil {
		log.Warn().Err(err).
			Uint64("playlistID", result.ID).
			Msg("Failed to initially populate smart playlist")
		// Return the playlist anyway, just warn about the population failure
	} else if refreshed != nil {
		result = refreshed
	}

	log.Info().
		Uint64("playlistID", result.ID).
		Str("name", name).
		Uint64("userID", userID).
		Int("itemCount", result.Data.ItemList.ItemCount).
		Msg("Smart playlist created successfully")

	return result, nil
}

// UpdateSmartPlaylistCriteria updates the criteria for a smart playlist
func (s *userPlaylistService) UpdateSmartPlaylistCriteria(ctx context.Context, playlistID uint64, criteria map[string]interface{}) (*models.MediaItem[*mediatypes.Playlist], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("playlistID", playlistID).
		Interface("criteria", criteria).
		Msg("Updating smart playlist criteria")

	// Get the playlist
	playlist, err := s.GetByID(ctx, playlistID)
	if err != nil {
		return nil, fmt.Errorf("failed to update smart playlist criteria: %w", err)
	}

	// Verify this is a smart playlist
	if !playlist.Data.IsSmart {
		log.Error().
			Uint64("playlistID", playlistID).
			Msg("Cannot update criteria of non-smart playlist")
		return nil, errors.New("cannot update criteria of non-smart playlist")
	}

	// Verify user has permission to modify this playlist
	userID := ctx.Value("userID").(uint64)
	if !s.hasWritePermission(ctx, userID, playlist) {
		log.Warn().
			Uint64("playlistID", playlistID).
			Uint64("ownerID", playlist.Data.OwnerID).
			Uint64("requestingUserID", userID).
			Msg("User attempting to modify a playlist without permission")
		return nil, errors.New("you don't have permission to modify this playlist")
	}

	// Update the criteria
	playlist.Data.SmartCriteria = criteria
	playlist.Data.AutoUpdateTime = time.Now()
	playlist.Data.LastModified = time.Now()
	playlist.Data.ModifiedBy = userID

	// Update the playlist
	updated, err := s.Update(ctx, playlist)
	if err != nil {
		log.Error().Err(err).
			Uint64("playlistID", playlistID).
			Msg("Failed to update smart playlist criteria")
		return nil, fmt.Errorf("failed to update smart playlist criteria: %w", err)
	}

	// Refresh the playlist to apply the new criteria
	refreshed, err := s.RefreshSmartPlaylist(ctx, playlistID)
	if err != nil {
		log.Warn().Err(err).
			Uint64("playlistID", playlistID).
			Msg("Failed to refresh smart playlist after criteria update")
		// Return the updated playlist anyway
		return updated, nil
	}

	log.Info().
		Uint64("playlistID", playlistID).
		Int("itemCount", refreshed.Data.ItemCount).
		Msg("Smart playlist criteria updated and refreshed successfully")

	return refreshed, nil
}

// RefreshSmartPlaylist updates a smart playlist based on its criteria
func (s *userPlaylistService) RefreshSmartPlaylist(ctx context.Context, playlistID uint64) (*models.MediaItem[*mediatypes.Playlist], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("playlistID", playlistID).
		Msg("Refreshing smart playlist")

	// Get the playlist
	playlist, err := s.GetByID(ctx, playlistID)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh smart playlist: %w", err)
	}

	// Verify this is a smart playlist
	if !playlist.Data.IsSmart {
		log.Error().
			Uint64("playlistID", playlistID).
			Msg("Cannot refresh non-smart playlist")
		return nil, errors.New("cannot refresh non-smart playlist")
	}

	// Verify user has permission to read this playlist (refresh is a read operation)
	userID := ctx.Value("userID").(uint64)
	if !s.hasPlaylistReadPermission(ctx, userID, playlist) {
		log.Warn().
			Uint64("playlistID", playlistID).
			Uint64("ownerID", playlist.Data.OwnerID).
			Uint64("requestingUserID", userID).
			Msg("User attempting to refresh a playlist without permission")
		return nil, errors.New("you don't have permission to refresh this playlist")
	}

	// Get the criteria
	// criteria := playlist.Data.SmartCriteria

	// In a real implementation, this would:
	// 1. Translate the criteria into a database query
	// 2. Execute the query to find all matching media items
	// 3. Replace the playlist items with the query results
	// 4. Update the playlist metadata

	// For now, just simulate the refresh by adding a note to the description
	now := time.Now()
	playlist.Data.Details.Description = fmt.Sprintf("%s\n\nLast refreshed: %s",
		playlist.Data.Details.Description, now.Format(time.RFC3339))
	playlist.Data.AutoUpdateTime = now
	playlist.Data.LastModified = now
	playlist.Data.ModifiedBy = userID

	// Update the playlist
	updated, err := s.Update(ctx, playlist)
	if err != nil {
		log.Error().Err(err).
			Uint64("playlistID", playlistID).
			Msg("Failed to update playlist after refresh")
		return nil, fmt.Errorf("failed to update playlist after refresh: %w", err)
	}

	log.Info().
		Uint64("playlistID", playlistID).
		Msg("Smart playlist refreshed successfully")

	return updated, nil
}

// Playlist sharing and collaboration

// SharePlaylistWithUser shares a playlist with another user
func (s *userPlaylistService) ShareWithUser(ctx context.Context, playlistID uint64, targetUserID uint64, permissionLevel string) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("playlistID", playlistID).
		Uint64("targetUserID", targetUserID).
		Str("permissionLevel", permissionLevel).
		Msg("Sharing playlist with user")

	// Get the playlist
	playlist, err := s.GetByID(ctx, playlistID)
	if err != nil {
		return fmt.Errorf("failed to share playlist: %w", err)
	}

	// Verify user has permission to share this playlist (only owner can share)
	userID := ctx.Value("userID").(uint64)
	if playlist.Data.OwnerID != userID {
		log.Warn().
			Uint64("playlistID", playlistID).
			Uint64("ownerID", playlist.Data.OwnerID).
			Uint64("requestingUserID", userID).
			Msg("User attempting to share a playlist they don't own")
		return errors.New("only the owner can share a playlist")
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
	// if playlist.Data.SharedWith == nil {
	// 	playlist.Data.SharedWith = []models.ListCollaborator{collaborator}
	// } else {
	// 	// Check if already shared
	// 	alreadyShared := false
	// 	for i, collab := range playlist.Data.SharedWith {
	// 		if collab.UserID == targetUserID {
	// 			alreadyShared = true
	// 			// Update permission level if it's different
	// 			if collab.PermissionLevel != permissionLevel {
	// 				playlist.Data.SharedWith[i].PermissionLevel = permissionLevel
	// 				playlist.Data.SharedWith[i].SharedAt = time.Now()
	// 				playlist.Data.SharedWith[i].SharedBy = userID
	// 			}
	// 			break
	// 		}
	// 	}
	//
	// 	if !alreadyShared {
	// 		playlist.Data.SharedWith = append(playlist.Data.SharedWith, collaborator)
	// 	}
	// }

	// Update the playlist
	_, err = s.Update(ctx, playlist)
	if err != nil {
		log.Error().Err(err).
			Uint64("playlistID", playlistID).
			Uint64("targetUserID", targetUserID).
			Msg("Failed to update playlist sharing information")
		return fmt.Errorf("failed to update playlist sharing information: %w", err)
	}

	log.Info().
		Uint64("playlistID", playlistID).
		Uint64("targetUserID", targetUserID).
		Str("permissionLevel", permissionLevel).
		Msg("Playlist shared successfully")

	return nil
}

// GetSharedPlaylists retrieves playlists shared with a user
func (s *userPlaylistService) GetSharedPlaylists(ctx context.Context, userID uint64) ([]*models.MediaItem[*mediatypes.Playlist], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Msg("Getting playlists shared with user")

	// TODO:
	// In a real implementation, this would query the database for playlists where
	// the user is in the SharedWith array
	// For now, we'll scan through all playlists to find ones shared with this user
	// allPlaylists, err := s.userRepo.GetAll(ctx, 1000, 0)
	// if err != nil {
	// 	log.Error().Err(err).
	// 		Uint64("userID", userID).
	// 		Msg("Failed to get all playlists")
	// 	return nil, fmt.Errorf("failed to get shared playlists: %w", err)
	// }
	//
	// var sharedPlaylists []*models.MediaItem[*mediatypes.Playlist]
	// for _, playlist := range allPlaylists {
	// 	// Skip playlists owned by this user (those are covered by GetUserPlaylists)
	// 	if playlist.Data.OwnerID == userID {
	// 		continue
	// 	}
	//
	// 	// Check if this playlist is shared with the user
	// 	for _, collab := range playlist.Data.SharedWith {
	// 		if collab.UserID == userID {
	// 			sharedPlaylists = append(sharedPlaylists, playlist)
	// 			break
	// 		}
	// 	}
	// }
	//
	// log.Info().
	// 	Uint64("userID", userID).
	// 	Int("count", len(sharedPlaylists)).
	// 	Msg("Retrieved playlists shared with user")
	//
	// return sharedPlaylists, nil
	return nil, nil
}

// GetPlaylistCollaborators retrieves the list of users a playlist is shared with
func (s *userPlaylistService) GetPlaylistCollaborators(ctx context.Context, playlistID uint64) ([]models.ListCollaborator, error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("playlistID", playlistID).
		Msg("Getting playlist collaborators")

	// Get the playlist
	playlist, err := s.GetByID(ctx, playlistID)
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist collaborators: %w", err)
	}

	// Verify user has permission to view this playlist
	userID := ctx.Value("userID").(uint64)
	if !s.hasPlaylistReadPermission(ctx, userID, playlist) {
		log.Warn().
			Uint64("playlistID", playlistID).
			Uint64("ownerID", playlist.Data.OwnerID).
			Uint64("requestingUserID", userID).
			Msg("User attempting to view collaborators without permission")
		return nil, errors.New("you don't have permission to view this playlist's collaborators")
	}

	// Return the shared with array (may be nil)
	if playlist.Data.SharedWith == nil {
		return []models.ListCollaborator{}, nil
	}

	// return playlist.Data.SharedWith, nil
	// TODO: implement all of the collaborator stuff
	return nil, nil
}

// RemovePlaylistCollaborator removes a user from the playlist's collaborators
func (s *userPlaylistService) RemovePlaylistCollaborator(ctx context.Context, playlistID uint64, collaboratorID uint64) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("playlistID", playlistID).
		Uint64("collaboratorID", collaboratorID).
		Msg("Removing playlist collaborator")

	// Get the playlist
	playlist, err := s.GetByID(ctx, playlistID)
	if err != nil {
		return fmt.Errorf("failed to remove playlist collaborator: %w", err)
	}

	// Verify user has permission to modify sharing (only owner can do this)
	userID := ctx.Value("userID").(uint64)
	if playlist.Data.OwnerID != userID {
		log.Warn().
			Uint64("playlistID", playlistID).
			Uint64("ownerID", playlist.Data.OwnerID).
			Uint64("requestingUserID", userID).
			Msg("User attempting to modify playlist sharing without permission")
		return errors.New("only the owner can modify playlist sharing")
	}

	// Check if the playlist has any collaborators
	if playlist.Data.SharedWith == nil || len(playlist.Data.SharedWith) == 0 {
		log.Info().
			Uint64("playlistID", playlistID).
			Msg("Playlist has no collaborators to remove")
		return nil
	}

	// Find and remove the collaborator
	// var newCollaborators []models.ListCollaborator
	found := false
	// for _, collab := range playlist.Data.SharedWith {
	// if collab.UserID != collaboratorID {
	// 	newCollaborators = append(newCollaborators, collab)
	// } else {
	// 	found = true
	// }
	// }

	if !found {
		log.Info().
			Uint64("playlistID", playlistID).
			Uint64("collaboratorID", collaboratorID).
			Msg("Collaborator not found in playlist")
		return nil
	}

	// Update the playlist with the new collaborators list
	// playlist.Data.SharedWith = newCollaborators
	_, err = s.Update(ctx, playlist)
	if err != nil {
		log.Error().Err(err).
			Uint64("playlistID", playlistID).
			Uint64("collaboratorID", collaboratorID).
			Msg("Failed to update playlist after removing collaborator")
		return fmt.Errorf("failed to update playlist after removing collaborator: %w", err)
	}

	log.Info().
		Uint64("playlistID", playlistID).
		Uint64("collaboratorID", collaboratorID).
		Msg("Playlist collaborator removed successfully")

	return nil
}

// Sync with media clients

// SyncPlaylistToClients synchronizes a playlist to specified media clients
func (s *userPlaylistService) SyncPlaylistToClients(ctx context.Context, playlistID uint64, clientIDs []uint64) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("playlistID", playlistID).
		Interface("clientIDs", clientIDs).
		Msg("Syncing playlist to clients")

	// This uses the same approach as the base service's SyncPlaylist method
	// But could be extended with user-specific validation and tracking
	// return s.SyncPlaylist(ctx, playlistID, clientIDs)
	return nil
}

// GetPlaylistSyncStatus retrieves the sync status of a playlist across clients
func (s *userPlaylistService) GetPlaylistSyncStatus(ctx context.Context, playlistID uint64) (*models.PlaylistSyncStatus, error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("playlistID", playlistID).
		Msg("Getting playlist sync status")

	// Get the playlist
	playlist, err := s.GetByID(ctx, playlistID)
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist sync status: %w", err)
	}

	// Verify user has permission to view this playlist
	userID := ctx.Value("userID").(uint64)
	if !s.hasPlaylistReadPermission(ctx, userID, playlist) {
		log.Warn().
			Uint64("playlistID", playlistID).
			Uint64("ownerID", playlist.Data.OwnerID).
			Uint64("requestingUserID", userID).
			Msg("User attempting to view playlist sync status without permission")
		return nil, errors.New("you don't have permission to view this playlist's sync status")
	}

	// Create a sync status object
	// status := &models.PlaylistSyncStatus{
	// 	PlaylistID:   playlistID,
	// 	LastSynced:   playlist.Data.LastSynced,
	// 	ClientStates: make(map[uint64]models.ClientSyncState),
	// }
	//
	// // Add state information for each client
	// for _, state := range playlist.Data.SyncClientStates {
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
	// 	} else if state.LastSynced.Before(playlist.Data.LastModified) {
	// 		clientState.SyncStatus = "out_of_sync"
	// 	} else {
	// 		clientState.SyncStatus = "in_sync"
	// 	}
	//
	// 	status.ClientStates[state.ClientID] = clientState
	// }

	// return status, nil
	// TODO: Implement and test playlist sync status
	return nil, nil

}

func (s *playlistService) Delete(ctx context.Context, id uint64) error {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("id", id).
		Msg("Deleting playlist")

	// First verify this is a playlist and the user has permission
	playlist, err := s.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete playlist: %w", err)
	}

	// TODO: Verify the current user has permission to delete this playlist
	// This would check if the current user ID matches playlist.Data.ItemList.OwnerID

	// Use the user service
	err = s.repo.Delete(ctx, id)
	if err != nil {
		log.Error().Err(err).
			Uint64("id", id).
			Msg("Failed to delete playlist")
		return fmt.Errorf("failed to delete playlist: %w", err)
	}

	log.Info().
		Uint64("id", id).
		Str("title", playlist.Title).
		Msg("Playlist deleted successfully")

	return nil
}

// Helper functions

// hasPlaylistReadPermission checks if the user has read permission for a playlist
func (s *userPlaylistService) hasPlaylistReadPermission(ctx context.Context, userID uint64, playlist *models.MediaItem[*mediatypes.Playlist]) bool {
	// The owner always has read permission
	if playlist.Data.OwnerID == userID {
		return true
	}

	// Check if the playlist is shared with this user
	for _, collab := range playlist.Data.SharedWith {
		if collab.UserID == userID {
			// Any permission level (read or write) allows reading
			return true
		}
	}

	// No permission found
	return false
}

// hasPlaylistWritePermission checks if the user has write permission for a playlist
func (s *userPlaylistService) hasWritePermission(ctx context.Context, userID uint64, playlist *models.MediaItem[*mediatypes.Playlist]) bool {
	// The owner always has write permission
	if playlist.Data.OwnerID == userID {
		return true
	}

	// Check if the playlist is shared with this user with write permission
	for _, collab := range playlist.Data.SharedWith {
		if collab.UserID == userID && collab.PermissionLevel == "write" {
			return true
		}
	}

	// No write permission found
	return false
}

func (s *userPlaylistService) GetCollaborators(userID string, playlistID string) ([]UserCollaboration, error) {
	var collabs []UserCollaboration
	// note yet implemented

	return collabs, nil

}
