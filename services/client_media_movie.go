package services

import (
	"context"
	"errors"
	"sort"

	"suasor/client"
	"suasor/client/media"
	"suasor/client/media/providers"
	mediatypes "suasor/client/media/types"
	"suasor/client/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/utils"
)

var ErrUnsupportedFeature = errors.New("feature not supported by this media client")

// ClientMovieService defines operations for interacting with movie clients
type ClientMovieService[T types.ClientConfig] interface {
	GetMovieByID(ctx context.Context, userID uint64, clientID uint64, movieID string) (*models.MediaItem[*mediatypes.Movie], error)
	GetMoviesByGenre(ctx context.Context, userID uint64, genre string) ([]*models.MediaItem[*mediatypes.Movie], error)
	GetMoviesByYear(ctx context.Context, userID uint64, year int) ([]*models.MediaItem[*mediatypes.Movie], error)
	GetMoviesByActor(ctx context.Context, userID uint64, actor string) ([]*models.MediaItem[*mediatypes.Movie], error)
	GetMoviesByDirector(ctx context.Context, userID uint64, director string) ([]*models.MediaItem[*mediatypes.Movie], error)
	GetMoviesByRating(ctx context.Context, userID uint64, minRating, maxRating float64) ([]*models.MediaItem[*mediatypes.Movie], error)
	GetLatestMoviesByAdded(ctx context.Context, userID uint64, count int) ([]*models.MediaItem[*mediatypes.Movie], error)
	GetPopularMovies(ctx context.Context, userID uint64, count int) ([]*models.MediaItem[*mediatypes.Movie], error)
	GetTopRatedMovies(ctx context.Context, userID uint64, count int) ([]*models.MediaItem[*mediatypes.Movie], error)
	SearchMovies(ctx context.Context, userID uint64, query string) ([]*models.MediaItem[*mediatypes.Movie], error)
}

type mediaMovieService[T types.ClientMediaConfig] struct {
	repo    repository.ClientRepository[T]
	factory *client.ClientFactoryService
}

// NewClientMovieService creates a new media movie service
func NewClientMovieService[T types.ClientMediaConfig](
	repo repository.ClientRepository[T],
	factory *client.ClientFactoryService,
) ClientMovieService[T] {
	return &mediaMovieService[T]{
		repo:    repo,
		factory: factory,
	}
}

// getMovieClients gets all movie clients for a user
func (s *mediaMovieService[T]) getMovieClients(ctx context.Context, userID uint64) ([]media.ClientMedia, error) {
	repo := s.repo
	// Get all media clients for the user
	clients, err := repo.GetByCategory(ctx, types.ClientCategoryMedia, userID)
	if err != nil {
		return nil, err
	}

	var movieClients []media.ClientMedia

	// Filter and instantiate clients that support movies
	for _, clientConfig := range clients {
		if clientConfig.Config.Data.SupportsMovies() {
			clientId := clientConfig.GetID()
			client, err := s.factory.GetClient(ctx, clientId, clientConfig.Config.Data)
			if err != nil {
				// Log error but continue with other clients
				continue
			}
			movieClients = append(movieClients, client.(media.ClientMedia))
		}
	}

	return movieClients, nil
}

// getSpecificMovieClient gets a specific movie client
func (s *mediaMovieService[T]) getSpecificMovieClient(ctx context.Context, userID, clientID uint64) (media.ClientMedia, error) {
	log := utils.LoggerFromContext(ctx)

	// TODO: Should see if the factory has already loaded the client. If not, then load it
	clientConfig, err := (s.repo).GetByID(ctx, clientID)
	if err != nil {
		return nil, err
	}
	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientType", clientConfig.Config.Data.GetType().String()).
		Msg("Retrieved client config")

	if !clientConfig.Config.Data.SupportsMovies() {
		log.Warn().
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("clientType", clientConfig.Config.Data.GetType().String()).
			Msg("Client does not support movies")
		return nil, ErrUnsupportedFeature
	}

	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientType", clientConfig.Config.Data.GetType().String()).
		Msg("Client supports movies")

	client, err := s.factory.GetClient(ctx, clientID, clientConfig.Config.Data)
	if err != nil {
		return nil, err
	}
	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientType", clientConfig.Config.Data.GetType().String()).
		Msg("Retrieved client")
	return client.(media.ClientMedia), nil
}

func (s *mediaMovieService[T]) GetMovieByID(ctx context.Context, userID uint64, clientID uint64, movieID string) (*models.MediaItem[*mediatypes.Movie], error) {
	client, err := s.getSpecificMovieClient(ctx, userID, clientID)
	log := utils.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Msg("Retrieved client")
	if err != nil {
		return nil, err
	}
	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("movieID", movieID).
		Msg("Retrieving movie")

	movieProvider, ok := client.(providers.MovieProvider)
	if !ok {
		log.Warn().
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("movieID", movieID).
			Msg("Client does not support movies")
		return nil, ErrUnsupportedFeature
	}

	movie, err := movieProvider.GetMovieByID(ctx, movieID)
	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("movieID", movieID).
		Msg("Retrieved movie")
	if err != nil {
		return nil, err
	}
	return movie, nil
}

func (s *mediaMovieService[T]) GetMoviesByGenre(ctx context.Context, userID uint64, genre string) ([]*models.MediaItem[*mediatypes.Movie], error) {
	clients, err := s.getMovieClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allMovies []*models.MediaItem[*mediatypes.Movie]

	// Query each client and aggregate results
	for _, client := range clients {
		movieProvider, ok := client.(providers.MovieProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Genre: genre,
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

func (s *mediaMovieService[T]) GetMoviesByYear(ctx context.Context, userID uint64, year int) ([]*models.MediaItem[*mediatypes.Movie], error) {
	clients, err := s.getMovieClients(ctx, userID)
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
			Year: year,
		}

		movies, err := movieProvider.GetMovies(ctx, options)
		if err != nil {
			continue
		}

		allMovies = append(allMovies, movies...)
	}

	return allMovies, nil
}

func (s *mediaMovieService[T]) GetMoviesByActor(ctx context.Context, userID uint64, actor string) ([]*models.MediaItem[*mediatypes.Movie], error) {
	clients, err := s.getMovieClients(ctx, userID)
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

func (s *mediaMovieService[T]) GetMoviesByDirector(ctx context.Context, userID uint64, director string) ([]*models.MediaItem[*mediatypes.Movie], error) {
	clients, err := s.getMovieClients(ctx, userID)
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

func (s *mediaMovieService[T]) GetMoviesByRating(ctx context.Context, userID uint64, minRating, maxRating float64) ([]*models.MediaItem[*mediatypes.Movie], error) {
	clients, err := s.getMovieClients(ctx, userID)
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

func (s *mediaMovieService[T]) GetLatestMoviesByAdded(ctx context.Context, userID uint64, count int) ([]*models.MediaItem[*mediatypes.Movie], error) {
	clients, err := s.getMovieClients(ctx, userID)
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

func (s *mediaMovieService[T]) GetPopularMovies(ctx context.Context, userID uint64, count int) ([]*models.MediaItem[*mediatypes.Movie], error) {
	clients, err := s.getMovieClients(ctx, userID)
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

func (s *mediaMovieService[T]) GetTopRatedMovies(ctx context.Context, userID uint64, count int) ([]*models.MediaItem[*mediatypes.Movie], error) {
	clients, err := s.getMovieClients(ctx, userID)
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

func (s *mediaMovieService[T]) SearchMovies(ctx context.Context, userID uint64, query string) ([]*models.MediaItem[*mediatypes.Movie], error) {
	clients, err := s.getMovieClients(ctx, userID)
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
