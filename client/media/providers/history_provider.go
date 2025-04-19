package providers

import (
	"context"
	"suasor/client/media/types"
	"suasor/types/models"
)

// HistoryProvider defines watch and play history capabilities
type HistoryProvider[T types.MediaData] interface {
	SupportsHistory() bool
	GetPlayHistory(ctx context.Context, options *types.QueryOptions) ([]*models.UserMediaItemData[T], error)
}
