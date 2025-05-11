# VDM2-Bank Makefile

# Variables
BUILD_DIR = ./build
BINARY_NAME = vdm2-bank
TEST_DIR = ./tests

# Build
.PHONY: build
build:
	@echo "Building VDM2-Bank API..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/api

# Run
.PHONY: run
run:
	@echo "Running VDM2-Bank API..."
	@go run ./cmd/api/main.go

# Run with Docker Compose
.PHONY: up
up:
	@echo "Starting services with Docker Compose..."
	@docker-compose up --build -d

# Stop Docker Compose
.PHONY: down
down:
	@echo "Stopping Docker Compose services..."
	@docker-compose down

# Test API
.PHONY: test-api
test-api:
	@echo "Running API tests..."
	@$(TEST_DIR)/client/run_tests.sh

# Test all
.PHONY: test
test:
	@echo "Running all tests..."
	@go test ./internal/... -v

# Integration tests
.PHONY: test-integration
test-integration:
	@echo "Running integration tests..."
	@docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit

# Clean
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@docker-compose down -v
	@docker-compose -f docker-compose.test.yml down -v

# Full test cycle
.PHONY: test-all
test-all: test test-integration test-api

# Help
.PHONY: help
help:
	@echo "VDM2-Bank Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build           Build the API binary"
	@echo "  make run             Run the API locally"
	@echo "  make up              Start all services with Docker Compose"
	@echo "  make down            Stop all Docker Compose services"
	@echo "  make test            Run unit tests"
	@echo "  make test-api        Run API client tests"
	@echo "  make test-integration Run integration tests"
	@echo "  make test-all        Run all tests"
	@echo "  make clean           Clean up generated files and containers"
	@echo "  make help            Show this help message"

# Default target
.DEFAULT_GOAL := help