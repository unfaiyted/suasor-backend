package metadata

import (
	"context"
	"suasor/client"
	"suasor/client/types"
)

// MetadataClient interface defines the operations that a metadata provider must support
type MetadataClient interface {
	client.Client

	// Movie metadata methods
	SupportsMovieMetadata() bool
	GetMovie(ctx context.Context, id string) (*Movie, error)
	SearchMovies(ctx context.Context, query string) ([]*Movie, error)
	GetMovieRecommendations(ctx context.Context, movieID string) ([]*Movie, error)
	GetPopularMovies(ctx context.Context) ([]*Movie, error)
	GetTrendingMovies(ctx context.Context) ([]*Movie, error)

	// TV Series metadata methods
	SupportsTVMetadata() bool
	GetTVShow(ctx context.Context, id string) (*TVShow, error)
	SearchTVShows(ctx context.Context, query string) ([]*TVShow, error)
	GetTVShowRecommendations(ctx context.Context, tvShowID string) ([]*TVShow, error)
	GetPopularTVShows(ctx context.Context) ([]*TVShow, error)
	GetTrendingTVShows(ctx context.Context) ([]*TVShow, error)
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
}

// BaseMetadataClient provides a base implementation of the MetadataClient interface
type BaseMetadataClient struct {
	client.BaseClient
	Config types.ClientConfig
}

// NewBaseMetadataClient creates a new BaseMetadataClient
func NewBaseMetadataClient(config types.ClientConfig) *BaseMetadataClient {
	return &BaseMetadataClient{
		BaseClient: *client.NewBaseClient(),
		Config:     config,
	}
}

// GetConfig returns the client configuration
func (c *BaseMetadataClient) GetConfig() types.ClientConfig {
	return c.Config
}

// Default implementations that return false or empty results
func (c *BaseMetadataClient) SupportsMovieMetadata() bool {
	return false
}

func (c *BaseMetadataClient) GetMovie(ctx context.Context, id string) (*Movie, error) {
	return nil, client.ErrNotImplemented
}

func (c *BaseMetadataClient) SearchMovies(ctx context.Context, query string) ([]*Movie, error) {
	return nil, client.ErrNotImplemented
}

func (c *BaseMetadataClient) GetMovieRecommendations(ctx context.Context, movieID string) ([]*Movie, error) {
	return nil, client.ErrNotImplemented
}

func (c *BaseMetadataClient) GetPopularMovies(ctx context.Context) ([]*Movie, error) {
	return nil, client.ErrNotImplemented
}

func (c *BaseMetadataClient) GetTrendingMovies(ctx context.Context) ([]*Movie, error) {
	return nil, client.ErrNotImplemented
}

func (c *BaseMetadataClient) SupportsTVMetadata() bool {
	return false
}

func (c *BaseMetadataClient) GetTVShow(ctx context.Context, id string) (*TVShow, error) {
	return nil, client.ErrNotImplemented
}

func (c *BaseMetadataClient) SearchTVShows(ctx context.Context, query string) ([]*TVShow, error) {
	return nil, client.ErrNotImplemented
}

func (c *BaseMetadataClient) GetTVShowRecommendations(ctx context.Context, tvShowID string) ([]*TVShow, error) {
	return nil, client.ErrNotImplemented
}

func (c *BaseMetadataClient) GetPopularTVShows(ctx context.Context) ([]*TVShow, error) {
	return nil, client.ErrNotImplemented
}

func (c *BaseMetadataClient) GetTrendingTVShows(ctx context.Context) ([]*TVShow, error) {
	return nil, client.ErrNotImplemented
}

func (c *BaseMetadataClient) GetTVSeason(ctx context.Context, tvShowID string, seasonNumber int) (*TVSeason, error) {
	return nil, client.ErrNotImplemented
}

func (c *BaseMetadataClient) GetTVEpisode(ctx context.Context, tvShowID string, seasonNumber int, episodeNumber int) (*TVEpisode, error) {
	return nil, client.ErrNotImplemented
}

func (c *BaseMetadataClient) SupportsPersonMetadata() bool {
	return false
}

func (c *BaseMetadataClient) GetPerson(ctx context.Context, id string) (*Person, error) {
	return nil, client.ErrNotImplemented
}

func (c *BaseMetadataClient) SearchPeople(ctx context.Context, query string) ([]*Person, error) {
	return nil, client.ErrNotImplemented
}

func (c *BaseMetadataClient) GetPersonMovieCredits(ctx context.Context, personID string) ([]*MovieCredit, error) {
	return nil, client.ErrNotImplemented
}

func (c *BaseMetadataClient) GetPersonTVCredits(ctx context.Context, personID string) ([]*TVCredit, error) {
	return nil, client.ErrNotImplemented
}

func (c *BaseMetadataClient) SupportsCollectionMetadata() bool {
	return false
}

func (c *BaseMetadataClient) GetCollection(ctx context.Context, id string) (*Collection, error) {
	return nil, client.ErrNotImplemented
}

func (c *BaseMetadataClient) SearchCollections(ctx context.Context, query string) ([]*Collection, error) {
	return nil, client.ErrNotImplemented
}