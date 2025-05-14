package subsonic

import (
	"context"
	"fmt"
	gosonic "github.com/supersonic-app/go-subsonic/subsonic"
	"strconv"
	"strings"

	"net/url"
	"suasor/clients/media"
	"suasor/clients/media/types"
	mediatypes "suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/logger"
)

// GetStreamURL returns the URL to stream a music track
func (c *SubsonicClient) GetStreamURL(ctx context.Context, trackID string) (string, error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("trackID", trackID).
		Msg("Generating stream URL for track")

	// Create query parameters
	params := url.Values{}
	params.Add("id", trackID)
	params.Add("f", "xml")
	params.Add("v", "1.15.0")
	params.Add("c", "suasor")
	params.Add("u", c.subsonicConfig().GetUsername())
	params.Add("p", c.subsonicConfig().GetPassword())

	streamURL := fmt.Sprintf("%s/rest/stream.view?%s",
		c.subsonicConfig().GetBaseURL(), params.Encode())

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

	// Create query parameters
	params := url.Values{}
	params.Add("id", coverArtID)
	params.Add("f", "xml")
	params.Add("v", "1.15.0")
	params.Add("c", "suasor")
	params.Add("u", c.subsonicConfig().GetUsername())
	params.Add("p", c.subsonicConfig().GetPassword())

	return fmt.Sprintf("%s/rest/getCoverArt.view?%s",
		c.subsonicConfig().GetBaseURL(), params.Encode())
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

	mediaItem := models.NewMediaItem[*mediatypes.Track](
		track)

	return mediaItem, nil
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

	mediaItem := models.NewMediaItem[*mediatypes.Album](
		album)

	mediaItem.SetClientInfo(client.GetClientID(), client.GetClientType(), item.ID)

	return mediaItem, nil
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

	mediaItem := models.NewMediaItem[*mediatypes.Artist](
		artist)
	mediaItem.SetClientInfo(client.GetClientID(), client.GetClientType(), item.ID)

	return mediaItem, nil
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

	mediaItem := models.NewMediaItem[*mediatypes.Playlist](
		playlist,
	)
	mediaItem.SetClientInfo(client.GetClientID(), client.GetClientType(), item.ID)

	return mediaItem, nil
}

// Helper function to check if any typed filter is present in the options
func hasAnyTypedFilter(options *types.QueryOptions) bool {
	if options == nil {
		return false
	}

	return options.ClientAlbumID != "" ||
		options.ClientArtistID != "" ||
		options.Genre != "" ||
		options.Year != 0
}

// Helper function to build a search query string from options
func buildQueryString(options *types.QueryOptions) string {
	if options == nil {
		return ""
	}

	var parts []string

	// Add the main query
	if options.Query != "" {
		parts = append(parts, options.Query)
	}

	// Add artist filter
	if options.ClientArtistID != "" {
		parts = append(parts, "artist:"+options.ClientArtistID)
	}

	// Add album filter
	if options.ClientAlbumID != "" {
		parts = append(parts, "album:"+options.ClientAlbumID)
	}

	// Add genre filter
	if options.Genre != "" {
		parts = append(parts, "genre:"+options.Genre)
	}

	// Add year filter
	if options.Year != 0 {
		parts = append(parts, "year:"+strconv.Itoa(options.Year))
	}

	return strings.Join(parts, " ")
}
