package subsonic

import (
	"context"
	t "suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/logger"
)

// GetHistory returns user's play history if available
func (c *SubsonicClient) GetHistory(ctx context.Context, options *t.QueryOptions) ([]*models.UserMediaItemData[t.MediaData], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Msg("History retrieval not supported by Subsonic")

	// Return empty slice to indicate no history
	return []*models.UserMediaItemData[t.MediaData]{}, nil
}

