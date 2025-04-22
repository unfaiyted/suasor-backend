package services

import (
	"context"
	"fmt"
	"suasor/clients/media/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/utils/logger"
)

// CoreMusicService defines the interface for music-related operations
type CoreMusicService interface {
	// Track-related operations
	GetTracksByAlbumID(ctx context.Context, albumID uint64) ([]*models.MediaItem[*types.Track], error)
	GetTracksByArtistID(ctx context.Context, artistID uint64) ([]*models.MediaItem[*types.Track], error)
	GetTracksInPlaylist(ctx context.Context, playlistID uint64) ([]*models.MediaItem[*types.Track], error)
	GetMostPlayedTracks(ctx context.Context, limit int) ([]*models.MediaItem[*types.Track], error)
	GetRecentlyAddedTracks(ctx context.Context, days int, limit int) ([]*models.MediaItem[*types.Track], error)
	GetTracksByGenre(ctx context.Context, genre string, limit int) ([]*models.MediaItem[*types.Track], error)
	GetSimilarTracks(ctx context.Context, trackID uint64, limit int) ([]*models.MediaItem[*types.Track], error)

	// Album-related operations
	GetAlbumsByArtistID(ctx context.Context, artistID uint64) ([]*models.MediaItem[*types.Album], error)
	GetAlbumWithTracks(ctx context.Context, albumID uint64) (*models.MediaItem[*types.Album], []*models.MediaItem[*types.Track], error)
	GetRecentlyAddedAlbums(ctx context.Context, days int, limit int) ([]*models.MediaItem[*types.Album], error)
	GetMostPlayedAlbums(ctx context.Context, limit int) ([]*models.MediaItem[*types.Album], error)
	GetAlbumsByGenre(ctx context.Context, genre string, limit int) ([]*models.MediaItem[*types.Album], error)

	// Artist-related operations
	GetArtistWithAlbums(ctx context.Context, artistID uint64) (*models.MediaItem[*types.Artist], []*models.MediaItem[*types.Album], error)
	GetTopArtists(ctx context.Context, limit int) ([]*models.MediaItem[*types.Artist], error)
	GetArtistsByGenre(ctx context.Context, genre string, limit int) ([]*models.MediaItem[*types.Artist], error)

	// Search operations
	SearchMusicLibrary(ctx context.Context, query types.QueryOptions) (*models.MediaItemList, error)
}

// coreMusicService implements the CoreMusicService interface
type coreMusicService struct {
	musicRepo repository.MusicRepository
}

// NewCoreMusicService creates a new core music service
func NewCoreMusicService(musicRepo repository.MusicRepository) CoreMusicService {
	return &coreMusicService{
		musicRepo: musicRepo,
	}
}

// GetTracksByAlbumID gets all tracks for a specific album
func (s *coreMusicService) GetTracksByAlbumID(ctx context.Context, albumID uint64) ([]*models.MediaItem[*types.Track], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("albumID", albumID).
		Msg("Getting tracks by album ID")

	tracks, err := s.musicRepo.GetTracksByAlbumID(ctx, albumID)
	if err != nil {
		log.Error().Err(err).
			Uint64("albumID", albumID).
			Msg("Failed to get tracks for album")
		return nil, fmt.Errorf("failed to get tracks: %w", err)
	}

	return tracks, nil
}

// GetTracksByArtistID gets all tracks by a specific artist
func (s *coreMusicService) GetTracksByArtistID(ctx context.Context, artistID uint64) ([]*models.MediaItem[*types.Track], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("artistID", artistID).
		Msg("Getting tracks by artist ID")

	tracks, err := s.musicRepo.GetTracksByArtistID(ctx, artistID)
	if err != nil {
		log.Error().Err(err).
			Uint64("artistID", artistID).
			Msg("Failed to get tracks for artist")
		return nil, fmt.Errorf("failed to get tracks: %w", err)
	}

	return tracks, nil
}

// GetTracksInPlaylist gets all tracks in a specific playlist
func (s *coreMusicService) GetTracksInPlaylist(ctx context.Context, playlistID uint64) ([]*models.MediaItem[*types.Track], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("playlistID", playlistID).
		Msg("Getting tracks in playlist")

	tracks, err := s.musicRepo.GetTracksInPlaylist(ctx, playlistID)
	if err != nil {
		log.Error().Err(err).
			Uint64("playlistID", playlistID).
			Msg("Failed to get tracks in playlist")
		return nil, fmt.Errorf("failed to get tracks in playlist: %w", err)
	}

	return tracks, nil
}

// GetMostPlayedTracks gets the most played tracks
func (s *coreMusicService) GetMostPlayedTracks(ctx context.Context, limit int) ([]*models.MediaItem[*types.Track], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Int("limit", limit).
		Msg("Getting most played tracks")

	tracks, err := s.musicRepo.GetMostPlayedTracks(ctx, limit)
	if err != nil {
		log.Error().Err(err).
			Int("limit", limit).
			Msg("Failed to get most played tracks")
		return nil, fmt.Errorf("failed to get most played tracks: %w", err)
	}

	return tracks, nil
}

// GetRecentlyAddedTracks gets recently added tracks
func (s *coreMusicService) GetRecentlyAddedTracks(ctx context.Context, days int, limit int) ([]*models.MediaItem[*types.Track], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Int("days", days).
		Int("limit", limit).
		Msg("Getting recently added tracks")

	tracks, err := s.musicRepo.GetRecentlyAddedTracks(ctx, days, limit)
	if err != nil {
		log.Error().Err(err).
			Int("days", days).
			Int("limit", limit).
			Msg("Failed to get recently added tracks")
		return nil, fmt.Errorf("failed to get recently added tracks: %w", err)
	}

	return tracks, nil
}

// GetTracksByGenre gets tracks by genre
func (s *coreMusicService) GetTracksByGenre(ctx context.Context, genre string, limit int) ([]*models.MediaItem[*types.Track], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Str("genre", genre).
		Int("limit", limit).
		Msg("Getting tracks by genre")

	tracks, err := s.musicRepo.GetTracksByGenre(ctx, genre, limit)
	if err != nil {
		log.Error().Err(err).
			Str("genre", genre).
			Int("limit", limit).
			Msg("Failed to get tracks by genre")
		return nil, fmt.Errorf("failed to get tracks by genre: %w", err)
	}

	return tracks, nil
}

// GetSimilarTracks gets tracks similar to a given track
func (s *coreMusicService) GetSimilarTracks(ctx context.Context, trackID uint64, limit int) ([]*models.MediaItem[*types.Track], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("trackID", trackID).
		Int("limit", limit).
		Msg("Getting similar tracks")

	tracks, err := s.musicRepo.GetSimilarTracks(ctx, trackID, limit)
	if err != nil {
		log.Error().Err(err).
			Uint64("trackID", trackID).
			Int("limit", limit).
			Msg("Failed to get similar tracks")
		return nil, fmt.Errorf("failed to get similar tracks: %w", err)
	}

	return tracks, nil
}

// GetAlbumsByArtistID gets all albums by a specific artist
func (s *coreMusicService) GetAlbumsByArtistID(ctx context.Context, artistID uint64) ([]*models.MediaItem[*types.Album], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("artistID", artistID).
		Msg("Getting albums by artist ID")

	albums, err := s.musicRepo.GetAlbumsByArtistID(ctx, artistID)
	if err != nil {
		log.Error().Err(err).
			Uint64("artistID", artistID).
			Msg("Failed to get albums for artist")
		return nil, fmt.Errorf("failed to get albums: %w", err)
	}

	return albums, nil
}

// GetAlbumWithTracks gets an album and all its tracks
func (s *coreMusicService) GetAlbumWithTracks(ctx context.Context, albumID uint64) (*models.MediaItem[*types.Album], []*models.MediaItem[*types.Track], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("albumID", albumID).
		Msg("Getting album with tracks")

	album, tracks, err := s.musicRepo.GetAlbumWithTracks(ctx, albumID)
	if err != nil {
		log.Error().Err(err).
			Uint64("albumID", albumID).
			Msg("Failed to get album with tracks")
		return nil, nil, fmt.Errorf("failed to get album with tracks: %w", err)
	}

	return album, tracks, nil
}

// GetRecentlyAddedAlbums gets recently added albums
func (s *coreMusicService) GetRecentlyAddedAlbums(ctx context.Context, days int, limit int) ([]*models.MediaItem[*types.Album], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Int("days", days).
		Int("limit", limit).
		Msg("Getting recently added albums")

	albums, err := s.musicRepo.GetRecentlyAddedAlbums(ctx, days, limit)
	if err != nil {
		log.Error().Err(err).
			Int("days", days).
			Int("limit", limit).
			Msg("Failed to get recently added albums")
		return nil, fmt.Errorf("failed to get recently added albums: %w", err)
	}

	return albums, nil
}

// GetMostPlayedAlbums gets the most played albums
func (s *coreMusicService) GetMostPlayedAlbums(ctx context.Context, limit int) ([]*models.MediaItem[*types.Album], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Int("limit", limit).
		Msg("Getting most played albums")

	albums, err := s.musicRepo.GetMostPlayedAlbums(ctx, limit)
	if err != nil {
		log.Error().Err(err).
			Int("limit", limit).
			Msg("Failed to get most played albums")
		return nil, fmt.Errorf("failed to get most played albums: %w", err)
	}

	return albums, nil
}

// GetAlbumsByGenre gets albums by genre
func (s *coreMusicService) GetAlbumsByGenre(ctx context.Context, genre string, limit int) ([]*models.MediaItem[*types.Album], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Str("genre", genre).
		Int("limit", limit).
		Msg("Getting albums by genre")

	albums, err := s.musicRepo.GetAlbumsByGenre(ctx, genre, limit)
	if err != nil {
		log.Error().Err(err).
			Str("genre", genre).
			Int("limit", limit).
			Msg("Failed to get albums by genre")
		return nil, fmt.Errorf("failed to get albums by genre: %w", err)
	}

	return albums, nil
}

// GetArtistWithAlbums gets an artist and all their albums
func (s *coreMusicService) GetArtistWithAlbums(ctx context.Context, artistID uint64) (*models.MediaItem[*types.Artist], []*models.MediaItem[*types.Album], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("artistID", artistID).
		Msg("Getting artist with albums")

	artist, albums, err := s.musicRepo.GetArtistWithAlbums(ctx, artistID)
	if err != nil {
		log.Error().Err(err).
			Uint64("artistID", artistID).
			Msg("Failed to get artist with albums")
		return nil, nil, fmt.Errorf("failed to get artist with albums: %w", err)
	}

	return artist, albums, nil
}

// GetTopArtists gets the top artists
func (s *coreMusicService) GetTopArtists(ctx context.Context, limit int) ([]*models.MediaItem[*types.Artist], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Int("limit", limit).
		Msg("Getting top artists")

	artists, err := s.musicRepo.GetTopArtists(ctx, limit)
	if err != nil {
		log.Error().Err(err).
			Int("limit", limit).
			Msg("Failed to get top artists")
		return nil, fmt.Errorf("failed to get top artists: %w", err)
	}

	return artists, nil
}

// GetArtistsByGenre gets artists by genre
func (s *coreMusicService) GetArtistsByGenre(ctx context.Context, genre string, limit int) ([]*models.MediaItem[*types.Artist], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Str("genre", genre).
		Int("limit", limit).
		Msg("Getting artists by genre")

	artists, err := s.musicRepo.GetArtistsByGenre(ctx, genre, limit)
	if err != nil {
		log.Error().Err(err).
			Str("genre", genre).
			Int("limit", limit).
			Msg("Failed to get artists by genre")
		return nil, fmt.Errorf("failed to get artists by genre: %w", err)
	}

	return artists, nil
}

// SearchMusicLibrary performs a comprehensive search across all music items
func (s *coreMusicService) SearchMusicLibrary(ctx context.Context, query types.QueryOptions) (*models.MediaItemList, error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Str("query", query.Query).
		Msg("Searching music library")

	results, err := s.musicRepo.SearchMusicLibrary(ctx, query)
	if err != nil {
		log.Error().Err(err).
			Str("query", query.Query).
			Msg("Failed to search music library")
		return nil, fmt.Errorf("failed to search music library: %w", err)
	}

	return results, nil
}
