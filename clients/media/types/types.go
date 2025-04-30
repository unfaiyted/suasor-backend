package types

import (
	"reflect"
	"time"
)

// ExternalID represents an ID from an external source
type Rating struct {
	Source string  `json:"source"` // e.g., "tmdb", "imdb", "trakt", "tvdb"
	Value  float32 `json:"value"`  // The actual ID
	// For sources that might have how many people voted on an item
	Votes int `json:"votes,omitempty"`
}

type Ratings []Rating

func (ratings Ratings) GetRating(source string) float32 {
	for _, rating := range ratings {
		if rating.Source == source {
			return rating.Value
		}
	}
	return 0
}

func (ratings Ratings) GetRatingVotes(source string) int {
	for _, rating := range ratings {
		if rating.Source == source {
			return rating.Votes
		}
	}
	return 0
}

// AddOrUpdate adds a new ExternalID or updates an existing one with the same source
func (ids *ExternalIDs) AddOrUpdate(source string, id string) {
	for i, extID := range *ids {
		if extID.Source == source {
			// Update existing ID
			(*ids)[i].ID = id
			return
		}
	}
	// Add new ID if not found
	*ids = append(*ids, ExternalID{Source: source, ID: id})
}

// AddOrUpdateRating adds a new Rating or updates an existing one with the same source
func (ratings *Ratings) AddOrUpdateRating(source string, value float32, votes int) {
	for i, rating := range *ratings {
		if rating.Source == source {
			// Update existing rating
			(*ratings)[i].Value = value
			(*ratings)[i].Votes = votes
			return
		}
	}
	// Add new rating if not found
	*ratings = append(*ratings, Rating{Source: source, Value: value, Votes: votes})
}

// MediaDetails contains common metadata fields for all media types
type MediaDetails struct {
	Title         string      `json:"title"`
	Description   string      `json:"description,omitempty"`
	ReleaseDate   time.Time   `json:"releaseDate,omitempty"`
	ReleaseYear   int         `json:"releaseYear,omitempty"`
	AddedAt       time.Time   `json:"addedAt,omitempty"`
	UpdatedAt     time.Time   `json:"updatedAt,omitempty"`
	Genres        []string    `json:"genres,omitempty"`
	Tags          []string    `json:"tags,omitempty"`
	Studio        string      `json:"studio,omitempty"`
	ExternalIDs   ExternalIDs `json:"externalIDs,omitempty"`
	ContentRating string      `json:"contentRating,omitempty"`
	Language      string      `json:"language,omitempty"`
	Ratings       Ratings     `json:"ratings,omitempty"`
	UserRating    float32     `json:"userRating,omitempty"`
	Artwork       Artwork     `json:"artwork,omitempty"`
	Duration      int64       `json:"durationSeconds,omitempty"` // Changed from time.Duration to int64 for Swagger compatibility
	IsFavorite    bool        `json:"isFavorite,omitempty"`
}

type MediaType string

const (
	MediaTypeMovie      MediaType = "movie"
	MediaTypeSeries     MediaType = "series"
	MediaTypeSeason     MediaType = "season"
	MediaTypeEpisode    MediaType = "episode"
	MediaTypeArtist     MediaType = "artist"
	MediaTypeAlbum      MediaType = "album"
	MediaTypeTrack      MediaType = "track"
	MediaTypePlaylist   MediaType = "playlist"
	MediaTypeCollection MediaType = "collection"
	MediaTypeAll        MediaType = "all"
	MediaTypeUnknown    MediaType = "unknown"
)

func GetMediaTypeFromTypeName(ofType any) MediaType {
	// First try direct type assertion which is more reliable
	switch ofType.(type) {
	case *Movie:
		return MediaTypeMovie
	case *Series:
		return MediaTypeSeries
	case *Season:
		return MediaTypeSeason
	case *Episode:
		return MediaTypeEpisode
	case *Artist:
		return MediaTypeArtist
	case *Album:
		return MediaTypeAlbum
	case *Track:
		return MediaTypeTrack
	case *Playlist:
		return MediaTypePlaylist
	case *Collection:
		return MediaTypeCollection
	}

	// If type assertion fails, try to extract type name
	t := reflect.TypeOf(ofType)
	if t == nil {
		return MediaTypeUnknown
	}

	// Normalize the type name by removing package paths and pointer symbols
	typeName := t.String()
	// Handle pointer types
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		typeName = t.Name()
	} else {
		typeName = t.Name()
	}

	// Now match on just the type name, not the full path
	switch typeName {
	case "Movie":
		return MediaTypeMovie
	case "Series":
		return MediaTypeSeries
	case "Season":
		return MediaTypeSeason
	case "Episode":
		return MediaTypeEpisode
	case "Artist":
		return MediaTypeArtist
	case "Album":
		return MediaTypeAlbum
	case "Track":
		return MediaTypeTrack
	case "Playlist":
		return MediaTypePlaylist
	case "Collection":
		return MediaTypeCollection
	default:
		return MediaTypeUnknown
	}
}
