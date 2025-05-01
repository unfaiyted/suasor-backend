package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/logger"
)

// MusicRepository interface defines operations specific to music in the database
// This repository handles specialized music queries and operations that work with
// the relationships between tracks, albums, and artists
type MusicRepository interface {
	// Track-related operations
	GetTracksByAlbumID(ctx context.Context, albumID uint64) ([]*models.MediaItem[*types.Track], error)
	GetTracksByArtistID(ctx context.Context, artistID uint64) ([]*models.MediaItem[*types.Track], error)
	GetTracksInPlaylist(ctx context.Context, playlistID uint64) ([]*models.MediaItem[*types.Track], error)
	GetMostPlayedTracks(ctx context.Context, limit int) ([]*models.MediaItem[*types.Track], error)
	GetRecentlyAddedTracks(ctx context.Context, days int, limit int) ([]*models.MediaItem[*types.Track], error)

	// Album-related operations
	GetAlbumsByArtistID(ctx context.Context, artistID uint64) ([]*models.MediaItem[*types.Album], error)
	GetAlbumWithTracks(ctx context.Context, albumID uint64) (*models.MediaItem[*types.Album], []*models.MediaItem[*types.Track], error)
	GetRecentlyAddedAlbums(ctx context.Context, days int, limit int) ([]*models.MediaItem[*types.Album], error)
	GetMostPlayedAlbums(ctx context.Context, limit int) ([]*models.MediaItem[*types.Album], error)

	// Artist-related operations
	GetArtistWithAlbums(ctx context.Context, artistID uint64) (*models.MediaItem[*types.Artist], []*models.MediaItem[*types.Album], error)
	GetTopArtists(ctx context.Context, limit int) ([]*models.MediaItem[*types.Artist], error)

	// Genre and attribute-based operations
	GetTracksByGenre(ctx context.Context, genre string, limit int) ([]*models.MediaItem[*types.Track], error)
	GetAlbumsByGenre(ctx context.Context, genre string, limit int) ([]*models.MediaItem[*types.Album], error)
	GetArtistsByGenre(ctx context.Context, genre string, limit int) ([]*models.MediaItem[*types.Artist], error)

	GetTrackByTitleAndArtistName(ctx context.Context, title string, artistName string) (*models.MediaItem[*types.Track], error)
	GetAlbumByTitleAndArtistName(ctx context.Context, title string, artistName string) (*models.MediaItem[*types.Album], error)

	// Music external_IDs often have somthing like
	// [{ "id": "1234", "source": "musicbrainzartist"},{"id": "1234", "source":"discogartist" }]
	// These ids may show up on the album or track details
	GetArtistByExternalIDs(ctx context.Context, externalIDs types.ExternalIDs) (*models.MediaItem[*types.Artist], error)
	GetArtistTracksByExternalIDs(ctx context.Context, externalIDs types.ExternalIDs) ([]*models.MediaItem[*types.Track], error)
	GetArtistAlbumsByExternalIDs(ctx context.Context, externalIDs types.ExternalIDs) ([]*models.MediaItem[*types.Album], error)

	// Advanced search operations
	SearchMusicLibrary(ctx context.Context, query types.QueryOptions) (*models.MediaItemList, error)
	GetSimilarTracks(ctx context.Context, trackID uint64, limit int) ([]*models.MediaItem[*types.Track], error)
}

// musicRepository implements the MusicRepository interface
type musicRepository struct {
	db         *gorm.DB
	trackRepo  CoreMediaItemRepository[*types.Track]
	albumRepo  CoreMediaItemRepository[*types.Album]
	artistRepo CoreMediaItemRepository[*types.Artist]
}

// NewMusicRepository creates a new music repository
func NewMusicRepository(
	db *gorm.DB,
	trackRepo CoreMediaItemRepository[*types.Track],
	albumRepo CoreMediaItemRepository[*types.Album],
	artistRepo CoreMediaItemRepository[*types.Artist],
) MusicRepository {
	return &musicRepository{
		db:         db,
		trackRepo:  trackRepo,
		albumRepo:  albumRepo,
		artistRepo: artistRepo,
	}
}

// GetTracksByAlbumID retrieves all tracks for a specific album
func (r *musicRepository) GetTracksByAlbumID(ctx context.Context, albumID uint64) ([]*models.MediaItem[*types.Track], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("albumID", albumID).
		Msg("Getting tracks by album ID")

	var tracks []*models.MediaItem[*types.Track]
	if err := r.db.WithContext(ctx).
		Where("type = ?", types.MediaTypeTrack).
		Where("data->>'albumID' = ?", fmt.Sprint(albumID)).
		Find(&tracks).Error; err != nil {
		return nil, fmt.Errorf("failed to get tracks by album ID: %w", err)
	}

	return tracks, nil
}

// GetTracksByArtistID retrieves all tracks by a specific artist
func (r *musicRepository) GetTracksByArtistID(ctx context.Context, artistID uint64) ([]*models.MediaItem[*types.Track], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("artistID", artistID).
		Msg("Getting tracks by artist ID")

	var tracks []*models.MediaItem[*types.Track]
	if err := r.db.WithContext(ctx).
		Where("type = ?", types.MediaTypeTrack).
		Where("data->>'artistID' = ?", fmt.Sprint(artistID)).
		Find(&tracks).Error; err != nil {
		return nil, fmt.Errorf("failed to get tracks by artist ID: %w", err)
	}

	return tracks, nil
}

// GetTracksInPlaylist retrieves all tracks in a specific playlist
func (r *musicRepository) GetTracksInPlaylist(ctx context.Context, playlistID uint64) ([]*models.MediaItem[*types.Track], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("playlistID", playlistID).
		Msg("Getting tracks in playlist")

	// First get the playlist
	var playlist models.MediaItem[*types.Playlist]
	if err := r.db.WithContext(ctx).
		Where("id = ? AND type = ?", playlistID, types.MediaTypePlaylist).
		First(&playlist).Error; err != nil {
		return nil, fmt.Errorf("failed to get playlist: %w", err)
	}

	itemList := playlist.GetData().GetItemList()

	// Get the track IDs from the playlist data
	if len(itemList.Items) == 0 {
		return []*models.MediaItem[*types.Track]{}, nil
	}

	// Extract the track IDs
	var trackIDs []uint64
	for _, item := range itemList.Items {
		if item.Type != types.MediaTypeTrack {
			continue
		}
		trackIDs = append(trackIDs, item.ItemID)
	}

	// Get the tracks
	var tracks []*models.MediaItem[*types.Track]
	if err := r.db.WithContext(ctx).
		Where("id IN ? AND type = ?", trackIDs, types.MediaTypeTrack).
		Find(&tracks).Error; err != nil {
		return nil, fmt.Errorf("failed to get tracks in playlist: %w", err)
	}

	return tracks, nil
}

// GetMostPlayedTracks retrieves the most played tracks
func (r *musicRepository) GetMostPlayedTracks(ctx context.Context, limit int) ([]*models.MediaItem[*types.Track], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Int("limit", limit).
		Msg("Getting most played tracks")

	var tracks []*models.MediaItem[*types.Track]
	query := r.db.WithContext(ctx).
		Where("type = ?", types.MediaTypeTrack).
		Order("(data->>'playCount')::int DESC NULLS LAST")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&tracks).Error; err != nil {
		return nil, fmt.Errorf("failed to get most played tracks: %w", err)
	}

	return tracks, nil
}

// GetRecentlyAddedTracks retrieves recently added tracks
func (r *musicRepository) GetRecentlyAddedTracks(ctx context.Context, days int, limit int) ([]*models.MediaItem[*types.Track], error) {
	return r.trackRepo.GetRecentItems(ctx, days, limit)
}

// GetAlbumsByArtistID retrieves all albums by a specific artist
func (r *musicRepository) GetAlbumsByArtistID(ctx context.Context, artistID uint64) ([]*models.MediaItem[*types.Album], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("artistID", artistID).
		Msg("Getting albums by artist ID")

	var albums []*models.MediaItem[*types.Album]
	if err := r.db.WithContext(ctx).
		Where("type = ?", types.MediaTypeAlbum).
		Where("data->>'artistID' = ?", fmt.Sprint(artistID)).
		Find(&albums).Error; err != nil {
		return nil, fmt.Errorf("failed to get albums by artist ID: %w", err)
	}

	return albums, nil
}

// GetAlbumWithTracks retrieves an album and all its tracks
func (r *musicRepository) GetAlbumWithTracks(ctx context.Context, albumID uint64) (*models.MediaItem[*types.Album], []*models.MediaItem[*types.Track], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("albumID", albumID).
		Msg("Getting album with tracks")

	// Get the album
	album, err := r.albumRepo.GetByID(ctx, albumID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get album: %w", err)
	}

	// Get the tracks
	tracks, err := r.GetTracksByAlbumID(ctx, albumID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get tracks for album: %w", err)
	}

	return album, tracks, nil
}

// GetRecentlyAddedAlbums retrieves recently added albums
func (r *musicRepository) GetRecentlyAddedAlbums(ctx context.Context, days int, limit int) ([]*models.MediaItem[*types.Album], error) {
	return r.albumRepo.GetRecentItems(ctx, days, limit)
}

// GetMostPlayedAlbums retrieves the most played albums
func (r *musicRepository) GetMostPlayedAlbums(ctx context.Context, limit int) ([]*models.MediaItem[*types.Album], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Int("limit", limit).
		Msg("Getting most played albums")

	var albums []*models.MediaItem[*types.Album]
	query := r.db.WithContext(ctx).
		Where("type = ?", types.MediaTypeAlbum).
		Order("(data->>'playCount')::int DESC NULLS LAST")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&albums).Error; err != nil {
		return nil, fmt.Errorf("failed to get most played albums: %w", err)
	}

	return albums, nil
}

// GetArtistWithAlbums retrieves an artist and all their albums
func (r *musicRepository) GetArtistWithAlbums(ctx context.Context, artistID uint64) (*models.MediaItem[*types.Artist], []*models.MediaItem[*types.Album], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("artistID", artistID).
		Msg("Getting artist with albums")

	// Get the artist
	artist, err := r.artistRepo.GetByID(ctx, artistID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get artist: %w", err)
	}

	// Get the albums
	albums, err := r.GetAlbumsByArtistID(ctx, artistID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get albums for artist: %w", err)
	}

	return artist, albums, nil
}

// GetTopArtists retrieves the top artists based on play count
func (r *musicRepository) GetTopArtists(ctx context.Context, limit int) ([]*models.MediaItem[*types.Artist], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Int("limit", limit).
		Msg("Getting top artists")

	var artists []*models.MediaItem[*types.Artist]
	query := r.db.WithContext(ctx).
		Where("type = ?", types.MediaTypeArtist).
		Order("(data->>'playCount')::int DESC NULLS LAST")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&artists).Error; err != nil {
		return nil, fmt.Errorf("failed to get top artists: %w", err)
	}

	return artists, nil
}

// GetTracksByGenre retrieves tracks by genre
func (r *musicRepository) GetTracksByGenre(ctx context.Context, genre string, limit int) ([]*models.MediaItem[*types.Track], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Str("genre", genre).
		Int("limit", limit).
		Msg("Getting tracks by genre")

	var tracks []*models.MediaItem[*types.Track]
	query := r.db.WithContext(ctx).
		Where("type = ?", types.MediaTypeTrack).
		Where("data->'genres' ? ?", genre)

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&tracks).Error; err != nil {
		return nil, fmt.Errorf("failed to get tracks by genre: %w", err)
	}

	return tracks, nil
}

// GetAlbumsByGenre retrieves albums by genre
func (r *musicRepository) GetAlbumsByGenre(ctx context.Context, genre string, limit int) ([]*models.MediaItem[*types.Album], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Str("genre", genre).
		Int("limit", limit).
		Msg("Getting albums by genre")

	var albums []*models.MediaItem[*types.Album]
	query := r.db.WithContext(ctx).
		Where("type = ?", types.MediaTypeAlbum).
		Where("data->'genres' ? ?", genre)

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&albums).Error; err != nil {
		return nil, fmt.Errorf("failed to get albums by genre: %w", err)
	}

	return albums, nil
}

// GetArtistsByGenre retrieves artists by genre
func (r *musicRepository) GetArtistsByGenre(ctx context.Context, genre string, limit int) ([]*models.MediaItem[*types.Artist], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Str("genre", genre).
		Int("limit", limit).
		Msg("Getting artists by genre")

	var artists []*models.MediaItem[*types.Artist]
	query := r.db.WithContext(ctx).
		Where("type = ?", types.MediaTypeArtist).
		Where("data->'genres' ? ?", genre)

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&artists).Error; err != nil {
		return nil, fmt.Errorf("failed to get artists by genre: %w", err)
	}

	return artists, nil
}

// SearchMusicLibrary performs a comprehensive search across all music items
func (r *musicRepository) SearchMusicLibrary(ctx context.Context, query types.QueryOptions) (*models.MediaItemList, error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Str("query", query.Query).
		Msg("Searching music library")

	musicTypes := []types.MediaType{
		types.MediaTypeTrack,
		types.MediaTypeAlbum,
		types.MediaTypeArtist,
	}

	// Build the query for searching across all music types
	dbQuery := r.db.WithContext(ctx).
		Where("type IN ?", musicTypes)

	if query.Query != "" {
		dbQuery = dbQuery.Where(
			"title ILIKE ? OR data->>'artistName' ILIKE ? OR data->>'albumName' ILIKE ?",
			"%"+query.Query+"%", "%"+query.Query+"%", "%"+query.Query+"%",
		)
	}

	// Execute separate queries for each type to populate the MediaItems struct
	var mediaItems models.MediaItemList = models.MediaItemList{}

	// Find artists
	var artists []*models.MediaItem[*types.Artist]
	if err := dbQuery.Where("type = ?", types.MediaTypeArtist).Find(&artists).Error; err != nil {
		return nil, fmt.Errorf("failed to search artists: %w", err)
	}
	mediaItems.AddArtistList(artists)

	// Find albums
	var albums []*models.MediaItem[*types.Album]
	if err := dbQuery.Where("type = ?", types.MediaTypeAlbum).Find(&albums).Error; err != nil {
		return nil, fmt.Errorf("failed to search albums: %w", err)
	}
	mediaItems.AddAlbumList(albums)

	// Find tracks
	var tracks []*models.MediaItem[*types.Track]
	if err := dbQuery.Where("type = ?", types.MediaTypeTrack).Find(&tracks).Error; err != nil {
		return nil, fmt.Errorf("failed to search tracks: %w", err)
	}
	mediaItems.AddTrackList(tracks)

	return &mediaItems, nil
}

// GetSimilarTracks finds tracks similar to a given track based on attributes
func (r *musicRepository) GetSimilarTracks(ctx context.Context, trackID uint64, limit int) ([]*models.MediaItem[*types.Track], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("trackID", trackID).
		Int("limit", limit).
		Msg("Getting similar tracks")

	// First get the source track
	sourceTrack, err := r.trackRepo.GetByID(ctx, trackID)
	if err != nil {
		return nil, fmt.Errorf("failed to get source track: %w", err)
	}

	// Get tracks with similar genres, ignoring the source track
	var tracks []*models.MediaItem[*types.Track]
	query := r.db.WithContext(ctx).
		Where("type = ?", types.MediaTypeTrack).
		Where("id != ?", trackID).
		Where("data->>'artistID' = ? OR data->>'albumID' = ?",
			fmt.Sprint(sourceTrack.Data.ArtistID),
			fmt.Sprint(sourceTrack.Data.AlbumID))

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&tracks).Error; err != nil {
		return nil, fmt.Errorf("failed to get similar tracks: %w", err)
	}

	return tracks, nil
}

// GetTrackByTitleAndArtist gets a track by title and artist ID
func (r *musicRepository) GetTrackByTitleAndArtistName(ctx context.Context, title string, artistName string) (*models.MediaItem[*types.Track], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Str("title", title).
		Str("artistName", artistName).
		Msg("Getting track by title and artist")

	// First get the track
	var track *models.MediaItem[*types.Track]
	query := r.db.WithContext(ctx).
		Where("type = ?", types.MediaTypeTrack).
		Where("data->>'title' = ?", title).
		Where("data->>'artistID' = ?", fmt.Sprint(artistName))

	if err := query.First(&track).Error; err != nil {
		return nil, fmt.Errorf("failed to get track by title and artist: %w", err)
	}

	// TODO: Advanced similarity calculations like we did with Title/Year with Movies

	return track, nil
}

// GetAlbumByTitleAndArtistName gets an album by title and artist name
func (r *musicRepository) GetAlbumByTitleAndArtistName(ctx context.Context, title string, artistName string) (*models.MediaItem[*types.Album], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Str("title", title).
		Str("artistName", artistName).
		Msg("Getting album by title and artist")

	var album *models.MediaItem[*types.Album]
	if err := r.db.WithContext(ctx).
		Where("type = ?", types.MediaTypeAlbum).
		Where("data->>'title' = ?", title).
		Where("data->>'artistName' = ?", artistName).
		First(&album).Error; err != nil {
		return nil, fmt.Errorf("failed to get album by title and artist: %w", err)
	}

	return album, nil
}

// GetArtistByExternalIDs gets an artist by external IDs
func (r *musicRepository) GetArtistByExternalIDs(ctx context.Context, externalIDs types.ExternalIDs) (*models.MediaItem[*types.Artist], error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Msg("Retrieving media item by external IDs")

	var items []*models.MediaItem[*types.Artist]

	if len(externalIDs) == 0 {
		return nil, fmt.Errorf("no external IDs provided")
	}

	// Start building the query
	db := r.db.WithContext(ctx)
	db.Where("type = ?", types.MediaTypeArtist)
	// For the first external ID, use Where; for subsequent IDs, use Or
	for i, externalID := range externalIDs {
		jsonPattern := fmt.Sprintf(`[{"source":"%s","id":"%s"}]`, externalID.Source, externalID.ID)

		if i == 0 {
			db = db.Where("external_ids @> ?", jsonPattern)
		} else {
			db = db.Or("external_ids @> ?", jsonPattern)
		}
	}

	// Execute the query
	if err := db.Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get media items by external IDs: %w", err)
	}

	log.Info().
		Int("count", len(items)).
		Msg("Media items retrieved successfully")

	if len(items) == 0 {
		return nil, fmt.Errorf("no media item found matching external IDs")
	}

	return items[0], nil
}

// GetArtistTracksByExternalIDs gets tracks by external IDs
func (r *musicRepository) GetArtistTracksByExternalIDs(ctx context.Context, externalIDs types.ExternalIDs) ([]*models.MediaItem[*types.Track], error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Msg("Retrieving artist tracks by external IDs")

	var tracks []*models.MediaItem[*types.Track]

	if len(externalIDs) == 0 {
		return nil, fmt.Errorf("no external IDs provided")
	}

	// Start building the query
	db := r.db.WithContext(ctx)

	// For the first external ID, use Where; for subsequent IDs, use Or
	for i, externalID := range externalIDs {
		jsonPattern := fmt.Sprintf(`[{"source":"%s","id":"%s"}]`, externalID.Source, externalID.ID)

		if i == 0 {
			db = db.Where("external_ids @> ?", jsonPattern)
		} else {
			db = db.Or("external_ids @> ?", jsonPattern)
		}
	}

	// Execute the query
	if err := db.Find(&tracks).Error; err != nil {
		return nil, fmt.Errorf("failed to get artist tracks by external IDs: %w", err)
	}

	log.Info().
		Int("count", len(tracks)).
		Msg("Artist tracks retrieved successfully")

	if len(tracks) == 0 {
		return nil, fmt.Errorf("no artist track found matching external IDs")
	}

	return tracks, nil
}

// GetArtistAlbumsByExternalIDs gets albums by external IDs
func (r *musicRepository) GetArtistAlbumsByExternalIDs(ctx context.Context, externalIDs types.ExternalIDs) ([]*models.MediaItem[*types.Album], error) {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Msg("Retrieving artist albums by external IDs")

	var albums []*models.MediaItem[*types.Album]

	if len(externalIDs) == 0 {
		return nil, fmt.Errorf("no external IDs provided")
	}

	// Start building the query
	db := r.db.WithContext(ctx)

	// For the first external ID, use Where; for subsequent IDs, use Or
	for i, externalID := range externalIDs {
		jsonPattern := fmt.Sprintf(`[{"source":"%s","id":"%s"}]`, externalID.Source, externalID.ID)

		if i == 0 {
			db = db.Where("external_ids @> ?", jsonPattern)
		} else {
			db = db.Or("external_ids @> ?", jsonPattern)
		}
	}

	// Execute the query
	if err := db.Find(&albums).Error; err != nil {
		return nil, fmt.Errorf("failed to get artist albums by external IDs: %w", err)
	}

	log.Info().
		Int("count", len(albums)).
		Msg("Artist albums retrieved successfully")

	if len(albums) == 0 {
		return nil, fmt.Errorf("no artist album found matching external IDs")
	}

	return albums, nil

}
