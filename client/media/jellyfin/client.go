package jellyfin

import (
	"context"
	"fmt"
	"strings"
	"time"

	jellyfin "github.com/sj14/jellyfin-go/api"
	"suasor/client/media/interfaces"
	"suasor/models"
	"suasor/utils"
)

// Configuration holds Jellyfin connection settings
type Configuration struct {
	BaseURL  string
	ApiKey   string
	Username string
	UserID   string
}

// JellyfinClient implements the MediaContentProvider interface
type JellyfinClient struct {
	interfaces.BaseMediaClient
	client *jellyfin.APIClient
	config Configuration
}

// Convert a single string to a string slice
func stringToSlice(s string) []string {
	if s == "" {
		return []string{}
	}
	return []string{s}
}

// Safely convert BaseItemKind to string
func baseItemKindToString(kind jellyfin.BaseItemKind) string {
	return string(kind)
}

// extractProviderIDs adds external IDs from the Jellyfin provider IDs map to the metadata
func extractProviderIDs(providerIds *map[string]string, externalIDs *interfaces.ExternalIDs) {
	if providerIds == nil {
		return
	}

	// Common media identifier mappings
	idMappings := map[string]string{
		"Imdb":              "imdb",
		"Tmdb":              "tmdb",
		"Tvdb":              "tvdb",
		"MusicBrainzTrack":  "musicbrainz",
		"MusicBrainzAlbum":  "musicbrainz",
		"MusicBrainzArtist": "musicbrainz",
	}

	// Extract all available IDs based on the mappings
	for jellyfinKey, externalKey := range idMappings {
		if id, ok := (*providerIds)[jellyfinKey]; ok {
			externalIDs.AddOrUpdate(externalKey, id)
		}
	}

}

// getQueryParameters extracts common query parameters from QueryOptions
// and converts them to the format expected by the Jellyfin API
func (j *JellyfinClient) getQueryParameters(options *interfaces.QueryOptions) (limit, startIndex *int32, sortBy []jellyfin.ItemSortBy, sortOrder []jellyfin.SortOrder) {

	// Default values
	defaultLimit := int32(100)
	defaultOffset := int32(0)
	limit = &defaultLimit
	startIndex = &defaultOffset

	if options != nil {
		if options.Limit > 0 {
			limitVal := int32(options.Limit)
			limit = &limitVal
		}
		if options.Offset > 0 {
			offsetVal := int32(options.Offset)
			startIndex = &offsetVal
		}
		if options.Sort != "" {
			// sortBy = &options.Sort
			sortBy = []jellyfin.ItemSortBy{jellyfin.ItemSortBy(options.Sort)}
			if options.SortOrder == "desc" {
				sortOrder = []jellyfin.SortOrder{jellyfin.SORTORDER_DESCENDING}
			} else {
				sortOrder = []jellyfin.SortOrder{jellyfin.SORTORDER_ASCENDING}
			}
		}
	}
	return
}

// Helper function to convert Jellyfin item to internal Collection type
func (j *JellyfinClient) convertToCollection(ctx context.Context, item *jellyfin.BaseItemDto) (interfaces.Collection, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	// Validate required fields
	if item == nil {
		return interfaces.Collection{}, fmt.Errorf("cannot convert nil item to collection")
	}

	if item.Id == nil || *item.Id == "" {
		return interfaces.Collection{}, fmt.Errorf("collection is missing required ID field")
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
	collection := interfaces.Collection{
		MediaItem: interfaces.MediaItem{
			Metadata: interfaces.MediaMetadata{
				Title:       title,
				Description: description,
				Artwork:     j.getArtworkURLs(item),
			},
			ExternalID: *item.Id,
			Type:       "collection",
			ClientID:   j.ClientID,
			ClientType: string(j.ClientType),
		},
		ItemCount: itemCount,
	}

	// Add potential year if available
	if item.ProductionYear.IsSet() {
		collection.Metadata.ReleaseYear = int(*item.ProductionYear.Get())
	}

	// Add community rating if available
	if item.CommunityRating.IsSet() {
		collection.Metadata.Ratings = append(collection.Metadata.Ratings, interfaces.Rating{
			Source: "jellyfin",
			Value:  float32(*item.CommunityRating.Get()),
		})
	}

	// Add user rating if available
	if item.UserData.IsSet() && item.UserData.Get().Rating.IsSet() {
		collection.Metadata.UserRating = float32(*item.UserData.Get().Rating.Get())
	}

	// Handle genres if available
	if item.Genres != nil {
		collection.Metadata.Genres = item.Genres
	}

	// Extract provider IDs if available
	extractProviderIDs(&item.ProviderIds, &collection.Metadata.ExternalIDs)

	log.Debug().
		Str("collectionID", *item.Id).
		Str("collectionName", collection.Metadata.Title).
		Int("itemCount", collection.ItemCount).
		Msg("Successfully converted Jellyfin item to collection")

	return collection, nil
}

// Helper function to convert Jellyfin item to internal Episode type
func (j *JellyfinClient) convertToEpisode(ctx context.Context, item *jellyfin.BaseItemDto) (interfaces.Episode, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	// Validate required fields
	if item == nil {
		return interfaces.Episode{}, fmt.Errorf("cannot convert nil item to episode")
	}

	if item.Id == nil || *item.Id == "" {
		return interfaces.Episode{}, fmt.Errorf("episode is missing required ID field")
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
	episode := interfaces.Episode{
		MediaItem: interfaces.MediaItem{
			Metadata: interfaces.MediaMetadata{
				Title:       title,
				Description: description,
				Artwork:     j.getArtworkURLs(item),
				Duration:    duration,
			},
			Type:       "episode",
			ClientID:   j.ClientID,
			ClientType: string(j.ClientType),
			ExternalID: *item.Id,
		},
		Number:       episodeNumber,
		SeasonNumber: seasonNumber,
		ShowTitle:    showTitle,
	}

	// Safely set IDs if available
	if item.SeriesId.IsSet() {
		episode.ShowID = *item.SeriesId.Get()
	}

	if item.SeasonId.IsSet() {
		episode.SeasonID = *item.SeasonId.Get()
	}

	// Add air date if available
	if item.PremiereDate.IsSet() {
		episode.Metadata.ReleaseDate = *item.PremiereDate.Get()
	}

	// Add community rating if available
	if item.CommunityRating.IsSet() {
		episode.Metadata.Ratings = append(episode.Metadata.Ratings, interfaces.Rating{
			Source: "jellyfin",
			Value:  float32(*item.CommunityRating.Get()),
		})
	}

	// Add user rating if available
	if item.UserData.IsSet() && item.UserData.Get().Rating.IsSet() {
		episode.Metadata.UserRating = float32(*item.UserData.Get().Rating.Get())
	}

	// Extract provider IDs if available
	extractProviderIDs(&item.ProviderIds, &episode.Metadata.ExternalIDs)

	log.Debug().
		Str("episodeID", *item.Id).
		Str("episodeName", episode.Metadata.Title).
		Int64("episodeNumber", episode.Number).
		Int("seasonNumber", episode.SeasonNumber).
		Msg("Successfully converted Jellyfin item to episode")

	return episode, nil
}

// NewJellyfinClient creates a new Jellyfin client instance
func NewJellyfinClient(ctx context.Context, clientID uint64, config interface{}) (interfaces.MediaContentProvider, error) {
	cfg, ok := config.(Configuration)
	if !ok {
		return nil, fmt.Errorf("invalid configuration for Jellyfin client")
	}

	// Create API client configuration
	apiConfig := &jellyfin.Configuration{
		Servers:       jellyfin.ServerConfigurations{{URL: cfg.BaseURL}},
		DefaultHeader: map[string]string{"Authorization": fmt.Sprintf(`MediaBrowser Token="%s"`, cfg.ApiKey)},
	}

	client := jellyfin.NewAPIClient(apiConfig)

	jellyfinClient := &JellyfinClient{
		BaseMediaClient: interfaces.BaseMediaClient{
			ClientID:   clientID,
			ClientType: models.MediaClientTypeJellyfin,
		},
		client: client,
		config: cfg,
	}

	// Resolve user ID if username is provided
	if cfg.Username != "" && cfg.UserID == "" {
		if err := jellyfinClient.resolveUserID(ctx); err != nil {
			// Log but don't fail - some operations might work without a user ID
			log := utils.LoggerFromContext(ctx)
			log.Warn().
				Err(err).
				Str("username", cfg.Username).
				Msg("Failed to resolve Jellyfin user ID, some operations may be limited")
		}
	}
	return jellyfinClient, nil
}

// Register the provider factory
func init() {
	interfaces.RegisterProvider(models.MediaClientTypeJellyfin, NewJellyfinClient)
}

// Capability methods
func (j *JellyfinClient) SupportsMovies() bool      { return true }
func (j *JellyfinClient) SupportsTVShows() bool     { return true }
func (j *JellyfinClient) SupportsMusic() bool       { return true }
func (j *JellyfinClient) SupportsPlaylists() bool   { return true }
func (j *JellyfinClient) SupportsCollections() bool { return true }

// resolveUserID resolves the user ID from the username
func (j *JellyfinClient) resolveUserID(ctx context.Context) error {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Str("username", j.config.Username).
		Msg("Resolving Jellyfin user ID from username")

		// Get the list of public users
	publicUsersReq := j.client.UserAPI.GetUsers(ctx)
	users, resp, err := publicUsersReq.Execute()

	if err != nil {
		log.Error().
			Err(err).
			Str("username", j.config.Username).
			Msg("Failed to fetch Jellyfin users")
		return fmt.Errorf("failed to fetch users: %w", err)
	}

	log.Debug().
		Int("statusCode", resp.StatusCode).
		Int("userCount", len(users)).
		Msg("Retrieved public users from Jellyfin")

	// Find the user with matching username
	for _, user := range users {
		if user.Name.IsSet() {
			if strings.EqualFold(*user.Name.Get(), j.config.Username) {
				j.config.UserID = *user.Id
				log.Info().
					Str("username", j.config.Username).
					Str("userID", j.config.UserID).
					Msg("Successfully resolved Jellyfin user ID")
				return nil
			}
		}
	}

	log.Warn().
		Str("username", j.config.Username).
		Msg("Could not find matching user in Jellyfin")
	return fmt.Errorf("user '%s' not found in Jellyfin", j.config.Username)
}

// GetMovies retrieves movies from the Jellyfin server
func (j *JellyfinClient) GetMovies(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.Movie, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.ClientID).
		Str("clientType", string(j.ClientType)).
		Str("baseURL", j.config.BaseURL).
		Msg("Retrieving movies from Jellyfin server")

	// Set up query parameters

	limit, startIndex, sortBy, sortOrder := j.getQueryParameters(options)

	// Include movie type in the query
	includeItemTypes := []jellyfin.BaseItemKind{jellyfin.BASEITEMKIND_MOVIE}
	mediaTypes := []jellyfin.MediaType{jellyfin.MEDIATYPE_VIDEO}
	// Call the Jellyfin API
	log.Debug().Msg("Making API request to Jellyfin server")
	itemsReq := j.client.ItemsAPI.GetItems(ctx).
		IncludeItemTypes(includeItemTypes).
		IsMovie(true).
		Recursive(true).
		MediaTypes(mediaTypes).
		Limit(*limit).
		StartIndex(*startIndex).
		SortBy(sortBy).
		SortOrder(sortOrder)

	log.Debug().
		Int32("Limit", *limit).
		Int32("StartIndex", *startIndex).
		Str("IncludeItemTypes", baseItemKindToString(includeItemTypes[0])).
		Bool("Recursive", true).
		Msg("Api Request with options")

	result, resp, err := itemsReq.Execute()

	log.Debug().
		Interface("responseItems", result.Items).
		Msg("Full response data from Jellyfin API")

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.BaseURL).
			Str("apiEndpoint", "/Items").
			Int("statusCode", 0).
			Msg("Failed to fetch movies from Jellyfin")
		return nil, fmt.Errorf("failed to fetch movies: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(result.Items)).
		Int("totalRecordCount", int(*result.TotalRecordCount)).
		Msg("Successfully retrieved movies from Jellyfin")

	// Convert results to expected format
	movies := make([]interfaces.Movie, 0)

	for _, item := range result.Items {
		log.Info().
			Str("itemType", baseItemKindToString(*item.Type)).
			Msg("Processing item")
		if *item.Type == jellyfin.BASEITEMKIND_MOVIE {
			movie, err := j.convertToMovie(ctx, &item)
			if err != nil {
				// Log error but continue
				log.Warn().
					Err(err).
					Str("movieID", *item.Id).
					Str("movieName", *item.Name.Get()).
					Msg("Error converting Jellyfin item to movie format")
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

// Helper function to convert Jellyfin item to internal TVShow type
func (j *JellyfinClient) convertToTVShow(ctx context.Context, item *jellyfin.BaseItemDto) (interfaces.TVShow, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	// Validate required fields
	if item == nil {
		return interfaces.TVShow{}, fmt.Errorf("cannot convert nil item to TV show")
	}

	if item.Id == nil || *item.Id == "" {
		return interfaces.TVShow{}, fmt.Errorf("TV show is missing required ID field")
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
	show := interfaces.TVShow{
		MediaItem: interfaces.MediaItem{
			Metadata: interfaces.MediaMetadata{
				Title:       title,
				Description: description,
				ReleaseYear: releaseYear,
				Genres:      genres,
				Artwork:     j.getArtworkURLs(item),
				Duration:    duration,
			},
			ExternalID: *item.Id,
			Type:       "tvshow",
			ClientID:   j.ClientID,
			ClientType: string(j.ClientType),
		},
		SeasonCount: seasonCount,
		Status:      status,
	}

	// Set SeriesStudio if available
	if item.SeriesStudio.IsSet() {
		show.Network = *item.SeriesStudio.Get()
	}

	// Extract provider IDs if available
	extractProviderIDs(&item.ProviderIds, &show.Metadata.ExternalIDs)

	// Set ratings if available
	if item.CommunityRating.IsSet() {
		show.Metadata.Ratings = append(show.Metadata.Ratings, interfaces.Rating{
			Source: "jellyfin",
			Value:  float32(*item.CommunityRating.Get()),
		})
	}

	// Set user rating if available
	if item.UserData.IsSet() && item.UserData.Get().Rating.IsSet() {
		show.Metadata.UserRating = float32(*item.UserData.Get().Rating.Get())
	}

	log.Debug().
		Str("showID", *item.Id).
		Str("showName", show.Metadata.Title).
		Int("seasonCount", show.SeasonCount).
		Msg("Successfully converted Jellyfin item to TV show")

	return show, nil
}

// Helper function to convert Jellyfin item to internal Movie type
func (j *JellyfinClient) convertToMovie(ctx context.Context, item *jellyfin.BaseItemDto) (interfaces.Movie, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	// Validate required fields
	if item == nil {
		return interfaces.Movie{}, fmt.Errorf("cannot convert nil item to movie")
	}

	if item.Id == nil || *item.Id == "" {
		return interfaces.Movie{}, fmt.Errorf("movie is missing required ID field")
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
	ratings := interfaces.Ratings{}

	// Safely add community rating if available
	if item.CommunityRating.IsSet() {
		ratings = append(ratings, interfaces.Rating{
			Source: "jellyfin",
			Value:  float32(*item.CommunityRating.Get()),
		})
	}

	// Build movie object
	movie := interfaces.Movie{
		MediaItem: interfaces.MediaItem{
			Metadata: interfaces.MediaMetadata{
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
			ExternalID: *item.Id,
			Type:       "movie",
			ClientID:   j.ClientID,
			ClientType: string(j.ClientType),
		},
	}

	// Set user rating if available
	if item.UserData.IsSet() && item.UserData.Get().Rating.IsSet() {
		movie.Metadata.UserRating = float32(*item.UserData.Get().Rating.Get())
	} else {
		log.Debug().
			Str("movieID", *item.Id).
			Msg("Movie has no user data, skipping user rating")
	}

	// Extract provider IDs if available
	extractProviderIDs(&item.ProviderIds, &movie.Metadata.ExternalIDs)

	log.Debug().
		Str("movieID", *item.Id).
		Str("movieTitle", movie.Metadata.Title).
		Int("year", movie.Metadata.ReleaseYear).
		Msg("Successfully converted Jellyfin item to movie")

	return movie, nil
}

// Helper to get a string value from a pointer, with a default empty string if nil
func getStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

// Helper to get artwork URLs for an item
func (j *JellyfinClient) getArtworkURLs(item *jellyfin.BaseItemDto) interfaces.Artwork {
	imageURLs := interfaces.Artwork{}

	if item == nil || item.Id == nil {
		return imageURLs
	}

	baseURL := strings.TrimSuffix(j.config.BaseURL, "/")
	itemID := *item.Id

	// Primary image (poster)
	if item.ImageTags != nil {
		if tag, ok := (item.ImageTags)["Primary"]; ok {
			imageURLs.Poster = fmt.Sprintf("%s/Items/%s/Images/Primary?tag=%s", baseURL, itemID, tag)
		}
	}

	// Backdrop image
	if item.BackdropImageTags != nil && len(item.BackdropImageTags) > 0 {
		imageURLs.Background = fmt.Sprintf("%s/Items/%s/Images/Backdrop?tag=%s", baseURL, itemID, item.BackdropImageTags[0])
	}

	// Other image types
	if item.ImageTags != nil {
		if tag, ok := (item.ImageTags)["Logo"]; ok {
			imageURLs.Logo = fmt.Sprintf("%s/Items/%s/Images/Logo?tag=%s", baseURL, itemID, tag)
		}

		if tag, ok := (item.ImageTags)["Thumb"]; ok {
			imageURLs.Thumbnail = fmt.Sprintf("%s/Items/%s/Images/Thumb?tag=%s", baseURL, itemID, tag)
		}

		if tag, ok := (item.ImageTags)["Banner"]; ok {
			imageURLs.Banner = fmt.Sprintf("%s/Items/%s/Images/Banner?tag=%s", baseURL, itemID, tag)
		}
	}

	return imageURLs
}

// GetMovieByID retrieves a specific movie by ID
func (j *JellyfinClient) GetMovieByID(ctx context.Context, id string) (interfaces.Movie, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.ClientID).
		Str("clientType", string(j.ClientType)).
		Str("movieID", id).
		Str("baseURL", j.config.BaseURL).
		Msg("Retrieving specific movie from Jellyfin server")

	// Set up query parameters
	includeItemTypes := []jellyfin.BaseItemKind{jellyfin.BASEITEMKIND_MOVIE}

	ids := id
	// fields := "ProductionYear,PremiereDate,ChannelMappingInfo,DateCreated,Genres,IndexOptions,HomePageUrl,Overview,ParentId,Path,ProviderIds,Studios,SortName"

	// Call the Jellyfin API
	log.Debug().
		Str("movieID", id).
		Msg("Making API request to Jellyfin server")

	fields := []jellyfin.ItemFields{
		jellyfin.ITEMFIELDS_DATE_CREATED,
		jellyfin.ITEMFIELDS_GENRES,
		jellyfin.ITEMFIELDS_PROVIDER_IDS,
		jellyfin.ITEMFIELDS_ORIGINAL_TITLE,
		jellyfin.ITEMFIELDS_AIR_TIME,
		jellyfin.ITEMFIELDS_EXTERNAL_URLS,
		jellyfin.ITEMFIELDS_STUDIOS,
	}

	itemsReq := j.client.ItemsAPI.GetItems(ctx).
		Ids(stringToSlice(ids)).
		IncludeItemTypes(includeItemTypes).
		Fields(fields)

	result, resp, err := itemsReq.Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.BaseURL).
			Str("apiEndpoint", "/Items").
			Str("movieID", id).
			Int("statusCode", 0).
			Msg("Failed to fetch movie from Jellyfin")
		return interfaces.Movie{}, fmt.Errorf("failed to fetch movie: %w", err)
	}

	// Check if any items were returned
	if len(result.Items) == 0 {
		log.Error().
			Str("movieID", id).
			Int("statusCode", resp.StatusCode).
			Msg("No movie found with the specified ID")
		return interfaces.Movie{}, fmt.Errorf("movie with ID %s not found", id)
	}

	item := result.Items[0]

	// Double-check that the returned item is a movie
	if *item.Type != jellyfin.BASEITEMKIND_MOVIE {
		log.Error().
			Str("movieID", id).
			Str("actualType", string(*item.Type.Ptr())).
			Msg("Item with specified ID is not a movie")
		return interfaces.Movie{}, fmt.Errorf("item with ID %s is not a movie", id)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Str("movieID", id).
		Str("movieName", *item.Name.Get()).
		Msg("Successfully retrieved movie from Jellyfin")

	movie, err := j.convertToMovie(ctx, &item)
	if err != nil {
		log.Error().
			Err(err).
			Str("movieID", id).
			Str("movieName", *item.Name.Get()).
			Msg("Error converting Jellyfin item to movie format")
		return interfaces.Movie{}, fmt.Errorf("error converting movie data: %w", err)
	}

	log.Debug().
		Str("movieID", id).
		Str("movieName", movie.Metadata.Title).
		Int("year", movie.Metadata.ReleaseYear).
		Msg("Successfully returned movie data")

	return movie, nil
}

// GetCollections retrieves collections from the Jellyfin server
func (j *JellyfinClient) GetCollections(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.Collection, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.ClientID).
		Str("clientType", string(j.ClientType)).
		Str("baseURL", j.config.BaseURL).
		Msg("Retrieving collections from Jellyfin server")

	// Set up query parameters
	includeItemTypes := []jellyfin.BaseItemKind{jellyfin.BASEITEMKIND_BOX_SET}

	limit, startIndex, sortBy, sortOrder := j.getQueryParameters(options)

	// Call the Jellyfin API
	log.Debug().Msg("Making API request to Jellyfin server for collections")
	itemsReq := j.client.ItemsAPI.GetItems(ctx).
		IncludeItemTypes(includeItemTypes).
		Recursive(true).
		Limit(*limit).
		StartIndex(*startIndex).
		SortBy(sortBy).
		SortOrder(sortOrder)

	result, resp, err := itemsReq.Execute()

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.BaseURL).
			Str("apiEndpoint", "/Items").
			Int("statusCode", 0).
			Msg("Failed to fetch collections from Jellyfin")
		return nil, fmt.Errorf("failed to fetch collections: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(result.Items)).
		Int("totalRecordCount", int(*result.TotalRecordCount)).
		Msg("Successfully retrieved collections from Jellyfin")

	// Convert results to expected format
	collections := make([]interfaces.Collection, 0)
	for _, item := range result.Items {
		if *item.Type == "BoxSet" {
			collection, err := j.convertToCollection(ctx, &item)
			if err != nil {
				// Log error but continue
				log.Warn().
					Err(err).
					Str("collectionID", *item.Id).
					Msg("Error converting Jellyfin item to collection format")
				continue
			}
			collections = append(collections, collection)
		}
	}

	log.Info().
		Int("collectionsReturned", len(collections)).
		Msg("Completed GetCollections request")

	return collections, nil
}

// GetTVShows retrieves TV shows from the Jellyfin server

func (j *JellyfinClient) GetTVShows(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.TVShow, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.ClientID).
		Str("clientType", string(j.ClientType)).
		Str("baseURL", j.config.BaseURL).
		Msg("Retrieving TV shows from Jellyfin server")

	// Set up query parameters
	includeItemTypes := []jellyfin.BaseItemKind{jellyfin.BASEITEMKIND_SERIES}

	limit, startIndex, sortBy, sortOrder := j.getQueryParameters(options)

	// Call the Jellyfin API
	log.Debug().Msg("Making API request to Jellyfin server for TV shows")
	itemsReq := j.client.ItemsAPI.GetItems(ctx).
		IncludeItemTypes(includeItemTypes).
		Recursive(true).
		Limit(*limit).
		StartIndex(*startIndex).
		SortBy(sortBy).
		SortOrder(sortOrder)

	result, resp, err := itemsReq.Execute()

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.BaseURL).
			Str("apiEndpoint", "/Items").
			Int("statusCode", 0).
			Msg("Failed to fetch TV shows from Jellyfin")
		return nil, fmt.Errorf("failed to fetch TV shows: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(result.Items)).
		Int("totalRecordCount", int(*result.TotalRecordCount)).
		Msg("Successfully retrieved TV shows from Jellyfin")

	// Convert results to expected format
	shows := make([]interfaces.TVShow, 0)
	for _, item := range result.Items {
		if *item.Type == "Series" {
			show, err := j.convertToTVShow(ctx, &item)
			if err != nil {
				// Log error but continue
				log.Warn().
					Err(err).
					Str("showID", *item.Id).
					Str("showName", *item.Name.Get()).
					Msg("Error converting Jellyfin item to TV show format")
				continue
			}
			shows = append(shows, show)
		}
	}

	log.Info().
		Int("showsReturned", len(shows)).
		Msg("Completed GetTVShows request")

	return shows, nil
}

// Helper function to get int value from pointer with default 0 if nil
func getInt32Value(ptr *int32) int {
	if ptr == nil {
		return 0
	}
	return int(*ptr)
}

// Helper function to get duration from ticks pointer
func getDurationFromTicks(ticks *int64) time.Duration {
	if ticks == nil {
		return 0
	}
	return time.Duration(*ticks/10000000) * time.Second
}

// GetTVShowByID retrieves a specific TV show by ID
func (j *JellyfinClient) GetTVShowByID(ctx context.Context, id string) (interfaces.TVShow, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.ClientID).
		Str("clientType", string(j.ClientType)).
		Str("showID", id).
		Str("baseURL", j.config.BaseURL).
		Msg("Retrieving specific TV show from Jellyfin server")

	// Set up query parameters
	ids := id

	// Call the Jellyfin API
	log.Debug().
		Str("showID", id).
		Msg("Making API request to Jellyfin server")

	itemsReq := j.client.ItemsAPI.GetItems(ctx).Ids(stringToSlice(ids))

	result, resp, err := itemsReq.Execute()

	log.Debug().
		Interface("responseItems", result.Items).
		Msg("Full response data from Jellyfin API")

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.BaseURL).
			Str("apiEndpoint", "/Items").
			Str("showID", id).
			Int("statusCode", 0).
			Msg("Failed to fetch TV show from Jellyfin")
		return interfaces.TVShow{}, fmt.Errorf("failed to fetch TV show: %w", err)
	}

	// Check if any items were returned
	if len(result.Items) == 0 {
		log.Error().
			Str("showID", id).
			Int("statusCode", resp.StatusCode).
			Msg("No TV show found with the specified ID")
		return interfaces.TVShow{}, fmt.Errorf("TV show with ID %s not found", id)
	}

	item := result.Items[0]

	// Double-check that the returned item is a TV show
	if *item.Type != "Series" {
		log.Error().
			Str("showID", id).
			Str("actualType", string(*item.Type.Ptr())).
			Msg("Item with specified ID is not a TV show")
		return interfaces.TVShow{}, fmt.Errorf("item with ID %s is not a TV show", id)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Str("showID", id).
		Str("showName", *item.Name.Get()).
		Msg("Successfully retrieved TV show from Jellyfin")

	show, err := j.convertToTVShow(ctx, &item)
	if err != nil {
		log.Error().
			Err(err).
			Str("showID", id).
			Str("showName", *item.Name.Get()).
			Msg("Error converting Jellyfin item to TV show format")
		return interfaces.TVShow{}, fmt.Errorf("error converting TV show data: %w", err)
	}

	log.Debug().
		Str("showID", id).
		Str("showName", show.Metadata.Title).
		Int("seasonCount", show.SeasonCount).
		Msg("Successfully returned TV show data")

	return show, nil
}

// Helper function to convert Jellyfin item to internal Season type
func (j *JellyfinClient) convertToSeason(ctx context.Context, item *jellyfin.BaseItemDto, showID string) (interfaces.Season, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	// Validate required fields
	if item == nil {
		return interfaces.Season{}, fmt.Errorf("cannot convert nil item to season")
	}

	if item.Id == nil || *item.Id == "" {
		return interfaces.Season{}, fmt.Errorf("season is missing required ID field")
	}

	// Safely get name or fallback to empty string
	title := ""
	if item.Name.IsSet() {
		title = *item.Name.Get()
	}

	log.Debug().
		Str("seasonID", *item.Id).
		Str("seasonName", title).
		Str("showID", showID).
		Msg("Converting Jellyfin item to season format")

	// Safely handle optional fields
	description := ""
	if item.Overview.IsSet() {
		description = *item.Overview.Get()
	}

	// Safely handle season number
	seasonNumber := 0
	if item.IndexNumber.IsSet() {
		seasonNumber = int(*item.IndexNumber.Get())
	}

	// Safely handle episode count
	episodeCount := 0
	if item.ChildCount.IsSet() {
		episodeCount = int(*item.ChildCount.Get())
	}

	// Create the basic season object
	season := interfaces.Season{
		MediaItem: interfaces.MediaItem{
			Metadata: interfaces.MediaMetadata{
				Title:       title,
				Description: description,
				Artwork:     j.getArtworkURLs(item),
			},
			Type:       "season",
			ClientID:   j.ClientID,
			ClientType: string(j.ClientType),
			ExternalID: *item.Id,
		},
		ParentID:     showID,
		Number:       seasonNumber,
		EpisodeCount: episodeCount,
	}

	// Safely set release date if available
	if item.PremiereDate.IsSet() {
		season.ReleaseDate = *item.PremiereDate.Get()
	}

	// Add community rating if available
	if item.CommunityRating.IsSet() {
		season.Metadata.Ratings = append(season.Metadata.Ratings, interfaces.Rating{
			Source: "jellyfin",
			Value:  float32(*item.CommunityRating.Get()),
		})
	}

	// Add user rating if available
	if item.UserData.IsSet() && item.UserData.Get().Rating.IsSet() {
		season.Metadata.UserRating = float32(*item.UserData.Get().Rating.Get())
	}

	// Extract provider IDs if available
	extractProviderIDs(&item.ProviderIds, &season.Metadata.ExternalIDs)

	log.Debug().
		Str("seasonID", *item.Id).
		Str("seasonName", season.Metadata.Title).
		Int("seasonNumber", season.Number).
		Int("episodeCount", season.EpisodeCount).
		Msg("Successfully converted Jellyfin item to season")

	return season, nil
}

// GetTVShowSeasons retrieves seasons for a TV show
func (j *JellyfinClient) GetTVShowSeasons(ctx context.Context, showID string) ([]interfaces.Season, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.ClientID).
		Str("clientType", string(j.ClientType)).
		Str("showID", showID).
		Str("baseURL", j.config.BaseURL).
		Msg("Retrieving seasons for TV show from Jellyfin server")

	// Call the Jellyfin API
	log.Debug().
		Str("showID", showID).
		Msg("Making API request to Jellyfin server for TV show seasons")

	seasonsReq := j.client.TvShowsAPI.GetSeasons(ctx, showID).
		EnableImages(true).
		EnableUserData(true).
		UserId(j.config.UserID)
	result, resp, err := seasonsReq.Execute()

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.BaseURL).
			Str("apiEndpoint", "/Shows/"+showID+"/Seasons").
			Str("showID", showID).
			Int("statusCode", 0).
			Msg("Failed to fetch seasons for TV show from Jellyfin")
		return nil, fmt.Errorf("failed to fetch seasons: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("seasonCount", len(result.Items)).
		Str("showID", showID).
		Msg("Successfully retrieved seasons for TV show from Jellyfin")

	seasons := make([]interfaces.Season, 0)
	for _, item := range result.Items {
		if *item.Type == "Season" {
			season, err := j.convertToSeason(ctx, &item, showID)
			if err != nil {
				// Log error but continue
				log.Warn().
					Err(err).
					Str("seasonID", *item.Id).
					Str("showID", showID).
					Msg("Error converting Jellyfin item to season format")
				continue
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
func (j *JellyfinClient) GetTVShowEpisodes(ctx context.Context, showID string, seasonNumber int) ([]interfaces.Episode, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.ClientID).
		Str("clientType", string(j.ClientType)).
		Str("showID", showID).
		Int("seasonNumber", seasonNumber).
		Str("baseURL", j.config.BaseURL).
		Msg("Retrieving episodes for TV show season from Jellyfin server")

	seasonNum := int32(seasonNumber)

	// Call the Jellyfin API
	log.Debug().
		Str("showID", showID).
		Int("seasonNumber", seasonNumber).
		Msg("Making API request to Jellyfin server for TV show episodes")

	episodesReq := j.client.TvShowsAPI.GetEpisodes(ctx, showID).Season(seasonNum).UserId(j.config.UserID)
	result, resp, err := episodesReq.Execute()

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.BaseURL).
			Str("apiEndpoint", "/Shows/"+showID+"/Episodes").
			Str("showID", showID).
			Int("seasonNumber", seasonNumber).
			Int("statusCode", 0).
			Msg("Failed to fetch episodes for TV show season from Jellyfin")
		return nil, fmt.Errorf("failed to fetch episodes: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("episodeCount", len(result.Items)).
		Str("showID", showID).
		Int("seasonNumber", seasonNumber).
		Msg("Successfully retrieved episodes for TV show season from Jellyfin")

	episodes := make([]interfaces.Episode, 0)
	for _, item := range result.Items {
		if *item.Type == "Episode" {
			episode, err := j.convertToEpisode(ctx, &item)
			if err != nil {
				// Log error but continue
				log.Warn().
					Err(err).
					Str("episodeID", *item.Id).
					Str("showID", showID).
					Int("seasonNumber", seasonNumber).
					Msg("Error converting Jellyfin item to episode format")
				continue
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
func (j *JellyfinClient) GetEpisodeByID(ctx context.Context, id string) (interfaces.Episode, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.ClientID).
		Str("clientType", string(j.ClientType)).
		Str("episodeID", id).
		Str("baseURL", j.config.BaseURL).
		Msg("Retrieving specific episode from Jellyfin server")

	// Set up query parameters
	ids := id

	// Call the Jellyfin API
	log.Debug().
		Str("episodeID", id).
		Msg("Making API request to Jellyfin server")

	itemsReq := j.client.ItemsAPI.GetItems(ctx).Ids(stringToSlice(ids))

	result, resp, err := itemsReq.Execute()

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.BaseURL).
			Str("apiEndpoint", "/Items").
			Str("episodeID", id).
			Int("statusCode", 0).
			Msg("Failed to fetch episode from Jellyfin")
		return interfaces.Episode{}, fmt.Errorf("failed to fetch episode: %w", err)
	}

	// Check if any items were returned
	if len(result.Items) == 0 {
		log.Error().
			Str("episodeID", id).
			Int("statusCode", resp.StatusCode).
			Msg("No episode found with the specified ID")
		return interfaces.Episode{}, fmt.Errorf("episode with ID %s not found", id)
	}

	item := result.Items[0]

	// Double-check that the returned item is an episode
	if *item.Type != jellyfin.BASEITEMKIND_EPISODE {
		log.Error().
			Str("episodeID", id).
			Str("actualType", baseItemKindToString(*item.Type)).
			Msg("Item with specified ID is not an episode")
		return interfaces.Episode{}, fmt.Errorf("item with ID %s is not an episode", id)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Str("episodeID", id).
		Str("episodeName", *item.Name.Get()).
		Msg("Successfully retrieved episode from Jellyfin")

	episode, err := j.convertToEpisode(ctx, &item)
	if err != nil {
		log.Error().
			Err(err).
			Str("episodeID", id).
			Msg("Error converting Jellyfin item to episode format")
		return interfaces.Episode{}, fmt.Errorf("error converting episode data: %w", err)
	}

	log.Debug().
		Str("episodeID", id).
		Str("episodeName", episode.Metadata.Title).
		Int64("episodeNumber", episode.Number).
		Int("seasonNumber", episode.SeasonNumber).
		Msg("Successfully returned episode data")

	return episode, nil
}

// GetMusic retrieves music tracks from the Jellyfin server
func (j *JellyfinClient) GetMusic(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.MusicTrack, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.ClientID).
		Str("clientType", string(j.ClientType)).
		Str("baseURL", j.config.BaseURL).
		Msg("Retrieving music tracks from Jellyfin server")

	// Set up query parameters
	includeItemTypes := []jellyfin.BaseItemKind{jellyfin.BASEITEMKIND_AUDIO}

	limit, startIndex, sortBy, sortOrder := j.getQueryParameters(options)

	// Call the Jellyfin API
	log.Debug().Msg("Making API request to Jellyfin server for music tracks")
	itemsReq := j.client.ItemsAPI.GetItems(ctx).
		IncludeItemTypes(includeItemTypes).
		Recursive(true).
		Limit(*limit).
		StartIndex(*startIndex).
		SortBy(sortBy).
		SortOrder(sortOrder).
		UserId(j.config.UserID)

	result, resp, err := itemsReq.Execute()

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.BaseURL).
			Str("apiEndpoint", "/Items").
			Int("statusCode", 0).
			Msg("Failed to fetch music tracks from Jellyfin")
		return nil, fmt.Errorf("failed to fetch music tracks: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(result.Items)).
		Int("totalRecordCount", int(*result.TotalRecordCount)).
		Msg("Successfully retrieved music tracks from Jellyfin")

	// Convert results to expected format
	tracks := make([]interfaces.MusicTrack, 0)
	for _, item := range result.Items {
		if *item.Type == "Audio" {
			track := interfaces.MusicTrack{
				MediaItem: interfaces.MediaItem{
					Metadata: interfaces.MediaMetadata{
						Title:       *item.Name.Get(),
						Description: *item.Overview.Get(),
						Duration:    getDurationFromTicks(item.RunTimeTicks.Get()),
						Artwork:     j.getArtworkURLs(&item),
					},
					ExternalID: *item.Id,
					Type:       "track",
					ClientID:   j.ClientID,
					ClientType: string(j.ClientType),
				},
				Number: int(*item.IndexNumber.Get()),
			}

			// Set album info if available
			if item.AlbumId.IsSet() {
				track.AlbumID = *item.AlbumId.Get()
			}
			if item.Album.IsSet() {
				track.AlbumName = *item.Album.Get()
			}

			// Add artist information if available
			if item.ArtistItems != nil && len(item.ArtistItems) > 0 {
				track.ArtistID = *item.ArtistItems[0].Id
				track.ArtistName = *item.ArtistItems[0].Name.Get()
			}

			extractProviderIDs(&item.ProviderIds, &track.Metadata.ExternalIDs)

			tracks = append(tracks, track)
		}
	}

	log.Info().
		Int("tracksReturned", len(tracks)).
		Msg("Completed GetMusic request")

	return tracks, nil
}

// GetMusicArtists retrieves music artists from the Jellyfin server
func (j *JellyfinClient) GetMusicArtists(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.MusicArtist, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.ClientID).
		Str("clientType", string(j.ClientType)).
		Str("baseURL", j.config.BaseURL).
		Msg("Retrieving music artists from Jellyfin server")

	limit, startIndex, sortBy, sortOrder := j.getQueryParameters(options)

	// Call the Jellyfin API
	log.Debug().Msg("Making API request to Jellyfin server for music artists")
	artistReq := j.client.ArtistsAPI.GetArtists(ctx).
		Limit(*limit).
		StartIndex(*startIndex).
		SortBy(sortBy).
		SortOrder(sortOrder).
		UserId(j.config.UserID)

	result, resp, err := artistReq.Execute()

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.BaseURL).
			Str("apiEndpoint", "/Artists").
			Int("statusCode", 0).
			Msg("Failed to fetch music artists from Jellyfin")
		return nil, fmt.Errorf("failed to fetch music artists: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(result.Items)).
		Int("totalRecordCount", int(*result.TotalRecordCount)).
		Msg("Successfully retrieved music artists from Jellyfin")

	// Convert results to expected format
	artists := make([]interfaces.MusicArtist, 0)
	for _, item := range result.Items {
		artist := interfaces.MusicArtist{
			MediaItem: interfaces.MediaItem{
				Metadata: interfaces.MediaMetadata{
					Title:       *item.Name.Get(),
					Description: *item.Overview.Get(),
					Artwork:     j.getArtworkURLs(&item),
					Genres:      item.Genres,
				},
				ExternalID: *item.Id,
				Type:       "artist",
				ClientID:   j.ClientID,
				ClientType: string(j.ClientType),
			},
		}

		extractProviderIDs(&item.ProviderIds, &artist.Metadata.ExternalIDs)

		artists = append(artists, artist)
	}

	log.Info().
		Int("artistsReturned", len(artists)).
		Msg("Completed GetMusicArtists request")

	return artists, nil
}

// GetMusicAlbums retrieves music albums from the Jellyfin server
func (j *JellyfinClient) GetMusicAlbums(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.MusicAlbum, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.ClientID).
		Str("clientType", string(j.ClientType)).
		Str("baseURL", j.config.BaseURL).
		Msg("Retrieving music albums from Jellyfin server")

	// Set up query parameters
	includeItemTypes := []jellyfin.BaseItemKind{jellyfin.BASEITEMKIND_MUSIC_ALBUM}

	limit, startIndex, sortBy, sortOrder := j.getQueryParameters(options)

	// Call the Jellyfin API
	log.Debug().Msg("Making API request to Jellyfin server for music albums")
	itemsReq := j.client.ItemsAPI.GetItems(ctx).
		IncludeItemTypes(includeItemTypes).
		Recursive(true).
		Limit(*limit).
		StartIndex(*startIndex).
		SortBy(sortBy).
		SortOrder(sortOrder).
		UserId(j.config.UserID)

	result, resp, err := itemsReq.Execute()

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.BaseURL).
			Str("apiEndpoint", "/Items").
			Int("statusCode", 0).
			Msg("Failed to fetch music albums from Jellyfin")
		return nil, fmt.Errorf("failed to fetch music albums: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(result.Items)).
		Int("totalRecordCount", int(*result.TotalRecordCount)).
		Msg("Successfully retrieved music albums from Jellyfin")

	// Convert results to expected format
	albums := make([]interfaces.MusicAlbum, 0)
	for _, item := range result.Items {
		album := interfaces.MusicAlbum{
			MediaItem: interfaces.MediaItem{
				Metadata: interfaces.MediaMetadata{
					Title:       *item.Name.Get(),
					Description: *item.Overview.Get(),
					ReleaseYear: int(*item.ProductionYear.Get()),
					Genres:      item.Genres,
					Artwork:     j.getArtworkURLs(&item),
				},
				Type:       "album",
				ExternalID: *item.Id,
				ClientID:   j.ClientID,
				ClientType: string(j.ClientType),
			},
			TrackCount: int(*item.ChildCount.Get()),
		}

		// Set album artist if available
		if item.AlbumArtist.IsSet() {
			album.ArtistName = *item.AlbumArtist.Get()
		}

		extractProviderIDs(&item.ProviderIds, &album.Metadata.ExternalIDs)

		albums = append(albums, album)
	}

	log.Info().
		Int("albumsReturned", len(albums)).
		Msg("Completed GetMusicAlbums request")

	return albums, nil
}

// GetMusicTrackByID retrieves a specific music track by ID
func (j *JellyfinClient) GetMusicTrackByID(ctx context.Context, id string) (interfaces.MusicTrack, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.ClientID).
		Str("clientType", string(j.ClientType)).
		Str("trackID", id).
		Str("baseURL", j.config.BaseURL).
		Msg("Retrieving specific music track from Jellyfin server")

	// Set up query parameters
	ids := id

	// Call the Jellyfin API
	log.Debug().
		Str("trackID", id).
		Msg("Making API request to Jellyfin server")

	itemsReq := j.client.ItemsAPI.GetItems(ctx)

	itemsReq.Ids(stringToSlice(ids))

	result, resp, err := itemsReq.Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.BaseURL).
			Str("apiEndpoint", "/Items").
			Str("trackID", id).
			Int("statusCode", 0).
			Msg("Failed to fetch music track from Jellyfin")
		return interfaces.MusicTrack{}, fmt.Errorf("failed to fetch music track: %w", err)
	}

	// Check if any items were returned
	if len(result.Items) == 0 {
		log.Error().
			Str("trackID", id).
			Int("statusCode", resp.StatusCode).
			Msg("No music track found with the specified ID")
		return interfaces.MusicTrack{}, fmt.Errorf("music track with ID %s not found", id)
	}

	item := result.Items[0]

	// Double-check that the returned item is a music track
	if *item.Type != "Audio" {
		log.Error().
			Str("trackID", id).
			Str("actualType", baseItemKindToString(*item.Type)).
			Msg("Item with specified ID is not a music track")
		return interfaces.MusicTrack{}, fmt.Errorf("item with ID %s is not a music track", id)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Str("trackID", id).
		Str("trackName", *item.Name.Get()).
		Msg("Successfully retrieved music track from Jellyfin")

	track := interfaces.MusicTrack{
		MediaItem: interfaces.MediaItem{
			Metadata: interfaces.MediaMetadata{
				Title:       *item.Name.Get(),
				Description: *item.Overview.Get(),
				Duration:    getDurationFromTicks(item.RunTimeTicks.Get()),
				Artwork:     j.getArtworkURLs(&item),
			},
			ExternalID: *item.Id,
			Type:       "track",
			ClientID:   j.ClientID,
			ClientType: string(j.ClientType),
		},
		Number: int(*item.IndexNumber.Get()),
	}

	// Set album info if available
	if item.AlbumId.IsSet() {
		track.AlbumID = *item.AlbumId.Get()
	}
	if item.Album.IsSet() {
		track.AlbumName = *item.Album.Get()
	}

	// Add artist information if available
	if item.ArtistItems != nil && len(item.ArtistItems) > 0 {
		track.ArtistID = *item.ArtistItems[0].Id
		track.ArtistName = *item.ArtistItems[0].Name.Get()
	}

	// Extract provider IDs
	extractProviderIDs(&item.ProviderIds, &track.Metadata.ExternalIDs)

	log.Debug().
		Str("trackID", id).
		Str("trackName", track.Metadata.Title).
		Int("trackNumber", track.Number).
		Msg("Successfully returned music track data")

	return track, nil
}

// GetPlaylists retrieves playlists from the Jellyfin server
func (j *JellyfinClient) GetPlaylists(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.Playlist, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.ClientID).
		Str("clientType", string(j.ClientType)).
		Str("baseURL", j.config.BaseURL).
		Msg("Retrieving playlists from Jellyfin server")

	// Set up query parameters
	includeItemTypes := []jellyfin.BaseItemKind{jellyfin.BASEITEMKIND_PLAYLIST}
	recursive := true

	limit, startIndex, sortBy, sortOrder := j.getQueryParameters(options)

	// Call the Jellyfin API
	log.Debug().Msg("Making API request to Jellyfin server for playlists")

	itemsRequest := j.client.ItemsAPI.GetItems(ctx).
		IncludeItemTypes(includeItemTypes).
		Recursive(recursive).
		Limit(*limit).
		StartIndex(*startIndex).
		SortBy(sortBy).
		SortOrder(sortOrder)

	result, resp, err := itemsRequest.Execute()

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.BaseURL).
			Str("apiEndpoint", "/Items").
			Int("statusCode", 0).
			Msg("Failed to fetch playlists from Jellyfin")
		return nil, fmt.Errorf("failed to fetch playlists: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(result.Items)).
		Int("totalRecordCount", int(*result.TotalRecordCount)).
		Msg("Successfully retrieved playlists from Jellyfin")

	// Convert results to expected format
	playlists := make([]interfaces.Playlist, 0)
	for _, item := range result.Items {
		if *item.Type == "Playlist" {
			playlist := interfaces.Playlist{
				MediaItem: interfaces.MediaItem{
					Metadata: interfaces.MediaMetadata{
						Title:       *item.Name.Get(),
						Description: *item.Overview.Get(),
						Artwork:     j.getArtworkURLs(&item),
					},
					ExternalID: *item.Id,
					Type:       "playlist",
					ClientID:   j.ClientID,
					ClientType: string(j.ClientType),
				},
				ItemCount: int(*item.ChildCount.Get()),
				IsPublic:  true, // Assume public by default in Jellyfin
			}
			playlists = append(playlists, playlist)
		}
	}

	log.Info().
		Int("playlistsReturned", len(playlists)).
		Msg("Completed GetPlaylists request")

	return playlists, nil
}

// GetWatchHistory retrieves watch history for the current user
func (j *JellyfinClient) GetWatchHistory(ctx context.Context, options *interfaces.QueryOptions) ([]interfaces.WatchHistoryItem, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.ClientID).
		Str("clientType", string(j.ClientType)).
		Str("baseURL", j.config.BaseURL).
		Msg("Retrieving watch history from Jellyfin server")

	limit, startIndex, sortBy, sortOrder := j.getQueryParameters(options)

	// Call the Jellyfin API to get resumed items
	log.Debug().Msg("Making API request to Jellyfin server for resume items")
	watchedItemsReq := j.client.ItemsAPI.GetItems(ctx).
		Limit(*limit).
		StartIndex(*startIndex).
		SortBy(sortBy).
		SortOrder(sortOrder).
		UserId(j.config.UserID).
		IsPlayed(true)

	result, resp, err := watchedItemsReq.Execute()

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.BaseURL).
			Str("apiEndpoint", "/UserItems/Resume").
			Int("statusCode", 0).
			Msg("Failed to fetch watch history from Jellyfin")
		return nil, fmt.Errorf("failed to fetch watch history: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(result.Items)).
		Int("totalRecordCount", int(*result.TotalRecordCount)).
		Msg("Successfully retrieved watch history from Jellyfin")

	// Convert results to expected format
	historyItems := make([]interfaces.WatchHistoryItem, 0)
	for _, item := range result.Items {

		userDataReq := j.client.ItemsAPI.GetItemUserData(ctx, *item.Id)
		userData, resp, err := userDataReq.Execute()

		if err != nil {
			continue
		}

		log.Info().
			Int("statusCode", resp.StatusCode).
			Int32("playCount", userData.GetPlayCount()).
			Msg("Successfully retrieved user item data from Jellyfin")

		historyItem := interfaces.WatchHistoryItem{
			MediaItem: interfaces.MediaItem{
				Metadata: interfaces.MediaMetadata{
					Title:       *item.Name.Get(),
					Description: *item.Overview.Get(),
					Artwork:     j.getArtworkURLs(&item),
				},
				ExternalID: *item.Id,
				ClientID:   j.ClientID,
				ClientType: string(j.ClientType),
			},
			PlayedPercentage: *userData.PlayedPercentage.Get(),
			LastWatchedAt:    *userData.LastPlayedDate.Get(), // Default to now if not available
		}

		// Set type based on item type
		switch *item.Type {
		case "Movie":
			historyItem.MediaItem.Type = "movie"
		case "Series":
			historyItem.MediaItem.Type = "tvshow"
		case "Episode":
			historyItem.MediaItem.Type = "episode"

			// Add additional episode info if available
			if item.SeriesName.IsSet() {
				historyItem.SeriesName = *item.SeriesName.Get()
			}
			if item.ParentIndexNumber.IsSet() {
				historyItem.SeasonNumber = int(*item.ParentIndexNumber.Get())
			}
			if item.IndexNumber.IsSet() {
				historyItem.EpisodeNumber = int(*item.IndexNumber.Get())
			}
		}

		// Set last played date if available
		if userData.LastPlayedDate.IsSet() {
			historyItem.LastWatchedAt = *userData.LastPlayedDate.Get()
		}

		historyItems = append(historyItems, historyItem)
	}

	log.Info().
		Int("historyItemsReturned", len(historyItems)).
		Msg("Completed GetWatchHistory request")

	return historyItems, nil
}

// GetMusicGenres retrieves music genres from the Jellyfin server
func (j *JellyfinClient) GetMusicGenres(ctx context.Context) ([]string, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.ClientID).
		Str("clientType", string(j.ClientType)).
		Str("baseURL", j.config.BaseURL).
		Msg("Retrieving music genres from Jellyfin server")

	// Call the Jellyfin API
	log.Debug().Msg("Making API request to Jellyfin server for music genres")
	musicGenresReq := j.client.MusicGenresAPI.GetMusicGenres(ctx)
	result, resp, err := musicGenresReq.Execute()

	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.BaseURL).
			Str("apiEndpoint", "/MusicGenres").
			Int("statusCode", 0).
			Msg("Failed to fetch music genres from Jellyfin")
		return nil, fmt.Errorf("failed to fetch music genres: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(result.Items)).
		Int("totalRecordCount", int(*result.TotalRecordCount)).
		Msg("Successfully retrieved music genres from Jellyfin")

	// Convert results to expected format
	genres := make([]string, 0, len(result.Items))
	for _, item := range result.Items {
		if item.Name.Get() != nil {
			genres = append(genres, *item.Name.Get())
		}
	}

	log.Info().
		Int("genresReturned", len(genres)).
		Msg("Completed GetMusicGenres request")

	return genres, nil
}

// GetMovieGenres retrieves movie genres from the Jellyfin server
func (j *JellyfinClient) GetMovieGenres(ctx context.Context) ([]string, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", j.ClientID).
		Str("clientType", string(j.ClientType)).
		Str("baseURL", j.config.BaseURL).
		Msg("Retrieving movie genres from Jellyfin server")

	// Set up query parameters to get only movie genres
	includeItemTypes := []jellyfin.BaseItemKind{jellyfin.BASEITEMKIND_MOVIE}
	// Call the Jellyfin API
	log.Debug().Msg("Making API request to Jellyfin server for movie genres")
	genresReq := j.client.GenresAPI.GetGenres(ctx)

	genresReq.IncludeItemTypes(includeItemTypes)
	result, resp, err := genresReq.Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", j.config.BaseURL).
			Str("apiEndpoint", "/Genres").
			Int("statusCode", 0).
			Msg("Failed to fetch movie genres from Jellyfin")
		return nil, fmt.Errorf("failed to fetch movie genres: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int("totalItems", len(result.Items)).
		Int("totalRecordCount", int(*result.TotalRecordCount)).
		Msg("Successfully retrieved movie genres from Jellyfin")

	// Convert results to expected format
	genres := make([]string, 0, len(result.Items))
	for _, item := range result.Items {
		if item.Name.Get() != nil {
			genres = append(genres, *item.Name.Get())
		}
	}

	log.Info().
		Int("genresReturned", len(genres)).
		Msg("Completed GetMovieGenres request")

	return genres, nil
}
