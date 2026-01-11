# VDM2-Bank

A Go-based backend service providing authentication, account management, and transaction processing for financial applications. This project follows a layered, modular architecture built with the Gin framework, featuring PostgreSQL for persistent storage, Redis for caching and rate limiting, and OAuth2 support.

## Features

- **User Authentication**: Email/password (bcrypt) and Google OAuth2
- **Account Management**: CRUD operations on user accounts, balance retrieval
- **Transactions**: Record movements, funds transfer with ACID guarantees
- **Caching & Rate Limiting**: Redis-backed session storage, balance cache, API rate limiting
- **Observability**: Structured logging (Zap), Prometheus metrics, Swagger documentation
- **Security**: JWT-based auth, input validation, secrets management
- **Deployment**: Docker, Kubernetes-ready, CI/CD pipeline with linting, testing, and security scans

## Getting Started

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 13+
- Redis 6+

### Installation

1. Clone the repository:

   ```bash
   git clone https://VDM2-BankBE.git
   cd VDM2-Bank
   ```

2. Copy and customize configuration:

   ```bash
   cp configs/config.yaml.example configs/config.yaml
   # Edit database URL, JWT secret, OAuth credentials
   ```

3. Build and run with Docker Compose:

   ```bash
   docker-compose up --build -d
   ```

4. Apply database migrations:

   ```bash
   docker-compose exec api ./service migrate up
   ```

5. The API will be available at `http://localhost:8080/api/v1`.
6. Swagger documentation will be available at `http://localhost:8080/swagger/index.html`.

## Project Structure

```plaintext
├── cmd/api                # Application entrypoint and docs
├── internal/
│   ├── config             # Configuration loader
│   ├── router             # Gin route definitions
│   ├── handler            # HTTP handlers
│   ├── service            # Business logic
│   ├── repository         # GORM database interactions
│   ├── model              # Go structs for DB models
│   ├── middleware         # Logging, auth, recovery
│   └── util               # Helpers (errors, pagination)
├── pkg/
│   ├── oauth              # Google OAuth client
│   └── cache              # Redis wrappers
├── configs                # YAML config files
├── migrations             # SQL migrations
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── README.md
```

## Configuration

All configuration values are stored in `configs/config.yaml`. Environment-specific overrides can be loaded via `.env` files.

## Running Tests

- **Unit Tests**:

  ```bash
  go test ./internal/...
  ```

- **Integration Tests**:

  ```bash
  docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit
  ```

## API Endpoints

Base URL: `/api/v1`

### Authentication
- `POST /auth/signup` - Register via email & bcrypt-hashed password
- `POST /auth/login` - Email/password login → issue JWT
- `GET /auth/google` - Redirect to Google OAuth consent
- `GET /auth/google/callback` - Handle OAuth callback

### Accounts
- `GET /accounts/balance` - Get account balance (DB + Redis cache)
- `GET /accounts/movements` - List transaction history
- `POST /accounts/movements` - Create a new movement

### Transfers
- `POST /transfers` - Funds transfer (wrapped in DB transaction)
- `GET /transfers` - List account transfers

Detailed API documentation is available via Swagger at `/swagger/index.html`.

## OpenAPI-first workflow (contract = reality)

This repository now includes a **multi-file OpenAPI 3.0.x** contract extracted from the existing Gin router and handlers.

- **Source of truth**: `api/src/openapi.yaml` (plus `api/src/paths/*` and `api/src/components/*`)
- **Bundled artifact (generated)**: `api/dist/openapi.bundle.yaml`
- **Generated Go server types + Gin interface (generated)**: `internal/generated/api.gen.go` (via `oapi-codegen`)
- **Generated client SDKs (generated)**:
  - TypeScript: `clients/ts`
  - Kotlin: `clients/kotlin`
  - Swift: `clients/swift`
  - Dart/Flutter: `clients/dart`

### How routes are discovered (baseline extraction)

Baseline OpenAPI was derived from:

- Gin routes: `internal/router/router.go`
- Handler request/response shapes: `internal/handler/*.go`
- Error envelope: `internal/util/errors.go`
- Auth middleware behavior: `internal/middleware/auth_middleware.go`

If behavior is unclear for a field/response, the spec contains **TODO** notes in `description` while keeping generation consistent.

### Regenerate server + SDKs

Prereqs:

- Go 1.21+ (CI uses Go 1.22)
- Node.js 20+ (for Redocly + openapi-generator wrapper)
- Java 17+ (required by OpenAPI Generator)

Run:

```bash
make generate
```

### Adding a new endpoint going forward (spec-first)

1. Edit OpenAPI under `api/src/` (add/modify path + schemas).
2. Run `make generate`.
3. Implement the generated `generated.ServerInterface` (see `cmd/server/main.go` for minimal wiring).
4. Keep CI green: the workflow expects generation to be deterministic and the working tree clean after `make generate`.

### Authentication (JWT + PASETO in contract)

- The OpenAPI contract documents two bearer schemes:
  - `BearerJWT` (implemented today)
  - `BearerPASETO` (**documented**, TODO to implement server-side validation)
- Current production auth enforcement is JWT-only via `internal/middleware/auth_middleware.go`.

## Contributing

1. Fork the repository
2. Create a new feature branch (`git checkout -b feature/YourFeature`)
3. Commit your changes (`git commit -m 'Add YourFeature'`)
4. Push to the branch (`git push origin feature/YourFeature`)
5. Open a Pull Request

## License

© 2025 VDMSquare