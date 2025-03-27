package radarr

import (
	"context"
	"fmt"

	radarr "github.com/devopsarr/radarr-go/radarr"
	"suasor/client/automation/types"
	"suasor/utils"
)

func (r *RadarrClient) GetTags(ctx context.Context) ([]types.Tag, error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", r.ClientID).
		Str("clientType", string(r.ClientType)).
		Msg("Retrieving tags from Radarr")

	tags, resp, err := r.client.TagAPI.ListTag(ctx).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to fetch tags from Radarr")
		return nil, fmt.Errorf("failed to fetch tags: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("tagCount", len(tags)).
		Msg("Successfully retrieved tags from Radarr")

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

func (r *RadarrClient) CreateTag(ctx context.Context, tagName string) (types.Tag, error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", r.ClientID).
		Str("clientType", string(r.ClientType)).
		Str("tagName", tagName).
		Msg("Creating new tag in Radarr")

	newTag := radarr.NewTagResource()
	newTag.SetLabel(tagName)

	createdTag, resp, err := r.client.TagAPI.CreateTag(ctx).TagResource(*newTag).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("tagName", tagName).
			Msg("Failed to create tag in Radarr")
		return types.Tag{}, fmt.Errorf("failed to create tag: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int32("tagID", createdTag.GetId()).
		Str("tagName", createdTag.GetLabel()).
		Msg("Successfully created tag in Radarr")

	return types.Tag{
		ID:   int64(createdTag.GetId()),
		Name: createdTag.GetLabel(),
	}, nil
}
