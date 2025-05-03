# Container-Based Dependency Injection

---
Created: 2025-05-03
Last Updated: 2025-05-03
Update Frequency: As needed when DI system changes
Owner: Backend Team
---

## Overview

Suasor uses a sophisticated container-based dependency injection (DI) system to manage component lifecycle, resolve dependencies, and promote loose coupling between components. This document explains the container implementation, features, and usage patterns.

## Architecture

The DI container is implemented in the `di/container` package and provides:

1. **Type-Safe Dependency Resolution**: Using Go generics for compile-time type safety
2. **Factory Pattern Support**: Lazy initialization of components
3. **Circular Dependency Detection**: Prevention of infinite dependency loops
4. **Thread Safety**: Concurrent access to the container
5. **Singleton Support**: Single-instance components
6. **Error Handling**: Proper reporting for dependency resolution failures
7. **Timeout Protection**: Protection against factory hangs

### Core Components

#### Container

The main container implementation:

```go
// Container is a dependency injection container
type Container struct {
    components map[reflect.Type]any
    factories  map[reflect.Type]func(c *Container) any
    mutex      sync.RWMutex
    // Add dependency tracking to detect circular dependencies
    resolutionStack map[reflect.Type]bool
    ctx             context.Context
    logger          zerolog.Logger
}
```

#### Helper Functions

The container provides helper functions for type-safe access:

```go
// GetTyped retrieves a component by its type with proper type safety
func GetTyped[T any](c *Container) (T, error)

// MustGet retrieves a component or panics if not found
func MustGet[T any](c *Container) T

// RegisterFactory registers a factory function for creating components
func RegisterFactory[T any](c *Container, factory func(c *Container) T)

// RegisterSingleton registers a factory that will be invoked only once
func RegisterSingleton[T any](c *Container, factory func(c *Container) T)
```

## Key Features

### Type-Safe Registration and Resolution

The container provides type-safe registration and resolution of dependencies using Go generics:

```go
// Register a component
container.RegisterFactory[services.UserService](c, 
    func(c *container.Container) services.UserService {
        userRepo := container.MustGet[repository.UserRepository](c)
        return services.NewUserService(userRepo)
    })

// Resolve a component
userService := container.MustGet[services.UserService](c)
```

### Circular Dependency Detection

The container detects circular dependencies to prevent infinite loops:

```go
// Check for circular dependencies
circularDependency := c.resolutionStack[t]
if circularDependency {
    dependencies := []string{}
    
    // Collect dependency stack
    for dep := range c.resolutionStack {
        if c.resolutionStack[dep] {
            dependencies = append(dependencies, dep.String())
        }
    }
    
    // Log the circular dependency error
    c.logger.Error().
        Type("type", t).
        Strs("dependency_stack", dependencies).
        Msg("⚠️ CIRCULAR DEPENDENCY DETECTED ⚠️")
        
    return nil, fmt.Errorf("circular dependency detected for type %v", t)
}
```

### Factory Timeout Protection

The container includes protection against hanging factory functions:

```go
// Call the factory with timeout to avoid infinite hangs
done := make(chan bool, 1)
var componentResult any
var factoryError error

go func() {
    defer func() {
        if r := recover(); r != nil {
            // Handle panic in factory
            factoryError = fmt.Errorf("factory panicked: %v", r)
        }
        done <- true
    }()

    componentResult = factory(c)
}()

// Wait with timeout to prevent infinite hang
select {
case <-done:
    // Factory completed successfully
case <-time.After(5 * time.Second):
    // Factory is taking too long - likely a hang
    return nil, fmt.Errorf("timeout resolving dependency")
}
```

### Thread Safety

All container operations are thread-safe using appropriate locking mechanisms:

```go
// Thread-safe registration
func (c *Container) Register(component any) {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    t := reflect.TypeOf(component)
    c.components[t] = component
}

// Thread-safe resolution with optimistic read-first approach
func (c *Container) Get(t reflect.Type) (any, error) {
    // First check if component exists with read lock
    c.mutex.RLock()
    component, exists := c.components[t]
    c.mutex.RUnlock()

    if exists {
        return component, nil
    }
    
    // Component creation logic with proper locking...
}
```

## Usage Patterns

### Basic Component Registration

```go
// Register a simple component
db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
container.Register(db)
```

### Factory Registration

```go
// Register a component with dependencies
container.RegisterFactory[services.UserService](c, 
    func(c *container.Container) services.UserService {
        userRepo := container.MustGet[repository.UserRepository](c)
        configService := container.MustGet[services.ConfigService](c)
        return services.NewUserService(userRepo, configService)
    })
```

### Singleton Registration

```go
// Register a singleton component
container.RegisterSingleton[services.ConfigService](c,
    func(c *container.Container) services.ConfigService {
        return services.NewConfigService()
    })
```

### Interface Implementation Registration

```go
// Register a concrete implementation for an interface
container.RegisterImpl[repository.UserRepository](c, 
    repository.NewUserRepository(db))
```

### Type Key Registration

```go
// Register a component with a string type key
container.RegisterTypedFactory[services.ClientService](c, 
    "plex",
    func(c *container.Container) services.ClientService {
        return services.NewPlexClientService()
    })

// Retrieve by type key
plexService, err := container.GetByTypeKey[services.ClientService](c, "plex")
```

## Application Bootstrap

The application bootstraps the DI container in the `di/register.go` file:

```go
// Initialize registers all dependencies in the container
func RegisterAppContainers(ctx context.Context, db *gorm.DB, configService services.ConfigService) *container.Container {
    // Create a new container with the provided context
    c := container.NewContainerWithContext(ctx)
    
    // Register core dependencies
    RegisterCore(ctx, c, db, configService)
    
    // Register repositories
    direpos.RegisterRepositories(ctx, c)
    
    // Register services
    diservices.RegisterServices(ctx, c)
    
    // Register handlers
    dihandlers.RegisterHandlers(ctx, c)
    
    return c
}
```

## Dependency Organization

Dependencies are organized into logical groups:

1. **Core Dependencies**: Database, configuration, and other core infrastructure
2. **Repository Dependencies**: Data access components
3. **Service Dependencies**: Business logic components
4. **Handler Dependencies**: API endpoint handlers

Each group has a dedicated registration function:

```go
// Register repositories
func RegisterRepositories(ctx context.Context, c *container.Container) {
    registerSystemRepositories(ctx, c)
    registerClientRepositories(ctx, c)
    RegisterMediaRepositories(ctx, c)
    registerMediaListRepositories(ctx, c)
    registerJobRepositories(ctx, c)
}

// Register services
func RegisterServices(ctx context.Context, c *container.Container) {
    registerSystemServices(ctx, c)
    registerClientServices(ctx, c)
    registerMediaItemServices(ctx, c)
    registerMediaDataServices(ctx, c)
    registerListServices(ctx, c)
    registerJobServices(ctx, c)
    registerSearchService(ctx, c)
    registerRecommendationService(ctx, c)
}
```

## Three-Pronged Architecture Integration

The container integrates with Suasor's three-pronged architecture by organizing dependencies into core, user, and client layers:

```go
// Register core media item repositories
container.RegisterFactory[repository.CoreMediaItemRepository[*types.Movie]](c, 
    func(c *container.Container) repository.CoreMediaItemRepository[*types.Movie] {
        db := container.MustGet[*gorm.DB](c)
        return repository.NewMediaItemRepository[*types.Movie](db)
    })

// Register user media item repositories (depends on core)
container.RegisterFactory[repository.UserMediaItemRepository[*types.Movie]](c,
    func(c *container.Container) repository.UserMediaItemRepository[*types.Movie] {
        db := container.MustGet[*gorm.DB](c)
        coreRepo := container.MustGet[repository.CoreMediaItemRepository[*types.Movie]](c)
        return repository.NewUserMediaItemRepository[*types.Movie](db, coreRepo)
    })

// Register client media item repositories (depends on user)
container.RegisterFactory[repository.ClientMediaItemRepository[*types.Movie]](c,
    func(c *container.Container) repository.ClientMediaItemRepository[*types.Movie] {
        db := container.MustGet[*gorm.DB](c)
        userRepo := container.MustGet[repository.UserMediaItemRepository[*types.Movie]](c)
        return repository.NewClientMediaItemRepository[*types.Movie](db, userRepo)
    })
```

## Generic Components

The container supports generic components, which is crucial for Suasor's type-safe media handlers:

```go
// Register a generic repository for all client types
func registerClientRepository[T clienttypes.ClientConfig](c *container.Container, db *gorm.DB) {
    container.RegisterFactory[repository.ClientRepository[T]](c, 
        func(c *container.Container) repository.ClientRepository[T] {
            return repository.NewClientRepository[T](db)
        })
}

// Use it for multiple client types
registerClientRepository[*clienttypes.EmbyConfig](c, db)
registerClientRepository[*clienttypes.JellyfinConfig](c, db)
registerClientRepository[*clienttypes.PlexConfig](c, db)
```

## Best Practices

### Dependency Registration

1. **Register in Logical Groups**: Organize registrations by functional area
2. **Use Factories for Complex Dependencies**: Use factory functions for components with dependencies
3. **Use Singletons for Stateless Services**: Register stateless services as singletons
4. **Register Interfaces, Not Implementations**: When possible, register interfaces to promote loose coupling

### Dependency Resolution

1. **Resolve Late**: Resolve dependencies at the latest possible moment
2. **Use Error Handling**: Prefer `GetTyped` over `MustGet` when appropriate
3. **Consider Circular Dependencies**: Design components to avoid circular dependencies
4. **Use Type Keys for Dynamic Resolution**: Use `GetByTypeKey` for dynamically resolved components

## Common Patterns

### Bundle Pattern

Suasor uses a bundle pattern to group related dependencies:

```go
// Register a bundle of client repositories
container.RegisterFactory[repobundles.ClientRepositories](c, 
    func(c *container.Container) repobundles.ClientRepositories {
        embyRepo := container.MustGet[repository.ClientRepository[*clienttypes.EmbyConfig]](c)
        jellyfinRepo := container.MustGet[repository.ClientRepository[*clienttypes.JellyfinConfig]](c)
        // More repositories...
        
        return repobundles.NewClientRepositories(
            embyRepo, jellyfinRepo, // More repositories...
        )
    })
```

### Factory Pattern

Components are created through factories to control instantiation:

```go
// Client factory example
container.RegisterFactory[services.ClientFactoryService](c, 
    func(c *container.Container) services.ClientFactoryService {
        clientRepos := container.MustGet[repobundles.ClientRepositories](c)
        return services.NewClientFactoryService(
            clientRepos.EmbyRepository(),
            clientRepos.JellyfinRepository(),
            // More repositories...
        )
    })
```

## Troubleshooting

### Common Issues

1. **Circular Dependencies**: Detected when components depend on each other
   - **Solution**: Refactor to break the dependency cycle

2. **Missing Dependencies**: Component not found during resolution
   - **Solution**: Check registration functions to ensure the component is registered

3. **Type Mismatches**: Component type doesn't match the expected type
   - **Solution**: Ensure proper type parameters are used in registration and resolution

4. **Factory Timeouts**: Factory function takes too long to execute
   - **Solution**: Optimize factory functions, ensure no blocking operations

### Debugging

The container includes extensive logging to help debug dependency issues:

```go
c.logger.Debug().
    Type("type", t).
    Strs("resolution_stack", depStack).
    Msg("Resolving dependency")
```

Set the log level to debug to see detailed dependency resolution information.

## Conclusion

Suasor's container-based dependency injection system provides a robust foundation for managing component dependencies. It promotes loose coupling, improves testability, and ensures type safety through Go's generics system. The system is designed to handle the complex dependencies required by Suasor's three-pronged architecture and support for multiple client types.