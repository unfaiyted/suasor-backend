/*
 * Emby Server REST API
 *
 * Explore the Emby Server API
 *
 */
package embyclient

type GetDirectoryContents struct {
	Username string `json:"Username,omitempty"`
	Password string `json:"Password,omitempty"`
}
