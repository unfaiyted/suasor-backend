package emby

import (
	"context"
	"fmt"

	"github.com/antihax/optional"
	"suasor/clients/media/types"
	embyclient "suasor/internal/clients/embyAPI"
	"suasor/types/models"
	"suasor/utils/logger"
)

func (e *EmbyClient) SupportsHistory() bool { return true }

// GetWatchHistory retrieves watch history from the Emby server
func (e *EmbyClient) GetPlayHistory(ctx context.Context, options *types.QueryOptions) (*models.MediaItemDataList, error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.GetClientID()).
		Str("clientType", string(e.GetClientType())).
		Msg("Retrieving watch history from Emby server")

	if e.embyConfig().UserID == "" {
		return nil, fmt.Errorf("user ID is required for watch history")
	}

	// queryParams := embyclient.ItemsServiceApiGetUsersByUseridItemsOpts{
	// 	IsPlayed:  optional.NewBool(true),
	// 	Recursive: optional.NewBool(true),
	// }

	userID := e.getUserID()

	queryParams := embyclient.ItemsServiceApiGetUsersByUseridItemsOpts{
		Recursive:        optional.NewBool(true),
		IncludeItemTypes: optional.NewString("Movie,Episode,Audio,Playlist,Series,Season"),
		IsPlayed:         optional.NewBool(true),
		EnableUserData:   optional.NewBool(true),
		Fields:           optional.NewString("PrimaryImageAspectRatio,BasicSyncInfo,CanDelete,Container,DateCreated,PremiereDate,Genres,MediaSourceCount,MediaSources,Overview,ParentId,Path,SortName,Studios,Taglines,ProviderIds,CommunityRating,CriticRating,UserData"),
	}

	// Apply options for pagination
	if options != nil {
		if options.Limit > 0 {
			queryParams.Limit = optional.NewInt32(int32(options.Limit))
		}
		if options.Offset > 0 {
			queryParams.StartIndex = optional.NewInt32(int32(options.Offset))
		}
	}
	queryParams.Limit = optional.NewInt32(100)
	queryParams.StartIndex = optional.NewInt32(0)
	// Apply options
	// ApplyClientQueryOptions(ctx, &queryParams, options)

	log.Debug().
		Int32("limit", queryParams.Limit.Value()).
		Int32("offset", queryParams.StartIndex.Value()).
		Str("userID", userID).
		Msg("Applying query options")

	// Call the Emby API
	items, resp, err := e.client.ItemsServiceApi.GetUsersByUseridItems(ctx, userID, &queryParams)
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", e.embyConfig().GetBaseURL()).
			// Str("apiEndpoint", "/Users/"+e.embyConfig().UserID+"/Items").
			Msg("Failed to fetch watch history from Emby")
		return nil, fmt.Errorf("failed to fetch watch history: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(items.Items)).
		Int("totalRecordCount", int(items.TotalRecordCount)).
		Msg("Successfully retrieved watch history from Emby")

	history, err := GetMixedMediaItemsData(e, ctx, items.Items)
	if err != nil {
		return nil, err
	}

	return history, nil
}
