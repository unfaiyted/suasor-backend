// interfaces/media_client.go
package media

import (
	client "suasor/client/types"
)

// MediaClient defines basic client information that all providers must implement
type MediaClient interface {
	GetClientID() uint64
	GetClientType() client.MediaClientType
}
