# Backend Implementation in Golang

This document describes the end-to-end design, implementation, and operational considerations for the Go-based backend of our service.

---

## 1. Folder Structure

```
cmd/
└── api/                       # Entry point: main.go

internal/
├── config/                    # Load env/YAML (Viper)
├── router/                    # Gin router setup
├── handler/                   # HTTP handlers → services
├── service/                   # Business logic
├── repository/                # DB access via GORM
├── model/                     # GORM models with `gorm` & `json` tags
├── middleware/                # JWT auth, logging, recovery
└── util/                      # Error helpers, pagination

pkg/
└── oauth/                     # Google OAuth2 client wrapper
└── cache/                     # Redis client wrappers

configs/
└── config.yaml                # Configuration defaults per environment

go.mod                         # Module definitions
go.sum                         # Dependency versions
```

---

## 2. Component Responsibilities

- **`model/`**  
  GORM-tagged Go structs reflecting tables (e.g. `User`, `Account`, `Movement`, etc.).

- **`repository/`**  
  Encapsulates all database interactions using `*gorm.DB`, including transactions.

- **`service/`**  
  Orchestrates business logic, enforces rules, composes multiple repository calls atomically.

- **`handler/`**  
  HTTP handlers using Gin that parse requests, call services, and write JSON responses.

- **`middleware/`**  
  Cross-cutting concerns: JWT validation, request logging (Zap/Logrus), panic recovery implemented as Gin middleware.

- **`router/`**  
  Wire up routes to handlers using Gin router and inject middleware.

- **`config/`**  
  Load and validate environment-specific settings (Viper).

- **`pkg/oauth` & `pkg/cache`**  
  Reusable packages for Google OAuth2 flows and Redis operations respectively.

---

## 3. Dependency Management

```go
module github.com/yourorg/yourapp

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1     // HTTP framework
    github.com/joho/godotenv v1.5.1
    github.com/pkg/errors v0.9.1         // error wrapping
    gorm.io/driver/postgres v1.4.0
    gorm.io/gorm v1.27.0
    github.com/golang-jwt/jwt/v5 v5.0.0   // JWT handling
    github.com/go-redis/redis/v8 v8.11.5
    golang.org/x/crypto/bcrypt v0.13.0
    golang.org/x/oauth2 v0.11.0
    golang.org/x/oauth2/google v0.11.0
    go.uber.org/zap v1.24.0              // structured logging
    github.com/prometheus/client_golang v1.15.0 // metrics
    github.com/gin-contrib/cors v1.4.0    // CORS middleware
    github.com/gin-contrib/zap v0.1.0     // Zap logger for Gin
)
```

- **Config**: Viper + optional `.env` via godotenv for local overrides.  
- **DI**: Initialize `*gorm.DB`, `*redis.Client`, and inject into repositories/services in `cmd/api/main.go`.

---

## 4. API Endpoints

| Method | Path                          | Handler                       | Description                                |
|--------|-------------------------------|-------------------------------|--------------------------------------------|
| POST   | `/auth/signup`                | `AuthHandler.SignUp`          | Register via email & bcrypt-hashed pass    |
| POST   | `/auth/login`                 | `AuthHandler.Login`           | Email/password login → issue JWT           |
| GET    | `/auth/google`                | `AuthHandler.GoogleAuth`      | Redirect to Google OAuth consent           |
| GET    | `/auth/google/callback`       | `AuthHandler.GoogleCallback`  | Handle OAuth callback, upsert user         |
| GET    | `/account/balance`            | `AccountHandler.Balance`      | Get account balance (DB + Redis cache)     |
| GET    | `/account/movements`          | `MovementHandler.List`        | List transaction history                   |
| POST   | `/transfer`                   | `TransferHandler.Transfer`    | Funds transfer (wrapped in DB transaction) |

---

## 5. API Call Examples

### 5.1 POST `/auth/signup`
**Request**:
```json
{
  "email": "user@example.com",
  "password": "secureP@ssw0rd",
  "username": "jdoe",
  "first_name": "John",
  "last_name": "Doe",
  "fiscal_code": "RSSMRA85T10A562S"
}
```
**Response** (`201 Created`):
```json
{
  "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "email": "user@example.com",
  "username": "jdoe",
  "first_name": "John",
  "last_name": "Doe",
  "fiscal_code": "RSSMRA85T10A562S",
  "created_at": "2025-05-01T12:34:56Z"
}
```

### 5.2 POST `/auth/login`
**Request**:
```json
{
  "email": "user@example.com",
  "password": "secureP@ssw0rd"
}
```
**Response** (`200 OK`):
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600
}
```

*... remaining call examples unchanged ...*

---

## 6. Testing Strategy

- **Unit Tests**  
  - Repositories: in-memory SQLite  
  - Services: mock repos via interfaces  
  - Handlers: `httptest` + table-driven scenarios

- **Integration Tests**  
  - Docker Compose: Postgres + Redis  
  - End-to-end API tests  
  - OAuth mock server
---


***********************************************************
## 7. Layered Architecture Diagram

```plain
┌──────────────┐
│ cmd/api/main │
└───────┬──────┘
        ↓
┌───────────────┐       ┌─────────────┐
│ config/init   │──────▶│ middleware  │
└───────────────┘       └─────────────┘
        ↓                      │
┌───────────────┐              ↓
│ router        │──────────► handlers
└───────────────┘              ↓
        ↓                      ↓
┌───────────────┐──────────► services
│ repositories  │              ↓
└───────────────┘              ↓
        ↓                      ↓
┌───────────────┐──────────► Postgres / Redis
│ pkg/oauth     │
│ pkg/cache     │
└───────────────┘
```

---

## 8. Robust Error Handling

- **Custom Error Type**  
  ```go
  type APIError struct {
      Code    int    `json:"code"`
      Message string `json:"message"`
  }
  func (e APIError) Error() string { return e.Message }
  ```

- **Middleware**  
  Catch panics and return:
  ```json
  {
    "error": { "code": 500, "message": "internal server error" }
  }
  ```

- **Wrapping**  
  Use `pkg/errors` for stack traces:  
  ```go
  return errors.Wrap(err, "failed to create user")
  ```

---

## 9. Structured Logging

- Adopt Zap or Logrus for JSON-formatted, leveled logs.  
- Include fields: `request_id`, `user_id`, `latency_ms`.  
- Example:
  ```go
  logger.Info("login attempt",
      zap.String("email", req.Email),
      zap.String("request_id", ctx.Value("reqID").(string)),
  )
  ```

---

## 10. Configuration Management

- Centralize via Viper; support `dev`, `staging`, `prod` profiles.  
- Validate required vars on startup; exit if missing critical keys (DB URL, JWT secret).

---

## 11. Database Migrations

- Use `golang-migrate/migrate` with SQL-based migrations.  
- Embed migrations in binary via Go’s `//go:embed` or `go-bindata`.

---

## 12. Security & Hardening

- **Rate Limiting**:  
  Redis-backed counters or `golang.org/x/time/rate`.  
- **Input Validation**:  
  `go-playground/validator` on request DTOs.  
- **Secrets Management**:  
  Integrate Vault or AWS Secrets Manager; avoid plaintext.

---

## 13. Observability

- **Metrics**:  
  Prometheus client for counters (requests, errors) and histograms (DB latency).  
- **Tracing**:  
  OpenTelemetry to trace across layers.  
- **Endpoints**:  
  Expose `/metrics` and `/debug/pprof` (behind auth in non-prod).

---

## 14. Deployment & CI/CD

- **Dockerfile**:  
  ```Dockerfile
  FROM golang:1.21-alpine AS builder
  WORKDIR /app
  COPY . .
  RUN go mod download && \
      CGO_ENABLED=0 go build -o service ./cmd/api

  FROM scratch
  COPY --from=builder /app/service /service
  ENTRYPOINT ["/service"]
  ```

- **Kubernetes**:  
  - Readiness/liveness probes  
  - Resource requests/limits  
  - ConfigMaps & Secrets

- **CI Pipeline**:  
  1. `go fmt`, `go vet`, `golangci-lint`  
  2. Unit & integration tests  
  3. Build & push Docker image  
  4. Security scans (Trivy, Snyk)

---

## 15. Performance & Scaling

- **DB Pooling**:  
  ```go
  db.DB().SetMaxOpenConns(50)
  db.DB().SetConnMaxLifetime(5 * time.Minute)
  ```

- **Caching**:  
  Redis for hot reads (e.g. balances).  
- **Profiling**:  
  Periodic `pprof` CPU & memory snapshots.

---

## 16. Code Organization & Style

- Follow [Effective Go](https://go.dev/doc/effective_go.html).  
- Keep handlers thin; business logic in `service/`.  
- Table-driven tests for handler & repository functions.

---


## 17. Data Models

### 17.1 PostgreSQL Schema & Go Structs

```sql
-- users table
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email TEXT NOT NULL UNIQUE,
  username TEXT NOT NULL UNIQUE,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  fiscal_code TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- accounts table
CREATE TABLE accounts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id),
  balance NUMERIC(18,2) NOT NULL DEFAULT 0,
  currency TEXT NOT NULL DEFAULT 'EUR',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- movements table
CREATE TABLE movements (
  id BIGSERIAL PRIMARY KEY,
  account_id UUID NOT NULL REFERENCES accounts(id),
  amount NUMERIC(18,2) NOT NULL,
  type TEXT NOT NULL CHECK (type IN ('credit','debit')),
  description TEXT,
  occurred_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- oauth_tokens table
CREATE TABLE oauth_tokens (
  user_id UUID PRIMARY KEY REFERENCES users(id),
  provider TEXT NOT NULL,
  access_token TEXT NOT NULL,
  refresh_token TEXT,
  expires_at TIMESTAMPTZ
);

-- transfers table
CREATE TABLE transfers (
  id BIGSERIAL PRIMARY KEY,
  from_account UUID NOT NULL REFERENCES accounts(id),
  to_account UUID NOT NULL REFERENCES accounts(id),
  amount NUMERIC(18,2) NOT NULL,
  status TEXT NOT NULL CHECK (status IN ('pending','completed','failed')),
  initiated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  completed_at TIMESTAMPTZ
);
```

```go
// internal/model/models.go

type User struct {
    ID          uuid.UUID ` + "`gorm:"type:uuid;primaryKey" json:"id"`" + `
    Email       string    ` + "`gorm:"uniqueIndex;not null" json:"email"`" + `
    Username    string    ` + "`gorm:"uniqueIndex;not null" json:"username"`" + `
    FirstName   string    ` + "`gorm:"not null" json:"first_name"`" + `
    LastName    string    ` + "`gorm:"not null" json:"last_name"`" + `
    FiscalCode  string    ` + "`gorm:"uniqueIndex;not null" json:"fiscal_code"`" + `
    PasswordHash string   ` + "`gorm:"not null" json:"-"`" + `
    CreatedAt   time.Time ` + "`json:"created_at"`" + `
    UpdatedAt   time.Time ` + "`json:"updated_at"`" + `
}

type Account struct { ... }
// (other models unchanged)
```

### 17.2 Redis Key Patterns & Data Structures

| Use Case                | Key Pattern                          | Type         | TTL         | Description                          |
|-------------------------|--------------------------------------|--------------|-------------|--------------------------------------|
| Session store           | `sess:{session_id}`                  | JSON string  | 24h         | Payload (user_id, expiry)            |
| Rate limiting           | `rl:{user_id}:{route}`               | Integer      | 1h          | Request count per user+route         |
| Balance cache           | `acct:balance:{account_id}`          | Decimal str  | 1m          | Hot cache for balance reads          |
| OTP / MFA code          | `otp:{user_id}:{purpose}`            | String       | 5m          | One-time verification codes          |
| OAuth state             | `oauth:state:{state_token}`          | String       | 10m         | CSRF protection during OAuth flows   |
| Feature flags           | `flags:{user_id}`                    | JSON map     | —           | Per-user feature toggles             |

```go
// pkg/cache/redis.go

type RedisClient struct {
    client *redis.Client
}

func (r *RedisClient) SetSession(ctx context.Context, sessID string, payload interface{}) error {
    data, _ := json.Marshal(payload)
    return r.client.Set(ctx, "sess:"+sessID, data, 24*time.Hour).Err()
}

func (r *RedisClient) IncrRateLimit(ctx context.Context, userID, route string, window time.Duration) (int64, error) {
    key := fmt.Sprintf("rl:%s:%s", userID, route)
    count, err := r.client.Incr(ctx, key).Result()
    if err != nil {
        return 0, err
    }
    if count == 1 {
        r.client.Expire(ctx, key, window)
    }
    return count, nil
}

func (r *RedisClient) GetBalanceCache(ctx context.Context, accountID uuid.UUID) (decimal.Decimal, error) {
    val, err := r.client.Get(ctx, "acct:balance:"+accountID.String()).Result()
    if err != nil {
        return decimal.Zero, err
    }
    return decimal.NewFromString(val)
}
```

---

© 2025 VDMSquare