package handlers

import (
	"github.com/gin-gonic/gin"
	// "suasor/services"
	// "suasor/types/responses"
)

// SearchHandler handles all search operations
type SearchHandler struct {
	// Probably going to need to create a search service that handles most of the logic.
}

// NewSearchHandler creates a new search handler
func NewSearchHandler() *SearchHandler {
	return &SearchHandler{}
}

// SearchMediaItems handles searching for media items
func (h *SearchHandler) Search(c *gin.Context) {
	// We should be able to search for media by a query string
	// The search process should look for local media first, then search the connected clients for items.
	// The collected clients each have seperate search methods by type.
	// We should be able to search for movies, series, and music.
	// the search result should have them seperated by type.

	// We will also search configured metadata sources to see if they have any results that were missed.
	// Initally we have tmbd setup for movies, tvdb for series, and musicbrainz will be added for music.

	// We will return a SearchResponse struct that will have a list of items in a list seperated by type.
	// Any new items will be added to our MediaItem database records.
	// We need to make sure we know if the item is in the users library or not.

	// We should save the user's search history to the database to keep track of what they are searching for and tie it back to possible recommendations based on recent searches. The query, and maybe a version of the results.
	// Considerations: this may be slow, we may not want thit to come as a single call. We may want the frontend to do more of the work to break up the search so that it seems faster. It may make sense to have multiple requests that we can break up and show users the data on the frontend as each request comes back.

	// It might be good to use methods that only search one datasource at a time and we update the frontend as they come back. think about it.
}

// GetRecentSearches returns a list of recent searches

// GetUserSearches returns a list of searches for a specific user

// SearchMediaItems searches for media items based on a query

// SearchMediaClients searches a users media clients based on a query

// SearchMediaClientsByType searches a users media clients based on a query and media type

// SearchMetadataClients searches a users metadata clients based on a query

// SearchAutomationClients searches a users automation clients based on a query

// SearchAll does every possible search combined. Might be really slow. Not ideal for user interaction, but maybe background tasks

// GetSearchSuggestions returns a list of search suggestions based on a query
