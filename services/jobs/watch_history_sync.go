package jobs

import (
	"context"
	"fmt"
	"log"
	"suasor/client"
	"suasor/client/media"
	"suasor/client/media/providers"
	mediatypes "suasor/client/media/types"
	clienttypes "suasor/client/types"
	"suasor/repository"
	"suasor/services/scheduler"
	"suasor/types/models"
	"suasor/utils"
	"time"
)

// MediaClientInfo is defined in common.go

// WatchHistorySyncJob synchronizes watched media history from external clients
type WatchHistorySyncJob struct {
	jobRepo         repository.JobRepository
	userRepo        repository.UserRepository
	configRepo      repository.UserConfigRepository
	historyRepo     repository.MediaPlayHistoryRepository
	movieRepo       repository.MediaItemRepository[*mediatypes.Movie]
	seriesRepo      repository.MediaItemRepository[*mediatypes.Series]
	episodeRepo     repository.MediaItemRepository[*mediatypes.Episode]
	musicRepo       repository.MediaItemRepository[*mediatypes.Track]
	clientRepos     repository.ClientRepositoryCollection
	clientFactories *client.ClientFactoryService
}

// NewWatchHistorySyncJob creates a new watch history sync job
func NewWatchHistorySyncJob(
	jobRepo repository.JobRepository,
	userRepo repository.UserRepository,
	configRepo repository.UserConfigRepository,
	historyRepo repository.MediaPlayHistoryRepository,
	movieRepo repository.MediaItemRepository[*mediatypes.Movie],
	seriesRepo repository.MediaItemRepository[*mediatypes.Series],
	episodeRepo repository.MediaItemRepository[*mediatypes.Episode],
	musicRepo repository.MediaItemRepository[*mediatypes.Track],
	clientRepos repository.ClientRepositoryCollection,
	clientFactories *client.ClientFactoryService,
) *WatchHistorySyncJob {
	return &WatchHistorySyncJob{
		jobRepo:         jobRepo,
		userRepo:        userRepo,
		configRepo:      configRepo,
		historyRepo:     historyRepo,
		movieRepo:       movieRepo,
		seriesRepo:      seriesRepo,
		episodeRepo:     episodeRepo,
		musicRepo:       musicRepo,
		clientRepos:     clientRepos,
		clientFactories: clientFactories,
	}
}

// Name returns the unique name of the job
func (j *WatchHistorySyncJob) Name() string {
	return "system.history.sync"
}

// Schedule returns when the job should next run
func (j *WatchHistorySyncJob) Schedule() time.Duration {
	// Default to checking daily
	return 24 * time.Hour
}

// Execute runs the watch history sync job
func (j *WatchHistorySyncJob) Execute(ctx context.Context) error {
	log.Println("Starting watch history sync job")

	// Get all users
	users, err := j.userRepo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("error getting users: %w", err)
	}

	// Process each user
	for _, user := range users {
		if err := j.processUserHistory(ctx, user); err != nil {
			log.Printf("Error processing history for user %s: %v", user.Username, err)
			// Continue with other users even if one fails
			continue
		}
	}

	log.Println("Watch history sync job completed")
	return nil
}

// processUserHistory syncs watch history for a single user
func (j *WatchHistorySyncJob) processUserHistory(ctx context.Context, user models.User) error {
	// Skip inactive users
	if !user.Active {
		log.Printf("Skipping inactive user: %s", user.Username)
		return nil
	}

	// Get user configuration
	_, err := j.configRepo.GetUserConfig(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("error getting user config: %w", err)
	}

	// Create a job run record for this user
	now := time.Now()
	jobRun := &models.JobRun{
		JobName:   j.Name(),
		JobType:   models.JobTypeSync,
		Status:    models.JobStatusRunning,
		StartTime: &now,
		UserID:    &user.ID,
		Metadata:  fmt.Sprintf(`{"userId":%d,"username":"%s","type":"watchHistory"}`, user.ID, user.Username),
	}

	if err := j.jobRepo.CreateJobRun(ctx, jobRun); err != nil {
		log.Printf("Error creating job run record: %v", err)
		return err
	}

	// Get all media clients for the user
	clients, err := j.getUserMediaClients(ctx, user.ID)
	if err != nil {
		errorMsg := fmt.Sprintf("Error getting media clients: %v", err)
		j.completeJobRun(ctx, jobRun.ID, models.JobStatusFailed, errorMsg)
		return fmt.Errorf(errorMsg)
	}

	if len(clients) == 0 {
		j.jobRepo.UpdateJobProgress(ctx, jobRun.ID, 100, "No media clients found")
		j.completeJobRun(ctx, jobRun.ID, models.JobStatusCompleted, "No media clients found")
		return nil
	}

	j.jobRepo.UpdateJobProgress(ctx, jobRun.ID, 10, fmt.Sprintf("Found %d media clients", len(clients)))

	// Process each client
	totalClients := len(clients)
	processedClients := 0
	var lastError error

	for i, client := range clients {
		// Update progress
		progress := 10 + int(float64(i)/float64(totalClients)*80.0)
		j.jobRepo.UpdateJobProgress(ctx, jobRun.ID, progress,
			fmt.Sprintf("Processing client %d/%d: %s", i+1, totalClients, client.Name))

		// Sync watch history for this client
		err := j.syncClientHistory(ctx, user.ID, client, jobRun.ID)
		if err != nil {
			log.Printf("Error syncing history for client %s: %v", client.Name, err)
			lastError = err
			continue
		}

		processedClients++
	}

	// Complete the job
	if lastError != nil {
		errorMsg := fmt.Sprintf("Completed with errors: %v", lastError)
		j.jobRepo.UpdateJobProgress(ctx, jobRun.ID, 100,
			fmt.Sprintf("Processed %d/%d clients with errors", processedClients, totalClients))
		j.completeJobRun(ctx, jobRun.ID, models.JobStatusFailed, errorMsg)
		return lastError
	}

	j.jobRepo.UpdateJobProgress(ctx, jobRun.ID, 100,
		fmt.Sprintf("Successfully processed %d/%d clients", processedClients, totalClients))
	j.completeJobRun(ctx, jobRun.ID, models.JobStatusCompleted, "")
	return nil
}

// completeJobRun finalizes a job run with status and error info
func (j *WatchHistorySyncJob) completeJobRun(ctx context.Context, jobRunID uint64, status models.JobStatus, errorMsg string) {
	if err := j.jobRepo.CompleteJobRun(ctx, jobRunID, status, errorMsg); err != nil {
		log.Printf("Error completing job run: %v", err)
	}
}

// SetupWatchHistorySyncSchedule creates or updates a watch history sync schedule for a user
func (j *WatchHistorySyncJob) SetupWatchHistorySyncSchedule(ctx context.Context, userID uint64, frequency string) error {
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

// RunManualSync runs the watch history sync job manually for a specific user
func (j *WatchHistorySyncJob) RunManualSync(ctx context.Context, userID uint64) error {
	// Get the user
	user, err := j.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("error getting user: %w", err)
	}

	if user == nil {
		return fmt.Errorf("user not found: %d", userID)
	}

	// Run the sync job for this user
	return j.processUserHistory(ctx, *user)
}

// getUserMediaClients returns all media clients for a user
func (j *WatchHistorySyncJob) getUserMediaClients(ctx context.Context, userID uint64) ([]MediaClientInfo, error) {
	log := utils.LoggerFromContext(ctx)
	log.Info().Uint64("userID", userID).Msg("Getting media clients for user")

	var clients []MediaClientInfo

	// Get all media client types from the repository collection
	mediaCategoryClients := j.clientRepos.GetAllByCategory(clienttypes.ClientCategoryMedia)

	// Emby clients
	embyClients, err := mediaCategoryClients.EmbyRepo.GetByUserID(ctx, userID)
	if err != nil {
		log.Warn().Err(err).Msg("Error getting Emby clients")
	} else {
		for _, c := range embyClients {
			clients = append(clients, MediaClientInfo{
				ClientID:   c.ID,
				ClientType: clienttypes.MediaClientTypeEmby,
				Name:       c.Name,
				UserID:     userID,
			})
		}
	}

	// Jellyfin clients
	jellyfinClients, err := mediaCategoryClients.JellyfinRepo.GetByUserID(ctx, userID)
	if err != nil {
		log.Warn().Err(err).Msg("Error getting Jellyfin clients")
	} else {
		for _, c := range jellyfinClients {
			clients = append(clients, MediaClientInfo{
				ClientID:   c.ID,
				ClientType: clienttypes.MediaClientTypeJellyfin,
				Name:       c.Name,
				UserID:     userID,
			})
		}
	}

	// Plex clients
	plexClients, err := mediaCategoryClients.PlexRepo.GetByUserID(ctx, userID)
	if err != nil {
		log.Warn().Err(err).Msg("Error getting Plex clients")
	} else {
		for _, c := range plexClients {
			clients = append(clients, MediaClientInfo{
				ClientID:   c.ID,
				ClientType: clienttypes.MediaClientTypePlex,
				Name:       c.Name,
				UserID:     userID,
			})
		}
	}

	// Subsonic clients (primarily for music)
	subsonicClients, err := mediaCategoryClients.SubsonicRepo.GetByUserID(ctx, userID)
	if err != nil {
		log.Warn().Err(err).Msg("Error getting Subsonic clients")
	} else {
		for _, c := range subsonicClients {
			clients = append(clients, MediaClientInfo{
				ClientID:   c.ID,
				ClientType: clienttypes.MediaClientTypeSubsonic,
				Name:       c.Name,
				UserID:     userID,
			})
		}
	}

	log.Info().Int("clientCount", len(clients)).Msg("Found media clients")
	return clients, nil
}

// syncClientHistory syncs watch history for a specific client
func (j *WatchHistorySyncJob) syncClientHistory(ctx context.Context, userID uint64, client MediaClientInfo, jobRunID uint64) error {
	// Log the start of synchronization
	log := utils.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", client.ClientID).
		Str("clientName", client.Name).
		Str("clientType", string(client.ClientType)).
		Msg("Syncing watch history")

	// Get the client using client factory
	mediaClient, err := j.getMediaClient(ctx, client.ClientID, string(client.ClientType))
	if err != nil {
		return fmt.Errorf("failed to get media client: %w", err)
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 20, "Fetching play history from client")

	// Check if client supports play history
	historyProvider, ok := mediaClient.(providers.HistoryProvider)
	if !ok {
		return fmt.Errorf("client doesn't support play history")
	}

	// Skip if client doesn't support history
	if !historyProvider.SupportsHistory() {
		log.Warn().
			Str("clientName", client.Name).
			Msg("Client doesn't support history - skipping")
		return nil
	}

	// Get play history items from the client
	playHistory, err := historyProvider.GetPlayHistory(ctx, &mediatypes.QueryOptions{
		Limit:     100, // Limit to 100 most recent items
		Sort:      "lastPlayedDate",
		SortOrder: "desc",
	})

	if err != nil {
		return fmt.Errorf("failed to get play history: %w", err)
	}

	// Update progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 40,
		fmt.Sprintf("Processing %d history items from client", len(playHistory)))

	// Process each history item
	processedItems := 0
	for i, historyItem := range playHistory {
		// Update detailed progress periodically
		if i%10 == 0 {
			progress := 40 + int(float64(i)/float64(len(playHistory))*50.0)
			j.jobRepo.UpdateJobProgress(ctx, jobRunID, progress,
				fmt.Sprintf("Processed %d/%d history items", i, len(playHistory)))
		}

		// Skip invalid items
		if historyItem.Item == nil || historyItem.Item.Data == nil {
			log.Warn().Msg("Skipping invalid history item with no data")
			continue
		}

		// Create/update history based on media type
		switch historyItem.Item.Type {
		case mediatypes.MediaTypeMovie:
			if err := j.processMovieHistory(ctx, userID, client.ClientID, historyItem); err != nil {
				log.Warn().Err(err).Msg("Error processing movie history")
				continue
			}
		case mediatypes.MediaTypeSeries:
			if err := j.processSeriesHistory(ctx, userID, client.ClientID, historyItem); err != nil {
				log.Warn().Err(err).Msg("Error processing series history")
				continue
			}
		case mediatypes.MediaTypeEpisode:
			if err := j.processEpisodeHistory(ctx, userID, client.ClientID, historyItem); err != nil {
				log.Warn().Err(err).Msg("Error processing episode history")
				continue
			}
		case mediatypes.MediaTypeTrack:
			if err := j.processMusicHistory(ctx, userID, client.ClientID, historyItem); err != nil {
				log.Warn().Err(err).Msg("Error processing music history")
				continue
			}
		default:
			log.Debug().
				Str("mediaType", string(historyItem.Item.Type)).
				Msg("Unsupported media type in history - skipping")
			continue
		}

		processedItems++
	}

	// Updated progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 90,
		fmt.Sprintf("Successfully processed %d/%d history items", processedItems, len(playHistory)))

	return nil
}

// processMovieHistory processes a movie history item and updates the database
func (j *WatchHistorySyncJob) processMovieHistory(ctx context.Context, userID, clientID uint64, historyItem models.MediaPlayHistory[mediatypes.MediaData]) error {
	log := utils.LoggerFromContext(ctx)

	// Get client item ID from the media item
	var clientItemID string
	for _, cid := range historyItem.Item.ClientIDs {
		if cid.ID == clientID {
			clientItemID = cid.ItemID
			break
		}
	}

	if clientItemID == "" {
		return fmt.Errorf("no client item ID found for movie history")
	}

	// Look up the movie in our database
	movieItem, err := j.movieRepo.GetByClientItemID(ctx, clientItemID, clientID)
	if err != nil {
		// If we can't find the movie, it might not be synced yet
		log.Warn().
			Str("clientItemID", clientItemID).
			Uint64("clientID", clientID).
			Msg("Movie not found in database for history - consider running media sync first")
		return err
	}

	// Create history record
	historyRecord := models.MediaPlayHistory[*mediatypes.Movie]{
		MediaItemID:      movieItem.ID,
		Type:             string(mediatypes.MediaTypeMovie),
		WatchedAt:        historyItem.WatchedAt,
		LastWatchedAt:    historyItem.LastWatchedAt,
		PlayedPercentage: historyItem.PlayedPercentage,
		PlayCount:        historyItem.PlayCount,
		PositionSeconds:  historyItem.PositionSeconds,
		DurationSeconds:  historyItem.DurationSeconds,
		Completed:        historyItem.PlayedPercentage >= 90, // Mark as completed if watched 90% or more
	}

	// Associate with the movie item
	historyRecord.Associate(movieItem)

	// Save to database
	err = j.historyRepo.CreateHistory(ctx, &historyRecord)
	if err != nil {
		return fmt.Errorf("failed to save movie history: %w", err)
	}

	log.Debug().
		Str("movieTitle", movieItem.Title).
		Time("watchedAt", historyItem.WatchedAt).
		Float64("percentage", historyItem.PlayedPercentage).
		Msg("Saved movie watch history")

	return nil
}

// processSeriesHistory processes a series history item and updates the database
func (j *WatchHistorySyncJob) processSeriesHistory(ctx context.Context, userID, clientID uint64, historyItem models.MediaPlayHistory[mediatypes.MediaData]) error {
	log := utils.LoggerFromContext(ctx)

	// Get client item ID from the media item
	var clientItemID string
	for _, cid := range historyItem.Item.ClientIDs {
		if cid.ID == clientID {
			clientItemID = cid.ItemID
			break
		}
	}

	if clientItemID == "" {
		return fmt.Errorf("no client item ID found for series history")
	}

	// Look up the series in our database
	seriesItem, err := j.seriesRepo.GetByClientItemID(ctx, clientItemID, clientID)
	if err != nil {
		// If we can't find the series, it might not be synced yet
		log.Warn().
			Str("clientItemID", clientItemID).
			Uint64("clientID", clientID).
			Msg("Series not found in database for history - consider running media sync first")
		return err
	}

	// Create history record
	historyRecord := models.MediaPlayHistory[*mediatypes.Series]{
		MediaItemID:      seriesItem.ID,
		Type:             string(mediatypes.MediaTypeSeries),
		WatchedAt:        historyItem.WatchedAt,
		LastWatchedAt:    historyItem.LastWatchedAt,
		PlayedPercentage: historyItem.PlayedPercentage,
		PlayCount:        historyItem.PlayCount,
		PositionSeconds:  historyItem.PositionSeconds,
		DurationSeconds:  historyItem.DurationSeconds,
		Completed:        historyItem.PlayedPercentage >= 90, // Mark as completed if watched 90% or more
	}

	// Associate with the series item
	historyRecord.Associate(seriesItem)

	// Save to database
	err = j.historyRepo.CreateHistory(ctx, &historyRecord)
	if err != nil {
		return fmt.Errorf("failed to save series history: %w", err)
	}

	log.Debug().
		Str("seriesTitle", seriesItem.Title).
		Time("watchedAt", historyItem.WatchedAt).
		Float64("percentage", historyItem.PlayedPercentage).
		Msg("Saved series watch history")

	return nil
}

// processEpisodeHistory processes an episode history item and updates the database
func (j *WatchHistorySyncJob) processEpisodeHistory(ctx context.Context, userID, clientID uint64, historyItem models.MediaPlayHistory[mediatypes.MediaData]) error {
	log := utils.LoggerFromContext(ctx)

	// Get client item ID from the media item
	var clientItemID string
	for _, cid := range historyItem.Item.ClientIDs {
		if cid.ID == clientID {
			clientItemID = cid.ItemID
			break
		}
	}

	if clientItemID == "" {
		return fmt.Errorf("no client item ID found for episode history")
	}

	// Look up the episode in our database
	episodeItem, err := j.episodeRepo.GetByClientItemID(ctx, clientItemID, clientID)
	if err != nil {
		// If we can't find the episode, it might not be synced yet
		log.Warn().
			Str("clientItemID", clientItemID).
			Uint64("clientID", clientID).
			Msg("Episode not found in database for history - consider running media sync first")
		return err
	}

	// Create history record
	historyRecord := models.MediaPlayHistory[*mediatypes.Episode]{
		MediaItemID:      episodeItem.ID,
		Type:             string(mediatypes.MediaTypeEpisode),
		WatchedAt:        historyItem.WatchedAt,
		LastWatchedAt:    historyItem.LastWatchedAt,
		PlayedPercentage: historyItem.PlayedPercentage,
		PlayCount:        historyItem.PlayCount,
		PositionSeconds:  historyItem.PositionSeconds,
		DurationSeconds:  historyItem.DurationSeconds,
		Completed:        historyItem.PlayedPercentage >= 90, // Mark as completed if watched 90% or more
	}

	// Associate with the episode item
	historyRecord.Associate(episodeItem)

	// Save to database
	err = j.historyRepo.CreateHistory(ctx, &historyRecord)
	if err != nil {
		return fmt.Errorf("failed to save episode history: %w", err)
	}

	log.Debug().
		Str("episodeTitle", episodeItem.Title).
		Time("watchedAt", historyItem.WatchedAt).
		Float64("percentage", historyItem.PlayedPercentage).
		Msg("Saved episode watch history")

	return nil
}

// processMusicHistory processes a music track history item and updates the database
func (j *WatchHistorySyncJob) processMusicHistory(ctx context.Context, userID, clientID uint64, historyItem models.MediaPlayHistory[mediatypes.MediaData]) error {
	log := utils.LoggerFromContext(ctx)

	// Get client item ID from the media item
	var clientItemID string
	for _, cid := range historyItem.Item.ClientIDs {
		if cid.ID == clientID {
			clientItemID = cid.ItemID
			break
		}
	}

	if clientItemID == "" {
		return fmt.Errorf("no client item ID found for music history")
	}

	// Look up the track in our database
	trackItem, err := j.musicRepo.GetByClientItemID(ctx, clientItemID, clientID)
	if err != nil {
		// If we can't find the track, it might not be synced yet
		log.Warn().
			Str("clientItemID", clientItemID).
			Uint64("clientID", clientID).
			Msg("Track not found in database for history - consider running media sync first")
		return err
	}

	// Create history record
	historyRecord := models.MediaPlayHistory[*mediatypes.Track]{
		MediaItemID:      trackItem.ID,
		Type:             string(mediatypes.MediaTypeTrack),
		WatchedAt:        historyItem.WatchedAt,
		LastWatchedAt:    historyItem.LastWatchedAt,
		PlayedPercentage: historyItem.PlayedPercentage,
		PlayCount:        historyItem.PlayCount,
		PositionSeconds:  historyItem.PositionSeconds,
		DurationSeconds:  historyItem.DurationSeconds,
		Completed:        historyItem.PlayedPercentage >= 90, // Mark as completed if listened to 90% or more
	}

	// Associate with the track item
	historyRecord.Associate(trackItem)

	// Save to database
	err = j.historyRepo.CreateHistory(ctx, &historyRecord)
	if err != nil {
		return fmt.Errorf("failed to save track history: %w", err)
	}

	log.Debug().
		Str("trackTitle", trackItem.Title).
		Time("playedAt", historyItem.WatchedAt).
		Float64("percentage", historyItem.PlayedPercentage).
		Msg("Saved track play history")

	return nil
}

// getMediaClient gets a media client from the client factory
func (j *WatchHistorySyncJob) getMediaClient(ctx context.Context, clientID uint64, clientType string) (media.MediaClient, error) {
	log := utils.LoggerFromContext(ctx)

	// Get the client config from the repository
	var clientConfig clienttypes.ClientConfig

	mediaCategoryClients := j.clientRepos.GetAllByCategory(clienttypes.ClientCategoryMedia)

	switch clienttypes.MediaClientType(clientType) {
	case clienttypes.MediaClientTypeEmby:
		c, err := mediaCategoryClients.EmbyRepo.GetByID(ctx, clientID)
		if err != nil {
			return nil, fmt.Errorf("failed to get emby client: %w", err)
		}
		clientConfig = c.GetConfig()
	case clienttypes.MediaClientTypeJellyfin:
		c, err := mediaCategoryClients.JellyfinRepo.GetByID(ctx, clientID)
		if err != nil {
			return nil, fmt.Errorf("failed to get jellyfin client: %w", err)
		}
		clientConfig = c.GetConfig()
	case clienttypes.MediaClientTypePlex:
		c, err := mediaCategoryClients.PlexRepo.GetByID(ctx, clientID)
		if err != nil {
			return nil, fmt.Errorf("failed to get plex client: %w", err)
		}
		clientConfig = c.GetConfig()
	case clienttypes.MediaClientTypeSubsonic:
		c, err := mediaCategoryClients.SubsonicRepo.GetByID(ctx, clientID)
		if err != nil {
			return nil, fmt.Errorf("failed to get subsonic client: %w", err)
		}
		clientConfig = c.GetConfig()
	default:
		return nil, fmt.Errorf("unsupported client type: %s", clientType)
	}

	// Get the client instance
	clientInstance, err := j.clientFactories.GetClient(ctx, clientID, clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	// Cast to media client
	mediaClient, ok := clientInstance.(media.MediaClient)
	if !ok {
		return nil, fmt.Errorf("client is not a media client")
	}

	log.Info().
		Uint64("clientID", clientID).
		Str("clientType", clientType).
		Msg("Successfully retrieved and initialized media client")

	return mediaClient, nil
}

