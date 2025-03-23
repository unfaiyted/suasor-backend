package interfaces

import (
	"context"
)

// MusicArtist represents a music artist
type Artist struct {
	Details        MediaMetadata
	Albums         []string `json:"albumIDs,omitempty"`
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

func (t Track) GetDetails() MediaMetadata { return t.Details }
func (t Track) GetMediaType() MediaType   { return MEDIATYPE_TRACK }

func (a Album) GetDetails() MediaMetadata { return a.Details }
func (a Album) GetMediaType() MediaType   { return MEDIATYPE_ALBUM }

func (a Artist) GetDetails() MediaMetadata { return a.Details }
func (a Artist) GetMediaType() MediaType   { return MEDIATYPE_ARTIST }

// MusicProvider defines music-related capabilities
type MusicProvider interface {
	SupportsMusic() bool
	GetMusic(ctx context.Context, options *QueryOptions) ([]MediaItem[Track], error)
	GetMusicArtists(ctx context.Context, options *QueryOptions) ([]MediaItem[Artist], error)
	GetMusicAlbums(ctx context.Context, options *QueryOptions) ([]MediaItem[Album], error)
	GetMusicTrackByID(ctx context.Context, id string) (MediaItem[Track], error)
	GetMusicGenres(ctx context.Context) ([]string, error)
}
