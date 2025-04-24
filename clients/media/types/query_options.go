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
	Favorites       bool      `json:"favorites,omitempty"`       // Filter to favorites only
	Genre           string    `json:"genre,omitempty"`           // Filter by genre
	Year            int       `json:"year,omitempty"`            // Filter by release year
	Actor           string    `json:"actor,omitempty"`           // Filter by actor name/ID
	Director        string    `json:"director,omitempty"`        // Filter by director name/ID
	Studio          string    `json:"studio,omitempty"`          // Filter by studio
	Creator         string    `json:"creator,omitempty"`         // Filter by content creator
	Role            string    `json:"role,omitempty"`            // Filter by role (actor, director, etc.)
	MediaType       MediaType `json:"mediaType,omitempty"`       // Filter by media type (movie, show, music, etc.)
	ContentRating   string    `json:"contentRating,omitempty"`   // Filter by content rating (PG, PG-13, etc.)
	Tags            []string  `json:"tags,omitempty"`            // Filter by tags
	RecentlyAdded   bool      `json:"recentlyAdded,omitempty"`   // Filter to recently added items
	RecentlyPlayed  bool      `json:"recentlyPlayed,omitempty"`  // Filter to recently played items
	Watched         bool      `json:"watched,omitempty"`         // Filter to watched items
	Watchlist       bool      `json:"watchlist,omitempty"`       // Filter to watchlist items
	DateAddedAfter  time.Time `json:"dateAddedAfter,omitempty"`  // Filter by date added after
	DateAddedBefore time.Time `json:"dateAddedBefore,omitempty"` // Filter by date added before
	ReleasedAfter   time.Time `json:"releasedAfter,omitempty"`   // Filter by release date after
	ReleasedBefore  time.Time `json:"releasedBefore,omitempty"`  // Filter by release date before
	PlayedAfter     time.Time `json:"playedAfter,omitempty"`     // Filter by played date after
	PlayedBefore    time.Time `json:"playedBefore,omitempty"`    // Filter by played date before
	MinimumRating   float32   `json:"minimumRating,omitempty"`   // Filter by minimum rating
	MaximumRating   float32   `json:"maximumRating,omitempty"`   // Filter by maximum rating (10 is the highest)

	// TODO: Add normalized rating logic . Get ratings from all clients and external sources (tmdb,imdb)
	// going to scale all of them to a range of 0-100 and then take the average of these ratings.
	// MinimumNormalizedRating float32 `json:"minimumNormalizedRating,omitempty"` // Filter by minimum normalized rating
	IsPublic         bool   `json:"isPublic,omitempty"`         // Filter by public status
	OwnerID          uint64 `json:"ownerID,omitempty"`          // Filter by owner ID
	ClientID         uint64 `json:"clientID,omitempty"`         // Filter by client ID
	PersonID         uint64 `json:"personID,omitempty"`         // Filter by person ID
	ItemIDs          string `json:"itemIDs,omitempty"`          // Filter by external ID (emby, jellyfin, plex, etc.)
	ExternalSourceID string `json:"externalSourceID,omitempty"` // Filter by external source ID (TMDB, IMDB, etc.)
}

func (opts *QueryOptions) WithOwnerID(ownerID uint64) *QueryOptions {
	opts.OwnerID = ownerID
	return opts
}

func (opts *QueryOptions) WithClientID(clientID uint64) *QueryOptions {
	opts.ClientID = clientID
	return opts
}

func (opts *QueryOptions) WithPersonID(personID uint64) *QueryOptions {
	opts.PersonID = personID
	return opts
}

func (opts *QueryOptions) WithItemIDs(itemIDs string) *QueryOptions {
	opts.ItemIDs = itemIDs
	return opts
}

func (opts *QueryOptions) WithExternalSourceID(externalSourceID string) *QueryOptions {
	opts.ExternalSourceID = externalSourceID
	return opts
}

func (opts *QueryOptions) WithMediaType(mediaType MediaType) *QueryOptions {
	opts.MediaType = mediaType
	return opts
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
	case "watched":
		return opts.Watched
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
	case "watched":
		if opts.Watched {
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
