package providers

import (
	"context"
	"suasor/client/media/types"
	"suasor/types/models"
)

// HistoryProvider defines watch and play history capabilities
type HistoryProvider interface {
	SupportsHistory() bool
	GetPlayHistory(ctx context.Context, options *types.QueryOptions) ([]models.MediaPlayHistory[types.MediaData], error)
}
