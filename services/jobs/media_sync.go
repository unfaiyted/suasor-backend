package jobs

import (
	"context"
	"fmt"
	"log"
	"suasor/client/media"
	mediatypes "suasor/client/media/types"
	"suasor/repository"
	"suasor/types/models"
	"time"
)

// We'll use MediaClientInfo from favorites_sync.go

// MediaSyncJob synchronizes media from external clients to the local database
type MediaSyncJob struct {
	jobRepo     repository.JobRepository
	userRepo    repository.UserRepository
	configRepo  repository.UserConfigRepository
	clientRepo  interface{} // Used to be repository.ClientRepository[types.MediaClientConfig]
	movieRepo   repository.MediaItemRepository[*mediatypes.Movie]
	seriesRepo  repository.MediaItemRepository[*mediatypes.Series]
	musicRepo   repository.MediaItemRepository[*mediatypes.Track]
	historyRepo repository.MediaPlayHistoryRepository
	db          interface{} // Used to create repositories
}

// NewMediaSyncJob creates a new media sync job
func NewMediaSyncJob(
	jobRepo repository.JobRepository,
	userRepo repository.UserRepository,
	configRepo repository.UserConfigRepository,
	clientRepo interface{}, // Use interface{} instead of specific repository type
	movieRepo repository.MediaItemRepository[*mediatypes.Movie],
	seriesRepo repository.MediaItemRepository[*mediatypes.Series],
	musicRepo repository.MediaItemRepository[*mediatypes.Track],
	historyRepo repository.MediaPlayHistoryRepository,
	db interface{}, // Database connection
) *MediaSyncJob {
	return &MediaSyncJob{
		jobRepo:     jobRepo,
		userRepo:    userRepo,
		configRepo:  configRepo,
		clientRepo:  clientRepo,
		movieRepo:   movieRepo,
		seriesRepo:  seriesRepo,
		musicRepo:   musicRepo,
		historyRepo: historyRepo,
		db:          db,
	}
}

// Name returns the unique name of the job
func (j *MediaSyncJob) Name() string {
	return "system.media.sync"
}

// Schedule returns when the job should next run
func (j *MediaSyncJob) Schedule() time.Duration {
	// Default to checking daily
	return 24 * time.Hour
}

// Execute runs the media sync job for all users and all clients
func (j *MediaSyncJob) Execute(ctx context.Context) error {
	log.Println("Starting media sync job")

	// Create a job run record for the overall sync
	now := time.Now()
	jobRun := &models.JobRun{
		JobName:       j.Name(),
		JobType:       models.JobTypeSync,
		Status:        models.JobStatusRunning,
		StartTime:     &now,
		Progress:      0,
		StatusMessage: "Starting media sync for all users",
	}

	err := j.jobRepo.CreateJobRun(ctx, jobRun)
	if err != nil {
		return fmt.Errorf("error creating job run record: %w", err)
	}

	// Get all users and filter active ones
	allUsers, err := j.userRepo.FindAll(ctx)
	if err != nil {
		j.jobRepo.CompleteJobRun(ctx, jobRun.ID, models.JobStatusFailed, fmt.Sprintf("Failed to get users: %v", err))
		return fmt.Errorf("failed to get users: %w", err)
	}
	
	// Filter active users
	var users []models.User
	for _, user := range allUsers {
		if user.Active {
			users = append(users, user)
		}
	}

	if len(users) == 0 {
		j.jobRepo.CompleteJobRun(ctx, jobRun.ID, models.JobStatusCompleted, "No users to process")
		return nil
	}

	// Get user configs to find their clients
	totalUsers := len(users)
	j.jobRepo.UpdateJobProgress(ctx, jobRun.ID, 5, fmt.Sprintf("Found %d users to process", totalUsers))

	// Create a context with timeout for the entire sync operation
	syncCtx, cancel := context.WithTimeout(context.Background(), 24*time.Hour)
	defer cancel()

	go func() {
		userCount := 0
		for _, user := range users {
			// Check for cancellation
			if syncCtx.Err() != nil {
				j.jobRepo.CompleteJobRun(syncCtx, jobRun.ID, models.JobStatusFailed, "Job was cancelled")
				return
			}

			// Get user config (for future use with client configs)
			_, err := j.configRepo.GetUserConfig(syncCtx, user.ID)
			if err != nil {
				log.Printf("Failed to get config for user %d: %v", user.ID, err)
				continue
			}

			// Get user's media clients
			clients, err := j.getUserMediaClients(syncCtx, user.ID)
			if err != nil {
				log.Printf("Failed to get media clients for user %d: %v", user.ID, err)
				continue
			}
			
			// Process each media client
			for _, client := range clients {
				// Process each media type the client may support
				mediaTypes := []string{"movie", "series", "music", "playlist", "collection"}
				for _, mediaType := range mediaTypes {
					// Launch sync for this client and media type
					_ = j.SyncUserMediaFromClient(syncCtx, user.ID, client.ClientID, mediaType)
					// We don't wait for each sync to complete as they run in goroutines
				}
			}

			userCount++
			progress := 5 + int(float64(userCount)/float64(totalUsers)*90.0)
			j.jobRepo.UpdateJobProgress(syncCtx, jobRun.ID, progress, fmt.Sprintf("Processed %d/%d users", userCount, totalUsers))
		}

		// Mark the overall job as complete
		j.jobRepo.UpdateJobProgress(syncCtx, jobRun.ID, 100, "Media sync completed for all users")
		j.jobRepo.CompleteJobRun(syncCtx, jobRun.ID, models.JobStatusCompleted, "")
	}()

	return nil
}

// getUserMediaClients returns all media clients for a user
func (j *MediaSyncJob) getUserMediaClients(ctx context.Context, userID uint64) ([]MediaClientInfo, error) {
	// This would query the client repository to get all media clients for the user
	// For now, we'll return a stub
	
	// Get clients from the database
	// In a real implementation, you would query the client repository
	// For now, we'll return an empty slice
	return []MediaClientInfo{}, nil
}

// SyncUserMediaFromClient synchronizes media from a client to the local database
func (j *MediaSyncJob) SyncUserMediaFromClient(ctx context.Context, userID, clientID uint64, mediaType string) error {
	log.Printf("Syncing %s media for user %d from client %d", mediaType, userID, clientID)
	
	// Create a job run record
	now := time.Now()
	jobRun := &models.JobRun{
		JobName:       j.Name(),
		JobType:       models.JobTypeSync,
		Status:        models.JobStatusRunning,
		StartTime:     &now,
		UserID:        &userID,
		Progress:      0,
		StatusMessage: fmt.Sprintf("Starting %s sync for client %d", mediaType, clientID),
		Metadata:      fmt.Sprintf(`{"userId":%d,"clientId":%d,"mediaType":"%s"}`, userID, clientID, mediaType),
	}
	
	err := j.jobRepo.CreateJobRun(ctx, jobRun)
	if err != nil {
		return fmt.Errorf("error creating job run record: %w", err)
	}
	
	// Run the sync process
	go func() {
		// Set up a new context with timeout
		syncCtx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()
		
		// Get the client
		j.jobRepo.UpdateJobProgress(syncCtx, jobRun.ID, 5, "Getting media client")
		client, err := j.getMediaClient(syncCtx, clientID)
		if err != nil {
			errorMsg := fmt.Sprintf("Failed to get media client: %v", err)
			j.jobRepo.CompleteJobRun(syncCtx, jobRun.ID, models.JobStatusFailed, errorMsg)
			log.Printf("Error: %s", errorMsg)
			return
		}
		
		// Check if the client supports this media type
		supported := j.clientSupportsMediaType(client, mediaType)
		if !supported {
			msg := fmt.Sprintf("Client does not support media type: %s", mediaType)
			j.jobRepo.CompleteJobRun(syncCtx, jobRun.ID, models.JobStatusCompleted, msg)
			log.Printf("Note: %s", msg)
			return
		}
		
		j.jobRepo.UpdateJobProgress(syncCtx, jobRun.ID, 10, fmt.Sprintf("Fetching %s items from client", mediaType))
		
		// Sync the appropriate media type
		switch mediaType {
		case "movie":
			err = j.syncMovies(syncCtx, jobRun.ID, client)
		case "series":
			err = j.syncSeries(syncCtx, jobRun.ID, client)
		case "music":
			err = j.syncMusic(syncCtx, jobRun.ID, client)
		case "playlist":
			err = j.syncPlaylists(syncCtx, jobRun.ID, client)
		case "collection":
			err = j.syncCollections(syncCtx, jobRun.ID, client)
		default:
			err = fmt.Errorf("unsupported media type: %s", mediaType)
		}
		
		if err != nil {
			errorMsg := fmt.Sprintf("Failed to sync %s: %v", mediaType, err)
			j.jobRepo.CompleteJobRun(syncCtx, jobRun.ID, models.JobStatusFailed, errorMsg)
			log.Printf("Error: %s", errorMsg)
			return
		}
		
		// Complete the job
		j.jobRepo.UpdateJobProgress(syncCtx, jobRun.ID, 100, fmt.Sprintf("%s sync completed successfully", mediaType))
		j.jobRepo.CompleteJobRun(syncCtx, jobRun.ID, models.JobStatusCompleted, "")
		
		log.Printf("%s sync for user %d from client %d completed", mediaType, userID, clientID)
	}()
	
	return nil
}

// getMediaClient retrieves the media client interface from the factory
func (j *MediaSyncJob) getMediaClient(ctx context.Context, clientID uint64) (media.MediaClient, error) {
	// In a real implementation, we would get the client from the database
	// For now, we'll just return an error
	
	// Use the client factory service to get or create the client
	// Currently config doesn't implement ClientConfig, so we'll use a placeholder
	// until proper implementation is available
	// For now, return an error to indicate this is not fully implemented
	return nil, fmt.Errorf("media client implementation not completed")
}

// clientSupportsMediaType checks if the client supports the specified media type
func (j *MediaSyncJob) clientSupportsMediaType(client media.MediaClient, mediaType string) bool {
	if client == nil {
		return false
	}
	
	switch mediaType {
	case "movie":
		return client.SupportsMovies()
	case "series":
		return client.SupportsSeries()
	case "music":
		return client.SupportsMusic()
	case "playlist":
		return client.SupportsPlaylists()
	case "collection":
		return client.SupportsCollections()
	default:
		return false
	}
}

// syncMovies syncs movies from the given client
func (j *MediaSyncJob) syncMovies(ctx context.Context, jobRunID uint64, client media.MediaClient) error {
	// In a full implementation we would need to cast client to the proper provider
	// and get movies from client
	
	// For now, return a placeholder since we haven't fully implemented all interfaces
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 15, "Fetching movies from client")
	movies := []models.MediaItem[*mediatypes.Movie]{}
	
	// Placeholder for getting movies from client
	// movies, err := movieProvider.GetMovies(ctx, &mediatypes.QueryOptions{})
	// if err != nil {
	//     return fmt.Errorf("failed to get movies: %w", err)
	// }
	
	// Update progress and total items count
	totalItems := len(movies)
	j.jobRepo.SetJobTotalItems(ctx, jobRunID, totalItems)
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 20, fmt.Sprintf("Found %d movies to process", totalItems))
	
	if totalItems == 0 {
		return nil // No movies to process
	}
	
	// In a real implementation, we would process movies in batches
	// For now, just simulate progress
	
	// Simulate processing
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 50, "Processing movies")
	time.Sleep(100 * time.Millisecond)
	
	// Simulate completion
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 80, "Finishing movie sync")
	
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 95, "Finalizing movie sync")
	return nil
}

// syncSeries syncs TV series from the given client
func (j *MediaSyncJob) syncSeries(ctx context.Context, jobRunID uint64, client media.MediaClient) error {
	// In a full implementation we would need to cast client to the proper provider
	// and get series from client
	
	// For now, return a placeholder since we haven't fully implemented all interfaces
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 15, "Fetching series from client")
	series := []models.MediaItem[*mediatypes.Series]{}
	
	// Placeholder for getting series from client
	// series, err := seriesProvider.GetSeries(ctx, &mediatypes.QueryOptions{})
	// if err != nil {
	//     return fmt.Errorf("failed to get series: %w", err)
	// }
	
	// Update progress and total items count
	totalItems := len(series)
	j.jobRepo.SetJobTotalItems(ctx, jobRunID, totalItems)
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 20, fmt.Sprintf("Found %d series to process", totalItems))
	
	if totalItems == 0 {
		return nil // No series to process
	}
	
	// In a real implementation, we would process series in batches
	// For now, just simulate progress
	
	// Simulate processing
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 50, "Processing series")
	time.Sleep(100 * time.Millisecond)
	
	// Simulate completion
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 80, "Finishing series sync")
	
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 95, "Finalizing series sync")
	return nil
}

// syncMusic syncs music tracks from the given client
func (j *MediaSyncJob) syncMusic(ctx context.Context, jobRunID uint64, client media.MediaClient) error {
	// In a full implementation we would need to cast client to the proper provider
	// and get music from client
	
	// For now, return a placeholder since we haven't fully implemented all interfaces
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 15, "Fetching music tracks from client")
	tracks := []models.MediaItem[*mediatypes.Track]{}
	
	// Placeholder for getting tracks from client
	// tracks, err := musicProvider.GetMusic(ctx, &mediatypes.QueryOptions{})
	// if err != nil {
	//     return fmt.Errorf("failed to get music tracks: %w", err)
	// }
	
	// Update progress and total items count
	totalItems := len(tracks)
	j.jobRepo.SetJobTotalItems(ctx, jobRunID, totalItems)
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 20, fmt.Sprintf("Found %d tracks to process", totalItems))
	
	if totalItems == 0 {
		return nil // No tracks to process
	}
	
	// In a real implementation, we would process tracks in batches
	// For now, just simulate progress
	
	// Simulate processing
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 50, "Processing tracks")
	time.Sleep(100 * time.Millisecond)
	
	// Simulate completion
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 80, "Finishing track sync")
	
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 95, "Finalizing music sync")
	return nil
}

// syncPlaylists syncs playlists from the given client
func (j *MediaSyncJob) syncPlaylists(ctx context.Context, jobRunID uint64, client media.MediaClient) error {
	// In a full implementation we would need to cast client to the proper provider
	// and get playlists from client
	
	// For now, return a placeholder since we haven't fully implemented all interfaces
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 15, "Fetching playlists from client")
	playlists := []models.MediaItem[*mediatypes.Playlist]{}
	
	// Placeholder for getting playlists from client
	// playlists, err := playlistProvider.GetPlaylists(ctx, &mediatypes.QueryOptions{})
	// if err != nil {
	//     return fmt.Errorf("failed to get playlists: %w", err)
	// }
	
	// Update progress and total items count
	totalItems := len(playlists)
	j.jobRepo.SetJobTotalItems(ctx, jobRunID, totalItems)
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 20, fmt.Sprintf("Found %d playlists to process", totalItems))
	
	if totalItems == 0 {
		return nil // No playlists to process
	}
	
	// Log playlists since we don't have a repository for them yet
	log.Printf("Found %d playlists", totalItems)
	for i, playlist := range playlists {
		log.Printf("Playlist %d: %s, Items: %d", i+1, playlist.Data.Details.Title, playlist.Data.ItemCount)
	}
	
	// Create a note that we need to implement playlist storage
	log.Printf("NOTE: Playlist storage not yet implemented - %d playlists found but not stored", totalItems)
	
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 95, "Finalizing playlist sync")
	return nil
}

// syncCollections syncs collections from the given client
func (j *MediaSyncJob) syncCollections(ctx context.Context, jobRunID uint64, client media.MediaClient) error {
	// In a full implementation we would need to cast client to the proper provider
	// and get collections from client
	
	// For now, return a placeholder since we haven't fully implemented all interfaces
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 15, "Fetching collections from client")
	collections := []models.MediaItem[*mediatypes.Collection]{}
	
	// Placeholder for getting collections from client
	// collections, err := collectionProvider.GetCollections(ctx, &mediatypes.QueryOptions{})
	// if err != nil {
	//     return fmt.Errorf("failed to get collections: %w", err)
	// }
	
	// Update progress and total items count
	totalItems := len(collections)
	j.jobRepo.SetJobTotalItems(ctx, jobRunID, totalItems)
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 20, fmt.Sprintf("Found %d collections to process", totalItems))
	
	if totalItems == 0 {
		return nil // No collections to process
	}
	
	// Log collections since we don't have a repository for them yet
	log.Printf("Found %d collections", totalItems)
	for i, collection := range collections {
		log.Printf("Collection %d: %s, Type: %s, Items: %d", 
			i+1, collection.Data.Details.Title, collection.Data.CollectionType, collection.Data.ItemCount)
	}
	
	// Create a note that we need to implement collection storage
	log.Printf("NOTE: Collection storage not yet implemented - %d collections found but not stored", totalItems)
	
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 95, "Finalizing collection sync")
	return nil
}