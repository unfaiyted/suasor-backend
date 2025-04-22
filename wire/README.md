# Wire Dependency Injection for Suasor

This directory contains Google Wire-based dependency injection for the Suasor backend. Wire is a code generation tool that automates dependency injection.

## Benefits

1. **Type Safety**: All dependencies are checked at compile time
2. **Explicit Dependencies**: The generated code shows exactly how dependencies are constructed
3. **No Reflection**: Everything is explicit Go code, no runtime reflection magic
4. **Reduced Boilerplate**: No need to write repetitive registration code

## Getting Started

1. Install Wire:
```
go install github.com/google/wire/cmd/wire@latest
```

2. Run Wire to generate the dependency injection code:
```
cd wire
go run github.com/google/wire/cmd/wire
```

This will produce a `medialists_gen.go` file with the generated code.

## How to Use Wire

1. **Define Provider Functions**: Create functions that return your components

2. **Create Injector Functions**: These define what you want to create and what dependencies are needed

3. **Run Wire**: Generate the implementation of your injector functions

4. **Use the Generated Code**: Import and use the generated functions in your application

## Example

The file `example.go` demonstrates how to use the Wire-generated code in your application. It shows how to:

1. Call the injector function to get fully initialized handlers
2. Use those handlers to set up your routes
3. Let Wire handle all the dependency complexity

## Expanding the System

To add more components to the Wire system:

1. Add provider functions for the new components
2. Update the injector functions to include the new providers
3. Run Wire to regenerate the code

## Best Practices

1. Keep provider functions simple and focused
2. Group related providers in provider sets
3. Use interfaces for dependencies to allow easy mocking in tests
4. Keep Wire-specific code separate from business logic

## Migrating Existing Code

To migrate existing code to Wire:

1. Start with a small, well-defined part of your application (like we did with MediaLists)
2. Create provider functions for each component
3. Replace manual registration with Wire-generated initialization
4. Gradually expand to cover more of your application