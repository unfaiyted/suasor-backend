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
	"suasor/utils"
)

// syncMovies syncs movies from the client to the database
func (j *MediaSyncJob) syncMovies(ctx context.Context, clientMedia media.ClientMedia, jobRunID uint64, clientID uint64) error {
	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 10, "Fetching movies from client")

	// Check if client supports movies
	movieProvider, ok := clientMedia.(providers.MovieProvider)
	if !ok {
		return fmt.Errorf("client doesn't support movies")
	}

	// Get all movies from the client
	clientType := clientMedia.(clients.Client).GetClientType().AsClientMediaType()
	movies, err := movieProvider.GetMovies(ctx, &mediatypes.QueryOptions{})
	if err != nil {
		return fmt.Errorf("failed to get movies: %w", err)
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 50, fmt.Sprintf("Processing %d movies", len(movies)))

	// Process movies in batches to avoid memory issues
	batchSize := 50
	totalMovies := len(movies)
	processedMovies := 0

	for i := 0; i < totalMovies; i += batchSize {
		end := i + batchSize
		if end > totalMovies {
			end = totalMovies
		}

		movieBatch := movies[i:end]
		err := j.processMovieBatch(ctx, movieBatch, clientID, clientType)
		if err != nil {
			return fmt.Errorf("failed to process movie batch: %w", err)
		}

		processedMovies += len(movieBatch)
		progress := 50 + int(float64(processedMovies)/float64(totalMovies)*50.0)
		j.jobRepo.UpdateJobProgress(ctx, jobRunID, progress, fmt.Sprintf("Processed %d/%d movies", processedMovies, totalMovies))
	}

	// Update job progress
	j.jobRepo.UpdateJobProgress(ctx, jobRunID, 100, fmt.Sprintf("Synced %d movies", totalMovies))

	return nil
}

// processMovieBatch processes a batch of movies and saves them to the database
func (j *MediaSyncJob) processMovieBatch(ctx context.Context, movies []*models.MediaItem[*mediatypes.Movie], clientID uint64, clientType clienttypes.ClientMediaType) error {
	for _, movie := range movies {
		// Skip if movie has no client ID information
		if len(movie.SyncClients) == 0 {
			log.Printf("Skipping movie with no client IDs: %s", movie.Data.Details.Title)
			continue
		}

		// Get the client ID and item ID for lookup
		clientItemID := ""
		for _, cid := range movie.SyncClients {
			if cid.ID == clientID {
				clientItemID = cid.ItemID
				break
			}
		}

		if clientItemID == "" {
			log.Printf("No matching client item ID found for movie: %s", movie.Data.Details.Title)
			continue
		}

		// Check if the movie already exists in the database
		existingMovie, err := j.itemRepos.MovieUserRepo().GetByClientItemID(ctx, clientID, clientItemID)

		if err != nil || existingMovie == nil {
			log.Printf("No matching client item ID found for movie: %s", movie.Data.Details.Title)

			approvedSources := []string{"imdb", "tmdb", "tvdb"}
			filteredExternalIDs := make([]mediatypes.ExternalID, 0, len(movie.Data.Details.ExternalIDs))
			for _, externalID := range movie.Data.Details.ExternalIDs {
				if utils.Contains(approvedSources, externalID.Source) {
					filteredExternalIDs = append(filteredExternalIDs, externalID)
				}
			}

			existingMovie, err = j.itemRepos.MovieUserRepo().GetByExternalIDs(ctx, movie.Data.Details.ExternalIDs)
		}
		// If we cant find it by client Item Id we should check by Title+Year
		if err != nil || existingMovie == nil {
			log.Printf("No matching client item ID found for movie: %s", movie.Data.Details.Title)
			existingMovie, err = j.itemRepos.MovieUserRepo().GetByTitleAndYear(ctx, clientID, movie.Data.Details.Title, movie.Data.Details.ReleaseYear)
		}
		if err == nil {
			log.Printf("Found movie by Title+Year: %s", movie.Data.Details.Title)
			// Movie exists, update it
			existingMovie.Merge(movie)

			// Save the updated movie
			log.Printf("Updating movie: %s", movie.Data.Details.Title)
			_, err = j.itemRepos.MovieUserRepo().Update(ctx, existingMovie)
			if err != nil {
				log.Printf("Error updating movie: %v", err)
				continue
			}
		} else {
			// Movie doesn't exist, create it
			// Set top level title and release fields
			movie.Title = movie.Data.Details.Title
			movie.ReleaseDate = movie.Data.Details.ReleaseDate
			movie.ReleaseYear = movie.Data.Details.ReleaseYear

			// Create the movie
			_, err = j.itemRepos.MovieUserRepo().Create(ctx, movie)
			if err != nil {
				log.Printf("Error creating movie: %v", err)
				continue
			}
		}
	}

	return nil
}
