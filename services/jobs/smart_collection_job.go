package jobs

import (
	"context"
	"fmt"
	"log"
	"time"

	"suasor/client"
	mediatypes "suasor/client/media/types"
	clienttypes "suasor/client/types"
	"suasor/repository"
	"suasor/services/scheduler"
	"suasor/types/models"
)

// SmartCollectionJob automatically organizes media into collections
type SmartCollectionJob struct {
	jobRepo       repository.JobRepository
	userRepo      repository.UserRepository
	configRepo    repository.UserConfigRepository
	movieRepo     repository.MediaItemRepository[*mediatypes.Movie]
	seriesRepo    repository.MediaItemRepository[*mediatypes.Series]
	musicRepo     repository.MediaItemRepository[*mediatypes.Track]
	clientRepos   map[clienttypes.ClientMediaType]interface{}
	clientFactory *client.ClientFactoryService
	aiService     interface{} // Using interface{} to avoid import cycles
}

// NewSmartCollectionJob creates a new smart collection job
func NewSmartCollectionJob(
	jobRepo repository.JobRepository,
	userRepo repository.UserRepository,
	configRepo repository.UserConfigRepository,
	movieRepo repository.MediaItemRepository[*mediatypes.Movie],
	seriesRepo repository.MediaItemRepository[*mediatypes.Series],
	musicRepo repository.MediaItemRepository[*mediatypes.Track],
	embyRepo interface{},
	jellyfinRepo interface{},
	plexRepo interface{},
	subsonicRepo interface{},
	clientFactory *client.ClientFactoryService,
	aiService interface{},
) *SmartCollectionJob {
	clientRepos := map[clienttypes.ClientMediaType]interface{}{
		clienttypes.ClientMediaTypeEmby:     embyRepo,
		clienttypes.ClientMediaTypeJellyfin: jellyfinRepo,
		clienttypes.ClientMediaTypePlex:     plexRepo,
		clienttypes.ClientMediaTypeSubsonic: subsonicRepo,
	}

	return &SmartCollectionJob{
		jobRepo:       jobRepo,
		userRepo:      userRepo,
		configRepo:    configRepo,
		movieRepo:     movieRepo,
		seriesRepo:    seriesRepo,
		musicRepo:     musicRepo,
		clientRepos:   clientRepos,
		clientFactory: clientFactory,
		aiService:     aiService,
	}
}

// Name returns the unique name of the job
func (j *SmartCollectionJob) Name() string {
	return "system.smart.collections"
}

// Schedule returns when the job should next run
func (j *SmartCollectionJob) Schedule() time.Duration {
	// Run weekly by default
	return 7 * 24 * time.Hour
}

// Execute runs the smart collection job
func (j *SmartCollectionJob) Execute(ctx context.Context) error {
	log.Println("Starting smart collection job")

	// Get all users
	users, err := j.userRepo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("error getting users: %w", err)
	}

	// Process each user
	for _, user := range users {
		if err := j.processUserCollections(ctx, user); err != nil {
			log.Printf("Error processing collections for user %s: %v", user.Username, err)
			// Continue with other users even if one fails
			continue
		}
	}

	log.Println("Smart collection job completed")
	return nil
}

// processUserCollections processes smart collections for a single user
func (j *SmartCollectionJob) processUserCollections(ctx context.Context, user models.User) error {
	// Skip inactive users
	if !user.Active {
		log.Printf("Skipping inactive user: %s", user.Username)
		return nil
	}

	// Get user configuration
	config, err := j.configRepo.GetUserConfig(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("error getting user config: %w", err)
	}

	// Check if smart collections are enabled for the user
	if !config.SmartCollectionsEnabled {
		log.Printf("Smart collections not enabled for user: %s", user.Username)
		return nil
	}

	// Create a job run record for this user
	now := time.Now()
	jobRun := &models.JobRun{
		JobName:   j.Name(),
		JobType:   models.JobTypeSystem,
		Status:    models.JobStatusRunning,
		StartTime: &now,
		UserID:    &user.ID,
		Metadata:  fmt.Sprintf(`{"userId":%d,"username":"%s","type":"smartCollections"}`, user.ID, user.Username),
	}

	if err := j.jobRepo.CreateJobRun(ctx, jobRun); err != nil {
		log.Printf("Error creating job run record: %v", err)
		return err
	}

	// Get the user's clients
	clients, err := j.getUserClientMedias(ctx, user.ID)
	if err != nil {
		j.completeJobRun(ctx, jobRun.ID, models.JobStatusFailed, fmt.Sprintf("Error getting media clients: %v", err))
		return err
	}

	// Process each collection type
	var jobError error
	collectionStats := map[string]int{
		"genreCollections":        0,
		"directorCollections":     0,
		"actorCollections":        0,
		"seasonalCollections":     0,
		"customCollections":       0,
		"aiGeneratedCollections":  0,
		"totalCollectionsCreated": 0,
		"totalCollectionsUpdated": 0,
	}

	// Process genre collections
	genreStats, err := j.processGenreCollections(ctx, user.ID, clients)
	if err != nil {
		log.Printf("Error processing genre collections: %v", err)
		jobError = err
		// Continue with other collection types even if one fails
	} else {
		collectionStats["genreCollections"] = genreStats.created + genreStats.updated
		collectionStats["totalCollectionsCreated"] += genreStats.created
		collectionStats["totalCollectionsUpdated"] += genreStats.updated
	}

	// Process director collections
	directorStats, err := j.processDirectorCollections(ctx, user.ID, clients)
	if err != nil {
		log.Printf("Error processing director collections: %v", err)
		if jobError == nil {
			jobError = err
		}
	} else {
		collectionStats["directorCollections"] = directorStats.created + directorStats.updated
		collectionStats["totalCollectionsCreated"] += directorStats.created
		collectionStats["totalCollectionsUpdated"] += directorStats.updated
	}

	// Process actor collections
	actorStats, err := j.processActorCollections(ctx, user.ID, clients)
	if err != nil {
		log.Printf("Error processing actor collections: %v", err)
		if jobError == nil {
			jobError = err
		}
	} else {
		collectionStats["actorCollections"] = actorStats.created + actorStats.updated
		collectionStats["totalCollectionsCreated"] += actorStats.created
		collectionStats["totalCollectionsUpdated"] += actorStats.updated
	}

	// Process seasonal collections
	seasonalStats, err := j.processSeasonalCollections(ctx, user.ID, clients)
	if err != nil {
		log.Printf("Error processing seasonal collections: %v", err)
		if jobError == nil {
			jobError = err
		}
	} else {
		collectionStats["seasonalCollections"] = seasonalStats.created + seasonalStats.updated
		collectionStats["totalCollectionsCreated"] += seasonalStats.created
		collectionStats["totalCollectionsUpdated"] += seasonalStats.updated
	}

	// Process AI-generated collections
	aiStats, err := j.processAIGeneratedCollections(ctx, user.ID, clients)
	if err != nil {
		log.Printf("Error processing AI-generated collections: %v", err)
		if jobError == nil {
			jobError = err
		}
	} else {
		collectionStats["aiGeneratedCollections"] = aiStats.created + aiStats.updated
		collectionStats["totalCollectionsCreated"] += aiStats.created
		collectionStats["totalCollectionsUpdated"] += aiStats.updated
	}

	// Complete the job
	status := models.JobStatusCompleted
	errorMessage := ""
	if jobError != nil {
		status = models.JobStatusFailed
		errorMessage = jobError.Error()
	}

	j.completeJobRun(ctx, jobRun.ID, status, errorMessage)
	return jobError
}

// completeJobRun finalizes a job run with status and error info
func (j *SmartCollectionJob) completeJobRun(ctx context.Context, jobRunID uint64, status models.JobStatus, errorMsg string) {
	if err := j.jobRepo.CompleteJobRun(ctx, jobRunID, status, errorMsg); err != nil {
		log.Printf("Error completing job run: %v", err)
	}
}

// SmartCollectionClientInfo holds information about a media client for collection operations
type SmartCollectionClientInfo struct {
	ClientID   uint64
	ClientType clienttypes.ClientMediaType
	Name       string
	Connection interface{} // The client connection
}

// findMediaItemsByClientID finds items from the given client in the array of client IDs
func findMediaItemsByClientID[T mediatypes.MediaData](items []*models.MediaItem[T], clientID uint64) []*models.MediaItem[T] {
	var result []*models.MediaItem[T]

	for _, item := range items {
		// Check if the item has a client ID matching the requested client
		for _, cid := range item.SyncClients {
			if cid.ID == clientID {
				result = append(result, item)
				break // Only add once even if multiple matches
			}
		}
	}

	return result
}

// CollectionStats holds statistics about collection creation/updates
type CollectionStats struct {
	created      int
	updated      int
	itemsAdded   int
	itemsRemoved int
}

// getUserClientMedias returns all media clients for a user
func (j *SmartCollectionJob) getUserClientMedias(ctx context.Context, userID uint64) ([]SmartCollectionClientInfo, error) {
	// In a real implementation, we would:
	// 1. Query each client repository for clients belonging to this user
	// 2. Create client connections for each client
	// 3. Return the list of clients

	// Mock implementation
	return []SmartCollectionClientInfo{
		{
			ClientID:   1,
			ClientType: clienttypes.ClientMediaTypeEmby,
			Name:       "Home Emby Server",
			Connection: nil, // Would be an actual client in a real implementation
		},
		{
			ClientID:   2,
			ClientType: clienttypes.ClientMediaTypePlex,
			Name:       "Home Plex Server",
			Connection: nil,
		},
	}, nil
}

// processGenreCollections creates and updates genre-based collections
func (j *SmartCollectionJob) processGenreCollections(ctx context.Context, userID uint64, clients []SmartCollectionClientInfo) (CollectionStats, error) {
	stats := CollectionStats{}
	log.Printf("Processing genre collections for user %d", userID)

	// In a real implementation, we would:
	// 1. Get the user's library of movies and TV shows
	// 2. Extract all genres
	// 3. For each significant genre, create or update a collection
	// 4. Add appropriate items to each collection

	// Would use the new ClientIDs array to find movies by client:
	// movies, err := j.movieRepo.GetByUserID(ctx, userID)
	// For each client, get the movies that belong to that client
	// clientMovies := findMediaItemsByClientID(movies, client.ClientID)
	// Extract genres from these movies
	// Create a collection for each major genre

	// Mock implementation
	genres := []string{"Action", "Comedy", "Drama", "Sci-Fi", "Horror"}

	for _, genre := range genres {
		// For each client
		for _, client := range clients {
			// Check if collection exists
			exists, err := j.collectionExistsInClient(ctx, client, fmt.Sprintf("%s Movies", genre))
			if err != nil {
				log.Printf("Error checking if collection exists: %v", err)
				continue
			}

			if exists {
				// Update existing collection
				stats.updated++
			} else {
				// Create new collection
				stats.created++
			}

			// Would find all movies with this genre for this client:
			// movies, err := j.movieRepo.GetByUserID(ctx, userID)
			// clientMovies := findMediaItemsByClientID(movies, client.ClientID)
			// Get genre movies from clientMovies
			// For each movie, add to collection using the ClientIDs array to get the client-specific ID

			// Mock adding items to the collection
			stats.itemsAdded += 5 + (len(genre) % 10) // Just a random number for mock purposes
		}
	}

	log.Printf("Processed genre collections: created %d, updated %d", stats.created, stats.updated)
	return stats, nil
}

// processDirectorCollections creates and updates director-based collections
func (j *SmartCollectionJob) processDirectorCollections(ctx context.Context, userID uint64, clients []SmartCollectionClientInfo) (CollectionStats, error) {
	stats := CollectionStats{}
	log.Printf("Processing director collections for user %d", userID)

	// Similar to genre collections, but for directors
	// Mock implementation
	directors := []string{"Christopher Nolan", "Steven Spielberg", "Martin Scorsese"}

	for _, director := range directors {
		for _, client := range clients {
			exists, err := j.collectionExistsInClient(ctx, client, fmt.Sprintf("%s Collection", director))
			if err != nil {
				log.Printf("Error checking if collection exists: %v", err)
				continue
			}

			if exists {
				stats.updated++
			} else {
				stats.created++
			}

			stats.itemsAdded += 3 + (len(director) % 5)
		}
	}

	log.Printf("Processed director collections: created %d, updated %d", stats.created, stats.updated)
	return stats, nil
}

// processActorCollections creates and updates actor-based collections
func (j *SmartCollectionJob) processActorCollections(ctx context.Context, userID uint64, clients []SmartCollectionClientInfo) (CollectionStats, error) {
	stats := CollectionStats{}
	log.Printf("Processing actor collections for user %d", userID)

	// Similar to director collections, but for actors
	// Mock implementation
	actors := []string{"Tom Hanks", "Meryl Streep", "Leonardo DiCaprio"}

	for _, actor := range actors {
		for _, client := range clients {
			exists, err := j.collectionExistsInClient(ctx, client, fmt.Sprintf("%s Collection", actor))
			if err != nil {
				log.Printf("Error checking if collection exists: %v", err)
				continue
			}

			if exists {
				stats.updated++
			} else {
				stats.created++
			}

			stats.itemsAdded += 4 + (len(actor) % 6)
		}
	}

	log.Printf("Processed actor collections: created %d, updated %d", stats.created, stats.updated)
	return stats, nil
}

// processSeasonalCollections creates and updates seasonal collections
func (j *SmartCollectionJob) processSeasonalCollections(ctx context.Context, userID uint64, clients []SmartCollectionClientInfo) (CollectionStats, error) {
	stats := CollectionStats{}
	log.Printf("Processing seasonal collections for user %d", userID)

	// In a real implementation, we would:
	// 1. Determine which seasonal collections are appropriate based on current date
	// 2. Create or update collections for upcoming seasons/holidays
	// 3. Find appropriate content for each collection

	// Mock implementation - determine current season
	now := time.Now()
	month := now.Month()

	var seasonalCollections []string

	// Add seasonal collections based on current month
	if month >= 9 && month <= 10 {
		seasonalCollections = append(seasonalCollections, "Halloween Favorites")
	}
	if month >= 11 && month <= 12 {
		seasonalCollections = append(seasonalCollections, "Holiday Classics")
	}
	if month >= 3 && month <= 5 {
		seasonalCollections = append(seasonalCollections, "Spring Watching")
	}
	if month >= 6 && month <= 8 {
		seasonalCollections = append(seasonalCollections, "Summer Blockbusters")
	}

	// Always include some standard seasonal collections
	seasonalCollections = append(seasonalCollections, "Oscar Winners")

	for _, collection := range seasonalCollections {
		for _, client := range clients {
			exists, err := j.collectionExistsInClient(ctx, client, collection)
			if err != nil {
				log.Printf("Error checking if collection exists: %v", err)
				continue
			}

			if exists {
				stats.updated++
			} else {
				stats.created++
			}

			stats.itemsAdded += 5 + (len(collection) % 8)
		}
	}

	log.Printf("Processed seasonal collections: created %d, updated %d", stats.created, stats.updated)
	return stats, nil
}

// processAIGeneratedCollections creates and updates AI-generated collections
func (j *SmartCollectionJob) processAIGeneratedCollections(ctx context.Context, userID uint64, clients []SmartCollectionClientInfo) (CollectionStats, error) {
	stats := CollectionStats{}
	log.Printf("Processing AI-generated collections for user %d", userID)

	// In a real implementation, we would:
	// 1. Get the user's watch history and preferences
	// 2. Use AI to suggest interesting collection themes
	// 3. Create or update collections based on these themes
	// 4. Find appropriate content for each collection

	// Mock implementation
	aiCollections := []string{
		"Movies That Will Make You Think",
		"Hidden Gems You Might Have Missed",
		"If You Liked Inception...",
		"Perfect Weekend Binges",
	}

	for _, collection := range aiCollections {
		for _, client := range clients {
			exists, err := j.collectionExistsInClient(ctx, client, collection)
			if err != nil {
				log.Printf("Error checking if collection exists: %v", err)
				continue
			}

			if exists {
				stats.updated++
			} else {
				stats.created++
			}

			stats.itemsAdded += 4 + (len(collection) % 7)
		}
	}

	log.Printf("Processed AI-generated collections: created %d, updated %d", stats.created, stats.updated)
	return stats, nil
}

// collectionExistsInClient checks if a collection exists in a client
func (j *SmartCollectionJob) collectionExistsInClient(ctx context.Context, client SmartCollectionClientInfo, collectionName string) (bool, error) {
	// In a real implementation, we would:
	// 1. Use the client connection to query for the collection
	// 2. Return true if it exists, false otherwise

	// Mock implementation - randomly return true or false
	return (collectionName[0] % 2) == 0, nil
}

// createCollectionInClient creates a collection in a client
func (j *SmartCollectionJob) createCollectionInClient(ctx context.Context, client SmartCollectionClientInfo, collectionName string, items []string) error {
	// In a real implementation, we would:
	// 1. Use the client connection to create the collection
	// 2. Add the specified items to the collection
	//
	// Note: In the actual implementation, we would need to:
	// 1. Get the media items from the repository
	// 2. Extract the client-specific item IDs from each media item's ClientIDs array
	// 3. Create the collection with those client-specific IDs
	//
	// Example implementation sketch:
	// var mediaItems []*models.MediaItem[*mediatypes.Movie]
	// for _, itemID := range items {
	//     // Get the media item (this is simplified)
	//     item, err := j.movieRepo.GetByID(ctx, itemID)
	//     if err != nil {
	//         continue
	//     }
	//     mediaItems = append(mediaItems, item)
	// }
	//
	// // Extract client-specific IDs
	// var clientItemIDs []string
	// for _, item := range mediaItems {
	//     for _, cid := range item.ClientIDs {
	//         if cid.ID == client.ClientID {
	//             clientItemIDs = append(clientItemIDs, cid.ItemID)
	//             break
	//         }
	//     }
	// }
	// Then use clientItemIDs to create the collection in the client

	log.Printf("Created collection '%s' in client %s", collectionName, client.Name)
	return nil
}

// updateCollectionInClient updates a collection in a client
func (j *SmartCollectionJob) updateCollectionInClient(ctx context.Context, client SmartCollectionClientInfo, collectionName string, addItems []string, removeItems []string) error {
	// In a real implementation, we would:
	// 1. Use the client connection to update the collection
	// 2. Add new items and remove old ones
	//
	// Note: In the actual implementation, we would need to:
	// 1. Get the media items from the repository for both add and remove lists
	// 2. Extract the client-specific item IDs from each media item's ClientIDs array
	// 3. Update the collection with those client-specific IDs
	//
	// Example implementation sketch:
	// // Process add items
	// var addMediaItems []*models.MediaItem[*mediatypes.Movie]
	// for _, itemID := range addItems {
	//     item, err := j.movieRepo.GetByID(ctx, itemID)
	//     if err != nil {
	//         continue
	//     }
	//     addMediaItems = append(addMediaItems, item)
	// }
	//
	// // Extract client-specific IDs for items to add
	// var addClientItemIDs []string
	// for _, item := range addMediaItems {
	//     for _, cid := range item.ClientIDs {
	//         if cid.ID == client.ClientID {
	//             addClientItemIDs = append(addClientItemIDs, cid.ItemID)
	//             break
	//         }
	//     }
	// }
	//
	// // Process remove items (similar to above)
	// // ...
	//
	// Then use addClientItemIDs and removeClientItemIDs to update the collection

	log.Printf("Updated collection '%s' in client %s", collectionName, client.Name)
	return nil
}

// SetupSmartCollectionSchedule creates or updates a smart collection schedule
func (j *SmartCollectionJob) SetupSmartCollectionSchedule(ctx context.Context, frequency string) error {
	// Check if job already exists
	existing, err := j.jobRepo.GetJobSchedule(ctx, j.Name())
	if err != nil {
		return fmt.Errorf("error checking for existing job: %w", err)
	}

	// If job exists, update it
	if existing != nil {
		existing.Frequency = frequency
		existing.Enabled = frequency != string(scheduler.FrequencyManual)
		return j.jobRepo.UpdateJobSchedule(ctx, existing)
	}

	// Create a new job schedule
	schedule := &models.JobSchedule{
		JobName:     j.Name(),
		JobType:     models.JobTypeSystem,
		Frequency:   frequency,
		Enabled:     frequency != string(scheduler.FrequencyManual),
		LastRunTime: nil, // Never run yet
	}

	return j.jobRepo.CreateJobSchedule(ctx, schedule)
}

// RunManualCollectionUpdate runs the smart collection job manually
func (j *SmartCollectionJob) RunManualCollectionUpdate(ctx context.Context) error {
	return j.Execute(ctx)
}

// CreateCustomCollection creates a custom collection based on specific criteria
func (j *SmartCollectionJob) CreateCustomCollection(ctx context.Context, userID uint64, collectionName string, criteria map[string]interface{}) error {
	log.Printf("Creating custom collection '%s' for user %d", collectionName, userID)

	// In a real implementation, we would:
	// 1. Get the user's clients
	// 2. Find media items matching the criteria
	// 3. Create the collection in each client
	// 4. Add the matching items to the collection

	// Mock implementation
	return nil
}

// GetCollectionSuggestions gets AI-generated collection suggestions for a user
func (j *SmartCollectionJob) GetCollectionSuggestions(ctx context.Context, userID uint64) ([]map[string]interface{}, error) {
	// In a real implementation, we would:
	// 1. Get the user's watch history and preferences
	// 2. Use AI to generate collection suggestions
	// 3. Return the suggestions with descriptions and sample content

	// Mock implementation
	return []map[string]interface{}{
		{
			"name":               "Sci-Fi Adventures",
			"description":        "A collection of science fiction movies with exploration themes",
			"sampleItems":        []string{"Interstellar", "The Martian", "Arrival"},
			"estimatedItemCount": 15,
		},
		{
			"name":               "Crime Thrillers",
			"description":        "Suspenseful crime movies and series with unexpected twists",
			"sampleItems":        []string{"Silence of the Lambs", "Se7en", "The Departed"},
			"estimatedItemCount": 12,
		},
	}, nil
}

