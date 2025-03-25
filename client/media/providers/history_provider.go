package providers

import (
	"context"
	"suasor/client/media/types"
)

// WatchHistoryProvider defines watch history capabilities
type WatchHistoryProvider interface {
	SupportsWatchHistory() bool
	GetWatchHistory(ctx context.Context, options *types.QueryOptions) ([]types.WatchHistoryItem[types.MediaData], error)
}
