// examples/di_container_example.go
package main

import (
	"context"
	"fmt"
	"log"
	"suasor/app"
	"suasor/client"
	"suasor/client/media/types"
	clienttypes "suasor/client/types"
	"suasor/database"
	"suasor/repository"
	"suasor/services"
	"time"

	"gorm.io/gorm"
)

func main() {
	// This example demonstrates how to create and use the dependency injection container
	// in the Suasor backend application. It illustrates the three-pronged architecture
	// and shows how dependencies are wired together.

	// 1. Initialize database connection
	db, err := database.ConnectDatabase("sqlite::memory:", false)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 2. Create configuration service
	configService := services.NewConfigService(repository.NewConfigRepository(db))

	// 3. Initialize the dependency injection container
	deps := app.InitializeDependencies(db, configService)

	// 4. Example: Using the container to access services and repositories
	fmt.Println("=== Dependency Injection Container Example ===")
	
	// The container provides a centralized way to access all dependencies
	container := deps.Container()

	// Create a context for our operations
	ctx := context.Background()

	// Example 1: Using the three-pronged architecture for movie services
	exampleUsingThreeProngedArchitecture(ctx, deps)

	// Example 2: Using client services
	exampleUsingClientServices(ctx, deps)

	// Example 3: Adding a new movie through the repository layer
	exampleAddingNewMovie(ctx, deps)

	// Example 4: Using specialized media handlers
	exampleUsingSpecializedHandlers(ctx, deps)
	
	// Example 5: Using the container directly
	exampleUsingContainerDirectly(ctx, container)

	fmt.Println("\n=== End of Example ===")
}

// Example 1: Using the three-pronged architecture for movie services
func exampleUsingThreeProngedArchitecture(ctx context.Context, deps *app.AppDependencies) {
	fmt.Println("\n--- Example 1: Three-Pronged Architecture ---")

	// The three-pronged architecture consists of:
	// 1. Core layer - basic CRUD operations on media items
	// 2. User layer - extends core with user-specific operations
	// 3. Client layer - extends user with client-specific operations

	// Get the movie services from each layer
	coreMovieService := deps.CoreMediaItemServices.MovieCoreService()
	userMovieService := deps.UserMediaItemServices.MovieUserService()
	clientMovieService := deps.ClientMediaItemServices.MovieClientService()

	// Core layer operations (database only)
	movie := &types.Movie{
		BaseItem: types.BaseItem{
			Title: "The Matrix",
			Year:  1999,
		},
	}

	// Example: Create a movie using the core service
	fmt.Println("Creating a movie using Core layer...")
	_, err := coreMovieService.Create(ctx, movie)
	if err != nil {
		fmt.Printf("Error creating movie: %v\n", err)
	}

	// User layer operations (extends core with user-specific data)
	fmt.Println("Accessing movie with User layer...")
	// This would typically include user rating, watch status, etc.
	userMovieService.GetByID(ctx, 1, 1) // userID = 1, movieID = 1

	// Client layer operations (extends user with client-specific data)
	fmt.Println("Accessing movie with Client layer...")
	// This would include client-specific IDs, paths, etc.
	clientMovieService.GetByID(ctx, 1, 1, 1) // userID = 1, clientID = 1, movieID = 1
}

// Example 2: Using client services
func exampleUsingClientServices(ctx context.Context, deps *app.AppDependencies) {
	fmt.Println("\n--- Example 2: Client Services ---")

	// Client services manage integrations with external media servers and services
	// The DI container provides typed access to each client service

	// Get the Jellyfin client service
	jellyfinService := deps.ClientServices.JellyfinService()

	// Example: Getting all Jellyfin clients for a user
	userID := uint(1)
	clients, err := jellyfinService.GetAllForUser(ctx, userID)
	if err != nil {
		fmt.Printf("Error getting Jellyfin clients: %v\n", err)
	} else {
		fmt.Printf("Found %d Jellyfin clients for user\n", len(clients))
	}

	// Example: Creating a new Jellyfin client
	newClient := &clienttypes.JellyfinConfig{
		BaseClientConfig: clienttypes.BaseClientConfig{
			Name:   "Home Jellyfin",
			URL:    "http://jellyfin.local:8096",
			UserID: userID,
		},
		Username: "admin",
		Token:    "abc123",
	}

	fmt.Println("Creating a new Jellyfin client...")
	_, err = jellyfinService.Create(ctx, newClient)
	if err != nil {
		fmt.Printf("Error creating Jellyfin client: %v\n", err)
	}

	// The ClientFactoryService creates appropriate client instances
	fmt.Println("The ClientFactoryService creates appropriate API clients...")
	clientFactory := deps.ClientFactoryService

	// Example: Getting a client by ID (this would create an actual API client)
	if len(clients) > 0 {
		clientID := clients[0].ID
		_, err := clientFactory.GetJellyfinClient(ctx, clientID)
		if err != nil {
			fmt.Printf("Error getting Jellyfin client instance: %v\n", err)
		}
	}
}

// Example 3: Adding a new movie through the repository layer
func exampleAddingNewMovie(ctx context.Context, deps *app.AppDependencies) {
	fmt.Println("\n--- Example 3: Repository Layer ---")

	// The repository layer provides data access
	// We can access repositories directly from the DI container

	// Get the movie repository
	movieRepo := deps.CoreMediaItemRepositories.MovieRepo()

	// Create a new movie
	movie := &types.Movie{
		BaseItem: types.BaseItem{
			Title:        "Inception",
			Year:         2010,
			Overview:     "A thief who steals corporate secrets through the use of dream-sharing technology.",
			RuntimeTicks: 148 * 60 * 1000, // 148 minutes in milliseconds
		},
	}

	fmt.Println("Creating a movie using the repository...")
	_, err := movieRepo.Create(ctx, movie)
	if err != nil {
		fmt.Printf("Error creating movie: %v\n", err)
	}

	// Now fetch the movie we just created
	movies, err := movieRepo.GetAll(ctx, &repository.QueryOptions{
		Limit: 10,
	})
	if err != nil {
		fmt.Printf("Error fetching movies: %v\n", err)
	} else {
		fmt.Printf("Found %d movies\n", len(movies))
		for _, m := range movies {
			fmt.Printf("Movie: %s (%d)\n", m.Title, m.Year)
		}
	}

	// Repositories also support user-specific data
	userMovieRepo := deps.UserRepositoryFactories.MovieUserRepo()
	fmt.Println("User repositories add user-specific operations...")
	_, err = userMovieRepo.GetAllForUser(ctx, 1, &repository.QueryOptions{
		Limit: 10,
	})
	if err != nil {
		fmt.Printf("Error fetching user movies: %v\n", err)
	}

	// Client repositories add client-specific operations
	clientMovieRepo := deps.ClientRepositoryFactories.MovieClientRepo()
	fmt.Println("Client repositories add client-specific operations...")
	_, err = clientMovieRepo.GetAllForClient(ctx, 1, &repository.QueryOptions{
		Limit: 10,
	})
	if err != nil {
		fmt.Printf("Error fetching client movies: %v\n", err)
	}
}

// Example 4: Using specialized media handlers
func exampleUsingSpecializedHandlers(ctx context.Context, deps *app.AppDependencies) {
	fmt.Println("\n--- Example 4: Specialized Handlers ---")

	// The DI container provides access to specialized handlers
	// These handlers implement specific API endpoints
	
	// Music handler example
	musicHandler := deps.SpecializedMediaHandlers.MusicHandler()
	fmt.Println("Music handler provides specialized music-related API endpoints")

	// Series handler example
	seriesHandler := deps.SpecializedMediaHandlers.SeriesSpecificHandler()
	fmt.Println("Series handler provides series-specific API endpoints")

	// Client-specific media handlers
	movieHandler := deps.ClientMediaHandlers.JellyfinMovieHandler()
	fmt.Println("Client-specific handlers provide client-type specific endpoints")

	// These handlers would be used in the router to handle HTTP requests
	fmt.Println("These handlers are used in the router to handle HTTP requests")

	// Just referencing these variables to avoid unused variable warnings
	_ = musicHandler
	_ = seriesHandler
	_ = movieHandler
}

// The following shows how to register services with a scheduler
// This is a simplified example - the actual implementation uses a more complex job system
func exampleWithScheduler(deps *app.AppDependencies) {
	fmt.Println("\n--- Example 5: Scheduler Integration ---")

	// Create a simple scheduler (this is simplified)
	type SimpleScheduler struct {
		jobs map[string]func()
	}

	scheduler := &SimpleScheduler{
		jobs: make(map[string]func()),
	}

	// Register a job to run every hour
	registerJob := func(name string, interval time.Duration, job func()) {
		scheduler.jobs[name] = job
		fmt.Printf("Registered job '%s' to run every %v\n", name, interval)
	}

	// Example: Register a job to synchronize media library
	mediaService := deps.CoreMediaItemServices.MovieCoreService()
	registerJob("media-sync", 1*time.Hour, func() {
		fmt.Println("Running media sync job...")
		// This would call the service methods to sync media
		_ = mediaService // Just to avoid unused variable warning
	})

	// Example: Register a job to clean up old data
	registerJob("cleanup", 24*time.Hour, func() {
		fmt.Println("Running cleanup job...")
		// This would clean up old data
	})

	// In a real application, the scheduler would run these jobs according to their schedules
	fmt.Println("In a real app, the scheduler would run jobs at specified intervals")
}

// Example 5: Using the container directly to access services
func exampleUsingContainerDirectly(ctx context.Context, c *container.Container) {
	fmt.Println("\n--- Example 5: Direct Container Access ---")
	
	// Get services directly from the container using the generic Get method
	healthService, err := container.GetTyped[services.HealthService](c)
	if err != nil {
		fmt.Printf("Error getting health service: %v\n", err)
	} else {
		fmt.Println("Successfully retrieved health service from container")
		status := healthService.Check(ctx)
		fmt.Printf("Health check status: %v\n", status.IsHealthy)
	}
	
	// Get repository directly from the container
	movieRepo, err := container.GetTyped[repository.MediaItemRepository[*types.Movie]](c)
	if err != nil {
		fmt.Printf("Error getting movie repository: %v\n", err)
	} else {
		fmt.Println("Successfully retrieved movie repository from container")
		// Use the repository
		movies, err := movieRepo.GetAll(ctx, &repository.QueryOptions{Limit: 5})
		if err != nil {
			fmt.Printf("Error fetching movies: %v\n", err)
		} else {
			fmt.Printf("Found %d movies using direct container access\n", len(movies))
		}
	}
	
	// You can also use MustGet if you're confident the service exists
	// This will panic if the service doesn't exist
	fmt.Println("MustGet can be used when you're confident the dependency exists")
	try := func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Recovered from panic: %v\n", r)
			}
		}()
		// This might panic if the service doesn't exist
		configService := container.MustGet[services.ConfigService](c)
		fmt.Printf("Got config service: %v\n", configService != nil)
	}
	try()
}