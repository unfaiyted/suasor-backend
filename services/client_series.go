package services

import (
	"context"
	"fmt"
	"time"

	"suasor/clients"
	"suasor/clients/media/providers"
	mediatypes "suasor/clients/media/types"
	"suasor/clients/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/utils/logger"
)

// ClientSeriesService defines service methods for TV series
type ClientSeriesService[T types.ClientConfig] interface {
	GetSeriesByID(ctx context.Context, clientID uint64, seriesID string) (*models.MediaItem[*mediatypes.Series], error)
	GetSeasonByID(ctx context.Context, clientID uint64, seasonID string) (*models.MediaItem[*mediatypes.Season], error)
	GetEpisodeByID(ctx context.Context, clientID uint64, episodeID string) (*models.MediaItem[*mediatypes.Episode], error)

	GetSeriesByName(ctx context.Context, clientID uint64, name string) ([]*models.MediaItem[*mediatypes.Series], error)
	GetSeriesByGenre(ctx context.Context, clientID uint64, genre string) ([]*models.MediaItem[*mediatypes.Series], error)
	GetRecentlyAdded(ctx context.Context, clientID uint64, count int) ([]*models.MediaItem[*mediatypes.Series], error)
	GetOnGoing(ctx context.Context, clientID uint64, count int) ([]*models.MediaItem[*mediatypes.Series], error)
	GetRecentEpisodes(ctx context.Context, clientID uint64, count int) ([]*models.MediaItem[*mediatypes.Episode], error)
	GetSeriesByNetwork(ctx context.Context, clientID uint64, network string) ([]*models.MediaItem[*mediatypes.Series], error)

	GetSeriesByYear(ctx context.Context, clientID uint64, year int) ([]*models.MediaItem[*mediatypes.Series], error)
	GetSeriesByActor(ctx context.Context, clientID uint64, actor string) ([]*models.MediaItem[*mediatypes.Series], error)
	GetSeriesByCreator(ctx context.Context, clientID uint64, creator string) ([]*models.MediaItem[*mediatypes.Series], error)
	GetSeriesByRating(ctx context.Context, clientID uint64, minRating, maxRating float64) ([]*models.MediaItem[*mediatypes.Series], error)
	GetLatestSeriesByAdded(ctx context.Context, clientID uint64, count int) ([]*models.MediaItem[*mediatypes.Series], error)
	GetPopularSeries(ctx context.Context, clientID uint64, count int) ([]*models.MediaItem[*mediatypes.Series], error)
	GetTopRatedSeries(ctx context.Context, clientID uint64, count int) ([]*models.MediaItem[*mediatypes.Series], error)
	SearchSeries(ctx context.Context, clientID uint64, query *mediatypes.QueryOptions) ([]*models.MediaItem[*mediatypes.Series], error)
	GetSeasonsBySeriesID(ctx context.Context, clientID uint64, seriesID string) ([]*models.MediaItem[*mediatypes.Season], error)
	GetEpisodesBySeriesID(ctx context.Context, clientID uint64, seriesID string) ([]*models.MediaItem[*mediatypes.Episode], error)
	GetEpisodesBySeasonNbr(ctx context.Context, clientID uint64, seriesID string, seasonNumber int) ([]*models.MediaItem[*mediatypes.Episode], error)
}

type clientSeriesService[T types.ClientMediaConfig] struct {
	clientRepo    repository.ClientRepository[T]
	clientFactory *clients.ClientProviderFactoryService
}

// NewClientSeriesService creates a new TV series service
func NewClientSeriesService[T types.ClientMediaConfig](
	clientRepo repository.ClientRepository[T],
	factory *clients.ClientProviderFactoryService,
) ClientSeriesService[T] {
	return &clientSeriesService[T]{
		clientRepo:    clientRepo,
		clientFactory: factory,
	}
}

// getSeriesProviders gets all providers that support TV shows
func (s *clientSeriesService[T]) getSeriesProviders(ctx context.Context, userID uint64) ([]providers.SeriesProvider, error) {
	// Get all media providers for the user
	clients, err := s.clientRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var providers []providers.SeriesProvider

	// Filter and instantiate providers that support TV shows
	for _, clientConfig := range clients {
		if clientConfig.Config.SupportsSeries() {
			provider, err := s.clientFactory.GetSeriesProvider(ctx, clientConfig.ID, clientConfig.Config)
			if err != nil {
				// Log error but continue with other providers
				continue
			}
			providers = append(providers, provider)
		}
	}

	return providers, nil
}

// getSeriesProvider gets a specific TV shows client
func (s *clientSeriesService[T]) getSeriesProvider(ctx context.Context, clientID uint64) (providers.SeriesProvider, error) {
	log := logger.LoggerFromContext(ctx)

	clientConfig, err := s.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	if !clientConfig.Config.SupportsSeries() {
		log.Warn().
			Uint64("clientID", clientID).
			Str("clientType", clientConfig.Config.GetType().String()).
			Msg("Client does not support TV shows")
		return nil, ErrUnsupportedFeature
	}

	provider, err := s.clientFactory.GetSeriesProvider(ctx, clientID, clientConfig.Config)
	if err != nil {
		return nil, err
	}

	return provider, nil
}

// GetSeriesByName searches for TV series by name across all providers
func (s *clientSeriesService[T]) GetSeriesByName(ctx context.Context, clientID uint64, name string) ([]*models.MediaItem[*mediatypes.Series], error) {
	log := logger.LoggerFromContext(ctx)
	provider, err := s.getSeriesProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	options := &mediatypes.QueryOptions{
		Query: name, // Use the standard query field instead of filters
	}

	series, err := provider.GetSeries(ctx, options)
	if err != nil {
		// Log error but continue with other providers
		log.Warn().
			Err(err).
			Uint64("clientID", clientID).
			Str("name", name).
			Msg("Error getting series by name from client")
	}

	return series, nil
}

// GetSeriesByGenre gets TV series by genre
func (s *clientSeriesService[T]) GetSeriesByGenre(ctx context.Context, clientID uint64, genre string) ([]*models.MediaItem[*mediatypes.Series], error) {
	log := logger.LoggerFromContext(ctx)
	provider, err := s.getSeriesProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}
	options := &mediatypes.QueryOptions{
		Genre: genre,
	}

	series, err := provider.GetSeries(ctx, options)
	if err != nil || len(series) == 0 {
		// Log error but continue with other providers
		log.Warn().
			Err(err).
			Uint64("clientID", clientID).
			Str("genre", genre).
			Msg("Error getting series by genre from client")

	}

	return series, nil
}

// GetSeriesByID gets a specific TV series by ID
func (s *clientSeriesService[T]) GetSeriesByID(ctx context.Context, clientID uint64, seriesID string) (*models.MediaItem[*mediatypes.Series], error) {
	client, err := s.getSeriesProvider(ctx, clientID)
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

	return series, nil
}

// GetSeasonByID gets a specific season by ID
func (s *clientSeriesService[T]) GetSeasonByID(ctx context.Context, clientID uint64, seasonID string) (*models.MediaItem[*mediatypes.Season], error) {
	// Note: The SeriesProvider interface doesn't have a GetSeasonByID method
	// This implementation assumes season IDs are prefixed with "show_ID-season_number"
	// This is a simplification, and in a real implementation, you might need a different approach

	// Check if the client exists and supports series
	_, err := s.getSeriesProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	// For simplicity, just return a placeholder error
	// In a real implementation, you'd extract the show ID and season number
	// from the seasonID and use GetSeriesSeasons to find the right season
	return nil, fmt.Errorf("getting season by ID not directly supported, use show ID + season number instead")
}

// GetEpisodeByID gets a specific episode by ID
func (s *clientSeriesService[T]) GetEpisodeByID(ctx context.Context, clientID uint64, episodeID string) (*models.MediaItem[*mediatypes.Episode], error) {
	client, err := s.getSeriesProvider(ctx, clientID)
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

	return episode, nil
}

// GetRecentlyAdded gets recently added TV series
func (s *clientSeriesService[T]) GetRecentlyAdded(ctx context.Context, clientID uint64, count int) ([]*models.MediaItem[*mediatypes.Series], error) {
	log := logger.LoggerFromContext(ctx)
	provider, err := s.getSeriesProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	// Cut-off date (e.g., last 30 days)
	cutoffDate := time.Now().AddDate(0, 0, -30)

	options := &mediatypes.QueryOptions{
		RecentlyAdded:  true,
		DateAddedAfter: &cutoffDate,
		Sort:           "dateAdded",
		SortOrder:      mediatypes.SortOrderDesc,
		Limit:          count,
	}

	series, err := provider.GetSeries(ctx, options)
	if err != nil || len(series) == 0 {
		// Log error but continue with other providers
		log.Warn().
			Err(err).
			Uint64("clientID", clientID).
			Msg("Error getting recently added series from client")
	}

	// Sort by date added (newest first)
	// This assumes the Series type has an AddedAt field
	// If not, you'll need to adjust this code

	return series, nil
}

// GetOnGoing gets currently ongoing TV series
func (s *clientSeriesService[T]) GetOnGoing(ctx context.Context, clientID uint64, count int) ([]*models.MediaItem[*mediatypes.Series], error) {
	log := logger.LoggerFromContext(ctx)
	provider, err := s.getSeriesProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	// Use the MediaType field to specify "ongoing" series
	options := &mediatypes.QueryOptions{
		MediaType: "ongoing",
		Limit:     count,
	}

	series, err := provider.GetSeries(ctx, options)
	if err != nil || len(series) == 0 {
		// Log error but continue with other providers
		log.Warn().
			Err(err).
			Uint64("clientID", clientID).
			Msg("Error getting ongoing series from client")
	}

	return series, nil
}

// GetRecentEpisodes gets recently aired episodes
func (s *clientSeriesService[T]) GetRecentEpisodes(ctx context.Context, clientID uint64, count int) ([]*models.MediaItem[*mediatypes.Episode], error) {
	log := logger.LoggerFromContext(ctx)
	provider, err := s.getSeriesProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	// Cut-off date (e.g., last 7 days)
	cutoffDate := time.Now().AddDate(0, 0, -7)

	// First, get recently updated series
	options := &mediatypes.QueryOptions{
		ReleasedAfter: &cutoffDate,
		Sort:          "dateAdded",
		SortOrder:     mediatypes.SortOrderDesc,
		Limit:         10, // Get a reasonable number of shows
	}

	series, err := provider.GetSeries(ctx, options)

	// For each series, get episodes from the latest season
	for _, show := range series {
		// Get the series ID from the ClientIDs
		// We could use the first client ID as a simplification
		if len(show.SyncClients) == 0 {
			continue
		}
		showID := show.SyncClients[0].ItemID

		// Get seasons for this series
		seasons, err := provider.GetSeriesSeasons(ctx, showID)
		if err != nil || len(seasons) == 0 {
			continue
		}

		// Sort seasons by number to get the latest one
		// This is a simple implementation
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
		episodes, err := provider.GetSeriesEpisodesBySeasonNbr(ctx, showID, latestSeason.Data.Number)
		if err != nil || len(episodes) == 0 {
			// Log error but continue with other providers
			log.Warn().
				Err(err).
				Uint64("clientID", clientID).
				Str("showID", showID).
				Msg("Error getting episodes for series from client")
		}

		return episodes, nil

	}
	return nil, nil

}

// GetSeriesByNetwork gets TV series by network
func (s *clientSeriesService[T]) GetSeriesByNetwork(ctx context.Context, clientID uint64, network string) ([]*models.MediaItem[*mediatypes.Series], error) {
	log := logger.LoggerFromContext(ctx)
	provider, err := s.getSeriesProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	// Use the Studio field to filter by network
	options := &mediatypes.QueryOptions{
		Studio: network,
	}

	series, err := provider.GetSeries(ctx, options)
	if err != nil || len(series) == 0 {
		// Log error but continue with other providers
		log.Warn().
			Err(err).
			Uint64("clientID", clientID).
			Str("network", network).
			Msg("Error getting series by network from client")
	}

	return series, nil
}

// GetSeriesByYear gets TV series by release year
func (s *clientSeriesService[T]) GetSeriesByYear(ctx context.Context, clientID uint64, year int) ([]*models.MediaItem[*mediatypes.Series], error) {
	log := logger.LoggerFromContext(ctx)
	provider, err := s.getSeriesProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	// Use the Year field to filter
	options := &mediatypes.QueryOptions{
		Year: year,
	}

	series, err := provider.GetSeries(ctx, options)
	if err != nil {
		// Log error but continue with other providers
		log.Warn().
			Err(err).
			Uint64("clientID", clientID).
			Int("year", year).
			Msg("Error getting series by year from client")
	}

	return series, nil
}

// GetSeriesByActor gets TV series by actor name/ID
func (s *clientSeriesService[T]) GetSeriesByActor(ctx context.Context, clientID uint64, actor string) ([]*models.MediaItem[*mediatypes.Series], error) {
	log := logger.LoggerFromContext(ctx)
	provider, err := s.getSeriesProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	// Use the Actor field to filter
	options := &mediatypes.QueryOptions{
		Actor: actor,
	}

	series, err := provider.GetSeries(ctx, options)
	if err != nil {
		// Log error but continue with other providers
		log.Warn().
			Err(err).
			Uint64("clientID", clientID).
			Str("actor", actor).
			Msg("Error getting series by actor from client")
		return nil, err
	}
	return series, nil

}

// GetSeriesByCreator gets TV series by creator/director
func (s *clientSeriesService[T]) GetSeriesByCreator(ctx context.Context, clientID uint64, creator string) ([]*models.MediaItem[*mediatypes.Series], error) {
	log := logger.LoggerFromContext(ctx)
	provider, err := s.getSeriesProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	// Use the Creator field to filter
	options := &mediatypes.QueryOptions{
		Creator: creator,
	}

	series, err := provider.GetSeries(ctx, options)
	if err != nil {
		// Log error but continue with other providers
		log.Warn().
			Err(err).
			Uint64("clientID", clientID).
			Str("creator", creator).
			Msg("Error getting series by creator from client")
		return nil, err

	}

	return series, nil
}

// GetSeriesByRating gets TV series with ratings in the specified range
func (s *clientSeriesService[T]) GetSeriesByRating(ctx context.Context, clientID uint64, minRating, maxRating float64) ([]*models.MediaItem[*mediatypes.Series], error) {
	log := logger.LoggerFromContext(ctx)
	provider, err := s.getSeriesProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	// Get series with minimal filtering (could be optimized if providers support rating filters)
	options := &mediatypes.QueryOptions{
		Limit: 100, // Get a reasonable number of shows
	}

	series, err := provider.GetSeries(ctx, options)
	if err != nil {
		// Log error but continue with other providers
		log.Warn().
			Err(err).
			Uint64("clientID", clientID).
			Msg("Error getting series with rating from client")
	}

	return series, nil
}

// GetLatestSeriesByAdded gets recently added series
func (s *clientSeriesService[T]) GetLatestSeriesByAdded(ctx context.Context, clientID uint64, count int) ([]*models.MediaItem[*mediatypes.Series], error) {
	log := logger.LoggerFromContext(ctx)
	provider, err := s.getSeriesProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	// Cut-off date (e.g., last 30 days)
	cutoffDate := time.Now().AddDate(0, 0, -30)

	options := &mediatypes.QueryOptions{
		DateAddedAfter: &cutoffDate,
		Sort:           "dateAdded",
		SortOrder:      mediatypes.SortOrderDesc,
		Limit:          count,
	}

	series, err := provider.GetSeries(ctx, options)
	if err != nil {
		// Log error but continue with other providers
		log.Warn().
			Err(err).
			Uint64("clientID", clientID).
			Msg("Error getting recently added series from client")
		return nil, err
	}

	// Sort all series by added date (newest first)
	return series, nil
}

// GetPopularSeries gets popular series based on play count, ratings, etc.
func (s *clientSeriesService[T]) GetPopularSeries(ctx context.Context, clientID uint64, count int) ([]*models.MediaItem[*mediatypes.Series], error) {
	log := logger.LoggerFromContext(ctx)
	provider, err := s.getSeriesProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	// Use custom query parameters for popularity
	options := &mediatypes.QueryOptions{
		Sort:      "popularity", // Assuming providers support this sort method
		SortOrder: mediatypes.SortOrderDesc,
		Limit:     count,
	}

	series, err := provider.GetSeries(ctx, options)
	if err != nil || len(series) == 0 {
		log.Warn().
			Err(err).
			Uint64("clientID", clientID).
			Msg("Error getting popular series from client")
		return nil, err
	}

	return series, nil
}

// GetTopRatedSeries gets series with the highest ratings
func (s *clientSeriesService[T]) GetTopRatedSeries(ctx context.Context, clientID uint64, count int) ([]*models.MediaItem[*mediatypes.Series], error) {
	log := logger.LoggerFromContext(ctx)
	provider, err := s.getSeriesProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	// Get series sorted by rating
	options := &mediatypes.QueryOptions{
		Sort:      "rating", // Assuming providers support rating sort
		SortOrder: mediatypes.SortOrderDesc,
		Limit:     count,
	}

	series, err := provider.GetSeries(ctx, options)
	if err != nil {
		log.Warn().
			Err(err).
			Uint64("clientID", clientID).
			Msg("Error getting top rated series from client")
		return nil, err
	}

	// We could sort all series by rating here if needed
	// But for now, rely on the client sorting

	return series, nil
}

// SearchSeries searches for series by name/title across all providers
func (s *clientSeriesService[T]) SearchSeries(ctx context.Context, clientID uint64, query *mediatypes.QueryOptions) ([]*models.MediaItem[*mediatypes.Series], error) {
	log := logger.LoggerFromContext(ctx)
	// This is essentially the same as GetSeriesByName, but we'll make it explicit
	provider, err := s.getSeriesProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	series, err := provider.GetSeries(ctx, query)
	if err != nil {
		log.Warn().
			Err(err).
			Uint64("clientID", clientID).
			Msg("Error searching series from client")
		return nil, err
	}

	return series, nil
}

// GetSeasonsBySeriesID gets all seasons for a specific series
func (s *clientSeriesService[T]) GetSeasonsBySeriesID(ctx context.Context, clientID uint64, seriesID string) ([]*models.MediaItem[*mediatypes.Season], error) {
	client, err := s.getSeriesProvider(ctx, clientID)
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

// GetEpisodesBySeriesID gets all episodes for a specific series
func (s *clientSeriesService[T]) GetEpisodesBySeriesID(ctx context.Context, clientID uint64, seriesID string) ([]*models.MediaItem[*mediatypes.Episode], error) {
	log := logger.LoggerFromContext(ctx)
	
	client, err := s.getSeriesProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	showProvider, ok := client.(providers.SeriesProvider)
	if !ok {
		return nil, fmt.Errorf("client does not implement series provider interface")
	}

	// First, get all seasons for this series
	seasons, err := showProvider.GetSeriesSeasons(ctx, seriesID)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", clientID).
			Str("seriesID", seriesID).
			Msg("Failed to get seasons for series")
		return nil, fmt.Errorf("failed to get seasons for series: %w", err)
	}

	if len(seasons) == 0 {
		log.Warn().
			Uint64("clientID", clientID).
			Str("seriesID", seriesID).
			Msg("No seasons found for series")
		return []*models.MediaItem[*mediatypes.Episode]{}, nil
	}

	// Get episodes for each season and combine them
	var allEpisodes []*models.MediaItem[*mediatypes.Episode]
	for _, season := range seasons {
		// Skip season 0 (specials) for now as it might not exist in all clients
		if season.Data.Number == 0 {
			continue
		}

		seasonEpisodes, err := showProvider.GetSeriesEpisodesBySeasonNbr(ctx, seriesID, season.Data.Number)
		if err != nil {
			log.Warn().
				Err(err).
				Uint64("clientID", clientID).
				Str("seriesID", seriesID).
				Int("seasonNumber", season.Data.Number).
				Msg("Failed to get episodes for season, continuing with other seasons")
			continue
		}

		allEpisodes = append(allEpisodes, seasonEpisodes...)
	}

	return allEpisodes, nil
}

// GetEpisodesBySeasonID gets all episodes for a specific season
func (s *clientSeriesService[T]) GetEpisodesBySeasonNbr(ctx context.Context, clientID uint64, seriesID string, seasonNumber int) ([]*models.MediaItem[*mediatypes.Episode], error) {
	client, err := s.getSeriesProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	showProvider, ok := client.(providers.SeriesProvider)
	if !ok {
		return nil, fmt.Errorf("client does not implement series provider interface")
	}

	// Get episodes for this series
	episodes, err := showProvider.GetSeriesEpisodesBySeasonNbr(ctx, seriesID, seasonNumber)
	if err != nil {
		return nil, err
	}

	return episodes, nil
}
