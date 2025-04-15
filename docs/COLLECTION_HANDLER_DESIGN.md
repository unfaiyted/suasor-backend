# Collection Handler Design Pattern

This document outlines the three-pronged architecture for collection handlers in the Suasor backend, following the same pattern established for series, music, and playlist handlers.

## Three-Pronged Architecture for Collections

The collection handlers implement a three-pronged architecture:

1. **Core Layer** (`core_collections.go`)
   - Basic database operations for collections
   - Provides data access without user context
   - Used by general public APIs

2. **User Layer** (`user_collections.go`)
   - User-specific operations for collections (creation, modification, etc.)
   - Extends core operations with user context
   - Used by personalized user APIs

3. **Client Layer** (`client_media_collection.go`)
   - Operations that interact with external collection libraries (Plex, Jellyfin, etc.)
   - Bridge between the application and media servers
   - Used by client-specific APIs

## Handler Organization

### 1. Core Collection Handler (`CoreCollectionHandler`)

- Responsible for basic database operations
- Provides read access to collections in the database
- Example endpoints:
  - `GET /collections` - Get all collections
  - `GET /collections/:id` - Get collection by ID
  - `GET /collections/:id/items` - Get items in a collection
  - `GET /collections/public` - Get public collections

### 2. User Collection Handler (`UserCollectionHandler`)

- Responsible for user-specific collection operations
- Manages collection creation, modification, and item management
- Example endpoints:
  - `GET /user/collections` - Get user's collections
  - `POST /user/collections` - Create a new collection
  - `PUT /user/collections/:id` - Update a collection
  - `DELETE /user/collections/:id` - Delete a collection
  - `POST /user/collections/:id/items` - Add an item to a collection

### 3. Client Collection Handler (`MediaClientCollectionHandler`)

- Responsible for retrieving collections from external media clients
- Manages access to Plex, Jellyfin, Emby collection libraries
- Example endpoints:
  - `GET /clients/media/{clientID}/collections` - Get collections from client
  - `GET /clients/media/{clientID}/collections/:id` - Get collection by ID from client
  - `GET /clients/media/{clientID}/collections/:id/items` - Get items in a client collection

## Implementation Details

### CoreCollectionHandler

- Uses `services.MediaItemService[*mediatypes.Collection]` and `services.CollectionService`
- Handles database queries and responses
- Focuses on read-only operations for public collection data

### UserCollectionHandler

- Uses `services.UserMediaItemService[*mediatypes.Collection]` and `services.CollectionService`
- Handles user-specific collection operations
- Manages collection creation, modification, and item management

### MediaClientCollectionHandler

- Uses `services.MediaClientCollectionService[T]`
- Handles communication with external collection libraries
- Converts client-specific collection data to common formats

## Router Configuration

The collection routes are organized in `router/collections.go`:

```go
// Core collection routes
collections := rg.Group("/collections")
{
  collections.GET("", coreCollectionHandler.GetAll)
  // ...
}

// User-specific collection routes
userCollections := rg.Group("/user/collections")
{
  userCollections.GET("", userCollectionHandler.GetUserCollections)
  // ...
}

// Client-specific routes are in router/media.go
// /clients/media/{clientID}/collections/...
```

## Usage Guidelines

- For basic collection data operations, use `CoreCollectionHandler`
- For user-specific collection operations, use `UserCollectionHandler`
- For client-specific collection operations, use `MediaClientCollectionHandler`
- Keep handler methods focused on their specific layer
- Ensure clear boundaries between layers

## Integration with App Dependencies

To fully implement this pattern, the app dependencies need to be updated:

```go
// In app/interfaces.go
type MediaItemServices interface {
  // Core services
  CoreCollectionService() services.MediaItemService[*mediatypes.Collection]
  
  // User services
  UserCollectionService() services.UserMediaItemService[*mediatypes.Collection]

  // Client services (accessed via MediaClientServices)
}

type CollectionServices interface {
  // Collection-specific operations beyond basic media item functions
  AddItemToCollection(ctx context.Context, collectionID, itemID uint64, itemType mediatypes.MediaType) error
  RemoveItemFromCollection(ctx context.Context, collectionID, itemID uint64) error
  GetCollectionItems(ctx context.Context, collectionID uint64) ([]mediatypes.MediaItem, error)
  GetCollectionsByUser(ctx context.Context, userID uint64) ([]*models.MediaItem[*mediatypes.Collection], error)
  // etc.
}
```

This three-pronged architecture provides clear separation of concerns, improves code organization, and ensures scalability as more collection features are added.