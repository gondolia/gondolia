# Architektur-Übersicht

## Grundprinzipien

### 1. Domain-Driven Design (DDD)

Jeder Microservice repräsentiert eine **Bounded Context**:

```
┌─────────────────────────────────────────────────────────────────┐
│                      BOUNDED CONTEXTS                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │  Identity   │  │   Catalog   │  │  Commerce   │             │
│  │   Context   │  │   Context   │  │   Context   │             │
│  │             │  │             │  │             │             │
│  │ - User      │  │ - Product   │  │ - Cart      │             │
│  │ - Company   │  │ - Category  │  │ - Order     │             │
│  │ - Auth      │  │ - Price     │  │ - Quote     │             │
│  │ - Permission│  │ - Search    │  │ - Payment   │             │
│  └─────────────┘  └─────────────┘  └─────────────┘             │
│                                                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │  Inventory  │  │  Shipping   │  │Integration  │             │
│  │   Context   │  │   Context   │  │   Context   │             │
│  │             │  │             │  │             │             │
│  │ - Stock     │  │ - Address   │  │ - SAP       │             │
│  │ - Plant     │  │ - Zone      │  │ - ERP Sync  │             │
│  │ - Zone      │  │ - PLZ       │  │ - Webhooks  │             │
│  │ - Availabil.│  │ - Carrier   │  │             │             │
│  └─────────────┘  └─────────────┘  └─────────────┘             │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 2. Event-Driven Architecture

Services kommunizieren primär über **Events**:

```
┌─────────────┐     ┌─────────────────┐     ┌─────────────┐
│   Order     │────▶│   Message Bus   │────▶│  Inventory  │
│   Service   │     │    (Kafka)      │     │   Service   │
└─────────────┘     └─────────────────┘     └─────────────┘
      │                     │                      │
      │   OrderCreated      │                      │
      │   ─────────────▶    │   ─────────────▶     │
      │                     │   StockReserved      │
      │                     │                      │
      │   OrderConfirmed    │                      │
      │   ◀─────────────    │   ◀─────────────     │
      │                     │                      │
```

### 3. API-First Design

Alle Services definieren ihre API via **OpenAPI 3.0** oder **Protocol Buffers**:

```yaml
# Beispiel: order-service/api/v1/order.proto
service OrderService {
  rpc CreateOrder(CreateOrderRequest) returns (Order);
  rpc GetOrder(GetOrderRequest) returns (Order);
  rpc ListOrders(ListOrdersRequest) returns (ListOrdersResponse);
}
```

---

## Kommunikationsmuster

### Synchron (Request/Response)

Für direkte Abfragen, z.B. Produktdetails:

```
Frontend ──HTTP──▶ API Gateway ──gRPC──▶ Catalog Service
                                              │
                                              ▼
                                         PostgreSQL
```

**Verwendung:**
- Produktabfragen
- Authentifizierung
- Preisberechnungen

### Asynchron (Event-Driven)

Für Zustandsänderungen und Benachrichtigungen:

```
Order Service                        Message Bus                    Other Services
     │                                    │                              │
     │  OrderCreated Event                │                              │
     │  ──────────────────────────────▶   │                              │
     │                                    │   ──────────────────────▶    │
     │                                    │   (fan-out to subscribers)   │
     │                                    │                              │
     │                                    │   InventoryService: Reserve  │
     │                                    │   NotificationService: Email │
     │                                    │   AnalyticsService: Track    │
```

**Verwendung:**
- Bestellungen
- Lagerbestandsänderungen
- Benachrichtigungen

---

## Datenbank-Strategie

### Option A: Database per Service (Empfohlen)

Jeder Service hat seine **eigene Datenbank**:

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Identity   │     │   Catalog   │     │    Order    │
│   Service   │     │   Service   │     │   Service   │
└──────┬──────┘     └──────┬──────┘     └──────┬──────┘
       │                   │                   │
       ▼                   ▼                   ▼
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  identity   │     │   catalog   │     │   orders    │
│     DB      │     │     DB      │     │     DB      │
│ (PostgreSQL)│     │ (PostgreSQL)│     │ (PostgreSQL)│
└─────────────┘     └─────────────┘     └─────────────┘
```

**Vorteile:**
- Unabhängige Skalierung
- Isolierte Failures
- Technologie-Freiheit pro Service

**Nachteile:**
- Komplexere Joins (über Events)
- Eventual Consistency

### Option B: Shared Database

Alle Services teilen sich eine Datenbank:

**Vorteile:**
- Einfachere Joins
- Sofortige Konsistenz

**Nachteile:**
- Tight Coupling
- Single Point of Failure
- Schwieriger zu skalieren

---

## Multi-Tenancy

Das System unterstützt **mehrere Mandanten** (Tenants):

```
┌─────────────────────────────────────────────────────────────────┐
│                         API Gateway                              │
│                                                                  │
│  Request Header: X-Tenant-ID: tenant-abc                        │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Every Service                               │
│                                                                  │
│  - Tenant ID in JWT Token                                       │
│  - All queries filtered by tenant_id                            │
│  - Separate indexes per tenant (optional)                       │
└─────────────────────────────────────────────────────────────────┘
```

### Strategien

| Strategie | Beschreibung | Verwendung |
|-----------|--------------|------------|
| **Row-Level** | `WHERE tenant_id = ?` | Standard |
| **Schema-Level** | Separate Schemas pro Tenant | Hohe Isolation |
| **Database-Level** | Separate DBs pro Tenant | Enterprise |

**Empfehlung:** Row-Level für V3.0, Schema-Level als Option

---

## Sicherheit

### Authentication & Authorization

```
┌─────────────────────────────────────────────────────────────────┐
│                         Frontend                                 │
│                                                                  │
│  1. Login Request (username, password)                          │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Identity Service                            │
│                                                                  │
│  2. Validate credentials                                         │
│  3. Generate JWT (access + refresh token)                       │
│  4. Return tokens                                                │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                         Frontend                                 │
│                                                                  │
│  5. Store tokens                                                 │
│  6. Include access token in requests                            │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                        API Gateway                               │
│                                                                  │
│  7. Validate JWT signature                                       │
│  8. Extract user/tenant info                                     │
│  9. Forward to service                                           │
└─────────────────────────────────────────────────────────────────┘
```

### JWT Token Struktur

```json
{
  "sub": "user-123",
  "tenant_id": "tenant-abc",
  "company_id": "company-456",
  "roles": ["customer", "company_admin"],
  "permissions": ["catalog.view", "order.create", "cart.manage"],
  "exp": 1706018400
}
```

---

## Skalierung

### Horizontal Scaling

```
                    ┌─────────────────┐
                    │  Load Balancer  │
                    └────────┬────────┘
                             │
         ┌───────────────────┼───────────────────┐
         │                   │                   │
         ▼                   ▼                   ▼
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Catalog   │     │   Catalog   │     │   Catalog   │
│  Service 1  │     │  Service 2  │     │  Service 3  │
└─────────────┘     └─────────────┘     └─────────────┘
         │                   │                   │
         └───────────────────┼───────────────────┘
                             │
                    ┌────────▼────────┐
                    │    PostgreSQL   │
                    │   (Read Replicas)│
                    └─────────────────┘
```

### Auto-Scaling Rules (Kubernetes)

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: catalog-service-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: catalog-service
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

---

## Monitoring & Observability

### Three Pillars

```
┌─────────────────────────────────────────────────────────────────┐
│                         Observability                            │
├─────────────────┬─────────────────┬─────────────────────────────┤
│                 │                 │                             │
│     Metrics     │     Logs        │     Traces                  │
│   (Prometheus)  │   (Loki/ELK)    │   (Jaeger)                  │
│                 │                 │                             │
│ - Request rate  │ - Error logs    │ - Request flow             │
│ - Latency       │ - Audit logs    │ - Cross-service            │
│ - Error rate    │ - Debug logs    │ - Bottlenecks              │
│ - Saturation    │                 │                             │
│                 │                 │                             │
└─────────────────┴─────────────────┴─────────────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │     Grafana     │
                    │   Dashboards    │
                    └─────────────────┘
```

---

## Deployment Strategy

### GitOps mit ArgoCD

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   GitHub    │────▶│   ArgoCD    │────▶│ Kubernetes  │
│    Repo     │     │             │     │   Cluster   │
└─────────────┘     └─────────────┘     └─────────────┘
      │                                        │
      │  1. Push to main                       │
      │  2. ArgoCD detects change              │
      │  3. Sync to cluster                    │
      │                                        ▼
      │                              ┌─────────────────┐
      │                              │   Production    │
      │                              │   Environment   │
      │                              └─────────────────┘
```

### Environments

| Environment | Zweck | Branch |
|-------------|-------|--------|
| Development | Entwicklung | feature/* |
| Staging | Testing | develop |
| Production | Live | main |

---

## Weiterführende Dokumentation

### Grundlagen & Prinzipien
- [B2B Self-Service Prinzipien](./b2b-self-service-principles.md) - Leitprinzipien für Requirements Engineering (80/20 Regel, Durchlaufzeit, Direktintegration)

### Architektur
- [Service Layer Structure](./service-layer-structure.md) - 4-Layer-Architektur für Backend Services
- [Category Architecture](./category-architecture.md) - Kategorien als First-Class Entität mit Hierarchie, i18n und SEO

### Integrationen
- [Datahub Konzept](./datahub-concept.md) - Zentrale Integrationsschicht für konfigurationsbasierte Schnittstellen
- [Externe Integrationen](./external-integrations-analysis.md) - Analyse aller 24 externen Services (Payment, PIM, Search, Auth, etc.)
- [Akeneo Import Analyse](./akeneo-import-analysis.md) - Produktimport aus Akeneo PIM mit Business Logic Analyse
- [Search Engine Evaluation](./search-engine-evaluation.md) - Algolia vs Meilisearch vs Elasticsearch (Empfehlung: Meilisearch)
- [SAP Integration](./sap-integration.md) - Specification Pattern für SAP-Schnittstellen
- [SAP Interface Analyse](./sap-interface-analysis.md) - Analyse aller 17 SAP-Interfaces aus V2

### Frontend & Portale
- [Support Portal Konzept](./support-portal-concept.md) - Customer Journey Tracking, Alert Dashboard, Trennung von Admin
- [Internationalisierung (i18n)](./i18n-concept.md) - Mehrsprachigkeit ohne fremdsprachige Fragmente
- [Analytics Konzept](./analytics-concept.md) - Kundenverhalten-Analyse, Conversion Funnel, PowerBI Export
- [CMS Konzept](./cms-concept.md) - Headless CMS (Payload), ersetzt Statamic + Typo3
- [KI-Assistenten Konzept](./ai-assistant-concept.md) - Self-hosted LLM, RAG, Sales Assistant, Produktberatung

### Weitere Dokumentation
- [Service-Übersicht](../services/README.md) - Alle Microservices im Detail
- [Observability](../observability/README.md) - Tracing, Logging, Metrics
