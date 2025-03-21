/*
 * Emby Server REST API
 *
 * Explore the Emby Server API
 *
 */
package embyclient

// Class ParentalRating  
type ParentalRating struct {
	// The name.
	Name string `json:"Name,omitempty"`
	// The value.
	Value int32 `json:"Value,omitempty"`
}
