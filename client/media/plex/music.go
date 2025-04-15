package plex

import (
	"context"
	"fmt"
	"strconv"
	"suasor/client/media/types"
	"suasor/types/models"
	"suasor/utils"

	"github.com/LukeHagar/plexgo"
	"github.com/LukeHagar/plexgo/models/operations"
)

// GetMusic retrieves music tracks from Plex
func (c *PlexClient) GetMusic(ctx context.Context, options *types.QueryOptions) ([]models.MediaItem[*types.Track], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Msg("Retrieving music tracks from Plex server")

	// Find the music library section
	log.Debug().Msg("Finding music library section")
	musicSectionKey, err := c.findLibrarySectionByType(ctx, "artist")
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("Failed to find music library section")
		return nil, err
	}

	if musicSectionKey == "" {
		log.Info().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("No music library section found in Plex")
		return []models.MediaItem[*types.Track]{}, nil
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
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Int("sectionKey", sectionKey).
			Msg("Failed to get music artists from Plex")
		return nil, fmt.Errorf("failed to get music artists: %w", err)
	}

	var tracks []models.MediaItem[*types.Track]

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

						for _, item := range tracksRes.Object.MediaContainer.Metadata {
							if *item.Type != "track" {
								continue
							}

							track := models.MediaItem[*types.Track]{
								Data: &types.Track{
									Details:    c.createChildMetadataFromPlexItem(&item),
									Number:     *item.Index,
									ArtistName: artist.Title,
									AlbumName:  *album.Title,
								},
							}
							track.Data.AddSyncClient(c.ClientID, *item.RatingKey, *item.ParentRatingKey)
							track.SetClientInfo(c.ClientID, c.ClientType, *item.RatingKey)
							tracks = append(tracks, track)

							log.Debug().
								Str("trackID", *item.RatingKey).
								Str("trackTitle", *item.Title).
								Int("trackNumber", *item.Index).
								Msg("Added track to result list")

							// Limit number of tracks to avoid too large responses
							if len(tracks) >= 100 {
								log.Info().
									Int("trackCount", len(tracks)).
									Msg("Reached track limit (100), returning results")
								return tracks, nil
							}
						}
					}
				}
			}
		}
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("tracksReturned", len(tracks)).
		Msg("Completed GetMusic request")

	return tracks, nil
}

// GetMusicArtists retrieves music artists from Plex
func (c *PlexClient) GetMusicArtists(ctx context.Context, options *types.QueryOptions) ([]models.MediaItem[*types.Artist], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Msg("Retrieving music artists from Plex server")

	// Find the music library section
	log.Debug().Msg("Finding music library section")
	musicSectionKey, err := c.findLibrarySectionByType(ctx, "artist")
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("Failed to find music library section")
		return nil, err
	}

	if musicSectionKey == "" {
		log.Info().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("No music library section found in Plex")
		return []models.MediaItem[*types.Artist]{}, nil
	}

	sectionKey, _ := strconv.Atoi(musicSectionKey)

	log.Debug().
		Int("sectionKey", sectionKey).
		Msg("Making API request to Plex server for music artists")

	res, err := c.plexAPI.Library.GetLibraryItems(ctx, operations.GetLibraryItemsRequest{SectionKey: sectionKey})
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Int("sectionKey", sectionKey).
			Msg("Failed to get music artists from Plex")
		return nil, fmt.Errorf("failed to get music artists: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil {
		log.Info().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("No music artists found in Plex")
		return []models.MediaItem[*types.Artist]{}, nil
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("totalItems", len(res.Object.MediaContainer.Metadata)).
		Msg("Successfully retrieved music artists from Plex")

	artists := make([]models.MediaItem[*types.Artist], 0, len(res.Object.MediaContainer.Metadata))
	for _, item := range res.Object.MediaContainer.Metadata {
		if item.Type != "artist" {
			continue
		}

		artist := models.MediaItem[*types.Artist]{
			Data: &types.Artist{
				Details: c.createMetadataFromPlexItem(&item),
			},
		}

		artist.SetClientInfo(c.ClientID, c.ClientType, item.RatingKey)

		artists = append(artists, artist)

		log.Debug().
			Str("artistID", item.RatingKey).
			Str("artistName", item.Title).
			Msg("Added artist to result list")
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("artistsReturned", len(artists)).
		Msg("Completed GetMusicArtists request")

	return artists, nil
}

// GetMusicAlbums retrieves music albums from Plex
func (c *PlexClient) GetMusicAlbums(ctx context.Context, options *types.QueryOptions) ([]models.MediaItem[*types.Album], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Msg("Retrieving music albums from Plex server")

	// Find the music library section
	log.Debug().Msg("Finding music library section")
	musicSectionKey, err := c.findLibrarySectionByType(ctx, "artist")
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("Failed to find music library section")
		return nil, err
	}

	if musicSectionKey == "" {
		log.Info().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("No music library section found in Plex")
		return []models.MediaItem[*types.Album]{}, nil
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
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Int("sectionKey", sectionKey).
			Msg("Failed to get music artists from Plex")
		return nil, fmt.Errorf("failed to get music artists: %w", err)
	}

	var albums []models.MediaItem[*types.Album]

	// Loop through artists to get their albums
	if res.Object.MediaContainer != nil && res.Object.MediaContainer.Metadata != nil {
		log.Debug().
			Int("artistCount", len(res.Object.MediaContainer.Metadata)).
			Msg("Processing artists to find albums")

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

			if albumsRes.Object.MediaContainer != nil && albumsRes.Object.MediaContainer.Metadata != nil {
				log.Debug().
					Str("artistID", artist.RatingKey).
					Str("artistName", artist.Title).
					Int("albumCount", len(albumsRes.Object.MediaContainer.Metadata)).
					Msg("Processing albums for artist")

				for _, item := range albumsRes.Object.MediaContainer.Metadata {
					if *item.Type != "album" {
						continue
					}

					album := models.MediaItem[*types.Album]{
						Data: &types.Album{
							Details:    c.createChildMetadataFromPlexItem(&item),
							ArtistName: artist.Title,
							TrackCount: *item.LeafCount,
						},
					}

					album.Data.AddSyncClient(c.ClientID, *item.RatingKey)
					album.SetClientInfo(c.ClientID, c.ClientType, *item.RatingKey)
					albums = append(albums, album)

					log.Debug().
						Str("albumID", *item.RatingKey).
						Str("albumName", *item.Title).
						Int("trackCount", *item.LeafCount).
						Msg("Added album to result list")
				}
			}
		}
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("albumsReturned", len(albums)).
		Msg("Completed GetMusicAlbums request")

	return albums, nil
}

// GetMusicTrackByID retrieves a specific music track by ID
func (c *PlexClient) GetMusicTrackByID(ctx context.Context, id string) (models.MediaItem[*types.Track], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
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
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("trackID", id).
			Msg("Failed to get music track from Plex")
		return models.MediaItem[*types.Track]{}, fmt.Errorf("failed to get music track: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil || len(res.Object.MediaContainer.Metadata) == 0 {
		log.Error().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("trackID", id).
			Msg("Music track not found in Plex")
		return models.MediaItem[*types.Track]{}, fmt.Errorf("music track not found")
	}

	item := res.Object.MediaContainer.Metadata[0]
	if item.Type != "track" {
		log.Error().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("trackID", id).
			Str("actualType", item.Type).
			Msg("Item retrieved is not a music track")
		return models.MediaItem[*types.Track]{}, fmt.Errorf("item is not a music track")
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
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("trackID", id).
		Str("trackTitle", item.Title).
		Str("albumName", albumName).
		Str("artistName", artistName).
		Msg("Successfully retrieved music track from Plex")

	track := models.MediaItem[*types.Track]{
		Data: &types.Track{
			AlbumName:  albumName,
			ArtistName: artistName,
			Number:     int(*item.Index),
			Details:    c.createMediaDetailsFromPlexItem(&item),
		},
	}

	track.Data.AddSyncClient(c.ClientID, *item.ParentRatingKey, artistID)
	track.SetClientInfo(c.ClientID, c.ClientType, item.RatingKey)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("trackID", id).
		Str("trackTitle", track.Data.Details.Title).
		Int("trackNumber", track.Data.Number).
		Msg("Successfully converted music track data")

	return track, nil
}

// GetMusicGenres retrieves music genres from Plex
func (c *PlexClient) GetMusicGenres(ctx context.Context) ([]string, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Msg("Retrieving music genres from Plex server")

	// Find the music library section
	log.Debug().Msg("Finding music library section")
	musicSectionKey, err := c.findLibrarySectionByType(ctx, "artist")
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("Failed to find music library section")
		return nil, err
	}

	if musicSectionKey == "" {
		log.Info().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
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
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
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
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("genresFound", len(genres)).
		Msg("Successfully retrieved music genres from Plex")

	return genres, nil
}
