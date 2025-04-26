package media

import (
	"context"
	"errors"
	"fmt"
	"suasor/clients"
	media "suasor/clients/media/types"
	types "suasor/clients/types"
	models "suasor/types/models"
	"suasor/types/responses"
	"suasor/utils/logger"
)

var ErrFeatureNotSupported = errors.New("feature not supported by this media client")

// ClientMedia defines basic client information that all providers must implement

type ClientMedia interface {
	clients.Client
	SupportsMovies() bool
	SupportsSeries() bool
	SupportsMusic() bool
	SupportsPlaylists() bool
	SupportsCollections() bool
	SupportsHistory() bool

	GetRegistry() *ClientItemRegistry

	Search(ctx context.Context, options *media.QueryOptions) (responses.SearchResults, error)
	AsGenericClient() clients.Client
}

type clientMedia struct {
	clients.Client
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

	// Create a new client with the provided config
	client := clients.NewClient(clientID, clientType.AsCategory(), config)

	return &clientMedia{
		Client:       client,
		config:       &config,
		ItemRegistry: itemRegistry,
		ClientType:   clientType,
	}, nil
}

// Default caity implementations (all false by default)
func (m *clientMedia) SupportsMovies() bool      { return false }
func (m *clientMedia) SupportsSeries() bool      { return false }
func (m *clientMedia) SupportsMusic() bool       { return false }
func (m *clientMedia) SupportsPlaylists() bool   { return false }
func (m *clientMedia) SupportsCollections() bool { return false }
func (m *clientMedia) SupportsHistory() bool     { return false }

func (b *clientMedia) GetRegistry() *ClientItemRegistry {
	return b.ItemRegistry
}

func (b *clientMedia) AsGenericClient() clients.Client {
	return b
}

// Embed in your clients to provide default behavior
func (b *clientMedia) GetMovies(ctx context.Context, options *media.QueryOptions) ([]*models.MediaItem[*media.Movie], error) {
	return nil, ErrFeatureNotSupported
}

func (b *clientMedia) GetSeries(ctx context.Context, options *media.QueryOptions) ([]*models.MediaItem[*media.Series], error) {
	return nil, ErrFeatureNotSupported
}

func (b *clientMedia) GetMusic(ctx context.Context, options *media.QueryOptions) ([]*models.MediaItem[*media.Track], error) {
	return nil, ErrFeatureNotSupported
}

func (b *clientMedia) GetPlaylists(ctx context.Context, options *media.QueryOptions) ([]*models.MediaItem[*media.Playlist], error) {
	return nil, ErrFeatureNotSupported
}

func (b *clientMedia) GetCollections(ctx context.Context, options *media.QueryOptions) ([]*models.MediaItem[*media.Collection], error) {
	return nil, ErrFeatureNotSupported
}

func (b *clientMedia) GetHistory(ctx context.Context, options *media.QueryOptions) ([]*models.UserMediaItemData[media.MediaData], error) {
	return nil, ErrFeatureNotSupported
}

func (b *clientMedia) ToMediaItem(ctx context.Context, item media.MediaData, itemID string) (models.MediaItem[media.MediaData], error) {
	if item == nil {
		return models.MediaItem[media.MediaData]{}, fmt.Errorf("cannot convert nil item to media item")
	}

	mediaItem := models.MediaItem[media.MediaData]{
		Data: item,
		Type: item.GetMediaType(),
	}
	mediaItem.SetClientInfo(b.GetClientID(), b.GetClientType(), itemID)

	return mediaItem, nil
}

func (b *clientMedia) TestConnection(ctx context.Context) (bool, error) {
	return false, nil
}

// TODO: This implementation should work, but it also isnt going to be as fast as if we
// write a specific implementation for each type as needed on the client side. This is
// because we might be able to reduce the number of requests to the client by doing a
// more generic search at the client level.
// This should work on for all clients and not need any special implemenations for them
func (b *clientMedia) Search(ctx context.Context, options *media.QueryOptions) (responses.SearchResults, error) {
	log := logger.LoggerFromContext(ctx)
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
