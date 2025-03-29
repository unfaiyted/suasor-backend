package services

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/joho/godotenv"
	"suasor/client/media"
	mediatypes "suasor/client/media/types"
	"suasor/client/types"
	"suasor/types/models"
)

func init() {
	// Try to load .env file from several possible locations
	locations := []string{
		".env",          // Current directory
		"../../../.env", // Project root
		filepath.Join(os.Getenv("HOME"), "media_movie_test.env"), // Home directory
	}

	for _, location := range locations {
		err := godotenv.Load(location)
		if err == nil {
			break
		}
	}
}

// MockMediaClientConfig implements MediaClientConfig for testing
type MockMediaClientConfig struct {
	SupportsMoviesVal bool
}

func (m MockMediaClientConfig) GetCategory() types.ClientCategory {
	return types.ClientCategory("media")
}

func (m MockMediaClientConfig) SupportsMovies() bool {
	return m.SupportsMoviesVal
}

func (m MockMediaClientConfig) SupportsTVShows() bool {
	return false
}

func (m MockMediaClientConfig) SupportsMusic() bool {
	return false
}

func (m MockMediaClientConfig) GetConfigType() string {
	return "mock"
}

// Add the missing GetClientType method
func (m MockMediaClientConfig) GetClientType() string {
	return "mock"
}

// MockMediaClient implements both media.MediaClient and providers.MovieProvider
type MockMediaClient struct {
	mock.Mock
	SupportsHistoryVal bool
	clientID           uint64
}

func (m *MockMediaClient) GetID() uint64 {
	return m.clientID
}

func (m *MockMediaClient) GetName() string {
	return "MockClient"
}

func (m *MockMediaClient) SupportsCollections() bool {
	return false // Change to true if needed for your tests
}

func (m *MockMediaClient) SupportsHistory() bool {
	return m.SupportsHistoryVal
}

// Add this method to implement the MediaClient interface fully
func (m *MockMediaClient) SupportsMovies() bool {
	return true
}

// MovieProvider methods
func (m *MockMediaClient) GetMovies(ctx context.Context, options *mediatypes.QueryOptions) ([]models.MediaItem[mediatypes.Movie], error) {
	args := m.Called(ctx, options)
	return args.Get(0).([]models.MediaItem[mediatypes.Movie]), args.Error(1)
}

func (m *MockMediaClient) GetMovieByID(ctx context.Context, id string) (models.MediaItem[mediatypes.Movie], error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.MediaItem[mediatypes.Movie]), args.Error(1)
}

func (m *MockMediaClient) GetMovieGenres(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

// MockClientRepository implements repository.ClientRepository
type MockClientRepository struct {
	clients map[uint64]*models.Client[types.MediaClientConfig]
}

func NewMockClientRepository() *MockClientRepository {
	return &MockClientRepository{
		clients: make(map[uint64]*models.Client[types.MediaClientConfig]),
	}
}

func (m *MockClientRepository) AddClient(client *models.Client[types.MediaClientConfig]) {
	m.clients[client.ID] = client
}

func (m *MockClientRepository) GetByID(ctx context.Context, id uint64, userID uint64) (*models.Client[types.MediaClientConfig], error) {
	client, exists := m.clients[id]
	if !exists || client.UserID != userID {
		return nil, assert.AnError
	}
	return client, nil
}

func (m *MockClientRepository) GetByCategory(ctx context.Context, category string, userID uint64) ([]models.Client[types.MediaClientConfig], error) {
	var results []models.Client[types.MediaClientConfig]
	for _, client := range m.clients {
		if client.UserID == userID {
			results = append(results, *client)
		}
	}
	return results, nil
}

// Fix the Create method signature
func (m *MockClientRepository) Create(ctx context.Context, client models.Client[types.MediaClientConfig]) (*models.Client[types.MediaClientConfig], error) {
	clientCopy := client
	m.clients[client.ID] = &clientCopy
	return &clientCopy, nil
}

// MockClientFactory implements media.ClientFactory
type MockClientFactory struct {
	clients map[uint64]media.MediaClient
}

func NewMockClientFactory() *MockClientFactory {
	return &MockClientFactory{
		clients: make(map[uint64]media.MediaClient),
	}
}

func (m *MockClientFactory) AddClient(clientID uint64, client media.MediaClient) {
	m.clients[clientID] = client
}

func (m *MockClientFactory) CreateMediaClient(ctx context.Context, clientID uint64, config types.MediaClientConfig) (media.MediaClient, error) {
	client, exists := m.clients[clientID]
	if !exists {
		return nil, assert.AnError
	}
	return client, nil
}

func TestMediaMovieServiceIntegration(t *testing.T) {
	// Skip if not running integration tests
	if os.Getenv("INTEGRATION") != "true" {
		t.Skip("Skipping integration test. Set INTEGRATION=true to run")
	}

	// Create context
	ctx := context.Background()

	// Set up test data
	userID := uint64(1)
	clientID := uint64(1)

	// Create a test movie
	testMovie := models.MediaItem[mediatypes.Movie]{
		ExternalID: "movie123",
		ClientID:   clientID,
		Data: mediatypes.Movie{
			Cast: []mediatypes.Person{},
			Crew: []mediatypes.Person{},
			Details: mediatypes.MediaDetails{
				Title:       "Test Movie",
				ReleaseYear: 2020,
				Artwork: mediatypes.Artwork{
					Poster: "http://example.com/poster.jpg",
				},
				AddedAt:     time.Now(),
				Description: "A test movie plot",
				Genres:      []string{"Action", "Drama"},
				Ratings:     mediatypes.Ratings{mediatypes.Rating{Source: "tmdb", Value: 8.5, Votes: 10}},
			},
		},
	}

	// Create test movies list
	testMovies := []models.MediaItem[mediatypes.Movie]{testMovie}

	// Create mock client, repository, and factory FIRST (before using them)
	mockClient := &MockMediaClient{clientID: clientID}
	mockClient.SupportsHistoryVal = true

	mockRepo := NewMockClientRepository()
	mockFactory := NewMockClientFactory()

	// Now we can use the mocks
	mockFactory.AddClient(clientID, mockClient)

	// Create the service (only once)
	service := NewMediaMovieService(mockRepo, mockFactory)

	// Set up expectations
	mockClient.On("GetMovieByID", mock.Anything, "movie123").Return(testMovie, nil)
	mockClient.On("GetMovies", mock.Anything, mock.Anything).Return(testMovies, nil)
	mockClient.On("GetMovieGenres", mock.Anything).Return([]string{"Action", "Drama", "Comedy"}, nil)

	// Create client repository and add test client
	mockRepo.AddClient(&models.Client[types.MediaClientConfig]{
		BaseModel: models.BaseModel{
			ID: clientID,
		},
		UserID: userID,
		Config: models.ClientConfigWrapper[types.MediaClientConfig]{
			Data: MockMediaClientConfig{
				SupportsMoviesVal: true,
			},
		},
	})

	// Run all test cases
	t.Run("TestGetMovieByID", func(t *testing.T) {
		movie, err := service.GetMovieByID(ctx, userID, clientID, "movie123")
		require.NoError(t, err)
		assert.Equal(t, "Test Movie", movie.Data.Details.Title)
		assert.Equal(t, 2020, movie.Data.Details.ReleaseYear)
	})

	t.Run("TestGetMoviesByGenre", func(t *testing.T) {
		movies, err := service.GetMoviesByGenre(ctx, userID, "Action")
		require.NoError(t, err)
		assert.Len(t, movies, 1)
		assert.Equal(t, "Test Movie", movies[0].Data.Details.Title)
	})

	t.Run("TestGetMoviesByYear", func(t *testing.T) {
		movies, err := service.GetMoviesByYear(ctx, userID, 2020)
		require.NoError(t, err)
		assert.Len(t, movies, 1)
		assert.Equal(t, "Test Movie", movies[0].Data.Details.Title)
	})

	t.Run("TestGetMoviesByActor", func(t *testing.T) {
		movies, err := service.GetMoviesByActor(ctx, userID, "Actor One")
		require.NoError(t, err)
		assert.Len(t, movies, 1)
		assert.Equal(t, "Test Movie", movies[0].Data.Details.Title)
	})

	t.Run("TestGetMoviesByDirector", func(t *testing.T) {
		movies, err := service.GetMoviesByDirector(ctx, userID, "Director One")
		require.NoError(t, err)
		assert.Len(t, movies, 1)
		assert.Equal(t, "Test Movie", movies[0].Data.Details.Title)
	})

	t.Run("TestGetMoviesByRating", func(t *testing.T) {
		movies, err := service.GetMoviesByRating(ctx, userID, 8.0, 9.0)
		require.NoError(t, err)
		assert.Len(t, movies, 1)
		assert.Equal(t, "Test Movie", movies[0].Data.Details.Title)
	})

	t.Run("TestGetLatestMoviesByAdded", func(t *testing.T) {
		movies, err := service.GetLatestMoviesByAdded(ctx, userID, 10)
		require.NoError(t, err)
		assert.Len(t, movies, 1)
		assert.Equal(t, "Test Movie", movies[0].Data.Details.Title)
	})

	t.Run("TestGetPopularMovies", func(t *testing.T) {
		movies, err := service.GetPopularMovies(ctx, userID, 10)
		require.NoError(t, err)
		assert.Len(t, movies, 1)
		assert.Equal(t, "Test Movie", movies[0].Data.Details.Title)
	})

	t.Run("TestGetTopRatedMovies", func(t *testing.T) {
		movies, err := service.GetTopRatedMovies(ctx, userID, 10)
		require.NoError(t, err)
		assert.Len(t, movies, 1)
		assert.Equal(t, "Test Movie", movies[0].Data.Details.Title)
	})

	t.Run("TestSearchMovies", func(t *testing.T) {
		movies, err := service.SearchMovies(ctx, userID, "Test")
		require.NoError(t, err)
		assert.Len(t, movies, 1)
		assert.Equal(t, "Test Movie", movies[0].Data.Details.Title)
	})
}
