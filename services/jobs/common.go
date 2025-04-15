package jobs

import (
	clienttypes "suasor/client/types"
)

// ClientMediaInfo is a structure to store media client information
type ClientMediaInfo struct {
	ClientID   uint64
	ClientType clienttypes.ClientMediaType
	Name       string
	UserID     uint64
}