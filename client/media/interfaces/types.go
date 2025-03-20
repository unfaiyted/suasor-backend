package interfaces

import (
	"time"
)

// ExternalID represents an ID from an external source
type ExternalID struct {
	Source string `json:"source"` // e.g., "tmdb", "imdb", "trakt", "tvdb"
	ID     string `json:"id"`     // The actual ID
}

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
	Title         string      `json:"title"`
	Description   string      `json:"description,omitempty"`
	ReleaseDate   time.Time   `json:"releaseDate,omitempty"`
	ReleaseYear   int         `json:"releaseYear,omitempty"`
	AddedAt       time.Time   `json:"addedAt,omitempty"`
	UpdatedAt     time.Time   `json:"updatedAt,omitempty"`
	Genres        []string    `json:"genres,omitempty"`
	Tags          []string    `json:"tags,omitempty"`
	Studios       []string    `json:"studios,omitempty"`
	ExternalIDs   ExternalIDs `json:"externalIds,omitempty"`
	ContentRating string      `json:"contentRating,omitempty"`
	Rating        float64     `json:"rating,omitempty"` // 0-10 scale
	UserRating    float64     `json:"userRating,omitempty"`
	Artwork       Artwork     `json:"artwork,omitempty"`
	Duration      int         `json:"durationSeconds,omitempty"`
}

// MediaItem is the base type for all media items
type MediaItem struct {
	ID          string        `json:"id"`
	Type        string        `json:"type"` // "movie", "tvshow", "episode", "music"
	ClientID    uint64        `json:"clientId"`
	ClientType  string        `json:"clientType"` // "plex", "jellyfin", etc.
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
	ID                string      `json:"id"`
	Number            int         `json:"seasonNumber"`
	Title             string      `json:"title,omitempty"`
	Overview          string      `json:"overview,omitempty"`
	EpisodeCount      int         `json:"episodeCount"`
	Artwork           Artwork     `json:"artwork,omitempty"`
	ReleaseDate       time.Time   `json:"releaseDate,omitempty"`
	ExternalParentIDs ExternalIDs `json:"externalIds,omitempty"`
}

// Episode represents a TV episode
type Episode struct {
	MediaItem
	Number            int64       `json:"number"`
	ShowID            string      `json:"showId"`
	SeasonID          string      `json:"seasonId"`
	ExternalParentIDs ExternalIDs `json:"externalIds,omitempty"`
	SeasonNumber      int         `json:"seasonNumber"`
	ShowTitle         string      `json:"showTitle,omitempty"`
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
	Albums         []string `json:"albumIds,omitempty"`
	Biography      string   `json:"biography,omitempty"`
	SimilarArtists []string `json:"similarArtists,omitempty"`
}

// MusicAlbum represents a music album
type MusicAlbum struct {
	MediaItem
	ArtistID          string      `json:"artistId"`
	ArtistName        string      `json:"artistName"`
	ExternalAlbumIDs  ExternalIDs `json:"externalAlbumIds,omitempty"`
	ExternalArtistIDs ExternalIDs `json:"externalArtistIds,omitempty"`
	TrackCount        int         `json:"trackCount"`
}

// MusicTrack represents a music track
type MusicTrack struct {
	MediaItem
	AlbumID           string      `json:"albumId"`
	ArtistID          string      `json:"artistId"`
	AlbumTitle        string      `json:"albumTitle,omitempty"`
	ExternalAlbumIDs  ExternalIDs `json:"externalAlbumIds,omitempty"`
	ExternalArtistIDs ExternalIDs `json:"externalArtistIds,omitempty"`

	ArtistName string `json:"artistName,omitempty"`
	Number     int    `json:"trackNumber,omitempty"`
	DiscNumber int    `json:"discNumber,omitempty"`
	Composer   string `json:"composer,omitempty"`
	Lyrics     string `json:"lyrics,omitempty"`
}

// Collection represents a collection of media items
type Collection struct {
	MediaItem
	ItemIDs        []string `json:"itemIds"`
	ItemCount      int      `json:"itemCount"`
	CollectionType string   `json:"collectionType"` // e.g., "movie", "tvshow"
}

// Playlist represents a user-created playlist of media items
type Playlist struct {
	MediaItem
	ItemIDs   []string `json:"itemIds"`
	ItemCount int      `json:"itemCount"`
	Owner     string   `json:"owner,omitempty"`
	IsPublic  bool     `json:"isPublic"`
}

// WatchHistoryItem represents an item in watch history
type WatchHistoryItem struct {
	ItemID          string    `json:"itemId"`
	ItemType        string    `json:"itemType"` // "movie", "episode"
	Title           string    `json:"title"`
	WatchedAt       time.Time `json:"watchedAt"`
	PositionSeconds int       `json:"positionSeconds"`
	DurationSeconds int       `json:"durationSeconds"`
	Completed       bool      `json:"completed"`
	ClientID        uint64    `json:"clientId"`
	ClientType      string    `json:"clientType"`
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
