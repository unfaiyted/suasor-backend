// interfaces/media_client.go
package interfaces

import (
	"suasor/client/media/types"
)

// MediaClient defines basic client information that all providers must implement
type MediaClient interface {
	GetClientID() uint64
	GetClientType() types.MediaClientType
}
