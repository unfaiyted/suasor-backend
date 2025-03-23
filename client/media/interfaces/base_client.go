// interfaces/base_client.go
package interfaces

import (
	"context"
	"errors"
)

var ErrFeatureNotSupported = errors.New("feature not supported by this media client")

// BaseMediaClient provides common behavior for all media clients
type BaseMediaClient struct {
	ClientID   uint64
	ClientType MediaClientType
}

// Get client information
func (b *BaseMediaClient) GetClientID() uint64            { return b.ClientID }
func (b *BaseMediaClient) GetClientType() MediaClientType { return b.ClientType }

// Default capability implementations (all false by default)
func (b *BaseMediaClient) SupportsMovies() bool      { return false }
func (b *BaseMediaClient) SupportsTVShows() bool     { return false }
func (b *BaseMediaClient) SupportsMusic() bool       { return false }
func (b *BaseMediaClient) SupportsPlaylists() bool   { return false }
func (b *BaseMediaClient) SupportsCollections() bool { return false }

// Default error implementation for unsupported features
// Embed in your clients to provide default behavior
func (b *BaseMediaClient) GetMovies(ctx context.Context, options *QueryOptions) ([]Movie, error) {
	return nil, ErrFeatureNotSupported
}
