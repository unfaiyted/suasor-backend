package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"suasor/client"
	mediatypes "suasor/client/media/types"
	clienttypes "suasor/client/types"
	"suasor/repository"
	"suasor/services/scheduler"
	"suasor/types/models"
	"suasor/utils"
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
	clientRepos     repository.ClientRepositoryCollection
	clientFactories *client.ClientFactoryService
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
	clientRepos repository.ClientRepositoryCollection,
	clientFactories *client.ClientFactoryService,
) *RecommendationJob {
	return &RecommendationJob{
		jobRepo:         jobRepo,
		userRepo:        userRepo,
		configRepo:      configRepo,
		movieRepo:       movieRepo,
		seriesRepo:      seriesRepo,
		musicRepo:       musicRepo,
		historyRepo:     historyRepo,
		clientRepos:     clientRepos,
		clientFactories: clientFactories,
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
	log := utils.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", user.ID).
		Uint64("jobRunID", jobRunID).
		Msg("Generating movie recommendations for user")

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 10, "Analyzing user movie preferences")

	// Get recent user movie history
	recentMovies, err := j.historyRepo.GetRecentUserMovieHistory(ctx, user.ID, 20)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving recent movie history")
		return fmt.Errorf("error retrieving movie history: %w", err)
	}

	// Get user's movies (for determining if recommended movies are already in library)
	userMovies, err := j.movieRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving user movies")
		return fmt.Errorf("error retrieving user movies: %w", err)
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 20, "Building user preference profile")

	// Build user preference profile based on watch history
	preferenceProfile := j.BuildUserMoviePreferences(recentMovies, userMovies, config)

	// Create a map of existing movies in library
	inLibraryMap := make(map[string]bool)

	// Also create a map of watched movies by title-year for easier lookup
	watchedMap := make(map[string]bool)

	for _, movie := range userMovies {
		if movie.Data != nil && movie.Data.Details.Title != "" {
			key := fmt.Sprintf("%s-%d", movie.Data.Details.Title, movie.Data.Details.ReleaseYear)
			inLibraryMap[key] = true

			// If we have watch history for this movie, mark it as watched
			if _, watched := preferenceProfile.WatchedMovies[movie.ID]; watched {
				watchedMap[key] = true
			}
		}
	}

	// Generate recommendation strategies based on user profile
	var recommendations []*models.Recommendation

	// Decide if we should use AI recommendations
	useAI := config.RecommendationSyncEnabled && j.aiClientService != nil

	if useAI {
		// Generate AI recommendations if enabled and available
		j.jobRepo.UpdateJobProgress(ctx, jobRunID, 30, "Generating AI-powered recommendations")

		aiRecs, err := j.generateAIMovieRecommendations(ctx, user.ID, preferenceProfile, config, watchedMap)
		if err != nil {
			// Log error but continue with traditional methods
			log.Error().Err(err).Msg("Error generating AI recommendations, falling back to traditional methods")
		} else {
			recommendations = append(recommendations, aiRecs...)
		}
	}

	// If AI recommendations are disabled or failed, or if we still need more recommendations,
	// use traditional methods as fallback

	if !useAI || len(recommendations) < 5 {
		// Update job progress
		j.jobRepo.UpdateJobProgress(ctx, jobRunID, 40, "Generating recommendations based on genres")

		// 1. Based on preferred genres
		genreBasedRecs := j.generateGenreBasedRecommendations(ctx, user.ID, preferenceProfile, inLibraryMap)
		recommendations = append(recommendations, genreBasedRecs...)

		// Update job progress
		j.jobRepo.UpdateJobProgress(ctx, jobRunID, 60, "Generating recommendations based on similar content")

		// 2. Based on similar content to what they've watched
		similarContentRecs := j.generateSimilarContentRecommendations(ctx, user.ID, recentMovies, inLibraryMap)
		recommendations = append(recommendations, similarContentRecs...)

		// Update job progress
		j.jobRepo.UpdateJobProgress(ctx, jobRunID, 70, "Generating recommendations based on popularity")

		// 3. Popular content they haven't seen
		popularRecs := j.generatePopularRecommendations(ctx, user.ID, inLibraryMap)
		recommendations = append(recommendations, popularRecs...)
	}

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
			log.Error().Err(err).Msg("Error creating batch recommendations")
			return fmt.Errorf("error saving recommendations: %w", err)
		}
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 100, fmt.Sprintf("Generated %d movie recommendations", len(finalRecs)))

	log.Info().
		Uint64("userID", user.ID).
		Int("recommendationCount", len(finalRecs)).
		Msg("Movie recommendations generated successfully")

	return nil
}

// generateAIMovieRecommendations uses AI to generate personalized movie recommendations
func (j *RecommendationJob) generateAIMovieRecommendations(
	ctx context.Context,
	userID uint64,
	profile *UserPreferenceProfile,
	config *models.UserConfig,
	watchedMap map[string]bool) ([]*models.Recommendation, error) {

	log := utils.LoggerFromContext(ctx)

	//TODO: Fix this to properly pull the clientRepos

	aiService, ok := j.clientRepos.(clienttypes.AiClient)
	if !ok {
		return nil, fmt.Errorf("AI client service is not of the expected type")
	}

	// Prepare AI recommendation request
	filters := map[string]interface{}{}

	// Add favorite genres with weights
	if len(profile.FavoriteGenres) > 0 {
		// Sort genres by weight
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

		// Take top 5 genres
		favoriteGenres := []string{}
		for i := 0; i < 5 && i < len(topGenres); i++ {
			favoriteGenres = append(favoriteGenres, topGenres[i].genre)
		}

		filters["favoriteGenres"] = favoriteGenres
	}

	// Add favorite actors
	if len(profile.FavoriteActors) > 0 {
		// Sort actors by weight
		type actorWeight struct {
			actor  string
			weight float32
		}

		var topActors []actorWeight
		for actor, weight := range profile.FavoriteActors {
			topActors = append(topActors, actorWeight{actor, weight})
		}

		// Sort by weight descending
		sort.Slice(topActors, func(i, j int) bool {
			return topActors[i].weight > topActors[j].weight
		})

		// Take top 5 actors
		favoriteActors := []string{}
		for i := 0; i < 5 && i < len(topActors); i++ {
			favoriteActors = append(favoriteActors, topActors[i].actor)
		}

		filters["favoriteActors"] = favoriteActors
	}

	// Add favorite directors
	if len(profile.FavoriteDirectors) > 0 {
		// Sort directors by weight
		type directorWeight struct {
			director string
			weight   float32
		}

		var topDirectors []directorWeight
		for director, weight := range profile.FavoriteDirectors {
			topDirectors = append(topDirectors, directorWeight{director, weight})
		}

		// Sort by weight descending
		sort.Slice(topDirectors, func(i, j int) bool {
			return topDirectors[i].weight > topDirectors[j].weight
		})

		// Take top 3 directors
		favoriteDirectors := []string{}
		for i := 0; i < 3 && i < len(topDirectors); i++ {
			favoriteDirectors = append(favoriteDirectors, topDirectors[i].director)
		}

		filters["favoriteDirectors"] = favoriteDirectors
	}

	// Add top rated/favorite content
	if len(profile.TopRatedContent) > 0 {
		// Take up to 5 top rated items
		topCount := 5
		if len(profile.TopRatedContent) < topCount {
			topCount = len(profile.TopRatedContent)
		}

		filters["favoriteMovies"] = profile.TopRatedContent[:topCount]
	}

	// Add excluded genres
	if len(profile.ExcludedGenres) > 0 {
		filters["excludedGenres"] = profile.ExcludedGenres
	}

	// Add preferred content rating range if set
	if profile.ContentRatingRange[0] != "" || profile.ContentRatingRange[1] != "" {
		filters["contentRatingRange"] = profile.ContentRatingRange
	}

	// Add preferred release years
	filters["releaseYearRange"] = profile.PreferredReleaseYears

	// Add preferred languages if we have any with significant weight
	if len(profile.PreferredLanguages) > 0 {
		// Only include languages with weight > 1.0
		var significantLanguages []string
		for lang, weight := range profile.PreferredLanguages {
			if weight > 1.0 {
				significantLanguages = append(significantLanguages, lang)
			}
		}

		if len(significantLanguages) > 0 {
			filters["preferredLanguages"] = significantLanguages
		}
	}

	// Option to exclude already watched content
	excludeWatched := config.RecommendationIncludeWatched
	filters["excludeWatched"] = excludeWatched

	// Call the AI service to get recommendations
	log.Info().
		Uint64("userID", userID).
		Interface("filters", filters).
		Msg("Requesting AI movie recommendations")

	// Request 10 recommendations
	aiRecommendations, err := aiService.GetRecommendations(ctx, "movie", filters, 10)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get AI recommendations")
		return nil, err
	}

	// Process AI recommendations into our recommendation format
	var recommendations []*models.Recommendation

	for i, rec := range aiRecommendations {
		// Extract details from the recommendation
		title, _ := rec["title"].(string)
		if title == "" {
			continue // Skip recommendations without a title
		}

		// Extract year (could be string or number in the JSON)
		var year int
		if yearVal, ok := rec["year"].(float64); ok {
			year = int(yearVal)
		} else if yearStr, ok := rec["year"].(string); ok {
			if y, err := strconv.Atoi(yearStr); err == nil {
				year = y
			}
		}

		// Create a unique key for this movie to check against watched/library maps
		key := fmt.Sprintf("%s-%d", title, year)

		// Skip if we've seen this and are excluding watched content
		if excludeWatched && watchedMap[key] {
			continue
		}

		// Create a media item for this recommendation
		mediaItem, err := j.findOrCreateMovieItem(ctx, title, year, userID)
		if err != nil {
			log.Error().Err(err).
				Str("title", title).
				Int("year", year).
				Msg("Error creating movie item for AI recommendation")
			continue
		}

		// Extract confidence (if available)
		confidence := float32(0.9) // Default high confidence for AI recommendations
		if conf, ok := rec["confidence"].(float64); ok {
			confidence = float32(conf)
		}

		// Extract reason (if available)
		reason := "AI-powered personalized recommendation"
		if r, ok := rec["reason"].(string); ok && r != "" {
			reason = r
		}

		// Additional metadata
		metadata := map[string]interface{}{
			"source": "ai",
			"rank":   i + 1,
		}

		// Add genres if available
		if genres, ok := rec["genres"].([]interface{}); ok {
			genreStrs := make([]string, 0, len(genres))
			for _, g := range genres {
				if gs, ok := g.(string); ok {
					genreStrs = append(genreStrs, gs)
				}
			}
			metadata["genres"] = genreStrs
		}

		// Determine if it's in the user's library
		isInLibrary := watchedMap[key]

		// Create the recommendation
		recommendation := &models.Recommendation{
			UserID:      userID,
			MediaItemID: mediaItem.ID,
			MediaType:   string(mediatypes.MediaTypeMovie),
			Source:      models.RecommendationSourceAI,
			Reason:      reason,
			Confidence:  confidence,
			InLibrary:   isInLibrary,
			Viewed:      watchedMap[key],
			Active:      true,
			Metadata:    makeMetadataJson(metadata),
		}

		recommendations = append(recommendations, recommendation)
	}

	log.Info().
		Uint64("userID", userID).
		Int("aiRecommendationsCount", len(recommendations)).
		Msg("AI movie recommendations generated")

	return recommendations, nil
}

// generateSystemMovieRecommendations creates movie recommendations using system algorithms
func (j *RecommendationJob) generateSystemMovieRecommendations(ctx context.Context, user models.User, config *models.UserConfig, jobRunID uint64) error {
	log.Printf("Generating system movie recommendations for user %s", user.Username)
	return j.generateMovieRecommendations(ctx, user, config, jobRunID)
}

// BuildUserMusicPreferences builds a user preference profile for music
func (j *RecommendationJob) BuildUserMusicPreferences(
	recentTracks []models.MediaPlayHistory[*mediatypes.Track],
	userTracks []*models.MediaItem[*mediatypes.Track],
	config *models.UserConfig) *UserPreferenceProfile {

	profile := &UserPreferenceProfile{
		FavoriteGenres:     make(map[string]float32),
		WatchedMovies:      make(map[uint64]bool), // Reusing this field for track IDs
		PreferredLanguages: make(map[string]float32),
		MinRating:          0.0, // No minimum rating for music
	}

	// No specific release year constraints for music by default
	currentYear := time.Now().Year()
	profile.PreferredReleaseYears = [2]int{currentYear - 50, currentYear}

	// Add excluded genres from user preferences if available
	if config.ExcludedGenres != nil && len(config.ExcludedGenres.Music) > 0 {
		profile.ExcludedGenres = config.ExcludedGenres.Music
	}

	// Gather favorite artists and tracks
	favoriteArtists := make(map[string]float32)
	topRatedTracks := make([]MusicSummary, 0)

	// Process play history to build preference profile
	for _, history := range recentTracks {
		if history.Item == nil || history.Item.Data == nil {
			continue
		}

		track := history.Item.Data

		// Mark track as played
		profile.WatchedMovies[history.MediaItemID] = true

		// Weight based on how recently and how many times played
		weight := float32(1.0)
		if history.PlayCount > 0 {
			weight += float32(history.PlayCount) * 0.5 // More plays = higher weight
		}

		// More recent plays have higher weight
		if !history.LastPlayedAt.IsZero() {
			daysAgo := time.Since(history.LastPlayedAt).Hours() / 24
			if daysAgo < 7 {
				weight += 1.0 // Played within last week
			} else if daysAgo < 30 {
				weight += 0.5 // Played within last month
			}
		}

		// Process track genres
		for _, genre := range track.Details.Genres {
			profile.FavoriteGenres[genre] += weight
		}

		// Add to top rated if user liked it
		if history.IsFavorite || history.PlayCount > 2 {
			// Create detailed rating information
			// Get the overall rating from the user rating or the first available rating
			var overallRating float32
			if track.Details.UserRating > 0 {
				overallRating = track.Details.UserRating
			} else if len(track.Details.Ratings) > 0 {
				overallRating = track.Details.Ratings[0].Value
			}

			detailedRating := &RatingDetails{
				Overall:    overallRating,
				MaxValue:   10.0,
				Source:     "user",
				Timestamp:  history.LastPlayedAt.Unix(),
				Categories: map[string]float32{},
			}

			// Add any detailed ratings if available
			if track.Details.Ratings != nil {
				for _, rating := range track.Details.Ratings {
					detailedRating.Categories[rating.Source] = float32(rating.Value)
				}
			}

			topRatedTracks = append(topRatedTracks, MusicSummary{
				Title:          track.Details.Title,
				Artist:         track.ArtistName,
				Album:          track.AlbumName,
				Year:           track.Details.ReleaseYear,
				Genres:         track.Details.Genres,
				DetailedRating: detailedRating,
				PlayCount:      int(history.PlayCount),
				IsFavorite:     history.IsFavorite,
				LastPlayDate:   history.LastPlayedAt.Unix(),
				DurationSec:    track.Duration,
				// UserTags:       track.Tags,
			})
		}

		// Add artist to favorites with weight
		if track.ArtistName != "" {
			favoriteArtists[track.ArtistName] += weight
		}
	}

	// Add user's explicitly preferred genres if set in their config
	if config.PreferredGenres != nil && len(config.PreferredGenres.Music) > 0 {
		for _, genre := range config.PreferredGenres.Music {
			// Give extra weight to explicitly preferred genres
			profile.FavoriteGenres[genre] += 3.0
		}
	}

	// Store the top rated tracks for AI recommendations
	profile.TopRatedContent = make([]MovieSummary, len(topRatedTracks))
	for i, track := range topRatedTracks {
		profile.TopRatedContent[i] = MovieSummary{
			Title:      track.Title,
			Year:       track.Year,
			Genres:     track.Genres,
			IsFavorite: track.IsFavorite,
			PlayCount:  track.PlayCount,
		}
	}

	return profile
}

// BuildUserSeriesPreferences builds a user preference profile for TV shows
func (j *RecommendationJob) BuildUserSeriesPreferences(
	recentSeries []models.MediaPlayHistory[*mediatypes.Series],
	userSeries []*models.MediaItem[*mediatypes.Series],
	config *models.UserConfig) *UserPreferenceProfile {

	profile := &UserPreferenceProfile{
		FavoriteGenres:     make(map[string]float32),
		FavoriteActors:     make(map[string]float32),
		WatchedMovies:      make(map[uint64]bool), // Reusing this field for series IDs
		PreferredLanguages: make(map[string]float32),
		MinRating:          5.0, // Default minimum rating
	}

	// Set default release year range to past 10 years (TV shows often run longer than movies)
	currentYear := time.Now().Year()
	profile.PreferredReleaseYears = [2]int{currentYear - 10, currentYear}

	// Add excluded genres from user preferences if available
	if config.ExcludedGenres != nil && len(config.ExcludedGenres.Series) > 0 {
		profile.ExcludedGenres = config.ExcludedGenres.Series
	}

	// User may have preferred content rating limits
	if config.MinContentRating != "" || config.MaxContentRating != "" {
		profile.ContentRatingRange = [2]string{config.MinContentRating, config.MaxContentRating}
	}

	// Handle preferred age of content
	if config.RecommendationMaxAge > 0 {
		minYear := currentYear - config.RecommendationMaxAge
		profile.PreferredReleaseYears[0] = minYear
	}

	// Gather favorite series
	topRatedSeries := make([]SeriesSummary, 0)

	// Process watch history to build preference profile
	for _, history := range recentSeries {
		if history.Item == nil || history.Item.Data == nil {
			continue
		}

		series := history.Item.Data

		// Mark series as watched
		profile.WatchedMovies[history.MediaItemID] = true

		// Weight based on how recently and how many times played
		weight := float32(1.0)
		if history.PlayCount > 0 {
			weight += float32(history.PlayCount) * 0.5 // More plays = higher weight
		}

		// More recent watches have higher weight
		if !history.LastPlayedAt.IsZero() {
			daysAgo := time.Since(history.LastPlayedAt).Hours() / 24
			if daysAgo < 7 {
				weight += 1.0 // Watched within last week
			} else if daysAgo < 30 {
				weight += 0.5 // Watched within last month
			}
		}

		// Process genres
		for _, genre := range series.Details.Genres {
			profile.FavoriteGenres[genre] += weight
		}

		// Add to top rated if user liked it
		if history.IsFavorite || history.PlayCount > 1 {
			// Create detailed rating information
			detailedRating := &RatingDetails{
				Overall:    float32(series.Rating),
				MaxValue:   10.0,
				Source:     "user",
				Timestamp:  history.LastPlayedAt.Unix(),
				Categories: map[string]float32{},
			}

			// Add any detailed ratings if available
			if series.Details.Ratings != nil {
				for _, rating := range series.Details.Ratings {
					detailedRating.Categories[rating.Source] = float32(rating.Value)
				}
			}

			topRatedSeries = append(topRatedSeries, SeriesSummary{
				Title:          series.Details.Title,
				Year:           series.ReleaseYear,
				Genres:         series.Genres,
				Rating:         float32(series.Rating),
				DetailedRating: detailedRating,
				Seasons:        series.SeasonCount,
				Status:         series.Status,
				IsFavorite:     history.IsFavorite,
				// EpisodesWatched is no longer available in MediaPlayHistory
				// Using PlayCount as a substitute measure of engagement
				EpisodesWatched: int(history.PlayCount),
				TotalEpisodes:   series.EpisodeCount,
				LastWatchDate:   history.LastPlayedAt.Unix(),
				// UserTags:        series.Tags,
			})
		}

		// Add language preference if available
		if series.Details.Language != "" {
			profile.PreferredLanguages[series.Details.Language] += weight
		}
	}

	// Add user's explicitly preferred genres if set in their config
	if config.PreferredGenres != nil && len(config.PreferredGenres.Series) > 0 {
		for _, genre := range config.PreferredGenres.Series {
			// Give extra weight to explicitly preferred genres
			profile.FavoriteGenres[genre] += 3.0
		}
	}

	// Store the top rated series for AI recommendations
	profile.TopRatedContent = make([]MovieSummary, len(topRatedSeries))
	for i, series := range topRatedSeries {
		profile.TopRatedContent[i] = MovieSummary{
			Title:  series.Title,
			Year:   series.Year,
			Genres: series.Genres,
			// Rating:     series.Rating,
			IsFavorite: series.IsFavorite,
		}
	}

	return profile
}

// generateSeriesRecommendations creates TV series recommendations for a user
func (j *RecommendationJob) generateSeriesRecommendations(ctx context.Context, user models.User, config *models.UserConfig, jobRunID uint64) error {
	log := utils.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", user.ID).
		Uint64("jobRunID", jobRunID).
		Msg("Generating series recommendations for user")

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 10, "Analyzing user TV show preferences")

	// Get recent user series watch history
	recentSeries, err := j.historyRepo.GetRecentUserSeriesHistory(ctx, user.ID, 20)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving recent series history")
		return fmt.Errorf("error retrieving series history: %w", err)
	}

	// Get user's series (for determining if recommended series are already in library)
	userSeries, err := j.seriesRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving user series")
		return fmt.Errorf("error retrieving user series: %w", err)
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 25, "Building user TV preference profile")

	// Build preference profile
	preferenceProfile := j.BuildUserSeriesPreferences(recentSeries, userSeries, config)

	// Create maps for library and watched status lookups
	inLibraryMap := make(map[string]bool)
	watchedMap := make(map[string]bool)

	for _, series := range userSeries {
		if series.Data != nil && series.Data.Details.Title != "" {
			key := fmt.Sprintf("%s", series.Data.Details.Title)
			inLibraryMap[key] = true

			// If we have watch history for this series, mark it as watched
			if _, watched := preferenceProfile.WatchedMovies[series.ID]; watched {
				watchedMap[key] = true
			}
		}
	}

	// Generate recommendation strategies based on user profile
	var recommendations []*models.Recommendation

	// Decide if we should use AI recommendations
	useAI := config.RecommendationSyncEnabled && j.aiClientService != nil

	if useAI {
		// Generate AI recommendations if enabled and available
		j.jobRepo.UpdateJobProgress(ctx, jobRunID, 30, "Generating AI-powered TV show recommendations")

		aiRecs, err := j.generateAISeriesRecommendations(ctx, user.ID, preferenceProfile, config, watchedMap)
		if err != nil {
			// Log error but continue with traditional methods
			log.Error().Err(err).Msg("Error generating AI recommendations for TV shows, falling back to traditional methods")
		} else {
			recommendations = append(recommendations, aiRecs...)
		}
	}

	// If AI recommendations are disabled, failed, or we need more recommendations,
	// use traditional method as fallback
	if !useAI || len(recommendations) < 5 {
		// Update job progress
		j.jobRepo.UpdateJobProgress(ctx, jobRunID, 50, "Creating series recommendations using traditional methods")

		// Generate sample series recommendations based on genres and popular shows
		// This is a simplified implementation - in a real system, you would have more
		// sophisticated algorithms that consider user preferences and watching patterns
		seriesRecommendations := []struct {
			Title      string
			Year       int
			Genre      string
			Rating     float32
			Confidence float32
			Reason     string
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
			{
				Title:      "Succession",
				Year:       2018,
				Genre:      "Drama",
				Rating:     8.9,
				Confidence: 0.85,
				Reason:     "Award-winning drama about family dynamics",
			},
			{
				Title:      "Wednesday",
				Year:       2022,
				Genre:      "Comedy",
				Rating:     8.2,
				Confidence: 0.75,
				Reason:     "Popular supernatural comedy",
			},
			{
				Title:      "The Boys",
				Year:       2019,
				Genre:      "Action",
				Rating:     8.7,
				Confidence: 0.8,
				Reason:     "Subversive take on the superhero genre",
			},
		}

		// Process all recommended series
		for _, rec := range seriesRecommendations {
			// Check if already in library
			isInLibrary := inLibraryMap[rec.Title]

			// Skip if we've already watched this series
			if watchedMap[rec.Title] {
				continue
			}

			// Create media item if not exists
			mediaItem, err := j.findOrCreateSeriesItem(ctx, rec.Title, rec.Year, rec.Genre, user.ID)
			if err != nil {
				log.Error().Err(err).
					Str("title", rec.Title).
					Msg("Error creating series item")
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
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 80, "Filtering and ranking series recommendations")

	// Filter out any duplicates and limit total recommendations
	finalRecs := j.FilterAndRankRecommendations(recommendations, 8)

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 90, "Saving series recommendations")

	// Save all recommendations in batch
	if len(finalRecs) > 0 {
		if err := j.jobRepo.BatchCreateRecommendations(ctx, finalRecs); err != nil {
			log.Error().Err(err).Msg("Error creating batch series recommendations")
			return fmt.Errorf("error saving series recommendations: %w", err)
		}
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 100, fmt.Sprintf("Generated %d series recommendations", len(finalRecs)))

	log.Info().
		Uint64("userID", user.ID).
		Int("recommendationCount", len(finalRecs)).
		Msg("Series recommendations generated successfully")

	return nil
}

// generateAISeriesRecommendations uses AI to generate personalized TV series recommendations
func (j *RecommendationJob) generateAISeriesRecommendations(
	ctx context.Context,
	userID uint64,
	profile *UserPreferenceProfile,
	config *models.UserConfig,
	watchedMap map[string]bool) ([]*models.Recommendation, error) {

	log := utils.LoggerFromContext(ctx)

	// We need to cast AI client service to the proper type
	aiService, ok := j.aiClientService.(clienttypes.AiClient)
	if !ok {
		return nil, fmt.Errorf("AI client service is not of the expected type")
	}

	// Prepare AI recommendation request
	filters := map[string]interface{}{}

	// Add favorite genres with weights
	if len(profile.FavoriteGenres) > 0 {
		// Sort genres by weight
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

		// Take top 5 genres
		favoriteGenres := []string{}
		for i := 0; i < 5 && i < len(topGenres); i++ {
			favoriteGenres = append(favoriteGenres, topGenres[i].genre)
		}

		filters["favoriteGenres"] = favoriteGenres
	}

	// Add top rated/favorite content
	if len(profile.TopRatedContent) > 0 {
		// Take up to 5 top rated items
		topCount := 5
		if len(profile.TopRatedContent) < topCount {
			topCount = len(profile.TopRatedContent)
		}

		filters["favoriteSeries"] = profile.TopRatedContent[:topCount]
	}

	// Add excluded genres
	if len(profile.ExcludedGenres) > 0 {
		filters["excludedGenres"] = profile.ExcludedGenres
	}

	// Add preferred content rating range if set
	if profile.ContentRatingRange[0] != "" || profile.ContentRatingRange[1] != "" {
		filters["contentRatingRange"] = profile.ContentRatingRange
	}

	// Add preferred release years
	filters["releaseYearRange"] = profile.PreferredReleaseYears

	// Option to exclude already watched content
	excludeWatched := !config.RecommendationIncludeWatched
	filters["excludeWatched"] = excludeWatched

	// Call the AI service to get recommendations
	log.Info().
		Uint64("userID", userID).
		Interface("filters", filters).
		Msg("Requesting AI series recommendations")

	// Request 8 recommendations
	aiRecommendations, err := aiService.GetRecommendations(ctx, "series", filters, 8)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get AI series recommendations")
		return nil, err
	}

	// Process AI recommendations into our recommendation format
	var recommendations []*models.Recommendation

	for i, rec := range aiRecommendations {
		// Extract details from the recommendation
		title, _ := rec["title"].(string)
		if title == "" {
			continue // Skip recommendations without a title
		}

		// Extract year (could be string or number in the JSON)
		var year int
		if yearVal, ok := rec["year"].(float64); ok {
			year = int(yearVal)
		} else if yearStr, ok := rec["year"].(string); ok {
			if y, err := strconv.Atoi(yearStr); err == nil {
				year = y
			}
		}

		// Extract genre if available
		genre := "Drama" // Default genre
		if genreVal, ok := rec["genre"].(string); ok && genreVal != "" {
			genre = genreVal
		} else if genres, ok := rec["genres"].([]interface{}); ok && len(genres) > 0 {
			if g, ok := genres[0].(string); ok {
				genre = g
			}
		}

		// Skip if we've seen this and are excluding watched content
		if excludeWatched && watchedMap[title] {
			continue
		}

		// Create a media item for this recommendation
		mediaItem, err := j.findOrCreateSeriesItem(ctx, title, year, genre, userID)
		if err != nil {
			log.Error().Err(err).
				Str("title", title).
				Int("year", year).
				Msg("Error creating series item for AI recommendation")
			continue
		}

		// Extract confidence (if available)
		confidence := float32(0.9) // Default high confidence for AI recommendations
		if conf, ok := rec["confidence"].(float64); ok {
			confidence = float32(conf)
		}

		// Extract reason (if available)
		reason := "AI-powered personalized TV series recommendation"
		if r, ok := rec["reason"].(string); ok && r != "" {
			reason = r
		}

		// Additional metadata
		metadata := map[string]interface{}{
			"source": "ai",
			"rank":   i + 1,
		}

		// Add genres if available
		if genres, ok := rec["genres"].([]interface{}); ok {
			genreStrs := make([]string, 0, len(genres))
			for _, g := range genres {
				if gs, ok := g.(string); ok {
					genreStrs = append(genreStrs, gs)
				}
			}
			metadata["genres"] = genreStrs
		}

		// Determine if it's in the user's library - use watchedMap since inLibraryMap isn't available in this context
		isInLibrary := watchedMap[title]

		// Create recommendation
		recommendation := &models.Recommendation{
			UserID:      userID,
			MediaItemID: mediaItem.ID,
			MediaType:   string(mediatypes.MediaTypeSeries),
			Source:      models.RecommendationSourceAI,
			Reason:      reason,
			Confidence:  confidence,
			InLibrary:   isInLibrary,
			Viewed:      watchedMap[title],
			Active:      true,
			Metadata:    makeMetadataJson(metadata),
		}

		recommendations = append(recommendations, recommendation)
	}

	log.Info().
		Uint64("userID", userID).
		Int("aiRecommendationsCount", len(recommendations)).
		Msg("AI series recommendations generated")

	return recommendations, nil
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
	log := utils.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", user.ID).
		Uint64("jobRunID", jobRunID).
		Msg("Generating music recommendations for user")

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 10, "Analyzing user music preferences")

	// Get recent user music play history
	recentTracks, err := j.historyRepo.GetRecentUserMusicHistory(ctx, user.ID, 20)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving recent music history")
		return fmt.Errorf("error retrieving music history: %w", err)
	}

	// Get user's music (for determining if recommended tracks are already in library)
	userMusic, err := j.musicRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving user music")
		return fmt.Errorf("error retrieving user music: %w", err)
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 25, "Building user music preference profile")

	// Build preference profile
	preferenceProfile := j.BuildUserMusicPreferences(recentTracks, userMusic, config)

	// Create maps for library and played status lookups
	inLibraryMap := make(map[string]bool)
	playedMap := make(map[string]bool)

	for _, track := range userMusic {
		if track.Data != nil && track.Data.Details.Title != "" && track.Data.ArtistName != "" {
			key := fmt.Sprintf("%s-%s", track.Data.ArtistName, track.Data.Details.Title)
			inLibraryMap[key] = true

			// If we have play history for this track, mark it as played
			if _, played := preferenceProfile.WatchedMovies[track.ID]; played {
				playedMap[key] = true
			}
		}
	}

	// Generate recommendation strategies based on user profile
	var recommendations []*models.Recommendation

	// Decide if we should use AI recommendations
	useAI := config.RecommendationSyncEnabled && j.aiClientService != nil

	if useAI {
		// Generate AI recommendations if enabled and available
		j.jobRepo.UpdateJobProgress(ctx, jobRunID, 30, "Generating AI-powered music recommendations")

		aiRecs, err := j.generateAIMusicRecommendations(ctx, user.ID, preferenceProfile, config, playedMap)
		if err != nil {
			// Log error but continue with traditional methods
			log.Error().Err(err).Msg("Error generating AI recommendations for music, falling back to traditional methods")
		} else {
			recommendations = append(recommendations, aiRecs...)
		}
	}

	// If AI recommendations are disabled, failed, or we need more recommendations,
	// use traditional method as fallback
	if !useAI || len(recommendations) < 5 {
		// Update job progress
		j.jobRepo.UpdateJobProgress(ctx, jobRunID, 50, "Creating music recommendations using traditional methods")

		// Generate sample music recommendations based on genres and popular tracks
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
			{
				Title:      "Midnight City",
				Artist:     "M83",
				Album:      "Hurry Up, We're Dreaming",
				Year:       2011,
				Genre:      "Electronic",
				Confidence: 0.8,
				Reason:     "Popular electronic track with wide appeal",
			},
			{
				Title:      "Hey Ya!",
				Artist:     "OutKast",
				Album:      "Speakerboxxx/The Love Below",
				Year:       2003,
				Genre:      "Hip-Hop",
				Confidence: 0.9,
				Reason:     "Classic hip-hop track with wide appeal",
			},
			{
				Title:      "Bohemian Rhapsody",
				Artist:     "Queen",
				Album:      "A Night at the Opera",
				Year:       1975,
				Genre:      "Rock",
				Confidence: 0.95,
				Reason:     "One of the most iconic rock songs of all time",
			},
		}

		// Process all recommended tracks
		for _, rec := range musicRecommendations {
			// Check if already in library
			key := fmt.Sprintf("%s-%s", rec.Artist, rec.Title)
			isInLibrary := inLibraryMap[key]

			// Skip if already played and we're excluding played content
			if !config.RecommendationIncludeWatched && playedMap[key] {
				continue
			}

			// Create media item if not exists
			mediaItem, err := j.findOrCreateMusicItem(ctx, rec.Title, rec.Artist, rec.Album, rec.Year, rec.Genre, user.ID)
			if err != nil {
				log.Error().Err(err).
					Str("title", rec.Title).
					Str("artist", rec.Artist).
					Msg("Error creating music item")
				continue
			}

			// Skip if user has already played this track
			viewed, _ := j.historyRepo.HasUserViewedMedia(ctx, user.ID, mediaItem.ID)
			if viewed && !config.RecommendationIncludeWatched {
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
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 80, "Filtering and ranking music recommendations")

	// Filter out any duplicates and limit total recommendations
	finalRecs := j.FilterAndRankRecommendations(recommendations, 8)

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 90, "Saving music recommendations")

	// Save all recommendations in batch
	if len(finalRecs) > 0 {
		if err := j.jobRepo.BatchCreateRecommendations(ctx, finalRecs); err != nil {
			log.Error().Err(err).Msg("Error creating batch music recommendations")
			return fmt.Errorf("error saving music recommendations: %w", err)
		}
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 100, fmt.Sprintf("Generated %d music recommendations", len(finalRecs)))

	log.Info().
		Uint64("userID", user.ID).
		Int("recommendationCount", len(finalRecs)).
		Msg("Music recommendations generated successfully")

	return nil
}

// generateAIMusicRecommendations uses AI to generate personalized music recommendations
func (j *RecommendationJob) generateAIMusicRecommendations(
	ctx context.Context,
	userID uint64,
	profile *UserPreferenceProfile,
	config *models.UserConfig,
	playedMap map[string]bool) ([]*models.Recommendation, error) {

	log := utils.LoggerFromContext(ctx)

	// We need to cast AI client service to the proper type
	aiService, ok := j.aiClientService.(clienttypes.AiClient)
	if !ok {
		return nil, fmt.Errorf("AI client service is not of the expected type")
	}

	// Prepare AI recommendation request
	filters := map[string]interface{}{}

	// Add favorite genres with weights
	if len(profile.FavoriteGenres) > 0 {
		// Sort genres by weight
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

		// Take top 5 genres
		favoriteGenres := []string{}
		for i := 0; i < 5 && i < len(topGenres); i++ {
			favoriteGenres = append(favoriteGenres, topGenres[i].genre)
		}

		filters["favoriteGenres"] = favoriteGenres
	}

	// Add top rated/favorite content
	if len(profile.TopRatedContent) > 0 {
		// Take up to 5 top rated items
		topCount := 5
		if len(profile.TopRatedContent) < topCount {
			topCount = len(profile.TopRatedContent)
		}

		filters["favoriteMusic"] = profile.TopRatedContent[:topCount]
	}

	// Add excluded genres
	if len(profile.ExcludedGenres) > 0 {
		filters["excludedGenres"] = profile.ExcludedGenres
	}

	// Add preferred release years
	filters["releaseYearRange"] = profile.PreferredReleaseYears

	// Option to exclude already watched content
	excludePlayed := !config.RecommendationIncludeWatched
	filters["excludePlayed"] = excludePlayed

	// Call the AI service to get recommendations
	log.Info().
		Uint64("userID", userID).
		Interface("filters", filters).
		Msg("Requesting AI music recommendations")

	// Request 8 recommendations
	aiRecommendations, err := aiService.GetRecommendations(ctx, "music", filters, 8)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get AI music recommendations")
		return nil, err
	}

	// Process AI recommendations into our recommendation format
	var recommendations []*models.Recommendation

	for i, rec := range aiRecommendations {
		// Extract details from the recommendation
		title, _ := rec["title"].(string)
		if title == "" {
			continue // Skip recommendations without a title
		}

		// Extract artist (required for music)
		artist, _ := rec["artist"].(string)
		if artist == "" {
			continue // Skip recommendations without an artist
		}

		// Extract album (optional)
		album, _ := rec["album"].(string)

		// Extract year (could be string or number in the JSON)
		var year int
		if yearVal, ok := rec["year"].(float64); ok {
			year = int(yearVal)
		} else if yearStr, ok := rec["year"].(string); ok {
			if y, err := strconv.Atoi(yearStr); err == nil {
				year = y
			}
		}

		// Extract genre if available
		genre := "Pop" // Default genre
		if genreVal, ok := rec["genre"].(string); ok && genreVal != "" {
			genre = genreVal
		} else if genres, ok := rec["genres"].([]interface{}); ok && len(genres) > 0 {
			if g, ok := genres[0].(string); ok {
				genre = g
			}
		}

		// Create a unique key for this track to check against played map
		key := fmt.Sprintf("%s-%s", artist, title)

		// Skip if we've played this and are excluding played content
		if excludePlayed && playedMap[key] {
			continue
		}

		// Create a media item for this recommendation
		mediaItem, err := j.findOrCreateMusicItem(ctx, title, artist, album, year, genre, userID)
		if err != nil {
			log.Error().Err(err).
				Str("title", title).
				Str("artist", artist).
				Msg("Error creating music item for AI recommendation")
			continue
		}

		// Extract confidence (if available)
		confidence := float32(0.9) // Default high confidence for AI recommendations
		if conf, ok := rec["confidence"].(float64); ok {
			confidence = float32(conf)
		}

		// Extract reason (if available)
		reason := "AI-powered personalized music recommendation"
		if r, ok := rec["reason"].(string); ok && r != "" {
			reason = r
		}

		// Additional metadata
		metadata := map[string]interface{}{
			"source": "ai",
			"rank":   i + 1,
		}

		// Add genres if available
		if genres, ok := rec["genres"].([]interface{}); ok {
			genreStrs := make([]string, 0, len(genres))
			for _, g := range genres {
				if gs, ok := g.(string); ok {
					genreStrs = append(genreStrs, gs)
				}
			}
			metadata["genres"] = genreStrs
		}

		// Add album if available
		if album != "" {
			metadata["album"] = album
		}

		// Determine if it's in the user's library
		isInLibrary := playedMap[key]

		// Create recommendation
		recommendation := &models.Recommendation{
			UserID:      userID,
			MediaItemID: mediaItem.ID,
			MediaType:   string(mediatypes.MediaTypeTrack),
			Source:      models.RecommendationSourceAI,
			Reason:      reason,
			Confidence:  confidence,
			InLibrary:   isInLibrary,
			Viewed:      playedMap[key],
			Active:      true,
			Metadata:    makeMetadataJson(metadata),
		}

		recommendations = append(recommendations, recommendation)
	}

	log.Info().
		Uint64("userID", userID).
		Int("aiRecommendationsCount", len(recommendations)).
		Msg("AI music recommendations generated")

	return recommendations, nil
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

// MovieSummary contains a summary of a movie for recommendation purposes
type MovieSummary struct {
	Title             string         `json:"title"`
	Year              int            `json:"year"`
	Genres            []string       `json:"genres,omitempty"`
	DetailedRating    *RatingDetails `json:"detailedRating,omitempty"` // Enhanced rating information
	PlayCount         int            `json:"playCount,omitempty"`
	IsFavorite        bool           `json:"isFavorite,omitempty"`
	CompletionPercent float32        `json:"completionPercent,omitempty"`
	WatchDate         int64          `json:"watchDate,omitempty"` // Unix timestamp of last watch
	UserTags          []string       `json:"userTags,omitempty"`  // Custom tags applied by the user
}

// SeriesSummary contains a summary of a TV series for recommendation purposes
type SeriesSummary struct {
	Title           string         `json:"title"`
	Year            int            `json:"year"`
	Genres          []string       `json:"genres,omitempty"`
	Rating          float32        `json:"rating,omitempty"`         // Basic rating for compatibility
	DetailedRating  *RatingDetails `json:"detailedRating,omitempty"` // Enhanced rating information
	Seasons         int            `json:"seasons,omitempty"`
	Status          string         `json:"status,omitempty"` // e.g., "Ended", "Continuing"
	IsFavorite      bool           `json:"isFavorite,omitempty"`
	EpisodesWatched int            `json:"episodesWatched,omitempty"` // Number of episodes watched
	TotalEpisodes   int            `json:"totalEpisodes,omitempty"`   // Total episodes in the series
	LastWatchDate   int64          `json:"lastWatchDate,omitempty"`   // Unix timestamp of last watch
	UserTags        []string       `json:"userTags,omitempty"`        // Custom tags applied by the user
}

// RatingDetails contains detailed rating information
type RatingDetails struct {
	Overall    float32            `json:"overall,omitempty"`    // Overall rating (0-10)
	Categories map[string]float32 `json:"categories,omitempty"` // Ratings for specific categories like "Acting", "Story", etc.
	Source     string             `json:"source,omitempty"`     // Source of the rating (user, aggregated, external service)
	MaxValue   float32            `json:"maxValue,omitempty"`   // Maximum possible rating value (default: 10)
	UserCount  int                `json:"userCount,omitempty"`  // Number of users who rated (for aggregated ratings)
	Timestamp  int64              `json:"timestamp,omitempty"`  // When the rating was last updated
}

// MusicSummary contains a summary of a music track/artist for recommendation purposes
type MusicSummary struct {
	Title          string         `json:"title"`
	Artist         string         `json:"artist"`
	Album          string         `json:"album,omitempty"`
	Year           int            `json:"year,omitempty"`
	Genres         []string       `json:"genres,omitempty"`
	Rating         float32        `json:"rating,omitempty"`         // Basic rating for compatibility
	DetailedRating *RatingDetails `json:"detailedRating,omitempty"` // Enhanced rating information
	PlayCount      int            `json:"playCount,omitempty"`
	IsFavorite     bool           `json:"isFavorite,omitempty"`
	LastPlayDate   int64          `json:"lastPlayDate,omitempty"` // Unix timestamp of last play
	DurationSec    int            `json:"durationSec,omitempty"`  // Duration in seconds
	UserTags       []string       `json:"userTags,omitempty"`     // Custom tags applied by the user
}

// UserPreferenceProfile represents a user's media preferences
type UserPreferenceProfile struct {
	// Movie preferences
	FavoriteGenres        map[string]float32 // Genre -> weight
	FavoriteActors        map[string]float32 // Actor -> weight
	FavoriteDirectors     map[string]float32 // Director -> weight
	WatchedMovies         map[uint64]bool    // MediaItemID -> watched
	ExcludedGenres        []string           // Genres to exclude
	PreferredReleaseYears [2]int             // [min, max] years
	ContentRatingRange    [2]string          // [min, max] content ratings
	MinRating             float32            // Minimum rating (0-10)
	PreferredLanguages    map[string]float32 // Language -> weight
	TopRatedContent       []MovieSummary     // Top rated/favorite content
	RecentlyWatched       []MovieSummary     // Recently watched content

	// Recommended content behavior
	ExcludeWatched       bool // Whether to exclude content already watched
	IncludeSimilarItems  bool // Whether to include items similar to favorites
	UseAIRecommendations bool // Whether to use AI for recommendations
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
		profile.ContentRatingRange = [2]string{config.MinContentRating, config.MaxContentRating}
	}

	// Handle preferred age of content
	if config.RecommendationMaxAge > 0 {
		minYear := currentYear - config.RecommendationMaxAge
		profile.PreferredReleaseYears[0] = minYear
	}

	// Gather favorite movies by rating
	topRatedMovies := make([]MovieSummary, 0)

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
		if !history.LastPlayedAt.IsZero() {
			daysAgo := time.Since(history.LastPlayedAt).Hours() / 24
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

		// Add to top rated if user liked it
		if history.IsFavorite || history.PlayCount > 1 {
			// Create detailed rating information
			// Get the overall rating from the user rating or the first available rating
			var overallRating float32
			if movie.Details.UserRating > 0 {
				overallRating = movie.Details.UserRating
			} else if len(movie.Details.Ratings) > 0 {
				overallRating = movie.Details.Ratings[0].Value
			}

			detailedRating := &RatingDetails{
				Overall:    overallRating,
				MaxValue:   10.0,
				Source:     "user",
				Timestamp:  history.LastPlayedAt.Unix(),
				Categories: map[string]float32{},
			}

			// Add any detailed ratings if available
			if movie.Details.Ratings != nil {
				for _, rating := range movie.Details.Ratings {
					detailedRating.Categories[rating.Source] = float32(rating.Value)
				}
			}

			topRatedMovies = append(topRatedMovies, MovieSummary{
				Title:  movie.Details.Title,
				Year:   movie.Details.ReleaseYear,
				Genres: movie.Details.Genres,
				// Maintain backwards compatibility with the basic rating
				DetailedRating:    detailedRating,
				PlayCount:         int(history.PlayCount),
				IsFavorite:        history.IsFavorite,
				CompletionPercent: float32(history.PlayedPercentage),
				WatchDate:         history.LastPlayedAt.Unix(),
				// UserTags:          movie.Tags,
			})
		}

		// Add language preference if available
		if movie.Details.Language != "" {
			profile.PreferredLanguages[movie.Details.Language] += weight
		}
	}

	// Add user's explicitly preferred genres if set in their config
	if config.PreferredGenres != nil && len(config.PreferredGenres.Movies) > 0 {
		for _, genre := range config.PreferredGenres.Movies {
			// Give extra weight to explicitly preferred genres
			profile.FavoriteGenres[genre] += 3.0
		}
	}

	// Store the top rated movies for use with AI-based recommendations
	profile.TopRatedContent = topRatedMovies

	// Find most frequently watched actors & directors by counting occurrences
	// This helps identify patterns that might not be evident from simply counting favorites
	for _, movie := range userMovies {
		if movie.Data == nil {
			continue
		}

		// Skip if this movie hasn't been watched
		if _, watched := profile.WatchedMovies[movie.ID]; !watched {
			continue
		}

		// Process cast
		for _, person := range movie.Data.Cast {
			if person.Role == "actor" {
				profile.FavoriteActors[person.Name] += 0.5 // Lower weight for just being in library
			} else if person.Role == "director" {
				profile.FavoriteDirectors[person.Name] += 0.5
			}
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
