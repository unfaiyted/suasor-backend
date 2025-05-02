// Package repository provides data access layer implementations for the application.
package repository

// MediaItemRepository represents the base repository for media items.
// This focuses on generic operations that are not specifically tied to clients or users.
// It provides the core functionality for working with media items directly.
//
// Relationships with other repositories:
// - MediaItemRepository: Core operations on media items without client or user associations
// - ClientMediaItemRepository: Operations for media items linked to specific clients
// - UserMediaItemRepository: Operations for media items owned by users (playlists, collections)
//
// This three-tier approach allows for clear separation of concerns while maintaining
// a single database table for all media items.

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"regexp"
	"strings"
	"suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/logger"
	"time"
)

// CoreMediaItemRepository defines the interface for generic media item operations
// This focuses solely on the media items themselves without user or client associations
type CoreMediaItemRepository[T types.MediaData] interface {
	GetAll(ctx context.Context, limit int, offset int, publicOnly bool) ([]*models.MediaItem[T], error)
	GetByID(ctx context.Context, id uint64) (*models.MediaItem[T], error)
	GetByIDs(ctx context.Context, ids []uint64) ([]*models.MediaItem[T], error)
	GetByUserID(ctx context.Context, userID uint64, limit int, offset int) ([]*models.MediaItem[T], error)
	GetByClientItemID(ctx context.Context, clientID uint64, clientItemID string) (*models.MediaItem[T], error)

	GetMixedMediaItemsByIDs(ctx context.Context, ids []uint64) (*models.MediaItemList, error)

	// Query operations
	GetByType(ctx context.Context, mediaType types.MediaType) ([]*models.MediaItem[T], error)
	GetByTitleAndYear(ctx context.Context, clientID uint64, title string, year int) (*models.MediaItem[T], error)
	GetByTitle(ctx context.Context, clientID uint64, title string) (*models.MediaItem[T], error)
	GetByExternalID(ctx context.Context, source string, externalID string) (*models.MediaItem[T], error)
	Search(ctx context.Context, query types.QueryOptions) ([]*models.MediaItem[T], error)

	// Specialized queries
	GetRecentItems(ctx context.Context, days int, limit int) ([]*models.MediaItem[T], error)
	GetPopularItems(ctx context.Context, limit int) ([]*models.MediaItem[T], error)
	GetItemsByAttributes(ctx context.Context, attributes map[string]interface{}, limit int) ([]*models.MediaItem[T], error)
}

type mediaItemRepository[T types.MediaData] struct {
	db *gorm.DB
}

// NewMediaItemRepository creates a new media item repository
func NewMediaItemRepository[T types.MediaData](db *gorm.DB) CoreMediaItemRepository[T] {
	return &mediaItemRepository[T]{db: db}
}

// GetByID retrieves a media item by its ID
func (r *mediaItemRepository[T]) GetByID(ctx context.Context, id uint64) (*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("id", id).
		Msg("Getting media item by ID")

	var item models.MediaItem[T]
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&item).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("media item not found")
		}
		return nil, fmt.Errorf("failed to get media item: %w", err)
	}
	return &item, nil
}

// GetMediaItemsByIDs retrieves multiple media items by their IDs
func (r *mediaItemRepository[T]) GetByIDs(ctx context.Context, ids []uint64) ([]*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Int("count", len(ids)).
		Msg("Getting media items by IDs")

	var items []*models.MediaItem[T]
	if err := r.db.WithContext(ctx).
		Where("id IN ?", ids).
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get media items by IDs: %w", err)
	}

	return items, nil
}

// GetByType retrieves all media items of a specific type
func (r *mediaItemRepository[T]) GetByType(ctx context.Context, mediaType types.MediaType) ([]*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Str("type", string(mediaType)).
		Msg("Getting media items by type")

	var items []*models.MediaItem[T]
	if err := r.db.WithContext(ctx).
		Where("type = ?", mediaType).
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get media items by type: %w", err)
	}

	return items, nil
}

// GetByExternalID retrieves a media item by an external ID
func (r *mediaItemRepository[T]) GetByExternalID(ctx context.Context, source string, externalID string) (*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Str("source", source).
		Str("externalID", externalID).
		Msg("Getting media item by external ID")

	var items []*models.MediaItem[T]

	// Use JSON contains operator to find items where externalIDs contains an entry with the given source and ID
	query := r.db.WithContext(ctx).
		Where("external_ids @> ?", fmt.Sprintf(`[{"source":"%s","id":"%s"}]`, source, externalID)).
		Find(&items)

	if err := query.Error; err != nil {
		return nil, fmt.Errorf("failed to get media item by external ID: %w", err)
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("media item not found")
	}

	// Return the first match
	return items[0], nil
}

// Search finds media items based on a query string
func (r *mediaItemRepository[T]) Search(ctx context.Context, query types.QueryOptions) ([]*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Str("query", query.Query).
		Str("type", string(query.MediaType)).
		Int("limit", query.Limit).
		Int("offset", query.Offset).
		Msg("Searching media items")

	dbQuery := r.db.WithContext(ctx)

	// Add type filter if provided
	if query.MediaType != "" {
		dbQuery = dbQuery.Where("type = ?", query.MediaType)
	}

	// Add search condition
	if query.Query != "" {
		// Use ILIKE for case-insensitive search in PostgreSQL
		// TODOL: user paramater string
		dbQuery = dbQuery.Where("title ILIKE ?", "%"+query.Query+"%")
	}

	// Add pagination
	if query.Limit > 0 {
		dbQuery = dbQuery.Limit(query.Limit)
	}

	if query.Offset > 0 {
		dbQuery = dbQuery.Offset(query.Offset)
	}

	// Order by most recently created
	dbQuery = dbQuery.Order("created_at DESC")

	var items []*models.MediaItem[T]
	if err := dbQuery.Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to search media items: %w", err)
	}

	return items, nil
}

// GetRecentItems retrieves recently added items of a specific type
func (r *mediaItemRepository[T]) GetRecentItems(ctx context.Context, days int, limit int) ([]*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)

	var zero T
	mediaType := types.GetMediaTypeFromTypeName(zero)

	log.Debug().
		Str("type", string(mediaType)).
		Int("days", days).
		Int("limit", limit).
		Msg("Getting recent media items")

	var items []*models.MediaItem[T]

	// Calculate the cutoff date
	cutoffDate := time.Now().AddDate(0, 0, -days)

	dbQuery := r.db.WithContext(ctx).
		Where("type = ?", mediaType).
		Where("created_at >= ?", cutoffDate)

	// Add limit if provided
	if limit > 0 {
		dbQuery = dbQuery.Limit(limit)
	}

	// Order by most recently created
	dbQuery = dbQuery.Order("created_at DESC")

	if err := dbQuery.Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get recent media items: %w", err)
	}

	return items, nil
}

// GetPopularItems retrieves popular items of a specific type
// Note: This implementation assumes a "play_count" or similar field in the data JSON
// You may need to adapt this based on your actual schema
func (r *mediaItemRepository[T]) GetPopularItems(ctx context.Context, limit int) ([]*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)

	var zero T
	mediaType := types.GetMediaTypeFromTypeName(zero)

	log.Debug().
		Str("type", string(mediaType)).
		Int("limit", limit).
		Msg("Getting popular media items")

	var items []*models.MediaItem[T]

	dbQuery := r.db.WithContext(ctx).
		Where("type = ?", mediaType)

	// Add an order by play_count or a similar metric from the JSON data
	// This is PostgreSQL-specific JSON path syntax
	dbQuery = dbQuery.Order("(data->>'playCount')::int DESC NULLS LAST")

	// Add limit if provided
	if limit > 0 {
		dbQuery = dbQuery.Limit(limit)
	}

	if err := dbQuery.Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get popular media items: %w", err)
	}

	return items, nil
}

// GetItemsByAttributes retrieves items matching specific attributes
func (r *mediaItemRepository[T]) GetItemsByAttributes(ctx context.Context, attributes map[string]interface{}, limit int) ([]*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Int("attributeCount", len(attributes)).
		Int("limit", limit).
		Msg("Getting media items by attributes")

	dbQuery := r.db.WithContext(ctx)

	// Add filters for each attribute
	for key, value := range attributes {
		// For JSON attributes, use the PostgreSQL JSON operators
		if key == "genre" || key == "tags" || key == "categories" {
			// Use the @> operator for array containment
			dbQuery = dbQuery.Where(fmt.Sprintf("data->'%s' @> ?", key), fmt.Sprintf(`["%v"]`, value))
		} else if key == "year" || key == "runtime" || key == "rating" {
			// These are likely numeric fields
			dbQuery = dbQuery.Where(fmt.Sprintf("data->>'%s' = ?", key), fmt.Sprintf("%v", value))
		} else {
			// For other fields, use direct column matching if it's a column, or JSON path if it's in the data
			if key == "type" || key == "title" {
				dbQuery = dbQuery.Where(fmt.Sprintf("%s = ?", key), value)
			} else {
				dbQuery = dbQuery.Where(fmt.Sprintf("data->>'%s' = ?", key), fmt.Sprintf("%v", value))
			}
		}
	}

	// Add limit if provided
	if limit > 0 {
		dbQuery = dbQuery.Limit(limit)
	}

	var items []*models.MediaItem[T]
	if err := dbQuery.Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get media items by attributes: %w", err)
	}

	return items, nil
}

func (r *mediaItemRepository[T]) GetMixedMediaItemsByIDs(ctx context.Context, ids []uint64) (*models.MediaItemList, error) {
	// Fetch movies
	movies, err := fetchMediaItemsByType[*types.Movie](ctx, r.db, ids, types.MediaTypeMovie)
	if err != nil {
		return nil, err
	}
	series, err := fetchMediaItemsByType[*types.Series](ctx, r.db, ids, types.MediaTypeSeries)
	if err != nil {
		return nil, err
	}
	episodes, err := fetchMediaItemsByType[*types.Episode](ctx, r.db, ids, types.MediaTypeEpisode)
	if err != nil {
		return nil, err
	}
	seasons, err := fetchMediaItemsByType[*types.Season](ctx, r.db, ids, types.MediaTypeSeason)
	if err != nil {
		return nil, err
	}
	tracks, err := fetchMediaItemsByType[*types.Track](ctx, r.db, ids, types.MediaTypeTrack)
	if err != nil {
		return nil, err
	}
	albums, err := fetchMediaItemsByType[*types.Album](ctx, r.db, ids, types.MediaTypeAlbum)
	if err != nil {
		return nil, err
	}
	artists, err := fetchMediaItemsByType[*types.Artist](ctx, r.db, ids, types.MediaTypeArtist)
	if err != nil {
		return nil, err
	}
	playlists, err := fetchMediaItemsByType[*types.Playlist](ctx, r.db, ids, types.MediaTypePlaylist)
	if err != nil {
		return nil, err
	}
	collections, err := fetchMediaItemsByType[*types.Collection](ctx, r.db, ids, types.MediaTypeCollection)
	if err != nil {
		return nil, err
	}
	var mediaItems models.MediaItemList

	mediaItems.AddMovieList(movies)
	mediaItems.AddSeriesList(series)
	mediaItems.AddSeasonList(seasons)
	mediaItems.AddEpisodeList(episodes)
	mediaItems.AddAlbumList(albums)
	mediaItems.AddArtistList(artists)
	mediaItems.AddTrackList(tracks)
	mediaItems.AddPlaylistList(playlists)
	mediaItems.AddCollectionList(collections)

	return &mediaItems, nil

}

func (r *mediaItemRepository[T]) GetAll(ctx context.Context, limit int, offset int, publicOnly bool) ([]*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)

	var zero T
	mediaType := types.GetMediaTypeFromTypeName(zero)

	log.Debug().
		Int("limit", limit).
		Int("offset", offset).
		Str("mediaType", string(mediaType)).
		Bool("publicOnly", publicOnly).
		Msg("Getting all media items")

	var items []*models.MediaItem[T]

	dbQuery := r.db.WithContext(ctx)

	// Add limit if provided
	if limit > 0 {
		dbQuery = dbQuery.Limit(limit)
	}

	// Add offset if provided
	if offset > 0 {
		dbQuery = dbQuery.Offset(offset)
	}
	dbQuery = dbQuery.Order("created_at DESC")
	if publicOnly {
		//TODO: validate this is correct path to this public indicator
		dbQuery = dbQuery.Where("is_public = ?", true)
	}
	if mediaType != types.MediaTypeUnknown || mediaType != types.MediaTypeAll {
		dbQuery = dbQuery.Where("type = ?", mediaType)
	}

	if err := dbQuery.Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get all media items: %w", err)
	}

	return items, nil
}

func (r *mediaItemRepository[T]) GetByClientItemID(ctx context.Context, clientID uint64, clientItemID string) (*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Str("clientItemID", clientItemID).
		Uint64("clientID", clientID).
		Msg("Getting media item by client item ID")

	var items []*models.MediaItem[T]

	// Use JSON containment operator to find the media item where SyncClients contains
	// the specified client ID and item ID
	queryArray := fmt.Sprintf(`[{"clientID":%d,"itemID":"%s"}]`, clientID, clientItemID)
	query := fmt.Sprintf(`{"clientID":%d,"itemID":"%s"}`, clientID, clientItemID)

	if err := r.db.WithContext(ctx).
		Where("sync_clients @> ?::jsonb OR sync_clients @> ?::jsonb", queryArray, query).
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get media item by client item ID: %w", err)
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("media item not found")
	}

	// Return the first match (should be unique, but we're being defensive)
	return items[0], nil
}

func (r *mediaItemRepository[T]) GetByUserID(ctx context.Context, userID uint64, limit int, offset int) ([]*models.MediaItem[T], error) {
	var items []*models.MediaItem[T]
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("userID", userID).
		Msg("Getting media items by user ID")

	var zero T
	mediaType := types.GetMediaTypeFromTypeName(zero)

	if mediaType == types.MediaTypePlaylist || mediaType == types.MediaTypeCollection {

		// Should for now be limited to user-owned playlists and collections
		query := r.db.WithContext(ctx).
			Where("type IN (?) AND data->'itemList'->>'ownerID' = ?", mediaType, userID)

		if limit > 0 {
			query = query.Limit(limit)
		}
		if offset > 0 {
			query = query.Offset(offset)
		}

		if err := query.Find(&items).Error; err != nil {
			log.Error().Err(err).Msg("Failed to get media items")
			return nil, fmt.Errorf("failed to get media items for user: %w", err)
		}

		log.Info().
			Int("count", len(items)).
			Msg("Media items retrieved successfully")

		return items, nil
	}
	return nil, fmt.Errorf("media type not supported")

}

func (r *mediaItemRepository[T]) GetByTitleAndYear(ctx context.Context, clientID uint64, title string, year int) (*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("clientID", clientID).
		Str("title", title).
		Int("year", year).
		Msg("Getting media items by title and year")

	var zero T
	mediaType := types.GetMediaTypeFromTypeName(zero)

	// Format year as string since JSON extraction with ->> always returns text
	yearStr := fmt.Sprintf("%d", year)

	// Normalize the search title
	normalizedSearchTitle := normalizeTitle(title)
	log.Debug().Str("normalizedTitle", normalizedSearchTitle).Msg("Normalized search title")

	// Split into words for word-by-word search
	searchWords := strings.Fields(normalizedSearchTitle)
	significantWords := make([]string, 0, len(searchWords))
	for _, word := range searchWords {
		if len(word) > 2 && !isCommonWord(word) {
			significantWords = append(significantWords, word)
		}
	}

	log.Debug().Strs("significantWords", significantWords).Msg("Significant words for search")

	// 1. First try exact match (highest confidence)
	exactQuery := r.db.WithContext(ctx).
		Where("type = ?", mediaType).
		Where("data->'details'->>'releaseYear' = ?", yearStr)

	// Add title variants
	exactCondition := r.db.Where("FALSE")
	for _, variant := range generateTitleVariants(title) {
		exactCondition = exactCondition.Or("LOWER(data->'details'->>'title') = LOWER(?)", variant)
	}

	var items []*models.MediaItem[T]
	if err := exactQuery.Where(exactCondition).Find(&items).Error; err == nil && len(items) > 0 {
		log.Debug().Str("title", items[0].Data.GetDetails().Title).Msg("Found exact title match")
		return items[0], nil
	}

	// 2. Try word-by-word search with the significant words
	if len(significantWords) > 0 {
		wordQuery := r.db.WithContext(ctx).
			Where("type = ?", mediaType).
			Where("data->'details'->>'releaseYear' = ?", yearStr)

		// Build word conditions
		wordCondition := r.db.Where("FALSE")
		for _, word := range significantWords {
			if len(word) >= 3 { // Only use words of reasonable length
				wordCondition = wordCondition.Or(
					"LOWER(data->'details'->>'title') LIKE ?",
					"%"+strings.ToLower(word)+"%")
			}
		}

		if err := wordQuery.Where(wordCondition).Find(&items).Error; err != nil {
			log.Error().Err(err).Msg("Error searching by words")
		} else {
			log.Debug().Int("count", len(items)).Msg("Items found by word search")
		}

		// If we found items, score them and find the best match
		if len(items) > 0 {
			bestScore := 0.0
			var bestMatch *models.MediaItem[T]

			for _, item := range items {
				itemTitle := item.Data.GetDetails().Title
				normalizedItemTitle := normalizeTitle(itemTitle)

				// Calculate match score (0.0 to 1.0)
				score := calculateTitleSimilarity(normalizedSearchTitle, normalizedItemTitle, searchWords)

				log.Debug().
					Str("title", itemTitle).
					Str("normalized", normalizedItemTitle).
					Float64("score", score).
					Msg("Title similarity score")

				if score > bestScore {
					bestScore = score
					bestMatch = item
				}
			}

			// Only return if we have a reasonable match (score above threshold)
			const MATCH_THRESHOLD = 0.99 // 70% similarity required
			if bestMatch != nil && bestScore >= MATCH_THRESHOLD {
				log.Debug().
					Str("matchedTitle", bestMatch.Data.GetDetails().Title).
					Float64("score", bestScore).
					Msg("Found best match with acceptable score")
				return bestMatch, nil
			} else if bestMatch != nil {
				log.Debug().
					Str("bestTitle", bestMatch.Data.GetDetails().Title).
					Float64("score", bestScore).
					Float64("threshold", MATCH_THRESHOLD).
					Msg("Best match below acceptable threshold")
			}
		}
	}

	return nil, fmt.Errorf("no media item found matching title '%s' and year %d", title, year)
}

// calculateTitleSimilarity returns a score from 0.0 to 1.0 indicating how similar two titles are
func calculateTitleSimilarity(title1, title2 string, title1Words []string) float64 {
	// 1. Exact match is perfect score
	if title1 == title2 {
		return 1.0
	}

	// 2. Check word overlap
	title2Words := strings.Fields(title2)

	// Empty titles can't be compared
	if len(title1Words) == 0 || len(title2Words) == 0 {
		return 0.0
	}

	// Count matching words
	matchCount := 0
	for _, word1 := range title1Words {
		if len(word1) < 3 || isCommonWord(word1) {
			continue // Skip short/common words
		}

		for _, word2 := range title2Words {
			if word1 == word2 || strings.Contains(word2, word1) || strings.Contains(word1, word2) {
				matchCount++
				break
			}
		}
	}

	// Calculate percentage of significant words that matched
	significantWordCount := 0
	for _, w := range title1Words {
		if len(w) >= 3 && !isCommonWord(w) {
			significantWordCount++
		}
	}

	if significantWordCount == 0 {
		return 0.0
	}

	wordMatchScore := float64(matchCount) / float64(significantWordCount)

	// 3. Penalize length difference
	lengthRatio := 1.0
	if len(title1) > len(title2) {
		lengthRatio = float64(len(title2)) / float64(len(title1))
	} else {
		lengthRatio = float64(len(title1)) / float64(len(title2))
	}

	// 4. Consider if one is a substring of the other
	substringBonus := 0.0
	if strings.Contains(title1, title2) || strings.Contains(title2, title1) {
		substringBonus = 0.2 // Add 20% bonus for substring match
	}

	// Calculate final score (word match + length similarity + substring bonus)
	score := (wordMatchScore * 0.7) + (lengthRatio * 0.2) + substringBonus
	if score > 1.0 {
		score = 1.0 // Cap at 1.0
	}

	return score
}

// normalizeTitle removes all punctuation and normalizes spacing
func normalizeTitle(title string) string {
	// Convert to lowercase
	normalized := strings.ToLower(title)

	// Replace common fractions
	normalized = strings.ReplaceAll(normalized, "½", "1/2")
	normalized = strings.ReplaceAll(normalized, "¼", "1/4")
	normalized = strings.ReplaceAll(normalized, "¾", "3/4")

	// Remove ALL punctuation
	punctRegex := regexp.MustCompile(`[^\w\s]`)
	normalized = punctRegex.ReplaceAllString(normalized, " ")

	// Normalize spaces (collapse multiple spaces into one)
	spaceRegex := regexp.MustCompile(`\s+`)
	normalized = spaceRegex.ReplaceAllString(normalized, " ")

	return strings.TrimSpace(normalized)
}

// isCommonWord returns true if the word is a common word to ignore
func isCommonWord(word string) bool {
	commonWords := map[string]bool{
		"the": true, "and": true, "but": true, "for": true,
		"not": true, "you": true, "one": true, "with": true,
	}
	return commonWords[strings.ToLower(word)]
}

// generateTitleVariants creates different variants of a title to handle article placement
func generateTitleVariants(title string) []string {
	title = strings.TrimSpace(title)
	if title == "" {
		return []string{}
	}

	variants := []string{title} // Original title is always included

	// Handle articles at the beginning: "The Matrix" -> "Matrix, The"
	for _, article := range []string{"The ", "A ", "An "} {
		if strings.HasPrefix(strings.ToUpper(title), strings.ToUpper(article)) {
			remainder := title[len(article):]
			variants = append(variants, remainder+", "+strings.TrimSpace(article))
		}
	}

	// Handle articles at the end: "Matrix, The" -> "The Matrix"
	for _, article := range []string{", The", ", A", ", An"} {
		if strings.HasSuffix(strings.ToUpper(title), strings.ToUpper(article)) {
			baseTitle := title[:len(title)-len(article)]
			articleForFront := strings.TrimPrefix(article, ", ")
			variants = append(variants, articleForFront+" "+baseTitle)
		}
	}

	return variants
}

// GetByTitle returns a media item array by title
func (r *mediaItemRepository[T]) GetByTitle(ctx context.Context, clientID uint64, title string) (*models.MediaItem[T], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("clientID", clientID).
		Str("title", title).
		Msg("Getting media items by title")

	var zero T
	mediaType := types.GetMediaTypeFromTypeName(zero)

	var items []*models.MediaItem[T]
	if err := r.db.WithContext(ctx).
		Where("type = ?", mediaType).
		Where("data->'details'->>'title' = ?", title).
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get media items by title: %w", err)
	}

	return items[0], nil
}
