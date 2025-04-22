package lidarr

import (
	"context"
	"fmt"

	lidarr "github.com/devopsarr/lidarr-go/lidarr"
	"suasor/clients/automation/types"
	"suasor/utils/logger"
)

func (l *LidarrClient) GetTags(ctx context.Context) ([]types.Tag, error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", l.ClientID).
		Str("clientType", string(l.ClientType)).
		Msg("Retrieving tags from Lidarr")

	tags, resp, err := l.client.TagAPI.ListTag(ctx).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to fetch tags from Lidarr")
		return nil, fmt.Errorf("failed to fetch tags: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("tagCount", len(tags)).
		Msg("Successfully retrieved tags from Lidarr")

	// Convert to our internal representation
	result := make([]types.Tag, 0, len(tags))
	for _, tag := range tags {
		result = append(result, types.Tag{
			ID:   int64(tag.GetId()),
			Name: tag.GetLabel(),
		})
	}

	return result, nil
}

func (l *LidarrClient) CreateTag(ctx context.Context, tagName string) (types.Tag, error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", l.ClientID).
		Str("clientType", string(l.ClientType)).
		Str("tagName", tagName).
		Msg("Creating new tag in Lidarr")

	newTag := lidarr.NewTagResource()
	newTag.SetLabel(tagName)

	createdTag, resp, err := l.client.TagAPI.CreateTag(ctx).TagResource(*newTag).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("tagName", tagName).
			Msg("Failed to create tag in Lidarr")

		return types.Tag{}, fmt.Errorf("failed to create tag: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int32("tagID", createdTag.GetId()).
		Str("tagName", createdTag.GetLabel()).
		Msg("Successfully created tag in Lidarr")

	return types.Tag{
		ID:   int64(createdTag.GetId()),
		Name: createdTag.GetLabel(),
	}, nil

}
