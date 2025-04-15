# Music Handler Design Pattern

This document outlines the three-pronged architecture for music handlers in the Suasor backend, following the same pattern established for series handlers.

## Three-Pronged Architecture for Music

The music handlers implement a three-pronged architecture:

1. **Core Layer** (`core_music.go`)
   - Basic database operations for tracks, albums, and artists
   - Provides data access without user context
   - Used by general public APIs

2. **User Layer** (`user_music.go`)
   - User-specific operations for music (favorites, ratings, play history, etc.)
   - Extends core operations with user context
   - Used by personalized user APIs

3. **Client Layer** (`client_media_music.go`)
   - Operations that interact with external music libraries (Plex, Jellyfin, etc.)
   - Bridge between the application and media servers
   - Used by client-specific APIs

## Handler Organization

### 1. Core Music Handler (`CoreMusicHandler`)

- Responsible for basic database operations
- Manages tracks, albums, and artists in the core database
- Example endpoints:
  - `GET /music/tracks/top` - Get top tracks
  - `GET /music/albums/:id/tracks` - Get tracks for an album
  - `GET /music/artists/:id/albums` - Get albums for an artist

### 2. User Music Handler (`UserMusicHandler`)

- Responsible for user-specific music operations
- Manages user favorites, play history, and ratings
- Example endpoints:
  - `GET /user/music/tracks/favorites` - Get user favorite tracks
  - `GET /user/music/albums/favorites` - Get user favorite albums
  - `GET /user/music/tracks/recently-played` - Get recently played tracks

### 3. Client Music Handler (`ClientMediaMusicHandler`)

- Responsible for retrieving music from external media clients
- Manages access to Plex, Jellyfin, Emby music libraries
- Example endpoints:
  - `GET /clients/media/{clientID}/music/tracks` - Get tracks from client
  - `GET /clients/media/{clientID}/music/albums` - Get albums from client
  - `GET /clients/media/{clientID}/music/artists` - Get artists from client

## Implementation Details

### CoreMusicHandler

- Uses `services.MediaItemService[*mediatypes.Track/Album/Artist]`
- Handles database queries and responses
- Focuses on read-only operations for public music data

### UserMusicHandler

- Uses `services.UserMediaItemService[*mediatypes.Track/Album/Artist]`
- Handles user-specific music operations
- Manages user preferences and music history

### ClientMediaMusicHandler

- Uses `services.ClientMediaMusicService[T]`
- Handles communication with external music libraries
- Converts client-specific data to common formats

## Router Configuration

The music routes should be organized in `router/music.go`:

```go
// Core music routes
tracks := rg.Group("/music/tracks")
{
  tracks.GET("/top", coreMusicHandler.GetTopTracks)
  // ...
}

// User-specific music routes
userTracks := rg.Group("/user/music/tracks")
{
  userTracks.GET("/favorites", userMusicHandler.GetFavoriteTracks)
  // ...
}

// Client-specific routes are in router/media.go
// /clients/media/{clientID}/music/...
```

## Usage Guidelines

- For basic music data operations, use `CoreMusicHandler`
- For user-specific music operations, use `UserMusicHandler`
- For client-specific music operations, use `ClientMediaMusicHandler`
- Keep handler methods focused on their specific layer
- Ensure clear boundaries between layers

## Integration with App Dependencies

To fully implement this pattern, the app dependencies need to be updated:

```go
// In app/interfaces.go
type MediaItemServices interface {
  // Core services
  CoreTrackService() services.MediaItemService[*mediatypes.Track]
  CoreAlbumService() services.MediaItemService[*mediatypes.Album]
  CoreArtistService() services.MediaItemService[*mediatypes.Artist]
  
  // User services
  UserTrackService() services.UserMediaItemService[*mediatypes.Track]
  UserAlbumService() services.UserMediaItemService[*mediatypes.Album]
  UserArtistService() services.UserMediaItemService[*mediatypes.Artist]

  // Client services (accessed via ClientMediaServices)
}
```

This three-pronged architecture provides clear separation of concerns and ensures scalability as more music features are added.
