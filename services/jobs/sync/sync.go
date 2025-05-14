package sync

import (
	"context"
	"fmt"
	"log"
	"suasor/clients"
	"suasor/clients/media"
	"suasor/clients/media/providers"
	mediatypes "suasor/clients/media/types"
	clienttypes "suasor/clients/types"
	"suasor/repository"
	repobundles "suasor/repository/bundles"
	servicebundles "suasor/services/bundles"
	"suasor/services/scheduler"
	"suasor/types/models"
	"suasor/utils/logger"
	"time"
)

// MediaSyncJob handles syncing of media items from clients
type MediaSyncJob struct {
	jobRepo             repository.JobRepository
	userRepo            repository.UserRepository
	userConfigRepo      repository.UserConfigRepository
	clientRepos         repobundles.ClientRepositories
	dataRepos           repobundles.UserMediaDataRepositories
	clientItemRepos     repobundles.ClientMediaItemRepositories
	clientMusicServices servicebundles.ClientMusicServices
	itemRepos           repobundles.UserMediaItemRepositories
	clientFactories     *clients.ClientProviderFactoryService
}

// NewMediaSyncJob creates a new media sync job
func NewMediaSyncJob(
	jobRepo repository.JobRepository,
	userRepo repository.UserRepository,
	userConfigRepo repository.UserConfigRepository,
	clientRepos repobundles.ClientRepositories,
	dataRepos repobundles.UserMediaDataRepositories,
	clientItemRepos repobundles.ClientMediaItemRepositories,
	clientMusicServices servicebundles.ClientMusicServices,
	itemRepos repobundles.UserMediaItemRepositories,
	clientFactories *clients.ClientProviderFactoryService,
) *MediaSyncJob {
	return &MediaSyncJob{
		jobRepo:             jobRepo,
		userRepo:            userRepo,
		userConfigRepo:      userConfigRepo,
		clientRepos:         clientRepos,
		dataRepos:           dataRepos,
		clientItemRepos:     clientItemRepos,
		clientMusicServices: clientMusicServices,
		itemRepos:           itemRepos,
		clientFactories:     clientFactories,
	}
}

// Schedule returns how often the job should run (default)
func (j *MediaSyncJob) Schedule() time.Duration {
	// Default to checking hourly
	return 1 * time.Hour
}

// Name returns the unique name of the job
func (j *MediaSyncJob) Name() string {
	// Make sure we always return a valid name even if struct is empty
	if j == nil || j.jobRepo == nil {
		return "system.media.sync"
	}
	return "system.media.sync"
}

// Execute runs the job
func (j *MediaSyncJob) Execute(ctx context.Context) error {
	log.Println("Starting media sync job")

	// Check if job is properly initialized
	if j == nil || j.jobRepo == nil {
		log.Println("MediaSyncJob not properly initialized, using stub implementation")
		log.Println("Media sync job completed (no-op)")
		return nil
	}

	// Check if any sync jobs are scheduled and due
	syncJobs, err := j.jobRepo.GetMediaSyncJobsByUser(ctx, 0) // Get all sync jobs
	if err != nil {
		return fmt.Errorf("failed to get media sync jobs: %w", err)
	}

	// Process each sync job
	for _, syncJob := range syncJobs {
		// Check if this job is enabled and due to run
		if !syncJob.Enabled || !j.isDue(syncJob) {
			continue
		}

		// Run the sync job
		err := j.runSyncJob(ctx, syncJob)
		if err != nil {
			log.Printf("Error running sync job %d: %v", syncJob.ID, err)
			// Continue with other jobs even if one fails
			continue
		}

		// Update last run time
		now := time.Now()
		syncJob.LastSyncTime = &now
		err = j.jobRepo.UpdateMediaSyncJob(ctx, &syncJob)
		if err != nil {
			log.Printf("Error updating sync job last run time: %v", err)
		}
	}

	log.Println("Media sync job completed")
	return nil
}

// RunManualSync runs a manual sync for a specific user and client
func (j *MediaSyncJob) RunManualSync(ctx context.Context, userID uint64, clientID uint64, syncType models.SyncType) error {
	// First, we need to determine the client type
	// We'll need to get it from the client record in the database

	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("syncType", string(syncType)).
		Msg("Running manual sync job")

	// Validate input parameters
	if userID == 0 {
		return fmt.Errorf("invalid user ID: cannot be zero")
	}

	if clientID == 0 {
		return fmt.Errorf("invalid client ID: cannot be zero")
	}

	if !syncType.IsValid() {
		return fmt.Errorf("invalid sync type value: %s", syncType)
	}

	// Special case for "full" sync, which syncs from all clients
	if syncType == models.SyncTypeFull {
		return j.RunFullSync(ctx, userID)
	}

	// Get a list of possible client types to check
	clientTypes := []clienttypes.ClientType{
		clienttypes.ClientTypeEmby,
		clienttypes.ClientTypeJellyfin,
		clienttypes.ClientTypePlex,
		clienttypes.ClientTypeSubsonic,
	}

	// Try to find which type this client is
	var clientType clienttypes.ClientType
	for _, cType := range clientTypes {
		// Try to get the config for this client type
		config, ownerID, err := j.getClientConfig(ctx, clientID)
		if err == nil && config != nil {
			// Found the client type
			log.Info().
				Uint64("clientID", clientID).
				Uint64("ownerID", ownerID).
				Str("clientType", string(cType)).
				Msg("Found client type")
			clientType = cType
			break
		}
	}

	// If we couldn't determine the client type, return an error
	if clientType == "" {
		return fmt.Errorf("couldn't determine client type for clientID=%d", clientID)
	}

	log.Info().
		Str("clientType", string(clientType)).
		Msg("Determined client type for manual sync")

	// Create a temporary sync job
	syncJob := models.MediaSyncJob{
		UserID:     userID,
		ClientID:   clientID,
		ClientType: clientType,
		SyncType:   syncType,
	}

	// Run the sync job
	return j.runSyncJob(ctx, syncJob)
}

// RunFullSync syncs all supported media types from all of the user's clients
func (j *MediaSyncJob) RunFullSync(ctx context.Context, userID uint64) error {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", userID).
		Msg("Running full sync for all user clients")

	// Create a job run record for tracking progress
	now := time.Now()
	jobRun := &models.JobRun{
		JobName:   fmt.Sprintf("%s.full", j.Name()),
		JobType:   models.JobTypeSync,
		Status:    models.JobStatusRunning,
		StartTime: &now,
		UserID:    &userID,
		Metadata:  fmt.Sprintf(`{"userID":%d,"syncType":"full"}`, userID),
	}

	err := j.jobRepo.CreateJobRun(ctx, jobRun)
	if err != nil {
		return fmt.Errorf("failed to create job run record: %w", err)
	}

	// Get all clients for this user
	clientList, err := j.clientRepos.GetAllMediaClientsForUser(ctx, userID)
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to get media clients: %v", err)
		j.completeJobRun(ctx, jobRun.ID, models.JobStatusFailed, errorMsg)
		return fmt.Errorf(errorMsg)
	}

	if clientList.GetTotal() == 0 {
		errorMsg := "No media clients found for user"
		j.completeJobRun(ctx, jobRun.ID, models.JobStatusFailed, errorMsg)
		return fmt.Errorf(errorMsg)
	}

	log.Info().
		Int("clientCount", clientList.GetTotal()).
		Msg("Found media clients for user")

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRun.ID, 10, fmt.Sprintf("Found %d clients", clientList.GetTotal()))

	// Create a sync helper for more advanced sync operations
	syncHelper := NewListSyncHelper(j)

	// Process each client for list sync
	i := 0
	for clientID, client := range clientList.Emby {
		i++
		processClientPlaylistJob[*clienttypes.EmbyConfig](ctx, j, jobRun, userID, client)
		// Update progress
		progress := 10 + int(float64(i+1)/float64(clientList.GetTotal())*40.0)
		j.jobRepo.UpdateJobProgress(ctx, jobRun.ID, progress,
			fmt.Sprintf("Processed client %d/%d", i+1, clientList.GetTotal()))

		log.Info().
			Uint64("clientID", clientID).
			Msg("Processed client")
	}

	// Special sync for items that need multiple clients
	if clientList.GetTotal() >= 2 && syncHelper != nil {
		// // Update progress
		// j.jobRepo.UpdateJobProgress(ctx, jobRun.ID, 60, "Starting cross-client syncs")
		//
		// // For each pair of clients, sync lists from one to the other
		// for i := 0; i < len(userClients); i++ {
		// 	for j := i + 1; j < len(userClients); j++ {
		// 		sourceClientInfo := userClients[i]
		// 		targetClientInfo := userClients[j]
		//
		// 		log.Info().
		// 			Uint64("sourceClientID", sourceClientInfo.GetID()).
		// 			Uint64("targetClientID", targetClientInfo.GetID()).
		// 			Msg("Syncing lists between clients")
		//
		// 		// Get client connections
		// 		sourceClient, _, err := j.getClientMedia(ctx, sourceClientInfo.GetID())
		// 		if err != nil {
		// 			log.Error().
		// 				Err(err).
		// 				Uint64("clientID", sourceClientInfo.GetID()).
		// 				Msg("Failed to get source client connection, skipping")
		// 			continue
		// 		}
		//
		// 		targetClient, _, err := j.getClientMedia(ctx, targetClientInfo.GetID())
		// 		if err != nil {
		// 			log.Error().
		// 				Err(err).
		// 				Uint64("clientID", targetClientInfo.GetID()).
		// 				Msg("Failed to get target client connection, skipping")
		// 			continue
		// 		}
		//
		// 		// Check if both clients support playlists
		// 		sourcePlaylistProvider, sourceOk := sourceClient.(providers.ListProvider[mediatypes.ListData])
		// 		targetPlaylistProvider, targetOk := targetClient.(providers.ListProvider[mediatypes.ListData])
		//
		// 		if sourceOk && targetOk {
		// 			// Use the sync helper to sync lists between clients
		// 			syncOptions := &SyncOptions{
		// 				MediaTypes: []mediatypes.MediaType{mediatypes.MediaTypePlaylist},
		// 			}
		//
		// 			err = syncHelper.SyncLists(ctx, sourceClient, targetClient, syncOptions)
		// 			if err != nil {
		// 				log.Error().
		// 					Err(err).
		// 					Uint64("sourceClientID", sourceClientInfo.GetID()).
		// 					Uint64("targetClientID", targetClientInfo.GetID()).
		// 					Msg("Error syncing lists between clients")
		// 			}
		// 		}
		// 	}
		// }
	}

	// Also run other sync types for completeness
	syncTypes := []models.SyncType{
		models.SyncTypeMovies,
		models.SyncTypeSeries,
		models.SyncTypeMusic,
		models.SyncTypeHistory,
	}

	// Sync each type for each client
	for i, syncType := range syncTypes {
		log.Info().
			Str("syncType", string(syncType)).
			Msg("Starting sync for media type")

		for _, clientInfo := range clientList.Emby {
			// Create sub-job
			syncJob := models.MediaSyncJob{
				UserID:     userID,
				ClientID:   clientInfo.GetID(),
				ClientType: clientInfo.GetClientType(),
				SyncType:   syncType,
			}

			// Run the sync
			err = j.runSyncJob(ctx, syncJob)
			if err != nil {
				log.Error().
					Err(err).
					Uint64("clientID", clientInfo.GetID()).
					Str("syncType", string(syncType)).
					Msg("Error running sync job")
			}
		}

		// Update progress
		progress := 60 + int(float64(i+1)/float64(len(syncTypes))*40.0)
		j.jobRepo.UpdateJobProgress(ctx, jobRun.ID, progress,
			fmt.Sprintf("Completed %s sync", syncType))
	}

	// Complete the job
	j.completeJobRun(ctx, jobRun.ID, models.JobStatusCompleted, "Full sync completed")
	return nil
}

// SyncUserMediaFromClient runs a sync job for a specific user and client
// This is an alias for RunManualSync for backward compatibility
func (j *MediaSyncJob) SyncUserMediaFromClient(ctx context.Context, userID uint64, clientID uint64, syncType models.SyncType) error {
	return j.RunManualSync(ctx, userID, clientID, syncType)
}

// completeJobRun marks a job run as completed with the given status and error message
func (j *MediaSyncJob) completeJobRun(ctx context.Context, jobRunID uint64, status models.JobStatus, errorMsg string) {
	if err := j.jobRepo.CompleteJobRun(ctx, jobRunID, status, errorMsg); err != nil {
		log.Printf("Error completing job run: %v", err)
	}
}

func (j *MediaSyncJob) getClientConfig(ctx context.Context, clientID uint64) (clienttypes.ClientConfig, uint64, error) {
	log := logger.LoggerFromContext(ctx)

	// Validate input parameters
	if clientID == 0 {
		return nil, 0, fmt.Errorf("invalid client ID: cannot be zero")
	}

	log.Info().
		Uint64("clientID", clientID).
		Msg("Retrieving client config from database")

	// Get client config from database
	clientList, err := j.clientRepos.GetAllMediaClients(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get media clients: %w", err)
	}

	config := clientList.GetClientConfig(clientID)
	ownerID := clientList.GetClientOwnerID(clientID)

	// Validate that config is not nil
	if config == nil {
		return nil, 0, fmt.Errorf("retrieved nil config for clientID=%d", clientID)
	}

	return config, ownerID, nil
}

// isDue checks if a sync job is due to run
func (j *MediaSyncJob) isDue(job models.MediaSyncJob) bool {
	// If no last sync time, it's always due
	if job.LastSyncTime == nil || job.LastSyncTime.IsZero() {
		return true
	}

	// Parse the frequency
	var duration time.Duration
	freq := scheduler.Frequency(job.Frequency)

	switch freq {
	case scheduler.FrequencyDaily:
		duration = 24 * time.Hour
	case scheduler.FrequencyWeekly:
		duration = 7 * 24 * time.Hour
	case scheduler.FrequencyMonthly:
		// Approximate month as 30 days
		duration = 30 * 24 * time.Hour
	case scheduler.FrequencyManual:
		// Manual jobs are never automatically due
		return false
	default:
		// Default to daily if frequency is unknown
		duration = 24 * time.Hour
	}

	// Check if enough time has passed since the last run
	return time.Since(*job.LastSyncTime) >= duration
}

// runSyncJob executes a media sync job
func (j *MediaSyncJob) runSyncJob(ctx context.Context, syncJob models.MediaSyncJob) error {
	// Create a job run record
	now := time.Now()
	jobRun := &models.JobRun{
		JobName:   fmt.Sprintf("%s.%s", j.Name(), syncJob.SyncType),
		JobType:   models.JobTypeSync,
		Status:    models.JobStatusRunning,
		StartTime: &now,
		UserID:    &syncJob.UserID,
		Metadata:  fmt.Sprintf(`{"clientID":%d,"mediaType":"%s"}`, syncJob.ClientID, syncJob.SyncType),
	}

	// Save the job run
	err := j.jobRepo.CreateJobRun(ctx, jobRun)
	if err != nil {
		return fmt.Errorf("failed to create job run: %w", err)
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRun.ID, 0, "Starting media sync")

	// Get the client from the database
	clientMedia, ownerID, err := j.getClientMedia(ctx, syncJob.ClientID)

	if err != nil {
		errorMsg := fmt.Sprintf("Failed to get media client: %v", err)
		j.completeJobRun(ctx, jobRun.ID, models.JobStatusFailed, errorMsg)
		return fmt.Errorf(errorMsg)
	}

	// Process different media types
	var syncError error

	// Normalize media type to handle both singular and plural forms
	switch syncJob.SyncType {
	case models.SyncTypeMovies:
		syncError = j.syncMovies(ctx, clientMedia, jobRun.ID, syncJob.ClientID)
	case models.SyncTypeSeries:
		syncError = j.syncSeries(ctx, clientMedia, jobRun.ID, syncJob.ClientID)
	case models.SyncTypeMusic:
		syncError = j.syncMusic(ctx, clientMedia, jobRun.ID, syncJob.ClientID)
	case models.SyncTypeHistory:
		syncError = j.syncHistory(ctx, clientMedia, ownerID, jobRun.ID, syncJob.ClientID)
	case models.SyncTypeFavorites:
		// syncError = j.syncFavorites(ctx, clientMedia, jobRun.ID, syncJob.ClientID)
	case models.SyncTypeCollections:
		// syncError = j.syncCollections(ctx, clientMedia, jobRun.ID, syncJob.ClientID)
	case models.SyncTypePlaylists:
		syncError = j.syncPlaylists(ctx, clientMedia, jobRun.ID, syncJob.ClientID)
	default:
		syncError = fmt.Errorf("unsupported media type: %s", syncJob.SyncType)
	}

	// Complete the job run
	status := models.JobStatusCompleted
	errorMessage := ""
	if syncError != nil {
		status = models.JobStatusFailed
		errorMessage = syncError.Error()
	}

	j.completeJobRun(ctx, jobRun.ID, status, errorMessage)
	return syncError
}

// syncPlaylists syncs playlists from the client to the database
func (j *MediaSyncJob) syncPlaylists(ctx context.Context, clientMedia media.ClientMedia, jobRunID uint64, clientID uint64) error {
	log := logger.LoggerFromContext(ctx)

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 10, "Fetching playlists from client")

	// Check if client supports playlists
	playlistProvider, ok := clientMedia.(providers.PlaylistProvider)
	if !ok {
		return fmt.Errorf("client doesn't support playlists")
	}

	if !playlistProvider.SupportsPlaylists() {
		return fmt.Errorf("client doesn't support playlists")
	}

	// Get all playlists from the client
	clientType := clientMedia.(clients.Client).GetClientType().AsClientMediaType()
	log.Info().
		Uint64("clientID", clientID).
		Str("clientType", string(clientType)).
		Msg("Fetching all playlists from client")

	playlists, err := playlistProvider.SearchPlaylists(ctx, &mediatypes.QueryOptions{})
	if err != nil {
		return fmt.Errorf("failed to get playlists: %w", err)
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 30, fmt.Sprintf("Processing %d playlists", len(playlists)))
	log.Info().
		Int("playlistCount", len(playlists)).
		Msg("Retrieved playlists from client")

	// Process playlists
	batchSize := 10
	totalPlaylists := len(playlists)
	processedPlaylists := 0

	// Check if we should use the ListSyncHelper for more comprehensive sync
	// or fall back to the legacy processing method
	if helper := NewListSyncHelper(j); helper != nil {
		// Use the modern sync helper approach for better sync capabilities
		log.Info().Msg("Using ListSyncHelper for playlist sync")

		// Process each playlist
		for _, playlist := range playlists {
			// Check if playlist already exists in the local database
			var clientItemID string
			for _, cid := range playlist.SyncClients {
				if cid.ID == clientID {
					clientItemID = cid.ItemID
					break
				}
			}

			if clientItemID == "" {
				log.Warn().
					Str("playlistTitle", playlist.Title).
					Msg("Could not determine client item ID for playlist, skipping")
				continue
			}

			log.Debug().
				Str("playlistTitle", playlist.Title).
				Str("clientItemID", clientItemID).
				Msg("Processing playlist")

			// Fetch playlist items from the client
			playlistItems, err := playlistProvider.GetPlaylistItems(ctx, clientItemID)
			if err != nil {
				log.Warn().
					Err(err).
					Str("playlistTitle", playlist.Title).
					Str("clientItemID", clientItemID).
					Msg("Error fetching playlist items, continuing with next playlist")
				continue
			}

			log.Debug().
				Str("playlistTitle", playlist.Title).
				Int("itemCount", playlistItems.Len()).
				Msg("Retrieved playlist items")

			// Update the playlist with its items in the local database
			// First check if it already exists
			existingPlaylist, err := j.itemRepos.PlaylistUserRepo().GetByClientItemID(ctx, clientID, clientItemID)
			if err != nil || existingPlaylist == nil {
				// Try to find by title as fallback
				// Try to find by title as fallback but handle possible empty results
				existingPlaylist2, err2 := j.itemRepos.PlaylistUserRepo().GetByTitle(ctx, clientID, playlist.Title)
				if err2 == nil && existingPlaylist2 != nil {
					existingPlaylist = existingPlaylist2
					err = nil
				}
			}

			if err == nil && existingPlaylist != nil {
				// Update existing playlist
				log.Info().
					Str("playlistTitle", playlist.Title).
					Msg("Updating existing playlist in database")

				// Merge the data but keep our existing relationships
				existingPlaylist.Merge(playlist)

				// Update the playlist items
				if err := j.processPlaylistItems(ctx, existingPlaylist, playlistItems, clientID); err != nil {
					log.Error().
						Err(err).
						Str("playlistTitle", playlist.Title).
						Msg("Error updating playlist items")
				}

				// Save the updated playlist
				_, err = j.itemRepos.PlaylistUserRepo().Update(ctx, existingPlaylist)
				if err != nil {
					log.Error().
						Err(err).
						Str("playlistTitle", playlist.Title).
						Msg("Error saving updated playlist")
				}
			} else {
				// Create new playlist in database
				log.Info().
					Str("playlistTitle", playlist.Title).
					Msg("Creating new playlist in database")

				// Set top level title and other required fields
				playlist.Title = playlist.Data.ItemList.Details.Title

				// Process the playlist items
				if err := j.processPlaylistItems(ctx, playlist, playlistItems, clientID); err != nil {
					log.Error().
						Err(err).
						Str("playlistTitle", playlist.Title).
						Msg("Error processing playlist items for new playlist")
				}

				// Create the playlist in the database
				_, err = j.itemRepos.PlaylistUserRepo().Create(ctx, playlist)
				if err != nil {
					log.Error().
						Err(err).
						Str("playlistTitle", playlist.Title).
						Msg("Error creating new playlist")
				}
			}

			processedPlaylists++
			progress := 30 + int(float64(processedPlaylists)/float64(totalPlaylists)*70.0)
			j.jobRepo.UpdateJobProgress(ctx, jobRunID, progress, fmt.Sprintf("Processed %d/%d playlists", processedPlaylists, totalPlaylists))
		}
	} else {
		// Use the original batch processing method as fallback
		log.Info().Msg("Using legacy batch processing for playlist sync")

		for i := 0; i < totalPlaylists; i += batchSize {
			end := i + batchSize
			if end > totalPlaylists {
				end = totalPlaylists
			}

			playlistBatch := playlists[i:end]
			err := j.processPlaylistBatch(ctx, playlistBatch, playlistProvider, clientID, clientType)
			if err != nil {
				log.Error().
					Err(err).
					Int("batchStart", i).
					Int("batchEnd", end).
					Msg("Failed to process playlist batch")
				continue
			}

			processedPlaylists += len(playlistBatch)
			progress := 30 + int(float64(processedPlaylists)/float64(totalPlaylists)*70.0)
			j.jobRepo.UpdateJobProgress(ctx, jobRunID, progress, fmt.Sprintf("Processed %d/%d playlists", processedPlaylists, totalPlaylists))
		}
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 100, fmt.Sprintf("Synced %d playlists", totalPlaylists))

	return nil
}

// processPlaylistBatch processes a batch of playlists and syncs their items
func (j *MediaSyncJob) processPlaylistBatch(ctx context.Context, playlists []*models.MediaItem[*mediatypes.Playlist], provider providers.PlaylistProvider, clientID uint64, clientType clienttypes.ClientMediaType) error {
	for _, playlist := range playlists {
		// Skip if playlist has no client ID information
		if len(playlist.SyncClients) == 0 {
			log.Printf("Skipping playlist with no client IDs: %s", playlist.Data.ItemList.Details.Title)
			continue
		}

		// Get the client ID and item ID for lookup
		clientItemID := ""
		for _, cid := range playlist.SyncClients {
			if cid.ID == clientID {
				clientItemID = cid.ItemID
				break
			}
		}

		if clientItemID == "" {
			log.Printf("No matching client item ID found for playlist: %s", playlist.Data.ItemList.Details.Title)
			continue
		}

		// Fetch the playlist items
		playlistItems, err := provider.GetPlaylistItems(ctx, clientItemID)
		if err != nil {
			log.Printf("Error getting items for playlist %s: %v", playlist.Data.ItemList.Details.Title, err)
			continue
		}

		// Update the playlist with its items
		if err := j.processPlaylistItems(ctx, playlist, playlistItems, clientID); err != nil {
			log.Printf("Error processing items for playlist %s: %v", playlist.Data.ItemList.Details.Title, err)
			continue
		}

		// Check if the playlist already exists in the database
		existingPlaylist, err := j.itemRepos.PlaylistUserRepo().GetByClientItemID(ctx, clientID, clientItemID)

		if err != nil || existingPlaylist == nil {
			// Try to find by title - not ideal but a fallback
			// Try to find by title as fallback but handle possible empty results
			existingPlaylist2, err2 := j.itemRepos.PlaylistUserRepo().GetByTitle(ctx, clientID, playlist.Data.ItemList.Details.Title)
			if err2 == nil && existingPlaylist2 != nil {
				existingPlaylist = existingPlaylist2
				err = nil
			}
		}

		if err == nil && existingPlaylist != nil {
			// Playlist exists, update it
			existingPlaylist.Merge(playlist)

			// Save the updated playlist
			log.Printf("Updating playlist: %s", playlist.Data.ItemList.Details.Title)
			_, err = j.itemRepos.PlaylistUserRepo().Update(ctx, existingPlaylist)
			if err != nil {
				log.Printf("Error updating playlist: %v", err)
				continue
			}
		} else {
			// Playlist doesn't exist, create it
			// Set top level title and other required fields
			playlist.Title = playlist.Data.ItemList.Details.Title

			// Create the playlist
			_, err = j.itemRepos.PlaylistUserRepo().Create(ctx, playlist)
			if err != nil {
				log.Printf("Error creating playlist: %v", err)
				continue
			}
		}
	}

	return nil
}

// processPlaylistItems processes items in a playlist and links them to media items
func (j *MediaSyncJob) processPlaylistItems(ctx context.Context, playlist *models.MediaItem[*mediatypes.Playlist], playlistItems *models.MediaItemList[*mediatypes.Playlist], clientID uint64) error {
	log := logger.LoggerFromContext(ctx)
	// Get or create the sync state for this client
	_, exists := playlist.SyncClients.GetSyncStatus(clientID)
	if !exists {
		log.Error().
			Uint64("clientID", clientID).
			Str("clientListID", playlist.SyncClients.GetClientItemID(clientID)).
			Msg("Client not found in sync state")
		return fmt.Errorf("client not found in sync state")
	}

	playlist.SyncClients.UpdateSyncStatus(clientID, models.SyncStatusPending)

	playlistItems.ForEach(func(itemID string, mediaType mediatypes.MediaType, item any) bool {
		typedItem, ok := item.(*models.MediaItem[mediatypes.MediaData])
		if !ok {
			log.Warn().
				Str("itemID", itemID).
				Msg("Could not convert item to media item, skipping")
			return true
		}
		// Find the client-specific ID for this item
		itemClientID := ""
		for _, cid := range typedItem.SyncClients {
			if cid.ID == clientID {
				itemClientID = cid.ItemID
				break
			}
		}

		if itemClientID == "" {
			log.Warn().
				Str("itemID", itemID).
				Str("itemName", typedItem.GetTitle()).
				Msg("Could not determine client item ID for item, skipping")
			return true
		}

		// Add this item to the list of client items
		return true
	})

	// Update the sync state with the new items
	playlist.Data.ItemList.LastSynced = time.Now()
	playlist.Data.ItemList.ItemCount = len(playlist.Data.ItemList.Items)

	return nil
}

// findOrCreateMediaItem finds or creates a media item in the database
func (j *MediaSyncJob) findOrCreateMediaItem(ctx context.Context, item *models.MediaItem[*mediatypes.Playlist], clientID uint64) (*models.MediaItem[*mediatypes.Playlist], error) {
	// Try to find the item by client ID first
	clientItemID := ""
	for _, cid := range item.SyncClients {
		if cid.ID == clientID {
			clientItemID = cid.ItemID
			break
		}
	}

	if clientItemID == "" {
		return nil, fmt.Errorf("no client item ID found")
	}

	// Try to find by client item ID
	existingItem, err := j.clientItemRepos.PlaylistClientRepo().GetByClientItemID(ctx, clientID, clientItemID)
	if err == nil && existingItem != nil {
		return existingItem, nil
	}

	// If not found, try to find by external IDs
	if item.Data.GetDetails().ExternalIDs != nil && len(item.Data.GetDetails().ExternalIDs) > 0 {
		existingItem, err = j.clientItemRepos.PlaylistClientRepo().GetByExternalIDs(ctx, item.Data.GetDetails().ExternalIDs)
		if err == nil && existingItem != nil {
			return existingItem, nil
		}
	}

	// If still not found, create a new one
	switch item.Type {
	case mediatypes.MediaTypePlaylist:
		// Create a new playlist item
		newItem := models.NewMediaItem[*mediatypes.Playlist](item.Data)
		newItem.Title = item.Data.GetDetails().Title
		newItem.SyncClients = item.SyncClients

		created, err := j.itemRepos.PlaylistUserRepo().Create(ctx, newItem)
		if err != nil {
			return nil, fmt.Errorf("failed to create playlist: %w", err)
		}

		return created, nil

	default:
		// Playlists should be of type MediaTypePlaylist
		return nil, fmt.Errorf("unsupported media type: %s", item.Type)
	}
}

func (j *MediaSyncJob) getClientMedia(ctx context.Context, clientID uint64) (media.ClientMedia, uint64, error) {
	log := logger.LoggerFromContext(ctx)

	// Validate input parameters
	if clientID == 0 {
		return nil, 0, fmt.Errorf("invalid client ID: cannot be zero")
	}

	log.Info().
		Uint64("clientID", clientID).
		Msg("Getting client media")

	// Use the type to get the client config by id
	clientConfig, userID, err := j.getClientConfig(ctx, clientID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get client config: %w", err)
	}

	// Validate client config is not nil before proceeding
	if clientConfig == nil {
		return nil, 0, fmt.Errorf("client config is nil for clientID=%d, userID=%d", clientID, userID)
	}

	// Cast media client from generic client
	client, err := j.clientFactories.GetClient(ctx, clientID, clientConfig)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get client: %w", err)
	}

	// Cast to media client
	clientMedia, ok := client.(media.ClientMedia)
	if !ok {
		return nil, 0, fmt.Errorf("client is not a media client")
	}

	return clientMedia, userID, nil
}

func processClientPlaylistJob[T clienttypes.ClientConfig](ctx context.Context, j *MediaSyncJob, jobRun *models.JobRun, userID uint64, clientInfo *models.Client[T]) {
	log := logger.LoggerFromContext(ctx)
	// First sync playlists
	log.Info().
		Uint64("clientID", clientInfo.GetID()).
		Str("clientName", clientInfo.GetName()).
		Msg("Syncing playlists from client")

	// Get the client connection
	_, _, err := j.getClientMedia(ctx, clientInfo.GetID())
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", clientInfo.GetID()).
			Msg("Failed to get client connection, skipping")
		return
	}

	// Create sub-job for playlist sync
	syncJob := models.MediaSyncJob{
		UserID:     userID,
		ClientID:   clientInfo.GetID(),
		ClientType: clientInfo.GetClientType(),
		SyncType:   models.SyncTypePlaylists,
	}

	// Run the playlist sync
	err = j.runSyncJob(ctx, syncJob)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", clientInfo.GetID()).
			Msg("Error syncing playlists from client")
	}

}
