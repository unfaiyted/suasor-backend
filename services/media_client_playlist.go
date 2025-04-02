package services

import (
	"context"
	"errors"
	"sort"

	"suasor/client"
	"suasor/client/media"
	"suasor/client/media/providers"
	mediatypes "suasor/client/media/types"
	"suasor/client/types"
	"suasor/repository"
	"suasor/types/models"
	"suasor/utils"
)

// MediaClientPlaylistService defines operations for interacting with playlist clients
type MediaClientPlaylistService[T types.ClientConfig] interface {
	GetPlaylistByID(ctx context.Context, userID uint64, clientID uint64, playlistID string) (*models.MediaItem[mediatypes.Playlist], error)
	GetPlaylists(ctx context.Context, userID uint64, count int) ([]models.MediaItem[mediatypes.Playlist], error)
	CreatePlaylist(ctx context.Context, userID uint64, clientID uint64, name string, description string) (*models.MediaItem[mediatypes.Playlist], error)
	UpdatePlaylist(ctx context.Context, userID uint64, clientID uint64, playlistID string, name string, description string) (*models.MediaItem[mediatypes.Playlist], error)
	DeletePlaylist(ctx context.Context, userID uint64, clientID uint64, playlistID string) error
	AddItemToPlaylist(ctx context.Context, userID uint64, clientID uint64, playlistID string, itemID string) error
	RemoveItemFromPlaylist(ctx context.Context, userID uint64, clientID uint64, playlistID string, itemID string) error
	SearchPlaylists(ctx context.Context, userID uint64, query string) ([]models.MediaItem[mediatypes.Playlist], error)
}

type mediaPlaylistService[T types.MediaClientConfig] struct {
	repo    repository.ClientRepository[T]
	factory *client.ClientFactoryService
}

// NewMediaClientPlaylistService creates a new media playlist service
func NewMediaClientPlaylistService[T types.MediaClientConfig](
	repo repository.ClientRepository[T],
	factory *client.ClientFactoryService,
) MediaClientPlaylistService[T] {
	return &mediaPlaylistService[T]{
		repo:    repo,
		factory: factory,
	}
}

// getPlaylistClients gets all playlist clients for a user
func (s *mediaPlaylistService[T]) getPlaylistClients(ctx context.Context, userID uint64) ([]media.MediaClient, error) {
	repo := s.repo
	// Get all media clients for the user
	clients, err := repo.GetByCategory(ctx, types.ClientCategoryMedia, userID)
	if err != nil {
		return nil, err
	}

	var playlistClients []media.MediaClient

	// Filter and instantiate clients that support playlists
	for _, clientConfig := range clients {
		if clientConfig.Config.Data.SupportsPlaylists() {
			clientId := clientConfig.GetID()
			client, err := s.factory.GetClient(ctx, clientId, clientConfig.Config.Data)
			if err != nil {
				// Log error but continue with other clients
				continue
			}
			playlistClients = append(playlistClients, client.(media.MediaClient))
		}
	}

	return playlistClients, nil
}

// getSpecificPlaylistClient gets a specific playlist client
func (s *mediaPlaylistService[T]) getSpecificPlaylistClient(ctx context.Context, userID, clientID uint64) (media.MediaClient, error) {
	log := utils.LoggerFromContext(ctx)

	clientConfig, err := (s.repo).GetByID(ctx, clientID, userID)
	if err != nil {
		return nil, err
	}
	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientType", clientConfig.Config.Data.GetClientType().String()).
		Msg("Retrieved client config")

	if !clientConfig.Config.Data.SupportsPlaylists() {
		log.Warn().
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("clientType", clientConfig.Config.Data.GetClientType().String()).
			Msg("Client does not support playlists")
		return nil, ErrUnsupportedFeature
	}

	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientType", clientConfig.Config.Data.GetClientType().String()).
		Msg("Client supports playlists")

	client, err := s.factory.GetClient(ctx, clientID, clientConfig.Config.Data)
	if err != nil {
		return nil, err
	}
	log.Debug().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("clientType", clientConfig.Config.Data.GetClientType().String()).
		Msg("Retrieved client")
	return client.(media.MediaClient), nil
}

func (s *mediaPlaylistService[T]) GetPlaylistByID(ctx context.Context, userID uint64, clientID uint64, playlistID string) (*models.MediaItem[mediatypes.Playlist], error) {
	client, err := s.getSpecificPlaylistClient(ctx, userID, clientID)
	log := utils.LoggerFromContext(ctx)
	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Msg("Retrieved client")
	if err != nil {
		return nil, err
	}
	log.Info().
		Uint64("userID", userID).
		Uint64("clientID", clientID).
		Str("playlistID", playlistID).
		Msg("Retrieving playlist")

	playlistProvider, ok := client.(providers.PlaylistProvider)
	if !ok {
		log.Warn().
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("playlistID", playlistID).
			Msg("Client does not support playlists")
		return nil, ErrUnsupportedFeature
	}

	// Check if the client supports getting playlist by ID
	if !playlistProvider.SupportsPlaylists() {
		log.Warn().
			Uint64("userID", userID).
			Uint64("clientID", clientID).
			Str("playlistID", playlistID).
			Msg("Client does not support playlists")
		return nil, ErrUnsupportedFeature
	}

	// Get all playlists and find by ID
	options := &mediatypes.QueryOptions{
		Filters: map[string]string{
			"id": playlistID,
		},
	}

	playlists, err := playlistProvider.GetPlaylists(ctx, options)
	if err != nil {
		return nil, err
	}

	// Check if we found any playlists
	if len(playlists) == 0 {
		return nil, errors.New("playlist not found")
	}

	// Return the first matching playlist
	return &playlists[0], nil
}

func (s *mediaPlaylistService[T]) GetPlaylists(ctx context.Context, userID uint64, count int) ([]models.MediaItem[mediatypes.Playlist], error) {
	clients, err := s.getPlaylistClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allPlaylists []models.MediaItem[mediatypes.Playlist]

	for _, client := range clients {
		playlistProvider, ok := client.(providers.PlaylistProvider)
		if !ok || !playlistProvider.SupportsPlaylists() {
			continue
		}

		options := &mediatypes.QueryOptions{
			Limit: count,
		}

		playlists, err := playlistProvider.GetPlaylists(ctx, options)
		if err != nil {
			continue
		}

		allPlaylists = append(allPlaylists, playlists...)
	}

	// Sort by added date
	sort.Slice(allPlaylists, func(i, j int) bool {
		return allPlaylists[i].Data.GetDetails().AddedAt.After(allPlaylists[j].Data.GetDetails().AddedAt)
	})

	// Limit to requested count if specified
	if count > 0 && len(allPlaylists) > count {
		allPlaylists = allPlaylists[:count]
	}

	return allPlaylists, nil
}

func (s *mediaPlaylistService[T]) CreatePlaylist(ctx context.Context, userID uint64, clientID uint64, name string, description string) (*models.MediaItem[mediatypes.Playlist], error) {
	// This is a placeholder. In actual implementations, clients would need to implement
	// a method for creating playlists which isn't currently in the PlaylistProvider interface.
	// For now, we'll return an error indicating this feature isn't implemented.
	return nil, errors.New("create playlist not implemented")
}

func (s *mediaPlaylistService[T]) UpdatePlaylist(ctx context.Context, userID uint64, clientID uint64, playlistID string, name string, description string) (*models.MediaItem[mediatypes.Playlist], error) {
	// This is a placeholder. In actual implementations, clients would need to implement
	// a method for updating playlists which isn't currently in the PlaylistProvider interface.
	// For now, we'll return an error indicating this feature isn't implemented.
	return nil, errors.New("update playlist not implemented")
}

func (s *mediaPlaylistService[T]) DeletePlaylist(ctx context.Context, userID uint64, clientID uint64, playlistID string) error {
	// This is a placeholder. In actual implementations, clients would need to implement
	// a method for deleting playlists which isn't currently in the PlaylistProvider interface.
	// For now, we'll return an error indicating this feature isn't implemented.
	return errors.New("delete playlist not implemented")
}

func (s *mediaPlaylistService[T]) AddItemToPlaylist(ctx context.Context, userID uint64, clientID uint64, playlistID string, itemID string) error {
	// This is a placeholder. In actual implementations, clients would need to implement
	// a method for adding items to playlists which isn't currently in the PlaylistProvider interface.
	// For now, we'll return an error indicating this feature isn't implemented.
	return errors.New("add item to playlist not implemented")
}

func (s *mediaPlaylistService[T]) RemoveItemFromPlaylist(ctx context.Context, userID uint64, clientID uint64, playlistID string, itemID string) error {
	// This is a placeholder. In actual implementations, clients would need to implement
	// a method for removing items from playlists which isn't currently in the PlaylistProvider interface.
	// For now, we'll return an error indicating this feature isn't implemented.
	return errors.New("remove item from playlist not implemented")
}

func (s *mediaPlaylistService[T]) SearchPlaylists(ctx context.Context, userID uint64, query string) ([]models.MediaItem[mediatypes.Playlist], error) {
	clients, err := s.getPlaylistClients(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allPlaylists []models.MediaItem[mediatypes.Playlist]

	for _, client := range clients {
		playlistProvider, ok := client.(providers.PlaylistProvider)
		if !ok || !playlistProvider.SupportsPlaylists() {
			continue
		}

		options := &mediatypes.QueryOptions{
			Query: query,
		}

		playlists, err := playlistProvider.GetPlaylists(ctx, options)
		if err != nil {
			continue
		}

		allPlaylists = append(allPlaylists, playlists...)
	}

	return allPlaylists, nil
}