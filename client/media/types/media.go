package types

import (
	"time"
)

// MusicArtist represents a music artist
type Artist struct {
	Details        MediaDetails
	Albums         []string `json:"albumIDs,omitempty"`
	AlbumCount     int      `json:"albumCount"`
	Biography      string   `json:"biography,omitempty"`
	SimilarArtists []string `json:"similarArtists,omitempty"`
}

// MusicAlbum represents a music album
type Album struct {
	Details    MediaDetails
	ArtistID   string `json:"artistID"`
	ArtistName string `json:"artistName"`
	TrackCount int    `json:"trackCount"`
}

// MusicTrack represents a music track
type Track struct {
	Details    MediaDetails
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

// Season represents a TV season
type Season struct {
	Details      MediaDetails
	Number       int       `json:"seasonNumber"`
	Title        string    `json:"title,omitempty"`
	Overview     string    `json:"overview,omitempty"`
	EpisodeCount int       `json:"episodeCount"`
	Artwork      Artwork   `json:"artwork,omitempty"`
	ReleaseDate  time.Time `json:"releaseDate,omitempty"`
	SeriesName   string    `json:"seriesName,omitempty"`
	SeriesID     string    `json:"seriesID"`
}

// Episode represents a TV episode
type Episode struct {
	Details      MediaDetails
	Number       int64  `json:"number"`
	ShowID       string `json:"showID"`
	SeasonID     string `json:"seasonID"`
	SeasonNumber int    `json:"seasonNumber"`
	ShowTitle    string `json:"showTitle,omitempty"`
}

// Series represents a TV series
type Series struct {
	Details       MediaDetails
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

// Collection represents a collection of media items
type Collection struct {
	Details        MediaDetails
	ItemIDs        []string `json:"itemIDs"`
	ItemCount      int      `json:"itemCount"`
	CollectionType string   `json:"collectionType"` // e.g., "movie", "tvshow"
}

// Playlist represents a user-created playlist of media items
type Playlist struct {
	Details   MediaDetails
	ItemIDs   []string `json:"itemIDs"`
	ItemCount int      `json:"itemCount"`
	Owner     string   `json:"owner,omitempty"`
	IsPublic  bool     `json:"isPublic"`
}

// Movie represents a movie item
type Movie struct {
	Details      MediaDetails
	Cast         []Person `json:"cast,omitempty"`
	Crew         []Person `json:"crew,omitempty"`
	TrailerURL   string   `json:"trailerUrl,omitempty"`
	Resolution   string   `json:"resolution,omitempty"` // e.g., "4K", "1080p"
	VideoCodec   string   `json:"videoCodec,omitempty"`
	AudioCodec   string   `json:"audioCodec,omitempty"`
	SubtitleURLs []string `json:"subtitleUrls,omitempty"`
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

type MediaData interface {
	isMediaData()
	GetDetails() MediaDetails
	GetMediaType() MediaType
}

func (Movie) isMediaData()      {}
func (Series) isMediaData()     {}
func (Episode) isMediaData()    {}
func (Track) isMediaData()      {}
func (Artist) isMediaData()     {}
func (Album) isMediaData()      {}
func (Season) isMediaData()     {}
func (Collection) isMediaData() {}
func (Playlist) isMediaData()   {}

func (t Track) GetDetails() MediaDetails { return t.Details }
func (t Track) GetMediaType() MediaType  { return MEDIATYPE_TRACK }

func (a Album) GetDetails() MediaDetails { return a.Details }
func (a Album) GetMediaType() MediaType  { return MEDIATYPE_ALBUM }

func (a Artist) GetDetails() MediaDetails { return a.Details }
func (a Artist) GetMediaType() MediaType  { return MEDIATYPE_ARTIST }

// Then in each media type
func (m Movie) GetDetails() MediaDetails { return m.Details }
func (m Movie) GetMediaType() MediaType  { return MEDIATYPE_MOVIE }

func (c Collection) GetDetails() MediaDetails { return c.Details }
func (c Collection) GetMediaType() MediaType  { return MEDIATYPE_COLLECTION }

func (p Playlist) GetDetails() MediaDetails { return p.Details }
func (p Playlist) GetMediaType() MediaType  { return MEDIATYPE_PLAYLIST }

func (t Series) GetDetails() MediaDetails { return t.Details }
func (t Series) GetMediaType() MediaType  { return MEDIATYPE_SHOW }

func (s Season) GetDetails() MediaDetails { return s.Details }
func (s Season) GetMediaType() MediaType  { return MEDIATYPE_SEASON }

func (e Episode) GetDetails() MediaDetails { return e.Details }
func (e Episode) GetMediaType() MediaType  { return MEDIATYPE_EPISODE }
