package radarr

import (
	"context"
	"fmt"

	"suasor/clients/automation/types"
	"suasor/utils/logger"
)

func (r *RadarrClient) GetQualityProfiles(ctx context.Context) ([]types.QualityProfile, error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", r.ClientID).
		Str("clientType", string(r.ClientType)).
		Msg("Retrieving quality profiles from Radarr")

	profiles, resp, err := r.client.QualityProfileAPI.ListQualityProfile(ctx).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to fetch quality profiles from Radarr")
		return nil, fmt.Errorf("failed to fetch quality profiles: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("profileCount", len(profiles)).
		Msg("Successfully retrieved quality profiles from Radarr")

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
