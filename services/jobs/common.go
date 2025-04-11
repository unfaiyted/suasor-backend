package jobs

import (
	clienttypes "suasor/client/types"
)

// MediaClientInfo is a structure to store media client information
type MediaClientInfo struct {
	ClientID   uint64
	ClientType clienttypes.MediaClientType
	Name       string
	UserID     uint64
}