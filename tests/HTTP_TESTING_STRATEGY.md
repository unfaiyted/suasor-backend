# HTTP Testing Strategy for Suasor

## Overview

This document outlines the testing strategy for Suasor using `.http` files with scripts. The approach allows for:

1. **Modular Testing**: Break tests into logical units
2. **Test Reusability**: Import and reuse test components 
3. **Test Orchestration**: Run tests in sequence with dependencies
4. **Validation**: Use scripts to validate responses
5. **Test Flow Control**: Skip or replay tests based on conditions
6. **State Management**: Share data between requests

## File Organization

Tests are organized into the following categories:

1. **Core Tests**
   - `_core/auth.http` - Authentication-related tests
   - `_core/users.http` - User management tests
   - `_core/config.http` - Configuration tests
   - `_core/health.http` - Health checks

2. **Feature Tests**
   - `client/*.http` - Client-related tests
   - `media/*.http` - Media-related tests 
   - `recommendation/*.http` - Recommendation tests
   - `ai/*.http` - AI-related tests
   - `jobs/*.http` - Job-related tests

3. **Integration Tests**
   - `integration/*.http` - End-to-end tests across features

4. **Test Runners**
   - `_runners/all.http` - Run all tests
   - `_runners/core.http` - Run core tests
   - `_runners/media.http` - Run media tests
   - `_runners/quick.http` - Run a subset of tests for quick validation

## Test Structure

Each test file should follow this structure:

```http
### Test File Description

# Global variables for this file
@baseVariable = value

### Test Name 1
# @name testName1
METHOD {{baseUrl}}/path
Content-Type: application/json
Authorization: Bearer {{tokenVariable}}

{
  "property": "value"
}

> {%
  // Validation script
  client.test("Status should be 200", function() {
    client.assert(response.status === 200, "Response status is not 200");
  });
  
  // Save response data for later tests
  client.global.set("savedVariable", response.body.data.id);
%}

### Test Name 2
# @name testName2
// Test that depends on testName1
```

## Script Usage

### Pre-request Scripts

Use pre-request scripts (prefixed with `<`) to:
- Set up test data
- Skip tests based on conditions
- Modify request parameters

Example:
```http
< {%
  if (!client.global.get("authToken")) {
    request.skip();
  }
%}
```

### Post-request Scripts

Use post-request scripts (prefixed with `>`) to:
- Validate responses
- Store data for subsequent requests
- Replay requests with modified parameters
- Log test information

Example:
```http
> {%
  client.test("Response status is 200", function() {
    client.assert(response.status === 200);
  });
  
  client.global.set("userId", response.body.data.id);
  
  client.log("User created with ID: " + response.body.data.id);
%}
```

## Variable Management

### Environment Variables

Store sensitive data in environment files:
- `http-client.env.json` - Define environment variables
- `http-client.private.env.json` - Store sensitive data (not committed to version control)

### Test Variables

- `client.global.set/get` - Variables available across all requests
- `client.test.set/get` - Variables available within the current test file
- `request.variables.set/get` - Variables available only for the current request

## Test Automation

1. **Test Dependencies**: Use the `run` command to execute tests in order
2. **Conditional Testing**: Use pre-request scripts to skip tests based on conditions
3. **Retry Logic**: Use `request.replay()` to retry requests with modified parameters
4. **Validation**: Use post-request scripts to validate responses
5. **Error Handling**: Use try/catch blocks for robust error handling

## Common Test Patterns

### Authentication Flow

```http
### Login
# @name login
POST {{baseUrl}}/auth/login
Content-Type: application/json

{
  "email": "{{TEST_USER}}",
  "password": "{{TEST_PASSWORD}}"
}

> {%
  client.global.set("authToken", response.body.data.accessToken);
%}

### Use Authentication
GET {{baseUrl}}/some/protected/endpoint
Authorization: Bearer {{authToken}}
```

### Sequential Tests

```http
# Run authentication first
run #login

# Then run other tests that require authentication
run #createResource
run #getResource
run #updateResource
run #deleteResource
```

### Test Data Validation

```http
> {%
  client.test("Response has valid structure", function() {
    client.assert(response.body.hasOwnProperty("data"), "Missing data property");
    client.assert(Array.isArray(response.body.data), "Data is not an array");
    client.assert(response.body.data.length > 0, "Data array is empty");
  });
%}
```

## Best Practices

1. **Naming Convention**: Use descriptive names for tests
2. **Test Independence**: Design tests to be run independently when possible
3. **Clean Up**: Clean up test data after tests complete
4. **Idempotence**: Design tests to be repeatable without side effects
5. **Readability**: Add comments to explain complex test logic
6. **Reusability**: Create reusable test components for common operations
7. **Security**: Don't hardcode sensitive data - use environment variables
8. **Documentation**: Document expected behaviors and edge cases