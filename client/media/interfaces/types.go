package interfaces

import (
	"time"
)

// ExternalID represents an ID from an external source
type ExternalID struct {
	Source string `json:"source"` // e.g., "tmdb", "imdb", "trakt", "tvdb"
	ID     string `json:"id"`     // The actual ID
}

// ExternalID represents an ID from an external source
type Rating struct {
	Source string  `json:"source"` // e.g., "tmdb", "imdb", "trakt", "tvdb"
	Value  float32 `json:"value"`  // The actual ID
	// For sources that might have how many people voted on an item
	Votes int `json:"votes,omitempty"`
}

type Ratings []Rating

// ExternalIDs is a collection of IDs from different sources
type ExternalIDs []ExternalID

// GetID returns the ID for a specific source, empty string if not found
func (ids ExternalIDs) GetID(source string) string {
	for _, id := range ids {
		if id.Source == source {
			return id.ID
		}
	}
	return ""
}

// Add method to add/update exisiting IDs by source

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

// Artwork holds different types of artwork
type Artwork struct {
	Poster     string `json:"poster,omitempty"`
	Background string `json:"background,omitempty"`
	Banner     string `json:"banner,omitempty"`
	Thumbnail  string `json:"thumbnail,omitempty"`
	Logo       string `json:"logo,omitempty"`
}

// Person represents someone involved with the media
type Person struct {
	Name      string `json:"name"`
	Role      string `json:"role,omitempty"`      // e.g., "Director", "Actor"
	Character string `json:"character,omitempty"` // For actors
	Photo     string `json:"photo,omitempty"`
}

// MediaMetadata contains common metadata fields for all media types
type MediaMetadata struct {
	Title         string        `json:"title"`
	Description   string        `json:"description,omitempty"`
	ReleaseDate   time.Time     `json:"releaseDate,omitempty"`
	ReleaseYear   int           `json:"releaseYear,omitempty"`
	AddedAt       time.Time     `json:"addedAt,omitempty"`
	UpdatedAt     time.Time     `json:"updatedAt,omitempty"`
	Genres        []string      `json:"genres,omitempty"`
	Tags          []string      `json:"tags,omitempty"`
	Studios       []string      `json:"studios,omitempty"`
	ExternalIDs   ExternalIDs   `json:"externalIDs,omitempty"`
	ContentRating string        `json:"contentRating,omitempty"`
	Ratings       Ratings       `json:"ratings,omitempty"`
	UserRating    float32       `json:"userRating,omitempty"`
	Artwork       Artwork       `json:"artwork,omitempty"`
	Duration      time.Duration `json:"durationSeconds,omitempty"`
}

type MediaData interface {
	isMediaData()
	GetDetails() MediaMetadata
	GetMediaType() MediaType
}

func (Movie) isMediaData()      {}
func (TVShow) isMediaData()     {}
func (Episode) isMediaData()    {}
func (Track) isMediaData()      {}
func (Artist) isMediaData()     {}
func (Album) isMediaData()      {}
func (Season) isMediaData()     {}
func (Collection) isMediaData() {}
func (Playlist) isMediaData()   {}

type MediaType string

const (
	MEDIATYPE_MOVIE      MediaType = "movie"
	MEDIATYPE_SHOW       MediaType = "show"
	MEDIATYPE_SEASON     MediaType = "season"
	MEDIATYPE_EPISODE    MediaType = "episode"
	MEDIATYPE_ARTIST     MediaType = "artist"
	MEDIATYPE_ALBUM      MediaType = "album"
	MEDIATYPE_TRACK      MediaType = "track"
	MEDIATYPE_PLAYLIST   MediaType = "playlist"
	MEDIATYPE_COLLECTION MediaType = "collection"
)

// MediaItem is the base type for all media items
type MediaItem[T MediaData] struct {
	ID          uint64          `json:"ID" gorm:"primaryKey"` // internal ID
	ExternalID  string          `json:"externalID" gorm:"index"`
	ClientID    uint64          `json:"clientID"  gorm:"index"` // internal ClientID
	ClientType  MediaClientType `json:"clientType"`             // internal Client Type "plex", "jellyfin", etc.
	Type        string          `json:"type"`                   // "movie", "tvshow", "episode", "music","playlist","artist"
	StreamURL   string          `json:"streamUrl,omitempty"`
	DownloadURL string          `json:"downloadUrl,omitempty"`
	Data        T
}

// Implement this interface for MediaItem[T]
func (m *MediaItem[MediaData]) SetClientInfo(clientID uint64, clientType MediaClientType, clientItemKey string) {
	m.ClientID = clientID
	m.ClientType = clientType
	m.ExternalID = clientItemKey
}

func (m *MediaItem[MediaData]) GetData() MediaData {
	return m.Data
}

// QueryOptions provides parameters for filtering and pagination
type QueryOptions struct {
	Limit                int               `json:"limit,omitempty"`
	Offset               int               `json:"offset,omitempty"`
	Sort                 string            `json:"sort,omitempty"`
	SortOrder            string            `json:"sortOrder,omitempty"` // "asc" or "desc"
	Filters              map[string]string `json:"filters,omitempty"`
	IncludeWatchProgress bool              `json:"includeWatchProgress,omitempty"`
}
