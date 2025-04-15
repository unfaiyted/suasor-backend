package types

import (
	"time"
)

// MusicArtist represents a music artist
type Artist struct {
	Details        MediaDetails
	Albums         []*Album `json:"albums,omitempty"`
	AlbumIDs       []uint64 `json:"albumIDs,omitempty"`
	AlbumCount     int      `json:"albumCount"`
	Biography      string   `json:"biography,omitempty"`
	SimilarArtists []string `json:"similarArtists,omitempty"`
}

type SyncClient struct {
	// ID of the client that this external ID belongs to (optional for service IDs like TMDB)
	ID uint64 `json:"clientId,omitempty"`
	// The actual ID value in the external system
	ItemID string `json:"itemId"`
}

type SyncClients []SyncClient

func (s SyncClients) AddClient(clientID uint64, itemID string) {
	// check if client ID already exists
	found := false
	for i, cID := range s {
		if cID.ID == clientID {
			// Update existing ID
			s[i].ItemID = itemID
			found = true
			break
		}
	}
	if !found {
		// Add new ID if not found
		s = append(s, SyncClient{
			ID:     clientID,
			ItemID: itemID,
		})
	}
}

func (s SyncClients) GetClientItemID(clientID uint64) string {
	for _, cID := range s {
		if cID.ID == clientID {
			return cID.ItemID
		}
	}
	return ""
}

// Season represents a TV season
type Season struct {
	Details      MediaDetails
	Number       int         `json:"seasonNumber"`
	Title        string      `json:"title,omitempty"`
	Overview     string      `json:"overview,omitempty"`
	EpisodeCount int         `json:"episodeCount"`
	Episodes     []*Episode  `json:"episodes,omitempty"`
	EpisodeIDs   []uint64    `json:"episodeIDs,omitempty"`
	Artwork      Artwork     `json:"artwork,omitempty"`
	ReleaseDate  time.Time   `json:"releaseDate,omitempty"`
	SeriesName   string      `json:"seriesName,omitempty"`
	SeriesID     uint64      `json:"seriesID"`
	SyncSeries   SyncClients `json:"syncSeries,omitempty"`
	Credits      Credits     `json:"credits,omitempty"`
}

// Series represents a TV series
type Series struct {
	Details       MediaDetails
	Seasons       []*Season `json:"seasons,omitempty"`
	EpisodeCount  int       `json:"episodeCount"`
	SeasonCount   int       `json:"seasonCount"`
	ReleaseYear   int       `json:"releaseYear"`
	ContentRating string    `json:"contentRating"`
	Rating        float64   `json:"rating"`
	Network       string    `json:"network,omitempty"`
	Status        string    `json:"status,omitempty"` // e.g., "Ended", "Continuing"
	Genres        []string  `json:"genres,omitempty"`
	Credits       Credits   `json:"credits,omitempty"`
}

// Movie represents a movie item
type Movie struct {
	Details      MediaDetails
	Credits      Credits  `json:"credits,omitempty"`
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

	IsCast    bool `json:"isCast,omitempty"`
	IsCrew    bool `json:"isCrew,omitempty"`
	IsGuest   bool `json:"isGuest,omitempty"`
	IsCreator bool `json:"isCreator,omitempty"`
	IsArtist  bool `json:"isArtist,omitempty"`
}

type Credits []Person

func (c Credits) GetCast() []Person {
	var cast []Person
	for _, person := range c {
		if person.IsCast {
			cast = append(cast, person)
		}
	}
	return cast
}

func (c Credits) GetCrew() []Person {
	var crew []Person
	for _, person := range c {
		if person.IsCrew {
			crew = append(crew, person)
		}
	}
	return crew
}

func (c Credits) GetGuests() []Person {
	var guests []Person
	for _, person := range c {
		if person.IsGuest {
			guests = append(guests, person)
		}
	}
	return guests
}

func (c Credits) GetCreators() []Person {
	var creators []Person
	for _, person := range c {
		if person.IsCreator {
			creators = append(creators, person)
		}
	}
	return creators
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
func (t Track) GetMediaType() MediaType  { return MediaTypeTrack }

func (a Album) GetDetails() MediaDetails { return a.Details }
func (a Album) GetMediaType() MediaType  { return MediaTypeAlbum }

func (a Artist) GetDetails() MediaDetails { return a.Details }
func (a Artist) GetMediaType() MediaType  { return MediaTypeArtist }

// Then in each media type
func (m Movie) GetDetails() MediaDetails { return m.Details }
func (m Movie) GetMediaType() MediaType  { return MediaTypeMovie }

func (c Collection) GetDetails() MediaDetails { return c.Details }
func (c Collection) GetMediaType() MediaType  { return MediaTypeCollection }

func (p Playlist) GetDetails() MediaDetails { return p.Details }
func (p Playlist) GetMediaType() MediaType  { return MediaTypePlaylist }

func (t Series) GetDetails() MediaDetails { return t.Details }
func (t Series) GetMediaType() MediaType  { return MediaTypeSeries }

func (s Season) GetDetails() MediaDetails { return s.Details }
func (s Season) GetMediaType() MediaType  { return MediaTypeSeason }

func (e Episode) GetDetails() MediaDetails { return e.Details }
func (e Episode) GetMediaType() MediaType  { return MediaTypeEpisode }
