package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	metadataTypes "suasor/clients/metadata/types"
	"time"
)

// CachedMedia represents metadata stored from external providers
type CachedMedia struct {
	ID           uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	ExternalID   string    `json:"externalId" gorm:"index;size:255"`     // ID from external provider
	MediaType    string    `json:"mediaType" gorm:"type:varchar(50)"`    // movie, tvshow, etc.
	ProviderType string    `json:"providerType" gorm:"type:varchar(50)"` // tmdb, tvdb, etc.
	Title        string    `json:"title" gorm:"size:255"`
	ReleaseDate  time.Time `json:"releaseDate"`
	Data         []byte    `json:"data" gorm:"type:jsonb"` // Complete metadata JSON
	CreatedAt    time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
	ExpiresAt    time.Time `json:"expiresAt"` // When this cache entry should expire
}

// MetadataRepository defines operations for metadata cache storage
type MetadataRepository interface {
	// Movie metadata operations
	CacheMovie(ctx context.Context, movie *metadataTypes.Movie, providerType string, expiresIn time.Duration) error
	GetCachedMovie(ctx context.Context, externalID string, providerType string) (*metadataTypes.Movie, error)
	GetCachedMoviesByReleaseDate(ctx context.Context, startDate, endDate time.Time) ([]*metadataTypes.Movie, error)

	// TV show metadata operations
	CacheTVShow(ctx context.Context, tvshow *metadataTypes.TVShow, providerType string, expiresIn time.Duration) error
	GetCachedTVShow(ctx context.Context, externalID string, providerType string) (*metadataTypes.TVShow, error)
	GetCachedTVShowsByReleaseDate(ctx context.Context, startDate, endDate time.Time) ([]*metadataTypes.TVShow, error)

	// Popular/trending operations
	CachePopularMovies(ctx context.Context, movies []*metadataTypes.Movie, providerType string, expiresIn time.Duration) error
	GetCachedPopularMovies(ctx context.Context, providerType string) ([]*metadataTypes.Movie, error)

	CachePopularTVShows(ctx context.Context, shows []*metadataTypes.TVShow, providerType string, expiresIn time.Duration) error
	GetCachedPopularTVShows(ctx context.Context, providerType string) ([]*metadataTypes.TVShow, error)

	// Cache maintenance
	CleanExpiredCache(ctx context.Context) (int, error)
}

type metadataRepository struct {
	db *gorm.DB
}

// NewMetadataRepository creates a new metadata repository
func NewMetadataRepository(db *gorm.DB) MetadataRepository {
	return &metadataRepository{
		db: db,
	}
}

// CacheMovie stores movie metadata in the cache
func (r *metadataRepository) CacheMovie(ctx context.Context, movie *metadataTypes.Movie, providerType string, expiresIn time.Duration) error {
	// Convert movie data to JSON
	movieJSON, err := json.Marshal(movie)
	if err != nil {
		return fmt.Errorf("failed to marshal movie data: %w", err)
	}

	// Set expiration time
	expiresAt := time.Now().Add(expiresIn)

	// Parse release date - handle empty strings
	var releaseDate time.Time
	if movie.ReleaseDate != "" {
		parsed, err := time.Parse("2006-01-02", movie.ReleaseDate)
		if err == nil {
			releaseDate = parsed
		}
	}

	// Create or update cache entry
	cache := CachedMedia{
		ExternalID:   movie.ID,
		MediaType:    "movie",
		ProviderType: providerType,
		Title:        movie.Title,
		ReleaseDate:  releaseDate,
		Data:         movieJSON,
		ExpiresAt:    expiresAt,
	}

	// Check if entry already exists
	var existing CachedMedia
	result := r.db.WithContext(ctx).
		Where("external_id = ? AND provider_type = ? AND media_type = ?",
			movie.ID, providerType, "movie").
		First(&existing)

	if result.Error == nil {
		// Update existing entry
		cache.ID = existing.ID
		cache.CreatedAt = existing.CreatedAt
		if err := r.db.WithContext(ctx).Save(&cache).Error; err != nil {
			return fmt.Errorf("failed to update movie cache: %w", err)
		}
	} else {
		// Create new entry
		if err := r.db.WithContext(ctx).Create(&cache).Error; err != nil {
			return fmt.Errorf("failed to create movie cache: %w", err)
		}
	}

	return nil
}

// GetCachedMovie retrieves a movie from the cache
func (r *metadataRepository) GetCachedMovie(ctx context.Context, externalID string, providerType string) (*metadataTypes.Movie, error) {
	var cache CachedMedia

	if err := r.db.WithContext(ctx).
		Where("external_id = ? AND provider_type = ? AND media_type = ? AND expires_at > ?",
			externalID, providerType, "movie", time.Now()).
		First(&cache).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("movie not found in cache or cache expired")
		}
		return nil, fmt.Errorf("failed to get movie from cache: %w", err)
	}

	// Unmarshal data
	var movie metadataTypes.Movie
	if err := json.Unmarshal(cache.Data, &movie); err != nil {
		return nil, fmt.Errorf("failed to unmarshal movie data: %w", err)
	}

	return &movie, nil
}

// GetCachedMoviesByReleaseDate retrieves movies with release dates in a specific range
func (r *metadataRepository) GetCachedMoviesByReleaseDate(ctx context.Context, startDate, endDate time.Time) ([]*metadataTypes.Movie, error) {
	var caches []CachedMedia

	if err := r.db.WithContext(ctx).
		Where("media_type = ? AND release_date BETWEEN ? AND ? AND expires_at > ?",
			"movie", startDate, endDate, time.Now()).
		Find(&caches).Error; err != nil {
		return nil, fmt.Errorf("failed to get movies by release date: %w", err)
	}

	movies := make([]*metadataTypes.Movie, 0, len(caches))
	for _, cache := range caches {
		var movie metadataTypes.Movie
		if err := json.Unmarshal(cache.Data, &movie); err != nil {
			// Log error but continue with other movies
			continue
		}
		movies = append(movies, &movie)
	}

	return movies, nil
}

// CacheTVShow stores TV show metadata in the cache
func (r *metadataRepository) CacheTVShow(ctx context.Context, tvshow *metadataTypes.TVShow, providerType string, expiresIn time.Duration) error {
	// Convert tvshow data to JSON
	tvshowJSON, err := json.Marshal(tvshow)
	if err != nil {
		return fmt.Errorf("failed to marshal TV show data: %w", err)
	}

	// Set expiration time
	expiresAt := time.Now().Add(expiresIn)

	// Parse release date - handle empty strings
	var firstAirDate time.Time
	if tvshow.FirstAirDate != "" {
		parsed, err := time.Parse("2006-01-02", tvshow.FirstAirDate)
		if err == nil {
			firstAirDate = parsed
		}
	}

	// Create or update cache entry
	cache := CachedMedia{
		ExternalID:   tvshow.ID,
		MediaType:    "tvshow",
		ProviderType: providerType,
		Title:        tvshow.Name,
		ReleaseDate:  firstAirDate,
		Data:         tvshowJSON,
		ExpiresAt:    expiresAt,
	}

	// Check if entry already exists
	var existing CachedMedia
	result := r.db.WithContext(ctx).
		Where("external_id = ? AND provider_type = ? AND media_type = ?",
			tvshow.ID, providerType, "tvshow").
		First(&existing)

	if result.Error == nil {
		// Update existing entry
		cache.ID = existing.ID
		cache.CreatedAt = existing.CreatedAt
		if err := r.db.WithContext(ctx).Save(&cache).Error; err != nil {
			return fmt.Errorf("failed to update TV show cache: %w", err)
		}
	} else {
		// Create new entry
		if err := r.db.WithContext(ctx).Create(&cache).Error; err != nil {
			return fmt.Errorf("failed to create TV show cache: %w", err)
		}
	}

	return nil
}

// GetCachedTVShow retrieves a TV show from the cache
func (r *metadataRepository) GetCachedTVShow(ctx context.Context, externalID string, providerType string) (*metadataTypes.TVShow, error) {
	var cache CachedMedia

	if err := r.db.WithContext(ctx).
		Where("external_id = ? AND provider_type = ? AND media_type = ? AND expires_at > ?",
			externalID, providerType, "tvshow", time.Now()).
		First(&cache).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("TV show not found in cache or cache expired")
		}
		return nil, fmt.Errorf("failed to get TV show from cache: %w", err)
	}

	// Unmarshal data
	var tvshow metadataTypes.TVShow
	if err := json.Unmarshal(cache.Data, &tvshow); err != nil {
		return nil, fmt.Errorf("failed to unmarshal TV show data: %w", err)
	}

	return &tvshow, nil
}

// GetCachedTVShowsByReleaseDate retrieves TV shows with release dates in a specific range
func (r *metadataRepository) GetCachedTVShowsByReleaseDate(ctx context.Context, startDate, endDate time.Time) ([]*metadataTypes.TVShow, error) {
	var caches []CachedMedia

	if err := r.db.WithContext(ctx).
		Where("media_type = ? AND release_date BETWEEN ? AND ? AND expires_at > ?",
			"tvshow", startDate, endDate, time.Now()).
		Find(&caches).Error; err != nil {
		return nil, fmt.Errorf("failed to get TV shows by release date: %w", err)
	}

	tvshows := make([]*metadataTypes.TVShow, 0, len(caches))
	for _, cache := range caches {
		var tvshow metadataTypes.TVShow
		if err := json.Unmarshal(cache.Data, &tvshow); err != nil {
			// Log error but continue with other TV shows
			continue
		}
		tvshows = append(tvshows, &tvshow)
	}

	return tvshows, nil
}

// CachePopularMovies stores popular movies in the cache
func (r *metadataRepository) CachePopularMovies(ctx context.Context, movies []*metadataTypes.Movie, providerType string, expiresIn time.Duration) error {
	// First, delete existing popular movies for this provider
	if err := r.db.WithContext(ctx).
		Where("provider_type = ? AND media_type = ? AND external_id LIKE ?",
			providerType, "movie", "popular-%").
		Delete(&CachedMedia{}).Error; err != nil {
		return fmt.Errorf("failed to clean existing popular movies: %w", err)
	}

	// Set expiration time
	expiresAt := time.Now().Add(expiresIn)

	// Store each movie with a special "popular-" prefix
	for i, movie := range movies {
		// Convert movie data to JSON
		movieJSON, err := json.Marshal(movie)
		if err != nil {
			return fmt.Errorf("failed to marshal movie data: %w", err)
		}

		// Parse release date - handle empty strings
		var releaseDate time.Time
		if movie.ReleaseDate != "" {
			parsed, err := time.Parse("2006-01-02", movie.ReleaseDate)
			if err == nil {
				releaseDate = parsed
			}
		}

		// Create cache entry with rank information in the ID
		cache := CachedMedia{
			ExternalID:   fmt.Sprintf("popular-%d-%s", i+1, movie.ID),
			MediaType:    "movie",
			ProviderType: providerType,
			Title:        movie.Title,
			ReleaseDate:  releaseDate,
			Data:         movieJSON,
			ExpiresAt:    expiresAt,
		}

		if err := r.db.WithContext(ctx).Create(&cache).Error; err != nil {
			return fmt.Errorf("failed to cache popular movie: %w", err)
		}
	}

	return nil
}

// GetCachedPopularMovies retrieves popular movies from the cache
func (r *metadataRepository) GetCachedPopularMovies(ctx context.Context, providerType string) ([]*metadataTypes.Movie, error) {
	var caches []CachedMedia

	if err := r.db.WithContext(ctx).
		Where("provider_type = ? AND media_type = ? AND external_id LIKE ? AND expires_at > ?",
			providerType, "movie", "popular-%", time.Now()).
		Order("external_id").
		Find(&caches).Error; err != nil {
		return nil, fmt.Errorf("failed to get popular movies: %w", err)
	}

	movies := make([]*metadataTypes.Movie, 0, len(caches))
	for _, cache := range caches {
		var movie metadataTypes.Movie
		if err := json.Unmarshal(cache.Data, &movie); err != nil {
			// Log error but continue with other movies
			continue
		}
		movies = append(movies, &movie)
	}

	return movies, nil
}

// CachePopularTVShows stores popular TV shows in the cache
func (r *metadataRepository) CachePopularTVShows(ctx context.Context, shows []*metadataTypes.TVShow, providerType string, expiresIn time.Duration) error {
	// First, delete existing popular TV shows for this provider
	if err := r.db.WithContext(ctx).
		Where("provider_type = ? AND media_type = ? AND external_id LIKE ?",
			providerType, "tvshow", "popular-%").
		Delete(&CachedMedia{}).Error; err != nil {
		return fmt.Errorf("failed to clean existing popular TV shows: %w", err)
	}

	// Set expiration time
	expiresAt := time.Now().Add(expiresIn)

	// Store each TV show with a special "popular-" prefix
	for i, show := range shows {
		// Convert TV show data to JSON
		showJSON, err := json.Marshal(show)
		if err != nil {
			return fmt.Errorf("failed to marshal TV show data: %w", err)
		}

		// Parse release date - handle empty strings
		var firstAirDate time.Time
		if show.FirstAirDate != "" {
			parsed, err := time.Parse("2006-01-02", show.FirstAirDate)
			if err == nil {
				firstAirDate = parsed
			}
		}

		// Create cache entry with rank information in the ID
		cache := CachedMedia{
			ExternalID:   fmt.Sprintf("popular-%d-%s", i+1, show.ID),
			MediaType:    "tvshow",
			ProviderType: providerType,
			Title:        show.Name,
			ReleaseDate:  firstAirDate,
			Data:         showJSON,
			ExpiresAt:    expiresAt,
		}

		if err := r.db.WithContext(ctx).Create(&cache).Error; err != nil {
			return fmt.Errorf("failed to cache popular TV show: %w", err)
		}
	}

	return nil
}

// GetCachedPopularTVShows retrieves popular TV shows from the cache
func (r *metadataRepository) GetCachedPopularTVShows(ctx context.Context, providerType string) ([]*metadataTypes.TVShow, error) {
	var caches []CachedMedia

	if err := r.db.WithContext(ctx).
		Where("provider_type = ? AND media_type = ? AND external_id LIKE ? AND expires_at > ?",
			providerType, "tvshow", "popular-%", time.Now()).
		Order("external_id").
		Find(&caches).Error; err != nil {
		return nil, fmt.Errorf("failed to get popular TV shows: %w", err)
	}

	shows := make([]*metadataTypes.TVShow, 0, len(caches))
	for _, cache := range caches {
		var show metadataTypes.TVShow
		if err := json.Unmarshal(cache.Data, &show); err != nil {
			// Log error but continue with other TV shows
			continue
		}
		shows = append(shows, &show)
	}

	return shows, nil
}

// CleanExpiredCache removes expired cache entries
func (r *metadataRepository) CleanExpiredCache(ctx context.Context) (int, error) {
	result := r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&CachedMedia{})

	if result.Error != nil {
		return 0, fmt.Errorf("failed to clean expired cache: %w", result.Error)
	}

	return int(result.RowsAffected), nil
}

