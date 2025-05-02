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
	historyProvider, ok := clientMedia.(providers.HistoryProvider)
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
			Msg("Client reports it doesn't support history - skipping")
		return nil
	}

	// Get watch history items from the client
	historyList, err := historyProvider.GetPlayHistory(ctx, &mediatypes.QueryOptions{})
	if err != nil {
		log.Error().
			Err(err).
			Uint64("clientID", clientID).
			Str("clientType", string(clientMedia.GetClientType())).
			Msg("Failed to get watch history from client")
		return fmt.Errorf("failed to get watch history: %w", err)
	}

	if historyList == nil {
		log.Warn().
			Uint64("clientID", clientID).
			Str("clientType", string(clientMedia.GetClientType())).
			Msg("No history data returned from client")
		return nil
	}

	// Count total items to process
	totalItems := historyList.GetTotalItems()
	log.Info().
		Int("totalItems", totalItems).
		Msg("Processing watch history items")

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 30,
		fmt.Sprintf("Processing %d history items", totalItems))

	// Process movies
	processedItems := 0
	if len(historyList.Movies) > 0 {
		log.Info().
			Int("count", len(historyList.Movies)).
			Msg("Processing movie history items")

		for uuid, movieHistory := range historyList.Movies {
			if err := j.processMovieHistory(ctx, clientID, movieHistory); err != nil {
				log.Warn().
					Err(err).
					Str("uuid", uuid).
					Str("title", movieHistory.Item.Title).
					Msg("Error processing movie history")
				continue
			}
			processedItems++

			// Update progress periodically (every 10 items)
			if processedItems%10 == 0 {
				progress := 30 + int(float64(processedItems)/float64(totalItems)*70.0)
				j.jobRepo.UpdateJobProgress(ctx, jobRunID, progress,
					fmt.Sprintf("Processed %d/%d history items", processedItems, totalItems))
			}
		}
	}

	// Process series
	if len(historyList.Series) > 0 {
		log.Info().
			Int("count", len(historyList.Series)).
			Msg("Processing series history items")

		for uuid, seriesHistory := range historyList.Series {
			if err := j.processSeriesHistory(ctx, clientID, seriesHistory); err != nil {
				log.Warn().
					Err(err).
					Str("uuid", uuid).
					Str("title", seriesHistory.Item.Title).
					Msg("Error processing series history")
				continue
			}
			processedItems++

			// Update progress periodically
			if processedItems%10 == 0 {
				progress := 30 + int(float64(processedItems)/float64(totalItems)*70.0)
				j.jobRepo.UpdateJobProgress(ctx, jobRunID, progress,
					fmt.Sprintf("Processed %d/%d history items", processedItems, totalItems))
			}
		}
	}

	// Process episodes
	if len(historyList.Episodes) > 0 {
		log.Info().
			Int("count", len(historyList.Episodes)).
			Msg("Processing episode history items")

		for uuid, episodeHistory := range historyList.Episodes {
			if err := j.processEpisodeHistory(ctx, clientID, episodeHistory); err != nil {
				log.Warn().
					Err(err).
					Str("uuid", uuid).
					Str("title", episodeHistory.Item.Title).
					Msg("Error processing episode history")
				continue
			}
			processedItems++

			// Update progress periodically
			if processedItems%10 == 0 {
				progress := 30 + int(float64(processedItems)/float64(totalItems)*70.0)
				j.jobRepo.UpdateJobProgress(ctx, jobRunID, progress,
					fmt.Sprintf("Processed %d/%d history items", processedItems, totalItems))
			}
		}
	}

	// Process tracks
	if len(historyList.Tracks) > 0 {
		log.Info().
			Int("count", len(historyList.Tracks)).
			Msg("Processing track history items")

		for uuid, trackHistory := range historyList.Tracks {
			if err := j.processTrackHistory(ctx, clientID, trackHistory); err != nil {
				log.Warn().
					Err(err).
					Str("uuid", uuid).
					Str("title", trackHistory.Item.Title).
					Msg("Error processing track history")
				continue
			}
			processedItems++

			// Update progress periodically
			if processedItems%10 == 0 {
				progress := 30 + int(float64(processedItems)/float64(totalItems)*70.0)
				j.jobRepo.UpdateJobProgress(ctx, jobRunID, progress,
					fmt.Sprintf("Processed %d/%d history items", processedItems, totalItems))
			}
		}
	}

	// Process albums if present
	if len(historyList.Albums) > 0 {
		log.Info().
			Int("count", len(historyList.Albums)).
			Msg("Processing album history items")

		for uuid, albumHistory := range historyList.Albums {
			if err := j.processAlbumHistory(ctx, clientID, albumHistory); err != nil {
				log.Warn().
					Err(err).
					Str("uuid", uuid).
					Str("title", albumHistory.Item.Title).
					Msg("Error processing album history")
				continue
			}
			processedItems++

			// Update progress periodically
			if processedItems%10 == 0 {
				progress := 30 + int(float64(processedItems)/float64(totalItems)*70.0)
				j.jobRepo.UpdateJobProgress(ctx, jobRunID, progress,
					fmt.Sprintf("Processed %d/%d history items", processedItems, totalItems))
			}
		}
	}

	// Update job progress to complete
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 100,
		fmt.Sprintf("Synced %d history items", totalItems))

	return nil
}

// processMovieHistory processes a movie history item and saves it to the database
func (j *MediaSyncJob) processMovieHistory(ctx context.Context, clientID uint64, historyItem *models.UserMediaItemData[*mediatypes.Movie]) error {
	log := logger.LoggerFromContext(ctx)

	// Get client item ID from the media item
	clientItemID, exists := historyItem.Item.GetClientItemID(clientID)
	if !exists {
		log.Debug().
			Uint64("clientID", clientID).
			Str("title", historyItem.Item.Title).
			Msg("No client item ID found for movie history - skipping")
		return fmt.Errorf("no client item ID found for movie history")
	}
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

// processSeriesHistory processes a series history item and saves it to the database
func (j *MediaSyncJob) processSeriesHistory(ctx context.Context, clientID uint64, historyItem *models.UserMediaItemData[*mediatypes.Series]) error {
	log := logger.LoggerFromContext(ctx)

	// Get client item ID from the media item
	clientItemID, exists := historyItem.Item.GetClientItemID(clientID)
	if !exists {
		log.Debug().
			Uint64("clientID", clientID).
			Str("title", historyItem.Item.Title).
			Msg("No client item ID found for series history - skipping")
		return fmt.Errorf("no client item ID found for series history")
	}
	if clientItemID == "" {
		log.Debug().
			Uint64("clientID", clientID).
			Str("title", historyItem.Item.Title).
			Msg("No client item ID found for series history - attempting to fetch from client")

		// Try to get series from client if needed
		item, _, err := j.fetchSeriesFromClient(ctx, clientID, historyItem.Item)
		if err != nil {
			return fmt.Errorf("failed to fetch series from client: %w", err)
		}

		// Now we have the series in our database, create the history record
		historyRecord := models.UserMediaItemData[*mediatypes.Series]{
			MediaItemID:      item.ID,
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
		historyRecord.Associate(item)

		// Save to database
		newHistoryItem, err := j.dataRepos.SeriesDataRepo().Create(ctx, &historyRecord)
		if err != nil {
			return fmt.Errorf("failed to save series history: %w", err)
		}

		log.Debug().
			Uint64("historyID", newHistoryItem.ID).
			Time("watchedAt", newHistoryItem.PlayedAt).
			Float64("percentage", newHistoryItem.PlayedPercentage).
			Msg("Saved series watch history (new item)")

		return nil
	}

	// Look up the series in our database
	existingSeries, err := j.itemRepos.SeriesUserRepo().GetByClientItemID(ctx, clientID, clientItemID)
	if err != nil {
		log.Debug().
			Str("clientItemID", clientItemID).
			Uint64("clientID", clientID).
			Msg("Series not found in database - attempting to fetch from client")

		// Try to get item from client
		item, _, err := j.fetchSeriesFromClient(ctx, clientID, historyItem.Item)
		if err != nil {
			return fmt.Errorf("failed to fetch series from client: %w", err)
		}

		// Now we have the series in our database, create the history record
		historyRecord := models.UserMediaItemData[*mediatypes.Series]{
			MediaItemID:      item.ID,
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
		historyRecord.Associate(item)

		// Save to database
		newHistoryItem, err := j.dataRepos.SeriesDataRepo().Create(ctx, &historyRecord)
		if err != nil {
			return fmt.Errorf("failed to save series history: %w", err)
		}

		log.Debug().
			Uint64("historyID", newHistoryItem.ID).
			Time("watchedAt", newHistoryItem.PlayedAt).
			Float64("percentage", newHistoryItem.PlayedPercentage).
			Msg("Saved series watch history (new item)")

		return nil
	}

	// Series exists in our database, create history record
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
		Msg("Saved series watch history (existing item)")

	return nil
}

// processEpisodeHistory processes an episode history item and saves it to the database
func (j *MediaSyncJob) processEpisodeHistory(ctx context.Context, clientID uint64, historyItem *models.UserMediaItemData[*mediatypes.Episode]) error {
	log := logger.LoggerFromContext(ctx)

	// Get client item ID from the media item
	clientItemID, exists := historyItem.Item.GetClientItemID(clientID)
	if !exists {
		log.Debug().
			Uint64("clientID", clientID).
			Str("title", historyItem.Item.Title).
			Msg("No client item ID found for episode history - skipping")
		return fmt.Errorf("no client item ID found for episode history")
	}
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
func (j *MediaSyncJob) processTrackHistory(ctx context.Context, clientID uint64, historyItem *models.UserMediaItemData[*mediatypes.Track]) error {
	log := logger.LoggerFromContext(ctx)

	// Get client item ID from the media item
	clientItemID, exists := historyItem.Item.GetClientItemID(clientID)
	if !exists {
		log.Debug().
			Uint64("clientID", clientID).
			Str("title", historyItem.Item.Title).
			Msg("No client item ID found for track history - skipping")
		return fmt.Errorf("no client item ID found for track history")
	}

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

func (j *MediaSyncJob) processAlbumHistory(ctx context.Context, clientID uint64, historyItem *models.UserMediaItemData[*mediatypes.Album]) error {
	log := logger.LoggerFromContext(ctx)

	// Get client item ID from the media item
	clientItemID, exists := historyItem.Item.GetClientItemID(clientID)
	if !exists {
		log.Debug().
			Uint64("clientID", clientID).
			Str("title", historyItem.Item.Title).
			Msg("No client item ID found for album history - skipping")
		return fmt.Errorf("no client item ID found for album history")
	}

	if clientItemID == "" {
		return fmt.Errorf("no client item ID found")
	}

	// Look up the album in our database
	existingAlbum, err := j.itemRepos.AlbumUserRepo().GetByClientItemID(ctx, clientID, clientItemID)
	if err != nil {
		log.Warn().
			Str("clientItemID", clientItemID).
			Uint64("clientID", clientID).
			Msg("Album not found in database for history - consider running music sync first")
		return err
	}

	// Create history record
	historyRecord := models.UserMediaItemData[*mediatypes.Album]{
		MediaItemID:      existingAlbum.ID,
		Type:             mediatypes.MediaTypeAlbum,
		PlayedAt:         historyItem.PlayedAt,
		LastPlayedAt:     historyItem.LastPlayedAt,
		PlayedPercentage: historyItem.PlayedPercentage,
		PlayCount:        historyItem.PlayCount,
		PositionSeconds:  historyItem.PositionSeconds,
		DurationSeconds:  historyItem.DurationSeconds,
		Completed:        historyItem.PlayedPercentage >= 90, // Mark as completed if listened to 90% or more
	}

	// Associate with the album item
	historyRecord.Associate(existingAlbum)

	// Save to database
	item, err := j.dataRepos.AlbumDataRepo().Create(ctx, &historyRecord)
	if err != nil {
		return fmt.Errorf("failed to save album history: %w", err)
	}

	log.Debug().
		Str("albumTitle", existingAlbum.Title).
		Time("playedAt", item.PlayedAt).
		Uint64("historyID", item.ID).
		Float64("percentage", historyItem.PlayedPercentage).
		Msg("Saved album play history")

	return nil
}

// fetchMovieFromClient fetches a movie from a client and stores it in the database
func (j *MediaSyncJob) fetchMovieFromClient(ctx context.Context, clientID uint64, item *models.MediaItem[*mediatypes.Movie]) (*models.MediaItem[*mediatypes.Movie], *mediatypes.Movie, error) {
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
	clientItemID, exists := item.GetClientItemID(clientID)
	if !exists {
		return nil, nil, fmt.Errorf("no client item ID found")
	}
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

// fetchSeriesFromClient fetches a series from a client and stores it in the database
func (j *MediaSyncJob) fetchSeriesFromClient(ctx context.Context, clientID uint64, item *models.MediaItem[*mediatypes.Series]) (*models.MediaItem[*mediatypes.Series], *mediatypes.Series, error) {
	log := logger.LoggerFromContext(ctx)

	// Get client info
	clientMedia, err := j.getClientMedia(ctx, clientID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get client: %w", err)
	}

	// Check if client supports series
	seriesProvider, ok := clientMedia.(providers.SeriesProvider)
	if !ok {
		return nil, nil, fmt.Errorf("client doesn't support series")
	}

	// Get client item ID from the media item
	clientItemID, exists := item.GetClientItemID(clientID)
	if !exists {
		return nil, nil, fmt.Errorf("no client item ID found")
	}
	if clientItemID == "" {
		return nil, nil, fmt.Errorf("no client item ID found")
	}

	// Fetch the series from the client
	series, err := seriesProvider.GetSeriesByID(ctx, clientItemID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get series from client: %w", err)
	}

	// Set top level title and release fields
	series.Title = series.Data.Details.Title
	series.ReleaseDate = series.Data.Details.ReleaseDate
	series.ReleaseYear = series.Data.Details.ReleaseYear

	// Create the series in our database
	savedSeries, err := j.itemRepos.SeriesUserRepo().Create(ctx, series)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to save series: %w", err)
	}

	log.Info().
		Str("title", series.Data.Details.Title).
		Str("clientItemID", clientItemID).
		Uint64("clientID", clientID).
		Msg("Successfully fetched and saved series from client")

	return savedSeries, series.Data, nil
}
