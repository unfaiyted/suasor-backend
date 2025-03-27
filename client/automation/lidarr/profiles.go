package lidarr

import (
	"context"
	"fmt"

	"suasor/client/automation/types"
	"suasor/utils"
)

func (l *LidarrClient) GetQualityProfiles(ctx context.Context) ([]types.QualityProfile, error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", l.ClientID).
		Str("clientType", string(l.ClientType)).
		Msg("Retrieving quality profiles from Lidarr")

	profiles, resp, err := l.client.QualityProfileAPI.ListQualityProfile(ctx).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to fetch quality profiles from Lidarr")
		return nil, fmt.Errorf("failed to fetch quality profiles: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("profileCount", len(profiles)).
		Msg("Successfully retrieved quality profiles from Lidarr")

	// Convert to our internal representation
	result := make([]types.QualityProfile, 0, len(profiles))
	for _, profile := range profiles {
		result = append(result, types.QualityProfile{
			ID:   int64(profile.GetId()),
			Name: profile.GetName(),
		})
	}

	return result, nil
}

func (l *LidarrClient) GetMetadataProfiles(ctx context.Context) ([]types.MetadataProfile, error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", l.ClientID).
		Str("clientType", string(l.ClientType)).
		Msg("Retrieving metadata profiles from Lidarr")

	profiles, resp, err := l.client.MetadataProfileAPI.ListMetadataProfile(ctx).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to fetch metadata profiles from Lidarr")
		return nil, fmt.Errorf("failed to fetch metadata profiles: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("profileCount", len(profiles)).
		Msg("Successfully retrieved metadata profiles from Lidarr")

	// Convert to our internal representation
	result := make([]types.MetadataProfile, 0, len(profiles))
	for _, profile := range profiles {
		result = append(result, types.MetadataProfile{
			ID:   profile.GetId(),
			Name: profile.GetName(),
		})
	}

	return result, nil
}
