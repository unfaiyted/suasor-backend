package services

import (
	"context"
	"errors"
	"sort"

	"suasor/clients"
	"suasor/clients/media/providers"
	mediatypes "suasor/clients/media/types"
	"suasor/clients/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/utils/logger"
)

var ErrUnsupportedFeature = errors.New("feature not supported by this media client")

// ClientMovieService defines operations for interacting with movie clients
type ClientMovieService[T types.ClientConfig] interface {
	// MovieService[T]
	GetClientMovieByItemID(ctx context.Context, clientID uint64, itemID string) (*models.MediaItem[*mediatypes.Movie], error)
	GetClientMoviesByGenre(ctx context.Context, clientID uint64, genre string) ([]*models.MediaItem[*mediatypes.Movie], error)
	GetClientMoviesByYear(ctx context.Context, clientID uint64, year int) ([]*models.MediaItem[*mediatypes.Movie], error)
	GetClientMoviesByActor(ctx context.Context, clientID uint64, actor string) ([]*models.MediaItem[*mediatypes.Movie], error)
	GetClientMoviesByDirector(ctx context.Context, clientID uint64, director string) ([]*models.MediaItem[*mediatypes.Movie], error)
	GetClientMoviesByRating(ctx context.Context, clientID uint64, minRating, maxRating float64) ([]*models.MediaItem[*mediatypes.Movie], error)
	GetClientMoviesLatestByAdded(ctx context.Context, clientID uint64, count int) ([]*models.MediaItem[*mediatypes.Movie], error)
	GetClientPopularMovies(ctx context.Context, clientID uint64, count int) ([]*models.MediaItem[*mediatypes.Movie], error)
	GetClientTopRatedMovies(ctx context.Context, clientID uint64, count int) ([]*models.MediaItem[*mediatypes.Movie], error)
	SearchClientMovies(ctx context.Context, clientID uint64, query *mediatypes.QueryOptions) ([]*models.MediaItem[*mediatypes.Movie], error)
	SearchAllClientsMovies(ctx context.Context, query *mediatypes.QueryOptions) ([]*models.MediaItem[*mediatypes.Movie], error)
}

type clientMovieService[T types.ClientMediaConfig] struct {
	clientRepo    repository.ClientRepository[T]
	clientFactory *clients.ClientProviderFactoryService
}

// NewClientMovieService creates a new media movie service
func NewClientMovieService[T types.ClientMediaConfig](
	clientRepo repository.ClientRepository[T],
	factory *clients.ClientProviderFactoryService,
) ClientMovieService[T] {
	return &clientMovieService[T]{
		clientRepo:    clientRepo,
		clientFactory: factory,
	}
}

func (s *clientMovieService[T]) GetClientMovieByItemID(ctx context.Context, clientID uint64, movieID string) (*models.MediaItem[*mediatypes.Movie], error) {
	// Get the client
	log := logger.LoggerFromContext(ctx)

	client, err := s.clientRepo.GetByID(ctx, clientID)
	log.Info().
		Uint64("clientID", clientID).
		Msg("Retrieved client")
	if err != nil {
		return nil, err
	}
	log.Info().
		Uint64("clientID", clientID).
		Str("movieID", movieID).
		Msg("Retrieving movie")

	provider, err := s.clientFactory.GetMovieProvider(ctx, clientID, client.Config.Data)

	movie, err := provider.GetMovieByID(ctx, movieID)
	log.Info().
		Uint64("clientID", clientID).
		Str("movieID", movieID).
		Msg("Retrieved movie")
	if err != nil {
		return nil, err
	}
	return movie, nil
}
func (s *clientMovieService[T]) GetMoviesByGenre(ctx context.Context, userID uint64, genre string) ([]*models.MediaItem[*mediatypes.Movie], error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("userID", userID).
		Str("genre", genre).
		Msg("Retrieving movies by genre")
	clients, err := s.getMovieProviders(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allMovies []*models.MediaItem[*mediatypes.Movie]

	log.Info().
		Uint64("userID", userID).
		Str("genre", genre).
		Int("count", len(allMovies)).
		Msg("Movies retrieved successfully")
	// Query each client and aggregate results
	for _, client := range clients {
		movieProvider, ok := client.(providers.MovieProvider)
		if !ok {
			log.Warn().
				Uint64("userID", userID).
				Str("genre", genre).
				Msg("Client does not support movies")
			continue
		}

		log.Debug().
			Uint64("userID", userID).
			Str("genre", genre).
			Msg("Retrieving movies by genre")
		options := &mediatypes.QueryOptions{
			Genre: genre,
		}

		movies, err := movieProvider.GetMovies(ctx, options)
		if err != nil {
			// Log error but continue with other clients
			continue
		}

		log.Debug().
			Uint64("userID", userID).
			Str("genre", genre).
			Int("count", len(movies)).
			Msg("Movies retrieved successfully")
		allMovies = append(allMovies, movies...)
	}

	return allMovies, nil
}
func (s *clientMovieService[T]) GetMoviesByYear(ctx context.Context, userID uint64, year int) ([]*models.MediaItem[*mediatypes.Movie], error) {
	providers, err := s.getMovieProviders(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allMovies []*models.MediaItem[*mediatypes.Movie]

	for _, provider := range providers {

		options := &mediatypes.QueryOptions{
			Year: year,
		}

		// TODO: run this in parallel
		movies, err := provider.GetMovies(ctx, options)
		if err != nil {
			continue
		}

		allMovies = append(allMovies, movies...)
	}

	return allMovies, nil
}
func (s *clientMovieService[T]) GetMoviesByActor(ctx context.Context, userID uint64, actor string) ([]*models.MediaItem[*mediatypes.Movie], error) {
	clients, err := s.getMovieProviders(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allMovies []*models.MediaItem[*mediatypes.Movie]

	for _, client := range clients {
		movieProvider, ok := client.(providers.MovieProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Actor: actor,
		}

		movies, err := movieProvider.GetMovies(ctx, options)
		if err != nil {
			continue
		}

		allMovies = append(allMovies, movies...)
	}

	return allMovies, nil
}
func (s *clientMovieService[T]) GetMoviesByDirector(ctx context.Context, userID uint64, director string) ([]*models.MediaItem[*mediatypes.Movie], error) {
	clients, err := s.getMovieProviders(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allMovies []*models.MediaItem[*mediatypes.Movie]

	for _, client := range clients {
		movieProvider, ok := client.(providers.MovieProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Director: director,
		}

		movies, err := movieProvider.GetMovies(ctx, options)
		if err != nil {
			continue
		}

		allMovies = append(allMovies, movies...)
	}

	return allMovies, nil
}
func (s *clientMovieService[T]) GetMoviesByRating(ctx context.Context, userID uint64, minRating, maxRating float64) ([]*models.MediaItem[*mediatypes.Movie], error) {
	clients, err := s.getMovieProviders(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allMovies []*models.MediaItem[*mediatypes.Movie]

	for _, client := range clients {
		movieProvider, ok := client.(providers.MovieProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			MinimumRating: float32(minRating),
			// Note: maxRating doesn't have a typed field yet, will need to be added if needed
		}

		movies, err := movieProvider.GetMovies(ctx, options)
		if err != nil {
			continue
		}

		allMovies = append(allMovies, movies...)
	}

	return allMovies, nil
}
func (s *clientMovieService[T]) GetLatestMoviesByAdded(ctx context.Context, userID uint64, count int) ([]*models.MediaItem[*mediatypes.Movie], error) {
	clients, err := s.getMovieProviders(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allMovies []*models.MediaItem[*mediatypes.Movie]

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
func (s *clientMovieService[T]) GetTopRatedMovies(ctx context.Context, userID uint64, count int) ([]*models.MediaItem[*mediatypes.Movie], error) {
	clients, err := s.getMovieProviders(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allMovies []*models.MediaItem[*mediatypes.Movie]

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
func (s *clientMovieService[T]) SearchClientMovies(ctx context.Context, clientID uint64, query *mediatypes.QueryOptions) ([]*models.MediaItem[*mediatypes.Movie], error) {
	client, err := s.clientRepo.GetByID(ctx, clientID)
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("clientID", clientID).
		Msg("Retrieved client")
	if err != nil {
		return nil, err
	}
	log.Info().
		Uint64("clientID", clientID).
		Msg("Searching movies")

	movieProvider, err := s.clientFactory.GetMovieProvider(ctx, clientID, client.Config.Data)
	if err != nil {
		return nil, err
	}

	movies, err := movieProvider.GetMovies(ctx, query)
	if err != nil {
		return nil, err
	}

	return movies, nil
}
func (s *clientMovieService[T]) GetMovieByClientItemID(ctx context.Context, clientID uint64, movieID string) (*models.MediaItem[*mediatypes.Movie], error) {
	client, err := s.clientRepo.GetByID(ctx, clientID)
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("clientID", clientID).
		Msg("Retrieved client")
	if err != nil {
		return nil, err
	}
	log.Info().
		Uint64("clientID", clientID).
		Str("movieID", movieID).
		Msg("Retrieving movie")

	movieProvider, err := s.clientFactory.GetMovieProvider(ctx, clientID, client.Config.Data)
	if err != nil {
		return nil, err
	}

	movie, err := movieProvider.GetMovieByID(ctx, movieID)
	log.Info().
		Uint64("clientID", clientID).
		Str("movieID", movieID).
		Msg("Retrieved movie")
	if err != nil {
		return nil, err
	}
	return movie, nil
}
func (s *clientMovieService[T]) GetClientMoviesByActor(ctx context.Context, clientID uint64, actor string) ([]*models.MediaItem[*mediatypes.Movie], error) {
	client, err := s.clientRepo.GetByID(ctx, clientID)
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("clientID", clientID).
		Msg("Retrieved client")
	if err != nil {
		return nil, err
	}
	log.Info().
		Uint64("clientID", clientID).
		Str("actor", actor).
		Msg("Retrieving movies by actor")

	movieProvider, err := s.clientFactory.GetMovieProvider(ctx, clientID, client.Config.Data)
	if err != nil {
		return nil, err
	}

	movies, err := movieProvider.GetMovies(ctx, &mediatypes.QueryOptions{
		Actor: actor,
	})
	if err != nil {
		return nil, err
	}

	return movies, nil
}
func (s *clientMovieService[T]) GetClientMoviesByDirector(ctx context.Context, clientID uint64, director string) ([]*models.MediaItem[*mediatypes.Movie], error) {
	client, err := s.clientRepo.GetByID(ctx, clientID)
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("clientID", clientID).
		Msg("Retrieved client")
	if err != nil {
		return nil, err
	}
	log.Info().
		Uint64("clientID", clientID).
		Str("director", director).
		Msg("Retrieving movies by director")

	movieProvider, err := s.clientFactory.GetMovieProvider(ctx, clientID, client.Config.Data)
	if err != nil {
		return nil, err
	}

	movies, err := movieProvider.GetMovies(ctx, &mediatypes.QueryOptions{
		Director: director,
	})
	if err != nil {
		return nil, err
	}

	return movies, nil
}
func (s *clientMovieService[T]) GetClientMoviesByGenre(ctx context.Context, clientID uint64, genre string) ([]*models.MediaItem[*mediatypes.Movie], error) {
	client, err := s.clientRepo.GetByID(ctx, clientID)
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("clientID", clientID).
		Msg("Retrieved client")
	if err != nil {
		return nil, err
	}
	log.Info().
		Uint64("clientID", clientID).
		Str("genre", genre).
		Msg("Retrieving movies by genre")

	movieProvider, err := s.clientFactory.GetMovieProvider(ctx, clientID, client.Config.Data)
	if err != nil {
		return nil, err
	}

	movies, err := movieProvider.GetMovies(ctx, &mediatypes.QueryOptions{
		Genre: genre,
	})
	if err != nil {
		return nil, err
	}

	return movies, nil
}
func (s *clientMovieService[T]) GetClientMoviesByRating(ctx context.Context, clientID uint64, minRating, maxRating float64) ([]*models.MediaItem[*mediatypes.Movie], error) {
	client, err := s.clientRepo.GetByID(ctx, clientID)
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("clientID", clientID).
		Msg("Retrieved client")
	if err != nil {
		return nil, err
	}
	log.Info().
		Uint64("clientID", clientID).
		Msg("Retrieving movies by rating")

	movieProvider, err := s.clientFactory.GetMovieProvider(ctx, clientID, client.Config.Data)
	if err != nil {
		return nil, err
	}

	movies, err := movieProvider.GetMovies(ctx, &mediatypes.QueryOptions{
		MinimumRating: float32(minRating),
		MaximumRating: float32(maxRating),
	})
	if err != nil {
		return nil, err
	}

	return movies, nil
}
func (s *clientMovieService[T]) GetClientMoviesByYear(ctx context.Context, clientID uint64, year int) ([]*models.MediaItem[*mediatypes.Movie], error) {
	client, err := s.clientRepo.GetByID(ctx, clientID)
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("clientID", clientID).
		Msg("Retrieved client")
	if err != nil {
		return nil, err
	}
	log.Info().
		Uint64("clientID", clientID).
		Int("year", year).
		Msg("Retrieving movies by year")

	movieProvider, err := s.clientFactory.GetMovieProvider(ctx, clientID, client.Config.Data)
	if err != nil {
		return nil, err
	}

	movies, err := movieProvider.GetMovies(ctx, &mediatypes.QueryOptions{
		Year: year,
	})
	if err != nil {
		return nil, err
	}

	return movies, nil
}
func (s *clientMovieService[T]) GetClientMoviesLatestByAdded(ctx context.Context, clientID uint64, count int) ([]*models.MediaItem[*mediatypes.Movie], error) {
	client, err := s.clientRepo.GetByID(ctx, clientID)
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("clientID", clientID).
		Msg("Retrieved client")
	if err != nil {
		return nil, err
	}
	log.Info().
		Uint64("clientID", clientID).
		Int("count", count).
		Msg("Retrieving movies by latest added")

	movieProvider, err := s.clientFactory.GetMovieProvider(ctx, clientID, client.Config.Data)
	if err != nil {
		return nil, err
	}

	movies, err := movieProvider.GetMovies(ctx, &mediatypes.QueryOptions{
		Sort:      "added",
		SortOrder: mediatypes.SortOrderDesc,
		Limit:     count,
	})
	if err != nil {
		return nil, err
	}

	return movies, nil
}
func (s *clientMovieService[T]) GetClientPopularMovies(ctx context.Context, clientID uint64, count int) ([]*models.MediaItem[*mediatypes.Movie], error) {
	client, err := s.clientRepo.GetByID(ctx, clientID)
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("clientID", clientID).
		Msg("Retrieved client")
	if err != nil {
		return nil, err
	}
	log.Info().
		Uint64("clientID", clientID).
		Int("count", count).
		Msg("Retrieving movies by popularity")

	movieProvider, err := s.clientFactory.GetMovieProvider(ctx, clientID, client.Config.Data)
	if err != nil {
		return nil, err
	}

	movies, err := movieProvider.GetMovies(ctx, &mediatypes.QueryOptions{
		Sort:      "popularity",
		SortOrder: mediatypes.SortOrderDesc,
		Limit:     count,
	})
	if err != nil {
		return nil, err
	}

	return movies, nil
}
func (s *clientMovieService[T]) GetClientTopRatedMovies(ctx context.Context, clientID uint64, count int) ([]*models.MediaItem[*mediatypes.Movie], error) {
	client, err := s.clientRepo.GetByID(ctx, clientID)
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("clientID", clientID).
		Msg("Retrieved client")
	if err != nil {
		return nil, err
	}
	log.Info().
		Uint64("clientID", clientID).
		Int("count", count).
		Msg("Retrieving movies by rating")

	movieProvider, err := s.clientFactory.GetMovieProvider(ctx, clientID, client.Config.Data)
	if err != nil {
		return nil, err
	}

	movies, err := movieProvider.GetMovies(ctx, &mediatypes.QueryOptions{
		Sort:      "rating",
		SortOrder: mediatypes.SortOrderDesc,
		Limit:     count,
	})
	if err != nil {
		return nil, err
	}

	return movies, nil
}
func (s *clientMovieService[T]) SearchAllClientsMovies(ctx context.Context, query *mediatypes.QueryOptions) ([]*models.MediaItem[*mediatypes.Movie], error) {
	clients, err := s.getMovieProviders(ctx, 0)
	if err != nil {
		return nil, err
	}

	var allMovies []*models.MediaItem[*mediatypes.Movie]

	for _, client := range clients {
		movieProvider, ok := client.(providers.MovieProvider)
		if !ok {
			continue
		}

		movies, err := movieProvider.GetMovies(ctx, query)
		if err != nil {
			continue
		}

		allMovies = append(allMovies, movies...)
	}

	return allMovies, nil
}

// getMovieClients gets all movie clients for a user
func (s *clientMovieService[T]) getMovieProviders(ctx context.Context, userID uint64) ([]providers.MovieProvider, error) {
	// Get all media clients for the user
	clients, err := s.clientRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var movieProviders []providers.MovieProvider

	// Filter and instantiate clients that support movies
	for _, clientConfig := range clients {
		if clientConfig.Config.Data.SupportsMovies() {
			clientID := clientConfig.GetID()
			provider, err := s.clientFactory.GetMovieProvider(ctx, clientID, clientConfig.Config.Data)
			if err != nil {
				// Log error but continue with other clients
				continue
			}
			movieProviders = append(movieProviders, provider)
		}
	}

	return movieProviders, nil
}
