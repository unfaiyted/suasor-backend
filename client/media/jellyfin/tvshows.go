package jellyfin

import (
	"context"
	"fmt"

	jellyfin "github.com/sj14/jellyfin-go/api"
	t "suasor/client/media/types"
	"suasor/utils"
)

func (j *JellyfinClient) GetTVShows(ctx context.Context, options *t.QueryOptions) ([]t.MediaItem[t.TVShow], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.ClientID).
		Str("clientType", string(j.ClientType)).
		Str("baseURL", j.config.BaseURL).
		Msg("Retrieving TV shows from Jellyfin server")

	// Set up query parameters
	includeItemTypes := []jellyfin.BaseItemKind{jellyfin.BASEITEMKIND_SERIES}

	limit, startIndex, sortBy, sortOrder := j.getQueryParameters(options)

	// Call the Jellyfin API
	log.Debug().Msg("Making API request to Jellyfin server for TV shows")
	itemsReq := j.client.ItemsAPI.GetItems(ctx).
		IncludeItemTypes(includeItemTypes).
		Recursive(true).
		Limit(*limit).
		StartIndex(*startIndex).
		SortBy(sortBy).
		SortOrder(sortOrder)

	result, resp, err := itemsReq.Execute()

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.BaseURL).
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
	shows := make([]t.MediaItem[t.TVShow], 0)
	for _, item := range result.Items {
		if *item.Type == "Series" {
			show, err := j.convertToTVShow(ctx, &item)
			if err != nil {
				// Log error but continue
				log.Warn().
					Err(err).
					Str("showID", *item.Id).
					Str("showName", *item.Name.Get()).
					Msg("Error converting Jellyfin item to TV show format")
				continue
			}
			shows = append(shows, show)
		}
	}

	log.Info().
		Int("showsReturned", len(shows)).
		Msg("Completed GetTVShows request")

	return shows, nil
}

func (j *JellyfinClient) GetTVShowByID(ctx context.Context, id string) (t.MediaItem[t.TVShow], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.ClientID).
		Str("clientType", string(j.ClientType)).
		Str("showID", id).
		Str("baseURL", j.config.BaseURL).
		Msg("Retrieving specific TV show from Jellyfin server")

	// Set up query parameters
	ids := id

	// Call the Jellyfin API
	log.Debug().
		Str("showID", id).
		Msg("Making API request to Jellyfin server")

	itemsReq := j.client.ItemsAPI.GetItems(ctx).Ids(stringToSlice(ids))

	result, resp, err := itemsReq.Execute()

	log.Debug().
		Interface("responseItems", result.Items).
		Msg("Full response data from Jellyfin API")

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.BaseURL).
			Str("apiEndpoint", "/Items").
			Str("showID", id).
			Int("statusCode", 0).
			Msg("Failed to fetch TV show from Jellyfin")
		return t.MediaItem[t.TVShow]{}, fmt.Errorf("failed to fetch TV show: %w", err)
	}

	// Check if any items were returned
	if len(result.Items) == 0 {
		log.Error().
			Str("showID", id).
			Int("statusCode", resp.StatusCode).
			Msg("No TV show found with the specified ID")
		return t.MediaItem[t.TVShow]{}, fmt.Errorf("TV show with ID %s not found", id)
	}

	item := result.Items[0]

	// Double-check that the returned item is a TV show
	if *item.Type != "Series" {
		log.Error().
			Str("showID", id).
			Str("actualType", string(*item.Type.Ptr())).
			Msg("Item with specified ID is not a TV show")
		return t.MediaItem[t.TVShow]{}, fmt.Errorf("item with ID %s is not a TV show", id)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Str("showID", id).
		Str("showName", *item.Name.Get()).
		Msg("Successfully retrieved TV show from Jellyfin")

	show, err := j.convertToTVShow(ctx, &item)
	if err != nil {
		log.Error().
			Err(err).
			Str("showID", id).
			Str("showName", *item.Name.Get()).
			Msg("Error converting Jellyfin item to TV show format")
		return t.MediaItem[t.TVShow]{}, fmt.Errorf("error converting TV show data: %w", err)
	}

	log.Debug().
		Str("showID", id).
		Str("showName", show.Data.Details.Title).
		Int("seasonCount", show.Data.SeasonCount).
		Msg("Successfully returned TV show data")

	return show, nil
}

func (j *JellyfinClient) GetTVShowSeasons(ctx context.Context, showID string) ([]t.MediaItem[t.Season], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.ClientID).
		Str("clientType", string(j.ClientType)).
		Str("showID", showID).
		Str("baseURL", j.config.BaseURL).
		Msg("Retrieving seasons for TV show from Jellyfin server")

	// Call the Jellyfin API
	log.Debug().
		Str("showID", showID).
		Msg("Making API request to Jellyfin server for TV show seasons")

	seasonsReq := j.client.TvShowsAPI.GetSeasons(ctx, showID).
		EnableImages(true).
		EnableUserData(true).
		UserId(j.config.UserID)
	result, resp, err := seasonsReq.Execute()

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.BaseURL).
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

	seasons := make([]t.MediaItem[t.Season], 0)
	for _, item := range result.Items {
		if *item.Type == "Season" {
			season, err := j.convertToSeason(ctx, &item)
			if err != nil {
				// Log error but continue
				log.Warn().
					Err(err).
					Str("seasonID", *item.Id).
					Str("showID", showID).
					Msg("Error converting Jellyfin item to season format")
				continue
			}
			seasons = append(seasons, season)
		}
	}

	log.Info().
		Int("seasonsReturned", len(seasons)).
		Str("showID", showID).
		Msg("Completed GetTVShowSeasons request")

	return seasons, nil
}

func (j *JellyfinClient) GetTVShowEpisodes(ctx context.Context, showID string, seasonNumber int) ([]t.MediaItem[t.Episode], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.ClientID).
		Str("clientType", string(j.ClientType)).
		Str("showID", showID).
		Int("seasonNumber", seasonNumber).
		Str("baseURL", j.config.BaseURL).
		Msg("Retrieving episodes for TV show season from Jellyfin server")

	seasonNum := int32(seasonNumber)

	// Call the Jellyfin API
	log.Debug().
		Str("showID", showID).
		Int("seasonNumber", seasonNumber).
		Msg("Making API request to Jellyfin server for TV show episodes")

	episodesReq := j.client.TvShowsAPI.GetEpisodes(ctx, showID).Season(seasonNum).UserId(j.config.UserID)
	result, resp, err := episodesReq.Execute()

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.BaseURL).
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

	episodes := make([]t.MediaItem[t.Episode], 0)
	for _, item := range result.Items {
		if *item.Type == "Episode" {
			episode, err := j.convertToEpisode(ctx, &item)
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
		Msg("Completed GetTVShowEpisodes request")

	return episodes, nil
}

func (j *JellyfinClient) GetEpisodeByID(ctx context.Context, id string) (t.MediaItem[t.Episode], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.ClientID).
		Str("clientType", string(j.ClientType)).
		Str("episodeID", id).
		Str("baseURL", j.config.BaseURL).
		Msg("Retrieving specific episode from Jellyfin server")

	// Set up query parameters
	ids := id

	// Call the Jellyfin API
	log.Debug().
		Str("episodeID", id).
		Msg("Making API request to Jellyfin server")

	itemsReq := j.client.ItemsAPI.GetItems(ctx).Ids(stringToSlice(ids))

	result, resp, err := itemsReq.Execute()

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.BaseURL).
			Str("apiEndpoint", "/Items").
			Str("episodeID", id).
			Int("statusCode", 0).
			Msg("Failed to fetch episode from Jellyfin")
		return t.MediaItem[t.Episode]{}, fmt.Errorf("failed to fetch episode: %w", err)
	}

	// Check if any items were returned
	if len(result.Items) == 0 {
		log.Error().
			Str("episodeID", id).
			Int("statusCode", resp.StatusCode).
			Msg("No episode found with the specified ID")
		return t.MediaItem[t.Episode]{}, fmt.Errorf("episode with ID %s not found", id)
	}

	item := result.Items[0]

	// Double-check that the returned item is an episode
	if *item.Type != jellyfin.BASEITEMKIND_EPISODE {
		log.Error().
			Str("episodeID", id).
			Str("actualType", baseItemKindToString(*item.Type)).
			Msg("Item with specified ID is not an episode")
		return t.MediaItem[t.Episode]{}, fmt.Errorf("item with ID %s is not an episode", id)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Str("episodeID", id).
		Str("episodeName", *item.Name.Get()).
		Msg("Successfully retrieved episode from Jellyfin")

	episode, err := j.convertToEpisode(ctx, &item)
	if err != nil {
		log.Error().
			Err(err).
			Str("episodeID", id).
			Msg("Error converting Jellyfin item to episode format")
		return t.MediaItem[t.Episode]{}, fmt.Errorf("error converting episode data: %w", err)
	}

	log.Debug().
		Str("episodeID", id).
		Str("episodeName", episode.Data.Details.Title).
		Int64("episodeNumber", episode.Data.Number).
		Int("seasonNumber", episode.Data.SeasonNumber).
		Msg("Successfully returned episode data")

	return episode, nil
}
