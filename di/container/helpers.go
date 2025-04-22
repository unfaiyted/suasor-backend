// container/helpers.go
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

	c.logger.Debug().Type("type", t).Msg("GetTyped called")

	component, err := c.Get(t)
	if err != nil {
		c.logger.Debug().
			Type("type", t).
			Err(err).
			Msg("GetTyped failed")
		return zero, err
	}

	c.logger.Debug().Type("type", t).Msg("GetTyped succeeded")
	return component.(T), nil
}

// MustGet retrieves a component or panics if not found
func MustGet[T any](c *Container) T {
	t := reflect.TypeOf((*T)(nil)).Elem()
	c.logger.Debug().Type("type", t).Msg("MustGet called")

	result, err := GetTyped[T](c)
	if err != nil {
		c.logger.Error().
			Type("type", t).
			Err(err).
			Msg("MustGet failed")
		panic(fmt.Sprintf("Failed to get component: %v", err))
	}

	c.logger.Debug().Type("type", t).Msg("MustGet succeeded")
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

// RegisterTypedFactory registers a factory function with a specific type key
func RegisterTypedFactory[T any](c *Container, typeKey string, factory func(c *Container) T) {
	// Create a type that can be used to retrieve the component
	typeKeyType := reflect.TypeOf(typeKey)

	c.RegisterFactory(typeKeyType, func(c *Container) any {
		return factory(c)
	})
}

// GetByTypeKey retrieves a component by a type key string
func GetByTypeKey[T any](c *Container, typeKey string) (T, error) {
	var zero T

	// Create a type that matches what was registered
	typeKeyType := reflect.TypeOf(typeKey)

	c.logger.Debug().
		Str("type_key", typeKey).
		Type("type_key_type", typeKeyType).
		Msg("GetByTypeKey called")

	component, err := c.Get(typeKeyType)
	if err != nil {
		c.logger.Debug().
			Str("type_key", typeKey).
			Err(err).
			Msg("GetByTypeKey failed")
		return zero, err
	}

	typedComponent, ok := component.(T)
	if !ok {
		c.logger.Debug().
			Str("type_key", typeKey).
			Type("expected_type", reflect.TypeOf((*T)(nil)).Elem()).
			Type("actual_type", reflect.TypeOf(component)).
			Msg("GetByTypeKey type mismatch")
		return zero, fmt.Errorf("component with key %s is not of expected type", typeKey)
	}

	c.logger.Debug().
		Str("type_key", typeKey).
		Type("type", reflect.TypeOf(typedComponent)).
		Msg("GetByTypeKey succeeded")
	return typedComponent, nil
}

// GetMultiple retrieves multiple components by their type keys
func GetMultiple[T any](c *Container, typeKeys []string) ([]T, error) {
	var results []T

	c.logger.Debug().
		Strs("type_keys", typeKeys).
		Type("component_type", reflect.TypeOf((*T)(nil)).Elem()).
		Msg("GetMultiple called")

	for _, key := range typeKeys {
		component, err := GetByTypeKey[T](c, key)
		if err != nil {
			c.logger.Debug().
				Str("failed_key", key).
				Strs("type_keys", typeKeys).
				Err(err).
				Msg("GetMultiple failed on key")
			return nil, err
		}

		results = append(results, component)
	}

	c.logger.Debug().
		Strs("type_keys", typeKeys).
		Int("results_count", len(results)).
		Msg("GetMultiple succeeded")
	return results, nil
}
