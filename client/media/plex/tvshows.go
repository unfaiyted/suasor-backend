package plex

import (
	"context"
	"fmt"
	"strconv"
	"suasor/client/media/types"
	"suasor/types/models"
	"suasor/utils"
	"time"

	"github.com/LukeHagar/plexgo"
	"github.com/LukeHagar/plexgo/models/operations"
)

// GetTVShows retrieves TV shows from Plex
func (c *PlexClient) GetTVShows(ctx context.Context, options *types.QueryOptions) ([]models.MediaItem[types.TVShow], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("baseURL", c.baseURL).
		Msg("Retrieving TV shows from Plex server")

	// First, find the TV show library section
	log.Debug().Msg("Finding TV show library section")
	tvSectionKey, err := c.findLibrarySectionByType(ctx, "show")
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("Failed to find TV show library section")
		return nil, err
	}

	if tvSectionKey == "" {
		log.Info().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("No TV show library section found in Plex")
		return []models.MediaItem[types.TVShow]{}, nil
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
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Int("sectionKey", sectionKey).
			Msg("Failed to get TV shows from Plex")
		return nil, fmt.Errorf("failed to get TV shows: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil {
		log.Info().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("No TV shows found in Plex")
		return []models.MediaItem[types.TVShow]{}, nil
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("totalItems", len(res.Object.MediaContainer.Metadata)).
		Msg("Successfully retrieved TV shows from Plex")

	shows := make([]models.MediaItem[types.TVShow], 0, len(res.Object.MediaContainer.Metadata))
	for _, item := range res.Object.MediaContainer.Metadata {
		if item.Type != "show" {
			continue
		}

		show := models.MediaItem[types.TVShow]{
			Data: types.TVShow{Details: c.createMetadataFromPlexItem(&item)},
		}

		show.SetClientInfo(c.ClientID, c.ClientType, item.RatingKey)

		if item.Rating != nil {
			show.Data.Rating = float64(*item.Rating)
		}
		if item.Year != nil {
			show.Data.ReleaseYear = *item.Year
		}
		if item.ContentRating != nil {
			show.Data.ContentRating = *item.ContentRating
		}
		if item.ChildCount != nil {
			show.Data.SeasonCount = *item.ChildCount
		}
		if item.LeafCount != nil {
			show.Data.EpisodeCount = int(*item.LeafCount)
		}

		if item.Genre != nil {
			show.Data.Genres = make([]string, 0, len(item.Genre))
			for _, genre := range item.Genre {
				if genre.Tag != nil {
					show.Data.Genres = append(show.Data.Genres, *genre.Tag)
				}
			}
		}

		shows = append(shows, show)

		log.Debug().
			Str("showID", item.RatingKey).
			Str("showTitle", item.Title).
			Int("seasonCount", show.Data.SeasonCount).
			Msg("Added TV show to result list")
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("showsReturned", len(shows)).
		Msg("Completed GetTVShows request")

	return shows, nil
}

// GetTVShowSeasons retrieves seasons for a specific TV show
func (c *PlexClient) GetTVShowSeasons(ctx context.Context, showID string) ([]models.MediaItem[types.Season], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("showID", showID).
		Str("baseURL", c.baseURL).
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
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("showID", showID).
			Msg("Failed to get TV show seasons from Plex")
		return nil, fmt.Errorf("failed to get TV show seasons: %w", err)
	}

	if childRes.Object.MediaContainer == nil || childRes.Object.MediaContainer.Metadata == nil {
		log.Info().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("showID", showID).
			Msg("No seasons found for TV show in Plex")
		return []models.MediaItem[types.Season]{}, nil
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("showID", showID).
		Int("totalItems", len(childRes.Object.MediaContainer.Metadata)).
		Msg("Successfully retrieved seasons for TV show from Plex")

	seasons := make([]models.MediaItem[types.Season], 0, len(childRes.Object.MediaContainer.Metadata))
	for _, item := range childRes.Object.MediaContainer.Metadata {
		if *item.Type != "season" {
			continue
		}

		season := models.MediaItem[types.Season]{
			ExternalID: *item.RatingKey,
			Data: types.Season{
				EpisodeCount: *item.LeafCount,
				Number:       *item.Index,
				Details: types.MediaMetadata{
					Description: *item.Summary,
					Title:       *item.Title,
					Artwork: types.Artwork{
						Thumbnail: c.makeFullURL(*item.Thumb),
					},
					ExternalIDs: types.ExternalIDs{types.ExternalID{
						Source: "plex",
						ID:     *item.RatingKey,
					}},
					UpdatedAt: time.Unix(int64(*item.UpdatedAt), 0),
					AddedAt:   time.Unix(int64(*item.AddedAt), 0),
				},
			},
		}

		season.SetClientInfo(c.ClientID, c.ClientType, *item.RatingKey)

		seasons = append(seasons, season)

		log.Debug().
			Str("seasonID", *item.RatingKey).
			Str("seasonTitle", *item.Title).
			Int("seasonNumber", *item.Index).
			Int("episodeCount", *item.LeafCount).
			Msg("Added season to result list")
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("showID", showID).
		Int("seasonsReturned", len(seasons)).
		Msg("Completed GetTVShowSeasons request")

	return seasons, nil
}

// GetTVShowEpisodes retrieves episodes for a specific season of a TV show
func (c *PlexClient) GetTVShowEpisodes(ctx context.Context, showID string, seasonNumber int) ([]models.MediaItem[types.Episode], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("showID", showID).
		Int("seasonNumber", seasonNumber).
		Str("baseURL", c.baseURL).
		Msg("Retrieving episodes for TV show season from Plex server")

	// First get all seasons
	log.Debug().
		Str("showID", showID).
		Int("seasonNumber", seasonNumber).
		Msg("Getting seasons for the TV show")

	seasons, err := c.GetTVShowSeasons(ctx, showID)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
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
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("showID", showID).
			Int("seasonNumber", seasonNumber).
			Msg("Season not found for TV show in Plex")
		return []models.MediaItem[types.Episode]{}, nil
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
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("showID", showID).
			Int("seasonNumber", seasonNumber).
			Str("seasonID", seasonID).
			Msg("Failed to get TV show episodes from Plex")
		return nil, fmt.Errorf("failed to get TV show episodes: %w", err)
	}

	if childRes.Object.MediaContainer == nil || childRes.Object.MediaContainer.Metadata == nil {
		log.Info().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("showID", showID).
			Int("seasonNumber", seasonNumber).
			Msg("No episodes found for TV show season in Plex")
		return []models.MediaItem[types.Episode]{}, nil
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("showID", showID).
		Int("seasonNumber", seasonNumber).
		Int("totalItems", len(childRes.Object.MediaContainer.Metadata)).
		Msg("Successfully retrieved episodes for TV show season from Plex")

	episodes := make([]models.MediaItem[types.Episode], 0, len(childRes.Object.MediaContainer.Metadata))
	for _, item := range childRes.Object.MediaContainer.Metadata {
		if *item.Type != "episode" {
			continue
		}

		episode := models.MediaItem[types.Episode]{

			ExternalID: *item.RatingKey,
			Data: types.Episode{
				ShowID:       showID,
				Number:       int64(*item.Index),
				SeasonNumber: int(*item.ParentIndex),
				SeasonID:     *item.ParentKey,
				Details: types.MediaMetadata{
					Description: *item.Summary,
					Title:       *item.Title,
					Artwork: types.Artwork{
						Thumbnail: c.makeFullURL(*item.Thumb),
					},
					UpdatedAt: time.Unix(int64(*item.UpdatedAt), 0),
					AddedAt:   time.Unix(int64(*item.AddedAt), 0),
				},
			},
		}
		episode.SetClientInfo(c.ClientID, c.ClientType, *item.RatingKey)

		// Add studio if available
		if item.ParentStudio != nil {
			episode.Data.Details.Studios = []string{*item.ParentStudio}
		}

		episodes = append(episodes, episode)

		log.Debug().
			Str("episodeID", *item.RatingKey).
			Str("episodeTitle", *item.Title).
			Int("seasonNumber", episode.Data.SeasonNumber).
			Int64("episodeNumber", episode.Data.Number).
			Msg("Added episode to result list")
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("showID", showID).
		Int("seasonNumber", seasonNumber).
		Int("episodesReturned", len(episodes)).
		Msg("Completed GetTVShowEpisodes request")

	return episodes, nil
}

// GetTVShowByID retrieves a specific TV show by ID
func (c *PlexClient) GetTVShowByID(ctx context.Context, id string) (models.MediaItem[types.TVShow], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("showID", id).
		Str("baseURL", c.baseURL).
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
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("showID", id).
			Msg("Failed to get TV show from Plex")
		return models.MediaItem[types.TVShow]{}, fmt.Errorf("failed to get TV show: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil || len(res.Object.MediaContainer.Metadata) == 0 {
		log.Error().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("showID", id).
			Msg("TV show not found in Plex")
		return models.MediaItem[types.TVShow]{}, fmt.Errorf("TV show not found")
	}

	item := res.Object.MediaContainer.Metadata[0]
	if item.Type != "show" {
		log.Error().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("showID", id).
			Str("actualType", item.Type).
			Msg("Item retrieved is not a TV show")
		return models.MediaItem[types.TVShow]{}, fmt.Errorf("item is not a TV show")
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("showID", id).
		Str("showTitle", item.Title).
		Msg("Successfully retrieved TV show from Plex")

	show := models.MediaItem[types.TVShow]{
		Data: types.TVShow{
			Details: c.createMediaMetadataFromPlexItem(&item),
		},
	}
	show.SetClientInfo(c.ClientID, c.ClientType, item.RatingKey)

	if item.Rating != nil {
		show.Data.Rating = float64(*item.Rating)
	}
	if item.ContentRating != nil {
		show.Data.ContentRating = *item.ContentRating
	}
	if item.ChildCount != nil {
		show.Data.SeasonCount = *item.ChildCount
	}
	if item.LeafCount != nil {
		show.Data.EpisodeCount = int(*item.LeafCount)
	}

	if item.Genre != nil {
		show.Data.Genres = make([]string, 0, len(item.Genre))
		for _, genre := range item.Genre {
			show.Data.Genres = append(show.Data.Genres, genre.Tag)
		}
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("showID", id).
		Str("showTitle", show.Data.Details.Title).
		Int("seasonCount", show.Data.SeasonCount).
		Msg("Successfully converted TV show data")

	return show, nil
}

// GetEpisodeByID retrieves a specific episode by ID
func (c *PlexClient) GetEpisodeByID(ctx context.Context, id string) (models.MediaItem[types.Episode], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("episodeID", id).
		Msg("Retrieving specific episode from Plex server")

	ratingKey, _ := strconv.Atoi(id)
	int64RatingKey := int64(ratingKey)

	res, err := c.plexAPI.Library.GetMediaMetaData(ctx, operations.GetMediaMetaDataRequest{
		RatingKey: int64RatingKey,
	})
	if err != nil {
		log.Error().Err(err).Str("episodeID", id).Msg("Failed to get episode from Plex")
		return models.MediaItem[types.Episode]{}, fmt.Errorf("failed to get episode: %w", err)
	}

	if res.Object.MediaContainer == nil ||
		res.Object.MediaContainer.Metadata == nil ||
		len(res.Object.MediaContainer.Metadata) == 0 {
		log.Error().Str("episodeID", id).Msg("Episode not found in Plex")
		return models.MediaItem[types.Episode]{}, fmt.Errorf("episode not found")
	}

	item := res.Object.MediaContainer.Metadata[0]
	if item.Type != "episode" {
		log.Error().Str("episodeID", id).Str("actualType", item.Type).Msg("Item retrieved is not an episode")
		return models.MediaItem[types.Episode]{}, fmt.Errorf("item is not an episode")
	}

	episode := models.MediaItem[types.Episode]{
		Data: types.Episode{
			Details: c.createMediaMetadataFromPlexItem(&item),
			Number:  int64(*item.Index),
		},
	}

	episode.SetClientInfo(c.ClientID, c.ClientType, item.RatingKey)

	// Add season number if available
	if item.ParentIndex != nil {
		episode.Data.SeasonNumber = int(*item.ParentIndex)
	}

	// Add show ID if available (via grandparentRatingKey)
	if item.GrandparentRatingKey != nil {
		episode.Data.ShowID = *item.GrandparentRatingKey
	}

	// Add studio if available
	if item.Studio != nil {
		episode.Data.Details.Studios = []string{*item.Studio}
	}

	log.Info().
		Str("episodeID", id).
		Str("episodeTitle", episode.Data.Details.Title).
		Int("seasonNumber", episode.Data.SeasonNumber).
		Int64("episodeNumber", episode.Data.Number).
		Msg("Successfully retrieved episode")

	return episode, nil
}
