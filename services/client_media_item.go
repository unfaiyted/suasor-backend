package services

import (
	"context"
	"fmt"
	"time"

	"suasor/clients"
	mediaclient "suasor/clients/media"
	"suasor/clients/media/providers"
	"suasor/clients/media/types"
	clienttypes "suasor/clients/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/utils/logger"
)

// ClientMediaItemService defines the interface for client-associated media item operations
// This service extends CoreMediaItemService with operations specific to media items
// that are linked to external clients like Plex, Emby, etc.
type ClientMediaItemService[T clienttypes.ClientMediaConfig, U types.MediaData] interface {
	// Include all core service methods
	UserMediaItemService[U]

	// Client-specific operations
	GetByClientID(ctx context.Context, clientID uint64, mediaType types.MediaType, limit int, offset int) ([]*models.MediaItem[U], error)

	SearchAcrossClients(ctx context.Context, query types.QueryOptions, clientIDs []uint64) (map[uint64][]*models.MediaItem[U], error)
	SearchClient(ctx context.Context, clientID uint64, options types.QueryOptions) ([]*models.MediaItem[U], error)

	DeleteClientItem(ctx context.Context, clientID uint64, itemID string) error

	// Sync operations
	// Multi-client operations
	GetByMultipleClients(ctx context.Context, clientIDs []uint64) ([]*models.MediaItem[U], error)
	SyncItemBetweenClients(ctx context.Context, itemID uint64, sourceClientID uint64, targetClientID uint64, targetItemID string) error
}

// clientMediaItemService implements ClientMediaItemService
type clientMediaItemService[T clienttypes.ClientMediaConfig, U types.MediaData] struct {
	UserMediaItemService[U] // Embed the user service
	clientRepo              repository.ClientRepository[T]
	itemRepo                repository.ClientMediaItemRepository[U]
	clientFactory           *clients.ClientProviderFactoryService
}

// NewClientMediaItemService creates a new client-associated media item service
func NewClientMediaItemService[T clienttypes.ClientMediaConfig, U types.MediaData](
	userService UserMediaItemService[U],
	clientRepo repository.ClientRepository[T],
	itemRepo repository.ClientMediaItemRepository[U],
	clientFactory *clients.ClientProviderFactoryService,
) ClientMediaItemService[T, U] {
	return &clientMediaItemService[T, U]{
		UserMediaItemService: userService,
		clientRepo:           clientRepo,
		itemRepo:             itemRepo,
		clientFactory:        clientFactory,
	}
}

// getClientProvider gets the appropriate provider for the media type
func (s *clientMediaItemService[T, U]) getClientProvider(ctx context.Context, clientID uint64) (interface{}, error) {
	log := logger.LoggerFromContext(ctx)

	// Get the client configuration from repository
	client, err := s.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		log.Error().Err(err).Uint64("clientID", clientID).Msg("Failed to get client configuration")
		return nil, fmt.Errorf("failed to get client configuration: %w", err)
	}

	// Get a zero value of U to determine the media type
	var zero U
	mediaType := types.GetMediaTypeFromTypeName(zero)

	log.Debug().
		Uint64("clientID", clientID).
		Str("mediaType", string(mediaType)).
		Msg("Getting client provider for media type")

	// Return the appropriate provider based on media type
	switch mediaType {
	case types.MediaTypeMovie:
		provider, err := s.clientFactory.GetMovieProvider(ctx, clientID, client.Config)
		if err != nil {
			return nil, err
		}
		return provider, nil
	case types.MediaTypeSeries, types.MediaTypeSeason, types.MediaTypeEpisode:
		provider, err := s.clientFactory.GetSeriesProvider(ctx, clientID, client.Config)
		if err != nil {
			return nil, err
		}
		return provider, nil
	case types.MediaTypeArtist, types.MediaTypeAlbum, types.MediaTypeTrack:
		provider, err := s.clientFactory.GetMusicProvider(ctx, clientID, client.Config)
		if err != nil {
			return nil, err
		}
		return provider, nil
	case types.MediaTypePlaylist:
		provider, err := s.clientFactory.GetPlaylistProvider(ctx, clientID, client.Config)
		if err != nil {
			return nil, err
		}
		return provider, nil
	case types.MediaTypeCollection:
		provider, err := s.clientFactory.GetCollectionProvider(ctx, clientID, client.Config)
		if err != nil {
			return nil, err
		}
		return provider, nil
	default:
		return nil, fmt.Errorf("provider, unsupported media type: %s", mediaType)
	}
}

// fetchItemFromClient fetches a media item from the client provider
func (s *clientMediaItemService[T, U]) fetchItemFromClient(ctx context.Context, clientID uint64, itemID string) (*models.MediaItem[U], error) {
	log := logger.LoggerFromContext(ctx)

	provider, err := s.getClientProvider(ctx, clientID)
	if err != nil {
		return nil, err
	}

	var zero U
	mediaType := types.GetMediaTypeFromTypeName(zero)

	log.Debug().
		Uint64("clientID", clientID).
		Str("itemID", itemID).
		Str("mediaType", string(mediaType)).
		Msg("Fetching item from client")

	var mediaItem *models.MediaItem[U]
	// Call the appropriate provider method based on media type
	switch mediaType {
	case types.MediaTypeMovie:
		movieProvider, _ := provider.(providers.MovieProvider)

		movie, err := movieProvider.GetMovieByID(ctx, itemID)
		if err != nil {
			return nil, err
		}

		// Convert from, to generic types
		mediaItem = models.NewMediaItemCopy[*types.Movie, U](movie)

	case types.MediaTypeSeries, types.MediaTypeSeason, types.MediaTypeEpisode:
		seriesProvider, _ := provider.(providers.SeriesProvider)
		options := &types.QueryOptions{
			MediaType: mediaType,
			ItemIDs:   itemID,
		}

		if mediaType == types.MediaTypeSeries {
			series, err := seriesProvider.GetSeries(ctx, options)
			if err != nil || len(series) == 0 {
				return nil, err
			}
			mediaItem = models.NewMediaItemCopy[*types.Series, U](series[0])
		} else if mediaType == types.MediaTypeEpisode {
			// Handle episode fetching
			episode, err := seriesProvider.GetEpisodeByID(ctx, itemID)
			if err != nil {
				return nil, err
			}
			mediaItem = models.NewMediaItemCopy[*types.Episode, U](episode)
		}

	case types.MediaTypeArtist, types.MediaTypeAlbum, types.MediaTypeTrack:
		musicProvider, _ := provider.(providers.MusicProvider)
		options := &types.QueryOptions{
			MediaType: mediaType,
			ItemIDs:   itemID,
		}

		if mediaType == types.MediaTypeArtist {
			artists, err := musicProvider.GetMusicArtists(ctx, options)
			if err != nil || len(artists) == 0 {
				return nil, err
			}
			mediaItem = models.NewMediaItemCopy[*types.Artist, U](artists[0])
		} else if mediaType == types.MediaTypeAlbum {
			albums, err := musicProvider.GetMusicAlbums(ctx, options)
			if err != nil || len(albums) == 0 {
				return nil, err
			}
			mediaItem = models.NewMediaItemCopy[*types.Album, U](albums[0])
		} else if mediaType == types.MediaTypeTrack {
			track, err := musicProvider.GetMusicTrackByID(ctx, itemID)
			if err != nil {
				return nil, err
			}
			mediaItem = models.NewMediaItemCopy[*types.Track, U](track)
		}

	default:
		return nil, fmt.Errorf("unsupported media type: %s", mediaType)
	}

	return mediaItem, nil
}

// saveOrUpdateMediaItem saves a new media item or updates an existing one
func (s *clientMediaItemService[T, U]) saveOrUpdateMediaItem(ctx context.Context, mediaItem *models.MediaItem[U]) (*models.MediaItem[U], error) {
	log := logger.LoggerFromContext(ctx)

	clientID := mediaItem.SyncClients[0].ID
	clientItemID, exists := mediaItem.GetClientItemID(clientID)
	if !exists {
		return nil, fmt.Errorf("client item ID not found")
	}

	// First check if the item already exists in our database
	existingItem, err := s.itemRepo.GetByClientItemID(ctx, clientID, clientItemID)
	if err == nil && existingItem != nil {
		// Item exists, update it
		log.Debug().
			Uint64("clientID", clientID).
			Str("clientItemID", clientItemID).
			Uint64("itemID", existingItem.ID).
			Msg("Updating existing media item")

		// Update with new data
		existingItem.Data = mediaItem.Data
		// Add additional fields to update if needed

		// Save the updated item
		updatedItem, err := s.itemRepo.Update(ctx, existingItem)
		if err != nil {
			log.Error().Err(err).
				Uint64("itemID", existingItem.ID).
				Msg("Failed to update media item")
			return nil, fmt.Errorf("failed to update media item: %w", err)
		}

		return updatedItem, nil
	}

	// Item doesn't exist, create a new one
	log.Debug().
		Uint64("clientID", clientID).
		Str("clientItemID", clientItemID).
		Msg("Creating new media item")

	newItem, err := s.itemRepo.Create(ctx, mediaItem)
	if err != nil {
		log.Error().Err(err).
			Uint64("clientID", clientID).
			Str("clientItemID", clientItemID).
			Msg("Failed to create media item")
		return nil, fmt.Errorf("failed to create media item: %w", err)
	}

	return newItem, nil
}

// GetByClientID retrieves all media items associated with a specific client
func (s *clientMediaItemService[T, U]) GetByClientID(ctx context.Context, clientID uint64, mediaType types.MediaType, limit int, offset int) ([]*models.MediaItem[U], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("clientID", clientID).
		Str("mediaType", string(mediaType)).
		Int("limit", limit).
		Int("offset", offset).
		Msg("Getting media items by client ID")

	// Get items from the client based on the media type
	provider, err := s.getClientProvider(ctx, clientID)
	if err != nil {
		log.Error().Err(err).
			Uint64("clientID", clientID).
			Str("mediaType", string(mediaType)).
			Msg("Failed to get client provider")
		return nil, fmt.Errorf("failed to get client provider: %w", err)
	}

	var clientItems []*models.MediaItem[U]
	options := &types.QueryOptions{
		MediaType: mediaType,
		Limit:     limit,
		Offset:    offset,
	}

	// Use the appropriate provider method based on media type
	switch mediaType {
	case types.MediaTypeMovie:
		movieProvider, ok := provider.(providers.MovieProvider)
		if !ok {
			return nil, fmt.Errorf("provider does not support movies")
		}
		movies, err := movieProvider.GetMovies(ctx, options)
		if err != nil {
			log.Error().Err(err).
				Uint64("clientID", clientID).
				Msg("Failed to get movies from client")
			return nil, fmt.Errorf("failed to get movies from client: %w", err)
		}
		// Convert from specific type to generic type
		for _, movie := range movies {
			item := models.NewMediaItemCopy[*types.Movie, U](movie)
			clientItems = append(clientItems, item)
		}

	case types.MediaTypeSeries:
		seriesProvider, ok := provider.(providers.SeriesProvider)
		if !ok {
			return nil, fmt.Errorf("provider does not support series")
		}
		series, err := seriesProvider.GetSeries(ctx, options)
		if err != nil {
			log.Error().Err(err).
				Uint64("clientID", clientID).
				Msg("Failed to get series from client")
			return nil, fmt.Errorf("failed to get series from client: %w", err)
		}
		// Convert from specific type to generic type
		for _, show := range series {
			item := models.NewMediaItemCopy[*types.Series, U](show)
			clientItems = append(clientItems, item)
		}

	case types.MediaTypeEpisode:
		// seriesProvider, ok := provider.(providers.SeriesProvider)
		// if !ok {
		// 	return nil, fmt.Errorf("provider does not support episodes")
		// }
		// Episodes require a series ID, so we'd need to get series first and then episodes
		// For this implementation, we'll return an error

		return nil, fmt.Errorf("direct episode retrieval not supported, please use series ID to get episodes")

	case types.MediaTypeArtist:
		musicProvider, ok := provider.(providers.MusicProvider)
		if !ok {
			return nil, fmt.Errorf("provider does not support artists")
		}
		artists, err := musicProvider.GetMusicArtists(ctx, options)
		if err != nil {
			log.Error().Err(err).
				Uint64("clientID", clientID).
				Msg("Failed to get artists from client")
			return nil, fmt.Errorf("failed to get artists from client: %w", err)
		}
		// Convert from specific type to generic type
		for _, artist := range artists {
			item := models.NewMediaItemCopy[*types.Artist, U](artist)
			clientItems = append(clientItems, item)
		}

	case types.MediaTypeAlbum:
		musicProvider, ok := provider.(providers.MusicProvider)
		if !ok {
			return nil, fmt.Errorf("provider does not support albums")
		}
		albums, err := musicProvider.GetMusicAlbums(ctx, options)
		if err != nil {
			log.Error().Err(err).
				Uint64("clientID", clientID).
				Msg("Failed to get albums from client")
			return nil, fmt.Errorf("failed to get albums from client: %w", err)
		}
		// Convert from specific type to generic type
		for _, album := range albums {
			item := models.NewMediaItemCopy[*types.Album, U](album)
			clientItems = append(clientItems, item)
		}

	case types.MediaTypeTrack:
		// musicProvider, ok := provider.(providers.MusicProvider)
		// if !ok {
		// 	return nil, fmt.Errorf("provider does not support tracks")
		// }
		// tracks, err := musicProvider.GetMusicTracks(ctx, options)
		// Tracks usually require an album ID, similar to episodes requiring a series ID
		return nil, fmt.Errorf("direct track retrieval not supported, please use album ID to get tracks")
	case types.MediaTypePlaylist:
		playlistProvider, ok := provider.(providers.PlaylistProvider)
		if !ok {
			return nil, fmt.Errorf("provider does not support playlists")
		}
		playlists, err := playlistProvider.SearchPlaylists(ctx, options)
		if err != nil {
			log.Error().Err(err).
				Uint64("clientID", clientID).
				Msg("Failed to get playlists from client")
			return nil, fmt.Errorf("failed to get playlists from client: %w", err)
		}
		// Convert from specific type to generic type
		for _, playlist := range playlists {
			item := models.NewMediaItemCopy[*types.Playlist, U](playlist)
			clientItems = append(clientItems, item)
		}

	case types.MediaTypeCollection:
		collectionProvider, ok := provider.(providers.CollectionProvider)
		if !ok {
			return nil, fmt.Errorf("provider does not support collections")
		}
		collections, err := collectionProvider.SearchCollections(ctx, options)
		if err != nil {
			log.Error().Err(err).
				Uint64("clientID", clientID).
				Msg("Failed to get collections from client")
			return nil, fmt.Errorf("failed to get collections from client: %w", err)
		}
		// Convert from specific type to generic type
		for _, collection := range collections {
			item := models.NewMediaItemCopy[*types.Collection, U](collection)
			clientItems = append(clientItems, item)
		}

	default:
		return nil, fmt.Errorf("unsupported media type: %s", mediaType)
	}

	log.Debug().
		Uint64("clientID", clientID).
		Int("clientItemsCount", len(clientItems)).
		Msg("Retrieved media items from client")

	if len(clientItems) == 0 {
		return []*models.MediaItem[U]{}, nil
	}

	// Create slices to hold client item IDs and title/year combinations for efficient DB lookups
	clientItemIDs := make([]string, 0, len(clientItems))
	titleYearKeys := make(map[string]struct{})

	for _, item := range clientItems {
		clientItemID, exists := item.GetClientItemID(clientID)
		if exists && clientItemID != "" {
			clientItemIDs = append(clientItemIDs, clientItemID)
		}

		// Also prepare title+year+mediaType keys for secondary matching
		if item.Title != "" && item.ReleaseYear > 0 {
			key := fmt.Sprintf("%s_%d_%s", item.Title, item.ReleaseYear, item.Type)
			titleYearKeys[key] = struct{}{}
		}
	}

	// Get client information (type) for adding to sync clients
	client, err := s.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		log.Error().Err(err).
			Uint64("clientID", clientID).
			Msg("Failed to get client information")
		return nil, fmt.Errorf("failed to get client information: %w", err)
	}
	clientType := client.Config.GetType()

	// Fetch only the relevant items from the database by client item IDs
	existingItemsMap := make(map[string]*models.MediaItem[U])
	if len(clientItemIDs) > 0 {
		// TODO: Implement a batch lookup method in the repository
		// For now, we'll use individual lookups
		for _, clientItemID := range clientItemIDs {
			dbItem, err := s.itemRepo.GetByClientItemID(ctx, clientID, clientItemID)
			if err == nil && dbItem != nil {
				existingItemsMap[clientItemID] = dbItem
			}
		}
	}

	// Fetch items that might match by title+year+mediaType
	titleYearTypeMap := make(map[string]*models.MediaItem[U])
	if len(titleYearKeys) > 0 {
		// TODO: Implement a batch lookup method in the repository
		// For now, we'll need to get items that could potentially match
		// This would be much more efficient with a specific DB query
		// Get all items of this media type
		mediaTypeItems, err := s.itemRepo.GetByType(ctx, mediaType)
		if err == nil {
			for _, item := range mediaTypeItems {
				if item.Title != "" && item.ReleaseYear > 0 {
					key := fmt.Sprintf("%s_%d_%s", item.Title, item.ReleaseYear, item.Type)
					if _, exists := titleYearKeys[key]; exists {
						titleYearTypeMap[key] = item
					}
				}
			}
		}
	}

	// Process each client item to match with database items or create new ones
	resultItems := make([]*models.MediaItem[U], 0, len(clientItems))

	for _, clientItem := range clientItems {
		// Extract client item ID
		clientItemID, exists := clientItem.GetClientItemID(clientID)
		if !exists || clientItemID == "" {
			// Skip items that don't have a client item ID
			log.Warn().
				Uint64("clientID", clientID).
				Str("title", clientItem.Title).
				Msg("Client item missing client item ID, skipping")
			continue
		}

		// Check if we already have this item in the database by client ID
		dbItem, found := existingItemsMap[clientItemID]
		if found {
			// Item exists, update it with new data from client
			dbItem.Data = clientItem.Data
			dbItem.Title = clientItem.Title
			dbItem.ReleaseDate = clientItem.ReleaseDate
			dbItem.ReleaseYear = clientItem.ReleaseYear
			dbItem.StreamURL = clientItem.StreamURL
			dbItem.DownloadURL = clientItem.DownloadURL
			dbItem.UpdatedAt = time.Now()

			// Save the updated item
			updatedItem, err := s.itemRepo.Update(ctx, dbItem)
			if err != nil {
				log.Error().Err(err).
					Uint64("itemID", dbItem.ID).
					Str("clientItemID", clientItemID).
					Msg("Failed to update existing item")
				continue
			}

			resultItems = append(resultItems, updatedItem)
			continue
		}

		// Item not found by client ID, try to match by title+year+mediaType
		titleYearKey := fmt.Sprintf("%s_%d_%s", clientItem.Title, clientItem.ReleaseYear, clientItem.Type)
		dbItemByMeta, foundByMeta := titleYearTypeMap[titleYearKey]

		if foundByMeta {
			// Found a match by metadata, add this client to the sync clients
			log.Debug().
				Uint64("clientID", clientID).
				Str("clientItemID", clientItemID).
				Str("title", clientItem.Title).
				Int("year", clientItem.ReleaseYear).
				Msg("Found matching item by metadata, adding client to sync clients")

			// Add this client to the sync clients
			dbItemByMeta.AddSyncClient(clientID, clientType, clientItemID)

			// Update other fields as needed
			dbItemByMeta.Data = clientItem.Data
			dbItemByMeta.UpdatedAt = time.Now()

			// Save the updated item
			updatedItem, err := s.itemRepo.Update(ctx, dbItemByMeta)
			if err != nil {
				log.Error().Err(err).
					Uint64("itemID", dbItemByMeta.ID).
					Str("clientItemID", clientItemID).
					Msg("Failed to update item when adding client to sync clients")
				continue
			}

			resultItems = append(resultItems, updatedItem)
			continue
		}

		// Item not found by either method, create a new one
		log.Debug().
			Uint64("clientID", clientID).
			Str("clientItemID", clientItemID).
			Str("title", clientItem.Title).
			Msg("Creating new media item")

		// Ensure the client info is set
		clientItem.SyncClients = models.SyncClients{} // Clear any existing sync clients
		clientItem.AddSyncClient(clientID, clientType, clientItemID)

		// Create the new item
		newItem, err := s.itemRepo.Create(ctx, clientItem)
		if err != nil {
			log.Error().Err(err).
				Str("clientItemID", clientItemID).
				Str("title", clientItem.Title).
				Msg("Failed to create new media item")
			continue
		}

		resultItems = append(resultItems, newItem)
	}

	log.Info().
		Uint64("clientID", clientID).
		Int("resultCount", len(resultItems)).
		Msg("Successfully processed media items from client")

	return resultItems, nil
}

// GetByClientItemID retrieves a media item by its client-specific ID
func (s *clientMediaItemService[T, U]) GetByClientItemID(ctx context.Context, clientID uint64, itemID string) (*models.MediaItem[U], error) {

	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Str("itemID", itemID).
		Uint64("clientID", clientID).
		Msg("CLIENT: Getting media item by client item ID")

	// First try to get item from the database
	dbItem, err := s.itemRepo.GetByClientItemID(ctx, clientID, itemID)
	if err != nil || dbItem == nil {
		log.Info().
			Str("itemID", itemID).
			Uint64("clientID", clientID).
			Msg("Item not found in database, fetching from client")

		// Item not found in DB, fetch from client
		clientItem, err := s.fetchItemFromClient(ctx, clientID, itemID)
		if err != nil {
			log.Error().Err(err).
				Str("itemID", itemID).
				Uint64("clientID", clientID).
				Msg("Failed to fetch item from client")
			return nil, err
		}

		// Save the item to database
		savedItem, err := s.saveOrUpdateMediaItem(ctx, clientItem)
		if err != nil {
			log.Error().Err(err).
				Str("itemID", itemID).
				Uint64("clientID", clientID).
				Msg("Failed to save item to database")
			return nil, err
		}

		return savedItem, nil
	}

	log.Debug().
		Str("itemID", itemID).
		Uint64("clientID", clientID).
		Uint64("id", dbItem.ID).
		Msg("Media item retrieved from database")

	return dbItem, nil
}

// GetByMultipleClients retrieves all media items associated with any of the specified clients
func (s *clientMediaItemService[T, U]) GetByMultipleClients(ctx context.Context, clientIDs []uint64) ([]*models.MediaItem[U], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Interface("clientIDs", clientIDs).
		Msg("Getting media items by multiple clients")

	// Get items from database
	results, err := s.itemRepo.GetByMultipleClients(ctx, clientIDs)
	if err != nil {
		log.Error().Err(err).
			Interface("clientIDs", clientIDs).
			Msg("Failed to get media items by multiple clients")
		return nil, err
	}

	log.Info().
		Interface("clientIDs", clientIDs).
		Int("count", len(results)).
		Msg("Media items retrieved by multiple clients")

	return results, nil
}

// SearchAcrossClients searches for media items across multiple clients
func (s *clientMediaItemService[T, U]) SearchAcrossClients(ctx context.Context, query types.QueryOptions, clientIDs []uint64) (map[uint64][]*models.MediaItem[U], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Str("query", query.Query).
		Interface("clientIDs", clientIDs).
		Str("type", string(query.MediaType)).
		Msg("Searching media items across clients")

	results := make(map[uint64][]*models.MediaItem[U])

	// Check each client
	for _, clientID := range clientIDs {
		// First, get client configuration
		client, err := s.clientRepo.GetByID(ctx, clientID)
		if err != nil {
			log.Error().Err(err).Uint64("clientID", clientID).Msg("Failed to get client configuration")
			continue
		}

		// Get client instance
		clientInstance, err := s.clientFactory.GetClient(ctx, clientID, client.Config)
		if err != nil {
			log.Error().Err(err).Uint64("clientID", clientID).Msg("Failed to get client instance")
			continue
		}

		// Ensure it's a media client
		mediaClient, ok := clientInstance.(mediaclient.ClientMedia)
		if !ok {
			log.Error().Uint64("clientID", clientID).Msg("Client is not a media client")
			continue
		}

		// Search directly through the client
		searchResults, err := mediaClient.Search(ctx, &query)
		if err != nil {
			log.Error().Err(err).Uint64("clientID", clientID).Msg("Failed to search through client")
			continue
		}

		// Transform the results to the expected type
		// This will depend on the media type we're searching for
		var clientResults []*models.MediaItem[U]

		// Process the search results based on media type
		switch query.MediaType {
		case types.MediaTypeMovie:
			for _, movie := range searchResults.Movies {
				// Type assertion to convert to the expected type
				movieData, ok := interface{}(movie.Data).(U)
				if !ok {
					continue
				}

				mediaItem := models.NewMediaItem(movie.Type, movieData)
				mediaItem.SetClientInfo(clientID, client.Config.GetType(), movie.SyncClients[0].ItemID)

				clientResults = append(clientResults, mediaItem)
			}
		case types.MediaTypeSeries:
			for _, series := range searchResults.Series {
				// Type assertion
				seriesData, ok := interface{}(series.Data).(U)
				if !ok {
					continue
				}

				mediaItem := models.NewMediaItem(series.Type, seriesData)
				mediaItem.SetClientInfo(clientID, client.Config.GetType(), series.SyncClients[0].ItemID)

				clientResults = append(clientResults, mediaItem)
			}
			// Handle other media types similarly
		}

		// Store results for this client
		results[clientID] = clientResults

		// Save results to database for future queries
		for _, item := range clientResults {
			_, err := s.saveOrUpdateMediaItem(ctx, item)
			if err != nil {
				log.Error().Err(err).
					Uint64("clientID", clientID).
					Msg("Failed to save search result to database")
			}
		}
	}

	log.Info().
		Str("query", query.Query).
		Interface("clientIDs", clientIDs).
		Str("type", string(query.MediaType)).
		Int("clientCount", len(results)).
		Msg("Media items found across clients")

	return results, nil
}

// SyncItemBetweenClients creates or updates a mapping between a media item and a target client
func (s *clientMediaItemService[T, U]) SyncItemBetweenClients(ctx context.Context, itemID uint64, sourceClientID uint64, targetClientID uint64, targetItemID string) error {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("itemID", itemID).
		Uint64("sourceClientID", sourceClientID).
		Uint64("targetClientID", targetClientID).
		Str("targetItemID", targetItemID).
		Msg("Syncing item between clients")

	// Get the source item
	// sourceItem, err := s.CoreMediaItemService.GetByID(ctx, itemID)
	// if err != nil {
	// 	log.Error().Err(err).
	// 		Uint64("itemID", itemID).
	// 		Msg("Failed to get source item")
	// 	return fmt.Errorf("failed to get source item: %w", err)
	// }

	// Get the target item from client
	targetItem, err := s.fetchItemFromClient(ctx, targetClientID, targetItemID)
	if err != nil {
		log.Error().Err(err).
			Uint64("targetClientID", targetClientID).
			Str("targetItemID", targetItemID).
			Msg("Failed to fetch target item from client")
		return fmt.Errorf("failed to fetch target item: %w", err)
	}

	// Save the target item to database
	savedTargetItem, err := s.saveOrUpdateMediaItem(ctx, targetItem)
	if err != nil {
		log.Error().Err(err).
			Uint64("targetClientID", targetClientID).
			Str("targetItemID", targetItemID).
			Msg("Failed to save target item")
		return fmt.Errorf("failed to save target item: %w", err)
	}

	// Create a mapping between the source and target items
	err = s.itemRepo.SyncItemBetweenClients(ctx, itemID, sourceClientID, targetClientID, targetItemID)
	if err != nil {
		log.Error().Err(err).
			Uint64("itemID", itemID).
			Uint64("sourceClientID", sourceClientID).
			Uint64("targetClientID", targetClientID).
			Str("targetItemID", targetItemID).
			Msg("Failed to sync item between clients")
		return fmt.Errorf("failed to sync item between clients: %w", err)
	}

	log.Info().
		Uint64("itemID", itemID).
		Uint64("sourceClientID", sourceClientID).
		Uint64("targetClientID", targetClientID).
		Uint64("targetItemID", savedTargetItem.ID).
		Msg("Item synced between clients successfully")

	return nil
}

func (s *clientMediaItemService[T, U]) DeleteClientItem(ctx context.Context, clientID uint64, itemID string) error {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("clientID", clientID).
		Str("itemID", itemID).
		Msg("Deleting client item")

	// Delete from database
	err := s.itemRepo.DeleteClientItem(ctx, clientID, itemID)
	if err != nil {
		log.Error().Err(err).
			Uint64("clientID", clientID).
			Str("itemID", itemID).
			Msg("Failed to delete client item from database")
		return err
	}

	log.Info().
		Uint64("clientID", clientID).
		Str("itemID", itemID).
		Msg("Client item deleted successfully")

	return nil
}

func (s *clientMediaItemService[T, U]) SearchClient(ctx context.Context, clientID uint64, options types.QueryOptions) ([]*models.MediaItem[U], error) {
	log := logger.LoggerFromContext(ctx)
	log.Debug().
		Uint64("clientID", clientID).
		Interface("options", options).
		Msg("Searching media items in client")

	mediaType := types.GetMediaType[U]()
	// Create query options and ensure clientID is set
	options.WithClientID(clientID)
	options.WithMediaType(mediaType)

	// Get client configuration
	client, err := s.clientRepo.GetByID(ctx, clientID)
	if err != nil {
		log.Error().Err(err).Uint64("clientID", clientID).Msg("Failed to get client configuration")
		return nil, fmt.Errorf("failed to get client configuration: %w", err)
	}

	// Get client instance
	clientInstance, err := s.clientFactory.GetClient(ctx, clientID, client.Config)
	if err != nil {
		log.Error().Err(err).Uint64("clientID", clientID).Msg("Failed to get client instance")
		return nil, fmt.Errorf("failed to get client instance: %w", err)
	}

	// Ensure it's a media client
	mediaClient, ok := clientInstance.(mediaclient.ClientMedia)
	if !ok {
		log.Error().Uint64("clientID", clientID).Msg("Client is not a media client")
		return nil, fmt.Errorf("client is not a media client")
	}

	// Perform the search through the client
	searchResults, err := mediaClient.Search(ctx, &options)
	if err != nil {
		log.Error().Err(err).Uint64("clientID", clientID).Msg("Failed to search through client")
		return nil, fmt.Errorf("failed to search through client: %w", err)
	}

	// Process the search results based on media type
	var clientResults []*models.MediaItem[U]
	clientType := client.Config.GetType()

	switch options.MediaType {
	case types.MediaTypeMovie:
		for _, movie := range searchResults.Movies {
			// Try to convert the data to the expected type
			movieData, ok := interface{}(movie.Data).(U)
			if !ok {
				log.Warn().
					Uint64("clientID", clientID).
					Str("itemID", movie.SyncClients[0].ItemID).
					Msg("Failed to convert movie data to expected type, skipping")
				continue
			}

			mediaItem := models.NewMediaItem(movie.Type, movieData)
			mediaItem.Title = movie.Title
			mediaItem.ReleaseDate = movie.ReleaseDate
			mediaItem.ReleaseYear = movie.ReleaseYear
			mediaItem.UpdatedAt = time.Now()
			mediaItem.SetClientInfo(clientID, clientType, movie.SyncClients[0].ItemID)

			clientResults = append(clientResults, mediaItem)
		}
	case types.MediaTypeSeries:
		for _, series := range searchResults.Series {
			seriesData, ok := interface{}(series.Data).(U)
			if !ok {
				log.Warn().
					Uint64("clientID", clientID).
					Str("itemID", series.SyncClients[0].ItemID).
					Msg("Failed to convert series data to expected type, skipping")
				continue
			}

			mediaItem := models.NewMediaItem(series.Type, seriesData)
			mediaItem.Title = series.Title
			mediaItem.ReleaseDate = series.ReleaseDate
			mediaItem.ReleaseYear = series.ReleaseYear
			mediaItem.UpdatedAt = time.Now()
			mediaItem.SetClientInfo(clientID, clientType, series.SyncClients[0].ItemID)

			clientResults = append(clientResults, mediaItem)
		}
	case types.MediaTypeArtist:
		for _, artist := range searchResults.Artists {
			artistData, ok := interface{}(artist.Data).(U)
			if !ok {
				log.Warn().
					Uint64("clientID", clientID).
					Str("itemID", artist.SyncClients[0].ItemID).
					Msg("Failed to convert artist data to expected type, skipping")
				continue
			}

			mediaItem := models.NewMediaItem(artist.Type, artistData)
			mediaItem.Title = artist.Title
			mediaItem.UpdatedAt = time.Now()
			mediaItem.SetClientInfo(clientID, clientType, artist.SyncClients[0].ItemID)

			clientResults = append(clientResults, mediaItem)
		}
	case types.MediaTypeAlbum:
		for _, album := range searchResults.Albums {
			albumData, ok := interface{}(album.Data).(U)
			if !ok {
				log.Warn().
					Uint64("clientID", clientID).
					Str("itemID", album.SyncClients[0].ItemID).
					Msg("Failed to convert album data to expected type, skipping")
				continue
			}

			mediaItem := models.NewMediaItem(album.Type, albumData)
			mediaItem.Title = album.Title
			mediaItem.ReleaseDate = album.ReleaseDate
			mediaItem.ReleaseYear = album.ReleaseYear
			mediaItem.UpdatedAt = time.Now()
			mediaItem.SetClientInfo(clientID, clientType, album.SyncClients[0].ItemID)

			clientResults = append(clientResults, mediaItem)
		}
	case types.MediaTypeTrack:
		for _, track := range searchResults.Tracks {
			trackData, ok := interface{}(track.Data).(U)
			if !ok {
				log.Warn().
					Uint64("clientID", clientID).
					Str("itemID", track.SyncClients[0].ItemID).
					Msg("Failed to convert track data to expected type, skipping")
				continue
			}

			mediaItem := models.NewMediaItem(track.Type, trackData)
			mediaItem.Title = track.Title
			mediaItem.ReleaseDate = track.ReleaseDate
			mediaItem.ReleaseYear = track.ReleaseYear
			mediaItem.UpdatedAt = time.Now()
			mediaItem.SetClientInfo(clientID, clientType, track.SyncClients[0].ItemID)

			clientResults = append(clientResults, mediaItem)
		}
	case types.MediaTypePlaylist:
		for _, playlist := range searchResults.Playlists {
			playlistData, ok := interface{}(playlist.Data).(U)
			if !ok {
				log.Warn().
					Uint64("clientID", clientID).
					Str("itemID", playlist.SyncClients[0].ItemID).
					Msg("Failed to convert playlist data to expected type, skipping")
				continue
			}

			mediaItem := models.NewMediaItem(playlist.Type, playlistData)
			mediaItem.Title = playlist.Title
			mediaItem.UpdatedAt = time.Now()
			mediaItem.SetClientInfo(clientID, clientType, playlist.SyncClients[0].ItemID)

			clientResults = append(clientResults, mediaItem)
		}
	case types.MediaTypeCollection:
		for _, collection := range searchResults.Collections {
			collectionData, ok := interface{}(collection.Data).(U)
			if !ok {
				log.Warn().
					Uint64("clientID", clientID).
					Str("itemID", collection.SyncClients[0].ItemID).
					Msg("Failed to convert collection data to expected type, skipping")
				continue
			}

			mediaItem := models.NewMediaItem(collection.Type, collectionData)
			mediaItem.Title = collection.Title
			mediaItem.UpdatedAt = time.Now()
			mediaItem.SetClientInfo(clientID, clientType, collection.SyncClients[0].ItemID)

			clientResults = append(clientResults, mediaItem)
		}
	}

	// Process and merge the results with database items
	if len(clientResults) == 0 {
		log.Info().
			Uint64("clientID", clientID).
			Msg("No results found in client search")
		return []*models.MediaItem[U]{}, nil
	}

	// Create a map of client item IDs for efficient lookups
	clientItemIDMap := make(map[string]*models.MediaItem[U])
	titleYearMap := make(map[string]*models.MediaItem[U])

	for _, item := range clientResults {
		clientItemID, exists := item.GetClientItemID(clientID)
		if exists && clientItemID != "" {
			clientItemIDMap[clientItemID] = item
		}

		// Also prepare title+year+mediaType keys for secondary matching
		if item.Title != "" && item.ReleaseYear > 0 {
			key := fmt.Sprintf("%s_%d_%s", item.Title, item.ReleaseYear, item.Type)
			titleYearMap[key] = item
		}
	}

	// Fetch database items matching these client items
	dbItemMap := make(map[string]*models.MediaItem[U])
	dbTitleYearMap := make(map[string]*models.MediaItem[U])

	// Get database items by client item IDs
	for clientItemID := range clientItemIDMap {
		dbItem, err := s.itemRepo.GetByClientItemID(ctx, clientID, clientItemID)
		if err == nil && dbItem != nil {
			dbItemMap[clientItemID] = dbItem
		}
	}

	// Fetch items by title+year+mediaType for secondary matching
	mediaTypeItems, err := s.itemRepo.GetByType(ctx, options.MediaType)
	if err == nil {
		for _, item := range mediaTypeItems {
			if item.Title != "" && item.ReleaseYear > 0 {
				key := fmt.Sprintf("%s_%d_%s", item.Title, item.ReleaseYear, item.Type)
				dbTitleYearMap[key] = item
			}
		}
	}

	// Final result items
	resultItems := make([]*models.MediaItem[U], 0, len(clientResults))

	// Process each client item
	for _, clientItem := range clientResults {
		clientItemID, exists := clientItem.GetClientItemID(clientID)
		if !exists || clientItemID == "" {
			// Skip items that don't have a client item ID
			log.Warn().
				Uint64("clientID", clientID).
				Str("title", clientItem.Title).
				Msg("Client item missing client item ID, skipping")
			continue
		}

		// Check if we already have this item in the database by client ID
		dbItem, found := dbItemMap[clientItemID]
		if found {
			// Item exists, update it with new data from client
			dbItem.Data = clientItem.Data
			dbItem.Title = clientItem.Title
			dbItem.ReleaseDate = clientItem.ReleaseDate
			dbItem.ReleaseYear = clientItem.ReleaseYear
			dbItem.StreamURL = clientItem.StreamURL
			dbItem.DownloadURL = clientItem.DownloadURL
			dbItem.UpdatedAt = time.Now()

			// Save the updated item
			updatedItem, err := s.itemRepo.Update(ctx, dbItem)
			if err != nil {
				log.Error().Err(err).
					Uint64("itemID", dbItem.ID).
					Str("clientItemID", clientItemID).
					Msg("Failed to update existing item")
				continue
			}

			resultItems = append(resultItems, updatedItem)
			continue
		}

		// Item not found by client ID, try to match by title+year+mediaType
		titleYearKey := fmt.Sprintf("%s_%d_%s", clientItem.Title, clientItem.ReleaseYear, clientItem.Type)
		dbItemByMeta, foundByMeta := dbTitleYearMap[titleYearKey]

		if foundByMeta {
			// Found a match by metadata, add this client to the sync clients
			log.Debug().
				Uint64("clientID", clientID).
				Str("clientItemID", clientItemID).
				Str("title", clientItem.Title).
				Int("year", clientItem.ReleaseYear).
				Msg("Found matching item by metadata, adding client to sync clients")

			// Add this client to the sync clients
			dbItemByMeta.AddSyncClient(clientID, clientType, clientItemID)

			// Update other fields as needed
			dbItemByMeta.Data = clientItem.Data
			dbItemByMeta.UpdatedAt = time.Now()

			// Save the updated item
			updatedItem, err := s.itemRepo.Update(ctx, dbItemByMeta)
			if err != nil {
				log.Error().Err(err).
					Uint64("itemID", dbItemByMeta.ID).
					Str("clientItemID", clientItemID).
					Msg("Failed to update item when adding client to sync clients")
				continue
			}

			resultItems = append(resultItems, updatedItem)
			continue
		}

		// Item not found by either method, create a new one
		log.Debug().
			Uint64("clientID", clientID).
			Str("clientItemID", clientItemID).
			Str("title", clientItem.Title).
			Msg("Creating new media item")

		// Ensure the client info is set
		clientItem.SyncClients = models.SyncClients{} // Clear any existing sync clients
		clientItem.AddSyncClient(clientID, clientType, clientItemID)

		// Create the new item
		newItem, err := s.itemRepo.Create(ctx, clientItem)
		if err != nil {
			log.Error().Err(err).
				Str("clientItemID", clientItemID).
				Str("title", clientItem.Title).
				Msg("Failed to create new media item")
			continue
		}

		resultItems = append(resultItems, newItem)
	}

	log.Info().
		Uint64("clientID", clientID).
		Int("resultCount", len(resultItems)).
		Msg("Successfully processed search results from client")

	return resultItems, nil
}
