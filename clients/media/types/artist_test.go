package types

import (
	"testing"
)

func TestArtistMergeWithNilAlbums(t *testing.T) {
	// Create an artist with nil Albums
	artist1 := &Artist{
		Details: &MediaDetails{
			Title: "Test Artist",
		},
		// Albums is nil
	}

	// Create another artist with valid Albums
	artist2 := &Artist{
		Details: &MediaDetails{
			Title: "Test Artist Updated",
		},
		Albums: AlbumEntries{
			{
				AlbumID:   101,
				AlbumName: "Album 1",
				TrackIDs:  []uint64{1001, 1002, 1003},
			},
			{
				AlbumID:   102,
				AlbumName: "Album 2",
				TrackIDs:  []uint64{2001, 2002},
			},
		},
	}

	// Test merging artist2 into artist1 (nil Albums)
	artist1.Merge(artist2)

	// Check if the Albums was properly initialized
	if artist1.Albums == nil {
		t.Error("Albums should not be nil after merge")
	}

	// Check if the data was properly merged
	if len(artist1.Albums) != 2 {
		t.Errorf("Expected 2 albums, got %d", len(artist1.Albums))
	}

	// Check if the tracks were properly added
	album1 := artist1.GetAlbumByID(101)
	if album1 == nil {
		t.Error("Expected to find album with ID 101")
	} else if len(album1.TrackIDs) != 3 {
		t.Errorf("Expected 3 tracks in album 1, got %d", len(album1.TrackIDs))
	}

	album2 := artist1.GetAlbumByID(102)
	if album2 == nil {
		t.Error("Expected to find album with ID 102")
	} else if len(album2.TrackIDs) != 2 {
		t.Errorf("Expected 2 tracks in album 2, got %d", len(album2.TrackIDs))
	}

	// Check if the album count was updated
	if artist1.AlbumCount != 2 {
		t.Errorf("Expected album count to be 2, got %d", artist1.AlbumCount)
	}

	// Check if the track count was updated
	if artist1.TrackCount != 5 {
		t.Errorf("Expected track count to be 5, got %d", artist1.TrackCount)
	}
}

func TestArtistMergeWithNilOtherAlbums(t *testing.T) {
	// Create an artist with valid Albums
	artist1 := &Artist{
		Details: &MediaDetails{
			Title: "Test Artist",
		},
		Albums: AlbumEntries{
			{
				AlbumID:   101,
				AlbumName: "Album 1",
				TrackIDs:  []uint64{1001, 1002, 1003},
			},
			{
				AlbumID:   102,
				AlbumName: "Album 2",
				TrackIDs:  []uint64{2001, 2002},
			},
		},
		AlbumCount: 2,
		TrackCount: 5,
	}

	// Create another artist with nil Albums
	artist2 := &Artist{
		Details: &MediaDetails{
			Title: "Test Artist Updated",
		},
		// Albums is nil
	}

	// Test merging artist2 (nil Albums) into artist1
	artist1.Merge(artist2)

	// Check if the data remained unchanged
	if len(artist1.Albums) != 2 {
		t.Errorf("Expected 2 albums, got %d", len(artist1.Albums))
	}

	// Check if the tracks remained unchanged
	album1 := artist1.GetAlbumByID(101)
	if album1 == nil {
		t.Error("Expected to find album with ID 101")
	} else if len(album1.TrackIDs) != 3 {
		t.Errorf("Expected 3 tracks in album 1, got %d", len(album1.TrackIDs))
	}

	album2 := artist1.GetAlbumByID(102)
	if album2 == nil {
		t.Error("Expected to find album with ID 102")
	} else if len(album2.TrackIDs) != 2 {
		t.Errorf("Expected 2 tracks in album 2, got %d", len(album2.TrackIDs))
	}

	// Check if the album count remained unchanged
	if artist1.AlbumCount != 2 {
		t.Errorf("Expected album count to be 2, got %d", artist1.AlbumCount)
	}

	// Check if the track count remained unchanged
	if artist1.TrackCount != 5 {
		t.Errorf("Expected track count to be 5, got %d", artist1.TrackCount)
	}
}

func TestArtistWithNilAlbums(t *testing.T) {
	// Create an artist with nil Albums
	artist := &Artist{
		Details: &MediaDetails{
			Title: "Test Artist",
		},
		// Albums is nil
	}

	// Test GetOrderedAlbums
	albums := artist.GetOrderedAlbums()
	if len(albums) != 0 {
		t.Errorf("Expected empty albums list, got %v", albums)
	}

	// Test GetTrackIDsByAlbum
	trackIDs := artist.GetTrackIDsByAlbum(101)
	if len(trackIDs) != 0 {
		t.Errorf("Expected empty track IDs list, got %v", trackIDs)
	}

	// Test GetAllTrackIDs
	allTrackIDs := artist.GetAllTrackIDs()
	if len(allTrackIDs) != 0 {
		t.Errorf("Expected empty all track IDs list, got %v", allTrackIDs)
	}

	// Test CalculateTrackCount
	trackCount := artist.CalculateTrackCount()
	if trackCount != 0 {
		t.Errorf("Expected track count to be 0, got %d", trackCount)
	}
}

func TestArtistMergeBiographyAndOtherFields(t *testing.T) {
	// Create an artist with minimal details
	artist1 := &Artist{
		Details: &MediaDetails{
			Title: "Test Artist",
		},
	}

	// Create another artist with additional details
	artist2 := &Artist{
		Details: &MediaDetails{
			Title: "Test Artist",
		},
		Biography: "This is a test artist biography",
		Genres:    []string{"Rock", "Pop"},
		SimilarArtists: []ArtistReference{
			{Name: "Similar Artist 1", ID: 201},
			{Name: "Similar Artist 2", ID: 202},
		},
		StartYear: 1990,
		EndYear:   2020,
		Rating:    4.5,
		Credits: Credits{
			{Name: "Member 1", Role: "Vocals", IsCast: true},
			{Name: "Member 2", Role: "Guitar", IsCast: true},
		},
	}

	// Test merging artist2 into artist1
	artist1.Merge(artist2)

	// Check if biography was merged
	if artist1.Biography != "This is a test artist biography" {
		t.Errorf("Expected biography to be merged, got %s", artist1.Biography)
	}

	// Check if genres were merged
	if len(artist1.Genres) != 2 {
		t.Errorf("Expected 2 genres, got %d", len(artist1.Genres))
	} else if artist1.Genres[0] != "Rock" || artist1.Genres[1] != "Pop" {
		t.Errorf("Expected genres Rock and Pop, got %v", artist1.Genres)
	}

	// Check if similar artists were merged
	if len(artist1.SimilarArtists) != 2 {
		t.Errorf("Expected 2 similar artists, got %d", len(artist1.SimilarArtists))
	}

	// Check if years were merged
	if artist1.StartYear != 1990 {
		t.Errorf("Expected start year 1990, got %d", artist1.StartYear)
	}
	if artist1.EndYear != 2020 {
		t.Errorf("Expected end year 2020, got %d", artist1.EndYear)
	}

	// Check if rating was merged
	if artist1.Rating != 4.5 {
		t.Errorf("Expected rating 4.5, got %f", artist1.Rating)
	}

	// Check if credits were merged
	if len(artist1.Credits) != 2 {
		t.Errorf("Expected 2 credit entries, got %d", len(artist1.Credits))
	}
}

func TestArtistAddAlbumTrackIDs(t *testing.T) {
	// Create a new artist
	artist := &Artist{
		Details: &MediaDetails{
			Title: "Test Artist",
		},
	}

	// Add tracks to a new album
	artist.AddAlbumTrackIDs(101, "Album 1", []uint64{1001, 1002, 1003})

	// Check if the album was added
	if len(artist.Albums) != 1 {
		t.Errorf("Expected 1 album, got %d", len(artist.Albums))
	}

	// Check if the track count was updated
	if artist.TrackCount != 3 {
		t.Errorf("Expected track count to be 3, got %d", artist.TrackCount)
	}

	// Add more tracks to the same album
	artist.AddAlbumTrackIDs(101, "Album 1", []uint64{1004, 1005})

	// Check if the tracks were added without duplicates
	album := artist.GetAlbumByID(101)
	if album == nil {
		t.Error("Expected to find album with ID 101")
	} else if len(album.TrackIDs) != 5 {
		t.Errorf("Expected 5 tracks in album, got %d", len(album.TrackIDs))
	}

	// Add tracks to a new album
	artist.AddAlbumTrackIDs(102, "Album 2", []uint64{2001, 2002})

	// Check if the album count was updated
	if artist.AlbumCount != 2 {
		t.Errorf("Expected album count to be 2, got %d", artist.AlbumCount)
	}

	// Check if the track count was updated
	if artist.TrackCount != 7 {
		t.Errorf("Expected track count to be 7, got %d", artist.TrackCount)
	}
}

func TestArtistMergeTrackIDsByAlbum(t *testing.T) {
	// Create a new artist
	artist := &Artist{
		Details: &MediaDetails{
			Title: "Test Artist",
		},
		Albums: AlbumEntries{
			{
				AlbumID:   101,
				AlbumName: "Album 1",
				TrackIDs:  []uint64{1001, 1002},
			},
		},
	}

	// Merge track IDs into existing album
	artist.MergeTrackIDsByAlbum(101, "Album 1", []uint64{1003, 1004})

	// Check if tracks were merged
	album := artist.GetAlbumByID(101)
	if album == nil {
		t.Error("Expected to find album with ID 101")
	} else if len(album.TrackIDs) != 4 {
		t.Errorf("Expected 4 tracks in album, got %d", len(album.TrackIDs))
	}

	// Merge track IDs into new album
	artist.MergeTrackIDsByAlbum(102, "Album 2", []uint64{2001, 2002})

	// Check if new album was created
	if len(artist.Albums) != 2 {
		t.Errorf("Expected 2 albums, got %d", len(artist.Albums))
	}

	// Check if track counts were updated
	if artist.TrackCount != 6 {
		t.Errorf("Expected track count to be 6, got %d", artist.TrackCount)
	}
}