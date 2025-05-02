package plex

import (
	"encoding/json"
	"fmt"

	"github.com/LukeHagar/plexgo/models/operations"
	"io"
)

// ParsePlexUsersResponse parses a GetUsersResponse into a structured PlexResponse
func ParsePlexUsersResponse(response *operations.GetUsersResponse) (*PlexResponse, error) {
	// If Body is already populated, use it
	if len(response.Body) > 0 {
		var plexResponse PlexResponse
		err := json.Unmarshal(response.Body, &plexResponse)
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
	MediaContainer MediaContainer `json:"MediaContainer"`
}

// MediaContainer represents the container holding Plex users and metadata
type MediaContainer struct {
	FriendlyName      string `json:"friendlyName"`
	Identifier        string `json:"identifier"`
	MachineIdentifier string `json:"machineIdentifier"`
	TotalSize         int    `json:"totalSize"`
	Size              int    `json:"size"`
	Users             []User `json:"User"`
}

// User represents a Plex user account
type User struct {
	ID                        int64    `json:"id"`
	Title                     string   `json:"title"`
	Username                  string   `json:"username"`
	Email                     string   `json:"email"`
	RecommendationsPlaylistID string   `json:"recommendationsPlaylistId"`
	Thumb                     string   `json:"thumb"`
	Protected                 int      `json:"protected"`
	Home                      int      `json:"home"`
	AllowTuners               int      `json:"allowTuners"`
	AllowSync                 int      `json:"allowSync"`
	AllowCameraUpload         int      `json:"allowCameraUpload"`
	AllowChannels             int      `json:"allowChannels"`
	AllowSubtitleAdmin        int      `json:"allowSubtitleAdmin"`
	FilterAll                 string   `json:"filterAll"`
	FilterMovies              string   `json:"filterMovies"`
	FilterMusic               string   `json:"filterMusic"`
	FilterPhotos              string   `json:"filterPhotos"`
	FilterTelevision          string   `json:"filterTelevision"`
	Restricted                int      `json:"restricted"`
	Servers                   []Server `json:"Server"`
}

// Server represents a Plex server associated with a user
type Server struct {
	ID                int    `json:"id"`
	ServerID          int    `json:"serverId"`
	MachineIdentifier string `json:"machineIdentifier"`
	Name              string `json:"name"`
	LastSeenAt        int64  `json:"lastSeenAt"`
	NumLibraries      int    `json:"numLibraries"`
	AllLibraries      int    `json:"allLibraries"`
	Owned             int    `json:"owned"`
	Pending           int    `json:"pending"`
}
