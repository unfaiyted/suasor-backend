# Media Handler Design Pattern

This document outlines the three-pronged architecture for media handlers in the Suasor backend.

## Three-Pronged Architecture

The backend implements a three-pronged architecture for handling media items:

1. **Core Layer**: Basic database operations for media items
2. **User Layer**: User-specific operations and data for media items
3. **Client Layer**: Operations that interact with external media clients (Plex, Jellyfin, etc.)

## Handler Organization

For each media type (movies, series, music, etc.), we have three handler types:

### 1. Core Handlers (`CoreXXXHandler`)

- Focus on basic database operations
- Operate on the application's internal database
- Endpoints usually have patterns like `/movies`, `/series`, etc.
- Use `CoreMediaItemService` for data access
- Example: `CoreSeriesHandler` for TV series in the database

### 2. User Handlers (`UserXXXHandler`)

- Focus on user-specific operations (favorites, watch history, ratings)
- Extend core operations with user-specific context
- Endpoints usually have patterns like `/user/movies`, `/user/series`, etc.
- Use `UserMediaItemService` for data access
- Example: `UserSeriesHandler` for managing user's series preferences

### 3. Client Handlers (`MediaClientXXXHandler`)

- Focus on operations that interact with external media clients
- Bridge between the application and media servers
- Endpoints usually have patterns like `/clients/media/{clientID}/movies`, etc.
- Use `MediaClientXXXService` for data access
- Example: `MediaClientSeriesHandler` for accessing series from Plex/Jellyfin

## Handler Responsibilities

Each handler type has specific responsibilities:

### CoreXXXHandler

- CRUD operations on media items
- Basic filtering (by genre, year, etc.)
- Search operations across the database
- Metadata operations

### UserXXXHandler

- User favorites management
- Watch history tracking
- User ratings and reviews
- User collections management
- Personalization features

### MediaClientXXXHandler

- Retrieve media from external clients
- Client-specific filters and search
- Client synchronization
- Client-specific features

## Implementation Guidance

When implementing a new feature:

1. Determine which layer(s) it belongs to
2. Implement service methods in the appropriate service layer
3. Add handler methods in the corresponding handler
4. Register routes in the appropriate router file

For cross-layer features, start with the core implementation and extend to user/client as needed.

## Example: Series Handlers

- **CoreSeriesHandler**: Basic series database operations
  - Get series by ID
  - List all series
  - Get seasons/episodes
  - Filter by metadata

- **UserSeriesHandler**: User-specific series operations
  - Get favorite series
  - Get watched series
  - Get series in watchlist
  - Update user data for series

- **ClientSeriesHandler**: Client-specific series operations
  - Retrieve series from a specific client
  - Client-specific search/filtering
  - Get seasons/episodes from client

## Routing Examples

```
# Core routes
/series/{id} - Get series by ID
/series/genre/{genre} - Get series by genre

# User routes
/user/series/favorites - Get user favorite series
/user/series/watchlist - Get user watchlist series

# Client routes
/clients/media/{clientID}/series/{seriesID} - Get series from specific client
/clients/media/{clientID}/series/search - Search series on specific client
```

## Service Dependencies

- **CoreXXXHandler** depends on `CoreMediaItemService`
- **UserXXXHandler** depends on `UserMediaItemService` (which extends CoreMediaItemService)
- **ClientMediaXXXHandler** depends on `ClientMediaXXXService`
