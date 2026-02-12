# ADR-001: Provider Pattern für Gondolia Open Source Platform

**Status:** Akzeptiert  
**Datum:** 2026-02-13  
**Autoren:** Architecture Team  
**Tags:** architecture, open-source, extensibility, provider-pattern

---

## Kontext & Problem

Webshop V3 wird als Open Source Projekt **Gondolia** veröffentlicht — eine B2B E-Commerce Platform basierend auf Go Microservices und Next.js Frontend.

### Ausgangslage

- **Identity Service** ist implementiert (Sprint 1+2), weitere Services (Catalog, Cart, Order, Inventory, Shipping, Notification, SAP-Integration) folgen
- Das System integriert **24+ externe Services** (SAP, Akeneo, Saferpay, Azure AD, Meilisearch, Ardis, etc.)
- Es gibt firmenspezifische Logik (company-specific SAP-Konfigurationen, RWDM, Egger Möbelplaner, ComNorm) die **nicht öffentlich** werden darf

### Problem

Wir brauchen eine Architektur die:

1. **Generisch genug** ist für eine Open Source Community
2. **Erweiterbar** ist ohne den Core zu forken
3. **Firmenspezifisches schützt** (Private Repos für proprietäre Integrationen)
4. **Einfach genug** ist, dass externe Entwickler schnell Provider schreiben können
5. **Konsistent** ist mit dem bestehenden Code-Stil (Interface-basierte Repositories, DI via Constructor Injection)

### Identifizierte Integrationspunkte

Aus der Analyse des bestehenden Codes und der Architektur-Dokumente:

| Kategorie | Beispiel-Implementierungen | Anzahl Interfaces |
|-----------|---------------------------|-------------------|
| ERP | SAP R/3, SAP S/4HANA, Microsoft Dynamics | 17 SAP-BAPIs |
| PIM | Akeneo, Pimcore, Salsify | 6 Operationen |
| Payment | Saferpay, Stripe, Adyen | 4 Operationen |
| Search | Meilisearch, Algolia, Elasticsearch | 5 Operationen |
| Auth/SSO | Azure AD (SAML), Keycloak, Auth0 | 3 Flows |
| Storage | Azure Blob, AWS S3, MinIO | 4 Operationen |
| Fulfillment | Ardis, DHL, Swiss Post | 4 Operationen |
| Notification | SMTP, Mailgun, SES, Push | 3 Kanäle |
| CRM | MS Dynamics CRM, Salesforce | 2 Operationen |
| B2B Procurement | RWDM, PunchCommerce, ComNorm | Pro Partner |
| Monitoring | Sentry, OpenTelemetry | Standardisiert |

---

## Untersuchte Alternativen

### 1. Medusa.js (TypeScript)

**Pattern:** Module/Provider System mit DI Container

```typescript
// Medusa: Provider als Module mit DI
class StripePaymentProvider extends AbstractPaymentProvider {
  async initiatePayment(context) { ... }
  async capturePayment(paymentId) { ... }
}

export default {
  services: [StripePaymentProvider],
}
```

**Vorteile:**
- Einfaches Konzept
- Provider als austauschbare Module
- Gute DX (Developer Experience)

**Nachteile:**
- TypeScript-spezifisch (Decorators, DI Container)
- Module-System ist relativ komplex (Loaders, Links, Workflows)
- Kein natives Plugin-Discovery (npm packages)

### 2. Vendure (TypeScript)

**Pattern:** Plugin System mit DI (NestJS)

```typescript
// Vendure: Plugins konfigurieren Provider
@VendurePlugin({
  providers: [MyPaymentHandler],
  configuration: config => {
    config.paymentOptions.paymentMethodHandlers.push(myHandler);
    return config;
  },
})
export class MyPaymentPlugin {}
```

**Vorteile:**
- Mächtiges Plugin-System
- Lifecycle Hooks
- Gute Typisierung

**Nachteile:**
- Sehr NestJS-spezifisch
- Plugin-Konfiguration ist verschachtelt und komplex
- Overhead durch DI-Framework

### 3. Saleor (Python/GraphQL)

**Pattern:** App/Plugin System mit Webhooks

```python
# Saleor: Plugins als separate Apps mit Webhooks
@webhook_handler("ORDER_CREATED")
def handle_order(payload):
    # Provider-Logik
    pass
```

**Vorteile:**
- Lose Kopplung über Webhooks
- Provider können in beliebiger Sprache sein
- Gute Isolation

**Nachteile:**
- Latenz durch HTTP-Calls
- Komplexes Deployment (separate Apps)
- Schwieriger zu debuggen

### 4. Go-nativer Ansatz: Interface + Registry (gewählt)

**Pattern:** Go Interfaces + funktionale Registry + `init()` Registration

```go
// Interfaces im Core, Implementierungen als separate Packages
type PaymentProvider interface {
    Initialize(ctx context.Context, req PaymentRequest) (*PaymentSession, error)
    Capture(ctx context.Context, sessionID string) (*PaymentResult, error)
}

// Registration via init()
func init() {
    provider.Register("payment", "saferpay", NewSaferpayProvider)
}
```

**Vorteile:**
- Idiomatisch Go (Interface Satisfaction, kein Framework-Overhead)
- Compile-Time Type Safety
- Einfaches Wiring via Import-Side-Effects (`_ "github.com/gondolia/gondolia-sap"`)
- Kein DI-Framework nötig
- Passt zum bestehenden Code-Stil (Identity Service nutzt bereits Interface-basierte Repositories)

**Nachteile:**
- Kein dynamisches Plugin-Loading (bewusste Entscheidung: Sicherheit > Dynamik)
- Provider müssen zur Compile-Time bekannt sein

---

## Entscheidung

**Wir verwenden ein Go-natives Provider Pattern mit Interface + Registry + init()-basierter Registration.**

### Begründung

1. **Konsistenz:** Der Identity Service nutzt bereits Interface-basierte Repositories (`UserRepository`, `CompanyRepository`, etc.) mit Constructor Injection. Das Provider Pattern erweitert dieses bewährte Muster auf externe Integrationen.

2. **Einfachheit:** Go Interfaces sind implizit — ein Typ erfüllt ein Interface automatisch wenn er die Methoden implementiert. Kein Registrierungs-Boilerplate wie bei Medusa/Vendure.

3. **Sicherheit:** Compile-Time Verification statt Runtime-Discovery. Keine dynamischen Plugins die Sicherheitslücken öffnen könnten.

4. **Separation of Concerns:** Interfaces im Public Repo, firmenspezifische Implementierungen im Private Repo. Verbindung über Go Module Import.

5. **Community-Freundlich:** Ein externer Entwickler kann einen Provider schreiben indem er einfach ein Go Interface implementiert — keine Framework-Kenntnisse nötig.

---

## Package-Struktur

### Übersicht

```
github.com/gondolia/gondolia              (Public — Core Platform)
├── provider/                              # Provider-Interfaces & Registry
│   ├── provider.go                        # Basis-Typen, Registry
│   ├── erp/                               # ERPProvider Interface
│   │   ├── erp.go                         # Interface-Definition
│   │   └── noop/                          # Noop-Implementierung
│   │       └── noop.go
│   ├── pim/                               # PIMProvider Interface
│   │   ├── pim.go
│   │   └── noop/
│   │       └── noop.go
│   ├── payment/                           # PaymentProvider Interface
│   │   ├── payment.go
│   │   └── noop/
│   │       └── noop.go
│   ├── search/                            # SearchProvider Interface
│   │   ├── search.go
│   │   └── noop/
│   │       └── noop.go
│   ├── auth/                              # AuthProvider Interface (SSO)
│   │   ├── auth.go
│   │   └── noop/
│   │       └── noop.go
│   ├── storage/                           # StorageProvider Interface
│   │   ├── storage.go
│   │   └── noop/
│   │       └── noop.go
│   ├── fulfillment/                       # FulfillmentProvider Interface
│   │   ├── fulfillment.go
│   │   └── noop/
│   │       └── noop.go
│   ├── notification/                      # NotificationProvider Interface
│   │   ├── notification.go
│   │   └── noop/
│   │       └── noop.go
│   ├── crm/                               # CRMProvider Interface
│   │   ├── crm.go
│   │   └── noop/
│   │       └── noop.go
│   └── tax/                               # TaxProvider Interface
│       ├── tax.go
│       └── noop/
│           └── noop.go
├── services/                              # Microservices
│   ├── identity/                          # ✅ Implementiert
│   ├── catalog/                           # ⏳ Sprint 3
│   ├── cart/
│   ├── order/
│   ├── inventory/
│   ├── shipping/
│   ├── notification/
│   └── gateway/
├── pkg/                                   # Shared Libraries
│   ├── errors/
│   ├── logging/
│   └── telemetry/
└── docs/

github.com/gondolia/gondolia-sap          (Public — SAP ERP Provider)
├── provider.go                            # SAP ERPProvider Implementierung
├── spec/                                  # Specification Pattern Framework
├── specs/                                 # Konkrete SAP-Spezifikationen
├── transformer/                           # SAP-Datentransformationen
└── internal/
    ├── soap/                              # SOAP Client
    └── odata/                             # OData Client (S/4HANA)

github.com/gondolia/gondolia-meilisearch   (Public — Meilisearch Provider)
├── provider.go                            # Meilisearch SearchProvider
└── internal/
    └── client/

github.com/gondolia/gondolia-saferpay      (Public — Saferpay Provider)
├── provider.go                            # Saferpay PaymentProvider
└── internal/
    └── client/

github.com/mycompany/gondolia-kuratle        (Private — company-specific)
├── providers.go                           # Registration aller Acme-Provider
├── sap/                                   # SAP-Konfiguration (VKOrg, Werke, etc.)
├── akeneo/                                # Akeneo PIMProvider + Attribut-Mapping
├── azure/                                 # Azure AD AuthProvider
├── ardis/                                 # Ardis FulfillmentProvider
├── rwdm/                                  # RWDM B2B-Integration
├── egger/                                 # Egger Möbelplaner
└── config/                                # company-specific Konfiguration
```

### Wo liegt was?

| Artefakt | Package | Begründung |
|----------|---------|------------|
| **Interfaces** | `gondolia/provider/{type}/` | Im Core, damit alle Provider dagegen implementieren |
| **Noop-Implementierungen** | `gondolia/provider/{type}/noop/` | Defaultwerte für nicht-konfigurierte Provider |
| **Registry** | `gondolia/provider/` | Zentral, von allen Services genutzt |
| **Generische Provider** (SAP, Meilisearch, Saferpay) | Eigene Public Repos | Wiederverwendbar, unabhängig versioniert |
| **Firmenspezifisches** | Private Repo | Geschützt, enthält Konfiguration + Custom Logic |

---

## Konkrete Interface-Definitionen

### Basis-Typen & Registry

```go
// provider/provider.go
package provider

import (
    "context"
    "fmt"
    "sync"
)

// Metadata beschreibt einen Provider
type Metadata struct {
    Name        string            // z.B. "saferpay"
    DisplayName string            // z.B. "Saferpay (SIX Payment Services)"
    Category    string            // z.B. "payment"
    Version     string            // z.B. "1.0.0"
    Description string
    ConfigSpec  []ConfigField     // Welche Konfiguration der Provider braucht
}

// ConfigField beschreibt ein Konfigurationsfeld
type ConfigField struct {
    Key         string
    Type        string // "string", "int", "bool", "secret"
    Required    bool
    Default     any
    Description string
}

// ProviderFactory erstellt eine Provider-Instanz aus Konfiguration
type ProviderFactory[T any] func(config map[string]any) (T, error)

// --- Global Registry ---

var (
    registry = make(map[string]map[string]any) // category -> name -> factory
    metadata = make(map[string]map[string]Metadata)
    mu       sync.RWMutex
)

// Register registriert eine Provider-Factory
func Register[T any](category, name string, meta Metadata, factory ProviderFactory[T]) {
    mu.Lock()
    defer mu.Unlock()

    if registry[category] == nil {
        registry[category] = make(map[string]any)
        metadata[category] = make(map[string]Metadata)
    }
    registry[category][name] = factory
    metadata[category][name] = meta
}

// Get holt eine Provider-Factory
func Get[T any](category, name string) (ProviderFactory[T], error) {
    mu.RLock()
    defer mu.RUnlock()

    cat, ok := registry[category]
    if !ok {
        return nil, fmt.Errorf("unknown provider category: %s", category)
    }
    factory, ok := cat[name]
    if !ok {
        return nil, fmt.Errorf("unknown provider: %s.%s", category, name)
    }
    f, ok := factory.(ProviderFactory[T])
    if !ok {
        return nil, fmt.Errorf("provider %s.%s has wrong type", category, name)
    }
    return f, nil
}

// List gibt alle registrierten Provider einer Kategorie zurück
func List(category string) []Metadata {
    mu.RLock()
    defer mu.RUnlock()

    var result []Metadata
    for _, m := range metadata[category] {
        result = append(result, m)
    }
    return result
}

// ListAll gibt alle registrierten Provider zurück
func ListAll() map[string][]Metadata {
    mu.RLock()
    defer mu.RUnlock()

    result := make(map[string][]Metadata)
    for cat, providers := range metadata {
        for _, m := range providers {
            result[cat] = append(result[cat], m)
        }
    }
    return result
}
```

### ERPProvider

```go
// provider/erp/erp.go
package erp

import (
    "context"
    "time"
)

// ERPProvider abstrahiert die Kommunikation mit einem ERP-System (SAP, Dynamics, etc.)
type ERPProvider interface {
    // --- Order Management ---

    // CreateOrder überträgt eine Bestellung an das ERP-System
    CreateOrder(ctx context.Context, req CreateOrderRequest) (*CreateOrderResult, error)

    // SimulateOrder berechnet Preise und Verfügbarkeit ohne Commit
    SimulateOrder(ctx context.Context, req SimulateOrderRequest) (*SimulateOrderResult, error)

    // GetOrderStatus ruft den Bestellstatus ab
    GetOrderStatus(ctx context.Context, orderID string) (*OrderStatus, error)

    // --- Inventory ---

    // GetProductAvailability ruft Lagerbestände ab
    GetProductAvailability(ctx context.Context, skus []string) ([]ProductStock, error)

    // --- Pricing ---

    // GetTierPrices ruft Staffelpreise für Kunde/Produkt ab
    GetTierPrices(ctx context.Context, req TierPriceRequest) ([]TierPrice, error)

    // --- Company/Customer Data ---

    // SyncCompany synchronisiert Firmendaten vom ERP
    SyncCompany(ctx context.Context, erpCustomerID string) (*CompanyData, error)

    // GetCompanyAddresses ruft Adressen einer Firma ab
    GetCompanyAddresses(ctx context.Context, erpCustomerID string) ([]Address, error)

    // --- Reports ---

    // GetOrderHistory ruft Bestellhistorie ab
    GetOrderHistory(ctx context.Context, req ReportFilter) ([]OrderReport, error)

    // GetShipmentHistory ruft Lieferhistorie ab
    GetShipmentHistory(ctx context.Context, req ReportFilter) ([]ShipmentReport, error)

    // GetInvoiceHistory ruft Rechnungshistorie ab
    GetInvoiceHistory(ctx context.Context, req ReportFilter) ([]InvoiceReport, error)

    // --- Metadata ---

    // Metadata gibt Informationen über den Provider zurück
    Metadata() Metadata
}

// --- Request/Response Types ---

type CreateOrderRequest struct {
    TenantConfig TenantConfig
    Order        Order
    Customer     Customer
    ShipTo       Address
    BillTo       Address
}

type CreateOrderResult struct {
    ERPOrderNumber string
    Items          []OrderItemResult
    Messages       []Message
}

type SimulateOrderRequest struct {
    TenantConfig TenantConfig
    Items        []SimulateItem
    Customer     Customer
    ShipTo       Address
    DesiredDate  *time.Time
}

type SimulateOrderResult struct {
    Items    []SimulatedItem
    Totals   Totals
    Schedule []DeliverySchedule
    Messages []Message
}

type OrderStatus struct {
    ERPOrderNumber string
    Status         string
    Items          []OrderItemStatus
}

type ProductStock struct {
    SKU          string
    PlantCode    string
    PlantName    string
    Quantity     float64
    Unit         string
    LeadTimeDays int
}

type TierPriceRequest struct {
    CustomerID string
    SKUs       []string
    Currency   string
}

type TierPrice struct {
    SKU       string
    MinQty    float64
    Price     float64
    Currency  string
    ValidFrom time.Time
    ValidTo   time.Time
}

type ReportFilter struct {
    CustomerID string
    DateFrom   time.Time
    DateTo     time.Time
    Limit      int
    Offset     int
}

type TenantConfig struct {
    SalesOrg    string // SAP: VKORG
    DistChannel string // SAP: VTWEG
    Division    string // SAP: SPART
    Currency    string
    Language    string
}

type Metadata struct {
    Name        string
    Version     string
    Protocol    string // "soap", "odata", "rest"
    Capabilities []string
}

// Message repräsentiert eine ERP-Nachricht
type Message struct {
    Type    string // "success", "warning", "error", "info"
    Code    string
    Message string
}
```

### PIMProvider

```go
// provider/pim/pim.go
package pim

import (
    "context"
    "io"
    "time"
)

// PIMProvider abstrahiert die Kommunikation mit einem PIM-System (Akeneo, Pimcore, etc.)
type PIMProvider interface {
    // --- Products ---

    // FetchProducts ruft Produkte ab (mit Cursor-Pagination)
    FetchProducts(ctx context.Context, filter ProductFilter) (*ProductPage, error)

    // FetchProduct ruft ein einzelnes Produkt ab
    FetchProduct(ctx context.Context, identifier string) (*Product, error)

    // --- Categories ---

    // FetchCategories ruft Kategorien ab
    FetchCategories(ctx context.Context) ([]Category, error)

    // --- Attributes ---

    // FetchAttributes ruft Attribut-Definitionen ab
    FetchAttributes(ctx context.Context) ([]Attribute, error)

    // --- Assets/Media ---

    // DownloadAsset lädt eine Asset-Datei herunter
    DownloadAsset(ctx context.Context, assetCode string) (io.ReadCloser, string, error)

    // --- Metadata ---

    Metadata() Metadata
}

type ProductFilter struct {
    UpdatedSince *time.Time
    Families     []string
    Categories   []string
    Cursor       string // Für Pagination
    Limit        int
}

type ProductPage struct {
    Products   []Product
    NextCursor string
    TotalCount int
}

type Product struct {
    Identifier string
    Family     string
    Categories []string
    Enabled    bool
    Values     map[string][]AttributeValue // Attribut-Name -> lokalisierte Werte
    Created    time.Time
    Updated    time.Time
}

type AttributeValue struct {
    Locale string
    Scope  string
    Data   any
}

type Category struct {
    Code   string
    Parent string
    Labels map[string]string // Locale -> Label
}

type Attribute struct {
    Code          string
    Type          string // "text", "number", "select", "media", etc.
    Group         string
    Localizable   bool
    Scopable      bool
    Labels        map[string]string
}

type Metadata struct {
    Name    string
    Version string
}
```

### PaymentProvider

```go
// provider/payment/payment.go
package payment

import (
    "context"
)

// PaymentProvider abstrahiert Payment-Gateways (Saferpay, Stripe, Adyen, etc.)
type PaymentProvider interface {
    // Initialize startet eine Payment-Session
    Initialize(ctx context.Context, req InitializeRequest) (*PaymentSession, error)

    // Authorize prüft und autorisiert eine Zahlung
    Authorize(ctx context.Context, sessionID string) (*AuthorizationResult, error)

    // Capture erfasst eine autorisierte Zahlung
    Capture(ctx context.Context, transactionID string, amount *Amount) (*CaptureResult, error)

    // Cancel storniert eine autorisierte Zahlung
    Cancel(ctx context.Context, transactionID string) error

    // Refund erstattet eine erfasste Zahlung (teil- oder vollständig)
    Refund(ctx context.Context, transactionID string, amount Amount) (*RefundResult, error)

    // HandleWebhook verarbeitet Provider-spezifische Webhooks
    HandleWebhook(ctx context.Context, payload []byte, headers map[string]string) (*WebhookEvent, error)

    // Metadata gibt Provider-Informationen zurück
    Metadata() Metadata
}

type InitializeRequest struct {
    OrderID        string
    Amount         Amount
    Currency       string
    Description    string
    ReturnURL      string
    WebhookURL     string
    PaymentMethods []string // z.B. ["VISA", "MASTERCARD", "TWINT"]
    CustomerEmail  string
    Metadata       map[string]string
}

type PaymentSession struct {
    SessionID   string
    RedirectURL string
    Token       string
    ExpiresAt   string
}

type AuthorizationResult struct {
    TransactionID  string
    Status         string // "authorized", "failed", "pending"
    Amount         Amount
    PaymentMethod  string
    CardDisplay    string // "xxxx xxxx xxxx 1234"
    LiabilityShift bool
    Raw            map[string]any // Provider-spezifische Rohdaten
}

type CaptureResult struct {
    TransactionID string
    Status        string
    CapturedAt    string
}

type RefundResult struct {
    RefundID string
    Status   string
    Amount   Amount
}

type Amount struct {
    Value    int64  // Cents
    Currency string // ISO 4217
}

type WebhookEvent struct {
    Type          string // "payment.authorized", "payment.captured", "payment.failed"
    TransactionID string
    OrderID       string
    Data          map[string]any
}

type Metadata struct {
    Name             string
    SupportedMethods []string
    TestMode         bool
}
```

### SearchProvider

```go
// provider/search/search.go
package search

import "context"

// SearchProvider abstrahiert Suchmaschinen (Meilisearch, Algolia, Elasticsearch, etc.)
type SearchProvider interface {
    // --- Indexing ---

    // IndexDocuments indexiert Dokumente in einen Index
    IndexDocuments(ctx context.Context, index string, documents []Document) (*TaskResult, error)

    // DeleteDocuments löscht Dokumente aus einem Index
    DeleteDocuments(ctx context.Context, index string, ids []string) (*TaskResult, error)

    // ConfigureIndex konfiguriert Index-Einstellungen
    ConfigureIndex(ctx context.Context, index string, config IndexConfig) error

    // --- Search ---

    // Search führt eine Suche aus
    Search(ctx context.Context, index string, query SearchQuery) (*SearchResult, error)

    // --- Management ---

    // CreateIndex erstellt einen neuen Index
    CreateIndex(ctx context.Context, index string, primaryKey string) error

    // DeleteIndex löscht einen Index
    DeleteIndex(ctx context.Context, index string) error

    // GetTaskStatus prüft den Status eines asynchronen Tasks
    GetTaskStatus(ctx context.Context, taskID string) (*TaskResult, error)

    // Health prüft die Verfügbarkeit der Suchmaschine
    Health(ctx context.Context) error

    // Metadata gibt Provider-Informationen zurück
    Metadata() Metadata
}

type Document map[string]any

type SearchQuery struct {
    Query      string
    Filters    []Filter
    Facets     []string
    Sort       []string
    Offset     int
    Limit      int
    Highlight  []string
}

type Filter struct {
    Field    string
    Operator string // "=", "!=", ">", "<", ">=", "<=", "IN", "NOT IN"
    Value    any
}

type SearchResult struct {
    Hits             []Document
    TotalHits        int
    Facets           map[string]map[string]int
    ProcessingTimeMs int
}

type IndexConfig struct {
    SearchableAttributes []string
    FilterableAttributes []string
    SortableAttributes   []string
    Synonyms             map[string][]string
    StopWords            []string
    TypoTolerance        *TypoTolerance
}

type TypoTolerance struct {
    Enabled bool
    MinWordSizeForTypos map[string]int
}

type TaskResult struct {
    TaskID string
    Status string // "enqueued", "processing", "succeeded", "failed"
    Error  string
}

type Metadata struct {
    Name     string
    Version  string
    Features []string // z.B. ["facets", "typo-tolerance", "synonyms", "geo-search"]
}
```

### AuthProvider (SSO)

```go
// provider/auth/auth.go
package auth

import "context"

// AuthProvider abstrahiert SSO/Identity Provider (Azure AD, Keycloak, Auth0, etc.)
type AuthProvider interface {
    // GetAuthURL gibt die URL für den SSO-Login zurück
    GetAuthURL(ctx context.Context, state string, redirectURL string) (string, error)

    // HandleCallback verarbeitet den SSO-Callback und gibt User-Infos zurück
    HandleCallback(ctx context.Context, code string, state string) (*SSOUser, error)

    // ValidateToken validiert einen SSO-Token (für API-basierte Flows)
    ValidateToken(ctx context.Context, token string) (*SSOUser, error)

    // GetUserInfo ruft Benutzerinformationen vom Provider ab
    GetUserInfo(ctx context.Context, accessToken string) (*SSOUser, error)

    // Metadata gibt Provider-Informationen zurück
    Metadata() Metadata
}

type SSOUser struct {
    ExternalID string
    Email      string
    FirstName  string
    LastName   string
    Groups     []string
    Attributes map[string]string
    RawClaims  map[string]any
}

type Metadata struct {
    Name     string
    Protocol string // "saml", "oidc", "oauth2"
    Issuer   string
}
```

### StorageProvider

```go
// provider/storage/storage.go
package storage

import (
    "context"
    "io"
    "time"
)

// StorageProvider abstrahiert Object Storage (S3, Azure Blob, MinIO, lokales Filesystem)
type StorageProvider interface {
    // Upload lädt eine Datei hoch
    Upload(ctx context.Context, path string, reader io.Reader, opts UploadOptions) (*FileInfo, error)

    // Download lädt eine Datei herunter
    Download(ctx context.Context, path string) (io.ReadCloser, *FileInfo, error)

    // Delete löscht eine Datei
    Delete(ctx context.Context, path string) error

    // Exists prüft ob eine Datei existiert
    Exists(ctx context.Context, path string) (bool, error)

    // GetSignedURL gibt eine temporäre URL für direkten Zugriff zurück
    GetSignedURL(ctx context.Context, path string, expiry time.Duration) (string, error)

    // List listet Dateien in einem Pfad
    List(ctx context.Context, prefix string, opts ListOptions) ([]FileInfo, error)

    // Metadata gibt Provider-Informationen zurück
    Metadata() Metadata
}

type UploadOptions struct {
    ContentType string
    Metadata    map[string]string
    ACL         string // "private", "public-read"
}

type ListOptions struct {
    MaxKeys int
    Cursor  string
}

type FileInfo struct {
    Path         string
    Size         int64
    ContentType  string
    LastModified time.Time
    Metadata     map[string]string
    ETag         string
}

type Metadata struct {
    Name   string
    Region string
    Bucket string
}
```

### FulfillmentProvider

```go
// provider/fulfillment/fulfillment.go
package fulfillment

import (
    "context"
    "time"
)

// FulfillmentProvider abstrahiert Logistik/Versand-Systeme (Ardis, DHL, Swiss Post, etc.)
type FulfillmentProvider interface {
    // CreateShipment erstellt einen Versandauftrag
    CreateShipment(ctx context.Context, req ShipmentRequest) (*ShipmentResult, error)

    // GetShipmentStatus ruft den Versandstatus ab
    GetShipmentStatus(ctx context.Context, shipmentID string) (*ShipmentStatus, error)

    // CancelShipment storniert einen Versandauftrag
    CancelShipment(ctx context.Context, shipmentID string) error

    // GetTrackingURL gibt die Tracking-URL zurück
    GetTrackingURL(ctx context.Context, trackingNumber string) (string, error)

    // CalculateShipping berechnet Versandkosten
    CalculateShipping(ctx context.Context, req ShippingCalcRequest) ([]ShippingOption, error)

    // Metadata gibt Provider-Informationen zurück
    Metadata() Metadata
}

type ShipmentRequest struct {
    OrderID  string
    From     Address
    To       Address
    Packages []Package
    Service  string // z.B. "standard", "express"
    Notes    string
}

type ShipmentResult struct {
    ShipmentID     string
    TrackingNumber string
    Label          []byte // PDF Label
    EstimatedDate  *time.Time
}

type ShipmentStatus struct {
    ShipmentID     string
    Status         string // "created", "picked_up", "in_transit", "delivered", "failed"
    TrackingNumber string
    Events         []TrackingEvent
}

type TrackingEvent struct {
    Timestamp   time.Time
    Status      string
    Location    string
    Description string
}

type ShippingCalcRequest struct {
    From     Address
    To       Address
    Packages []Package
}

type ShippingOption struct {
    Service       string
    Name          string
    Price         int64
    Currency      string
    EstimatedDays int
}

type Address struct {
    Name       string
    Street     string
    PostalCode string
    City       string
    Country    string // ISO 3166-1 alpha-2
}

type Package struct {
    WeightGrams int
    LengthCm    int
    WidthCm     int
    HeightCm    int
}

type Metadata struct {
    Name     string
    Carriers []string
}
```

### NotificationProvider

```go
// provider/notification/notification.go
package notification

import "context"

// NotificationProvider abstrahiert Benachrichtigungskanäle (Email, SMS, Push, etc.)
type NotificationProvider interface {
    // Send sendet eine Benachrichtigung
    Send(ctx context.Context, msg Message) (*SendResult, error)

    // SendBatch sendet mehrere Benachrichtigungen
    SendBatch(ctx context.Context, msgs []Message) ([]SendResult, error)

    // Channels gibt die unterstützten Kanäle zurück
    Channels() []string

    // Metadata gibt Provider-Informationen zurück
    Metadata() Metadata
}

type Message struct {
    Channel    string // "email", "sms", "push"
    To         []Recipient
    From       string
    Subject    string
    Body       string // HTML für Email, Plain-Text für SMS/Push
    BodyPlain  string // Plain-Text Fallback für Email
    Template   string // Template-ID (optional)
    TemplateData map[string]any
    Metadata   map[string]string
    Attachments []Attachment
}

type Recipient struct {
    Address string // Email, Telefonnummer oder Device-Token
    Name    string
}

type Attachment struct {
    Filename    string
    ContentType string
    Data        []byte
}

type SendResult struct {
    MessageID string
    Status    string // "sent", "queued", "failed"
    Error     string
}

type Metadata struct {
    Name     string
    Channels []string
}
```

### CRMProvider

```go
// provider/crm/crm.go
package crm

import "context"

// CRMProvider abstrahiert CRM-Systeme (MS Dynamics, Salesforce, etc.)
type CRMProvider interface {
    // SyncContact synchronisiert einen Kontakt mit dem CRM
    SyncContact(ctx context.Context, contact Contact) (*SyncResult, error)

    // GetAccount ruft Firmendaten aus dem CRM ab
    GetAccount(ctx context.Context, accountID string) (*Account, error)

    // ListAccounts ruft eine Liste von Firmen ab
    ListAccounts(ctx context.Context, filter AccountFilter) ([]Account, error)

    // Metadata gibt Provider-Informationen zurück
    Metadata() Metadata
}

type Contact struct {
    ExternalID string
    Email      string
    FirstName  string
    LastName   string
    Phone      string
    Company    string
    Attributes map[string]any
}

type Account struct {
    ExternalID    string
    Name          string
    ERPCustomerID string
    Attributes    map[string]any
}

type AccountFilter struct {
    Query  string
    Limit  int
    Offset int
}

type SyncResult struct {
    ExternalID string
    Action     string // "created", "updated", "unchanged"
}

type Metadata struct {
    Name string
}
```

### TaxProvider

```go
// provider/tax/tax.go
package tax

import "context"

// TaxProvider abstrahiert Steuerberechnungen (intern, Avalara, TaxJar, etc.)
type TaxProvider interface {
    // CalculateTax berechnet Steuern für eine Liste von Positionen
    CalculateTax(ctx context.Context, req TaxRequest) (*TaxResult, error)

    // Metadata gibt Provider-Informationen zurück
    Metadata() Metadata
}

type TaxRequest struct {
    Country  string // ISO 3166-1 alpha-2
    Region   string
    Currency string
    Items    []TaxItem
}

type TaxItem struct {
    SKU       string
    Quantity  float64
    UnitPrice int64 // Cents
    TaxCode   string
}

type TaxResult struct {
    Items    []TaxItemResult
    TotalTax int64
}

type TaxItemResult struct {
    SKU      string
    TaxRate  float64 // z.B. 0.081 für 8.1%
    TaxAmount int64
}

type Metadata struct {
    Name string
}
```

---

## Wiring & Registration

### Wie registriert ein Provider sich?

Jeder Provider registriert sich via `init()`:

```go
// github.com/gondolia/gondolia-saferpay/provider.go
package saferpay

import (
    "github.com/gondolia/gondolia/provider"
    "github.com/gondolia/gondolia/provider/payment"
)

func init() {
    provider.Register[payment.PaymentProvider]("payment", "saferpay",
        provider.Metadata{
            Name:        "saferpay",
            DisplayName: "Saferpay (SIX Payment Services)",
            Category:    "payment",
            Version:     "1.0.0",
            Description: "Payment gateway for Switzerland and Europe",
            ConfigSpec: []provider.ConfigField{
                {Key: "customer_id", Type: "string", Required: true, Description: "Saferpay Customer ID"},
                {Key: "terminal_id", Type: "string", Required: true, Description: "Saferpay Terminal ID"},
                {Key: "api_secret", Type: "secret", Required: true, Description: "Saferpay API Secret"},
                {Key: "test_mode", Type: "bool", Default: false, Description: "Use Saferpay test environment"},
            },
        },
        NewProvider,
    )
}

// NewProvider erstellt eine neue Saferpay-Instanz
func NewProvider(config map[string]any) (payment.PaymentProvider, error) {
    customerID, _ := config["customer_id"].(string)
    terminalID, _ := config["terminal_id"].(string)
    apiSecret, _ := config["api_secret"].(string)
    testMode, _ := config["test_mode"].(bool)

    if customerID == "" || terminalID == "" || apiSecret == "" {
        return nil, fmt.Errorf("saferpay: customer_id, terminal_id and api_secret are required")
    }

    baseURL := "https://www.saferpay.com/api"
    if testMode {
        baseURL = "https://test.saferpay.com/api"
    }

    return &Provider{
        customerID: customerID,
        terminalID: terminalID,
        apiSecret:  apiSecret,
        baseURL:    baseURL,
        client:     &http.Client{Timeout: 30 * time.Second},
    }, nil
}
```

### Wie bindet ein Private Repo seine Provider ein?

Das Private Repo importiert Provider-Packages und die `init()`-Funktionen registrieren automatisch:

```go
// github.com/mycompany/gondolia-kuratle/providers.go
package acme

// Blank Imports registrieren alle company-specificn Provider
import (
    _ "github.com/gondolia/gondolia-sap"         // SAP ERPProvider
    _ "github.com/gondolia/gondolia-meilisearch"  // Meilisearch SearchProvider
    _ "github.com/gondolia/gondolia-saferpay"     // Saferpay PaymentProvider

    _ "github.com/mycompany/gondolia-kuratle/akeneo" // company-specificr Akeneo Provider
    _ "github.com/mycompany/gondolia-kuratle/azure"  // Acme Azure AD Provider
    _ "github.com/mycompany/gondolia-kuratle/ardis"  // Acme Ardis Provider
)
```

### Wie wird ein Provider in einem Service genutzt?

```go
// services/order/cmd/server/main.go
package main

import (
    "github.com/gondolia/gondolia/provider"
    "github.com/gondolia/gondolia/provider/erp"
    "github.com/gondolia/gondolia/provider/payment"

    // Provider-Registrierung via Blank Import
    _ "github.com/gondolia/gondolia-sap"
    _ "github.com/gondolia/gondolia-saferpay"
)

func main() {
    cfg := loadConfig()

    // ERP Provider erstellen
    erpFactory, err := provider.Get[erp.ERPProvider]("erp", cfg.ERPProvider)
    if err != nil {
        log.Fatal("ERP provider not found:", err)
    }
    erpProvider, err := erpFactory(cfg.ERPConfig)
    if err != nil {
        log.Fatal("Failed to create ERP provider:", err)
    }

    // Payment Provider erstellen
    paymentFactory, err := provider.Get[payment.PaymentProvider]("payment", cfg.PaymentProvider)
    if err != nil {
        log.Fatal("Payment provider not found:", err)
    }
    paymentProvider, err := paymentFactory(cfg.PaymentConfig)
    if err != nil {
        log.Fatal("Failed to create payment provider:", err)
    }

    // Services erstellen (Dependency Injection wie im Identity Service)
    orderService := service.NewOrderService(
        orderRepo,
        erpProvider,
        paymentProvider,
        // ...
    )

    // ...
}
```

### Konfiguration

```yaml
# config.yaml
providers:
  erp:
    name: sap
    config:
      endpoint: "https://sap.example.com:443/sap/bc/srt/rfc/sap"
      client: "100"
      username: "${SAP_USERNAME}"
      password: "${SAP_PASSWORD}"
      sales_org: "1000"
      dist_channel: "10"
      division: "00"

  payment:
    name: saferpay
    config:
      customer_id: "${SAFERPAY_CUSTOMER_ID}"
      terminal_id: "${SAFERPAY_TERMINAL_ID}"
      api_secret: "${SAFERPAY_API_SECRET}"
      test_mode: false

  search:
    name: meilisearch
    config:
      url: "http://meilisearch:7700"
      master_key: "${MEILI_MASTER_KEY}"

  storage:
    name: s3
    config:
      region: "eu-central-1"
      bucket: "gondolia-assets"
      access_key: "${AWS_ACCESS_KEY}"
      secret_key: "${AWS_SECRET_KEY}"

  notification:
    name: smtp
    config:
      host: "smtp.example.com"
      port: 587
      username: "${SMTP_USERNAME}"
      password: "${SMTP_PASSWORD}"
      from: "shop@example.com"
```

---

## Beispiel: Einen neuen Payment-Provider schreiben

So würde ein externer Entwickler einen Stripe-Provider für Gondolia schreiben:

### 1. Neues Go-Modul erstellen

```bash
mkdir gondolia-stripe
cd gondolia-stripe
go mod init github.com/myname/gondolia-stripe
go get github.com/gondolia/gondolia
go get github.com/stripe/stripe-go/v76
```

### 2. Provider implementieren

```go
// provider.go
package stripe

import (
    "context"
    "fmt"

    "github.com/gondolia/gondolia/provider"
    "github.com/gondolia/gondolia/provider/payment"
    stripego "github.com/stripe/stripe-go/v76"
    "github.com/stripe/stripe-go/v76/checkout/session"
)

func init() {
    provider.Register[payment.PaymentProvider]("payment", "stripe",
        provider.Metadata{
            Name:        "stripe",
            DisplayName: "Stripe",
            Category:    "payment",
            Version:     "1.0.0",
            Description: "Stripe payment gateway",
            ConfigSpec: []provider.ConfigField{
                {Key: "secret_key", Type: "secret", Required: true},
                {Key: "webhook_secret", Type: "secret", Required: true},
            },
        },
        NewProvider,
    )
}

type Provider struct {
    secretKey     string
    webhookSecret string
}

func NewProvider(config map[string]any) (payment.PaymentProvider, error) {
    secretKey, _ := config["secret_key"].(string)
    if secretKey == "" {
        return nil, fmt.Errorf("stripe: secret_key is required")
    }
    stripego.Key = secretKey

    webhookSecret, _ := config["webhook_secret"].(string)
    return &Provider{
        secretKey:     secretKey,
        webhookSecret: webhookSecret,
    }, nil
}

func (p *Provider) Initialize(ctx context.Context, req payment.InitializeRequest) (*payment.PaymentSession, error) {
    params := &stripego.CheckoutSessionParams{
        Mode: stripego.String(string(stripego.CheckoutSessionModePayment)),
        LineItems: []*stripego.CheckoutSessionLineItemParams{
            {
                PriceData: &stripego.CheckoutSessionLineItemPriceDataParams{
                    Currency:   stripego.String(req.Currency),
                    UnitAmount: stripego.Int64(req.Amount.Value),
                    ProductData: &stripego.CheckoutSessionLineItemPriceDataProductDataParams{
                        Name: stripego.String(req.Description),
                    },
                },
                Quantity: stripego.Int64(1),
            },
        },
        SuccessURL: stripego.String(req.ReturnURL),
        CancelURL:  stripego.String(req.ReturnURL),
    }

    s, err := session.New(params)
    if err != nil {
        return nil, fmt.Errorf("stripe: create session: %w", err)
    }

    return &payment.PaymentSession{
        SessionID:   s.ID,
        RedirectURL: s.URL,
    }, nil
}

func (p *Provider) Authorize(ctx context.Context, sessionID string) (*payment.AuthorizationResult, error) {
    s, err := session.Get(sessionID, nil)
    if err != nil {
        return nil, fmt.Errorf("stripe: get session: %w", err)
    }

    status := "failed"
    if s.PaymentStatus == stripego.CheckoutSessionPaymentStatusPaid {
        status = "authorized"
    }

    return &payment.AuthorizationResult{
        TransactionID: s.PaymentIntent.ID,
        Status:        status,
        Amount:        payment.Amount{Value: s.AmountTotal, Currency: string(s.Currency)},
        PaymentMethod: string(s.PaymentMethodTypes[0]),
    }, nil
}

func (p *Provider) Capture(ctx context.Context, transactionID string, amount *payment.Amount) (*payment.CaptureResult, error) {
    // Stripe auto-captures by default
    return &payment.CaptureResult{
        TransactionID: transactionID,
        Status:        "captured",
    }, nil
}

func (p *Provider) Cancel(ctx context.Context, transactionID string) error {
    // Cancel via PaymentIntent
    return nil
}

func (p *Provider) Refund(ctx context.Context, transactionID string, amount payment.Amount) (*payment.RefundResult, error) {
    // Stripe Refund API
    return nil, nil
}

func (p *Provider) HandleWebhook(ctx context.Context, payload []byte, headers map[string]string) (*payment.WebhookEvent, error) {
    // Verify and parse Stripe webhook
    return nil, nil
}

func (p *Provider) Metadata() payment.Metadata {
    return payment.Metadata{
        Name:             "stripe",
        SupportedMethods: []string{"card", "sepa_debit", "ideal", "bancontact", "giropay"},
        TestMode:         false,
    }
}
```

### 3. Verwenden

```go
// In der Applikation:
import _ "github.com/myname/gondolia-stripe"

// Das war's! Der Provider ist jetzt über die Registry verfügbar.
// In der config.yaml:
// providers:
//   payment:
//     name: stripe
//     config:
//       secret_key: "${STRIPE_SECRET_KEY}"
//       webhook_secret: "${STRIPE_WEBHOOK_SECRET}"
```

---

## Konsequenzen

### Positive Konsequenzen

1. **Klare Trennung Public/Private:** Interfaces und generische Provider sind Open Source, firmenspezifische Implementierungen bleiben privat
2. **Einfache Erweiterbarkeit:** Neue Provider = neues Go-Modul das ein Interface implementiert
3. **Type Safety:** Compile-Time Verification aller Provider-Verträge
4. **Konsistenz:** Alle Services nutzen dasselbe Pattern für externe Integrationen
5. **Testbarkeit:** Noop-Provider für Tests, Interface-basierte Mocks
6. **Community-freundlich:** Kein Framework-Wissen nötig, Standard-Go-Patterns
7. **Backward-Compatible mit bestehendem Code:** Der Identity Service nutzt bereits Interface-basierte Repositories — das Provider Pattern erweitert dieses Muster

### Negative Konsequenzen

1. **Kein dynamisches Plugin-Loading:** Provider müssen zur Compile-Time bekannt sein (bewusste Entscheidung)
2. **Blank Imports:** `_ "github.com/..."` ist nicht sofort intuitiv für Go-Neulinge
3. **Global Registry:** Verwendet Package-Level State (`init()`), was in sehr großen Systemen unübersichtlich werden könnte
4. **Generics-Nutzung:** `provider.Register[T]` und `provider.Get[T]` erfordern Go 1.18+ (kein Problem für neues Projekt)

### Risiken

| Risiko | Mitigation |
|--------|------------|
| Interfaces zu breit/eng | Review durch Community, iterative Anpassung, Breaking Changes nur in Major Versions |
| Registry-Konflikte bei gleichen Namen | Namenskonvention: `{category}.{vendor}.{name}`, Panic bei Duplikaten |
| Performance durch Interface-Indirection | Minimal, Go Interfaces sind extrem effizient |
| Zu viele kleine Repos | Standardisiertes Repo-Template, Monorepo für offizielle Provider als Alternative |

### Migration vom bestehenden Code

Der bestehende Identity Service muss **nicht** geändert werden. Das Provider Pattern wird für neue Services (ab Catalog Service, Sprint 3) eingeführt. Bestehende Architektur-Konzepte (SAP Specification Pattern, Datahub, etc.) werden als Provider-Implementierungen abgebildet:

- `sap-integration.md` Specification Pattern → `gondolia-sap` ERPProvider
- `external-integrations-analysis.md` Provider Registry → `provider/` Package
- `akeneo-import-analysis.md` Import Pipeline → `gondolia-akeneo` PIMProvider
- `search-engine-evaluation.md` Meilisearch → `gondolia-meilisearch` SearchProvider

---

## Nächste Schritte

1. [ ] `provider/` Package mit Registry und allen Interfaces implementieren
2. [ ] Noop-Implementierungen für alle Provider-Typen
3. [ ] Erstes konkretes Provider-Paket: `gondolia-meilisearch` (für Catalog Service Sprint 3)
4. [ ] Catalog Service mit SearchProvider-Integration
5. [ ] `gondolia-sap` als eigenes Public Repo (Specification Pattern aus `sap-integration.md`)
6. [ ] Template-Repository für Community-Provider (`gondolia-provider-template`)
7. [ ] Dokumentation: "Writing a Gondolia Provider" Guide

---

## Referenzen

- [Medusa.js Provider System](https://docs.medusajs.com/learn/fundamentals/modules/providers)
- [Vendure Plugin System](https://docs.vendure.io/guides/developer-guide/plugins/)
- [Saleor App System](https://docs.saleor.io/docs/3.x/developer/extending/apps/overview)
- [Go Interface Best Practices](https://go.dev/wiki/CodeReviewComments#interfaces)
- [SAP Integration Architecture](./sap-integration.md)
- [External Integrations Analysis](./external-integrations-analysis.md)
- [Akeneo Import Analysis](./akeneo-import-analysis.md)
- [Search Engine Evaluation](./search-engine-evaluation.md)
