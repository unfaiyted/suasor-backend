package types

import (
	"fmt"
	"strings"
	"time"
)

type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
	SortOrderNone SortOrder = ""
)

type SortType string

const (
	SortTypeCreatedAt SortType = "created_at"
	SortTypeUpdatedAt SortType = "updated_at"
	SortTypeAddedAt   SortType = "added_at"
)

// QueryOptions provides parameters for filtering and pagination
// using typed fields for all filters to ensure type safety
type QueryOptions struct {
	Limit                int       `json:"limit,omitempty"`
	Offset               int       `json:"offset,omitempty"`
	Sort                 SortType  `json:"sort,omitempty"`
	SortOrder            SortOrder `json:"sortOrder,omitempty"` // "asc" or "desc"
	Query                string    `json:"query,omitempty"`
	IncludeWatchProgress bool      `json:"includeWatchProgress,omitempty"`

	// Common, typed query filters
	Favorites        bool      `json:"favorites,omitempty"`        // Filter to favorites only
	Genre            string    `json:"genre,omitempty"`            // Filter by genre
	Year             int       `json:"year,omitempty"`             // Filter by release year
	Actor            string    `json:"actor,omitempty"`            // Filter by actor name/ID
	Director         string    `json:"director,omitempty"`         // Filter by director name/ID
	Studio           string    `json:"studio,omitempty"`           // Filter by studio
	Creator          string    `json:"creator,omitempty"`          // Filter by content creator
	MediaType        MediaType `json:"mediaType,omitempty"`        // Filter by media type (movie, show, music, etc.)
	ContentRating    string    `json:"contentRating,omitempty"`    // Filter by content rating (PG, PG-13, etc.)
	Tags             []string  `json:"tags,omitempty"`             // Filter by tags
	RecentlyAdded    bool      `json:"recentlyAdded,omitempty"`    // Filter to recently added items
	RecentlyPlayed   bool      `json:"recentlyPlayed,omitempty"`   // Filter to recently played items
	Unwatched        bool      `json:"unwatched,omitempty"`        // Filter to unwatched items
	DateAddedAfter   time.Time `json:"dateAddedAfter,omitempty"`   // Filter by date added after
	DateAddedBefore  time.Time `json:"dateAddedBefore,omitempty"`  // Filter by date added before
	ReleasedAfter    time.Time `json:"releasedAfter,omitempty"`    // Filter by release date after
	ReleasedBefore   time.Time `json:"releasedBefore,omitempty"`   // Filter by release date before
	PlayedAfter      time.Time `json:"playedAfter,omitempty"`      // Filter by played date after
	PlayedBefore     time.Time `json:"playedBefore,omitempty"`     // Filter by played date before
	MinimumRating    float32   `json:"minimumRating,omitempty"`    // Filter by minimum rating
	OwnerID          uint64    `json:"ownerId,omitempty"`          // Filter by owner ID
	ItemIDs          string    `json:"itemIds,omitempty"`          // Filter by external ID (emby, jellyfin, plex, etc.)
	ExternalSourceID string    `json:"externalSourceID,omitempty"` // Filter by external source ID (TMDB, IMDB, etc.)
}

// HasFilter checks if a specific filter is set
func (opts *QueryOptions) HasFilter(filterName string) bool {
	switch filterName {
	case "favorites":
		return opts.Favorites
	case "genre":
		return opts.Genre != ""
	case "year":
		return opts.Year > 0
	case "actor":
		return opts.Actor != ""
	case "director":
		return opts.Director != ""
	case "studio":
		return opts.Studio != ""
	case "creator":
		return opts.Creator != ""
	case "mediaType":
		return opts.MediaType != ""
	case "contentRating":
		return opts.ContentRating != ""
	case "tags":
		return len(opts.Tags) > 0
	case "recentlyAdded":
		return opts.RecentlyAdded
	case "recentlyPlayed":
		return opts.RecentlyPlayed
	case "unwatched":
		return opts.Unwatched
	case "dateAddedAfter":
		return !opts.DateAddedAfter.IsZero()
	case "dateAddedBefore":
		return !opts.DateAddedBefore.IsZero()
	case "releasedAfter":
		return !opts.ReleasedAfter.IsZero()
	case "releasedBefore":
		return !opts.ReleasedBefore.IsZero()
	case "playedAfter":
		return !opts.PlayedAfter.IsZero()
	case "playedBefore":
		return !opts.PlayedBefore.IsZero()
	case "minimumRating":
		return opts.MinimumRating > 0
	case "externalSourceID":
		return opts.ExternalSourceID != ""
	default:
		return false
	}
}

// GetFilterValue returns the string value of a filter
func (opts *QueryOptions) GetFilterValue(filterName string) string {
	switch filterName {
	case "favorites":
		if opts.Favorites {
			return "true"
		}
		return ""
	case "genre":
		return opts.Genre
	case "year":
		if opts.Year > 0 {
			return fmt.Sprintf("%d", opts.Year)
		}
		return ""
	case "actor":
		return opts.Actor
	case "director":
		return opts.Director
	case "studio":
		return opts.Studio
	case "creator":
		return opts.Creator
	case "mediaType":
		return string(opts.MediaType)
	case "contentRating":
		return opts.ContentRating
	case "tags":
		if len(opts.Tags) > 0 {
			return strings.Join(opts.Tags, ",")
		}
		return ""
	case "recentlyAdded":
		if opts.RecentlyAdded {
			return "true"
		}
		return ""
	case "recentlyPlayed":
		if opts.RecentlyPlayed {
			return "true"
		}
		return ""
	case "unwatched":
		if opts.Unwatched {
			return "true"
		}
		return ""
	case "minimumRating":
		if opts.MinimumRating > 0 {
			return fmt.Sprintf("%.1f", opts.MinimumRating)
		}
		return ""
	case "externalSourceID":
		return opts.ExternalSourceID
	default:
		return ""
	}
}
