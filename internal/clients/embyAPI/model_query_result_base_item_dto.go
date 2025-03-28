/*
 * Emby Server REST API
 *
 * Explore the Emby Server API
 *
 */
package embyclient

type QueryResultBaseItemDto struct {
	Items []BaseItemDto `json:"Items,omitempty"`
	TotalRecordCount int32 `json:"TotalRecordCount,omitempty"`
}
