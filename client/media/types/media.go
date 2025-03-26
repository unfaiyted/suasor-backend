package types

import (
	"time"
)

// MusicArtist represents a music artist
type Artist struct {
	Details        MediaMetadata
	Albums         []string `json:"albumIDs,omitempty"`
	AlbumCount     int      `json:"albumCount"`
	Biography      string   `json:"biography,omitempty"`
	SimilarArtists []string `json:"similarArtists,omitempty"`
}

// MusicAlbum represents a music album
type Album struct {
	Details    MediaMetadata
	ArtistID   string `json:"artistID"`
	ArtistName string `json:"artistName"`
	TrackCount int    `json:"trackCount"`
}

// MusicTrack represents a music track
type Track struct {
	Details    MediaMetadata
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
	Details      MediaMetadata
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
	Details      MediaMetadata
	Number       int64  `json:"number"`
	ShowID       string `json:"showID"`
	SeasonID     string `json:"seasonID"`
	SeasonNumber int    `json:"seasonNumber"`
	ShowTitle    string `json:"showTitle,omitempty"`
}

// TVShow represents a TV series
type TVShow struct {
	Details       MediaMetadata
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
	Details        MediaMetadata
	ItemIDs        []string `json:"itemIDs"`
	ItemCount      int      `json:"itemCount"`
	CollectionType string   `json:"collectionType"` // e.g., "movie", "tvshow"
}

// Playlist represents a user-created playlist of media items
type Playlist struct {
	Details   MediaMetadata
	ItemIDs   []string `json:"itemIDs"`
	ItemCount int      `json:"itemCount"`
	Owner     string   `json:"owner,omitempty"`
	IsPublic  bool     `json:"isPublic"`
}

// HistoryItem represents an item in watch history
type MediaPlayHistory[T MediaData] struct {
	Item             MediaItem[T]
	Type             string    `json:"type"` // "movie", "episode" , "show","season"
	WatchedAt        time.Time `json:"watchedAt"`
	LastWatchedAt    time.Time `json:"lastWatchedAt"`
	IsFavorite       bool      `json:"isFavorite,omitempty"`
	PlayedPercentage float64   `json:"playedPercentage,omitempty"`
	PlayCount        int32     `json:"playCount,omitempty"`
	PositionSeconds  int       `json:"positionSeconds"`
	DurationSeconds  int       `json:"durationSeconds"`
	Completed        bool      `json:"completed"`
}

// Movie represents a movie item
type Movie struct {
	Details      MediaMetadata
	Cast         []Person `json:"cast,omitempty"`
	Crew         []Person `json:"crew,omitempty"`
	TrailerURL   string   `json:"trailerUrl,omitempty"`
	Resolution   string   `json:"resolution,omitempty"` // e.g., "4K", "1080p"
	VideoCodec   string   `json:"videoCodec,omitempty"`
	AudioCodec   string   `json:"audioCodec,omitempty"`
	SubtitleURLs []string `json:"subtitleUrls,omitempty"`
}

// MediaItem is the base type for all media items
type MediaItem[T MediaData] struct {
	ID          uint64          `json:"ID" gorm:"primaryKey"` // internal ID
	ExternalID  string          `json:"externalID" gorm:"index"`
	ClientID    uint64          `json:"clientID"  gorm:"index"` // internal ClientID
	ClientType  MediaClientType `json:"clientType"`             // internal Client Type "plex", "jellyfin", etc.
	Type        MediaType       `json:"type"`                   // "movie", "tvshow", "episode", "music","playlist","artist"
	StreamURL   string          `json:"streamUrl,omitempty"`
	DownloadURL string          `json:"downloadUrl,omitempty"`
	Data        T
}

func (m *MediaItem[T]) SetData(i *MediaItem[T], data T) {
	i.Data = data
}

func (m *MediaItem[T]) AsEpisode() (MediaItem[Episode], bool) {
	if m.Type != MEDIATYPE_EPISODE {
		return MediaItem[Episode]{}, false
	}
	episode, ok := any(m).(MediaItem[Episode])

	return episode, ok
}

func (m *MediaItem[T]) AsMovie() (MediaItem[Movie], bool) {
	if m.Type != MEDIATYPE_MOVIE {
		return MediaItem[Movie]{}, false
	}
	movie, ok := any(m).(MediaItem[Movie])

	return movie, ok
}

func (m *MediaItem[T]) AsTVShow() (MediaItem[TVShow], bool) {
	if m.Type != MEDIATYPE_SHOW {
		return MediaItem[TVShow]{}, false
	}
	show, ok := any(m).(MediaItem[TVShow])

	return show, ok
}

func (m *MediaItem[T]) AsSeason() (MediaItem[Season], bool) {
	if m.Type != MEDIATYPE_SEASON {
		return MediaItem[Season]{}, false
	}
	season, ok := any(m).(MediaItem[Season])

	return season, ok
}

func (m *MediaItem[T]) AsTrack() (MediaItem[Track], bool) {
	if m.Type != MEDIATYPE_TRACK {
		return MediaItem[Track]{}, false
	}
	track, ok := any(m).(MediaItem[Track])

	return track, ok
}

func (m *MediaItem[T]) AsAlbum() (MediaItem[Album], bool) {
	if m.Type != MEDIATYPE_ALBUM {
		return MediaItem[Album]{}, false
	}
	album, ok := any(m).(MediaItem[Album])

	return album, ok
}

func (m *MediaItem[T]) AsArtist() (MediaItem[Artist], bool) {
	if m.Type != MEDIATYPE_ARTIST {
		return MediaItem[Artist]{}, false
	}
	artist, ok := any(m).(MediaItem[Artist])

	return artist, ok
}

func (m *MediaItem[T]) AsCollection() (MediaItem[Collection], bool) {
	if m.Type != MEDIATYPE_COLLECTION {
		return MediaItem[Collection]{}, false
	}
	collection, ok := any(m).(MediaItem[Collection])

	return collection, ok
}

func (m *MediaItem[T]) AsPlaylist() (MediaItem[Playlist], bool) {
	if m.Type != MEDIATYPE_PLAYLIST {
		return MediaItem[Playlist]{}, false
	}
	playlist, ok := any(m).(MediaItem[Playlist])

	return playlist, ok
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

// Implement this interface for MediaItem[T]
func (m *MediaItem[MediaData]) SetClientInfo(clientID uint64, clientType MediaClientType, clientItemKey string) {
	m.ClientID = clientID
	m.ClientType = clientType
	m.ExternalID = clientItemKey
}

func (m *MediaItem[MediaData]) GetData() MediaData {
	return m.Data
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

func (t Track) GetDetails() MediaMetadata { return t.Details }
func (t Track) GetMediaType() MediaType   { return MEDIATYPE_TRACK }

func (a Album) GetDetails() MediaMetadata { return a.Details }
func (a Album) GetMediaType() MediaType   { return MEDIATYPE_ALBUM }

func (a Artist) GetDetails() MediaMetadata { return a.Details }
func (a Artist) GetMediaType() MediaType   { return MEDIATYPE_ARTIST }

// Then in each media type
func (m Movie) GetDetails() MediaMetadata { return m.Details }
func (m Movie) GetMediaType() MediaType   { return MEDIATYPE_MOVIE }

func (c Collection) GetDetails() MediaMetadata { return c.Details }
func (c Collection) GetMediaType() MediaType   { return MEDIATYPE_COLLECTION }

func (p Playlist) GetDetails() MediaMetadata { return p.Details }
func (p Playlist) GetMediaType() MediaType   { return MEDIATYPE_PLAYLIST }

func (t TVShow) GetDetails() MediaMetadata { return t.Details }
func (t TVShow) GetMediaType() MediaType   { return MEDIATYPE_SHOW }

func (s Season) GetDetails() MediaMetadata { return s.Details }
func (s Season) GetMediaType() MediaType   { return MEDIATYPE_SEASON }

func (e Episode) GetDetails() MediaMetadata { return e.Details }
func (e Episode) GetMediaType() MediaType   { return MEDIATYPE_EPISODE }
