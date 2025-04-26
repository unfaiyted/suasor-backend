// tvshows.go
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

// GetSeriess retrieves TV shows from the Emby server
func (e *EmbyClient) GetSeries(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Series], error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.GetClientID()).
		Str("clientType", string(e.GetClientType())).
		Msg("Retrieving TV shows from Emby server")

	queryParams := embyclient.ItemsServiceApiGetItemsOpts{
		IncludeItemTypes: optional.NewString("Series"),
		Recursive:        optional.NewBool(true),
	}

	ApplyClientQueryOptions(&queryParams, options)

	items, resp, err := e.client.ItemsServiceApi.GetItems(ctx, &queryParams)
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", e.embyConfig().GetBaseURL()).
			Str("apiEndpoint", "/Items").
			Msg("Failed to fetch TV shows from Emby")
		return nil, fmt.Errorf("failed to fetch TV shows: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(items.Items)).
		Int("totalRecordCount", int(items.TotalRecordCount)).
		Msg("Successfully retrieved TV shows from Emby")

	shows := make([]*models.MediaItem[*types.Series], 0)
	for _, item := range items.Items {
		if item.Type_ == "Series" {
			itemShow, err := GetItem[*types.Series](ctx, e, &item)
			show, err := GetMediaItem[*types.Series](ctx, e, itemShow, item.Id)
			if err != nil {
				log.Warn().
					Err(err).
					Str("showID", item.Id).
					Str("showName", item.Name).
					Msg("Error converting Emby item to TV show format")
				continue
			}
			shows = append(shows, show)
		}
	}

	return shows, nil
}

// GetSeriesByID retrieves a specific TV show by ID
func (e *EmbyClient) GetSeriesByID(ctx context.Context, id string) (*models.MediaItem[*types.Series], error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.GetClientID()).
		Str("clientType", string(e.GetClientType())).
		Str("showID", id).
		Msg("Retrieving specific TV show from Emby server")

	queryParams := embyclient.ItemsServiceApiGetItemsOpts{
		Ids:              optional.NewString(id),
		IncludeItemTypes: optional.NewString("Series"),
	}

	items, resp, err := e.client.ItemsServiceApi.GetItems(ctx, &queryParams)
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", e.embyConfig().GetBaseURL()).
			Str("apiEndpoint", "/Items").
			Str("showID", id).
			Msg("Failed to fetch TV show from Emby")
		return &models.MediaItem[*types.Series]{}, fmt.Errorf("failed to fetch TV show: %w", err)
	}

	if len(items.Items) == 0 {
		log.Error().
			Str("showID", id).
			Int("statusCode", resp.StatusCode).
			Msg("No TV show found with the specified ID")
		return &models.MediaItem[*types.Series]{}, fmt.Errorf("TV show with ID %s not found", id)
	}

	item := items.Items[0]
	if item.Type_ != "Series" {
		log.Error().
			Str("showID", id).
			Str("actualType", item.Type_).
			Msg("Item with specified ID is not a TV show")
		return &models.MediaItem[*types.Series]{}, fmt.Errorf("item with ID %s is not a TV show", id)
	}

	itemT, err := GetItem[*types.Series](ctx, e, &item)
	if err != nil {
		return &models.MediaItem[*types.Series]{}, err
	}
	show, err := GetMediaItem[*types.Series](ctx, e, itemT, item.Id)
	if err != nil {
		return &models.MediaItem[*types.Series]{}, err
	}

	return show, nil
}

// GetSeriesSeasons retrieves seasons for a TV show
func (e *EmbyClient) GetSeriesSeasons(ctx context.Context, showID string) ([]*models.MediaItem[*types.Season], error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.GetClientID()).
		Str("clientType", string(e.GetClientType())).
		Str("showID", showID).
		Msg("Retrieving seasons for TV show from Emby server")

	opts := embyclient.TvShowsServiceApiGetShowsByIdSeasonsOpts{
		EnableImages:   optional.NewBool(true),
		EnableUserData: optional.NewBool(true),
	}

	result, resp, err := e.client.TvShowsServiceApi.GetShowsByIdSeasons(ctx, e.embyConfig().UserID, showID, &opts)
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", e.embyConfig().GetBaseURL()).
			Str("apiEndpoint", "/Shows/"+showID+"/Seasons").
			Str("showID", showID).
			Msg("Failed to fetch seasons for TV show from Emby")
		return nil, fmt.Errorf("failed to fetch seasons: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("seasonCount", len(result.Items)).
		Str("showID", showID).
		Msg("Successfully retrieved seasons for TV show from Emby")

	seasons := make([]*models.MediaItem[*types.Season], 0)
	for _, item := range result.Items {
		if item.Type_ == "Season" {
			itemSeason, err := GetItem[*types.Season](ctx, e, &item)
			season, err := GetMediaItem[*types.Season](ctx, e, itemSeason, item.Id)
			if err != nil {
				log.Warn().
					Err(err).
					Str("seasonID", item.Id).
					Str("seasonName", item.Name).
					Msg("Error converting Emby item to season format")
				continue
			}
			seasons = append(seasons, season)
		}
	}

	return seasons, nil
}

// GetSeriesEpisodes retrieves episodes for a season
func (e *EmbyClient) GetSeriesEpisodes(ctx context.Context, showID string, seasonNumber int) ([]*models.MediaItem[*types.Episode], error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.GetClientID()).
		Str("clientType", string(e.GetClientType())).
		Str("showID", showID).
		Int("seasonNumber", seasonNumber).
		Msg("Retrieving episodes for TV show season from Emby server")

	queryParams := embyclient.TvShowsServiceApiGetShowsByIdEpisodesOpts{
		IncludeItemTypes: optional.NewString("Episode"),
		Recursive:        optional.NewBool(true),
	}

	items, _, err := e.client.TvShowsServiceApi.GetShowsByIdEpisodes(ctx, showID, &queryParams)
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", e.embyConfig().GetBaseURL()).
			Str("apiEndpoint", "/Shows/"+showID+"/Episodes").
			Str("showID", showID).
			Int("seasonNumber", seasonNumber).
			Msg("Failed to fetch episodes for TV show season from Emby")
		return nil, fmt.Errorf("failed to fetch episodes: %w", err)
	}

	mediaItemEpisodes := make([]*models.MediaItem[*types.Episode], 0)
	for _, item := range items.Items {
		if item.Type_ == "Episode" && int(item.ParentIndexNumber) == seasonNumber {
			itemEpisode, err := GetItem[*types.Episode](ctx, e, &item)
			episode, err := GetMediaItem[*types.Episode](ctx, e, itemEpisode, item.Id)
			if err != nil {
				log.Warn().
					Err(err).
					Str("episodeID", item.Id).
					Str("episodeName", item.Name).
					Msg("Error converting Emby item to episode format")
				continue
			}
			if err != nil {
				log.Warn().
					Err(err).
					Str("episodeID", item.Id).
					Str("episodeName", item.Name).
					Msg("Error converting Emby item to episode format")
				continue
			}
			mediaItemEpisodes = append(mediaItemEpisodes, episode)
		}
	}

	return mediaItemEpisodes, nil
}

// GetEpisodeByID retrieves a specific episode by ID
func (e *EmbyClient) GetEpisodeByID(ctx context.Context, id string) (*models.MediaItem[*types.Episode], error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.GetClientID()).
		Str("clientType", string(e.GetClientType())).
		Str("episodeID", id).
		Msg("Retrieving specific episode from Emby server")

	queryParams := embyclient.ItemsServiceApiGetItemsOpts{
		Ids:              optional.NewString(id),
		IncludeItemTypes: optional.NewString("Episode"),
	}

	items, resp, err := e.client.ItemsServiceApi.GetItems(ctx, &queryParams)
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", e.embyConfig().GetBaseURL()).
			Str("apiEndpoint", "/Items").
			Str("episodeID", id).
			Msg("Failed to fetch episode from Emby")
		return &models.MediaItem[*types.Episode]{}, fmt.Errorf("failed to fetch episode: %w", err)
	}

	if len(items.Items) == 0 {
		log.Error().
			Str("episodeID", id).
			Int("statusCode", resp.StatusCode).
			Msg("No episode found with the specified ID")
		return &models.MediaItem[*types.Episode]{}, fmt.Errorf("episode with ID %s not found", id)
	}

	item := items.Items[0]
	if item.Type_ != "Episode" {
		log.Error().
			Str("episodeID", id).
			Str("actualType", item.Type_).
			Msg("Item with specified ID is not an episode")
		return &models.MediaItem[*types.Episode]{}, fmt.Errorf("item with ID %s is not an episode", id)
	}

	itemT, err := GetItem[*types.Episode](ctx, e, &item)
	if err != nil {
		return &models.MediaItem[*types.Episode]{}, err
	}
	episode, err := GetMediaItem[*types.Episode](ctx, e, itemT, item.Id)
	if err != nil {
		return &models.MediaItem[*types.Episode]{}, err
	}

	return episode, nil
}
