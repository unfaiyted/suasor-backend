package jellyfin

import (
	"context"
	"fmt"
	"time"

	jellyfin "github.com/sj14/jellyfin-go/api"
	t "suasor/client/media/types"
	"suasor/types/models"
	"suasor/utils"
)

// Helper function to convert Jellyfin item to internal Collection type
func (j *JellyfinClient) convertToCollection(ctx context.Context, item *jellyfin.BaseItemDto) (models.MediaItem[t.Collection], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	// Validate required fields
	if item == nil {
		return models.MediaItem[t.Collection]{}, fmt.Errorf("cannot convert nil item to collection")
	}

	if item.Id == nil || *item.Id == "" {
		return models.MediaItem[t.Collection]{}, fmt.Errorf("collection is missing required ID field")
	}

	// Safely get name or fallback to empty string
	title := ""
	if item.Name.IsSet() {
		title = *item.Name.Get()
	}

	log.Debug().
		Str("collectionID", *item.Id).
		Str("collectionName", title).
		Msg("Converting Jellyfin item to collection format")

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

	// Build collection object
	collection := models.MediaItem[t.Collection]{
		Data: t.Collection{
			Details: t.MediaDetails{
				Title:       title,
				Description: description,
				Artwork:     j.getArtworkURLs(item),
			},
			ItemCount: itemCount,
		},
		Type: "collection",
	}
	collection.SetClientInfo(j.ClientID, j.ClientType, *item.Id)

	// Add potential year if available
	if item.ProductionYear.IsSet() {
		collection.Data.Details.ReleaseYear = int(*item.ProductionYear.Get())
	}

	// Add community rating if available
	if item.CommunityRating.IsSet() {
		collection.Data.Details.Ratings = append(collection.Data.Details.Ratings, t.Rating{
			Source: "jellyfin",
			Value:  float32(*item.CommunityRating.Get()),
		})
	}

	// Add user rating if available
	if item.UserData.IsSet() && item.UserData.Get().Rating.IsSet() {
		collection.Data.Details.UserRating = float32(*item.UserData.Get().Rating.Get())
	}

	// Handle genres if available
	if item.Genres != nil {
		collection.Data.Details.Genres = item.Genres
	}

	// Extract provider IDs if available
	extractProviderIDs(&item.ProviderIds, &collection.Data.Details.ExternalIDs)

	log.Debug().
		Str("collectionID", *item.Id).
		Str("collectionName", collection.Data.Details.Title).
		Int("itemCount", collection.Data.ItemCount).
		Msg("Successfully converted Jellyfin item to collection")

	return collection, nil
}

// Helper function to convert Jellyfin item to internal Episode type
func (j *JellyfinClient) convertToEpisode(ctx context.Context, item *jellyfin.BaseItemDto) (models.MediaItem[t.Episode], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	// Validate required fields
	if item == nil {
		return models.MediaItem[t.Episode]{}, fmt.Errorf("cannot convert nil item to episode")
	}

	if item.Id == nil || *item.Id == "" {
		return models.MediaItem[t.Episode]{}, fmt.Errorf("episode is missing required ID field")
	}

	// Safely get name or fallback to empty string
	title := ""
	if item.Name.IsSet() {
		title = *item.Name.Get()
	}

	log.Debug().
		Str("episodeID", *item.Id).
		Str("episodeName", title).
		Msg("Converting Jellyfin item to episode format")

	// Safely handle optional fields
	description := ""
	if item.Overview.IsSet() {
		description = *item.Overview.Get()
	}

	// Safely handle duration
	var duration time.Duration
	if item.RunTimeTicks.IsSet() {
		duration = time.Duration(*item.RunTimeTicks.Get()/10000000) * time.Second
	}

	// Safely handle episode number
	var episodeNumber int64
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

	// Create the basic episode object
	episode := models.MediaItem[t.Episode]{
		Data: t.Episode{
			Details: t.MediaDetails{
				Title:       title,
				Description: description,
				Artwork:     j.getArtworkURLs(item),
				Duration:    duration,
			},
			Number:       episodeNumber,
			SeasonNumber: seasonNumber,
			ShowTitle:    showTitle,
		},
		Type: "episode",
	}

	episode.SetClientInfo(j.ClientID, j.ClientType, *item.Id)

	// Safely set IDs if available
	if item.SeriesId.IsSet() {
		episode.Data.ShowID = *item.SeriesId.Get()
	}

	if item.SeasonId.IsSet() {
		episode.Data.SeasonID = *item.SeasonId.Get()
	}

	// Add air date if available
	if item.PremiereDate.IsSet() {
		episode.Data.Details.ReleaseDate = *item.PremiereDate.Get()
	}

	// Add community rating if available
	if item.CommunityRating.IsSet() {
		episode.Data.Details.Ratings = append(episode.Data.Details.Ratings, t.Rating{
			Source: "jellyfin",
			Value:  float32(*item.CommunityRating.Get()),
		})
	}

	// Add user rating if available
	if item.UserData.IsSet() && item.UserData.Get().Rating.IsSet() {
		episode.Data.Details.UserRating = float32(*item.UserData.Get().Rating.Get())
	}

	// Extract provider IDs if available
	extractProviderIDs(&item.ProviderIds, &episode.Data.Details.ExternalIDs)

	log.Debug().
		Str("episodeID", *item.Id).
		Str("episodeName", episode.Data.Details.Title).
		Int64("episodeNumber", episode.Data.Number).
		Int("seasonNumber", episode.Data.SeasonNumber).
		Msg("Successfully converted Jellyfin item to episode")

	return episode, nil
}

func (j *JellyfinClient) convertToSeries(ctx context.Context, item *jellyfin.BaseItemDto) (models.MediaItem[t.Series], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	// Validate required fields
	if item == nil {
		return models.MediaItem[t.Series]{}, fmt.Errorf("cannot convert nil item to TV show")
	}

	if item.Id == nil || *item.Id == "" {
		return models.MediaItem[t.Series]{}, fmt.Errorf("TV show is missing required ID field")
	}

	// Safely get name or fallback to empty string
	title := ""
	if item.Name.IsSet() {
		title = *item.Name.Get()
	}

	log.Debug().
		Str("showID", *item.Id).
		Str("showName", title).
		Msg("Converting Jellyfin item to TV show format")

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
	var duration time.Duration
	if item.RunTimeTicks.IsSet() {
		duration = time.Duration(*item.RunTimeTicks.Get()/10000000) * time.Second
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

	// Build TV show object
	show := models.MediaItem[t.Series]{
		Data: t.Series{
			Details: t.MediaDetails{
				Title:       title,
				Description: description,
				ReleaseYear: releaseYear,
				Genres:      genres,
				Artwork:     j.getArtworkURLs(item),
				Duration:    duration,
			},
			Status:      status,
			SeasonCount: seasonCount,
		},
		Type: "tvshow",
	}

	// ClientID:   j.ClientID,
	// 			ExternalID: *item.Id,
	// 			ClientType: string(j.ClientType),
	show.SetClientInfo(j.ClientID, j.ClientType, *item.Id)

	// Set SeriesStudio if available
	if item.SeriesStudio.IsSet() {
		show.Data.Network = *item.SeriesStudio.Get()
	}

	// Extract provider IDs if available
	extractProviderIDs(&item.ProviderIds, &show.Data.Details.ExternalIDs)

	// Set ratings if available
	if item.CommunityRating.IsSet() {
		show.Data.Details.Ratings = append(show.Data.Details.Ratings, t.Rating{
			Source: "jellyfin",
			Value:  float32(*item.CommunityRating.Get()),
		})
	}

	// Set user rating if available
	if item.UserData.IsSet() && item.UserData.Get().Rating.IsSet() {
		show.Data.Details.UserRating = float32(*item.UserData.Get().Rating.Get())
	}

	log.Debug().
		Str("showID", *item.Id).
		Str("showName", show.Data.Details.Title).
		Int("seasonCount", show.Data.SeasonCount).
		Msg("Successfully converted Jellyfin item to TV show")

	return show, nil
}

// Helper function to convert Jellyfin item to internal Movie type
func (j *JellyfinClient) convertToMovie(ctx context.Context, item *jellyfin.BaseItemDto) (models.MediaItem[t.Movie], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	// Validate required fields
	if item == nil {
		return models.MediaItem[t.Movie]{}, fmt.Errorf("cannot convert nil item to movie")
	}

	if item.Id == nil || *item.Id == "" {
		return models.MediaItem[t.Movie]{}, fmt.Errorf("movie is missing required ID field")
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
	var duration time.Duration
	if item.RunTimeTicks.IsSet() {
		duration = time.Duration(*item.RunTimeTicks.Get()/10000000) * time.Second
	}

	// Initialize ratings
	ratings := t.Ratings{}

	// Safely add community rating if available
	if item.CommunityRating.IsSet() {
		ratings = append(ratings, t.Rating{
			Source: "jellyfin",
			Value:  float32(*item.CommunityRating.Get()),
		})
	}

	// Build movie object
	movie := models.MediaItem[t.Movie]{
		Data: t.Movie{
			Details: t.MediaDetails{
				Title:         title,
				Description:   description,
				ReleaseDate:   releaseDate,
				ReleaseYear:   releaseYear,
				ContentRating: contentRating,
				Genres:        genres,
				Artwork:       j.getArtworkURLs(item),
				Duration:      duration,
				Ratings:       ratings,
			},
		},
		Type: t.MediaTypeMovie,
	}

	movie.SetClientInfo(j.ClientID, j.ClientType, *item.Id)

	// Set user rating if available
	if item.UserData.IsSet() && item.UserData.Get().Rating.IsSet() {
		movie.Data.Details.UserRating = float32(*item.UserData.Get().Rating.Get())
	} else {
		log.Debug().
			Str("movieID", *item.Id).
			Msg("Movie has no user data, skipping user rating")
	}

	// Extract provider IDs if available
	extractProviderIDs(&item.ProviderIds, &movie.Data.Details.ExternalIDs)

	log.Debug().
		Str("movieID", *item.Id).
		Str("movieTitle", movie.Data.Details.Title).
		Int("year", movie.Data.Details.ReleaseYear).
		Msg("Successfully converted Jellyfin item to movie")

	return movie, nil
}

func (j *JellyfinClient) convertToAlbum(ctx context.Context, item *jellyfin.BaseItemDto) (models.MediaItem[t.Album], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	// Validate required fields
	if item == nil {
		return models.MediaItem[t.Album]{}, fmt.Errorf("cannot convert nil item to album")
	}

	if item.Id == nil || *item.Id == "" {
		return models.MediaItem[t.Album]{}, fmt.Errorf("album is missing required ID field")
	}

	// Safely get name or fallback to empty string
	title := ""
	if item.Name.IsSet() {
		title = *item.Name.Get()
	}

	log.Debug().
		Str("albumID", *item.Id).
		Str("albumName", title).
		Msg("Converting Jellyfin item to album format")

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

	// Build album object
	album := models.MediaItem[t.Album]{
		Data: t.Album{
			Details: t.MediaDetails{
				Title:       title,
				Description: description,
				ReleaseYear: releaseYear,
				Artwork:     j.getArtworkURLs(item),
			},
			ArtistName: artistName,
			TrackCount: trackCount,
		},
		Type: "album",
	}

	album.SetClientInfo(j.ClientID, j.ClientType, *item.Id)

	// Add user rating if available
	if item.UserData.IsSet() && item.UserData.Get().Rating.IsSet() {
		album.Data.Details.UserRating = float32(*item.UserData.Get().Rating.Get())
	}

	// Add community rating if available
	if item.CommunityRating.IsSet() {
		album.Data.Details.Ratings = append(album.Data.Details.Ratings, t.Rating{
			Source: "jellyfin",
			Value:  float32(*item.CommunityRating.Get()),
		})
	}

	// Extract provider IDs if available
	extractProviderIDs(&item.ProviderIds, &album.Data.Details.ExternalIDs)

	log.Debug().
		Str("albumID", *item.Id).
		Str("albumName", album.Data.Details.Title).
		Str("artistName", album.Data.ArtistName).
		Int("trackCount", album.Data.TrackCount).
		Msg("Successfully converted Jellyfin item to album")

	return album, nil
}

// Helper function to convert Jellyfin item to internal Season type
func (j *JellyfinClient) convertToSeason(ctx context.Context, item *jellyfin.BaseItemDto) (models.MediaItem[t.Season], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	// Validate required fields
	if item == nil {
		return models.MediaItem[t.Season]{}, fmt.Errorf("cannot convert nil item to season")
	}

	if item.Id == nil || *item.Id == "" {
		return models.MediaItem[t.Season]{}, fmt.Errorf("season is missing required ID field")
	}

	// Safely get name or fallback to empty string
	title := ""
	if item.Name.IsSet() {
		title = *item.Name.Get()
	}

	log.Debug().
		Str("seasonID", *item.Id).
		Str("seasonName", title).
		Msg("Converting Jellyfin item to season format")

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

	// Safely handle series ID
	seriesID := ""
	if item.SeriesId.IsSet() {
		seriesID = *item.SeriesId.Get()
	}

	// Build season object
	season := models.MediaItem[t.Season]{
		Data: t.Season{
			Details: t.MediaDetails{
				Title:       title,
				Description: description,
				Artwork:     j.getArtworkURLs(item),
			},
			Number:       seasonNumber,
			EpisodeCount: episodeCount,
			SeriesName:   seriesName,
			SeriesID:     seriesID,
		},
		Type: "season",
	}

	season.SetClientInfo(j.ClientID, j.ClientType, *item.Id)

	// Add release year if available
	if item.ProductionYear.IsSet() {
		season.Data.Details.ReleaseYear = int(*item.ProductionYear.Get())
	}

	// Add premiere date if available
	if item.PremiereDate.IsSet() {
		season.Data.Details.ReleaseDate = *item.PremiereDate.Get()
	}

	// Add community rating if available
	if item.CommunityRating.IsSet() {
		season.Data.Details.Ratings = append(season.Data.Details.Ratings, t.Rating{
			Source: "jellyfin",
			Value:  float32(*item.CommunityRating.Get()),
		})
	}

	// Add user rating if available
	if item.UserData.IsSet() && item.UserData.Get().Rating.IsSet() {
		season.Data.Details.UserRating = float32(*item.UserData.Get().Rating.Get())
	}

	// Extract provider IDs if available
	extractProviderIDs(&item.ProviderIds, &season.Data.Details.ExternalIDs)

	log.Debug().
		Str("seasonID", *item.Id).
		Str("seasonName", season.Data.Details.Title).
		Int("seasonNumber", season.Data.Number).
		Int("episodeCount", season.Data.EpisodeCount).
		Msg("Successfully converted Jellyfin item to season")

	return season, nil
}

func (j *JellyfinClient) convertByItemType(ctx context.Context, item *jellyfin.BaseItemDto) (models.MediaItem[t.MediaData], error) {

	if item == nil {
		return models.MediaItem[t.MediaData]{}, fmt.Errorf("cannot convert nil item to item")
	}

	var result models.MediaItem[t.MediaData]
	// var err error

	switch *item.Type {
	case jellyfin.BASEITEMKIND_MOVIE:
		movie, convErr := j.convertToMovie(ctx, item)
		if convErr != nil {
			return models.MediaItem[t.MediaData]{}, convErr
		}
		// Convert to generic MediaData
		result.Type = movie.Type
		result.ClientID = movie.ClientID
		result.ExternalID = movie.ExternalID
		result.ClientType = movie.ClientType
		result.Data = movie.Data

	case jellyfin.BASEITEMKIND_EPISODE:
		episode, convErr := j.convertToEpisode(ctx, item)
		if convErr != nil {
			return models.MediaItem[t.MediaData]{}, convErr
		}
		result.Type = episode.Type
		result.ClientID = episode.ClientID
		result.ExternalID = episode.ExternalID
		result.ClientType = episode.ClientType
		result.Data = episode.Data
	case jellyfin.BASEITEMKIND_MUSIC_ALBUM:
		album, convErr := j.convertToAlbum(ctx, item)
		if convErr != nil {
			return models.MediaItem[t.MediaData]{}, convErr
		}
		result.Type = album.Type
		result.ClientID = album.ClientID
		result.ExternalID = album.ExternalID
		result.ClientType = album.ClientType
		result.Data = album.Data
	case jellyfin.BASEITEMKIND_SERIES:
		tvShow, convErr := j.convertToSeries(ctx, item)
		if convErr != nil {
			return models.MediaItem[t.MediaData]{}, convErr
		}
		result.Type = tvShow.Type
		result.ClientID = tvShow.ClientID
		result.ExternalID = tvShow.ExternalID
		result.ClientType = tvShow.ClientType
		result.Data = tvShow.Data
	case jellyfin.BASEITEMKIND_SEASON:
		season, convErr := j.convertToSeason(ctx, item)
		if convErr != nil {
			return models.MediaItem[t.MediaData]{}, convErr
		}
		result.Type = season.Type
		result.ClientID = season.ClientID
		result.ExternalID = season.ExternalID
		result.ClientType = season.ClientType
		result.Data = season.Data
	case jellyfin.BASEITEMKIND_AUDIO:
		artist, convErr := j.convertToArtist(ctx, item)
		if convErr != nil {
			return models.MediaItem[t.MediaData]{}, convErr
		}
		result.Type = artist.Type
		result.ClientID = artist.ClientID
		result.ExternalID = artist.ExternalID
		result.ClientType = artist.ClientType
		result.Data = artist.Data
	case jellyfin.BASEITEMKIND_PLAYLIST:
		playlist, convErr := j.convertToPlaylist(ctx, item)
		if convErr != nil {
			return models.MediaItem[t.MediaData]{}, convErr
		}
		result.Type = playlist.Type
		result.ClientID = playlist.ClientID
		result.ExternalID = playlist.ExternalID
		result.ClientType = playlist.ClientType
		result.Data = playlist.Data
	case jellyfin.BASEITEMKIND_COLLECTION_FOLDER:
		collection, convErr := j.convertToCollection(ctx, item)
		if convErr != nil {
			return models.MediaItem[t.MediaData]{}, convErr
		}
		result.Type = collection.Type
		result.ClientID = collection.ClientID
		result.ExternalID = collection.ExternalID
		result.ClientType = collection.ClientType
		result.Data = collection.Data
	default:
		return models.MediaItem[t.MediaData]{}, fmt.Errorf("item type not supported")
	}
	return result, nil

}

//convertToArtist

func (j *JellyfinClient) convertToArtist(ctx context.Context, item *jellyfin.BaseItemDto) (models.MediaItem[t.Artist], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	// Validate required fields
	if item == nil {
		return models.MediaItem[t.Artist]{}, fmt.Errorf("cannot convert nil item to artist")
	}

	if item.Id == nil || *item.Id == "" {
		return models.MediaItem[t.Artist]{}, fmt.Errorf("artist is missing required ID field")
	}

	// Safely get name or fallback to empty string
	name := ""
	if item.Name.IsSet() {
		name = *item.Name.Get()
	}

	log.Debug().
		Str("artistID", *item.Id).
		Str("artistName", name).
		Msg("Converting Jellyfin item to artist format")

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

	// Build artist object
	artist := models.MediaItem[t.Artist]{
		Data: t.Artist{
			Details: t.MediaDetails{
				Title:       name,
				Description: description,
				Genres:      genres,
				Artwork:     j.getArtworkURLs(item),
			},
			AlbumCount: albumCount,
		},
		Type: "artist",
	}

	artist.SetClientInfo(j.ClientID, j.ClientType, *item.Id)

	// Add user rating if available
	if item.UserData.IsSet() && item.UserData.Get().Rating.IsSet() {
		artist.Data.Details.UserRating = float32(*item.UserData.Get().Rating.Get())
	}

	// Add community rating if available
	if item.CommunityRating.IsSet() {
		artist.Data.Details.Ratings = append(artist.Data.Details.Ratings, t.Rating{
			Source: "jellyfin",
			Value:  float32(*item.CommunityRating.Get()),
		})
	}

	// Extract provider IDs if available
	extractProviderIDs(&item.ProviderIds, &artist.Data.Details.ExternalIDs)

	log.Debug().
		Str("artistID", *item.Id).
		Str("artistName", artist.Data.Details.Title).
		Int("albumCount", len(artist.Data.Albums)).
		Msg("Successfully converted Jellyfin item to artist")

	return artist, nil
}

func (j *JellyfinClient) convertToPlaylist(ctx context.Context, item *jellyfin.BaseItemDto) (models.MediaItem[t.Playlist], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	// Validate required fields
	if item == nil {
		return models.MediaItem[t.Playlist]{}, fmt.Errorf("cannot convert nil item to playlist")
	}

	if item.Id == nil || *item.Id == "" {
		return models.MediaItem[t.Playlist]{}, fmt.Errorf("playlist is missing required ID field")
	}

	// Safely get name or fallback to empty string
	title := ""
	if item.Name.IsSet() {
		title = *item.Name.Get()
	}

	log.Debug().
		Str("playlistID", *item.Id).
		Str("playlistName", title).
		Msg("Converting Jellyfin item to playlist format")

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

	// Build playlist object
	playlist := models.MediaItem[t.Playlist]{
		Data: t.Playlist{
			Details: t.MediaDetails{
				Title:       title,
				Description: description,
				Artwork:     j.getArtworkURLs(item),
			},
			ItemCount: itemCount,
			IsPublic:  true, // Assume public by default in Jellyfin
		},
		Type: "playlist",
	}

	playlist.SetClientInfo(j.ClientID, j.ClientType, *item.Id)

	// Add community rating if available
	if item.CommunityRating.IsSet() {
		playlist.Data.Details.Ratings = append(playlist.Data.Details.Ratings, t.Rating{
			Source: "jellyfin",
			Value:  float32(*item.CommunityRating.Get()),
		})
	}

	// Add user rating if available
	if item.UserData.IsSet() && item.UserData.Get().Rating.IsSet() {
		playlist.Data.Details.UserRating = float32(*item.UserData.Get().Rating.Get())
	}

	log.Debug().
		Str("playlistID", *item.Id).
		Str("playlistName", playlist.Data.Details.Title).
		Int("itemCount", playlist.Data.ItemCount).
		Msg("Successfully converted Jellyfin item to playlist")

	return playlist, nil
}

// func (j *JellyfinClient) convertToSeason(ctx context.Context, item *jellyfin.BaseItemDto, showID string) (models.MediaItem[t.Season], error) {
// 	// Get logger from context
// 	log := utils.LoggerFromContext(ctx)
//
// 	// Validate required fields
// 	if item == nil {
// 		return models.MediaItem[t.Season]{}, fmt.Errorf("cannot convert nil item to season")
// 	}
//
// 	if item.Id == nil || *item.Id == "" {
// 		return models.MediaItem[t.Season]{}, fmt.Errorf("season is missing required ID field")
// 	}
//
// 	// Safely get name or fallback to empty string
// 	title := ""
// 	if item.Name.IsSet() {
// 		title = *item.Name.Get()
// 	}
//
// 	log.Debug().
// 		Str("seasonID", *item.Id).
// 		Str("seasonName", title).
// 		Str("showID", showID).
// 		Msg("Converting Jellyfin item to season format")
//
// 	// Safely handle optional fields
// 	description := ""
// 	if item.Overview.IsSet() {
// 		description = *item.Overview.Get()
// 	}
//
// 	// Safely handle season number
// 	seasonNumber := 0
// 	if item.IndexNumber.IsSet() {
// 		seasonNumber = int(*item.IndexNumber.Get())
// 	}
//
// 	// Safely handle episode count
// 	episodeCount := 0
// 	if item.ChildCount.IsSet() {
// 		episodeCount = int(*item.ChildCount.Get())
// 	}
//
// 	// Create the basic season object
// 	season := models.MediaItem[t.Season]{
// 		Data: t.Season{
// 			Details: t.MediaMetadata{
// 				Title:       title,
// 				Description: description,
// 				Artwork:     j.getArtworkURLs(item),
// 			},
//
// 			ParentID:     showID,
// 			Number:       seasonNumber,
// 			EpisodeCount: episodeCount,
// 		},
// 		Type: "season",
// 	}
//
// 	season.SetClientInfo(j.ClientID, j.ClientType, *item.Id)
//
// 	// Safely set release date if available
// 	if item.PremiereDate.IsSet() {
// 		season.Data.ReleaseDate = *item.PremiereDate.Get()
// 	}
//
// 	// Add community rating if available
// 	if item.CommunityRating.IsSet() {
// 		season.Data.Details.Ratings = append(season.Data.Details.Ratings, t.Rating{
// 			Source: "jellyfin",
// 			Value:  float32(*item.CommunityRating.Get()),
// 		})
// 	}
//
// 	// Add user rating if available
// 	if item.UserData.IsSet() && item.UserData.Get().Rating.IsSet() {
// 		season.Data.Details.UserRating = float32(*item.UserData.Get().Rating.Get())
// 	}
//
// 	// Extract provider IDs if available
// 	extractProviderIDs(&item.ProviderIds, &season.Data.Details.ExternalIDs)
//
// 	log.Debug().
// 		Str("seasonID", *item.Id).
// 		Str("seasonName", season.Data.Details.Title).
// 		Int("seasonNumber", season.Data.Number).
// 		Int("episodeCount", season.Data.EpisodeCount).
// 		Msg("Successfully converted Jellyfin item to season")
//
// 	return season, nil
// }
