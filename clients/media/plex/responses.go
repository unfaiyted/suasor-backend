package plex

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"context"
	"github.com/unfaiyted/plexgo/models/operations"
	"io"
	"suasor/utils/logger"
)

// ParsePlexUsersResponse parses a GetUsersResponse into a structured PlexResponse
func ParsePlexUsersResponse(ctx context.Context, response *operations.GetUsersResponse) (*PlexResponse, error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Str("responseBody", string(response.Body)).
		Msg("Parsing Plex response")
	// If Body is already populated, use it
	if len(response.Body) > 0 {
		var plexResponse PlexResponse
		err := xml.Unmarshal(response.Body, &plexResponse)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal Plex response: %w", err)
		}
		return &plexResponse, nil
	}

	// If Body is not populated but RawResponse is available, read from it
	if response.RawResponse != nil && response.RawResponse.Body != nil {
		// Create a new reader from the RawResponse.Body to avoid closing it
		// (in case the caller needs it later)
		bodyBytes, err := io.ReadAll(response.RawResponse.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		var plexResponse PlexResponse
		err = json.Unmarshal(bodyBytes, &plexResponse)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal Plex response: %w", err)
		}
		return &plexResponse, nil
	}

	return nil, fmt.Errorf("no response body available to parse")
}

// ParsePlexSearchResponse parses a raw HTTP response from the Plex search API
func ParsePlexSearchResponse(ctx context.Context, rawResponse *http.Response) (*PlexSearchResponse, error) {
	log := logger.LoggerFromContext(ctx)

	if rawResponse == nil || rawResponse.Body == nil {
		return nil, fmt.Errorf("no response body available to parse")
	}

	// Read the body
	bodyBytes, err := io.ReadAll(rawResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Log the raw response for debugging
	bodyString := string(bodyBytes)
	log.Debug().
		Str("responseBody", bodyString).
		Msg("Parsing Plex search response")

	// Create a new SearchResponse
	var searchResponse PlexSearchResponse

	// Determine content type
	contentType := rawResponse.Header.Get("Content-Type")

	if strings.Contains(contentType, "xml") || strings.HasPrefix(bodyString, "<?xml") || strings.HasPrefix(bodyString, "<MediaContainer") {
		// Parse as XML
		err = xml.Unmarshal(bodyBytes, &searchResponse)
		if err != nil {
			log.Warn().
				Err(err).
				Msg("Failed to unmarshal Plex search response as XML, trying JSON")

			// If XML fails, try JSON as a fallback
			err = json.Unmarshal(bodyBytes, &searchResponse)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal Plex search response as JSON after XML failed: %w", err)
			}
		}
	} else {
		// Assume JSON
		err = json.Unmarshal(bodyBytes, &searchResponse)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to unmarshal Plex search response as JSON")

			// If we have detailed error information, log it
			if jsonErr, ok := err.(*json.UnmarshalTypeError); ok {
				log.Error().
					Str("expected", jsonErr.Type.String()).
					Str("got", jsonErr.Value).
					Str("field", jsonErr.Field).
					Int64("offset", jsonErr.Offset).
					Msg("JSON type error details")
			}

			return nil, fmt.Errorf("failed to unmarshal Plex search response: %w", err)
		}
	}

	return &searchResponse, nil
}

// PlexResponse represents the top-level response from Plex
type PlexResponse struct {
	MediaContainer MediaContainer `xml:"MediaContainer" json:"MediaContainer"`
}

// MediaContainer represents the container holding Plex users and metadata
type MediaContainer struct {
	FriendlyName      string `xml:"friendlyName,attr" json:"friendlyName"`
	Identifier        string `xml:"identifier,attr" json:"identifier"`
	MachineIdentifier string `xml:"machineIdentifier,attr" json:"machineIdentifier"`
	TotalSize         int    `xml:"totalSize,attr" json:"totalSize"`
	Size              int    `xml:"size,attr" json:"size"`
	Users             []User `xml:"User" json:"User"`
}

// User represents a Plex user account
type User struct {
	ID                        int64    `xml:"id,attr" json:"id"`
	Title                     string   `xml:"title,attr" json:"title"`
	Username                  string   `xml:"username,attr" json:"username"`
	Email                     string   `xml:"email,attr" json:"email"`
	RecommendationsPlaylistID string   `xml:"recommendationsPlaylistId,attr" json:"recommendationsPlaylistId"`
	Thumb                     string   `xml:"thumb,attr" json:"thumb"`
	Protected                 int      `xml:"protected,attr" json:"protected"`
	Home                      int      `xml:"home,attr" json:"home"`
	AllowTuners               int      `xml:"allowTuners,attr" json:"allowTuners"`
	AllowSync                 int      `xml:"allowSync,attr" json:"allowSync"`
	AllowCameraUpload         int      `xml:"allowCameraUpload,attr" json:"allowCameraUpload"`
	AllowChannels             int      `xml:"allowChannels,attr" json:"allowChannels"`
	AllowSubtitleAdmin        int      `xml:"allowSubtitleAdmin,attr" json:"allowSubtitleAdmin"`
	FilterAll                 string   `xml:"filterAll,attr" json:"filterAll"`
	FilterMovies              string   `xml:"filterMovies,attr" json:"filterMovies"`
	FilterMusic               string   `xml:"filterMusic,attr" json:"filterMusic"`
	FilterPhotos              string   `xml:"filterPhotos,attr" json:"filterPhotos"`
	FilterTelevision          string   `xml:"filterTelevision,attr" json:"filterTelevision"`
	Restricted                int      `xml:"restricted,attr" json:"restricted"`
	Servers                   []Server `xml:"Server" json:"Server"`
}

// Server represents a Plex server associated with a user
type Server struct {
	ID                int    `xml:"id,attr" json:"id"`
	ServerID          int    `xml:"serverId,attr" json:"serverId"`
	MachineIdentifier string `xml:"machineIdentifier,attr" json:"machineIdentifier"`
	Name              string `xml:"name,attr" json:"name"`
	LastSeenAt        int64  `xml:"lastSeenAt,attr" json:"lastSeenAt"`
	NumLibraries      int    `xml:"numLibraries,attr" json:"numLibraries"`
	AllLibraries      int    `xml:"allLibraries,attr" json:"allLibraries"`
	Owned             int    `xml:"owned,attr" json:"owned"`
	Pending           int    `xml:"pending,attr" json:"pending"`
}

// PlexSearchResponse represents the top-level response from Plex search API
type PlexSearchResponse struct {
	MediaContainer SearchMediaContainer `xml:"MediaContainer" json:"MediaContainer"`
}

// SearchMediaContainer represents the container holding Plex search results
type SearchMediaContainer struct {
	Size       int         `xml:"size,attr" json:"size"`
	TotalSize  int         `xml:"totalSize,attr" json:"totalSize"`
	Identifier string      `xml:"identifier,attr" json:"identifier"`
	Hub        []SearchHub `xml:"Hub" json:"Hub"`
}

// SearchHub represents a category of search results (movies, shows, etc.)
type SearchHub struct {
	Title    string           `xml:"title,attr" json:"title"`
	Type     string           `xml:"type,attr" json:"type"`
	HubKey   string           `xml:"hubKey,attr" json:"hubKey"`
	Size     int              `xml:"size,attr" json:"size"`
	Metadata []SearchMetadata `xml:"Metadata" json:"Metadata"`
}

// SearchMetadata represents a single item in the search results
type SearchMetadata struct {
	RatingKey             string      `xml:"ratingKey,attr" json:"ratingKey"`
	Key                   string      `xml:"key,attr" json:"key"`
	GUID                  string      `xml:"guid,attr" json:"guid"`
	Type                  string      `xml:"type,attr" json:"type"`
	Title                 string      `xml:"title,attr" json:"title"`
	LibrarySectionTitle   string      `xml:"librarySectionTitle,attr" json:"librarySectionTitle"`
	LibrarySectionID      int         `xml:"librarySectionID,attr" json:"librarySectionID"`
	LibrarySectionKey     string      `xml:"librarySectionKey,attr" json:"librarySectionKey"`
	ContentRating         string      `xml:"contentRating,attr" json:"contentRating"`
	Summary               string      `xml:"summary,attr" json:"summary"`
	Rating                interface{} `xml:"rating,attr" json:"rating"`
	Year                  interface{} `xml:"year,attr" json:"year"`
	Thumb                 string      `xml:"thumb,attr" json:"thumb"`
	Art                   string      `xml:"art,attr" json:"art"`
	Duration              interface{} `xml:"duration,attr" json:"duration"`
	OriginallyAvailableAt string      `xml:"originallyAvailableAt,attr" json:"originallyAvailableAt"`
	AddedAt               interface{} `xml:"addedAt,attr" json:"addedAt"`
	UpdatedAt             interface{} `xml:"updatedAt,attr" json:"updatedAt"`

	// Fields specific to TV shows
	ParentRatingKey string `xml:"parentRatingKey,attr" json:"parentRatingKey"`
	ParentKey       string `xml:"parentKey,attr" json:"parentKey"`
	ParentTitle     string `xml:"parentTitle,attr" json:"parentTitle"`

	// Fields specific to episodes & tracks
	GrandparentRatingKey string `xml:"grandparentRatingKey,attr" json:"grandparentRatingKey"`
	GrandparentKey       string `xml:"grandparentKey,attr" json:"grandparentKey"`
	GrandparentTitle     string `xml:"grandparentTitle,attr" json:"grandparentTitle"`

	// Additional fields from the response
	Score          interface{} `xml:"score,attr" json:"score"`
	TagLine        string      `xml:"tagline,attr" json:"tagline"`
	Slug           string      `xml:"slug,attr" json:"slug"`
	Studio         string      `xml:"studio,attr" json:"studio"`
	AudienceRating interface{} `xml:"audienceRating,attr" json:"audienceRating"`
	ViewCount      interface{} `xml:"viewCount,attr" json:"viewCount"`
	LastViewedAt   interface{} `xml:"lastViewedAt,attr" json:"lastViewedAt"`
}

// GetYear returns the year as an integer
func (s SearchMetadata) GetYear() int {
	if s.Year == nil {
		return 0
	}

	// Handle different types
	switch v := s.Year.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case string:
		year, err := strconv.Atoi(v)
		if err != nil {
			return 0
		}
		return year
	default:
		// Try converting to string
		yearStr := fmt.Sprintf("%v", s.Year)
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			return 0
		}
		return year
	}
}
