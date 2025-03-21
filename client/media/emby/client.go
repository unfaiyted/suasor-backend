package emby

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/antihax/optional"
	"suasor/client/media/interfaces"
	embyclient "suasor/internal/clients/embyAPI"
	"suasor/models"
	"suasor/utils"
)

// Configuration holds Emby connection settings
type Configuration struct {
	BaseURL  string
	ApiKey   string
	Username string // For user-specific queries
	UserID   string
}

// EmbyClient implements the MediaContentProvider interface
type EmbyClient struct {
	interfaces.BaseMediaClient
	client *embyclient.APIClient
	config Configuration
}

// NewEmbyClient creates a new Emby client instance
func NewEmbyClient(ctx context.Context, clientID uint64, config interface{}) (interfaces.MediaContentProvider, error) {
	cfg, ok := config.(Configuration)
	if !ok {
		return nil, fmt.Errorf("invalid configuration for Emby client")
	}

	// Create API client configuration
	apiConfig := embyclient.NewConfiguration()
	apiConfig.BasePath = cfg.BaseURL

	// Set up API key in default headers
	apiConfig.DefaultHeader = map[string]string{
		"X-Emby-Token": cfg.ApiKey,
	}

	client := embyclient.NewAPIClient(apiConfig)

	embyClient := &EmbyClient{
		BaseMediaClient: interfaces.BaseMediaClient{
			ClientID:   clientID,
			ClientType: models.MediaClientTypeEmby,
		},
		client: client,
		config: cfg,
	}

	// Resolve user ID if username is provided
	if cfg.Username != "" && cfg.UserID == "" {
		if err := embyClient.resolveUserID(ctx); err != nil {
			// Log but don't fail - some operations might work without a user ID
			log := utils.LoggerFromContext(ctx)
			log.Warn().
				Err(err).
				Str("username", cfg.Username).
				Msg("Failed to resolve Emby user ID, some operations may be limited")
		}
	}
	return embyClient, nil
}

// Register the provider factory
func init() {
	interfaces.RegisterProvider(models.MediaClientTypeEmby, NewEmbyClient)
}

// Capability methods
func (e *EmbyClient) SupportsMovies() bool      { return true }
func (e *EmbyClient) SupportsTVShows() bool     { return true }
func (e *EmbyClient) SupportsMusic() bool       { return true }
func (e *EmbyClient) SupportsPlaylists() bool   { return true }
func (e *EmbyClient) SupportsCollections() bool { return true }

// Add this new method to EmbyClient to resolve user ID
func (e *EmbyClient) resolveUserID(ctx context.Context) error {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Str("username", e.config.Username).
		Msg("Resolving Emby user ID from username")

	// Get the list of public users
	users, resp, err := e.client.UserServiceApi.GetUsersPublic(ctx)
	if err != nil {
		log.Error().
			Err(err).
			Str("username", e.config.Username).
			Msg("Failed to fetch Emby users")
		return fmt.Errorf("failed to fetch users: %w", err)
	}

	log.Debug().
		Int("statusCode", resp.StatusCode).
		Int("userCount", len(users)).
		Msg("Retrieved public users from Emby")

	// Find the user with matching username
	for _, user := range users {
		if strings.EqualFold(user.Name, e.config.Username) {
			e.config.UserID = user.Id
			log.Info().
				Str("username", e.config.Username).
				Str("userID", e.config.UserID).
				Msg("Successfully resolved Emby user ID")
			return nil
		}
	}

	log.Warn().
		Str("username", e.config.Username).
		Msg("Could not find matching user in Emby")
	return fmt.Errorf("user '%s' not found in Emby", e.config.Username)
}

// GetMovies retrieves movies from the Emby server
func (e *EmbyClient) GetMovies(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.Movie, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.ClientID).
		Str("clientType", string(e.ClientType)).
		Str("baseURL", e.config.BaseURL).
		Msg("Retrieving movies from Emby server")

	// Create URL parameters
	// url.Values{}
	queryParams := embyclient.ItemsServiceApiGetItemsOpts{
		IncludeItemTypes: optional.NewString("Movie"),
		Recursive:        optional.NewBool(true),
	}

	if options != nil {
		if options.Limit > 0 {
			queryParams.Limit = optional.NewInt32(int32(options.Limit))
		}
		if options.Offset > 0 {
			queryParams.StartIndex = optional.NewInt32(int32(options.Offset))
		}
		// Add sorting if provided
		if options.Sort != "" {
			queryParams.SortBy = optional.NewString(options.Sort)
			if options.SortOrder == "desc" {
				queryParams.SortOrder = optional.NewString("Descending")
			} else {
				queryParams.SortOrder = optional.NewString("Ascending")
			}
		}

		log.Debug().
			Int("limit", options.Limit).
			Int("offset", options.Offset).
			Str("sort", options.Sort).
			Str("sortOrder", options.SortOrder).
			Msg("Applied query options")
	}

	// Call the Emby API
	log.Debug().Msg("Making API request to Emby server")
	items, resp, err := e.client.ItemsServiceApi.GetItems(ctx, &queryParams)
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", e.config.BaseURL).
			Str("apiEndpoint", "/Items").
			Int("statusCode", 0).
			Msg("Failed to fetch movies from Emby")
		return nil, fmt.Errorf("failed to fetch movies: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(items.Items)).
		Int("totalRecordCount", int(items.TotalRecordCount)).
		Msg("Successfully retrieved movies from Emby")

	// Convert results to expected format
	movies := make([]interfaces.Movie, 0)

	for _, item := range items.Items {
		if item.Type_ == "Movie" {
			movie, err := e.convertToMovie(ctx, &item)
			if err != nil {
				// Log error
				log.Warn().
					Err(err).
					Str("movieID", item.Id).
					Str("movieName", item.Name).
					Msg("Error converting Emby item to movie format")
				continue
			}
			movies = append(movies, movie)
		}
	}

	log.Info().
		Int("moviesReturned", len(movies)).
		Msg("Completed GetMovies request")

	return movies, nil
}

// GetMovieByID retrieves a specific movie by ID
func (e *EmbyClient) GetMovieByID(ctx context.Context, id string) (interfaces.Movie, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.ClientID).
		Str("clientType", string(e.ClientType)).
		Str("movieID", id).
		Str("baseURL", e.config.BaseURL).
		Msg("Retrieving specific movie from Emby server")

	// Create query parameters - specifically filtering for movie type
	queryParams := embyclient.ItemsServiceApiGetItemsOpts{
		Ids:              optional.NewString(id),
		IncludeItemTypes: optional.NewString("Movie"), // Ensure we only get movies
		Fields:           optional.NewString("ProductionYear,PremiereDate,ChannelMappingInfo,DateCreated,Genres,IndexOptions,HomePageUrl,Overview,ParentId,Path,ProviderIds,Studios,SortName"),

		// * @param "Fields" (optional.String) -  Optional. Specify additional fields of information to return in the output. This allows multiple, comma delimeted. Options: Budget, Chapters, DateCreated, Genres, HomePageUrl, IndexOptions, MediaStreams, Overview, ParentId, Path, People, ProviderIds, PrimaryImageAspectRatio, Revenue, SortName, Studios, Taglines
	}

	// Call the Emby API
	log.Debug().
		Str("movieID", id).
		Msg("Making API request to Emby server")

	items, resp, err := e.client.ItemsServiceApi.GetItems(ctx, &queryParams)

	log.Info().Any("item", items)

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", e.config.BaseURL).
			Str("apiEndpoint", "/Items").
			Str("movieID", id).
			Int("statusCode", 0).
			Msg("Failed to fetch movie from Emby")
		return interfaces.Movie{}, fmt.Errorf("failed to fetch movie: %w", err)
	}

	// Check if any items were returned
	if len(items.Items) == 0 {
		log.Error().
			Str("movieID", id).
			Int("statusCode", resp.StatusCode).
			Msg("No movie found with the specified ID")
		return interfaces.Movie{}, fmt.Errorf("movie with ID %s not found", id)
	}

	item := items.Items[0]

	// Double-check that the returned item is a movie
	if item.Type_ != "Movie" {
		log.Error().
			Str("movieID", id).
			Str("actualType", item.Type_).
			Msg("Item with specified ID is not a movie")
		return interfaces.Movie{}, fmt.Errorf("item with ID %s is not a movie", id)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Str("movieID", id).
		Str("movieName", item.Name).
		Msg("Successfully retrieved movie from Emby")

	movie, err := e.convertToMovie(ctx, &item)
	if err != nil {
		log.Error().
			Err(err).
			Str("movieID", id).
			Str("movieName", item.Name).
			Msg("Error converting Emby item to movie format")
		return interfaces.Movie{}, fmt.Errorf("error converting movie data: %w", err)
	}

	log.Debug().
		Str("movieID", id).
		Str("movieName", movie.Metadata.Title).
		Int("year", movie.Metadata.ReleaseYear).
		Int32("releaseYear", item.ProductionYear).
		Msg("Successfully returned movie data")

	return movie, nil
}

// Helper function to convert Emby item to internal Movie type

// Helper function to convert Emby item to internal Movie type
func (e *EmbyClient) convertToMovie(ctx context.Context, item *embyclient.BaseItemDto) (interfaces.Movie, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	// Validate required fields
	if item == nil {
		return interfaces.Movie{}, fmt.Errorf("cannot convert nil item to movie")
	}

	if item.Id == "" {
		return interfaces.Movie{}, fmt.Errorf("movie is missing required ID field")
	}

	log.Debug().
		Str("movieID", item.Id).
		Str("movieName", item.Name).
		Int32("releaseYear", item.ProductionYear).
		Str("releaseDate", item.PremiereDate.Format("2006-01-02")).
		Msg("Converting Emby item to movie format")

	var officialRating int

	// Handle empty or non-numeric rating safely
	if item.OfficialRating != "" {
		// Try to convert to integer, but don't fail if it's not a number
		val, err := strconv.Atoi(item.OfficialRating)
		if err == nil {
			officialRating = val
		} else {
			log.Debug().
				Str("movieID", item.Id).
				Str("rating", item.OfficialRating).
				Err(err).
				Msg("Non-numeric rating found, defaulting to 0")
		}
	}

	// Determine release year from either ProductionYear or PremiereDate
	releaseYear := int(item.ProductionYear)
	if releaseYear == 0 && !item.PremiereDate.IsZero() {
		releaseYear = item.PremiereDate.Year()
		log.Debug().
			Str("movieID", item.Id).
			Str("premiereDate", item.PremiereDate.Format("2006-01-02")).
			Int("extractedYear", releaseYear).
			Msg("Using year from premiere date instead of production year")
	}

	// Build movie object with safe handling of optional fields
	movie := interfaces.Movie{
		MediaItem: interfaces.MediaItem{
			Metadata: interfaces.MediaMetadata{
				Title:       item.Name,
				Description: item.Overview,
				ReleaseDate: item.PremiereDate,
				ReleaseYear: releaseYear,
				Genres:      item.Genres,
				Artwork:     e.getArtworkURLs(item),
				Duration:    time.Duration(item.RunTimeTicks/10000000) * time.Second,
				Ratings: interfaces.Ratings{
					interfaces.Rating{
						Source: "emby",
						Value:  float32(officialRating),
					},
				},
			},
			ExternalID: item.Id,
			Type:       "movie",
			ClientID:   e.ClientID,
			ClientType: string(e.ClientType),
		},
	}

	// Only set UserRating if UserData is not nil
	if item.UserData != nil {
		movie.Metadata.UserRating = float32(item.UserData.Rating)
	} else {
		log.Debug().
			Str("movieID", item.Id).
			Msg("Movie has no user data, skipping user rating")
	}

	// Extract provider IDs if available
	if item.ProviderIds != nil {
		ids := *item.ProviderIds
		if imdbID, ok := ids["Imdb"]; ok {
			movie.Metadata.ExternalIDs.AddOrUpdate("imdb", imdbID)
		}
		if tmdbID, ok := ids["Tmdb"]; ok {
			movie.Metadata.ExternalIDs.AddOrUpdate("tmdb", tmdbID)
		}
	}

	log.Debug().
		Str("movieID", item.Id).
		Str("movieTitle", movie.Metadata.Title).
		Int("year", movie.Metadata.ReleaseYear).
		Msg("Successfully converted Emby item to movie")

	return movie, nil
}

// Helper to get image URLs for an item

// Helper to get image URLs for an item with proper nil checks
func (e *EmbyClient) getArtworkURLs(item *embyclient.BaseItemDto) interfaces.Artwork {
	imageURLs := interfaces.Artwork{}

	if item == nil {
		return imageURLs
	}

	baseURL := strings.TrimSuffix(e.config.BaseURL, "/")

	// Primary image (poster) - with nil check
	if item.ImageTags != nil {
		if tag, ok := item.ImageTags["Primary"]; ok {
			imageURLs.Poster = fmt.Sprintf("%s/Items/%s/Images/Primary?tag=%s", baseURL, item.Id, tag)
		}
	}

	// Backdrop image - with nil and length check
	if item.BackdropImageTags != nil && len(item.BackdropImageTags) > 0 {
		imageURLs.Background = fmt.Sprintf("%s/Items/%s/Images/Backdrop?tag=%s", baseURL, item.Id, item.BackdropImageTags[0])
	}

	// Other image types - with nil check
	if item.ImageTags != nil {
		if tag, ok := item.ImageTags["Logo"]; ok {
			imageURLs.Logo = fmt.Sprintf("%s/Items/%s/Images/Logo?tag=%s", baseURL, item.Id, tag)
		}

		if tag, ok := item.ImageTags["Thumb"]; ok {
			imageURLs.Thumbnail = fmt.Sprintf("%s/Items/%s/Images/Thumb?tag=%s", baseURL, item.Id, tag)
		}

		if tag, ok := item.ImageTags["Banner"]; ok {
			imageURLs.Banner = fmt.Sprintf("%s/Items/%s/Images/Banner?tag=%s", baseURL, item.Id, tag)
		}
	}

	return imageURLs
}

// GetCollections retrieves collections from the Emby server
func (e *EmbyClient) GetCollections(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.Collection, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.ClientID).
		Str("clientType", string(e.ClientType)).
		Str("baseURL", e.config.BaseURL).
		Msg("Retrieving collections from Emby server")

	// Create URL parameters
	queryParams := embyclient.ItemsServiceApiGetItemsOpts{
		IncludeItemTypes: optional.NewString("BoxSet"),
		Recursive:        optional.NewBool(true),
	}

	if options != nil {
		if options.Limit > 0 {
			queryParams.Limit = optional.NewInt32(int32(options.Limit))
		}
		if options.Offset > 0 {
			queryParams.StartIndex = optional.NewInt32(int32(options.Offset))
		}
		// Add sorting if provided
		if options.Sort != "" {
			queryParams.SortBy = optional.NewString(options.Sort)
			if options.SortOrder == "desc" {
				queryParams.SortOrder = optional.NewString("Descending")
			} else {
				queryParams.SortOrder = optional.NewString("Ascending")
			}
		}

		log.Debug().
			Int("limit", options.Limit).
			Int("offset", options.Offset).
			Str("sort", options.Sort).
			Str("sortOrder", options.SortOrder).
			Msg("Applied query options")
	}

	// Call the Emby API
	log.Debug().Msg("Making API request to Emby server for collections")
	items, resp, err := e.client.ItemsServiceApi.GetItems(ctx, &queryParams)
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", e.config.BaseURL).
			Str("apiEndpoint", "/Items").
			Int("statusCode", 0).
			Msg("Failed to fetch collections from Emby")
		return nil, fmt.Errorf("failed to fetch collections: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(items.Items)).
		Int("totalRecordCount", int(items.TotalRecordCount)).
		Msg("Successfully retrieved collections from Emby")

	// Convert results to expected format
	collections := make([]interfaces.Collection, 0)
	for _, item := range items.Items {
		if item.Type_ == "BoxSet" {
			collection := interfaces.Collection{
				MediaItem: interfaces.MediaItem{
					Metadata: interfaces.MediaMetadata{
						Title:       item.Name,
						Description: item.Overview,
						Artwork:     e.getArtworkURLs(&item),
					},
					ExternalID: item.Id,
					Type:       "collection",
					ClientID:   e.ClientID,
					ClientType: string(e.ClientType),
				},
			}
			collections = append(collections, collection)
		}
	}

	log.Info().
		Int("collectionsReturned", len(collections)).
		Msg("Completed GetCollections request")

	return collections, nil
}

// GetTVShows retrieves TV shows from the Emby server
func (e *EmbyClient) GetTVShows(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.TVShow, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.ClientID).
		Str("clientType", string(e.ClientType)).
		Str("baseURL", e.config.BaseURL).
		Msg("Retrieving TV shows from Emby server")

	queryParams := embyclient.ItemsServiceApiGetItemsOpts{
		IncludeItemTypes: optional.NewString("Series"),
		Recursive:        optional.NewBool(true),
	}

	if options != nil {
		if options.Limit > 0 {
			queryParams.Limit = optional.NewInt32(int32(options.Limit))
		}
		if options.Offset > 0 {
			queryParams.StartIndex = optional.NewInt32(int32(options.Offset))
		}
		if options.Sort != "" {
			queryParams.SortBy = optional.NewString(options.Sort)
			if options.SortOrder == "desc" {
				queryParams.SortOrder = optional.NewString("Descending")
			} else {
				queryParams.SortOrder = optional.NewString("Ascending")
			}
		}

		log.Debug().
			Int("limit", options.Limit).
			Int("offset", options.Offset).
			Str("sort", options.Sort).
			Str("sortOrder", options.SortOrder).
			Msg("Applied query options")
	}

	log.Debug().Msg("Making API request to Emby server for TV shows")
	items, resp, err := e.client.ItemsServiceApi.GetItems(ctx, &queryParams)
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", e.config.BaseURL).
			Str("apiEndpoint", "/Items").
			Int("statusCode", 0).
			Msg("Failed to fetch TV shows from Emby")
		return nil, fmt.Errorf("failed to fetch TV shows: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(items.Items)).
		Int("totalRecordCount", int(items.TotalRecordCount)).
		Msg("Successfully retrieved TV shows from Emby")

	shows := make([]interfaces.TVShow, 0)
	for _, item := range items.Items {
		if item.Type_ == "Series" {
			show := interfaces.TVShow{
				MediaItem: interfaces.MediaItem{
					Metadata: interfaces.MediaMetadata{
						Title:       item.Name,
						Description: item.Overview,
						ReleaseYear: int(item.ProductionYear),
						Genres:      item.Genres,
						Artwork:     e.getArtworkURLs(&item),
						Duration:    time.Duration(item.RunTimeTicks/10000000) * time.Second,
					},
					ExternalID: item.Id,
					Type:       "tvshow",
					ClientID:   e.ClientID,
					ClientType: string(e.ClientType),
				},
				SeasonCount: int(item.ChildCount),
				Status:      item.Status,
				Network:     item.SeriesStudio,
			}

			// Extract provider IDs if available
			if item.ProviderIds != nil {
				ids := *item.ProviderIds
				if imdbID, ok := ids["Imdb"]; ok {
					show.Metadata.ExternalIDs.AddOrUpdate("imdb", imdbID)
				}
				if tmdbID, ok := ids["Tmdb"]; ok {
					show.Metadata.ExternalIDs.AddOrUpdate("tmdb", tmdbID)
				}
				if tvdbID, ok := ids["Tvdb"]; ok {
					show.Metadata.ExternalIDs.AddOrUpdate("tvdb", tvdbID)
				}
			}

			shows = append(shows, show)
		}
	}

	log.Info().
		Int("showsReturned", len(shows)).
		Msg("Completed GetTVShows request")

	return shows, nil
}

// GetTVShowByID retrieves a specific TV show by ID
func (e *EmbyClient) GetTVShowByID(ctx context.Context, id string) (interfaces.TVShow, error) {
	items, _, err := e.client.ItemsServiceApi.GetItems(ctx, &embyclient.ItemsServiceApiGetItemsOpts{Ids: optional.NewString(id)})
	if err != nil {
		return interfaces.TVShow{}, fmt.Errorf("failed to fetch TV show: %w", err)
	}

	if len(items.Items) == 0 {
		return interfaces.TVShow{}, fmt.Errorf("TV show with ID %s not found", id)
	}

	item := items.Items[0]
	if item.Type_ != "Series" {
		return interfaces.TVShow{}, fmt.Errorf("item with ID %s is not a TV show", id)
	}

	show := interfaces.TVShow{
		MediaItem: interfaces.MediaItem{
			Metadata: interfaces.MediaMetadata{
				Title:       item.Name,
				Description: item.Overview,
				ReleaseYear: int(item.ProductionYear),
				Genres:      item.Genres,
				Artwork:     e.getArtworkURLs(&item),
				Duration:    time.Duration(item.RunTimeTicks/10000000) * time.Second,
			},
			ExternalID: item.Id,
			Type:       "tvshow",
			ClientID:   e.ClientID,
			ClientType: string(e.ClientType),
		},
		SeasonCount: int(item.ChildCount),
		Status:      item.Status,
		Network:     item.SeriesStudio,
	}

	// Extract provider IDs if available
	if item.ProviderIds != nil {
		ids := *item.ProviderIds
		if imdbID, ok := ids["Imdb"]; ok {
			show.Metadata.ExternalIDs.AddOrUpdate("imdb", imdbID)
		}
		if tmdbID, ok := ids["Tmdb"]; ok {
			show.Metadata.ExternalIDs.AddOrUpdate("tmdb", tmdbID)
		}
		if tvdbID, ok := ids["Tvdb"]; ok {
			show.Metadata.ExternalIDs.AddOrUpdate("tvdb", tvdbID)
		}
	}

	return show, nil
}

// GetTVShowSeasons retrieves seasons for a TV show
func (e *EmbyClient) GetTVShowSeasons(ctx context.Context, showID string) ([]interfaces.Season, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.ClientID).
		Str("clientType", string(e.ClientType)).
		Str("showID", showID).
		Str("baseURL", e.config.BaseURL).
		Msg("Retrieving seasons for TV show from Emby server")

	// Create query parameters
	opts := embyclient.TvShowsServiceApiGetShowsByIdSeasonsOpts{
		EnableImages:   optional.NewBool(true),
		EnableUserData: optional.NewBool(true),
	}

	log.Debug().
		Str("showID", showID).
		Bool("enableImages", true).
		Bool("enableUserData", true).
		Msg("Making API request to Emby server for TV show seasons")

	// Call the Emby API
	result, resp, err := e.client.TvShowsServiceApi.GetShowsByIdSeasons(ctx, e.config.UserID, showID, &opts)
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", e.config.BaseURL).
			Str("apiEndpoint", "/Shows/"+showID+"/Seasons").
			Str("showID", showID).
			Int("statusCode", 0).
			Msg("Failed to fetch seasons for TV show from Emby")
		return nil, fmt.Errorf("failed to fetch seasons: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("seasonCount", len(result.Items)).
		Str("showID", showID).
		Msg("Successfully retrieved seasons for TV show from Emby")

	seasons := make([]interfaces.Season, 0)
	for _, item := range result.Items {
		if item.Type_ == "Season" {
			season := interfaces.Season{
				MediaItem: interfaces.MediaItem{
					Metadata: interfaces.MediaMetadata{
						Title:       item.Name,
						Description: item.Overview,
						Artwork:     e.getArtworkURLs(&item),
					},
					Type:       "season",
					ClientID:   e.ClientID,
					ClientType: string(e.ClientType),
					ExternalID: item.Id,
				},
				ParentID:     showID,
				Number:       int(item.IndexNumber),
				EpisodeCount: int(item.ChildCount),
			}

			if !item.PremiereDate.IsZero() {
				season.ReleaseDate = item.PremiereDate
			}

			seasons = append(seasons, season)
		}
	}

	log.Info().
		Int("seasonsReturned", len(seasons)).
		Str("showID", showID).
		Msg("Completed GetTVShowSeasons request")

	return seasons, nil
}

// GetTVShowEpisodes retrieves episodes for a season
func (e *EmbyClient) GetTVShowEpisodes(ctx context.Context, showID string, seasonNumber int) ([]interfaces.Episode, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", e.ClientID).
		Str("clientType", string(e.ClientType)).
		Str("showID", showID).
		Int("seasonNumber", seasonNumber).
		Str("baseURL", e.config.BaseURL).
		Msg("Retrieving episodes for TV show season from Emby server")

	// Need to fetch the items from ItemsApi since the ShowsApi doesn't return the full data
	queryParams := embyclient.TvShowsServiceApiGetShowsByIdEpisodesOpts{
		// ParentId:         optional.NewString(showID),
		IncludeItemTypes: optional.NewString("Episode"),
		Recursive:        optional.NewBool(true),
	}

	// queryParams.Filters = optional.NewString("ParentIndexNumber=" + strconv.Itoa(seasonNumber))

	log.Debug().
		Str("showID", showID).
		Int("seasonNumber", seasonNumber).
		Msg("Making API request to Emby server for TV show episodes")

	items, _, err := e.client.TvShowsServiceApi.GetShowsByIdEpisodes(ctx, showID, &queryParams)

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", e.config.BaseURL).
			Str("apiEndpoint", "/Items").
			Str("showID", showID).
			Int("seasonNumber", seasonNumber).
			Int("statusCode", 0).
			Msg("Failed to fetch episodes for TV show season from Emby")
		return nil, fmt.Errorf("failed to fetch episodes: %w", err)
	}

	log.Info().
		// Int("statusCode", resp.StatusCode).
		Int("episodeCount", len(items.Items)).
		Str("showID", showID).
		Int("seasonNumber", seasonNumber).
		Msg("Successfully retrieved episodes for TV show season from Emby")

	episodes := make([]interfaces.Episode, 0)
	for _, item := range items.Items {
		if item.Type_ == "Episode" {
			episode := interfaces.Episode{
				MediaItem: interfaces.MediaItem{
					Metadata: interfaces.MediaMetadata{
						Title:       item.Name,
						Description: item.Overview,
						Artwork:     e.getArtworkURLs(&item),
						Duration:    time.Duration(item.RunTimeTicks/10000000) * time.Second,
					},
					Type:       "episode",
					ClientID:   e.ClientID,
					ClientType: string(e.ClientType),
					ExternalID: item.Id,
				},
				Number:       int64(item.IndexNumber),
				ShowID:       showID,
				SeasonID:     item.SeasonId,
				SeasonNumber: seasonNumber,
				ShowTitle:    item.SeriesName,
			}

			// Add external IDs
			if item.ProviderIds != nil {
				ids := *item.ProviderIds
				if imdbID, ok := ids["Imdb"]; ok {
					episode.Metadata.ExternalIDs.AddOrUpdate("imdb", imdbID)
				}
				if tmdbID, ok := ids["Tmdb"]; ok {
					episode.Metadata.ExternalIDs.AddOrUpdate("tmdb", tmdbID)
				}
				if tvdbID, ok := ids["Tvdb"]; ok {
					episode.Metadata.ExternalIDs.AddOrUpdate("tvdb", tvdbID)
				}
			}

			episodes = append(episodes, episode)
		}
	}

	log.Info().
		Int("episodesReturned", len(episodes)).
		Str("showID", showID).
		Int("seasonNumber", seasonNumber).
		Msg("Completed GetTVShowEpisodes request")

	return episodes, nil
}

// GetEpisodeByID retrieves a specific episode by ID
func (e *EmbyClient) GetEpisodeByID(ctx context.Context, id string) (interfaces.Episode, error) {
	items, _, err := e.client.ItemsServiceApi.GetItems(ctx, &embyclient.ItemsServiceApiGetItemsOpts{Ids: optional.NewString(id)})
	if err != nil {
		return interfaces.Episode{}, fmt.Errorf("failed to fetch episode: %w", err)
	}

	if len(items.Items) == 0 {
		return interfaces.Episode{}, fmt.Errorf("episode with ID %s not found", id)
	}

	item := items.Items[0]
	if item.Type_ != "Episode" {
		return interfaces.Episode{}, fmt.Errorf("item with ID %s is not an episode", id)
	}

	episode := interfaces.Episode{
		MediaItem: interfaces.MediaItem{
			Metadata: interfaces.MediaMetadata{
				Title:       item.Name,
				Description: item.Overview,
				Artwork:     e.getArtworkURLs(&item),
				Duration:    time.Duration(item.RunTimeTicks/10000000) * time.Second,
			},
			Type:       "episode",
			ClientID:   e.ClientID,
			ClientType: string(e.ClientType),
			ExternalID: item.Id,
		},
		Number:       int64(item.IndexNumber),
		ShowID:       item.SeriesId,
		SeasonID:     item.SeasonId,
		SeasonNumber: int(item.ParentIndexNumber),
		ShowTitle:    item.SeriesName,
	}

	// Add external IDs
	if item.ProviderIds != nil {
		ids := *item.ProviderIds
		if imdbID, ok := ids["Imdb"]; ok {
			episode.Metadata.ExternalIDs.AddOrUpdate("imdb", imdbID)
		}
		if tmdbID, ok := ids["Tmdb"]; ok {
			episode.Metadata.ExternalIDs.AddOrUpdate("tmdb", tmdbID)
		}
		if tvdbID, ok := ids["Tvdb"]; ok {
			episode.Metadata.ExternalIDs.AddOrUpdate("tvdb", tvdbID)
		}
	}

	return episode, nil
}

// GetMusic retrieves music tracks from the Emby server
func (e *EmbyClient) GetMusic(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.MusicTrack, error) {
	queryParams := embyclient.ItemsServiceApiGetItemsOpts{
		IncludeItemTypes: optional.NewString("Audio"),
		Recursive:        optional.NewBool(true),
	}

	// if e.config.UserID != "" {
	// 	queryParams.UserId = optional.NewString(e.config.UserID)
	// }

	if options != nil {
		if options.Limit > 0 {
			queryParams.Limit = optional.NewInt32(int32(options.Limit))
		}
		if options.Offset > 0 {
			queryParams.StartIndex = optional.NewInt32(int32(options.Offset))
		}
		if options.Sort != "" {
			queryParams.SortBy = optional.NewString(options.Sort)
			if options.SortOrder == "desc" {
				queryParams.SortOrder = optional.NewString("Descending")
			} else {
				queryParams.SortOrder = optional.NewString("Ascending")
			}
		}
	}

	items, _, err := e.client.ItemsServiceApi.GetItems(ctx, &queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch music tracks: %w", err)
	}

	tracks := make([]interfaces.MusicTrack, 0)
	for _, item := range items.Items {
		track := interfaces.MusicTrack{
			MediaItem: interfaces.MediaItem{
				Metadata: interfaces.MediaMetadata{
					Title:       item.Name,
					Description: item.Overview,
					Duration:    time.Duration(item.RunTimeTicks/10000000) * time.Second,
					Artwork:     e.getArtworkURLs(&item),
				},
				ExternalID: item.Id,
				Type:       "track",
				ClientID:   e.ClientID,
				ClientType: string(e.ClientType),
			},
			Number:    int(item.IndexNumber),
			AlbumID:   item.AlbumId,
			AlbumName: item.Album,
		}

		// Add artist information if available
		if len(item.ArtistItems) > 0 {
			track.ArtistID = item.ArtistItems[0].Id
			track.ArtistName = item.ArtistItems[0].Name
		}

		// Extract provider IDs
		if item.ProviderIds != nil {
			ids := *item.ProviderIds
			if musicbrainzID, ok := ids["MusicBrainzTrack"]; ok {
				track.Metadata.ExternalIDs.AddOrUpdate("musicbrainz", musicbrainzID)
			}
		}

		tracks = append(tracks, track)
	}

	return tracks, nil
}

// GetMusicArtists retrieves music artists from the Emby server
func (e *EmbyClient) GetMusicArtists(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.MusicArtist, error) {
	opts := embyclient.ArtistsServiceApiGetArtistsOpts{}
	opts.Recursive = optional.NewBool(true)

	// if e.config.UserID != "" {
	// 	opts.UserId = optional.NewString(e.config.UserID)
	// }

	if options != nil {
		if options.Limit > 0 {
			opts.Limit = optional.NewInt32(int32(options.Limit))
		}
		if options.Offset > 0 {
			opts.StartIndex = optional.NewInt32(int32(options.Offset))
		}
		if options.Sort != "" {
			opts.SortBy = optional.NewString(options.Sort)
			if options.SortOrder == "desc" {
				opts.SortOrder = optional.NewString("Descending")
			} else {
				opts.SortOrder = optional.NewString("Ascending")
			}
		}
	}

	result, _, err := e.client.ArtistsServiceApi.GetArtists(ctx, &opts)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch music artists: %w", err)
	}

	artists := make([]interfaces.MusicArtist, 0)
	for _, item := range result.Items {
		artist := interfaces.MusicArtist{
			MediaItem: interfaces.MediaItem{
				Metadata: interfaces.MediaMetadata{
					Title:       item.Name,
					Description: item.Overview,
					Artwork:     e.getArtworkURLs(&item),
					Genres:      item.Genres,
				},
				ExternalID: item.Id,
				Type:       "artist",
				ClientID:   e.ClientID,
				ClientType: string(e.ClientType),
			},
		}

		// Extract provider IDs if available
		if item.ProviderIds != nil {
			ids := *item.ProviderIds
			if musicbrainzID, ok := ids["MusicBrainzArtist"]; ok {
				artist.Metadata.ExternalIDs.AddOrUpdate("musicbrainz", musicbrainzID)
			}
		}

		artists = append(artists, artist)
	}

	return artists, nil
}

// GetMusicAlbums retrieves music albums from the Emby server
func (e *EmbyClient) GetMusicAlbums(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.MusicAlbum, error) {
	queryParams := embyclient.ItemsServiceApiGetItemsOpts{
		IncludeItemTypes: optional.NewString("MusicAlbum"),
		Recursive:        optional.NewBool(true),
	}
	//
	// if e.config.UserID != "" {
	// 	queryParams.UserId = optional.NewString(e.config.UserID)
	// }

	if options != nil {
		if options.Limit > 0 {
			queryParams.Limit = optional.NewInt32(int32(options.Limit))
		}
		if options.Offset > 0 {
			queryParams.StartIndex = optional.NewInt32(int32(options.Offset))
		}
		if options.Sort != "" {
			queryParams.SortBy = optional.NewString(options.Sort)
			if options.SortOrder == "desc" {
				queryParams.SortOrder = optional.NewString("Descending")
			} else {
				queryParams.SortOrder = optional.NewString("Ascending")
			}
		}
	}

	items, _, err := e.client.ItemsServiceApi.GetItems(ctx, &queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch music albums: %w", err)
	}

	albums := make([]interfaces.MusicAlbum, 0)
	for _, item := range items.Items {
		album := interfaces.MusicAlbum{
			MediaItem: interfaces.MediaItem{
				Metadata: interfaces.MediaMetadata{
					Title:       item.Name,
					Description: item.Overview,
					ReleaseYear: int(item.ProductionYear),
					Genres:      item.Genres,
					Artwork:     e.getArtworkURLs(&item),
				},
				Type:       "album",
				ExternalID: item.Id,
				ClientID:   e.ClientID,
				ClientType: string(e.ClientType),
			},
			ArtistName: item.AlbumArtist,
			TrackCount: int(item.ChildCount),
		}

		// Extract provider IDs
		if item.ProviderIds != nil {
			ids := *item.ProviderIds
			if musicbrainzID, ok := ids["MusicBrainzAlbum"]; ok {
				album.Metadata.ExternalIDs.AddOrUpdate("musicbrainz", musicbrainzID)
			}
		}

		albums = append(albums, album)
	}

	return albums, nil
}

// GetMusicTrackByID retrieves a specific music track by ID
func (e *EmbyClient) GetMusicTrackByID(ctx context.Context, id string) (interfaces.MusicTrack, error) {
	items, _, err := e.client.ItemsServiceApi.GetItems(ctx, &embyclient.ItemsServiceApiGetItemsOpts{Ids: optional.NewString(id)})
	if err != nil {
		return interfaces.MusicTrack{}, fmt.Errorf("failed to fetch music track: %w", err)
	}

	if len(items.Items) == 0 {
		return interfaces.MusicTrack{}, fmt.Errorf("music track with ID %s not found", id)
	}

	item := items.Items[0]
	if item.Type_ != "Audio" {
		return interfaces.MusicTrack{}, fmt.Errorf("item with ID %s is not a music track", id)
	}

	track := interfaces.MusicTrack{
		MediaItem: interfaces.MediaItem{
			Metadata: interfaces.MediaMetadata{
				Title:       item.Name,
				Description: item.Overview,
				Duration:    time.Duration(item.RunTimeTicks/10000000) * time.Second,
				Artwork:     e.getArtworkURLs(&item),
			},
			ExternalID: item.Id,
			Type:       "track",
			ClientID:   e.ClientID,
			ClientType: string(e.ClientType),
		},
		Number:    int(item.IndexNumber),
		AlbumID:   item.AlbumId,
		AlbumName: item.Album,
	}

	// Add artist information if available
	if len(item.ArtistItems) > 0 {
		track.ArtistID = item.ArtistItems[0].Id
		track.ArtistName = item.ArtistItems[0].Name
	}

	// Extract provider IDs
	if item.ProviderIds != nil {
		ids := *item.ProviderIds
		if musicbrainzID, ok := ids["MusicBrainzTrack"]; ok {
			track.Metadata.ExternalIDs.AddOrUpdate("musicbrainz", musicbrainzID)
		}
	}

	return track, nil
}

// GetPlaylists retrieves playlists from the Emby server
func (e *EmbyClient) GetPlaylists(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.Playlist, error) {
	queryParams := embyclient.ItemsServiceApiGetItemsOpts{
		IncludeItemTypes: optional.NewString("Playlist"),
		Recursive:        optional.NewBool(true),
	}

	// if e.config.UserID != "" {
	// 	queryParams.UserId = optional.NewString(e.config.UserID)
	// }

	if options != nil {
		if options.Limit > 0 {
			queryParams.Limit = optional.NewInt32(int32(options.Limit))
		}
		if options.Offset > 0 {
			queryParams.StartIndex = optional.NewInt32(int32(options.Offset))
		}
	}

	items, _, err := e.client.ItemsServiceApi.GetItems(ctx, &queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch playlists: %w", err)
	}

	playlists := make([]interfaces.Playlist, 0)
	for _, item := range items.Items {
		if item.Type_ == "Playlist" {
			playlist := interfaces.Playlist{
				MediaItem: interfaces.MediaItem{
					Metadata: interfaces.MediaMetadata{
						Title:       item.Name,
						Description: item.Overview,
						Artwork:     e.getArtworkURLs(&item),
					},
					ExternalID: item.Id,
					Type:       "playlist",
					ClientID:   e.ClientID,
					ClientType: string(e.ClientType),
				},
				ItemCount: int(item.ChildCount),
				// Owner:     item.UserData.Key,
				IsPublic: true, // Assume public by default in Emby
			}
			playlists = append(playlists, playlist)
		}
	}

	return playlists, nil
}

// GetWatchHistory retrieves watch history from the Emby server

func (e *EmbyClient) GetWatchHistory(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.WatchHistoryItem, error) {
	queryParams := embyclient.ItemsServiceApiGetUsersByUseridItemsOpts{
		IsPlayed:  optional.NewBool(true),
		Recursive: optional.NewBool(true),
	}

	// if e.config.UserID != "" {
	// 	queryParams.UserId = optional.NewString(e.config.UserID)
	// }

	// Apply options for pagination, etc.
	if options != nil {
		if options.Limit > 0 {
			queryParams.Limit = optional.NewInt32(int32(options.Limit))
		}
		if options.Offset > 0 {
			queryParams.StartIndex = optional.NewInt32(int32(options.Offset))
		}
	}

	items, _, err := e.client.ItemsServiceApi.GetUsersByUseridItems(ctx, e.config.UserID, &queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch watch history: %w", err)
	}

	history := make([]interfaces.WatchHistoryItem, 0)
	for _, item := range items.Items {
		watchItem := interfaces.WatchHistoryItem{
			MediaItem: interfaces.MediaItem{
				Metadata: interfaces.MediaMetadata{
					Title: item.Name,
				},
				// ID:         item.Id,
				ClientID:   e.ClientID,
				ClientType: string(e.ClientType),
			},
			ItemType:        string(item.Type_),
			WatchedAt:       item.UserData.LastPlayedDate,
			IsFavorite:      item.UserData.IsFavorite,
			PlayCount:       item.UserData.PlayCount,
			PositionSeconds: int(item.UserData.PlaybackPositionTicks),
		}
		history = append(history, watchItem)
	}

	return history, nil
}

// GetMusicGenres retrieves music genres from the Emby server

// GetMusicGenres retrieves music genres from the Emby server
func (e *EmbyClient) GetMusicGenres(ctx context.Context) ([]string, error) {
	opts := embyclient.MusicGenresServiceApiGetMusicgenresOpts{}
	// if e.config.UserID != "" {
	// 	opts.UserId = optional.NewString(e.config.UserID)
	// }

	result, _, err := e.client.MusicGenresServiceApi.GetMusicgenres(ctx, &opts)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch music genres: %w", err)
	}

	genres := make([]string, 0, len(result.Items))
	for _, item := range result.Items {
		genres = append(genres, item.Name)
	}

	return genres, nil
}

// GetMovieGenres retrieves movie genres from the Emby server
func (e *EmbyClient) GetMovieGenres(ctx context.Context) ([]string, error) {
	opts := embyclient.GenresServiceApiGetGenresOpts{IsMovie: optional.NewBool(true)}

	result, _, err := e.client.GenresServiceApi.GetGenres(ctx, &opts)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch movie genres: %w", err)
	}

	genres := make([]string, 0, len(result.Items))
	for _, item := range result.Items {
		genres = append(genres, item.Name)
	}

	return genres, nil
}
