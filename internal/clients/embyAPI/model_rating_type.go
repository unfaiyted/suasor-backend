/*
 * Emby Server REST API
 *
 * Explore the Emby Server API
 *
 */
package embyclient

type RatingType string

// List of RatingType
const (
	SCORE_RatingType RatingType = "Score"
	LIKES_RatingType RatingType = "Likes"
)
