// sync_list_example.go
package sync

import (
	"context"
	"fmt"
	
	"suasor/clients/media"
	"suasor/clients/media/providers"
	mediatypes "suasor/clients/media/types"
	"suasor/utils/logger"
)

// ListSyncExample demonstrates how to use the adapter pattern for syncing lists
type ListSyncExample struct {
	// Dependencies would be injected here
}

// SyncPlaylists syncs playlists between clients using the adapter pattern
func (j *ListSyncExample) SyncPlaylists(
	ctx context.Context,
	sourceClient media.ClientMedia,
	targetClient media.ClientMedia,
) error {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("sourceClientID", sourceClient.GetClientID()).
		Uint64("targetClientID", targetClient.GetClientID()).
		Msg("Syncing playlists between clients using adapters")
	
	// Check if clients support playlists
	sourcePlaylistProvider, ok := sourceClient.(providers.PlaylistProvider)
	if !ok || !sourcePlaylistProvider.SupportsPlaylists() {
		return fmt.Errorf("source client doesn't support playlists")
	}
	
	targetPlaylistProvider, ok := targetClient.(providers.PlaylistProvider)
	if !ok || !targetPlaylistProvider.SupportsPlaylists() {
		return fmt.Errorf("target client doesn't support playlists")
	}
	
	// Create playlist adapters
	sourceAdapter := providers.NewPlaylistListAdapter(sourcePlaylistProvider)
	targetAdapter := providers.NewPlaylistListAdapter(targetPlaylistProvider)
	
	// Create the sync adapter for playlists
	syncAdapter := providers.NewListSyncAdapter(sourceAdapter, targetAdapter)
	
	// Sync playlists using the adapter
	err := syncAdapter.SyncLists(ctx, &mediatypes.QueryOptions{})
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to sync playlists")
		return err
	}
	
	log.Info().
		Msg("Successfully synced playlists")
	
	return nil
}

// SyncCollections syncs collections between clients using the adapter pattern
func (j *ListSyncExample) SyncCollections(
	ctx context.Context,
	sourceClient media.ClientMedia,
	targetClient media.ClientMedia,
) error {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("sourceClientID", sourceClient.GetClientID()).
		Uint64("targetClientID", targetClient.GetClientID()).
		Msg("Syncing collections between clients using adapters")
	
	// Check if clients support collections
	sourceCollectionProvider, ok := sourceClient.(providers.CollectionProvider)
	if !ok || !sourceCollectionProvider.SupportsCollections() {
		return fmt.Errorf("source client doesn't support collections")
	}
	
	targetCollectionProvider, ok := targetClient.(providers.CollectionProvider)
	if !ok || !targetCollectionProvider.SupportsCollections() {
		return fmt.Errorf("target client doesn't support collections")
	}
	
	// Create collection adapters
	sourceAdapter := providers.NewCollectionListAdapter(sourceCollectionProvider)
	targetAdapter := providers.NewCollectionListAdapter(targetCollectionProvider)
	
	// Create the sync adapter for collections
	syncAdapter := providers.NewListSyncAdapter(sourceAdapter, targetAdapter)
	
	// Sync collections using the adapter
	err := syncAdapter.SyncLists(ctx, &mediatypes.QueryOptions{})
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to sync collections")
		return err
	}
	
	log.Info().
		Msg("Successfully synced collections")
	
	return nil
}

// SyncAllLists syncs both playlists and collections between clients
func (j *ListSyncExample) SyncAllLists(
	ctx context.Context,
	sourceClient media.ClientMedia,
	targetClient media.ClientMedia,
) error {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("sourceClientID", sourceClient.GetClientID()).
		Uint64("targetClientID", targetClient.GetClientID()).
		Msg("Syncing all lists between clients")
	
	// Try to sync playlists
	err := j.SyncPlaylists(ctx, sourceClient, targetClient)
	if err != nil {
		log.Warn().
			Err(err).
			Msg("Failed to sync playlists, continuing with collections")
	}
	
	// Try to sync collections
	err = j.SyncCollections(ctx, sourceClient, targetClient)
	if err != nil {
		log.Warn().
			Err(err).
			Msg("Failed to sync collections")
	}
	
	return nil
}

// For working with mixed list types, we can use the TypedListProvider approach:

// SyncMixedLists demonstrates using specialized type adapters for each list type
// This avoids the need to directly work with the ListData interface
func (j *ListSyncExample) SyncMixedLists(
	ctx context.Context,
	sourceClient media.ClientMedia,
	targetClient media.ClientMedia,
) error {
	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("sourceClientID", sourceClient.GetClientID()).
		Uint64("targetClientID", targetClient.GetClientID()).
		Msg("Syncing mixed lists between clients")
	
	// Get the appropriate adapter methods
	sourcePlaylistAdapter, hasSourcePlaylists := getPlaylistAdapter(sourceClient)
	targetPlaylistAdapter, hasTargetPlaylists := getPlaylistAdapter(targetClient)
	
	sourceCollectionAdapter, hasSourceCollections := getCollectionAdapter(sourceClient)
	targetCollectionAdapter, hasTargetCollections := getCollectionAdapter(targetClient)
	
	// If both clients support playlists, sync them
	if hasSourcePlaylists && hasTargetPlaylists {
		// Create the sync adapter for playlists
		playlistSyncAdapter := providers.NewListSyncAdapter(sourcePlaylistAdapter, targetPlaylistAdapter)
		
		// Sync playlists
		err := playlistSyncAdapter.SyncLists(ctx, nil)
		if err != nil {
			log.Warn().
				Err(err).
				Msg("Failed to sync playlists")
		}
	}
	
	// If both clients support collections, sync them
	if hasSourceCollections && hasTargetCollections {
		// Create the sync adapter for collections
		collectionSyncAdapter := providers.NewListSyncAdapter(sourceCollectionAdapter, targetCollectionAdapter)
		
		// Sync collections
		err := collectionSyncAdapter.SyncLists(ctx, nil)
		if err != nil {
			log.Warn().
				Err(err).
				Msg("Failed to sync collections")
		}
	}
	
	return nil
}

// Helper functions to get type-specific adapters

// getPlaylistAdapter returns a playlist list adapter if supported
func getPlaylistAdapter(client media.ClientMedia) (providers.ListProvider[*mediatypes.Playlist], bool) {
	// First check if the client directly provides an adapter
	if provider, ok := client.(interface {
		AsPlaylistListProvider() providers.ListProvider[*mediatypes.Playlist]
	}); ok {
		return provider.AsPlaylistListProvider(), true
	}
	
	// Otherwise check if it implements PlaylistProvider
	if provider, ok := client.(providers.PlaylistProvider); ok && provider.SupportsPlaylists() {
		return providers.NewPlaylistListAdapter(provider), true
	}
	
	return nil, false
}

// getCollectionAdapter returns a collection list adapter if supported
func getCollectionAdapter(client media.ClientMedia) (providers.ListProvider[*mediatypes.Collection], bool) {
	// First check if the client directly provides an adapter
	if provider, ok := client.(interface {
		AsCollectionListProvider() providers.ListProvider[*mediatypes.Collection]
	}); ok {
		return provider.AsCollectionListProvider(), true
	}
	
	// Otherwise check if it implements CollectionProvider
	if provider, ok := client.(providers.CollectionProvider); ok && provider.SupportsCollections() {
		return providers.NewCollectionListAdapter(provider), true
	}
	
	return nil, false
}