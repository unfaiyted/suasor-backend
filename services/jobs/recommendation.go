package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	mediatypes "suasor/client/media/types"
	"suasor/repository"
	"suasor/services/scheduler"
	"suasor/types/models"
	"time"
)

// RecommendationJob creates recommendations for users based on their preferences
type RecommendationJob struct {
	jobRepo         repository.JobRepository
	userRepo        repository.UserRepository
	configRepo      repository.UserConfigRepository
	movieRepo       repository.MediaItemRepository[*mediatypes.Movie]
	seriesRepo      repository.MediaItemRepository[*mediatypes.Series]
	musicRepo       repository.MediaItemRepository[*mediatypes.Track]
	historyRepo     repository.MediaPlayHistoryRepository
	aiClientService interface{} // Using interface{} to avoid import cycles
}

// NewRecommendationJob creates a new recommendation job
func NewRecommendationJob(
	jobRepo repository.JobRepository,
	userRepo repository.UserRepository,
	configRepo repository.UserConfigRepository,
	movieRepo repository.MediaItemRepository[*mediatypes.Movie],
	seriesRepo repository.MediaItemRepository[*mediatypes.Series],
	musicRepo repository.MediaItemRepository[*mediatypes.Track],
	historyRepo repository.MediaPlayHistoryRepository,
	aiClientService interface{},
) *RecommendationJob {
	return &RecommendationJob{
		jobRepo:         jobRepo,
		userRepo:        userRepo,
		configRepo:      configRepo,
		movieRepo:       movieRepo,
		seriesRepo:      seriesRepo,
		musicRepo:       musicRepo,
		historyRepo:     historyRepo,
		aiClientService: aiClientService,
	}
}

// Name returns the unique name of the job
func (j *RecommendationJob) Name() string {
	return "system.recommendation"
}

// Schedule returns when the job should next run
func (j *RecommendationJob) Schedule() time.Duration {
	// Default to checking daily
	return 24 * time.Hour
}

// Execute runs the recommendation job
func (j *RecommendationJob) Execute(ctx context.Context) error {
	log.Println("Starting recommendation job")

	// Get all users
	users, err := j.userRepo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("error getting users: %w", err)
	}

	// Process each user
	for _, user := range users {
		if err := j.processUserRecommendations(ctx, user); err != nil {
			log.Printf("Error processing recommendations for user %s: %v", user.Username, err)
			// Continue with other users even if one fails
			continue
		}
	}

	log.Println("Recommendation job completed")
	return nil
}

// processUserRecommendations generates recommendations for a single user
func (j *RecommendationJob) processUserRecommendations(ctx context.Context, user models.User) error {
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

	// Check if user has automated recommendations enabled
	if !config.RecommendationSyncEnabled {
		log.Printf("Recommendation sync not enabled for user: %s", user.Username)
		return nil
	}

	// Check if it's time to generate recommendations based on frequency
	shouldRun := j.shouldRunForUser(ctx, user.ID, config.RecommendationSyncFrequency)
	if !shouldRun {
		log.Printf("Not time to run recommendations for user: %s", user.Username)
		return nil
	}

	log.Printf("Generating recommendations for user: %s", user.Username)

	// Create a job run record for this user
	now := time.Now()
	jobRun := &models.JobRun{
		JobName:   j.Name(),
		JobType:   models.JobTypeRecommendation,
		Status:    models.JobStatusRunning,
		StartTime: &now,
		UserID:    &user.ID,
		Metadata:  fmt.Sprintf(`{"userId":%d,"username":"%s"}`, user.ID, user.Username),
	}

	if err := j.jobRepo.CreateJobRun(ctx, jobRun); err != nil {
		log.Printf("Error creating job run record: %v", err)
		return err
	}

	// Process recommendations for different content types
	var jobError error

	// Generate movie recommendations if enabled
	if j.IsContentTypeEnabled(config.RecommendationContentTypes, "movie") {
		if err := j.generateMovieRecommendations(ctx, user, config, jobRun.ID); err != nil {
			log.Printf("Error generating movie recommendations: %v", err)
			jobError = err
		}
	}

	// Generate TV show recommendations if enabled
	if j.IsContentTypeEnabled(config.RecommendationContentTypes, "series") {
		if err := j.generateSeriesRecommendations(ctx, user, config, jobRun.ID); err != nil {
			log.Printf("Error generating series recommendations: %v", err)
			if jobError == nil {
				jobError = err
			}
		}
	}

	// Generate music recommendations if enabled
	if j.IsContentTypeEnabled(config.RecommendationContentTypes, "music") {
		if err := j.generateMusicRecommendations(ctx, user, config, jobRun.ID); err != nil {
			log.Printf("Error generating music recommendations: %v", err)
			if jobError == nil {
				jobError = err
			}
		}
	}

	// Set job status based on outcome
	status := models.JobStatusCompleted
	errorMessage := ""
	if jobError != nil {
		status = models.JobStatusFailed
		errorMessage = jobError.Error()
	}

	if err := j.jobRepo.CompleteJobRun(ctx, jobRun.ID, status, errorMessage); err != nil {
		log.Printf("Error completing job run: %v", err)
	}

	// Update the job schedule's last run time
	jobName := fmt.Sprintf("%s.user.%d", j.Name(), user.ID)
	if err := j.jobRepo.UpdateJobLastRunTime(ctx, jobName, now); err != nil {
		log.Printf("Error updating job last run time: %v", err)
	}

	return jobError
}

// shouldRunForUser determines if recommendations should be generated for a user
func (j *RecommendationJob) shouldRunForUser(ctx context.Context, userID uint64, frequency string) bool {
	// Convert to scheduler.Frequency
	freq := scheduler.Frequency(frequency)

	// Manual frequency means never auto-run
	if freq == scheduler.FrequencyManual {
		return false
	}

	// Get the last run time for this user
	jobName := fmt.Sprintf("%s.user.%d", j.Name(), userID)
	schedule, err := j.jobRepo.GetJobSchedule(ctx, jobName)
	if err != nil {
		log.Printf("Error getting job schedule for user %d: %v", userID, err)
		// If we can't get the schedule, assume it should run
		return true
	}

	// If no schedule exists or it has never run, it should run
	if schedule == nil || schedule.LastRunTime == nil {
		return true
	}

	// Check if enough time has passed since the last run
	return freq.ShouldRunNow(*schedule.LastRunTime)
}

// isContentTypeEnabled checks if a content type is enabled in the comma-separated list
// IsContentTypeEnabled checks if a content type is enabled in the comma-separated list
func (j *RecommendationJob) IsContentTypeEnabled(contentTypes string, contentType string) bool {
	if contentTypes == "" {
		// If no content types are specified, assume all are enabled
		return true
	}
	
	// Split the comma-separated list
	types := strings.Split(contentTypes, ",")
	
	// Check if the content type is in the list
	for _, t := range types {
		if strings.TrimSpace(t) == contentType {
			return true
		}
	}
	
	return false
}

// generateMovieRecommendations creates movie recommendations for a user
func (j *RecommendationJob) generateMovieRecommendations(ctx context.Context, user models.User, config *models.UserConfig, jobRunID uint64) error {
	log.Printf("Generating movie recommendations for user %s", user.Username)

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 10, "Analyzing user movie preferences")

	// Get recent user movie history
	recentMovies, err := j.historyRepo.GetRecentUserMovieHistory(ctx, user.ID, 20)
	if err != nil {
		log.Printf("Error retrieving recent movie history: %v", err)
		return fmt.Errorf("error retrieving movie history: %w", err)
	}

	// Get user's movies (for determining if recommended movies are already in library)
	userMovies, err := j.movieRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		log.Printf("Error retrieving user movies: %v", err)
		return fmt.Errorf("error retrieving user movies: %w", err)
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 20, "Building user preference profile")

	// Build user preference profile based on watch history
	preferenceProfile := j.BuildUserMoviePreferences(recentMovies, userMovies, config)
	
	// Create a map of existing movies in library
	inLibraryMap := make(map[string]bool)
	for _, movie := range userMovies {
		if movie.Data != nil && movie.Data.Details.Title != "" {
			key := fmt.Sprintf("%s-%d", movie.Data.Details.Title, movie.Data.Details.ReleaseYear)
			inLibraryMap[key] = true
		}
	}

	// Generate recommendation strategies based on user profile
	var recommendations []*models.Recommendation

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 30, "Generating recommendations based on genres")

	// 1. Based on preferred genres
	genreBasedRecs := j.generateGenreBasedRecommendations(ctx, user.ID, preferenceProfile, inLibraryMap)
	recommendations = append(recommendations, genreBasedRecs...)

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 50, "Generating recommendations based on similar content")

	// 2. Based on similar content to what they've watched
	similarContentRecs := j.generateSimilarContentRecommendations(ctx, user.ID, recentMovies, inLibraryMap)
	recommendations = append(recommendations, similarContentRecs...)

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 70, "Generating recommendations based on popularity")

	// 3. Popular content they haven't seen
	popularRecs := j.generatePopularRecommendations(ctx, user.ID, inLibraryMap)
	recommendations = append(recommendations, popularRecs...)

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 80, "Filtering and ranking recommendations")

	// Filter out any duplicates and limit total recommendations
	finalRecs := j.FilterAndRankRecommendations(recommendations, 15)

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 90, "Saving recommendations")

	// Set the job run ID for all recommendations
	for _, rec := range finalRecs {
		rec.JobRunID = &jobRunID
	}

	// Save all recommendations in batch
	if len(finalRecs) > 0 {
		if err := j.jobRepo.BatchCreateRecommendations(ctx, finalRecs); err != nil {
			log.Printf("Error creating batch recommendations: %v", err)
			return fmt.Errorf("error saving recommendations: %w", err)
		}
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 100, fmt.Sprintf("Generated %d movie recommendations", len(finalRecs)))

	return nil
}

// generateSystemMovieRecommendations creates movie recommendations using system algorithms
func (j *RecommendationJob) generateSystemMovieRecommendations(ctx context.Context, user models.User, config *models.UserConfig, jobRunID uint64) error {
	log.Printf("Generating system movie recommendations for user %s", user.Username)
	return j.generateMovieRecommendations(ctx, user, config, jobRunID)
}

// generateSeriesRecommendations creates TV series recommendations for a user
func (j *RecommendationJob) generateSeriesRecommendations(ctx context.Context, user models.User, config *models.UserConfig, jobRunID uint64) error {
	log.Printf("Generating series recommendations for user %s", user.Username)
	
	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 10, "Analyzing user TV show preferences")
	
	// Get recent user series watch history
	_, err := j.historyRepo.GetRecentUserSeriesHistory(ctx, user.ID, 20)
	if err != nil {
		log.Printf("Error retrieving recent series history: %v", err)
		return fmt.Errorf("error retrieving series history: %w", err)
	}
	
	// Get user's series (for determining if recommended series are already in library)
	userSeries, err := j.seriesRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		log.Printf("Error retrieving user series: %v", err)
		return fmt.Errorf("error retrieving user series: %w", err)
	}
	
	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 25, "Building user TV preference profile")
	
	// Create a map of existing series in library
	inLibraryMap := make(map[string]bool)
	for _, series := range userSeries {
		if series.Data != nil && series.Data.Details.Title != "" {
			key := fmt.Sprintf("%s", series.Data.Details.Title)
			inLibraryMap[key] = true
		}
	}
	
	// Generate sample series recommendations based on genres
	// In a real implementation, this would be more sophisticated and leverage
	// the user's watch history, preferences, and possibly external recommendation APIs
	seriesRecommendations := []struct {
		Title       string
		Year        int
		Genre       string
		Rating      float32
		Confidence  float32
		Reason      string
	}{
		{
			Title:      "Stranger Things",
			Year:       2016,
			Genre:      "Sci-Fi",
			Rating:     8.7,
			Confidence: 0.85,
			Reason:     "Popular sci-fi series with high ratings",
		},
		{
			Title:      "The Crown",
			Year:       2016,
			Genre:      "Drama",
			Rating:     8.7,
			Confidence: 0.8,
			Reason:     "Critically acclaimed historical drama",
		},
		{
			Title:      "Ted Lasso",
			Year:       2020,
			Genre:      "Comedy",
			Rating:     8.8,
			Confidence: 0.9,
			Reason:     "Award-winning comedy series",
		},
		{
			Title:      "The Last of Us",
			Year:       2023,
			Genre:      "Drama",
			Rating:     8.8,
			Confidence: 0.85,
			Reason:     "Popular adaptation of a video game",
		},
		{
			Title:      "The Mandalorian",
			Year:       2019,
			Genre:      "Sci-Fi",
			Rating:     8.7,
			Confidence: 0.8,
			Reason:     "Star Wars universe series with wide appeal",
		},
	}
	
	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 50, "Creating series recommendations")
	
	// Process all recommended series
	var recommendations []*models.Recommendation
	
	for _, rec := range seriesRecommendations {
		// Check if already in library
		isInLibrary := inLibraryMap[rec.Title]
		
		// Create media item if not exists
		mediaItem, err := j.findOrCreateSeriesItem(ctx, rec.Title, rec.Year, rec.Genre, user.ID)
		if err != nil {
			log.Printf("Error creating series item: %v", err)
			continue
		}
		
		// Skip if user has already viewed this series
		viewed, _ := j.historyRepo.HasUserViewedMedia(ctx, user.ID, mediaItem.ID)
		if viewed {
			continue
		}
		
		// Check if we already have this recommendation
		existingRec, _ := j.jobRepo.GetRecommendationByMediaItem(ctx, user.ID, mediaItem.ID)
		if existingRec != nil && existingRec.Active {
			continue
		}
		
		// Create recommendation
		recommendation := &models.Recommendation{
			UserID:      user.ID,
			MediaItemID: mediaItem.ID,
			MediaType:   string(mediatypes.MediaTypeSeries),
			Source:      models.RecommendationSourceSystem,
			Reason:      rec.Reason,
			Confidence:  rec.Confidence,
			InLibrary:   isInLibrary,
			Viewed:      viewed,
			Active:      true,
			JobRunID:    &jobRunID,
			Metadata:    makeMetadataJson(map[string]interface{}{"genre": rec.Genre, "rating": rec.Rating}),
		}
		
		recommendations = append(recommendations, recommendation)
	}
	
	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 75, "Saving series recommendations")
	
	// Save all recommendations in batch
	if len(recommendations) > 0 {
		if err := j.jobRepo.BatchCreateRecommendations(ctx, recommendations); err != nil {
			log.Printf("Error creating batch recommendations: %v", err)
			return fmt.Errorf("error saving series recommendations: %w", err)
		}
	}
	
	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 100, fmt.Sprintf("Generated %d series recommendations", len(recommendations)))
	
	return nil
}

// findOrCreateSeriesItem finds a series in the database or creates it if it doesn't exist
func (j *RecommendationJob) findOrCreateSeriesItem(ctx context.Context, title string, year int, genre string, userID uint64) (*models.MediaItem[*mediatypes.Series], error) {
	// Create a series with the correct structure
	details := mediatypes.MediaDetails{
		Title:       title,
		ReleaseYear: year,
		Description: "Generated series recommendation",
	}
	
	series := &mediatypes.Series{
		Details: details,
		Genres:  []string{genre},
	}
	
	// Try to find an existing series with same title
	userSeries, err := j.seriesRepo.GetByUserID(ctx, userID)
	if err == nil {
		for _, existing := range userSeries {
			if existing.Data != nil && existing.Data.Details.Title == title {
				return existing, nil
			}
		}
	}
	
	// If not found, create a new placeholder item
	mediaItem := models.MediaItem[*mediatypes.Series]{
		Type: mediatypes.MediaTypeSeries,
		Data: series,
		ClientIDs: []models.ClientID{
			{
				ID:     0, // Special ID for recommendation engine
				Type:   "recommendation",
				ItemID: fmt.Sprintf("recommendation-%s", title),
			},
		},
		ExternalIDs: []models.ExternalID{},
	}
	
	// For a placeholder item, use a mock ID
	mediaItem.ID = uint64(time.Now().UnixNano() % 100000000)
	return &mediaItem, nil
}

// generateMusicRecommendations creates music recommendations for a user
func (j *RecommendationJob) generateMusicRecommendations(ctx context.Context, user models.User, config *models.UserConfig, jobRunID uint64) error {
	log.Printf("Generating music recommendations for user %s", user.Username)
	
	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 10, "Analyzing user music preferences")
	
	// Get recent user music play history
	_, err := j.historyRepo.GetRecentUserMusicHistory(ctx, user.ID, 20)
	if err != nil {
		log.Printf("Error retrieving recent music history: %v", err)
		return fmt.Errorf("error retrieving music history: %w", err)
	}
	
	// Get user's music (for determining if recommended tracks are already in library)
	userMusic, err := j.musicRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		log.Printf("Error retrieving user music: %v", err)
		return fmt.Errorf("error retrieving user music: %w", err)
	}
	
	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 25, "Building user music preference profile")
	
	// Create a map of existing tracks in library
	inLibraryMap := make(map[string]bool)
	for _, track := range userMusic {
		if track.Data != nil && track.Data.Details.Title != "" {
			key := fmt.Sprintf("%s-%s", track.Data.ArtistName, track.Data.Details.Title)
			inLibraryMap[key] = true
		}
	}
	
	// Generate sample music recommendations
	// In a real implementation, this would use the user's listening history,
	// genre preferences, and other factors to generate personalized recommendations
	musicRecommendations := []struct {
		Title      string
		Artist     string
		Album      string
		Year       int
		Genre      string
		Confidence float32
		Reason     string
	}{
		{
			Title:      "Dreams",
			Artist:     "Fleetwood Mac",
			Album:      "Rumours",
			Year:       1977,
			Genre:      "Rock",
			Confidence: 0.85,
			Reason:     "Classic rock track with wide appeal",
		},
		{
			Title:      "Blinding Lights",
			Artist:     "The Weeknd",
			Album:      "After Hours",
			Year:       2020,
			Genre:      "Pop",
			Confidence: 0.9,
			Reason:     "Popular modern pop hit",
		},
		{
			Title:      "Redbone",
			Artist:     "Childish Gambino",
			Album:      "Awaken, My Love!",
			Year:       2016,
			Genre:      "R&B",
			Confidence: 0.8,
			Reason:     "Acclaimed R&B track with crossover appeal",
		},
		{
			Title:      "Come As You Are",
			Artist:     "Nirvana",
			Album:      "Nevermind",
			Year:       1991,
			Genre:      "Alternative",
			Confidence: 0.75,
			Reason:     "Influential alternative track",
		},
		{
			Title:      "Alright",
			Artist:     "Kendrick Lamar",
			Album:      "To Pimp a Butterfly",
			Year:       2015,
			Genre:      "Hip-Hop",
			Confidence: 0.85,
			Reason:     "Critically acclaimed hip-hop track",
		},
	}
	
	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 50, "Creating music recommendations")
	
	// Process all recommended tracks
	var recommendations []*models.Recommendation
	
	for _, rec := range musicRecommendations {
		// Check if already in library
		key := fmt.Sprintf("%s-%s", rec.Artist, rec.Title)
		isInLibrary := inLibraryMap[key]
		
		// Create media item if not exists
		mediaItem, err := j.findOrCreateMusicItem(ctx, rec.Title, rec.Artist, rec.Album, rec.Year, rec.Genre, user.ID)
		if err != nil {
			log.Printf("Error creating music item: %v", err)
			continue
		}
		
		// Skip if user has already played this track
		viewed, _ := j.historyRepo.HasUserViewedMedia(ctx, user.ID, mediaItem.ID)
		if viewed {
			continue
		}
		
		// Check if we already have this recommendation
		existingRec, _ := j.jobRepo.GetRecommendationByMediaItem(ctx, user.ID, mediaItem.ID)
		if existingRec != nil && existingRec.Active {
			continue
		}
		
		// Create recommendation
		recommendation := &models.Recommendation{
			UserID:      user.ID,
			MediaItemID: mediaItem.ID,
			MediaType:   string(mediatypes.MediaTypeTrack),
			Source:      models.RecommendationSourceSystem,
			Reason:      rec.Reason,
			Confidence:  rec.Confidence,
			InLibrary:   isInLibrary,
			Viewed:      viewed,
			Active:      true,
			JobRunID:    &jobRunID,
			Metadata:    makeMetadataJson(map[string]interface{}{"genre": rec.Genre, "year": rec.Year}),
		}
		
		recommendations = append(recommendations, recommendation)
	}
	
	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 75, "Saving music recommendations")
	
	// Save all recommendations in batch
	if len(recommendations) > 0 {
		if err := j.jobRepo.BatchCreateRecommendations(ctx, recommendations); err != nil {
			log.Printf("Error creating batch recommendations: %v", err)
			return fmt.Errorf("error saving music recommendations: %w", err)
		}
	}
	
	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 100, fmt.Sprintf("Generated %d music recommendations", len(recommendations)))
	
	return nil
}

// findOrCreateMusicItem finds a music track in the database or creates it if it doesn't exist
func (j *RecommendationJob) findOrCreateMusicItem(ctx context.Context, title, artist, album string, year int, genre string, userID uint64) (*models.MediaItem[*mediatypes.Track], error) {
	// Create a track with the correct structure
	details := mediatypes.MediaDetails{
		Title:       title,
		ReleaseYear: year,
		Genres:      []string{genre},
	}
	
	track := &mediatypes.Track{
		Details:    details,
		ArtistName: artist,
		AlbumName:  album,
	}
	
	// Try to find an existing track with same title and artist
	userTracks, err := j.musicRepo.GetByUserID(ctx, userID)
	if err == nil {
		for _, existing := range userTracks {
			if existing.Data != nil && 
			   existing.Data.Details.Title == title && 
			   existing.Data.ArtistName == artist {
				return existing, nil
			}
		}
	}
	
	// If not found, create a new placeholder item
	mediaItem := models.MediaItem[*mediatypes.Track]{
		Type:        mediatypes.MediaTypeTrack,
		Data:        track,
		ClientIDs:   []models.ClientID{},
		ExternalIDs: []models.ExternalID{},
	}
	
	// Add client ID for recommendation engine
	mediaItem.AddClientID(0, "recommendation", fmt.Sprintf("recommendation-%s-%s", artist, title))
	
	// Add external ID for recommendation
	mediaItem.AddExternalID("recommendation", fmt.Sprintf("recommendation-%s-%s", artist, title))
	
	// For a placeholder item, use a mock ID
	mediaItem.ID = uint64(time.Now().UnixNano() % 100000000)
	return &mediaItem, nil
}

// UserPreferenceProfile represents a user's media preferences
type UserPreferenceProfile struct {
	// Movie preferences
	FavoriteGenres     map[string]float32 // Genre -> weight
	FavoriteActors     map[string]float32 // Actor -> weight
	FavoriteDirectors  map[string]float32 // Director -> weight
	WatchedMovies      map[uint64]bool    // MediaItemID -> watched
	ExcludedGenres     []string           // Genres to exclude
	PreferredReleaseYears [2]int          // [min, max] years
	MinRating         float32             // Minimum rating (0-10)
	PreferredLanguages map[string]float32 // Language -> weight
}

// BuildUserMoviePreferences builds a user preference profile based on watch history and config
func (j *RecommendationJob) BuildUserMoviePreferences(
	recentMovies []models.MediaPlayHistory[*mediatypes.Movie],
	userMovies []*models.MediaItem[*mediatypes.Movie],
	config *models.UserConfig) *UserPreferenceProfile {
	
	profile := &UserPreferenceProfile{
		FavoriteGenres:     make(map[string]float32),
		FavoriteActors:     make(map[string]float32),
		FavoriteDirectors:  make(map[string]float32),
		WatchedMovies:      make(map[uint64]bool),
		PreferredLanguages: make(map[string]float32),
		MinRating:          5.0, // Default minimum rating
	}
	
	// Set default release year range to past 15 years
	currentYear := time.Now().Year()
	profile.PreferredReleaseYears = [2]int{currentYear - 15, currentYear}
	
	// Add excluded genres from user preferences if available
	if config.ExcludedGenres != nil && len(config.ExcludedGenres.Movies) > 0 {
		profile.ExcludedGenres = config.ExcludedGenres.Movies
	}
	
	// User may have preferred content rating limits
	if config.MinContentRating != "" || config.MaxContentRating != "" {
		// Logic for content rating would depend on your rating system
	}
	
	// Handle preferred age of content
	if config.RecommendationMaxAge > 0 {
		minYear := currentYear - config.RecommendationMaxAge
		profile.PreferredReleaseYears[0] = minYear
	}
	
	// Process watch history to build preference profile
	for _, history := range recentMovies {
		if history.Item == nil || history.Item.Data == nil {
			continue
		}
		
		movie := history.Item.Data
		
		// Mark movie as watched
		profile.WatchedMovies[history.MediaItemID] = true
		
		// Weight based on how recently and how many times played
		weight := float32(1.0)
		if history.PlayCount > 0 {
			weight += float32(history.PlayCount) * 0.5 // More plays = higher weight
		}
		
		// More recent watches have higher weight
		if !history.LastWatchedAt.IsZero() {
			daysAgo := time.Since(history.LastWatchedAt).Hours() / 24
			if daysAgo < 7 {
				weight += 1.0 // Watched within last week
			} else if daysAgo < 30 {
				weight += 0.5 // Watched within last month
			}
		}
		
		// Process genres
		for _, genre := range movie.Details.Genres {
			profile.FavoriteGenres[genre] += weight
		}
		
		// Process cast
		for _, person := range movie.Cast {
			if person.Role == "actor" {
				profile.FavoriteActors[person.Name] += weight
			} else if person.Role == "director" {
				profile.FavoriteDirectors[person.Name] += weight
			}
		}
		
		// We could process language if it was available in the Details struct
		// For now, this is a placeholder for future language preference tracking
	}
	
	// Add user's preferred genres if set in their config
	if config.PreferredGenres != nil && len(config.PreferredGenres.Movies) > 0 {
		for _, genre := range config.PreferredGenres.Movies {
			// Give extra weight to explicitly preferred genres
			profile.FavoriteGenres[genre] += 3.0
		}
	}
	
	return profile
}

// generateGenreBasedRecommendations creates recommendations based on user's genre preferences
func (j *RecommendationJob) generateGenreBasedRecommendations(
	ctx context.Context, 
	userID uint64, 
	profile *UserPreferenceProfile,
	inLibraryMap map[string]bool) []*models.Recommendation {
	
	recommendations := []*models.Recommendation{}
	
	// Get top 3 genres
	type genreWeight struct {
		genre  string
		weight float32
	}
	
	var topGenres []genreWeight
	for genre, weight := range profile.FavoriteGenres {
		topGenres = append(topGenres, genreWeight{genre, weight})
	}
	
	// Sort by weight descending
	sort.Slice(topGenres, func(i, j int) bool {
		return topGenres[i].weight > topGenres[j].weight
	})
	
	// Limit to top 3 genres
	genreCount := 3
	if len(topGenres) < genreCount {
		genreCount = len(topGenres)
	}
	
	// For each top genre, find movies that match
	for i := 0; i < genreCount; i++ {
		if i >= len(topGenres) {
			break
		}
		
		genre := topGenres[i].genre
		
		// This is a placeholder for querying a movie database or recommendation service
		// In a real implementation, you would query your database for movies matching this genre
		// that the user hasn't watched yet
		movies := j.getSystemRecommendedMoviesByGenre(ctx, genre, profile)
		
		for _, movie := range movies {
			// Check if already in library
			key := fmt.Sprintf("%s-%d", movie.Title, movie.Year)
			isInLibrary := inLibraryMap[key]
			
			// Create media item if not exists (simplified here)
			mediaItem, err := j.findOrCreateMovieItem(ctx, movie.Title, movie.Year, userID)
			if err != nil {
				log.Printf("Error creating movie item: %v", err)
				continue
			}
			
			// Skip if user has already viewed this movie
			viewed, _ := j.historyRepo.HasUserViewedMedia(ctx, userID, mediaItem.ID)
			if viewed {
				continue
			}
			
			// Check if we already have this recommendation for this user
			existingRec, _ := j.jobRepo.GetRecommendationByMediaItem(ctx, userID, mediaItem.ID)
			if existingRec != nil && existingRec.Active {
				// Skip if already recommended and active
				continue
			}
			
			// Calculate confidence score based on genre weight
			confidence := topGenres[i].weight / 10.0 // Normalize to 0-1 range
			if confidence > 1.0 {
				confidence = 1.0
			}
			
			// Create recommendation
			recommendation := &models.Recommendation{
				UserID:      userID,
				MediaItemID: mediaItem.ID,
				MediaType:   string(mediatypes.MediaTypeMovie),
				Source:      models.RecommendationSourceSystem,
				Reason:      fmt.Sprintf("Based on your interest in %s movies", genre),
				Confidence:  confidence,
				InLibrary:   isInLibrary,
				Viewed:      viewed,
				Active:      true,
				Metadata:    makeMetadataJson(map[string]interface{}{"matchedGenre": genre}),
			}
			
			recommendations = append(recommendations, recommendation)
		}
	}
	
	return recommendations
}

// generateSimilarContentRecommendations creates recommendations based on similar content to what user has watched
func (j *RecommendationJob) generateSimilarContentRecommendations(
	ctx context.Context, 
	userID uint64, 
	recentMovies []models.MediaPlayHistory[*mediatypes.Movie],
	inLibraryMap map[string]bool) []*models.Recommendation {
	
	recommendations := []*models.Recommendation{}
	
	// Take up to 5 most recent movies to find similar content
	recentCount := 5
	if len(recentMovies) < recentCount {
		recentCount = len(recentMovies)
	}
	
	for i := 0; i < recentCount; i++ {
		if i >= len(recentMovies) || recentMovies[i].Item == nil || recentMovies[i].Item.Data == nil {
			continue
		}
		
		recentMovie := recentMovies[i].Item.Data
		
		// In a real implementation, you would call a recommendation service or database
		// to find similar movies based on this recent watch
		similarMovies := j.getSystemSimilarMovies(ctx, recentMovie)
		
		for _, movie := range similarMovies {
			// Check if already in library
			key := fmt.Sprintf("%s-%d", movie.Title, movie.Year)
			isInLibrary := inLibraryMap[key]
			
			// Create media item if not exists
			mediaItem, err := j.findOrCreateMovieItem(ctx, movie.Title, movie.Year, userID)
			if err != nil {
				log.Printf("Error creating movie item: %v", err)
				continue
			}
			
			// Skip if user has already viewed this movie
			viewed, _ := j.historyRepo.HasUserViewedMedia(ctx, userID, mediaItem.ID)
			if viewed {
				continue
			}
			
			// Check if we already have this recommendation
			existingRec, _ := j.jobRepo.GetRecommendationByMediaItem(ctx, userID, mediaItem.ID)
			if existingRec != nil && existingRec.Active {
				continue
			}
			
			// Create recommendation
			recommendation := &models.Recommendation{
				UserID:      userID,
				MediaItemID: mediaItem.ID,
				MediaType:   string(mediatypes.MediaTypeMovie),
				Source:      models.RecommendationSourceSystem,
				Reason:      fmt.Sprintf("Similar to %s, which you watched recently", recentMovie.Details.Title),
				Confidence:  0.8, // Higher confidence for similar content
				InLibrary:   isInLibrary,
				Viewed:      viewed,
				Active:      true,
				Metadata:    makeMetadataJson(map[string]interface{}{"similarTo": recentMovie.Details.Title}),
			}
			
			recommendations = append(recommendations, recommendation)
		}
	}
	
	return recommendations
}

// generatePopularRecommendations creates recommendations based on popular content
func (j *RecommendationJob) generatePopularRecommendations(
	ctx context.Context, 
	userID uint64,
	inLibraryMap map[string]bool) []*models.Recommendation {
	
	recommendations := []*models.Recommendation{}
	
	// Get popular movies the user hasn't seen yet
	popularMovies := j.getSystemPopularMovies(ctx)
	
	for _, movie := range popularMovies {
		// Check if already in library
		key := fmt.Sprintf("%s-%d", movie.Title, movie.Year)
		isInLibrary := inLibraryMap[key]
		
		// Create media item if not exists
		mediaItem, err := j.findOrCreateMovieItem(ctx, movie.Title, movie.Year, userID)
		if err != nil {
			log.Printf("Error creating movie item: %v", err)
			continue
		}
		
		// Skip if user has already viewed this movie
		viewed, _ := j.historyRepo.HasUserViewedMedia(ctx, userID, mediaItem.ID)
		if viewed {
			continue
		}
		
		// Check if we already have this recommendation
		existingRec, _ := j.jobRepo.GetRecommendationByMediaItem(ctx, userID, mediaItem.ID)
		if existingRec != nil && existingRec.Active {
			continue
		}
		
		// Create recommendation
		recommendation := &models.Recommendation{
			UserID:      userID,
			MediaItemID: mediaItem.ID,
			MediaType:   string(mediatypes.MediaTypeMovie),
			Source:      models.RecommendationSourceSystem,
			Reason:      "Popular movie you might enjoy",
			Confidence:  0.7, // Lower confidence for general popularity
			InLibrary:   isInLibrary,
			Viewed:      viewed,
			Active:      true,
			Metadata:    makeMetadataJson(map[string]interface{}{"popularityRank": movie.PopularityRank}),
		}
		
		recommendations = append(recommendations, recommendation)
	}
	
	return recommendations
}

// FilterAndRankRecommendations removes duplicates and ranks recommendations
func (j *RecommendationJob) FilterAndRankRecommendations(recommendations []*models.Recommendation, limit int) []*models.Recommendation {
	// Create a map to deduplicate by media item ID
	uniqueRecs := make(map[uint64]*models.Recommendation)
	
	for _, rec := range recommendations {
		// If we already have this recommendation, keep the one with higher confidence
		if existing, found := uniqueRecs[rec.MediaItemID]; found {
			if rec.Confidence > existing.Confidence {
				uniqueRecs[rec.MediaItemID] = rec
			}
		} else {
			uniqueRecs[rec.MediaItemID] = rec
		}
	}
	
	// Convert map to slice
	var filteredRecs []*models.Recommendation
	for _, rec := range uniqueRecs {
		filteredRecs = append(filteredRecs, rec)
	}
	
	// Sort by confidence descending
	sort.Slice(filteredRecs, func(i, j int) bool {
		return filteredRecs[i].Confidence > filteredRecs[j].Confidence
	})
	
	// Limit to requested number
	if len(filteredRecs) > limit {
		filteredRecs = filteredRecs[:limit]
	}
	
	return filteredRecs
}

// MovieRecommendation represents a movie recommendation from the system
type MovieRecommendation struct {
	Title          string
	Year           int
	PopularityRank int
	Rating         float32
	Genres         []string
}

// getSystemRecommendedMoviesByGenre gets system recommended movies by genre
// This is a placeholder that would be replaced with actual database queries
func (j *RecommendationJob) getSystemRecommendedMoviesByGenre(ctx context.Context, genre string, profile *UserPreferenceProfile) []MovieRecommendation {
	// In a real implementation, this would query your database or external API
	// For now, return some placeholders based on genre
	movies := []MovieRecommendation{}
	
	// Sample movies based on genre
	switch genre {
	case "Action":
		movies = append(movies, 
			MovieRecommendation{Title: "John Wick", Year: 2014, Rating: 7.4, Genres: []string{"Action", "Thriller"}},
			MovieRecommendation{Title: "Die Hard", Year: 1988, Rating: 8.2, Genres: []string{"Action", "Thriller"}},
			MovieRecommendation{Title: "Mad Max: Fury Road", Year: 2015, Rating: 8.1, Genres: []string{"Action", "Adventure"}})
	case "Drama":
		movies = append(movies, 
			MovieRecommendation{Title: "The Shawshank Redemption", Year: 1994, Rating: 9.3, Genres: []string{"Drama"}},
			MovieRecommendation{Title: "Forrest Gump", Year: 1994, Rating: 8.8, Genres: []string{"Drama", "Romance"}},
			MovieRecommendation{Title: "The Godfather", Year: 1972, Rating: 9.2, Genres: []string{"Drama", "Crime"}})
	case "Sci-Fi":
		movies = append(movies, 
			MovieRecommendation{Title: "Blade Runner 2049", Year: 2017, Rating: 8.0, Genres: []string{"Sci-Fi", "Drama"}},
			MovieRecommendation{Title: "Arrival", Year: 2016, Rating: 7.9, Genres: []string{"Sci-Fi", "Drama"}},
			MovieRecommendation{Title: "Dune", Year: 2021, Rating: 8.0, Genres: []string{"Sci-Fi", "Adventure"}})
	case "Comedy":
		movies = append(movies, 
			MovieRecommendation{Title: "Superbad", Year: 2007, Rating: 7.6, Genres: []string{"Comedy"}},
			MovieRecommendation{Title: "The Grand Budapest Hotel", Year: 2014, Rating: 8.1, Genres: []string{"Comedy", "Drama"}},
			MovieRecommendation{Title: "Bridesmaids", Year: 2011, Rating: 6.8, Genres: []string{"Comedy", "Romance"}})
	default:
		// Default recommendations if genre doesn't match
		movies = append(movies, 
			MovieRecommendation{Title: "The Matrix", Year: 1999, Rating: 8.7, Genres: []string{"Action", "Sci-Fi"}},
			MovieRecommendation{Title: "Inception", Year: 2010, Rating: 8.8, Genres: []string{"Sci-Fi", "Action"}},
			MovieRecommendation{Title: "Interstellar", Year: 2014, Rating: 8.6, Genres: []string{"Sci-Fi", "Drama"}})
	}
	
	return movies
}

// getSystemSimilarMovies gets movies similar to the provided movie
// This is a placeholder that would be replaced with actual similarity logic
func (j *RecommendationJob) getSystemSimilarMovies(ctx context.Context, movie *mediatypes.Movie) []MovieRecommendation {
	// In a real implementation, this would use collaborative filtering, content-based filtering,
	// or other recommendation techniques to find similar movies
	
	// For now, return some placeholders based on the movie's genre
	// Simple placeholder logic based on movie title
	if len(movie.Details.Genres) > 0 {
		return j.getSystemRecommendedMoviesByGenre(ctx, movie.Details.Genres[0], nil)
	}
	
	// Default recommendations
	return []MovieRecommendation{
		{Title: "The Dark Knight", Year: 2008, Rating: 9.0, Genres: []string{"Action", "Crime", "Drama"}},
		{Title: "Pulp Fiction", Year: 1994, Rating: 8.9, Genres: []string{"Crime", "Drama"}},
		{Title: "Fight Club", Year: 1999, Rating: 8.8, Genres: []string{"Drama"}},
	}
}

// getSystemPopularMovies gets popular movies from the system
// This is a placeholder that would be replaced with actual database queries
func (j *RecommendationJob) getSystemPopularMovies(ctx context.Context) []MovieRecommendation {
	// In a real implementation, this would query your database or external API for popular movies
	return []MovieRecommendation{
		{Title: "Oppenheimer", Year: 2023, PopularityRank: 1, Rating: 8.5, Genres: []string{"Biography", "Drama", "History"}},
		{Title: "Barbie", Year: 2023, PopularityRank: 2, Rating: 7.0, Genres: []string{"Adventure", "Comedy", "Fantasy"}},
		{Title: "Mission: Impossible - Dead Reckoning", Year: 2023, PopularityRank: 3, Rating: 7.8, Genres: []string{"Action", "Adventure", "Thriller"}},
		{Title: "Guardians of the Galaxy Vol. 3", Year: 2023, PopularityRank: 4, Rating: 8.0, Genres: []string{"Action", "Adventure", "Comedy"}},
		{Title: "Spider-Man: Across the Spider-Verse", Year: 2023, PopularityRank: 5, Rating: 8.7, Genres: []string{"Animation", "Action", "Adventure"}},
	}
}

// findOrCreateMovieItem finds a movie in the database or creates it if it doesn't exist
func (j *RecommendationJob) findOrCreateMovieItem(ctx context.Context, title string, year int, userID uint64) (*models.MediaItem[*mediatypes.Movie], error) {
	// Create a placeholder movie with the correct structure
	details := mediatypes.MediaDetails{
		Title:       title,
		ReleaseYear: year,
		Description: "Generated recommendation",
	}

	details.Genres = []string{"Unknown"} // Default genre
	
	movie := &mediatypes.Movie{
		Details: details,
	}

	// Try to find an existing movie with same title and year
	// This is a simplified approach - in a real system, you might use
	// a more sophisticated matching algorithm or external IDs
	
	// Search in user's media items first
	userMovies, err := j.movieRepo.GetByUserID(ctx, userID)
	if err == nil {
		for _, existing := range userMovies {
			if existing.Data != nil && 
			   existing.Data.Details.Title == title && 
			   existing.Data.Details.ReleaseYear == year {
				return existing, nil
			}
		}
	}
	
	// If not found, create a new placeholder item
	mediaItem := models.MediaItem[*mediatypes.Movie]{
		Type:        mediatypes.MediaTypeMovie,
		Data:        movie,
		ClientIDs:   []models.ClientID{},
		ExternalIDs: []models.ExternalID{},
	}
	
	// Add client ID for recommendation engine
	mediaItem.AddClientID(0, "recommendation", fmt.Sprintf("recommendation-%s-%d", title, year))
	
	// Add external ID for recommendation
	mediaItem.AddExternalID("recommendation", fmt.Sprintf("recommendation-%s-%d", title, year))
	
	// In a real implementation, you would:
	// 1. Create the movie in the database
	// 2. Return the created item with its new ID
	
	// For now, use a placeholder ID
	mediaItem.ID = uint64(time.Now().UnixNano() % 100000000) // Mock ID generation
	return &mediaItem, nil
}

// hasUserViewed checks if a user has already viewed/played a media item
func (j *RecommendationJob) hasUserViewed(ctx context.Context, userID, mediaItemID uint64) (bool, error) {
	return j.historyRepo.HasUserViewedMedia(ctx, userID, mediaItemID)
}

// isInUserLibrary checks if a media item is in the user's library
func (j *RecommendationJob) isInUserLibrary(ctx context.Context, userID, mediaItemID uint64, mediaType string) (bool, error) {
	// This depends on your implementation, but generally you would:
	// 1. Get all media items for the user's clients
	// 2. Check if the media item is in the collection
	
	if mediaType == string(mediatypes.MediaTypeMovie) {
		userMovies, err := j.movieRepo.GetByUserID(ctx, userID)
		if err != nil {
			return false, err
		}
		
		for _, movie := range userMovies {
			if movie.ID == mediaItemID {
				return true, nil
			}
		}
	}
	
	return false, nil
}

// Helper to build filters for AI movie recommendations
func buildMovieFilters(config *models.UserConfig, recentMovies []map[string]interface{}) map[string]interface{} {
	filters := map[string]interface{}{}

	// Add preferred genres if available
	if config.PreferredGenres != nil && len(config.PreferredGenres.Movies) > 0 {
		filters["preferredGenres"] = config.PreferredGenres.Movies
	}

	// Add excluded genres if available
	if config.ExcludedGenres != nil && len(config.ExcludedGenres.Movies) > 0 {
		filters["excludedGenres"] = config.ExcludedGenres.Movies
	}

	// Add content rating preferences
	if config.MinContentRating != "" {
		filters["minContentRating"] = config.MinContentRating
	}
	if config.MaxContentRating != "" {
		filters["maxContentRating"] = config.MaxContentRating
	}

	// Add age preference
	if config.RecommendationMaxAge > 0 {
		filters["maxYearsOld"] = config.RecommendationMaxAge
	}

	// Add excluded keywords
	if config.ExcludedKeywords != "" {
		filters["excludedKeywords"] = config.ExcludedKeywords
	}

	// Add recommendation strategy
	if config.RecommendationStrategy != "" {
		filters["strategy"] = config.RecommendationStrategy
	}

	// Add recently watched movies
	if len(recentMovies) > 0 {
		filters["recentlyWatched"] = recentMovies
	}

	return filters
}

// Helper to serialize metadata to JSON
func makeMetadataJson(data map[string]interface{}) string {
	if data == nil {
		return "{}"
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "{}"
	}

	return string(jsonBytes)
}

// UpdateUserRecommendationSchedule updates a user's recommendation schedule based on their config
func (j *RecommendationJob) UpdateUserRecommendationSchedule(ctx context.Context, userID uint64) error {
	// Get user configuration
	config, err := j.configRepo.GetUserConfig(ctx, userID)
	if err != nil {
		return fmt.Errorf("error getting user config: %w", err)
	}

	frequency := config.RecommendationSyncFrequency

	// Create or update the user's job schedule
	return j.createUserRecommendationJob(ctx, userID, frequency)
}

// createUserRecommendationJob creates or updates a job schedule for a user
func (j *RecommendationJob) createUserRecommendationJob(ctx context.Context, userID uint64, frequency string) error {
	jobName := fmt.Sprintf("%s.user.%d", j.Name(), userID)

	// Check if the job already exists
	existing, err := j.jobRepo.GetJobSchedule(ctx, jobName)
	if err != nil {
		return fmt.Errorf("error checking for existing job: %w", err)
	}

	// If the job exists, update it
	if existing != nil {
		existing.Frequency = frequency
		existing.Enabled = frequency != string(scheduler.FrequencyManual)
		return j.jobRepo.UpdateJobSchedule(ctx, existing)
	}

	// Create a new job schedule
	schedule := &models.JobSchedule{
		JobName:     jobName,
		JobType:     models.JobTypeRecommendation,
		Frequency:   frequency,
		Enabled:     frequency != string(scheduler.FrequencyManual),
		UserID:      &userID,
		LastRunTime: nil, // Never run yet
	}

	return j.jobRepo.CreateJobSchedule(ctx, schedule)
}

// SetupMediaSyncJob creates or updates a media sync job for a user and client
func (j *RecommendationJob) SetupMediaSyncJob(ctx context.Context, userID, clientID uint64, clientType string, mediaType string, frequency string) error {
	// Check if sync job already exists
	syncJobs, err := j.jobRepo.GetMediaSyncJobsByUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("error checking for existing sync jobs: %w", err)
	}

	// Look for matching job
	var existingJob *models.MediaSyncJob
	for i := range syncJobs {
		if syncJobs[i].ClientID == clientID && syncJobs[i].MediaType == mediaType {
			existingJob = &syncJobs[i]
			break
		}
	}

	// If job exists, update it
	if existingJob != nil {
		existingJob.Frequency = frequency
		existingJob.Enabled = frequency != string(scheduler.FrequencyManual)
		return j.jobRepo.UpdateMediaSyncJob(ctx, existingJob)
	}

	// Create new sync job
	syncJob := &models.MediaSyncJob{
		UserID:     userID,
		ClientID:   clientID,
		ClientType: clientType,
		MediaType:  mediaType,
		Frequency:  frequency,
		Enabled:    frequency != string(scheduler.FrequencyManual),
	}

	return j.jobRepo.CreateMediaSyncJob(ctx, syncJob)
}