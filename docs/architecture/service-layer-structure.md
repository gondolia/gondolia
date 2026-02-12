# Service Layer Structure

Alle Backend Services im Webshop V3 folgen einer **konsistenten 4-Layer-Architektur**.

---

## Verzeichnisstruktur

```
services/
  {service-name}/
    cmd/
      main.go                    # Entry point, DI Setup, Server Start
    internal/
      {domain}/
        handler.go               # HTTP/gRPC Handler (Controller)
        service.go               # Business Logic
        repository.go            # Database Access
        models.go                # Data Structures, DTOs
        errors.go                # Domain-spezifische Fehler
        {domain}_test.go         # Unit Tests
      middleware/
        auth.go                  # JWT Authentication
        logging.go               # Request Logging
        tracing.go               # OpenTelemetry Integration
    api/
      openapi.yaml               # OpenAPI 3.0 Spec
      proto/                     # Protocol Buffers (für gRPC)
    migrations/
      000001_init.up.sql
      000001_init.down.sql
```

### Beispiel: Identity Service

```
services/identity/
  cmd/
    main.go
  internal/
    user/
      handler.go
      service.go
      repository.go
      models.go
      errors.go
      user_test.go
    company/
      handler.go
      service.go
      repository.go
      models.go
    auth/
      handler.go
      service.go
      token.go
    middleware/
      auth.go
      tenant.go
  api/
    openapi.yaml
  migrations/
    000001_create_users.up.sql
    000001_create_users.down.sql
```

---

## Layer-Verantwortlichkeiten

### 1. Handler Layer (`handler.go`)

**Zweck:** HTTP/gRPC Request-Handling

**Verantwortlichkeiten:**
- Request Parameters parsen (Path, Query, Body)
- Input-Validierung (Format, Required Fields)
- User/Tenant aus JWT-Context extrahieren
- Service Layer aufrufen
- Response formatieren (JSON/Protobuf)
- HTTP Status Codes setzen
- Tracing Spans starten

**Beispiel:**

```go
// internal/user/handler.go
package user

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/gondolia/gondolia/pkg/telemetry"
)

type Handler struct {
    service Service
    logger  *zap.Logger
}

func NewHandler(service Service, logger *zap.Logger) *Handler {
    return &Handler{service: service, logger: logger}
}

func (h *Handler) GetUser(c *gin.Context) {
    ctx, span := telemetry.Start(c.Request.Context(), "Handler.GetUser")
    defer span.End()

    userID := c.Param("id")
    if userID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
        return
    }

    user, err := h.service.GetByID(ctx, userID)
    if err != nil {
        h.handleError(c, err)
        return
    }

    c.JSON(http.StatusOK, user)
}

func (h *Handler) CreateUser(c *gin.Context) {
    ctx, span := telemetry.Start(c.Request.Context(), "Handler.CreateUser")
    defer span.End()

    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
        return
    }

    // Tenant aus JWT Context
    tenantID := c.GetString("tenant_id")

    user, err := h.service.Create(ctx, tenantID, &req)
    if err != nil {
        h.handleError(c, err)
        return
    }

    c.JSON(http.StatusCreated, user)
}

func (h *Handler) handleError(c *gin.Context, err error) {
    switch {
    case errors.Is(err, ErrUserNotFound):
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
    case errors.Is(err, ErrEmailAlreadyExists):
        c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
    case errors.Is(err, ErrValidation):
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    default:
        h.logger.Error("internal error", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
    }
}
```

**Best Practices:**
- Handler **dünn** halten (minimale Logik)
- **Keine Business Logic** im Handler
- Konsistente Error Responses
- Logging mit Context
- Tracing Spans für Observability

---

### 2. Service Layer (`service.go`)

**Zweck:** Business Logic implementieren

**Verantwortlichkeiten:**
- Business Rules validieren
- Mehrere Repositories koordinieren
- Datenbank-Transaktionen managen
- Domain Logic implementieren
- Externe Services aufrufen
- Side Effects triggern (Events, Emails)

**Beispiel:**

```go
// internal/user/service.go
package user

import (
    "context"
    "github.com/jmoiron/sqlx"
    "github.com/gondolia/gondolia/pkg/telemetry"
)

// Interface für Testbarkeit
type Service interface {
    GetByID(ctx context.Context, id string) (*User, error)
    Create(ctx context.Context, tenantID string, req *CreateUserRequest) (*User, error)
    Update(ctx context.Context, id string, req *UpdateUserRequest) (*User, error)
    Delete(ctx context.Context, id string) error
}

type service struct {
    repo       Repository
    db         *sqlx.DB
    companyRepo company.Repository
    eventBus   events.Publisher
    logger     *zap.Logger
}

func NewService(
    repo Repository,
    db *sqlx.DB,
    companyRepo company.Repository,
    eventBus events.Publisher,
    logger *zap.Logger,
) Service {
    return &service{
        repo:        repo,
        db:          db,
        companyRepo: companyRepo,
        eventBus:    eventBus,
        logger:      logger,
    }
}

func (s *service) Create(ctx context.Context, tenantID string, req *CreateUserRequest) (*User, error) {
    ctx, span := telemetry.Start(ctx, "Service.CreateUser")
    defer span.End()

    // Business Validierung
    if err := s.validateCreateRequest(req); err != nil {
        return nil, err
    }

    // Email-Duplikat prüfen
    existing, _ := s.repo.GetByEmail(ctx, tenantID, req.Email)
    if existing != nil {
        return nil, ErrEmailAlreadyExists
    }

    // Company existiert?
    if req.CompanyID != "" {
        company, err := s.companyRepo.GetByID(ctx, req.CompanyID)
        if err != nil {
            return nil, ErrCompanyNotFound
        }
        if company.TenantID != tenantID {
            return nil, ErrCompanyNotFound // Tenant Isolation
        }
    }

    // Transaktion starten
    tx, err := s.db.BeginTxx(ctx, nil)
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()

    // User erstellen
    user := &User{
        ID:        uuid.New().String(),
        TenantID:  tenantID,
        Email:     req.Email,
        FirstName: req.FirstName,
        LastName:  req.LastName,
        CompanyID: req.CompanyID,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    if err := s.repo.CreateTx(ctx, tx, user); err != nil {
        return nil, err
    }

    // Commit
    if err := tx.Commit(); err != nil {
        return nil, err
    }

    // Event publishen (async)
    go s.eventBus.Publish(ctx, "user.created", UserCreatedEvent{
        UserID:   user.ID,
        TenantID: tenantID,
        Email:    user.Email,
    })

    return user, nil
}

func (s *service) validateCreateRequest(req *CreateUserRequest) error {
    if req.Email == "" {
        return fmt.Errorf("%w: email required", ErrValidation)
    }
    if !isValidEmail(req.Email) {
        return fmt.Errorf("%w: invalid email format", ErrValidation)
    }
    if req.FirstName == "" {
        return fmt.Errorf("%w: first_name required", ErrValidation)
    }
    return nil
}
```

**Best Practices:**
- Service als **Interface** definieren (Dependency Injection, Testing)
- **Transaktionen** für zusammengehörige Operationen
- Business Rules **hier** implementieren, nicht im Handler/Repository
- Externe Aufrufe **hier** koordinieren
- Events/Side Effects **nach** erfolgreichem Commit
- Keine HTTP-spezifischen Konzepte (Status Codes, etc.)

---

### 3. Repository Layer (`repository.go`)

**Zweck:** Datenbank-Zugriff abstrahieren

**Verantwortlichkeiten:**
- SQL Queries ausführen
- Struct Mapping (Rows → Structs)
- Prepared Statements nutzen
- Transaktionen unterstützen
- **KEINE Business Logic**

**Beispiel:**

```go
// internal/user/repository.go
package user

import (
    "context"
    "database/sql"
    "github.com/jmoiron/sqlx"
)

type Repository interface {
    GetByID(ctx context.Context, id string) (*User, error)
    GetByEmail(ctx context.Context, tenantID, email string) (*User, error)
    List(ctx context.Context, tenantID string, filter ListFilter) ([]*User, error)
    Create(ctx context.Context, user *User) error
    CreateTx(ctx context.Context, tx *sqlx.Tx, user *User) error
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id string) error
}

type repository struct {
    db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
    return &repository{db: db}
}

func (r *repository) GetByID(ctx context.Context, id string) (*User, error) {
    var user User
    err := r.db.GetContext(ctx, &user, `
        SELECT id, tenant_id, email, first_name, last_name,
               company_id, created_at, updated_at
        FROM users
        WHERE id = $1 AND deleted_at IS NULL
    `, id)

    if err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrUserNotFound
        }
        return nil, err
    }
    return &user, nil
}

func (r *repository) GetByEmail(ctx context.Context, tenantID, email string) (*User, error) {
    var user User
    err := r.db.GetContext(ctx, &user, `
        SELECT id, tenant_id, email, first_name, last_name,
               company_id, created_at, updated_at
        FROM users
        WHERE tenant_id = $1 AND email = $2 AND deleted_at IS NULL
    `, tenantID, email)

    if err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrUserNotFound
        }
        return nil, err
    }
    return &user, nil
}

func (r *repository) CreateTx(ctx context.Context, tx *sqlx.Tx, user *User) error {
    _, err := tx.ExecContext(ctx, `
        INSERT INTO users (id, tenant_id, email, first_name, last_name,
                          company_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `, user.ID, user.TenantID, user.Email, user.FirstName,
       user.LastName, user.CompanyID, user.CreatedAt, user.UpdatedAt)

    return err
}

func (r *repository) List(ctx context.Context, tenantID string, filter ListFilter) ([]*User, error) {
    query := `
        SELECT id, tenant_id, email, first_name, last_name,
               company_id, created_at, updated_at
        FROM users
        WHERE tenant_id = $1 AND deleted_at IS NULL
    `
    args := []interface{}{tenantID}
    argIndex := 2

    // Optional: Company Filter
    if filter.CompanyID != "" {
        query += fmt.Sprintf(" AND company_id = $%d", argIndex)
        args = append(args, filter.CompanyID)
        argIndex++
    }

    // Pagination
    query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
    args = append(args, filter.Limit, filter.Offset)

    var users []*User
    err := r.db.SelectContext(ctx, &users, query, args...)
    if err != nil {
        return nil, err
    }
    return users, nil
}
```

**Best Practices:**
- **Parameterized Queries** (SQL Injection Prevention)
- **sqlx** für Struct Mapping verwenden
- Domain-spezifische Fehler zurückgeben (`ErrUserNotFound`)
- Transaktionen unterstützen (`*sqlx.Tx` Parameter)
- `FOR UPDATE` bei Concurrent Access
- Queries **nur hier**, nicht im Service

---

### 4. Models Layer (`models.go`)

**Zweck:** Datenstrukturen definieren

**Enthält:**
- Entity Structs (mit DB-Tags)
- Request/Response DTOs
- Domain-spezifische Fehler

**Beispiel:**

```go
// internal/user/models.go
package user

import "time"

// Entity (Datenbank)
type User struct {
    ID        string    `db:"id" json:"id"`
    TenantID  string    `db:"tenant_id" json:"tenant_id"`
    Email     string    `db:"email" json:"email"`
    FirstName string    `db:"first_name" json:"first_name"`
    LastName  string    `db:"last_name" json:"last_name"`
    CompanyID string    `db:"company_id" json:"company_id,omitempty"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
    UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
    DeletedAt *time.Time `db:"deleted_at" json:"-"`
}

// Request DTOs
type CreateUserRequest struct {
    Email     string `json:"email" binding:"required,email"`
    FirstName string `json:"first_name" binding:"required"`
    LastName  string `json:"last_name" binding:"required"`
    CompanyID string `json:"company_id,omitempty"`
}

type UpdateUserRequest struct {
    FirstName *string `json:"first_name,omitempty"`
    LastName  *string `json:"last_name,omitempty"`
    CompanyID *string `json:"company_id,omitempty"`
}

// Filter für List-Operationen
type ListFilter struct {
    CompanyID string
    Limit     int
    Offset    int
}

// Response DTOs (falls anders als Entity)
type UserResponse struct {
    ID        string `json:"id"`
    Email     string `json:"email"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    FullName  string `json:"full_name"`
}

func (u *User) ToResponse() *UserResponse {
    return &UserResponse{
        ID:        u.ID,
        Email:     u.Email,
        FirstName: u.FirstName,
        LastName:  u.LastName,
        FullName:  u.FirstName + " " + u.LastName,
    }
}
```

**errors.go:**

```go
// internal/user/errors.go
package user

import "errors"

var (
    ErrUserNotFound      = errors.New("user not found")
    ErrEmailAlreadyExists = errors.New("email already exists")
    ErrCompanyNotFound   = errors.New("company not found")
    ErrValidation        = errors.New("validation error")
    ErrUnauthorized      = errors.New("unauthorized")
)
```

---

## Dependency Injection (main.go)

```go
// services/identity/cmd/main.go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/jmoiron/sqlx"
    "go.uber.org/zap"

    "github.com/gondolia/gondolia/services/identity/internal/user"
    "github.com/gondolia/gondolia/services/identity/internal/company"
    "github.com/gondolia/gondolia/pkg/telemetry"
)

func main() {
    // Logger
    logger, _ := zap.NewProduction()
    defer logger.Sync()

    // Telemetry
    shutdown := telemetry.Init("identity-service")
    defer shutdown()

    // Database
    db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
    if err != nil {
        logger.Fatal("failed to connect to database", zap.Error(err))
    }
    defer db.Close()

    // Event Bus
    eventBus := events.NewKafkaPublisher(os.Getenv("KAFKA_BROKERS"))

    // Repositories
    userRepo := user.NewRepository(db)
    companyRepo := company.NewRepository(db)

    // Services (mit Dependencies)
    userService := user.NewService(userRepo, db, companyRepo, eventBus, logger)
    companyService := company.NewService(companyRepo, db, logger)

    // Handlers
    userHandler := user.NewHandler(userService, logger)
    companyHandler := company.NewHandler(companyService, logger)

    // Router
    router := gin.New()
    router.Use(middleware.Tracing())
    router.Use(middleware.Logging(logger))
    router.Use(middleware.Auth(os.Getenv("JWT_SECRET")))

    // Routes
    v1 := router.Group("/api/v1")
    {
        users := v1.Group("/users")
        {
            users.GET("/:id", userHandler.GetUser)
            users.POST("", userHandler.CreateUser)
            users.PUT("/:id", userHandler.UpdateUser)
            users.DELETE("/:id", userHandler.DeleteUser)
        }

        companies := v1.Group("/companies")
        {
            companies.GET("/:id", companyHandler.GetCompany)
            companies.POST("", companyHandler.CreateCompany)
        }
    }

    // Start Server
    logger.Info("starting identity-service", zap.String("port", "8080"))
    if err := router.Run(":8080"); err != nil {
        logger.Fatal("failed to start server", zap.Error(err))
    }
}
```

---

## Zusammenfassung: Layer-Regeln

| Layer | Darf aufrufen | Darf NICHT |
|-------|---------------|------------|
| **Handler** | Service | Repository, DB direkt |
| **Service** | Repository, andere Services, Event Bus | Handler, HTTP-Konzepte |
| **Repository** | Datenbank | Service, Handler, Business Logic |
| **Models** | - | Keine Logik, nur Datenstrukturen |

**Flow:**

```
HTTP Request
    ↓
Handler (Parse, Validate, Auth Context)
    ↓
Service (Business Logic, Transactions, Coordination)
    ↓
Repository (SQL, Mapping)
    ↓
Database
```

---

## Service-Splitting: Im Zweifel teilen!

### Grundprinzip

> **Lieber mehrere kleine, fokussierte Services als ein großer, komplexer.**

Die Business-Logik im Webshop ist teilweise komplex (Preisberechnung, ZZON, SAP-Integration). Um Übersichtlichkeit und Wartbarkeit zu gewährleisten:

- **Jeder Service sollte EINE klare Verantwortung haben**
- **Im Zweifel: Service aufteilen**
- **Komplexität verteilen statt konzentrieren**

### Wann splitten?

| Indikator | Schwellwert | Aktion |
|-----------|-------------|--------|
| Domains im `internal/` | > 3-4 | Splitten erwägen |
| Zeilen in `service.go` | > 500 | Splitten oder Subservices |
| Unterschiedliche Skalierung | Ja | Definitiv splitten |
| Team-Ownership unklar | Ja | Splitten |
| Deployment-Zyklen unterschiedlich | Ja | Splitten |

### Beispiele für Splitting

#### Beispiel 1: Order-Domäne

```
# ❌ FALSCH: Zu viel in einem Service
services/order/
  internal/
    order/           # Bestellungen erstellen
    quote/           # Angebote
    payment/         # Zahlungsabwicklung
    invoice/         # Rechnungserstellung
    returns/         # Retouren
    pricing/         # Preisberechnung

# ✅ RICHTIG: Fokussierte Services
services/order/              # Bestellungen erstellen & verwalten
  internal/
    order/
    orderitem/

services/quote/              # Angebotserstellung (eigene komplexe Logik)
  internal/
    quote/
    quoteitem/
    pricing/                 # Quote-spezifische Preislogik

services/payment/            # Zahlungsabwicklung
  internal/
    payment/
    transaction/
    refund/

services/billing/            # Rechnungen & Gutschriften
  internal/
    invoice/
    creditnote/

services/returns/            # Retouren-Management
  internal/
    return/
    returnitem/
```

#### Beispiel 2: Inventory-Domäne

```
# ❌ FALSCH: Zu viele Verantwortlichkeiten
services/inventory/
  internal/
    stock/           # Lagerbestand
    plant/           # Werke
    zone/            # Lieferzonen
    zzon/            # PLZ → Werk Zuordnung
    availability/    # Verfügbarkeit
    reservation/     # Reservierungen

# ✅ RICHTIG: Getrennte Concerns
services/inventory/          # Kernbestand
  internal/
    stock/
    reservation/

services/plant/              # Werke & Zonen
  internal/
    plant/
    zone/
    zzon/                    # Komplexe ZZON-Logik isoliert

services/availability/       # Verfügbarkeitsberechnung
  internal/
    availability/
    calculation/
```

#### Beispiel 3: Identity-Domäne

```
# ⚠️ GRENZFALL: Könnte so bleiben oder gesplittet werden
services/identity/
  internal/
    user/            # Benutzer
    company/         # Firmen
    auth/            # Authentifizierung
    permission/      # Berechtigungen

# Optional splitten wenn Auth komplex wird:
services/identity/           # User & Company Management
  internal/
    user/
    company/

services/auth/               # Authentifizierung & Token
  internal/
    auth/
    token/
    session/

services/authorization/      # Berechtigungen (RBAC/ABAC)
  internal/
    permission/
    role/
    policy/
```

### Service-Kommunikation nach Splitting

#### Synchron (gRPC) - Für direkte Abfragen

```go
// order-service ruft inventory-service
func (s *orderService) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*Order, error) {
    // Verfügbarkeit prüfen via gRPC
    availResp, err := s.inventoryClient.CheckAvailability(ctx, &CheckAvailabilityRequest{
        ProductID: req.ProductID,
        Quantity:  req.Quantity,
        PlantID:   req.PlantID,
    })
    if err != nil {
        return nil, err
    }
    if !availResp.Available {
        return nil, ErrProductNotAvailable
    }

    // Order erstellen...
}
```

#### Asynchron (Kafka) - Für Events/Zustandsänderungen

```go
// order-service published Event
func (s *orderService) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*Order, error) {
    // Order erstellen...
    order, err := s.repo.Create(ctx, order)
    if err != nil {
        return nil, err
    }

    // Event publishen - andere Services reagieren
    s.eventBus.Publish(ctx, "order.created", OrderCreatedEvent{
        OrderID:    order.ID,
        CustomerID: order.CustomerID,
        Items:      order.Items,
        TotalAmount: order.TotalAmount,
    })

    return order, nil
}

// inventory-service konsumiert Event
func (s *inventoryService) HandleOrderCreated(ctx context.Context, event OrderCreatedEvent) error {
    // Bestand reservieren
    for _, item := range event.Items {
        if err := s.reserveStock(ctx, item.ProductID, item.Quantity, event.OrderID); err != nil {
            // Compensation Event publishen
            s.eventBus.Publish(ctx, "stock.reservation.failed", StockReservationFailedEvent{
                OrderID: event.OrderID,
                Reason:  err.Error(),
            })
            return err
        }
    }

    s.eventBus.Publish(ctx, "stock.reserved", StockReservedEvent{
        OrderID: event.OrderID,
    })
    return nil
}
```

### Event-Choreographie vs. Orchestrierung

#### Choreographie (Empfohlen für die meisten Fälle)

Services reagieren selbstständig auf Events:

```
Order Service                    Kafka                         Other Services
     │                             │                                │
     │  order.created              │                                │
     │  ─────────────────────────▶ │                                │
     │                             │  ─────────────────────────▶    │
     │                             │                                │
     │                             │  inventory: reserve stock      │
     │                             │  payment: initiate payment     │
     │                             │  notification: send email      │
     │                             │                                │
     │  stock.reserved             │                                │
     │  ◀───────────────────────── │  ◀─────────────────────────    │
     │                             │                                │
```

#### Orchestrierung (Für komplexe Workflows)

Ein Saga-Coordinator steuert den Ablauf:

```
Saga Coordinator                 Services
     │                              │
     │  1. Create Order             │
     │  ──────────────────────────▶ │ Order Service
     │                              │
     │  2. Reserve Stock            │
     │  ──────────────────────────▶ │ Inventory Service
     │                              │
     │  3. Process Payment          │
     │  ──────────────────────────▶ │ Payment Service
     │                              │
     │  4. Confirm Order            │
     │  ──────────────────────────▶ │ Order Service
     │                              │
```

### Datenbank-Strategie bei Splitting

Jeder Service hat seine **eigene Datenbank/Schema**:

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Order     │     │  Inventory  │     │   Payment   │
│   Service   │     │   Service   │     │   Service   │
└──────┬──────┘     └──────┬──────┘     └──────┬──────┘
       │                   │                   │
       ▼                   ▼                   ▼
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   orders    │     │  inventory  │     │  payments   │
│   schema    │     │   schema    │     │   schema    │
└─────────────┘     └─────────────┘     └─────────────┘
       │                   │                   │
       └───────────────────┴───────────────────┘
                           │
                    ┌──────▼──────┐
                    │ PostgreSQL  │
                    │   Cluster   │
                    └─────────────┘
```

**Wichtig:**
- **Keine Cross-Schema Joins** - Daten über gRPC/Events holen
- **Eventual Consistency** akzeptieren
- **Idempotente Event-Handler** implementieren

### Checkliste vor dem Splitten

- [ ] Klare Bounded Context Grenzen definiert?
- [ ] API Contract zwischen Services definiert?
- [ ] Event-Schema dokumentiert?
- [ ] Fehlerbehandlung & Compensation geplant?
- [ ] Monitoring & Tracing vorbereitet?
- [ ] Deployment-Pipeline angepasst?

---

## Testing

### Unit Tests

```go
// internal/user/service_test.go
func TestService_Create(t *testing.T) {
    mockRepo := &MockRepository{}
    mockEventBus := &MockEventBus{}
    service := NewService(mockRepo, nil, nil, mockEventBus, zap.NewNop())

    t.Run("successful creation", func(t *testing.T) {
        mockRepo.On("GetByEmail", mock.Anything, "tenant-1", "test@example.com").
            Return(nil, ErrUserNotFound)
        mockRepo.On("CreateTx", mock.Anything, mock.Anything, mock.Anything).
            Return(nil)

        user, err := service.Create(context.Background(), "tenant-1", &CreateUserRequest{
            Email:     "test@example.com",
            FirstName: "Test",
            LastName:  "User",
        })

        assert.NoError(t, err)
        assert.Equal(t, "test@example.com", user.Email)
    })

    t.Run("email already exists", func(t *testing.T) {
        mockRepo.On("GetByEmail", mock.Anything, "tenant-1", "existing@example.com").
            Return(&User{Email: "existing@example.com"}, nil)

        _, err := service.Create(context.Background(), "tenant-1", &CreateUserRequest{
            Email:     "existing@example.com",
            FirstName: "Test",
            LastName:  "User",
        })

        assert.ErrorIs(t, err, ErrEmailAlreadyExists)
    })
}
```

---

## Weiterführende Dokumentation

- [Architektur-Übersicht](./overview.md)
- [Observability Guide](../observability/README.md)
- [Service-Übersicht](../services/README.md)
