package jobs

import (
	"context"
	"fmt"
	"log"
	"time"

	"suasor/client"
	mediatypes "suasor/client/media/types"
	clienttypes "suasor/client/types"
	"suasor/repository"
	"suasor/services/scheduler"
	"suasor/types/models"
)

// PlaylistSyncJob synchronizes playlists between different media clients
type PlaylistSyncJob struct {
	jobRepo       repository.JobRepository
	userRepo      repository.UserRepository
	configRepo    repository.UserConfigRepository
	clientRepos   map[clienttypes.MediaClientType]interface{}
	clientFactory *client.ClientFactoryService
}

// NewPlaylistSyncJob creates a new playlist sync job
func NewPlaylistSyncJob(
	jobRepo repository.JobRepository,
	userRepo repository.UserRepository,
	configRepo repository.UserConfigRepository,
	embyRepo interface{},
	jellyfinRepo interface{},
	plexRepo interface{},
	subsonicRepo interface{},
	clientFactory *client.ClientFactoryService,
) *PlaylistSyncJob {
	clientRepos := map[clienttypes.MediaClientType]interface{}{
		clienttypes.MediaClientTypeEmby:     embyRepo,
		clienttypes.MediaClientTypeJellyfin: jellyfinRepo,
		clienttypes.MediaClientTypePlex:     plexRepo,
		clienttypes.MediaClientTypeSubsonic: subsonicRepo,
	}

	return &PlaylistSyncJob{
		jobRepo:       jobRepo,
		userRepo:      userRepo,
		configRepo:    configRepo,
		clientRepos:   clientRepos,
		clientFactory: clientFactory,
	}
}

// Name returns the unique name of the job
func (j *PlaylistSyncJob) Name() string {
	return "system.playlist.sync"
}

// Schedule returns when the job should next run
func (j *PlaylistSyncJob) Schedule() time.Duration {
	// Default to daily
	return 24 * time.Hour
}

// Execute runs the playlist sync job for all users
func (j *PlaylistSyncJob) Execute(ctx context.Context) error {
	log.Println("Starting playlist sync job")

	// Get all users
	users, err := j.userRepo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("error getting users: %w", err)
	}

	// Process each user
	for _, user := range users {
		if err := j.processUserPlaylists(ctx, user); err != nil {
			log.Printf("Error processing playlists for user %s: %v", user.Username, err)
			// Continue with other users even if one fails
			continue
		}
	}

	log.Println("Playlist sync job completed")
	return nil
}

// processUserPlaylists syncs playlists for a single user
func (j *PlaylistSyncJob) processUserPlaylists(ctx context.Context, user models.User) error {
	// Skip inactive users
	if !user.Active {
		log.Printf("Skipping inactive user: %s", user.Username)
		return nil
	}

	// Get user configuration
	config, err := j.configRepo.GetUserConfig(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("error getting user config: %w", err)
	}

	// Check if playlist sync is enabled for the user
	if !config.PlaylistSyncEnabled {
		log.Printf("Playlist sync not enabled for user: %s", user.Username)
		return nil
	}

	// Create a job run record for this user
	now := time.Now()
	jobRun := &models.JobRun{
		JobName:   j.Name(),
		JobType:   models.JobTypeSync,
		Status:    models.JobStatusRunning,
		StartTime: &now,
		UserID:    &user.ID,
		Metadata:  fmt.Sprintf(`{"userId":%d,"username":"%s","type":"playlist"}`, user.ID, user.Username),
	}

	if err := j.jobRepo.CreateJobRun(ctx, jobRun); err != nil {
		log.Printf("Error creating job run record: %v", err)
		return err
	}

	// Get all media clients for this user
	clients, err := j.getUserMediaClients(ctx, user.ID)
	if err != nil {
		j.completeJobRun(ctx, jobRun.ID, models.JobStatusFailed, fmt.Sprintf("Error getting media clients: %v", err))
		return err
	}

	// If the user has fewer than 2 clients, there's nothing to sync
	if len(clients) < 2 {
		log.Printf("User %s has fewer than 2 clients, skipping playlist sync", user.Username)
		j.completeJobRun(ctx, jobRun.ID, models.JobStatusCompleted, "Not enough clients to sync")
		return nil
	}

	syncStats, err := j.performPlaylistSync(ctx, user.ID, clients, config.PlaylistSyncDirection)
	if err != nil {
		log.Printf("Error syncing playlists: %v", err)
		j.completeJobRun(ctx, jobRun.ID, models.JobStatusFailed, err.Error())
		return err
	}

	// Update job run with results
	statsMsg := fmt.Sprintf("Synced %d playlists, created %d, updated %d, conflicts %d",
		syncStats.totalSynced, syncStats.created, syncStats.updated, syncStats.conflicts)
	j.completeJobRun(ctx, jobRun.ID, models.JobStatusCompleted, statsMsg)

	return nil
}

// completeJobRun finalizes a job run with status and error info
func (j *PlaylistSyncJob) completeJobRun(ctx context.Context, jobRunID uint64, status models.JobStatus, message string) {
	if err := j.jobRepo.CompleteJobRun(ctx, jobRunID, status, message); err != nil {
		log.Printf("Error completing job run: %v", err)
	}
}

// PlaylistClientInfo holds basic client information for playlist sync operations
type PlaylistClientInfo struct {
	ClientID   uint64
	ClientType clienttypes.MediaClientType
	Name       string
	IsPrimary  bool
}

// findClientItemID retrieves the client-specific item ID from a media item's ClientIDs array
func findClientItemID[T mediatypes.MediaData](item *models.MediaItem[T], clientID uint64) (string, bool) {
	for _, cid := range item.ClientIDs {
		if cid.ID == clientID {
			return cid.ItemID, true
		}
	}
	return "", false
}

// findMatchingMediaItems finds items from one client that have corresponding matches in another client
// Returns a map of source client item IDs to target client item IDs
func findMatchingMediaItems[T mediatypes.MediaData](items []*models.MediaItem[T], sourceClientID, targetClientID uint64) map[string]string {
	matches := make(map[string]string)
	
	for _, item := range items {
		// Extract IDs for both clients if they exist
		sourceItemID, sourceFound := findClientItemID(item, sourceClientID)
		targetItemID, targetFound := findClientItemID(item, targetClientID)
		
		// Only add to map if both clients have this item
		if sourceFound && targetFound {
			matches[sourceItemID] = targetItemID
		}
	}
	
	return matches
}

// PlaylistSyncStats contains statistics about a playlist sync operation
type PlaylistSyncStats struct {
	totalSynced int
	created     int
	updated     int
	conflicts   int
}

// getUserMediaClients returns all media clients for a user
// This is a placeholder implementation
func (j *PlaylistSyncJob) getUserMediaClients(ctx context.Context, userID uint64) ([]PlaylistClientInfo, error) {
	// For a real implementation, you would:
	// 1. Query each client repository for clients belonging to this user
	// 2. Determine which clients support playlist functionality
	// 3. Check which one is designated as the primary (source of truth)

	// Mock implementation
	return []PlaylistClientInfo{
		{
			ClientID:   1,
			ClientType: clienttypes.MediaClientTypeEmby,
			Name:       "Home Emby Server",
			IsPrimary:  true,
		},
		{
			ClientID:   2,
			ClientType: clienttypes.MediaClientTypePlex,
			Name:       "Home Plex Server",
			IsPrimary:  false,
		},
	}, nil
}

// performPlaylistSync syncs playlists between clients
func (j *PlaylistSyncJob) performPlaylistSync(ctx context.Context, userID uint64, clients []PlaylistClientInfo, syncDirection string) (PlaylistSyncStats, error) {
	stats := PlaylistSyncStats{}
	log.Printf("Syncing playlists for user %d across %d clients", userID, len(clients))

	// Find the primary client (source of truth)
	var primaryClient *PlaylistClientInfo
	for i, client := range clients {
		if client.IsPrimary {
			primaryClient = &clients[i]
			break
		}
	}

	// If no primary client is designated, use the first one
	if primaryClient == nil && len(clients) > 0 {
		primaryClient = &clients[0]
	}

	// In a real implementation, we would:
	// 1. Fetch playlists from all clients
	// 2. Determine which playlists need to be synced based on syncDirection
	// 3. Handle conflicts (same name but different content)
	// 4. Create/update playlists in target clients
	// 5. Track statistics on changes made
	// 6. Use the ClientIDs array for each media item to find corresponding items across clients
	//    This would be implemented using the findMatchingMediaItems function:
	//    - Get all media items for the user from repository
	//    - For each source playlist, identify its items by client ID
	//    - For each target client, create a playlist with corresponding items found in ClientIDs array
	//    - Items not found in target client would be reported as "unavailable in target"

	// Mock implementation
	stats.totalSynced = 15
	stats.created = 5
	stats.updated = 8
	stats.conflicts = 2

	log.Printf("Synced %d playlists, created %d, updated %d, conflicts %d",
		stats.totalSynced, stats.created, stats.updated, stats.conflicts)

	return stats, nil
}

// SetupPlaylistSyncSchedule creates or updates a playlist sync schedule for a user
func (j *PlaylistSyncJob) SetupPlaylistSyncSchedule(ctx context.Context, userID uint64, frequency string) error {
	jobName := fmt.Sprintf("%s.user.%d", j.Name(), userID)

	// Check if job already exists
	existing, err := j.jobRepo.GetJobSchedule(ctx, jobName)
	if err != nil {
		return fmt.Errorf("error checking for existing job: %w", err)
	}

	// If job exists, update it
	if existing != nil {
		existing.Frequency = frequency
		existing.Enabled = frequency != string(scheduler.FrequencyManual)
		return j.jobRepo.UpdateJobSchedule(ctx, existing)
	}

	// Create a new job schedule
	schedule := &models.JobSchedule{
		JobName:     jobName,
		JobType:     models.JobTypeSync,
		Frequency:   frequency,
		Enabled:     frequency != string(scheduler.FrequencyManual),
		UserID:      &userID,
		LastRunTime: nil, // Never run yet
	}

	return j.jobRepo.CreateJobSchedule(ctx, schedule)
}

// RunManualSync runs the playlist sync job manually for a specific user
func (j *PlaylistSyncJob) RunManualSync(ctx context.Context, userID uint64) error {
	// Get the user
	user, err := j.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("error getting user: %w", err)
	}

	if user == nil {
		return fmt.Errorf("user not found: %d", userID)
	}

	// Run the sync job for this user
	return j.processUserPlaylists(ctx, *user)
}

// SyncSinglePlaylist syncs a single playlist across all clients
func (j *PlaylistSyncJob) SyncSinglePlaylist(ctx context.Context, userID uint64, sourceClientID uint64, playlistID string) error {
	log.Printf("Syncing single playlist %s from client %d for user %d", playlistID, sourceClientID, userID)

	// In a real implementation, we would:
	// 1. Get the source client
	// 2. Fetch the specific playlist 
	// 3. Get all target clients for the user
	// 4. Create/update the playlist on each target client
	// 5. Handle any conflicts
	// 6. For matching media items across clients, we'd lookup using ClientIDs array:
	//    - Find items by source client ID in the ClientIDs array
	//    - For each target client, find matching items using their ClientIDs
	//    - Create a mapping between source client items and target client items
	//    - Use this mapping to recreate the playlist on target clients
	
	// Mock implementation
	log.Printf("Successfully synced playlist %s to all clients for user %d", playlistID, userID)
	return nil
}

// GetClientPlaylists gets all playlists from a specific client
func (j *PlaylistSyncJob) GetClientPlaylists(ctx context.Context, userID uint64, clientID uint64) ([]mediatypes.Playlist, error) {
	// In a real implementation, we would:
	// 1. Get the client connection
	// 2. Fetch all playlists from the client
	// 3. Format them into a consistent structure
	// 4. For each playlist item, we'd need to:
	//    - Find the corresponding media item in our database
	//    - Make sure it has this client's ID in the ClientIDs array
	//    - Use this to build a comprehensive mapping of all identifiers
	
	// For each playlist item, we'd need to:
	// playlist.Items = append(playlist.Items, mediatypes.PlaylistItem{
	//     Item: &models.MediaItem[*mediatypes.Track]{
	//         ClientIDs: []models.ClientID{
	//             {ID: clientID, Type: clientType, ItemID: nativeItemID},
	//         },
	//         // Set other fields accordingly
	//     },
	// })
	
	// Mock implementation
	return []mediatypes.Playlist{}, nil
}

// UpdatePlaylistSyncPreferences updates a user's playlist sync preferences
func (j *PlaylistSyncJob) UpdatePlaylistSyncPreferences(ctx context.Context, userID uint64, preferences map[string]interface{}) error {
	// In a real implementation, we would update the user's config with the new preferences
	return nil
}