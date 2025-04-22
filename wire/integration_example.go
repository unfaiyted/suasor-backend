package wire

// import (
// 	"context"
// 	"log"
// 	"suasor/app"
//
// 	"github.com/gin-gonic/gin"
// 	"gorm.io/gorm"
// 	"suasor/services"
// )
//
// // IntegrateWithApplication shows how to use wire-generated handlers
// // with the existing application
// func IntegrateWithApplication() {
// 	// Create a new router
// 	router := gin.Default()
//
// 	// Initialize the application handlers
// 	handlers, err := InitializeAllHandlers(context.Background())
// 	if err != nil {
// 		log.Fatalf("Failed to initialize handlers: %v", err)
// 	}
//
// 	// Register the system handlers
// 	registerSystemHandlers(router, handlers)
//
// 	// Register the media handlers
// 	registerMediaHandlers(router, handlers)
//
// 	// Start the server
// 	router.Run(":8080")
// }
//
// // registerSystemHandlers registers the system handlers with the router
// func registerSystemHandlers(router *gin.Engine, handlers ApplicationHandlers) {
// 	// Auth routes
// 	authRoutes := router.Group("/auth")
// 	{
// 		authRoutes.POST("/login", handlers.System.AuthHandler.Login)
// 		authRoutes.POST("/refresh", handlers.System.AuthHandler.RefreshToken)
// 		authRoutes.POST("/logout", handlers.System.AuthHandler.Logout)
// 	}
//
// 	// User routes
// 	userRoutes := router.Group("/users")
// 	{
// 		userRoutes.GET("/me", handlers.System.UserHandler.GetCurrentUser)
// 		userRoutes.PUT("/me", handlers.System.UserHandler.UpdateCurrentUser)
// 	}
//
// 	// Config routes
// 	configRoutes := router.Group("/config")
// 	{
// 		configRoutes.GET("/", handlers.System.ConfigHandler.GetConfig)
// 		configRoutes.PUT("/", handlers.System.ConfigHandler.UpdateConfig)
// 	}
//
// 	// Health routes
// 	router.GET("/health", handlers.System.HealthHandler.CheckHealth)
//
// 	// Search routes
// 	searchRoutes := router.Group("/search")
// 	{
// 		searchRoutes.GET("/", handlers.System.SearchHandler.Search)
// 		searchRoutes.GET("/suggestions", handlers.System.SearchHandler.GetSuggestions)
// 	}
//
// 	// People routes
// 	peopleRoutes := router.Group("/people")
// 	{
// 		peopleRoutes.GET("/:id", handlers.Media.PeopleHandler.GetPersonByID)
// 		peopleRoutes.GET("/search", handlers.Media.PeopleHandler.SearchPeople)
// 	}
// }
//
// // registerMediaHandlers registers the media handlers with the router
// func registerMediaHandlers(router *gin.Engine, handlers ApplicationHandlers) {
// 	// Movies routes
// 	movieRoutes := router.Group("/movies")
// 	{
// 		movieRoutes.GET("/", handlers.Media.MovieHandler.GetAll)
// 		movieRoutes.GET("/:id", handlers.Media.MovieHandler.GetByID)
// 	}
//
// 	// Series routes
// 	seriesRoutes := router.Group("/series")
// 	{
// 		seriesRoutes.GET("/", handlers.Media.SeriesHandler.GetAll)
// 		seriesRoutes.GET("/:id", handlers.Media.SeriesHandler.GetByID)
// 	}
//
// 	// Playlists routes
// 	playlistRoutes := router.Group("/playlists")
// 	{
// 		playlistRoutes.GET("/", handlers.Media.PlaylistHandler.GetAll)
// 		playlistRoutes.GET("/:id", handlers.Media.PlaylistHandler.GetByID)
// 	}
// }
//
// // ReplaceAppInitialization shows how to replace the current app initialization
// // with the wire-generated handlers
// func ReplaceAppInitialization() *app.AppDependencies {
// 	// Create a new context
// 	ctx := context.Background()
//
// 	// Initialize all handlers
// 	handlers, err := InitializeAllHandlers(ctx)
// 	if err != nil {
// 		log.Fatalf("Failed to initialize handlers: %v", err)
// 	}
//
// 	// In a real implementation, we would create the database connection and config service
// 	// For now, this is just a simplified example
// 	var db *gorm.DB
// 	var configService services.ConfigService
//
// 	// Create the app dependencies
// 	return &app.AppDependencies{
// 		// Set handlers
// 		// Note: In a real implementation, you would set these properly
// 	}
// }

