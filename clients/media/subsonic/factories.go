package subsonic

import (
	"context"
	"time"

	gosonic "github.com/supersonic-app/go-subsonic/subsonic"
	"suasor/clients/media"
	"suasor/clients/media/types"
	"suasor/di/container"
	"suasor/utils/logger"
)

// RegisterMediaItemFactories registers all media item factories for Subsonic
func RegisterMediaItemFactories(c *container.Container) {
	registry := container.MustGet[media.ClientItemRegistry](c)

	// Register track factory
	media.RegisterFactory[*SubsonicClient, *gosonic.Child, *types.Track](
		&registry,
		func(client *SubsonicClient, ctx context.Context, item *gosonic.Child) (*types.Track, error) {
			return client.trackFactory(ctx, item)
		},
	)

	// Register album factory
	media.RegisterFactory[*SubsonicClient, *gosonic.AlbumID3, *types.Album](
		&registry,
		func(client *SubsonicClient, ctx context.Context, item *gosonic.AlbumID3) (*types.Album, error) {
			return client.albumFactory(ctx, item)
		},
	)

	// Register artist factory
	media.RegisterFactory[*SubsonicClient, *gosonic.ArtistID3, *types.Artist](
		&registry,
		func(client *SubsonicClient, ctx context.Context, item *gosonic.ArtistID3) (*types.Artist, error) {
			return client.artistFactory(ctx, item)
		},
	)

	// Register playlist factory
	media.RegisterFactory[*SubsonicClient, *gosonic.Playlist, *types.Playlist](
		&registry,
		func(client *SubsonicClient, ctx context.Context, item *gosonic.Playlist) (*types.Playlist, error) {
			return client.playlistFactory(ctx, item)
		},
	)
}

// Factory function for Track
func (c *SubsonicClient) trackFactory(ctx context.Context, song *gosonic.Child) (*types.Track, error) {
	log := logger.LoggerFromContext(ctx)

	if song.ID == "" {
		return nil, nil
	}

	log.Debug().
		Str("trackID", song.ID).
		Str("trackName", song.Title).
		Str("artist", song.Artist).
		Msg("Converting Subsonic item to track format")

	// Duration is in seconds, convert to int64
	duration := time.Duration(song.Duration) * time.Second

	track := &types.Track{
		Details: &types.MediaDetails{
			Title:       song.Title,
			Duration:    int64(duration.Seconds()),
			ReleaseYear: song.Year,
			Genres:      []string{song.Genre},
			Artwork: types.Artwork{
				Poster: c.GetCoverArtURL(song.CoverArt),
			},
		},
		ArtistName: song.Artist,
		AlbumName:  song.Album,
		Number:     song.Track,
	}

	// Add external IDs if available
	if song.MusicBrainzID != "" {
		track.Details.ExternalIDs.AddOrUpdate("musicbrainz", song.MusicBrainzID)
	}

	log.Debug().
		Str("trackID", song.ID).
		Str("trackTitle", track.Details.Title).
		Str("artist", track.ArtistName).
		Str("album", track.AlbumName).
		Msg("Successfully converted Subsonic item to track")

	return track, nil
}

// Factory function for Album
func (c *SubsonicClient) albumFactory(ctx context.Context, album *gosonic.AlbumID3) (*types.Album, error) {
	log := logger.LoggerFromContext(ctx)

	if album.ID == "" {
		return nil, nil
	}

	log.Debug().
		Str("albumID", album.ID).
		Str("albumName", album.Name).
		Str("artist", album.Artist).
		Msg("Converting Subsonic item to album format")

	musicAlbum := &types.Album{
		Details: &types.MediaDetails{
			Title:       album.Name,
			ReleaseYear: album.Year,
			Duration:    int64(album.Duration),
			Genres:      []string{album.Genre},
			Artwork: types.Artwork{
				Poster: c.GetCoverArtURL(album.CoverArt),
			},
		},
		ArtistName: album.Artist,
		TrackCount: album.SongCount,
	}

	log.Debug().
		Str("albumID", album.ID).
		Str("albumTitle", musicAlbum.Details.Title).
		Str("artist", musicAlbum.ArtistName).
		Msg("Successfully converted Subsonic item to album")

	return musicAlbum, nil
}

// Factory function for Artist
func (c *SubsonicClient) artistFactory(ctx context.Context, artist *gosonic.ArtistID3) (*types.Artist, error) {
	log := logger.LoggerFromContext(ctx)

	if artist.ID == "" {
		return nil, nil
	}

	log.Debug().
		Str("artistID", artist.ID).
		Str("artistName", artist.Name).
		Msg("Converting Subsonic item to artist format")

	musicArtist := &types.Artist{
		Details: &types.MediaDetails{
			Title: artist.Name,
			Artwork: types.Artwork{
				Poster: c.GetCoverArtURL(artist.CoverArt),
			},
		},
	}

	log.Debug().
		Str("artistID", artist.ID).
		Str("artistName", musicArtist.Details.Title).
		Msg("Successfully converted Subsonic item to artist")

	return musicArtist, nil
}

// Factory function for Playlist
func (c *SubsonicClient) playlistFactory(ctx context.Context, pl *gosonic.Playlist) (*types.Playlist, error) {
	log := logger.LoggerFromContext(ctx)

	if pl.ID == "" {
		return nil, nil
	}

	log.Debug().
		Str("playlistID", pl.ID).
		Str("playlistName", pl.Name).
		Msg("Converting Subsonic item to playlist format")

	playlist := &types.Playlist{
		ItemList: types.ItemList{
			Details: &types.MediaDetails{
				Title:       pl.Name,
				Description: pl.Comment,
				Duration:    int64(pl.Duration),
				Artwork: types.Artwork{
					Poster: c.GetCoverArtURL(pl.CoverArt),
				},
			},
			ItemCount: pl.SongCount,
			IsPublic:  pl.Public,
		},
	}

	log.Debug().
		Str("playlistID", pl.ID).
		Str("playlistTitle", playlist.Details.Title).
		Int("songCount", playlist.ItemCount).
		Msg("Successfully converted Subsonic item to playlist")

	return playlist, nil
}
