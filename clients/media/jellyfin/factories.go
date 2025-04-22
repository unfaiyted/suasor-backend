package jellyfin

import (
	"context"
	"fmt"
	"time"

	jellyfin "github.com/sj14/jellyfin-go/api"
	media "suasor/clients/media"
	"suasor/clients/media/types"
	"suasor/di/container"
	"suasor/utils/logger"
)

// RegisterMediaItemFactories registers all media item factories for Jellyfin
func RegisterMediaItemFactories(c *container.Container) {

	registry := container.MustGet[media.ClientItemRegistry](c)
	// Register all the media factories for Jellyfin
	media.RegisterFactory[*JellyfinClient, *jellyfin.BaseItemDto, *types.Movie](
		&registry,
		func(client *JellyfinClient, ctx context.Context, item *jellyfin.BaseItemDto) (*types.Movie, error) {
			return client.movieFactory(ctx, item)
		},
	)

	media.RegisterFactory[*JellyfinClient, *jellyfin.BaseItemDto, *types.Track](
		&registry,
		func(client *JellyfinClient, ctx context.Context, item *jellyfin.BaseItemDto) (*types.Track, error) {
			return client.trackFactory(ctx, item)
		},
	)

	media.RegisterFactory[*JellyfinClient, *jellyfin.BaseItemDto, *types.Artist](
		&registry,
		func(client *JellyfinClient, ctx context.Context, item *jellyfin.BaseItemDto) (*types.Artist, error) {
			return client.artistFactory(ctx, item)
		},
	)

	media.RegisterFactory[*JellyfinClient, *jellyfin.BaseItemDto, *types.Album](
		&registry,
		func(client *JellyfinClient, ctx context.Context, item *jellyfin.BaseItemDto) (*types.Album, error) {
			return client.albumFactory(ctx, item)
		},
	)

	media.RegisterFactory[*JellyfinClient, *jellyfin.BaseItemDto, *types.Playlist](
		&registry,
		func(client *JellyfinClient, ctx context.Context, item *jellyfin.BaseItemDto) (*types.Playlist, error) {
			return client.playlistFactory(ctx, item)
		},
	)

	media.RegisterFactory[*JellyfinClient, *jellyfin.BaseItemDto, *types.Series](
		&registry,
		func(client *JellyfinClient, ctx context.Context, item *jellyfin.BaseItemDto) (*types.Series, error) {
			return client.seriesFactory(ctx, item)
		},
	)

	media.RegisterFactory[*JellyfinClient, *jellyfin.BaseItemDto, *types.Season](
		&registry,
		func(client *JellyfinClient, ctx context.Context, item *jellyfin.BaseItemDto) (*types.Season, error) {
			return client.seasonFactory(ctx, item)
		},
	)

	media.RegisterFactory[*JellyfinClient, *jellyfin.BaseItemDto, *types.Episode](
		&registry,
		func(client *JellyfinClient, ctx context.Context, item *jellyfin.BaseItemDto) (*types.Episode, error) {
			return client.episodeFactory(ctx, item)
		},
	)

	media.RegisterFactory[*JellyfinClient, *jellyfin.BaseItemDto, *types.Collection](
		&registry,
		func(client *JellyfinClient, ctx context.Context, item *jellyfin.BaseItemDto) (*types.Collection, error) {
			return client.collectionFactory(ctx, item)
		},
	)
}

// Factory function for Movie
func (j *JellyfinClient) movieFactory(ctx context.Context, item *jellyfin.BaseItemDto) (*types.Movie, error) {
	log := logger.LoggerFromContext(ctx)

	if item.Id == nil || *item.Id == "" {
		return nil, fmt.Errorf("movie is missing required ID field")
	}

	// Safely get name or fallback to empty string
	title := ""
	if item.Name.IsSet() {
		title = *item.Name.Get()
	}

	log.Debug().
		Str("movieID", *item.Id).
		Str("movieName", title).
		Msg("Converting Jellyfin item to movie format")

	// Safely handle optional fields
	description := ""
	if item.Overview.IsSet() {
		description = *item.Overview.Get()
	}

	contentRating := ""
	if item.OfficialRating.IsSet() {
		contentRating = *item.OfficialRating.Get()
	}

	// Determine release year from either ProductionYear or PremiereDate
	var releaseYear int
	var releaseDate time.Time

	if item.ProductionYear.IsSet() {
		releaseYear = int(*item.ProductionYear.Get())
	}

	if item.PremiereDate.IsSet() {
		releaseDate = *item.PremiereDate.Get()
		if releaseYear == 0 {
			releaseYear = releaseDate.Year()
			log.Debug().
				Str("movieID", *item.Id).
				Str("premiereDate", releaseDate.Format("2006-01-02")).
				Int("extractedYear", releaseYear).
				Msg("Using year from premiere date instead of production year")
		}
	}

	// Extract genres
	var genres []string
	if item.Genres != nil {
		genres = item.Genres
	}

	// Calculate duration
	var durationSecs int64 = 0
	if item.RunTimeTicks.IsSet() {
		durationSecs = int64(*item.RunTimeTicks.Get() / 10000000)
	}

	// Initialize ratings
	ratings := types.Ratings{}

	// Safely add community rating if available
	if item.CommunityRating.IsSet() {
		ratings = append(ratings, types.Rating{
			Source: "jellyfin",
			Value:  float32(*item.CommunityRating.Get()),
		})
	}

	movie := &types.Movie{
		Details: types.MediaDetails{
			Title:         title,
			Description:   description,
			ReleaseDate:   releaseDate,
			ReleaseYear:   releaseYear,
			ContentRating: contentRating,
			Genres:        genres,
			Artwork:       *j.getArtworkURLs(item),
			Duration:      durationSecs,
			Ratings:       ratings,
		},
	}

	// Set user rating if available
	if item.UserData.IsSet() && item.UserData.Get().Rating.IsSet() {
		movie.Details.UserRating = float32(*item.UserData.Get().Rating.Get())
		// Set favorite status if available

		isFavorite := *item.UserData.Get().IsFavorite
		movie.Details.IsFavorite = isFavorite

	} else {
		log.Debug().
			Str("movieID", *item.Id).
			Msg("Movie has no user data, skipping user rating")
	}

	// Extract provider IDs if available
	extractProviderIDs(&item.ProviderIds, &movie.Details.ExternalIDs)

	log.Debug().
		Str("movieID", *item.Id).
		Str("movieTitle", movie.Details.Title).
		Int("year", movie.Details.ReleaseYear).
		Msg("Successfully converted Jellyfin item to movie")

	return movie, nil
}

// Factory function for Track
func (j *JellyfinClient) trackFactory(ctx context.Context, item *jellyfin.BaseItemDto) (*types.Track, error) {
	if item.Id == nil || *item.Id == "" {
		return nil, fmt.Errorf("track is missing required ID field")
	}

	// Safely get name or fallback to empty string
	title := ""
	if item.Name.IsSet() {
		title = *item.Name.Get()
	}

	// Safely handle optional fields
	description := ""
	if item.Overview.IsSet() {
		description = *item.Overview.Get()
	}

	// Calculate duration
	var durationSecs int64 = 0
	if item.RunTimeTicks.IsSet() {
		durationSecs = int64(*item.RunTimeTicks.Get() / 10000000)
	}

	// Safely handle track number
	trackNumber := 0
	if item.IndexNumber.IsSet() {
		trackNumber = int(*item.IndexNumber.Get())
	}

	// Safely handle album name
	albumName := ""
	if item.Album.IsSet() {
		albumName = *item.Album.Get()
	}

	track := &types.Track{
		Details: types.MediaDetails{
			Title:       title,
			Description: description,
			Duration:    durationSecs,
			Artwork:     *j.getArtworkURLs(item),
		},
		Number:    trackNumber,
		AlbumName: albumName,
	}

	// Set album ID if available
	if item.AlbumId.IsSet() {
		track.SyncAlbum.AddClient(j.ClientID, *item.AlbumId.Get())
	}

	// Add artist information if available
	if item.AlbumArtists != nil && len(item.AlbumArtists) > 0 {
		artistID := *item.AlbumArtists[0].Id
		name := *item.AlbumArtists[0].Name.Get()
		track.AddSyncClient(j.ClientID, *item.AlbumId.Get(), artistID)
		track.ArtistName = name
	}

	// Extract provider IDs
	extractProviderIDs(&item.ProviderIds, &track.Details.ExternalIDs)

	return track, nil
}

// Factory function for Artist
func (j *JellyfinClient) artistFactory(ctx context.Context, item *jellyfin.BaseItemDto) (*types.Artist, error) {
	if item.Id == nil || *item.Id == "" {
		return nil, fmt.Errorf("artist is missing required ID field")
	}

	// Safely get name or fallback to empty string
	name := ""
	if item.Name.IsSet() {
		name = *item.Name.Get()
	}

	// Safely handle optional fields
	description := ""
	if item.Overview.IsSet() {
		description = *item.Overview.Get()
	}

	// Safely handle album count
	albumCount := 0
	if item.ChildCount.IsSet() {
		albumCount = int(*item.ChildCount.Get())
	}

	// Safely handle genres
	var genres []string
	if item.Genres != nil {
		genres = item.Genres
	}

	artist := &types.Artist{
		Details: types.MediaDetails{
			Title:       name,
			Description: description,
			Genres:      genres,
			Artwork:     *j.getArtworkURLs(item),
		},
		AlbumCount: albumCount,
	}

	// Add user rating if available
	if item.UserData.IsSet() && item.UserData.Get().Rating.IsSet() {
		artist.Details.UserRating = float32(*item.UserData.Get().Rating.Get())
	}

	// Add community rating if available
	if item.CommunityRating.IsSet() {
		artist.Details.Ratings = append(artist.Details.Ratings, types.Rating{
			Source: "jellyfin",
			Value:  float32(*item.CommunityRating.Get()),
		})
	}

	// Extract provider IDs if available
	extractProviderIDs(&item.ProviderIds, &artist.Details.ExternalIDs)

	return artist, nil
}

// Factory function for Album
func (j *JellyfinClient) albumFactory(ctx context.Context, item *jellyfin.BaseItemDto) (*types.Album, error) {
	if item.Id == nil || *item.Id == "" {
		return nil, fmt.Errorf("album is missing required ID field")
	}

	// Safely get name or fallback to empty string
	title := ""
	if item.Name.IsSet() {
		title = *item.Name.Get()
	}

	// Safely handle optional fields
	description := ""
	if item.Overview.IsSet() {
		description = *item.Overview.Get()
	}

	// Safely handle artist name
	artistName := ""
	if item.AlbumArtist.IsSet() {
		artistName = *item.AlbumArtist.Get()
	}

	// Safely handle release year
	releaseYear := 0
	if item.ProductionYear.IsSet() {
		releaseYear = int(*item.ProductionYear.Get())
	}

	// Safely handle track count
	trackCount := 0
	if item.ChildCount.IsSet() {
		trackCount = int(*item.ChildCount.Get())
	}

	// Safely handle genres
	var genres []string
	if item.Genres != nil {
		genres = item.Genres
	}

	album := &types.Album{
		Details: types.MediaDetails{
			Title:       title,
			Description: description,
			ReleaseYear: releaseYear,
			Genres:      genres,
			Artwork:     *j.getArtworkURLs(item),
		},
		ArtistName: artistName,
		TrackCount: trackCount,
	}

	// Add user rating if available
	if item.UserData.IsSet() && item.UserData.Get().Rating.IsSet() {
		album.Details.UserRating = float32(*item.UserData.Get().Rating.Get())
	}

	// Add community rating if available
	if item.CommunityRating.IsSet() {
		album.Details.Ratings = append(album.Details.Ratings, types.Rating{
			Source: "jellyfin",
			Value:  float32(*item.CommunityRating.Get()),
		})
	}

	// Extract provider IDs if available
	extractProviderIDs(&item.ProviderIds, &album.Details.ExternalIDs)

	return album, nil
}

// Factory function for Playlist
func (j *JellyfinClient) playlistFactory(ctx context.Context, item *jellyfin.BaseItemDto) (*types.Playlist, error) {
	if item.Id == nil || *item.Id == "" {
		return nil, fmt.Errorf("playlist is missing required ID field")
	}

	// Safely get name or fallback to empty string
	title := ""
	if item.Name.IsSet() {
		title = *item.Name.Get()
	}

	// Safely handle optional fields
	description := ""
	if item.Overview.IsSet() {
		description = *item.Overview.Get()
	}

	// Safely handle item count
	itemCount := 0
	if item.ChildCount.IsSet() {
		itemCount = int(*item.ChildCount.Get())
	}

	playlist := &types.Playlist{
		ItemList: types.ItemList{
			Details: types.MediaDetails{
				Title:       title,
				Description: description,
				Artwork:     *j.getArtworkURLs(item),
			},
			ItemCount: itemCount,
			IsPublic:  true, // Assume public by default in Jellyfin
		},
	}

	// Add community rating if available
	if item.CommunityRating.IsSet() {
		playlist.Details.Ratings = append(playlist.Details.Ratings, types.Rating{
			Source: "jellyfin",
			Value:  float32(*item.CommunityRating.Get()),
		})
	}

	// Add user rating if available
	if item.UserData.IsSet() && item.UserData.Get().Rating.IsSet() {
		playlist.Details.UserRating = float32(*item.UserData.Get().Rating.Get())
	}

	return playlist, nil
}

// Factory function for Series
func (j *JellyfinClient) seriesFactory(ctx context.Context, item *jellyfin.BaseItemDto) (*types.Series, error) {
	if item.Id == nil || *item.Id == "" {
		return nil, fmt.Errorf("TV show is missing required ID field")
	}

	// Safely get name or fallback to empty string
	title := ""
	if item.Name.IsSet() {
		title = *item.Name.Get()
	}

	// Safely handle optional fields
	description := ""
	if item.Overview.IsSet() {
		description = *item.Overview.Get()
	}

	// Default values
	releaseYear := 0
	if item.ProductionYear.IsSet() {
		releaseYear = int(*item.ProductionYear.Get())
	}

	// Safely handle genres
	var genres []string
	if item.Genres != nil {
		genres = item.Genres
	}

	// Safely handle duration
	var durationSecs int64 = 0
	if item.RunTimeTicks.IsSet() {
		durationSecs = int64(*item.RunTimeTicks.Get() / 10000000)
	}

	// Safely handle season count
	seasonCount := 0
	if item.ChildCount.IsSet() {
		seasonCount = int(*item.ChildCount.Get())
	}

	// Safely handle status
	status := ""
	if item.Status.IsSet() {
		status = *item.Status.Get()
	}

	series := &types.Series{
		Details: types.MediaDetails{
			Title:       title,
			Description: description,
			ReleaseYear: releaseYear,
			Genres:      genres,
			Artwork:     *j.getArtworkURLs(item),
			Duration:    durationSecs,
		},
		Status:      status,
		SeasonCount: seasonCount,
	}

	// Set SeriesStudio if available
	if item.SeriesStudio.IsSet() {
		series.Network = *item.SeriesStudio.Get()
	}

	// Extract provider IDs if available
	extractProviderIDs(&item.ProviderIds, &series.Details.ExternalIDs)

	// Set ratings if available
	if item.CommunityRating.IsSet() {
		series.Details.Ratings = append(series.Details.Ratings, types.Rating{
			Source: "jellyfin",
			Value:  float32(*item.CommunityRating.Get()),
		})
	}

	// Set user rating if available
	if item.UserData.IsSet() && item.UserData.Get().Rating.IsSet() {
		series.Details.UserRating = float32(*item.UserData.Get().Rating.Get())
	}

	return series, nil
}

// Factory function for Season
func (j *JellyfinClient) seasonFactory(ctx context.Context, item *jellyfin.BaseItemDto) (*types.Season, error) {
	if item.Id == nil || *item.Id == "" {
		return nil, fmt.Errorf("season is missing required ID field")
	}

	// Safely get name or fallback to empty string
	title := ""
	if item.Name.IsSet() {
		title = *item.Name.Get()
	}

	// Safely handle optional fields
	description := ""
	if item.Overview.IsSet() {
		description = *item.Overview.Get()
	}

	// Safely handle season number
	seasonNumber := 0
	if item.IndexNumber.IsSet() {
		seasonNumber = int(*item.IndexNumber.Get())
	}

	// Safely handle episode count
	episodeCount := 0
	if item.ChildCount.IsSet() {
		episodeCount = int(*item.ChildCount.Get())
	}

	// Safely handle series name
	seriesName := ""
	if item.SeriesName.IsSet() {
		seriesName = *item.SeriesName.Get()
	}

	season := &types.Season{
		Details: types.MediaDetails{
			Title:       title,
			Description: description,
			Artwork:     *j.getArtworkURLs(item),
		},
		Number:       seasonNumber,
		EpisodeCount: episodeCount,
		SeriesName:   seriesName,
	}

	// Safely handle series ID
	if item.SeriesId.IsSet() {
		season.SyncSeries.AddClient(j.ClientID, *item.SeriesId.Get())
	}

	// Add release year if available
	if item.ProductionYear.IsSet() {
		season.Details.ReleaseYear = int(*item.ProductionYear.Get())
	}

	// Add premiere date if available
	if item.PremiereDate.IsSet() {
		season.Details.ReleaseDate = *item.PremiereDate.Get()
	}

	// Add community rating if available
	if item.CommunityRating.IsSet() {
		season.Details.Ratings = append(season.Details.Ratings, types.Rating{
			Source: "jellyfin",
			Value:  float32(*item.CommunityRating.Get()),
		})
	}

	// Add user rating if available
	if item.UserData.IsSet() && item.UserData.Get().Rating.IsSet() {
		season.Details.UserRating = float32(*item.UserData.Get().Rating.Get())
	}

	// Extract provider IDs if available
	extractProviderIDs(&item.ProviderIds, &season.Details.ExternalIDs)

	return season, nil
}

// Factory function for Episode
func (j *JellyfinClient) episodeFactory(ctx context.Context, item *jellyfin.BaseItemDto) (*types.Episode, error) {
	if item.Id == nil || *item.Id == "" {
		return nil, fmt.Errorf("episode is missing required ID field")
	}

	// Safely get name or fallback to empty string
	title := ""
	if item.Name.IsSet() {
		title = *item.Name.Get()
	}

	// Safely handle optional fields
	description := ""
	if item.Overview.IsSet() {
		description = *item.Overview.Get()
	}

	// Safely handle duration
	var durationSecs int64 = 0
	if item.RunTimeTicks.IsSet() {
		durationSecs = int64(*item.RunTimeTicks.Get() / 10000000)
	}

	// Safely handle episode number
	episodeNumber := int64(0)
	if item.IndexNumber.IsSet() {
		episodeNumber = int64(*item.IndexNumber.Get())
	}

	// Safely handle season number
	seasonNumber := 0
	if item.ParentIndexNumber.IsSet() {
		seasonNumber = int(*item.ParentIndexNumber.Get())
	}

	// Safely handle show title
	showTitle := ""
	if item.SeriesName.IsSet() {
		showTitle = *item.SeriesName.Get()
	}

	episode := &types.Episode{
		Details: types.MediaDetails{
			Title:       title,
			Description: description,
			Artwork:     *j.getArtworkURLs(item),
			Duration:    durationSecs,
		},
		Number:       episodeNumber,
		SeasonNumber: seasonNumber,
		ShowTitle:    showTitle,
	}

	// Safely set IDs if available
	if item.SeriesId.IsSet() && item.SeasonId.IsSet() {
		episode.AddSyncClient(j.ClientID, *item.SeasonId.Get(), *item.SeriesId.Get())
	}

	// Add air date if available
	if item.PremiereDate.IsSet() {
		episode.Details.ReleaseDate = *item.PremiereDate.Get()
	}

	// Add community rating if available
	if item.CommunityRating.IsSet() {
		episode.Details.Ratings = append(episode.Details.Ratings, types.Rating{
			Source: "jellyfin",
			Value:  float32(*item.CommunityRating.Get()),
		})
	}

	// Add user rating if available
	if item.UserData.IsSet() && item.UserData.Get().Rating.IsSet() {
		episode.Details.UserRating = float32(*item.UserData.Get().Rating.Get())
	}

	// Extract provider IDs if available
	extractProviderIDs(&item.ProviderIds, &episode.Details.ExternalIDs)

	return episode, nil
}

// Factory function for Collection
func (j *JellyfinClient) collectionFactory(ctx context.Context, item *jellyfin.BaseItemDto) (*types.Collection, error) {
	if item.Id == nil || *item.Id == "" {
		return nil, fmt.Errorf("collection is missing required ID field")
	}

	// Safely get name or fallback to empty string
	title := ""
	if item.Name.IsSet() {
		title = *item.Name.Get()
	}

	// Safely handle optional fields
	description := ""
	if item.Overview.IsSet() {
		description = *item.Overview.Get()
	}

	// Safely handle item count
	itemCount := 0
	if item.ChildCount.IsSet() {
		itemCount = int(*item.ChildCount.Get())
	}

	collection := &types.Collection{
		ItemList: types.ItemList{
			Details: types.MediaDetails{
				Title:       title,
				Description: description,
				Artwork:     *j.getArtworkURLs(item),
			},
			ItemCount: itemCount,
			IsPublic:  true, // Assume public by default in Jellyfin
		},
	}

	// Add potential year if available
	if item.ProductionYear.IsSet() {
		collection.Details.ReleaseYear = int(*item.ProductionYear.Get())
	}

	// Add community rating if available
	if item.CommunityRating.IsSet() {
		collection.Details.Ratings = append(collection.Details.Ratings, types.Rating{
			Source: "jellyfin",
			Value:  float32(*item.CommunityRating.Get()),
		})
	}

	// Add user rating if available
	if item.UserData.IsSet() && item.UserData.Get().Rating.IsSet() {
		collection.Details.UserRating = float32(*item.UserData.Get().Rating.Get())
	}

	// Handle genres if available
	if item.Genres != nil {
		collection.Details.Genres = item.Genres
	}

	// Extract provider IDs if available
	extractProviderIDs(&item.ProviderIds, &collection.Details.ExternalIDs)

	return collection, nil
}
