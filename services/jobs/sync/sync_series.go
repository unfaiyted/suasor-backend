package sync

import (
	"context"
	"fmt"
	"log"
	"suasor/clients"
	"suasor/clients/media"
	"suasor/clients/media/providers"
	mediatypes "suasor/clients/media/types"
	clienttypes "suasor/clients/types"
	"suasor/types/models"
	"suasor/utils/logger"
)

// syncSeries syncs TV series from the client to the database
func (j *MediaSyncJob) syncSeries(ctx context.Context, clientMedia media.ClientMedia, jobRunID uint64, clientID uint64) error {
	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 10, "Fetching series from client")

	// Check if client supports series
	seriesProvider, ok := clientMedia.(providers.SeriesProvider)
	if !ok {
		return fmt.Errorf("client doesn't support series")
	}

	// Get all series from the client
	clientType := clientMedia.(clients.Client).GetClientType()
	series, err := seriesProvider.GetSeries(ctx, &mediatypes.QueryOptions{})
	if err != nil {
		return fmt.Errorf("failed to get series: %w", err)
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 50, fmt.Sprintf("Processing %d series", len(series)))

	// Process series in batches to avoid memory issues
	batchSize := 50
	totalSeries := len(series)
	processedSeries := 0

	for i := 0; i < totalSeries; i += batchSize {
		end := i + batchSize
		if end > totalSeries {
			end = totalSeries
		}

		seriesBatch := series[i:end]
		err := j.processSeriesBatch(ctx, seriesBatch, clientID, clientType)
		if err != nil {
			return fmt.Errorf("failed to process series batch: %w", err)
		}

		processedSeries += len(seriesBatch)
		progress := 50 + int(float64(processedSeries)/float64(totalSeries)*50.0)
		j.jobRepo.UpdateJobProgress(ctx, jobRunID, progress, fmt.Sprintf("Processed %d/%d series", processedSeries, totalSeries))
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 100, fmt.Sprintf("Synced %d series", totalSeries))

	return nil
}

// syncEpisodes syncs TV episodes from the client to the database
func (j *MediaSyncJob) syncEpisodes(ctx context.Context, clientMedia media.ClientMedia, jobRunID uint64, clientID uint64) error {
	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 10, "Fetching episodes from client")

	// Check if client supports episodes
	seriesProvider, ok := clientMedia.(providers.SeriesProvider)
	if !ok {
		return fmt.Errorf("client doesn't support episodes")
	}

	// Get all episodes from the client
	clientType := clientMedia.(clients.Client).GetClientType()

	// Initialize a slice to hold all episodes
	var allEpisodes []*models.MediaItem[*mediatypes.Episode]

	// First get all series
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 15, "Fetching series list")
	allSeries, err := seriesProvider.GetSeries(ctx, &mediatypes.QueryOptions{})
	if err != nil {
		return fmt.Errorf("failed to get series list: %w", err)
	}

	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 20, fmt.Sprintf("Found %d series, fetching episodes", len(allSeries)))

	// Set total items for tracking progress
	totalSeries := len(allSeries)
	j.jobRepo.SetJobTotalItems(ctx, jobRunID, totalSeries)
	processedSeries := 0

	// For each series, get episodes
	for _, series := range allSeries {
		if series.Data == nil || len(series.SyncClients) == 0 {
			// Skip series with no data or no client ID
			log.Printf("Skipping series with missing data")
			continue
		}

		// Find the client item ID for this series
		var seriesID string
		for _, cid := range series.SyncClients {
			if cid.ID == clientID {
				seriesID = cid.ItemID
				break
			}
		}

		if seriesID == "" {
			log.Printf("No matching client item ID found for series: %s", series.Data.Details.Title)
			continue
		}

		// Get seasons for this series
		seasons, err := seriesProvider.GetSeriesSeasons(ctx, seriesID)
		if err != nil {
			log.Printf("Error getting seasons for series %s: %v", series.Data.Details.Title, err)
			continue
		}

		// For each season, get episodes
		for _, season := range seasons {
			if season.Data == nil {
				continue
			}

			seasonNumber := season.Data.Number
			episodes, err := seriesProvider.GetSeriesEpisodesBySeasonNbr(ctx, seriesID, seasonNumber)
			if err != nil {
				log.Printf("Error getting episodes for series %s season %d: %v",
					series.Data.Details.Title, seasonNumber, err)
				continue
			}

			allEpisodes = append(allEpisodes, episodes...)
		}

		// Update progress
		processedSeries++
		progress := 20 + int(float64(processedSeries)/float64(totalSeries)*30.0)
		j.jobRepo.UpdateJobProgress(ctx, jobRunID, progress,
			fmt.Sprintf("Processed %d/%d series, found %d episodes",
				processedSeries, totalSeries, len(allEpisodes)))
		j.jobRepo.IncrementJobProcessedItems(ctx, jobRunID, 1)
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 50,
		fmt.Sprintf("Processing %d episodes", len(allEpisodes)))

	// Process episodes in batches to avoid memory issues
	batchSize := 100
	totalEpisodes := len(allEpisodes)
	processedEpisodes := 0

	for i := 0; i < totalEpisodes; i += batchSize {
		end := i + batchSize
		if end > totalEpisodes {
			end = totalEpisodes
		}

		episodeBatch := allEpisodes[i:end]
		processedEpisodesBatch, err := j.processEpisodeBatch(ctx, episodeBatch, clientID, clientType)
		if err != nil {
			return fmt.Errorf("failed to process episode batch: %w", err)
		}

		processedEpisodes += len(processedEpisodesBatch)
		progress := 50 + int(float64(processedEpisodes)/float64(totalEpisodes)*50.0)
		j.jobRepo.UpdateJobProgress(ctx, jobRunID, progress,
			fmt.Sprintf("Processed %d/%d episodes", processedEpisodes, totalEpisodes))
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 100,
		fmt.Sprintf("Synced %d episodes from %d series", totalEpisodes, totalSeries))

	return nil
}

// processSeriesBatch processes a batch of series and saves them to the database
func (j *MediaSyncJob) processSeriesBatch(ctx context.Context, series []*models.MediaItem[*mediatypes.Series], clientID uint64, clientType clienttypes.ClientType) error {
	// Try to get a series provider for this client to fetch season details
	clientMedia, err := j.getClientMedia(ctx, clientID, clientType)
	if err != nil {
		// Just log the error but continue processing with what we have
		log.Printf("Failed to get media client for season details: %v", err)
	}

	// Cast to series provider if possible
	var seriesProvider providers.SeriesProvider
	if clientMedia != nil {
		if sp, ok := clientMedia.(providers.SeriesProvider); ok {
			seriesProvider = sp
		}
	}

	// processedSeries := make([]*models.MediaItem[*mediatypes.Series], 0, len(series))
	for _, s := range series {

		clientSeriesID, exists := s.GetClientItemID(clientID)
		if !exists {
			log.Printf("Skipping series with no client IDs: %s", s.GetData().Details.Title)
			continue
		}

		// Get Seasons and episodes for this series from the provider
		seasons, err := seriesProvider.GetSeriesSeasons(ctx, clientSeriesID)
		if err != nil {
			log.Printf("Error getting seasons for series %s: %v", s.GetData().Details.Title, err)
		}

		// Get the data for  seasons/episodes updates them to the database and returns
		// the processed seasons/episodes
		processedSeasons, err := j.processSeasonBatch(ctx, seasons, clientSeriesID, clientID, clientType)

		for _, season := range processedSeasons {
			log.Printf("Season: %d", season.GetData().Number)
			seasonNumber := season.GetData().Number
			seasonID := season.ID

			log.Printf("Adding season: number=%d, ID=%d, title=%s",
				seasonNumber, seasonID, season.GetData().Details.Title)

			// Add season episodes first
			s.Data.AddSeasonEpisodeIDs(season.GetData())

			// Now assign the season ID to the corresponding season number
			s.Data.SetSeasonID(seasonNumber, seasonID)
		}

		// Check if the series already exists in the database
		existingSeries, err := j.itemRepos.SeriesUserRepo().GetByClientItemID(ctx, clientID, clientSeriesID)

		// Check by external services IDs
		if err != nil || existingSeries == nil {
			existingSeries, err = j.itemRepos.SeriesUserRepo().GetByExternalIDs(ctx, s.GetData().Details.ExternalIDs)
		}

		// Check by title and year
		if err != nil || existingSeries == nil {
			existingSeries, err = j.itemRepos.SeriesUserRepo().GetByTitleAndYear(ctx, clientID, s.GetData().Details.Title, s.GetData().Details.ReleaseYear)
		}

		if err == nil {
			// Series exists, update it
			existingSeries.Merge(s)
			existingSeries.Data.Merge(s.Data)
			updatedSeries, err := j.itemRepos.SeriesUserRepo().Update(ctx, existingSeries)
			j.updateSeasonsEpisodesShowIDs(ctx, updatedSeries, clientID, clientType)
			if err != nil {
				log.Printf("Error processing season batch: %v", err)
			}

		} else {
			// Series doesn't exist, create it
			// Set top level title and release fields
			s.Title = s.Data.Details.Title
			s.ReleaseDate = s.Data.Details.ReleaseDate
			s.ReleaseYear = s.Data.Details.ReleaseYear

			// Make sure we have genres data initialized
			if s.Data.Genres == nil {
				s.Data.Genres = []string{}
			}

			// Create the series
			createdSeries, err := j.itemRepos.SeriesUserRepo().Create(ctx, s)
			if err != nil {
				log.Printf("Error creating series: %v", err)
				continue
			}
			// update seasons/episodes showIDs
			j.updateSeasonsEpisodesShowIDs(ctx, createdSeries, clientID, clientType)
		}
	}

	return nil
}

func (j *MediaSyncJob) processSeasonBatch(
	ctx context.Context,
	seasons []*models.MediaItem[*mediatypes.Season],
	clientSeriesID string,
	clientID uint64, clientType clienttypes.ClientType) ([]*models.MediaItem[*mediatypes.Season], error) {

	clientMedia, err := j.getClientMedia(ctx, clientID, clientType)
	if err != nil {
		// Just log the error but continue processing with what we have
		log.Printf("Failed to get media client for season details: %v", err)
	}
	// Cast to series provider if possible
	var seriesProvider providers.SeriesProvider
	if clientMedia != nil {
		if sp, ok := clientMedia.(providers.SeriesProvider); ok {
			seriesProvider = sp
		}
	}
	processedSeasons := make([]*models.MediaItem[*mediatypes.Season], 0, len(seasons))
	// seasonEpisodes := make(map[uint][]*models.MediaItem[*mediatypes.Episode])
	for _, season := range seasons {

		seasonClientItemID, exists := season.GetClientItemID(clientID)
		if !exists {
			log.Printf("Skipping season with no client IDs: %s", season.GetData().Details.Title)
			continue
		}

		// Check if the season already exists in the database
		existingSeason, err := j.itemRepos.SeasonUserRepo().GetByClientItemID(ctx, clientID, seasonClientItemID)
		if err != nil || existingSeason == nil {
			// existingSeason, err = j.itemRepos.SeasonUserRepo().GetByExternalIDs(ctx, season.GetData().Details.ExternalIDs)
		}
		if err != nil && existingSeason == nil {
			// TODO: Season Might want Title and Date for better matching.
			// existingSeason, err = j.itemRepos.SeasonUserRepo().GetByTitleAndYear(ctx, clientID, season.GetData().Details.Title, season.GetData().Details.ReleaseYear)
		}

		_, exists = season.GetClientItemID(clientID)
		if !exists {
			log.Printf("Skipping season with no client IDs: %s", season.GetData().Details.Title)
			continue
		}
		episodes, err := seriesProvider.GetSeriesEpisodesBySeasonNbr(ctx, clientSeriesID, season.GetData().Number)
		if err != nil || len(episodes) == 0 {
			log.Printf("Error getting episodes for series %s season %d: %v",
				season.GetData().Details.Title, season.GetData().Number, err)
		}
		processedEpisodes, err := j.processEpisodeBatch(ctx, episodes, clientID, clientType)
		if err != nil {
			log.Printf("Error processing episode batch: %v", err)
		}

		// get EpisodeIDs
		seasonEpisodes := make([]uint64, 0, len(processedEpisodes))
		for _, episode := range processedEpisodes {
			seasonEpisodes = append(seasonEpisodes, episode.ID)
		}

		// Add debug logging to check season numbers
		log.Printf("Processing season number: %d for series %s", season.GetData().Number, season.GetData().Details.Title)

		// Update exisiting season or save new one.
		if existingSeason != nil {
			log.Printf("Existing season found with number: %d", existingSeason.GetData().Number)
			existingSeason.Merge(season)
			existingSeason.Data.MergeEpisodeIDs(seasonEpisodes)

			// Ensure the season number is preserved and not overwritten
			if existingSeason.GetData().Number != season.GetData().Number && season.GetData().Number > 0 {
				log.Printf("Updating season number from %d to %d", existingSeason.GetData().Number, season.GetData().Number)
				existingSeason.GetData().Number = season.GetData().Number
			}

			updatedSeason, err := j.itemRepos.SeasonUserRepo().Update(ctx, existingSeason)
			if err != nil {
				log.Printf("Error updating season: %v", err)
				continue
			}
			processedSeasons = append(processedSeasons, updatedSeason)
		} else {
			// Ensure season has a valid number
			if season.GetData().Number == 0 {
				log.Printf("WARNING: Season has zero number, forcing to 1")
				season.GetData().Number = 1
			}

			season.Data.MergeEpisodeIDs(seasonEpisodes)
			newSeason, err := j.itemRepos.SeasonUserRepo().Create(ctx, season)
			if err != nil {
				log.Printf("Error creating season: %v", err)
				continue
			}
			processedSeasons = append(processedSeasons, newSeason)
		}

	}

	return processedSeasons, nil
}

// processEpisodeBatch processes a batch of episodes and saves them to the database
func (j *MediaSyncJob) processEpisodeBatch(
	ctx context.Context,
	episodes []*models.MediaItem[*mediatypes.Episode],
	clientID uint64, clientType clienttypes.ClientType) ([]*models.MediaItem[*mediatypes.Episode], error) {

	processedEpisodes := make([]*models.MediaItem[*mediatypes.Episode], 0, len(episodes))
	// Finding existing episodes, try to match on them like we do for movies
	for _, episode := range episodes {
		episodeClientID, exists := episode.GetClientItemID(clientID)
		if !exists {
			log.Printf("Skipping episode with no client IDs: %s", episode.GetData().Details.Title)
			continue
		}
		existingEpisode, err := j.itemRepos.EpisodeUserRepo().GetByClientItemID(ctx, clientID, episodeClientID)
		if err != nil || existingEpisode == nil {
			existingEpisode, err = j.itemRepos.EpisodeUserRepo().GetByExternalIDs(ctx, episode.GetData().Details.ExternalIDs)
		}
		if err != nil && existingEpisode == nil {
			// TODO: Episode Might want Title and Date for better matching.
			existingEpisode, err = j.itemRepos.EpisodeUserRepo().GetByTitleAndYear(ctx, clientID, episode.GetData().Details.Title, episode.GetData().Details.ReleaseYear)
		}

		// Update exisiting episode or save new one.
		if existingEpisode != nil {
			existingEpisode.Merge(episode)
			// TODO: Episode specific logic for merging.

			updatedEpisode, err := j.itemRepos.EpisodeUserRepo().Update(ctx, existingEpisode)
			if err != nil {
				log.Printf("Error updating episode: %v", err)
				continue
			}
			processedEpisodes = append(processedEpisodes, updatedEpisode)
		} else {
			_, err = j.itemRepos.EpisodeUserRepo().Create(ctx, episode)
			if err != nil {
				log.Printf("Error creating episode: %v", err)
				continue
			}
			processedEpisodes = append(processedEpisodes, episode)
		}
	}

	return processedEpisodes, nil
}

func (j *MediaSyncJob) updateSeasonsEpisodesShowIDs(
	ctx context.Context,
	series *models.MediaItem[*mediatypes.Series],
	clientID uint64, clientType clienttypes.ClientType) {

	log := logger.LoggerFromContext(ctx)
	log.Info().
		Uint64("clientID", clientID).
		Str("clientType", string(clientType)).
		Msg("Updating seasons/episodes showIDs")

	// Guard against no seasons
	if len(series.Data.Seasons) == 0 {
		log.Warn().Msg("Series has no seasons")
		return
	}

	// Debug logging of all seasons and their IDs
	log.Debug().Msg("Before processing - Season information:")
	for i, season := range series.Data.Seasons {
		log.Debug().
			Int("index", i).
			Int("seasonNumber", season.SeasonNumber).
			Uint64("seasonID", season.SeasonID).
			Int("episodeCount", len(season.EpisodeIDs)).
			Msg("Season info")
	}

	// Loop through each season entry
	for _, season := range series.Data.Seasons {
		// Skip if season ID is not set
		if season.SeasonID == 0 {
			log.Warn().
				Int("seasonNumber", season.SeasonNumber).
				Msg("Season has no ID, skipping")
			continue
		}

		// Update season record
		seasonRecord, err := j.itemRepos.SeasonUserRepo().GetByID(ctx, season.SeasonID)
		if err != nil {
			log.Error().
				Err(err).
				Uint64("seasonID", season.SeasonID).
				Int("seasonNumber", season.SeasonNumber).
				Uint64("clientID", clientID).
				Str("clientType", string(clientType)).
				Msg("Failed to get season by ID")
			continue
		}

		// Set the series ID on the season
		seasonRecord.GetData().SetSeriesID(series.ID)
		_, err = j.itemRepos.SeasonUserRepo().Update(ctx, seasonRecord)
		if err != nil {
			log.Error().
				Err(err).
				Uint64("seasonID", season.SeasonID).
				Int("seasonNumber", season.SeasonNumber).
				Uint64("clientID", clientID).
				Str("clientType", string(clientType)).
				Msg("Failed to update season")
			continue
		}

		// Update all episodes for this season
		for _, episodeID := range season.EpisodeIDs {
			episode, err := j.itemRepos.EpisodeUserRepo().GetByID(ctx, episodeID)
			if err != nil {
				log.Error().
					Err(err).
					Uint64("episodeID", episodeID).
					Int("seasonNumber", season.SeasonNumber).
					Uint64("clientID", clientID).
					Str("clientType", string(clientType)).
					Msg("Failed to get episode by ID")
				continue
			}

			// Set season and series IDs on the episode
			episode.GetData().SetSeriesID(series.ID)
			episode.GetData().SetSeasonNumber(season.SeasonNumber)
			episode.GetData().SetSeasonID(season.SeasonID)

			log.Debug().
				Int("seasonNumber", season.SeasonNumber).
				Uint64("episodeID", episodeID).
				Uint64("seasonID", season.SeasonID).
				Msg("Setting season ID for episode")

			_, err = j.itemRepos.EpisodeUserRepo().Update(ctx, episode)
			if err != nil {
				log.Error().
					Err(err).
					Uint64("episodeID", episodeID).
					Uint64("clientID", clientID).
					Str("clientType", string(clientType)).
					Msg("Failed to update episode")
				continue
			}
		}
	}
}
