# Media Clients Architecture

## Overview

This document outlines the updated architecture for media clients in the Suasor backend, focusing on the implementation of Emby, Jellyfin, Plex, and Subsonic clients. The design emphasizes a factory pattern for media item conversion, improved type safety, and code organization.

## Core Components

### 1. Client Structure

Each media client implementation follows this common structure:

- **client.go**: Main client implementation that handles connection, authentication, and core functionality
- **factories.go**: Implements factory methods for converting provider-specific items to our standard models
- **helpers.go**: Utility functions for item conversion, query parameter mapping, and other client-specific operations
- **query.go** (for some clients): Shadow types for API query parameters

### 2. Factory Pattern Implementation

The system uses a generic factory pattern for media item conversion:

- **Registration**: Each client registers factory functions for converting their native item types to our standard models
- **Type Safety**: Uses generics to enforce type safety during conversion
- **Centralized Registry**: All factories are stored in a central registry for runtime lookup

### 3. Key Interfaces and Types

- **ClientMedia**: Base interface that all media clients must implement
- **ClientItemRegistry**: Registry for storing and retrieving factory functions
- **MediaFactory**: Generic type for factory functions with context-aware conversion

## Factory Pattern Details

### Registration Process

Factory functions are registered using the `RegisterFactory` function in a dedicated `RegisterMediaItemFactories` function:

```go
// Register all media item factories for a specific client
func RegisterMediaItemFactories(c *container.Container) {
    registry := container.MustGet[media.ClientItemRegistry](c)
    
    // Register factory for each media type
    media.RegisterFactory[*ClientType, *NativeItemType, *StandardModelType](
        &registry,
        func(client *ClientType, ctx context.Context, item *NativeItemType) (*StandardModelType, error) {
            return client.itemTypeFactory(ctx, item)
        },
    )
}
```

### Factory Implementation

Each client implements factory methods for every supported media type:

```go
// Example factory function for Movie
func (c *ClientType) movieFactory(ctx context.Context, item *NativeItemType) (*types.Movie, error) {
    // Convert native item to standard Movie model
    movie := &types.Movie{
        Details: types.MediaDetails{
            Title: item.Name,
            // ... other field mappings
        },
    }
    
    return movie, nil
}
```

### Item Conversion

Items are converted using the `ConvertTo` function, which looks up the appropriate factory:

```go
// Convert any item using the registered factory
movieItem, err := media.ConvertTo[*ClientType, *NativeItemType, *types.Movie](
    client, ctx, nativeItem)
```

## Helper Functions

The helpers.go file contains utility functions for:

1. **Simplified Conversion**: Higher-level functions for common conversion patterns
   ```go
   func GetItem[T types.MediaData](ctx, client, item) (T, error)
   func GetMediaItem[T types.MediaData](ctx, client, item, itemID) (*models.MediaItem[T], error)
   func GetMediaItemList[T types.MediaData](ctx, client, items) (*[]*models.MediaItem[T], error)
   ```

2. **Mixed Media Handling**: Functions for handling lists of different media types
   ```go
   func GetMixedMediaItems(ctx, client, items) (*models.MediaItemList, error)
   func GetMixedMediaItemsData(ctx, client, items) (*models.MediaItemDataList, error)
   ```

3. **Query Parameter Mapping**: Converting between our standard query options and client-specific parameters
   ```go
   func ApplyClientQueryOptions(queryParams, options) 
   ```

## Client-Specific Shadow Types

Some clients (like Jellyfin) implement shadow structs to work around limitations in the API libraries:

```go
// Shadow struct with direct field access
type ClientQueryOptions struct {
    UserId          *string
    SearchTerm      *string
    IncludeItemTypes *[]ItemType
    // ... other fields
}

// Methods to convert between our types and API types
func (o *ClientQueryOptions) ToItemsRequest() *api.GetItemsRequest {
    // ... convert shadow struct to API request
}

func (o *ClientQueryOptions) FromQueryOptions(options *types.QueryOptions) *ClientQueryOptions {
    // ... convert our standard options to shadow options
}
```

## Benefits of This Approach

1. **Type Safety**: The generic factory pattern provides compile-time type checking
2. **Separation of Concerns**: Clear separation between client functionality and conversion logic
3. **Consistent Patterns**: Standardized approach across all media clients
4. **Extensibility**: Easy to add support for new media types or clients
5. **Testability**: Conversion logic can be tested independently of API calls
6. **Maintainability**: Code organization makes it easier to understand and update
7. **API Independence**: Shadow types decouple our code from third-party API libraries

## Implementation for Each Client

Each client implementation follows this common pattern but adapts to the specific API requirements and data structures of its service:

### Emby Implementation

- Uses the Emby API client to interact with Emby servers
- Implements factory methods for all media types
- Provides helper functions for item conversion
- Maps query parameters to Emby API calls

### Jellyfin Implementation

- Uses the Jellyfin API client with similar patterns to Emby
- Implements shadow types to work around API library limitations
- Provides method-based query parameter mapping
- Follows the same factory pattern for item conversion

### Plex Implementation

- Adapts the pattern to Plex's unique API structure
- Handles Plex's specific media organization model
- Maps between Plex's item types and our standard models
- Implements factory methods for all supported media types

### Subsonic Implementation

- Implements the pattern for Subsonic's more limited media types
- Focuses on music-related media types (artists, albums, tracks)
- Adapts to Subsonic's simpler API structure
- Provides specialized handling for Subsonic's unique features

## Migration from Converters to Factories

The architecture has been migrated from using direct converter functions to factory methods:

- **Before**: Direct conversion functions in converters.go
  ```go
  func convertToMovie(item *nativeItem) *types.Movie { ... }
  ```

- **After**: Factory methods registered at initialization and used via the registry
  ```go
  func (c *Client) movieFactory(ctx context.Context, item *nativeItem) (*types.Movie, error) { ... }
  ```

This change improves code organization, type safety, and maintainability while enabling more consistent error handling and logging.

## Common Migration Steps

When migrating a client from the old approach to the new factory pattern:

1. Create a factories.go file with factory methods for each media type
2. Register these factories using the RegisterMediaItemFactories function
3. Update the client initialization to create and store the registry
4. Implement helper functions for common conversion patterns
5. Refactor existing converters to use the factory pattern
6. Update the client implementation to use the new conversion approach

The result is a more maintainable, type-safe implementation that follows a consistent pattern across all media clients.