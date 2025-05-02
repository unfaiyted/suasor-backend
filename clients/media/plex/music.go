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

// GetMusic retrieves music tracks from Plex
func (c *PlexClient) GetMusicTracks(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Track], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Msg("Retrieving music tracks from Plex server")

	// Find the music library section
	log.Debug().Msg("Finding music library section")
	musicSectionKey, err := c.findLibrarySectionByType(ctx, "artist")
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Msg("Failed to find music library section")
		return nil, err
	}

	if musicSectionKey == "" {
		log.Info().
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Msg("No music library section found in Plex")
		return nil, nil
	}

	// For tracks, we need to traverse the hierarchy: artists > albums > tracks
	sectionKey, _ := strconv.Atoi(musicSectionKey)

	log.Debug().
		Int("sectionKey", sectionKey).
		Msg("Making API request to Plex server for music artists")
	res, err := c.plexAPI.Library.GetLibraryItems(ctx, operations.GetLibraryItemsRequest{SectionKey: sectionKey})
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Int("sectionKey", sectionKey).
			Msg("Failed to get music artists from Plex")
		return nil, fmt.Errorf("failed to get music artists: %w", err)
	}

	var tracks []*models.MediaItem[*types.Track]

	// Loop through artists
	if res.Object.MediaContainer != nil && res.Object.MediaContainer.Metadata != nil {
		log.Debug().
			Int("artistCount", len(res.Object.MediaContainer.Metadata)).
			Msg("Processing artists to find music tracks")

		for _, artist := range res.Object.MediaContainer.Metadata {
			artistKey, _ := strconv.Atoi(artist.RatingKey)
			float64ArtistKey := float64(artistKey)

			log.Debug().
				Str("artistID", artist.RatingKey).
				Str("artistName", artist.Title).
				Msg("Getting albums for artist")

			albumsRes, err := c.plexAPI.Library.GetMetadataChildren(ctx, float64ArtistKey, plexgo.String("Stream"))
			if err != nil {
				log.Warn().
					Err(err).
					Str("artistID", artist.RatingKey).
					Str("artistName", artist.Title).
					Msg("Failed to get albums for artist, skipping")
				continue
			}

			// Loop through albums
			if albumsRes.Object.MediaContainer != nil && albumsRes.Object.MediaContainer.Metadata != nil {
				log.Debug().
					Str("artistID", artist.RatingKey).
					Str("artistName", artist.Title).
					Int("albumCount", len(albumsRes.Object.MediaContainer.Metadata)).
					Msg("Processing albums to find tracks")

				for _, album := range albumsRes.Object.MediaContainer.Metadata {
					albumKey, _ := strconv.Atoi(*album.RatingKey)

					float64AlbumKey := float64(albumKey)

					log.Debug().
						Str("albumID", *album.RatingKey).
						Str("albumName", *album.Title).
						Msg("Getting tracks for album")

					tracksRes, err := c.plexAPI.Library.GetMetadataChildren(ctx, float64AlbumKey, plexgo.String("Stream"))
					if err != nil {
						log.Warn().
							Err(err).
							Str("albumID", *album.RatingKey).
							Str("albumName", *album.Title).
							Msg("Failed to get tracks for album, skipping")
						continue
					}

					// Loop through tracks
					if tracksRes.Object.MediaContainer != nil && tracksRes.Object.MediaContainer.Metadata != nil {
						log.Debug().
							Str("albumID", *album.RatingKey).
							Str("albumName", *album.Title).
							Int("trackCount", len(tracksRes.Object.MediaContainer.Metadata)).
							Msg("Processing tracks for album")

						tracks, err := GetChildMediaItemsList[*types.Track](ctx, c, tracksRes.Object.MediaContainer.Metadata)
						if err != nil {
							log.Error().
								Err(err).
								Uint64("clientID", c.GetClientID()).
								Str("clientType", string(c.GetClientType())).
								Str("albumID", *album.RatingKey).
								Str("albumName", *album.Title).
								Msg("Failed to get tracks for album, skipping")
							continue
						}
						// Limit number of tracks for now
						if len(tracks) >= 1000 {
							log.Info().
								Int("trackCount", len(tracks)).
								Msg("Reached track limit (1000), returning results")
							return tracks, nil
						}
					}
				}
			}
		}
	}

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Int("tracksReturned", len(tracks)).
		Msg("Completed GetMusic request")

	return tracks, nil
}

// GetMusicArtists retrieves music artists from Plex
func (c *PlexClient) GetMusicArtists(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Artist], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Msg("Retrieving music artists from Plex server")

	// Find the music library section
	log.Debug().Msg("Finding music library section")
	musicSectionKey, err := c.findLibrarySectionByType(ctx, "artist")
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Msg("Failed to find music library section")
		return nil, err
	}

	if musicSectionKey == "" {
		log.Info().
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Msg("No music library section found in Plex")
		return nil, nil
	}

	sectionKey, _ := strconv.Atoi(musicSectionKey)

	log.Debug().
		Int("sectionKey", sectionKey).
		Msg("Making API request to Plex server for music artists")

	res, err := c.plexAPI.Library.GetLibraryItems(ctx, operations.GetLibraryItemsRequest{SectionKey: sectionKey})
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Int("sectionKey", sectionKey).
			Msg("Failed to get music artists from Plex")
		return nil, fmt.Errorf("failed to get music artists: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil {
		log.Info().
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Msg("No music artists found in Plex")
		return nil, nil
	}

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Int("totalItems", len(res.Object.MediaContainer.Metadata)).
		Msg("Successfully retrieved music artists from Plex")

	artists, err := GetMediaItemList[*types.Artist](ctx, c, res.Object.MediaContainer.Metadata)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Msg("Failed to get music artists from Plex")
		return nil, fmt.Errorf("failed to get music artists: %w", err)
	}

	return artists, nil
}

// GetMusicAlbums retrieves music albums from Plex
func (c *PlexClient) GetMusicAlbums(ctx context.Context, options *types.QueryOptions) ([]*models.MediaItem[*types.Album], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Msg("Retrieving music albums from Plex server")

	// Find the music library section
	log.Debug().Msg("Finding music library section")
	musicSectionKey, err := c.findLibrarySectionByType(ctx, "artist")
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Msg("Failed to find music library section")
		return nil, err
	}

	if musicSectionKey == "" {
		log.Info().
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Msg("No music library section found in Plex")
		return nil, nil
	}

	// For albums, we need to traverse artists first
	sectionKey, _ := strconv.Atoi(musicSectionKey)

	log.Debug().
		Int("sectionKey", sectionKey).
		Msg("Making API request to Plex server for music artists")

	res, err := c.plexAPI.Library.GetLibraryItems(ctx, operations.GetLibraryItemsRequest{SectionKey: sectionKey})
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Int("sectionKey", sectionKey).
			Msg("Failed to get music artists from Plex")
		return nil, fmt.Errorf("failed to get music artists: %w", err)
	}

	artists, err := GetMediaItemList[*types.Artist](ctx, c, res.Object.MediaContainer.Metadata)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Int("sectionKey", sectionKey).
			Msg("Failed to get music albums from Plex")
		return nil, fmt.Errorf("failed to get music albums: %w", err)
	}

	for _, artist := range artists {
		artistKey, _ := strconv.Atoi(artist.SyncClients.GetClientItemID(c.GetClientID()))
		float64ArtistKey := float64(artistKey)

		log.Debug().
			Int("artistID", artistKey).
			Str("artistName", artist.Title).
			Msg("Getting albums for artist")

		albumsRes, err := c.plexAPI.Library.GetMetadataChildren(ctx, float64ArtistKey, plexgo.String("Stream"))
		if err != nil {
			log.Warn().
				Err(err).
				Int("artistID", artistKey).
				Str("artistName", artist.Title).
				Msg("Failed to get albums for artist, skipping")
			continue
		}

		if albumsRes.Object.MediaContainer != nil && albumsRes.Object.MediaContainer.Metadata != nil {
			log.Debug().
				Int("artistID", artistKey).
				Str("artistName", artist.Title).
				Int("albumCount", len(albumsRes.Object.MediaContainer.Metadata)).
				Msg("Processing albums for artist")

			albums, err := GetChildMediaItemsList[*types.Album](ctx, c, albumsRes.Object.MediaContainer.Metadata)
			if err != nil {
				log.Error().
					Err(err).
					Int("artistID", artistKey).
					Str("artistName", artist.Title).
					Msg("Failed to get albums for artist, skipping")
				continue

			}
			return albums, nil
		}

	}
	return nil, nil
}

// GetMusicTrackByID retrieves a specific music track by ID
func (c *PlexClient) GetMusicTrackByID(ctx context.Context, id string) (*models.MediaItem[*types.Track], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("trackID", id).
		Msg("Retrieving specific music track from Plex server")

	ratingKey, _ := strconv.Atoi(id)
	int64RatingKey := int64(ratingKey)

	log.Debug().
		Str("trackID", id).
		Int64("ratingKey", int64RatingKey).
		Msg("Making API request to Plex server for music track")

	res, err := c.plexAPI.Library.GetMediaMetaData(ctx, operations.GetMediaMetaDataRequest{RatingKey: int64RatingKey})
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Str("trackID", id).
			Msg("Failed to get music track from Plex")
		return nil, fmt.Errorf("failed to get music track: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil || len(res.Object.MediaContainer.Metadata) == 0 {
		log.Error().
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Str("trackID", id).
			Msg("Music track not found in Plex")
		return nil, fmt.Errorf("music track not found")
	}

	item := res.Object.MediaContainer.Metadata[0]
	if item.Type != "track" {
		log.Error().
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Str("trackID", id).
			Str("actualType", item.Type).
			Msg("Item retrieved is not a music track")
		return nil, fmt.Errorf("item is not a music track")

	}

	// Get album and artist info
	var albumName string
	var artistID string
	var artistName string

	if item.ParentRatingKey != nil {
		// Get album info
		albumKey, _ := strconv.Atoi(*item.ParentRatingKey)
		int64AlbumKey := int64(albumKey)

		log.Debug().
			Str("albumID", *item.ParentRatingKey).
			Int64("albumKey", int64AlbumKey).
			Msg("Getting parent album information")

		albumRes, err := c.plexAPI.Library.GetMediaMetaData(ctx, operations.GetMediaMetaDataRequest{
			RatingKey: int64AlbumKey,
		})

		if err == nil && albumRes.Object.MediaContainer != nil &&
			albumRes.Object.MediaContainer.Metadata != nil &&
			len(albumRes.Object.MediaContainer.Metadata) > 0 {

			albumName = albumRes.Object.MediaContainer.Metadata[0].Title

			// Get artist info if available
			if albumRes.Object.MediaContainer.Metadata[0].ParentRatingKey != nil {
				artistID = *albumRes.Object.MediaContainer.Metadata[0].ParentRatingKey

				log.Debug().
					Str("artistID", artistID).
					Msg("Getting parent artist information")

				artistKey, _ := strconv.Atoi(artistID)
				int64ArtistKey := int64(artistKey)
				artistRes, err := c.plexAPI.Library.GetMediaMetaData(ctx, operations.GetMediaMetaDataRequest{
					RatingKey: int64ArtistKey,
				})

				if err == nil && artistRes.Object.MediaContainer != nil &&
					artistRes.Object.MediaContainer.Metadata != nil &&
					len(artistRes.Object.MediaContainer.Metadata) > 0 {

					artistName = artistRes.Object.MediaContainer.Metadata[0].Title
					log.Debug().
						Str("artistName", artistName).
						Msg("Retrieved artist name")
				}
			}
		}
	}

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("trackID", id).
		Str("trackTitle", item.Title).
		Str("albumName", albumName).
		Str("artistName", artistName).
		Msg("Successfully retrieved music track from Plex")

	itemTrack, err := GetItemFromMetadata[*types.Track](ctx, c, &item)
	track, err := GetMediaItem[*types.Track](ctx, c, itemTrack, item.RatingKey)

	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Str("trackID", id).
			Msg("Error converting Plex item to track format")
		return nil, fmt.Errorf("error converting track data: %w", err)
	}

	track.Data.AlbumName = albumName
	track.Data.ArtistName = artistName
	track.Data.Number = int(*item.Index)
	// track.Data.AddSyncClient(c.GetClientID(), *item.ParentRatingKey, artistID)
	track.SetClientInfo(c.GetClientID(), c.GetClientType(), item.RatingKey)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("trackID", id).
		Str("trackTitle", track.Data.Details.Title).
		Int("trackNumber", track.Data.Number).
		Msg("Successfully converted music track data")

	return track, nil
}

func (c *PlexClient) GetMusicArtistByID(ctx context.Context, id string) (*models.MediaItem[*types.Artist], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("artistID", id).
		Msg("Retrieving specific music artist from Plex server")

	ratingKey, _ := strconv.Atoi(id)
	int64RatingKey := int64(ratingKey)

	log.Debug().
		Str("artistID", id).
		Int64("ratingKey", int64RatingKey).
		Msg("Making API request to Plex server for music artist")

	res, err := c.plexAPI.Library.GetMediaMetaData(ctx, operations.GetMediaMetaDataRequest{RatingKey: int64RatingKey})
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Str("artistID", id).
			Msg("Failed to get music artist from Plex")
		return nil, fmt.Errorf("failed to get music artist: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil || len(res.Object.MediaContainer.Metadata) == 0 {
		log.Error().
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Str("artistID", id).
			Msg("Music artist not found in Plex")
		return nil, fmt.Errorf("music artist not found")
	}

	item := res.Object.MediaContainer.Metadata[0]
	if item.Type != "artist" {
		log.Error().
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Str("artistID", id).
			Str("actualType", item.Type).
			Msg("Item retrieved is not a music artist")
		return nil, fmt.Errorf("item is not a music artist")
	}

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("artistID", id).
		Str("artistName", item.Title).
		Msg("Successfully retrieved music artist from Plex")

	itemArtist, err := GetItemFromMetadata[*types.Artist](ctx, c, &item)
	artist, err := GetMediaItem[*types.Artist](ctx, c, itemArtist, item.RatingKey)

	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Str("artistID", id).
			Msg("Error converting Plex item to artist format")
		return nil, fmt.Errorf("error converting artist data: %w", err)
	}

	// Set client info
	artist.SetClientInfo(c.GetClientID(), c.GetClientType(), item.RatingKey)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("artistID", id).
		Str("artistTitle", artist.Data.Details.Title).
		Msg("Successfully converted music artist data")

	return artist, nil
}

func (c *PlexClient) GetMusicAlbumByID(ctx context.Context, id string) (*models.MediaItem[*types.Album], error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("albumID", id).
		Msg("Retrieving specific music album from Plex server")

	ratingKey, _ := strconv.Atoi(id)
	int64RatingKey := int64(ratingKey)

	log.Debug().
		Str("albumID", id).
		Int64("ratingKey", int64RatingKey).
		Msg("Making API request to Plex server for music album")

	res, err := c.plexAPI.Library.GetMediaMetaData(ctx, operations.GetMediaMetaDataRequest{RatingKey: int64RatingKey})
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Str("albumID", id).
			Msg("Failed to get music album from Plex")
		return nil, fmt.Errorf("failed to get music album: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil || len(res.Object.MediaContainer.Metadata) == 0 {
		log.Error().
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Str("albumID", id).
			Msg("Music album not found in Plex")
		return nil, fmt.Errorf("music album not found")
	}

	item := res.Object.MediaContainer.Metadata[0]
	if item.Type != "album" {
		log.Error().
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Str("albumID", id).
			Str("actualType", item.Type).
			Msg("Item retrieved is not a music album")
		return nil, fmt.Errorf("item is not a music album")
	}

	// Get artist info if available
	var artistID string
	var artistName string

	if item.ParentRatingKey != nil {
		artistID = *item.ParentRatingKey
		artistKey, _ := strconv.Atoi(artistID)
		int64ArtistKey := int64(artistKey)

		log.Debug().
			Str("artistID", artistID).
			Msg("Getting parent artist information")

		artistRes, err := c.plexAPI.Library.GetMediaMetaData(ctx, operations.GetMediaMetaDataRequest{
			RatingKey: int64ArtistKey,
		})

		if err == nil && artistRes.Object.MediaContainer != nil &&
			artistRes.Object.MediaContainer.Metadata != nil &&
			len(artistRes.Object.MediaContainer.Metadata) > 0 {

			artistName = artistRes.Object.MediaContainer.Metadata[0].Title
			log.Debug().
				Str("artistName", artistName).
				Msg("Retrieved artist name")
		}
	}

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("albumID", id).
		Str("albumTitle", item.Title).
		Str("artistName", artistName).
		Msg("Successfully retrieved music album from Plex")

	itemAlbum, err := GetItemFromMetadata[*types.Album](ctx, c, &item)
	album, err := GetMediaItem[*types.Album](ctx, c, itemAlbum, item.RatingKey)

	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Str("albumID", id).
			Msg("Error converting Plex item to album format")
		return nil, fmt.Errorf("error converting album data: %w", err)
	}

	// Set artist information
	album.Data.ArtistName = artistName
	// album.Data.AddSyncClient(c.GetClientID(), artistID, "")
	album.SetClientInfo(c.GetClientID(), c.GetClientType(), item.RatingKey)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Str("albumID", id).
		Str("albumTitle", album.Data.Details.Title).
		Str("artistName", album.Data.ArtistName).
		Msg("Successfully converted music album data")

	return album, nil
}

// GetMusicGenres retrieves music genres from Plex
func (c *PlexClient) GetMusicGenres(ctx context.Context) ([]string, error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Msg("Retrieving music genres from Plex server")

	// Find the music library section
	log.Debug().Msg("Finding music library section")
	musicSectionKey, err := c.findLibrarySectionByType(ctx, "artist")
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Msg("Failed to find music library section")
		return nil, err
	}

	if musicSectionKey == "" {
		log.Info().
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Msg("No music library section found in Plex")
		return []string{}, nil
	}

	// Get genres from the library items
	sectionKey, _ := strconv.Atoi(musicSectionKey)

	log.Debug().
		Int("sectionKey", sectionKey).
		Msg("Making API request to Plex server for music content")

	res, err := c.plexAPI.Library.GetLibraryItems(ctx, operations.GetLibraryItemsRequest{
		IncludeMeta: operations.GetLibraryItemsQueryParamIncludeMetaEnable.ToPointer(),
		Tag:         "all",
		Type:        operations.GetLibraryItemsQueryParamTypeAudio,
		SectionKey:  sectionKey})

	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.GetClientID()).
			Str("clientType", string(c.GetClientType())).
			Int("sectionKey", sectionKey).
			Msg("Failed to get music content from Plex")
		return nil, fmt.Errorf("failed to get music content: %w", err)
	}

	genreMap := make(map[string]bool)
	if res.Object.MediaContainer != nil && res.Object.MediaContainer.Metadata != nil {
		log.Debug().
			Int("contentCount", len(res.Object.MediaContainer.Metadata)).
			Msg("Extracting genres from music content")

		for _, item := range res.Object.MediaContainer.Metadata {
			if item.Genre != nil {
				for _, genre := range item.Genre {
					if genre.Tag != nil {
						genreMap[*genre.Tag] = true
						log.Debug().
							Str("genre", *genre.Tag).
							Msg("Found music genre")
					}
				}
			}
		}
	}

	genres := make([]string, 0, len(genreMap))
	for genre := range genreMap {
		genres = append(genres, genre)
	}

	log.Info().
		Uint64("clientID", c.GetClientID()).
		Str("clientType", string(c.GetClientType())).
		Int("genresFound", len(genres)).
		Msg("Successfully retrieved music genres from Plex")

	return genres, nil
}
