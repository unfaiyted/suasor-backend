package services

import (
	"context"
	"fmt"
	apprepos "suasor/app/repository"
	"suasor/client"
	mediaclient "suasor/client/media"
	"suasor/client/media/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/types/responses"
	"suasor/utils"
)

// SearchService provides search capabilities across the application
type SearchService interface {
	// SearchAll searches all sources for media based on the query options
	SearchAll(ctx context.Context, userID uint64, options types.QueryOptions) (responses.SearchResults, error)

	// SearchMedia searches local database media items
	SearchMedia(ctx context.Context, userID uint64, options types.QueryOptions) (responses.SearchResults, error)

	// SearchClientMedias searches all media clients for a user
	SearchClientMedias(ctx context.Context, userID uint64, options types.QueryOptions) (responses.SearchResults, error)

	// SearchMetadataClients searches all metadata clients for a user
	SearchMetadataClients(ctx context.Context, userID uint64, options types.QueryOptions) (responses.SearchResults, error)

	// GetRecentSearches gets a user's recent searches
	GetRecentSearches(ctx context.Context, userID uint64, limit int) ([]models.SearchHistory, error)

	// GetTrendingSearches gets popular searches across all users
	GetTrendingSearches(ctx context.Context, limit int) ([]models.SearchHistory, error)

	// GetSearchSuggestions gets search suggestions based on partial input
	GetSearchSuggestions(ctx context.Context, partialQuery string, limit int) ([]string, error)
}

// searchService implements SearchService
type searchService struct {
	searchRepo           repository.SearchRepository
	clientRepos          apprepos.ClientRepositories
	itemRepos            apprepos.CoreMediaItemRepositories
	personRepo           repository.PersonRepository
	clientFactoryService *client.ClientFactoryService
}

// NewSearchService creates a new search service instance
func NewSearchService(
	searchRepo repository.SearchRepository,
	clientRepos apprepos.ClientRepositories,
	itemRepos apprepos.CoreMediaItemRepositories,
	personRepo repository.PersonRepository,
	clientFactoryService *client.ClientFactoryService,
) SearchService {
	return &searchService{
		searchRepo:           searchRepo,
		clientRepos:          clientRepos,
		itemRepos:            itemRepos,
		personRepo:           personRepo,
		clientFactoryService: clientFactoryService,
	}
}

// SearchAll searches all available sources for media based on query options
func (s *searchService) SearchAll(ctx context.Context, userID uint64, options types.QueryOptions) (responses.SearchResults, error) {
	log := utils.LoggerFromContext(ctx)
	log.Info().Str("query", options.Query).Msg("Performing search across all sources")

	// Search local database first
	dbResults, err := s.SearchMedia(ctx, userID, options)
	if err != nil {
		log.Error().Err(err).Msg("Error searching local database")
		// Continue with other searches despite error
	}

	// Search media clients
	clientResults, err := s.SearchClientMedias(ctx, userID, options)
	if err != nil {
		log.Error().Err(err).Msg("Error searching media clients")
		// Continue with other searches despite error
	}

	// Search metadata clients
	metadataResults, err := s.SearchMetadataClients(ctx, userID, options)
	if err != nil {
		log.Error().Err(err).Msg("Error searching metadata clients")
		// Continue with other searches despite error
	}

	// Merge results, ensuring we don't have duplicates
	results := s.mergeSearchResults(dbResults, clientResults, metadataResults)

	// Save the search to history
	_, err = s.searchRepo.SaveSearchHistory(ctx, userID, options.Query, results.TotalCount)
	if err != nil {
		log.Error().Err(err).Msg("Failed to save search history")
		// Continue despite error
	}

	return results, nil
}

// SearchMedia searches local database media items
func (s *searchService) SearchMedia(ctx context.Context, userID uint64, options types.QueryOptions) (responses.SearchResults, error) {
	log := utils.LoggerFromContext(ctx)
	log.Info().Str("query", options.Query).Msg("Searching local database")

	results := responses.SearchResults{}

	// Search movies
	movies, err := s.itemRepos.MovieRepo().Search(ctx, options)
	if err != nil {
		log.Error().Err(err).Msg("Error searching movies")
	} else {
		results.Movies = movies
	}

	// Search series
	series, err := s.itemRepos.SeriesRepo().Search(ctx, options)
	if err != nil {
		log.Error().Err(err).Msg("Error searching series")
	} else {
		results.Series = series
	}

	// Search episodes if requested
	if options.MediaType == "episode" || options.MediaType == "" {
		episodes, err := s.itemRepos.EpisodeRepo().Search(ctx, options)
		if err != nil {
			log.Error().Err(err).Msg("Error searching episodes")
		} else {
			results.Episodes = episodes
		}
	}

	// Search music if requested
	if options.MediaType == "music" || options.MediaType == "" {
		// Search tracks
		tracks, err := s.itemRepos.TrackRepo().Search(ctx, options)
		if err != nil {
			log.Error().Err(err).Msg("Error searching tracks")
		} else {
			results.Tracks = tracks
		}

		// Search albums
		albums, err := s.itemRepos.AlbumRepo().Search(ctx, options)
		if err != nil {
			log.Error().Err(err).Msg("Error searching albums")
		} else {
			results.Albums = albums
		}

		// Search artists
		artists, err := s.itemRepos.ArtistRepo().Search(ctx, options)
		if err != nil {
			log.Error().Err(err).Msg("Error searching artists")
		} else {
			results.Artists = artists
		}
	}

	// Search collections
	collections, err := s.itemRepos.CollectionRepo().Search(ctx, options)
	if err != nil {
		log.Error().Err(err).Msg("Error searching collections")
	} else {
		results.Collections = collections
	}

	// Search playlists
	playlists, err := s.itemRepos.PlaylistRepo().Search(ctx, options)
	if err != nil {
		log.Error().Err(err).Msg("Error searching playlists")
	} else {
		results.Playlists = playlists
	}

	// Search people
	if options.MediaType == "person" || options.MediaType == "" {
		people, err := s.personRepo.Search(ctx, options)
		if err != nil {
			log.Error().Err(err).Msg("Error searching people")
		} else {
			results.People = people
		}
	}

	// Calculate total count
	results.TotalCount = s.calculateTotalCount(results)

	return results, nil
}

// SearchClientMedias searches media clients for a user
func (s *searchService) SearchClientMedias(ctx context.Context, userID uint64, options types.QueryOptions) (responses.SearchResults, error) {
	log := utils.LoggerFromContext(ctx)
	log.Info().Uint64("userID", userID).Str("query", options.Query).Msg("Searching media clients")

	// Get all client configs for the user
	clients, err := s.clientRepos.GetAllMediaClientsForUser(ctx, userID)
	if err != nil {
		return responses.SearchResults{}, fmt.Errorf("failed to get client configs: %w", err)
	}

	// TODO: Loop through all the clients and search them all in parallel
	// get client Connection details
	client, err := s.clientFactoryService.GetClient(ctx, clients.Emby[0].ID, clients.Emby[0].GetConfig())
	if err != nil {
		return responses.SearchResults{}, fmt.Errorf("failed to get client: %w", err)
	}

	clientMedia := client.(mediaclient.ClientMedia)

	clientMedia.Search(ctx, &options)

	results := responses.SearchResults{}

	// TODO: Implement parallel search across all media clients
	// This would involve getting client instances, checking their capabilities,
	// and calling search methods on each based on media type

	// For now, this is a placeholder implementation

	return results, nil
}

// SearchMetadataClients searches metadata clients for a user
func (s *searchService) SearchMetadataClients(ctx context.Context, userID uint64, options types.QueryOptions) (responses.SearchResults, error) {
	log := utils.LoggerFromContext(ctx)
	log.Info().Uint64("userID", userID).Str("query", options.Query).Msg("Searching metadata clients")

	// Get metadata client configs for the user
	// TODO: Implement this functionality

	results := responses.SearchResults{}

	// TODO: Implement parallel search across all metadata clients
	// This would involve getting metadata client instances, checking their capabilities,
	// and calling search methods on each based on media type

	// For now, this is a placeholder implementation

	return results, nil
}

// GetRecentSearches gets a user's recent searches
func (s *searchService) GetRecentSearches(ctx context.Context, userID uint64, limit int) ([]models.SearchHistory, error) {
	return s.searchRepo.GetRecentSearches(ctx, userID, limit)
}

// GetTrendingSearches gets popular searches across all users
func (s *searchService) GetTrendingSearches(ctx context.Context, limit int) ([]models.SearchHistory, error) {
	return s.searchRepo.GetTrendingSearches(ctx, limit)
}

// GetSearchSuggestions gets search suggestions based on partial input
func (s *searchService) GetSearchSuggestions(ctx context.Context, partialQuery string, limit int) ([]string, error) {
	return s.searchRepo.GetSearchSuggestions(ctx, partialQuery, limit)
}

// mergeSearchResults combines results from multiple sources and removes duplicates
func (s *searchService) mergeSearchResults(results ...responses.SearchResults) responses.SearchResults {
	merged := responses.SearchResults{}

	// TODO: Implement proper merging with duplicate detection
	// For now, just concatenate all results
	for _, r := range results {
		merged.Movies = append(merged.Movies, r.Movies...)
		merged.Series = append(merged.Series, r.Series...)
		merged.Episodes = append(merged.Episodes, r.Episodes...)
		merged.Tracks = append(merged.Tracks, r.Tracks...)
		merged.Albums = append(merged.Albums, r.Albums...)
		merged.Artists = append(merged.Artists, r.Artists...)
		merged.Collections = append(merged.Collections, r.Collections...)
		merged.Playlists = append(merged.Playlists, r.Playlists...)
		merged.People = append(merged.People, r.People...)
	}

	// Calculate total count
	merged.TotalCount = s.calculateTotalCount(merged)

	return merged
}

// calculateTotalCount counts all items in the search results
func (s *searchService) calculateTotalCount(results responses.SearchResults) int {
	total := 0
	total += len(results.Movies)
	total += len(results.Series)
	total += len(results.Episodes)
	total += len(results.Tracks)
	total += len(results.Albums)
	total += len(results.Artists)
	total += len(results.Collections)
	total += len(results.Playlists)
	total += len(results.People)
	return total
}
