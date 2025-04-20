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
	"suasor/client/media/providers"
	mediatypes "suasor/client/media/types"
	"suasor/client/types"
	"suasor/types/models"
)

func init() {
	// Try to load .env file from several possible locations
	locations := []string{
		".env",        // Current directory
		"../.env",     // Project root
		"../../.env",  // Root directory
		filepath.Join(os.Getenv("HOME"), "media_playlist_test.env"), // Home directory
	}

	for _, location := range locations {
		err := godotenv.Load(location)
		if err == nil {
			break
		}
	}
}

// MockClientMediaConfig implements ClientMediaConfig for testing
type MockPlaylistClientMediaConfig struct {
	SupportsPlaylistsVal bool
}

func (m MockPlaylistClientMediaConfig) GetCategory() types.ClientCategory {
	return types.ClientCategory("media")
}

func (m MockPlaylistClientMediaConfig) SupportsPlaylists() bool {
	return m.SupportsPlaylistsVal
}

func (m MockPlaylistClientMediaConfig) SupportsMovies() bool {
	return false
}

func (m MockPlaylistClientMediaConfig) SupportsTVShows() bool {
	return false
}

func (m MockPlaylistClientMediaConfig) SupportsMusic() bool {
	return false
}

func (m MockPlaylistClientMediaConfig) GetConfigType() string {
	return "mock"
}

// Add the missing GetClientType method
func (m MockPlaylistClientMediaConfig) GetClientType() types.ClientMediaType {
	return types.ClientMediaType("mock")
}

// MockPlaylistClient implements both media.ClientMedia and providers.PlaylistProvider
type MockPlaylistClient struct {
	mock.Mock
	clientID uint64
}

func (m *MockPlaylistClient) GetID() uint64 {
	return m.clientID
}

func (m *MockPlaylistClient) GetName() string {
	return "MockClient"
}

func (m *MockPlaylistClient) SupportsCollections() bool {
	return false
}

func (m *MockPlaylistClient) SupportsHistory() bool {
	return false
}

// Add this method to implement the ClientMedia interface fully
func (m *MockPlaylistClient) SupportsMovies() bool {
	return false
}

func (m *MockPlaylistClient) SupportsPlaylists() bool {
	return true
}

// PlaylistProvider methods
func (m *MockPlaylistClient) GetPlaylists(ctx context.Context, options *mediatypes.QueryOptions) ([]models.MediaItem[mediatypes.Playlist], error) {
	args := m.Called(ctx, options)
	return args.Get(0).([]models.MediaItem[mediatypes.Playlist]), args.Error(1)
}

// MockClientRepository implements repository.ClientRepository
type MockPlaylistClientRepository struct {
	clients map[uint64]*models.Client[types.ClientMediaConfig]
}

func NewMockPlaylistClientRepository() *MockPlaylistClientRepository {
	return &MockPlaylistClientRepository{
		clients: make(map[uint64]*models.Client[types.ClientMediaConfig]),
	}
}

func (m *MockPlaylistClientRepository) AddClient(client *models.Client[types.ClientMediaConfig]) {
	m.clients[client.ID] = client
}

func (m *MockPlaylistClientRepository) GetByID(ctx context.Context, id uint64, userID uint64) (*models.Client[types.ClientMediaConfig], error) {
	client, exists := m.clients[id]
	if !exists || client.UserID != userID {
		return nil, assert.AnError
	}
	return client, nil
}

func (m *MockPlaylistClientRepository) GetByCategory(ctx context.Context, category string, userID uint64) ([]models.Client[types.ClientMediaConfig], error) {
	var results []models.Client[types.ClientMediaConfig]
	for _, client := range m.clients {
		if client.UserID == userID {
			results = append(results, *client)
		}
	}
	return results, nil
}

// Fix the Create method signature
func (m *MockPlaylistClientRepository) Create(ctx context.Context, client models.Client[types.ClientMediaConfig]) (*models.Client[types.ClientMediaConfig], error) {
	clientCopy := client
	m.clients[client.ID] = &clientCopy
	return &clientCopy, nil
}

// MockClientFactory implements client factory
type MockPlaylistClientFactoryService struct {
	clients map[uint64]media.ClientMedia
}

func NewMockPlaylistClientFactoryService() *MockPlaylistClientFactoryService {
	return &MockPlaylistClientFactoryService{
		clients: make(map[uint64]media.ClientMedia),
	}
}

func (m *MockPlaylistClientFactoryService) AddClient(clientID uint64, client media.ClientMedia) {
	m.clients[clientID] = client
}

func (m *MockPlaylistClientFactoryService) GetClient(ctx context.Context, clientID uint64, config types.ClientConfig) (interface{}, error) {
	client, exists := m.clients[clientID]
	if !exists {
		return nil, assert.AnError
	}
	return client, nil
}

func TestMediaPlaylistServiceIntegration(t *testing.T) {
	// Skip if not running integration tests
	if os.Getenv("INTEGRATION") != "true" {
		t.Skip("Skipping integration test. Set INTEGRATION=true to run")
	}

	// Create context
	ctx := context.Background()

	// Set up test data
	userID := uint64(1)
	clientID := uint64(1)

	// Create a test playlist
	testPlaylist := models.MediaItem[mediatypes.Playlist]{
		ExternalID: "playlist123",
		ClientID:   clientID,
		Data: mediatypes.Playlist{
			Details: mediatypes.MediaDetails{
				Title:       "Test Playlist",
				Description: "A test playlist",
				AddedAt:     time.Now(),
				Artwork: mediatypes.Artwork{
					Poster: "http://example.com/poster.jpg",
				},
			},
			ItemIDs:   []string{"item1", "item2", "item3"},
			ItemCount: 3,
			Owner:     "testuser",
			IsPublic:  true,
		},
	}

	// Create test playlists list
	testPlaylists := []models.MediaItem[mediatypes.Playlist]{testPlaylist}

	// Create mock client, repository, and factory
	mockClient := &MockPlaylistClient{clientID: clientID}
	mockRepo := NewMockPlaylistClientRepository()
	mockFactory := NewMockPlaylistClientFactoryService()

	// Add client to factory
	mockFactory.AddClient(clientID, mockClient)

	// Create the service
	service := NewClientMediaPlaylistService(mockRepo, mockFactory)

	// Set up expectations
	mockClient.On("GetPlaylists", mock.Anything, mock.Anything).Return(testPlaylists, nil)

	// Create client repository and add test client
	mockRepo.AddClient(&models.Client[types.ClientMediaConfig]{
		BaseModel: models.BaseModel{
			ID: clientID,
		},
		UserID: userID,
		Config: models.ClientConfigWrapper[types.ClientMediaConfig]{
			Data: MockPlaylistClientMediaConfig{
				SupportsPlaylistsVal: true,
			},
		},
	})

	// Run all test cases
	t.Run("TestGetPlaylistByID", func(t *testing.T) {
		// Setup the mock to return a specific playlist when filtering by ID
		mockClient.On("GetPlaylists", mock.Anything, mock.MatchedBy(func(options *mediatypes.QueryOptions) bool {
			return options != nil && options.ExternalSourceID == "playlist123"
		})).Return(testPlaylists, nil)

		playlist, err := service.GetPlaylistByID(ctx, userID, clientID, "playlist123")
		require.NoError(t, err)
		assert.Equal(t, "Test Playlist", playlist.Data.Details.Title)
		assert.Equal(t, 3, playlist.Data.ItemCount)
		assert.Equal(t, "testuser", playlist.Data.Owner)
	})

	t.Run("TestGetPlaylists", func(t *testing.T) {
		playlists, err := service.GetPlaylists(ctx, userID, 10)
		require.NoError(t, err)
		assert.Len(t, playlists, 1)
		assert.Equal(t, "Test Playlist", playlists[0].Data.Details.Title)
		assert.Equal(t, 3, playlists[0].Data.ItemCount)
	})

	t.Run("TestSearchPlaylists", func(t *testing.T) {
		// Setup mock for search query
		mockClient.On("GetPlaylists", mock.Anything, mock.MatchedBy(func(options *mediatypes.QueryOptions) bool {
			return options != nil && options.Query == "Test"
		})).Return(testPlaylists, nil)

		playlists, err := service.SearchPlaylists(ctx, userID, "Test")
		require.NoError(t, err)
		assert.Len(t, playlists, 1)
		assert.Equal(t, "Test Playlist", playlists[0].Data.Details.Title)
	})

	// Test the methods that are expected to return "not implemented" for now
	t.Run("TestCreatePlaylist", func(t *testing.T) {
		_, err := service.CreatePlaylist(ctx, userID, clientID, "New Playlist", "Description")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not implemented")
	})

	t.Run("TestUpdatePlaylist", func(t *testing.T) {
		_, err := service.UpdatePlaylist(ctx, userID, clientID, "playlist123", "Updated Title", "Updated Description")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not implemented")
	})

	t.Run("TestDeletePlaylist", func(t *testing.T) {
		err := service.DeletePlaylist(ctx, userID, clientID, "playlist123")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not implemented")
	})

	t.Run("TestAddItemToPlaylist", func(t *testing.T) {
		err := service.AddItemToPlaylist(ctx, userID, clientID, "playlist123", "newitem")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not implemented")
	})

	t.Run("TestRemoveItemFromPlaylist", func(t *testing.T) {
		err := service.RemoveItemFromPlaylist(ctx, userID, clientID, "playlist123", "item1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not implemented")
	})
}

// Let's add integration tests for all media client types that support playlists
func TestClientMediaPlaylistsIntegration(t *testing.T) {
	// Skip if not running integration tests
	if os.Getenv("INTEGRATION") != "true" {
		t.Skip("Skipping integration test. Set INTEGRATION=true to run")
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create test cases for each media client type
	testCases := []struct {
		name          string
		testEnvPrefix string
		newClientFunc func(context.Context, uint64, interface{}) (media.ClientMedia, error)
		configFactory func(baseURL, apiKey, user, password string) interface{}
	}{
		{
			name:          "Emby",
			testEnvPrefix: "EMBY_TEST",
			newClientFunc: func(ctx context.Context, id uint64, config interface{}) (media.ClientMedia, error) {
				return media.NewEmbyClient(ctx, id, config.(types.EmbyConfig))
			},
			configFactory: func(baseURL, apiKey, user, password string) interface{} {
				return types.EmbyConfig{
					BaseClientMediaConfig: types.BaseClientMediaConfig{
						BaseURL: baseURL,
						APIKey:  apiKey,
					},
					Username: user,
					Password: password,
				}
			},
		},
		{
			name:          "Jellyfin",
			testEnvPrefix: "JELLYFIN_TEST",
			newClientFunc: func(ctx context.Context, id uint64, config interface{}) (media.ClientMedia, error) {
				return media.NewJellyfinClient(ctx, id, config.(types.JellyfinConfig))
			},
			configFactory: func(baseURL, apiKey, user, password string) interface{} {
				return types.JellyfinConfig{
					BaseClientMediaConfig: types.BaseClientMediaConfig{
						BaseURL: baseURL,
						APIKey:  apiKey,
					},
					Username: user,
					Password: password,
				}
			},
		},
		{
			name:          "Plex",
			testEnvPrefix: "PLEX_TEST",
			newClientFunc: func(ctx context.Context, id uint64, config interface{}) (media.ClientMedia, error) {
				return media.NewPlexClient(ctx, id, config.(types.PlexConfig))
			},
			configFactory: func(baseURL, apiKey, user, password string) interface{} {
				return types.PlexConfig{
					BaseClientMediaConfig: types.BaseClientMediaConfig{
						BaseURL: baseURL,
						APIKey:  apiKey,
					},
					Username: user,
					Password: password,
				}
			},
		},
	}

	// Run tests for each client type
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Get test credentials from environment
			baseURL := os.Getenv(tc.testEnvPrefix + "_URL")
			apiKey := os.Getenv(tc.testEnvPrefix + "_API_KEY")
			user := os.Getenv(tc.testEnvPrefix + "_USER")
			password := os.Getenv(tc.testEnvPrefix + "_PASSWORD")

			if baseURL == "" || apiKey == "" || user == "" {
				t.Skipf("Skipping %s test - missing environment variables", tc.name)
				return
			}

			// Create client configuration
			config := tc.configFactory(baseURL, apiKey, user, password)

			// Initialize client with ID 1
			client, err := tc.newClientFunc(ctx, 1, config)
			if err != nil {
				t.Skipf("Failed to initialize %s client: %v", tc.name, err)
				return
			}

			// Check if client supports playlists
			playlistProvider, ok := client.(providers.PlaylistProvider)
			if !ok {
				t.Skipf("%s client does not implement PlaylistProvider interface", tc.name)
				return
			}

			if !playlistProvider.SupportsPlaylists() {
				t.Skipf("%s client does not support playlists", tc.name)
				return
			}

			// Test getting playlists
			playlists, err := playlistProvider.GetPlaylists(ctx, &mediatypes.QueryOptions{
				Limit: 5,
			})

			if err != nil {
				t.Errorf("Error getting playlists from %s: %v", tc.name, err)
				return
			}

			// Log number of playlists found
			t.Logf("Found %d playlists in %s", len(playlists), tc.name)

			// Verify playlist structure if any were found
			if len(playlists) > 0 {
				playlist := playlists[0]
				
				// Log details of first playlist
				t.Logf("First playlist: %s (ID: %s, Items: %d)", 
					playlist.Data.Details.Title, 
					playlist.ExternalID, 
					playlist.Data.ItemCount)

				// Test the MediaItem structure to ensure it conforms to our expected format
				assert.NotEmpty(t, playlist.ExternalID, "Playlist should have an external ID")
				assert.NotEmpty(t, playlist.Data.Details.Title, "Playlist should have a title")
				assert.True(t, playlist.Data.ItemCount >= 0, "Playlist should have a valid item count")
				assert.NotEqual(t, time.Time{}, playlist.Data.Details.AddedAt, "Playlist should have an added date")
				
				// If there are items, verify their structure
				if playlist.Data.ItemCount > 0 && len(playlist.Data.ItemIDs) > 0 {
					assert.NotEmpty(t, playlist.Data.ItemIDs[0], "Playlist items should have valid IDs")
				}
			}
		})
	}
}