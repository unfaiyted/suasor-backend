// emby_client_integration_test.go
package emby

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"suasor/client/media"
	"suasor/client/media/providers"
	"suasor/client/media/types"

	"github.com/joho/godotenv"
)

func init() {
	// Try to load .env file from several possible locations
	locations := []string{
		".env",          // Current directory
		"../../../.env", // Project root
		filepath.Join(os.Getenv("HOME"), "emby_test.env"), // Home directory
	}

	for _, location := range locations {
		err := godotenv.Load(location)
		if err == nil {
			fmt.Printf("Loaded environment from: %s\n", location)
			break
		}
	}
}

// Integration test for EmbyClient
// To run these tests:
// EMBY_TEST_URL=http://your-server:8096 EMBY_TEST_API_KEY=your-api-key EMBY_TEST_USER_ID=your-user-id go test -v -tags=integration

func TestEmbyClientIntegration(t *testing.T) {
	// Skip if not running integration tests or missing environment variables
	if os.Getenv("INTEGRATION") != "true" {
		t.Skip("Skipping integration test. Set INTEGRATION=true to run")
	}

	// Get test credentials from environment
	baseURL := os.Getenv("EMBY_TEST_URL")
	apiKey := os.Getenv("EMBY_TEST_API_KEY")
	user := os.Getenv("EMBY_TEST_USER")

	if baseURL == "" || apiKey == "" || user == "" {
		t.Fatal("Missing required environment variables for integration test")
	}

	// Create client configuration
	config := types.EmbyConfig{
		BaseURL:  baseURL,
		APIKey:   apiKey,
		Username: user,
	}

	ctx := context.Background()
	// Initialize client
	client, err := NewEmbyClient(ctx, 1, config)
	require.NoError(t, err)
	require.NotNil(t, client)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Run all test cases
	if movieProvider, ok := media.AsMovieProvider(client); ok {
		t.Run("TestMovieProvider", func(t *testing.T) {
			testGetMovies(t, ctx, movieProvider)
		})
	} else {
		t.Log("Client does not support MovieProvider interface")
	}

	if tvProvider, ok := media.AsTVShowProvider(client); ok {
		t.Run("TestTVShowProvider", func(t *testing.T) {
			testGetTVShowEpisodes(t, ctx, tvProvider)
			testGetTVShows(t, ctx, tvProvider)
			testGetTVShowSeasons(t, ctx, tvProvider)
		})
	} else {
		t.Log("Client does not support TVShowProvider interface")
	}

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

	if collectionProvider, ok := media.AsCollectionProvider(client); ok {
		t.Run("TestCollectionProvider", func(t *testing.T) {
			testGetCollections(t, ctx, collectionProvider)
		})
	} else {
		t.Log("Client does not support CollectionProvider interface")
	}

	if watchHistoryProvider, ok := media.AsHistoryProvider(client); ok {
		t.Run("TestWatchHistoryProvider", func(t *testing.T) {
			testGetWatchHistory(t, ctx, watchHistoryProvider)
		})
	} else {
		t.Log("Client does not support WatchHistoryProvider interface")
	}
}

// Test getting movies from Emby
func testGetMovies(t *testing.T, ctx context.Context, client providers.MovieProvider) {
	// Get movies with limit
	options := &types.QueryOptions{
		Limit: 10,
		Sort:  "SortName",
	}

	movies, err := client.GetMovies(ctx, options)
	require.NoError(t, err)

	// Validate results
	assert.NotEmpty(t, movies, "Expected to get at least one movie")
	if len(movies) > 0 {
		movie := movies[0]
		t.Logf("Got movie: %s (ID: %s)", movie.Data.Details.Title, movie.ExternalID)

		// Verify movie has expected fields
		assert.NotEmpty(t, movie.ExternalID)
		assert.NotEmpty(t, movie.Data.Details.Title)
		assert.NotEmpty(t, movie.Data.Details.Artwork.Poster, "Expected movie to have a poster image")
	}
}

// Test getting a specific movie by ID
func testGetSpecificMovie(t *testing.T, ctx context.Context, client providers.MovieProvider) {
	// First get a list of movies to get a valid ID
	movies, err := client.GetMovies(ctx, &types.QueryOptions{Limit: 1})
	require.NoError(t, err)
	require.NotEmpty(t, movies, "Need at least one movie to test GetMovieByID")

	movieID := movies[0].ExternalID

	// Get the specific movie
	movie, err := client.GetMovieByID(ctx, movieID)
	require.NoError(t, err)

	// Validate the result
	assert.Equal(t, movieID, movie.ExternalID)
	assert.NotEmpty(t, movie.Data.Details.Title)
	assert.NotEmpty(t, movie.Data.Details.ReleaseYear)

}

// Test getting TV shows
func testGetTVShows(t *testing.T, ctx context.Context, client providers.TVShowProvider) {
	shows, err := client.GetTVShows(ctx, &types.QueryOptions{Limit: 5})
	require.NoError(t, err)

	if len(shows) > 0 {
		t.Logf("Got %d TV shows", len(shows))
		show := shows[0]
		assert.NotEmpty(t, show.ExternalID)
		assert.NotEmpty(t, show.Data.Details.Title)
		assert.NotEmpty(t, show.Data.Details.Artwork.Poster, "Expected TV show to have a poster image")

		// Test GetTVShowByID with the first show
		showByID, err := client.GetTVShowByID(ctx, show.ExternalID)
		require.NoError(t, err)
		assert.Equal(t, show.ExternalID, showByID.ExternalID)
	} else {
		t.Log("No TV shows found in library to test")
	}
}

// Test getting TV show seasons
func testGetTVShowSeasons(t *testing.T, ctx context.Context, client providers.TVShowProvider) {
	// Get a TV show first
	shows, err := client.GetTVShows(ctx, &types.QueryOptions{Limit: 1})
	if err != nil || len(shows) == 0 {
		t.Skip("No TV shows available to test seasons")
	}

	showID := shows[0].ExternalID

	// Get seasons for the show
	seasons, err := client.GetTVShowSeasons(ctx, showID)
	require.NoError(t, err)

	if len(seasons) > 0 {
		t.Logf("Got %d seasons for show '%s'", len(seasons), shows[0].Data.Details.Title)
		season := seasons[0]
		assert.NotEmpty(t, season.ExternalID)
		assert.NotEmpty(t, season.Data.Details.Title)
		assert.Greater(t, season.Data.Number, 0, "Season number should be greater than 0")
	} else {
		t.Log("No seasons found for the TV show")
	}
}

// Test getting TV show episodes
func testGetTVShowEpisodes(t *testing.T, ctx context.Context, client providers.TVShowProvider) {
	// Get a TV show first
	shows, err := client.GetTVShows(ctx, &types.QueryOptions{Limit: 1})
	if err != nil || len(shows) == 0 {
		t.Skip("No TV shows available to test episodes")
	}

	showID := shows[0].ExternalID

	// Get seasons
	seasons, err := client.GetTVShowSeasons(ctx, showID)
	if err != nil || len(seasons) == 0 {
		t.Skip("No seasons available to test episodes")
	}

	// Get episodes for the first season
	episodes, err := client.GetTVShowEpisodes(ctx, showID, seasons[0].Data.Number)
	require.NoError(t, err)

	if len(episodes) > 0 {
		t.Logf("Got %d episodes for season %d of show '%s'",
			len(episodes), seasons[0].Data.Number, shows[0].Data.Details.Title)

		episode := episodes[0]
		assert.NotEmpty(t, episode.ExternalID)
		assert.NotEmpty(t, episode.Data.Details.Title)
		assert.Equal(t, showID, episode.Data.ShowID)

		// Test GetEpisodeByID
		episodeByID, err := client.GetEpisodeByID(ctx, episode.ExternalID)
		require.NoError(t, err)
		assert.Equal(t, episode.ExternalID, episodeByID.ExternalID)
	} else {
		t.Log("No episodes found for the season")
	}
}

// Test getting collections
func testGetCollections(t *testing.T, ctx context.Context, client providers.CollectionProvider) {
	collections, err := client.GetCollections(ctx, &types.QueryOptions{Limit: 5})
	require.NoError(t, err)

	if len(collections) > 0 {
		t.Logf("Got %d collections", len(collections))
		collection := collections[0]
		assert.NotEmpty(t, collection.ExternalID)
		assert.NotEmpty(t, collection.Data.Details.Title)
	} else {
		t.Log("No collections found in library")
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
		assert.NotEmpty(t, albums[0].ExternalID)
		assert.NotEmpty(t, albums[0].Data.Details.Title)
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
		assert.NotEmpty(t, playlists[0].ExternalID)
		assert.NotEmpty(t, playlists[0].Data.Details.Title)
	} else {
		t.Log("No playlists found in library")
	}
}

// Test getting watch history
func testGetWatchHistory(t *testing.T, ctx context.Context, client providers.HistoryProvider) {
	history, err := client.GetPlayHistory(ctx, &types.QueryOptions{Limit: 10})
	require.NoError(t, err)

	t.Logf("Got %d watch history items", len(history))
	if len(history) > 0 {

		if movieData, ok := history[0].Item.Data.(types.Movie); ok {
			t.Logf("Movie Title: %s", movieData.Details.Title)
			assert.NotEmpty(t, movieData.Details.Title)
		}

		if episodeData, ok := history[0].Item.Data.(types.Episode); ok {
			t.Logf("Episode Title: %s", episodeData.Details.Title)
			assert.NotEmpty(t, episodeData.Details.Title)
		}

		assert.NotEqual(t, time.Time{}, history[0].WatchedAt, "Expected watch date to be set")
	}
}

// Test getting genres
func testGetGenres(t *testing.T, ctx context.Context, client providers.MediaContentProvider) {
	// Movie genres
	movieGenres, err := client.GetMovieGenres(ctx)
	require.NoError(t, err)
	t.Logf("Got %d movie genres", len(movieGenres))
	if len(movieGenres) > 0 {
		t.Logf("Some movie genres: %v", movieGenres[:min(3, len(movieGenres))])
	}

	// Music genres
	musicGenres, err := client.GetMusicGenres(ctx)
	require.NoError(t, err)
	t.Logf("Got %d music genres", len(musicGenres))
	if len(musicGenres) > 0 {
		t.Logf("Some music genres: %v", musicGenres[:min(3, len(musicGenres))])
	}
}

// Helper function for min value
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
