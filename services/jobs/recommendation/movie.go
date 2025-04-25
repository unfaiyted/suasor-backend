package recommendation

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"strings"
	aitypes "suasor/clients/ai/types"
	mediatypes "suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils"
	"suasor/utils/logger"
	"time"
)

// generateMovieRecommendations generates movie recommendations for a user
func (j *RecommendationJob) generateMovieRecommendations(ctx context.Context, jobRunID uint64, user models.User, preferenceProfile *UserPreferenceProfile, config *models.UserConfig) error {
	log := logger.LoggerFromContext(ctx)
	log.Info().Msg("Generating movie recommendations")

	// Get an AI client for the user
	aiClient, err := j.getAIClient(ctx, user.ID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get AI client, skipping movie recommendations")
		return err
	}

	// Create prompt with user preferences
	prompt := j.buildMovieRecommendationPrompt(user, preferenceProfile)
	log.Debug().Str("prompt", prompt).Msg("Generated AI prompt for movie recommendations")

	// Call AI model for recommendations
	model := "claude-3-opus-20240229" // Default model
	if config.AIModelPreferences != nil && config.AIModelPreferences.DefaultModelForRecommendations != "" {
		model = config.AIModelPreferences.DefaultModelForRecommendations
	}

	// Prepare the system message
	systemPrompt := "You are a movie recommendation expert. Your goal is to provide personalized movie recommendations based on the user's preferences, watch history, and specified criteria. Provide detailed, thoughtful explanations for why each movie would appeal to this specific user."

	// Set response format as JSON
	responseFormat := map[string]any{
		"type": "json_object",
	}

	// Prepare request options
	options := map[string]interface{}{
		"temperature":     0.7,
		"max_tokens":      4000,
		"top_p":           0.9,
		"response_format": responseFormat,
	}

	// Call the AI model
	resp, err := aiClient.GenerateContent(ctx, systemPrompt, prompt, model, options)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate recommendations from AI")
		return err
	}

	// Process AI response (parse JSON)
	var recommendations []MovieRecommendation
	var aiResponse map[string]interface{}

	// Extract the content from the response
	if resp.Text != "" {
		err = json.Unmarshal([]byte(resp.Text), &aiResponse)
		if err != nil {
			log.Error().Err(err).Str("responseText", resp.Text).Msg("Failed to parse AI response JSON")
			return err
		}

		// Process recommendations from the parsed JSON
		if recommendationsData, ok := aiResponse["recommendations"].([]interface{}); ok {
			for _, recData := range recommendationsData {
				if recMap, ok := recData.(map[string]interface{}); ok {
					recommendation := MovieRecommendation{
						Title:         utils.GetStringFromMap(recMap, "title", ""),
						Year:          utils.GetIntFromMap(recMap, "year", 0),
						Score:         utils.GetFloatFromMap(recMap, "score", 0),
						Reasoning:     utils.GetStringFromMap(recMap, "reasoning", ""),
						RecommendedBy: "AI",
						Timestamp:     time.Now(),
					}

					// Extract genres
					if genresData, ok := recMap["genres"].([]interface{}); ok {
						for _, g := range genresData {
							if genre, ok := g.(string); ok {
								recommendation.Genres = append(recommendation.Genres, genre)
							}
						}
					}

					// Extract similar movies
					if similarData, ok := recMap["similarToMovies"].([]interface{}); ok {
						for _, s := range similarData {
							if similar, ok := s.(string); ok {
								recommendation.SimilarToMovies = append(recommendation.SimilarToMovies, similar)
							}
						}
					}

					// Extract matching actors
					if actorsData, ok := recMap["matchesActors"].([]interface{}); ok {
						for _, a := range actorsData {
							if actor, ok := a.(string); ok {
								recommendation.MatchesActors = append(recommendation.MatchesActors, actor)
							}
						}
					}

					// Extract matching directors
					if directorsData, ok := recMap["matchesDirectors"].([]interface{}); ok {
						for _, d := range directorsData {
							if director, ok := d.(string); ok {
								recommendation.MatchesDirectors = append(recommendation.MatchesDirectors, director)
							}
						}
					}

					// Extract matching genres
					if matchGenresData, ok := recMap["matchesGenres"].([]interface{}); ok {
						for _, g := range matchGenresData {
							if genre, ok := g.(string); ok {
								recommendation.MatchesGenres = append(recommendation.MatchesGenres, genre)
							}
						}
					}

					recommendations = append(recommendations, recommendation)
				}
			}
		}
	}

	// Store recommendations in the database
	// This would typically involve setting up a table for recommendations
	// and storing the data there. For now, we'll just log them.
	log.Info().Int("count", len(recommendations)).Msg("Generated movie recommendations")

	// Convert recommendations to database model and store them
	var modelRecommendations []*models.Recommendation
	for _, rec := range recommendations {
		modelRec := &models.Recommendation{
			UserID:           user.ID,
			MediaType:        "movie",
			Title:            rec.Title,
			Year:             rec.Year,
			Genres:           rec.Genres,
			Confidence:       rec.Score,
			Reasoning:        rec.Reasoning,
			SimilarItems:     rec.SimilarToMovies,
			MatchesActors:    rec.MatchesActors,
			MatchesDirectors: rec.MatchesDirectors,
			MatchesGenres:    rec.MatchesGenres,
			RecommendedBy:    rec.RecommendedBy,
			JobRunID:         jobRunID,
			CreatedAt:        rec.Timestamp,
		}
		modelRecommendations = append(modelRecommendations, modelRec)
	}

	// Store recommendations in the database
	if len(modelRecommendations) > 0 {
		err = j.recommendationRepo.CreateMany(ctx, modelRecommendations)
		if err != nil {
			log.Error().Err(err).Msg("Failed to store movie recommendations")
			return err
		}
		log.Info().
			Int("count", len(modelRecommendations)).
			Msg("Stored movie recommendations in database")
	}

	return nil
}

// buildMovieRecommendationPrompt creates a prompt for the AI model based on user preferences
func (j *RecommendationJob) buildMovieRecommendationPrompt(user models.User, profile *UserPreferenceProfile) string {
	var prompt strings.Builder

	// Start with instruction
	prompt.WriteString("Generate personalized movie recommendations for a user based on their preferences and watch history.\n\n")

	// Add user context
	prompt.WriteString("# User Information\n")
	prompt.WriteString(fmt.Sprintf("User ID: %d\n", user.ID))
	prompt.WriteString(fmt.Sprintf("Username: %s\n\n", user.Username))

	// Add movie preferences
	prompt.WriteString("# Movie Preferences\n")

	// Provide favorite genres with weights
	prompt.WriteString("## Favorite Genres\n")
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

		// Take top 10 genres at most
		count := len(topGenres)
		if count > 10 {
			count = 10
		}

		for i := 0; i < count; i++ {
			prompt.WriteString(fmt.Sprintf("- %s (weight: %.2f)\n", topGenres[i].genre, topGenres[i].weight))
		}
	} else {
		prompt.WriteString("No favorite genres identified yet\n")
	}
	prompt.WriteString("\n")

	// Provide favorite actors with weights
	prompt.WriteString("## Favorite Actors\n")
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

		// Take top 10 actors at most
		count := len(topActors)
		if count > 10 {
			count = 10
		}

		for i := 0; i < count; i++ {
			prompt.WriteString(fmt.Sprintf("- %s (weight: %.2f)\n", topActors[i].actor, topActors[i].weight))
		}
	} else {
		prompt.WriteString("No favorite actors identified yet\n")
	}
	prompt.WriteString("\n")

	// Provide favorite directors with weights
	prompt.WriteString("## Favorite Directors\n")
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

		// Take top 5 directors at most
		count := len(topDirectors)
		if count > 5 {
			count = 5
		}

		for i := 0; i < count; i++ {
			prompt.WriteString(fmt.Sprintf("- %s (weight: %.2f)\n", topDirectors[i].director, topDirectors[i].weight))
		}
	} else {
		prompt.WriteString("No favorite directors identified yet\n")
	}
	prompt.WriteString("\n")

	// Add recently watched movies
	prompt.WriteString("## Recently Watched Movies\n")
	if len(profile.RecentMovies) > 0 {
		// Take up to 10 most recent movies
		count := len(profile.RecentMovies)
		if count > 10 {
			count = 10
		}

		for i := 0; i < count; i++ {
			movie := profile.RecentMovies[i]
			rating := ""
			if movie.Rating > 0 {
				rating = fmt.Sprintf(" - Rated: %.1f/5", movie.Rating)
			}

			genres := ""
			if len(movie.Genres) > 0 {
				genres = fmt.Sprintf(" - Genres: %s", strings.Join(movie.Genres, ", "))
			}

			prompt.WriteString(fmt.Sprintf("- %s (%d)%s%s\n",
				movie.Title,
				movie.Year,
				rating,
				genres))
		}
	} else {
		prompt.WriteString("No recently watched movies\n")
	}
	prompt.WriteString("\n")

	// Add highly rated movies
	prompt.WriteString("## Highly Rated Movies\n")
	if len(profile.TopRatedMovies) > 0 {
		// Take up to 10 highest rated movies
		count := len(profile.TopRatedMovies)
		if count > 10 {
			count = 10
		}

		for i := 0; i < count; i++ {
			movie := profile.TopRatedMovies[i]
			rating := "N/A"
			if movie.Rating > 0 {
				rating = fmt.Sprintf("%.1f/5", movie.Rating)
			}

			genres := ""
			if len(movie.Genres) > 0 {
				genres = fmt.Sprintf(" - Genres: %s", strings.Join(movie.Genres, ", "))
			}

			prompt.WriteString(fmt.Sprintf("- %s (%d) - Rated: %s%s\n",
				movie.Title,
				movie.Year,
				rating,
				genres))
		}
	} else {
		prompt.WriteString("No highly rated movies\n")
	}
	prompt.WriteString("\n")

	// Add watch time patterns
	prompt.WriteString("## Watch Time Patterns\n")
	if len(profile.MovieWatchTimes) > 0 {
		// Convert hour to counts for display
		type hourCount struct {
			hour  string
			count int64
		}

		var patterns []hourCount
		for hour, counts := range profile.MovieWatchTimes {
			var totalCount int64
			for _, count := range counts {
				totalCount += count
			}
			patterns = append(patterns, hourCount{hour, totalCount})
		}

		// Sort by count descending
		sort.Slice(patterns, func(i, j int) bool {
			return patterns[i].count > patterns[j].count
		})

		// Take top 5 most active hours
		count := len(patterns)
		if count > 5 {
			count = 5
		}

		for i := 0; i < count; i++ {
			prompt.WriteString(fmt.Sprintf("- %s: %d views\n", patterns[i].hour, patterns[i].count))
		}
	} else {
		prompt.WriteString("No watch time patterns identified yet\n")
	}
	prompt.WriteString("\n")

	// Add activity level
	if movieActivity, ok := profile.OverallActivityLevel["movie"]; ok {
		prompt.WriteString(fmt.Sprintf("## Activity Level: %.2f (0-1 scale, where 1 is very active)\n\n", movieActivity))
	}

	// Add exclusion criteria
	prompt.WriteString("# Exclusion Criteria\n")

	// Add already watched movies
	prompt.WriteString("## Already Watched Movies\n")
	prompt.WriteString(fmt.Sprintf("User has watched %d unique movies (IDs not shown for brevity)\n\n", len(profile.WatchedMovieIDs)))

	// Add excluded genres if any
	if len(profile.ExcludedMovieGenres) > 0 {
		prompt.WriteString("## Excluded Genres\n")
		for _, genre := range profile.ExcludedMovieGenres {
			prompt.WriteString(fmt.Sprintf("- %s\n", genre))
		}
		prompt.WriteString("\n")
	}

	// Add request for recommendations in specific format
	prompt.WriteString(`# Recommendation Request
Please provide 5-10 movie recommendations based on the user's preferences and watch history. For each recommendation, include:

1. Title
2. Year of release
3. Genres (as an array)
4. A score between 0-1 indicating how well it matches the user's preferences
5. A detailed explanation of why this movie would appeal to the user based on their specific preferences
6. Similar movies they've already watched that influenced this recommendation (as an array)
7. Actors in the movie that match their preferences (as an array)
8. Directors in the movie that match their preferences (as an array)
9. Genres in the movie that match their preferences (as an array)

Format your response as a valid JSON object with a "recommendations" array containing these details for each movie.

Example structure:
{
  "recommendations": [
    {
      "title": "Movie Title",
      "year": 2023,
      "genres": ["Action", "Sci-Fi"],
      "score": 0.92,
      "reasoning": "This movie would appeal to the user because...",
      "similarToMovies": ["Similar Movie 1", "Similar Movie 2"],
      "matchesActors": ["Actor 1", "Actor 2"],
      "matchesDirectors": ["Director 1"],
      "matchesGenres": ["Action", "Sci-Fi"]
    }
  ]
}
`)

	return prompt.String()
}

// generateAIMovieRecommendations uses AI to generate personalized movie recommendations
func (j *RecommendationJob) generateAIMovieRecommendations(
	ctx context.Context,
	userID uint64,
	profile *UserPreferenceProfile,
	config *models.UserConfig,
	watchedMap map[string]bool) ([]*models.Recommendation, error) {

	log := logger.LoggerFromContext(ctx)

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

	// TODO: Implement UserConfig to this to ensure we are picking up the users specified
	// Genres to include/exclude and give them a stronger weight.

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
			UserID:           userID,
			MediaType:        "movie",
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
		Msg("Generated AI movie recommendations")

	return recommendations, nil
}

// processMovieHistory analyzes movie watch history to build preferences
func (j *RecommendationJob) processMovieHistory(ctx context.Context, profile *UserPreferenceProfile, histories []*models.UserMediaItemData[*mediatypes.Movie]) {
	log := logger.LoggerFromContext(ctx)

	// Maps for processing
	recentMovies := []MovieSummary{}
	highRatedMovies := []MovieSummary{}
	watchTimes := make(map[int]int)
	watchDays := make(map[string]int)
	movieGenreWeights := make(map[string]float32)
	actorWeights := make(map[string]float32)
	directorWeights := make(map[string]float32)

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
			movie, err := j.itemRepos.MovieRepo().GetByID(ctx, history.MediaItemID)
			if err != nil || movie == nil || movie.Data == nil {
				continue
			}

			// Process movie details
			movieData := movie.Data

			// Get credits if available
			var credits []*models.Credit
			if j.creditRepo != nil {
				credits, _ = j.getCreditsForMediaItem(ctx, movie.ID)
			}

			// Build movie summary
			summary := MovieSummary{
				ID:            movie.ID,
				Title:         movie.Title,
				Year:          movie.ReleaseYear,
				Genres:        movieData.Details.Genres,
				Rating:        history.UserRating,
				WatchCount:    int(history.PlayCount),
				LastWatchDate: history.LastPlayedAt.Unix(),
			}

			// Add duration if available
			if movieData.Details.Duration > 0 {
				summary.DurationMin = int(movieData.Details.Duration)
			}

			// Add director if available from credits
			directors := GetCrewByRole(credits, "Director")
			if len(directors) > 0 {
				summary.Director = directors[0].Name
				// Add to director weights
				directorWeights[directors[0].Name] += weight
			}

			// Add cast if available from credits
			cast := GetCastFromCredits(credits, 5)
			castNames := ExtractNamesFromCredits(cast)
			summary.Cast = castNames

			// Process cast weights
			for _, actor := range castNames {
				actorWeights[actor] += weight
			}

			// Add detailed rating if available
			if history.UserRating > 0 {
				summary.DetailedRating = &RatingDetails{
					Overall:   history.UserRating,
					Timestamp: history.LastPlayedAt.Unix(),
				}
			}

			// Process genres with weighted scoring
			if len(movieData.Details.Genres) > 0 {
				for _, genre := range movieData.Details.Genres {
					movieGenreWeights[genre] += weight
				}
			}

			// Add to recent movies list
			recentMovies = append(recentMovies, summary)

			// Add to high-rated movies if applicable
			if history.UserRating > 3.5 {
				highRatedMovies = append(highRatedMovies, summary)
				summary.IsFavorite = history.UserRating > 4.0
			}
		}
	}

	// Sort recent movies by watch date (newest first)
	sort.Slice(recentMovies, func(i, j int) bool {
		return recentMovies[i].LastWatchDate > recentMovies[j].LastWatchDate
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
	profile.FavoriteActors = actorWeights
	profile.FavoriteDirectors = directorWeights

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

		// Scale: 0 = no activity, 1 = very active (5+ movies per month)
		activityScore := float32(recentCount) / 5.0
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
