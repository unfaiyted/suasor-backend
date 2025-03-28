// interfaces/base_client.go
package interfaces

import (
	"errors"
	client "suasor/client/types"
)

var ErrFeatureNotSupported = errors.New("feature not supported by this media client")

// BaseClient provides common behavior for all media clients
type BaseClient struct {
	ClientID uint64
	Category client.ClientCategory
}

// Get client information
func (b *BaseClient) GetClientID() uint64                { return b.ClientID }
func (b *BaseClient) GetCategory() client.ClientCategory { return b.Category }
