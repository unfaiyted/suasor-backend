/*
 * Emby Server REST API
 *
 * Explore the Emby Server API
 *
 */
package embyclient

type SyncTarget struct {
	// The name.
	Name string `json:"Name,omitempty"`
	// The identifier.
	Id string `json:"Id,omitempty"`
}
