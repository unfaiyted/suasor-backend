# Three-Pronged Media Data Architecture

## Overview

The Suasor backend implements a "three-pronged" architecture for handling media data. This design separates concerns into three layers:

1. **Core Layer**: Manages basic storage and retrieval of media data
2. **User Layer**: Handles user-specific operations and data
3. **Client Layer**: Interfaces with external media clients (Plex, Jellyfin, etc.)

This architecture is supported by a comprehensive dependency injection system that simplifies component wiring and promotes clean separation of concerns.

## Key Components

### MediaDataFactory

The `MediaDataFactory` is the central component for creating properly configured repositories, services, and handlers for media data operations. It provides factory methods for all three layers of the architecture.

```go
// MediaDataFactory is a factory for creating components in the three-pronged architecture
type MediaDataFactory struct {
    db               *gorm.DB
    clientFactory    *client.ClientFactoryService
    coreRepositories CoreMediaItemRepositories
}
```

Factory methods include:
- `CreateCoreRepositories()`: Creates repositories for basic media operations
- `CreateUserRepositories()`: Creates repositories for user-specific operations
- `CreateClientRepositories()`: Creates repositories for client-specific operations
- `CreateCoreServices()`: Creates core layer services
- `CreateUserServices()`: Creates user layer services
- `CreateClientServices()`: Creates client layer services
- `CreateCoreHandlers()`: Creates handlers for the core layer
- `CreateUserHandlers()`: Creates handlers for the user layer
- `CreateClientHandlers()`: Creates handlers for the client layer

### ServiceRegistrar

The `ServiceRegistrar` manages the registration and initialization of all services in the system. It follows a specific order to ensure dependencies are properly initialized.

```go
// ServiceRegistrar is responsible for registering all services in the dependency injection system
type ServiceRegistrar struct {
    db            *gorm.DB
    clientFactory *client.ClientFactoryService
    dependencies  *AppDependencies
}
```

Key methods:
- `RegisterCoreServices()`: Initializes core infrastructure services
- `RegisterRepositories()`: Sets up all repositories
- `RegisterMediaDataServices()`: Creates all media data services
- `RegisterMediaDataHandlers()`: Creates all media data handlers
- `RegisterAllServices()`: Orchestrates the entire initialization process

### AppDependencies

The `AppDependencies` struct holds all application components and serves as the central registry for dependency injection.

```go
// AppDependencies contains all application dependencies
// It uses a clean three-pronged architecture for media data dependencies
type AppDependencies struct {
    // Database connection
    db *gorm.DB

    // Core infrastructure repositories and services
    SystemRepositories
    UserRepositories
    JobRepositories

    // Three-pronged architecture for media data
    // Repository layer
    CoreMediaItemRepositories         // Base storage layer
    CoreUserMediaItemDataRepositories // Core-User Data storage layer
    UserRepositoryFactories           // User-specific storage layer
    UserDataFactories                 // User-specific user data storage layer
    ClientRepositoryFactories         // Client-specific storage layer
    ClientUserDataRepositories        // Client-specific user data storage layer

    // Service layer
    CoreMediaItemServices         // Core business logic
    CoreUserMediaItemDataServices // Core-User Data business logic
    UserMediaItemServices         // User-specific business logic
    UserMediaItemDataServices     // User-specific user data logic
    ClientMediaItemServices       // Client-specific business logic
    ClientUserMediaItemDataServices // Client-specific user data logic
    MediaCollectionServices         // Collection/playlist specialized services

    // Handler layer (presentation)
    CoreMediaItemHandlers     // Core API endpoints
    CoreMediaItemDataHandlers // Core-Data API endpoints
    UserMediaItemHandlers     // User-specific API endpoints
    ClientMediaItemHandlers   // Client-specific API endpoints
    SpecializedMediaHandlers  // Domain-specific API endpoints

    // Standard services and handlers
    UserServices
    SystemServices
    ClientServices
    // ... additional fields
}
```

## Dependency Injection Flow

### Initialization Order

1. Core services are initialized first (`RegisterCoreServices`)
2. Repositories are initialized (`RegisterRepositories`)
3. Media data services with the three-pronged approach are created (`RegisterMediaDataServices`)
4. Media data handlers are initialized (`RegisterMediaDataHandlers`)
5. Standard handlers are created (`RegisterStandardHandlers`)

### Layer Organization

Each layer builds on the previous layer:

1. **Core Layer**: Most basic operations (CRUD)
   - Repositories: `CoreMediaItemRepositories`
   - Services: `CoreMediaItemServices`
   - Handlers: `CoreMediaItemHandlers`

2. **User Layer**: Extends core with user-specific operations
   - Repositories: `UserRepositoryFactories`, `UserDataFactories`
   - Services: `UserMediaItemServices`, `UserMediaItemDataServices`
   - Handlers: `UserMediaItemHandlers`

3. **Client Layer**: Extends user with client-specific operations
   - Repositories: `ClientRepositoryFactories`, `ClientUserDataRepositories`
   - Services: `ClientMediaItemServices`, `ClientUserMediaItemDataServices`
   - Handlers: `ClientMediaItemHandlers`

## Interface Design

The architecture uses Go's type parameters (generics) extensively to maintain type safety across the system while reducing code duplication.

### Repository Interfaces

```go
// Core-layer repository
type MediaItemRepository[T types.MediaData] interface {
    Create(ctx context.Context, item models.MediaItem[T]) (*models.MediaItem[T], error)
    Update(ctx context.Context, item models.MediaItem[T]) (*models.MediaItem[T], error)
    GetByID(ctx context.Context, id uint64) (*models.MediaItem[T], error)
    // ... additional methods
}

// User-layer repository
type UserMediaItemRepository[T types.MediaData] interface {
    // Extends core with user-specific operations
    GetByUserID(ctx context.Context, userID uint64) ([]*models.MediaItem[T], error)
    // ... additional methods
}

// Client-layer repository
type ClientMediaItemRepository[T types.MediaData] interface {
    // Extends user with client-specific operations
    GetByClientItemID(ctx context.Context, clientItemID string, clientID uint64) (*models.MediaItem[T], error)
    // ... additional methods
}
```

### Service Interfaces

```go
// Core-layer service
type CoreMediaItemService[T types.MediaData] interface {
    Create(ctx context.Context, item models.MediaItem[T]) (*models.MediaItem[T], error)
    Update(ctx context.Context, item models.MediaItem[T]) (*models.MediaItem[T], error)
    GetByID(ctx context.Context, id uint64) (*models.MediaItem[T], error)
    // ... additional methods
}

// User-layer service
type UserMediaItemService[T types.MediaData] interface {
    // Extends core with user-specific operations
    GetForUser(ctx context.Context, userID uint64, limit int) ([]*models.MediaItem[T], error)
    // ... additional methods
}

// Client-layer service
type ClientMediaItemService[T types.MediaData] interface {
    // Extends user with client-specific operations
    SyncFromClient(ctx context.Context, clientID uint64) error
    // ... additional methods
}
```

### Handler Interfaces

```go
// Core-layer handler
type CoreUserMediaItemDataHandler[T types.MediaData] struct {
    service CoreMediaItemService[T]
}

// User-layer handler
type UserUserMediaItemDataHandler[T types.MediaData] struct {
    service UserMediaItemDataService[T]
    coreHandler *CoreUserMediaItemDataHandler[T]
}

// Client-layer handler
type ClientUserMediaItemDataHandler[T types.MediaData] struct {
    service ClientUserMediaItemDataService[T]
    userHandler *UserUserMediaItemDataHandler[T]
}
```

## Using the Architecture

### Adding a New Media Type

1. Define the media type struct implementing the `types.MediaData` interface
2. Register the type in all three-pronged repositories/services/handlers
3. Update the `MediaDataFactory` methods to include the new type

### Creating a New Service

1. Define the service interface with appropriate type parameters
2. Implement the interface for each layer (core, user, client)
3. Register the implementation in `ServiceRegistrar`
4. Add the service to the `AppDependencies` struct

### Extending the Architecture

The three-pronged approach can be extended with additional specialized services:

1. Create domain-specific interfaces extending the base interfaces
2. Implement the specialized services
3. Register them in the `MediaDataFactory` and `ServiceRegistrar`
4. Add appropriate fields to `AppDependencies`

## Best Practices

1. **Layer Separation**: Each layer should only depend on its own layer and the layer below it
2. **Interface Design**: Define clean interfaces before implementing
3. **Factory Methods**: Use factory methods for creating properly configured components
4. **Registration Order**: Follow the proper initialization order in the `ServiceRegistrar`
5. **Type Safety**: Leverage Go's type parameters for type-safe generic code
6. **Dependency Injection**: Components should receive their dependencies through constructors
7. **Testing**: Each layer can be tested independently with mock implementations

## Example: Media Collection Flow

To understand how this architecture works in practice, consider the flow for retrieving a user's media collections:

1. The router calls the appropriate handler method in `ClientMediaItemHandlers`
2. The client handler delegates to the user handler after client-specific processing
3. The user handler delegates to the core handler after user-specific processing
4. The core handler uses its service to retrieve the data
5. The data flows back up through the chain, with each layer adding its specific processing
6. The final result is returned to the client

This layered approach ensures separation of concerns while allowing each layer to focus on its specific responsibilities.