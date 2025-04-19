// interfaces/media_client.go
package media

import (
	"context"
	"errors"
	"fmt"
	"suasor/client"
	media "suasor/client/media/types"
	types "suasor/client/types"
	models "suasor/types/models"
	"suasor/types/responses"
	"suasor/utils"
)

var ErrFeatureNotSupported = errors.New("feature not supported by this media client")

// ClientMedia defines basic client information that all providers must implement

type ClientMedia interface {
	client.Client
	SupportsMovies() bool
	SupportsSeries() bool
	SupportsMusic() bool
	SupportsPlaylists() bool
	SupportsCollections() bool
	SupportsHistory() bool

	GetRegistry() *ClientItemRegistry

	Search(ctx context.Context, options *media.QueryOptions) (responses.SearchResults, error)
}

type BaseClientMedia struct {
	client.BaseClient
	ClientType   types.ClientMediaType
	ItemRegistry *ClientItemRegistry
	config       *types.ClientMediaConfig
}

func NewClientMedia(
	ctx context.Context,
	clientID uint64,
	clientType types.ClientMediaType,
	itemRegistry *ClientItemRegistry,
	config types.ClientMediaConfig) (ClientMedia, error) {
	return &BaseClientMedia{
		BaseClient: client.BaseClient{
			ClientID: clientID,
			Category: clientType.AsCategory(),
		},

		config:       &config,
		ItemRegistry: itemRegistry,
		ClientType:   clientType,
	}, nil
}

// Default caity implementations (all false by default)
func (m *BaseClientMedia) SupportsMovies() bool      { return false }
func (m *BaseClientMedia) SupportsSeries() bool      { return false }
func (m *BaseClientMedia) SupportsMusic() bool       { return false }
func (m *BaseClientMedia) SupportsPlaylists() bool   { return false }
func (m *BaseClientMedia) SupportsCollections() bool { return false }
func (m *BaseClientMedia) SupportsHistory() bool     { return false }

func (b *BaseClientMedia) GetRegistry() *ClientItemRegistry {
	return b.ItemRegistry
}

// Embed in your clients to provide default behavior
func (b *BaseClientMedia) GetMovies(ctx context.Context, options *media.QueryOptions) ([]*models.MediaItem[*media.Movie], error) {
	return nil, ErrFeatureNotSupported
}

func (b *BaseClientMedia) GetSeries(ctx context.Context, options *media.QueryOptions) ([]*models.MediaItem[*media.Series], error) {
	return nil, ErrFeatureNotSupported
}

func (b *BaseClientMedia) GetMusic(ctx context.Context, options *media.QueryOptions) ([]*models.MediaItem[*media.Track], error) {
	return nil, ErrFeatureNotSupported
}

func (b *BaseClientMedia) GetPlaylists(ctx context.Context, options *media.QueryOptions) ([]*models.MediaItem[*media.Playlist], error) {
	return nil, ErrFeatureNotSupported
}

func (b *BaseClientMedia) GetCollections(ctx context.Context, options *media.QueryOptions) ([]*models.MediaItem[*media.Collection], error) {
	return nil, ErrFeatureNotSupported
}

func (b *BaseClientMedia) GetHistory(ctx context.Context, options *media.QueryOptions) ([]*models.UserMediaItemData[media.MediaData], error) {
	return nil, ErrFeatureNotSupported
}

func (b *BaseClientMedia) ToMediaItem(ctx context.Context, item media.MediaData, itemID string) (models.MediaItem[media.MediaData], error) {
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

// func (b *BaseClientMedia) ToMediaItemEpisode(ctx context.Context, item media.Episode, itemID string) (models.MediaItem[media.Episode], error) {
// 	mediaItem := models.MediaItem[media.Episode]{
// 		Data: item,
// 		Type: item.GetMediaType(),
// 	}
// 	mediaItem.SetClientInfo(b.ClientID, b.ClientType, itemID)
//
// 	return mediaItem, nil
// }
//
// func (b *BaseClientMedia) ToMediaItemSeason(ctx context.Context, item media.Season, itemID string) (models.MediaItem[media.Season], error) {
// 	mediaItem := models.MediaItem[media.Season]{
// 		Data: item,
// 		Type: item.GetMediaType(),
// 	}
// 	mediaItem.SetClientInfo(b.ClientID, b.ClientType, itemID)
//
// 	return mediaItem, nil
// }
//
// func (b *BaseClientMedia) ToMediaItemPlaylist(ctx context.Context, item media.Playlist, itemID string) (models.MediaItem[media.Playlist], error) {
// 	mediaItem := models.MediaItem[media.Playlist]{
// 		Data: item,
// 		Type: item.GetMediaType(),
// 	}
// 	mediaItem.SetClientInfo(b.ClientID, b.ClientType, itemID)
//
// 	return mediaItem, nil
// }
//
// func (b *BaseClientMedia) ToMediaItemCollection(ctx context.Context, item media.Collection, itemID string) (models.MediaItem[media.Collection], error) {
// 	mediaItem := models.MediaItem[media.Collection]{
// 		Data: item,
// 		Type: item.GetMediaType(),
// 	}
// 	mediaItem.SetClientInfo(b.ClientID, b.ClientType, itemID)
//
// 	return mediaItem, nil
// }
//
// func (b *BaseClientMedia) ToMediaItemMovie(ctx context.Context, item media.Movie, itemID string) (models.MediaItem[media.Movie], error) {
// 	mediaItem := models.MediaItem[media.Movie]{
// 		Data: item,
// 		Type: item.GetMediaType(),
// 	}
// 	mediaItem.SetClientInfo(b.ClientID, b.ClientType, itemID)
//
// 	return mediaItem, nil
// }
//
// func (b *BaseClientMedia) ToMediaItemSeries(ctx context.Context, item media.Series, itemID string) (models.MediaItem[media.Series], error) {
// 	mediaItem := models.MediaItem[media.Series]{
// 		Data: item,
// 		Type: item.GetMediaType(),
// 	}
//
// 	mediaItem.SetClientInfo(b.ClientID, b.ClientType, itemID)
//
// 	return mediaItem, nil
// }
//
// func (b *BaseClientMedia) ToMediaItemTrack(ctx context.Context, item media.Track, itemID string) (models.MediaItem[media.Track], error) {
// 	mediaItem := models.MediaItem[media.Track]{
// 		Data: item,
// 		Type: item.GetMediaType(),
// 	}
// 	mediaItem.SetClientInfo(b.ClientID, b.ClientType, itemID)
//
// 	return mediaItem, nil
// }
//
// func (b *BaseClientMedia) ToMediaItemAlbum(ctx context.Context, item media.Album, itemID string) (models.MediaItem[media.Album], error) {
// 	mediaItem := models.MediaItem[media.Album]{
// 		Data: item,
// 		Type: item.GetMediaType(),
// 	}
// 	mediaItem.SetClientInfo(b.ClientID, b.ClientType, itemID)
//
// 	return mediaItem, nil
// }
//
// func (b *BaseClientMedia) ToMediaItemArtist(ctx context.Context, item media.Artist, itemID string) (models.MediaItem[media.Artist], error) {
// 	mediaItem := models.MediaItem[media.Artist]{
// 		Data: item,
// 		Type: item.GetMediaType(),
// 	}
// 	mediaItem.SetClientInfo(b.ClientID, b.ClientType, itemID)
//
// 	return mediaItem, nil
// }

func (b *BaseClientMedia) TestConnection(ctx context.Context) (bool, error) {
	return false, nil
}

// TODO: This implementation should work, but it also isnt going to be as fast as if we
// write a specific implementation for each type as needed on the client side. This is
// because we might be able to reduce the number of requests to the client by doing a
// more generic search at the client level.
// This should work on for all clients and not need any special implemenations for them
func (b *BaseClientMedia) Search(ctx context.Context, options *media.QueryOptions) (responses.SearchResults, error) {
	log := utils.LoggerFromContext(ctx)
	log.Info().Str("query", options.Query).Msg("Searching media items")

	var results responses.SearchResults

	// check if mediaType is empty and if it is set it to ALL
	if options.MediaType == "" {
		options.MediaType = media.MediaTypeAll
	}

	switch options.MediaType {
	case media.MediaTypeMovie:
		movies, err := b.GetMovies(ctx, options)
		if err != nil {
			return results, err
		}
		results.Movies = movies
	case media.MediaTypeSeries:
		series, err := b.GetSeries(ctx, options)
		if err != nil {
			return results, err
		}
		results.Series = series
	case media.MediaTypeTrack, media.MediaTypeAlbum, media.MediaTypeArtist:
		b.GetMusic(ctx, options)
	// case media.MediaTypePerson:
	// results.People = b.GetPeople(ctx, options)
	case media.MediaTypePlaylist:
		playlists, err := b.GetPlaylists(ctx, options)
		if err != nil {
			return results, err
		}
		results.Playlists = playlists
	case media.MediaTypeCollection:
		collections, err := b.GetCollections(ctx, options)
		if err != nil {
			return results, err
		}
		results.Collections = collections
	default:
		movies, err := b.GetMovies(ctx, options)
		if err != nil {
			return results, err
		}
		results.Movies = movies
		series, err := b.GetSeries(ctx, options)
		if err != nil {
			return results, err
		}
		results.Series = series
		b.GetMusic(ctx, options)
		playlists, err := b.GetPlaylists(ctx, options)
		if err != nil {
			return results, err
		}
		results.Playlists = playlists
		collections, err := b.GetCollections(ctx, options)
		if err != nil {
			return results, err
		}
		results.Collections = collections
	}

	return results, nil
}
