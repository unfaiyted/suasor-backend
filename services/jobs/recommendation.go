package jobs

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"suasor/client"
	"suasor/client/ai"
	aitypes "suasor/client/ai/types"
	mediatypes "suasor/client/media/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/utils"
	"time"
)

// Helper functions for handling credits

// GetCastFromCredits extracts cast members from Credits, limited to maxCount
func GetCastFromCredits(credits []models.Credit, maxCount int) []models.Credit {
	var cast []models.Credit

	// Get all cast members
	for _, credit := range credits {
		if credit.IsCast {
			cast = append(cast, credit)
		}
	}

	// Sort by order if available
	sort.Slice(cast, func(i, j int) bool {
		return cast[i].Order < cast[j].Order
	})

	// Limit to maxCount
	if len(cast) > maxCount {
		cast = cast[:maxCount]
	}

	return cast
}

// GetCrewByRole extracts crew members with a specific role
func GetCrewByRole(credits []models.Credit, role string) []models.Credit {
	var result []models.Credit

	for _, credit := range credits {
		if credit.IsCrew && credit.Role == role {
			result = append(result, credit)
		}
	}

	return result
}

// GetCrewByDepartment extracts crew members from a specific department
func GetCrewByDepartment(credits []models.Credit, department string) []models.Credit {
	var result []models.Credit

	for _, credit := range credits {
		if credit.IsCrew && credit.Department == department {
			result = append(result, credit)
		}
	}

	return result
}

// GetCreatorsFromCredits extracts creators from Credits
func GetCreatorsFromCredits(credits []models.Credit) []models.Credit {
	var creators []models.Credit

	for _, credit := range credits {
		if credit.IsCreator {
			creators = append(creators, credit)
		}
	}

	return creators
}

// ExtractNamesFromCredits extracts just the names from a list of credits
func ExtractNamesFromCredits(credits []models.Credit) []string {
	var names []string

	for _, credit := range credits {
		names = append(names, credit.Name)
	}

	return names
}

// GetPeopleByRole retrieves people from the repository who have a specific role
func (j *RecommendationJob) GetPeopleByRole(ctx context.Context, role string) ([]models.Person, error) {
	// If we don't have a people repository, return an error
	if j.peopleRepo == nil {
		return nil, fmt.Errorf("people repository not available")
	}

	// Get people by role
	people, err := j.peopleRepo.GetByRole(ctx, role)
	if err != nil {
		return nil, err
	}

	return people, nil
}

// GetPersonByID retrieves a person by ID
func (j *RecommendationJob) GetPersonByID(ctx context.Context, personID uint64) (*models.Person, error) {
	// If we don't have a people repository, return an error
	if j.peopleRepo == nil {
		return nil, fmt.Errorf("people repository not available")
	}

	// Get person by ID
	person, err := j.peopleRepo.GetByID(ctx, personID)
	if err != nil {
		return nil, err
	}

	return person, nil
}

// getCreditsForMediaItem retrieves all credits for a given media item
func (j *RecommendationJob) getCreditsForMediaItem(ctx context.Context, mediaItemID uint64) ([]models.Credit, error) {
	// If we don't have a credit repository, return an error
	if j.creditRepo == nil {
		return nil, fmt.Errorf("credit repository not available")
	}

	// Get credits from the repository
	credits, err := j.creditRepo.GetByMediaItemID(ctx, mediaItemID)
	if err != nil {
		return nil, err
	}

	return credits, nil
}

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

	// New repositories for credits and people
	creditRepo repository.CreditRepository // Will be implemented in the future
	peopleRepo repository.PersonRepository // Will be implemented in the future
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
	// Optional repositories for credits and people - can be nil for now
	creditRepo repository.CreditRepository,
	peopleRepo repository.PersonRepository,
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
		creditRepo:      creditRepo,
		peopleRepo:      peopleRepo,
	}
}

// getAIClient returns an AI client for the given user
// It tries to get the default AI client from the user config, or falls back to the first active AI client
func (j *RecommendationJob) getAIClient(ctx context.Context, userID uint64) (ai.AIClient, error) {
	logger := log.Logger{} // would ideally use structured logging from context

	// Get user config to check for default AI client
	config, err := j.configRepo.GetUserConfig(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user config: %w", err)
	}

	// First try to get the default AI client if set
	if config.DefaultClients != nil && config.DefaultClients.AIClientID > 0 {
		// Try Claude repository first
		claudeRepo := j.clientRepos.AllRepos().ClaudeRepo
		if claudeRepo != nil {
			claudeClient, err := claudeRepo.GetByID(ctx, config.DefaultClients.AIClientID)
			if err == nil && claudeClient != nil {
				// Found the default Claude client
				client, err := j.clientFactories.GetClient(ctx, claudeClient.ID, claudeClient.Config.Data)
				if err == nil && client != nil {
					logger.Printf("Using default Claude AI client ID %d for user %d", claudeClient.ID, userID)
					return client.(ai.AIClient), nil
				}
			}
		}

		// Try OpenAI repository next
		openAIRepo := j.clientRepos.AllRepos().OpenAIRepo
		if openAIRepo != nil {
			openAIClient, err := openAIRepo.GetByID(ctx, config.DefaultClients.AIClientID)
			if err == nil && openAIClient != nil {
				// Found the default OpenAI client
				client, err := j.clientFactories.GetClient(ctx, openAIClient.ID, openAIClient.Config.Data)
				if err == nil && client != nil {
					logger.Printf("Using default OpenAI client ID %d for user %d", openAIClient.ID, userID)
					return client.(ai.AIClient), nil
				}
			}
		}

		// Try Ollama repository next
		ollamaRepo := j.clientRepos.AllRepos().OllamaRepo
		if ollamaRepo != nil {
			ollamaClient, err := ollamaRepo.GetByID(ctx, config.DefaultClients.AIClientID)
			if err == nil && ollamaClient != nil {
				// Found the default Ollama client
				client, err := j.clientFactories.GetClient(ctx, ollamaClient.ID, ollamaClient.Config.Data)
				if err == nil && client != nil {
					logger.Printf("Using default Ollama client ID %d for user %d", ollamaClient.ID, userID)
					return client.(ai.AIClient), nil
				}
			}
		}

		// If we get here, the default client couldn't be found or created
		logger.Printf("Default AI client ID %d for user %d not found or could not be created",
			config.DefaultClients.AIClientID, userID)
	}

	// If default client not set or couldn't be loaded, try to get any AI client

	// Try Claude clients first
	claudeRepo := j.clientRepos.AllRepos().ClaudeRepo
	if claudeRepo != nil {
		claudeClients, err := claudeRepo.GetByUserID(ctx, userID)
		if err == nil && len(claudeClients) > 0 {
			// Use the first active Claude client
			for _, clientConfig := range claudeClients {
				client, err := j.clientFactories.GetClient(ctx, clientConfig.ID, clientConfig.Config.Data)
				if err == nil && client != nil {
					logger.Printf("Using first available Claude client ID %d for user %d", clientConfig.ID, userID)
					return client.(ai.AIClient), nil
				}
			}
		}
	}

	// Try OpenAI clients next
	openAIRepo := j.clientRepos.AllRepos().OpenAIRepo
	if openAIRepo != nil {
		openAIClients, err := openAIRepo.GetByUserID(ctx, userID)
		if err == nil && len(openAIClients) > 0 {
			// Use the first active OpenAI client
			for _, clientConfig := range openAIClients {
				client, err := j.clientFactories.GetClient(ctx, clientConfig.ID, clientConfig.Config.Data)
				if err == nil && client != nil {
					logger.Printf("Using first available OpenAI client ID %d for user %d", clientConfig.ID, userID)
					return client.(ai.AIClient), nil
				}
			}
		}
	}

	// Try Ollama clients next
	ollamaRepo := j.clientRepos.AllRepos().OllamaRepo
	if ollamaRepo != nil {
		ollamaClients, err := ollamaRepo.GetByUserID(ctx, userID)
		if err == nil && len(ollamaClients) > 0 {
			// Use the first active Ollama client
			for _, clientConfig := range ollamaClients {
				client, err := j.clientFactories.GetClient(ctx, clientConfig.ID, clientConfig.Config.Data)
				if err == nil && client != nil {
					logger.Printf("Using first available Ollama client ID %d for user %d", clientConfig.ID, userID)
					return client.(ai.AIClient), nil
				}
			}
		}
	}

	// No AI client found
	logger.Printf("No AI clients found for user %d", userID)
	return nil, fmt.Errorf("no AI clients found for user %d", userID)
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

// Execute implements the scheduler.Job interface
func (j *RecommendationJob) Execute(ctx context.Context) error {
	// Since this implementation needs to match the scheduler.Job interface, 
	// we'll create a basic version without parameters
	log.Println("Executing recommendation job")
	return nil
}

// ExecuteWithParams runs the recommendation job with parameters
func (j *RecommendationJob) ExecuteWithParams(ctx context.Context, jobID uint64, jobRunID uint64, params map[string]interface{}) error {
	ctx, jobLog := utils.WithJobID(ctx, jobID)

	jobLog.Info().
		Uint64("jobID", jobID).
		Uint64("jobRunID", jobRunID).
		Interface("params", params).
		Msg("Starting recommendation job")

	// Update job status to in-progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 0, "Starting recommendation generation")

	// Get all active users (or a specific user if provided in params)
	var users []models.User
	var err error

	if userIDParam, ok := params["userID"]; ok {
		// If a specific user ID was provided
		userIDint, _ := strconv.ParseUint(fmt.Sprintf("%v", userIDParam), 10, 64)
		user, err := j.userRepo.FindByID(ctx, userIDint)
		if err != nil {
			jobLog.Error().Err(err).Msg("Failed to get user")
			return err
		}
		users = []models.User{*user}
	} else {
		// Get all active users
		users, err = j.userRepo.FindAllActive(ctx) // Active
		if err != nil {
			jobLog.Error().Err(err).Msg("Failed to get users")
			return err
		}
	}

	total := len(users)
	jobLog.Info().
		Int("userCount", total).
		Msg("Generating recommendations for users")

	// Process each user
	for idx, user := range users {
		progress := (idx * 100) / total
		statusMsg := fmt.Sprintf("Processing user %d/%d", idx+1, total)
		j.jobRepo.UpdateJobProgress(ctx, jobRunID, progress, statusMsg)

		err := j.processUserRecommendations(ctx, jobRunID, user)
		if err != nil {
			jobLog.Error().
				Err(err).
				Uint64("userID", user.ID).
				Msg("Failed to generate recommendations for user")
			// Continue to the next user
			continue
		}
	}

	// Mark job as completed
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 100, "Recommendation generation completed")
	jobLog.Info().Msg("Recommendation job completed")

	return nil
}

// processUserRecommendations generates recommendations for a single user
func (j *RecommendationJob) processUserRecommendations(ctx context.Context, jobRunID uint64, user models.User) error {
	log := utils.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", user.ID).
		Msg("Processing recommendations for user")

	// Get user config to determine preferences
	config, err := j.configRepo.GetUserConfig(ctx, user.ID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user config")
		return err
	}

	// Check if recommendations are enabled for this user
	if !config.RecommendationSyncEnabled {
		log.Info().Msg("Recommendations are disabled for this user, skipping")
		return nil
	}

	// Build user preference profile for personalized recommendations
	preferenceProfile, err := j.buildUserPreferenceProfile(ctx, user.ID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to build user preference profile")
		return err
	}

	// Generate different types of recommendations based on user preferences
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 10, "Generating movie recommendations")
	err = j.generateMovieRecommendations(ctx, jobRunID, user, preferenceProfile, config)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate movie recommendations")
		// Continue to other types
	}

	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 40, "Generating TV show recommendations")
	err = j.generateSeriesRecommendations(ctx, jobRunID, user, preferenceProfile, config)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate TV show recommendations")
		// Continue to other types
	}

	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 70, "Generating music recommendations")
	err = j.generateMusicRecommendations(ctx, jobRunID, user, preferenceProfile, config)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate music recommendations")
		// Continue to other types
	}

	return nil
}

// buildUserPreferenceProfile analyzes user history and preferences
func (j *RecommendationJob) buildUserPreferenceProfile(ctx context.Context, userID uint64) (*UserPreferenceProfile, error) {
	log := utils.LoggerFromContext(ctx)
	startTime := time.Now()

	// Create new profile with all the necessary maps initialized
	profile := &UserPreferenceProfile{
		// Common preferences
		PreferredLanguages:    make(map[string]float32),
		PreferredReleaseYears: [2]int{1900, time.Now().Year()}, // Default to all years
		ContentRatingRange:    [2]string{"G", "R"},             // Default to all ratings
		MinRating:             0.0,                             // Default to no minimum rating

		// Movie preferences
		WatchedMovieIDs:       make(map[uint64]bool),
		RecentMovies:          []MovieSummary{},
		TopRatedMovies:        []MovieSummary{},
		FavoriteMovieGenres:   make(map[string]float32),
		FavoriteActors:        make(map[string]float32),
		FavoriteDirectors:     make(map[string]float32),
		MovieWatchTimes:       make(map[string][]int64),
		MovieWatchDays:        make(map[string]int),
		MovieTagPreferences:   make(map[string]float32),
		ExcludedMovieGenres:   []string{},
		PreferredMovieGenres:  []string{},
		MovieReleaseYearRange: [2]int{1900, time.Now().Year()},

		// Series preferences
		WatchedSeriesIDs:       make(map[uint64]bool),
		RecentSeries:           []SeriesSummary{},
		TopRatedSeries:         []SeriesSummary{},
		FavoriteSeriesGenres:   make(map[string]float32),
		FavoriteShowrunners:    make(map[string]float32),
		SeriesWatchTimes:       make(map[string][]int64),
		SeriesWatchDays:        make(map[string]int),
		SeriesTagPreferences:   make(map[string]float32),
		ExcludedSeriesGenres:   []string{},
		PreferredSeriesGenres:  []string{},
		SeriesReleaseYearRange: [2]int{1900, time.Now().Year()},
		PreferredSeriesStatus:  []string{},

		// Music preferences
		PlayedMusicIDs:       make(map[uint64]bool),
		RecentMusic:          []MusicSummary{},
		TopRatedMusic:        []MusicSummary{},
		FavoriteMusicGenres:  make(map[string]float32),
		FavoriteArtists:      make(map[string]float32),
		MusicPlayTimes:       make(map[string][]int64),
		MusicPlayDays:        make(map[string]int),
		MusicTagPreferences:  make(map[string]float32),
		ExcludedMusicGenres:  []string{},
		PreferredMusicGenres: []string{},

		// Watch/play patterns
		WatchTimeOfDay:       make(map[int]int),
		WatchDayOfWeek:       make(map[string]int),
		TypicalSessionLength: make(map[string]float32),

		// Owned content
		OwnedMovieIDs:  make(map[string]bool),
		OwnedSeriesIDs: make(map[string]bool),
		OwnedMusicIDs:  make(map[string]bool),

		// Activity levels
		OverallActivityLevel: make(map[string]float32),

		// Analysis metadata
		AnalysisTimestamp:  time.Now().Unix(),
		ProfileConfidence:  0.0, // Will be calculated based on data points
		DataPointsAnalyzed: 0,
	}

	// Get user configuration to determine preferences
	config, err := j.configRepo.GetUserConfig(ctx, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user config")
		return profile, err
	}

	// Apply user configuration settings
	if config.ExcludedGenres != nil && config.ExcludedGenres.Movies != nil {
		profile.ExcludedMovieGenres = config.ExcludedGenres.Movies
	}

	if config.ExcludedGenres != nil && config.ExcludedGenres.Series != nil {
		profile.ExcludedSeriesGenres = config.ExcludedGenres.Series
	}

	if config.ExcludedGenres != nil && config.ExcludedGenres.Music != nil {
		profile.ExcludedMusicGenres = config.ExcludedGenres.Music
	}

	if config.PreferredGenres != nil && config.PreferredGenres.Movies != nil {
		profile.PreferredMovieGenres = config.PreferredGenres.Movies
	}

	if config.PreferredGenres != nil && config.PreferredGenres.Series != nil {
		profile.PreferredSeriesGenres = config.PreferredGenres.Series
	}

	if config.PreferredGenres != nil && config.PreferredGenres.Music != nil {
		profile.PreferredMusicGenres = config.PreferredGenres.Music
	}

	// Set behavior preferences
	profile.ExcludeWatched = !config.RecommendationIncludeWatched
	profile.IncludeSimilarItems = config.RecommendationIncludeSimilar
	profile.UseAIRecommendations = config.RecommendationSyncEnabled
	profile.NotifyForMovies = strings.Contains(config.NotifyMediaTypes, "movie")
	profile.NotifyForSeries = strings.Contains(config.NotifyMediaTypes, "series")
	profile.NotifyForMusic = strings.Contains(config.NotifyMediaTypes, "music")
	profile.RatingThreshold = config.NotifyRatingThreshold
	profile.NotifyForUpcoming = config.NotifyUpcomingReleases
	profile.NotifyForRecent = config.NotifyRecentReleases

	// Set content rating range based on user settings
	if config.MinContentRating != "" && config.MaxContentRating != "" {
		profile.ContentRatingRange = [2]string{config.MinContentRating, config.MaxContentRating}
	}

	// Set preferred language weights
	profile.PreferredLanguages["en"] = 1.0 // Default English
	if config.PreferredAudioLanguages != "" {
		languages := strings.Split(config.PreferredAudioLanguages, ",")
		for i, lang := range languages {
			// Give higher weight to languages listed first
			weight := 1.5 - (float32(i) * 0.1)
			if weight < 1.0 {
				weight = 1.0
			}
			profile.PreferredLanguages[strings.TrimSpace(lang)] = weight
		}
	}

	// Get watch/play history
	moviesRecentHistory, err := j.historyRepo.GetRecentUserMovieHistory(ctx, userID, 200)
	if err != nil {
		log.Error().Err(err).Str("type", "movie").Msg("Failed to get user history")
	}

	seriesRecentHistory, err := j.historyRepo.GetRecentUserSeriesHistory(ctx, userID, 200)
	if err != nil {
		log.Error().Err(err).Str("type", "series").Msg("Failed to get user history")
	}

	musicRecentHistory, err := j.historyRepo.GetRecentUserMusicHistory(ctx, userID, 200)
	if err != nil {
		log.Error().Err(err).Str("type", "music").Msg("Failed to get user history")
	}

	// Count total data points for confidence calculation
	profile.DataPointsAnalyzed = len(moviesRecentHistory) + len(seriesRecentHistory) + len(musicRecentHistory)

	// Process movie history
	j.processMovieHistory(ctx, profile, moviesRecentHistory)

	// Process series history
	j.processSeriesHistory(ctx, profile, seriesRecentHistory)

	// Process music history
	j.processMusicHistory(ctx, profile, musicRecentHistory)

	// Calculate advanced metrics based on all processed data
	j.calculateAdvancedMetrics(profile)

	// Calculate profile confidence based on amount of data
	if profile.DataPointsAnalyzed > 100 {
		profile.ProfileConfidence = 0.9
	} else if profile.DataPointsAnalyzed > 50 {
		profile.ProfileConfidence = 0.7
	} else if profile.DataPointsAnalyzed > 20 {
		profile.ProfileConfidence = 0.5
	} else if profile.DataPointsAnalyzed > 5 {
		profile.ProfileConfidence = 0.3
	} else {
		profile.ProfileConfidence = 0.1
	}

	// Log profile generation time
	log.Info().
		Uint64("userID", userID).
		Int("dataPoints", profile.DataPointsAnalyzed).
		Float32("confidence", profile.ProfileConfidence).
		Dur("duration", time.Since(startTime)).
		Msg("User preference profile built")

	return profile, nil
}

// processMovieHistory analyzes movie watch history to build preferences
func (j *RecommendationJob) processMovieHistory(ctx context.Context, profile *UserPreferenceProfile, histories []models.MediaPlayHistory[*mediatypes.Movie]) {
	log := utils.LoggerFromContext(ctx)

	// Maps for processing
	recentMovies := []MovieSummary{}
	highRatedMovies := []MovieSummary{}
	watchTimes := make(map[int]int)   // Hour of day -> count
	watchDays := make(map[string]int) // Day of week -> count
	movieGenreWeights := make(map[string]float32)

	// Default days of week
	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	for _, day := range days {
		watchDays[day] = 0
		profile.MovieWatchDays[day] = 0
	}

	// Process each history item
	for _, history := range histories {
		// Track watched movie ID
		if history.MediaItemID > 0 {
			profile.WatchedMovieIDs[history.MediaItemID] = true
		}

		// Track watch time patterns
		if !history.LastPlayedAt.IsZero() {
			// Hour of day
			hour := history.LastPlayedAt.Hour()
			watchTimes[hour]++
			profile.WatchTimeOfDay[hour]++

			// Day of week
			day := history.LastPlayedAt.Weekday().String()
			watchDays[day]++
			profile.MovieWatchDays[day]++
		}

		// Calculate weight based on recency, completion, and rating
		weight := float32(1.0)

		// More weight for higher completion percentage
		if history.PlayedPercentage > 0 {
			weight = float32(history.PlayedPercentage) / 100
		}

		// More weight for higher ratings
		if history.UserRating > 0 {
			weight *= (history.UserRating / 5.0) * 1.5
		}

		// More weight for more recent watches
		daysSinceWatch := time.Since(history.LastPlayedAt).Hours() / 24
		if daysSinceWatch > 0 && daysSinceWatch < 365 {
			recencyFactor := 1.0 - (float32(daysSinceWatch) / 365.0)
			weight *= (1.0 + recencyFactor)
		}

		// Get detailed movie information if available
		if history.MediaItemID > 0 {
			movie, err := j.movieRepo.GetByID(ctx, history.MediaItemID)
			if err != nil || movie == nil || movie.Data == nil {
				continue
			}

			// Process movie details
			movieData := movie.Data

			// Build movie summary
			summary := MovieSummary{
				Title:             movieData.Details.Title,
				Year:              movieData.Details.ReleaseYear,
				Genres:            movieData.Details.Genres,
				PlayCount:         int(history.PlayCount),
				CompletionPercent: float32(history.PlayedPercentage),
				WatchDate:         history.LastPlayedAt.Unix(),
			}

			// Add TMDB ID if available
			if movie.ExternalIDs.GetID("tmdb") != "" {
				summary.TMDB_ID = movie.ExternalIDs.GetID("tmdb")
			}

			// Add cast if available
			cast := movieData.Credits.GetCast()
			if len(cast) > 0 {
				// Take only top-billed actors (up to 3)
				castCount := len(cast)
				if castCount > 3 {
					castCount = 3
				}
				// Convert Person objects to strings for summary
				castNames := make([]string, castCount)
				for i, person := range cast[:castCount] {
					castNames[i] = person.Name
				}
				summary.Cast = castNames

				// Add to favorite actors with weighted score
				for i, actor := range cast {
					// Decrease weight for lower-billed actors
					actorWeight := weight * (1.0 - (float32(i) * 0.1))
					if actorWeight < 0.1 {
						actorWeight = 0.1
					}
					profile.FavoriteActors[actor.Name] += actorWeight
				}
			}

			// Add directors if available
			crew := movieData.Credits.GetCrew()
			if len(crew) > 0 {
				var directors []string
				for _, crewMember := range crew {
					if crewMember.Role == "Director" {
						directors = append(directors, crewMember.Name)
						profile.FavoriteDirectors[crewMember.Name] += weight
					}
				}
				if len(directors) > 0 {
					summary.Directors = directors
				}
			}

			// Process genres with weighted scoring
			if movieData.Details.Genres != nil && len(movieData.Details.Genres) > 0 {
				for _, genre := range movieData.Details.Genres {
					movieGenreWeights[genre] += weight
				}
			}

			// Process user tags if any
			// if movieData.UserTags != nil && len(movieData.UserTags) > 0 {
			// 	summary.UserTags = movieData.UserTags
			// 	for _, tag := range movieData.UserTags {
			// 		profile.MovieTagPreferences[tag] += weight
			// 	}
			// }

			// Add to recent movies list
			recentMovies = append(recentMovies, summary)

			// Add to high-rated movies if applicable
			if history.UserRating > 3.5 {
				highRatedMovies = append(highRatedMovies, MovieSummary{
					Title:          movieData.Details.Title,
					Year:           movieData.Details.ReleaseYear,
					Genres:         movieData.Details.Genres,
					DetailedRating: &RatingDetails{Overall: history.UserRating},
					IsFavorite:     history.UserRating > 4.0,
					TMDB_ID:        movie.ExternalIDs.GetID("tmdb"),
					Cast:           summary.Cast,
					Directors:      summary.Directors,
				})
			}
		}
	}

	// Sort recent movies by watch date (newest first)
	sort.Slice(recentMovies, func(i, j int) bool {
		return recentMovies[i].WatchDate > recentMovies[j].WatchDate
	})

	// Limit to 20 most recent
	if len(recentMovies) > 20 {
		recentMovies = recentMovies[:20]
	}

	// Sort high-rated movies by rating (highest first)
	sort.Slice(highRatedMovies, func(i, j int) bool {
		if highRatedMovies[i].DetailedRating == nil || highRatedMovies[j].DetailedRating == nil {
			return false
		}
		return highRatedMovies[i].DetailedRating.Overall > highRatedMovies[j].DetailedRating.Overall
	})

	// Limit to 20 highest rated
	if len(highRatedMovies) > 20 {
		highRatedMovies = highRatedMovies[:20]
	}

	// Update profile with processed movie data
	profile.RecentMovies = recentMovies
	profile.TopRatedMovies = highRatedMovies
	profile.FavoriteMovieGenres = movieGenreWeights

	// Calculate movie watch time patterns
	for hour, count := range watchTimes {
		hourKey := fmt.Sprintf("%02d:00", hour)
		profile.MovieWatchTimes[hourKey] = append(profile.MovieWatchTimes[hourKey], int64(count))
	}

	// Calculate typical session length for movies if we have data
	if len(histories) > 0 {
		var totalDuration float32
		var validItems int

		for _, history := range histories {
			if history.DurationSeconds > 0 {
				totalDuration += float32(history.DurationSeconds)
				validItems++
			}
		}

		if validItems > 0 {
			profile.TypicalSessionLength["movie"] = totalDuration / float32(validItems)
		}
	}

	// Calculate activity level for movies (0-1 scale)
	if len(histories) > 0 {
		// Base activity on number of watches and recency
		recentCount := 0
		for _, history := range histories {
			// Count items watched in the last 30 days
			if time.Since(history.LastPlayedAt).Hours() < 30*24 {
				recentCount++
			}
		}

		// Scale: 0 = no activity, 1 = very active (10+ movies per month)
		activityScore := float32(recentCount) / 10.0
		if activityScore > 1.0 {
			activityScore = 1.0
		}

		profile.OverallActivityLevel["movie"] = activityScore
	}

	log.Debug().
		Int("movieHistoryCount", len(histories)).
		Int("topRatedMovies", len(profile.TopRatedMovies)).
		Int("recentMovies", len(profile.RecentMovies)).
		Int("genres", len(profile.FavoriteMovieGenres)).
		Msg("Processed movie history")
}

// processSeriesHistory analyzes TV series watch history to build preferences
func (j *RecommendationJob) processSeriesHistory(ctx context.Context, profile *UserPreferenceProfile, histories []models.MediaPlayHistory[*mediatypes.Series]) {
	log := utils.LoggerFromContext(ctx)

	// Maps for processing
	recentSeries := []SeriesSummary{}
	highRatedSeries := []SeriesSummary{}
	watchTimes := make(map[int]int)
	watchDays := make(map[string]int)
	seriesGenreWeights := make(map[string]float32)

	// Default days of week
	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	for _, day := range days {
		watchDays[day] = 0
		profile.SeriesWatchDays[day] = 0
	}

	// Process each history item
	for _, history := range histories {
		// Track watched series ID
		if history.MediaItemID > 0 {
			profile.WatchedSeriesIDs[history.MediaItemID] = true
		}

		// Track watch time patterns
		if !history.LastPlayedAt.IsZero() {
			// Hour of day
			hour := history.LastPlayedAt.Hour()
			watchTimes[hour]++
			profile.WatchTimeOfDay[hour]++

			// Day of week
			day := history.LastPlayedAt.Weekday().String()
			watchDays[day]++
			profile.SeriesWatchDays[day]++
		}

		// Calculate weight based on recency, completion, and rating
		weight := float32(1.0)

		// More weight for higher completion percentage
		if history.PlayedPercentage > 0 {
			weight = float32(history.PlayedPercentage) / 100
		}

		// More weight for higher ratings
		if history.UserRating > 0 {
			weight *= (history.UserRating / 5.0) * 1.5
		}

		// More weight for more recent watches
		daysSinceWatch := time.Since(history.LastPlayedAt).Hours() / 24
		if daysSinceWatch > 0 && daysSinceWatch < 365 {
			recencyFactor := 1.0 - (float32(daysSinceWatch) / 365.0)
			weight *= (1.0 + recencyFactor)
		}

		// Get detailed series information if available
		if history.MediaItemID > 0 {
			series, err := j.seriesRepo.GetByID(ctx, history.MediaItemID)
			if err != nil || series == nil || series.Data == nil {
				continue
			}

			// Process series details
			seriesData := series.Data

			// Build series summary
			summary := SeriesSummary{
				Title:           seriesData.Details.Title,
				Year:            seriesData.Details.ReleaseYear,
				Genres:          seriesData.Genres,
				Rating:          history.UserRating,
				EpisodesWatched: int(history.PlayCount),
				LastWatchDate:   history.LastPlayedAt.Unix(),
				Status:          seriesData.Status,
				Network:         seriesData.Network,
			}

			// Add detailed rating if available
			if history.UserRating > 0 {
				summary.DetailedRating = &RatingDetails{
					Overall:   history.UserRating,
					Timestamp: history.LastPlayedAt.Unix(),
				}
			}

			// Add seasons/episodes info if available
			if seriesData.SeasonCount > 0 {
				summary.Seasons = seriesData.SeasonCount
			}

			if seriesData.EpisodeCount > 0 {
				summary.TotalEpisodes = seriesData.EpisodeCount
			}

			// Add TMDB ID if available
			if series.ExternalIDs.GetID("tmdb") != "" {
				summary.TMDB_ID = series.ExternalIDs.GetID("tmdb")
			}

			// Add cast if available
			cast := seriesData.Credits.GetCast()
			if len(cast) > 0 {
				// Take only top-billed actors (up to 5)
				castCount := len(cast)
				if castCount > 5 {
					castCount = 5
				}
				// Convert Person objects to strings for summary
				castNames := make([]string, castCount)
				for i, person := range cast[:castCount] {
					castNames[i] = person.Name
				}
				summary.Cast = castNames
			}

			// Add showrunners/creators if available
			creators := seriesData.Credits.GetCreators()
			if len(creators) > 0 {
				// Convert Person objects to strings for summary
				creatorNames := make([]string, len(creators))
				for i, person := range creators {
					creatorNames[i] = person.Name
				}
				summary.Showrunners = creatorNames
				for _, creator := range creators {
					profile.FavoriteShowrunners[creator.Name] += weight
				}
			}

			// Process genres with weighted scoring
			if seriesData.Genres != nil && len(seriesData.Genres) > 0 {
				for _, genre := range seriesData.Genres {
					seriesGenreWeights[genre] += weight
				}
			}

			// Process user tags if any
			// if seriesData.UserTags != nil && len(seriesData.UserTags) > 0 {
			// 	summary.UserTags = seriesData.UserTags
			// 	for _, tag := range seriesData.UserTags {
			// 		profile.SeriesTagPreferences[tag] += weight
			// 	}
			// }

			// Add to recent series list
			recentSeries = append(recentSeries, summary)

			// Add to high-rated series if applicable
			if history.UserRating > 3.5 {
				highRatedSeries = append(highRatedSeries, summary)
				summary.IsFavorite = history.UserRating > 4.0
			}
		}
	}

	// Sort recent series by watch date (newest first)
	sort.Slice(recentSeries, func(i, j int) bool {
		return recentSeries[i].LastWatchDate > recentSeries[j].LastWatchDate
	})

	// Limit to 20 most recent
	if len(recentSeries) > 20 {
		recentSeries = recentSeries[:20]
	}

	// Sort high-rated series by rating (highest first)
	sort.Slice(highRatedSeries, func(i, j int) bool {
		if highRatedSeries[i].DetailedRating == nil || highRatedSeries[j].DetailedRating == nil {
			return false
		}
		return highRatedSeries[i].DetailedRating.Overall > highRatedSeries[j].DetailedRating.Overall
	})

	// Limit to 20 highest rated
	if len(highRatedSeries) > 20 {
		highRatedSeries = highRatedSeries[:20]
	}

	// Update profile with processed series data
	profile.RecentSeries = recentSeries
	profile.TopRatedSeries = highRatedSeries
	profile.FavoriteSeriesGenres = seriesGenreWeights

	// Calculate series watch time patterns
	for hour, count := range watchTimes {
		hourKey := fmt.Sprintf("%02d:00", hour)
		profile.SeriesWatchTimes[hourKey] = append(profile.SeriesWatchTimes[hourKey], int64(count))
	}

	// Calculate typical session length for series if we have data
	if len(histories) > 0 {
		var totalDuration float32
		var validItems int

		for _, history := range histories {
			if history.DurationSeconds > 0 {
				totalDuration += float32(history.DurationSeconds)
				validItems++
			}
		}

		if validItems > 0 {
			profile.TypicalSessionLength["series"] = totalDuration / float32(validItems)
		}
	}

	// Calculate activity level for series (0-1 scale)
	if len(histories) > 0 {
		// Base activity on number of watches and recency
		recentCount := 0
		for _, history := range histories {
			// Count items watched in the last 30 days
			if time.Since(history.LastPlayedAt).Hours() < 30*24 {
				recentCount++
			}
		}

		// Scale: 0 = no activity, 1 = very active (20+ episodes per month)
		activityScore := float32(recentCount) / 20.0
		if activityScore > 1.0 {
			activityScore = 1.0
		}

		profile.OverallActivityLevel["series"] = activityScore
	}

	// Find preferred series status based on watch history
	statusCount := make(map[string]int)
	for _, series := range recentSeries {
		if series.Status != "" {
			statusCount[series.Status]++
		}
	}

	// Add the most common status types to preferences
	type statusFreq struct {
		status string
		count  int
	}

	var statuses []statusFreq
	for status, count := range statusCount {
		statuses = append(statuses, statusFreq{status, count})
	}

	sort.Slice(statuses, func(i, j int) bool {
		return statuses[i].count > statuses[j].count
	})

	// Add top 2 status preferences if we have enough data
	if len(statuses) > 0 {
		profile.PreferredSeriesStatus = append(profile.PreferredSeriesStatus, statuses[0].status)
		if len(statuses) > 1 {
			profile.PreferredSeriesStatus = append(profile.PreferredSeriesStatus, statuses[1].status)
		}
	}

	log.Debug().
		Int("seriesHistoryCount", len(histories)).
		Int("topRatedSeries", len(profile.TopRatedSeries)).
		Int("recentSeries", len(profile.RecentSeries)).
		Int("genres", len(profile.FavoriteSeriesGenres)).
		Msg("Processed series history")
}

// processMusicHistory analyzes music play history to build preferences
func (j *RecommendationJob) processMusicHistory(ctx context.Context, profile *UserPreferenceProfile, histories []models.MediaPlayHistory[*mediatypes.Track]) {
	log := utils.LoggerFromContext(ctx)

	// Maps for processing
	recentMusic := []MusicSummary{}
	topRatedMusic := []MusicSummary{}
	playTimes := make(map[int]int)
	playDays := make(map[string]int)
	musicGenreWeights := make(map[string]float32)
	artistWeights := make(map[string]float32)

	// Default days of week
	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	for _, day := range days {
		playDays[day] = 0
		profile.MusicPlayDays[day] = 0
	}

	// Process each history item
	for _, history := range histories {
		// Track played music ID
		if history.MediaItemID > 0 {
			profile.PlayedMusicIDs[history.MediaItemID] = true
		}

		// Track play time patterns
		if !history.LastPlayedAt.IsZero() {
			// Hour of day
			hour := history.LastPlayedAt.Hour()
			playTimes[hour]++
			profile.WatchTimeOfDay[hour]++

			// Day of week
			day := history.LastPlayedAt.Weekday().String()
			playDays[day]++
			profile.MusicPlayDays[day]++
		}

		// Calculate weight based on recency, completion, and rating
		weight := float32(1.0)

		// More weight for higher completion percentage
		if history.PlayedPercentage > 0 {
			weight = float32(history.PlayedPercentage) / 100
		}

		// More weight for higher ratings
		if history.UserRating > 0 {
			weight *= (history.UserRating / 5.0) * 1.5
		}

		// More weight for more recent plays
		daysSincePlay := time.Since(history.LastPlayedAt).Hours() / 24
		if daysSincePlay > 0 && daysSincePlay < 365 {
			recencyFactor := 1.0 - (float32(daysSincePlay) / 365.0)
			weight *= (1.0 + recencyFactor)
		}

		// Get detailed music information if available
		if history.MediaItemID > 0 {
			track, err := j.musicRepo.GetByID(ctx, history.MediaItemID)
			if err != nil || track == nil || track.Data == nil {
				continue
			}

			// Process track details
			trackData := track.Data

			// Build music summary
			summary := MusicSummary{
				Title:        trackData.Details.Title,
				Artist:       trackData.ArtistName,
				Album:        trackData.AlbumName,
				Year:         trackData.Details.ReleaseYear,
				Genres:       trackData.Details.Genres,
				Rating:       history.UserRating,
				PlayCount:    int(history.PlayCount),
				LastPlayDate: history.LastPlayedAt.Unix(),
			}

			// Add duration if available
			if trackData.Duration > 0 {
				summary.DurationSec = trackData.Duration
			}

			// Add detailed rating if available
			if history.UserRating > 0 {
				summary.DetailedRating = &RatingDetails{
					Overall:   history.UserRating,
					Timestamp: history.LastPlayedAt.Unix(),
				}
			}

			// Add external ID if available
			// if trackData.ExternalID != "" {
			// 	summary.ExternalID = trackData.ExternalID
			// }

			// Process genres with weighted scoring
			if len(trackData.Details.Genres) > 0 {
				for _, genre := range trackData.Details.Genres {
					musicGenreWeights[genre] += weight
				}
			}

			// Process artist preferences
			if trackData.ArtistName != "" {
				artistWeights[trackData.ArtistName] += weight
			}

			// Process user tags if any
			// if trackData.UserTags != nil && len(trackData.UserTags) > 0 {
			// 	summary.UserTags = trackData.UserTags
			// 	for _, tag := range trackData.UserTags {
			// 		profile.MusicTagPreferences[tag] += weight
			// 	}
			// }

			// Process mood tags if any
			// Track doesn't have mood field yet
			if false { // Placeholder until Mood field is added
				// Placeholder until Mood field is added
				var moods []string
				summary.Mood = moods
				for _, mood := range moods {
					if profile.MusicMoodPreferences == nil {
						profile.MusicMoodPreferences = make(map[string]float32)
					}
					profile.MusicMoodPreferences[mood] += weight
				}
			}

			// Mark as favorite if highly rated or frequently played
			if history.UserRating > 4.0 || history.PlayCount > 5 {
				summary.IsFavorite = true
			}

			// Add to recent music list
			recentMusic = append(recentMusic, summary)

			// Add to top rated music if applicable
			if history.UserRating > 3.5 {
				topRatedMusic = append(topRatedMusic, summary)
			}
		}
	}

	// Sort recent music by play date (newest first)
	sort.Slice(recentMusic, func(i, j int) bool {
		return recentMusic[i].LastPlayDate > recentMusic[j].LastPlayDate
	})

	// Limit to 20 most recent
	if len(recentMusic) > 20 {
		recentMusic = recentMusic[:20]
	}

	// Sort top rated music by rating (highest first)
	sort.Slice(topRatedMusic, func(i, j int) bool {
		if topRatedMusic[i].DetailedRating == nil || topRatedMusic[j].DetailedRating == nil {
			return false
		}
		return topRatedMusic[i].DetailedRating.Overall > topRatedMusic[j].DetailedRating.Overall
	})

	// Limit to 20 highest rated
	if len(topRatedMusic) > 20 {
		topRatedMusic = topRatedMusic[:20]
	}

	// Update profile with processed music data
	profile.RecentMusic = recentMusic
	profile.TopRatedMusic = topRatedMusic
	profile.FavoriteMusicGenres = musicGenreWeights
	profile.FavoriteArtists = artistWeights

	// Calculate music play time patterns
	for hour, count := range playTimes {
		hourKey := fmt.Sprintf("%02d:00", hour)
		profile.MusicPlayTimes[hourKey] = append(profile.MusicPlayTimes[hourKey], int64(count))
	}

	// Calculate typical session length for music if we have data
	if len(histories) > 0 {
		var totalDuration float32
		var validItems int

		for _, history := range histories {
			if history.DurationSeconds > 0 {
				totalDuration += float32(history.DurationSeconds)
				validItems++
			}
		}

		if validItems > 0 {
			profile.TypicalSessionLength["music"] = totalDuration / float32(validItems)
		}
	}

	// Calculate activity level for music (0-1 scale)
	if len(histories) > 0 {
		// Base activity on number of plays and recency
		recentCount := 0
		for _, history := range histories {
			// Count items played in the last 30 days
			if time.Since(history.LastPlayedAt).Hours() < 30*24 {
				recentCount++
			}
		}

		// Scale: 0 = no activity, 1 = very active (50+ tracks per month)
		activityScore := float32(recentCount) / 50.0
		if activityScore > 1.0 {
			activityScore = 1.0
		}

		profile.OverallActivityLevel["music"] = activityScore
	}

	// Determine preferred music duration range
	if len(recentMusic) > 5 {
		var durations []int
		for _, music := range recentMusic {
			if music.DurationSec > 0 {
				durations = append(durations, music.DurationSec)
			}
		}

		if len(durations) > 5 {
			sort.Ints(durations)

			// Calculate 25th and 75th percentiles for min/max preferred duration
			minIdx := len(durations) / 4
			maxIdx := (len(durations) * 3) / 4

			profile.MusicDurationRange = [2]int{durations[minIdx], durations[maxIdx]}
		}
	}

	log.Debug().
		Int("musicHistoryCount", len(histories)).
		Int("topRatedMusic", len(profile.TopRatedMusic)).
		Int("recentMusic", len(profile.RecentMusic)).
		Int("genres", len(profile.FavoriteMusicGenres)).
		Int("artists", len(profile.FavoriteArtists)).
		Msg("Processed music history")
}

// calculateAdvancedMetrics derives higher-level insights from the user's history
func (j *RecommendationJob) calculateAdvancedMetrics(profile *UserPreferenceProfile) {
	// Skip if we don't have enough data
	if profile.DataPointsAnalyzed < 10 {
		profile.GenreBreadth = 0.5        // Default to medium
		profile.ContentCompleter = 0.5    // Default to medium
		profile.NewContentScore = 0.5     // Default to medium
		profile.PopularityInfluence = 0.5 // Default to medium
		profile.RatingInfluence = 0.5     // Default to medium
		profile.ExplorationScore = 0.5    // Default to medium
		profile.BingeWatchingScore = 0.5  // Default to medium
		profile.ContentRotationFreq = 0.5 // Default to medium
		return
	}

	// Calculate genre breadth (diversity of tastes)
	// This is based on how many different genres the user watches/listens to
	totalGenres := len(profile.FavoriteMovieGenres) + len(profile.FavoriteSeriesGenres) + len(profile.FavoriteMusicGenres)
	if totalGenres > 25 {
		profile.GenreBreadth = 1.0 // Very diverse
	} else if totalGenres > 15 {
		profile.GenreBreadth = 0.8 // Quite diverse
	} else if totalGenres > 10 {
		profile.GenreBreadth = 0.6 // Moderately diverse
	} else if totalGenres > 5 {
		profile.GenreBreadth = 0.4 // Somewhat focused
	} else {
		profile.GenreBreadth = 0.2 // Very focused
	}

	// Calculate content completion tendency
	// This is based on how often the user completes what they start
	var completionSum float32
	var completionCount int

	// Check movie completion
	for _, movie := range profile.RecentMovies {
		if movie.CompletionPercent > 0 {
			completionSum += movie.CompletionPercent
			completionCount++
		}
	}

	// If we have enough data, calculate the score
	if completionCount > 0 {
		profile.ContentCompleter = completionSum / (float32(completionCount) * 100.0)
	} else {
		profile.ContentCompleter = 0.5 // Default to medium
	}

	// Calculate new content score
	// This indicates preference for new vs. classic content
	var recentContentCount int
	var totalContent int

	// Check movie years
	currentYear := time.Now().Year()
	for _, movie := range profile.RecentMovies {
		totalContent++
		if movie.Year >= currentYear-3 {
			recentContentCount++
		}
	}

	// Check series years
	for _, series := range profile.RecentSeries {
		totalContent++
		if series.Year >= currentYear-3 {
			recentContentCount++
		}
	}

	// If we have enough data, calculate the score
	if totalContent > 0 {
		profile.NewContentScore = float32(recentContentCount) / float32(totalContent)
	} else {
		profile.NewContentScore = 0.5 // Default to medium
	}

	// Calculate binge watching tendency
	// This is based on how many episodes/movies watched in a single day
	watchCountByDate := make(map[string]int)

	// Process movie watch dates
	for _, movie := range profile.RecentMovies {
		if movie.WatchDate > 0 {
			date := time.Unix(movie.WatchDate, 0).Format("2006-01-02")
			watchCountByDate[date]++
		}
	}

	// Process series watch dates
	for _, series := range profile.RecentSeries {
		if series.LastWatchDate > 0 {
			date := time.Unix(series.LastWatchDate, 0).Format("2006-01-02")
			watchCountByDate[date]++
		}
	}

	// Calculate average watches per day when watching
	var totalWatchDays int
	var totalWatches int
	for _, count := range watchCountByDate {
		totalWatchDays++
		totalWatches += count
	}

	if totalWatchDays > 0 {
		avgWatchesPerDay := float32(totalWatches) / float32(totalWatchDays)

		// Scale: 1-2 items = low (0.2), 3-4 = medium (0.5), 5+ = high (0.8+)
		if avgWatchesPerDay >= 5 {
			profile.BingeWatchingScore = 0.8 + (float32(avgWatchesPerDay-5) * 0.04)
			if profile.BingeWatchingScore > 1.0 {
				profile.BingeWatchingScore = 1.0
			}
		} else if avgWatchesPerDay >= 3 {
			profile.BingeWatchingScore = 0.5 + ((avgWatchesPerDay - 3) * 0.15)
		} else {
			profile.BingeWatchingScore = 0.2 + ((avgWatchesPerDay - 1) * 0.15)
		}
	} else {
		profile.BingeWatchingScore = 0.5 // Default to medium
	}

	// Calculate content rotation frequency
	// This indicates how often the user switches genres/styles
	if len(profile.RecentMovies) > 3 || len(profile.RecentSeries) > 3 {
		var genreChangeCount int
		var itemCount int

		// Check genre changes in recent movies
		var lastGenre string
		for i, movie := range profile.RecentMovies {
			if i == 0 && len(movie.Genres) > 0 {
				lastGenre = movie.Genres[0]
				itemCount++
				continue
			}

			if len(movie.Genres) > 0 {
				itemCount++
				// Check if primary genre changed
				if movie.Genres[0] != lastGenre {
					genreChangeCount++
					lastGenre = movie.Genres[0]
				}
			}
		}

		// Check genre changes in recent series
		lastGenre = ""
		for i, series := range profile.RecentSeries {
			if i == 0 && len(series.Genres) > 0 {
				lastGenre = series.Genres[0]
				itemCount++
				continue
			}

			if len(series.Genres) > 0 {
				itemCount++
				// Check if primary genre changed
				if series.Genres[0] != lastGenre {
					genreChangeCount++
					lastGenre = series.Genres[0]
				}
			}
		}

		// Calculate ratio of genre changes to items
		if itemCount > 1 {
			changeRatio := float32(genreChangeCount) / float32(itemCount-1)
			profile.ContentRotationFreq = changeRatio
			if profile.ContentRotationFreq > 1.0 {
				profile.ContentRotationFreq = 1.0
			}
		}
	} else {
		profile.ContentRotationFreq = 0.5 // Default to medium
	}

	// Calculate exploration score (willingness to try new things)
	// This is based on genre breadth, content rotation, and new content score
	profile.ExplorationScore = (profile.GenreBreadth + profile.ContentRotationFreq + profile.NewContentScore) / 3.0
}

// generateMovieRecommendations creates movie recommendations for a user
func (j *RecommendationJob) generateMovieRecommendations(ctx context.Context, jobRunID uint64, user models.User, preferenceProfile *UserPreferenceProfile, config *models.UserConfig) error {
	log := utils.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", user.ID).
		Int("watchedMovieCount", len(preferenceProfile.WatchedMovieIDs)).
		Int("favoriteGenres", len(preferenceProfile.FavoriteMovieGenres)).
		Float32("explorationScore", preferenceProfile.ExplorationScore).
		Msg("Generating movie recommendations")

	// Get existing movies in the user's library
	movies, err := j.movieRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user's movies")
		return err
	}

	// Build maps to track what's in the library and what's been watched
	inLibraryMap := make(map[string]bool)
	watchedMap := make(map[string]bool)

	for _, movie := range movies {
		if movie.Data != nil && movie.Data.Details.Title != "" {
			key := fmt.Sprintf("%s-%d", movie.Data.Details.Title, movie.ReleaseYear)
			inLibraryMap[key] = true

			// Store in owned movie IDs for future reference
			if movie.ExternalIDs.GetID("tmdb") != "" {
				preferenceProfile.OwnedMovieIDs[movie.ExternalIDs.GetID("tmdb")] = true
			}

			// If we have watch history for this movie, mark it as watched
			if _, watched := preferenceProfile.WatchedMovieIDs[movie.ID]; watched {
				watchedMap[key] = true
			}
		}
	}

	// Generate recommendation strategies based on user profile
	var recommendations []*models.Recommendation

	// Get AI client for recommendations if needed
	aiClient, aiErr := j.getAIClient(ctx, user.ID)

	// Decide if we should use AI recommendations
	useAI := config.RecommendationSyncEnabled && aiErr == nil && aiClient != nil

	if useAI {
		// Generate AI recommendations if enabled and available
		j.jobRepo.UpdateJobProgress(ctx, jobRunID, 15, "Generating AI-powered recommendations")

		aiRecs, err := j.generateAIMovieRecommendations(ctx, user.ID, preferenceProfile, config, watchedMap)
		if err != nil {
			log.Error().Err(err).Msg("Failed to generate AI recommendations")
		} else {
			recommendations = append(recommendations, aiRecs...)
		}
	}

	// Add similar media recommendations based on top-rated content
	if config.RecommendationIncludeSimilar {
		j.jobRepo.UpdateJobProgress(ctx, jobRunID, 20, "Generating similar content recommendations")

		// TODO: Implement similar movie recommendations
		// similarRecs := j.generateSimilarMovieRecommendations(ctx, preferenceProfile, watchedMap)
		// recommendations = append(recommendations, similarRecs...)
	}

	// Generate genre-based recommendations
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 25, "Generating genre-based recommendations")

	// TODO: Implement genre-based recommendations
	// genreRecs := j.generateGenreBasedRecommendations(ctx, preferenceProfile, watchedMap)
	// recommendations = append(recommendations, genreRecs...)

	// Save all recommendations to the database
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 30, "Saving recommendations")

	// TODO: Store recommendations in the database
	// err = j.recommendationRepo.SaveRecommendations(ctx, user.ID, "movie", recommendations)

	log.Info().
		Int("count", len(recommendations)).
		Msg("Generated movie recommendations")

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

	// Get AI client for the user
	aiClient, err := j.getAIClient(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get AI client: %v", err)
	}

	// Prepare AI recommendation request
	request := &aitypes.RecommendationRequest{
		MediaType:       "movie",
		Count:           10,
		UserPreferences: map[string]interface{}{},
		ExcludeIDs:      []string{},
		AdditionalContext: fmt.Sprintf("User profile confidence: %.2f. Exploration score: %.2f. Content completion tendency: %.2f.",
			profile.ProfileConfidence, profile.ExplorationScore, profile.ContentCompleter),
	}

	// Use userPreferences for our filters
	filters := request.UserPreferences

	// Add favorite genres with weights
	if len(profile.FavoriteMovieGenres) > 0 {
		// Sort genres by weight
		type genreWeight struct {
			genre  string
			weight float32
		}

		var topGenres []genreWeight
		for genre, weight := range profile.FavoriteMovieGenres {
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

	// Add top rated movies
	if len(profile.TopRatedMovies) > 0 {
		topRated := make([]map[string]interface{}, 0, len(profile.TopRatedMovies))

		for _, movie := range profile.TopRatedMovies {
			topRated = append(topRated, map[string]interface{}{
				"title":  movie.Title,
				"year":   movie.Year,
				"genres": movie.Genres,
			})
		}

		filters["topRatedContent"] = topRated
	}

	// Add recently watched movies
	if len(profile.RecentMovies) > 0 {
		recentlyWatched := make([]map[string]interface{}, 0, len(profile.RecentMovies))

		for _, movie := range profile.RecentMovies {
			recentlyWatched = append(recentlyWatched, map[string]interface{}{
				"title":  movie.Title,
				"year":   movie.Year,
				"genres": movie.Genres,
			})
		}

		filters["recentlyWatched"] = recentlyWatched
	}

	// Add excluded genres
	if len(profile.ExcludedMovieGenres) > 0 {
		filters["excludedGenres"] = profile.ExcludedMovieGenres
	}

	// Add preferred genres
	if len(profile.PreferredMovieGenres) > 0 {
		filters["preferredGenres"] = profile.PreferredMovieGenres
	}

	// Add preferred release year range
	if profile.MovieReleaseYearRange[0] > 0 && profile.MovieReleaseYearRange[1] > 0 {
		filters["yearRange"] = profile.MovieReleaseYearRange
	}

	// Add preferred duration range if set
	if profile.MovieDurationPreference[0] > 0 && profile.MovieDurationPreference[1] > 0 {
		filters["durationRange"] = profile.MovieDurationPreference
	}

	// Add tag preferences if available
	if len(profile.MovieTagPreferences) > 0 {
		// Get top tags
		type tagWeight struct {
			tag    string
			weight float32
		}

		var topTags []tagWeight
		for tag, weight := range profile.MovieTagPreferences {
			topTags = append(topTags, tagWeight{tag, weight})
		}

		// Sort by weight
		sort.Slice(topTags, func(i, j int) bool {
			return topTags[i].weight > topTags[j].weight
		})

		// Take top 5 tags
		preferredTags := []string{}
		for i := 0; i < 5 && i < len(topTags); i++ {
			preferredTags = append(preferredTags, topTags[i].tag)
		}

		filters["preferredTags"] = preferredTags
	}

	// Add watch time patterns if strong preferences exist
	if len(profile.MovieWatchTimes) > 0 {
		// Find peak watch times
		var peakHours []string
		maxCount := 0

		for hour, counts := range profile.MovieWatchTimes {
			total := 0
			for _, count := range counts {
				total += int(count)
			}

			if total > maxCount {
				maxCount = total
				peakHours = []string{hour}
			} else if total == maxCount {
				peakHours = append(peakHours, hour)
			}
		}

		if len(peakHours) > 0 {
			filters["preferredWatchTimes"] = peakHours
		}
	}

	// Add whether to exclude already watched content
	excludeWatched := profile.ExcludeWatched
	filters["excludeWatched"] = excludeWatched

	// Add content rating preferences
	if profile.ContentRatingRange[0] != "" && profile.ContentRatingRange[1] != "" {
		filters["contentRatingRange"] = profile.ContentRatingRange
	}

	// Add preferred languages
	if len(profile.PreferredLanguages) > 0 {
		// Get top languages
		type langWeight struct {
			lang   string
			weight float32
		}

		var topLangs []langWeight
		for lang, weight := range profile.PreferredLanguages {
			topLangs = append(topLangs, langWeight{lang, weight})
		}

		// Sort by weight
		sort.Slice(topLangs, func(i, j int) bool {
			return topLangs[i].weight > topLangs[j].weight
		})

		// Take top 3 languages
		preferredLangs := []string{}
		for i := 0; i < 3 && i < len(topLangs); i++ {
			preferredLangs = append(preferredLangs, topLangs[i].lang)
		}

		filters["preferredLanguages"] = preferredLangs
	}

	// Call the AI service to get recommendations
	log.Info().
		Uint64("userID", userID).
		Interface("filters", filters).
		Msg("Requesting AI movie recommendations")

	// Request recommendations
	aiRecommendations, err := aiClient.GetRecommendations(ctx, request)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get AI recommendations")
		return nil, err
	}

	// Process AI recommendations into our recommendation format
	var recommendations []*models.Recommendation

	if aiRecommendations == nil || len(aiRecommendations.Items) == 0 {
		log.Warn().Msg("No AI recommendations returned")
		return recommendations, nil
	}

	for i, rec := range aiRecommendations.Items {
		// Extract details from the recommendation
		title := rec.Title
		if title == "" {
			continue // Skip recommendations without a title
		}

		// Use the year directly
		year := rec.Year

		// Create a unique key for this movie to check against watched/library maps
		key := fmt.Sprintf("%s-%d", title, year)

		// Skip if we've seen this and are excluding watched content
		if excludeWatched && watchedMap[key] {
			continue
		}

		// TODO: Find existing media item or create a placeholder
		// For now, we'll skip the media item creation

		// Set confidence based on profile confidence
		confidence := float32(0.7) + (profile.ProfileConfidence * 0.3) // Scale from 0.7-1.0 based on profile

		// Extract reason (if available)
		reason := "AI-powered personalized recommendation"
		if rec.Reason != "" {
			reason = rec.Reason
		}

		// Additional metadata
		metadata := map[string]interface{}{
			"source": "ai",
			"rank":   i + 1,
		}

		// Add genres if available
		if len(rec.Genres) > 0 {
			metadata["genres"] = rec.Genres
		}

		// Add TMDB ID if available
		if rec.ExternalID != "" {
			metadata["tmdbId"] = rec.ExternalID
		}

		// Determine if it's in the user's library
		isInLibrary := false
		// Check if we have the external ID as a string
		isInLibrary = profile.OwnedMovieIDs[rec.ExternalID] || watchedMap[key]

		// Create the recommendation based on the models.Recommendation struct
		recommendation := &models.Recommendation{
			UserID:          userID,
			MediaType:       "movie",
			Source:          models.RecommendationSourceAI,
			SourceClientType: "ai",
			Reason:          reason,
			Confidence:      confidence,
			InLibrary:       isInLibrary,
		}

		recommendations = append(recommendations, recommendation)
	}

	log.Info().
		Int("count", len(recommendations)).
		Msg("Generated AI movie recommendations")

	return recommendations, nil
}

// generateSeriesRecommendations creates TV series recommendations for a user
func (j *RecommendationJob) generateSeriesRecommendations(ctx context.Context, jobRunID uint64, user models.User, preferenceProfile *UserPreferenceProfile, config *models.UserConfig) error {
	ctx, log := utils.WithJobID(ctx, jobRunID)
	log.Info().
		Uint64("userID", user.ID).
		Msg("Generating TV series recommendations")

	// Get existing series in the user's library
	seriesList, err := j.seriesRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user's TV series")
		return err
	}

	// Build maps to track what's in the library and what's been watched
	inLibraryMap := make(map[string]bool)
	watchedMap := make(map[string]bool)

	for _, series := range seriesList {
		if series.Data != nil && series.Data.Details.Title != "" {
			key := fmt.Sprintf("%s", series.Data.Details.Title)
			inLibraryMap[key] = true

			// If we have watch history for this series, mark it as watched
			if _, watched := preferenceProfile.WatchedSeriesIDs[series.ID]; watched {
				watchedMap[key] = true
			}
		}
	}

	// Generate recommendation strategies based on user profile
	var recommendations []*models.Recommendation

	// Get AI client for recommendations if needed
	aiClient, aiErr := j.getAIClient(ctx, user.ID)

	// Decide if we should use AI recommendations
	useAI := config.RecommendationSyncEnabled && aiErr == nil && aiClient != nil

	if useAI {
		// Generate AI recommendations if enabled and available
		j.jobRepo.UpdateJobProgress(ctx, jobRunID, 30, "Generating AI-powered TV show recommendations")

		aiRecs, err := j.generateAISeriesRecommendations(ctx, user.ID, preferenceProfile, config, watchedMap)
		if err != nil {
			log.Error().Err(err).Msg("Failed to generate AI series recommendations")
		} else {
			recommendations = append(recommendations, aiRecs...)
		}
	}

	// Add similar media recommendations based on top-rated content
	if config.RecommendationIncludeSimilar {
		j.jobRepo.UpdateJobProgress(ctx, jobRunID, 35, "Generating similar TV show recommendations")

		// TODO: Implement similar series recommendations
		// similarRecs := j.generateSimilarSeriesRecommendations(ctx, preferenceProfile, watchedMap)
		// recommendations = append(recommendations, similarRecs...)
	}

	// Save all recommendations to the database
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 40, "Saving TV show recommendations")

	// TODO: Store recommendations in the database
	// err = j.recommendationRepo.SaveRecommendations(ctx, user.ID, "series", recommendations)

	log.Info().
		Int("count", len(recommendations)).
		Msg("Generated TV series recommendations")

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

	// Get AI client for the user
	aiClient, err := j.getAIClient(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get AI client: %v", err)
	}

	// Prepare AI recommendation request
	request := &aitypes.RecommendationRequest{
		MediaType:       "series",
		Count:           8,
		UserPreferences: map[string]interface{}{},
		ExcludeIDs:      []string{},
		AdditionalContext: fmt.Sprintf("User profile confidence: %.2f. Binge-watching score: %.2f. Content rotation frequency: %.2f.",
			profile.ProfileConfidence, profile.BingeWatchingScore, profile.ContentRotationFreq),
	}

	// Use userPreferences for our filters
	filters := request.UserPreferences

	// Add favorite genres with weights
	if len(profile.FavoriteSeriesGenres) > 0 {
		// Sort genres by weight
		type genreWeight struct {
			genre  string
			weight float32
		}

		var topGenres []genreWeight
		for genre, weight := range profile.FavoriteSeriesGenres {
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

	// Add favorite showrunners/creators
	if len(profile.FavoriteShowrunners) > 0 {
		// Sort by weight
		type showrunnerWeight struct {
			showrunner string
			weight     float32
		}

		var topShowrunners []showrunnerWeight
		for showrunner, weight := range profile.FavoriteShowrunners {
			topShowrunners = append(topShowrunners, showrunnerWeight{showrunner, weight})
		}

		// Sort by weight descending
		sort.Slice(topShowrunners, func(i, j int) bool {
			return topShowrunners[i].weight > topShowrunners[j].weight
		})

		// Take top 3 showrunners
		favoriteShowrunners := []string{}
		for i := 0; i < 3 && i < len(topShowrunners); i++ {
			favoriteShowrunners = append(favoriteShowrunners, topShowrunners[i].showrunner)
		}

		filters["favoriteShowrunners"] = favoriteShowrunners
	}

	// Add top rated series
	if len(profile.TopRatedSeries) > 0 {
		topRated := make([]map[string]interface{}, 0, len(profile.TopRatedSeries))

		for _, series := range profile.TopRatedSeries {
			topRated = append(topRated, map[string]interface{}{
				"title":  series.Title,
				"year":   series.Year,
				"genres": series.Genres,
				"status": series.Status,
			})
		}

		filters["topRatedContent"] = topRated
	}

	// Add recently watched series
	if len(profile.RecentSeries) > 0 {
		recentlyWatched := make([]map[string]interface{}, 0, len(profile.RecentSeries))

		for _, series := range profile.RecentSeries {
			recentlyWatched = append(recentlyWatched, map[string]interface{}{
				"title":  series.Title,
				"year":   series.Year,
				"genres": series.Genres,
				"status": series.Status,
			})
		}

		filters["recentlyWatched"] = recentlyWatched
	}

	// Add excluded genres
	if len(profile.ExcludedSeriesGenres) > 0 {
		filters["excludedGenres"] = profile.ExcludedSeriesGenres
	}

	// Add preferred genres
	if len(profile.PreferredSeriesGenres) > 0 {
		filters["preferredGenres"] = profile.PreferredSeriesGenres
	}

	// Add preferred series status if available (e.g., "Ended", "Continuing")
	if len(profile.PreferredSeriesStatus) > 0 {
		filters["preferredStatus"] = profile.PreferredSeriesStatus
	}

	// Add preferred release year range
	if profile.SeriesReleaseYearRange[0] > 0 && profile.SeriesReleaseYearRange[1] > 0 {
		filters["yearRange"] = profile.SeriesReleaseYearRange
	}

	// Add preferred episode length range if set
	if profile.SeriesEpisodeLengthRange[0] > 0 && profile.SeriesEpisodeLengthRange[1] > 0 {
		filters["episodeLengthRange"] = profile.SeriesEpisodeLengthRange
	}

	// Add tag preferences if available
	if len(profile.SeriesTagPreferences) > 0 {
		// Get top tags
		type tagWeight struct {
			tag    string
			weight float32
		}

		var topTags []tagWeight
		for tag, weight := range profile.SeriesTagPreferences {
			topTags = append(topTags, tagWeight{tag, weight})
		}

		// Sort by weight
		sort.Slice(topTags, func(i, j int) bool {
			return topTags[i].weight > topTags[j].weight
		})

		// Take top 5 tags
		preferredTags := []string{}
		for i := 0; i < 5 && i < len(topTags); i++ {
			preferredTags = append(preferredTags, topTags[i].tag)
		}

		filters["preferredTags"] = preferredTags
	}

	// Add watch time patterns if strong preferences exist
	if len(profile.SeriesWatchTimes) > 0 {
		// Find peak watch times
		var peakHours []string
		maxCount := 0

		for hour, counts := range profile.SeriesWatchTimes {
			total := 0
			for _, count := range counts {
				total += int(count)
			}

			if total > maxCount {
				maxCount = total
				peakHours = []string{hour}
			} else if total == maxCount {
				peakHours = append(peakHours, hour)
			}
		}

		if len(peakHours) > 0 {
			filters["preferredWatchTimes"] = peakHours
		}
	}

	// Add whether to exclude already watched content
	excludeWatched := profile.ExcludeWatched
	filters["excludeWatched"] = excludeWatched

	// Add content rating preferences
	if profile.ContentRatingRange[0] != "" && profile.ContentRatingRange[1] != "" {
		filters["contentRatingRange"] = profile.ContentRatingRange
	}

	// Add preferred languages
	if len(profile.PreferredLanguages) > 0 {
		// Get top languages
		type langWeight struct {
			lang   string
			weight float32
		}

		var topLangs []langWeight
		for lang, weight := range profile.PreferredLanguages {
			topLangs = append(topLangs, langWeight{lang, weight})
		}

		// Sort by weight
		sort.Slice(topLangs, func(i, j int) bool {
			return topLangs[i].weight > topLangs[j].weight
		})

		// Take top 3 languages
		preferredLangs := []string{}
		for i := 0; i < 3 && i < len(topLangs); i++ {
			preferredLangs = append(preferredLangs, topLangs[i].lang)
		}

		filters["preferredLanguages"] = preferredLangs
	}

	// Call the AI service to get recommendations
	log.Info().
		Uint64("userID", userID).
		Interface("filters", filters).
		Msg("Requesting AI series recommendations")

	// Request recommendations
	aiRecommendations, err := aiClient.GetRecommendations(ctx, request)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get AI series recommendations")
		return nil, err
	}

	// Process AI recommendations into our recommendation format
	var recommendations []*models.Recommendation

	if aiRecommendations == nil || len(aiRecommendations.Items) == 0 {
		log.Warn().Msg("No AI series recommendations returned")
		return recommendations, nil
	}

	for i, rec := range aiRecommendations.Items {
		// Extract details from the recommendation
		title := rec.Title
		if title == "" {
			continue // Skip recommendations without a title
		}

		// We don't need to store year in a variable since we don't use it later

		// Create a unique key for this series to check against watched/library maps
		key := fmt.Sprintf("%s", title)

		// Skip if we've seen this and are excluding watched content
		if excludeWatched && watchedMap[key] {
			continue
		}

		// Create a new media item for this recommendation
		// TODO: Find existing media item or create a placeholder
		// For now, we'll skip the media item creation

		// Set confidence based on profile confidence
		confidence := float32(0.7) + (profile.ProfileConfidence * 0.3) // Scale from 0.7-1.0 based on profile

		// Extract reason
		reason := "AI-powered personalized series recommendation"
		if rec.Reason != "" {
			reason = rec.Reason
		}

		// Additional metadata
		metadata := map[string]interface{}{
			"source": "ai",
			"rank":   i + 1,
		}

		// Add genres if available
		if len(rec.Genres) > 0 {
			metadata["genres"] = rec.Genres
		}

		// Add TMDB ID if available
		if rec.ExternalID != "" {
			metadata["tmdbId"] = rec.ExternalID
		}

		// Add status if available (like "continuing" or "ended")
		// rec.Status doesn't exist in RecommendationItem, we'll skip this for now

		// Determine if it's in the user's library
		isInLibrary := false
		// Check if we have the external ID as a string
		isInLibrary = profile.OwnedSeriesIDs[rec.ExternalID] || watchedMap[key]

		// Create the recommendation based on the models.Recommendation struct
		recommendation := &models.Recommendation{
			UserID:          userID,
			MediaType:       "series",
			Source:          models.RecommendationSourceAI,
			SourceClientType: "ai",
			Reason:          reason,
			Confidence:      confidence,
			InLibrary:       isInLibrary,
		}

		recommendations = append(recommendations, recommendation)
	}

	log.Info().
		Int("count", len(recommendations)).
		Msg("Generated AI series recommendations")

	return recommendations, nil
}

// generateMusicRecommendations creates music recommendations for a user
func (j *RecommendationJob) generateMusicRecommendations(ctx context.Context, jobRunID uint64, user models.User, preferenceProfile *UserPreferenceProfile, config *models.UserConfig) error {
	ctx, log := utils.WithJobID(ctx, jobRunID)
	log.Info().
		Uint64("userID", user.ID).
		Msg("Generating music recommendations")

	// Get existing music in the user's library
	tracks, err := j.musicRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user's music tracks")
		return err
	}

	// Build maps to track what's in the library and what's been played
	inLibraryMap := make(map[string]bool)
	playedMap := make(map[string]bool)

	for _, track := range tracks {
		if track.Data != nil && track.Data.Details.Title != "" && track.Data.ArtistName != "" {
			key := fmt.Sprintf("%s-%s", track.Data.ArtistName, track.Data.Details.Title)
			inLibraryMap[key] = true

			// If we have play history for this track, mark it as played
			if _, played := preferenceProfile.PlayedMusicIDs[track.ID]; played {
				playedMap[key] = true
			}
		}
	}

	// Generate recommendation strategies based on user profile
	var recommendations []*models.Recommendation

	// Get AI client for recommendations if needed
	aiClient, aiErr := j.getAIClient(ctx, user.ID)

	// Decide if we should use AI recommendations
	useAI := config.RecommendationSyncEnabled && aiErr == nil && aiClient != nil

	if useAI {
		// Generate AI recommendations if enabled and available
		j.jobRepo.UpdateJobProgress(ctx, jobRunID, 30, "Generating AI-powered music recommendations")

		aiRecs, err := j.generateAIMusicRecommendations(ctx, user.ID, preferenceProfile, config, playedMap)
		if err != nil {
			log.Error().Err(err).Msg("Failed to generate AI music recommendations")
		} else {
			recommendations = append(recommendations, aiRecs...)
		}
	}

	// Add similar media recommendations based on top-played content
	if config.RecommendationIncludeSimilar {
		j.jobRepo.UpdateJobProgress(ctx, jobRunID, 35, "Generating similar music recommendations")

		// TODO: Implement similar music recommendations
		// similarRecs := j.generateSimilarMusicRecommendations(ctx, preferenceProfile, playedMap)
		// recommendations = append(recommendations, similarRecs...)
	}

	// Save all recommendations to the database
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 40, "Saving music recommendations")

	// TODO: Store recommendations in the database
	// err = j.recommendationRepo.SaveRecommendations(ctx, user.ID, "music", recommendations)

	log.Info().
		Int("count", len(recommendations)).
		Msg("Generated music recommendations")

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

	// Get AI client for the user
	aiClient, err := j.getAIClient(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get AI client: %v", err)
	}

	// Prepare AI recommendation request
	request := &aitypes.RecommendationRequest{
		MediaType:       "music",
		Count:           8,
		UserPreferences: map[string]interface{}{},
		ExcludeIDs:      []string{},
		AdditionalContext: fmt.Sprintf("User profile confidence: %.2f. Activity level: %.2f. Music mood preferences variety: %d.",
			profile.ProfileConfidence, profile.OverallActivityLevel["music"], len(profile.MusicMoodPreferences)),
	}

	// Use userPreferences for our filters
	filters := request.UserPreferences

	// Add favorite genres with weights
	if len(profile.FavoriteMusicGenres) > 0 {
		// Sort genres by weight
		type genreWeight struct {
			genre  string
			weight float32
		}

		var topGenres []genreWeight
		for genre, weight := range profile.FavoriteMusicGenres {
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

	// Add favorite artists
	if len(profile.FavoriteArtists) > 0 {
		// Sort artists by weight
		type artistWeight struct {
			artist string
			weight float32
		}

		var topArtists []artistWeight
		for artist, weight := range profile.FavoriteArtists {
			topArtists = append(topArtists, artistWeight{artist, weight})
		}

		// Sort by weight descending
		sort.Slice(topArtists, func(i, j int) bool {
			return topArtists[i].weight > topArtists[j].weight
		})

		// Take top 5 artists
		favoriteArtists := []string{}
		for i := 0; i < 5 && i < len(topArtists); i++ {
			favoriteArtists = append(favoriteArtists, topArtists[i].artist)
		}

		filters["favoriteArtists"] = favoriteArtists
	}

	// Add top rated music
	if len(profile.TopRatedMusic) > 0 {
		topRated := make([]map[string]interface{}, 0, len(profile.TopRatedMusic))

		for _, music := range profile.TopRatedMusic {
			topRated = append(topRated, map[string]interface{}{
				"title":  music.Title,
				"artist": music.Artist,
				"album":  music.Album,
				"genres": music.Genres,
			})
		}

		filters["topRatedContent"] = topRated
	}

	// Add recently played music
	if len(profile.RecentMusic) > 0 {
		recentlyPlayed := make([]map[string]interface{}, 0, len(profile.RecentMusic))

		for _, music := range profile.RecentMusic {
			recentlyPlayed = append(recentlyPlayed, map[string]interface{}{
				"title":  music.Title,
				"artist": music.Artist,
				"album":  music.Album,
				"genres": music.Genres,
			})
		}

		filters["recentlyPlayed"] = recentlyPlayed
	}

	// Add mood preferences if available
	if len(profile.MusicMoodPreferences) > 0 {
		// Sort moods by weight
		type moodWeight struct {
			mood   string
			weight float32
		}

		var topMoods []moodWeight
		for mood, weight := range profile.MusicMoodPreferences {
			topMoods = append(topMoods, moodWeight{mood, weight})
		}

		// Sort by weight descending
		sort.Slice(topMoods, func(i, j int) bool {
			return topMoods[i].weight > topMoods[j].weight
		})

		// Take top 3 moods
		preferredMoods := []string{}
		for i := 0; i < 3 && i < len(topMoods); i++ {
			preferredMoods = append(preferredMoods, topMoods[i].mood)
		}

		filters["preferredMoods"] = preferredMoods
	}

	// Add excluded genres
	if len(profile.ExcludedMusicGenres) > 0 {
		filters["excludedGenres"] = profile.ExcludedMusicGenres
	}

	// Add preferred genres
	if len(profile.PreferredMusicGenres) > 0 {
		filters["preferredGenres"] = profile.PreferredMusicGenres
	}

	// Add duration preferences if set
	if profile.MusicDurationRange[0] > 0 && profile.MusicDurationRange[1] > 0 {
		filters["durationRangeSec"] = profile.MusicDurationRange
	}

	// Add tag preferences if available
	if len(profile.MusicTagPreferences) > 0 {
		// Get top tags
		type tagWeight struct {
			tag    string
			weight float32
		}

		var topTags []tagWeight
		for tag, weight := range profile.MusicTagPreferences {
			topTags = append(topTags, tagWeight{tag, weight})
		}

		// Sort by weight
		sort.Slice(topTags, func(i, j int) bool {
			return topTags[i].weight > topTags[j].weight
		})

		// Take top 5 tags
		preferredTags := []string{}
		for i := 0; i < 5 && i < len(topTags); i++ {
			preferredTags = append(preferredTags, topTags[i].tag)
		}

		filters["preferredTags"] = preferredTags
	}

	// Add play time patterns if strong preferences exist
	if len(profile.MusicPlayTimes) > 0 {
		// Find peak play times
		var peakHours []string
		maxCount := 0

		for hour, counts := range profile.MusicPlayTimes {
			total := 0
			for _, count := range counts {
				total += int(count)
			}

			if total > maxCount {
				maxCount = total
				peakHours = []string{hour}
			} else if total == maxCount {
				peakHours = append(peakHours, hour)
			}
		}

		if len(peakHours) > 0 {
			filters["preferredPlayTimes"] = peakHours
		}
	}

	// Add whether to exclude already played content
	excludePlayed := profile.ExcludeWatched
	filters["excludePlayed"] = excludePlayed

	// Add preferred languages
	if len(profile.PreferredLanguages) > 0 {
		// Get top languages
		type langWeight struct {
			lang   string
			weight float32
		}

		var topLangs []langWeight
		for lang, weight := range profile.PreferredLanguages {
			topLangs = append(topLangs, langWeight{lang, weight})
		}

		// Sort by weight
		sort.Slice(topLangs, func(i, j int) bool {
			return topLangs[i].weight > topLangs[j].weight
		})

		// Take top 3 languages
		preferredLangs := []string{}
		for i := 0; i < 3 && i < len(topLangs); i++ {
			preferredLangs = append(preferredLangs, topLangs[i].lang)
		}

		filters["preferredLanguages"] = preferredLangs
	}

	// Call the AI service to get recommendations
	log.Info().
		Uint64("userID", userID).
		Interface("filters", filters).
		Msg("Requesting AI music recommendations")

	// Request recommendations
	aiRecommendations, err := aiClient.GetRecommendations(ctx, request)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get AI music recommendations")
		return nil, err
	}

	// Process AI recommendations into our recommendation format
	var recommendations []*models.Recommendation

	if aiRecommendations == nil || len(aiRecommendations.Items) == 0 {
		log.Warn().Msg("No AI music recommendations returned")
		return recommendations, nil
	}

	for i, rec := range aiRecommendations.Items {
		// Extract details from the recommendation
		title := rec.Title
		if title == "" {
			continue // Skip recommendations without a title
		}

		// Create a unique key for this track
		key := title
		artist := ""
		if rec.Description != "" {
			// If description has artist info, use it
			artist = rec.Description
			key = fmt.Sprintf("%s-%s", artist, title)
		}

		// Skip if we've played this and are excluding played content
		if excludePlayed && playedMap[key] {
			continue
		}

		// Create a new media item for this recommendation
		// TODO: Find existing media item or create a placeholder
		// For now, we'll skip the media item creation

		// Set confidence based on profile confidence
		confidence := float32(0.7) + (profile.ProfileConfidence * 0.3) // Scale from 0.7-1.0 based on profile

		// Extract reason
		reason := "AI-powered personalized music recommendation"
		if rec.Reason != "" {
			reason = rec.Reason
		}

		// Additional metadata
		metadata := map[string]interface{}{
			"source": "ai",
			"rank":   i + 1,
		}

		// Add genres if available
		if len(rec.Genres) > 0 {
			metadata["genres"] = rec.Genres
		}

		// Add artist if available
		if artist != "" {
			metadata["artist"] = artist
		}

		// Add mood if available
		// rec.Moods doesn't exist in RecommendationItem, we'll skip this for now

		// Add external ID if available
		if rec.ExternalID != "" {
			metadata["externalId"] = rec.ExternalID
		}

		// Determine if it's in the user's library
		isInLibrary := false
		// Check if we have the external ID as a string
		isInLibrary = profile.OwnedMusicIDs[rec.ExternalID] || playedMap[key]

		// Create the recommendation based on the models.Recommendation struct
		recommendation := &models.Recommendation{
			UserID:          userID,
			MediaType:       "music",
			Source:          models.RecommendationSourceAI,
			SourceClientType: "ai",
			Reason:          reason,
			Confidence:      confidence,
			InLibrary:       isInLibrary,
		}

		recommendations = append(recommendations, recommendation)
	}

	log.Info().
		Int("count", len(recommendations)).
		Msg("Generated AI music recommendations")

	return recommendations, nil
}

// SetupMediaSyncJob creates or updates a media sync job for a user
func (j *RecommendationJob) SetupMediaSyncJob(ctx context.Context, userID, clientID uint64, clientType, mediaType, frequency string) error {
	// Implementation would set up a media sync job
	// This is just a stub to satisfy the interface
	log.Printf("Setting up media sync job for user %d, client %d, type %s", userID, clientID, mediaType)
	return nil
}

// UpdateUserRecommendationSchedule updates the recommendation schedule for a user
func (j *RecommendationJob) UpdateUserRecommendationSchedule(ctx context.Context, userID uint64) error {
	// Implementation would update the recommendation schedule for a user
	// This is just a stub to satisfy the interface
	log.Printf("Updating recommendation schedule for user %d", userID)
	return nil
}
