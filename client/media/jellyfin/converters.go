package jellyfin

import (
	"context"
	"fmt"
	"time"

	jellyfin "github.com/sj14/jellyfin-go/api"
	t "suasor/client/media/types"
	"suasor/utils"
)

// Helper function to convert Jellyfin item to internal Collection type
func (j *JellyfinClient) convertToCollection(ctx context.Context, item *jellyfin.BaseItemDto) (t.MediaItem[t.Collection], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	// Validate required fields
	if item == nil {
		return t.MediaItem[t.Collection]{}, fmt.Errorf("cannot convert nil item to collection")
	}

	if item.Id == nil || *item.Id == "" {
		return t.MediaItem[t.Collection]{}, fmt.Errorf("collection is missing required ID field")
	}

	// Safely get name or fallback to empty string
	title := ""
	if item.Name.IsSet() {
		title = *item.Name.Get()
	}

	log.Debug().
		Str("collectionID", *item.Id).
		Str("collectionName", title).
		Msg("Converting Jellyfin item to collection format")

	// Safely handle optional fields
	description := ""
	if item.Overview.IsSet() {
		description = *item.Overview.Get()
	}

	// Safely handle item count
	itemCount := 0
	if item.ChildCount.IsSet() {
		itemCount = int(*item.ChildCount.Get())
	}

	// Build collection object
	collection := t.MediaItem[t.Collection]{
		Data: t.Collection{
			Details: t.MediaMetadata{
				Title:       title,
				Description: description,
				Artwork:     j.getArtworkURLs(item),
			},
			ItemCount: itemCount,
		},
		Type: "collection",
	}
	collection.SetClientInfo(j.ClientID, j.ClientType, *item.Id)

	// Add potential year if available
	if item.ProductionYear.IsSet() {
		collection.Data.Details.ReleaseYear = int(*item.ProductionYear.Get())
	}

	// Add community rating if available
	if item.CommunityRating.IsSet() {
		collection.Data.Details.Ratings = append(collection.Data.Details.Ratings, t.Rating{
			Source: "jellyfin",
			Value:  float32(*item.CommunityRating.Get()),
		})
	}

	// Add user rating if available
	if item.UserData.IsSet() && item.UserData.Get().Rating.IsSet() {
		collection.Data.Details.UserRating = float32(*item.UserData.Get().Rating.Get())
	}

	// Handle genres if available
	if item.Genres != nil {
		collection.Data.Details.Genres = item.Genres
	}

	// Extract provider IDs if available
	extractProviderIDs(&item.ProviderIds, &collection.Data.Details.ExternalIDs)

	log.Debug().
		Str("collectionID", *item.Id).
		Str("collectionName", collection.Data.Details.Title).
		Int("itemCount", collection.Data.ItemCount).
		Msg("Successfully converted Jellyfin item to collection")

	return collection, nil
}

// Helper function to convert Jellyfin item to internal Episode type
func (j *JellyfinClient) convertToEpisode(ctx context.Context, item *jellyfin.BaseItemDto) (t.MediaItem[t.Episode], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	// Validate required fields
	if item == nil {
		return t.MediaItem[t.Episode]{}, fmt.Errorf("cannot convert nil item to episode")
	}

	if item.Id == nil || *item.Id == "" {
		return t.MediaItem[t.Episode]{}, fmt.Errorf("episode is missing required ID field")
	}

	// Safely get name or fallback to empty string
	title := ""
	if item.Name.IsSet() {
		title = *item.Name.Get()
	}

	log.Debug().
		Str("episodeID", *item.Id).
		Str("episodeName", title).
		Msg("Converting Jellyfin item to episode format")

	// Safely handle optional fields
	description := ""
	if item.Overview.IsSet() {
		description = *item.Overview.Get()
	}

	// Safely handle duration
	var duration time.Duration
	if item.RunTimeTicks.IsSet() {
		duration = time.Duration(*item.RunTimeTicks.Get()/10000000) * time.Second
	}

	// Safely handle episode number
	var episodeNumber int64
	if item.IndexNumber.IsSet() {
		episodeNumber = int64(*item.IndexNumber.Get())
	}

	// Safely handle season number
	seasonNumber := 0
	if item.ParentIndexNumber.IsSet() {
		seasonNumber = int(*item.ParentIndexNumber.Get())
	}

	// Safely handle show title
	showTitle := ""
	if item.SeriesName.IsSet() {
		showTitle = *item.SeriesName.Get()
	}

	// Create the basic episode object
	episode := t.MediaItem[t.Episode]{
		Data: t.Episode{
			Details: t.MediaMetadata{
				Title:       title,
				Description: description,
				Artwork:     j.getArtworkURLs(item),
				Duration:    duration,
			},
			Number:       episodeNumber,
			SeasonNumber: seasonNumber,
			ShowTitle:    showTitle,
		},
		Type: "episode",
	}

	episode.SetClientInfo(j.ClientID, j.ClientType, *item.Id)

	// Safely set IDs if available
	if item.SeriesId.IsSet() {
		episode.Data.ShowID = *item.SeriesId.Get()
	}

	if item.SeasonId.IsSet() {
		episode.Data.SeasonID = *item.SeasonId.Get()
	}

	// Add air date if available
	if item.PremiereDate.IsSet() {
		episode.Data.Details.ReleaseDate = *item.PremiereDate.Get()
	}

	// Add community rating if available
	if item.CommunityRating.IsSet() {
		episode.Data.Details.Ratings = append(episode.Data.Details.Ratings, t.Rating{
			Source: "jellyfin",
			Value:  float32(*item.CommunityRating.Get()),
		})
	}

	// Add user rating if available
	if item.UserData.IsSet() && item.UserData.Get().Rating.IsSet() {
		episode.Data.Details.UserRating = float32(*item.UserData.Get().Rating.Get())
	}

	// Extract provider IDs if available
	extractProviderIDs(&item.ProviderIds, &episode.Data.Details.ExternalIDs)

	log.Debug().
		Str("episodeID", *item.Id).
		Str("episodeName", episode.Data.Details.Title).
		Int64("episodeNumber", episode.Data.Number).
		Int("seasonNumber", episode.Data.SeasonNumber).
		Msg("Successfully converted Jellyfin item to episode")

	return episode, nil
}

func (j *JellyfinClient) convertToTVShow(ctx context.Context, item *jellyfin.BaseItemDto) (t.MediaItem[t.TVShow], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	// Validate required fields
	if item == nil {
		return t.MediaItem[t.TVShow]{}, fmt.Errorf("cannot convert nil item to TV show")
	}

	if item.Id == nil || *item.Id == "" {
		return t.MediaItem[t.TVShow]{}, fmt.Errorf("TV show is missing required ID field")
	}

	// Safely get name or fallback to empty string
	title := ""
	if item.Name.IsSet() {
		title = *item.Name.Get()
	}

	log.Debug().
		Str("showID", *item.Id).
		Str("showName", title).
		Msg("Converting Jellyfin item to TV show format")

	// Safely handle optional fields
	description := ""
	if item.Overview.IsSet() {
		description = *item.Overview.Get()
	}

	// Default values
	releaseYear := 0
	if item.ProductionYear.IsSet() {
		releaseYear = int(*item.ProductionYear.Get())
	}

	// Safely handle genres
	var genres []string
	if item.Genres != nil {
		genres = item.Genres
	}

	// Safely handle duration
	var duration time.Duration
	if item.RunTimeTicks.IsSet() {
		duration = time.Duration(*item.RunTimeTicks.Get()/10000000) * time.Second
	}

	// Safely handle season count
	seasonCount := 0
	if item.ChildCount.IsSet() {
		seasonCount = int(*item.ChildCount.Get())
	}

	// Safely handle status
	status := ""
	if item.Status.IsSet() {
		status = *item.Status.Get()
	}

	// Build TV show object
	show := t.MediaItem[t.TVShow]{
		Data: t.TVShow{
			Details: t.MediaMetadata{
				Title:       title,
				Description: description,
				ReleaseYear: releaseYear,
				Genres:      genres,
				Artwork:     j.getArtworkURLs(item),
				Duration:    duration,
			},
			Status:      status,
			SeasonCount: seasonCount,
		},
		Type: "tvshow",
	}

	// ClientID:   j.ClientID,
	// 			ExternalID: *item.Id,
	// 			ClientType: string(j.ClientType),
	show.SetClientInfo(j.ClientID, j.ClientType, *item.Id)

	// Set SeriesStudio if available
	if item.SeriesStudio.IsSet() {
		show.Data.Network = *item.SeriesStudio.Get()
	}

	// Extract provider IDs if available
	extractProviderIDs(&item.ProviderIds, &show.Data.Details.ExternalIDs)

	// Set ratings if available
	if item.CommunityRating.IsSet() {
		show.Data.Details.Ratings = append(show.Data.Details.Ratings, t.Rating{
			Source: "jellyfin",
			Value:  float32(*item.CommunityRating.Get()),
		})
	}

	// Set user rating if available
	if item.UserData.IsSet() && item.UserData.Get().Rating.IsSet() {
		show.Data.Details.UserRating = float32(*item.UserData.Get().Rating.Get())
	}

	log.Debug().
		Str("showID", *item.Id).
		Str("showName", show.Data.Details.Title).
		Int("seasonCount", show.Data.SeasonCount).
		Msg("Successfully converted Jellyfin item to TV show")

	return show, nil
}

// Helper function to convert Jellyfin item to internal Movie type
func (j *JellyfinClient) convertToMovie(ctx context.Context, item *jellyfin.BaseItemDto) (t.MediaItem[t.Movie], error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	// Validate required fields
	if item == nil {
		return t.MediaItem[t.Movie]{}, fmt.Errorf("cannot convert nil item to movie")
	}

	if item.Id == nil || *item.Id == "" {
		return t.MediaItem[t.Movie]{}, fmt.Errorf("movie is missing required ID field")
	}

	// Safely get name or fallback to empty string
	title := ""
	if item.Name.IsSet() {
		title = *item.Name.Get()
	}

	log.Debug().
		Str("movieID", *item.Id).
		Str("movieName", title).
		Msg("Converting Jellyfin item to movie format")

	// Safely handle optional fields
	description := ""
	if item.Overview.IsSet() {
		description = *item.Overview.Get()
	}

	contentRating := ""
	if item.OfficialRating.IsSet() {
		contentRating = *item.OfficialRating.Get()
	}

	// Determine release year from either ProductionYear or PremiereDate
	var releaseYear int
	var releaseDate time.Time

	if item.ProductionYear.IsSet() {
		releaseYear = int(*item.ProductionYear.Get())
	}

	if item.PremiereDate.IsSet() {
		releaseDate = *item.PremiereDate.Get()
		if releaseYear == 0 {
			releaseYear = releaseDate.Year()
			log.Debug().
				Str("movieID", *item.Id).
				Str("premiereDate", releaseDate.Format("2006-01-02")).
				Int("extractedYear", releaseYear).
				Msg("Using year from premiere date instead of production year")
		}
	}

	// Extract genres
	var genres []string
	if item.Genres != nil {
		genres = item.Genres
	}

	// Calculate duration
	var duration time.Duration
	if item.RunTimeTicks.IsSet() {
		duration = time.Duration(*item.RunTimeTicks.Get()/10000000) * time.Second
	}

	// Initialize ratings
	ratings := t.Ratings{}

	// Safely add community rating if available
	if item.CommunityRating.IsSet() {
		ratings = append(ratings, t.Rating{
			Source: "jellyfin",
			Value:  float32(*item.CommunityRating.Get()),
		})
	}

	// Build movie object
	movie := t.MediaItem[t.Movie]{
		Data: t.Movie{
			Details: t.MediaMetadata{
				Title:         title,
				Description:   description,
				ReleaseDate:   releaseDate,
				ReleaseYear:   releaseYear,
				ContentRating: contentRating,
				Genres:        genres,
				Artwork:       j.getArtworkURLs(item),
				Duration:      duration,
				Ratings:       ratings,
			},
		},
		Type: t.MEDIATYPE_MOVIE,
	}

	movie.SetClientInfo(j.ClientID, j.ClientType, *item.Id)

	// Set user rating if available
	if item.UserData.IsSet() && item.UserData.Get().Rating.IsSet() {
		movie.Data.Details.UserRating = float32(*item.UserData.Get().Rating.Get())
	} else {
		log.Debug().
			Str("movieID", *item.Id).
			Msg("Movie has no user data, skipping user rating")
	}

	// Extract provider IDs if available
	extractProviderIDs(&item.ProviderIds, &movie.Data.Details.ExternalIDs)

	log.Debug().
		Str("movieID", *item.Id).
		Str("movieTitle", movie.Data.Details.Title).
		Int("year", movie.Data.Details.ReleaseYear).
		Msg("Successfully converted Jellyfin item to movie")

	return movie, nil
}
