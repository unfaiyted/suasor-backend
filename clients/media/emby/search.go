package emby

import (
	"context"
	"fmt"

	"github.com/antihax/optional"
	"suasor/clients/media/types"
	embyclient "suasor/internal/clients/embyAPI"
	"suasor/types/responses"
	"suasor/utils/logger"
)

func (e *EmbyClient) SupportsSearch() bool { return true }

// Search for media items in Emby
func (e *EmbyClient) Search(ctx context.Context, options *types.QueryOptions) (responses.SearchResults, error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.GetClientID()).
		Str("clientType", string(e.GetClientType())).
		Str("query", options.Query).
		Str("mediaType", string(options.MediaType)).
		Msg("Searching media items in Emby")

	// Initialize the result container
	var results responses.SearchResults

	// Get user ID
	userID := e.getUserID()
	if userID == "" {
		log.Error().Msg("User ID is required for Emby queries but was not provided or resolved")
		return results, fmt.Errorf("failed to search Emby: missing user ID")
	}

	// Create base query parameters
	queryParams := embyclient.ItemsServiceApiGetItemsOpts{
		Recursive: optional.NewBool(true),
		UserId:    optional.NewString(userID),
		Fields:    optional.NewString("PrimaryImageAspectRatio,BasicSyncInfo,CanDelete,Container,DateCreated,PremiereDate,ProductionYear,Genres,MediaSourceCount,MediaSources,Overview,ParentId,Path,SortName,Studios,Taglines"),
	}

	// Add search term
	if options.Query != "" {
		queryParams.SearchTerm = optional.NewString(options.Query)
	}

	// Apply additional options (filters, sorting, pagination)
	ApplyClientQueryOptions(ctx, &queryParams, options)

	// Filter by media type if specified, otherwise search all supported types
	if options.MediaType != "" && options.MediaType != types.MediaTypeAll {
		// Set appropriate include types for the specific media type
		switch options.MediaType {
		case types.MediaTypeMovie:
			queryParams.IncludeItemTypes = optional.NewString("Movie")
		case types.MediaTypeSeries:
			queryParams.IncludeItemTypes = optional.NewString("Series")
		case types.MediaTypeEpisode:
			queryParams.IncludeItemTypes = optional.NewString("Episode")
		case types.MediaTypeSeason:
			queryParams.IncludeItemTypes = optional.NewString("Season")
		case types.MediaTypeArtist:
			queryParams.IncludeItemTypes = optional.NewString("MusicArtist")
		case types.MediaTypeAlbum:
			queryParams.IncludeItemTypes = optional.NewString("MusicAlbum")
		case types.MediaTypeTrack:
			queryParams.IncludeItemTypes = optional.NewString("Audio")
		case types.MediaTypePlaylist:
			queryParams.IncludeItemTypes = optional.NewString("Playlist")
		case types.MediaTypeCollection:
			queryParams.IncludeItemTypes = optional.NewString("BoxSet")
		default:
			log.Warn().
				Str("mediaType", string(options.MediaType)).
				Msg("Unsupported media type for Emby search, will search all types")
		}
	} else {
		// When searching all types, include all supported media types
		queryParams.IncludeItemTypes = optional.NewString("Movie,Series,Episode,Season,MusicArtist,MusicAlbum,Audio,Playlist,BoxSet")
	}

	// Call the Emby API
	response, resp, err := e.client.ItemsServiceApi.GetItems(ctx, &queryParams)
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", e.embyConfig().GetBaseURL()).
			Str("apiEndpoint", "/Items").
			Msg("Failed to search items from Emby")
		return results, fmt.Errorf("failed to search items: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(response.Items)).
		Int("totalRecordCount", int(response.TotalRecordCount)).
		Msg("Successfully searched items from Emby")

	// Process the results
	for _, item := range response.Items {
		switch item.Type_ {
		case "Movie":
			itemMovie, err := GetItem[*types.Movie](ctx, e, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", item.Id).
					Str("itemName", item.Name).
					Msg("Error converting Emby item to movie format")
				continue
			}
			mediaItemMovie, err := GetMediaItem[*types.Movie](ctx, e, itemMovie, item.Id)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", item.Id).
					Str("itemName", item.Name).
					Msg("Error creating media item for movie")
				continue
			}
			results.Movies = append(results.Movies, mediaItemMovie)

		case "Series":
			itemSeries, err := GetItem[*types.Series](ctx, e, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", item.Id).
					Str("itemName", item.Name).
					Msg("Error converting Emby item to series format")
				continue
			}
			mediaItemSeries, err := GetMediaItem[*types.Series](ctx, e, itemSeries, item.Id)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", item.Id).
					Str("itemName", item.Name).
					Msg("Error creating media item for series")
				continue
			}
			results.Series = append(results.Series, mediaItemSeries)

		case "Episode":
			itemEpisode, err := GetItem[*types.Episode](ctx, e, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", item.Id).
					Str("itemName", item.Name).
					Msg("Error converting Emby item to episode format")
				continue
			}
			mediaItemEpisode, err := GetMediaItem[*types.Episode](ctx, e, itemEpisode, item.Id)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", item.Id).
					Str("itemName", item.Name).
					Msg("Error creating media item for episode")
				continue
			}
			results.Episodes = append(results.Episodes, mediaItemEpisode)

		case "Season":
			itemSeason, err := GetItem[*types.Season](ctx, e, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", item.Id).
					Str("itemName", item.Name).
					Msg("Error converting Emby item to season format")
				continue
			}
			mediaItemSeason, err := GetMediaItem[*types.Season](ctx, e, itemSeason, item.Id)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", item.Id).
					Str("itemName", item.Name).
					Msg("Error creating media item for season")
				continue
			}
			results.Seasons = append(results.Seasons, mediaItemSeason)

		case "MusicArtist":
			itemArtist, err := GetItem[*types.Artist](ctx, e, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", item.Id).
					Str("itemName", item.Name).
					Msg("Error converting Emby item to artist format")
				continue
			}
			mediaItemArtist, err := GetMediaItem[*types.Artist](ctx, e, itemArtist, item.Id)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", item.Id).
					Str("itemName", item.Name).
					Msg("Error creating media item for artist")
				continue
			}
			results.Artists = append(results.Artists, mediaItemArtist)

		case "MusicAlbum":
			itemAlbum, err := GetItem[*types.Album](ctx, e, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", item.Id).
					Str("itemName", item.Name).
					Msg("Error converting Emby item to album format")
				continue
			}
			mediaItemAlbum, err := GetMediaItem[*types.Album](ctx, e, itemAlbum, item.Id)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", item.Id).
					Str("itemName", item.Name).
					Msg("Error creating media item for album")
				continue
			}
			results.Albums = append(results.Albums, mediaItemAlbum)

		case "Audio":
			itemTrack, err := GetItem[*types.Track](ctx, e, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", item.Id).
					Str("itemName", item.Name).
					Msg("Error converting Emby item to track format")
				continue
			}
			mediaItemTrack, err := GetMediaItem[*types.Track](ctx, e, itemTrack, item.Id)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", item.Id).
					Str("itemName", item.Name).
					Msg("Error creating media item for track")
				continue
			}
			results.Tracks = append(results.Tracks, mediaItemTrack)

		case "Playlist":
			itemPlaylist, err := GetItem[*types.Playlist](ctx, e, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", item.Id).
					Str("itemName", item.Name).
					Msg("Error converting Emby item to playlist format")
				continue
			}
			mediaItemPlaylist, err := GetMediaItem[*types.Playlist](ctx, e, itemPlaylist, item.Id)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", item.Id).
					Str("itemName", item.Name).
					Msg("Error creating media item for playlist")
				continue
			}
			results.Playlists = append(results.Playlists, mediaItemPlaylist)

		case "BoxSet":
			itemCollection, err := GetItem[*types.Collection](ctx, e, &item)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", item.Id).
					Str("itemName", item.Name).
					Msg("Error converting Emby item to collection format")
				continue
			}
			mediaItemCollection, err := GetMediaItem[*types.Collection](ctx, e, itemCollection, item.Id)
			if err != nil {
				log.Warn().
					Err(err).
					Str("itemID", item.Id).
					Str("itemName", item.Name).
					Msg("Error creating media item for collection")
				continue
			}
			results.Collections = append(results.Collections, mediaItemCollection)

		default:
			log.Debug().
				Str("itemType", item.Type_).
				Str("itemID", item.Id).
				Str("itemName", item.Name).
				Msg("Skipping unsupported item type in search results")
		}
	}

	log.Info().
		Int("moviesCount", len(results.Movies)).
		Int("seriesCount", len(results.Series)).
		Int("episodesCount", len(results.Episodes)).
		Int("artistsCount", len(results.Artists)).
		Int("albumsCount", len(results.Albums)).
		Int("tracksCount", len(results.Tracks)).
		Int("playlistsCount", len(results.Playlists)).
		Int("collectionsCount", len(results.Collections)).
		Msg("Completed search in Emby server")

	return results, nil
}
