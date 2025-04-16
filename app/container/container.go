// container/container.go
package container

import (
	"fmt"
	"reflect"
	"sync"
)

// Container is a dependency injection container
type Container struct {
	components map[reflect.Type]any
	factories  map[reflect.Type]func(c *Container) any
	mutex      sync.RWMutex
}

// NewContainer creates a new DI container
func NewContainer() *Container {
	return &Container{
		components: make(map[reflect.Type]any),
		factories:  make(map[reflect.Type]func(c *Container) any),
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
	// (simplified for brevity)
	c.mutex.RLock()
	component, exists := c.components[t]
	c.mutex.RUnlock()

	if exists {
		return component, nil
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	factory, factoryExists := c.factories[t]
	if !factoryExists {
		return nil, fmt.Errorf("no component registered for type %v", t)
	}

	component = factory(c)
	c.components[t] = component
	return component, nil
}
