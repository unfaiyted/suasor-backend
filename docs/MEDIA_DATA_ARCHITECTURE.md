# Three-Pronged Media Data Architecture

## Overview

The Three-Pronged Architecture is a layered approach to handling media data in Suasor. It consists of three primary layers:

1. **Core Layer** - Basic CRUD operations on media items
2. **User Layer** - User-specific operations (extends Core)
3. **Client Layer** - Client-specific operations (extends User)

This architecture follows the principle of "inheritance by composition" where each higher layer extends the capabilities of the layers below it.

## Key Components

### Repositories

- **CoreMediaItemRepository**: Basic CRUD operations for media items
- **UserMediaItemRepository**: User-specific operations (e.g., favorites, watched status)
- **ClientMediaItemRepository**: Client-specific operations (e.g., specific to Jellyfin, Plex)

### Services

- **CoreMediaItemService**: Basic CRUD service operations
- **UserMediaItemService**: User-specific service operations
- **ClientMediaItemService**: Client-specific service operations

### Handlers

- **CoreUserMediaItemDataHandler**: Core layer endpoint handlers
- **UserUserMediaItemDataHandler**: User layer endpoint handlers
- **ClientUserMediaItemDataHandler**: Client layer endpoint handlers

## Factory Pattern

The architecture uses the `MediaDataFactory` to create properly configured components. This factory:

- Creates repositories, services, and handlers for each layer
- Ensures proper initialization and composition
- Maintains type safety through generics

## Dependency Injection

The architecture leverages structured dependency injection through:

- **ServiceRegistrar**: Handles registration of all components
- **Initialize**: Creates and configures all dependencies
- **AppDependencies**: Contains all application dependencies

## Media Types

The architecture supports various media types through generics:

- Movie
- Series
- Episode
- Track
- Album
- Artist
- Collection
- Playlist

## Implementation Details

### Creating a New Component

To implement a new media type:

1. Add the media type to the CoreRepositories, UserRepositoryFactories, and ClientRepositoryFactories interfaces
2. Add the media type to the CoreMediaItemServices, UserMediaItemServices, and ClientMediaItemServices interfaces
3. Add the media type to the CoreMediaItemHandlers, UserMediaItemHandlers, and ClientMediaItemHandlers interfaces
4. Update the implementations to include the new components

### Inheritance by Composition Example

The core idea is that each layer extends the functionality of the previous layer by composition:

```go
// User layer extending Core layer
type UserMediaItemService[T mediatypes.MediaData] struct {
    coreService CoreMediaItemService[T]
    repository  UserMediaItemRepository[T]
}

// Client layer extending User layer
type ClientMediaItemService[T mediatypes.MediaData] struct {
    userService UserMediaItemService[T]
    repository  ClientMediaItemRepository[T]
}
```

### Specialized Domain Handlers

For specialized functionality beyond the basic CRUD operations:

- **CoreMovieHandler**: Movie-specific operations
- **CoreCollectionHandler**: Collection-specific operations
- **CorePlaylistHandler**: Playlist-specific operations
- **CoreMusicHandler**: Music-specific operations

## Benefits

1. **Type Safety**: Uses Go generics for type-safe components
2. **Separation of Concerns**: Clear separation between core, user, and client responsibilities
3. **Maintainability**: Standardized approach makes code more maintainable
4. **Extensibility**: Easy to add new media types or functionality
5. **Reduced Code Duplication**: Shared functionality in base layers
6. **Testability**: Clear interfaces make components more testable

## Usage Examples

### Query from Core Layer

```go
// Get all movies
movies, err := coreMovieService.GetAll(ctx, limit, offset)
```

### Query from User Layer

```go
// Get user's favorite movies
favorites, err := userMovieService.GetUserFavorites(ctx, userId, limit, offset)
```

### Query from Client Layer

```go
// Get movies from a specific Jellyfin server
jellyfinMovies, err := clientMovieService.GetFromClient(ctx, clientId, limit, offset)
```