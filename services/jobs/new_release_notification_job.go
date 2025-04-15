package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	mediatypes "suasor/client/media/types"
	"suasor/client/metadata"
	"suasor/client/types"
	"suasor/repository"
	"suasor/services/scheduler"
	"suasor/types/models"
	"suasor/utils"
)

// NewReleaseNotificationJob identifies and notifies users about new releases
type NewReleaseNotificationJob struct {
	jobRepo            repository.JobRepository
	userRepo           repository.UserRepository
	configRepo         repository.UserConfigRepository
	movieRepo          repository.ClientMediaItemRepository[*mediatypes.Movie]
	seriesRepo         repository.ClientMediaItemRepository[*mediatypes.Series]
	musicRepo          repository.ClientMediaItemRepository[*mediatypes.Track]
	userMovieDataRepo  repository.UserMediaItemDataRepository[*mediatypes.Movie]
	userSeriesDataRepo repository.UserMediaItemDataRepository[*mediatypes.Series]
	userMusicDataRepo  repository.UserMediaItemDataRepository[*mediatypes.Track]
	metadataClient     interface{} // Using interface{} to avoid import cycles
}

// NewNewReleaseNotificationJob creates a new release notification job
func NewNewReleaseNotificationJob(
	jobRepo repository.JobRepository,
	userRepo repository.UserRepository,
	configRepo repository.UserConfigRepository,
	movieRepo repository.ClientMediaItemRepository[*mediatypes.Movie],
	seriesRepo repository.ClientMediaItemRepository[*mediatypes.Series],
	musicRepo repository.ClientMediaItemRepository[*mediatypes.Track],
	userMovieDataRepo repository.UserMediaItemDataRepository[*mediatypes.Movie],
	userSeriesDataRepo repository.UserMediaItemDataRepository[*mediatypes.Series],
	userMusicDataRepo repository.UserMediaItemDataRepository[*mediatypes.Track],
	metadataClient interface{},
) *NewReleaseNotificationJob {
	return &NewReleaseNotificationJob{
		jobRepo:            jobRepo,
		userRepo:           userRepo,
		configRepo:         configRepo,
		movieRepo:          movieRepo,
		seriesRepo:         seriesRepo,
		musicRepo:          musicRepo,
		userMovieDataRepo:  userMovieDataRepo,
		userSeriesDataRepo: userSeriesDataRepo,
		userMusicDataRepo:  userMusicDataRepo,
		metadataClient:     metadataClient,
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
	log := utils.LoggerFromContext(ctx)
	log.Info().Msg("Starting new release notification job")

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
		log.Error().Err(err).Msg("Error creating job run record")
		return err
	}

	// Set the total number of steps - fetch, user processing, notifications
	if err := j.jobRepo.SetJobTotalItems(ctx, jobRun.ID, 3); err != nil {
		log.Warn().Err(err).Msg("Failed to set total items for job")
	}

	// Step 1: Fetch new releases from metadata providers
	j.jobRepo.UpdateJobProgress(ctx, jobRun.ID, 25, "Searching for new releases")
	newReleases, err := j.fetchNewReleases(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching new releases")
		j.completeJobRun(ctx, jobRun.ID, models.JobStatusFailed, fmt.Sprintf("Failed to fetch new releases: %v", err))
		return err
	}

	// Update progress after first step
	j.jobRepo.IncrementJobProcessedItems(ctx, jobRun.ID, 1)

	// If no new releases were found, finish the job early
	if len(newReleases.Movies) == 0 && len(newReleases.Series) == 0 && len(newReleases.Music) == 0 {
		log.Info().Msg("No new releases found, completing job")
		notificationStats := NotificationStats{}
		statsJSON, _ := json.Marshal(notificationStats)
		j.completeJobRun(ctx, jobRun.ID, models.JobStatusCompleted, string(statsJSON))
		return nil
	}

	log.Info().
		Int("movieCount", len(newReleases.Movies)).
		Int("seriesCount", len(newReleases.Series)).
		Int("musicCount", len(newReleases.Music)).
		Msg("Found new releases")

	// Step 2: Get all active users with preferences
	j.jobRepo.UpdateJobProgress(ctx, jobRun.ID, 50, "Processing user preferences")
	users, err := j.userRepo.FindAllActive(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error getting active users")
		j.completeJobRun(ctx, jobRun.ID, models.JobStatusFailed, fmt.Sprintf("Failed to get active users: %v", err))
		return err
	}

	// Update progress after second step
	j.jobRepo.IncrementJobProcessedItems(ctx, jobRun.ID, 1)

	// Step 3: Process each user and create notifications
	j.jobRepo.UpdateJobProgress(ctx, jobRun.ID, 75, "Creating notifications")

	// Track overall notification statistics
	totalStats := NotificationStats{}

	// Process each user in parallel or batches
	for _, user := range users {
		// Get user config for preferences
		config, err := j.configRepo.GetUserConfig(ctx, user.ID)
		if err != nil {
			log.Error().Err(err).Uint64("userId", user.ID).Msg("Error getting user config")
			continue // Skip this user, but continue with others
		}

		// Skip users who have disabled notifications
		if !config.NotificationsEnabled || !config.NewReleaseNotificationsEnabled {
			log.Debug().Uint64("userId", user.ID).Msg("User has notifications disabled")
			continue
		}

		// Process notifications for this user
		userStats, err := j.processUserNotifications(ctx, user, config, newReleases)
		if err != nil {
			log.Error().Err(err).Uint64("userId", user.ID).Msg("Error processing user notifications")
			continue // Skip this user, but continue with others
		}

		// Update overall statistics
		totalStats.MovieNotifications += userStats.MovieNotifications
		totalStats.SeriesNotifications += userStats.SeriesNotifications
		totalStats.MusicNotifications += userStats.MusicNotifications
		totalStats.TotalNotifications += userStats.TotalNotifications

		if userStats.TotalNotifications > 0 {
			totalStats.UsersNotified++
		}
	}

	// Update progress after third step
	j.jobRepo.IncrementJobProcessedItems(ctx, jobRun.ID, 1)

	// Convert stats to map for backwards compatibility
	statsMap := map[string]int{
		"usersNotified":       totalStats.UsersNotified,
		"movieNotifications":  totalStats.MovieNotifications,
		"seriesNotifications": totalStats.SeriesNotifications,
		"musicNotifications":  totalStats.MusicNotifications,
		"totalNotifications":  totalStats.TotalNotifications,
	}

	// Complete the job
	statsJSON, _ := json.Marshal(statsMap)
	j.completeJobRun(ctx, jobRun.ID, models.JobStatusCompleted, string(statsJSON))

	log.Info().
		Int("usersNotified", totalStats.UsersNotified).
		Int("totalNotifications", totalStats.TotalNotifications).
		Msg("New release notification job completed")
	return nil
}

// completeJobRun finalizes a job run with status and results
func (j *NewReleaseNotificationJob) completeJobRun(ctx context.Context, jobRunID uint64, status models.JobStatus, message string) {
	log := utils.LoggerFromContext(ctx)
	if err := j.jobRepo.CompleteJobRun(ctx, jobRunID, status, message); err != nil {
		log.Printf("Error completing job run: %v", err)
	}
}

// SetupNewReleaseNotificationSchedule creates or updates a notification schedule for a user
func (j *NewReleaseNotificationJob) SetupNewReleaseNotificationSchedule(ctx context.Context, userID uint64, frequency string) error {
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
		JobType:     models.JobTypeNotification,
		Frequency:   frequency,
		Enabled:     frequency != string(scheduler.FrequencyManual),
		UserID:      &userID,
		LastRunTime: nil, // Never run yet
	}

	return j.jobRepo.CreateJobSchedule(ctx, schedule)
}

// processUserNotifications processes notifications for a single user
func (j *NewReleaseNotificationJob) processUserNotifications(ctx context.Context, user models.User, config *models.UserConfig, newReleases NewReleases) (NotificationStats, error) {
	log := utils.LoggerFromContext(ctx)
	stats := NotificationStats{}

	// Build the user preference profile
	profile, err := j.buildUserPreferenceProfile(ctx, user.ID, config)
	if err != nil {
		return stats, fmt.Errorf("error building user preference profile: %w", err)
	}

	// Prepare notifications list
	var notifications []UserNotification

	// Process movie notifications
	if profile.NotifyForMovies && len(newReleases.Movies) > 0 {
		movieNotifications, err := j.processMovieNotifications(ctx, user.ID, profile, newReleases.Movies)
		if err != nil {
			log.Error().Err(err).Uint64("userId", user.ID).Msg("Error processing movie notifications")
		} else {
			notifications = append(notifications, movieNotifications...)
			stats.MovieNotifications += len(movieNotifications)
		}
	}

	// Process series notifications
	if profile.NotifyForSeries && len(newReleases.Series) > 0 {
		seriesNotifications, err := j.processSeriesNotifications(ctx, user.ID, profile, newReleases.Series)
		if err != nil {
			log.Error().Err(err).Uint64("userId", user.ID).Msg("Error processing series notifications")
		} else {
			notifications = append(notifications, seriesNotifications...)
			stats.SeriesNotifications += len(seriesNotifications)
		}
	}

	// Process music notifications
	if profile.NotifyForMusic && len(newReleases.Music) > 0 {
		musicNotifications, err := j.processMusicNotifications(ctx, user.ID, profile, newReleases.Music)
		if err != nil {
			log.Error().Err(err).Uint64("userId", user.ID).Msg("Error processing music notifications")
		} else {
			notifications = append(notifications, musicNotifications...)
			stats.MusicNotifications += len(musicNotifications)
		}
	}

	// Sort and filter notifications
	filteredNotifications := j.filterAndPrioritizeNotifications(notifications, profile.MaxNotifications)

	// Save notifications
	if len(filteredNotifications) > 0 {
		if err := j.saveUserNotifications(ctx, filteredNotifications); err != nil {
			log.Error().Err(err).Uint64("userId", user.ID).Msg("Error saving notifications")
			return stats, err
		}

		// Update final count
		stats.TotalNotifications = len(filteredNotifications)
	}

	log.Info().
		Uint64("userId", user.ID).
		Int("movieNotifications", stats.MovieNotifications).
		Int("seriesNotifications", stats.SeriesNotifications).
		Int("musicNotifications", stats.MusicNotifications).
		Int("totalNotifications", stats.TotalNotifications).
		Msg("Processed user notifications")

	return stats, nil
}

// buildUserPreferenceProfile builds a profile of user preferences for notifications
func (j *NewReleaseNotificationJob) buildUserPreferenceProfile(ctx context.Context, userID uint64, config *models.UserConfig) (*UserPreferenceProfile, error) {
	log := utils.LoggerFromContext(ctx)

	profile := &UserPreferenceProfile{
		// Default to notifying for all content types
		NotifyForMovies: true,
		NotifyForSeries: true,
		NotifyForMusic:  true,

		// Initialize maps
		FavoriteActors:      make(map[string]float32),
		FavoriteDirectors:   make(map[string]float32),
		FavoriteShowrunners: make(map[string]float32),
		FavoriteArtists:     make(map[string]float32),
		OwnedMovieIDs:       make(map[string]bool),
		OwnedSeriesIDs:      make(map[string]bool),
		OwnedMusicIDs:       make(map[string]bool),

		// Default notification settings
		RatingThreshold:   5.0, // Default minimum rating
		MaxNotifications:  10,  // Default max notifications
		NotifyForUpcoming: true,
		NotifyForRecent:   true,
	}

	// Apply user configuration settings for content types
	if config.NotifyMediaTypes != "" {
		mediaTypes := strings.Split(config.NotifyMediaTypes, ",")

		// Reset all to false first
		profile.NotifyForMovies = false
		profile.NotifyForSeries = false
		profile.NotifyForMusic = false

		// Set the ones that are specified
		for _, mediaType := range mediaTypes {
			mediaType = strings.TrimSpace(mediaType)
			switch mediaType {
			case "movie":
				profile.NotifyForMovies = true
			case "series":
				profile.NotifyForSeries = true
			case "music":
				profile.NotifyForMusic = true
			}
		}
	}

	// Apply user configuration settings for genres
	if config.PreferredGenres != nil {
		profile.PreferredMovieGenres = config.PreferredGenres.Movies
		profile.PreferredSeriesGenres = config.PreferredGenres.Series
		profile.PreferredMusicGenres = config.PreferredGenres.Music
	}

	if config.ExcludedGenres != nil {
		profile.ExcludedMovieGenres = config.ExcludedGenres.Movies
		profile.ExcludedSeriesGenres = config.ExcludedGenres.Series
		profile.ExcludedMusicGenres = config.ExcludedGenres.Music
	}

	// Apply user notification settings
	if config.NotifyRatingThreshold > 0 {
		profile.RatingThreshold = config.NotifyRatingThreshold
	}

	if config.MaxNotificationsPerDay > 0 {
		profile.MaxNotifications = config.MaxNotificationsPerDay
	}

	profile.NotifyForUpcoming = config.NotifyUpcomingReleases
	profile.NotifyForRecent = config.NotifyRecentReleases

	// Load user's media library contents
	if profile.NotifyForMovies {
		// Get user's movies
		userMovies, err := j.movieRepo.GetByUserID(ctx, userID)
		if err != nil {
			log.Error().Err(err).Uint64("userId", userID).Msg("Error getting user movies")
		} else {
			// Process movies to extract favorites and owned content
			for _, movie := range userMovies {
				if movie.Data == nil {
					continue
				}

				// Add to owned movie IDs
				if movie.Data.Details.ExternalIDs != nil {
					for _, externalID := range movie.Data.Details.ExternalIDs {
						profile.OwnedMovieIDs[externalID.ID] = true
					}
				}

				// Process cast for favorite actors and directors
				for _, person := range movie.Data.Credits.GetCast() {
					if strings.EqualFold(person.Role, "actor") {
						profile.FavoriteActors[person.Name] += 1.0
					} else if strings.EqualFold(person.Role, "director") {
						profile.FavoriteDirectors[person.Name] += 1.0
					}
				}
			}
		}
	}

	if profile.NotifyForSeries {
		// Get user's series
		userSeries, err := j.seriesRepo.GetByUserID(ctx, userID)
		if err != nil {
			log.Error().Err(err).Uint64("userId", userID).Msg("Error getting user series")
		} else {
			// Process series to extract favorites and owned content
			for _, series := range userSeries {
				if series.Data == nil {
					continue
				}

				// Add to owned series IDs
				if series.Data.Details.ExternalIDs != nil {
					for _, externalID := range series.Data.Details.ExternalIDs {
						profile.OwnedSeriesIDs[externalID.ID] = true
					}
				}
			}
		}
	}

	if profile.NotifyForMusic {
		// Get user's music
		userMusic, err := j.musicRepo.GetByUserID(ctx, userID)
		if err != nil {
			log.Error().Err(err).Uint64("userId", userID).Msg("Error getting user music")
		} else {
			// Process music to extract favorites and owned content
			for _, track := range userMusic {
				if track.Data == nil {
					continue
				}

				// Add artist to favorites
				if track.Data.ArtistName != "" {
					profile.FavoriteArtists[track.Data.ArtistName] += 1.0
				}
			}
		}
	}

	return profile, nil
}

// processMovieNotifications processes movie notifications for a user
func (j *NewReleaseNotificationJob) processMovieNotifications(ctx context.Context, userID uint64, profile *UserPreferenceProfile, movies []NewRelease) ([]UserNotification, error) {
	log := utils.LoggerFromContext(ctx)
	var notifications []UserNotification

	// Filter movies based on user preferences
	for _, movie := range movies {
		// Skip if the movie is not recent or upcoming (based on user preferences)
		if !j.shouldNotifyForRelease(movie.ReleaseDate, profile) {
			continue
		}

		// Skip if the movie is already in the user's library
		if profile.OwnedMovieIDs[movie.ExternalID] {
			continue
		}

		// Skip movies with ratings below threshold
		if movie.Rating < profile.RatingThreshold {
			continue
		}

		// Check for excluded genres
		excluded := false
		for _, genre := range movie.Genres {
			for _, excludedGenre := range profile.ExcludedMovieGenres {
				if strings.EqualFold(genre, excludedGenre) {
					excluded = true
					break
				}
			}
			if excluded {
				break
			}
		}

		if excluded {
			continue
		}

		// Check for preferred genres (if specified)
		if len(profile.PreferredMovieGenres) > 0 {
			preferred := false
			for _, genre := range movie.Genres {
				for _, preferredGenre := range profile.PreferredMovieGenres {
					if strings.EqualFold(genre, preferredGenre) {
						preferred = true
						break
					}
				}
				if preferred {
					break
				}
			}

			// Skip if no preferred genres match
			if !preferred {
				continue
			}
		}

		// Calculate priority based on user preferences
		priority := 3 // Default priority

		// Increase priority for favorite creators
		for _, creator := range movie.Creators {
			if weight, found := profile.FavoriteDirectors[creator]; found {
				priority += int(weight)
				break
			}
		}

		// Adjust priority based on rating
		if movie.Rating >= 8.0 {
			priority++
		}

		// Cap priority at 5
		if priority > 5 {
			priority = 5
		}

		// Create notification
		notificationType := NotificationTypeNewRelease
		if movie.ReleaseDate.After(time.Now()) {
			notificationType = NotificationTypeUpcoming
		}

		notification := UserNotification{
			UserID:      userID,
			Title:       "New Movie: " + movie.Title,
			Message:     j.formatMovieNotificationMessage(movie),
			Type:        notificationType,
			ContentType: "movie",
			ContentID:   movie.ID,
			ImageURL:    movie.ImageURL,
			ActionURL:   "", // Should be filled in with a URL to view the movie details
			Created:     time.Now(),
			Expires:     time.Now().Add(7 * 24 * time.Hour), // Expire after 1 week
			Priority:    priority,
			Read:        false,
			Dismissed:   false,
		}

		notifications = append(notifications, notification)
	}

	log.Info().
		Uint64("userId", userID).
		Int("notificationCount", len(notifications)).
		Int("totalMovies", len(movies)).
		Msg("Processed movie notifications")

	return notifications, nil
}

// processSeriesNotifications processes series notifications for a user
func (j *NewReleaseNotificationJob) processSeriesNotifications(ctx context.Context, userID uint64, profile *UserPreferenceProfile, series []NewRelease) ([]UserNotification, error) {
	log := utils.LoggerFromContext(ctx)
	var notifications []UserNotification

	// Filter series based on user preferences
	for _, show := range series {
		// Skip if the show is not recent or upcoming (based on user preferences)
		if !j.shouldNotifyForRelease(show.ReleaseDate, profile) {
			continue
		}

		// Skip if the show is already in the user's library
		if profile.OwnedSeriesIDs[show.ExternalID] {
			continue
		}

		// Skip shows with ratings below threshold
		if show.Rating < profile.RatingThreshold {
			continue
		}

		// Check for excluded genres
		excluded := false
		for _, genre := range show.Genres {
			for _, excludedGenre := range profile.ExcludedSeriesGenres {
				if strings.EqualFold(genre, excludedGenre) {
					excluded = true
					break
				}
			}
			if excluded {
				break
			}
		}

		if excluded {
			continue
		}

		// Check for preferred genres (if specified)
		if len(profile.PreferredSeriesGenres) > 0 {
			preferred := false
			for _, genre := range show.Genres {
				for _, preferredGenre := range profile.PreferredSeriesGenres {
					if strings.EqualFold(genre, preferredGenre) {
						preferred = true
						break
					}
				}
				if preferred {
					break
				}
			}

			// Skip if no preferred genres match
			if !preferred {
				continue
			}
		}

		// Calculate priority based on user preferences
		priority := 3 // Default priority

		// Increase priority for favorite creators
		for _, creator := range show.Creators {
			if weight, found := profile.FavoriteShowrunners[creator]; found {
				priority += int(weight)
				break
			}
		}

		// Adjust priority based on rating
		if show.Rating >= 8.0 {
			priority++
		}

		// Cap priority at 5
		if priority > 5 {
			priority = 5
		}

		// Create notification
		notificationType := NotificationTypeNewRelease
		if show.ReleaseDate.After(time.Now()) {
			notificationType = NotificationTypeUpcoming
		}

		notification := UserNotification{
			UserID:      userID,
			Title:       "New Series: " + show.Title,
			Message:     j.formatSeriesNotificationMessage(show),
			Type:        notificationType,
			ContentType: "series",
			ContentID:   show.ID,
			ImageURL:    show.ImageURL,
			ActionURL:   "", // Should be filled in with a URL to view the series details
			Created:     time.Now(),
			Expires:     time.Now().Add(7 * 24 * time.Hour), // Expire after 1 week
			Priority:    priority,
			Read:        false,
			Dismissed:   false,
		}

		notifications = append(notifications, notification)
	}

	log.Info().
		Uint64("userId", userID).
		Int("notificationCount", len(notifications)).
		Int("totalSeries", len(series)).
		Msg("Processed series notifications")

	return notifications, nil
}

// processMusicNotifications processes music notifications for a user
func (j *NewReleaseNotificationJob) processMusicNotifications(ctx context.Context, userID uint64, profile *UserPreferenceProfile, music []NewRelease) ([]UserNotification, error) {
	log := utils.LoggerFromContext(ctx)
	var notifications []UserNotification

	// Filter music based on user preferences
	for _, track := range music {
		// Skip if the track is not recent or upcoming (based on user preferences)
		if !j.shouldNotifyForRelease(track.ReleaseDate, profile) {
			continue
		}

		// Skip if the track is already in the user's library
		if profile.OwnedMusicIDs[track.ExternalID] {
			continue
		}

		// Check for excluded genres
		excluded := false
		for _, genre := range track.Genres {
			for _, excludedGenre := range profile.ExcludedMusicGenres {
				if strings.EqualFold(genre, excludedGenre) {
					excluded = true
					break
				}
			}
			if excluded {
				break
			}
		}

		if excluded {
			continue
		}

		// Check for preferred genres (if specified)
		if len(profile.PreferredMusicGenres) > 0 {
			preferred := false
			for _, genre := range track.Genres {
				for _, preferredGenre := range profile.PreferredMusicGenres {
					if strings.EqualFold(genre, preferredGenre) {
						preferred = true
						break
					}
				}
				if preferred {
					break
				}
			}

			// Skip if no preferred genres match
			if !preferred {
				continue
			}
		}

		// Calculate priority based on user preferences
		priority := 3 // Default priority

		// Increase priority for favorite artists
		for _, creator := range track.Creators {
			if weight, found := profile.FavoriteArtists[creator]; found {
				priority += int(weight)
				break
			}
		}

		// Cap priority at 5
		if priority > 5 {
			priority = 5
		}

		// Create notification
		notificationType := NotificationTypeNewRelease
		if track.ReleaseDate.After(time.Now()) {
			notificationType = NotificationTypeUpcoming
		}

		notification := UserNotification{
			UserID:      userID,
			Title:       "New Music: " + track.Title,
			Message:     j.formatMusicNotificationMessage(track),
			Type:        notificationType,
			ContentType: "music",
			ContentID:   track.ID,
			ImageURL:    track.ImageURL,
			ActionURL:   "", // Should be filled in with a URL to view the music details
			Created:     time.Now(),
			Expires:     time.Now().Add(7 * 24 * time.Hour), // Expire after 1 week
			Priority:    priority,
			Read:        false,
			Dismissed:   false,
		}

		notifications = append(notifications, notification)
	}

	log.Info().
		Uint64("userId", userID).
		Int("notificationCount", len(notifications)).
		Int("totalMusic", len(music)).
		Msg("Processed music notifications")

	return notifications, nil
}

// shouldNotifyForRelease checks if a release date should trigger a notification based on user preferences
func (j *NewReleaseNotificationJob) shouldNotifyForRelease(releaseDate time.Time, profile *UserPreferenceProfile) bool {
	now := time.Now()

	// Check if it's an upcoming release
	if releaseDate.After(now) {
		// Only notify if user wants notifications for upcoming
		return profile.NotifyForUpcoming
	}

	// It's a recent release, check how recent
	daysSinceRelease := now.Sub(releaseDate).Hours() / 24

	// We consider anything in the last 14 days as recent
	isRecent := daysSinceRelease <= 14

	// Only notify if user wants notifications for recent releases
	return isRecent && profile.NotifyForRecent
}

// filterAndPrioritizeNotifications sorts and limits notifications
func (j *NewReleaseNotificationJob) filterAndPrioritizeNotifications(notifications []UserNotification, maxCount int) []UserNotification {
	// If we don't have more than the max, return all
	if len(notifications) <= maxCount {
		return notifications
	}

	// Sort by priority (higher first) and then by release date (newer first)
	sort.Slice(notifications, func(i, j int) bool {
		if notifications[i].Priority != notifications[j].Priority {
			return notifications[i].Priority > notifications[j].Priority
		}
		return notifications[i].Created.After(notifications[j].Created)
	})

	// Return only the top maxCount notifications
	return notifications[:maxCount]
}

// formatMovieNotificationMessage formats a notification message for a movie
func (j *NewReleaseNotificationJob) formatMovieNotificationMessage(movie NewRelease) string {
	var sb strings.Builder

	// Determine if upcoming or recently released
	if movie.ReleaseDate.After(time.Now()) {
		daysUntil := int(movie.ReleaseDate.Sub(time.Now()).Hours() / 24)

		if daysUntil <= 1 {
			sb.WriteString("Coming tomorrow! ")
		} else {
			sb.WriteString(fmt.Sprintf("Coming in %d days! ", daysUntil))
		}
	} else {
		daysSince := int(time.Since(movie.ReleaseDate).Hours() / 24)

		if daysSince <= 1 {
			sb.WriteString("Just released! ")
		} else {
			sb.WriteString(fmt.Sprintf("Released %d days ago. ", daysSince))
		}
	}

	// Add description if available
	if movie.Description != "" {
		// Truncate to a reasonable length
		if len(movie.Description) > 100 {
			sb.WriteString(movie.Description[:100] + "...")
		} else {
			sb.WriteString(movie.Description)
		}
	}

	// Add genres
	if len(movie.Genres) > 0 {
		sb.WriteString(" Genres: ")
		for i, genre := range movie.Genres {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(genre)

			// Limit to 3 genres
			if i == 2 {
				break
			}
		}
	}

	// Add rating if good
	if movie.Rating >= 7.0 {
		sb.WriteString(fmt.Sprintf(" Rating: %.1f/10", movie.Rating))
	}

	// Add directors if available
	if len(movie.Creators) > 0 {
		sb.WriteString(". Director: ")
		if len(movie.Creators) > 1 {
			sb.WriteString(movie.Creators[0] + " and others")
		} else {
			sb.WriteString(movie.Creators[0])
		}
	}

	return sb.String()
}

// formatSeriesNotificationMessage formats a notification message for a series
func (j *NewReleaseNotificationJob) formatSeriesNotificationMessage(series NewRelease) string {
	var sb strings.Builder

	// Determine if upcoming or recently released
	if series.ReleaseDate.After(time.Now()) {
		daysUntil := int(series.ReleaseDate.Sub(time.Now()).Hours() / 24)

		if daysUntil <= 1 {
			sb.WriteString("Premiering tomorrow! ")
		} else {
			sb.WriteString(fmt.Sprintf("Premiering in %d days! ", daysUntil))
		}
	} else {
		daysSince := int(time.Since(series.ReleaseDate).Hours() / 24)

		if daysSince <= 1 {
			sb.WriteString("Just premiered! ")
		} else {
			sb.WriteString(fmt.Sprintf("Premiered %d days ago. ", daysSince))
		}
	}

	// Add description if available
	if series.Description != "" {
		// Truncate to a reasonable length
		if len(series.Description) > 100 {
			sb.WriteString(series.Description[:100] + "...")
		} else {
			sb.WriteString(series.Description)
		}
	}

	// Add genres
	if len(series.Genres) > 0 {
		sb.WriteString(" Genres: ")
		for i, genre := range series.Genres {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(genre)

			// Limit to 3 genres
			if i == 2 {
				break
			}
		}
	}

	// Add rating if good
	if series.Rating >= 7.0 {
		sb.WriteString(fmt.Sprintf(" Rating: %.1f/10", series.Rating))
	}

	// Add creators if available
	if len(series.Creators) > 0 {
		sb.WriteString(". Creator: ")
		if len(series.Creators) > 1 {
			sb.WriteString(series.Creators[0] + " and others")
		} else {
			sb.WriteString(series.Creators[0])
		}
	}

	return sb.String()
}

// formatMusicNotificationMessage formats a notification message for music
func (j *NewReleaseNotificationJob) formatMusicNotificationMessage(music NewRelease) string {
	var sb strings.Builder

	// Determine if upcoming or recently released
	if music.ReleaseDate.After(time.Now()) {
		daysUntil := int(music.ReleaseDate.Sub(time.Now()).Hours() / 24)

		if daysUntil <= 1 {
			sb.WriteString("Releasing tomorrow! ")
		} else {
			sb.WriteString(fmt.Sprintf("Releasing in %d days! ", daysUntil))
		}
	} else {
		daysSince := int(time.Since(music.ReleaseDate).Hours() / 24)

		if daysSince <= 1 {
			sb.WriteString("Just released! ")
		} else {
			sb.WriteString(fmt.Sprintf("Released %d days ago. ", daysSince))
		}
	}

	// Add description if available
	if music.Description != "" {
		// Truncate to a reasonable length
		if len(music.Description) > 100 {
			sb.WriteString(music.Description[:100] + "...")
		} else {
			sb.WriteString(music.Description)
		}
	}

	// Add genres
	if len(music.Genres) > 0 {
		sb.WriteString(" Genres: ")
		for i, genre := range music.Genres {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(genre)

			// Limit to 3 genres
			if i == 2 {
				break
			}
		}
	}

	// Add artists if available
	if len(music.Creators) > 0 {
		sb.WriteString(". Artist: ")
		if len(music.Creators) > 1 {
			sb.WriteString(music.Creators[0] + " and others")
		} else {
			sb.WriteString(music.Creators[0])
		}
	}

	return sb.String()
}

// saveUserNotifications saves notifications for a user
func (j *NewReleaseNotificationJob) saveUserNotifications(ctx context.Context, notifications []UserNotification) error {
	log := utils.LoggerFromContext(ctx)

	// In a real implementation, this would save to a database or other notification system
	// For now, we'll just log the notifications

	for _, notification := range notifications {
		log.Info().
			Uint64("userId", notification.UserID).
			Str("title", notification.Title).
			Str("message", notification.Message).
			Str("type", string(notification.Type)).
			Str("contentType", notification.ContentType).
			Str("contentId", notification.ContentID).
			Int("priority", notification.Priority).
			Msg("Would save notification")
	}

	// In a production implementation, we would save these notifications to a database
	// or send them to a notification service

	return nil
}

// fetchNewReleases fetches new releases from metadata providers
func (j *NewReleaseNotificationJob) fetchNewReleases(ctx context.Context) (NewReleases, error) {
	log := utils.LoggerFromContext(ctx)
	log.Info().Msg("Fetching new releases from metadata providers")

	// Initialize the result
	var result NewReleases

	// Get metadata client
	metadataClient, ok := j.metadataClient.(metadata.MetadataClient)
	if !ok {
		return result, fmt.Errorf("metadata client not properly initialized")
	}

	// Configure how far ahead to look for upcoming movies/shows (7 days by default)
	upcomingDays := 7

	// Configure how recently released items to include (7 days by default)
	recentReleaseDays := 7

	// 1. Fetch upcoming and new movies
	if metadataClient.SupportsMovieMetadata() {
		// Get upcoming movies
		upcomingMovies, err := metadataClient.GetUpcomingMovies(ctx, upcomingDays)
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch upcoming movies")
			// Continue with other content types
		} else {
			log.Info().Int("count", len(upcomingMovies)).Msg("Found upcoming movies")

			// Convert to our NewRelease format
			for _, movie := range upcomingMovies {
				// Parse release date
				releaseDate, err := time.Parse("2006-01-02", movie.ReleaseDate)
				if err != nil {
					// If we can't parse the date, use today
					releaseDate = time.Now()
				}

				// Check if this is recent enough
				daysSinceRelease := time.Since(releaseDate).Hours() / 24
				if daysSinceRelease > float64(recentReleaseDays) {
					continue // Skip older releases
				}

				// Extract genres
				var genres []string
				for _, genre := range movie.Genres {
					genres = append(genres, genre.Name)
				}

				// Extract creators (directors)
				var creators []string
				if movie.Credits.Crew != nil {
					for _, crewMember := range movie.Credits.Crew {
						if crewMember.Job == "Director" {
							creators = append(creators, crewMember.Name)
						}
					}
				}

				// Create the new release
				newRelease := NewRelease{
					ID:          fmt.Sprintf("movie-%s", movie.ID),
					ExternalID:  movie.ID,
					Title:       movie.Title,
					Description: movie.Overview,
					ReleaseDate: releaseDate,
					MediaType:   "movie",
					Genres:      genres,
					Creators:    creators,
					Rating:      movie.VoteAverage,
					Source:      "tmdb", // Assuming TMDB for now
					ImageURL:    movie.PosterPath,
					Metadata:    movie,
				}

				result.Movies = append(result.Movies, newRelease)
			}
		}

		// Get now playing movies
		nowPlayingMovies, err := metadataClient.GetNowPlayingMovies(ctx, recentReleaseDays)
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch now playing movies")
		} else {
			log.Info().Int("count", len(nowPlayingMovies)).Msg("Found now playing movies")

			// Convert to our NewRelease format
			for _, movie := range nowPlayingMovies {
				// Check if we already have this movie in the upcoming list
				alreadyAdded := false
				for _, existingMovie := range result.Movies {
					if existingMovie.ExternalID == movie.ID {
						alreadyAdded = true
						break
					}
				}

				if alreadyAdded {
					continue
				}

				// Parse release date
				releaseDate, err := time.Parse("2006-01-02", movie.ReleaseDate)
				if err != nil {
					// If we can't parse the date, use today
					releaseDate = time.Now()
				}

				// Check if this is recent enough
				daysSinceRelease := time.Since(releaseDate).Hours() / 24
				if daysSinceRelease > float64(recentReleaseDays) {
					continue // Skip older releases
				}

				// Extract genres
				var genres []string
				for _, genre := range movie.Genres {
					genres = append(genres, genre.Name)
				}

				// Extract creators (directors)
				var creators []string
				if movie.Credits.Crew != nil {
					for _, crewMember := range movie.Credits.Crew {
						if crewMember.Job == "Director" {
							creators = append(creators, crewMember.Name)
						}
					}
				}

				// Create the new release
				newRelease := NewRelease{
					ID:          fmt.Sprintf("movie-%s", movie.ID),
					ExternalID:  movie.ID,
					Title:       movie.Title,
					Description: movie.Overview,
					ReleaseDate: releaseDate,
					MediaType:   "movie",
					Genres:      genres,
					Creators:    creators,
					Rating:      movie.VoteAverage,
					Source:      "tmdb", // Assuming TMDB for now
					ImageURL:    movie.PosterPath,
					Metadata:    movie,
				}

				result.Movies = append(result.Movies, newRelease)
			}
		}
	}

	// 2. Fetch new TV shows and seasons
	if metadataClient.SupportsTVMetadata() {
		// Get recent TV shows
		recentTVShows, err := metadataClient.GetRecentTVShows(ctx, recentReleaseDays+upcomingDays)
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch recent TV shows")
		} else {
			log.Info().Int("count", len(recentTVShows)).Msg("Found recent TV shows")

			// Convert to our NewRelease format
			for _, tvShow := range recentTVShows {
				// Parse first air date
				firstAirDate, err := time.Parse("2006-01-02", tvShow.FirstAirDate)
				if err != nil {
					// If we can't parse the date, use today
					firstAirDate = time.Now()
				}

				// Check if this is recent enough
				daysSinceRelease := time.Since(firstAirDate).Hours() / 24
				if daysSinceRelease > float64(recentReleaseDays) && daysSinceRelease < -float64(upcomingDays) {
					continue // Skip if not within our window
				}

				// Extract genres
				var genres []string
				for _, genre := range tvShow.Genres {
					genres = append(genres, genre.Name)
				}

				// Extract creators (show creators)
				var creators []string
				for _, creator := range tvShow.CreatedBy {
					creators = append(creators, creator.Name)
				}

				// Create the new release
				newRelease := NewRelease{
					ID:          fmt.Sprintf("tvshow-%s", tvShow.ID),
					ExternalID:  tvShow.ID,
					Title:       tvShow.Name,
					Description: tvShow.Overview,
					ReleaseDate: firstAirDate,
					MediaType:   "series",
					Genres:      genres,
					Creators:    creators,
					Rating:      tvShow.VoteAverage,
					Source:      "tmdb", // Assuming TMDB for now
					ImageURL:    tvShow.PosterPath,
					Metadata:    tvShow,
				}

				result.Series = append(result.Series, newRelease)
			}
		}

		// Check for trending shows as well
		trendingTVShows, err := metadataClient.GetTrendingTVShows(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch trending TV shows")
		} else {
			// Filter to only include recent trending shows
			for _, tvShow := range trendingTVShows {
				// Check if we already have this show
				alreadyAdded := false
				for _, existingShow := range result.Series {
					if existingShow.ExternalID == tvShow.ID {
						alreadyAdded = true
						break
					}
				}

				if alreadyAdded {
					continue
				}

				// Parse first air date
				firstAirDate, err := time.Parse("2006-01-02", tvShow.FirstAirDate)
				if err != nil {
					continue // Skip if no valid date
				}

				// Only include if it's a recent show (within the last month)
				daysSinceRelease := time.Since(firstAirDate).Hours() / 24
				if daysSinceRelease > 30 {
					continue
				}

				// Extract genres
				var genres []string
				for _, genre := range tvShow.Genres {
					genres = append(genres, genre.Name)
				}

				// Extract creators (show creators)
				var creators []string
				for _, creator := range tvShow.CreatedBy {
					creators = append(creators, creator.Name)
				}

				// Create the new release
				newRelease := NewRelease{
					ID:          fmt.Sprintf("tvshow-%s", tvShow.ID),
					ExternalID:  tvShow.ID,
					Title:       tvShow.Name,
					Description: tvShow.Overview,
					ReleaseDate: firstAirDate,
					MediaType:   "series",
					Genres:      genres,
					Creators:    creators,
					Rating:      tvShow.VoteAverage,
					Source:      "tmdb", // Assuming TMDB for now
					ImageURL:    tvShow.PosterPath,
					Metadata:    tvShow,
				}

				result.Series = append(result.Series, newRelease)
			}
		}
	}

	// 3. For music, we would integrate with a music metadata provider
	// Currently, we don't have a music metadata provider integrated,
	// so we'll leave this section empty for now

	log.Info().
		Int("movies", len(result.Movies)).
		Int("series", len(result.Series)).
		Int("music", len(result.Music)).
		Msg("Completed fetching new releases")

	return result, nil
}

// getMetadataClient gets a metadata client for a specific provider
func (j *NewReleaseNotificationJob) getMetadataClient(ctx context.Context, providerType types.ClientType) (interface{}, error) {
	log := utils.LoggerFromContext(ctx)

	// Here we would typically:
	// 1. Check if the requested provider type is supported
	// 2. Get its configuration from the database
	// 3. Instantiate a new client using the configuration

	// For now, we're just returning the already-injected metadata client
	if j.metadataClient == nil {
		log.Error().Str("providerType", string(providerType)).Msg("No metadata client available")
		return nil, fmt.Errorf("no metadata client available for provider type: %s", providerType)
	}

	return j.metadataClient, nil
}
