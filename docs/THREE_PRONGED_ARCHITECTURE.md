# Three-Pronged Architecture in Suasor

This document outlines the clean Three-Pronged Architecture pattern implemented in the Suasor backend, focusing on how it's integrated into the dependency injection system.

## Overview

The Three-Pronged Architecture consists of three distinct layers:

1. **Core Layer**: Basic CRUD operations on media items without client or user associations
2. **User Layer**: User-specific operations (extends Core)
3. **Client Layer**: Client-specific operations (extends User)

This approach allows for clear separation of concerns while maintaining a single database table structure.

## Design Principles

- **Clean Interfaces**: Each layer has a clearly defined interface
- **Hierarchical Structure**: Higher layers build on lower layers
- **Type Safety**: Go generics ensure type correctness across layers
- **Factory Pattern**: Factory methods create properly configured components
- **No Legacy Code**: Clean design without backward compatibility compromises
- **Domain-Driven Design**: Components are organized by domain (Movies, Series, etc.)

## Implementation Details

### Repository Layer

The repository layer is structured with three types of repositories:

```
MediaItemRepository[T] (Core)
  ↑
UserMediaItemRepository[T] (User)
  ↑
ClientMediaItemRepository[T] (Client)
```

Each repository focuses on its specific concerns:
- **Core**: Basic CRUD operations
- **User**: User-owned and user-specific operations
- **Client**: Operations for media items linked to external clients

### Service Layer

The service layer follows the same pattern:

```
CoreMediaItemService[T]
  ↑
UserMediaItemService[T]
  ↑
ClientMediaItemService[T]
```

Each service encapsulates its layer's behavior and delegates to the lower layer when needed.

### Handler Layer

The handler layer completes the pattern:

```
CoreMediaItemHandler[T]
  ↑
UserMediaItemHandler[T]
  ↑
ClientMediaItemHandler[T]
```

Each handler exposes endpoints appropriate to its layer and delegates to the lower layer when needed.

## Dependency Injection Integration

The dependency injection system uses a factory-based approach to create and wire together the three-pronged components:

### The ThreeProngedFactory

The `ThreeProngedFactory` is the centerpiece of our architecture, providing methods to create and configure all components:

```go
// Create the factory
factory := NewThreeProngedFactory(db, clientFactory)

// Create repositories for all media types
coreRepos := factory.CreateCoreRepositories()
userRepos := factory.CreateUserRepositories()
clientRepos := factory.CreateClientRepositories()

// Create services for all media types
coreServices := factory.CreateCoreServices(coreRepos)
userServices := factory.CreateUserServices(coreServices, userRepos)
clientServices := factory.CreateClientServices(coreServices, clientRepos)

// Create handlers for all media types
coreHandlers := factory.CreateCoreHandlers(coreServices)
userHandlers := factory.CreateUserHandlers(userServices, coreHandlers)
clientHandlers := factory.CreateClientHandlers(clientServices, userHandlers)

// Create specialized handlers for specific domains
specializedHandlers := factory.CreateSpecializedMediaHandlers(
    coreServices, userServices, clientServices,
    musicHandler, seriesHandler, playlistHandler, collectionHandler)

// Create specialized services for collections and playlists
collectionServices := factory.CreateMediaCollectionServices(
    coreServices, userServices, clientServices,
    coreCollectionService, userCollectionService, clientCollectionService,
    playlistService)
```

### Accessing Components Through AppDependencies

All components are accessible through the `AppDependencies` struct:

```go
// Core layer
deps.CoreRepositories.MovieRepo()
deps.CoreMediaItemServices.MovieCoreService()
deps.CoreMediaItemHandlers.MovieCoreHandler()

// User layer
deps.UserRepositoryFactories.MovieUserRepo()
deps.UserMediaItemServices.MovieUserService()
deps.UserMediaItemHandlers.MovieUserHandler()

// Client layer
deps.ClientRepositoryFactories.MovieClientRepo()
deps.ClientMediaItemServices.MovieClientService()
deps.ClientMediaItemHandlers.MovieClientHandler()
```

## Media Type Specialization

The architecture is specialized for each media type:

- Movies
- Series
- Episodes
- Tracks
- Albums
- Artists
- Collections
- Playlists

Each type gets its own set of repositories, services, and handlers at all three layers.

## Specialized Services and Handlers

In addition to the three-pronged architecture, we have specialized components for specific domains:

```go
// Specialized media handlers
deps.SpecializedMediaHandlers.MusicHandler()
deps.SpecializedMediaHandlers.SeriesSpecificHandler()
deps.SpecializedMediaHandlers.PlaylistSpecificHandler()
deps.SpecializedMediaHandlers.CollectionSpecificHandler()

// Specialized collection services
deps.MediaCollectionServices.CoreCollectionService()
deps.MediaCollectionServices.UserCollectionService()
deps.MediaCollectionServices.ClientCollectionService()
deps.MediaCollectionServices.PlaylistService()
```

## Benefits

1. **Clear Separation of Concerns**: Each layer has a specific focus
2. **Code Reuse**: Higher layers extend lower layers, reducing duplication
3. **Flexibility**: Changes to one layer don't necessarily affect others
4. **Type Safety**: Generic parameters ensure type correctness
5. **Factory Pattern**: Easy creation of properly configured components
6. **Clean Design**: No legacy code or backward compatibility compromises
7. **Testability**: Each layer can be tested independently

## Usage Examples

### Accessing Core Operations

```go
// Get a movie by ID (core operation)
deps.CoreMediaItemServices.MovieCoreService().GetByID(ctx, movieID)
```

### Accessing User Operations

```go
// Get user-owned playlists (user operation)
deps.UserMediaItemServices.PlaylistUserService().GetUserContent(ctx, userID, 10)
```

### Accessing Client Operations

```go
// Sync client-specific metadata (client operation)
deps.ClientMediaItemServices.MovieClientService().SyncClientItemData(ctx, itemID, clientID, clientItemData)
```

### Using the Three-Pronged Handlers in Routers

```go
// Register routes for all three layers
func RegisterMovieRoutes(router *gin.RouterGroup, deps *app.AppDependencies) {
    // Core routes (accessible to all)
    router.GET("/movies/:id", deps.CoreMediaItemHandlers.MovieCoreHandler().GetByID)
    
    // User routes (require user authentication)
    authRouter := router.Group("/user")
    authRouter.Use(middleware.RequireAuth())
    authRouter.GET("/movies/favorites", deps.UserMediaItemHandlers.MovieUserHandler().GetFavorites)
    
    // Client routes (require client authentication)
    clientRouter := router.Group("/client/:clientID")
    clientRouter.Use(middleware.RequireClientAuth())
    clientRouter.POST("/movies/sync", deps.ClientMediaItemHandlers.MovieClientHandler().SyncClientItemData)
}
```