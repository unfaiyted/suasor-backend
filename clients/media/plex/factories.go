package plex

import (
	"context"
	"fmt"
	"strconv"

	"github.com/LukeHagar/plexgo/models/operations"
	media "suasor/clients/media"
	"suasor/clients/media/types"
	"suasor/di/container"
	"suasor/utils/logger"
)

// RegisterMediaItemFactories registers all media item factories for Plex
func RegisterMediaItemFactories(c *container.Container) {
	registry := container.MustGet[media.ClientItemRegistry](c)

	// Register factory for GetLibraryItems response type
	media.RegisterFactory[*PlexClient, *operations.GetLibraryItemsMetadata, *types.Movie](
		&registry,
		func(client *PlexClient, ctx context.Context, item *operations.GetLibraryItemsMetadata) (*types.Movie, error) {
			return client.movieFactory(ctx, item)
		},
	)

	media.RegisterFactory[*PlexClient, *operations.GetLibraryItemsMetadata, *types.Series](
		&registry,
		func(client *PlexClient, ctx context.Context, item *operations.GetLibraryItemsMetadata) (*types.Series, error) {
			return client.seriesFactory(ctx, item)
		},
	)

	media.RegisterFactory[*PlexClient, *operations.GetMetadataChildrenMetadata, *types.Season](
		&registry,
		func(client *PlexClient, ctx context.Context, item *operations.GetMetadataChildrenMetadata) (*types.Season, error) {
			return client.seasonFactory(ctx, item)
		},
	)

	media.RegisterFactory[*PlexClient, *operations.GetMediaMetaDataMetadata, *types.Episode](
		&registry,
		func(client *PlexClient, ctx context.Context, item *operations.GetMediaMetaDataMetadata) (*types.Episode, error) {
			return client.episodeFactory(ctx, item)
		},
	)

	media.RegisterFactory[*PlexClient, *operations.GetLibraryItemsMetadata, *types.Collection](
		&registry,
		func(client *PlexClient, ctx context.Context, item *operations.GetLibraryItemsMetadata) (*types.Collection, error) {
			return client.collectionFactory(ctx, item)
		},
	)

	media.RegisterFactory[*PlexClient, *operations.GetPlaylistMetadata, *types.Playlist](
		&registry,
		func(client *PlexClient, ctx context.Context, item *operations.GetPlaylistMetadata) (*types.Playlist, error) {
			return client.playlistFactory(ctx, item)
		},
	)

	media.RegisterFactory[*PlexClient, *operations.GetMetadataChildrenMetadata, *types.Album](
		&registry,
		func(client *PlexClient, ctx context.Context, item *operations.GetMetadataChildrenMetadata) (*types.Album, error) {
			return client.albumFactory(ctx, item)
		},
	)

	media.RegisterFactory[*PlexClient, *operations.GetMediaMetaDataMetadata, *types.Artist](
		&registry,
		func(client *PlexClient, ctx context.Context, item *operations.GetMediaMetaDataMetadata) (*types.Artist, error) {
			return client.artistFactory(ctx, item)
		},
	)

	media.RegisterFactory[*PlexClient, *operations.GetLibraryItemsMetadata, *types.Track](
		&registry,
		func(client *PlexClient, ctx context.Context, item *operations.GetLibraryItemsMetadata) (*types.Track, error) {
			return client.trackFactory(ctx, item)
		},
	)

	// Register factories for GetMediaMetaData response type if needed
	// Add more factory registrations for other Plex-specific types as needed
}

// Factory function for Movie
func (c *PlexClient) movieFactory(ctx context.Context, item *operations.GetLibraryItemsMetadata) (*types.Movie, error) {
	log := logger.LoggerFromContext(ctx)

	if item.RatingKey == "" {
		return nil, fmt.Errorf("movie is missing required ID field (RatingKey)")
	}

	log.Debug().
		Str("movieID", item.RatingKey).
		Str("movieTitle", item.Title).
		Msg("Converting Plex item to movie format")

	// Create base metadata using helper
	metadata := c.createDetailsFromLibraryMetadata(item)

	// Create the movie object
	movie := &types.Movie{
		Details: &metadata,
	}

	log.Debug().
		Str("movieID", item.RatingKey).
		Str("movieTitle", movie.Details.Title).
		Int("year", movie.Details.ReleaseYear).
		Msg("Successfully converted Plex item to movie")

	return movie, nil
}

// Factory function for Series
func (c *PlexClient) seriesFactory(ctx context.Context, item *operations.GetLibraryItemsMetadata) (*types.Series, error) {
	log := logger.LoggerFromContext(ctx)

	if item.RatingKey == "" {
		return nil, fmt.Errorf("series is missing required ID field (RatingKey)")
	}

	log.Debug().
		Str("seriesID", item.RatingKey).
		Str("seriesTitle", item.Title).
		Msg("Converting Plex item to series format")

	// Create base metadata using helper
	metadata := c.createDetailsFromLibraryMetadata(item)

	// Create the series object
	series := &types.Series{
		Details: &metadata,
	}

	if item.Rating != nil {
		series.Rating = float64(*item.Rating)
	}

	if item.Studio != nil {
		series.Network = *item.Studio
	}

	if item.Rating != nil {
		series.Rating = float64(*item.Rating)
	}
	if item.Year != nil {
		series.ReleaseYear = *item.Year
	}
	if item.ContentRating != nil {
		series.ContentRating = *item.ContentRating
	}
	if item.ChildCount != nil {
		series.SeasonCount = *item.ChildCount
	}

	if item.LeafCount != nil {
		series.EpisodeCount = int(*item.LeafCount)
	}

	if item.Genre != nil {
		series.Genres = make([]string, 0, len(item.Genre))
		for _, genre := range item.Genre {
			if genre.Tag != nil {
				series.Genres = append(series.Genres, *genre.Tag)
			}
		}
	}

	// Add status if available
	// if item.Status != nil {
	// 	series.Status = *item.Status
	// }

	log.Debug().
		Str("seriesID", item.RatingKey).
		Str("seriesTitle", series.Details.Title).
		Int("seasonCount", series.SeasonCount).
		Msg("Successfully converted Plex item to series")

	return series, nil
}

// Factory function for Season
func (c *PlexClient) seasonFactory(ctx context.Context, item *operations.GetMetadataChildrenMetadata) (*types.Season, error) {
	log := logger.LoggerFromContext(ctx)

	if item.RatingKey == nil || *item.RatingKey == "" {
		return nil, fmt.Errorf("season is missing required ID field (RatingKey)")
	}

	log.Debug().
		Str("seasonID", *item.RatingKey).
		Str("seasonTitle", *item.Title).
		Msg("Converting Plex item to season format")

	// Build season object
	season := &types.Season{
		Details: c.createDetailsFromMetadataChildren(item),
		Number:  int(*item.Index),
	}

	// Add specific season fields
	if item.Index != nil {
		season.Number = int(*item.Index)
	}

	if item.LeafCount != nil {
		season.EpisodeCount = int(*item.LeafCount)
	}

	if item.ParentTitle != nil {
		season.SeriesName = *item.ParentTitle
	}

	// Add parent series ID for synchronization
	if item.ParentRatingKey != nil {
		// TODO:
	}

	log.Debug().
		Str("seasonID", *item.RatingKey).
		Str("seasonTitle", season.Details.Title).
		Int("seasonNumber", season.Number).
		Int("episodeCount", season.EpisodeCount).
		Msg("Successfully converted Plex item to season")

	return season, nil
}

// Factory function for Episode
func (c *PlexClient) episodeFactory(ctx context.Context, item *operations.GetMediaMetaDataMetadata) (*types.Episode, error) {
	log := logger.LoggerFromContext(ctx)

	if item.RatingKey == "" {
		return nil, fmt.Errorf("episode is missing required ID field (RatingKey)")
	}

	log.Debug().
		Str("episodeID", item.RatingKey).
		Str("episodeTitle", item.Title).
		Msg("Converting Plex item to episode format")

	// Build episode object
	episode := &types.Episode{
		Details: c.createDetailsFromMediaMetadata(item),
	}

	// Add specific episode fields
	if item.Index != nil {
		episode.Number = int64(*item.Index)
	}

	if item.ParentIndex != nil {
		episode.SeasonNumber = int(*item.ParentIndex)
	}

	if item.ParentIndex != nil {
		episode.SeasonNumber = int(*item.ParentIndex)
	}

	// Add show ID if available (via grandparentRatingKey)
	if item.GrandparentRatingKey != nil {
		// episode.AddSyncClient(c.GetClientID(), *item.GrandparentRatingKey, *item.ParentRatingKey)
	}

	// Add studio if available
	if item.Studio != nil {
		episode.Details.Studio = *item.Studio
	}

	parentRatingKey, err := strconv.Atoi(*item.ParentRatingKey)
	int64RatingKey := int64(parentRatingKey)

	opts := operations.GetMediaMetaDataRequest{
		RatingKey: int64RatingKey,
	}
	// Get Season Info to get Grandparent title
	seasonRes, err := c.plexAPI.Library.GetMediaMetaData(ctx, opts)
	seasonMetadata := seasonRes.Object.MediaContainer.Metadata[0]
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Msg("Failed to get season info from Plex")
		return nil, fmt.Errorf("failed to get season info: %w", err)
	}

	if seasonMetadata.GrandparentTitle != nil {
		episode.ShowTitle = *item.GrandparentTitle
	}

	// Add parent IDs for synchronization
	if item.ParentRatingKey != nil && item.GrandparentRatingKey != nil {
		// episode.AddSyncClient(c.GetClientID(), *item.ParentRatingKey, *item.GrandparentRatingKey)
	}

	log.Debug().
		Str("episodeID", item.RatingKey).
		Str("episodeTitle", episode.Details.Title).
		Int64("episodeNumber", episode.Number).
		Int("seasonNumber", episode.SeasonNumber).
		Str("showTitle", episode.ShowTitle).
		Msg("Successfully converted Plex item to episode")

	return episode, nil
}

// Factory function for Collection
func (c *PlexClient) collectionFactory(ctx context.Context, item *operations.GetLibraryItemsMetadata) (*types.Collection, error) {
	log := logger.LoggerFromContext(ctx)

	if item.RatingKey == "" {
		return nil, fmt.Errorf("collection is missing required ID field (RatingKey)")
	}

	log.Debug().
		Str("collectionID", item.RatingKey).
		Str("collectionTitle", item.Title).
		Msg("Converting Plex item to collection format")

	// Create the collection object
	collection := &types.Collection{
		ItemList: types.ItemList{
			// TODO: Get the items for the list itself
			Items: []types.ListItem{},
			Details: &types.MediaDetails{
				Title:       item.Title,
				Description: item.Summary,
				Artwork: types.Artwork{
					Thumbnail: c.makeFullURL(*item.Thumb),
				},
			},
			IsPublic: true, // Plex collections are typically public
		},
	}

	// TODO: Needs to fix the sync client states
	// collection.ItemList.SyncClientStates
	// Set item count if available
	if item.ChildCount != nil {
		collection.ItemCount = int(*item.ChildCount)
	}

	log.Debug().
		Str("collectionID", item.RatingKey).
		Str("collectionTitle", collection.Details.Title).
		Int("itemCount", collection.ItemCount).
		Msg("Successfully converted Plex item to collection")

	return collection, nil
}

// Factory function for Playlist
func (c *PlexClient) playlistFactory(ctx context.Context, item *operations.GetPlaylistMetadata) (*types.Playlist, error) {
	log := logger.LoggerFromContext(ctx)

	if *item.RatingKey == "" {
		return nil, fmt.Errorf("playlist is missing required ID field (RatingKey)")
	}

	log.Debug().
		Str("playlistID", *item.RatingKey).
		Str("playlistTitle", *item.Title).
		Msg("Converting Plex item to playlist format")

	// Create the playlist object
	playlist := &types.Playlist{
		ItemList: types.ItemList{
			Details: &types.MediaDetails{
				Title:       *item.Title,
				Description: *item.Summary,
				Artwork:     types.Artwork{
					// Thumbnail: c.makeFullURL(*item.Thumb),
				},
				ExternalIDs: types.ExternalIDs{types.ExternalID{
					Source: "plex",
					ID:     *item.RatingKey,
				}},
			},
			IsPublic: true, // Assume public by default in Plex
		},
	}

	// Set item count if available
	if item.LeafCount != nil {
		playlist.ItemCount = int(*item.LeafCount)
	}

	log.Debug().
		Str("playlistID", *item.RatingKey).
		Str("playlistTitle", playlist.Details.Title).
		Int("itemCount", playlist.ItemCount).
		Msg("Successfully converted Plex item to playlist")

	return playlist, nil
}

// Factory function for Album
func (c *PlexClient) albumFactory(ctx context.Context, item *operations.GetMetadataChildrenMetadata) (*types.Album, error) {
	log := logger.LoggerFromContext(ctx)

	if *item.RatingKey == "" {
		return nil, fmt.Errorf("album is missing required ID field (RatingKey)")
	}

	log.Debug().
		Str("albumID", *item.RatingKey).
		Str("albumTitle", *item.Title).
		Msg("Converting Plex item to album format")

	// Create base metadata
	metadata := c.createDetailsFromMetadataChildren(item)

	// Create the album object
	album := &types.Album{
		Details: metadata,
	}

	// Add specific album fields
	if item.ParentTitle != nil {
		album.ArtistName = *item.ParentTitle
	}

	if item.LeafCount != nil {
		album.TrackCount = int(*item.LeafCount)
	}

	log.Debug().
		Str("albumID", *item.RatingKey).
		Str("albumTitle", album.Details.Title).
		Str("artistName", album.ArtistName).
		Int("trackCount", album.TrackCount).
		Msg("Successfully converted Plex item to album")

	return album, nil
}

// Factory function for Artist
func (c *PlexClient) artistFactory(ctx context.Context, item *operations.GetMediaMetaDataMetadata) (*types.Artist, error) {
	log := logger.LoggerFromContext(ctx)

	if item.RatingKey == "" {
		return nil, fmt.Errorf("artist is missing required ID field (RatingKey)")
	}

	log.Debug().
		Str("artistID", item.RatingKey).
		Str("artistName", item.Title).
		Msg("Converting Plex item to artist format")

	// Create base metadata
	metadata := c.createDetailsFromMediaMetadata(item)

	// Create the artist object
	artist := &types.Artist{
		Details: metadata,
	}

	// Add specific artist fields
	if item.ChildCount != nil {
		artist.AlbumCount = int(*item.ChildCount)
	}

	log.Debug().
		Str("artistID", item.RatingKey).
		Str("artistName", artist.Details.Title).
		Int("albumCount", artist.AlbumCount).
		Msg("Successfully converted Plex item to artist")

	return artist, nil
}

// Factory function for Track
func (c *PlexClient) trackFactory(ctx context.Context, item *operations.GetLibraryItemsMetadata) (*types.Track, error) {
	log := logger.LoggerFromContext(ctx)

	if item.RatingKey == "" {
		return nil, fmt.Errorf("track is missing required ID field (RatingKey)")
	}

	log.Debug().
		Str("trackID", item.RatingKey).
		Str("trackTitle", item.Title).
		Msg("Converting Plex item to track format")

	// Create base metadata
	metadata := c.createDetailsFromLibraryMetadata(item)

	// Create the track object
	track := &types.Track{
		Details: &metadata,
	}

	// Add specific track fields
	if item.Index != nil {
		track.Number = int(*item.Index)
	}

	if item.ParentTitle != nil {
		track.AlbumName = *item.ParentTitle
	}

	if item.GrandparentTitle != nil {
		track.ArtistName = *item.GrandparentTitle
	}

	// Add synchronization IDs
	if item.ParentRatingKey != nil {
		// track.SyncAlbum.AddClient(c.GetClientID(), *item.ParentRatingKey)
	}

	if item.ParentRatingKey != nil && item.GrandparentRatingKey != nil {
		// track.AddSyncClient(c.GetClientID(), *item.ParentRatingKey, *item.GrandparentRatingKey)
	}

	log.Debug().
		Str("trackID", item.RatingKey).
		Str("trackTitle", track.Details.Title).
		Int("trackNumber", track.Number).
		Str("albumName", track.AlbumName).
		Str("artistName", track.ArtistName).
		Msg("Successfully converted Plex item to track")

	return track, nil
}
