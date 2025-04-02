// services/media_client_music.go
package services

import (
	"context"
	"fmt"
	"suasor/client"
	"suasor/client/media"
	"suasor/client/media/providers"
	mediatypes "suasor/client/media/types"
	"suasor/client/types"
	"suasor/repository"
	models "suasor/types/models"
	"suasor/utils"
)

// MusicSearchResults contains search results for music
type MusicSearchResults struct {
	Artists []models.MediaItem[mediatypes.Artist] `json:"artists"`
	Albums  []models.MediaItem[mediatypes.Album]  `json:"albums"`
	Tracks  []models.MediaItem[mediatypes.Track]  `json:"tracks"`
}

// MediaClientMusicService defines the music client service interface
type MediaClientMusicService[T types.MediaClientConfig] interface {
	GetTrackByID(ctx context.Context, userID, clientID uint64, trackID string) (models.MediaItem[mediatypes.Track], error)
	GetAlbumByID(ctx context.Context, userID, clientID uint64, albumID string) (models.MediaItem[mediatypes.Album], error)
	GetArtistByID(ctx context.Context, userID, clientID uint64, artistID string) (models.MediaItem[mediatypes.Artist], error)
	GetTracksByAlbum(ctx context.Context, userID, clientID uint64, albumID string) ([]models.MediaItem[mediatypes.Track], error)
	GetAlbumsByArtist(ctx context.Context, userID, clientID uint64, artistID string) ([]models.MediaItem[mediatypes.Album], error)
	GetArtistsByGenre(ctx context.Context, userID uint64, genre string) ([]models.MediaItem[mediatypes.Artist], error)
	GetAlbumsByGenre(ctx context.Context, userID uint64, genre string) ([]models.MediaItem[mediatypes.Album], error)
	GetTracksByGenre(ctx context.Context, userID uint64, genre string) ([]models.MediaItem[mediatypes.Track], error)
	GetAlbumsByYear(ctx context.Context, userID uint64, year int) ([]models.MediaItem[mediatypes.Album], error)
	GetLatestAlbumsByAdded(ctx context.Context, userID uint64, count int) ([]models.MediaItem[mediatypes.Album], error)
	GetPopularAlbums(ctx context.Context, userID uint64, count int) ([]models.MediaItem[mediatypes.Album], error)
	GetPopularArtists(ctx context.Context, userID uint64, count int) ([]models.MediaItem[mediatypes.Artist], error)
	SearchMusic(ctx context.Context, userID uint64, query string) (MusicSearchResults, error)
}

// MediaClientMusicServiceImpl implements the MediaClientMusicService interface
type MediaClientMusicServiceImpl[T types.MediaClientConfig] struct {
	clientRepo repository.ClientRepository[T]
	factory    *client.ClientFactoryService
}

// NewMediaClientMusicService creates a new media client music service
func NewMediaClientMusicService[T types.MediaClientConfig](
	clientRepo repository.ClientRepository[T],
	factory *client.ClientFactoryService,
) MediaClientMusicService[T] {
	return &MediaClientMusicServiceImpl[T]{
		clientRepo: clientRepo,
		factory:    factory,
	}
}

// getMusicClients gets all music clients for a user
func (s *MediaClientMusicServiceImpl[T]) getMusicClients(ctx context.Context, userID uint64) ([]media.MediaClient, error) {
	repo := s.clientRepo
	// Get all media clients for the user
	clients, err := repo.GetByCategory(ctx, types.ClientCategoryMedia, userID)
	if err != nil {
		return nil, err
	}

	var musicClients []media.MediaClient

	// Filter and instantiate clients that support music
	for _, clientConfig := range clients {
		if clientConfig.Config.Data.SupportsMusic() {
			clientId := clientConfig.GetID()
			client, err := s.factory.GetClient(ctx, clientId, clientConfig.Config.Data.GetType())
			if err != nil {
				// Log error but continue with other clients
				continue
			}
			musicClients = append(musicClients, client.(media.MediaClient))
		}
	}

	return musicClients, nil
}

// getSpecificMusicClient gets a specific music client
func (s *MediaClientMusicServiceImpl[T]) getSpecificMusicClient(ctx context.Context, userID, clientID uint64) (media.MediaClient, error) {
	log := utils.LoggerFromContext(ctx)

	clientConfig, err := (s.clientRepo).GetByID(ctx, clientID, userID)
	if err != nil {
		return nil, err
	}
	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientType", clientConfig.Config.Data.GetType().String()).
		Msg("Retrieved client config")

	if !clientConfig.Config.Data.SupportsMusic() {
		log.Warn().
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("clientType", clientConfig.Config.Data.GetType().String()).
			Msg("Client does not support music")
		return nil, fmt.Errorf("client %d does not support music", clientID)
	}

	client, err := s.factory.GetClient(ctx, clientID, clientConfig.Config.Data.GetType())
	if err != nil {
		log.Error().
			Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("clientType", clientConfig.Config.Data.GetType().String()).
			Msg("Failed to get client")
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	return client.(media.MediaClient), nil
}

// GetTrackByID retrieves a specific music track by ID
func (s *MediaClientMusicServiceImpl[T]) GetTrackByID(ctx context.Context, userID, clientID uint64, trackID string) (models.MediaItem[mediatypes.Track], error) {
	log := utils.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("trackID", trackID).
		Msg("Getting track by ID")

	client, err := s.getSpecificMusicClient(ctx, userID, clientID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Msg("Failed to get client")
		return models.MediaItem[mediatypes.Track]{}, fmt.Errorf("failed to get client: %w", err)
	}

	musicClient, ok := client.(providers.MusicProvider)
	if !ok {
		log.Error().
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("clientType", "media client").
			Msg("Client does not implement music provider interface")
		return models.MediaItem[mediatypes.Track]{}, fmt.Errorf("client does not implement music provider interface")
	}

	return musicClient.GetMusicTrackByID(ctx, trackID)
}

// GetAlbumByID retrieves a specific music album by ID
func (s *MediaClientMusicServiceImpl[T]) GetAlbumByID(ctx context.Context, userID, clientID uint64, albumID string) (models.MediaItem[mediatypes.Album], error) {
	log := utils.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("albumID", albumID).
		Msg("Getting album by ID")

	client, err := s.getSpecificMusicClient(ctx, userID, clientID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Msg("Failed to get client")
		return models.MediaItem[mediatypes.Album]{}, fmt.Errorf("failed to get client: %w", err)
	}

	musicClient, ok := client.(providers.MusicProvider)
	if !ok {
		log.Error().
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("clientType", "media client").
			Msg("Client does not implement music provider interface")
		return models.MediaItem[mediatypes.Album]{}, fmt.Errorf("client does not implement music provider interface")
	}

	options := &mediatypes.QueryOptions{
		Filters: map[string]string{
			"id": albumID,
		},
	}

	albums, err := musicClient.GetMusicAlbums(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("albumID", albumID).
			Msg("Failed to get album")
		return models.MediaItem[mediatypes.Album]{}, fmt.Errorf("failed to get album: %w", err)
	}

	if len(albums) == 0 {
		log.Error().
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("albumID", albumID).
			Msg("Album not found")
		return models.MediaItem[mediatypes.Album]{}, fmt.Errorf("album not found")
	}

	return albums[0], nil
}

// GetArtistByID retrieves a specific music artist by ID
func (s *MediaClientMusicServiceImpl[T]) GetArtistByID(ctx context.Context, userID, clientID uint64, artistID string) (models.MediaItem[mediatypes.Artist], error) {
	log := utils.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("artistID", artistID).
		Msg("Getting artist by ID")

	client, err := s.getSpecificMusicClient(ctx, userID, clientID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Msg("Failed to get client")
		return models.MediaItem[mediatypes.Artist]{}, fmt.Errorf("failed to get client: %w", err)
	}

	musicClient, ok := client.(providers.MusicProvider)
	if !ok {
		log.Error().
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("clientType", "media client").
			Msg("Client does not implement music provider interface")
		return models.MediaItem[mediatypes.Artist]{}, fmt.Errorf("client does not implement music provider interface")
	}

	options := &mediatypes.QueryOptions{
		Filters: map[string]string{
			"id": artistID,
		},
	}

	artists, err := musicClient.GetMusicArtists(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("artistID", artistID).
			Msg("Failed to get artist")
		return models.MediaItem[mediatypes.Artist]{}, fmt.Errorf("failed to get artist: %w", err)
	}

	if len(artists) == 0 {
		log.Error().
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("artistID", artistID).
			Msg("Artist not found")
		return models.MediaItem[mediatypes.Artist]{}, fmt.Errorf("artist not found")
	}

	return artists[0], nil
}

// GetTracksByAlbum retrieves all tracks for a specific album
func (s *MediaClientMusicServiceImpl[T]) GetTracksByAlbum(ctx context.Context, userID, clientID uint64, albumID string) ([]models.MediaItem[mediatypes.Track], error) {
	log := utils.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("albumID", albumID).
		Msg("Getting tracks by album")

	client, err := s.getSpecificMusicClient(ctx, userID, clientID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Msg("Failed to get client")
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	musicClient, ok := client.(providers.MusicProvider)
	if !ok {
		log.Error().
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("clientType", "media client").
			Msg("Client does not implement music provider interface")
		return nil, fmt.Errorf("client does not implement music provider interface")
	}

	options := &mediatypes.QueryOptions{
		Filters: map[string]string{
			"albumID": albumID,
		},
	}

	tracks, err := musicClient.GetMusic(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("albumID", albumID).
			Msg("Failed to get tracks by album")
		return nil, fmt.Errorf("failed to get tracks by album: %w", err)
	}

	return tracks, nil
}

// GetAlbumsByArtist retrieves all albums for a specific artist
func (s *MediaClientMusicServiceImpl[T]) GetAlbumsByArtist(ctx context.Context, userID, clientID uint64, artistID string) ([]models.MediaItem[mediatypes.Album], error) {
	log := utils.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("artistID", artistID).
		Msg("Getting albums by artist")

	client, err := s.getSpecificMusicClient(ctx, userID, clientID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Msg("Failed to get client")
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	musicClient, ok := client.(providers.MusicProvider)
	if !ok {
		log.Error().
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("clientType", "media client").
			Msg("Client does not implement music provider interface")
		return nil, fmt.Errorf("client does not implement music provider interface")
	}

	options := &mediatypes.QueryOptions{
		Filters: map[string]string{
			"artistID": artistID,
		},
	}

	albums, err := musicClient.GetMusicAlbums(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("artistID", artistID).
			Msg("Failed to get albums by artist")
		return nil, fmt.Errorf("failed to get albums by artist: %w", err)
	}

	return albums, nil
}

// GetArtistsByGenre retrieves all artists for a specific genre
func (s *MediaClientMusicServiceImpl[T]) GetArtistsByGenre(ctx context.Context, userID uint64, genre string) ([]models.MediaItem[mediatypes.Artist], error) {
	log := utils.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", userID).
		Str("genre", genre).
		Msg("Getting artists by genre")

	musicClients, err := s.getMusicClients(ctx, userID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to get music clients")
		return nil, fmt.Errorf("failed to get music clients: %w", err)
	}

	var allArtists []models.MediaItem[mediatypes.Artist]

	for _, client := range musicClients {
		musicClient, ok := client.(providers.MusicProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Filters: map[string]string{
				"genre": genre,
			},
		}

		artists, err := musicClient.GetMusicArtists(ctx, options)
		if err != nil {
			log.Warn().Err(err).
				Uint64("userID", userID).
				Uint64("clientID", client.GetClientID()).
				Str("clientType", "media client").
				Str("genre", genre).
				Msg("Failed to get artists by genre from client, continuing with others")
			continue
		}

		allArtists = append(allArtists, artists...)
	}

	return allArtists, nil
}

// GetAlbumsByGenre retrieves all albums for a specific genre
func (s *MediaClientMusicServiceImpl[T]) GetAlbumsByGenre(ctx context.Context, userID uint64, genre string) ([]models.MediaItem[mediatypes.Album], error) {
	log := utils.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", userID).
		Str("genre", genre).
		Msg("Getting albums by genre")

	musicClients, err := s.getMusicClients(ctx, userID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to get music clients")
		return nil, fmt.Errorf("failed to get music clients: %w", err)
	}

	var allAlbums []models.MediaItem[mediatypes.Album]

	for _, client := range musicClients {
		musicClient, ok := client.(providers.MusicProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Filters: map[string]string{
				"genre": genre,
			},
		}

		albums, err := musicClient.GetMusicAlbums(ctx, options)
		if err != nil {
			log.Warn().Err(err).
				Uint64("userID", userID).
				Uint64("clientID", client.GetClientID()).
				Str("clientType", "media client").
				Str("genre", genre).
				Msg("Failed to get albums by genre from client, continuing with others")
			continue
		}

		allAlbums = append(allAlbums, albums...)
	}

	return allAlbums, nil
}

// GetTracksByGenre retrieves all tracks for a specific genre
func (s *MediaClientMusicServiceImpl[T]) GetTracksByGenre(ctx context.Context, userID uint64, genre string) ([]models.MediaItem[mediatypes.Track], error) {
	log := utils.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", userID).
		Str("genre", genre).
		Msg("Getting tracks by genre")

	musicClients, err := s.getMusicClients(ctx, userID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to get music clients")
		return nil, fmt.Errorf("failed to get music clients: %w", err)
	}

	var allTracks []models.MediaItem[mediatypes.Track]

	for _, client := range musicClients {
		musicClient, ok := client.(providers.MusicProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Filters: map[string]string{
				"genre": genre,
			},
		}

		tracks, err := musicClient.GetMusic(ctx, options)
		if err != nil {
			log.Warn().Err(err).
				Uint64("userID", userID).
				Uint64("clientID", client.GetClientID()).
				Str("clientType", "media client").
				Str("genre", genre).
				Msg("Failed to get tracks by genre from client, continuing with others")
			continue
		}

		allTracks = append(allTracks, tracks...)
	}

	return allTracks, nil
}

// GetAlbumsByYear retrieves all albums for a specific release year
func (s *MediaClientMusicServiceImpl[T]) GetAlbumsByYear(ctx context.Context, userID uint64, year int) ([]models.MediaItem[mediatypes.Album], error) {
	log := utils.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", userID).
		Int("year", year).
		Msg("Getting albums by year")

	musicClients, err := s.getMusicClients(ctx, userID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to get music clients")
		return nil, fmt.Errorf("failed to get music clients: %w", err)
	}

	var allAlbums []models.MediaItem[mediatypes.Album]

	for _, client := range musicClients {
		musicClient, ok := client.(providers.MusicProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Filters: map[string]string{
				"releaseYear": fmt.Sprintf("%d", year),
			},
		}

		albums, err := musicClient.GetMusicAlbums(ctx, options)
		if err != nil {
			log.Warn().Err(err).
				Uint64("userID", userID).
				Uint64("clientID", client.GetClientID()).
				Str("clientType", "media client").
				Int("year", year).
				Msg("Failed to get albums by year from client, continuing with others")
			continue
		}

		allAlbums = append(allAlbums, albums...)
	}

	return allAlbums, nil
}

// GetLatestAlbumsByAdded retrieves the most recently added albums
func (s *MediaClientMusicServiceImpl[T]) GetLatestAlbumsByAdded(ctx context.Context, userID uint64, count int) ([]models.MediaItem[mediatypes.Album], error) {
	log := utils.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", userID).
		Int("count", count).
		Msg("Getting latest albums by added date")

	musicClients, err := s.getMusicClients(ctx, userID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to get music clients")
		return nil, fmt.Errorf("failed to get music clients: %w", err)
	}

	var allAlbums []models.MediaItem[mediatypes.Album]

	for _, client := range musicClients {
		musicClient, ok := client.(providers.MusicProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Limit:     count,
			Sort:      "addedAt",
			SortOrder: mediatypes.SortOrderDesc,
		}

		albums, err := musicClient.GetMusicAlbums(ctx, options)
		if err != nil {
			log.Warn().Err(err).
				Uint64("userID", userID).
				Uint64("clientID", client.GetClientID()).
				Str("clientType", "media client").
				Int("count", count).
				Msg("Failed to get latest albums from client, continuing with others")
			continue
		}

		allAlbums = append(allAlbums, albums...)
	}

	// We should handle limiting and sorting across all clients
	// TODO: Implement proper sorting by addedAt across all clients

	// Simple limit
	if len(allAlbums) > count {
		allAlbums = allAlbums[:count]
	}

	return allAlbums, nil
}

// GetPopularAlbums retrieves the most popular albums
func (s *MediaClientMusicServiceImpl[T]) GetPopularAlbums(ctx context.Context, userID uint64, count int) ([]models.MediaItem[mediatypes.Album], error) {
	log := utils.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", userID).
		Int("count", count).
		Msg("Getting popular albums")

	musicClients, err := s.getMusicClients(ctx, userID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to get music clients")
		return nil, fmt.Errorf("failed to get music clients: %w", err)
	}

	var allAlbums []models.MediaItem[mediatypes.Album]

	for _, client := range musicClients {
		musicClient, ok := client.(providers.MusicProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Limit:     count,
			Sort:      "popularity", // This might vary based on client implementation
			SortOrder: mediatypes.SortOrderDesc,
		}

		albums, err := musicClient.GetMusicAlbums(ctx, options)
		if err != nil {
			log.Warn().Err(err).
				Uint64("userID", userID).
				Uint64("clientID", client.GetClientID()).
				Str("clientType", "media client").
				Int("count", count).
				Msg("Failed to get popular albums from client, continuing with others")
			continue
		}

		allAlbums = append(allAlbums, albums...)
	}

	// We should handle limiting and sorting across all clients
	// TODO: Implement proper sorting by popularity across all clients

	// Simple limit
	if len(allAlbums) > count {
		allAlbums = allAlbums[:count]
	}

	return allAlbums, nil
}

// GetPopularArtists retrieves the most popular artists
func (s *MediaClientMusicServiceImpl[T]) GetPopularArtists(ctx context.Context, userID uint64, count int) ([]models.MediaItem[mediatypes.Artist], error) {
	log := utils.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", userID).
		Int("count", count).
		Msg("Getting popular artists")

	musicClients, err := s.getMusicClients(ctx, userID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to get music clients")
		return nil, fmt.Errorf("failed to get music clients: %w", err)
	}

	var allArtists []models.MediaItem[mediatypes.Artist]

	for _, client := range musicClients {
		musicClient, ok := client.(providers.MusicProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Limit:     count,
			Sort:      "popularity", // This might vary based on client implementation
			SortOrder: mediatypes.SortOrderDesc,
		}

		artists, err := musicClient.GetMusicArtists(ctx, options)
		if err != nil {
			log.Warn().Err(err).
				Uint64("userID", userID).
				Uint64("clientID", client.GetClientID()).
				Str("clientType", "media client").
				Int("count", count).
				Msg("Failed to get popular artists from client, continuing with others")
			continue
		}

		allArtists = append(allArtists, artists...)
	}

	// We should handle limiting and sorting across all clients
	// TODO: Implement proper sorting by popularity across all clients

	// Simple limit
	if len(allArtists) > count {
		allArtists = allArtists[:count]
	}

	return allArtists, nil
}

// SearchMusic searches for music (artists, albums, tracks) across all clients
func (s *MediaClientMusicServiceImpl[T]) SearchMusic(ctx context.Context, userID uint64, query string) (MusicSearchResults, error) {
	log := utils.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", userID).
		Str("query", query).
		Msg("Searching music")

	musicClients, err := s.getMusicClients(ctx, userID)
	if err != nil {
		log.Error().Err(err).
			Uint64("userID", userID).
			Msg("Failed to get music clients")
		return MusicSearchResults{}, fmt.Errorf("failed to get music clients: %w", err)
	}

	var results MusicSearchResults

	for _, client := range musicClients {
		musicClient, ok := client.(providers.MusicProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Query: query,
		}

		// Get artists
		artists, err := musicClient.GetMusicArtists(ctx, options)
		if err != nil {
			log.Warn().Err(err).
				Uint64("userID", userID).
				Uint64("clientID", client.GetClientID()).
				Str("clientType", "media client").
				Str("query", query).
				Msg("Failed to search artists from client, continuing with others")
		} else {
			results.Artists = append(results.Artists, artists...)
		}

		// Get albums
		albums, err := musicClient.GetMusicAlbums(ctx, options)
		if err != nil {
			log.Warn().Err(err).
				Uint64("userID", userID).
				Uint64("clientID", client.GetClientID()).
				Str("clientType", "media client").
				Str("query", query).
				Msg("Failed to search albums from client, continuing with others")
		} else {
			results.Albums = append(results.Albums, albums...)
		}

		// Get tracks
		tracks, err := musicClient.GetMusic(ctx, options)
		if err != nil {
			log.Warn().Err(err).
				Uint64("userID", userID).
				Uint64("clientID", client.GetClientID()).
				Str("clientType", "media client").
				Str("query", query).
				Msg("Failed to search tracks from client, continuing with others")
		} else {
			results.Tracks = append(results.Tracks, tracks...)
		}
	}

	log.Info().
		Uint64("userID", userID).
		Str("query", query).
		Int("artistsCount", len(results.Artists)).
		Int("albumsCount", len(results.Albums)).
		Int("tracksCount", len(results.Tracks)).
		Msg("Music search completed successfully")

	return results, nil
}