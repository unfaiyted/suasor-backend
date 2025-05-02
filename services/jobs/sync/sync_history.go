package sync

import (
	"context"
	"fmt"
	"suasor/clients/media"
	"suasor/clients/media/providers"
	mediatypes "suasor/clients/media/types"
	"suasor/types/models"
	"suasor/utils/logger"
)

// syncHistory syncs watch history from a client to the database
func (j *MediaSyncJob) syncHistory(ctx context.Context, clientMedia media.ClientMedia, jobRunID uint64, clientID uint64) error {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 10, "Fetching watch history from client")

	// Check if client supports history
	historyProvider, ok := clientMedia.(providers.HistoryProvider[mediatypes.MediaData])
	if !ok {
		log.Warn().
			Uint64("clientID", clientID).
			Str("clientType", string(clientMedia.GetClientType())).
			Msg("Client doesn't support history")
		return fmt.Errorf("client doesn't support history")
	}

	// Skip if client specifically reports it doesn't support history
	if !historyProvider.SupportsHistory() {
		log.Warn().
			Uint64("clientID", clientID).
			Str("clientType", string(clientMedia.GetClientType())).
			Msg("Client reports it doesn't support  history - skipping")
		return nil
	}

	// Get watch history items from the client
	playHistory, err := historyProvider.GetPlayHistory(ctx, &mediatypes.QueryOptions{})
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", clientID).
			Str("clientType", string(clientMedia.GetClientType())).
			Msg("Failed to get watch history from client")
		return fmt.Errorf("failed to get watch history: %w", err)
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 30,
		fmt.Sprintf("Processing %d history items", len(playHistory)))

	// Process history items in batches to avoid memory issues
	batchSize := 50
	totalItems := len(playHistory)
	processedItems := 0

	for i := 0; i < totalItems; i += batchSize {
		end := i + batchSize
		if end > totalItems {
			end = totalItems
		}

		historyBatch := playHistory[i:end]
		err := j.processHistoryBatch(ctx, historyBatch, clientID, jobRunID)
		if err != nil {
			log.Error().
				Err(err).
				Uint64("clientID", clientID).
				Str("clientType", string(clientMedia.GetClientType())).
				Int("batchStart", i).
				Int("batchEnd", end).
				Msg("Failed to process history batch")
			continue // Continue processing other batches even if one fails
		}

		processedItems += len(historyBatch)
		progress := 30 + int(float64(processedItems)/float64(totalItems)*70.0)
		j.jobRepo.UpdateJobProgress(ctx, jobRunID, progress,
			fmt.Sprintf("Processed %d/%d history items", processedItems, totalItems))
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 100,
		fmt.Sprintf("Synced %d history items", totalItems))

	return nil
}

// processHistoryBatch processes a batch of history items and saves them to the database
func (j *MediaSyncJob) processHistoryBatch(ctx context.Context, historyItems []*models.UserMediaItemData[mediatypes.MediaData], clientID uint64, jobRunID uint64) error {
	log := logger.LoggerFromContext(ctx)

	for _, historyItem := range historyItems {
		// Skip invalid items
		if historyItem == nil || historyItem.Item == nil {
			log.Warn().Msg("Skipping invalid history item with no data")
			continue
		}

		// Create/update history based on media type
		switch historyItem.Item.Type {
		case mediatypes.MediaTypeMovie:
			if err := j.processMovieHistory(ctx, clientID, historyItem); err != nil {
				log.Warn().
					Err(err).
					Str("mediaType", string(historyItem.Item.Type)).
					Str("title", historyItem.Item.Title).
					Msg("Error processing movie history")
				continue
			}

		case mediatypes.MediaTypeSeries:
			if err := j.processSeriesHistory(ctx, clientID, historyItem); err != nil {
				log.Warn().
					Err(err).
					Str("mediaType", string(historyItem.Item.Type)).
					Str("title", historyItem.Item.Title).
					Msg("Error processing series history")
				continue
			}

		case mediatypes.MediaTypeEpisode:
			if err := j.processEpisodeHistory(ctx, clientID, historyItem); err != nil {
				log.Warn().
					Err(err).
					Str("mediaType", string(historyItem.Item.Type)).
					Str("title", historyItem.Item.Title).
					Msg("Error processing episode history")
				continue
			}

		case mediatypes.MediaTypeTrack:
			if err := j.processMusicHistory(ctx, clientID, historyItem); err != nil {
				log.Warn().
					Err(err).
					Str("mediaType", string(historyItem.Item.Type)).
					Str("title", historyItem.Item.Title).
					Msg("Error processing music history")
				continue
			}

		default:
			log.Debug().
				Str("mediaType", string(historyItem.Item.Type)).
				Str("title", historyItem.Item.Title).
				Msg("Unsupported media type in history - skipping")
			continue
		}
	}

	return nil
}

// processMovieHistory processes a movie history item and saves it to the database
func (j *MediaSyncJob) processMovieHistory(ctx context.Context, clientID uint64, historyItem *models.UserMediaItemData[mediatypes.MediaData]) error {
	log := logger.LoggerFromContext(ctx)

	// Get client item ID from the media item
	clientItemID := getClientItemID(historyItem.Item, clientID)
	if clientItemID == "" {
		log.Debug().
			Uint64("clientID", clientID).
			Str("title", historyItem.Item.Title).
			Msg("No client item ID found for movie history - attempting to fetch from client")

		// Try to get updated item from client
		// This path will be invoked if we have history for an item that isn't yet in our database
		item, _, err := j.fetchMovieFromClient(ctx, clientID, historyItem.Item)
		if err != nil {
			return fmt.Errorf("failed to fetch movie from client: %w", err)
		}

		// Now we have the movie in our database, create the history record
		historyRecord := models.UserMediaItemData[*mediatypes.Movie]{
			MediaItemID:      item.ID,
			Type:             mediatypes.MediaTypeMovie,
			PlayedAt:         historyItem.PlayedAt,
			LastPlayedAt:     historyItem.LastPlayedAt,
			PlayedPercentage: historyItem.PlayedPercentage,
			PlayCount:        historyItem.PlayCount,
			PositionSeconds:  historyItem.PositionSeconds,
			DurationSeconds:  historyItem.DurationSeconds,
			Completed:        historyItem.PlayedPercentage >= 90, // Mark as completed if watched 90% or more
		}

		// Associate with the movie item
		historyRecord.Associate(item)

		// Save to database
		newHistoryItem, err := j.dataRepos.MovieDataRepo().Create(ctx, &historyRecord)
		if err != nil {
			return fmt.Errorf("failed to save movie history: %w", err)
		}

		log.Debug().
			Uint64("historyID", newHistoryItem.ID).
			Time("watchedAt", newHistoryItem.PlayedAt).
			Float64("percentage", newHistoryItem.PlayedPercentage).
			Msg("Saved movie watch history (new item)")

		return nil
	}

	// Look up the movie in our database
	existingMovie, err := j.itemRepos.MovieUserRepo().GetByClientItemID(ctx, clientID, clientItemID)
	if err != nil {
		log.Debug().
			Str("clientItemID", clientItemID).
			Uint64("clientID", clientID).
			Msg("Movie not found in database - attempting to fetch from client")

		// Try to get item from client
		item, _, err := j.fetchMovieFromClient(ctx, clientID, historyItem.Item)
		if err != nil {
			return fmt.Errorf("failed to fetch movie from client: %w", err)
		}

		// Now we have the movie in our database, create the history record
		historyRecord := models.UserMediaItemData[*mediatypes.Movie]{
			MediaItemID:      item.ID,
			Type:             mediatypes.MediaTypeMovie,
			PlayedAt:         historyItem.PlayedAt,
			LastPlayedAt:     historyItem.LastPlayedAt,
			PlayedPercentage: historyItem.PlayedPercentage,
			PlayCount:        historyItem.PlayCount,
			PositionSeconds:  historyItem.PositionSeconds,
			DurationSeconds:  historyItem.DurationSeconds,
			Completed:        historyItem.PlayedPercentage >= 90, // Mark as completed if watched 90% or more
		}

		// Associate with the movie item
		historyRecord.Associate(item)

		// Save to database
		newHistoryItem, err := j.dataRepos.MovieDataRepo().Create(ctx, &historyRecord)
		if err != nil {
			return fmt.Errorf("failed to save movie history: %w", err)
		}

		log.Debug().
			Uint64("historyID", newHistoryItem.ID).
			Time("watchedAt", newHistoryItem.PlayedAt).
			Float64("percentage", newHistoryItem.PlayedPercentage).
			Msg("Saved movie watch history (new item)")

		return nil
	}

	// Movie exists in our database, create history record
	historyRecord := models.UserMediaItemData[*mediatypes.Movie]{
		MediaItemID:      existingMovie.ID,
		Type:             mediatypes.MediaTypeMovie,
		PlayedAt:         historyItem.PlayedAt,
		LastPlayedAt:     historyItem.LastPlayedAt,
		PlayedPercentage: historyItem.PlayedPercentage,
		PlayCount:        historyItem.PlayCount,
		PositionSeconds:  historyItem.PositionSeconds,
		DurationSeconds:  historyItem.DurationSeconds,
		Completed:        historyItem.PlayedPercentage >= 90, // Mark as completed if watched 90% or more
	}

	// Associate with the movie item
	historyRecord.Associate(existingMovie)

	// Save to database
	item, err := j.dataRepos.MovieDataRepo().Create(ctx, &historyRecord)
	if err != nil {
		return fmt.Errorf("failed to save movie history: %w", err)
	}

	log.Debug().
		Uint64("historyID", item.ID).
		Time("watchedAt", item.PlayedAt).
		Float64("percentage", item.PlayedPercentage).
		Msg("Saved movie watch history (existing item)")

	return nil
}

// fetchMovieFromClient fetches a movie from a client and stores it in the database
func (j *MediaSyncJob) fetchMovieFromClient(ctx context.Context, clientID uint64, item *models.MediaItem[mediatypes.MediaData]) (*models.MediaItem[*mediatypes.Movie], *mediatypes.Movie, error) {
	log := logger.LoggerFromContext(ctx)

	// Get client info
	clientMedia, err := j.getClientMedia(ctx, clientID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get client: %w", err)
	}

	// Check if client supports movies
	movieProvider, ok := clientMedia.(providers.MovieProvider)
	if !ok {
		return nil, nil, fmt.Errorf("client doesn't support movies")
	}

	// Get client item ID from the media item
	clientItemID := getClientItemID(item, clientID)
	if clientItemID == "" {
		return nil, nil, fmt.Errorf("no client item ID found")
	}

	// Fetch the movie from the client
	movie, err := movieProvider.GetMovieByID(ctx, clientItemID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get movie from client: %w", err)
	}

	// Set top level title and release fields
	movie.Title = movie.Data.Details.Title
	movie.ReleaseDate = movie.Data.Details.ReleaseDate
	movie.ReleaseYear = movie.Data.Details.ReleaseYear

	// Create the movie in our database
	savedMovie, err := j.itemRepos.MovieUserRepo().Create(ctx, movie)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to save movie: %w", err)
	}

	log.Info().
		Str("title", movie.Data.Details.Title).
		Str("clientItemID", clientItemID).
		Uint64("clientID", clientID).
		Msg("Successfully fetched and saved movie from client")

	return savedMovie, movie.Data, nil
}

// processSeriesHistory processes a series history item and saves it to the database
func (j *MediaSyncJob) processSeriesHistory(ctx context.Context, clientID uint64, historyItem *models.UserMediaItemData[mediatypes.MediaData]) error {
	log := logger.LoggerFromContext(ctx)

	// Get client item ID from the media item
	clientItemID := getClientItemID(historyItem.Item, clientID)
	if clientItemID == "" {
		log.Debug().
			Uint64("clientID", clientID).
			Str("title", historyItem.Item.Title).
			Msg("No client item ID found for series history - attempting to fetch from client")

		// Try to get series from client if needed
		// TODO: Implement fetchSeriesFromClient similar to fetchMovieFromClient
		return fmt.Errorf("fetching series from client not yet implemented")
	}

	// Look up the series in our database
	existingSeries, err := j.itemRepos.SeriesUserRepo().GetByClientItemID(ctx, clientID, clientItemID)
	if err != nil {
		log.Warn().
			Str("clientItemID", clientItemID).
			Uint64("clientID", clientID).
			Msg("Series not found in database for history - consider running series sync first")
		return err
	}

	// Create history record
	historyRecord := models.UserMediaItemData[*mediatypes.Series]{
		MediaItemID:      existingSeries.ID,
		Type:             mediatypes.MediaTypeSeries,
		PlayedAt:         historyItem.PlayedAt,
		LastPlayedAt:     historyItem.LastPlayedAt,
		PlayedPercentage: historyItem.PlayedPercentage,
		PlayCount:        historyItem.PlayCount,
		PositionSeconds:  historyItem.PositionSeconds,
		DurationSeconds:  historyItem.DurationSeconds,
		Completed:        historyItem.PlayedPercentage >= 90, // Mark as completed if watched 90% or more
	}

	// Associate with the series item
	historyRecord.Associate(existingSeries)

	// Save to database
	item, err := j.dataRepos.SeriesDataRepo().Create(ctx, &historyRecord)
	if err != nil {
		return fmt.Errorf("failed to save series history: %w", err)
	}

	log.Debug().
		Str("seriesTitle", existingSeries.Title).
		Time("watchedAt", item.PlayedAt).
		Uint64("historyID", item.ID).
		Float64("percentage", historyItem.PlayedPercentage).
		Msg("Saved series watch history")

	return nil
}

// processEpisodeHistory processes an episode history item and saves it to the database
func (j *MediaSyncJob) processEpisodeHistory(ctx context.Context, clientID uint64, historyItem *models.UserMediaItemData[mediatypes.MediaData]) error {
	log := logger.LoggerFromContext(ctx)

	// Get client item ID from the media item
	clientItemID := getClientItemID(historyItem.Item, clientID)
	if clientItemID == "" {
		log.Debug().
			Uint64("clientID", clientID).
			Str("title", historyItem.Item.Title).
			Msg("No client item ID found for episode history - skipping")
		return fmt.Errorf("no client item ID found for episode history")
	}

	// Look up the episode in our database
	existingEpisode, err := j.itemRepos.EpisodeUserRepo().GetByClientItemID(ctx, clientID, clientItemID)
	if err != nil {
		log.Warn().
			Str("clientItemID", clientItemID).
			Uint64("clientID", clientID).
			Msg("Episode not found in database for history - consider running series sync first")
		return err
	}

	// Create history record
	historyRecord := models.UserMediaItemData[*mediatypes.Episode]{
		MediaItemID:      existingEpisode.ID,
		Type:             mediatypes.MediaTypeEpisode,
		PlayedAt:         historyItem.PlayedAt,
		LastPlayedAt:     historyItem.LastPlayedAt,
		PlayedPercentage: historyItem.PlayedPercentage,
		PlayCount:        historyItem.PlayCount,
		PositionSeconds:  historyItem.PositionSeconds,
		DurationSeconds:  historyItem.DurationSeconds,
		Completed:        historyItem.PlayedPercentage >= 90, // Mark as completed if watched 90% or more
	}

	// Associate with the episode item
	historyRecord.Associate(existingEpisode)

	// Save to database
	item, err := j.dataRepos.EpisodeDataRepo().Create(ctx, &historyRecord)
	if err != nil {
		return fmt.Errorf("failed to save episode history: %w", err)
	}

	log.Debug().
		Str("episodeTitle", existingEpisode.Title).
		Time("watchedAt", item.PlayedAt).
		Uint64("historyID", item.ID).
		Float64("percentage", historyItem.PlayedPercentage).
		Msg("Saved episode watch history")

	return nil
}

// processMusicHistory processes a music track history item and saves it to the database
func (j *MediaSyncJob) processMusicHistory(ctx context.Context, clientID uint64, historyItem *models.UserMediaItemData[mediatypes.MediaData]) error {
	log := logger.LoggerFromContext(ctx)

	// Get client item ID from the media item
	clientItemID := getClientItemID(historyItem.Item, clientID)
	if clientItemID == "" {
		log.Debug().
			Uint64("clientID", clientID).
			Str("title", historyItem.Item.Title).
			Msg("No client item ID found for track history - skipping")
		return fmt.Errorf("no client item ID found for track history")
	}

	// Look up the track in our database
	existingTrack, err := j.itemRepos.TrackUserRepo().GetByClientItemID(ctx, clientID, clientItemID)
	if err != nil {
		log.Warn().
			Str("clientItemID", clientItemID).
			Uint64("clientID", clientID).
			Msg("Track not found in database for history - consider running music sync first")
		return err
	}

	// Create history record
	historyRecord := models.UserMediaItemData[*mediatypes.Track]{
		MediaItemID:      existingTrack.ID,
		Type:             mediatypes.MediaTypeTrack,
		PlayedAt:         historyItem.PlayedAt,
		LastPlayedAt:     historyItem.LastPlayedAt,
		PlayedPercentage: historyItem.PlayedPercentage,
		PlayCount:        historyItem.PlayCount,
		PositionSeconds:  historyItem.PositionSeconds,
		DurationSeconds:  historyItem.DurationSeconds,
		Completed:        historyItem.PlayedPercentage >= 90, // Mark as completed if listened to 90% or more
	}

	// Associate with the track item
	historyRecord.Associate(existingTrack)

	// Save to database
	item, err := j.dataRepos.TrackDataRepo().Create(ctx, &historyRecord)
	if err != nil {
		return fmt.Errorf("failed to save track history: %w", err)
	}

	log.Debug().
		Str("trackTitle", existingTrack.Title).
		Time("playedAt", item.PlayedAt).
		Uint64("historyID", item.ID).
		Float64("percentage", historyItem.PlayedPercentage).
		Msg("Saved track play history")

	return nil
}

// getClientItemID helper function to get client item ID from a media item
func getClientItemID(item *models.MediaItem[mediatypes.MediaData], clientID uint64) string {
	if item == nil || len(item.SyncClients) == 0 {
		return ""
	}

	for _, cid := range item.SyncClients {
		if cid.ID == clientID {
			return cid.ItemID
		}
	}

	return ""
}
