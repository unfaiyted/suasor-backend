package providers

import (
	"context"
	"suasor/client/media/types"
)

// HistoryProvider defines watch and play history capabilities
type HistoryProvider interface {
	SupportsHistory() bool
	GetPlayHistory(ctx context.Context, options *types.QueryOptions) ([]types.MediaPlayHistory[types.MediaData], error)
}
