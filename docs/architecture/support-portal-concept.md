# Support Portal Konzept

## Executive Summary

Dieses Dokument beschreibt die Trennung von **Admin-Backend** und **Support-Portal** für Webshop V3, mit Fokus auf Customer Journey Tracking und optimierte Support-Workflows.

---

## 1. Problemanalyse (V2 Nova)

### 1.1 Aktuelle Situation

```
┌─────────────────────────────────────────────────────────────────┐
│                    NOVA BACKEND (V2)                             │
│                                                                  │
│  43 Resources │ 13 Actions │ 6 Filters │ 10 Metrics             │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                    ALLE NUTZER                            │   │
│  │                                                           │   │
│  │  • Admins (Konfiguration, Settings)                      │   │
│  │  • Support (Kundenbetreuung)                             │   │
│  │  • Produktmanager (Katalog)                              │   │
│  │  • Buchhaltung (Bestellungen)                            │   │
│  │                                                           │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

### 1.2 Identifizierte Probleme

| Problem | Auswirkung | Betroffene |
|---------|------------|------------|
| **Zu viele Resources** | Navigation überladen, 43 Menüpunkte | Alle |
| **Keine Kundenübersicht** | Kein Gesamtbild der Customer Journey | Support |
| **SAP-Fehler als XML** | Nicht verständlich für Support | Support |
| **Manuelle Prozesse** | Viele Klicks für einfache Aufgaben | Support |
| **Keine Timeline** | Aktivitäten nicht chronologisch | Support |
| **Versteckte Quotes** | Nur über Kunden-Detail erreichbar | Support |
| **Keine Alerts** | Keine Benachrichtigung bei Problemen | Support |

### 1.3 Support-Anforderungen (Neu)

```
Was Support BRAUCHT:
├── Schnelle Kundensuche (Name, Email, SAP-Nr, Bestell-Nr)
├── Customer Journey auf einen Blick
├── Aktuelle Probleme/Fehler sofort sehen
├── Aktionen mit einem Klick (Login als Kunde, Status ändern)
├── Verständliche Fehlermeldungen (kein XML)
└── Proaktive Alerts bei Problemen

Was Support NICHT braucht:
├── Produktkatalog-Verwaltung
├── System-Konfiguration
├── Shipping/Payment Setup
├── Tenant-Verwaltung
└── CMS/Content Management
```

---

## 2. Architektur-Empfehlung

### 2.1 Getrennte Anwendungen

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              V3 BACKEND ARCHITEKTUR                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────┐          ┌─────────────────────────┐          │
│  │     ADMIN PORTAL        │          │    SUPPORT PORTAL       │          │
│  │                         │          │                         │          │
│  │  • Systemkonfiguration  │          │  • Kundenübersicht      │          │
│  │  • Produktverwaltung    │          │  • Customer Journey     │          │
│  │  • Katalog & Preise     │          │  • Bestellungen         │          │
│  │  • Shipping/Payment     │          │  • Quick Actions        │          │
│  │  • Tenant Settings      │          │  • Problem Dashboard    │          │
│  │  • User Management      │          │  • Live Search          │          │
│  │                         │          │                         │          │
│  │  Nutzer: Admins, PMs    │          │  Nutzer: Support Team   │          │
│  └───────────┬─────────────┘          └───────────┬─────────────┘          │
│              │                                    │                         │
│              └──────────────┬─────────────────────┘                         │
│                             │                                               │
│                             ▼                                               │
│              ┌─────────────────────────────┐                               │
│              │       API GATEWAY           │                               │
│              │                             │                               │
│              │  • Authentication (JWT)     │                               │
│              │  • Role-Based Access        │                               │
│              │  • Rate Limiting            │                               │
│              └─────────────┬───────────────┘                               │
│                            │                                               │
│         ┌──────────────────┼──────────────────┐                           │
│         ▼                  ▼                  ▼                           │
│  ┌─────────────┐   ┌─────────────┐   ┌─────────────┐                     │
│  │  Identity   │   │   Catalog   │   │    Order    │   ...               │
│  │  Service    │   │   Service   │   │   Service   │                     │
│  └─────────────┘   └─────────────┘   └─────────────┘                     │
│                                                                           │
└───────────────────────────────────────────────────────────────────────────┘
```

### 2.2 Vorteile der Trennung

| Aspekt | Admin Portal | Support Portal |
|--------|--------------|----------------|
| **Fokus** | Konfiguration & Daten | Kundenbetreuung |
| **Komplexität** | Hoch (Power User) | Niedrig (Effizienz) |
| **Navigation** | Hierarchisch | Aufgabenbasiert |
| **Daten** | Vollzugriff | Nur kundenrelevant |
| **Updates** | Selten (stabil) | Häufig (UX) |
| **Tech Stack** | Next.js + Refine | Next.js (custom) |

---

## 3. Support Portal Design

### 3.1 Hauptbereiche

```
┌─────────────────────────────────────────────────────────────────┐
│  🔍 Suche: Kunde, Bestellung, Firma...            [Agent: Max] │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐  │
│  │Dashboard│ │ Kunden  │ │Bestellg.│ │Probleme │ │ Firmen  │  │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘  │
│                                                                 │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│                     [ HAUPTINHALT ]                             │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 3.2 Dashboard (Startseite)

```
┌─────────────────────────────────────────────────────────────────┐
│                        SUPPORT DASHBOARD                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────────┐  ┌──────────────────┐  ┌────────────────┐ │
│  │  🔴 5 Probleme   │  │  📦 23 Bestell.  │  │  👥 142 Aktiv  │ │
│  │  Sofort handeln  │  │  Heute           │  │  Online jetzt  │ │
│  └──────────────────┘  └──────────────────┘  └────────────────┘ │
│                                                                  │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │ AKTUELLE PROBLEME                               [Alle →]    ││
│  ├─────────────────────────────────────────────────────────────┤│
│  │ 🔴 SAP-Fehler    │ Bestellung #12345 │ Müller AG │ vor 5m   ││
│  │ 🟠 Zahlung       │ Bestellung #12340 │ Meier GmbH│ vor 15m  ││
│  │ 🟡 Lager         │ Produkt ABC-123   │ 3 Kunden  │ vor 1h   ││
│  └─────────────────────────────────────────────────────────────┘│
│                                                                  │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │ LETZTE AKTIVITÄTEN                                          ││
│  ├─────────────────────────────────────────────────────────────┤│
│  │ 10:45 │ Hans Müller │ Bestellung aufgegeben │ CHF 1'234.50  ││
│  │ 10:42 │ Anna Meier  │ Warenkorb erstellt    │ 5 Artikel     ││
│  │ 10:38 │ Peter Huber │ Eingeloggt            │ Firma XY AG   ││
│  └─────────────────────────────────────────────────────────────┘│
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### 3.3 Customer Journey View

**Das Herzstück des Support Portals:**

```
┌─────────────────────────────────────────────────────────────────┐
│  KUNDE: Hans Müller                      [Als Kunde einloggen] │
│  hans.mueller@example.com │ SAP: 123456 │ Firma: Müller AG      │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐   │
│  │ Journey │ │Bestellg.│ │Warenkörbe│ │ Merkliste│ │ Tickets │   │
│  │    ●    │ │   12    │ │    2    │ │   15    │ │    1    │   │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘   │
│                                                                  │
│  ═══════════════════════════════════════════════════════════════│
│                                                                  │
│  CUSTOMER JOURNEY TIMELINE                                       │
│  ─────────────────────────────────────────────────────────────  │
│                                                                  │
│  HEUTE                                                           │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │ 10:45  📦 BESTELLUNG #12345 aufgegeben                      ││
│  │        └─ 5 Artikel │ CHF 1'234.50 │ Lieferung: Express     ││
│  │        └─ Status: ⏳ Warte auf SAP-Bestätigung              ││
│  │        └─ [Details] [SAP neu senden] [Status ändern]        ││
│  └─────────────────────────────────────────────────────────────┘│
│  ┌─────────────────────────────────────────────────────────────┐│
│  │ 10:30  🛒 WARENKORB aktualisiert                            ││
│  │        └─ +2 Artikel hinzugefügt │ Total: CHF 1'234.50      ││
│  └─────────────────────────────────────────────────────────────┘│
│  ┌─────────────────────────────────────────────────────────────┐│
│  │ 10:15  🔐 EINGELOGGT                                        ││
│  │        └─ IP: 192.168.1.100 │ Browser: Chrome │ Desktop     ││
│  └─────────────────────────────────────────────────────────────┘│
│                                                                  │
│  GESTERN                                                         │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │ 16:20  ❤️ MERKLISTE: Produkt "Eiche Laminat" hinzugefügt   ││
│  └─────────────────────────────────────────────────────────────┘│
│  ┌─────────────────────────────────────────────────────────────┐│
│  │ 14:45  🔍 SUCHE: "laminat eiche 8mm"                        ││
│  │        └─ 23 Ergebnisse │ 3 Produkte angesehen              ││
│  └─────────────────────────────────────────────────────────────┘│
│                                                                  │
│  LETZTE WOCHE                                                    │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │ Mo 15.01  📦 BESTELLUNG #12300 geliefert                    ││
│  │           └─ 3 Artikel │ CHF 890.00 │ ✅ Bezahlt            ││
│  └─────────────────────────────────────────────────────────────┘│
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### 3.4 Event-Typen für Journey

| Event | Icon | Beschreibung | Datenquelle |
|-------|------|--------------|-------------|
| `login` | 🔐 | Kunde eingeloggt | Auth Service |
| `logout` | 🚪 | Kunde ausgeloggt | Auth Service |
| `search` | 🔍 | Suchanfrage | Search Service |
| `product_view` | 👁️ | Produkt angesehen | Analytics |
| `cart_add` | 🛒 | Warenkorb hinzugefügt | Cart Service |
| `cart_remove` | ➖ | Warenkorb entfernt | Cart Service |
| `cart_update` | 🔄 | Menge geändert | Cart Service |
| `wishlist_add` | ❤️ | Merkliste hinzugefügt | Catalog Service |
| `quote_created` | 📋 | Angebot erstellt | Order Service |
| `quote_simulated` | 💰 | Preis berechnet | SAP Service |
| `order_placed` | 📦 | Bestellung aufgegeben | Order Service |
| `order_paid` | 💳 | Bezahlt | Payment Service |
| `order_shipped` | 🚚 | Versendet | SAP Event |
| `order_delivered` | ✅ | Geliefert | SAP Event |
| `order_error` | 🔴 | Fehler | SAP Service |
| `sap_error` | ⚠️ | SAP-Fehlermeldung (gemappt) | SAP Service |
| `support_contact` | 📞 | Support kontaktiert | Ticket System |

> **Hinweis zu SAP-Fehlern:** Die originalen SAP-Meldungen (z.B. "Kunde ist nicht kreditwürdig") werden dem Kunden NICHT angezeigt. Stattdessen sieht der Kunde eine kundenfreundliche, lokalisierte Meldung. Support-Mitarbeiter können jedoch die Original-SAP-Meldung in der Customer Journey einsehen.
> Siehe: [SAP Error Message Mapping](./sap-integration.md#sap-error-message-mapping)

---

## 4. Datenmodell

### 4.1 Customer Journey Events

```sql
CREATE TABLE customer_journey_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    customer_id UUID NOT NULL,
    company_id UUID,

    -- Event Identifikation
    event_type VARCHAR(50) NOT NULL,      -- 'order_placed', 'cart_add', etc.
    event_category VARCHAR(30) NOT NULL,  -- 'order', 'cart', 'auth', 'search'

    -- Event Details
    title VARCHAR(255) NOT NULL,          -- Kurzbeschreibung
    description TEXT,                     -- Detailtext
    metadata JSONB,                       -- Flexible Zusatzdaten

    -- Referenzen
    reference_type VARCHAR(50),           -- 'order', 'product', 'cart'
    reference_id UUID,                    -- ID des referenzierten Objekts

    -- Context
    session_id VARCHAR(100),
    ip_address INET,
    user_agent TEXT,
    device_type VARCHAR(20),              -- 'desktop', 'mobile', 'tablet'

    -- Severity (für Probleme)
    severity VARCHAR(10),                 -- 'info', 'warning', 'error', 'critical'

    -- Timestamps
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Indexes
    CONSTRAINT fk_customer FOREIGN KEY (customer_id) REFERENCES customers(id),
    CONSTRAINT fk_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);

-- Performance Indexes
CREATE INDEX idx_journey_customer ON customer_journey_events(customer_id, occurred_at DESC);
CREATE INDEX idx_journey_tenant_date ON customer_journey_events(tenant_id, occurred_at DESC);
CREATE INDEX idx_journey_severity ON customer_journey_events(severity, occurred_at DESC)
    WHERE severity IS NOT NULL;
CREATE INDEX idx_journey_reference ON customer_journey_events(reference_type, reference_id);
```

### 4.2 Metadata Beispiele

```json
// Event: order_placed
{
    "order_number": "WS-12345",
    "item_count": 5,
    "total_amount": 1234.50,
    "currency": "CHF",
    "shipping_method": "express",
    "payment_method": "invoice"
}

// Event: search
{
    "query": "laminat eiche 8mm",
    "results_count": 23,
    "filters_applied": {
        "category": "bodenbelaege",
        "thickness": "8mm"
    },
    "clicked_results": ["p_100096", "p_100097"]
}

// Event: order_error
{
    "order_id": "uuid-123",
    "error_code": "SAP_TIMEOUT",
    "error_message": "SAP nicht erreichbar",
    "retry_count": 2,
    "sap_function": "Z_BAPI_SALESORDER_CREATE"
}
```

### 4.3 Problem/Alert Tracking

```sql
CREATE TABLE support_alerts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,

    -- Alert Typ
    alert_type VARCHAR(50) NOT NULL,      -- 'sap_error', 'payment_failed', 'stock_out'
    severity VARCHAR(10) NOT NULL,        -- 'warning', 'error', 'critical'
    status VARCHAR(20) NOT NULL DEFAULT 'open',  -- 'open', 'in_progress', 'resolved', 'ignored'

    -- Betroffene Entitäten
    customer_id UUID,
    company_id UUID,
    order_id UUID,
    product_id UUID,

    -- Alert Details
    title VARCHAR(255) NOT NULL,
    description TEXT,
    metadata JSONB,

    -- Resolution
    resolved_by UUID,                     -- Support User ID
    resolved_at TIMESTAMPTZ,
    resolution_note TEXT,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_alerts_status ON support_alerts(tenant_id, status, severity, created_at DESC);
```

---

## 5. API Design

### 5.1 Support Portal API Endpoints

```yaml
# Customer Journey
GET  /api/support/v1/customers/{id}/journey
     ?from=2024-01-01&to=2024-01-31
     &event_types=order,cart,auth
     &limit=50&offset=0

GET  /api/support/v1/customers/{id}/summary
     # Aggregierte Daten: Bestellungen, Umsatz, letzte Aktivität

# Suche
GET  /api/support/v1/search
     ?q=müller
     &type=customer,order,company
     &limit=20

# Alerts/Probleme
GET  /api/support/v1/alerts
     ?status=open
     &severity=error,critical
     &limit=50

PATCH /api/support/v1/alerts/{id}
      # Status ändern, Resolution Note

# Quick Actions
POST /api/support/v1/customers/{id}/login-as
     # Generiert temporären Login-Link

POST /api/support/v1/orders/{id}/retry-sap
     # SAP Export erneut versuchen

PATCH /api/support/v1/orders/{id}/status
      # Bestellstatus ändern
```

### 5.2 Response Beispiel: Customer Journey

```json
{
  "customer": {
    "id": "uuid-123",
    "name": "Hans Müller",
    "email": "hans.mueller@example.com",
    "sap_number": "123456",
    "company": {
      "id": "uuid-456",
      "name": "Müller AG"
    }
  },
  "summary": {
    "total_orders": 12,
    "total_revenue": 15234.50,
    "last_order_at": "2024-01-15T10:45:00Z",
    "last_login_at": "2024-01-15T10:15:00Z",
    "open_carts": 2,
    "wishlist_items": 15
  },
  "events": [
    {
      "id": "evt-001",
      "type": "order_placed",
      "category": "order",
      "title": "Bestellung #12345 aufgegeben",
      "severity": null,
      "occurred_at": "2024-01-15T10:45:00Z",
      "metadata": {
        "order_number": "WS-12345",
        "item_count": 5,
        "total_amount": 1234.50
      },
      "reference": {
        "type": "order",
        "id": "uuid-order-123"
      },
      "actions": [
        {"name": "view_order", "label": "Details anzeigen"},
        {"name": "retry_sap", "label": "SAP neu senden", "enabled": true}
      ]
    },
    {
      "id": "evt-002",
      "type": "cart_update",
      "category": "cart",
      "title": "Warenkorb aktualisiert",
      "occurred_at": "2024-01-15T10:30:00Z",
      "metadata": {
        "items_added": 2,
        "cart_total": 1234.50
      }
    }
  ],
  "pagination": {
    "total": 156,
    "limit": 50,
    "offset": 0,
    "has_more": true
  }
}
```

---

## 6. Event Collection

### 6.1 Event Sources

```
┌─────────────────────────────────────────────────────────────────┐
│                      EVENT COLLECTION                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │  Identity   │  │   Catalog   │  │    Cart     │             │
│  │  Service    │  │   Service   │  │   Service   │             │
│  │             │  │             │  │             │             │
│  │ • login     │  │ • view      │  │ • add       │             │
│  │ • logout    │  │ • search    │  │ • remove    │             │
│  │ • register  │  │ • wishlist  │  │ • update    │             │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘             │
│         │                │                │                     │
│         └────────────────┼────────────────┘                     │
│                          ▼                                      │
│              ┌─────────────────────┐                           │
│              │       KAFKA         │                           │
│              │                     │                           │
│              │ topic: customer.*   │                           │
│              └──────────┬──────────┘                           │
│                         │                                       │
│                         ▼                                       │
│              ┌─────────────────────┐                           │
│              │   Journey Service   │                           │
│              │                     │                           │
│              │ • Event Aggregation │                           │
│              │ • Timeline Building │                           │
│              │ • Alert Generation  │                           │
│              └──────────┬──────────┘                           │
│                         │                                       │
│                         ▼                                       │
│              ┌─────────────────────┐                           │
│              │    PostgreSQL       │                           │
│              │                     │                           │
│              │ • journey_events    │                           │
│              │ • support_alerts    │                           │
│              └─────────────────────┘                           │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 6.2 Event Publishing (Services)

```go
// Jeder Service publisht Events
type JourneyEventPublisher interface {
    Publish(ctx context.Context, event *JourneyEvent) error
}

type JourneyEvent struct {
    TenantID     string                 `json:"tenant_id"`
    CustomerID   string                 `json:"customer_id"`
    CompanyID    *string                `json:"company_id,omitempty"`
    EventType    string                 `json:"event_type"`
    Category     string                 `json:"category"`
    Title        string                 `json:"title"`
    Description  *string                `json:"description,omitempty"`
    Metadata     map[string]interface{} `json:"metadata,omitempty"`
    ReferenceType *string               `json:"reference_type,omitempty"`
    ReferenceID  *string                `json:"reference_id,omitempty"`
    Severity     *string                `json:"severity,omitempty"`
    SessionID    *string                `json:"session_id,omitempty"`
    OccurredAt   time.Time              `json:"occurred_at"`
}

// Beispiel: Order Service
func (s *OrderService) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*Order, error) {
    order, err := s.repo.Create(ctx, req)
    if err != nil {
        return nil, err
    }

    // Journey Event publishen
    s.journeyPublisher.Publish(ctx, &JourneyEvent{
        TenantID:      order.TenantID,
        CustomerID:    order.CustomerID,
        CompanyID:     &order.CompanyID,
        EventType:     "order_placed",
        Category:      "order",
        Title:         fmt.Sprintf("Bestellung #%s aufgegeben", order.OrderNumber),
        Metadata: map[string]interface{}{
            "order_number":    order.OrderNumber,
            "item_count":      len(order.Items),
            "total_amount":    order.Total,
            "shipping_method": order.ShippingMethod,
        },
        ReferenceType: stringPtr("order"),
        ReferenceID:   &order.ID,
        OccurredAt:    time.Now(),
    })

    return order, nil
}
```

### 6.3 Alert Generation

```go
// Journey Service: Automatische Alert-Generierung
func (s *JourneyService) ProcessEvent(ctx context.Context, event *JourneyEvent) error {
    // 1. Event speichern
    if err := s.repo.SaveEvent(ctx, event); err != nil {
        return err
    }

    // 2. Alert-Regeln prüfen
    if alert := s.checkAlertRules(event); alert != nil {
        if err := s.repo.CreateAlert(ctx, alert); err != nil {
            return err
        }

        // WebSocket Notification an Support Portal
        s.notifier.NotifySupport(ctx, alert)
    }

    return nil
}

// Alert Rules
var alertRules = []AlertRule{
    {
        EventType: "order_error",
        Severity:  "error",
        Title:     "Bestellung fehlgeschlagen",
    },
    {
        EventType: "payment_failed",
        Severity:  "warning",
        Title:     "Zahlung fehlgeschlagen",
    },
    {
        EventType: "sap_timeout",
        Severity:  "critical",
        Title:     "SAP nicht erreichbar",
    },
}
```

---

## 7. UI Components

### 7.1 Global Search

```tsx
// components/GlobalSearch.tsx
interface SearchResult {
  type: 'customer' | 'order' | 'company' | 'product';
  id: string;
  title: string;
  subtitle: string;
  badges?: Badge[];
}

function GlobalSearch() {
  const [query, setQuery] = useState('');
  const { data: results } = useSearch(query);

  return (
    <Command>
      <CommandInput
        placeholder="Kunde, Bestellung, Firma suchen..."
        value={query}
        onChange={setQuery}
      />
      <CommandList>
        <CommandGroup heading="Kunden">
          {results?.customers.map(c => (
            <CommandItem key={c.id}>
              <UserIcon />
              <span>{c.name}</span>
              <span className="text-muted">{c.email}</span>
              {c.hasOpenIssues && <Badge variant="destructive">Problem</Badge>}
            </CommandItem>
          ))}
        </CommandGroup>
        <CommandGroup heading="Bestellungen">
          {results?.orders.map(o => (
            <CommandItem key={o.id}>
              <PackageIcon />
              <span>#{o.orderNumber}</span>
              <span className="text-muted">{o.customerName}</span>
              <Badge>{o.status}</Badge>
            </CommandItem>
          ))}
        </CommandGroup>
      </CommandList>
    </Command>
  );
}
```

### 7.2 Journey Timeline

```tsx
// components/JourneyTimeline.tsx
interface TimelineEvent {
  id: string;
  type: string;
  category: string;
  title: string;
  occurredAt: Date;
  metadata: Record<string, any>;
  severity?: 'info' | 'warning' | 'error' | 'critical';
  actions?: Action[];
}

function JourneyTimeline({ customerId }: { customerId: string }) {
  const { data, fetchNextPage, hasNextPage } = useInfiniteJourney(customerId);

  return (
    <div className="space-y-4">
      {data?.pages.map(page =>
        page.events.map(event => (
          <TimelineItem key={event.id} event={event} />
        ))
      )}

      {hasNextPage && (
        <Button onClick={fetchNextPage}>Mehr laden</Button>
      )}
    </div>
  );
}

function TimelineItem({ event }: { event: TimelineEvent }) {
  const Icon = eventIcons[event.type];

  return (
    <div className={cn(
      "flex gap-4 p-4 rounded-lg border",
      event.severity === 'error' && "border-red-500 bg-red-50",
      event.severity === 'warning' && "border-yellow-500 bg-yellow-50"
    )}>
      <div className="flex-shrink-0">
        <Icon className="w-5 h-5" />
      </div>

      <div className="flex-1 space-y-1">
        <div className="flex items-center justify-between">
          <span className="font-medium">{event.title}</span>
          <time className="text-sm text-muted-foreground">
            {formatRelativeTime(event.occurredAt)}
          </time>
        </div>

        {event.metadata && (
          <EventMetadata type={event.type} data={event.metadata} />
        )}

        {event.actions && event.actions.length > 0 && (
          <div className="flex gap-2 pt-2">
            {event.actions.map(action => (
              <Button
                key={action.name}
                variant="outline"
                size="sm"
                onClick={() => executeAction(action)}
                disabled={!action.enabled}
              >
                {action.label}
              </Button>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
```

### 7.3 Alert Dashboard

```tsx
// components/AlertDashboard.tsx
function AlertDashboard() {
  const { data: alerts } = useAlerts({ status: 'open' });

  const critical = alerts?.filter(a => a.severity === 'critical') || [];
  const errors = alerts?.filter(a => a.severity === 'error') || [];
  const warnings = alerts?.filter(a => a.severity === 'warning') || [];

  return (
    <div className="space-y-6">
      {/* Stats */}
      <div className="grid grid-cols-3 gap-4">
        <StatCard
          title="Kritisch"
          value={critical.length}
          icon={<AlertCircle className="text-red-500" />}
          variant="destructive"
        />
        <StatCard
          title="Fehler"
          value={errors.length}
          icon={<XCircle className="text-orange-500" />}
          variant="warning"
        />
        <StatCard
          title="Warnungen"
          value={warnings.length}
          icon={<AlertTriangle className="text-yellow-500" />}
        />
      </div>

      {/* Alert List */}
      <Card>
        <CardHeader>
          <CardTitle>Aktuelle Probleme</CardTitle>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Status</TableHead>
                <TableHead>Problem</TableHead>
                <TableHead>Kunde/Firma</TableHead>
                <TableHead>Zeit</TableHead>
                <TableHead>Aktionen</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {alerts?.map(alert => (
                <AlertRow key={alert.id} alert={alert} />
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
}
```

---

## 8. Berechtigungen

### 8.1 Rollen

| Rolle | Admin Portal | Support Portal | Beschreibung |
|-------|--------------|----------------|--------------|
| `super_admin` | ✅ Vollzugriff | ✅ Vollzugriff | System-Administrator |
| `admin` | ✅ Vollzugriff | ❌ Kein Zugriff | Tenant-Administrator |
| `product_manager` | ✅ Katalog only | ❌ Kein Zugriff | Produktverwaltung |
| `support_lead` | ❌ Kein Zugriff | ✅ Vollzugriff | Support-Teamleiter |
| `support_agent` | ❌ Kein Zugriff | ✅ Eingeschränkt | Support-Mitarbeiter |

### 8.2 Support-Berechtigungen

```yaml
support_agent:
  can_view:
    - customers
    - orders
    - journey_events
    - alerts

  can_execute:
    - login_as_customer      # Mit Logging
    - change_order_status    # Nur bestimmte Status
    - retry_sap_export
    - resolve_alert

  cannot:
    - delete_customer
    - delete_order
    - change_prices
    - access_admin_settings

support_lead:
  extends: support_agent
  can_execute:
    - delete_customer        # Mit Bestätigung
    - bulk_operations
    - export_data
    - manage_support_agents
```

### 8.3 Audit Trail

```go
// Alle Support-Aktionen werden geloggt
type SupportAuditLog struct {
    ID           string
    AgentID      string    // Support-Mitarbeiter
    Action       string    // "login_as_customer", "change_order_status"
    TargetType   string    // "customer", "order"
    TargetID     string
    OldValue     *string   // Vorheriger Wert (JSON)
    NewValue     *string   // Neuer Wert (JSON)
    Reason       *string   // Optionale Begründung
    IPAddress    string
    CreatedAt    time.Time
}
```

---

## 9. Tech Stack Empfehlung

### 9.1 Support Portal Frontend

```yaml
Framework: Next.js 14 (App Router)
UI Library: shadcn/ui + Tailwind CSS
State: TanStack Query (React Query)
Forms: React Hook Form + Zod
Tables: TanStack Table
Charts: Recharts
Real-time: WebSocket (für Alerts)
Search: Cmdk (Command palette)
```

### 9.2 Backend Services

```yaml
# Neuer Service für Support
services/
  support/
    internal/
      journey/       # Journey Event Handling
      alert/         # Alert Management
      search/        # Global Search
      action/        # Quick Actions
```

---

## 10. Migration Roadmap

### Phase 1: Foundation (Woche 1-2)

```
□ Journey Service implementieren
□ Event-Schema definieren
□ Kafka Topics einrichten
□ PostgreSQL Tabellen erstellen
□ Basis-API Endpoints
```

### Phase 2: Event Collection (Woche 3-4)

```
□ Identity Service: Auth Events
□ Cart Service: Cart Events
□ Order Service: Order Events
□ SAP Service: Error Events
□ Event Aggregation
```

### Phase 3: Support Portal UI (Woche 5-7)

```
□ Next.js Projekt Setup
□ Authentication/Authorization
□ Global Search
□ Customer Journey View
□ Alert Dashboard
□ Quick Actions
```

### Phase 4: Advanced Features (Woche 8+)

```
□ Real-time WebSocket Alerts
□ Bulk Operations
□ Export/Reports
□ Keyboard Shortcuts
□ Mobile Responsive
```

---

## 11. Zusammenfassung

### Vorteile der Trennung

| Aspekt | Vorher (Nova) | Nachher (Getrennt) |
|--------|---------------|-------------------|
| **Fokus** | Alles für Alle | Spezialisiert |
| **UX** | CRUD-basiert | Workflow-basiert |
| **Performance** | 43 Resources laden | Nur relevante Daten |
| **Onboarding** | Komplex | Einfach |
| **Wartung** | Monolithisch | Unabhängig |

### Key Features Support Portal

1. **Customer Journey Timeline** - Alle Aktivitäten chronologisch
2. **Global Search** - Kunde, Bestellung, Firma mit einem Tastendruck
3. **Alert Dashboard** - Probleme sofort sehen
4. **Quick Actions** - Login als Kunde, Status ändern mit einem Klick
5. **Verständliche Fehler** - Kein XML, klare Meldungen
6. **Real-time Updates** - WebSocket für neue Alerts

### Empfehlung

**Ja, die Trennung macht Sinn.** Support und Admin haben fundamental unterschiedliche Bedürfnisse:

- **Admin** = Konfiguration, Daten pflegen, Power User
- **Support** = Kunden helfen, schnell reagieren, Effizienz

---

## 12. V3 Frontend Architektur (Einheitlicher Stack)

### Kein PHP in V3

V3 verwendet **ausschließlich Next.js** für alle Frontends. Nova/Laravel wird nicht übernommen.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         V3 FRONTEND MONOREPO                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  packages/                                                                   │
│  └── ui/                         # Shared Component Library                 │
│      ├── components/             # Button, Table, Form, Card, etc.          │
│      ├── hooks/                  # useAuth, useApi, useTenant               │
│      └── styles/                 # Tailwind Config, Theme                   │
│                                                                              │
│  apps/                                                                       │
│  ├── shop/                       # B2B Webshop (Kunden)                     │
│  │   └── Next.js 14                                                         │
│  │                                                                           │
│  ├── admin/                      # Admin Portal (Produktmanager, Admins)   │
│  │   └── Next.js 14 + Refine.dev                                           │
│  │                                                                           │
│  └── support/                    # Support Portal (Customer Service)        │
│      └── Next.js 14 + Custom UI                                             │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              GO BACKEND                                      │
│                                                                              │
│  API Gateway → Identity │ Catalog │ Cart │ Order │ Support │ ...           │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Tech Stack Übersicht

| Komponente | V2 (Legacy) | V3 (Neu) |
|------------|-------------|----------|
| **Shop Frontend** | Vue.js | Next.js 14 |
| **Admin Portal** | Laravel Nova (PHP) | Next.js + Refine.dev |
| **Support Portal** | (Teil von Nova) | Next.js (Custom) |
| **Backend** | Laravel (PHP) | Go Microservices |
| **UI Components** | Tailwind + Custom | shadcn/ui (Shared) |

### Vorteile einheitlicher Stack

1. **Ein Frontend-Team** - Nur TypeScript/React Skills nötig
2. **Shared Components** - UI Library für alle Apps
3. **Shared Types** - API Types einmal definiert, überall genutzt
4. **Einheitliche Tooling** - ESLint, Prettier, Testing
5. **Kein PHP** - Keine PHP-Infrastruktur mehr nötig

### Admin Portal mit Refine.dev

Für CRUD-intensive Admin-Funktionen empfehlen wir **Refine.dev**:

```tsx
// apps/admin/src/resources/products.tsx
import { List, useTable, EditButton } from "@refinedev/antd";

export const ProductList = () => {
  const { tableProps } = useTable({
    resource: "products",
    syncWithLocation: true,
  });

  return (
    <List>
      <Table {...tableProps} rowKey="id">
        <Table.Column dataIndex="sku" title="SKU" sorter />
        <Table.Column dataIndex="name" title="Name" />
        <Table.Column dataIndex="price" title="Preis" />
        <Table.Column
          title="Aktionen"
          render={(_, record) => <EditButton recordItemId={record.id} />}
        />
      </Table>
    </List>
  );
};
```

**Refine Features:**
- Auto-generated CRUD
- Data Provider für REST/GraphQL
- Auth Provider Integration
- i18n Support
- Audit Logs
- Access Control
