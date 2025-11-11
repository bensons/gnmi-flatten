.PHONY: help build clean test install run fmt vet all

# Binary name
BINARY_NAME=gnmi-flatten

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet
GOINSTALL=$(GOCMD) install

# Default target
all: build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -buildvcs=false -o $(BINARY_NAME) -v

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Install the binary to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	$(GOINSTALL)

# Run with sample file
run: build
	@echo "Running $(BINARY_NAME) with test-sample.json..."
	./$(BINARY_NAME) -file test-sample.json

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

# Run go vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...

# Show help
help:
	@echo "Available targets:"
	@echo "  make build    - Build the binary"
	@echo "  make clean    - Remove build artifacts"
	@echo "  make test     - Run tests"
	@echo "  make install  - Install binary to GOPATH/bin"
	@echo "  make run      - Build and run with sample.json"
	@echo "  make fmt      - Format code"
	@echo "  make vet      - Run go vet"
	@echo "  make all      - Build the binary (default)"
	@echo "  make help     - Show this help message"

