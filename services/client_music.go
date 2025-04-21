package services

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"suasor/client"
	"suasor/client/media/providers"
	mediatypes "suasor/client/media/types"
	"suasor/client/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/utils"
)

// ClientMusicService defines operations for interacting with music clients
type ClientMusicService[T types.ClientConfig] interface {
	GetClientAlbumByID(ctx context.Context, clientID uint64, albumID string) (models.MediaItem[*mediatypes.Album], error)
	GetClientArtistByID(ctx context.Context, clientID uint64, artistID string) (models.MediaItem[*mediatypes.Artist], error)
	GetClientTrackByID(ctx context.Context, clientID uint64, trackID string) (models.MediaItem[*mediatypes.Track], error)

	GetClientSimilarTracks(ctx context.Context, clientID uint64, trackID string, limit int) ([]*models.MediaItem[*mediatypes.Track], error)
	GetClientSimilarArtists(ctx context.Context, clientID uint64, artistID string, limit int) ([]*models.MediaItem[*mediatypes.Artist], error)

	GetClientArtistsByGenre(ctx context.Context, clientID uint64, genre string) ([]*models.MediaItem[*mediatypes.Artist], error)
	GetClientAlbumsByGenre(ctx context.Context, clientID uint64, genre string) ([]*models.MediaItem[*mediatypes.Album], error)
	GetClientTracksByGenre(ctx context.Context, clientID uint64, genre string) ([]*models.MediaItem[*mediatypes.Track], error)

	GetClientAlbumsByYear(ctx context.Context, clientID uint64, year int) ([]*models.MediaItem[*mediatypes.Album], error)
	GetClientRandomAlbums(ctx context.Context, clientID uint64, limit int) ([]*models.MediaItem[*mediatypes.Album], error)

	GetClientTopAlbums(ctx context.Context, clientID uint64, limit int) ([]*models.MediaItem[*mediatypes.Album], error)
	GetClientTopArtists(ctx context.Context, clientID uint64, limit int) ([]*models.MediaItem[*mediatypes.Artist], error)
	GetClientTopTracks(ctx context.Context, clientID uint64, limit int) ([]*models.MediaItem[*mediatypes.Track], error)

	GetClientTracksByAlbum(ctx context.Context, clientID uint64, albumID string) ([]*models.MediaItem[*mediatypes.Track], error)
	GetClientAlbumsByArtist(ctx context.Context, clientID uint64, artistID string) ([]*models.MediaItem[*mediatypes.Album], error)

	GetClientRecentlyPlayedTracks(ctx context.Context, clientID uint64, limit int) ([]*models.MediaItem[*mediatypes.Track], error)
	GetClientRecentlyAddedTracks(ctx context.Context, clientID uint64, limit int) ([]*models.MediaItem[*mediatypes.Track], error)
	GetClientRecentlyAddedAlbums(ctx context.Context, clientID uint64, limit int) ([]*models.MediaItem[*mediatypes.Album], error)

	GetClientFavoriteAlbums(ctx context.Context, clientID uint64, limit int) ([]*models.MediaItem[*mediatypes.Album], error)
	GetClientFavoriteArtists(ctx context.Context, clientID uint64, limit int) ([]*models.MediaItem[*mediatypes.Artist], error)
	GetClientFavoriteTracks(ctx context.Context, clientID uint64, limit int) ([]*models.MediaItem[*mediatypes.Track], error)

	SearchMusic(ctx context.Context, clientID uint64, query *mediatypes.QueryOptions) (*MusicSearchResults, error)
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
	clientRepo    repository.ClientRepository[T]
	clientFactory *client.ClientFactoryService
}

func NewClientMusicService[T types.ClientMediaConfig](
	clientRepo repository.ClientRepository[T],
	clientFactory *client.ClientFactoryService,
) ClientMusicService[T] {
	return &mediaMusicService[T]{
		clientRepo:    clientRepo,
		clientFactory: clientFactory,
	}
}

// getMusicProviders gets all music clients for a user
func (s *mediaMusicService[T]) getMusicProviders(ctx context.Context, userID uint64) ([]providers.MusicProvider, error) {
	// GetClient all media clients for the user that support music
	clients, err := s.clientRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	log := utils.LoggerFromContext(ctx)
	log.Debug().Msg("GetClientting music clients")

	var musicClients []providers.MusicProvider

	// Filter for clients that support music and instantiate them
	for _, clientConfig := range clients {
		if clientConfig.Config.Data.SupportsMusic() {
			log.Debug().
				Uint64("clientID", clientConfig.ID).
				Str("clientType", clientConfig.Config.Data.GetClientType().String()).
				Msg("Found music-supporting client")

			provider, err := s.clientFactory.GetMusicProvider(ctx, clientConfig.ID, clientConfig.Config.Data)
			if err != nil {
				log.Error().
					Err(err).
					Uint64("clientID", clientConfig.ID).
					Msg("Failed to instantiate music client")
				continue
			}

			musicClients = append(musicClients, provider)
		}
	}

	log.Debug().
		Int("musicClientCount", len(musicClients)).
		Msg("Retrieved music clients")

	return musicClients, nil
}

// getMusicProvider gets a specific music client
func (s *mediaMusicService[T]) getMusicProvider(ctx context.Context, clientID uint64) (providers.MusicProvider, error) {
	log := utils.LoggerFromContext(ctx)

	clientConfig, err := s.clientRepo.GetByID(ctx, clientID)
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
			Str("clientType", clientConfig.Config.Data.GetClientType().String()).
			Msg("Client does not support music")
		return nil, errors.New("client does not support music")
	}

	provider, err := s.clientFactory.GetMusicProvider(ctx, clientID, clientConfig.Config.Data)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", clientID).
			Str("clientType", clientConfig.Config.Data.GetClientType().String()).
			Msg("Failed to instantiate music client")
		return nil, err
	}

	return provider, nil
}

func (s *mediaMusicService[T]) GetClientAlbumByID(ctx context.Context, clientID uint64, albumID string) (models.MediaItem[*mediatypes.Album], error) {
	log := utils.LoggerFromContext(ctx)

	// GetClient the specified client
	provider, err := s.getMusicProvider(ctx, clientID)
	if err != nil {
		return models.MediaItem[*mediatypes.Album]{}, err
	}
	options := &mediatypes.QueryOptions{
		ExternalSourceID: albumID,
	}

	albums, err := provider.GetMusicAlbums(ctx, options)
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

func (s *mediaMusicService[T]) GetClientArtistByID(ctx context.Context, clientID uint64, artistID string) (models.MediaItem[*mediatypes.Artist], error) {
	log := utils.LoggerFromContext(ctx)

	// GetClient the specified client
	client, err := s.getMusicProvider(ctx, clientID)
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

func (s *mediaMusicService[T]) GetClientTrackByID(ctx context.Context, clientID uint64, trackID string) (models.MediaItem[*mediatypes.Track], error) {
	log := utils.LoggerFromContext(ctx)

	// GetClient the specified client
	client, err := s.getMusicProvider(ctx, clientID)
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

func (s *mediaMusicService[T]) GetClientRecentlyAddedAlbums(ctx context.Context, clientID uint64, limit int) ([]*models.MediaItem[*mediatypes.Album], error) {
	// GetClient music clients
	provider, err := s.getMusicProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	log := utils.LoggerFromContext(ctx)
	var allAlbums []*models.MediaItem[*mediatypes.Album]

	// Configure for recently added
	options := &mediatypes.QueryOptions{
		RecentlyAdded: true,
		Sort:          "dateAdded",
		SortOrder:     mediatypes.SortOrderDesc,
		Limit:         limit,
	}

	albums, err := provider.GetMusicAlbums(ctx, options)
	if err != nil {
		// Log error but continue with other clients
		log.Warn().
			Err(err).
			Uint64("clientID", clientID).
			Msg("Error getting recently added albums from client")
	}

	allAlbums = append(allAlbums, albums...)

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

func (s *mediaMusicService[T]) GetClientRecentlyPlayedTracks(ctx context.Context, clientID uint64, limit int) ([]*models.MediaItem[*mediatypes.Track], error) {
	// GetClient music clients
	provider, err := s.getMusicProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Int("limit", limit).
		Msg("GetClientting recently played tracks across clients")

	var allTracks []*models.MediaItem[*mediatypes.Track]

	// Cut-off date (e.g., last 30 days)
	cutoffDate := time.Now().AddDate(0, 0, -30)

	// Configure for recently played
	options := &mediatypes.QueryOptions{
		RecentlyPlayed: true,
		PlayedAfter:    cutoffDate,
		Sort:           "datePlayed",
		SortOrder:      mediatypes.SortOrderDesc,
		Limit:          limit,
	}

	tracks, err := provider.GetMusic(ctx, options)
	if err != nil {
		// Log error but continue with other clients
		log.Warn().
			Err(err).
			Uint64("clientID", clientID).
			Msg("Error getting recently played tracks from client")
	}

	log.Debug().
		Int("trackCount", len(allTracks)).
		Msg("Retrieved recently played tracks")

	return tracks, nil
}

func (s *mediaMusicService[T]) GetClientAlbumsByGenre(ctx context.Context, clientID uint64, genre string) ([]*models.MediaItem[*mediatypes.Album], error) {
	// GetClient music clients
	provider, err := s.getMusicProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Str("genre", genre).
		Msg("GetClientting albums by genre across clients")

	options := &mediatypes.QueryOptions{
		Genre: genre,
	}
	albums, err := provider.GetMusicAlbums(ctx, options)
	if err != nil {
		// Log error but continue with other clients
		log.Warn().
			Err(err).
			Uint64("clientID", clientID).
			Str("genre", genre).
			Msg("Error getting albums by genre from client")
		return nil, err
	}

	return albums, nil
}

func (s *mediaMusicService[T]) GetClientArtistsByGenre(ctx context.Context, clientID uint64, genre string) ([]*models.MediaItem[*mediatypes.Artist], error) {
	// GetClient music clients
	provider, err := s.getMusicProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Str("genre", genre).
		Msg("GetClientting artists by genre across clients")

	options := &mediatypes.QueryOptions{
		Genre: genre,
	}

	artists, err := provider.GetMusicArtists(ctx, options)
	if err != nil {
		// Log error but continue with other clients
		log.Warn().
			Err(err).
			Str("genre", genre).
			Msg("Error getting artists by genre from client")
	}

	log.Debug().
		Str("genre", genre).
		Msg("Retrieved artists by genre")

	return artists, nil
}

func (s *mediaMusicService[T]) GetClientRandomAlbums(ctx context.Context, clientID uint64, limit int) ([]*models.MediaItem[*mediatypes.Album], error) {
	// GetClient music clients
	provider, err := s.getMusicProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Int("limit", limit).
		Msg("GetClientting random albums across clients")

	// No specific options for random, some clients might support "random" as a sort
	options := &mediatypes.QueryOptions{
		Sort:  "random", // Some clients might support this
		Limit: limit,
	}

	albums, err := provider.GetMusicAlbums(ctx, options)
	if err != nil {
		// Log error but continue with other clients
		log.Warn().
			Err(err).
			Uint64("clientID", clientID).
			Msg("Error getting random albums from client")
	}

	log.Debug().
		Msg("Retrieved random albums")

	return albums, nil
}

// GetClientTopAlbums retrieves top albums across all clients
func (s *mediaMusicService[T]) GetClientTopAlbums(ctx context.Context, clientID uint64, limit int) ([]*models.MediaItem[*mediatypes.Album], error) {
	// GetClient music clients
	provider, err := s.getMusicProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Int("limit", limit).
		Uint64("clientID", clientID).
		Msg("GetClientting top albums across clients")

	// Configure for top albums (by play count or rating)
	options := &mediatypes.QueryOptions{
		Sort:      "playCount", // or "rating" depending on what the client supports
		SortOrder: mediatypes.SortOrderDesc,
		Limit:     limit,
	}

	albums, err := provider.GetMusicAlbums(ctx, options)
	if err != nil {
		// Try a different sort if the first fails
		options.Sort = "rating"
		albums, err = provider.GetMusicAlbums(ctx, options)
		if err != nil {
			// Log error but continue with other clients
			log.Warn().
				Err(err).
				Uint64("clientID", clientID).
				Msg("Error getting top albums from client")
		}
	}

	log.Debug().
		Msg("Retrieved top albums")

	return albums, nil
}

// GetClientTopAlbumsForClient retrieves top albums from a specific client
func (s *mediaMusicService[T]) GetClientTopAlbumsForClient(ctx context.Context, clientID uint64, limit int) ([]*models.MediaItem[*mediatypes.Album], error) {
	log := utils.LoggerFromContext(ctx)

	// GetClient the specific client
	provider, err := s.getMusicProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	// Configure for top albums
	options := &mediatypes.QueryOptions{
		Sort:      "playCount", // or "rating" depending on what the client supports
		SortOrder: mediatypes.SortOrderDesc,
		Limit:     limit,
	}

	albums, err := provider.GetMusicAlbums(ctx, options)
	if err != nil {
		// Try a different sort if the first fails
		options.Sort = "rating"
		albums, err = provider.GetMusicAlbums(ctx, options)
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

func (s *mediaMusicService[T]) SearchMusic(ctx context.Context, clientID uint64, query *mediatypes.QueryOptions) (*MusicSearchResults, error) {
	// This is a generalized search that returns mixed results (tracks, albums, artists)
	// GetClient music clients
	clients, err := s.getMusicProviders(ctx, clientID)
	if err != nil {
		return nil, err
	}

	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Int("clientCount", len(clients)).
		Str("query", query.Query).
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
		if query.Limit == 0 {
			query.Limit = 10
		}

		tracks, err := musicClient.GetMusic(ctx, query)
		if err == nil {
			// Add track results
			results.Tracks = append(results.Tracks, tracks...)
		}

		// Search for albums (limit is 5)
		query.Limit = 5
		albums, err := musicClient.GetMusicAlbums(ctx, query)

		if err == nil {
			// Add album results
			results.Albums = append(results.Albums, albums...)
		}

		artists, err := musicClient.GetMusicArtists(ctx, query)
		if err == nil {
			// Add artist results
			results.Artists = append(results.Artists, artists...)
		}
	}

	log.Debug().
		Int("artistsCount", len(results.Artists)).
		Int("albumsCount", len(results.Albums)).
		Int("tracksCount", len(results.Tracks)).
		Str("query", query.Query).
		Msg("Retrieved search results")

	return &results, nil
}

// GetClientTracksByAlbum retrieves all tracks from a specific album
func (s *mediaMusicService[T]) GetClientTracksByAlbum(ctx context.Context, clientID uint64, albumID string) ([]*models.MediaItem[*mediatypes.Track], error) {
	log := utils.LoggerFromContext(ctx)

	// GetClient the specific client
	client, err := s.getMusicProvider(ctx, clientID)
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

	// GetClient tracks associated with the album
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

// GetClientAlbumsByArtist retrieves all albums by a specific artist
func (s *mediaMusicService[T]) GetClientAlbumsByArtist(ctx context.Context, clientID uint64, artistID string) ([]*models.MediaItem[*mediatypes.Album], error) {
	log := utils.LoggerFromContext(ctx)

	// GetClient the specific client
	client, err := s.getMusicProvider(ctx, clientID)
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

	// GetClient albums associated with the artist
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

// GetClientTracksByGenre retrieves tracks by genre across all clients
func (s *mediaMusicService[T]) GetClientTracksByGenre(ctx context.Context, clientID uint64, genre string) ([]*models.MediaItem[*mediatypes.Track], error) {
	// GetClient music clients
	provider, err := s.getMusicProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	log := utils.LoggerFromContext(ctx)

	// Configure for genre filter
	options := &mediatypes.QueryOptions{
		Genre: genre,
		Limit: 50, // Reasonable limit per client
	}

	tracks, err := provider.GetMusic(ctx, options)
	if err != nil {
		// Log error but continue with other clients
		log.Warn().
			Err(err).
			Uint64("clientID", clientID).
			Str("genre", genre).
			Msg("Error getting tracks by genre from client")
	}

	log.Debug().
		Str("genre", genre).
		Msg("Retrieved tracks by genre")

	return tracks, nil
}

// GetClientAlbumsByYear retrieves albums released in a specific year
func (s *mediaMusicService[T]) GetClientAlbumsByYear(ctx context.Context, clientID uint64, year int) ([]*models.MediaItem[*mediatypes.Album], error) {
	// GetClient music clients
	provider, err := s.getMusicProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	log := utils.LoggerFromContext(ctx)

	// Configure for year filter
	options := &mediatypes.QueryOptions{
		Year:  year,
		Limit: 50, // Reasonable limit per client
	}

	albums, err := provider.GetMusicAlbums(ctx, options)

	if err != nil {
		// Log error but continue with other clients
		log.Warn().
			Err(err).
			Uint64("clientID", clientID).
			Int("year", year).
			Msg("Error getting albums by year from client")
	}

	log.Debug().
		Int("year", year).
		Msg("Retrieved albums by year")

	return albums, nil
}

// GetClientPopularArtists retrieves popular artists based on play count or ratings
func (s *mediaMusicService[T]) GetClientTopArtists(ctx context.Context, clientID uint64, limit int) ([]*models.MediaItem[*mediatypes.Artist], error) {
	// GetClient music clients
	provider, err := s.getMusicProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	log := utils.LoggerFromContext(ctx)
	var allArtists []*models.MediaItem[*mediatypes.Artist]

	// Configure for popularity sort
	options := &mediatypes.QueryOptions{
		Sort:      "popularity", // This should be supported by the client
		SortOrder: mediatypes.SortOrderDesc,
		Limit:     limit,
	}

	artists, err := provider.GetMusicArtists(ctx, options)
	if err != nil {
		// Try a different approach if the first fails
		options.Sort = "playCount"
		artists, err = provider.GetMusicArtists(ctx, options)
		if err != nil {
			// Log error but continue with other clients
			log.Warn().
				Err(err).
				Msg("Error getting popular artists from client")
		}
	}

	// Sort by some popularity metric (if available)
	// This would require the Artist type to have such a field

	// Limit to requested count

	log.Debug().
		Int("artistCount", len(allArtists)).
		Msg("Retrieved popular artists")

	return artists, nil
}

// GetClientFavoriteArtists retrieves favorite artists from a specific client
func (s *mediaMusicService[T]) GetClientFavoriteArtists(ctx context.Context, clientID uint64, limit int) ([]*models.MediaItem[*mediatypes.Artist], error) {
	log := utils.LoggerFromContext(ctx)

	// GetClient the specific client
	client, err := s.getMusicProvider(ctx, clientID)
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

// GetClientTopTracks retrieves top tracks from a specific client
func (s *mediaMusicService[T]) GetClientTopTracks(ctx context.Context, clientID uint64, limit int) ([]*models.MediaItem[*mediatypes.Track], error) {
	log := utils.LoggerFromContext(ctx)

	// GetClient the specific client
	client, err := s.getMusicProvider(ctx, clientID)
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

// GetClientRecentlyAddedTracks retrieves recently added tracks from a specific client
func (s *mediaMusicService[T]) GetClientRecentlyAddedTracks(ctx context.Context, clientID uint64, limit int) ([]*models.MediaItem[*mediatypes.Track], error) {
	log := utils.LoggerFromContext(ctx)

	// GetClient the specific client
	client, err := s.getMusicProvider(ctx, clientID)
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

func (s *mediaMusicService[T]) GetClientSimilarTracks(ctx context.Context, clientID uint64, trackID string, limit int) ([]*models.MediaItem[*mediatypes.Track], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("clientID", clientID).
		Str("trackID", trackID).
		Int("limit", limit).
		Msg("GetClientting similar tracks")

	// GetClient the specific client
	client, err := s.getMusicProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	// Check if client implements MusicProvider
	musicClient, ok := client.(providers.MusicProvider)
	if !ok {
		return nil, fmt.Errorf("client does not implement music provider interface")
	}

	// Configure query options
	options := &mediatypes.QueryOptions{
		ExternalSourceID: trackID,
		Limit:            limit,
	}

	tracks, err := musicClient.GetMusic(ctx, options)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", clientID).
			Str("trackID", trackID).
			Msg("Failed to retrieve similar tracks")
		return nil, err
	}

	log.Debug().
		Int("trackCount", len(tracks)).
		Uint64("clientID", clientID).
		Str("trackID", trackID).
		Msg("Retrieved similar tracks from client")

	return tracks, nil
}

// GetClientFavoriteTracks godoc
// @Summary GetClient favorite tracks
// @Description Retrieves the user's favorite tracks from a client
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param limit query int false "Maximum number of tracks to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Track]] "Favorite tracks retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /clients/media/{clientID}/music/tracks/favorites [get]
func (s *mediaMusicService[T]) GetClientFavoriteTracks(ctx context.Context, clientID uint64, limit int) ([]*models.MediaItem[*mediatypes.Track], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("clientID", clientID).
		Int("limit", limit).
		Msg("GetClientting favorite tracks")

	// GetClient the specific client
	client, err := s.getMusicProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	// Check if client implements MusicProvider
	musicClient, ok := client.(providers.MusicProvider)
	if !ok {
		return nil, fmt.Errorf("client does not implement music provider interface")
	}

	// Configure query options
	options := &mediatypes.QueryOptions{
		Favorites: true,
		Limit:     limit,
	}

	tracks, err := musicClient.GetMusic(ctx, options)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", clientID).
			Msg("Failed to retrieve favorite tracks")
		return nil, err
	}

	log.Debug().
		Int("trackCount", len(tracks)).
		Uint64("clientID", clientID).
		Msg("Retrieved favorite tracks from client")

	return tracks, nil
}

// GetClientFavoriteAlbums godoc
// @Summary GetClient favorite albums
// @Description Retrieves the user's favorite albums from a client
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param limit query int false "Maximum number of albums to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Album]] "Favorite albums retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /clients/media/{clientID}/music/albums/favorites [get]
func (s *mediaMusicService[T]) GetClientFavoriteAlbums(ctx context.Context, clientID uint64, limit int) ([]*models.MediaItem[*mediatypes.Album], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("clientID", clientID).
		Int("limit", limit).
		Msg("GetClientting favorite albums")

	// GetClient the specific client
	client, err := s.getMusicProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	// Check if client implements MusicProvider
	musicClient, ok := client.(providers.MusicProvider)
	if !ok {
		return nil, fmt.Errorf("client does not implement music provider interface")
	}

	// Configure query options
	options := &mediatypes.QueryOptions{
		Favorites: true,
		Limit:     limit,
	}

	albums, err := musicClient.GetMusicAlbums(ctx, options)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", clientID).
			Msg("Failed to retrieve favorite albums")
		return nil, err
	}

	log.Debug().
		Int("albumCount", len(albums)).
		Uint64("clientID", clientID).
		Msg("Retrieved favorite albums from client")

	return albums, nil
}

// GetClientSimiarArtists godoc
// @Summary GetClient similar artists
// @Description Retrieves artists similar to a specific artist from a client
// @Tags music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientID path int true "Client ID"
// @Param artistID path string true "Artist ID"
// @Param limit query int false "Maximum number of artists to return (default 10)"
// @Success 200 {object} responses.APIResponse[[]models.MediaItem[mediatypes.Artist]] "Similar artists retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[error] "Invalid request"
// @Failure 401 {object} responses.ErrorResponse[error] "Unauthorized"
// @Failure 500 {object} responses.ErrorResponse[error] "Server error"
// @Router /clients/media/{clientID}/music/artists/{artistID}/similar [get]
func (s *mediaMusicService[T]) GetClientSimilarArtists(ctx context.Context, clientID uint64, artistID string, limit int) ([]*models.MediaItem[*mediatypes.Artist], error) {
	log := utils.LoggerFromContext(ctx)
	log.Debug().
		Uint64("clientID", clientID).
		Str("artistID", artistID).
		Int("limit", limit).
		Msg("GetClientting similar artists")

	// GetClient the specific client
	client, err := s.getMusicProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	// Check if client implements MusicProvider
	musicClient, ok := client.(providers.MusicProvider)
	if !ok {
		return nil, fmt.Errorf("client does not implement music provider interface")
	}

	// Configure query options
	options := &mediatypes.QueryOptions{
		ExternalSourceID: artistID,
		Limit:            limit,
	}

	artists, err := musicClient.GetMusicArtists(ctx, options)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", clientID).
			Str("artistID", artistID).
			Msg("Failed to retrieve similar artists")
		return nil, err
	}

	log.Debug().
		Int("artistCount", len(artists)).
		Uint64("clientID", clientID).
		Str("artistID", artistID).
		Msg("Retrieved similar artists from client")

	return artists, nil
}
