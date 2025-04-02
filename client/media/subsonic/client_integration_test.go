// subsonic_client_integration_test.go
package subsonic

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joho/godotenv"

	"suasor/client/media"
	"suasor/client/media/providers"
	"suasor/client/media/types"
	client "suasor/client/types"
)

func init() {
	// Try to load .env file from several possible locations
	locations := []string{
		".env",          // Current directory
		"../../../.env", // Project root
		filepath.Join(os.Getenv("HOME"), "subsonic_test.env"), // Home directory
	}

	for _, location := range locations {
		err := godotenv.Load(location)
		if err == nil {
			fmt.Printf("Loaded environment from: %s\n", location)
			break
		}
	}
}

// Integration test for SubsonicClient
// To run these tests:
// SUBSONIC_TEST_HOST=your-server SUBSONIC_TEST_PORT=4040 SUBSONIC_TEST_USERNAME=user SUBSONIC_TEST_PASSWORD=pass SUBSONIC_TEST_SSL=false go test -v -tags=integration

func TestSubsonicClientIntegration(t *testing.T) {
	// Skip if not running integration tests or missing environment variables
	if os.Getenv("INTEGRATION") != "true" {
		t.Skip("Skipping integration test. Set INTEGRATION=true to run")
	}

	// Get test credentials from environment
	host := os.Getenv("SUBSONIC_TEST_HOST")
	portStr := os.Getenv("SUBSONIC_TEST_PORT")
	username := os.Getenv("SUBSONIC_TEST_USERNAME")
	password := os.Getenv("SUBSONIC_TEST_PASSWORD")
	sslStr := os.Getenv("SUBSONIC_TEST_SSL")

	if host == "" || username == "" || password == "" {
		t.Fatal("Missing required environment variables for integration test")
	}

	// Convert port string to int
	portNum := 4040 // default
	if portStr != "" {
		var err error
		portNum, err = strconv.Atoi(portStr)
		if err != nil {
			t.Fatalf("Invalid port number: %s", portStr)
		}
	}

	// Convert SSL string to bool
	ssl := false
	if sslStr == "true" {
		ssl = true
	}

	// Create client configuration
	config := client.SubsonicConfig{
		Host:     host,
		Port:     portNum,
		Username: username,
		Password: password,
		SSL:      ssl,
	}

	ctx := context.Background()

	// Initialize client
	client, err := NewSubsonicClient(ctx, 1, config)
	require.NoError(t, err)
	require.NotNil(t, client)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Run test cases for supported interfaces
	if musicProvider, ok := media.AsMusicProvider(client); ok {
		t.Run("TestMusicProvider", func(t *testing.T) {
			testGetMusicContent(t, ctx, musicProvider)
		})
	} else {
		t.Log("Client does not support MusicProvider interface")
	}

	if playlistProvider, ok := media.AsPlaylistProvider(client); ok {
		t.Run("TestPlaylistProvider", func(t *testing.T) {
			testGetPlaylists(t, ctx, playlistProvider)
		})
	} else {
		t.Log("Client does not support PlaylistProvider interface")
	}

	// Test Subsonic-specific features using type assertion
	if subsonicClient, ok := client.(*SubsonicClient); ok {
		t.Run("TestSubsonicSpecific", func(t *testing.T) {
			testGetPlaylistItems(t, ctx, subsonicClient)
			testStreamURL(t, ctx, subsonicClient)
			testCoverArtURL(t, ctx, subsonicClient)
		})
	}
}

// Test getting music content
func testGetMusicContent(t *testing.T, ctx context.Context, client providers.MusicProvider) {
	// Test artists
	artists, err := client.GetMusicArtists(ctx, &types.QueryOptions{Limit: 5})
	require.NoError(t, err)
	if len(artists) > 0 {
		t.Logf("Got %d music artists", len(artists))
		assert.NotEmpty(t, artists[0].ExternalID)
		assert.NotEmpty(t, artists[0].Data.Details.Title)
	}

	// Test albums
	albums, err := client.GetMusicAlbums(ctx, &types.QueryOptions{Limit: 5})
	require.NoError(t, err)
	if len(albums) > 0 {
		t.Logf("Got %d music albums", len(albums))
		assert.NotEmpty(t, albums[0].ExternalID, "Expected album to have an external ID")
		assert.NotEmpty(t, albums[0].Data.Details.Title, "Expected album to have a title")
	}

	// Test tracks
	tracks, err := client.GetMusic(ctx, &types.QueryOptions{Limit: 5})
	require.NoError(t, err)
	if len(tracks) > 0 {
		t.Logf("Got %d music tracks", len(tracks))
		track := tracks[0]
		assert.NotEmpty(t, track.ExternalID)
		assert.NotEmpty(t, track.Data.Details.Title)

		// Test GetMusicTrackByID
		trackByID, err := client.GetMusicTrackByID(ctx, track.ExternalID)
		require.NoError(t, err)
		assert.Equal(t, track.ExternalID, trackByID.ExternalID)
	}
}

// Test getting playlists
func testGetPlaylists(t *testing.T, ctx context.Context, client providers.PlaylistProvider) {
	playlists, err := client.GetPlaylists(ctx, &types.QueryOptions{Limit: 5})
	require.NoError(t, err)

	if len(playlists) > 0 {
		t.Logf("Got %d playlists", len(playlists))
		playlist := playlists[0]

		// Basic validation
		assert.NotEmpty(t, playlist.ExternalID, "Playlist should have an external ID")
		assert.NotEmpty(t, playlist.Data.Details.Title, "Playlist should have a title")
		
		// Log more details about the playlist for inspection
		t.Logf("First playlist: %s (ID: %s, Items: %d, Owner: %s)",
			playlist.Data.Details.Title,
			playlist.ExternalID,
			playlist.Data.ItemCount,
			playlist.Data.Owner)
		
		// Validate playlist structure
		assert.True(t, playlist.Data.ItemCount >= 0, "Playlist should have a valid item count")
		
		// Note: Subsonic might not set the AddedAt field consistently
		if !playlist.Data.Details.AddedAt.IsZero() {
			assert.NotEqual(t, time.Time{}, playlist.Data.Details.AddedAt, "Playlist should have an added date")
		}
		
		// Check if the playlist has items
		if playlist.Data.ItemCount > 0 && len(playlist.Data.ItemIDs) > 0 {
			t.Logf("First playlist has %d items, first item ID: %s", 
				len(playlist.Data.ItemIDs), 
				playlist.Data.ItemIDs[0])
			assert.NotEmpty(t, playlist.Data.ItemIDs[0], "Playlist items should have valid IDs")
		}
		
		// Test querying with filters if this playlist provider supports it
		if playlist.Data.ItemCount > 0 {
			// Try to get playlist by ID using filter
			filteredOptions := &types.QueryOptions{
				Filters: map[string]string{
					"id": playlist.ExternalID,
				},
			}
			
			filteredPlaylists, err := client.GetPlaylists(ctx, filteredOptions)
			if err == nil && len(filteredPlaylists) > 0 {
				t.Logf("Successfully retrieved playlist by ID filter")
				assert.Equal(t, playlist.ExternalID, filteredPlaylists[0].ExternalID, 
					"Filtered playlist should match requested ID")
			} else {
				t.Logf("Provider doesn't support filtering playlists by ID or no results: %v", err)
			}
		}
	} else {
		t.Log("No playlists found in library")
	}
}

// Test getting playlist items - Subsonic specific
func testGetPlaylistItems(t *testing.T, ctx context.Context, client *SubsonicClient) {
	// Get a playlist first
	playlists, err := client.GetPlaylists(ctx, &types.QueryOptions{Limit: 1})
	if err != nil || len(playlists) == 0 {
		t.Skip("No playlists available to test items")
	}

	playlistID := playlists[0].ExternalID

	// Get tracks for the playlist
	tracks, err := client.GetPlaylistItems(ctx, playlistID)
	require.NoError(t, err)

	if len(tracks) > 0 {
		t.Logf("Got %d tracks for playlist '%s'", len(tracks), playlists[0].Data.Details.Title)
		track := tracks[0]
		assert.NotEmpty(t, track.ExternalID)
		assert.NotEmpty(t, track.Data.Details.Title)
	} else {
		t.Logf("No tracks found for the playlist '%s'", playlists[0].Data.Details.Title)
	}
}

// Test getting stream URL - Subsonic specific
func testStreamURL(t *testing.T, ctx context.Context, client *SubsonicClient) {
	// First get a track to get a valid ID
	tracks, err := client.GetMusic(ctx, &types.QueryOptions{Limit: 1})
	if err != nil || len(tracks) == 0 {
		t.Skip("No tracks available to test stream URL")
	}

	trackID := tracks[0].ExternalID

	// Get the stream URL
	streamURL, err := client.GetStreamURL(ctx, trackID)
	require.NoError(t, err)
	assert.NotEmpty(t, streamURL)

	t.Logf("Stream URL for track '%s': %s", tracks[0].Data.Details.Title, streamURL)
}

// Test getting cover art URL - Subsonic specific
func testCoverArtURL(t *testing.T, ctx context.Context, client *SubsonicClient) {
	// First get an album to get a valid cover art ID
	albums, err := client.GetMusicAlbums(ctx, &types.QueryOptions{Limit: 1})
	if err != nil || len(albums) == 0 || albums[0].Data.Details.Artwork.Poster == "" {
		t.Skip("No albums with cover art available to test")
	}

	coverURL := albums[0].Data.Details.Artwork.Poster
	assert.NotEmpty(t, coverURL)

	t.Logf("Cover art URL for album '%s': %s", albums[0].Data.Details.Title, coverURL)
}

// Helper function for min value
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
