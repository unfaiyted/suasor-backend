package services

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"suasor/client/media"
	"suasor/client/media/providers"
	mediatypes "suasor/client/media/types"
	"suasor/client/types"
	"suasor/repository"
	"suasor/types/models"
)

var ErrUnsupportedFeature = errors.New("feature not supported by this media client")

// MediaClientMovieService defines operations for interacting with movie clients
type MediaClientMovieService interface {
	GetMovieByID(ctx context.Context, userID uint64, clientID uint64, movieID string) (*models.MediaItem[mediatypes.Movie], error)
	GetMoviesByGenre(ctx context.Context, userID uint64, genre string) ([]models.MediaItem[mediatypes.Movie], error)
	GetMoviesByYear(ctx context.Context, userID uint64, year int) ([]models.MediaItem[mediatypes.Movie], error)
	GetMoviesByActor(ctx context.Context, userID uint64, actor string) ([]models.MediaItem[mediatypes.Movie], error)
	GetMoviesByDirector(ctx context.Context, userID uint64, director string) ([]models.MediaItem[mediatypes.Movie], error)
	GetMoviesByRating(ctx context.Context, userID uint64, minRating, maxRating float64) ([]models.MediaItem[mediatypes.Movie], error)
	GetLatestMoviesByAdded(ctx context.Context, userID uint64, count int) ([]models.MediaItem[mediatypes.Movie], error)
	GetPopularMovies(ctx context.Context, userID uint64, count int) ([]models.MediaItem[mediatypes.Movie], error)
	GetTopRatedMovies(ctx context.Context, userID uint64, count int) ([]models.MediaItem[mediatypes.Movie], error)
	SearchMovies(ctx context.Context, userID uint64, query string) ([]models.MediaItem[mediatypes.Movie], error)
}

type mediaMovieService struct {
	clientRepo    repository.ClientRepository[types.MediaClientConfig]
	clientFactory media.ClientFactory
}

// NewMediaClientMovieService creates a new media movie service
func NewMediaClientMovieService(
	clientRepo repository.ClientRepository[types.MediaClientConfig],
	clientFactory media.ClientFactory,
) MediaClientMovieService {
	return &mediaMovieService{
		clientRepo:    clientRepo,
		clientFactory: clientFactory,
	}
}

// getMovieClients gets all movie clients for a user
func (s *mediaMovieService) getMovieClients(ctx context.Context, userID uint64) ([]media.MediaClient, error) {
	// Get all media clients for the user
	clients, err := s.clientRepo.GetByCategory(ctx, types.ClientCategoryMedia, userID)
	if err != nil {
		return nil, err
	}

	var movieClients []media.MediaClient

	// Filter and instantiate clients that support movies
	for _, clientConfig := range clients {
		if clientConfig.Config.Data.SupportsMovies() {
			clientId := clientConfig.GetID()
			client, err := s.clientFactory.GetMediaClient(ctx, clientId, clientConfig.Config.Data)
			if err != nil {
				// Log error but continue with other clients
				continue
			}
			movieClients = append(movieClients, client)
		}
	}

	return movieClients, nil
}

// getSpecificMovieClient gets a specific movie client
func (s *mediaMovieService) getSpecificMovieClient(ctx context.Context, userID, clientID uint64) (media.MediaClient, error) {
	clientConfig, err := s.clientRepo.GetByID(ctx, clientID, userID)
	if err != nil {
		return nil, err
	}

	if !clientConfig.Config.Data.SupportsMovies() {
		return nil, ErrUnsupportedFeature
	}

	return s.clientFactory.GetMediaClient(ctx, clientID, clientConfig.Config.Data)
}

func (s *mediaMovieService) GetMovieByID(ctx context.Context, userID uint64, clientID uint64, movieID string) (*models.MediaItem[mediatypes.Movie], error) {
	client, err := s.getSpecificMovieClient(ctx, userID, clientID)
	if err != nil {
		return nil, err
	}

	// This is a simplified example - actual implementation would depend on your client interface
	movieProvider, ok := client.(providers.MovieProvider)
	if !ok {
		return nil, ErrUnsupportedFeature
	}

	movie, err := movieProvider.GetMovieByID(ctx, movieID)
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (s *mediaMovieService) GetMoviesByGenre(ctx context.Context, userID uint64, genre string) ([]models.MediaItem[mediatypes.Movie], error) {
	clients, err := s.getMovieClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allMovies []models.MediaItem[mediatypes.Movie]

	// Query each client and aggregate results
	for _, client := range clients {
		movieProvider, ok := client.(providers.MovieProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Filters: map[string]string{
				"genre": genre,
			},
		}

		movies, err := movieProvider.GetMovies(ctx, options)
		if err != nil {
			// Log error but continue with other clients
			continue
		}

		allMovies = append(allMovies, movies...)
	}

	return allMovies, nil
}

func (s *mediaMovieService) GetMoviesByYear(ctx context.Context, userID uint64, year int) ([]models.MediaItem[mediatypes.Movie], error) {
	clients, err := s.getMovieClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allMovies []models.MediaItem[mediatypes.Movie]

	for _, client := range clients {
		movieProvider, ok := client.(providers.MovieProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Filters: map[string]string{
				"year": strconv.Itoa(year),
			},
		}

		movies, err := movieProvider.GetMovies(ctx, options)
		if err != nil {
			continue
		}

		allMovies = append(allMovies, movies...)
	}

	return allMovies, nil
}

func (s *mediaMovieService) GetMoviesByActor(ctx context.Context, userID uint64, actor string) ([]models.MediaItem[mediatypes.Movie], error) {
	clients, err := s.getMovieClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allMovies []models.MediaItem[mediatypes.Movie]

	for _, client := range clients {
		movieProvider, ok := client.(providers.MovieProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Filters: map[string]string{
				"actor": actor,
			},
		}

		movies, err := movieProvider.GetMovies(ctx, options)
		if err != nil {
			continue
		}

		allMovies = append(allMovies, movies...)
	}

	return allMovies, nil
}

func (s *mediaMovieService) GetMoviesByDirector(ctx context.Context, userID uint64, director string) ([]models.MediaItem[mediatypes.Movie], error) {
	clients, err := s.getMovieClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allMovies []models.MediaItem[mediatypes.Movie]

	for _, client := range clients {
		movieProvider, ok := client.(providers.MovieProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Filters: map[string]string{
				"director": director,
			},
		}

		movies, err := movieProvider.GetMovies(ctx, options)
		if err != nil {
			continue
		}

		allMovies = append(allMovies, movies...)
	}

	return allMovies, nil
}

func (s *mediaMovieService) GetMoviesByRating(ctx context.Context, userID uint64, minRating, maxRating float64) ([]models.MediaItem[mediatypes.Movie], error) {
	clients, err := s.getMovieClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allMovies []models.MediaItem[mediatypes.Movie]

	for _, client := range clients {
		movieProvider, ok := client.(providers.MovieProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Filters: map[string]string{
				"minRating": strconv.FormatFloat(minRating, 'f', -1, 64),
				"maxRating": strconv.FormatFloat(maxRating, 'f', -1, 64),
			},
		}

		movies, err := movieProvider.GetMovies(ctx, options)
		if err != nil {
			continue
		}

		allMovies = append(allMovies, movies...)
	}

	return allMovies, nil
}

func (s *mediaMovieService) GetLatestMoviesByAdded(ctx context.Context, userID uint64, count int) ([]models.MediaItem[mediatypes.Movie], error) {
	clients, err := s.getMovieClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allMovies []models.MediaItem[mediatypes.Movie]

	for _, client := range clients {
		movieProvider, ok := client.(providers.MovieProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Sort:      "added",
			SortOrder: mediatypes.SortOrderDesc,
			Limit:     count,
		}

		movies, err := movieProvider.GetMovies(ctx, options)
		if err != nil {
			continue
		}

		allMovies = append(allMovies, movies...)
	}

	// Sort by added date again since we're combining from multiple sources
	sort.Slice(allMovies, func(i, j int) bool {
		return allMovies[i].Data.GetDetails().AddedAt.After(allMovies[j].Data.GetDetails().AddedAt)
	})

	// Limit to requested count
	if len(allMovies) > count {
		allMovies = allMovies[:count]
	}

	return allMovies, nil
}

func (s *mediaMovieService) GetPopularMovies(ctx context.Context, userID uint64, count int) ([]models.MediaItem[mediatypes.Movie], error) {
	clients, err := s.getMovieClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allMovies []models.MediaItem[mediatypes.Movie]

	for _, client := range clients {
		movieProvider, ok := client.(providers.MovieProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Sort:      "popularity",
			SortOrder: mediatypes.SortOrderDesc,
			Limit:     count,
		}

		movies, err := movieProvider.GetMovies(ctx, options)
		if err != nil {
			continue
		}

		allMovies = append(allMovies, movies...)
	}

	// Sort by popularity
	// sort.Slice(allMovies, func(i, j int) bool {
	// 	return allMovies[i].Data.GetDetails().Popularity > allMovies[j].Data.GetDetails().Popularity
	// })

	// Limit to requested count
	if len(allMovies) > count {
		allMovies = allMovies[:count]
	}

	return allMovies, nil
}

func (s *mediaMovieService) GetTopRatedMovies(ctx context.Context, userID uint64, count int) ([]models.MediaItem[mediatypes.Movie], error) {
	clients, err := s.getMovieClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allMovies []models.MediaItem[mediatypes.Movie]

	for _, client := range clients {
		movieProvider, ok := client.(providers.MovieProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Sort:      "rating",
			SortOrder: mediatypes.SortOrderDesc,
			Limit:     count,
		}

		movies, err := movieProvider.GetMovies(ctx, options)
		if err != nil {
			continue
		}

		allMovies = append(allMovies, movies...)
	}

	// Sort by rating
	// sort.Slice(allMovies, func(i, j int) bool {
	// 	return allMovies[i].Data.GetDetails().Rating > allMovies[j].Data.GetDetails().Rating
	// })

	// Limit to requested count
	if len(allMovies) > count {
		allMovies = allMovies[:count]
	}

	return allMovies, nil
}

func (s *mediaMovieService) SearchMovies(ctx context.Context, userID uint64, query string) ([]models.MediaItem[mediatypes.Movie], error) {
	clients, err := s.getMovieClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allMovies []models.MediaItem[mediatypes.Movie]

	for _, client := range clients {
		movieProvider, ok := client.(providers.MovieProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Query: query,
		}

		movies, err := movieProvider.GetMovies(ctx, options)
		if err != nil {
			continue
		}

		allMovies = append(allMovies, movies...)
	}

	return allMovies, nil
}
