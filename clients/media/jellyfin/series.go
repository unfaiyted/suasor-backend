package jellyfin

import (
	"context"
	"fmt"

	jellyfin "github.com/sj14/jellyfin-go/api"
	"strings"
	t "suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/logger"
)

func (j *JellyfinClient) SupportsSeries() bool { return true }

func (j *JellyfinClient) GetSeries(ctx context.Context, options *t.QueryOptions) ([]*models.MediaItem[*t.Series], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.GetClientID()).
		Str("clientType", string(j.GetClientType())).
		Str("baseURL", j.config.GetBaseURL()).
		Msg("Retrieving TV shows from Jellyfin server")

	// Set up query parameters
	includeItemTypes := []jellyfin.BaseItemKind{jellyfin.BASEITEMKIND_SERIES}

	// Call the Jellyfin API
	log.Debug().Msg("Making API request to Jellyfin server for TV shows")
	itemsReq := j.client.ItemsAPI.GetItems(ctx).
		IncludeItemTypes(includeItemTypes).
		Recursive(true)

	// Set user ID first if available to ensure it's never nil
	if j.getUserID() != "" {
		itemsReq.UserId(j.getUserID())
	}

	// Then apply any additional options
	if queryOptions := NewJellyfinQueryOptions(ctx, options); queryOptions != nil {
		queryOptions.SetItemsRequest(ctx, &itemsReq)
	}

	result, resp, err := itemsReq.Execute()

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.GetBaseURL()).
			Str("apiEndpoint", "/Items").
			Int("statusCode", 0).
			Msg("Failed to fetch TV shows from Jellyfin")
		return nil, fmt.Errorf("failed to fetch TV shows: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(result.Items)).
		Int("totalRecordCount", int(*result.TotalRecordCount)).
		Msg("Successfully retrieved TV shows from Jellyfin")

	// Convert results to expected format
	shows := make([]*models.MediaItem[*t.Series], 0)
	for _, item := range result.Items {
		if *item.Type == "Series" {
			itemSeries, err := GetItem[*t.Series](ctx, j, &item)
			series, err := GetMediaItem[*t.Series](ctx, j, itemSeries, *item.Id)
			if err != nil {
				// Log error but continue
				log.Warn().
					Err(err).
					Str("showID", *item.Id).
					Str("showName", *item.Name.Get()).
					Msg("Error converting Jellyfin item to TV show format")
				continue
			}
			shows = append(shows, series)
		}
	}

	log.Info().
		Int("showsReturned", len(shows)).
		Msg("Completed GetSeriess request")

	return shows, nil
}

func (j *JellyfinClient) GetSeriesByID(ctx context.Context, id string) (models.MediaItem[*t.Series], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.GetClientID()).
		Str("clientType", string(j.GetClientType())).
		Str("showID", id).
		Str("baseURL", j.config.GetBaseURL()).
		Msg("Retrieving specific TV show from Jellyfin server")

	// Set up query parameters
	ids := id

	// Call the Jellyfin API
	log.Debug().
		Str("showID", id).
		Msg("Making API request to Jellyfin server")

	itemsReq := j.client.ItemsAPI.GetItems(ctx).Ids(strings.Split(ids, ","))

	// Set user ID if available
	if j.getUserID() != "" {
		itemsReq.UserId(j.getUserID())
	}

	result, resp, err := itemsReq.Execute()

	log.Debug().
		Interface("responseItems", result.Items).
		Msg("Full response data from Jellyfin API")

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.GetBaseURL()).
			Str("apiEndpoint", "/Items").
			Str("showID", id).
			Int("statusCode", 0).
			Msg("Failed to fetch TV show from Jellyfin")
		return models.MediaItem[*t.Series]{}, fmt.Errorf("failed to fetch TV show: %w", err)
	}

	// Check if any items were returned
	if len(result.Items) == 0 {
		log.Error().
			Str("showID", id).
			Int("statusCode", resp.StatusCode).
			Msg("No TV show found with the specified ID")
		return models.MediaItem[*t.Series]{}, fmt.Errorf("TV show with ID %s not found", id)
	}

	item := result.Items[0]

	// Double-check that the returned item is a TV show
	if *item.Type != "Series" {
		log.Error().
			Str("showID", id).
			Str("actualType", string(*item.Type.Ptr())).
			Msg("Item with specified ID is not a TV show")
		return models.MediaItem[*t.Series]{}, fmt.Errorf("item with ID %s is not a TV show", id)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Str("showID", id).
		Str("showName", *item.Name.Get()).
		Msg("Successfully retrieved TV show from Jellyfin")

	itemSeries, err := GetItem[*t.Series](ctx, j, &item)
	series, err := GetMediaItem[*t.Series](ctx, j, itemSeries, *item.Id)
	if err != nil {
		log.Error().
			Err(err).
			Str("showID", id).
			Str("showName", *item.Name.Get()).
			Msg("Error converting Jellyfin item to TV show format")
		return models.MediaItem[*t.Series]{}, fmt.Errorf("error converting TV show data: %w", err)
	}

	log.Debug().
		Str("showID", id).
		Str("showName", series.Data.Details.Title).
		Int("seasonCount", series.Data.SeasonCount).
		Msg("Successfully returned TV show data")

	return *series, nil
}

func (j *JellyfinClient) GetSeriesSeasons(ctx context.Context, showID string) ([]models.MediaItem[*t.Season], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.GetClientID()).
		Str("clientType", string(j.GetClientType())).
		Str("showID", showID).
		Str("baseURL", j.config.GetBaseURL()).
		Msg("Retrieving seasons for TV show from Jellyfin server")

	// Call the Jellyfin API
	log.Debug().
		Str("showID", showID).
		Msg("Making API request to Jellyfin server for TV show seasons")

	seasonsReq := j.client.TvShowsAPI.GetSeasons(ctx, showID).
		EnableImages(true).
		EnableUserData(true)

	// Set user ID if available
	if j.getUserID() != "" {
		seasonsReq.UserId(j.getUserID())
	}
	result, resp, err := seasonsReq.Execute()

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.GetBaseURL()).
			Str("apiEndpoint", "/Shows/"+showID+"/Seasons").
			Str("showID", showID).
			Int("statusCode", 0).
			Msg("Failed to fetch seasons for TV show from Jellyfin")
		return nil, fmt.Errorf("failed to fetch seasons: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("seasonCount", len(result.Items)).
		Str("showID", showID).
		Msg("Successfully retrieved seasons for TV show from Jellyfin")

	seasons := make([]models.MediaItem[*t.Season], 0)
	for _, item := range result.Items {
		if *item.Type == "Season" {
			itemSeason, err := GetItem[*t.Season](ctx, j, &item)
			season, err := GetMediaItem[*t.Season](ctx, j, itemSeason, *item.Id)
			if err != nil {
				// Log error but continue
				log.Warn().
					Err(err).
					Str("seasonID", *item.Id).
					Str("showID", showID).
					Msg("Error converting Jellyfin item to season format")
				continue
			}
			seasons = append(seasons, *season)
		}
	}

	log.Info().
		Int("seasonsReturned", len(seasons)).
		Str("showID", showID).
		Msg("Completed GetSeriesSeasons request")

	return seasons, nil
}

func (j *JellyfinClient) GetSeriesEpisodes(ctx context.Context, showID string, seasonNumber int) ([]*models.MediaItem[*t.Episode], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.GetClientID()).
		Str("clientType", string(j.GetClientType())).
		Str("showID", showID).
		Int("seasonNumber", seasonNumber).
		Str("baseURL", j.config.GetBaseURL()).
		Msg("Retrieving episodes for TV show season from Jellyfin server")

	seasonNum := int32(seasonNumber)

	// Call the Jellyfin API
	log.Debug().
		Str("showID", showID).
		Int("seasonNumber", seasonNumber).
		Msg("Making API request to Jellyfin server for TV show episodes")

	episodesReq := j.client.TvShowsAPI.GetEpisodes(ctx, showID).Season(seasonNum)

	// Set user ID if available
	if j.getUserID() != "" {
		episodesReq.UserId(j.getUserID())
	}
	result, resp, err := episodesReq.Execute()

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.GetBaseURL()).
			Str("apiEndpoint", "/Shows/"+showID+"/Episodes").
			Str("showID", showID).
			Int("seasonNumber", seasonNumber).
			Int("statusCode", 0).
			Msg("Failed to fetch episodes for TV show season from Jellyfin")
		return nil, fmt.Errorf("failed to fetch episodes: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("episodeCount", len(result.Items)).
		Str("showID", showID).
		Int("seasonNumber", seasonNumber).
		Msg("Successfully retrieved episodes for TV show season from Jellyfin")

	episodes := make([]*models.MediaItem[*t.Episode], 0)
	for _, item := range result.Items {
		if *item.Type == "Episode" {
			itemEpisode, err := GetItem[*t.Episode](ctx, j, &item)
			episode, err := GetMediaItem[*t.Episode](ctx, j, itemEpisode, *item.Id)
			if err != nil {
				// Log error but continue
				log.Warn().
					Err(err).
					Str("episodeID", *item.Id).
					Str("showID", showID).
					Int("seasonNumber", seasonNumber).
					Msg("Error converting Jellyfin item to episode format")
				continue
			}
			episodes = append(episodes, episode)
		}
	}

	log.Info().
		Int("episodesReturned", len(episodes)).
		Str("showID", showID).
		Int("seasonNumber", seasonNumber).
		Msg("Completed GetSeriesEpisodes request")

	return episodes, nil
}

func (j *JellyfinClient) GetEpisodeByID(ctx context.Context, id string) (*models.MediaItem[*t.Episode], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.GetClientID()).
		Str("clientType", string(j.GetClientType())).
		Str("episodeID", id).
		Str("baseURL", j.config.GetBaseURL()).
		Msg("Retrieving specific episode from Jellyfin server")

	// Set up query parameters
	ids := id

	// Call the Jellyfin API
	log.Debug().
		Str("episodeID", id).
		Msg("Making API request to Jellyfin server")

	itemsReq := j.client.ItemsAPI.GetItems(ctx).Ids(strings.Split(ids, ","))

	// Set user ID if available
	if j.getUserID() != "" {
		itemsReq.UserId(j.getUserID())
	}

	result, resp, err := itemsReq.Execute()

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.GetBaseURL()).
			Str("apiEndpoint", "/Items").
			Str("episodeID", id).
			Int("statusCode", 0).
			Msg("Failed to fetch episode from Jellyfin")
		return nil, fmt.Errorf("failed to fetch episode: %w", err)
	}

	// Check if any items were returned
	if len(result.Items) == 0 {
		log.Error().
			Str("episodeID", id).
			Int("statusCode", resp.StatusCode).
			Msg("No episode found with the specified ID")
		return nil, fmt.Errorf("episode with ID %s not found", id)
	}

	item := result.Items[0]

	// Double-check that the returned item is an episode
	if *item.Type != jellyfin.BASEITEMKIND_EPISODE {
		log.Error().
			Str("episodeID", id).
			Str("actualType", string(*item.Type)).
			Msg("Item with specified ID is not an episode")
		return nil, fmt.Errorf("item with ID %s is not an episode", id)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Str("episodeID", id).
		Str("episodeName", *item.Name.Get()).
		Msg("Successfully retrieved episode from Jellyfin")

	itemEpisode, err := GetItem[*t.Episode](ctx, j, &item)
	episode, err := GetMediaItem[*t.Episode](ctx, j, itemEpisode, *item.Id)
	if err != nil {
		log.Error().
			Err(err).
			Str("episodeID", id).
			Msg("Error converting Jellyfin item to episode format")
		return nil, fmt.Errorf("error converting episode data: %w", err)
	}

	log.Debug().
		Str("episodeID", id).
		Str("episodeName", episode.Data.Details.Title).
		Int64("episodeNumber", episode.Data.Number).
		Int("seasonNumber", episode.Data.SeasonNumber).
		Msg("Successfully returned episode data")

	return episode, nil
}
