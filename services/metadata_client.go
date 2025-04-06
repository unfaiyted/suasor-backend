package services

import (
	"context"
	"fmt"
	"suasor/client"
	"suasor/client/metadata"
	"suasor/client/types"
	"suasor/repository"
	"suasor/types/responses"
)

// MetadataClientService provides operations for metadata clients
type MetadataClientService[T types.MetadataClientConfig] struct {
	clientService *clientService[T]
}

// NewMetadataClientService creates a new MetadataClientService
func NewMetadataClientService[T types.MetadataClientConfig](factory *client.ClientFactoryService, repo repository.ClientRepository[T]) *MetadataClientService[T] {
	return &MetadataClientService[T]{
		clientService: NewClientService(factory, repo),
	}
}

// GetMetadataClient returns a metadata client instance for the given client ID
func (s *MetadataClientService[T]) GetMetadataClient(ctx context.Context, clientID uint64) (metadata.MetadataClient, error) {
	// Get the client configuration
	clientModel, err := s.clientService.GetByID(ctx, clientID, 0) // 0 for userID as metadata clients may be system-wide
	if err != nil {
		return nil, err
	}

	// Create client instance using the factory
	clientInstance, err := s.clientService.factory.GetClient(ctx, clientID, clientModel.Config.Data)
	if err != nil {
		return nil, err
	}

	// Check if it's a metadata client
	metadataClient, ok := clientInstance.(metadata.MetadataClient)
	if !ok {
		return nil, fmt.Errorf("client with ID %d is not a metadata client", clientID)
	}

	return metadataClient, nil
}

// Movies

// GetMovie retrieves a movie by ID
func (s *MetadataClientService[T]) GetMovie(ctx context.Context, clientID uint64, movieID string) (*responses.MetadataMovieResponse, error) {
	client, err := s.GetMetadataClient(ctx, clientID)
	if err != nil {
		return nil, err
	}

	if !client.SupportsMovieMetadata() {
		return nil, fmt.Errorf("client does not support movie metadata")
	}

	movie, err := client.GetMovie(ctx, movieID)
	if err != nil {
		return nil, err
	}

	return &responses.MetadataMovieResponse{
		Data:    movie,
		Success: true,
	}, nil
}

// SearchMovies searches for movies by query
func (s *MetadataClientService[T]) SearchMovies(ctx context.Context, clientID uint64, query string) (*responses.MetadataMovieSearchResponse, error) {
	client, err := s.GetMetadataClient(ctx, clientID)
	if err != nil {
		return nil, err
	}

	if !client.SupportsMovieMetadata() {
		return nil, fmt.Errorf("client does not support movie metadata")
	}

	movies, err := client.SearchMovies(ctx, query)
	if err != nil {
		return nil, err
	}

	return &responses.MetadataMovieSearchResponse{
		Data:    movies,
		Success: true,
	}, nil
}

// GetMovieRecommendations gets movie recommendations based on a movie ID
func (s *MetadataClientService[T]) GetMovieRecommendations(ctx context.Context, clientID uint64, movieID string) (*responses.MetadataMovieSearchResponse, error) {
	client, err := s.GetMetadataClient(ctx, clientID)
	if err != nil {
		return nil, err
	}

	if !client.SupportsMovieMetadata() {
		return nil, fmt.Errorf("client does not support movie metadata")
	}

	movies, err := client.GetMovieRecommendations(ctx, movieID)
	if err != nil {
		return nil, err
	}

	return &responses.MetadataMovieSearchResponse{
		Data:    movies,
		Success: true,
	}, nil
}

// GetPopularMovies gets popular movies
func (s *MetadataClientService[T]) GetPopularMovies(ctx context.Context, clientID uint64) (*responses.MetadataMovieSearchResponse, error) {
	client, err := s.GetMetadataClient(ctx, clientID)
	if err != nil {
		return nil, err
	}

	if !client.SupportsMovieMetadata() {
		return nil, fmt.Errorf("client does not support movie metadata")
	}

	movies, err := client.GetPopularMovies(ctx)
	if err != nil {
		return nil, err
	}

	return &responses.MetadataMovieSearchResponse{
		Data:    movies,
		Success: true,
	}, nil
}

// GetTrendingMovies gets trending movies
func (s *MetadataClientService[T]) GetTrendingMovies(ctx context.Context, clientID uint64) (*responses.MetadataMovieSearchResponse, error) {
	client, err := s.GetMetadataClient(ctx, clientID)
	if err != nil {
		return nil, err
	}

	if !client.SupportsMovieMetadata() {
		return nil, fmt.Errorf("client does not support movie metadata")
	}

	movies, err := client.GetTrendingMovies(ctx)
	if err != nil {
		return nil, err
	}

	return &responses.MetadataMovieSearchResponse{
		Data:    movies,
		Success: true,
	}, nil
}

// TV Shows

// GetTVShow retrieves a TV show by ID
func (s *MetadataClientService[T]) GetTVShow(ctx context.Context, clientID uint64, tvShowID string) (*responses.MetadataTVShowResponse, error) {
	client, err := s.GetMetadataClient(ctx, clientID)
	if err != nil {
		return nil, err
	}

	if !client.SupportsTVMetadata() {
		return nil, fmt.Errorf("client does not support TV metadata")
	}

	tvShow, err := client.GetTVShow(ctx, tvShowID)
	if err != nil {
		return nil, err
	}

	return &responses.MetadataTVShowResponse{
		Data:    tvShow,
		Success: true,
	}, nil
}

// SearchTVShows searches for TV shows by query
func (s *MetadataClientService[T]) SearchTVShows(ctx context.Context, clientID uint64, query string) (*responses.MetadataTVShowSearchResponse, error) {
	client, err := s.GetMetadataClient(ctx, clientID)
	if err != nil {
		return nil, err
	}

	if !client.SupportsTVMetadata() {
		return nil, fmt.Errorf("client does not support TV metadata")
	}

	tvShows, err := client.SearchTVShows(ctx, query)
	if err != nil {
		return nil, err
	}

	return &responses.MetadataTVShowSearchResponse{
		Data:    tvShows,
		Success: true,
	}, nil
}

// GetTVShowRecommendations gets TV show recommendations based on a TV show ID
func (s *MetadataClientService[T]) GetTVShowRecommendations(ctx context.Context, clientID uint64, tvShowID string) (*responses.MetadataTVShowSearchResponse, error) {
	client, err := s.GetMetadataClient(ctx, clientID)
	if err != nil {
		return nil, err
	}

	if !client.SupportsTVMetadata() {
		return nil, fmt.Errorf("client does not support TV metadata")
	}

	tvShows, err := client.GetTVShowRecommendations(ctx, tvShowID)
	if err != nil {
		return nil, err
	}

	return &responses.MetadataTVShowSearchResponse{
		Data:    tvShows,
		Success: true,
	}, nil
}

// GetPopularTVShows gets popular TV shows
func (s *MetadataClientService[T]) GetPopularTVShows(ctx context.Context, clientID uint64) (*responses.MetadataTVShowSearchResponse, error) {
	client, err := s.GetMetadataClient(ctx, clientID)
	if err != nil {
		return nil, err
	}

	if !client.SupportsTVMetadata() {
		return nil, fmt.Errorf("client does not support TV metadata")
	}

	tvShows, err := client.GetPopularTVShows(ctx)
	if err != nil {
		return nil, err
	}

	return &responses.MetadataTVShowSearchResponse{
		Data:    tvShows,
		Success: true,
	}, nil
}

// GetTrendingTVShows gets trending TV shows
func (s *MetadataClientService[T]) GetTrendingTVShows(ctx context.Context, clientID uint64) (*responses.MetadataTVShowSearchResponse, error) {
	client, err := s.GetMetadataClient(ctx, clientID)
	if err != nil {
		return nil, err
	}

	if !client.SupportsTVMetadata() {
		return nil, fmt.Errorf("client does not support TV metadata")
	}

	tvShows, err := client.GetTrendingTVShows(ctx)
	if err != nil {
		return nil, err
	}

	return &responses.MetadataTVShowSearchResponse{
		Data:    tvShows,
		Success: true,
	}, nil
}

// GetTVSeason retrieves a TV season by show ID and season number
func (s *MetadataClientService[T]) GetTVSeason(ctx context.Context, clientID uint64, tvShowID string, seasonNumber int) (*responses.MetadataTVSeasonResponse, error) {
	client, err := s.GetMetadataClient(ctx, clientID)
	if err != nil {
		return nil, err
	}

	if !client.SupportsTVMetadata() {
		return nil, fmt.Errorf("client does not support TV metadata")
	}

	season, err := client.GetTVSeason(ctx, tvShowID, seasonNumber)
	if err != nil {
		return nil, err
	}

	return &responses.MetadataTVSeasonResponse{
		Data:    season,
		Success: true,
	}, nil
}

// GetTVEpisode retrieves a TV episode by show ID, season number, and episode number
func (s *MetadataClientService[T]) GetTVEpisode(ctx context.Context, clientID uint64, tvShowID string, seasonNumber, episodeNumber int) (*responses.MetadataTVEpisodeResponse, error) {
	client, err := s.GetMetadataClient(ctx, clientID)
	if err != nil {
		return nil, err
	}

	if !client.SupportsTVMetadata() {
		return nil, fmt.Errorf("client does not support TV metadata")
	}

	episode, err := client.GetTVEpisode(ctx, tvShowID, seasonNumber, episodeNumber)
	if err != nil {
		return nil, err
	}

	return &responses.MetadataTVEpisodeResponse{
		Data:    episode,
		Success: true,
	}, nil
}

// People

// GetPerson retrieves a person by ID
func (s *MetadataClientService[T]) GetPerson(ctx context.Context, clientID uint64, personID string) (*responses.MetadataPersonResponse, error) {
	client, err := s.GetMetadataClient(ctx, clientID)
	if err != nil {
		return nil, err
	}

	if !client.SupportsPersonMetadata() {
		return nil, fmt.Errorf("client does not support person metadata")
	}

	person, err := client.GetPerson(ctx, personID)
	if err != nil {
		return nil, err
	}

	return &responses.MetadataPersonResponse{
		Data:    person,
		Success: true,
	}, nil
}

// SearchPeople searches for people by query
func (s *MetadataClientService[T]) SearchPeople(ctx context.Context, clientID uint64, query string) (*responses.MetadataPersonSearchResponse, error) {
	client, err := s.GetMetadataClient(ctx, clientID)
	if err != nil {
		return nil, err
	}

	if !client.SupportsPersonMetadata() {
		return nil, fmt.Errorf("client does not support person metadata")
	}

	people, err := client.SearchPeople(ctx, query)
	if err != nil {
		return nil, err
	}

	return &responses.MetadataPersonSearchResponse{
		Data:    people,
		Success: true,
	}, nil
}

// Collections

// GetCollection retrieves a collection by ID
func (s *MetadataClientService[T]) GetCollection(ctx context.Context, clientID uint64, collectionID string) (*responses.MetadataCollectionResponse, error) {
	client, err := s.GetMetadataClient(ctx, clientID)
	if err != nil {
		return nil, err
	}

	if !client.SupportsCollectionMetadata() {
		return nil, fmt.Errorf("client does not support collection metadata")
	}

	collection, err := client.GetCollection(ctx, collectionID)
	if err != nil {
		return nil, err
	}

	return &responses.MetadataCollectionResponse{
		Data:    collection,
		Success: true,
	}, nil
}

// SearchCollections searches for collections by query
func (s *MetadataClientService[T]) SearchCollections(ctx context.Context, clientID uint64, query string) (*responses.MetadataCollectionSearchResponse, error) {
	client, err := s.GetMetadataClient(ctx, clientID)
	if err != nil {
		return nil, err
	}

	if !client.SupportsCollectionMetadata() {
		return nil, fmt.Errorf("client does not support collection metadata")
	}

	collections, err := client.SearchCollections(ctx, query)
	if err != nil {
		return nil, err
	}

	return &responses.MetadataCollectionSearchResponse{
		Data:    collections,
		Success: true,
	}, nil
}