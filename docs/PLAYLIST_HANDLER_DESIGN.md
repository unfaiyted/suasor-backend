# Playlist Handler Design Pattern

This document outlines the three-pronged architecture for playlist handlers in the Suasor backend, following the same pattern established for series and music handlers.

## Three-Pronged Architecture for Playlists

The playlist handlers implement a three-pronged architecture:

1. **Core Layer** (`core_playlists.go`)
   - Basic database operations for playlists
   - Provides data access without user context
   - Used by general public APIs

2. **User Layer** (`user_playlists.go`)
   - User-specific operations for playlists (creation, modification, etc.)
   - Extends core operations with user context
   - Used by personalized user APIs

3. **Client Layer** (`client_media_playlist.go`)
   - Operations that interact with external playlist libraries (Plex, Jellyfin, etc.)
   - Bridge between the application and media servers
   - Used by client-specific APIs

## Handler Organization

### 1. Core Playlist Handler (`CorePlaylistHandler`)

- Responsible for basic database operations
- Provides read access to playlists in the database
- Example endpoints:
  - `GET /playlists` - Get all playlists
  - `GET /playlists/:id` - Get playlist by ID
  - `GET /playlists/:id/tracks` - Get tracks in a playlist

### 2. User Playlist Handler (`UserPlaylistHandler`)

- Responsible for user-specific playlist operations
- Manages playlist creation, modification, and track management
- Example endpoints:
  - `GET /user/playlists` - Get user's playlists
  - `POST /user/playlists` - Create a new playlist
  - `PUT /user/playlists/:id` - Update a playlist
  - `DELETE /user/playlists/:id` - Delete a playlist
  - `POST /user/playlists/:id/tracks` - Add a track to a playlist

### 3. Client Playlist Handler (`ClientMediaPlaylistHandler`)

- Responsible for retrieving playlists from external media clients
- Manages access to Plex, Jellyfin, Emby playlist libraries
- Example endpoints:
  - `GET /clients/media/{clientID}/playlists` - Get playlists from client
  - `GET /clients/media/{clientID}/playlists/:id` - Get playlist by ID from client
  - `GET /clients/media/{clientID}/playlists/:id/tracks` - Get tracks in a client playlist

## Implementation Details

### CorePlaylistHandler

- Uses `services.MediaItemService[*mediatypes.Playlist]` and `services.PlaylistService`
- Handles database queries and responses
- Focuses on read-only operations for public playlist data

### UserPlaylistHandler

- Uses `services.UserMediaItemService[*mediatypes.Playlist]` and `services.PlaylistService`
- Handles user-specific playlist operations
- Manages playlist creation, modification, and track management

### ClientMediaPlaylistHandler

- Uses `services.ClientMediaPlaylistService[T]`
- Handles communication with external playlist libraries
- Converts client-specific playlist data to common formats

## Router Configuration

The playlist routes are organized in `router/playlists.go`:

```go
// Core playlist routes
playlists := rg.Group("/playlists")
{
  playlists.GET("", corePlaylistHandler.GetAll)
  // ...
}

// User-specific playlist routes
userPlaylists := rg.Group("/user/playlists")
{
  userPlaylists.GET("", userPlaylistHandler.GetUserPlaylists)
  // ...
}

// Client-specific routes are in router/media.go
// /clients/media/{clientID}/playlists/...
```

## Usage Guidelines

- For basic playlist data operations, use `CorePlaylistHandler`
- For user-specific playlist operations, use `UserPlaylistHandler`
- For client-specific playlist operations, use `ClientMediaPlaylistHandler`
- Keep handler methods focused on their specific layer
- Ensure clear boundaries between layers

## Integration with App Dependencies

To fully implement this pattern, the app dependencies need to be updated:

```go
// In app/interfaces.go
type MediaItemServices interface {
  // Core services
  CorePlaylistService() services.MediaItemService[*mediatypes.Playlist]
  
  // User services
  UserPlaylistService() services.UserMediaItemService[*mediatypes.Playlist]

  // Client services (accessed via ClientMediaServices)
}

type PlaylistServices interface {
  // Playlist-specific operations beyond basic media item functions
  AddTrackToPlaylist(ctx context.Context, playlistID, trackID uint64) error
  RemoveTrackFromPlaylist(ctx context.Context, playlistID, trackID uint64) error
  GetPlaylistsByUser(ctx context.Context, userID uint64) ([]*models.MediaItem[*mediatypes.Playlist], error)
  // etc.
}
```

This three-pronged architecture provides clear separation of concerns, improves code organization, and ensures scalability as more playlist features are added.
