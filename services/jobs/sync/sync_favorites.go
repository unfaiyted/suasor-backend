package sync

import (
	"context"
	"fmt"
	"log"
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

// FavoritesSyncJob synchronizes favorite/liked media from external clients
type FavoritesSyncJob struct {
	jobRepo         repository.JobRepository
	userRepo        repository.UserRepository
	userConfigRepo  repository.UserConfigRepository
	clientRepos     repobundles.ClientRepositories
	dataRepos       repobundles.UserMediaDataRepositories
	clientItemRepos repobundles.ClientMediaItemRepositories
	itemRepos       repobundles.UserMediaItemRepositories
	clientFactories *clients.ClientProviderFactoryService
}

// NewFavoritesSyncJob creates a new favorites sync job
func NewFavoritesSyncJob(
	jobRepo repository.JobRepository,
	userRepo repository.UserRepository,
	userConfigRepo repository.UserConfigRepository,
	clientRepos repobundles.ClientRepositories,
	dataRepos repobundles.UserMediaDataRepositories,
	clientItemRepos repobundles.ClientMediaItemRepositories,
	itemRepos repobundles.UserMediaItemRepositories,
	clientFactories *clients.ClientProviderFactoryService,
) *FavoritesSyncJob {
	return &FavoritesSyncJob{
		jobRepo:         jobRepo,
		userRepo:        userRepo,
		userConfigRepo:  userConfigRepo,
		clientRepos:     clientRepos,
		dataRepos:       dataRepos,
		clientItemRepos: clientItemRepos,
		itemRepos:       itemRepos,
		clientFactories: clientFactories,
	}
}

// Name returns the unique name of the job
func (j *FavoritesSyncJob) Name() string {
	// Make sure we always return a valid name even if struct is empty
	if j == nil || j.jobRepo == nil {
		return "system.favorites.sync"
	}
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

	// Check if job is properly initialized
	if j == nil || j.userRepo == nil || j.jobRepo == nil {
		log.Println("FavoritesSyncJob not properly initialized, using stub implementation")
		log.Println("Favorites sync job completed (no-op)")
		return nil
	}

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
	config, err := j.userConfigRepo.GetUserConfig(ctx, user.ID)
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
	clients, err := j.getUserClientMedias(ctx, user.ID)
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

// getUserClientMedias returns all media clients for a user
func (j *FavoritesSyncJob) getUserClientMedias(ctx context.Context, userID uint64) ([]ClientMediaInfo, error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().Uint64("userID", userID).Msg("Getting media clients for user")

	var clients []ClientMediaInfo

	// Get all media client types from the repository collection

	// Emby clients
	embyClients, err := j.clientRepos.EmbyRepo().GetByUserID(ctx, userID)
	if err != nil {
		log.Warn().Err(err).Msg("Error getting Emby clients")
	} else {
		for _, c := range embyClients {
			clients = append(clients, ClientMediaInfo{
				ClientID:   c.ID,
				ClientType: clienttypes.ClientMediaTypeEmby,
				Name:       c.Name,
				UserID:     userID,
			})
		}
	}

	// Jellyfin clients
	jellyfinClients, err := j.clientRepos.JellyfinRepo().GetByUserID(ctx, userID)
	if err != nil {
		log.Warn().Err(err).Msg("Error getting Jellyfin clients")
	} else {
		for _, c := range jellyfinClients {
			clients = append(clients, ClientMediaInfo{
				ClientID:   c.ID,
				ClientType: clienttypes.ClientMediaTypeJellyfin,
				Name:       c.Name,
				UserID:     userID,
			})
		}
	}

	// Plex clients
	plexClients, err := j.clientRepos.PlexRepo().GetByUserID(ctx, userID)
	if err != nil {
		log.Warn().Err(err).Msg("Error getting Plex clients")
	} else {
		for _, c := range plexClients {
			clients = append(clients, ClientMediaInfo{
				ClientID:   c.ID,
				ClientType: clienttypes.ClientMediaTypePlex,
				Name:       c.Name,
				UserID:     userID,
			})
		}
	}

	// Subsonic clients (primarily for music)
	subsonicClients, err := j.clientRepos.SubsonicRepo().GetByUserID(ctx, userID)
	if err != nil {
		log.Warn().Err(err).Msg("Error getting Subsonic clients")
	} else {
		for _, c := range subsonicClients {
			clients = append(clients, ClientMediaInfo{
				ClientID:   c.ID,
				ClientType: clienttypes.ClientMediaTypeSubsonic,
				Name:       c.Name,
				UserID:     userID,
			})
		}
	}

	log.Info().Int("clientCount", len(clients)).Msg("Found media clients")
	return clients, nil
}

// syncClientFavorites syncs favorites for a specific client
func (j *FavoritesSyncJob) syncClientFavorites(ctx context.Context, userID uint64, client ClientMediaInfo, jobRunID uint64) error {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", client.ClientID).
		Str("clientName", client.Name).
		Str("clientType", string(client.ClientType)).
		Msg("Syncing favorites")

	// Get the client using client factory
	clientMedia, err := j.getClientMedia(ctx, client.ClientID, string(client.ClientType))
	if err != nil {
		return fmt.Errorf("failed to get media client: %w", err)
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 20, "Fetching favorites from client")

	// Process different media types
	totalProcessed := 0

	// Sync movie favorites
	movieCount, err := j.syncMovieFavorites(ctx, userID, client.ClientID, clientMedia, jobRunID)
	if err != nil {
		log.Warn().Err(err).Msg("Error syncing movie favorites")
	}
	totalProcessed += movieCount

	// Sync series favorites
	seriesCount, err := j.syncSeriesFavorites(ctx, userID, client.ClientID, clientMedia, jobRunID)
	if err != nil {
		log.Warn().Err(err).Msg("Error syncing series favorites")
	}
	totalProcessed += seriesCount

	// Sync episode favorites
	episodeCount, err := j.syncEpisodeFavorites(ctx, userID, client.ClientID, clientMedia, jobRunID)
	if err != nil {
		log.Warn().Err(err).Msg("Error syncing episode favorites")
	}
	totalProcessed += episodeCount

	// Sync music favorites
	musicCount, err := j.syncMusicFavorites(ctx, userID, client.ClientID, clientMedia, jobRunID)
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
func (j *FavoritesSyncJob) syncMovieFavorites(ctx context.Context, userID, clientID uint64, clientMedia media.ClientMedia, jobRunID uint64) (int, error) {
	log := logger.LoggerFromContext(ctx)

	// Check if client supports movies
	movieProvider, ok := clientMedia.(providers.MovieProvider)
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
		existingMovie, err := j.itemRepos.MovieUserRepo().GetByClientItemID(ctx, clientID, clientItemID)
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
		hasViewed, err := j.dataRepos.MovieDataRepo().HasUserViewedMedia(ctx, userID, existingMovie.ID)

		// Create or update history
		if hasViewed {
			// Get existing history (we'll need to implement a method to get history by user and media item)
			playCount, _ := j.dataRepos.MovieDataRepo().GetItemPlayCount(ctx, userID, existingMovie.ID)

			// Create new history record (timestamp is needed since we don't have a way to update existing record)
			historyRecord := models.UserMediaItemData[*mediatypes.Movie]{
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
			item, err := j.dataRepos.MovieDataRepo().Create(ctx, &historyRecord)
			if err != nil {
				log.Warn().
					Err(err).
					Str("movieTitle", existingMovie.Title).
					Msg("Failed to update movie favorite status in history")
				continue
			}
			log.Debug().
				Str("movieTitle", existingMovie.Title).
				Time("watchedAt", item.PlayedAt).
				Float64("percentage", item.PlayedPercentage).
				Msg("Saved movie favorite status in history")
		} else {
			// Create new history record for a favorited but never watched item
			historyRecord := models.UserMediaItemData[*mediatypes.Movie]{
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
			item, err := j.dataRepos.MovieDataRepo().Create(ctx, &historyRecord)
			if err != nil {
				log.Warn().
					Err(err).
					Str("movieTitle", existingMovie.Title).
					Msg("Failed to create movie favorite status in history")
				continue
			}
			log.Debug().
				Str("movieTitle", existingMovie.Title).
				Time("watchedAt", item.PlayedAt).
				Float64("percentage", item.PlayedPercentage).
				Msg("Saved movie favorite status in history")
		}

		// Also update the movie's favorite status in the MediaItem
		existingMovie.Data.Details.IsFavorite = true
		_, err = j.itemRepos.MovieUserRepo().Update(ctx, existingMovie)
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
func (j *FavoritesSyncJob) syncSeriesFavorites(ctx context.Context, userID, clientID uint64, clientMedia media.ClientMedia, jobRunID uint64) (int, error) {
	log := logger.LoggerFromContext(ctx)

	// Check if client supports series
	seriesProvider, ok := clientMedia.(providers.SeriesProvider)
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
		existingSeries, err := j.itemRepos.SeriesUserRepo().GetByClientItemID(ctx, clientID, clientItemID)
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
		hasViewed, err := j.dataRepos.SeriesDataRepo().HasUserViewedMedia(ctx, userID, existingSeries.ID)

		// Create or update history
		if hasViewed {
			// Get existing play count
			playCount, _ := j.dataRepos.SeriesDataRepo().GetItemPlayCount(ctx, userID, existingSeries.ID)

			// Create new history record (timestamp is needed since we don't have a way to update existing record)
			historyRecord := models.UserMediaItemData[*mediatypes.Series]{
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
			item, err := j.dataRepos.SeriesDataRepo().Create(ctx, &historyRecord)
			if err != nil {
				log.Warn().
					Err(err).
					Str("seriesTitle", existingSeries.Title).
					Msg("Failed to update series favorite status in history")
				continue
			}

			log.Debug().
				Str("seriesTitle", existingSeries.Title).
				Time("watchedAt", item.PlayedAt).
				Float64("percentage", item.PlayedPercentage).
				Msg("Saved series favorite status in history")
		} else {
			// Create new history record for a favorited but never watched item
			historyRecord := models.UserMediaItemData[*mediatypes.Series]{
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
			item, err := j.dataRepos.SeriesDataRepo().Create(ctx, &historyRecord)
			if err != nil {
				log.Warn().
					Err(err).
					Str("seriesTitle", existingSeries.Title).
					Msg("Failed to create series favorite status in history")
				continue
			}
			log.Debug().
				Str("seriesTitle", existingSeries.Title).
				Time("watchedAt", item.PlayedAt).
				Float64("percentage", item.PlayedPercentage).
				Msg("Saved series favorite status in history")
		}

		// Also update the series's favorite status in the MediaItem
		existingSeries.Data.Details.IsFavorite = true
		_, err = j.itemRepos.SeriesUserRepo().Update(ctx, existingSeries)
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
func (j *FavoritesSyncJob) syncEpisodeFavorites(ctx context.Context, userID, clientID uint64, clientMedia media.ClientMedia, jobRunID uint64) (int, error) {
	// For simplicity, since episodes are usually favorited via series, we'll skip separate implementation
	// But this should follow similar patterns to movies and series

	return 0, nil
}

// syncMusicFavorites syncs favorite music tracks
func (j *FavoritesSyncJob) syncMusicFavorites(ctx context.Context, userID, clientID uint64, clientMedia media.ClientMedia, jobRunID uint64) (int, error) {
	log := logger.LoggerFromContext(ctx)

	// Check if client supports music
	musicProvider, ok := clientMedia.(providers.MusicProvider)
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
		existingTrack, err := j.itemRepos.TrackUserRepo().GetByClientItemID(ctx, clientID, clientItemID)
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
		hasPlayed, err := j.dataRepos.TrackDataRepo().HasUserViewedMedia(ctx, userID, existingTrack.ID)

		// Create or update history
		if hasPlayed {
			// Get existing play count
			playCount, _ := j.dataRepos.TrackDataRepo().GetItemPlayCount(ctx, userID, existingTrack.ID)

			// Create new history record (timestamp is needed since we don't have a way to update existing record)
			historyRecord := models.UserMediaItemData[*mediatypes.Track]{
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
			item, err := j.dataRepos.TrackDataRepo().Create(ctx, &historyRecord)
			if err != nil {
				log.Warn().
					Err(err).
					Str("trackTitle", existingTrack.Title).
					Msg("Failed to update track favorite status in history")
				continue
			}
			log.Debug().
				Str("trackTitle", existingTrack.Title).
				Time("playedAt", item.PlayedAt).
				Float64("percentage", item.PlayedPercentage).
				Msg("Saved track favorite status in history")
		} else {
			// Create new history record for a favorited but never played item
			historyRecord := models.UserMediaItemData[*mediatypes.Track]{
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
			item, err := j.dataRepos.TrackDataRepo().Create(ctx, &historyRecord)
			if err != nil {
				log.Warn().
					Err(err).
					Str("trackTitle", existingTrack.Title).
					Msg("Failed to create track favorite status in history")
				continue
			}
			log.Debug().
				Str("trackTitle", existingTrack.Title).
				Time("playedAt", item.PlayedAt).
				Float64("percentage", item.PlayedPercentage).
				Msg("Saved track favorite status in history")
		}

		// Also update the track's favorite status in the MediaItem
		existingTrack.Data.Details.IsFavorite = true
		_, err = j.itemRepos.TrackUserRepo().Update(ctx, existingTrack)
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

// getClientMedia gets a media client from the client factory
func (j *FavoritesSyncJob) getClientMedia(ctx context.Context, clientID uint64, clientType string) (media.ClientMedia, error) {
	log := logger.LoggerFromContext(ctx)

	// Get the client config from the repository
	var clientConfig clienttypes.ClientConfig

	switch clienttypes.ClientMediaType(clientType) {
	case clienttypes.ClientMediaTypeEmby:
		c, err := j.clientRepos.EmbyRepo().GetByID(ctx, clientID)
		if err != nil {
			return nil, fmt.Errorf("failed to get emby client: %w", err)
		}
		clientConfig = c.GetConfig()
	case clienttypes.ClientMediaTypeJellyfin:
		c, err := j.clientRepos.JellyfinRepo().GetByID(ctx, clientID)
		if err != nil {
			return nil, fmt.Errorf("failed to get jellyfin client: %w", err)
		}
		clientConfig = c.GetConfig()
	case clienttypes.ClientMediaTypePlex:
		c, err := j.clientRepos.PlexRepo().GetByID(ctx, clientID)
		if err != nil {
			return nil, fmt.Errorf("failed to get plex client: %w", err)
		}
		clientConfig = c.GetConfig()
	case clienttypes.ClientMediaTypeSubsonic:
		c, err := j.clientRepos.SubsonicRepo().GetByID(ctx, clientID)
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
	clientMedia, ok := clientInstance.(media.ClientMedia)
	if !ok {
		return nil, fmt.Errorf("client is not a media client")
	}

	log.Info().
		Uint64("clientID", clientID).
		Str("clientType", clientType).
		Msg("Successfully retrieved and initialized media client")

	return clientMedia, nil
}
