// container/container.go
package container

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"
	
	"github.com/rs/zerolog"
	"suasor/utils/logger"
)

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

// NewContainer creates a new DI container
func NewContainer() *Container {
	ctx := context.Background()
	return &Container{
		components:      make(map[reflect.Type]any),
		factories:       make(map[reflect.Type]func(c *Container) any),
		resolutionStack: make(map[reflect.Type]bool),
		ctx:             ctx,
		logger:          logger.LoggerFromContext(ctx),
	}
}

// NewContainerWithContext creates a new DI container with a specified context
func NewContainerWithContext(ctx context.Context) *Container {
	return &Container{
		components:      make(map[reflect.Type]any),
		factories:       make(map[reflect.Type]func(c *Container) any),
		resolutionStack: make(map[reflect.Type]bool),
		ctx:             ctx,
		logger:          logger.LoggerFromContext(ctx),
	}
}

// Register adds a component to the container
func (c *Container) Register(component any) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	t := reflect.TypeOf(component)
	c.components[t] = component
}

// RegisterInterface registers a component under an interface type
func (c *Container) RegisterInterface(interfaceType reflect.Type, component any) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.components[interfaceType] = component
}

// RegisterFactory registers a factory function for creating a component
func (c *Container) RegisterFactory(interfaceType reflect.Type, factory func(c *Container) any) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.factories[interfaceType] = factory
}

// Get retrieves a component by type
func (c *Container) Get(t reflect.Type) (any, error) {
	// Implementation with proper locking and factory invocation
	// fmt.Printf("Container.Get called for type: %v\n", t)

	// First check if component exists with read lock
	c.mutex.RLock()
	component, exists := c.components[t]
	c.mutex.RUnlock()

	if exists {
		// fmt.Printf("Found existing component for type: %v\n", t)
		return component, nil
	}

	c.logger.Debug().Type("type", t).Msg("Component not found, checking for factory")

	// Find a factory for this component type
	var factory func(c *Container) any
	var factoryExists bool

	// Get factory with a read lock
	c.mutex.RLock()
	factory, factoryExists = c.factories[t]

	// Check for circular dependencies under read lock
	circularDependency := c.resolutionStack[t]
	c.mutex.RUnlock()

	// Handle circular dependency
	if circularDependency {
		dependencies := []string{}
		
		// Collect dependency stack under read lock
		c.mutex.RLock()
		for dep := range c.resolutionStack {
			if c.resolutionStack[dep] {
				dependencies = append(dependencies, dep.String())
			}
		}
		c.mutex.RUnlock()
		
		// Log the circular dependency error
		c.logger.Error().
			Type("type", t).
			Strs("dependency_stack", dependencies).
			Msg("⚠️ CIRCULAR DEPENDENCY DETECTED ⚠️")
			
		return nil, fmt.Errorf("circular dependency detected for type %v", t)
	}

	// No factory found
	if !factoryExists {
		c.logger.Debug().Type("type", t).Msg("No factory found for type")
		return nil, fmt.Errorf("no component registered for type %v", t)
	}

	// Try to create the component with factory
	// First check if another thread has created it in the meantime
	c.mutex.Lock()
	component, exists = c.components[t]
	if exists {
		c.logger.Debug().Type("type", t).Msg("Component found after re-check")
		c.mutex.Unlock()
		return component, nil
	}

	// Mark as being resolved to detect circular dependencies
	c.resolutionStack[t] = true
	// Release the lock before calling factory to prevent deadlocks
	c.mutex.Unlock()

	// fmt.Printf("STARTING: Creating component using factory for type: %v\n", t)

	// Log current resolution stack for debugging (under read lock)
	depStack := []string{}
	c.mutex.RLock()
	for dep := range c.resolutionStack {
		if c.resolutionStack[dep] {
			depStack = append(depStack, dep.String())
		}
	}
	c.mutex.RUnlock()
	
	c.logger.Debug().
		Type("type", t).
		Strs("resolution_stack", depStack).
		Msg("Resolving dependency")

	// Call the factory with timeout to avoid infinite hangs
	done := make(chan bool, 1)
	var componentResult any
	var factoryError error

	go func() {
		defer func() {
			if r := recover(); r != nil {
				c.logger.Error().
					Type("type", t).
					Interface("panic", r).
					Msg("PANIC in factory")
				factoryError = fmt.Errorf("factory for %v panicked: %v", t, r)
			}
			done <- true
		}()

		componentResult = factory(c)
	}()

	// Wait with timeout to prevent infinite hang
	select {
	case <-done:
		// Factory completed successfully or with panic
		c.logger.Debug().Type("type", t).Msg("Factory completed")
	case <-time.After(5 * time.Second):
		// Factory is taking too long - likely a hang
		components := []string{}
		factories := []string{}
		
		// Provide detailed diagnostics of what's in the container (under read lock)
		c.mutex.RLock()
		for compType := range c.components {
			components = append(components, compType.String())
		}
		for factoryType := range c.factories {
			factories = append(factories, factoryType.String())
		}
		c.mutex.RUnlock()

		// Log detailed timeout information
		c.logger.Error().
			Type("type", t).
			Strs("registered_components", components).
			Strs("registered_factories", factories).
			Strs("resolution_stack", depStack).
			Msg("⚠️ TIMEOUT: Factory is taking too long (possible hang or deadlock) ⚠️")

		// Update resolution stack
		c.mutex.Lock()
		c.resolutionStack[t] = false
		c.mutex.Unlock()

		return nil, fmt.Errorf("timeout resolving dependency %v - possible deadlock", t)
	}

	// If there was a panic or other error in the factory
	if factoryError != nil {
		// Update resolution stack
		c.mutex.Lock()
		c.resolutionStack[t] = false
		c.mutex.Unlock()
		return nil, factoryError
	}

	component = componentResult
	// fmt.Printf("FINISHED: Component created for type: %v\n", t)

	// Now store the component in the map
	c.mutex.Lock()
	// Check again if another thread has created the component while we were waiting
	existingComponent, exists := c.components[t]
	if exists {
		// Another thread beat us to it, use that one
		c.resolutionStack[t] = false
		c.mutex.Unlock()
		return existingComponent, nil
	}

	// Remove from resolution stack and store our new component
	c.resolutionStack[t] = false
	c.components[t] = component
	c.mutex.Unlock()

	return component, nil
}
