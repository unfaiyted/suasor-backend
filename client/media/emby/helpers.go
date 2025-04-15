// helpers.go
package emby

import (
	"fmt"
	"strings"
	"time"

	"github.com/antihax/optional"
	"suasor/client/media/types"
	embyclient "suasor/internal/clients/embyAPI"
)

// applyQueryOptions applies the common query options to Emby API parameters
func applyQueryOptions(queryParams *embyclient.ItemsServiceApiGetItemsOpts, options *types.QueryOptions) {
	if options == nil {
		return
	}

	if options.ItemIDs != "" {
		queryParams.Ids = optional.NewString(options.ItemIDs)
	}

	if options.Limit > 0 {
		queryParams.Limit = optional.NewInt32(int32(options.Limit))
	}

	if options.Offset > 0 {
		queryParams.StartIndex = optional.NewInt32(int32(options.Offset))
	}

	if options.Sort != "" {
		// TODO: Look into mapping the SortBy to emby definitions
		queryParams.SortBy = optional.NewString(string(options.Sort))
		if options.SortOrder == "desc" {
			queryParams.SortOrder = optional.NewString("Descending")
		} else {
			queryParams.SortOrder = optional.NewString("Ascending")
		}
	}

	// Apply search term (should be outside the filters check)
	if options.Query != "" {
		queryParams.SearchTerm = optional.NewString(options.Query)
		// Also enable recursive search when searching
		if !queryParams.Recursive.IsSet() {
			queryParams.Recursive = optional.NewBool(true)
		}
		// Increase the limit for search results if not explicitly set
		if options.Limit <= 0 && !queryParams.Limit.IsSet() {
			queryParams.Limit = optional.NewInt32(50) // Higher default for searches
		}
	}

	// Apply filters - use typed fields

	// Media type filter
	if options.MediaType != "" {
		queryParams.IncludeItemTypes = optional.NewString(string(options.MediaType))
	}

	// Genre filter
	if options.Genre != "" {
		queryParams.Genres = optional.NewString(options.Genre)
	}

	// Favorite filter
	if options.Favorites {
		queryParams.IsFavorite = optional.NewBool(true)
	}

	// Year filter
	if options.Year > 0 {
		queryParams.Years = optional.NewString(fmt.Sprintf("%d", options.Year))
	}

	// Person filters
	if options.Actor != "" {
		queryParams.Person = optional.NewString(options.Actor)
	}

	if options.Director != "" {
		queryParams.Person = optional.NewString(options.Director)
	}

	// Creator filter
	if options.Creator != "" {
		queryParams.Person = optional.NewString(options.Creator)
	}

	// Apply more advanced filters

	// Content rating filter
	if options.ContentRating != "" {
		queryParams.OfficialRatings = optional.NewString(options.ContentRating)
	}

	// Tags filter
	if len(options.Tags) > 0 {
		queryParams.Tags = optional.NewString(strings.Join(options.Tags, ","))
	}

	// Recently added filter
	if options.RecentlyAdded {
		queryParams.SortBy = optional.NewString("DateCreated,SortName")
		queryParams.SortOrder = optional.NewString("Descending")
	}

	// Recently played filter
	if options.RecentlyPlayed {
		queryParams.SortBy = optional.NewString("DatePlayed,SortName")
		queryParams.SortOrder = optional.NewString("Descending")
	}

	// Unwatched filter
	if options.Watched {
		queryParams.IsPlayed = optional.NewBool(true)
	}

	// Date filters
	if !options.DateAddedAfter.IsZero() {
		queryParams.MinDateLastSaved = optional.NewString(options.DateAddedAfter.Format(time.RFC3339))
	}

	if !options.DateAddedBefore.IsZero() {
		queryParams.MaxPremiereDate = optional.NewString(options.DateAddedBefore.Format(time.RFC3339))
	}

	if !options.ReleasedAfter.IsZero() {
		queryParams.MinPremiereDate = optional.NewString(options.ReleasedAfter.Format(time.RFC3339))
	}

	if !options.ReleasedBefore.IsZero() {
		queryParams.MaxPremiereDate = optional.NewString(options.ReleasedBefore.Format(time.RFC3339))
	}

	// Rating filter
	if options.MinimumRating > 0 {
		queryParams.MinCommunityRating = optional.NewFloat64(float64(options.MinimumRating))
	}

	// Debug logging removed to avoid logger dependency
}
