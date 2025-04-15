# User Media Item Data Service Design

This document outlines the service layer design for User Media Item Data, which follows the three-pronged architecture pattern.

## Overview

The User Media Item Data service layer is responsible for handling the business logic around user-specific data related to media items, such as watch history, favorites, ratings, and playback progress. This layer sits between the handlers (API controllers) and the repositories (data access).

## Three-Pronged Architecture

The service layer follows the same three-pronged architecture as the repository layer:

1. **Core Layer**: Handles basic operations that are common to all user media item data
2. **User Layer**: Manages user-specific operations like favorites, ratings, and watch history
3. **Client Layer**: Handles client-specific operations including synchronization with external media services

Each layer builds upon the previous one through composition, providing a clean separation of concerns while ensuring code reuse.

## Service Interfaces

### Core Layer - `CoreUserMediaItemDataService[T]`

The core service provides fundamental operations:

```go
type CoreUserMediaItemDataService[T types.MediaData] interface {
    Create(ctx context.Context, data *models.UserMediaItemData[T]) (*models.UserMediaItemData[T], error)
    GetByID(ctx context.Context, id uint64) (*models.UserMediaItemData[T], error)
    Update(ctx context.Context, data *models.UserMediaItemData[T]) (*models.UserMediaItemData[T], error)
    Delete(ctx context.Context, id uint64) error
    GetByUserIDAndMediaItemID(ctx context.Context, userID, mediaItemID uint64) (*models.UserMediaItemData[T], error)
    HasUserMediaItemData(ctx context.Context, userID, mediaItemID uint64) (bool, error)
}
```

### User Layer - `UserMediaItemDataService[T]`

The user service extends the core service with user-centric operations:

```go
type UserMediaItemDataService[T types.MediaData] interface {
    // Embed core service methods
    CoreUserMediaItemDataService[T]

    // User-specific methods
    GetUserHistory(ctx context.Context, userID uint64, limit, offset int, mediaType *types.MediaType) ([]*models.UserMediaItemData[T], error)
    GetRecentHistory(ctx context.Context, userID uint64, limit int, mediaType *types.MediaType) ([]*models.UserMediaItemData[T], error)
    GetUserPlayHistory(ctx context.Context, userID uint64, limit, offset int, mediaType *types.MediaType, completed *bool) ([]*models.UserMediaItemData[T], error)
    GetContinueWatching(ctx context.Context, userID uint64, limit int, mediaType *types.MediaType) ([]*models.UserMediaItemData[T], error)
    RecordPlay(ctx context.Context, data *models.UserMediaItemData[T]) (*models.UserMediaItemData[T], error)
    ToggleFavorite(ctx context.Context, mediaItemID, userID uint64, favorite bool) (*models.UserMediaItemData[T], error)
    UpdateRating(ctx context.Context, mediaItemID, userID uint64, rating float32) (*models.UserMediaItemData[T], error)
    GetFavorites(ctx context.Context, userID uint64, limit, offset int) ([]*models.UserMediaItemData[T], error)
    ClearUserHistory(ctx context.Context, userID uint64) error
}
```

### Client Layer - `ClientUserMediaItemDataService[T]`

The client service extends the user service with client-specific operations:

```go
type ClientUserMediaItemDataService[T types.MediaData] interface {
    // Embed user service methods
    UserMediaItemDataService[T]

    // Client-specific methods
    SyncClientItemData(ctx context.Context, userID uint64, clientID uint64, items []models.UserMediaItemData[T]) error
    GetClientItemData(ctx context.Context, userID uint64, clientID uint64, since *string) ([]*models.UserMediaItemData[T], error)
    GetByClientID(ctx context.Context, userID uint64, clientID uint64, clientItemID string) (*models.UserMediaItemData[T], error)
    RecordClientPlay(ctx context.Context, userID uint64, clientID uint64, clientItemID string, data *models.UserMediaItemData[T]) (*models.UserMediaItemData[T], error)
    GetPlaybackState(ctx context.Context, userID uint64, clientID uint64, clientItemID string) (*models.UserMediaItemData[T], error)
    UpdatePlaybackState(ctx context.Context, userID uint64, clientID uint64, clientItemID string, position int, duration int, percentage float64) (*models.UserMediaItemData[T], error)
}
```

## Service Implementation

Each service implementation follows a consistent pattern:

1. Takes repository dependencies via constructor injection
2. Delegates to the appropriate repository methods
3. Adds business logic, validation, and logging
4. Implements interface contract fully

### Example:

```go
// userMediaItemDataService implements UserMediaItemDataService
type userMediaItemDataService[T types.MediaData] struct {
    coreService CoreUserMediaItemDataService[T]
    repo        repository.UserMediaItemDataRepository[T]
}

// NewUserMediaItemDataService creates a new user media item data service
func NewUserMediaItemDataService[T types.MediaData](
    coreService CoreUserMediaItemDataService[T],
    repo repository.UserMediaItemDataRepository[T],
) UserMediaItemDataService[T] {
    return &userMediaItemDataService[T]{
        coreService: coreService,
        repo:        repo,
    }
}
```

## Service Factory

A factory pattern is provided to simplify the creation of services with proper dependencies:

```go
// UserMediaItemDataServiceFactory creates services for user media item data
type UserMediaItemDataServiceFactory struct {
    db *gorm.DB
}

// CreateClientService creates a client user media item data service
func (f *UserMediaItemDataServiceFactory) CreateClientService[T types.MediaData]() services.ClientUserMediaItemDataService[T] {
    userService := f.CreateUserService[T]()
    repo := repository.NewClientUserMediaItemDataRepository[T](f.db)
    return services.NewClientUserMediaItemDataService(userService, repo)
}

// Specialized factory methods
func (f *UserMediaItemDataServiceFactory) CreateMovieDataService() services.ClientUserMediaItemDataService[*types.Movie] {
    return f.CreateClientService[*types.Movie]()
}
```

## Type Safety with Generics

The service layer uses generics with `[T types.MediaData]` constraint to ensure type safety across different media types. This allows the services to work with any media type that satisfies the `MediaData` interface.

## Usage Example

```go
// Create the factory
factory := factory.NewUserMediaItemDataServiceFactory(db)

// Create type-specific services
movieService := factory.CreateMovieDataService()
seriesService := factory.CreateSeriesDataService()
musicService := factory.CreateMusicDataService()

// Use the services
favorites, err := movieService.GetFavorites(ctx, userID, 10, 0)
watchHistory, err := seriesService.GetUserHistory(ctx, userID, 20, 0, &types.MediaTypeSeries)
```

## Benefits

1. **Separation of Concerns**: Each layer has a clear responsibility
2. **Code Reuse**: Layers build on each other through composition
3. **Type Safety**: Generics ensure type safety across different media types
4. **Testability**: Each layer can be tested independently
5. **Extensibility**: New functionality can be added to the appropriate layer
6. **Consistency**: The three-pronged approach is consistent with other parts of the system
7. **Dependency Injection**: Services and repositories are composed using DI for better testability
