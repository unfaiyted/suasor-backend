.PHONY: swag build run docker-build docker-run test pretty-test claude-example movie-recommendations

# Variables
DOCKER_IMAGE := suasor 
GO_VERSION := 1.23
ALPINE_VERSION := 3.19

# Local development commands
swag:
	swag init

build: swag
	CGO_ENABLED=0 go build -o main main.go

run: swag
	go run main.go

test:
	go test ./... -v

pretty-test:
	gotestsome ./...

# AI client examples
claude-example:
	@echo "Running Claude AI client example..."
	@if [ -z "$$CLAUDE_API_KEY" ]; then \
		echo "Error: CLAUDE_API_KEY environment variable is not set"; \
		echo "Usage: CLAUDE_API_KEY=your-api-key make claude-example"; \
		exit 1; \
	fi
	go run examples/claude_client_example.go

movie-recommendations:
	@echo "Running Movie Recommendations example..."
	@echo "Using environment variables from .env file"
	go run examples/movie_recommendations.go
# Docker commands
docker-build: swag
	docker build -t $(DOCKER_IMAGE) .

docker-run: docker-build
	docker run -p 8080:8080 $(DOCKER_IMAGE)
