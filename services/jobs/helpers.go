package jobs

//
// import (
// 	"slices"
// 	mediatypes "suasor/clients/media/types"
// 	"suasor/types/models"
// )
//
// func MergeMediaItem[T mediatypes.MediaData](existingItem *models.MediaItem[T], newItem *models.MediaItem[T]) *models.MediaItem[T] {
//
// 	// Merge sync clients
// 	existingItem.SyncClients.Merge(newItem.SyncClients)
// 	existingItem.ExternalIDs.Merge(newItem.ExternalIDs)
//
// 	existingDetails := existingItem.Data.GetDetails()
// 	newDetails := newItem.Data.GetDetails()
// 	existingDetails.ExternalIDs = existingItem.ExternalIDs
//
// 	// Update data fields
// 	if existingDetails.Title == "" {
// 		existingDetails.Title = newDetails.Title
// 	}
// 	if existingDetails.Description == "" {
// 		existingDetails.Description = newDetails.Description
// 	}
// 	if existingDetails.ContentRating == "" {
// 		existingDetails.ContentRating = newDetails.ContentRating
// 	}
// 	if existingDetails.ContentRating == "" {
// 		existingDetails.ContentRating = newDetails.ContentRating
// 	}
// 	if existingDetails.Studio == "" {
// 		existingDetails.Studio = newDetails.Studio
// 	}
//
// 	existingDetails.Genres = mergeStringArray(existingDetails.Genres, newDetails.Genres)
// 	existingDetails.Ratings = mergeRatings(existingDetails.Ratings, newDetails.Ratings)
//
// 	// Artworks
// 	if existingDetails.Artwork.Poster == "" {
// 		existingDetails.Artwork.Poster = newDetails.Artwork.Poster
// 	}
// 	if existingDetails.Artwork.Banner == "" {
// 		existingDetails.Artwork.Banner = newDetails.Artwork.Banner
// 	}
// 	if existingDetails.Artwork.Thumbnail == "" {
// 		existingDetails.Artwork.Thumbnail = newDetails.Artwork.Thumbnail
// 	}
// 	if existingDetails.Artwork.Logo == "" {
// 		existingDetails.Artwork.Logo = newDetails.Artwork.Logo
// 	}
//
// 	if existingDetails.ReleaseYear == 0 {
// 		existingDetails.ReleaseYear = newDetails.ReleaseYear
// 	}
// 	if existingDetails.ReleaseDate.IsZero() {
// 		existingDetails.ReleaseDate = newDetails.ReleaseDate
// 	}
//
// 	if existingItem.Title == "" {
// 		existingItem.Title = newDetails.Title
// 	}
// 	if existingItem.ReleaseYear == 0 {
// 		existingItem.ReleaseYear = newDetails.ReleaseYear
// 	}
// 	if existingItem.ReleaseDate.IsZero() {
// 		existingItem.ReleaseDate = newDetails.ReleaseDate
// 	}
//
// 	return existingItem
// }
//
// // Helpers
// func mergeStringArray(genres []string, newGenres []string) []string {
// 	for _, newGenre := range newGenres {
// 		if !slices.Contains(genres, newGenre) {
// 			genres = append(genres, newGenre)
// 		}
// 	}
// 	return genres
// }
//
// func mergeRatings(ratings mediatypes.Ratings, newRatings mediatypes.Ratings) mediatypes.Ratings {
// 	for _, newRating := range newRatings {
// 		if !slices.Contains(ratings, newRating) {
// 			ratings = append(ratings, newRating)
// 		}
// 	}
// 	return ratings
// }
