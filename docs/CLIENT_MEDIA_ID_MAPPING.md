# Client Media ID Mapping

This document describes how client-specific media item IDs are stored and retrieved in the system.

## Overview

Our system integrates with multiple external media clients (Plex, Emby, Jellyfin, etc.). Each client has its own ID system for media items. To provide a seamless user experience, we maintain mappings between our internal media items and the client-specific IDs.

## Storage Structure

Client IDs are stored in the `media_items` table as a JSONB field called `sync_clients`. This field contains an array of `SyncClient` objects:

```go
type SyncClient struct {
    // ID of the client that this external ID belongs to
    ID uint64 `json:"clientId,omitempty"`
    // Type of client this ID belongs to
    Type client.ClientType `json:"clientType,omitempty" gorm:"type:varchar(50)"`
    // The actual ID value in the external system
    ItemID string `json:"itemId"`
}
```

This structure allows each media item to have multiple client IDs, one for each client system it exists in.

## Key Differences: SyncClients vs ExternalIDs

It's important to understand the distinction between `SyncClients` and `ExternalIDs`:

1. **SyncClients** - Used for mapping to our integrated media clients (Plex, Emby, Jellyfin, etc.)
   - These are systems we actively sync with
   - Client-specific and tied to a specific client instance
   - Example: A movie's ID in a user's Plex server

2. **ExternalIDs** - Used for mapping to external metadata providers (TMDB, IMDB, etc.)
   - These are external reference IDs not tied to any specific client
   - Global identifiers that don't change between systems
   - Example: IMDB ID "tt0111161" for "The Shawshank Redemption"

## Client Media Item Helper

To simplify working with client media items, we've introduced a `ClientMediaItemHelper` that provides methods for mapping between client IDs and internal IDs:

```go
// ClientMediaItemHelper provides helper methods for working with client media items
type ClientMediaItemHelper struct {
    db *gorm.DB
}
```

### Key Methods

1. **GetMediaItemByClientID** - Finds a media item by client ID
   ```go
   func (h *ClientMediaItemHelper) GetMediaItemByClientID(ctx context.Context, clientID uint64, clientItemID string) (uint64, error)
   ```

2. **GetMediaItemByClientIDAndType** - Finds a media item by client ID and media type
   ```go
   func (h *ClientMediaItemHelper) GetMediaItemByClientIDAndType(ctx context.Context, clientID uint64, clientItemID string, mediaType types.MediaType) (uint64, error)
   ```

3. **GetMediaItemsByClientIDs** - Batch retrieval of multiple media items by client IDs
   ```go
   func (h *ClientMediaItemHelper) GetMediaItemsByClientIDs(ctx context.Context, clientID uint64, clientItemIDs []string) (map[string]uint64, error)
   ```

4. **GetOrCreateMediaItemMapping** - Ensures a mapping exists, creating one if needed
   ```go
   func (h *ClientMediaItemHelper) GetOrCreateMediaItemMapping[T types.MediaData](ctx context.Context, clientID uint64, clientType types.ClientType, clientItemID string, mediaType types.MediaType, title string, data T) (uint64, error)
   ```

5. **SyncClientsList** - Gets all media items with client IDs for a specific client
   ```go
   func (h *ClientMediaItemHelper) SyncClientsList(ctx context.Context, clientID uint64) (map[string]uint64, error)
   ```

## JSONB Queries

The helper uses PostgreSQL JSONB queries to efficiently search within the SyncClients array:

```sql
SELECT id FROM media_items 
WHERE sync_clients @> '[{"id": 123, "itemId": "client_item_abc"}]'::jsonb
```

This query finds all media items where the `sync_clients` array contains an object with the specified client ID and item ID.

## Usage Examples

### Finding a Media Item by Client ID

```go
// Create the helper
helper := NewClientMediaItemHelper(db)

// Get internal media item ID from client ID
internalID, err := helper.GetMediaItemByClientID(ctx, clientID, clientItemID)
if err != nil {
    // Handle not found or error
}

// Use the internal ID
mediaItem, err := repo.GetMediaItemByID(ctx, internalID)
```

### Syncing Client Data

```go
// Get all client IDs for synchronized media items
clientMapping, err := helper.SyncClientsList(ctx, clientID)

// clientMapping is a map from client item IDs to internal IDs
// { "client_item_123": 456, "client_item_456": 789 }

// Use this mapping for efficient syncing
for clientItemID, internalID := range clientMapping {
    // Process each mapped item
}
```

## Best Practices

1. **Always use the helper** - The helper abstracts away the complexity of JSONB queries and provides a consistent interface

2. **Consider batch operations** - When processing multiple items, use batch methods to reduce database queries

3. **Handle not found gracefully** - Create new mappings when appropriate instead of returning errors

4. **Use transactions** - When creating or updating multiple mappings, wrap operations in a transaction

5. **Add proper indexes** - Ensure your database has GIN indexes on the JSONB fields for performance:
   ```sql
   CREATE INDEX idx_media_items_sync_clients ON media_items USING GIN (sync_clients);
   ```