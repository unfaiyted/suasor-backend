package services

import (
	"context"
	"sort"
	"strconv"
	"suasor/client"
	"suasor/client/media"
	"suasor/client/media/providers"
	mediatypes "suasor/client/media/types"
	clienttypes "suasor/client/types"
	"suasor/repository"
	"suasor/types/models"
)

// MediaClientSeriesService defines operations for interacting with TV show clients
type MediaClientSeriesService[T interface {
	clienttypes.MediaClientConfig
}] interface {
	GetSeriesByID(ctx context.Context, userID uint64, clientID uint64, seriesID string) (*models.MediaItem[mediatypes.Series], error)
	GetSeriesByName(ctx context.Context, userID uint64, name string) ([]models.MediaItem[mediatypes.Series], error)
	GetSeasonsBySeriesID(ctx context.Context, userID uint64, clientID uint64, seriesID string) ([]models.MediaItem[mediatypes.Season], error)
	GetSeriesByGenre(ctx context.Context, userID uint64, genre string) ([]models.MediaItem[mediatypes.Series], error)
	GetSeriesByYear(ctx context.Context, userID uint64, year int) ([]models.MediaItem[mediatypes.Series], error)
	GetSeriesByActor(ctx context.Context, userID uint64, actor string) ([]models.MediaItem[mediatypes.Series], error)
	GetSeriesByCreator(ctx context.Context, userID uint64, creator string) ([]models.MediaItem[mediatypes.Series], error)
	GetSeriesByRating(ctx context.Context, userID uint64, minRating, maxRating float64) ([]models.MediaItem[mediatypes.Series], error)
	GetLatestSeriesByAdded(ctx context.Context, userID uint64, count int) ([]models.MediaItem[mediatypes.Series], error)
	GetPopularSeries(ctx context.Context, userID uint64, count int) ([]models.MediaItem[mediatypes.Series], error)
	GetTopRatedSeries(ctx context.Context, userID uint64, count int) ([]models.MediaItem[mediatypes.Series], error)
	SearchSeries(ctx context.Context, userID uint64, query string) ([]models.MediaItem[mediatypes.Series], error)
}

type mediaSeriesService[T interface {
	clienttypes.MediaClientConfig
}] struct {
	clientRepo    repository.ClientRepository[T]
	clientFactory *client.ClientFactoryService
}

// NewMediaClientSeriesService creates a new media TV show service
func NewMediaClientSeriesService[T interface {
	clienttypes.MediaClientConfig
}](
	clientRepo repository.ClientRepository[T],
	clientFactory *client.ClientFactoryService,
) MediaClientSeriesService[T] {
	return &mediaSeriesService[T]{
		clientRepo:    clientRepo,
		clientFactory: clientFactory,
	}
}

// getSeriesClients gets all TV show clients for a user
func (s *mediaSeriesService[T]) getSeriesClients(ctx context.Context, userID uint64) ([]media.MediaClient, error) {
	// Get all media clients for the user
	clients, err := s.clientRepo.GetByCategory(ctx, clienttypes.ClientCategoryMedia, userID)
	if err != nil {
		return nil, err
	}

	var showClients []media.MediaClient

	// Filter and instantiate clients that support TV shows
	for _, clientConfig := range clients {
		// Check if the config supports series
		if config, ok := any(clientConfig.Config.Data).(clienttypes.MediaClientConfig); ok && config.SupportsSeries() {
			clientId := clientConfig.GetID()
			client, err := s.clientFactory.GetClient(ctx, clientId, clientConfig.Config.Data.GetType())
			if err != nil {
				// Log error but continue with other clients
				continue
			}

			if mediaClient, ok := client.(media.MediaClient); ok {
				showClients = append(showClients, mediaClient)
			}
		}
	}

	return showClients, nil
}

// getSpecificSeriesClient gets a specific TV show client
func (s *mediaSeriesService[T]) getSpecificSeriesClient(ctx context.Context, userID, clientID uint64) (media.MediaClient, error) {
	clientConfig, err := s.clientRepo.GetByID(ctx, clientID, userID)
	if err != nil {
		return nil, err
	}

	// Check if the config supports series
	mediaConfig, ok := any(clientConfig.Config.Data).(clienttypes.MediaClientConfig)
	if !ok || !mediaConfig.SupportsSeries() {
		return nil, ErrUnsupportedFeature
	}

	client, err := s.clientFactory.GetClient(ctx, clientID, clientConfig.Config.Data.GetType())
	if err != nil {
		return nil, err
	}

	showClient, ok := client.(media.MediaClient)
	if !ok {
		return nil, ErrUnsupportedFeature
	}
	return showClient, nil
}

func (s *mediaSeriesService[T]) GetSeriesByID(ctx context.Context, userID uint64, clientID uint64, seriesID string) (*models.MediaItem[mediatypes.Series], error) {
	client, err := s.getSpecificSeriesClient(ctx, userID, clientID)
	if err != nil {
		return nil, err
	}

	showProvider, ok := client.(providers.SeriesProvider)
	if !ok {
		return nil, ErrUnsupportedFeature
	}

	series, err := showProvider.GetSeriesByID(ctx, seriesID)
	if err != nil {
		return nil, err
	}
	return &series, nil
}

func (s *mediaSeriesService[T]) GetSeriesByName(ctx context.Context, userID uint64, name string) ([]models.MediaItem[mediatypes.Series], error) {
	clients, err := s.getSeriesClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allSeries []models.MediaItem[mediatypes.Series]

	for _, client := range clients {
		showProvider, ok := client.(providers.SeriesProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Filters: map[string]string{
				"name": name,
			},
		}

		series, err := showProvider.GetSeries(ctx, options)
		if err != nil {
			continue
		}

		allSeries = append(allSeries, series...)
	}

	return allSeries, nil
}

func (s *mediaSeriesService[T]) GetSeasonsBySeriesID(ctx context.Context, userID uint64, clientID uint64, seriesID string) ([]models.MediaItem[mediatypes.Season], error) {
	client, err := s.getSpecificSeriesClient(ctx, userID, clientID)
	if err != nil {
		return nil, err
	}

	showProvider, ok := client.(providers.SeriesProvider)
	if !ok {
		return nil, ErrUnsupportedFeature
	}

	seasons, err := showProvider.GetSeriesSeasons(ctx, seriesID)
	if err != nil {
		return nil, err
	}

	return seasons, nil
}

func (s *mediaSeriesService[T]) GetSeriesByGenre(ctx context.Context, userID uint64, genre string) ([]models.MediaItem[mediatypes.Series], error) {
	clients, err := s.getSeriesClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allSeries []models.MediaItem[mediatypes.Series]

	for _, client := range clients {
		showProvider, ok := client.(providers.SeriesProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Filters: map[string]string{
				"genre": genre,
			},
		}

		shows, err := showProvider.GetSeries(ctx, options)
		if err != nil {
			continue
		}

		allSeries = append(allSeries, shows...)
	}

	return allSeries, nil
}

func (s *mediaSeriesService[T]) GetSeriesByYear(ctx context.Context, userID uint64, year int) ([]models.MediaItem[mediatypes.Series], error) {
	clients, err := s.getSeriesClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allSeries []models.MediaItem[mediatypes.Series]

	for _, client := range clients {
		showProvider, ok := client.(providers.SeriesProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Filters: map[string]string{
				"year": strconv.Itoa(year),
			},
		}

		shows, err := showProvider.GetSeries(ctx, options)
		if err != nil {
			continue
		}

		allSeries = append(allSeries, shows...)
	}

	return allSeries, nil
}

func (s *mediaSeriesService[T]) GetSeriesByActor(ctx context.Context, userID uint64, actor string) ([]models.MediaItem[mediatypes.Series], error) {
	clients, err := s.getSeriesClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allSeries []models.MediaItem[mediatypes.Series]

	for _, client := range clients {
		showProvider, ok := client.(providers.SeriesProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Filters: map[string]string{
				"actor": actor,
			},
		}

		shows, err := showProvider.GetSeries(ctx, options)
		if err != nil {
			continue
		}

		allSeries = append(allSeries, shows...)
	}

	return allSeries, nil
}

func (s *mediaSeriesService[T]) GetSeriesByCreator(ctx context.Context, userID uint64, creator string) ([]models.MediaItem[mediatypes.Series], error) {
	clients, err := s.getSeriesClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allSeries []models.MediaItem[mediatypes.Series]

	for _, client := range clients {
		showProvider, ok := client.(providers.SeriesProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Filters: map[string]string{
				"creator": creator,
			},
		}

		shows, err := showProvider.GetSeries(ctx, options)
		if err != nil {
			continue
		}

		allSeries = append(allSeries, shows...)
	}

	return allSeries, nil
}

func (s *mediaSeriesService[T]) GetSeriesByRating(ctx context.Context, userID uint64, minRating, maxRating float64) ([]models.MediaItem[mediatypes.Series], error) {
	clients, err := s.getSeriesClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allSeries []models.MediaItem[mediatypes.Series]

	for _, client := range clients {
		showProvider, ok := client.(providers.SeriesProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Filters: map[string]string{
				"minRating": strconv.FormatFloat(minRating, 'f', -1, 64),
				"maxRating": strconv.FormatFloat(maxRating, 'f', -1, 64),
			},
		}

		shows, err := showProvider.GetSeries(ctx, options)
		if err != nil {
			continue
		}

		allSeries = append(allSeries, shows...)
	}

	return allSeries, nil
}

func (s *mediaSeriesService[T]) GetLatestSeriesByAdded(ctx context.Context, userID uint64, count int) ([]models.MediaItem[mediatypes.Series], error) {
	clients, err := s.getSeriesClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allSeries []models.MediaItem[mediatypes.Series]

	for _, client := range clients {
		showProvider, ok := client.(providers.SeriesProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Sort:      "added",
			SortOrder: mediatypes.SortOrderDesc,
			Limit:     count,
		}

		shows, err := showProvider.GetSeries(ctx, options)
		if err != nil {
			continue
		}

		allSeries = append(allSeries, shows...)
	}

	// Sort by added date again since we're combining from multiple sources
	sort.Slice(allSeries, func(i, j int) bool {
		return allSeries[i].Data.GetDetails().AddedAt.After(allSeries[j].Data.GetDetails().AddedAt)
	})

	// Limit to requested count
	if len(allSeries) > count {
		allSeries = allSeries[:count]
	}

	return allSeries, nil
}

func (s *mediaSeriesService[T]) GetPopularSeries(ctx context.Context, userID uint64, count int) ([]models.MediaItem[mediatypes.Series], error) {
	clients, err := s.getSeriesClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allSeries []models.MediaItem[mediatypes.Series]

	for _, client := range clients {
		showProvider, ok := client.(providers.SeriesProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Sort:      "popularity",
			SortOrder: mediatypes.SortOrderDesc,
			Limit:     count,
		}

		shows, err := showProvider.GetSeries(ctx, options)
		if err != nil {
			continue
		}

		allSeries = append(allSeries, shows...)
	}

	// Limit to requested count
	if len(allSeries) > count {
		allSeries = allSeries[:count]
	}

	return allSeries, nil
}

func (s *mediaSeriesService[T]) GetTopRatedSeries(ctx context.Context, userID uint64, count int) ([]models.MediaItem[mediatypes.Series], error) {
	clients, err := s.getSeriesClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allSeries []models.MediaItem[mediatypes.Series]

	for _, client := range clients {
		showProvider, ok := client.(providers.SeriesProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Sort:      "rating",
			SortOrder: mediatypes.SortOrderDesc,
			Limit:     count,
		}

		shows, err := showProvider.GetSeries(ctx, options)
		if err != nil {
			continue
		}

		allSeries = append(allSeries, shows...)
	}

	// Limit to requested count
	if len(allSeries) > count {
		allSeries = allSeries[:count]
	}

	return allSeries, nil
}

func (s *mediaSeriesService[T]) SearchSeries(ctx context.Context, userID uint64, query string) ([]models.MediaItem[mediatypes.Series], error) {
	clients, err := s.getSeriesClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allSeries []models.MediaItem[mediatypes.Series]

	for _, client := range clients {
		showProvider, ok := client.(providers.SeriesProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Query: query,
		}

		shows, err := showProvider.GetSeries(ctx, options)
		if err != nil {
			continue
		}

		allSeries = append(allSeries, shows...)
	}

	return allSeries, nil
}
