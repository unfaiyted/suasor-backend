package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	mediatypes "suasor/client/media/types"
	"suasor/repository"
	"suasor/services/scheduler"
	"suasor/types/models"
)

// NewReleaseNotificationJob identifies and notifies users about new releases
type NewReleaseNotificationJob struct {
	jobRepo       repository.JobRepository
	userRepo      repository.UserRepository
	configRepo    repository.UserConfigRepository
	movieRepo     repository.MediaItemRepository[*mediatypes.Movie]
	seriesRepo    repository.MediaItemRepository[*mediatypes.Series]
	musicRepo     repository.MediaItemRepository[*mediatypes.Track]
	historyRepo   repository.MediaPlayHistoryRepository
	metadataClient interface{} // Using interface{} to avoid import cycles
}

// NewNewReleaseNotificationJob creates a new release notification job
func NewNewReleaseNotificationJob(
	jobRepo repository.JobRepository,
	userRepo repository.UserRepository,
	configRepo repository.UserConfigRepository,
	movieRepo repository.MediaItemRepository[*mediatypes.Movie],
	seriesRepo repository.MediaItemRepository[*mediatypes.Series],
	musicRepo repository.MediaItemRepository[*mediatypes.Track],
	historyRepo repository.MediaPlayHistoryRepository,
	metadataClient interface{},
) *NewReleaseNotificationJob {
	return &NewReleaseNotificationJob{
		jobRepo:        jobRepo,
		userRepo:       userRepo,
		configRepo:     configRepo,
		movieRepo:      movieRepo,
		seriesRepo:     seriesRepo,
		musicRepo:      musicRepo,
		historyRepo:    historyRepo,
		metadataClient: metadataClient,
	}
}

// Name returns the unique name of the job
func (j *NewReleaseNotificationJob) Name() string {
	return "system.newrelease.notification"
}

// Schedule returns when the job should next run
func (j *NewReleaseNotificationJob) Schedule() time.Duration {
	// Run daily by default
	return 24 * time.Hour
}

// Execute runs the new release notification job
func (j *NewReleaseNotificationJob) Execute(ctx context.Context) error {
	log.Println("Starting new release notification job")

	// Create a job run record
	now := time.Now()
	jobRun := &models.JobRun{
		JobName:   j.Name(),
		JobType:   models.JobTypeNotification,
		Status:    models.JobStatusRunning,
		StartTime: &now,
		Metadata:  fmt.Sprintf(`{"type":"newReleaseNotification","startTime":"%s"}`, now.Format(time.RFC3339)),
	}

	if err := j.jobRepo.CreateJobRun(ctx, jobRun); err != nil {
		log.Printf("Error creating job run record: %v", err)
		return err
	}

	// Fetch new releases from metadata providers
	newReleases, err := j.fetchNewReleases(ctx)
	if err != nil {
		j.completeJobRun(ctx, jobRun.ID, models.JobStatusFailed, fmt.Sprintf("Error fetching new releases: %v", err))
		return err
	}

	log.Printf("Found %d new movie releases, %d new series releases, and %d new music releases",
		len(newReleases.movies), len(newReleases.series), len(newReleases.music))

	// Get all users
	users, err := j.userRepo.FindAll(ctx)
	if err != nil {
		j.completeJobRun(ctx, jobRun.ID, models.JobStatusFailed, fmt.Sprintf("Error getting users: %v", err))
		return err
	}

	notificationStats := map[string]int{
		"usersNotified":      0,
		"movieNotifications": 0,
		"seriesNotifications": 0,
		"musicNotifications": 0,
		"totalNotifications": 0,
	}

	// Process each user
	for _, user := range users {
		if !user.Active {
			continue
		}

		// Get user configuration
		config, err := j.configRepo.GetUserConfig(ctx, user.ID)
		if err != nil {
			log.Printf("Error getting config for user %d: %v", user.ID, err)
			continue
		}

		if !config.NewReleaseNotificationsEnabled {
			continue
		}

		userStats, err := j.processUserNotifications(ctx, user, config, newReleases)
		if err != nil {
			log.Printf("Error processing notifications for user %d: %v", user.ID, err)
			continue
		}

		// Update stats
		if userStats.totalNotifications > 0 {
			notificationStats["usersNotified"]++
		}
		notificationStats["movieNotifications"] += userStats.movieNotifications
		notificationStats["seriesNotifications"] += userStats.seriesNotifications
		notificationStats["musicNotifications"] += userStats.musicNotifications
		notificationStats["totalNotifications"] += userStats.totalNotifications
	}

	// Complete the job
	statsJSON, _ := json.Marshal(notificationStats)
	j.completeJobRun(ctx, jobRun.ID, models.JobStatusCompleted, string(statsJSON))

	log.Println("New release notification job completed")
	return nil
}

// NewReleases holds new release information by media type
type NewReleases struct {
	movies []NewRelease
	series []NewRelease
	music  []NewRelease
}

// NewRelease holds information about a single new release
type NewRelease struct {
	ID          string
	Title       string
	ReleaseDate time.Time
	MediaType   string
	Genres      []string
	Creators    []string // Directors for movies, showrunners for series, artists for music
}

// NotificationStats holds statistics about notifications sent
type NotificationStats struct {
	movieNotifications  int
	seriesNotifications int
	musicNotifications  int
	totalNotifications  int
}

// completeJobRun finalizes a job run with status and results
func (j *NewReleaseNotificationJob) completeJobRun(ctx context.Context, jobRunID uint64, status models.JobStatus, message string) {
	if err := j.jobRepo.CompleteJobRun(ctx, jobRunID, status, message); err != nil {
		log.Printf("Error completing job run: %v", err)
	}
}

// fetchNewReleases fetches new releases from metadata providers
func (j *NewReleaseNotificationJob) fetchNewReleases(ctx context.Context) (NewReleases, error) {
	releases := NewReleases{}
	
	// In a real implementation, we would:
	// 1. Use the metadata client to fetch recent and upcoming releases
	// 2. Filter to include only releases from the past week and upcoming month
	// 3. Format into structured data for processing

	// Mock implementation
	now := time.Now()
	
	// Add some mock movie releases
	releases.movies = []NewRelease{
		{
			ID:          "movie1",
			Title:       "The Avengers: New Era",
			ReleaseDate: now.AddDate(0, 0, -5), // 5 days ago
			MediaType:   "movie",
			Genres:      []string{"Action", "Sci-Fi"},
			Creators:    []string{"Joe Russo", "Anthony Russo"},
		},
		{
			ID:          "movie2",
			Title:       "Space Odyssey 2023",
			ReleaseDate: now.AddDate(0, 0, 7), // 7 days in the future
			MediaType:   "movie",
			Genres:      []string{"Sci-Fi", "Drama"},
			Creators:    []string{"Christopher Nolan"},
		},
		{
			ID:          "movie3",
			Title:       "Comedy Central",
			ReleaseDate: now.AddDate(0, 0, -2), // 2 days ago
			MediaType:   "movie",
			Genres:      []string{"Comedy"},
			Creators:    []string{"Judd Apatow"},
		},
	}
	
	// Add some mock series releases
	releases.series = []NewRelease{
		{
			ID:          "series1",
			Title:       "Stranger Things 5",
			ReleaseDate: now.AddDate(0, 0, 14), // 14 days in the future
			MediaType:   "series",
			Genres:      []string{"Sci-Fi", "Horror", "Drama"},
			Creators:    []string{"Duffer Brothers"},
		},
		{
			ID:          "series2",
			Title:       "The Mandalorian Season 4",
			ReleaseDate: now.AddDate(0, 0, -3), // 3 days ago
			MediaType:   "series",
			Genres:      []string{"Sci-Fi", "Action", "Western"},
			Creators:    []string{"Jon Favreau"},
		},
	}
	
	// Add some mock music releases
	releases.music = []NewRelease{
		{
			ID:          "music1",
			Title:       "New Horizons",
			ReleaseDate: now.AddDate(0, 0, -1), // 1 day ago
			MediaType:   "album",
			Genres:      []string{"Rock", "Alternative"},
			Creators:    []string{"Imagine Dragons"},
		},
		{
			ID:          "music2",
			Title:       "Echoes of Tomorrow",
			ReleaseDate: now.AddDate(0, 0, 5), // 5 days in the future
			MediaType:   "album",
			Genres:      []string{"Electronic", "Ambient"},
			Creators:    []string{"Daft Punk"},
		},
	}

	return releases, nil
}

// processUserNotifications processes new release notifications for a user
func (j *NewReleaseNotificationJob) processUserNotifications(ctx context.Context, user models.User, config *models.UserConfig, newReleases NewReleases) (NotificationStats, error) {
	stats := NotificationStats{}
	log.Printf("Processing new release notifications for user %s", user.Username)

	// Get user preferences and history
	preferences, err := j.getUserPreferences(ctx, user.ID)
	if err != nil {
		return stats, fmt.Errorf("error getting user preferences: %w", err)
	}

	// Process movie notifications if enabled
	if j.isMediaTypeEnabled(config.NewReleaseMediaTypes, "movie") {
		for _, release := range newReleases.movies {
			if j.shouldNotifyUserAboutRelease(release, preferences) {
				j.sendNewReleaseNotification(ctx, user.ID, release)
				stats.movieNotifications++
				stats.totalNotifications++
			}
		}
	}

	// Process series notifications if enabled
	if j.isMediaTypeEnabled(config.NewReleaseMediaTypes, "series") {
		for _, release := range newReleases.series {
			if j.shouldNotifyUserAboutRelease(release, preferences) {
				j.sendNewReleaseNotification(ctx, user.ID, release)
				stats.seriesNotifications++
				stats.totalNotifications++
			}
		}
	}

	// Process music notifications if enabled
	if j.isMediaTypeEnabled(config.NewReleaseMediaTypes, "music") {
		for _, release := range newReleases.music {
			if j.shouldNotifyUserAboutRelease(release, preferences) {
				j.sendNewReleaseNotification(ctx, user.ID, release)
				stats.musicNotifications++
				stats.totalNotifications++
			}
		}
	}

	log.Printf("Sent %d notifications to user %s", stats.totalNotifications, user.Username)
	return stats, nil
}

// UserPreferences holds a user's preferences for recommendations
type UserPreferences struct {
	favoriteGenres    map[string][]string
	favoriteCreators  []string
	blacklistedGenres map[string][]string
	recentlyWatched   []string
}

// getUserPreferences gets a user's preferences for recommendations
func (j *NewReleaseNotificationJob) getUserPreferences(ctx context.Context, userID uint64) (UserPreferences, error) {
	preferences := UserPreferences{
		favoriteGenres: map[string][]string{},
		favoriteCreators: []string{},
		blacklistedGenres: map[string][]string{},
		recentlyWatched: []string{},
	}

	// In a real implementation, we would:
	// 1. Get the user's config to extract explicit preferences
	// 2. Analyze watch history to infer preferences
	// 3. Structure the data for easy comparison

	// Mock implementation
	preferences.favoriteGenres["movie"] = []string{"Action", "Sci-Fi", "Comedy"}
	preferences.favoriteGenres["series"] = []string{"Drama", "Sci-Fi", "Crime"}
	preferences.favoriteGenres["music"] = []string{"Rock", "Electronic"}
	
	preferences.favoriteCreators = []string{
		"Christopher Nolan",
		"Duffer Brothers",
		"Imagine Dragons",
	}
	
	preferences.blacklistedGenres["movie"] = []string{"Horror"}
	preferences.blacklistedGenres["series"] = []string{}
	preferences.blacklistedGenres["music"] = []string{"Country"}
	
	preferences.recentlyWatched = []string{
		"The Dark Knight",
		"Stranger Things",
		"Breaking Bad",
	}

	return preferences, nil
}

// isMediaTypeEnabled checks if a media type is enabled in the comma-separated list
func (j *NewReleaseNotificationJob) isMediaTypeEnabled(mediaTypes string, mediaType string) bool {
	// TODO: Implement this properly with string parsing
	// For now, always return true for testing
	return true
}

// shouldNotifyUserAboutRelease checks if a user should be notified about a release
func (j *NewReleaseNotificationJob) shouldNotifyUserAboutRelease(release NewRelease, preferences UserPreferences) bool {
	// In a real implementation, we would:
	// 1. Check if the release matches the user's genre preferences
	// 2. Check if the release is by one of the user's favorite creators
	// 3. Check if the release is related to content the user has watched
	// 4. Apply other filtering rules based on user preferences

	// Mock implementation - simplified matching
	
	// Check if any genre matches preferences
	matchesGenre := false
	for _, genre := range release.Genres {
		for _, favGenre := range preferences.favoriteGenres[release.MediaType] {
			if genre == favGenre {
				matchesGenre = true
				break
			}
		}
		if matchesGenre {
			break
		}
	}
	
	// Check if any creator matches preferences
	matchesCreator := false
	for _, creator := range release.Creators {
		for _, favCreator := range preferences.favoriteCreators {
			if creator == favCreator {
				matchesCreator = true
				break
			}
		}
		if matchesCreator {
			break
		}
	}
	
	// Check if any genre is blacklisted
	isBlacklisted := false
	for _, genre := range release.Genres {
		for _, blacklistedGenre := range preferences.blacklistedGenres[release.MediaType] {
			if genre == blacklistedGenre {
				isBlacklisted = true
				break
			}
		}
		if isBlacklisted {
			break
		}
	}
	
	// Return true if it matches preferences and isn't blacklisted
	return (matchesGenre || matchesCreator) && !isBlacklisted
}

// sendNewReleaseNotification sends a notification about a new release to a user
func (j *NewReleaseNotificationJob) sendNewReleaseNotification(ctx context.Context, userID uint64, release NewRelease) {
	// In a real implementation, we would:
	// 1. Create a notification record in the database
	// 2. Send an email or push notification if configured
	// 3. Track which notifications have been sent to avoid duplicates

	log.Printf("Sending notification to user %d about '%s'", userID, release.Title)

	// Build a notification message
	message := fmt.Sprintf("New %s release: %s", release.MediaType, release.Title)
	if !release.ReleaseDate.After(time.Now()) {
		message += " is now available"
	} else {
		daysUntil := int(release.ReleaseDate.Sub(time.Now()).Hours() / 24)
		message += fmt.Sprintf(" is coming in %d days", daysUntil)
	}

	// In a real implementation, we would store this notification
	log.Printf("Notification message: %s", message)
}

// SetupNewReleaseNotificationSchedule creates or updates a new release notification schedule
func (j *NewReleaseNotificationJob) SetupNewReleaseNotificationSchedule(ctx context.Context, frequency string) error {
	// Check if job already exists
	existing, err := j.jobRepo.GetJobSchedule(ctx, j.Name())
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
		JobName:     j.Name(),
		JobType:     models.JobTypeNotification,
		Frequency:   frequency,
		Enabled:     frequency != string(scheduler.FrequencyManual),
		LastRunTime: nil, // Never run yet
	}

	return j.jobRepo.CreateJobSchedule(ctx, schedule)
}

// RunManualNewReleaseCheck runs the new release notification job manually
func (j *NewReleaseNotificationJob) RunManualNewReleaseCheck(ctx context.Context) error {
	return j.Execute(ctx)
}

// GetUpcomingReleases gets a list of upcoming releases matching user preferences
func (j *NewReleaseNotificationJob) GetUpcomingReleases(ctx context.Context, userID uint64, mediaType string, days int) ([]map[string]interface{}, error) {
	log.Printf("Getting upcoming %s releases for user %d in the next %d days", mediaType, userID, days)

	// In a real implementation, we would:
	// 1. Get the user's preferences
	// 2. Fetch upcoming releases from metadata providers
	// 3. Filter to include only those matching preferences
	// 4. Format them into a list for display

	// Mock implementation
	releases := []map[string]interface{}{}
	
	now := time.Now()
	
	if mediaType == "movie" || mediaType == "all" {
		releases = append(releases, map[string]interface{}{
			"id": "movie123",
			"title": "The Matrix Resurrection",
			"releaseDate": now.AddDate(0, 0, 10).Format(time.RFC3339),
			"mediaType": "movie",
			"genres": []string{"Sci-Fi", "Action"},
			"directors": []string{"Lana Wachowski"},
			"description": "Neo returns to face a new threat.",
		})
	}
	
	if mediaType == "series" || mediaType == "all" {
		releases = append(releases, map[string]interface{}{
			"id": "series456",
			"title": "House of the Dragon Season 2",
			"releaseDate": now.AddDate(0, 0, 15).Format(time.RFC3339),
			"mediaType": "series",
			"genres": []string{"Fantasy", "Drama"},
			"creators": []string{"Ryan Condal", "George R.R. Martin"},
			"description": "The Targaryen civil war continues.",
		})
	}
	
	if mediaType == "music" || mediaType == "all" {
		releases = append(releases, map[string]interface{}{
			"id": "music789",
			"title": "Echoes of Tomorrow",
			"releaseDate": now.AddDate(0, 0, 5).Format(time.RFC3339),
			"mediaType": "album",
			"genres": []string{"Electronic", "Ambient"},
			"artists": []string{"Daft Punk"},
			"description": "A groundbreaking return after years of silence.",
		})
	}

	return releases, nil
}

// UpdateUserNotificationPreferences updates a user's new release notification preferences
func (j *NewReleaseNotificationJob) UpdateUserNotificationPreferences(ctx context.Context, userID uint64, preferences map[string]interface{}) error {
	// In a real implementation, we would update the user's config with the new preferences
	log.Printf("Updated notification preferences for user %d", userID)
	return nil
}