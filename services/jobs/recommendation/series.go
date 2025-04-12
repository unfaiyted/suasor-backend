package recommendation

import (
	"context"
	"fmt"
	"sort"
	"time"

	aitypes "suasor/client/ai/types"
	mediatypes "suasor/client/media/types"
	"suasor/types/models"
	"suasor/utils"
)

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
	excludeWatched := profile.ExcludePlayed
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
		// TODO: Need to look at the series IDs we are using and how we want to
		// handle this better
		isInLibrary = profile.OwnedSeriesIDs[rec.ExternalID] || watchedMap[key]

		// Create the recommendation based on the models.Recommendation struct
		recommendation := &models.Recommendation{
			UserID:           userID,
			MediaType:        "series",
			Source:           models.RecommendationSourceAI,
			SourceClientType: "ai",
			Reasoning:        reason,
			Confidence:       confidence,
			InLibrary:        isInLibrary,
		}

		recommendations = append(recommendations, recommendation)
	}

	log.Info().
		Int("count", len(recommendations)).
		Msg("Generated AI series recommendations")

	return recommendations, nil
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
	// Other data structures for series processing

	// Default days of week
	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	for _, day := range days {
		watchDays[day] = 0
		profile.SeriesWatchDays[day] = 0
	}

	// Process each history item (implementation would be similar to processMovieHistory)
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
				ID:            series.ID,
				Title:         seriesData.Details.Title,
				Year:          seriesData.ReleaseYear,
				Genres:        seriesData.Genres,
				Rating:        history.UserRating,
				LastWatchDate: history.LastPlayedAt.Unix(),
			}

			// Add status if available
			if seriesData.Status != "" {
				summary.Status = seriesData.Status
			}

			// Add detailed rating if available
			if history.UserRating > 0 {
				summary.DetailedRating = &RatingDetails{
					Overall:   history.UserRating,
					Timestamp: history.LastPlayedAt.Unix(),
				}
			}

			// Process genres with weighted scoring
			if seriesData.Genres != nil && len(seriesData.Genres) > 0 {
				for _, genre := range seriesData.Genres {
					seriesGenreWeights[genre] += weight
				}
			}

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
