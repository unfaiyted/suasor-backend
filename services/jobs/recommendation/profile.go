package recommendation

import (
	"context"
	"fmt"
	"suasor/utils"
	"time"
)

// buildUserPreferenceProfile analyzes a user's media consumption to build a preference profile
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
		WatchTimeOfDay:        make(map[int]int),
		TypicalSessionLength:  make(map[string]float32),
		OverallActivityLevel:  make(map[string]float32),

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
		MusicMoodPreferences: make(map[string]float32),
		ExcludedMusicGenres:  []string{},
		PreferredMusicGenres: []string{},
	}

	// Get user's movie watch history
	movieHistory, err := j.userMovieDataRepo.GetUserHistory(ctx, userID, 100, 0)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get movie history")
		return nil, fmt.Errorf("failed to get movie history: %w", err)
	}

	// Process movie history
	j.processMovieHistory(ctx, profile, movieHistory)

	// Get user's series watch history
	seriesHistory, err := j.userSeriesDataRepo.GetUserHistory(ctx, userID, 100, 0)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get series history")
		// Don't return - continue with what we have
	} else {
		// Process series history
		j.processSeriesHistory(ctx, profile, seriesHistory)
	}

	// Get user's music play history
	musicHistory, err := j.userMusicDataRepo.GetUserHistory(ctx, userID, 100, 0)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get music history")
		// Don't return - continue with what we have
	} else {
		// Process music history using implementation from music.go
		j.processMusicHistory(ctx, profile, musicHistory)
	}

	// Calculate advanced metrics across all media types
	j.calculateAdvancedMetrics(profile)

	log.Info().
		Int("movies", len(profile.RecentMovies)).
		Int("series", len(profile.RecentSeries)).
		Int("music", len(profile.RecentMusic)).
		Dur("duration", time.Since(startTime)).
		Msg("Built user preference profile")

	return profile, nil
}

// calculateAdvancedMetrics computes cross-media metrics based on history
func (j *RecommendationJob) calculateAdvancedMetrics(profile *UserPreferenceProfile) {
	// Skip if we don't have enough data
	if len(profile.RecentMovies) == 0 && len(profile.RecentSeries) == 0 && len(profile.RecentMusic) == 0 {
		return
	}

	// Calculate preferred genres across all media types
	// This helps identify if a user has consistent preferences across different media

	// Calculate overall activity patterns
	// This can show when the user is most active, useful for scheduling recommendations

	// Add movie genres to preferred genres if they're highly weighted
	for genre, weight := range profile.FavoriteMovieGenres {
		if weight > 2.0 && !contains(profile.PreferredMovieGenres, genre) {
			profile.PreferredMovieGenres = append(profile.PreferredMovieGenres, genre)
		}
	}

	// Similarly for excluded genres (those with negative weights)
	for genre, weight := range profile.FavoriteMovieGenres {
		if weight < 0.2 && len(profile.FavoriteMovieGenres) > 5 && !contains(profile.ExcludedMovieGenres, genre) {
			profile.ExcludedMovieGenres = append(profile.ExcludedMovieGenres, genre)
		}
	}

	// Add series genres to preferred genres if they're highly weighted
	for genre, weight := range profile.FavoriteSeriesGenres {
		if weight > 2.0 && !contains(profile.PreferredSeriesGenres, genre) {
			profile.PreferredSeriesGenres = append(profile.PreferredSeriesGenres, genre)
		}
	}

	// Similarly for excluded series genres
	for genre, weight := range profile.FavoriteSeriesGenres {
		if weight < 0.2 && len(profile.FavoriteSeriesGenres) > 5 && !contains(profile.ExcludedSeriesGenres, genre) {
			profile.ExcludedSeriesGenres = append(profile.ExcludedSeriesGenres, genre)
		}
	}

	// Add music genres to preferred genres if they're highly weighted
	for genre, weight := range profile.FavoriteMusicGenres {
		if weight > 2.0 && !contains(profile.PreferredMusicGenres, genre) {
			profile.PreferredMusicGenres = append(profile.PreferredMusicGenres, genre)
		}
	}

	// Similarly for excluded music genres
	for genre, weight := range profile.FavoriteMusicGenres {
		if weight < 0.2 && len(profile.FavoriteMusicGenres) > 5 && !contains(profile.ExcludedMusicGenres, genre) {
			profile.ExcludedMusicGenres = append(profile.ExcludedMusicGenres, genre)
		}
	}
}

// Helper function to check if a string is in a slice
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// UpdateUserRecommendationSchedule updates the recommendation schedule for a user
func (j *RecommendationJob) UpdateUserRecommendationSchedule(ctx context.Context, userID uint64) error {
	log := utils.LoggerFromContext(ctx)
	// Implementation would update the recommendation schedule for a user
	// This is just a stub to satisfy the interface
	log.Info().Msg("Updating recommendation schedule for user")
	return nil
}
