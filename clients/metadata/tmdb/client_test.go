package tmdb

import (
	"context"
	"suasor/clients/types"
	"testing"
)

func TestNewClient(t *testing.T) {
	config := &types.TMDBConfig{
		BaseMetadataClientConfig: types.BaseMetadataClientConfig{
			BaseClientConfig: types.BaseClientConfig{
				Name:         "TMDB Test",
				BaseURL:      "https://api.themoviedb.org/3",
				Enabled:      true,
				ValidateConn: true,
			},
			SupportsMovies:      true,
			SupportsTV:          true,
			SupportsPersons:     true,
			SupportsCollections: true,
		},
		ApiKey: "test-api-key",
	}

	_, err := NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create TMDB client: %v", err)
	}
}

// Additional tests would be added here for each TMDB client method
// For example:
// func TestGetMovie(t *testing.T) { ... }
// func TestSearchMovies(t *testing.T) { ... }
// etc.

// These tests will likely use mocking to avoid making actual API calls during testing.
// For now, we're just setting up the structure.

