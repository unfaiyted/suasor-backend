package metadata

import (
	"context"
	"suasor/clients"
	"suasor/clients/types"
)

// clientMetadata interface defines the operations that a metadata provider must support
type ClientMetadata interface {
	clients.Client

	// Movie metadata methods
	SupportsMovieMetadata() bool
	GetMovie(ctx context.Context, id string) (*Movie, error)
	SearchMovies(ctx context.Context, query string) ([]*Movie, error)
	GetMovieRecommendations(ctx context.Context, movieID string) ([]*Movie, error)
	GetPopularMovies(ctx context.Context) ([]*Movie, error)
	GetTrendingMovies(ctx context.Context) ([]*Movie, error)
	GetUpcomingMovies(ctx context.Context, daysAhead int) ([]*Movie, error)
	GetNowPlayingMovies(ctx context.Context, daysPast int) ([]*Movie, error)

	// TV Series metadata methods
	SupportsTVMetadata() bool
	GetTVShow(ctx context.Context, id string) (*TVShow, error)
	SearchTVShows(ctx context.Context, query string) ([]*TVShow, error)
	GetTVShowRecommendations(ctx context.Context, tvShowID string) ([]*TVShow, error)
	GetPopularTVShows(ctx context.Context) ([]*TVShow, error)
	GetTrendingTVShows(ctx context.Context) ([]*TVShow, error)
	GetRecentTVShows(ctx context.Context, daysWindow int) ([]*TVShow, error)
	GetTVSeason(ctx context.Context, tvShowID string, seasonNumber int) (*TVSeason, error)
	GetTVEpisode(ctx context.Context, tvShowID string, seasonNumber int, episodeNumber int) (*TVEpisode, error)

	// Person metadata methods
	SupportsPersonMetadata() bool
	GetPerson(ctx context.Context, id string) (*Person, error)
	SearchPeople(ctx context.Context, query string) ([]*Person, error)
	GetPersonMovieCredits(ctx context.Context, personID string) ([]*MovieCredit, error)
	GetPersonTVCredits(ctx context.Context, personID string) ([]*TVCredit, error)

	// Collection metadata methods
	SupportsCollectionMetadata() bool
	GetCollection(ctx context.Context, id string) (*Collection, error)
	SearchCollections(ctx context.Context, query string) ([]*Collection, error)

	GetMetadataConfig() types.ClientMetadataConfig
}

// clientMetadata provides a base implementation of the clientMetadata interface
type clientMetadata struct {
	clients.Client
	config *types.ClientMetadataConfig
}

// NewclientMetadata creates a new clientMetadata
func NewClientMetadata(ctx, clientID uint64, config types.ClientMetadataConfig) (ClientMetadata, error) {
	// Create a new client with the provided config
	client := clients.NewClient(clientID, config.GetCategory(), config)
	return &clientMetadata{
		Client: client,
		config: &config,
	}, nil
}

// Default implementations that return false or empty results
func (c *clientMetadata) SupportsMovieMetadata() bool {
	return false
}

func (c *clientMetadata) GetMovie(ctx context.Context, id string) (*Movie, error) {
	return nil, clients.ErrNotImplemented
}

func (c *clientMetadata) SearchMovies(ctx context.Context, query string) ([]*Movie, error) {
	return nil, clients.ErrNotImplemented
}

func (c *clientMetadata) GetMovieRecommendations(ctx context.Context, movieID string) ([]*Movie, error) {
	return nil, clients.ErrNotImplemented
}

func (c *clientMetadata) GetPopularMovies(ctx context.Context) ([]*Movie, error) {
	return nil, clients.ErrNotImplemented
}

func (c *clientMetadata) GetTrendingMovies(ctx context.Context) ([]*Movie, error) {
	return nil, clients.ErrNotImplemented
}

func (c *clientMetadata) GetUpcomingMovies(ctx context.Context, daysAhead int) ([]*Movie, error) {
	return nil, clients.ErrNotImplemented
}

func (c *clientMetadata) GetNowPlayingMovies(ctx context.Context, daysPast int) ([]*Movie, error) {
	return nil, clients.ErrNotImplemented
}

func (c *clientMetadata) SupportsTVMetadata() bool {
	return false
}

func (c *clientMetadata) GetTVShow(ctx context.Context, id string) (*TVShow, error) {
	return nil, clients.ErrNotImplemented
}

func (c *clientMetadata) SearchTVShows(ctx context.Context, query string) ([]*TVShow, error) {
	return nil, clients.ErrNotImplemented
}

func (c *clientMetadata) GetTVShowRecommendations(ctx context.Context, tvShowID string) ([]*TVShow, error) {
	return nil, clients.ErrNotImplemented
}

func (c *clientMetadata) GetPopularTVShows(ctx context.Context) ([]*TVShow, error) {
	return nil, clients.ErrNotImplemented
}

func (c *clientMetadata) GetTrendingTVShows(ctx context.Context) ([]*TVShow, error) {
	return nil, clients.ErrNotImplemented
}

func (c *clientMetadata) GetRecentTVShows(ctx context.Context, daysWindow int) ([]*TVShow, error) {
	return nil, clients.ErrNotImplemented
}

func (c *clientMetadata) GetTVSeason(ctx context.Context, tvShowID string, seasonNumber int) (*TVSeason, error) {
	return nil, clients.ErrNotImplemented
}

func (c *clientMetadata) GetTVEpisode(ctx context.Context, tvShowID string, seasonNumber int, episodeNumber int) (*TVEpisode, error) {
	return nil, clients.ErrNotImplemented
}

func (c *clientMetadata) SupportsPersonMetadata() bool {
	return false
}

func (c *clientMetadata) GetPerson(ctx context.Context, id string) (*Person, error) {
	return nil, clients.ErrNotImplemented
}

func (c *clientMetadata) SearchPeople(ctx context.Context, query string) ([]*Person, error) {
	return nil, clients.ErrNotImplemented
}

func (c *clientMetadata) GetPersonMovieCredits(ctx context.Context, personID string) ([]*MovieCredit, error) {
	return nil, clients.ErrNotImplemented
}

func (c *clientMetadata) GetPersonTVCredits(ctx context.Context, personID string) ([]*TVCredit, error) {
	return nil, clients.ErrNotImplemented
}

func (c *clientMetadata) SupportsCollectionMetadata() bool {
	return false
}

func (c *clientMetadata) GetCollection(ctx context.Context, id string) (*Collection, error) {
	return nil, clients.ErrNotImplemented
}

func (c *clientMetadata) SearchCollections(ctx context.Context, query string) ([]*Collection, error) {
	return nil, clients.ErrNotImplemented
}

func (c *clientMetadata) GetMetadataConfig() types.ClientMetadataConfig {
	return *c.config
}
