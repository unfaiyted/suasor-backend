/*
 * Emby Server REST API
 *
 * Explore the Emby Server API
 *
 */
package embyclient

type PlaylistsPlaylistCreationResult struct {
	Id string `json:"Id,omitempty"`
	Name string `json:"Name,omitempty"`
	ItemAddedCount int32 `json:"ItemAddedCount,omitempty"`
}
