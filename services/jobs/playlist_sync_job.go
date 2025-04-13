package jobs

import (
	"context"
	"fmt"
	"log"
	"time"

	"suasor/client"
	"suasor/client/media"
	"suasor/client/media/providers"
	mediatypes "suasor/client/media/types"
	clienttypes "suasor/client/types"

	"suasor/repository"
	"suasor/services/scheduler"
	"suasor/types/models"
)

// PlaylistSyncJob synchronizes playlists between different media clients
type PlaylistSyncJob struct {
	jobRepo          repository.JobRepository
	userRepo         repository.UserRepository
	configRepo       repository.UserConfigRepository
	clientRepos      map[clienttypes.MediaClientType]interface{}
	clientFactory    *client.ClientFactoryService
	mediaItemRepo    repository.MediaItemRepository[mediatypes.MediaData]
	mediaHistoryRepo repository.MediaPlayHistoryRepository
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

// findClientItemID retrieves the client-specific item ID from a media item's ClientIDs array
func findClientItemID[T mediatypes.MediaData](item *models.MediaItem[T], clientID uint64) (string, bool) {
	for _, cid := range item.ClientIDs {
		if cid.ID == clientID {
			return cid.ItemID, true
		}
	}
	return "", false
}

// findMatchingTargetItem looks up a media item in the target client given its source client ID
func (j *PlaylistSyncJob) findMatchingTargetItem(
	ctx context.Context,
	sourceClientID uint64,
	sourceItemID string,
	targetClientID uint64,
) (string, error) {
	if j.mediaItemRepo == nil {
		return "", fmt.Errorf("mediaItemRepo not initialized")
	}

	// First, find the media item using the source client's information
	sourceItem, err := j.mediaItemRepo.GetByClientItemID(ctx, sourceItemID, sourceClientID)
	if err != nil {
		return "", fmt.Errorf("source item not found: %w", err)
	}

	// If we have the media item, check if it has an ID for the target client
	for _, clientID := range sourceItem.ClientIDs {
		if clientID.ID == targetClientID {
			// Found the matching ID in the target client
			return clientID.ItemID, nil
		}
	}

	return "", fmt.Errorf("no matching item ID in target client")
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
	logger := log.Logger{} // Ideally use structured logging from context
	logger.Printf("Syncing playlists for user %d across %d clients", userID, len(clients))

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

	// Get all user media clients that support playlists
	playlistClients := make(map[uint64]media.MediaClient)
	for _, clientInfo := range clients {
		client, err := j.getMediaClient(ctx, userID, clientInfo.ClientID)
		if err != nil {
			logger.Printf("Error getting client %d: %v", clientInfo.ClientID, err)
			continue
		}

		playlistProvider, ok := client.(providers.PlaylistProvider)
		if !ok || !playlistProvider.SupportsPlaylists() {
			logger.Printf("Client %d does not support playlists", clientInfo.ClientID)
			continue
		}

		playlistClients[clientInfo.ClientID] = client
	}

	if len(playlistClients) < 2 {
		logger.Printf("Not enough playlist clients to sync")
		return stats, nil
	}

	// 1. Fetch all playlists from all clients
	clientPlaylists := make(map[uint64][]models.MediaItem[*mediatypes.Playlist])
	for clientID, client := range playlistClients {
		provider := client.(providers.PlaylistProvider)
		playlists, err := provider.GetPlaylists(ctx, &mediatypes.QueryOptions{})
		if err != nil {
			logger.Printf("Error fetching playlists from client %d: %v", clientID, err)
			continue
		}
		clientPlaylists[clientID] = playlists
	}

	// 2. Handle the different sync directions
	// "primary-to-clients": Primary client is source of truth
	// "clients-to-primary": Changes in clients override primary
	// "bidirectional": Most recent changes win
	switch syncDirection {
	case "primary-to-clients":
		if primaryClient == nil {
			return stats, fmt.Errorf("no primary client found for primary-to-clients sync")
		}
		stats = j.syncPrimaryToClients(ctx, userID, *primaryClient, clients, clientPlaylists, playlistClients)
	case "clients-to-primary":
		if primaryClient == nil {
			return stats, fmt.Errorf("no primary client found for clients-to-primary sync")
		}
		stats = j.syncClientsToPrimary(ctx, userID, *primaryClient, clients, clientPlaylists, playlistClients)
	case "bidirectional":
		stats = j.syncBidirectional(ctx, userID, clients, clientPlaylists, playlistClients)
	default:
		return stats, fmt.Errorf("unknown sync direction: %s", syncDirection)
	}

	logger.Printf("Synced %d playlists, created %d, updated %d, conflicts %d",
		stats.totalSynced, stats.created, stats.updated, stats.conflicts)

	return stats, nil
}

// getMediaClient gets a media client by ID
func (j *PlaylistSyncJob) getMediaClient(ctx context.Context, userID, clientID uint64) (media.MediaClient, error) {
	// This is a simplified implementation
	// In a real implementation, we would use the client factory to get the client
	clientType := clienttypes.MediaClientTypePlex // Just an example
	repo := j.clientRepos[clientType].(repository.ClientRepository[clienttypes.MediaClientConfig])

	clientConfig, err := repo.GetByID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	client, err := j.clientFactory.GetClient(ctx, clientID, clientConfig.Config.Data)
	if err != nil {
		return nil, err
	}

	return client.(media.MediaClient), nil
}

// syncPrimaryToClients syncs playlists from the primary client to all other clients
func (j *PlaylistSyncJob) syncPrimaryToClients(
	ctx context.Context,
	userID uint64,
	primaryClient PlaylistClientInfo,
	clients []PlaylistClientInfo,
	clientPlaylists map[uint64][]models.MediaItem[*mediatypes.Playlist],
	playlistClients map[uint64]media.MediaClient,
) PlaylistSyncStats {
	stats := PlaylistSyncStats{}
	logger := log.Logger{} // Ideally use structured logging from context

	// Get primary client playlists
	primaryPlaylists, ok := clientPlaylists[primaryClient.ClientID]
	if !ok {
		logger.Printf("No playlists found for primary client %d", primaryClient.ClientID)
		return stats
	}

	// Get matching items map for target clients
	for _, playlist := range primaryPlaylists {
		// For each target client
		for _, clientInfo := range clients {
			if clientInfo.ClientID == primaryClient.ClientID {
				continue // Skip primary client
			}

			targetClient, ok := playlistClients[clientInfo.ClientID]
			if !ok {
				continue
			}

			targetProvider := targetClient.(providers.PlaylistProvider)

			// Find matching playlist on target or create it
			targetPlaylists, ok := clientPlaylists[clientInfo.ClientID]
			if !ok {
				continue
			}

			var targetPlaylist *models.MediaItem[*mediatypes.Playlist]
			for i, p := range targetPlaylists {
				if p.Data.Details.Title == playlist.Data.Details.Title {
					targetPlaylist = &targetPlaylists[i]
					break
				}
			}

			if targetPlaylist == nil {
				// Create new playlist on target
				newPlaylist, err := targetProvider.CreatePlaylist(ctx,
					playlist.Data.Details.Title,
					playlist.Data.Details.Description)
				if err != nil {
					logger.Printf("Error creating playlist on client %d: %v", clientInfo.ClientID, err)
					continue
				}
				targetPlaylist = newPlaylist
				stats.created++
			}

			// Sync playlist items
			_, err := j.syncPlaylistItems(ctx, userID, playlist, *targetPlaylist,
				primaryClient.ClientID, clientInfo.ClientID, targetProvider)
			if err != nil {
				logger.Printf("Error syncing playlist items: %v", err)
			} else {
				stats.updated++
				stats.totalSynced++
			}
		}
	}

	return stats
}

// syncClientsToPrimary syncs playlists from clients to the primary client
func (j *PlaylistSyncJob) syncClientsToPrimary(
	ctx context.Context,
	userID uint64,
	primaryClient PlaylistClientInfo,
	clients []PlaylistClientInfo,
	clientPlaylists map[uint64][]models.MediaItem[*mediatypes.Playlist],
	playlistClients map[uint64]media.MediaClient,
) PlaylistSyncStats {
	stats := PlaylistSyncStats{}
	logger := log.Logger{} // Ideally use structured logging from context

	primaryMediaClient, ok := playlistClients[primaryClient.ClientID]
	if !ok {
		logger.Printf("Primary client %d not found in playlist clients", primaryClient.ClientID)
		return stats
	}

	primaryProvider := primaryMediaClient.(providers.PlaylistProvider)
	primaryPlaylistMap := make(map[string]*models.MediaItem[*mediatypes.Playlist])

	// Create a map of primary playlist titles for easier lookup
	primaryPlaylists, ok := clientPlaylists[primaryClient.ClientID]
	if ok {
		for i, p := range primaryPlaylists {
			primaryPlaylistMap[p.Data.Details.Title] = &primaryPlaylists[i]
		}
	}

	// For each client (except primary)
	for _, clientInfo := range clients {
		if clientInfo.ClientID == primaryClient.ClientID {
			continue
		}

		clientPlaylists, ok := clientPlaylists[clientInfo.ClientID]
		if !ok {
			continue
		}

		// For each playlist in this client
		for _, playlist := range clientPlaylists {
			// Check if this playlist exists in primary
			primaryPlaylist, exists := primaryPlaylistMap[playlist.Data.Details.Title]

			if !exists {
				// Create new playlist on primary
				newPlaylist, err := primaryProvider.CreatePlaylist(ctx,
					playlist.Data.Details.Title,
					playlist.Data.Details.Description)
				if err != nil {
					logger.Printf("Error creating playlist on primary client: %v", err)
					continue
				}
				primaryPlaylist = newPlaylist
				stats.created++
			}

			// Sync playlist items
			_, err := j.syncPlaylistItems(ctx, userID, playlist, *primaryPlaylist,
				clientInfo.ClientID, primaryClient.ClientID, primaryProvider)
			if err != nil {
				logger.Printf("Error syncing playlist items: %v", err)
			} else {
				stats.updated++
				stats.totalSynced++
			}
		}
	}

	return stats
}

// syncBidirectional syncs playlists between all clients, with most recent changes winning
func (j *PlaylistSyncJob) syncBidirectional(
	ctx context.Context,
	userID uint64,
	clients []PlaylistClientInfo,
	clientPlaylists map[uint64][]models.MediaItem[*mediatypes.Playlist],
	playlistClients map[uint64]media.MediaClient,
) PlaylistSyncStats {
	stats := PlaylistSyncStats{}
	logger := log.Logger{} // Ideally use structured logging from context

	// Create a map of all playlists by title for conflict resolution
	type PlaylistVersion struct {
		ClientID     uint64
		Playlist     models.MediaItem[*mediatypes.Playlist]
		LastModified time.Time
	}

	playlistVersions := make(map[string][]PlaylistVersion)

	// Group playlists by title across all clients
	for clientID, playlists := range clientPlaylists {
		for _, playlist := range playlists {
			title := playlist.Data.Details.Title
			playlistVersions[title] = append(playlistVersions[title], PlaylistVersion{
				ClientID:     clientID,
				Playlist:     playlist,
				LastModified: playlist.Data.LastModified,
			})
		}
	}

	// For each playlist title, determine the most recent version
	for _, versions := range playlistVersions {
		if len(versions) <= 1 {
			continue // No need to sync if only one version exists
		}

		// Find the most recently modified version
		var newestVersion PlaylistVersion
		for i, version := range versions {
			if i == 0 || version.LastModified.After(newestVersion.LastModified) {
				newestVersion = version
			}
		}

		// Use the newest version as source of truth
		for _, version := range versions {
			if version.ClientID == newestVersion.ClientID {
				continue // Skip the source version
			}

			targetClient, ok := playlistClients[version.ClientID]
			if !ok {
				continue
			}

			targetProvider := targetClient.(providers.PlaylistProvider)

			// Sync playlist items from newest to this version
			_, err := j.syncPlaylistItems(ctx, userID, newestVersion.Playlist, version.Playlist,
				newestVersion.ClientID, version.ClientID, targetProvider)
			if err != nil {
				logger.Printf("Error syncing playlist items: %v", err)
			} else {
				stats.updated++
				stats.totalSynced++
			}
		}
	}

	return stats
}

// syncPlaylistItems syncs items from source playlist to target playlist
func (j *PlaylistSyncJob) syncPlaylistItems(
	ctx context.Context,
	userID uint64,
	sourcePlaylist models.MediaItem[*mediatypes.Playlist],
	targetPlaylist models.MediaItem[*mediatypes.Playlist],
	sourceClientID uint64,
	targetClientID uint64,
	targetProvider providers.PlaylistProvider,
) (int, error) {
	logger := log.Logger{} // Ideally use structured logging from context

	// Get the target playlist's client-specific ID by finding it in the ClientIDs array
	var targetPlaylistID string
	for _, cid := range targetPlaylist.ClientIDs {
		if cid.ID == targetClientID {
			targetPlaylistID = cid.ItemID
			break
		}
	}

	// Get the source playlist's client-specific ID
	var sourcePlaylistID string
	for _, cid := range sourcePlaylist.ClientIDs {
		if cid.ID == sourceClientID {
			sourcePlaylistID = cid.ItemID
			break
		}
	}

	logger.Printf("Syncing items from playlist %s (client %d) to playlist %s (client %d)",
		sourcePlaylistID, sourceClientID, targetPlaylistID, targetClientID)

	// Get the source items - using the most appropriate method based on what's available
	var sourceItems []string

	// First, check if we have a SyncClientState for the source client
	if sourcePlaylist.Data.SyncClientStates != nil {
		sourceState := sourcePlaylist.Data.SyncClientStates.GetSyncClientState(sourceClientID)
		if sourceState != nil && len(sourceState.ItemIDs) > 0 {
			// Use client-specific item IDs from the sync state
			sourceItems = sourceState.ItemIDs
			logger.Printf("Using %d items from sync client state for client %d", len(sourceItems), sourceClientID)
		}
	}

	// If no items found in sync state, check the Items array
	if len(sourceItems) == 0 && len(sourcePlaylist.Data.Items) > 0 {
		for _, item := range sourcePlaylist.Data.Items {
			sourceItems = append(sourceItems, item.ItemID)
		}
		logger.Printf("Using %d items from PlaylistItems for source playlist", len(sourceItems))
	}

	// If still no items, fall back to the legacy ItemIDs array
	if len(sourceItems) == 0 && len(sourcePlaylist.Data.ItemIDs) > 0 {
		for _, itemID := range sourcePlaylist.Data.ItemIDs {
			sourceItems = append(sourceItems, fmt.Sprintf("%d", itemID))
		}
		logger.Printf("Using %d items from legacy ItemIDs for source playlist", len(sourceItems))
	}

	// For each source item, find its corresponding ID in the target client
	syncCount := 0
	var targetItems []string

	for _, sourceItemID := range sourceItems {
		// Find the target client's ID for this item
		targetItemID, err := j.findMatchingTargetItem(ctx, sourceClientID, sourceItemID, targetClientID)
		if err != nil {
			logger.Printf("Could not find matching item for %s in target client: %v", sourceItemID, err)
			continue
		}

		// Add item to target playlist using the target client's playlist ID format
		if targetPlaylistID == "" {
			logger.Printf("Warning: Empty target playlist ID for client %d", targetClientID)
			continue
		}

		// Add the item to the target playlist on the client
		err = targetProvider.AddItemToPlaylist(ctx, targetPlaylistID, targetItemID)
		if err != nil {
			logger.Printf("Error adding item to target playlist: %v", err)
			continue
		}

		// Add to our list of successfully synced target items
		targetItems = append(targetItems, targetItemID)
		syncCount++

		// Record this change in the playlist items history
		now := time.Now()

		// Update the Items array if it exists
		if len(targetPlaylist.Data.Items) > 0 {
			// Find or create the item in the target playlist's Items array
			found := false
			for i, item := range targetPlaylist.Data.Items {
				if item.ItemID == targetItemID {
					// Update existing item
					targetPlaylist.Data.Items[i].LastChanged = now
					targetPlaylist.Data.Items[i].ChangeHistory = append(targetPlaylist.Data.Items[i].ChangeHistory,
						mediatypes.ChangeRecord{
							ClientID:   sourceClientID,
							ItemID:     sourceItemID,
							ChangeType: "sync",
							Timestamp:  now,
						})
					found = true
					break
				}
			}

			if !found {
				// Add new item to the playlist
				targetPlaylist.Data.Items = append(targetPlaylist.Data.Items, mediatypes.PlaylistItem{
					ItemID:      targetItemID,
					Position:    len(targetPlaylist.Data.Items),
					LastChanged: now,
					ChangeHistory: []mediatypes.ChangeRecord{
						{
							ClientID:   sourceClientID,
							ItemID:     sourceItemID,
							ChangeType: "sync",
							Timestamp:  now,
						},
					},
				})
			}
		}
	}

	// Update the SyncClientStates to store the latest client-specific IDs
	// This ensures we have a record of which items are on each client
	if targetPlaylist.Data.SyncClientStates == nil {
		targetPlaylist.Data.SyncClientStates = mediatypes.SyncClientStates{}
	}

	// Store the target client's item IDs
	targetPlaylist.Data.SyncClientStates.AddOrUpdateSyncClientState(
		targetClientID,
		targetItems,
		targetPlaylistID,
	)

	// Update last synced timestamp for both playlists
	now := time.Now()
	sourcePlaylist.Data.LastSynced = now
	targetPlaylist.Data.LastSynced = now

	return syncCount, nil
}

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
	logger := log.Logger{} // Ideally use structured logging from context
	logger.Printf("Syncing single playlist %s from client %d for user %d", playlistID, sourceClientID, userID)

	// Get user configuration to determine sync direction
	config, err := j.configRepo.GetUserConfig(ctx, userID)
	if err != nil {
		return fmt.Errorf("error getting user config: %w", err)
	}

	if !config.PlaylistSyncEnabled {
		return fmt.Errorf("playlist sync not enabled for user %d", userID)
	}

	// Get all clients for this user
	clients, err := j.getUserMediaClients(ctx, userID)
	if err != nil {
		return fmt.Errorf("error getting media clients: %w", err)
	}

	// Find the source client info
	var sourceClientInfo *PlaylistClientInfo
	for i, client := range clients {
		if client.ClientID == sourceClientID {
			sourceClientInfo = &clients[i]
			break
		}
	}

	if sourceClientInfo == nil {
		return fmt.Errorf("source client %d not found for user %d", sourceClientID, userID)
	}

	// Get source client
	sourceClient, err := j.getMediaClient(ctx, userID, sourceClientID)
	if err != nil {
		return fmt.Errorf("error getting source client: %w", err)
	}

	sourceProvider, ok := sourceClient.(providers.PlaylistProvider)
	if !ok || !sourceProvider.SupportsPlaylists() {
		return fmt.Errorf("source client %d does not support playlists", sourceClientID)
	}

	// Get source playlist
	options := &mediatypes.QueryOptions{
		ExternalSourceID: playlistID,
	}

	sourcePlaylists, err := sourceProvider.GetPlaylists(ctx, options)
	if err != nil {
		return fmt.Errorf("error getting source playlist: %w", err)
	}

	if len(sourcePlaylists) == 0 {
		return fmt.Errorf("playlist %s not found in client %d", playlistID, sourceClientID)
	}

	sourcePlaylist := sourcePlaylists[0]

	// Update the SyncClientState for this playlist if needed
	if sourcePlaylist.Data.SyncClientStates == nil {
		sourcePlaylist.Data.SyncClientStates = mediatypes.SyncClientStates{}
	}

	// Get the playlist items for this source playlist
	playlistItems, err := sourceProvider.GetPlaylistItems(ctx, playlistID, nil)
	if err != nil {
		logger.Printf("Error getting playlist items for source playlist: %v", err)
		// Continue with empty items rather than failing completely
	} else {
		// Extract item IDs for syncing
		var sourceItemIDs []string
		for _, item := range playlistItems {
			// Get the client-specific ID for this item
			clientItemID, found := findClientItemID(&item, sourceClientID)
			if found {
				sourceItemIDs = append(sourceItemIDs, clientItemID)
			}
		}

		// Store these in the source playlist's SyncClientState
		sourcePlaylist.Data.SyncClientStates.AddOrUpdateSyncClientState(
			sourceClientID,
			sourceItemIDs,
			playlistID,
		)
	}

	// For each target client
	for _, clientInfo := range clients {
		if clientInfo.ClientID == sourceClientID {
			continue // Skip source client
		}

		// Get target client
		targetClient, err := j.getMediaClient(ctx, userID, clientInfo.ClientID)
		if err != nil {
			logger.Printf("Error getting target client %d: %v", clientInfo.ClientID, err)
			continue
		}

		targetProvider, ok := targetClient.(providers.PlaylistProvider)
		if !ok || !targetProvider.SupportsPlaylists() {
			logger.Printf("Target client %d does not support playlists", clientInfo.ClientID)
			continue
		}

		// Check if playlist already exists in target
		targetPlaylists, err := targetProvider.GetPlaylists(ctx, &mediatypes.QueryOptions{})
		if err != nil {
			logger.Printf("Error getting playlists from target client %d: %v", clientInfo.ClientID, err)
			continue
		}

		var targetPlaylist *models.MediaItem[*mediatypes.Playlist]
		for i, p := range targetPlaylists {
			if p.Data.Details.Title == sourcePlaylist.Data.Details.Title {
				targetPlaylist = &targetPlaylists[i]
				break
			}
		}

		// Create or update target playlist
		if targetPlaylist == nil {
			// Create new playlist on target
			newPlaylist, err := targetProvider.CreatePlaylist(ctx,
				sourcePlaylist.Data.Details.Title,
				sourcePlaylist.Data.Details.Description)
			if err != nil {
				logger.Printf("Error creating playlist on client %d: %v", clientInfo.ClientID, err)
				continue
			}
			targetPlaylist = newPlaylist
		}

		// Get the client-specific playlist ID for target by finding it in ClientIDs
		var targetPlaylistID string
		for _, cid := range targetPlaylist.ClientIDs {
			if cid.ID == clientInfo.ClientID {
				targetPlaylistID = cid.ItemID
				break
			}
		}

		// Map source items to target items and sync
		syncCount, err := j.syncPlaylistItems(ctx, userID, sourcePlaylist, *targetPlaylist,
			sourceClientID, clientInfo.ClientID, targetProvider)
		if err != nil {
			logger.Printf("Error syncing playlist items: %v", err)
			continue
		}

		logger.Printf("Successfully synced %d items from playlist %s (client ID: %s) to client %d (playlist ID: %s) for user %d",
			syncCount, sourcePlaylist.Data.Details.Title, playlistID, clientInfo.ClientID, targetPlaylistID, userID)
	}

	return nil
}

// InitSyncServices initializes additional services needed for sync
func (j *PlaylistSyncJob) InitSyncServices(
	mediaItemRepo repository.MediaItemRepository[mediatypes.MediaData],
	mediaHistoryRepo repository.MediaPlayHistoryRepository,
) {
	// Store these repositories for use in sync operations
	j.mediaItemRepo = mediaItemRepo
	j.mediaHistoryRepo = mediaHistoryRepo
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

	// Mock implementation
	return []mediatypes.Playlist{}, nil
}

// resolvePlaylistConflicts handles conflicts between different versions of the same playlist
// It determines which version should win based on sync direction and modification timestamps
func (j *PlaylistSyncJob) resolvePlaylistConflicts(
	ctx context.Context,
	sourcePlaylist models.MediaItem[*mediatypes.Playlist],
	targetPlaylist models.MediaItem[*mediatypes.Playlist],
	syncDirection string,
) (bool, error) {
	// Returns true if source playlist should override target playlist

	// If no sync direction specified, default to bidirectional
	if syncDirection == "" {
		syncDirection = "bidirectional"
	}

	switch syncDirection {
	case "primary-to-clients":
		// If primary is source, always override
		return true, nil
	case "clients-to-primary":
		// If primary is target, always override
		return false, nil
	case "bidirectional":
		// Most recent changes win
		return sourcePlaylist.Data.LastModified.After(targetPlaylist.Data.LastModified), nil
	default:
		return false, fmt.Errorf("unknown sync direction: %s", syncDirection)
	}
}

// UpdatePlaylistSyncPreferences updates a user's playlist sync preferences
func (j *PlaylistSyncJob) UpdatePlaylistSyncPreferences(ctx context.Context, userID uint64, preferences map[string]interface{}) error {
	// Get current user config
	config, err := j.configRepo.GetUserConfig(ctx, userID)
	if err != nil {
		return fmt.Errorf("error getting user config: %w", err)
	}

	// Update playlist sync preferences
	if enabled, ok := preferences["enabled"].(bool); ok {
		config.PlaylistSyncEnabled = enabled
	}

	if direction, ok := preferences["direction"].(string); ok {
		if direction == "primary-to-clients" || direction == "clients-to-primary" || direction == "bidirectional" {
			config.PlaylistSyncDirection = direction
		} else {
			return fmt.Errorf("invalid sync direction: %s", direction)
		}
	}

	// Update frequency if provided
	if frequency, ok := preferences["frequency"].(string); ok {
		// Set up the job schedule with the new frequency
		err = j.SetupPlaylistSyncSchedule(ctx, userID, frequency)
		if err != nil {
			return fmt.Errorf("error updating sync schedule: %w", err)
		}
	}

	// Save the updated config
	err = j.configRepo.SaveUserConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("error updating user config: %w", err)
	}

	return nil
}

// // PlaylistClientInfo holds information about a media client that supports playlists
// type PlaylistClientInfo struct {
// 	ClientID   uint64
// 	ClientType clienttypes.MediaClientType
// 	Name       string
// 	IsPrimary  bool
// }
//
// // PlaylistSyncStats tracks statistics about a playlist sync operation
// type PlaylistSyncStats struct {
// 	totalSynced int
// 	created     int
// 	updated     int
// 	conflicts   int
// }

