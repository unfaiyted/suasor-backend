package subsonic

import (
	gosonic "github.com/supersonic-app/go-subsonic/subsonic"
	t "suasor/client/media/types"
	"suasor/models"
	"time"
)

func (c *SubsonicClient) convertChildToTrack(song gosonic.Child) models.MediaItem[t.Track] {
	// Convert duration from seconds to time.Duration
	duration := time.Duration(song.Duration) * time.Second

	track := models.MediaItem[t.Track]{
		Type: "music",
		Data: t.Track{
			Details: t.MediaMetadata{
				Title:       song.Title,
				Duration:    duration,
				ReleaseYear: song.Year, // Use ReleaseYear instead of Year
				Genres:      []string{song.Genre},
				Artwork: t.Artwork{
					Poster: song.CoverArt, // Will be replaced with full URL later
				},
			},
			ArtistName: song.Artist,
			AlbumTitle: song.Album,
			Number:     song.Track, // Number field instead of TrackNumber
		},
	}

	track.SetClientInfo(c.ClientID, c.ClientType, song.ID)
	return track
}
