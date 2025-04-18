package media

import (
	"context"
	"fmt"
	"reflect"

	"suasor/client/media/types"
)

var GlobalMediaRegistry = NewClientItemRegistry()

func init() {
	// Register all the media factories for Emby
	RegisterFactory[*EmbyClient, *embyclient.BaseItemDto, *types.Album](
		GlobalMediaRegistry,
		func(client *EmbyClient, ctx context.Context, item *embyclient.BaseItemDto) (*types.Album, error) {
			return client.albumFactory(ctx, item)
		},
	)

	RegisterFactory[*EmbyClient, *embyclient.BaseItemDto, *types.Movie](
		GlobalMediaRegistry,
		func(client *EmbyClient, ctx context.Context, item *embyclient.BaseItemDto) (*types.Movie, error) {
			return client.movieFactory(ctx, item)
		},
	)

	// Register other factories for other media types...
}

// MediaFactory is a generic factory function type
type MediaFactory[C ClientMedia, I any, O types.MediaData] func(client C, ctx context.Context, item I) (O, error)

// MediaClientRegistry manages factories for different clients and media types
type ClientItemRegistry struct {
	factories map[reflect.Type]map[reflect.Type]map[reflect.Type]any
}

// NewMediaClientRegistry creates a new empty registry
func NewClientItemRegistry() *ClientItemRegistry {
	return &ClientItemRegistry{
		factories: make(map[reflect.Type]map[reflect.Type]map[reflect.Type]any),
	}
}

// RegisterFactory adds a factory function to the registry
func RegisterFactory[C ClientMedia, I any, O types.MediaData](
	registry *ClientItemRegistry,
	factory MediaFactory[C, I, O],
) {
	var (
		clientType = reflect.TypeOf((*C)(nil)).Elem()
		inputType  = reflect.TypeOf((*I)(nil)).Elem()
		outputType = reflect.TypeOf((*O)(nil)).Elem()
	)

	// Initialize nested maps if needed
	if _, ok := registry.factories[clientType]; !ok {
		registry.factories[clientType] = make(map[reflect.Type]map[reflect.Type]any)
	}
	if _, ok := registry.factories[clientType][inputType]; !ok {
		registry.factories[clientType][inputType] = make(map[reflect.Type]any)
	}

	// Store the factory
	registry.factories[clientType][inputType][outputType] = factory
}

// GetFactory retrieves a factory function for specific types
func GetFactory[C ClientMedia, I any, O types.MediaData](
	registry *ClientItemRegistry,
) (MediaFactory[C, I, O], error) {
	var (
		clientType = reflect.TypeOf((*C)(nil)).Elem()
		inputType  = reflect.TypeOf((*I)(nil)).Elem()
		outputType = reflect.TypeOf((*O)(nil)).Elem()
		zero       O
	)

	// Look up the factory
	if clientFactories, ok := registry.factories[clientType]; ok {
		if inputFactories, ok := clientFactories[inputType]; ok {
			if factory, ok := inputFactories[outputType]; ok {
				return factory.(MediaFactory[C, I, O]), nil
			}
		}
	}

	return func(C, context.Context, I) (O, error) {
		return zero, fmt.Errorf("no factory registered for client type %v, input type %v, output type %v",
			clientType, inputType, outputType)
	}, fmt.Errorf("factory not found")
}

// ConvertTo converts item to the desired output type using the appropriate factory
func ConvertTo[C ClientMedia, I any, O types.MediaData](
	registry *ClientItemRegistry,
	client C,
	ctx context.Context,
	item I,
) (O, error) {
	factory, err := GetFactory[C, I, O](registry)
	if err != nil {
		var zero O
		return zero, err
	}
	return factory(client, ctx, item)
}

// In a separate file:
// album := GetMediaItem[*types.Album](registry, embyClient, ctx, itemDto)
