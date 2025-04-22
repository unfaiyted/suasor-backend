# Migration Plan: Moving to Wire

This document outlines a step-by-step plan to migrate the current manual DI system to Google Wire.

## Wire v0.6.0 Limitations and Workarounds

Wire v0.6.0 has limited support for Go generics, particularly with:

1. **Multi-parameter generic types**: Types like `ClientMediaItemHandler[T, U]` where multiple type parameters are used
2. **AST handling**: Wire's AST processing cannot handle some Go 1.18+ syntax constructs

### Our Workaround Approach

We've implemented a hybrid approach to work around these limitations:

1. **Wire for Non-Generic Components**: Use Wire to generate dependency wiring for non-generic components
   - AuthHandler, ConfigHandler, HealthHandler, UserHandler
   - These are generated with individual Wire injector functions

2. **Manual Wiring for Generic Components**: Manually implement wiring for generic components
   - Media handlers, client handlers, specialized handlers
   - This is done in the `InitializeAllHandlers` function

3. **Hybrid Container Structure**: Combine Wire-generated and manually wired components
   - Wire generates handlers for non-generic components
   - Manual code initializes and wires together generic components
   - Both are combined in the final application structure

## Phase 1: Proof of Concept

1. **Start with Media Lists**: Implement Wire for just the media lists functionality
   - Create provider functions for all related components
   - Test that the generated code works correctly
   - Compare with the existing implementation

2. **Integrate with Existing System**: Make the Wire-generated code work alongside the existing system
   - Create a bridge between Wire-generated components and manually registered ones
   - Test that they work together properly

## Phase 2: Gradual Migration

3. **Migrate Core Services**: Move the core service layer to Wire
   - Create provider functions for all core services
   - Update injector functions to include these providers
   - Test that everything still works

4. **Migrate User Services**: Move the user service layer to Wire
   - Create provider functions for all user services
   - Update injector functions to include these providers
   - Test again

5. **Migrate Client Services**: Move the client service layer to Wire
   - Create provider functions for all client services
   - Update injector functions to include these providers
   - Test with multiple client types

## Phase 3: Complete Migration

6. **Migrate Repositories**: Move all repositories to Wire
   - Create provider functions for all repositories
   - Update injector functions as needed
   - Test database operations

7. **Migrate Handlers**: Move all handlers to Wire
   - Create provider functions for all handlers
   - Update injector functions as needed
   - Test API endpoints

8. **Migrate Main Application**: Create a top-level injector function
   - Create an injector that builds the entire application
   - Replace the current initialization code
   - Test the entire application

## Phase 4: Clean-up and Optimization

9. **Remove Old DI Code**: Once everything is migrated, remove the old code
   - Remove old registrations
   - Clean up any bridge code
   - Test the application

10. **Optimize Provider Structure**: Refine the provider organization
    - Group related providers into provider sets
    - Optimize dependencies between providers
    - Look for places to simplify the graph

## Phase 5: Documentation and Training

11. **Update Documentation**: Document the new DI system
    - Update code comments
    - Update README files
    - Create examples

12. **Train Team Members**: Ensure everyone understands the new system
    - Explain the benefits of Wire
    - Show how to add new components
    - Show how to troubleshoot issues

## Best Practices During Migration

1. **Keep Changes Small**: Migrate one component at a time
2. **Maintain Test Coverage**: Keep running tests after each change
3. **Use Feature Flags**: If possible, use feature flags to switch between old and new DI systems
4. **Document As You Go**: Keep track of changes and decisions
5. **Regularly Check for Regressions**: Make sure existing functionality still works

## Future Plan for Full Wire Integration

As Wire continues to evolve and support for generics improves, our plan is:

1. **Monitor Wire Development**: Keep an eye on new Wire releases that add better support for generics
2. **Gradual Replacement**: When a Wire version fully supporting our generic types is available, gradually replace our manual wiring with Wire-generated code
3. **Version Constraints**: Explicitly specify Wire version constraints in our module to ensure consistent behavior

This hybrid approach allows us to benefit from Wire's DI capabilities for simpler components while still maintaining a functional and performant system for complex generic components.