# Media Handlers Design Pattern

This document outlines the unified three-pronged architecture for all media type handlers in the Suasor backend, including series, music, and playlists.

## Three-Pronged Architecture Overview

All media handlers in the Suasor backend follow a consistent three-pronged architecture:

1. **Core Layer** (`core_*.go`)
   - Basic database operations for media items
   - Provides data access without user context
   - Used by general public APIs

2. **User Layer** (`user_*.go`)
   - User-specific operations for media items
   - Extends core operations with user context
   - Used by personalized user APIs

3. **Client Layer** (`client_media_*.go`)
   - Operations that interact with external media libraries
   - Bridge between the application and media servers
   - Used by client-specific APIs

## File Naming Conventions

To maintain clear and consistent code organization, we follow these file naming conventions:

- **Core handlers**: `core_*.go` (e.g., `core_series.go`, `core_music.go`, `core_playlists.go`)
- **User handlers**: `user_*.go` (e.g., `user_series.go`, `user_music.go`, `user_playlists.go`) 
- **Client handlers**: `client_media_*.go` (e.g., `client_media_series.go`, `client_media_music.go`, `client_media_playlist.go`)

## Routes Organization

Routes are organized to mirror this three-pronged architecture:

- **Core routes**: `/media-type/...` (e.g., `/series/...`, `/music/tracks/...`, `/playlists/...`)
- **User routes**: `/user/media-type/...` (e.g., `/user/series/...`, `/user/music/tracks/...`, `/user/playlists/...`)
- **Client routes**: `/clients/media/{clientID}/media-type/...` (e.g., `/clients/media/{clientID}/series/...`)

## Implementation by Media Type

### Series

1. **Core Series Handler** (`core_series.go`)
   - Basic database operations for series, seasons, and episodes
   - Example endpoints: `/series/{id}`, `/series/{id}/seasons`

2. **User Series Handler** (`user_series.go`) 
   - User-specific operations like favorites, watched status
   - Example endpoints: `/user/series/favorites`, `/user/series/watchlist`

3. **Client Series Handler** (`client_media_series.go`)
   - Operations to retrieve series from external clients
   - Example endpoints: `/clients/media/{clientID}/series/{seriesID}`

### Music

1. **Core Music Handler** (`core_music.go`)
   - Basic database operations for tracks, albums, and artists
   - Example endpoints: `/music/tracks/top`, `/music/albums/{id}/tracks`

2. **User Music Handler** (`user_music.go`)
   - User-specific operations like favorites, recently played
   - Example endpoints: `/user/music/tracks/favorites`, `/user/music/tracks/recently-played`

3. **Client Music Handler** (`client_media_music.go`) 
   - Operations to retrieve music from external clients
   - Example endpoints: `/clients/media/{clientID}/music/tracks`

### Playlists

1. **Core Playlist Handler** (`core_playlists.go`)
   - Basic database operations for playlists
   - Example endpoints: `/playlists/{id}`, `/playlists/{id}/tracks`

2. **User Playlist Handler** (`user_playlists.go`)
   - User-specific operations like creation, modification
   - Example endpoints: `/user/playlists`, `/user/playlists/{id}/tracks`

3. **Client Playlist Handler** (`client_media_playlist.go`)
   - Operations to retrieve playlists from external clients
   - Example endpoints: `/clients/media/{clientID}/playlists`

### Collections

1. **Core Collection Handler** (`core_collections.go`)
   - Basic database operations for collections
   - Example endpoints: `/collections/{id}`, `/collections/{id}/items`, `/collections/public`

2. **User Collection Handler** (`user_collections.go`)
   - User-specific operations like creation, modification, item management
   - Example endpoints: `/user/collections`, `/user/collections/{id}/items`

3. **Client Collection Handler** (`client_media_collection.go`)
   - Operations to retrieve collections from external clients
   - Example endpoints: `/clients/media/{clientID}/collections`

### Movies

1. **Core Movie Handler** (`core_movie.go`)
   - Basic database operations for movies
   - Example endpoints: `/movies/{id}`, `/movies/genre/:genre`, `/movies/top-rated`

2. **User Movie Handler** (`user_movie.go`)
   - User-specific operations like favorites, watchlist, ratings
   - Example endpoints: `/user/movies/favorites`, `/user/movies/watchlist`

3. **Client Movie Handler** (`client_media_movie.go`)
   - Operations to retrieve movies from external clients
   - Example endpoints: `/clients/media/{clientID}/movies`

## Service Dependencies

Each handler type has specific service dependencies:

- **Core handlers** use `services.MediaItemService[T]` or specialized core services
- **User handlers** use `services.UserMediaItemService[T]` or specialized user services
- **Client handlers** use `services.ClientMediaTypeService[T]` (e.g., `MediaClientSeriesService[T]`)

## Benefits of This Architecture

This three-pronged architecture provides several benefits:

1. **Clear separation of concerns**: Each layer has a specific responsibility
2. **Improved code organization**: Consistent patterns make the codebase easier to navigate
3. **Better maintainability**: Changes to one layer don't affect others
4. **Easier testing**: Each layer can be tested independently
5. **Scalability**: New media types can easily follow the same pattern
6. **Developer experience**: Consistent patterns make onboarding new developers easier

## Applying to New Media Types

When implementing a new media type, follow these steps:

1. Create three handler files following the naming convention:
   - `core_<media-type>.go`
   - `user_<media-type>.go`
   - `client_media_<media-type>.go`

2. Create a router file `<media-type>.go` with the three route groups

3. Update app dependencies to expose the required services

4. Register the new routes in the router pipeline

This consistent pattern ensures that all media types are implemented in a uniform way, making the codebase more maintainable and easier to understand.
