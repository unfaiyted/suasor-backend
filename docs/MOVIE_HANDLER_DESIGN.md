# Movie Handler Design Pattern

This document outlines the three-pronged architecture for movie handlers in the Suasor backend, following the same pattern established for series, music, playlists, and collections.

## Three-Pronged Architecture for Movies

The movie handlers implement a three-pronged architecture:

1. **Core Layer** (`core_movie.go`)
   - Basic database operations for movies
   - Provides data access without user context
   - Used by general public APIs

2. **User Layer** (`user_movie.go`)
   - User-specific operations for movies (favorites, watchlist, ratings)
   - Extends core operations with user context
   - Used by personalized user APIs

3. **Client Layer** (`client_media_movie.go`)
   - Operations that interact with external movie libraries (Plex, Jellyfin, etc.)
   - Bridge between the application and media servers
   - Used by client-specific APIs

## Handler Organization

### 1. Core Movie Handler (`CoreMovieHandler`)

- Responsible for basic database operations
- Provides read access to movies in the database
- Example endpoints:
  - `GET /movies` - Get all movies
  - `GET /movies/:id` - Get movie by ID
  - `GET /movies/genre/:genre` - Get movies by genre
  - `GET /movies/top-rated` - Get top rated movies

### 2. User Movie Handler (`UserMovieHandler`)

- Responsible for user-specific movie operations
- Manages user favorites, watched status, watchlist, and ratings
- Example endpoints:
  - `GET /user/movies/favorites` - Get user favorite movies
  - `GET /user/movies/watched` - Get user watched movies
  - `GET /user/movies/watchlist` - Get user watchlist movies
  - `PATCH /user/movies/:id` - Update user data for a movie

### 3. Client Movie Handler (`ClientMediaMovieHandler`)

- Responsible for retrieving movies from external media clients
- Manages access to Plex, Jellyfin, Emby movie libraries
- Example endpoints:
  - `GET /clients/media/{clientID}/movies` - Get movies from client
  - `GET /clients/media/{clientID}/movies/:id` - Get movie by ID from client
  - `GET /clients/media/{clientID}/movies/genre/:genre` - Get movies by genre from client

## Implementation Details

### CoreMovieHandler

- Uses `services.MediaItemService[*mediatypes.Movie]`
- Handles database queries and responses
- Focuses on read-only operations for public movie data

### UserMovieHandler

- Uses `services.UserMediaItemService[*mediatypes.Movie]`
- Handles user-specific movie operations
- Manages user preferences for movies

### ClientMediaMovieHandler

- Uses `services.ClientMediaMovieService[T]`
- Handles communication with external movie libraries
- Converts client-specific movie data to common formats

## Router Configuration

The movie routes are organized in `router/movies.go`:

```go
// Core movie routes
movies := rg.Group("/movies")
{
  movies.GET("", coreMovieHandler.GetAll)
  // ...
}

// User-specific movie routes
userMovies := rg.Group("/user/movies")
{
  userMovies.GET("/favorites", userMovieHandler.GetFavoriteMovies)
  // ...
}

// Client-specific routes are in router/media.go
// /clients/media/{clientID}/movies/...
```

## Usage Guidelines

- For basic movie data operations, use `CoreMovieHandler`
- For user-specific movie operations, use `UserMovieHandler`
- For client-specific movie operations, use `ClientMediaMovieHandler`
- Keep handler methods focused on their specific layer
- Ensure clear boundaries between layers

## Integration with App Dependencies

To fully implement this pattern, the app dependencies need to be updated:

```go
// In app/interfaces.go
type MediaItemServices interface {
  // Core services
  CoreMovieService() services.MediaItemService[*mediatypes.Movie]
  
  // User services
  UserMovieService() services.UserMediaItemService[*mediatypes.Movie]

  // Client services (accessed via ClientMediaServices)
}
```

This three-pronged architecture provides clear separation of concerns, improves code organization, and ensures scalability as more movie features are added.
