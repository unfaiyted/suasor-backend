package wire

import (
	"context"
	"fmt"
	"log"
)

// ExampleHybridWireApproach demonstrates how to use our hybrid approach
// combining Wire-generated handlers for non-generic types and
// manually wired handlers for generic types
func ExampleHybridWireApproach() {
	// Create a context
	ctx := context.Background()

	// You can initialize individual non-generic handlers using Wire
	authHandler, err := InitializeAuthHandler(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize auth handler: %v", err)
	}
	fmt.Println("Auth handler initialized successfully")

	configHandler, err := InitializeConfigHandler(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize config handler: %v", err)
	}
	fmt.Println("Config handler initialized successfully")

	healthHandler, err := InitializeHealthHandler(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize health handler: %v", err)
	}
	fmt.Println("Health handler initialized successfully")

	userHandler, err := InitializeUserHandler(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize user handler: %v", err)
	}
	fmt.Println("User handler initialized successfully")

	// Or you can initialize all handlers at once (combining Wire and manual wiring)
	allHandlers, err := InitializeAllHandlers(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize all handlers: %v", err)
	}

	// Access the system handlers (Wire-generated)
	_ = allHandlers.System.AuthHandler       // From Wire
	_ = allHandlers.System.ConfigHandler     // From Wire
	_ = allHandlers.System.UserHandler       // From Wire
	_ = allHandlers.System.HealthHandler     // From Wire
	_ = allHandlers.System.SearchHandler     // Manually wired
	_ = allHandlers.System.JobHandler        // Manually wired
	_ = allHandlers.System.UserConfigHandler // Manually wired

	// Access the manually wired generic handlers
	_ = allHandlers.Media.MovieHandler
	_ = allHandlers.Media.SeriesHandler 
	_ = allHandlers.MediaData.MovieDataHandler
	_ = allHandlers.Clients.EmbyHandler
	_ = allHandlers.Specialized.RecommendationHandler

	fmt.Println("All handlers initialized successfully")
}

