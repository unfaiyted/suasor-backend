# Suasor Backend

Suasor is a comprehensive media management and AI integration platform that provides a unified interface for interacting with various media clients and AI services.

## Features

- **Media Client Integration**: Connect to Plex, Jellyfin, Emby, and Subsonic media servers
- **Automation Integration**: Support for Radarr, Sonarr, and Lidarr
- **AI Integration**: Connect to Claude AI, OpenAI, and Ollama for media recommendations and analysis
- **REST API**: Comprehensive API for all operations with Swagger documentation
- **User Management**: Authentication, roles, and user-specific configurations

## Prerequisites

- Go 1.20 or higher
- PostgreSQL database
- Docker (optional)
- Make (optional)

## Project Structure

```
├── app/          # Application components and dependency injection
├── client/       # Client implementations for external services
│   ├── ai/       # AI service clients (Claude, OpenAI, etc.)
│   ├── automation/ # Automation clients (Radarr, Sonarr, etc.)
│   ├── media/    # Media server clients (Plex, Jellyfin, etc.)
├── config/       # Configuration files
├── database/     # Database connection and models
├── docs/         # Swagger documentation
├── handlers/     # HTTP request handlers
├── repository/   # Data access layer
├── router/       # HTTP routing and middleware
├── services/     # Business logic layer
├── types/        # Type definitions
│   ├── models/   # Data models
│   ├── requests/ # Request types
│   ├── responses/ # Response types
├── utils/        # Utility functions
```

## Getting Started

### Local Development

1. Clone the repository

```bash
git clone https://github.com/unfaiyted/suasor.git
cd suasor/backend
```

2. Install dependencies

```bash
go mod download
```

3. Set up environment variables or config file

The application uses a configuration file located at `config/app.config.json`.

4. Run the application

```bash
make run
```

The API will be available at `http://localhost:8080`

### Using Docker

```bash
# Build the Docker image
make docker-build

# Run the container
make docker-run
```

## API Documentation

The API documentation is available through Swagger UI when the application is running:

```
http://localhost:8080/swagger/index.html
```

### Key Endpoints

- `GET /api/v1/health` - Health check endpoint
- `GET /api/v1/auth/login` - User authentication
- `GET /api/v1/clients/{clientType}` - Get clients by type
- `GET /api/v1/swagger/index.html` - API documentation (Swagger UI)

## Development

### Generate Swagger Documentation

```bash
make swag
```

This runs `swag init --exclude ./internal/**` to generate Swagger documentation while excluding the internal API client files.

### Adding New Endpoints

1. Define the route in `router/router.go` or respective router file
2. Create handler in `handlers/`
3. Update Swagger documentation using comments
4. Generate new Swagger docs with `make swag`

### Running Tests

```bash
# Run all tests
make test

# Run pretty tests (requires gotestsome)
make pretty-test

# Run integration tests (requires configured clients)
INTEGRATION=true make test
```

### Running Examples

```bash
# Run Claude AI client example
CLAUDE_API_KEY=your-api-key make claude-example

# Run movie recommendations example
make movie-recommendations
```

## Architecture

The application follows a clean architecture with:

- HTTP handlers in the `handlers` package
- Business logic in the `services` package
- Data access in the `repository` package
- External clients in the `client` package

Dependency injection is used to ensure loose coupling and testability.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.