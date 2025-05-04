// sync_playlists_adapter.go
package sync

import (
	"context"
	"fmt"
	"time"
	
	"suasor/clients"
	"suasor/clients/media"
	"suasor/clients/media/providers"
	mediatypes "suasor/clients/media/types"
	clienttypes "suasor/clients/types"
	"suasor/repository"
	repobundles "suasor/repository/bundles"
	"suasor/services/scheduler"
	"suasor/types/models"
	"suasor/utils/logger"
)

// AdaptedPlaylistSyncJob is an improved version of the playlist sync job that uses the adapter pattern
type AdaptedPlaylistSyncJob struct {
	jobRepo            repository.JobRepository
	userRepo           repository.UserRepository
	configRepo         repository.UserConfigRepository
	clientRepos        repobundles.ClientRepositories
	clientFactory      *clients.ClientProviderFactoryService
	mediaItemRepo      repository.ClientMediaItemRepository[mediatypes.MediaData]
	userMediaDataRepo  repository.UserMediaItemDataRepository[mediatypes.MediaData]
	userMovieDataRepo  repository.UserMediaItemDataRepository[*mediatypes.Movie]
	userSeriesDataRepo repository.UserMediaItemDataRepository[*mediatypes.Series]
	userMusicDataRepo  repository.UserMediaItemDataRepository[*mediatypes.Track]
}

// NewAdaptedPlaylistSyncJob creates a new playlist sync job using the adapter pattern
func NewAdaptedPlaylistSyncJob(
	jobRepo repository.JobRepository,
	userRepo repository.UserRepository,
	configRepo repository.UserConfigRepository,
	clientRepos repobundles.ClientRepositories,
	clientFactory *clients.ClientProviderFactoryService,
	userMovieDataRepo repository.UserMediaItemDataRepository[*mediatypes.Movie],
	userSeriesDataRepo repository.UserMediaItemDataRepository[*mediatypes.Series],
	userMusicDataRepo repository.UserMediaItemDataRepository[*mediatypes.Track],
) *AdaptedPlaylistSyncJob {
	return &AdaptedPlaylistSyncJob{
		jobRepo:            jobRepo,
		userRepo:           userRepo,
		configRepo:         configRepo,
		clientRepos:        clientRepos,
		clientFactory:      clientFactory,
		userMovieDataRepo:  userMovieDataRepo,
		userSeriesDataRepo: userSeriesDataRepo,
		userMusicDataRepo:  userMusicDataRepo,
	}
}

// Name returns the unique name of the job
func (j *AdaptedPlaylistSyncJob) Name() string {
	return "system.playlist.sync.adapted"
}

// Schedule returns when the job should next run
func (j *AdaptedPlaylistSyncJob) Schedule() time.Duration {
	// Default to daily
	return 24 * time.Hour
}

// Execute runs the playlist sync job for all users
func (j *AdaptedPlaylistSyncJob) Execute(ctx context.Context) error {
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Starting adapted playlist sync job")
	
	// Get all users
	users, err := j.userRepo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("error getting users: %w", err)
	}
	
	// Process each user
	for _, user := range users {
		if err := j.processUserPlaylists(ctx, user); err != nil {
			log.Error().
				Err(err).
				Str("username", user.Username).
				Msg("Error processing playlists for user")
			// Continue with other users even if one fails
			continue
		}
	}
	
	log.Info().Msg("Adapted playlist sync job completed")
	return nil
}

// processUserPlaylists syncs playlists for a single user using adapters
func (j *AdaptedPlaylistSyncJob) processUserPlaylists(ctx context.Context, user models.User) error {
	log := logger.LoggerFromContext(ctx)
	
	// Skip inactive users
	if !user.Active {
		log.Info().
			Str("username", user.Username).
			Msg("Skipping inactive user")
		return nil
	}
	
	// Get user configuration
	config, err := j.configRepo.GetUserConfig(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("error getting user config: %w", err)
	}
	
	// Check if playlist sync is enabled for the user
	if !config.PlaylistSyncEnabled {
		log.Info().
			Str("username", user.Username).
			Msg("Playlist sync not enabled for user")
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
		log.Error().
			Err(err).
			Msg("Error creating job run record")
		return err
	}
	
	// Get all media clients for this user
	clients, err := j.getUserPlaylistClients(ctx, user.ID)
	if err != nil {
		j.completeJobRun(ctx, jobRun.ID, models.JobStatusFailed, fmt.Sprintf("Error getting media clients: %v", err))
		return err
	}
	
	// If the user has fewer than 2 clients, there's nothing to sync
	if len(clients) < 2 {
		log.Info().
			Str("username", user.Username).
			Msg("User has fewer than 2 clients, skipping playlist sync")
		j.completeJobRun(ctx, jobRun.ID, models.JobStatusCompleted, "Not enough clients to sync")
		return nil
	}
	
	// Determine the primary client (source of truth)
	var primaryClient *PlaylistClientInfo
	for i, client := range clients {
		if client.IsPrimary {
			primaryClient = clients[i]
			break
		}
	}
	
	// If no primary client is designated, use the first one
	if primaryClient == nil && len(clients) > 0 {
		primaryClient = clients[0]
	}
	
	// Perform the sync based on the specified direction
	result, err := j.performAdaptedPlaylistSync(ctx, user.ID, clients, primaryClient, config.PlaylistSyncDirection)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Error syncing playlists")
		j.completeJobRun(ctx, jobRun.ID, models.JobStatusFailed, err.Error())
		return err
	}
	
	// Update job run with results
	resultMsg := fmt.Sprintf("Synced %d playlists, created %d, updated %d",
		result.TotalSynced, result.Created, result.Updated)
	j.completeJobRun(ctx, jobRun.ID, models.JobStatusCompleted, resultMsg)
	
	return nil
}

// PlaylistSyncResult holds statistics about a playlist sync operation
type PlaylistSyncResult struct {
	TotalSynced int
	Created     int
	Updated     int
	Errors      int
}

// getUserPlaylistClients returns all media clients for a user that support playlists
func (j *AdaptedPlaylistSyncJob) getUserPlaylistClients(ctx context.Context, userID uint64) ([]*PlaylistClientInfo, error) {
	log := logger.LoggerFromContext(ctx)
	
	// Logic to get clients from repositories
	// This would be implemented according to your application's architecture
	
	// For demonstration purposes, we'll return a placeholder result
	log.Info().
		Uint64("userID", userID).
		Msg("Getting playlist clients for user")
	
	// In a real implementation, you would:
	// 1. Query your client repositories for all clients belonging to this user
	// 2. Check which ones support playlists
	// 3. Return information about them
	
	// Example placeholder implementation
	return []*PlaylistClientInfo{
		{
			ClientID:   1,
			ClientType: clienttypes.ClientMediaTypeEmby,
			Name:       "Home Emby Server",
			IsPrimary:  true,
		},
		{
			ClientID:   2,
			ClientType: clienttypes.ClientMediaTypePlex,
			Name:       "Home Plex Server",
			IsPrimary:  false,
		},
	}, nil
}

// performAdaptedPlaylistSync syncs playlists between clients using the adapter pattern
func (j *AdaptedPlaylistSyncJob) performAdaptedPlaylistSync(
	ctx context.Context,
	userID uint64,
	clients []*PlaylistClientInfo,
	primaryClient *PlaylistClientInfo,
	syncDirection string,
) (*PlaylistSyncResult, error) {
	log := logger.LoggerFromContext(ctx)
	result := &PlaylistSyncResult{}
	
	log.Info().
		Uint64("userID", userID).
		Int("clientCount", len(clients)).
		Str("syncDirection", syncDirection).
		Msg("Syncing playlists using adapter pattern")
	
	// Get the client media objects for each client info
	playlistProviders := make(map[uint64]providers.PlaylistProvider)
	for _, clientInfo := range clients {
		client, err := j.getClientMedia(ctx, userID, clientInfo.ClientID)
		if err != nil {
			log.Error().
				Err(err).
				Uint64("clientID", clientInfo.ClientID).
				Msg("Error getting client media")
			continue
		}
		
		// Check if this client supports playlists
		playlistProvider, ok := client.(providers.PlaylistProvider)
		if !ok || !playlistProvider.SupportsPlaylists() {
			log.Warn().
				Uint64("clientID", clientInfo.ClientID).
				Msg("Client doesn't support playlists")
			continue
		}
		
		playlistProviders[clientInfo.ClientID] = playlistProvider
	}
	
	// Ensure we have at least two clients that support playlists
	if len(playlistProviders) < 2 {
		return result, fmt.Errorf("not enough clients support playlists (found %d)", len(playlistProviders))
	}
	
	// Create adapters for all playlist providers
	playlistAdapters := make(map[uint64]providers.ListProvider[*mediatypes.Playlist])
	for clientID, provider := range playlistProviders {
		playlistAdapters[clientID] = providers.NewPlaylistListAdapter(provider)
	}
	
	// Handle different sync directions
	switch syncDirection {
	case "primary-to-clients":
		// Sync from primary client to all others
		if primaryClient == nil {
			return result, fmt.Errorf("no primary client found for primary-to-clients sync")
		}
		
		primaryAdapter, ok := playlistAdapters[primaryClient.ClientID]
		if !ok {
			return result, fmt.Errorf("primary client doesn't support playlists")
		}
		
		for clientID, adapter := range playlistAdapters {
			if clientID == primaryClient.ClientID {
				continue // Skip primary client
			}
			
			// Create sync adapter and perform sync
			syncAdapter := providers.NewListSyncAdapter(primaryAdapter, adapter)
			err := syncAdapter.SyncLists(ctx, nil)
			if err != nil {
				log.Error().
					Err(err).
					Uint64("sourceClientID", primaryClient.ClientID).
					Uint64("targetClientID", clientID).
					Msg("Error syncing playlists from primary to client")
				result.Errors++
				continue
			}
			
			// Update stats (approximate since we don't track individual changes in this version)
			result.TotalSynced++
		}
		
	case "clients-to-primary":
		// Sync from all clients to the primary
		if primaryClient == nil {
			return result, fmt.Errorf("no primary client found for clients-to-primary sync")
		}
		
		primaryAdapter, ok := playlistAdapters[primaryClient.ClientID]
		if !ok {
			return result, fmt.Errorf("primary client doesn't support playlists")
		}
		
		for clientID, adapter := range playlistAdapters {
			if clientID == primaryClient.ClientID {
				continue // Skip primary client
			}
			
			// Create sync adapter and perform sync
			syncAdapter := providers.NewListSyncAdapter(adapter, primaryAdapter)
			err := syncAdapter.SyncLists(ctx, nil)
			if err != nil {
				log.Error().
					Err(err).
					Uint64("sourceClientID", clientID).
					Uint64("targetClientID", primaryClient.ClientID).
					Msg("Error syncing playlists from client to primary")
				result.Errors++
				continue
			}
			
			// Update stats (approximate since we don't track individual changes in this version)
			result.TotalSynced++
		}
		
	case "bidirectional":
		// Sync between all clients in both directions
		// For bidirectional sync, we need to consider timestamps to resolve conflicts
		
		// For this simplified implementation, we'll just sync each pair of clients
		clientIDs := make([]uint64, 0, len(playlistAdapters))
		for clientID := range playlistAdapters {
			clientIDs = append(clientIDs, clientID)
		}
		
		// Sync each pair of clients
		for i := 0; i < len(clientIDs); i++ {
			for j := i + 1; j < len(clientIDs); j++ {
				sourceID := clientIDs[i]
				targetID := clientIDs[j]
				
				// Create sync adapter and perform sync (source to target)
				syncAdapter1 := providers.NewListSyncAdapter(playlistAdapters[sourceID], playlistAdapters[targetID])
				err := syncAdapter1.SyncLists(ctx, nil)
				if err != nil {
					log.Error().
						Err(err).
						Uint64("sourceClientID", sourceID).
						Uint64("targetClientID", targetID).
						Msg("Error syncing playlists from source to target")
					result.Errors++
					continue
				}
				
				// Create sync adapter and perform sync (target to source)
				syncAdapter2 := providers.NewListSyncAdapter(playlistAdapters[targetID], playlistAdapters[sourceID])
				err = syncAdapter2.SyncLists(ctx, nil)
				if err != nil {
					log.Error().
						Err(err).
						Uint64("sourceClientID", targetID).
						Uint64("targetClientID", sourceID).
						Msg("Error syncing playlists from target to source")
					result.Errors++
					continue
				}
				
				// Update stats (approximate since we don't track individual changes in this version)
				result.TotalSynced += 2
			}
		}
		
	default:
		return result, fmt.Errorf("unknown sync direction: %s", syncDirection)
	}
	
	return result, nil
}

// completeJobRun updates a job run with its final status and message
func (j *AdaptedPlaylistSyncJob) completeJobRun(ctx context.Context, jobRunID uint64, status models.JobStatus, message string) {
	if err := j.jobRepo.CompleteJobRun(ctx, jobRunID, status, message); err != nil {
		logger.LoggerFromContext(ctx).Error().
			Err(err).
			Uint64("jobRunID", jobRunID).
			Msg("Error completing job run")
	}
}

// getClientMedia gets a media client by ID
func (j *AdaptedPlaylistSyncJob) getClientMedia(ctx context.Context, userID uint64, clientID uint64) (media.ClientMedia, error) {
	// In a real implementation, this would get the client from your client factory
	// This is a simplified placeholder that would need to be replaced with your actual implementation
	
	// Example implementation (commented out):
	// clientType := clienttypes.ClientMediaTypePlex // Just an example
	// repo := j.clientRepos[clientType].(repository.ClientRepository[clienttypes.ClientMediaConfig])
	//
	// clientConfig, err := repo.GetByID(ctx, clientID)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// client, err := j.clientFactory.GetClient(ctx, clientID, clientConfig.Config.Data)
	// if err != nil {
	// 	return nil, err
	// }
	// return client.(media.ClientMedia), nil
	
	// For demonstration, let's return an error to indicate this needs implementation
	return nil, fmt.Errorf("getClientMedia not fully implemented - needs to be connected to your client repository and factory")
}

// SetupPlaylistSyncSchedule sets up the job schedule for playlist syncing
func (j *AdaptedPlaylistSyncJob) SetupPlaylistSyncSchedule(ctx context.Context, userID uint64, frequency string) error {
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
func (j *AdaptedPlaylistSyncJob) RunManualSync(ctx context.Context, userID uint64) error {
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
func (j *AdaptedPlaylistSyncJob) SyncSinglePlaylist(ctx context.Context, userID uint64, sourceClientID uint64, playlistID string) error {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", userID).
		Uint64("sourceClientID", sourceClientID).
		Str("playlistID", playlistID).
		Msg("Syncing single playlist across clients using adapters")
	
	// Get all media clients for this user
	clients, err := j.getUserPlaylistClients(ctx, userID)
	if err != nil {
		return fmt.Errorf("error getting media clients: %w", err)
	}
	
	// Find the source client info
	var sourceClientInfo *PlaylistClientInfo
	for i, client := range clients {
		if client.ClientID == sourceClientID {
			sourceClientInfo = clients[i]
			break
		}
	}
	
	if sourceClientInfo == nil {
		return fmt.Errorf("source client %d not found for user %d", sourceClientID, userID)
	}
	
	// Get source client
	sourceClient, err := j.getClientMedia(ctx, userID, sourceClientID)
	if err != nil {
		return fmt.Errorf("error getting source client: %w", err)
	}
	
	// Check if it supports playlists
	sourcePlaylistProvider, ok := sourceClient.(providers.PlaylistProvider)
	if !ok || !sourcePlaylistProvider.SupportsPlaylists() {
		return fmt.Errorf("source client %d does not support playlists", sourceClientID)
	}
	
	// Create source adapter
	sourceAdapter := providers.NewPlaylistListAdapter(sourcePlaylistProvider)
	
	// For each target client
	for _, clientInfo := range clients {
		if clientInfo.ClientID == sourceClientID {
			continue // Skip source client
		}
		
		// Get target client
		targetClient, err := j.getClientMedia(ctx, userID, clientInfo.ClientID)
		if err != nil {
			log.Error().
				Err(err).
				Uint64("targetClientID", clientInfo.ClientID).
				Msg("Error getting target client, skipping")
			continue
		}
		
		// Check if it supports playlists
		targetPlaylistProvider, ok := targetClient.(providers.PlaylistProvider)
		if !ok || !targetPlaylistProvider.SupportsPlaylists() {
			log.Warn().
				Uint64("targetClientID", clientInfo.ClientID).
				Msg("Target client doesn't support playlists, skipping")
			continue
		}
		
		// Create target adapter
		targetAdapter := providers.NewPlaylistListAdapter(targetPlaylistProvider)
		
		// Get the source playlist
		sourcePlaylists, err := sourceAdapter.SearchLists(ctx, &mediatypes.QueryOptions{
			ExternalSourceID: playlistID,
		})
		if err != nil || len(sourcePlaylists) == 0 {
			log.Error().
				Err(err).
				Str("playlistID", playlistID).
				Msg("Failed to find source playlist")
			continue
		}
		sourcePlaylist := sourcePlaylists[0]
		
		// Find matching playlist in target by title
		targetPlaylists, err := targetAdapter.SearchLists(ctx, &mediatypes.QueryOptions{})
		if err != nil {
			log.Error().
				Err(err).
				Uint64("targetClientID", clientInfo.ClientID).
				Msg("Failed to search playlists in target client")
			continue
		}
		
		var targetPlaylist *models.MediaItem[*mediatypes.Playlist]
		var targetPlaylistID string
		
		// Look for a matching playlist by title
		for _, playlist := range targetPlaylists {
			if playlist.Title == sourcePlaylist.Title {
				targetPlaylist = playlist
				
				// Get the client-specific ID
				for _, id := range targetPlaylist.SyncClients {
					if id.ID == clientInfo.ClientID {
						targetPlaylistID = id.ItemID
						break
					}
				}
				break
			}
		}
		
		// If no matching playlist found, create one
		if targetPlaylist == nil {
			newPlaylist, err := targetAdapter.CreateList(ctx, sourcePlaylist.Title, 
				sourcePlaylist.Data.ItemList.Details.Description)
			if err != nil {
				log.Error().
					Err(err).
					Uint64("targetClientID", clientInfo.ClientID).
					Msg("Failed to create playlist in target client")
				continue
			}
			targetPlaylist = newPlaylist
			
			// Get the client-specific ID
			for _, id := range targetPlaylist.SyncClients {
				if id.ID == clientInfo.ClientID {
					targetPlaylistID = id.ItemID
					break
				}
			}
		}
		
		// Create a sync adapter for just these two playlists
		syncAdapter := providers.NewListSyncAdapter(sourceAdapter, targetAdapter)
		
		// Sync items from the source playlist to the target playlist
		// Get the source playlist items
		sourceItems, err := sourceAdapter.GetListItems(ctx, playlistID, nil)
		if err != nil {
			log.Error().
				Err(err).
				Str("playlistID", playlistID).
				Msg("Failed to get source playlist items")
			continue
		}
		
		// Get the target playlist items
		targetItems, err := targetAdapter.GetListItems(ctx, targetPlaylistID, nil)
		if err != nil {
			// If error, assume empty playlist
			targetItems = []*models.MediaItem[*mediatypes.Playlist]{}
		}
		
		// Create a map of target items by title for quick lookup
		targetItemsByTitle := make(map[string]bool)
		for _, item := range targetItems {
			targetItemsByTitle[item.Title] = true
		}
		
		// Add each source item to target if not already present
		for _, sourceItem := range sourceItems {
			if targetItemsByTitle[sourceItem.Title] {
				// Skip items that already exist
				continue
			}
			
			// Get the item ID
			var sourceItemID string
			for _, id := range sourceItem.SyncClients {
				if id.ID == sourceClientID {
					sourceItemID = id.ItemID
					break
				}
			}
			
			// Add to target
			err = targetAdapter.AddItemList(ctx, targetPlaylistID, sourceItemID)
			if err != nil {
				log.Error().
					Err(err).
					Str("targetPlaylistID", targetPlaylistID).
					Str("sourceItemID", sourceItemID).
					Msg("Failed to add item to target playlist")
				continue
			}
		}
		
		log.Info().
			Uint64("targetClientID", clientInfo.ClientID).
			Str("playlistTitle", sourcePlaylist.Title).
			Msg("Successfully synced playlist to target client")
	}
	
	return nil
}