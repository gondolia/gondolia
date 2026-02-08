# Identity Service

Authentication and Authorization service for the Webshop V3 platform.

## Features

- Multi-tenant authentication
- JWT-based access and refresh tokens
- User management with company assignments
- Role-based access control (RBAC)
- Password reset flow
- User invitation flow

## Quick Start

### Prerequisites

- Go 1.25+
- k3d (Kubernetes in Docker)
- Skaffold
- Make

### 1. Start Development Environment (from project root)

```bash
# First time setup
make setup

# Start Kubernetes dev environment
make dev
```

This will:
- Build the Identity Service Docker image
- Deploy to local k3d cluster
- Deploy PostgreSQL database
- Set up port forwarding (Service: localhost:8081, DB: localhost:5432)

### 2. Initialize Database (in a new terminal)

```bash
make db-setup
```

This runs migrations and seeds test data.

### 3. Test the API

```bash
# Login
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: demo" \
  -d '{"email":"admin@demo.local","password":"admin123"}'
```

## Test Credentials

| User | Email | Password | Role |
|------|-------|----------|------|
| Admin | admin@demo.local | admin123 | Administrator |
| User | user@demo.local | test123 | Benutzer |

**Tenant:** `demo` (required in `X-Tenant-ID` header)

---

## Testing

### Run All Tests

```bash
cd services/identity
go test ./... -v
```

### Run Specific Test Packages

```bash
# Auth tests (JWT, password hashing)
go test ./internal/auth/... -v

# Service tests (business logic)
go test ./internal/service/... -v
```

### Test Coverage

```bash
go test ./... -cover
```

### Test Structure

```
internal/
├── auth/
│   ├── password_test.go     # Password hashing, validation
│   └── jwt_test.go          # JWT generation, validation, expiry
├── service/
│   └── auth_service_test.go # Login, logout, refresh, password reset
└── repository/
    └── mocks/
        └── mocks.go         # Mock implementations for testing
```

### Test Cases

#### Password Tests (`auth/password_test.go`)

| Test | Description |
|------|-------------|
| `TestHashPassword` | Verifies bcrypt hashing works correctly |
| `TestVerifyPassword` | Tests password verification against hashes |
| `TestValidatePasswordStrength` | Ensures password requirements are enforced |
| `TestGenerateSecureToken` | Verifies cryptographic token generation |
| `TestHashToken` | Tests SHA256 token hashing for storage |

#### JWT Tests (`auth/jwt_test.go`)

| Test | Description |
|------|-------------|
| `TestJWTManager_GenerateAndValidateAccessToken` | Full access token lifecycle |
| `TestJWTManager_GenerateAndValidateRefreshToken` | Full refresh token lifecycle |
| `TestJWTManager_InvalidToken` | Handles invalid/tampered tokens |
| `TestJWTManager_WrongSecret` | Rejects tokens signed with wrong secret |
| `TestJWTManager_TokenExpiry` | Validates expired tokens are rejected |
| `TestDefaultTokenConfig` | Verifies default configuration values |

#### Auth Service Tests (`service/auth_service_test.go`)

| Test | Description |
|------|-------------|
| `TestAuthService_Login_Success` | Successful login returns tokens |
| `TestAuthService_Login_WrongPassword` | Wrong password returns error |
| `TestAuthService_Login_UserNotFound` | Non-existent user returns error |
| `TestAuthService_Login_InactiveUser` | Inactive user cannot login |
| `TestAuthService_Login_SSOOnlyUser` | SSO-only user cannot use password |
| `TestAuthService_RefreshToken_Success` | Token refresh works correctly |
| `TestAuthService_RefreshToken_InvalidToken` | Invalid refresh token rejected |
| `TestAuthService_Logout_Success` | Logout revokes refresh token |
| `TestAuthService_GetCurrentUser_Success` | Returns user with company context |
| `TestAuthService_ForgotPassword_Success` | Generates password reset token |
| `TestAuthService_ForgotPassword_NonexistentUser` | Silent fail for security |
| `TestAuthService_ResetPassword_Success` | Password reset works |
| `TestAuthService_ResetPassword_WeakPassword` | Weak passwords rejected |
| `TestAuthService_ResetPassword_InvalidToken` | Invalid reset token rejected |

---

## API Endpoints

### Authentication (Public)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/login` | Login with email/password |
| POST | `/api/v1/auth/refresh` | Refresh access token |
| POST | `/api/v1/auth/forgot-password` | Request password reset |
| POST | `/api/v1/auth/reset-password` | Reset password with token |
| GET | `/api/v1/auth/invitations/:token` | Validate invitation |
| POST | `/api/v1/auth/invitations/:token/accept` | Accept invitation |

### Authentication (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/logout` | Logout (revoke token) |
| GET | `/api/v1/auth/me` | Get current user |
| POST | `/api/v1/auth/switch-company` | Switch company context |

### Users

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/users` | List users |
| POST | `/api/v1/users` | Create user |
| GET | `/api/v1/users/:id` | Get user |
| PUT | `/api/v1/users/:id` | Update user |
| DELETE | `/api/v1/users/:id` | Delete user |
| POST | `/api/v1/users/:id/activate` | Activate user |
| POST | `/api/v1/users/:id/deactivate` | Deactivate user |
| POST | `/api/v1/users/invite` | Invite user to company |

### Companies

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/companies` | List companies |
| POST | `/api/v1/companies` | Create company |
| GET | `/api/v1/companies/:id` | Get company |
| PUT | `/api/v1/companies/:id` | Update company |
| DELETE | `/api/v1/companies/:id` | Delete company |
| GET | `/api/v1/companies/:id/users` | List company users |
| POST | `/api/v1/companies/:id/users` | Add user to company |
| PUT | `/api/v1/companies/:id/users/:userId` | Update user role |
| DELETE | `/api/v1/companies/:id/users/:userId` | Remove user |

### Roles

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/roles` | List roles |
| POST | `/api/v1/roles` | Create role |
| GET | `/api/v1/roles/:id` | Get role |
| PUT | `/api/v1/roles/:id` | Update role |
| DELETE | `/api/v1/roles/:id` | Delete role |

### Health

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health/live` | Liveness probe |
| GET | `/health/ready` | Readiness probe |

---

## API Examples

### Login

```bash
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: demo" \
  -d '{"email":"admin@demo.local","password":"admin123"}'
```

Response:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "token_type": "Bearer",
  "expires_in": 900
}
```

### Get Current User

```bash
curl http://localhost:8081/api/v1/auth/me \
  -H "Authorization: Bearer <access_token>" \
  -H "X-Tenant-ID: demo"
```

Response:
```json
{
  "user": {
    "id": "00000000-0000-0000-0000-000000000005",
    "email": "admin@demo.local",
    "firstname": "Admin",
    "lastname": "User",
    "is_active": true,
    "is_salesmaster": true
  },
  "company": {
    "id": "00000000-0000-0000-0000-000000000002",
    "name": "Demo Company GmbH",
    "sap_company_number": "1000"
  },
  "role": {
    "name": "Administrator",
    "permissions": {
      "company.manage": true,
      "sales.create-order": true
    }
  },
  "permissions": ["company.manage", "sales.create-order", ...]
}
```

### List Users

```bash
curl http://localhost:8081/api/v1/users \
  -H "Authorization: Bearer <access_token>" \
  -H "X-Tenant-ID: demo"
```

### Refresh Token

```bash
curl -X POST http://localhost:8081/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: demo" \
  -d '{"refresh_token":"<refresh_token>"}'
```

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `HTTP_PORT` | 8080 | HTTP server port |
| `GRPC_PORT` | 9090 | gRPC server port |
| `DATABASE_HOST` | localhost | PostgreSQL host |
| `DATABASE_PORT` | 5432 | PostgreSQL port |
| `DATABASE_NAME` | identity | Database name |
| `DATABASE_USER` | postgres | Database user |
| `DATABASE_PASSWORD` | postgres | Database password |
| `JWT_ACCESS_SECRET` | *required* | Secret for access tokens (min 32 chars) |
| `JWT_REFRESH_SECRET` | *required* | Secret for refresh tokens (min 32 chars) |
| `JWT_ACCESS_TOKEN_EXPIRY` | 15m | Access token lifetime |
| `JWT_REFRESH_TOKEN_EXPIRY` | 168h | Refresh token lifetime (7 days) |

---

## Project Structure

```
identity/
├── cmd/
│   ├── server/              # Main server entrypoint
│   └── migrate/             # Migration and seed tool
├── internal/
│   ├── auth/                # JWT and password utilities
│   │   ├── jwt.go           # JWT token generation/validation
│   │   ├── jwt_test.go      # JWT tests
│   │   ├── password.go      # Password hashing
│   │   └── password_test.go # Password tests
│   ├── config/              # Configuration loading
│   ├── domain/              # Domain models and errors
│   │   ├── user.go          # User entity
│   │   ├── company.go       # Company entity
│   │   ├── role.go          # Role entity with permissions
│   │   ├── auth.go          # Auth-related models
│   │   └── errors.go        # Domain errors
│   ├── handler/             # HTTP handlers
│   │   ├── auth.go          # Auth endpoints
│   │   ├── user.go          # User endpoints
│   │   ├── company.go       # Company endpoints
│   │   ├── role.go          # Role endpoints
│   │   └── health.go        # Health checks
│   ├── middleware/          # HTTP middleware
│   │   ├── auth.go          # JWT authentication
│   │   └── tenant.go        # Tenant extraction
│   ├── repository/          # Data access layer
│   │   ├── interfaces.go    # Repository interfaces
│   │   ├── mocks/           # Mock implementations
│   │   └── postgres/        # PostgreSQL implementations
│   └── service/             # Business logic
│       ├── auth_service.go      # Authentication logic
│       ├── auth_service_test.go # Auth service tests
│       ├── user_service.go      # User management
│       ├── company_service.go   # Company management
│       └── role_service.go      # Role management
├── migrations/              # SQL migration files
├── Dockerfile               # Production Dockerfile
├── Dockerfile.debug         # Debug Dockerfile with Delve
├── Makefile                 # Local development commands
└── go.mod                   # Go module definition
```

---

## Database Schema

### Tables

| Table | Description |
|-------|-------------|
| `tenants` | Multi-tenant configuration |
| `companies` | Companies within tenants |
| `roles` | Roles with permissions |
| `users` | User accounts |
| `user_companies` | User-Company-Role assignments |
| `refresh_tokens` | Active refresh tokens |
| `password_resets` | Password reset tokens |
| `authentication_logs` | Audit log for auth events |

### Migrations

```bash
# Run migrations (from project root with port-forward active)
make db-setup

# Or manually
cd services/identity
go run ./cmd/migrate -command=up
go run ./cmd/migrate -command=seed
```

---

## Events Published

| Event | Trigger | Payload |
|-------|---------|---------|
| `user.created` | New user registered | User ID, Company ID |
| `user.updated` | User modified | User ID, Changed fields |
| `user.deleted` | User removed | User ID |
| `company.updated` | Company modified | Company ID |

---

## Development Workflow

### Making Changes

1. Edit code in `services/identity/`
2. Skaffold automatically rebuilds and redeploys
3. Test changes via `curl` or API client

### Running Tests

```bash
cd services/identity
go test ./... -v
```

### Debugging

```bash
# From project root
make debug
```

This starts the service with Delve debugger on port 40000.

### Viewing Logs

```bash
# From project root
make logs
```

### Observability

```bash
# Open Jaeger (traces) and Grafana (logs)
make observe
```

- Jaeger: http://localhost:16686
- Grafana: http://localhost:3000
