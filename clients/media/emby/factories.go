package emby

import (
	"context"
	"fmt"
	"time"

	"suasor/clients/media"
	"suasor/clients/media/types"
	"suasor/di/container"
	embyclient "suasor/internal/clients/embyAPI"
	"suasor/utils/logger"
)

// RegisterMediaItemFactories registers all media item factories for Emby
func RegisterMediaItemFactories(c *container.Container) {

	registry := container.MustGet[media.ClientItemRegistry](c)
	// Register all the media factories for Emby
	media.RegisterFactory[*EmbyClient, *embyclient.BaseItemDto, *types.Movie](
		&registry,
		func(client *EmbyClient, ctx context.Context, item *embyclient.BaseItemDto) (*types.Movie, error) {
			return client.movieFactory(ctx, item)
		},
	)

	media.RegisterFactory[*EmbyClient, *embyclient.BaseItemDto, *types.Track](
		&registry,
		func(client *EmbyClient, ctx context.Context, item *embyclient.BaseItemDto) (*types.Track, error) {
			return client.trackFactory(ctx, item)
		},
	)

	media.RegisterFactory[*EmbyClient, *embyclient.BaseItemDto, *types.Artist](
		&registry,
		func(client *EmbyClient, ctx context.Context, item *embyclient.BaseItemDto) (*types.Artist, error) {
			return client.artistFactory(ctx, item)
		},
	)

	media.RegisterFactory[*EmbyClient, *embyclient.BaseItemDto, *types.Album](
		&registry,
		func(client *EmbyClient, ctx context.Context, item *embyclient.BaseItemDto) (*types.Album, error) {
			return client.albumFactory(ctx, item)
		},
	)

	media.RegisterFactory[*EmbyClient, *embyclient.PlaylistsPlaylistCreationResult, *types.Playlist](
		&registry,
		func(client *EmbyClient, ctx context.Context, item *embyclient.PlaylistsPlaylistCreationResult) (*types.Playlist, error) {
			return client.playlistFactory(ctx, item)
		},
	)

	media.RegisterFactory[*EmbyClient, *embyclient.BaseItemDto, *types.Playlist](
		&registry,
		func(client *EmbyClient, ctx context.Context, item *embyclient.BaseItemDto) (*types.Playlist, error) {
			return client.playlistBaseItemFactory(ctx, item)
		},
	)

	media.RegisterFactory[*EmbyClient, *embyclient.BaseItemDto, *types.Series](
		&registry,
		func(client *EmbyClient, ctx context.Context, item *embyclient.BaseItemDto) (*types.Series, error) {
			return client.seriesFactory(ctx, item)
		},
	)

	media.RegisterFactory[*EmbyClient, *embyclient.BaseItemDto, *types.Season](
		&registry,
		func(client *EmbyClient, ctx context.Context, item *embyclient.BaseItemDto) (*types.Season, error) {
			return client.seasonFactory(ctx, item)
		},
	)

	media.RegisterFactory[*EmbyClient, *embyclient.BaseItemDto, *types.Episode](
		&registry,
		func(client *EmbyClient, ctx context.Context, item *embyclient.BaseItemDto) (*types.Episode, error) {
			return client.episodeFactory(ctx, item)
		},
	)

	media.RegisterFactory[*EmbyClient, *embyclient.BaseItemDto, *types.Collection](
		&registry,
		func(client *EmbyClient, ctx context.Context, item *embyclient.BaseItemDto) (*types.Collection, error) {
			return client.collectionFactory(ctx, item)
		},
	)

}

// Factory function for Movie
func (e *EmbyClient) movieFactory(ctx context.Context, item *embyclient.BaseItemDto) (*types.Movie, error) {
	log := logger.LoggerFromContext(ctx)

	if item.Id == "" {
		return nil, fmt.Errorf("movie is missing required ID field")
	}

	log.Debug().
		Str("movieID", item.Id).
		Str("movieName", item.Name).
		Int32("releaseYear", item.ProductionYear).
		Str("releaseDate", item.PremiereDate.Format("2006-01-02")).
		Msg("Converting Emby item to movie format")

	// Determine release year from either ProductionYear or PremiereDate
	releaseYear := int(item.ProductionYear)
	if releaseYear == 0 && !item.PremiereDate.IsZero() {
		releaseYear = item.PremiereDate.Year()
		log.Debug().
			Str("movieID", item.Id).
			Str("premiereDate", item.PremiereDate.Format("2006-01-02")).
			Int("extractedYear", releaseYear).
			Msg("Using year from premiere date instead of production year")
	}

	movie := &types.Movie{

		Details: &types.MediaDetails{
			Title:         item.Name,
			Description:   item.Overview,
			ReleaseDate:   item.PremiereDate,
			ReleaseYear:   releaseYear,
			AddedAt:       time.Now(),
			UpdatedAt:     time.Now(),
			ContentRating: item.OfficialRating,
			Genres:        item.Genres,
			Studio:        item.SeriesStudio,
			Artwork:       e.getArtworkURLs(item),
			Duration:      int64(item.RunTimeTicks / 10000000),
			Ratings:       getRatingsFromItem(item),
		},
	}

	// Extract provider IDs if available
	if item.ProviderIds != nil {
		movie.Details.ExternalIDs = convertToExternalIDs(item.ProviderIds)
	}

	log.Debug().
		Str("movieID", item.Id).
		Str("movieTitle", movie.Details.Title).
		Int("year", movie.Details.ReleaseYear).
		Msg("Successfully converted Emby item to movie")

	return movie, nil
}

// Factory function for Track
func (e *EmbyClient) trackFactory(ctx context.Context, item *embyclient.BaseItemDto) (*types.Track, error) {
	track := &types.Track{
		Number:     int(item.IndexNumber),
		AlbumName:  item.Album,
		DiscNumber: int(item.ParentIndexNumber),
		Details: &types.MediaDetails{
			Title:         item.Name,
			ReleaseDate:   item.PremiereDate,
			ReleaseYear:   int(item.PremiereDate.Year()),
			Description:   item.Overview,
			AddedAt:       time.Now(),
			UpdatedAt:     time.Now(),
			ContentRating: item.OfficialRating,
			Duration:      int64(item.RunTimeTicks / 10000000),
			Artwork:       e.getArtworkURLs(item),
			Ratings:       getRatingsFromItem(item),
		},
	}

	// track.SyncAlbum.AddClient(e.GetClientID(), item.Id)

	// Add artist information if available
	if len(item.ArtistItems) > 0 {
		// track.AddSyncClient(e.GetClientID(), item.AlbumId, item.ArtistItems[0].Id)
		track.ArtistName = item.ArtistItems[0].Name
	}

	// Extract provider IDs
	if item.ProviderIds != nil {
		track.Details.ExternalIDs = convertToExternalIDs(item.ProviderIds)
	}

	return track, nil
}

// Factory function for Artist
func (e *EmbyClient) artistFactory(ctx context.Context, item *embyclient.BaseItemDto) (*types.Artist, error) {
	artist := &types.Artist{
		Details: &types.MediaDetails{
			Title:         item.Name,
			Description:   item.Overview,
			AddedAt:       time.Now(),
			UpdatedAt:     time.Now(),
			ContentRating: item.OfficialRating,

			Artwork: e.getArtworkURLs(item),
			Genres:  item.Genres,
			Ratings: getRatingsFromItem(item),
		},
	}

	// Extract provider IDs if available
	if item.ProviderIds != nil {
		artist.Details.ExternalIDs = convertToExternalIDs(item.ProviderIds)
	}

	return artist, nil
}

// Factory function for Album
func (e *EmbyClient) albumFactory(ctx context.Context, item *embyclient.BaseItemDto) (*types.Album, error) {
	album := &types.Album{
		Details: &types.MediaDetails{
			Title:       item.Name,
			Description: item.Overview,
			ReleaseYear: int(item.PremiereDate.Year()),
			ReleaseDate: item.PremiereDate,
			AddedAt:     time.Now(),
			UpdatedAt:   time.Now(),
			Genres:      item.Genres,
			Studio:      item.SeriesStudio,
			Artwork:     e.getArtworkURLs(item),
			Ratings:     getRatingsFromItem(item),
		},
		ArtistName: item.AlbumArtist,
		TrackCount: int(item.ChildCount),
	}

	// Extract provider IDs
	if item.ProviderIds != nil {
		album.Details.ExternalIDs = convertToExternalIDs(item.ProviderIds)
	}

	return album, nil
}

// Factory function for Playlist
func (e *EmbyClient) playlistFactory(ctx context.Context, item *embyclient.PlaylistsPlaylistCreationResult) (*types.Playlist, error) {
	playlist := &types.Playlist{
		ItemList: types.ItemList{
			Details: &types.MediaDetails{
				Title:       item.Name,
				Description: "Emby Playlist",
				AddedAt:     time.Now(),
				UpdatedAt:   time.Now(),
				// Artwork:     e.getArtworkURLs(item),
			},
			ItemCount: int(item.ItemAddedCount),
			IsPublic:  true, // Assume public by default in Emby
		},
	}

	return playlist, nil
}

func (e *EmbyClient) playlistBaseItemFactory(ctx context.Context, item *embyclient.BaseItemDto) (*types.Playlist, error) {
	playlist := &types.Playlist{
		ItemList: types.ItemList{
			Details: &types.MediaDetails{
				Title:       item.Name,
				Description: item.Overview,
				Artwork:     e.getArtworkURLs(item),
			},
			ItemCount: int(item.ChildCount),
			IsPublic:  true, // Assume public by default in Emby
		},
	}

	return playlist, nil
}

// Factory function for Series
func (e *EmbyClient) seriesFactory(ctx context.Context, item *embyclient.BaseItemDto) (*types.Series, error) {
	series := &types.Series{

		ContentRating: item.OfficialRating,
		Details: &types.MediaDetails{
			Title:         item.Name,
			Description:   item.Overview,
			ReleaseDate:   item.PremiereDate,
			ReleaseYear:   int(item.PremiereDate.Year()),
			AddedAt:       time.Now(),
			UpdatedAt:     time.Now(),
			Genres:        item.Genres,
			ContentRating: item.OfficialRating,
			Studio:        item.SeriesStudio,
			Artwork:       e.getArtworkURLs(item),
			Duration:      int64(item.RunTimeTicks / 10000000),
			Ratings:       getRatingsFromItem(item),
		},
		ReleaseYear: int(item.PremiereDate.Year()),
		SeasonCount: int(item.ChildCount),
		Status:      item.Status,
		Network:     item.SeriesStudio,
	}

	// Extract provider IDs if available
	if item.ProviderIds != nil {
		series.Details.ExternalIDs = convertToExternalIDs(item.ProviderIds)
	}

	return series, nil
}

// Factory function for Season
func (e *EmbyClient) seasonFactory(ctx context.Context, item *embyclient.BaseItemDto) (*types.Season, error) {
	season := &types.Season{
		Details: &types.MediaDetails{
			Title:         item.Name,
			Description:   item.Overview,
			ReleaseDate:   item.PremiereDate,
			ReleaseYear:   int(item.PremiereDate.Year()),
			AddedAt:       time.Now(),
			ContentRating: item.OfficialRating,
			UpdatedAt:     time.Now(),
			Artwork:       e.getArtworkURLs(item),
			Ratings:       getRatingsFromItem(item),
		},
		Number:       int(item.IndexNumber),
		EpisodeCount: int(item.ChildCount),
	}

	// season.SyncSeries.AddClient(e.GetClientID(), item.ParentId)

	if !item.PremiereDate.IsZero() {
		season.ReleaseDate = item.PremiereDate
	}

	return season, nil
}

// Factory function for Episode
func (e *EmbyClient) episodeFactory(ctx context.Context, item *embyclient.BaseItemDto) (*types.Episode, error) {
	episode := &types.Episode{
		Details: &types.MediaDetails{
			ReleaseYear:   int(item.PremiereDate.Year()),
			ReleaseDate:   item.PremiereDate,
			Title:         item.Name,
			Description:   item.Overview,
			AddedAt:       time.Now(),
			UpdatedAt:     time.Now(),
			ContentRating: item.OfficialRating,
			Genres:        item.Genres,
			Artwork:       e.getArtworkURLs(item),
			Duration:      int64(item.RunTimeTicks / 10000000),
			Ratings:       getRatingsFromItem(item),
		},
		Number:       int64(item.IndexNumber),
		SeasonNumber: int(item.ParentIndexNumber),
		ShowTitle:    item.SeriesName,
	}

	// episode.AddSyncClient(e.GetClientID(), item.SeasonId, item.SeriesId)

	// Add external IDs
	if item.ProviderIds != nil {
		episode.Details.ExternalIDs = convertToExternalIDs(item.ProviderIds)
	}

	return episode, nil
}

// Factory function for Collection
func (e *EmbyClient) collectionFactory(ctx context.Context, item *embyclient.BaseItemDto) (*types.Collection, error) {
	collection := &types.Collection{
		ItemList: types.ItemList{
			Details: &types.MediaDetails{

				Title:       item.Name,
				Description: item.Overview,
				AddedAt:     time.Now(),
				UpdatedAt:   time.Now(),

				Artwork: e.getArtworkURLs(item),
			},
			ItemCount: int(item.ChildCount),
			IsPublic:  true, // Assume public by default in Emby
		},
	}

	return collection, nil
}

func getRatingsFromItem(item *embyclient.BaseItemDto) types.Ratings {
	return types.Ratings{
		types.Rating{
			Source: "emby-community",
			Value:  float32(item.CommunityRating),
		},
		types.Rating{
			Source: "emby-critic",
			Value:  float32(item.CriticRating),
		},
	}
}
