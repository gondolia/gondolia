# Catalog Service

Product catalog service for the Gondolia B2B E-Commerce platform.

## Features

- **Products**: Manage products with multi-language support (i18n)
- **Categories**: Hierarchical category tree
- **Pricing**: B2B contract pricing with tier support
- **PIM Integration**: Sync from external PIM systems (Akeneo, Pimcore, etc.)
- **Search**: Product search via search provider (Meilisearch, Algolia, etc.)
- **Multi-tenant**: Full tenant isolation

## Architecture

Follows the same 4-layer architecture as identity service:

```
Handler → Service → Repository → Domain
```

### Layers

1. **Domain** (`internal/domain/`): Pure business models and errors
2. **Repository** (`internal/repository/`): Data access interfaces + PostgreSQL implementations
3. **Service** (`internal/service/`): Business logic
4. **Handler** (`internal/handler/`): HTTP/REST endpoints (Gin)

## API Endpoints

### Products

- `GET /api/v1/products` - List products (paginated, filterable)
- `GET /api/v1/products/:id` - Get product by ID
- `POST /api/v1/products` - Create product
- `PUT /api/v1/products/:id` - Update product
- `DELETE /api/v1/products/:id` - Delete product (soft delete)

### Categories

- `GET /api/v1/categories` - Get category tree
- `GET /api/v1/categories/list` - List categories (paginated)
- `GET /api/v1/categories/:id` - Get category by ID
- `POST /api/v1/categories` - Create category
- `PUT /api/v1/categories/:id` - Update category
- `DELETE /api/v1/categories/:id` - Delete category

### Prices

- `GET /api/v1/products/:productId/prices` - Get prices for product
- `POST /api/v1/products/:productId/prices` - Create price for product
- `PUT /api/v1/prices/:id` - Update price
- `DELETE /api/v1/prices/:id` - Delete price

### Search

- `GET /api/v1/search?q=...&filters=...` - Search products

### Sync

- `POST /api/v1/sync/pim?full=true` - Trigger PIM sync

## Domain Models

### Product

```go
type Product struct {
    ID          uuid.UUID
    TenantID    uuid.UUID
    SKU         string
    Name        map[string]string    // locale -> name
    Description map[string]string    // locale -> description
    CategoryIDs []uuid.UUID
    Attributes  []ProductAttribute
    Status      ProductStatus        // draft|active|archived
    Images      []ProductImage
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### Category

```go
type Category struct {
    ID        uuid.UUID
    TenantID  uuid.UUID
    Code      string
    ParentID  *uuid.UUID
    Name      map[string]string    // locale -> name
    SortOrder int
    Active    bool
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### Price

```go
type Price struct {
    ID              uuid.UUID
    TenantID        uuid.UUID
    ProductID       uuid.UUID
    CustomerGroupID *uuid.UUID       // nil = base price
    MinQuantity     int
    Price           float64
    Currency        string
    ValidFrom       *time.Time
    ValidTo         *time.Time
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

## Environment Variables

```bash
SERVICE_NAME=catalog-service
HTTP_PORT=8081
GRPC_PORT=9091

# Database
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=catalog
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# PIM Provider
PIM_PROVIDER=mock          # mock|akeneo|pimcore
PIM_URL=
PIM_API_KEY=

# Search Provider
SEARCH_PROVIDER=mock       # mock|meilisearch|algolia
SEARCH_URL=
SEARCH_API_KEY=

# CORS
ALLOWED_ORIGINS=http://localhost:3000
```

## Development

### Run locally

```bash
# Start dependencies
docker compose up -d postgres redis

# Run migrations
cd services/catalog/cmd/migrate
go run main.go

# Start service
cd services/catalog/cmd/server
go run main.go
```

### Run tests

```bash
go test ./services/catalog/...
```

### Build Docker image

```bash
docker build -t catalog-service -f services/catalog/Dockerfile .
```

## Database Schema

### products

- `id` UUID PRIMARY KEY
- `tenant_id` UUID NOT NULL
- `sku` VARCHAR(100) NOT NULL
- `name` JSONB (i18n)
- `description` JSONB (i18n)
- `category_ids` UUID[]
- `attributes` JSONB
- `status` VARCHAR(20)
- `images` JSONB
- `pim_identifier` VARCHAR(255)
- `last_synced_at` TIMESTAMP
- `created_at` TIMESTAMP
- `updated_at` TIMESTAMP
- `deleted_at` TIMESTAMP

**Indexes:**
- `tenant_id, sku` (UNIQUE)
- `tenant_id`
- `status`
- `category_ids` (GIN)

### categories

- `id` UUID PRIMARY KEY
- `tenant_id` UUID NOT NULL
- `code` VARCHAR(100) NOT NULL
- `parent_id` UUID (FK → categories)
- `name` JSONB (i18n)
- `sort_order` INT
- `active` BOOLEAN
- `pim_code` VARCHAR(255)
- `last_synced_at` TIMESTAMP
- `created_at` TIMESTAMP
- `updated_at` TIMESTAMP
- `deleted_at` TIMESTAMP

**Indexes:**
- `tenant_id, code` (UNIQUE)
- `parent_id`

### prices

- `id` UUID PRIMARY KEY
- `tenant_id` UUID NOT NULL
- `product_id` UUID NOT NULL (FK → products)
- `customer_group_id` UUID
- `min_quantity` INT
- `price` DECIMAL(10,2)
- `currency` CHAR(3)
- `valid_from` TIMESTAMP
- `valid_to` TIMESTAMP
- `created_at` TIMESTAMP
- `updated_at` TIMESTAMP
- `deleted_at` TIMESTAMP

**Indexes:**
- `product_id`
- `customer_group_id`
- `valid_from, valid_to`

## Provider Integration

### PIM Provider

The service uses the `provider/pim` interface to sync products and categories from external PIM systems.

Supported providers:
- **Mock**: For development/testing
- **Akeneo**: Enterprise PIM
- **Pimcore**: Open-source PIM

### Search Provider

The service uses the `provider/search` interface to index and search products.

Supported providers:
- **Mock**: For development/testing
- **Meilisearch**: Fast, typo-tolerant search
- **Algolia**: Hosted search service

## Testing

The service includes comprehensive unit tests following the same patterns as the identity service:

```bash
go test -v ./services/catalog/internal/service/...
```

Coverage:
- Product service: CRUD operations, validation
- Category service: Tree operations, circular reference prevention
- Price service: Overlap detection, date range validation

## Docker Compose

The service is integrated into the main `docker-compose.yml`:

```yaml
catalog:
  build:
    context: .
    dockerfile: services/catalog/Dockerfile
  environment:
    DATABASE_NAME: catalog
    HTTP_PORT: "8081"
  labels:
    - "traefik.http.routers.catalog.rule=PathPrefix(`/api/v1/products`)"
```

## Next Steps

1. Implement actual PIM provider (Akeneo/Pimcore)
2. Implement actual search provider (Meilisearch/Algolia)
3. Add authentication middleware (integrate with identity service)
4. Add authorization (permission checks)
5. Implement image upload/storage
6. Add GraphQL API
7. Add caching layer (Redis)
8. Add Kafka events for product changes
9. Performance optimization (batch operations)
10. Add OpenTelemetry tracing
