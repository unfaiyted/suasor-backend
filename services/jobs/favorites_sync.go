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
	"suasor/utils"
)

// Using MediaClientInfo from common.go

// FavoritesSyncJob synchronizes favorite/liked media from external clients
type FavoritesSyncJob struct {
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

// NewFavoritesSyncJob creates a new favorites sync job
func NewFavoritesSyncJob(
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
) *FavoritesSyncJob {
	return &FavoritesSyncJob{
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
func (j *FavoritesSyncJob) Name() string {
	return "system.favorites.sync"
}

// Schedule returns when the job should next run
func (j *FavoritesSyncJob) Schedule() time.Duration {
	// Default to checking daily
	return 24 * time.Hour
}

// Execute runs the favorites sync job
func (j *FavoritesSyncJob) Execute(ctx context.Context) error {
	log.Println("Starting favorites sync job")

	// Get all users
	users, err := j.userRepo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("error getting users: %w", err)
	}

	// Process each user
	for _, user := range users {
		if err := j.processUserFavorites(ctx, user); err != nil {
			log.Printf("Error processing favorites for user %s: %v", user.Username, err)
			// Continue with other users even if one fails
			continue
		}
	}

	log.Println("Favorites sync job completed")
	return nil
}

// processUserFavorites syncs favorites for a single user
func (j *FavoritesSyncJob) processUserFavorites(ctx context.Context, user models.User) error {
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

	// Check if favorites sync is enabled for the user
	// First check if sync notifications are enabled as a proxy for sync being enabled
	if !config.NotifyOnSync {
		log.Printf("Favorites sync not enabled for user: %s", user.Username)
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
		Metadata:  fmt.Sprintf(`{"userId":%d,"username":"%s","type":"favorites"}`, user.ID, user.Username),
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

		// Sync favorites for this client
		err := j.syncClientFavorites(ctx, user.ID, client, jobRun.ID)
		if err != nil {
			log.Printf("Error syncing favorites for client %s: %v", client.Name, err)
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
func (j *FavoritesSyncJob) completeJobRun(ctx context.Context, jobRunID uint64, status models.JobStatus, errorMsg string) {
	if err := j.jobRepo.CompleteJobRun(ctx, jobRunID, status, errorMsg); err != nil {
		log.Printf("Error completing job run: %v", err)
	}
}

// SetupFavoritesSyncSchedule creates or updates a favorites sync schedule for a user
func (j *FavoritesSyncJob) SetupFavoritesSyncSchedule(ctx context.Context, userID uint64, frequency string) error {
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

// RunManualSync runs the favorites sync job manually for a specific user
func (j *FavoritesSyncJob) RunManualSync(ctx context.Context, userID uint64) error {
	// Get the user
	user, err := j.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("error getting user: %w", err)
	}

	if user == nil {
		return fmt.Errorf("user not found: %d", userID)
	}

	// Run the sync job for this user
	return j.processUserFavorites(ctx, *user)
}

// getUserMediaClients returns all media clients for a user
func (j *FavoritesSyncJob) getUserMediaClients(ctx context.Context, userID uint64) ([]MediaClientInfo, error) {
	log := utils.LoggerFromContext(ctx)
	log.Info().Uint64("userID", userID).Msg("Getting media clients for user")

	var clients []MediaClientInfo

	// Get all media client types from the repository collection
	mediaCategoryClients := j.clientRepos.GetAllByCategory(ctx, clienttypes.ClientCategoryMedia)

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

// syncClientFavorites syncs favorites for a specific client
func (j *FavoritesSyncJob) syncClientFavorites(ctx context.Context, userID uint64, client MediaClientInfo, jobRunID uint64) error {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", client.ClientID).
		Str("clientName", client.Name).
		Str("clientType", string(client.ClientType)).
		Msg("Syncing favorites")

	// Get the client using client factory
	mediaClient, err := j.getMediaClient(ctx, client.ClientID, string(client.ClientType))
	if err != nil {
		return fmt.Errorf("failed to get media client: %w", err)
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 20, "Fetching favorites from client")

	// Process different media types
	totalProcessed := 0

	// Sync movie favorites
	movieCount, err := j.syncMovieFavorites(ctx, userID, client.ClientID, mediaClient, jobRunID)
	if err != nil {
		log.Warn().Err(err).Msg("Error syncing movie favorites")
	}
	totalProcessed += movieCount

	// Sync series favorites
	seriesCount, err := j.syncSeriesFavorites(ctx, userID, client.ClientID, mediaClient, jobRunID)
	if err != nil {
		log.Warn().Err(err).Msg("Error syncing series favorites")
	}
	totalProcessed += seriesCount

	// Sync episode favorites
	episodeCount, err := j.syncEpisodeFavorites(ctx, userID, client.ClientID, mediaClient, jobRunID)
	if err != nil {
		log.Warn().Err(err).Msg("Error syncing episode favorites")
	}
	totalProcessed += episodeCount

	// Sync music favorites
	musicCount, err := j.syncMusicFavorites(ctx, userID, client.ClientID, mediaClient, jobRunID)
	if err != nil {
		log.Warn().Err(err).Msg("Error syncing music favorites")
	}
	totalProcessed += musicCount

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 90,
		fmt.Sprintf("Synced %d favorite items", totalProcessed))

	return nil
}

// syncMovieFavorites syncs favorite movies
func (j *FavoritesSyncJob) syncMovieFavorites(ctx context.Context, userID, clientID uint64, mediaClient media.MediaClient, jobRunID uint64) (int, error) {
	log := utils.LoggerFromContext(ctx)

	// Check if client supports movies
	movieProvider, ok := mediaClient.(providers.MovieProvider)
	if !ok {
		log.Debug().Msg("Client doesn't support movies, skipping movie favorites")
		return 0, nil
	}

	// Get favorite movies
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 30, "Fetching favorite movies")

	// Get movies from client with favorites filter
	movies, err := movieProvider.GetMovies(ctx, &mediatypes.QueryOptions{
		Limit:     100,
		Favorites: true,
	})

	if err != nil {
		return 0, fmt.Errorf("failed to get favorite movies: %w", err)
	}

	// Update progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 40,
		fmt.Sprintf("Processing %d favorite movies", len(movies)))

	// Track count of processed items
	processed := 0

	// Process each favorite movie
	for _, movie := range movies {
		// Get client item ID for this client
		var clientItemID string
		for _, cid := range movie.SyncClients {
			if cid.ID == clientID {
				clientItemID = cid.ItemID
				break
			}
		}

		if clientItemID == "" {
			log.Warn().
				Str("movieTitle", movie.Data.Details.Title).
				Msg("No matching client item ID found for movie")
			continue
		}

		// Get the movie from database
		existingMovie, err := j.movieRepo.GetByClientItemID(ctx, clientItemID, clientID)
		if err != nil {
			log.Warn().
				Err(err).
				Str("clientItemID", clientItemID).
				Str("movieTitle", movie.Data.Details.Title).
				Msg("Movie not found in database")
			continue
		}

		// Create or update history record to mark as favorite
		// First, check if there's already a history record
		hasViewed, err := j.historyRepo.HasUserViewedMedia(ctx, userID, existingMovie.ID)

		// Create or update history
		if hasViewed {
			// Get existing history (we'll need to implement a method to get history by user and media item)
			playCount, _ := j.historyRepo.GetItemPlayCount(ctx, userID, existingMovie.ID)

			// Create new history record (timestamp is needed since we don't have a way to update existing record)
			historyRecord := models.MediaPlayHistory[*mediatypes.Movie]{
				MediaItemID:      existingMovie.ID,
				Type:             mediatypes.MediaTypeMovie,
				PlayedAt:         time.Now(),
				LastPlayedAt:     time.Now(),
				IsFavorite:       true,
				PlayCount:        int32(playCount),
				PlayedPercentage: 100, // Assume watched if favorited
				Completed:        true,
			}

			// Associate with the movie item
			historyRecord.Associate(existingMovie)

			// Save to database
			err = j.historyRepo.CreateHistory(ctx, &historyRecord)
			if err != nil {
				log.Warn().
					Err(err).
					Str("movieTitle", existingMovie.Title).
					Msg("Failed to update movie favorite status in history")
				continue
			}
		} else {
			// Create new history record for a favorited but never watched item
			historyRecord := models.MediaPlayHistory[*mediatypes.Movie]{
				MediaItemID:      existingMovie.ID,
				Type:             mediatypes.MediaTypeMovie,
				PlayedAt:         time.Now(),
				LastPlayedAt:     time.Now(),
				IsFavorite:       true,
				PlayCount:        1,   // First view
				PlayedPercentage: 100, // Assume watched if favorited
				Completed:        true,
			}

			// Associate with the movie item
			historyRecord.Associate(existingMovie)

			// Save to database
			err = j.historyRepo.CreateHistory(ctx, &historyRecord)
			if err != nil {
				log.Warn().
					Err(err).
					Str("movieTitle", existingMovie.Title).
					Msg("Failed to create movie favorite status in history")
				continue
			}
		}

		// Also update the movie's favorite status in the MediaItem
		existingMovie.Data.Details.IsFavorite = true
		_, err = j.movieRepo.Update(ctx, *existingMovie)
		if err != nil {
			log.Warn().
				Err(err).
				Str("movieTitle", existingMovie.Title).
				Msg("Failed to update movie favorite status")
			continue
		}

		log.Debug().
			Str("movieTitle", existingMovie.Title).
			Msg("Successfully synced movie favorite status")

		processed++
	}

	return processed, nil
}

// syncSeriesFavorites syncs favorite series
func (j *FavoritesSyncJob) syncSeriesFavorites(ctx context.Context, userID, clientID uint64, mediaClient media.MediaClient, jobRunID uint64) (int, error) {
	log := utils.LoggerFromContext(ctx)

	// Check if client supports series
	seriesProvider, ok := mediaClient.(providers.SeriesProvider)
	if !ok {
		log.Debug().Msg("Client doesn't support series, skipping series favorites")
		return 0, nil
	}

	// Get favorite series
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 50, "Fetching favorite series")

	// Get series from client with favorites filter
	series, err := seriesProvider.GetSeries(ctx, &mediatypes.QueryOptions{
		Limit:     100,
		Favorites: true,
	})

	if err != nil {
		return 0, fmt.Errorf("failed to get favorite series: %w", err)
	}

	// Update progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 60,
		fmt.Sprintf("Processing %d favorite series", len(series)))

	// Track count of processed items
	processed := 0

	// Process each favorite series
	for _, s := range series {
		// Get client item ID for this client
		var clientItemID string
		for _, cid := range s.SyncClients {
			if cid.ID == clientID {
				clientItemID = cid.ItemID
				break
			}
		}

		if clientItemID == "" {
			log.Warn().
				Str("seriesTitle", s.Data.Details.Title).
				Msg("No matching client item ID found for series")
			continue
		}

		// Get the series from database
		existingSeries, err := j.seriesRepo.GetByClientItemID(ctx, clientItemID, clientID)
		if err != nil {
			log.Warn().
				Err(err).
				Str("clientItemID", clientItemID).
				Str("seriesTitle", s.Data.Details.Title).
				Msg("Series not found in database")
			continue
		}

		// Create or update history record to mark as favorite
		// First, check if there's already a history record
		hasViewed, err := j.historyRepo.HasUserViewedMedia(ctx, userID, existingSeries.ID)

		// Create or update history
		if hasViewed {
			// Get existing play count
			playCount, _ := j.historyRepo.GetItemPlayCount(ctx, userID, existingSeries.ID)

			// Create new history record (timestamp is needed since we don't have a way to update existing record)
			historyRecord := models.MediaPlayHistory[*mediatypes.Series]{
				MediaItemID:      existingSeries.ID,
				Type:             mediatypes.MediaTypeSeries,
				PlayedAt:         time.Now(),
				LastPlayedAt:     time.Now(),
				IsFavorite:       true,
				PlayCount:        int32(playCount),
				PlayedPercentage: 100, // Assume watched if favorited
				Completed:        true,
			}

			// Associate with the series item
			historyRecord.Associate(existingSeries)

			// Save to database
			err = j.historyRepo.CreateHistory(ctx, &historyRecord)
			if err != nil {
				log.Warn().
					Err(err).
					Str("seriesTitle", existingSeries.Title).
					Msg("Failed to update series favorite status in history")
				continue
			}
		} else {
			// Create new history record for a favorited but never watched item
			historyRecord := models.MediaPlayHistory[*mediatypes.Series]{
				MediaItemID:      existingSeries.ID,
				Type:             mediatypes.MediaTypeSeries,
				PlayedAt:         time.Now(),
				LastPlayedAt:     time.Now(),
				IsFavorite:       true,
				PlayCount:        1,   // First view
				PlayedPercentage: 100, // Assume watched if favorited
				Completed:        true,
			}

			// Associate with the series item
			historyRecord.Associate(existingSeries)

			// Save to database
			err = j.historyRepo.CreateHistory(ctx, &historyRecord)
			if err != nil {
				log.Warn().
					Err(err).
					Str("seriesTitle", existingSeries.Title).
					Msg("Failed to create series favorite status in history")
				continue
			}
		}

		// Also update the series's favorite status in the MediaItem
		existingSeries.Data.Details.IsFavorite = true
		_, err = j.seriesRepo.Update(ctx, *existingSeries)
		if err != nil {
			log.Warn().
				Err(err).
				Str("seriesTitle", existingSeries.Title).
				Msg("Failed to update series favorite status")
			continue
		}

		log.Debug().
			Str("seriesTitle", existingSeries.Title).
			Msg("Successfully synced series favorite status")

		processed++
	}

	return processed, nil
}

// syncEpisodeFavorites syncs favorite episodes
func (j *FavoritesSyncJob) syncEpisodeFavorites(ctx context.Context, userID, clientID uint64, mediaClient media.MediaClient, jobRunID uint64) (int, error) {
	// For simplicity, since episodes are usually favorited via series, we'll skip separate implementation
	// But this should follow similar patterns to movies and series

	return 0, nil
}

// syncMusicFavorites syncs favorite music tracks
func (j *FavoritesSyncJob) syncMusicFavorites(ctx context.Context, userID, clientID uint64, mediaClient media.MediaClient, jobRunID uint64) (int, error) {
	log := utils.LoggerFromContext(ctx)

	// Check if client supports music
	musicProvider, ok := mediaClient.(providers.MusicProvider)
	if !ok {
		log.Debug().Msg("Client doesn't support music, skipping music favorites")
		return 0, nil
	}

	// Get favorite tracks
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 70, "Fetching favorite music")

	// Get tracks from client with favorites filter
	tracks, err := musicProvider.GetMusic(ctx, &mediatypes.QueryOptions{
		Limit:     100,
		Favorites: true,
	})

	if err != nil {
		return 0, fmt.Errorf("failed to get favorite tracks: %w", err)
	}

	// Update progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 80,
		fmt.Sprintf("Processing %d favorite tracks", len(tracks)))

	// Track count of processed items
	processed := 0

	// Process each favorite track
	for _, track := range tracks {
		// Get client item ID for this client
		var clientItemID string
		for _, cid := range track.SyncClients {
			if cid.ID == clientID {
				clientItemID = cid.ItemID
				break
			}
		}

		if clientItemID == "" {
			log.Warn().
				Str("trackTitle", track.Data.Details.Title).
				Msg("No matching client item ID found for track")
			continue
		}

		// Get the track from database
		existingTrack, err := j.musicRepo.GetByClientItemID(ctx, clientItemID, clientID)
		if err != nil {
			log.Warn().
				Err(err).
				Str("clientItemID", clientItemID).
				Str("trackTitle", track.Data.Details.Title).
				Msg("Track not found in database")
			continue
		}

		// Create or update history record to mark as favorite
		// First, check if there's already a history record
		hasPlayed, err := j.historyRepo.HasUserViewedMedia(ctx, userID, existingTrack.ID)

		// Create or update history
		if hasPlayed {
			// Get existing play count
			playCount, _ := j.historyRepo.GetItemPlayCount(ctx, userID, existingTrack.ID)

			// Create new history record (timestamp is needed since we don't have a way to update existing record)
			historyRecord := models.MediaPlayHistory[*mediatypes.Track]{
				MediaItemID:      existingTrack.ID,
				Type:             mediatypes.MediaTypeTrack,
				PlayedAt:         time.Now(), // For music, this is "played at"
				LastPlayedAt:     time.Now(),
				IsFavorite:       true,
				PlayCount:        int32(playCount),
				PlayedPercentage: 100, // Assume listened to completion if favorited
				Completed:        true,
			}

			// Associate with the track item
			historyRecord.Associate(existingTrack)

			// Save to database
			err = j.historyRepo.CreateHistory(ctx, &historyRecord)
			if err != nil {
				log.Warn().
					Err(err).
					Str("trackTitle", existingTrack.Title).
					Msg("Failed to update track favorite status in history")
				continue
			}
		} else {
			// Create new history record for a favorited but never played item
			historyRecord := models.MediaPlayHistory[*mediatypes.Track]{
				MediaItemID:      existingTrack.ID,
				Type:             mediatypes.MediaTypeTrack,
				PlayedAt:         time.Now(),
				LastPlayedAt:     time.Now(),
				IsFavorite:       true,
				PlayCount:        1,   // First play
				PlayedPercentage: 100, // Assume listened to completion if favorited
				Completed:        true,
			}

			// Associate with the track item
			historyRecord.Associate(existingTrack)

			// Save to database
			err = j.historyRepo.CreateHistory(ctx, &historyRecord)
			if err != nil {
				log.Warn().
					Err(err).
					Str("trackTitle", existingTrack.Title).
					Msg("Failed to create track favorite status in history")
				continue
			}
		}

		// Also update the track's favorite status in the MediaItem
		existingTrack.Data.Details.IsFavorite = true
		_, err = j.musicRepo.Update(ctx, *existingTrack)
		if err != nil {
			log.Warn().
				Err(err).
				Str("trackTitle", existingTrack.Title).
				Msg("Failed to update track favorite status")
			continue
		}

		log.Debug().
			Str("trackTitle", existingTrack.Title).
			Msg("Successfully synced track favorite status")

		processed++
	}

	return processed, nil
}

// getMediaClient gets a media client from the client factory
func (j *FavoritesSyncJob) getMediaClient(ctx context.Context, clientID uint64, clientType string) (media.MediaClient, error) {
	log := utils.LoggerFromContext(ctx)

	// Get the client config from the repository
	var clientConfig clienttypes.ClientConfig

	mediaCategoryClients := j.clientRepos.GetAllByCategory(ctx, clienttypes.ClientCategoryMedia)

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
