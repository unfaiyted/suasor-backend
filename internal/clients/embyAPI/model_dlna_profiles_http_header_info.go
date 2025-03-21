/*
 * Emby Server REST API
 *
 * Explore the Emby Server API
 *
 */
package embyclient

type DlnaProfilesHttpHeaderInfo struct {
	Name string `json:"Name,omitempty"`
	Value string `json:"Value,omitempty"`
	Match *DlnaProfilesHeaderMatchType `json:"Match,omitempty"`
}
