# Dependency Injection in Suasor

This document explains the dependency injection (DI) pattern used in the Suasor backend, including the three-pronged architecture and how dependencies are wired together.

## Table of Contents

- [Overview](#overview)
- [Container Structure](#container-structure)
- [Three-Pronged Architecture](#three-pronged-architecture)
- [Dependency Initialization](#dependency-initialization)
- [Usage Examples](#usage-examples)
- [Adding New Dependencies](#adding-new-dependencies)
- [Testing with the DI Container](#testing-with-the-di-container)

## Overview

Suasor uses a dependency injection container to manage and wire together the various components of the application. This approach offers several benefits:

- **Loose coupling** between components
- **Improved testability** through easier mocking
- **Centralized dependency management**
- **Consistent object lifecycle**
- **Type-safe access** to services and repositories

The DI container is implemented as the `AppDependencies` struct in `app/dependencies.go`, which holds references to all services, repositories, and handlers used throughout the application.

## Container Structure

The container is organized into several main sections:

### Core Infrastructure

- **Database connection**: `db *gorm.DB`
- **System repositories**: Configuration, health checks, etc.
- **User repositories**: User accounts, sessions, etc.
- **Job repositories**: Background jobs and scheduling

### Three-Pronged Architecture Components

The media handling uses a three-pronged architecture:

#### Repository Layer
- **Core media repositories**: Base storage layer (`CoreMediaItemRepositories`)
- **User repositories**: User-specific storage (`UserRepositoryFactories`)
- **Client repositories**: Client-specific storage (`ClientRepositoryFactories`)
- **Data repositories**: User and client data storage (`CoreUserMediaItemDataRepositories`, `UserDataFactories`, `ClientUserDataRepositories`)

#### Service Layer
- **Core services**: Base business logic (`CoreMediaItemServices`)
- **User services**: User-specific logic (`UserMediaItemServices`, `CoreUserMediaItemDataServices`)
- **Client services**: Client-specific logic (`ClientMediaItemServices`, `ClientUserMediaItemDataServices`)
- **Specialized services**: Domain-specific logic (`MediaCollectionServices`)

#### Handler Layer (API Presentation)
- **Core handlers**: Base API endpoints (`CoreMediaItemHandlers`)
- **User handlers**: User-specific endpoints (`UserMediaItemHandlers`, `CoreMediaItemDataHandlers`)
- **Client handlers**: Client-specific endpoints (`ClientMediaItemHandlers`)
- **Specialized handlers**: Domain-specific endpoints (`SpecializedMediaHandlers`)

### Standard Application Components

- **User services**: Authentication, user profiles, etc.
- **System services**: Configuration, health checks, etc.
- **Client services**: External client management
- **Job services**: Background job processing

## Three-Pronged Architecture

The three-pronged architecture is a core design pattern in Suasor that organizes media handling into three layers:

### 1. Core Layer

The foundation layer that provides:
- Basic data structures and models
- Database CRUD operations
- Data validation and integrity checks
- Media type handling (movies, series, music, etc.)

### 2. User Layer

Extends the core layer with user-specific features:
- User ownership and access control
- User preferences and settings
- User-specific data (ratings, watch status, etc.)
- Collection management

### 3. Client Layer

Extends the user layer with client-specific features:
- Client connection handling (Jellyfin, Emby, Plex, etc.)
- Client-specific IDs and mappings
- External API integration
- Synchronization logic

Each layer builds upon the previous one, creating a hierarchy of functionality. This architecture allows for flexibility and maintainability by:

- Isolating concerns at each layer
- Allowing new client types to be added without modifying core logic
- Enabling testing at each layer independently
- Providing clear boundaries between components

## Dependency Initialization

Dependencies are initialized in `app/init_dependencies.go` through the `InitializeDependencies` function. The initialization process follows these steps:

1. Create the container struct with the database connection
2. Initialize system repositories and services
3. Initialize user repositories and services
4. Create the client factory service
5. Initialize core repositories using the media data factory
6. Initialize user repositories using the factory
7. Initialize client repositories using the factory
8. Create services for each layer (core, user, client)
9. Initialize specialized services (collections, playlists)
10. Create handlers for each layer and domain

The initialization uses dependency injection throughout, passing dependencies to constructors rather than allowing components to create their own dependencies.

## Usage Examples

### Accessing Services

```go
// Get a service from the container
jellyfinService := deps.ClientServices.JellyfinService()

// Use the service
clients, err := jellyfinService.GetAllForUser(ctx, userID)
```

### Using the Three-Pronged Architecture

```go
// Core layer (database only)
coreMovieService := deps.CoreMediaItemServices.MovieCoreService()
movie, err := coreMovieService.GetByID(ctx, movieID)

// User layer (user-specific data)
userMovieService := deps.UserMediaItemServices.MovieUserService()
userMovie, err := userMovieService.GetByID(ctx, userID, movieID)

// Client layer (client-specific data)
clientMovieService := deps.ClientMediaItemServices.MovieClientService()
clientMovie, err := clientMovieService.GetByID(ctx, userID, clientID, movieID)
```

### Accessing Repositories

```go
// Direct repository access
movieRepo := deps.CoreMediaItemRepositories.MovieRepo()
allMovies, err := movieRepo.GetAll(ctx, &repository.QueryOptions{Limit: 10})
```

## Adding New Dependencies

To add new dependencies to the container:

1. Define the interface in `app/interfaces.go`
2. Create the concrete implementation in an appropriate file (e.g., `app/client_impl.go`)
3. Add the dependency to the `AppDependencies` struct in `app/dependencies.go`
4. Initialize the dependency in `app/init_dependencies.go`

Example for adding a new client type:

```go
// 1. Define the interface
type NewClientServices interface {
    NewClientService() services.ClientService[*types.NewClientConfig]
}

// 2. Create the implementation
type newClientServicesImpl struct {
    newClientService services.ClientService[*types.NewClientConfig]
}

func (s *newClientServicesImpl) NewClientService() services.ClientService[*types.NewClientConfig] {
    return s.newClientService
}

// 3. Add to AppDependencies
type AppDependencies struct {
    // ... existing dependencies
    NewClientServices NewClientServices
}

// 4. Initialize in InitializeDependencies
func InitializeDependencies(db *gorm.DB, configService services.ConfigService) *AppDependencies {
    // ... existing initialization
    
    // Create the repository
    newClientRepo := repository.NewClientRepository[*types.NewClientConfig](db)
    
    // Initialize the service
    deps.NewClientServices = &newClientServicesImpl{
        newClientService: services.NewClientService[*types.NewClientConfig](
            deps.ClientFactoryService, 
            newClientRepo),
    }
    
    return deps
}
```

## Testing with the DI Container

For testing, you can create a simplified version of the container with mocked dependencies:

```go
func createTestDependencies() *app.AppDependencies {
    // Create a mock DB
    db, _ := database.ConnectDatabase("sqlite::memory:", false)
    
    // Create a mock config service
    configRepo := &mocks.MockConfigRepository{}
    configService := services.NewConfigService(configRepo)
    
    // Initialize dependencies
    deps := app.InitializeDependencies(db, configService)
    
    // Replace specific dependencies with mocks if needed
    mockUserRepo := &mocks.MockUserRepository{}
    deps.UserRepositories = &app.userRepositoriesImpl{
        userRepo: mockUserRepo,
        // Keep other repositories...
    }
    
    return deps
}

func TestSomeFunctionality(t *testing.T) {
    deps := createTestDependencies()
    
    // Use the container in your tests
    service := deps.UserServices.UserService()
    
    // Test the service...
}
```

The DI container approach makes it easy to replace real implementations with mocks for testing.