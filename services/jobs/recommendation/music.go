package recommendation

import (
	"context"
	"fmt"
	"sort"
	aitypes "suasor/client/ai/types"
	mediatypes "suasor/client/media/types"
	"suasor/types/models"
	"suasor/utils"
	"time"
)

// processMusicHistory analyzes music play history to build preferences
// This is a full implementation of music history processing
func (j *RecommendationJob) processMusicHistoryImpl(ctx context.Context, profile *UserPreferenceProfile, histories []models.UserMediaItemData[*mediatypes.Track]) {
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
			track, err := j.itemRepos.TrackRepo().GetByID(ctx, history.MediaItemID)
			if err != nil || track == nil || track.Data == nil {
				continue
			}

			// Process track details
			trackData := track.Data

			// Build music summary
			summary := MusicSummary{
				ID:           track.ID,
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
				summary.DurationSec = int(trackData.Duration)
			}

			// Add detailed rating if available
			if history.UserRating > 0 {
				summary.DetailedRating = &RatingDetails{
					Overall:   history.UserRating,
					Timestamp: history.LastPlayedAt.Unix(),
				}
			}

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

		// Scale: 0 = no activity, 1 = very active (30+ tracks per month)
		activityScore := float32(recentCount) / 30.0
		if activityScore > 1.0 {
			activityScore = 1.0
		}

		profile.OverallActivityLevel["music"] = activityScore
	}

	// Process preferred genres
	if len(musicGenreWeights) > 0 {
		// Sort genres by weight
		type genreWeight struct {
			genre  string
			weight float32
		}

		var genres []genreWeight
		for genre, weight := range musicGenreWeights {
			genres = append(genres, genreWeight{genre, weight})
		}

		sort.Slice(genres, func(i, j int) bool {
			return genres[i].weight > genres[j].weight
		})

		// Add top genres to preferred genres
		count := len(genres)
		if count > 5 {
			count = 5
		}

		for i := 0; i < count; i++ {
			if genres[i].weight > 1.5 {
				profile.PreferredMusicGenres = append(profile.PreferredMusicGenres, genres[i].genre)
			}
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

// generateMusicRecommendations generates music recommendations for a user
// This is a skeleton for future implementation
func (j *RecommendationJob) generateMusicRecommendations(ctx context.Context, jobRunID uint64, user models.User, preferenceProfile *UserPreferenceProfile, config *models.UserConfig) error {
	log := utils.LoggerFromContext(ctx)
	log.Info().Msg("Music recommendations not yet implemented")

	// Placeholder for future implementation
	// Similar to the movie recommendations implementation

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
	excludePlayed := profile.ExcludePlayed
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
		// TODO: Implement moods logic
		// TODO: Need to look at how we handle music from an artist/album/track

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
			UserID:           userID,
			MediaType:        "music",
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
		Msg("Generated AI music recommendations")

	return recommendations, nil
}

// processMusicHistory analyzes music play history to build preferences
func (j *RecommendationJob) processMusicHistory(ctx context.Context, profile *UserPreferenceProfile, histories []*models.UserMediaItemData[*mediatypes.Track]) {
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
			track, err := j.itemRepos.TrackRepo().GetByID(ctx, history.MediaItemID)
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
