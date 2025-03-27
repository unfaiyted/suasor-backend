package sonarr

import (
	"context"
	"fmt"

	sonarr "github.com/devopsarr/sonarr-go/sonarr"
	"suasor/utils"

	"suasor/client/automation/types"
)

func (s *SonarrClient) GetTags(ctx context.Context) ([]types.Tag, error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", s.ClientID).
		Str("clientType", string(s.ClientType)).
		Msg("Retrieving tags from Sonarr")

	tags, resp, err := s.client.TagAPI.ListTag(ctx).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to fetch tags from Sonarr")
		return nil, fmt.Errorf("failed to fetch tags: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("tagCount", len(tags)).
		Msg("Successfully retrieved tags from Sonarr")

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

func (s *SonarrClient) CreateTag(ctx context.Context, tagName string) (types.Tag, error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", s.ClientID).
		Str("clientType", string(s.ClientType)).
		Str("tagName", tagName).
		Msg("Creating new tag in Sonarr")

	newTag := sonarr.NewTagResource()
	newTag.SetLabel(tagName)

	createdTag, resp, err := s.client.TagAPI.CreateTag(ctx).TagResource(*newTag).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("tagName", tagName).
			Msg("Failed to create tag in Sonarr")
		return types.Tag{}, fmt.Errorf("failed to create tag: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int32("tagID", createdTag.GetId()).
		Str("tagName", createdTag.GetLabel()).
		Msg("Successfully created tag in Sonarr")

	return types.Tag{
		ID:   int64(createdTag.GetId()),
		Name: createdTag.GetLabel(),
	}, nil
}
