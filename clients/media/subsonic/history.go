package subsonic

import (
	"context"
	t "suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/logger"
)

// Subsonic doesn't have a formal way to get history, but we can implement a basic version
func (c *SubsonicClient) SupportsHistory() bool {
	return false
}

// GetHistory returns user's play history if available
func (c *SubsonicClient) GetHistory(ctx context.Context, options *t.QueryOptions) ([]*models.UserMediaItemData[t.MediaData], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Msg("History retrieval not supported by Subsonic")

	// Return empty slice to indicate no history
	return []*models.UserMediaItemData[t.MediaData]{}, nil
}
