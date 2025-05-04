// lists.go
package emby

import (
	"context"
	"suasor/clients/media/types"
	"suasor/types/models"
)

func (e *EmbyClient) GetListItems(ctx context.Context, client *EmbyClient, listID string, options *types.QueryOptions) ([]*models.MediaItem[types.ListData], error) {
	return nil, nil
}
