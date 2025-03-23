// plex/client.go
package plex

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"suasor/client/media/interfaces"
	"suasor/models"
	"suasor/utils"
	"time"

	"github.com/LukeHagar/plexgo"
	"github.com/LukeHagar/plexgo/models/operations"
)

// init is automatically called when package is imported
func init() {
	interfaces.RegisterProvider(models.MediaClientTypePlex, NewPlexClient)
}

// PlexClient implements MediaContentProvider for Plex
type PlexClient struct {
	interfaces.BaseMediaClient
	config     models.PlexConfig
	httpClient *http.Client
	baseURL    string
	plexAPI    *plexgo.PlexAPI
}

// NewPlexClient creates a new Plex client
func NewPlexClient(ctx context.Context, clientID uint64, config any) (interfaces.MediaContentProvider, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", clientID).
		Str("clientType", string(models.MediaClientTypePlex)).
		Msg("Creating new Plex client")

	plexConfig, ok := config.(models.PlexConfig)
	if !ok {
		log.Error().
			Uint64("clientID", clientID).
			Str("clientType", string(models.MediaClientTypePlex)).
			Msg("Invalid Plex configuration")
		return nil, fmt.Errorf("invalid Plex configuration")
	}

	log.Debug().
		Uint64("clientID", clientID).
		Str("host", plexConfig.Host).
		Msg("Initializing Plex API client")

	// Initialize the Plex API client
	plexAPI := plexgo.New(
		plexgo.WithSecurity(plexConfig.Token),
		plexgo.WithServerURL(plexConfig.Host),
	)

	client := &PlexClient{
		BaseMediaClient: interfaces.BaseMediaClient{
			ClientID:   clientID,
			ClientType: models.MediaClientTypePlex,
		},
		config:  plexConfig,
		plexAPI: plexAPI,
		baseURL: plexConfig.Host,
	}

	log.Info().
		Uint64("clientID", clientID).
		Str("clientType", string(models.MediaClientTypePlex)).
		Str("host", plexConfig.Host).
		Msg("Successfully created Plex client")

	return client, nil
}

// Capability methods
func (c *PlexClient) SupportsMovies() bool      { return true }
func (c *PlexClient) SupportsTVShows() bool     { return true }
func (c *PlexClient) SupportsMusic() bool       { return true }
func (c *PlexClient) SupportsPlaylists() bool   { return true }
func (c *PlexClient) SupportsCollections() bool { return true }

// Helper functions for common operations

// makeFullURL creates a complete URL from a resource path
func (c *PlexClient) makeFullURL(resourcePath string) string {
	if resourcePath == "" {
		return ""
	}

	if strings.HasPrefix(resourcePath, "http") {
		return resourcePath
	}

	return fmt.Sprintf("%s%s", c.baseURL, resourcePath)
}

// findLibrarySectionByType returns the section key for the specified type
func (c *PlexClient) findLibrarySectionByType(ctx context.Context, sectionType string) (string, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Debug().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("sectionType", sectionType).
		Msg("Finding library section by type")

	libraries, err := c.plexAPI.Library.GetAllLibraries(ctx)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("sectionType", sectionType).
			Msg("Failed to get libraries from Plex")
		return "", fmt.Errorf("failed to get libraries: %w", err)
	}

	log.Debug().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("libraryCount", len(libraries.Object.MediaContainer.Directory)).
		Msg("Retrieved libraries from Plex")

	for _, dir := range libraries.Object.MediaContainer.Directory {
		if dir.Type == sectionType {
			log.Debug().
				Uint64("clientID", c.ClientID).
				Str("clientType", string(c.ClientType)).
				Str("sectionType", sectionType).
				Str("sectionKey", dir.Key).
				Str("sectionTitle", dir.Title).
				Msg("Found matching library section")
			return dir.Key, nil
		}
	}

	log.Debug().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("sectionType", sectionType).
		Msg("No matching library section found")

	return "", nil
}

func (c *PlexClient) createChildMetadataFromPlexItem(item *operations.GetMetadataChildrenMetadata) interfaces.MediaMetadata {
	metadata := interfaces.MediaMetadata{
		Title:       *item.Title,
		Description: *item.Summary,
		Artwork: interfaces.Artwork{
			Thumbnail: c.makeFullURL(*item.Thumb),
		},
		ExternalIDs: interfaces.ExternalIDs{interfaces.ExternalID{
			Source: "plex",
			ID:     *item.RatingKey,
		}},
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

// createMetadataFromPlexItem creates a MediaMetadata from a Plex item
func (c *PlexClient) createMetadataFromPlexItem(item *operations.GetLibraryItemsMetadata) interfaces.MediaMetadata {
	metadata := interfaces.MediaMetadata{
		Title:       item.Title,
		Description: item.Summary,
		Artwork: interfaces.Artwork{
			Thumbnail: c.makeFullURL(*item.Thumb),
		},
		ExternalIDs: interfaces.ExternalIDs{interfaces.ExternalID{
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
		metadata.Ratings = interfaces.Ratings{
			interfaces.Rating{
				Source: "plex",
				Value:  float32(*item.Rating),
			},
		}
	}
	if item.Duration != nil {
		metadata.Duration = time.Duration(*item.Duration) * time.Millisecond
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

// createMediaMetadataFromPlexItem creates a MediaMetadata from a Plex item
func (c *PlexClient) createMediaMetadataFromPlexItem(item *operations.GetMediaMetaDataMetadata) interfaces.MediaMetadata {
	metadata := interfaces.MediaMetadata{
		Title:       item.Title,
		Description: item.Summary,
		Artwork:     interfaces.Artwork{
			// Thumbnail: c.makeFullURL(*item.Thumb),
		},
		ExternalIDs: interfaces.ExternalIDs{interfaces.ExternalID{
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
		metadata.Ratings = interfaces.Ratings{
			interfaces.Rating{
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

// GetPlaylists retrieves playlists from Plex
func (c *PlexClient) GetPlaylists(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.Playlist, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("baseURL", c.baseURL).
		Msg("Retrieving playlists from Plex server")

	log.Debug().Msg("Making API request to Plex server for playlists")
	res, err := c.plexAPI.Playlists.GetPlaylists(ctx, operations.PlaylistTypeAudio.ToPointer(), operations.QueryParamSmartOne.ToPointer())
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("baseURL", c.baseURL).
			Msg("Failed to get playlists from Plex")
		return nil, fmt.Errorf("failed to get playlists: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil {
		log.Info().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("No playlists found in Plex")
		return []interfaces.Playlist{}, nil
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("totalItems", len(res.Object.MediaContainer.Metadata)).
		Msg("Successfully retrieved playlists from Plex")

	playlists := make([]interfaces.Playlist, 0, len(res.Object.MediaContainer.Metadata))
	for _, item := range res.Object.MediaContainer.Metadata {
		playlist := interfaces.Playlist{
			MediaItem: interfaces.MediaItem{
				ExternalID: *item.RatingKey,
				ClientID:   c.ClientID,
				ClientType: string(c.ClientType),
				Metadata: interfaces.MediaMetadata{
					Description: *item.Summary,
					Title:       *item.Title,
					Artwork:     interfaces.Artwork{
						// Thumbnail: c.makeFullURL(*item.Thumb),
					},
					ExternalIDs: interfaces.ExternalIDs{interfaces.ExternalID{
						Source: "plex",
						ID:     *item.RatingKey,
					}},
					UpdatedAt: time.Unix(int64(*item.UpdatedAt), 0),
					AddedAt:   time.Unix(int64(*item.AddedAt), 0),
				},
			},
		}
		c.BaseMediaClient.AddClientInfo(&playlist.MediaItem)
		playlists = append(playlists, playlist)
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("playlistsReturned", len(playlists)).
		Msg("Completed GetPlaylists request")

	return playlists, nil
}

// GetCollections retrieves collections from Plex
func (c *PlexClient) GetCollections(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.Collection, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("baseURL", c.baseURL).
		Msg("Retrieving collections from Plex server")

	log.Debug().Msg("Making API request to Plex server for collections")
	res, err := c.plexAPI.Library.GetAllLibraries(ctx)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("baseURL", c.baseURL).
			Msg("Failed to get collections from Plex")
		return nil, fmt.Errorf("failed to get collections: %w", err)
	}

	directories := res.Object.MediaContainer.GetDirectory()

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("totalDirectories", len(directories)).
		Msg("Successfully retrieved library directories from Plex")

	collections := make([]interfaces.Collection, 0, len(directories))

	for _, dir := range directories {
		collection := interfaces.Collection{
			MediaItem: interfaces.MediaItem{
				ExternalID: dir.Key,
				ClientID:   c.ClientID,
				ClientType: string(c.ClientType),
				Metadata: interfaces.MediaMetadata{
					Title: dir.Title,
					Artwork: interfaces.Artwork{
						Thumbnail: c.makeFullURL(dir.Thumb),
					},
					ExternalIDs: interfaces.ExternalIDs{interfaces.ExternalID{
						Source: "plex",
						ID:     dir.Key,
					}},
				},
			},
		}

		// // Only add these fields if available
		// if dir.LeafCount != nil {
		// 	collection.ItemCount = *dir.LeafCount
		// }

		c.BaseMediaClient.AddClientInfo(&collection.MediaItem)
		collections = append(collections, collection)

		log.Debug().
			Str("collectionID", dir.Key).
			Str("collectionName", dir.Title).
			Msg("Added collection to result list")
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("collectionsReturned", len(collections)).
		Msg("Completed GetCollections request")

	return collections, nil
}

// GetMovies retrieves movies from Plex
func (c *PlexClient) GetMovies(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.Movie, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("baseURL", c.baseURL).
		Msg("Retrieving movies from Plex server")

	// First, find the movie library section
	log.Debug().Msg("Finding movie library section")
	movieSectionKey, err := c.findLibrarySectionByType(ctx, "movie")
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("Failed to find movie library section")
		return nil, err
	}

	if movieSectionKey == "" {
		log.Info().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("No movie library section found in Plex")
		return []interfaces.Movie{}, nil
	}

	// Get movies from the movie section
	sectionKey, _ := strconv.Atoi(movieSectionKey)
	log.Debug().
		Int("sectionKey", sectionKey).
		Msg("Making API request to Plex server for movies")

	res, err := c.plexAPI.Library.GetLibraryItems(ctx, operations.GetLibraryItemsRequest{
		Tag:         "all",
		Type:        operations.GetLibraryItemsQueryParamTypeMovie,
		SectionKey:  sectionKey,
		IncludeMeta: operations.GetLibraryItemsQueryParamIncludeMetaEnable.ToPointer(),
	})

	// TODO: this interface here is wrong, need differernt type
	// log.Debug().Interface("response", res).Msg("Response from Plex")

	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Int("sectionKey", sectionKey).
			Msg("Failed to get movies from Plex")
		return nil, fmt.Errorf("failed to get movies: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil {
		log.Info().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("No movies found in Plex")
		return []interfaces.Movie{}, nil
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("totalItems", len(res.Object.MediaContainer.Metadata)).
		Msg("Successfully retrieved movies from Plex")

	movies := make([]interfaces.Movie, 0, len(res.Object.MediaContainer.Metadata))
	for _, item := range res.Object.MediaContainer.Metadata {
		if item.Type != "movie" {
			continue
		}

		movie := interfaces.Movie{
			MediaItem: interfaces.MediaItem{
				ExternalID: item.RatingKey,
				ClientID:   c.ClientID,
				ClientType: string(c.ClientType),
				Metadata:   c.createMetadataFromPlexItem(&item),
			},
		}

		c.BaseMediaClient.AddClientInfo(&movie.MediaItem)
		movies = append(movies, movie)

		log.Debug().
			Str("movieID", item.RatingKey).
			Str("movieTitle", item.Title).
			Msg("Added movie to result list")
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("moviesReturned", len(movies)).
		Msg("Completed GetMovies request")

	return movies, nil
}

// GetTVShows retrieves TV shows from Plex
func (c *PlexClient) GetTVShows(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.TVShow, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("baseURL", c.baseURL).
		Msg("Retrieving TV shows from Plex server")

	// First, find the TV show library section
	log.Debug().Msg("Finding TV show library section")
	tvSectionKey, err := c.findLibrarySectionByType(ctx, "show")
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("Failed to find TV show library section")
		return nil, err
	}

	if tvSectionKey == "" {
		log.Info().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("No TV show library section found in Plex")
		return []interfaces.TVShow{}, nil
	}

	// Get TV shows from the TV section
	sectionKey, _ := strconv.Atoi(tvSectionKey)
	log.Debug().
		Int("sectionKey", sectionKey).
		Msg("Making API request to Plex server for TV shows")

	res, err := c.plexAPI.Library.GetLibraryItems(ctx, operations.GetLibraryItemsRequest{
		IncludeMeta: operations.GetLibraryItemsQueryParamIncludeMetaEnable.ToPointer(),
		Tag:         "all",
		Type:        operations.GetLibraryItemsQueryParamTypeTvShow,
		SectionKey:  sectionKey,
	})
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Int("sectionKey", sectionKey).
			Msg("Failed to get TV shows from Plex")
		return nil, fmt.Errorf("failed to get TV shows: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil {
		log.Info().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("No TV shows found in Plex")
		return []interfaces.TVShow{}, nil
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("totalItems", len(res.Object.MediaContainer.Metadata)).
		Msg("Successfully retrieved TV shows from Plex")

	shows := make([]interfaces.TVShow, 0, len(res.Object.MediaContainer.Metadata))
	for _, item := range res.Object.MediaContainer.Metadata {
		if item.Type != "show" {
			continue
		}

		show := interfaces.TVShow{
			MediaItem: interfaces.MediaItem{
				ExternalID: item.RatingKey,
				ClientID:   c.ClientID,
				ClientType: string(c.ClientType),
				Metadata:   c.createMetadataFromPlexItem(&item),
			},
		}

		if item.Rating != nil {
			show.Rating = float64(*item.Rating)
		}
		if item.Year != nil {
			show.ReleaseYear = *item.Year
		}
		if item.ContentRating != nil {
			show.ContentRating = *item.ContentRating
		}
		if item.ChildCount != nil {
			show.SeasonCount = *item.ChildCount
		}
		if item.LeafCount != nil {
			show.EpisodeCount = int(*item.LeafCount)
		}

		if item.Genre != nil {
			show.Genres = make([]string, 0, len(item.Genre))
			for _, genre := range item.Genre {
				if genre.Tag != nil {
					show.Genres = append(show.Genres, *genre.Tag)
				}
			}
		}

		c.BaseMediaClient.AddClientInfo(&show.MediaItem)
		shows = append(shows, show)

		log.Debug().
			Str("showID", item.RatingKey).
			Str("showTitle", item.Title).
			Int("seasonCount", show.SeasonCount).
			Msg("Added TV show to result list")
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("showsReturned", len(shows)).
		Msg("Completed GetTVShows request")

	return shows, nil
}

// GetTVShowSeasons retrieves seasons for a specific TV show
func (c *PlexClient) GetTVShowSeasons(ctx context.Context, showID string) ([]interfaces.Season, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("showID", showID).
		Str("baseURL", c.baseURL).
		Msg("Retrieving seasons for TV show from Plex server")

	ratingKey, _ := strconv.Atoi(showID)
	float64RatingKey := float64(ratingKey)

	log.Debug().
		Str("showID", showID).
		Float64("ratingKey", float64RatingKey).
		Msg("Making API request to Plex server for TV show seasons")

	childRes, err := c.plexAPI.Library.GetMetadataChildren(ctx, float64RatingKey, plexgo.String("Stream"))
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("showID", showID).
			Msg("Failed to get TV show seasons from Plex")
		return nil, fmt.Errorf("failed to get TV show seasons: %w", err)
	}

	if childRes.Object.MediaContainer == nil || childRes.Object.MediaContainer.Metadata == nil {
		log.Info().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("showID", showID).
			Msg("No seasons found for TV show in Plex")
		return []interfaces.Season{}, nil
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("showID", showID).
		Int("totalItems", len(childRes.Object.MediaContainer.Metadata)).
		Msg("Successfully retrieved seasons for TV show from Plex")

	seasons := make([]interfaces.Season, 0, len(childRes.Object.MediaContainer.Metadata))
	for _, item := range childRes.Object.MediaContainer.Metadata {
		if *item.Type != "season" {
			continue
		}

		season := interfaces.Season{
			MediaItem: interfaces.MediaItem{
				ExternalID: *item.RatingKey,
				ClientID:   c.ClientID,
				ClientType: string(c.ClientType),
				Metadata: interfaces.MediaMetadata{
					Description: *item.Summary,
					Title:       *item.Title,
					Artwork: interfaces.Artwork{
						Thumbnail: c.makeFullURL(*item.Thumb),
					},
					ExternalIDs: interfaces.ExternalIDs{interfaces.ExternalID{
						Source: "plex",
						ID:     *item.RatingKey,
					}},
					UpdatedAt: time.Unix(int64(*item.UpdatedAt), 0),
					AddedAt:   time.Unix(int64(*item.AddedAt), 0),
				},
			},
			EpisodeCount: *item.LeafCount,
			Number:       *item.Index,
		}

		c.BaseMediaClient.AddClientInfo(&season.MediaItem)
		seasons = append(seasons, season)

		log.Debug().
			Str("seasonID", *item.RatingKey).
			Str("seasonTitle", *item.Title).
			Int("seasonNumber", *item.Index).
			Int("episodeCount", *item.LeafCount).
			Msg("Added season to result list")
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("showID", showID).
		Int("seasonsReturned", len(seasons)).
		Msg("Completed GetTVShowSeasons request")

	return seasons, nil
}

// GetTVShowEpisodes retrieves episodes for a specific season of a TV show
func (c *PlexClient) GetTVShowEpisodes(ctx context.Context, showID string, seasonNumber int) ([]interfaces.Episode, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("showID", showID).
		Int("seasonNumber", seasonNumber).
		Str("baseURL", c.baseURL).
		Msg("Retrieving episodes for TV show season from Plex server")

	// First get all seasons
	log.Debug().
		Str("showID", showID).
		Int("seasonNumber", seasonNumber).
		Msg("Getting seasons for the TV show")

	seasons, err := c.GetTVShowSeasons(ctx, showID)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("showID", showID).
			Int("seasonNumber", seasonNumber).
			Msg("Failed to get seasons for TV show")
		return nil, err
	}

	var seasonID string
	for _, season := range seasons {
		if season.Number == seasonNumber {
			for _, externalID := range season.MediaItem.Metadata.ExternalIDs {
				if externalID.Source == "plex" {
					seasonID = externalID.ID
					break
				}
			}
			break
		}
	}

	if seasonID == "" {
		log.Info().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("showID", showID).
			Int("seasonNumber", seasonNumber).
			Msg("Season not found for TV show in Plex")
		return []interfaces.Episode{}, nil
	}

	ratingKey, _ := strconv.Atoi(seasonID)
	float64RatingKey := float64(ratingKey)

	log.Debug().
		Str("seasonID", seasonID).
		Float64("ratingKey", float64RatingKey).
		Msg("Making API request to Plex server for TV show episodes")

	childRes, err := c.plexAPI.Library.GetMetadataChildren(ctx, float64RatingKey, plexgo.String("Stream"))
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("showID", showID).
			Int("seasonNumber", seasonNumber).
			Str("seasonID", seasonID).
			Msg("Failed to get TV show episodes from Plex")
		return nil, fmt.Errorf("failed to get TV show episodes: %w", err)
	}

	if childRes.Object.MediaContainer == nil || childRes.Object.MediaContainer.Metadata == nil {
		log.Info().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("showID", showID).
			Int("seasonNumber", seasonNumber).
			Msg("No episodes found for TV show season in Plex")
		return []interfaces.Episode{}, nil
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("showID", showID).
		Int("seasonNumber", seasonNumber).
		Int("totalItems", len(childRes.Object.MediaContainer.Metadata)).
		Msg("Successfully retrieved episodes for TV show season from Plex")

	episodes := make([]interfaces.Episode, 0, len(childRes.Object.MediaContainer.Metadata))
	for _, item := range childRes.Object.MediaContainer.Metadata {
		if *item.Type != "episode" {
			continue
		}

		episode := interfaces.Episode{
			ShowID:   showID,
			SeasonID: *item.ParentKey,
			MediaItem: interfaces.MediaItem{
				ExternalID: *item.RatingKey,
				ClientID:   c.ClientID,
				ClientType: string(c.ClientType),
				Metadata: interfaces.MediaMetadata{
					Description: *item.Summary,
					Title:       *item.Title,
					Artwork: interfaces.Artwork{
						Thumbnail: c.makeFullURL(*item.Thumb),
					},
					UpdatedAt: time.Unix(int64(*item.UpdatedAt), 0),
					AddedAt:   time.Unix(int64(*item.AddedAt), 0),
				},
			},
			Number:       int64(*item.Index),
			SeasonNumber: int(*item.ParentIndex),
		}

		// Add studio if available
		if item.ParentStudio != nil {
			episode.MediaItem.Metadata.Studios = []string{*item.ParentStudio}
		}

		c.BaseMediaClient.AddClientInfo(&episode.MediaItem)
		episodes = append(episodes, episode)

		log.Debug().
			Str("episodeID", *item.RatingKey).
			Str("episodeTitle", *item.Title).
			Int("seasonNumber", episode.SeasonNumber).
			Int64("episodeNumber", episode.Number).
			Msg("Added episode to result list")
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("showID", showID).
		Int("seasonNumber", seasonNumber).
		Int("episodesReturned", len(episodes)).
		Msg("Completed GetTVShowEpisodes request")

	return episodes, nil
}

// GetMusic retrieves music tracks from Plex
func (c *PlexClient) GetMusic(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.MusicTrack, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("baseURL", c.baseURL).
		Msg("Retrieving music tracks from Plex server")

	// Find the music library section
	log.Debug().Msg("Finding music library section")
	musicSectionKey, err := c.findLibrarySectionByType(ctx, "artist")
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("Failed to find music library section")
		return nil, err
	}

	if musicSectionKey == "" {
		log.Info().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("No music library section found in Plex")
		return []interfaces.MusicTrack{}, nil
	}

	// For tracks, we need to traverse the hierarchy: artists > albums > tracks
	sectionKey, _ := strconv.Atoi(musicSectionKey)
	float64SectionKey := float64(sectionKey)

	log.Debug().
		Int("sectionKey", sectionKey).
		Msg("Making API request to Plex server for music artists")

	res, err := c.plexAPI.Library.GetMetadataChildren(ctx, float64SectionKey, plexgo.String("Stream"))
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Int("sectionKey", sectionKey).
			Msg("Failed to get music artists from Plex")
		return nil, fmt.Errorf("failed to get music artists: %w", err)
	}

	var tracks []interfaces.MusicTrack

	// Loop through artists
	if res.Object.MediaContainer != nil && res.Object.MediaContainer.Metadata != nil {
		log.Debug().
			Int("artistCount", len(res.Object.MediaContainer.Metadata)).
			Msg("Processing artists to find music tracks")

		for _, artist := range res.Object.MediaContainer.Metadata {
			artistKey, _ := strconv.Atoi(*artist.RatingKey)
			float64ArtistKey := float64(artistKey)

			log.Debug().
				Str("artistID", *artist.RatingKey).
				Str("artistName", *artist.Title).
				Msg("Getting albums for artist")

			albumsRes, err := c.plexAPI.Library.GetMetadataChildren(ctx, float64ArtistKey, plexgo.String("Stream"))
			if err != nil {
				log.Warn().
					Err(err).
					Str("artistID", *artist.RatingKey).
					Str("artistName", *artist.Title).
					Msg("Failed to get albums for artist, skipping")
				continue
			}

			// Loop through albums
			if albumsRes.Object.MediaContainer != nil && albumsRes.Object.MediaContainer.Metadata != nil {
				log.Debug().
					Str("artistID", *artist.RatingKey).
					Str("artistName", *artist.Title).
					Int("albumCount", len(albumsRes.Object.MediaContainer.Metadata)).
					Msg("Processing albums to find tracks")

				for _, album := range albumsRes.Object.MediaContainer.Metadata {
					albumKey, _ := strconv.Atoi(*album.RatingKey)

					float64AlbumKey := float64(albumKey)

					log.Debug().
						Str("albumID", *album.RatingKey).
						Str("albumName", *album.Title).
						Msg("Getting tracks for album")

					tracksRes, err := c.plexAPI.Library.GetMetadataChildren(ctx, float64AlbumKey, plexgo.String("Stream"))
					if err != nil {
						log.Warn().
							Err(err).
							Str("albumID", *album.RatingKey).
							Str("albumName", *album.Title).
							Msg("Failed to get tracks for album, skipping")
						continue
					}

					// Loop through tracks
					if tracksRes.Object.MediaContainer != nil && tracksRes.Object.MediaContainer.Metadata != nil {
						log.Debug().
							Str("albumID", *album.RatingKey).
							Str("albumName", *album.Title).
							Int("trackCount", len(tracksRes.Object.MediaContainer.Metadata)).
							Msg("Processing tracks for album")

						for _, item := range tracksRes.Object.MediaContainer.Metadata {
							if *item.Type != "track" {
								continue
							}

							track := interfaces.MusicTrack{
								MediaItem: interfaces.MediaItem{
									ExternalID: *item.RatingKey,
									ClientID:   c.ClientID,
									ClientType: string(c.ClientType),
									Metadata:   c.createChildMetadataFromPlexItem(&item),
								},
								Number:     *item.Index,
								ArtistID:   *artist.RatingKey,
								ArtistName: *artist.Title,
								AlbumID:    *album.RatingKey,
								AlbumName:  *album.Title,
							}
							c.BaseMediaClient.AddClientInfo(&track.MediaItem)
							tracks = append(tracks, track)

							log.Debug().
								Str("trackID", *item.RatingKey).
								Str("trackTitle", *item.Title).
								Int("trackNumber", *item.Index).
								Msg("Added track to result list")

							// Limit number of tracks to avoid too large responses
							if len(tracks) >= 100 {
								log.Info().
									Int("trackCount", len(tracks)).
									Msg("Reached track limit (100), returning results")
								return tracks, nil
							}
						}
					}
				}
			}
		}
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("tracksReturned", len(tracks)).
		Msg("Completed GetMusic request")

	return tracks, nil
}

// GetMusicArtists retrieves music artists from Plex
func (c *PlexClient) GetMusicArtists(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.MusicArtist, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("baseURL", c.baseURL).
		Msg("Retrieving music artists from Plex server")

	// Find the music library section
	log.Debug().Msg("Finding music library section")
	musicSectionKey, err := c.findLibrarySectionByType(ctx, "artist")
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("Failed to find music library section")
		return nil, err
	}

	if musicSectionKey == "" {
		log.Info().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("No music library section found in Plex")
		return []interfaces.MusicArtist{}, nil
	}

	sectionKey, _ := strconv.Atoi(musicSectionKey)
	float64SectionKey := float64(sectionKey)

	log.Debug().
		Int("sectionKey", sectionKey).
		Msg("Making API request to Plex server for music artists")

	res, err := c.plexAPI.Library.GetMetadataChildren(ctx, float64SectionKey, plexgo.String("Stream"))
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Int("sectionKey", sectionKey).
			Msg("Failed to get music artists from Plex")
		return nil, fmt.Errorf("failed to get music artists: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil {
		log.Info().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("No music artists found in Plex")
		return []interfaces.MusicArtist{}, nil
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("totalItems", len(res.Object.MediaContainer.Metadata)).
		Msg("Successfully retrieved music artists from Plex")

	artists := make([]interfaces.MusicArtist, 0, len(res.Object.MediaContainer.Metadata))
	for _, item := range res.Object.MediaContainer.Metadata {
		if *item.Type != "artist" {
			continue
		}

		artist := interfaces.MusicArtist{
			MediaItem: interfaces.MediaItem{
				ExternalID: *item.RatingKey,
				ClientID:   c.ClientID,
				ClientType: string(c.ClientType),
				Metadata:   c.createChildMetadataFromPlexItem(&item),
			},
		}

		c.BaseMediaClient.AddClientInfo(&artist.MediaItem)
		artists = append(artists, artist)

		log.Debug().
			Str("artistID", *item.RatingKey).
			Str("artistName", *item.Title).
			Msg("Added artist to result list")
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("artistsReturned", len(artists)).
		Msg("Completed GetMusicArtists request")

	return artists, nil
}

// GetMusicAlbums retrieves music albums from Plex
func (c *PlexClient) GetMusicAlbums(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.MusicAlbum, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("baseURL", c.baseURL).
		Msg("Retrieving music albums from Plex server")

	// Find the music library section
	log.Debug().Msg("Finding music library section")
	musicSectionKey, err := c.findLibrarySectionByType(ctx, "artist")
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("Failed to find music library section")
		return nil, err
	}

	if musicSectionKey == "" {
		log.Info().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("No music library section found in Plex")
		return []interfaces.MusicAlbum{}, nil
	}

	// For albums, we need to traverse artists first
	sectionKey, _ := strconv.Atoi(musicSectionKey)
	float64SectionKey := float64(sectionKey)

	log.Debug().
		Int("sectionKey", sectionKey).
		Msg("Making API request to Plex server for music artists")

	res, err := c.plexAPI.Library.GetMetadataChildren(ctx, float64SectionKey, plexgo.String("Stream"))
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Int("sectionKey", sectionKey).
			Msg("Failed to get music artists from Plex")
		return nil, fmt.Errorf("failed to get music artists: %w", err)
	}

	var albums []interfaces.MusicAlbum

	// Loop through artists to get their albums
	if res.Object.MediaContainer != nil && res.Object.MediaContainer.Metadata != nil {
		log.Debug().
			Int("artistCount", len(res.Object.MediaContainer.Metadata)).
			Msg("Processing artists to find albums")

		for _, artist := range res.Object.MediaContainer.Metadata {
			artistKey, _ := strconv.Atoi(*artist.RatingKey)
			float64ArtistKey := float64(artistKey)

			log.Debug().
				Str("artistID", *artist.RatingKey).
				Str("artistName", *artist.Title).
				Msg("Getting albums for artist")

			albumsRes, err := c.plexAPI.Library.GetMetadataChildren(ctx, float64ArtistKey, plexgo.String("Stream"))
			if err != nil {
				log.Warn().
					Err(err).
					Str("artistID", *artist.RatingKey).
					Str("artistName", *artist.Title).
					Msg("Failed to get albums for artist, skipping")
				continue
			}

			if albumsRes.Object.MediaContainer != nil && albumsRes.Object.MediaContainer.Metadata != nil {
				log.Debug().
					Str("artistID", *artist.RatingKey).
					Str("artistName", *artist.Title).
					Int("albumCount", len(albumsRes.Object.MediaContainer.Metadata)).
					Msg("Processing albums for artist")

				for _, item := range albumsRes.Object.MediaContainer.Metadata {
					if *item.Type != "album" {
						continue
					}

					album := interfaces.MusicAlbum{
						MediaItem: interfaces.MediaItem{
							ExternalID: *item.RatingKey,
							ClientID:   c.ClientID,
							ClientType: string(c.ClientType),
							Metadata:   c.createChildMetadataFromPlexItem(&item),
						},
						ArtistID:   *artist.RatingKey,
						ArtistName: *artist.Title,
						TrackCount: *item.LeafCount,
					}

					c.BaseMediaClient.AddClientInfo(&album.MediaItem)
					albums = append(albums, album)

					log.Debug().
						Str("albumID", *item.RatingKey).
						Str("albumName", *item.Title).
						Int("trackCount", *item.LeafCount).
						Msg("Added album to result list")
				}
			}
		}
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("albumsReturned", len(albums)).
		Msg("Completed GetMusicAlbums request")

	return albums, nil
}

// GetWatchHistory retrieves watch history from Plex
func (c *PlexClient) GetWatchHistory(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.WatchHistoryItem, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("baseURL", c.baseURL).
		Msg("Retrieving watch history from Plex server")

	log.Warn().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Msg("Watch history retrieval not yet implemented for Plex")

	// This would require querying Plex for watch history
	return []interfaces.WatchHistoryItem{}, fmt.Errorf("Watch history retrieval not yet implemented for Plex")
}

// GetMovieByID retrieves a specific movie by ID
func (c *PlexClient) GetMovieByID(ctx context.Context, id string) (interfaces.Movie, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("movieID", id).
		Str("baseURL", c.baseURL).
		Msg("Retrieving specific movie from Plex server")

	ratingKey, _ := strconv.Atoi(id)
	int64RatingKey := int64(ratingKey)

	log.Debug().
		Str("movieID", id).
		Int64("ratingKey", int64RatingKey).
		Msg("Making API request to Plex server for movie")

	res, err := c.plexAPI.Library.GetMediaMetaData(ctx, operations.GetMediaMetaDataRequest{
		RatingKey: int64RatingKey,
	})
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("movieID", id).
			Msg("Failed to get movie from Plex")
		return interfaces.Movie{}, fmt.Errorf("failed to get movie: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil || len(res.Object.MediaContainer.Metadata) == 0 {
		log.Error().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("movieID", id).
			Msg("Movie not found in Plex")
		return interfaces.Movie{}, fmt.Errorf("movie not found")
	}

	item := res.Object.MediaContainer.Metadata[0]
	if item.Type != "movie" {
		log.Error().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("movieID", id).
			Str("actualType", item.Type).
			Msg("Item retrieved is not a movie")
		return interfaces.Movie{}, fmt.Errorf("item is not a movie")
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("movieID", id).
		Str("movieTitle", item.Title).
		Msg("Successfully retrieved movie from Plex")

	movie := interfaces.Movie{
		MediaItem: interfaces.MediaItem{
			ExternalID: item.RatingKey,
			ClientID:   c.ClientID,
			ClientType: string(c.ClientType),
			Metadata:   c.createMediaMetadataFromPlexItem(&item),
		},
	}

	c.BaseMediaClient.AddClientInfo(&movie.MediaItem)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("movieID", id).
		Str("movieTitle", movie.Metadata.Title).
		Msg("Successfully converted movie data")

	return movie, nil
}

// GetTVShowByID retrieves a specific TV show by ID
func (c *PlexClient) GetTVShowByID(ctx context.Context, id string) (interfaces.TVShow, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("showID", id).
		Str("baseURL", c.baseURL).
		Msg("Retrieving specific TV show from Plex server")

	ratingKey, _ := strconv.Atoi(id)
	int64RatingKey := int64(ratingKey)

	log.Debug().
		Str("showID", id).
		Int64("ratingKey", int64RatingKey).
		Msg("Making API request to Plex server for TV show")

	res, err := c.plexAPI.Library.GetMediaMetaData(ctx, operations.GetMediaMetaDataRequest{
		RatingKey: int64RatingKey,
	})

	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("showID", id).
			Msg("Failed to get TV show from Plex")
		return interfaces.TVShow{}, fmt.Errorf("failed to get TV show: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil || len(res.Object.MediaContainer.Metadata) == 0 {
		log.Error().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("showID", id).
			Msg("TV show not found in Plex")
		return interfaces.TVShow{}, fmt.Errorf("TV show not found")
	}

	item := res.Object.MediaContainer.Metadata[0]
	if item.Type != "show" {
		log.Error().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("showID", id).
			Str("actualType", item.Type).
			Msg("Item retrieved is not a TV show")
		return interfaces.TVShow{}, fmt.Errorf("item is not a TV show")
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("showID", id).
		Str("showTitle", item.Title).
		Msg("Successfully retrieved TV show from Plex")

	show := interfaces.TVShow{
		MediaItem: interfaces.MediaItem{
			ExternalID: item.RatingKey,
			ClientID:   c.ClientID,
			ClientType: string(c.ClientType),
			Metadata:   c.createMediaMetadataFromPlexItem(&item),
		},
	}

	if item.Rating != nil {
		show.Rating = float64(*item.Rating)
	}
	if item.ContentRating != nil {
		show.ContentRating = *item.ContentRating
	}
	if item.ChildCount != nil {
		show.SeasonCount = *item.ChildCount
	}
	if item.LeafCount != nil {
		show.EpisodeCount = int(*item.LeafCount)
	}

	if item.Genre != nil {
		show.Genres = make([]string, 0, len(item.Genre))
		for _, genre := range item.Genre {
			show.Genres = append(show.Genres, genre.Tag)
		}
	}

	c.BaseMediaClient.AddClientInfo(&show.MediaItem)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("showID", id).
		Str("showTitle", show.Metadata.Title).
		Int("seasonCount", show.SeasonCount).
		Msg("Successfully converted TV show data")

	return show, nil
}

// GetEpisodeByID retrieves a specific episode by ID
func (c *PlexClient) GetEpisodeByID(ctx context.Context, id string) (interfaces.Episode, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("episodeID", id).
		Msg("Retrieving specific episode from Plex server")

	ratingKey, _ := strconv.Atoi(id)
	int64RatingKey := int64(ratingKey)

	res, err := c.plexAPI.Library.GetMediaMetaData(ctx, operations.GetMediaMetaDataRequest{
		RatingKey: int64RatingKey,
	})
	if err != nil {
		log.Error().Err(err).Str("episodeID", id).Msg("Failed to get episode from Plex")
		return interfaces.Episode{}, fmt.Errorf("failed to get episode: %w", err)
	}

	if res.Object.MediaContainer == nil ||
		res.Object.MediaContainer.Metadata == nil ||
		len(res.Object.MediaContainer.Metadata) == 0 {
		log.Error().Str("episodeID", id).Msg("Episode not found in Plex")
		return interfaces.Episode{}, fmt.Errorf("episode not found")
	}

	item := res.Object.MediaContainer.Metadata[0]
	if item.Type != "episode" {
		log.Error().Str("episodeID", id).Str("actualType", item.Type).Msg("Item retrieved is not an episode")
		return interfaces.Episode{}, fmt.Errorf("item is not an episode")
	}

	episode := interfaces.Episode{
		MediaItem: interfaces.MediaItem{
			ExternalID: item.RatingKey,
			ClientID:   c.ClientID,
			ClientType: string(c.ClientType),
			Metadata:   c.createMediaMetadataFromPlexItem(&item),
		},
		Number: int64(*item.Index),
	}

	// Add season number if available
	if item.ParentIndex != nil {
		episode.SeasonNumber = int(*item.ParentIndex)
	}

	// Add show ID if available (via grandparentRatingKey)
	if item.GrandparentRatingKey != nil {
		episode.ShowID = *item.GrandparentRatingKey
	}

	// Add studio if available
	if item.Studio != nil {
		episode.MediaItem.Metadata.Studios = []string{*item.Studio}
	}

	c.BaseMediaClient.AddClientInfo(&episode.MediaItem)

	log.Info().
		Str("episodeID", id).
		Str("episodeTitle", episode.Metadata.Title).
		Int("seasonNumber", episode.SeasonNumber).
		Int64("episodeNumber", episode.Number).
		Msg("Successfully retrieved episode")

	return episode, nil
}

// GetMusicTrackByID retrieves a specific music track by ID
func (c *PlexClient) GetMusicTrackByID(ctx context.Context, id string) (interfaces.MusicTrack, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("trackID", id).
		Str("baseURL", c.baseURL).
		Msg("Retrieving specific music track from Plex server")

	ratingKey, _ := strconv.Atoi(id)
	int64RatingKey := int64(ratingKey)

	log.Debug().
		Str("trackID", id).
		Int64("ratingKey", int64RatingKey).
		Msg("Making API request to Plex server for music track")

	res, err := c.plexAPI.Library.GetMediaMetaData(ctx, operations.GetMediaMetaDataRequest{RatingKey: int64RatingKey})
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("trackID", id).
			Msg("Failed to get music track from Plex")
		return interfaces.MusicTrack{}, fmt.Errorf("failed to get music track: %w", err)
	}

	if res.Object.MediaContainer == nil || res.Object.MediaContainer.Metadata == nil || len(res.Object.MediaContainer.Metadata) == 0 {
		log.Error().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("trackID", id).
			Msg("Music track not found in Plex")
		return interfaces.MusicTrack{}, fmt.Errorf("music track not found")
	}

	item := res.Object.MediaContainer.Metadata[0]
	if item.Type != "track" {
		log.Error().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Str("trackID", id).
			Str("actualType", item.Type).
			Msg("Item retrieved is not a music track")
		return interfaces.MusicTrack{}, fmt.Errorf("item is not a music track")
	}

	// Get album and artist info
	var albumName string
	var artistID string
	var artistName string

	if item.ParentRatingKey != nil {
		// Get album info
		albumKey, _ := strconv.Atoi(*item.ParentRatingKey)
		int64AlbumKey := int64(albumKey)

		log.Debug().
			Str("albumID", *item.ParentRatingKey).
			Int64("albumKey", int64AlbumKey).
			Msg("Getting parent album information")

		albumRes, err := c.plexAPI.Library.GetMediaMetaData(ctx, operations.GetMediaMetaDataRequest{
			RatingKey: int64AlbumKey,
		})

		if err == nil && albumRes.Object.MediaContainer != nil &&
			albumRes.Object.MediaContainer.Metadata != nil &&
			len(albumRes.Object.MediaContainer.Metadata) > 0 {

			albumName = albumRes.Object.MediaContainer.Metadata[0].Title

			// Get artist info if available
			if albumRes.Object.MediaContainer.Metadata[0].ParentRatingKey != nil {
				artistID = *albumRes.Object.MediaContainer.Metadata[0].ParentRatingKey

				log.Debug().
					Str("artistID", artistID).
					Msg("Getting parent artist information")

				artistKey, _ := strconv.Atoi(artistID)
				int64ArtistKey := int64(artistKey)
				artistRes, err := c.plexAPI.Library.GetMediaMetaData(ctx, operations.GetMediaMetaDataRequest{
					RatingKey: int64ArtistKey,
				})

				if err == nil && artistRes.Object.MediaContainer != nil &&
					artistRes.Object.MediaContainer.Metadata != nil &&
					len(artistRes.Object.MediaContainer.Metadata) > 0 {

					artistName = artistRes.Object.MediaContainer.Metadata[0].Title
					log.Debug().
						Str("artistName", artistName).
						Msg("Retrieved artist name")
				}
			}
		}
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("trackID", id).
		Str("trackTitle", item.Title).
		Str("albumName", albumName).
		Str("artistName", artistName).
		Msg("Successfully retrieved music track from Plex")

	track := interfaces.MusicTrack{
		MediaItem: interfaces.MediaItem{
			ExternalID: item.RatingKey,
			ClientID:   c.ClientID,
			ClientType: string(c.ClientType),
			Metadata:   c.createMediaMetadataFromPlexItem(&item),
		},
		AlbumName:  albumName,
		ArtistName: artistName,
		ArtistID:   artistID,
		AlbumID:    *item.ParentRatingKey,
		Number:     int(*item.Index),
	}

	c.BaseMediaClient.AddClientInfo(&track.MediaItem)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("trackID", id).
		Str("trackTitle", track.Metadata.Title).
		Int("trackNumber", track.Number).
		Msg("Successfully converted music track data")

	return track, nil
}

// GetMusicGenres retrieves music genres from Plex
func (c *PlexClient) GetMusicGenres(ctx context.Context) ([]string, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("baseURL", c.baseURL).
		Msg("Retrieving music genres from Plex server")

	// Find the music library section
	log.Debug().Msg("Finding music library section")
	musicSectionKey, err := c.findLibrarySectionByType(ctx, "artist")
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("Failed to find music library section")
		return nil, err
	}

	if musicSectionKey == "" {
		log.Info().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("No music library section found in Plex")
		return []string{}, nil
	}

	// Get genres from the library items
	sectionKey, _ := strconv.Atoi(musicSectionKey)

	log.Debug().
		Int("sectionKey", sectionKey).
		Msg("Making API request to Plex server for music content")

	res, err := c.plexAPI.Library.GetLibraryItems(ctx, operations.GetLibraryItemsRequest{
		IncludeMeta: operations.GetLibraryItemsQueryParamIncludeMetaEnable.ToPointer(),
		Tag:         "all",
		Type:        operations.GetLibraryItemsQueryParamTypeAudio,
		SectionKey:  sectionKey})

	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Int("sectionKey", sectionKey).
			Msg("Failed to get music content from Plex")
		return nil, fmt.Errorf("failed to get music content: %w", err)
	}

	genreMap := make(map[string]bool)
	if res.Object.MediaContainer != nil && res.Object.MediaContainer.Metadata != nil {
		log.Debug().
			Int("contentCount", len(res.Object.MediaContainer.Metadata)).
			Msg("Extracting genres from music content")

		for _, item := range res.Object.MediaContainer.Metadata {
			if item.Genre != nil {
				for _, genre := range item.Genre {
					if genre.Tag != nil {
						genreMap[*genre.Tag] = true
						log.Debug().
							Str("genre", *genre.Tag).
							Msg("Found music genre")
					}
				}
			}
		}
	}

	genres := make([]string, 0, len(genreMap))
	for genre := range genreMap {
		genres = append(genres, genre)
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("genresFound", len(genres)).
		Msg("Successfully retrieved music genres from Plex")

	return genres, nil
}

// GetMovieGenres retrieves movie genres from Plex
func (c *PlexClient) GetMovieGenres(ctx context.Context) ([]string, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Str("baseURL", c.baseURL).
		Msg("Retrieving movie genres from Plex server")

	// Find the movie library section
	log.Debug().Msg("Finding movie library section")
	movieSectionKey, err := c.findLibrarySectionByType(ctx, "movie")
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("Failed to find movie library section")
		return nil, err
	}

	if movieSectionKey == "" {
		log.Info().
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Msg("No movie library section found in Plex")
		return []string{}, nil
	}

	// Get genres directly from the genre endpoint
	sectionKey, _ := strconv.Atoi(movieSectionKey)

	log.Debug().
		Int("sectionKey", sectionKey).
		Msg("Making API request to Plex server for movie genres")

	res, err := c.plexAPI.Library.GetGenresLibrary(ctx, sectionKey, operations.GetGenresLibraryQueryParamTypeMovie)
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", c.ClientID).
			Str("clientType", string(c.ClientType)).
			Int("sectionKey", sectionKey).
			Msg("Failed to get movie genres from Plex")
		return nil, fmt.Errorf("failed to get movie genres: %w", err)
	}

	genreMap := make(map[string]bool)
	if res.Object.MediaContainer != nil {
		directories := res.Object.MediaContainer.GetDirectory()
		log.Debug().
			Int("genreCount", len(directories)).
			Msg("Extracting genres from directories")

		for _, item := range directories {
			genreMap[item.Title] = true
			log.Debug().
				Str("genre", item.Title).
				Msg("Found movie genre")
		}
	}

	genres := make([]string, 0, len(genreMap))
	for genre := range genreMap {
		genres = append(genres, genre)
	}

	log.Info().
		Uint64("clientID", c.ClientID).
		Str("clientType", string(c.ClientType)).
		Int("genresFound", len(genres)).
		Msg("Successfully retrieved movie genres from Plex")

	return genres, nil
}

