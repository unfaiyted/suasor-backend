package responses

import (
	"suasor/clients/media/types"
	"suasor/types/models"
	"time"
)

// SearchResults holds results from a search operation
type SearchResults struct {
	Movies      []*models.MediaItem[*types.Movie]      `json:"movies"`
	Series      []*models.MediaItem[*types.Series]     `json:"series"`
	Episodes    []*models.MediaItem[*types.Episode]    `json:"episodes,omitempty"`
	Tracks      []*models.MediaItem[*types.Track]      `json:"tracks,omitempty"`
	Albums      []*models.MediaItem[*types.Album]      `json:"albums,omitempty"`
	Artists     []*models.MediaItem[*types.Artist]     `json:"artists,omitempty"`
	Collections []*models.MediaItem[*types.Collection] `json:"collections,omitempty"`
	Playlists   []*models.MediaItem[*types.Playlist]   `json:"playlists,omitempty"`
	People      []*models.Person                       `json:"people,omitempty"`
	TotalCount  int                                    `json:"totalCount"`
}

// SearchResponse is the response for search operations
type SearchResponse struct {
	Success bool `json:"success"`
	Results SearchResults
}

// RecentSearchHistoryItem represents a single search history item in responses
type RecentSearchHistoryItem struct {
	ID          uint64    `json:"id"`
	Query       string    `json:"query"`
	ResultCount int       `json:"resultCount"`
	SearchedAt  time.Time `json:"searchedAt"`
}

// RecentSearchesResponse is the response for recent searches
type RecentSearchesResponse struct {
	Success  bool                      `json:"success"`
	Searches []RecentSearchHistoryItem `json:"searches"`
}

// TrendingSearchItem represents a trending search with count
type TrendingSearchItem struct {
	Query       string `json:"query"`
	SearchCount int    `json:"searchCount"`
}

// TrendingSearchesResponse is the response for trending searches
type TrendingSearchesResponse struct {
	Success  bool                 `json:"success"`
	Searches []TrendingSearchItem `json:"searches"`
}

// SearchSuggestionsResponse is the response for search suggestions
type SearchSuggestionsResponse struct {
	Success     bool     `json:"success"`
	Suggestions []string `json:"suggestions"`
}

// ConvertToSearchResponse converts service search results to API response format
func ConvertToSearchResponse(results SearchResults) SearchResponse {
	return SearchResponse{
		Success: true,
		Results: results,
	}
}

// ConvertToRecentSearchesResponse converts search history to API response format
func ConvertToRecentSearchesResponse(searches []models.SearchHistory) RecentSearchesResponse {
	response := RecentSearchesResponse{
		Success:  true,
		Searches: make([]RecentSearchHistoryItem, 0, len(searches)),
	}

	for _, search := range searches {
		response.Searches = append(response.Searches, RecentSearchHistoryItem{
			ID:          search.ID,
			Query:       search.Query,
			ResultCount: search.ResultCount,
			SearchedAt:  search.SearchedAt,
		})
	}

	return response
}

// ConvertToTrendingSearchesResponse converts trending search data to API response format
func ConvertToTrendingSearchesResponse(searches []models.SearchHistory) TrendingSearchesResponse {
	response := TrendingSearchesResponse{
		Success:  true,
		Searches: make([]TrendingSearchItem, 0, len(searches)),
	}

	// The trending searches already have count information from the repository
	for _, search := range searches {
		response.Searches = append(response.Searches, TrendingSearchItem{
			Query:       search.Query,
			SearchCount: search.ResultCount, // ResultCount is repurposed as the count in trending searches
		})
	}

	return response
}
