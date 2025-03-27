package lidarr

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
	"suasor/client/automation/providers"
	auto "suasor/client/automation/types"
	"suasor/client/types"
	"suasor/types/models"
	"suasor/types/requests"

	logger "suasor/utils"
)

func init() {
	// Try to load .env file from several possible locations
	locations := []string{
		".env",          // Current directory
		"../../../.env", // Project root
		filepath.Join(os.Getenv("HOME"), "lidarr_test.env"), // Home directory
	}

	for _, location := range locations {
		err := godotenv.Load(location)
		if err == nil {
			fmt.Printf("Loaded environment from: %s\n", location)
			break
		}
	}
}

// Integration test for LidarrClient
// To run these tests:
// INTEGRATION=true LIDARR_TEST_URL=http://your-server:8686 LIDARR_TEST_API_KEY=your-api-key go test -v -tags=integration

func TestLidarrClientIntegration(t *testing.T) {
	// Skip if not running integration tests or missing environment variables
	if os.Getenv("INTEGRATION") != "true" {
		t.Skip("Skipping integration test. Set INTEGRATION=true to run")
	}

	// Get test credentials from environment
	baseURL := os.Getenv("LIDARR_TEST_URL")
	apiKey := os.Getenv("LIDARR_TEST_API_KEY")

	if baseURL == "" || apiKey == "" {
		t.Fatal("Missing required environment variables for integration test")
	}

	// Create client configuration
	config := types.JellyfinConfig{
		BaseURL: baseURL,
		APIKey:  apiKey,
	}

	logger.Initialize()
	ctx := context.Background()

	// Initialize client
	client, err := NewLidarrClient(ctx, 1, config)
	require.NoError(t, err)
	require.NotNil(t, client)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	provider, ok := client.(providers.AutomationProvider)
	require.True(t, ok)

	// Run all test cases
	t.Run("TestGetSystemStatus", func(t *testing.T) {
		testGetSystemStatus(t, ctx, provider)
	})

	t.Run("TestGetLibraryItems", func(t *testing.T) {
		testGetLibraryItems(t, ctx, provider)
	})

	t.Run("TestGetAndSearchMedia", func(t *testing.T) {
		testGetAndSearchMedia(t, ctx, provider)
	})

	t.Run("TestQualityProfiles", func(t *testing.T) {
		testGetQualityProfiles(t, ctx, provider)
	})

	t.Run("TestMetadataProfiles", func(t *testing.T) {
		testGetMetadataProfiles(t, ctx, provider)
	})

	t.Run("TestTags", func(t *testing.T) {
		testGetAndCreateTags(t, ctx, provider)
	})

	t.Run("TestCalendar", func(t *testing.T) {
		testGetCalendar(t, ctx, provider)
	})

	t.Run("TestCommands", func(t *testing.T) {
		testExecuteCommand(t, ctx, provider)
	})

	// Add/Update/Delete operations should be in a specific test
	// that can be run separately to avoid modifying data
	if os.Getenv("TEST_MODIFICATIONS") == "true" {
		t.Run("TestAddUpdateDeleteMedia", func(t *testing.T) {
			testAddUpdateDeleteMedia(t, ctx, provider)
		})
	}
}

// Test getting system status from Lidarr
func testGetSystemStatus(t *testing.T, ctx context.Context, client providers.AutomationProvider) {
	status, err := client.GetSystemStatus(ctx)
	require.NoError(t, err)

	// Validate results
	assert.NotEmpty(t, status.Version, "Expected to get a version string")
	t.Logf("Lidarr version: %s", status.Version)
	assert.NotEmpty(t, status.OsName, "Expected to get an OS name")
}

// Test getting library items from Lidarr
func testGetLibraryItems(t *testing.T, ctx context.Context, client providers.AutomationProvider) {
	// Get library items with limit
	options := &auto.LibraryQueryOptions{
		Limit: 10,
	}

	artists, err := client.GetLibraryItems(ctx, options)
	require.NoError(t, err)

	// Log the count even if empty
	t.Logf("Found %d artists in library", len(artists))

	if len(artists) > 0 {
		artist := artists[0]
		t.Logf("Got artist: %s (ID: %d)", artist.Title, artist.ID)

		// Verify artist has expected fields
		assert.NotZero(t, artist.ID)
		assert.NotEmpty(t, artist.Title)
		assert.NotEmpty(t, artist.Path, "Expected artist to have a path")
	}
}

// Test getting and searching for media
func testGetAndSearchMedia(t *testing.T, ctx context.Context, client providers.AutomationProvider) {
	// First, check if we have any artists in the library
	artists, err := client.GetLibraryItems(ctx, &auto.LibraryQueryOptions{Limit: 1})
	require.NoError(t, err)

	if len(artists) > 0 {
		// Test getting a specific artist by ID
		artistID := int64(artists[0].ID)
		artist, err := client.GetMediaByID(ctx, artistID)
		require.NoError(t, err)

		// Validate the result
		assert.Equal(t, artistID, int64(artist.ID))
		assert.NotEmpty(t, artist.Title)
	} else {
		t.Log("No artists in library to test GetMediaByID")
	}

	// Test search functionality
	searchTerm := "Taylor Swift" // A commonly available artist
	searchResults, err := client.SearchMedia(ctx, searchTerm, nil)
	require.NoError(t, err)

	t.Logf("Search for '%s' returned %d results", searchTerm, len(searchResults))
	if len(searchResults) > 0 {
		result := searchResults[0]
		assert.NotEmpty(t, result.Title)
	}
}

// Test getting quality profiles
func testGetQualityProfiles(t *testing.T, ctx context.Context, client providers.AutomationProvider) {
	profiles, err := client.GetQualityProfiles(ctx)
	require.NoError(t, err)

	// Validate results
	assert.NotEmpty(t, profiles, "Expected to get at least one quality profile")
	t.Logf("Found %d quality profiles", len(profiles))

	if len(profiles) > 0 {
		profile := profiles[0]
		assert.NotZero(t, profile.ID)
		assert.NotEmpty(t, profile.Name)
		t.Logf("Quality profile: %s (ID: %d)", profile.Name, profile.ID)
	}
}

// Test getting metadata profiles (specific to Lidarr)
func testGetMetadataProfiles(t *testing.T, ctx context.Context, client providers.AutomationProvider) {
	profiles, err := client.GetMetadataProfiles(ctx)
	require.NoError(t, err)

	// Validate results
	assert.NotEmpty(t, profiles, "Expected to get at least one metadata profile")
	t.Logf("Found %d metadata profiles", len(profiles))

	if len(profiles) > 0 {
		profile := profiles[0]
		assert.NotZero(t, profile.ID)
		assert.NotEmpty(t, profile.Name)
		t.Logf("Metadata profile: %s (ID: %d)", profile.Name, profile.ID)
	}
}

// Test getting and creating tags
func testGetAndCreateTags(t *testing.T, ctx context.Context, client providers.AutomationProvider) {
	// Get existing tags
	tags, err := client.GetTags(ctx)
	require.NoError(t, err)

	t.Logf("Found %d tags", len(tags))

	// Only test tag creation if explicitly enabled
	if os.Getenv("TEST_TAG_CREATION") == "true" {
		// Create a new tag with a unique name
		tagName := fmt.Sprintf("test-tag-%d", time.Now().Unix())
		newTag, err := client.CreateTag(ctx, tagName)
		require.NoError(t, err)

		assert.Equal(t, tagName, newTag.Name)
		assert.NotZero(t, newTag.ID)
		t.Logf("Created new tag: %s (ID: %d)", newTag.Name, newTag.ID)

		// Verify the tag was created by getting all tags again
		updatedTags, err := client.GetTags(ctx)
		require.NoError(t, err)
		assert.Greater(t, len(updatedTags), len(tags), "Tag list should have grown")
	}
}

// Test getting calendar
func testGetCalendar(t *testing.T, ctx context.Context, client providers.AutomationProvider) {
	// Get calendar for the next 30 days
	now := time.Now()
	end := now.AddDate(0, 0, 30) // 30 days from now

	calendar, err := client.GetCalendar(ctx, now, end)
	require.NoError(t, err)

	t.Logf("Found %d items in calendar for the next 30 days", len(calendar))

	if len(calendar) > 0 {
		item := calendar[0]
		assert.NotEmpty(t, item.Title)
		assert.NotZero(t, item.ID)
	}
}

// Test executing a command
func testExecuteCommand(t *testing.T, ctx context.Context, client providers.AutomationProvider) {
	// Only run if explicitly enabled
	if os.Getenv("TEST_COMMANDS") != "true" {
		t.Skip("Skipping command execution. Set TEST_COMMANDS=true to run")
		return
	}

	// Execute a simple command like RefreshMonitoredDownloads
	command := auto.Command{
		Name: "RefreshMonitoredDownloads",
	}

	result, err := client.ExecuteCommand(ctx, command)
	require.NoError(t, err)

	assert.Equal(t, command.Name, result.Name)
	assert.NotZero(t, result.ID)
	assert.NotEmpty(t, result.Status)
	t.Logf("Command %s executed with ID %d, status: %s", result.Name, result.ID, result.Status)
}

// Test adding, updating, and deleting media
// This test is potentially destructive and should be run with caution
func testAddUpdateDeleteMedia(t *testing.T, ctx context.Context, client providers.AutomationProvider) {
	// First, get quality profiles to use a valid profile ID
	profiles, err := client.GetQualityProfiles(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, profiles, "Need at least one quality profile")

	// Get metadata profiles
	metadataProfiles, err := client.GetMetadataProfiles(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, metadataProfiles, "Need at least one metadata profile")

	// Search for an artist to add
	searchResults, err := client.SearchMedia(ctx, "Coldplay", nil)
	require.NoError(t, err)
	require.NotEmpty(t, searchResults, "Need at least one search result")

	// Find a result that has a MusicBrainzID
	var searchItem models.AutomationMediaItem[auto.AutomationData]
	// for _, result := range searchResults {
	// if result.Data.MusicBrainzID != "" {
	// searchItem = result
	// break
	// }
	// }
	// require.NotEmpty(t, searchItem.Data.MusicBrainzID, "Need a search result with MusicBrainzID")

	// Prepare add request
	addRequest := requests.AutomationMediaAddRequest{
		Title:             searchItem.Title,
		QualityProfileID:  profiles[0].ID,
		MetadataProfileID: metadataProfiles[0].ID,
		// MusicBrainzID:     searchItem.Data.MusicBrainzID,
		Monitored:      true,
		Path:           "/music", // This should be a valid path in your Lidarr setup
		SearchForMedia: false,
	}

	// Add the artist
	addedArtist, err := client.AddMedia(ctx, addRequest)
	require.NoError(t, err)
	t.Logf("Added artist: %s (ID: %d)", addedArtist.Title, addedArtist.ID)

	// Update the artist
	updateRequest := requests.AutomationMediaUpdateRequest{
		Monitored: false,
	}

	updatedArtist, err := client.UpdateMedia(ctx, int64(addedArtist.ID), updateRequest)
	require.NoError(t, err)
	assert.Equal(t, addedArtist.ID, updatedArtist.ID)
	assert.False(t, updatedArtist.Monitored)
	t.Logf("Updated artist monitoring status to: %v", updatedArtist.Monitored)

	// Delete the artist
	err = client.DeleteMedia(ctx, int64(addedArtist.ID))
	require.NoError(t, err)
	t.Logf("Deleted artist with ID: %d", addedArtist.ID)
}
