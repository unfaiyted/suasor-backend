// app/container/helpers.go
package container

import (
	"fmt"
	"reflect"
	"sync"
)

// GetTyped retrieves a component by its type with proper type safety
func GetTyped[T any](c *Container) (T, error) {
	var zero T
	t := reflect.TypeOf((*T)(nil)).Elem()
	fmt.Printf("GetTyped called for type: %v\n", t)
	
	component, err := c.Get(t)
	if err != nil {
		fmt.Printf("GetTyped failed for type %v with error: %v\n", t, err)
		return zero, err
	}
	
	fmt.Printf("GetTyped succeeded for type: %v\n", t)
	return component.(T), nil
}

// MustGet retrieves a component or panics if not found
func MustGet[T any](c *Container) T {
	t := reflect.TypeOf((*T)(nil)).Elem()
	fmt.Printf("MustGet called for type: %v\n", t)
	
	result, err := GetTyped[T](c)
	if err != nil {
		fmt.Printf("MustGet failed for type %v with error: %v\n", t, err)
		panic(fmt.Sprintf("Failed to get component: %v", err))
	}
	
	fmt.Printf("MustGet succeeded for type: %v\n", t)
	return result
}

// RegisterImpl registers an implementation for an interface
func RegisterImpl[I any](c *Container, impl I) {
	t := reflect.TypeOf((*I)(nil)).Elem()
	c.RegisterInterface(t, impl)
}

// RegisterFactory registers a factory function for creating components
func RegisterFactory[T any](c *Container, factory func(c *Container) T) {
	t := reflect.TypeOf((*T)(nil)).Elem()

	c.RegisterFactory(t, func(c *Container) any {
		return factory(c)
	})
}

// RegisterSingleton registers a factory that will be invoked only once
func RegisterSingleton[T any](c *Container, factory func(c *Container) T) {
	var instance T
	var once sync.Once

	RegisterFactory(c, func(c *Container) T {
		once.Do(func() {
			instance = factory(c)
		})
		return instance
	})
}

