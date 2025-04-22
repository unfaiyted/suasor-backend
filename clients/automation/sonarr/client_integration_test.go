// client_integration_test.go
package sonarr

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	config "suasor/clients/types"

	"github.com/joho/godotenv"
	"suasor/clients/automation/providers"
	"suasor/clients/automation/types"
	"suasor/types/requests"
	logger "suasor/utils/logger"
)

func init() {
	// Try to load .env file from several possible locations
	locations := []string{
		".env",          // Current directory
		"../../../.env", // Project root
		filepath.Join(os.Getenv("HOME"), "sonarr_test.env"), // Home directory
	}

	for _, location := range locations {
		err := godotenv.Load(location)
		if err == nil {
			fmt.Printf("Loaded environment from: %s\n", location)
			break
		}
	}
}

// Integration test for SonarrClient
// To run these tests:
// INTEGRATION=true SONARR_TEST_URL=http://your-server:8989 SONARR_TEST_API_KEY=your-api-key go test -v -tags=integration

func TestSonarrClientIntegration(t *testing.T) {
	// Skip if not running integration tests or missing environment variables
	if os.Getenv("INTEGRATION") != "true" {
		t.Skip("Skipping integration test. Set INTEGRATION=true to run")
	}

	// Get test credentials from environment
	baseURL := os.Getenv("SONARR_TEST_URL")
	apiKey := os.Getenv("SONARR_TEST_API_KEY")

	if baseURL == "" || apiKey == "" {
		t.Fatal("Missing required environment variables for integration test")
	}

	// Create client configuration
	config := config.SonarrConfig{
		BaseURL: baseURL,
		APIKey:  apiKey,
	}

	logger.Initialize()
	ctx := context.Background()

	// Initialize client
	client, err := NewSonarrClient(ctx, 1, config)
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

// Test getting system status from Sonarr
func testGetSystemStatus(t *testing.T, ctx context.Context, client providers.AutomationProvider) {
	status, err := client.GetSystemStatus(ctx)
	require.NoError(t, err)

	// Validate results
	assert.NotEmpty(t, status.Version, "Expected to get a version string")
	t.Logf("Sonarr version: %s", status.Version)
	assert.NotEmpty(t, status.OsName, "Expected to get an OS name")
}

// Test getting library items from Sonarr
func testGetLibraryItems(t *testing.T, ctx context.Context, client providers.AutomationProvider) {
	// Get library items with limit
	options := &types.LibraryQueryOptions{
		Limit: 10,
	}

	series, err := client.GetLibraryItems(ctx, options)
	require.NoError(t, err)

	// Log the count even if empty
	t.Logf("Found %d TV shows in library", len(series))

	if len(series) > 0 {
		show := series[0]
		t.Logf("Got TV show: %s (ID: %d)", show.Title, show.ID)

		// Verify show has expected fields
		assert.NotZero(t, show.ID)
		assert.NotEmpty(t, show.Title)
		assert.NotEmpty(t, show.Path, "Expected TV show to have a path")
	}
}

// Test getting and searching for media
func testGetAndSearchMedia(t *testing.T, ctx context.Context, client providers.AutomationProvider) {
	// First, check if we have any TV shows in the library
	series, err := client.GetLibraryItems(ctx, &types.LibraryQueryOptions{Limit: 1})
	require.NoError(t, err)

	if len(series) > 0 {
		// Test getting a specific TV show by ID
		seriesID := int64(series[0].ID)
		show, err := client.GetMediaByID(ctx, seriesID)
		require.NoError(t, err)

		// Validate the result
		assert.Equal(t, seriesID, int64(show.ID))
		assert.NotEmpty(t, show.Title)
		assert.NotZero(t, show.Year)
	} else {
		t.Log("No TV shows in library to test GetMediaByID")
	}

	// Test search functionality
	searchTerm := "Breaking Bad" // A commonly available TV show
	searchResults, err := client.SearchMedia(ctx, searchTerm, nil)
	require.NoError(t, err)

	t.Logf("Search for '%s' returned %d results", searchTerm, len(searchResults))
	if len(searchResults) > 0 {
		result := searchResults[0]
		assert.NotEmpty(t, result.Title)
		assert.True(t, result.Year > 0, "Expected year to be set")
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
	command := types.Command{
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

	// Search for a TV show to add
	searchResults, err := client.SearchMedia(ctx, "New Girl", nil)
	require.NoError(t, err)
	require.NotEmpty(t, searchResults, "Need at least one search result")

	// Prepare add request
	searchItem := searchResults[0]
	addRequest := requests.AutomationMediaAddRequest{
		Title:            searchItem.Title,
		Year:             int(searchItem.Year),
		QualityProfileID: profiles[0].ID,
		TVDBID:           int64(searchItem.ID), // This needs to be the actual TVDB ID
		Monitored:        true,
		Path:             "/tv", // This should be a valid path in your Sonarr setup
	}

	// Add the TV show
	addedShow, err := client.AddMedia(ctx, addRequest)
	require.NoError(t, err)
	t.Logf("Added TV show: %s (ID: %d)", addedShow.Title, addedShow.ID)

	// Update the TV show
	updateRequest := requests.AutomationMediaUpdateRequest{
		Monitored: false,
	}

	updatedShow, err := client.UpdateMedia(ctx, int64(addedShow.ID), updateRequest)
	require.NoError(t, err)
	assert.Equal(t, addedShow.ID, updatedShow.ID)
	assert.False(t, updatedShow.Monitored)
	t.Logf("Updated TV show monitoring status to: %v", updatedShow.Monitored)

	// Delete the TV show
	err = client.DeleteMedia(ctx, int64(addedShow.ID))
	require.NoError(t, err)
	t.Logf("Deleted TV show with ID: %d", addedShow.ID)
}
