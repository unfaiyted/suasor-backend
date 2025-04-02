// helpers.go
package emby

import (
	"github.com/antihax/optional"
	"suasor/client/media/types"
	embyclient "suasor/internal/clients/embyAPI"
)

// applyQueryOptions applies the common query options to Emby API parameters
func applyQueryOptions(queryParams *embyclient.ItemsServiceApiGetItemsOpts, options *types.QueryOptions) {
	if options == nil {
		return
	}

	if options.Limit > 0 {
		queryParams.Limit = optional.NewInt32(int32(options.Limit))
	}

	if options.Offset > 0 {
		queryParams.StartIndex = optional.NewInt32(int32(options.Offset))
	}

	if options.Sort != "" {
		queryParams.SortBy = optional.NewString(options.Sort)
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

	// Apply filters
	if options.Filters != nil {
		// Media type filter
		if mediaType, ok := options.Filters["mediaType"]; ok {
			queryParams.IncludeItemTypes = optional.NewString(mediaType)
		}

		// Genre filter
		if genre, ok := options.Filters["genre"]; ok {
			queryParams.Genres = optional.NewString(genre)
		}

		// Favorite filter
		if favorite, ok := options.Filters["isFavorite"]; ok && favorite == "true" {
			queryParams.IsFavorite = optional.NewBool(true)
		}

		// Year filter
		if year, ok := options.Filters["year"]; ok {
			queryParams.Years = optional.NewString(year)
		}

		// Person filter (actor or creator)
		if actor, ok := options.Filters["actor"]; ok {
			queryParams.Person = optional.NewString(actor)
		}

		if director, ok := options.Filters["director"]; ok {
			queryParams.Person = optional.NewString(director)
		}

		if creator, ok := options.Filters["creator"]; ok {
			queryParams.Person = optional.NewString(creator)
		}
	}

	// Debug logging removed to avoid logger dependency
}
