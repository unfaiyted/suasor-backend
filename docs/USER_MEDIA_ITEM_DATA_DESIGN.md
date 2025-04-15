# User Media Item Data Repository Design

This document outlines the design of the User Media Item Data repository, which follows a three-pronged architecture for handling user-specific data related to media items.

## Overview

The User Media Item Data repository is responsible for managing user-specific data related to media items, such as:

- Watch history
- Favorites
- Ratings
- Playback progress
- Play counts
- Completed status

The repository is designed using a three-pronged architecture to separate concerns and provide a clean interface for each layer:

1. **Core Layer**: Handles basic CRUD operations and database access
2. **User Layer**: Manages user-specific operations like favorites, ratings, and history
3. **Client Layer**: Handles client-specific operations including synchronization with external media systems

## Data Model

The primary data model is `UserMediaItemData[T]`, a generic type where `T` satisfies the `types.MediaData` interface. This allows type-safe operations for different media types (movies, series, music, etc.).

```go
type UserMediaItemData[T types.MediaData] struct {
    ID               uint64          `json:"id" gorm:"primaryKey;autoIncrement"`
    UserID           uint64          `json:"userId" gorm:"index"`
    MediaItemID      uint64          `json:"mediaItemId" gorm:"index"`
    Item             *MediaItem[T]   `json:"item" gorm:"-"`
    Type             types.MediaType `json:"type" gorm:"type:varchar(50)"`
    PlayedAt         time.Time       `json:"playedAt" gorm:"index"`
    LastPlayedAt     time.Time       `json:"lastPlayedAt" gorm:"index"`
    IsFavorite       bool            `json:"isFavorite,omitempty"`
    IsDisliked       bool            `json:"isDisliked,omitempty"`
    UserRating       float32         `json:"userRating,omitempty"`
    Watchlist        bool            `json:"watchlist,omitempty"`
    PlayedPercentage float64         `json:"playedPercentage,omitempty"`
    PlayCount        int32           `json:"playCount,omitempty"`
    PositionSeconds  int             `json:"positionSeconds"`
    DurationSeconds  int             `json:"durationSeconds"`
    Completed        bool            `json:"completed"`
    CreatedAt        time.Time       `json:"createdAt" gorm:"autoCreateTime"`
    UpdatedAt        time.Time       `json:"updatedAt" gorm:"autoUpdateTime"`
}
```

## Repository Structure

### Core Layer - `CoreUserMediaItemDataRepository[T]`

The core layer provides basic CRUD operations for the `UserMediaItemData[T]` model:

- `Create`: Creates a new user media item data entry
- `GetByID`: Retrieves a specific entry by ID
- `Update`: Updates an existing entry
- `Delete`: Removes a specific entry
- `GetByUserIDAndMediaItemID`: Retrieves data for a specific user and media item
- `HasUserMediaItemData`: Checks if a user has data for a specific media item

### User Layer - `UserMediaItemDataRepository[T]`

The user layer provides user-focused operations:

- `GetUserHistory`: Retrieves a user's media history
- `GetRecentHistory`: Retrieves a user's recent media history
- `GetUserPlayHistory`: Retrieves play history with optional filtering
- `GetContinueWatching`: Retrieves items that a user has started but not completed
- `RecordPlay`: Records a new play event
- `ToggleFavorite`: Marks or unmarks a media item as a favorite
- `UpdateRating`: Sets a user's rating for a media item
- `GetFavorites`: Retrieves favorite media items for a user
- `ClearUserHistory`: Removes all data for a user

### Client Layer - `ClientUserMediaItemDataRepository[T]`

The client layer handles client-specific operations for integrating with external media services:

- `SyncClientItemData`: Synchronizes user media item data from an external client
- `GetClientItemData`: Retrieves data for synchronization with a client
- `GetByClientID`: Retrieves a user media item data entry by client ID
- `RecordClientPlay`: Records a play event from a client
- `MapClientMediaItemToInternal`: Maps a client media item to an internal media item
- `GetPlaybackState`: Retrieves the current playback state for a client item
- `UpdatePlaybackState`: Updates the playback state for a client item

### Facade - `UserMediaItemDataRepository[T]`

The main repository interface (`UserMediaItemDataRepository[T]`) acts as a facade, combining the functionality of all three layers into a single interface. This provides a unified API for consumers while maintaining internal separation of concerns.

## Design Principles

1. **Type Safety**: Uses generics with the `[T types.MediaData]` constraint to ensure type safety across media types
2. **Separation of Concerns**: Each layer focuses on a specific aspect of user media data
3. **Single Responsibility**: Each repository and method has a clearly defined responsibility
4. **Dependency Injection**: Components are composed using dependency injection for easier testing and flexibility
5. **Error Handling**: Consistent error handling and wrapping for better debugging

## Usage Example

```go
// Create the repository with the correct media type
movieRepo := repository.NewUserMediaItemDataRepository[*types.Movie](db)
seriesRepo := repository.NewUserMediaItemDataRepository[*types.Series](db)
musicRepo := repository.NewUserMediaItemDataRepository[*types.Track](db)

// User operations
favorites, err := movieRepo.GetFavorites(ctx, userID, 10, 0)
history, err := seriesRepo.GetUserHistory(ctx, userID, 20, 0, &types.MediaTypeSeries)
musicHistory, err := musicRepo.GetRecentHistory(ctx, userID, 10, &types.MediaTypeTrack)

// Client operations
err := movieRepo.SyncClientItemData(ctx, userID, "plex", clientItems)
state, err := seriesRepo.GetPlaybackState(ctx, userID, "emby", "series123")
```

## Benefits of Three-Pronged Architecture

1. **Maintainability**: Each layer can be modified independently without affecting the others
2. **Testability**: Smaller, focused components are easier to test
3. **Flexibility**: New functionality can be added to the appropriate layer without disrupting others
4. **Scalability**: Can adapt to different database technologies or client types by modifying specific layers
5. **Type Safety**: Consistent use of generics ensures type safety throughout the repository
