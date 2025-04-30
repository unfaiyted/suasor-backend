package models

import (
	"slices"
	mediatypes "suasor/clients/media/types"
)

// Helpers
func mergeStringArray(genres []string, newGenres []string) []string {
	for _, newGenre := range newGenres {
		if !slices.Contains(genres, newGenre) {
			genres = append(genres, newGenre)
		}
	}
	return genres
}

func mergeRatings(ratings mediatypes.Ratings, newRatings mediatypes.Ratings) mediatypes.Ratings {
	for _, newRating := range newRatings {
		if !slices.Contains(ratings, newRating) {
			ratings = append(ratings, newRating)
		}
	}
	return ratings
}
