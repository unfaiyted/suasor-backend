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

// MediaItem is the base type for all media items
type MediaItem struct {
	ID          uint64        `json:"ID" gorm:"primaryKey"` // internal ID
	ExternalID  string        `json:"externalID" gorm:"index"`
	ClientID    uint64        `json:"clientID"  gorm:"index"` // internal ClientID
	ClientType  string        `json:"clientType"`             // internal Client Type "plex", "jellyfin", etc.
	Type        string        `json:"type"`                   // "movie", "tvshow", "episode", "music","playlist","artist"
	Metadata    MediaMetadata `json:"metadata"`
	StreamURL   string        `json:"streamUrl,omitempty"`
	DownloadURL string        `json:"downloadUrl,omitempty"`
}

// Movie represents a movie item
type Movie struct {
	MediaItem
	Cast         []Person `json:"cast,omitempty"`
	Crew         []Person `json:"crew,omitempty"`
	TrailerURL   string   `json:"trailerUrl,omitempty"`
	Resolution   string   `json:"resolution,omitempty"` // e.g., "4K", "1080p"
	VideoCodec   string   `json:"videoCodec,omitempty"`
	AudioCodec   string   `json:"audioCodec,omitempty"`
	SubtitleURLs []string `json:"subtitleUrls,omitempty"`
}

// Season represents a TV season
type Season struct {
	MediaItem
	Number       int       `json:"seasonNumber"`
	Title        string    `json:"title,omitempty"`
	Overview     string    `json:"overview,omitempty"`
	EpisodeCount int       `json:"episodeCount"`
	Artwork      Artwork   `json:"artwork,omitempty"`
	ReleaseDate  time.Time `json:"releaseDate,omitempty"`
	ParentID     string    `json:"parentID,omitempty"`
}

// Episode represents a TV episode
type Episode struct {
	MediaItem
	Number       int64  `json:"number"`
	ShowID       string `json:"showID"`
	SeasonID     string `json:"seasonID"`
	SeasonNumber int    `json:"seasonNumber"`
	ShowTitle    string `json:"showTitle,omitempty"`
}

// TVShow represents a TV series
type TVShow struct {
	MediaItem
	Seasons       []Season `json:"seasons,omitempty"`
	EpisodeCount  int      `json:"episodeCount"`
	SeasonCount   int      `json:"seasonCount"`
	ReleaseYear   int      `json:"releaseYear"`
	ContentRating string   `json:"contentRating"`
	Rating        float64  `json:"rating"`
	Network       string   `json:"network,omitempty"`
	Status        string   `json:"status,omitempty"` // e.g., "Ended", "Continuing"
	Genres        []string `json:"genres,omitempty"`
}

// MusicArtist represents a music artist
type MusicArtist struct {
	MediaItem
	Albums         []string `json:"albumIDs,omitempty"`
	Biography      string   `json:"biography,omitempty"`
	SimilarArtists []string `json:"similarArtists,omitempty"`
}

// MusicAlbum represents a music album
type MusicAlbum struct {
	MediaItem
	ArtistID   string `json:"artistID"`
	ArtistName string `json:"artistName"`
	TrackCount int    `json:"trackCount"`
}

// MusicTrack represents a music track
type MusicTrack struct {
	MediaItem
	AlbumID    string `json:"albumID"`
	ArtistID   string `json:"artistID"`
	AlbumName  string `json:"albumName"`
	AlbumTitle string `json:"albumTitle,omitempty"`

	ArtistName string `json:"artistName,omitempty"`
	Number     int    `json:"trackNumber,omitempty"`
	DiscNumber int    `json:"discNumber,omitempty"`
	Composer   string `json:"composer,omitempty"`
	Lyrics     string `json:"lyrics,omitempty"`
}

// Collection represents a collection of media items
type Collection struct {
	MediaItem
	ItemIDs        []string `json:"itemIDs"`
	ItemCount      int      `json:"itemCount"`
	CollectionType string   `json:"collectionType"` // e.g., "movie", "tvshow"
}

// Playlist represents a user-created playlist of media items
type Playlist struct {
	MediaItem
	ItemIDs   []string `json:"itemIDs"`
	ItemCount int      `json:"itemCount"`
	Owner     string   `json:"owner,omitempty"`
	IsPublic  bool     `json:"isPublic"`
}

// WatchHistoryItem represents an item in watch history
type WatchHistoryItem struct {
	MediaItem
	ItemType         string    `json:"itemType"` // "movie", "episode" , "show","season"
	WatchedAt        time.Time `json:"watchedAt"`
	LastWatchedAt    time.Time `json:"lastWatchedAt"`
	IsFavorite       bool      `json:"isFavorite,omitempty"`
	PlayedPercentage float64   `json:"playedPercentage,omitempty"`
	PlayCount        int32     `json:"playCount,omitempty"`
	PositionSeconds  int       `json:"positionSeconds"`
	DurationSeconds  int       `json:"durationSeconds"`
	Completed        bool      `json:"completed"`
	SeriesName       string    `json:"seriesName,omiempty"`
	SeasonNumber     int       `json:"seasonNumber,omitempty"`
	EpisodeNumber    int       `json:"episodeNumber,omitempty"`
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
