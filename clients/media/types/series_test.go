package types

import (
	"testing"
)

func TestSeriesMergeWithNilSeasons(t *testing.T) {
	// Create a series with nil Seasons
	series1 := &Series{
		Details: &MediaDetails{
			Title:       "Test Series",
			ReleaseYear: 2020,
		},
		// Seasons is nil
	}

	// Create another series with valid Seasons
	series2 := &Series{
		Details: &MediaDetails{
			Title:       "Test Series Updated",
			ReleaseYear: 2020,
		},
		Seasons: SeasonEntries{
			{
				SeasonNumber: 1,
				SeasonID:     1,
				EpisodeIDs:   []uint64{101, 102, 103},
			},
			{
				SeasonNumber: 2,
				SeasonID:     2,
				EpisodeIDs:   []uint64{201, 202},
			},
		},
	}

	// Test merging series2 into series1 (nil Seasons)
	series1.Merge(series2)

	// Check if the Seasons was properly initialized
	if series1.Seasons == nil {
		t.Error("Seasons should not be nil after merge")
	}

	// Check if the data was properly merged
	if len(series1.Seasons) != 2 {
		t.Errorf("Expected 2 seasons, got %d", len(series1.Seasons))
	}

	// Check if the episodes were properly added
	season1 := series1.GetSeasonByNumber(1)
	if season1 == nil {
		t.Error("Expected to find season 1")
	} else if len(season1.EpisodeIDs) != 3 {
		t.Errorf("Expected 3 episodes in season 1, got %d", len(season1.EpisodeIDs))
	}

	season2 := series1.GetSeasonByNumber(2)
	if season2 == nil {
		t.Error("Expected to find season 2")
	} else if len(season2.EpisodeIDs) != 2 {
		t.Errorf("Expected 2 episodes in season 2, got %d", len(season2.EpisodeIDs))
	}

	// Check if the season count was updated
	if series1.SeasonCount != 2 {
		t.Errorf("Expected season count to be 2, got %d", series1.SeasonCount)
	}

	// Check if the episode count was updated
	if series1.EpisodeCount != 5 {
		t.Errorf("Expected episode count to be 5, got %d", series1.EpisodeCount)
	}
}

func TestSeriesMergeWithNilOtherSeasons(t *testing.T) {
	// Create a series with valid Seasons
	series1 := &Series{
		Details: &MediaDetails{
			Title:       "Test Series",
			ReleaseYear: 2020,
		},
		Seasons: SeasonEntries{
			{
				SeasonNumber: 1,
				SeasonID:     1,
				EpisodeIDs:   []uint64{101, 102, 103},
			},
			{
				SeasonNumber: 2,
				SeasonID:     2,
				EpisodeIDs:   []uint64{201, 202},
			},
		},
		EpisodeCount: 5,
		SeasonCount:  2,
	}

	// Create another series with nil Seasons
	series2 := &Series{
		Details: &MediaDetails{
			Title:       "Test Series Updated",
			ReleaseYear: 2020,
		},
		// Seasons is nil
	}

	// Test merging series2 (nil Seasons) into series1
	series1.Merge(series2)

	// Check if the data remained unchanged
	if len(series1.Seasons) != 2 {
		t.Errorf("Expected 2 seasons, got %d", len(series1.Seasons))
	}

	// Check if the episodes remained unchanged
	season1 := series1.GetSeasonByNumber(1)
	if season1 == nil {
		t.Error("Expected to find season 1")
	} else if len(season1.EpisodeIDs) != 3 {
		t.Errorf("Expected 3 episodes in season 1, got %d", len(season1.EpisodeIDs))
	}

	season2 := series1.GetSeasonByNumber(2)
	if season2 == nil {
		t.Error("Expected to find season 2")
	} else if len(season2.EpisodeIDs) != 2 {
		t.Errorf("Expected 2 episodes in season 2, got %d", len(season2.EpisodeIDs))
	}

	// Check if the season count remained unchanged
	if series1.SeasonCount != 2 {
		t.Errorf("Expected season count to be 2, got %d", series1.SeasonCount)
	}

	// Check if the episode count remained unchanged
	if series1.EpisodeCount != 5 {
		t.Errorf("Expected episode count to be 5, got %d", series1.EpisodeCount)
	}
}

func TestSeriesWithNilSeasons(t *testing.T) {
	// Create a series with nil Seasons
	series := &Series{
		Details: &MediaDetails{
			Title:       "Test Series",
			ReleaseYear: 2020,
		},
		// Seasons is nil
	}

	// Test GetOrderedSeasons
	seasons := series.GetOrderedSeasons()
	if len(seasons) != 0 {
		t.Errorf("Expected empty seasons list, got %v", seasons)
	}

	// Test GetEpisodeIDsBySeason
	episodeIDs := series.GetEpisodeIDsBySeason(1)
	if len(episodeIDs) != 0 {
		t.Errorf("Expected empty episode IDs list, got %v", episodeIDs)
	}

	// Test GetAllEpisodeIDs
	allEpisodeIDs := series.GetAllEpisodeIDs()
	if len(allEpisodeIDs) != 0 {
		t.Errorf("Expected empty all episode IDs list, got %v", allEpisodeIDs)
	}

	// Test CalculateEpisodeCount
	episodeCount := series.CalculateEpisodeCount()
	if episodeCount != 0 {
		t.Errorf("Expected episode count to be 0, got %d", episodeCount)
	}
}