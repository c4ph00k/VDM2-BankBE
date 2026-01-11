.PHONY: test gen-mocks

test:
	go test ./...

gen-mocks:
	go generate ./...

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

# Generate Swagger docs
.PHONY: swagger
swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/api/main.go -o cmd/api/docs --parseDependency --parseInternal --parseDepth 3

# ------------------------------
# OpenAPI-first pipeline (OAS 3.0)
# ------------------------------

OPENAPI_SRC := api/src/openapi.yaml
OPENAPI_BUNDLE := api/dist/openapi.bundle.yaml
OAPI_CODEGEN_VERSION := v2.4.0

.PHONY: api-lint
api-lint:
	@echo "Linting OpenAPI..."
	@npm ci
	@npx --no-install redocly lint --config .redocly.yaml $(OPENAPI_SRC)

.PHONY: api-bundle
api-bundle:
	@echo "Bundling OpenAPI -> $(OPENAPI_BUNDLE)..."
	@npm ci
	@node -e "require('fs').mkdirSync('api/dist',{recursive:true})"
	@npx --no-install redocly bundle --config .redocly.yaml $(OPENAPI_SRC) -o $(OPENAPI_BUNDLE)

.PHONY: gen-server
gen-server: api-bundle
	@echo "Generating Go server types + Gin interface..."
	@node -e "require('fs').mkdirSync('internal/generated',{recursive:true})"
	@go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION) \
		-config api/oapi-codegen.yaml \
		$(OPENAPI_BUNDLE) > internal/generated/api.gen.go

.PHONY: gen-clients
gen-clients: api-bundle
	@echo "Generating client SDKs (TS/Kotlin/Swift/Dart)..."
	@npm ci
	@npx --no-install openapi-generator-cli generate -c clients/ts/openapi-generator-config.yaml
	@npx --no-install openapi-generator-cli generate -c clients/kotlin/openapi-generator-config.yaml
	@npx --no-install openapi-generator-cli generate -c clients/swift/openapi-generator-config.yaml
	@npx --no-install openapi-generator-cli generate -c clients/dart/openapi-generator-config.yaml

.PHONY: generate
generate: api-lint api-bundle gen-server gen-clients
	@echo "Generation complete."

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
	@echo "  make swagger         Generate Swagger API documentation"
	@echo "  make clean           Clean up generated files and containers"
	@echo "  make help            Show this help message"

# Default target
.DEFAULT_GOAL := help