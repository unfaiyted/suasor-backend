# User Media Item Data Handlers Design

This document outlines the design of the handlers for User Media Item Data, implementing the three-pronged architecture pattern at the API layer.

## Overview

The User Media Item Data handlers are responsible for exposing API endpoints that allow clients to interact with user-specific data related to media items, such as watch history, favorites, ratings, and playback progress. Following the three-pronged architecture approach used in the repository and service layers, the handlers are organized into three distinct layers:

1. **Core Layer**: Basic CRUD operations and core functionality
2. **User Layer**: User-specific operations like history, favorites, and ratings
3. **Client Layer**: Client-specific operations for synchronization with external media systems

## Handler Structure

### Core Layer - `CoreUserMediaItemDataHandler[T]`

The core handler provides fundamental API endpoints for basic CRUD operations:

- `GET /user-media-data/:id` - Get a specific user media item data entry
- `GET /user-media-data/check` - Check if a user has data for a specific media item
- `GET /user-media-data/user-media` - Get user media item data for a specific user and media item
- `DELETE /user-media-data/:id` - Delete a specific user media item data entry

### User Layer - `UserUserMediaItemDataHandler[T]`

The user handler extends the core handler with user-centric endpoints:

- `GET /user-media-data/history` - Get a user's media play history
- `GET /user-media-data/continue-watching` - Get media items that a user has started but not completed
- `GET /user-media-data/recent` - Get a user's recent media history
- `POST /user-media-data/record` - Record a new play event
- `PUT /user-media-data/media/:mediaItemId/favorite` - Toggle favorite status for a media item
- `PUT /user-media-data/media/:mediaItemId/rating` - Update user rating for a media item
- `GET /user-media-data/favorites` - Get a user's favorite media items
- `DELETE /user-media-data/clear` - Clear a user's play history

The user layer also includes forwarding methods for all core layer endpoints, ensuring a complete API surface.

### Client Layer - `ClientUserMediaItemDataHandler[T]`

The client handler extends the user handler with client-specific endpoints:

- `POST /user-media-data/client/:clientId/sync` - Synchronize user media item data from an external client
- `GET /user-media-data/client/:clientId` - Get user media item data for synchronization with a client
- `GET /user-media-data/client/:clientId/item/:clientItemId` - Get user media item data by client ID
- `POST /user-media-data/client/:clientId/item/:clientItemId/play` - Record a play event from a client
- `GET /user-media-data/client/:clientId/item/:clientItemId/state` - Get playback state for a client item
- `PUT /user-media-data/client/:clientId/item/:clientItemId/state` - Update playback state for a client item

The client layer includes forwarding methods for all user and core layer endpoints, ensuring a complete API surface.

## Media Type Specialization

The handlers use Go's generics to provide type safety for different media types:

```go
// Generic handler type
type ClientUserMediaItemDataHandler[T types.MediaData] struct {
    service services.ClientUserMediaItemDataService[T]
    userHandler *UserUserMediaItemDataHandler[T]
}
```

Media-type specific handlers can be created for movies, series, episodes, tracks, etc.:

```go
// Movie-specific handler
movieHandler := NewClientUserMediaItemDataHandler[*types.Movie](
    service, userHandler)

// Series-specific handler
seriesHandler := NewClientUserMediaItemDataHandler[*types.Series](
    service, userHandler)
```

## Router Implementation

The router layer dynamically selects the appropriate handler based on the media type parameter:

```go
handlerMap := map[string]UserMediaItemDataHandlerInterface{
    "movies": mediaHandlers.MovieHandler(),
    "series": mediaHandlers.SeriesHandler(),
    "tracks": mediaHandlers.TrackHandler(),
    // ...
}

getHandler := func(c *gin.Context) UserMediaItemDataHandlerInterface {
    mediaType := c.Param("mediaType")
    handler, exists := handlerMap[mediaType]
    if !exists {
        // Default to movie handler if type not specified or invalid
        return mediaHandlers.MovieHandler()
    }
    return handler
}
```

Routes are organized in a hierarchical structure that mirrors the three-pronged architecture:

```
/user-media-data                     # Base path
  /:id                               # Core endpoints
  /check
  /user-media

  /history                           # User endpoints
  /continue-watching
  /favorites
  ...

  /client/:clientId                  # Client endpoints
    /sync
    /item/:clientItemId
    /item/:clientItemId/play
    ...

  /movies                            # Media-type specific endpoints
    /history
    /favorites
    ...
```

## Dependency Injection

The handlers follow a dependency injection pattern, where each layer receives its dependencies via constructor injection:

```go
// Core Layer
func NewCoreUserMediaItemDataHandler[T types.MediaData](
    service services.CoreUserMediaItemDataService[T],
) *CoreUserMediaItemDataHandler[T]

// User Layer
func NewUserUserMediaItemDataHandler[T types.MediaData](
    service services.UserMediaItemDataService[T],
    coreHandler *CoreUserMediaItemDataHandler[T],
) *UserUserMediaItemDataHandler[T]

// Client Layer
func NewClientUserMediaItemDataHandler[T types.MediaData](
    service services.ClientUserMediaItemDataService[T],
    userHandler *UserUserMediaItemDataHandler[T],
) *ClientUserMediaItemDataHandler[T]
```

This approach:
1. Makes testing easier by allowing mock dependencies
2. Ensures clear separation of concerns
3. Makes dependencies explicit
4. Facilitates composability

## Benefits of the Three-Pronged Approach at the Handler Level

1. **API Versioning**: The three-pronged approach makes it easier to evolve APIs over time, as each layer can be versioned independently.

2. **Consistent Authorization**: Different authorization rules can be applied to different layers (e.g., client operations might require additional privileges).

3. **Documentation Organization**: API documentation is naturally organized into logical groups that match the business domains.

4. **Simplified Testing**: Each handler layer can be tested independently, with mock services injected for isolation.

5. **Code Organization**: The separation makes it easier to locate relevant code when working on a specific feature.

6. **Unified Interface**: The client layer provides a complete API surface through method forwarding, so clients only need to work with a single handler type.

## Usage Example

```go
// Create services
coreService := services.NewCoreUserMediaItemDataService[*types.Movie](repo)
userService := services.NewUserMediaItemDataService[*types.Movie](coreService, repo)
clientService := services.NewClientUserMediaItemDataService[*types.Movie](userService, repo)

// Create handlers
coreHandler := handlers.NewCoreUserMediaItemDataHandler[*types.Movie](coreService)
userHandler := handlers.NewUserUserMediaItemDataHandler[*types.Movie](userService, coreHandler)
clientHandler := handlers.NewClientUserMediaItemDataHandler[*types.Movie](clientService, userHandler)

// Register routes
router.GET("/user-media-data/history", userHandler.GetMediaPlayHistory)
router.POST("/user-media-data/client/:clientId/sync", clientHandler.SyncClientItemData)
```

## Conclusion

The three-pronged architecture for User Media Item Data handlers provides a clean, maintainable, and extensible approach to building APIs. By separating concerns into core, user, and client layers, the architecture facilitates independent evolution of each layer while ensuring a consistent and comprehensive API surface.
