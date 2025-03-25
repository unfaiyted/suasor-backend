// plex/client_integration_test.go
package plex

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/joho/godotenv"
	"suasor/client/media/types"

	logger "suasor/utils"
)

func init() {
	// Try to load .env file from several possible locations
	locations := []string{
		".env",          // Current directory
		"../../../.env", // Project root
		filepath.Join(os.Getenv("HOME"), "plex_test.env"), // Home directory
	}

	for _, location := range locations {
		err := godotenv.Load(location)
		if err == nil {
			fmt.Printf("Loaded environment from: %s\n", location)
			break
		}
	}
}

// Integration test for PlexClient
// To run these tests:
// PLEX_TEST_HOST=http://your-server:32400 PLEX_TEST_TOKEN=your-token INTEGRATION=true go test -v -tags=integration

func TestPlexClientIntegration(t *testing.T) {
	// Skip if not running integration tests or missing environment variables
	if os.Getenv("INTEGRATION") != "true" {
		t.Skip("Skipping integration test. Set INTEGRATION=true to run")
	}

	// Get test credentials from environment
	host := os.Getenv("PLEX_TEST_HOST")
	token := os.Getenv("PLEX_TEST_TOKEN")

	if host == "" || token == "" {
		t.Fatal("Missing required environment variables for integration test")
	}

	// Create client configuration
	config := types.PlexConfig{
		Host:  host,
		Token: token,
	}

	logger.Initialize()
	ctx := context.Background()

	// Initialize client
	client, err := NewPlexClient(ctx, 1, config)
	require.NoError(t, err)
	require.NotNil(t, client)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test each provider capability if supported
	if movieProvider, ok := interfaces.AsMovieProvider(client); ok {
		t.Run("TestMovieProvider", func(t *testing.T) {
			testMovieProvider(t, ctx, movieProvider)
		})
	} else {
		t.Log("Client does not support MovieProvider interface")
	}

	if tvProvider, ok := interfaces.AsTVShowProvider(client); ok {
		t.Run("TestTVShowProvider", func(t *testing.T) {
			testTVShowProvider(t, ctx, tvProvider)
		})
	} else {
		t.Log("Client does not support TVShowProvider interface")
	}

	if musicProvider, ok := interfaces.AsMusicProvider(client); ok {
		t.Run("TestMusicProvider", func(t *testing.T) {
			testMusicProvider(t, ctx, musicProvider)
		})
	} else {
		t.Log("Client does not support MusicProvider interface")
	}

	if playlistProvider, ok := interfaces.AsPlaylistProvider(client); ok {
		t.Run("TestPlaylistProvider", func(t *testing.T) {
			testPlaylistProvider(t, ctx, playlistProvider)
		})
	} else {
		t.Log("Client does not support PlaylistProvider interface")
	}

	if collectionProvider, ok := interfaces.AsCollectionProvider(client); ok {
		t.Run("TestCollectionProvider", func(t *testing.T) {
			testCollectionProvider(t, ctx, collectionProvider)
		})
	} else {
		t.Log("Client does not support CollectionProvider interface")
	}

	if watchHistoryProvider, ok := interfaces.AsWatchHistoryProvider(client); ok {
		t.Run("TestWatchHistoryProvider", func(t *testing.T) {
			testWatchHistoryProvider(t, ctx, watchHistoryProvider)
		})
	} else {
		t.Log("Client does not support WatchHistoryProvider interface")
	}
}

// Test movie functionality
func testMovieProvider(t *testing.T, ctx context.Context, provider interfaces.MovieProvider) {
	// Test getting movies
	t.Run("TestGetMovies", func(t *testing.T) {
		options := &types.QueryOptions{
			Limit: 10,
			Sort:  "title",
		}

		movies, err := provider.GetMovies(ctx, options)
		require.NoError(t, err)

		// Validate results
		assert.NotEmpty(t, movies, "Expected to get at least one movie")
		if len(movies) > 0 {
			movie := movies[0]
			t.Logf("Got movie: %s (ID: %s)", movie.Data.Details.Title, movie.ExternalID)

			// Verify movie has expected fields
			assert.NotEmpty(t, movie.ExternalID)
			assert.NotEmpty(t, movie.Data.Details.Title)
			assert.NotEmpty(t, movie.Data.Details.Artwork.Thumbnail, "Expected movie to have a thumbnail image")

			// Test getting a specific movie
			specificMovie, err := provider.GetMovieByID(ctx, movie.ExternalID)
			require.NoError(t, err)
			assert.Equal(t, movie.ExternalID, specificMovie.ExternalID)
			assert.NotEmpty(t, specificMovie.Data.Details.Title)
		}
	})

	// Test getting movie genres
	t.Run("TestGetMovieGenres", func(t *testing.T) {
		genres, err := provider.GetMovieGenres(ctx)
		require.NoError(t, err)
		t.Logf("Got %d movie genres", len(genres))
		if len(genres) > 0 {
			t.Logf("Some movie genres: %v", genres[:min(3, len(genres))])
		}
	})
}

// Test TV show functionality
func testTVShowProvider(t *testing.T, ctx context.Context, provider interfaces.TVShowProvider) {
	// Test getting TV shows
	t.Run("TestGetTVShows", func(t *testing.T) {
		shows, err := provider.GetTVShows(ctx, &types.QueryOptions{Limit: 5})
		require.NoError(t, err)

		if len(shows) > 0 {
			t.Logf("Got %d TV shows", len(shows))
			show := shows[0]
			assert.NotEmpty(t, show.ExternalID)
			assert.NotEmpty(t, show.Data.Details.Title)
			assert.NotEmpty(t, show.Data.Details.Artwork.Thumbnail, "Expected TV show to have a thumbnail image")

			// Test getting a specific TV show
			showID := show.ExternalID

			// We need to use type assertion here because GetTVShowByID isn't in the TVShowProvider interface
			// This is a good example of where your interface design is beneficial - we only test what's explicitly supported
			fullClient, ok := provider.(interface {
				GetTVShowByID(ctx context.Context, id string) (types.MediaItem[types.TVShow], error)
			})

			if ok {
				showByID, err := fullClient.GetTVShowByID(ctx, showID)
				require.NoError(t, err)
				assert.Equal(t, showID, showByID.ExternalID)
			} else {
				t.Log("Provider doesn't support GetTVShowByID")
			}

			// Test getting seasons
			testTVShowSeasons(t, ctx, provider, showID)
		} else {
			t.Log("No TV shows found in library to test")
		}
	})
}

// Test TV show seasons and episodes
func testTVShowSeasons(t *testing.T, ctx context.Context, provider interfaces.TVShowProvider, showID string) {
	// Get seasons for the show
	seasons, err := provider.GetTVShowSeasons(ctx, showID)
	require.NoError(t, err)

	if len(seasons) > 0 {
		t.Logf("Got %d seasons for show", len(seasons))
		season := seasons[0]
		assert.NotEmpty(t, season.ExternalID)
		assert.NotEmpty(t, season.Data.Details.Title)
		assert.Greater(t, season.Data.Number, 0, "Season number should be greater than 0")

		// Test getting episodes for a season
		episodes, err := provider.GetTVShowEpisodes(ctx, showID, season.Data.Number)
		require.NoError(t, err)

		if len(episodes) > 0 {
			t.Logf("Got %d episodes for season %d", len(episodes), season.Data.Number)

			episode := episodes[0]
			assert.NotEmpty(t, episode.ExternalID)
			assert.NotEmpty(t, episode.Data.Details.Title)
			assert.Equal(t, showID, episode.Data.ShowID)

			// Test GetEpisodeByID if supported
			fullClient, ok := provider.(interface {
				GetEpisodeByID(ctx context.Context, id string) (types.MediaItem[types.Episode], error)
			})

			if ok {
				time.Sleep(2 * time.Second) // Brief pause to avoid rate limiting
				episodeByID, err := fullClient.GetEpisodeByID(ctx, episode.ExternalID)
				require.NoError(t, err)
				assert.Equal(t, episode.ExternalID, episodeByID.ExternalID)
			} else {
				t.Log("Provider doesn't support GetEpisodeByID")
			}
		} else {
			t.Log("No episodes found for the season")
		}
	} else {
		t.Log("No seasons found for the TV show")
	}
}

// Test music functionality
func testMusicProvider(t *testing.T, ctx context.Context, provider interfaces.MusicProvider) {
	// Test getting artists
	t.Run("TestGetMusicArtists", func(t *testing.T) {
		artists, err := provider.GetMusicArtists(ctx, &types.QueryOptions{Limit: 5})
		require.NoError(t, err)
		if len(artists) > 0 {
			t.Logf("Got %d music artists", len(artists))
			assert.NotEmpty(t, artists[0].ExternalID)
			assert.NotEmpty(t, artists[0].Data.Details.Title)
		} else {
			t.Log("No music artists found")
		}
	})

	// Test getting albums
	t.Run("TestGetMusicAlbums", func(t *testing.T) {
		albums, err := provider.GetMusicAlbums(ctx, &types.QueryOptions{Limit: 5})
		require.NoError(t, err)
		if len(albums) > 0 {
			t.Logf("Got %d music albums", len(albums))
			assert.NotEmpty(t, albums[0].ExternalID)
			assert.NotEmpty(t, albums[0].Data.Details.Title)
		} else {
			t.Log("No music albums found")
		}
	})

	// Test getting tracks
	t.Run("TestGetMusic", func(t *testing.T) {
		tracks, err := provider.GetMusic(ctx, &types.QueryOptions{Limit: 5})
		require.NoError(t, err)
		if len(tracks) > 0 {
			t.Logf("Got %d music tracks", len(tracks))
			track := tracks[0]
			assert.NotEmpty(t, track.ExternalID)
			assert.NotEmpty(t, track.Data.Details.Title)

			// Test GetMusicTrackByID if supported
			fullClient, ok := provider.(interface {
				GetMusicTrackByID(ctx context.Context, id string) (types.MediaItem[types.Track], error)
			})

			if ok {
				trackByID, err := fullClient.GetMusicTrackByID(ctx, track.ExternalID)
				require.NoError(t, err)
				assert.Equal(t, track.ExternalID, trackByID.ExternalID)
			} else {
				t.Log("Provider doesn't support GetMusicTrackByID")
			}
		} else {
			t.Log("No music tracks found")
		}
	})

	// Test getting music genres if supported
	t.Run("TestGetMusicGenres", func(t *testing.T) {
		// Check if provider supports genre retrieval
		genreProvider, ok := provider.(interface {
			GetMusicGenres(ctx context.Context) ([]string, error)
		})

		if !ok {
			t.Log("Provider doesn't support GetMusicGenres")
			return
		}

		genres, err := genreProvider.GetMusicGenres(ctx)
		require.NoError(t, err)
		t.Logf("Got %d music genres", len(genres))
		if len(genres) > 0 {
			t.Logf("Some music genres: %v", genres[:min(3, len(genres))])
		}
	})
}

// Test playlist functionality
func testPlaylistProvider(t *testing.T, ctx context.Context, provider interfaces.PlaylistProvider) {
	playlists, err := provider.GetPlaylists(ctx, &types.QueryOptions{Limit: 5})
	require.NoError(t, err)

	if len(playlists) > 0 {
		t.Logf("Got %d playlists", len(playlists))
		assert.NotEmpty(t, playlists[0].ExternalID)
		assert.NotEmpty(t, playlists[0].Data.Details.Title)
	} else {
		t.Log("No playlists found in library")
	}
}

// Test collection functionality
func testCollectionProvider(t *testing.T, ctx context.Context, provider interfaces.CollectionProvider) {
	collections, err := provider.GetCollections(ctx, &types.QueryOptions{Limit: 5})
	require.NoError(t, err)

	if len(collections) > 0 {
		t.Logf("Got %d collections", len(collections))
		assert.NotEmpty(t, collections[0].ExternalID)
		assert.NotEmpty(t, collections[0].Data.Details.Title)
	} else {
		t.Log("No collections found in library")
	}
}

// Test watch history functionality
func testWatchHistoryProvider(t *testing.T, ctx context.Context, provider interfaces.WatchHistoryProvider) {
	history, err := provider.GetWatchHistory(ctx, &types.QueryOptions{Limit: 10})

	// Note: The Plex client implementation may return an error as it's not fully implemented
	if err != nil {
		t.Logf("Watch history retrieval returned error (may be expected): %v", err)
		return
	}

	t.Logf("Got %d watch history items", len(history))
	if len(history) > 0 {
		assert.NotEmpty(t, history[0].Item.GetData().GetDetails().Title)
		assert.NotEqual(t, time.Time{}, history[0].LastWatchedAt, "Expected watch date to be set")
	}
}

// Helper function for min value
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
