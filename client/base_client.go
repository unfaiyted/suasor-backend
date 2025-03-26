// interfaces/base_client.go
package interfaces

import (
	"context"
	"errors"
	"fmt"
	"suasor/client/media/types"
)

var ErrFeatureNotSupported = errors.New("feature not supported by this media client")

// BaseMediaClient provides common behavior for all media clients
type BaseMediaClient struct {
	ClientID   uint64
	ClientType types.MediaClientType
}

// Get client information
func (b *BaseMediaClient) GetClientID() uint64                  { return b.ClientID }
func (b *BaseMediaClient) GetClientType() types.MediaClientType { return b.ClientType }

// Default capability implementations (all false by default)
func (b *BaseMediaClient) SupportsMovies() bool      { return false }
func (b *BaseMediaClient) SupportsTVShows() bool     { return false }
func (b *BaseMediaClient) SupportsMusic() bool       { return false }
func (b *BaseMediaClient) SupportsPlaylists() bool   { return false }
func (b *BaseMediaClient) SupportsCollections() bool { return false }
func (b *BaseMediaClient) SupportsHistory() bool     { return false }

// Default error implementation for unsupported features
// Embed in your clients to provide default behavior
func (b *BaseMediaClient) GetMovies(ctx context.Context, options *types.QueryOptions) ([]types.MediaItem[types.Movie], error) {
	return nil, ErrFeatureNotSupported
}

func (b *BaseMediaClient) GetTVShows(ctx context.Context, options *types.QueryOptions) ([]types.MediaItem[types.TVShow], error) {
	return nil, ErrFeatureNotSupported
}

func (b *BaseMediaClient) GetMusic(ctx context.Context, options *types.QueryOptions) ([]types.MediaItem[types.Track], error) {
	return nil, ErrFeatureNotSupported
}

func (b *BaseMediaClient) GetPlaylists(ctx context.Context, options *types.QueryOptions) ([]types.MediaItem[types.Playlist], error) {
	return nil, ErrFeatureNotSupported
}

func (b *BaseMediaClient) GetCollections(ctx context.Context, options *types.QueryOptions) ([]types.MediaItem[types.Collection], error) {
	return nil, ErrFeatureNotSupported
}

func (b *BaseMediaClient) GetHistory(ctx context.Context, options *types.QueryOptions) ([]types.MediaPlayHistory[types.MediaData], error) {
	return nil, ErrFeatureNotSupported
}

func (b *BaseMediaClient) ToMediaItem(ctx context.Context, item types.MediaData, itemID string) (types.MediaItem[types.MediaData], error) {
	if item == nil {
		return types.MediaItem[types.MediaData]{}, fmt.Errorf("cannot convert nil item to media item")
	}

	mediaItem := types.MediaItem[types.MediaData]{
		Data: item,
		Type: item.GetMediaType(),
	}
	mediaItem.SetClientInfo(b.ClientID, b.ClientType, itemID)

	return mediaItem, nil
}

func (b *BaseMediaClient) ToMediaItemEpisode(ctx context.Context, item types.Episode, itemID string) (types.MediaItem[types.Episode], error) {
	mediaItem := types.MediaItem[types.Episode]{
		Data: item,
		Type: item.GetMediaType(),
	}
	mediaItem.SetClientInfo(b.ClientID, b.ClientType, itemID)

	return mediaItem, nil
}

func (b *BaseMediaClient) ToMediaItemSeason(ctx context.Context, item types.Season, itemID string) (types.MediaItem[types.Season], error) {
	mediaItem := types.MediaItem[types.Season]{
		Data: item,
		Type: item.GetMediaType(),
	}
	mediaItem.SetClientInfo(b.ClientID, b.ClientType, itemID)

	return mediaItem, nil
}

func (b *BaseMediaClient) ToMediaItemPlaylist(ctx context.Context, item types.Playlist, itemID string) (types.MediaItem[types.Playlist], error) {
	mediaItem := types.MediaItem[types.Playlist]{
		Data: item,
		Type: item.GetMediaType(),
	}
	mediaItem.SetClientInfo(b.ClientID, b.ClientType, itemID)

	return mediaItem, nil
}

func (b *BaseMediaClient) ToMediaItemCollection(ctx context.Context, item types.Collection, itemID string) (types.MediaItem[types.Collection], error) {
	mediaItem := types.MediaItem[types.Collection]{
		Data: item,
		Type: item.GetMediaType(),
	}
	mediaItem.SetClientInfo(b.ClientID, b.ClientType, itemID)

	return mediaItem, nil
}

func (b *BaseMediaClient) ToMediaItemMovie(ctx context.Context, item types.Movie, itemID string) (types.MediaItem[types.Movie], error) {
	mediaItem := types.MediaItem[types.Movie]{
		Data: item,
		Type: item.GetMediaType(),
	}
	mediaItem.SetClientInfo(b.ClientID, b.ClientType, itemID)

	return mediaItem, nil
}

func (b *BaseMediaClient) ToMediaItemSeries(ctx context.Context, item types.TVShow, itemID string) (types.MediaItem[types.TVShow], error) {
	mediaItem := types.MediaItem[types.TVShow]{
		Data: item,
		Type: item.GetMediaType(),
	}
	mediaItem.SetClientInfo(b.ClientID, b.ClientType, itemID)

	return mediaItem, nil
}

func (b *BaseMediaClient) ToMediaItemTrack(ctx context.Context, item types.Track, itemID string) (types.MediaItem[types.Track], error) {
	mediaItem := types.MediaItem[types.Track]{
		Data: item,
		Type: item.GetMediaType(),
	}
	mediaItem.SetClientInfo(b.ClientID, b.ClientType, itemID)

	return mediaItem, nil
}

func (b *BaseMediaClient) ToMediaItemAlbum(ctx context.Context, item types.Album, itemID string) (types.MediaItem[types.Album], error) {
	mediaItem := types.MediaItem[types.Album]{
		Data: item,
		Type: item.GetMediaType(),
	}
	mediaItem.SetClientInfo(b.ClientID, b.ClientType, itemID)

	return mediaItem, nil
}

func (b *BaseMediaClient) ToMediaItemArtist(ctx context.Context, item types.Artist, itemID string) (types.MediaItem[types.Artist], error) {
	mediaItem := types.MediaItem[types.Artist]{
		Data: item,
		Type: item.GetMediaType(),
	}
	mediaItem.SetClientInfo(b.ClientID, b.ClientType, itemID)

	return mediaItem, nil
}
