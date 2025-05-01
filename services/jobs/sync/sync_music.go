package sync

import (
	"context"
	"fmt"
	"suasor/clients"
	"suasor/clients/media"
	"suasor/clients/media/providers"
	mediatypes "suasor/clients/media/types"
)

// syncMusic syncs music tracks from the client to the database
func (j *MediaSyncJob) syncMusic(ctx context.Context, clientMedia media.ClientMedia, jobRunID uint64, clientID uint64) error {
	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 10, "Fetching music from client")

	// Check if client supports music
	musicProvider, ok := clientMedia.(providers.MusicProvider)
	if !ok {
		return fmt.Errorf("client doesn't support music")
	}

	// Get all tracks from the client
	clientType := clientMedia.(clients.Client).GetClientType()
	tracks, err := musicProvider.GetMusic(ctx, &mediatypes.QueryOptions{Limit: 100, Offset: 0})
	if err != nil {
		return fmt.Errorf("failed to get tracks: %w", err)
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 50, fmt.Sprintf("Processing %d tracks", len(tracks)))

	// Process tracks in batches to avoid memory issues
	batchSize := 100
	totalTracks := len(tracks)
	processedTracks := 0

	for i := 0; i < totalTracks; i += batchSize {
		end := i + batchSize
		if end > totalTracks {
			end = totalTracks
		}

		trackBatch := tracks[i:end]
		processedTrackBatch, err := j.processIndependentTrackBatch(ctx, trackBatch, clientID, clientType)
		if err != nil {
			return fmt.Errorf("failed to process track batch: %w", err)
		}

		processedTracks += len(processedTrackBatch)
		progress := 50 + int(float64(processedTracks)/float64(totalTracks)*50.0)
		j.jobRepo.UpdateJobProgress(ctx, jobRunID, progress, fmt.Sprintf("Processed %d/%d tracks", processedTracks, totalTracks))
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 100, fmt.Sprintf("Synced %d tracks", totalTracks))

	return nil
}

// syncAlbums syncs music albums from the client to the database
func (j *MediaSyncJob) syncAlbums(ctx context.Context, clientMedia media.ClientMedia, jobRunID uint64, clientID uint64) error {
	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 10, "Fetching albums from client")

	// Check if client supports albums
	musicProvider, ok := clientMedia.(providers.MusicProvider)
	if !ok {
		return fmt.Errorf("client doesn't support albums")
	}

	// Get all albums from the client
	clientType := clientMedia.(clients.Client).GetClientType()
	albums, err := musicProvider.GetMusicAlbums(ctx, &mediatypes.QueryOptions{})
	if err != nil {
		return fmt.Errorf("failed to get albums: %w", err)
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 50, fmt.Sprintf("Processing %d albums", len(albums)))

	// Process albums in batches to avoid memory issues
	batchSize := 50
	totalAlbums := len(albums)
	processedAlbums := 0

	for i := 0; i < totalAlbums; i += batchSize {
		end := i + batchSize
		if end > totalAlbums {
			end = totalAlbums
		}

		albumBatch := albums[i:end]
		albumsWithTracks, err := j.processIndependentAlbumBatch(ctx, albumBatch, clientID, clientType)
		if err != nil {
			return fmt.Errorf("failed to process album batch: %w", err)
		}

		processedAlbums += len(albumsWithTracks)
		progress := 50 + int(float64(processedAlbums)/float64(totalAlbums)*50.0)
		j.jobRepo.UpdateJobProgress(ctx, jobRunID, progress, fmt.Sprintf("Processed %d/%d albums", processedAlbums, totalAlbums))
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 100, fmt.Sprintf("Synced %d albums", totalAlbums))

	return nil
}

// syncArtists syncs music artists from the client to the database
func (j *MediaSyncJob) syncArtists(ctx context.Context, clientMedia media.ClientMedia, jobRunID uint64, clientID uint64) error {
	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 10, "Fetching artists from client")

	// Check if client supports artists
	musicProvider, ok := clientMedia.(providers.MusicProvider)
	if !ok {
		return fmt.Errorf("client doesn't support artists")
	}

	// Get all artists from the client
	clientType := clientMedia.(clients.Client).GetClientType()
	artists, err := musicProvider.GetMusicArtists(ctx, &mediatypes.QueryOptions{})
	if err != nil {
		return fmt.Errorf("failed to get artists: %w", err)
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 50, fmt.Sprintf("Processing %d artists", len(artists)))

	// Process artists in batches to avoid memory issues
	batchSize := 50
	totalArtists := len(artists)
	processedArtists := 0

	for i := 0; i < totalArtists; i += batchSize {
		end := i + batchSize
		if end > totalArtists {
			end = totalArtists
		}

		artistBatch := artists[i:end]
		err := j.processArtistBatch(ctx, artistBatch, clientID, clientType)
		if err != nil {
			return fmt.Errorf("failed to process artist batch: %w", err)
		}

		processedArtists += len(artistBatch)
		progress := 50 + int(float64(processedArtists)/float64(totalArtists)*50.0)
		j.jobRepo.UpdateJobProgress(ctx, jobRunID, progress, fmt.Sprintf("Processed %d/%d artists", processedArtists, totalArtists))
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 100, fmt.Sprintf("Synced %d artists", totalArtists))

	return nil
}
