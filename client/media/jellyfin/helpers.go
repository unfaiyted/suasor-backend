package jellyfin

import (
	jellyfin "github.com/sj14/jellyfin-go/api"
	t "suasor/client/media/types"
	"time"
)

// Helper function to get int value from pointer with default 0 if nil
func getInt32Value(ptr *int32) int {
	if ptr == nil {
		return 0
	}
	return int(*ptr)
}

// Helper function to get duration from ticks pointer
func getDurationFromTicks(ticks *int64) time.Duration {
	if ticks == nil {
		return 0
	}
	return time.Duration(*ticks/10000000) * time.Second
}

// Helper to get a string value from a pointer, with a default empty string if nil
func getStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

// Convert a single string to a string slice
func stringToSlice(s string) []string {
	if s == "" {
		return []string{}
	}
	return []string{s}
}

// Safely convert BaseItemKind to string
func baseItemKindToString(kind jellyfin.BaseItemKind) string {
	return string(kind)
}

// extractProviderIDs adds external IDs from the Jellyfin provider IDs map to the metadata
func extractProviderIDs(providerIds *map[string]string, externalIDs *t.ExternalIDs) {
	if providerIds == nil {
		return
	}

	// Common media identifier mappings
	idMappings := map[string]string{
		"Imdb":              "imdb",
		"Tmdb":              "tmdb",
		"Tvdb":              "tvdb",
		"MusicBrainzTrack":  "musicbrainz",
		"MusicBrainzAlbum":  "musicbrainz",
		"MusicBrainzArtist": "musicbrainz",
	}

	// Extract all available IDs based on the mappings
	for jellyfinKey, externalKey := range idMappings {
		if id, ok := (*providerIds)[jellyfinKey]; ok {
			externalIDs.AddOrUpdate(externalKey, id)
		}
	}

}

// Helper function for min value
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
