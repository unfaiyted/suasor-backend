// container/container.go
package container

import (
	"fmt"
	"reflect"
	"sync"
	"time"
)

// Container is a dependency injection container
type Container struct {
	components map[reflect.Type]any
	factories  map[reflect.Type]func(c *Container) any
	mutex      sync.RWMutex
	// Add dependency tracking to detect circular dependencies
	resolutionStack map[reflect.Type]bool
}

// NewContainer creates a new DI container
func NewContainer() *Container {
	return &Container{
		components:      make(map[reflect.Type]any),
		factories:       make(map[reflect.Type]func(c *Container) any),
		resolutionStack: make(map[reflect.Type]bool),
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
	fmt.Printf("Container.Get called for type: %v\n", t)

	// First check if component exists with read lock
	c.mutex.RLock()
	component, exists := c.components[t]
	c.mutex.RUnlock()

	if exists {
		fmt.Printf("Found existing component for type: %v\n", t)
		return component, nil
	}

	fmt.Printf("Component not found, checking for factory for type: %v\n", t)

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
		fmt.Printf("⚠️ CIRCULAR DEPENDENCY DETECTED trying to resolve type: %v ⚠️\n", t)
		// Print current stack under read lock
		c.mutex.RLock()
		fmt.Println("Dependency resolution stack:")
		for dep := range c.resolutionStack {
			if c.resolutionStack[dep] {
				fmt.Printf("- %v\n", dep)
			}
		}
		c.mutex.RUnlock()
		return nil, fmt.Errorf("circular dependency detected for type %v", t)
	}

	// No factory found
	if !factoryExists {
		fmt.Printf("No factory found for type: %v\n", t)
		return nil, fmt.Errorf("no component registered for type %v", t)
	}

	// Try to create the component with factory
	// First check if another thread has created it in the meantime
	c.mutex.Lock()
	component, exists = c.components[t]
	if exists {
		fmt.Printf("Component found after re-check for type: %v\n", t)
		c.mutex.Unlock()
		return component, nil
	}

	// Mark as being resolved to detect circular dependencies
	c.resolutionStack[t] = true
	// Release the lock before calling factory to prevent deadlocks
	c.mutex.Unlock()

	fmt.Printf("STARTING: Creating component using factory for type: %v\n", t)

	// Print current resolution stack for debugging (under read lock)
	c.mutex.RLock()
	fmt.Println("Current resolution stack:")
	for dep := range c.resolutionStack {
		if c.resolutionStack[dep] {
			fmt.Printf("- %v\n", dep)
		}
	}
	c.mutex.RUnlock()

	// Call the factory with timeout to avoid infinite hangs
	done := make(chan bool, 1)
	var componentResult any
	var factoryError error

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("PANIC in factory for %v: %v\n", t, r)
				factoryError = fmt.Errorf("factory for %v panicked: %v", t, r)
			}
			done <- true
		}()

		fmt.Printf("Calling factory for type: %v\n", t)
		componentResult = factory(c)
		fmt.Printf("Factory returned for type: %v\n", t)
	}()

	// Wait with timeout to prevent infinite hang
	select {
	case <-done:
		// Factory completed successfully or with panic
		fmt.Printf("Factory for %v completed\n", t)
	case <-time.After(5 * time.Second):
		// Factory is taking too long - likely a hang
		fmt.Printf("⚠️ TIMEOUT: Factory for %v is taking too long (possible hang or deadlock) ⚠️\n", t)
		
		// Provide detailed diagnostics of what's in the container (under read lock)
		c.mutex.RLock()
		fmt.Println("Registered components:")
		for compType := range c.components {
			fmt.Printf("- %v\n", compType)
		}

		fmt.Println("Registered factories:")
		for factoryType := range c.factories {
			fmt.Printf("- %v\n", factoryType)
		}
		c.mutex.RUnlock()

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
	fmt.Printf("FINISHED: Component created for type: %v\n", t)

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
