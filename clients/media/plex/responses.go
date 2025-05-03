package plex

import (
	"encoding/json"
	"encoding/xml"
	"fmt"

	"context"
	"github.com/LukeHagar/plexgo/models/operations"
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
