package services

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"suasor/client"
	"suasor/client/media"
	"suasor/client/media/providers"
	mediatypes "suasor/client/media/types"
	"suasor/client/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/utils"
)

// ClientMusicService defines operations for interacting with music clients
type ClientMusicService[T types.ClientConfig] interface {
	GetAlbumByID(ctx context.Context, userID uint64, clientID uint64, albumID string) (models.MediaItem[*mediatypes.Album], error)
	GetArtistByID(ctx context.Context, userID uint64, clientID uint64, artistID string) (models.MediaItem[*mediatypes.Artist], error)
	GetTrackByID(ctx context.Context, userID uint64, clientID uint64, trackID string) (models.MediaItem[*mediatypes.Track], error)
	GetRecentlyAddedAlbums(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[*mediatypes.Album], error)
	GetRecentlyPlayedTracks(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[*mediatypes.Track], error)
	GetAlbumsByGenre(ctx context.Context, userID uint64, genre string) ([]*models.MediaItem[*mediatypes.Album], error)
	GetArtistsByGenre(ctx context.Context, userID uint64, genre string) ([]*models.MediaItem[*mediatypes.Artist], error)
	GetRandomAlbums(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[*mediatypes.Album], error)
	GetTopAlbumsForClient(ctx context.Context, userID uint64, clientID uint64, limit int) ([]*models.MediaItem[*mediatypes.Album], error)
	GetTopAlbums(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[*mediatypes.Album], error)
	GetPopularAlbums(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[*mediatypes.Album], error)
	GetPopularArtists(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[*mediatypes.Artist], error)
	GetTopArtists(ctx context.Context, userID uint64, clientID uint64, limit int) ([]*models.MediaItem[*mediatypes.Artist], error)
	GetFavoriteArtists(ctx context.Context, userID uint64, clientID uint64, limit int) ([]*models.MediaItem[*mediatypes.Artist], error)
	GetTracksByAlbum(ctx context.Context, userID uint64, clientID uint64, albumID string) ([]*models.MediaItem[*mediatypes.Track], error)
	GetAlbumsByArtist(ctx context.Context, userID uint64, clientID uint64, artistID string) ([]*models.MediaItem[*mediatypes.Album], error)
	GetTracksByGenre(ctx context.Context, userID uint64, genre string) ([]*models.MediaItem[*mediatypes.Track], error)
	GetAlbumsByYear(ctx context.Context, userID uint64, year int) ([]*models.MediaItem[*mediatypes.Album], error)
	GetLatestAlbumsByAdded(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[*mediatypes.Album], error)
	GetTopTracks(ctx context.Context, userID uint64, clientID uint64, limit int) ([]*models.MediaItem[*mediatypes.Track], error)
	GetRecentlyAddedTracks(ctx context.Context, userID uint64, clientID uint64, limit int) ([]*models.MediaItem[*mediatypes.Track], error)
	SearchMusic(ctx context.Context, userID uint64, query string) (MusicSearchResults, error)
}

// MusicSearchResult is a wrapper for music search results
type MusicSearchResult struct {
	Type     string
	ClientID uint64
	Data     mediatypes.MediaData
}

// MusicSearchResults is a structured container for search results by type
type MusicSearchResults struct {
	Artists []*models.MediaItem[*mediatypes.Artist]
	Albums  []*models.MediaItem[*mediatypes.Album]
	Tracks  []*models.MediaItem[*mediatypes.Track]
}

type mediaMusicService[T types.ClientMediaConfig] struct {
	repo    repository.ClientRepository[T]
	factory *client.ClientFactoryService
}

func NewClientMusicService[T types.ClientMediaConfig](
	repo repository.ClientRepository[T],
	factory *client.ClientFactoryService,
) ClientMusicService[T] {
	return &mediaMusicService[T]{
		repo:    repo,
		factory: factory,
	}
}

// getMusicClients gets all music clients for a user
func (s *mediaMusicService[T]) getMusicClients(ctx context.Context, userID uint64) ([]media.ClientMedia, error) {
	// Get all media clients for the user that support music
	clients, err := s.repo.GetByCategory(ctx, types.ClientCategoryMedia, userID)
	if err != nil {
		return nil, err
	}

	log := utils.LoggerFromContext(ctx)
	log.Debug().Msg("Getting music clients")

	var musicClients []media.ClientMedia

	// Filter for clients that support music and instantiate them
	for _, clientConfig := range clients {
		if clientConfig.Config.Data.SupportsMusic() {
			log.Debug().
				Uint64("clientID", clientConfig.ID).
				Str("clientType", clientConfig.Config.Data.GetType().String()).
				Msg("Found music-supporting client")

			clientId := clientConfig.GetID()
			client, err := s.factory.GetClient(ctx, clientId, clientConfig.Config.Data)
			if err != nil {
				log.Error().
					Err(err).
					Uint64("clientID", clientId).
					Msg("Failed to instantiate music client")
				continue
			}

			musicClients = append(musicClients, client.(media.ClientMedia))
		}
	}

	log.Debug().
		Int("musicClientCount", len(musicClients)).
		Msg("Retrieved music clients")

	return musicClients, nil
}

// getSpecificMusicClient gets a specific music client
func (s *mediaMusicService[T]) getSpecificMusicClient(ctx context.Context, userID, clientID uint64) (media.ClientMedia, error) {
	log := utils.LoggerFromContext(ctx)

	clientConfig, err := s.repo.GetByID(ctx, clientID)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", clientID).
			Msg("Failed to get client config")
		return nil, err
	}

	if !clientConfig.Config.Data.SupportsMusic() {
		log.Warn().
			Uint64("clientID", clientID).
			Str("clientType", clientConfig.Config.Data.GetType().String()).
			Msg("Client does not support music")
		return nil, errors.New("client does not support music")
	}

	client, err := s.factory.GetClient(ctx, clientID, clientConfig.Config.Data)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", clientID).
			Str("clientType", clientConfig.Config.Data.GetType().String()).
			Msg("Failed to instantiate music client")
		return nil, err
	}

	return client.(media.ClientMedia), nil
}

func (s *mediaMusicService[T]) GetAlbumByID(ctx context.Context, userID uint64, clientID uint64, albumID string) (models.MediaItem[*mediatypes.Album], error) {
	log := utils.LoggerFromContext(ctx)

	// Get the specified client
	client, err := s.getSpecificMusicClient(ctx, userID, clientID)
	if err != nil {
		return models.MediaItem[*mediatypes.Album]{}, err
	}

	// Check if client implements MusicProvider
	musicClient, ok := client.(providers.MusicProvider)
	if !ok {
		log.Warn().
			Uint64("clientID", clientID).
			Str("clientType", "media client").
			Msg("Client does not implement music provider interface")
		return models.MediaItem[*mediatypes.Album]{}, fmt.Errorf("client does not implement music provider interface")
	}

	options := &mediatypes.QueryOptions{
		ExternalSourceID: albumID,
	}

	albums, err := musicClient.GetMusicAlbums(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Uint64("clientID", clientID).
			Str("albumID", albumID).
			Msg("Failed to retrieve album")
		return models.MediaItem[*mediatypes.Album]{}, err
	}

	if len(albums) == 0 {
		log.Warn().
			Uint64("clientID", clientID).
			Str("albumID", albumID).
			Msg("Album not found")
		return models.MediaItem[*mediatypes.Album]{}, fmt.Errorf("album not found")
	}

	return *albums[0], nil
}

func (s *mediaMusicService[T]) GetArtistByID(ctx context.Context, userID uint64, clientID uint64, artistID string) (models.MediaItem[*mediatypes.Artist], error) {
	log := utils.LoggerFromContext(ctx)

	// Get the specified client
	client, err := s.getSpecificMusicClient(ctx, userID, clientID)
	if err != nil {
		return models.MediaItem[*mediatypes.Artist]{}, err
	}

	// Check if client implements MusicProvider
	musicClient, ok := client.(providers.MusicProvider)
	if !ok {
		log.Warn().
			Uint64("clientID", clientID).
			Str("clientType", "media client").
			Msg("Client does not implement music provider interface")
		return models.MediaItem[*mediatypes.Artist]{}, fmt.Errorf("client does not implement music provider interface")
	}

	options := &mediatypes.QueryOptions{
		ExternalSourceID: artistID,
	}

	artists, err := musicClient.GetMusicArtists(ctx, options)
	if err != nil {
		log.Error().Err(err).
			Uint64("clientID", clientID).
			Str("artistID", artistID).
			Msg("Failed to retrieve artist")
		return models.MediaItem[*mediatypes.Artist]{}, err
	}

	if len(artists) == 0 {
		log.Warn().
			Uint64("clientID", clientID).
			Str("artistID", artistID).
			Msg("Artist not found")
		return models.MediaItem[*mediatypes.Artist]{}, fmt.Errorf("artist not found")
	}

	return *artists[0], nil
}

func (s *mediaMusicService[T]) GetTrackByID(ctx context.Context, userID uint64, clientID uint64, trackID string) (models.MediaItem[*mediatypes.Track], error) {
	log := utils.LoggerFromContext(ctx)

	// Get the specified client
	client, err := s.getSpecificMusicClient(ctx, userID, clientID)
	if err != nil {
		return models.MediaItem[*mediatypes.Track]{}, err
	}

	// Check if client implements MusicProvider
	musicClient, ok := client.(providers.MusicProvider)
	if !ok {
		log.Warn().
			Uint64("clientID", clientID).
			Str("clientType", "media client").
			Msg("Client does not implement music provider interface")
		return models.MediaItem[*mediatypes.Track]{}, fmt.Errorf("client does not implement music provider interface")
	}

	// Call client directly for specific track
	track, err := musicClient.GetMusicTrackByID(ctx, trackID)
	if err != nil {
		log.Error().Err(err).
			Uint64("clientID", clientID).
			Str("trackID", trackID).
			Msg("Failed to retrieve track")
		return models.MediaItem[*mediatypes.Track]{}, err
	}

	return *track, nil
}

func (s *mediaMusicService[T]) GetRecentlyAddedAlbums(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[*mediatypes.Album], error) {
	// Get music clients
	clients, err := s.getMusicClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Int("clientCount", len(clients)).
		Int("limit", limit).
		Msg("Getting recently added albums across clients")

	var allAlbums []*models.MediaItem[*mediatypes.Album]

	for _, client := range clients {
		musicClient, ok := client.(providers.MusicProvider)
		if !ok {
			continue
		}

		// Configure for recently added
		options := &mediatypes.QueryOptions{
			RecentlyAdded: true,
			Sort:          "dateAdded",
			SortOrder:     mediatypes.SortOrderDesc,
			Limit:         limit,
		}

		albums, err := musicClient.GetMusicAlbums(ctx, options)
		if err != nil {
			// Log error but continue with other clients
			log.Warn().
				Err(err).
				Uint64("clientID", client.GetClientID()).
				Msg("Error getting recently added albums from client")
			continue
		}

		allAlbums = append(allAlbums, albums...)
	}

	// Sort by date added (newest first)
	sort.Slice(allAlbums, func(i, j int) bool {
		return allAlbums[i].Data.Details.AddedAt.After(allAlbums[j].Data.Details.AddedAt)
	})

	// Apply limit
	if len(allAlbums) > limit {
		allAlbums = allAlbums[:limit]
	}

	log.Debug().
		Int("albumCount", len(allAlbums)).
		Msg("Retrieved recently added albums")

	return allAlbums, nil
}

func (s *mediaMusicService[T]) GetRecentlyPlayedTracks(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[*mediatypes.Track], error) {
	// Get music clients
	clients, err := s.getMusicClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Int("clientCount", len(clients)).
		Int("limit", limit).
		Msg("Getting recently played tracks across clients")

	var allTracks []*models.MediaItem[*mediatypes.Track]

	// Cut-off date (e.g., last 30 days)
	cutoffDate := time.Now().AddDate(0, 0, -30)

	for _, client := range clients {
		musicClient, ok := client.(providers.MusicProvider)
		if !ok {
			continue
		}

		// Configure for recently played
		options := &mediatypes.QueryOptions{
			RecentlyPlayed: true,
			PlayedAfter:    cutoffDate,
			Sort:           "datePlayed",
			SortOrder:      mediatypes.SortOrderDesc,
			Limit:          limit,
		}

		tracks, err := musicClient.GetMusic(ctx, options)
		if err != nil {
			// Log error but continue with other clients
			log.Warn().
				Err(err).
				Uint64("clientID", client.GetClientID()).
				Msg("Error getting recently played tracks from client")
			continue
		}

		allTracks = append(allTracks, tracks...)
	}

	// Sort by date played (newest first) if tracks have played date
	// This would require tracks to have a PlayedAt field which isn't modeled yet
	// For now, rely on the client's sorting

	// Apply limit
	if len(allTracks) > limit {
		allTracks = allTracks[:limit]
	}

	log.Debug().
		Int("trackCount", len(allTracks)).
		Msg("Retrieved recently played tracks")

	return allTracks, nil
}

func (s *mediaMusicService[T]) GetAlbumsByGenre(ctx context.Context, userID uint64, genre string) ([]*models.MediaItem[*mediatypes.Album], error) {
	// Get music clients
	clients, err := s.getMusicClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Int("clientCount", len(clients)).
		Str("genre", genre).
		Msg("Getting albums by genre across clients")

	var allAlbums []*models.MediaItem[*mediatypes.Album]

	for _, client := range clients {
		musicClient, ok := client.(providers.MusicProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Genre: genre,
		}

		albums, err := musicClient.GetMusicAlbums(ctx, options)
		if err != nil {
			// Log error but continue with other clients
			log.Warn().
				Err(err).
				Uint64("clientID", client.GetClientID()).
				Str("genre", genre).
				Msg("Error getting albums by genre from client")
			continue
		}

		allAlbums = append(allAlbums, albums...)
	}

	log.Debug().
		Int("albumCount", len(allAlbums)).
		Str("genre", genre).
		Msg("Retrieved albums by genre")

	return allAlbums, nil
}

func (s *mediaMusicService[T]) GetArtistsByGenre(ctx context.Context, userID uint64, genre string) ([]*models.MediaItem[*mediatypes.Artist], error) {
	// Get music clients
	clients, err := s.getMusicClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Int("clientCount", len(clients)).
		Str("genre", genre).
		Msg("Getting artists by genre across clients")

	var allArtists []*models.MediaItem[*mediatypes.Artist]

	for _, client := range clients {
		musicClient, ok := client.(providers.MusicProvider)
		if !ok {
			continue
		}

		options := &mediatypes.QueryOptions{
			Genre: genre,
		}

		artists, err := musicClient.GetMusicArtists(ctx, options)
		if err != nil {
			// Log error but continue with other clients
			log.Warn().
				Err(err).
				Uint64("clientID", client.GetClientID()).
				Str("genre", genre).
				Msg("Error getting artists by genre from client")
			continue
		}

		allArtists = append(allArtists, artists...)
	}

	log.Debug().
		Int("artistCount", len(allArtists)).
		Str("genre", genre).
		Msg("Retrieved artists by genre")

	return allArtists, nil
}

func (s *mediaMusicService[T]) GetRandomAlbums(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[*mediatypes.Album], error) {
	// Get music clients
	clients, err := s.getMusicClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Int("clientCount", len(clients)).
		Int("limit", limit).
		Msg("Getting random albums across clients")

	var allAlbums []*models.MediaItem[*mediatypes.Album]

	// Calculate limit per client to evenly distribute
	clientLimit := limit / len(clients)
	if clientLimit < 1 {
		clientLimit = 1
	}

	for _, client := range clients {
		musicClient, ok := client.(providers.MusicProvider)
		if !ok {
			continue
		}

		// No specific options for random, some clients might support "random" as a sort
		options := &mediatypes.QueryOptions{
			Sort:  "random", // Some clients might support this
			Limit: clientLimit,
		}

		albums, err := musicClient.GetMusicAlbums(ctx, options)
		if err != nil {
			// Log error but continue with other clients
			log.Warn().
				Err(err).
				Uint64("clientID", client.GetClientID()).
				Msg("Error getting random albums from client")
			continue
		}

		allAlbums = append(allAlbums, albums...)
	}

	// Shuffle the albums to randomize across clients
	// In a production environment, you might want a better shuffling algorithm
	/*
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(allAlbums), func(i, j int) {
			allAlbums[i], allAlbums[j] = allAlbums[j], allAlbums[i]
		})
	*/

	// Apply limit
	if len(allAlbums) > limit {
		allAlbums = allAlbums[:limit]
	}

	log.Debug().
		Int("albumCount", len(allAlbums)).
		Msg("Retrieved random albums")

	return allAlbums, nil
}

// GetTopAlbums retrieves top albums across all clients
func (s *mediaMusicService[T]) GetTopAlbums(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[*mediatypes.Album], error) {
	// Get music clients
	clients, err := s.getMusicClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Int("clientCount", len(clients)).
		Int("limit", limit).
		Msg("Getting top albums across clients")

	var allAlbums []*models.MediaItem[*mediatypes.Album]

	for _, client := range clients {
		musicClient, ok := client.(providers.MusicProvider)
		if !ok {
			continue
		}

		// Configure for top albums (by play count or rating)
		options := &mediatypes.QueryOptions{
			Sort:      "playCount", // or "rating" depending on what the client supports
			SortOrder: mediatypes.SortOrderDesc,
			Limit:     limit,
		}

		albums, err := musicClient.GetMusicAlbums(ctx, options)
		if err != nil {
			// Try a different sort if the first fails
			options.Sort = "rating"
			albums, err = musicClient.GetMusicAlbums(ctx, options)
			if err != nil {
				// Log error but continue with other clients
				log.Warn().
					Err(err).
					Uint64("clientID", client.GetClientID()).
					Msg("Error getting top albums from client")
				continue
			}
		}

		allAlbums = append(allAlbums, albums...)
	}

	// Sort by rating or playCount if available
	// This would require the album model to have these fields

	// Apply limit
	if len(allAlbums) > limit {
		allAlbums = allAlbums[:limit]
	}

	log.Debug().
		Int("albumCount", len(allAlbums)).
		Msg("Retrieved top albums")

	return allAlbums, nil
}

// GetTopAlbumsForClient retrieves top albums from a specific client
func (s *mediaMusicService[T]) GetTopAlbumsForClient(ctx context.Context, userID uint64, clientID uint64, limit int) ([]*models.MediaItem[*mediatypes.Album], error) {
	log := utils.LoggerFromContext(ctx)

	// Get the specific client
	client, err := s.getSpecificMusicClient(ctx, userID, clientID)
	if err != nil {
		return nil, err
	}

	// Check if client implements MusicProvider
	musicClient, ok := client.(providers.MusicProvider)
	if !ok {
		return nil, fmt.Errorf("client does not implement music provider interface")
	}

	// Configure for top albums
	options := &mediatypes.QueryOptions{
		Sort:      "playCount", // or "rating" depending on what the client supports
		SortOrder: mediatypes.SortOrderDesc,
		Limit:     limit,
	}

	albums, err := musicClient.GetMusicAlbums(ctx, options)
	if err != nil {
		// Try a different sort if the first fails
		options.Sort = "rating"
		albums, err = musicClient.GetMusicAlbums(ctx, options)
		if err != nil {
			log.Error().
				Err(err).
				Uint64("clientID", clientID).
				Msg("Failed to retrieve top albums")
			return nil, err
		}
	}

	log.Debug().
		Int("albumCount", len(albums)).
		Uint64("clientID", clientID).
		Msg("Retrieved top albums from client")

	return albums, nil
}

func (s *mediaMusicService[T]) SearchMusic(ctx context.Context, userID uint64, query string) (MusicSearchResults, error) {
	// This is a generalized search that returns mixed results (tracks, albums, artists)
	// Get music clients
	clients, err := s.getMusicClients(ctx, userID)
	if err != nil {
		return MusicSearchResults{}, err
	}

	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Int("clientCount", len(clients)).
		Str("query", query).
		Msg("Searching music across clients")

	results := MusicSearchResults{
		Artists: []*models.MediaItem[*mediatypes.Artist]{},
		Albums:  []*models.MediaItem[*mediatypes.Album]{},
		Tracks:  []*models.MediaItem[*mediatypes.Track]{},
	}

	for _, client := range clients {
		musicClient, ok := client.(providers.MusicProvider)
		if !ok {
			continue
		}

		// Search for tracks
		trackOptions := &mediatypes.QueryOptions{
			Query: query,
			Limit: 10, // Limit per type
		}

		tracks, err := musicClient.GetMusic(ctx, trackOptions)
		if err == nil {
			// Add track results
			results.Tracks = append(results.Tracks, tracks...)
		}

		// Search for albums
		albumOptions := &mediatypes.QueryOptions{
			Query: query,
			Limit: 5, // Fewer albums than tracks
		}

		albums, err := musicClient.GetMusicAlbums(ctx, albumOptions)
		if err == nil {
			// Add album results
			results.Albums = append(results.Albums, albums...)
		}

		// Search for artists
		artistOptions := &mediatypes.QueryOptions{
			Query: query,
			Limit: 5, // Fewer artists than tracks
		}

		artists, err := musicClient.GetMusicArtists(ctx, artistOptions)
		if err == nil {
			// Add artist results
			results.Artists = append(results.Artists, artists...)
		}
	}

	log.Debug().
		Int("artistsCount", len(results.Artists)).
		Int("albumsCount", len(results.Albums)).
		Int("tracksCount", len(results.Tracks)).
		Str("query", query).
		Msg("Retrieved search results")

	return results, nil
}

// GetTracksByAlbum retrieves all tracks from a specific album
func (s *mediaMusicService[T]) GetTracksByAlbum(ctx context.Context, userID uint64, clientID uint64, albumID string) ([]*models.MediaItem[*mediatypes.Track], error) {
	log := utils.LoggerFromContext(ctx)

	// Get the specific client
	client, err := s.getSpecificMusicClient(ctx, userID, clientID)
	if err != nil {
		return nil, err
	}

	// Check if client implements MusicProvider
	musicClient, ok := client.(providers.MusicProvider)
	if !ok {
		return nil, fmt.Errorf("client does not implement music provider interface")
	}

	// Configure query options to filter by album ID
	options := &mediatypes.QueryOptions{
		ExternalSourceID: albumID, // Use the album ID as a filter
	}

	// Get tracks associated with the album
	tracks, err := musicClient.GetMusic(ctx, options)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", clientID).
			Str("albumID", albumID).
			Msg("Failed to retrieve tracks for album")
		return nil, err
	}

	log.Debug().
		Int("trackCount", len(tracks)).
		Str("albumID", albumID).
		Msg("Retrieved tracks for album")

	return tracks, nil
}

// GetAlbumsByArtist retrieves all albums by a specific artist
func (s *mediaMusicService[T]) GetAlbumsByArtist(ctx context.Context, userID uint64, clientID uint64, artistID string) ([]*models.MediaItem[*mediatypes.Album], error) {
	log := utils.LoggerFromContext(ctx)

	// Get the specific client
	client, err := s.getSpecificMusicClient(ctx, userID, clientID)
	if err != nil {
		return nil, err
	}

	// Check if client implements MusicProvider
	musicClient, ok := client.(providers.MusicProvider)
	if !ok {
		return nil, fmt.Errorf("client does not implement music provider interface")
	}

	// Configure query options to filter by artist ID
	options := &mediatypes.QueryOptions{
		ExternalSourceID: artistID, // This might need to be adapted based on the client's API
	}

	// Get albums associated with the artist
	albums, err := musicClient.GetMusicAlbums(ctx, options)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", clientID).
			Str("artistID", artistID).
			Msg("Failed to retrieve albums for artist")
		return nil, err
	}

	log.Debug().
		Int("albumCount", len(albums)).
		Str("artistID", artistID).
		Msg("Retrieved albums for artist")

	return albums, nil
}

// GetTracksByGenre retrieves tracks by genre across all clients
func (s *mediaMusicService[T]) GetTracksByGenre(ctx context.Context, userID uint64, genre string) ([]*models.MediaItem[*mediatypes.Track], error) {
	// Get music clients
	clients, err := s.getMusicClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	log := utils.LoggerFromContext(ctx)
	var allTracks []*models.MediaItem[*mediatypes.Track]

	for _, client := range clients {
		musicClient, ok := client.(providers.MusicProvider)
		if !ok {
			continue
		}

		// Configure for genre filter
		options := &mediatypes.QueryOptions{
			Genre: genre,
			Limit: 50, // Reasonable limit per client
		}

		tracks, err := musicClient.GetMusic(ctx, options)
		if err != nil {
			// Log error but continue with other clients
			log.Warn().
				Err(err).
				Uint64("clientID", client.GetClientID()).
				Str("genre", genre).
				Msg("Error getting tracks by genre from client")
			continue
		}

		allTracks = append(allTracks, tracks...)
	}

	log.Debug().
		Int("trackCount", len(allTracks)).
		Str("genre", genre).
		Msg("Retrieved tracks by genre")

	return allTracks, nil
}

// GetAlbumsByYear retrieves albums released in a specific year
func (s *mediaMusicService[T]) GetAlbumsByYear(ctx context.Context, userID uint64, year int) ([]*models.MediaItem[*mediatypes.Album], error) {
	// Get music clients
	clients, err := s.getMusicClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	log := utils.LoggerFromContext(ctx)
	var allAlbums []*models.MediaItem[*mediatypes.Album]

	for _, client := range clients {
		musicClient, ok := client.(providers.MusicProvider)
		if !ok {
			continue
		}

		// Configure for year filter
		options := &mediatypes.QueryOptions{
			Year:  year,
			Limit: 50, // Reasonable limit per client
		}

		albums, err := musicClient.GetMusicAlbums(ctx, options)
		if err != nil {
			// Log error but continue with other clients
			log.Warn().
				Err(err).
				Uint64("clientID", client.GetClientID()).
				Int("year", year).
				Msg("Error getting albums by year from client")
			continue
		}

		allAlbums = append(allAlbums, albums...)
	}

	log.Debug().
		Int("albumCount", len(allAlbums)).
		Int("year", year).
		Msg("Retrieved albums by year")

	return allAlbums, nil
}

// GetLatestAlbumsByAdded retrieves the most recently added albums
// This is an alias for GetRecentlyAddedAlbums to maintain API compatibility
func (s *mediaMusicService[T]) GetLatestAlbumsByAdded(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[*mediatypes.Album], error) {
	return s.GetRecentlyAddedAlbums(ctx, userID, limit)
}

// GetPopularAlbums retrieves popular albums based on play count or ratings
func (s *mediaMusicService[T]) GetPopularAlbums(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[*mediatypes.Album], error) {
	// In practice, this is similar to GetTopAlbums but might have a different sorting algorithm
	// For simplicity, we'll just use the GetTopAlbums method for now
	return s.GetTopAlbums(ctx, userID, limit)
}

// GetPopularArtists retrieves popular artists based on play count or ratings
func (s *mediaMusicService[T]) GetPopularArtists(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[*mediatypes.Artist], error) {
	// Get music clients
	clients, err := s.getMusicClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	log := utils.LoggerFromContext(ctx)
	var allArtists []*models.MediaItem[*mediatypes.Artist]

	for _, client := range clients {
		musicClient, ok := client.(providers.MusicProvider)
		if !ok {
			continue
		}

		// Configure for popularity sort
		options := &mediatypes.QueryOptions{
			Sort:      "popularity", // This should be supported by the client
			SortOrder: mediatypes.SortOrderDesc,
			Limit:     limit,
		}

		artists, err := musicClient.GetMusicArtists(ctx, options)
		if err != nil {
			// Try a different approach if the first fails
			options.Sort = "playCount"
			artists, err = musicClient.GetMusicArtists(ctx, options)
			if err != nil {
				// Log error but continue with other clients
				log.Warn().
					Err(err).
					Uint64("clientID", client.GetClientID()).
					Msg("Error getting popular artists from client")
				continue
			}
		}

		allArtists = append(allArtists, artists...)
	}

	// Sort by some popularity metric (if available)
	// This would require the Artist type to have such a field

	// Limit to requested count
	if len(allArtists) > limit {
		allArtists = allArtists[:limit]
	}

	log.Debug().
		Int("artistCount", len(allArtists)).
		Msg("Retrieved popular artists")

	return allArtists, nil
}

// GetTopArtists retrieves top artists from a specific client
func (s *mediaMusicService[T]) GetTopArtists(ctx context.Context, userID uint64, clientID uint64, limit int) ([]*models.MediaItem[*mediatypes.Artist], error) {
	log := utils.LoggerFromContext(ctx)

	// Get the specific client
	client, err := s.getSpecificMusicClient(ctx, userID, clientID)
	if err != nil {
		return nil, err
	}

	// Check if client implements MusicProvider
	musicClient, ok := client.(providers.MusicProvider)
	if !ok {
		return nil, fmt.Errorf("client does not implement music provider interface")
	}

	// Configure for top artists
	options := &mediatypes.QueryOptions{
		Sort:      "popularity", // This should be supported by the client
		SortOrder: mediatypes.SortOrderDesc,
		Limit:     limit,
	}

	artists, err := musicClient.GetMusicArtists(ctx, options)
	if err != nil {
		// Try a different approach if the first fails
		options.Sort = "playCount"
		artists, err = musicClient.GetMusicArtists(ctx, options)
		if err != nil {
			log.Error().
				Err(err).
				Uint64("clientID", clientID).
				Msg("Failed to retrieve top artists")
			return nil, err
		}
	}

	log.Debug().
		Int("artistCount", len(artists)).
		Uint64("clientID", clientID).
		Msg("Retrieved top artists from client")

	return artists, nil
}

// GetFavoriteArtists retrieves favorite artists from a specific client
func (s *mediaMusicService[T]) GetFavoriteArtists(ctx context.Context, userID uint64, clientID uint64, limit int) ([]*models.MediaItem[*mediatypes.Artist], error) {
	log := utils.LoggerFromContext(ctx)

	// Get the specific client
	client, err := s.getSpecificMusicClient(ctx, userID, clientID)
	if err != nil {
		return nil, err
	}

	// Check if client implements MusicProvider
	musicClient, ok := client.(providers.MusicProvider)
	if !ok {
		return nil, fmt.Errorf("client does not implement music provider interface")
	}

	// Configure for favorite artists
	options := &mediatypes.QueryOptions{
		Favorites: true,
		Limit:     limit,
	}

	artists, err := musicClient.GetMusicArtists(ctx, options)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", clientID).
			Msg("Failed to retrieve favorite artists")
		return nil, err
	}

	log.Debug().
		Int("artistCount", len(artists)).
		Uint64("clientID", clientID).
		Msg("Retrieved favorite artists from client")

	return artists, nil
}

// GetTopTracks retrieves top tracks from a specific client
func (s *mediaMusicService[T]) GetTopTracks(ctx context.Context, userID uint64, clientID uint64, limit int) ([]*models.MediaItem[*mediatypes.Track], error) {
	log := utils.LoggerFromContext(ctx)

	// Get the specific client
	client, err := s.getSpecificMusicClient(ctx, userID, clientID)
	if err != nil {
		return nil, err
	}

	// Check if client implements MusicProvider
	musicClient, ok := client.(providers.MusicProvider)
	if !ok {
		return nil, fmt.Errorf("client does not implement music provider interface")
	}

	// Configure for top tracks
	options := &mediatypes.QueryOptions{
		Sort:      "playCount", // This should be supported by the client
		SortOrder: mediatypes.SortOrderDesc,
		Limit:     limit,
	}

	tracks, err := musicClient.GetMusic(ctx, options)
	if err != nil {
		// Try a different approach if the first fails
		options.Sort = "popularity"
		tracks, err = musicClient.GetMusic(ctx, options)
		if err != nil {
			log.Error().
				Err(err).
				Uint64("clientID", clientID).
				Msg("Failed to retrieve top tracks")
			return nil, err
		}
	}

	log.Debug().
		Int("trackCount", len(tracks)).
		Uint64("clientID", clientID).
		Msg("Retrieved top tracks from client")

	return tracks, nil
}

// GetRecentlyAddedTracks retrieves recently added tracks from a specific client
func (s *mediaMusicService[T]) GetRecentlyAddedTracks(ctx context.Context, userID uint64, clientID uint64, limit int) ([]*models.MediaItem[*mediatypes.Track], error) {
	log := utils.LoggerFromContext(ctx)

	// Get the specific client
	client, err := s.getSpecificMusicClient(ctx, userID, clientID)
	if err != nil {
		return nil, err
	}

	// Check if client implements MusicProvider
	musicClient, ok := client.(providers.MusicProvider)
	if !ok {
		return nil, fmt.Errorf("client does not implement music provider interface")
	}

	// Configure for recently added tracks
	options := &mediatypes.QueryOptions{
		RecentlyAdded: true,
		Sort:          "dateAdded",
		SortOrder:     mediatypes.SortOrderDesc,
		Limit:         limit,
	}

	tracks, err := musicClient.GetMusic(ctx, options)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", clientID).
			Msg("Failed to retrieve recently added tracks")
		return nil, err
	}

	log.Debug().
		Int("trackCount", len(tracks)).
		Uint64("clientID", clientID).
		Msg("Retrieved recently added tracks from client")

	return tracks, nil
}
