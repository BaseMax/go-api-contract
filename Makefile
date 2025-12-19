.PHONY: build test clean install run-example help

# Build the binary
build:
	@echo "Building apicontract..."
	@go build -o apicontract ./cmd/apicontract
	@echo "✓ Build complete: ./apicontract"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run the test script
test-all: build
	@echo "Running full test suite..."
	@./test.sh

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f apicontract
	@echo "✓ Clean complete"

# Install the binary to GOPATH/bin
install:
	@echo "Installing apicontract..."
	@go install ./cmd/apicontract
	@echo "✓ Install complete"

# Run example
run-example: build
	@echo "Running example..."
	@./apicontract -spec examples/openapi-local.yaml \
		-endpoint http://localhost:8080/posts/1 \
		-method GET \
		-verbose

# Start test server
start-server:
	@echo "Starting test server on http://localhost:8080..."
	@go run examples/test-server.go

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✓ Format complete"

# Run linter
lint:
	@echo "Running linter..."
	@go vet ./...
	@echo "✓ Lint complete"

# Get dependencies
deps:
	@echo "Getting dependencies..."
	@go mod download
	@go mod tidy
	@echo "✓ Dependencies updated"

# Help
help:
	@echo "Available targets:"
	@echo "  build        - Build the binary"
	@echo "  test         - Run Go tests"
	@echo "  test-all     - Run full test suite"
	@echo "  clean        - Remove build artifacts"
	@echo "  install      - Install to GOPATH/bin"
	@echo "  run-example  - Run an example validation"
	@echo "  start-server - Start the test server"
	@echo "  fmt          - Format code"
	@echo "  lint         - Run linter"
	@echo "  deps         - Update dependencies"
	@echo "  help         - Show this help message"
