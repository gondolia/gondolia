# Hersteller-Entität (Manufacturer)

Stand: 2026-02-14

## Übersicht

Ein **Hersteller (Manufacturer)** ist eine eigenständige Entität im Katalog-Service. Produkte haben einen Hersteller — oder keinen. Hersteller sind relevant für:

- **Katalog & Navigation:** Kunden filtern und browsen nach Hersteller
- **Suche:** Hersteller als Facette in der Produktsuche, Hersteller selbst durchsuchbar
- **Frontend:** Hersteller-Übersichtsseiten, Hersteller-Detailseiten, Hersteller-Anzeige auf Produktdetailseiten
- **ERP/SAP:** Mapping auf SAP-Lieferanten-/Hersteller-Stammdaten
- **PIM:** Import aus Akeneo/PIM-Systemen

### Branchenbeispiele

| Branche | Typische Hersteller | Herstelleranzahl pro Shop |
|---------|--------------------|-----------------------------|
| Industriebedarf | Bosch Rexroth, Festo, SKF, Parker Hannifin, Siemens | 200–500 |
| Elektro/Elektronik | Lapp Kabel, Phoenix Contact, Wago, Osram, Weidmüller | 150–400 |
| Arbeitsschutz | Uvex, 3M, Ansell, Honeywell, Dräger, Moldex | 100–300 |
| Labor/Chemie | Merck, Carl Roth, VWR, Eppendorf, Sartorius, IKA | 200–600 |
| Verpackung | Tesa, Beiersdorf, Rajapack, Kiefel, Multivac | 80–200 |
| IT/Hardware | Dell, HPE, Cisco, APC, Lenovo, Intel | 100–300 |

In einem typischen B2B-Shop hat ein Tenant zwischen **100 und 600 aktive Hersteller** mit jeweils **10 bis 5'000 Produkten**.

---

## 1. Entität: Manufacturer

### Felder

| Feld | Typ | Pflicht | Beschreibung |
|------|-----|---------|--------------|
| `id` | UUID | ✅ | Primärschlüssel |
| `tenant_id` | UUID | ✅ | Mandantenzugehörigkeit (Multi-Tenancy) |
| `code` | string | ✅ | Eindeutiger, maschinenlesbarer Code pro Tenant (z.B. `bosch-rexroth`, `3m`) |
| `name` | map[string]string | ✅ | Mehrsprachiger Name (`{"de": "Bosch Rexroth", "fr": "Bosch Rexroth"}`) |
| `slug` | map[string]string | ✅ | URL-Slug pro Sprache (`{"de": "bosch-rexroth", "fr": "bosch-rexroth"}`) |
| `description` | map[string]string | ❌ | Mehrsprachige Beschreibung (Freitext, Markdown) |
| `logo_url` | string | ❌ | URL zum Hersteller-Logo (aus Asset Storage) |
| `website` | string | ❌ | Offizielle Website des Herstellers |
| `country` | string | ❌ | Herkunftsland (ISO 3166-1 Alpha-2, z.B. `DE`, `US`, `CH`) |
| `status` | enum | ✅ | `active` oder `inactive` |
| `meta_title` | map[string]string | ❌ | SEO: Meta-Title pro Sprache |
| `meta_description` | map[string]string | ❌ | SEO: Meta-Description pro Sprache |
| `contact_email` | string | ❌ | Interne Kontakt-E-Mail (nicht öffentlich) |
| `contact_phone` | string | ❌ | Internes Kontakttelefon (nicht öffentlich) |
| `contact_person` | string | ❌ | Interner Ansprechpartner (nicht öffentlich) |
| `sort_order` | int | ❌ | Manuelle Sortierung (0 = automatisch alphabetisch) |
| `pim_identifier` | string | ❌ | Externe ID aus Akeneo/PIM |
| `erp_identifier` | string | ❌ | Externe ID aus SAP/ERP |
| `last_synced_at` | timestamp | ❌ | Letzter PIM/ERP-Sync |
| `created_at` | timestamp | ✅ | Erstellt am |
| `updated_at` | timestamp | ✅ | Zuletzt geändert |
| `deleted_at` | timestamp | ❌ | Soft-Delete Zeitstempel |

**Hinweise:**
- `code` ist pro Tenant eindeutig und dient als stabiler Identifier für Imports und API-Zugriffe
- `slug` wird automatisch aus `name` generiert, kann aber manuell überschrieben werden
- Kontaktdaten sind **rein intern** (Admin-Bereich) und werden nie an das Frontend ausgeliefert
- `pim_identifier` und `erp_identifier` ermöglichen die Zuordnung zu externen Systemen

---

## 2. Beziehung zu Produkten

### Kardinalität

```
Manufacturer 1 ──── 0..* Product
```

- Ein **Produkt** hat **genau einen Hersteller** — oder **keinen** (`manufacturer_id` ist nullable)
- Ein **Hersteller** hat **beliebig viele Produkte**
- `manufacturer_id` ist ein Foreign Key in der `products`-Tabelle

### Löschverhalten

**Soft-Delete mit Constraint:**

Ein Hersteller kann nur soft-deleted werden (Feld `deleted_at` gesetzt). Das Verhalten bei Löschung:

| Szenario | Verhalten |
|----------|-----------|
| Hersteller hat 0 aktive Produkte | Soft-Delete sofort möglich |
| Hersteller hat aktive Produkte | Warnung im Admin-UI: «Hersteller hat X aktive Produkte. Wirklich deaktivieren?» |
| Soft-Delete durchgeführt | `deleted_at` wird gesetzt. Produkte behalten ihre `manufacturer_id`. |
| Produkt-Anzeige nach Soft-Delete | Hersteller-Name wird weiterhin angezeigt. Kein Link zur Hersteller-Seite. |
| Hersteller-Seite nach Soft-Delete | 404 — Seite nicht mehr erreichbar |
| Wiederherstellung | `deleted_at` auf NULL setzen → Hersteller wieder aktiv |

**Kein Hard-Delete.** Kein Reassign. Produkte bleiben dem gelöschten Hersteller zugeordnet. Im Frontend wird der Hersteller-Name angezeigt, aber die Hersteller-Seite ist nicht mehr erreichbar.

**Begründung:** In einem B2B-Kontext verschwinden Hersteller nicht einfach — sie werden aufgekauft, umbennant oder stellen die Produktion ein. Die historische Zuordnung bleibt wichtig für Bestellhistorie, SAP-Mapping und Reporting.

---

## 3. Datenmodell

### SQL-Schema

```sql
-- Migration: 000003_create_manufacturers.up.sql

CREATE TABLE manufacturers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    code VARCHAR(100) NOT NULL,
    name JSONB NOT NULL,                        -- {"de": "Bosch Rexroth", "fr": "Bosch Rexroth"}
    slug JSONB NOT NULL,                        -- {"de": "bosch-rexroth", "fr": "bosch-rexroth"}
    description JSONB,                          -- {"de": "...", "fr": "..."}
    logo_url VARCHAR(500),
    website VARCHAR(500),
    country VARCHAR(2),                         -- ISO 3166-1 Alpha-2
    status VARCHAR(20) NOT NULL DEFAULT 'active',

    -- SEO
    meta_title JSONB,
    meta_description JSONB,

    -- Interne Kontaktdaten
    contact_email VARCHAR(255),
    contact_phone VARCHAR(50),
    contact_person VARCHAR(255),

    -- Sortierung
    sort_order INTEGER NOT NULL DEFAULT 0,

    -- Externe Systeme
    pim_identifier VARCHAR(255),
    erp_identifier VARCHAR(255),
    last_synced_at TIMESTAMPTZ,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    -- Constraints
    UNIQUE(tenant_id, code)
);

-- Indizes
CREATE INDEX idx_manufacturers_tenant_status ON manufacturers(tenant_id, status) WHERE deleted_at IS NULL;
CREATE INDEX idx_manufacturers_tenant_slug ON manufacturers USING GIN (slug) WHERE deleted_at IS NULL;
CREATE INDEX idx_manufacturers_pim ON manufacturers(tenant_id, pim_identifier) WHERE pim_identifier IS NOT NULL;
CREATE INDEX idx_manufacturers_erp ON manufacturers(tenant_id, erp_identifier) WHERE erp_identifier IS NOT NULL;

-- FK auf products
ALTER TABLE products ADD COLUMN manufacturer_id UUID REFERENCES manufacturers(id);
CREATE INDEX idx_products_manufacturer ON products(manufacturer_id) WHERE manufacturer_id IS NOT NULL;
```

```sql
-- Migration: 000003_create_manufacturers.down.sql

ALTER TABLE products DROP COLUMN IF EXISTS manufacturer_id;
DROP TABLE IF EXISTS manufacturers;
```

### Go Domain Struct

```go
// internal/manufacturer/models.go
package manufacturer

import (
    "time"
    "github.com/google/uuid"
)

type ManufacturerStatus string

const (
    ManufacturerStatusActive   ManufacturerStatus = "active"
    ManufacturerStatusInactive ManufacturerStatus = "inactive"
)

type Manufacturer struct {
    ID              uuid.UUID          `json:"id" db:"id"`
    TenantID        uuid.UUID          `json:"tenant_id" db:"tenant_id"`
    Code            string             `json:"code" db:"code"`
    Name            map[string]string  `json:"name" db:"name"`
    Slug            map[string]string  `json:"slug" db:"slug"`
    Description     map[string]string  `json:"description,omitempty" db:"description"`
    LogoURL         string             `json:"logo_url,omitempty" db:"logo_url"`
    Website         string             `json:"website,omitempty" db:"website"`
    Country         string             `json:"country,omitempty" db:"country"`
    Status          ManufacturerStatus `json:"status" db:"status"`

    // SEO
    MetaTitle       map[string]string  `json:"meta_title,omitempty" db:"meta_title"`
    MetaDescription map[string]string  `json:"meta_description,omitempty" db:"meta_description"`

    // Interne Kontaktdaten (nie im Public API Response)
    ContactEmail    string             `json:"-" db:"contact_email"`
    ContactPhone    string             `json:"-" db:"contact_phone"`
    ContactPerson   string             `json:"-" db:"contact_person"`

    // Sortierung
    SortOrder       int                `json:"sort_order" db:"sort_order"`

    // Externe Systeme
    PIMIdentifier   *string            `json:"pim_identifier,omitempty" db:"pim_identifier"`
    ERPIdentifier   *string            `json:"erp_identifier,omitempty" db:"erp_identifier"`
    LastSyncedAt    *time.Time         `json:"last_synced_at,omitempty" db:"last_synced_at"`

    // Timestamps
    CreatedAt       time.Time          `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time          `json:"updated_at" db:"updated_at"`
    DeletedAt       *time.Time         `json:"-" db:"deleted_at"`

    // Berechnete Felder (nicht in DB)
    ProductCount    int                `json:"product_count,omitempty" db:"-"`
}

// --- Request DTOs ---

type CreateManufacturerRequest struct {
    Code            string            `json:"code" binding:"required,min=1,max=100"`
    Name            map[string]string `json:"name" binding:"required"`
    Description     map[string]string `json:"description,omitempty"`
    LogoURL         string            `json:"logo_url,omitempty"`
    Website         string            `json:"website,omitempty"`
    Country         string            `json:"country,omitempty" binding:"omitempty,len=2"`
    Status          ManufacturerStatus `json:"status,omitempty"`
    MetaTitle       map[string]string `json:"meta_title,omitempty"`
    MetaDescription map[string]string `json:"meta_description,omitempty"`
    ContactEmail    string            `json:"contact_email,omitempty"`
    ContactPhone    string            `json:"contact_phone,omitempty"`
    ContactPerson   string            `json:"contact_person,omitempty"`
    SortOrder       int               `json:"sort_order,omitempty"`
    PIMIdentifier   *string           `json:"pim_identifier,omitempty"`
    ERPIdentifier   *string           `json:"erp_identifier,omitempty"`
}

type UpdateManufacturerRequest struct {
    Name            map[string]string  `json:"name,omitempty"`
    Description     map[string]string  `json:"description,omitempty"`
    LogoURL         *string            `json:"logo_url,omitempty"`
    Website         *string            `json:"website,omitempty"`
    Country         *string            `json:"country,omitempty"`
    Status          *ManufacturerStatus `json:"status,omitempty"`
    MetaTitle       map[string]string  `json:"meta_title,omitempty"`
    MetaDescription map[string]string  `json:"meta_description,omitempty"`
    ContactEmail    *string            `json:"contact_email,omitempty"`
    ContactPhone    *string            `json:"contact_phone,omitempty"`
    ContactPerson   *string            `json:"contact_person,omitempty"`
    SortOrder       *int               `json:"sort_order,omitempty"`
    PIMIdentifier   *string            `json:"pim_identifier,omitempty"`
    ERPIdentifier   *string            `json:"erp_identifier,omitempty"`
}

type ManufacturerFilter struct {
    TenantID uuid.UUID
    Status   *ManufacturerStatus
    Country  *string
    Search   *string // Suche in Name und Code
    Limit    int
    Offset   int
}

// --- Response DTOs ---

type ManufacturerListResponse struct {
    ID           uuid.UUID         `json:"id"`
    Code         string            `json:"code"`
    Name         map[string]string `json:"name"`
    Slug         map[string]string `json:"slug"`
    LogoURL      string            `json:"logo_url,omitempty"`
    Country      string            `json:"country,omitempty"`
    Status       ManufacturerStatus `json:"status"`
    ProductCount int               `json:"product_count"`
}

type ManufacturerDetailResponse struct {
    ID              uuid.UUID          `json:"id"`
    Code            string             `json:"code"`
    Name            map[string]string  `json:"name"`
    Slug            map[string]string  `json:"slug"`
    Description     map[string]string  `json:"description,omitempty"`
    LogoURL         string             `json:"logo_url,omitempty"`
    Website         string             `json:"website,omitempty"`
    Country         string             `json:"country,omitempty"`
    Status          ManufacturerStatus `json:"status"`
    MetaTitle       map[string]string  `json:"meta_title,omitempty"`
    MetaDescription map[string]string  `json:"meta_description,omitempty"`
    ProductCount    int                `json:"product_count"`
    CreatedAt       time.Time          `json:"created_at"`
    UpdatedAt       time.Time          `json:"updated_at"`
}

// Admin-Response enthält zusätzlich Kontaktdaten und externe IDs
type ManufacturerAdminResponse struct {
    ManufacturerDetailResponse
    ContactEmail  string     `json:"contact_email,omitempty"`
    ContactPhone  string     `json:"contact_phone,omitempty"`
    ContactPerson string     `json:"contact_person,omitempty"`
    PIMIdentifier *string    `json:"pim_identifier,omitempty"`
    ERPIdentifier *string    `json:"erp_identifier,omitempty"`
    LastSyncedAt  *time.Time `json:"last_synced_at,omitempty"`
}
```

### Erweiterung des Product-Structs

```go
// internal/domain/product.go — Erweiterung

type Product struct {
    // ... bestehende Felder ...

    // Hersteller-Beziehung
    ManufacturerID *uuid.UUID    `json:"manufacturer_id,omitempty" db:"manufacturer_id"`
    Manufacturer   *Manufacturer `json:"manufacturer,omitempty" db:"-"` // Lazy-loaded, nicht in products-Tabelle
}

type ProductFilter struct {
    // ... bestehende Felder ...
    ManufacturerID *uuid.UUID // Filter: Nur Produkte dieses Herstellers
}
```

**Lazy-Loading:** Das `Manufacturer`-Objekt wird **nicht** automatisch bei jedem Produkt-Query geladen. Es wird nur mitgeliefert, wenn:

1. Die Produktdetail-API aufgerufen wird (`GET /api/v1/products/:id`)
2. Explizit per Query-Parameter angefordert: `?include=manufacturer`
3. In der Produktliste über einen separaten Batch-Query (N+1-Vermeidung)

**Implementierung im Repository:**

```go
// Produkte laden, dann Hersteller in einem Batch-Query nachladen
func (r *repository) enrichWithManufacturers(ctx context.Context, products []*Product) error {
    // Alle einzigartigen ManufacturerIDs sammeln
    ids := make(map[uuid.UUID]bool)
    for _, p := range products {
        if p.ManufacturerID != nil {
            ids[*p.ManufacturerID] = true
        }
    }
    if len(ids) == 0 {
        return nil
    }

    // Ein Query für alle Hersteller
    manufacturers, err := r.manufacturerRepo.GetByIDs(ctx, mapKeys(ids))
    if err != nil {
        return err
    }

    // Zuordnen
    mfMap := make(map[uuid.UUID]*Manufacturer)
    for _, m := range manufacturers {
        mfMap[m.ID] = m
    }
    for _, p := range products {
        if p.ManufacturerID != nil {
            p.Manufacturer = mfMap[*p.ManufacturerID]
        }
    }
    return nil
}
```

---

## 4. API-Design

### Public API (Storefront)

#### Hersteller-Liste

```
GET /api/v1/manufacturers?search=bosch&country=DE&limit=20&offset=0
```

**Response:**

```json
{
    "data": [
        {
            "id": "uuid",
            "code": "bosch-rexroth",
            "name": {"de": "Bosch Rexroth", "fr": "Bosch Rexroth"},
            "slug": {"de": "bosch-rexroth", "fr": "bosch-rexroth"},
            "logo_url": "https://assets.example.com/logos/bosch-rexroth.svg",
            "country": "DE",
            "status": "active",
            "product_count": 347
        }
    ],
    "pagination": {
        "total": 1,
        "limit": 20,
        "offset": 0
    }
}
```

**Query-Parameter:**

| Parameter | Typ | Beschreibung |
|-----------|-----|--------------|
| `search` | string | Volltextsuche in Name und Code |
| `country` | string | Filter nach Herkunftsland (ISO Alpha-2) |
| `letter` | string | Filter nach Anfangsbuchstabe (`A`, `B`, ...) |
| `limit` | int | Max. Ergebnisse (Default: 20, Max: 100) |
| `offset` | int | Pagination-Offset |
| `sort` | string | Sortierung: `name` (Default), `product_count`, `sort_order` |

#### Hersteller-Detail

```
GET /api/v1/manufacturers/:id
GET /api/v1/manufacturers/by-slug/:slug    # Alternativ per Slug (für SEO-URLs)
```

**Response:**

```json
{
    "id": "uuid",
    "code": "festo",
    "name": {"de": "Festo SE & Co. KG", "fr": "Festo SE & Co. KG"},
    "slug": {"de": "festo", "fr": "festo"},
    "description": {"de": "Festo ist ein weltweit führender Anbieter von Automatisierungstechnik..."},
    "logo_url": "https://assets.example.com/logos/festo.svg",
    "website": "https://www.festo.com",
    "country": "DE",
    "status": "active",
    "meta_title": {"de": "Festo Produkte — Pneumatik & Automatisierung"},
    "meta_description": {"de": "Festo Produkte online bestellen: Pneumatikzylinder, Ventile, Antriebe und Zubehör."},
    "product_count": 1243,
    "created_at": "2026-01-15T10:30:00Z",
    "updated_at": "2026-02-10T14:22:00Z"
}
```

#### Produkte eines Herstellers

```
GET /api/v1/manufacturers/:id/products?limit=20&offset=0&category_id=uuid&status=active
```

Nutzt intern denselben `ProductFilter` wie die Products-API, erweitert um `manufacturer_id`.

#### Filter in der Products-API

```
GET /api/v1/products?manufacturer_id=uuid&category_id=uuid&search=zylinder
```

Der bestehende `ProductFilter` wird um `manufacturer_id` erweitert (siehe Datenmodell-Erweiterung oben).

### Admin API

```
POST   /api/v1/admin/manufacturers           # Erstellen
GET    /api/v1/admin/manufacturers            # Liste (inkl. inactive, inkl. Kontaktdaten)
GET    /api/v1/admin/manufacturers/:id        # Detail (inkl. Kontaktdaten, PIM/ERP-IDs)
PUT    /api/v1/admin/manufacturers/:id        # Aktualisieren
DELETE /api/v1/admin/manufacturers/:id        # Soft-Delete
POST   /api/v1/admin/manufacturers/:id/restore # Wiederherstellen
```

**Create-Request:**

```json
{
    "code": "phoenix-contact",
    "name": {"de": "Phoenix Contact", "fr": "Phoenix Contact"},
    "description": {"de": "Hersteller von industrieller Verbindungstechnik..."},
    "logo_url": "https://assets.example.com/logos/phoenix-contact.svg",
    "website": "https://www.phoenixcontact.com",
    "country": "DE",
    "status": "active",
    "contact_email": "vertrieb@phoenixcontact.de",
    "contact_person": "Max Mustermann",
    "erp_identifier": "LIEFERANT-00472"
}
```

**Slug wird automatisch generiert** aus dem Namen. Bei Konflikten wird ein Suffix angehängt (`phoenix-contact-2`).

### Events

| Event | Trigger | Konsumenten |
|-------|---------|-------------|
| `manufacturer.created` | Neuer Hersteller | Meilisearch-Indexer |
| `manufacturer.updated` | Hersteller geändert | Meilisearch-Indexer, Cache-Invalidation |
| `manufacturer.deleted` | Soft-Delete | Meilisearch-Indexer (aus Index entfernen) |
| `manufacturer.restored` | Wiederherstellung | Meilisearch-Indexer |

---

## 5. Service-Architektur

Der Hersteller wird als **Domain innerhalb des Catalog-Service** implementiert — kein eigener Microservice. Begründung: Hersteller sind eng mit dem Produktkatalog verknüpft und haben keine unabhängigen Skalierungsanforderungen.

### Verzeichnisstruktur

```
services/catalog/
  internal/
    manufacturer/
      handler.go          # HTTP Handler (Public + Admin)
      service.go           # Business Logic
      repository.go        # DB-Zugriff
      models.go            # Structs, DTOs, Errors
      errors.go            # Domain-spezifische Fehler
      manufacturer_test.go # Unit Tests
    product/
      ...                  # Bestehend — erweitert um ManufacturerID
    domain/
      product.go           # Erweitert um ManufacturerID + Manufacturer
  migrations/
    000003_create_manufacturers.up.sql
    000003_create_manufacturers.down.sql
```

### Error-Definitionen

```go
// internal/manufacturer/errors.go
package manufacturer

import "errors"

var (
    ErrManufacturerNotFound    = errors.New("manufacturer not found")
    ErrCodeAlreadyExists       = errors.New("manufacturer code already exists")
    ErrSlugAlreadyExists       = errors.New("manufacturer slug already exists")
    ErrManufacturerHasProducts = errors.New("manufacturer has active products")
    ErrValidation              = errors.New("validation error")
)
```

---

## 6. Frontend

### Hersteller-Übersichtsseite

**Route:** `/hersteller` (bzw. `/fabricants`, `/manufacturers` — je nach Sprache)

**Darstellung:**

```
┌─────────────────────────────────────────────────────────┐
│ Hersteller                                              │
│                                                         │
│ [Suche: ___________________________]                    │
│                                                         │
│ A B C D E F G H I J K L M N O P Q R S T U V W X Y Z   │
│                                                         │
│ ── A ──────────────────────────────────────────────     │
│ ┌──────────┐ ┌──────────┐ ┌──────────┐                 │
│ │ [Logo]   │ │ [Logo]   │ │ [Logo]   │                 │
│ │ 3M       │ │ Ansell   │ │ APC      │                 │
│ │ 156 Prod.│ │ 89 Prod. │ │ 42 Prod. │                 │
│ └──────────┘ └──────────┘ └──────────┘                 │
│                                                         │
│ ── B ──────────────────────────────────────────────     │
│ ┌──────────┐ ┌──────────┐ ┌──────────┐                 │
│ │ [Logo]   │ │ [Logo]   │ │ [Logo]   │                 │
│ │ Bosch    │ │ Bosch    │ │ Beiers-  │                 │
│ │ Rexroth  │ │ 478 Prod.│ │ dorf     │                 │
│ │ 347 Prod.│ │          │ │ 23 Prod. │                 │
│ └──────────┘ └──────────┘ └──────────┘                 │
│                                                         │
│ ...                                                     │
└─────────────────────────────────────────────────────────┘
```

**Anforderungen:**
- Alphabetische Gruppierung mit Buchstaben-Schnellnavigation (Ankerlinks)
- Hersteller-Kacheln: Logo, Name, Produktanzahl
- Suchfeld mit Live-Filterung (clientseitig, da max. ~600 Hersteller)
- Klick auf Kachel → Hersteller-Detailseite
- Nur Hersteller mit mindestens 1 aktiven Produkt anzeigen

### Hersteller-Detailseite

**Route:** `/hersteller/bosch-rexroth` (Slug-basiert, SEO-freundlich)

```
┌─────────────────────────────────────────────────────────┐
│ [Logo]  Bosch Rexroth                                   │
│                                                         │
│ Bosch Rexroth ist ein weltweit führender Anbieter       │
│ von Antriebs- und Steuerungstechnologien...             │
│                                                         │
│ 🌐 www.bosch-rexroth.com    📍 Deutschland             │
│                                                         │
│ ── 347 Produkte ───────────────────────────────────     │
│                                                         │
│ [Kategorie-Filter ▾]  [Sortierung ▾]  [Suche: ___]     │
│                                                         │
│ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐       │
│ │ [Bild]  │ │ [Bild]  │ │ [Bild]  │ │ [Bild]  │       │
│ │ Hydraul.│ │ Linear- │ │ Ventil  │ │ Pumpe   │       │
│ │ zylinder│ │ führung │ │ 4WE6    │ │ A10VSO  │       │
│ │ CHF 485 │ │ CHF 120 │ │ CHF 340 │ │ CHF 2890│       │
│ └─────────┘ └─────────┘ └─────────┘ └─────────┘       │
│                                                         │
│ [1] [2] [3] ... [18]                                    │
└─────────────────────────────────────────────────────────┘
```

**Anforderungen:**
- Header: Logo, Name, Beschreibung, Website-Link, Herkunftsland
- Produktliste des Herstellers (paginiert, filterbar, sortierbar)
- Kategorie-Filter (nur Kategorien die Produkte dieses Herstellers enthalten)
- SEO: Eigene Meta-Tags, Canonical URL, strukturierte Daten (Schema.org `Organization`)

### Hersteller-Filter in der Produkt-Sidebar

```
┌─────────────────────┐
│ Hersteller          │
│ [Suche: ______]     │
│ ☑ Bosch Rexroth (347) │
│ ☐ Festo (1243)      │
│ ☐ SKF (89)          │
│ ☐ Parker (156)      │
│ [Alle anzeigen ▾]   │
└─────────────────────┘
```

**Anforderungen:**
- Checkbox-Filter in der Produkt-Sidebar
- Produktanzahl pro Hersteller (basierend auf aktuellem Filter-Kontext)
- Suche innerhalb der Hersteller-Liste (bei >10 Herstellern)
- Multi-Select: Mehrere Hersteller gleichzeitig filterbar
- Top-5 anzeigen, Rest unter «Alle anzeigen» einklappbar

### Hersteller-Anzeige auf der Produktdetailseite

```
┌─────────────────────────────────────────┐
│ Hydraulikzylinder CDT3                  │
│                                         │
│ Hersteller: [Logo] Bosch Rexroth →     │
│                                         │
│ ...                                     │
└─────────────────────────────────────────┘
```

**Anforderungen:**
- Hersteller-Name + Logo (klein) in den Produktdetails
- Link zur Hersteller-Detailseite
- Position: Im oberen Bereich der Produktdetails, unter dem Produktnamen

---

## 7. Suche (Meilisearch)

### Hersteller als Facette in der Produktsuche

Im Meilisearch-Produktindex wird der Hersteller als filterbares und facettierbares Attribut indexiert:

```json
{
    "id": "product-uuid",
    "name": "Hydraulikzylinder CDT3",
    "manufacturer_id": "manufacturer-uuid",
    "manufacturer_name": "Bosch Rexroth",
    "manufacturer_code": "bosch-rexroth",
    "searchable_text": "Hydraulikzylinder CDT3 Bosch Rexroth ...",
    "..."
}
```

**Meilisearch-Konfiguration:**

```json
{
    "filterableAttributes": ["manufacturer_id", "manufacturer_code", "category_ids", "status"],
    "facets": ["manufacturer_code"],
    "searchableAttributes": ["name", "sku", "manufacturer_name", "searchable_text"]
}
```

**Facetten-Response:**

```json
{
    "hits": [...],
    "facetDistribution": {
        "manufacturer_code": {
            "bosch-rexroth": 347,
            "festo": 1243,
            "skf": 89
        }
    }
}
```

### Hersteller selbst durchsuchbar

Eigener Meilisearch-Index `manufacturers`:

```json
{
    "id": "manufacturer-uuid",
    "code": "bosch-rexroth",
    "name_de": "Bosch Rexroth",
    "name_fr": "Bosch Rexroth",
    "description_de": "Weltweit führender Anbieter von Antriebs- und Steuerungstechnologien...",
    "country": "DE",
    "product_count": 347
}
```

**Anwendungsfälle:**
- Globale Suchleiste: «bosch» findet sowohl Produkte als auch den Hersteller «Bosch Rexroth»
- Autosuggest: Bei Eingabe «fes» erscheint «Festo» als Hersteller-Vorschlag (neben Produkten)

---

## 8. ERP/SAP-Integration

### Hersteller-Mapping

In SAP gibt es kein einheitliches «Hersteller»-Objekt. Je nach SAP-Modul und Kundenkonfiguration kommen verschiedene Strukturen in Frage:

| SAP-Objekt | Einsatz | Mapping |
|------------|---------|---------|
| **Lieferantenstamm (Kreditor)** | Wenn Hersteller = Lieferant | `erp_identifier` → SAP Kreditor-Nr. |
| **Hersteller-Stammdaten (MM)** | SAP Materialstamm Feld MARA-MFRNR | `erp_identifier` → SAP Hersteller-Nr. |
| **Geschäftspartner (BP)** | S/4HANA Business Partner | `erp_identifier` → BP-Nummer mit Rolle «Hersteller» |
| **Klassifizierungsmerkmal** | Hersteller als Merkmal am Material | `code` → Merkmalswert |

**Sync-Richtung:**

| Richtung | Beschreibung |
|----------|-------------|
| **SAP → Gondolia** | Hersteller-Stammdaten werden aus SAP importiert (initiale Befüllung, laufender Abgleich) |
| **Gondolia → SAP** | Bei Produktbestellungen wird die Hersteller-Referenz mitgegeben (SAP-Materialnummer enthält bereits Hersteller-Zuordnung) |

**IDoc/RFC-Mapping:**

```
SAP Hersteller-Stammdaten → Gondolia:
  MFRNR (Herstellernummer)     → erp_identifier
  NAME1 (Name)                 → name["de"]
  LAND1 (Land)                 → country
  STRAS (Strasse)              → (nicht gemappt, optional als Kontaktdaten)
```

### PIM-Sync (Akeneo)

Akeneo hat ein natives Konzept für «Reference Entities», das für Hersteller genutzt wird:

```
Akeneo Reference Entity "manufacturer" → Gondolia:
  code                                  → code
  label (per locale)                    → name
  description (per locale)              → description
  image                                 → logo_url (nach Asset-Upload)
  website (Custom Attribute)            → website
  country (Custom Attribute)            → country
```

**Sync-Ablauf:**

1. Akeneo-Webhook oder geplanter Job löst Sync aus
2. Catalog-Service empfängt Hersteller-Daten
3. Upsert per `pim_identifier`: Existiert der Hersteller? → Update. Neu? → Create.
4. `last_synced_at` wird aktualisiert
5. Event `manufacturer.updated` oder `manufacturer.created` wird publiziert
6. Meilisearch-Index wird aktualisiert

**Import-Logik:**

```go
func (s *service) SyncFromPIM(ctx context.Context, tenantID uuid.UUID, data PIMManufacturerData) (*Manufacturer, error) {
    existing, err := s.repo.GetByPIMIdentifier(ctx, tenantID, data.PIMIdentifier)
    if err != nil && !errors.Is(err, ErrManufacturerNotFound) {
        return nil, err
    }

    if existing != nil {
        // Update
        existing.Name = data.Name
        existing.Description = data.Description
        existing.LogoURL = data.LogoURL
        existing.Website = data.Website
        existing.Country = data.Country
        existing.LastSyncedAt = timePtr(time.Now())
        existing.UpdatedAt = time.Now()
        return existing, s.repo.Update(ctx, existing)
    }

    // Create
    manufacturer := &Manufacturer{
        ID:            uuid.New(),
        TenantID:      tenantID,
        Code:          data.Code,
        Name:          data.Name,
        Slug:          generateSlug(data.Name),
        Description:   data.Description,
        LogoURL:       data.LogoURL,
        Website:       data.Website,
        Country:       data.Country,
        Status:        ManufacturerStatusActive,
        PIMIdentifier: &data.PIMIdentifier,
        LastSyncedAt:  timePtr(time.Now()),
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
    }
    return manufacturer, s.repo.Create(ctx, manufacturer)
}
```

---

## 9. Priorisierung

| Phase | Aufgabe | Aufwand |
|-------|---------|---------|
| **Phase 1** | DB-Migration, Domain-Struct, Repository, Service, CRUD-Endpoints (Admin + Public) | 3–4 Tage |
| **Phase 2** | Product-Erweiterung (`manufacturer_id`), Lazy-Loading, Filter in Products-API | 2–3 Tage |
| **Phase 3** | Meilisearch-Integration (Facette + eigener Index) | 1–2 Tage |
| **Phase 4** | Frontend: Übersichtsseite, Detailseite, Filter-Sidebar, Produktdetail-Anzeige | 3–5 Tage |
| **Phase 5** | PIM-Sync (Akeneo), SAP-Mapping | 2–3 Tage |

**Gesamtaufwand:** ~2–3 Wochen

---

## Offene Fragen

1. **Logo-Upload:** Werden Hersteller-Logos über den bestehenden Asset-Service verwaltet oder direkt als URL gespeichert? Gibt es Anforderungen an Bildformate (SVG bevorzugt)?
2. **Akeneo-Struktur:** Ist in Akeneo bereits eine Reference Entity für Hersteller eingerichtet? Welche Attribute sind vorhanden?
3. **SAP-Hersteller:** Welches SAP-Objekt wird aktuell für Hersteller-Stammdaten verwendet (Kreditor, MFRNR, Business Partner)?
4. **Multi-Select:** Sollen Kunden in der Suche nach mehreren Herstellern gleichzeitig filtern können? (Empfehlung: ja)
5. **Hersteller ohne Produkte:** Sollen Hersteller auf der Übersichtsseite angezeigt werden, auch wenn sie (noch) keine aktiven Produkte haben? (Empfehlung: nein)
6. **Hersteller-Hierarchie:** Gibt es Konzernstrukturen die abgebildet werden müssen? (z.B. Bosch → Bosch Rexroth, Bosch Power Tools) (Empfehlung: Nein, flache Liste reicht. Bei Bedarf später als `parent_id` erweiterbar.)
