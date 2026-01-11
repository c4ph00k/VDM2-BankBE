# Testing (Unit Tests)

This backend uses **Gin** and is structured as:

- **Routing**: `internal/router/router.go` registers OpenAPI-generated routes from `internal/generated` via `generated.RegisterHandlersWithOptions(...)`.
- **Handlers (HTTP layer)**: `internal/handler/*` depend on `internal/service` interfaces.
- **Middleware**: `internal/middleware/*` enforces auth based on OpenAPI security markers (`generated.BearerJWTScopes`, `generated.BearerPASETOScopes`).
- **Services (business logic)**: `internal/service/*` depend on repository/store boundaries and are unit-tested without real DB/Redis/network.
- **Repositories (DB layer)**: `internal/repository/*` (GORM) are **not** exercised with real DB in unit tests.
- **Error envelope**: `internal/util/errors.go` controls error-to-HTTP mapping via `util.HandleError`.

## Test architecture

- **Gin handler tests** (endpoint-level): build a Gin engine using the generated router + middleware and hit endpoints via `httptest`.
  - Service dependencies are mocked with **GoMock**.
- **Service tests**: call service methods directly.
  - Repository/store dependencies are mocked with **GoMock**.
- **Middleware tests**: exercise auth middleware behavior (missing token, wrong scheme, invalid/expired token, valid JWT/PASETO).
- **Token tests**: verify current JWT + PASETO validation behavior used by `AuthService.VerifyToken`.

Test helpers live in `internal/testutil/`.

## Mock generation (GoMock + mockgen)

Mocks are generated with `mockgen` via `//go:generate` directives placed immediately above each interface:

- `internal/service/service.go` (handler→service boundary)
- `internal/repository/repository.go` (service→repository boundary)
- `internal/service/deps.go` (service→external boundaries: cache/OAuth/transaction runner)

Mocks live alongside their package:

- `internal/service/mocks/`
- `internal/repository/mocks/`

## How to run tests

From `VDM2-BankBE/`:

```bash
make test
```

Or:

```bash
go test ./...
```

## How to (re)generate mocks

From `VDM2-BankBE/`:

```bash
make gen-mocks
```

Or:

```bash
go generate ./...
```

## Adding a new unit test

- **Handler test**: add a new table-driven test in `internal/handler/*_test.go` and mock the required service(s).
- **Service test**: add a new table-driven test in `internal/service/*_test.go` and mock required repository/store interfaces.
- Use `internal/testutil` helpers for building requests, decoding responses, and generating deterministic tokens.

## CI expectations

CI runs:

- `go generate ./...` and fails if it produces a diff (mocks must be up to date)
- `go test ./...`

