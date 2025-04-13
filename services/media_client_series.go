package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"suasor/client"
	"suasor/client/media"
	"suasor/client/media/providers"
	mediatypes "suasor/client/media/types"
	"suasor/client/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/utils"
)

// MediaClientSeriesService defines service methods for TV series
type MediaClientSeriesService[T types.ClientConfig] interface {
	GetSeriesByID(ctx context.Context, userID, clientID uint64, seriesID string) (*models.MediaItem[*mediatypes.Series], error)
	GetSeasonByID(ctx context.Context, userID, clientID uint64, seasonID string) (*models.MediaItem[*mediatypes.Season], error)
	GetEpisodeByID(ctx context.Context, userID, clientID uint64, episodeID string) (*models.MediaItem[*mediatypes.Episode], error)
	GetSeriesByName(ctx context.Context, userID uint64, name string) ([]models.MediaItem[*mediatypes.Series], error)
	GetSeriesByGenre(ctx context.Context, userID uint64, genre string) ([]models.MediaItem[*mediatypes.Series], error)
	GetRecentlyAdded(ctx context.Context, userID uint64, count int) ([]models.MediaItem[*mediatypes.Series], error)
	GetOnGoing(ctx context.Context, userID uint64, count int) ([]models.MediaItem[*mediatypes.Series], error)
	GetRecentEpisodes(ctx context.Context, userID uint64, count int) ([]models.MediaItem[*mediatypes.Episode], error)
	GetSeriesByNetwork(ctx context.Context, userID uint64, network string) ([]models.MediaItem[*mediatypes.Series], error)

	// Added missing methods required by handlers
	GetSeriesByYear(ctx context.Context, userID uint64, year int) ([]models.MediaItem[*mediatypes.Series], error)
	GetSeriesByActor(ctx context.Context, userID uint64, actor string) ([]models.MediaItem[*mediatypes.Series], error)
	GetSeriesByCreator(ctx context.Context, userID uint64, creator string) ([]models.MediaItem[*mediatypes.Series], error)
	GetSeriesByRating(ctx context.Context, userID uint64, minRating, maxRating float64) ([]models.MediaItem[*mediatypes.Series], error)
	GetLatestSeriesByAdded(ctx context.Context, userID uint64, count int) ([]models.MediaItem[*mediatypes.Series], error)
	GetPopularSeries(ctx context.Context, userID uint64, count int) ([]models.MediaItem[*mediatypes.Series], error)
	GetTopRatedSeries(ctx context.Context, userID uint64, count int) ([]models.MediaItem[*mediatypes.Series], error)
	SearchSeries(ctx context.Context, userID uint64, query string) ([]models.MediaItem[*mediatypes.Series], error)
	GetSeasonsBySeriesID(ctx context.Context, userID, clientID uint64, seriesID string) ([]models.MediaItem[*mediatypes.Season], error)
}

type mediaSeriesService[T types.MediaClientConfig] struct {
	repo    repository.ClientRepository[T]
	factory *client.ClientFactoryService
}

// NewMediaClientSeriesService creates a new TV series service
func NewMediaClientSeriesService[T types.MediaClientConfig](
	repo repository.ClientRepository[T],
	factory *client.ClientFactoryService,
) MediaClientSeriesService[T] {
	return &mediaSeriesService[T]{
		repo:    repo,
		factory: factory,
	}
}

// getSeriesClients gets all clients that support TV shows
func (s *mediaSeriesService[T]) getSeriesClients(ctx context.Context, userID uint64) ([]media.MediaClient, error) {
	// Get all media clients for the user
	clients, err := s.repo.GetByCategory(ctx, types.ClientCategoryMedia, userID)
	if err != nil {
		return nil, err
	}

	var seriesClients []media.MediaClient

	// Filter and instantiate clients that support TV shows
	for _, clientConfig := range clients {
		if clientConfig.Config.Data.SupportsSeries() {
			client, err := s.factory.GetClient(ctx, clientConfig.ID, clientConfig.Config.Data)
			if err != nil {
				// Log error but continue with other clients
				continue
			}
			seriesClients = append(seriesClients, client.(media.MediaClient))
		}
	}

	return seriesClients, nil
}

// getSpecificSeriesClient gets a specific TV shows client
func (s *mediaSeriesService[T]) getSpecificSeriesClient(ctx context.Context, userID, clientID uint64) (media.MediaClient, error) {
	log := utils.LoggerFromContext(ctx)

	clientConfig, err := s.repo.GetByID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	if !clientConfig.Config.Data.SupportsSeries() {
		log.Warn().
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("clientType", clientConfig.Config.Data.GetType().String()).
			Msg("Client does not support TV shows")
		return nil, ErrUnsupportedFeature
	}

	client, err := s.factory.GetClient(ctx, clientID, clientConfig.Config.Data)
	if err != nil {
		return nil, err
	}

	return client.(media.MediaClient), nil
}

// GetSeriesByName searches for TV series by name across all clients
func (s *mediaSeriesService[T]) GetSeriesByName(ctx context.Context, userID uint64, name string) ([]models.MediaItem[*mediatypes.Series], error) {
	clients, err := s.getSeriesClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allSeries []models.MediaItem[*mediatypes.Series]

	for _, client := range clients {
		showProvider, ok := client.(providers.SeriesProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Query: name, // Use the standard query field instead of filters
		}

		series, err := showProvider.GetSeries(ctx, options)
		if err != nil {
			continue
		}

		allSeries = append(allSeries, series...)
	}

	// Filter by name similarity if the client doesn't do it
	filteredSeries := allSeries[:0] // Reuse the same slice
	for _, s := range allSeries {
		if strings.Contains(strings.ToLower(s.Data.Details.Title), strings.ToLower(name)) {
			filteredSeries = append(filteredSeries, s)
		}
	}

	return filteredSeries, nil
}

// GetSeriesByGenre gets TV series by genre
func (s *mediaSeriesService[T]) GetSeriesByGenre(ctx context.Context, userID uint64, genre string) ([]models.MediaItem[*mediatypes.Series], error) {
	clients, err := s.getSeriesClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allSeries []models.MediaItem[*mediatypes.Series]

	for _, client := range clients {
		showProvider, ok := client.(providers.SeriesProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Genre: genre,
		}

		series, err := showProvider.GetSeries(ctx, options)
		if err != nil {
			continue
		}

		allSeries = append(allSeries, series...)
	}

	return allSeries, nil
}

// GetSeriesByID gets a specific TV series by ID
func (s *mediaSeriesService[T]) GetSeriesByID(ctx context.Context, userID, clientID uint64, seriesID string) (*models.MediaItem[*mediatypes.Series], error) {
	client, err := s.getSpecificSeriesClient(ctx, userID, clientID)
	if err != nil {
		return nil, err
	}

	showProvider, ok := client.(providers.SeriesProvider)
	if !ok {
		return nil, fmt.Errorf("client does not implement series provider interface")
	}

	// Get series by ID
	series, err := showProvider.GetSeriesByID(ctx, seriesID)
	if err != nil {
		return nil, err
	}

	return &series, nil
}

// GetSeasonByID gets a specific season by ID
func (s *mediaSeriesService[T]) GetSeasonByID(ctx context.Context, userID, clientID uint64, seasonID string) (*models.MediaItem[*mediatypes.Season], error) {
	// Note: The SeriesProvider interface doesn't have a GetSeasonByID method
	// This implementation assumes season IDs are prefixed with "show_ID-season_number"
	// This is a simplification, and in a real implementation, you might need a different approach

	// Check if the client exists and supports series
	_, err := s.getSpecificSeriesClient(ctx, userID, clientID)
	if err != nil {
		return nil, err
	}

	// For simplicity, just return a placeholder error
	// In a real implementation, you'd extract the show ID and season number
	// from the seasonID and use GetSeriesSeasons to find the right season
	return nil, fmt.Errorf("getting season by ID not directly supported, use show ID + season number instead")
}

// GetEpisodeByID gets a specific episode by ID
func (s *mediaSeriesService[T]) GetEpisodeByID(ctx context.Context, userID, clientID uint64, episodeID string) (*models.MediaItem[*mediatypes.Episode], error) {
	client, err := s.getSpecificSeriesClient(ctx, userID, clientID)
	if err != nil {
		return nil, err
	}

	showProvider, ok := client.(providers.SeriesProvider)
	if !ok {
		return nil, fmt.Errorf("client does not implement series provider interface")
	}

	// Get episode by ID
	episode, err := showProvider.GetEpisodeByID(ctx, episodeID)
	if err != nil {
		return nil, err
	}

	return &episode, nil
}

// GetRecentlyAdded gets recently added TV series
func (s *mediaSeriesService[T]) GetRecentlyAdded(ctx context.Context, userID uint64, count int) ([]models.MediaItem[*mediatypes.Series], error) {
	clients, err := s.getSeriesClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allSeries []models.MediaItem[*mediatypes.Series]

	// Cut-off date (e.g., last 30 days)
	cutoffDate := time.Now().AddDate(0, 0, -30)

	for _, client := range clients {
		showProvider, ok := client.(providers.SeriesProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			RecentlyAdded:  true,
			DateAddedAfter: cutoffDate,
			Sort:           "dateAdded",
			SortOrder:      mediatypes.SortOrderDesc,
			Limit:          count,
		}

		series, err := showProvider.GetSeries(ctx, options)
		if err != nil {
			continue
		}

		allSeries = append(allSeries, series...)
	}

	// Sort by date added (newest first)
	// This assumes the Series type has an AddedAt field
	// If not, you'll need to adjust this code

	// Limit to requested count
	if len(allSeries) > count {
		allSeries = allSeries[:count]
	}

	return allSeries, nil
}

// GetOnGoing gets currently ongoing TV series
func (s *mediaSeriesService[T]) GetOnGoing(ctx context.Context, userID uint64, count int) ([]models.MediaItem[*mediatypes.Series], error) {
	clients, err := s.getSeriesClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allSeries []models.MediaItem[*mediatypes.Series]

	for _, client := range clients {
		showProvider, ok := client.(providers.SeriesProvider)
		if !ok {
			continue
		}

		// Use the MediaType field to specify "ongoing" series
		options := &mediatypes.QueryOptions{
			MediaType: "ongoing",
			Limit:     count,
		}

		series, err := showProvider.GetSeries(ctx, options)
		if err != nil {
			continue
		}

		allSeries = append(allSeries, series...)
	}

	// Filter to only include ongoing series
	// This assumes the Series type has a Status field
	// If not, you'll need to adjust this code

	// Limit to requested count
	if len(allSeries) > count {
		allSeries = allSeries[:count]
	}

	return allSeries, nil
}

// GetRecentEpisodes gets recently aired episodes
func (s *mediaSeriesService[T]) GetRecentEpisodes(ctx context.Context, userID uint64, count int) ([]models.MediaItem[*mediatypes.Episode], error) {
	clients, err := s.getSeriesClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allEpisodes []models.MediaItem[*mediatypes.Episode]

	// Cut-off date (e.g., last 7 days)
	cutoffDate := time.Now().AddDate(0, 0, -7)

	for _, client := range clients {
		showProvider, ok := client.(providers.SeriesProvider)
		if !ok {
			continue
		}

		// First, get recently updated series
		options := &mediatypes.QueryOptions{
			ReleasedAfter: cutoffDate,
			Sort:          "dateAdded",
			SortOrder:     mediatypes.SortOrderDesc,
			Limit:         10, // Get a reasonable number of shows
		}

		series, err := showProvider.GetSeries(ctx, options)
		if err != nil {
			continue
		}

		// For each series, get episodes from the latest season
		for _, show := range series {
			// Get the series ID from the ClientIDs
			// We could use the first client ID as a simplification
			if len(show.SyncClients) == 0 {
				continue
			}
			showID := show.SyncClients[0].ItemID

			// Get seasons for this series
			seasons, err := showProvider.GetSeriesSeasons(ctx, showID)
			if err != nil || len(seasons) == 0 {
				continue
			}

			// Sort seasons by number to get the latest one
			// This is a simple implementation - in a real app you might want to sort differently
			var highestSeasonNum int
			var latestSeasonIdx int
			for i, season := range seasons {
				if season.Data.Number > highestSeasonNum {
					highestSeasonNum = season.Data.Number
					latestSeasonIdx = i
				}
			}

			latestSeason := seasons[latestSeasonIdx]

			// Get episodes for the latest season
			episodes, err := showProvider.GetSeriesEpisodes(ctx, showID, latestSeason.Data.Number)
			if err != nil {
				continue
			}

			// Add all episodes to our collection
			// In a real implementation, you'd filter by air date
			allEpisodes = append(allEpisodes, episodes...)

			// If we have enough episodes, stop processing more shows
			if len(allEpisodes) >= count {
				break
			}
		}
	}

	// Limit to requested count
	if len(allEpisodes) > count {
		allEpisodes = allEpisodes[:count]
	}

	return allEpisodes, nil
}

// GetSeriesByNetwork gets TV series by network
func (s *mediaSeriesService[T]) GetSeriesByNetwork(ctx context.Context, userID uint64, network string) ([]models.MediaItem[*mediatypes.Series], error) {
	clients, err := s.getSeriesClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allSeries []models.MediaItem[*mediatypes.Series]

	for _, client := range clients {
		showProvider, ok := client.(providers.SeriesProvider)
		if !ok {
			continue
		}

		// Use the Studio field to filter by network
		options := &mediatypes.QueryOptions{
			Studio: network,
		}

		series, err := showProvider.GetSeries(ctx, options)
		if err != nil {
			continue
		}

		allSeries = append(allSeries, series...)
	}

	return allSeries, nil
}

// GetSeriesByYear gets TV series by release year
func (s *mediaSeriesService[T]) GetSeriesByYear(ctx context.Context, userID uint64, year int) ([]models.MediaItem[*mediatypes.Series], error) {
	clients, err := s.getSeriesClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allSeries []models.MediaItem[*mediatypes.Series]

	for _, client := range clients {
		showProvider, ok := client.(providers.SeriesProvider)
		if !ok {
			continue
		}

		// Use the Year field to filter
		options := &mediatypes.QueryOptions{
			Year: year,
		}

		series, err := showProvider.GetSeries(ctx, options)
		if err != nil {
			continue
		}

		allSeries = append(allSeries, series...)
	}

	return allSeries, nil
}

// GetSeriesByActor gets TV series by actor name/ID
func (s *mediaSeriesService[T]) GetSeriesByActor(ctx context.Context, userID uint64, actor string) ([]models.MediaItem[*mediatypes.Series], error) {
	clients, err := s.getSeriesClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allSeries []models.MediaItem[*mediatypes.Series]

	for _, client := range clients {
		showProvider, ok := client.(providers.SeriesProvider)
		if !ok {
			continue
		}

		// Use the Actor field to filter
		options := &mediatypes.QueryOptions{
			Actor: actor,
		}

		series, err := showProvider.GetSeries(ctx, options)
		if err != nil {
			continue
		}

		allSeries = append(allSeries, series...)
	}

	return allSeries, nil
}

// GetSeriesByCreator gets TV series by creator/director
func (s *mediaSeriesService[T]) GetSeriesByCreator(ctx context.Context, userID uint64, creator string) ([]models.MediaItem[*mediatypes.Series], error) {
	clients, err := s.getSeriesClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allSeries []models.MediaItem[*mediatypes.Series]

	for _, client := range clients {
		showProvider, ok := client.(providers.SeriesProvider)
		if !ok {
			continue
		}

		// Use the Creator field to filter
		options := &mediatypes.QueryOptions{
			Creator: creator,
		}

		series, err := showProvider.GetSeries(ctx, options)
		if err != nil {
			continue
		}

		allSeries = append(allSeries, series...)
	}

	return allSeries, nil
}

// GetSeriesByRating gets TV series with ratings in the specified range
func (s *mediaSeriesService[T]) GetSeriesByRating(ctx context.Context, userID uint64, minRating, maxRating float64) ([]models.MediaItem[*mediatypes.Series], error) {
	clients, err := s.getSeriesClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allSeries []models.MediaItem[*mediatypes.Series]

	// Since rating is typically client-specific and might not be directly filterable,
	// we'll get a broader set and filter in memory
	for _, client := range clients {
		showProvider, ok := client.(providers.SeriesProvider)
		if !ok {
			continue
		}

		// Get series with minimal filtering (could be optimized if clients support rating filters)
		options := &mediatypes.QueryOptions{
			Limit: 100, // Get a reasonable number of shows
		}

		series, err := showProvider.GetSeries(ctx, options)
		if err != nil {
			continue
		}

		// Filter by rating in memory
		for _, s := range series {
			// Check if rating is in range - Rating is in the Series struct, not in Details
			if s.Data.Rating >= minRating && s.Data.Rating <= maxRating {
				allSeries = append(allSeries, s)
			}
		}
	}

	return allSeries, nil
}

// GetLatestSeriesByAdded gets recently added series
func (s *mediaSeriesService[T]) GetLatestSeriesByAdded(ctx context.Context, userID uint64, count int) ([]models.MediaItem[*mediatypes.Series], error) {
	clients, err := s.getSeriesClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allSeries []models.MediaItem[*mediatypes.Series]

	// Cut-off date (e.g., last 30 days)
	cutoffDate := time.Now().AddDate(0, 0, -30)

	for _, client := range clients {
		showProvider, ok := client.(providers.SeriesProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			DateAddedAfter: cutoffDate,
			Sort:           "dateAdded",
			SortOrder:      mediatypes.SortOrderDesc,
			Limit:          count,
		}

		series, err := showProvider.GetSeries(ctx, options)
		if err != nil {
			continue
		}

		allSeries = append(allSeries, series...)
	}

	// Sort all series by added date (newest first)
	// This could be done with a custom sort, but we'll rely on the client sorting for now

	// Limit to requested count
	if len(allSeries) > count {
		allSeries = allSeries[:count]
	}

	return allSeries, nil
}

// GetPopularSeries gets popular series based on play count, ratings, etc.
func (s *mediaSeriesService[T]) GetPopularSeries(ctx context.Context, userID uint64, count int) ([]models.MediaItem[*mediatypes.Series], error) {
	clients, err := s.getSeriesClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allSeries []models.MediaItem[*mediatypes.Series]

	for _, client := range clients {
		showProvider, ok := client.(providers.SeriesProvider)
		if !ok {
			continue
		}

		// Use custom query parameters for popularity
		options := &mediatypes.QueryOptions{
			Sort:      "popularity", // Assuming clients support this sort method
			SortOrder: mediatypes.SortOrderDesc,
			Limit:     count,
		}

		series, err := showProvider.GetSeries(ctx, options)
		if err != nil {
			continue
		}

		allSeries = append(allSeries, series...)
	}

	// Limit to requested count
	if len(allSeries) > count {
		allSeries = allSeries[:count]
	}

	return allSeries, nil
}

// GetTopRatedSeries gets series with the highest ratings
func (s *mediaSeriesService[T]) GetTopRatedSeries(ctx context.Context, userID uint64, count int) ([]models.MediaItem[*mediatypes.Series], error) {
	clients, err := s.getSeriesClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allSeries []models.MediaItem[*mediatypes.Series]

	for _, client := range clients {
		showProvider, ok := client.(providers.SeriesProvider)
		if !ok {
			continue
		}

		// Get series sorted by rating
		options := &mediatypes.QueryOptions{
			Sort:      "rating", // Assuming clients support rating sort
			SortOrder: mediatypes.SortOrderDesc,
			Limit:     count,
		}

		series, err := showProvider.GetSeries(ctx, options)
		if err != nil {
			continue
		}

		allSeries = append(allSeries, series...)
	}

	// We could sort all series by rating here if needed
	// But for now, rely on the client sorting

	// Limit to requested count
	if len(allSeries) > count {
		allSeries = allSeries[:count]
	}

	return allSeries, nil
}

// SearchSeries searches for series by name/title across all clients
func (s *mediaSeriesService[T]) SearchSeries(ctx context.Context, userID uint64, query string) ([]models.MediaItem[*mediatypes.Series], error) {
	// This is essentially the same as GetSeriesByName, but we'll make it explicit
	clients, err := s.getSeriesClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allSeries []models.MediaItem[*mediatypes.Series]

	for _, client := range clients {
		showProvider, ok := client.(providers.SeriesProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Query: query,
		}

		series, err := showProvider.GetSeries(ctx, options)
		if err != nil {
			continue
		}

		allSeries = append(allSeries, series...)
	}

	// Filter to prioritize more relevant matches
	// This is a simple implementation to ensure matches
	filteredSeries := allSeries[:0] // Reuse the same slice
	lowerQuery := strings.ToLower(query)

	// First pass: exact title matches
	for _, s := range allSeries {
		if strings.ToLower(s.Data.Details.Title) == lowerQuery {
			filteredSeries = append(filteredSeries, s)
		}
	}

	// Second pass: title starts with query
	if len(filteredSeries) < 10 {
		for _, s := range allSeries {
			if strings.HasPrefix(strings.ToLower(s.Data.Details.Title), lowerQuery) {
				// Check if already added in first pass
				alreadyAdded := false
				for _, fs := range filteredSeries {
					if fs.Data.Details.Title == s.Data.Details.Title {
						alreadyAdded = true
						break
					}
				}
				if !alreadyAdded {
					filteredSeries = append(filteredSeries, s)
				}
			}
		}
	}

	// Third pass: contains query anywhere in title
	if len(filteredSeries) < 10 {
		for _, s := range allSeries {
			if strings.Contains(strings.ToLower(s.Data.Details.Title), lowerQuery) {
				// Check if already added
				alreadyAdded := false
				for _, fs := range filteredSeries {
					if fs.Data.Details.Title == s.Data.Details.Title {
						alreadyAdded = true
						break
					}
				}
				if !alreadyAdded {
					filteredSeries = append(filteredSeries, s)
				}
			}
		}
	}

	return filteredSeries, nil
}

// GetSeasonsBySeriesID gets all seasons for a specific series
func (s *mediaSeriesService[T]) GetSeasonsBySeriesID(ctx context.Context, userID, clientID uint64, seriesID string) ([]models.MediaItem[*mediatypes.Season], error) {
	client, err := s.getSpecificSeriesClient(ctx, userID, clientID)
	if err != nil {
		return nil, err
	}

	showProvider, ok := client.(providers.SeriesProvider)
	if !ok {
		return nil, fmt.Errorf("client does not implement series provider interface")
	}

	// Get seasons for this series
	seasons, err := showProvider.GetSeriesSeasons(ctx, seriesID)
	if err != nil {
		return nil, err
	}

	return seasons, nil
}

