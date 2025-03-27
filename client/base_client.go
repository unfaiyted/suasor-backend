// interfaces/base_client.go
package interfaces

import (
	"errors"
	client "suasor/client/types"
)

var ErrFeatureNotSupported = errors.New("feature not supported by this media client")

// BaseClient provides common behavior for all media clients
type BaseClient struct {
	ClientID   uint64
	ClientType client.ClientType
}

// Get client information
func (b *BaseClient) GetClientID() uint64              { return b.ClientID }
func (b *BaseClient) GetClientType() client.ClientType { return b.ClientType }
