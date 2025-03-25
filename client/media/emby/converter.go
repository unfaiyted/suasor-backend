package emby

import (
	"context"
	"fmt"
	"time"

	"suasor/client/media/types"
	embyclient "suasor/internal/clients/embyAPI"
	"suasor/utils"
)

func (e *EmbyClient) convertToWatchHistoryItem(item *embyclient.BaseItemDto) (types.WatchHistoryItem[types.MediaData], error) {
	if item == nil || item.UserData == nil {
		return types.WatchHistoryItem[types.MediaData]{}, fmt.Errorf("cannot convert nil item to watch history item")
	}

	// TODO: Needs to properly handle other types
	var watchData types.MediaData
	// TODO: find emby types list
	if item.Type_ == "Movie" {
		watchData = types.Movie{}
	}

	watchItem := types.WatchHistoryItem[types.MediaData]{
		Item: types.MediaItem[types.MediaData]{
			Data: watchData,
		},
		ItemType:        string(item.Type_),
		WatchedAt:       item.UserData.LastPlayedDate,
		IsFavorite:      item.UserData.IsFavorite,
		PlayCount:       item.UserData.PlayCount,
		PositionSeconds: int(item.UserData.PlaybackPositionTicks / 10000000),
	}
	watchItem.Item.SetClientInfo(e.ClientID, e.ClientType, item.Id)

	return watchItem, nil
}

// Helper function to convert Emby item to internal Movie type
func (e *EmbyClient) convertToMovie(ctx context.Context, item *embyclient.BaseItemDto) (types.MediaItem[types.Movie], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	// Validate required fields
	if item == nil {
		return types.MediaItem[types.Movie]{}, fmt.Errorf("cannot convert nil item to movie")
	}

	if item.Id == "" {
		return types.MediaItem[types.Movie]{}, fmt.Errorf("movie is missing required ID field")
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
	movie := types.MediaItem[types.Movie]{
		Data: types.Movie{
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
		},
		Type: "movie",
	}
	movie.SetClientInfo(e.ClientID, e.ClientType, item.Id)

	// Only set UserRating if UserData is not nil
	if item.UserData != nil {
		movie.Data.Details.UserRating = float32(item.UserData.Rating)
	} else {
		log.Debug().
			Str("movieID", item.Id).
			Msg("Movie has no user data, skipping user rating")
	}

	// Extract provider IDs if available
	if item.ProviderIds != nil {
		ids := *item.ProviderIds
		if imdbID, ok := ids["Imdb"]; ok {
			movie.Data.Details.ExternalIDs.AddOrUpdate("imdb", imdbID)
		}
		if tmdbID, ok := ids["Tmdb"]; ok {
			movie.Data.Details.ExternalIDs.AddOrUpdate("tmdb", tmdbID)
		}
	}

	log.Debug().
		Str("movieID", item.Id).
		Str("movieTitle", movie.Data.Details.Title).
		Int("year", movie.Data.Details.ReleaseYear).
		Msg("Successfully converted Emby item to movie")

	return movie, nil
}

func (e *EmbyClient) convertToTrack(item *embyclient.BaseItemDto) (types.MediaItem[types.Track], error) {
	if item == nil {
		return types.MediaItem[types.Track]{}, fmt.Errorf("cannot convert nil item to music track")
	}

	track := types.MediaItem[types.Track]{
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

func (e *EmbyClient) convertToMusicArtist(item *embyclient.BaseItemDto) (types.MediaItem[types.Artist], error) {
	if item == nil {
		return types.MediaItem[types.Artist]{}, fmt.Errorf("cannot convert nil item to music artist")
	}

	artist := types.MediaItem[types.Artist]{
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

func (e *EmbyClient) convertToAlbum(item *embyclient.BaseItemDto) (types.MediaItem[types.Album], error) {
	if item == nil {
		return types.MediaItem[types.Album]{}, fmt.Errorf("cannot convert nil item to music album")
	}

	album := types.MediaItem[types.Album]{
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

func (e *EmbyClient) convertToPlaylist(item *embyclient.BaseItemDto) (types.MediaItem[types.Playlist], error) {
	if item == nil {
		return types.MediaItem[types.Playlist]{}, fmt.Errorf("cannot convert nil item to playlist")
	}

	playlist := types.MediaItem[types.Playlist]{
		Data: types.Playlist{
			Details: types.MediaMetadata{
				Title:       item.Name,
				Description: item.Overview,
				Artwork:     e.getArtworkURLs(item),
			},
			ItemCount: int(item.ChildCount),
			IsPublic:  true, // Assume public by default in Emby
		},
		Type: "playlist",
	}
	playlist.SetClientInfo(e.ClientID, e.ClientType, item.Id)

	return playlist, nil
}

// Note: This is a duplicate of the first function and should be removed

func (e *EmbyClient) convertToTVShow(item *embyclient.BaseItemDto) (types.MediaItem[types.TVShow], error) {
	if item == nil {
		return types.MediaItem[types.TVShow]{}, fmt.Errorf("cannot convert nil item to TV show")
	}

	show := types.MediaItem[types.TVShow]{
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

func (e *EmbyClient) convertToSeason(item *embyclient.BaseItemDto, showID string) (types.MediaItem[types.Season], error) {
	if item == nil {
		return types.MediaItem[types.Season]{}, fmt.Errorf("cannot convert nil item to season")
	}

	season := types.MediaItem[types.Season]{
		Data: types.Season{
			Details: types.MediaMetadata{
				Title:       item.Name,
				Description: item.Overview,
				Artwork:     e.getArtworkURLs(item),
			},
			ParentID:     showID,
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

func (e *EmbyClient) convertToEpisode(item *embyclient.BaseItemDto, showID string, seasonNumber int) (types.MediaItem[types.Episode], error) {
	if item == nil {
		return types.MediaItem[types.Episode]{}, fmt.Errorf("cannot convert nil item to episode")
	}

	episode := types.MediaItem[types.Episode]{
		Data: types.Episode{
			Details: types.MediaMetadata{
				Title:       item.Name,
				Description: item.Overview,
				Artwork:     e.getArtworkURLs(item),
				Duration:    time.Duration(item.RunTimeTicks/10000000) * time.Second,
			},
			Number:       int64(item.IndexNumber),
			ShowID:       showID,
			SeasonID:     item.SeasonId,
			SeasonNumber: seasonNumber,
			ShowTitle:    item.SeriesName,
		},
		Type: "episode",
	}
	episode.SetClientInfo(e.ClientID, e.ClientType, item.Id)

	// Add external IDs
	if item.ProviderIds != nil {
		ids := *item.ProviderIds
		if imdbID, ok := ids["Imdb"]; ok {
			episode.Data.Details.ExternalIDs.AddOrUpdate("imdb", imdbID)
		}
		if tmdbID, ok := ids["Tmdb"]; ok {
			episode.Data.Details.ExternalIDs.AddOrUpdate("tmdb", tmdbID)
		}
		if tvdbID, ok := ids["Tvdb"]; ok {
			episode.Data.Details.ExternalIDs.AddOrUpdate("tvdb", tvdbID)
		}
	}

	return episode, nil
}
