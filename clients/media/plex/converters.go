package plex

import (
	"suasor/clients/media/types"
	"time"

	"github.com/LukeHagar/plexgo/models/operations"
)

func (c *PlexClient) createDetailsFromMetadataChildren(item *operations.GetMetadataChildrenMetadata) types.MediaDetails {
	metadata := types.MediaDetails{
		Title:       *item.Title,
		Description: *item.Summary,
		Artwork: types.Artwork{
			Thumbnail: c.makeFullURL(*item.Thumb),
		},
		ExternalIDs: types.ExternalIDs{types.ExternalID{
			Source: "plex",
			ID:     *item.RatingKey,
		}},
		UpdatedAt: time.Unix(int64(*item.UpdatedAt), 0),
		AddedAt:   time.Unix(int64(*item.AddedAt), 0),
	}

	// Add optional fields if present
	if item.UpdatedAt != nil {
		metadata.UpdatedAt = time.Unix(int64(*item.UpdatedAt), 0)
	}
	if *item.AddedAt != 0 {
		metadata.AddedAt = time.Unix(int64(*item.AddedAt), 0)
	}
	if item.ParentYear != nil {
		metadata.ReleaseYear = *item.ParentYear
	}

	return metadata
}

// createMetadataFromPlexItem creates a MediaDetails from a Plex item
func (c *PlexClient) createDetailsFromLibraryMetadata(item *operations.GetLibraryItemsMetadata) types.MediaDetails {
	metadata := types.MediaDetails{
		Title:       item.Title,
		Description: item.Summary,
		Artwork: types.Artwork{
			Thumbnail: c.makeFullURL(*item.Thumb),
		},
		ExternalIDs: types.ExternalIDs{types.ExternalID{
			Source: "plex",
			ID:     item.RatingKey,
		}},
	}

	// Add optional fields if present
	if item.UpdatedAt != nil {
		metadata.UpdatedAt = time.Unix(int64(*item.UpdatedAt), 0)
	}
	if item.AddedAt != 0 {
		metadata.AddedAt = time.Unix(int64(item.AddedAt), 0)
	}
	if item.Year != nil {
		metadata.ReleaseYear = *item.Year
	}
	if item.Rating != nil {
		metadata.Ratings = types.Ratings{
			types.Rating{
				Source: "plex",
				Value:  float32(*item.Rating),
			},
		}
	}
	if item.Duration != nil {
		duration := time.Duration(*item.Duration) * time.Millisecond
		metadata.Duration = int64(duration.Seconds())
	}
	if item.Studio != nil {
		metadata.Studios = []string{*item.Studio}
	}
	if item.ContentRating != nil {
		metadata.ContentRating = *item.ContentRating
	}

	// Add genres if present
	if item.Genre != nil {
		metadata.Genres = make([]string, 0, len(item.Genre))
		for _, genre := range item.Genre {
			if genre.Tag != nil {
				metadata.Genres = append(metadata.Genres, *genre.Tag)
			}
		}
	}

	return metadata
}

// createMediaDetailsFromPlexItem creates a MediaDetails from a Plex item
func (c *PlexClient) createDetailsFromMediaMetadata(item *operations.GetMediaMetaDataMetadata) types.MediaDetails {
	metadata := types.MediaDetails{
		Title:       item.Title,
		Description: item.Summary,
		Artwork:     types.Artwork{
			// Thumbnail: c.makeFullURL(*item.Thumb),
		},
		ExternalIDs: types.ExternalIDs{types.ExternalID{
			Source: "plex",
			ID:     item.RatingKey,
		}},
	}

	// Add optional fields if present
	if item.AddedAt != 0 {
		metadata.AddedAt = time.Unix(int64(item.AddedAt), 0)
	}

	metadata.UpdatedAt = time.Unix(int64(item.UpdatedAt), 0)
	metadata.ReleaseYear = item.Year

	if item.Rating != nil {
		metadata.Ratings = types.Ratings{
			types.Rating{
				Source: "plex",
				Value:  float32(*item.Rating),
			},
		}
	}
	if item.Studio != nil {
		metadata.Studios = []string{*item.Studio}
	}
	if item.ContentRating != nil {
		metadata.ContentRating = *item.ContentRating
	}

	// Add genres if present
	if item.Genre != nil {
		metadata.Genres = make([]string, 0, len(item.Genre))
		for _, genre := range item.Genre {
			metadata.Genres = append(metadata.Genres, genre.Tag)
		}
	}

	return metadata
}
