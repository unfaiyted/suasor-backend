package jobs

import (
	"context"
	"fmt"
	"log"
	"strings"
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

// MediaSyncJob handles syncing of media items from clients
type MediaSyncJob struct {
	jobRepo         repository.JobRepository
	userRepo        repository.UserRepository
	configRepo      repository.UserConfigRepository
	movieRepo       repository.MediaItemRepository[*mediatypes.Movie]
	seriesRepo      repository.MediaItemRepository[*mediatypes.Series]
	episodeRepo     repository.MediaItemRepository[*mediatypes.Episode]
	musicRepo       repository.MediaItemRepository[*mediatypes.Track]
	albumRepo       repository.MediaItemRepository[*mediatypes.Album]
	artistRepo      repository.MediaItemRepository[*mediatypes.Artist]
	clientRepos     repository.ClientRepositoryCollection
	clientFactories *client.ClientFactoryService
}

// NewMediaSyncJob creates a new media sync job
func NewMediaSyncJob(
	jobRepo repository.JobRepository,
	userRepo repository.UserRepository,
	configRepo repository.UserConfigRepository,
	movieRepo repository.MediaItemRepository[*mediatypes.Movie],
	seriesRepo repository.MediaItemRepository[*mediatypes.Series],
	episodeRepo repository.MediaItemRepository[*mediatypes.Episode],
	musicRepo repository.MediaItemRepository[*mediatypes.Track],
	albumRepo repository.MediaItemRepository[*mediatypes.Album],
	artistRepo repository.MediaItemRepository[*mediatypes.Artist],
	clientRepos repository.ClientRepositoryCollection,
	clientFactories *client.ClientFactoryService,
) *MediaSyncJob {
	return &MediaSyncJob{
		jobRepo:         jobRepo,
		userRepo:        userRepo,
		configRepo:      configRepo,
		movieRepo:       movieRepo,
		seriesRepo:      seriesRepo,
		episodeRepo:     episodeRepo,
		musicRepo:       musicRepo,
		albumRepo:       albumRepo,
		artistRepo:      artistRepo,
		clientRepos:     clientRepos,
		clientFactories: clientFactories,
	}
}

// Schedule returns how often the job should run (default)
func (j *MediaSyncJob) Schedule() time.Duration {
	// Default to checking hourly
	return 1 * time.Hour
}

// Name returns the unique name of the job
func (j *MediaSyncJob) Name() string {
	return "system.media.sync"
}

// Execute runs the job
func (j *MediaSyncJob) Execute(ctx context.Context) error {
	log.Println("Starting media sync job")

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
func (j *MediaSyncJob) RunManualSync(ctx context.Context, userID uint64, clientID uint64, mediaType string) error {
	// First, we need to determine the client type
	// We'll need to get it from the client record in the database

	log := utils.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("mediaType", mediaType).
		Msg("Running manual media sync job")

	// Get a list of possible client types to check
	clientTypes := []string{"emby", "jellyfin", "plex", "subsonic"}

	// Try to find which type this client is
	var clientType string
	for _, cType := range clientTypes {
		// Try to get the config for this client type
		_, err := j.getClientConfig(ctx, clientID, cType)
		if err == nil {
			// Found the client type
			clientType = cType
			break
		}
	}

	// If we couldn't determine the client type, return an error
	if clientType == "" {
		return fmt.Errorf("couldn't determine client type for clientID=%d", clientID)
	}

	log.Info().
		Str("clientType", clientType).
		Msg("Determined client type for manual sync")

	// Create a temporary sync job
	syncJob := models.MediaSyncJob{
		UserID:     userID,
		ClientID:   clientID,
		ClientType: clientType,
		MediaType:  mediaType,
	}

	// Run the sync job
	return j.runSyncJob(ctx, syncJob)
}

// SyncUserMediaFromClient runs a sync job for a specific user and client
// This is an alias for RunManualSync for backward compatibility
func (j *MediaSyncJob) SyncUserMediaFromClient(ctx context.Context, userID uint64, clientID uint64, mediaType string) error {
	return j.RunManualSync(ctx, userID, clientID, mediaType)
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
		JobName:   fmt.Sprintf("%s.%s", j.Name(), syncJob.MediaType),
		JobType:   models.JobTypeSync,
		Status:    models.JobStatusRunning,
		StartTime: &now,
		UserID:    &syncJob.UserID,
		Metadata:  fmt.Sprintf(`{"clientId":%d,"mediaType":"%s"}`, syncJob.ClientID, syncJob.MediaType),
	}

	// Save the job run
	err := j.jobRepo.CreateJobRun(ctx, jobRun)
	if err != nil {
		return fmt.Errorf("failed to create job run: %w", err)
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRun.ID, 0, "Starting media sync")

	// Get the client from the database
	mediaClient, err := j.getMediaClient(ctx, syncJob.ClientID, syncJob.ClientType)
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to get media client: %v", err)
		j.completeJobRun(ctx, jobRun.ID, models.JobStatusFailed, errorMsg)
		return fmt.Errorf(errorMsg)
	}

	// Process different media types
	var syncError error
	mediaType := strings.ToLower(syncJob.MediaType)

	// Normalize media type to handle both singular and plural forms
	switch mediaType {
	case "movie", "movies":
		syncError = j.syncMovies(ctx, mediaClient, jobRun.ID, syncJob.ClientID)
	case "series", "serie", "tvshows", "tvshow", "tv", "shows", "show":
		syncError = j.syncSeries(ctx, mediaClient, jobRun.ID, syncJob.ClientID)
	case "episode", "episodes":
		syncError = j.syncEpisodes(ctx, mediaClient, jobRun.ID, syncJob.ClientID)
	case "music", "tracks", "track", "songs", "song":
		syncError = j.syncMusic(ctx, mediaClient, jobRun.ID, syncJob.ClientID)
	case "artist", "artists":
		syncError = j.syncArtists(ctx, mediaClient, jobRun.ID, syncJob.ClientID)
	case "album", "albums":
		syncError = j.syncAlbums(ctx, mediaClient, jobRun.ID, syncJob.ClientID)
	default:
		syncError = fmt.Errorf("unsupported media type: %s", syncJob.MediaType)
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

// completeJobRun marks a job run as completed with the given status and error message
func (j *MediaSyncJob) completeJobRun(ctx context.Context, jobRunID uint64, status models.JobStatus, errorMsg string) {
	if err := j.jobRepo.CompleteJobRun(ctx, jobRunID, status, errorMsg); err != nil {
		log.Printf("Error completing job run: %v", err)
	}
}

func (j *MediaSyncJob) getClientConfig(ctx context.Context, clientID uint64, clientType string) (clienttypes.ClientConfig, error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", clientID).
		Msg("Retrieving client config from database")

	// Get client config from database
	var config clienttypes.ClientConfig

	clientRepos := j.clientRepos.GetAllByCategory(clienttypes.ClientCategoryMedia)
	switch clientType {
	case "emby":
		c, err := clientRepos.EmbyRepo.GetByID(ctx, clientID)
		if err != nil {
			return nil, fmt.Errorf("failed to get emby client: %w", err)
		}
		config = c.GetConfig()
	case "jellyfin":
		c, err := clientRepos.JellyfinRepo.GetByID(ctx, clientID)
		if err != nil {
			return nil, fmt.Errorf("failed to get jellyfin client: %w", err)
		}
		config = c.GetConfig()
	case "plex":
		c, err := clientRepos.PlexRepo.GetByID(ctx, clientID)
		if err != nil {
			return nil, fmt.Errorf("failed to get plex client: %w", err)
		}
		config = c.GetConfig()
	case "subsonic":
		c, err := clientRepos.SubsonicRepo.GetByID(ctx, clientID)
		if err != nil {
			return nil, fmt.Errorf("failed to get subsonic client: %w", err)
		}
		config = c.GetConfig()
	default:
		return nil, fmt.Errorf("unsupported client type: %s", clientType)
	}

	return config, nil
}

// getMediaClient gets a media client from the database and initializes it
func (j *MediaSyncJob) getMediaClient(ctx context.Context, clientID uint64, clientType string) (media.MediaClient, error) {
	// Use the type to get the client config by id
	clientConfig, err := j.getClientConfig(ctx, clientID, clientType)
	if err != nil {
		return nil, fmt.Errorf("failed to get client config: %w", err)
	}

	// Validate client config is not nil before proceeding
	if clientConfig == nil {
		return nil, fmt.Errorf("client config is nil for clientID=%d, clientType=%s", clientID, clientType)
	}

	// Cast media client from generic client
	client, err := j.clientFactories.GetClient(ctx, clientID, clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	// Cast to media client
	mediaClient, ok := client.(media.MediaClient)
	if !ok {
		return nil, fmt.Errorf("client is not a media client")
	}

	return mediaClient, nil
}

// syncMovies syncs movies from the client to the database
func (j *MediaSyncJob) syncMovies(ctx context.Context, mediaClient media.MediaClient, jobRunID uint64, clientID uint64) error {
	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 10, "Fetching movies from client")

	// Check if client supports movies
	movieProvider, ok := mediaClient.(providers.MovieProvider)
	if !ok {
		return fmt.Errorf("client doesn't support movies")
	}

	// Get all movies from the client
	clientType := mediaClient.(client.Client).GetType().AsMediaClientType()
	movies, err := movieProvider.GetMovies(ctx, &mediatypes.QueryOptions{})
	if err != nil {
		return fmt.Errorf("failed to get movies: %w", err)
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 50, fmt.Sprintf("Processing %d movies", len(movies)))

	// Process movies in batches to avoid memory issues
	batchSize := 50
	totalMovies := len(movies)
	processedMovies := 0

	for i := 0; i < totalMovies; i += batchSize {
		end := i + batchSize
		if end > totalMovies {
			end = totalMovies
		}

		movieBatch := movies[i:end]
		err := j.processMovieBatch(ctx, movieBatch, clientID, clientType)
		if err != nil {
			return fmt.Errorf("failed to process movie batch: %w", err)
		}

		processedMovies += len(movieBatch)
		progress := 50 + int(float64(processedMovies)/float64(totalMovies)*50.0)
		j.jobRepo.UpdateJobProgress(ctx, jobRunID, progress, fmt.Sprintf("Processed %d/%d movies", processedMovies, totalMovies))
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 100, fmt.Sprintf("Synced %d movies", totalMovies))

	return nil
}

// syncSeries syncs TV series from the client to the database
func (j *MediaSyncJob) syncSeries(ctx context.Context, mediaClient media.MediaClient, jobRunID uint64, clientID uint64) error {
	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 10, "Fetching series from client")

	// Check if client supports series
	seriesProvider, ok := mediaClient.(providers.SeriesProvider)
	if !ok {
		return fmt.Errorf("client doesn't support series")
	}

	// Get all series from the client
	clientType := mediaClient.(client.Client).GetType().AsMediaClientType()
	series, err := seriesProvider.GetSeries(ctx, &mediatypes.QueryOptions{})
	if err != nil {
		return fmt.Errorf("failed to get series: %w", err)
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 50, fmt.Sprintf("Processing %d series", len(series)))

	// Process series in batches to avoid memory issues
	batchSize := 50
	totalSeries := len(series)
	processedSeries := 0

	for i := 0; i < totalSeries; i += batchSize {
		end := i + batchSize
		if end > totalSeries {
			end = totalSeries
		}

		seriesBatch := series[i:end]
		err := j.processSeriesBatch(ctx, seriesBatch, clientID, clientType)
		if err != nil {
			return fmt.Errorf("failed to process series batch: %w", err)
		}

		processedSeries += len(seriesBatch)
		progress := 50 + int(float64(processedSeries)/float64(totalSeries)*50.0)
		j.jobRepo.UpdateJobProgress(ctx, jobRunID, progress, fmt.Sprintf("Processed %d/%d series", processedSeries, totalSeries))
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 100, fmt.Sprintf("Synced %d series", totalSeries))

	return nil
}

// syncEpisodes syncs TV episodes from the client to the database
func (j *MediaSyncJob) syncEpisodes(ctx context.Context, mediaClient media.MediaClient, jobRunID uint64, clientID uint64) error {
	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 10, "Fetching episodes from client")

	// Check if client supports episodes
	seriesProvider, ok := mediaClient.(providers.SeriesProvider)
	if !ok {
		return fmt.Errorf("client doesn't support episodes")
	}

	// Get all episodes from the client
	clientType := mediaClient.(client.Client).GetType().AsMediaClientType()
	
	// Initialize a slice to hold all episodes
	var allEpisodes []models.MediaItem[*mediatypes.Episode]
	
	// First get all series
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 15, "Fetching series list")
	allSeries, err := seriesProvider.GetSeries(ctx, &mediatypes.QueryOptions{})
	if err != nil {
		return fmt.Errorf("failed to get series list: %w", err)
	}
	
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 20, fmt.Sprintf("Found %d series, fetching episodes", len(allSeries)))
	
	// Set total items for tracking progress
	totalSeries := len(allSeries)
	j.jobRepo.SetJobTotalItems(ctx, jobRunID, totalSeries)
	processedSeries := 0
	
	// For each series, get episodes
	for _, series := range allSeries {
		if series.Data == nil || len(series.ClientIDs) == 0 {
			// Skip series with no data or no client ID
			log.Printf("Skipping series with missing data")
			continue
		}
		
		// Find the client item ID for this series
		var seriesID string
		for _, cid := range series.ClientIDs {
			if cid.ID == clientID {
				seriesID = cid.ItemID
				break
			}
		}
		
		if seriesID == "" {
			log.Printf("No matching client item ID found for series: %s", series.Data.Details.Title)
			continue
		}
		
		// Get seasons for this series
		seasons, err := seriesProvider.GetSeriesSeasons(ctx, seriesID)
		if err != nil {
			log.Printf("Error getting seasons for series %s: %v", series.Data.Details.Title, err)
			continue
		}
		
		// For each season, get episodes
		for _, season := range seasons {
			if season.Data == nil {
				continue
			}
			
			seasonNumber := season.Data.Number
			episodes, err := seriesProvider.GetSeriesEpisodes(ctx, seriesID, seasonNumber)
			if err != nil {
				log.Printf("Error getting episodes for series %s season %d: %v", 
					series.Data.Details.Title, seasonNumber, err)
				continue
			}
			
			allEpisodes = append(allEpisodes, episodes...)
		}
		
		// Update progress
		processedSeries++
		progress := 20 + int(float64(processedSeries)/float64(totalSeries)*30.0)
		j.jobRepo.UpdateJobProgress(ctx, jobRunID, progress, 
			fmt.Sprintf("Processed %d/%d series, found %d episodes", 
				processedSeries, totalSeries, len(allEpisodes)))
		j.jobRepo.IncrementJobProcessedItems(ctx, jobRunID, 1)
	}
	
	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 50, 
		fmt.Sprintf("Processing %d episodes", len(allEpisodes)))
	
	// Process episodes in batches to avoid memory issues
	batchSize := 100
	totalEpisodes := len(allEpisodes)
	processedEpisodes := 0
	
	for i := 0; i < totalEpisodes; i += batchSize {
		end := i + batchSize
		if end > totalEpisodes {
			end = totalEpisodes
		}
		
		episodeBatch := allEpisodes[i:end]
		err := j.processEpisodeBatch(ctx, episodeBatch, clientID, clientType)
		if err != nil {
			return fmt.Errorf("failed to process episode batch: %w", err)
		}
		
		processedEpisodes += len(episodeBatch)
		progress := 50 + int(float64(processedEpisodes)/float64(totalEpisodes)*50.0)
		j.jobRepo.UpdateJobProgress(ctx, jobRunID, progress, 
			fmt.Sprintf("Processed %d/%d episodes", processedEpisodes, totalEpisodes))
	}
	
	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 100, 
		fmt.Sprintf("Synced %d episodes from %d series", totalEpisodes, totalSeries))

	return nil
}

// syncMusic syncs music tracks from the client to the database
func (j *MediaSyncJob) syncMusic(ctx context.Context, mediaClient media.MediaClient, jobRunID uint64, clientID uint64) error {
	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 10, "Fetching music from client")

	// Check if client supports music
	musicProvider, ok := mediaClient.(providers.MusicProvider)
	if !ok {
		return fmt.Errorf("client doesn't support music")
	}

	// Get all tracks from the client
	clientType := mediaClient.(client.Client).GetType().AsMediaClientType()
	tracks, err := musicProvider.GetMusic(ctx, &mediatypes.QueryOptions{})
	if err != nil {
		return fmt.Errorf("failed to get tracks: %w", err)
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 50, fmt.Sprintf("Processing %d tracks", len(tracks)))

	// Process tracks in batches to avoid memory issues
	batchSize := 100
	totalTracks := len(tracks)
	processedTracks := 0

	for i := 0; i < totalTracks; i += batchSize {
		end := i + batchSize
		if end > totalTracks {
			end = totalTracks
		}

		trackBatch := tracks[i:end]
		err := j.processTrackBatch(ctx, trackBatch, clientID, clientType)
		if err != nil {
			return fmt.Errorf("failed to process track batch: %w", err)
		}

		processedTracks += len(trackBatch)
		progress := 50 + int(float64(processedTracks)/float64(totalTracks)*50.0)
		j.jobRepo.UpdateJobProgress(ctx, jobRunID, progress, fmt.Sprintf("Processed %d/%d tracks", processedTracks, totalTracks))
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 100, fmt.Sprintf("Synced %d tracks", totalTracks))

	return nil
}

// syncAlbums syncs music albums from the client to the database
func (j *MediaSyncJob) syncAlbums(ctx context.Context, mediaClient media.MediaClient, jobRunID uint64, clientID uint64) error {
	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 10, "Fetching albums from client")

	// Check if client supports albums
	musicProvider, ok := mediaClient.(providers.MusicProvider)
	if !ok {
		return fmt.Errorf("client doesn't support albums")
	}

	// Get all albums from the client
	clientType := mediaClient.(client.Client).GetType().AsMediaClientType()
	albums, err := musicProvider.GetMusicAlbums(ctx, &mediatypes.QueryOptions{})
	if err != nil {
		return fmt.Errorf("failed to get albums: %w", err)
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 50, fmt.Sprintf("Processing %d albums", len(albums)))

	// Process albums in batches to avoid memory issues
	batchSize := 50
	totalAlbums := len(albums)
	processedAlbums := 0

	for i := 0; i < totalAlbums; i += batchSize {
		end := i + batchSize
		if end > totalAlbums {
			end = totalAlbums
		}

		albumBatch := albums[i:end]
		err := j.processAlbumBatch(ctx, albumBatch, clientID, clientType)
		if err != nil {
			return fmt.Errorf("failed to process album batch: %w", err)
		}

		processedAlbums += len(albumBatch)
		progress := 50 + int(float64(processedAlbums)/float64(totalAlbums)*50.0)
		j.jobRepo.UpdateJobProgress(ctx, jobRunID, progress, fmt.Sprintf("Processed %d/%d albums", processedAlbums, totalAlbums))
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 100, fmt.Sprintf("Synced %d albums", totalAlbums))

	return nil
}

// syncArtists syncs music artists from the client to the database
func (j *MediaSyncJob) syncArtists(ctx context.Context, mediaClient media.MediaClient, jobRunID uint64, clientID uint64) error {
	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 10, "Fetching artists from client")

	// Check if client supports artists
	musicProvider, ok := mediaClient.(providers.MusicProvider)
	if !ok {
		return fmt.Errorf("client doesn't support artists")
	}

	// Get all artists from the client
	clientType := mediaClient.(client.Client).GetType().AsMediaClientType()
	artists, err := musicProvider.GetMusicArtists(ctx, &mediatypes.QueryOptions{})
	if err != nil {
		return fmt.Errorf("failed to get artists: %w", err)
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 50, fmt.Sprintf("Processing %d artists", len(artists)))

	// Process artists in batches to avoid memory issues
	batchSize := 50
	totalArtists := len(artists)
	processedArtists := 0

	for i := 0; i < totalArtists; i += batchSize {
		end := i + batchSize
		if end > totalArtists {
			end = totalArtists
		}

		artistBatch := artists[i:end]
		err := j.processArtistBatch(ctx, artistBatch, clientID, clientType)
		if err != nil {
			return fmt.Errorf("failed to process artist batch: %w", err)
		}

		processedArtists += len(artistBatch)
		progress := 50 + int(float64(processedArtists)/float64(totalArtists)*50.0)
		j.jobRepo.UpdateJobProgress(ctx, jobRunID, progress, fmt.Sprintf("Processed %d/%d artists", processedArtists, totalArtists))
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 100, fmt.Sprintf("Synced %d artists", totalArtists))

	return nil
}

// processMovieBatch processes a batch of movies and saves them to the database
func (j *MediaSyncJob) processMovieBatch(ctx context.Context, movies []models.MediaItem[*mediatypes.Movie], clientID uint64, clientType clienttypes.MediaClientType) error {
	for _, movie := range movies {
		// Skip if movie has no client ID information
		if len(movie.ClientIDs) == 0 {
			log.Printf("Skipping movie with no client IDs: %s", movie.Data.Details.Title)
			continue
		}

		// Get the client ID and item ID for lookup
		clientItemID := ""
		for _, cid := range movie.ClientIDs {
			if cid.ID == clientID {
				clientItemID = cid.ItemID
				break
			}
		}

		if clientItemID == "" {
			log.Printf("No matching client item ID found for movie: %s", movie.Data.Details.Title)
			continue
		}

		// Check if the movie already exists in the database
		existingMovie, err := j.movieRepo.GetByClientItemID(ctx, clientItemID, clientID)
		if err == nil {
			// Movie exists, update it
			// Merge client IDs
			for _, cid := range movie.ClientIDs {
				found := false
				for i, existingCid := range existingMovie.ClientIDs {
					if existingCid.ID == cid.ID && existingCid.Type == cid.Type {
						// Update existing entry if needed
						existingMovie.ClientIDs[i].ItemID = cid.ItemID
						found = true
						break
					}
				}
				if !found {
					// Add new client ID
					existingMovie.ClientIDs = append(existingMovie.ClientIDs, cid)
				}
			}

			// Merge external IDs
			for _, extID := range movie.ExternalIDs {
				found := false
				for i, existingExtID := range existingMovie.ExternalIDs {
					if existingExtID.Source == extID.Source {
						// Update existing entry
						existingMovie.ExternalIDs[i].ID = extID.ID
						found = true
						break
					}
				}
				if !found {
					// Add new external ID
					existingMovie.ExternalIDs = append(existingMovie.ExternalIDs, extID)
				}
			}

			// Update data fields
			existingMovie.Data.Details = movie.Data.Details
			existingMovie.Title = movie.Data.Details.Title
			existingMovie.ReleaseDate = movie.Data.Details.ReleaseDate
			existingMovie.ReleaseYear = movie.Data.Details.ReleaseYear

			// Save the updated movie
			_, err = j.movieRepo.Update(ctx, *existingMovie)
			if err != nil {
				log.Printf("Error updating movie: %v", err)
				continue
			}
		} else {
			// Movie doesn't exist, create it
			// Set top level title and release fields
			movie.Title = movie.Data.Details.Title
			movie.ReleaseDate = movie.Data.Details.ReleaseDate
			movie.ReleaseYear = movie.Data.Details.ReleaseYear

			// Create the movie
			_, err = j.movieRepo.Create(ctx, movie)
			if err != nil {
				log.Printf("Error creating movie: %v", err)
				continue
			}
		}
	}

	return nil
}

// processSeriesBatch processes a batch of series and saves them to the database
func (j *MediaSyncJob) processSeriesBatch(ctx context.Context, series []models.MediaItem[*mediatypes.Series], clientID uint64, clientType clienttypes.MediaClientType) error {
	// Try to get a series provider for this client to fetch season details
	mediaClient, err := j.getMediaClient(ctx, clientID, clientType.String())
	if err != nil {
		// Just log the error but continue processing with what we have
		log.Printf("Failed to get media client for season details: %v", err)
	}
	
	// Cast to series provider if possible
	var seriesProvider providers.SeriesProvider
	if mediaClient != nil {
		if sp, ok := mediaClient.(providers.SeriesProvider); ok {
			seriesProvider = sp
		}
	}
	for _, s := range series {
		// Skip if series has no client ID information
		if len(s.ClientIDs) == 0 {
			log.Printf("Skipping series with no client IDs: %s", s.Data.Details.Title)
			continue
		}

		// Get the client ID and item ID for lookup
		clientItemID := ""
		for _, cid := range s.ClientIDs {
			if cid.ID == clientID {
				clientItemID = cid.ItemID
				break
			}
		}

		if clientItemID == "" {
			log.Printf("No matching client item ID found for series: %s", s.Data.Details.Title)
			continue
		}

		// Check if the series already exists in the database
		existingSeries, err := j.seriesRepo.GetByClientItemID(ctx, clientItemID, clientID)
		if err == nil {
			// Series exists, update it
			// Merge client IDs
			for _, cid := range s.ClientIDs {
				found := false
				for i, existingCid := range existingSeries.ClientIDs {
					if existingCid.ID == cid.ID && existingCid.Type == cid.Type {
						// Update existing entry if needed
						existingSeries.ClientIDs[i].ItemID = cid.ItemID
						found = true
						break
					}
				}
				if !found {
					// Add new client ID
					existingSeries.ClientIDs = append(existingSeries.ClientIDs, cid)
				}
			}

			// Merge external IDs
			for _, extID := range s.ExternalIDs {
				found := false
				for i, existingExtID := range existingSeries.ExternalIDs {
					if existingExtID.Source == extID.Source {
						// Update existing entry
						existingSeries.ExternalIDs[i].ID = extID.ID
						found = true
						break
					}
				}
				if !found {
					// Add new external ID
					existingSeries.ExternalIDs = append(existingSeries.ExternalIDs, extID)
				}
			}

			// Update data fields
			existingSeries.Data.Details = s.Data.Details
			existingSeries.Title = s.Data.Details.Title
			existingSeries.ReleaseDate = s.Data.Details.ReleaseDate
			existingSeries.ReleaseYear = s.Data.Details.ReleaseYear
			
			// Update additional series-specific fields
			existingSeries.Data.Genres = s.Data.Genres
			existingSeries.Data.Network = s.Data.Network
			existingSeries.Data.Status = s.Data.Status
			existingSeries.Data.ContentRating = s.Data.ContentRating
			existingSeries.Data.Rating = s.Data.Rating

			// Update seasons if available
			if len(s.Data.Seasons) > 0 {
				existingSeries.Data.Seasons = s.Data.Seasons
				existingSeries.Data.SeasonCount = s.Data.SeasonCount
			} else if seriesProvider != nil {
				// Try to fetch seasons if they're not already loaded
				var seriesID string
				for _, cid := range s.ClientIDs {
					if cid.ID == clientID {
						seriesID = cid.ItemID
						break
					}
				}
				
				if seriesID != "" {
					// Fetch seasons for this series
					seasons, err := seriesProvider.GetSeriesSeasons(ctx, seriesID)
					if err == nil && len(seasons) > 0 {
						// Convert to Season type from pointer
						seriesSeasons := make([]mediatypes.Season, 0, len(seasons))
						for _, season := range seasons {
							if season.Data != nil {
								seriesSeasons = append(seriesSeasons, *season.Data)
							}
						}
						
						existingSeries.Data.Seasons = seriesSeasons
						existingSeries.Data.SeasonCount = len(seriesSeasons)
						
						// Update episode count by summing episode counts from seasons
						totalEpisodes := 0
						for _, season := range seriesSeasons {
							totalEpisodes += season.EpisodeCount
						}
						
						if totalEpisodes > 0 {
							existingSeries.Data.EpisodeCount = totalEpisodes
						}
					}
				}
			}

			// Update episode count
			if s.Data.EpisodeCount > 0 {
				existingSeries.Data.EpisodeCount = s.Data.EpisodeCount
			}

			// Save the updated series
			_, err = j.seriesRepo.Update(ctx, *existingSeries)
			if err != nil {
				log.Printf("Error updating series: %v", err)
				continue
			}
		} else {
			// Series doesn't exist, create it
			// Set top level title and release fields
			s.Title = s.Data.Details.Title
			s.ReleaseDate = s.Data.Details.ReleaseDate
			s.ReleaseYear = s.Data.Details.ReleaseYear
			
			// Make sure we have genres data initialized
			if s.Data.Genres == nil {
				s.Data.Genres = []string{}
			}
			
			if s.Data.Seasons == nil {
				s.Data.Seasons = []mediatypes.Season{}
			}

			// Create the series
			_, err = j.seriesRepo.Create(ctx, s)
			if err != nil {
				log.Printf("Error creating series: %v", err)
				continue
			}
		}
	}

	return nil
}

// processEpisodeBatch processes a batch of episodes and saves them to the database
func (j *MediaSyncJob) processEpisodeBatch(ctx context.Context, episodes []models.MediaItem[*mediatypes.Episode], clientID uint64, clientType clienttypes.MediaClientType) error {
	for _, episode := range episodes {
		// Skip if episode has no client ID information
		if len(episode.ClientIDs) == 0 {
			log.Printf("Skipping episode with no client IDs: %s", episode.Data.Details.Title)
			continue
		}

		// Get the client ID and item ID for lookup
		clientItemID := ""
		for _, cid := range episode.ClientIDs {
			if cid.ID == clientID {
				clientItemID = cid.ItemID
				break
			}
		}

		if clientItemID == "" {
			log.Printf("No matching client item ID found for episode: %s", episode.Data.Details.Title)
			continue
		}

		// Check if the episode already exists in the database
		existingEpisode, err := j.episodeRepo.GetByClientItemID(ctx, clientItemID, clientID)
		if err == nil {
			// Episode exists, update it
			// Merge client IDs
			for _, cid := range episode.ClientIDs {
				found := false
				for i, existingCid := range existingEpisode.ClientIDs {
					if existingCid.ID == cid.ID && existingCid.Type == cid.Type {
						// Update existing entry if needed
						existingEpisode.ClientIDs[i].ItemID = cid.ItemID
						found = true
						break
					}
				}
				if !found {
					// Add new client ID
					existingEpisode.ClientIDs = append(existingEpisode.ClientIDs, cid)
				}
			}

			// Merge external IDs
			for _, extID := range episode.ExternalIDs {
				found := false
				for i, existingExtID := range existingEpisode.ExternalIDs {
					if existingExtID.Source == extID.Source {
						// Update existing entry
						existingEpisode.ExternalIDs[i].ID = extID.ID
						found = true
						break
					}
				}
				if !found {
					// Add new external ID
					existingEpisode.ExternalIDs = append(existingEpisode.ExternalIDs, extID)
				}
			}

			// Update data fields
			existingEpisode.Data.Details = episode.Data.Details
			existingEpisode.Title = episode.Data.Details.Title
			existingEpisode.ReleaseDate = episode.Data.Details.ReleaseDate
			existingEpisode.ReleaseYear = episode.Data.Details.ReleaseYear

			// Save the updated episode
			_, err = j.episodeRepo.Update(ctx, *existingEpisode)
			if err != nil {
				log.Printf("Error updating episode: %v", err)
				continue
			}
		} else {
			// Episode doesn't exist, create it
			// Set top level title and release fields
			episode.Title = episode.Data.Details.Title
			episode.ReleaseDate = episode.Data.Details.ReleaseDate
			episode.ReleaseYear = episode.Data.Details.ReleaseYear

			// Create the episode
			_, err = j.episodeRepo.Create(ctx, episode)
			if err != nil {
				log.Printf("Error creating episode: %v", err)
				continue
			}
		}
	}

	return nil
}

// processTrackBatch processes a batch of music tracks and saves them to the database
func (j *MediaSyncJob) processTrackBatch(ctx context.Context, tracks []models.MediaItem[*mediatypes.Track], clientID uint64, clientType clienttypes.MediaClientType) error {
	for _, track := range tracks {
		// Skip if track has no client ID information
		if len(track.ClientIDs) == 0 {
			log.Printf("Skipping track with no client IDs: %s", track.Data.Details.Title)
			continue
		}

		// Get the client ID and item ID for lookup
		clientItemID := ""
		for _, cid := range track.ClientIDs {
			if cid.ID == clientID {
				clientItemID = cid.ItemID
				break
			}
		}

		if clientItemID == "" {
			log.Printf("No matching client item ID found for track: %s", track.Data.Details.Title)
			continue
		}

		// Check if the track already exists in the database
		existingTrack, err := j.musicRepo.GetByClientItemID(ctx, clientItemID, clientID)
		if err == nil {
			// Track exists, update it
			// Merge client IDs
			for _, cid := range track.ClientIDs {
				found := false
				for i, existingCid := range existingTrack.ClientIDs {
					if existingCid.ID == cid.ID && existingCid.Type == cid.Type {
						// Update existing entry if needed
						existingTrack.ClientIDs[i].ItemID = cid.ItemID
						found = true
						break
					}
				}
				if !found {
					// Add new client ID
					existingTrack.ClientIDs = append(existingTrack.ClientIDs, cid)
				}
			}

			// Merge external IDs
			for _, extID := range track.ExternalIDs {
				found := false
				for i, existingExtID := range existingTrack.ExternalIDs {
					if existingExtID.Source == extID.Source {
						// Update existing entry
						existingTrack.ExternalIDs[i].ID = extID.ID
						found = true
						break
					}
				}
				if !found {
					// Add new external ID
					existingTrack.ExternalIDs = append(existingTrack.ExternalIDs, extID)
				}
			}

			// Update data fields
			existingTrack.Data = track.Data
			existingTrack.Title = track.Data.Details.Title
			existingTrack.ReleaseDate = track.Data.Details.ReleaseDate
			existingTrack.ReleaseYear = track.Data.Details.ReleaseYear

			// Save the updated track
			_, err = j.musicRepo.Update(ctx, *existingTrack)
			if err != nil {
				log.Printf("Error updating track: %v", err)
				continue
			}
		} else {
			// Track doesn't exist, create it
			// Set top level title and release fields
			track.Title = track.Data.Details.Title
			track.ReleaseDate = track.Data.Details.ReleaseDate
			track.ReleaseYear = track.Data.Details.ReleaseYear

			// Create the track
			_, err = j.musicRepo.Create(ctx, track)
			if err != nil {
				log.Printf("Error creating track: %v", err)
				continue
			}
		}
	}

	return nil
}

// processAlbumBatch processes a batch of music albums and saves them to the database
func (j *MediaSyncJob) processAlbumBatch(ctx context.Context, albums []models.MediaItem[*mediatypes.Album], clientID uint64, clientType clienttypes.MediaClientType) error {
	for _, album := range albums {
		// Skip if album has no client ID information
		if len(album.ClientIDs) == 0 {
			log.Printf("Skipping album with no client IDs: %s", album.Data.Details.Title)
			continue
		}

		// Get the client ID and item ID for lookup
		clientItemID := ""
		for _, cid := range album.ClientIDs {
			if cid.ID == clientID {
				clientItemID = cid.ItemID
				break
			}
		}

		if clientItemID == "" {
			log.Printf("No matching client item ID found for album: %s", album.Data.Details.Title)
			continue
		}

		// Check if the album already exists in the database
		existingAlbum, err := j.albumRepo.GetByClientItemID(ctx, clientItemID, clientID)
		if err == nil {
			// Album exists, update it
			// Merge client IDs
			for _, cid := range album.ClientIDs {
				found := false
				for i, existingCid := range existingAlbum.ClientIDs {
					if existingCid.ID == cid.ID && existingCid.Type == cid.Type {
						// Update existing entry if needed
						existingAlbum.ClientIDs[i].ItemID = cid.ItemID
						found = true
						break
					}
				}
				if !found {
					// Add new client ID
					existingAlbum.ClientIDs = append(existingAlbum.ClientIDs, cid)
				}
			}

			// Merge external IDs
			for _, extID := range album.ExternalIDs {
				found := false
				for i, existingExtID := range existingAlbum.ExternalIDs {
					if existingExtID.Source == extID.Source {
						// Update existing entry
						existingAlbum.ExternalIDs[i].ID = extID.ID
						found = true
						break
					}
				}
				if !found {
					// Add new external ID
					existingAlbum.ExternalIDs = append(existingAlbum.ExternalIDs, extID)
				}
			}

			// Update data fields
			existingAlbum.Data = album.Data
			existingAlbum.Title = album.Data.Details.Title
			existingAlbum.ReleaseDate = album.Data.Details.ReleaseDate
			existingAlbum.ReleaseYear = album.Data.Details.ReleaseYear

			// Save the updated album
			_, err = j.albumRepo.Update(ctx, *existingAlbum)
			if err != nil {
				log.Printf("Error updating album: %v", err)
				continue
			}
		} else {
			// Album doesn't exist, create it
			// Set top level title and release fields
			album.Title = album.Data.Details.Title
			album.ReleaseDate = album.Data.Details.ReleaseDate
			album.ReleaseYear = album.Data.Details.ReleaseYear

			// Create the album
			_, err = j.albumRepo.Create(ctx, album)
			if err != nil {
				log.Printf("Error creating album: %v", err)
				continue
			}
		}
	}

	return nil
}

// processArtistBatch processes a batch of music artists and saves them to the database
func (j *MediaSyncJob) processArtistBatch(ctx context.Context, artists []models.MediaItem[*mediatypes.Artist], clientID uint64, clientType clienttypes.MediaClientType) error {
	for _, artist := range artists {
		// Skip if artist has no client ID information
		if len(artist.ClientIDs) == 0 {
			log.Printf("Skipping artist with no client IDs: %s", artist.Data.Details.Title)
			continue
		}

		// Get the client ID and item ID for lookup
		clientItemID := ""
		for _, cid := range artist.ClientIDs {
			if cid.ID == clientID {
				clientItemID = cid.ItemID
				break
			}
		}

		if clientItemID == "" {
			log.Printf("No matching client item ID found for artist: %s", artist.Data.Details.Title)
			continue
		}

		// Check if the artist already exists in the database
		existingArtist, err := j.artistRepo.GetByClientItemID(ctx, clientItemID, clientID)
		if err == nil {
			// Artist exists, update it
			// Merge client IDs
			for _, cid := range artist.ClientIDs {
				found := false
				for i, existingCid := range existingArtist.ClientIDs {
					if existingCid.ID == cid.ID && existingCid.Type == cid.Type {
						// Update existing entry if needed
						existingArtist.ClientIDs[i].ItemID = cid.ItemID
						found = true
						break
					}
				}
				if !found {
					// Add new client ID
					existingArtist.ClientIDs = append(existingArtist.ClientIDs, cid)
				}
			}

			// Merge external IDs
			for _, extID := range artist.ExternalIDs {
				found := false
				for i, existingExtID := range existingArtist.ExternalIDs {
					if existingExtID.Source == extID.Source {
						// Update existing entry
						existingArtist.ExternalIDs[i].ID = extID.ID
						found = true
						break
					}
				}
				if !found {
					// Add new external ID
					existingArtist.ExternalIDs = append(existingArtist.ExternalIDs, extID)
				}
			}

			// Update data fields
			existingArtist.Data = artist.Data
			existingArtist.Title = artist.Data.Details.Title

			// Save the updated artist
			_, err = j.artistRepo.Update(ctx, *existingArtist)
			if err != nil {
				log.Printf("Error updating artist: %v", err)
				continue
			}
		} else {
			// Artist doesn't exist, create it
			// Set top level title and release fields
			artist.Title = artist.Data.Details.Title

			// Create the artist
			_, err = j.artistRepo.Create(ctx, artist)
			if err != nil {
				log.Printf("Error creating artist: %v", err)
				continue
			}
		}
	}

	return nil
}
