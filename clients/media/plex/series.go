package plex

import (
	"context"
	"fmt"
	"strconv"
	"suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/logger"

	"github.com/LukeHagar/plexgo"
	"github.com/LukeHagar/plexgo/models/operations"
)

// GetSeriess retrieves TV shows from Plex
func (c *PlexClient) GetSeries(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Series], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Msg("Retrieving TV shows from Plex server")

	// First, find the TV show library section
	log.Debug().Msg("Finding TV show library section")
	tvSectionKey, err := c.findLibrarySectionByType(ctx, "show")
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Msg("Failed to find TV show library section")
		return nil, err
	}

	if tvSectionKey == "" {
		log.Info().
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Msg("No TV show library section found in Plex")
		return nil, nil
	}

	// Get TV shows from the TV section
	sectionKey, _ := strconv.Atoi(tvSectionKey)
	log.Debug().
		Int("sectionKey", sectionKey).
		Msg("Making API request to Plex server for TV shows")

	res, err := c.plexAPI.Library.GetLibraryItems(ctx, operations.GetLibraryItemsRequest{
		IncludeMeta: operations.GetLibraryItemsQueryParamIncludeMetaEnable.ToPointer(),
		Tag:         "all",
		Type:        operations.GetLibraryItemsQueryParamTypeTvShow,
		SectionKey:  sectionKey,
	})

	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Int("sectionKey", sectionKey).
			Msg("Failed to get TV shows from Plex")
		return nil, fmt.Errorf("failed to get TV shows: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil {
		log.Info().
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Msg("No TV shows found in Plex")
		return nil, nil
	}

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Int("totalItems", len(res.Object.MediaContainer.Metadata)).
		Msg("Successfully retrieved TV shows from Plex")

	series, err := GetMediaItemList[*types.Series](ctx, c, res.Object.MediaContainer.Metadata)

	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Msg("Failed to get TV shows from Plex")
		return nil, fmt.Errorf("failed to get TV shows: %w", err)
	}

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Int("showsReturned", len(series)).
		Msg("Completed GetSeries request")

	return series, nil
}

// GetSeriesSeasons retrieves seasons for a specific TV show
func (c *PlexClient) GetSeriesSeasons(ctx context.Context, showID string) ([]*models.MediaItem[*types.Season], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("showID", showID).
		Msg("Retrieving seasons for TV show from Plex server")

	ratingKey, _ := strconv.Atoi(showID)
	float64RatingKey := float64(ratingKey)

	log.Debug().
		Str("showID", showID).
		Float64("ratingKey", float64RatingKey).
		Msg("Making API request to Plex server for TV show seasons")

	childRes, err := c.plexAPI.Library.GetMetadataChildren(ctx, float64RatingKey, plexgo.String("Stream"))
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Str("showID", showID).
			Msg("Failed to get TV show seasons from Plex")
		return nil, fmt.Errorf("failed to get TV show seasons: %w", err)
	}

	if childRes.Object.MediaContainer == nil || childRes.Object.MediaContainer.Metadata == nil {
		log.Info().
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Str("showID", showID).
			Msg("No seasons found for TV show in Plex")
		return nil, nil
	}

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("showID", showID).
		Int("totalItems", len(childRes.Object.MediaContainer.Metadata)).
		Msg("Successfully retrieved seasons for TV show from Plex")

	seasons, err := GetChildMediaItemsList[*types.Season](ctx, c, childRes.Object.MediaContainer.Metadata)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("showID", showID).
		Int("seasonsReturned", len(seasons)).
		Msg("Completed GetSeriesSeasons request")

	return seasons, nil
}

// GetSeriesEpisodes retrieves episodes for a specific season of a TV show
func (c *PlexClient) GetSeriesEpisodes(ctx context.Context, showID string, seasonNumber int) ([]*models.MediaItem[*types.Episode], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("showID", showID).
		Int("seasonNumber", seasonNumber).
		Msg("Retrieving episodes for TV show season from Plex server")

	// First get all seasons
	log.Debug().
		Str("showID", showID).
		Int("seasonNumber", seasonNumber).
		Msg("Getting seasons for the TV show")

	seasons, err := c.GetSeriesSeasons(ctx, showID)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Str("showID", showID).
			Int("seasonNumber", seasonNumber).
			Msg("Failed to get seasons for TV show")
		return nil, err
	}

	var seasonID string
	for _, season := range seasons {
		if season.Data.Number == seasonNumber {
			for _, externalID := range season.Data.Details.ExternalIDs {
				if externalID.Source == "plex" {
					seasonID = externalID.ID
					break
				}
			}
			break
		}
	}

	if seasonID == "" {
		log.Info().
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Str("showID", showID).
			Int("seasonNumber", seasonNumber).
			Msg("Season not found for TV show in Plex")
		return nil, nil
	}

	ratingKey, _ := strconv.Atoi(seasonID)
	float64RatingKey := float64(ratingKey)

	log.Debug().
		Str("seasonID", seasonID).
		Float64("ratingKey", float64RatingKey).
		Msg("Making API request to Plex server for TV show episodes")

	childRes, err := c.plexAPI.Library.GetMetadataChildren(ctx, float64RatingKey, plexgo.String("Stream"))
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Str("showID", showID).
			Int("seasonNumber", seasonNumber).
			Str("seasonID", seasonID).
			Msg("Failed to get TV show episodes from Plex")
		return nil, fmt.Errorf("failed to get TV show episodes: %w", err)
	}

	if childRes.Object.MediaContainer == nil || childRes.Object.MediaContainer.Metadata == nil {
		log.Info().
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Str("showID", showID).
			Int("seasonNumber", seasonNumber).
			Msg("No episodes found for TV show season in Plex")
		return nil, nil
	}

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("showID", showID).
		Int("seasonNumber", seasonNumber).
		Int("totalItems", len(childRes.Object.MediaContainer.Metadata)).
		Msg("Successfully retrieved episodes for TV show season from Plex")

	//GetMetadataChildrenMetadata
	episodes, err := GetChildMediaItemsList[*types.Episode](ctx, c, childRes.Object.MediaContainer.Metadata)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("showID", showID).
		Int("seasonNumber", seasonNumber).
		Int("episodesReturned", len(episodes)).
		Msg("Completed GetSeriesEpisodes request")

	return episodes, nil
}

// GetSeriesByID retrieves a specific TV show by ID
func (c *PlexClient) GetSeriesByID(ctx context.Context, id string) (*models.MediaItem[*types.Series], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("showID", id).
		Msg("Retrieving specific TV show from Plex server")

	ratingKey, _ := strconv.Atoi(id)
	int64RatingKey := int64(ratingKey)

	log.Debug().
		Str("showID", id).
		Int64("ratingKey", int64RatingKey).
		Msg("Making API request to Plex server for TV show")

	res, err := c.plexAPI.Library.GetMediaMetaData(ctx, operations.GetMediaMetaDataRequest{
		RatingKey: int64RatingKey,
	})

	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Str("showID", id).
			Msg("Failed to get TV show from Plex")
		return nil, fmt.Errorf("failed to get TV show: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil || len(res.Object.MediaContainer.Metadata) == 0 {
		log.Error().
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Str("showID", id).
			Msg("TV show not found in Plex")
		return nil, fmt.Errorf("TV show not found")
	}

	item := res.Object.MediaContainer.Metadata[0]
	if item.Type != "show" {
		log.Error().
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Str("showID", id).
			Str("actualType", item.Type).
			Msg("Item retrieved is not a TV show")
		return nil, fmt.Errorf("item is not a TV show")
	}

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("showID", id).
		Str("showTitle", item.Title).
		Msg("Successfully retrieved TV show from Plex")

	itemSeries, err := GetItemFromMetadata[*types.Series](ctx, c, &item)
	series, err := GetMediaItem[*types.Series](ctx, c, itemSeries, item.RatingKey)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("showID", id).
		Str("showTitle", series.Data.Details.Title).
		Int("seasonCount", series.Data.SeasonCount).
		Msg("Successfully converted TV show data")

	return series, nil
}

// GetEpisodeByID retrieves a specific episode by ID
func (c *PlexClient) GetEpisodeByID(ctx context.Context, id string) (*models.MediaItem[*types.Episode], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("episodeID", id).
		Msg("Retrieving specific episode from Plex server")

	ratingKey, _ := strconv.Atoi(id)
	int64RatingKey := int64(ratingKey)

	res, err := c.plexAPI.Library.GetMediaMetaData(ctx, operations.GetMediaMetaDataRequest{
		RatingKey: int64RatingKey,
	})
	if err != nil {
		log.Error().Err(err).Str("episodeID", id).Msg("Failed to get episode from Plex")
		return nil, fmt.Errorf("failed to get episode: %w", err)
	}

	if res.Object.MediaContainer == nil ||
		res.Object.MediaContainer.Metadata == nil ||
		len(res.Object.MediaContainer.Metadata) == 0 {
		log.Error().Str("episodeID", id).Msg("Episode not found in Plex")
		return nil, fmt.Errorf("episode not found")
	}

	item := res.Object.MediaContainer.Metadata[0]
	if item.Type != "episode" {
		log.Error().Str("episodeID", id).Str("actualType", item.Type).Msg("Item retrieved is not an episode")
		return nil, fmt.Errorf("item is not an episode")
	}

	itemEpisode, err := GetItemFromMetadata[*types.Episode](ctx, c, &item)
	episode, err := GetMediaItem[*types.Episode](ctx, c, itemEpisode, item.RatingKey)

	log.Info().
		Str("episodeID", id).
		Str("episodeTitle", episode.Data.Details.Title).
		Int("seasonNumber", episode.Data.SeasonNumber).
		Int64("episodeNumber", episode.Data.Number).
		Msg("Successfully retrieved episode")

	return episode, nil
}
