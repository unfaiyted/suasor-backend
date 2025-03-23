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
	"suasor/client/media/interfaces"
	"suasor/models"
	logger "suasor/utils"
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
	config := models.SubsonicConfig{
		Host:     host,
		Port:     portNum,
		Username: username,
		Password: password,
		SSL:      ssl,
	}

	logger.Initialize()
	ctx := context.Background()

	// Initialize client
	client := NewSubsonicClient(1, config)
	require.NotNil(t, client)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Run all test cases that are supported by Subsonic
	t.Run("TestGetMusic", func(t *testing.T) {
		testGetMusic(t, ctx, client)
	})

	t.Run("TestGetMusicTrackByID", func(t *testing.T) {
		testGetMusicTrackByID(t, ctx, client)
	})

	t.Run("TestGetMusicArtists", func(t *testing.T) {
		testGetMusicArtists(t, ctx, client)
	})

	t.Run("TestGetMusicAlbums", func(t *testing.T) {
		testGetMusicAlbums(t, ctx, client)
	})

	t.Run("TestGetMusicGenres", func(t *testing.T) {
		testGetMusicGenres(t, ctx, client)
	})

	t.Run("TestGetPlaylists", func(t *testing.T) {
		testGetPlaylists(t, ctx, client)
	})

	t.Run("TestGetPlaylistItems", func(t *testing.T) {
		testGetPlaylistItems(t, ctx, client)
	})

	t.Run("TestStreamURL", func(t *testing.T) {
		testStreamURL(t, ctx, client)
	})

	t.Run("TestCoverArtURL", func(t *testing.T) {
		testCoverArtURL(t, ctx, client)
	})
}

// Test getting music tracks from Subsonic
func testGetMusic(t *testing.T, ctx context.Context, client interfaces.MediaContentProvider) {
	// Get tracks with limit
	options := &interfaces.QueryOptions{
		Limit: 10,
	}

	tracks, err := client.GetMusic(ctx, options)
	require.NoError(t, err)

	// Validate results
	assert.NotEmpty(t, tracks, "Expected to get at least one music track")
	if len(tracks) > 0 {
		track := tracks[0]
		t.Logf("Got track: %s (ID: %s)", track.Metadata.Title, track.ExternalID)

		// Verify track has expected fields
		assert.NotEmpty(t, track.ExternalID)
		assert.NotEmpty(t, track.Metadata.Title)
		assert.NotEmpty(t, track.ArtistName, "Expected track to have artist name")
	}
}

// Test getting a specific music track by ID
func testGetMusicTrackByID(t *testing.T, ctx context.Context, client interfaces.MediaContentProvider) {
	// First get a list of tracks to get a valid ID
	tracks, err := client.GetMusic(ctx, &interfaces.QueryOptions{Limit: 1})
	require.NoError(t, err)
	require.NotEmpty(t, tracks, "Need at least one track to test GetMusicTrackByID")

	trackID := tracks[0].ExternalID

	// Get the specific track
	track, err := client.GetMusicTrackByID(ctx, trackID)
	require.NoError(t, err)

	// Validate the result
	assert.Equal(t, trackID, track.ExternalID)
	assert.NotEmpty(t, track.Metadata.Title)
	assert.NotEmpty(t, track.ArtistName)
}

// Test getting music artists
func testGetMusicArtists(t *testing.T, ctx context.Context, client interfaces.MediaContentProvider) {
	artists, err := client.GetMusicArtists(ctx, &interfaces.QueryOptions{Limit: 5})
	require.NoError(t, err)

	if len(artists) > 0 {
		t.Logf("Got %d artists", len(artists))
		artist := artists[0]
		assert.NotEmpty(t, artist.ExternalID)
		assert.NotEmpty(t, artist.Metadata.Title)
	} else {
		t.Log("No artists found in library")
	}
}

// Test getting music albums
func testGetMusicAlbums(t *testing.T, ctx context.Context, client interfaces.MediaContentProvider) {
	albums, err := client.GetMusicAlbums(ctx, &interfaces.QueryOptions{Limit: 5})
	require.NoError(t, err)

	if len(albums) > 0 {
		t.Logf("Got %d albums", len(albums))
		album := albums[0]
		assert.NotEmpty(t, album.ExternalID)
		assert.NotEmpty(t, album.Metadata.Title)
		assert.NotEmpty(t, album.ArtistName, "Expected album to have artist name")
		assert.GreaterOrEqual(t, album.TrackCount, 0, "Track count should be at least 0")
	} else {
		t.Log("No albums found in library")
	}
}

// Test getting music genres
func testGetMusicGenres(t *testing.T, ctx context.Context, client interfaces.MediaContentProvider) {
	genres, err := client.GetMusicGenres(ctx)
	require.NoError(t, err)

	t.Logf("Got %d music genres", len(genres))
	if len(genres) > 0 {
		t.Logf("Some music genres: %v", genres[:min(3, len(genres))])
	}
}

// Test getting playlists
func testGetPlaylists(t *testing.T, ctx context.Context, client interfaces.MediaContentProvider) {
	playlists, err := client.GetPlaylists(ctx, &interfaces.QueryOptions{Limit: 5})
	require.NoError(t, err)

	if len(playlists) > 0 {
		t.Logf("Got %d playlists", len(playlists))
		playlist := playlists[0]
		assert.NotEmpty(t, playlist.ExternalID)
		assert.NotEmpty(t, playlist.Metadata.Title)
	} else {
		t.Log("No playlists found in library")
	}
}

// Test getting playlist items
func testGetPlaylistItems(t *testing.T, ctx context.Context, client *SubsonicClient) {
	// Get a playlist first
	playlists, err := client.GetPlaylists(ctx, &interfaces.QueryOptions{Limit: 1})
	if err != nil || len(playlists) == 0 {
		t.Skip("No playlists available to test items")
	}

	playlistID := playlists[0].ExternalID

	// Get tracks for the playlist
	tracks, err := client.GetPlaylistItems(ctx, playlistID)
	require.NoError(t, err)

	if len(tracks) > 0 {
		t.Logf("Got %d tracks for playlist '%s'", len(tracks), playlists[0].Metadata.Title)
		track := tracks[0]
		assert.NotEmpty(t, track.ExternalID)
		assert.NotEmpty(t, track.Metadata.Title)
		assert.NotEmpty(t, track.ArtistName)
	} else {
		t.Logf("No tracks found for the playlist '%s'", playlists[0].Metadata.Title)
	}
}

// Test getting stream URL
func testStreamURL(t *testing.T, ctx context.Context, client *SubsonicClient) {
	// First get a track to get a valid ID
	tracks, err := client.GetMusic(ctx, &interfaces.QueryOptions{Limit: 1})
	if err != nil || len(tracks) == 0 {
		t.Skip("No tracks available to test stream URL")
	}

	trackID := tracks[0].ExternalID

	// Get the stream URL
	streamURL, err := client.GetStreamURL(ctx, trackID)
	require.NoError(t, err)
	assert.NotEmpty(t, streamURL)

	t.Logf("Stream URL for track '%s': %s", tracks[0].Metadata.Title, streamURL)
}

// Test getting cover art URL
func testCoverArtURL(t *testing.T, ctx context.Context, client *SubsonicClient) {
	// First get an album to get a valid cover art ID
	albums, err := client.GetMusicAlbums(ctx, &interfaces.QueryOptions{Limit: 1})
	if err != nil || len(albums) == 0 || albums[0].Metadata.Artwork.Poster == "" {
		t.Skip("No albums with cover art available to test")
	}

	coverURL := albums[0].Metadata.Artwork.Poster
	assert.NotEmpty(t, coverURL)

	t.Logf("Cover art URL for album '%s': %s", albums[0].Metadata.Title, coverURL)
}

// Helper function for min value
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
