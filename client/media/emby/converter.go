package emby

import (
	"context"
	"fmt"
	"time"

	"suasor/client/media/types"
	embyclient "suasor/internal/clients/embyAPI"
	"suasor/types/models"
	"suasor/utils"
)

func (e *EmbyClient) convertToWatchHistoryItem(ctx context.Context, item *embyclient.BaseItemDto) (models.MediaPlayHistory[types.MediaData], error) {
	if item == nil || item.UserData == nil {
		return models.MediaPlayHistory[types.MediaData]{}, fmt.Errorf("cannot convert nil item to watch history item")
	}

	// TODO: Needs to properly handle other types
	var watchData models.MediaItem[types.MediaData]
	watchData.SetClientInfo(e.ClientID, e.ClientType, item.Id)
	// TODO: find emby types list
	if item.Type_ == "Movie" {
		mediaItemMovie, err := e.convertToMovie(ctx, item)
		if err != nil {
			return models.MediaPlayHistory[types.MediaData]{}, err
		}
		watchData.SetData(&watchData, mediaItemMovie.GetData())
	} else if item.Type_ == "Episode" {
		mediaItemEpisode, err := e.convertToEpisode(item)
		if err != nil {
			return models.MediaPlayHistory[types.MediaData]{}, err
		}
		watchData.SetData(&watchData, mediaItemEpisode.GetData())
	} else if item.Type_ == "Audio" {
		mediaItemMusic, err := e.convertToTrack(item)
		if err != nil {
			return models.MediaPlayHistory[types.MediaData]{}, err
		}
		watchData.SetData(&watchData, mediaItemMusic.GetData())
	} else if item.Type_ == "Playlist" {
		mediaItemPlaylist, err := e.convertToPlaylist(item)
		if err != nil {
			return models.MediaPlayHistory[types.MediaData]{}, err
		}
		watchData.SetData(&watchData, mediaItemPlaylist.GetData())
	} else if item.Type_ == "Series" {
		mediaItemTVShow, err := e.convertToTVShow(item)
		if err != nil {
			return models.MediaPlayHistory[types.MediaData]{}, err
		}
		watchData.SetData(&watchData, mediaItemTVShow.GetData())
	} else if item.Type_ == "Season" {
		mediaItemSeason, err := e.convertToSeason(item, item.ParentId)
		if err != nil {
			return models.MediaPlayHistory[types.MediaData]{}, err
		}
		watchData.SetData(&watchData, mediaItemSeason.GetData())
	} else if item.Type_ == "Collection" {
		mediaItemCollection, err := e.convertToCollection(item)
		if err != nil {
			return models.MediaPlayHistory[types.MediaData]{}, err
		}
		watchData.SetData(&watchData, mediaItemCollection.GetData())
	}

	watchItem := models.MediaPlayHistory[types.MediaData]{
		Item:            &watchData,
		Type:            string(item.Type_),
		WatchedAt:       item.UserData.LastPlayedDate,
		IsFavorite:      item.UserData.IsFavorite,
		PlayCount:       item.UserData.PlayCount,
		PositionSeconds: int(item.UserData.PlaybackPositionTicks / 10000000),
	}
	watchItem.Item.SetClientInfo(e.ClientID, e.ClientType, item.Id)

	return watchItem, nil
}

// Helper function to convert Emby item to internal Movie type
func (e *EmbyClient) convertToMovie(ctx context.Context, item *embyclient.BaseItemDto) (models.MediaItem[types.Movie], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	// Validate required fields
	if item == nil {
		return models.MediaItem[types.Movie]{}, fmt.Errorf("cannot convert nil item to movie")
	}

	if item.Id == "" {
		return models.MediaItem[types.Movie]{}, fmt.Errorf("movie is missing required ID field")
	}

	log.Debug().
		Str("movieID", item.Id).
		Str("movieName", item.Name).
		Int32("releaseYear", item.ProductionYear).
		Str("releaseDate", item.PremiereDate.Format("2006-01-02")).
		Msg("Converting Emby item to movie format")

	// Handle empty or non-numeric rating safely
	if item.OfficialRating != "" {
		// Try to convert to integer, but don't fail if it's not a number
		log.Debug().
			Str("movieID", item.Id).
			Str("rating", item.OfficialRating).
			Msg("Content rating found.")
	}

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

	// Build movie object with safe handling of optional fields
	movie := types.Movie{
		Details: types.MediaMetadata{
			Title:       item.Name,
			Description: item.Overview,
			ReleaseDate: item.PremiereDate,
			ReleaseYear: releaseYear,
			Genres:      item.Genres,
			Artwork:     e.getArtworkURLs(item),
			Duration:    time.Duration(item.RunTimeTicks/10000000) * time.Second,
			Ratings: types.Ratings{
				types.Rating{
					Source: "emby",
					Value:  float32(item.CommunityRating),
				},
			},
		},
	}

	mediaItemMovie := models.MediaItem[types.Movie]{
		Data: movie,
		Type: movie.GetMediaType(),
	}
	mediaItemMovie.SetClientInfo(e.ClientID, e.ClientType, item.Id)

	// Only set UserRating if UserData is not nil
	if item.UserData != nil {
		movie.Details.UserRating = float32(item.UserData.Rating)
	} else {
		log.Debug().
			Str("movieID", item.Id).
			Msg("Movie has no user data, skipping user rating")
	}

	// Extract provider IDs if available
	if item.ProviderIds != nil {
		ids := *item.ProviderIds
		if imdbID, ok := ids["Imdb"]; ok {
			movie.Details.ExternalIDs.AddOrUpdate("imdb", imdbID)
		}
		if tmdbID, ok := ids["Tmdb"]; ok {
			movie.Details.ExternalIDs.AddOrUpdate("tmdb", tmdbID)
		}
	}

	log.Debug().
		Str("movieID", item.Id).
		Str("movieTitle", movie.Details.Title).
		Int("year", movie.Details.ReleaseYear).
		Msg("Successfully converted Emby item to movie")

	return mediaItemMovie, nil
}

func (e *EmbyClient) convertToTrack(item *embyclient.BaseItemDto) (models.MediaItem[types.Track], error) {
	if item == nil {
		return models.MediaItem[types.Track]{}, fmt.Errorf("cannot convert nil item to music track")
	}

	track := models.MediaItem[types.Track]{
		Data: types.Track{
			Details: types.MediaMetadata{
				Title:       item.Name,
				Description: item.Overview,
				Duration:    time.Duration(item.RunTimeTicks/10000000) * time.Second,
				Artwork:     e.getArtworkURLs(item),
			},
			Number:    int(item.IndexNumber),
			AlbumID:   item.AlbumId,
			AlbumName: item.Album,
		},
		Type: "track",
	}
	track.SetClientInfo(e.ClientID, e.ClientType, item.Id)

	// Add artist information if available
	if len(item.ArtistItems) > 0 {
		track.Data.ArtistID = item.ArtistItems[0].Id
		track.Data.ArtistName = item.ArtistItems[0].Name
	}

	// Extract provider IDs
	if item.ProviderIds != nil {
		ids := *item.ProviderIds
		if musicbrainzID, ok := ids["MusicBrainzTrack"]; ok {
			track.Data.Details.ExternalIDs.AddOrUpdate("musicbrainz", musicbrainzID)
		}
	}

	return track, nil
}

func (e *EmbyClient) convertToMusicArtist(item *embyclient.BaseItemDto) (models.MediaItem[types.Artist], error) {
	if item == nil {
		return models.MediaItem[types.Artist]{}, fmt.Errorf("cannot convert nil item to music artist")
	}

	artist := models.MediaItem[types.Artist]{
		Data: types.Artist{
			Details: types.MediaMetadata{
				Title:       item.Name,
				Description: item.Overview,
				Artwork:     e.getArtworkURLs(item),
				Genres:      item.Genres,
			},
		},
		Type: "artist",
	}
	artist.SetClientInfo(e.ClientID, e.ClientType, item.Id)

	// Extract provider IDs if available
	if item.ProviderIds != nil {
		ids := *item.ProviderIds
		if musicbrainzID, ok := ids["MusicBrainzArtist"]; ok {
			artist.Data.Details.ExternalIDs.AddOrUpdate("musicbrainz", musicbrainzID)
		}
	}

	return artist, nil
}

func (e *EmbyClient) convertToAlbum(item *embyclient.BaseItemDto) (models.MediaItem[types.Album], error) {
	if item == nil {
		return models.MediaItem[types.Album]{}, fmt.Errorf("cannot convert nil item to music album")
	}

	album := models.MediaItem[types.Album]{
		Data: types.Album{
			Details: types.MediaMetadata{
				Title:       item.Name,
				Description: item.Overview,
				ReleaseYear: int(item.ProductionYear),
				Genres:      item.Genres,
				Artwork:     e.getArtworkURLs(item),
			},
			ArtistName: item.AlbumArtist,
			TrackCount: int(item.ChildCount),
		},
		Type: "album",
	}
	album.SetClientInfo(e.ClientID, e.ClientType, item.Id)

	// Extract provider IDs
	if item.ProviderIds != nil {
		ids := *item.ProviderIds
		if musicbrainzID, ok := ids["MusicBrainzAlbum"]; ok {
			album.Data.Details.ExternalIDs.AddOrUpdate("musicbrainz", musicbrainzID)
		}
	}

	return album, nil
}

func (e *EmbyClient) convertToPlaylist(item *embyclient.BaseItemDto) (models.MediaItem[types.Playlist], error) {
	if item == nil {
		return models.MediaItem[types.Playlist]{}, fmt.Errorf("cannot convert nil item to playlist")
	}

	playlist := types.Playlist{
		Details: types.MediaMetadata{
			Title:       item.Name,
			Description: item.Overview,
			Artwork:     e.getArtworkURLs(item),
		},
		ItemCount: int(item.ChildCount),
		IsPublic:  true, // Assume public by default in Emby
	}
	mediaItemPlaylist := models.MediaItem[types.Playlist]{
		Data: playlist,
		Type: playlist.GetMediaType(),
	}
	mediaItemPlaylist.SetClientInfo(e.ClientID, e.ClientType, item.Id)

	return mediaItemPlaylist, nil
}

// Note: This is a duplicate of the first function and should be removed

func (e *EmbyClient) convertToTVShow(item *embyclient.BaseItemDto) (models.MediaItem[types.TVShow], error) {
	if item == nil {
		return models.MediaItem[types.TVShow]{}, fmt.Errorf("cannot convert nil item to TV show")
	}

	show := models.MediaItem[types.TVShow]{
		Data: types.TVShow{
			Details: types.MediaMetadata{
				Title:       item.Name,
				Description: item.Overview,
				ReleaseYear: int(item.ProductionYear),
				Genres:      item.Genres,
				Artwork:     e.getArtworkURLs(item),
				Duration:    time.Duration(item.RunTimeTicks/10000000) * time.Second,
			},
			SeasonCount: int(item.ChildCount),
			Status:      item.Status,
			Network:     item.SeriesStudio,
		},
		Type: "tvshow",
	}
	show.SetClientInfo(e.ClientID, e.ClientType, item.Id)

	// Extract provider IDs if available
	if item.ProviderIds != nil {
		ids := *item.ProviderIds
		if imdbID, ok := ids["Imdb"]; ok {
			show.Data.Details.ExternalIDs.AddOrUpdate("imdb", imdbID)
		}
		if tmdbID, ok := ids["Tmdb"]; ok {
			show.Data.Details.ExternalIDs.AddOrUpdate("tmdb", tmdbID)
		}
		if tvdbID, ok := ids["Tvdb"]; ok {
			show.Data.Details.ExternalIDs.AddOrUpdate("tvdb", tvdbID)
		}
	}

	return show, nil
}

func (e *EmbyClient) convertToSeason(item *embyclient.BaseItemDto, showID string) (models.MediaItem[types.Season], error) {
	if item == nil {
		return models.MediaItem[types.Season]{}, fmt.Errorf("cannot convert nil item to season")
	}

	season := models.MediaItem[types.Season]{
		Data: types.Season{
			Details: types.MediaMetadata{
				Title:       item.Name,
				Description: item.Overview,
				Artwork:     e.getArtworkURLs(item),
			},
			SeriesID:     showID,
			Number:       int(item.IndexNumber),
			EpisodeCount: int(item.ChildCount),
		},
		Type: "season",
	}
	season.SetClientInfo(e.ClientID, e.ClientType, item.Id)

	if !item.PremiereDate.IsZero() {
		season.Data.ReleaseDate = item.PremiereDate
	}

	return season, nil
}

func (e *EmbyClient) convertToEpisode(item *embyclient.BaseItemDto) (models.MediaItem[types.Episode], error) {
	if item == nil {
		return models.MediaItem[types.Episode]{}, fmt.Errorf("cannot convert nil item to episode")
	}

	episode := types.Episode{
		Details: types.MediaMetadata{
			Title:       item.Name,
			Description: item.Overview,
			Artwork:     e.getArtworkURLs(item),
			Duration:    time.Duration(item.RunTimeTicks/10000000) * time.Second,
		},
		Number:       int64(item.IndexNumber),
		ShowID:       item.SeriesId,
		SeasonID:     item.SeasonId,
		SeasonNumber: int(item.ParentIndexNumber),
		ShowTitle:    item.SeriesName,
	}

	mediaItemEpisde := models.MediaItem[types.Episode]{
		Data: episode,
		Type: episode.GetMediaType(),
	}
	mediaItemEpisde.SetClientInfo(e.ClientID, e.ClientType, item.Id)

	// Add external IDs
	if item.ProviderIds != nil {
		ids := *item.ProviderIds
		if imdbID, ok := ids["Imdb"]; ok {
			episode.Details.ExternalIDs.AddOrUpdate("imdb", imdbID)
		}
		if tmdbID, ok := ids["Tmdb"]; ok {
			episode.Details.ExternalIDs.AddOrUpdate("tmdb", tmdbID)
		}
		if tvdbID, ok := ids["Tvdb"]; ok {
			episode.Details.ExternalIDs.AddOrUpdate("tvdb", tvdbID)
		}
	}

	return mediaItemEpisde, nil
}
