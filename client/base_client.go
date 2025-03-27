// interfaces/base_client.go
package interfaces

import (
	"context"
	"errors"
	"fmt"
	media "suasor/client/media/types"
	client "suasor/client/types"
	"suasor/types/models"
)

var ErrFeatureNotSupported = errors.New("feature not supported by this media client")

// BaseMediaClient provides common behavior for all media clients
type BaseMediaClient struct {
	ClientID   uint64
	ClientType client.MediaClientType
}

// Get client information
func (b *BaseMediaClient) GetClientID() uint64                   { return b.ClientID }
func (b *BaseMediaClient) GetClientType() client.MediaClientType { return b.ClientType }

// Default capability implementations (all false by default)
func (b *BaseMediaClient) SupportsMovies() bool      { return false }
func (b *BaseMediaClient) SupportsTVShows() bool     { return false }
func (b *BaseMediaClient) SupportsMusic() bool       { return false }
func (b *BaseMediaClient) SupportsPlaylists() bool   { return false }
func (b *BaseMediaClient) SupportsCollections() bool { return false }
func (b *BaseMediaClient) SupportsHistory() bool     { return false }

// Default error implementation for unsupported features
// Embed in your clients to provide default behavior
func (b *BaseMediaClient) GetMovies(ctx context.Context, options *media.QueryOptions) ([]models.MediaItem[media.Movie], error) {
	return nil, ErrFeatureNotSupported
}

func (b *BaseMediaClient) GetTVShows(ctx context.Context, options *media.QueryOptions) ([]models.MediaItem[media.TVShow], error) {
	return nil, ErrFeatureNotSupported
}

func (b *BaseMediaClient) GetMusic(ctx context.Context, options *media.QueryOptions) ([]models.MediaItem[media.Track], error) {
	return nil, ErrFeatureNotSupported
}

func (b *BaseMediaClient) GetPlaylists(ctx context.Context, options *media.QueryOptions) ([]models.MediaItem[media.Playlist], error) {
	return nil, ErrFeatureNotSupported
}

func (b *BaseMediaClient) GetCollections(ctx context.Context, options *media.QueryOptions) ([]models.MediaItem[media.Collection], error) {
	return nil, ErrFeatureNotSupported
}

func (b *BaseMediaClient) GetHistory(ctx context.Context, options *media.QueryOptions) ([]models.MediaPlayHistory[media.MediaData], error) {
	return nil, ErrFeatureNotSupported
}

func (b *BaseMediaClient) ToMediaItem(ctx context.Context, item media.MediaData, itemID string) (models.MediaItem[media.MediaData], error) {
	if item == nil {
		return models.MediaItem[media.MediaData]{}, fmt.Errorf("cannot convert nil item to media item")
	}

	mediaItem := models.MediaItem[media.MediaData]{
		Data: item,
		Type: item.GetMediaType(),
	}
	mediaItem.SetClientInfo(b.ClientID, b.ClientType, itemID)

	return mediaItem, nil
}

func (b *BaseMediaClient) ToMediaItemEpisode(ctx context.Context, item media.Episode, itemID string) (models.MediaItem[media.Episode], error) {
	mediaItem := models.MediaItem[media.Episode]{
		Data: item,
		Type: item.GetMediaType(),
	}
	mediaItem.SetClientInfo(b.ClientID, b.ClientType, itemID)

	return mediaItem, nil
}

func (b *BaseMediaClient) ToMediaItemSeason(ctx context.Context, item media.Season, itemID string) (models.MediaItem[media.Season], error) {
	mediaItem := models.MediaItem[media.Season]{
		Data: item,
		Type: item.GetMediaType(),
	}
	mediaItem.SetClientInfo(b.ClientID, b.ClientType, itemID)

	return mediaItem, nil
}

func (b *BaseMediaClient) ToMediaItemPlaylist(ctx context.Context, item media.Playlist, itemID string) (models.MediaItem[media.Playlist], error) {
	mediaItem := models.MediaItem[media.Playlist]{
		Data: item,
		Type: item.GetMediaType(),
	}
	mediaItem.SetClientInfo(b.ClientID, b.ClientType, itemID)

	return mediaItem, nil
}

func (b *BaseMediaClient) ToMediaItemCollection(ctx context.Context, item media.Collection, itemID string) (models.MediaItem[media.Collection], error) {
	mediaItem := models.MediaItem[media.Collection]{
		Data: item,
		Type: item.GetMediaType(),
	}
	mediaItem.SetClientInfo(b.ClientID, b.ClientType, itemID)

	return mediaItem, nil
}

func (b *BaseMediaClient) ToMediaItemMovie(ctx context.Context, item media.Movie, itemID string) (models.MediaItem[media.Movie], error) {
	mediaItem := models.MediaItem[media.Movie]{
		Data: item,
		Type: item.GetMediaType(),
	}
	mediaItem.SetClientInfo(b.ClientID, b.ClientType, itemID)

	return mediaItem, nil
}

func (b *BaseMediaClient) ToMediaItemSeries(ctx context.Context, item media.TVShow, itemID string) (models.MediaItem[media.TVShow], error) {
	mediaItem := models.MediaItem[media.TVShow]{
		Data: item,
		Type: item.GetMediaType(),
	}
	mediaItem.SetClientInfo(b.ClientID, b.ClientType, itemID)

	return mediaItem, nil
}

func (b *BaseMediaClient) ToMediaItemTrack(ctx context.Context, item media.Track, itemID string) (models.MediaItem[media.Track], error) {
	mediaItem := models.MediaItem[media.Track]{
		Data: item,
		Type: item.GetMediaType(),
	}
	mediaItem.SetClientInfo(b.ClientID, b.ClientType, itemID)

	return mediaItem, nil
}

func (b *BaseMediaClient) ToMediaItemAlbum(ctx context.Context, item media.Album, itemID string) (models.MediaItem[media.Album], error) {
	mediaItem := models.MediaItem[media.Album]{
		Data: item,
		Type: item.GetMediaType(),
	}
	mediaItem.SetClientInfo(b.ClientID, b.ClientType, itemID)

	return mediaItem, nil
}

func (b *BaseMediaClient) ToMediaItemArtist(ctx context.Context, item media.Artist, itemID string) (models.MediaItem[media.Artist], error) {
	mediaItem := models.MediaItem[media.Artist]{
		Data: item,
		Type: item.GetMediaType(),
	}
	mediaItem.SetClientInfo(b.ClientID, b.ClientType, itemID)

	return mediaItem, nil
}
