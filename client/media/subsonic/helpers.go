package subsonic

import (
	"context"
	"fmt"
	gosonic "github.com/supersonic-app/go-subsonic/subsonic"
	"net/url"
	media "suasor/client/media"
	mediatypes "suasor/client/media/types"
	types "suasor/client/types"
	"suasor/types/models"
	"suasor/utils"
)

// GetStreamURL returns the URL to stream a music track
func (c *SubsonicClient) GetStreamURL(ctx context.Context, trackID string) (string, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("trackID", trackID).
		Msg("Generating stream URL for track")

	subsonicConfig := c.Config.(*types.SubsonicConfig)

	// Create query parameters
	params := url.Values{}
	params.Add("id", trackID)
	params.Add("f", "xml")
	params.Add("v", "1.15.0")
	params.Add("c", "suasor")
	params.Add("u", subsonicConfig.Username)
	params.Add("p", subsonicConfig.Password)

	streamURL := fmt.Sprintf("%s/rest/stream.view?%s",
		subsonicConfig.BaseURL, params.Encode())

	log.Debug().
		Str("trackID", trackID).
		Str("streamURL", streamURL).
		Msg("Generated stream URL for track")

	return streamURL, nil
}

// GetCoverArtURL returns the URL to download cover art
func (c *SubsonicClient) GetCoverArtURL(coverArtID string) string {
	if coverArtID == "" {
		return ""
	}

	subsonicConfig := c.Config.(*types.SubsonicConfig)

	// Create query parameters
	params := url.Values{}
	params.Add("id", coverArtID)
	params.Add("f", "xml")
	params.Add("v", "1.15.0")
	params.Add("c", "suasor")
	params.Add("u", subsonicConfig.Username)
	params.Add("p", subsonicConfig.Password)

	return fmt.Sprintf("%s/rest/getCoverArt.view?%s",
		subsonicConfig.BaseURL, params.Encode())
}

// Generic helper functions to work with the factory pattern

// GetTrack converts a Subsonic Child object to a Track using the factory pattern
func GetTrack(
	ctx context.Context,
	client *SubsonicClient,
	item *gosonic.Child,
) (*mediatypes.Track, error) {
	return media.ConvertTo[*SubsonicClient, *gosonic.Child, *mediatypes.Track](
		client, ctx, item)
}

// GetAlbum converts a Subsonic Album object to an Album using the factory pattern
func GetAlbum(
	ctx context.Context,
	client *SubsonicClient,
	item *gosonic.AlbumID3,
) (*mediatypes.Album, error) {
	return media.ConvertTo[*SubsonicClient, *gosonic.AlbumID3, *mediatypes.Album](
		client, ctx, item)
}

// GetArtist converts a Subsonic Artist object to an Artist using the factory pattern
func GetArtist(
	ctx context.Context,
	client *SubsonicClient,
	item *gosonic.ArtistID3,
) (*mediatypes.Artist, error) {
	return media.ConvertTo[*SubsonicClient, *gosonic.ArtistID3, *mediatypes.Artist](
		client, ctx, item)
}

// GetPlaylist converts a Subsonic Playlist object to a Playlist using the factory pattern
func GetPlaylist(
	ctx context.Context,
	client *SubsonicClient,
	item *gosonic.Playlist,
) (*mediatypes.Playlist, error) {
	return media.ConvertTo[*SubsonicClient, *gosonic.Playlist, *mediatypes.Playlist](
		client, ctx, item)
}

// Helper function to convert a Subsonic Child to a MediaItem Track
func GetTrackItem(
	ctx context.Context,
	client *SubsonicClient,
	item *gosonic.Child,
) (*models.MediaItem[*mediatypes.Track], error) {
	track, err := GetTrack(ctx, client, item)
	if err != nil {
		return nil, err
	}

	mediaItem := models.MediaItem[*mediatypes.Track]{
		Data: track,
		Type: track.GetMediaType(),
	}
	mediaItem.SetClientInfo(client.ClientID, client.ClientType, item.ID)

	return &mediaItem, nil
}

// Helper function to convert a Subsonic Album to a MediaItem Album
func GetAlbumItem(
	ctx context.Context,
	client *SubsonicClient,
	item *gosonic.AlbumID3,
) (*models.MediaItem[*mediatypes.Album], error) {
	album, err := GetAlbum(ctx, client, item)
	if err != nil {
		return nil, err
	}

	mediaItem := models.MediaItem[*mediatypes.Album]{
		Data: album,
		Type: album.GetMediaType(),
	}
	mediaItem.SetClientInfo(client.ClientID, client.ClientType, item.ID)

	return &mediaItem, nil
}

// Helper function to convert a Subsonic Artist to a MediaItem Artist
func GetArtistItem(
	ctx context.Context,
	client *SubsonicClient,
	item *gosonic.ArtistID3,
) (*models.MediaItem[*mediatypes.Artist], error) {
	artist, err := GetArtist(ctx, client, item)
	if err != nil {
		return nil, err
	}

	mediaItem := models.MediaItem[*mediatypes.Artist]{
		Data: artist,
		Type: artist.GetMediaType(),
	}
	mediaItem.SetClientInfo(client.ClientID, client.ClientType, item.ID)

	return &mediaItem, nil
}

// Helper function to convert a Subsonic Playlist to a MediaItem Playlist
func GetPlaylistItem(
	ctx context.Context,
	client *SubsonicClient,
	item *gosonic.Playlist,
) (*models.MediaItem[*mediatypes.Playlist], error) {
	playlist, err := GetPlaylist(ctx, client, item)
	if err != nil {
		return nil, err
	}

	mediaItem := models.MediaItem[*mediatypes.Playlist]{
		Data: playlist,
		Type: playlist.GetMediaType(),
	}
	mediaItem.SetClientInfo(client.ClientID, client.ClientType, item.ID)

	return &mediaItem, nil
}
