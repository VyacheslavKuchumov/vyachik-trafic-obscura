.PHONY: all build-client build-server run-client run-server clean

# Default target
all: build-client build-server

# Build targets
build-client:
	@echo "Building client..."
	@go build -o bin/client ./cmd/client
	@echo "Client built: bin/client"

build-server:
	@echo "Building server..."
	@go build -o bin/server ./cmd/server
	@echo "Server built: bin/server"

# Run targets (with auto-build)
run-client: build-client
	@echo "Running client..."
	./bin/client

run-server: build-server
	@echo "Running server..."
	./bin/server

# Development (hot reload) - optional
dev-client:
	@echo "Starting client in dev mode..."
	@go run ./cmd/client

dev-server:
	@echo "Starting server in dev mode..."
	@go run ./cmd/server

# Clean up
clean:
	@echo "Cleaning up..."
	@rm -rf bin/
	@echo "Cleaned bin directory"